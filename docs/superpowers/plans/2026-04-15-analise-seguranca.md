# Relatório de Análise de Segurança — Abditum

**Data:** 2026-04-15
**Escopo:** Análise estática completa do código-fonte (`internal/crypto`, `internal/storage`, `internal/vault`, `cmd/abditum`, `internal/tui`, `go.mod`)
**Metodologia:** Revisão manual por domínio, executada em 4 agentes paralelos independentes
**Nenhuma correção foi aplicada** — este documento é apenas o resultado da análise

---

## Sumário Executivo

O Abditum é um gerenciador de senhas TUI offline (sem rede), escrito em Go. A criptografia central (AES-256-GCM + Argon2id) é sólida, mas a **gestão de memória sensível** tem falhas que podem expor a senha mestre e segredos armazenados em dumps de memória, arquivos swap e core dumps. Há também dois findings críticos de **design** que tornam ineficazes os mecanismos de proteção existentes.

| Severidade | Qtd |
|------------|-----|
| CRÍTICO    | 5   |
| ALTO       | 10  |
| MÉDIO      | 9   |
| BAIXO      | 5   |
| INFORMATIVO | 12  |

---

## CRÍTICO

### C-01 — Clipboard não limpo em crash/SIGKILL

**Arquivo:** `cmd/abditum/main.go:36`

```go
defer clipboard.WriteAll("")
```

O `defer` só executa em saída normal do processo. SIGKILL, crash por panic não recuperado, ou `os.Exit` direto deixam a senha copiada no clipboard indefinidamente.

**Risco:** Senha mestre ou segredo copiado permanece acessível a qualquer processo/usuário após encerramento inesperado da aplicação.

**Correção sugerida:** Registrar handler de sinal (`SIGTERM`, `SIGINT`) que escreva `""` no clipboard antes de `os.Exit`. Para SIGKILL não há solução completa em userspace — documentar a limitação.

---

### C-02 — Feature de auto-clear do clipboard nunca conectada

**Arquivos:** `internal/vault/entities.go:53`, `internal/vault/manager.go:204`

O campo `tempoLimparAreaTransferenciaSegundos` existe, é validado e persistido — mas **nenhuma goroutine ou ticker** o lê para agendar a limpeza do clipboard.

**Risco:** O usuário confia que o clipboard é limpo automaticamente após N segundos; isso nunca acontece. A feature está silenciosamente quebrada.

**Correção sugerida:** Implementar a goroutine de limpeza no `Manager` ou no componente TUI que executa a cópia, usando `time.AfterFunc`.

---

### CRIT-01 — `Manager.senha` nunca populado; Lock() não zera a senha real

**Arquivo:** `internal/vault/manager.go:14,23-28`

```go
type Manager struct {
    senha []byte  // declarado mas nunca atribuído
    ...
}
```

A senha real vive em `FileRepository.password` (camada de storage). `Manager.Lock()` zera `Manager.senha` (que é sempre `nil`/vazio) sem tocar `FileRepository.password`.

**Risco:** Após `Lock()`, a senha mestre ainda está viva em memória no `FileRepository`. Um dump de memória após bloqueio do cofre expõe a senha.

**Correção sugerida:** `Manager.Lock()` deve chamar explicitamente `repository.WipePassword()` (método a criar) que zere e descarte o slice interno.

---

### CRIT-02 — `FileRepository.password` sem ciclo de wipe

**Arquivo:** `internal/storage/repository.go:29,159`

```go
type FileRepository struct {
    password []byte  // nunca zerado
    ...
}
```

`UpdatePassword` substitui o slice sem zerar o anterior. Não existe método `Close()`/`Wipe()`. O GC pode mover o slice para outra região de memória antes de coletar, multiplicando cópias da senha.

**Risco:** A senha mestre pode permanecer em memória por tempo indeterminado após mudança de senha ou encerramento.

**Correção sugerida:** Usar `crypto.Wipe()` antes de toda substituição de `password`. Implementar `FileRepository.Close()` que zere o campo. Considerar `mlock` para impedir paginação para swap.

---

### CRIT-03 — Serialização cria cópias não zeráveis de todos os segredos

**Arquivo:** `internal/vault/serialization.go:139-148`

