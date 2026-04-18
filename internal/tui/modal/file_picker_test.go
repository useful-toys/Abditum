package modal_test

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/tui/testdata"
)

// fakeFileInfo implementa os.FileInfo com campos fixos.
type fakeFileInfo struct {
	name    string
	size    int64
	modTime time.Time
	isDir   bool
}

func (f fakeFileInfo) Name() string       { return f.name }
func (f fakeFileInfo) Size() int64        { return f.size }
func (f fakeFileInfo) Mode() os.FileMode  { return 0o644 }
func (f fakeFileInfo) ModTime() time.Time { return f.modTime }
func (f fakeFileInfo) IsDir() bool        { return f.isDir }
func (f fakeFileInfo) Sys() any           { return nil }

// fakeDirEntry implementa os.DirEntry wrappando fakeFileInfo.
type fakeDirEntry struct{ info fakeFileInfo }

func (e fakeDirEntry) Name() string               { return e.info.name }
func (e fakeDirEntry) IsDir() bool                { return e.info.isDir }
func (e fakeDirEntry) Type() fs.FileMode          { return 0 }
func (e fakeDirEntry) Info() (fs.FileInfo, error) { return e.info, nil }

// dir cria uma fakeDirEntry de diretório.
func dir(name string) os.DirEntry {
	return fakeDirEntry{fakeFileInfo{name: name, isDir: true}}
}

// file cria uma fakeDirEntry de arquivo com tamanho e data fixos.
func file(name string, size int64, modTime time.Time) os.DirEntry {
	return fakeDirEntry{fakeFileInfo{name: name, size: size, modTime: modTime, isDir: false}}
}

// Datas fixas para golden files determinísticos.
var (
	date20250315, _ = time.Parse("2006-01-02 15:04", "2025-03-15 14:32")
	date20250102, _ = time.Parse("2006-01-02 15:04", "2025-01-02 09:15")
	date20250404, _ = time.Parse("2006-01-02 15:04", "2025-04-04 18:47")
	date20250401, _ = time.Parse("2006-01-02 15:04", "2025-04-01 10:00")
)

// makeTestReadDir retorna um readDir para a árvore hipotética de testes.
// Estrutura: / → home → usuario → documentos/downloads/projetos/fotos
// Arquivos .abditum em /home/usuario/projetos/abditum/ e /home/.../contratos/2025/
// Retorna fs.ErrPermission para /home/usuario/documentos/contratos/2024/
func makeTestReadDir() func(string) ([]os.DirEntry, error) {
	table := map[string][]os.DirEntry{
		"/":                                  {dir("home")},
		"/home":                              {dir("usuario")},
		"/home/usuario":                      {dir("documentos"), dir("downloads"), dir("projetos"), dir("fotos")},
		"/home/usuario/documentos":           {dir("contratos"), dir("relatorios")},
		"/home/usuario/documentos/contratos": {dir("2024"), dir("2025")},
		// 2024 retorna permissão negada — ver abaixo
		"/home/usuario/documentos/contratos/2025": {
			file("cofre.abditum", 512_000, date20250401),
		},
		"/home/usuario/documentos/relatorios":  {},
		"/home/usuario/downloads":              {dir("instaladores"), dir("temporarios")},
		"/home/usuario/downloads/instaladores": {},
		"/home/usuario/downloads/temporarios":  {},
		"/home/usuario/projetos":               {dir("abditum"), dir("site")},
		"/home/usuario/projetos/abditum": {
			dir("docs"),
			dir("src"),
			file("database.abditum", 25_800_000, date20250315),
			file("config.abditum", 1_229, date20250102),
			file("backup.abditum", 18_400_000, date20250404),
		},
		"/home/usuario/projetos/abditum/docs": {},
		"/home/usuario/projetos/abditum/src":  {},
		"/home/usuario/projetos/site":         {},
		"/home/usuario/fotos":                 {},
	}
	return func(path string) ([]os.DirEntry, error) {
		// Normalizar separadores para compatibilidade Unix/Windows nos testes.
		path = filepath.ToSlash(filepath.Clean(path))
		if path == "/home/usuario/documentos/contratos/2024" {
			return nil, fs.ErrPermission
		}
		entries, ok := table[path]
		if !ok {
			return nil, fs.ErrNotExist
		}
		return entries, nil
	}
}

