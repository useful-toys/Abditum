# TUI Root Model Design

## Visão Geral

O TUI Abditum usa uma arquitetura onde o RootModel é o orquestrador central que coordena Views isoladas e independentes. Views não acessam estado de outras views ou do root — comunicação via métodos diretos ou eventos quando necessário.

---

## Conceitos Centrais

### Frame (Tela Principal)

A tela principal é dividida em:
- **Header** (2 linhas fixas)
- **Work Area** (área de trabalho)
- **Message Bar** (1 linha)
- **Action Bar** (1 linha)

### Work Area States

| WorkArea | childViews |
|---------|-----------|
| WorkAreaWelcome | Welcome |
| WorkAreaSettings | Settings |
| WorkAreaVault | Tree + Detail (split) |
| WorkAreaTemplates | List + Detail (split) |

---

## Interfaces

### ChildView

- Componente da tela principal
- Coexiste em paralelo com outras ChildViews na mesma Work Area
- Espaço fixo dedicado pelo root (não negociado)
- Referenciado diretamente via ponteiro

```go
type ChildView interface {
    Render(height, width int, theme Theme) string
    HandleKey(msg tea.KeyMsg) tea.Cmd
    HandleEvent(event any)
    HandleTeaMsg(msg tea.Msg)
}
```

### ModalView

- Stack sobre a tela principal
- Só interage com quem a apresentou (parent)
- Espaço máximo, centralizado pelo root

```go
type ModalView interface {
    Render(maxHeight, maxWidth int, theme Theme) string
    HandleKey(msg tea.KeyMsg) tea.Cmd
}
```

---

## Modal System

### Intent

```go
type Intent int

const (
    IntentConfirm Intent = iota  // ação principal
    IntentCancel                // cancelar
    IntentOther                // alternativa
)
```

### ModalOption

```go
type ModalOption struct {
    Keys   []string      // teclas que disparam (ex: []string{"s", "Enter"})
    Label  string       // rótulo (ex: "Salvar")
    Intent Intent      // intenção
    Action func() tea.Cmd  // executado pelo HandleKey para decidir cmd
}
```

### Fluxo de Modal

1. Root push modal na stack
2. Root renderiza modal (passa dimensões)
3. Root passa eventos para modal no topo
4. Modal retorna Cmd se quiser — root executa
5. Cmd produz `ModalSubmit{Intent, Data}` se precisar

O root é transparente — não sabe como modal funciona internamente.

---

## RootModel

### Estado

```go
type RootModel struct {
    // Dimensões
    width  int
    height int

    // Estado da aplicação
    workArea     WorkArea // qual WorkArea está ativa
    focusedChild ChildView // ponteiro para ChildView com foco (input)

    // Tema
    theme *Theme

    // Vault
    vaultManager *vault.Manager

    // Views (tipos dos subpackages)
    welcome      *welcome.Welcome    // de tui/welcome
    settings    *settings.Settings // de tui/settings
    vaultTree    *secret.Tree      // de tui/secret
    secretDetail *secret.Detail    // de tui/secret
    templateList *template.List     // de tui/template
    templateDetail *template.Detail // de tui/template

    // Modals stack
    modals []ModalView

    // Timers
    lastActionAt time.Time
}
```

### Responsabilidades

1. **Dimensões**
   - Armazena width/height do terminal
   - Se dimensões desconhecidas, exibe "Aguarde..."

2. **Coordenação de Render**
   ```
   Render():
     se width == 0 || height == 0:
         → "Aguarde..."
     senão:
         → render all ChildViews → render modalStack
   ```

3. **Dispatch de Teclas**
   ```
   HandleKey(tea.KeyMsg):
     se modalStack não vazia:
         → modal no topo .HandleKey()
     senão se focusedChild != nil:
         → focusedChild.HandleKey()
     senão:
         → root trata
   ```

4. **Event Routing**
   | Mensagem | Destino |
   |---------|--------|
   | WindowSizeMsg | root |
   | domain events (secretAdded, etc) | todas as ChildViews |
   | tickMsg | todas as ChildViews |
   | ModalSubmit{Intent, Data} | root processa se preciso |

5. **Modal Submit**
   - Modal retorna Cmd que produz `ModalSubmit{Intent, Data}`
   - Root não conhece lógica interna do modal
   - Cmd pode ser nil (modal simples)
   - Root executa o Cmd retornado pelo modal transparentemente

6. **Ações (futuro)**
   - ActionManager armazenará actions globais, reusáveis, específicas
   - Verifica `enabled()` antes de executar
   - Por enquanto: não implementar

---

## Estrutura de Arquivos

```
tui/
├── view.go              # interfaces ChildView, ModalView + Theme + WorkArea
├── root.go             # RootModel
├── welcome/
│   └── welcome_view.go   # type Welcome
├── settings/
│   └── settings_view.go # type Settings
├── secret/
│   ├── tree_view.go    # type Tree
│   └── detail_view.go  # type Detail
├── template/
│   ├── list_view.go   # type List
│   └── detail_view.go # type Detail
└── modal/
    ├── modal.go       # Intent + ModalOption + baseModal
    ├── password.go   # passwordModal
    ├── confirm.go  # confirmModal
    ├── filepicker.go # filepickerModal
    └── help.go     # helpModal
```

---

## Keyboard Flow

```
tea.KeyPressMsg
    ↓
Root.HandleKey()
    ├── modalStack não vazia?
    │   └── yes: modal.top.HandleKey() → executa cmd
    ├── focusedChild != nil?
    │   └── yes: child.HandleKey() → executa cmd
    └── root trata
```

---

## Exemplos de WorkArea

### WorkAreaWelcome

```
┌────────────────────────────────┐
│ Header (2 linhas)               │
├────────────────────────────────┤
│                                │
│ welcomeView.Render(h, w, theme)  │
│                                │
├────────────────────────────────┤
│ msg bar (1 linha)              │
├────────────────────────────────┤
│ action bar (1 linha)            │
└────────────────────────────────┘
```

### WorkAreaVault

```
┌────────────────────────────────┐
│ Header (2 linhas)               │
├──────────────┬─────────────────┤
│ vaultTree  │ secretDetail    │
│ Render    │ Render        │
│ (w/2, h) │ (w-w/2, h)   │
├──────────────┴─────────────────┤
│ msg bar                      │
├─────────────────────────── │
│ action bar                  │
└─────────────────────────── │
```

---

## Pendentes

- [ ] ActionManager (não implementar agora)
- [ ] MessageManager (não implementar agora)
- [ ] Comunicação entre ChildViews (listeners)
- [ ] Event naming convention

---

## Decisões Tomadas

- ChildViews referenciadas diretamente via ponteiro (não usa ID)
- ModalView não guarda estado de root — comunicação via Cmd/msg
- Modal retorna Cmd opcional, root executa transparentemente
- ModalOption.Action() executado pelo HandleKey do modal para decidir cmd
- Subpackages organizados por domínio (welcome, settings, secret, template, modal)
- ActionManager e MessageManager: não implementar nesta fase