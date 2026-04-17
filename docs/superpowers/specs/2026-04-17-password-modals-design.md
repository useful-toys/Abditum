# Design — Modais de Senha (PasswordEntry e PasswordCreate)

**Data:** 2026-04-17  
**Status:** Aprovado para implementação  
**Relacionados:** `golden/tui-spec-dialog-senha.md`, `golden/tui-design-system.md`

---

## Contexto

O projeto não possui nenhum componente de entrada de texto na TUI. Este documento especifica a implementação dos dois modais de senha — `PasswordEntry` (abrir cofre) e `PasswordCreate` (criar/alterar senha mestra) — e as mudanças de infraestrutura necessárias para suportá-los.

---

## Decisões de Design

| Decisão | Escolha | Justificativa |
|---|---|---|
| Componente de input | Campo customizado `[]byte` + `crypto.Wipe` | `bubbles/textinput` não suporta máscara fixa nem wipe de memória; codebase é 100% custom |
| Máscara de senha | `•` por caractere digitado | Feedback natural; ocultar comprimento tem custo UX sem ganho de segurança proporcional |
| Cursor | Real (`tea.Cursor`) | Experiência nativa; requer extensão do pipeline no `RootModel` |
| Placeholder | Nenhum | Dica vem pela barra de mensagens — campo vazio com fundo `surface.input` |
| Contador de tentativas | Barra de mensagens (tipo Info) | Mantém altura fixa do modal; responsabilidade do orquestrador |
| Força de senha | Heurística simples (4 critérios) | Sem dependência externa; suficiente para feedback visual |

---

## Arquitetura

### Novos arquivos

```
internal/tui/modal/
├── password_field.go         # PasswordField — buffer seguro + renderização
├── password_entry_modal.go   # PasswordEntry — modal de 1 campo
├── password_create_modal.go  # PasswordCreate — modal de 2 campos
└── password_strength.go      # StrengthScore — heurística de força de senha
```

### Arquivos modificados

```
internal/tui/modal.go                 # Interface ModalView: +Cursor()
internal/tui/root.go                  # RootModel.View(): propagar cursor do modal ativo
internal/tui/modal/confirm_modal.go   # +Cursor() retorna nil
internal/tui/modal/help_modal.go      # +Cursor() retorna nil
```

---

## Componente: `PasswordField`

Struct com responsabilidade única: gerenciar o buffer de uma senha e renderizar as duas linhas do campo (label + área digitável).

### Estrutura

```go
type PasswordField struct {
    label   string  // "Senha", "Nova senha", "Confirmação"
    value   []byte  // buffer real — nunca exposto como string
}
```

### API pública

```go
func NewPasswordField(label string) *PasswordField
func (f *PasswordField) Value() []byte          // retorna cópia — nunca o slice interno
func (f *PasswordField) Len() int               // comprimento atual
func (f *PasswordField) Clear()                 // zera e realoca buffer
func (f *PasswordField) Wipe()                  // internal/crypto.Wipe(value) + Clear
func (f *PasswordField) HandleKey(msg tea.KeyMsg) bool  // true se consumiu a tecla
func (f *PasswordField) Render(innerWidth int, focused bool, theme *design.Theme) string
```

### Teclas tratadas pelo `HandleKey`

| Tecla | Efeito |
|---|---|
| Rune imprimível | Append ao buffer |
| `Backspace` | Remove último byte |
| Qualquer outra (`←` `→` `Del` `Home` `End` etc.) | Não consome (retorna false) |

> **Decisão de design:** navegação interna (←→, Home, End, Delete) não é suportada. O cursor é sempre posicionado ao final do conteúdo. Para corrigir um erro, o usuário usa `Backspace` para apagar caractere a caractere. Isso é o comportamento padrão em campos de senha — simples, previsível, e evita que a posição do cursor revele informação sobre o comprimento da senha.

### Renderização

`Render` retorna **duas linhas** separadas por `\n`:

**Linha 1 — label:**
- Focado: `accent.primary` bold
- Inativo: `text.secondary`

**Linha 2 — área digitável** (largura = `innerWidth`, fundo `surface.input`):
- Vazio: espaços brancos (fundo `surface.input` visível)
- Com conteúdo: `f.Len()` bullets `•` em `text.secondary` + espaços até a largura total

O modal **não usa o retorno de `Render`** para calcular a posição do cursor — usa `f.Len()` diretamente.

