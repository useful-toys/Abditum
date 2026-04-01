# Design System — Abditum TUI

> Definições visuais fundamentais para o pacote `internal/tui`.  
> Complementa `tui-design.md` (layout e interação) e `tui-elm-architecture.md` (arquitetura).
>
> **Wireframes com aplicação prática dos tokens:** ver [`tui-specification.md`](tui-specification.md)

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
- **Consistência entre contextos:** a mesma cor semântica é usada em mensagens de aviso, modais de alerta, e demais elementos com o mesmo significado — nunca um contexto usa uma cor semântica diferente para o mesmo tipo de informação.

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

| Atributo | Efeito visual | Uso comum |
|---|---|---|
| **Bold** | Texto mais brilhante e/ou espesso | Títulos, opções default, labels em foco |
| *Italic* | Texto inclinado (suporte varia por terminal) | Texto auxiliar, placeholders, descrições |
| Dim | Texto com brilho reduzido | Itens desabilitados, conteúdo secundário |
| Underline | Sublinhado | Reservado — uso pontual (links, atalhos) |
| ~~Strikethrough~~ | Texto riscado | Itens marcados para remoção |
| Normal | Sem atributo | Corpo de texto, valores, itens de lista |

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

| Elemento | Estilo sugerido | Cor da borda | Princípio |
|---|---|---|---|
| Painel inativo | Rounded | `border.default` | Cantos arredondados — visual moderno e leve |
| Painel com foco | Rounded | `border.focused` | Mesma forma, cor diferente — foco por cor |
| Modal (semântico) | Rounded | Cor do tipo | Borda colorida comunica tipo/semântica |
| Modal (neutro) | Rounded | `border.default` ou `border.focused` | Dependente do contexto |
| Separador | `│` simples | `border.default` | Linha simples para discreção |

### Princípios

- **Um único estilo de canto** (provavelmente Rounded) para consistência. Evitar misturar Single e Double na mesma interface.
- **Diferenciação por cor, não por estilo.** Painel ativo = mesma borda, cor diferente. Mais sutil que trocar de Single para Double.
- **Bordas são discretas.** O conteúdo é protagonista — bordas enquadram sem competir.
- **Título na borda.** Painéis e modais podem ter título integrado à borda superior (ex: `╭─ Cofre ───────╮`). Lipgloss suporta esse padrão.

---

## Ícones e Símbolos

Vocabulário de caracteres Unicode usados como ícones na interface. Usa-se Unicode básico (não Nerd Fonts) para máxima compatibilidade com terminais.

### Navegação em árvore ou hierarquia

| Símbolo | Uso |
|---|---|
| `▸` | Item recolhido/não expandido (U+25B8 BLACK RIGHT-POINTING SMALL TRIANGLE) |
| `▾` | Item expandido (U+25BE BLACK DOWN-POINTING SMALL TRIANGLE) |
| `·` | Item folha (U+00B7 MIDDLE DOT) |

### Estados de itens

| Símbolo | Semântica |
|---|---|
| `★` | Favorito (U+2605 BLACK STAR) — exibido em `accent.secondary` |
| `☆` | Não favorito — quando necessário mostrar ambos os estados (U+2606 WHITE STAR) |
| `✕` | Marcado para exclusão (U+2715 MULTIPLICATION X) |
| `+` | Adicionado na sessão atual (U+002B PLUS SIGN) — exibido em `semantic.info` |
| `~` | Modificado na sessão atual (U+007E TILDE) — exibido em `semantic.info` |
| `•` | Alterações não salvas no nível do cofre — indicador no header (U+2022 BULLET) |

### Mensagens (barra de mensagens)

