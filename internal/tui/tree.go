package tui

import (
	"strings"
	"unicode/utf8"

	"github.com/useful-toys/abditum/internal/vault"
)

// treeNode represents a visible item in the vault tree.
type treeNode struct {
	kind    nodeKind
	depth   int // indentation level (0 = root)
	pasta   *vault.Pasta
	segredo *vault.Segredo
	// folder state
	expanded bool
	// computed prefix and label
	prefix string // e.g. "▼ " "▶ " "▷ " "● " "★ " "✦ " "✎ " "✗ "
	label  string
	estado vault.EstadoSessao
}

type nodeKind int

const (
	nodeFolder nodeKind = iota
	nodeSecret
	nodeFavorites // virtual folder
	nodeFavItem   // item inside Favorites
)

// treeModel tracks the tree panel state.
type treeModel struct {
	nodes       []treeNode
	cursor      int // index into nodes
	scrollTop   int // first visible node index
	expanded    map[*vault.Pasta]bool
	favExpanded bool
}

func newTreeModel() treeModel {
	return treeModel{
		expanded: make(map[*vault.Pasta]bool),
	}
}

// rebuild rebuilds the flat visible-node list from the vault.
func (t *treeModel) rebuild(cofre *vault.Cofre) {
	t.nodes = t.nodes[:0]

	// Collect favorites
	var favs []*vault.Segredo
	t.collectFavorites(cofre.PastaGeral(), &favs)

	if len(favs) > 0 {
		t.nodes = append(t.nodes, treeNode{
			kind:     nodeFavorites,
			depth:    0,
			expanded: t.favExpanded,
			prefix:   favPrefix(t.favExpanded),
			label:    "Favoritos",
		})
		if t.favExpanded {
			for _, s := range favs {
				prefix := secretPrefix(s)
				t.nodes = append(t.nodes, treeNode{
					kind:    nodeFavItem,
					depth:   1,
					segredo: s,
					prefix:  prefix,
					label:   s.Nome(),
					estado:  s.EstadoSessao(),
				})
			}
		}
	}

	// Build main tree from root folder
	t.buildFolder(cofre.PastaGeral(), 0)
}

func (t *treeModel) collectFavorites(p *vault.Pasta, result *[]*vault.Segredo) {
	for _, s := range p.Segredos() {
		if s.Favorito() && s.EstadoSessao() != vault.EstadoExcluido {
			*result = append(*result, s)
		}
	}
	for _, sub := range p.Subpastas() {
		t.collectFavorites(sub, result)
	}
}

func (t *treeModel) buildFolder(p *vault.Pasta, depth int) {
	hasChildren := len(p.Subpastas()) > 0 || len(activeSecrets(p)) > 0

	var prefix string
	exp := t.expanded[p]
	if !hasChildren {
		prefix = "▷ "
	} else if exp {
		prefix = "▼ "
	} else {
		prefix = "▶ "
	}

	t.nodes = append(t.nodes, treeNode{
		kind:     nodeFolder,
		depth:    depth,
		pasta:    p,
		expanded: exp,
		prefix:   prefix,
		label:    p.Nome(),
	})

	if !exp {
		return
	}

	for _, sub := range p.Subpastas() {
		t.buildFolder(sub, depth+1)
	}
	for _, s := range p.Segredos() {
		if s.EstadoSessao() == vault.EstadoExcluido {
			// still show with ✗ styling
		}
		prefix := secretPrefix(s)
		t.nodes = append(t.nodes, treeNode{
			kind:    nodeSecret,
			depth:   depth + 1,
			segredo: s,
			prefix:  prefix,
			label:   s.Nome(),
			estado:  s.EstadoSessao(),
		})
	}
}

func activeSecrets(p *vault.Pasta) []*vault.Segredo {
	var result []*vault.Segredo
	for _, s := range p.Segredos() {
		if s.EstadoSessao() != vault.EstadoExcluido {
			result = append(result, s)
		}
	}
	return result
}

func countActive(p *vault.Pasta) int {
	count := len(activeSecrets(p))
	for _, sub := range p.Subpastas() {
		count += countActive(sub)
	}
	return count
}

func secretPrefix(s *vault.Segredo) string {
	switch s.EstadoSessao() {
	case vault.EstadoIncluido:
		return "✦ "
	case vault.EstadoModificado:
		return "✎ "
	case vault.EstadoExcluido:
		return "✗ "
	}
	if s.Favorito() {
		return "★ "
	}
	return "● "
}

func favPrefix(expanded bool) string {
	if expanded {
		return "▼ "
	}
	return "▶ "
}

// expandFolder toggles a folder's expanded state and rebuilds.
func (t *treeModel) toggleFolder(cofre *vault.Cofre) {
	if t.cursor < 0 || t.cursor >= len(t.nodes) {
		return
	}
	node := t.nodes[t.cursor]
	switch node.kind {
	case nodeFolder:
		t.expanded[node.pasta] = !t.expanded[node.pasta]
	case nodeFavorites:
		t.favExpanded = !t.favExpanded
	}
	t.rebuild(cofre)
	t.clampScroll(0)
}

