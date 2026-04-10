# Arquitetura de Testes de UI com Golden Files — `internal/tui`

> Descreve os padrões, convenções e decisões arquiteturais que governam como os testes de aparência e comportamento dos componentes TUI são organizados e executados.

---

## Sumário

- [Contexto e Motivação](#contexto-e-motivação)
- [Estrutura de Arquivos](#estrutura-de-arquivos)
- [As Duas Categorias de Teste](#as-duas-categorias-de-teste)
- [Golden Files: Dois Níveis de Validação](#golden-files-dois-níveis-de-validação)
  - [Nível 1 — `.txt.golden`: Layout e Estrutura](#nível-1--txtgolden-layout-e-estrutura)
  - [Nível 2 — `.json.golden`: Estilos e Cores](#nível-2--jsongolden-estilos-e-cores)
- [Convenção de Nomenclatura dos Arquivos Golden](#convenção-de-nomenclatura-dos-arquivos-golden)
- [O Subpacote `testdata`](#o-subpacote-testdata)
  - [`StyleTransition` e Formato de Tupla](#styletransition-e-formato-de-tupla)
  - [Canonicidade do Parser](#canonicidade-do-parser)
- [O Flag `-update`](#o-flag--update)
- [O Ciclo de Vida dos Golden Files](#o-ciclo-de-vida-dos-golden-files)
- [Helpers de Teste: Padrão `checkOrUpdate*Golden`](#helpers-de-teste-padrão-checkOrupdate-golden)
  - [Variantes por componente](#variantes-por-componente)
- [Cobertura por Componente](#cobertura-por-componente)
  - [Dimensões de variação](#dimensões-de-variação)
  - [Catálogo de cenários por componente](#catálogo-de-cenários-por-componente)
- [Fixtures e Construtores de Teste](#fixtures-e-construtores-de-teste)
- [Testes Comportamentais com Dados da `json.golden`](#testes-comportamentais-com-dados-da-jsongolden)
- [Testes Comportamentais Puros (`Update()`)](#testes-comportamentais-puros-update)
- [Testes de Integração de Fluxo](#testes-de-integração-de-fluxo)
- [Separação entre Testes Visuais e Comportamentais](#separação-entre-testes-visuais-e-comportamentais)
- [Decisões Arquiteturais Relevantes](#decisões-arquiteturais-relevantes)

---

## Contexto e Motivação

Componentes TUI produzem `string` com códigos ANSI embutidos. Testes convencionais baseados em `strings.Contains()` são insuficientes para capturar as categorias de regressão mais comuns em interfaces terminais:

- Espaçamento incorreto entre elementos
- Cor ou estilo aplicados na posição errada
- Truncamento ou quebra de linha inesperada
- Alinhamento de bordas e símbolos

O caso motivador documentado é o bug de espaçamento em `DecisionDialog.renderActionBar()`: a string `"backup──────"` em vez de `"backup ──────"` passou em todos os testes baseados em `strings.Contains()` porque nenhum verificava o espaço entre o rótulo e o traço.

A solução adotada é uma arquitetura de **golden files com dois níveis**: um nível captura a estrutura visual (layout puro, sem cores), outro captura as transições de estilo (cores e fontes, sem texto). Os dois juntos cobrem a totalidade do output visível sem depender dos códigos ANSI brutos.

---

## Estrutura de Arquivos

```
internal/tui/
├── testdata/
│   ├── ansiparser.go         — subpacote testdata: parser ANSI + serializer JSON
│   ├── ansiparser_test.go    — testes do parser
│   └── golden/               — todos os arquivos golden do pacote tui
│       ├── {component}-{variant}-{size}.txt.golden
│       └── {component}-{variant}-{size}.json.golden
├── messages_test.go          — declara flag -update; helpers goldenPath + checkOrUpdateGolden
├── actions_test.go           — testes golden de commandbar
├── welcome_test.go           — testes golden de welcome
├── decision_test.go          — fixtures PoC + testes golden de decision
├── help_test.go              — testes golden de helpModal
├── passwordentry_test.go     — testes golden + comportamentais de passwordEntry
├── passwordcreate_test.go    — testes golden + comportamentais de passwordCreate
├── filepicker_test.go        — testes golden + comportamentais de filePicker
├── root_test.go              — testes golden + comportamentais de rootModel
└── exit_flow_integration_test.go — testes de integração de fluxo (sem golden)
```

Todos os arquivos `*_test.go` pertencem ao mesmo pacote `tui` (white-box). Os helpers são funções privadas declaradas nos arquivos de teste e compartilhadas por pertencerem ao mesmo pacote de compilação de teste.

---

## As Duas Categorias de Teste

Cada arquivo de teste de componente contém, em geral, duas categorias de teste:

**1. Testes de aparência (golden):** Chamam `View()` em um estado determinístico e comparam o output com arquivos gravados em disco. Detectam regressões visuais byte a byte.

**2. Testes comportamentais (unitários):** Chamam `Update()` com mensagens específicas e verificam o estado resultante ou os comandos emitidos. Detectam regressões de lógica de negócio.

As duas categorias coexistem nos mesmos arquivos de teste, separadas por comentários de bloco.

---

## Golden Files: Dois Níveis de Validação

### Nível 1 — `.txt.golden`: Layout e Estrutura

O arquivo `.txt.golden` armazena o output visual do componente com **todos os códigos ANSI removidos**, produzindo texto puro. É o "desenho" canônico da tela — bordas, spacing, quebras de linha e conteúdo textual.

A conversão é feita pela função `stripANSI()`, que aplica uma regex `\x1b\[[0-9;]*m` sobre o output bruto e descarta todas as sequências de escape SGR.

**O que este nível detecta:**
- Espaçamento ausente ou a mais entre elementos
- Caracteres de borda trocados ou faltando
- Wrapping de texto inesperado
- Número de linhas diferente do esperado
- Conteúdo textual diferente do esperado

**O que este nível não captura:**
- Troca de cor sem troca de conteúdo
- Remoção ou adição de `bold`, `italic`, `underline`

### Nível 2 — `.json.golden`: Estilos e Cores

O arquivo `.json.golden` armazena as **transições de estilo visual** extraídas do output ANSI original, no formato de um array JSON de tuplas compactas. Cada tupla registra linha, coluna, cor de foreground, cor de background e atributos de fonte no ponto onde o estado visual muda.

O formato de cada entrada é `[linha, coluna, fg_hex|null, bg_hex|null, [estilos]]`.

**O que este nível detecta:**
- Cor trocada em qualquer posição
- `bold`, `italic`, `underline` adicionados ou removidos
- Cor de borda usando token errado (e.g., `semantic.error` em vez de `border.focused`)

**O que este nível não captura:**
- Espaçamento ou conteúdo textual

Os dois níveis são complementares e sempre criados em par para o mesmo cenário.

---

## Convenção de Nomenclatura dos Arquivos Golden

O padrão é:

```
{component}-{variant}-{size}.{ext}.golden
```

- `component` — nome do tipo/componente em minúsculas: `commandbar`, `messages`, `decision`, `help`, `passwordentry`, `passwordcreate`, `filepicker`, `welcome`, `root`, `flow-create-vault`, `flow-open-vault`
- `variant` — descritor do cenário; pode ser composto com hifens: `success`, `destructive-1action-short`, `15actions-mid`, `open-withfiles`, `welcome-initial`
- `size` — largura do terminal para componentes de linha única (`30`, `60`, `80`); dimensão completa `WxH` para componentes bidimensionais (`80x24`, `30x16`)
- `ext` — `txt` ou `json`

Exemplos:

| Arquivo | Componente | Variante | Tamanho |
|---|---|---|---|
| `messages-success-30.txt.golden` | message bar | success | 30 colunas |
| `decision-destructive-1action-short-60x24.json.golden` | decision dialog | destrutivo, 1 ação, título curto | 60×24 |
| `help-15actions-bottom-30x16.txt.golden` | help modal | 15 ações, scroll no fim | 30×16 |
| `filepicker-open-withfiles-80x24.json.golden` | file picker | modo open, com arquivos | 80×24 |
| `welcome-cyberpunk-80.txt.golden` | welcome screen | tema cyberpunk | 80 colunas |

Essa convenção é **estritamente descritiva**: o nome do arquivo deve identificar o cenário sem necessidade de ler o conteúdo.

---

## O Subpacote `testdata`

O diretório `testdata/` é um pacote Go independente (`package testdata`), não um diretório de dados estáticos. Contém:

- `ansiparser.go` — parser SGR, tipo `StyleTransition`, funções `ParseANSIStyle()` e `MarshalStyleTransitions()`
- `ansiparser_test.go` — testes unitários do próprio parser

Os arquivos de teste do pacote `tui` importam este subpacote com alias: `testdatapkg "github.com/useful-toys/abditum/internal/tui/testdata"`.

### `StyleTransition` e Formato de Tupla

`StyleTransition` representa um ponto de mudança de estado visual:

| Campo | Tipo | Semântica |
|---|---|---|
| `Line` | `int` | Linha (0-indexed) onde a mudança ocorre |
| `Col` | `int` | Coluna (0-indexed) onde a mudança ocorre |
| `FG` | `*string` | Hex da cor do foreground, ou `null` se default |
| `BG` | `*string` | Hex da cor do background, ou `null` se não aplicado |
| `Style` | `[]string` | Atributos ativos: `"bold"`, `"italic"`, `"underline"`, `"strikethrough"`, `"faint"`, `"blink"`, `"reverse"` |

A serialização usa `MarshalJSON()` customizado que produz o formato de tupla compacta `[line, col, fg, bg, [styles]]`. O serializer externo `MarshalStyleTransitions()` produz um array com indentação externa (um item por linha), mas cada tupla colapsada em uma linha — maximizando legibilidade do diff.

### Canonicidade do Parser

O parser é **canônico e agnóstico às sequências SGR brutas**: diferentes sequências que produzem o mesmo estado visual são normalizadas para a mesma tupla. Sequências redundantes (que não mudam o estado visual) são silenciadas. Múltiplas sequências SGR na mesma posição de coluna são colapsadas em uma única tupla com o estado final.

Exemplos de equivalência tratados:
- Cores 16-bit, 256-bit e true color normalizadas para hex
- Ordem de SGR intercambiável (`\x1b[1m\x1b[38;5;208m` = `\x1b[38;5;208m\x1b[1m`)
- Código `0` (reset) zera fg, bg e todos os estilos

O resultado é que os arquivos `.json.golden` podem ser comparados byte-a-byte sem falsos positivos causados por refatorações que preservam o output visual mas alteram os códigos SGR internos.

---

## O Flag `-update`

O flag `-update` é declarado **uma única vez** em `messages_test.go`:

```
var update = flag.Bool("update", false, "regenerate golden files")
```

Como todos os arquivos `*_test.go` do pacote `tui` compilam juntos, esta variável é acessível em todos eles. O flag é documentado no comentário de cada função de golden test com o comando de uso:

```
// Usage: go test ./internal/tui/... -run TestXxx_Golden -update
```

Há três modos de operação:

| Condição | Comportamento |
|---|---|
| Arquivo não existe (primeira execução) | Cria o arquivo e considera o teste aprovado (bootstrap de baseline) |
| Arquivo existe, sem `-update` | Compara byte a byte; falha com diff se diferente |
| Arquivo existe, com `-update` | Sobrescreve incondicionalmente; teste passa sempre |

Esta convenção evita que a primeira execução falhe por ausência de arquivos e permite regeneração controlada após mudanças visuais intencionais.

---

## O Ciclo de Vida dos Golden Files

```
1. NOVO COMPONENTE
   └── go test ./internal/tui/...
       → golden ausente: cria arquivo e passa
       → developer commit: golden files vão para git

2. DESENVOLVIMENTO NORMAL (sem mudança visual)
   └── go test ./internal/tui/...
       → compara byte a byte: PASS

3. MUDANÇA VISUAL INTENCIONAL (nova cor, novo spacing)
   └── go test ./internal/tui/... -run TestXxx_Golden -update
       → regenera baselines
       → developer revisa diff no git: git diff testdata/golden/
       → commit com mensagem descretiva

4. REGRESSÃO VISUAL (sem -update)
   └── go test ./internal/tui/...
       → mismatch: FAIL com want/got
       → developer investiga
```

Os arquivos golden são versionados no repositório. O diff dos arquivos é a evidência humana de mudanças visuais — cada modificação de output passa por revisão via git, tornando regressões inadvertidas impossíveis de passar silenciosamente.

---

## Helpers de Teste: Padrão `checkOrUpdate*Golden`

Todos os helpers de golden seguem o mesmo contrato:

- Recebem caminho do arquivo e o conteúdo obtido
- Se `-update` estiver ativo: cria o diretório (se necessário) e grava incondicionalmente
- Se o arquivo não existir: cria — bootstrap automático de baseline
- Se o arquivo existir e conteúdo diferir: `t.Errorf()` com diff legível `want:\n...\ngot:\n...`

A função base `checkOrUpdateGolden(t, path, got string)` está em `messages_test.go`.

### Variantes por componente

Cada componente que usa nomenclatura diferente do padrão genérico define sua própria função de helper de path:

| Função helper | Arquivo | Particularidade |
|---|---|---|
| `goldenPath(component, variant, width, ext)` | `messages_test.go` | Padrão genérico `{component}-{variant}-{width}.{ext}.golden` |
| `decisionGoldenPath(variant, ext)` | `decision_test.go` | Sem campo de largura separado — width já está embutido na variante (`-30x24`) |
| `helpGoldenPath(variant, ext)` | `help_test.go` | Sem campo de largura separado — width já está na variante |

Os helpers de path específicos são acompanhados de helpers `checkOrUpdate*Golden` locais que internalizam o caminho, simplificando as chamadas dentro das funções de teste. O comportamento de update/bootstrap é idêntico ao genérico.

---

## Cobertura por Componente

### Dimensões de variação

Cada componente é testado ao longo de múltiplas dimensões ortogonais:

| Dimensão | Descrição | Exemplos |
|---|---|---|
| **Largura de terminal** | Largura estreita e normal | 30, 60, 80 colunas |
| **Altura de terminal** | Para componentes scrolláveis | `x12`, `x16`, `x24` |
| **Variante de conteúdo** | Comprimento de título/corpo | `short`, `long` |
| **Número de ações** | Quantidade de botões/ações | `1action`, `2action`, `3action` |
| **Severidade/Tema** | Estado semântico | `destructive`, `error`, `alert`, `informative`, `neutral` |
| **Estado de dados** | Preenchimento de campos, scroll | `initial`, `filled`, `empty`, `top`, `mid`, `bottom` |
| **Modo operacional** | Open/Save, com/sem arquivos | `open-withfiles`, `save-fieldempty` |
| **Tema visual** | Tokyo Night / Cyberpunk | `tokyo-night`, `cyberpunk` |

### Catálogo de cenários por componente

| Componente | Arquivo de teste | Cenários golden | Arquivos golden |
|---|---|---|---|
| `commandbar` | `actions_test.go` | 5 variantes × 2 larguras | 20 (10 txt + 10 json) |
| `messages` | `messages_test.go` | 6 `MsgKind` × 2 larguras | 24 (12 txt + 12 json) |
| `welcome` | `welcome_test.go` | 2 temas × 1 largura | 4 (2 txt + 2 json) |
| `decision` | `decision_test.go` | 10 cenários × 2 larguras | 40 (20 txt + 20 json) |
| `help` | `help_test.go` | 4 configurações × 2 larguras | 16 (8 txt + 8 json) |
| `passwordentry` | `passwordentry_test.go` | 2 estados | 4 (2 txt + 2 json) |
| `passwordcreate` | `passwordcreate_test.go` | 3 estados | 6 (3 txt + 3 json) |
| `filepicker` | `filepicker_test.go` | ~8 cenários | 16 (8 txt + 8 json) |
| `root` | `root_test.go` | 4 cenários | 8 (4 txt + 4 json) |
| `flow-open-vault` / `flow-create-vault` | flow tests | 1 cenário por flow | 4 (2 txt + 2 json) |

**Total aproximado:** 142 arquivos golden no diretório `testdata/golden/`.

---

## Fixtures e Construtores de Teste

Para componentes com muitas combinações (especialmente `DecisionDialog`), os arquivos de teste definem **construtores de fixture nomeados**:

- `pocKey1()` através de `pocKeyF()` — 15 funções correspondendo a 15 combinações da matriz PoC (5 severidades × 3 contagens de ações), cada uma chamando `SetSize(80, 24)`
- `newDestructive1Short(w int)`, `newAlert3Long(w int)`, etc. — construtores paramétricos por largura, usados exclusivamente nos testes golden para garantir renderização determinística em múltiplas larguras

Essa separação entre "fixtures PoC" (estado 80×24 para testes comportamentais) e "construtores golden" (largura parametrizável para testes de aparência) evita duplicação e mantém os cenários comportamentais independentes de variações de largura.

Para `helpModal`, helpers `help3actions()`, `help15actions()` e `helpGroupLabel()` produzem as listas de ações que alimentam múltiplos cenários de scroll — incluindo um sentinela `-1` para indicar "scroll máximo" calculado automaticamente.

---

## Testes Comportamentais com Dados da `json.golden`

Além de simplesmente gravar e comparar, o parser ANSI é usado diretamente em **testes comportamentais** que fazem asserções sobre a estrutura de estilo extraída do output, sem comparação com arquivo em disco.

O padrão é:

1. Chamar `View()` no modelo configurado
2. Chamar `testdatapkg.ParseANSIStyle(out)` para extrair transitions
3. Inspecionar a slice de transições para verificar presença/ausência de atributos em linhas específicas

Exemplos presentes no código:

- `TestPasswordEntryModal_ConfirmarDisabledWhenEmpty` — quando o campo de senha está vazio, nenhuma transição com `"bold"` deve aparecer nas últimas 3 linhas (área de ação). Verifica D-PE-03.
- `TestPasswordEntryModal_ConfirmarActiveWhenFilled` — quando o campo tem valor, deve existir pelo menos uma transição com `"bold"` nas últimas 3 linhas. Verifica D-PE-04.

Este uso do parser para asserções programáticas (fora do contexto de comparação com arquivo) é uma decisão deliberada: permite verificar invariantes de estilo dependentes de estado sem necessidade de um golden file por estado.

---

## Testes Comportamentais Puros (`Update()`)

Cada arquivo de componente contém uma seção de testes unitários para `Update()`. Estes testes:

- **Não usam golden files**
- Verificam a emissão correta de `tea.Cmd` e transformações de estado
- Cobrem casos como: navegação de cursor, validação em tempo real, rejeição de entrada inválida, ciclo de foco entre campos

Exemplos de categorias de `Update()` testadas:

| Componente | O que é testado |
|---|---|
| `DecisionDialog` | Enter dispara ação padrão; Esc dispara cancel; teclas explícitas (`"m"`, `"a"`, `"t"`) disparam ações nomeadas; tecla desconhecida retorna nil |
| `helpModal` | Down aumenta scroll; Up não vai abaixo de 0; F1 e Esc fecham o modal |
| `filePickerModal` | Tab alterna painel; Down/Up movem cursor; Enter em diretório expande nó; ESC emite `popModalMsg`; modo Save vs Open têm comportamento Tab diferente |
| `passwordEntryModal` | Enter com campo vazio retorna nil; Enter com valor emite `pwdEnteredMsg`; Esc emite `flowCancelledMsg` |
| `passwordCreateModal` | Tab alterna entre campo senha e confirmação; mismatch mantém Enter bloqueado; match emite `pwdCreatedMsg` |
| `rootModel` | Push/pop de modais; routing de `modalResult` para `activeFlow`; `startFlowMsg` limpa modais órfãos; `WindowSizeMsg` não propaga para modais |

Estes testes são implementados com `t.Run` table-driven onde a combinação de casos é alta (filepicker tem 18 casos behaviorais mapeando a um "D-07 matrix").

---

## Testes de Integração de Fluxo

`exit_flow_integration_test.go` contém testes que exercitam fluxos de saída completos sem golden files. Estes testes:

- Simulam sequências de mensagens no `rootModel`
- Verificam estados finais após múltiplos `Update()` encadeados
- Usam stubs de `vaultSaver` para injetar erros de I/O
- Verificam que `pushModalMsg` corretos são emitidos nos pontos certos do fluxo

A separação em arquivo dedicado sinaliza que estes são testes de **integração de lógica de orquestração**, não de aparência, e não dependem de renderização visual.

---

## Separação entre Testes Visuais e Comportamentais

A separação é física (por blocos comentados dentro dos arquivos) e semântica:

| Aspecto | Testes Visuais (golden) | Testes Comportamentais (unitários) |
|---|---|---|
| **O que verificam** | Output de `View()` | Output de `Update()` + estado interno |
| **Determinismo** | Dependem de estado fixo, SetSize determinístico | Dependem apenas de mensagens injetadas |
| **Regressões capturadas** | Espaçamento, cor, layout, wrapping | Lógica, roteamento, emissão de comandos |
| **Arquivos em disco** | Sim — `testdata/golden/` | Não |
| **Atualização** | `go test -update` | Sem necessidade de atualização |

Esta separação garante que uma falha de golden seja sempre uma questão visual (intencional ou não), jamais confundida com uma falha lógica.

---

## Decisões Arquiteturais Relevantes

| Decisão | Consequência |
|---|---|
| Dois arquivos golden por cenário (`.txt` + `.json`) | Cada nível captura uma dimensão diferente; nenhum falso positivo entre estrutura e estilo |
| Parser ANSI como subpacote Go (não script externo) | Testável, importável, sem dependência de ferramentas externas |
| Canonicidade do parser (estado visual, não SGR) | Refatorações que preservam aparência não invalidam golden files |
| Auto-criação na primeira execução | Zero fricção para adicionar novos cenários; sem falhas de CI em novos componentes |
| Flag `-update` compartilhado via variável de pacote | Único ponto de controle; regeneração coordenada de todos os golden do pacote |
| Fixtures nomeadas por cenário (`pocKey1..F`, `newDestructive1Short(w)`) | Cenários de decision dialog legíveis sem necessidade de comentários; largura parametrizável para golden |
| Helpers de path específicos por componente (`decisionGoldenPath`, `helpGoldenPath`) | Variantes com dimensões embutidas no nome não precisam de campo de tamanho separado |
| Parser ANSI reutilizado para asserções programáticas (D-PE-03/04) | Testes de estado condicional (`bold` quando ativo, sem `bold` quando inativo) sem golden por estado |
| Nomenclatura `{component}-{variant}-{size}` | Diff legível no git; identificação imediata do cenário por nome de arquivo |
| Testes de integração em arquivo dedicado sem golden | Sinaliza visualmente que fluxo de saída é lógica, não aparência |
