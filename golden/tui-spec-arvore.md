# Especificação Visual — Árvore de Segredos

> Painel esquerdo, busca e ações na árvore.
> Parte da [Especificação Visual](tui-specification.md).
>
> **Documentos relacionados:**
> - [`tui-design-system.md`](tui-design-system.md) — fundações visuais
> - [`tui-spec-detalhe.md`](tui-spec-detalhe.md) — painel direito (detalhe do segredo)

### Painel Esquerdo: Árvore

**Contexto:** Área de trabalho — Modo Cofre.
**Largura:** ~35% da área de trabalho.
**Responsabilidade:** Exibir a hierarquia de pastas e segredos; permitir navegação e seleção do item a detalhar no painel direito.

**Wireframe (Modo Cofre — scroll ativo, segredo selecionado, painel com foco):**

```
  ▼ Favoritos          (2) ↑
      ★ Bradesco              │
      ★ Gmail                 │
  ▼ Geral              (8)  ■
    ▼ Sites e Apps     (5)  │
      ● Gmail           <╡      ← special.highlight + bold (item selecionado)
      ● YouTube              │
      ● Facebook             │
  ▼ Financeiro         (3)  │
    ● Nubank                 ↓
```

> `↑`/`↓` indicam conteúdo além da área visível; `■` é o thumb proporcional na posição `│`; `<╡` marca o item sendo detalhado no painel direito. `<╡` e scroll (`↑`/`↓`/`■`) ocupam a mesma coluna — o separador entre painéis. Quando `<╡` coincide com um indicador de scroll na mesma linha, `<╡` tem prioridade (o indicador de scroll é suprimido naquela linha).

**Wireframe (item marcado para exclusão — selecionado):**

```
    ▼ Sites e Apps     (5)  │
      ✗ Gmail           <╡      ← special.highlight; `semantic.warning` + strikethrough
      ● YouTube              │
```

**Wireframe (cofre vazio):**

```
  ▷ Geral              (0)  │   ← special.highlight (pasta raiz selecionada)
                             │
                             │
```

Painel direito exibe placeholder "Cofre vazio" centralizado quando o cofre não tem nenhum segredo.

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Nome de item (normal) | `text.primary` | — |
| Fundo de item selecionado | `special.highlight` | — |
| Nome de item selecionado | `text.primary` | **bold** |
| `▼ ▶ ▷` — prefixos de pasta | `text.secondary` | — |
| `●` — prefixo de segredo | `text.secondary` | — |
| `★` — prefixo de segredo favoritado | `accent.secondary` | — |
| `★` — prefixo de itens dentro de `▼ Favoritos` | `accent.secondary` | — |
| Nome da pasta virtual `Favoritos` | `accent.primary` | **bold** |
| Contadores `(n)` | `text.secondary` | — |
| Nome de segredo marcado para exclusão | `semantic.warning` | ~~strikethrough~~ |
| `✗` — prefixo de segredo marcado para exclusão | `semantic.warning` | — |
| Nome de segredo recém-criado (não salvo) | `semantic.warning` | — |
| `✦` — prefixo de segredo recém-criado | `semantic.warning` | — |
| Nome de segredo modificado (não salvo) | `semantic.warning` | — |
| `✎` — prefixo de segredo modificado | `semantic.warning` | — |
| Nome de item desabilitado | `text.disabled` | dim |
| `│` separador — painel com foco | `border.focused` | — |
| `│` separador — painel sem foco | `border.default` | — |
| `<╡` conector de seleção no separador | `accent.primary` | — |
| `↑` / `↓` indicadores de scroll no `│` | `text.secondary` | — |
| `■` thumb de scroll no `│` | `text.secondary` | — |

#### Estados dos componentes

| Componente | Estado | Condição |
|---|---|---|
| `Favoritos` | visível, expandível (`▼/▶`) | ≥ 1 segredo favoritado |
| `Favoritos` | oculta | 0 segredos favoritados |
| Pasta ou segredo | `special.highlight` + texto **bold** | Cursor posicionado sobre o item |
| Pasta com filhos, expandida | prefixo `▼` em `text.secondary` | Pasta não-vazia, aberta |
| Pasta com filhos, recolhida | prefixo `▶` em `text.secondary` | Pasta não-vazia, fechada |
| Pasta sem filhos | prefixo `▷` em `text.secondary` | Pasta vazia |
| Segredo (folha, limpo) | prefixo `●` em `text.secondary` | Segredo sem alterações pendentes |
| Segredo recém-criado | prefixo `✦` em `semantic.warning` + texto `semantic.warning` | Criado em memória, ainda não salvo em disco |
| Segredo modificado | prefixo `✎` em `semantic.warning` + texto `semantic.warning` | Editado em memória, ainda não salvo em disco |
| Segredo marcado para exclusão | prefixo `✗` em `semantic.warning` + texto `semantic.warning` + ~~strikethrough~~ | Marcado para exclusão, ainda não salvo |
| `<╡` no separador | visível | Foco da árvore está sobre um segredo |
| `<╡` no separador | ausente — `│` normal | Nenhum segredo exibido no painel direito |
| `↑`/`↓`/`■` no `│` | visível | Conteúdo excede a área visível do painel |
| Painel esquerdo | placeholder "Cofre vazio" à direita | Cofre sem nenhum segredo |

> **`<╡` × `■`:** quando o item selecionado coincide com a posição do thumb, `<╡` tem prioridade — mesma regra do DS para sobreposição em bordas.

#### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Painel recebe foco | Dica de campo | `• ↑↓ para navegar` |
| `Favoritos` (a pasta) selecionada | Dica de campo | `• Pasta virtual — segredos permanecem na localização original` |

#### Eventos

**Navegação:**

**Navegação — movimento linear:**

| Evento | Efeito na árvore |
|---|---|
| Cursor desce uma linha | Foco move para o próximo item visível (respeitando expand/collapse); se já está no último item, não move |
| Cursor sobe uma linha | Foco move para o item anterior visível; se já está no primeiro item, não move |
| Cursor vai ao primeiro item | Foco move para o topo absoluto da árvore (primeiro item da lista, independente do scroll) |
| Cursor vai ao último item | Foco move para o último item visível da árvore |
| Scroll desce uma página | Janela desliza viewport − 1 linhas para baixo; cursor vai para o item no topo da nova janela se estava fora dela |
| Scroll sobe uma página | Janela desliza viewport − 1 linhas para cima; cursor vai para o item no fundo da nova janela se estava fora dela |

**Navegação — movimento hierárquico:**

| Evento | Efeito na árvore |
|---|---|
| Avançar sobre pasta recolhida (`▶`) | Pasta expandida; filhos tornam-se visíveis; prefixo `▶` → `▼`; foco salta para o primeiro filho visível (subpasta ou segredo) |
| Avançar sobre pasta expandida (`▼`) | Foco desce para o primeiro filho da pasta |
| Avançar sobre pasta vazia (`▷`) | Sem efeito — pasta vazia não tem filhos para expandir |
| Avançar sobre segredo | Sem efeito de navegação na árvore — painel direito já exibe o detalhe pelo foco |
| Recuar sobre filho de pasta | Foco sobe para a pasta pai |
| Recuar sobre pasta expandida | Pasta recolhida; prefixo `▼` → `▶`; foco permanece na pasta |
| Recuar sobre pasta raiz (`Geral`) recolhida | Sem efeito — sem pai disponível |
| Recuar sobre pasta raiz (`Geral`) expandida | Pasta recolhida; foco permanece na pasta raiz |

**Navegação — foco entre painéis:**

| Evento | Efeito na árvore |
|---|---|
| Foco alternado para painel direito | `│` muda de `border.focused` para `border.default`; barra de comandos exibe ações do painel direito |
| Foco recebido do painel direito | `│` muda de `border.default` para `border.focused`; barra de comandos exibe ações da árvore; cursor de campo vai para o item que estava com foco quando a árvore perdeu foco |

**Navegação — scroll visual:**

| Evento | Efeito na árvore |
|---|---|
| Item em foco sai da área visível (scroll para cima) | Janela rola automaticamente para manter o item em foco visível |
| Item em foco sai da área visível (scroll para baixo) | Janela rola automaticamente para manter o item em foco visível |
| Conteúdo total cabe na área visível | Indicadores `↑`/`↓`/`■` desaparecem do `│` |
| Conteúdo total não cabe na área visível | `↑` aparece se há conteúdo acima; `↓` aparece se há conteúdo abaixo; `■` posicionado proporcionalmente |

**Navegação — mouse:**

| Evento | Efeito na árvore |
|---|---|
| Clique em item | Foco move para o item clicado (mesmo efeito de cursor com `↑`/`↓`) |
| Clique no prefixo `▶` ou `▼` | Pasta expande/recolhe — mesmo efeito de `→`/`←` sobre pasta |
| Clique no prefixo `▷` | Sem efeito |
| Scroll do mouse para cima/baixo | Janela desliza; cursor acompanha se sair da área visível |
| Clique em item dentro de `Favoritos` | Foco move para o atalho dentro de `Favoritos`; painel direito exibe o segredo referenciado |

**Navegação — `Favoritos`:**

| Evento | Efeito na árvore |
|---|---|
| Foco entra em `Favoritos` (pasta virtual) | Painel direito mantém último segredo exibido; barra exibe dica "Pasta virtual — segredos permanecem na localização original" |
| `Favoritos` expandida | Atalhos dos segredos favoritados tornam-se visíveis; prefixo `▶` → `▼` |
| `Favoritos` recolhida | Atalhos ocultados; prefixo `▼` → `▶` |
| Foco em atalho dentro de `Favoritos` | Painel direito exibe o detalhe do segredo referenciado; `<╡` aparece na linha do atalho |

**Segredo — criação e duplicação:**

| Evento | Efeito na árvore |
|---|---|
| Novo segredo criado (foco em pasta) | Nó `✦ <novo>` inserido no final da pasta em foco; foco salta para o novo nó; contador da pasta e ancestrais +1 |
| Novo segredo criado (foco em segredo) | Nó `✦ <novo>` inserido imediatamente abaixo do segredo em foco; foco salta para o novo nó; contador da pasta e ancestrais +1 |
| Segredo duplicado | Nó `✦ <nome> (2)` inserido imediatamente abaixo do segredo original; foco salta para o duplicado; contador da pasta e ancestrais +1 |

**Segredo — edição de conteúdo:**

| Evento | Efeito na árvore |
|---|---|
| Nome do segredo alterado | Nome do nó atualizado imediatamente; se era `●`, prefixo muda para `✎`; se já era `✦`, permanece `✦` |
| Campo ou observação editado | Prefixo muda de `●` para `✎` (apenas se `EstadoOriginal`; `✦` permanece `✦`) |

**Segredo — exclusão e restauração:**

