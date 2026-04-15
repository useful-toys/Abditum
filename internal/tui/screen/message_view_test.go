package screen

import (
	"testing"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
)

func TestMessageLineView_ZeroValue(t *testing.T) {
	var v MessageLineView
	// Zero value deve renderizar sem pânico e retornar linha de borda com largura correta.
	output := v.Render(80, design.TokyoNight)
	if lipgloss.Width(output) == 0 {
		t.Error("Render de zero value retornou string vazia, esperado linha de borda")
	}
}

func TestMessageLineView_SetSuccess(t *testing.T) {
	var v MessageLineView
	v.SetSuccess("Cofre salvo")

	if v.current.Kind != design.MsgSuccess {
		t.Errorf("SetSuccess: Kind = %d, want %d", v.current.Kind, design.MsgSuccess)
	}
	if v.current.Text != "Cofre salvo" {
		t.Errorf("SetSuccess: Text = %q, want %q", v.current.Text, "Cofre salvo")
	}
	if v.ttl != design.MsgSuccess.DefaultTTL() {
		t.Errorf("SetSuccess: ttl = %d, want %d", v.ttl, design.MsgSuccess.DefaultTTL())
	}
}

func TestMessageLineView_SetBusy_ResetsSpinner(t *testing.T) {
	var v MessageLineView
	v.current.SpinnerFrame = 3 // simular que já havia animação em curso
	v.SetBusy("Salvando...")

	if v.current.SpinnerFrame != 0 {
		t.Errorf("SetBusy deve zerar SpinnerFrame, got %d", v.current.SpinnerFrame)
	}
	if v.current.Kind != design.MsgBusy {
		t.Errorf("SetBusy: Kind = %d, want %d", v.current.Kind, design.MsgBusy)
	}
	if v.ttl != 0 {
		t.Errorf("SetBusy: ttl = %d, want 0 (permanente)", v.ttl)
	}
}

func TestMessageLineView_Clear(t *testing.T) {
	var v MessageLineView
	v.SetSuccess("algo")
	v.Clear()

	var zero design.Message
	if v.current != zero {
		t.Errorf("Clear: current = %+v, want zero value", v.current)
	}
	if v.ttl != 0 {
		t.Errorf("Clear: ttl = %d, want 0", v.ttl)
	}
}

func TestMessageLineView_SetWarning(t *testing.T) {
	var v MessageLineView
	v.SetWarning("atenção")
	if v.current.Kind != design.MsgWarning {
		t.Errorf("SetWarning: Kind = %d, want %d", v.current.Kind, design.MsgWarning)
	}
	if v.ttl != 5 {
		t.Errorf("SetWarning: ttl = %d, want 5", v.ttl)
	}
}

func TestMessageLineView_SetError(t *testing.T) {
	var v MessageLineView
	v.SetError("falha")
	if v.current.Kind != design.MsgError {
		t.Errorf("SetError: Kind = %d, want %d", v.current.Kind, design.MsgError)
	}
	if v.ttl != 5 {
		t.Errorf("SetError: ttl = %d, want 5", v.ttl)
	}
}

func TestMessageLineView_SetInfo(t *testing.T) {
	var v MessageLineView
	v.SetInfo("info")
	if v.current.Kind != design.MsgInfo {
		t.Errorf("SetInfo: Kind = %d, want %d", v.current.Kind, design.MsgInfo)
	}
}

func TestMessageLineView_SetHintField(t *testing.T) {
	var v MessageLineView
	v.SetHintField("pressione Tab")
	if v.current.Kind != design.MsgHintField {
		t.Errorf("SetHintField: Kind = %d, want %d", v.current.Kind, design.MsgHintField)
	}
	if v.ttl != 0 {
		t.Errorf("SetHintField: ttl = %d, want 0 (permanente)", v.ttl)
	}
}

