package screen

import (
	"strings"

	"charm.land/lipgloss/v2"
	"github.com/useful-toys/abditum/internal/tui/actions"
	"github.com/useful-toys/abditum/internal/tui/design"
)

// ActionLineView renderiza a linha de ações disponíveis no contexto atual.
// Não implementa ChildView — é um renderizador stateless chamado diretamente pelo root.
// O zero value é válido e produz uma linha com apenas espaços na largura correta.
type ActionLineView struct{}

// NewActionLineView cria uma nova instância da linha de ações.
func NewActionLineView() *ActionLineView {
	return &ActionLineView{}
}

// Render retorna a linha de ações com exatamente `width` colunas.
//
// Layout:
//
//	[2 espaços][ação₁][ · ][ação₂][ · ]…[padding][F1 Ajuda]
//
// A âncora F1 é identificada por design.Shortcuts.Help e fixada à direita.
// Ações que não cabem no espaço disponível são descartadas (as de maior Priority, mais à direita).
// `acts` deve estar pré-ordenada por Priority crescente (menor = mais à esquerda).
func (v *ActionLineView) Render(width int, theme *design.Theme, acts []actions.Action) string {
	const (
		prefixCols = 2 // 2 espaços à esquerda
		anchorCols = 8 // reservado para "F1 Ajuda" ou espaços quando âncora ausente
		minPadding = 1 // pelo menos 1 espaço entre ações normais e âncora
	)

	// Separar âncora (F1) das demais ações.
	var anchor *actions.Action
	var normal []actions.Action
	for i := range acts {
		a := acts[i]
		if len(a.Keys) > 0 &&
			a.Keys[0].Code == design.Shortcuts.Help.Code &&
			a.Keys[0].Mod == design.Shortcuts.Help.Mod {
			anchor = &a
		} else {
			normal = append(normal, a)
		}
	}

	// Espaço disponível para ações normais: total menos prefixo, padding mínimo e âncora.
	availableCols := width - prefixCols - minPadding - anchorCols

	// Renderizar ações normais que cabem no espaço disponível.
	// Usar design.RenderAction e design.ActionSeparator para estilização consistente.
	separator := design.ActionSeparator(theme)

	type renderedActionItem struct {
		text  string
		width int
	}
	var renderedNormal []renderedActionItem
	usedCols := 0

	for _, a := range normal {
		// Pular ações que não estão visíveis ou não têm tecla
		if !a.Visible || len(a.Keys) == 0 {
			continue
		}

		// Renderizar ação usando design helper
		rendered := design.RenderAction(a.Keys[0].Label, a.Label, theme)
		needed := rendered.Width

		// Contar espaço do separador se não é a primeira ação
		if len(renderedNormal) > 0 {
			needed += separator.Width
		}

		// Parar se não cabe no espaço disponível
		if usedCols+needed > availableCols {
			break
		}

		renderedNormal = append(renderedNormal, renderedActionItem{
			text:  rendered.Text,
			width: rendered.Width,
		})
		usedCols += needed
	}

	// Montar bloco de ações normais com separadores.
	var normalBuilder strings.Builder
	for i, item := range renderedNormal {
		if i > 0 {
			normalBuilder.WriteString(separator.Text)
		}
		normalBuilder.WriteString(item.text)
	}
	normalText := normalBuilder.String()

	// Calcular padding entre ações normais e âncora.
	// Usar lipgloss.Width para contar com espaçamento correto
	paddingCols := width - prefixCols - lipgloss.Width(normalText) - anchorCols
	if paddingCols < minPadding {
		paddingCols = minPadding
	}

	// Renderizar âncora ou preencher com espaços quando ausente.
	var anchorText string
	if anchor != nil && len(anchor.Keys) > 0 {
		renderedAnchor := design.RenderAction(anchor.Keys[0].Label, anchor.Label, theme)
		anchorText = renderedAnchor.Text
	} else {
		anchorText = strings.Repeat(" ", anchorCols)
	}

	return strings.Repeat(" ", prefixCols) +
		normalText +
		strings.Repeat(" ", paddingCols) +
		anchorText
}
