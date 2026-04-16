package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tea "charm.land/bubbletea/v2"
	"github.com/atotto/clipboard"

	"github.com/useful-toys/abditum/internal/tui"
)

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
	flag.StringVar(&vaultPath, "vault", "", "Path to the Abditum vault file")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	defer clipboard.WriteAll("") //nolint:errcheck

	// Ensure terminal is restored on exit (in case of panic or unexpected exit)
	defer restoreTerminal()

	root := tui.NewRootModel(tui.WithVersion(version))
	// Setup all actions (system and application)
	setupActions(root)
	// Enable alternate screen (cleaner TUI) and mouse support for interactive components
	p := tea.NewProgram(root,
		tea.WithContext(ctx),
	)
	_, err := p.Run()
	return err
}

// restoreTerminal ensures the terminal is in a clean state.
// This is called on defer to handle unexpected exits or panics.
func restoreTerminal() {
	// Restore cursor visibility
	fmt.Fprint(os.Stdout, "\033[?25h")
	// Disable alternate screen if still active
	fmt.Fprint(os.Stdout, "\033[?1049l")
	// Clear any remaining content
	fmt.Fprint(os.Stdout, "\033[2J\033[H")
	os.Stdout.Sync()
}
