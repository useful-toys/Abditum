package tui

import (
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/vault"
)

// createVaultFlow implements the state machine for creating a new vault.
// States: stateCheckDirty -> statePickFile -> stateCheckOverwrite -> statePwdCreate -> stateSaveNew -> done.
type createVaultFlow struct {
	state      int // current state in the state machine
	width      int
	height     int
	theme      *Theme
	mgr        *vault.Manager
	messages   *MessageManager
	actions    *ActionManager
	targetPath string // path selected by user for new vault
}

// State constants for createVaultFlow (reuses common states)
const (
	stateCheckOverwrite = iota + 10 // offset to avoid collision with openVaultFlow states
	statePwdCreate
	stateSaveNew
)

// newCreateVaultFlow creates and initializes a createVaultFlow.
func newCreateVaultFlow(mgr *vault.Manager, messages *MessageManager, actions *ActionManager, theme *Theme) *createVaultFlow {
	return &createVaultFlow{
		state:    stateCheckDirty,
		theme:    theme,
		mgr:      mgr,
		messages: messages,
		actions:  actions,
	}
}

// Init starts the flow. If there are unsaved changes, prompt the user.
// Otherwise, proceed directly to file selection.
func (f *createVaultFlow) Init() tea.Cmd {
	if f.mgr != nil && f.mgr.IsModified() {
		f.state = stateCheckDirty
		// Push confirmation dialog for unsaved changes
		return func() tea.Msg {
			return Acknowledge(SeverityNeutral, "Alterações não salvas",
				"Existem alterações não salvas. Criar um novo cofre descartará as mudanças.", nil)
		}
	}
	f.state = statePickFile
	// Push file picker modal (in save mode)
	return func() tea.Msg {
		return pushModalMsg{modal: &filePickerModal{mode: FilePickerFile}}
	}
}

// Update processes messages from modals and transitions states.
func (f *createVaultFlow) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case filePickerResult:
		if msg.Cancelled {
			return endFlow()
		}
		f.targetPath = msg.Path
		f.state = stateCheckOverwrite
		// Check if file exists
		if _, err := os.Stat(f.targetPath); err == nil {
			// File exists - show overwrite confirmation
			return func() tea.Msg {
				return Decision(SeverityDestructive, "Arquivo existe",
					"Um cofre já existe neste caminho. Deseja sobrescrever?",
					DecisionAction{Key: "Enter", Label: "Sobrescrever", Default: true},
					nil,
					DecisionAction{Key: "Esc", Label: "Voltar"})
			}
		}
		// File doesn't exist - proceed to password creation
		f.state = statePwdCreate
		return func() tea.Msg {
			return pushModalMsg{modal: &passwordCreateModal{}}
		}

	case pwdCreatedMsg:
		f.state = stateSaveNew
		password := msg.Password
		path := f.targetPath
		// Create and save the vault in a background command
		return func() tea.Msg {
			// Create a new empty vault
			newVault := vault.NovoCofre()
			newVault.InicializarConteudoPadrao()

			// Attempt to save the vault
			err := storage.SaveNew(path, newVault, password)
			crypto.Wipe(password)

			if err != nil {
				f.messages.Show(MsgError, "Não foi possível salvar o cofre.", 5, false)
				return endFlow()
			}
			// Success
			return vaultOpenedMsg{Path: path}
		}

	case flowCancelledMsg:
		return endFlow()

	default:
		return nil
	}
}

// View is not used in flow-driven interface; modals handle visualization.
func (f *createVaultFlow) View(width, height int) string {
	return ""
}
