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
//
// Flow:
//  1. Init: check for external file modifications (background I/O).
//  2. If external mod detected: show "Conflito de Modificação" decision dialog.
//  3. If user chooses "Sobrescrever e sair": proceed to save.
//  4. If user chooses "Voltar": end flow (no save, no exit).
//  5. On successful save: tea.Quit.
//  6. On save error: show error message, end flow (stay in app).
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
)

// extModDetectedMsg - emitted when an external file modification is detected during
// save-and-exit. Routed to activeFlow via modalResult.
type extModDetectedMsg struct{}

func (extModDetectedMsg) isModalResult() {}

// extModOverwriteMsg - emitted by the "Conflito de Modificação" decision dialog
// when the user chooses "Sobrescrever e sair".
type extModOverwriteMsg struct{}

func (extModOverwriteMsg) isModalResult() {}

// extModCancelMsg - emitted by the "Conflito de Modificação" decision dialog
// when the user chooses "Voltar" (cancel the save-and-exit).
type extModCancelMsg struct{}

func (extModCancelMsg) isModalResult() {}

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
	switch msg.(type) {
	case saveAndExitReadyMsg:
		// No external modification detected — proceed to save
		return f.doSave()

	case extModDetectedMsg:
		// External modification detected: ask user what to do
		return func() tea.Msg {
			return Decision(SeverityDestructive, "Conflito de Modificação",
				"O arquivo do cofre foi modificado externamente desde a última vez que foi carregado.\n\nSobrescrever descartará as alterações externas.",
				DecisionAction{Key: "Enter", Label: "Sobrescrever e sair", Default: true,
					Cmd: func() tea.Msg { return extModOverwriteMsg{} }},
				nil,
				DecisionAction{Key: "Esc", Label: "Voltar",
					Cmd: func() tea.Msg { return extModCancelMsg{} }})
		}

	case extModOverwriteMsg:
		// User confirmed overwrite: proceed to save
		return f.doSave()

	case extModCancelMsg:
		// User cancelled: end flow without saving or exiting
		return endFlow()

	case saveAndExitOKMsg:
		// Save succeeded: quit the application
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