// makeTestReadDirManyFiles retorna readDir com 12 arquivos .abditum em /home/usuario/projetos/abditum/
// Usado nos golden files de scroll de arquivos.
func makeTestReadDirManyFiles() func(string) ([]os.DirEntry, error) {
	base := makeTestReadDir()
	manyFiles := []os.DirEntry{
		dir("docs"), dir("src"),
		file("arquivo01.abditum", 1_000, date20250315),
		file("arquivo02.abditum", 2_000, date20250315),
		file("arquivo03.abditum", 3_000, date20250315),
		file("arquivo04.abditum", 4_000, date20250315),
		file("arquivo05.abditum", 5_000, date20250315),
		file("arquivo06.abditum", 6_000, date20250315),
		file("arquivo07.abditum", 7_000, date20250315),
		file("arquivo08.abditum", 8_000, date20250315),
		file("arquivo09.abditum", 9_000, date20250315),
		file("arquivo10.abditum", 10_000, date20250315),
		file("arquivo11.abditum", 11_000, date20250315),
		file("arquivo12.abditum", 12_000, date20250315),
	}
	return func(path string) ([]os.DirEntry, error) {
		// Normalizar separadores para compatibilidade Unix/Windows nos testes.
		path = filepath.ToSlash(filepath.Clean(path))
		if path == "/home/usuario/projetos/abditum" {
			return manyFiles, nil
		}
		return base(path)
	}
}

// fixedTimeFmt retorna uma função de formatação de tempo determinística.
func fixedTimeFmt(tm time.Time) string { return tm.Format("02/01/06 15:04") }

// newOpenPicker cria FilePicker Open com filesystem fictício.
func newOpenPicker() *modal.FilePickerModal {
	m := modal.NewFilePicker(modal.FilePickerOptions{
		Mode:       modal.FilePickerOpen,
		Extension:  ".abditum",
		InitialDir: "/home/usuario/projetos/abditum",
		OnResult:   func(string) tea.Cmd { return nil },
	})
	m.SetReadDirForTest(makeTestReadDir())
	m.SetTimeFmtForTest(fixedTimeFmt)
	m.RebuildForTest("/home/usuario/projetos/abditum")
	return m
}

// newSavePicker cria FilePicker Save com filesystem fictício.
func newSavePicker(suggested string) *modal.FilePickerModal {
	m := modal.NewFilePicker(modal.FilePickerOptions{
		Mode:       modal.FilePickerSave,
		Extension:  ".abditum",
		InitialDir: "/home/usuario/projetos/abditum",
		Suggested:  suggested,
		OnResult:   func(string) tea.Cmd { return nil },
	})
	m.SetReadDirForTest(makeTestReadDir())
	m.SetTimeFmtForTest(fixedTimeFmt)
	m.RebuildForTest("/home/usuario/projetos/abditum")
	return m
}

// keyPress cria um tea.KeyMsg para tecla com texto.
func keyPress(code tea.Key) tea.KeyMsg {
	return tea.KeyPressMsg(code)
}

// trackingMsgCtrl implementa tui.MessageController registrando o último método
// e texto chamados — usado em testes de comportamento do FilePicker.
type trackingMsgCtrl struct {
	lastMethod  string
	lastText    string
	errorCalled bool
}

func (s *trackingMsgCtrl) SetHintField(text string) { s.lastMethod = "HintField"; s.lastText = text }
func (s *trackingMsgCtrl) SetError(text string) {
	s.lastMethod = "Error"
	s.lastText = text
	s.errorCalled = true
}
func (s *trackingMsgCtrl) SetWarning(text string)   { s.lastMethod = "Warning"; s.lastText = text }
func (s *trackingMsgCtrl) SetBusy(text string)      {}
func (s *trackingMsgCtrl) SetSuccess(text string)   {}
func (s *trackingMsgCtrl) SetInfo(text string)      {}
func (s *trackingMsgCtrl) SetHintUsage(text string) {}
func (s *trackingMsgCtrl) Clear()                   {}

var _ tui.MessageController = (*trackingMsgCtrl)(nil)

func TestFilePicker_BuildTreeChain_ExpandsAncestors(t *testing.T) {
	m := newOpenPicker()
	// treeCursor deve apontar para /home/usuario/projetos/abditum
	if m.TreeCursorPath() != "/home/usuario/projetos/abditum" {
		t.Errorf("treeCursor path = %q, want /home/usuario/projetos/abditum", m.TreeCursorPath())
	}
}

