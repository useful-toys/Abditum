package tui

import (
	"errors"
	"os"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/vault"
)

// ─────────────────────────────────────────────────────────────────────────────
// Exit flow integration tests — Spec compliance for Fluxos 3, 4 and 5.
//
// Fluxo 3: Ctrl+Q on welcome screen (no vault) → confirmation dialog
// Fluxo 4: Ctrl+Q on vault area (clean/unmodified) → confirmation dialog
// Fluxo 5: Ctrl+Q on vault area with unsaved changes → save/discard/back dialog
// saveAndExitFlow: no external mod → save → quit
// saveAndExitFlow: external mod detected → conflict dialog
// saveAndExitFlow: conflict dialog "Sobrescrever e sair" → save → quit
// saveAndExitFlow: conflict dialog "Voltar" → end flow (stay in app)
// saveAndExitFlow: save error → error shown, end flow (stay in app)
// ─────────────────────────────────────────────────────────────────────────────

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

// unwrapDecisionCmd calls a Ctrl+Q cmd chain twice to reach the pushModalMsg:
// outer cmd → Decision() factory (tea.Cmd) → pushModalMsg.
func unwrapDecisionPush(t *testing.T, cmd tea.Cmd) pushModalMsg {
	t.Helper()
	if cmd == nil {
		t.Fatal("expected cmd, got nil")
	}
	inner := cmd()
	innerCmd, ok := inner.(tea.Cmd)
	if !ok {
		t.Fatalf("expected Decision() factory (tea.Cmd), got %T", inner)
	}
	msg := innerCmd()
	push, ok := msg.(pushModalMsg)
	if !ok {
		t.Fatalf("expected pushModalMsg, got %T", msg)
	}
	return push
}

// mockVaultSaver is a vaultSaver that records calls and optionally returns an error.
type mockVaultSaver struct {
	salvarCalled bool
	salvarErr    error
}

func (m *mockVaultSaver) Salvar() error {
	m.salvarCalled = true
	return m.salvarErr
}

// newSaveAndExitFlowWithMock creates a saveAndExitFlow with a mockVaultSaver injected.
// Uses the internal vaultSaver interface — only valid inside the tui package.
func newSaveAndExitFlowWithMock(saver vaultSaver, path string, meta storage.FileMetadata, messages *MessageManager) *saveAndExitFlow {
	return &saveAndExitFlow{
		state:    stateCheckExtMod,
		mgr:      saver,
		path:     path,
		metadata: meta,
		messages: messages,
		theme:    ThemeTokyoNight,
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Fluxo 3: Ctrl+Q on welcome screen (no vault open)
// ─────────────────────────────────────────────────────────────────────────────

// TestCtrlQ_Fluxo3_NoVault verifies that Ctrl+Q on the welcome screen (no vault
// open) pushes an exit confirmation dialog rather than quitting immediately.
func TestCtrlQ_Fluxo3_NoVault(t *testing.T) {
	m := NewRootModel()
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	if m.area != workAreaWelcome {
		t.Fatal("precondition: rootModel should be in workAreaWelcome")
	}

	_, cmd := m.Update(makeKeyPress("ctrl+q"))
	unwrapDecisionPush(t, cmd) // panics with t.Fatal if not pushModalMsg
}

// ─────────────────────────────────────────────────────────────────────────────
// Fluxo 4: Ctrl+Q on vault area with no unsaved changes
// ─────────────────────────────────────────────────────────────────────────────

// TestCtrlQ_Fluxo4_CleanVault verifies that Ctrl+Q on a vault area with no
// unsaved changes also shows a confirmation dialog (not an immediate quit).
func TestCtrlQ_Fluxo4_CleanVault(t *testing.T) {
	m := NewRootModel()
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Simulate vault opened; m.mgr remains nil (no unsaved changes path)
	m.Update(vaultOpenedMsg{Path: "/tmp/test.abditum"})
	if m.area != workAreaVault {
		t.Fatal("precondition: rootModel should be in workAreaVault")
	}

	_, cmd := m.Update(makeKeyPress("ctrl+q"))
	unwrapDecisionPush(t, cmd)
}

// ─────────────────────────────────────────────────────────────────────────────
// Fluxo 5: Ctrl+Q with unsaved changes
// ─────────────────────────────────────────────────────────────────────────────

// TestCtrlQ_Fluxo5_UnsavedChanges verifies that Ctrl+Q with a real vault.Manager
// that has unsaved changes shows a "Salvar / Descartar / Voltar" dialog.
func TestCtrlQ_Fluxo5_UnsavedChanges(t *testing.T) {
	m := NewRootModel()
	m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})

	// Create a real vault.Manager with a modified cofre
	cofre := vault.NovoCofre()
	repo := &mockVaultRepo{}
	mgr := vault.NewManager(cofre, repo)

	// Mark vault as modified via the public API: create a template
	if _, err := mgr.CriarModelo("test-template", []vault.CampoModelo{}); err != nil {
		t.Fatalf("CriarModelo failed: %v", err)
	}
	m.mgr = mgr

	if !m.mgr.IsModified() {
		t.Fatal("precondition: manager should report IsModified()=true")
	}

	_, cmd := m.Update(makeKeyPress("ctrl+q"))
	unwrapDecisionPush(t, cmd)
}

