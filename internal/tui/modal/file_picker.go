package modal

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/bubbles/v2/textinput"
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
