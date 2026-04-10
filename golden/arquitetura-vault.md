# Arquitetura — Pacote `internal/vault`

## 1. Responsabilidade do Pacote

O pacote `vault` implementa a camada de domínio do Abditum. Ele define todas as entidades do modelo de domínio, as regras de negócio que governam seus comportamentos, e o serviço de aplicação (`Manager`) que orquestra operações sobre essas entidades. O pacote é agnóstico a persistência: não sabe como os dados são escritos no disco nem como os campos são criptografados — ambos são delegados para pacotes externos via interfaces e injeção de dependência.

---

## 2. Modelo de Domínio

### 2.1. Hierarquia de Entidades

O domínio é organizado em uma hierarquia de agregação com o `Cofre` como raiz:

```
Cofre
├── Configuracoes
├── []ModeloSegredo
│   └── []CampoModelo
└── Pasta (Pasta Geral — raiz da árvore)
    ├── []Pasta (subpastas recursivas)
    └── []Segredo
        ├── []CampoSegredo
        └── CampoSegredo (observacao — campo especial)
```

**`Cofre`** é o agregado raiz. Toda mutação passa por ele ou é coordenada por ele via `Manager`. Toda persistência é atômica sobre o `Cofre` inteiro — não há salvamento parcial de entidades individuais.

**`Pasta`** é um container hierárquico. Pode conter subpastas e segredos. A raiz da árvore é a `pastaGeral`, que possui proteção especial contra renomeação e exclusão. A hierarquia é uma árvore n-ária ordenada, e a ordem dos itens dentro de cada pasta é explicitamente controlável pelo usuário.

**`Segredo`** representa uma credencial ou informação confidencial. Sua identidade é composta pelo par `(pasta, nome)` — o nome deve ser único dentro da pasta pai. Contém uma lista de campos de usuário e um campo `observacao` separado e especial.

**`CampoSegredo`** representa um campo individual dentro de um segredo. Tem nome, tipo de visibilidade e valor. O valor é armazenado como `[]byte` (não `string`) para permitir zeragem segura da memória. Sua identidade é sua posição (índice) na lista de campos.

**`ModeloSegredo`** é um template reutilizável que define a estrutura de campos (`[]CampoModelo`) para um tipo de segredo. Sua identidade é o nome, que deve ser globalmente único no cofre.

**`Configuracoes`** é um value object (struct sem ponteiros) que armazena os timers operacionais: tempo de bloqueio por inatividade, tempo para ocultar campos sensíveis e tempo para limpar a área de transferência.

### 2.2. Tipos Enumerados

**`TipoCampo`** define a visibilidade de um campo:
- `TipoCampoComum`: exibido normalmente.
- `TipoCampoSensivel`: oculto por padrão, revelado temporariamente sob demanda.

**`EstadoSessao`** rastreia o ciclo de vida de um segredo em relação ao arquivo persistido. É uma flag de sessão — nunca serializada para disco:
- `EstadoOriginal`: segredo carregado do arquivo, sem modificações na sessão atual.
- `EstadoIncluido`: segredo criado na sessão atual, ainda não existente no arquivo.
- `EstadoModificado`: segredo existia no arquivo, mas foi alterado na sessão.
- `EstadoExcluido`: segredo marcado para exclusão — removido fisicamente apenas ao salvar.

### 2.3. Referências Bidiretoras

As entidades mantêm referências ao pai para navegação eficiente:
- `Segredo.pasta` aponta para a `Pasta` pai.
- `Pasta.pai` aponta para a `Pasta` pai (`nil` para `pastaGeral`).

Essas referências **não são serializadas**. Após a desserialização, elas são reconstruídas em uma passagem única pela função `popularReferencias`, que percorre a árvore recursivamente e reconecta todos os ponteiros.

---

## 3. Convenções de Encapsulamento

### 3.1. Campos Privados com Getters Exportados

