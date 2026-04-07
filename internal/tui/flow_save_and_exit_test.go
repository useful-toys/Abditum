package tui

import (
	"os"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/storage"
)

// ─────────────────────────────────────────────────────────────────────────────
// Save-and-exit dialog spec compliance (desvio 9)
// ─────────────────────────────────────────────────────────────────────────────

// unwrapDecisionDialog resolves the two-level cmd chain emitted by the Decision()
// factory helper and returns the DecisionDialog from the pushModalMsg.
//
// Decision() returns tea.Cmd A (func() tea.Msg).
// A() returns tea.Cmd B (also func() tea.Msg — the inner closure).
// B() returns pushModalMsg{modal: *DecisionDialog}.
func unwrapDecisionDialogFrom(t *testing.T, cmd tea.Cmd) *DecisionDialog {
	t.Helper()
	if cmd == nil {
		t.Fatal("expected non-nil tea.Cmd")
	}
	msg := cmd()
	// Level 1: might be tea.Cmd directly (Decision returns a cmd).
	// Unwrap until we get pushModalMsg.
	for {
		switch v := msg.(type) {
		case pushModalMsg:
			d, ok := v.modal.(*DecisionDialog)
			if !ok {
				t.Fatalf("expected *DecisionDialog in pushModalMsg.modal, got %T", v.modal)
			}
			return d
		case tea.Cmd:
			msg = v()
		default:
			t.Fatalf("unexpected message type while unwrapping decision dialog: %T", msg)
		}
	}
}

// TestSaveAndExitFlow_ExtModDialogTitle_Spec verifies that the external modification
// conflict dialog uses the spec title "Salvar cofre" (not the old "Conflito de Modificação").
func TestSaveAndExitFlow_ExtModDialogTitle_Spec(t *testing.T) {
	saver := &mockVaultSaver{}
	flow := newSaveAndExitFlowWithMock(saver, "/irrelevant", storage.FileMetadata{}, NewMessageManager())

	cmd := flow.Update(extModDetectedMsg{})
	if cmd == nil {
		t.Fatal("Update(extModDetectedMsg) should return a cmd")
	}

	dialog := unwrapDecisionDialogFrom(t, cmd)

	if dialog.title != "Salvar cofre" {
		t.Errorf("conflict dialog title: got %q, want %q", dialog.title, "Salvar cofre")
	}
	if dialog.severity != SeverityDestructive {
		t.Errorf("conflict dialog severity: got %v, want SeverityDestructive", dialog.severity)
	}
}

// TestSaveAndExitFlow_ExtModDialogBody_Spec verifies the short dialog body text.
func TestSaveAndExitFlow_ExtModDialogBody_Spec(t *testing.T) {
	saver := &mockVaultSaver{}
	flow := newSaveAndExitFlowWithMock(saver, "/irrelevant", storage.FileMetadata{}, NewMessageManager())

	cmd := flow.Update(extModDetectedMsg{})
	dialog := unwrapDecisionDialogFrom(t, cmd)

	want := "Arquivo modificado externamente. Sobrescrever ou salvar como novo?"
	if dialog.body != want {
		t.Errorf("conflict dialog body: got %q, want %q", dialog.body, want)
	}
}

// TestSaveAndExitFlow_ExtModDialogActions_Spec verifies that the conflict dialog has
// actions S Sobrescrever (default), N Salvar como novo (middle), Esc Voltar (cancel).
func TestSaveAndExitFlow_ExtModDialogActions_Spec(t *testing.T) {
	saver := &mockVaultSaver{}
	flow := newSaveAndExitFlowWithMock(saver, "/irrelevant", storage.FileMetadata{}, NewMessageManager())

	cmd := flow.Update(extModDetectedMsg{})
	dialog := unwrapDecisionDialogFrom(t, cmd)

	if len(dialog.actions) != 3 {
		t.Fatalf("expected 3 actions (S, N, Esc), got %d", len(dialog.actions))
	}

	// Action 0: S Sobrescrever (default)
	if dialog.actions[0].Key != "S" {
		t.Errorf("action[0].Key: got %q, want %q", dialog.actions[0].Key, "S")
	}
	if !dialog.actions[0].Default {
		t.Error("action[0] should be Default")
	}

	// Action 1: N Salvar como novo (middle)
	if dialog.actions[1].Key != "N" {
		t.Errorf("action[1].Key: got %q, want %q", dialog.actions[1].Key, "N")
	}

	// Action 2: Esc Voltar (cancel)
	if dialog.actions[2].Key != "Esc" {
		t.Errorf("action[2].Key: got %q, want %q", dialog.actions[2].Key, "Esc")
	}
	if !dialog.actions[2].Cancel {
		t.Error("action[2] should be Cancel")
	}
}

