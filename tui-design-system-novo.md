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

Anti-padrões documentam o que **não deve ser feito** na interface do Abditum. Cada item descreve um padrão de implementação incorreto, por que viola os princípios deste design system ou os requisitos do produto, e qual consequência concreta ele gera para o usuário.

> **Regra de uso:** toda decisão de implementação que contradiga um item desta seção deve ser justificada explicitamente na especificação que a adota. Sem justificativa documentada, o anti-padrão prevalece como proibição.

---

### Segurança Visual

#### Revelação Passiva de Campo Sensível *(Crítico)*
**O que é:** campos sensíveis revelam seu valor por foco, hover ou Tab — sem ação explícita do usuário.
**Por que é errado:** viola o princípio *Segurança como experiência*. Foco é navegação, não autorização. A revelação exige gesto intencional.
**Consequência:** dado sensível exposto em ambiente público sem que o usuário perceba.

---

#### Mascaramento Simbólico Sem Substância *(Alto)*
**O que é:** `••••••••` é exibido na lista ou detalhe, mas o campo permanece copiável por atalho sem nenhum feedback — a máscara é puramente visual.
**Por que é errado:** cria falsa sensação de proteção. O risco real fica invisível ao próprio usuário.
**Consequência:** dado sensível copiado e exposto em ambiente compartilhado sem percepção do usuário.

---

#### Paridade Visual Entre Campo Sensível e Campo Comum *(Alto)*
**O que é:** campos sensíveis e campos comuns têm a mesma aparência no painel de detalhe — sem `◉`, sem máscara, sem distinção de estilo.
**Por que é errado:** o usuário não sabe quais campos requerem atenção antes de revelar. A distinção visual é parte da proteção.
**Consequência:** revelação acidental de campo sensível que o usuário não sabia ser sensível; ou proteção ignorada por ser imperceptível.

---

#### Countdown Invisível da Clipboard *(Médio)*
**O que é:** ao copiar um campo, a mensagem de sucesso aparece e desaparece, mas não há indicação de quanto tempo resta até a limpeza automática.
**Por que é errado:** o requisito define TTL configurável para clipboard. Sem indicação, o usuário não sabe se o dado ainda está disponível — ou se já foi limpo.
**Consequência:** dado sensível persiste na clipboard além do esperado; ou o usuário tenta colar em outra aplicação e estranha o comportamento sem entender o motivo.

---

#### Exportação Sem Cerimônia *(Crítico)*
**O que é:** a operação de exportar — que produz um arquivo **não criptografado** com todo o conteúdo do cofre — tem o mesmo tratamento visual de qualquer ação rotineira.
**Por que é errado:** o requisito exige aviso explícito sobre riscos. O design system tem severidade `destrutivo` precisamente para este caso — irreversível em termos de exposição de dados.
**Consequência:** usuário exporta inadvertidamente para local inseguro sem compreender as implicações.

---

#### Dirty State Apenas no Cabeçalho *(Crítico)*
**O que é:** o indicador de alterações não salvas (`•`) aparece somente no cabeçalho global, sem os prefixos `✦ ✎ ✗` por item na árvore.
**Por que é errado:** viola o princípio *Estado sempre visível*. O usuário sabe que *algo* mudou, mas não *o quê*. Antes de salvar um cofre com dados sensíveis, ele deve poder auditar o que está prestes a persistir.
**Consequência:** salvamento de exclusão acidental não percebida; ou perda de modificação por falta de rastreabilidade visual por item.

---

### Estado e Feedback

#### Silêncio Após Operação Crítica *(Alto)*
**O que é:** salvar, bloquear ou exportar o cofre não produz nenhuma mensagem de confirmação. A operação "simplesmente acontece".
**Por que é errado:** viola o princípio *Feedback imediato*. Para um cofre de senhas, ausência de feedback após salvar é indistinguível de falha silenciosa.
**Consequência:** usuário fecha o terminal acreditando que salvou quando não salvou; ou repete a operação desnecessariamente por insegurança.

---