| Símbolo | Semântica | Nota |
|---|---|---|
| `✓` | Sucesso (U+2713 CHECK MARK) | Alternativa a ✅ para máxima compatibilidade |
| `⚠` | Aviso (U+26A0 WARNING SIGN) | Atenção requerida |
| `✗` | Erro (U+2717 BALLOT X) | Alternativa a ❌ |
| `◐ ◓ ◑ ◒` | Progresso/ativo | Spinner — 4 frames |
| `•` | Informação (U+2022 BULLET) | Alternativa a 💡 |

> **Emoji vs Unicode:** os emojis (`✅ ⚠️ ❌ 💡`) são visualmente mais ricos mas ocupam 2 colunas em muitos terminais e podem não renderizar em todos os ambientes. Os símbolos Unicode acima são fallback de 1 coluna. A decisão emoji vs Unicode será tomada na implementação com testes em terminais reais.

### Tipos de diálogo (semântico)

| Símbolo | Tipo | Semântica |
|---|---|---|
| `?` | Question | Decisão neutra (U+003F) |
| `⚠` | Alert/Warning | Ação potencialmente destrutiva (U+26A0) |
| `ℹ` | Info | Informação (U+2139 INFORMATION SOURCE) |

### Campos com conteúdo sensível

| Símbolo | Uso |
|---|---|
| `•` | Caractere de substituição — repetido para preencher o espaço do valor oculto (U+2022 BULLET), ex: `••••••••` |
| `◉` | Indicador de que o campo pode ser revelado — exibido como sufixo do **label** do campo, não do valor (U+25C9 FISHEYE) |

> **Nota:** `•` como máscara não é um ícone de estado — é um caractere de substituição repetido, semanticamente distinto do `•` de alterações não salvas no header.

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
| **Normal** | `text.primary` | `surface.base` | — | `border.default` | Item, campo, painel inativo |
| **Com foco** | `text.primary` | `surface.base` | — | `border.focused` | Painel ativo |
| **Selecionado (cursor)** | `text.primary` | `special.highlight` | **Bold** | — | Item sob cursor em árvore ou lista |
| **Desabilitado** | `text.disabled` | `surface.base` | Dim | — | Ação indisponível |
| **Marcado para exclusão** | `special.muted` | `surface.base` | ~~Strikethrough~~ | — | Item com `✕` |
| **Favorito** | `text.primary` | `surface.base` | — | — | Item com `★` em `accent.secondary` |
| **Adicionado (sessão)** | `text.primary` | `surface.base` | — | — | Item com `+` em `semantic.info` |
| **Modificado (sessão)** | `text.primary` | `surface.base` | — | — | Item com `~` em `semantic.info` |
| **Campo sensível revelado** | `text.primary` | `surface.base` | — | — | Valor temporariamente visível — sem diferenciação de cor |
| **Pasta virtual / somente leitura** | `text.secondary` | `surface.base` | *Italic* | — | Pasta Favoritos — não editável |
| **Erro inline** | `semantic.error` | `surface.raised` | — | — | Validação com falha |

> **Nota:** TUIs não têm estado "pressionado" (pressed). Confirmação de ação é comunicada por mudança de contexto ou mensagem na barra de status.

### Transições

Em TUI, estados mudam **instantaneamente** — sem animação nem fade. A única animação é o spinner `MsgBusy` (1fps). Transições suaves não são viáveis em terminais.

### Princípios

- **Foco é por borda, seleção é por fundo.** Dois conceitos distintos: foco indica *qual painel* recebe input; seleção indica *qual item* dentro do painel é o alvo.
- **Nunca depender só de cor.** Itens marcados para exclusão usam cor + strikethrough + símbolo `✕`. Itens favoritos usam cor + símbolo `★`. Redundância garante legibilidade em terminais com cores limitadas.
- **Dim é preferível a hidden.** Itens desabilitados devem ser visíveis (dim) para que o usuário saiba que existem — invisibilidade causa confusão.

---

## Componentes Modais

### Princípios gerais