**Cálculo de `innerWidth`:** o modal passa `modalWidth - 2 - 2*design.DialogPaddingH` ao chamar `Render`. Para modal de 50 colunas: `50 - 2 - 4 = 44`. Os golden files de `PasswordField` usam `innerWidth = 44` — mas são testados isoladamente com tamanho `46x2` (44 de conteúdo + 2 de padding lateral somados pelo test harness para simular o contexto do modal).

### Wireframe do campo

```
  Nova senha                    ← label (accent.primary bold se focado)
  •••••••••••••                 ← bullets + espaços (fundo surface.input)
                                  ↑ cursor real aqui
```

---

## Modal: `PasswordEntry`

**Título:** `Senha mestra`  
**Borda:** `border.focused`  
**Largura:** 50 colunas  
**Altura:** fixa (5 linhas de corpo + 2 bordas = 7 linhas totais)

### Wireframes

**Campo vazio (ação default bloqueada):**
```
╭── Senha mestra ────────────────────────────╮
│                                            │
│  Senha                                     │
│                                            │
│                                            │
╰── Enter Confirmar ──────────── Esc Cancelar ──╯
         ↑ text.disabled
```

**Com conteúdo (ação default ativa):**
```
╭── Senha mestra ────────────────────────────╮
│                                            │
│  Senha                                     │
│  ••••••••••••                              │
│                                            │
╰── Enter Confirmar ──────────── Esc Cancelar ──╯
         ↑ accent.primary bold
```

### Estados

| Componente | Estado | Condição |
|---|---|---|
| Campo `Senha` | sempre visível, sempre focado | campo único |
| Ação `Enter Confirmar` | bloqueada `text.disabled` | campo vazio |
| Ação `Enter Confirmar` | ativa `accent.primary` bold | campo não vazio |
| Ação `Esc Cancelar` | sempre ativa | — |

### Mensagens na barra

| Contexto | Tipo | Texto |
|---|---|---|
| Diálogo abre | Dica | `• Digite a senha para desbloquear o cofre` |
| Tentativa ≥ 2 (emitida pelo orquestrador) | Info | `ℹ Tentativa N de M` |
| Diálogo fecha | — | Barra limpa *(orquestrador assume)* |

### Construtor e callbacks

```go
func NewPasswordEntryModal(
    mc        tui.MessageController,
    onConfirm func(password []byte) tea.Cmd,
    onCancel  func() tea.Cmd,
) *PasswordEntryModal
```

- `mc` — injetado para que o modal emita dicas e erros diretamente na barra de status via `mc.SetHint(...)` / `mc.SetError(...)`.
- `onConfirm` — chamado com `field.Value()` quando o usuário pressiona `Enter`. O modal **não sabe** se a senha foi correta; o orquestrador decide fechar, incrementar tentativas, ou fechar por esgotamento.
- `onCancel` — chamado quando o usuário pressiona `Esc`.

O modal emite a dica inicial (`• Digite a senha para desbloquear o cofre`) via `mc` ao ser criado.

Para indicar senha incorreta, o orquestrador chama:

```go
func (m *PasswordEntryModal) NotifyWrongPassword()  // limpa campo; modal não emite mensagem
```

A mensagem `ℹ Tentativa N de M` é emitida pelo orquestrador (não pelo modal).

### Teclado

| Tecla | Efeito | Condição |
|---|---|---|
| `Enter` | Chama `onConfirm(field.Value())` | Campo não vazio |
| `Esc` | `field.Wipe()` + `onCancel()` | — |
| `Tab` | No-op | Campo único |

---

## Modal: `PasswordCreate`

**Título:** `Definir senha mestra`  
**Borda:** `border.focused`  
**Largura:** 50 colunas  
**Altura:** variável (medidor aparece/desaparece conforme conteúdo de `Nova senha`)

### Wireframes

**Estado inicial — foco em "Nova senha", campos vazios:**
```
╭── Definir senha mestra ───────────────────╮
│                                           │
│  Nova senha                               │
│                                           │
│                                           │
│  Confirmação                              │
│                                           │
│                                           │
╰── Enter Confirmar ──────────── Esc Cancelar ──╯
         ↑ text.disabled
```

