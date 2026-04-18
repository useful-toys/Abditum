# Teclas Implícitas Enter/Esc em ModalOption — Plano de Implementação

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Tornar `ModalOption.Keys` opcional, fazendo `KeyHandler` e `DialogFrame` injetarem Enter/Esc automaticamente como teclas implícitas da primeira e última option, eliminando a repetição de declaração em callers.

**Architecture:** A lógica de teclas implícitas é centralizada no `KeyHandler.Handle()` (despacho) e no `DialogFrame` (footer/largura). A lógica manual e duplicada em `confirm_modal.go` é removida. `ModalOption` não muda — apenas seu uso.

**Tech Stack:** Go, bubbletea v2, lipgloss v2

---

## Mapa de arquivos

| Arquivo | Mudança |
|---|---|
| `internal/tui/modal/modal_base.go` | Atualizar comentário de `Keys` |
| `internal/tui/modal/key_handler.go` | Adicionar despacho implícito de Enter/Esc |
| `internal/tui/modal/key_handler_test.go` | Adicionar testes para o comportamento implícito |
| `internal/tui/modal/frame.go` | Usar key implícita quando `Keys` está vazio (footer e largura) |
| `internal/tui/modal/frame_test.go` | Adicionar teste de footer com `Keys` vazio |
| `internal/tui/modal/confirm_modal.go` | Remover lógica manual de injeção de Esc |
| `internal/tui/modal/confirm_modal_test.go` | Atualizar/adicionar testes para o novo comportamento |

---

## Tarefa 1: Atualizar comentário de `ModalOption.Keys`

**Arquivos:**
- Modificar: `internal/tui/modal/modal_base.go`

- [ ] **Passo 1: Atualizar o comentário do campo `Keys`**

Substituir o comentário atual:

```go
// ModalOption representa uma ação disponível ao usuário dentro de um modal.
type ModalOption struct {
	// Keys lista as teclas que ativam esta opção.
	// Keys[0].Label é exibido no rodapé do diálogo.
	// Outras Keys são aliases funcionais (ex: Enter como alias para "S Sobrescrever").
	Keys []design.Key
	// Label é o texto exibido ao usuário descrevendo a ação.
	Label string
	// Action é a função executada quando a opção é escolhida.
	Action func() tea.Cmd
}
```

Pelo novo:

```go
// ModalOption representa uma ação disponível ao usuário dentro de um modal.
type ModalOption struct {
	// Keys lista as teclas que ativam esta opção.
	// Keys[0].Label é exibido no rodapé do diálogo.
	// Outras Keys são aliases funcionais (ex: Enter como alias para "S Sobrescrever").
	//
	// Keys é opcional (nil ou vazio). Quando omitido, teclas implícitas são aplicadas
	// pelo KeyHandler e pelo DialogFrame:
	//   - Primeira option → Enter
	//   - Última option   → Esc
	//   - Option única    → Enter e Esc
	// Quando Keys está preenchido, as teclas implícitas são adicionadas como aliases.
	Keys []design.Key
	// Label é o texto exibido ao usuário descrevendo a ação.
	Label string
	// Action é a função executada quando a opção é escolhida.
	Action func() tea.Cmd
}
```

- [ ] **Passo 2: Executar os testes para garantir que nada quebrou**

```
go test ./internal/tui/modal/...
```

Esperado: PASS (só foi mudado um comentário).

- [ ] **Passo 3: Commit**

```
git add internal/tui/modal/modal_base.go
git commit -m "docs(modal): documenta que ModalOption.Keys é opcional com teclas implícitas"
```

---

## Tarefa 2: Adicionar despacho implícito de Enter/Esc no `KeyHandler`

**Arquivos:**
- Modificar: `internal/tui/modal/key_handler.go`
- Modificar: `internal/tui/modal/key_handler_test.go`

### Passo a passo

- [ ] **Passo 1: Escrever os testes que falham primeiro**

Adicionar ao final de `key_handler_test.go`:

