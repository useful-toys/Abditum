package screen

import (
	"path/filepath"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/vault"
)

// TabSpec descreve uma aba do cabeçalho.
type TabSpec struct {
	Label string
	Area  design.WorkArea
}

// DefaultTabs é a lista fixa de abas da aplicação, na ordem de exibição.
var DefaultTabs = []TabSpec{
	{Label: "Cofre", Area: design.WorkAreaVault},
	{Label: "Modelos", Area: design.WorkAreaTemplates},
	{Label: "Config", Area: design.WorkAreaSettings},
}

// tabSpacing é o número de espaços entre abas consecutivas na linha 1.
const tabSpacing = 2

// tabHitArea registra o intervalo horizontal de uma aba para detecção de clique.
type tabHitArea struct {
	area     design.WorkArea
	colStart int // coluna inicial (inclusive)
	colEnd   int // coluna final (exclusive)
}

// HeaderView renderiza o cabeçalho fixo de 2 linhas da aplicação.
// É um consumidor passivo: o RootModel injeta contexto via setters antes de cada render.
// Não possui eventos internos, filas nem timers.
type HeaderView struct {
	vault        *vault.Manager  // nil = sem cofre (estado boas-vindas)
	searchQuery  *string         // nil = busca inativa; non-nil = ativa (pode ser "")
	activeMode   design.WorkArea // relevante só quando vault != nil
	tabPositions []tabHitArea    // posições calculadas no último Render, para mouse
}

// NewHeaderView cria uma nova instância do cabeçalho.
func NewHeaderView() *HeaderView {
	return &HeaderView{}
}

// SetVault informa qual cofre está aberto.
// nil = sem cofre — renderiza apenas "  Abditum" sem abas.
func (v *HeaderView) SetVault(m *vault.Manager) {
	v.vault = m
}

// SetSearchQuery informa o estado de busca.
// nil = busca inativa. Non-nil = ativa (string pode ser "").
func (v *HeaderView) SetSearchQuery(q *string) {
	v.searchQuery = q
}

// SetActiveMode informa qual aba está ativa.
// Só tem efeito visual quando vault != nil.
func (v *HeaderView) SetActiveMode(mode design.WorkArea) {
	v.activeMode = mode
}

// Render retorna as 2 linhas do cabeçalho concatenadas com "\n",
// cada linha com exatamente `width` colunas.
// Atualiza v.tabPositions como efeito colateral para uso em Update.
func (v *HeaderView) Render(height, width int, theme *design.Theme) string {
	vaultName := vaultDisplayName(v.vault)
	isDirty := v.vault != nil && v.vault.IsModified()

	var tabs []TabSpec
	if v.vault != nil {
		tabs = DefaultTabs
	}

	// Linha 1: título + abas
	line1, _ := RenderTitleLine(vaultName, isDirty, tabs, v.activeMode, width, theme)

	// Calcular posições das abas para detecção de clique (efeito colateral).
	v.tabPositions = computeTabPositions(tabs, v.activeMode, width, theme)

	// Linha 2: separador com conector da aba ativa e busca (se ativa)
	line2, _ := RenderSeparatorLine(tabs, v.activeMode, v.searchQuery, width, theme)

	return line1 + "\n" + line2
}

// HandleKey não processa teclas nesta view.
func (v *HeaderView) HandleKey(msg tea.KeyMsg) tea.Cmd { return nil }

// HandleEvent não processa eventos externos nesta view.
func (v *HeaderView) HandleEvent(event any) {}

// HandleTeaMsg não processa mensagens do framework nesta view.
func (v *HeaderView) HandleTeaMsg(msg tea.Msg) tea.Cmd { return nil }

// Update detecta cliques de mouse nas linhas Y=0 ou Y=1 e emite WorkAreaChangedMsg.
// Não atualiza v.activeMode — o RootModel faz isso via SetActiveMode.
func (v *HeaderView) Update(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	case tea.MouseClickMsg:
		if m.Button == tea.MouseLeft && (m.Y == 0 || m.Y == 1) {
			for _, hit := range v.tabPositions {
				if m.X >= hit.colStart && m.X < hit.colEnd {
					area := hit.area
					return func() tea.Msg {
						return WorkAreaChangedMsg{Area: area}
					}
				}
			}
		}
	}
	return nil
}

