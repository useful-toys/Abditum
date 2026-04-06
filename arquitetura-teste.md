# Arquitetura de Testes de View para Componentes de Tela

## Objetivo

Validar que `View()` de componentes (modais, diálogos, telas) renderizam **texto, espaçamento, bordas e cores** corretos, detectando regressões visuais sutis (faltas de espaço, cores erradas, alinhamentos incorretos).

Essa arquitetura é **genérica** e aplicável a qualquer componente que implemente `View() string`. Exemplos: `DecisionDialog`, `HelpDialog`, `PromptDialog`, futuras modais ou telas full-screen.

### Escopo

Golden files validam **apenas o conteúdo do componente retornado por `View()`**, não o entorno do terminal ou o contexto externo:

```
╭── ⚠  Excluir segredo ─────────────────────────────────────╮
│                                                              │
│  Gmail será excluído permanentemente. Esta ação não pode    │
│  ser desfeita.                                               │
│                                                              │
╰── Enter Excluir ───────────────────────── Esc Cancelar ──╯
```

**Exemplo (DecisionDialog):** O `View()` retorna exatamente essas linhas. O posicionamento na tela (centrado, canto, padding externo) é responsabilidade de `rootModel.Render()` via `lipgloss.Place()` e não é testado aqui — seria imprevisível variar de terminal para terminal.

Para componentes full-screen (sem envolvimento em `Place()`), `View()` retorna o conteúdo completo esperado na tela — todo o espaço reservado.

## Problema

Testes unitários convencionais (`strings.Contains`, asserções de propriedade) não capturam regressões visuais:
- Espaçamento errado entre elementos
- Cores incorretas em posições específicas
- Alinhamento quebrado em wrapping de texto
- Caracteres de borda ou estrutura desalinhados

**Exemplo real:** O bug de espaçamento em `DecisionDialog.renderActionBar()` passou em todos os testes existentes porque nenhum verificava se havia espaço após o último token de ação. O output visualmente ficava "backup──────" em vez de "backup ──────".

## Solução: Golden Files com Dois Níveis

### Nível 1: `.txt.golden` — Layout Visual

**Arquivo:** `{package}/testdata/golden/{component}-{variant}-{size}.txt.golden`

**Conteúdo:** Output ANSI cru de `View()` — o conteúdo exato que o componente renderiza.

**Valida:**
- Caracteres de borda, símbolos, ícones
- Espaçamento exato entre elementos
- Wrapping de texto e alinhamento
- Estrutura visual completa

**Não valida:**
- Padding ou entorno externo (responsabilidade de quem chama `View()`)
- Posicionamento relativo à tela (ex: centralizado via `lipgloss.Place()`)

**Exemplos:**
```
# Para DecisionDialog (modal)
decision-destructive-1action-80x24.txt.golden
decision-destructive-2action-80x24.txt.golden
decision-destructive-3action-80x24.txt.golden
decision-error-1action-80x24.txt.golden
decision-neutral-2action-30x25.txt.golden

# Para HelpDialog (futuro)
help-default-80x24.txt.golden
help-long-content-50x25.txt.golden

# Para PromptDialog (futuro)
prompt-text-input-60x5.txt.golden
prompt-password-input-60x5.txt.golden
```

O padrão é: `{component}-{variant}-{size}.txt.golden`
- `component`: nome do tipo/componente (`decision`, `help`, `prompt`)
- `variant`: diferenciador (severidade, tipo, estilo, etc.)
- `size`: terminal size (`80x24`, `30x25`, etc.)

### Nível 2: `.json.golden` — Estilo por Posição

**Arquivo:** `{package}/testdata/golden/{component}-{variant}-{size}.json.golden`

**Conteúdo:** Lista de tuplas registrando **transições visuais de estilo** — APENAS onde a cor/fonte **muda realmente**.

**IMPORTANTE:** O `.json.golden` é **canônico, compacto e agnóstico a SGR**.

#### Regra de Ouro
Registrar tupla **APENAS quando o estado `(fg, bg, style)` muda**:
- Se `(fg, bg, style)` são idênticos ao estado anterior → **não registrar**
- Se uma linha começa em estado padrão (`#default`, `null`, `[]`) → **não registrar** a tupla de inicialização, a menos que haja mudança naquela linha
- Registrar tupla a cada mudança real

#### Agnóstico a SGR
Independente de quantos `\x1b[...m` (SGR codes) existirem no output ANSI, ou sua ordem/posição na string — registra APENAS o resultado visual final.