**Digitando em "Nova senha" — medidor aparece:**
```
╭── Definir senha mestra ───────────────────╮
│                                           │
│  Nova senha                               │
│  ••••••••••••                             │
│                                           │
│  Confirmação                              │
│                                           │
│                                           │
│  Força: ████████░░ Boa                    │
│                                           │
╰── Enter Confirmar ──────────── Esc Cancelar ──╯
         ↑ text.disabled (Confirmação vazia)
```
*(medidor: `████` = semantic.success ou warning; `░░` = text.disabled)*

**Ambos preenchidos, senhas conferem — foco em "Confirmação":**
```
╭── Definir senha mestra ───────────────────╮
│                                           │
│  Nova senha                               │
│  ••••••••••••                             │
│                                           │
│  Confirmação                              │
│  ••••••••••••                             │
│                                           │
│  Força: ████████░░ Boa                    │
│                                           │
╰── Enter Confirmar ──────────── Esc Cancelar ──╯
         ↑ accent.primary bold
```

### Estados

| Componente | Estado | Condição |
|---|---|---|
| Campo `Nova senha` | sempre visível | — |
| Campo `Confirmação` | sempre visível | — |
| Medidor de força | visível (+ linha em branco antes) | `Nova senha` não vazio |
| Medidor de força | oculto | `Nova senha` vazio |
| Ação `Enter Confirmar` | bloqueada | qualquer campo vazio **ou** senhas divergentes |
| Ação `Enter Confirmar` | ativa | ambos não vazios **e** `bytes.Equal(fieldNew.Value(), fieldConfirm.Value())` |
| Ação `Esc Cancelar` | sempre ativa | — |

### Mensagens na barra

| Contexto | Tipo | Texto |
|---|---|---|
| Diálogo abre / foco em `Nova senha` | Dica | `• A senha mestra protege todo o cofre — use 12+ caracteres` |
| Foco em `Confirmação` (senhas iguais ou vazia) | Dica | `• Redigite a senha para confirmar` |
| Foco em `Confirmação` (senhas divergentes) | Erro 5s | `✕ As senhas não conferem — digite novamente` |
| Digitação em `Confirmação` → senhas divergem | Erro 5s | `✕ As senhas não conferem — digite novamente` |
| Digitação em `Confirmação` → senhas conferem | Dica | `• Redigite a senha para confirmar` |
| `Enter` → senhas divergentes | Erro 5s | `✕ As senhas não conferem — digite novamente` |
| Diálogo fecha | — | Barra limpa *(orquestrador assume)* |

### Construtor e callbacks

```go
func NewPasswordCreateModal(
    mc        tui.MessageController,
    onConfirm func(password []byte) tea.Cmd,
    onCancel  func() tea.Cmd,
) *PasswordCreateModal
```

- `mc` — injetado para emitir dicas e erros na barra de status.
- `onConfirm` — chamado com `fieldNew.Value()` quando ambos os campos estão preenchidos, iguais e o usuário pressiona `Enter`.
- `onCancel` — chamado quando o usuário pressiona `Esc`.

O modal emite a dica inicial (`• A senha mestra protege todo o cofre — use 12+ caracteres`) via `mc` ao ser criado.

### Teclado

| Tecla | Efeito | Condição |
|---|---|---|
| `Tab` | Alterna foco `Nova senha` ↔ `Confirmação` | — |
| `Enter` | `onConfirm(fieldNew.Value())` + wipe de ambos | Ambos não vazios e iguais |
| `Esc` | Wipe de ambos + `onCancel()` | — |

### Validação em tempo real

A cada tecla em `Confirmação`:
- `bytes.Equal(fieldNew.Value(), fieldConfirm.Value())` → ativa ação default + dica na barra
- Divergentes → bloqueia ação default + erro na barra (TTL 5s)

Ao abandonar `Confirmação` via `Tab` com senhas divergentes: foco volta para `Nova senha`, erro exibido na barra.

---

## Componente: `StrengthScore` (`password_strength.go`)

Função pura — zero dependências externas.

### Heurística (4 critérios, 1 ponto cada)

| Critério | Condição |
|---|---|
| Comprimento | `len(password) >= 12` |
| Maiúscula | contém pelo menos 1 rune `unicode.IsUpper` |
| Número | contém pelo menos 1 rune `unicode.IsDigit` |
| Símbolo | contém pelo menos 1 rune `!@#$%^&*()-_=+[]{}|;:,.<>?/~` |

### Níveis

