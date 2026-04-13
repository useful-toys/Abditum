# Arquitetura — Composição de Modelos da TUI (`internal/tui`)

> Descreve os padrões, convenções e decisões arquiteturais que governam a composição do `rootModel` com `childModel`, `modalView` e `flowHandler` na TUI do Abditum.

## Sumário

- [Visão Geral](#visão-geral)
- [O `rootModel` como Único `tea.Model`](#o-rootmodel-como-único-teamodel)
- [A Interface `childModel`](#a-interface-childmodel)
  - [Contrato de Dimensionamento](#contrato-de-dimensionamento)
  - [Campos Concretos, Nunca Interface](#campos-concretos-nunca-interface)
- [Máquina de Estados — `workArea`](#máquina-de-estados--workarea)
  - [`liveWorkChildren()` e `activeChild()`](#liveworkchildren-e-activechild)
  - [Transições de Área](#transições-de-área)
- [Despacho de `Update`](#despacho-de-update)
  - [Ordem de Prioridade em `KeyPressMsg`](#ordem-de-prioridade-em-keypressmsg)
  - [Broadcast para Eventos de Domínio](#broadcast-para-eventos-de-domínio)
- [Composição de `View` — Layout da Tela](#composição-de-view--layout-da-tela)
  - [Estrutura de Linhas](#estrutura-de-linhas)
  - [Layouts Lado a Lado](#layouts-lado-a-lado)
  - [Propagação de `SetSize`](#propagação-de-setsize)
- [A Stack de Modais](#a-stack-de-modais)
  - [A Interface `modalView`](#a-interface-modalview)
  - [Push e Pop](#push-e-pop)
  - [Bloqueio de Input](#bloqueio-de-input)
  - [Roteamento de `modalResult`](#roteamento-de-modalresult)
  - [Renderização — Apenas o Topo](#renderização--apenas-o-topo)
  - [Barra de Comandos Durante Modais](#barra-de-comandos-durante-modais)
- [Tipos de Modal](#tipos-de-modal)
  - [`modalModel` — Modal Genérico](#modalmodel--modal-genérico)
  - [`DecisionDialog` — Decisão com Severidade](#decisiondialog--decisão-com-severidade)
  - [`helpModal`, `filePickerModal`, Modais de Senha](#helpmodal-filepickermodal-modais-de-senha)
- [Flows — Fluxos de Tela Invisíveis](#flows--fluxos-de-tela-invisíveis)
  - [A Interface `flowHandler`](#a-interface-flowhandler)
  - [Ciclo de Vida de um Flow](#ciclo-de-vida-de-um-flow)
  - [Fast-path de CLI](#fast-path-de-cli)
- [Comunicação Entre Camadas](#comunicação-entre-camadas)
  - [Princípio da Mensagem Unidirecional](#princípio-da-mensagem-unidirecional)
  - [Fábricas de Comando em `mutations.go`](#fábricas-de-comando-em-mutationsgo)
  - [Serviços Compartilhados](#serviços-compartilhados)
- [O `ActionManager`](#o-actionmanager)
  - [Registro e Escopo](#registro-e-escopo)
  - [Despacho e Prioridade de Dono](#despacho-e-prioridade-de-dono)
  - [Barra de Comandos Visível](#barra-de-comandos-visível)
- [Decisões Arquiteturais Relevantes](#decisões-arquiteturais-relevantes)

## Visão Geral

A TUI do Abditum usa o framework Bubble Tea, que exige um único `tea.Model` raiz. Em vez de encadear diretamente modelos Bubble Tea, a arquitetura introduz três abstrações próprias — `childModel`, `modalView` e `flowHandler` — que o `rootModel` coordena como componentes internos.

O resultado é uma composição em três dimensões:

- **Horizontal (área de trabalho):** qual conjunto de `childModel`s ocupa a tela principal, determinado pela máquina de estados `workArea`.
- **Vertical (stack de modais):** quais `modalView`s estão empilhadas sobre a área de trabalho, recebendo input com prioridade.
- **Temporal (flow ativo):** qual `flowHandler` está orquestrando uma sequência de passos multi-etapa, empurrando e consumindo modais.

O `rootModel` é o único proprietário e árbitro dessas três dimensões.

## O `rootModel` como Único `tea.Model`

O `rootModel` é a única struct que implementa `tea.Model` no pacote. Toda a inicialização, dispatch de mensagens e renderização começa e termina nele.

Seus campos de estado principais são:

- `area workArea` — qual área de trabalho está ativa
- `width, height int` — dimensões atuais do terminal
- `theme *Theme` — tema visual ativo
- `mgr *vault.Manager` — gerenciador de cofre (nil enquanto nenhum cofre está aberto)
- `isDirty bool` — se há modificações não salvas
- `modals []modalView` — stack LIFO de modais ativos
- `activeFlow flowHandler` — flow em execução, ou nil
- `actions *ActionManager` — gerenciador de ações de teclado
- `messages *MessageManager` — gerenciador de mensagens na barra de status

Os campos dos filhos são ponteiros concretos para cada modelo de área de trabalho (`welcome`, `vaultTree`, `secretDetail`, `templateList`, `templateDetail`, `settings`). Apenas os ponteiros não-nil correspondentes à área ativa serão alocados em um dado momento.

## A Interface `childModel`

```go
type childModel interface {
    Update(tea.Msg) tea.Cmd
    View() string
    SetSize(width, height int)
    ApplyTheme(*Theme)
}
```

Todo modelo de área de trabalho implementa este contrato. `ApplyTheme` é obrigatório — não é opcional como em `modalView`. Isso garante que o compilador force todos os filhos a responder à alternância de tema, sem possibilidade de esquecimento silencioso.

### Contrato de Dimensionamento

`SetSize` é sempre chamado pelo `rootModel` imediatamente antes da renderização. O filho recebe as dimensões exatas que deve preencher — não limites máximos. Se `View()` for chamado com dimensões zero, o comportamento é indefinido (crash ou layout corrompido).

### Campos Concretos, Nunca Interface

Os campos do `rootModel` que armazenam filhos são declarados como ponteiros concretos (`*welcomeModel`, `*vaultTreeModel`, etc.), nunca como `childModel`. Isso é uma decisão de segurança de tipos: um `childModel((*welcomeModel)(nil))` seria não-nil em Go, quebrando verificações `if m.welcome != nil`. Com ponteiros concretos, a verificação de nil é sempre confiável.

## Máquina de Estados — `workArea`

A variável `area workArea` (tipo iota) controla qual conjunto de filhos está visível e recebe input:

| Estado | Filhos renderizados | Uso |
|---|---|---|
| `workAreaWelcome` | `welcome` (tela inteira) | Tela inicial, sem cofre aberto |
| `workAreaVault` | `vaultTree` (esquerda) + `secretDetail` (direita) | Cofre aberto |
| `workAreaTemplates` | `templateList` (esquerda) + `templateDetail` (direita) | Navegação de templates |
| `workAreaSettings` | `settings` (tela inteira) | Configurações |

Somente o `rootModel` altera `area`. Filhos nunca transicionam diretamente entre áreas — emitem mensagens que o `rootModel` interpreta para realizar a transição.

### `liveWorkChildren()` e `activeChild()`

`liveWorkChildren()` retorna todos os ponteiros de filhos não-nil como uma slice de `childModel`. É usado para propagar `SetSize`, `ApplyTheme` e mensagens de broadcast para todos os filhos existentes, independente de qual está ativo.

`activeChild()` retorna o filho correspondente ao `area` atual. É o único filho que recebe input de teclado no despacho de `Update`.

A distinção é importante: `liveWorkChildren()` garante que filhos em background (como `secretDetail` quando `vaultTree` tem o foco) ainda recebem eventos de domínio e atualizações de tema.

### Transições de Área

Quando o cofre é aberto (`vaultOpenedMsg`), o `rootModel` executa `enterVault()`:

- `m.area = workAreaVault`
- `m.welcome = nil` — permite coleta de lixo
- Aloca e dimensiona `m.vaultTree` e `m.secretDetail`

A transição inversa (fechar cofre) zera os ponteiros de vault e recria `welcome`. Em nenhum momento dois conjuntos de filhos de área diferente coexistem além do tempo necessário para a transição.

## Despacho de `Update`

### Ordem de Prioridade em `KeyPressMsg`

O `Update` do `rootModel` processa `tea.KeyPressMsg` em uma cadeia de prioridade estrita:

1. **Trava de emergência** (`ctrl+alt+shift+q`) — sempre primeiro, wipe de memória imediato
2. **`ctrl+q` direto** — inicia `saveAndExitFlow` ou mostra `DecisionDialog` conforme estado
3. **`f12` direto** — alterna o tema, chama `applyTheme()`
4. **Help modal prioritário** — se o topo da stack for `*helpModal`, ele consome toda a entrada (incluindo ESC e F1)
5. **`ActionManager.Dispatch`** — tenta `ScopeGlobal` sempre; tenta `ScopeLocal` apenas se `!inFlowOrModal`
6. **Topo da stack de modais** — `m.modals[len-1].Update(msg)`
7. **Flow ativo** — `m.activeFlow.Update(msg)` (fallback quando modal não consumiu)
8. **Filho ativo** — `m.activeChild().Update(msg)`

Cada etapa pode interceptar a mensagem e retornar sem avançar para a próxima. As etapas 1–5 têm precedência incondicional sobre modais e filhos.

### Broadcast para Eventos de Domínio

Mensagens que não são `KeyPressMsg` — eventos de domínio como `secretAddedMsg`, `vaultSavedMsg`, `templateRenamedMsg` — são enviadas via `broadcast()` para **todos** os `liveWorkChildren()` e **todos** os modais simultaneamente. Nenhum filho de background fica desatualizado quando o estado do cofre muda.

O `rootModel` ainda processa o evento antes do broadcast para atualizar seu próprio estado (como `isDirty`, `vaultPath`).

## Composição de `View` — Layout da Tela

### Estrutura de Linhas

O frame de tela é montado por `renderFrame`, que recebe opcionalmente o modal do topo:

```
Linha 0–1:    header (2 linhas)
Linha 2–N:    área de trabalho (height − 4 linhas)
Linha N+1:    barra de mensagem (1 linha)
Linha N+2:    barra de comandos (1 linha)
```

O conteúdo da área de trabalho varia conforme `area`. Se um modal é passado, ele é sobreposto àrea de trabalho via `lipgloss.Place` com alinhamento centrado, deixando header e barras intactas.

O `rootModel.View()` retorna um `tea.View` com `AltScreen = true` e `BackgroundColor = m.theme.SurfaceBase`. Isso garante que o fundo do terminal coincide com a superfície base do tema ativo.

### Layouts Lado a Lado

Para `workAreaVault` e `workAreaTemplates`, dois filhos são renderizados lado a lado com `lipgloss.JoinHorizontal`. A divisão é exatamente `width/2` e `width - width/2`, garantindo preenchimento total sem gap de arredondamento.

Cada filho recebe seus próprios limites exatos via `SetSize` antes de renderizar.

### Propagação de `SetSize`

Em `tea.WindowSizeMsg`:
- `rootModel` atualiza `m.width` e `m.height`
- Chama `SetSize` em todos os `liveWorkChildren()` com as dimensões calculadas para cada região

Modais **não** são dimensionados em `WindowSizeMsg`. Eles recebem `SetAvailableSize(maxWidth, workH)` just-in-time dentro de `renderFrame`, imediatamente antes de `modal.View()` ser chamado. Isso significa que modais sempre têm as dimensões corretas do frame atual sem estado de tamanho desatualizado.

## A Stack de Modais

A stack de modais é o mecanismo de sobreposição de contexto da TUI. Enquanto houver pelo menos um modal na stack, o contexto de interação pertence ao modal do topo.

### A Interface `modalView`

```go
type modalView interface {
    Update(tea.Msg) tea.Cmd
    View() string
    Shortcuts() []Shortcut
    SetAvailableSize(maxWidth, maxHeight int)
}
```

Diferenças em relação a `childModel`:

- `Shortcuts()` retorna os atalhos de teclado específicos deste modal, exibidos na barra de comandos enquanto ele está no topo
- `SetAvailableSize` recebe limites **máximos**, não dimensões exatas — o modal pode renderizar menor e se autocentrar
- Não inclui `ApplyTheme` — o protocolo de tema é opcional para modais (ver [Arquitetura de Tema](arquitetura-tui-tema.md))

### Push e Pop

**Push:** qualquer componente retorna um `tea.Cmd` que emite `pushModalMsg{modal: m}`. O `rootModel` appenda o modal à slice. Se o modal implementa `Init() tea.Cmd`, `Init` é chamado imediatamente após o push.

**Pop:** o próprio modal emite `popModalMsg{}` como `tea.Cmd` quando seu ciclo de vida termina (confirmação, cancelamento, ESC). O `rootModel` trunca a slice em `len-1`. Como a emissão é assíncrona (via runtime do Bubble Tea), não há risco de pop durante o próprio Update do modal.

A stack pode ter múltiplos modais sobrepostos. Por exemplo, um `DecisionDialog` pode ser empurrado sobre um `filePickerModal` que já está na stack. Cada pop desvela o modal anterior, sem que o `rootModel` precise conhecer a sequência.

### Bloqueio de Input

O flag `inFlowOrModal` é calculado antes do despacho:

```go
inFlowOrModal := m.activeFlow != nil || len(m.modals) > 0
```

Quando `inFlowOrModal == true`, o `ActionManager.Dispatch` pula todas as ações com `ScopeLocal`. Somente ações `ScopeGlobal` (como F1 e F12) continuam disparando. Isso bloqueia automaticamente todas as ações de área de trabalho sem que cada modal precise implementar lógica de bloqueio própria.

### Roteamento de `modalResult`

Modais comunicam seu resultado ao flow orquestrador por meio de mensagens que implementam a interface marcadora `modalResult`. Exemplos: `passwordEntryResult`, `filePickerResult`, `overwriteConfirmedMsg`.

O `rootModel` identifica `modalResult` no `Update` e redireciona exclusivamente para `m.activeFlow.Update(msg)`. Filhos de área de trabalho nunca recebem resultados de modal. Essa separação garante que o flow seja o único responsável por interpretar e reagir ao resultado.

### Renderização — Apenas o Topo

Apenas `modals[len-1]` (o modal mais recente) é renderizado e recebe input. Modais abaixo do topo ficam suspensos — seus estados são preservados, mas não atualizam e não renderizam até voltarem ao topo.

### Barra de Comandos Durante Modais

Quando há um modal ativo, a barra de comandos exibe `modal.Shortcuts()` em vez das ações do `ActionManager`. Isso dá ao modal controle total sobre o contexto de ajuda visual sem depender do sistema de ações globais.

## Tipos de Modal

### `modalModel` — Modal Genérico

`modalModel` é um modal de propósito geral para diálogos informativos simples (mensagens, confirmações). Não tem lógica própria de resultado — usa callbacks passados na construção. É o tipo de modal mais simples da codebase.

### `DecisionDialog` — Decisão com Severidade

`DecisionDialog` é um modal especializado para diálogos de confirmação com impacto semântico. Sua severidade (`Neutral`, `Informative`, `Alert`, `Error`, `Destructive`) controla a cor da borda e do destaque das ações. Isso comunica visualmente o risco da operação sem texto adicional.

O `DecisionDialog` usa box drawing Unicode desenhado manualmente, não a API de bordas do lipgloss. Isso permite embutir a barra de ações na borda inferior do modal, formando uma unidade visual coesa. Cada tecla de ação (Enter, ESC, letras customizadas) é renderizada inline no rodapé da caixa.

Dois construtores cobrem os casos típicos: `Acknowledge` (uma ação de confirmação) e `Decision` (ação padrão + ações intermediárias + cancelamento). Ambos retornam `tea.Cmd` que empurra o modal na stack.

### `helpModal`, `filePickerModal`, Modais de Senha

`helpModal` é o único modal com tratamento especial no despacho — quando está no topo, ele intercepta toda entrada de teclado antes do `ActionManager`. Isso impede que atalhos globais como F12 disparem acidentalmente enquanto o usuário lê a ajuda.

`filePickerModal` é o único modal que implementa o protocolo opcional `ApplyTheme(*Theme)`, pois é instanciado uma única vez por flow e pode sobreviver à alternância de tema.

Os modais de senha (`passwordEntryModal`, `passwordCreateModal`) não recebem `*Theme` e usam constantes hardcoded de `tokens.go`. Essa inconsistência está documentada em [Arquitetura de Tema](arquitetura-tui-tema.md).

## Flows — Fluxos de Tela Invisíveis

### A Interface `flowHandler`

```go
type flowHandler interface {
    Init() tea.Cmd
    Update(tea.Msg) tea.Cmd
}
```

Flows não têm `View()`. Eles são invisíveis — toda a UI durante um flow é realizada por modais que o flow em si empurra na stack. O flow age como orquestrador de uma sequência de passos, não como renderizador.

### Ciclo de Vida de um Flow

O `rootModel` gerencia o ciclo de vida do flow ativo via mensagens:

- `startFlowMsg{flow}` — limpa `m.modals`, define `m.activeFlow = flow`, chama `flow.Init()`
- `endFlowMsg` — define `m.activeFlow = nil`

Apenas um flow pode estar ativo por vez. Os três flows concretos são:

**`openVaultFlow`**: orquestra a abertura de um cofre existente. Em sequência: exibe `FilePicker` → recebe `filePickerResult` → exibe `PasswordEntry` → recebe `passwordEntryResult` → executa carregamento em background → emite `vaultOpenedMsg`.

**`createVaultFlow`**: orquestra a criação de um novo cofre. Inclui um passo adicional de verificação de sobrescrita (`overwriteConfirmedMsg`/`overwriteCancelledMsg`) e um passo de avaliação da força da senha (`weakPwdProceedMsg`/`weakPwdReviseMsg`).

**`saveAndExitFlow`**: orquestra o salvamento e saída quando há modificações. Verifica modificação externa do arquivo em background antes de prosseguir. É o único flow que não começa por escolha do usuário, mas é acionado diretamente pelo handler de `ctrl+q` no `rootModel` quando `isDirty == true`.

Flows registram suas próprias ações no `ActionManager` enquanto ativos e as limpam ao encerrar (`actions.ClearOwned(f)`).

### Fast-path de CLI

Quando o programa é iniciado com um caminho de arquivo como argumento, `openVaultFlow` verifica se `cliPath` está definido. Em caso positivo, pula o passo de `FilePicker` e vai diretamente para `PasswordEntry`. Isso é implementado dentro do próprio flow, sem lógica especial no `rootModel`.

## Comunicação Entre Camadas

### Princípio da Mensagem Unidirecional

Nenhum filho, modal ou flow chama métodos diretamente no `rootModel`. Toda comunicação é unidirecional, via `tea.Cmd` retornando `tea.Msg`:

```
componente.Update(msg) → tea.Cmd → tea.Msg → rootModel.Update(msg)
```

Isso garante que o `rootModel` seja o único ponto de mutação de estado de nível superior. Filhos não têm referência ao `rootModel`.

### Fábricas de Comando em `mutations.go`

`mutations.go` centraliza fábricas de `tea.Cmd` para operações simples sobre o cofre:

```go
func softDeleteSecretCmd(id string) tea.Cmd {
    return func() tea.Msg { return secretDeletedMsg{id: id} }
}
```

Filhos chamam essas fábricas e retornam o `tea.Cmd` resultante de seus próprios `Update`. O runtime do Bubble Tea executa o comando, emite a mensagem, e o `rootModel` a recebe para atualização de estado e broadcast.

Operações de I/O pesado (carregar vault, salvar vault) usam o mesmo padrão mas com execução real em goroutine implícita do runtime.

### Serviços Compartilhados

`ActionManager` e `MessageManager` são passados a todos os filhos e flows como referências compartilhadas no momento da construção. Filhos os usam diretamente — não há bus de eventos intermediário. `MessageManager.Show()` e `ActionManager.Register()` são chamadas diretas via closures capturadas.

## O `ActionManager`

### Registro e Escopo

Qualquer struct (rootModel, filho, flow) registra ações passando a si mesmo como owner:

```go
actions.Register(owner, Action{Keys: ["f5"], Scope: ScopeLocal, ...})
```

Cada `Action` tem:
- `Keys []string` — teclas gatilho
- `Scope ActionScope` — `ScopeLocal` (bloqueado durante flow/modal) ou `ScopeGlobal` (sempre)
- `Enabled func() bool` — closure dinâmica verificada no momento do dispatch
- `Handler func() tea.Cmd` — closure que retorna o comando a executar
- `Priority int` — posição na barra de comandos (maior = mais à esquerda)
- `HideFromBar bool` — exibido apenas no modal de ajuda, não na barra de comandos

Closures em `Enabled` capturam o estado do owner diretamente. Por exemplo, as ações F5/F6 do `rootModel` têm `Enabled: func() bool { return m.area == workAreaWelcome }`, o que as desabilita automaticamente quando um cofre está aberto, sem verificação adicional em nenhum outro lugar.

### Despacho e Prioridade de Dono

`Dispatch(key, inFlowOrModal)` percorre os owners em ordem de registro, dando prioridade ao `activeOwner`. Um filho pode chamar `actions.SetActiveOwner(self)` para que suas ações tenham prioridade sobre as do `rootModel` para a mesma tecla.

### Barra de Comandos Visível

`Visible()` retorna todas as ações com `Enabled() == true` e `HideFromBar == false`, ordenadas por `Priority` decrescente. A tecla F1 é ancorada à direita por convenção. A barra de comandos reflete automaticamente o contexto — um filho pode habilitar ou desabilitar ações dinamicamente via closures sem notificar o `rootModel`.

## Decisões Arquiteturais Relevantes

| Decisão | Consequência |
|---|---|
| `rootModel` é o único `tea.Model` | Sem composição de `tea.Model` aninhados; todo dispatch passa por um único ponto de controle |
| Campos de filhos como ponteiros concretos | Nil check confiável; ausência do "typed nil trap" de interfaces |
| `ApplyTheme` obrigatório em `childModel`, opcional em `modalView` | Filhos de área de trabalho sempre respondem à troca de tema; modais legados funcionam sem modificação |
| Stack de modais LIFO com somente o topo ativo | Filhos não precisam saber se há modal ativo; o bloqueio de input é centralizado via `inFlowOrModal` |
| `modalResult` roteado exclusivamente ao `activeFlow` | Flows são os únicos consumidores de resultados de modal; filhos nunca interferem |
| `SetAvailableSize` just-in-time para modais | Modais nunca ficam com dimensões stale após resize |
| Flows sem `View()` | Separação clara entre orquestração (flow) e apresentação (modal); flows são testáveis sem renderização |
| `ActionManager` centralizado com closures de `Enabled` | Barra de comandos e bloqueio de ações refletem estado real sem mensagens extras |
| Comunicação via `tea.Cmd`/`tea.Msg` unidirecional | Filhos e flows não têm dependência do `rootModel`; estado do nível superior só muta no `rootModel.Update` |
| `broadcast()` para eventos de domínio | Filhos de background mantêm estado coerente sem precisar de sincronização explícita |