// expandFolderIfCollapsed expands a folder without toggling.
func (t *treeModel) expandFolder(cofre *vault.Cofre) {
	if t.cursor < 0 || t.cursor >= len(t.nodes) {
		return
	}
	node := t.nodes[t.cursor]
	switch node.kind {
	case nodeFolder:
		if !t.expanded[node.pasta] {
			t.expanded[node.pasta] = true
			t.rebuild(cofre)
			// move cursor to first child
			if t.cursor+1 < len(t.nodes) {
				t.cursor++
			}
		} else if t.cursor+1 < len(t.nodes) {
			t.cursor++
		}
	case nodeFavorites:
		if !t.favExpanded {
			t.favExpanded = true
			t.rebuild(cofre)
			if t.cursor+1 < len(t.nodes) {
				t.cursor++
			}
		}
	}
	t.clampScroll(0)
}

// collapseFolder collapses current folder or moves to parent.
func (t *treeModel) collapseFolder(cofre *vault.Cofre) {
	if t.cursor < 0 || t.cursor >= len(t.nodes) {
		return
	}
	node := t.nodes[t.cursor]
	switch node.kind {
	case nodeFolder:
		if t.expanded[node.pasta] {
			t.expanded[node.pasta] = false
			t.rebuild(cofre)
		} else {
			// move to parent
			t.moveToParent()
		}
	case nodeFavorites:
		if t.favExpanded {
			t.favExpanded = false
			t.rebuild(cofre)
		}
	default:
		t.moveToParent()
	}
	t.clampScroll(0)
}

func (t *treeModel) moveToParent() {
	if t.cursor <= 0 {
		return
	}
	cur := t.nodes[t.cursor]
	for i := t.cursor - 1; i >= 0; i-- {
		n := t.nodes[i]
		if n.depth < cur.depth && (n.kind == nodeFolder || n.kind == nodeFavorites) {
			t.cursor = i
			return
		}
	}
}

func (t *treeModel) moveUp() {
	if t.cursor > 0 {
		t.cursor--
	}
	t.clampScroll(0)
}

func (t *treeModel) moveDown() {
	if t.cursor < len(t.nodes)-1 {
		t.cursor++
	}
	t.clampScroll(0)
}

func (t *treeModel) moveHome() {
	t.cursor = 0
	t.clampScroll(0)
}

func (t *treeModel) moveEnd() {
	if len(t.nodes) > 0 {
		t.cursor = len(t.nodes) - 1
	}
	t.clampScroll(0)
}

// clampScroll ensures cursor is visible within [scrollTop, scrollTop+height).
func (t *treeModel) clampScroll(height int) {
	if height <= 0 {
		return
	}
	if t.cursor < t.scrollTop {
		t.scrollTop = t.cursor
	}
	if t.cursor >= t.scrollTop+height {
		t.scrollTop = t.cursor - height + 1
	}
	if t.scrollTop < 0 {
		t.scrollTop = 0
	}
}

// selectedSecret returns the segredo under the cursor, or nil if cursor is on a folder.
func (t *treeModel) selectedSecret() *vault.Segredo {
	if t.cursor < 0 || t.cursor >= len(t.nodes) {
		return nil
	}
	n := t.nodes[t.cursor]
	if n.kind == nodeSecret || n.kind == nodeFavItem {
		return n.segredo
	}
	return nil
}

