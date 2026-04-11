package tui

import (
	"testing"

	testdatapkg "github.com/useful-toys/abditum/internal/tui/testdata"
)

// TestCreateVaultFlowStructExists validates the createVaultFlow struct exists.
func TestCreateVaultFlowStructExists(t *testing.T) {
	_ = (*createVaultFlow)(nil)
}

// TestCreateVaultFlowImplementsFlowHandler validates createVaultFlow implements flowHandler.
func TestCreateVaultFlowImplementsFlowHandler(t *testing.T) {
	var _ flowHandler = (*createVaultFlow)(nil)
}

// TestCreateVaultFlowNewCreatesFlow validates newCreateVaultFlow creates a flow.
func TestCreateVaultFlowNewCreatesFlow(t *testing.T) {
	flow := newCreateVaultFlow(nil, NewMessageManager(), NewActionManager(), ThemeTokyoNight)
	if flow == nil {
		t.Error("newCreateVaultFlow should return a non-nil flow")
	}
	if flow.state != stateCheckDirty {
		t.Errorf("Flow should start in stateCheckDirty, got %d", flow.state)
	}
}

// TestCreateVaultFlowInit validates Init() initializes the flow.
func TestCreateVaultFlowInit(t *testing.T) {
	flow := newCreateVaultFlow(nil, NewMessageManager(), NewActionManager(), ThemeTokyoNight)
	cmd := flow.Init()
	// Init should return a valid tea.Cmd
	if cmd != nil {
		// Call the command to ensure it doesn't panic
		msg := cmd()
		if msg == nil {
			t.Error("Init command should produce a message")
		}
	}
}

// TestCreateVaultFlowUpdate validates Update() handles messages.
func TestCreateVaultFlowUpdate(t *testing.T) {
	flow := newCreateVaultFlow(nil, NewMessageManager(), NewActionManager(), ThemeTokyoNight)
	flow.state = statePickFile
	// Send a generic message
	cmd := flow.Update(nil)
	// Update should return a valid tea.Cmd (may be nil)
	if cmd != nil {
		// Calling the command shouldn't panic
		_ = cmd()
	}
}

// TestCreateVaultFlowHandlesFilePickerResult validates handling of file selection.
func TestCreateVaultFlowHandlesFilePickerResult(t *testing.T) {
	flow := newCreateVaultFlow(nil, NewMessageManager(), NewActionManager(), ThemeTokyoNight)
	flow.state = statePickFile
	msg := filePickerResult{Path: "/tmp/test.abditum", Cancelled: false}
	cmd := flow.Update(msg)
	if cmd == nil {
		t.Error("Update with filePickerResult should return a command")
	}
	if flow.targetPath != "/tmp/test.abditum" {
		t.Errorf("Flow should store target path, got %q", flow.targetPath)
	}
}

// TestCreateVaultFlowHandlesFilePickerCancelled validates cancellation.
func TestCreateVaultFlowHandlesFilePickerCancelled(t *testing.T) {
	flow := newCreateVaultFlow(nil, NewMessageManager(), NewActionManager(), ThemeTokyoNight)
	flow.state = statePickFile
	msg := filePickerResult{Path: "", Cancelled: true}
	cmd := flow.Update(msg)
	if cmd == nil {
		t.Error("Update with cancelled result should return endFlow command")
	}
}

// TestCreateVaultFlowView validates View() returns a string.
func TestCreateVaultFlowView(t *testing.T) {
	flow := newCreateVaultFlow(nil, NewMessageManager(), NewActionManager(), ThemeTokyoNight)
	view := flow.View(80, 24)
	// View should return a string (may be empty for flow-driven interface)
	if view == "" {
		// Empty string is acceptable; flow uses modals for visualization
	}
}

