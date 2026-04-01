# Design da TUI — Abditum

> Decisões visuais e de interação para o pacote `internal/tui`.  
> Complementa `tui-elm-architecture.md` (arquitetura) e `fluxos.md` (comportamento).

---

## Áreas da Tela

A interface é dividida em quatro zonas horizontais empilhadas verticalmente. A work area ocupa todo o espaço restante entre as outras três.

| Zona | Altura | Conteúdo |
|---|---|---|
| **Header** | 1 linha | Nome da aplicação, nome do cofre aberto (se houver), indicador de alterações não salvas |
| **Barra de mensagens** | 1 linha | Feedback tipado via `MessageManager` — sobreposta à work area, não ocupa linha dedicada |
| **Work area** | Restante | Área principal — childModel ativo (boas-vindas, cofre, modelos, configurações) |
| **Command bar** | 1 linha | Ações disponíveis no momento ou atalhos do modal ativo |

> **Status de comprometimento:** a work area é decisão fechada. Header, barra de mensagens e command bar são bastante prováveis mas não comprometidas — a estrutura exata será definida na implementação.

### Header

Faixa no topo da tela com informações de contexto global. Exibe o nome da aplicação e, quando um cofre está aberto, o nome do arquivo e um indicador visual se houver alterações não salvas (ex: `•` ou `[modificado]`). Detalhes visuais (cor, separadores, formatação) serão definidos na implementação.

### Work area

Zona central onde o conteúdo principal é renderizado. Exibe um de quatro modos mutuamente exclusivos:

| Modo | Conteúdo |
|---|---|
| **Boas-vindas** | ASCII art com logo, sem interação além de ações globais |
| **Cofre** | Dois painéis lado a lado — árvore de navegação (esquerda) e detalhe do segredo selecionado (direita) |
| **Modelos** | Dois painéis — lista de modelos (esquerda) e detalhe do modelo (direita) |
| **Configurações** | Painel único ocupando toda a área |

A work area **não muda durante fluxos** — o usuário vê o conteúdo atual com modais sobrepostos. A transição para outro modo só ocorre após um fluxo concluir com sucesso (ex: abrir cofre → transição de boas-vindas para cofre).

### Command bar

Faixa na última linha da tela que mostra as ações disponíveis no contexto atual. Cada ação é exibida como `[tecla] label` (ex: `[F1] Abrir`, `[?] Ajuda`, `[ctrl+Q] Sair`).

**Conteúdo dinâmico:** a command bar alterna entre duas fontes de conteúdo:

- **Sem modal ativo** — exibe ações do `ActionManager.Visible()`, filtradas por `Enabled() == true` e priorizadas pelo espaço disponível. Se houver mais ações do que cabem na largura, exibe as mais relevantes (critério de prioridade será definido na implementação).
- **Com modal ativo** — exibe os atalhos do modal do topo da stack (via `Shortcuts()`). Exemplos: `[Enter] Confirmar` / `[Esc] Cancelar` para modais de senha; `[S] Sim` / `[N] Não` / `[Esc] Cancelar` para confirmações.

**Help modal:** a command bar mostra permanentemente o atalho `[?] Ajuda` (registrado como `ScopeGlobal`). Ao pressionar `?`, abre-se um modal de help que lista **todas** as ações registradas no `ActionManager`, agrupadas por categoria (`Group`). Para cada ação, exibe: todas as teclas alternativas, label e descrição. O help funciona inclusive sobre outros modais (é `ScopeGlobal`) e é fechado via ESC. Isso permite que o usuário descubra funcionalidades sem memorizar atalhos — a command bar mostra o essencial, o help modal mostra tudo.

---

## Barra de Mensagens

### Espaço disponível

A barra de mensagens ocupa **1 linha** com aproximadamente **95% da largura do terminal**. É renderizada **sobre** um elemento existente da interface (provavelmente a última linha da work area), não como linha adicional — o layout do frame não reserva espaço dedicado.

Se o texto exceder a largura disponível, é truncado com reticências (`…`). Não há quebra de linha nem expansão vertical.

### Anatomia da mensagem

```
◐ Salvando cofre...
├─┤ └───────────────┘
 │        texto
emoji
```

Cada mensagem é composta por:
- **Emoji** — indica o tipo visualmente. Para `MsgBusy`, o emoji é animado (rotação a cada segundo).
- **Texto** — descritivo, conciso, em linguagem neutra (sem jargão técnico).