// Actions retorna nil — HeaderView não possui actions próprias.
func (v *HeaderView) Actions() []actions.Action { return nil }

// --- Helpers de render exportados ---

// RenderTab retorna a representação visual de uma aba na linha 1.
//
// Inativa: "╭ Label ╮" — ╭╮ em theme.Border.Default, label em theme.Text.Secondary.
// Ativa:   "╭──────╮" — ╭╮ e ─ em theme.Border.Default (mesmo ─ que a linha separadora).
//
// A largura total é idêntica nos dois estados.
func RenderTab(label string, active bool, theme *design.Theme) (string, int) {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Default))
	tl := borderStyle.Render(design.SymCornerTL)
	tr := borderStyle.Render(design.SymCornerTR)

	// A largura interna é: 1 espaço + len(label) + 1 espaço = lipgloss.Width(label) + 2
	innerWidth := lipgloss.Width(label) + 2

	var inner string
	if active {
		fill := strings.Repeat(design.SymBorderH, innerWidth)
		inner = borderStyle.Render(fill)
	} else {
		labelStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))
		inner = " " + labelStyle.Render(label) + " "
	}

	rendered := tl + inner + tr
	return rendered, lipgloss.Width(rendered)
}

// RenderTabConnector retorna o fragmento da linha 2 que "suspende" a aba ativa.
//
// Formato: ╯[espaço + Label + espaço]╰
//   - ╯ e ╰ em theme.Border.Default, sem fundo
//   - conteúdo interno com fundo theme.Special.Highlight
//   - Label em theme.Accent.Primary bold
func RenderTabConnector(label string, theme *design.Theme) (string, int) {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Default))
	br := borderStyle.Render(design.SymCornerBR) // ╯
	bl := borderStyle.Render(design.SymCornerBL) // ╰

	labelStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Accent.Primary)).
		Bold(true).
		Background(lipgloss.Color(theme.Special.Highlight))
	spaceStyle := lipgloss.NewStyle().Background(lipgloss.Color(theme.Special.Highlight))

	inner := spaceStyle.Render(" ") + labelStyle.Render(label) + spaceStyle.Render(" ")

	rendered := br + inner + bl
	return rendered, lipgloss.Width(rendered)
}

// RenderTitleLine monta a linha 1 completa com exatamente `width` colunas.
//
// Sem cofre (vaultName == ""):
//
//	"  Abditum" + padding até width
//
// Com cofre:
//
//	"  Abditum · nome •   [tab1]  [tab2]  [tab3]"
func RenderTitleLine(
	vaultName string,
	isDirty bool,
	tabs []TabSpec,
	activeMode design.WorkArea,
	width int,
	theme *design.Theme,
) (string, int) {
	appStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Accent.Primary)).Bold(true)
	appName := "  " + appStyle.Render("Abditum")

	if vaultName == "" {
		// Sem cofre: só o nome da aplicação, preenchido até width
		line := appName + strings.Repeat(" ", max(0, width-lipgloss.Width(appName)))
		return line, lipgloss.Width(line)
	}

	// Bloco de abas
	tabsBlock, tabsWidth := renderTabsBlock(tabs, activeMode, theme)

	// Cálculo de espaço disponível para o nome do cofre
	// Prefixo: "  Abditum · " = 2 espaços + "Abditum" (7) + " · " (3) = 12 colunas
	prefixWidth := lipgloss.Width("  Abditum · ")
	dirtyWidth := 0
	if isDirty {
		dirtyWidth = lipgloss.Width(" " + design.SymBullet)
	}
	const paddingMin = 1
	disponivel := width - prefixWidth - dirtyWidth - paddingMin - tabsWidth

	// Truncar nome do cofre se necessário
	name := truncateRight(vaultName, disponivel)

	// Montar prefixo completo
	sepStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Default))
	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))
	dirtyStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Semantic.Warning))

	prefix := appName +
		" " + sepStyle.Render(design.SymHeaderSep) + " " +
		nameStyle.Render(name)

	if isDirty {
		prefix += " " + dirtyStyle.Render(design.SymBullet)
	}

	// Padding entre nome/dirty e abas
	usedLeft := lipgloss.Width(prefix)
	padding := width - usedLeft - tabsWidth
	if padding < paddingMin {
		padding = paddingMin
	}

	line := prefix + strings.Repeat(" ", padding) + tabsBlock
	// Garantir largura exata
	lineWidth := lipgloss.Width(line)
	if lineWidth < width {
		line += strings.Repeat(" ", width-lineWidth)
	}
	return line, lipgloss.Width(line)
}

