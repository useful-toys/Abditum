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
	stateStrengthCheck // intermediate: weak password decision pending
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
					DecisionAction{Key: "Enter", Label: "Sobrescrever", Default: true,
						Cmd: func() tea.Msg { return overwriteConfirmedMsg{} }},
					nil,
					DecisionAction{Key: "Esc", Label: "Voltar",
						Cmd: func() tea.Msg { return overwriteCancelledMsg{} }})
			}
		}
		// File doesn't exist - proceed to password creation
		f.state = statePwdCreate
		return func() tea.Msg {
			return pushModalMsg{modal: &passwordCreateModal{}}
		}

	case pwdCreatedMsg:
		password := msg.Password
		// Spec Gap 2.1: Gate on password strength before proceeding to save
		if crypto.EvaluatePasswordStrength(password) < crypto.StrengthStrong {
			f.state = stateStrengthCheck
			return func() tea.Msg {
				return Decision(SeverityAlert, "Senha fraca",
					"A senha não atende aos critérios de segurança recomendados.\n\nDeseja prosseguir mesmo assim ou revisar a senha?",
					DecisionAction{Key: "R", Label: "Revisar", Default: true,
						Cmd: func() tea.Msg {
							crypto.Wipe(password) // discard weak password on revise
							return weakPwdReviseMsg{}
						}},
					nil,
					DecisionAction{Key: "Esc", Label: "Prosseguir",
						Cmd: func() tea.Msg { return weakPwdProceedMsg{Password: password} }})
			}
		}
		// Strong password: proceed directly to save
		return f.saveVault(password)

	case weakPwdProceedMsg:
		// User chose to proceed despite weak password
		return f.saveVault(msg.Password)

	case weakPwdReviseMsg:
		// User chose to revise - return to password creation modal
		f.state = statePwdCreate
		return func() tea.Msg {
			return pushModalMsg{modal: &passwordCreateModal{}}
		}

	case overwriteConfirmedMsg:
		// User confirmed overwrite - proceed to password creation
		f.state = statePwdCreate
		return func() tea.Msg {
			return pushModalMsg{modal: &passwordCreateModal{}}
		}

	case overwriteCancelledMsg:
		// User chose "Voltar" - return to file picker
		f.state = statePickFile
		return func() tea.Msg {
			return pushModalMsg{modal: &filePickerModal{mode: FilePickerFile}}
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

// saveVault creates and persists a new vault at f.targetPath using the provided
// password. It transitions to stateSaveNew and runs the I/O in a background Cmd.
func (f *createVaultFlow) saveVault(password []byte) tea.Cmd {
	f.state = stateSaveNew
	path := f.targetPath
	return func() tea.Msg {
		newVault := vault.NovoCofre()
		newVault.InicializarConteudoPadrao()
		err := storage.SaveNew(path, newVault, password)
		crypto.Wipe(password)
		if err != nil {
			f.messages.Show(MsgError, "Não foi possível salvar o cofre.", 5, false)
			return endFlow()
		}
		return vaultOpenedMsg{Path: path}
	}
}