### Tipos de mensagem

| Tipo | Emoji | Cor do texto | Hex | Swatch | Uso |
|---|---|---|---|---|---|
| `MsgInfo` | ✅ | Verde | `#9ece6a` | <span style="color:#9ece6a">██</span> | Operação concluída com sucesso |
| `MsgWarn` | ⚠️ | Amarelo | `#e0af68` | <span style="color:#e0af68">██</span> | Atenção — bloqueio iminente, conflito detectado |
| `MsgError` | ❌ | Vermelho (bold) | `#f7768e` | <span style="color:#f7768e">██</span> | Falha — salvamento, corrupção, operação impossível |
| `MsgBusy` | ◐ ◓ ◑ ◒ | Azul | `#7aa2f7` | <span style="color:#7aa2f7">██</span> | Operação em andamento — spinner rotativo |
| `MsgHint` | 💡 | Cinza (itálico) | `#565f89` | <span style="color:#565f89">██</span> | Explicação contextual — descrição de campo, dica de uso |

A cor é aplicada ao **texto inteiro** (emoji + mensagem) via lipgloss. As cores acima pertencem à paleta Tokyo Night, consistente com o gradiente do ASCII art (`tui-elm-architecture.md` §D-14).

### Comportamento temporal

Cada mensagem tem dois parâmetros que controlam quando desaparece:

| Parâmetro | Efeito |
|---|---|
| `ttlSeconds` | Tempo em segundos até a mensagem expirar automaticamente. `0` = permanente (até ser substituída). |
| `clearOnInput` | Se `true`, a mensagem desaparece no próximo evento de teclado ou mouse. |

Combinações típicas:

| Situação | TTL | clearOnInput | Exemplo |
|---|---|---|---|
| Sucesso de operação | 2–3s | `false` | "Favoritado", "Copiado para a área de transferência" |
| Erro recuperável | 5s | `false` | "Falha ao salvar — arquivo em uso" |
| Progresso | 0 | `false` | "Salvando cofre..." — permanece até conclusão |
| Aviso de bloqueio | 0 | `true` | "Cofre será bloqueado em breve" — some ao interagir |
| Hint de campo | 0 | `false` | Descrição do campo em foco — substituído ao navegar |

### Prioridade e sobreposição

A barra exibe **uma única mensagem por vez**. Não há fila nem stack — a última chamada `Show()` sobrescreve a anterior.

Isso gera um comportamento natural quando combinado com o tick de 1 segundo:

```
Usuário inativo por 45s (75% do timeout)
    → tick: Show(MsgWarn, "Bloqueio iminente", 0, true)     ⚠️ Bloqueio iminente
Usuário copia um campo
    → HandleInput() limpa o warning (clearOnInput)
    → rootModel: Show(MsgInfo, "Copiado", 3, false)         ✅ Copiado
3 segundos depois
    → Tick() expira a mensagem                               (vazio)
1 segundo depois, se inatividade continua
    → tick: Show(MsgWarn, "Bloqueio iminente", 0, true)     ⚠️ Bloqueio iminente
```

O warning de bloqueio é **re-emitido naturalmente** pelo tick a cada segundo enquanto a condição persistir — sem mecanismo de prioridade explícito.

### Animação do spinner (`MsgBusy`)

O spinner avança **1 frame por segundo**, sincronizado com o tick global:

```
Segundo 0:  ◐ Salvando cofre...
Segundo 1:  ◓ Salvando cofre...
Segundo 2:  ◑ Salvando cofre...
Segundo 3:  ◒ Salvando cofre...
Segundo 4:  ◐ Salvando cofre...    (ciclo reinicia)
```

A frequência de 1fps é adequada para uma barra de texto. Não existe modal de progresso — `MsgBusy` é o único indicador de operação em andamento. O `activeFlow` já bloqueia input local (teclas caem no flow, que as ignora), tornando um overlay de progresso redundante.

---

## Modais

### Princípios gerais

- Modais são painéis sobrepostos **acima** de todo o frame, centralizados via `lipgloss.Place()`.
- O conteúdo por trás do modal permanece visível mas **não recebe input** — apenas o modal do topo recebe teclas e mouse.
- Modais se auto-dimensionam pelo conteúdo — não recebem tamanho alocado do compositor.
- A **command bar** troca de conteúdo durante modal: exibe os atalhos do modal ativo (`Shortcuts()`) em vez das ações do `ActionManager`.

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

