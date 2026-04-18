package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"

	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/tui/operation"
)

// Bubble Tea v2 manages the alternate screen lifecycle automatically via the
// AltScreen field on tea.View — no manual terminal restoration is needed.

// version is injected at build time via -ldflags "-X main.version=$(git describe --tags --always)"
// In local builds without tags, defaults to "dev"
// Never hardcoded in source — always injected or defaults to dev
var version = "dev"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal: could not start Abditum")
		os.Exit(1)
	}
}

func run() error {
	var vaultPath string
	flag.StringVar(&vaultPath, "vault", "", "Caminho para o arquivo de cofre Abditum")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	defer clipboard.WriteAll("") //nolint:errcheck

	root := tui.NewRootModel(tui.WithVersion(version))
	// Setup all actions (system and application)
	setupActions(root)
	if vaultPath != "" {
		if cmd := buildVaultCmd(vaultPath, root); cmd != nil {
			root.SetInitialCommand(cmd)
		}
	}
	p := tea.NewProgram(root,
		tea.WithContext(ctx),
	)
	_, err := p.Run()
	return err
}

// buildVaultCmd decide qual operação disparar com base no caminho --vault.
//   - Arquivo existe → Fluxo 1 (Abrir cofre)
//   - Arquivo não existe mas o diretório pai existe → Fluxo 2 (Criar cofre)
//   - Qualquer outro caso → nil (tela normal, sem operação automática)
//
// Recebe root para acessar MessageController, que não pode ser nil.
func buildVaultCmd(vaultPath string, root *tui.RootModel) tea.Cmd {
	info, err := os.Stat(vaultPath)
	if err == nil && !info.IsDir() {
		// Arquivo existe: abrir o cofre diretamente.
		return tui.StartOperation(
			operation.NewAbrirCofreOperation(root.MessageController(), nil, vaultPath),
		)
	}
	if os.IsNotExist(err) {
		// Arquivo não existe: verificar se o diretório pai existe para criar.
		dir := filepath.Dir(vaultPath)
		if dirInfo, dirErr := os.Stat(dir); dirErr == nil && dirInfo.IsDir() {
			return tui.StartOperation(
				operation.NewCriarCofreOperation(root.MessageController(), nil, vaultPath),
			)
		}
	}
	// Caminho inválido (ex: diretório pai inexistente): iniciar normalmente sem operação.
	return nil
}
