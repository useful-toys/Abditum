package settings

import (
	"strings"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/testdata"
	"github.com/useful-toys/abditum/internal/vault"
)

// stubMC é um MessageController mínimo para capturar hints e erros nos testes.
type stubMC struct {
	lastHint  string
	lastError string
}

func (s *stubMC) SetHintField(text string) { s.lastHint = text }
func (s *stubMC) SetError(text string)     { s.lastError = text }

// novoManager cria um vault.Manager com cofre padrão para uso nos testes.
func novoManager() *vault.Manager {
	cofre := vault.NovoCofre()
	if err := cofre.InicializarConteudoPadrao(); err != nil {
		panic("falha ao inicializar cofre: " + err.Error())
	}
	return vault.NewManagerForTest(cofre, "meu-cofre.abd")
}

// keyDown retorna uma tecla ↓.
func keyDown() tea.KeyMsg { return tea.KeyPressMsg{Code: tea.KeyDown} }

// keyUp retorna uma tecla ↑.
func keyUp() tea.KeyMsg { return tea.KeyPressMsg{Code: tea.KeyUp} }

// keyEnter retorna uma tecla Enter.
func keyEnter() tea.KeyMsg { return tea.KeyPressMsg{Code: tea.KeyEnter} }

// keyEsc retorna uma tecla Esc.
func keyEsc() tea.KeyMsg { return tea.KeyPressMsg{Code: tea.KeyEscape} }

// keyChar retorna a tecla do caractere fornecido.
func keyChar(r rune) tea.KeyMsg { return tea.KeyPressMsg{Code: r} }

// keyBackspace retorna a tecla Backspace.
func keyBackspace() tea.KeyMsg { return tea.KeyPressMsg{Code: tea.KeyBackspace} }

// --- Testes Golden ---

// TestSettingsView_Golden_Navegacao verifica o layout da tela no estado padrão de navegação
// (cofre aberto, foco no primeiro item selecionável).
func TestSettingsView_Golden_Navegacao(t *testing.T) {
	testdata.TestRenderManaged(t, "settings", "navegacao", []string{"80x20"},
		func(w, h int, theme *design.Theme) string {
			vm := novoManager()
			v := NewSettingsView(vm, &stubMC{}, "v0.1.0")
			return v.Render(h, w, theme)
		},
	)
}

// TestSettingsView_Golden_EdicaoNumerica verifica o layout com campo numérico em edição.
func TestSettingsView_Golden_EdicaoNumerica(t *testing.T) {
	testdata.TestRenderManaged(t, "settings", "edicao-numerica", []string{"80x20"},
		func(w, h int, theme *design.Theme) string {
			vm := novoManager()
			v := NewSettingsView(vm, &stubMC{}, "v0.1.0")
			// Mover até o campo "Bloqueio por inatividade" (índice 1 nos selecionáveis)
			v.HandleKey(keyDown()) // → Bloqueio por inatividade
			// Entrar em modo de edição
			v.HandleKey(keyEnter())
			// Digitar "120"
			v.HandleKey(keyChar('1'))
			v.HandleKey(keyChar('2'))
			v.HandleKey(keyChar('0'))
			return v.Render(h, w, theme)
		},
	)
}

// TestSettingsView_Golden_SemCofre verifica o layout quando nenhum cofre está aberto.
func TestSettingsView_Golden_SemCofre(t *testing.T) {
	testdata.TestRenderManaged(t, "settings", "sem-cofre", []string{"80x20"},
		func(w, h int, theme *design.Theme) string {
			v := NewSettingsView(nil, &stubMC{}, "v0.1.0")
			return v.Render(h, w, theme)
		},
	)
}

// --- Testes de navegação ---

// TestSettingsView_Navegacao_DesceEmpessoa verifica que ↓ move o foco para baixo.
func TestSettingsView_Navegacao_DesceEmpessoa(t *testing.T) {
	v := NewSettingsView(nil, &stubMC{}, "v")
	inicial := v.cursor
	v.HandleKey(keyDown())
	if v.cursor != inicial+1 {
		t.Errorf("↓ esperava cursor=%d, obteve %d", inicial+1, v.cursor)
	}
}

// TestSettingsView_Navegacao_SubeDoZeroVaiParaUltimo verifica wrapping de ↑ no primeiro item.
func TestSettingsView_Navegacao_SubeDoZeroVaiParaUltimo(t *testing.T) {
	v := NewSettingsView(nil, &stubMC{}, "v")
	if v.cursor != 0 {
		t.Fatalf("cursor inicial esperado 0, obteve %d", v.cursor)
	}
	v.HandleKey(keyUp())
	sel := v.selecionaveis()
	ultimo := len(sel) - 1
	if v.cursor != ultimo {
		t.Errorf("↑ no primeiro item: esperava cursor=%d (último), obteve %d", ultimo, v.cursor)
	}
}