- Modais são painéis sobrepostos **acima** de todo o frame, centralizados na tela.
- O conteúdo por trás do modal permanece visível mas **não recebe input** — apenas o modal do topo recebe eventos.
- Modais se auto-dimensionam pelo conteúdo — não recebem tamanho alocado.
- A **interface de status** (ex: command bar) muda durante modal: exibe atalhos do modal ativo em vez das ações de contexto.

### Navegação padrão

| Tecla | Comportamento |
|---|---|
| **ENTER** | Aciona a opção marcada como `Default` |
| **ESC** | Aciona a opção marcada como `Cancel`. Se não houver opção cancel, fecha o modal (dismiss) |
| Atalho da opção | Aciona diretamente a opção correspondente |

### Tipos de modal

#### Confirmação e perguntas

Apresentam um **título**, uma **mensagem** explicativa (opcional), e um conjunto finito de opções. Cada opção tem label, atalhos, e um Cmd a executar.

##### Tipo semântico (`DialogType`)

O tipo semântico determina o emoji e a cor base do modal, comunicando a natureza da decisão antes de o usuário ler o conteúdo.

| Tipo | Emoji | Semântica | Cor base |
|---|---|---|---|
| Question | ❓ | Decisão neutra | Azul (`#7aa2f7`) |
| Alert | ⚠️ | Ação potencialmente destrutiva | Amarelo (`#e0af68`) |
| Info | ℹ️ | Informação que requer reconhecimento | Ciano (`#7dcfff`) |

A cor base é aplicada à **borda ou título** do modal (definição visual exata adiada para fase de implementação). O emoji é exibido junto ao título.

##### Variantes pré-definidas

| Factory | Opções | Default (ENTER) | Cancel (ESC) | Padrão de uso |
|---|---|---|---|---|
| `Confirm` | Sim / Não | Sim | Não | Confirmações binárias |
| `ConfirmOrCancel` | Sim / Não / Cancelar | Sim | Cancelar | Decisões com escape |
| `Ask` | Opções customizadas | Configurável | Configurável | Múltiplas alternativas |

##### Composição visual

- **Título** — frase curta, exibida com o emoji do tipo semântico (ex: "⚠️ Confirmação").
- **Mensagem** — texto explicativo. Pode ser vazio se o título for autoexplicativo.
- **Opções** — exibidas como botões ou labels com atalhos.
- A opção padrão é visualmente destacada (bold ou cor).
- Atalhos aparecem na interface de status.

Modais de confirmação comunicam a decisão a quem os abriu via callback ou similar mecanismo.

#### Mensagem informativa

Apresenta um título e texto descritivo, com tipo semântico (`DialogType`). Sem opções — apenas dismiss.

| Factory | Dismiss | Tipo |
|---|---|---|
| `Message` | ESC ou ENTER | Qualquer `DialogType` |

**Composição visual:**

- Título com emoji do tipo semântico (ex: "ℹ️ Informação", "⚠️ Aviso").
- Texto livre abaixo — pode ter múltiplas linhas.
- Cor da borda/título segue o tipo.
- Interface de status mostra apenas: OK / Fechar.

#### Entrada de dados sensíveis

Campos cujo conteúdo deve ser protegido visualmente.

| Variante | Campos | Validação |
|---|---|---|
| Single | 1 campo | Nenhuma |
| Dual | 2 campos (repetição) | Verifica correspondência |

**Composição visual:**

- Título descritivo.
- Campo(s) com máscara — caracteres substituídos por símbolo (ex: `•`).
- **Cor:** mesma que campos normais (`text.primary`). A máscara já comunica o caráter sensível.
- Interface de status: Confirmar / Cancelar.

**Nota:** Dados sensíveis devem nunca transitar por log ou broadcast.

#### Seleção de arquivo

Navegação e seleção de arquivo do sistema de arquivos.

| Modo | Função |
|---|---|
| Open | Selecionar arquivo existente |
| Save | Escolher destino para escrita |

**Composição visual:**

