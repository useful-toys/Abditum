package settings

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)

// Limites e passo dos campos numéricos, conforme golden/requisitos.md.
const (
	minBloqueioSegundos  = 61 // bloqueio por inatividade: mínimo estritamente acima de 60 s
	minOcultarSegundos   = 3  // ocultação de campo sensível: mínimo estritamente acima de 2 s
	minClipboardSegundos = 11 // limpeza de clipboard: mínimo estritamente acima de 10 s
	passoSegundos        = 5  // passo de ajuste rápido com +/-
)

// Colunas de referência para alinhamento dos itens (visual, em caracteres).
const (
	indentGrupo  = 30 // coluna onde os cabeçalhos de grupo começam
	indentItem   = 30 // coluna de início dos itens (largura do prefixo)
	colunaValor  = 62 // coluna onde o valor exibido começa
	larguraInput = 4  // largura visual do campo de entrada numérica
)

// tipoItem classifica o comportamento de um settingItem.
type tipoItem int

const (
	tipoGrupo          tipoItem = iota // cabeçalho de grupo, não selecionável
	tipoNumerico                       // campo timer editável com entrada numérica
	tipoTema                           // somente focável; troca via F12 global
	tipoSomenteLeitura                 // somente exibição, não interativo
)

// settingItem representa um item na lista de configurações.
type settingItem struct {
	label     string
	descricao string   // exibida sob o item quando focado
	tipo      tipoItem
	valor     int    // valor atual, usado em tipoNumerico
	textoVal  string // valor textual, usado em tipoTema e tipoSomenteLeitura
	minimo    int    // valor mínimo aceito (inclusive), usado em tipoNumerico
	campo     string // identificador do campo para aplicar ao vault: "bloqueio", "ocultar", "clipboard"
}

// messageController é a interface mínima da barra de mensagens necessária à SettingsView.
// Subconjunto de tui.MessageController para evitar import cycle entre os pacotes.
type messageController interface {
	SetHintField(text string)
	SetError(text string)
}

// SettingsView exibe as opções de configuração da aplicação.
type SettingsView struct {
	vm      *vault.Manager
	mc      messageController
	version string
	items   []settingItem

	// cursor é o índice entre os itens selecionáveis (não inclui tipoGrupo).
	cursor   int
	editMode bool   // true quando um campo numérico está em edição inline
	editBuf  string // dígitos digitados durante a edição
	editSnap int    // valor antes da edição, para cancelamento com Esc restaurar

	// Posição do cursor real do terminal na work area (linha, coluna),
	// atualizada em Render e consultada em Cursor().
	cursorWorkAreaRow int
	cursorWorkAreaCol int
}

// NewSettingsView cria uma nova instância da tela de configurações.
// vm pode ser nil quando nenhum cofre estiver aberto; mc deve ser não-nil.
func NewSettingsView(vm *vault.Manager, mc messageController, version string) *SettingsView {
	v := &SettingsView{vm: vm, mc: mc, version: version}
	v.rebuildItems()
	return v
}

// rebuildItems reconstrói a lista de itens usando os valores atuais do cofre.
func (v *SettingsView) rebuildItems() {
	var bloqueio, ocultar, clipboard int
	temaAtual := "–"
	nomeArquivo := "–"

	if v.vm != nil {
		cofre := v.vm.Vault()
		if cofre != nil {
			cfg := cofre.Configuracoes()
			bloqueio = cfg.TempoBloqueioSegundos()
			ocultar = cfg.TempoOcultarSegundos()
			clipboard = cfg.TempoLimparTransferenciaSegundos()
			if t := cfg.TemaVisual(); t != "" {
				temaAtual = t
			}
		}
		if fp := v.vm.FilePath(); fp != "" {
			nomeArquivo = filepath.Base(fp)
		}
	}

	v.items = []settingItem{
		{label: "Aparência", tipo: tipoGrupo},
		{
			label:     "Tema visual",
			descricao: "Tema aplicado ao cofre atual.",
			tipo:      tipoTema,
			textoVal:  temaAtual,
		},
		{label: "Segurança", tipo: tipoGrupo},
		{
			label:     "Bloqueio por inatividade",
			descricao: "Tempo de bloqueio automático por inatividade.",
			tipo:      tipoNumerico,
			valor:     bloqueio,
			minimo:    minBloqueioSegundos,
			campo:     "bloqueio",
		},
		{
			label:     "Ocultar campo sensível",
			descricao: "Tempo até ocultar automaticamente um campo sensível revelado.",
			tipo:      tipoNumerico,
			valor:     ocultar,
			minimo:    minOcultarSegundos,
			campo:     "ocultar",
		},
		{
			label:     "Limpar área de transferência",
			descricao: "Tempo até limpar automaticamente o conteúdo copiado.",
			tipo:      tipoNumerico,
			valor:     clipboard,
			minimo:    minClipboardSegundos,
			campo:     "clipboard",
		},
		{label: "Sobre", tipo: tipoGrupo},
		{
			label:    "Versão",
			tipo:     tipoSomenteLeitura,
			textoVal: v.version,
		},
		{
			label:    "Arquivo do cofre",
			tipo:     tipoSomenteLeitura,
			textoVal: nomeArquivo,
		},
	}
}