Exemplos de equivalência (TODOS geram mesma tupla):
```
\x1b[1m\x1b[38;5;208m          → [linha, col, "#d08c00", null, ["bold"]]
\x1b[38;5;208m\x1b[1m          → [linha, col, "#d08c00", null, ["bold"]]
\x1b[38;2;208;140;0m\x1b[1m    → [linha, col, "#d08c00", null, ["bold"]]
```

Todos produzem **visualmente**: texto laranja bold. O `.json.golden` registra exatamente isso — a cor hex normalizada + o style.

#### Formato de cada tupla

```
[linha, coluna, fg_hex|null, bg_hex|null, style_array]
```

- **linha** — número da linha no output (0-indexed)
- **coluna** — posição do caractere onde o estilo **muda visualmente**
- **fg_hex** — cor do texto em hex (ex: `"#a0a0a0"`, `"#d08c00"`), ou `null` se default terminal
- **bg_hex** — cor de fundo em hex, ou `null` se não aplicado
- **style_array** — array de strings com estilos ativos: `["bold"]`, `["italic", "underline"]`, `[]` para nenhum

#### Exemplo: DecisionDialog 3-action (80×24)

Saída visual:
```
╭── ⚠  Excluir segredo ─────────────────────────────────────╮
│                                                              │
│  Gmail será excluído permanentemente. Esta ação não pode    │
│  ser desfeita.                                               │
╰── Enter Excluir ───────────────────────── Esc Cancelar ──╯
```

Golden file (apenas transições reais):
```json
[
  [0, 0, "#a0a0a0", null, []],
  [0, 4, "#d08c00", null, ["bold"]],
  [0, 24, "#a0a0a0", null, []],
  [4, 4, "#d08c00", null, ["bold"]],
  [4, 14, "#a0a0a0", null, []]
]
```

