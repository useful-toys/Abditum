package operation

import (
	"errors"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/vault"
)

// abrirEstado representa em qual etapa do fluxo de abertura estamos.
type abrirEstado int

const (
	abrindoInformandoCaminho abrirEstado = iota // FilePicker modo Open
	abrindoInformandoSenha                      // PasswordEntryModal
	abrindoAbrindo                              // IO assíncrono
)

// abrirAvancaMsg é a mensagem interna de transição de estado.
type abrirAvancaMsg struct {
	estado abrirEstado
}

// abrirCofreResultMsg carrega o resultado da tentativa de abrir o cofre.
type abrirCofreResultMsg struct {
	manager *vault.Manager
	err     error
}

// AbrirCofreOperation implementa o Fluxo 1 (Abrir Cofre Existente).
//
// Se caminhoInicial != "", valida o cabeçalho do arquivo e vai diretamente
// para a entrada de senha. Caso contrário, executa o fluxo GUI completo.
type AbrirCofreOperation struct {
	notifier tui.MessageController
	saver    vaultSaver
	guard    *guardCofreAlterado // não-nil apenas no fluxo GUI com cofre alterado
	caminho  string              // caminho do cofre selecionado
	senha    []byte              // senha informada pelo usuário
}

// NewAbrirCofreOperation cria a operation com o vault.Manager concreto.
// manager pode ser nil quando nenhum cofre está carregado.
// caminhoInicial pode ser "" para fluxo GUI completo, ou um caminho CLI.
func NewAbrirCofreOperation(
	notifier tui.MessageController,
	manager *vault.Manager,
	caminhoInicial string,
) *AbrirCofreOperation {
	var saver vaultSaver
	if manager != nil {
		saver = manager
	}
	return &AbrirCofreOperation{
		notifier: notifier,
		saver:    saver,
		caminho:  caminhoInicial,
	}
}

// newAbrirCofreOperationFromSaver é usada nos testes (mesmo pacote).
func newAbrirCofreOperationFromSaver(
	notifier tui.MessageController,
	saver vaultSaver,
	caminhoInicial string,
) *AbrirCofreOperation {
	return &AbrirCofreOperation{
		notifier: notifier,
		saver:    saver,
		caminho:  caminhoInicial,
	}
}

// Init inicia o fluxo de abertura.
func (a *AbrirCofreOperation) Init() tea.Cmd {
	if a.caminho != "" {
		// Entrada via CLI: validar cabeçalho antes de pedir senha
		if err := storage.ValidateHeader(a.caminho); err != nil {
			a.notifier.SetError(erroDeAberturaCategoria(err))
			return tui.OperationCompleted()
		}
		return a.abrirModalSenha()
	}
	// Fluxo GUI completo: verificar cofre alterado primeiro
	a.guard = novoGuardCofreAlterado(
		a.notifier,
		a.saver,
		func() tea.Cmd {
			return func() tea.Msg { return abrirAvancaMsg{estado: abrindoInformandoCaminho} }
		},
		func() tea.Cmd { return tui.OperationCompleted() },
	)
	return a.guard.Init()
}

// Update trata as mensagens internas da AbrirCofreOperation.
func (a *AbrirCofreOperation) Update(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	// Mensagens do guard — delegar se guard ativo
	case guardSaveMsg, guardSaveResultMsg, guardDiscardMsg, guardCancelMsg:
		if a.guard != nil {
			return a.guard.Update(msg)
		}
		return nil

	case abrirAvancaMsg:
		switch m.estado {
		case abrindoInformandoCaminho:
			return a.abrirFilePicker()
		case abrindoInformandoSenha:
			return a.abrirModalSenha()
		case abrindoAbrindo:
			a.notifier.SetBusy("Abrindo cofre...")
			caminho := a.caminho
			senha := a.senha
			// Carregar o cofre de forma assíncrona para não bloquear a UI
			return func() tea.Msg {
				repo := storage.NewFileRepositoryForOpen(caminho, senha)
				cofre, err := repo.Carregar()
				if err != nil {
					return abrirCofreResultMsg{err: err}
				}
				manager := vault.NewManager(cofre, repo)
				return abrirCofreResultMsg{manager: manager}
			}
		}

	case abrirCofreResultMsg:
		if m.err != nil {
			a.notifier.SetError(erroDeAberturaCategoria(m.err))
			// Senha incorreta: reabrir modal de senha para nova tentativa
			if isErrAutenticacao(m.err) {
				return a.abrirModalSenha()
			}
			// Arquivo corrompido ou inválido: voltar ao seletor de arquivo
			return a.abrirFilePicker()
		}
		a.notifier.Clear()
		return func() tea.Msg { return tui.VaultOpenedMsg{Manager: m.manager} }
	}

	return nil
}

