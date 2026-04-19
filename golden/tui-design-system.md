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
- [Barra de Mensagens](#barra-de-mensagens)
- [Barra de Comandos](#barra-de-comandos)
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
| | `surface.raised` | Fundo dos painéis laterais e das janelas que abrem sobre a tela, incluindo todos os diálogos modais | `#24283b` <span style="background:#24283b;color:#24283b">██</span> | `#1a1a2e` <span style="background:#1a1a2e;color:#1a1a2e">██</span> |
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
| `▌` | Cursor de bloco (não usado — cursor real do terminal usado) | 1 | Block Elements |
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

Diálogos são janelas sobrepostas modais que capturam o foco para uma interação isolada. A área de trabalho permanece visível porém inativa — sem escurecimento de overlay. Os diálogos operam em pilha: apenas o topo recebe input; ao fechar, o foco retorna ao elemento anterior.

**Fundo:** todos os diálogos usam `surface.raised` como cor de fundo de todas as linhas do corpo (bordas laterais e conteúdo). Isso os distingue visualmente da `surface.base` da tela por trás.

A aplicação utiliza quatro tipos — Notificação, Confirmação, Ajuda e Funcional. A anatomia comum (borda superior com título, corpo, borda de ações com até 3 ações), o dimensionamento, o sistema de scroll, a identidade visual, a severidade e as convenções de teclado de diálogos são especificados em [Especificação Visual — Diálogos](tui-spec-dialogos.md#diálogos). O contrato de cada tipo e as instâncias concretas também residem na especificação:

- [Anatomia, dimensionamento e identidade visual](tui-spec-dialogos.md#diálogos) — estrutura comum a todos os tipos
- [Notificação, Confirmação, Ajuda, Funcional](tui-spec-dialogos.md#notificação) — contrato de cada tipo
- [Catálogo de Diálogos de Decisão](tui-spec-dialogos.md#catálogo-de-diálogos-de-decisão) — instâncias concretas
- [Diálogos de Senha](tui-spec-dialog-senha.md) — PasswordEntry, PasswordCreate
- [FilePicker](tui-spec-dialog-filepicker.md) — seleção de arquivo (Open, Save)
- [Ajuda](tui-spec-dialog-help.md) — referência de atalhos

## Barra de Mensagens

A barra de mensagens é a zona fixa entre a área de trabalho e a barra de comandos (conforme [Dimensionamento e Layout](#dimensionamento-e-layout)). Exibe uma mensagem por vez — nova mensagem substitui a anterior imediatamente. Não há fila nem pilha. A borda `─` é permanente; quando há mensagem, símbolo e texto são embutidos na borda.

Anatomia, dimensionamento, identidade visual, eventos e ciclo de vida especificados na [Especificação Visual — Barra de Mensagens](tui-spec-barras.md#barra-de-mensagens).

## Barra de Comandos

A barra de comandos é a última linha da tela (conforme [Dimensionamento e Layout](#dimensionamento-e-layout)). Exibe as ações acionáveis por teclado no contexto atual — o usuário nunca precisa adivinhar o que pode fazer. A âncora `F1 Ajuda` é fixa na extrema direita; ações de menor prioridade são removidas quando falta espaço.

Anatomia, dimensionamento, identidade visual, sistema de ações e eventos especificados na [Especificação Visual — Barra de Comandos](tui-spec-barras.md#barra-de-comandos).

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
| Campo de entrada (ativo) | Campo com cursor | Fundo `surface.input` + cursor real do terminal em `text.primary` + label em `accent.primary` **bold** |
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

A aplicação é operada inteiramente por teclado. Esta seção define a notação para representar teclas na documentação e na interface, as convenções semânticas que governam o significado de cada tecla, e a política de escopos e reservas. O mapeamento completo por tela está na [Especificação Visual](tui-specification.md).

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