func TestFormatFileSize(t *testing.T) {
	cases := []struct {
		bytes int64
		want  string
	}{
		{25_800_000, "24.6 MB"},
		{1_229, "1.2 KB"},
		{2_000_000_000, "1.9 GB"},
		{1_024, "1.0 KB"},
	}
	for _, c := range cases {
		got := modal.FormatFileSizeForTest(c.bytes)
		if got != c.want {
			t.Errorf("formatFileSize(%d) = %q, want %q", c.bytes, got, c.want)
		}
	}
}

func TestFilePicker_Render_OpenTreeInitial(t *testing.T) {
	m := newOpenPicker()
	testdata.TestRenderManaged(t, "file_picker", "open_tree_initial", []string{"88x30"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestFilePicker_Render_OpenFilesNoScroll(t *testing.T) {
	m := newOpenPicker()
	// Mover foco para arquivos — Tab é no-op com HandleKey stub, mas o teste fica pronto para quando HandleKey for implementado
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	testdata.TestRenderManaged(t, "file_picker", "open_files_noscroll", []string{"88x30"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestFilePicker_Render_OpenEmptyDir(t *testing.T) {
	m := modal.NewFilePicker(modal.FilePickerOptions{
		Mode:       modal.FilePickerOpen,
		Extension:  ".abditum",
		InitialDir: "/home/usuario/downloads/temporarios",
		OnResult:   func(string) tea.Cmd { return nil },
	})
	m.SetReadDirForTest(makeTestReadDir())
	m.SetTimeFmtForTest(fixedTimeFmt)
	m.RebuildForTest("/home/usuario/downloads/temporarios")
	testdata.TestRenderManaged(t, "file_picker", "open_empty_dir", []string{"88x30"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestFilePicker_Render_SaveNameEmpty(t *testing.T) {
	m := newSavePicker("")
	// Tab x2: no-op com stub, mas teste pronto para quando HandleKey for implementado
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	testdata.TestRenderManaged(t, "file_picker", "save_name_empty", []string{"88x30"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestFilePicker_Render_SaveNameFilled(t *testing.T) {
	m := newSavePicker("meu-cofre")
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	testdata.TestRenderManaged(t, "file_picker", "save_name_filled", []string{"88x30"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

// newOpenPickerManyFiles cria FilePicker Open com 12 arquivos em /home/usuario/projetos/abditum/
// para golden files de scroll.
func newOpenPickerManyFiles() *modal.FilePickerModal {
	m := modal.NewFilePicker(modal.FilePickerOptions{
		Mode:       modal.FilePickerOpen,
		Extension:  ".abditum",
		InitialDir: "/home/usuario/projetos/abditum",
		OnResult:   func(string) tea.Cmd { return nil },
	})
	m.SetReadDirForTest(makeTestReadDirManyFiles())
	m.SetTimeFmtForTest(fixedTimeFmt)
	m.RebuildForTest("/home/usuario/projetos/abditum")
	return m
}

func TestFilePicker_Render_OpenFilesScrollTop(t *testing.T) {
	m := newOpenPickerManyFiles()
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab}) // foco para arquivos (no-op with stub)
	// fileScroll=0, cursor=0 — estado inicial já é scroll top
	testdata.TestRenderManaged(t, "file_picker", "open_files_scroll_top", []string{"88x30"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestFilePicker_Render_OpenFilesScrollMid(t *testing.T) {
	m := newOpenPickerManyFiles()
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	// Navegar para o meio da lista (6 downs — no-op with stub)
	for i := 0; i < 6; i++ {
		m.HandleKey(tea.KeyPressMsg{Code: tea.KeyDown})
	}
	testdata.TestRenderManaged(t, "file_picker", "open_files_scroll_mid", []string{"88x30"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestFilePicker_Render_OpenFilesScrollEnd(t *testing.T) {
	m := newOpenPickerManyFiles()
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	// End key (no-op with stub)
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnd})
	testdata.TestRenderManaged(t, "file_picker", "open_files_scroll_end", []string{"88x30"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestFilePicker_Render_OpenTreeScrollTop(t *testing.T) {
	m := newOpenPicker()
	// Expandir toda a árvore (varios Rights) — no-op with HandleKey stub
	for i := 0; i < 20; i++ {
		m.HandleKey(tea.KeyPressMsg{Code: tea.KeyRight})
		m.HandleKey(tea.KeyPressMsg{Code: tea.KeyDown})
	}
	// Voltar ao topo — no-op with HandleKey stub
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyHome})
	testdata.TestRenderManaged(t, "file_picker", "open_tree_scroll_top", []string{"88x14"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestFilePicker_Render_OpenTreeScrollMid(t *testing.T) {
	m := newOpenPicker()
	for i := 0; i < 20; i++ {
		m.HandleKey(tea.KeyPressMsg{Code: tea.KeyRight})
		m.HandleKey(tea.KeyPressMsg{Code: tea.KeyDown})
	}
	// Ir para o meio — no-op with HandleKey stub
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyHome})
	for i := 0; i < 7; i++ {
		m.HandleKey(tea.KeyPressMsg{Code: tea.KeyDown})
	}
	testdata.TestRenderManaged(t, "file_picker", "open_tree_scroll_mid", []string{"88x14"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestFilePicker_Render_OpenTreeScrollEnd(t *testing.T) {
	m := newOpenPicker()
	for i := 0; i < 20; i++ {
		m.HandleKey(tea.KeyPressMsg{Code: tea.KeyRight})
		m.HandleKey(tea.KeyPressMsg{Code: tea.KeyDown})
	}
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnd})
	testdata.TestRenderManaged(t, "file_picker", "open_tree_scroll_end", []string{"88x14"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

// --- Open behavior ---

func TestFilePicker_Open_Enter_OnFile(t *testing.T) {
	var resultPath string
	m := modal.NewFilePicker(modal.FilePickerOptions{
		Mode:       modal.FilePickerOpen,
		Extension:  ".abditum",
		InitialDir: "/home/usuario/projetos/abditum",
		OnResult:   func(p string) tea.Cmd { resultPath = p; return nil },
	})
	m.SetReadDirForTest(makeTestReadDir())
	m.SetTimeFmtForTest(fixedTimeFmt)
	m.RebuildForTest("/home/usuario/projetos/abditum")

	m.Render(30, 88, design.TokyoNight)
	m.HandleKey(tea.KeyPressMsg{}) // prime hintEmitted

	// Mover foco para arquivos e confirmar
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnter})

	want := "/home/usuario/projetos/abditum/database.abditum"
	if filepath.ToSlash(resultPath) != want {
		t.Errorf("OnResult path = %q, want %q", resultPath, want)
	}
}

func TestFilePicker_Open_Esc_Cancels(t *testing.T) {
	var resultPath = "not-called"
	m := modal.NewFilePicker(modal.FilePickerOptions{
		Mode:       modal.FilePickerOpen,
		Extension:  ".abditum",
		InitialDir: "/home/usuario/projetos/abditum",
		OnResult:   func(p string) tea.Cmd { resultPath = p; return nil },
	})
	m.SetReadDirForTest(makeTestReadDir())
	m.SetTimeFmtForTest(fixedTimeFmt)
	m.RebuildForTest("/home/usuario/projetos/abditum")
	m.Render(30, 88, design.TokyoNight)
	m.HandleKey(tea.KeyPressMsg{}) // prime hintEmitted
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEscape})

	if resultPath != "" {
		t.Errorf("OnResult path = %q, want empty string (cancelled)", resultPath)
	}
}

// --- Save behavior ---

func TestFilePicker_Save_Enter_WithName(t *testing.T) {
	var resultPath string
	m := modal.NewFilePicker(modal.FilePickerOptions{
		Mode:       modal.FilePickerSave,
		Extension:  ".abditum",
		InitialDir: "/home/usuario/projetos/abditum",
		OnResult:   func(p string) tea.Cmd { resultPath = p; return nil },
	})
	m.SetReadDirForTest(makeTestReadDir())
	m.SetTimeFmtForTest(fixedTimeFmt)
	m.RebuildForTest("/home/usuario/projetos/abditum")
	m.Render(30, 88, design.TokyoNight)
	m.HandleKey(tea.KeyPressMsg{}) // prime hintEmitted

	// Tab x2 para campo
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	// Digitar nome
	for _, r := range "meu-cofre" {
		m.HandleKey(tea.KeyPressMsg{Code: rune(r), Text: string(r)})
	}
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnter})

	want := "/home/usuario/projetos/abditum/meu-cofre.abditum"
	if filepath.ToSlash(resultPath) != want {
		t.Errorf("OnResult = %q, want %q", resultPath, want)
	}
}

func TestFilePicker_Save_Enter_WithExtension(t *testing.T) {
	var resultPath string
	m := modal.NewFilePicker(modal.FilePickerOptions{
		Mode:       modal.FilePickerSave,
		Extension:  ".abditum",
		InitialDir: "/home/usuario/projetos/abditum",
		OnResult:   func(p string) tea.Cmd { resultPath = p; return nil },
	})
	m.SetReadDirForTest(makeTestReadDir())
	m.SetTimeFmtForTest(fixedTimeFmt)
	m.RebuildForTest("/home/usuario/projetos/abditum")
	m.Render(30, 88, design.TokyoNight)
	m.HandleKey(tea.KeyPressMsg{}) // prime hintEmitted
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	for _, r := range "meu-cofre.abditum" {
		m.HandleKey(tea.KeyPressMsg{Code: rune(r), Text: string(r)})
	}
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnter})

	want := "/home/usuario/projetos/abditum/meu-cofre.abditum"
	if filepath.ToSlash(resultPath) != want {
		t.Errorf("OnResult = %q, want %q (no duplicate extension)", resultPath, want)
	}
}

func TestFilePicker_Save_Enter_EmptyField(t *testing.T) {
	called := false
	m := modal.NewFilePicker(modal.FilePickerOptions{
		Mode:       modal.FilePickerSave,
		Extension:  ".abditum",
		InitialDir: "/home/usuario/projetos/abditum",
		OnResult:   func(p string) tea.Cmd { called = true; return nil },
	})
	m.SetReadDirForTest(makeTestReadDir())
	m.SetTimeFmtForTest(fixedTimeFmt)
	m.RebuildForTest("/home/usuario/projetos/abditum")
	m.Render(30, 88, design.TokyoNight)
	m.HandleKey(tea.KeyPressMsg{}) // prime hintEmitted
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	// Enter sem digitar
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnter})
	if called {
		t.Error("OnResult should not be called when name field is empty")
	}
}

// --- Tree navigation ---

func TestFilePicker_Tree_Left_CollapsesNode(t *testing.T) {
	m := newOpenPicker()
	m.Render(30, 88, design.TokyoNight)
	m.HandleKey(tea.KeyPressMsg{}) // prime hintEmitted

	if m.TreeCursorPath() != "/home/usuario/projetos/abditum" {
		t.Fatalf("expected cursor at abditum, got %q", m.TreeCursorPath())
	}

	// Expandir abditum para ter filhos visíveis, depois recolher
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyRight})
	hasDocs := false
	for _, p := range m.VisibleNodePaths() {
		if p == "/home/usuario/projetos/abditum/docs" {
			hasDocs = true
		}
	}
	if !hasDocs {
		t.Fatal("docs not visible after expanding abditum — precondition failed")
	}

	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyLeft})

	for _, p := range m.VisibleNodePaths() {
		if p == "/home/usuario/projetos/abditum/docs" {
			t.Error("docs still visible after collapsing abditum")
		}
	}
}

func TestFilePicker_Tree_Left_AtRoot(t *testing.T) {
	m := newOpenPicker()
	m.Render(30, 88, design.TokyoNight)
	m.HandleKey(tea.KeyPressMsg{}) // prime hintEmitted
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyHome})
	before := m.TreeCursorPath()
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyLeft})
	after := m.TreeCursorPath()
	if before != after {
		t.Errorf("Left at root changed cursor: %q → %q", before, after)
	}
}

func TestFilePicker_Tree_Right_EmptyFolder(t *testing.T) {
	m := newOpenPicker()
	m.Render(30, 88, design.TokyoNight)
	m.HandleKey(tea.KeyPressMsg{}) // prime hintEmitted
	// Navegar para /home/usuario/fotos (sem subdiretórios)
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyEnd})
	before := m.VisibleNodePaths()
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyRight})
	after := m.VisibleNodePaths()
	if len(before) != len(after) {
		t.Error("Right on empty folder changed visible nodes")
	}
}

// --- Tab cycling, invalid chars, permission ---

func TestFilePicker_Tab_Cycles_Open(t *testing.T) {
	m := newOpenPicker()
	m.Render(30, 88, design.TokyoNight)
	m.HandleKey(tea.KeyPressMsg{}) // prime hintEmitted
	if m.FocusPanel() != 0 {
		t.Fatal("initial focus should be 0 (tree)")
	}
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.FocusPanel() != 1 {
		t.Errorf("after 1 Tab: focus = %d, want 1", m.FocusPanel())
	}
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.FocusPanel() != 0 {
		t.Errorf("after 2 Tab: focus = %d, want 0", m.FocusPanel())
	}
}

func TestFilePicker_Tab_Cycles_Save(t *testing.T) {
	m := newSavePicker("")
	m.Render(30, 88, design.TokyoNight)
	m.HandleKey(tea.KeyPressMsg{}) // prime hintEmitted
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.FocusPanel() != 1 {
		t.Errorf("after 1 Tab: focus = %d, want 1", m.FocusPanel())
	}
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.FocusPanel() != 2 {
		t.Errorf("after 2 Tab: focus = %d, want 2", m.FocusPanel())
	}
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	if m.FocusPanel() != 0 {
		t.Errorf("after 3 Tab: focus = %d, want 0", m.FocusPanel())
	}
}

func TestFilePicker_InvalidChars_Blocked(t *testing.T) {
	m := newSavePicker("")
	m.Render(30, 88, design.TokyoNight)
	m.HandleKey(tea.KeyPressMsg{}) // prime hintEmitted
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab}) // foco no campo
	m.HandleKey(tea.KeyPressMsg{Code: '/', Text: "/"})
	if m.NameFieldValue() != "" {
		t.Errorf("invalid char '/' should be blocked, got %q", m.NameFieldValue())
	}
}

