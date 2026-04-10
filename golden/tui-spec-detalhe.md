# Especificação Visual — Detalhe do Segredo

> Painel direito: modos Leitura, Edição de Valores e Edição de Estrutura.
> Parte da [Especificação Visual](tui-specification.md).
>
> **Documentos relacionados:**
> - [`tui-design-system.md`](tui-design-system.md) — fundações visuais
> - [`tui-spec-arvore.md`](tui-spec-arvore.md) — painel esquerdo (árvore de segredos)

### Painel Direito: Detalhe do Segredo — Modo Leitura

**Contexto:** Área de trabalho — Modo Cofre.
**Largura:** ~65% da área de trabalho.
**Responsabilidade:** Exibir o nome, o caminho de pastas, os campos e a observação do segredo selecionado na árvore; permitir navegação entre campos, cópia de valores e reveal de campos sensíveis.

> Este documento especifica apenas o **modo leitura**. O modo edição de valores é especificado em [Modo Edição de Valores](#painel-direito-detalhe-do-segredo--modo-edição-de-valores). O modo edição de estrutura é especificado em [Modo Edição de Estrutura](#painel-direito-detalhe-do-segredo--modo-edição-de-estrutura).

---

#### Anatomia do painel

```
  Nome do Segredo                          Geral › Sites › Gmail ↑
  ──────────────────────────────────────────────────────────────  │
  Rótulo do campo 1                                               ■
  Valor do campo 1                                                │
                                                                  │
  Rótulo do campo 2                                               │
  Valor do campo 2                                                │
                                                                  │
  ╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌  ↓
  Texto da observação...
```

**Linha 1 — cabeçalho do segredo:**
- Esquerda: nome do segredo em `text.primary` **bold**
- Direita: breadcrumb com caminho completo de pastas em `text.secondary` — formato `Pasta › Subpasta › ...`; truncado à esquerda com `…` se não couber na linha. `★` aparece entre o nome e o breadcrumb quando o segredo é favoritado, em `accent.secondary`
- O breadcrumb mostra o caminho até o segredo, excluindo o nome do segredo

**Linha 2 — separador:**
- `─` em `border.default` por toda a largura do painel (exceto a coluna reservada ao scroll)

**Área de campos:**
- Cada campo ocupa dois segmentos: **rótulo** (linha própria, `text.secondary`) e **valor** (linha(s) seguinte(s), `text.primary`)
- Uma linha em branco separa campos consecutivos
- Campos sensíveis exibem o valor mascarado com `••••••••` em `text.secondary`; ao serem revelados, o valor real aparece em `text.primary`
- Campos com valor vazio: o rótulo é exibido normalmente, a linha do valor fica em branco

**Separador da Observação:**
- `╌` (U+254C) em `border.default`, ocupando toda a largura — omitido quando a Observação está vazia
- A Observação não tem rótulo; o separador e a posição final comunicam o que é

**Trilha de scroll:**
- Última coluna do painel reservada para `↑`/`↓`/`■` em `text.secondary`
- Reservada mesmo quando não há scroll (evita deslocamento de conteúdo ao ativar)

---

#### Wireframes

**Painel sem foco — segredo com campos variados:**

```
  Gmail ★                              Geral › Sites e Apps
  ──────────────────────────────────────────────────────────
  URL
  https://accounts.google.com/login

  Usuário
  fulano@gmail.com

  Senha
  ••••••••••

  Token 2FA
  ••••••••

  ╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌
  Conta pessoal principal — criada em 2018.
```

> Sem foco: nenhum campo destacado. O `★` aparece entre o nome e o breadcrumb quando o segredo é favoritado.

**Painel com foco — cursor em campo comum:**

```
  Gmail ★                              Geral › Sites e Apps
  ──────────────────────────────────────────────────────────
  URL
  https://accounts.google.com/login

  Usuário                                                     ← special.highlight no bloco
  fulano@gmail.com

  Senha
  ••••••••••

  ╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌
  Conta pessoal principal.
```

> O bloco inteiro do campo em foco (rótulo + valor + linha em branco) recebe `special.highlight`. Barra de comandos (campo comum): `Enter Editar · ⌃S Favoritar · ⌃C Copiar · Tab Árvore · F1 Ajuda`

**Painel com foco — cursor em campo sensível:**

```
  Gmail ★                              Geral › Sites e Apps
  ──────────────────────────────────────────────────────────
  URL
  https://accounts.google.com/login

  Usuário
  fulano@gmail.com

  Senha                                                       ← special.highlight no bloco
  ••••••••••

  ╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌
  Conta pessoal principal.
```

> Barra de comandos (campo sensível mascarado): `Enter Editar · ⌃S Favoritar · ⌃C Copiar · ⌃R Revelar · Tab Árvore · F1 Ajuda`

**Campo sensível — estado de dica (1º `⌃R`):**

```
  Gmail ★                              Geral › Sites e Apps
  ──────────────────────────────────────────────────────────
  URL
  https://accounts.google.com/login

  Usuário
  fulano@gmail.com

  Senha                                                       ← special.highlight
  min•••••••••••••                                            ← 3 chars revelados + •• mascarados

  ╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌
  Conta pessoal principal.
```

> Barra de comandos (dica ativa): `Enter Editar · ⌃S Favoritar · ⌃C Copiar · ⌃R Mostrar tudo · Tab Árvore · F1 Ajuda`

**Campo sensível — revelado completamente (2º `⌃R`):**

```
  Gmail ★                              Geral › Sites e Apps
  ──────────────────────────────────────────────────────────
  URL
  https://accounts.google.com/login

  Usuário
  fulano@gmail.com

  Senha                                                       ← special.highlight
  minha-senha-secreta-123

  ╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌
  Conta pessoal principal.
```

> Barra de comandos (revelado): `Enter Editar · ⌃S Favoritar · ⌃C Copiar · ⌃R Ocultar · Tab Árvore · F1 Ajuda`

**Scroll ativo:**

```
  Gmail ★                              Geral › Sites e Apps ↑
  ──────────────────────────────────────────────────────────  │
  URL                                                         ■
  https://accounts.google.com/login/v2/identifier?hl=pt-BR   │
                                                              │
  Usuário                                                     │
  fulano@gmail.com                                            │
                                                              ↓
```

> Trilha de scroll: `↑` quando há conteúdo acima, `↓` quando há abaixo, `■` na posição proporcional do thumb. A coluna da trilha é sempre reservada — o conteúdo não se desloca ao ativar o scroll.

**Valor longo com quebra de linha:**

```
  Passos de acesso
  1. Acesse https://accounts.google.com
  2. Clique em "Fazer login com o Google"
  3. Confirme o dispositivo no app

```

> Valores multilinha recebem word-wrap; cada linha do valor ocupa a largura disponível (exceto a coluna do scroll). O campo continua sendo tratado como uma unidade de foco — o bloco inteiro recebe highlight.

**Placeholders:**

```
  (sem segredo selecionado)
  ─────────────────────────────────────────────────────────────────


               Selecione um segredo para ver os detalhes


```

```
  (cofre vazio)
  ─────────────────────────────────────────────────────────────────


                           Cofre vazio


```

> Textos em `text.secondary` *italic*, centralizados na área de conteúdo.

**Segredo sem Observação (separador omitido):**

```
  API Key — Stripe                            Geral › Financeiro
  ──────────────────────────────────────────────────────────────
  Serviço
  Stripe

  Chave
  ••••••••••

```

> Quando a Observação está vazia, o separador `╌╌╌` é omitido. Não há linha em branco extra no final.

**Breadcrumb truncado (caminho longo):**

```
  Gmail ★          … › Projetos › Cliente ABC › Acessos › Gmail
  ──────────────────────────────────────────────────────────────
```

> O breadcrumb é truncado à esquerda com `…` quando o caminho completo não cabe. O nome do segredo e o `★` nunca são truncados.

---

#### Mapa de teclas

| Tecla | Efeito | Condição |
|---|---|---|
| `↑` / `↓` | Move cursor para o campo anterior / próximo | Painel com foco |
| `Home` | Vai ao primeiro campo | Painel com foco |
| `End` | Vai ao último campo (Observação, se não vazia) | Painel com foco |
| `PgUp` / `PgDn` | Scroll por página (viewport − 1 linhas) | Painel com foco |
| `Enter` | Entra no modo edição do campo em foco | Painel com foco |
| `⌃S` | Favoritar / Desfavoritar segredo | Painel com foco |
| `⌃R` | 1º toque: revela dica (3 primeiros chars); 2º toque: revela valor completo; 3º toque: re-mascara | Painel com foco; campo sensível em foco |
| `⌃C` | Copiar valor do campo para clipboard; agenda limpeza da clipboard se campo sensível | Painel com foco; qualquer campo |
| `Tab` | Foco → painel esquerdo (árvore) | Painel com foco |

> `⌃R` não tem efeito quando o campo em foco é comum — a barra de comandos omite a ação `Revelar` nesses casos.

---

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Nome do segredo | `text.primary` | **bold** |
| `★` favorito | `accent.secondary` | — |
| Breadcrumb de pasta | `text.secondary` | — |
| Separador `───` cabeçalho | `border.default` | — |
| Rótulo de campo | `text.secondary` | **bold** |
| Valor de campo comum | `text.primary` | — |
| Valor de campo — URL | `text.link` | — |
| Valor de campo sensível — mascarado `••••••••` | `text.secondary` | — |
| Valor de campo sensível — dica (`min••••`) | `text.secondary` | — |
| Fundo do campo em foco | `special.highlight` | — |
| Separador `╌╌╌` da Observação | `border.default` | — |
| Texto da Observação | `text.primary` | — |
| Placeholders | `text.secondary` | *italic* |
| `│` separador vertical — painel com foco | `border.focused` | — |
| `│` separador vertical — painel sem foco | `border.default` | — |
| `↑`/`↓`/`■` trilha de scroll | `text.secondary` | — |

---

#### Estados dos componentes

| Componente | Estado | Condição |
|---|---|---|
| Painel | placeholder "Selecione…" | Cofre tem segredos; nenhum segredo foi selecionado ainda na sessão |
| Painel | placeholder "Cofre vazio" | Cofre sem nenhum segredo |
| Painel | segredo exibido (último selecionado) | Cursor da árvore em pasta — painel mantém o último segredo exibido |
| Painel | segredo exibido (atual) | Cursor da árvore em segredo |
| Cursor de campo | ausente | Painel sem foco |
| Cursor de campo | `special.highlight` no bloco do campo | Painel com foco |
| `★` | visível no cabeçalho, entre nome e breadcrumb | Segredo favoritado |
| `★` | ausente | Segredo não favoritado |
| Campo sensível | mascarado `••••••••` | Estado inicial ao exibir qualquer segredo |
| Campo sensível | dica (3 primeiros chars + `••`) | 1º `⌃R`; campo ainda em foco; timeout não expirou |
| Campo sensível | revelado (valor completo) | 2º `⌃R`; campo ainda em foco; timeout não expirou |
| Campo sensível revelado | re-mascarado | Timeout expirou; segredo diferente selecionado; foco saiu do campo |
| Separador `╌╌╌` | visível | Observação não vazia |
| Separador `╌╌╌` | omitido | Observação vazia |
| Trilha de scroll | `↑`/`↓`/`■` ativos | Conteúdo excede a área visível |
| Trilha de scroll | coluna reservada, vazia | Conteúdo cabe na área visível |

---

#### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Painel recebe foco | Dica | `• Navegue com ↑↓ e copie com ⌃C` |
| Campo sensível selecionado | Dica | `• ⌃R Revelar · ⌃C Copiar` |
| `⌃C` copia valor | Sucesso (5s) | `✓ [Rótulo do campo] copiado para a área de transferência` |

---

#### Eventos

| Evento | Efeito |
|---|---|
| Segredo selecionado na árvore | Conteúdo atualizado; campos revelados re-mascarados; cursor vai ao primeiro campo; `<╡` aparece no separador |
| Painel recebe foco (`Tab`) | Cursor de campo aparece no campo anteriormente ativo, ou no primeiro campo se nunca focado |
| `↑` / `↓` | Cursor move para o campo anterior / próximo; scroll automático se necessário |
| `Home` / `End` | Cursor vai ao primeiro / último campo; scroll automático |
| `PgUp` / `PgDn` | Scroll por página |
| `Enter` | Entra no modo edição do campo em foco |
| `⌃S` | Segredo favoritado → desfavoritado (ou vice-versa); `★` no cabeçalho do painel atualiza imediatamente; árvore atualiza em segundo plano |
| `⌃R` em campo sensível mascarado | Campo entra em estado de dica (3 primeiros chars); barra muda para `⌃R Mostrar tudo · ⌃R Ocultar` |
| `⌃R` em campo sensível com dica | Campo revelado completamente; barra muda para `⌃R Ocultar` |
| `⌃R` em campo sensível revelado | Campo re-mascarado; barra volta para `⌃R Revelar` |
| `↑` / `↓` saindo de campo sensível revelado | Campo re-mascarado silenciosamente antes de mover o cursor |
| `Tab` com campo sensível revelado | Campo re-mascarado silenciosamente; foco transferido para a árvore |
| Timeout de reveal expira | Campo re-mascarado silenciosamente; sem mensagem |
| Segredo diferente selecionado | Todos os campos revelados re-mascarados; cursor vai ao primeiro campo |

---

#### Comportamento

- **Cursor somente com foco** — o cursor de campo (highlight no bloco) aparece apenas quando o painel tem foco; sem foco, o conteúdo é exibido sem destaque
- **Bloco de campo** — o campo em foco compreende: linha do rótulo + linha(s) do valor + linha em branco de separação; todo o bloco recebe `special.highlight`
- **`Enter` entra no modo edição** — disponível em qualquer campo com foco; aciona o modo edição de valores (especificado separadamente)
- **`⌃R` contextual** — disponível apenas com campo sensível em foco; cicla entre três estados: mascarado → dica (3 primeiros chars) → completo → mascarado. Não aparece na barra quando o campo em foco é comum
- **Re-mascaramento ao sair do campo** — ao mover o cursor para outro campo (`↑`/`↓`/`Home`/`End`) ou ao transferir o foco para a árvore (`Tab`), qualquer campo sensível que estiver em estado de dica ou revelado é re-mascarado silenciosamente antes da movimentação
- **Campos sensíveis sempre iniciam mascarados** — incluindo segredos já visitados anteriormente na sessão
- **Reveal timeout** — configurável nas Configurações; ao expirar, o campo é re-mascarado silenciosamente (sem mensagem na barra). Ao trocar de segredo, todos os reveals são cancelados imediatamente
- **URLs** — valores identificados como URL recebem `text.link`, diferenciados visualmente de texto puro
- **Observação — word-wrap** — o texto da Observação quebra na largura disponível (exceto a coluna do scroll); pode ocupar múltiplas linhas; o painel inteiro é scrollável
- **Scroll** — a última coluna do painel é sempre reservada para a trilha de scroll, mesmo quando não há overflow — o conteúdo não se desloca ao ativar o scroll (ver [DS — Scroll em diálogos](tui-design-system.md#scroll-em-diálogos))
- **`<╡` e trilha de scroll são independentes** — `<╡` aparece no separador vertical esquerdo e indica qual item da árvore está sendo detalhado; a trilha de scroll aparece na margem direita e reflete o scroll do conteúdo do painel. Um não afeta o outro
- **Posição do cursor ao retornar o foco** — ao receber foco via `Tab` novamente, o cursor vai ao campo que estava ativo antes de o foco sair; se nunca focado, vai ao primeiro campo
- **Breadcrumb — truncamento** — o breadcrumb é truncado à esquerda com `…` se o caminho completo não couber; o nome do segredo e o `★` nunca são truncados

---

### Painel Direito: Detalhe do Segredo — Modo Edição de Valores

**Contexto:** Área de trabalho — Modo Cofre. Ativado quando o usuário pressiona `Enter` sobre um campo no painel de detalhe em Modo Leitura.
**Largura:** ~65% da área de trabalho (igual ao Modo Leitura).
**Responsabilidade:** Permitir editar o valor de cada campo do segredo individualmente, com persistência imediata por campo, sem estado global pendente.

> O modo edição de estrutura (renomear campos, adicionar/remover campos, reordenar) é especificado em [Modo Edição de Estrutura](#painel-direito-detalhe-do-segredo--modo-edição-de-estrutura).

---

#### Anatomia do modo

O Modo Edição de Valores é uma camada sobre o Modo Leitura. O layout do painel (cabeçalho, separador, campos, observação, scroll) permanece o mesmo — o que muda são:

1. **Indicador de modo** — `[editando]` em `accent.primary` **bold** aparece no cabeçalho, após o nome do segredo e antes do `★`/breadcrumb
2. **Cursor de campo** — continua sendo `special.highlight` no bloco, como no Modo Leitura; o input se abre sobre o campo em foco
3. **Input inline** — quando um campo está em edição, o valor é substituído por um campo de texto editável na mesma posição; o input ocupa a largura total do painel (exceto a coluna de scroll)
4. **Barra de comandos** — muda conforme o estado: cursor de campo sem input aberto, ou input aberto

---

#### Anatomia do cabeçalho em edição

```
  Gmail [editando] ★                     Geral › Sites e Apps
  ──────────────────────────────────────────────────────────
```

- Nome do segredo: `text.primary` **bold** (igual ao Modo Leitura)
- `[editando]`: `accent.primary` **bold**, separado do nome por um espaço
- `★` e breadcrumb: inalterados

---

#### Wireframes

**Cursor no campo, sem input aberto (campo comum):**

```
  Gmail [editando] ★               Geral › Sites e Apps
  ──────────────────────────────────────────────────────
  URL

  Usuário                                                ← special.highlight no bloco
  fulano@gmail.com

  Senha
  ••••••••••

  ╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌
  Conta pessoal principal.
```

> Barra: `Enter Editar campo · ⌃N Renomear · ⌃S Favoritar · Tab Árvore · Esc Sair da edição · F1 Ajuda`

**Input aberto — campo comum:**

```
  Gmail [editando] ★               Geral › Sites e Apps
  ──────────────────────────────────────────────────────
  URL

  Usuário                                                ← special.highlight no bloco
  ░fulano@gmail.com▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░

  Senha
  ••••••••••

  ╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌
  Conta pessoal principal.
```

> `░` marca o fundo do input (`input.background`); `▌` é o cursor de texto. O input substitui visualmente a linha do valor; o rótulo permanece acima. Barra: `Enter Confirmar · Esc Cancelar campo · F1 Ajuda`

**Input aberto — campo sensível (revelado automaticamente):**

```
  Gmail [editando] ★               Geral › Sites e Apps
  ──────────────────────────────────────────────────────
  URL

  Usuário
  fulano@gmail.com

  Senha                                                  ← special.highlight no bloco
  ░minha-senha-secreta-123▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░

  ╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌
  Conta pessoal principal.
```

> Ao abrir o input de campo sensível, o valor é revelado automaticamente em texto claro dentro do input. Ao fechar o input (`Enter` ou `Esc`), o campo é re-mascarado imediatamente. Barra: `Enter Confirmar · Esc Cancelar campo · F1 Ajuda`

**Renomear segredo — input no cabeçalho (`⌃N`):**

```
  ░Gmail▌░░░░░░░░░░  [editando] ★        Geral › Sites e Apps
  ──────────────────────────────────────────────────────────
  URL
  https://accounts.google.com/login

  Usuário
  fulano@gmail.com

  Senha
  ••••••••••
```

> O input do nome abre inline no cabeçalho, substituindo o nome do segredo; `[editando]`, `★` e breadcrumb permanecem à direita. Nenhum campo da lista está em foco enquanto o input do nome está aberto. Barra: `Enter Confirmar nome · Esc Cancelar · F1 Ajuda`

**Validação — nome duplicado:**

```
  ░Gmail▌░░░░░░░░░░  [editando] ★        Geral › Sites e Apps
  ──────────────────────────────────────────────────────────
```

> Barra de mensagens (erro): `✗ Já existe um segredo com esse nome nesta pasta` — input permanece aberto; o valor não é persistido.

**Cursor no campo, sem input — campo sensível:**

```
  Gmail [editando] ★               Geral › Sites e Apps
  ──────────────────────────────────────────────────────
  URL

  Usuário
  fulano@gmail.com

  Senha                                                  ← special.highlight no bloco
  ••••••••••

  ╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌
  Conta pessoal principal.
```

> Campo sensível permanece mascarado enquanto não há input aberto. Barra: `Enter Editar campo · ⌃N Renomear · ⌃S Favoritar · Tab Árvore · Esc Sair da edição · F1 Ajuda`

---

#### Mapa de teclas

**Com cursor de campo, sem input aberto:**

| Tecla | Efeito | Condição |
|---|---|---|
| `↑` / `↓` | Move cursor para o campo anterior / próximo (sem abrir input) | — |
| `Home` / `End` | Cursor vai ao primeiro / último campo | — |
| `Enter` | Abre input inline no campo em foco | — |
| `⌃N` | Abre input inline no cabeçalho (renomear segredo) | — |
| `⌃S` | Favoritar / Desfavoritar segredo | — |
| `Tab` | Foco → árvore; sai do modo edição | — |
| `Esc` | Sai do modo edição; retorna ao Modo Leitura | — |

**Com input de campo aberto:**

| Tecla | Efeito |
|---|---|
| Texto / Backspace / Delete | Edita o valor no input |
| `Enter` | Persiste o valor; fecha o input; cursor permanece no campo |
| `↑` | Persiste o valor implicitamente; fecha o input; move cursor para o campo anterior |
| `↓` | Persiste o valor implicitamente; fecha o input; move cursor para o próximo campo |
| `Esc` | Cancela; restaura o valor anterior; fecha o input; cursor permanece no campo |

**Com input do nome aberto (`⌃N`):**

| Tecla | Efeito |
|---|---|
| Texto / Backspace / Delete | Edita o nome no input |
| `Enter` | Valida e persiste o nome; fecha o input; retorna ao cursor de campo |
| `Esc` | Cancela; restaura o nome anterior; fecha o input |

> `Tab` com input de campo aberto: persiste o valor implicitamente, fecha o input, foco vai para a árvore e sai do modo edição.
> `Tab` com input do nome aberto: cancela o nome (sem persistir), foco vai para a árvore e sai do modo edição.

---

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Nome do segredo (cabeçalho) | `text.primary` | **bold** |
| `[editando]` | `accent.primary` | **bold** |
| `★` favorito | `accent.secondary` | — |
| Breadcrumb de pasta | `text.secondary` | — |
| Fundo do campo em foco (sem input) | `special.highlight` | — |
| Fundo do input aberto | `input.background` | — |
| Texto dentro do input | `text.primary` | — |
| Cursor de texto no input | terminal padrão | — |
| Rótulo de campo | `text.secondary` | **bold** |
| Valor de campo comum (sem input) | `text.primary` | — |
| Valor de campo sensível mascarado (sem input) | `text.secondary` | — |
| Separador `───` cabeçalho | `border.default` | — |
| Separador `╌╌╌` da Observação | `border.default` | — |

---

#### Estados dos componentes

| Componente | Estado | Condição |
|---|---|---|
| Indicador `[editando]` | visível no cabeçalho | Modo edição de valores ativo |
| Cursor de campo | `special.highlight` no bloco | Sempre (modo edição tem foco implícito) |
| Input de campo | ausente | Cursor de campo sem edição ativa |
| Input de campo | aberto sobre a linha do valor | `Enter` pressionado sobre o campo |
| Campo sensível | mascarado `••••••••` | Input fechado |
| Campo sensível | revelado (texto claro no input) | Input aberto |
| Campo sensível | re-mascarado | Input fechado após `Enter` ou `Esc` |
| Input do nome | ausente | `⌃N` não pressionado |
| Input do nome | aberto no cabeçalho | `⌃N` pressionado |
| Cursor de campo da lista | ausente | Input do nome aberto |

---

#### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Modo edição ativado | Dica | `• Enter para editar um campo · Esc para sair` |
| Campo confirmado (`Enter` ou `↑`/`↓` implícito) | Sucesso (3s) | `✓ [Rótulo do campo] salvo` |
| Nome duplicado ao confirmar | Erro | `✗ Já existe um segredo com esse nome nesta pasta` |
| Campo confirmado — campo sensível | Sucesso (3s) | `✓ [Rótulo do campo] salvo` |

---

#### Eventos

| Evento | Efeito |
|---|---|
| `Enter` no Modo Leitura sobre um campo | Modo edição de valores ativado; indicador `[editando]` aparece; input abre no campo em foco |
| `↑` / `↓` sem input aberto | Cursor de campo move; sem efeito colateral |
| `↑` / `↓` com input aberto | Valor persistido implicitamente; input fechado; cursor move para o campo anterior/próximo |
| `Enter` com input aberto | Valor persistido; input fechado; cursor permanece no campo; mensagem de sucesso exibida |
| `Esc` com input aberto | Valor descartado; valor anterior restaurado; input fechado; cursor permanece no campo |
| `Tab` com input aberto | Valor persistido implicitamente; input fechado; foco vai para a árvore; modo edição encerrado |
| `Tab` sem input aberto | Foco vai para a árvore; modo edição encerrado |
| `Esc` sem input aberto | Modo edição encerrado; retorna ao Modo Leitura; indicador `[editando]` removido |
| `⌃N` | Input do nome abre no cabeçalho; cursor de campo da lista some |
| `Enter` com input do nome aberto | Nome validado; se válido: persistido, input fechado, cursor de campo da lista retorna; se inválido: mensagem de erro, input permanece |
| `Esc` com input do nome aberto | Nome descartado; nome anterior restaurado; input fechado; cursor de campo da lista retorna |
| `Tab` com input do nome aberto | Nome descartado (sem persistir); foco vai para a árvore; modo edição encerrado |
| Campo sensível: input abre | Valor revelado automaticamente em texto claro no input |
| Campo sensível: input fecha | Campo re-mascarado imediatamente |
| `⌃Q` (sair da aplicação) | Modo edição encerrado sem diálogo de confirmação (persistência imediata por campo elimina estado pendente) |

---

#### Comportamento

- **Persistência imediata por campo** — cada campo é salvo ao confirmar (`Enter` ou movimento implícito com `↑`/`↓`/`Tab`); não há estado de "edição pendente" global. `⌃Q` pode sair sem diálogo de confirmação relacionado ao modo edição
- **Input inline** — o input abre na mesma posição da linha do valor, substituindo-a visualmente; o rótulo permanece acima; a estrutura do painel não se desloca
- **Campo sensível revelado no input** — ao abrir o input de um campo sensível, o valor real é exibido em texto claro para permitir edição; ao fechar o input (por qualquer tecla), o campo é re-mascarado imediatamente, independentemente do resultado (confirmado ou cancelado)
- **`⌃R` indisponível no modo edição** — o ciclo de reveal do Modo Leitura não se aplica; o reveal ocorre automaticamente ao abrir o input
- **`⌃C` indisponível no modo edição** — cópia de campo não está disponível enquanto o modo edição está ativo
- **Navegação sem abrir input** — `↑`/`↓`/`Home`/`End` movem o cursor entre campos sem abrir o input, igual ao Modo Leitura; o input só abre com `Enter` explícito
- **Input do nome (`⌃N`) é independente do cursor de campo da lista** — enquanto o input do nome está aberto, nenhum campo da lista está em foco; ao fechar o input do nome, o cursor retorna ao campo que estava em foco antes de `⌃N`
- **Validação do nome** — o nome não pode ser vazio; não pode duplicar o nome de outro segredo na mesma pasta; a validação ocorre ao pressionar `Enter` no input do nome; erros mantêm o input aberto
- **Sair do modo edição** — `Esc` sem input aberto ou `Tab` encerram o modo edição; o indicador `[editando]` é removido; o painel retorna ao Modo Leitura com o mesmo campo em foco
- **Scroll** — o comportamento de scroll é idêntico ao Modo Leitura; a coluna da trilha é sempre reservada

---

### Painel Direito: Detalhe do Segredo — Modo Edição de Estrutura

**Contexto:** Área de trabalho — Modo Cofre. Ativado quando o usuário pressiona `⌃E` na árvore, no painel em Modo Leitura ou no painel em Modo Edição de Valores.
**Largura:** ~65% da área de trabalho (igual ao Modo Leitura).
**Responsabilidade:** Permitir alterar a estrutura dos campos do segredo — renomear rótulos, inserir campos, excluir campos e reordenar campos. Valores dos campos não são editados neste modo.

> Restrições do domínio que este modo deve respeitar:
> - A **Observação** é não-deletável, não-renomeável e não-movível — ocupa sempre a última posição e é excluída da navegação do cursor neste modo
> - O **tipo** de um campo (`texto` / `texto_sensivel`) não pode ser alterado após criação — apenas na inserção
> - Nomes de campo **não têm restrição de unicidade**

---

#### Anatomia do modo

O Modo Edição de Estrutura é uma camada sobre o painel de detalhe. O layout permanece o mesmo (cabeçalho, separador, campos, observação, scroll). O que muda:

1. **Indicador de modo** — `[estrutura]` em `accent.primary` **bold** no cabeçalho, no lugar de `[editando]`
2. **Cursor de campo** — `special.highlight` no bloco do campo em foco, como nos outros modos; o cursor navega apenas entre campos editáveis (Observação excluída)
3. **Rótulo em destaque** — o rótulo do campo em foco recebe ênfase adicional (`text.primary` **bold**) para comunicar que é o alvo das ações de estrutura
4. **Input inline de rótulo** — quando um rótulo está em edição, o texto do rótulo é substituído por um input na mesma linha
5. **Barra de comandos** — exibe as ações do modo estrutura

---

#### Anatomia do cabeçalho em modo estrutura

```
  Gmail [estrutura] ★                    Geral › Sites e Apps
  ──────────────────────────────────────────────────────────
```

- Nome do segredo: `text.primary` **bold**
- `[estrutura]`: `accent.primary` **bold**, separado do nome por um espaço
- `★` e breadcrumb: inalterados

---

#### Wireframes

**Cursor no campo, sem input aberto:**

```
  Gmail [estrutura] ★              Geral › Sites e Apps
  ──────────────────────────────────────────────────────
  URL
  https://accounts.google.com/login

  Usuário                                                ← special.highlight no bloco; rótulo bold
  fulano@gmail.com

  Senha
  ••••••••••

  ╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌╌
  Conta pessoal principal.
```

> Barra: `Enter Renomear · !↑ Mover cima · !↓ Mover baixo · !Ins Inserir · !Del Excluir · Tab Árvore · Esc Sair · F1 Ajuda`
> Observação não tem cursor de foco — está visível mas excluída da navegação do modo estrutura.

**Input de rótulo aberto (`Enter`):**

```
  Gmail [estrutura] ★              Geral › Sites e Apps
  ──────────────────────────────────────────────────────
  URL
  https://accounts.google.com/login

  ░Usuário▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  ← input inline na linha do rótulo
  fulano@gmail.com

  Senha
  ••••••••••
```

> `░` marca o fundo do input (`input.background`); `▌` é o cursor de texto. O valor do campo permanece visível abaixo (leitura, sem alteração). Barra: `Enter Confirmar · Esc Cancelar · F1 Ajuda`

**Input de rótulo aberto — campo sensível:**

```
  Gmail [estrutura] ★              Geral › Sites e Apps
  ──────────────────────────────────────────────────────
  URL
  https://accounts.google.com/login

  Usuário
  fulano@gmail.com

  ░Senha▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  ← input do rótulo
  ••••••••••                                             ← valor permanece mascarado
```

> Campo sensível permanece mascarado no modo estrutura — não há reveal automático ao editar o rótulo.

**Inserção de novo campo (`!Ins`):**

```
  Gmail [estrutura] ★              Geral › Sites e Apps
  ──────────────────────────────────────────────────────
  URL
  https://accounts.google.com/login

  Usuário                                                ← campo com foco antes de !Ins
  fulano@gmail.com

  ░▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  [texto] ⌃T  ← novo campo inserido abaixo; input vazio + badge de tipo
                                                         ← valor vazio (campo novo)
  Senha
  ••••••••••
```

> O novo campo é inserido imediatamente abaixo do campo em foco e acima da Observação (se o foco estiver no último campo editável, o novo campo é inserido entre ele e a Observação). O input do rótulo abre automaticamente com o cursor. O badge `[texto]` indica o tipo atual; `⌃T` alterna entre `[texto]` e `[sensível]` enquanto o input está aberto. Barra: `Enter Confirmar · ⌃T Tipo · Esc Cancelar · F1 Ajuda`

**Badge de tipo alternado para sensível:**

```
  ░▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  [sensível] ⌃T
```

> Após `⌃T`, o badge muda para `[sensível]`. O campo ainda não tem rótulo nem valor. `Enter` confirma nome e tipo.

**Reordenar campo (`!↑` / `!↓`):**

```
  Gmail [estrutura] ★              Geral › Sites e Apps
  ──────────────────────────────────────────────────────
  URL
  https://accounts.google.com/login

  Senha                                                  ← campo movido para cima com !↑ (era abaixo de Usuário)
  ••••••••••

  Usuário                                                ← special.highlight — campo em foco, foi deslocado para baixo
  fulano@gmail.com
```

> A reordenação é imediata e visível — o bloco do campo em foco se desloca e o cursor acompanha. O foco permanece no campo que foi movido.

---

#### Mapa de teclas

**Com cursor de campo, sem input aberto:**

| Tecla | Efeito | Condição |
|---|---|---|
| `↑` / `↓` | Move cursor para o campo anterior / próximo | Apenas entre campos editáveis (Observação excluída) |
| `Home` / `End` | Cursor vai ao primeiro / último campo editável | — |
| `Enter` | Abre input inline no rótulo do campo em foco | — |
| `!↑` | Move o campo em foco uma posição acima | Sem efeito no primeiro campo editável |
| `!↓` | Move o campo em foco uma posição abaixo | Sem efeito no último campo editável (antes da Observação) |
| `!Ins` | Insere novo campo abaixo do campo em foco; input do rótulo abre automaticamente | — |
| `!Del` | Exclui o campo em foco imediatamente e irreversivelmente | — |
| `Tab` | Foco → árvore; sai do modo estrutura | — |
| `Esc` | Sai do modo estrutura; retorna ao Modo Leitura | — |
| `⌃E` | — (sem efeito — já está no modo estrutura) | — |

**Com input de rótulo aberto (`Enter` ou via `!Ins`):**

| Tecla | Efeito |
|---|---|
| Texto / Backspace / Delete | Edita o nome do rótulo |
| `⌃T` | Alterna o tipo do campo entre `texto` e `texto_sensivel` (apenas disponível em inserção — ver nota) |
| `Enter` | Valida e persiste o rótulo (e tipo, se inserção); fecha input; cursor permanece no campo |
| `Esc` | Cancela; restaura o rótulo anterior (ou descarta inserção); fecha input |
| `↑` | Persiste implicitamente; fecha input; move cursor para o campo anterior |
| `↓` | Persiste implicitamente; fecha input; move cursor para o próximo campo |
| `Tab` | Persiste implicitamente; fecha input; foco vai para a árvore; sai do modo estrutura |

> **`⌃T` (toggle de tipo) só está disponível durante a inserção** (`!Ins`). Em renomeação de campo existente, o tipo é imutável — `⌃T` não tem efeito e o badge de tipo não é exibido.

---

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Nome do segredo (cabeçalho) | `text.primary` | **bold** |
| `[estrutura]` | `accent.primary` | **bold** |
| `★` favorito | `accent.secondary` | — |
| Breadcrumb de pasta | `text.secondary` | — |
| Fundo do campo em foco (sem input) | `special.highlight` | — |
| Rótulo do campo em foco (sem input) | `text.primary` | **bold** |
| Rótulo dos campos fora do foco | `text.secondary` | **bold** |
| Fundo do input de rótulo | `input.background` | — |
| Texto dentro do input de rótulo | `text.primary` | — |
| Cursor de texto no input | terminal padrão | — |
| Badge de tipo `[texto]` / `[sensível]` | `text.secondary` | — |
| Valores dos campos (leitura) | inalterados do Modo Leitura | — |
| Separador `───` cabeçalho | `border.default` | — |
| Separador `╌╌╌` da Observação | `border.default` | — |
| Observação (texto) | `text.secondary` | *italic* (diferenciada do modo leitura para comunicar inatividade) |

> A Observação recebe `text.secondary` *italic* no modo estrutura para sinalizar visualmente que está excluída da navegação e das ações.

---

#### Estados dos componentes

| Componente | Estado | Condição |
|---|---|---|
| Indicador `[estrutura]` | visível no cabeçalho | Modo estrutura ativo |
| Cursor de campo | `special.highlight` no bloco | Sempre (modo estrutura tem foco implícito) |
| Cursor de campo | ausente na Observação | Observação nunca recebe foco no modo estrutura |
| Rótulo do campo em foco | `text.primary` **bold** | — |
| Input de rótulo | ausente | `Enter` não pressionado |
| Input de rótulo | aberto sobre a linha do rótulo | `Enter` pressionado; ou `!Ins` executado |
| Badge `[texto]` / `[sensível]` | visível à direita do input | Apenas durante inserção (`!Ins`) |
| Badge `[texto]` / `[sensível]` | ausente | Renomeação de campo existente |
| Observação | visível, não focável, `text.secondary` *italic* | Sempre no modo estrutura |
| Campo sensível | mascarado `••••••••` | Sempre no modo estrutura (sem reveal) |
| Campo recém-inserido | input do rótulo aberto, vazio | Imediatamente após `!Ins` |

---

#### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Modo estrutura ativado | Dica | `• Enter para renomear · !Ins inserir · !Del excluir · !↑↓ mover` |
| Rótulo renomeado confirmado | Sucesso (3s) | `✓ Campo renomeado` |
| Campo inserido | Sucesso (3s) | `✓ Campo "[nome]" adicionado` |
| Campo excluído | Sucesso (3s) | `✓ Campo "[nome]" excluído` |
| Rótulo vazio ao confirmar | Erro | `✗ O nome do campo não pode ser vazio` |
| `!Del` no único campo editável | Erro | `✗ O segredo deve ter pelo menos um campo` |
| `!↑` no primeiro campo | — | Sem mensagem — ação sem efeito silenciosa |
| `!↓` no último campo editável | — | Sem mensagem — ação sem efeito silenciosa |

---

#### Eventos

| Evento | Efeito |
|---|---|
| `⌃E` no Modo Leitura | Modo estrutura ativado; indicador `[estrutura]` aparece; cursor vai ao primeiro campo editável |
| `⌃E` no Modo Edição de Valores | Modo valores encerrado (sem persistência pendente — imediata); modo estrutura ativado |
| `⌃E` na árvore | Painel recebe foco; modo estrutura ativado; cursor vai ao primeiro campo editável |
| `↑` / `↓` sem input aberto | Cursor move entre campos editáveis (Observação ignorada) |
| `Enter` sem input aberto | Input do rótulo abre no campo em foco |
| `Enter` com input aberto | Rótulo validado; se válido: persistido, input fechado, cursor permanece; se inválido (vazio): mensagem de erro, input permanece |
| `Esc` com input aberto (renomeação) | Rótulo descartado; rótulo anterior restaurado; input fechado |
| `Esc` com input aberto (inserção) | Campo recém-inserido descartado; cursor retorna ao campo que estava em foco antes de `!Ins` |
| `!↑` | Campo em foco sobe uma posição; cursor acompanha; persistido imediatamente |
| `!↓` | Campo em foco desce uma posição; cursor acompanha; persistido imediatamente |
| `!Ins` | Novo campo inserido abaixo do campo em foco (tipo `texto`); input do rótulo abre automaticamente com cursor; badge `[texto]` visível |
| `⌃T` com input de inserção aberto | Tipo alterna entre `texto` e `texto_sensivel`; badge atualiza imediatamente |
| `Enter` com input de inserção | Rótulo e tipo confirmados; campo inserido definitivamente; input fechado; cursor no novo campo |
| `!Del` | Campo em foco excluído imediatamente; cursor vai ao campo seguinte (ou anterior se era o último editável) |
| `Esc` sem input aberto | Modo estrutura encerrado; retorna ao Modo Leitura; indicador `[estrutura]` removido |
| `Tab` sem input aberto | Foco vai para a árvore; modo estrutura encerrado |
| `Tab` com input aberto | Rótulo persistido implicitamente; input fechado; foco vai para a árvore; modo encerrado |
| `⌃Q` | Saída da aplicação; persiste o que já foi confirmado (imediato por operação) |

---

#### Comportamento

- **Persistência imediata por operação** — cada ação confirmada (renomear, inserir, mover, excluir) persiste em memória imediatamente; não há um "cancelar tudo" ao sair do modo. `Esc` só cancela o input atualmente aberto, não as operações já confirmadas
- **Observação excluída da navegação** — o cursor de campo nunca vai para a Observação no modo estrutura; `↑`/`↓`/`Home`/`End` ignoram a Observação; `!↓` no último campo editável não tem efeito (não pode ultrapassar a Observação)
- **Tipo imutável em campos existentes** — `⌃T` só funciona durante a inserção de novo campo (`!Ins`); o badge de tipo só é exibido nesse contexto; em renomeação, o tipo não é alterável e o badge não aparece
- **`!Del` é irreversível** — a exclusão ocorre imediatamente ao pressionar `!Del`, sem confirmação; o campo e seu valor são descartados; se o segredo tiver apenas um campo editável, a exclusão é bloqueada com mensagem de erro
- **`!Del` move o cursor** — após excluir, o cursor vai para o campo seguinte; se era o último campo editável, vai para o anterior
- **Input inline de rótulo** — o input substitui visualmente a linha do rótulo; o valor do campo permanece visível abaixo em modo leitura durante a edição do rótulo (o modo estrutura não altera valores)
- **Campo sensível permanece mascarado** — no modo estrutura, campos sensíveis exibem `••••••••`; não há reveal nem `⌃R`
- **Inserção abaixo do foco, acima da Observação** — se o foco está no último campo editável, o novo campo é inserido imediatamente antes da Observação; se o foco está em outro campo, é inserido imediatamente abaixo do campo em foco
- **Troca de modo** — `⌃E` no Modo Edição de Valores troca para o modo estrutura sem diálogo; a persistência imediata do modo valores garante que não há dado pendente a perder
- **Sair do modo** — `Esc` sem input aberto ou `Tab` encerram o modo estrutura; o indicador `[estrutura]` é removido; o painel retorna ao Modo Leitura
- **Scroll** — idêntico ao Modo Leitura; a coluna da trilha é sempre reservada

---

## Telas