```go
func TestKeyHandler_ImplicitEnter_FirstOption_NoKeys(t *testing.T) {
	// Primeira option sem Keys: Enter deve disparar sua ação.
	called := false
	opts := []ModalOption{
		{
			Label:  "Confirmar",
			Action: func() tea.Cmd { called = true; return nil },
		},
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Cancelar",
			Action: func() tea.Cmd { return nil },
		},
	}
	h := KeyHandler{Options: opts}
	_, handled := h.Handle(makeSpecialKeyMsg(tea.KeyEnter))
	if !handled {
		t.Error("Handle(Enter): handled = false, want true")
	}
	if !called {
		t.Error("Handle(Enter): action da primeira option não foi chamada")
	}
}

func TestKeyHandler_ImplicitEsc_LastOption_NoKeys(t *testing.T) {
	// Última option sem Keys: Esc deve disparar sua ação.
	called := false
	opts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Confirmar",
			Action: func() tea.Cmd { return nil },
		},
		{
			Label:  "Cancelar",
			Action: func() tea.Cmd { called = true; return nil },
		},
	}
	h := KeyHandler{Options: opts}
	_, handled := h.Handle(makeSpecialKeyMsg(tea.KeyEscape))
	if !handled {
		t.Error("Handle(Esc): handled = false, want true")
	}
	if !called {
		t.Error("Handle(Esc): action da última option não foi chamada")
	}
}

func TestKeyHandler_ImplicitBoth_SingleOption_NoKeys(t *testing.T) {
	// Option única sem Keys: Enter e Esc devem disparar a mesma ação.
	callCount := 0
	opts := []ModalOption{
		{
			Label:  "OK",
			Action: func() tea.Cmd { callCount++; return nil },
		},
	}
	h := KeyHandler{Options: opts}

	h.Handle(makeSpecialKeyMsg(tea.KeyEnter))
	h.Handle(makeSpecialKeyMsg(tea.KeyEscape))
	if callCount != 2 {
		t.Errorf("Single option sem Keys: callCount = %d, want 2 (Enter e Esc ambos devem disparar)", callCount)
	}
}

func TestKeyHandler_ImplicitEnter_AddsToExplicitKeys(t *testing.T) {
	// Primeira option com Keys: [letter('s')].
	// Enter deve ser adicionado como alias — action chamada por 's' e por Enter.
	callCount := 0
	opts := []ModalOption{
		{
			Keys:   []design.Key{design.Letter('s')},
			Label:  "Sim",
			Action: func() tea.Cmd { callCount++; return nil },
		},
		{
			Keys:   []design.Key{design.Letter('n')},
			Label:  "Não",
			Action: func() tea.Cmd { return nil },
		},
	}
	h := KeyHandler{Options: opts}

	h.Handle(tea.KeyPressMsg{Code: 's'})   // tecla explícita
	h.Handle(makeSpecialKeyMsg(tea.KeyEnter)) // alias implícito
	if callCount != 2 {
		t.Errorf("First option with explicit key + implicit Enter: callCount = %d, want 2", callCount)
	}
}

func TestKeyHandler_ImplicitKeys_DoNotOverrideExplicit(t *testing.T) {
	// Se Enter já está declarado explicitamente, não deve ser disparado duas vezes.
	callCount := 0
	opts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "OK",
			Action: func() tea.Cmd { callCount++; return nil },
		},
	}
	h := KeyHandler{Options: opts}

	h.Handle(makeSpecialKeyMsg(tea.KeyEnter))
	if callCount != 1 {
		t.Errorf("Enter explícito + implícito: callCount = %d, want 1 (não deve duplicar)", callCount)
	}
}
```

- [ ] **Passo 2: Executar os testes para confirmar que falham**

```
go test ./internal/tui/modal/... -run TestKeyHandler_Implicit
```

Esperado: FAIL — os casos implícitos ainda não estão implementados.

- [ ] **Passo 3: Implementar o despacho implícito em `key_handler.go`**

Substituir o método `Handle` completo:

```go
// Handle processa a tecla fornecida.
//
// Retorna (cmd, true) se a tecla foi consumida — execução de ação ou movimento de scroll.
// Retorna (nil, false) se a tecla não foi reconhecida.
//
// Ordem de despacho:
//  1. Opções: itera Options, compara com cada Key em opt.Keys usando key.Matches(msg).
//     No primeiro match, executa opt.Action() e retorna (cmd, true).
//  2. Teclas implícitas: Enter → primeira option; Esc → última option.
//     Aplicadas apenas se não consumidas no passo 1.
//  3. Scroll (apenas se Scroll != nil):
//     ↑ → Scroll.Up(), ↓ → Scroll.Down()
//     PgUp → Scroll.PageUp(), PgDn → Scroll.PageDown()
//     Home → Scroll.Home(), End → Scroll.End()
//     Após atualizar o estado, retorna (nil, true).
func (h *KeyHandler) Handle(msg tea.KeyMsg) (tea.Cmd, bool) {
	// 1. Despachar ações registradas explicitamente.
	for _, opt := range h.Options {
		for _, k := range opt.Keys {
			if k.Matches(msg) {
				return opt.Action(), true
			}
		}
	}

	// 2. Teclas implícitas: Enter → primeira option; Esc → última option.
	if len(h.Options) > 0 {
		if design.Keys.Enter.Matches(msg) {
			return h.Options[0].Action(), true
		}
		if design.Keys.Esc.Matches(msg) {
			return h.Options[len(h.Options)-1].Action(), true
		}
	}

	// 3. Navegar scroll (se configurado).
	if h.Scroll == nil {
		return nil, false
	}
	switch {
	case design.Keys.Up.Matches(msg):
		h.Scroll.Up()
		return nil, true
	case design.Keys.Down.Matches(msg):
		h.Scroll.Down()
		return nil, true
	case design.Keys.PgUp.Matches(msg):
		h.Scroll.PageUp()
		return nil, true
	case design.Keys.PgDn.Matches(msg):
		h.Scroll.PageDown()
		return nil, true
	case design.Keys.Home.Matches(msg):
		h.Scroll.Home()
		return nil, true
	case design.Keys.End.Matches(msg):
		h.Scroll.End()
		return nil, true
	}
	return nil, false
}
```

- [ ] **Passo 4: Executar todos os testes do pacote modal**

```
go test ./internal/tui/modal/...
```

Esperado: PASS. Verificar especialmente:
- `TestKeyHandler_UnrecognizedKey_ReturnsNotHandled` → **este teste vai mudar de comportamento**: antes, Esc com apenas Enter registrado retornava `false`; agora, Esc implícito dispara a única option. Esse teste precisará ser atualizado (ver Passo 5).

- [ ] **Passo 5: Atualizar `TestKeyHandler_UnrecognizedKey_ReturnsNotHandled`**

O teste atual verifica que Esc não é tratado quando apenas Enter está registrado. Com as teclas implícitas, Esc agora dispara a única option. O teste precisa refletir o novo comportamento:

```go
func TestKeyHandler_UnrecognizedKey_ReturnsNotHandled(t *testing.T) {
	// Tecla completamente fora do mapeamento (nem Enter/Esc nem explícitas).
	h := KeyHandler{Options: []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Tab},
			Label:  "OK",
			Action: func() tea.Cmd { return nil },
		},
	}}
	// 'a' não é Enter, Esc nem Tab — não deve ser tratada.
	_, handled := h.Handle(tea.KeyPressMsg{Code: 'a'})
	if handled {
		t.Error("Handle('a' quando apenas Tab registrado): handled = true, want false")
	}
}
```

- [ ] **Passo 6: Executar todos os testes novamente**

```
go test ./internal/tui/modal/...
```

Esperado: PASS em todos.

- [ ] **Passo 7: Commit**

```
git add internal/tui/modal/key_handler.go internal/tui/modal/key_handler_test.go
git commit -m "feat(modal): KeyHandler despacha Enter/Esc implícitos para primeira e última option"
```

---

## Tarefa 3: Atualizar `DialogFrame` para usar key implícita quando `Keys` está vazio

**Arquivos:**
- Modificar: `internal/tui/modal/frame.go`
- Modificar: `internal/tui/modal/frame_test.go`

### Passo a passo

- [ ] **Passo 1: Escrever o teste que falha primeiro**