Todos os campos de todas as entidades são privados ao pacote (letra minúscula). O acesso externo ocorre exclusivamente por getters exportados. Isso garante que nenhum consumidor externo do pacote possa colocar uma entidade em estado inválido atribuindo diretamente a um campo.

### 3.2. Cópias Defensivas

Todo getter que retorna um slice retorna uma cópia do slice interno. O consumidor pode modificar a cópia retornada sem afetar o estado interno da entidade. Isso é aplicado em `Cofre.Modelos()`, `Pasta.Subpastas()`, `Pasta.Segredos()` e `Segredo.Campos()`.

Adicionalmente, `Cofre.Modelos()` aplica ordenação alfabética sobre a cópia defensiva durante o próprio getter, sem alterar o slice interno.

### 3.3. Métodos de Acesso com Semântica Explícita

O getter `CampoSegredo.ValorComoString()` é documentado como fronteira de segurança: a conversão `[]byte` → `string` é irreversível, pois strings em Go são imutáveis e não podem ser zeradas da memória. Esse getter deve ser chamado apenas quando o usuário explicitamente requisita a exibição do valor.

`Segredo.Observacao()` retorna um `string` diretamente porque o campo observação é sempre do tipo comum (não sensível).

---

## 4. Padrão Manager — Serviço de Aplicação

### 4.1. Responsabilidade

`Manager` é o único ponto de entrada para mutações no estado do cofre. A TUI, e qualquer outro consumidor, interage com o domínio exclusivamente via métodos do `Manager`. As entidades expõem seus dados de forma navegável, mas suas funções de mutação são privadas ao pacote — acessíveis apenas pelo `Manager` e pelo próprio `Cofre`.

### 4.2. Estado Mantido pelo Manager

O `Manager` mantém:
- Referência para o `Cofre` gerenciado.
- Referência para o repositório (`RepositorioCofre`) usado para persistência.
- A senha mestra (`[]byte`) para criptografia, guardada em memória enquanto o cofre está desbloqueado.
- O flag de estado de bloqueio (`bloqueado`).

### 4.3. Verificação de Bloqueio

Todos os métodos do `Manager` que acessam ou mutam dados verificam `m.bloqueado` como primeira instrução e retornam `ErrCofreBloqueado` imediatamente se o cofre estiver bloqueado.

### 4.4. Separação Validar–Mutar

O padrão de implementação de qualquer operação é:

1. Verificar se o cofre está bloqueado.
2. Executar toda a validação necessária (sem modificar nenhuma entidade).
3. Se e somente se a validação for bem-sucedida, executar a mutação.
4. Atualizar o estado global (`cofre.marcarModificado()` e, quando aplicável, `segredo.marcarModificacao()`).

Esse padrão garante que qualquer falha de validação deixa o cofre no estado exatamente igual ao anterior à chamada — não há estados parcialmente mutados.

### 4.5. Detecção de No-Op

Operações de renomeação e reposicionamento verificam se o novo valor difere do atual antes de marcar o cofre como modificado. Se a operação não produz mudança real (renomear para o mesmo nome, mover para a posição atual), o método retorna `nil` sem alterar `cofre.modificado` nem `segredo.dataUltimaModificacao`.

---

## 5. Ciclo de Vida dos Segredos

### 5.1. Soft Delete

A exclusão de segredos é reversível até o momento do salvamento. Ao excluir, o `EstadoSessao` do segredo é alterado para `EstadoExcluido`, mas o segredo permanece presente no slice `pasta.segredos`. Isso permite a operação de restauração.

A transição de estado ao excluir depende do estado atual:
- `EstadoOriginal` ou `EstadoModificado` → `EstadoExcluido`.
- `EstadoIncluido` → o segredo é removido fisicamente do slice imediatamente, pois nunca existiu no arquivo persistido.

A restauração (`RestaurarSegredo`) transiciona `EstadoExcluido` → `EstadoModificado`.

### 5.2. Visibilidade do Estado Excluído

