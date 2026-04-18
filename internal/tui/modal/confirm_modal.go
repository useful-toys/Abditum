package modal

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// ConfirmModal exibe um diálogo de confirmação com título, mensagem, severidade e ações.
// Implementa tui.ModalView. Criado via NewConfirmModal ou NewConfirmModalSeverity.
type ConfirmModal struct {
	severity design.Severity
	title    string
	message  string
	options  []ModalOption
	keys     KeyHandler // despacha teclas das opções; sem scroll
}

// NewConfirmModal cria um ConfirmModal de severidade Neutra com as opções fornecidas.
// opts define as ações disponíveis — o caller injeta os closures corretos.
// Convenção: 1ª opção é a ação principal (Enter); última é o cancelamento (Esc).
func NewConfirmModal(title, message string, opts []ModalOption) *ConfirmModal {
	return NewConfirmModalSeverity(design.SeverityNeutral, title, message, opts)
}

// NewConfirmModalSeverity cria um ConfirmModal com severidade visual explícita.
// Se há apenas 1 opção, ESC é adicionado automaticamente como alias para disparar a mesma ação.
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

// Render constrói um DialogFrame com cores e símbolo derivados da severidade,
// e passa o corpo (mensagem com padding) para o frame renderizar.
func (m *ConfirmModal) Render(maxHeight, maxWidth int, theme *design.Theme) string {
	// Corpo: padding vertical acima e abaixo da mensagem.
	padding := strings.Repeat("\n", design.DialogPaddingV)
	body := padding + m.message + padding

	frame := DialogFrame{
		Title:           m.title,
		TitleColor:      theme.Text.Primary,
		Symbol:          m.severity.Symbol(),
		SymbolColor:     m.severity.BorderColor(theme),
		BorderColor:     m.severity.BorderColor(theme),
		Options:         m.options,
		DefaultKeyColor: m.severity.DefaultKeyColor(theme),
		Scroll:          nil,
	}
	return frame.Render(body, maxWidth, theme)
}

// HandleKey delega para m.keys.Handle(msg).
func (m *ConfirmModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if cmd, handled := m.keys.Handle(msg); handled {
		return cmd
	}
	return nil
}

// Update processa mensagens Bubble Tea. Delega para HandleKey em tea.KeyMsg.
func (m *ConfirmModal) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		return m.HandleKey(key)
	}
	return nil
}

// Cursor retorna nil — ConfirmModal não tem campo de texto com cursor.
func (m *ConfirmModal) Cursor(_, _ int) *tea.Cursor {
	return nil
}
