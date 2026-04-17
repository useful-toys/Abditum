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

func (f fakeFileInfo) Name() string      { return f.name }
func (f fakeFileInfo) Size() int64       { return f.size }
func (f fakeFileInfo) Mode() os.FileMode { return 0o644 }
func (f fakeFileInfo) ModTime() time.Time { return f.modTime }
func (f fakeFileInfo) IsDir() bool       { return f.isDir }
func (f fakeFileInfo) Sys() any          { return nil }

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
		"/": {dir("home")},
		"/home": {dir("usuario")},
		"/home/usuario": {dir("documentos"), dir("downloads"), dir("projetos"), dir("fotos")},
		"/home/usuario/documentos": {dir("contratos"), dir("relatorios")},
		"/home/usuario/documentos/contratos": {dir("2024"), dir("2025")},
		// 2024 retorna permissão negada — ver abaixo
		"/home/usuario/documentos/contratos/2025": {
			file("cofre.abditum", 512_000, date20250401),
		},
		"/home/usuario/documentos/relatorios": {},
		"/home/usuario/downloads": {dir("instaladores"), dir("temporarios")},
		"/home/usuario/downloads/instaladores": {},
		"/home/usuario/downloads/temporarios":  {},
		"/home/usuario/projetos": {dir("abditum"), dir("site")},
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
	lastMethod string
	lastText   string
}

func (s *trackingMsgCtrl) SetHintField(text string)  { s.lastMethod = "HintField"; s.lastText = text }
func (s *trackingMsgCtrl) SetError(text string)      { s.lastMethod = "Error"; s.lastText = text }
func (s *trackingMsgCtrl) SetWarning(text string)    { s.lastMethod = "Warning"; s.lastText = text }
func (s *trackingMsgCtrl) SetBusy(text string)       {}
func (s *trackingMsgCtrl) SetSuccess(text string)    {}
func (s *trackingMsgCtrl) SetInfo(text string)       {}
func (s *trackingMsgCtrl) SetHintUsage(text string)  {}
func (s *trackingMsgCtrl) Clear()                    {}

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
	m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
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
	m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	testdata.TestRenderManaged(t, "file_picker", "save_name_empty", []string{"88x30"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

func TestFilePicker_Render_SaveNameFilled(t *testing.T) {
	m := newSavePicker("meu-cofre")
	m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	m.Update(tea.KeyPressMsg{Code: tea.KeyTab})
	testdata.TestRenderManaged(t, "file_picker", "save_name_filled", []string{"88x30"},
		func(w, h int, theme *design.Theme) string {
			return m.Render(h, w, theme)
		})
}