#### Spinner Sem Resolução *(Médio)*
**O que é:** a mensagem `◐ Carregando…` é exibida durante uma operação mas nunca é substituída por `✓ Sucesso` ou `✕ Erro` ao término.
**Por que é errado:** o design system define que `MsgBusy` é sempre substituído explicitamente ao terminar. Spinner sem resolução é ambiguidade de estado — o usuário não sabe se pode interagir.
**Consequência:** interação prematura antes da operação concluir; ou espera indefinida por algo que já terminou.

---

#### Fila de Mensagens *(Médio)*
**O que é:** múltiplas mensagens são enfileiradas e aparecem em sequência na barra de mensagens, ou aparecem sobrepostas.
**Por que é errado:** o design system é explícito: *"uma mensagem por vez — nova mensagem substitui a anterior imediatamente. Não há fila nem pilha."* Fila cria ordem temporal que o usuário não está preparado para rastrear.
**Consequência:** mensagens de contextos diferentes sem correspondência clara com as ações que as geraram.

---

#### Contador de Segredos Defasado *(Médio)*
**O que é:** a contagem de segredos de cada pasta na árvore não atualiza em tempo real — o número só muda ao salvar ou recarregar.
**Por que é errado:** o requisito especifica que a contagem deve refletir segredos ativos (não marcados para exclusão) em tempo real de sessão. Contador defasado mente sobre o estado atual.
**Consequência:** decisões tomadas com base em dados incorretos; perda de confiança na interface.

---

#### Modo Ativo Sem Indicador *(Alto)*
**O que é:** a aplicação entra em modo de busca, edição ou reordenação sem que nenhum indicador visual persista — barra de comandos e cabeçalho não mudam.
**Por que é errado:** em TUI, modo ativo é contexto invisível por natureza. Sem indicador persistente, o usuário digita no modo errado.
**Consequência:** dados sobrescritos acidentalmente; ações disparadas sem intenção; desorientação em cofres com muitos itens.

---

### Navegação e Teclado

#### Enter Polissêmico Sem Distinção de Modo *(Médio)*
**O que é:** `Enter` sobre uma pasta expande a pasta; sobre um segredo abre o modo de edição — mas o cursor não diferencia visualmente "navegação" de "edição prestes a iniciar".
**Por que é errado:** as consequências diferem radicalmente: expandir pasta é reversível com uma tecla; entrar em modo de edição pode acionar mudanças. A distinção deve ser visual, não apenas semântica.
**Consequência:** abertura acidental do modo de edição ao tentar apenas visualizar; modificação de dados sem intenção.

---

#### Cursor Saltando ao Topo Após Operação *(Médio)*
**O que é:** após excluir, mover ou reordenar um item, o cursor retorna ao topo da lista independentemente de onde estava.
**Por que é errado:** viola o princípio *Estabilidade espacial*. Operações devem preservar o contexto de posição ou mover o cursor ao item logicamente mais próximo.
**Consequência:** re-navegação obrigatória após cada operação; experiência lenta e frustrante em cofres com muitas pastas e segredos.

---

#### Setas com Semântica Dupla Imperceptível *(Baixo)*
**O que é:** `←` e `→` expandem/recolhem pastas na árvore e também navegam entre botões em diálogos — sem indicador visual do comportamento ativo.
**Por que é errado:** o design system reconhece a ambiguidade contextual das setas. Sem indicador de modo (barra de comandos diferenciada, indicador de foco), o usuário não sabe qual comportamento esperar.
**Consequência:** tentativa de expandir pasta fecha um diálogo acidentalmente; ou tentativa de navegar entre botões dispara recolhimento de pasta ao fundo.

---

### Diálogos e Confirmações

#### Fadiga de Confirmação *(Médio)*
**O que é:** toda ação — inclusive operações benignas como renomear um campo ou favoritar um segredo — exibe um diálogo de confirmação com `Enter / Esc`.
**Por que é errado:** quando toda ação pede confirmação, o usuário aprende a apertar `Enter` reflexivamente. A confirmação que deveria proteger ações graves (exportar, excluir pasta com n segredos) fica com o mesmo peso visual de trocar um nome.
**Consequência:** confirmação de destruição de dados pelo mesmo reflexo condicionado usado para confirmar renomeação.

---

