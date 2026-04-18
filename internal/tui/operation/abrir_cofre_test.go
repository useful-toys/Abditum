package operation

import (
	"errors"
	"testing"

	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/tui"
)

// --- Init ---

func TestAbrirCofre_Init_SemCofre_AbreFilePicker(t *testing.T) {
	op := newAbrirCofreOperationFromSaver(&stubNotifier{}, nil, "")
	cmd := op.Init()
	msg := execCmd(cmd)
	switch v := msg.(type) {
	case tui.OpenModalMsg:
		// OK
	case abrirAvancaMsg:
		cmd2 := op.Update(v)
		msg2 := execCmd(cmd2)
		if _, ok := msg2.(tui.OpenModalMsg); !ok {
			t.Errorf("Init sem cofre: esperado OpenModalMsg (FilePicker), obteve %T", msg2)
		}
	default:
		t.Errorf("Init sem cofre: esperado OpenModalMsg ou abrirAvancaMsg, obteve %T", msg)
	}
}

func TestAbrirCofre_Init_CofreAlterado_AbreGuardModal(t *testing.T) {
	op := newAbrirCofreOperationFromSaver(&stubNotifier{}, &stubManager{isModified: true}, "")
	cmd := op.Init()
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init cofre alterado: esperado OpenModalMsg (guard), obteve %T", msg)
	}
}

func TestAbrirCofre_Init_ComCaminhoInicial_HeaderInvalido_Erro(t *testing.T) {
	n := &stubNotifier{}
	// Caminho inexistente — ValidateHeader retornará erro de IO
	op := newAbrirCofreOperationFromSaver(n, nil, "/caminho/inexistente.abditum")
	cmd := op.Init()
	msg := execCmd(cmd)
	if n.lastMethod != "SetError" {
		t.Errorf("header inválido: esperado SetError, obteve %q", n.lastMethod)
	}
	if _, ok := msg.(tui.OperationCompletedMsg); !ok {
		t.Errorf("header inválido: esperado OperationCompletedMsg, obteve %T", msg)
	}
}

// --- Update: abertura ---

func TestAbrirCofre_Update_AbrindoEstado_SetsBusy(t *testing.T) {
	n := &stubNotifier{}
	op := newAbrirCofreOperationFromSaver(n, nil, "")
	op.caminho = "/tmp/cofre.abditum"
	op.senha = []byte("senha")

	cmd := op.Update(abrirAvancaMsg{estado: abrindoAbrindo})
	if n.lastMethod != "SetBusy" {
		t.Errorf("abrindoAbrindo: esperado SetBusy, obteve %q", n.lastMethod)
	}
	if cmd == nil {
		t.Error("abrindoAbrindo: esperado cmd não-nil")
	}
}

func TestAbrirCofre_Update_ResultMsg_SenhaErrada_VoltaParaSenha(t *testing.T) {
	n := &stubNotifier{}
	op := newAbrirCofreOperationFromSaver(n, nil, "")

	cmd := op.Update(abrirCofreResultMsg{err: crypto.ErrAuthFailed})
	if n.lastMethod != "SetError" {
		t.Errorf("senha errada: esperado SetError, obteve %q", n.lastMethod)
	}
	// Deve reabrir o modal de senha
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("senha errada: esperado OpenModalMsg (senha), obteve %T", msg)
	}
}

func TestAbrirCofre_Update_ResultMsg_Corrompido_VoltaParaCaminho(t *testing.T) {
	n := &stubNotifier{}
	op := newAbrirCofreOperationFromSaver(n, nil, "")

	cmd := op.Update(abrirCofreResultMsg{err: storage.ErrCorrupted})
	if n.lastMethod != "SetError" {
		t.Errorf("corrompido: esperado SetError, obteve %q", n.lastMethod)
	}
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("corrompido: esperado OpenModalMsg (picker), obteve %T", msg)
	}
}

func TestAbrirCofre_Update_ResultMsg_Sucesso_EmiteVaultOpened(t *testing.T) {
	n := &stubNotifier{}
	op := newAbrirCofreOperationFromSaver(n, nil, "")

	cmd := op.Update(abrirCofreResultMsg{err: nil})
	if n.lastMethod != "Clear" {
		t.Errorf("sucesso: esperado Clear, obteve %q", n.lastMethod)
	}
	if cmd == nil {
		t.Error("sucesso: esperado cmd não-nil")
	}
}

func TestAbrirCofre_Update_MensagemDesconhecida_RetornaNil(t *testing.T) {
	op := newAbrirCofreOperationFromSaver(&stubNotifier{}, nil, "")
	type outraMsg struct{}
	if cmd := op.Update(outraMsg{}); cmd != nil {
		t.Error("Update(outraMsg): esperado nil")
	}
}

// --- Error classification ---

func TestErroDeAberturaCategoria_AuthFailed(t *testing.T) {
	msg := erroDeAberturaCategoria(crypto.ErrAuthFailed)
	if msg == "" {
		t.Error("ErrAuthFailed: esperado mensagem não-vazia")
	}
}

func TestErroDeAberturaCategoria_Corrompido(t *testing.T) {
	msg := erroDeAberturaCategoria(storage.ErrCorrupted)
	if msg == "" {
		t.Error("ErrCorrupted: esperado mensagem não-vazia")
	}
}

func TestErroDeAberturaCategoria_ErroGenerico(t *testing.T) {
	msg := erroDeAberturaCategoria(errors.New("algum erro"))
	if msg == "" {
		t.Error("erro genérico: esperado mensagem não-vazia")
	}
}
