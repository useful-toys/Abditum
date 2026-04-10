# Arquitetura — Pacote `internal/storage`

## Responsabilidade

O pacote `storage` é a única camada do sistema responsável por persistir e recuperar cofres cifrados em disco. Ele traduz entre a representação em memória (`vault.Cofre`) e o formato binário `.abditum`, aplicando derivação de chave, cifração autenticada, escrita atômica e detecção de modificação externa. Nenhuma outra camada conhece o formato do arquivo.

## Formato Binário `.abditum`

Ver [especificacao-arquivo-cofre.md §1](especificacao-arquivo-cofre.md) para especificação completa do layout do arquivo (cabeçalho de 49 bytes fixos, magic, versão, salt, nonce, payload), o uso de AAD (§2) e detalhes criptográficos.

Detalhes específicos de implementação no pacote `storage`:

- Magic e versão são armazenados como constantes no código (`MagicSize`, `HeaderSize`, etc.).
- Validação de magic (`ErrInvalidMagic`) ocorre no parse inicial sem tentativa de derivação de chave.
- O byte de versão determina qual profile Argon2id será ativado via `ProfileForVersion`.

## Derivação de Chave e Perfis de Versão

Ver [especificacao-arquivo-cofre.md §3](especificacao-arquivo-cofre.md) para conceitos de Argon2id, salt, nonce e ciclo de vida, incluindo política de versionamento.

Implementação em `storage`:

Os parâmetros Argon2id **não são armazenados no arquivo**. O número de versão no cabeçalho é a única informação presente; os parâmetros correspondentes são obtidos internamente via `ProfileForVersion`, que consulta um mapa estático `formatProfiles`. Cada versão define um `FormatProfile` com os campos `Time`, `Memory`, `Threads` e `KeyLen`.

Quando a versão lida do arquivo não consta no mapa, `ProfileForVersion` retorna `ErrVersionTooNew`. A chave derivada é wiped via `crypto.Wipe` imediatamente após o uso.

## Protocolo de Escrita Atômica

### Criação de novo cofre — `SaveNew`

Quando o cofre ainda não existe em disco, a escrita é feita diretamente no caminho destino. Não há arquivo preexistente a proteger, portanto o protocolo de renomeação não é necessário. O arquivo é criado com permissões `0600` (somente dono pode ler e escrever).

### Salvamento de cofre existente — `Save`

Para preservar o arquivo original caso o processo seja interrompido, o salvamento segue um protocolo de cinco etapas com cadeia de backups:

1. Escrever o novo conteúdo em `<vault>.tmp`.
2. Se `<vault>.bak` existir, renomeá-lo para `<vault>.bak2` (sobrescreve qualquer `.bak2` anterior).
3. Renomear `<vault>` para `<vault>.bak`.
4. Renomear `<vault>.tmp` para `<vault>` de forma atômica.
5. Em caso de sucesso: remover `<vault>.bak2` (best-effort). Em caso de falha após o passo 1: remover `<vault>.tmp` (best-effort) e tentar restaurar `<vault>.bak`.

O invariante preservado é que o arquivo `<vault>` nunca fica em estado parcialmente escrito. O pior cenário em caso de interrupção é o arquivo `.tmp` órfão, que é resolvido pelo `RecoverOrphans`.

### Atomicidade por Sistema Operacional

A renomeação da etapa 4 requer semântica atômica diferente por plataforma, implementada com build tags:

- **Unix/Linux/macOS:** `os.Rename` mapeia para `rename(2)` do POSIX, que é atômico quando origem e destino estão no mesmo filesystem. Por isso o `.tmp` é criado no **mesmo diretório** que o vault (nunca em `os.TempDir()`).
- **Windows:** `os.Rename` não é atômico ao substituir arquivo existente, pois internamente faz `DeleteFile` + `MoveFile`. O código usa `MoveFileEx` com a flag `MOVEFILE_REPLACE_EXISTING`, que realiza a substituição como operação única no MFT do NTFS.

## Recuperação de Órfãos — `RecoverOrphans`

`RecoverOrphans` deve ser chamado uma vez na inicialização do programa, antes de qualquer `Load` ou `Save`. Ele verifica se existe um arquivo `<vault>.tmp` e, caso exista, o remove.

Um `.tmp` órfão indica que o processo anterior foi interrompido após a escrita do arquivo temporário mas antes da renomeação atômica. Como o vault original não foi tocado nesse ponto, o `.tmp` pode ser descartado com segurança.

`RecoverOrphans` **não** tenta restaurar o vault automaticamente a partir de `.bak` em caso de corrupção. Essa decisão é deliberada: restauração automática poderia provocar perda de dados silenciosa em cenários de corrupção não relacionados ao processo de salvamento. O erro é retornado ao chamador para que o usuário decida.

## Detecção de Modificação Externa — `DetectExternalChange`

