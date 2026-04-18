package operation

import (
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/useful-toys/abditum/internal/crypto"
	"github.com/useful-toys/abditum/internal/storage"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/design"
	"github.com/useful-toys/abditum/internal/tui/modal"
	"github.com/useful-toys/abditum/internal/vault"
)

// criarEstado representa em qual etapa do fluxo de criação estamos.
type criarEstado int

const (
	criandoInformandoCaminho      criarEstado = iota // FilePicker modo Save
	criandoConfirmandoSobrescrita                    // ConfirmModal "arquivo já existe"
	criandoInformandoSenha                           // PasswordCreateModal
	criandoAvaliacaoSenhaFraca                       // ConfirmModal "senha fraca"
	criandoCriando                                   // IO assíncrono
)

// criarAvancaMsg é a mensagem interna de transição de estado.
type criarAvancaMsg struct {
	estado criarEstado
}

// criarCofreResultMsg carrega o resultado da criação do cofre.
type criarCofreResultMsg struct {
	manager *vault.Manager
	err     error
}

// CriarCofreOperation implementa o Fluxo 2 (Criar Novo Cofre).
//
// Se caminhoInicial != "", pula o guard e o picker e vai direto para a entrada de senha.
// Se caminhoInicial == "", executa o fluxo completo via GUI.
type CriarCofreOperation struct {
	notifier tui.MessageController
	saver    vaultSaver
	guard    *guardCofreAlterado // não-nil apenas no fluxo GUI com cofre alterado
	caminho  string              // caminho destino selecionado
	senha    []byte              // senha informada
}

// NewCriarCofreOperation cria a operation com o vault.Manager concreto.
// manager pode ser nil quando nenhum cofre está carregado.
// caminhoInicial pode ser "" para fluxo GUI completo.
func NewCriarCofreOperation(
	notifier tui.MessageController,
	manager *vault.Manager,
	caminhoInicial string,
) *CriarCofreOperation {
	var saver vaultSaver
	if manager != nil {
		saver = manager
	}
	return &CriarCofreOperation{
		notifier: notifier,
		saver:    saver,
		caminho:  caminhoInicial,
	}
}

// newCriarCofreOperationFromSaver é usada nos testes (mesmo pacote).
func newCriarCofreOperationFromSaver(
	notifier tui.MessageController,
	saver vaultSaver,
	caminhoInicial string,
) *CriarCofreOperation {
	return &CriarCofreOperation{
		notifier: notifier,
		saver:    saver,
		caminho:  caminhoInicial,
	}
}

// Init inicia o fluxo de criação.
func (c *CriarCofreOperation) Init() tea.Cmd {
	if c.caminho != "" {
		// Entrada via CLI: pular guard e picker, ir direto para senha
		return c.abrirModalSenha()
	}
	// Fluxo GUI completo: verificar cofre alterado primeiro
	c.guard = novoGuardCofreAlterado(
		c.notifier,
		c.saver,
		func() tea.Cmd {
			return func() tea.Msg { return criarAvancaMsg{estado: criandoInformandoCaminho} }
		},
		func() tea.Cmd { return tui.OperationCompleted() },
	)
	return c.guard.Init()
}

// Update trata as mensagens internas da CriarCofreOperation.
func (c *CriarCofreOperation) Update(msg tea.Msg) tea.Cmd {
	switch m := msg.(type) {
	// Mensagens do guard — delegar se guard ativo
	case guardSaveMsg, guardSaveResultMsg, guardDiscardMsg, guardCancelMsg:
		if c.guard != nil {
			return c.guard.Update(msg)
		}
		return nil

	case criarAvancaMsg:
		switch m.estado {
		case criandoInformandoCaminho:
			return c.abrirFilePicker()
		case criandoConfirmandoSobrescrita:
			return c.abrirModalSobrescrita()
		case criandoInformandoSenha:
			return c.abrirModalSenha()
		case criandoAvaliacaoSenhaFraca:
			return c.abrirModalSenhaFraca()
		case criandoCriando:
			c.notifier.SetBusy("Criando cofre...")
			caminho := c.caminho
			senha := c.senha
			return func() tea.Msg {
				cofre := vault.NovoCofre()
				if err := cofre.InicializarConteudoPadrao(); err != nil {
					return criarCofreResultMsg{err: err}
				}
				repo := storage.NewFileRepositoryForCreate(caminho, senha)
				manager := vault.NewManager(cofre, repo)
				if err := manager.Salvar(false); err != nil {
					return criarCofreResultMsg{err: err}
				}
				return criarCofreResultMsg{manager: manager}
			}
		}

	case criarCofreResultMsg:
		if m.err != nil {
			c.notifier.SetError(m.err.Error())
			return tui.OperationCompleted()
		}
		c.notifier.Clear()
		return func() tea.Msg { return tui.VaultOpenedMsg{Manager: m.manager} }
	}

	return nil
}