O getter `Pasta.Segredos()` inclui todos os segredos, inclusive os marcados como `EstadoExcluido`. A filtragem de segredos excluídos para a interface é responsabilidade da TUI. A operação `Buscar` e `ListarFavoritos`, no entanto, excluem automaticamente segredos com `EstadoExcluido`.

### 5.3. Transições de EstadoSessao

| Estado Inicial  | Operação              | Estado Final    |
|-----------------|-----------------------|-----------------|
| `Original`      | Qualquer mutação      | `Modificado`    |
| `Incluido`      | Qualquer mutação      | `Incluido`      |
| `Modificado`    | Qualquer mutação      | `Modificado`    |
| `Original`      | Excluir               | `Excluido`      |
| `Modificado`    | Excluir               | `Excluido`      |
| `Incluido`      | Excluir               | Removido fisicamente |
| `Excluido`      | Restaurar             | `Modificado`    |

Operações estruturais (mover, reposicionar) **não alteram** `EstadoSessao` — são tratadas como reorganização de navegação, não como mutação de conteúdo.

Favoritar também **não altera** `EstadoSessao` nem `dataUltimaModificacao` — é classificado como preferência de navegação, não como edição de conteúdo.

---

## 6. Operações de Pasta

### 6.1. Pasta Geral

A `pastaGeral` é a raiz da árvore e possui dois níveis de proteção:
- Não pode ser renomeada (`ErrPastaGeralProtected`).
- Não pode ser excluída (`ErrPastaGeralNaoExcluivel`).

O `Manager` verifica diretamente a referência de identidade do ponteiro para aplicar essa proteção.

### 6.2. Posicionamento

A inserção de pastas e segredos usa indexação 0-based. A posição `len(lista)` é válida e significa "inserir ao final". Posições negativas ou maiores que `len` são rejeitadas com `ErrPosicaoInvalida`.

### 6.3. Detecção de Ciclo em Movimentação

Ao mover uma pasta para um novo pai (`MoverPasta`), o sistema verifica se o destino é um descendente da pasta sendo movida, o que criaria um ciclo na hierarquia. A verificação percorre recursivamente os ancestrais do destino até a raiz. Se um ciclo for detectado, `ErrCycleDetected` é retornado e nenhuma mutação ocorre.

### 6.4. Exclusão de Pasta com Promoção

Ao excluir uma pasta, todos os seus filhos (subpastas e segredos) são promovidos para a pasta pai. A semântica de conflito de nomes é:

- **Subpasta com nome coincidente no pai**: os conteúdos das duas pastas são mesclados recursivamente; a pasta coincidente no pai absorve os filhos da pasta deletada.
- **Segredo com nome coincidente no pai**: o segredo promovido é renomeado automaticamente com um sufixo numérico no formato `"Nome (N)"`, onde N incrementa a partir de 1 até não haver conflito. Cada renomeação é registrada em um `Renomeacao` struct e retornada ao chamador para exibição informativa na TUI.

Segredos com `EstadoExcluido` são promovidos mantendo esse estado.

---

## 7. Operações de Template (ModeloSegredo)

### 7.1. Invariantes

- O nome de um `ModeloSegredo` deve ser globalmente único no cofre.
- O nome `"Observação"` é reservado e não pode ser usado como nome de campo em templates.
- Templates são sempre exibidos em ordem alfabética. O getter `Cofre.Modelos()` aplica a ordenação a cada chamada sobre a cópia defensiva.

### 7.2. Exclusão sem Restrições

Templates podem ser excluídos a qualquer momento, independente de qualquer segredo referenciar ou não o template. Não há verificação de uso antes da exclusão.

### 7.3. Campos de Template

Campos de template (`CampoModelo`) podem ser adicionados em posição arbitrária, removidos por índice, e reordenados por índice de origem e destino. Todas essas operações validam os índices antes de mutar.

---

## 8. Campo Observação