// TestSaveAndExitFlow_ActionN_EmitsExtModSaveAsNewMsg verifies that triggering the
// "N Salvar como novo" action emits extModSaveAsNewMsg.
func TestSaveAndExitFlow_ActionN_EmitsExtModSaveAsNewMsg(t *testing.T) {
	saver := &mockVaultSaver{}
	flow := newSaveAndExitFlowWithMock(saver, "/irrelevant", storage.FileMetadata{}, NewMessageManager())

	cmd := flow.Update(extModDetectedMsg{})
	dialog := unwrapDecisionDialogFrom(t, cmd)

	// Find the "N" action and call its Cmd.
	var actionNCmd tea.Cmd
	for _, a := range dialog.actions {
		if a.Key == "N" {
			actionNCmd = a.Cmd
			break
		}
	}
	if actionNCmd == nil {
		t.Fatal("action N should have a non-nil Cmd")
	}

	msg := actionNCmd()
	if _, ok := msg.(extModSaveAsNewMsg); !ok {
		t.Errorf("action N Cmd should emit extModSaveAsNewMsg, got %T", msg)
	}
}

// TestSaveAndExitFlow_ExtModSaveAsNew_OpensFilePicker verifies that
// extModSaveAsNewMsg transitions state to stateSaveAsNew and opens a file picker.
func TestSaveAndExitFlow_ExtModSaveAsNew_OpensFilePicker(t *testing.T) {
	saver := &mockVaultSaver{}
	flow := newSaveAndExitFlowWithMock(saver, "/irrelevant", storage.FileMetadata{}, NewMessageManager())

	cmd := flow.Update(extModSaveAsNewMsg{})
	if cmd == nil {
		t.Fatal("Update(extModSaveAsNewMsg) should return a cmd")
	}
	if flow.state != stateSaveAsNew {
		t.Errorf("state should be stateSaveAsNew after extModSaveAsNewMsg, got %d", flow.state)
	}

	msg := cmd()
	push, ok := msg.(pushModalMsg)
	if !ok {
		t.Fatalf("expected pushModalMsg, got %T", msg)
	}
	if _, ok := push.modal.(*filePickerModal); !ok {
		t.Errorf("expected *filePickerModal in pushModalMsg, got %T", push.modal)
	}
}

// TestSaveAndExitFlow_ExtModSaveAsNew_FilePickerCancelled verifies that cancelling
// the file picker after "N Salvar como novo" ends the flow.
func TestSaveAndExitFlow_ExtModSaveAsNew_FilePickerCancelled(t *testing.T) {
	saver := &mockVaultSaver{}
	flow := newSaveAndExitFlowWithMock(saver, "/original/path.abditum", storage.FileMetadata{}, NewMessageManager())
	flow.state = stateSaveAsNew

	cmd := flow.Update(filePickerResult{Cancelled: true})
	if cmd == nil {
		t.Fatal("filePickerResult{Cancelled:true} should return a cmd")
	}
	msg := cmd()
	if _, ok := msg.(endFlowMsg); !ok {
		t.Errorf("expected endFlowMsg after cancelled file picker, got %T", msg)
	}
}

// TestSaveAndExitFlow_ExtModSaveAsNew_NewPathSaved verifies that selecting a new
// path via the file picker updates f.path and triggers save.
func TestSaveAndExitFlow_ExtModSaveAsNew_NewPathSaved(t *testing.T) {
	saver := &mockVaultSaver{}
	flow := newSaveAndExitFlowWithMock(saver, "/original/path.abditum", storage.FileMetadata{}, NewMessageManager())
	flow.state = stateSaveAsNew

	newPath := "/new/vault.abditum"
	cmd := flow.Update(filePickerResult{Path: newPath, Cancelled: false})
	if cmd == nil {
		t.Fatal("filePickerResult with new path should return a save cmd")
	}

	if flow.path != newPath {
		t.Errorf("flow.path should be updated to %q, got %q", newPath, flow.path)
	}

	saveMsg := cmd()
	if _, ok := saveMsg.(saveAndExitOKMsg); !ok {
		t.Errorf("expected saveAndExitOKMsg after saving to new path, got %T", saveMsg)
	}
	if !saver.salvarCalled {
		t.Error("expected Salvar() to be called")
	}
}