// TestSettingsView_Navegacao_DesceDoUltimoVaiParaZero verifica wrapping de ↓ no último item.
func TestSettingsView_Navegacao_DesceDoUltimoVaiParaZero(t *testing.T) {
	v := NewSettingsView(nil, &stubMC{}, "v")
	sel := v.selecionaveis()
	v.cursor = len(sel) - 1
	v.HandleKey(keyDown())
	if v.cursor != 0 {
		t.Errorf("↓ no último item: esperava cursor=0, obteve %d", v.cursor)
	}
}

// TestSettingsView_Navegacao_GruposNaoRecebemFoco garante que cabeçalhos de grupo são pulados.
func TestSettingsView_Navegacao_GruposNaoRecebemFoco(t *testing.T) {
	v := NewSettingsView(nil, &stubMC{}, "v")
	for _, sel := range v.selecionaveis() {
		if v.items[sel].tipo == tipoGrupo {
			t.Errorf("índice %d é tipoGrupo mas aparece nos selecionáveis", sel)
		}
	}
}

// --- Testes do item de tema ---

// TestSettingsView_TemaFocado_HintCorreto verifica que focar o item de tema emite o hint certo.
func TestSettingsView_TemaFocado_HintCorreto(t *testing.T) {
	mc := &stubMC{}
	v := NewSettingsView(nil, mc, "v")
	// O item de tema é o primeiro selecionável (índice 0).
	if v.cursor != 0 {
		t.Fatalf("cursor inicial esperado 0, obteve %d", v.cursor)
	}
	v.emitirHintFoco()
	if !strings.Contains(mc.lastHint, "F12") {
		t.Errorf("hint do tema deveria mencionar F12, obteve %q", mc.lastHint)
	}
}

// TestSettingsView_TemaFocado_EnterIgnorado garante que Enter sobre o tema não ativa edição.
func TestSettingsView_TemaFocado_EnterIgnorado(t *testing.T) {
	v := NewSettingsView(nil, &stubMC{}, "v")
	// cursor=0 → tema
	v.HandleKey(keyEnter())
	if v.editMode {
		t.Error("Enter sobre item de tema não deve ativar editMode")
	}
}

// TestSettingsView_Tema_RefleteThemeAtual verifica que syncTema atualiza o item de tema.
func TestSettingsView_Tema_RefleteThemeAtual(t *testing.T) {
	v := NewSettingsView(nil, &stubMC{}, "v")
	theme := design.Cyberpunk
	v.syncTema(theme)
	for _, it := range v.items {
		if it.tipo == tipoTema {
			if it.textoVal != theme.Name {
				t.Errorf("item de tema deveria exibir %q, exibe %q", theme.Name, it.textoVal)
			}
			return
		}
	}
	t.Error("item de tema não encontrado")
}

// --- Testes de edição numérica ---

// TestSettingsView_EdicaoNumerica_EnterAtiva verifica que Enter sobre numérico ativa editMode.
func TestSettingsView_EdicaoNumerica_EnterAtiva(t *testing.T) {
	v := NewSettingsView(nil, &stubMC{}, "v")
	v.HandleKey(keyDown()) // → Bloqueio por inatividade (tipoNumerico)
	v.HandleKey(keyEnter())
	if !v.editMode {
		t.Error("Enter sobre tipoNumerico deveria ativar editMode")
	}
}

// TestSettingsView_EdicaoNumerica_SomenteDigitosAceitos garante que apenas dígitos são aceitos.
func TestSettingsView_EdicaoNumerica_SomenteDigitosAceitos(t *testing.T) {
	v := NewSettingsView(nil, &stubMC{}, "v")
	v.HandleKey(keyDown())
	v.HandleKey(keyEnter())
	// Limpar o buffer pré-preenchido
	for range v.editBuf {
		v.HandleKey(keyBackspace())
	}
	v.HandleKey(keyChar('5'))
	v.HandleKey(keyChar('a'))   // letra: ignorada
	v.HandleKey(keyChar('+'))   // símbolo: ignorado
	if v.editBuf != "5" {
		t.Errorf("buffer deveria ser %q, obteve %q", "5", v.editBuf)
	}
}

// TestSettingsView_EdicaoNumerica_BackspaceRemoveUltimoDigito verifica Backspace no buffer.
func TestSettingsView_EdicaoNumerica_BackspaceRemoveUltimoDigito(t *testing.T) {
	v := NewSettingsView(nil, &stubMC{}, "v")
	v.HandleKey(keyDown())
	v.HandleKey(keyEnter())
	for range v.editBuf {
		v.HandleKey(keyBackspace())
	}
	v.HandleKey(keyChar('9'))
	v.HandleKey(keyChar('0'))
	v.HandleKey(keyBackspace())
	if v.editBuf != "9" {
		t.Errorf("após Backspace buffer deveria ser %q, obteve %q", "9", v.editBuf)
	}
}

