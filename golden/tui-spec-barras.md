# Especificação Visual — Barras

> Barra de Comandos e Barra de Mensagens.
> Parte da [Especificação Visual](tui-specification.md).
>
> **Documento de fundação:**
> - [`tui-design-system.md`](tui-design-system.md) — fundações visuais

## Barra de Comandos

A barra de comandos é a última linha da tela (conforme [DS — Dimensionamento e Layout](tui-design-system.md#dimensionamento-e-layout)). Exibe as ações acionáveis por teclado no contexto atual — o usuário nunca precisa adivinhar o que pode fazer.

**Princípio de conteúdo:** a barra exibe apenas ações de caso de uso (F-keys, atalhos de domínio, `⌃S`). Teclas de navegação universais — `↑↓`, `←→`, `Tab`, `Enter`, `Esc` — são senso comum em TUI e não são exibidas. Exceção: diálogos exibem ações internas específicas do contexto.

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
| `HideFromBar = false` | Exibida (se `Enabled`) | Listada |

Além destes:

- **Prioridade** — valor numérico. Maior prioridade → mais à esquerda. Quando o espaço é insuficiente, ações de menor prioridade são removidas primeiro.
- **Grupo** — valor numérico. Usado exclusivamente no modal de Ajuda para organizar ações. Grupos renderizados em ordem numérica crescente. Dentro de cada grupo, ações ordenadas por `Prioridade`. Não afeta a barra de comandos.
- **Label do grupo** — string registrada por grupo (ex: grupo 1 → "Navegação"). Exibido como título de seção no Help em `text.secondary` **bold**.

Regras de layout:

- **`F1 Ajuda` sempre visível** — âncora fixa na extrema direita; o cálculo de espaço desconta `F1 Ajuda` antes de distribuir as demais ações.
- **Ações desabilitadas desaparecem** — `Enabled = false` remove da barra (não fica dim). A ação continua listada no Help.
- **Ações de confirmação/cancelamento** (`Enter`/`Esc`) já estão na borda inferior do diálogo — não são duplicadas na barra.

### Teclado

A barra funciona como guia de descoberta — o usuário vê quais teclas estão disponíveis no contexto atual. As teclas são atribuídas conforme as [Convenções Semânticas](tui-design-system.md#convenções-semânticas) e os [Escopos](tui-design-system.md#escopos) definidos no DS.

### Eventos

| Evento | Reação da barra |
|---|---|
| Área de trabalho com foco | Exibe ações do painel ativo (árvore ou detalhe) |
| Troca de foco entre painéis | Atualiza para ações do painel que recebe foco |
| Diálogo de decisão aberto | Vazia (apenas `F1 Ajuda`) |
| Diálogo funcional aberto | Exibe ações internas do diálogo |
| Diálogo fecha (pop da pilha) | Volta para ações do contexto anterior |
| Terminal redimensionado | Recalcula ações visíveis (prioridade governa corte) |

## Barra de Mensagens

A barra de mensagens é a zona fixa entre a área de trabalho e a barra de comandos (conforme [DS — Dimensionamento e Layout](tui-design-system.md#dimensionamento-e-layout)). Exibe uma mensagem por vez — nova mensagem substitui a anterior imediatamente. Não há fila nem pilha.

### Anatomia

A barra é uma linha de borda `─` na largura total do terminal. Quando há mensagem ativa, símbolo e texto são embutidos na borda, seguindo a mesma regra de composição da [Borda Superior de diálogos](tui-spec-dialogos.md#anatomia-comum).

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

**Severidade:**

A barra de mensagens utiliza o mesmo sistema de [severidade dos diálogos](tui-spec-dialogos.md#severidade). Cada tipo de mensagem herda o token semântico correspondente — a mesma paleta que governa bordas e símbolos de diálogos governa a cor da mensagem inteira.

| Severidade | Tipo de mensagem | Símbolo | Token | Atributo |
|---|---|---|---|---|
| Erro | Erro | `✕` | `semantic.error` | **bold** |
| Alerta | Alerta | `⚠` | `semantic.warning` | — |
| Informativo | Informação | `ℹ` | `semantic.info` | — |
| Neutro | Sucesso | `✓` | `semantic.success` | — |

> Não existe tipo de mensagem "Destrutivo" na barra — ações destrutivas são sempre comunicadas por diálogos.

**Tipos não-semânticos:**

Além das severidades, a barra suporta tipos utilitários sem correspondência com severidade de diálogos:

| Tipo | Símbolo | Token | Atributo |
|---|---|---|---|
| Ocupado (spinner) | `◐ ◓ ◑ ◒` | `accent.primary` | — |
| Dica de campo | `•` | `text.secondary` | *italic* |
| Dica de uso | `•` | `text.secondary` | *italic* |

**Regras de cor:**

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

> O ciclo de vida da barra em diálogos funcionais (mensagem de contexto ao abrir, dica por campo, limpeza ao fechar) é contrato do orquestrador — documentado em [Diálogos — Funcional](tui-spec-dialogos.md#funcional).

