# Design System — Abditum TUI

> Definições visuais fundamentais para o pacote `internal/tui`.  
> Complementa `tui-specification.md` (wireframes e comportamento) e `tui-elm-architecture.md` (arquitetura).
>
> **Wireframes com aplicação prática dos tokens:** ver [`tui-specification.md`](tui-specification.md)

---

## Paleta de Cores

A paleta é organizada por **papel funcional** — cada papel define *para que* a cor é usada, não qual cor concreta. Isso garante que trocar de tema é uma operação isolada: mudar os valores hex sem alterar lógica ou estrutura.

### Regras de aplicação

- **Texto sobre superfícies:** `text.primary` sobre `surface.base` deve ter contraste mínimo legível.
- **Foco de painel é implícito:** a command bar muda para refletir as ações do painel ativo — não existem bordas de foco em painéis. `border.focused` é reservado para diálogos modais.
- **Semânticas são reservadas:** cores semânticas aparecem somente para comunicar estado — nunca como decoração.
- **Consistência entre contextos:** a mesma cor semântica é usada em mensagens de aviso, modais de alerta, e demais elementos com o mesmo significado.

### Papéis e tokens

| Categoria | Papel | Uso | Tokyo Night | Cyberpunk |
|---|---|---|---|---|
| **Superfícies** | `surface.base` | Fundo principal | `#1a1b26` <span style="background:#1a1b26;color:#1a1b26">██</span> | `#0a0a1a` <span style="background:#0a0a1a;color:#0a0a1a">██</span> |
| | `surface.raised` | Painéis, modais, elevações | `#24283b` <span style="background:#24283b;color:#24283b">██</span> | `#1a1a2e` <span style="background:#1a1a2e;color:#1a1a2e">██</span> |
| | `surface.overlay` | Tooltips, menus, overlays | `#414868` <span style="background:#414868;color:#414868">██</span> | `#2a2a3e` <span style="background:#2a2a3e;color:#2a2a3e">██</span> |
| **Texto** | `text.primary` | Texto principal, labels | `#a9b1d6` <span style="color:#a9b1d6">██</span> | `#e0e0ff` <span style="color:#e0e0ff">██</span> |
| | `text.secondary` | Descrições, hints, placeholders | `#565f89` <span style="color:#565f89">██</span> | `#8888aa` <span style="color:#8888aa">██</span> |
| | `text.disabled` | Itens indisponíveis | `#3b4261` <span style="color:#3b4261">██</span> | `#444466` <span style="color:#444466">██</span> |
| **Bordas** | `border.default` | Bordas de diálogos neutros, separadores | `#414868` <span style="color:#414868">██</span> | `#3a3a5c` <span style="color:#3a3a5c">██</span> |
| | `border.focused` | Borda de diálogo modal ativo | `#7aa2f7` <span style="color:#7aa2f7">██</span> | `#ff2975` <span style="color:#ff2975">██</span> |
| **Interação** | `accent.primary` | Cursor, item selecionado, ação principal | `#7aa2f7` <span style="color:#7aa2f7">██</span> | `#ff2975` <span style="color:#ff2975">██</span> |
| | `accent.secondary` | Favoritos, decoração sutil | `#bb9af7` <span style="color:#bb9af7">██</span> | `#00fff5` <span style="color:#00fff5">██</span> |
| **Semânticas** | `semantic.success` | Confirmação, operação ok | `#9ece6a` <span style="color:#9ece6a">██</span> | `#05ffa1` <span style="color:#05ffa1">██</span> |
| | `semantic.warning` | Ação irreversível, bloqueio iminente | `#e0af68` <span style="color:#e0af68">██</span> | `#ffe900` <span style="color:#ffe900">██</span> |
| | `semantic.error` | Falha, exclusão | `#f7768e` <span style="color:#f7768e">██</span> | `#ff3860` <span style="color:#ff3860">██</span> |
| | `semantic.info` | Informação neutra, indicadores de sessão | `#7dcfff` <span style="color:#7dcfff">██</span> | `#00b4d8` <span style="color:#00b4d8">██</span> |
| **Especiais** | `special.muted` | Itens marcados para exclusão | `#565f89` <span style="color:#565f89">██</span> | `#666688` <span style="color:#666688">██</span> |
| | `special.highlight` | Fundo de item selecionado | `#283457` <span style="background:#283457;color:#a9b1d6">██</span> | `#2a1533` <span style="background:#2a1533;color:#e0e0ff">██</span> |

### Gradiente do logo

