package tui

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/vault"
)

// openVaultSaveBeforeMsg is emitted by the "S Salvar" action in the dirty-check
// Decision dialog when opening a vault. The flow saves the current vault and
// then continues to the file picker.
type openVaultSaveBeforeMsg struct{}

func (openVaultSaveBeforeMsg) isModalResult() {}

// openVaultFlow implements the state machine for opening an existing vault.
// States: stateCheckDirty -> statePickFile -> statePwdEntry -> statePreload -> done.
// If cliPath is set (CLI fast-path), skip file picker and go directly to password entry.
type openVaultFlow struct {
	state           int // current state in the state machine
	width           int
	height          int
	theme           *Theme
	mgr             *vault.Manager
	messages        *MessageManager
	actions         *ActionManager
	pickedPath      string // path selected by user
	cliPath         string // path provided via CLI (for fast-path)
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

// Init starts the flow. If there are unsaved changes, prompt the user with a
// Decision dialog (SeverityAlert) offering to save, discard, or cancel (spec:
// Desvio 3 — Abrir cofre dirty-check). If CLI path is set (fast-path), skip
// file selection and go directly to password entry. Otherwise, proceed to file
// selection.
func (f *openVaultFlow) Init() tea.Cmd {
	if f.mgr != nil && f.mgr.IsModified() {
		f.state = stateCheckDirty
		// Push Decision dialog for unsaved changes (spec: Desvio 3)
		return Decision(SeverityAlert, "Abrir cofre",
			"Cofre modificado. Salvar ou descartar?",
			DecisionAction{Key: "S", Label: "Salvar", Default: true,
				Cmd: func() tea.Msg { return openVaultSaveBeforeMsg{} }},
			[]DecisionAction{
				{Key: "D", Label: "Descartar",
					Cmd: FilePicker("Abrir cofre", FilePickerOpen, ".abditum", f.messages, f.theme)},
			},
			DecisionAction{Key: "Esc", Label: "Voltar"})
	}

	// CLI fast-path: skip file picker if cliPath is set
	if f.cliPath != "" {
		f.pickedPath = f.cliPath
		f.state = statePwdEntry
		f.passwordAttempt = 0
		return PasswordEntry("Abrir cofre")
	}

	f.state = statePickFile
	// Push file picker modal
	return FilePicker("Abrir cofre", FilePickerOpen, ".abditum", f.messages, f.theme)
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
		return PasswordEntry("Abrir cofre")

	case openVaultSaveBeforeMsg:
		// Save current vault, then proceed to file picker.
		// Capture pointers into locals to avoid closure aliasing.
		mgr := f.mgr
		messages := f.messages
		return func() tea.Msg {
			if err := mgr.Salvar(); err != nil {
				messages.Show(MsgError, "Não foi possível salvar o cofre.", 5, false)
				return endFlowMsg{}
			}
			return FilePicker("Abrir cofre", FilePickerOpen, ".abditum", messages, nil)()
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
						return endFlowMsg{}
					}
					// Wrong password (attempt < 5): show Acknowledge dialog (spec: Desvio 4)
					f.state = statePwdEntry
					return Acknowledge(SeverityError, "Abrir cofre",
						"Senha incorreta. Necessário tentar novamente.",
						PasswordEntry("Abrir cofre"))()
				}
				// File errors: show Acknowledge dialog (spec: Desvio 5).
				// Two cases: invalid format/version vs corrupted/generic.
				f.state = statePickFile
				if errors.Is(err, storage.ErrInvalidMagic) || errors.Is(err, storage.ErrVersionTooNew) {
					return Acknowledge(SeverityError, "Abrir cofre",
						"Arquivo inválido ou versão não suportada. Necessário corrigir.",
						func() tea.Msg { return FilePicker("Abrir cofre", FilePickerOpen, ".abditum", f.messages, f.theme)() })()
				}
				// ErrCorrupted and generic errors
				return Acknowledge(SeverityError, "Abrir cofre",
					"Arquivo corrompido ou inválido. Necessário fechar.",
					func() tea.Msg { return FilePicker("Abrir cofre", FilePickerOpen, ".abditum", f.messages, f.theme)() })()
			}
			// Success: store metadata and notify rootModel
			f.vaultMetadata = metadata
			f.state = stateDone
			if cofre != nil && f.mgr != nil {
				// Store the cofre in the manager (this will be done by rootModel actually)
				_ = cofre // just to acknowledge we received it
			}
			return vaultOpenedMsg{Path: path, Metadata: metadata}
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