// TestSettingsView_EdicaoNumerica_EscRestauraSemSalvar verifica que Esc cancela sem persistir.
func TestSettingsView_EdicaoNumerica_EscRestauraSemSalvar(t *testing.T) {
	vm := novoManager()
	v := NewSettingsView(vm, &stubMC{}, "v")
	v.HandleKey(keyDown()) // → Bloqueio
	v.HandleKey(keyEnter())

	valorOriginal := v.items[v.cursorIndex()].valor
	for range v.editBuf {
		v.HandleKey(keyBackspace())
	}
	v.HandleKey(keyChar('9'))
	v.HandleKey(keyChar('9'))
	v.HandleKey(keyChar('9'))
	v.HandleKey(keyEsc())

	if v.editMode {
		t.Error("Esc deveria sair do editMode")
	}
	valorApos := v.items[v.cursorIndex()].valor
	if valorApos != valorOriginal {
		t.Errorf("Esc deveria restaurar o valor original %d, obteve %d", valorOriginal, valorApos)
	}
}

// TestSettingsView_EdicaoNumerica_ValorValidoConfirmaEAplica verifica confirmação com valor válido.
func TestSettingsView_EdicaoNumerica_ValorValidoConfirmaEAplica(t *testing.T) {
	vm := novoManager()
	mc := &stubMC{}
	v := NewSettingsView(vm, mc, "v")
	v.HandleKey(keyDown()) // → Bloqueio por inatividade
	v.HandleKey(keyEnter())
	for range v.editBuf {
		v.HandleKey(keyBackspace())
	}
	// Digitar 120 (> 60, portanto válido)
	v.HandleKey(keyChar('1'))
	v.HandleKey(keyChar('2'))
	v.HandleKey(keyChar('0'))
	v.HandleKey(keyEnter())

	if v.editMode {
		t.Error("Enter com valor válido deveria sair do editMode")
	}
	if mc.lastError != "" {
		t.Errorf("não deveria haver erro, obteve %q", mc.lastError)
	}
	// Confirmar que o domínio foi atualizado
	cfg := vm.Vault().Configuracoes()
	if cfg.TempoBloqueioSegundos() != 120 {
		t.Errorf("domínio deveria ter bloqueio=120, obteve %d", cfg.TempoBloqueioSegundos())
	}
}

// TestSettingsView_EdicaoNumerica_ValorForaDaRangeExibeErroEMantemEdicao verifica rejeição.
func TestSettingsView_EdicaoNumerica_ValorForaDaRangeExibeErroEMantemEdicao(t *testing.T) {
	mc := &stubMC{}
	v := NewSettingsView(nil, mc, "v")
	v.HandleKey(keyDown()) // → Bloqueio (mínimo 61)
	v.HandleKey(keyEnter())
	for range v.editBuf {
		v.HandleKey(keyBackspace())
	}
	// Digitar 10 (< 61, inválido)
	v.HandleKey(keyChar('1'))
	v.HandleKey(keyChar('0'))
	v.HandleKey(keyEnter())

	if !v.editMode {
		t.Error("valor inválido deveria manter editMode")
	}
	if mc.lastError == "" {
		t.Error("deveria haver mensagem de erro")
	}
}

// --- Testes de hints ---

// TestSettingsView_Hint_NumericoFocado verifica hint de campo numérico focado.
func TestSettingsView_Hint_NumericoFocado(t *testing.T) {
	mc := &stubMC{}
	v := NewSettingsView(nil, mc, "v")
	v.HandleKey(keyDown()) // → tipoNumerico
	v.emitirHintFoco()
	if !strings.Contains(mc.lastHint, "Enter") || !strings.Contains(mc.lastHint, "+") {
		t.Errorf("hint numérico deveria mencionar Enter e +/-, obteve %q", mc.lastHint)
	}
}

// TestSettingsView_Hint_EmEdicao verifica hint durante edição.
func TestSettingsView_Hint_EmEdicao(t *testing.T) {
	mc := &stubMC{}
	v := NewSettingsView(nil, mc, "v")
	v.HandleKey(keyDown())
	v.HandleKey(keyEnter()) // entra em editMode → deve emitir hint de edição
	if !strings.Contains(mc.lastHint, "Esc") {
		t.Errorf("hint de edição deveria mencionar Esc, obteve %q", mc.lastHint)
	}
}

// --- Testes de ajuste rápido ---