// renderTree renders the tree panel into `height` lines of `width` columns.
// separatorCol: the column index of the │ separator between tree and detail.
func renderTree(st styles, t treeModel, height, width, separatorCol int, detailSecret *vault.Segredo, focused bool) string {
	t.clampScroll(height)

	lines := make([]string, height)

	sepStyle := st.BorderDefault
	if focused {
		sepStyle = st.BorderFocused
	}

	for row := 0; row < height; row++ {
		nodeIdx := t.scrollTop + row
		lineWidth := width - 1 // leave room for separator col

		if nodeIdx >= len(t.nodes) {
			lines[row] = strings.Repeat(" ", lineWidth) + sepStyle.Render("│")
			continue
		}

		node := t.nodes[nodeIdx]
		indent := strings.Repeat("  ", node.depth)
		prefixStr := node.prefix

		// Format: indent + prefix + label + counter(folder)
		var content string
		if node.kind == nodeFolder {
			total := countActive(node.pasta)
			counter := "(" + itoa(total) + ")"
			labelPart := node.label
			// available width for label
			avail := lineWidth - len([]rune(indent)) - len([]rune(prefixStr)) - len([]rune(counter)) - 1
			if avail < 0 {
				avail = 0
			}
			labelT := truncateRight(labelPart, avail)
			// Pad between label and counter
			gap := lineWidth - len([]rune(indent)) - len([]rune(prefixStr)) - len([]rune(labelT)) - len([]rune(counter))
			if gap < 1 {
				gap = 1
			}
			content = indent +
				st.TextSecondary.Render(prefixStr) +
				st.TextPrimary.Render(labelT) +
				strings.Repeat(" ", gap) +
				st.TextSecondary.Render(counter)
		} else if node.kind == nodeFavorites {
			total := countFavorites(t)
			counter := "(" + itoa(total) + ")"
			labelT := truncateRight(node.label, lineWidth-len([]rune(indent))-len([]rune(prefixStr))-len([]rune(counter))-1)
			gap := lineWidth - len([]rune(indent)) - len([]rune(prefixStr)) - len([]rune(labelT)) - len([]rune(counter))
			if gap < 1 {
				gap = 1
			}
			content = indent +
				st.TextSecondary.Render(prefixStr) +
				st.AccentPrimary.Bold(true).Render(labelT) +
				strings.Repeat(" ", gap) +
				st.TextSecondary.Render(counter)
		} else {
			// secret node
			s := node.segredo
			avail := lineWidth - len([]rune(indent)) - len([]rune(prefixStr))
			if avail < 0 {
				avail = 0
			}
			labelT := truncateRight(node.label, avail)

			var prefixRendered string
			var labelRendered string
			switch s.EstadoSessao() {
			case vault.EstadoIncluido:
				prefixRendered = st.SemanticWarning.Render(prefixStr)
				labelRendered = st.SemanticWarning.Render(labelT)
			case vault.EstadoModificado:
				prefixRendered = st.SemanticWarning.Render(prefixStr)
				labelRendered = st.SemanticWarning.Render(labelT)
			case vault.EstadoExcluido:
				prefixRendered = st.SemanticWarning.Render(prefixStr)
				labelRendered = st.SemanticWarning.Strikethrough(true).Render(labelT)
			default:
				if s.Favorito() {
					prefixRendered = st.AccentSecondary.Render(prefixStr)
				} else {
					prefixRendered = st.TextSecondary.Render(prefixStr)
				}
				labelRendered = st.TextPrimary.Render(labelT)
			}
			content = indent + prefixRendered + labelRendered
		}

		// Apply selection highlight
		if nodeIdx == t.cursor {
			// Pad content to lineWidth for full-width highlight
			visLen := visibleLen(content)
			if visLen < lineWidth {
				content += strings.Repeat(" ", lineWidth-visLen)
			}
			content = st.Selected.Render(stripANSI(content))
		}

		// Determine separator character for this row
		sepChar := "│"
		// Check if this row shows the selected secret connecting to detail
		if detailSecret != nil &&
			nodeIdx == t.cursor &&
			(t.nodes[nodeIdx].kind == nodeSecret || t.nodes[nodeIdx].kind == nodeFavItem) {
			sepChar = "<╡"
		}
		// Scroll indicators override separator (except <╡)
		if sepChar == "│" {
			if len(t.nodes) > height {
				if row == 0 && t.scrollTop > 0 {
					sepChar = "↑"
				} else if row == height-1 && t.scrollTop+height < len(t.nodes) {
					sepChar = "↓"
				} else {
					// thumb position
					thumbRow := thumbPosition(t.scrollTop, len(t.nodes)-height, height)
					if row == thumbRow {
						sepChar = "■"
					}
				}
			}
		}

		// Pad content to lineWidth
		visLen := visibleLen(content)
		if visLen < lineWidth {
			content += strings.Repeat(" ", lineWidth-visLen)
		}

		lines[row] = content + sepStyle.Render(sepChar)
	}

	return strings.Join(lines, "\n")
}

func countFavorites(t treeModel) int {
	count := 0
	for _, n := range t.nodes {
		if n.kind == nodeFavItem {
			count++
		}
	}
	return count
}

func thumbPosition(scrollTop, maxScroll, height int) int {
	if maxScroll <= 0 {
		return 0
	}
	pos := int(float64(scrollTop) / float64(maxScroll) * float64(height-1))
	if pos >= height {
		pos = height - 1
	}
	return pos
}

// visibleLen estimates the visible column width of a string with ANSI escapes.
// Approximation: counts runes outside escape sequences.
func visibleLen(s string) int {
	inEscape := false
	count := 0
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		count += utf8.RuneLen(r) / utf8.RuneLen(r) // count each rune as 1 col
	}
	return count
}

// stripANSI removes ANSI escape sequences from s.
func stripANSI(s string) string {
	var b strings.Builder
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
			continue
		}
		if inEscape {
			if r == 'm' {
				inEscape = false
			}
			continue
		}
		b.WriteRune(r)
	}
	return b.String()
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	negative := n < 0
	if negative {
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if negative {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