**Explicação:**
- `[0, 0, ...]`: Linha 0 inicia com cinza (#a0a0a0), sem bold — bordas (`╭──`)
- `[0, 4, ...]`: Na coluna 4, **muda** para laranja bold — começa o título (`⚠  Excluir segredo`)
- `[0, 24, ...]`: Na coluna 24, **volta** a cinza sem bold — fim do título, dashes de preenchimento
- `[4, 4, ...]`: Linha 4 (ações). Coluna 4 **muda** para laranja bold — primeira ação (`Enter Excluir`)
- `[4, 14, ...]`: Coluna 14 **volta** a cinza — fim da ação, dashes e próxima ação em cinza

**O que NÃO registra:**
- `[1, 0, ...]`, `[2, 0, ...]`, `[3, 0, ...]`: Linhas 1-3 começam em estado padrão e não têm mudanças → omitidas
- `[4, 0, ...]`: Linha 4 começa em estado padrão → não registra (mas há mudança em coluna 4, que é registrada)

#### Validação

Valida:
- Cor correta em cada posição de mudança visual
- Estilos (bold, italic, underline, etc.) ativos/inativos corretamente
- Nenhuma cor faltando ou desordenada

Não valida:
- Ordem ou quantidade de SGRs no output ANSI — só o resultado visual
- Estados redundantes (idênticos ao anterior)

**VSCode Integration:** Cores hex são renderizadas com color swatch na gutter — legibilidade visual direta. Abra o arquivo `.json.golden` no VSCode e veja as cores reais ao lado das tuplas.

---

## Matriz de Cobertura

A matriz é **específica para cada componente**. Esta fase cobre **apenas Camada 1** — golden files visuais (`.txt.golden` + `.json.golden`) em **3 tamanhos de terminal** (30, 60, 80 cols) + testes de `Update()`.

### DecisionDialog — Golden Files (Camada 1)

**5 severidades × 3 layouts × 3 tamanhos = 45 pares**

| Severidade | 1-action (30/60/80) | 2-action (30/60/80) | 3-action (30/60/80) |
|-----------|---------------------|---------------------|---------------------|
| Destructive | ✓✓✓ | ✓✓✓ | ✓✓✓ |
| Error | ✓✓✓ | ✓✓✓ | ✓✓✓ |
| Alert | ✓✓✓ | ✓✓✓ | ✓✓✓ |
| Informative | ✓✓✓ | ✓✓✓ | ✓✓✓ |
| Neutral | ✓✓✓ | ✓✓✓ | ✓✓✓ |

**Total:** 45 pares `.txt.golden` + `.json.golden` = 90 arquivos

### DecisionDialog — Testes de Update

**8 cenários de estado** (Enter/Esc/tecla explícita/tecla desconhecida × variações de layout)

### HelpDialog — Golden Files (Camada 1)

**2 variantes × 3 tamanhos = 6 pares**

| Variante | 30 cols | 60 cols | 80 cols |
|----------|---------|---------|---------|
| Poucas ações (sem scroll) | ✓ | ✓ | ✓ |
| Muitas ações (com scroll) | ✓ | ✓ | ✓ |

**Total:** 6 pares = 12 arquivos

### HelpDialog — Testes de Update

**8 cenários de estado** (Esc/F1/Up/Down/PgUp/PgDown/Home/End × boundary conditions)

### Linha de Mensagem — Golden Files (Camada 1)

**6 kinds × 3 tamanhos = 18 pares**

| Kind | 30 cols | 60 cols | 80 cols |
|------|---------|---------|---------|
| Success | ✓ | ✓ | ✓ |
| Info | ✓ | ✓ | ✓ |
| Warn | ✓ | ✓ | ✓ |
| Error | ✓ | ✓ | ✓ |
| Busy (frame 0) | ✓ | ✓ | ✓ |
| Hint | ✓ | ✓ | ✓ |

**Total:** 18 pares = 36 arquivos

### Linha de Ações — Golden Files (Camada 1)

**3 cenários × 3 tamanhos = 9 pares**

| Cenário | 30 cols | 60 cols | 80 cols |
|---------|---------|---------|---------|
| Sem ações | ✓ | ✓ | ✓ |
| 2 ações | ✓ | ✓ | ✓ |
| 5 ações + F1 anchor | ✓ | ✓ | ✓ |

**Total:** 9 pares = 18 arquivos

### Total Geral — Camada 1

| Componente | Golden Pairs | Arquivos | Testes Update |
|-----------|-------------|----------|---------------|
| DecisionDialog | 45 | 90 | 8 cenários |
| HelpDialog | 6 | 12 | 8 cenários |
| Linha de Mensagem | 18 | 36 | — (função pura) |
| Linha de Ações | 9 | 18 | — (função pura) |
| **Total** | **78 pares** | **156 arquivos** | **16 cenários** |

### Por que 3 tamanhos (30, 60, 80)

| Tamanho | O que captura |
|---------|--------------|
| **30 cols** | Extremo apertado — wrapping agressivo, truncamento de texto, tokens de ação comprimidos, box width mínimo |
| **60 cols** | Intermediário — wrapping moderado, layout balanceado, comportamento típico de terminais médios |
| **80 cols** | Padrão — layout "ideal", sem wrapping, espaçamento completo |

Cada par de golden files em tamanho diferente é uma **fotografia do layout naquele ponto exato**. Se uma mudança no código alterar o wrapping, o spacing, ou a responsividade, o golden file daquele tamanho específico será o primeiro a quebrar — dando um feedback preciso sobre onde o problema ocorreu.

---

## Implementação

### 1. Estrutura de Diretórios

**Padrão genérico:**
```
{package}/testdata/golden/
├── {component}-{variant}-{size}.txt.golden
├── {component}-{variant}-{size}.json.golden
├── ...
```

**Exemplo para DecisionDialog:**
```
internal/tui/testdata/golden/
├── decision-destructive-1action-80x24.txt.golden
├── decision-destructive-1action-80x24.json.golden
├── decision-destructive-2action-80x24.txt.golden
├── decision-destructive-2action-80x24.json.golden
├── decision-destructive-3action-80x24.txt.golden
├── decision-destructive-3action-80x24.json.golden
├── decision-error-1action-80x24.txt.golden
├── decision-error-1action-80x24.json.golden
├── ... (15 pares severidade×layout)
├── decision-destructive-1action-30x25.txt.golden
├── decision-destructive-1action-30x25.json.golden
├── ... (9 pares layout×tamanho)
```

### 2. Parser SGR Mínimo (Reutilizável)

Código genérico em `internal/tui/testdata/parser.go` ou similar (~100 linhas):

```go
// StyleTransition representa uma transição de estilo (cor/fonte muda)
type StyleTransition struct {
    Line   int      // número da linha (0-indexed)
    Col    int      // coluna onde muda (0-indexed)
    FG     string   // cor foreground em hex, ou "" se default
    BG     string   // cor background em hex, ou "" se default
    Style  []string // ["bold", "italic", ...], ou []
}

// ParseANSIStyle extrai transições de estilo do output ANSI
// Registra APENAS onde o estado (FG, BG, Style) muda realmente
func ParseANSIStyle(output string) []StyleTransition { ... }
```

Responsabilidades:
- Caminha caractere por caractere pelo output
- Quando encontra `\x1b[...m` (SGR), atualiza estado (fg, bg, bold, italic, etc.)
- Quando o estado **muda** (comparado ao anterior), registra uma tupla
- Normaliza codes 256-color e truecolor para hex
- Conta linhas (quebras `\n`) e colunas
- **Não registra** transições para estado idêntico ao anterior

### 3. Golden Test Runner (Genérico)

Padrão reutilizável em testes de cada componente:

```go
func TestComponentViewGolden(t *testing.T) {
    tests := []struct {
        name      string
        component tea.Model // ou interface específica com View()
        golden    string    // ex: "destructive-1action-80x24"
    }{
        // Preencher conforme matriz do componente
        {"Destructive 1-action", pocKey1(), "decision-destructive-1action-80x24"},
        {"Destructive 2-action", pocKey2(), "decision-destructive-2action-80x24"},
        ...
        {"Layout 1-action 30x25", pocKey1_30x25(), "decision-destructive-1action-30x25"},
        ...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            out := tt.component.View()
            
            // Nível 1: .txt.golden
            txtPath := filepath.Join("testdata/golden", tt.golden+".txt.golden")
            if *flagUpdate {
                os.WriteFile(txtPath, []byte(out), 0644)
            } else {
                expected, _ := os.ReadFile(txtPath)
                if string(expected) != out {
                    t.Errorf("View() output mismatch:\nexpected:\n%s\n\ngot:\n%s", expected, out)
                }
            }
            
            // Nível 2: .json.golden
            styles := ParseANSIStyle(out)
            jsonPath := filepath.Join("testdata/golden", tt.golden+".json.golden")
            jsonBytes, _ := json.MarshalIndent(styles, "", "  ")
            if *flagUpdate {
                os.WriteFile(jsonPath, jsonBytes, 0644)
            } else {
                expectedJSON, _ := os.ReadFile(jsonPath)
                if string(expectedJSON) != string(jsonBytes) {
                    t.Errorf("Style mismatch:\nexpected:\n%s\n\ngot:\n%s", expectedJSON, jsonBytes)
                }
            }
        })
    }
}
```

### 4. Testes de `Update()` — Exploração de Estado

Além de `View()`, o método `Update()` de cada componente interativo deve ser testado explorando **todas as transições de estado possíveis**.

#### DecisionDialog

**Teclas válidas:**
- `Enter` → executa ação com `Default: true`, retorna cmd non-nil
- `Esc` → executa ação com `Cancel: true`, retorna cmd non-nil
- Tecla explícita (ex: `M`, `A`) → executa ação correspondente, retorna cmd non-nil
- Tecla desconhecida → nil (nenhuma ação)

**Matriz de testes Update:**
| Entrada | Estado inicial | Resultado esperado |
|---------|---------------|-------------------|
| `Enter` | 1-action (só Default) | cmd != nil, ação Default executada |
| `Enter` | 2-action (Default + Cancel) | cmd != nil, ação Default executada |
| `Enter` | 3-action (Default + Neutral + Cancel) | cmd != nil, ação Default executada |
| `Esc` | 2-action | cmd != nil, ação Cancel executada |
| `Esc` | 1-action (sem Cancel) | cmd == nil (só pop) |
| `M` | ação com Key="M" | cmd != nil, ação M executada |
| `X` | nenhuma ação com Key="X" | nil (ignorada) |
| `Tab` | qualquer | nil (não existe foco por Tab) |

#### HelpDialog

**Teclas válidas:**
- `Esc` / `F1` → popModal
- `Up` / `Down` → scroll ±1 (clamp 0..max)
- `PgUp` / `PgDown` → scroll ±contentHeight
- `Home` / `End` → scroll 0 / max

**Matriz de testes Update:**
| Entrada | Estado inicial | Resultado esperado |
|---------|---------------|-------------------|
| `Esc` | scroll=0 | popModal |
| `Up` | scroll=0 | scroll=0 (clamp) |
| `Up` | scroll=5 | scroll=4 |
| `Down` | scroll=max | scroll=max (clamp) |
| `Down` | scroll=0 | scroll=1 |
| `PgDown` | scroll=0 | scroll=contentHeight |
| `Home` | scroll=max | scroll=0 |
| `End` | scroll=0 | scroll=max |

#### MessageBar e CommandBar

São **funções puras** (`RenderMessageBar(msg, width)` e `RenderCommandBar(width)`) — não têm `Update()`. Seus testes são cobertos inteiramente pelos golden files de `View()`.

### 5. Flag de Regeneração

```bash
# Em package de testes, adicionar flag global:
var flagUpdate = flag.Bool("update", false, "regenerate golden files")

# Usar em cada teste golden
```

Rodagem:
```bash
go test ./internal/tui/... -update
```

Regenera ambos os golden files (`.txt.golden` + `.json.golden`) quando há mudanças intencionais.

---

## Workflow

### Desenvolvimento normal (sem mudanças visuais)
```bash
go test ./internal/tui/...
# Compara output contra golden files
# Se não bater: teste falha, mostra diff
```

### Após mudança visual intencional (ex: mudar cor, spacing)
```bash
# Revisa mudança no código
# Regenera golden files:
go test ./internal/tui/... -update

# Valida visualmente nos arquivos:
# - .txt.golden: abre no terminal, vê layout
# - .json.golden: abre no VSCode, vê cores com swatches

# Commita ambos os arquivos:
git add testdata/golden/
git commit -m "test: update golden files after visual change"
```

### Bug detectado (ex: espaçamento faltando)
```bash
# Teste falha com diff claro:
# "expected: 'Excluir ──'  got: 'Excluir──'"

# Corrige o bug em decision.go
go test ./internal/tui/... -update  # regen golden
git add internal/tui/decision.go testdata/golden/
git commit -m "fix: add trailing space in renderActionBar"
```

---

## Estilos Suportados

Do lipgloss v2, os estilos inline capturados:

| Estilo | Chave JSON | Uso típico |
|--------|-----------|-----------|
| Bold | `"bold"` | Títulos, teclas default |
| Italic | `"italic"` | Notas, hints |
| Underline | `"underline"` | Links, atalhos (sem style/color) |
| Strikethrough | `"strikethrough"` | Itens removidos |
| Faint | `"faint"` | Texto secundário desabilitado |
| Blink | `"blink"` | Alertas (raro em TUI moderno) |
| Reverse | `"reverse"` | Seleção invertida |

**Nota:** `Underline` pode ter sub-atributos (`underline_style`: Single, Double, Curly, Dotted, Dashed; `underline_color`). Inicialmente registramos só `"underline"` como booleano. Se necessário, expandir para objetos.

---

## Benefícios

| Aspecto | Benefício |
|--------|-----------|
| **Texto** | Captura espaçamento, wrapping, alinhamento exatos |
| **Cores** | Valida cores certas por severidade |
| **Fonts** | Bold, italic, underline ativos/inativos |
| **Regressões** | Qualquer mudança quebra teste → revisão consciente |
| **Manutenção** | Golden files legíveis (texto + JSON com color swatches) |
| **Integração** | `-update` flag simples para regenerar |

---

## Casos de Uso

### 1. Regressão de espaçamento (como o bug atual)
```
Encontrado em .txt.golden:
- expected: "Abrir backup ──────"
+ got:      "Abrir backup──────"
→ Teste falha imediatamente
```

### 2. Regressão de cor
```
Encontrado em .json.golden:
- expected: [0, 4, "#d08c00", null, ["bold"]]
+ got:      [0, 4, "#a0a0a0", null, ["bold"]]
→ Teste falha, diff visual no VSCode (color swatches)
```

### 3. Mudança intencional de severidade (ex: novo estilo Warning)
```
Dev muda ColorWarn em colors.go
Tests quebram (esperado)
Dev roda: go test -update
VSCode mostra novas cores nos .json.golden
Dev valida visualmente e commita
```

---

## Escopo Futuro

- **HelpDialog:** Aplicar mesma arquitetura (golden files para layout + cores)
- **PromptDialog:** Golden files com diferentes variantes de input
- **Futuros modais/telas:** Usar parser SGR genérico, adaptar matriz de cobertura
- **Propriedades estruturais:** Adicionar invariantes (ex: "cada linha tem exatamente N colunas", "bordas sempre presentes")
- **Screenshot tests:** Se necessário capturar imagens PNG para validação visual completa


