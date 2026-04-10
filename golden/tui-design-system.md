# Design System — Abditum TUI

> Fundações visuais e padrões transversais para o pacote `internal/tui`.
> Define princípios, tokens, estados e padrões que governam toda decisão de UI.

## Fronteira deste documento

Este documento define **fundações** e **padrões reutilizáveis** — o que cada peça visual é, como se comporta em abstrato e como peças se combinam em situações recorrentes.

A composição dessas peças em telas, wireframes e fluxos concretos pertence ao documento de especificação:
- [`tui-specification-novo.md`](tui-specification-novo.md) — telas, wireframes de componentes e fluxos visuais

> **Regra de governança:** toda decisão de UI/UX deve ser compatível com os princípios definidos aqui. Em conflito entre especificação local e princípio, o princípio prevalece.

---

## Sumário

- [O Terminal como Meio](#o-terminal-como-meio)
- [Princípios](#princípios)
- [Paleta de Cores](#paleta-de-cores)
- [Tipografia](#tipografia)
- [Bordas](#bordas)
- [Dimensionamento e Layout](#dimensionamento-e-layout)
- [Ícones e Símbolos](#ícones-e-símbolos)
- [Estados Visuais](#estados-visuais)
- [Diálogos](#diálogos)
  - [Anatomia Comum](#anatomia-comum)
  - [Apresentação e Pilha](#apresentação-e-pilha)
  - [Dimensionamento](#dimensionamento)
  - [Scroll](#scroll)
  - [Identidade Visual](#identidade-visual)
    - [Severidade](#severidade)
  - [Teclado](#teclado)
  - [Notificação](#notificação)
  - [Confirmação](#confirmação)
  - [Ajuda](#ajuda)
  - [Funcional](#funcional)
- [Barra de Mensagens](#barra-de-mensagens)
  - [Anatomia](#anatomia-1)
  - [Dimensionamento](#dimensionamento-1)
  - [Identidade Visual](#identidade-visual-1)
  - [Teclado](#teclado-1)
  - [Eventos](#eventos)
  - [Ciclo de Vida](#ciclo-de-vida)
  - [Referência](#referência)
- [Barra de Comandos](#barra-de-comandos)
  - [Anatomia](#anatomia-2)
  - [Dimensionamento](#dimensionamento-2)
  - [Identidade Visual](#identidade-visual-2)
  - [Ações](#ações)
  - [Teclado](#teclado-2)
  - [Eventos](#eventos-1)
  - [Referência](#referência-1)
- [Foco e Navegação](#foco-e-navegação)
- [Teclado](#teclado-3)
  - [Notação](#notação)
  - [Convenções Semânticas](#convenções-semânticas)
  - [Escopos](#escopos)
  - [Atalhos Globais](#atalhos-globais)
  - [Regiões de Teclas de Função](#regiões-de-teclas-de-função)
- [Acessibilidade](#acessibilidade)


---

## O Terminal como Meio

**O que o terminal oferece:**
- Grade fixa de caracteres monospaced — alinhamento perfeito é gratuito
- Atributos ANSI: bold (universal), dim (amplo), italic (parcial), strikethrough (parcial), underline (amplo)
- Cores: true color em terminais modernos, 256 cores em legados, `NO_COLOR` como contrato de acessibilidade
- Teclado como canal de entrada primário — cada tecla é um evento discreto
- Mouse como canal secundário — clique e scroll, sem hover real nem drag contínuo
- Renderização de texto Unicode (BMP seguro; glifos de largura ambígua e Nerd Fonts são um risco)

**O que o terminal não tem:**
- Pixel independente, subpixel rendering, fontes customizadas, tamanhos de texto
- Z-index real, transparência, sombras, gradientes, bordas arredondadas reais
- Hover state, animação suave, transições visuais
- Layout flexível — a posição de cada caractere é absoluta na grade

**Consequências para o design:**
- A estrutura visual é construída por espaço em branco, alinhamento, separadores ASCII e hierarquia tipográfica — não por bordas decorativas nem containers visuais
- `bold` é o único destaque tipográfico universalmente confiável; `italic` e `strikethrough` precisam de reforço visual (detalhes na seção [Tipografia](#tipografia))
- Nenhum estado crítico pode depender exclusivamente de cor — cada estado usa pelo menos duas camadas de comunicação (detalhes na seção [Acessibilidade](#acessibilidade))
- Símbolos são escolhidos por clareza semântica e previsibilidade de renderização, não por estética (detalhes na seção [Ícones e Símbolos](#ícones-e-símbolos))
- O teclado é o caminho primário; toda ação acionável por teclado deve ser descobrível e executável também por mouse

---

## Princípios

Todos os princípios operam dentro do perímetro definido pelo terminal. Não há hierarquia entre eles — são compromissos simultâneos que a interface deve honrar. Quando dois princípios tensionam, a resolução é pelo contexto da tela específica, documentada na especificação.

### Identidade

- **Segurança como experiência:** segurança não é um recurso técnico invisível — é algo que o usuário deve *sentir* na interface. Operações sensíveis (revelar senha, exportar, sobrescrever, excluir) parecem deliberadas, com confirmação proporcional ao risco. Campos sensíveis são ocultos por padrão. A interface nunca expõe dados protegidos sem ação explícita.
- **Discrição e portabilidade:** a interface não chama atenção em ambientes públicos ou compartilhados. O visual é contido. Nenhum dado sensível aparece fora do contexto controlado pelo usuário. A aplicação não deixa rastros — não persiste estado fora do arquivo do cofre.
- **Controle total do usuário:** o usuário decide quando salvar, quando revelar, quando exportar. A aplicação não toma decisões irreversíveis em nome dele. Alterações permanecem reversíveis até o salvamento explícito. A única exceção é a alteração de senha mestra, que é imediata por necessidade criptográfica.
- **Simplicidade com profundidade:** a interface expõe primeiro o essencial — abrir, navegar, copiar. Complexidade (edição de estrutura, reordenação, busca, configurações) aparece apenas quando o usuário a procura. Um iniciante consegue usar o cofre em 30 segundos; um usuário avançado tem atalhos para tudo.

### Experiência

- **Hierarquia da informação:** o usuário distingue rapidamente contexto global (qual cofre, qual pasta, qual segredo), seleção atual, detalhe exibido e ações disponíveis. A importância relativa dos elementos é comunicada por posição, peso tipográfico e cor — nunca apenas por cor.
- **Estado sempre visível:** seleção, alterações pendentes, itens modificados, bloqueios, erros e processamento são perceptíveis sem exigir memorização do último comando executado. O estado do cofre (dirty/clean) está sempre no cabeçalho. O estado de cada segredo (adicionado/modificado/excluído) está junto ao item na árvore.
- **Feedback imediato:** toda ação relevante produz resposta visível — mudança de contexto, atualização do item, mensagem transitória ou indicador de progresso. Ausência de feedback é um defeito.
- **Reversibilidade por padrão:** ações destrutivas pedem confirmação. Ações de alto impacto oferecem cancelamento claro. Exclusão de segredos é uma marcação reversível até o salvamento.
- **Consistência de interação:** a mesma tecla, o mesmo símbolo e o mesmo tratamento visual mantêm o mesmo significado em toda a aplicação. `Enter` sempre avança ou aprofunda — confirma em diálogos, seleciona/expande na árvore, ativa/confirma em edição. `Esc` sempre retrocede ou abandona — fecha modais, cancela edição, sai de modos. O vetor direcional é consistente mesmo quando a ação concreta varia por escopo. Exceções devem ser documentadas e justificadas na especificação.
- **Estabilidade espacial:** cabeçalho, árvore, detalhe e barra de comandos permanecem em posições previsíveis entre estados. O layout não "pula" quando o conteúdo muda. Isso preserva memória muscular e reduz carga cognitiva.

---

## Paleta de Cores

A paleta é organizada por **papel funcional** — cada papel define *para que* a cor é usada, não qual cor concreta. Isso garante que trocar de tema é uma operação isolada: mudar os valores hex sem alterar lógica ou estrutura.

### Princípios da paleta

- **Papéis não são intercambiáveis:** mesmo quando dois tokens compartilham o mesmo valor, o papel funcional continua sendo diferente.
- **Semânticas não são decorativas:** `semantic.*` existe para comunicar estado, nunca para ornamentar a interface.
- **Contraste é obrigatório:** textos e sinais críticos precisam continuar legíveis sobre suas superfícies previstas.

### Papéis e tokens

| Categoria | Papel | Uso | Tokyo Night | Cyberpunk |
|---|---|---|---|---|
| **Superfícies** | `surface.base` | Cor de fundo da tela inteira | `#1a1b26` <span style="background:#1a1b26;color:#1a1b26">██</span> | `#0a0a1a` <span style="background:#0a0a1a;color:#0a0a1a">██</span> |
| | `surface.raised` | Fundo dos painéis laterais e das janelas que abrem sobre a tela | `#24283b` <span style="background:#24283b;color:#24283b">██</span> | `#1a1a2e` <span style="background:#1a1a2e;color:#1a1a2e">██</span> |
| | `surface.input` | Fundo dos campos de texto dentro de diálogos — tom rebaixado que delimita a área digitável | `#1e1f2e` <span style="background:#1e1f2e;color:#1e1f2e">██</span> | `#0e0e22` <span style="background:#0e0e22;color:#0e0e22">██</span> |
| **Texto** | `text.primary` | Texto normal — nomes de segredos, títulos de campos, conteúdo legível | `#a9b1d6` <span style="color:#a9b1d6">██</span> | `#e0e0ff` <span style="color:#e0e0ff">██</span> |
| | `text.secondary` | Texto de apoio — descrições de segredos, texto dentro de campos vazios, atalhos na barra inferior | `#565f89` <span style="color:#565f89">██</span> | `#8888aa` <span style="color:#8888aa">██</span> |
| | `text.disabled` | Texto de opções que não podem ser usadas no momento | `#3b4261` <span style="color:#3b4261">██</span> | `#444466` <span style="color:#444466">██</span> |
| | `text.link` | URLs e referências externas (tela Sobre) | `#7aa2f7` <span style="color:#7aa2f7">██</span> | `#ff2975` <span style="color:#ff2975">██</span> |
| **Bordas** | `border.default` | Linhas que dividem painéis, bordas de janelas informativas (ajuda, seleção de itens, navegação de arquivos) | `#414868` <span style="color:#414868\">██</span> | `#3a3a5c` <span style="color:#3a3a5c\">██</span> |
| | `border.focused` | Borda do painel ativo, de janelas de entrada (senhas, textos) e de diálogos com severidade neutra. Diálogos com severidade não-neutra usam `semantic.*` — ver [Diálogos](#diálogos) | `#7aa2f7` <span style="color:#7aa2f7\">██</span> | `#ff2975` <span style="color:#ff2975\">██</span> |
| **Interação** | `accent.primary` | Barra de seleção na lista, cursor de navegação, botão principal de ação | `#7aa2f7` <span style="color:#7aa2f7\">██</span> | `#ff2975` <span style="color:#ff2975\">██</span> |
| | `accent.secondary` | Ícone de favorito (★), nomes de pastas na navegação de arquivos | `#bb9af7` <span style="color:#bb9af7\">██</span> | `#00fff5` <span style="color:#00fff5\">██</span> |
| **Semânticas** | `semantic.success` | Operação concluída com sucesso, configuração ligada (ON) | `#9ece6a` <span style="color:#9ece6a\">██</span> | `#05ffa1` <span style="color:#05ffa1\">██</span> |
| | `semantic.warning` | Alerta antes de ação permanente, aviso de bloqueio por tentativas erradas, prefixos de estado dirty (`✦ ✎ ✗`) | `#e0af68` <span style="color:#e0af68\">██</span> | `#ffe900` <span style="color:#ffe900\">██</span> |
| | `semantic.error` | Erro de operação, senha incorreta, borda de diálogos destrutivos | `#f7768e` <span style="color:#f7768e\">██</span> | `#ff3860` <span style="color:#ff3860\">██</span> |
| | `semantic.info` | Informação contextual | `#7dcfff` <span style="color:#7dcfff\">██</span> | `#00b4d8` <span style="color:#00b4d8\">██</span> |
| | `semantic.off` | Configuração desligada (OFF) | `#737aa2` <span style="color:#737aa2\">██</span> | `#9999cc` <span style="color:#9999cc\">██</span> |
| **Especiais** | `special.muted` | Texto esmaecido — uso pontual em contextos que precisam de cor apagada sem conotação semântica | `#8690b5` <span style="color:#8690b5\">██</span> | `#666688` <span style="color:#666688\">██</span> |
| | `special.highlight` | Fundo colorido atrás do item selecionado na lista | `#283457` <span style="background:#283457;color:#a9b1d6">██</span> | `#2a1533` <span style="background:#2a1533;color:#e0e0ff">██</span> |
| | `special.match` | Trecho de texto que corresponde ao termo digitado na busca | `#f7c67a` <span style="color:#f7c67a">██</span> | `#ffc107` <span style="color:#ffc107">██</span> |

### Notas de contraste

> **`special.muted`:** usado para texto com aparência "apagada" sem conotação semântica específica. Contraste intencional abaixo do normal — adequado para conteúdo secundário pontual, não para informação crítica.

> **Aliases de valor:** `text.link` = `accent.primary` em hex. O alias documenta intenção — autores de temas podem divergir os valores quando precisarem distinguir link de ação primária.

> **Bordas de modais semânticos:** modais com severidade não-neutra usam diretamente os tokens `semantic.warning`, `semantic.info` ou `semantic.error` como cor de borda — não existe token `border.*` separado para casos semânticos. A severidade do diálogo governa a borda.

### Gradiente do logo

| Linha | Tokyo Night | Cyberpunk |
|---|---|---|
| 1 | `#9d7cd8` <span style="color:#9d7cd8">██</span> | `#ff2975` <span style="color:#ff2975">██</span> |
| 2 | `#89ddff` <span style="color:#89ddff">██</span> | `#b026ff` <span style="color:#b026ff">██</span> |
| 3 | `#7aa2f7` <span style="color:#7aa2f7">██</span> | `#00fff5` <span style="color:#00fff5">██</span> |
| 4 | `#7dcfff` <span style="color:#7dcfff">██</span> | `#05ffa1` <span style="color:#05ffa1">██</span> |
| 5 | `#bb9af7` <span style="color:#bb9af7">██</span> | `#ff2975` <span style="color:#ff2975">██</span> |

### Temas

> **Ambos os temas são suportados simultaneamente.** O usuário seleciona o tema ativo nas Configurações; `F12` alterna rapidamente entre os dois sem abrir um menu.

---

## Tipografia

Em TUI não existem fontes nem tamanhos; a tipografia disponível é o conjunto de atributos ANSI que o terminal realmente suporta. O papel do design system é definir **quando** usar esses atributos e como degradar quando eles falharem.

### Atributos e fallback

| Atributo | Suporte | Fallback | Uso principal |
|---|---|---|---|
| **Bold** | Universal | — | Títulos, cursor selecionado, ação default |
| Dim / Faint | Amplo | Cor já comunica o estado | Itens desabilitados, conteúdo secundário |
| *Italic* | Parcial | `text.secondary` já diferencia | Hints, pastas virtuais, textos auxiliares |
| Underline | Amplo | — | Uso pontual |
| ~~Strikethrough~~ | Parcial | `✗` + `special.muted` preservam o sentido | Itens marcados para exclusão |
| Blink | Inconsistente | Não usar | Nenhum |

### Combinações previstas

| Combinação | Uso |
|---|---|
| Bold + cor semântica | Título de modal de alerta ou informação |
| Dim + strikethrough | Item excluído, com `✗` como reforço |
| Italic + `text.secondary` | Hints e textos auxiliares |

> **Regra prática:** `bold` é o único destaque tipográfico realmente confiável. `italic` e ~~strikethrough~~ sempre precisam de reforço visual; `blink` não é usado.

---

## Bordas

A interface é minimalista: bordas aparecem apenas em modais e separadores. Painéis, listas e blocos de conteúdo são organizados por espaço, alinhamento e hierarquia tipográfica.

### Aplicação

| Elemento | Estilo | Token | Observação |
|---|---|---|---|
| Modal neutro | Rounded (`╭╮╰╯│─`) | `border.default` | Diálogo sem urgência semântica |
| Modal semântico | Rounded (`╭╮╰╯│─`) | `semantic.*` ou `accent.*` | Cor reforça o tipo do diálogo |
| Separador vertical | `│` | `border.default` | Divide painéis lado a lado |
| Separador horizontal | `─` | `border.default` | Separa grupos ou seções |
| Junction em T | `├` `┬` `┴` `┤` | `border.default` | Ponto onde separadores se encontram entre si ou se ligam a bordas internas de painéis (ex: FilePicker) |

> **Regra prática:** Rounded é o único estilo de caixa adotado. Separadores são linhas; a interface evita boxes decorativos. Junctions em T são usados exclusivamente para conectar separadores internos — nunca como ornamento.

---

## Dimensionamento e Layout

O design system define cores, tipografia, bordas e símbolos — mas sem dimensionamento, a implementação precisa inventar proporções. Esta seção estabelece as fundações de tamanho e espaço.

### Terminal suportado

| Parâmetro | Valor | Justificativa |
|---|---|---|
| Largura mínima | 80 colunas | Padrão POSIX; cabe em qualquer multiplexer |
| Altura mínima | 24 linhas | Padrão POSIX; garante cabeçalho + área útil + barra de comandos |
| Abaixo do mínimo | Degradação sem crash | Truncamento com `…`; a aplicação nunca quebra em terminais pequenos |

### Zonas verticais

A interface é dividida em quatro zonas empilhadas verticalmente:

| Zona | Altura | Conteúdo |
|---|---|---|
| **Cabeçalho** | 2 linhas | Nome da app, nome do cofre, indicador não salvo, abas de modo |
| **Área de trabalho** | Restante | Conteúdo do modo ativo (cofre, modelos, configurações, boas-vindas) |
| **Barra de mensagens** | 1 linha | Borda separadora `─` com mensagem embutida — quando há mensagem, o texto substitui o trecho central da borda |
| **Barra de comandos** | 1 linha | Ações do contexto ativo |

### Proporções de painel

Para modos com dois painéis (Cofre, Modelos):

| Painel | Proporção | Papel |
|---|---|---|
| Esquerdo (árvore / lista) | ~35% | Navegação e seleção |
| Direito (detalhe) | ~65% | Conteúdo do item selecionado |

A proporção é aproximada — a implementação pode ajustar em ±5% para alinhamento estético ou para acomodar terminais muito largos.

> **Dimensionamento de componentes:** as barras de mensagens e comandos possuem anatomia e dimensionamento detalhados nas seções dedicadas [Barra de Mensagens](#barra-de-mensagens) e [Barra de Comandos](#barra-de-comandos). Os diálogos têm dimensionamento próprio em [Diálogos — Dimensionamento](#dimensionamento).

## Ícones e Símbolos

Inventário completo dos caracteres Unicode usados pela interface.

### Critérios de seleção

> **Restrições:** BMP apenas (U+0000–U+FFFF) — emojis e Nerd Fonts excluídos. Todos os símbolos ocupam 1 coluna, exceto `<╡` (2 colunas, por composição). Símbolos de largura ambígua (locale-dependente) são evitados. Semântica prevalece sobre estética: `✗` (exclusão) e `✕` (erro) são distintos e ambos necessários.

### Inventário

O contexto de uso detalhado de cada símbolo está na seção onde ele é consumido (Sobreposição, Mensagens, Estados Visuais, especificação de telas).

| Símbolo | Nome semântico | Colunas | Bloco Unicode |
|---|---|---|---|
| `▶` | Pasta recolhida | 1 | Geometric Shapes |
| `▼` | Pasta expandida | 1 | Geometric Shapes |
| `▷` | Pasta vazia | 1 | Geometric Shapes |
| `●` | Item folha | 1 | Geometric Shapes |
| `★` | Favorito | 1 | Misc. Symbols |
| `✗` | Marcado para exclusão | 1 | Dingbats |
| `✦` | Recém-criado (não salvo) | 1 | Dingbats |
| `✎` | Modificado (não salvo) | 1 | Dingbats |
| `•` | Indicador contextual (ver nota) | 1 | Latin Supplement |
| `◉` | Campo revelável | 1 | Geometric Shapes |
| `✓` | Sucesso | 1 | Dingbats |
| `ℹ` | Informação | 1 | Letterlike Symbols |
| `⚠` | Alerta / aviso | 1 | Misc. Symbols |
| `✕` | Erro | 1 | Dingbats |
| `F1` | Atalho de ajuda | — | tecla de função |
| `◐ ◓ ◑ ◒` | Spinner de atividade | 1 | Geometric Shapes |
| `▌` | Cursor de campo | 1 | Block Elements |
| `↑` `↓` | Indicação de scroll (direção) | 1 | Arrows |
| `■` | Thumb de scroll (posição) | 1 | Geometric Shapes |
| `─` `│` | Separadores | 1 | Box Drawing |
| `├` `┬` `┴` `┤` | Junctions em T — pontos onde separadores se encontram ou se ligam a bordas de painéis internos | 1 | Box Drawing |
| `·` | Separador do cabeçalho | 1 | Latin Supplement |
| `╭╮╰╯` | Cantos arredondados | 1 | Box Drawing |
| `<╡` | Conector árvore → detalhe | 1+1 | Basic Latin + Box Drawing |
| `…` | Truncamento | 1 | Latin Supplement |
| `••••` | Máscara de conteúdo sensível | 1/cada | Latin Supplement |

> **`•` reutilizado:** aparece como indicador de alterações pendentes no cabeçalho, como marcador de dica contextual na barra de mensagens, como marcador de dica de uso, e como caractere de máscara em campos sensíveis. A distinção é sempre pelo contexto visual — nunca coexistem na mesma região.

> **`◐ ◓ ◑ ◒` e largura ambígua:** especificados como 1 coluna neste inventário. Em terminais com fontes que tratam Geometric Shapes como largura ambígua (dependente de locale), podem ser renderizados em 2 colunas — causando jitter na mensagem adjacente. Ver anti-padrão [Largura de Spinner Variante Entre Frames](#layout-e-estrutura).

---

## Estados Visuais

Estados visuais definem como o mesmo elemento muda de aparência conforme o contexto.

### Matriz resumida

| Estado | Tratamento visual |
|---|---|
| Normal | `text.primary` sobre `surface.base` |
| Selecionado | `special.highlight` + **bold** |
| Desabilitado | `text.disabled` + dim |
| Marcado para exclusão | `semantic.warning` + `✗` + ~~strikethrough~~ |
| Recém-criado (não salvo) | `✦` + texto `semantic.warning` |
| Modificado (não salvo) | `✎` + texto `semantic.warning` |
| Favorito | `★` em `accent.secondary` |
| Pasta virtual / leitura | `text.secondary` + *italic* |
| Campo sensível revelado | mesmo estilo do texto normal; a diferença é o valor exposto |
| Erro inline | `semantic.error` |

### Regras de transição

- Foco da área ativa é indicado pelo separador vertical e pela barra de comandos (ver [Foco e Navegação](#foco-e-navegação)) — não por borda ao redor do painel.
- TUI não tem estado "pressionado"; confirmação vem por mudança de contexto ou mensagem.
- Transições são instantâneas. A única animação prevista é o spinner `MsgBusy`.

## Diálogos

Diálogos são janelas sobrepostas que capturam o foco da aplicação para uma interação isolada — uma decisão, um reconhecimento, uma entrada de dados ou uma consulta de referência. Enquanto um diálogo estiver aberto, a área de trabalho permanece visível porém inativa.

A aplicação utiliza quatro tipos de diálogo, cada um com propósito, anatomia e regras de comportamento distintos:

| Tipo | Propósito | Seção |
|---|---|---|
| [Notificação](#notificação) | Informar um fato que exige reconhecimento | [▸](#notificação) |
| [Confirmação](#confirmação) | Solicitar uma escolha explícita do usuário | [▸](#confirmação) |
| [Ajuda](#ajuda) | Exibir referência de atalhos (somente leitura) | [▸](#ajuda) |
| [Funcional](#funcional) | Capturar entrada de dados com campos interativos | [▸](#funcional) |

Instâncias concretas de cada tipo — com títulos, mensagens e ações específicas — são documentadas na [Especificação Visual — Diálogos de Decisão](tui-specification.md#catálogo-de-diálogos-de-decisão) e na [Especificação Visual — Diálogos Funcionais](tui-specification.md#diálogos-funcionais).

### Anatomia Comum

Todo diálogo é composto por três regiões estruturais, desenhadas com bordas arredondadas (`╭╮╰╯│─`):

```text
╭── <símbolo>  <título> ─────────────────────╮  ← Borda Superior
│                                            │
│  <conteúdo do corpo>                       │  ← Corpo
│                                            │
╰── <Tecla Label> ──── <Tecla Label> ────────╯  ← Borda de Ações
```

**Borda Superior (título):**
- Estrutura char a char:
  - **Sem símbolo:** `╭──` (canto + 2× borda) + ` ` (1 espaço) + título + ` ` (1 espaço) + preenchimento `─`×N + `╮` (canto).
  - **Com símbolo:** `╭──` (canto + 2× borda) + ` ` (1 espaço) + símbolo + `  ` (2 espaços) + título + ` ` (1 espaço) + preenchimento `─`×N + `╮` (canto).
- O título ocupa a borda superior a partir da 5ª coluna (após `╭── `), preservando os caracteres de canto de ambos os lados. O preenchimento `─` garante pelo menos 1 caractere de borda antes do `╮`.
- Quando a severidade é Neutro ou o tipo de diálogo não usa severidade (Ajuda, Funcional), o símbolo é omitido.
- **Truncamento:** se o título excede o espaço disponível na borda (largura máxima do diálogo − cantos − espaçamento), ele é truncado com `…`.
- O título descreve o fluxo ou ação principal (ex: `Salvar cofre`, `Senha mestra`, `Ajuda`). Capitalizado conforme o nome, sem artigos desnecessários.

**Corpo:**
- Bordas laterais `│` delimitam o conteúdo.
- Padding interno: 2 colunas horizontais (entre `│` e o texto). Padding vertical de 1 linha (acima e abaixo do conteúdo) aplica-se a Notificação, Confirmação e Ajuda; diálogos Funcionais usam **0 linhas** de padding vertical — o conteúdo denso e interativo ocupa todo o espaço disponível.
- O conteúdo varia por tipo de diálogo: texto estático (Notificação, Confirmação), tabela de atalhos (Ajuda) ou campos interativos (Funcional).
- Diálogos Funcionais podem conter **divisores internos** que segmentam o corpo em regiões. A separação pode ser horizontal, vertical ou ambas:
  - **Horizontal:** `─` conectado às bordas laterais por T junctions (`├` à esquerda, `┤` à direita).
  - **Vertical:** `│` conectado às bordas horizontais por T junctions (`┬` no topo, `┴` na base).
  - **Cruzçamento:** `┼` onde um divisor horizontal e um vertical se cruzam.

**Borda de Ações (rodapé):**
- Estrutura char a char: `╰` (canto) + `─` (1× borda) + ações + preenchimento `─`×N (pelo menos 1) + `╯` (canto).
- Cada ação é representada como: ` ` (espaço) + tecla + ` ` (espaço) + label + ` ` (espaço). Exemplo: ` Enter Salvar `.
- Ações são separadas entre si por segmentos de preenchimento `─`.
- Layout varia conforme a quantidade de ações:
  - **1 ação:** alinhada à direita. Preenchimento `─` ocupa todo o espaço à esquerda.
    ```text
    ╰──────────────────────────── Enter OK ─╯
    ```
  - **2 ações:** principal à esquerda, cancelamento à direita. Preenchimento `─` entre elas.
    ```text
    ╰─ Enter Confirmar ──────── Esc Cancelar ─╯
    ```
  - **3 ações:** principal à esquerda, secundária ao centro, cancelamento à direita. Preenchimento `─` distribuído entre elas.
    ```text
    ╰─ S Sobrescrever ── N Como novo ── Esc Voltar ─╯
    ```
- Limite máximo: 3 ações. Diálogos com 4+ ações na borda são um [anti-padrão](tui-design-system-anti-patterns.md#diálogos-e-confirmações) ("Borda como Menu").

### Apresentação e Pilha

- O diálogo centraliza-se horizontal e verticalmente sobre a tela inteira.
- O conteúdo abaixo permanece visível, mas inativo (sem escurecimento de overlay).
- Apenas o elemento do topo da pilha recebe input; os inferiores permanecem montados, porém congelados.
- Ao fechar o elemento do topo, o foco retorna ao elemento imediatamente anterior na pilha.

### Dimensionamento

- **Largura mínima:** suficiente para acomodar o título e as ações da borda inferior, ou no mínimo 20 colunas.
- **Largura máxima:** até 95% da largura do terminal.
- **Largura fixa:** diálogos funcionais específicos podem definir largura fixa (ex: PasswordEntry = 50 colunas). A largura fixa é documentada na especificação de cada subtipo.
- **Altura:** determinada pelo contorno do conteúdo, sem espaços vazios exagerados.
- **Padding interno:** 2 colunas laterais. Padding vertical de 1 linha para Notificação, Confirmação e Ajuda; **0 linhas** para Funcional.

### Scroll

Quando o conteúdo do corpo excede o espaço disponível:

- **Largura excedida:** word-wrap quebra linhas mantendo integridade de palavras.
- **Altura excedida:** ativa-se scroll vertical com indicadores visuais na borda lateral direita.
- A borda superior e a borda de ações nunca participam do scroll — permanecem sempre fixas.
- **Pré-condição:** o scroll só é renderizável se o diálogo possuir pelo menos 5 linhas internas (excluindo borda superior e borda de ações). Abaixo desse mínimo, o conteúdo é truncado sem indicadores de scroll.
- **Navegação por teclado:** teclas direcionais (`↑`/`↓`), `PgUp`/`PgDn`, `Home`/`End`.

**Composição da borda lateral direita com scroll:**

A borda lateral direita do corpo é composta por 3 elementos, cada um ocupando posições fixas:

| Elemento | Posição | Caractere | Descrição |
|---|---|---|---|
| Seta superior | 1ª linha do corpo | `↑` | Indica conteúdo acima do viewport. Substitui o `│` da borda |
| Thumb | Entre a 2ª e a penúltima linha do corpo | `■` | Posição relativa do viewport no conteúdo total. Nunca sobrepõe as setas |
| Seta inferior | Última linha do corpo | `↓` | Indica conteúdo abaixo do viewport. Substitui o `│` da borda |

- As setas `↑` e `↓` ocupam sempre a primeira e a última linha do corpo, respectivamente.
- O thumb `■` é posicionado proporcionalmente entre a 2ª linha e a penúltima linha do corpo — ele **nunca** é desenhado sobre a posição de uma seta.
- Nas linhas onde nenhum indicador está presente, a borda permanece `│`.

Wireframe ilustrando o scroll ativo (5 linhas internas, com padding vertical):

```text
╭── ⚠  Título do Diálogo ────────────────────╮
│                                            ↑
│  Primeira linha do conteúdo longo.         ■
│  Segunda linha mostrando limite excedido.  │
│  Terceira linha com mais informações.      │
│                                            ↓
╰── Enter Salvar ── A Alt ──── Esc Cancelar ─╯
```

### Identidade Visual

Regras visuais padrão aplicadas a **todos** os diálogos. Cada tipo documenta apenas as variações.

> Caracteres estruturais: ver [Anatomia Comum](#anatomia-comum).

| Elemento | Token | Atributo | Observação |
|---|---|---|---|
| Bordas e cantos | Determinado pela severidade ou pelo tipo | — | Notificação/Confirmação: token da [Severidade](#severidade). Ajuda: `border.default`. Funcional: `border.focused` |
| Símbolo na borda superior | Determinado pela severidade | — | `⚠`, `✕` ou `ℹ` conforme severidade. Omitido em Neutro, Ajuda e Funcional |
| Título | `text.primary` | **bold** | Descreve o fluxo ou ação principal |
| Texto do corpo | `text.primary` | — | — |
| Tecla da ação default (`Enter`) | Token da tecla default da severidade | **bold** | Notificação/Confirmação: ver [Severidade](#severidade). Ajuda/Funcional: `accent.primary` |
| Teclas de ações secundárias e cancelamento | Segue o token de borda | — | — |

#### Severidade

Severidade governa o tratamento visual — borda, símbolo e cor da tecla default — aplicado exclusivamente aos diálogos de **Notificação** e **Confirmação**. Diálogos de Ajuda e Funcional não utilizam severidade.

| Severidade | Símbolo | Token de borda | Token da tecla default | Quando usar |
|---|---|---|---|---|
| Destrutivo | `⚠` | `semantic.warning` | `semantic.error` | Ação irreversível ou com perda de dados |
| Erro | `✕` | `semantic.error` | `accent.primary` | Falha ocorrida, condição irrecuperável |
| Alerta | `⚠` | `semantic.warning` | `accent.primary` | Situação importante mas recuperável |
| Informativo | `ℹ` | `semantic.info` | `accent.primary` | Informação que requer atenção |
| Neutro | — | `border.focused` | `accent.primary` | Operação rotineira, sem urgência |

> **Nota:** severidades Destrutivo e Alerta compartilham o símbolo `⚠` e o token de borda `semantic.warning`. A distinção visual está na tecla default: `semantic.error` (vermelho) para destrutivo, `accent.primary` para alerta. Isso reforça que o perigo está na *ação*, não apenas na *situação*.

Cada tipo de diálogo é documentado com o seguinte template de sub-seções: **Quando usar**, **Anatomia**, **Variações Visuais**, **Ações**, **Teclado**, **Barra de Comandos**, **Barra de Mensagens**, **Exemplo Visual** e **Referência**.

### Teclado

Convenções de teclado aplicadas a todos os diálogos. Cada tipo documenta apenas as variações.

**Teclas implícitas — `Enter` e `Esc`:**

Todo diálogo possui duas teclas implícitas que não precisam ser declaradas pelas ações:

| Tecla | Papel implícito |
|---|---|
| `Enter` | Executa a ação **principal** (a da extrema esquerda, ou a única ação) |
| `Esc` | Executa a ação de **cancelamento** (a da extrema direita, ou a única ação) |

Comportamento conforme a quantidade de ações:

| Ações | `Enter` | `Esc` |
|---|---|---|
| 1 ação | Executa a ação única | Mesmo efeito que `Enter` — executa a ação única |
| 2 ações | Executa a ação da esquerda (principal) | Executa a ação da direita (cancelamento) |
| 3 ações | Executa a ação da esquerda (principal) | Executa a ação da direita (cancelamento) |

**Teclas explícitas — letras de atalho:**

Além das teclas implícitas, cada ação pode declarar uma tecla de atalho (tipicamente a primeira letra da label). A tecla declarada é exibida na borda de ações antes da label (ex: `S Sobrescrever`). As teclas implícitas `Enter` e `Esc` continuam funcionando mesmo quando a ação possui tecla explícita. Exemplo:

```text
╰─ S Sobrescrever ── N Como novo ── Esc Voltar ─╯
```

- `S` → Sobrescrever (tecla explícita)
- `Enter` → Sobrescrever (tecla implícita — ação principal)
- `N` → Como novo (tecla explícita)
- `Esc` → Voltar (tecla implícita — cancelamento)

A ação secundária (centro, quando presente) **sempre** precisa declarar sua tecla explícita — não possui tecla implícita.

**Diálogos Funcionais — exceção do `Enter`:**

Em diálogos funcionais, `Enter` pode ter comportamento contextual dependendo do campo em foco (ex: submeter um campo, selecionar um item em lista). Entretanto, em algum estado do diálogo o `Enter` **deve** acionar a confirmação do diálogo. `Esc` sempre cancela e fecha o diálogo, sem exceção.

### Notificação

**Quando usar:** o usuário precisa tomar ciência de um fato — uma falha, um alerta ou uma informação relevante. Não há decisão a tomar; apenas reconhecimento.

**Anatomia:**

| Região | Presença | Conteúdo |
|---|---|---|
| Borda Superior | Obrigatória | Símbolo de severidade + título |
| Corpo | Obrigatório | Apenas afirmação. Sem pergunta. Frases terminam com ponto final |
| Borda de Ações | Obrigatória | Exatamente 1 ação, alinhada à direita |

**Redação do corpo:** afirmação concisa e direta. Referências a itens específicos em aspas simples. Exemplos:
- `Arquivo corrompido ou inválido. Necessário fechar.`
- `Senhas não conferem. Necessário digitar novamente.`
- `Arquivo inválido ou versão não suportada. Necessário corrigir.`

**Variações Visuais:** sem variações — segue integralmente a [Identidade Visual](#identidade-visual) geral com severidade.

**Ações:**
- Borda de Ações: exclusivamente `Enter OK`, alinhada à direita.
- Nenhuma ação secundária. Nenhuma ação de cancelamento.

**Teclado:** sem variações — segue integralmente as convenções de [Teclado](#teclado) geral (1 ação: `Enter` e `Esc` ambos fecham).

**Barra de Comandos:** vazia. Ações do diálogo não se repetem na barra.

**Barra de Mensagens:** limpa durante toda a exibição do diálogo.

**Exemplo Visual:**

```text
╭── ✕  Arquivo corrompido ───────────────╮
│                                        │
│  Arquivo corrompido ou inválido.       │
│  Necessário fechar.                    │
│                                        │
╰────────────────────────────── Enter OK ╯
```

**Referência:** instâncias concretas no [Catálogo de Diálogos de Decisão](tui-specification.md#catálogo-de-diálogos-de-decisão).

### Confirmação

**Quando usar:** o usuário precisa fazer uma escolha explícita que confirma, bifurca ou cancela um fluxo — salvar, descartar, sobrescrever, excluir.

**Anatomia:**

| Região | Presença | Conteúdo |
|---|---|---|
| Borda Superior | Obrigatória | Símbolo de severidade (quando não Neutro) + título |
| Corpo | Obrigatório | Afirmação de contexto (opcional, terminada em ponto) + pergunta objetiva (terminada em `?`) |
| Borda de Ações | Obrigatória | 2 ou 3 ações |

**Redação do corpo:** fato opcional seguido de pergunta concisa que apresenta as opções de decisão. A pergunta não menciona a opção `Voltar` (Esc). Referências a itens específicos em aspas simples. Exemplos:
- `Sair do Abditum?`
- `Cofre modificado. Salvar ou descartar?`
- `Arquivo modificado externamente. Sobrescrever?`
- `'Gmail' será excluído permanentemente. Continuar?`

**Variações Visuais:** sem variações — segue integralmente a [Identidade Visual](#identidade-visual) geral com severidade.

**Ações:**
- Borda de Ações: 2 a 3 ações.
  - Ação principal à esquerda (ex: `Enter Salvar`, `S Sobrescrever`).
  - Ação secundária ao centro, quando presente (ex: `N Salvar como novo`).
  - `Esc Cancelar` (ou `Esc Voltar`) sempre na extrema direita.
- Todas as ações ficam ativas simultaneamente — não há validação condicional.

**Teclado:** sem variações — segue integralmente as convenções de [Teclado](#teclado) geral (2–3 ações com teclas explícitas).

**Barra de Comandos:** vazia. Decisões ficam exclusivamente na borda de ações.

**Barra de Mensagens:** limpa durante toda a exibição do diálogo.

**Exemplo Visual:**

```text
╭── ⚠  Salvar cofre ─────────────────────────────╮
│                                                │
│  Arquivo modificado externamente.              │
│  Sobrescrever ou salvar como novo?             │
│                                                │
╰── S Sobrescrever ── N Como novo ─ Esc Voltar ──╯
```

**Referência:** instâncias concretas no [Catálogo de Diálogos de Decisão](tui-specification.md#catálogo-de-diálogos-de-decisão).

### Ajuda

**Quando usar:** o usuário precisa consultar os atalhos de teclado disponíveis no contexto atual. Acionado por `F1`. Diálogo somente leitura, sem impacto no estado da aplicação.

**Anatomia:**

| Região | Presença | Conteúdo |
|---|---|---|
| Borda Superior | Obrigatória | Título `Ajuda` em **bold**, sem símbolo |
| Corpo | Obrigatório | Tabela de atalhos organizada por contexto (seções com cabeçalho). Scroll ativado quando o conteúdo excede a altura |
| Borda de Ações | Obrigatória | Exatamente 1 ação: `Esc Fechar`, alinhada à direita |

**Variações Visuais:**
- Não usa severidade. Borda em `border.default`.
- Nomes das teclas de atalho no corpo em `text.primary`; descrições em `text.secondary`.

**Ações:**
- Borda de Ações: exclusivamente `Esc Fechar`, alinhada à direita.

**Teclado:** sem variações — segue integralmente as convenções de [Teclado](#teclado) geral (1 ação: `Esc` fecha). Teclas de scroll (`↑`/`↓`, `PgUp`/`PgDn`, `Home`/`End`) ativas quando há conteúdo excedente.

**Barra de Comandos:** pode exibir ações auxiliares do contexto (ex: `F12` para troca de tema), sem repetir a ação da borda.

**Barra de Mensagens:** limpa durante toda a exibição do diálogo.

**Exemplo Visual:**

```text
╭── Ajuda ──────────────────────────────╮
│ Árvore                                ↑
│ F2       Renomear arquivo atual       ■
│ Ctrl+N   Novo arquivo no diretório    │
│ Ctrl+D   Marcar para exclusão         ↓
╰──────────────────────────── Esc Fechar ╯
```

**Referência:** especificação completa em [Especificação Visual — Help](tui-specification.md#help).

### Funcional

**Quando usar:** o usuário precisa fornecer dados por meio de campos interativos — entrada de senha, seleção de arquivo, criação de senha com confirmação. Diferente dos diálogos de decisão, o Funcional captura input estruturado.

**Anatomia:**

| Região | Presença | Conteúdo |
|---|---|---|
| Borda Superior | Obrigatória | Título em **bold**, sem símbolo de severidade |
| Corpo | Obrigatório | Campos de entrada (`input`), labels, contadores e outros componentes interativos. Conteúdo varia por subtipo. Sem padding vertical (0 linhas acima e abaixo) |
| Divisores internos | Opcional | Separadores horizontais (`─` com `├` / `┤`), verticais (`│` com `┬` / `┴`) ou ambos (`┼` no cruzamento), segmentando o corpo em regiões distintas |
| Borda de Ações | Obrigatória | Ação de confirmação + ação de cancelamento (2 ações) |

**Variações Visuais:**
- Não usa severidade. Borda em `border.focused`.
- Labels de campo: `accent.primary` + **bold** quando o campo está ativo; `text.secondary` quando inativo.
- Área de campo de entrada: fundo `surface.input`.
- Máscara de senha: caracteres `●` em `text.secondary`, com comprimento fixo (não revela o tamanho real da senha).
- Cursor no campo ativo: `▌` em `text.primary`.
- Ação default (`Enter`): estado condicional — `text.disabled` / dim enquanto houver validações pendentes; `accent.primary` + **bold** ao satisfazer condições.

**Ações:**
- Borda de Ações: 2 ações.
  - Ação de confirmação à esquerda (ex: `Enter Confirmar`).
  - `Esc Cancelar` à direita.
- Estado condicional do `Enter` descrito em [Variações Visuais](#identidade-visual) acima.

**Teclado:**
- `Enter` e `Esc` conforme [Teclado](#teclado) geral, com a exceção de que `Enter` pode ter comportamento contextual por campo (ver [Teclado — Diálogos Funcionais](#teclado)).
- `Tab` / `Shift+Tab`: navega entre campos (quando há múltiplos campos).
- Teclas de edição: digitação, `Backspace`, `Del` — comportamento padrão de campo de texto.

**Barra de Comandos:** exibe ações auxiliares específicas do diálogo (ex: `Tab Campo seguinte`, `Del Limpar linha`). Ações da borda de ações **nunca se repetem** na barra de comandos.

**Barra de Mensagens:** território de uso exclusivo do diálogo funcional.
- Ao abrir o diálogo: dica de campo (`•`) orientando a ação esperada.
- Durante a interação: dica atualizada conforme o campo em foco.
- Após erro de validação: mensagem de erro (`✕`) exibida até correção ou troca de campo.
- Ao fechar o diálogo: barra limpa (responsabilidade do orquestrador).

**Exemplo Visual:**

```text
╭── Alterar senha mestra ───────────────╮
│  Senha atual:   ••••••••              │
│  Nova senha:    ▌                     │
╰── Enter Confirmar ────── Esc Cancelar ╯
```

Exemplo com divisores internos (FilePicker simplificado):

```text
╭── Abrir cofre ──────────┬─────────────╮
│  📁 Documentos          │ cofre.abdt  │
│  📁 Projetos            │ notas.abdt  │
│    📄 cofre.abdt        │             │
├─────────────────────────┴─────────────┤
│  Arquivo: cofre.abdt                  │
╰── Enter Abrir ──────── Esc Cancelar ──╯
```

**Subtipos conhecidos:**

| Subtipo | Propósito | Campos | Referência |
|---|---|---|---|
| PasswordEntry | Entrada de senha para abrir cofre | 1 campo (senha) | [spec](tui-specification.md#passwordentry) |
| PasswordCreate | Criação ou alteração de senha mestra | 2–3 campos (atual, nova, confirmação) | [spec](tui-specification.md#passwordcreate) |
| FilePicker | Seleção de arquivo (abrir ou salvar) | Árvore de diretórios + campo de nome | [spec](tui-specification.md#filepicker) |

Cada subtipo tem anatomia interna, estados e validações específicas documentadas na especificação visual.

## Barra de Mensagens

A barra de mensagens é a linha entre a área de trabalho e a barra de comandos (conforme [Dimensionamento e Layout](#dimensionamento-e-layout)). Exibe uma mensagem por vez — nova mensagem substitui a anterior imediatamente. Não há fila nem pilha.

### Anatomia

A barra é uma linha de borda `─` na largura total do terminal. Quando há mensagem ativa, símbolo e texto são embutidos na borda, seguindo a mesma regra de composição da [Borda Superior de diálogos](#anatomia-comum).

**Estrutura char a char:**

- **Com símbolo:** `───` (3× borda) + ` ` (1 espaço) + símbolo + `  ` (2 espaços) + texto + ` ` (1 espaço) + preenchimento `─`×N (pelo menos 1).
- **Sem símbolo:** `───` (3× borda) + ` ` (1 espaço) + texto + ` ` (1 espaço) + preenchimento `─`×N (pelo menos 1).
- **Sem mensagem:** `─` repetido na largura total do terminal.

Símbolos possíveis: `✓` sucesso · `ℹ` informação · `⚠` alerta · `✕` erro · `◐◓◑◒` spinner · `•` dica. Cada símbolo ocupa 1 coluna. Todos os tipos atuais possuem símbolo — o caso sem símbolo é previsto para extensibilidade.

**Truncamento:** se o texto excede o espaço disponível, é truncado com `…`. Largura máxima do texto: largura do terminal − 9 (com símbolo) ou largura do terminal − 6 (sem símbolo).

**Wireframe — com símbolo:**

```
───␣✓␣␣Cofre salvo␣───────────────────────────────────────────────────────────
```

**Wireframe — sem símbolo:**

```
───␣Cofre salvo␣──────────────────────────────────────────────────────────────
```

**Wireframe — sem mensagem (idle):**

```
──────────────────────────────────────────────────────────────────────────────
```

> Legenda: `␣` = espaço.

### Dimensionamento

| Parâmetro | Valor |
|---|---|
| Altura | 1 linha fixa |
| Largura | 100% da largura do terminal |
| Prefixo | 3 colunas de borda `───` + 1 espaço |
| Símbolo | 1 coluna |
| Espaçamento símbolo → texto | 2 espaços |
| Espaço após texto | 1 espaço |
| Sufixo mínimo | 1 coluna de borda `─` |
| Largura máxima do texto | largura do terminal − 9 (com símbolo) · largura do terminal − 6 (sem símbolo) |

### Identidade Visual

#### Severidade

A barra de mensagens utiliza o mesmo sistema de severidade dos [diálogos](#severidade). Cada tipo de mensagem herda o token semântico correspondente — a mesma paleta que governa bordas e símbolos de diálogos governa a cor da mensagem inteira.

| Severidade | Tipo de mensagem | Símbolo | Token | Atributo |
|---|---|---|---|---|
| Erro | Erro | `✕` | `semantic.error` | **bold** |
| Alerta | Alerta | `⚠` | `semantic.warning` | — |
| Informativo | Informação | `ℹ` | `semantic.info` | — |
| Neutro | Sucesso | `✓` | `semantic.success` | — |

> Não existe tipo de mensagem "Destrutivo" na barra — ações destrutivas são sempre comunicadas por diálogos.

#### Tipos não-semânticos

Além das severidades, a barra suporta tipos utilitários sem correspondência com severidade de diálogos:

| Tipo | Símbolo | Token | Atributo |
|---|---|---|---|
| Ocupado (spinner) | `◐ ◓ ◑ ◒` | `accent.primary` | — |
| Dica de campo | `•` | `text.secondary` | *italic* |
| Dica de uso | `•` | `text.secondary` | *italic* |

#### Regras de cor

- Token se aplica à mensagem inteira — símbolo e texto usam o mesmo token de cor. Não há distinção de cor entre o símbolo e o conteúdo textual dentro de uma mesma mensagem.
- Borda `─` sempre em `border.default`, independente do tipo de mensagem.

### Teclado

A barra de mensagens é passiva — não aceita interação por teclado.

### Eventos

| Evento | Reação da barra |
|---|---|
| Nenhuma mensagem ativa | Borda `─` contínua |
| Orquestrador emite mensagem | Símbolo + texto embutidos na borda |
| Nova mensagem emitida | Substitui imediatamente a mensagem anterior |
| TTL expira | Mensagem desaparece → borda `─` |
| Diálogo funcional abre | Dica de campo (`•`) orientando a ação esperada |
| Campo recebe foco (diálogo funcional) | Dica atualizada conforme o campo |
| Erro de validação (diálogo funcional) | Mensagem de erro (`✕`) até correção ou troca de campo |
| Diálogo fecha | Barra limpa → borda `─` |

### Ciclo de Vida

| Tipo | TTL padrão | Dismissal padrão |
|---|---|---|
| Sucesso | 5 s | Expiração |
| Informação | 5 s | Expiração |
| Alerta | 5 s | Expiração |
| Erro | 5 s | Expiração |
| Ocupado (spinner) | Sem TTL | Substituição explícita por Sucesso, Erro ou Alerta ao término da operação |
| Dica de campo | Permanente | Troca de campo ou substituição por outro tipo |
| Dica de uso | Permanente | Substituição por qualquer outro tipo |

- O caller pode sobrescrever TTL e trigger de dismissal conforme o contexto.
- Spinner avança 1 frame/segundo sincronizado com tick global.

> O ciclo de vida da barra em diálogos funcionais (mensagem de contexto ao abrir, dica por campo, limpeza ao fechar) é contrato do orquestrador — documentado em [Diálogos — Funcional](#funcional).

### Referência

Instâncias concretas e wireframes expandidos na [especificação de telas](tui-specification.md):

- [Barra de Mensagens](tui-specification.md#barra-de-mensagens) — tokens, estados, eventos, comportamento

## Barra de Comandos

A barra de comandos é a última linha da tela (conforme [Dimensionamento e Layout](#dimensionamento-e-layout)). Exibe as ações acionáveis por teclado no contexto atual — o usuário nunca precisa adivinhar o que pode fazer.

### Anatomia

A barra é uma linha de texto na largura total do terminal. Ações são distribuídas à esquerda; a âncora `F1 Ajuda` é fixa à direita. O espaço restante é preenchido com espaços.

**Estrutura char a char:**

` ` (2 espaços) + ação₁ + ` · ` (separador) + ação₂ + ` · ` + … + ` `×N (preenchimento, pelo menos 1) + `F1 Ajuda`

Cada ação é composta por: TECLA + ` ` (1 espaço) + Label. Exemplo: `^S Salvar`, `Del Excluir`.

- **Prefixo:** 2 espaços fixos. A primeira ação começa na 3ª coluna.
- **Separador:** ` · ` (espaço + middle dot + espaço = 3 colunas) entre ações adjacentes.
- **Preenchimento:** espaços entre a última ação e a âncora. Mínimo 1 espaço.
- **Âncora:** `F1 Ajuda` (8 colunas) — sempre na extrema direita, nunca removida.

**Truncamento por prioridade:** quando não há espaço para todas as ações, as de menor prioridade são removidas primeiro (ver [Ações](#ações)). A âncora `F1 Ajuda` nunca é sacrificada.

**Wireframe — estado normal:**

```
␣␣^I␣Novo␣·␣^E␣Editar␣·␣Del␣Excluir␣·␣^S␣Salvar␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣F1␣Ajuda
```

**Wireframe — espaço restrito (ações truncadas por prioridade):**

```
␣␣^I␣Novo␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣F1␣Ajuda
```

**Wireframe — diálogo de decisão ativo (vazia):**

```
␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣␣F1␣Ajuda
```

> Legenda: `␣` = espaço; `·` = separador (caractere real da barra).

### Dimensionamento

| Parâmetro | Valor |
|---|---|
| Altura | 1 linha fixa |
| Largura | 100% da largura do terminal |
| Prefixo | 2 espaços |
| Formato de ação | TECLA + 1 espaço + Label |
| Separador entre ações | ` · ` (3 colunas) |
| Preenchimento | espaços (mínimo 1 coluna) |
| Âncora | `F1 Ajuda` (8 colunas, extrema direita) |
| Espaço disponível para ações | largura do terminal − 2 (prefixo) − 8 (âncora) − 1 (preenchimento mínimo) |

### Identidade Visual

| Elemento | Token | Atributo |
|---|---|---|
| Tecla da ação (ex: `^S`) | `accent.primary` | **bold** |
| Label da ação (ex: `Salvar`) | `text.primary` | — |
| Separador ` · ` | `text.secondary` | — |
| Âncora `F1 Ajuda` — tecla | `accent.primary` | **bold** |
| Âncora `F1 Ajuda` — label | `text.primary` | — |

### Ações

Cada ação registrada no contexto ativo possui atributos que controlam sua apresentação:

| Atributo | Efeito na barra | Efeito no Help |
|---|---|---|
| `Enabled = true` | Exibida com estilo normal | Listada |
| `Enabled = false` | Oculta | Listada |
| `HideFromBar = true` | Oculta | Listada |

Regras de layout:

- **Prioridade** governa ordenação: maior prioridade → mais à esquerda.
- **Espaço insuficiente:** ações de menor prioridade removidas primeiro.
- **`F1 Ajuda` sempre visível** — âncora fixa na extrema direita; o cálculo de espaço desconta `F1 Ajuda` antes de distribuir as demais ações.
- **Ações desabilitadas desaparecem** — `Enabled = false` remove da barra (não fica dim). A ação continua listada no Help.
- **Ações de confirmação/cancelamento** (`Enter`/`Esc`) já estão na borda inferior do diálogo — não são duplicadas na barra.

### Teclado

A barra funciona como guia de descoberta — o usuário vê quais teclas estão disponíveis no contexto atual. As teclas são atribuídas conforme as [Convenções Semânticas](#convenções-semânticas) e os [Escopos](#escopos) definidos em [Teclado](#teclado-3).

### Eventos

| Evento | Reação da barra |
|---|---|
| Área de trabalho com foco | Exibe ações do painel ativo (árvore ou detalhe) |
| Troca de foco entre painéis | Atualiza para ações do painel que recebe foco |
| Diálogo de decisão aberto | Vazia (apenas `F1 Ajuda`) |
| Diálogo funcional aberto | Exibe ações internas do diálogo |
| Diálogo fecha (pop da pilha) | Volta para ações do contexto anterior |
| Terminal redimensionado | Recalcula ações visíveis (prioridade governa corte) |

### Referência

Instâncias concretas e wireframes expandidos na [especificação de telas](tui-specification.md):

- [Barra de Comandos](tui-specification.md#barra-de-comandos) — anatomia detalhada, atributos, eventos, comportamento

## Foco e Navegação

A interface sempre possui exatamente um elemento com foco. Todas as ações — exibidas na barra de comandos e acionadas por teclado — são relativas ao elemento focado. O foco determina **o quê** o usuário está manipulando; as teclas determinam **como** (ver [Teclado — Convenções Semânticas](#convenções-semânticas)).

O modelo de foco é o mesmo na interface principal e nos diálogos funcionais — as regras abaixo se aplicam uniformemente.

### Conceito de foco

A interface é organizada em **áreas** que contêm **elementos**:

| Contexto | Áreas | Elementos |
|---|---|---|
| Interface principal | Painéis (árvore, detalhe) | Nós, itens de lista, campos editáveis |
| Diálogo funcional | Regiões do diálogo | Campos de entrada, listas internas |

Em ambos os contextos:

- Exatamente um elemento possui foco a qualquer momento.
- `Tab` / `⇧Tab` alternam o foco entre áreas. Setas movem entre elementos dentro da área.
- O ciclo entre áreas é circular. Áreas vazias ou sem conteúdo interativo são puladas.

O foco governa dois reflexos imediatos na interface:

- **Ações disponíveis:** a [Barra de Comandos](#barra-de-comandos) exibe as ações do elemento ou área focada. Trocar o foco atualiza as ações visíveis.
- **Dica contextual:** a [Barra de Mensagens](#barra-de-mensagens) pode exibir uma dica (`•`) associada ao elemento focado — orientando a ação esperada ou descrevendo o campo.

### Indicação visual

A distinção visual do foco é **consistente** — o elemento focado sempre recebe destaque, independente do contexto:

| Contexto | Elemento focado | Tratamento visual |
|---|---|---|
| Árvore / lista | Nó ou item selecionado | `special.highlight` + **bold** |
| Árvore → separador | Linha do item selecionado | `│` substituído por `<╡` em `accent.primary` |
| Campo de entrada (ativo) | Campo com cursor | Fundo `surface.input` + cursor `▌` em `text.primary` + label em `accent.primary` **bold** |
| Campo de entrada (inativo) | Campo sem foco | Fundo `surface.input` + label em `text.secondary` |
| Área ativa | Área com foco | Identificada pela **barra de comandos** (exibe ações da área focada) |

- A interface principal não usa bordas ao redor de painéis — existe apenas um separador vertical `│` em `border.default` entre árvore e detalhe.
- Campos não possuem borda — a área digitável é delimitada pelo fundo `surface.input` (tom rebaixado em relação ao `surface.raised` do diálogo).
- Placeholder em `text.secondary` + *italic* — desaparece ao digitar.
- Em **NO_COLOR**: o fundo `surface.input` pode ser perdido; cursor + label em **bold** permanecem como indicadores de foco suficientes.

### Navegação

A navegação é nativamente por teclado, mas o mouse deve ser suportado para ativar foco em qualquer elemento válido.

**Teclado:**
- `Tab` / `⇧Tab` movem o foco entre áreas (painéis ou regiões do diálogo).
- Setas `↑↓←→` movem o foco entre elementos dentro da área ativa.
- `Home` / `End` movem ao primeiro / último elemento visível.
- As convenções semânticas completas estão em [Teclado](#teclado-3).

**Mouse:**
- Clique em um elemento válido transfere o foco imediatamente para ele.
- Toda ação acionável por teclado deve ser executável também por mouse.

### Validação de campos

- Erro de validação: exibido na barra de mensagens (tipo Erro), não inline — os formulários são simples o suficiente para mostrar um erro por vez.

## Teclado

A aplicação é operada inteiramente por teclado. Esta seção define a notação para representar teclas na documentação e na interface, as convenções semânticas que governam o significado de cada tecla, e a política de escopos e reservas. O mapeamento completo por tela está na [especificação de telas](tui-specification-novo.md).

### Notação

Convenção de representação textual de teclas e modificadores — usada neste documento, na barra de comandos e no diálogo de Ajuda:

| Modificador | Notação | Unicode |
|---|---|---|
| `Ctrl` | `⌃` | U+2303 |
| `Shift` | `⇧` | U+21E7 |
| `Alt` | `!` | — |

Teclas especiais são escritas por extenso: `Enter`, `Esc`, `Tab`, `Del`, `Ins`, `Home`, `End`, `PgUp`, `PgDn`. Combinações são concatenadas sem espaço: `⌃Q`, `⌃!⇧Q`, `⇧F6`.

### Convenções Semânticas

Significado fixo das teclas estruturais — válido em toda a aplicação. Componentes específicos (diálogos, barras) documentam apenas variações.

| Tecla | Significado |
|---|---|
| `Enter` | Confirma, avança ou aprofunda — confirma em diálogos, seleciona/expande na árvore, inicia/termina edição de campo |
| `Esc` | Interrompe, retrocede ou abandona — fecha diálogo, cancela edição, sai de modo |
| `Tab` | Avança foco para o próximo bloco — próximo painel (modo leitura) ou próximo campo (modo edição) |
| `⇧Tab` | Retorna foco para o bloco anterior |
| `↑` `↓` `←` `→` | Navegação direcional em listas, árvores e campos |
| `Home` / `End` | Primeiro / último item visível, ou início / fim de linha em campos |
| `PgUp` / `PgDn` | Scroll por página (viewport − 1) |
| `Ins` | Inserção / criação no contexto do foco |
| `Del` | Exclusão no contexto do foco |

- Se uma tecla precisa ter significado diferente em dois contextos, a especificação deve documentar e justificar a exceção.
- Teclas de navegação universais (`↑↓←→`, `Tab`, `Home`, `End`, `PgUp`, `PgDn`) não aparecem na barra de comandos — são senso comum em TUI. Exceção: diálogos podem exibir opções explicitamente.

### Escopos

Escopos mais específicos sobrepõem os mais gerais quando ambos estão ativos:

| Escopo | Descrição | Exemplo |
|---|---|---|
| **Global** | Funciona em qualquer contexto | `F1`, `F12`, `⌃Q`, `⌃!⇧Q` |
| **Área de trabalho** | Quando a área de trabalho tem foco (sem diálogo) | `F2`–`F11`, `⇧F6`, `⇧F7`, `⌃F7` |
| **Diálogo** | Enquanto um diálogo está no topo da pilha | `Enter`, `Esc`, `Tab` |
| **Contextual** | Ações específicas do item ou campo com foco | `Ins`, `Del`, `⌃<letra>` |

### Atalhos Globais

As 4 teclas que funcionam em qualquer contexto da aplicação:

| Tecla | Ação | Notas |
|---|---|---|
| `F1` | Abrir / fechar Ajuda | |
| `F12` | Alternar tema | Não exibida na barra de comandos |
| `⌃Q` | Sair da aplicação | Com confirmação quando há alterações |
| `⌃!⇧Q` | Bloquear cofre | Emergencial — descarta alterações sem confirmação. Atalho "complicado" para evitar acidentes |

### Regiões de Teclas de Função

As teclas de função são reservadas por grupos, seguindo a ergonomia do teclado físico:

| Região | Uso |
|---|---|
| `F2` a `F4` | Seleção de áreas de trabalho (Cofre, Modelos, Configurações) |
| `F5` a `F8` | Ações de persistência (criar, abrir, salvar, recarregar) |
| `F9` a `F11` | Ações complementares (exportar, importar, alterar senha mestra) |

A atribuição específica de cada tecla a fluxos individuais está na [especificação de telas](tui-specification.md).

## Acessibilidade

### NO_COLOR e modo monocromático

Quando `$NO_COLOR` está definido (ou o terminal informa que não suporta cores), `lipgloss` remove todas as cores. A interface deve permanecer totalmente funcional.

**Princípio:** nenhum estado crítico pode depender exclusivamente de cor. Todo estado usa pelo menos duas camadas de comunicação (cor + símbolo, cor + atributo tipográfico, símbolo + posição).

**Matriz de fallback:**

| Estado visual | Com cor | Fallback NO_COLOR |
|---|---|---|
| Item selecionado | `special.highlight` + **bold** | **bold** |
| Aba ativa | `special.highlight` + **bold** | **bold** + borda `╭───╮` |
| Badge `⚠ Fraca` | `semantic.warning` | `⚠ Fraca` (símbolo preserva semântica) |
| Badge `✓ Forte` | `semantic.success` | `✓ Forte` |
| Config "ativado" | `semantic.success` | `ativado` (texto preserva estado) |
| Config "desativado" | `semantic.off` | `desativado` (texto preserva estado) |
| Dirty `•` | `semantic.warning` | `•` (símbolo preserva estado) |
| Busca match | `special.match` + **bold** | **bold** |
| Exclusão `✗` | `semantic.warning` + strikethrough | `✗` + strikethrough |
| Recém-criado `✦` | `semantic.warning` | `✦` (símbolo preserva estado) |
| Modificado `✎` | `semantic.warning` | `✎` (símbolo preserva estado) |
| Favorito `★` | `accent.secondary` | `★` (símbolo preserva semântica) |
| Máscara `••••••••` | `text.secondary` | `text.primary` |
| Borda de modal | `semantic.*` / `border.*` | Borda presente — tipo distinguido por símbolo no título |
| Campo de input | `surface.input` | `surface.base` — cursor + label em **bold** preservam foco |

