# Especificação Visual — Cabeçalho

> Contexto global: nome da app, cofre, indicador dirty e abas de modo.
> Parte da [Especificação Visual](tui-specification.md).
>
> **Documento de fundação:**
> - [`tui-design-system.md`](tui-design-system.md) — fundações visuais

## Cabeçalho

**Responsabilidade:** contexto global — qual aplicação, qual cofre, se há alterações pendentes e qual modo está ativo na área de trabalho.
**Posição:** linhas 1–2 da tela (zona Cabeçalho do [DS — Dimensionamento](tui-design-system.md#dimensionamento-e-layout)).
**Altura fixa:** 2 linhas.

### Anatomia

| Linha | Papel | Conteúdo |
|---|---|---|
| **1 — Título** | Contexto + navegação | Nome da app, `·` separador, nome do cofre, `•` dirty, abas de modo à direita |
| **2 — Separadora** | Divisa cabeçalho ↔ área de trabalho | Linha `─` full-width; a aba ativa "pousa" nesta linha via `╯ Texto ╰` |

**Três estados estruturais:**

| Estado | Linha 1 | Linha 2 | Abas |
|---|---|---|---|
| Sem cofre (boas-vindas) | Apenas nome da app | Separador simples, sem conectores | Ocultas |
| Cofre aberto | Nome da app `·` cofre `•` + abas | Separador com aba ativa suspensa | Visíveis (3) |
| Busca ativa | Nome da app `·` cofre `•` + abas | Campo de busca à esquerda + aba ativa suspensa à direita | Visíveis (3) |

> Wireframes ilustrativos a 80 colunas. A largura real acompanha o terminal.

#### Sem cofre (Boas-vindas)

```
  Abditum
──────────────────────────────────────────────────────────────────────────────────
```

Sem nome de cofre, sem indicador dirty, sem abas. A linha separadora é contínua.

#### Cofre aberto — anatomia base

> Estado impossível em operação normal (sempre há um modo ativo). Mostrado para ilustrar a posição de todos os elementos antes de qualquer aba estar ativa.

**Sem alterações:**

```
  Abditum · cofre                          ╭ Cofre ╮  ╭ Modelos ╮  ╭ Config ╮
──────────────────────────────────────────────────────────────────────────────────
```

**Com alterações não salvas:**

```
  Abditum · cofre •                         ╭ Cofre ╮  ╭ Modelos ╮  ╭ Config ╮
──────────────────────────────────────────────────────────────────────────────────
```

O `•` aparece imediatamente após o nome do cofre, em `semantic.warning`. Desaparece após salvamento bem-sucedido.

#### Modo Cofre ativo

```
  Abditum · cofre •                      ╭───────╮  ╭ Modelos ╮  ╭ Config ╮
─────────────────────────────────────────╯ Cofre ╰──────────────────────────────
```

A aba ativa na linha 1 substitui o texto por `─` (`╭───────╮`), mantendo a mesma largura da versão inativa (`╭ Cofre ╮`). Na linha 2, o texto desce para o gap entre `╯` e `╰`, que se alinham verticalmente com `╭` e `╮` da linha 1 — conectando visualmente a aba à área de trabalho abaixo.

#### Modo Modelos ativo

```
  Abditum · cofre                          ╭ Cofre ╮  ╭─────────╮  ╭ Config ╮
──────────────────────────────────────────────────────╯ Modelos ╰────────────────
```

#### Modo Configurações ativo

```
  Abditum · cofre                           ╭ Cofre ╮  ╭ Modelos ╮  ╭────────╮
────────────────────────────────────────────────────────────────────╯ Config ╰──
```

A aba mais à direita pode encostar na borda do terminal — `╰` ocupa a última coluna, sem `─` posterior.

> **Nota:** a aba Configurações é referida como "Config" nos wireframes por economia de espaço. O texto completo na implementação é `Config`.

#### Modo busca ativo

Ativo enquanto o campo de busca estiver aberto (ver [Busca de Segredos](tui-spec-arvore.md#busca-de-segredos)). Disponível apenas no Modo Cofre com cofre aberto.

A linha separadora (linha 2) é substituída pelo campo de busca. A aba ativa permanece suspensa à direita na mesma linha, sem alteração de posição ou estilo.

**Campo aberto, sem query (recém-ativado):**

```
  Abditum · cofre •                      ╭───────╮  ╭ Modelos ╮  ╭ Config ╮
 ─ Busca: ────────────────────────────────╯ Cofre ╰──────────────────────────
```

**Campo aberto, com query:**

```
  Abditum · cofre •                      ╭───────╮  ╭ Modelos ╮  ╭ Config ╮
 ─ Busca: gmail ──────────────────────────╯ Cofre ╰──────────────────────────
```

> **Exceção de layout documentada:** a linha separadora do cabeçalho tem papel estrutural fixo no DS (divisa cabeçalho ↔ área de trabalho). Durante o modo busca, essa linha assume papel adicional de display do campo de busca. Exceção justificada pelo princípio **Hierarquia da Informação** — o campo imediatamente acima da árvore cria relação visual direta entre query e resultado — e pelo princípio **O Terminal como Meio** — espaço vertical é recurso escasso. Escopo-limitada ao Modo Cofre com busca ativa.

#### Mecânica visual da aba ativa

A transformação de aba inativa → ativa ocorre em duas linhas simultâneas:

| Linha | Aba inativa | Aba ativa |
|---|---|---|
| **1** | `╭ Texto ╮` (borda + texto) | `╭──────╮` (borda + preenchimento `─`) |
| **2** | `─────────` (separador contínuo) | `╯ Texto ╰` (gap com texto sobre `special.highlight`) |

Regras de alinhamento:

- A largura total da aba é **idêntica** nos estados ativo e inativo
- `╯` alinha-se verticalmente com `╭` da linha acima
- `╰` alinha-se verticalmente com `╮` da linha acima
- O conteúdo entre `╯` e `╰` (espaço + texto + espaço) tem fundo `special.highlight`
- As bordas `╭╮╯╰` e o preenchimento `─` usam sempre `border.default`, independente do estado

### Identidade Visual

| Elemento | Token | Atributo |
|---|---|---|
| `Abditum` (nome da app) | `accent.primary` | **bold** |
| `·` separador nome/cofre | `border.default` | — |
| Nome do cofre (radical, sem `.abditum`) | `text.secondary` | — |
| `•` indicador não salvo | `semantic.warning` | — |
| Bordas das abas (`╭╮╯╰─`) — ativa e inativa | `border.default` | — |
| Aba ativa — fundo (gap entre `╯` e `╰`) | `special.highlight` | — |
| Aba ativa — texto | `accent.primary` | **bold** |
| Aba inativa — texto | `text.secondary` | — |
| Linha separadora | `border.default` | — |
| `─ Busca: ` rótulo (modo busca) | `border.default` | — |
| Texto da query (modo busca) | `accent.primary` | **bold** |
| `─` preenchimento (modo busca) | `border.default` | — |

### Dimensionamento

**Truncamento do nome do cofre:**

O espaço disponível para o nome do cofre é limitado — as abas ocupam largura fixa à direita. O componente calcula o espaço em tempo real.

> **Extensão `.abditum` é omitida** — a app só trabalha com este formato, então a extensão é redundante. O nome exibido é o radical do arquivo (ex: `cofre.abditum` → `cofre`).

**Fórmula:**

```
prefixo  = "  Abditum · "                             (12 colunas)
dirty    = " •"  se IsDirty(), ou ""                   (2 ou 0 colunas)
abas     = bloco de abas + espaços entre elas           (largura fixa, ~32 colunas)
padding  = mín. 1 coluna entre nome/dirty e abas

disponível = largura_terminal − prefixo − dirty − abas − padding
```

**Algoritmo:**

1. Se o nome completo (radical sem extensão) cabe → exibir como está
2. Se não cabe → truncar com `…`: `{nome[0..n]}…` onde `n` é calculado para caber
3. Se nem 1 caractere + `…` (2 colunas) cabe → exibir apenas `…`

**Prioridade de cessão de espaço:**

| Prioridade | Elemento | Comportamento |
|---|---|---|
| 1ª (cede primeiro) | Nome do cofre | Truncado conforme algoritmo acima |
| 2ª | Separador `·` e indicador `•` | Preservados enquanto houver espaço |
| 3ª (nunca cede) | Abas | Largura fixa, nunca truncadas |

**Wireframe — nome truncado (terminal ~80 colunas, modo Cofre):**

```
  Abditum · meu-cofre-pessoa… •          ╭───────╮  ╭ Modelos ╮  ╭ Config ╮
─────────────────────────────────────────╯ Cofre ╰──────────────────────────────
```

O radical `meu-cofre-pessoal` foi truncado para `meu-cofre-pessoa…`.

**Campo de busca na linha separadora:**

| Elemento | Largura | Notas |
|---|---|---|
| `─ Busca: ` (prefixo fixo) | 9 colunas | `─` + espaço + `Busca:` + espaço |
| Texto da query | variável | Em `accent.primary` **bold** |
| `─` preenchimento | restante − largura da aba ativa − 2 (margem direita mínima) | Preenche até a aba |
| Aba ativa (`╯ Texto ╰`) | igual ao estado normal | Posição e estilo inalterados |

- **Query longa:** truncada à **esquerda** com `…` — a parte mais recente da query fica sempre visível
- A largura disponível para a query é calculada em tempo real e recalculada a cada resize do terminal

### Eventos

| Evento | Mudança visual |
|---|---|
| Cofre aberto com sucesso | Aparece `·` nome do cofre e as 3 abas |
| Cofre fechado / bloqueado | Desaparece nome do cofre e abas; volta ao estado boas-vindas |
| Alteração em memória (`IsDirty() = true`) | Aparece `•` em `semantic.warning` |
| Salvamento bem-sucedido (`IsDirty() = false`) | Desaparece `•` |
| Navegação entre modos (Cofre / Modelos / Config) | Aba ativa muda; nova aba suspensa na linha separadora |
| Terminal redimensionado | Nome do cofre recalcula truncamento |

### Comportamento

- **Abas clicáveis** — mouse troca o modo ativo ao clicar no texto ou na borda da aba (área de hit inclui linhas 1 e 2 da aba)
- **Navegação por teclado** — `F2` Cofre, `F3` Modelos, `F4` Config (escopo Área de trabalho — só ativas com cofre aberto)
- **Indicador dirty** — aparece/desaparece imediatamente conforme `IsDirty()`, sem animação
- **Truncamento dinâmico** — recalculado a cada renderização (resize do terminal, mudança de modo ativo, cofre aberto/fechado)
