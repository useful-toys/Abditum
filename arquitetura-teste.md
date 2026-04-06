# Arquitetura de Testes Golden para Componentes TUI

## Objetivo

Validar que `View()` de componentes TUI (modais, diálogos, linhas de status) renderizam **texto, espaçamento, bordas e cores** corretos, detectando regressões visuais sutis que testes convencionais não conseguem capturar.

Esta arquitetura é **genérica** e aplicável a qualquer componente que implemente uma função de renderização que retorne `string` (preferencialmente ANSI).

### Aplicabilidade

Exemplos de componentes que se beneficiam desta arquitetura:
- `DecisionDialog` (modal de decisão)
- `HelpDialog` (modal de ajuda)
- `RenderMessageBar` (linha de mensagem)
- `RenderCommandBar` (linha de ações/comandos)
- Futuros modais, telas full-screen, ou elementos de UI independentes.

### Escopo dos Golden Files

Golden files validam **apenas o conteúdo visual do componente retornado pela sua função de renderização**, não o entorno do terminal ou o contexto externo (como o componente é posicionado na tela). O posicionamento é responsabilidade de quem chama (ex: `rootModel.Render()` via `lipgloss.Place()`) e não é testado aqui, pois seria imprevisível variar de terminal para terminal.

**Exemplo (DecisionDialog):** O `View()` retorna apenas o box do diálogo, ilustrando os elementos visuais que serão validados por esta arquitetura de testes:

```
╭── ⚠  Excluir segredo ─────────────────────────────────────╮
│                                                              │
│  Gmail será excluído permanentemente. Esta ação não pode    │
│  ser desfeita.                                               │
│                                                              │
╰── Enter Excluir ───────────────────────── Esc Cancelar ──╯
```

Neste exemplo de `DecisionDialog`:
- O **título** (`⚠ Excluir segredo`) é renderizado com **fonte bold** e sua cor (vermelho no caso de "⚠") e a cor da borda superior (`╭──`) se alinham à **severidade da mensagem**.
- As **ações** na barra inferior (`Enter Excluir`, `Esc Cancelar`) são **destacadas** (com suas respectivas teclas em bold e cor de destaque) e corretamente espaçadas.

Esta arquitetura visa validar que **evoluções no código não causarão regressões** nesses aspectos visuais cruciais, garantindo a consistência da UI textual.

Para componentes full-screen (sem envolvimento em `lipgloss.Place()`), a função de renderização retorna o conteúdo completo esperado na tela — todo o espaço reservado.

## Problema: Limitações dos Testes Atuais

Testes unitários convencionais, baseados em `strings.Contains()` ou asserções de propriedades simples, são insuficientes para garantir a fidelidade visual de interfaces textuais. Eles não capturam regressões como:
- Espaçamento incorreto entre elementos (`"texto──"` vs `"texto ──"`)
- Cores ou estilos errados em posições específicas
- Alinhamento quebrado ou wrapping de texto inesperado
- Caracteres de borda ou símbolos desalinhados

**Exemplo real:** O bug de espaçamento na `DecisionDialog.renderActionBar()` (commit `6325967`) passou em todos os testes existentes porque nenhum verificava a presença de um espaço crucial. O output visualmente ficava `"backup──────"` em vez de `"backup ──────"`.

## Solução: Golden Files com Dois Níveis de Validação

Para capturar tanto a estrutura visual quanto o estilo de forma robusta e canônica, utilizamos uma abordagem de dois arquivos golden por caso de teste. Cada arquivo possui uma responsabilidade distinta e complementar:

1.  **`.txt.golden` (Estrutura Visual Canônica):** Valida o layout, espaçamento e conteúdo textual, removendo toda a formatação ANSI para focar no "desenho" puro da tela.
2.  **`.json.golden` (Estilo Visual Canônico):** Valida cores, estilos de fonte (bold, itálico), etc., de forma agnóstica aos códigos ANSI exatos, focando no resultado visual da formatação.

### Nível 1: `.txt.golden` — O Desenho Canônico do Layout

**Arquivo:** `{package}/testdata/golden/{component}-{variant}-{size}.txt.golden`