| Tipo | Emoji | Cor base | Hex (Tokyo Night) | Swatch | Uso típico |
|---|---|---|---|---|---|
| `DialogQuestion` | ❓ | Azul | `#7aa2f7` | <span style="color:#7aa2f7">██</span> | Escolhas neutras: salvar/descartar/cancelar, sobrescrever |
| `DialogAlert` | ⚠️ | Amarelo | `#e0af68` | <span style="color:#e0af68">██</span> | Ações destrutivas ou irreversíveis: excluir, descartar |
| `DialogInfo` | ℹ️ | Ciano | `#7dcfff` | <span style="color:#7dcfff">██</span> | Informações que pedem reconhecimento: política, aviso |

A cor base é aplicada à **borda ou título** do modal (definição visual exata adiada para fase de implementação). O emoji é exibido junto ao título.

##### Variantes pré-definidas

| Factory | Opções | Default (ENTER) | Cancel (ESC) | Caso de uso |
|---|---|---|---|---|
| `Confirm` | Sim / Não | Sim | Não | F23 excluir segredo, F31 excluir pasta, F35 excluir modelo |
| `ConfirmOrCancel` | Sim / Não / Cancelar | Sim | Cancelar | F5 sair com alterações (salvar/descartar/voltar) |
| `Ask` | Opções customizadas | Configurável | Configurável | F8 conflito externo (sobrescrever/salvar como/voltar) |

##### Composição visual

- **Título** — frase curta, exibida com o emoji do `DialogType` (ex: "⚠️ Excluir segredo").
- **Mensagem** — texto explicativo abaixo do título. Pode ser vazio se o título for autoexplicativo. Pode ter múltiplas linhas.
- **Opções** — exibidas como botões ou labels com atalhos entre parênteses.
- A opção `Default` é visualmente destacada (ex: bold ou cor diferente).
- Atalhos de cada opção aparecem na command bar via `Shortcuts()`.

**Mecanismo:** callback-based — o contexto da decisão já é conhecido no momento da abertura. Adequado para children e flows.

#### Mensagem informativa

Apresenta um título e texto descritivo, com tipo semântico (`DialogType`). Sem opções — apenas dismiss.

| Factory | Dismiss | Caso de uso |
|---|---|---|
| `Message` | ESC ou ENTER | F2 senha fraca (aviso — `DialogAlert`), F13 política de mesclagem (`DialogInfo`) |

**Composição visual:**

- Título com emoji do tipo (ex: "ℹ️ Política de mesclagem", "⚠️ Senha fraca").
- Texto livre abaixo — pode ter múltiplas linhas.
- Cor da borda/título segue o `DialogType`.
- Command bar mostra apenas `[Enter] OK` / `[Esc] Fechar`.

#### Entrada de senha

Campos de senha com máscara (caracteres ocultos). Dados sensíveis (`[]byte`) — nunca transitam via callback.

| Factory | Campos | Validação interna | Caso de uso |
|---|---|---|---|
| `PasswordEntry` | 1 campo | Nenhuma | F1 abrir cofre |
| `PasswordCreate` | 2 campos (senha + confirmação) | Verifica se coincidem | F2 criar cofre, F11 alterar senha |

**Composição visual:**

- Título descritivo (ex: "Senha mestra", "Definir senha mestra").
- Campo(s) com máscara — caracteres substituídos por `•` ou `*`.
- `PasswordCreate`: se os campos não coincidem ao confirmar, mostra erro inline sem fechar o modal. O usuário corrige e tenta novamente.
- Command bar: `[Enter] Confirmar` / `[Esc] Cancelar`.

**Mecanismo:** `modalResult` — roteado exclusivamente ao `activeFlow`. Bytes de senha nunca entram em broadcast.

#### File picker

Navegação de diretórios para escolher um arquivo.

| Factory | Modo | Caso de uso |
|---|---|---|
| `FilePicker` (open) | Selecionar arquivo existente | F1 abrir cofre, F13 importar |
| `FilePicker` (save) | Escolher destino para escrita | F2 criar cofre, F9 salvar como, F12 exportar |

**Composição visual:**

