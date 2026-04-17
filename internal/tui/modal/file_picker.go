package modal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// FilePickerMode controla o comportamento do picker.
type FilePickerMode int

const (
	FilePickerOpen FilePickerMode = iota // abrir arquivo existente
	FilePickerSave                       // salvar / nomear novo arquivo
)

// FilePickerOptions são os parâmetros de construção do FilePicker.
type FilePickerOptions struct {
	Mode       FilePickerMode
	Extension  string                   // ex: ".abditum" — inclui o ponto
	InitialDir string                   // "" → CWD → fallback ~
	Suggested  string                   // Save: valor inicial do campo Arquivo:
	OnResult   func(path string) tea.Cmd // path="" se cancelado
	Messages   tui.MessageController    // nil tolerado
}

// treeNode representa uma entrada de diretório na árvore lazy.
type treeNode struct {
	path       string
	name       string
	depth      int
	expanded   bool
	loaded     bool
	children   []*treeNode
	hasSubdirs bool
}

// visibleNode é uma entrada achatada da árvore — visível no painel Estrutura.
type visibleNode struct {
	node *treeNode
}

// FilePickerModal implementa tui.ModalView para seleção de arquivo.
type FilePickerModal struct {
	mode     FilePickerMode
	ext      string
	onResult func(string) tea.Cmd
	messages tui.MessageController

	root         *treeNode
	visibleNodes []visibleNode
	treeCursor   int
	treeScroll   int

	currentPath string
	files       []string
	fileInfos   []os.FileInfo
	fileCursor  int // -1 quando vazio
	fileScroll  int

	focusPanel int // 0=árvore 1=arquivos 2=campo nome (Save apenas)

	nameField textinput.Model

	lastMaxHeight int
	lastMaxWidth  int

	readDir     func(string) ([]os.DirEntry, error)
	hintEmitted bool
	timeFmt     func(time.Time) string

	// fallbackWarning: emitido no primeiro Render() se InitialDir não era acessível
	fallbackWarning string
}

// SetReadDirForTest injeta filesystem fictício — usado exclusivamente em testes.
func (m *FilePickerModal) SetReadDirForTest(fn func(string) ([]os.DirEntry, error)) {
	m.readDir = fn
}

// SetTimeFmtForTest injeta formatação de tempo fixa — usado exclusivamente em testes.
func (m *FilePickerModal) SetTimeFmtForTest(fn func(time.Time) string) {
	m.timeFmt = fn
}

// dirRead chama m.readDir se injetado, caso contrário os.ReadDir.
func (m *FilePickerModal) dirRead(path string) ([]os.DirEntry, error) {
	if m.readDir != nil {
		return m.readDir(path)
	}
	return os.ReadDir(path)
}

// formatTime formata um time.Time usando m.timeFmt se injetado, caso contrário local.
func (m *FilePickerModal) formatTime(t time.Time) string {
	if m.timeFmt != nil {
		return m.timeFmt(t)
	}
	return t.Format("02/01/06 15:04")
}

// NewFilePicker cria e inicializa o modal.
// Chamar tui.OpenModal(NewFilePicker(opts)) para exibir.
func NewFilePicker(opts FilePickerOptions) *FilePickerModal {
	nf := textinput.New()
	nf.Prompt = ""
	if opts.Suggested != "" {
		nf.SetValue(opts.Suggested)
	}

	m := &FilePickerModal{
		mode:       opts.Mode,
		ext:        opts.Extension,
		onResult:   opts.OnResult,
		messages:   opts.Messages,
		nameField:  nf,
		fileCursor: -1,
	}

	// Resolver InitialDir
	dir := opts.InitialDir
	if dir == "" {
		if cwd, err := os.Getwd(); err == nil {
			dir = cwd
		}
	}
	if dir == "" || !dirExists(dir) {
		if home, err := os.UserHomeDir(); err == nil {
			if dir != "" {
				m.fallbackWarning = design.SymWarning + " Diretório atual inacessível — navegando para home"
			}
			dir = home
		}
	}
	m.currentPath = dir

	// Construir árvore e carregar arquivos
	m.root = m.buildTreeChain(dir)
	m.buildVisibleNodes()
	// Posicionar treeCursor no nó de dir
	for i, vn := range m.visibleNodes {
		if vn.node.path == dir {
			m.treeCursor = i
			break
		}
	}
	m.loadFiles(dir)

	return m
}