O campo `observacao` existe em todo `Segredo` e é tratado como entidade separada dos campos de usuário. Não aparece no slice `campos` retornado por `Segredo.Campos()`. É serializado separadamente como campo `observacao` no JSON. Seu nome `"Observação"` é reservado e não pode ser usado como nome de campo em templates nem em campos de usuário. Ele é sempre do tipo `TipoCampoComum`. Seu valor é zerável pela operação de bloqueio. O comprimento máximo do seu valor é 1000 caracteres (`ErrObservacaoMuitoLonga`).

---

## 9. Salvamento — Comprometimento em Três Fases

A operação `Salvar` implementa um protocolo de comprometimento que garante que falhas de persistência não causem perda de dados ou inconsistência no estado em memória:

**Fase 1 — Snapshot**: Cria-se uma cópia profunda do `Cofre` em memória. Durante a cópia, segredos com `EstadoExcluido` são filtrados — o snapshot não os contém. O vault em uso permanece intacto durante toda essa fase.

**Fase 2 — Persistência**: O snapshot é entregue ao repositório para escrita em disco. Se a persistência falhar, o método retorna o erro e o vault em memória permanece exatamente como estava antes da chamada.

**Fase 3 — Finalização**: Somente após confirmação de sucesso na fase 2, os segredos com `EstadoExcluido` são removidos fisicamente do vault em memória e `cofre.modificado` é resetado para `false`.

---

## 10. Bloqueio e Zeramento de Memória Sensível

A operação `Lock()` executa o seguinte protocolo de limpeza:

1. Zera a senha mestra (`[]byte`) usando `crypto.Wipe` e atribui `nil`.
2. Percorre recursivamente toda a árvore de pastas e, para cada segredo, chama `zerarValoresSensiveis()`, que zera todos os campos do tipo `TipoCampoSensivel` e o campo `observacao` — ambos com `crypto.Wipe` seguido de atribuição `nil`.
3. Atribui `nil` à referência do `Cofre` no `Manager`.
4. Seta `bloqueado = true`.

Após o bloqueio, `Vault()` retorna `nil` e todos os métodos de mutação retornam `ErrCofreBloqueado`.

---

## 11. Busca e Listagem

**`Buscar`** realiza busca case-insensitive por substring na árvore completa de pastas. Pesquisa o nome do segredo, o valor da observação, os nomes de todos os campos (independente do tipo) e os valores dos campos do tipo `TipoCampoComum`. Valores de campos `TipoCampoSensivel` são excluídos da busca por segurança. Segredos com `EstadoExcluido` são excluídos dos resultados.

**`ListarFavoritos`** retorna todos os segredos com `favorito = true` via DFS sobre a árvore completa, excluindo segredos com `EstadoExcluido`.

---

## 12. Serialização

### 12.1. Formato

O domínio serializa para JSON via structs intermediários (`DTO`) privadas ao arquivo `serialization.go`. As entidades de domínio não têm tags `json` — a separação entre o modelo de domínio e o modelo de serialização é total.

### 12.2. Convenções de Serialização

- O campo `pasta_geral.nome` é sempre escrito como `"Geral"` no JSON independente do nome em memória.
- Valores de campo (`CampoSegredo.valor`) são serializados como strings UTF-8 simples, sem encoding Base64.
- Ao desserializar, todos os segredos recebem `EstadoOriginal` — o estado de sessão é sempre resetado ao carregar.
- `cofre.modificado` é sempre `false` após desserialização.
- O campo `version` passado a `DeserializarCofre` é reservado para compatibilidade futura de formato, permitindo seleção de campos legados por versão.

### 12.3. Validação na Desserialização

A desserialização valida que `pasta_geral` está presente e que `pasta_geral.nome == "Geral"`. Qualquer violação retorna um erro sentinel específico (`ErrPastaGeralAusente`, `ErrPastaGeralNomeInvalido`).

---

## 13. Interface de Repositório