| Pontos | Nível | Label | Cor do label | Cor da barra preenchida |
|---|---|---|---|---|
| 0–1 | Fraca | `⚠ Fraca` | `semantic.warning` | `semantic.warning` |
| 2–3 | Boa | `Boa` | `semantic.success` | `semantic.success` |
| 4 | Forte | `✓ Forte` | `semantic.success` | `semantic.success` |

### Barra de progresso

10 blocos Unicode `█` (preenchido) e espaço com fundo `text.disabled` (vazio).  
Proporção: `pontos / 4 * 10` blocos preenchidos (arredondado).

```
Força: ██████████ ✓ Forte    ← 4 pontos
Força: ████████░░ Boa        ← 3 pontos  
Força: █████░░░░░ Boa        ← 2 pontos
Força: ██░░░░░░░░ ⚠ Fraca    ← 1 ponto
```

### API

```go
type StrengthLevel int

const (
    StrengthWeak   StrengthLevel = iota // 0 ou 1 ponto — mesma apresentação visual
    StrengthFair                        // 2 ou 3 pontos
    StrengthStrong                      // 4 pontos
)

func EvaluateStrength(password []byte) StrengthLevel
func RenderStrengthMeter(password []byte, innerWidth int, theme *design.Theme) string
```

> **Nota:** 0 pontos e 1 ponto resultam no mesmo `StrengthWeak` — não há distinção visual entre eles. `EvaluateStrength` com senha vazia retorna `StrengthWeak`, mas `RenderStrengthMeter` nunca é chamada com senha vazia (o medidor só aparece quando `Nova senha` não é vazio).

---

## Extensão do Pipeline de Cursor

### Interface `ModalView` (tui/modal.go)

Adicionar método:

```go
type ModalView interface {
    Render(maxHeight, maxWidth int, theme *design.Theme) string
    HandleKey(msg tea.KeyMsg) tea.Cmd
    Update(msg tea.Msg) tea.Cmd
    Cursor(modalTopY, modalLeftX int) *tea.Cursor  // novo
}
```

`ConfirmModal` e `HelpModal` implementam retornando `nil`.

### `RootModel.View()` (tui/root.go)

Após renderizar o modal e compor o conteúdo, medir o conteúdo renderizado e consultar o cursor:

```go
if len(r.modals) > 0 {
    top := r.modals[len(r.modals)-1]
    modalH := r.height - 2
    modalContent := top.Render(modalH, r.width, r.theme)
    centeredModal := lipgloss.Place(r.width, modalH, lipgloss.Center, lipgloss.Center, modalContent)

    // Calcular posição do modal na tela para offset do cursor.
    modalW, modalActualH := lipgloss.Size(modalContent)
    topY := 1 + (modalH-modalActualH)/2  // 1 = offset Y do layer compositor
    leftX := (r.width - modalW) / 2

    if c := top.Cursor(topY, leftX); c != nil {
        v.Cursor = c
    }

    result := lipgloss.NewCompositor(
        lipgloss.NewLayer(base),
        lipgloss.NewLayer(centeredModal).Y(1).Z(1),
    ).Render()
    // ...
}
```

> **Nota:** `modalScreenPosition` não é uma função separada — o cálculo é feito inline em `View()` usando `lipgloss.Size(modalContent)` após renderizar, evitando duplicação de lógica.

### Cálculo de posição no modal

Cada modal de senha implementa `Cursor(topY, leftX int) *tea.Cursor` calculando a linha exata do campo focado dentro do body.

A fórmula geral:
```
cursorX = leftX + 1 (borda esquerda) + DialogPaddingH + f.Len()
cursorY = topY  + 1 (borda superior) + linhaDoFieldNoBody
```

`DialogPaddingH = 2` (constante em `design`).

**`PasswordEntry` — mapa de linhas do body (0-indexed):**

```
Linha 0: (vazia — padding)
Linha 1: label "Senha"
Linha 2: área digitável  ← linhaDoFieldNoBody = 2
Linha 3: (vazia — padding)
Linha 4: (vazia — padding inferior)
```

**`PasswordCreate` — mapa de linhas do body (0-indexed, sem medidor):**

