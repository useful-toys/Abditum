package main

import (
	"fmt"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
)

// pocMessagesModel is the standalone PoC model demonstrating the full
// message bar + command bar system (D-01 through D-06).
// No vault.Manager, no rootModel — completely self-contained.
type pocMessagesModel struct {
	width    int
	height   int
	actions  *tui.ActionManager
	msgs     *tui.MessageManager
	showHelp bool // simple toggle for F1 inline help
}

// newPocModel creates a pocMessagesModel with all 15 actions from the D-06 table registered.
func newPocModel() *pocMessagesModel {
	m := &pocMessagesModel{
		actions: tui.NewActionManager(),
		msgs:    tui.NewMessageManager(),
	}

	// D-06 action table — Group 1 "Mensagens"
	// 8 visible in command bar (HideFromBar: false), 5 shift+Fx hidden, ctrl+q + f1
	m.actions.Register(m,
		tui.Action{Keys: []string{"f2"}, Label: "Dica uso", Description: "Mostrar MsgHint permanente",
			Group: 1, Scope: tui.ScopeLocal, Priority: 90, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.msgs.Show(tui.MsgHint, "Dica de uso permanente", 0, false); return nil }},
		tui.Action{Keys: []string{"f3"}, Label: "Dica campo", Description: "Mostrar MsgHint de campo",
			Group: 1, Scope: tui.ScopeLocal, Priority: 80, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.msgs.Show(tui.MsgHint, "Dica de campo permanente", 0, false); return nil }},
		tui.Action{Keys: []string{"f4"}, Label: "Info", Description: "Mostrar MsgInfo",
			Group: 1, Scope: tui.ScopeLocal, Priority: 70, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.msgs.Show(tui.MsgInfo, "Informação neutra", 0, false); return nil }},
		tui.Action{Keys: []string{"f5"}, Label: "Alerta", Description: "Mostrar MsgWarn",
			Group: 1, Scope: tui.ScopeLocal, Priority: 60, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.msgs.Show(tui.MsgWarn, "Alerta de atenção", 0, false); return nil }},
		tui.Action{Keys: []string{"f6"}, Label: "Erro", Description: "Mostrar MsgError",
			Group: 1, Scope: tui.ScopeLocal, Priority: 50, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.msgs.Show(tui.MsgError, "Erro de operação", 0, false); return nil }},
		tui.Action{Keys: []string{"f7"}, Label: "Ocupado", Description: "Mostrar MsgBusy com spinner",
			Group: 1, Scope: tui.ScopeLocal, Priority: 40, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.msgs.Show(tui.MsgBusy, "Processando...", 0, false); return nil }},
		tui.Action{Keys: []string{"f8"}, Label: "Sucesso", Description: "Mostrar MsgSuccess",
			Group: 1, Scope: tui.ScopeLocal, Priority: 30, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.msgs.Show(tui.MsgSuccess, "Operação concluída", 0, false); return nil }},
		tui.Action{Keys: []string{"f9"}, Label: "Limpar", Description: "Limpar barra de mensagens",
			Group: 1, Scope: tui.ScopeLocal, Priority: 20, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.msgs.Clear(); return nil }},

		// shift+Fx variants — HideFromBar: true (visible in Help modal only)
		tui.Action{Keys: []string{"shift+f2"}, Label: "Dica uso 5s", Description: "MsgHint com TTL de 5s",
			Group: 1, Scope: tui.ScopeLocal, Priority: 89, HideFromBar: true,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.msgs.Show(tui.MsgHint, "Dica de uso (5s)", 5, false); return nil }},
		tui.Action{Keys: []string{"shift+f3"}, Label: "Dica campo 5s", Description: "MsgHint campo com TTL de 5s",
			Group: 1, Scope: tui.ScopeLocal, Priority: 79, HideFromBar: true,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.msgs.Show(tui.MsgHint, "Dica de campo (5s)", 5, false); return nil }},
		tui.Action{Keys: []string{"shift+f4"}, Label: "Info 5s", Description: "MsgInfo com TTL de 5s",
			Group: 1, Scope: tui.ScopeLocal, Priority: 69, HideFromBar: true,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.msgs.Show(tui.MsgInfo, "Informação (5s)", 5, false); return nil }},
		tui.Action{Keys: []string{"shift+f5"}, Label: "Alerta 5s", Description: "MsgWarn com TTL de 5s",
			Group: 1, Scope: tui.ScopeLocal, Priority: 59, HideFromBar: true,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.msgs.Show(tui.MsgWarn, "Alerta (5s)", 5, false); return nil }},
		tui.Action{Keys: []string{"shift+f6"}, Label: "Erro 5s", Description: "MsgError com TTL de 5s",
			Group: 1, Scope: tui.ScopeLocal, Priority: 49, HideFromBar: true,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.msgs.Show(tui.MsgError, "Erro (5s)", 5, false); return nil }},

		// Navigation actions
		tui.Action{Keys: []string{"ctrl+q"}, Label: "Sair", Description: "Sair do PoC",
			Group: 1, Scope: tui.ScopeLocal, Priority: 10, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { return tea.Quit }},
		tui.Action{Keys: []string{"f1"}, Label: "Ajuda", Description: "Mostrar/ocultar ajuda",
			Group: 1, Scope: tui.ScopeGlobal, Priority: 0, HideFromBar: false,
			Enabled: func() bool { return true },
			Handler: func() tea.Cmd { m.showHelp = !m.showHelp; return nil }},
	)

	m.actions.RegisterGroupLabel(1, "Mensagens")
	return m
}

