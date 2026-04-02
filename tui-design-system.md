# Design System — Abditum TUI

> Fundações visuais e padrões transversais para o pacote `internal/tui`.
> Define princípios, tokens, estados e padrões que governam toda decisão de UI.
>
> **Regra de fronteira:** este documento define *fundações* — o que cada peça visual é e como se comporta em abstrato.
> A composição dessas peças em telas, wireframes e fluxos concretos pertence aos documentos de especificação.
>
> **Documentos complementares:**
> - [`tui-specification.md`](tui-specification.md) — composição de telas, wireframes e fluxos visuais
> - [`tui-elm-architecture.md`](tui-elm-architecture.md) — arquitetura de componentes (Elm pattern)

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

> **Regra de governança:** toda decisão de UI/UX deste projeto deve ser compatível com estes princípios. Em caso de conflito entre uma especificação local e um princípio, o princípio prevalece e a especificação deve ser ajustada.

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
- **Consistência de interação:** a mesma tecla, o mesmo símbolo e o mesmo tratamento visual mantêm o mesmo significado em toda a aplicação. `Enter` sempre confirma. `Esc` sempre cancela ou fecha. `F16` sempre revela/oculta. Exceções devem ser documentadas e justificadas.
- **Estabilidade espacial:** cabeçalho, árvore, detalhe e barra de comandos permanecem em posições previsíveis entre estados. O layout não "pula" quando o conteúdo muda. Isso preserva memória muscular e reduz carga cognitiva.

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
| | `border.focused` | Borda do painel ativo, de janelas de entrada (senhas, textos) e de diálogos neutros (perguntas sem urgência). Diálogos com severidade usam `semantic.*` — ver Sobreposição | `#7aa2f7` <span style="color:#7aa2f7">██</span> | `#ff2975` <span style="color:#ff2975">██</span> |
| **Interação** | `accent.primary` | Barra de seleção na lista, cursor de navegação, botão principal de ação | `#7aa2f7` <span style="color:#7aa2f7">██</span> | `#ff2975` <span style="color:#ff2975">██</span> |
| | `accent.secondary` | Ícone de favorito (★), nomes de pastas na navegação de arquivos | `#bb9af7` <span style="color:#bb9af7">██</span> | `#00fff5` <span style="color:#00fff5">██</span> |
| **Semânticas** | `semantic.success` | Operação concluída com sucesso, configuração ligada (ON) | `#9ece6a` <span style="color:#9ece6a">██</span> | `#05ffa1` <span style="color:#05ffa1">██</span> |
| | `semantic.warning` | Alerta antes de ação permanente, aviso de bloqueio por tentativas erradas | `#e0af68` <span style="color:#e0af68">██</span> | `#ffe900` <span style="color:#ffe900">██</span> |
| | `semantic.error` | Erro de operação, senha incorreta, borda de diálogos destrutivos | `#f7768e` <span style="color:#f7768e">██</span> | `#ff3860` <span style="color:#ff3860">██</span> |
| | `semantic.info` | Informação contextual, marcadores de segredos novos (`+`) ou modificados (`~`) | `#7dcfff` <span style="color:#7dcfff">██</span> | `#00b4d8` <span style="color:#00b4d8">██</span> |
| | `semantic.off` | Configuração desligada (OFF) | `#737aa2` <span style="color:#737aa2">██</span> | `#9999cc` <span style="color:#9999cc">██</span> |
| **Especiais** | `special.muted` | Texto de segredos marcados para exclusão (riscado, esmaecido) | `#565f89` <span style="color:#565f89">██</span> | `#666688` <span style="color:#666688">██</span> |
| | `special.highlight` | Fundo colorido atrás do item selecionado na lista | `#283457` <span style="background:#283457;color:#a9b1d6">██</span> | `#2a1533` <span style="background:#2a1533;color:#e0e0ff">██</span> |
| | `special.match` | Trecho de texto que corresponde ao termo digitado na busca | `#f7c67a` <span style="color:#f7c67a">██</span> | `#ffc107` <span style="color:#ffc107">██</span> |

### Notas de contraste

> **`special.muted` sobre `special.highlight`:** em Tokyo Night, `#565f89` (texto de item excluído) sobre `#283457` (fundo selecionado) apresenta contraste reduzido (~2.5:1). O símbolo `✕` + strikethrough garantem legibilidade na ausência de contraste de cor suficiente — conforme o princípio de dupla camada. Validar em monitores com brilho reduzido.

> **Aliases de valor:** `text.link` = `accent.primary` em hex. O alias documenta intenção — autores de temas podem divergir os valores quando precisarem distinguir link de ação primária.