```go
"valor": string(c.valor),  // converte []byte → string
```

Strings em Go são imutáveis e não podem ser zeradas. Durante serialização (salvar cofre), todos os valores sensíveis são copiados para strings no heap sem controle de ciclo de vida.

**Risco:** Cada `Save()` deixa rastros não zeráveis de todos os segredos em memória até o GC coletar — sem garantia de quando ou se serão sobrescritos.

**Correção sugerida:** Serializar diretamente para `[]byte` (ex: `json.Marshal` com campos `json.RawMessage`) mantendo o buffer final zerável. Alternativamente, usar `encoding/json` com `RawMessage` para evitar a conversão para `string`.

---

## ALTO

### A-01 — `Wipe()` pode ser otimizada pelo compilador

**Arquivo:** `internal/crypto/memory.go:13-18`

O loop de zeroing não tem garantia de não ser eliminado pelo compilador Go (dead store elimination). `runtime.KeepAlive` ajuda, mas não é garantia formal para otimizações agressivas.

**Correção sugerida:** Usar `runtime/internal/sys.Clear` (Go 1.21+) ou `crypto/internal/alias.AnyOverlap` + assembly, ou ao menos `atomic.StoreUint8` para forçar efeito colateral observável.

---

### A-02 — Salt mínimo aceito é 1 byte

**Arquivo:** `internal/crypto/kdf.go:75-77`

A validação aceita qualquer salt com `len >= 1`. RFC 9106 recomenda mínimo de 16 bytes para Argon2.

**Correção sugerida:** Validar `len(salt) >= 16` e retornar erro descritivo.

---

### A-03 — `EvaluatePasswordStrength` usa `len()` em bytes, não runes

**Arquivo:** `internal/crypto/password.go:32,40`

`len("é")` retorna 2 (bytes UTF-8), não 1 (caractere). Senhas com caracteres Unicode multibyte são classificadas com comprimento maior do que o real em caracteres.

**Correção sugerida:** Usar `utf8.RuneCountInString()` e `[]rune(password)` para iteração.

---

### A-04 — `SecureAllocate` existe mas não é usado no caminho crítico

**Arquivo:** `internal/crypto/memory.go`

A função `SecureAllocate` (com `mlock`) está implementada mas não é usada para alocar a chave derivada pelo KDF nem para o nonce do AEAD.

**Risco:** A chave de criptografia pode ser paginada para swap.

**Correção sugerida:** Usar `SecureAllocate` para o slice da chave derivada em `kdf.go`.

---

### A-05 — Arquivos `.bak` e `.bak2` retidos permanentemente

**Arquivo:** `internal/storage/storage.go:163`

O código comenta que `.bak2` deve ser removido após operação bem-sucedida, mas a remoção nunca foi implementada. Backups acumulam indefinidamente com conteúdo cifrado do cofre em versões anteriores.

**Risco:** Proliferação de arquivos de backup com dados sensíveis (embora cifrados).

**Correção sugerida:** Implementar remoção do `.bak2` após `Save()` bem-sucedido.

---

### A-06 — `fileData` e `data` não zerados após uso

**Arquivo:** `internal/storage/storage.go:69,127,179`

Os buffers `fileData` (header + ciphertext) e `data` (arquivo completo lido do disco) não são zerados após uso. Ficam no heap até o GC.

**Correção sugerida:** `defer crypto.Wipe(fileData)` e `defer crypto.Wipe(data)` imediatamente após alocação.

---

### A-07 — Leitura dupla do arquivo apenas para extrair salt

**Arquivo:** `internal/storage/repository.go:96,131`

O arquivo do cofre é lido duas vezes: uma para extrair o salt e outra para descriptografar. Além de ineficiente, duplica o tempo em que os dados ficam em buffers não zerados.

**Correção sugerida:** Extrair salt e carregar o arquivo em uma única leitura.

---

### A-08 — `MoveFileEx` não é atômico com destino preexistente no Windows

**Arquivo:** `internal/storage/atomic_rename_windows.go:25`

`MoveFileEx` com destino existente não é atômico no Windows — há uma janela entre deletar o destino e mover o arquivo. `ReplaceFileW` (WinAPI) é a operação correta para rename atômico com destino preexistente.

