package operation

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/vault"
)

// guardSaveMsg dispara a ação de salvar no guard.
type guardSaveMsg struct {
	forced bool // true = forçar sobrescrita de modificação externa
}

// guardSaveResultMsg carrega o resultado da tentativa de salvar no guard.
type guardSaveResultMsg struct {
	err    error
	forced bool
}

// guardDiscardMsg sinaliza que o usuário optou por descartar alterações.
type guardDiscardMsg struct{}

// guardCancelMsg sinaliza que o usuário optou por voltar (cancelar).
type guardCancelMsg struct{}

// guardCofreAlterado verifica se há um cofre com alterações não salvas antes
// de prosseguir para outra operação. É um helper interno do pacote operation/,
// compartilhado pelos Fluxos 1, 2 e 5.
//
// Se não há cofre ou o cofre está inalterado, chama onProceder() diretamente.
// Se há alterações, abre modal de decisão (Salvar / Descartar / Voltar).
// Se salvar falha com ErrModifiedExternally, abre modal de conflito.
type guardCofreAlterado struct {
	saver      vaultSaver
	notifier   tui.MessageController
	onProceder func() tea.Cmd
	onAbortado func() tea.Cmd
}

// novoGuardCofreAlterado cria um guardCofreAlterado.
// saver pode ser nil quando nenhum cofre está carregado.
func novoGuardCofreAlterado(
	notifier tui.MessageController,
	saver vaultSaver,
	onProceder func() tea.Cmd,
	onAbortado func() tea.Cmd,
) *guardCofreAlterado {
	return &guardCofreAlterado{
		saver:      saver,
		notifier:   notifier,
		onProceder: onProceder,
		onAbortado: onAbortado,
	}
}

// Init inicia o guard. Se não há cofre ou não há alterações, dispara onProceder
// imediatamente. Caso contrário, abre o modal de decisão.
func (g *guardCofreAlterado) Init() tea.Cmd {
	if g.saver == nil || !g.saver.IsModified() {
		return g.onProceder()
	}
	return tui.OpenModal(g.buildModifiedModal())
}

// Update trata as mensagens internas do guard.
func (g *guardCofreAlterado) Update(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case guardSaveMsg:
		g.notifier.SetBusy("Salvando...")
		forced := m.forced
		return func() tea.Msg {
			return guardSaveResultMsg{err: g.saver.Salvar(forced), forced: forced}
		}

	case guardSaveResultMsg:
		if m.err == nil {
			g.notifier.Clear()
			return g.onProceder()
		}
		if !m.forced && errors.Is(m.err, vault.ErrModifiedExternally) {
			g.notifier.Clear()
			return tui.OpenModal(g.buildConflictModal())
		}
		g.notifier.SetError(m.err.Error())
		return g.onAbortado()

	case guardDiscardMsg:
		return g.onProceder()

	case guardCancelMsg:
		return g.onAbortado()
	}
	return nil
}

// buildModifiedModal cria o modal de decisão quando há alterações não salvas.
func (g *guardCofreAlterado) buildModifiedModal() *modal.ConfirmModal {
	return modal.NewConfirmModal(
		"Alterações não salvas",
		"Há alterações não salvas. O que deseja fazer?",
		[]modal.ModalOption{
			{
				Keys:  []design.Key{design.Keys.Enter},
				Label: "Salvar e prosseguir",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return guardSaveMsg{forced: false}
					})
				},
			},
			{
				Keys:  []design.Key{design.Letter('d')},
				Label: "Descartar e prosseguir",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return guardDiscardMsg{}
					})
				},
			},
			{
				Keys:  []design.Key{design.Keys.Esc},
				Label: "Voltar",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return guardCancelMsg{}
					})
				},
			},
		},
	)
}

// buildConflictModal cria o modal de resolução de conflito quando o arquivo foi
// modificado externamente enquanto o cofre estava aberto.
func (g *guardCofreAlterado) buildConflictModal() *modal.ConfirmModal {
	return modal.NewConfirmModal(
		"Conflito",
		"O arquivo foi modificado externamente. Deseja sobrescrever?",
		[]modal.ModalOption{
			{
				Keys:  []design.Key{design.Keys.Enter},
				Label: "Sobrescrever e prosseguir",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return guardSaveMsg{forced: true}
					})
				},
			},
			{
				Keys:  []design.Key{design.Keys.Esc},
				Label: "Voltar",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return guardCancelMsg{}
					})
				},
			},
		},
	)
}