| Evento | Efeito na árvore |
|---|---|
| Segredo marcado para exclusão | Prefixo → `✗`; texto `semantic.warning` + strikethrough; contador da pasta e ancestrais −1; se favoritado, some de `Favoritos` |
| Exclusão cancelada (restauração) | Prefixo original restaurado (`●`, `★`, `✦` ou `✎`); texto normal; contador da pasta e ancestrais +1; se era favoritado, volta a `Favoritos` |

**Segredo — favorito:**

| Evento | Efeito na árvore |
|---|---|
| Segredo favoritado | Prefixo `●` → `★` (se limpo); se já era `✦` ou `✎`, prefixo dirty mantido (ver regra de prioridade em Comportamento); `Favoritos` aparece se era a primeira marcação; atalho inserido em `Favoritos` |
| Segredo desfavoritado | Prefixo `★` → `●` (se limpo); atalho removido de `Favoritos`; `Favoritos` desaparece se contagem chegar a 0 |

**Segredo — reordenação e movimentação:**

| Evento | Efeito na árvore |
|---|---|
| Segredo subido uma posição na pasta | Nó sobe uma posição dentro da pasta; foco acompanha |
| Segredo descido uma posição na pasta | Nó desce uma posição dentro da pasta; foco acompanha |
| Segredo reposicionado para posição específica | Nó move para a nova posição dentro da pasta; foco acompanha |
| Segredo movido para outra pasta | Nó some da pasta de origem; aparece na pasta destino na posição especificada; foco acompanha o nó na nova posição; contadores de origem (−1) e destino (+1) e respectivos ancestrais atualizados |

**Pasta — criação e renomeação:**

| Evento | Efeito na árvore |
|---|---|
| Pasta criada | Nó `▷ <nome>` inserido na posição especificada dentro do pai; foco salta para o novo nó |
| Pasta renomeada | Nome do nó atualizado imediatamente |

**Pasta — reordenação e movimentação:**

| Evento | Efeito na árvore |
|---|---|
| Pasta subida uma posição | Nó sobe uma posição entre os irmãos; foco acompanha |
| Pasta descida uma posição | Nó desce uma posição entre os irmãos; foco acompanha |
| Pasta reposicionada para posição específica | Nó move para a nova posição entre os irmãos; foco acompanha |
| Pasta movida para outro pai | Nó some da posição atual; aparece dentro do novo pai; foco acompanha; hierarquia do novo pai atualizada |

**Pasta — exclusão:**

| Evento | Efeito na árvore |
|---|---|
| Pasta excluída (sem conflitos de nome) | Nó da pasta removido; subpastas e segredos promovidos ao pai na posição da pasta excluída; contadores do pai recalculados; foco vai para o primeiro filho promovido (ou para o pai, se pasta era vazia) |
| Pasta excluída (com conflitos de nome) | Idem acima; segredos com conflito de nome exibidos com nome renomeado (sufixo `(N)`); barra de mensagens exibe alerta com lista de renomeações |

**Cofre — persistência:**

| Evento | Efeito na árvore |
|---|---|
| Salvo com sucesso (mesmo arquivo) | Nós `✗` removidos fisicamente da árvore; prefixos `✦` e `✎` voltam a `●` ou `★` conforme o flag `favorito`; contadores recalculados; foco permanece no item atual |
| Salvo como (arquivo diferente) | Efeitos idênticos ao salvar com sucesso — a árvore não distingue o destino do arquivo |
| Salvo com outra senha | Efeitos idênticos ao salvar com sucesso — a árvore não conhece a chave de cifragem |
| Reverter alterações (recarregar do disco) | Árvore completamente reconstruída a partir do arquivo em disco: nós `✦` removidos (não existem no disco); nós `✎` voltam ao nome e prefixo originais (`●` ou `★`); nós `✗` voltam ao prefixo original (`●` ou `★`); contadores recalculados; se o item em foco ainda existe, foco permanece nele; se o item em foco era `✦` (deixou de existir), foco vai para a pasta pai; `Favoritos` reconstruída a partir dos dados do disco |

#### Comportamento