**Risco:** Corrupção do cofre se o processo for interrompido durante o save no Windows.

**Correção sugerida:** Substituir por `ReplaceFileW` via `syscall`.

---

### A-09 — `zerarValoresSensiveis()` ignora campos comuns

**Arquivo:** `internal/vault/entities.go:237-246`

`zerarValoresSensiveis()` só zera campos do tipo `TipoCampoSensivel`. URLs, usernames e notas (tipo `TipoCampoComum`) não são zerados no `Lock()`.

**Risco:** Dados potencialmente sensíveis (usernames, URLs) permanecem em memória após bloqueio.

**Correção sugerida:** Zerar todos os campos na função, independente do tipo.

---

### A-10 — Nil pointer panic em `DescerPastaNaPosicao`

**Arquivo:** `internal/vault/manager.go:598-609`

Se `pasta.pai == nil` (PastaGeral, raiz da árvore), a função não verifica antes de acessar `pasta.pai.subpastas`, causando panic.

**Risco:** Crash da aplicação por operação do usuário em condição de borda previsível.

**Correção sugerida:** Verificar `if pasta.pai == nil { return erro }` antes de acessar o pai.

---

## MÉDIO

### M-01 — Senha mestre sem `mlock` no `FileRepository`

**Arquivo:** `internal/storage/repository.go:27-29`

O campo `password []byte` do `FileRepository` não é alocado com `mlock`, podendo ser paginado para swap e persistir em disco.

---

### M-02 — Sem `prctl(PR_SET_DUMPABLE, 0)`

Não há chamada para desabilitar core dumps na inicialização. Um crash gera um core dump com toda a memória do processo exposta, incluindo a senha mestre e todos os segredos descriptografados.

**Correção sugerida:** Chamar `prctl(PR_SET_DUMPABLE, 0)` no início de `main()` em Linux. No Windows, `SetProcessMitigationPolicy`.

---

### M-03 — Erro de restauração de emergência silenciado

**Arquivo:** `internal/storage/storage.go:158`

```go
//nolint:errcheck
```

O erro da restauração de emergência do backup é silenciado. Se a restauração falhar, o usuário perde dados sem nenhum aviso.

**Correção sugerida:** Logar ou propagar o erro de restauração — nunca silenciar.

---

### M-04 — TOCTOU entre `Stat` e `ReadFile` em `DetectExternalChange`

**Arquivo:** `internal/storage/detect.go:29-48`

A função verifica o arquivo com `Stat` e depois lê com `ReadFile`. Entre as duas operações, o arquivo pode ser substituído por um atacante com acesso local (symlink attack / race condition).

**Correção sugerida:** Abrir o arquivo com `os.Open`, obter `Stat` do file descriptor aberto, e ler do mesmo fd.

---

### M-05 — Sem validação de `vaultPath` (path traversal)

**Arquivo:** `internal/storage/storage.go:28,93,178`

O caminho do cofre não é validado nem normalizado. Um valor com `../` poderia apontar para arquivo fora do diretório esperado.

**Correção sugerida:** Usar `filepath.Clean` e verificar que o path resultante está dentro do diretório de dados da aplicação.

---

### M-06 — Configurações desserializadas sem validação

**Arquivo:** `internal/vault/serialization.go:175`

Valores como `tempoLimparAreaTransferenciaSegundos` são carregados do arquivo sem verificar se são positivos ou dentro de um intervalo razoável.

**Correção sugerida:** Chamar a função de validação existente em `entities.go` após deserializar as configurações.

---

### M-07 — Roundtrip perde acento em `"Observacao"`

**Arquivo:** `internal/vault/serialization.go:245-249`

O campo é serializado como `"Observacao"` (sem acento) mas o nome canônico é `"Observação"`. Após um ciclo salvar/carregar, o nome do campo pode mudar dependendo do caminho de código.

---

### M-08 — Falha de `mlock` silenciosa

**Arquivo:** `internal/crypto/memory.go:47-58`

Quando `mlock` falha (ex: limite de ulimit excedido), o erro é ignorado sem aviso ao usuário. A aplicação continua como se a proteção estivesse ativa.

**Correção sugerida:** Logar um aviso não-fatal ao usuário informando que a proteção de memória não está disponível.

---