Adicionar ao final de `frame_test.go`:

```go
func TestDialogFrame_ImplicitKeys_NoKeysInOptions(t *testing.T) {
	// Options sem Keys declarados: footer deve exibir Enter/Esc automaticamente.
	theme := design.TokyoNight
	opts := []modal.ModalOption{
		{Label: "Confirmar", Action: func() tea.Cmd { return nil }},
		{Label: "Cancelar", Action: func() tea.Cmd { return nil }},
	}
	body := "Mensagem de confirmação"
	testdata.TestRenderManaged(t, "frame", "implicit_keys", []string{"60x8"},
		renderFrame("Diálogo", "", "", theme.Border.Focused, theme.Accent.Primary, opts, nil, body))
}
```

- [ ] **Passo 2: Executar o teste para confirmar que falha**

```
go test ./internal/tui/modal/... -run TestDialogFrame_ImplicitKeys
```

Esperado: FAIL — o footer mostra linha vazia ou omite ações.

- [ ] **Passo 3: Adicionar função auxiliar `implicitKey` em `frame.go`**

Adicionar **antes** de `renderBottomBorder`, após a função `calculateBodyWidth`:

```go
// implicitKey retorna a tecla implícita de uma option dado seu índice e o total de options.
//   - Índice 0 (primeira): Enter
//   - Índice total-1 (última): Esc
//   - Outros índices: zero value (ok=false)
//
// Mesmas regras do KeyHandler: a option única é tanto primeira quanto última — Enter é retornado
// (o label do footer exibe Enter para a option única, assim como no rodapé convencional).
func implicitKey(index, total int) (design.Key, bool) {
	isFirst := index == 0
	isLast := index == total-1
	switch {
	case isFirst:
		return design.Keys.Enter, true
	case isLast:
		return design.Keys.Esc, true
	}
	return design.Key{}, false
}
```

- [ ] **Passo 4: Atualizar `calculateBodyWidth` para usar a key implícita**

Substituir o bloco de cálculo de `actionWidth` dentro de `calculateBodyWidth`:

Antes:
```go
	// Largura das ações
	actionWidth := 3
	for _, opt := range f.Options {
		if len(opt.Keys) == 0 {
			continue
		}
		_, keyWidth := design.RenderDialogAction(opt.Keys[0].Label, opt.Label, f.BorderColor, theme)
		actionWidth += keyWidth + 4 + 3
	}
```

Depois:
```go
	// Largura das ações
	actionWidth := 3
	for i, opt := range f.Options {
		keyLabel := ""
		if len(opt.Keys) > 0 {
			keyLabel = opt.Keys[0].Label
		} else if k, ok := implicitKey(i, len(f.Options)); ok {
			keyLabel = k.Label
		} else {
			continue
		}
		_, keyWidth := design.RenderDialogAction(keyLabel, opt.Label, f.BorderColor, theme)
		actionWidth += keyWidth + 4 + 3
	}
```

- [ ] **Passo 5: Atualizar `renderBottomBorder` para usar a key implícita**

Substituir o bloco de construção de `rendered` dentro de `renderBottomBorder`:

Antes:
```go
	var rendered []renderedOpt
	for i, opt := range f.Options {
		if len(opt.Keys) == 0 {
			continue
		}
		keyColor := f.BorderColor
		if i == 0 {
			keyColor = f.DefaultKeyColor
		}
		text, w := design.RenderDialogAction(opt.Keys[0].Label, opt.Label, keyColor, theme)
		rendered = append(rendered, renderedOpt{text: text, width: w})
	}
```

Depois:
```go
	var rendered []renderedOpt
	for i, opt := range f.Options {
		keyLabel := ""
		if len(opt.Keys) > 0 {
			keyLabel = opt.Keys[0].Label
		} else if k, ok := implicitKey(i, len(f.Options)); ok {
			keyLabel = k.Label
		} else {
			continue
		}
		keyColor := f.BorderColor
		if i == 0 {
			keyColor = f.DefaultKeyColor
		}
		text, w := design.RenderDialogAction(keyLabel, opt.Label, keyColor, theme)
		rendered = append(rendered, renderedOpt{text: text, width: w})
	}
```

