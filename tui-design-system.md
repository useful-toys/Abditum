# Design System — Abditum TUI

> Definições visuais fundamentais para o pacote `internal/tui`.  
> Complementa `tui-design.md` (layout e interação) e `tui-elm-architecture.md` (arquitetura).

---

## Paleta de Cores

A paleta é organizada por **papel funcional**, não por nome de cor. Cada papel define *para que* a cor é usada — a cor concreta muda conforme o tema escolhido.

### Papéis funcionais

A TUI usa os seguintes papéis de cor:

| Categoria | Papel | Descrição |
|---|---|---|
| **Superfícies** | `surface.base` | Fundo principal da aplicação |
| | `surface.raised` | Fundo de painéis, modais, elementos elevados |
| | `surface.overlay` | Fundo de tooltips, menus flutuantes, overlays |
| **Texto** | `text.primary` | Texto principal — conteúdo, labels, títulos |
| | `text.secondary` | Texto auxiliar — descrições, placeholders, hints |
| | `text.disabled` | Texto desabilitado ou indisponível |
| **Bordas** | `border.default` | Bordas de painéis e separadores — estado normal |
| | `border.focused` | Borda do painel ou elemento com foco ativo |
| **Interação** | `accent.primary` | Cor principal de ação — elemento selecionado, cursor, destaque de foco |
| | `accent.secondary` | Cor secundária — informações complementares, links |
| **Semânticas** | `semantic.success` | Operação concluída, confirmação positiva |
| | `semantic.warning` | Atenção requerida, ação potencialmente perigosa |
| | `semantic.error` | Falha, ação destrutiva, erro |
| | `semantic.info` | Informação neutra, dica contextual |
| **Especiais** | `special.muted` | Itens apagados — marcados para exclusão, desabilitados |
| | `special.highlight` | Fundo de item selecionado em listas/árvore |

### Regras de aplicação

- **Texto sobre superfícies:** `text.primary` sobre `surface.base` deve ter contraste mínimo legível. Em TUI, isso é garantido naturalmente pela paleta (fundo escuro + texto claro).
- **Bordas indicam foco:** `border.focused` é a única indicação visual de qual painel está ativo — deve ser claramente distinta de `border.default`.
- **Semânticas são reservadas:** cores semânticas aparecem somente para comunicar estado (sucesso, erro, etc.) — nunca como decoração.
- **Consistência entre contextos:** a mesma cor semântica é usada em mensagens (`MsgWarn`), modais (`DialogAlert`), e demais elementos com o mesmo significado.

---

## Proposta A: Tokyo Night