// dirExists retorna true se path é um diretório acessível.
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// buildTreeChain constrói a árvore da raiz até targetDir com todos os ancestors expandidos.
// Funciona em Unix (raiz = "/") e Windows (raiz = "C:\") via filepath.VolumeName.
func (m *FilePickerModal) buildTreeChain(targetDir string) *treeNode {
	// Determinar raiz do filesystem
	vol := filepath.VolumeName(targetDir)
	rootPath := vol + string(filepath.Separator)

	root := &treeNode{
		path:  rootPath,
		name:  rootPath,
		depth: 0,
	}

	// Dividir targetDir em segmentos a partir da raiz
	rel, err := filepath.Rel(rootPath, targetDir)
	if err != nil || rel == "." {
		// targetDir é a raiz
		m.expandNodeWith(root)
		return root
	}

	parts := strings.Split(filepath.ToSlash(rel), "/")
	current := root
	currentPath := rootPath

	for i, part := range parts {
		if part == "" {
			continue
		}
		m.expandNodeWith(current)
		childPath := filepath.Join(currentPath, part)
		var found *treeNode
		for _, ch := range current.children {
			if ch.path == childPath {
				found = ch
				break
			}
		}
		if found == nil {
			// Criar nó intermediário se não encontrado (ex: leitura falhou parcialmente)
			found = &treeNode{
				path:  childPath,
				name:  part,
				depth: i + 1,
			}
			current.children = append(current.children, found)
		}
		found.expanded = true
		currentPath = childPath
		current = found
	}
	return root
}

// expandNodeWith carrega os filhos de node usando m.dirRead, marcando hasSubdirs.
// Se já está carregado (loaded=true), não relê o disco.
func (m *FilePickerModal) expandNodeWith(node *treeNode) {
	if node.loaded {
		node.expanded = true
		return
	}
	entries, err := m.dirRead(node.path)
	node.loaded = true
	if err != nil {
		return
	}
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		child := &treeNode{
			path:  filepath.Join(node.path, e.Name()),
			name:  e.Name(),
			depth: node.depth + 1,
		}
		node.children = append(node.children, child)
	}
	node.hasSubdirs = len(node.children) > 0
	node.expanded = true
}

// buildVisibleNodes achata a árvore em m.visibleNodes (DFS pré-ordem).
// Chamado após qualquer mudança de estado da árvore.
func (m *FilePickerModal) buildVisibleNodes() {
	m.visibleNodes = m.visibleNodes[:0]
	if m.root != nil {
		m.collectVisible(m.root)
	}
}

func (m *FilePickerModal) collectVisible(node *treeNode) {
	m.visibleNodes = append(m.visibleNodes, visibleNode{node: node})
	if node.expanded {
		for _, ch := range node.children {
			m.collectVisible(ch)
		}
	}
}

// loadFiles carrega os arquivos com m.ext presentes em dir no painel de arquivos.
// Atualiza m.files, m.fileInfos, m.currentPath, m.fileCursor, m.fileScroll.
func (m *FilePickerModal) loadFiles(dir string) {
	m.currentPath = dir
	m.files = m.files[:0]
	m.fileInfos = nil
	m.fileCursor = -1
	m.fileScroll = 0

	entries, err := m.dirRead(dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		if !strings.HasSuffix(e.Name(), m.ext) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		// Exibir nome sem extensão
		nameNoExt := strings.TrimSuffix(e.Name(), m.ext)
		m.files = append(m.files, nameNoExt)
		m.fileInfos = append(m.fileInfos, info)
	}
	if len(m.files) > 0 {
		m.fileCursor = 0
	}
}

// RebuildForTest reinicializa a árvore e arquivos para o diretório especificado.
// Usado exclusivamente em testes após SetReadDirForTest.
func (m *FilePickerModal) RebuildForTest(dir string) {
	m.fallbackWarning = "" // limpa aviso emitido durante construção com disco real
	dir = filepath.Clean(dir)
	m.currentPath = dir
	m.root = m.buildTreeChain(dir)
	m.buildVisibleNodes()
	for i, vn := range m.visibleNodes {
		if vn.node.path == dir {
			m.treeCursor = i
			break
		}
	}
	m.loadFiles(dir)
}

// TreeCursorPath retorna o path do nó sob treeCursor — usado em testes.
func (m *FilePickerModal) TreeCursorPath() string {
	if m.treeCursor < 0 || m.treeCursor >= len(m.visibleNodes) {
		return ""
	}
	return filepath.ToSlash(m.visibleNodes[m.treeCursor].node.path)
}

// FileCursor retorna o índice do arquivo selecionado no painel — usado em testes.
func (m *FilePickerModal) FileCursor() int { return m.fileCursor }