#### Uniformidade de Risco Visual *(Alto)*
**O que é:** o diálogo de "Excluir pasta com 47 segredos" e o de "Renomear pasta" têm a mesma aparência — mesma borda, mesmo símbolo (ou ausência), mesma cor da ação default.
**Por que é errado:** o design system define explicitamente o modelo Intenção × Severidade com tratamento visual diferenciado. Destrutivo usa `semantic.warning` de borda e `semantic.error` na tecla default. Rotineiro usa `border.focused` neutro.
**Consequência:** usuário não calibra a gravidade da ação a partir de pistas visuais; todas parecem igualmente reversíveis.

---

#### Pilha de Modais Sem Indicador de Profundidade *(Médio)*
**O que é:** uma operação abre um diálogo que ao ser confirmado abre outro, que por sua vez abre um terceiro — sem indicação de profundidade na pilha nem de como encerrar todo o fluxo de uma vez.
**Por que é errado:** `Esc` fecha apenas o elemento do topo. Com pilha de 3 níveis, o usuário precisa de 3 `Esc` sem saber quantos faltam.
**Consequência:** desorientação; fechamento acidental do contexto base ao pressionar `Esc` além da conta.

---

#### Ação Default Ausente Quando Inativa *(Médio)*
**O que é:** quando as condições do diálogo não estão satisfeitas (ex: campo de senha vazio), a ação default desaparece completamente da borda inferior.
**Por que é errado:** o design system especifica `text.disabled` sem bold — a ação deve estar presente mas inativa. Remover o elemento retira a âncora visual que orienta o usuário sobre qual ação estará disponível quando as condições forem satisfeitas.
**Consequência:** o usuário não sabe se falta preencher algo, se preencheu errado, ou se a ação simplesmente não existe naquele contexto.

---

#### Confirmação Assimétrica ao Sair *(Crítico)*
**O que é:** "Salvar e Sair" pede confirmação dupla, mas "Descartar e Sair" executa com um único `Enter` — a ação mais destrutiva tem menos fricção que a ação segura.
**Por que é errado:** fricção deve ser proporcional ao risco. Descartar é irreversível; salvar é a ação conservadora. Hierarquia de confirmação invertida.
**Consequência:** usuário que quer salvar passa por mais passos que o usuário que vai descartar — incentivo perverso que aumenta perdas acidentais de dados.

---

### Layout e Estrutura

#### Layout Saltitante por Conteúdo Variável *(Médio)*
**O que é:** ao selecionar um segredo com 2 campos, depois um com 8 campos, o painel de detalhe reposiciona elementos fixos (título do segredo, barra de comandos) para acomodar os campos extras.
**Por que é errado:** viola o princípio *Estabilidade espacial*. Conteúdo variável cresce dentro da área de trabalho — não empurra as zonas fixas de layout.
**Consequência:** o usuário perde ancoragem após cada seleção; precisa reaprender a posição dos elementos a cada troca de segredo.

---

#### Over-boxing *(Baixo)*
**O que é:** cada seção da interface (árvore, detalhe, campos, observação) é envolta em caixa própria com borda `╭╮╰╯` — a tela vira um grid de caixas dentro de caixas.
**Por que é errado:** o design system é explícito: *"bordas aparecem apenas em modais e separadores. Painéis são organizados por espaço, alinhamento e hierarquia tipográfica."* Over-boxing transplanta layout CSS para o terminal e falha em ambos os meios.
**Consequência:** ruído visual extremo; o usuário não distingue bordas estruturais (modais) de bordas decorativas; em terminais com suporte parcial, metade das bordas renderiza incorretamente.

---

#### Informação Densa Sem Hierarquia *(Baixo)*
**O que é:** no painel de detalhe, nome do segredo, labels de campos, valores e observação aparecem com o mesmo peso tipográfico e espaçamento.
**Por que é errado:** o design system define bold para títulos, `text.secondary` para labels, `◉` para campos reveláveis. Sem hierarquia, o usuário varre o painel inteiro para encontrar o campo que quer.
**Consequência:** tempo elevado para localizar e copiar o campo desejado; interface percebida como lista plana, não formulário estruturado.

---