- **Espelho do cofre** — a árvore é uma representação visual direta e sempre atualizada do estado do cofre. Qualquer mutação no cofre — independentemente de onde ou como foi originada — deve se refletir imediatamente na árvore. Não existe estado interno da árvore que contradiga o cofre.
- **Foco persiste sobre o mesmo elemento** — quando qualquer evento atualiza a árvore (reordenação, renomeação, movimentação, exclusão de outro item, salvar, reverter…), o foco permanece sobre o mesmo elemento, mesmo que sua posição na lista tenha mudado. O scroll se ajusta automaticamente para garantir que o elemento com foco esteja visível.
- **Foco ao remover o elemento focado** — se o evento for a remoção do próprio elemento com foco, o foco migra automaticamente seguindo a ordem de preferência: (1) elemento imediatamente abaixo na lista visível; (2) se não existir, elemento imediatamente acima; (3) se a lista ficou vazia, `▼ Geral` (pasta raiz, que nunca pode ser removida).
- **Seleção apenas por cor** — não há símbolo de cursor. A seleção é indicada exclusivamente pelo fundo `special.highlight`. Os prefixos (`▼ ▶ ▷ ● ★ ✦ ✎ ✗`) são estruturais e não mudam com a seleção
- **Detalhe automático** — o painel direito exibe o segredo que está com foco na árvore. Quando o foco está sobre uma pasta, o painel mantém o último segredo exibido. O detalhe não precisa ser "aberto" — é atualizado continuamente conforme o foco se move
- **Nome inicial de novo segredo** — `<novo>`; é o nome provisório que aparece no nó até que o usuário edite o campo Nome no painel de detalhes
- **Segredos com alterações pendentes** — três prefixos indicam estado não salvo, todos em `semantic.warning` (mesma semântica do `•` dirty no cabeçalho): `✦` recém-criado, `✎` modificado, `✗` marcado para exclusão (+ strikethrough). Todos desaparecem após `^S` bem-sucedido
- **`Favoritos` — posição e comportamento** — quando visível, é sempre o primeiro item da lista; se comporta como pasta normal (`▼/▶`); itens internos são atalhos para os segredos originais (os segredos permanecem na hierarquia de origem)
- **`Favoritos` — aparição e remoção** — o nó aparece instantaneamente quando o primeiro segredo é favoritado; desaparece instantaneamente quando o último segredo favoritado é desfavoritado. A atualização segue o princípio "Espelho do cofre" — a árvore reflete o estado do cofre imediatamente após a execução da ação
- **Foco preservado ao inserir/remover `Favoritos`** — quando o nó `Favoritos` aparece ou desaparece, a posição absoluta de todos os itens na lista desloca ±1. O foco permanece sobre o mesmo elemento lógico (identificado por identidade, não por índice). O scroll se ajusta automaticamente para manter o elemento em foco visível
- **Favorito com estado dirty** — o prefixo dirty (`✦`, `✎`, `✗`) substitui o `★` dentro de `Favoritos`; o `★` só aparece como prefixo quando o segredo está limpo. Prioridade de prefixo: `✗` > `✎` > `✦` > `★` > `●`. Segredo marcado para exclusão some imediatamente de `Favoritos` — permanece na hierarquia de origem com prefixo `✗`
- **Navegação linear ignora expand/collapse** — `↑`/`↓` navegam apenas entre itens *visíveis*; filhos de pastas recolhidas são invisíveis e portanto pulados
- **`→` sobre segredo é no-op** — segredos são folhas; avançar sobre eles não tem efeito (o detalhe já foi atualizado ao receber foco)
- **`←` tem dois comportamentos** — sobre pasta expandida, recolhe a pasta e foco permanece na pasta; sobre qualquer outro item (pasta recolhida, pasta vazia, segredo), sobe o foco para a pasta pai. Sobre a pasta raiz expandida, apenas recolhe
- **Foco ao retornar ao painel** — ao receber foco via Tab, o cursor restaura a posição anterior (não vai ao topo)
- **Scroll automático** — o viewport se ajusta automaticamente para manter o item em foco visível; nunca há item em foco fora da área visível
- **Scroll no separador** — o scroll da árvore é indicado por `↑`/`↓`/`■` no `│` (separador entre painéis). `<╡` e scroll ocupam a mesma coluna: `<╡` tem prioridade sobre `■` em caso de coincidência (ver [DS — Scroll em diálogos](tui-design-system.md#scroll-em-diálogos)). Quando `<╡` coincide com `↑` ou `↓`, `<╡` prevalece — a direção do scroll é implícita pela presença do outro indicador nas demais linhas
- **Indentação** — 2 espaços por nível de aninhamento

---

### Busca de Segredos

**Contexto de uso:** filtrar a árvore de segredos por texto livre no Modo Cofre.
**Escopo:** disponível apenas no **Modo Cofre**, com cofre aberto e foco no painel esquerdo (árvore). Nos modos Modelos e Configurações, `⌃F` e `F10` não têm efeito de busca. O campo de busca na linha separadora do cabeçalho **só aparece no Modo Cofre e apenas enquanto a busca estiver ativa** — nunca em outros modos, nunca na tela de boas-vindas.
**Modelo:** type-to-search — o campo na linha separadora do cabeçalho é display-only; o foco permanece na árvore durante toda a interação.

---

#### Ativação e saída

| Mecanismo | Efeito |
|---|---|
| `⌃F` ou `F10` com campo **fechado** | Campo abre na linha separadora; barra de mensagens exibe dica; barra de comandos muda para ações de busca |
| `⌃F` ou `F10` com campo **aberto** | Toggle: campo fecha; query descartada; árvore restaurada; barra restaurada ao estado anterior |
| `Esc` com campo aberto | Idêntico ao toggle com campo aberto; cursor retorna ao item que estava selecionado antes da busca |

> A busca **não pode ser ativada** com foco no painel direito (detalhe). O foco deve estar na árvore.

---

#### Mapa de teclas durante busca ativa

| Tecla | Efeito |
|---|---|
| Alfanumérica / símbolo imprimível | Acrescenta caractere à query; árvore filtra em tempo real |
| `Backspace` | Remove o último caractere da query |
| `Del` | Limpa toda a query de uma vez; campo permanece aberto e vazio; árvore restaurada completa |
| `↑` / `↓` | Navega entre os resultados visíveis na árvore filtrada |
| `Home` / `End` | Primeiro / último resultado visível |
| `PgUp` / `PgDn` | Scroll por página nos resultados |
| `Enter` com segredo selecionado | Abre detalhe no painel direito; campo permanece aberto |
| `Enter` com pasta selecionada | Expande / recolhe pasta; campo permanece aberto |
| `Tab` | Foco → painel direito (detalhe do item selecionado); campo permanece aberto e visível |
| `⌃F` / `F10` | Toggle: fecha o campo, descarta a query, restaura a árvore |
| `Esc` | Fecha o campo, descarta a query, restaura a árvore; cursor retorna ao item anterior |
| `F-keys` / `⌃Letra` | Ações normais da árvore (ActionManager) — **não alimentam a query** |

> **Regra de roteamento:** apenas teclas que produzem caracteres imprimíveis (Unicode printable) e `Backspace` são interceptadas pela busca enquanto o campo estiver aberto. Modificadores, F-keys e teclas de controle passam normalmente ao ActionManager.

---

#### Comportamento do filtro

- **Correspondência:** substring, case-insensitive, ignorando acentuação — conforme requisito funcional
- **Escopo da busca:** nome do segredo, nome de campo, valor de campo **comum**, observação
- **Excluído da busca:** valores de campos sensíveis (nomes de campos sensíveis participam normalmente)
- **Excluídos dos resultados:** segredos marcados para exclusão (`✗`)
- **Árvore compacta:** apenas pastas que contêm ≥ 1 resultado são exibidas; pastas sem resultados desaparecem completamente
- **Contadores de pasta durante filtro ativo:** formato `(N/Total)` — `N` = segredos que atendem à busca nessa pasta e subpastas; `Total` = total de segredos ativos nessa pasta e subpastas. Exemplo: `(2/6)` significa que 2 dos 6 segredos atendem à query. Quando `N = Total`, o contador volta ao formato simples `(N)` — sem barra. O formato `(N/Total)` só aparece durante busca ativa com query não vazia
- **Indicador visual de filtro ativo:** o painel esquerdo exibe `Filtro ativo` em `semantic.warning` + *italic*, alinhado à direita na primeira linha da área de trabalho, quando há query não vazia. Garante percepção do filtro mesmo que o cabeçalho esteja fora da viewport ou o foco esteja no painel direito
- **Match highlight:** o trecho de texto correspondente à query é exibido em `special.match` + **bold**
- **Query vazia:** campo aberto sem texto — árvore exibe tudo; contadores voltam ao formato `(N)`; indicador `Filtro ativo` não aparece
- **Persistência:** ao fechar o campo, a query é descartada e a árvore restaurada completa; o campo sempre abre vazio

---

#### Wireframes

**Campo aberto, sem query (recém-ativado):**

```
  Abditum · cofre •                      ╭───────╮  ╭ Modelos ╮  ╭ Config ╮
 ─ Busca: ────────────────────────────────╯ Cofre ╰──────────────────────────
  ▼ Favoritos          (2)  │
    ★ Bradesco         <╡
    ★ Gmail                 │
  ▼ Geral              (8)  │
    ▼ Sites            (5)  │
      ● Gmail               │
      ● YouTube             │
 ─ • Digite para filtrar os segredos ────────────────────────────────────────
  ⌃F Fechar · Del Limpar                                              F1 Ajuda
```

> Query vazia: árvore completa, contadores no formato `(N)`, sem indicador `Filtro ativo`.

**Campo aberto, com query — resultados encontrados:**

```
  Abditum · cofre •                      ╭───────╮  ╭ Modelos ╮  ╭ Config ╮
 ─ Busca: gmail ──────────────────────────╯ Cofre ╰──────────────────────────
  ▼ Favoritos        (1/2)  │              ← Filtro ativo
    ★ Gmail            <╡       ← match em special.match + bold
  ▼ Geral            (2/8)  │
    ▼ Sites          (2/5)  │
      ● Gmail               │
      ● Gmail Pro           │
 ─ ℹ 3 resultado(s) ─────────────────────────────────────────────────────────
  ⌃F Fechar · Del Limpar                                              F1 Ajuda
```

> `Filtro ativo` em `semantic.warning` + *italic*, alinhado à direita. `(1/2)` = 1 resultado dos 2 segredos em Favoritos. Quando `N = Total`, contador volta a `(N)`.

**Campo aberto, sem resultados:**

```
  Abditum · cofre •                      ╭───────╮  ╭ Modelos ╮  ╭ Config ╮
 ─ Busca: xyzxyz ─────────────────────────╯ Cofre ╰──────────────────────────
  ▷ Geral              (0)  │              ← Filtro ativo
                             │
                             │
 ─ ℹ Nenhum resultado ───────────────────────────────────────────────────────
  ⌃F Fechar · Del Limpar                                              F1 Ajuda
```

> Pasta raiz sempre visível, mesmo sem resultados. Indicador `Filtro ativo` permanece.

**Campo aberto, query longa (truncada à esquerda):**

```
  Abditum · cofre •                      ╭───────╮  ╭ Modelos ╮  ╭ Config ╮
 ─ Busca: …ail.google.com/conta ──────────╯ Cofre ╰──────────────────────────
```

> A parte mais recente da query (direita) fica sempre visível. `…` substitui os caracteres iniciais quando a query excede o espaço disponível.

---

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| `─ Busca: ` rótulo na linha separadora | `border.default` | — |
| Texto da query | `accent.primary` | **bold** |
| `─` preenchimento na linha separadora | `border.default` | — |
| Trecho de match na árvore | `special.match` | **bold** |
| Contador `(N/Total)` durante filtro ativo | `text.secondary` | — |
| Indicador `Filtro ativo` | `semantic.warning` | *italic* |

---

#### Estados dos componentes

| Componente | Estado | Condição |
|---|---|---|
| Campo de busca na linha separadora | oculto | Campo fechado — linha separadora normal |
| Campo de busca na linha separadora | visível, vazio | Campo aberto, query vazia |
| Campo de busca na linha separadora | visível, com texto | Query ativa (≥ 1 caractere) |
| Campo de busca na linha separadora | **nunca visível** fora do Modo Cofre | Modos Modelos, Configurações, Boas-vindas |
| Árvore | completa | Campo fechado **ou** campo aberto com query vazia |
| Árvore | filtrada (compacta) | Campo aberto com query ≥ 1 caractere |
| Pasta | visível | Contém ≥ 1 resultado direto ou indireto |
| Pasta | oculta | Não contém nenhum resultado |
| Pasta raiz | sempre visível | Mesmo sem resultados — exibe `(0)` e `▷` |
| Contador de pasta | formato `(N)` | Campo fechado **ou** query vazia **ou** `N = Total` |
| Contador de pasta | formato `(N/Total)` | Query ativa com ≥ 1 caractere e `N < Total` |
| Indicador `Filtro ativo` | visível, 1ª linha da área de trabalho, alinhado à direita | Query ativa com ≥ 1 caractere |
| Indicador `Filtro ativo` | oculto | Campo fechado ou query vazia |
| Trecho de match | `special.match` + **bold** | Substring correspondente à query |
| Barra de comandos | ações de busca (`⌃F Fechar · Del Limpar`) | Campo aberto |
| Barra de comandos | ações normais da árvore | Campo fechado |

---

#### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Campo abre (query vazia) | Dica de uso | `• Digite para filtrar os segredos` |
| Query ativa, com resultados | Informação | `ℹ N resultado(s)` |
| Query ativa, sem resultados | Informação | `ℹ Nenhum resultado` |
| `Backspace` apaga último caractere — query fica vazia | Dica de uso | `• Digite para filtrar os segredos` |
| `Del` limpa a query | Dica de uso | `• Digite para filtrar os segredos` |
| Campo fecha (`Esc`, `⌃F`, `F10`) | — | Barra restaurada ao estado anterior à busca |

---

#### Barra de comandos durante busca ativa

```
  ⌃F Fechar · Del Limpar                                              F1 Ajuda
```

As ações normais da árvore (ActionManager) ficam ocultas na barra enquanto o campo estiver aberto — o ActionManager continua processando suas teclas (`⌃Letra`, `F-keys`), mas a barra reflete apenas o contexto de busca.

---

#### Transições especiais

| Evento | Efeito |
|---|---|
| `⌃F` / `F10` — campo fechado | Campo abre; separadora substituída; barra muda; dica exibida |
| `⌃F` / `F10` — campo aberto | Campo fecha; query descartada; separadora restaurada; cursor volta ao item anterior; barra restaurada |
| `Esc` — campo aberto | Idêntico ao toggle com campo aberto |
| Digitação — query não vazia | Árvore filtra em tempo real; `ℹ N resultado(s)` atualiza a cada caractere |
| `Backspace` — query vazia após apagar | Árvore restaurada completa; campo permanece aberto; dica exibida |
| `Del` | Query limpa instantaneamente; campo permanece aberto; árvore restaurada; dica exibida |
| `Enter` — segredo selecionado | Detalhe atualizado no painel direito; campo permanece aberto |
| `Enter` — pasta selecionada | Pasta expande / recolhe; campo permanece aberto |
| `Tab` — foco na árvore | Foco vai para painel direito; campo permanece aberto e visível; type-to-search suspende até foco retornar à árvore |
| Foco retorna à árvore (`Tab` / clique) | Type-to-search retoma — teclas alfanuméricas voltam a alimentar a query |
| Terminal redimensionado | Largura disponível da query recalculada; truncamento com `…` reaplicado se necessário |

---


## Ações na Árvore de Segredos

Esta seção detalha as ações disponíveis ao interagir com a árvore de segredos (painel esquerdo do Modo Cofre) e seus respectivos atalhos de teclado. As regras gerais de navegação e atribuição de teclas são definidas no [Design System — Mapa de Teclas](tui-design-system.md#mapa-de-teclas).

### Navegação na árvore (geral)

| Tecla           | Ação                                     | Notas                                            |
|-----------------|------------------------------------------|--------------------------------------------------|
| `↑` / `↓`       | Mover cursor na lista / árvore           |                                                  |
| `Home` / `End`  | Mover para o primeiro / último item visível |                                                  |
| `PgUp` / `PgDn` | Rolar uma página para cima / baixo       |                                                  |
| `Tab`           | Alternar foco entre painéis              | Move o foco para o painel direito (Detalhe) e vice-versa. |

### Ações em pastas

| Tecla           | Ação                                     | Notas                                                                      |
|-----------------|------------------------------------------|----------------------------------------------------------------------------|
| `→`             | Expandir pasta                           |                                                                            |
| `←`             | Recolher pasta                           |                                                                            |
| `Enter`         | Expandir / Recolher pasta                | Quando o foco está em uma pasta, expande/contrai.                          |
| `Shift+Insert`  | Criar nova pasta                         | Cria uma nova pasta no mesmo nível da pasta focada ou dentro dela, se não houver nenhuma pasta focada. |
| `Ctrl+Shift+I`  | Criar nova pasta                         | Atalho alternativo para criar uma nova pasta.                              |
| `Delete`        | Remover pasta                            | Marca a pasta selecionada para remoção (reversível até o salvamento).      |

### Ações em segredos

A coluna **Favoritos** indica se a ação está disponível quando o cursor está na pasta virtual Favoritos. Ações indisponíveis ficam ocultas na barra de comandos — não aparecem desabilitadas.

| Tecla    | Ação                          | Favoritos | Notas                                                                      |
|----------|-------------------------------|-----------|----------------------------------------------------------------------------|
| `Enter`  | Focar no painel de detalhes   | ✓         | Comporta-se de forma similar ao `Tab` quando o foco está em um segredo.    |
| `Insert` | Novo segredo                  | —         | Indisponível: Favoritos é somente leitura, sem pasta real associada.       |
| `⌃I`     | Novo segredo                  | —         | Atalho alternativo — mesma restrição.                                      |
| `⌃E`     | Editar segredo                | ✓         | Opera no segredo real, independente da visão atual.                        |
| `⌃D`     | Duplicar segredo              | —         | Indisponível: destino ambíguo. Navegar até a pasta real para duplicar.     |
| `⌃M`     | Mover para outra pasta        | —         | Indisponível: mover a partir de pasta somente leitura não é permitido.     |
| `!↑`     | Mover para cima na lista      | —         | Indisponível: a ordem na Favoritos reflete a árvore real.                  |
| `!↓`     | Mover para baixo na lista     | —         | Indisponível: idem.                                                        |
| `⌃S`     | Desfavoritar segredo          | ✓ (só ⊖)  | Na Favoritos, o toggle só remove o favorito — o segredo some da lista imediatamente. Em pasta real, alterna entre favoritar e desfavoritar normalmente. |
| `⌃R`     | Revelar primeiro campo sensível | ✓       | Visível apenas se o segredo tiver pelo menos um campo sensível.            |
| `⌃C`     | Copiar primeiro campo sensível  | ✓       | Visível apenas se o segredo tiver pelo menos um campo sensível.            |
| `Delete` | Excluir segredo               | —         | Indisponível: exclusão direta a partir de pasta somente leitura não é permitida. |

#### ⌃D — Duplicar segredo

**Contexto:** foco na árvore com cursor em um segredo, em pasta real. Indisponível na pasta virtual Favoritos — o destino do duplicado seria ambíguo para o usuário; a operação deve ser realizada navegando até a pasta real do segredo.

**Comportamento:**
- Cria uma cópia do segredo com todos os campos, valores e histórico de modelo idênticos ao original.
- O novo segredo recebe automaticamente um nome único na mesma pasta com sufixo numérico — ex: `Gmail (1)` se `Gmail` já existe; `Gmail (2)` se `Gmail (1)` também já existe.
- O novo segredo é posicionado imediatamente **após o original** na lista da pasta.
- O novo segredo entra em estado `incluido`.
- O cursor da árvore permanece no segredo original após a duplicação — o usuário pode navegar para o novo com `↓`.
- A operação é instantânea, sem diálogo de confirmação.

**Feedback:** barra de mensagens exibe `✓ "[Nome original]" duplicado como "[Novo nome]"`.

**Referência:** [Fluxo 19 — Duplicar segredo](fluxos.md#fluxo-19--duplicar-segredo)

---

#### ⌃M — Mover para outra pasta

**Contexto:** foco na árvore com cursor em um segredo. Não disponível na pasta virtual Favoritos (a pasta Favoritos é somente leitura — mover deve ocorrer na pasta real).

**Modo de seleção inline:**

A árvore entra em **modo mover** — um estado visual distinto:
- O segredo em foco recebe um indicador de "em movimento" (ex: ícone `↷` ou destaque diferenciado em `accent.secondary`) e o cursor passa a navegar pela estrutura de pastas como destino.
- A barra de mensagens exibe `• Navegue até a pasta de destino e pressione Enter para confirmar`.
- A barra de comandos muda para: `Enter Mover aqui · Esc Cancelar`.
- O usuário navega com `↑↓←→` entre as pastas visíveis.
- Pastas que resultariam em conflito de nome (já contêm um segredo com o mesmo nome) são marcadas visualmente como inválidas — o cursor pode passar por elas, mas `Enter` sobre elas exibe mensagem de erro na barra e aguarda nova seleção.
- `Enter` sobre uma pasta válida confirma o movimento; o segredo é movido para a pasta de destino, o modo mover é encerrado e o cursor acompanha o segredo para a nova posição.
- `Esc` cancela o modo mover sem efeito colateral; o cursor retorna ao segredo original.

**Referência:** [Fluxo 25 — Mover segredo para outra pasta](fluxos.md#fluxo-25--mover-segredo-para-outra-pasta)

---

#### !↑ / !↓ — Reordenar segredo na lista

**Contexto:** foco na árvore com cursor em um segredo, dentro de uma pasta real (não Favoritos).

**Comportamento:**
- `!↑` desloca o segredo uma posição acima na lista da pasta atual; `!↓` desloca uma posição abaixo.
- A operação é instantânea e pode ser repetida sucessivamente.
- O cursor acompanha o segredo — após o deslocamento, o cursor permanece sobre o mesmo segredo na nova posição.
- Múltiplos deslocamentos antes de salvar resultam apenas no estado final — o histórico de movimentos intermediários é descartado.
- A operação não tem feedback de mensagem — o deslocamento visual imediato é o feedback.

**Limites:**
- `!↑` não tem efeito quando o segredo já está na primeira posição da pasta.
- `!↓` não tem efeito quando o segredo já está na última posição da pasta.
- Ambos ficam **ocultos na barra de comandos** e inativos quando o cursor está na pasta virtual Favoritos.

**Indicador de modo ativo:** a barra de status/cabeçalho não precisa de indicador de modo para reordenação direta — a operação é pontual e sem estado persistente.

**Referência:** [Fluxo 26 — Reordenar segredo dentro da mesma pasta](fluxos.md#fluxo-26--reordenar-segredo-dentro-da-mesma-pasta)

---

#### Barra de comandos contextualizada (árvore, cursor em segredo — completa)

A tabela abaixo consolida todas as variações da barra de comandos para segredos na árvore, incluindo os atalhos anteriores (`⌃R`, `⌃C`) e os novos (`⌃D`, `⌃M`, `!↑`, `!↓`).

| Condição | Barra de comandos |
|---|---|
| Pasta real — segredo sem campo sensível — posição intermediária | `Enter Detalhes · ⌃E Editar · ⌃D Duplicar · ⌃M Mover · !↑ !↓ Reordenar · ⌃S Favoritar · Del Excluir · F1 Ajuda` |
| Pasta real — segredo sem campo sensível — primeiro da lista | `Enter Detalhes · ⌃E Editar · ⌃D Duplicar · ⌃M Mover · !↓ Mover para baixo · ⌃S Favoritar · Del Excluir · F1 Ajuda` |
| Pasta real — segredo sem campo sensível — último da lista | `Enter Detalhes · ⌃E Editar · ⌃D Duplicar · ⌃M Mover · !↑ Mover para cima · ⌃S Favoritar · Del Excluir · F1 Ajuda` |
| Pasta real — segredo com campo sensível — reveal mascarado | `Enter Detalhes · ⌃E Editar · ⌃D Duplicar · ⌃M Mover · !↑ !↓ Reordenar · ⌃S Favoritar · ⌃R Revelar · ⌃C Copiar · Del Excluir · F1 Ajuda` |
| Pasta real — segredo com campo sensível — reveal com dica | `Enter Detalhes · ⌃E Editar · ⌃D Duplicar · ⌃M Mover · !↑ !↓ Reordenar · ⌃S Favoritar · ⌃R Mostrar tudo · ⌃C Copiar · Del Excluir · F1 Ajuda` |
| Pasta real — segredo com campo sensível — reveal completo | `Enter Detalhes · ⌃E Editar · ⌃D Duplicar · ⌃M Mover · !↑ !↓ Reordenar · ⌃S Favoritar · ⌃R Ocultar · ⌃C Copiar · Del Excluir · F1 Ajuda` |
| Pasta Favoritos — segredo sem campo sensível | `Enter Detalhes · ⌃E Editar · ⌃S Desfavoritar · F1 Ajuda` |
| Pasta Favoritos — segredo com campo sensível — reveal mascarado | `Enter Detalhes · ⌃E Editar · ⌃S Desfavoritar · ⌃R Revelar · ⌃C Copiar · F1 Ajuda` |
| Pasta Favoritos — segredo com campo sensível — reveal com dica | `Enter Detalhes · ⌃E Editar · ⌃S Desfavoritar · ⌃R Mostrar tudo · ⌃C Copiar · F1 Ajuda` |
| Pasta Favoritos — segredo com campo sensível — reveal completo | `Enter Detalhes · ⌃E Editar · ⌃S Desfavoritar · ⌃R Ocultar · ⌃C Copiar · F1 Ajuda` |
| Modo mover ativo (⌃M pressionado) | `Enter Mover aqui · Esc Cancelar` |

> **Nota sobre tamanho da barra:** as entradas acima são o conjunto completo de ações disponíveis. Em terminais estreitos, a barra de comandos trunca à direita — as ações mais prioritárias devem aparecer primeiro. A ordem na barra segue a frequência de uso esperada.

---

#### ⌃R e ⌃C na árvore — Atalhos de campo sensível

**Contexto:** foco na árvore com cursor em um segredo que possui pelo menos um campo sensível.

**Campo alvo:** sempre o **primeiro campo sensível** do segredo (menor índice de posição no tipo).

**Visibilidade dos atalhos:**
- `⌃R` e `⌃C` aparecem na barra de comandos **somente** quando o cursor da árvore está em um segredo com pelo menos um campo sensível.
- Quando o cursor está em uma pasta ou em um segredo sem campos sensíveis, os atalhos são omitidos da barra e não têm efeito.

##### Comportamento de ⌃R na árvore

- `⌃R` cicla o estado de reveal do primeiro campo sensível usando o **mesmo mecanismo de 3 estados do painel de detalhe**: mascarado → dica (3 primeiros chars + `••`) → completo → mascarado.
- O painel direito é aberto (ou atualizado) automaticamente exibindo o segredo com o campo sensível já no estado correspondente ao toque atual:
  - **1º toque:** painel exibe o campo sensível em estado de dica.
  - **2º toque:** painel exibe o campo sensível revelado completamente.
  - **3º toque:** campo re-mascarado; painel permanece aberto.
- As mesmas regras de re-mascaramento do painel se aplicam: trocar de segredo na árvore ou timeout expirado re-mascara o campo silenciosamente.
- A barra de comandos reflete o estado atual do reveal (igual ao painel):
  - Mascarado: `⌃R Revelar`
  - Dica ativa: `⌃R Mostrar tudo`
  - Revelado: `⌃R Ocultar`

##### Comportamento de ⌃C na árvore

- `⌃C` copia o valor **completo** do primeiro campo sensível para a clipboard — independentemente do estado de reveal atual (não é necessário revelar antes de copiar).
- Agenda limpeza automática da clipboard (mesmo comportamento do `⌃C` no painel de detalhe).
- O painel direito é aberto (ou atualizado) automaticamente exibindo o segredo, mas o estado de reveal do campo **não muda** — a cópia não desencadeia reveal.
- A barra de mensagens exibe confirmação: `✓ [Rótulo do campo] copiado para a área de transferência`.

---

