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
  - [Matriz Intenção × Severidade](#matriz-intenção--severidade)
  - [Exemplo: Confirmação × Destrutivo — Excluir segredo](#exemplo-confirmação--destrutivo--excluir-segredo)
  - [Exemplo: Confirmação × Neutro — Alterações não salvas](#exemplo-confirmação--neutro--alterações-não-salvas)
  - [Exemplo: Reconhecimento × Informativo — Conflito detectado](#exemplo-reconhecimento--informativo--conflito-detectado)
  - [Exemplo: Reconhecimento × Erro — Falha ao abrir cofre](#exemplo-reconhecimento--erro--falha-ao-abrir-cofre)
- [Diálogos Funcionais](#diálogos-funcionais)
  - [PasswordEntry](#passwordentry)
  - [PasswordCreate](#passwordcreate)
  - [FilePicker](#filepicker)
  - [Help](#help)

---

## Diálogos de Decisão

Todos os diálogos de decisão seguem a anatomia comum e o modelo bidimensional (Intenção × Severidade) definidos no [design system — Sobreposição](tui-design-system-novo.md#sobreposição). Esta seção define a matriz genérica e apresenta wireframes ilustrativos de combinações típicas.

---

### Matriz Intenção × Severidade

A tabela abaixo mostra como cada combinação se manifesta visualmente. Nem todas as combinações são igualmente comuns — as marcadas com **★** aparecem nos wireframes illustrativos abaixo.

| | Destrutivo (`⚠`) | Erro (`✗`) | Alerta (`⚠`) | Informativo (`ℹ`) | Neutro (—) |
|---|---|---|---|---|---|
| **Confirmação** | ★ Excluir item, sobrescrever arquivo | Raro (tentar novamente?) | Ação com efeito colateral | Confirmar operação longa | ★ Salvar/descartar/voltar |
| **Reconhecimento** | Ação executada e irreversível | ★ Falha ao abrir cofre | Operação parcialmente bem-sucedida | ★ Conflito detectado | Operação concluída |

**Como ler a matriz:** cada célula descreve o *cenário típico* — não o conteúdo literal do diálogo. O tratamento visual (borda, símbolo, cor da tecla) é determinado exclusivamente pela severidade; a intenção define a barra de ações (múltiplas opções vs. apenas OK).

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

```
╭── Senha mestra ─────────────────╮
│                                  │
│  Senha                           │  ← label em text.secondary
│  ••••••••••▌                     │  ← texto mascarado + cursor
│                                  │
│  Tentativa 2 de 5                │  ← contador em text.secondary
╰── Enter Confirmar ──────── Esc Cancelar ╯
```

| Elemento | Token | Atributo |
|---|---|---|
| Título `Senha mestra` | `text.primary` | **bold** |
| Label `Senha` | `text.secondary` | — |
| Máscara `••••••••••` | `text.secondary` | — |
| Cursor `▌` | `text.primary` | — |
| Contador `Tentativa 2 de 5` | `text.secondary` | — |
| Fundo do campo | `surface.input` | — |

**Comportamento:**
- Máscara de comprimento fixo (8 `•`) — não revela o tamanho real da senha
- Contador aparece a partir da segunda tentativa
- Campo único — `Tab` não faz nada dentro deste diálogo

---

### PasswordCreate

**Contexto de uso:** criação de senha (ao criar cofre ou alterar senha mestra).
**Token de borda:** `border.focused`

```
╭── Definir senha mestra ─────────╮
│                                  │
│  Nova senha                      │  ← label em accent.primary (campo ativo)
│  ••••••••••▌                     │  ← texto mascarado + cursor
│                                  │
│  Confirmação                     │  ← label em text.secondary (campo inativo)
│                                  │
│  Força: ████████░░ Boa           │  ← medidor de força
╰── Enter Confirmar ──────── Esc Cancelar ╯
```

| Elemento | Token | Atributo |
|---|---|---|
| Título `Definir senha mestra` | `text.primary` | **bold** |
| Label do campo ativo | `accent.primary` | **bold** |
| Label do campo inativo | `text.secondary` | — |
| Máscara | `text.secondary` | — |
| Cursor `▌` | `text.primary` | — |
| Medidor — preenchido | `semantic.success` ou `semantic.warning` | — |
| Medidor — vazio | `text.disabled` | — |
| Label de força `Boa` / `Forte` | `semantic.success` | — |
| Label de força `Fraca` | `semantic.warning` | — |
| Fundo dos campos | `surface.input` | — |

**Comportamento:**
- `Tab` alterna entre os campos `Nova senha` e `Confirmação`
- Medidor de força atualizado a cada tecla no campo `Nova senha`
- Se as senhas não conferem ao confirmar, borda do campo Confirmação muda para `semantic.error`

---

### FilePicker

**Contexto de uso:** abrir ou salvar arquivo do cofre.
**Token de borda:** `border.focused`

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
╰── Enter Selecionar ─────── Esc Cancelar ╯
```

| Elemento | Token | Atributo |
|---|---|---|
| Título (ex: `Abrir cofre`) | `text.primary` | **bold** |
| Separadores internos | `border.default` | — |
| Header de coluna (`Diretórios`, caminho) | `text.secondary` | — |
| Pasta selecionada (`► cofres/`) | `accent.primary` | **bold** |
| Pastas normais | `text.primary` | — |
| Nomes de pastas (ícone) | `accent.secondary` | — |
| Arquivos `.abditum` | `text.primary` | — |
| Label `Nome do arquivo` | `text.secondary` | — |
| Campo de nome | `surface.input` | — |

**Comportamento:**
- Painel esquerdo: árvore de diretórios navegável com `↑↓←→`
- Painel direito: lista de arquivos no diretório selecionado
- `Tab` alterna entre painel de diretórios, lista de arquivos e campo de nome (modo save)
- No modo open: campo de nome ausente; `Enter` sobre arquivo seleciona-o
- No modo save: campo de nome presente; `.abditum` adicionado automaticamente se omitido

---

### Help

**Contexto de uso:** tabela de atalhos do contexto ativo.
**Token de borda:** `border.default` (diálogo de consulta, não recebe entrada de texto)

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

| Elemento | Token | Atributo |
|---|---|---|
| Título `Ajuda — Atalhos e Ações` | `text.primary` | **bold** |
| Nome do grupo (`Navegação`, `Segredo`) | `text.secondary` | **bold** |
| Tecla (ex: `F16`) | `accent.primary` | — |
| Descrição da ação | `text.primary` | — |
| Indicador de scroll `↓ mais` | `text.secondary` | — |

**Comportamento:**
- Conteúdo gerado dinamicamente a partir dos atalhos do contexto ativo
- `↑↓` para scroll se o conteúdo exceder a área visível
- `Esc` fecha o modal

---

<!-- SEÇÕES FUTURAS — a preencher pela equipe -->

<!--
## Telas

### Boas-vindas
### Modo Cofre
### Modo Modelos
### Modo Configurações

## Componentes

### Cabeçalho
### Barra de Comandos
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