| Linha | Tokyo Night | Cyberpunk |
|---|---|---|
| 1 | `#9d7cd8` <span style="color:#9d7cd8">██</span> | `#ff2975` <span style="color:#ff2975">██</span> |
| 2 | `#89ddff` <span style="color:#89ddff">██</span> | `#b026ff` <span style="color:#b026ff">██</span> |
| 3 | `#7aa2f7` <span style="color:#7aa2f7">██</span> | `#00fff5` <span style="color:#00fff5">██</span> |
| 4 | `#7dcfff` <span style="color:#7dcfff">██</span> | `#05ffa1` <span style="color:#05ffa1">██</span> |
| 5 | `#bb9af7` <span style="color:#bb9af7">██</span> | `#ff2975` <span style="color:#ff2975">██</span> |

### Comparação

| Critério | Tokyo Night | Cyberpunk |
|---|---|---|
| **Conforto prolongado** | Excelente — dessaturada, tons frios e suaves | Moderado — neons cansam em sessões longas |
| **Legibilidade** | Alta — texto `#a9b1d6` sobre `#1a1b26` é equilibrado | Alta — texto `#e0e0ff` sobre `#0a0a1a` tem mais contraste |
| **Distinção semântica** | Clara — cores suficientemente distintas entre si | Muito clara — alta saturação torna diferenças óbvias |
| **Profissionalismo** | Alta — sóbria, familiar a devs (VS Code, IDEs) | Baixa — estética de entretenimento, pode parecer lúdica |
| **Adequação ao domínio** | Forte — ferramenta de segurança pede sobriedade | Fraca — neons contrastam com a seriedade de um cofre de senhas |
| **Expressividade do logo** | Elegante — gradiente suave violeta→ciano | Impactante — gradiente neon rosa→ciano→verde |
| **Acessibilidade** | Boa — contraste suficiente sem ser agressivo | Risco — neons podem ser problemáticos para sensibilidade visual |
| **Fidelidade em 256 cores** | **Boa** — tons dessaturados mapeiam bem para o cubo 256 | **Risco** — neons de alta saturação perdem intensidade no mapeamento 256, podendo parecer apagados |

### Decisão

> **Ambos os temas são suportados simultaneamente.** O usuário seleciona o tema ativo nas Configurações; F12 alterna rapidamente entre os dois sem abrir um menu.

A abstração por **papéis funcionais** é o que viabiliza isso: trocar de tema é uma operação isolada — mudar os valores hex em um único arquivo de estilos, sem alterar lógica ou estrutura.

---

## Tipografia

Em TUI não existem fontes nem tamanhos — o terminal usa fonte monoespaçada fixa. Os "pesos tipográficos" disponíveis são atributos ANSI que o lipgloss expõe. O suporte a esses atributos varia por terminal — a tabela abaixo documenta o comportamento esperado e o fallback de design para cada um.

### Atributos, suporte e fallback

| Atributo | Efeito visual | Suporte | Fallback se não suportado | Uso no Abditum |
|---|---|---|---|---|
| **Bold** | Texto mais brilhante e/ou espesso | Universal | — (degradação mínima) | Títulos, cursor selecionado, opção default |
| Dim / Faint | Brilho reduzido | Amplo; em alguns terminais exibe igual ao normal | Cor já comunica — perda tolerável | Itens desabilitados, conteúdo secundário |
| *Italic* | Texto inclinado | **Parcial** — vários terminais ignoram ou exibem como normal | `text.secondary` já comunica auxiliaridade | Pasta virtual (Favoritos), hints |
| Underline | Sublinhado | Amplo | — | Uso pontual — não está em uso ativo |
| ~~Strikethrough~~ | Texto riscado | **Parcial** — terminais legados ignoram | Símbolo `✕` + `special.muted` garantem legibilidade | Itens marcados para exclusão |
| Reverse | Inverte fg/bg | Amplo | — | Não está em uso ativo |
| Blink | Piscar | Disponível, mas frequentemente desativado pelo usuário ou terminal | Não usar para comunicar estado — decorativo | **Não usar** |

> **Nota sobre Bold:** em terminais que usam a mesma fonte para bold e normal (ex: alguns multiplexadores tmux), bold se manifesta apenas como cor mais brilhante, não como peso diferente. Isso é aceitável — bold continua sendo o atributo de destaque mais confiável.

> **Nota sobre Italic + Strikethrough:** ambos podem ser ignorados silenciosamente. Todo elemento que usa italic ou strikethrough deve ter um segundo diferenciador (cor, símbolo, ou estrutura) que preserve o significado na ausência do atributo.

### Combinações

Atributos podem ser combinados. Combinações previstas:

| Combinação | Uso |
|---|---|
| Bold + cor semântica | Título de modal — bold amarelo para alerta, bold ciano para info |
| Dim + strikethrough | Item marcado para exclusão e desabilitado — `✕` garante legibilidade sem esses atributos |
| Italic + `text.secondary` | Hints e pastas somente leitura — `text.secondary` é o diferenciador real |

### Princípios

- **Bold é o único destaque confiável.** Usar com moderação — se tudo for bold, nada é destaque.
- **Nunca depender de italic ou strikethrough como único diferenciador.** Esses atributos podem ser ignorados — cor e símbolo garantem legibilidade.
- **Dim indica existência sem relevância.** Prefira dim a invisível — o usuário precisa saber que o elemento existe.
- **Blink não é usado.** Frequentemente desativado; não é garantia de percepção.
- **Strikethrough tem significado semântico único** — "marcado para remoção". Não usar decorativamente.

---

## Bordas

A interface é **minimalista**: bordas são usadas apenas em **diálogos modais** e **separadores de linha**. Painéis, campos e listas não têm borda — espaço em branco e hierarquia tipográfica organizam o conteúdo.

### Estilos disponíveis

| Estilo | Caracteres | Exemplo |
|---|---|---|
| Rounded | `╭ ╮ ╰ ╯ │ ─` | `╭──────╮`<br>`│      │`<br>`╰──────╯` |
| Single | `┌ ┐ └ ┘ │ ─` | `┌──────┐`<br>`│      │`<br>`└──────┘` |
| Double | `╔ ╗ ╚ ╝ ║ ═` | `╔══════╗`<br>`║      ║`<br>`╚══════╝` |
| Thick | `┏ ┓ ┗ ┛ ┃ ━` | `┏━━━━━━┓`<br>`┃      ┃`<br>`┗━━━━━━┛` |
| Hidden | espaços | Sem borda visível — apenas padding |

### Aplicação

| Elemento | Estilo | Cor da borda | Princípio |
|---|---|---|---|
| Modal semântico | Rounded | Cor do tipo (`semantic.*` ou `accent.*`) | Borda colorida comunica o tipo do diálogo |
| Modal neutro | Rounded | `border.default` | Diálogos informativos sem urgência semântica |
| Separador vertical | `│` simples | `border.default` | Divide painéis side-by-side sem os envolver |
| Separador horizontal | `─` linha | `border.default` | Separa grupos de conteúdo ou seções |

### Princípios

- **Bordas apenas em diálogos.** Painéis, campos e listas não têm borda — espaço e tipografia estruturam o conteúdo sem envolver cada elemento num box.
- **Rounded é o único estilo usado.** Consistência visual em todos os diálogos — nunca misturar Single, Double ou Thick.
- **Cor da borda comunica semântica.** Em diálogos modais, a cor reforça o tipo (alerta, info, neutro) antes de o usuário ler o conteúdo.
- **Separadores são linhas, não boxes.** Áreas de conteúdo são divididas por `│` ou `─` — linhas que separam sem envolver.

---

## Ícones e Símbolos

Vocabulário de caracteres Unicode usados como ícones na interface. Usa-se Unicode básico (não Nerd Fonts) para máxima compatibilidade com terminais.

### Navegação em árvore ou hierarquia

| Símbolo | Uso |
|---|---|
| `▶` | Pasta com filhos, recolhida (U+25B6 BLACK RIGHT-POINTING TRIANGLE) |
| `▼` | Pasta com filhos, expandida (U+25BC BLACK DOWN-POINTING TRIANGLE) |
| `▷` | Pasta vazia — sem filhos (U+25B7 WHITE RIGHT-POINTING TRIANGLE) |
| `●` | Segredo normal — item folha (U+25CF BLACK CIRCLE) em `text.secondary` |
| `★` | Segredo favoritado — substitui `●` como prefixo (U+2605 BLACK STAR) em `accent.secondary` |

### Estados de itens

| Símbolo | Semântica |
|---|---|
| `★` | Segredo favoritado — **prefixo** do item na árvore (U+2605 BLACK STAR) em `accent.secondary`; substitui `●` |
| `✕` | Marcado para exclusão (U+2715 MULTIPLICATION X) |
| `+` | Adicionado na sessão atual (U+002B PLUS SIGN) — exibido em `semantic.info` |
| `~` | Modificado na sessão atual (U+007E TILDE) — exibido em `semantic.info` |
| `•` | Alterações não salvas no nível do cofre — indicador no header (U+2022 BULLET) |

### Mensagens (barra de mensagens)