- [ ] **Passo 6: Executar todos os testes**

```
go test ./internal/tui/modal/...
```

Esperado: PASS. O novo teste de golden file irá gerar o arquivo de referência na primeira execução.

- [ ] **Passo 7: Inspecionar o golden file gerado**

```
go test ./internal/tui/modal/... -run TestDialogFrame_ImplicitKeys -v
```

Verificar visualmente que o footer exibe `Enter Confirmar` à esquerda e `Esc Cancelar` à direita.

- [ ] **Passo 8: Commit**

```
git add internal/tui/modal/frame.go internal/tui/modal/frame_test.go
git add internal/tui/modal/testdata/
git commit -m "feat(modal): DialogFrame exibe Enter/Esc implícitos no rodapé quando Keys está vazio"
```

---

## Tarefa 4: Remover lógica manual de injeção de Esc de `confirm_modal.go`

**Arquivos:**
- Modificar: `internal/tui/modal/confirm_modal.go`
- Modificar: `internal/tui/modal/confirm_modal_test.go`

A lógica nas linhas 35-78 de `confirm_modal.go` injeta Esc diretamente em `Keys`. Com o `KeyHandler` tratando isso centralizadamente, essa lógica é redundante e deve ser removida.

### Passo a passo

- [ ] **Passo 1: Verificar quais testes de `confirm_modal_test.go` dependem da lógica manual**

Os seguintes testes testam o comportamento de Esc automático — continuarão passando pois o `KeyHandler` agora cobre isso:
- `TestConfirmModal_SingleOption_EscDispatchesAction`
- `TestConfirmModal_SingleOption_EnterStillWorks`
- `TestConfirmModal_SingleOption_EscAlreadyRegistered`
- `TestConfirmModal_SingleOption_WithCustomKeys`
- `TestConfirmModal_DifferentSeverities_AllWorkCorrectly`

O teste `TestConfirmModal_MultipleOptions_EscNotAutoAdded` **mudará de comportamento**: antes, Esc não era auto-adicionado com múltiplas options sem Esc explícita; agora, Esc implícito dispara a última option. Esse teste deve ser atualizado.

- [ ] **Passo 2: Atualizar `TestConfirmModal_MultipleOptions_EscNotAutoAdded`**

O teste atual verifica que Esc **não** funciona com múltiplas options sem Esc explícita. Com o novo design, Esc **sempre** dispara a última option (mesmo que essa option use Tab como tecla principal). Atualizar para refletir isso:

```go
func TestConfirmModal_MultipleOptions_EscDispatachesLastOption(t *testing.T) {
	// Com múltiplas opções, Esc dispara a última option (implicitamente),
	// mesmo que ela não declare Esc como tecla principal.
	outraCalled := false
	opts := []modal.ModalOption{
		{
			Keys:   []design.Key{design.Keys.Enter},
			Label:  "Confirmar",
			Action: func() tea.Cmd { return nil },
		},
		{
			Keys:   []design.Key{design.Keys.Tab},
			Label:  "Outra",
			Action: func() tea.Cmd { outraCalled = true; return nil },
		},
	}
	m := modal.NewConfirmModal("Título", "Mensagem", opts)

	// Esc dispara a última option (Tab Outra).
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEscape})
	if !outraCalled {
		t.Error("HandleKey(Esc) com múltiplas options: deve disparar a última option")
	}
}
```

- [ ] **Passo 3: Executar os testes antes de remover a lógica manual**

```
go test ./internal/tui/modal/... -run TestConfirmModal
```

Esperado: PASS — todos os testes devem passar com o teste atualizado.

- [ ] **Passo 4: Simplificar `NewConfirmModalSeverity` removendo a lógica manual**

Substituir a função `NewConfirmModalSeverity` completa:

Antes (linhas 29-82):
```go
func NewConfirmModalSeverity(severity design.Severity, title, message string, opts []ModalOption) *ConfirmModal {
	// Fazer uma cópia para não modificar o slice original do caller
	optsCopy := make([]ModalOption, len(opts))
	copy(optsCopy, opts)
	
	// Aplicar teclas implícitas quando Keys estiver vazio ou nil
	for i := range optsCopy {
		if optsCopy[i].Keys == nil || len(optsCopy[i].Keys) == 0 {
			// Determinar se é primeira, última ou ambas
			isFirst := i == 0
			isLast := i == len(optsCopy)-1
			
			switch {
			case isFirst && isLast:
				// Única opção: Enter e Esc (ambas ativam a mesma ação)
				optsCopy[i].Keys = []design.Key{design.Keys.Enter, design.Keys.Esc}
			case isFirst:
				// Primeira opção: Enter
				optsCopy[i].Keys = []design.Key{design.Keys.Enter}
			case isLast:
				// Última opção: Esc
				optsCopy[i].Keys = []design.Key{design.Keys.Esc}
			}
		}
	}
	
	m := &ConfirmModal{
		severity: severity,
		title:    title,
		message:  message,
		options:  optsCopy,
	}
	
	// Quando há apenas 1 ação, adiciona ESC como alias para disparar a mesma ação.
	// (Mantemos essa lógica existente para compatibilidade, embora agora seja redundante
	// para o caso de uma única opção com Keys vazio, pois já definimos Keys = [Enter, Esc] acima)
	if len(optsCopy) == 1 && optsCopy[0].Keys != nil {
		// Verifica se ESC já não está na lista de teclas
		hasEsc := false
		for _, k := range optsCopy[0].Keys {
			if k.Code == design.Keys.Esc.Code && k.Mod == design.Keys.Esc.Mod {
				hasEsc = true
				break
			}
		}
		// Se ESC ainda não está registrada, adiciona como alias
		if !hasEsc {
			optsCopy[0].Keys = append(optsCopy[0].Keys, design.Keys.Esc)
		}
	}
	
	m.keys = KeyHandler{Options: optsCopy}
	return m
}
```

Depois (versão simplificada):
```go
// NewConfirmModalSeverity cria um ConfirmModal com severidade visual explícita.
// As teclas Enter (primeira option) e Esc (última option) são gerenciadas automaticamente
// pelo KeyHandler — não é necessário declará-las em ModalOption.Keys.
func NewConfirmModalSeverity(severity design.Severity, title, message string, opts []ModalOption) *ConfirmModal {
	return &ConfirmModal{
		severity: severity,
		title:    title,
		message:  message,
		options:  opts,
		keys:     KeyHandler{Options: opts},
	}
}
```

Remover também o import de `"github.com/useful-toys/abditum/internal/tui/design"` se não for mais usado (verificar se `design.SeverityNeutral` e outros ainda são usados em `Render` — sim, são, então o import permanece).

- [ ] **Passo 5: Executar todos os testes do pacote modal**

```
go test ./internal/tui/modal/...
```

Esperado: PASS em todos.

- [ ] **Passo 6: Executar todos os testes do módulo**

```
go test ./...
```

Esperado: PASS. Verificar se algum caller externo ao pacote `modal` dependia do comportamento antigo.

- [ ] **Passo 7: Commit**

```
git add internal/tui/modal/confirm_modal.go internal/tui/modal/confirm_modal_test.go
git commit -m "refactor(modal): remove injeção manual de Esc em NewConfirmModalSeverity"
```

---

## Tarefa 5: Verificação final e testes de golden file

**Arquivos:**
- Nenhum arquivo novo — apenas execução e inspeção.

- [ ] **Passo 1: Executar suite completa**

```
go test ./...
```

Esperado: PASS em todos.

- [ ] **Passo 2: Verificar golden files dos confirm_modal existentes**

Se algum golden file de `confirm_modal` mudou (por causa da remoção da cópia de opts), regenerar:

```
go test ./internal/tui/modal/... -update
```

Inspecionar visualmente se o output é o mesmo de antes (a remoção da cópia não deve mudar o visual).

- [ ] **Passo 3: Commit final (se houver golden files atualizados)**

```
git add internal/tui/modal/testdata/
git commit -m "test(modal): atualiza golden files após remoção de cópia desnecessária em NewConfirmModalSeverity"
```
