package tui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/vault"
)

// vaultSaver is the minimal interface saveAndExitFlow needs from vault.Manager.
// Using an interface (rather than *vault.Manager directly) enables test injection.
type vaultSaver interface {
	Salvar() error
}

// saveAndExitFlow implements the state machine for saving and exiting.
// Triggered by Ctrl+Q when m.mgr != nil && m.mgr.IsModified().
//
// States:
//
//	stateCheckExtMod -> stateSaveAndExit -> stateDoneExit
//	stateCheckExtMod -> stateSaveAsNew (user chose "N Salvar como novo")
//
// Flow:
//  1. Init: check for external file modifications (background I/O).
//  2. If external mod detected: show "Salvar cofre" decision dialog (S/N/Esc).
//  3. If user chooses "S Sobrescrever": proceed to save.
//  4. If user chooses "N Salvar como novo": open file picker for alternate path.
//  5. If user chooses "Esc Voltar": end flow (no save, no exit).
//  6. On successful save: tea.Quit.
//  7. On save error: show error message, end flow (stay in app).
type saveAndExitFlow struct {
	state    int
	mgr      vaultSaver
	path     string
	metadata storage.FileMetadata
	messages *MessageManager
	theme    *Theme
}

// State constants for saveAndExitFlow (offset 20 to avoid collision).
const (
	stateCheckExtMod = iota + 20
	stateSaveAndExit
	stateDoneExit
	stateSaveAsNew // user chose "N Salvar como novo" — waiting for new path from file picker
)

// extModDetectedMsg - emitted when an external file modification is detected during
// save-and-exit. Routed to activeFlow via modalResult.
type extModDetectedMsg struct{}

func (extModDetectedMsg) isModalResult() {}

// extModOverwriteMsg - emitted by the conflict decision dialog when the user
// chooses "S Sobrescrever".
type extModOverwriteMsg struct{}

func (extModOverwriteMsg) isModalResult() {}

// extModCancelMsg - emitted by the conflict dialog when the user chooses "Voltar"
// (cancel the save-and-exit).
type extModCancelMsg struct{}

func (extModCancelMsg) isModalResult() {}

// extModSaveAsNewMsg - emitted by the conflict dialog when the user chooses
// "N Salvar como novo". Opens a file picker to choose an alternate save path.
type extModSaveAsNewMsg struct{}

func (extModSaveAsNewMsg) isModalResult() {}

// saveAndExitReadyMsg - internal signal that no external mod was detected;
// the flow may proceed directly to save. Not a modalResult.
type saveAndExitReadyMsg struct{}

// saveAndExitOKMsg - internal signal that Salvar() succeeded. Not a modalResult.
// rootModel routes this to activeFlow's Update, which returns tea.Quit.
type saveAndExitOKMsg struct{}

// newSaveAndExitFlow creates a saveAndExitFlow using the provided manager,
// vault path, and file metadata snapshot for external change detection.
// mgr must be non-nil and satisfy vaultSaver (vault.Manager does).
func newSaveAndExitFlow(mgr *vault.Manager, path string, metadata storage.FileMetadata, messages *MessageManager, theme *Theme) *saveAndExitFlow {
	return &saveAndExitFlow{
		state:    stateCheckExtMod,
		mgr:      mgr,
		path:     path,
		metadata: metadata,
		messages: messages,
		theme:    theme,
	}
}

// Init starts the flow by launching a background check for external modifications.
func (f *saveAndExitFlow) Init() tea.Cmd {
	f.state = stateCheckExtMod
	path := f.path
	metadata := f.metadata
	return func() tea.Msg {
		changed, err := storage.DetectExternalChange(path, metadata)
		if err != nil || changed {
			return extModDetectedMsg{}
		}
		return saveAndExitReadyMsg{}
	}
}

// Update processes flow messages and transitions states.
func (f *saveAndExitFlow) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case saveAndExitReadyMsg:
		// No external modification detected — proceed to save.
		_ = msg
		return f.doSave()

	case extModDetectedMsg:
		// Desvio 9: External modification detected — show "Salvar cofre" dialog with
		// three actions: S Sobrescrever, N Salvar como novo, Esc Voltar.
		_ = msg
		return func() tea.Msg {
			return Decision(SeverityDestructive, "Salvar cofre",
				"Arquivo modificado externamente. Sobrescrever ou salvar como novo?",
				DecisionAction{Key: "S", Label: "Sobrescrever", Default: true,
					Cmd: func() tea.Msg { return extModOverwriteMsg{} }},
				[]DecisionAction{
					{Key: "N", Label: "Salvar como novo",
						Cmd: func() tea.Msg { return extModSaveAsNewMsg{} }},
				},
				DecisionAction{Key: "Esc", Label: "Voltar",
					Cmd: func() tea.Msg { return extModCancelMsg{} }})
		}

	case extModSaveAsNewMsg:
		// User chose "N Salvar como novo": open file picker to choose alternate save path.
		_ = msg
		f.state = stateSaveAsNew
		return func() tea.Msg {
			return pushModalMsg{modal: &filePickerModal{mode: FilePickerSave}}
		}

	case filePickerResult:
		// Only handled when in stateSaveAsNew (user chose alternate path after external conflict).
		if f.state != stateSaveAsNew {
			return nil
		}
		if msg.Cancelled {
			return endFlow()
		}
		// Update path to the new location chosen by the user, then proceed to save.
		f.path = msg.Path
		return f.doSave()

	case extModOverwriteMsg:
		// User confirmed overwrite: proceed to save.
		_ = msg
		return f.doSave()

	case extModCancelMsg:
		// User cancelled: end flow without saving or exiting.
		_ = msg
		return endFlow()

	case saveAndExitOKMsg:
		// Save succeeded: quit the application.
		_ = msg
		f.state = stateDoneExit
		return tea.Quit

	default:
		return nil
	}
}

// doSave calls m.mgr.Salvar() in a background Cmd.
// On success it emits saveAndExitOKMsg which Update converts to tea.Quit.
// On failure it shows an error and ends the flow so the user can stay in the app.
func (f *saveAndExitFlow) doSave() tea.Cmd {
	f.state = stateSaveAndExit
	mgr := f.mgr
	messages := f.messages
	return func() tea.Msg {
		if err := mgr.Salvar(); err != nil {
			messages.Show(MsgError, "Não foi possível salvar o cofre antes de sair.", 5, false)
			return endFlowMsg{}
		}
		return saveAndExitOKMsg{}
	}
}