// selecionaveis retorna os índices em v.items dos itens que recebem foco (não são grupo).
func (v *SettingsView) selecionaveis() []int {
	var idx []int
	for i, it := range v.items {
		if it.tipo != tipoGrupo {
			idx = append(idx, i)
		}
	}
	return idx
}

// cursorIndex retorna o índice real em v.items do item com foco atual.
func (v *SettingsView) cursorIndex() int {
	sel := v.selecionaveis()
	if len(sel) == 0 {
		return 0
	}
	if v.cursor >= len(sel) {
		v.cursor = len(sel) - 1
	}
	return sel[v.cursor]
}

// moverCursor desloca o foco em delta posições com wrapping e emite o hint do novo item.
func (v *SettingsView) moverCursor(delta int) {
	sel := v.selecionaveis()
	if len(sel) == 0 {
		return
	}
	v.cursor = (v.cursor + delta + len(sel)) % len(sel)
	v.emitirHintFoco()
}

// emitirHintFoco envia à barra de mensagens o hint contextual do item com foco.
func (v *SettingsView) emitirHintFoco() {
	if v.mc == nil {
		return
	}
	idx := v.cursorIndex()
	switch v.items[idx].tipo {
	case tipoTema:
		v.mc.SetHintField("F12 para alternar tema visual")
	case tipoNumerico:
		v.mc.SetHintField("Enter edita · +/- altera o valor")
	default:
		v.mc.SetHintField("")
	}
}

// Render retorna a tela de configurações para as dimensões e tema fornecidos.
func (v *SettingsView) Render(height, width int, theme *design.Theme) string {
	v.syncTema(theme)

	lines, cursorLineOffset, cursorCol := v.renderLinhas(width, theme)
	nLinhas := len(lines)

	padTop := (height - nLinhas) / 2
	if padTop < 0 {
		padTop = 0
	}
	padBottom := height - nLinhas - padTop
	if padBottom < 0 {
		padBottom = 0
	}

	// Atualiza posição do cursor na work area para uso em Cursor().
	if v.editMode && cursorLineOffset >= 0 {
		v.cursorWorkAreaRow = padTop + cursorLineOffset
		v.cursorWorkAreaCol = cursorCol
	}

	blankLine := strings.Repeat(" ", width)
	var sb strings.Builder
	for i := 0; i < padTop; i++ {
		sb.WriteString(blankLine + "\n")
	}
	sb.WriteString(strings.Join(lines, "\n"))
	for i := 0; i < padBottom; i++ {
		sb.WriteString("\n" + blankLine)
	}

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Background(lipgloss.Color(theme.Surface.Base)).
		Render(sb.String())
}

// syncTema atualiza o item de tema visual com o nome do tema ativo.
func (v *SettingsView) syncTema(theme *design.Theme) {
	for i := range v.items {
		if v.items[i].tipo == tipoTema {
			nome := theme.Name
			if v.vm != nil {
				if cofre := v.vm.Vault(); cofre != nil {
					if cfgTema := cofre.Configuracoes().TemaVisual(); cfgTema != "" {
						nome = cfgTema
					}
				}
			}
			v.items[i].textoVal = nome
			return
		}
	}
}