**Conteúdo:** A representação **canônica da estrutura visual** do componente. Contém **apenas texto visível**, com todos os códigos de escape ANSI (cores, bold, etc.) removidos. O resultado é um arquivo de texto puro, legível por humanos, que isola o "desenho" da tela da sua formatação.

```
╭── ⚠  Excluir segredo ─────────────────────────────────────╮
│                                                              │
│  Gmail será excluído permanentemente. Esta ação não pode    │
│  ser desfeita.                                               │
│                                                              │
╰── Enter Excluir ───────────────────────── Esc Cancelar ──╯
```

**Propósito:** Validar a integridade do layout, o posicionamento dos elementos e a consistência do texto.

**Valida (O quê e Onde):**
- Caracteres de borda, símbolos e ícones.
- Espaçamento exato entre todos os elementos visíveis.
- Quebras de linha (wrapping) e alinhamento do texto.
- Estrutura geral, incluindo o número total de linhas e a largura visual de cada linha.

**NÃO Valida (Como):**
- **Cores:** Uma mudança na cor de um elemento (ex: de `Error` para `Success`) **não** quebraria este teste, pois a cor é removida.
- **Estilos de Fonte:** A aplicação ou remoção de `bold`, `italic`, `underline`, etc., **não** quebraria este teste.
- *A validação de estilo é responsabilidade exclusiva do Nível 2: `.json.golden`.*

**Exemplo de Regressão que este arquivo PEGA:**
O bug de espaçamento em `DecisionDialog.renderActionBar()` (`"backup──────"` em vez de `"backup ──────"`) teria sido capturado imediatamente, pois a sequência de caracteres visíveis é diferente.

**Convenção de Nomes para Golden Files:**
O padrão de nomes segue `{component}-{variant}-{size}.txt.golden` ou `.json.golden`:
- `component`: Nome do tipo/componente (ex: `decision`, `help`, `messages`, `commandbar`, `prompt`).
- `variant`: Um diferenciador para o cenário (ex: `destructive-1action`, `f1only`, `success`, `long-content`).
- `size`: As dimensões do terminal (ex: `80x24`, `30x25`). Para componentes de uma linha, apenas a largura (ex: `80`, `30`).

*Exemplos de nomes de arquivo para diferentes componentes:*
- **DecisionDialog:** `decision-destructive-1action-80x24.txt.golden`, `decision-error-2action-30x25.txt.golden`.
- **HelpDialog:** `help-fewactions-60x12.txt.golden`, `help-15actions-bottom-30x16.txt.golden`.
- **MessageBar:** `messages-success-60.txt.golden`, `messages-error-30.txt.golden`.
- **CommandBar:** `commandbar-typical-60.txt.golden`, `commandbar-many-30.txt.golden`.
- **PromptDialog (futuro):** `prompt-text-input-60x5.txt.golden`, `prompt-password-input-60x5.txt.golden`.


### Nível 2: `.json.golden` — O Estilo Canônico por Posição

**Arquivo:** `{package}/testdata/golden/{component}-{variant}-{size}.json.golden`

**Conteúdo:** A representação **canônica do estilo visual** do componente. É uma lista de tuplas que registram **transições de estilo visual**, ou seja, os pontos exatos (linha, coluna) onde a cor do foreground, background ou qualquer atributo de fonte (bold, italic, etc.) muda.

Este arquivo é gerado por um parser SGR que analisa o output ANSI **original** da função de renderização e o converte para uma forma padronizada. Ele é **agnóstico aos códigos ANSI brutos** — significa que diferentes sequências de escape SGR que produzem o mesmo efeito visual são normalizadas para a mesma tupla canônica.

**Propósito:** Validar a aplicação correta de cores e estilos, garantindo que a formatação visual esteja de acordo com a especificação, complementando a validação estrutural do `.txt.golden`.

**Valida (Como):**
- **Cores corretas** em cada posição de mudança (ex: a cor de `Error` é de fato a cor vermelha especificada).
- **Estilos de Fonte corretos** (bold, italic, underline, etc.) são aplicados e removidos nas posições exatas.

**NÃO Valida (O quê e Onde):**
- Espaçamento, alinhamento ou conteúdo textual (responsabilidade do Nível 1: `.txt.golden`).