O pacote define a interface `RepositorioCofre` com dois métodos:

- `Salvar(cofre *Cofre) error`: persiste o vault (snapshot) recebido.
- `Carregar() (*Cofre, error)`: carrega e retorna um vault do armazenamento.

A implementação concreta dessa interface está no pacote `internal/storage`. O `Manager` recebe a implementação via injeção de dependência no construtor `NewManager`. Isso permite que testes usem implementações mock do repositório sem depender de disco.

---

## 14. Erros

Todos os erros do pacote são sentinel errors (variáveis `var Err... = errors.New(...)`). O consumidor (TUI) usa `errors.Is()` para identificar o tipo de falha e exibir a mensagem adequada ao usuário. Nenhum erro inclui contexto dinâmico — erros de domínio descrevem a violação de regra, não os valores que causaram a violação.

---

## 15. Inicialização

**`NovoCofre()`** cria um `Cofre` vazio com `pastaGeral` inicializada e configurações padrão (5 min bloqueio, 15 s revelar, 30 s limpar área de transferência). O vault não é marcado como modificado.

**`InicializarConteudoPadrao()`** popula o cofre recém-criado com pastas e templates padrão. Por decisão explícita, essa operação **não** marca o vault como modificado — o conteúdo padrão não deve forçar um salvamento imediato após a criação do cofre.

---

## 16. Duplicação de Segredo

A duplicação cria uma cópia profunda do segredo com cópias independentes de todos os `[]byte` de valor de campo. O nome do duplicado segue o esquema `"Nome (2)"`, `"Nome (3)"`, etc., incrementando o sufixo até que não haja conflito no mesmo pasta pai. O duplicado recebe `EstadoIncluido`. O original não é alterado.

---

## 17. Regras Detalhadas por Operação

Esta seção especifica exatamente o que cada operação do `Manager` aceita, rejeita e produz como efeito colateral. A regra universal, omitida em cada entrada mas sempre aplicável, é: **se o cofre estiver bloqueado, qualquer operação retorna `ErrCofreBloqueado` antes de executar qualquer outra verificação**.

A convenção de efeito colateral "marca modificado" significa `cofre.marcarModificado()` (atualiza `cofre.modificado = true` e `cofre.dataUltimaModificacao`). "Marca modificação do segredo" significa `segredo.marcarModificacao()` (atualiza `segredo.dataUltimaModificacao`).

---

### 17.1. Operações de Cofre

#### `AlterarConfiguracoes(novasConfig)`

| Condição | Resultado |
|---|---|
| `tempoBloqueioInatividadeMinutos` ≤ 0 | `ErrConfigInvalida` |
| `tempoOcultarSegredoSegundos` ≤ 0 | `ErrConfigInvalida` |
| `tempoLimparAreaTransferenciaSegundos` ≤ 0 | `ErrConfigInvalida` |
| Todos os timers > 0 | Atualiza configurações; marca modificado |

#### `Lock()`

Não retorna erro. Se já bloqueado, é no-op. Caso contrário: zera senha mestra, zera campos sensíveis e observação de todos os segredos recursivamente, atribui `nil` ao cofre, seta `bloqueado = true`.

#### `Salvar()`

| Condição | Resultado |
|---|---|
| Repositório retorna erro | Retorna o erro; vault em memória inalterado |
| Sucesso | Remove `EstadoExcluido` fisicamente da memória; `cofre.modificado = false` |

#### `Buscar(consulta)`

Não retorna erro. Se bloqueado, retorna `nil`. Caso contrário, busca case-insensitive por substring em: nome do segredo, texto da observação, **nomes** de todos os campos (comuns e sensíveis), **valores** apenas dos campos `TipoCampoComum`. Exclui segredos com `EstadoExcluido`.

#### `ListarFavoritos()`

Não retorna erro. Se bloqueado, retorna `nil`. Caso contrário, retorna via DFS todos os segredos com `favorito = true`, excluindo `EstadoExcluido`.

