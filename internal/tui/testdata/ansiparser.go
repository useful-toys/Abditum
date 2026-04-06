package testdata

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// StyleTransition representa uma transição de estilo visual (cor/fonte muda).
// Registra APENAS onde o estado (FG, BG, Style) muda realmente — sem redundâncias.
type StyleTransition struct {
	Line  int      `json:"line"`  // número da linha (0-indexed)
	Col   int      `json:"col"`   // coluna onde muda (0-indexed)
	FG    *string  `json:"fg"`    // cor foreground em hex (ex: "#a0a0a0"), ou null se default
	BG    *string  `json:"bg"`    // cor background em hex, ou null se default
	Style []string `json:"style"` // ["bold", "italic", ...], ou [] para nenhum
}

// MarshalJSON serializa StyleTransition como uma tupla compacta de 5 elementos:
//
//	[line, col, fg_hex_or_null, bg_hex_or_null, [styles]]
//
// Isso corresponde ao formato especificado em arquitetura-teste.md § "Formato de cada tupla".
// Exemplo de saída: [0, 4, "#d08c00", null, ["bold"]]
func (st StyleTransition) MarshalJSON() ([]byte, error) {
	type tuple [5]any
	var fg, bg any
	if st.FG != nil {
		fg = *st.FG
	}
	if st.BG != nil {
		bg = *st.BG
	}
	// Garante que Style nunca seja null no JSON — usa [] em vez de null quando nil.
	style := st.Style
	if style == nil {
		style = []string{}
	}
	return json.Marshal(tuple{st.Line, st.Col, fg, bg, style})
}

// ansiState holds the current rendering state during ANSI sequence parsing.
type ansiState struct {
	fg    *string // cor foreground em hex, ou nil
	bg    *string // cor background em hex, ou nil
	style map[string]bool
}

// ParseANSIStyle extrai transições de estilo visual do output ANSI.
// Normaliza SGR codes para estado visual canônico (fg hex, bg hex, style set).
// Registra tupla APENAS quando o estado muda — ignora redundâncias e ordem de SGR.
func ParseANSIStyle(output string) []StyleTransition {
	var result []StyleTransition

	currentState := ansiState{
		fg:    nil,
		bg:    nil,
		style: make(map[string]bool),
	}

	// Função auxiliar: converte state em tupla para comparação
	stateKey := func(s ansiState) string {
		key := ""
		if s.fg != nil {
			key += fmt.Sprintf("fg:%s|", *s.fg)
		} else {
			key += "fg:nil|"
		}
		if s.bg != nil {
			key += fmt.Sprintf("bg:%s|", *s.bg)
		} else {
			key += "bg:nil|"
		}
		for _, st := range []string{"bold", "italic", "underline", "strikethrough", "faint", "blink", "reverse"} {
			if s.style[st] {
				key += st + ","
			}
		}
		return key
	}

	lastStateKey := stateKey(currentState)

	// Regex para encontrar escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[([0-9;]*)m`)

	line := 0
	col := 0
	textPos := 0 // tracks position in output past the last processed escape sequence

	for _, match := range ansiRegex.FindAllStringSubmatchIndex(output, -1) {
		codeStr := output[match[2]:match[3]]

		// Processa texto antes deste escape sequence, começando de textPos
		textBefore := output[textPos:match[0]]
		for _, ch := range textBefore {
			if ch == '\n' {
				line++
				col = 0
			} else {
				col++
			}
		}
		textPos = match[1] // avança posição para depois do escape completo

		// Processa o código SGR
		if codeStr == "" {
			codeStr = "0" // reset
		}

		codes := strings.Split(codeStr, ";")
		for _, codeNum := range codes {
			code, _ := strconv.Atoi(strings.TrimSpace(codeNum))
			applyCode(&currentState, code)
		}

		// Checa se estado mudou
		newStateKey := stateKey(currentState)
		if newStateKey != lastStateKey {
			// Registra transição
			transition := StyleTransition{
				Line:  line,
				Col:   col,
				FG:    currentState.fg,
				BG:    currentState.bg,
				Style: styleMapToArray(currentState.style),
			}
			// Se a última transição registrada está na mesma posição (line, col),
			// substitui em vez de adicionar — múltiplos SGR na mesma coluna resultam
			// em apenas uma transição (o estado final prevalece).
			if len(result) > 0 && result[len(result)-1].Line == line && result[len(result)-1].Col == col {
				result[len(result)-1] = transition
			} else {
				result = append(result, transition)
			}
			lastStateKey = newStateKey
		}
	}

	return result
}

