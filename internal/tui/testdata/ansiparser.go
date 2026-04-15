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

// marshalStyleChanges serializa uma fatia de StyleTransition no formato legível:
// array externo com indentação, cada tupla em uma única linha compacta.
//
// Exemplo:
//
//	[
//	  [0, 2, null, null, ["bold", "faint"]],
//	  [0, 9, "#d08c00", null, []]
//	]
//
// Isso é equivalente a json.MarshalIndent mas com cada tupla colapsada em uma linha,
// tornando o arquivo .json.golden fácil de ler e diff-friendly.
func marshalStyleChanges(transitions []StyleTransition) ([]byte, error) {
	if transitions == nil {
		return []byte("null"), nil
	}
	if len(transitions) == 0 {
		return []byte("[]"), nil
	}
	var buf strings.Builder
	buf.WriteString("[\n")
	for i, t := range transitions {
		b, err := json.Marshal(t) // compacto por MarshalJSON acima
		if err != nil {
			return nil, err
		}
		buf.WriteString("  ")
		buf.Write(b)
		if i < len(transitions)-1 {
			buf.WriteByte(',')
		}
		buf.WriteByte('\n')
	}
	buf.WriteString("]")
	return []byte(buf.String()), nil
}

// ansiState holds the current rendering state during ANSI sequence parsing.
type ansiState struct {
	fg    *string // cor foreground em hex, ou nil
	bg    *string // cor background em hex, ou nil
	style map[string]bool
}

// ansiToStyleChanges extrai transições de estilo visual do output ANSI.
// Normaliza SGR codes para estado visual canônico (fg hex, bg hex, style set).
// Registra tupla APENAS quando o estado muda — ignora redundâncias e ordem de SGR.
func ansiToStyleChanges(output string) []StyleTransition {
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
		applyCodes(&currentState, codes)

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

// applyCodes processa uma sequência de códigos SGR, tratando corretamente
// sequências multi-código como 38;2;R;G;B (truecolor) e 38;5;N (256-color).
func applyCodes(state *ansiState, codes []string) {
	for i := 0; i < len(codes); i++ {
		code, _ := strconv.Atoi(strings.TrimSpace(codes[i]))
		switch code {
		case 0: // Reset
			state.fg = nil
			state.bg = nil
			state.style = make(map[string]bool)
		case 1:
			state.style["bold"] = true
		case 2:
			state.style["faint"] = true
		case 3:
			state.style["italic"] = true
		case 4:
			state.style["underline"] = true
		case 5:
			state.style["blink"] = true
		case 7:
			state.style["reverse"] = true
		case 9:
			state.style["strikethrough"] = true
		case 21, 22:
			state.style["bold"] = false
			state.style["faint"] = false
		case 23:
			state.style["italic"] = false
		case 24:
			state.style["underline"] = false
		case 25:
			state.style["blink"] = false
		case 27:
			state.style["reverse"] = false
		case 29:
			state.style["strikethrough"] = false
		case 30, 31, 32, 33, 34, 35, 36, 37:
			color := colorCode16(code - 30)
			state.fg = &color
		case 38: // Foreground: 38;5;N (256-color) ou 38;2;R;G;B (truecolor)
			i = applyExtendedColor(codes, i, &state.fg)
		case 39:
			state.fg = nil
		case 40, 41, 42, 43, 44, 45, 46, 47:
			color := colorCode16(code - 40)
			state.bg = &color
		case 48: // Background: 48;5;N (256-color) ou 48;2;R;G;B (truecolor)
			i = applyExtendedColor(codes, i, &state.bg)
		case 49:
			state.bg = nil
		case 90, 91, 92, 93, 94, 95, 96, 97:
			color := colorCode16(code - 90 + 8)
			state.fg = &color
		case 100, 101, 102, 103, 104, 105, 106, 107:
			color := colorCode16(code - 100 + 8)
			state.bg = &color
		}
	}
}

// applyExtendedColor processa 256-color (5;N) ou truecolor (2;R;G;B) a partir
// da posição i (que aponta para o código 38 ou 48). Retorna o novo índice i.
func applyExtendedColor(codes []string, i int, target **string) int {
	if i+1 >= len(codes) {
		return i
	}
	mode, _ := strconv.Atoi(strings.TrimSpace(codes[i+1]))
	switch mode {
	case 5: // 256-color: 38;5;N
		if i+2 < len(codes) {
			n, _ := strconv.Atoi(strings.TrimSpace(codes[i+2]))
			color := colorCode256(n)
			*target = &color
			return i + 2
		}
	case 2: // Truecolor: 38;2;R;G;B
		if i+4 < len(codes) {
			r, _ := strconv.Atoi(strings.TrimSpace(codes[i+2]))
			g, _ := strconv.Atoi(strings.TrimSpace(codes[i+3]))
			b, _ := strconv.Atoi(strings.TrimSpace(codes[i+4]))
			color := fmt.Sprintf("#%02x%02x%02x", r, g, b)
			*target = &color
			return i + 4
		}
	}
	return i + 1
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

// colorCode256 converte código ANSI 256-color para hex.
// Cobre as 3 faixas: 0-15 (16-color), 16-231 (cubo 6×6×6), 232-255 (grayscale).
func colorCode256(code int) string {
	if code < 16 {
		return colorCode16(code)
	}
	if code < 232 {
		// Cubo 6×6×6: código 16-231
		idx := code - 16
		b := idx % 6
		g := (idx / 6) % 6
		r := idx / 36
		// Cada componente mapeia para 0, 95, 135, 175, 215, 255
		levels := [6]int{0, 95, 135, 175, 215, 255}
		return fmt.Sprintf("#%02x%02x%02x", levels[r], levels[g], levels[b])
	}
	// Grayscale: código 232-255
	gray := 8 + (code-232)*10
	return fmt.Sprintf("#%02x%02x%02x", gray, gray, gray)
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
