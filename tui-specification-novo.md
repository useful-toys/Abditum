# Especificação Visual — Abditum TUI

> Wireframes, layouts de componentes e fluxos visuais concretos.
> Cada tela e componente consome os padrões definidos no design system.
>
> **Documento de fundação:**
> - [`tui-design-system-novo.md`](tui-design-system-novo.md) — princípios, tokens, estados, padrões transversais

## Sumário

- [Atalhos da Aplicação](#atalhos-da-aplicação)
- [Diálogos de Decisão](#diálogos-de-decisão)
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
  - [Painel Esquerdo: Árvore](#painel-esquerdo-árvore)
  - [Busca de Segredos](#busca-de-segredos)
  - [Painel Direito: Detalhe do Segredo — Modo Leitura](#painel-direito-detalhe-do-segredo--modo-leitura)
- [Ações na Árvore de Segredos](#ações-na-árvore-de-segredos)
  - [⌃R e ⌃C na árvore — Atalhos de campo sensível](#r-e-c-na-árvore--atalhos-de-campo-sensível)
- [Telas](#telas)
  - [Boas-vindas](#boas-vindas)

---

## Atalhos da Aplicação

Este documento detalha as atribuições específicas de teclas para os fluxos e funções da aplicação. As políticas transversais de teclado e o agrupamento de teclas F por categoria de ação são definidos no [Design System — Mapa de Teclas](tui-design-system-novo.md#mapa-de-teclas).

### Atalhos Globais

| Tecla | Ação (Fluxo ou Função) | Escopo | Notas |
|---|---|---|---|
| `F1` | Abrir / fechar modal de Ajuda | Global | |
| `F12` | Alternar Tema | Global | |
| `Ctrl+Q` | Sair da Aplicação (Fluxos 3, 4, 5) | Global | Gerencia todas as saídas com as devidas confirmações |
| `Ctrl+Alt+Shift+Q` | Bloquear Cofre (Fluxo 6) | Global | Bloqueio emergencial, descarta alterações, sem confirmação. Atalho "complicado" para evitar acidentes. |

### Atalhos de Área de Trabalho (Fluxos Principais)

Os seguintes atalhos disparam os fluxos principais da aplicação quando a área de trabalho tem foco (sem diálogos abertos). Eles seguem os agrupamentos de teclas F definidos no Design System.

| Tecla | Ação (Fluxo) | Notas |
|---|---|---|
| `F2` | Modo Cofre (aba) | Só com cofre aberto |
| `F3` | Modo Modelos (aba) | Só com cofre aberto |
| `F4` | Modo Configurações (aba) | Abrange o Fluxo 14: Configurar o Cofre |
| `F5` | Criar Novo Cofre (Fluxo 2) | |
| `F6` | Abrir Cofre Existente (Fluxo 1) | |
| `Shift+F6` | Descartar Alterações e Recarregar Cofre (Fluxo 10) | Similaridade semântica com F6 |
| `F7` | Salvar Cofre no Arquivo Atual (Fluxo 8) | |
| `Shift+F7` | Salvar Cofre em Outro Arquivo (Fluxo 9) | |
| `Ctrl+F7` | Alterar Senha Mestra (Fluxo 11) | Implica salvamento |
| `F8` | (Livre) | Reservado para futuras ações de persistência |
| `F9` | Exportar Cofre (Fluxo 12) | |
| `Shift+F9` | Importar Cofre (Fluxo 13) | |
| `F10` | Busca de Segredos — abrir/fechar campo | Só com cofre aberto e foco na árvore; toggle |
| `F11` | (Livre) | |

> **Fluxo 7 — Aviso de Bloqueio Iminente por Inatividade:** É um fluxo iniciado pelo sistema, não requer um atalho manual do usuário.

---

## Diálogos de Decisão

Todos os diálogos de decisão seguem a anatomia comum e os padrões de interação definidos no [design system — Sobreposição](tui-design-system-novo.md#sobreposição), incluindo a [Referência Visual por Severidade](tui-design-system-novo.md#severidade) e as [Regras de Ações na Borda Inferior](tui-design-system-novo.md#ações-na-borda-inferior).

---

## Catálogo de Diálogos de Decisão

Esta seção lista todas as instâncias de diálogos de decisão da aplicação, especificando seu contexto, título, mensagem no corpo e ações na borda. A estrutura visual é definida na seção [Sobreposição](tui-design-system-novo.md#sobreposição) do Design System.

| Ação | Situação | Tipo | Título | Mensagem no Corpo | Ações na Borda |
|---|---|---|---|---|---|
| **Sair** | Sem alterações | Confirmação × Neutro | `Sair do Abditum` | `Sair do Abditum?` | `Enter Sair`, `Esc Voltar` |
| **Sair** | Com alterações | Confirmação × Alerta | `Sair do Abditum` | `Cofre modificado. Salvar ou descartar?` | `S Salvar`, `D Descartar`, `Esc Voltar` |
| **Salvar** | Conflito externo | Confirmação × Destrutivo | `Salvar cofre` | `Arquivo modificado externamente. Sobrescrever?` | `S Sobrescrever`, `Esc Voltar` |
| **Abrir cofre** | Falha (arquivo inválido) | Reconhecimento × Erro | `Abrir cofre` | `Arquivo corrompido ou inválido. Necessário fechar.` | `Enter OK` |
| **Abrir cofre** | Modificações não salvas | Confirmação × Alerta | `Abrir cofre` | `Cofre modificado. Salvar ou descartar?` | `S Salvar`, `D Descartar`, `Esc Voltar` |
| **Abrir cofre** | Caminho/Formato inválido | Reconhecimento × Erro | `Abrir cofre` | `Arquivo inválido ou versão não suportada. Necessário corrigir.` | `Enter OK` |
| **Abrir cofre** | Senha incorreta | Reconhecimento × Erro | `Abrir cofre` | `Senha incorreta. Necessário tentar novamente.` | `Enter OK` |
| **Criar novo cofre** | Modificações não salvas | Confirmação × Alerta | `Criar novo cofre` | `Cofre modificado. Salvar ou descartar?` | `S Salvar`, `D Descartar`, `Esc Voltar` |
| **Criar novo cofre** | Arquivo de destino existente | Confirmação × Alerta | `Criar novo cofre` | `Arquivo '[Nome]' já existe. Sobrescrever?` | `S Sobrescrever`, `I Outro caminho`, `Esc Voltar` |
| **Criar novo cofre** | Senhas não coincidem | Reconhecimento × Erro | `Criar novo cofre` | `Senhas não conferem. Necessário digitar novamente.` | `Enter OK` |
| **Criar novo cofre** | Senha fraca | Confirmação × Alerta | `Criar novo cofre` | `Senha é fraca. Prosseguir ou revisar?` | `P Prosseguir`, `R Revisar`, `Esc Voltar` |
| **Salvar cofre** | Conflito externo | Confirmação × Destrutivo | `Salvar cofre` | `Arquivo modificado externamente. Sobrescrever ou salvar como novo?` | `S Sobrescrever`, `N Salvar como novo`, `Esc Voltar` |
| **Salvar cofre como** | Destino é arquivo atual | Reconhecimento × Alerta | `Salvar cofre como` | `Destino não pode ser o arquivo atual. Necessário escolher outro.` | `Enter OK` |
| **Salvar cofre como** | Arquivo de destino existente | Confirmação × Alerta | `Salvar cofre como` | `Arquivo '[Nome]' já existe. Sobrescrever?` | `S Sobrescrever`, `I Outro caminho`, `Esc Voltar` |
| **Descartar e recarregar** | Arquivo modificado externamente | Confirmação × Destrutivo | `Descartar e recarregar` | `Cofre modificado externamente. Prosseguir com recarregamento?` | `P Prosseguir`, `Esc Voltar` |
| **Descartar e recarregar** | Confirmação de descarte | Confirmação × Destrutivo | `⚠ Descartar e recarregar` | `Todas as alterações serão descartadas. Continuar?` | `C Continuar`, `Esc Voltar` |
| **Alterar senha mestra** | Senhas não coincidem | Reconhecimento × Erro | `Alterar senha mestra` | `Senhas não conferem. Necessário digitar novamente.` | `Enter OK` |
| **Alterar senha mestra** | Senha fraca | Confirmação × Alerta | `Alterar senha mestra` | `Senha é fraca. Prosseguir ou revisar?` | `P Prosseguir`, `R Revisar`, `Esc Voltar` |
| **Alterar senha mestra** | Conflito externo | Confirmação × Destrutivo | `Alterar senha mestra` | `Arquivo modificado externamente. Sobrescrever?` | `S Sobrescrever`, `Esc Voltar` |
| **Exportar cofre** | Senha incorreta (reautenticação) | Reconhecimento × Erro | `Exportar cofre` | `Senha incorreta. Necessário tentar novamente.` | `Enter OK` |
| **Exportar cofre** | Riscos de segurança (não criptografado) | Confirmação × Alerta | `Exportar cofre` | `Arquivo não criptografado. Expor dados sensíveis?` | `E Exportar`, `Esc Voltar` |
| **Exportar cofre** | Arquivo de destino existente | Confirmação × Alerta | `Exportar cofre` | `Arquivo '[Nome]' já existe. Sobrescrever?` | `S Sobrescrever`, `I Outro caminho`, `Esc Voltar` |
| **Importar cofre** | Arquivo de intercâmbio inválido | Reconhecimento × Erro | `Importar cofre` | `Arquivo inválido ou sem Pasta Geral. Necessário corrigir.` | `Enter OK` |
| **Importar cofre** | Confirmação da política de mesclagem | Confirmação × Informativo | `Importar cofre` | `Pastas mescladas. Conflitos substituídos. Confirmar?` | `C Confirmar`, `Esc Voltar` |

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
╰── Enter Confirmar ──────────── Esc Cancelar ──╯
       ↑ text.disabled (bloqueado)
```

**Estado com digitação (ação default ativa):**

```
╭── Senha mestra ────────────────────────────╮
│                                            │
│  Senha                                     │
│  ░••••••••▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
│                                            │
╰── Enter Confirmar ──────────── Esc Cancelar ──╯
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
╰── Enter Confirmar ──────────── Esc Cancelar ──╯
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
| `Enter` → senha incorreta | Erro (5s) | `✕ Senha incorreta` |
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
╰── Enter Confirmar ──────────── Esc Cancelar ──╯
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
╰── Enter Confirmar ──────────── Esc Cancelar ──╯
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
╰── Enter Confirmar ──────────── Esc Cancelar ──╯
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
╰── Enter Confirmar ──────────── Esc Cancelar ──╯
       ↑ text.disabled (senhas divergem)
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
| Ação `Enter Confirmar` | bloqueada (`text.disabled`) | Campo `Nova senha` vazio **ou** campo `Confirmação` vazio **ou** senhas divergentes |
| Ação `Enter Confirmar` | ativa (`accent.primary` **bold**) | Ambos os campos não vazios **e** senhas conferem |
| Ação `Esc Cancelar` | sempre ativa | — |

> **Nota:** a verificação de igualdade entre as senhas ocorre **em tempo real** — a cada tecla no campo `Confirmação` e ao abandonar o campo (Tab ou mudança de foco). Se as senhas divergem, a ação default fica bloqueada e a barra de mensagens exibe erro no lugar da dica de campo.

**Mensagens:**

| Contexto | Tipo | Texto |
|---|---|---|
| Diálogo abre / foco em `Nova senha` (vazio ou válido) | Dica de campo | `• A senha mestra protege todo o cofre — use 12+ caracteres` |
| Foco em `Confirmação` (vazio ou válido) | Dica de campo | `• Redigite a senha para confirmar` |
| Foco em `Confirmação` (senhas divergentes) | Erro (5s) | `✕ As senhas não conferem — digite novamente` |
| Digitação em `Confirmação` (senhas divergentes) | Erro (5s) | `✕ As senhas não conferem — digite novamente` |
| `Enter` → senhas divergentes | Erro (5s) | `✕ As senhas não conferem — digite novamente` |
| Diálogo fecha (confirmação ou cancelamento) | — | Barra limpa *(orquestrador assume)* |

**Comportamento:**
- `Tab` alterna entre os campos `Nova senha` e `Confirmação`
- Medidor de força atualizado a cada tecla no campo `Nova senha`
- Máscara de comprimento fixo (8 `•`) — não revela o tamanho real da senha
- Validação de igualdade em tempo real: a cada tecla no campo `Confirmação` e ao abandonar o campo (Tab)
- Senhas divergentes: ação default bloqueada (`text.disabled`); barra de mensagens exibe erro (`✕`) no lugar da dica de campo; erro permanece até que as senhas confiram ou o campo seja limpo

**Transições especiais:**

| Evento | Efeito |
|---|---|
| Digitação em `Confirmação` torna senhas iguais | Erro na barra é substituído pela dica de campo; ação default muda para `accent.primary` **bold** |
| Digitação em `Confirmação` torna senhas diferentes | Dica de campo é substituída por erro (`✕`, TTL 5s); ação default volta para `text.disabled` |
| Abandonar `Confirmação` (Tab) com senhas divergentes | Erro exibido na barra; foco move para `Nova senha`; ação default bloqueada |
| Abandonar `Confirmação` (Tab) com senhas iguais | Dica exibida na barra; foco move para `Nova senha`; ação default ativa |

---

### FilePicker

**Contexto de uso:** abrir ou salvar arquivo do cofre.
**Token de borda:** `border.focused`
**Dimensionamento:** largura máxima do DS (70 colunas ou 80% do terminal, o menor); altura 80% do terminal. Proporção árvore/arquivos ~40/60.
**Diretório inicial:** determinado pelo fluxo orquestrador. Se não informado, CWD do processo. Se o CWD não existe ou não tem permissão de leitura, fallback para home do usuário (`~`).
**Nome sugerido (modo Save):** determinado pelo fluxo orquestrador. Se não informado, campo inicia vazio. O campo não possui placeholder.
**Filtro de extensão:** apenas arquivos com a extensão `<ext>` (parâmetro `extensao`) são exibidos no painel de arquivos. Não há campo de filtro editável. Arquivos e diretórios ocultos (nome iniciado com `.`) não são exibidos. A extensão é omitida na exibição dos nomes de arquivo (redundante — o filtro já restringe ao formato).
**Padding:** 2 colunas horizontal; **0 vertical** — exceção ao DS [Dimensionamento de diálogos](tui-design-system-novo.md#dimensionamento-de-diálogos). Justificativa: princípio "O Terminal como Meio" — espaço vertical é recurso escasso; o FilePicker é o diálogo mais denso da aplicação (header de caminho + 2 painéis + campo de nome no modo Save). As bordas `╭╮╰╯` e os headers internos (`Caminho:`, `Estrutura`, `Arquivos`, `Nome do arquivo`) criam contenção e separação suficientes sem padding vertical.

O FilePicker opera em dois modos — **Open** e **Save** — com wireframes e condições distintos. Ambos compartilham a mesma anatomia de painéis.

> Nos wireframes abaixo, `░` representa áreas com fundo `surface.input` (campos de entrada).

> **Decisão de layout:** o FilePicker usa separadores internos com junctions em T (`├┬┴┤`) e painéis lado a lado — estrutura que não se encaixa no modelo padrão de diálogos do DS. Esta configuração foi documentada como **exceção justificada** (ver [DS — Exceções ao dimensionamento](tui-design-system-novo.md#dimensionamento-de-diálogos)) e não promoveu uma subseção no DS porque: (1) o FilePicker é o único diálogo com essa complexidade; (2) é um padrão de SO consolidado, não um padrão reutilizável interno; (3) o mecanismo de exceção do DS cobre o caso. Se um segundo diálogo com painéis internos surgir, a exceção será promovida a subseção.

**Barra de comandos durante FilePicker:** enquanto o FilePicker está ativo, a barra de comandos exibe apenas as ações internas do diálogo (conforme regra geral de [Barra de Comandos durante diálogo ativo](#anatomia)). Ações de confirmação/cancelamento (`Enter`/`Esc`) já estão na borda do diálogo — não são duplicadas na barra.

```
  Tab Painel                                                                  F1 Ajuda
```

| Ação | Tecla | Descrição |
|---|---|---|
| Alternar painel | `Tab` | Cicla foco entre os painéis (Árvore → Arquivos no modo Open; Árvore → Arquivos → Campo Nome no modo Save) |
| Ajuda | `F1` | Abre o Help — âncora fixa |

---

#### Contrato de entrada e saída

**Entrada (parâmetros do orquestrador):**

| Parâmetro | Tipo | Obrigatório | Uso |
|---|---|---|---|
| `modo` | `Open \| Save` | Sim | Define título, ações e presença do campo de nome |
| `extensao` | `String` | Sim | Extensão filtrada e adicionada automaticamente ao salvar (ex: `".abditum"`, `".json"`). Deve incluir o ponto inicial. |
| `diretorio_inicial` | `PathBuf` | Não | Diretório onde o FilePicker abre. Default: CWD → fallback `~` |
| `nome_sugerido` | `String` | Não (modo Save) | Valor inicial do campo `Nome do arquivo`. Default: vazio |

**Saída (retorno ao orquestrador):**

| Resultado | Valor | Significado |
|---|---|---|
| Confirmado | `Some(PathBuf)` | Caminho completo do arquivo selecionado (modo Open) ou caminho de salvamento com extensão `<ext>` garantida (modo Save) |
| Cancelado | `None` | Usuário abandonou o diálogo via `Esc` |

---

#### FilePicker — Modo Open

**Título:** `Abrir cofre`
**Objetivo:** selecionar um arquivo `<ext>` existente.

**Wireframe (arquivo selecionado — ação default ativa, scroll em ambos os painéis):**

```
╭── Abrir cofre ─────────────────────────────────────────────────────╮
│  Caminho: /home/usuario/projetos/abditum                           │
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
╰── Enter Abrir ───────────────┴────────────────────── Esc Cancelar ──╯
       ↑ accent.primary + bold (desbloqueado)
```

> Scroll da árvore (`↑` `■` `↓`) substitui o `│` do separador entre painéis. Scroll dos arquivos (`↑` `■` `↓`) substitui o `│` da borda direita do modal. O `┴` na borda inferior marca a junção do separador com a base do diálogo. Metadados (tamanho + `dd/mm/aa HH:MM`) na mesma linha do nome.

**Wireframe (nenhum arquivo — ação default bloqueada, sem scroll):**

```
╭── Abrir cofre ─────────────────────────────────────────────────────╮
│  Caminho: /home/usuario/documentos                                 │
├─ Estrutura ──────────────────┬─ Arquivos ──────────────────────────┤
│  ▶ /                         │                                     │
│    ▼ usuario                 │  Nenhum cofre neste diretório       │
│      ▼ documentos            │                                     │
│        ▶ fotos               │                                     │
│        ▶ textos              │                                     │
│                              │                                     │
╰── Enter Abrir ───────────────┴────────────────────── Esc Cancelar ──╯
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
| Rótulo `Caminho:` | `text.secondary` | — |
| Valor do caminho | `text.primary` | — |
| Ação default (bloqueada) | `text.disabled` | — |
| Ação default (desbloqueada) | `accent.primary` | **bold** |

**Estados dos componentes:**

| Componente | Estado | Condição |
|---|---|---|
| Painel `Estrutura` (árvore) | sempre visível | — |
| Painel `Arquivos` (lista) | conteúdo visível | Pasta selecionada contém arquivos `<ext>` |
| Painel `Arquivos` (lista) | texto vazio | Pasta selecionada **não** contém arquivos `<ext>` |
| Rótulo `Caminho` | sempre visível, somente leitura | Atualiza ao navegar na árvore |
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
- **Painel de arquivos reflete o cursor da árvore:** ao mover o cursor (`↑↓`) entre pastas na árvore, o painel de arquivos atualiza imediatamente para mostrar os `<ext>` da pasta sob o cursor — não apenas ao expandir. O rótulo `Caminho` e o painel de arquivos acompanham a pasta com cursor, independente de ela estar expandida ou recolhida
- **Navegação por teclado na árvore:** `↑↓` navega entre pastas visíveis; `→` expande pasta recolhida; `←` recolhe pasta expandida; `Enter` avança foco para o primeiro arquivo no painel de arquivos (se a pasta sob o cursor contém `<ext>`; sem efeito se não contém); `Home`/`End` vai ao primeiro/último item visível; `PgUp`/`PgDn` scroll por página
- **Navegação por teclado nos arquivos:** `↑↓` navega entre arquivos; `Enter` confirma seleção (equivale à ação default); `Home`/`End` vai ao primeiro/último arquivo visível; `PgUp`/`PgDn` scroll por página
- Ao navegar para uma pasta na árvore, se ela contém arquivos `<ext>`, o primeiro é pré-selecionado automaticamente no painel de arquivos
- **Indicador de pasta vazia:** pastas sem subdiretórios visíveis usam `▷` conforme o DS — não são expansíveis. `→` não tem efeito sobre elas (nada a expandir). `Enter` segue a regra padrão — avança foco para o painel de arquivos se a pasta contém `<ext>`. `▷` indica ausência de subdiretórios expansíveis — não impede que a pasta contenha arquivos `<ext>` exibidos no painel de arquivos
- **Clique simples em pasta:** move cursor para a pasta (atualiza painel de arquivos e `Caminho`)
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

#### FilePicker — Modo Save

**Título:** `Salvar cofre`
**Objetivo:** escolher diretório e nome para salvar o arquivo do cofre.

**Wireframe (campo nome preenchido — ação default ativa):**

```
╭── Salvar cofre ────────────────────────────────────────────────────╮
│  Caminho: /home/usuario/projetos/abditum                           │
├─ Estrutura ──────────────────┬─ Arquivos ──────────────────────────┤
│  ▶ /                         │  ● database   25.8 MB 15/03/25 14:32│
│    ▼ usuario                 │  ● config       1.2 KB 02/01/25 09:15│
│      ▼ projetos              │                                     │
│        ▼ abditum             │                                     │
│          ▶ docs              │                                     │
│                              │                                     │
├──────────────────────────────┴─────────────────────────────────────┤
│  Nome do arquivo                                                   │
│  ░meu-cofre▌░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
╰── Enter Salvar ───────────────────────────────────────── Esc Cancelar ──╯
       ↑ accent.primary + bold (desbloqueado)
```

**Wireframe (campo nome vazio — ação default bloqueada):**

```
╭── Salvar cofre ────────────────────────────────────────────────────╮
│  Caminho: /home/usuario/projetos                                   │
├─ Estrutura ──────────────────┬─ Arquivos ──────────────────────────┤
│  ▶ /                         │  ● database   25.8 MB 15/03/25 14:32│
│    ▼ usuario                 │                                     │
│      ▼ projetos              │                                     │
│                              │                                     │
├──────────────────────────────┴─────────────────────────────────────┤
│  Nome do arquivo                                                   │
│  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░ │
╰── Enter Salvar ───────────────────────────────────────── Esc Cancelar ──╯
       ↑ text.disabled (bloqueado)
```

> Tokens de estrutura (título, headers, separadores, pasta, arquivo, metadados, Caminho, ação default) idênticos ao [Modo Open](#filepicker--modo-open). Exclusivos do Modo Save:

| Elemento | Token | Atributo |
|---|---|---|
| Label `Nome do arquivo` (campo ativo) | `accent.primary` | **bold** |
| Label `Nome do arquivo` (campo inativo) | `text.secondary` | — |
| Área do campo `░` | `surface.input` | — |
| Cursor `▌` | `text.primary` | — |

**Estados dos componentes:**

| Componente | Estado | Condição |
|---|---|---|
| Painel `Estrutura` (árvore) | sempre visível | — |
| Painel `Arquivos` (lista) | conteúdo visível | Pasta selecionada contém arquivos `<ext>` |
| Painel `Arquivos` (lista) | texto vazio | Pasta selecionada **não** contém arquivos `<ext>` |
| Rótulo `Caminho` | sempre visível, somente leitura | Atualiza ao navegar na árvore |
| Campo `Nome do arquivo` | sempre visível | — |
| Caracteres inválidos para filesystem (`/ \ : * ? " < > \|`) | bloqueados silenciosamente | Tecla não produz efeito — sem mensagem de erro |
| Extensão `<ext>` | adicionada automaticamente | Se o nome digitado não termina em `<ext>` |
| Ação `Enter Salvar` | bloqueada (`text.disabled`) | Campo `Nome do arquivo` vazio |
| Ação `Enter Salvar` | ativa (`accent.primary` **bold**) | Campo `Nome do arquivo` não vazio |
| Ação `Esc Cancelar` | sempre ativa | — |

> **Nota:** a validação de sobrescrita (arquivo já existe) é responsabilidade do fluxo que chamou o FilePicker, não do diálogo. O picker retorna o caminho completo; o fluxo abre diálogo de Confirmação × Destrutivo se necessário.

**Mensagens:**

| Contexto | Tipo | Texto |
|---|---|---|
| Diálogo abre / foco na árvore | Dica de campo | `• Navegue pelas pastas e escolha onde salvar` |
| Foco no painel de arquivos | Dica de campo | `• Arquivos existentes neste diretório` |
| Foco no campo `Nome do arquivo` (vazio) | Dica de campo | `• Digite o nome do arquivo — <ext> será adicionado automaticamente` |
| Foco no campo `Nome do arquivo` (preenchido) | Dica de campo | `• Confirme para salvar o cofre` |
| Diálogo fecha | — | Barra limpa *(orquestrador assume)* |

**Comportamento:**

- **Foco inicial:** árvore de diretórios (painel esquerdo)
- **Ordem do Tab:** Árvore → Arquivos → Campo `Nome do arquivo` → volta (3 stops)
- **Scroll:** cada painel tem scroll independente com indicadores `↑`/`↓`/`■` na borda direita do respectivo painel
- Navegação na árvore e painel de arquivos idêntica ao modo Open, com uma exceção: **`Enter` no painel de arquivos copia o nome (sem extensão) para o campo `Nome do arquivo` e move foco para o campo** — não confirma o diálogo. A confirmação requer `Enter` novamente (no campo ou em qualquer contexto com ação default ativa)
- No painel de arquivos: `↑↓` apenas destaca o arquivo (highlight) — **não** copia o nome para o campo. Somente `Enter` ou clique simples no arquivo copiam o nome (sem extensão) para o campo `Nome do arquivo`
- Ao navegar na árvore, o campo `Nome do arquivo` **não é limpo** — preserva o nome digitado
- Extensão `<ext>` é adicionada silenciosamente ao caminho de retorno, sem alterar o texto exibido no campo
- **Duplo-clique em pasta:** expande/recolhe (mesmo que `→`/`←`)
- **Duplo-clique em arquivo existente:** copia o nome para o campo `Nome do arquivo`
- Scroll do mouse, arquivos ocultos, caminho longo, permissões, fallback CWD, ordenação, indentação, formato de metadados e truncamento: idêntico ao [Modo Open](#filepicker--modo-open)

**Transições especiais:**

| Evento | Efeito |
|---|---|
| Clique simples em arquivo existente no painel | Nome copiado para campo `Nome do arquivo`; ação default muda para `accent.primary` **bold** |
| `Enter` no painel de arquivos | Nome copiado para campo `Nome do arquivo`; foco move para o campo. **Não** confirma o diálogo |
| `Enter` na árvore (pasta com `<ext>`) | Foco avança para o primeiro arquivo no painel de arquivos |
| `Enter` na árvore (pasta sem `<ext>`) | Sem efeito |
| `→` em pasta recolhida | Pasta expandida; cursor permanece na pasta |
| `←` em pasta expandida | Pasta recolhida; cursor permanece na pasta |
| Limpar campo `Nome do arquivo` | Ação default volta para `text.disabled` |
| `Enter` com campo preenchido | Diálogo fecha com caminho completo (diretório + nome + `<ext>`) |
| Tentar expandir pasta sem permissão | Erro na barra (`✕ Sem permissão para acessar <pasta>`); pasta permanece recolhida |

---

### Help

**Contexto de uso:** lista todas as ações do ActionManager, agrupadas. Acionado por `F1` em qualquer contexto.
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
│  Insert      Novo segredo                                        │
│  ^E          Editar segredo                                      │
│  Delete      Excluir segredo                                     │
│                                                                  │
│  Cofre                                                           │
│  ^S          Salvar cofre                                        │
│  ^Q          Sair (salva se necessário)                          │
│  F1          Esta ajuda                                          │
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

> **Nota:** os wireframes são snapshots ilustrativos. O conteúdo real é gerado dinamicamente pelo ActionManager a partir do contexto ativo.

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Título `Ajuda — Atalhos e Ações` | `text.primary` | **bold** |
| Label do grupo (`Navegação`, `Segredo`, `Cofre`) | `text.secondary` | **bold** |
| Tecla (ex: `Ctrl+R`, `Insert`, `^S`) | `accent.primary` | — |
| Descrição da ação | `text.primary` | — |
| Seta de scroll (`↑` / `↓` na borda direita) | `text.secondary` | — |
| Thumb de posição (`■` na borda direita) | `text.secondary` | — |
| Borda | `border.default` | — |

#### Estados dos componentes

| Componente | Estado | Condição |
|---|---|---|
| Conteúdo | sem scroll | Todas as ações cabem na viewport |
| Conteúdo | com scroll | Ações excedem a viewport — indicadores `↑`/`↓` e thumb `■` na borda direita (ver [DS — Scroll em diálogos](tui-design-system-novo.md#scroll-em-diálogos)) |
| `F1` na barra de comandos | oculto (`HideFromBar`) | Enquanto o Help estiver aberto |
| Barra de comandos | vazia | Help não registra ações internas na barra |

#### Eventos

| Evento | Efeito |
|---|---|
| `F1` pressionado (modal fechado) | Abre o modal; barra de comandos fica vazia; `F1` oculto |
| `F1` pressionado (modal aberto) | Fecha o modal; `F1` volta visível na barra |
| `Esc` | Fecha o modal; `F1` volta visível na barra |
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
| Busca ativa | Nome da app `·` cofre `•` + abas | Campo de busca à esquerda + aba ativa suspensa à direita | Visíveis (3) |

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

#### Modo busca ativo

Ativo enquanto o campo de busca estiver aberto (ver [Busca de Segredos](#busca-de-segredos)). Disponível apenas no Modo Cofre com cofre aberto.

A linha separadora (linha 2) é substituída pelo campo de busca. A aba ativa permanece suspensa à direita na mesma linha, sem alteração de posição ou estilo.

**Campo aberto, sem query (recém-ativado):**

```
  Abditum · cofre •                      ╭───────╮  ╭ Modelos ╮  ╭ Config ╮
 ─ Busca: ────────────────────────────────╯ Cofre ╰──────────────────────────
```

**Campo aberto, com query:**

```
  Abditum · cofre •                      ╭───────╮  ╭ Modelos ╮  ╭ Config ╮
 ─ Busca: gmail ──────────────────────────╯ Cofre ╰──────────────────────────
```

**Regras de layout do campo na linha separadora:**

| Elemento | Largura | Notas |
|---|---|---|
| `─ Busca: ` (prefixo fixo) | 9 colunas | `─` + espaço + `Busca:` + espaço |
| Texto da query | variável | Em `accent.primary` **bold** |
| `─` preenchimento | restante − largura da aba ativa − 2 (margem direita mínima) | Preenche até a aba |
| Aba ativa (`╯ Texto ╰`) | igual ao estado normal | Posição e estilo inalterados |

- **Query longa:** truncada à **esquerda** com `…` — a parte mais recente da query fica sempre visível
- A largura disponível para a query é calculada em tempo real e recalculada a cada resize do terminal

**Tokens exclusivos do modo busca na linha separadora:**

| Elemento | Token | Atributo |
|---|---|---|
| `─ Busca: ` rótulo | `border.default` | — |
| Texto da query | `accent.primary` | **bold** |
| `─` preenchimento | `border.default` | — |

> **Exceção de layout documentada:** a linha separadora do cabeçalho tem papel estrutural fixo no DS (divisa cabeçalho ↔ área de trabalho). Durante o modo busca, essa linha assume papel adicional de display do campo de busca. Exceção justificada pelo princípio **Hierarquia da Informação** — o campo imediatamente acima da árvore cria relação visual direta entre query e resultado — e pelo princípio **O Terminal como Meio** — espaço vertical é recurso escasso. Escopo-limitada ao Modo Cofre com busca ativa.

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
| Navegação entre modos (Cofre / Modelos / Config) | Aba ativa muda; nova aba suspensa na linha separadora |
| Terminal redimensionado | Nome do cofre recalcula truncamento |

---

#### Comportamento

- **Abas clicáveis** — mouse troca o modo ativo ao clicar no texto ou na borda da aba (área de hit inclui linhas 1 e 2 da aba)
- **Navegação por teclado** — `F2` Cofre, `F3` Modelos, `F4` Config (escopo Área de trabalho — só ativas com cofre aberto)
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

Cada ação na barra segue o formato: **TECLA Label** — tecla em `accent.primary` **bold**, label em `text.primary`. Ações separadas por `·` em `text.secondary`. A ação `F1` (Ajuda) é âncora fixa na extrema direita.

**Estado normal:**

```
  ^I Novo · ^E Editar · Del Excluir · ^S Salvar                              F1 Ajuda
```

**Com ação desabilitada (nenhum segredo selecionado):**

```
  ^I Novo · ^E Editar · ^S Salvar                                              F1 Ajuda
```

Ações com `Enabled = false` não aparecem na barra — só no modal de Ajuda. O espaço colapsa; separadores `·` são re-calculados entre ações visíveis.

**Durante diálogo ativo (apenas ações internas):**

```
  Tab Campos · F5 Revelar                                                    F1 Ajuda
```

Ações do ActionManager ficam ocultas. A barra mostra apenas as ações internas do diálogo do topo da pilha. Ações de confirmação/cancelamento (`Enter`/`Esc`) já estão na borda do diálogo — não são duplicadas na barra.

**Espaço restrito:**

```
  ^I Novo                                                                    F1 Ajuda
```

Ações de menor prioridade são ocultadas quando não há espaço. `F1` permanece sempre visível — é via Help que o usuário descobre as ações ocultas.

---

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Tecla da ação (ex: `Insert`) | `accent.primary` | **bold** |
| Label da ação (ex: `Novo`) | `text.primary` | — |
| Separador `·` | `text.secondary` | — |
| `F1` (Ajuda) | `accent.primary` | **bold** |

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
| Troca de foco entre painéis | Ações do painel que recebe foco ficam ativas |
| Seleção de item na árvore | Ações de item (editar, excluir, revelar) ficam `Enabled = true` — aparecem na barra |
| Nenhum item selecionado | Ações de item ficam `Enabled = false` — desaparecem da barra |
| Diálogo aberto (push na pilha) | Troca para ações internas do diálogo |
| Diálogo fechado (pop da pilha) | Volta para ações do ActionManager |
| Terminal redimensionado | Recalcula quais ações cabem (prioridade governa corte) |

---

#### Comportamento

- **Âncora `F1`** — reserva espaço fixo na extrema direita. O cálculo de espaço disponível desconta `F1 Ajuda` antes de distribuir as demais ações
- **Ações desabilitadas desaparecem da barra** — `Enabled = false` remove a ação da barra (não fica exibida como dim). A ação continua listada no Help
- **Diálogos de decisão** (confirmação/reconhecimento) — tipicamente não têm ações internas; a barra pode ficar vazia (apenas `F1 Ajuda`) enquanto o diálogo estiver aberto
- **Diálogos funcionais** (PasswordEntry, FilePicker etc.) — registram ações internas (Tab entre campos, revelar senha, etc.) que aparecem na barra
- **Truncamento** — se mesmo a ação de maior prioridade + `F1 Ajuda` não cabem, a barra mostra apenas `F1 Ajuda`

---

### Barra de Mensagens

**Responsabilidade:** comunicar feedback ao usuário — sucesso, erro, aviso, progresso, dicas.
**Posição:** 1 linha fixa entre a área de trabalho e a barra de comandos (zona Barra de mensagens do [DS — Dimensionamento](tui-design-system-novo.md#dimensionamento-e-layout)).
**Altura fixa:** 1 linha.
**Anatomia:** borda `─` contínua na largura total do terminal. Quando há mensagem, o texto (símbolo + `·` espaço + conteúdo) começa com 2 espaços de padding à esquerda (alinhado com o texto do cabeçalho), seguido de `─` até o fim da linha. O espaço entre símbolo e texto é sempre exatamente 1 caractere.

**Anatomia (exemplo — sucesso):**

```
── ✓ Gmail copiado para a área de transferência ────────────────────────────────
```

Todos os tipos seguem este padrão. Diferenças por tipo: `✓` sucesso · `✕` erro (**bold**) · `⚠` aviso · `◐◓◑◒` spinner · `•` dica (*italic*) · `ℹ` informação · sem mensagem (borda `─` contínua). Mensagem longa truncada com `…` no fim.

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
| Texto | truncado com `…` | Mensagem excede largura disponível (terminal − 2 padding − 2 borda mínima) |

#### Eventos

| Evento | Efeito |
|---|---|
| Operação concluída com sucesso | Exibe `✓` mensagem (`semantic.success`, TTL 5s) |
| Informação neutra | Exibe `ℹ` mensagem (`semantic.info`, TTL 5s) |
| Condição de alerta (ex: bloqueio iminente) | Exibe `⚠` mensagem (`semantic.warning`, permanente, desaparece com input) |
| Falha em operação | Exibe `✕` mensagem (`semantic.error` + bold, TTL 5s) |
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
- **Scroll no separador** — o scroll da árvore é indicado por `↑`/`↓`/`■` no `│` (separador entre painéis). `<╡` e scroll ocupam a mesma coluna: `<╡` tem prioridade sobre `■` em caso de coincidência (ver [DS — Scroll em diálogos](tui-design-system-novo.md#scroll-em-diálogos)). Quando `<╡` coincide com `↑` ou `↓`, `<╡` prevalece — a direção do scroll é implícita pela presença do outro indicador nas demais linhas
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

Esta seção detalha as ações disponíveis ao interagir com a árvore de segredos (painel esquerdo do Modo Cofre) e seus respectivos atalhos de teclado. As regras gerais de navegação e atribuição de teclas são definidas no [Design System — Mapa de Teclas](tui-design-system-novo.md#mapa-de-teclas).

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

| Tecla    | Ação                                     | Notas                                                                      |
|----------|------------------------------------------|----------------------------------------------------------------------------|
| `Enter`  | Focar no painel de detalhes do segredo   | Comporta-se de forma similar ao `Tab` quando o foco está em um segredo.    |
| `Insert` | Novo segredo                             | Cria um novo segredo na pasta atualmente focada.                           |
| `Ctrl+I` | Novo segredo                             | Atalho alternativo para criar um novo segredo.                             |
| `^E`     | Editar segredo                           | Entra no modo de edição para o segredo selecionado.                        |
| `⌃S`     | Favoritar / Desfavoritar segredo         | Toggle — alterna entre favoritado e não favoritado.                        |
| `⌃R`     | Revelar primeiro campo sensível          | Visível apenas se o segredo tiver pelo menos um campo sensível. Abre/atualiza o painel direito. |
| `⌃C`     | Copiar primeiro campo sensível           | Visível apenas se o segredo tiver pelo menos um campo sensível. Agenda limpeza da clipboard. |
| `Delete` | Excluir segredo                          | Marca o segredo selecionado para exclusão (reversível até o salvamento).   |

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

##### Barra de comandos contextualizada (árvore, cursor em segredo)

| Condição | Barra de comandos |
|---|---|
| Segredo sem campo sensível | `Enter Detalhes · ⌃E Editar · ⌃S Favoritar · Del Excluir · F1 Ajuda` |
| Segredo com campo sensível — reveal mascarado | `Enter Detalhes · ⌃E Editar · ⌃S Favoritar · ⌃R Revelar · ⌃C Copiar · Del Excluir · F1 Ajuda` |
| Segredo com campo sensível — reveal com dica | `Enter Detalhes · ⌃E Editar · ⌃S Favoritar · ⌃R Mostrar tudo · ⌃C Copiar · Del Excluir · F1 Ajuda` |
| Segredo com campo sensível — reveal completo | `Enter Detalhes · ⌃E Editar · ⌃S Favoritar · ⌃R Ocultar · ⌃C Copiar · Del Excluir · F1 Ajuda` |

---

### Painel Direito: Detalhe do Segredo — Modo Leitura

**Contexto:** Área de trabalho — Modo Cofre.
**Largura:** ~65% da área de trabalho.
**Responsabilidade:** Exibir o nome, o caminho de pastas, os campos e a observação do segredo selecionado na árvore; permitir navegação entre campos, cópia de valores e reveal de campos sensíveis.

> Este documento especifica apenas o **modo leitura**. O modo edição de valores e o modo edição de estrutura são especificados separadamente.

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
- **Scroll** — a última coluna do painel é sempre reservada para a trilha de scroll, mesmo quando não há overflow — o conteúdo não se desloca ao ativar o scroll (ver [DS — Scroll em diálogos](tui-design-system-novo.md#scroll-em-diálogos))
- **`<╡` e trilha de scroll são independentes** — `<╡` aparece no separador vertical esquerdo e indica qual item da árvore está sendo detalhado; a trilha de scroll aparece na margem direita e reflete o scroll do conteúdo do painel. Um não afeta o outro
- **Posição do cursor ao retornar o foco** — ao receber foco via `Tab` novamente, o cursor vai ao campo que estava ativo antes de o foco sair; se nunca focado, vai ao primeiro campo
- **Breadcrumb — truncamento** — o breadcrumb é truncado à esquerda com `…` se o caminho completo não couber; o nome do segredo e o `★` nunca são truncados

---

## Telas

### Boas-vindas

**Trigger:** Aplicação inicia sem cofre aberto, ou após fechar/bloquear cofre.  
**Interação:** Nenhuma — tela estática. Toda ação disponível via barra de comandos.

**Wireframe (área de trabalho — terminal 80 × 24):**

```
                                                                                
                                                                                
                                                                                
                   ___    __        ___ __                                      
                  /   |  / /_  ____/ (_) /___  ______ ___                       
                 / /| | / __ \/ __  / / __/ / / / __ `__ \                     
                / ___ |/ /_/ / /_/ / / /_/ /_/ / / / / / /                     
               /_/  |_/_.___/\__,_/_/\__/\__,_/_/ /_/ /_/                      
                                                                                
                             v0.1.0                                             
                                                                                
                                                                                
```

> Logo e versão centralizados via `lipgloss.Place()`. As linhas do logo recebem as cores do [DS — Gradiente do logo](tui-design-system-novo.md#gradiente-do-logo) — não representável neste wireframe monocromático.

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Logo (linhas 1–5) | DS — [Gradiente do logo](tui-design-system-novo.md#gradiente-do-logo) — por linha | — |
| Versão (ex: `v0.1.0`) | `text.secondary` | — |

> As cores do logo não são tokens nomeados — são os valores hexadecimais da tabela de gradiente do DS, aplicados por linha conforme o tema ativo.

#### Estados dos componentes

| Componente | Estado | Condição |
|---|---|---|
| Logo + versão | visível, centralizado | Tela ativa |
| Cabeçalho | sem abas | Nenhum cofre aberto — ver [Cabeçalho — Sem cofre](#sem-cofre-boas-vindas) |

#### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Tela entra em exibição | Dica de uso | `• Abra ou crie um cofre para começar` |

#### Eventos

| Evento | Efeito |
|---|---|
| Aplicação inicia sem cofre | Modo boas-vindas exibido |
| Cofre fechado | Tela boas-vindas exibida |
| Cofre bloqueado | Tela boas-vindas exibida; arquivo permanece em disco, requer nova autenticação |
| Terminal redimensionado | Logo e versão recentralizados |

#### Comportamento

- Logo e versão centralizados horizontal e verticalmente na área de trabalho via `lipgloss.Place()`
- As cores do logo acompanham o tema ativo — mudam instantaneamente com `F12`
- O cabeçalho não exibe abas neste modo (ver [Cabeçalho — Sem cofre](#sem-cofre-boas-vindas))
- **Versão dinâmica** — o texto exibido vem da string injetada em tempo de build via `-ldflags "-X main.version=$(git describe --tags --always)"`. Em builds locais sem tag, exibe `dev`. O valor **nunca** é hardcoded no fonte

---

<!-- SEÇÕES FUTURAS — a preencher pela equipe -->

<!--
## Telas (continuação)

### Modo Cofre
### Modo Modelos
### Modo Configurações

## Componentes

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
