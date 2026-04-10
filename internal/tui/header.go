package tui

import (
	"strings"
	"unicode/utf8"
)

// appMode representa o modo ativo da área de trabalho.
type appMode int

const (
	modeWelcome appMode = iota
	modeVault
	modeModels
	modeConfig
)

// renderHeader retorna as 2 linhas do cabeçalho.
// Linha 1: nome da app [· cofre [•]] [abas]
// Linha 2: separador ─ com aba ativa suspensa
func renderHeader(st styles, width int, mode appMode, vaultName string, dirty bool, searchQuery string, searchActive bool) string {
	if vaultName == "" {
		// Estado: sem cofre aberto (boas-vindas)
		line1 := st.AppName.Render("Abditum")
		line2 := st.BorderDefault.Render(strings.Repeat("─", width))
		return line1 + "\n" + line2
	}

	// Abas — larguras fixas
	const (
		tabVaultText   = " Cofre "
		tabModelsText  = " Modelos "
		tabConfigText  = " Config "
		tabVaultWidth  = len("╭") + len(tabVaultText) + len("╮")                     // 9
		tabModelsWidth = len("╭") + len(tabModelsText) + len("╮")                    // 11
		tabConfigWidth = len("╭") + len(tabConfigText) + len("╮")                    // 10
		tabsBlock      = tabVaultWidth + 2 + tabModelsWidth + 2 + tabConfigWidth + 2 // +2 espaços entre cada
	)
	// 2 spaces prefix + "Abditum · " prefix (10 chars) = 12 fixed cols
	// dirty indicator " •" if set = 2 cols
	// then we need: 1 padding + tabsBlock
	prefix := "  " + st.AppName.Render("Abditum") + st.BorderDefault.Render(" · ")
	prefixLen := 10 // "  Abditum · " (visual)

	dirtyStr := ""
	dirtyLen := 0
	if dirty {
		dirtyStr = " " + st.DirtyDot.Render("•")
		dirtyLen = 2
	}

	// calculate space for vault name
	padding := 1 // min 1 col between name and tabs block
	available := width - prefixLen - dirtyLen - padding - tabsBlock
	if available < 0 {
		available = 0
	}

	displayName := truncateRight(vaultName, available)
	vaultNameStr := st.CoffeeName.Render(displayName)

	// Build tab strings
	tabVault := renderTab(st, "Cofre", mode == modeVault)
	tabModels := renderTab(st, "Modelos", mode == modeModels)
	tabConfig := renderTab(st, "Config", mode == modeConfig)

	line1 := prefix + vaultNameStr + dirtyStr +
		strings.Repeat(" ", max(padding, width-prefixLen-utf8.RuneCountInString(displayName)-dirtyLen-tabsBlock-1)) +
		" " + tabVault + "  " + tabModels + "  " + tabConfig

	// Line 2: separator with active tab suspended
	line2 := renderSeparatorLine(st, width, mode, searchQuery, searchActive)

	return line1 + "\n" + line2
}

// renderTab renders a tab (active or inactive).
func renderTab(st styles, label string, active bool) string {
	if active {
		// Active: top shows ╭──────╮, bottom will show ╯ Text ╰
		inner := strings.Repeat("─", utf8.RuneCountInString(label)+2)
		return st.TabBorder.Render("╭" + inner + "╮")
	}
	return st.TabBorder.Render("╭") + st.TabInactive.Render(" "+label+" ") + st.TabBorder.Render("╮")
}

// renderSeparatorLine renders line 2 of the header (─── with active tab suspended).
func renderSeparatorLine(st styles, width int, mode appMode, searchQuery string, searchActive bool) string {
	if mode == modeWelcome {
		return st.BorderDefault.Render(strings.Repeat("─", width))
	}

	// Determine active tab label and its width on the line
	var activeLabel string
	var tabOffset int
	switch mode {
	case modeVault:
		activeLabel = "Cofre"
		tabOffset = computeTabOffset(width, 0) // first tab
	case modeModels:
		activeLabel = "Modelos"
		tabOffset = computeTabOffset(width, 1)
	case modeConfig:
		activeLabel = "Config"
		tabOffset = computeTabOffset(width, 2)
	}

	// Active-tab suspended text: ╰<space><label><space>╯
	activeTabRendered := st.TabBorder.Render("╯") +
		st.TabActive.Render(" "+activeLabel+" ") +
		st.TabBorder.Render("╰")
	activeTabVisLen := 1 + 1 + utf8.RuneCountInString(activeLabel) + 1 + 1 // ╯ + space + label + space + ╰

	if searchActive && mode == modeVault {
		// Search-mode line 2: ─ Busca: <query> ────── ╯ Cofre ╰ ──
		prefix := st.BorderDefault.Render("─ ") + st.TextSecondary.Render("Busca:") + st.BorderDefault.Render(" ")
		prefixLen := 9 // "─ Busca: "
		// space for query
		queryAvail := tabOffset - prefixLen - 1
		if queryAvail < 0 {
			queryAvail = 0
		}
		var queryStr string
		if searchQuery == "" {
			queryStr = strings.Repeat("─", queryAvail)
		} else {
			q := truncateLeft(searchQuery, queryAvail)
			queryStr = st.AccentPrimary.Bold(true).Render(q)
			fill := queryAvail - utf8.RuneCountInString(q)
			if fill > 0 {
				queryStr += st.BorderDefault.Render(strings.Repeat("─", fill))
			}
		}
		suffix := st.BorderDefault.Render(strings.Repeat("─", max(0, width-tabOffset-activeTabVisLen)))
		return prefix + queryStr + activeTabRendered + suffix
	}

	// Normal separator line
	leftDashes := st.BorderDefault.Render(strings.Repeat("─", tabOffset))
	rightDashes := st.BorderDefault.Render(strings.Repeat("─", max(0, width-tabOffset-activeTabVisLen)))
	return leftDashes + activeTabRendered + rightDashes
}

// computeTabOffset returns the column position (0-based) where the active tab starts on line 2.
// tabs are placed at the right side of the header, with 2-space separators.
func computeTabOffset(width, tabIdx int) int {
	// tab widths: Cofre=9, Modelos=11, Config=10 (╭ + space + label + space + ╮)
	tabWidths := [3]int{9, 11, 10}
	// total tabs block width = sum(widths) + 2*2 (two gaps of 2)
	total := tabWidths[0] + 2 + tabWidths[1] + 2 + tabWidths[2] + 2
	// right-aligned: starts at width - total
	start := width - total
	if start < 0 {
		start = 0
	}
	// offset to the requested tab
	off := start
	for i := 0; i < tabIdx; i++ {
		off += tabWidths[i] + 2
	}
	return off
}

// truncateRight truncates s to at most n rune widths, appending "…" if needed.
func truncateRight(s string, n int) string {
	if n <= 0 {
		return "…"
	}
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	if n <= 1 {
		return "…"
	}
	return string(runes[:n-1]) + "…"
}

// truncateLeft truncates s to at most n rune widths from the right, prepending "…" if needed.
func truncateLeft(s string, n int) string {
	if n <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	if n <= 1 {
		return "…"
	}
	return "…" + string(runes[len(runes)-(n-1):])
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