// abrirFilePicker abre o FilePicker no modo Save com extensão .abditum.
func (c *CriarCofreOperation) abrirFilePicker() tea.Cmd {
	return tui.OpenModal(modal.NewFilePicker(modal.FilePickerOptions{
		Mode:      modal.FilePickerSave,
		Extension: ".abditum",
		Messages:  c.notifier,
		OnResult: func(path string) tea.Cmd {
			if path == "" {
				return tui.OperationCompleted()
			}
			c.caminho = path
			if arquivoExiste(path) {
				return func() tea.Msg { return criarAvancaMsg{estado: criandoConfirmandoSobrescrita} }
			}
			return func() tea.Msg { return criarAvancaMsg{estado: criandoInformandoSenha} }
		},
	}))
}

// abrirModalSobrescrita abre o modal de confirmação de sobrescrita.
func (c *CriarCofreOperation) abrirModalSobrescrita() tea.Cmd {
	return tui.OpenModal(modal.NewConfirmModal(
		"Arquivo existente",
		"O arquivo já existe. Deseja sobrescrever?",
		[]modal.ModalOption{
			{
				Keys:  []design.Key{design.Keys.Enter},
				Label: "Sobrescrever",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return criarAvancaMsg{estado: criandoInformandoSenha}
					})
				},
			},
			{
				Keys:  []design.Key{design.Keys.Esc},
				Label: "Outro caminho",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return criarAvancaMsg{estado: criandoInformandoCaminho}
					})
				},
			},
		},
	))
}

// abrirModalSenha abre o PasswordCreateModal.
func (c *CriarCofreOperation) abrirModalSenha() tea.Cmd {
	return tui.OpenModal(modal.NewPasswordCreateModal(
		c.notifier,
		func(password []byte) tea.Cmd {
			c.senha = password
			if crypto.EvaluatePasswordStrength(password) == crypto.StrengthWeak {
				return tea.Batch(tui.CloseModal(), func() tea.Msg {
					return criarAvancaMsg{estado: criandoAvaliacaoSenhaFraca}
				})
			}
			return tea.Batch(tui.CloseModal(), func() tea.Msg {
				return criarAvancaMsg{estado: criandoCriando}
			})
		},
		func() tea.Cmd {
			// Cancelar: voltar ao picker se houve fluxo GUI, ou completar se entrada CLI
			if c.saver != nil || c.guard != nil {
				return tea.Batch(tui.CloseModal(), func() tea.Msg {
					return criarAvancaMsg{estado: criandoInformandoCaminho}
				})
			}
			return tea.Batch(tui.CloseModal(), tui.OperationCompleted())
		},
	))
}

// abrirModalSenhaFraca abre o modal de aviso de senha fraca.
func (c *CriarCofreOperation) abrirModalSenhaFraca() tea.Cmd {
	return tui.OpenModal(modal.NewConfirmModal(
		"Senha fraca",
		"A senha informada é fraca. Deseja prosseguir assim mesmo?",
		[]modal.ModalOption{
			{
				Keys:  []design.Key{design.Keys.Enter},
				Label: "Prosseguir",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return criarAvancaMsg{estado: criandoCriando}
					})
				},
			},
			{
				Keys:  []design.Key{design.Keys.Esc},
				Label: "Revisar",
				Action: func() tea.Cmd {
					return tea.Batch(tui.CloseModal(), func() tea.Msg {
						return criarAvancaMsg{estado: criandoInformandoSenha}
					})
				},
			},
		},
	))
}

// arquivoExiste reporta se o caminho aponta para um arquivo existente.
func arquivoExiste(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}