#### Truncamento Silencioso *(Médio)*
**O que é:** nomes longos de segredo, pasta ou campo são cortados na borda da coluna disponível sem o caractere `…`.
**Por que é errado:** o design system define `…` como pictograma obrigatório de truncamento. Sem ele, o usuário não sabe se leu o nome completo ou se há conteúdo além do visível.
**Consequência:** confusão de identidade entre credenciais com nomes similares (ex: "Gmail Pessoal" vs "Gmail Pessoal 2"); cópia da senha errada.

---

#### Pasta Virtual Indistinguível de Pasta Real *(Médio)*
**O que é:** a pasta virtual "Favoritos" aparece na árvore com o mesmo visual (`▶`) de pastas regulares — mesmo símbolo, mesmo estilo, mesmas ações na barra de comandos.
**Por que é errado:** o design system prevê `text.secondary + italic` para pastas virtuais/leitura. O requisito é explícito: pasta virtual é somente leitura. A restrição deve ser comunicada antes da tentativa de ação.
**Consequência:** usuário tenta criar segredo em Favoritos e recebe erro sem nenhuma pista visual prévia de que a operação era impossível.

---

#### Caractere de Largura Dupla Sem Contabilização de Colunas *(Alto)*
**O que é:** símbolos de largura ambígua ou caracteres que o terminal renderiza em 2 colunas (ex: alguns CJK, certos blocos Unicode) são usados sem compensar a coluna extra — o restante da linha fica deslocado para a direita em 1 coluna.
**Por que é errado:** o design system exige símbolos BMP de largura previsível de 1 coluna exatamente por isso. A largura de um caractere Unicode depende do locale, da fonte e do emulador de terminal — não da especificação Unicode. Código que assume 1 byte = 1 coluna ou 1 rune = 1 coluna produz desalinhamento silencioso, invisível em um terminal e quebrado em outro.
**Consequência:** bordas de modais não fecham, colunas de tabelas ficam tortas, separadores `│` aparecem fora de posição — tudo de forma irreproducível entre ambientes de desenvolvimento e produção.

---

#### Layout Não Recalculado Após Redimensionamento *(Crítico)*
**O que é:** ao redimensionar o terminal (janela menor ou maior), o conteúdo continua renderizado com as dimensões anteriores — textos transbordam para fora da área visível, separadores ficam curtos ou longos demais, modais saem da tela.
**Por que é errado:** o terminal pode ser redimensionado a qualquer momento. O evento `WindowSizeMsg` do Bubble Tea deve propagar as novas dimensões para todos os componentes que calculam largura ou altura. Qualquer componente que cacheia dimensões sem reagir ao resize produz layout quebrado.
**Consequência:** após redimensionar, a interface fica inutilizável até reiniciar a aplicação; bordas de modais ultrapassam os limites do terminal; conteúdo da árvore sobrescreve o painel de detalhe.

---

#### String Mais Longa que o Espaço Disponível Sem Truncamento *(Alto)*
**O que é:** um nome de segredo, valor de campo, mensagem de erro ou título de diálogo com comprimento superior ao espaço disponível é renderizado sem truncamento — o texto simplesmente transborda para a próxima linha ou para fora da área do componente.
**Por que é errado:** o design system define `…` como mecanismo obrigatório de truncamento para qualquer texto que exceda o espaço disponível. Transbordo de texto quebra o layout de grade fixa do terminal: uma linha que "vaza" desloca todas as linhas seguintes e pode corromper a estrutura visual inteira do componente.
**Consequência:** nome de segredo muito longo empurra o conector `<╡` para fora da coluna do separador; título de diálogo longo quebra a borda superior; valor de campo longo sobrescreve o label do campo abaixo.

---

#### Conteúdo Vertical Sem Scroll Quando Excede a Área *(Alto)*
**O que é:** um painel, diálogo ou seção com conteúdo que excede a altura disponível simplesmente corta o conteúdo no limite inferior — sem indicador de scroll, sem setas `↑↓`, sem thumb `■`.
**Por que é errado:** o design system define o padrão de scroll para exatamente este caso: indicadores na borda direita comunicam que há conteúdo além da viewport. Cortar silenciosamente é invisível — o usuário não sabe que existe mais conteúdo abaixo.
**Consequência:** campos do final de um segredo com muitos campos ficam inacessíveis; parte inferior de um diálogo longo — incluindo as ações de confirmação — desaparece; o usuário não consegue confirmar nem cancelar a operação.

