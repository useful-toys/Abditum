package modal

import (
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// ConfirmModal exibe um diálogo de confirmação com título, mensagem e ações de confirmar/cancelar.
// Implemente tui.ModalView e é criado via NewConfirmModal.
type ConfirmModal struct {
	// title é o cabeçalho exibido no topo do modal.
	title string
	// message é o texto descritivo exibido ao usuário.
	message string
	// options lista as ações disponíveis, ativadas por teclas específicas.
	options []ModalOption
}

// NewConfirmModal cria um ConfirmModal com as ações padrão de confirmar (Enter) e cancelar (Esc).
func NewConfirmModal(title, message string) *ConfirmModal {
	return &ConfirmModal{
		title:   title,
		message: message,
		options: []ModalOption{
			{
				Keys:   []design.Key{design.Keys.Enter},
				Label:  "Confirmar",
				Intent: IntentConfirm,
				Action: func() tea.Cmd {
					return func() tea.Msg { return tui.ModalReadyMsg{} }
				},
			},
			{
				Keys:   []design.Key{design.Keys.Esc},
				Label:  "Cancelar",
				Intent: IntentCancel,
				Action: func() tea.Cmd {
					return func() tea.Msg { return tui.CloseModalMsg{} }
				},
			},
		},
	}
}

// Render retorna o modal estilizado com borda, usando as cores do tema fornecido.
func (m *ConfirmModal) Render(maxHeight, maxWidth int, theme *design.Theme) string {
	content := m.title + "\n\n" + m.message
	style := lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder(), true).
		BorderForeground(lipgloss.Color(theme.Border.Default)).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Raised))
	return style.Render(content)
}

// HandleKey verifica se a tecla pressionada corresponde a alguma opção do modal e executa sua ação.
func (m *ConfirmModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	for _, opt := range m.options {
		for _, k := range opt.Keys {
			if k.Matches(msg) {
				return opt.Action()
			}
		}
	}
	return nil
}

// Update processa mensagens do Bubble Tea. ConfirmModal não reage a mensagens internas.
func (m *ConfirmModal) Update(msg tea.Msg) tea.Cmd {
	return nil
}
