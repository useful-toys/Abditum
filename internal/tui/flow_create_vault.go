package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/vault"
)

// createVaultFlow implements the state machine for creating a new vault.
// States: stateCheckDirty -> statePickFile -> stateCheckOverwrite -> statePwdCreate -> stateSaveNew -> done.
// If cliPath is set (CLI fast-path), skip dirty-check and file picker, go directly to password creation.
type createVaultFlow struct {
	state      int // current state in the state machine
	width      int
	height     int
	theme      *Theme
	mgr        *vault.Manager
	messages   *MessageManager
	actions    *ActionManager
	targetPath string // path selected by user for new vault
	cliPath    string // path provided via CLI (D-CLI-02 fast-path)
}

// State constants for createVaultFlow (reuses common states)
const (
	stateCheckOverwrite = iota + 10 // offset to avoid collision with openVaultFlow states
	statePwdCreate
	stateSaveNew
	stateStrengthCheck // intermediate: weak password decision pending
)

// createVaultSaveBeforeMsg — emitted by "S Salvar" in dirty-check dialog when
// creating a new vault. Triggers saving the current vault before opening the file picker.
type createVaultSaveBeforeMsg struct{}

func (createVaultSaveBeforeMsg) isModalResult() {}

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
	// D-CLI-02: CLI fast-path — skip dirty-check and file picker when cliPath is set
	if f.cliPath != "" {
		f.targetPath = f.cliPath
		f.state = statePwdCreate
		return func() tea.Msg {
			return pushModalMsg{modal: &passwordCreateModal{}}
		}
	}

	if f.mgr != nil && f.mgr.IsModified() {
		f.state = stateCheckDirty
		// Desvio 6: Decision(SeverityAlert) with save/discard/back actions (not Acknowledge).
		return func() tea.Msg {
			return Decision(SeverityAlert, "Criar novo cofre",
				"Cofre modificado. Salvar ou descartar?",
				DecisionAction{Key: "S", Label: "Salvar", Default: true,
					Cmd: func() tea.Msg { return createVaultSaveBeforeMsg{} }},
				[]DecisionAction{
					{Key: "D", Label: "Descartar",
						Cmd: FilePicker("Salvar cofre", FilePickerSave, ".abditum", f.messages)},
				},
				DecisionAction{Key: "Esc", Label: "Voltar"})
		}
	}
	f.state = statePickFile
	// Push file picker modal (in save mode)
	return FilePicker("Salvar cofre", FilePickerSave, ".abditum", f.messages)
}

// Update processes messages from modals and transitions states.
func (f *createVaultFlow) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case createVaultSaveBeforeMsg:
		// User chose "S Salvar" in dirty-check dialog: save current vault then open file picker.
		mgr := f.mgr
		messages := f.messages
		// D-SIG-02: emit MessageBusy synchronously before returning background Cmd
		messages.Show(MessageBusy, "Salvando cofre...", 0, false)
		return func() tea.Msg {
			if err := mgr.Salvar(); err != nil {
				messages.Show(MessageError, "Não foi possível salvar o cofre.", 5, false)
				return endFlowMsg{}
			}
			return FilePicker("Salvar cofre", FilePickerSave, ".abditum", messages)()
		}

	case filePickerResult:
		if msg.Cancelled {
			return endFlow()
		}
		f.targetPath = msg.Path
		f.state = stateCheckOverwrite
		// Check if file exists
		if _, err := os.Stat(f.targetPath); err == nil {
			// Desvio 7: File exists — show overwrite confirmation with interpolated filename,
			// SeverityAlert (not Destructive), key S (not Enter), and "I Outro caminho" middle action.
			baseName := strings.TrimSuffix(filepath.Base(f.targetPath), ".abditum")
			overwriteMsg := fmt.Sprintf("Arquivo '%s' já existe. Sobrescrever?", baseName)
			return func() tea.Msg {
				return Decision(SeverityAlert, "Criar novo cofre",
					overwriteMsg,
					DecisionAction{Key: "S", Label: "Sobrescrever", Default: true,
						Cmd: func() tea.Msg { return overwriteConfirmedMsg{} }},
					[]DecisionAction{
						{Key: "I", Label: "Outro caminho",
							Cmd: func() tea.Msg { return overwriteCancelledMsg{} }},
					},
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
			// Desvio 8: Default action is now P Prosseguir (not R Revisar).
			// Esc cancels the entire flow (endFlow), not proceeds.
			return func() tea.Msg {
				return Decision(SeverityAlert, "Criar novo cofre",
					"Senha é fraca. Prosseguir ou revisar?",
					DecisionAction{Key: "P", Label: "Prosseguir", Default: true,
						Cmd: func() tea.Msg { return weakPwdProceedMsg{Password: password} }},
					[]DecisionAction{
						{Key: "R", Label: "Revisar",
							Cmd: func() tea.Msg {
								crypto.Wipe(password) // discard weak password on revise
								return weakPwdReviseMsg{}
							}},
					},
					DecisionAction{Key: "Esc", Label: "Voltar",
						Cmd: func() tea.Msg {
							crypto.Wipe(password) // discard password, cancel entire flow
							return endFlowMsg{}
						}})
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
		// User chose "Voltar" or "Outro caminho" - return to file picker
		f.state = statePickFile
		return FilePicker("Salvar cofre", FilePickerSave, ".abditum", f.messages)

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
	// D-SIG-02: emit MessageBusy synchronously before returning background Cmd
	f.messages.Show(MessageBusy, "Criando cofre...", 0, false)
	return func() tea.Msg {
		newVault := vault.NovoCofre()
		newVault.InicializarConteudoPadrao()
		err := storage.SaveNew(path, newVault, password)
		crypto.Wipe(password)
		if err != nil {
			f.messages.Show(MessageError, "Não foi possível salvar o cofre.", 5, false)
			return endFlow()
		}
		// Compute file metadata for external change detection baseline
		metadata, err := storage.ComputeFileMetadata(path)
		if err != nil {
			// Metadata failure is non-fatal: vault was saved successfully
			metadata = storage.FileMetadata{}
		}
		// D-SUC-02: emit MessageSuccess before vaultOpenedMsg
		f.messages.Show(MessageSuccess, "✓ Cofre criado com sucesso", 5, false)
		return vaultOpenedMsg{Path: path, Metadata: metadata}
	}
}