---

#### Valor de Campo Mais Longo que a Área do Campo *(Médio)*
**O que é:** o valor de um campo de texto (especialmente campos sensíveis com senhas longas ou chaves de API) é exibido sem adaptação ao espaço disponível — sem truncamento, sem scroll horizontal, sem quebra de linha controlada.
**Por que é errado:** campos de detalhe têm largura fixa ditada pelo painel. Um valor mais longo que a área do campo transborda para fora dos limites do componente ou quebra a grade de alinhamento do restante do painel.
**Consequência:** senha longa sobrescreve o label do campo ao lado; chave de API com 128 caracteres empurra o separador `│` para fora da posição; o usuário não consegue ver nem copiar o valor completo sem revelar e fazer scroll.

---

#### Sangramento de Estilo ANSI *(Alto)*
**O que é:** uma sequência de estilo ANSI (cor, bold, dim) abre mas não é fechada com reset explícito — o estilo vaza para o conteúdo seguinte, colorindo ou engrossando texto que não deveria receber aquele atributo.
**Por que é errado:** no terminal, o estado de estilo é global e persistente até ser explicitamente resetado. Diferente de CSS (escopo por elemento), uma tag ANSI sem fechamento afeta tudo que vier depois — incluindo outros componentes, a barra de comandos e até o prompt do shell após a aplicação encerrar.
**Consequência:** item selecionado em `accent.primary` contamina o texto do painel de detalhe; erro em `semantic.error` deixa toda a barra de comandos vermelha; ao sair, o prompt do shell fica colorido com o último estilo aberto da TUI.

---

#### Cálculo de Largura com Sequências ANSI Incluídas *(Alto)*
**O que é:** ao calcular quantas colunas um texto ocupa (para truncar, alinhar ou construir bordas), o código usa o comprimento da string Go (`len(s)`) em vez da largura visual — que exclui os bytes das sequências de escape ANSI.
**Por que é errado:** uma string estilizada como `"\x1b[1mGmail\x1b[0m"` tem `len()` = 15 bytes, mas ocupa apenas 5 colunas visíveis. Usar `len()` produz alinhamento errado, truncamento precoce e bordas que parecem corretas no código mas ficam deslocadas na tela.
**Consequência:** o título de um modal com texto em bold fica desalinhado na borda superior; itens da árvore com indicadores coloridos (`✦`, `✎`) têm indentação incorreta; a barra de comandos parece ter espaço vazio sobrando ou faltando.

---

#### Layout Colapsando em Estado Vazio *(Médio)*
**O que é:** quando um painel fica sem conteúdo (cofre vazio, pasta vazia, resultado de busca sem correspondências), o componente colapsa para altura zero — fazendo o outro painel ou elemento adjacente expandir inesperadamente ou o separador `│` desaparecer.
**Por que é errado:** o painel deve manter sua proporção e estrutura independentemente do conteúdo. Estado vazio é um estado válido que precisa de altura mínima reservada e de uma mensagem de estado vazio — não de colapso silencioso de layout.
**Consequência:** ao esvaziar uma pasta, o painel esquerdo some e o painel de detalhe ocupa 100% da largura; o separador desaparece; ao criar um item, o layout salta de volta para a proporção 35/65 — violando estabilidade espacial.

---

#### Deslocamento de Conteúdo por Indicador Aparecendo ou Sumindo *(Médio)*
**O que é:** ao adicionar ou remover um indicador de estado junto ao nome de um item (ex: `✦` ao criar, `✎` ao editar, `★` ao favoritar), o nome do item se desloca horizontalmente — todos os outros itens sem indicador ficam numa coluna, os com indicador ficam numa coluna diferente.
**Por que é errado:** indica que o indicador não ocupa uma coluna reservada fixa. O layout correto reserva sempre 1 coluna para o indicador — exibindo espaço em branco quando não há indicador — de modo que o nome do item começa sempre na mesma coluna.
**Consequência:** ao favoritar um segredo, o nome "pula" 1 coluna para a direita; ao salvar e os indicadores `✦ ✎` desaparecem, todos os nomes "pulam" de volta — a árvore inteira trepida visualmente.

