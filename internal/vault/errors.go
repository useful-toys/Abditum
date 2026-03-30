// Package vault implements the domain layer for Abditum password manager.
package vault

import "errors"

// Sentinel errors for validation failures.
// TUI can use errors.Is() for type checking.
var (
	ErrNomeVazio              = errors.New("nome não pode ser vazio")
	ErrNomeMuitoLongo         = errors.New("nome excede 255 caracteres")
	ErrNameConflict           = errors.New("já existe item com este nome")
	ErrPastaGeralProtected    = errors.New("Pasta Geral não pode ser modificada")
	ErrPastaGeralNaoExcluivel = errors.New("Pasta Geral não pode ser excluída")
	ErrCycleDetected          = errors.New("operação criaria ciclo na hierarquia")
	ErrObservacaoReserved     = errors.New("nome 'Observação' é reservado")
	ErrConfigInvalida         = errors.New("configuração inválida")
	ErrPosicaoInvalida        = errors.New("posição inválida")
	ErrSegredoNaoEncontrado   = errors.New("segredo não encontrado")
	ErrCofreBloqueado         = errors.New("cofre está bloqueado")
	ErrDestinoInvalido        = errors.New("destino inválido para operação")
	ErrCampoInvalido          = errors.New("índice de campo inválido")
	ErrModeloEmUso            = errors.New("modelo está em uso por segredos")
	ErrNomeReservado          = errors.New("nome é reservado e não pode ser usado")
	ErrPastaInvalida          = errors.New("pasta inválida ou nula")
	ErrModeloInvalido         = errors.New("modelo inválido ou nulo")
	ErrSegredoInvalido        = errors.New("segredo inválido ou nulo")
	ErrSegredoJaExcluido      = errors.New("segredo já está marcado como excluído")
	ErrSegredoNaoExcluido     = errors.New("segredo não está marcado como excluído")
	ErrObservacaoMuitoLonga   = errors.New("observação excede 1000 caracteres")
)
