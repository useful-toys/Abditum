# Especificação Visual — FilePicker

> Diálogo funcional de seleção de arquivo (modos Open e Save).
> Parte da [Especificação Visual](tui-specification.md).
>
> **Documentos relacionados:**
> - [`tui-design-system.md`](tui-design-system.md) — fundações visuais
> - [`tui-spec-dialogos.md`](tui-spec-dialogos.md) — anatomia comum e tipos de diálogo

## FilePicker

**Contexto de uso:** abrir ou salvar arquivo do cofre.
**Token de borda:** `border.focused`
**Dimensionamento:** largura máxima do DS (70 colunas ou 80% do terminal, o menor); altura 80% do terminal. Proporção árvore/arquivos ~40/60.
**Diretório inicial:** determinado pelo fluxo orquestrador. Se não informado, CWD do processo. Se o CWD não existe ou não tem permissão de leitura, fallback para home do usuário (`~`).
**Nome sugerido (modo Save):** determinado pelo fluxo orquestrador. Se não informado, campo inicia vazio. O campo não possui placeholder.
**Filtro de extensão:** apenas arquivos com a extensão `<ext>` (parâmetro `extensao`) são exibidos no painel de arquivos. Não há campo de filtro editável. Arquivos e diretórios ocultos (nome iniciado com `.`) não são exibidos. A extensão é omitida na exibição dos nomes de arquivo (redundante — o filtro já restringe ao formato).
**Padding:** 2 colunas horizontal; **0 vertical** — exceção ao DS [Dimensionamento de diálogos](tui-design-system.md#dimensionamento-de-diálogos). Justificativa: princípio "O Terminal como Meio" — espaço vertical é recurso escasso; o FilePicker é o diálogo mais denso da aplicação (caminho + 2 painéis + campo `Arquivo:` no modo Save). As bordas `╭╮╰╯` e os headers internos (`Estrutura`, `Arquivos`) criam contenção e separação suficientes sem padding vertical.

O FilePicker opera em dois modos — **Open** e **Save** — com wireframes e condições distintos. Ambos compartilham a mesma anatomia de painéis.

> Nos wireframes abaixo, `░` representa áreas com fundo `surface.input` (campos de entrada).

> **Decisão de layout:** o FilePicker usa separadores internos com junctions em T (`├┬┴┤`) e painéis lado a lado — estrutura que não se encaixa no modelo padrão de diálogos do DS. Esta configuração foi documentada como **exceção justificada** (ver [DS — Exceções ao dimensionamento](tui-design-system.md#dimensionamento-de-diálogos)) e não promoveu uma subseção no DS porque: (1) o FilePicker é o único diálogo com essa complexidade; (2) é um padrão de SO consolidado, não um padrão reutilizável interno; (3) o mecanismo de exceção do DS cobre o caso. Se um segundo diálogo com painéis internos surgir, a exceção será promovida a subseção.

**Barra de comandos durante FilePicker:** enquanto o FilePicker está ativo, a barra de comandos exibe apenas as ações internas do diálogo (conforme regra geral de [Barra de Comandos durante diálogo ativo](tui-spec-barras.md#anatomia)). Ações de confirmação/cancelamento (`Enter`/`Esc`) já estão na borda do diálogo — não são duplicadas na barra.

```
  Tab Painel                                                                  F1 Ajuda
```

| Ação | Tecla | Descrição |
|---|---|---|
| Alternar painel | `Tab` | Cicla foco entre os painéis (Árvore → Arquivos no modo Open; Árvore → Arquivos → Campo Nome no modo Save) |
| Ajuda | `F1` | Abre o Help — âncora fixa |

---

### Contrato de entrada e saída

**Entrada (parâmetros do orquestrador):**

| Parâmetro | Tipo | Obrigatório | Uso |
|---|---|---|---|
| `modo` | `Open \| Save` | Sim | Define título, ações e presença do campo de nome |
| `extensao` | `String` | Sim | Extensão filtrada e adicionada automaticamente ao salvar (ex: `".abditum"`, `".json"`). Deve incluir o ponto inicial. |
| `diretorio_inicial` | `PathBuf` | Não | Diretório onde o FilePicker abre. Default: CWD → fallback `~` |
| `nome_sugerido` | `String` | Não (modo Save) | Valor inicial do campo `Arquivo:`. Default: vazio |

**Saída (retorno ao orquestrador):**

| Resultado | Valor | Significado |
|---|---|---|
| Confirmado | `Some(PathBuf)` | Caminho completo do arquivo selecionado (modo Open) ou caminho de salvamento com extensão `<ext>` garantida (modo Save) |
| Cancelado | `None` | Usuário abandonou o diálogo via `Esc` |

---

### FilePicker — Modo Open

**Título:** `Abrir cofre`
**Objetivo:** selecionar um arquivo `<ext>` existente.

**Wireframe (arquivo selecionado — ação default ativa, scroll em ambos os painéis):**

```
╭── Abrir cofre ─────────────────────────────────────────────────────╮
│  /home/usuario/projetos/abditum                                    │
├─ Estrutura ──────────────────┬─ Arquivos ──────────────────────────┤
│  ▶ /                         ↑  ● database   25.8 MB 15/03/25 14:32↑
│    ▼ usuario                 │  ● config       1.2 KB 02/01/25 09:15│
│      ▶ documentos            │  ● backup      18.4 MB 04/04/25 18:47│
│      ▼ projetos              │                                     │
│        ▶ site                │                                     │
│        ▼ abditum             ■                                     ■
│          ▶ docs              │                                     │
│          ▶ src               │                                     │
│        ▶ outros              │                                     │
│      ▶ downloads             ↓                                     ↓
╰── Enter Abrir ──────────────┴────────────────────── Esc Cancelar ──╯
       ↑ accent.primary + bold (desbloqueado)
```

> Scroll da árvore (`↑` `■` `↓`) substitui o `│` do separador entre painéis. Scroll dos arquivos (`↑` `■` `↓`) substitui o `│` da borda direita do modal. O `┴` na borda inferior marca a junção do separador com a base do diálogo. Metadados (tamanho + `dd/mm/aa HH:MM`) na mesma linha do nome.

**Wireframe (nenhum arquivo — ação default bloqueada, sem scroll):**

```
╭── Abrir cofre ─────────────────────────────────────────────────────╮
│  /home/usuario/documentos                                          │
├─ Estrutura ──────────────────┬─ Arquivos ──────────────────────────┤
│  ▶ /                         │                                     │
│    ▼ usuario                 │  Nenhum cofre neste diretório       │
│      ▼ documentos            │                                     │
│        ▶ fotos               │                                     │
│        ▶ textos              │                                     │
│                              │                                     │
╰── Enter Abrir ──────────────┴────────────────────── Esc Cancelar ──╯
       ↑ text.disabled (bloqueado)
```

| Elemento | Token | Atributo |
|---|---|---|
| Título `Abrir cofre` | `text.primary` | **bold** |
| Header `Estrutura` | `text.secondary` | **bold** |
| Header `Arquivos` | `text.secondary` | **bold** |
| Separadores internos (`├`, `┬`, `┴`, `─`, `│`) | `border.default` | — |
| Pasta selecionada na árvore | `accent.primary` | **bold** |
| Pasta não selecionada | `text.primary` | — |
| Indicador de pasta (`▶` recolhida, `▼` expandida, `▷` vazia) | `accent.secondary` | — |
| Arquivo selecionado no painel de arquivos | `special.highlight` (fundo) + `text.primary` | **bold** |
| Arquivo não selecionado | `text.primary` | — |
| Indicador de arquivo `●` | `text.secondary` | — |
| Nome do arquivo (sem extensão `<ext>`) | — | Extensão omitida na exibição — redundante com o filtro |
| Metadados (tamanho, data/hora) | `text.secondary` | — |
| Texto `Nenhum cofre neste diretório` | `text.secondary` | — |
| Valor do caminho | `text.secondary` | — |
| Ação default (bloqueada) | `text.disabled` | — |
| Ação default (desbloqueada) | `accent.primary` | **bold** |

**Estados dos componentes:**

| Componente | Estado | Condição |
|---|---|---|
| Painel `Estrutura` (árvore) | sempre visível | — |
| Painel `Arquivos` (lista) | conteúdo visível | Pasta selecionada contém arquivos `<ext>` |
| Painel `Arquivos` (lista) | texto vazio | Pasta selecionada **não** contém arquivos `<ext>` |
| Caminho (valor) | sempre visível, somente leitura | Atualiza ao navegar na árvore |
| Arquivo pré-selecionado no painel | selecionado | Primeiro `<ext>` da pasta, automaticamente ao entrar na pasta |
| Ação `Enter Abrir` | bloqueada (`text.disabled`) | Pasta sob cursor não contém arquivos `<ext>` |
| Ação `Enter Abrir` | ativa (`accent.primary` **bold**) | Pasta sob cursor contém `<ext>` (pré-seleção automática habilita a ação, mesmo com foco na árvore) |
| Ação `Esc Cancelar` | sempre ativa | — |

**Mensagens:**

| Contexto | Tipo | Texto |
|---|---|---|
| Diálogo abre / foco na árvore | Dica de campo | `• Navegue pelas pastas e selecione um cofre` |
| Foco no painel de arquivos (com seleção) | Dica de campo | `• Selecione o cofre para abrir` |
| Foco no painel de arquivos (painel vazio) | Dica de campo | `• Nenhum cofre neste diretório — navegue para outra pasta` |
| Diálogo fecha | — | Barra limpa *(orquestrador assume)* |

**Comportamento:**

- **Carregamento lazy:** a árvore não carrega todo o filesystem na abertura. Apenas o caminho até o diretório inicial é expandido. O conteúdo de cada pasta é lido sob demanda ao expandir — evita lentidão em filesystems grandes
- **Foco inicial:** árvore de diretórios (painel esquerdo)
- **Ordem do Tab:** Árvore → Arquivos → volta (2 stops)
- **Scroll:** cada painel tem scroll independente com indicadores `↑`/`↓`/`■` na borda direita do respectivo painel
- **Painel de arquivos reflete o cursor da árvore:** ao mover o cursor (`↑↓`) entre pastas na árvore, o painel de arquivos atualiza imediatamente para mostrar os `<ext>` da pasta sob o cursor — não apenas ao expandir. O caminho exibido e o painel de arquivos acompanham a pasta com cursor, independente de ela estar expandida ou recolhida
- **Navegação por teclado na árvore:** `↑↓` navega entre pastas visíveis; `→` expande pasta recolhida; `←` recolhe pasta expandida; `Enter` avança foco para o primeiro arquivo no painel de arquivos (se a pasta sob o cursor contém `<ext>`; sem efeito se não contém); `Home`/`End` vai ao primeiro/último item visível; `PgUp`/`PgDn` scroll por página
- **Navegação por teclado nos arquivos:** `↑↓` navega entre arquivos; `Enter` confirma seleção (equivale à ação default); `Home`/`End` vai ao primeiro/último arquivo visível; `PgUp`/`PgDn` scroll por página
- Ao navegar para uma pasta na árvore, se ela contém arquivos `<ext>`, o primeiro é pré-selecionado automaticamente no painel de arquivos
- **Indicador de pasta vazia:** pastas sem subdiretórios visíveis usam `▷` conforme o DS — não são expansíveis. `→` não tem efeito sobre elas (nada a expandir). `Enter` segue a regra padrão — avança foco para o painel de arquivos se a pasta contém `<ext>`. `▷` indica ausência de subdiretórios expansíveis — não impede que a pasta contenha arquivos `<ext>` exibidos no painel de arquivos
- **Clique simples em pasta:** move cursor para a pasta (atualiza painel de arquivos e caminho exibido)
- **Clique simples em arquivo:** seleciona o arquivo (highlight)
- **Duplo-clique em pasta:** expande/recolhe (mesmo que `→`/`←`)
- **Duplo-clique em arquivo:** confirma seleção (mesmo que ação default)
- **Scroll do mouse:** afeta o painel com foco
- **Arquivos e diretórios ocultos** (nome iniciado com `.`) não são exibidos
- **Caminho longo:** truncado no início com `…` (ex: `…/projetos/abditum`)
- **Diretórios sem permissão:** exibidos normalmente na árvore; ao tentar expandir, erro na barra (`✕ Sem permissão para acessar <pasta>`) e pasta permanece recolhida
- **Fallback de CWD:** se o CWD é inacessível, o FilePicker navega para home do usuário (`~`) e exibe mensagem informativa (`⚠ Diretório atual inacessível — navegando para home`)

**Ordenação:**

| Painel | Critério | Detalhes |
|---|---|---|
| Árvore (pastas) | Alfabético, case-insensitive | Ordem lexicográfica (`a` = `A`) |
| Arquivos | Alfabético, case-insensitive | Ordem lexicográfica pelo nome sem extensão |

**Indentação da árvore:** 2 espaços por nível de profundidade.

**Formato dos metadados:**

| Campo | Formato | Exemplo |
|---|---|---|
| Tamanho | `{valor} {unidade}` — base 1024, unidades KB/MB/GB, 1 casa decimal | `25.8 MB`, `1.2 KB`, `18.4 MB` |
| Data/hora | `dd/mm/aa HH:MM` — dígitos numéricos, locale local | `15/03/25 14:32` |

**Alinhamento dos metadados:** no painel de arquivos, os metadados são alinhados em colunas — tamanho alinhado à direita, data/hora em posição fixa. O nome do arquivo ocupa o espaço restante à esquerda. Isso facilita a leitura por scanning vertical.

**Comportamento na raiz:** `←` na pasta raiz (`/`) não tem efeito — a seleção permanece na raiz.

**Truncamento de metadados:** em terminais estreitos, os metadados são os primeiros a truncar (direita → esquerda). O nome do arquivo tem prioridade e só trunca se não houver espaço mesmo para ele.

**Transições especiais:**

| Evento | Efeito |
|---|---|
| Cursor move para pasta sem `<ext>` | Painel de arquivos mostra texto vazio; ação default muda para `text.disabled` |
| Cursor move para pasta com `<ext>` | Primeiro arquivo pré-selecionado; ação default muda para `accent.primary` **bold** |
| `Enter` no painel de arquivos | Diálogo fecha com o arquivo selecionado |
| `Enter` na árvore (pasta com `<ext>`) | Foco avança para o primeiro arquivo no painel de arquivos |
| `Enter` na árvore (pasta sem `<ext>`) | Sem efeito |
| `→` em pasta recolhida | Pasta expandida; cursor permanece na pasta |
| `←` em pasta expandida | Pasta recolhida; cursor permanece na pasta |
| `→` em pasta `▷` (vazia) | Sem efeito (nada a expandir) |
| Tentar expandir pasta sem permissão | Erro na barra (`✕ Sem permissão para acessar <pasta>`); pasta permanece recolhida |

---

### FilePicker — Modo Save

**Título:** `Salvar cofre`
**Objetivo:** escolher diretório e nome para salvar o arquivo do cofre.

**Wireframe (campo nome preenchido — ação default ativa):**

```
╭── Salvar cofre ────────────────────────────────────────────────────╮
│  /home/usuario/projetos/abditum                                    │
├─ Estrutura ──────────────────┬─ Arquivos ──────────────────────────┤
│  ▶ /                         │  ● database   25.8 MB 15/03/25 14:32│
│    ▼ usuario                 │  ● config       1.2 KB 02/01/25 09:15│
│      ▼ projetos              │                                     │
│        ▼ abditum             │                                     │
│          ▶ docs              │                                     │
│                              │                                     │
├──────────────────────────────┴─────────────────────────────────────┤
│  Arquivo: ░meu-cofre▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  │
╰── Enter Salvar ───────────────────────────────────────── Esc Cancelar ──╯
       ↑ accent.primary + bold (desbloqueado)
```

**Wireframe (campo nome vazio — ação default bloqueada):**

```
╭── Salvar cofre ────────────────────────────────────────────────────╮
│  /home/usuario/projetos                                            │
├─ Estrutura ──────────────────┬─ Arquivos ──────────────────────────┤
│  ▶ /                         │  ● database   25.8 MB 15/03/25 14:32│
│    ▼ usuario                 │                                     │
│      ▼ projetos              │                                     │
│                              │                                     │
├──────────────────────────────┴─────────────────────────────────────┤
│  Arquivo: ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  │
╰── Enter Salvar ───────────────────────────────────────── Esc Cancelar ──╯
       ↑ text.disabled (bloqueado)
```

> Tokens de estrutura (título, headers, separadores, pasta, arquivo, metadados, caminho, ação default) idênticos ao [Modo Open](#filepicker--modo-open). Exclusivos do Modo Save:

| Elemento | Token | Atributo |
|---|---|---|
| Rótulo `Arquivo:` (campo ativo) | `accent.primary` | **bold** |
| Rótulo `Arquivo:` (campo inativo) | `text.secondary` | — |
| Área do campo `░` | `surface.input` | — |
| Cursor `▌` | `text.primary` | — |

**Estados dos componentes:**

| Componente | Estado | Condição |
|---|---|---|
| Painel `Estrutura` (árvore) | sempre visível | — |
| Painel `Arquivos` (lista) | conteúdo visível | Pasta selecionada contém arquivos `<ext>` |
| Painel `Arquivos` (lista) | texto vazio | Pasta selecionada **não** contém arquivos `<ext>` |
| Caminho (valor) | sempre visível, somente leitura | Atualiza ao navegar na árvore |
| Campo `Arquivo:` | sempre visível | — |
| Caracteres inválidos para filesystem (`/ \ : * ? " < > \|`) | bloqueados silenciosamente | Tecla não produz efeito — sem mensagem de erro |
| Extensão `<ext>` | adicionada automaticamente | Se o nome digitado não termina em `<ext>` |
| Ação `Enter Salvar` | bloqueada (`text.disabled`) | Campo `Arquivo:` vazio |
| Ação `Enter Salvar` | ativa (`accent.primary` **bold**) | Campo `Arquivo:` não vazio |
| Ação `Esc Cancelar` | sempre ativa | — |

> **Nota:** a validação de sobrescrita (arquivo já existe) é responsabilidade do fluxo que chamou o FilePicker, não do diálogo. O picker retorna o caminho completo; o fluxo abre diálogo de Confirmação × Destrutivo se necessário.

**Mensagens:**

| Contexto | Tipo | Texto |
|---|---|---|
| Diálogo abre / foco na árvore | Dica de campo | `• Navegue pelas pastas e escolha onde salvar` |
| Foco no painel de arquivos | Dica de campo | `• Arquivos existentes neste diretório` |
| Foco no campo `Arquivo:` (vazio) | Dica de campo | `• Digite o nome do arquivo — <ext> será adicionado automaticamente` |
| Foco no campo `Arquivo:` (preenchido) | Dica de campo | `• Confirme para salvar o cofre` |
| Diálogo fecha | — | Barra limpa *(orquestrador assume)* |

**Comportamento:**

- **Foco inicial:** árvore de diretórios (painel esquerdo)
- **Ordem do Tab:** Árvore → Arquivos → Campo `Arquivo:` → volta (3 stops)
- **Scroll:** cada painel tem scroll independente com indicadores `↑`/`↓`/`■` na borda direita do respectivo painel
- Navegação na árvore e painel de arquivos idêntica ao modo Open, com uma exceção: **`Enter` no painel de arquivos copia o nome (sem extensão) para o campo `Arquivo:` e move foco para o campo** — não confirma o diálogo. A confirmação requer `Enter` novamente (no campo ou em qualquer contexto com ação default ativa)
- No painel de arquivos: `↑↓` apenas destaca o arquivo (highlight) — **não** copia o nome para o campo. Somente `Enter` ou clique simples no arquivo copiam o nome (sem extensão) para o campo `Arquivo:`
- Ao navegar na árvore, o campo `Arquivo:` **não é limpo** — preserva o nome digitado
- Extensão `<ext>` é adicionada silenciosamente ao caminho de retorno, sem alterar o texto exibido no campo
- **Duplo-clique em pasta:** expande/recolhe (mesmo que `→`/`←`)
- **Duplo-clique em arquivo existente:** copia o nome para o campo `Arquivo:`
- Scroll do mouse, arquivos ocultos, caminho longo, permissões, fallback CWD, ordenação, indentação, formato de metadados e truncamento: idêntico ao [Modo Open](#filepicker--modo-open)

**Transições especiais:**

| Evento | Efeito |
|---|---|
| Clique simples em arquivo existente no painel | Nome copiado para campo `Arquivo:`; ação default muda para `accent.primary` **bold** |
| `Enter` no painel de arquivos | Nome copiado para campo `Arquivo:`; foco move para o campo. **Não** confirma o diálogo |
| `Enter` na árvore (pasta com `<ext>`) | Foco avança para o primeiro arquivo no painel de arquivos |
| `Enter` na árvore (pasta sem `<ext>`) | Sem efeito |
| `→` em pasta recolhida | Pasta expandida; cursor permanece na pasta |
| `←` em pasta expandida | Pasta recolhida; cursor permanece na pasta |
| Limpar campo `Arquivo:` | Ação default volta para `text.disabled` |
| `Enter` com campo preenchido | Diálogo fecha com caminho completo (diretório + nome + `<ext>`) |
| Tentar expandir pasta sem permissão | Erro na barra (`✕ Sem permissão para acessar <pasta>`); pasta permanece recolhida |

---