Baseada na paleta [Tokyo Night](https://github.com/enkia/tokyo-night-vscode-theme) — tema escuro com tons predominantemente azuis e roxos. Projetada para conforto em uso prolongado: fundo azul-noite (não preto puro), texto acinzentado-azulado, destaques vibrantes mas dessaturados.

### Superfícies

| Papel | Hex | Swatch | Nota |
|---|---|---|---|
| `surface.base` | `#1a1b26` | <span style="background:#1a1b26;color:#1a1b26">██</span> | Azul-noite profundo — menos agressivo que preto puro |
| `surface.raised` | `#24283b` | <span style="background:#24283b;color:#24283b">██</span> | Elevação sutil — painéis, modais |
| `surface.overlay` | `#414868` | <span style="background:#414868;color:#414868">██</span> | Overlays — diferenciação clara do fundo |

### Texto

| Papel | Hex | Swatch | Nota |
|---|---|---|---|
| `text.primary` | `#a9b1d6` | <span style="color:#a9b1d6">██</span> | Azul-acinzentado claro — confortável para leitura |
| `text.secondary` | `#565f89` | <span style="color:#565f89">██</span> | Cinza muted — hints, descrições, placeholders |
| `text.disabled` | `#3b4261` | <span style="color:#3b4261">██</span> | Quase invisível — itens indisponíveis |

### Bordas

| Papel | Hex | Swatch | Nota |
|---|---|---|---|
| `border.default` | `#414868` | <span style="color:#414868">██</span> | Cinza-azulado — separadores, bordas sem foco |
| `border.focused` | `#7aa2f7` | <span style="color:#7aa2f7">██</span> | Azul vibrante — painel ativo, campo em edição |

### Interação

| Papel | Hex | Swatch | Nota |
|---|---|---|---|
| `accent.primary` | `#7aa2f7` | <span style="color:#7aa2f7">██</span> | Azul — cursor, item selecionado, ação principal |
| `accent.secondary` | `#bb9af7` | <span style="color:#bb9af7">██</span> | Lilás — informação complementar, decoração sutil |

### Semânticas

| Papel | Hex | Swatch | Nota |
|---|---|---|---|
| `semantic.success` | `#9ece6a` | <span style="color:#9ece6a">██</span> | Verde suave — confirmação, operação ok |
| `semantic.warning` | `#e0af68` | <span style="color:#e0af68">██</span> | Amarelo quente — bloqueio iminente, ação irreversível |
| `semantic.error` | `#f7768e` | <span style="color:#f7768e">██</span> | Rosa-avermelhado — falha, exclusão |
| `semantic.info` | `#7dcfff` | <span style="color:#7dcfff">██</span> | Ciano — informação neutra, reconhecimento |

### Especiais

| Papel | Hex | Swatch | Nota |
|---|---|---|---|
| `special.muted` | `#565f89` | <span style="color:#565f89">██</span> | Cinza — itens marcados para exclusão, desabilitados |
| `special.highlight` | `#283457` | <span style="background:#283457;color:#a9b1d6">██</span> | Azul escuro — fundo de item selecionado em listas |

### Gradiente do logo

| Linha | Hex | Swatch |
|---|---|---|
| 1 | `#9d7cd8` | <span style="color:#9d7cd8">██</span> |
| 2 | `#89ddff` | <span style="color:#89ddff">██</span> |
| 3 | `#7aa2f7` | <span style="color:#7aa2f7">██</span> |
| 4 | `#7dcfff` | <span style="color:#7dcfff">██</span> |
| 5 | `#bb9af7` | <span style="color:#bb9af7">██</span> |

### Personalidade

Sóbria, profissional, confortável. Transmite confiança e calma — adequada para uma ferramenta de segurança. Cores dessaturadas reduzem fadiga visual em sessões longas. O azul-noite como fundo evita o preto puro, que pode parecer "portal para o vazio" em terminais grandes.

---

## Proposta B: Cyberpunk

Inspirada na estética cyberpunk/synthwave — fundo muito escuro com acentos neon vibrantes. Alta saturação nos destaques, contraste dramático. Cores quentes (rosa, amarelo) dominam a interação, com ciano elétrico como contraponto frio.

### Superfícies

| Papel | Hex | Swatch | Nota |
|---|---|---|---|
| `surface.base` | `#0a0a1a` | <span style="background:#0a0a1a;color:#0a0a1a">██</span> | Preto-azulado profundo — noite digital |
| `surface.raised` | `#1a1a2e` | <span style="background:#1a1a2e;color:#1a1a2e">██</span> | Elevação com tom roxo sutil |
| `surface.overlay` | `#2a2a3e` | <span style="background:#2a2a3e;color:#2a2a3e">██</span> | Modais e overlays — violeta escuro |

### Texto

| Papel | Hex | Swatch | Nota |
|---|---|---|---|
| `text.primary` | `#e0e0ff` | <span style="color:#e0e0ff">██</span> | Lavanda claro — brilhante, futurístico |
| `text.secondary` | `#8888aa` | <span style="color:#8888aa">██</span> | Lilás apagado — hints, descrições |
| `text.disabled` | `#444466` | <span style="color:#444466">██</span> | Roxo escuro — quase fundido ao fundo |

### Bordas

| Papel | Hex | Swatch | Nota |
|---|---|---|---|
| `border.default` | `#3a3a5c` | <span style="color:#3a3a5c">██</span> | Roxo-acinzentado — separadores discretos |
| `border.focused` | `#ff2975` | <span style="color:#ff2975">██</span> | Rosa neon — foco impossível de ignorar |

### Interação

| Papel | Hex | Swatch | Nota |
|---|---|---|---|
| `accent.primary` | `#ff2975` | <span style="color:#ff2975">██</span> | Rosa-magenta neon — ação principal, cursor |
| `accent.secondary` | `#00fff5` | <span style="color:#00fff5">██</span> | Ciano elétrico — contraponto frio, informação |

### Semânticas

| Papel | Hex | Swatch | Nota |
|---|---|---|---|
| `semantic.success` | `#05ffa1` | <span style="color:#05ffa1">██</span> | Verde neon — brilhante, inequívoco |
| `semantic.warning` | `#ffe900` | <span style="color:#ffe900">██</span> | Amarelo elétrico — alerta visualmente urgente |
| `semantic.error` | `#ff3860` | <span style="color:#ff3860">██</span> | Vermelho quente — falha, perigo |
| `semantic.info` | `#00b4d8` | <span style="color:#00b4d8">██</span> | Ciano médio — informação, reconhecimento |

### Especiais

| Papel | Hex | Swatch | Nota |
|---|---|---|---|
| `special.muted` | `#666688` | <span style="color:#666688">██</span> | Lilás desbotado — itens apagados |
| `special.highlight` | `#2a1533` | <span style="background:#2a1533;color:#e0e0ff">██</span> | Magenta muito escuro — fundo de seleção |

### Gradiente do logo

| Linha | Hex | Swatch |
|---|---|---|
| 1 | `#ff2975` | <span style="color:#ff2975">██</span> |
| 2 | `#b026ff` | <span style="color:#b026ff">██</span> |
| 3 | `#00fff5` | <span style="color:#00fff5">██</span> |
| 4 | `#05ffa1` | <span style="color:#05ffa1">██</span> |
| 5 | `#ff2975` | <span style="color:#ff2975">██</span> |

### Personalidade

Ousada, energética, high-tech. Transmite poder e modernidade — como um terminal de hacker em filme de ficção científica. A alta saturação dos neons chama atenção mas pode causar fadiga em uso prolongado. O rosa neon como cor de foco é incomum e memorável, mas polarizante.

---

## Comparação

| Critério | Tokyo Night | Cyberpunk |
|---|---|---|
| **Conforto prolongado** | Excelente — dessaturada, tons frios e suaves | Moderado — neons cansam em sessões longas |
| **Legibilidade** | Alta — texto `#a9b1d6` sobre `#1a1b26` é equilibrado | Alta — texto `#e0e0ff` sobre `#0a0a1a` tem mais contraste |
| **Distinção semântica** | Clara — cores suficientemente distintas entre si | Muito clara — alta saturação torna diferenças óbvias |
| **Profissionalismo** | Alta — sóbria, familiar a devs (VS Code, IDEs) | Baixa — estética de entretenimento, pode parecer lúdica |
| **Adequação ao domínio** | Forte — ferramenta de segurança pede sobriedade | Fraca — neons contrastam com a seriedade de um cofre de senhas |
| **Expressividade do logo** | Elegante — gradiente suave violeta→ciano | Impactante — gradiente neon rosa→ciano→verde |
| **Acessibilidade** | Boa — contraste suficiente sem ser agressivo | Risco — neons podem ser problemáticos para sensibilidade visual |

---

## Decisão

> **Em aberto.** A decisão será tomada após avaliar ambas as propostas visualmente na implementação (Phase 5 stubs com cores reais no terminal).

Independente da escolha, a abstração por **papéis funcionais** garante que trocar de paleta é uma operação isolada — mudar os valores hex em um único arquivo de estilos, sem alterar lógica.

---

## Tipografia

Em TUI não existem fontes nem tamanhos — o terminal usa fonte monoespaçada fixa. Os "pesos tipográficos" disponíveis são atributos ANSI que o lipgloss expõe: **bold**, *italic*, dim, underline e ~~strikethrough~~.

### Atributos e quando usá-los

| Atributo | Efeito visual | Uso no Abditum |
|---|---|---|
| **Bold** | Texto mais brilhante e/ou espesso | Títulos de painéis, opção default em modais, label de campo em foco, nome da aplicação no header |
| *Italic* | Texto inclinado (suporte varia por terminal) | Hints (`MsgHint`), placeholders, descrições contextuais |
| Dim | Texto com brilho reduzido | Itens desabilitados, texto secundário quando `text.secondary` não for suficiente |
| Underline | Sublinhado | Reservado — uso pontual se necessário (ex: link ou atalho em texto corrido) |
| ~~Strikethrough~~ | Texto riscado | Segredos marcados para exclusão na árvore |
| Normal | Sem atributo | Corpo de texto, valores de campos, itens de lista |

### Combinações

Atributos podem ser combinados. Combinações previstas:

| Combinação | Uso |
|---|---|
| Bold + cor semântica | Título de modal com `DialogType` (ex: bold amarelo para `DialogAlert`) |
| Dim + strikethrough | Item marcado para exclusão e desabilitado simultaneamente |
| Italic + `text.secondary` | Hints e descrições — itálico reforça o caráter auxiliar |

### Princípios

- **Bold é o único destaque forte.** Usar com moderação — se tudo for bold, nada é destaque.
- **Dim é o oposto de bold.** Indica que o elemento existe mas não é relevante no momento.
- **Italic indica conteúdo auxiliar** — não é o dado em si, é uma explicação *sobre* o dado.
- **Strikethrough tem significado semântico único** — "marcado para remoção". Não usar decorativamente.
- **Underline é reserva.** Em TUI, underline pode ser confundido com cursor ou link. Evitar uso rotineiro.

---

## Bordas

Caracteres de box-drawing Unicode definem a linguagem visual de painéis, modais e separadores. Lipgloss oferece estilos predefinidos.

### Estilos disponíveis

| Estilo | Caracteres | Exemplo |
|---|---|---|
| Rounded | `╭ ╮ ╰ ╯ │ ─` | `╭──────╮`<br>`│      │`<br>`╰──────╯` |
| Single | `┌ ┐ └ ┘ │ ─` | `┌──────┐`<br>`│      │`<br>`└──────┘` |
| Double | `╔ ╗ ╚ ╝ ║ ═` | `╔══════╗`<br>`║      ║`<br>`╚══════╝` |
| Thick | `┏ ┓ ┗ ┛ ┃ ━` | `┏━━━━━━┓`<br>`┃      ┃`<br>`┗━━━━━━┛` |
| Hidden | espaços | Sem borda visível — apenas padding |

### Aplicação por elemento

| Elemento | Estilo sugerido | Cor da borda | Notas |
|---|---|---|---|
| Painel sem foco | Rounded | `border.default` | Cantos arredondados — visual moderno e leve |
| Painel com foco | Rounded | `border.focused` | Mesma forma, cor diferente — foco por cor, não por peso |
| Modal de confirmação | Rounded | Cor do `DialogType` | Borda colorida comunica tipo antes de ler conteúdo |
| Modal de senha | Rounded | `border.focused` | Neutro — a atenção é no campo, não na moldura |
| Modal de help | Rounded | `border.default` | Passivo — overlay informacional |
| Separador vertical (painéis) | `│` simples | `border.default` | Linha única entre árvore e detalhe |

### Princípios

- **Um único estilo de canto** (provavelmente Rounded) para consistência. Evitar misturar Single e Double na mesma interface.
- **Diferenciação por cor, não por estilo.** Painel ativo = mesma borda, cor diferente. Mais sutil que trocar de Single para Double.
- **Bordas são discretas.** O conteúdo é protagonista — bordas enquadram sem competir.
- **Título na borda.** Painéis e modais podem ter título integrado à borda superior (ex: `╭─ Cofre ───────╮`). Lipgloss suporta esse padrão.

---

## Ícones e Símbolos

Vocabulário de caracteres Unicode usados como ícones na interface. Usa-se Unicode básico (não Nerd Fonts) para máxima compatibilidade com terminais.

### Navegação na árvore

| Símbolo | Uso |
|---|---|
| `▸` | Pasta recolhida (U+25B8 BLACK RIGHT-POINTING SMALL TRIANGLE) |
| `▾` | Pasta expandida (U+25BE BLACK DOWN-POINTING SMALL TRIANGLE) |
| `·` | Segredo — item folha (U+00B7 MIDDLE DOT) |

### Estados de itens

| Símbolo | Uso |
|---|---|
| `★` | Favorito (U+2605 BLACK STAR) |
| `☆` | Não favorito — se necessário mostrar ambos (U+2606 WHITE STAR) |
| `✕` | Marcado para exclusão (U+2715 MULTIPLICATION X) |
| `•` | Alterações não salvas — indicador no header (U+2022 BULLET) |

### Mensagens (barra de mensagens)

| Símbolo | MsgKind | Uso |
|---|---|---|
| `✓` | `MsgInfo` | Sucesso (U+2713 CHECK MARK) — alternativa a ✅ se emoji não renderizar |
| `⚠` | `MsgWarn` | Atenção (U+26A0 WARNING SIGN) |
| `✗` | `MsgError` | Erro (U+2717 BALLOT X) — alternativa a ❌ |
| `◐ ◓ ◑ ◒` | `MsgBusy` | Spinner rotativo — 4 frames a 1fps |
| `•` | `MsgHint` | Hint (U+2022 BULLET) — alternativa a 💡 |

> **Emoji vs Unicode:** os emojis (`✅ ⚠️ ❌ 💡`) são visualmente mais ricos mas ocupam 2 colunas em muitos terminais e podem não renderizar em todos os ambientes. Os símbolos Unicode acima são fallback de 1 coluna. A decisão emoji vs Unicode será tomada na implementação com testes em terminais reais.

### Modais (tipo semântico)

| Símbolo | DialogType | Uso |
|---|---|---|
| `?` | `DialogQuestion` | Decisão neutra — alternativa a ❓ |
| `⚠` | `DialogAlert` | Ação destrutiva — mesmo símbolo do warning |
| `ℹ` | `DialogInfo` | Informação (U+2139 INFORMATION SOURCE) |

### Campos sensíveis

| Símbolo | Uso |
|---|---|
| `•` | Caractere de máscara de senha (U+2022 BULLET) — `••••••••` |
| `◉` | Campo revelável — indicador de que pode ser desmascarado (U+25C9 FISHEYE) |

### Scroll e navegação

| Símbolo | Uso |
|---|---|
| `↑` `↓` | Indicadores de scroll disponível (U+2191, U+2193) |
| `─` | Separador horizontal (U+2500 BOX DRAWINGS LIGHT HORIZONTAL) |
| `│` | Separador vertical (U+2502 BOX DRAWINGS LIGHT VERTICAL) |
| `…` | Texto truncado (U+2026 HORIZONTAL ELLIPSIS) |

### Princípios

- **Semântica antes de estética.** Cada símbolo tem um significado único — não reutilizar `★` para dois propósitos diferentes.
- **Fallback de 1 coluna.** Todo símbolo deve ter uma alternativa que ocupa exatamente 1 coluna de terminal, para layouts previsíveis.
- **Sem Nerd Fonts.** A TUI deve funcionar em qualquer terminal com suporte Unicode básico. Ícones elaborados (nerdfont glyphs) excluem usuários com configuração padrão.
- **Consistência com hierarquia tipográfica.** Símbolos complementam — bold para títulos, `★` para favorito, `✕` para exclusão. Nunca usar mais de um ícone por item.

---

## Estados Visuais

Definição de como elementos mudam visualmente conforme o estado de interação.

### Matriz de estados

| Estado | Cor do texto | Cor de fundo | Atributo | Borda | Exemplo |
|---|---|---|---|---|---|
| **Normal** | `text.primary` | `surface.base` | — | `border.default` | Item de lista, campo, painel inativo |
| **Focado** | `text.primary` | `surface.base` | — | `border.focused` | Painel ativo — borda muda de cor |
| **Selecionado (cursor)** | `text.primary` | `special.highlight` | **Bold** | — | Item sob cursor na árvore ou lista |
| **Ativo (pressionado)** | — | — | — | — | TUI não tem estado pressed |
| **Desabilitado** | `text.disabled` | `surface.base` | Dim | — | Ação indisponível na command bar |
| **Marcado para exclusão** | `special.muted` | `surface.base` | ~~Strikethrough~~ | — | Segredo com `✕` na árvore |
| **Favorito** | `text.primary` | `surface.base` | — | — | Item normal + `★` com `semantic.warning` ou `accent.secondary` |
| **Erro inline** | `semantic.error` | `surface.raised` | — | — | Mensagem de validação em modal de senha/texto |

### Transições

Em TUI, estados mudam **instantaneamente** — sem animação nem fade. A única animação é o spinner `MsgBusy` (1fps). Transições suaves não são viáveis em terminais.

### Princípios

- **Foco é por borda, seleção é por fundo.** Dois conceitos distintos: foco indica *qual painel* recebe input; seleção indica *qual item* dentro do painel é o alvo.
- **Nunca depender só de cor.** Itens marcados para exclusão usam cor + strikethrough + símbolo `✕`. Itens favoritos usam cor + símbolo `★`. Redundância garante legibilidade em terminais com cores limitadas.
- **Dim é preferível a hidden.** Itens desabilitados devem ser visíveis (dim) para que o usuário saiba que existem — invisibilidade causa confusão.

---

## Compatibilidade de Terminal

TUIs rodam em ambientes heterogêneos. O design system deve funcionar desde terminais modernos (24-bit color, todos os atributos) até terminais com capacidades reduzidas.

### Níveis de cor

| Nível | Cores | Terminais típicos | Suporte |
|---|---|---|---|
| **True Color (24-bit)** | 16 milhões | Windows Terminal, iTerm2, Alacritty, kitty, WezTerm, VS Code, GNOME Terminal (recente) | Alvo principal — hex exatos da paleta |
| **256 cores** | 216 cores + 24 cinzas | xterm-256color, tmux, Terminal.app (macOS), terminais SSH | Fallback obrigatório — cores mapeadas para o cubo 6×6×6 mais próximo |
| **16 cores (ANSI)** | 16 nomeadas | Consoles legados, SSH para servidores antigos, tty Linux | Fallback mínimo — funcional mas sem distinção fina |
| **Sem cor** | Monocromático | Pipes, redirecionamento, terminais muito antigos | Lipgloss desativa cor automaticamente (detecção via `$TERM` / `$NO_COLOR`) |

### Estratégia de fallback para 256 cores

As cores hex da paleta (Tokyo Night / Cyberpunk) são True Color. Em terminais 256-color, lipgloss converte automaticamente para o índice mais próximo no cubo XTerm. O resultado pode perder nuance — cores próximas podem colapsar para o mesmo índice.

**Cores em risco (Tokyo Night):**

| Papel | Hex exato | Índice 256 aproximado | Resultado visual |
|---|---|---|---|
| `surface.base` `#1a1b26` | 234 (`#1c1c1c`) | OK — escuro próximo |
| `surface.raised` `#24283b` | 236 (`#303030`) | OK — distinguível de base |
| `text.secondary` `#565f89` | 60 (`#5f5f87`) | OK — match aceitável |
| `border.default` `#414868` | 60 (`#5f5f87`) | Risco — colide com `text.secondary` |
| `special.highlight` `#283457` | 236 (`#303030`) | Risco — colide com `surface.raised` |

**Mitigação:** em decisões de design onde cor é a única diferenciação (ex: `border.default` vs `text.secondary`), garantir que exista também uma diferença estrutural (borda é box-drawing, texto é conteúdo) que sobreviva à colisão de cores.

### Atributos ANSI — matriz de suporte

| Atributo | Suporte | Risco | Fallback |
|---|---|---|---|
| **Bold** | Universal — todos os terminais | Nenhum | — |
| **Dim** (faint) | Amplo — falta em poucos | Baixo | Se ausente, terminal ignora (exibe normal) — aceitável, pois dim é reforço, não única diferenciação |
| **Italic** | Parcial — falha em: cmd.exe, ConHost (Windows legado), Terminal.app (macOS antigo), alguns terminais Linux sobre SSH | Médio | Texto italic aparece normal. Hints devem usar `text.secondary` (cor) como diferenciação primária, italic como reforço |
| **Underline** | Amplo | Baixo | Uso reservado — impacto mínimo se ausente |
| **Strikethrough** | Parcial — falha em: ConHost, Terminal.app, terminais mais antigos | Médio | Segredos excluídos devem ter `✕` + cor muted como diferenciação primária. Strikethrough é reforço visual |
| **Reverse** (inversão fg/bg) | Universal | Nenhum | Candidato alternativo para seleção de itens |
| **Foreground color** | Universal (ANSI 16), amplo (256), amplo (True Color) | Baixo | Lipgloss faz downgrade automático |
| **Background color** | Igual ao foreground | Baixo | Lipgloss faz downgrade automático |

### Restrições de caracteres Unicode

| Característica | Suporte | Risco | Fallback |
|---|---|---|---|
| **Box-drawing** (`─│╭╮╰╯`) | Universal em terminais gráficos | Falha rara — terminais sem Unicode | Lipgloss tem estilo `ASCII` (`-`, `\|`, `+`) |
| **Símbolos básicos** (`★✕▸▾·•`) | Amplo — presente na maioria das fontes monoespaçadas | Baixo | Alternativas ASCII: `*`, `x`, `>`, `v`, `.`, `o` |
| **Emoji** (`✅❌⚠️💡❓ℹ️`) | Parcial — largura inconsistente (1 ou 2 colunas), renderização varia entre terminais e fontes | Alto | Usar símbolos Unicode de 1 coluna (`✓✗⚠•?ℹ`) como fallback. Decisão emoji vs Unicode na implementação |
| **Nerd Fonts glyphs** | Requer instalação de fonte específica | Não usar | — |

### Largura de caractere — o problema dos 2 colunas

Alguns caracteres Unicode (especialmente emojis e CJK) ocupam **2 colunas** de terminal em vez de 1. Isso quebra layouts calculados se a contagem de colunas estiver errada.

| Caractere | Largura esperada | Largura real (varia) | Problema |
|---|---|---|---|
| `✅` | 1 | 2 (na maioria) | Desalinha colunas de tabela |
| `⚠️` | 1 | 1 ou 2 (inconsistente!) | Impossível calcular layout confiável |
| `★` | 1 | 1 (consistente) | Seguro |
| `◐` | 1 | 1 (consistente) | Seguro |

**Regra:** elementos usados em posições onde o alinhamento importa (command bar, colunas, prefixos de lista) devem usar **apenas caracteres de largura 1 garantida**. Emojis ficam restritos à barra de mensagens (onde desalinhamento de ±1 coluna é tolerável) ou são substituídos por símbolos.

### Princípios de compatibilidade

- **Degradação graceful:** a interface deve ser *funcional* em 256 cores e *usável* em 16 cores. True Color é preferência, não requisito.
- **Atributos como reforço, não como única diferenciação.** Se italic falhar, o hint ainda é visível via cor. Se strikethrough falhar, o `✕` e a cor muted ainda comunicam exclusão.
- **Testar nos 3 terminais de referência:** Windows Terminal (True Color), tmux sobre SSH (256), Terminal.app macOS (256 + italic limitado).
- **Respeitar `$NO_COLOR`.** Se a variável `NO_COLOR` estiver definida, desativar toda cor. Lipgloss/Bubble Tea v2 fazem isso automaticamente.
- **Largura segura.** Nunca usar emoji em posições de layout calculado. Reservar para texto livre onde ±1 coluna não afeta funcionalidade.