// Files retorna a lista de nomes de arquivo sem extensão — usado em testes.
func (m *FilePickerModal) Files() []string { return m.files }

// FocusPanel retorna o painel com foco (0=árvore,1=arquivos,2=campo) — usado em testes.
func (m *FilePickerModal) FocusPanel() int { return m.focusPanel }

// Render retorna a representação visual do modal.
func (m *FilePickerModal) Render(maxHeight, maxWidth int, theme *design.Theme) string {
	return "TODO"
}

// HandleKey processa eventos de teclado.
func (m *FilePickerModal) HandleKey(msg tea.KeyMsg) tea.Cmd {
	return nil
}

// Update processa mensagens do Bubble Tea.
func (m *FilePickerModal) Update(msg tea.Msg) tea.Cmd {
	if key, ok := msg.(tea.KeyMsg); ok {
		return m.HandleKey(key)
	}
	return nil
}

// Cursor retorna a posição do cursor real para o modal.
func (m *FilePickerModal) Cursor(topY, leftX int) *tea.Cursor {
	return nil
}

// formatFileSize formata bytes em KB/MB/GB (base 1024, 1 casa decimal).
func formatFileSize(bytes int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case bytes >= GB:
		return formatF(float64(bytes)/GB) + " GB"
	case bytes >= MB:
		return formatF(float64(bytes)/MB) + " MB"
	default:
		return formatF(float64(bytes)/KB) + " KB"
	}
}

func formatF(f float64) string {
	// %.1f sempre produz exatamente 1 casa decimal em Go (locale-independente)
	return fmt.Sprintf("%.1f", f)
}

// FormatFileSizeForTest expõe formatFileSize para testes externos.
func FormatFileSizeForTest(bytes int64) string { return formatFileSize(bytes) }

// padRight pads s até width colunas visuais (ANSI-aware via lipgloss.Width).
func padRight(s string, width int) string {
	w := lipgloss.Width(s)
	if w >= width {
		return s
	}
	return s + strings.Repeat(" ", width-w)
}

// modalDimensions calcula as dimensões do modal a partir do terminal.
func modalDimensions(maxHeight, maxWidth int, mode FilePickerMode) (modalW, innerW, treeW, filesW, visibleH int) {
	modalW = maxWidth * 80 / 100
	if modalW > 70 {
		modalW = 70
	}
	innerW = modalW - 2
	treeW = innerW * 40 / 100
	if treeW < 8 {
		treeW = 8
	}
	filesW = innerW - treeW - 1
	if filesW < 8 {
		filesW = 8
		// Terminais muito estreitos: recuar treeW para que a soma caiba.
		if treeW+filesW+1 > innerW {
			treeW = innerW - filesW - 1
			if treeW < 1 {
				treeW = 1
			}
		}
	}

	overhead := 3 // Open: borda sup + caminho + sep painéis
	if mode == FilePickerSave {
		overhead = 5 // + sep campo + campo Arquivo:
	}
	modalH := maxHeight * 8 / 10
	visibleH = modalH - overhead
	if visibleH < 3 {
		visibleH = 3
	}
	return
}

// renderTreeSepChar retorna o caractere do separador da árvore na linha lineIdx (0-based dentro do conteúdo).
// Substitui │ por ↑/■/↓ conforme scroll. total = len(visibleNodes), vp = visibleH.
func renderTreeSepChar(lineIdx, treeScroll, total, vp int, theme *design.Theme) string {
	ss := ScrollState{Offset: treeScroll, Total: total, Viewport: vp}
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Default))
	lineNum := lineIdx + 1 // 1-based
	switch {
	case lineNum == 1 && ss.CanScrollUp():
		s, _ := design.RenderScrollArrow(true, theme)
		return s
	case lineNum == vp && ss.CanScrollDown():
		s, _ := design.RenderScrollArrow(false, theme)
		return s
	case ss.ThumbLine() == lineNum:
		s, _ := design.RenderScrollThumb(theme)
		return s
	default:
		return borderStyle.Render(design.SymBorderV)
	}
}

// renderFileSepChar retorna o caractere da borda direita do modal na linha lineIdx.
// Usa theme.Border.Focused porque esta borda pertence ao contorno externo do modal,
// que é sempre renderizado como focused independentemente do painel ativo.
func renderFileSepChar(lineIdx, fileScroll, total, vp int, theme *design.Theme) string {
	ss := ScrollState{Offset: fileScroll, Total: total, Viewport: vp}
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Focused))
	lineNum := lineIdx + 1
	switch {
	case lineNum == 1 && ss.CanScrollUp():
		s, _ := design.RenderScrollArrow(true, theme)
		return s
	case lineNum == vp && ss.CanScrollDown():
		s, _ := design.RenderScrollArrow(false, theme)
		return s
	case ss.ThumbLine() == lineNum:
		s, _ := design.RenderScrollThumb(theme)
		return s
	default:
		return borderStyle.Render(design.SymBorderV)
	}
}