### M-09 — Sem limite de tamanho em valores de campos

**Arquivo:** `internal/vault/entities.go:371-378`

Não há validação de tamanho máximo para valores de campos. Um valor muito grande pode causar uso excessivo de memória durante criptografia/serialização.

---

## BAIXO

### B-01 — `go.mod` declara `go 1.26.1` (versão futura)

**Arquivo:** `go.mod:3`

Go 1.26.1 não existe. Isso pode causar erros em ferramentas que validam a versão do módulo.

**Correção sugerida:** Usar a versão Go mais recente estável (ex: `go 1.22`).

---

### B-02 — `charmbracelet/ultraviolet` referenciado por commit hash

**Arquivo:** `go.mod:17`

Dependência sem tag semântica, referenciada por hash de commit. Dificulta auditoria e atualização de segurança.

---

### B-03 — `golang.org/x/exp` desatualizado (18 meses)

**Arquivo:** `go.mod:31`

Pacote experimental com 18 meses de defasagem. Pode conter bugs corrigidos em versões mais recentes.

---

### B-04 — `aead.go` exige nonce externo

**Arquivo:** `internal/crypto/aead.go:171`

`SealWithAAD` exige que o chamador forneça o nonce, colocando a responsabilidade de geração segura fora do pacote crypto.

**Risco baixo:** A geração atual parece correta, mas o design convida a erros futuros.

---

### B-05 — `//nolint` sem justificativa

Vários `//nolint:errcheck` sem comentário explicando por que o erro é seguro de ignorar.

---

## INFORMATIVO

Estes itens não representam vulnerabilidades, mas são observações de qualidade ou design.

| ID   | Descrição |
|------|-----------|
| I-01 | Argon2id parametrizado corretamente (memory, iterations, parallelism) |
| I-02 | AES-256-GCM com nonce aleatório de 12 bytes — implementação correta |
| I-03 | Sem imports de `net/http`, `os/exec` ou qualquer pacote de rede — superfície de ataque mínima |
| I-04 | `recover.go` tem lógica de recuperação de backup bem estruturada |
| I-05 | Testes unitários existem para crypto e storage, mas cobrem poucos casos de borda |
| I-06 | `detect.go` usa hash do conteúdo para detectar mudança externa — abordagem robusta |
| I-07 | Separação clara entre domínios (crypto / storage / vault / tui) facilita auditoria |
| I-08 | Nomes em português no domínio de negócio são consistentes e legíveis |
| I-09 | `go.sum` presente e completo |
| I-10 | Sem uso de `unsafe` fora do necessário para `mlock` |
| I-11 | Logs de debug não presentes em código de produção |
| I-12 | Sem hardcoded secrets ou chaves no código-fonte |

---

## Priorização de Correções

### Fase 1 — Crítico (corrigir antes do primeiro release)

1. **CRIT-01 + CRIT-02:** Implementar `WipePassword()` no storage e conectar ao `Manager.Lock()`
2. **CRIT-03:** Refatorar serialização para evitar conversão de segredos para `string`
3. **C-02:** Implementar a goroutine de auto-clear do clipboard
4. **C-01:** Registrar handler de sinal para limpeza do clipboard

### Fase 2 — Alto (corrigir antes de distribuição pública)

5. **A-08:** Usar `ReplaceFileW` no Windows para rename atômico
6. **A-10:** Corrigir nil pointer panic em `DescerPastaNaPosicao`
7. **A-05:** Implementar remoção de `.bak2` após save bem-sucedido
8. **A-06:** Zerar buffers `fileData` e `data` após uso
9. **A-01 + A-04:** Usar `SecureAllocate` para chave derivada; fortalecer `Wipe()`

### Fase 3 — Médio/Baixo (qualidade e hardening)

10. **M-02:** Desabilitar core dumps na inicialização
11. **M-01:** `mlock` para `FileRepository.password`
12. **A-02:** Validar salt mínimo de 16 bytes
13. **A-03:** Usar `utf8.RuneCountInString` em `EvaluatePasswordStrength`
14. **B-01:** Corrigir versão do `go.mod`
15. Demais findings médios e baixos

---

*Relatório gerado por análise estática manual. Nenhum teste dinâmico (fuzzing, execução instrumentada) foi realizado.*