// renderLinhas constrói as linhas da tela e retorna o row offset do campo em edição (ou -1)
// e a coluna onde o cursor real deve ser posicionado.
func (v *SettingsView) renderLinhas(width int, theme *design.Theme) (lines []string, cursorRow, cursorCol int) {
	const (
		simboloFoco = "› " // U+203A + espaço — indicador de item focado
		semFoco     = "  " // dois espaços — mesmo alinhamento, sem símbolo
	)

	cursorRow = -1

	accentStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Accent.Primary)).
		Bold(true)
	grupoStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Primary)).
		Bold(true)
	secundarioStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Secondary))
	inputStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.Surface.Input)).
		Foreground(lipgloss.Color(theme.Text.Primary))
	highlightStyle := lipgloss.NewStyle().
		Background(lipgloss.Color(theme.Special.Highlight)).
		Foreground(lipgloss.Color(theme.Accent.Primary)).
		Bold(true)
	descrStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Secondary))

	curIdx := v.cursorIndex()

	// Título da tela.
	lines = append(lines, strings.Repeat(" ", indentGrupo)+"Configurações", "")

	prevWasGroup := false
	for i, it := range v.items {
		focado := (i == curIdx)

		switch it.tipo {
		case tipoGrupo:
			// Linha em branco antes de cada grupo (exceto o primeiro).
			if len(lines) > 2 || prevWasGroup {
				lines = append(lines, "")
			}
			lines = append(lines, strings.Repeat(" ", indentGrupo)+grupoStyle.Render(it.label))
			prevWasGroup = true

		default:
			prevWasGroup = false

			// Prefixo: símbolo de foco ou espaço equivalente.
			var prefixo string
			if focado {
				prefixo = strings.Repeat(" ", indentItem-2) + accentStyle.Render(simboloFoco)
			} else {
				prefixo = strings.Repeat(" ", indentItem) + semFoco
			}

			// Label: destacado quando focado.
			var labelStr string
			if focado {
				labelStr = highlightStyle.Render(it.label)
			} else {
				labelStr = it.label
			}

			// Valor: renderizado conforme tipo e estado.
			var valorStr string
			switch it.tipo {
			case tipoNumerico:
				if focado && v.editMode {
					// Campo em edição: fundo surface.input; cursor real ao final do buffer.
					conteudo := fmt.Sprintf("%-*s", larguraInput, v.editBuf)
					valorStr = inputStyle.Render(conteudo) + " s"
				} else if focado {
					valorStr = accentStyle.Render(fmt.Sprintf("%d s", it.valor))
				} else {
					valorStr = secundarioStyle.Render(fmt.Sprintf("%d s", it.valor))
				}
			default:
				valorStr = secundarioStyle.Render(it.textoVal)
			}

			// Alinha label e valor garantindo ao menos um espaço entre eles.
			prefixoW := lipgloss.Width(prefixo)
			labelW := lipgloss.Width(labelStr)
			espaco := colunaValor - (prefixoW + labelW)
			if espaco < 1 {
				espaco = 1
			}

			linha := prefixo + labelStr + strings.Repeat(" ", espaco) + valorStr
			if focado && v.editMode && it.tipo == tipoNumerico {
				// Registra a linha onde o cursor real deve ser posicionado.
				cursorRow = len(lines)
				// Coluna: início do valor + comprimento atual do buffer.
				cursorCol = colunaValor + len([]rune(v.editBuf))
			}
			lines = append(lines, linha)

			// Linha de descrição contextual, somente sob o item focado.
			if focado && it.descricao != "" {
				descLinha := strings.Repeat(" ", indentItem+2) + descrStyle.Render(it.descricao)
				lines = append(lines, descLinha, "")
			}
		}
	}

	return lines, cursorRow, cursorCol
}

// Cursor retorna a posição do cursor real do terminal para a work area de settings.
// baseRow e baseCol são as coordenadas absolutas do canto superior esquerdo da work area
// na tela (passadas pelo RootModel). Retorna nil quando não há campo em edição.
func (v *SettingsView) Cursor(baseRow, baseCol int) *tea.Cursor {
	if !v.editMode {
		return nil
	}
	return tea.NewCursor(baseCol+v.cursorWorkAreaCol, baseRow+v.cursorWorkAreaRow)
}

// HandleKey processa eventos de teclado conforme o estado atual.
func (v *SettingsView) HandleKey(msg tea.KeyMsg) tea.Cmd {
	if v.editMode {
		return v.handleKeyEdicao(msg)
	}
	return v.handleKeyNavegacao(msg)
}

// handleKeyNavegacao trata teclas no modo de navegação.
func (v *SettingsView) handleKeyNavegacao(msg tea.KeyMsg) tea.Cmd {
	switch {
	case design.Keys.Up.Matches(msg):
		v.moverCursor(-1)
	case design.Keys.Down.Matches(msg):
		v.moverCursor(+1)
	case design.Keys.Enter.Matches(msg):
		idx := v.cursorIndex()
		it := &v.items[idx]
		if it.tipo == tipoNumerico {
			v.editSnap = it.valor
			v.editBuf = strconv.Itoa(it.valor)
			v.editMode = true
			if v.mc != nil {
				v.mc.SetHintField("Enter confirma · Esc cancela")
			}
		}
	default:
		k := msg.Key()
		switch k.Code {
		case '+':
			v.ajustarNumerico(+passoSegundos)
		case '-':
			v.ajustarNumerico(-passoSegundos)
		}
	}
	return nil
}

