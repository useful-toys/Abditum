package tui

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/vault"
)

// openVaultFlow implements the state machine for opening an existing vault.
// States: stateCheckDirty -> statePickFile -> statePwdEntry -> statePreload -> done.
type openVaultFlow struct {
	state           int // current state in the state machine
	width           int
	height          int
	theme           *Theme
	mgr             *vault.Manager
	messages        *MessageManager
	actions         *ActionManager
	pickedPath      string // path selected by user
	vaultMetadata   storage.FileMetadata
	passwordAttempt int
}

// State constants for openVaultFlow
const (
	stateCheckDirty = iota
	statePickFile
	statePwdEntry
	statePreload
	stateDone
)

// newOpenVaultFlow creates and initializes an openVaultFlow.
func newOpenVaultFlow(mgr *vault.Manager, messages *MessageManager, actions *ActionManager, theme *Theme) *openVaultFlow {
	return &openVaultFlow{
		state:    stateCheckDirty,
		theme:    theme,
		mgr:      mgr,
		messages: messages,
		actions:  actions,
	}
}

// Init starts the flow. If there are unsaved changes, prompt the user.
// Otherwise, proceed directly to file selection.
func (f *openVaultFlow) Init() tea.Cmd {
	if f.mgr != nil && f.mgr.IsModified() {
		f.state = stateCheckDirty
		// Push confirmation dialog for unsaved changes
		return func() tea.Msg {
			return Acknowledge(SeverityNeutral, "Alterações não salvas",
				"Existem alterações não salvas. Abrir outro cofre descartará as mudanças.", nil)
		}
	}
	f.state = statePickFile
	// Push file picker modal
	return func() tea.Msg {
		return pushModalMsg{modal: &filePickerModal{}}
	}
}

// Update processes messages from modals and transitions states.
func (f *openVaultFlow) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case filePickerResult:
		if msg.Cancelled {
			return endFlow()
		}
		f.pickedPath = msg.Path
		f.state = statePwdEntry
		f.passwordAttempt = 0
		// Push password entry modal
		return func() tea.Msg {
			return pushModalMsg{modal: &passwordEntryModal{}}
		}

	case pwdEnteredMsg:
		f.state = statePreload
		password := msg.Password
		path := f.pickedPath
		// Attempt to load the vault in a background command
		return func() tea.Msg {
			// Call RecoverOrphans silently
			_ = storage.RecoverOrphans(path)
			// Call Load
			cofre, metadata, err := storage.Load(path, password)
			crypto.Wipe(password)
			if err != nil {
				if errors.Is(err, crypto.ErrAuthFailed) {
					f.passwordAttempt++
					if f.passwordAttempt >= 5 {
						f.messages.Show(MsgError, "Limite de tentativas excedido", 5, false)
						return endFlow()
					}
					// Return to password entry with retry
					f.state = statePwdEntry
					f.messages.Show(MsgWarn, "Senha incorreta. Tente novamente", 3, false)
					return pushModalMsg{modal: &passwordEntryModal{}}
				}
				// Other errors: show error dialog
				errMsg := "Não foi possível abrir o cofre. Arquivo corrompido ou inacessível."
				if errors.Is(err, storage.ErrInvalidMagic) {
					errMsg = "Arquivo não é um cofre válido."
				} else if errors.Is(err, storage.ErrVersionTooNew) {
					errMsg = "Versão do cofre não é suportada por esta versão do Abditum."
				} else if errors.Is(err, storage.ErrCorrupted) {
					errMsg = "Cofre corrompido e não pode ser recuperado."
				}
				f.messages.Show(MsgError, errMsg, 5, false)
				return endFlow()
			}
			// Success: store metadata and notify rootModel
			f.vaultMetadata = metadata
			f.state = stateDone
			if cofre != nil && f.mgr != nil {
				// Store the cofre in the manager (this will be done by rootModel actually)
				_ = cofre // just to acknowledge we received it
			}
			return vaultOpenedMsg{Path: path}
		}

	case flowCancelledMsg:
		return endFlow()

	default:
		return nil
	}
}

// View is not used in flow-driven interface; modals handle visualization.
func (f *openVaultFlow) View(width, height int) string {
	return ""
}