// TestSettingsView_AjusteRapido_MaisIncrementa verifica que + incrementa o campo numérico.
func TestSettingsView_AjusteRapido_MaisIncrementa(t *testing.T) {
	vm := novoManager()
	v := NewSettingsView(vm, &stubMC{}, "v")
	v.HandleKey(keyDown()) // → Bloqueio por inatividade
	antes := v.items[v.cursorIndex()].valor
	v.HandleKey(keyChar('+'))
	depois := v.items[v.cursorIndex()].valor
	if depois != antes+passoSegundos {
		t.Errorf("+ deveria incrementar %d → %d, obteve %d", antes, antes+passoSegundos, depois)
	}
}

// TestSettingsView_AjusteRapido_MenosDecrementa verifica que - decrementa respeitando o mínimo.
func TestSettingsView_AjusteRapido_MenosDecrementa(t *testing.T) {
	vm := novoManager()
	v := NewSettingsView(vm, &stubMC{}, "v")
	v.HandleKey(keyDown()) // → Bloqueio (default 300)
	v.HandleKey(keyChar('-'))
	depois := v.items[v.cursorIndex()].valor
	if depois != 300-passoSegundos {
		t.Errorf("- deveria decrementar 300 → %d, obteve %d", 300-passoSegundos, depois)
	}
}

// TestSettingsView_AjusteRapido_NaoDescaAbaixoDoMinimo verifica que o decremento respeita o mínimo.
func TestSettingsView_AjusteRapido_NaoDescaAbaixoDoMinimo(t *testing.T) {
	vm := novoManager()
	v := NewSettingsView(vm, &stubMC{}, "v")
	v.HandleKey(keyDown()) // → Bloqueio (mínimo 61)
	// Forçar valor próximo do mínimo
	idx := v.cursorIndex()
	v.items[idx].valor = minBloqueioSegundos + 1
	v.HandleKey(keyChar('-'))
	depois := v.items[v.cursorIndex()].valor
	if depois < minBloqueioSegundos {
		t.Errorf("valor não deveria descer abaixo de %d, obteve %d", minBloqueioSegundos, depois)
	}
}

// --- Testes de mutação no domínio ---

// TestSettingsView_Mutacao_AjusteRapidoPropagaAoDominio garante que + persiste no vault.
func TestSettingsView_Mutacao_AjusteRapidoPropagaAoDominio(t *testing.T) {
	vm := novoManager()
	v := NewSettingsView(vm, &stubMC{}, "v")
	v.HandleKey(keyDown()) // → Bloqueio
	valorAntes := vm.Vault().Configuracoes().TempoBloqueioSegundos()
	v.HandleKey(keyChar('+'))
	valorDepois := vm.Vault().Configuracoes().TempoBloqueioSegundos()
	if valorDepois != valorAntes+passoSegundos {
		t.Errorf("domínio deveria ter bloqueio=%d, obteve %d", valorAntes+passoSegundos, valorDepois)
	}
}

// TestSettingsView_Mutacao_CofresemAberto_NaoPanicaNemMuta garante comportamento sem vault.
func TestSettingsView_Mutacao_SemCofre_NaoPanicaNemMuta(t *testing.T) {
	// Nenhum panic deve ocorrer; HandleKey deve simplesmente não fazer nada no domínio.
	v := NewSettingsView(nil, &stubMC{}, "v")
	v.HandleKey(keyDown())
	v.HandleKey(keyChar('+'))
	v.HandleKey(keyChar('-'))
	v.HandleKey(keyEnter())
	for range v.editBuf {
		v.HandleKey(keyBackspace())
	}
	v.HandleKey(keyChar('1'))
	v.HandleKey(keyChar('2'))
	v.HandleKey(keyChar('0'))
	v.HandleKey(keyEnter())
	// Sem panic = sucesso
}

// --- Teste de cursor real ---

// TestSettingsView_Cursor_NilForaNaEdicao verifica que Cursor retorna nil fora do editMode.
func TestSettingsView_Cursor_NilForaNaEdicao(t *testing.T) {
	v := NewSettingsView(nil, &stubMC{}, "v")
	v.Render(20, 80, design.TokyoNight)
	if c := v.Cursor(2, 0); c != nil {
		t.Errorf("Cursor fora de editMode deveria ser nil, obteve %v", c)
	}
}

// TestSettingsView_Cursor_NaoNilEmEdicao verifica que Cursor não é nil durante editMode.
func TestSettingsView_Cursor_NaoNilEmEdicao(t *testing.T) {
	v := NewSettingsView(nil, &stubMC{}, "v")
	v.HandleKey(keyDown())
	v.HandleKey(keyEnter()) // entra em editMode
	v.Render(20, 80, design.TokyoNight)
	if c := v.Cursor(2, 0); c == nil {
		t.Error("Cursor em editMode deveria ser não-nil")
	}
}