---

### 17.2. Operações de Pasta

#### `CriarPasta(pai, nome, posicao)`

| Condição | Resultado |
|---|---|
| `nome` vazio | `ErrNomeVazio` |
| `len(nome)` > 255 bytes | `ErrNomeMuitoLongo` |
| Nome já existe entre subpastas do `pai` | `ErrNameConflict` |
| `posicao` < 0 ou `posicao` > `len(pai.subpastas)` | `ErrPosicaoInvalida` |
| `posicao` == `len(pai.subpastas)` | Válido; insere ao final |
| Sucesso | Cria pasta no `pai` na posição indicada; marca modificado |

#### `RenomearPasta(pasta, novoNome)`

| Condição | Resultado |
|---|---|
| `pasta` é `pastaGeral` (`pai == nil`) | `ErrPastaGeralProtected` |
| `novoNome` vazio | `ErrNomeVazio` |
| `len(novoNome)` > 255 bytes | `ErrNomeMuitoLongo` |
| `novoNome` igual ao de alguma irmã (exceto a própria) | `ErrNameConflict` |
| `novoNome` igual ao nome atual | No-op; não marca modificado |
| Sucesso | Renomeia; marca modificado |

#### `MoverPasta(pasta, destino)`

| Condição | Resultado |
|---|---|
| `pasta` é `pastaGeral` | `ErrPastaGeralProtected` |
| `pasta == destino` | `ErrDestinoInvalido` |
| `destino` é descendente de `pasta` | `ErrCycleDetected` |
| Nome de `pasta` já existe entre subpastas do `destino` | `ErrNameConflict` |
| Sucesso | Move `pasta` para o final das subpastas do `destino`; marca modificado |

#### `ReposicionarPasta(pasta, novaPosicao)`

| Condição | Resultado |
|---|---|
| `pasta` é `pastaGeral` (`pai == nil`) | `ErrPastaGeralProtected` |
| `novaPosicao` < 0 ou ≥ `len(pai.subpastas)` | `ErrPosicaoInvalida` |
| `novaPosicao` == posição atual | No-op; não marca modificado |
| Sucesso | Reposiciona dentro do mesmo pai; marca modificado |

#### `SubirPastaNaPosicao(pasta)`

| Condição | Resultado |
|---|---|
| Posição atual == 0 | No-op; retorna `nil` sem marcar modificado |
| Sucesso | Equivale a `ReposicionarPasta(pasta, posicaoAtual - 1)` |

#### `DescerPastaNaPosicao(pasta)`

| Condição | Resultado |
|---|---|
| Posição atual == último índice | No-op; retorna `nil` sem marcar modificado |
| Sucesso | Equivale a `ReposicionarPasta(pasta, posicaoAtual + 1)` |

#### `ExcluirPasta(pasta)`

| Condição | Resultado |
|---|---|
| `pasta` é `pastaGeral` (`pai == nil`) | `ErrPastaGeralNaoExcluivel` |
| Subpasta da excluída tem nome igual a subpasta existente no pai | Mescla conteúdos recursivamente (sem erro) |
| Segredo da excluída tem nome igual a segredo existente no pai | Renomeia com sufixo `"(N)"` e registra em `[]Renomeacao` retornado |
| Sucesso | Promove todas as subpastas e segredos ao pai; remove a pasta; marca modificado |

Segredos com `EstadoExcluido` são sempre promovidos mantendo esse estado, sem participar da busca de conflitos por nome.

---

### 17.3. Operações de Segredo

#### `CriarSegredo(pasta, nome, modelo)`

| Condição | Resultado |
|---|---|
| `nome` vazio | `ErrNomeVazio` |
| `len(nome)` > 255 bytes | `ErrNomeMuitoLongo` |
| `modelo == nil` | `ErrModeloInvalido` |
| Nome já existe na pasta | `ErrNameConflict` |
| Sucesso | Cria segredo com campos do modelo (valores vazios), `EstadoIncluido`, `favorito = false`; marca modificado e modificação do segredo |