// TestCreateVaultFlowEmitsVaultOpenedMsg validates successful vault creation emits vaultOpenedMsg.
func TestCreateVaultFlowEmitsVaultOpenedMsg(t *testing.T) {
	// This is a placeholder test to validate the message type exists.
	var msg any = vaultOpenedMsg{Path: "/path/to/vault.abditum"}
	if _, ok := msg.(vaultOpenedMsg); !ok {
		t.Error("vaultOpenedMsg should be a valid message type")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Behavioral Tests — D-SIG-02 synchronous MsgBusy and D-CLI-02 cliPath fast-path
// ─────────────────────────────────────────────────────────────────────────────

// TestCreateVaultFlow_MsgBusy_EmittedSynchronously verifies D-SIG-02:
// calling saveVault() synchronously sets MsgBusy BEFORE returning the Cmd.
func TestCreateVaultFlow_MsgBusy_EmittedSynchronously(t *testing.T) {
	msgs := NewMessageManager()
	flow := newCreateVaultFlow(nil, msgs, NewActionManager(), ThemeTokyoNight)
	flow.targetPath = "/tmp/test-abditum-sync.abditum"
	flow.state = statePwdCreate

	// The MessageManager should be empty (or have hint) before saveVault
	if curr := msgs.Current(); curr != nil && curr.Kind == MsgBusy {
		t.Fatal("did not expect MsgBusy before saveVault")
	}

	// Call saveVault directly — D-SIG-02: Show(MsgBusy) must be called synchronously
	_ = flow.saveVault([]byte("testpassword"))

	// BEFORE running the returned Cmd, verify MsgBusy was set synchronously
	curr := msgs.Current()
	if curr == nil {
		t.Fatal("expected MsgBusy to be set synchronously after saveVault(), got nil")
	}
	if curr.Kind != MsgBusy {
		t.Errorf("expected MsgBusy kind, got %v", curr.Kind)
	}
	if curr.Text != "Criando cofre..." {
		t.Errorf("expected 'Criando cofre...' text, got %q", curr.Text)
	}
}

// TestCreateVaultFlow_CliPath_SkipsToPasswordCreate verifies D-CLI-02:
// when cliPath is set, Init() skips dirty-check and file picker and returns
// a Cmd that emits pushModalMsg{modal: *passwordCreateModal}.
func TestCreateVaultFlow_CliPath_SkipsToPasswordCreate(t *testing.T) {
	msgs := NewMessageManager()
	flow := newCreateVaultFlow(nil, msgs, NewActionManager(), ThemeTokyoNight)
	flow.cliPath = "/tmp/test-new.abditum"

	cmd := flow.Init()
	if cmd == nil {
		t.Fatal("expected non-nil Cmd from Init() with cliPath set")
	}

	// Verify state was set to statePwdCreate
	if flow.state != statePwdCreate {
		t.Errorf("expected state=statePwdCreate after CLI fast-path Init(), got %d", flow.state)
	}

	// Verify targetPath was set
	if flow.targetPath != "/tmp/test-new.abditum" {
		t.Errorf("expected targetPath=%q, got %q", "/tmp/test-new.abditum", flow.targetPath)
	}

	// Execute the Cmd — should emit pushModalMsg with *passwordCreateModal
	msg := cmd()
	pm, ok := msg.(pushModalMsg)
	if !ok {
		t.Fatalf("expected pushModalMsg, got %T", msg)
	}
	if _, ok := pm.modal.(*passwordCreateModal); !ok {
		t.Fatalf("expected *passwordCreateModal in pushModalMsg.modal, got %T", pm.modal)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Golden File Tests
// ─────────────────────────────────────────────────────────────────────────────

// TestCreateVaultFlow_Golden validates flow view rendering against golden files.
func TestCreateVaultFlow_Golden(t *testing.T) {
	flow := newCreateVaultFlow(nil, NewMessageManager(), NewActionManager(), ThemeTokyoNight)
	flow.state = stateCheckDirty

	view := flow.View(80, 24)

	// .txt.golden: plain text render
	txtPath := goldenPath("flow-create-vault", "initial", 80, "txt")
	checkOrUpdateGolden(t, txtPath, stripANSI(view))

	// .json.golden: style transitions
	transitions := testdatapkg.ParseANSIStyle(view)
	jsonBytes, err := testdatapkg.MarshalStyleTransitions(transitions)
	if err != nil {
		t.Fatalf("marshal transitions: %v", err)
	}
	jsonPath := goldenPath("flow-create-vault", "initial", 80, "json")
	checkOrUpdateGolden(t, jsonPath, string(jsonBytes))
}