- Título descritivo.
- Navegação de diretórios.
- Modo save: campo para nome do arquivo.
- Interface de status: Selecionar / Cancelar.

#### Entrada de texto

Campo de texto livre com validação opcional.

| Tipo | Validação |
|---|---|
| Simples | Nenhuma |
| Com validação | Callback/função |

**Composição visual:**

- Título descritivo.
- Campo com placeholder opcional.
- Se validação falhar, mensagem de erro exibida — modal não fecha.
- Interface de status: Confirmar / Cancelar.

#### Seleção em lista

Lista de opções navegável e selecionável.

**Composição visual:**

- Título descritivo.
- Lista de itens com cursor.
- Item selecionado visualmente destacado.
- Interface de status: Selecionar / Cancelar / Navegar.

#### Help/Ajuda

Modal com referência de ações e atalhos disponíveis.

**Composição visual:**

- Ações agrupadas por categoria.
- Cada ação mostra: teclas, descrição.
- Pode ter scroll se a lista for longa.
- Interface de status: Fechar.
- Dismiss via tecla de escape.

### Stack de modais

Modais podem se sobrepor. Regras:

- O modal do topo recebe input. Modais abaixo permanecem vivos mas não recebem eventos.
- Modais abaixo não são atualizados até serem trazidos ao topo.
- A interface de status sempre reflete os atalhos do modal do **topo** da stack.

---

## Compatibilidade de Terminal

TUIs rodam em ambientes heterogêneos. O design system deve funcionar desde terminais modernos (24-bit color, todos os atributos) até terminais com capacidades reduzidas.

### Níveis de cor

| Nível | Cores | Contexto | Abordagem |
|---|---|---|---|
| **True Color (24-bit)** | 16 milhões | Terminais modernos | Alvo principal — cores exatas |
| **256 cores** | 216 cores + 24 cinzas | Terminais UNIX/Linux | Fallback — cores mapeadas para cubo próximo |
| **16 cores (ANSI)** | 16 nomeadas | Consoles legados | Fallback mínimo — sem distinção fina |
| **Monocromático** | Preto e branco | Pipes/redirecionamento | Sem cor |

### Estratégia de fallback

Quando degradando de True Color (24-bit) para 256 cores, cores próximas podem colapsar para o mesmo índice. Isso é tolerado quando existe diferenciação estrutural além da cor (ex: borda é box-drawing, texto é conteúdo semantic).

### Suporte a atributos ANSI

Atributos como bold, italic, dim, underline e strikethrough têm suporte variado:

- **Bold:** Universal
- **Dim:** Amplo, com fallback seguro (exibe normal)
- **Italic:** Parcial em alguns terminais
- **Underline:** Amplo
- **Strikethrough:** Parcial em alguns terminais
- **Cores:** Universal (ANSI 16), amplo (256), moderno (True Color)

**Princípio:** Use atributos como reforço de diferenciação visual, nunca como único indicador. Se italic falhar, a cor ou a estrutura devem comunicar a mesma intenção.

### Largura de caractere

Alguns caracteres Unicode (especialmente emojis) ocupam **2 colunas** de terminal em vez de 1, quebrando layouts calculados.

**Regra:** Caracteres em posições onde o alinhamento importa (tabelas, prefixos) devem ser de largura **1 coluna garantida**. Símbolos como `★`, `✓`, `✕`, `◐` são seguros. Emojis ficam restritos a texto livre onde desalinhamento é tolerável.

### Princípios de compatibilidade

- **Degradação graceful:** A interface deve funcionar em 256 cores e ser usável em 16 cores. True Color é preferência.
- **Atributos como reforço:** Nunca dependa de um único atributo para comunicar significado. Cor + estrutura + símbolo = redundância.
- **Testar em múltiplos ambientes:** Validar True Color, 256 cores e 16 cores.
- **Largura segura:** Emojis geram desalinhamento. Use símbolos Unicode ou texto em posições de layout crítico.
