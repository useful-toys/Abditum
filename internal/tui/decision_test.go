package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	testdatapkg "github.com/useful-toys/abditum/internal/tui/testdata"
)

// ─────────────────────────────────────────────────────────────────────────────
// PoC fixture constructors (15 combinations: 5 severities × 3 action counts)
// These match the PoC table in root.go (keys 1–9, a–f).
// Each calls SetAvailableSize(80, 24) so boxWidth() returns a concrete value.
// ─────────────────────────────────────────────────────────────────────────────

// ── Destrutivo ───────────────────────────────────────────────────────────────

func pocKey1() *DecisionDialog {
	d := NewDecisionDialog(SeverityDestructive, IntentionAcknowledge,
		"Exclusão concluída", "Gmail foi excluído permanentemente.",
		[]DecisionAction{
			{Key: "Enter", Label: "OK", Default: true},
		})
	return d
}

func pocKey2() *DecisionDialog {
	d := NewDecisionDialog(SeverityDestructive, IntentionConfirm,
		"Excluir segredo",
		"Gmail será excluído permanentemente. Esta ação não pode ser desfeita.",
		[]DecisionAction{
			{Key: "Enter", Label: "Excluir", Default: true},
			{Key: "Esc", Label: "Cancelar", Cancel: true},
		})
	return d
}

func pocKey3() *DecisionDialog {
	d := NewDecisionDialog(SeverityDestructive, IntentionConfirm,
		"Excluir pasta",
		"Financeiro e todos os seus segredos serão excluídos permanentemente.",
		[]DecisionAction{
			{Key: "Enter", Label: "Excluir", Default: true},
			{Key: "M", Label: "Mover conteúdo"},
			{Key: "Esc", Label: "Cancelar", Cancel: true},
		})
	return d
}

// ── Erro ─────────────────────────────────────────────────────────────────────

func pocKey4() *DecisionDialog {
	d := NewDecisionDialog(SeverityError, IntentionAcknowledge,
		"Falha ao salvar",
		"Não foi possível salvar o cofre. O arquivo pode estar em uso por outro processo.",
		[]DecisionAction{
			{Key: "Enter", Label: "OK", Default: true},
		})
	return d
}

func pocKey5() *DecisionDialog {
	d := NewDecisionDialog(SeverityError, IntentionConfirm,
		"Senha incorreta",
		"A senha está incorreta. O cofre não pôde ser aberto.",
		[]DecisionAction{
			{Key: "Enter", Label: "Tentar novamente", Default: true},
			{Key: "Esc", Label: "Cancelar", Cancel: true},
		})
	return d
}

func pocKey6() *DecisionDialog {
	d := NewDecisionDialog(SeverityError, IntentionConfirm,
		"Cofre corrompido",
		"O arquivo está corrompido. Deseja tentar recuperar a partir do backup?",
		[]DecisionAction{
			{Key: "Enter", Label: "Recuperar", Default: true},
			{Key: "A", Label: "Abrir backup"},
			{Key: "Esc", Label: "Cancelar", Cancel: true},
		})
	return d
}

// ── Alerta ───────────────────────────────────────────────────────────────────

func pocKey7() *DecisionDialog {
	d := NewDecisionDialog(SeverityAlert, IntentionAcknowledge,
		"Sessão bloqueada",
		"O cofre foi bloqueado após 5 minutos de inatividade.",
		[]DecisionAction{
			{Key: "Enter", Label: "OK", Default: true},
		})
	return d
}

func pocKey8() *DecisionDialog {
	d := NewDecisionDialog(SeverityAlert, IntentionConfirm,
		"Alterações não salvas",
		"Existem alterações não salvas. Sair irá descartá-las.",
		[]DecisionAction{
			{Key: "Enter", Label: "Descartar", Default: true},
			{Key: "Esc", Label: "Voltar", Cancel: true},
		})
	return d
}

func pocKey9() *DecisionDialog {
	d := NewDecisionDialog(SeverityAlert, IntentionConfirm,
		"Senha fraca",
		"A senha mestra é fraca e pode ser facilmente descoberta.",
		[]DecisionAction{
			{Key: "Enter", Label: "Usar assim mesmo", Default: true},
			{Key: "T", Label: "Trocar senha"},
			{Key: "Esc", Label: "Cancelar", Cancel: true},
		})
	return d
}