// renderTopBorder renderiza ╭── título ──╮
func (m *FilePickerModal) renderTopBorder(modalW int, theme *design.Theme) string {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Focused))

	title := "Abrir cofre"
	if m.mode == FilePickerSave {
		title = "Salvar cofre"
	}
	titleText, titleWidth := design.RenderDialogTitle(title, "", "", theme)
	fillRight := modalW - 2 - 1 - 1 - titleWidth - 1 // canto + dash + space + title + space
	if fillRight < 1 {
		fillRight = 1
	}

	tl := borderStyle.Render(design.SymCornerTL)
	tr := borderStyle.Render(design.SymCornerTR)
	dashL := borderStyle.Render(design.SymBorderH)
	dashR := borderStyle.Render(strings.Repeat(design.SymBorderH, fillRight))
	return tl + dashL + " " + titleText + " " + dashR + tr
}

// renderPathLine renderiza │ /path/to/dir │
// Trunca com … no início se o caminho for mais longo que innerW - 2.
func (m *FilePickerModal) renderPathLine(innerW int, theme *design.Theme) string {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Focused))
	pathStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))

	path := m.currentPath
	maxPathW := innerW - 2 // 1 espaço cada lado
	if lipgloss.Width(path) > maxPathW {
		// Truncar início com …
		for lipgloss.Width(design.SymEllipsis+path) > maxPathW && len(path) > 0 {
			// Remover pelo prefixo por rune
			_, size := []rune(path)[0], 1
			path = path[size:]
		}
		path = design.SymEllipsis + path
	}
	content := pathStyle.Render(padRight(" "+path, innerW-1) + " ")
	bgStyle := lipgloss.NewStyle().Background(lipgloss.Color(theme.Surface.Raised))
	return borderStyle.Render(design.SymBorderV) +
		bgStyle.Render(content) +
		borderStyle.Render(design.SymBorderV)
}

// renderPanelSeparator renderiza ├─ Estrutura ──┬─ Arquivos ──┤
func (m *FilePickerModal) renderPanelSeparator(innerW, treeW int, theme *design.Theme) string {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Default))
	headerStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(theme.Text.Secondary)).
		Bold(true)

	treeLabel := " " + headerStyle.Render("Estrutura") + " "
	treeLabelW := lipgloss.Width(treeLabel)
	treeDash := treeW - treeLabelW
	if treeDash < 1 {
		treeDash = 1
	}

	filesLabel := " " + headerStyle.Render("Arquivos") + " "
	filesLabelW := lipgloss.Width(filesLabel)
	filesW := innerW - treeW - 1 // -1 para o ┬
	filesDash := filesW - filesLabelW
	if filesDash < 1 {
		filesDash = 1
	}

	jL := borderStyle.Render(design.SymJunctionL)
	jT := borderStyle.Render(design.SymJunctionT)
	jR := borderStyle.Render(design.SymJunctionR)
	dash := func(n int) string {
		return borderStyle.Render(strings.Repeat(design.SymBorderH, n))
	}
	return jL + treeLabel + dash(treeDash) + jT + filesLabel + dash(filesDash) + jR
}

// renderTreeLine renderiza uma linha do painel de árvore.
// absIdx é o índice em visibleNodes; treeW é a largura do painel.
func (m *FilePickerModal) renderTreeLine(absIdx, treeW int, theme *design.Theme) string {
	if absIdx < 0 || absIdx >= len(m.visibleNodes) {
		return strings.Repeat(" ", treeW)
	}
	node := m.visibleNodes[absIdx].node

	// Ícone de pasta
	var icon string
	if node.depth == 0 {
		icon = "" // raiz: sem indicador
	} else if !node.loaded || !node.hasSubdirs {
		if node.loaded {
			icon = design.SymFolderEmpty
		} else {
			icon = design.SymFolderCollapsed // ainda não carregado: assume recolhido
		}
	} else if node.expanded {
		icon = design.SymFolderExpanded
	} else {
		icon = design.SymFolderCollapsed
	}

	iconStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Accent.Secondary))
	indent := strings.Repeat("  ", node.depth)

	var nameStr string
	isCursor := absIdx == m.treeCursor
	isActive := m.focusPanel == 0

	nameStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Primary))
	if isCursor {
		cursorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Accent.Primary)).
			Bold(true)
		if isActive {
			cursorStyle = cursorStyle.Background(lipgloss.Color(theme.Special.Highlight))
		}
		nameStr = cursorStyle.Render(node.name)
	} else {
		nameStr = nameStyle.Render(node.name)
	}

	var line string
	if icon == "" {
		line = indent + nameStr
	} else {
		line = indent + iconStyle.Render(icon) + " " + nameStr
	}
	return padRight(line, treeW)
}

