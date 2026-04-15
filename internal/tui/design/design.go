package design

// Theme agrupa todos os tokens visuais que definem a identidade visual da aplicação.
// Passe um valor de Theme para componentes de UI para garantir consistência visual.
type Theme struct {
	// Name identifica o tema, ex: "Tokyo Night".
	Name string
	// Surface define as cores de fundo para as camadas da interface.
	Surface SurfaceTokens
	// Text define as cores de texto para diferentes contextos.
	Text TextTokens
	// Accent define as cores de destaque para elementos interativos.
	Accent AccentTokens
	// Border define as cores de borda nos estados padrão e com foco.
	Border BorderTokens
	// Semantic define cores com significado semântico (sucesso, erro, aviso, etc).
	Semantic SemanticTokens
	// Special define tokens visuais de uso específico.
	Special SpecialTokens
	// LogoGradient são as 5 cores para as linhas do logo ASCII, de cima para baixo.
	LogoGradient [5]string
}

// SurfaceTokens define as cores de fundo para as camadas da interface.
type SurfaceTokens struct {
	// Base é a cor mais escura, usada como pano de fundo da tela principal.
	Base string
	// Raised é a cor de painéis e cards elevados acima da base.
	Raised string
	// Input é a cor de fundo de campos de entrada de texto.
	Input string
}

// TextTokens define as cores de texto para diferentes níveis de hierarquia visual.
type TextTokens struct {
	// Primary é a cor padrão para conteúdo de alta legibilidade.
	Primary string
	// Secondary é usada em textos de suporte e metadados.
	Secondary string
	// Disabled é usada em textos de elementos desabilitados.
	Disabled string
	// Link é a cor para textos clicáveis ou navegáveis.
	Link string
}

// AccentTokens define as cores de destaque do tema.
type AccentTokens struct {
	// Primary é a cor de ação principal, usada em botões e seleções ativas.
	Primary string
	// Secondary é uma cor complementar para destaque alternativo.
	Secondary string
}

// BorderTokens define as cores de borda para os estados de foco.
type BorderTokens struct {
	// Default é a cor de borda de elementos sem foco.
	Default string
	// Focused é a cor de borda do elemento atualmente focado.
	Focused string
}

// SemanticTokens define cores associadas a significados do sistema.
type SemanticTokens struct {
	// Success indica uma operação concluída com êxito.
	Success string
	// Warning indica uma situação que requer atenção.
	Warning string
	// Error indica uma falha ou estado inválido.
	Error string
	// Info é usada para mensagens informativas neutras.
	Info string
	// Off é usada em estados desativados ou indisponíveis.
	Off string
}

// SpecialTokens agrupa tokens visuais de uso específico.
type SpecialTokens struct {
	// Muted é para texto ou ícones com menor ênfase visual.
	Muted string
	// Highlight é a cor de fundo para itens selecionados em destaque.
	Highlight string
	// Match é usada para realçar trechos correspondentes em buscas.
	Match string
}

// TokyoNight é o tema padrão da aplicação, baseado na paleta Tokyo Night.
var TokyoNight = &Theme{
	Name:         "Tokyo Night",
	Surface:      SurfaceTokens{"#1a1b26", "#24283b", "#1e1f2e"},
	Text:         TextTokens{"#a9b1d6", "#565f89", "#3b4261", "#7aa2f7"},
	Accent:       AccentTokens{"#7aa2f7", "#bb9af7"},
	Border:       BorderTokens{"#414868", "#7aa2f7"},
	Semantic:     SemanticTokens{"#9ece6a", "#e0af68", "#f7768e", "#7dcfff", "#737aa2"},
	Special:      SpecialTokens{"#8690b5", "#283457", "#f7c67a"},
	LogoGradient: [5]string{"#9d7cd8", "#89ddff", "#7aa2f7", "#7dcfff", "#bb9af7"},
}

// Cyberpunk é um tema alternativo com paleta neon sobre fundo escuro.
var Cyberpunk = &Theme{
	Name:         "Cyberpunk",
	Surface:      SurfaceTokens{"#0a0a1a", "#1a1a2e", "#0e0e22"},
	Text:         TextTokens{"#e0e0ff", "#8888aa", "#444466", "#ff2975"},
	Accent:       AccentTokens{"#ff2975", "#00fff5"},
	Border:       BorderTokens{"#3a3a5c", "#ff2975"},
	Semantic:     SemanticTokens{"#05ffa1", "#ffe900", "#ff3860", "#00b4d8", "#9999cc"},
	Special:      SpecialTokens{"#666688", "#2a1533", "#ffc107"},
	LogoGradient: [5]string{"#ff2975", "#b026ff", "#00fff5", "#05ffa1", "#ff2975"},
}

// MinWidth é a largura mínima suportada do terminal em colunas (padrão POSIX).
const MinWidth = 80

// MinHeight é a altura mínima suportada do terminal em linhas (padrão POSIX).
const MinHeight = 24

// PanelTreeRatio é a proporção nominal do painel de navegação esquerdo em modos de dois painéis.
// O painel direito (detalhe) ocupa o restante. A implementação pode ajustar ±5%.
const PanelTreeRatio = 0.35

// HeaderHeight é a altura em linhas da região de cabeçalho da tela principal.
const HeaderHeight = 2

// MessageHeight é a altura em linhas da barra de mensagens de status.
const MessageHeight = 1

// ActionHeight é a altura em linhas da barra de ações do contexto atual.
const ActionHeight = 1