// ── Informativo ───────────────────────────────────────────────────────────────

func pocKeyA() *DecisionDialog {
	d := NewDecisionDialog(SeverityInformative, IntentionAcknowledge,
		"Cofre criado",
		"O cofre foi criado com sucesso em ~/documentos/pessoal.abditum.",
		[]DecisionAction{
			{Key: "Enter", Label: "OK", Default: true},
		})
	return d
}

func pocKeyB() *DecisionDialog {
	d := NewDecisionDialog(SeverityInformative, IntentionConfirm,
		"Conflito detectado",
		"O arquivo foi modificado externamente desde a última abertura.",
		[]DecisionAction{
			{Key: "Enter", Label: "Sobrescrever", Default: true},
			{Key: "Esc", Label: "Cancelar", Cancel: true},
		})
	return d
}

func pocKeyC() *DecisionDialog {
	d := NewDecisionDialog(SeverityInformative, IntentionConfirm,
		"Importação concluída",
		"12 segredos importados. 3 entradas já existentes foram atualizadas.",
		[]DecisionAction{
			{Key: "Enter", Label: "Ver detalhes", Default: true},
			{Key: "F", Label: "Fechar"},
			{Key: "Esc", Label: "OK", Cancel: true},
		})
	return d
}

// ── Neutro ────────────────────────────────────────────────────────────────────

func pocKeyD() *DecisionDialog {
	d := NewDecisionDialog(SeverityNeutral, IntentionAcknowledge,
		"Operação concluída",
		"A exportação foi salva em ~/documentos/backup-2026-04-05.json.",
		[]DecisionAction{
			{Key: "Enter", Label: "OK", Default: true},
		})
	return d
}

func pocKeyE() *DecisionDialog {
	d := NewDecisionDialog(SeverityNeutral, IntentionConfirm,
		"Sair do Abditum",
		"Deseja sair? Todas as alterações não salvas serão descartadas.",
		[]DecisionAction{
			{Key: "Enter", Label: "Sair", Default: true},
			{Key: "Esc", Label: "Cancelar", Cancel: true},
		})
	return d
}

func pocKeyF() *DecisionDialog {
	d := NewDecisionDialog(SeverityNeutral, IntentionConfirm,
		"Salvar cofre",
		"Deseja salvar as alterações antes de continuar?",
		[]DecisionAction{
			{Key: "Enter", Label: "Salvar", Default: true},
			{Key: "N", Label: "Não salvar"},
			{Key: "Esc", Label: "Voltar", Cancel: true},
		})
	return d
}

// ─────────────────────────────────────────────────────────────────────────────
// TestDecisionDialog_MatrixViewRenders — table-driven over all 15 fixtures
// ─────────────────────────────────────────────────────────────────────────────