> **Bordas de modais semânticos:** modais com urgência semântica (DialogAlert, DialogInfo) usam diretamente os tokens `semantic.warning`, `semantic.info` ou `semantic.error` como cor de borda — não existe token `border.*` separado para casos semânticos. A semântica do modal governa a borda.

### Gradiente do logo

| Linha | Tokyo Night | Cyberpunk |
|---|---|---|
| 1 | `#9d7cd8` <span style="color:#9d7cd8">██</span> | `#ff2975` <span style="color:#ff2975">██</span> |
| 2 | `#89ddff` <span style="color:#89ddff">██</span> | `#b026ff` <span style="color:#b026ff">██</span> |
| 3 | `#7aa2f7` <span style="color:#7aa2f7">██</span> | `#00fff5` <span style="color:#00fff5">██</span> |
| 4 | `#7dcfff` <span style="color:#7dcfff">██</span> | `#05ffa1` <span style="color:#05ffa1">██</span> |
| 5 | `#bb9af7` <span style="color:#bb9af7">██</span> | `#ff2975` <span style="color:#ff2975">██</span> |

### Decisão

> **Ambos os temas são suportados simultaneamente.** O usuário seleciona o tema ativo nas Configurações; F12 alterna rapidamente entre os dois sem abrir um menu.

## Tipografia

Em TUI não existem fontes nem tamanhos; a tipografia disponível é o conjunto de atributos ANSI que o terminal realmente suporta. O papel do design system é definir **quando** usar esses atributos e como degradar quando eles falharem.

### Atributos e fallback

| Atributo | Suporte | Fallback | Uso principal |
|---|---|---|---|
| **Bold** | Universal | — | Títulos, cursor selecionado, ação default |
| Dim / Faint | Amplo | Cor já comunica o estado | Itens desabilitados, conteúdo secundário |
| *Italic* | Parcial | `text.secondary` já diferencia | Hints, pastas virtuais, textos auxiliares |
| Underline | Amplo | — | Uso pontual |
| ~~Strikethrough~~ | Parcial | `✕` + `special.muted` preservam o sentido | Itens marcados para exclusão |
| Blink | Inconsistente | Não usar | Nenhum |

### Combinações previstas

| Combinação | Uso |
|---|---|
| Bold + cor semântica | Título de modal de alerta ou informação |
| Dim + strikethrough | Item excluído, com `✕` como reforço |
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

## Ícones e Símbolos

Inventário completo dos caracteres Unicode usados pela interface. A seleção privilegia símbolos do BMP (Basic Multilingual Plane), com largura previsível, em vez de glifos dependentes de fonte.

O contexto de uso detalhado de cada símbolo está na seção onde ele é consumido (Sobreposição, Mensagens, Estados Visuais, especificação de telas).

| Símbolo | Nome semântico | Colunas |
|---|---|---|
| `▶` | Pasta recolhida | 1 |
| `▼` | Pasta expandida | 1 |
| `▷` | Pasta vazia | 1 |
| `●` | Item folha | 1 |
| `★` | Favorito | 1 |
| `✕` | Marcado para exclusão | 1 |
| `+` | Adicionado na sessão | 1 |
| `~` | Modificado na sessão | 1 |
| `•` | Indicador contextual (ver nota) | 1 |
| `◉` | Campo revelável | 1 |
| `✓` | Sucesso | 1 |
| `ℹ` | Informação | 1 |
| `⚠` | Alerta / aviso | 1 |
| `✗` | Erro | 1 |
| `?` | Pergunta / decisão | 1 |
| `◐ ◓ ◑ ◒` | Spinner de atividade | 1 |
| `💡` | Dica de uso | 2 |
| `▌` | Cursor de campo | 1 |
| `↑` `↓` | Indicação de scroll | 1 |
| `─` `│` | Separadores | 1 |
| `╭╮╰╯` | Cantos arredondados (diálogos) | 1 |
| `<╡` | Conector árvore → detalhe | 1+1 |
| `…` | Truncamento | 1 |
| `••••` | Máscara de conteúdo sensível | 1/cada |

> **`•` reutilizado:** aparece como indicador de alterações pendentes no cabeçalho, como marcador de dica contextual na barra de mensagens, e como caractere de máscara em campos sensíveis. A distinção é sempre pelo contexto visual — nunca coexistem na mesma região.

> **`💡` é o único emoji previsto.** A stack Charm mede largura de exibição corretamente, mas ele deve ficar fora de cálculos manuais sensíveis a alinhamento.

---

## Estados Visuais

Estados visuais definem como o mesmo elemento muda de aparência conforme o contexto.

### Matriz resumida

