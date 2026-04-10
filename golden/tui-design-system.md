# Design System — Abditum TUI

> Fundações visuais e padrões transversais para o pacote `internal/tui`.
> Define princípios, tokens, estados e padrões que governam toda decisão de UI.

## Escopo deste documento

Este documento define **fundações** e **padrões reutilizáveis** — o que cada peça visual é, como se comporta em abstrato e como peças se combinam em situações recorrentes.

A composição dessas peças em telas, wireframes e fluxos concretos pertence ao documento de especificação:
- [`tui-specification-novo.md`](tui-specification-novo.md) — telas, wireframes de componentes e fluxos visuais

> **Regra de governança:** toda decisão de UI/UX deve ser compatível com os princípios definidos aqui. Em conflito entre especificação local e princípio, este último prevalece.

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
  - [Mensagens](#mensagens)
  - [Foco e Navegação](#foco-e-navegação)
  - [Mapa de Teclas](#mapa-de-teclas)
  - [Acessibilidade](#acessibilidade)
- [Diálogos](#diálogos)
  - [Especificação Geral](#especificação-geral)
  - [Diálogos de Decisão](#diálogos-de-decisão)
  - [Help](#help)
  - [Diálogos Funcionais](#diálogos-funcionais)

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
| | `border.focused` | Borda do painel ativo, de janelas de entrada (senhas, textos) e de diálogos com severidade neutra. Diálogos com severidade não-neutra usam `semantic.*` — ver [Sobreposição](#sobreposição) | `#7aa2f7` <span style="color:#7aa2f7\">██</span> | `#ff2975` <span style="color:#ff2975\">██</span> |
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
| **Área de trabalho** | Restante | Conteúdo do modo ativo (cofre, modelos, configurações, bem-vindo) |
| **Barra de mensagens** | 1 linha | Borda separadora `─` com mensagem embutida — quando há mensagem, o texto substitui o trecho central da borda |
| **Barra de comandos** | 1 linha | Ações do contexto ativo |

### Proporções de painel

Para modos com dois painéis (Cofre, Modelos):

| Painel | Proporção | Papel |
|---|---|---|
| Esquerdo (árvore / lista) | ~35% | Navegação e seleção |
| Direito (detalhe) | ~65% | Conteúdo do item selecionado |

A proporção é aproximada — a implementação pode ajustar em ±5% para alinhamento estético ou para acomodar terminais muito largos.

### Barra de mensagens

| Parâmetro | Valor |
|---|---|
| Altura | 1 linha fixa entre área de trabalho e barra de comandos |
| Anatomia | Borda `─` contínua; mensagem embutida após 2 espaços de padding esquerdo |
| Largura da borda | 100% da largura do terminal |
| Largura do texto | Largura do terminal − 2 (padding) − 2 (margem `─` direita mínima) |
| Truncamento | Com `…` quando o texto exceder o espaço disponível |
| Sem mensagem | Borda `─` contínua (separador visual permanente) |

---

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

- Foco de painel é indicado pelo separador vertical (ver [Foco e Navegação](#foco-e-navegação)) e pela barra de comandos — não por borda ao redor do painel.
- TUI não tem estado "pressionado"; confirmação vem por mudança de contexto ou mensagem.
- Transições são instantâneas. A única animação prevista é o spinner `MsgBusy`.

---

## Padrões

Padrões são regras de comportamento transversais — aplicam-se a múltiplas telas e componentes. Os documentos de especificação consomem estes padrões ao definir componentes e fluxos concretos.

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


### Redação de Mensagens

Este guia estabelece estilo e gramática para todas as comunicações textuais na interface, garantindo clareza, concisão e consistência.

#### Princípios Gerais

-   **Direta e Objetiva:** Vá direto ao ponto. Evite rodeios, jargões desnecessários e linguagem floreada.
-   **Clara e Unívoca:** A mensagem deve ser compreendida de imediato, sem ambiguidade.
-   **Acionável (quando aplicável):** Em caso de erro ou alerta, sugira um próximo passo ou aponte a causa.
-   **Contextual:** Adapte a mensagem ao estado da interface e ao conhecimento do usuário naquele momento do fluxo.
-   **Minimalista:** Respeite o espaço limitado do terminal.

#### Tom de Voz

-   **Formal-neutro:** Use uma voz técnica, mas acessível. Evite personificação, gírias ou excesso de exclamações.
-   **Foco no Usuário:** Use a segunda pessoa ("você" implícito ou explícito quando necessário) para direcionar dicas e ações. Ex: "Digite a senha para desbloquear"
-   **Afirmativo:** Prefira frases afirmativas.

#### Gramática e Estilo

-   **Capitalização:**
    -   **Início de frase:** Sempre maiúscula.
    -   **Nomes de itens:** Conforme o nome original (sem capitalização artificial).
    -   **Labels de campo/ação:** Conforme a UI (ex: "Salvar", "Nova senha").
-   **Pontuação:**
    -   **Mensagens curtas (barra):** Sem pontuação final (ponto, exclamação). Ex: `✓ Cofre salvo`
    -   **Mensagens longas (diálogos):** Use ponto final para encerrar frases completas.
    -   **Perguntas (diálogos):** Use ponto de interrogação.
-   **Nomes de itens em mensagens:** Se referenciar um item específico (ex: "Gmail"), use aspas simples `'Gmail'` para distingui-lo do texto da mensagem, ou `**bold**` se o contexto for de realce crítico no diálogo.

#### Estrutura por Tipo de Mensagem

##### 1. Títulos de Diálogo (ex: na borda superior)

-   **Padrão:** O título deve ser o nome do fluxo ou da ação principal.
-   **Formato:** `[Nome do Fluxo/Ação Principal]` (capitalizado conforme o nome, ex: "Sair do Abditum", "Definir senha mestra", "Abrir cofre").

##### 2. Mensagens no Corpo do Diálogo

-   **Padrão:**
    -   **Diálogos de Decisão (Confirmação):** Afirmação de um fato (opcional), seguida de uma pergunta concisa que apresenta as opções de decisão. A pergunta não menciona a opção `Voltar` (Esc).
    -   **Diálogos de Decisão (Reconhecimento):** Apenas uma afirmação. Não há pergunta.
-   **Formato:**
    -   **Confirmação:** Fato termina com ponto; pergunta com interrogação.
    -   **Reconhecimento:** Afirmação termina com ponto final.
-   **Exemplos (ATUALIZADOS PARA CONCISÃO MÁXIMA):**
    -   `Sair do Abditum?`
    -   `Cofre modificado. Salvar ou descartar?`
    -   `Arquivo modificado externamente. Sobrescrever?`
    -   `'Gmail' será excluído permanentemente. Continuar?`
    -   `Arquivo corrompido ou inválido. Necessário fechar.`

##### 3. Mensagens na Barra de Mensagens (inferior)

-   **Padrão:** Curto e reativo, `<símbolo> [texto]`.
-   **Formato:** Começa com maiúscula (após o símbolo), sem pontuação final.
-   **Exemplos:**
    -   `✓ Cofre salvo`
    -   `ℹ Arquivo já existe`
    -   `⚠ Senha fraca`
    -   `✕ Senha incorreta`
    -   `◐ Salvando cofre`
    -   `• Digite a senha para desbloquear`

---

### Foco e Navegação


O modelo de foco define como o usuário percebe e alterna entre áreas interativas da interface.

**Alternância com Tab:**

`Tab` é contextual — o comportamento depende do estado do painel:

| Contexto | `Tab` | `⇧Tab` |
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

- Campos não possuem borda — a área digitável é delimitada por um fundo `surface.input` (tom rebaixado em relação ao `surface.raised` do diálogo).
- Label do campo ativo em `accent.primary` + **bold**; labels dos campos inativos em `text.secondary`.
- Foco indicado pela presença do cursor `▌` em `text.primary` dentro do fundo `surface.input`.
- Placeholder em `text.secondary` + *italic* — desaparece ao digitar.
- Erro de validação: exibido na barra de mensagens (tipo Erro), não inline — formulários simples mostram um erro por vez.
- Em **NO_COLOR**: o fundo `surface.input` pode ser perdido; o cursor + label em **bold** permanecem como indicadores de foco suficientes.

---

### Mapa de Teclas

Esta seção define a **política de atribuição de teclas** — como atalhos são organizados, quais regras regem conflitos e quais teclas têm significado global. O mapeamento completo por tela está na [especificação de telas](tui-specification-novo.md).

### Representação Visual de Teclas e Modificadores

Para garantir consistência e clareza na documentação de atalhos, são adotadas as seguintes representações visuais compactas para teclas e modificadores:

| Tecla / Modificador | Representação Visual | Unicode | Notas |
|---|---|---|---|
| `Ctrl` | `⌃` | U+2303 (UP ARROWHEAD) | Usado para atalhos de controle. |
| `Shift` | `⇧` | U+21E7 (UPWARDS WHITE ARROW) | Usado para atalhos de modificação ou navegação. |
| `Alt` | `!` | (Caracter comum) | Usado para atalhos alternativos. |
| `Del` | `Del` | (Texto simples) | Tecla Delete. |
| `Ins` | `Ins` | (Texto simples) | Tecla Insert. |
| `PgUp` | `PgUp` | (Texto simples) | Page Up. |
| `PgDn` | `PgDn` | (Texto simples) | Page Down. |
| `Home` | `Home` | (Texto simples) | Início da linha/conteúdo. |
| `End` | `End` | (Texto simples) | Fim da linha/conteúdo. |
| `Esc` | `Esc` | (Texto simples) | Abandona ou retrocede. |
| `Enter` | `Enter` | (Texto simples) | Confirma ou avança. |
| `Tab` | `Tab` | (Texto simples) | Alterna foco ou campos. |

### Política de escopos e Ergonomia:

As teclas são atribuídas seguindo uma hierarquia de escopos e agrupamentos físicos no teclado, visando otimizar a memória muscular e evitar acionamentos acidentais para ações críticas. Escopos mais específicos sobrepõem os mais gerais quando ambos estão ativos.

| Escopo | Descrição | Exemplo |
|---|---|---|
| **Global** | Funciona em qualquer contexto da aplicação | `F1`, `F12`, `⌃Q`, `⌃!⇧Q` |
| **Área de trabalho** | Funciona quando a área de trabalho tem foco (sem diálogo aberto) | `F2-F11` (ações do cofre e modos), `⇧F6`, `⇧F7`, `⌃F7` |
| **Diálogo** | Funciona apenas enquanto um diálogo está no topo da pilha | `Enter`, `Esc`, `Tab` |
| **Contextual/Foco** | Ações específicas do item ou campo com foco | `Ins`, `Del`, `⌃<letra>` (para ações locais) |

### Regras de Consistência e Semântica de Modificadores

As representações visuais de teclas e modificadores seguem as definições da seção [Representação Visual de Teclas e Modificadores](#representação-visual-de-teclas-e-modificadores). As regras semânticas de uso são:

-   `Enter` sempre avança ou aprofunda: confirma em diálogos, seleciona/expande na árvore, ativa/confirma edição de campo.
-   `Esc` sempre retrocede ou abandona: fecha modal, cancela edição, sai de modo (busca, edição).
-   `Tab` / `⇧Tab` navegam entre painéis (modo leitura) ou campos (modo edição).
-   `↑↓←→` são para navegação direcional em listas, árvores e campos.
-   `Home` / `End` navegam ao primeiro / último item visível ou início/fim de linha em campos.
-   `PgUp` / `PgDn` realizam scroll por página (viewport − 1) em conteúdo com scroll.
-   `Ins`: Sugerido para ações de inserção/criação (no contexto do foco).
-   `Del`: Sugerido para ações de exclusão (no contexto do foco).
-   Se uma tecla precisa ter significado diferente em dois contextos, isso deve ser documentado e justificado na especificação.
-   Teclas de navegação universais (`↑↓←→`, `Tab`, `Home`, `End`, `PgUp`, `PgDn`) não aparecem na barra de comandos — são senso comum em TUI. Exceção: diálogos podem exibir opções explicitamente.

**Atalhos Globais (Aplicam-se em qualquer contexto):**

| Tecla | Ação (Função) | Notas |
|---|---|---|
| `F1` | Abrir / fechar modal de Ajuda | |
| `F12` | Alternar Tema | Ação pontual, sem necessidade de visibilidade permanente na barra de comandos |
| `⌃Q` | Sair da Aplicação | Gerencia todas as saídas com as devidas confirmações |
| `⌃!⇧Q` | Bloquear Cofre | Bloqueio emergencial, descarta alterações, sem confirmação. Atalho "complicado" para evitar acidentes. |

**Teclas de Área de Trabalho (Ativas quando a área de trabalho tem foco, sem diálogos):**

A atribuição específica de teclas a fluxos individuais é detalhada na [especificação de telas](tui-specification-novo.md), mas as teclas F são reservadas por grupos de ações, seguindo a ergonomia do teclado físico:

-   **`F2` a `F4`**: Reservadas para **seleção das áreas de trabalho** (Modo Cofre, Modelos, Configurações).
-   **`F5` a `F8`**: Reservadas para **ações de persistência do cofre** (criar, abrir, salvar, recarregar).
-   **`F9` a `F11`**: Reservadas para **ações complementares de gerenciamento do cofre** (exportar, importar, alterar senha mestra).

> **Fluxo 7 — Aviso de Bloqueio Iminente por Inatividade:** É um fluxo iniciado pelo sistema, não requer um atalho manual do usuário.


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

## Diálogos

Os diálogos são elementos de sobreposição que interrompem o fluxo principal para apresentar informação ou solicitar entrada do usuário. Existem três categorias:

- **Decisão** (reconhecimento e confirmação) — exibem uma informação e solicitam uma escolha
- **Help** — conteúdo de ajuda estruturado; reconhecimento especializado com scroll
- **Funcional** — diálogos interativos com campos de entrada (ex: senha, caminho de arquivo)

---

### Especificação Geral

Regras comuns a todos os tipos de diálogo.

**Apresentação:**

- Estilo de borda: Rounded (`╭╮╰╯│─`)
- Fundo interno: `surface.raised`
- O conteúdo abaixo permanece visível, mas inativo — sem escurecimento de overlay
- Posição: centrado horizontal e verticalmente na área de trabalho

| Parâmetro | Valor | Justificativa |
|---|---|---|
| Largura mínima | O necessário para caber título na borda superior e ações na borda inferior | Garante integridade da moldura |
| Largura máxima | 95% do terminal | Evita que o diálogo pareça "colado" nas bordas |
| Padding interno horizontal | 2 colunas | Respiro visual sem desperdício |
| Padding interno vertical | 1 linha | Respiro entre borda e conteúdo |

```
╭── Título ────────────────────────────────────────────────╮
│                                                          │
│  Conteúdo do diálogo (linha visível 1)                   ■
│  Conteúdo do diálogo (linha visível 2)                   │
│  Conteúdo do diálogo (linha visível 3)                   ↓
╰── Enter Confirmar ─────────────────────── Esc Cancelar ──╯
```

**Moldura:**

- **Borda superior** contém o título embutido, precedido pelo símbolo de severidade quando aplicável. Diálogos sem severidade (Help, Funcional) não exibem símbolo
- **Borda inferior** contém as ações de confirmação e cancelamento — nunca duplicadas na barra de comandos
- **Ação default** (associada a `Enter`): tecla + label em **bold**, coloridos com o token de destaque do tipo
- **Demais ações**: tecla + label na cor da borda, sem bold
- **Borda e título** usam o mesmo token de cor — definido pelo tipo do diálogo
- **Ações internas** (revelar senha, alternar campo, expandir diretório) aparecem exclusivamente na barra de comandos — nunca na borda
- **Teclas de navegação** (↑↓, →, ←, Tab) não aparecem na borda nem na barra de comandos — são senso comum em TUI

**Navegação (todos os tipos):**

| Tecla | Comportamento |
|---|---|
| `Enter` | Aciona a ação default |
| `Esc` | Aciona a ação de cancelamento; se não existir, fecha o diálogo |
| Atalho da opção | Aciona diretamente a opção correspondente |

**Sobreposição:**

- Apenas o diálogo do topo recebe input; os inferiores permanecem montados, porém congelados
- A barra de comandos reflete os atalhos do diálogo do topo enquanto ele estiver ativo
- Ao fechar o diálogo, o foco retorna ao elemento anterior na pilha (ou ao conteúdo base)

**Scroll:**

Quando há conteúdo fora da viewport, a borda direita comunica direção e posição usando três elementos:

- **Setas** (`↑` / `↓`) — substituem o primeiro e/ou último `│` da borda direita
- **Thumb** (`■`) — indica a posição relativa no conteúdo. Posição calculada: `round(scroll_offset / max_scroll × (linhas_borda - 1))`

| Posição do scroll | Borda direita (primeiro `│`) | Borda direita (último `│`) | Thumb `■` |
|---|---|---|---|
| Topo (mais conteúdo abaixo) | `│` (normal) | `↓` | Próximo ao topo |
| Meio (conteúdo acima e abaixo) | `↑` | `↓` | Proporcional à posição |
| Final (mais conteúdo acima) | `↑` | `│` (normal) | Próximo à base |
| Sem scroll (tudo visível) | `│` (normal) | `│` (normal) | Não aparece |

> **Prioridade de renderização:** se o thumb `■` coincide com a posição de uma seta, a seta prevalece — direção é mais importante que posição.

| Elemento | Token | Atributo |
|---|---|---|
| Seta de scroll (`↑` / `↓`) | `text.secondary` | — |
| Thumb de posição (`■`) | `text.secondary` | — |

| Tecla | Efeito |
|---|---|
| `↑` / `↓` | Move uma linha |
| `PgUp` / `PgDn` | Move uma página (viewport − 1 linhas) |
| `Home` / `End` | Vai ao início / fim do conteúdo |

- **Scroll do mouse** roda o conteúdo com foco
- **Clique na seta** (`↑`/`↓`) move uma linha
- **Drag do thumb** (`■`) não suportado — TUI não tem drag contínuo

> As bordas superior e inferior (título, ações) permanecem intactas — o scroll afeta apenas o conteúdo interno.

---

### Diálogos de Decisão

Exibem uma informação e solicitam uma escolha do usuário. Dois subtipos:

- **Reconhecimento** — o usuário toma ciência de uma informação (1 ação)
- **Confirmação** — o usuário decide entre continuar ou cancelar (2–3 ações)

**Severidade** — governa borda, símbolo e cor da ação default:

| Severidade | Símbolo | Token de borda | Token da ação default | Quando usar |
|---|---|---|---|---|
| Destrutivo | `⚠` | `semantic.warning` | `semantic.error` | Ação irreversível ou com perda de dados |
| Erro | `✕` | `semantic.error` | `accent.primary` | Falha ocorrida, condição irrecuperável |
| Alerta | `⚠` | `semantic.warning` | `accent.primary` | Situação importante mas recuperável |
| Informativo | `ℹ` | `semantic.info` | `accent.primary` | Informação que requer atenção |
| Neutro | — | `border.focused` | `accent.primary` | Operação rotineira, sem urgência |

> Destrutivo e Alerta compartilham símbolo (`⚠`) e borda (`semantic.warning`). A distinção está na ação default: `semantic.error` para destrutivo, `accent.primary` para alerta — o perigo está na *ação*, não na *situação*.

**Aparência:**

| Elemento | Token | Atributo |
|---|---|---|
| Símbolo na borda superior | token de borda da severidade | — |
| Título na borda superior | token de borda da severidade | **bold** |
| Mensagem no corpo | `text.primary` | — |
| Nomes referenciados no corpo | `text.primary` | **bold** |
| Ação default | token da ação default da severidade | **bold** |
| Demais ações | token de borda da severidade | — |

> Barra de mensagens **não utilizada** — a mensagem completa está no corpo. Barra de comandos **vazia** — ações estão na borda inferior.

**Ações na borda inferior:**

| Ações | Layout | Uso |
|---|---|---|
| **1 ação** | Alinhada à **direita** | Reconhecimento |
| **2 ações** | Default à **esquerda**, Cancelar à **direita** | Confirmação binária |
| **3 ações** | Default à **esquerda**, alternativa no meio, Cancelar à **direita** | Confirmação com alternativa |

> **Limite:** 3 ações é o máximo tolerado. Mais que isso indica falha de design — divida em etapas ou use seletor interno.

Wireframes:

```
╰── Enter OK ──────────────────────────────────╯
```
```
╰── Enter Excluir ───────────────────── Esc Cancelar ──╯
```
```
╰── Enter Salvar ── A Salvar como ─────── Esc Cancelar ──╯
```

**Scroll em diálogos de decisão:**

Quando a mensagem excede o espaço disponível, o dimensionamento segue em sequência:

1. **Largura cresce** até o máximo (70 colunas ou 80% do terminal, o menor)
2. **Word-wrap** — quebra em linhas respeitando limites de palavra
3. **Altura cresce** até o máximo (80% da altura da área de trabalho)
4. Se ainda exceder: **scroll vertical** (indicadores e navegação conforme Especificação Geral)

O diálogo é sempre aberto com a primeira linha visível.

---

### Help

Diálogo de reconhecimento especializado que exibe conteúdo de ajuda estruturado e tipicamente longo.

**Aparência:**

| Elemento | Token | Atributo |
|---|---|---|
| Título | `border.default` | **bold** |
| Borda | `border.default` | — |
| Conteúdo | `text.primary` | — |
| Ação `Fechar` | `border.default` | **bold** |

> Sem símbolo de severidade. Barra de comandos **vazia**.

**Ações:**

```
╰── Enter Fechar ──────────────────────────────╯
```

**Eventos:**

| Evento | Efeito |
|---|---|
| `Enter` | Fecha o diálogo |
| `↑` / `↓` / `PgUp` / `PgDn` / `Home` / `End` | Navegação do conteúdo — conforme Especificação Geral |

---

### Diálogos Funcionais

Diálogos interativos com campos de entrada que oferecem uma função específica (ex: PasswordEntry, PasswordCreate, FilePicker).

**Aparência:**

| Elemento | Token | Atributo |
|---|---|---|
| Título | `text.primary` | **bold** |
| Borda (entrada de texto) | `border.focused` | — |
| Borda (consulta) | `border.default` | — |
| Label do campo ativo | `accent.primary` | **bold** |
| Label do campo inativo | `text.secondary` | — |
| Ação default (ativa) | `accent.primary` | **bold** |
| Ação default (bloqueada) | `text.disabled` | — |
| Ação Cancelar | token de borda | — |

> Sem símbolo de severidade. Barra de comandos exibe ações internas (ex: revelar senha, alternar campo).

**Estados:**

A ação default fica **inativa** até as condições mínimas serem satisfeitas — para todos os meios de acionamento (`Enter`, mouse, atalho):

| Componente | Estado | Condição |
|---|---|---|
| Ação default | ativa (`accent.primary` **bold**) | Condições mínimas satisfeitas |
| Ação default | bloqueada (`text.disabled`) | Condições mínimas **não** satisfeitas |
| Ação Cancelar | sempre ativa | — |

> A especificação de cada diálogo funcional documenta suas condições de desbloqueio.

**Barra de mensagens:**

| Momento | Conteúdo | Tipo |
|---|---|---|
| Diálogo abre | Dica contextual do primeiro campo com foco | Dica de campo (`•` *italic*) |
| Foco entra em campo (vazio ou válido) | Dica descritiva sobre o campo | Dica de campo (`•` *italic*) |
| Foco entra em campo (com valor inválido) | Mensagem de erro explicando a invalidação | Erro (`✕` **bold**, TTL 5s) |
| Tentativa de confirmar com validação falha | Mensagem de erro; diálogo permanece aberto | Erro (`✕` **bold**, TTL 5s) |
| Diálogo fecha | Barra é limpa | — |

> Mensagens pós-fechamento (ex: "◐ Criando cofre…", "✓ Cofre aberto", "Operação cancelada") são responsabilidade do orquestrador — não do diálogo.