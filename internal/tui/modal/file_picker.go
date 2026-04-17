package modal

import (
	"os"
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

// buildTreeChain constrói a árvore lazy-loaded de diretórios.
func (m *FilePickerModal) buildTreeChain(dir string) *treeNode {
	return &treeNode{path: dir, name: dir, depth: 0}
}

// buildVisibleNodes achatada a árvore expandida em visibleNodes.
func (m *FilePickerModal) buildVisibleNodes() {}

// loadFiles carrega a lista de arquivos do currentPath.
func (m *FilePickerModal) loadFiles(dir string) {}

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