func TestFilePicker_PermissionDenied(t *testing.T) {
	// NOTE: Use trackingMsgCtrl (NOT stubMessageController — that's the no-op from password tests)
	stub := &trackingMsgCtrl{}
	m := modal.NewFilePicker(modal.FilePickerOptions{
		Mode:       modal.FilePickerOpen,
		Extension:  ".abditum",
		InitialDir: "/home/usuario/documentos/contratos",
		Messages:   stub,
		OnResult:   func(string) tea.Cmd { return nil },
	})
	m.SetReadDirForTest(makeTestReadDir())
	m.SetTimeFmtForTest(fixedTimeFmt)
	m.RebuildForTest("/home/usuario/documentos/contratos")
	m.Render(30, 88, design.TokyoNight)
	m.HandleKey(tea.KeyPressMsg{}) // prime hintEmitted

	// treeCursor começa em contratos; expandir para ver 2024 e 2025
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyRight})

	// Navegar para 2024 e tentar expandir
	for i := 0; i < 10; i++ {
		m.HandleKey(tea.KeyPressMsg{Code: tea.KeyDown})
		if m.TreeCursorPath() == "/home/usuario/documentos/contratos/2024" {
			break
		}
	}
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyRight})

	// SetError deve ter sido chamado (mesmo que SetHintField sobreescreva lastMethod)
	if !stub.errorCalled {
		t.Errorf("expected SetError to be called for permission denied, lastMethod=%q", stub.lastMethod)
	}
}

