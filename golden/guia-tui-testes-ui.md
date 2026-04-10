# Guia de Testes de UI com Golden Files — `internal/tui`

> Guia prático: como criar, rodar, atualizar e depurar testes visuais do pacote `tui`.
> Para entender a arquitetura e as decisões por trás do sistema, consulte `arquitetura-tui-testes-ui.md`.

---

## Sumário

- [Pré-requisitos](#pré-requisitos)
- [Estrutura mínima de um teste golden](#estrutura-mínima-de-um-teste-golden)
- [Múltiplos cenários no mesmo teste](#múltiplos-cenários-no-mesmo-teste)
- [Componente 2D: quando usar `WxH` no nome](#componente-2d-quando-usar-wxh-no-nome)
- [Criar golden files pela primeira vez](#criar-golden-files-pela-primeira-vez)
- [Rodar os testes](#rodar-os-testes)
- [Entender uma falha](#entender-uma-falha)
- [Atualizar golden files após mudança visual intencional](#atualizar-golden-files-após-mudança-visual-intencional)
- [Testar invariantes de estilo sem golden file](#testar-invariantes-de-estilo-sem-golden-file)
- [Erros comuns e como resolver](#erros-comuns-e-como-resolver)

---

## Pré-requisitos

O helper `checkOrUpdateGolden` e a variável `update` estão em `messages_test.go`. O import do parser ANSI é necessário em qualquer arquivo que use `ParseANSIStyle`:

```go
import (
    testdatapkg "github.com/useful-toys/abditum/internal/tui/testdata"
)
```

Como todos os `*_test.go` do pacote `tui` compilam em um único binário de teste, os helpers declarados em `messages_test.go` (`goldenPath`, `checkOrUpdateGolden`, `stripANSI`, `update`) já estão disponíveis em qualquer outro arquivo `*_test.go` do pacote — sem nenhum import adicional.

---

## Estrutura mínima de um teste golden

Todo teste golden do projeto segue este padrão exato:

```go
func TestMeuComponente_Golden(t *testing.T) {
    // 1. Construir o modelo em estado determinístico
    m := &meuComponente{titulo: "Exemplo"}
    m.Init()
    m.theme = ThemeTokyoNight
    m.SetSize(80, 24)

    // 2. Renderizar
    out := m.View()

    // 3. Nível 1 — layout: texto puro sem ANSI
    txtPath := goldenPath("meucomponente", "initial", 80, "txt")
    checkOrUpdateGolden(t, txtPath, stripANSI(out))

    // 4. Nível 2 — estilos: transições de cor e fonte
    transitions := testdatapkg.ParseANSIStyle(out)
    jsonBytes, err := testdatapkg.MarshalStyleTransitions(transitions)
    if err != nil {
        t.Fatalf("marshal transitions: %v", err)
    }
    jsonPath := goldenPath("meucomponente", "initial", 80, "json")
    checkOrUpdateGolden(t, jsonPath, string(jsonBytes))
}
```

**Regras invariáveis:**
- `SetSize` é obrigatório antes de `View()` para garantir dimensões determinísticas
- `stripANSI(out)` vai para `.txt.golden`; `ParseANSIStyle(out)` (sem strip) vai para `.json.golden`
- Os dois arquivos são sempre criados em par para o mesmo cenário

---

## Múltiplos cenários no mesmo teste

Quando o mesmo componente tem variações de estado, use table-driven com `t.Run`. Note o `tc := tc` para captura correta de variável de loop com `t.Parallel()` ou sub-testes independentes:

```go
func TestRenderMeuComponente_Golden(t *testing.T) {
    type testCase struct {
        variant string
        estado  string // ou qualquer campo que muda entre cenários
    }

    cases := []testCase{
        {"vazio", ""},
        {"preenchido", "algum valor"},
        {"erro", "valor inválido"},
    }
    widths := []int{30, 80}

    for _, tc := range cases {
        for _, w := range widths {
            tc := tc // captura de variável de loop
            w := w
            t.Run(fmt.Sprintf("%s-%d", tc.variant, w), func(t *testing.T) {
                m := &meuComponente{valor: tc.estado}
                m.Init()
                m.theme = ThemeTokyoNight
                m.SetSize(w, 24)

                out := m.View()

                txtPath := goldenPath("meucomponente", tc.variant, w, "txt")
                checkOrUpdateGolden(t, txtPath, stripANSI(out))

                transitions := testdatapkg.ParseANSIStyle(out)
                jsonBytes, err := testdatapkg.MarshalStyleTransitions(transitions)
                if err != nil {
                    t.Fatalf("marshal transitions: %v", err)
                }
                jsonPath := goldenPath("meucomponente", tc.variant, w, "json")
                checkOrUpdateGolden(t, jsonPath, string(jsonBytes))
            })
        }
    }
}
```

Isso produz automaticamente os arquivos:
- `testdata/golden/meucomponente-vazio-30.txt.golden`
- `testdata/golden/meucomponente-vazio-30.json.golden`
- `testdata/golden/meucomponente-vazio-80.txt.golden`
- … e assim por diante para todas as combinações.

---

## Componente 2D: quando usar `WxH` no nome

Componentes com uma única linha renderizada (message bar, command bar) usam apenas a largura no nome: `messages-success-30.txt.golden`.

Componentes com altura variável ou scrollável (decision dialog, filepicker, help modal) codificam as duas dimensões na variante. Nesse caso, use `decisionGoldenPath` (ou `helpGoldenPath`) em vez de `goldenPath`, pois o tamanho já está embutido no nome da variante:

```go
// variante já inclui dimensões: "destructive-1action-short-60x24"
checkOrUpdateDecisionGolden(t, "destructive-1action-short-60x24", "txt", stripANSI(out))
checkOrUpdateDecisionGolden(t, "destructive-1action-short-60x24", "json", string(jsonBytes))
```

Para novos componentes 2D, defina os helpers equivalentes no arquivo de teste do componente seguindo o mesmo padrão de `decision_test.go`.

---

## Criar golden files pela primeira vez

Ao escrever um novo teste golden, rode-o normalmente. Os arquivos são criados automaticamente na primeira execução sem falha:

```sh
# Cria os golden files do novo teste
go test ./internal/tui/... -run TestMeuComponente_Golden -v
```

Verifique o conteúdo antes de commitar — o que está sendo persistido é o baseline que todas as execuções futuras vão comparar:

```sh
# Inspecionar layout
cat internal/tui/testdata/golden/meucomponente-initial-80.txt.golden

# Inspecionar estilos
cat internal/tui/testdata/golden/meucomponente-initial-80.json.golden
```

O `.txt.golden` deve ser legível como texto — bordas, conteúdo, espaçamento visíveis sem escapes ANSI. O `.json.golden` deve ser um array de tuplas `[linha, col, fg, bg, [estilos]]`, uma por linha.

Confirme que o visual está correto, depois commite ambos os arquivos golden junto com o código do teste.

---

## Rodar os testes

```sh
# Roda toda a suite do pacote tui (inclui todos os golden)
go test ./internal/tui/...

# Roda apenas os testes golden de um componente
go test ./internal/tui/... -run TestRenderMessageBar_Golden

# Roda um sub-teste específico (e.g. apenas success a 30 colunas)
go test ./internal/tui/... -run "TestRenderMessageBar_Golden/success-30"

# Com output verboso — mostra cada sub-teste PASS/FAIL
go test ./internal/tui/... -run TestDecisionDialog_Golden -v

# Roda todos os golden de todos os componentes (função de prefixo)
go test ./internal/tui/... -run "_Golden"

# Roda o parser ANSI também
go test ./internal/tui/testdata/...
```

---

## Entender uma falha

Quando um golden test falha, a saída segue o formato:

```
--- FAIL: TestRenderMessageBar_Golden/success-30 (0.00s)
    messages_test.go:NNN: golden mismatch for testdata/golden/messages-success-30.txt.golden:
    want:
    ╭──────────────────────────────╮
    │ ✓  Cofre salvo               │
    ╰──────────────────────────────╯
    got:
    ╭──────────────────────────────╮
    │ ✓ Cofre salvo                │
    ╰──────────────────────────────╯
```

**Para diagnosticar:**

1. Identifique se a diferença está no `.txt.golden` (layout/espaçamento) ou no `.json.golden` (cor/estilo)
2. Para `.txt.golden`: procure espaços, caracteres de borda, wrapping, ou conteúdo textual diferente
3. Para `.json.golden`: procure transições em linhas/colunas diferentes, ou atributos de cor/fonte incorretos
4. Compare com o arquivo atual em disco para entender a intenção original:

```sh
cat internal/tui/testdata/golden/messages-success-30.txt.golden
```

Se a diferença for intencional (nova feature, refactor visual), atualize os goldens (ver seção abaixo). Se for uma regressão inadvertida, corrija o código.

---

## Atualizar golden files após mudança visual intencional

Após modificar intencionalmente a aparência de um componente, regenere os baselines com o flag `-update`:

```sh
# Regenerar todos os golden do pacote
go test ./internal/tui/... -update

# Regenerar apenas um componente (mais rápido)
go test ./internal/tui/... -run TestMeuComponente_Golden -update

# Regenerar apenas um sub-teste
go test ./internal/tui/... -run "TestMeuComponente_Golden/preenchido-80" -update
```

Após regenerar, **revise o diff** antes de commitar:

```sh
# Ver todas as mudanças nos golden files
git diff internal/tui/testdata/golden/

# Ver apenas mudanças de layout (txt)
git diff internal/tui/testdata/golden/*.txt.golden

# Ver apenas mudanças de estilo (json)
git diff internal/tui/testdata/golden/*.json.golden
```

**Checklist antes de commitar goldens atualizados:**
- [ ] O diff de `.txt.golden` reflete a mudança de layout que era esperada
- [ ] O diff de `.json.golden` reflete as mudanças de cor/estilo que eram esperadas
- [ ] Nenhum arquivo golden foi alterado além do componente que mudou
- [ ] O commit inclui tanto o código modificado quanto os golden files atualizados

---

## Testar invariantes de estilo sem golden file

Quando o estado de estilo depende de uma condição de runtime (e.g. botão habilitado vs. desabilitado), usar um golden file por estado seria excessivo. Nesses casos, use o parser ANSI diretamente para inspecionar o output:

```go
func TestMeuComponente_BotaoAtivo(t *testing.T) {
    // Estado com campo preenchido — botão deve estar em bold
    m := &meuComponente{valor: "preenchido"}
    m.Init()
    m.theme = ThemeTokyoNight
    m.SetSize(80, 24)

    out := m.View()
    transitions := testdatapkg.ParseANSIStyle(out)

    // Encontrar a última linha renderizada
    maxLine := 0
    for _, tr := range transitions {
        if tr.Line > maxLine {
            maxLine = tr.Line
        }
    }

    // Verificar que existe bold na área de ação (últimas 3 linhas)
    temBold := false
    for _, tr := range transitions {
        if tr.Line >= maxLine-2 {
            for _, s := range tr.Style {
                if s == "bold" {
                    temBold = true
                }
            }
        }
    }
    if !temBold {
        t.Error("esperado bold na área de ação quando campo preenchido")
    }
}

func TestMeuComponente_BotaoDesabilitado(t *testing.T) {
    // Estado com campo vazio — botão NÃO deve estar em bold
    m := &meuComponente{valor: ""}
    m.Init()
    m.theme = ThemeTokyoNight
    m.SetSize(80, 24)

    out := m.View()
    transitions := testdatapkg.ParseANSIStyle(out)

    maxLine := 0
    for _, tr := range transitions {
        if tr.Line > maxLine {
            maxLine = tr.Line
        }
    }

    for _, tr := range transitions {
        if tr.Line >= maxLine-2 {
            for _, s := range tr.Style {
                if s == "bold" {
                    t.Error("bold inesperado na área de ação quando campo vazio")
                }
            }
        }
    }
}
```

Outros atributos verificáveis via `StyleTransition`:
- `tr.FG` — ponteiro para string hex do foreground, e.g. `"#9ece6a"`; use `*tr.FG == "#f7768e"` para verificar cor específica
- `tr.BG` — ponteiro para string hex do background
- `tr.Style` — slice com qualquer combinação de `"bold"`, `"italic"`, `"underline"`, `"strikethrough"`, `"faint"`, `"blink"`, `"reverse"`

---

## Erros comuns e como resolver

### `golden mismatch` inesperado em CI, mas passa localmente

Causa mais comum: o golden foi gerado com uma versão diferente do lipgloss ou com uma locale diferente. 

Verifique se o arquivo em disco está commitado corretamente:

```sh
git status internal/tui/testdata/golden/
```

Se houver arquivos não commitados (gerados localmente mas não incluídos no commit anterior), commite-os:

```sh
git add internal/tui/testdata/golden/
git commit -m "tui: adicionar/atualizar golden files de meucomponente"
```

### Arquivo `.txt.golden` com caracteres de controle visíveis

Causa: `stripANSI` foi aplicado antes que o output estivesse completo, ou a regex não cobre todas as sequências. Verifique se o output de `View()` não contém sequências além de SGR (e.g. sequências de cursor `\x1b[?25h`). O `stripANSI` cobre apenas `\x1b[[0-9;]*m` — o padrão SVG do projeto não produz outras sequências.

### `.json.golden` está vazio (`[]`)

Causa: o componente está renderizando texto puro sem nenhum código ANSI — provavelmente o tema não foi injetado ou `SetSize` não foi chamado antes de `View()`. Verifique:

```go
m.theme = ThemeTokyoNight  // obrigatório
m.SetSize(80, 24)          // obrigatório antes de View()
```

### `panic: runtime error` ao chamar `View()` sem `Init()`

Quase todo componente do projeto inicializa campos internos em `Init()`. Sempre chame `m.Init()` antes de `m.View()` em testes.

### Golden test cria arquivo mas o teste falha na rodada seguinte

Causa: o conteúdo gerado não é determinístico — provavelmente `time.Now()`, um timestamp, ou um número aleatório está no output. Testes golden requerem que `View()` produza sempre o mesmo output dado o mesmo estado. Substitua valores não-determinísticos por campos fixos no estado do modelo para testes.