---

#### Largura de Spinner Variante Entre Frames *(Baixo)*
**O que é:** os frames do spinner `◐ ◓ ◑ ◒` têm larguras de renderização diferentes entre si em alguns terminais — causando leve deslocamento do texto adjacente a cada tick de animação.
**Por que é errado:** todos os quatro caracteres do inventário ocupam 1 coluna no BMP, mas terminais com fontes que tratam Geometric Shapes como "largura ambígua" podem renderizá-los em 2 colunas dependendo do locale. Se houver divergência entre frames, a mensagem adjacente ao spinner oscila 1 coluna por segundo.
**Consequência:** a mensagem `◐ Abrindo cofre…` pisca horizontalmente a cada frame do spinner; em terminais afetados, a mensagem parece trêmula mesmo sem nenhuma atualização de conteúdo.

---

#### Shift de Layout por Contador com Largura Variável *(Baixo)*
**O que é:** contagens numéricas exibidas na árvore (ex: número de segredos por pasta: `12`, `9`, `124`) não têm largura reservada fixada — quando o número muda de 1 para 2 dígitos ou de 2 para 3, o conteúdo adjacente se desloca.
**Por que é errado:** texto que muda de largura dinamicamente quebra o alinhamento de tudo que está à direita ou que usa a mesma coluna como referência. A largura máxima previsível deve ser reservada (ou o valor deve ser right-aligned num campo de largura fixa).
**Consequência:** ao salvar e um item ser removido, a contagem da pasta muda de `10` para `9` e o nome da pasta "pula" 1 coluna para a esquerda; ao importar vários segredos, as contagens saltam de 2 para 3 dígitos e toda a coluna de nomes se move.

---

#### Artefato Visual de Render Anterior *(Alto)*
**O que é:** ao fechar um modal, trocar de aba ou atualizar o conteúdo de um painel, parte do frame anterior permanece visível — caracteres de bordas, texto ou cores do estado anterior não foram sobrescritos pelo novo frame.
**Por que é errado:** o Bubble Tea renderiza por diff — apenas o que mudou é reescrito. Se o componente anterior ocupava mais linhas ou colunas que o novo, as células excedentes ficam com o conteúdo antigo. O componente precisa preencher explicitamente todo o espaço que reserva, incluindo linhas em branco ao final.
**Consequência:** ao fechar um diálogo de confirmação alto, os caracteres da borda inferior do diálogo ficam "flutuando" sobre o conteúdo de fundo; ao navegar de um segredo com 8 campos para um com 2 campos, os 6 campos extras do item anterior continuam visíveis.

---

#### Linha Final do Terminal Causando Scroll Involuntário *(Médio)*
**O que é:** um componente escreve na última coluna da última linha do terminal, fazendo o terminal interpretar isso como overflow e rolar o conteúdo uma linha para cima — deslocando toda a interface.
**Por que é errado:** em muitos emuladores de terminal, escrever na célula exata `(linhas, colunas)` — o canto inferior direito — aciona um scroll automático independente da intenção da aplicação. A barra de comandos, que ocupa a última linha, deve evitar escrever na última coluna.
**Consequência:** a barra de comandos "cai" uma linha, criando uma linha vazia no topo da interface; cada interação que re-renderiza a barra agrava o problema, empurrando o cabeçalho gradualmente para fora da tela.

---

#### Cursor Posicionado em Coluna Incorreta Durante Edição *(Alto)*
**O que é:** o cursor visual (bloco ou underline piscante) aparece numa coluna diferente de onde o próximo caractere será inserido — geralmente por causa de discrepância entre a largura em bytes usada para calcular a posição e a largura em colunas visuais do conteúdo já digitado.
**Por que é errado:** o terminal posiciona o cursor por colunas físicas. Se o código calcula `offset = len(input[:cursor])` em bytes em vez de somar as larguras de rune de cada caractere (especialmente com caracteres multibyte como acentos, `ü`, `ñ`, `ç` ou emojis), o cursor fica desalinhado em relação ao texto exibido — confundindo e dificultando a edição.
**Consequência:** ao editar uma senha com `ç` ou `á`, o cursor "pula" para a posição errada; o usuário vê o cursor sobre o caractere X mas ao digitar o novo caractere aparece na posição X+1 ou X-1; backspace apaga o caractere errado.