// RenderSeparatorLine monta a linha 2 completa com exatamente `width` colunas.
//
// Sem abas: SymBorderH repetido width vezes.
// Com abas, busca inativa: ─── + conector da aba ativa + ───
// Com abas, busca ativa: " ─ Busca: " + query + ─── + conector + ───
func RenderSeparatorLine(
	tabs []TabSpec,
	activeMode design.WorkArea,
	searchQuery *string,
	width int,
	theme *design.Theme,
) (string, int) {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Default))
	hBorder := func(n int) string {
		if n <= 0 {
			return ""
		}
		return borderStyle.Render(strings.Repeat(design.SymBorderH, n))
	}

	// Sem abas: linha de borda horizontal simples
	if len(tabs) == 0 {
		line := hBorder(width)
		return line, lipgloss.Width(line)
	}

	// Encontrar aba ativa e sua posição horizontal na linha 1
	activeLabel := ""
	activeColStart := 0
	col := 0
	for i, tab := range tabs {
		if i > 0 {
			col += tabSpacing
		}
		_, tw := RenderTab(tab.Label, tab.Area == activeMode, theme)
		if tab.Area == activeMode {
			activeLabel = tab.Label
			activeColStart = col
		}
		col += tw
	}

	connector, connectorWidth := RenderTabConnector(activeLabel, theme)

	// Posição do conector na linha 2: alinhado com ╭ da aba ativa na linha 1.
	// O prefixo "  Abditum · nome • " tem uma largura variável, mas as abas começam
	// num offset que depende do width e da largura do bloco de abas.
	// Para simplificar: o conector começa na mesma coluna que a aba ativa.
	// Calculamos o offset das abas em relação ao início da linha usando a mesma
	// lógica de renderTabsBlock — as abas ficam à direita na linha 1.
	_, tabsWidth := renderTabsBlock(tabs, activeMode, theme)
	tabsStartCol := width - tabsWidth
	connectorStartCol := tabsStartCol + activeColStart
	// Ajuste para "  " que precede a primeira aba (já incluído em tabsBlock se
	// usarmos tabSpacing entre abas; aqui as abas não têm prefixo, então o
	// connectorStartCol é direto).

	if searchQuery == nil {
		// Busca inativa: ─── + conector + ───
		leftCount := connectorStartCol
		rightCount := width - connectorStartCol - connectorWidth
		if rightCount < 0 {
			rightCount = 0
		}
		line := hBorder(leftCount) + connector + hBorder(rightCount)
		lineWidth := lipgloss.Width(line)
		if lineWidth < width {
			line += hBorder(width - lineWidth)
		}
		return line, lipgloss.Width(line)
	}

	// Busca ativa
	// O connector deve estar alinhado com a aba ativa na linha 1.
	// Prefixo: " ─ Busca: " (10 colunas = 1 espaço + 1 char borda + 1 espaço + "Busca:" + 1 espaço)
	searchPrefix := " " + design.SymBorderH + " Busca: "
	searchPrefixRendered := borderStyle.Render(searchPrefix)
	const searchPrefixCols = 10

	// Espaço entre prefixo da busca e início do connector
	prefixToConnector := connectorStartCol - searchPrefixCols

	//	query pode ocupar até prefixToConnector - 1 colunas (precisamos de pelo menos 1 borda antes do connector)
	querySpaceAvailable := prefixToConnector - 1
	if querySpaceAvailable < 1 {
		querySpaceAvailable = 1
	}
	query := truncateLeft(*searchQuery, querySpaceAvailable)
	queryStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Accent.Primary)).Bold(true)
	queryRendered := queryStyle.Render(query)
	queryWidth := lipgloss.Width(queryRendered)

	// Prefixo da busca + query + preenchimento até connector
	prefixQueryFill := searchPrefixRendered + queryRendered
	searchToConnectorWidth := searchPrefixCols + queryWidth
	fillBetween := connectorStartCol - searchToConnectorWidth
	if fillBetween < 1 {
		fillBetween = 1
	}

	line := prefixQueryFill + hBorder(fillBetween) + connector
	lineWidth := lipgloss.Width(line)
	if lineWidth < width {
		line += hBorder(width - lineWidth)
	}
	return line, lipgloss.Width(line)
}

