package operation

import (
	"errors"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/vault"
)

// --- Init ---

func TestGuard_Init_SemCofre_ChamaOnProceder(t *testing.T) {
	var chamado bool
	g := novoGuardCofreAlterado(
		&stubNotifier{}, nil,
		func() tea.Cmd { chamado = true; return nil },
		func() tea.Cmd { t.Error("onAbortado não deveria ser chamado"); return nil },
	)
	execCmd(g.Init())
	if !chamado {
		t.Error("Init sem cofre: onProceder não foi chamado")
	}
}

func TestGuard_Init_CofreInalterado_ChamaOnProceder(t *testing.T) {
	var chamado bool
	g := novoGuardCofreAlterado(
		&stubNotifier{},
		&stubManager{isModified: false},
		func() tea.Cmd { chamado = true; return nil },
		func() tea.Cmd { t.Error("onAbortado não deveria ser chamado"); return nil },
	)
	execCmd(g.Init())
	if !chamado {
		t.Error("Init cofre inalterado: onProceder não foi chamado")
	}
}

func TestGuard_Init_CofreAlterado_AbreModal(t *testing.T) {
	g := novoGuardCofreAlterado(
		&stubNotifier{},
		&stubManager{isModified: true},
		func() tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
	msg := execCmd(g.Init())
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init cofre alterado: esperado OpenModalMsg, obteve %T", msg)
	}
}

// --- Update: guardSaveMsg ---

func TestGuard_Update_SalvarSucesso_ChamaOnProceder(t *testing.T) {
	n := &stubNotifier{}
	var chamado bool
	g := novoGuardCofreAlterado(
		n,
		&stubManager{isModified: true},
		func() tea.Cmd { chamado = true; return nil },
		func() tea.Cmd { t.Error("onAbortado inesperado"); return nil },
	)
	cmd := g.Update(guardSaveMsg{forced: false})
	if n.lastMethod != "SetBusy" {
		t.Errorf("esperado SetBusy, obteve %q", n.lastMethod)
	}
	resultMsg := execCmd(cmd)
	execCmd(g.Update(resultMsg))
	if !chamado {
		t.Error("após salvar OK: onProceder não foi chamado")
	}
	if n.lastMethod != "Clear" {
		t.Errorf("após salvar OK: esperado Clear, obteve %q", n.lastMethod)
	}
}

func TestGuard_Update_SalvarErroGenerico_ChamaOnAbortado(t *testing.T) {
	n := &stubNotifier{}
	var abortado bool
	g := novoGuardCofreAlterado(
		n,
		&stubManager{isModified: true, salvarErr: errors.New("disco cheio")},
		func() tea.Cmd { t.Error("onProceder inesperado"); return nil },
		func() tea.Cmd { abortado = true; return nil },
	)
	cmd := g.Update(guardSaveMsg{forced: false})
	resultMsg := execCmd(cmd)
	execCmd(g.Update(resultMsg))
	if !abortado {
		t.Error("após erro genérico: onAbortado não foi chamado")
	}
	if n.lastMethod != "SetError" {
		t.Errorf("após erro genérico: esperado SetError, obteve %q", n.lastMethod)
	}
}

func TestGuard_Update_ModificadoExternamente_AbreModalConflito(t *testing.T) {
	n := &stubNotifier{}
	g := novoGuardCofreAlterado(
		n,
		&stubManager{isModified: true, salvarErr: vault.ErrModifiedExternally},
		func() tea.Cmd { return nil },
		func() tea.Cmd { return nil },
	)
	cmd := g.Update(guardSaveMsg{forced: false})
	resultMsg := execCmd(cmd)
	resultCmd := g.Update(resultMsg)
	msg := execCmd(resultCmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("ErrModifiedExternally: esperado OpenModalMsg (conflito), obteve %T", msg)
	}
	if n.lastMethod != "Clear" {
		t.Errorf("ErrModifiedExternally: esperado Clear, obteve %q", n.lastMethod)
	}
}

func TestGuard_Update_SalvarForcado_Sucesso_ChamaOnProceder(t *testing.T) {
	n := &stubNotifier{}
	var chamado bool
	g := novoGuardCofreAlterado(
		n,
		&stubManager{isModified: true},
		func() tea.Cmd { chamado = true; return nil },
		func() tea.Cmd { t.Error("onAbortado inesperado"); return nil },
	)
	cmd := g.Update(guardSaveMsg{forced: true})
	resultMsg := execCmd(cmd)
	execCmd(g.Update(resultMsg))
	if !chamado {
		t.Error("após salvar forçado OK: onProceder não foi chamado")
	}
}

func TestGuard_Update_SalvarForcado_Erro_ChamaOnAbortado(t *testing.T) {
	n := &stubNotifier{}
	var abortado bool
	m := &stubManager{isModified: true, salvarErr: errors.New("falha")}
	g := novoGuardCofreAlterado(
		n, m,
		func() tea.Cmd { t.Error("onProceder inesperado"); return nil },
		func() tea.Cmd { abortado = true; return nil },
	)
	cmd := g.Update(guardSaveMsg{forced: true})
	resultMsg := execCmd(cmd)
	execCmd(g.Update(resultMsg))
	if !abortado {
		t.Error("após forçado com erro: onAbortado não foi chamado")
	}
}

// --- Update: descarte ---

func TestGuard_Update_DescartarMsg_ChamaOnProceder(t *testing.T) {
	var chamado bool
	g := novoGuardCofreAlterado(
		&stubNotifier{},
		&stubManager{isModified: true},
		func() tea.Cmd { chamado = true; return nil },
		func() tea.Cmd { t.Error("onAbortado inesperado"); return nil },
	)
	execCmd(g.Update(guardDiscardMsg{}))
	if !chamado {
		t.Error("descartar: onProceder não foi chamado")
	}
}

// --- Update: cancelar ---

func TestGuard_Update_CancelarMsg_ChamaOnAbortado(t *testing.T) {
	var abortado bool
	g := novoGuardCofreAlterado(
		&stubNotifier{},
		&stubManager{isModified: true},
		func() tea.Cmd { t.Error("onProceder inesperado"); return nil },
		func() tea.Cmd { abortado = true; return nil },
	)
	execCmd(g.Update(guardCancelMsg{}))
	if !abortado {
		t.Error("cancelar: onAbortado não foi chamado")
	}
}