// --- Cursor ---

func TestFilePicker_Cursor_NilWhenNotSave(t *testing.T) {
	m := newOpenPicker()
	m.Render(30, 88, design.TokyoNight)
	if c := m.Cursor(0, 0); c != nil {
		t.Errorf("Cursor() = %v, want nil for Open mode", c)
	}
}

func TestFilePicker_Cursor_NilWhenNotFocused(t *testing.T) {
	m := newSavePicker("")
	m.Render(30, 88, design.TokyoNight)
	// foco=0 (árvore)
	if c := m.Cursor(0, 0); c != nil {
		t.Errorf("Cursor() = %v, want nil when focusPanel != 2", c)
	}
}

func TestFilePicker_Cursor_Position(t *testing.T) {
	m := newSavePicker("abc")
	m.Render(30, 88, design.TokyoNight)
	m.HandleKey(tea.KeyPressMsg{}) // prime hintEmitted
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab})
	m.HandleKey(tea.KeyPressMsg{Code: tea.KeyTab}) // foco no campo

	c := m.Cursor(5, 2)
	if c == nil {
		t.Fatal("Cursor() = nil, want non-nil when focusPanel == 2")
	}
	// Y = topY(5) + 4 + visibleH; visibleH for h=30, Save: (30*8/10) - 5 = 24-5=19; Y = 5+4+19 = 28
	// X = leftX(2) + 1 + 1 + 9 + nameField.Position()
	// "abc" suggested → position = 3
	wantY := 5 + 4 + 19
	wantX := 2 + 1 + 1 + 9 + 3
	if c.Position.Y != wantY {
		t.Errorf("Cursor.Y = %d, want %d", c.Position.Y, wantY)
	}
	if c.Position.X != wantX {
		t.Errorf("Cursor.X = %d, want %d", c.Position.X, wantX)
	}
}