- Título descritivo (ex: "Abrir cofre", "Salvar cofre como").
- Navegação de diretórios com listagem de arquivos.
- Modo save: campo para informar nome do arquivo. Extensão (ex: `.abditum`) adicionada automaticamente se omitida.
- Modo open: filtra por extensão quando aplicável.
- Command bar: `[Enter] Selecionar` / `[Esc] Cancelar` + atalhos de navegação.

**Mecanismo:** `modalResult` — roteado ao `activeFlow`.

#### Entrada de texto

Campo de texto livre com validação opcional.

| Factory | Validação | Caso de uso |
|---|---|---|
| `TextInput` | `validate func(string) error` — chamada a cada ENTER | F18 nome do segredo, F27 criar pasta, F28 renomear pasta, F33 renomear modelo |

**Composição visual:**

- Título descritivo (ex: "Nome do segredo", "Nova pasta").
- Placeholder opcional (texto sugestivo quando vazio).
- Se `validate` retornar erro, mensagem de erro exibida inline — modal não fecha.
- Command bar: `[Enter] Confirmar` / `[Esc] Cancelar`.

**Mecanismo:** `modalResult` — roteado ao `activeFlow`.

#### Seleção em lista

Lista de opções navegável com filtro.

| Factory | Navegação | Caso de uso |
|---|---|---|
| `Select` | ↑↓ para mover, ENTER para confirmar, ESC para cancelar | F18 escolher modelo, F25 mover segredo (pasta destino), F29 mover pasta |

**Composição visual:**

- Título descritivo (ex: "Escolher modelo", "Mover para pasta").
- Lista de itens com cursor de seleção.
- Item selecionado visualmente destacado.
- Command bar: `[Enter] Selecionar` / `[Esc] Cancelar` / `[↑↓] Navegar`.

**Mecanismo:** `modalResult` — roteado ao `activeFlow`.

#### Help

Exibe todas as ações registradas no `ActionManager`, agrupadas por `Group`.

| Factory | Acionamento | Scope |
|---|---|---|
| `helpModal` | `?` | `ScopeGlobal` — funciona mesmo durante flows ou sobre outros modais |

**Composição visual:**

- Ações agrupadas por `Group` (ex: "Segredo", "Pasta", "Cofre", "Navegação").
- Cada ação mostra: todas as teclas, label, description.
- Scroll vertical se a lista exceder a altura disponível.
- Command bar: `[Esc] Fechar`.
- Dismiss via ESC — sem ENTER (não há opção a selecionar).

### Stack de modais

Modais podem se sobrepor (ex: help aberto sobre confirmação). Regras:

- O modal do topo recebe input. Modais abaixo permanecem vivos mas não recebem teclas.
- Modais abaixo continuam recebendo mensagens de domínio via broadcast (ex: `tickMsg`).
- A command bar sempre reflete os atalhos do modal do **topo** da stack.

---

## Painéis Divididos

Nos modos **Cofre** e **Modelos**, a work area é dividida horizontalmente em dois painéis lado a lado:

| Modo | Painel esquerdo | Painel direito |
|---|---|---|
| Cofre | Árvore de navegação (pastas e segredos) | Detalhe do segredo selecionado |
| Modelos | Lista de modelos | Detalhe do modelo selecionado |

### Proporção

A proporção entre os painéis será definida na implementação. Opções em consideração:
- Proporção fixa (ex: 30/70 ou 1/3 + 2/3).
- Proporção adaptável à largura do terminal (mais espaço para a árvore em terminais largos).
- Redimensionamento pelo usuário não é previsto neste momento.

### Separador

Os painéis são visualmente separados por uma borda vertical ou espaço. O estilo exato (linha, cor, largura) será definido na implementação.

### Ausência de painel direito

Se nenhum segredo (ou modelo) estiver selecionado, o painel direito pode exibir um placeholder (texto contextual ou espaço vazio). Comportamento definido na fase que implementa o conteúdo.

---

## Foco entre Painéis

Quando dois painéis estão visíveis, apenas um tem **foco ativo** — recebe teclas de navegação e determina quais ações do `ActionManager` estão disponíveis (via `SetActiveOwner`).

### Indicação visual

O painel com foco é visualmente identificado. Opções em consideração:
- Borda do painel ativo com cor diferente (ex: azul) vs. borda do painel inativo (ex: cinza).
- Título do painel ativo em destaque (bold ou cor).
- Definição visual exata adiada para implementação.