// renderFileLine renderiza uma linha do painel de arquivos.
// absIdx é o índice em m.files; filesW é a largura do painel.
func (m *FilePickerModal) renderFileLine(absIdx, filesW int, theme *design.Theme) string {
	if absIdx < 0 || absIdx >= len(m.files) {
		return strings.Repeat(" ", filesW)
	}

	name := m.files[absIdx]
	info := m.fileInfos[absIdx]
	sizeStr := formatFileSize(info.Size())
	dateStr := m.formatTime(info.ModTime())

	// Larguras fixas
	sizeW := 7  // ex: "25.8 MB"
	dateW := 14 // "dd/mm/aa HH:MM"
	bulletW := 1
	sepW := 1

	// nameW disponível
	nameW := filesW - bulletW - sepW - sizeW - sepW - dateW
	if nameW < 4 {
		// Truncar date
		dateW = 0
		nameW = filesW - bulletW - sepW - sizeW
		if nameW < 4 {
			// Truncar size também
			sizeW = 0
			nameW = filesW - bulletW
		}
	}

	// Truncar nome se necessário
	if lipgloss.Width(name) > nameW && nameW > 1 {
		for lipgloss.Width(name+design.SymEllipsis) > nameW && len(name) > 0 {
			runes := []rune(name)
			name = string(runes[:len(runes)-1])
		}
		name = name + design.SymEllipsis
	}

	isSel := absIdx == m.fileCursor
	bulletStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))
	metaStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))

	var nameRendered string
	if isSel && m.focusPanel == 1 {
		nameRendered = lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Text.Primary)).
			Background(lipgloss.Color(theme.Special.Highlight)).
			Bold(true).
			Render(padRight(name, nameW))
	} else {
		nameRendered = lipgloss.NewStyle().
			Foreground(lipgloss.Color(theme.Text.Primary)).
			Render(padRight(name, nameW))
	}

	bullet := bulletStyle.Render(design.SymLeaf)
	line := bullet + " " + nameRendered
	if sizeW > 0 {
		line += " " + metaStyle.Render(padRight(sizeStr, sizeW))
	}
	if dateW > 0 {
		line += " " + metaStyle.Render(dateStr)
	}
	return padRight(line, filesW)
}

// renderEmptyFilesMessage renderiza a mensagem quando o painel de arquivos está vazio.
func (m *FilePickerModal) renderEmptyFilesMessage(filesW int, theme *design.Theme) string {
	var msg string
	if m.mode == FilePickerOpen {
		msg = "Nenhum cofre neste diretório"
	} else {
		msg = ""
	}
	style := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Text.Secondary))
	return padRight(style.Render(msg), filesW)
}

// renderContentLines monta as visibleH linhas de conteúdo dos dois painéis lado a lado.
func (m *FilePickerModal) renderContentLines(visibleH, treeW, filesW int, theme *design.Theme) []string {
	borderStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(theme.Border.Focused))
	bgStyle := lipgloss.NewStyle().Background(lipgloss.Color(theme.Surface.Raised))
	lBorder := borderStyle.Render(design.SymBorderV)

	lines := make([]string, visibleH)
	for i := 0; i < visibleH; i++ {
		treeIdx := m.treeScroll + i
		fileIdx := m.fileScroll + i

		treePart := bgStyle.Render(m.renderTreeLine(treeIdx, treeW, theme))
		sepChar := renderTreeSepChar(i, m.treeScroll, len(m.visibleNodes), visibleH, theme)

		var filePart string
		if len(m.files) == 0 {
			if i == 1 {
				filePart = bgStyle.Render(m.renderEmptyFilesMessage(filesW, theme))
			} else {
				filePart = bgStyle.Render(strings.Repeat(" ", filesW))
			}
		} else {
			filePart = bgStyle.Render(m.renderFileLine(fileIdx, filesW, theme))
		}
		rBorder := renderFileSepChar(i, m.fileScroll, len(m.files), visibleH, theme)

		lines[i] = lBorder + treePart + sepChar + filePart + rBorder
	}
	return lines
}
