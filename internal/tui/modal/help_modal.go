package modal

import (
	"sort"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// helpTitle é o título do modal de ajuda conforme tui-spec-dialog-help.md.
const helpTitle = "Ajuda — Atalhos e Ações"

// keyColumnWidth é a largura reservada para a coluna de teclas no corpo do HelpModal.
const keyColumnWidth = 14

// HelpModal exibe todas as actions registradas, agrupadas por ActionGroup.
// Suporta scroll quando o conteúdo excede o espaço disponível.
// Implementa tui.ModalView.
type HelpModal struct {
	actions []actions.Action
	groups  []actions.ActionGroup
	scroll  ScrollState // estado de scroll — começa em Offset=0
	keys    KeyHandler  // despacha scroll (↑↓PgUp/PgDn/Home/End) e Esc (fechar)
}

// NewHelpModal cria o HelpModal com as actions e grupos fornecidos.
// Scroll começa no topo (Offset = 0).
func NewHelpModal(acts []actions.Action, groups []actions.ActionGroup) *HelpModal {
	m := &HelpModal{
		actions: acts,
		groups:  groups,
	}
	closeOpts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Fechar",
			Intent: IntentCancel,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
	}
	m.keys = KeyHandler{
		Options: closeOpts,
		Scroll:  &m.scroll,
	}
	return m
}

// Render gera o corpo dinamicamente, fatia o viewport conforme scroll,
// e passa para DialogFrame.Render.
func (m *HelpModal) Render(maxHeight, maxWidth int, theme *design.Theme) string {
	allLines := m.buildAllLines(maxWidth, theme)

	// viewport = maxHeight - 2 (borda superior + borda de rodapé)
	viewport := maxHeight - 2
	if viewport < 1 {
		viewport = 1
	}

	// Atualizar estado de scroll.
	m.scroll.Total = len(allLines)
	m.scroll.Viewport = viewport

	// Fatiar linhas visíveis.
	start := m.scroll.Offset
	end := start + viewport
	if end > len(allLines) {
		end = len(allLines)
	}
	visibleLines := allLines[start:end]
	body := strings.Join(visibleLines, "\n")

	// Configurar o frame.
	closeOpts := []ModalOption{
		{
			Keys:   []design.Key{design.Keys.Esc},
			Label:  "Fechar",
			Intent: IntentCancel,
			Action: func() tea.Cmd { return tui.CloseModal() },
		},
	}
	var scrollPtr *ScrollState
	if m.scroll.Total > m.scroll.Viewport {
		scrollPtr = &m.scroll
	}

	frame := DialogFrame{
		Title:           helpTitle,
		TitleColor:      theme.Text.Primary,
		Symbol:          "",
		SymbolColor:     "",
		BorderColor:     theme.Border.Default,
		Options:         closeOpts,
		DefaultKeyColor: theme.Accent.Primary,
		Scroll:          scrollPtr,
	}
	return frame.Render(body, maxWidth, theme)
}

// buildAllLines gera todas as linhas de conteúdo do modal de ajuda.
// Grupos ordenados por ActionGroup.Order crescente; ações ordenadas por Action.Priority crescente.
// Linha em branco entre grupos (não antes do primeiro, não após o último).
func (m *HelpModal) buildAllLines(maxWidth int, theme *design.Theme) []string {
	// Ordenar grupos por Order (estável para empates).
	sortedGroups := make([]actions.ActionGroup, len(m.groups))
	copy(sortedGroups, m.groups)
	sort.SliceStable(sortedGroups, func(i, j int) bool {
		return sortedGroups[i].Order < sortedGroups[j].Order
	})

	// Mapear GroupID → actions, ordenadas por Priority.
	groupActions := make(map[string][]actions.Action)
	for _, a := range m.actions {
		groupActions[a.GroupID] = append(groupActions[a.GroupID], a)
	}
	for id := range groupActions {
		sort.SliceStable(groupActions[id], func(i, j int) bool {
			return groupActions[id][i].Priority < groupActions[id][j].Priority
		})
	}

	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Secondary)).
		Bold(true)
	keyStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Accent.Primary))
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Primary))

	var allLines []string
	for i, grp := range sortedGroups {
		if i > 0 {
			allLines = append(allLines, "") // linha em branco entre grupos
		}
		// Cabeçalho do grupo
		allLines = append(allLines, headerStyle.Render(grp.Label))

		// Ações do grupo
		for _, act := range groupActions[grp.ID] {
			keyLabel := ""
			if len(act.Keys) > 0 {
				keyLabel = act.Keys[0].Label
			}
			// Coluna de tecla: largura fixa de keyColumnWidth, pad com espaços
			keyRendered := keyStyle.Render(keyLabel)
			keyVisualWidth := lipgloss.Width(keyRendered)
			pad := keyColumnWidth - keyVisualWidth
			if pad < 1 {
				pad = 1
			}
			line := keyRendered + strings.Repeat(" ", pad) + descStyle.Render(act.Description)
			allLines = append(allLines, line)
		}
	}
	return allLines
}

// HandleKey delega para m.keys.Handle(msg).
func (m *HelpModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if cmd, handled := m.keys.Handle(msg); handled {
		return cmd
	}
	return nil
}

// Update processa mensagens Bubble Tea. Delega para HandleKey em tea.KeyMsg.
func (m *HelpModal) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		return m.HandleKey(key)
	}
	return nil
}