### Tecla de troca

A tecla que alterna foco entre os painéis será definida na implementação. Candidatas: `Tab`, `←`/`→`, ou outra tecla que não entre em conflito com ações de domínio. A troca de foco chama `ActionManager.SetActiveOwner(painel)`, fazendo com que a command bar atualize automaticamente.

---

## Paleta de Cores

A TUI usa a paleta **Tokyo Night** como base. As cores já aparecem em contextos específicos (mensagens, modais, ASCII art) e são consolidadas aqui como referência.

### Cores semânticas

| Uso | Cor | Hex | Swatch |
|---|---|---|---|
| Sucesso / Info | Verde | `#9ece6a` | <span style="color:#9ece6a">██</span> |
| Atenção / Alerta | Amarelo | `#e0af68` | <span style="color:#e0af68">██</span> |
| Erro / Destrutivo | Vermelho | `#f7768e` | <span style="color:#f7768e">██</span> |
| Ação principal / Pergunta | Azul | `#7aa2f7` | <span style="color:#7aa2f7">██</span> |
| Informação / Neutro | Ciano | `#7dcfff` | <span style="color:#7dcfff">██</span> |
| Hint / Secundário | Cinza | `#565f89` | <span style="color:#565f89">██</span> |

### Gradiente do logo

O ASCII art da tela de boas-vindas usa um gradiente de 5 cores, uma por linha:

| Linha | Cor | Hex | Swatch |
|---|---|---|---|
| 1 | Violeta | `#9d7cd8` | <span style="color:#9d7cd8">██</span> |
| 2 | Ciano claro | `#89ddff` | <span style="color:#89ddff">██</span> |
| 3 | Azul | `#7aa2f7` | <span style="color:#7aa2f7">██</span> |
| 4 | Ciano | `#7dcfff` | <span style="color:#7dcfff">██</span> |
| 5 | Lilás | `#bb9af7` | <span style="color:#bb9af7">██</span> |

### Cores de UI

Cores para bordas, texto normal, fundo e demais elementos de UI serão definidas na implementação, mantendo consistência com a paleta Tokyo Night. Candidatas:

| Elemento | Hex | Swatch | Papel na paleta Tokyo Night |
|---|---|---|---|
| Foreground (texto) | `#a9b1d6` | <span style="color:#a9b1d6">██</span> | Cor principal de texto — tom azulado claro |
| Background (fundo) | `#1a1b26` | <span style="background:#1a1b26;color:#1a1b26">██</span> | Fundo escuro profundo — azul-noite |
| Bordas inativas | `#414868` | <span style="color:#414868">██</span> | Cinza-azulado médio — separadores e bordas sem foco |
| Comentário / desabilitado | `#565f89` | <span style="color:#565f89">██</span> | Cinza muted — texto secundário, itens desabilitados |
| Seleção / destaque | `#283457` | <span style="color:#283457">██</span> | Azul escuro — fundo de item selecionado |
| Bordas ativas | `#7aa2f7` | <span style="color:#7aa2f7">██</span> | Azul vibrante — painel com foco |

### Referência: paleta Tokyo Night