---

#### Campo de Edição Sem Scroll Horizontal *(Alto)*
**O que é:** ao editar um campo cujo valor é mais longo que a largura visual do campo, o texto simplesmente é truncado ou overflow — sem que o conteúdo role horizontalmente para manter o cursor sempre visível.
**Por que é errado:** o estado editável de um campo não é o mesmo que o estado de exibição. Um campo de exibição pode truncar com `…`, mas um campo em modo de edição deve exibir sempre a janela do texto centrada no cursor. Sem scroll, o usuário perde o contexto do que está digitando assim que o valor ultrapassa a largura do campo.
**Consequência:** ao editar uma URL de 80 caracteres num campo de 30 colunas, o texto após a coluna 30 desaparece; o usuário não consegue ver o que digitou nem confirmar se está correto; mover o cursor com `←`/`→` navega "no escuro" enquanto o conteúdo da tela fica estático.

---

### Tipografia e Cor

#### Bold Inflacionado *(Baixo)*
**O que é:** bold é aplicado a títulos de seção, labels de campos, nomes de pastas, nomes de segredos, ações na barra e texto de atalho — tudo em bold.
**Por que é errado:** bold é o único destaque tipográfico universalmente confiável — sua força vem da raridade. Se tudo é bold, nada é bold.
**Consequência:** hierarquia tipográfica colapsa; a interface parece gritante e dificulta a leitura de escaneamento rápido.

---

#### Token Semântico com Uso Decorativo *(Baixo)*
**O que é:** `semantic.success` (verde) é usado para destacar o nome do segredo favorito; `semantic.warning` é usado como cor de dica por "parecer uma atenção".
**Por que é errado:** o design system é explícito: *"`semantic.*` existe para comunicar estado operacional, nunca para ornamentar a interface."* Uso decorativo dilui o significado dos tokens.
**Consequência:** usuário para de confiar nos indicadores semânticos; ignora mensagens de erro genuínas por acreditar que são "decoração".

---

#### Cor Hardcoded *(Alto)*
**O que é:** o código usa literais hex diretamente (ex: `#7aa2f7`) em vez dos tokens semânticos (ex: `accent.primary`, `border.focused`).
**Por que é errado:** a abstração de tokens é a base do sistema de temas. Hardcode significa que o segundo tema continua exibindo as cores do primeiro.
**Consequência:** troca de tema requer alteração de código; o tema Cyberpunk nunca funciona corretamente onde há hardcode.

---

#### Italic Sem Reforço de Cor *(Baixo)*
**O que é:** hints e textos de dica aparecem *apenas* em italic — sem a cor `text.secondary` exigida como complemento pelo design system.
**Por que é errado:** italic tem suporte parcial. Em terminais sem suporte (frequente no Windows/ConHost), o texto de dica fica visualmente idêntico ao conteúdo primary. A combinação `italic + text.secondary` é obrigatória para garantir legibilidade cruzada.
**Consequência:** em terminais sem suporte a italic, hints são indistinguíveis do conteúdo real — o usuário não sabe o que é dado e o que é orientação.

---

### Acessibilidade

#### Estado Comunicado Apenas por Cor *(Alto)*
**O que é:** estados como "modificado", "criado" e "excluído" na árvore são comunicados apenas por cor — os símbolos `✦ ✎ ✗` não são usados.
**Por que é errado:** viola o princípio fundamental de acessibilidade: *"nenhum estado crítico pode depender exclusivamente de cor."* Toda a seção de fallback NO_COLOR da matriz de acessibilidade existe para garantir este comportamento.
**Consequência:** em ambientes CI, terminais monocromáticos ou com `$NO_COLOR`, toda informação de estado de sessão desaparece.

---