| Estado | Tratamento visual |
|---|---|
| Normal | `text.primary` sobre `surface.base` |
| Selecionado | `special.highlight` + **bold** |
| Desabilitado | `text.disabled` + dim |
| Marcado para exclusão | `special.muted` + `✕` + ~~strikethrough~~ |
| Favorito | `★` em `accent.secondary` |
| Adicionado / modificado | `+` ou `~` em `semantic.info` |
| Pasta virtual / leitura | `text.secondary` + *italic* |
| Campo sensível revelado | mesmo estilo do texto normal; a diferença é o valor exposto |
| Erro inline | `semantic.error` |

### Regras de transição

- Foco de painel é indicado pelo separador vertical (ver Padrão: Foco e Navegação) e pela barra de comandos — não por borda ao redor do painel.
- TUI não tem estado "pressionado"; confirmação vem por mudança de contexto ou mensagem.
- Transições são instantâneas. A única animação prevista é o spinner `MsgBusy`.

---

## Padrões

Padrões são regras de comportamento transversais — aplicam-se a múltiplas telas e componentes. Os documentos de especificação consomem estes padrões ao definir componentes e fluxos concretos.

---

### Sobreposição

Elementos sobrepostos (modais, diálogos, seletores) seguem regras uniformes de apresentação e interação.

**Apresentação:**

- Centralizados horizontal e verticalmente sobre o conteúdo
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

**Tipos semânticos:**

Elementos sobrepostos carregam intenção semântica que governa o estilo de borda e título:

| Tipo | Símbolo | Papel | Token de borda |
|---|---|---|---|
| Neutro | — | Operação sem urgência | `border.focused` |
| Pergunta | `?` | Decisão neutra | `border.focused` |
| Alerta | `⚠` | Ação potencialmente destrutiva | `semantic.warning` |
| Informação | `ℹ` | Informação que requer reconhecimento | `semantic.info` |
| Erro | — | Validação com falha | `semantic.error` |

**Anatomia comum:**

Todo diálogo — de decisão ou funcional — segue a mesma estrutura de moldura:

```
╭── ⚠  Título ────────────────────╮  ← borda superior: símbolo + título em bold
│                                  │
│  (conteúdo interno do diálogo)   │
│                                  │
╰────────── S Ação ── Esc Cancelar ╯  ← borda inferior: ações alinhadas à direita
```

Regras da moldura:

- **Borda superior** contém o título embutido, precedido pelo símbolo semântico quando aplicável (`?`, `⚠`, `ℹ`)
- **Borda inferior** contém apenas ações de confirmação e cancelamento, alinhadas à direita
- **Ação default** (associada a `Enter`): tecla + label em **bold**, coloridos com o token de destaque do tipo de diálogo (ver tabela abaixo) — visualmente distinta das demais
- **Demais ações**: tecla + label na cor da borda, sem bold
- **Borda e título** usam o mesmo token — definido pela tabela de tipos semânticos
- **Ações internas** (revelar senha, alternar campo, expandir diretório) aparecem exclusivamente na barra de comandos — não na borda do diálogo
- **Teclas de navegação** (↑↓, →, ←, Tab) são de conhecimento amplo e não aparecem na borda
- A **barra de comandos** espelha as ações de confirmação/cancelamento e acrescenta as ações internas e de navegação do diálogo

Token da tecla de ação por tipo de diálogo:

| Tipo | Token da tecla principal | Motivo |
|---|---|---|
| Pergunta | `accent.primary` | Ação neutra, sem risco |
| Alerta | `semantic.error` | Ação destrutiva — cor de perigo reforça a gravidade |
| Informação | `accent.primary` | Apenas reconhecimento, sem risco |
| Neutro / Funcional | `accent.primary` | Padrão para diálogos sem urgência |

**Diálogos de decisão:**

DialogQuestion, DialogAlert e DialogInfo pedem uma decisão do usuário. O conteúdo interno é uma mensagem + contexto:

- **Mensagem** em `text.primary`; nomes de itens referenciados em **bold**

```
╭── ⚠  Excluir segredo ───────────╮
│                                  │
│  Gmail será excluído             │
│  permanentemente.                │
│  Esta ação não pode ser desfeita.│
│                                  │
╰────────── S Excluir ── N Cancelar ╯
```

```
╭── ?  Alterações não salvas ─────╮
│                                  │
│  Deseja salvar antes de sair?    │
│                                  │
╰── S Salvar ── N Descartar ── Esc Voltar ╯
```

```
╭── ℹ  Conflito detectado ───────╮
│                                  │
│  O arquivo foi modificado        │
│  externamente.                   │
│                                  │
╰───────────────────── Enter OK ──╯
```

**Diálogos funcionais:**

PasswordEntry, PasswordCreate, FilePicker e Help oferecem uma função específica em vez de uma decisão sim/não. Compartilham a mesma moldura (título na borda superior, ações na borda inferior) mas diferem no conteúdo interno.

Regras específicas:

