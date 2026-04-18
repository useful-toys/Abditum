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