```
Linha 0: (vazia)
Linha 1: label "Nova senha"
Linha 2: área digitável Nova senha   ← linhaDoFieldNoBody("Nova senha") = 2
Linha 3: (vazia)
Linha 4: (vazia)
Linha 5: label "Confirmação"
Linha 6: área digitável Confirmação  ← linhaDoFieldNoBody("Confirmação") = 6
Linha 7: (vazia)
```

**`PasswordCreate` — com medidor (body cresce 2 linhas ao final):**

```
Linha 0–7: igual ao sem medidor
Linha 8: medidor de força
Linha 9: (vazia após medidor)
```

As linhas de `Nova senha` (2) e `Confirmação` (6) **não mudam** quando o medidor aparece — ele é adicionado ao final. O modal implementa:

```go
func (m *PasswordCreateModal) Cursor(topY, leftX int) *tea.Cursor {
    if m.focused == fieldNew {
        y := topY + 1 + 2  // borda + linha 2
        x := leftX + 1 + design.DialogPaddingH + m.fieldNew.Len()
        return &tea.Cursor{X: x, Y: y}
    }
    // focused == fieldConfirm
    y := topY + 1 + 6  // borda + linha 6
    x := leftX + 1 + design.DialogPaddingH + m.fieldConfirm.Len()
    return &tea.Cursor{X: x, Y: y}
}
```

`Render` não retorna posição de cursor — sua responsabilidade é exclusivamente renderizar o conteúdo.

---

## Segurança de Memória

- `PasswordField.value` é `[]byte` — nunca convertido para `string` internamente
- `Value()` retorna **cópia** do slice — o caller é responsável pelo wipe após uso
- `Wipe()` chama `internal/crypto.Wipe(f.value)` (pacote `github.com/useful-toys/abditum/internal/crypto`) e redefine o slice
- `Esc` em qualquer modal chama `Wipe()` em todos os campos antes de fechar
- `Enter` em `PasswordCreate` faz wipe de `fieldConfirm` imediatamente após comparação; wipe de `fieldNew` é responsabilidade do orquestrador após uso do valor

---

## Testes com Golden Files

Cada componente e cada modal deve ter testes de renderização usando o framework de golden files do projeto (`internal/tui/testdata`). Os testes verificam dois aspectos simultaneamente:

- **`.golden.txt`** — texto visível sem ANSI: verifica alinhamento de bordas, título, conteúdo e ações
- **`.golden.json`** — transições de estilo: verifica cores, bold, italic de cada elemento

### Localização dos arquivos

```
internal/tui/modal/testdata/golden/
├── password_field-empty_focused-44x2.golden.txt
├── password_field-empty_focused-44x2.golden.json
├── password_field-empty_blurred-44x2.golden.txt
├── password_field-empty_blurred-44x2.golden.json
├── password_field-content_focused-44x2.golden.txt
├── password_field-content_focused-44x2.golden.json
├── password_field-content_blurred-44x2.golden.txt
├── password_field-content_blurred-44x2.golden.json
├── password_entry-empty-50x7.golden.txt
├── password_entry-empty-50x7.golden.json
├── password_entry-with_content-50x7.golden.txt
├── password_entry-with_content-50x7.golden.json
├── password_create-initial-50x9.golden.txt
├── password_create-initial-50x9.golden.json
├── password_create-with_meter-50x11.golden.txt
├── password_create-with_meter-50x11.golden.json
├── password_create-confirmed-50x11.golden.txt
├── password_create-confirmed-50x11.golden.json
├── password_create-mismatch-50x11.golden.txt
├── password_create-mismatch-50x11.golden.json
├── strength_meter-weak-44x1.golden.txt
├── strength_meter-weak-44x1.golden.json
├── strength_meter-fair-44x1.golden.txt
├── strength_meter-fair-44x1.golden.json
├── strength_meter-strong-44x1.golden.txt
└── strength_meter-strong-44x1.golden.json
```

### Variantes a cobrir

**`PasswordField`**

| Variante | Condição | Tamanho |
|---|---|---|
| `empty_focused` | Campo vazio, focado | `44x2` |
| `empty_blurred` | Campo vazio, sem foco | `44x2` |
| `content_focused` | 8 caracteres, focado | `44x2` |
| `content_blurred` | 8 caracteres, sem foco | `44x2` |

**`PasswordEntry`**

| Variante | Condição | Tamanho |
|---|---|---|
| `empty` | Campo vazio, ação bloqueada | `50x7` |
| `with_content` | Campo com conteúdo, ação ativa | `50x7` |