- **Título** sem símbolo semântico (não há urgência)
- **Borda** em `border.focused` para diálogos que recebem entrada de texto; `border.default` para diálogos de consulta
- A anatomia interna completa de cada um está documentada na especificação de telas

**PasswordEntry** — entrada de senha única (abrir cofre):

```
╭── Senha mestra ─────────────────╮
│                                  │
│  Senha                           │  ← label em text.secondary
│  ••••••••••▌                     │  ← texto mascarado + cursor
│                                  │
│  Tentativa 2 de 5                │  ← contador em text.secondary
╰────────── Enter Confirmar ── Esc Cancelar ╯
```

**PasswordCreate** — criação de senha com confirmação (criar cofre, alterar senha):

```
╭── Definir senha mestra ─────────╮
│                                  │
│  Nova senha                      │  ← label em accent.primary (campo ativo)
│  ••••••••••▌                     │  ← texto mascarado + cursor
│                                  │
│  Confirmação                     │  ← label em text.secondary (campo inativo)
│                                  │
│  Força: ████████░░ Boa           │  ← medidor de força
╰────────── Enter Confirmar ── Esc Cancelar ╯
```

**FilePicker** — navegação de diretórios com painéis (abrir/salvar):

```
╭── Abrir cofre ───────────────────────────────────────────────────╮
│  ── Diretórios ──────────────┬── /home/usuario/cofres ────────── │
│  ▸ /                         │  cofre.abditum                    │
│    ▾ usuario/                │  pessoal.abditum                  │
│      ► cofres/               │  trabalho.abditum                 │
│      Documents/              │                                   │
│  ────────────────────────────┴────────────────────────────────── │
│  Nome do arquivo               ← só no modo save                │
│  meu-cofre▌                    ← campo de nome com cursor       │
╰───────────────────────── Enter Selecionar ── Esc Cancelar ╯
```

**Help** — tabela de atalhos do contexto ativo:

```
╭── Ajuda — Atalhos e Ações ──────────────────────────────────────╮
│                                                                  │
│  Navegação                        ← grupo em text.secondary bold │
│  ↑↓        Mover cursor                                         │
│  → Enter   Expandir / selecionar                                 │
│  Tab       Alternar painéis                                      │
│                                                                  │
│  Segredo                          ← grupos agrupam por contexto  │
│  F16       Revelar / ocultar                                     │
│  F17       Copiar valor                                          │
│                                                        ↓ mais ── │
╰──────────────────────────────────────────────── Esc Fechar ╯
```

---

### Mensagens

A aplicação comunica feedback ao usuário por meio de uma mensagem exibida na barra de mensagens. Uma mensagem por vez — nova mensagem substitui a anterior imediatamente. Não há fila nem pilha.

**Posição:** sobreposta à última linha da área de trabalho — não reserva linha própria.

**Largura:** ~95% da largura do terminal. Trunca com `…` se necessário.

**Tipos de mensagem:**

| Tipo | Símbolo | Token | Atributo | TTL | Desaparece com input |
|---|---|---|---|---|---|
| Sucesso | `✓` | `semantic.success` | — | 2–3 s | Não |
| Informação | `ℹ` | `semantic.info` | — | 3 s | Não |
| Aviso | `⚠` | `semantic.warning` | — | Permanente | **Sim** |
| Erro | `✗` | `semantic.error` | **bold** | 5 s | Não |
| Ocupado (spinner) | `◐ ◓ ◑ ◒` | `accent.primary` | — | Permanente | Não |
| Dica de campo | `•` | `text.secondary` | *italic* | Permanente | Não |
| Dica de uso | `💡` | `text.secondary` | *italic* | Permanente | Não |

**Regras de comportamento:**

- Mensagem de **Aviso** é re-emitida a cada tick enquanto a condição persistir (ex: bloqueio iminente)
- **Ocupado** permanece até ser substituído por Sucesso ou Erro ao concluir a operação; spinner avança 1 frame/segundo sincronizado com tick global
- **Dica de campo** é substituída ao navegar para outro campo
- **Dica de uso** é substituída pela próxima mensagem de qualquer tipo
- `💡` ocupa 2 colunas — reservar espaço adequado em cálculos de truncamento

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
| Indicador `+` / `~` | `semantic.info` | `+` / `~` (símbolos preservam estado) |
| Busca match | `special.match` + **bold** | **bold** |
| Exclusão `✕` | `semantic.error` + strikethrough | `✕` + strikethrough |
| Favorito `★` | `accent.secondary` | `★` (símbolo preserva semântica) |
| Máscara `••••••••` | `text.secondary` | `text.primary` |
| Borda de modal | `semantic.*` / `border.*` | Borda presente — tipo distinguido por símbolo no título |
| Campo de input | `surface.input` | `surface.base` — borda ainda presente |
