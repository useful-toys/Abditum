package design

// WorkArea representa qual área de trabalho está ativa na tela principal.
// É usada por RootModel para decidir qual ChildView exibir e pelo HeaderView
// para renderizar a aba ativa no cabeçalho.
type WorkArea int

const (
	// WorkAreaWelcome exibe a tela de boas-vindas, para usuários sem cofre aberto.
	WorkAreaWelcome WorkArea = iota
	// WorkAreaSettings exibe as configurações da aplicação.
	WorkAreaSettings
	// WorkAreaVault exibe a área de gerenciamento do cofre de segredos.
	WorkAreaVault
	// WorkAreaTemplates exibe a área de gerenciamento de templates de segredos.
	WorkAreaTemplates
)