| Símbolo | Tipo | Semântica | Colunas |
|---|---|---|---|
| `✓` | MsgSuccess | Sucesso (U+2713 CHECK MARK) | 1 |
| `ℹ` | MsgInfo | Informação (U+2139 INFORMATION SOURCE) | 1 |
| `⚠` | MsgWarn | Aviso (U+26A0 WARNING SIGN) | 1 |
| `✗` | MsgError | Erro (U+2717 BALLOT X) | 1 |
| `◐ ◓ ◑ ◒` | MsgBusy | Spinner — 4 frames | 1 |
| `•` | MsgHint | Dica contextual de campo (U+2022 BULLET) | 1 |
| `💡` | MsgTip | Dica de uso (U+1F4A1 LIGHT BULB) | **2** |

> **Largura dupla (`💡`):** o cálculo de truncamento da barra de mensagens deve reservar 2 colunas para o ícone MsgTip. Todos os outros ícones de mensagem ocupam 1 coluna.

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
| **Normal** | `text.primary` | `surface.base` | — | — | Item, campo, painel |
| **Painel ativo** | — | — | — | — | Foco indicado pela command bar — sem borda |
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

- **Foco de painel é pela command bar, seleção é por fundo.** Qual painel recebe input é indicado pelas ações exibidas na command bar. Qual item dentro do painel está selecionado é indicado por `special.highlight`.
- **Nunca depender só de cor.** Itens marcados para exclusão usam cor + strikethrough + símbolo `✕`. Itens favoritos usam cor + símbolo `★`. Redundância garante legibilidade em terminais com cores limitadas.
- **Dim é preferível a hidden.** Itens desabilitados devem ser visíveis (dim) para que o usuário saiba que existem — invisibilidade causa confusão.

---

## Componentes Modais

### Princípios gerais

- Modais são painéis sobrepostos **acima** de todo o frame, centralizados na tela.
- O conteúdo por trás do modal permanece visível mas **não recebe input** — apenas o modal do topo recebe eventos.
- O conteúdo de fundo **não é escurecido** (sem dim/fade). Re-renderizar o frame com cores alteradas adicionaria complexidade sem ganho de usabilidade em TUI.
- Modais se auto-dimensionam pelo conteúdo — não recebem tamanho alocado.
- A **command bar** troca para os atalhos do modal ativo enquanto ele estiver aberto.

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
| **256 cores (8-bit)** | 216 cores + 24 cinzas | **Alvo principal** — terminais UNIX/Linux, SSH, tmux | Design validado neste nível |
| **True Color (24-bit)** | 16 milhões | Terminais modernos (Alacritty, Kitty, WezTerm) | Aprimoramento automático via lipgloss downsampling |
| **16 cores (ANSI)** | 16 nomeadas | Consoles legados | Fallback mínimo — distinção semântica degradada |
| **Monocromático** | Preto e branco | Pipes/redirecionamento | Sem cor — não é caso de uso interativo |

> **Decisão técnica:** 256 cores é o mínimo garantido. Tokens são especificados como hex (True Color) — o lipgloss faz downsampling automático para o perfil do terminal detectado. O design deve ser **validado em 256 cores** antes de ser aprovado.

### Estratégia de alvo 256 cores

Os tokens são definidos como hex True Color, mas devem ser escolhidos considerando sua correspondência no cubo de 256 cores (6×6×6 + 24 cinzas). Cores próximas podem colapsar para o mesmo índice quando mapeadas — isso é tolerado somente quando existe diferenciação estrutural além da cor (símbolo, atributo tipográfico, ou posição na tela).

**Risco principal:** pares de papéis que precisam ser distinguíveis — ex: `surface.base` vs `surface.raised`, ou `border.default` vs `border.focused` — devem produzir índices 256 visivelmente distintos.

### Largura de caractere

Alguns caracteres Unicode (especialmente emojis) ocupam **2 colunas** de terminal em vez de 1, quebrando layouts calculados.

**Regra:** Caracteres em posições onde o alinhamento importa (tabelas, prefixos) devem ser de largura **1 coluna garantida**. Símbolos como `★`, `✓`, `✕`, `◐` são seguros. Emojis ficam restritos a texto livre onde desalinhamento é tolerável.

### Princípios de compatibilidade

- **256 cores é o alvo de design.** True Color é um aprimoramento automático — não é baseline. Decisões visuais (contraste, distinção semântica) devem ser aprovadas em 256 cores.
- **Atributos como reforço:** Nunca dependa de um único atributo para comunicar significado. Cor + estrutura + símbolo = redundância.
- **Testar em múltiplos ambientes:** Validar obrigatoriamente em 256 cores (ex: `TERM=xterm-256color`); validar True Color para garantir que o aprimoramento não introduz regressões.
- **Largura segura:** Emojis geram desalinhamento. Use símbolos Unicode ou texto em posições de layout crítico.