func TestDecisionDialog_MatrixViewRenders(t *testing.T) {
	type fixture struct {
		name         string
		dialog       func() *DecisionDialog
		titleContain string
		actionLabel  string
	}

	fixtures := []fixture{
		{"Dest·1 Acknowledge", pocKey1, "Exclusão concluída", "OK"},
		{"Dest·2 Confirm", pocKey2, "Excluir segredo", "Excluir"},
		{"Dest·3 Confirm 3", pocKey3, "Excluir pasta", "Excluir"},
		{"Err·1 Acknowledge", pocKey4, "Falha ao salvar", "OK"},
		{"Err·2 Confirm", pocKey5, "Senha incorreta", "Tentar novamente"},
		{"Err·3 Confirm 3", pocKey6, "Cofre corrompido", "Recuperar"},
		{"Ale·1 Acknowledge", pocKey7, "Sessão bloqueada", "OK"},
		{"Ale·2 Confirm", pocKey8, "Alterações não salvas", "Descartar"},
		{"Ale·3 Confirm 3", pocKey9, "Senha fraca", "Usar assim mesmo"},
		{"Inf·1 Acknowledge", pocKeyA, "Cofre criado", "OK"},
		{"Inf·2 Confirm", pocKeyB, "Conflito detectado", "Sobrescrever"},
		{"Inf·3 Confirm 3", pocKeyC, "Importação concluída", "Ver detalhes"},
		{"Neu·1 Acknowledge", pocKeyD, "Operação concluída", "OK"},
		{"Neu·2 Confirm", pocKeyE, "Sair do Abditum", "Sair"},
		{"Neu·3 Confirm 3", pocKeyF, "Salvar cofre", "Salvar"},
	}

	for _, f := range fixtures {
		f := f // capture loop variable
		t.Run(f.name, func(t *testing.T) {
			d := f.dialog()
			out := d.View(80, 24, TokyoNight)

			if out == "" {
				t.Errorf("%s: View() returned empty string", f.name)
			}
			if !strings.Contains(out, f.titleContain) {
				t.Errorf("%s: View() missing title %q\ngot:\n%s", f.name, f.titleContain, out)
			}
			if !strings.Contains(out, f.actionLabel) {
				t.Errorf("%s: View() missing action label %q\ngot:\n%s", f.name, f.actionLabel, out)
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// TestDecisionDialog_SymbolPresence
// ─────────────────────────────────────────────────────────────────────────────

func TestDecisionDialog_SymbolPresence(t *testing.T) {
	tests := []struct {
		name       string
		d          *DecisionDialog
		wantSymbol string
		wantAbsent []string
	}{
		{
			name:       "Informativo contains ℹ",
			d:          pocKeyA(),
			wantSymbol: SymInfo,
		},
		{
			name:       "Alerta contains ⚠",
			d:          pocKey7(),
			wantSymbol: SymWarning,
		},
		{
			name:       "Destrutivo contains ⚠",
			d:          pocKey1(),
			wantSymbol: SymWarning,
		},
		{
			name:       "Erro contains ✕",
			d:          pocKey4(),
			wantSymbol: SymError,
		},
		{
			name:       "Neutro has no symbol",
			d:          pocKeyD(),
			wantSymbol: "",
			wantAbsent: []string{SymWarning, SymError, SymInfo},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.d.View(80, 24, TokyoNight)
			if tt.wantSymbol != "" && !strings.Contains(out, tt.wantSymbol) {
				t.Errorf("expected symbol %q in output\ngot:\n%s", tt.wantSymbol, out)
			}
			for _, absent := range tt.wantAbsent {
				if strings.Contains(out, absent) {
					t.Errorf("expected symbol %q to be ABSENT in Neutro output\ngot:\n%s", absent, out)
				}
			}
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// TestDecisionDialog_BorderChars
// ─────────────────────────────────────────────────────────────────────────────

func TestDecisionDialog_BorderChars(t *testing.T) {
	d := pocKey1()
	out := d.View(80, 24, TokyoNight)

	if !strings.Contains(out, "╭") {
		t.Errorf("View() missing top-left corner ╭\ngot:\n%s", out)
	}
	if !strings.Contains(out, "╰") {
		t.Errorf("View() missing bottom-left corner ╰\ngot:\n%s", out)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Interaction tests
// ─────────────────────────────────────────────────────────────────────────────

// TestDecisionDialog_EnterTriggersDefault: pocKeyF (Neutro 3-action: Salvar/Não salvar/Voltar),
// send KeyEnter, assert returned tea.Cmd is non-nil.
func TestDecisionDialog_EnterTriggersDefault(t *testing.T) {
	d := pocKeyF()
	cmd := d.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Error("Enter should return a non-nil tea.Cmd (popModal + optional user cmd)")
	}
}

// TestDecisionDialog_EscTriggersCancel: pocKey2 (Destrutivo 2-action: Excluir/Cancelar),
// send esc key, assert cmd is non-nil.
func TestDecisionDialog_EscTriggersCancel(t *testing.T) {
	d := pocKey2()
	cmd := d.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if cmd == nil {
		t.Error("Esc should return a non-nil tea.Cmd (popModal)")
	}
}

// TestDecisionDialog_UnknownKeyIgnored: pocKeyE (Neutro 2-action),
// send key "z", assert cmd is nil.
func TestDecisionDialog_UnknownKeyIgnored(t *testing.T) {
	d := pocKeyE()
	cmd := d.Update(tea.KeyPressMsg{Code: 'z'})
	if cmd != nil {
		t.Error("unknown key 'z' should return nil cmd")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Edge case tests
// ─────────────────────────────────────────────────────────────────────────────

// TestDecisionDialog_LongBodyWraps: body with 150+ chars, View(80,24),
// View() should contain \n (multi-line body).
func TestDecisionDialog_LongBodyWraps(t *testing.T) {
	longBody := "Esta é uma mensagem muito longa para testar o sistema de quebra de linha do corpo do diálogo que deve quebrar em múltiplas linhas quando excede a largura disponível da caixa do diálogo."
	if len([]rune(longBody)) < 150 {
		t.Fatal("test precondition: longBody must be 150+ chars")
	}
	d := NewDecisionDialog(SeverityNeutral, IntentionAcknowledge,
		"Título", longBody,
		[]DecisionAction{{Key: "Enter", Label: "OK", Default: true}})
	out := d.View(80, 24, TokyoNight)
	if !strings.Contains(out, "\n") {
		t.Error("View() with long body should contain newlines (multi-line wrapping)")
	}
}

// TestDecisionDialog_ShortBodyFits: single short body, View(50,10),
// View() is non-empty and does not panic.
func TestDecisionDialog_ShortBodyFits(t *testing.T) {
	d := NewDecisionDialog(SeverityNeutral, IntentionAcknowledge,
		"Título", "OK?",
		[]DecisionAction{{Key: "Enter", Label: "OK", Default: true}})
	out := d.View(50, 10, TokyoNight)
	if out == "" {
		t.Error("View() with short body should return non-empty string")
	}
}

// TestDecisionDialog_AcknowledgeHasNoCancel: pocKey1 (Destrutivo 1-action),
// output should NOT contain "Esc".
func TestDecisionDialog_AcknowledgeHasNoCancel(t *testing.T) {
	d := pocKey1()
	out := d.View(80, 24, TokyoNight)
	if strings.Contains(out, "Esc") {
		t.Errorf("Acknowledge dialog should not render 'Esc' in output\ngot:\n%s", out)
	}
}

// TestDecisionDialog_SmallSizeUsesMinWidth: pocKey4 with View(20,10) (below minimum),
// View() is non-empty and does not panic (falls back to boxWidth=40 floor).
func TestDecisionDialog_SmallSizeUsesMinWidth(t *testing.T) {
	d := pocKey4()
	out := d.View(20, 10, TokyoNight)
	if out == "" {
		t.Error("View() with small terminal size should return non-empty string (uses min width floor)")
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// Golden test helpers (decision-specific)
// ─────────────────────────────────────────────────────────────────────────────

// decisionGoldenPath returns the golden file path for a decision dialog scenario.
// The variant string already includes the width (e.g. "destructive-1action-short-30"),
// so we do NOT append a separate width field (unlike the generic goldenPath helper).
func decisionGoldenPath(variant, ext string) string {
	name := fmt.Sprintf("decision-%s.%s.golden", variant, ext)
	return filepath.Join("testdata", "golden", name)
}

// checkOrUpdateDecisionGolden compares output against the golden file for a decision scenario,
// or writes it if the -update flag is set or the file is missing (first run).
func checkOrUpdateDecisionGolden(t *testing.T, variant, ext, got string) {
	t.Helper()
	path := decisionGoldenPath(variant, ext)
	if *update {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("mkdirall %s: %v", filepath.Dir(path), err)
		}
		if err := os.WriteFile(path, []byte(got), 0644); err != nil {
			t.Fatalf("write golden %s: %v", path, err)
		}
		return
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		// First run: auto-generate golden file instead of failing.
		if err2 := os.MkdirAll(filepath.Dir(path), 0755); err2 != nil {
			t.Fatalf("mkdirall %s: %v", filepath.Dir(path), err2)
		}
		if err2 := os.WriteFile(path, []byte(got), 0644); err2 != nil {
			t.Fatalf("write golden %s: %v", path, err2)
		}
		return
	}
	if err != nil {
		t.Fatalf("read golden %s: %v", path, err)
	}
	if string(data) != got {
		t.Errorf("golden mismatch for %s:\nwant:\n%s\ngot:\n%s", path, string(data), got)
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// TestDecisionDialog_Golden — 10 scenarios × (txt + json) = 20 golden files
// ─────────────────────────────────────────────────────────────────────────────

func TestDecisionDialog_Golden(t *testing.T) {
	type testCase struct {
		variant string
		dialog  *DecisionDialog
	}

	// Helper: build a dialog constructor for a given width.
	// Each named scenario is rendered at both 30x24 and 60x24 to cover the full
	// 10-scenarios × 2-widths × 2-formats = 40-file matrix required by the spec.
	newDestructive1Short := func(w int) *DecisionDialog {
		d := NewDecisionDialog(SeverityDestructive, IntentionAcknowledge,
			"Excluir segredo",
			"Gmail será excluído permanentemente.",
			[]DecisionAction{{Key: "Enter", Label: "Excluir", Default: true}})
		return d
	}
	newDestructive2Long := func(w int) *DecisionDialog {
		d := NewDecisionDialog(SeverityDestructive, IntentionConfirm,
			"Excluir permanentemente este segredo do cofre atual?",
			"Esta ação não pode ser desfeita. Todos os dados associados serão removidos.",
			[]DecisionAction{
				{Key: "Enter", Label: "Excluir", Default: true},
				{Key: "Esc", Label: "Cancelar", Cancel: true},
			})
		return d
	}
	newError3Short := func(w int) *DecisionDialog {
		d := NewDecisionDialog(SeverityError, IntentionConfirm,
			"Cofre corrompido",
			"O arquivo está corrompido. Deseja tentar recuperar?",
			[]DecisionAction{
				{Key: "Enter", Label: "Recuperar", Default: true},
				{Key: "A", Label: "Abrir backup"},
				{Key: "Esc", Label: "Cancelar", Cancel: true},
			})
		return d
	}
	newError1Long := func(w int) *DecisionDialog {
		d := NewDecisionDialog(SeverityError, IntentionAcknowledge,
			"Erro crítico ao acessar o cofre — arquivo danificado",
			"Não foi possível decodificar o arquivo. Verifique se o disco está íntegro e tente novamente.",
			[]DecisionAction{{Key: "Enter", Label: "OK", Default: true}})
		return d
	}
	newAlert2Short := func(w int) *DecisionDialog {
		d := NewDecisionDialog(SeverityAlert, IntentionConfirm,
			"Sobrescrever?",
			"Já existe um segredo com este nome. Deseja substituir o existente?",
			[]DecisionAction{
				{Key: "Enter", Label: "Sim", Default: true},
				{Key: "Esc", Label: "Não", Cancel: true},
			})
		return d
	}
	newAlert3Long := func(w int) *DecisionDialog {
		d := NewDecisionDialog(SeverityAlert, IntentionConfirm,
			"Conflito de nome ao salvar novo segredo no cofre",
			"Um segredo chamado 'github-token' já existe.",
			[]DecisionAction{
				{Key: "Enter", Label: "Substituir", Default: true},
				{Key: "R", Label: "Renomear"},
				{Key: "Esc", Label: "Cancelar", Cancel: true},
			})
		return d
	}
	newInformative1Short := func(w int) *DecisionDialog {
		d := NewDecisionDialog(SeverityInformative, IntentionAcknowledge,
			"Dica",
			"Pressione Ctrl+N para criar um novo cofre.",
			[]DecisionAction{{Key: "Enter", Label: "Entendi", Default: true}})
		return d
	}
	newInformative2Long := func(w int) *DecisionDialog {
		d := NewDecisionDialog(SeverityInformative, IntentionConfirm,
			"Segredo copiado para a área de transferência com sucesso",
			"O conteúdo será limpo automaticamente em 30 segundos por segurança.",
			[]DecisionAction{
				{Key: "Enter", Label: "Copiar", Default: true},
				{Key: "Esc", Label: "Fechar", Cancel: true},
			})
		return d
	}
	newNeutral3Short := func(w int) *DecisionDialog {
		d := NewDecisionDialog(SeverityNeutral, IntentionConfirm,
			"Continuar?",
			"Tem certeza que deseja prosseguir?",
			[]DecisionAction{
				{Key: "Enter", Label: "Sim", Default: true},
				{Key: "N", Label: "Não"},
				{Key: "Esc", Label: "Talvez", Cancel: true},
			})
		return d
	}
	newNeutral2Long := func(w int) *DecisionDialog {
		d := NewDecisionDialog(SeverityNeutral, IntentionConfirm,
			"Confirmar alterações pendentes antes de fechar o cofre?",
			"Existem 3 modificações não salvas. Se fechar sem salvar, as alterações serão perdidas.",
			[]DecisionAction{
				{Key: "Enter", Label: "Confirmar", Default: true},
				{Key: "Esc", Label: "Cancelar", Cancel: true},
			})
		return d
	}

	cases := []testCase{
		// Scenario 1: Destructive 1-action short title — 2 widths
		{variant: "destructive-1action-short-30x24", dialog: newDestructive1Short(30)},
		{variant: "destructive-1action-short-60x24", dialog: newDestructive1Short(60)},
		// Scenario 2: Destructive 2-action long title+body — 2 widths
		{variant: "destructive-2action-long-30x24", dialog: newDestructive2Long(30)},
		{variant: "destructive-2action-long-60x24", dialog: newDestructive2Long(60)},
		// Scenario 3: Error 3-action short title — 2 widths
		{variant: "error-3action-short-30x24", dialog: newError3Short(30)},
		{variant: "error-3action-short-60x24", dialog: newError3Short(60)},
		// Scenario 4: Error 1-action long title+body — 2 widths
		{variant: "error-1action-long-30x24", dialog: newError1Long(30)},
		{variant: "error-1action-long-60x24", dialog: newError1Long(60)},
		// Scenario 5: Alert 2-action short title, 2-line body — 2 widths
		{variant: "alert-2action-short-30x24", dialog: newAlert2Short(30)},
		{variant: "alert-2action-short-60x24", dialog: newAlert2Short(60)},
		// Scenario 6: Alert 3-action long title, 1-line body — 2 widths
		{variant: "alert-3action-long-30x24", dialog: newAlert3Long(30)},
		{variant: "alert-3action-long-60x24", dialog: newAlert3Long(60)},
		// Scenario 7: Informative 1-action short title, 1-line body — 2 widths
		{variant: "informative-1action-short-30x24", dialog: newInformative1Short(30)},
		{variant: "informative-1action-short-60x24", dialog: newInformative1Short(60)},
		// Scenario 8: Informative 2-action long title, 2-line body — 2 widths
		{variant: "informative-2action-long-30x24", dialog: newInformative2Long(30)},
		{variant: "informative-2action-long-60x24", dialog: newInformative2Long(60)},
		// Scenario 9: Neutral 3-action short title, 1-line body — 2 widths
		{variant: "neutral-3action-short-30x24", dialog: newNeutral3Short(30)},
		{variant: "neutral-3action-short-60x24", dialog: newNeutral3Short(60)},
		// Scenario 10: Neutral 2-action long title, 2-line body — 2 widths
		{variant: "neutral-2action-long-30x24", dialog: newNeutral2Long(30)},
		{variant: "neutral-2action-long-60x24", dialog: newNeutral2Long(60)},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.variant, func(t *testing.T) {
			// Extract width and height from variant (e.g., "destructive-1action-short-30x24")
			var w, h int
			if strings.HasSuffix(tc.variant, "-30x24") {
				w, h = 30, 24
			} else if strings.HasSuffix(tc.variant, "-60x24") {
				w, h = 60, 24
			}
			out := tc.dialog.View(w, h)

			// .txt.golden: raw ANSI output
			checkOrUpdateDecisionGolden(t, tc.variant, "txt", stripANSI(out))

			// .json.golden: style transitions
			transitions := testdatapkg.ParseANSIStyle(out)
			jsonBytes, err := testdatapkg.MarshalStyleTransitions(transitions)
			if err != nil {
				t.Fatalf("marshal transitions: %v", err)
			}
			checkOrUpdateDecisionGolden(t, tc.variant, "json", string(jsonBytes))
		})
	}
}

// ─────────────────────────────────────────────────────────────────────────────
// TestDecisionDialog_Update_* — explicit key binding behavioral tests
// ─────────────────────────────────────────────────────────────────────────────

// TestDecisionDialog_Update_ExplicitKey_Destructive: pocKey3 has key "M" for "Mover conteúdo".
// Sending key "m" (lowercase) must return non-nil cmd (case-insensitive match per Update logic).
func TestDecisionDialog_Update_ExplicitKey_Destructive(t *testing.T) {
	d := pocKey3() // Destructive 3-action: Enter Excluir / M Mover conteúdo / Esc Cancelar
	cmd := d.Update(tea.KeyPressMsg{Code: 'm'})
	if cmd == nil {
		t.Error("explicit key 'm' (matching 'M' action) must return non-nil cmd")
	}
}

// TestDecisionDialog_Update_ExplicitKey_Error: pocKey6 has key "A" for "Abrir backup".
func TestDecisionDialog_Update_ExplicitKey_Error(t *testing.T) {
	d := pocKey6() // Error 3-action: Enter Recuperar / A Abrir backup / Esc Cancelar
	cmd := d.Update(tea.KeyPressMsg{Code: 'a'})
	if cmd == nil {
		t.Error("explicit key 'a' (matching 'A' action) must return non-nil cmd")
	}
}

// TestDecisionDialog_Update_ExplicitKey_Alert: pocKey9 has key "T" for "Trocar senha".
func TestDecisionDialog_Update_ExplicitKey_Alert(t *testing.T) {
	d := pocKey9() // Alert 3-action: Enter Usar assim mesmo / T Trocar senha / Esc Cancelar
	cmd := d.Update(tea.KeyPressMsg{Code: 't'})
	if cmd == nil {
		t.Error("explicit key 't' (matching 'T' action) must return non-nil cmd")
	}
}

// TestDecisionDialog_Update_ExplicitKey_Informative: pocKeyC has key "F" for "Fechar".
func TestDecisionDialog_Update_ExplicitKey_Informative(t *testing.T) {
	d := pocKeyC() // Informative 3-action: Enter Ver detalhes / F Fechar / Esc OK
	cmd := d.Update(tea.KeyPressMsg{Code: 'f'})
	if cmd == nil {
		t.Error("explicit key 'f' (matching 'F' action) must return non-nil cmd")
	}
}

// TestDecisionDialog_Update_EnterOnAcknowledge: pocKey1 (Acknowledge, single action).
// Enter must return non-nil cmd.
func TestDecisionDialog_Update_EnterOnAcknowledge(t *testing.T) {
	d := pocKey1()
	cmd := d.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if cmd == nil {
		t.Error("Enter on Acknowledge dialog must return non-nil cmd (pop modal)")
	}
}

// TestDecisionDialog_Update_EscOnAcknowledge: pocKey1 has no Cancel action.
// Esc must still return non-nil cmd (fallback pop modal path).
func TestDecisionDialog_Update_EscOnAcknowledge(t *testing.T) {
	d := pocKey1()
	cmd := d.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if cmd == nil {
		t.Error("Esc on Acknowledge dialog must return non-nil cmd (fallback pop modal)")
	}
}

// TestDecisionDialog_Update_CancelKeyNotTreatedAsExplicit: pocKey3 has "Esc Cancelar" as Cancel.
// Esc must trigger cancel path, not explicit key path. Return non-nil.
func TestDecisionDialog_Update_CancelKeyNotTreatedAsExplicit(t *testing.T) {
	d := pocKey3()
	cmd := d.Update(tea.KeyPressMsg{Code: tea.KeyEscape})
	if cmd == nil {
		t.Error("Esc must trigger cancel path and return non-nil cmd")
	}
}

// TestDecisionDialog_Update_UnknownExplicitKey: pocKey3 has M/Esc/Enter.
// Key "z" must return nil.
func TestDecisionDialog_Update_UnknownExplicitKey(t *testing.T) {
	d := pocKey3()
	cmd := d.Update(tea.KeyPressMsg{Code: 'z'})
	if cmd != nil {
		t.Error("unknown key 'z' on 3-action dialog must return nil cmd")
	}
}