#### Área de Clique de Um Caractere *(Baixo)*
**O que é:** o conector `<╡` — indicador do item selecionado na árvore — é o único alvo clicável para "selecionar e ver detalhe", com 2 colunas de largura em posição variável.
**Por que é errado:** o design system prevê mouse como canal secundário completo. Alvo de 2 colunas em posição dinâmica exige precisão impossível para uso fluido. O clique deve funcionar na linha inteira.
**Consequência:** usuário com mouse não consegue usar a aplicação de forma eficiente; cada clique requer precisão cirúrgica.

---

#### Mensagem de Erro Técnica Exposta ao Usuário *(Médio)*
**O que é:** erros de sistema (arquivo corrompido, falha de I/O) exibem a mensagem interna do runtime, como `"unexpected end of JSON at offset 1247"` ou o caminho completo do arquivo.
**Por que é errado:** viola o requisito de privacidade (*"ausência total de logs contendo caminhos de arquivos"*) e o princípio de que mensagens de erro são para o usuário, não para o sistema.
**Consequência:** exposição do caminho do arquivo (localização do cofre); usuário sem ação clara diante de mensagem técnica incompreensível.

---

### Ciclo de Vida do Cofre

#### Auto-save Silencioso *(Alto)*
**O que é:** a aplicação salva automaticamente após certas operações (ex: alteração de senha mestra) sem comunicar claramente que isso ocorreu nem distinguir visualmente do salvamento manual.
**Por que é errado:** a alteração de senha mestra é a *única* exceção explícita ao princípio *Controle total do usuário* — e por isso exige comunicação extra-clara, não mais silêncio.
**Consequência:** usuário tenta desfazer a alteração de senha achando que é reversível como tudo mais no cofre — e descobre que não é.

---

#### Conflito de Arquivo Ignorado ou Minimizado *(Crítico)*
**O que é:** ao detectar que o arquivo foi modificado externamente, a aplicação salva por cima sem aviso — ou exibe um diálogo tão breve que o usuário confirma reflexivamente.
**Por que é errado:** o requisito exige escolha explícita entre Sobrescrever / Salvar como novo / Cancelar. Trata-se de diálogo de severidade destrutiva — perda de dados de outra sessão é irreversível.
**Consequência:** dados salvos por outra instância ou backup restaurado destruídos silenciosamente.

---

#### Re-autenticação Durante a Sessão *(Alto)*
**O que é:** a aplicação solicita a senha mestra novamente para operações de salvamento ou exportação durante uma sessão já autenticada.
**Por que é errado:** o requisito é explícito: *"a senha é fornecida uma única vez ao abrir o cofre. Não há re-solicitação de senha para salvamento ou descarte."* Re-autenticação frequente é *security theater* — a senha já está em memória — e treina o usuário a digitar sua senha mestra em contextos inesperados.
**Consequência:** fricção sem ganho de segurança; usuário condicionado a digitar a senha sem questionar o contexto (vetor de phishing de UI).

---

#### Exclusão Que Desaparece Imediatamente *(Crítico)*
**O que é:** itens marcados para exclusão somem da lista imediatamente — sem `✗` + strikethrough — como se já tivessem sido deletados permanentemente.
**Por que é errado:** o requisito define exclusão como *marcação reversível até o salvamento*. O design system define tratamento visual específico: `semantic.warning + ✗ + strikethrough`. Sumir com o item remove a reversibilidade de fato, mesmo que tecnicamente ela ainda exista.
**Consequência:** usuário crê ter deletado permanentemente; não tenta desfazer; perde dados que poderiam ser recuperados.

---

#### Importação Sem Prévia de Impacto *(Crítico)*
**O que é:** ao importar um arquivo de intercâmbio, a aplicação executa a mesclagem diretamente — sobrescrevendo segredos conflitantes e criando pastas — sem mostrar ao usuário o que será alterado antes da confirmação.
**Por que é errado:** a política de importação (mescla de pastas, sobrescrita de segredos e modelos conflitantes) é agressiva e irreversível no arquivo. O impacto precisa ser comunicado *antes* da confirmação, não descoberto *depois*.
**Consequência:** perda de versões de segredos que o usuário não pretendia sobrescrever; auditoria manual do cofre inteiro necessária para entender o que mudou.