// --- Helpers internos ---

// vaultDisplayName extrai o nome de exibição do cofre a partir do caminho do arquivo.
// Retorna "" se m for nil.
func vaultDisplayName(m *vault.Manager) string {
	if m == nil {
		return ""
	}
	base := filepath.Base(m.FilePath())
	return strings.TrimSuffix(base, ".abditum")
}

// renderTabsBlock renderiza todas as abas em sequência, com tabSpacing entre elas.
// Retorna o bloco renderizado e sua largura visual total.
func renderTabsBlock(tabs []TabSpec, activeMode design.WorkArea, theme *design.Theme) (string, int) {
	var sb strings.Builder
	total := 0
	for i, tab := range tabs {
		if i > 0 {
			sb.WriteString(strings.Repeat(" ", tabSpacing))
			total += tabSpacing
		}
		rendered, w := RenderTab(tab.Label, tab.Area == activeMode, theme)
		sb.WriteString(rendered)
		total += w
	}
	return sb.String(), total
}

// computeTabPositions calcula os intervalos horizontais de cada aba na linha 1
// para uso na detecção de cliques do mouse.
func computeTabPositions(tabs []TabSpec, activeMode design.WorkArea, width int, theme *design.Theme) []tabHitArea {
	if len(tabs) == 0 {
		return nil
	}
	_, tabsWidth := renderTabsBlock(tabs, activeMode, theme)
	tabsStartCol := width - tabsWidth

	var positions []tabHitArea
	col := tabsStartCol
	for i, tab := range tabs {
		if i > 0 {
			col += tabSpacing
		}
		_, tw := RenderTab(tab.Label, tab.Area == activeMode, theme)
		positions = append(positions, tabHitArea{
			area:     tab.Area,
			colStart: col,
			colEnd:   col + tw,
		})
		col += tw
	}
	return positions
}

// truncateRight trunca s à direita para caber em maxCols colunas visuais.
// Adiciona design.SymEllipsis se necessário.
func truncateRight(s string, maxCols int) string {
	runes := []rune(s)
	if lipgloss.Width(s) <= maxCols {
		return s
	}
	if maxCols < 2 {
		return design.SymEllipsis
	}
	// Reduzir runa por runa até caber com ellipsis
	for len(runes) > 0 {
		candidate := string(runes) + design.SymEllipsis
		if lipgloss.Width(candidate) <= maxCols {
			return candidate
		}
		runes = runes[:len(runes)-1]
	}
	return design.SymEllipsis
}

// truncateLeft trunca s à esquerda para caber em maxCols colunas visuais.
// Adiciona design.SymEllipsis no início se necessário.
func truncateLeft(s string, maxCols int) string {
	runes := []rune(s)
	if lipgloss.Width(s) <= maxCols {
		return s
	}
	if maxCols < 2 {
		return design.SymEllipsis
	}
	// Reduzir do início até caber com ellipsis
	for len(runes) > 0 {
		candidate := design.SymEllipsis + string(runes)
		if lipgloss.Width(candidate) <= maxCols {
			return candidate
		}
		runes = runes[1:]
	}
	return design.SymEllipsis
}

// max retorna o maior de dois inteiros.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