// Init starts the per-second tick for MessageManager (spinner animation + TTL countdown).
func (m *pocMessagesModel) Init() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg { return tui.TickMsg(t) })
}

// Update handles window resize, ticks, and keyboard input.
func (m *pocMessagesModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tui.TickMsg:
		m.msgs.Tick()
		return m, tea.Tick(time.Second, func(t time.Time) tea.Msg { return tui.TickMsg(t) })

	case tea.KeyPressMsg:
		m.msgs.HandleInput()
		// ESC closes help overlay
		if m.showHelp && msg.String() == "esc" {
			m.showHelp = false
			return m, nil
		}
		if cmd := m.actions.Dispatch(msg.String(), m.showHelp); cmd != nil {
			return m, cmd
		}
		return m, nil
	}
	return m, nil
}

// View composes the 4-zone layout: header + work area + message bar + command bar.
func (m *pocMessagesModel) View() tea.View {
	if m.width == 0 || m.height == 0 {
		v := tea.NewView("Initializing...")
		v.AltScreen = true
		return v
	}

	headerStyle   := lipgloss.NewStyle().Width(m.width).Bold(true).Foreground(lipgloss.Color("#a9b1d6"))
	workAreaStyle := lipgloss.NewStyle().Width(m.width)

	const headerH = 1
	const msgBarH = 1
	const cmdBarH = 1
	workH := m.height - headerH - msgBarH - cmdBarH
	if workH < 0 {
		workH = 0
	}

	// Header
	header := headerStyle.Render("  Abditum — PoC Mensagens")

	// Work area: help overlay or centered label
	var workContent string
	if m.showHelp {
		workContent = m.renderHelp()
	} else {
		workContent = lipgloss.Place(m.width, workH, lipgloss.Center, lipgloss.Center,
			lipgloss.NewStyle().Foreground(lipgloss.Color("#565f89")).Render(
				"Teste das mensagens\n\n"+
					"F2–F9: exibir mensagens  ·  shift+F2–F6: variantes com TTL\n"+
					"F1: ajuda  ·  ctrl+Q: sair"))
	}
	workArea := workAreaStyle.Height(workH).Render(workContent)

	// Message bar (from Plan 01)
	msgBar := tui.RenderMessageBar(m.msgs.Current(), m.width)

	// Command bar (from Plan 02)
	var cmdBarContent string
	if m.showHelp {
		cmdBarContent = lipgloss.NewStyle().Foreground(lipgloss.Color("#565f89")).Render("  esc  Fechar ajuda")
	} else {
		cmdBarContent = m.actions.RenderCommandBar(m.width)
	}

	content := strings.Join([]string{header, workArea, msgBar, cmdBarContent}, "\n")
	v := tea.NewView(content)
	v.AltScreen = true
	return v
}

// renderHelp renders an inline help view listing all registered actions.
// Uses a simple string builder — no need for helpModal machinery in a PoC.
func (m *pocMessagesModel) renderHelp() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7aa2f7"))
	keyStyle   := lipgloss.NewStyle().Foreground(lipgloss.Color("#7aa2f7")).Bold(true)
	sepStyle   := lipgloss.NewStyle().Foreground(lipgloss.Color("#565f89"))
	dimStyle   := lipgloss.NewStyle().Foreground(lipgloss.Color("#565f89"))

	var b strings.Builder
	b.WriteString("\n")
	b.WriteString("  " + titleStyle.Render("Atalhos — Mensagens") + "\n")
	b.WriteString("  " + sepStyle.Render(strings.Repeat("─", 46)) + "\n\n")

	for _, act := range m.actions.All() {
		if len(act.Keys) == 0 {
			continue
		}
		hideTag := ""
		if act.HideFromBar {
			hideTag = dimStyle.Render(" (oculto da barra)")
		}
		b.WriteString(fmt.Sprintf("  %-14s  %s%s\n",
			keyStyle.Render(act.Keys[0]),
			act.Description,
			hideTag,
		))
	}
	b.WriteString("\n  " + dimStyle.Render("F1 ou Esc para fechar") + "\n")
	return b.String()
}