O pacote fornece mecanismo para detectar se o arquivo de vault foi modificado por processo externo (sincronizadores de nuvem, editors, scripts) entre um Load e a próxima interação do usuário.

A estratégia é em dois níveis:

1. **Fast path — tamanho:** compara o tamanho atual do arquivo com o registrado em `FileMetadata`. Uma divergência de tamanho sinaliza mudança imediatamente, com custo de um único `os.Stat` (sem leitura de conteúdo).
2. **Slow path — SHA-256:** se o tamanho coincide, lê o arquivo e compara o hash SHA-256 completo. Detecta substituições *in-place* que preservam o tamanho.

**Mtime é deliberadamente ignorado.** Serviços de sincronização em nuvem (Dropbox, iCloud, Google Drive) e sistemas de arquivos de rede podem modificar o mtime por razões que não correspondem a mudanças de conteúdo, gerando falsos positivos. A decisão de usar tamanho + hash está documentada em `arquitetura.md §7` e `historico-decisoes.md §17`.

`FileMetadata` é um struct com dois campos: `Size int64` e `Hash [32]byte`. Ele é capturado após cada `Load` ou `Save` e armazenado no `FileRepository` como baseline para comparações futuras.

## `FileRepository` — Implementação da Interface

`FileRepository` implementa a interface `vault.RepositorioCofre` (definida no pacote `vault` com os métodos `Salvar` e `Carregar`). A conformidade é verificada em tempo de compilação via `var _ vault.RepositorioCofre = (*FileRepository)(nil)`.

### Estado interno

O repositório mantém cinco campos privados:

- `path` — caminho absoluto do arquivo vault.
- `password` — referência à slice de bytes da senha master. O chamador é responsável por não wipe-ar a slice enquanto o repositório estiver em uso.
- `salt` — salt de 32 bytes extraído do cabeçalho. Preservado entre salvamentos para que a chave derivada permaneça estável.
- `isNew` — flag booleano que distingue o primeiro `Salvar` (usa `SaveNew`) de todos os seguintes (usa `Save` com protocolo atômico).
- `metadata` — snapshot de `FileMetadata` atualizado após cada `Salvar` ou `Carregar`.

### Ciclo de vida

Dois construtores cobrem os dois cenários de uso:

- `NewFileRepositoryForCreate` — para cofres que ainda não existem em disco. O campo `salt` começa nulo; na primeira chamada a `Salvar`, o `salt` é extraído do arquivo recém-criado e armazenado. A flag `isNew` é revertida para `false` após o primeiro salvamento bem-sucedido.
- `NewFileRepository` — para cofres já existentes. Recebe o salt lido do cabeçalho e o `FileMetadata` capturado pelo `Load` inicial, iniciando o repositório já pronto para salvamentos subsequentes.

O método `UpdatePassword` substitui a senha armazenada sem alterar o salt, permitindo que a TUI troque a senha sem recriar o repositório. A geração de novo salt em cenários de troca de senha é responsabilidade de camada superior.

## Erros Sentinela

O pacote declara três erros verificáveis via `errors.Is`:

- `ErrInvalidMagic` — arquivo muito curto ou magic bytes incorretos. Indica que o arquivo não é um vault `.abditum` válido.
- `ErrVersionTooNew` — versão do formato não consta no mapa de perfis. Indica que o binário precisa ser atualizado para abrir este arquivo.
- `ErrCorrupted` — decifração bem-sucedida, mas o JSON decifrado falha na validação estrutural (JSON inválido, `pasta_geral` ausente ou `pasta_geral.nome != "Geral"`). Indica corrupção de conteúdo além da proteção criptográfica.

Para senha errada ou adulteração de conteúdo, o AEAD falha com `crypto.ErrAuthFailed`, que é propagado diretamente sem ser redeclarado neste pacote.

## Segurança Operacional

- Arquivos são criados com permissão `0600`, impedindo leitura por outros usuários do sistema.
- JSON serializado do cofre é mantido em memória apenas o tempo necessário; bytes são wiped via `crypto.Wipe` com `defer` em todos os caminhos.
- A chave derivada (AES-256) também é wiped imediatamente após uso.
- O nonce é novo a cada escrita, eliminando reutilização de `(chave, nonce)` mesmo com a mesma senha.

## Convenções do Pacote

- Métodos da interface pública usam nomes em português (`Salvar`, `Carregar`), seguindo a convenção geral do projeto.
- Erros são sempre envolvidos com `fmt.Errorf("storage.<Função>: <contexto>: %w", err)`, garantindo contexto de origem na cadeia de erros.
- A interface `vault.RepositorioCofre` está definida no pacote `vault`, não em `storage`, seguindo o princípio de que dependências apontam de `storage` para `vault`, não o inverso.
- Constantes de layout do header (`MagicSize`, `SaltOffset`, `NonceOffset`, `HeaderSize`) são exportadas para uso nos testes e para verificação de integridade por outros componentes.