// mockVaultRepo implements vault.RepositorioCofre for tests.
type mockVaultRepo struct {
	salvarErr error
}

func (r *mockVaultRepo) Salvar(cofre *vault.Cofre) error { return r.salvarErr }
func (r *mockVaultRepo) Carregar() (*vault.Cofre, error) { return nil, errors.New("not implemented") }

// ─────────────────────────────────────────────────────────────────────────────
// saveAndExitFlow unit tests
// ─────────────────────────────────────────────────────────────────────────────

// TestSaveAndExitFlow_NoExtMod_SaveSuccess verifies the happy path:
// no external modification → save → saveAndExitOKMsg → tea.Quit.
func TestSaveAndExitFlow_NoExtMod_SaveSuccess(t *testing.T) {
	tmp, err := os.CreateTemp("", "abditum-test-*.abditum")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	tmp.Close()
	path := tmp.Name()
	defer os.Remove(path)

	meta, _ := storage.ComputeFileMetadata(path)
	saver := &mockVaultSaver{}
	flow := newSaveAndExitFlowWithMock(saver, path, meta, NewMessageManager())

	initCmd := flow.Init()
	if initCmd == nil {
		t.Fatal("Init should return a background I/O cmd")
	}
	initMsg := initCmd()

	if _, ok := initMsg.(saveAndExitReadyMsg); !ok {
		t.Fatalf("expected saveAndExitReadyMsg (no ext mod), got %T", initMsg)
	}

	saveCmd := flow.Update(initMsg)
	if saveCmd == nil {
		t.Fatal("Update(saveAndExitReadyMsg) should return a save cmd")
	}
	saveMsg := saveCmd()

	if !saver.salvarCalled {
		t.Error("expected Salvar() to be called")
	}
	if _, ok := saveMsg.(saveAndExitOKMsg); !ok {
		t.Fatalf("expected saveAndExitOKMsg after successful save, got %T", saveMsg)
	}

	quitCmd := flow.Update(saveMsg)
	if quitCmd == nil {
		t.Fatal("Update(saveAndExitOKMsg) should return tea.Quit cmd")
	}
	quitMsg := quitCmd()
	if _, ok := quitMsg.(tea.QuitMsg); !ok {
		t.Errorf("expected tea.QuitMsg, got %T", quitMsg)
	}
}