#### `ExcluirSegredo(segredo)`

| Condição | Resultado |
|---|---|
| `segredo == nil` | `ErrSegredoInvalido` |
| `estadoSessao == EstadoExcluido` | `ErrSegredoJaExcluido` |
| `estadoSessao == EstadoIncluido` | Remove fisicamente do slice da pasta; marca modificado (sem marcar modificação do segredo) |
| `estadoSessao == EstadoOriginal` ou `EstadoModificado` | Transiciona para `EstadoExcluido`; marca modificado e modificação do segredo |

#### `RestaurarSegredo(segredo)`

| Condição | Resultado |
|---|---|
| `segredo == nil` | `ErrSegredoInvalido` |
| `estadoSessao != EstadoExcluido` | `ErrSegredoNaoExcluido` |
| Sucesso | Transiciona para `EstadoModificado`; marca modificado e modificação do segredo |

#### `AlternarFavoritoSegredo(segredo)`

| Condição | Resultado |
|---|---|
| `segredo == nil` | `ErrSegredoInvalido` |
| `estadoSessao == EstadoExcluido` | `ErrSegredoJaExcluido` |
| Sucesso | Inverte `favorito`; marca **apenas** modificado — **não** marca modificação do segredo nem altera `estadoSessao` |

#### `DuplicarSegredo(segredo)`

| Condição | Resultado |
|---|---|
| `segredo == nil` | `ErrSegredoInvalido` |
| `estadoSessao == EstadoExcluido` | `ErrSegredoJaExcluido` |
| `segredo.pasta == nil` | `ErrPastaInvalida` |
| Sucesso | Cria cópia profunda com nome único (`"Nome (1)"`, `"Nome (2)"`, etc.), `EstadoIncluido`, `favorito = false`; marca modificado e modificação do duplicado |

#### `RenomearSegredo(segredo, novoNome)`

| Condição | Resultado |
|---|---|
| `novoNome` vazio | `ErrNomeVazio` |
| `len(novoNome)` > 255 bytes | `ErrNomeMuitoLongo` |
| Nome já existe na mesma pasta (outro segredo) | `ErrNameConflict` |
| `novoNome` igual ao nome atual | No-op; não marca modificado nem modificação do segredo |
| `estadoSessao == EstadoOriginal` | Transiciona para `EstadoModificado`; marca modificado e modificação do segredo |
| `estadoSessao == EstadoIncluido` ou `EstadoModificado` | Estado permanece; marca modificado e modificação do segredo |

#### `EditarCampoSegredo(segredo, indice, novoValor)`

| Condição | Resultado |
|---|---|
| `indice` < 0 ou ≥ `len(segredo.campos)` | `ErrCampoInvalido` |
| `novoValor` igual ao valor atual | No-op; não marca modificado nem modificação do segredo |
| `estadoSessao == EstadoOriginal` | Transiciona para `EstadoModificado`; marca modificado e modificação do segredo |
| `estadoSessao == EstadoIncluido` ou `EstadoModificado` | Estado permanece; marca modificado e modificação do segredo |

O valor é armazenado como cópia profunda (`[]byte` novo).

#### `EditarObservacao(segredo, novoTexto)`

| Condição | Resultado |
|---|---|
| `len(novoTexto)` > 1000 bytes | `ErrObservacaoMuitoLonga` |
| `novoTexto` igual ao texto atual | No-op; não marca modificado nem modificação do segredo |
| `estadoSessao == EstadoOriginal` | Transiciona para `EstadoModificado`; marca modificado e modificação do segredo |
| `estadoSessao == EstadoIncluido` ou `EstadoModificado` | Estado permanece; marca modificado e modificação do segredo |

#### `MoverSegredo(segredo, destino, posicao)`