**`PasswordCreate`**

| Variante | Condição | Tamanho |
|---|---|---|
| `initial` | Ambos vazios, sem medidor | `50x9` |
| `with_meter` | `Nova senha` preenchida, medidor visível, `Confirmação` vazia | `50x11` |
| `confirmed` | Ambos preenchidos, senhas iguais, ação ativa | `50x11` |
| `mismatch` | Ambos preenchidos, senhas divergentes, ação bloqueada | `50x11` |

**`StrengthMeter`**

| Variante | Condição | Tamanho |
|---|---|---|
| `weak` | 0–1 pontos | `44x1` |
| `fair` | 2–3 pontos | `44x1` |
| `strong` | 4 pontos | `44x1` |

### Padrão de teste

```go
func TestPasswordEntryModal_Empty(t *testing.T) {
    mc := &stubMessageController{}
    m := modal.NewPasswordEntryModal(mc,
        func(_ []byte) tea.Cmd { return nil },
        func() tea.Cmd { return nil },
    )
    testdata.TestRenderManaged(t, "password_entry", "empty", []string{"50x7"},
        func(w, h int, theme *design.Theme) string {
            return m.Render(h, w, theme)
        })
}
```

`stubMessageController` é um struct local nos testes que implementa `tui.MessageController` com métodos no-op.

### O que a IA deve verificar nos golden files

Após gerar os golden files com `-update-golden`, a IA **deve ler os arquivos `.golden.txt`** e verificar visualmente:

1. **Bordas** — `╭` `╮` `╰` `╯` `│` `─` alinhados corretamente; largura total = 50 colunas
2. **Título** — `Senha mestra` ou `Definir senha mestra` centralizado entre os `──` da borda superior
3. **Labels dos campos** — `Senha`, `Nova senha`, `Confirmação` com indentação de 2 espaços (`DialogPaddingH`)
4. **Área digitável** — largura correta (innerWidth = 44 colunas para modal de 50)
5. **Ações no rodapé** — `Enter Confirmar` à esquerda, `Esc Cancelar` à direita; `──` preenchendo o espaço entre elas
6. **Medidor de força** — `Força: ████...░░ Label` com comprimento total de 10 blocos
7. **Linhas em branco** — padding vertical acima e abaixo de cada campo, e antes do medidor

Após verificar o `.golden.txt`, verificar o `.golden.json`:

1. **Borda** — tokens `border.focused` (`#7aa2f7` Tokyo Night) em todas as posições de borda
2. **Label ativo** — `accent.primary` + `bold` no label do campo focado
3. **Label inativo** — `text.secondary` (sem bold) no label sem foco
4. **Máscara** `••••••••` — `text.secondary`
5. **Ação bloqueada** — `text.disabled` na ação `Enter Confirmar` quando campo vazio
6. **Ação ativa** — `accent.primary` + `bold` na ação `Enter Confirmar` quando desbloqueada
7. **Medidor preenchido** — `semantic.success` ou `semantic.warning` conforme nível
8. **Medidor vazio** — `text.disabled`
9. **Fundo do campo** — `surface.input` (`#1e1f2e` Tokyo Night) em toda a linha digitável

---

## Checklist de Implementação

- [ ] `password_strength.go` — `EvaluateStrength` + `RenderStrengthMeter`
- [ ] `password_field.go` — `PasswordField` com buffer, HandleKey, Render, Wipe
- [ ] `modal.go` — adicionar `Cursor()` à interface `ModalView`
- [ ] `confirm_modal.go` — implementar `Cursor() nil`
- [ ] `help_modal.go` — implementar `Cursor() nil`
- [ ] `password_entry_modal.go` — `PasswordEntryModal` completo
- [ ] `password_create_modal.go` — `PasswordCreateModal` completo
- [ ] `root.go` — propagação inline de cursor em `View()` (via `lipgloss.Size` + `top.Cursor(topY, leftX)`)
- [ ] Testes golden: `password_field` (4 variantes)
- [ ] Testes golden: `password_entry` (2 variantes)
- [ ] Testes golden: `password_create` (4 variantes)
- [ ] Testes golden: `strength_meter` (3 variantes)
- [ ] Verificação visual dos `.golden.txt` — bordas, título, labels, ações, medidor
- [ ] Verificação de estilos dos `.golden.json` — cores e atributos de cada elemento
