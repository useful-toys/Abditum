package tui

import (
	"strings"
	"testing"
)

// ─────────────────────────────────────────────────────────────────────────────
// PoC fixture constructors (15 combinations: 5 severities × 3 action counts)
// These match the PoC table in root.go (keys 1–9, a–f).
// Each calls SetSize(80, 24) so boxWidth() returns a concrete value.
// ─────────────────────────────────────────────────────────────────────────────

// ── Destrutivo ───────────────────────────────────────────────────────────────

func pocKey1() *DecisionDialog {
	d := NewDecisionDialog(SeverityDestructive, IntentionAcknowledge,
		"Exclusão concluída", "Gmail foi excluído permanentemente.",
		[]DecisionAction{
			{Key: "Enter", Label: "OK", Default: true},
		})
	d.SetSize(80, 24)
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
	d.SetSize(80, 24)
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
	d.SetSize(80, 24)
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
	d.SetSize(80, 24)
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
	d.SetSize(80, 24)
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
	d.SetSize(80, 24)
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
	d.SetSize(80, 24)
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
	d.SetSize(80, 24)
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
	d.SetSize(80, 24)
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
	d.SetSize(80, 24)
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
	d.SetSize(80, 24)
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
	d.SetSize(80, 24)
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
	d.SetSize(80, 24)
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
	d.SetSize(80, 24)
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
	d.SetSize(80, 24)
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
			out := d.View()

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
			wantSymbol: SymWarn,
		},
		{
			name:       "Destrutivo contains ⚠",
			d:          pocKey1(),
			wantSymbol: SymWarn,
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
			wantAbsent: []string{SymWarn, SymError, SymInfo},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := tt.d.View()
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
	out := d.View()

	if !strings.Contains(out, "╭") {
		t.Errorf("View() missing top-left corner ╭\ngot:\n%s", out)
	}
	if !strings.Contains(out, "╰") {
		t.Errorf("View() missing bottom-left corner ╰\ngot:\n%s", out)
	}
}