A [Tokyo Night](https://github.com/enkia/tokyo-night-vscode-theme) é um tema de cores escuro com tons predominantemente azuis e roxos. Principais características:

- **Fundo:** azul-noite muito escuro (`#1a1b26`), não preto puro — reduz contraste agressivo.
- **Texto:** azul-acinzentado claro (`#a9b1d6`) — confortável para leitura prolongada.
- **Destaques:** cores vibrantes mas dessaturadas — projetadas para serem distintas sem cansar.
- **Semântica consistente:** verde para sucesso, amarelo para warning, vermelho para erro, azul para ação primária, ciano para informação — convenções universais mantidas pela paleta.

O Abditum não implementa o tema completo — seleciona um subconjunto de cores que cobrem as necessidades de uma TUI com foco em legibilidade e distinção semântica.

---

## Tela de Boas-vindas

Exibida no modo `workAreaWelcome` — quando nenhum cofre está aberto. É a primeira tela que o usuário vê ao iniciar a aplicação.

### Conteúdo

- **Logo** — ASCII art de 5 linhas com o nome "Abditum", renderizado com gradiente de cores (ver paleta acima). Centralizado horizontal e verticalmente na work area.
- **Versão** — número da versão abaixo do logo (formato a definir).
- **Ações disponíveis** — a command bar exibe as ações relevantes no contexto de boas-vindas: `[F1] Abrir cofre`, `[F2] Criar cofre`, `[?] Ajuda`, `[ctrl+Q] Sair`.

### Interação

Sem navegação interna — a tela é estática. O usuário interage exclusivamente via ações da command bar. Abrir ou criar um cofre com sucesso transiciona para `workAreaVault`.

---

## Campos Sensíveis

Segredos contêm campos cujo valor é protegido visualmente (senhas, tokens, chaves de API).

### Máscara

Campos sensíveis são exibidos com máscara por padrão — caracteres substituídos por `•` (ex: `••••••••`). A quantidade de bullets pode ser fixa (não revelar o comprimento) ou proporcional ao conteúdo real — decisão adiada para implementação.

### Revelar

O usuário pode revelar temporariamente o valor de um campo sensível (ação F16). O campo revelado volta a ser mascarado automaticamente após o timeout configurado (`lastRevealAt` + duração no `vault.Manager`).

### Auto-hide

Quando o timer expira, o `rootModel` emite `fieldHideMsg{}` e o painel de detalhe oculta os campos revelados. O comportamento é coordenado com o lock: se o cofre é bloqueado, todos os campos são ocultados como parte do wipe de memória sensível.

---

## Indicador de Alterações Não Salvas

Quando o cofre tem alterações em memória que não foram gravadas em disco, o header exibe um indicador visual junto ao nome do arquivo.

### Formato

O formato exato será definido na implementação. Opções em consideração:
- Prefixo ou sufixo no nome do arquivo (ex: `• cofre.abditum` ou `cofre.abditum [modificado]`).
- Cor diferente no nome do arquivo.
- Combinação de ambos.

### Origem do estado

O `rootModel` consulta `vault.Manager` para saber se há alterações não salvas (ex: `mgr.IsDirty()`). O indicador é atualizado a cada re-render — sem mensagem especial.

---

## Navegação na Árvore

O painel esquerdo no modo **Cofre** exibe a hierarquia de pastas e segredos como uma árvore navegável.

### Estrutura visual

- Pastas e segredos são exibidos com indentação proporcional à profundidade na hierarquia.
- Pastas podem ser expandidas/recolhidas (expand/collapse). O estado expand/collapse é local ao `vaultTreeModel`.
- Distinção visual entre pastas e segredos: prefixo, ícone ou formatação diferente. Opções: `📁`/`🔑`, `▸`/`•`, ou apenas indentação + nome. Definição na implementação.

### Estados visuais de itens

| Estado | Indicação visual |
|---|---|
| Selecionado (cursor) | Destaque de fundo ou cor de texto — o item que tem foco |
| Favorito | Indicador visual (ex: `★` ou cor) |
| Marcado para exclusão | Estilo diferenciado (ex: ~~tachado~~, cor apagada, ou prefixo `✕`) |
| Pasta expandida | Prefixo `▾` ou similar |
| Pasta recolhida | Prefixo `▸` ou similar |

### Navegação

- `↑`/`↓` — mover cursor entre itens visíveis.
- `Enter` ou `→` — expandir pasta / selecionar segredo (exibir no painel de detalhe).
- `←` — recolher pasta ou voltar para a pasta pai.
- Detalhes da navegação por teclado serão definidos na implementação.

---

## Overlay de Modais

Modais são sobrepostos acima de todo o frame via `lipgloss.Place()`. O conteúdo da interface permanece visível por trás do modal.

### Fundo

O conteúdo por trás do modal **não é escurecido** (sem dim/fade) na versão inicial. TUIs têm capacidades limitadas de transparência — dim exigiria re-renderizar o frame com cores alteradas, o que adiciona complexidade. Decisão revisável na implementação se o contraste for insuficiente.

### Borda do modal

Modais têm borda visível para separação do conteúdo de fundo. Estilo (linha simples, dupla, arredondada) e cor (pode seguir o `DialogType` para confirmações) serão definidos na implementação.

### Posicionamento

Modais são centralizados horizontal e verticalmente no terminal. Não há posicionamento relativo ao conteúdo que disparou o modal.