// handleKeyEdicao trata teclas no modo de edição inline.
func (v *SettingsView) handleKeyEdicao(msg tea.KeyMsg) tea.Cmd {
	switch {
	case design.Keys.Enter.Matches(msg):
		v.confirmarEdicao()
	case design.Keys.Esc.Matches(msg):
		v.cancelarEdicao()
	default:
		k := msg.Key()
		if k.Code == tea.KeyBackspace && k.Mod == 0 {
			if len(v.editBuf) > 0 {
				runes := []rune(v.editBuf)
				v.editBuf = string(runes[:len(runes)-1])
			}
		} else if k.Code >= '0' && k.Code <= '9' && k.Mod == 0 {
			v.editBuf += string(rune(k.Code))
		}
		// Silenciosamente ignora qualquer outra tecla durante edição.
	}
	return nil
}

// confirmarEdicao valida o buffer e aplica o valor ao vault.
func (v *SettingsView) confirmarEdicao() {
	idx := v.cursorIndex()
	it := &v.items[idx]

	val, err := strconv.Atoi(v.editBuf)
	if err != nil || val < it.minimo {
		if v.mc != nil {
			v.mc.SetError(fmt.Sprintf("Mínimo: %d s", it.minimo))
		}
		return // permanece em editMode até o usuário corrigir ou cancelar
	}

	if v.vm != nil {
		if cofre := v.vm.Vault(); cofre != nil {
			cfg := cofre.Configuracoes()
			nova := v.novaConfig(cfg, it.campo, val)
			if err := v.vm.AlterarConfiguracoes(nova); err != nil {
				if v.mc != nil {
					v.mc.SetError(err.Error())
				}
				return
			}
		}
	}

	it.valor = val
	v.editMode = false
	v.editBuf = ""
	v.emitirHintFoco()
}

// cancelarEdicao restaura o valor original e sai do modo de edição.
func (v *SettingsView) cancelarEdicao() {
	idx := v.cursorIndex()
	v.items[idx].valor = v.editSnap
	v.editMode = false
	v.editBuf = ""
	v.emitirHintFoco()
}

// ajustarNumerico incrementa ou decrementa o campo numérico focado em passos fixos.
func (v *SettingsView) ajustarNumerico(delta int) {
	idx := v.cursorIndex()
	it := &v.items[idx]
	if it.tipo != tipoNumerico {
		return
	}
	novoValor := it.valor + delta
	if novoValor < it.minimo {
		novoValor = it.minimo
	}

	if v.vm != nil {
		if cofre := v.vm.Vault(); cofre != nil {
			cfg := cofre.Configuracoes()
			nova := v.novaConfig(cfg, it.campo, novoValor)
			if err := v.vm.AlterarConfiguracoes(nova); err != nil {
				if v.mc != nil {
					v.mc.SetError(err.Error())
				}
				return
			}
		}
	}
	it.valor = novoValor
}

// novaConfig constrói um vault.Configuracoes com o campo indicado substituído por novoValor.
func (v *SettingsView) novaConfig(cfg vault.Configuracoes, campo string, novoValor int) vault.Configuracoes {
	b := cfg.TempoBloqueioSegundos()
	o := cfg.TempoOcultarSegundos()
	c := cfg.TempoLimparTransferenciaSegundos()
	t := cfg.TemaVisual()
	switch campo {
	case "bloqueio":
		b = novoValor
	case "ocultar":
		o = novoValor
	case "clipboard":
		c = novoValor
	}
	return vault.NovasConfiguracoes(b, o, c, t)
}

// HandleEvent recebe eventos externos (ex: mudança de tema via F12).
// A sincronização de tema é feita em Render() via syncTema, sem necessidade de evento explícito.
func (v *SettingsView) HandleEvent(event any) {}

// HandleTeaMsg processa mensagens do framework Bubble Tea.
func (v *SettingsView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update não altera o estado desta view em resposta a mensagens genéricas.
func (v *SettingsView) Update(msg tea.Msg) tea.Cmd { return nil }

// Actions retorna as actions de teclado registradas no ActionManager para esta view.
func (v *SettingsView) Actions() []actions.Action {
	return []actions.Action{
		{
			Keys:        []design.Key{design.Keys.Up},
			Label:       "↑",
			Description: "Move o foco para o item anterior",
			Priority:    10,
			Visible:     false,
			OnExecute:   func() tea.Cmd { v.moverCursor(-1); return nil },
		},
		{
			Keys:        []design.Key{design.Keys.Down},
			Label:       "↓",
			Description: "Move o foco para o próximo item",
			Priority:    11,
			Visible:     false,
			OnExecute:   func() tea.Cmd { v.moverCursor(+1); return nil },
		},
	}
}