**Exemplo de Regressão que este arquivo PEGA:**
- Mudar a cor de um título de erro de vermelho para cinza quebraria este teste, mas não o `.txt.golden`.
- Remover o `bold` de uma tecla de atalho quebraria este teste, mas não o `.txt.golden`.

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

#### Convenção de Nomes para Golden Files
O padrão de nomes segue `{component}-{variant}-{size}.json.golden`:
- `component`: Nome do tipo/componente (ex: `decision`, `help`, `messages`, `commandbar`, `prompt`).
- `variant`: Um diferenciador para o cenário (ex: `destructive-1action`, `f1only`, `success`, `long-content`).
- `size`: As dimensões do terminal (ex: `80x24`, `30x25`). Para componentes de uma linha, apenas a largura (ex: `80`, `30`).

*Exemplos de nomes de arquivo para diferentes componentes:*
- **DecisionDialog:** `decision-destructive-1action-80x24.json.golden`, `decision-error-2action-30x25.json.golden`.
- **HelpDialog:** `help-fewactions-60x12.json.golden`, `help-15actions-bottom-30x16.json.golden`.
- **MessageBar:** `messages-success-60.json.golden`, `messages-error-30.json.golden`.
- **CommandBar:** `commandbar-typical-60.json.golden`, `commandbar-many-30.json.golden`.
- **PromptDialog (futuro):** `prompt-text-input-60x5.json.golden`, `prompt-password-input-60x5.json.golden`.

---

## Princípios de Cobertura

A definição da matriz de testes para cada componente deve seguir os seguintes princípios arquiteturais para garantir robustez sem cair na exaustividade de combinações cartesianas.

### 1. Tamanhos de Terminal
A preferência é testar com **2 tamanhos de largura**:
*   **30 colunas:** Cenário restrito (truncamento, wrapping agressivo).
*   **60 colunas:** Cenário padrão (layout ideal).
*   *Nota:* Outros tamanhos podem ser discutidos para casos com necessidades específicas de validação.

### 2. Variantes de Construção
Recomenda-se testar **todas as variantes de construção** da view definidas pelo componente.
*   *Exemplo:* No diálogo, testar para cada severidade.

### 3. Variantes de Popular o Modelo
Além das variantes de construção, deve-se testar **diferentes formas de popular o modelo** com dados para cobrir casos de uso complexos.
*   *Exemplo:* Na barra de ações, testar com várias actions, cada qual configurada de forma peculiar (ex: teclas especiais como F1, labels longos, ações desabilitadas).
*   *Exemplo:* No diálogo de ajuda, popular com grupos de ações variados para forçar scroll e testar a renderização de cabeçalhos de grupo.

### 4. Diversidade de Conteúdo e Estados do Modelo
Deve-se explorar a renderização com diferentes tipos de conteúdo textual e diferentes estados do modelo para validar o layout dinâmico e interativo:
*   **Títulos:** Curtos e longos (para testar truncamento ou ajuste).
*   **Corpo do Texto:**
    *   Linhas curtas.
    *   Múltiplas linhas (com separadores explícitos).
    *   Linhas longas que forçam a quebra automática (wrapping).
    *   Conteúdo extenso que força o surgimento de scroll.
*   **Estados de Scroll:** Para componentes com scroll, testar a representação visual em diferentes posições: início (topo), meio (indicador de posição) e fim (últimas linhas visíveis).
*   **Foco/Seleção:** Foco em diferentes ações ou itens da lista (se aplicável ao componente).

### Estratégia de Combinação
Não deve explorar todas as combinações possíveis, mas deve fazer uma **boa cobertura**.
*   *Evitar:* Matriz cartesiana completa.
*   *Adotar:* Seleção inteligente de cenários onde cada teste valida múltiplas dimensões simultaneamente.

### Matriz de Exemplo (Conceitual)

Esta é uma matriz de exemplo **ilustrativa** para um componente genérico (`GenericComponent`), demonstrando os princípios de cobertura. As matrizes reais para cada componente devem ser detalhadas em seus respectivos documentos de contexto ou plano.