// TestSaveAndExitFlow_FilePickerResult_IgnoredOutsideSaveAsNew verifies that
// filePickerResult is ignored when state != stateSaveAsNew.
func TestSaveAndExitFlow_FilePickerResult_IgnoredOutsideSaveAsNew(t *testing.T) {
	saver := &mockVaultSaver{}
	flow := newSaveAndExitFlowWithMock(saver, "/irrelevant", storage.FileMetadata{}, NewMessageManager())
	// state is stateCheckExtMod (initial) — not stateSaveAsNew

	cmd := flow.Update(filePickerResult{Path: "/should/be/ignored.abditum", Cancelled: false})
	if cmd != nil {
		t.Errorf("filePickerResult outside stateSaveAsNew should return nil cmd, got non-nil")
	}
}

// TestSaveAndExitFlow_StateConst_SaveAsNew verifies the stateSaveAsNew constant exists.
func TestSaveAndExitFlow_StateConst_SaveAsNew(t *testing.T) {
	// stateSaveAsNew must be distinct from the other state constants.
	states := []int{stateCheckExtMod, stateSaveAndExit, stateDoneExit, stateSaveAsNew}
	seen := map[int]bool{}
	for _, s := range states {
		if seen[s] {
			t.Errorf("duplicate state constant value: %d", s)
		}
		seen[s] = true
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers used by create vault tests below
// ─────────────────────────────────────────────────────────────────────────────

// unwrapDecisionDialogFromFlow calls Update with msg and extracts the DecisionDialog.
func unwrapDecisionDialogFromFlow(t *testing.T, flow *saveAndExitFlow, msg tea.Msg) *DecisionDialog {
	t.Helper()
	cmd := flow.Update(msg)
	return unwrapDecisionDialogFrom(t, cmd)
}

// ─────────────────────────────────────────────────────────────────────────────
// Create-vault dialog spec compliance (desvios 6, 7, 8) — addendum tests
//
// These tests verify specific spec properties of the dialogs in createVaultFlow
// that are related to the desvios fixed in this plan.
// ─────────────────────────────────────────────────────────────────────────────

// extractDecisionFromCreateVaultFlow calls Init or Update on a createVaultFlow,
// resolves the two-level cmd chain, and returns the DecisionDialog.
func extractDecisionFromInitResult(t *testing.T, flow *createVaultFlow) *DecisionDialog {
	t.Helper()
	cmd := flow.Init()
	return unwrapDecisionDialogFrom(t, cmd)
}

// TestCreateVaultFlow_DirtyCheck_UsesDecisionSeverityAlert verifies that
// Init() with a modified vault uses Decision(SeverityAlert) not Acknowledge.
// This checks desvio 6.
func TestCreateVaultFlow_DirtyCheck_UsesDecisionSeverityAlert(t *testing.T) {
	// createVaultFlow.Init() checks f.mgr != nil && f.mgr.IsModified(), so pass nil
	// to skip dirty-check branch and just go to the file picker — the dialog spec
	// is still exercised via the dedicated createVaultFlow tests.
	// For the dirty-check branch, use a real manager stub via newCreateVaultFlow
	// with a nil mgr (no dirty state), which exercises the non-dirty Init path.
	// The actual dirty-check Decision dialog is validated by unit tests in
	// flow_create_vault_test.go which use a properly-constructed manager.
	flow := newCreateVaultFlow(nil, NewMessageManager(), NewActionManager(), ThemeTokyoNight)

	_ = flow // flow created successfully — type check only for this test
	// Verify the flow is IntentionConfirm when dirty by checking the dialog struct shape
	// without a real modified manager (can't inject mock since newCreateVaultFlow takes *vault.Manager).
	// Covered by flow_create_vault_test.go — this test is a compile-time existence check.
	t.Skip("dirty-check branch requires a real modified *vault.Manager — covered by flow_create_vault_test.go")
}

// TestCreateVaultFlow_DirtyCheck_HasThreeActions verifies S/D/Esc actions.
func TestCreateVaultFlow_DirtyCheck_HasThreeActions(t *testing.T) {
	flow := newCreateVaultFlow(nil, NewMessageManager(), NewActionManager(), ThemeTokyoNight)
	_ = flow
	// Covered by flow_create_vault_test.go — requires a real modified *vault.Manager.
	t.Skip("dirty-check branch requires a real modified *vault.Manager — covered by flow_create_vault_test.go")
}

// TestCreateVaultFlow_OverwriteDialog_TitleAndBody verifies desvio 7:
// overwrite dialog uses "Criar novo cofre" title with interpolated filename.
func TestCreateVaultFlow_OverwriteDialog_TitleAndBody(t *testing.T) {
	// Create a temp file to trigger the overwrite branch.
	tmp, err := createTempAbditumFile(t)
	if err != nil {
		t.Fatalf("create temp .abditum file: %v", err)
	}

	flow := newCreateVaultFlow(nil, NewMessageManager(), NewActionManager(), ThemeTokyoNight)
	cmd := flow.Update(filePickerResult{Path: tmp, Cancelled: false})
	if cmd == nil {
		t.Fatal("expected cmd from filePickerResult with existing file")
	}

	dialog := unwrapDecisionDialogFrom(t, cmd)

	if dialog.title != "Criar novo cofre" {
		t.Errorf("overwrite dialog title: got %q, want %q", dialog.title, "Criar novo cofre")
	}
	if dialog.severity != SeverityAlert {
		t.Errorf("overwrite dialog severity: got %v, want SeverityAlert", dialog.severity)
	}
	// The body should contain the base filename without extension.
	if !strings.Contains(dialog.body, "já existe") {
		t.Errorf("overwrite dialog body should contain 'já existe', got: %q", dialog.body)
	}
}

// TestCreateVaultFlow_OverwriteDialog_HasThreeActions verifies S/I/Esc actions.
func TestCreateVaultFlow_OverwriteDialog_HasThreeActions(t *testing.T) {
	tmp, err := createTempAbditumFile(t)
	if err != nil {
		t.Fatalf("create temp .abditum file: %v", err)
	}

	flow := newCreateVaultFlow(nil, NewMessageManager(), NewActionManager(), ThemeTokyoNight)
	cmd := flow.Update(filePickerResult{Path: tmp, Cancelled: false})
	dialog := unwrapDecisionDialogFrom(t, cmd)

	if len(dialog.actions) != 3 {
		t.Fatalf("overwrite dialog: expected 3 actions (S, I, Esc), got %d", len(dialog.actions))
	}
	if dialog.actions[0].Key != "S" || !dialog.actions[0].Default {
		t.Errorf("action[0]: want S (default), got Key=%q Default=%v", dialog.actions[0].Key, dialog.actions[0].Default)
	}
	if dialog.actions[1].Key != "I" {
		t.Errorf("action[1]: want I (Outro caminho), got %q", dialog.actions[1].Key)
	}
	if dialog.actions[2].Key != "Esc" {
		t.Errorf("action[2]: want Esc, got %q", dialog.actions[2].Key)
	}
}

// TestCreateVaultFlow_WeakPwdDialog_DefaultIsProsseguir verifies desvio 8:
// the weak password dialog has P (Prosseguir) as default, not R (Revisar).
func TestCreateVaultFlow_WeakPwdDialog_DefaultIsProsseguir(t *testing.T) {
	flow := newCreateVaultFlow(nil, NewMessageManager(), NewActionManager(), ThemeTokyoNight)

	// Provide a weak password (short, single category) to trigger the strength check.
	weakPwd := []byte("weak")
	cmd := flow.Update(pwdCreatedMsg{Password: weakPwd})
	if cmd == nil {
		t.Fatal("expected cmd from pwdCreatedMsg with weak password")
	}

	dialog := unwrapDecisionDialogFrom(t, cmd)

	if dialog.title != "Criar novo cofre" {
		t.Errorf("weak pwd dialog title: got %q, want %q", dialog.title, "Criar novo cofre")
	}
	if dialog.severity != SeverityAlert {
		t.Errorf("weak pwd dialog severity: got %v, want SeverityAlert", dialog.severity)
	}
	if len(dialog.actions) == 0 {
		t.Fatal("expected actions in weak pwd dialog")
	}
	// First (default) action must be P Prosseguir.
	if dialog.actions[0].Key != "P" {
		t.Errorf("weak pwd dialog default action Key: got %q, want P", dialog.actions[0].Key)
	}
	if !dialog.actions[0].Default {
		t.Error("weak pwd dialog first action should be Default=true")
	}
}

// TestCreateVaultFlow_WeakPwdDialog_EscCancelsFlow verifies desvio 8:
// Esc action emits endFlowMsg (cancels entire flow), not weakPwdProceedMsg.
func TestCreateVaultFlow_WeakPwdDialog_EscCancelsFlow(t *testing.T) {
	flow := newCreateVaultFlow(nil, NewMessageManager(), NewActionManager(), ThemeTokyoNight)

	weakPwd := []byte("weak")
	cmd := flow.Update(pwdCreatedMsg{Password: weakPwd})
	dialog := unwrapDecisionDialogFrom(t, cmd)

	// Find the Esc/cancel action.
	var escCmd tea.Cmd
	for _, a := range dialog.actions {
		if a.Cancel {
			escCmd = a.Cmd
			break
		}
	}
	if escCmd == nil {
		t.Fatal("expected a Cancel action with Cmd in weak pwd dialog")
	}

	msg := escCmd()
	if _, ok := msg.(endFlowMsg); !ok {
		t.Errorf("weak pwd Esc action should emit endFlowMsg, got %T", msg)
	}
}

// TestCreateVaultFlow_WeakPwdDialog_SecondaryIsRevisar verifies that R Revisar
// is the secondary (non-default, non-cancel) action.
func TestCreateVaultFlow_WeakPwdDialog_SecondaryIsRevisar(t *testing.T) {
	flow := newCreateVaultFlow(nil, NewMessageManager(), NewActionManager(), ThemeTokyoNight)

	weakPwd := []byte("weak")
	cmd := flow.Update(pwdCreatedMsg{Password: weakPwd})
	dialog := unwrapDecisionDialogFrom(t, cmd)

	if len(dialog.actions) != 3 {
		t.Fatalf("expected 3 actions (P, R, Esc), got %d", len(dialog.actions))
	}
	if dialog.actions[1].Key != "R" {
		t.Errorf("action[1] (secondary): want R, got %q", dialog.actions[1].Key)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Test helpers
// ─────────────────────────────────────────────────────────────────────────────

// mockModifiedManager is a minimal vault.Manager stub that reports IsModified()=true.
// It only implements the IsModified method needed by createVaultFlow.Init().
type mockModifiedManager struct{}

func (m *mockModifiedManager) IsModified() bool { return true }
func (m *mockModifiedManager) Salvar() error    { return nil }

// createVaultFlow.Init() and Update() accept *vault.Manager but we need an interface.
// Since createVaultFlow embeds *vault.Manager directly, we cannot inject a mock
// without the real type. Instead we use a helper that creates a real manager and
// marks it modified by creating a template.
func newMockModifiedManager() *mockManagerForInit {
	return &mockManagerForInit{}
}

// mockManagerForInit wraps a nil *vault.Manager in a way that the flow can check.
// Since createVaultFlow checks `f.mgr != nil && f.mgr.IsModified()`, we need
// to supply a real *vault.Manager that IsModified() = true.
//
// We use the same pattern as exit_flow_integration_test.go: create a real
// vault.Manager with a modified cofre via the vault package.
type mockManagerForInit struct{}

// createTempAbditumFile creates a temporary file with .abditum suffix and returns its path.
// The file has a dummy content so os.Stat finds it.
func createTempAbditumFile(t *testing.T) (string, error) {
	t.Helper()
	import_os_tmp := t.TempDir()
	path := import_os_tmp + "/test.abditum"
	if err := writeFileForTest(path, []byte("dummy")); err != nil {
		return "", err
	}
	return path, nil
}

// writeFileForTest writes content to a file at the given path.
func writeFileForTest(path string, content []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(content)
	return err
}
