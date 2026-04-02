# Especificação Visual — Abditum TUI

> Wireframes, layouts de componentes e fluxos visuais concretos.
> Cada tela e componente consome os padrões definidos no design system.
>
> **Documento de fundação:**
> - [`tui-design-system-novo.md`](tui-design-system-novo.md) — princípios, tokens, estados, padrões transversais

## Fronteira deste documento

Este documento define **composições** — telas, wireframes de componentes e fluxos visuais concretos.

Quando algo é específico de uma tela ou componente (ex: quais campos o diálogo PasswordCreate exibe, como a árvore de pastas é indentada), pertence aqui. Quando é uma regra que se aplica a qualquer tela (ex: anatomia da moldura de diálogos, tokens de cor, regras de foco), pertence ao [design system](tui-design-system-novo.md).

> **Teste de fronteira:** se eu trocar o nome do item concreto e a regra continuar válida, é **padrão** e pertence ao design system. Se a regra só faz sentido para esta tela ou componente, pertence aqui.

---

## Sumário

- [Diálogos de Decisão](#diálogos-de-decisão)
  - [Referência Visual por Severidade](#referência-visual-por-severidade)
  - [Exemplo: Confirmação × Destrutivo — Excluir segredo](#exemplo-confirmação--destrutivo--excluir-segredo)
  - [Exemplo: Confirmação × Neutro — Alterações não salvas](#exemplo-confirmação--neutro--alterações-não-salvas)
  - [Exemplo: Reconhecimento × Informativo — Conflito detectado](#exemplo-reconhecimento--informativo--conflito-detectado)
  - [Exemplo: Reconhecimento × Erro — Falha ao abrir cofre](#exemplo-reconhecimento--erro--falha-ao-abrir-cofre)
- [Diálogos Funcionais](#diálogos-funcionais)
  - [PasswordEntry](#passwordentry)
  - [PasswordCreate](#passwordcreate)
  - [FilePicker](#filepicker)
    - [FilePicker — Modo Open](#filepicker--modo-open)
    - [FilePicker — Modo Save](#filepicker--modo-save)
  - [Help](#help)
- [Componentes](#componentes)
  - [Cabeçalho](#cabeçalho)
  - [Barra de Mensagens](#barra-de-mensagens)
  - [Barra de Comandos](#barra-de-comandos)

---

## Diálogos de Decisão

Todos os diálogos de decisão seguem a anatomia comum e o modelo bidimensional (Intenção × Severidade) definidos no [design system — Sobreposição](tui-design-system-novo.md#sobreposição). Esta seção define a tabela de referência visual completa e apresenta wireframes ilustrativos de combinações típicas.

---

### Referência Visual por Severidade

A severidade governa **todo** o tratamento visual da moldura. A intenção define apenas a barra de ações (múltiplas opções vs. `Enter OK`). A tabela abaixo é a referência completa — qualquer diálogo de decisão pode ser estilizado a partir dela, sem consultar os exemplos.

| Severidade | Símbolo | Token: borda e título | Token: tecla default | Token: demais ações | Título: atributo |
|---|---|---|---|---|---|
| **Destrutivo** | `⚠` | `semantic.warning` | `semantic.error` | `semantic.warning` | **bold** |
| **Erro** | `✗` | `semantic.error` | `accent.primary` | `semantic.error` | **bold** |
| **Alerta** | `⚠` | `semantic.warning` | `accent.primary` | `semantic.warning` | **bold** |
| **Informativo** | `ℹ` | `semantic.info` | `accent.primary` | `semantic.info` | **bold** |
| **Neutro** | — | `border.focused` | `accent.primary` | `border.focused` | **bold** |

**Barra de ações por intenção:**

| Intenção | Ação default (esquerda) | Ações intermediárias | Ação cancelar (direita) |
|---|---|---|---|
| **Confirmação** | Tecla + label em **bold**, token da tecla default | Tecla + label na cor da borda | `Esc` + label na cor da borda |
| **Reconhecimento** | `Enter OK` em **bold**, token `accent.primary` | — | — (`Esc` fecha, equivalente a OK) |

**Regras derivadas (recordatório):**

- Borda e título usam o **mesmo** token
- Tecla default sempre em **bold**; demais ações sem bold
- Símbolo precede o título na borda superior; Neutro não usa símbolo
- Mensagem interna em `text.primary`; nomes de itens referenciados em **bold**

---

### Exemplo: Confirmação × Destrutivo — Excluir segredo

**Intenção:** Confirmação | **Severidade:** Destrutivo (`⚠`)
**Token de borda:** `semantic.warning`
**Ação default:** `S Excluir` — token `semantic.error` + **bold** (ação destrutiva)

```
╭── ⚠  Excluir segredo ───────────╮
│                                  │
│  Gmail será excluído             │
│  permanentemente.                │
│  Esta ação não pode ser desfeita.│
│                                  │
╰── S Excluir ────────── N Cancelar ╯
```

| Elemento | Token | Atributo |
|---|---|---|
| Borda e título `⚠ Excluir segredo` | `semantic.warning` | **bold** (título) |
| Mensagem | `text.primary` | — |
| Nome do item (`Gmail`) | `text.primary` | **bold** |
| Tecla `S` + label `Excluir` | `semantic.error` | **bold** |
| Tecla `N` + label `Cancelar` | `semantic.warning` | — |

---

### Exemplo: Confirmação × Neutro — Alterações não salvas

**Intenção:** Confirmação | **Severidade:** Neutro (—)
**Token de borda:** `border.focused`
**Ação default:** `S Salvar` — token `accent.primary` + **bold**

```
╭── Alterações não salvas ────────╮
│                                  │
│  Deseja salvar antes de sair?    │
│                                  │
╰── S Salvar ── N Descartar ── Esc Voltar ╯
```

| Elemento | Token | Atributo |
|---|---|---|
| Borda e título `Alterações não salvas` | `border.focused` | **bold** (título) |
| Mensagem | `text.primary` | — |
| Tecla `S` + label `Salvar` | `accent.primary` | **bold** |
| Tecla `N` + label `Descartar` | `border.focused` | — |
| `Esc` + label `Voltar` | `border.focused` | — |

> **Nota:** severidade Neutro não usa símbolo — o título aparece sem prefixo.

---

### Exemplo: Reconhecimento × Informativo — Conflito detectado

**Intenção:** Reconhecimento | **Severidade:** Informativo (`ℹ`)
**Token de borda:** `semantic.info`
**Ação default:** `Enter OK` — token `accent.primary` + **bold**

```
╭── ℹ  Conflito detectado ───────╮
│                                  │
│  O arquivo foi modificado        │
│  externamente.                   │
│                                  │
╰── Enter OK ───────────────────╯
```

| Elemento | Token | Atributo |
|---|---|---|
| Borda e título `ℹ Conflito detectado` | `semantic.info` | **bold** (título) |
| Mensagem | `text.primary` | — |
| `Enter` + label `OK` | `accent.primary` | **bold** |

> **Nota:** diálogos de reconhecimento têm apenas uma ação (`Enter OK`). `Esc` também fecha o diálogo (equivalente a OK para reconhecimento).

---

### Exemplo: Reconhecimento × Erro — Falha ao abrir cofre

**Intenção:** Reconhecimento | **Severidade:** Erro (`✗`)
**Token de borda:** `semantic.error`
**Ação default:** `Enter OK` — token `accent.primary` + **bold**

```
╭── ✗  Falha ao abrir cofre ─────╮
│                                  │
│  O arquivo está corrompido ou    │
│  não é um cofre válido.          │
│                                  │
╰── Enter OK ───────────────────╯
```

| Elemento | Token | Atributo |
|---|---|---|
| Borda e título `✗ Falha ao abrir cofre` | `semantic.error` | **bold** (título) |
| Mensagem | `text.primary` | — |
| `Enter` + label `OK` | `accent.primary` | **bold** |

---

## Diálogos Funcionais

Todos os diálogos funcionais seguem a anatomia comum do [design system — Sobreposição](tui-design-system-novo.md#sobreposição), sem símbolo semântico no título. Esta seção especifica a anatomia interna de cada um.

---

### PasswordEntry

**Contexto de uso:** entrada de senha para abrir cofre.
**Token de borda:** `border.focused`
**Largura fixa:** 50 colunas

> Nos wireframes abaixo, `░` representa a área com fundo `surface.input` (campo de entrada). Na implementação real, o campo é uma área de fundo rebaixado sem hachura — conforme definido em [Campos de entrada de texto](tui-design-system-novo.md#foco-e-navegação).

**Estado inicial (campo vazio — ação default bloqueada):**

```
╭── Senha mestra ────────────────────────────╮
│                                            │
│  Senha                                     │
│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
│                                            │
╰── Enter Confirmar ──────────── Esc Cancelar ╯
       ↑ text.disabled (bloqueado)
```

**Estado com digitação (ação default ativa):**

```
╭── Senha mestra ────────────────────────────╮
│                                            │
│  Senha                                     │
│  ░••••••••▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
│                                            │
╰── Enter Confirmar ──────────── Esc Cancelar ╯
       ↑ accent.primary + bold (desbloqueado)
```

**Estado com contador de tentativas (a partir da 2ª):**

```
╭── Senha mestra ────────────────────────────╮
│                                            │
│  Senha                                     │
│  ░••••••••▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
│                                            │
│  Tentativa 2 de 5                          │
╰── Enter Confirmar ──────────── Esc Cancelar ╯
```

| Elemento | Token | Atributo |
|---|---|---|
| Título `Senha mestra` | `text.primary` | **bold** |
| Label `Senha` | `accent.primary` | **bold** (campo ativo, sempre — diálogo de campo único) |
| Área do campo `░` | `surface.input` | — |
| Placeholder (antes de digitar) | `text.secondary` | *italic* |
| Máscara `••••••••` | `text.secondary` | — |
| Cursor `▌` | `text.primary` | — |
| Contador `Tentativa 2 de 5` | `text.secondary` | — |
| Ação default (bloqueada) | `text.disabled` | — |
| Ação default (desbloqueada) | `accent.primary` | **bold** |

**Estados dos componentes:**

| Componente | Estado | Condição |
|---|---|---|
| Campo `Senha` | sempre visível, sempre com foco | Diálogo de campo único |
| Contador de tentativas | visível | Tentativa atual ≥ 2 |
| Contador de tentativas | oculto | Primeira tentativa |
| Ação `Enter Confirmar` | bloqueada (`text.disabled`) | Campo `Senha` vazio |
| Ação `Enter Confirmar` | ativa (`accent.primary` **bold**) | Campo `Senha` não vazio |
| Ação `Esc Cancelar` | sempre ativa | — |

**Mensagens:**

| Contexto | Tipo | Texto |
|---|---|---|
| Diálogo abre / foco no campo (vazio ou válido) | Dica de campo | `• Digite a senha para desbloquear o cofre` |
| `Enter` → senha incorreta | Erro (5s) | `✗ Senha incorreta` |
| Diálogo fecha (confirmação ou cancelamento) | — | Barra limpa *(orquestrador assume)* |

**Comportamento:**
- Máscara de comprimento fixo (8 `•`) — não revela o tamanho real da senha
- Campo único — `Tab` não faz nada dentro deste diálogo

**Transições especiais:**

| Evento | Efeito |
|---|---|
| `Enter` com senha incorreta | Campo limpo; ação default volta para `text.disabled`; contador incrementado |
| Tentativas esgotadas | Diálogo fecha automaticamente |

---

### PasswordCreate

**Contexto de uso:** criação de senha (ao criar cofre ou alterar senha mestra).
**Token de borda:** `border.focused`
**Largura fixa:** 50 colunas

**Estado inicial (foco no primeiro campo — ação default bloqueada):**

```
╭── Definir senha mestra ───────────────────╮
│                                            │
│  Nova senha                                │
│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
│                                            │
│  Confirmação                               │
│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
│                                            │
╰── Enter Confirmar ──────────── Esc Cancelar ╯
       ↑ text.disabled (bloqueado)
```

**Estado com digitação (primeiro campo ativo, medidor aparece — ação ainda bloqueada):**

```
╭── Definir senha mestra ───────────────────╮
│                                            │
│  Nova senha                                │
│  ░••••••••▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
│                                            │
│  Confirmação                               │
│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
│                                            │
│  Força: ████████░░ Boa                     │
│                                            │
╰── Enter Confirmar ──────────── Esc Cancelar ╯
       ↑ text.disabled (2º campo vazio)
```

**Estado com ambos campos preenchidos (ação default desbloqueada):**

```
╭── Definir senha mestra ───────────────────╮
│                                            │
│  Nova senha                                │
│  ░••••••••░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
│                                            │
│  Confirmação                               │
│  ░••••▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
│                                            │
│  Força: ████████░░ Boa                     │
│                                            │
╰── Enter Confirmar ──────────── Esc Cancelar ╯
       ↑ accent.primary + bold (desbloqueado)
```

| Elemento | Token | Atributo |
|---|---|---|
| Título `Definir senha mestra` | `text.primary` | **bold** |
| Label do campo ativo | `accent.primary` | **bold** |
| Label do campo inativo | `text.secondary` | — |
| Área do campo `░` | `surface.input` | — |
| Placeholder (antes de digitar) | `text.secondary` | *italic* |
| Máscara | `text.secondary` | — |
| Cursor `▌` | `text.primary` | — |
| Medidor — preenchido | `semantic.success` ou `semantic.warning` | — |
| Medidor — vazio | `text.disabled` | — |
| Label de força `Boa` / `Forte` | `semantic.success` | — |
| Label de força `Fraca` | `semantic.warning` | — |
| Ação default (bloqueada) | `text.disabled` | — |
| Ação default (desbloqueada) | `accent.primary` | **bold** |

**Estados dos componentes:**

| Componente | Estado | Condição |
|---|---|---|
| Campo `Nova senha` | sempre visível | — |
| Campo `Confirmação` | sempre visível | — |
| Medidor de força | visível | Campo `Nova senha` não vazio |
| Medidor de força | oculto | Campo `Nova senha` vazio |
| Linha em branco antes do medidor | visível | Medidor visível |
| Ação `Enter Confirmar` | bloqueada (`text.disabled`) | Campo `Nova senha` vazio **ou** campo `Confirmação` vazio |
| Ação `Enter Confirmar` | ativa (`accent.primary` **bold**) | Ambos os campos não vazios |
| Ação `Esc Cancelar` | sempre ativa | — |

> **Nota:** a verificação de igualdade entre as senhas ocorre **no momento do Enter** — não bloqueia a ação default. Se as senhas divergem, o erro é comunicado e o campo de confirmação é limpo.

**Mensagens:**

| Contexto | Tipo | Texto |
|---|---|---|
| Diálogo abre / foco em `Nova senha` (vazio ou válido) | Dica de campo | `• A senha mestra protege todo o cofre — use 12+ caracteres` |
| Foco em `Confirmação` (vazio ou válido) | Dica de campo | `• Redigite a senha para confirmar` |
| Foco em campo com erro prévio de divergência | Erro (5s) | `✗ As senhas não conferem — digite novamente` |
| `Enter` → senhas divergentes | Erro (5s) | `✗ As senhas não conferem — digite novamente` |
| Diálogo fecha (confirmação ou cancelamento) | — | Barra limpa *(orquestrador assume)* |

**Comportamento:**
- `Tab` alterna entre os campos `Nova senha` e `Confirmação`
- Medidor de força atualizado a cada tecla no campo `Nova senha`
- Máscara de comprimento fixo (8 `•`) — não revela o tamanho real da senha

**Transições especiais:**

| Evento | Efeito |
|---|---|
| `Enter` com senhas divergentes | Foco move para `Confirmação`; campo `Confirmação` limpo; ação default volta para `text.disabled` |

---

### FilePicker

**Contexto de uso:** abrir ou salvar arquivo do cofre.
**Token de borda:** `border.focused`
**Dimensionamento:** largura máxima do DS (70 colunas ou 80% do terminal, o menor); altura 80% do terminal. Proporção árvore/arquivos ~40/60.
**Diretório inicial:** CWD do processo.
**Filtro fixo:** apenas arquivos `*.abditum` são exibidos no painel de arquivos. Não há campo de filtro editável.

O FilePicker opera em dois modos — **Open** e **Save** — com wireframes e condições distintos. Ambos compartilham a mesma anatomia de painéis.

> Nos wireframes abaixo, `░` representa áreas com fundo `surface.input` (campos de entrada).

---

#### FilePicker — Modo Open

**Título:** `Abrir cofre`
**Objetivo:** selecionar um arquivo `.abditum` existente.

**Wireframe (arquivo selecionado — ação default ativa):**

```
╭── Abrir cofre ───────────────────────────────────────────────────╮
│  Caminho: /home/usuario/projetos/abditum                         │
├─ Estrutura ──────────────────┬─ Arquivos ────────────────────────┤
│  ▶ /                         │  ● database.abditum   25.8 MB 1h │
│    ▼ usuario/                │  ● config.abditum      1.2 KB 3d │
│      ▶ documentos/           │  ● backup.abditum     18.4 MB 1s │
│      ▼ projetos/             │                                   │
│        ▶ site/               │                                   │
│        ▼ abditum/            │                                   │
│          ▶ docs/             │                                   │
│          ▶ src/              │                                   │
│        ▶ outros/             │                                   │
│      ▶ downloads/            │                                   │
│                              │                                   │
╰── Enter Abrir ────────────────────────────────────── Esc Cancelar ╯
       ↑ accent.primary + bold (desbloqueado)
```

**Wireframe (nenhum arquivo — ação default bloqueada):**

```
╭── Abrir cofre ───────────────────────────────────────────────────╮
│  Caminho: /home/usuario/documentos                               │
├─ Estrutura ──────────────────┬─ Arquivos ────────────────────────┤
│  ▶ /                         │                                   │
│    ▼ usuario/                │  Nenhum arquivo .abditum          │
│      ▼ documentos/           │  neste diretório                  │
│        ▶ fotos/              │                                   │
│        ▶ textos/             │                                   │
│                              │                                   │
╰── Enter Abrir ────────────────────────────────────── Esc Cancelar ╯
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
| Indicador de pasta (`▶` recolhida, `▼` expandida) | `accent.secondary` | — |
| Arquivo selecionado no painel de arquivos | `special.highlight` (fundo) + `text.primary` | **bold** |
| Arquivo não selecionado | `text.primary` | — |
| Indicador de arquivo `●` | `text.secondary` | — |
| Metadados (tamanho, data relativa) | `text.secondary` | — |
| Texto `Nenhum arquivo .abditum` | `text.secondary` | — |
| Rótulo `Caminho:` | `text.secondary` | — |
| Valor do caminho | `text.primary` | — |
| Ação default (bloqueada) | `text.disabled` | — |
| Ação default (desbloqueada) | `accent.primary` | **bold** |

**Estados dos componentes:**

| Componente | Estado | Condição |
|---|---|---|
| Painel `Estrutura` (árvore) | sempre visível | — |
| Painel `Arquivos` (lista) | conteúdo visível | Pasta selecionada contém arquivos `.abditum` |
| Painel `Arquivos` (lista) | texto vazio | Pasta selecionada **não** contém arquivos `.abditum` |
| Rótulo `Caminho` | sempre visível, somente leitura | Atualiza ao navegar na árvore |
| Arquivo pré-selecionado no painel | selecionado | Primeiro `.abditum` da pasta, automaticamente ao entrar na pasta |
| Ação `Enter Abrir` | bloqueada (`text.disabled`) | Nenhum arquivo `.abditum` selecionado (pasta sem arquivos ou foco na árvore sem seleção à direita) |
| Ação `Enter Abrir` | ativa (`accent.primary` **bold**) | Um arquivo `.abditum` está selecionado no painel de arquivos |
| Ação `Esc Cancelar` | sempre ativa | — |

**Mensagens:**

| Contexto | Tipo | Texto |
|---|---|---|
| Diálogo abre / foco na árvore | Dica de campo | `• Navegue pelas pastas e selecione um cofre` |
| Foco no painel de arquivos (com seleção) | Dica de campo | `• Enter para abrir o cofre selecionado` |
| Foco no painel de arquivos (painel vazio) | Dica de campo | `• Nenhum cofre neste diretório — navegue para outra pasta` |
| Diálogo fecha | — | Barra limpa *(orquestrador assume)* |

**Comportamento:**

- `Tab` alterna entre árvore e painel de arquivos (2 stops)
- Na árvore: `↑↓` navega entre pastas; `→` ou `Enter` expande pasta recolhida; `←` recolhe pasta expandida; `Enter` em pasta expandida recolhe
- No painel de arquivos: `↑↓` navega entre arquivos; `Enter` confirma seleção (equivale à ação default)
- Ao expandir pasta na árvore, o painel de arquivos atualiza para mostrar os `.abditum` daquela pasta; primeiro arquivo pré-selecionado automaticamente
- Rótulo `Caminho` atualiza ao navegar na árvore

**Transições especiais:**

| Evento | Efeito |
|---|---|
| Navegar para pasta sem `.abditum` | Painel de arquivos mostra texto vazio; ação default muda para `text.disabled` |
| Navegar para pasta com `.abditum` | Primeiro arquivo pré-selecionado; ação default muda para `accent.primary` **bold** |
| `Enter` no painel de arquivos | Diálogo fecha com o arquivo selecionado |

---

#### FilePicker — Modo Save

**Título:** `Salvar cofre`
**Objetivo:** escolher diretório e nome para salvar o arquivo do cofre.

**Wireframe (campo nome preenchido — ação default ativa):**

```
╭── Salvar cofre ──────────────────────────────────────────────────╮
│  Caminho: /home/usuario/projetos/abditum                         │
├─ Estrutura ──────────────────┬─ Arquivos ────────────────────────┤
│  ▶ /                         │  ● database.abditum   25.8 MB 1h │
│    ▼ usuario/                │  ● config.abditum      1.2 KB 3d │
│      ▼ projetos/             │                                   │
│        ▼ abditum/            │                                   │
│          ▶ docs/             │                                   │
│                              │                                   │
├──────────────────────────────┴───────────────────────────────────┤
│  Nome do arquivo                                                 │
│  ░meu-cofre▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
╰── Enter Salvar ───────────────────────────────────── Esc Cancelar ╯
       ↑ accent.primary + bold (desbloqueado)
```

**Wireframe (campo nome vazio — ação default bloqueada):**

```
╭── Salvar cofre ──────────────────────────────────────────────────╮
│  Caminho: /home/usuario/projetos                                 │
├─ Estrutura ──────────────────┬─ Arquivos ────────────────────────┤
│  ▶ /                         │  ● database.abditum   25.8 MB 1h │
│    ▼ usuario/                │                                   │
│      ▼ projetos/             │                                   │
│                              │                                   │
├──────────────────────────────┴───────────────────────────────────┤
│  Nome do arquivo                                                 │
│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
╰── Enter Salvar ───────────────────────────────────── Esc Cancelar ╯
       ↑ text.disabled (bloqueado)
```

| Elemento | Token | Atributo |
|---|---|---|
| Título `Salvar cofre` | `text.primary` | **bold** |
| Header `Estrutura` | `text.secondary` | **bold** |
| Header `Arquivos` | `text.secondary` | **bold** |
| Separadores internos | `border.default` | — |
| Pasta selecionada na árvore | `accent.primary` | **bold** |
| Pasta não selecionada | `text.primary` | — |
| Indicador de pasta (`▶`/`▼`) | `accent.secondary` | — |
| Arquivo existente | `text.primary` | — |
| Indicador de arquivo `●` | `text.secondary` | — |
| Metadados | `text.secondary` | — |
| Rótulo `Caminho:` | `text.secondary` | — |
| Valor do caminho | `text.primary` | — |
| Label `Nome do arquivo` (campo ativo) | `accent.primary` | **bold** |
| Label `Nome do arquivo` (campo inativo) | `text.secondary` | — |
| Área do campo `░` | `surface.input` | — |
| Placeholder | `text.secondary` | *italic* |
| Cursor `▌` | `text.primary` | — |
| Ação default (bloqueada) | `text.disabled` | — |
| Ação default (desbloqueada) | `accent.primary` | **bold** |

**Estados dos componentes:**

| Componente | Estado | Condição |
|---|---|---|
| Painel `Estrutura` (árvore) | sempre visível | — |
| Painel `Arquivos` (lista) | conteúdo visível | Pasta selecionada contém arquivos `.abditum` |
| Painel `Arquivos` (lista) | texto vazio | Pasta selecionada **não** contém arquivos `.abditum` |
| Rótulo `Caminho` | sempre visível, somente leitura | Atualiza ao navegar na árvore |
| Campo `Nome do arquivo` | sempre visível | — |
| Extensão `.abditum` | adicionada automaticamente | Se o nome digitado não termina em `.abditum` |
| Ação `Enter Salvar` | bloqueada (`text.disabled`) | Campo `Nome do arquivo` vazio |
| Ação `Enter Salvar` | ativa (`accent.primary` **bold**) | Campo `Nome do arquivo` não vazio |
| Ação `Esc Cancelar` | sempre ativa | — |

> **Nota:** a validação de sobrescrita (arquivo já existe) é responsabilidade do fluxo que chamou o FilePicker, não do diálogo. O picker retorna o caminho completo; o fluxo abre diálogo de Confirmação × Destrutivo se necessário.

**Mensagens:**

| Contexto | Tipo | Texto |
|---|---|---|
| Diálogo abre / foco na árvore | Dica de campo | `• Navegue pelas pastas e escolha onde salvar` |
| Foco no painel de arquivos | Dica de campo | `• Arquivos existentes neste diretório` |
| Foco no campo `Nome do arquivo` (vazio) | Dica de campo | `• Digite o nome do arquivo — .abditum será adicionado automaticamente` |
| Foco no campo `Nome do arquivo` (preenchido) | Dica de campo | `• Enter para salvar o cofre` |
| Diálogo fecha | — | Barra limpa *(orquestrador assume)* |

**Comportamento:**

- `Tab` alterna entre árvore, painel de arquivos e campo de nome (3 stops)
- Navegação na árvore e painel de arquivos idêntica ao modo Open
- No painel de arquivos: selecionar um arquivo existente **copia o nome** para o campo `Nome do arquivo` (facilita sobrescrever)
- Ao navegar na árvore, o campo `Nome do arquivo` **não é limpo** — preserva o nome digitado
- Extensão `.abditum` é adicionada silenciosamente ao caminho de retorno, sem alterar o texto exibido no campo

**Transições especiais:**

| Evento | Efeito |
|---|---|
| Selecionar arquivo existente no painel | Nome copiado para campo `Nome do arquivo`; ação default muda para `accent.primary` **bold** |
| Limpar campo `Nome do arquivo` | Ação default volta para `text.disabled` |
| `Enter` com campo preenchido | Diálogo fecha com caminho completo (diretório + nome + `.abditum`) |

---

### Help

**Contexto de uso:** lista todas as ações do ActionManager, agrupadas. Acionado por `?` em qualquer contexto.
**Token de borda:** `border.default` (diálogo de consulta, não recebe entrada de texto)
**Dimensionamento:** largura máxima do DS; altura até 80% do terminal. Scroll quando o conteúdo excede a viewport.

**Wireframe (exemplo: Modo Cofre — segredo selecionado, sem scroll):**

```
╭── Ajuda — Atalhos e Ações ───────────────────────────────────────╮
│                                                                  │
│  Navegação                                                       │
│  ↑↓          Mover cursor na lista                               │
│  → / Enter   Expandir pasta ou selecionar segredo                │
│  ←           Recolher pasta ou subir para pasta pai              │
│  Tab         Alternar foco entre painéis                         │
│                                                                  │
│  Segredo                                                         │
│  Ctrl+R      Revelar / ocultar campo sensível                    │
│  Ctrl+C      Copiar valor para área de transferência             │
│  F21         Novo segredo                                        │
│  F22         Editar segredo                                      │
│  F23         Excluir segredo                                     │
│                                                                  │
│  Cofre                                                           │
│  ^S          Salvar cofre                                        │
│  F5          Sair (salva se necessário)                          │
│  ?           Esta ajuda                                          │
│                                                                  │
╰──────────────────────────────────────────────────── Esc Fechar ──╯
```

**Wireframe (exemplo: scroll — início do conteúdo, mais abaixo):**

```
╭── Ajuda — Atalhos e Ações ───────────────────────────────────────╮
│                                                                  ■
│  Navegação                                                       │
│  ↑↓          Mover cursor na lista                               │
│  → / Enter   Expandir pasta ou selecionar segredo                │
│  ←           Recolher pasta ou subir para pasta pai              │
│  Tab         Alternar foco entre painéis                         │
│                                                                  │
│  Segredo                                                         │
│  Ctrl+R      Revelar / ocultar campo sensível                    ↓
╰──────────────────────────────────────────────────── Esc Fechar ──╯
```

**Wireframe (exemplo: scroll — meio do conteúdo):**

```
╭── Ajuda — Atalhos e Ações ───────────────────────────────────────╮
│                                                                  ↑
│  Ctrl+C      Copiar valor para área de transferência             │
│  F21         Novo segredo                                        │
│  F22         Editar segredo                                      │
│  F23         Excluir segredo                                     ■
│                                                                  │
│  Cofre                                                           │
│  ^S          Salvar cofre                                        │
│  F5          Sair (salva se necessário)                          │
│  ?           Esta ajuda                                          ↓
╰──────────────────────────────────────────────────── Esc Fechar ──╯
```

> **Nota:** os wireframes são snapshots ilustrativos. O conteúdo real é gerado dinamicamente pelo ActionManager a partir do contexto ativo.

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Título `Ajuda — Atalhos e Ações` | `text.primary` | **bold** |
| Label do grupo (`Navegação`, `Segredo`, `Cofre`) | `text.secondary` | **bold** |
| Tecla (ex: `Ctrl+R`, `F21`, `^S`) | `accent.primary` | — |
| Descrição da ação | `text.primary` | — |
| Seta de scroll (`↑` / `↓` na borda direita) | `text.secondary` | — |
| Thumb de posição (`■` na borda direita) | `text.secondary` | — |
| Borda | `border.default` | — |

#### Estados dos componentes

| Componente | Estado | Condição |
|---|---|---|
| Conteúdo | sem scroll | Todas as ações cabem na viewport |
| Conteúdo | com scroll | Ações excedem a viewport — indicadores `↑`/`↓` e thumb `■` na borda direita (ver [DS — Scroll em diálogos](tui-design-system-novo.md#scroll-em-diálogos)) |
| `?` na barra de comandos | oculto (`HideFromBar`) | Enquanto o Help estiver aberto |
| Barra de comandos | vazia | Help não registra ações internas na barra |

#### Eventos

| Evento | Efeito |
|---|---|
| `?` pressionado (qualquer contexto) | Abre o modal; barra de comandos fica vazia; `?` oculto |
| `Esc` | Fecha o modal; `?` volta visível na barra |
| `↑` / `↓` | Scroll por linha (se conteúdo excede viewport) |
| `PgUp` / `PgDn` | Scroll por página (viewport − 1 linhas) |
| `Home` / `End` | Vai ao início / fim do conteúdo |

#### Comportamento

- **Conteúdo dinâmico** — gerado a partir de todas as ações registradas no ActionManager no momento da abertura
- **Agrupamento** — ações são organizadas pelo atributo numérico `Grupo`. Cada grupo tem um `Label` registrado (ex: 1 → "Navegação", 2 → "Segredo"). Grupos renderizados em ordem numérica crescente
- **Ordenação interna** — dentro de cada grupo, ações ordenadas por `Prioridade` (maior primeiro)
- **Scroll** — segue o padrão transversal do DS: indicadores `↑`/`↓` na borda direita, navegação por `↑↓` / `PgUp`/`PgDn` / `Home`/`End`
- **Borda inferior** — `Esc Fechar` sempre visível, independente do estado de scroll

---

## Componentes

### Cabeçalho

**Responsabilidade:** contexto global — qual aplicação, qual cofre, se há alterações pendentes e qual modo está ativo na área de trabalho.
**Posição:** linhas 1–2 da tela (zona Cabeçalho do [DS — Dimensionamento](tui-design-system-novo.md#dimensionamento-e-layout)).
**Altura fixa:** 2 linhas.

**Anatomia:**

| Linha | Papel | Conteúdo |
|---|---|---|
| **1 — Título** | Contexto + navegação | Nome da app, `·` separador, nome do cofre, `•` dirty, abas de modo à direita |
| **2 — Separadora** | Divisa cabeçalho ↔ área de trabalho | Linha `─` full-width; a aba ativa "pousa" nesta linha via `╯ Texto ╰` |

**Dois estados estruturais:**

| Estado | Linha 1 | Linha 2 | Abas |
|---|---|---|---|
| Sem cofre (boas-vindas) | Apenas nome da app | Separador simples, sem conectores | Ocultas |
| Cofre aberto | Nome da app `·` cofre `•` + abas | Separador com aba ativa suspensa | Visíveis (3) |

---

#### Sem cofre (Boas-vindas)

> Wireframes ilustrativos a 80 colunas. A largura real acompanha o terminal.

```
  Abditum
──────────────────────────────────────────────────────────────────────────────────
```

Sem nome de cofre, sem indicador dirty, sem abas. A linha separadora é contínua.

---

#### Cofre aberto — anatomia base

> Estado impossível em operação normal (sempre há um modo ativo). Mostrado para ilustrar a posição de todos os elementos antes de qualquer aba estar ativa.

**Sem alterações:**

```
  Abditum · cofre                          ╭ Cofre ╮  ╭ Modelos ╮  ╭ Config ╮
──────────────────────────────────────────────────────────────────────────────────
```

**Com alterações não salvas:**

```
  Abditum · cofre •                         ╭ Cofre ╮  ╭ Modelos ╮  ╭ Config ╮
──────────────────────────────────────────────────────────────────────────────────
```

O `•` aparece imediatamente após o nome do cofre, em `semantic.warning`. Desaparece após salvamento bem-sucedido.

---

#### Modo Cofre ativo

```
  Abditum · cofre •                      ╭───────╮  ╭ Modelos ╮  ╭ Config ╮
─────────────────────────────────────────╯ Cofre ╰──────────────────────────────
```

A aba ativa na linha 1 substitui o texto por `─` (`╭───────╮`), mantendo a mesma largura da versão inativa (`╭ Cofre ╮`). Na linha 2, o texto desce para o gap entre `╯` e `╰`, que se alinham verticalmente com `╭` e `╮` da linha 1 — conectando visualmente a aba à área de trabalho abaixo.

---

#### Modo Modelos ativo

```
  Abditum · cofre                          ╭ Cofre ╮  ╭─────────╮  ╭ Config ╮
──────────────────────────────────────────────────────╯ Modelos ╰────────────────
```

---

#### Modo Configurações ativo

```
  Abditum · cofre                           ╭ Cofre ╮  ╭ Modelos ╮  ╭────────╮
────────────────────────────────────────────────────────────────────╯ Config ╰──
```

A aba mais à direita pode encostar na borda do terminal — `╰` ocupa a última coluna, sem `─` posterior.

> **Nota:** a aba Configurações é referida como "Config" nos wireframes por economia de espaço. O texto completo na implementação é `Config`.

---

#### Mecânica visual da aba ativa

A transformação de aba inativa → ativa ocorre em duas linhas simultâneas:

| Linha | Aba inativa | Aba ativa |
|---|---|---|
| **1** | `╭ Texto ╮` (borda + texto) | `╭──────╮` (borda + preenchimento `─`) |
| **2** | `─────────` (separador contínuo) | `╯ Texto ╰` (gap com texto sobre `special.highlight`) |

Regras de alinhamento:

- A largura total da aba é **idêntica** nos estados ativo e inativo
- `╯` alinha-se verticalmente com `╭` da linha acima
- `╰` alinha-se verticalmente com `╮` da linha acima
- O conteúdo entre `╯` e `╰` (espaço + texto + espaço) tem fundo `special.highlight`
- As bordas `╭╮╯╰` e o preenchimento `─` usam sempre `border.default`, independente do estado

---

#### Truncamento do nome do cofre

O espaço disponível para o nome do cofre é limitado — as abas ocupam largura fixa à direita. O componente calcula o espaço em tempo real.

> **Extensão `.abditum` é omitida** — a app só trabalha com este formato, então a extensão é redundante. O nome exibido é o radical do arquivo (ex: `cofre.abditum` → `cofre`).

**Fórmula:**

```
prefixo  = "  Abditum · "                             (12 colunas)
dirty    = " •"  se IsDirty(), ou ""                   (2 ou 0 colunas)
abas     = bloco de abas + espaços entre elas           (largura fixa, ~32 colunas)
padding  = mín. 1 coluna entre nome/dirty e abas

disponível = largura_terminal − prefixo − dirty − abas − padding
```

**Algoritmo:**

1. Se o nome completo (radical sem extensão) cabe → exibir como está
2. Se não cabe → truncar com `…`: `{nome[0..n]}…` onde `n` é calculado para caber
3. Se nem 1 caractere + `…` (2 colunas) cabe → exibir apenas `…`

**Prioridade de cessão de espaço:**

| Prioridade | Elemento | Comportamento |
|---|---|---|
| 1ª (cede primeiro) | Nome do cofre | Truncado conforme algoritmo acima |
| 2ª | Separador `·` e indicador `•` | Preservados enquanto houver espaço |
| 3ª (nunca cede) | Abas | Largura fixa, nunca truncadas |

**Wireframe — nome truncado (terminal ~80 colunas, modo Cofre):**

```
  Abditum · meu-cofre-pessoa… •          ╭───────╮  ╭ Modelos ╮  ╭ Config ╮
─────────────────────────────────────────╯ Cofre ╰──────────────────────────────
```

O radical `meu-cofre-pessoal` foi truncado para `meu-cofre-pessoa…`.

---

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| `Abditum` (nome da app) | `accent.primary` | **bold** |
| `·` separador nome/cofre | `border.default` | — |
| Nome do cofre (radical, sem `.abditum`) | `text.secondary` | — |
| `•` indicador não salvo | `semantic.warning` | — |
| Bordas das abas (`╭╮╯╰─`) — ativa e inativa | `border.default` | — |
| Aba ativa — fundo (gap entre `╯` e `╰`) | `special.highlight` | — |
| Aba ativa — texto | `accent.primary` | **bold** |
| Aba inativa — texto | `text.secondary` | — |
| Linha separadora | `border.default` | — |

---

#### Eventos

| Evento | Mudança visual |
|---|---|
| Cofre aberto com sucesso | Aparece `·` nome do cofre e as 3 abas |
| Cofre fechado / bloqueado | Desaparece nome do cofre e abas; volta ao estado boas-vindas |
| Alteração em memória (`IsDirty() = true`) | Aparece `•` em `semantic.warning` |
| Salvamento bem-sucedido (`IsDirty() = false`) | Desaparece `•` |
| Navegação entre modos (`F1` Cofre / `F2` Modelos / `F3` Config, ou clique) | Aba ativa muda; nova aba suspensa na linha separadora |
| Terminal redimensionado | Nome do cofre recalcula truncamento |

---

#### Comportamento

- **Abas clicáveis** — mouse troca o modo ativo ao clicar no texto ou na borda da aba (área de hit inclui linhas 1 e 2 da aba)
- **Navegação por teclado** — `F1` Cofre, `F2` Modelos, `F3` Config (escopo Área de trabalho — só ativas com cofre aberto)
- **Indicador dirty** — aparece/desaparece imediatamente conforme `IsDirty()`, sem animação
- **Truncamento dinâmico** — recalculado a cada renderização (resize do terminal, mudança de modo ativo, cofre aberto/fechado)

---

### Barra de Comandos

**Responsabilidade:** exibir as ações disponíveis no contexto ativo — o usuário nunca precisa adivinhar o que pode fazer.
**Posição:** última linha da tela (zona Barra de comandos do [DS — Dimensionamento](tui-design-system-novo.md#dimensionamento-e-layout)).
**Altura fixa:** 1 linha.

**Princípio de conteúdo:** a barra exibe apenas ações de caso de uso (F-keys, atalhos de domínio, `^S`). Teclas de navegação universais — `↑↓`, `←→`, `Tab`, `Enter`, `Esc` — são senso comum em TUI e não são exibidas. Exceção: diálogos exibem ações internas específicas do contexto.

---

#### Anatomia

Cada ação na barra segue o formato: **TECLA Label** — tecla em `accent.primary` **bold**, label em `text.primary`. Ações separadas por `·` em `text.secondary`. A ação `?` (Ajuda) é âncora fixa na extrema direita.

**Estado normal:**

```
  F21 Novo · F22 Editar · F23 Excluir · ^S Salvar                                   ?
```

**Com ação desabilitada (nenhum segredo selecionado):**

```
  F21 Novo · F22 Editar · F23 Excluir · ^S Salvar                                   ?
```

`F23 Excluir` em `text.disabled` + dim. Permanece visível na posição — não colapsa.

**Durante diálogo ativo (apenas ações internas):**

```
  Tab Campos · F5 Revelar                                                            ?
```

Ações do ActionManager ficam ocultas. A barra mostra apenas as ações internas do diálogo do topo da pilha. Ações de confirmação/cancelamento (`Enter`/`Esc`) já estão na borda do diálogo — não são duplicadas na barra.

**Espaço restrito:**

```
  F21 Novo                                                                           ?
```

Ações de menor prioridade são ocultadas quando não há espaço. `?` permanece sempre visível — é via Help que o usuário descobre as ações ocultas.

---

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Tecla da ação (ex: `F21`) | `accent.primary` | **bold** |
| Label da ação (ex: `Novo`) | `text.primary` | — |

| Separador `·` | `text.secondary` | — |
| `?` (Ajuda) | `accent.primary` | **bold** |

---

#### Atributos das ações

Cada ação registrada no ActionManager possui atributos que controlam sua apresentação:

| Atributo | Efeito na barra | Efeito no Help |
|---|---|---|
| `Enabled = true` | Exibida com estilo normal | Listada |
| `Enabled = false` | **Não aparece** na barra | Listada |
| `HideFromBar = true` | **Não aparece** na barra | Listada |
| `HideFromBar = false` | Exibida (se `Enabled`) | Listada |

Além destes:

- **Prioridade** — valor numérico. Maior prioridade → mais à esquerda na barra. Quando o espaço é insuficiente, ações de menor prioridade são removidas primeiro
- **Grupo** — valor numérico. Usado exclusivamente no modal de Ajuda para organizar ações. Grupos renderizados em ordem numérica crescente. Dentro de cada grupo, ações ordenadas por `Prioridade`. Não afeta a barra de comandos
- **Label do grupo** — string registrada por grupo (ex: grupo 1 → "Navegação"). Exibido como título de seção no Help em `text.secondary` bold

---

#### Eventos

| Evento | Mudança na barra |
|---|---|
| Troca de foco entre painéis (`Tab` / `Shift+Tab`) | Ações do painel que recebe foco ficam ativas |
| Seleção de item na árvore | Ações de item (editar, excluir, revelar) ficam `Enabled = true` — aparecem na barra |
| Nenhum item selecionado | Ações de item ficam `Enabled = false` — desaparecem da barra |
| Diálogo aberto (push na pilha) | Troca para ações internas do diálogo |
| Diálogo fechado (pop da pilha) | Volta para ações do ActionManager |
| Terminal redimensionado | Recalcula quais ações cabem (prioridade governa corte) |

---

#### Comportamento

- **Âncora `?`** — reserva espaço fixo na extrema direita. O cálculo de espaço disponível desconta `?` antes de distribuir as demais ações
- **Ações desabilitadas desaparecem da barra** — `Enabled = false` remove a ação da barra (não fica exibida como dim). A ação continua listada no Help
- **Diálogos de decisão** (confirmação/reconhecimento) — tipicamente não têm ações internas; a barra pode ficar vazia (apenas `?`) enquanto o diálogo estiver aberto
- **Diálogos funcionais** (PasswordEntry, FilePicker etc.) — registram ações internas (Tab entre campos, revelar senha, etc.) que aparecem na barra
- **Truncamento** — se mesmo a ação de maior prioridade + `?` não cabem, a barra mostra apenas `?`

---

### Barra de Mensagens

**Responsabilidade:** comunicar feedback ao usuário — sucesso, erro, aviso, progresso, dicas.
**Posição:** 1 linha fixa entre a área de trabalho e a barra de comandos (zona Barra de mensagens do [DS — Dimensionamento](tui-design-system-novo.md#dimensionamento-e-layout)).
**Altura fixa:** 1 linha.
**Anatomia:** borda `─` contínua na largura total do terminal. Quando há mensagem, o texto (símbolo + conteúdo) começa com 2 espaços de padding à esquerda (alinhado com o texto do cabeçalho), seguido de `─` até o fim da linha.

**Wireframe (sem mensagem — borda separadora):**

```
────────────────────────────────────────────────────────────────────────────────
```

**Wireframe (sucesso):**

```
── ✓ Gmail copiado para a área de transferência ────────────────────────────────
   ↑ semantic.success
```

**Wireframe (erro):**

```
── ✗ Falha ao salvar — arquivo em uso por outro processo ───────────────────────
   ↑ semantic.error + bold
```

**Wireframe (aviso):**

```
── ⚠ Cofre será bloqueado em 15 segundos ──────────────────────────────────────
   ↑ semantic.warning
```

**Wireframe (spinner):**

```
── ◐ Salvando cofre... ─────────────────────────────────────────────────────────
   ↑ accent.primary
```

**Wireframe (dica de campo / dica de uso):**

```
── • Use Tab para alternar o foco entre os painéis ─────────────────────────────
   ↑ text.secondary + italic
```

**Wireframe (informação):**

```
── ℹ Cofre criado em /home/user/documentos/pessoal.abditum ─────────────────────
   ↑ semantic.info
```

**Wireframe (truncamento — mensagem excede largura disponível):**

```
── ✗ Erro ao importar arquivo: o formato do arquivo não é compatível com a v… ──
   ↑ semantic.error + bold                                            ↑ trunca com …
```

#### Tokens

Os tokens de cada tipo de mensagem são definidos no [DS — Mensagens](tui-design-system-novo.md#mensagens). Adicional:

| Elemento | Token | Atributo |
|---|---|---|
| Borda `─` (sem mensagem) | `border.default` | — |
| Borda `─` (com mensagem) | `border.default` | — |

> A cor da borda não muda conforme o tipo de mensagem — apenas o texto embutido usa o token semântico correspondente.

#### Estados dos componentes

| Componente | Estado | Condição |
|---|---|---|
| Borda | visível (sem texto) | Nenhuma mensagem ativa |
| Borda + mensagem | visível (texto embutido) | Mensagem ativa — tipo governa símbolo, cor e atributo |
| Texto | truncado com `…` | Mensagem excede ~95% da largura do terminal |

#### Eventos

| Evento | Efeito |
|---|---|
| Operação concluída com sucesso | Exibe `✓` mensagem (`semantic.success`, TTL 2–3s) |
| Informação neutra | Exibe `ℹ` mensagem (`semantic.info`, TTL 3s) |
| Condição de alerta (ex: bloqueio iminente) | Exibe `⚠` mensagem (`semantic.warning`, permanente, desaparece com input) |
| Falha em operação | Exibe `✗` mensagem (`semantic.error` + bold, TTL 5s) |
| Operação em andamento | Exibe spinner `◐◓◑◒` (`accent.primary`, permanente até sucesso/erro) |
| Campo recebe foco (diálogo funcional) | Exibe `•` dica de campo (`text.secondary` italic) |
| Aplicação emite dica proativa | Exibe `•` dica de uso (`text.secondary` italic) |
| TTL expira | Mensagem desaparece — volta à borda `─` |
| Nova mensagem emitida | Substitui imediatamente a mensagem anterior |
| Diálogo fecha | Barra é limpa — volta à borda `─` |

#### Comportamento

- **Borda permanente** — a borda `─` é sempre visível, funcionando como separador entre a área de trabalho e a barra de comandos. Contribui para a estabilidade espacial
- **Uma mensagem por vez** — nova mensagem substitui a anterior imediatamente. Não há fila nem pilha
- **Texto embutido** — o texto (símbolo + conteúdo) substitui o trecho central da borda, com `─` preenchendo os lados
- **Aviso re-emitido** — mensagens de aviso são re-emitidas a cada tick enquanto a condição persistir
- **Responsabilidade do orquestrador** — mensagens pós-fechamento de diálogo (ex: "✓ Cofre aberto") são emitidas pelo orquestrador, não pelo diálogo

---

<!-- SEÇÕES FUTURAS — a preencher pela equipe -->

<!--
## Telas

### Boas-vindas
### Modo Cofre
### Modo Modelos
### Modo Configurações

## Componentes

### Painel Esquerdo: Árvore
### Painel Direito: Detalhe do Segredo
### Painel Direito: Detalhe do Modelo

## Fluxos Visuais

### Criar cofre
### Abrir cofre
### Salvar cofre
### Bloquear cofre
### Alterar senha mestra
### Criar segredo
### Editar segredo
### Excluir segredo
### Buscar segredo
### Exportar cofre
### Importar cofre
-->