// TestSaveAndExitFlow_ExtModDetected verifies that when external modification
// is detected, a conflict Decision dialog is pushed.
func TestSaveAndExitFlow_ExtModDetected(t *testing.T) {
	tmp, err := os.CreateTemp("", "abditum-extmod-*.abditum")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	tmp.Close()
	path := tmp.Name()
	defer os.Remove(path)

	// Snapshot metadata before external modification
	meta, _ := storage.ComputeFileMetadata(path)

	// Externally modify the file so DetectExternalChange returns true
	if err := os.WriteFile(path, []byte("external modification"), 0600); err != nil {
		t.Fatalf("write external modification: %v", err)
	}

	saver := &mockVaultSaver{}
	flow := newSaveAndExitFlowWithMock(saver, path, meta, NewMessageManager())

	initCmd := flow.Init()
	if initCmd == nil {
		t.Fatal("Init should return a background I/O cmd")
	}
	initMsg := initCmd()

	if _, ok := initMsg.(extModDetectedMsg); !ok {
		t.Fatalf("expected extModDetectedMsg, got %T", initMsg)
	}

	dialogCmd := flow.Update(initMsg)
	if dialogCmd == nil {
		t.Fatal("Update(extModDetectedMsg) should return a dialog cmd")
	}
	// Decision() returns a tea.Cmd that emits pushModalMsg
	dialogMsg := dialogCmd()
	innerCmd, ok := dialogMsg.(tea.Cmd)
	if !ok {
		t.Fatalf("expected Decision() factory (tea.Cmd), got %T", dialogMsg)
	}
	if _, ok := innerCmd().(pushModalMsg); !ok {
		t.Error("expected pushModalMsg (conflict dialog)")
	}
}

// TestSaveAndExitFlow_ExtModOverwrite verifies that extModOverwriteMsg triggers save.
func TestSaveAndExitFlow_ExtModOverwrite(t *testing.T) {
	saver := &mockVaultSaver{}
	flow := newSaveAndExitFlowWithMock(saver, "/irrelevant", storage.FileMetadata{}, NewMessageManager())

	saveCmd := flow.Update(extModOverwriteMsg{})
	if saveCmd == nil {
		t.Fatal("Update(extModOverwriteMsg) should return a save cmd")
	}
	saveMsg := saveCmd()

	if !saver.salvarCalled {
		t.Error("expected Salvar() to be called after overwrite confirm")
	}
	if _, ok := saveMsg.(saveAndExitOKMsg); !ok {
		t.Fatalf("expected saveAndExitOKMsg, got %T", saveMsg)
	}
}

// TestSaveAndExitFlow_ExtModCancel verifies that extModCancelMsg ends the flow.
func TestSaveAndExitFlow_ExtModCancel(t *testing.T) {
	saver := &mockVaultSaver{}
	flow := newSaveAndExitFlowWithMock(saver, "/irrelevant", storage.FileMetadata{}, NewMessageManager())

	cancelCmd := flow.Update(extModCancelMsg{})
	if cancelCmd == nil {
		t.Fatal("Update(extModCancelMsg) should return endFlow cmd")
	}
	msg := cancelCmd()
	if _, ok := msg.(endFlowMsg); !ok {
		t.Errorf("expected endFlowMsg after cancel, got %T", msg)
	}
}

// TestSaveAndExitFlow_SaveError verifies that a save error ends the flow without quitting.
func TestSaveAndExitFlow_SaveError(t *testing.T) {
	saver := &mockVaultSaver{salvarErr: errors.New("disk full")}
	messages := NewMessageManager()
	flow := newSaveAndExitFlowWithMock(saver, "/irrelevant", storage.FileMetadata{}, messages)

	saveCmd := flow.Update(saveAndExitReadyMsg{})
	if saveCmd == nil {
		t.Fatal("Update(saveAndExitReadyMsg) should return a save cmd")
	}
	saveMsg := saveCmd()

	if _, ok := saveMsg.(endFlowMsg); !ok {
		t.Errorf("expected endFlowMsg after save error, got %T", saveMsg)
	}
	if saver.salvarCalled == false {
		t.Error("expected Salvar() to have been called")
	}
}
