package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"
	"github.com/useful-toys/abditum/internal/tui"
	"github.com/useful-toys/abditum/internal/vault"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal: could not start Abditum")
		os.Exit(1)
	}
}

func run() error {
	// Parse optional vault path argument
	var initialPath string
	if len(os.Args) > 1 {
		initialPath = os.Args[1]
	}

	// Phase 5: create an empty Manager stub (no vault loaded yet).
	// Phase 6 will open or create a real vault via modal flows.
	mgr := vault.NewManager(vault.NovoCofre(), nil)

	// Graceful shutdown on SIGTERM/SIGINT; clipboard always cleared on exit.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	defer clipboard.WriteAll("") //nolint:errcheck

	root := tui.NewRootModel(mgr, initialPath)
	p := tea.NewProgram(root, tea.WithContext(ctx))
	_, err := p.Run()
	return err
}