| # | Golden Name | Variant | Content Type | State | Width (30/60) | Captura |
|---|-------------|---------|--------------|-------|---------------|---------|
| 1 | `generic-default-short-30` | Default | Título Curto | Padrão | 30 | Layout básico, truncamento |
| 2 | `generic-default-long-60` | Default | Título Longo | Padrão | 60 | Wrapping, layout completo |
| 3 | `generic-error-multiline-30` | Error | Múltiplas Linhas | Padrão | 30 | Cor, símbolos, wrapping |
| 4 | `generic-warning-scroll-60` | Warning | Conteúdo extenso | Scroll Meio | 60 | Scrollbar, posição |

**Total Conceitual:** `N` pares de golden files (`.txt.golden` + `.json.golden`) + `M` cenários de `Update()`.

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

Padrão reutilizável em testes de cada componente, demonstrando a separação:

```go
func TestComponentViewGolden(t *testing.T) {
    // ... setup do loop de testes ...
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // 1. Gera o output ANSI completo uma vez
            outWithANSI := tt.component.View()
            
            // 2. Teste de Nível 1: Layout Canônico (.txt.golden)
            txtPath := filepath.Join("testdata/golden", tt.golden+".txt.golden")
            visibleText := stripANSI(outWithANSI) // Remove códigos de escape
            if *flagUpdate {
                os.WriteFile(txtPath, []byte(visibleText), 0644)
            } else {
                expected, _ := os.ReadFile(txtPath)
                if string(expected) != visibleText {
                    t.Errorf("Layout mismatch (.txt.golden):\nexpected:\n%s\n\ngot:\n%s", expected, visibleText)
                }
            }
            
            // 3. Teste de Nível 2: Estilo Canônico (.json.golden)
            styles := ParseANSIStyle(outWithANSI) // Processa o output original com ANSI
            jsonPath := filepath.Join("testdata/golden", tt.golden+".json.golden")
            jsonBytes, _ := json.MarshalIndent(styles, "", "  ")
            if *flagUpdate {
                os.WriteFile(jsonPath, jsonBytes, 0644)
            } else {
                expectedJSON, _ := os.ReadFile(jsonPath)
                if string(expectedJSON) != string(jsonBytes) {
                    t.Errorf("Style mismatch (.json.golden):\nexpected:\n%s\n\ngot:\n%s", expectedJSON, jsonBytes)
                }
            }
        })
    }
}

// stripANSI remove todas as sequências de escape ANSI de uma string.
func stripANSI(s string) string {
    re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
    return re.ReplaceAllString(s, "")
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

### 1. Regressão de espaçamento (como o bug de `renderActionBar`)
```
# Teste falha com diff claro no `.txt.golden`:
- expected: "Abrir backup ──────"
+ got:      "Abrir backup──────"
```
**Capturado por:** `.txt.golden` (validação estrutural).

### 2. Regressão de cor
```
# Teste falha com diff no `.json.golden`:
- expected: [0, 4, "#d08c00", null, ["bold"]]  // Laranja
+ got:      [0, 4, "#a0a0a0", null, ["bold"]]  // Cinza
```
**Capturado por:** `.json.golden` (validação de estilo). O `.txt.golden` passaria, pois o texto e o espaçamento não mudaram.

### 3. Mudança intencional de severidade (ex: novo estilo Warning)
```
# Dev muda ColorWarn em colors.go
# go test ./... falha em ambos os arquivos:
# - .txt.golden: quebra se o símbolo mudar (ex: ⚠ → !)
# - .json.golden: quebra porque a cor mudou de amarelo para a nova cor

# Dev regenera os arquivos:
go test ./internal/tui/... -update

# Dev valida a mudança visual nos arquivos gerados e commita.
```
**Capturado por:** Ambos. O `.txt.golden` pega a mudança do símbolo `⚠` e o `.json.golden` pega a mudança da cor.

---

## Escopo Futuro

- **HelpDialog:** Aplicar mesma arquitetura (golden files para layout + cores)
- **PromptDialog:** Golden files com diferentes variantes de input
- **Futuros modais/telas:** Usar parser SGR genérico, adaptar matriz de cobertura
- **Propriedades estruturais:** Adicionar invariantes (ex: "cada linha tem exatamente N colunas", "bordas sempre presentes")
- **Screenshot tests:** Se necessário capturar imagens PNG para validação visual completa