// applyCode aplica um código SGR ao estado
func applyCode(state *ansiState, code int) {
	switch code {
	case 0: // Reset
		state.fg = nil
		state.bg = nil
		state.style = make(map[string]bool)
	case 1: // Bold
		state.style["bold"] = true
	case 2: // Faint
		state.style["faint"] = true
	case 3: // Italic
		state.style["italic"] = true
	case 4: // Underline
		state.style["underline"] = true
	case 5: // Blink
		state.style["blink"] = true
	case 7: // Reverse
		state.style["reverse"] = true
	case 9: // Strikethrough
		state.style["strikethrough"] = true
	case 21, 22: // Normal (not bold/faint)
		state.style["bold"] = false
		state.style["faint"] = false
	case 23: // Normal (not italic)
		state.style["italic"] = false
	case 24: // Normal (not underline)
		state.style["underline"] = false
	case 25: // Normal (not blink)
		state.style["blink"] = false
	case 27: // Normal (not reverse)
		state.style["reverse"] = false
	case 29: // Normal (not strikethrough)
		state.style["strikethrough"] = false
	case 30, 31, 32, 33, 34, 35, 36, 37: // 16-color foreground
		color := colorCode16(code)
		state.fg = &color
	case 38: // 256-color or truecolor foreground (handled separately with next codes)
		// Será tratado em contexto de múltiplos códigos
	case 39: // Default foreground
		state.fg = nil
	case 40, 41, 42, 43, 44, 45, 46, 47: // 16-color background
		color := colorCode16(code - 10)
		state.bg = &color
	case 48: // 256-color or truecolor background (handled separately)
	case 49: // Default background
		state.bg = nil
	case 90, 91, 92, 93, 94, 95, 96, 97: // Bright foreground (16-color)
		color := colorCode16(code - 60 + 8)
		state.fg = &color
	case 100, 101, 102, 103, 104, 105, 106, 107: // Bright background (16-color)
		color := colorCode16(code - 100 + 8)
		state.bg = &color
	}
}

// colorCode16 converte código ANSI 16-color para hex
func colorCode16(code int) string {
	colors := []string{
		"#000000", // 0: black
		"#800000", // 1: red
		"#008000", // 2: green
		"#808000", // 3: yellow
		"#000080", // 4: blue
		"#800080", // 5: magenta
		"#008080", // 6: cyan
		"#c0c0c0", // 7: white
		"#808080", // 8: bright black
		"#ff0000", // 9: bright red
		"#00ff00", // 10: bright green
		"#ffff00", // 11: bright yellow
		"#0000ff", // 12: bright blue
		"#ff00ff", // 13: bright magenta
		"#00ffff", // 14: bright cyan
		"#ffffff", // 15: bright white
	}
	if code >= 0 && code < len(colors) {
		return colors[code]
	}
	return "#000000"
}

// styleMapToArray converte map de estilos para array ordenado
func styleMapToArray(styleMap map[string]bool) []string {
	styles := []string{"bold", "italic", "underline", "strikethrough", "faint", "blink", "reverse"}
	active := []string{}
	for _, s := range styles {
		if styleMap[s] {
			active = append(active, s)
		}
	}
	return active
}
