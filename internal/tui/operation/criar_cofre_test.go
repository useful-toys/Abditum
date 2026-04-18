package operation

import (
	"errors"
	"testing"

	"github.com/useful-toys/abditum/internal/tui"
)

// --- Init ---

func TestCriarCofre_Init_SemCofre_AbreFilePicker(t *testing.T) {
	op := newCriarCofreOperationFromSaver(&stubNotifier{}, nil, "")
	cmd := op.Init()
	msg := execCmd(cmd)
	// guard com saver nil → onProceder() imediato → emite criarAvancaMsg
	// que no Update abre FilePicker → OpenModalMsg
	switch v := msg.(type) {
	case tui.OpenModalMsg:
		// OK: diretamente abriu modal
	case criarAvancaMsg:
		// onProceder emitiu criarAvancaMsg — processar no Update
		cmd2 := op.Update(v)
		msg2 := execCmd(cmd2)
		if _, ok := msg2.(tui.OpenModalMsg); !ok {
			t.Errorf("Init sem cofre: esperado OpenModalMsg (FilePicker), obteve %T", msg2)
		}
	default:
		t.Errorf("Init sem cofre: esperado OpenModalMsg ou criarAvancaMsg, obteve %T", msg)
	}
}

func TestCriarCofre_Init_ComCaminhoInicial_AbrePasswordModal(t *testing.T) {
	op := newCriarCofreOperationFromSaver(&stubNotifier{}, nil, "/tmp/cofre.abditum")
	cmd := op.Init()
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init com caminho: esperado OpenModalMsg (PasswordCreate), obteve %T", msg)
	}
}

func TestCriarCofre_Init_CofreAlterado_AbreGuardModal(t *testing.T) {
	op := newCriarCofreOperationFromSaver(&stubNotifier{}, &stubManager{isModified: true}, "")
	cmd := op.Init()
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OpenModalMsg); !ok {
		t.Errorf("Init cofre alterado: esperado OpenModalMsg (guard), obteve %T", msg)
	}
}

// --- Update: criação ---

func TestCriarCofre_Update_CriandoEstado_SetsBusy(t *testing.T) {
	n := &stubNotifier{}
	op := newCriarCofreOperationFromSaver(n, nil, "/tmp/cofre.abditum")
	op.caminho = "/tmp/cofre.abditum"
	op.senha = []byte("SenhaForte123!")

	cmd := op.Update(criarAvancaMsg{estado: criandoCriando})
	if n.lastMethod != "SetBusy" {
		t.Errorf("criandoCriando: esperado SetBusy, obteve %q", n.lastMethod)
	}
	if cmd == nil {
		t.Error("criandoCriando: esperado cmd não-nil")
	}
}

func TestCriarCofre_Update_ResultMsg_Falha_EmiteSetErrorECompleta(t *testing.T) {
	n := &stubNotifier{}
	op := newCriarCofreOperationFromSaver(n, nil, "")

	cmd := op.Update(criarCofreResultMsg{err: errors.New("falha ao criar")})
	if n.lastMethod != "SetError" {
		t.Errorf("falha: esperado SetError, obteve %q", n.lastMethod)
	}
	msg := execCmd(cmd)
	if _, ok := msg.(tui.OperationCompletedMsg); !ok {
		t.Errorf("falha: esperado OperationCompletedMsg, obteve %T", msg)
	}
}

func TestCriarCofre_Update_ResultMsg_Sucesso_EmiteVaultOpened(t *testing.T) {
	n := &stubNotifier{}
	op := newCriarCofreOperationFromSaver(n, nil, "")

	cmd := op.Update(criarCofreResultMsg{err: nil})
	if n.lastMethod != "Clear" {
		t.Errorf("sucesso: esperado Clear, obteve %q", n.lastMethod)
	}
	if cmd == nil {
		t.Error("sucesso: esperado cmd não-nil")
	}
}

func TestCriarCofre_Update_MensagemDesconhecida_RetornaNil(t *testing.T) {
	op := newCriarCofreOperationFromSaver(&stubNotifier{}, nil, "")
	type outraMsg struct{}
	if cmd := op.Update(outraMsg{}); cmd != nil {
		t.Error("Update(outraMsg): esperado nil")
	}
}