// abrirFilePicker abre o FilePicker no modo Open com extensão .abditum.
// A construção do modal é diferida para dentro do Cmd, evitando efeitos colaterais
// (como emissão de hints) durante o processamento de mensagens no Update.
func (a *AbrirCofreOperation) abrirFilePicker() tea.Cmd {
	return func() tea.Msg {
		return tui.OpenModalMsg{Modal: modal.NewFilePicker(modal.FilePickerOptions{
			Mode:      modal.FilePickerOpen,
			Extension: ".abditum",
			Messages:  a.notifier,
			OnResult: func(path string) tea.Cmd {
				if path == "" {
					return tui.OperationCompleted()
				}
				// Validar o cabeçalho antes de pedir a senha — rápido, sem criptografia
				if err := storage.ValidateHeader(path); err != nil {
					a.notifier.SetError(erroDeAberturaCategoria(err))
					return a.abrirFilePicker()
				}
				a.caminho = path
				return func() tea.Msg { return abrirAvancaMsg{estado: abrindoInformandoSenha} }
			},
		})}
	}
}

// abrirModalSenha abre o PasswordEntryModal para coleta da senha mestra.
// A construção do modal é diferida para dentro do Cmd, evitando que SetHintField
// sobrescreva mensagens de erro emitidas antes de retornar o Cmd.
func (a *AbrirCofreOperation) abrirModalSenha() tea.Cmd {
	return func() tea.Msg {
		return tui.OpenModalMsg{Modal: modal.NewPasswordEntryModal(
			a.notifier,
			func(password []byte) tea.Cmd {
				a.senha = password
				return tea.Batch(tui.CloseModal(), func() tea.Msg {
					return abrirAvancaMsg{estado: abrindoAbrindo}
				})
			},
			func() tea.Cmd {
				// Cancelar senha: voltar ao seletor de arquivo
				return tea.Batch(tui.CloseModal(), func() tea.Msg {
					return abrirAvancaMsg{estado: abrindoInformandoCaminho}
				})
			},
		)}
	}
}

// erroDeAberturaCategoria converte erros técnicos em mensagens amigáveis ao usuário.
// Evita expor detalhes de implementação ou criptografia na interface.
func erroDeAberturaCategoria(err error) string {
	if errors.Is(err, crypto.ErrAuthFailed) {
		return "Senha incorreta ou arquivo corrompido."
	}
	if errors.Is(err, storage.ErrCorrupted) {
		return "Arquivo corrompido ou inválido."
	}
	if errors.Is(err, storage.ErrInvalidMagic) {
		return "O arquivo selecionado não é um cofre Abditum."
	}
	if errors.Is(err, storage.ErrVersionTooNew) {
		return "O cofre foi criado com uma versão mais recente do Abditum. Atualize o aplicativo."
	}
	return "Não foi possível abrir o cofre."
}

// isErrAutenticacao verifica se o erro indica falha de autenticação (senha incorreta).
// Usada para decidir se o usuário pode tentar novamente com outra senha.
func isErrAutenticacao(err error) bool {
	return errors.Is(err, crypto.ErrAuthFailed)
}
