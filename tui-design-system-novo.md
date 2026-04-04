# Design System — Abditum TUI

> Fundações visuais e padrões transversais para o pacote `internal/tui`.
> Define princípios, tokens, estados e padrões que governam toda decisão de UI.

## Fronteira deste documento

Este documento define **fundações** e **padrões reutilizáveis** — o que cada peça visual é, como se comporta em abstrato e como peças se combinam em situações recorrentes.

A composição dessas peças em telas, wireframes e fluxos concretos pertence ao documento de especificação:
- [`tui-specification-novo.md`](tui-specification-novo.md) — telas, wireframes de componentes e fluxos visuais

### Teste de fronteira

> *Se eu trocar o nome do item concreto (ex: "Gmail" por qualquer outro segredo) e a regra continuar válida, é **padrão** — fica neste documento. Se a regra só faz sentido para aquela tela ou componente específico, é **composição** — vai para a especificação.*

> **Regra de governança:** toda decisão de UI/UX deste projeto deve ser compatível com os princípios definidos aqui. Em caso de conflito entre uma especificação local e um princípio, o princípio prevalece e a especificação deve ser ajustada.

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
- [Padrões](#padrões)
  - [Sobreposição](#sobreposição)
  - [Mensagens](#mensagens)
  - [Foco e Navegação](#foco-e-navegação)
  - [Mapa de Teclas](#mapa-de-teclas)
  - [Acessibilidade](#acessibilidade)
- [Anti-padrões](#anti-padrões)
  - [Segurança Visual](#segurança-visual)
  - [Estado e Feedback](#estado-e-feedback-1)
  - [Navegação e Teclado](#navegação-e-teclado)
  - [Diálogos e Confirmações](#diálogos-e-confirmações)
  - [Layout e Estrutura](#layout-e-estrutura)
  - [Tipografia e Cor](#tipografia-e-cor)
  - [Acessibilidade](#acessibilidade-1)
  - [Ciclo de Vida do Cofre](#ciclo-de-vida-do-cofre)

---

## O Terminal como Meio

O Abditum é uma aplicação TUI. As propriedades do terminal não são restrições a superar — são o material com o qual trabalhamos. Todo o design opera dentro deste perímetro.

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
| **Bordas** | `border.default` | Linhas que dividem painéis, bordas de janelas informativas (ajuda, seleção de itens, navegação de arquivos) | `#414868` <span style="color:#414868">██</span> | `#3a3a5c` <span style="color:#3a3a5c">██</span> |
| | `border.focused` | Borda do painel ativo, de janelas de entrada (senhas, textos) e de diálogos com severidade neutra. Diálogos com severidade não-neutra usam `semantic.*` — ver [Sobreposição](#sobreposição) | `#7aa2f7` <span style="color:#7aa2f7">██</span> | `#ff2975` <span style="color:#ff2975">██</span> |
| **Interação** | `accent.primary` | Barra de seleção na lista, cursor de navegação, botão principal de ação | `#7aa2f7` <span style="color:#7aa2f7">██</span> | `#ff2975` <span style="color:#ff2975">██</span> |
| | `accent.secondary` | Ícone de favorito (★), nomes de pastas na navegação de arquivos | `#bb9af7` <span style="color:#bb9af7">██</span> | `#00fff5` <span style="color:#00fff5">██</span> |
| **Semânticas** | `semantic.success` | Operação concluída com sucesso, configuração ligada (ON) | `#9ece6a` <span style="color:#9ece6a">██</span> | `#05ffa1` <span style="color:#05ffa1">██</span> |
| | `semantic.warning` | Alerta antes de ação permanente, aviso de bloqueio por tentativas erradas, prefixos de estado dirty (`✦ ✎ ✗`) | `#e0af68` <span style="color:#e0af68">██</span> | `#ffe900` <span style="color:#ffe900">██</span> |
| | `semantic.error` | Erro de operação, senha incorreta, borda de diálogos destrutivos | `#f7768e` <span style="color:#f7768e">██</span> | `#ff3860` <span style="color:#ff3860">██</span> |
| | `semantic.info` | Informação contextual | `#7dcfff` <span style="color:#7dcfff">██</span> | `#00b4d8` <span style="color:#00b4d8">██</span> |
| | `semantic.off` | Configuração desligada (OFF) | `#737aa2` <span style="color:#737aa2">██</span> | `#9999cc` <span style="color:#9999cc">██</span> |
| **Especiais** | `special.muted` | Texto esmaecido — uso pontual em contextos que precisam de cor apagada sem conotação semântica | `#8690b5` <span style="color:#8690b5">██</span> | `#666688` <span style="color:#666688">██</span> |
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

> **Regra prática:** Rounded é o único estilo de caixa adotado. Separadores são linhas; a interface evita boxes decorativos.

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

### Dimensionamento de diálogos

| Parâmetro | Valor | Justificativa |
|---|---|---|
| Largura mínima | 40 colunas | Cabe título + 2 ações na borda inferior |
| Largura máxima | 70 colunas ou 80% do terminal (o menor) | Evita que o diálogo pareça "colado" nas bordas |
| Padding interno horizontal | 2 colunas | Respiro visual sem desperdício |
| Padding interno vertical | 1 linha | Respiro entre borda e conteúdo |
| Posição | Centrado horizontal e verticalmente | Previsibilidade — o olho do usuário vai sempre ao centro |

**Mensagens longas em diálogos:**

Quando a mensagem do diálogo excede o espaço disponível, o dimensionamento segue três passos em sequência:

1. **Largura cresce** até o máximo (70 colunas ou 80% do terminal, o menor)
2. **Word-wrap** — a mensagem quebra em linhas respeitando limites de palavra
3. **Altura cresce** até o máximo (80% da altura do terminal)
4. Se ainda exceder, **scroll vertical** — ver [Scroll em diálogos](#scroll-em-diálogos) abaixo

O diálogo nunca ultrapassa os limites de largura e altura máximos. A mensagem sempre inicia sem scroll — o conteúdo é totalmente visível quando possível.

### Scroll em diálogos

Padrão transversal — aplica-se a qualquer diálogo ou componente modal cujo conteúdo excede a área visível.

**Indicadores visuais — borda direita:**

Quando há conteúdo fora da viewport, a borda direita comunica direção e posição do scroll usando três elementos:

- **Setas** (`↑` / `↓`) — substituem o primeiro e/ou último `│` da borda direita para indicar que há conteúdo acima/abaixo
- **Thumb** (`■`) — um único `│` é substituído por `■` para indicar a posição relativa no conteúdo. Posição calculada: `round(scroll_offset / max_scroll × (linhas_borda - 1))`

| Posição do scroll | Borda direita (primeiro `│`) | Borda direita (último `│`) | Thumb `■` |
|---|---|---|---|
| Topo (mais conteúdo abaixo) | `│` (normal) | `↓` | Próximo ao topo |
| Meio (conteúdo acima e abaixo) | `↑` | `↓` | Proporcional à posição |
| Final (mais conteúdo acima) | `↑` | `│` (normal) | Próximo à base |
| Sem scroll (tudo visível) | `│` (normal) | `│` (normal) | Não aparece |

> **Prioridade de renderização:** se o thumb `■` coincide com a posição de uma seta (`↑`/`↓`), a seta prevalece — a direção é mais importante que a posição.

| Elemento | Token | Atributo |
|---|---|---|
| Seta de scroll (`↑` / `↓`) | `text.secondary` | — |
| Thumb de posição (`■`) | `text.secondary` | — |

**Navegação:**

| Tecla | Efeito |
|---|---|
| `↑` / `↓` | Move uma linha |
| `PgUp` / `PgDn` | Move uma página (viewport − 1 linhas) |
| `Home` / `End` | Vai ao início / fim do conteúdo |

> As bordas superior e inferior do diálogo (título, ações) permanecem intactas — o scroll afeta apenas o conteúdo interno.

### Barra de mensagens

| Parâmetro | Valor |
|---|---|
| Altura | 1 linha fixa entre área de trabalho e barra de comandos |
| Anatomia | Borda `─` contínua; mensagem embutida após 2 espaços de padding esquerdo |
| Largura da borda | 100% da largura do terminal |
| Largura do texto | Largura do terminal − 2 (padding) − 2 (`─` direita mínima) |
| Truncamento | Com `…` quando o texto exceder o espaço disponível |
| Sem mensagem | Borda `─` contínua (separador visual permanente) |

---

## Ícones e Símbolos

Inventário completo dos caracteres Unicode usados pela interface.

### Critérios de seleção

A escolha de cada símbolo segue restrições práticas do terminal como meio:

- **BMP apenas (U+0000–U+FFFF).** Caracteres do Basic Multilingual Plane têm suporte consistente em terminais Windows (ConHost, Windows Terminal), macOS (Terminal.app, iTerm2) e Linux (gnome-terminal, Alacritty, kitty). Caracteres fora do BMP (emojis, suplementares) dependem de fontes e renderizadores que não controlamos.
- **Largura previsível.** Todos os símbolos ocupam exatamente 1 coluna terminal, exceto `<╡` (2 colunas, por composição). Símbolos de "largura ambígua" no Unicode (que podem ser renderizados como 1 ou 2 colunas dependendo do locale) são evitados.
- **Sem emojis.** Emojis ocupam 2 colunas, dependem de fontes coloridas e têm renderização inconsistente entre terminais — especialmente no Windows. São excluídos do inventário.
- **Sem Nerd Fonts.** Glifos de Nerd Fonts (~U+E000–U+F8FF, Private Use Area) só existem em fontes instaladas pelo usuário. O Abditum não pode assumir que essas fontes estão disponíveis.
- **Semântica sobre estética.** Cada símbolo é escolhido pelo significado que comunica, não pela aparência. `✗` (exclusão) e `✕` (erro) são visualmente similares mas semanticamente distintos — ambos permanecem no inventário porque servem papéis diferentes.

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

- Foco de painel é indicado pelo separador vertical (ver [Foco e Navegação](#foco-e-navegação)) e pela barra de comandos — não por borda ao redor do painel.
- TUI não tem estado "pressionado"; confirmação vem por mudança de contexto ou mensagem.
- Transições são instantâneas. A única animação prevista é o spinner `MsgBusy`.

---

## Padrões

Padrões são regras de comportamento transversais — aplicam-se a múltiplas telas e componentes. Os documentos de especificação consomem estes padrões ao definir componentes e fluxos concretos.

---

### Sobreposição

Elementos sobrepostos (modais, diálogos, seletores) seguem regras uniformes de apresentação e interação.

**Apresentação:**

- Centralizados horizontal e verticalmente sobre o conteúdo (ver [Dimensionamento de diálogos](#dimensionamento-de-diálogos))
- Estilo de borda: Rounded (`╭╮╰╯│─`)
- Fundo interno: `surface.raised`
- O conteúdo abaixo permanece visível, mas inativo — sem escurecimento de overlay

**Foco e pilha:**

- Apenas o elemento do topo recebe input; os inferiores permanecem montados, porém congelados
- A barra de comandos reflete os atalhos do elemento do topo enquanto ele estiver ativo
- Ao fechar o elemento do topo, o foco retorna ao elemento anterior na pilha (ou ao conteúdo base)

**Navegação padrão:**

| Tecla | Comportamento |
|---|---|
| `Enter` | Aciona a opção default |
| `Esc` | Aciona a opção de cancelamento; se não existir, fecha o elemento |
| Atalho da opção | Aciona diretamente a opção correspondente |

**Modelo bidimensional: Intenção × Severidade**

Diálogos são classificados por duas dimensões ortogonais. Qualquer intenção pode combinar com qualquer severidade. A referência descritiva usa a forma: "diálogo de _intenção_ com severidade _severidade_" — por exemplo, "diálogo de confirmação com severidade destrutiva".

**Intenção** — o que o diálogo pede ao usuário:

| Intenção | Descrição | Barra de ações |
|---|---|---|
| Confirmação | Requer escolha entre opções (ex: Excluir / Cancelar) | Ação default + opcionais + cancelar |
| Reconhecimento | Requer apenas que o usuário tome ciência | `Enter OK` (única ação) |

**Severidade** — governa o tratamento visual (borda, símbolo, cor da tecla default):

| Severidade | Símbolo | Token de borda | Token da tecla default | Quando usar |
|---|---|---|---|---|
| Destrutivo | `⚠` | `semantic.warning` | `semantic.error` | Ação irreversível ou com perda de dados |
| Erro | `✕` | `semantic.error` | `accent.primary` | Falha ocorrida, condição irrecuperável |
| Alerta | `⚠` | `semantic.warning` | `accent.primary` | Situação importante mas recuperável |
| Informativo | `ℹ` | `semantic.info` | `accent.primary` | Informação que requer atenção |
| Neutro | — | `border.focused` | `accent.primary` | Operação rotineira, sem urgência |

> **Nota:** severidades Destrutivo e Alerta compartilham o símbolo `⚠` e o token de borda `semantic.warning`. A distinção visual está na tecla default: `semantic.error` (vermelho) para destrutivo, `accent.primary` para alerta. Isso reforça que o perigo está na *ação*, não apenas na *situação*.

**Anatomia comum:**

Todo diálogo — de decisão ou funcional — segue a mesma estrutura de moldura:

```
╭── ⚠  Título ────────────────────╮  ← borda superior: símbolo + título em bold
│                                  │
│  (conteúdo interno do diálogo)   │
│                                  │
╰── S Ação ────────── Esc Cancelar ╯  ← borda inferior: default à esquerda, cancelar à direita
```

Regras da moldura:

- **Borda superior** contém o título embutido, precedido pelo símbolo de severidade quando aplicável (`⚠`, `ℹ`, `✕`). Severidade Neutro não usa símbolo
- **Borda inferior** contém apenas ações de confirmação e cancelamento, alinhadas à direita
- **Ordem das ações:** a ação default (associada a `Enter`) fica sempre na primeira posição (mais à esquerda); a ação de cancelamento (`Esc`) fica sempre na última posição (mais à direita, junto à borda). Em diálogos com 3 ou mais ações, as intermediárias ficam entre default e cancelar
- **Ação default** (associada a `Enter`): tecla + label em **bold**, coloridos com o token de destaque da severidade (ver tabela de severidades acima) — visualmente distinta das demais
- **Demais ações**: tecla + label na cor da borda, sem bold
- **Borda e título** usam o mesmo token — definido pela tabela de tipos semânticos
- **Ações internas** (revelar senha, alternar campo, expandir diretório) aparecem exclusivamente na barra de comandos — não na borda do diálogo
- **Teclas de navegação** (↑↓, →, ←, Tab) são de conhecimento amplo e não aparecem na borda
- A **barra de comandos** exibe apenas as ações internas do diálogo (ex: Tab entre campos, revelar senha). Ações de confirmação e cancelamento já estão na borda inferior do diálogo — não são duplicadas na barra

**Diálogos de decisão (confirmação e reconhecimento):**

Diálogos de decisão pedem uma ação do usuário — seja uma escolha entre opções (confirmação) ou o reconhecimento de uma informação. O conteúdo interno é uma mensagem + contexto:

- **Mensagem** em `text.primary`; nomes de itens referenciados em **bold**
- A severidade define o tratamento visual da moldura e da tecla default conforme a tabela acima

> A matriz completa de combinações Intenção × Severidade com wireframes ilustrativos está documentada na [especificação de telas](tui-specification-novo.md#diálogos-de-decisão).

**Diálogos funcionais:**

PasswordEntry, PasswordCreate, FilePicker e Help oferecem uma função específica em vez de uma decisão sim/não. Compartilham a mesma moldura (título na borda superior, ações na borda inferior) mas diferem no conteúdo interno.

Regras específicas:

- **Título** sem símbolo semântico (não há urgência)
- **Borda** em `border.focused` para diálogos que recebem entrada de texto; `border.default` para diálogos de consulta

> A anatomia interna de cada diálogo funcional está documentada na [especificação de telas](tui-specification-novo.md#diálogos-funcionais).

**Ação default condicional (diálogos funcionais):**

Em diálogos funcionais com campos de entrada, a ação default (confirmar/continuar) fica **inativa** até que as condições mínimas do diálogo sejam satisfeitas. Isso se aplica a **todos os meios de acionamento**: `Enter`, mouse, atalho de teclado.

| Estado | Estilo da ação default | Comportamento |
|---|---|---|
| Condições **não** satisfeitas | `text.disabled`, sem bold | Tecla/mouse ignorados silenciosamente |
| Condições satisfeitas | Token normal da severidade + **bold** | Tecla/mouse acionam a ação |

A especificação de cada diálogo funcional documenta suas condições em uma tabela dedicada.

> A ação de cancelamento (`Esc`) permanece sempre ativa — o usuário pode abandonar o diálogo a qualquer momento.

**Barra de mensagens em diálogos:**

Diálogos funcionais usam a barra de mensagens para comunicar dicas e erros de validação. Diálogos de decisão (confirmação/reconhecimento) **não** usam a barra — a mensagem completa está no corpo do diálogo.

Ciclo de vida da barra durante um diálogo funcional:

| Momento | Conteúdo da barra | Tipo |
|---|---|---|
| Diálogo abre | Dica contextual do primeiro campo com foco | Dica de campo (`•` italic) |
| Foco entra em campo (branco ou válido) | Dica descritiva sobre o campo | Dica de campo (`•` italic) |
| Foco entra em campo (com valor inválido) | Mensagem de erro explicando a invalidação | Erro (`✕` bold, TTL 5s) |
| Tentativa de confirmar com validação falha | Mensagem de erro; diálogo permanece aberto | Erro (`✕` bold, TTL 5s) |
| Diálogo fecha (confirmação ou cancelamento) | Barra é limpa | — |

> **Separação de responsabilidade:** mensagens pós-fechamento (ex: "◐ Criando cofre…", "✓ Cofre aberto", "Operação cancelada") são responsabilidade do orquestrador — não do diálogo.

---

### Mensagens

A aplicação comunica feedback ao usuário por meio de uma mensagem exibida na barra de mensagens. Uma mensagem por vez — nova mensagem substitui a anterior imediatamente. Não há fila nem pilha.

**Posição:** sobreposta à última linha da área de trabalho — não reserva linha própria.

**Largura:** ~95% da largura do terminal. Trunca com `…` se necessário.

**Formato:** `<símbolo> <texto>` — exatamente 1 espaço entre o símbolo e o início do texto. O símbolo ocupa sempre 1 coluna; o espaço seguinte é fixo e não varia por tipo de mensagem.

**Tipos de mensagem:**

| Tipo | Símbolo | Token | Atributo |
|---|---|---|---|
| Sucesso | `✓` | `semantic.success` | — |
| Informação | `ℹ` | `semantic.info` | — |
| Alerta | `⚠` | `semantic.warning` | — |
| Erro | `✕` | `semantic.error` | **bold** |
| Ocupado (spinner) | `◐ ◓ ◑ ◒` | `accent.primary` | — |
| Dica de campo | `•` | `text.secondary` | *italic* |
| Dica de uso | `•` | `text.secondary` | *italic* |

> **Token se aplica à mensagem inteira** — símbolo e texto usam o mesmo token de cor. Não há distinção de cor entre o símbolo e o conteúdo textual dentro de uma mesma mensagem.

**Ciclo de vida:**

O ciclo de vida de cada mensagem é controlado pelo orquestrador que a emite. A tabela abaixo define os **defaults recomendados** — o caller pode sobrescrever TTL e trigger de dismissal conforme o contexto.

| Tipo | TTL padrão | Dismissal padrão |
|---|---|---|
| Sucesso | 5 s | Expiração |
| Informação | 5 s | Expiração |
| Alerta | 5 s | Expiração |
| Erro | 5 s | Expiração |
| Ocupado (spinner) | Sem TTL | Substituição explícita por Sucesso, Erro ou Alerta ao término da operação |
| Dica de campo | Permanente | Troca de campo ou substituição por outro tipo |
| Dica de uso | Permanente | Substituição por qualquer outro tipo |

**Regras de comportamento:**

- **Ocupado** spinner avança 1 frame/segundo sincronizado com tick global.

> O ciclo de vida da barra em diálogos funcionais (mensagem de contexto ao abrir, dica por campo, limpeza ao fechar) é contrato do orquestrador — documentado em [Barra de mensagens em diálogos](#sobreposição).

---

### Foco e Navegação

O modelo de foco define como o usuário percebe e alterna entre áreas interativas da interface.

**Alternância com Tab:**

`Tab` é contextual — o comportamento depende do estado do painel:

| Contexto | `Tab` | `Shift+Tab` |
|---|---|---|
| Modo leitura | Foco → próximo painel (árvore ↔ detalhe) | Foco → painel anterior |
| Modo edição (detalhe) | Foco → próximo campo editável | Foco → campo anterior |
| Modo edição, último campo | Foco → painel esquerdo (árvore) | Foco → campo anterior |
| Modo edição, primeiro campo | Foco → próximo campo | Foco → painel esquerdo (árvore) |

O ciclo entre painéis é circular. Painéis vazios ou sem conteúdo interativo são pulados.

**Indicação de foco:**

- A área de trabalho não tem painéis com borda — existe apenas um separador vertical (`│`) em `border.default` entre a árvore/lista (esquerda) e o detalhe (direita)
- **Conector `<╡`:** na linha do item selecionado na árvore, o separador `│` é substituído por `<╡` em `accent.primary` — amarra visualmente o item ao conteúdo detalhado à direita
- O painel ativo é identificado pela **barra de comandos**, que exibe as ações do painel com foco

**Teclado primeiro, mouse sempre:**

- Teclas de navegação direcional (`↑↓←→`) são o caminho primário para listas e árvores
- `Home` / `End` navegam ao primeiro / último item visível
- Toda ação acionável por teclado deve ser descobrível e executável também por mouse

**Campos de entrada de texto:**

- Campos não possuem borda — a área digitável é delimitada por um fundo `surface.input` (tom rebaixado em relação ao `surface.raised` do diálogo)
- Label do campo ativo em `accent.primary` + **bold**; label dos campos inativos em `text.secondary`
- Foco indicado pela presença do cursor `▌` em `text.primary` dentro do fundo `surface.input`
- Placeholder em `text.secondary` + *italic* — desaparece ao digitar
- Erro de validação: exibido na barra de mensagens (tipo Erro), não inline — os formulários são simples o suficiente para mostrar um erro por vez
- Em **NO_COLOR**: o fundo `surface.input` pode ser perdido; o cursor + label em **bold** permanecem como indicadores de foco suficientes

---

### Mapa de Teclas

Esta seção define a **política de atribuição de teclas** — como atalhos são organizados, quais regras regem conflitos e quais teclas têm significado global. O mapeamento completo por tela está na [especificação de telas](tui-specification-novo.md).

**Política de escopos:**

Cada tecla pertence a exatamente um escopo. Escopos mais específicos sobrepõem os mais gerais quando ambos estão ativos (ex: um diálogo sobrepõe a área de trabalho).

| Escopo | Descrição | Exemplo |
|---|---|---|
| **Global** | Funciona em qualquer contexto da aplicação | `F12` tema, `Esc` cancelar |
| **Área de trabalho** | Funciona quando a área de trabalho tem foco (sem diálogo aberto) | `^S` salvar, `^F` buscar |
| **Diálogo** | Funciona apenas enquanto um diálogo está no topo da pilha | `Enter` confirmar, atalho da opção |
| **Tela-específico** | Funciona apenas em uma tela ou painel particular | F-keys de segredo, pasta, modelo |

**Regras de consistência:**

- `Enter` sempre avança ou aprofunda; `Esc` sempre retrocede ou abandona. O vetor direcional é consistente mesmo quando a ação concreta varia por escopo (ver [Princípios — Consistência de interação](#experiência)).
- Se uma tecla precisa ter significado diferente em dois contextos, isso deve ser documentado e justificado na especificação.
- Teclas de navegação universais (`↑↓←→`, `Tab`, `Home`, `End`) não aparecem na barra de comandos — são senso comum em TUI. Exceção: diálogos exibem opções explicitamente.

**Teclas globais (definidas neste documento):**

| Tecla | Ação | Referência |
|---|---|---|
| `Enter` | Avança / aprofunda: confirma em diálogos, seleciona/expande na árvore, ativa/confirma edição de campo | [Sobreposição](#sobreposição), Princípios |
| `Esc` | Retrocede / abandona: fecha modal, cancela edição, sai de modo (busca, edição) | [Sobreposição](#sobreposição), Princípios |
| `Tab` / `Shift+Tab` | Navega entre painéis (modo leitura) ou campos (modo edição) | [Foco e Navegação](#foco-e-navegação) |
| `↑` `↓` | Navegação direcional em listas e árvores | [Foco e Navegação](#foco-e-navegação) |
| `←` `→` | Expandir/recolher pastas; alternar opções em diálogos | [Foco e Navegação](#foco-e-navegação) |
| `Home` / `End` | Primeiro / último item visível | [Foco e Navegação](#foco-e-navegação) |
| `PgUp` / `PgDn` | Scroll por página (viewport − 1) em conteúdo com scroll | [Scroll em diálogos](#scroll-em-diálogos) |
| `F12` | Alternar tema (Tokyo Night ↔ Cyberpunk) | [Temas](#temas) |
| `F1` | Abrir modal de Ajuda | — |

**Teclas de área de trabalho (definidas neste documento):**

| Tecla | Ação | Condição |
|---|---|---|
| `F2` | Modo Cofre (aba) | Só com cofre aberto |
| `F3` | Modo Modelos (aba) | Só com cofre aberto |
| `F4` | Modo Config (aba) | Só com cofre aberto |

> O mapeamento de F-keys por contexto funcional (segredos, pastas, modelos, cofre) é definido na especificação de telas.

**Ações ocultas da barra de comandos:**

Algumas ações globais não aparecem na barra de comandos — são registradas no ActionManager com o atributo "oculto da barra". Essas ações contínuam disponíveis por teclado e aparecem no modal de Ajuda (`F1`).

| Tecla | Ação | Justificativa |
|---|---|---|
| `F12` | Alternar tema | Ação pontual, sem necessidade de visibilidade permanente |

---

### Acessibilidade

#### NO_COLOR e modo monocromático

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

---

## Anti-padrões

Anti-padrões documentam o que **não deve ser feito** na interface do Abditum. Cada item lista o padrão incorreto, por que viola os princípios do design system, e qual consequência concreta afeta o usuário.

> **Regra de uso:** toda decisão de implementação que contradiga um item desta seção deve ser justificada explicitamente na especificação que a adota. Sem justificativa documentada, o anti-padrão prevalece como proibição.

---

### Segurança Visual

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Revelação Passiva de Sensível** *(Crítico)* | Campos sensíveis revelam por foco/Tab sem ação explícita | Dado sensível exposto sem percepção do usuário |
| **Máscara Apenas Visual** *(Alto)* | `••••••••` exibido mas copiável sem feedback | Proteção ilusória; dado sensível exposto em clipboard |
| **Campo Sensível Indistinguível** *(Alto)* | Campos sensíveis e comuns têm mesma aparência | Revelação acidental ou proteção ignorada |
| **Countdown Invisível** *(Médio)* | Cópia bem-sucedida mas sem indicação de TTL da clipboard | Usuário não sabe se o dado ainda está disponível |
| **Exportação Sem Cerimônia** *(Crítico)* | Exportação (arquivo não criptografado) com tratamento de ação rotineira | Usuário exporta para local inseguro sem compreender risco |
| **Dirty State Apenas Global** *(Crítico)* | Indicador `•` só no cabeçalho, sem `✦ ✎ ✗` por item | Usuário não consegue auditar o que será salvo |

---

### Estado e Feedback

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Silêncio Após Operação Crítica** *(Alto)* | Salvar/bloquear/exportar sem mensagem de confirmação | Indistinguível de falha silenciosa |
| **Spinner Sem Resolução** *(Médio)* | `◐ Carregando…` nunca substituído por `✓` ou `✕` | Usuário não sabe se pode interagir |
| **Fila de Mensagens** *(Médio)* | Múltiplas mensagens enfileiradas ou sobrepostas | Falta correspondência entre ação e mensagem |
| **Contador Defasado** *(Médio)* | Contagem de segredos não atualiza em tempo real | Decisões baseadas em dados incorretos |
| **Modo Ativo Sem Indicador** *(Alto)* | Busca/edição/reordenação sem indicador persistente na barra | Usuário digita no modo errado |

---

### Navegação e Teclado

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Enter Polissêmico** *(Médio)* | `Enter` expande ou edita sem distinção visual clara | Edição acidental ao tentar visualizar |
| **Cursor ao Topo Após Operação** *(Médio)* | Exclusão/reordenação retorna cursor ao topo da lista | Re-navegação obrigatória; experiência frustrante |
| **Setas com Semântica Dupla Invisível** *(Baixo)* | `←`/`→` expandem pastas E navegam diálogos sem indicador | Expansão/fechamento/navegação acidental |

---

### Diálogos e Confirmações

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Fadiga de Confirmação** *(Médio)* | Toda ação pede confirmação, inclusive benignas | Usuário aprende a apertar Enter reflexivamente |
| **Uniformidade de Risco Visual** *(Alto)* | "Excluir 47 segredos" e "Renomear pasta" têm mesma aparência | Usuário não calibra gravidade da ação |
| **Pilha de Modais Sem Profundidade** *(Médio)* | Modal abre modal abre modal sem indicação | Desorientação; fechamento acidental com Esc repetido |
| **Ação Default Ausente** *(Médio)* | Ação default desaparece quando inativa | Usuário não sabe o que falta preencher |
| **Confirmação Assimétrica** *(Crítico)* | "Salvar e Sair" pede dupla confirmação; "Descartar" não | Incentivo perverso aumenta perdas de dados |

---

### Layout e Estrutura

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Layout Saltitante** *(Médio)* | Elementos fixos reposicionam por conteúdo variável | Perda de ancoragem após cada seleção |
| **Over-boxing** *(Baixo)* | Toda seção envolta em borda; grade de boxes | Ruído visual; ambiguidade estrutural/decoração |
| **Informação Densa** *(Baixo)* | Nomes, labels, valores com mesmo peso tipográfico | Dificulta localização e varredura rápida |
| **Truncamento Ausente** *(Alto)* | Texto longo cortado sem `…` ou transborda | Confusão de identidade; layout corrompido |
| **Pasta Virtual Indistinguível** *(Médio)* | Favoritos parecem pastas normais | Usuário tenta criar item e recebe erro inesperado |
| **Caractere Largura Dupla** *(Alto)* | Símbolos ambíguos sem contabilização de colunas | Bordas não fecham; separadores desalinhados |
| **Resize Sem Recálculo** *(Crítico)* | Layout não atualiza ao redimensionar terminal | Interface inutilizável até reiniciar |
| **Conteúdo Sem Scroll** *(Alto)* | Painel/diálogo corta conteúdo sem `↑↓` nem thumb | Campos/ações finais inacessíveis |
| **Campo Maior que Área** *(Médio)* | Valor longo sem truncamento/scroll horizontal | Sobrescrita de labels; valor ilegível |
| **Sangramento ANSI** *(Alto)* | Estilo não resetado contamina conteúdo seguinte | Cores/estilos vaza para componentes vizinhos e shell |
| **Cálculo de Largura Errado** *(Alto)* | `len(s)` em vez de largura visual (ANSI excluído) | Desalinhamento de bordas e truncamento |
| **Layout Colapsa Vazio** *(Médio)* | Painel sem conteúdo desaparece | Separador desaparece; proporção quebra ao preenchimento |
| **Indicador Causa Deslocamento** *(Médio)* | `✦ ✎ ★` não em coluna fixa | Nomes "pulam" horizontalmente ao marcar/desmarcar |
| **Spinner com Largura Variável** *(Baixo)* | Frames do spinner ocupam 1 ou 2 colunas | Mensagem pisca horizontalmente |
| **Contador com Largura Dinâmica** *(Baixo)* | Número muda de 9 para 10 dígitos | Coluna inteira se desloca |
| **Artefato de Render Anterior** *(Alto)* | Caracteres/cores do frame antigo permanecem | Bordas flutuam; campos extras visíveis |
| **Última Linha Causa Scroll** *(Médio)* | Escrever em `(linhas, colunas)` aciona scroll | Barra de comandos "cai"; layout deslocado |
| **Cursor Desalinhado** *(Alto)* | Cursor em coluna errada durante edição (bytes vs runes) | Backspace apaga caractere errado |
| **Campo Edição Sem Scroll H** *(Alto)* | Campo longo truncado ou overflow sem scroll | Usuário não vê valor completo |

---

### Tipografia e Cor

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Bold Inflacionado** *(Baixo)* | Bold aplicado a tudo (títulos, labels, nomes, ações) | Hierarquia colapsa; interface gritante |
| **Token Semântico Decorativo** *(Baixo)* | `semantic.success` / `warning` usado para ornamento | Usuário para de confiar nos indicadores |
| **Cor Hardcoded** *(Alto)* | Hex literal em vez de tokens de tema | Segundo tema nunca funciona corretamente |
| **Italic Sem Cor** *(Baixo)* | Hints em italic apenas, sem `text.secondary` | Indistinguível do conteúdo em terminais sem italic |

---

### Acessibilidade

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Estado Apenas por Cor** *(Alto)* | `✦ ✎ ✗` não usados; apenas color diferencia | Em NO_COLOR, informação de estado desaparece |
| **Área de Clique Minúscula** *(Baixo)* | `<╡` (2 colunas dinâmicas) único alvo clicável | Mouse inutilizável; cada clique requer precisão |
| **Erro Técnico Exposto** *(Médio)* | Mensagens internas: "unexpected JSON at 1247" | Exposição de caminho de arquivo; usuário confuso |

---

### Ciclo de Vida do Cofre

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Auto-save Silencioso** *(Alto)* | Alteração de senha salva automaticamente sem feedback | Usuário acredita que é reversível |
| **Conflito de Arquivo Minimizado** *(Crítico)* | Arquivo modificado externamente sobrescrito sem aviso | Dados de outra sessão/backup destruídos |
| **Re-autenticação Durante Sessão** *(Alto)* | Senha mestra solicitada novamente em salvamento/exportação | Fricção ilegítima; treina digitação irrefletida |
| **Exclusão Desaparece Imediatamente** *(Crítico)* | Item marcado para exclusão some sem `✗` + strikethrough | Usuário crê ter deletado permanentemente |
| **Importação Sem Prévia de Impacto** *(Crítico)* | Mesclagem executada sem mostrar o que será sobrescrito | Perda de dados não intencionada |


