package design

// Símbolos de navegação em árvore.
const (
	// SymFolderCollapsed é o símbolo de pasta recolhida.
	SymFolderCollapsed = "▶" // U+25B6
	// SymFolderExpanded é o símbolo de pasta expandida.
	SymFolderExpanded = "▼" // U+25BC
	// SymFolderEmpty é o símbolo de pasta vazia.
	SymFolderEmpty = "▷" // U+25B7
	// SymLeaf é o símbolo de item folha (segredo).
	SymLeaf = "●" // U+25CF
)

// Símbolos de estado de item.
const (
	// SymFavorite é o símbolo de item favorito.
	SymFavorite = "★" // U+2605
	// SymDeleted é o símbolo de item marcado para exclusão.
	SymDeleted = "✗" // U+2717
	// SymCreated é o símbolo de item recém-criado, ainda não salvo.
	SymCreated = "✦" // U+2726
	// SymModified é o símbolo de item modificado, ainda não salvo.
	SymModified = "✎" // U+270E
)

// Símbolos semânticos — usados na barra de mensagens e em badges de estado.
const (
	// SymSuccess indica operação concluída com êxito.
	SymSuccess = "✓" // U+2713
	// SymInfo indica informação contextual.
	SymInfo = "ℹ" // U+2139
	// SymWarning indica alerta ou aviso.
	SymWarning = "⚠" // U+26A0
	// SymError indica falha ou estado inválido. Distinto de SymDeleted (✗).
	SymError = "✕" // U+2715
)

// Símbolos de elementos de UI.
const (
	// SymRevealable é o indicador de campo com conteúdo revelável.
	SymRevealable = "◉" // U+25C9
	// SymMask é o caractere de máscara para conteúdo sensível.
	// Mesmo glifo que SymBullet (U+2022); papéis semânticos distintos.
	SymMask = "•" // U+2022
	// SymCursor é o cursor em campos de entrada de texto.
	SymCursor = "▌" // U+258C
	// SymScrollUp é o indicador de scroll disponível para cima.
	SymScrollUp = "↑" // U+2191
	// SymScrollDown é o indicador de scroll disponível para baixo.
	SymScrollDown = "↓" // U+2193
	// SymScrollThumb é o thumb que indica a posição atual do scroll.
	SymScrollThumb = "■" // U+25A0
	// SymEllipsis é o caractere de truncamento de texto.
	SymEllipsis = "…" // U+2026
	// SymBullet é o indicador contextual — dirty marker no cabeçalho e marcador de dica.
	// Mesmo glifo que SymMask (U+2022); papéis semânticos distintos.
	SymBullet = "•" // U+2022
	// SymHeaderSep é o separador usado no cabeçalho entre elementos.
	SymHeaderSep = "·" // U+00B7
	// SymTreeConnector é o conector entre a linha selecionada na árvore e o painel de detalhe.
	// Ocupa 2 colunas (Basic Latin '<' + U+2561 '╡').
	SymTreeConnector = "<╡"
)

// Símbolos de estrutura — caracteres Box Drawing usados em bordas e separadores.
const (
	// SymBorderH é o separador horizontal.
	SymBorderH = "─" // U+2500
	// SymBorderV é o separador vertical.
	SymBorderV = "│" // U+2502
	// SymCornerTL é o canto superior esquerdo arredondado.
	SymCornerTL = "╭" // U+256D
	// SymCornerTR é o canto superior direito arredondado.
	SymCornerTR = "╮" // U+256E
	// SymCornerBL é o canto inferior esquerdo arredondado.
	SymCornerBL = "╰" // U+2570
	// SymCornerBR é o canto inferior direito arredondado.
	SymCornerBR = "╯" // U+256F
	// SymJunctionL é o ponto de junção em T voltado para a esquerda.
	SymJunctionL = "├" // U+251C
	// SymJunctionT é o ponto de junção em T voltado para cima.
	SymJunctionT = "┬" // U+252C
	// SymJunctionB é o ponto de junção em T voltado para baixo.
	SymJunctionB = "┴" // U+2534
	// SymJunctionR é o ponto de junção em T voltado para a direita.
	SymJunctionR = "┤" // U+2524
)

// SpinnerFrames são os 4 frames da animação de atividade, em ordem de exibição.
// Nota: caracteres Geometric Shapes podem renderizar em 2 colunas em alguns locales.
// Reserve 2 colunas no layout adjacente para evitar jitter.
var SpinnerFrames = [4]string{"◐", "◓", "◑", "◒"}

// SpinnerFrame retorna o frame do spinner para o índice de animação fornecido.
// O índice é reduzido com frame%4, portanto aceita qualquer valor inteiro não-negativo.
func SpinnerFrame(frame int) string {
	return SpinnerFrames[frame%4]
}