func TestMessageLineView_SetHintUsage(t *testing.T) {
	var v MessageLineView
	v.SetHintUsage("use ctrl+s para salvar")
	if v.current.Kind != design.MsgHintUsage {
		t.Errorf("SetHintUsage: Kind = %d, want %d", v.current.Kind, design.MsgHintUsage)
	}
	if v.ttl != 0 {
		t.Errorf("SetHintUsage: ttl = %d, want 0 (permanente)", v.ttl)
	}
}

func TestMessageLineView_Update_UnknownMsg(t *testing.T) {
	var v MessageLineView
	v.SetSuccess("algo")
	initialTTL := v.ttl

	type unknownMsg struct{}
	cmd := v.Update(unknownMsg{})

	if cmd != nil {
		t.Error("Update com msg desconhecida deve retornar nil cmd")
	}
	if v.ttl != initialTTL {
		t.Errorf("Update com msg desconhecida não deve alterar ttl: got %d, want %d", v.ttl, initialTTL)
	}
}

func TestMessageLineView_Render_WithMessage(t *testing.T) {
	var v MessageLineView
	v.SetSuccess("Cofre salvo")
	output := v.Render(80, design.TokyoNight)

	w := lipgloss.Width(output)
	if w != 80 {
		t.Errorf("Render: largura = %d, want 80", w)
	}
}

func TestMessageLineView_Render_ZeroValue_Width(t *testing.T) {
	var v MessageLineView
	output := v.Render(80, design.TokyoNight)

	w := lipgloss.Width(output)
	if w != 80 {
		t.Errorf("Render zero value: largura = %d, want 80", w)
	}
}

func TestMessageLineView_Render_ReturnsNoNewline(t *testing.T) {
	var v MessageLineView
	v.SetInfo("teste")
	output := v.Render(80, design.TokyoNight)

	for _, r := range output {
		if r == '\n' {
			t.Error("Render não deve conter newline — barra é linha única")
			break
		}
	}
}

// _testTick é um helper de teste que dispara a lógica de TickMsg diretamente,
// sem precisar importar package tui (o que causaria import cycle).
func (v *MessageLineView) _testTick() tea.Cmd {
	return v.tick()
}

func TestMessageLineView_SpinnerAdvances(t *testing.T) {
	var v MessageLineView
	v.SetBusy("carregando")

	for i := 1; i <= 8; i++ {
		v._testTick()
		want := i % 4
		if v.current.SpinnerFrame != want {
			t.Errorf("após %d ticks: SpinnerFrame = %d, want %d", i, v.current.SpinnerFrame, want)
		}
	}
}

func TestMessageLineView_TTL_Decrements(t *testing.T) {
	var v MessageLineView
	v.SetSuccess("ok") // ttl = 5

	for i := 4; i >= 1; i-- {
		v._testTick()
		if v.ttl != i {
			t.Errorf("após tick: ttl = %d, want %d", v.ttl, i)
		}
	}
	// Último tick: ttl chega a 0, mensagem é zerada.
	v._testTick()
	var zero design.Message
	if v.current != zero {
		t.Errorf("após ttl=0: current = %+v, want zero value", v.current)
	}
	if v.ttl != 0 {
		t.Errorf("após ttl=0: ttl = %d, want 0", v.ttl)
	}
}

func TestMessageLineView_BusyTTL_NeverExpires(t *testing.T) {
	var v MessageLineView
	v.SetBusy("operando")

	for i := 0; i < 10; i++ {
		v._testTick()
	}
	// Kind ainda deve ser MsgBusy — ttl=0 significa permanente.
	if v.current.Kind != design.MsgBusy {
		t.Errorf("MsgBusy não deve expirar: Kind = %d, want %d", v.current.Kind, design.MsgBusy)
	}
}

func TestMessageLineView_HintField_NeverExpires(t *testing.T) {
	var v MessageLineView
	v.SetHintField("pressione Enter")

	for i := 0; i < 10; i++ {
		v._testTick()
	}
	if v.current.Kind != design.MsgHintField {
		t.Errorf("MsgHintField não deve expirar: Kind = %d, want %d", v.current.Kind, design.MsgHintField)
	}
}
