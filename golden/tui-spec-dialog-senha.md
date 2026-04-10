# Especificação Visual — Diálogos de Senha

> PasswordEntry e PasswordCreate.
> Parte da [Especificação Visual](tui-specification.md).
>
> **Documentos relacionados:**
> - [`tui-design-system.md`](tui-design-system.md) — fundações visuais
> - [`tui-spec-dialogos.md`](tui-spec-dialogos.md) — anatomia comum e tipos de diálogo

## PasswordEntry

**Contexto de uso:** entrada de senha para abrir cofre.
**Token de borda:** `border.focused`
**Largura fixa:** 50 colunas

### Anatomia

> Nos wireframes abaixo, `░` representa a área com fundo `surface.input` (campo de entrada). Na implementação real, o campo é uma área de fundo rebaixado sem hachura — conforme definido em [Campos de entrada de texto](tui-design-system.md#foco-e-navegação).

**Estado inicial (campo vazio — ação default bloqueada):**

```
╭── Senha mestra ────────────────────────────╮
│                                            │
│  Senha                                     │
│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
│                                            │
╰── Enter Confirmar ───────────── Esc Cancelar ──╯
       ↑ text.disabled (bloqueado)
```

**Estado com digitação (ação default ativa):**

```
╭── Senha mestra ────────────────────────────╮
│                                            │
│  Senha                                     │
│  ░••••••••▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  │
│                                            │
╰── Enter Confirmar ───────── Esc Cancelar ──╯
       ↑ accent.primary + bold (desbloqueado)
```

**Estado com contador de tentativas (a partir da 2ª):**

```
╭── Senha mestra ────────────────────────────╮
│                                            │
│  Senha                                     │
│  ░••••••••▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  │
│                                            │
│  Tentativa 2 de 5                          │
╰── Enter Confirmar ───────── Esc Cancelar ──╯
```

### Identidade Visual

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

### Estados

| Componente | Estado | Condição |
|---|---|---|
| Campo `Senha` | sempre visível, sempre com foco | Diálogo de campo único |
| Contador de tentativas | visível | Tentativa atual ≥ 2 |
| Contador de tentativas | oculto | Primeira tentativa |
| Ação `Enter Confirmar` | bloqueada (`text.disabled`) | Campo `Senha` vazio |
| Ação `Enter Confirmar` | ativa (`accent.primary` **bold**) | Campo `Senha` não vazio |
| Ação `Esc Cancelar` | sempre ativa | — |

### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Diálogo abre / foco no campo (vazio ou válido) | Dica de campo | `• Digite a senha para desbloquear o cofre` |
| `Enter` → senha incorreta | Erro (5s) | `✕ Senha incorreta` |
| Diálogo fecha (confirmação ou cancelamento) | — | Barra limpa *(orquestrador assume)* |

### Eventos

| Evento | Efeito |
|---|---|
| `Enter` com senha incorreta | Campo limpo; ação default volta para `text.disabled`; contador incrementado |
| Tentativas esgotadas | Diálogo fecha automaticamente |

### Comportamento

- Máscara de comprimento fixo (8 `•`) — não revela o tamanho real da senha
- Campo único — `Tab` não faz nada dentro deste diálogo

## PasswordCreate

**Contexto de uso:** criação de senha (ao criar cofre ou alterar senha mestra).
**Token de borda:** `border.focused`
**Largura fixa:** 50 colunas

### Anatomia

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
╰── Enter Confirmar ───────────────── Esc Cancelar ──╯
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
╰── Enter Confirmar ───────────────── Esc Cancelar ──╯
       ↑ text.disabled (2º campo vazio)
```

**Estado com ambos campos preenchidos e senhas conferem (ação default desbloqueada):**

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
╰── Enter Confirmar ───────────────── Esc Cancelar ──╯
       ↑ accent.primary + bold (desbloqueado)
```

**Estado com senhas divergentes (ação default bloqueada — erro no campo):**

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
╰── Enter Confirmar ───────────────── Esc Cancelar ──╯
       ↑ text.disabled (senhas divergem)
```

### Identidade Visual

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

### Estados

| Componente | Estado | Condição |
|---|---|---|
| Campo `Nova senha` | sempre visível | — |
| Campo `Confirmação` | sempre visível | — |
| Medidor de força | visível | Campo `Nova senha` não vazio |
| Medidor de força | oculto | Campo `Nova senha` vazio |
| Linha em branco antes do medidor | visível | Medidor visível |
| Ação `Enter Confirmar` | bloqueada (`text.disabled`) | Campo `Nova senha` vazio **ou** campo `Confirmação` vazio **ou** senhas divergentes |
| Ação `Enter Confirmar` | ativa (`accent.primary` **bold**) | Ambos os campos não vazios **e** senhas conferem |
| Ação `Esc Cancelar` | sempre ativa | — |

> **Nota:** a verificação de igualdade entre as senhas ocorre **em tempo real** — a cada tecla no campo `Confirmação` e ao abandonar o campo (Tab ou mudança de foco). Se as senhas divergem, a ação default fica bloqueada e a barra de mensagens exibe erro no lugar da dica de campo.

### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Diálogo abre / foco em `Nova senha` (vazio ou válido) | Dica de campo | `• A senha mestra protege todo o cofre — use 12+ caracteres` |
| Foco em `Confirmação` (vazio ou válido) | Dica de campo | `• Redigite a senha para confirmar` |
| Foco em `Confirmação` (senhas divergentes) | Erro (5s) | `✕ As senhas não conferem — digite novamente` |
| Digitação em `Confirmação` (senhas divergentes) | Erro (5s) | `✕ As senhas não conferem — digite novamente` |
| `Enter` → senhas divergentes | Erro (5s) | `✕ As senhas não conferem — digite novamente` |
| Diálogo fecha (confirmação ou cancelamento) | — | Barra limpa *(orquestrador assume)* |

### Teclado

| Tecla | Efeito | Condição |
|---|---|---|
| `Tab` | Alterna foco entre `Nova senha` e `Confirmação` | — |
| `Enter` | Confirma criação da senha | Ambos preenchidos e senhas conferem |
| `Esc` | Cancela o diálogo | — |

### Eventos

| Evento | Efeito |
|---|---|
| Digitação em `Confirmação` torna senhas iguais | Erro na barra é substituído pela dica de campo; ação default muda para `accent.primary` **bold** |
| Digitação em `Confirmação` torna senhas diferentes | Dica de campo é substituída por erro (`✕`, TTL 5s); ação default volta para `text.disabled` |
| Abandonar `Confirmação` (Tab) com senhas divergentes | Erro exibido na barra; foco move para `Nova senha`; ação default bloqueada |
| Abandonar `Confirmação` (Tab) com senhas iguais | Dica exibida na barra; foco move para `Nova senha`; ação default ativa |

### Comportamento

- `Tab` alterna entre os campos `Nova senha` e `Confirmação`
- Medidor de força atualizado a cada tecla no campo `Nova senha`
- Máscara de comprimento fixo (8 `•`) — não revela o tamanho real da senha
- Validação de igualdade em tempo real: a cada tecla no campo `Confirmação` e ao abandonar o campo (Tab)
- Senhas divergentes: ação default bloqueada (`text.disabled`); barra de mensagens exibe erro (`✕`) no lugar da dica de campo; erro permanece até que as senhas confiram ou o campo seja limpo