| Condição | Resultado |
|---|---|
| `destino == nil` | `ErrPastaInvalida` |
| Nome já existe no `destino` (outro segredo) | `ErrNameConflict` |
| `posicao` < 0 ou ≥ `len(destino.segredos)` | Posição ignorada silenciosamente; segredo inserido ao final |
| Sucesso | Move para `destino` na posição indicada; marca **apenas** modificado — **não** altera `estadoSessao` |

#### `ReposicionarSegredo(segredo, novaPosicao)`

| Condição | Resultado |
|---|---|
| `segredo.pasta == nil` | `ErrPastaInvalida` |
| `novaPosicao` < 0 ou ≥ `len(pasta.segredos)` | `ErrPosicaoInvalida` |
| `novaPosicao` == posição atual | No-op; não marca modificado |
| Sucesso | Reposiciona dentro da mesma pasta; marca **apenas** modificado — **não** altera `estadoSessao` |

#### `SubirSegredoNaPosicao(segredo)`

| Condição | Resultado |
|---|---|
| Posição atual == 0 | No-op; retorna `nil` sem marcar modificado |
| Sucesso | Equivale a `ReposicionarSegredo(segredo, posicaoAtual - 1)` |

#### `DescerSegredoNaPosicao(segredo)`

| Condição | Resultado |
|---|---|
| Posição atual == último índice | No-op; retorna `nil` sem marcar modificado |
| Sucesso | Equivale a `ReposicionarSegredo(segredo, posicaoAtual + 1)` |

---

### 17.4. Operações de Template (`ModeloSegredo`)

#### `CriarModelo(nome, campos)`

| Condição | Resultado |
|---|---|
| `nome` vazio | `ErrNomeVazio` |
| `len(nome)` > 255 bytes | `ErrNomeMuitoLongo` |
| Nome já existe em outro template do cofre | `ErrNameConflict` |
| Algum campo em `campos` tem nome reservado (`"Observação"`, case-insensitive) | `ErrObservacaoReserved` |
| Sucesso | Insere template em posição alfabética; marca modificado |

#### `RenomearModelo(modelo, novoNome)`

| Condição | Resultado |
|---|---|
| `novoNome` vazio | `ErrNomeVazio` |
| `len(novoNome)` > 255 bytes | `ErrNomeMuitoLongo` |
| Nome já existe em outro template (exceto o próprio) | `ErrNameConflict` |
| `novoNome` igual ao nome atual | No-op; não marca modificado |
| Sucesso | Renomeia; marca modificado |

#### `ExcluirModelo(modelo)`

Sem validações de domínio além do estado de bloqueio. Remove o template do cofre; marca modificado.

#### `AdicionarCampo(modelo, nome, tipo, posicao)`

| Condição | Resultado |
|---|---|
| `nome` é reservado (`"Observação"`, case-insensitive) | `ErrObservacaoReserved` |
| `posicao` < 0 ou `posicao` > `len(modelo.campos)` | `ErrPosicaoInvalida` |
| `posicao` == `len(modelo.campos)` | Válido; insere ao final |
| Sucesso | Insere campo na posição; marca modificado |

Não valida nome vazio nem comprimento máximo de nome de campo.

#### `RemoverCampo(modelo, indice)`

| Condição | Resultado |
|---|---|
| `indice` < 0 ou ≥ `len(modelo.campos)` | `ErrCampoInvalido` |
| Sucesso | Remove o campo; marca modificado |

#### `ReordenarCampo(modelo, indiceOrigem, indiceDestino)`

| Condição | Resultado |
|---|---|
| `indiceOrigem` < 0 ou ≥ `len(modelo.campos)` | `ErrCampoInvalido` |
| `indiceDestino` < 0 ou ≥ `len(modelo.campos)` | `ErrCampoInvalido` |
| `indiceOrigem == indiceDestino` | Executa sem erro; comportamento é no-op efetivo |
| Sucesso | Move o campo; marca modificado |
