# Especificação Visual — Abditum TUI · Parte 2

> Suplemento de [`tui-specification.md`](tui-specification.md).
> Cobre lacunas identificadas na revisão: badge de força de senha, interface de busca, modo de reordenação,
> painel de modelo em edição, fluxos visuais adicionais e tabela mestre de atalhos.

---

## Badge de Força de Senha

**Contexto:** usado no modal `PasswordCreate` e nos fluxos de criação e alteração de senha mestra do cofre.
O badge aparece **imediatamente abaixo do campo "Nova senha"**, atualizado em tempo real a cada tecla.

> **Requisito:** `crypto.EvaluatePasswordStrength` retorna `StrengthWeak` / `StrengthStrong`.
> Senha fraca não bloqueia o submit — exibe aviso não bloqueante.

---

### Campo vazio

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
                    <span style="color:#7aa2f7">╭─────────────────────────────────────────╮</span>
                    <span style="color:#7aa2f7">│</span>  <span style="background:#24283b;color:#a9b1d6"><strong>Definir senha mestra</strong></span>                   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7"><strong>Nova senha</strong></span>                              <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">╭───────────────────────────────────╮</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">│</span> ▌                                  <span style="color:#7aa2f7">│</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">╰───────────────────────────────────╯</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#565f89">Confirmação</span>                             <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#414868">╭───────────────────────────────────╮</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#414868">╰───────────────────────────────────╯</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">╰─────────────────────────────────────────╯</span>
</pre>

Badge ausente enquanto o campo está vazio — não existe avaliação sem conteúdo.

---

### Senha fraca (`StrengthWeak`)

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
                    <span style="color:#7aa2f7">╭─────────────────────────────────────────╮</span>
                    <span style="color:#7aa2f7">│</span>  <span style="background:#24283b;color:#a9b1d6"><strong>Definir senha mestra</strong></span>                   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7"><strong>Nova senha</strong></span>                    <span style="color:#e0af68">⚠ Fraca</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">╭───────────────────────────────────╮</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">│</span> <span style="color:#565f89">••••••▌</span>                            <span style="color:#7aa2f7">│</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">╰───────────────────────────────────╯</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#e0af68"><em>Senha fraca — use letras, números e símbolos</em></span>  <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#565f89">Confirmação</span>                             <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#414868">╭───────────────────────────────────╮</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#414868">╰───────────────────────────────────╯</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">╰─────────────────────────────────────────╯</span>
</pre>

O badge `⚠ Fraca` aparece na mesma linha do label "Nova senha", alinhado à direita.
A mensagem inline é aviso não bloqueante — o usuário **pode** continuar e confirmar a senha fraca.

---

### Senha forte (`StrengthStrong`)

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7"><strong>Nova senha</strong></span>                  <span style="color:#9ece6a">✓ Forte</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">╭───────────────────────────────────╮</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">│</span> <span style="color:#565f89">••••••••••••▌</span>                      <span style="color:#7aa2f7">│</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">╰───────────────────────────────────╯</span>   <span style="color:#7aa2f7">│</span>
</pre>

Badge `✓ Forte` sem mensagem inline — não há nada a avisar quando a senha é forte.

---

### Tokens do Badge

| Elemento | Token | Atributo | Notas |
|---|---|---|---|
| Badge fraca `⚠ Fraca` | `semantic.warning` | — | Inline na linha do label, alinhado à direita |
| Badge forte `✓ Forte` | `semantic.success` | — | Inline na linha do label, alinhado à direita |
| Mensagem inline fraca | `semantic.warning` | *italic* | Abaixo do campo; desaparece quando vira Forte |
| Campo com senha fraca | `border.focused` | — | Borda NÃO muda para erro — fraca não é inválida |

> **Princípio:** força fraca é um aviso informativo, não um erro de validação. A borda do campo permanece `border.focused` (não vira `semantic.error`). Apenas a confirmação não-conferida usa `semantic.error`.

---

## Interface de Busca

**Acionamento:** `^F` em qualquer estado do modo Cofre (painel esquerdo ou direito com foco).
**Comportamento:** a busca transforma temporariamente a árvore em uma lista plana de resultados. O painel direito continua funcionando normalmente para o item selecionado. Pressionar `Esc` com o campo de busca vazio retorna à árvore normal.

---

### Campo de busca ativo, sem resultados ainda

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <strong><span style="color:#7aa2f7">Abditum</span></strong>  <span style="color:#414868">·</span>  <span style="color:#565f89">cofre.abditum</span>              <span style="background:#283457;color:#7aa2f7">╭────────╮</span>  <span style="color:#565f89">╭ Modelos ╮</span>  <span style="color:#565f89">╭ Configurações ╮</span>
 <span style="color:#414868">─────────────────────────────────────────────╯</span><span style="background:#283457;color:#7aa2f7"><strong> Cofre  </strong></span><span style="color:#414868">╰──────────────────────────────────────────</span>
  <span style="color:#7aa2f7">╭ Buscar ────────────────────────────╮</span>   <span style="color:#7aa2f7">│</span>
  <span style="color:#7aa2f7">│</span> ▌                                  <span style="color:#7aa2f7">│</span>   <span style="color:#7aa2f7">│</span>
  <span style="color:#7aa2f7">╰────────────────────────────────────╯</span>   <span style="color:#7aa2f7">│</span>
  <span style="color:#414868">────────────────────────────────────────┤</span>
                                          <span style="color:#414868">│</span>
    <span style="color:#565f89"><em>Digite para buscar segredos</em></span>           <span style="color:#414868">│</span>
                                          <span style="color:#414868">│</span>
</pre>

O campo de busca substitui a primeira linha do painel esquerdo. A área de resultados abaixo fica vazia com placeholder.

---

### Busca com resultados

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <strong><span style="color:#7aa2f7">Abditum</span></strong>  <span style="color:#414868">·</span>  <span style="color:#565f89">cofre.abditum</span>              <span style="background:#283457;color:#7aa2f7">╭────────╮</span>  <span style="color:#565f89">╭ Modelos ╮</span>  <span style="color:#565f89">╭ Configurações ╮</span>
 <span style="color:#414868">─────────────────────────────────────────────╯</span><span style="background:#283457;color:#7aa2f7"><strong> Cofre  </strong></span><span style="color:#414868">╰──────────────────────────────────────────</span>
  <span style="color:#7aa2f7">╭ Buscar ── 3 resultado(s) ──────────╮</span>   <span style="color:#7aa2f7">│</span>  Gmail                                           <span style="color:#bb9af7">★</span>
  <span style="color:#7aa2f7">│</span> gm▌                                <span style="color:#7aa2f7">│</span>   <span style="color:#7aa2f7">│</span>  <span style="color:#414868">────────────────────────────────────────────</span>
  <span style="color:#7aa2f7">╰────────────────────────────────────╯</span>   <span style="color:#7aa2f7">│</span>  <span style="color:#565f89">URL</span>           https://mail.google.com
  <span style="color:#414868">────────────────────────────────────────┤</span>  <span style="color:#565f89">Usuário</span>       fulano@gmail.com
<span style="background:#283457;color:#a9b1d6"><strong>  ► Gmail                              </strong></span><span style="color:#7aa2f7">&lt;╡</span>  <span style="color:#565f89">Senha</span>         <span style="color:#565f89">••••••••</span>              <span style="color:#7aa2f7">F16</span>
  <span style="color:#565f89">●</span> <span style="color:#565f89">Sites e Apps</span>  <span style="color:#7dcfff">gm</span>ail@empresa.com     <span style="color:#7aa2f7">│</span>
  <span style="color:#565f89">●</span> Team <span style="color:#7dcfff">Gm</span>ail                          <span style="color:#7aa2f7">│</span>  <span style="color:#565f89">Observação</span>    Conta pessoal
                                          <span style="color:#7aa2f7">│</span>
</pre>

**Detalhes do painel de resultados:**
- Contador no cabeçalho do campo: `N resultado(s)` em `text.secondary`
- Lista **plana** (sem hierarquia de pastas)
- Contexto da pasta exibido como `text.secondary` em itálico na segunda linha do item (quando o nome do segredo não inclui a query)
- Termos correspondentes destacados com `special.match` (bold) dentro do nome do item
- O `&lt;╡` conector e painel direito funcionam normalmente

---

### Busca sem resultados

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#7aa2f7">╭ Buscar ── 0 resultado(s) ──────────╮</span>   <span style="color:#7aa2f7">│</span>
  <span style="color:#7aa2f7">│</span> xyz123▌                            <span style="color:#7aa2f7">│</span>   <span style="color:#7aa2f7">│</span>
  <span style="color:#7aa2f7">╰────────────────────────────────────╯</span>   <span style="color:#7aa2f7">│</span>
  <span style="color:#414868">────────────────────────────────────────┤</span>
                                          <span style="color:#414868">│</span>
    <span style="color:#565f89"><em>Nenhum resultado para "xyz123"</em></span>        <span style="color:#414868">│</span>
                                          <span style="color:#414868">│</span>
</pre>

---

### Navegação na Busca

| Tecla | Ação |
|---|---|
| `^F` | Abre busca (foco vai para o campo) |
| Qualquer caractere | Filtra resultados em tempo real |
| `↑` / `↓` | Move cursor entre resultados |
| `Enter` | Seleciona resultado — abre no painel direito |
| `Tab` | Foco → painel direito (resultado selecionado) |
| `Esc` com texto no campo | Limpa o campo (mantém busca aberta) |
| `Esc` com campo vazio | Fecha busca — retorna à árvore normal |
| `Backspace` até vazio + `Esc` | Alternativa para fechar |

### Tokens da Busca

| Elemento | Token | Atributo |
|---|---|---|
| Borda do campo de busca | `border.focused` | — |
| Cabeçalho `╭ Buscar` | `border.focused` | — |
| Contador `N resultado(s)` | `text.secondary` | — |
| Item normal na lista | `text.primary` | — |
| Termo correspondente destacado | `special.match` | **bold** |
| Subtítulo de pasta (contexto) | `text.secondary` | *italic* |
| Item selecionado | `special.highlight` | **bold** |
| Separador `────` abaixo do campo | `border.default` | — |
| Placeholder / "Nenhum resultado" | `text.secondary` | *italic* |

> **Nota de busca:** o painel direito não é afetado pela busca — continua exibindo o segredo selecionado. O `&lt;╡` conector funciona normalmente.

> **Distinção `special.match` vs `semantic.info`:** `special.match` (dourado) é exclusivo para highlight de substring em resultados de busca. `semantic.info` (ciano) é exclusivo para indicadores de estado de sessão (`+`, `~`). Ambos podem coexistir no mesmo item de árvore — cores distintas evitam ambiguidade.

> **Nota de segurança:** valores de campos sensíveis **nunca** participam da busca — apenas seus nomes. O texto digitado no campo de busca não é logado nem persiste após fechar.

---

## Painel de Modelo: Estados de Edição

A seção Área de Trabalho: Modelos define o layout geral, idêntico ao modo Cofre (35/65). Esta seção especifica os estados de edição do painel direito de modelo.

---

### Modelo em leitura

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  Cartão de Crédito
  <span style="color:#414868">──────────────────────────────────────────────────────────</span>
  <span style="color:#565f89">Campo</span>               <span style="color:#565f89">Tipo</span>          <span style="color:#565f89">Posição</span>
  <span style="color:#414868">──────────────────────────────────────────────────────────</span>
  Número              texto              1
  Titular             texto              2
  Validade            texto              3
  CVV                 <span style="color:#565f89">sensível</span>           4
  Banco               texto              5
</pre>

| Elemento | Token | Atributo |
|---|---|---|
| Título do modelo | `text.primary` | **bold** |
| Cabeçalhos de coluna (`Campo`, `Tipo`, `Posição`) | `text.secondary` | — |
| Valor de campo "sensível" | `text.secondary` | — |
| Valor de campo "texto" | `text.primary` | — |

---

### Modelo em edição

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  Cartão de Crédito  <span style="color:#e0af68">•</span>
  <span style="color:#414868">──────────────────────────────────────────────────────────</span>
  <span style="color:#565f89">Campo</span>               <span style="color:#565f89">Tipo</span>          <span style="color:#565f89">Posição</span>
  <span style="color:#414868">──────────────────────────────────────────────────────────</span>
  <span style="color:#7aa2f7"><strong>Número</strong></span>              <span style="color:#3b4261">texto</span>          <span style="color:#3b4261">1</span>
  <span style="color:#414868">╭─────────────────────────────────╮</span>
  <span style="color:#565f89">│</span> Número▌                          <span style="color:#7aa2f7">│</span>
  <span style="color:#7aa2f7">╰─────────────────────────────────╯</span>
  <span style="color:#3b4261">Titular             texto              2</span>
  <span style="color:#3b4261">Validade            texto              3</span>
  <span style="color:#3b4261">CVV                 sensível           4</span>
  <span style="color:#3b4261">Banco               texto              5</span>
</pre>

Campo ativo: expandido para input inline. Outros campos: `text.disabled` + dim.
`•` no título sinaliza alterações pendentes.

---

### Adicionando novo campo

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#3b4261">Número              texto              1</span>
  <span style="color:#3b4261">Titular             texto              2</span>
  <span style="color:#3b4261">Validade            texto              3</span>
  <span style="color:#3b4261">CVV                 sensível           4</span>
  <span style="color:#3b4261">Banco               texto              5</span>
  <span style="color:#414868">──────────────────────────────────────────────────────────</span>
  <span style="color:#7aa2f7"><strong>Novo campo</strong></span>
  <span style="color:#7aa2f7">╭──────────────────────────╮</span>  <span style="color:#565f89">Tipo:</span>  <span style="background:#283457;color:#a9b1d6"><strong> texto </strong></span> · sensível
  <span style="color:#7aa2f7">│</span> ▌                         <span style="color:#7aa2f7">│</span>
  <span style="color:#7aa2f7">╰──────────────────────────╯</span>
</pre>

Novo campo aparece na última posição. O seletor de tipo usa `←` / `→` para alternar entre `texto` e `sensível`. A opção ativa tem fundo `special.highlight`. `Enter` confirma; `Esc` descarta.

---

### Navegação no Modelo

| Tecla | Ação |
|---|---|
| `F36` | Adiciona novo campo (foco vai para o campo de nome) |
| `↑` / `↓` | Navega entre campos existentes |
| `Enter` sobre campo | Edita nome do campo inline |
| `←` / `→` sobre tipo | Alterna tipo (apenas ao adicionar — tipo não pode ser alterado depois) |
| `F39` | Exclui campo selecionado (com `DialogAlert`) |
| `^S` | Salva alterações do modelo |
| `Esc` | Descarta alterações — `DialogAlert` se houver mudanças |

---

## Indicadores de Estado de Sessão na Árvore

Wireframe completo mostrando todos os estados de sessão coexistindo na mesma árvore.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#bb9af7">★</span> Favoritos              <span style="color:#565f89">(2)</span>
  ▼ Geral                 <span style="color:#565f89">(8)</span>
    ▼ Sites e Apps        <span style="color:#565f89">(5)</span>
  <span style="background:#283457;color:#a9b1d6"><strong>  ► Google               (2) </strong></span>
      <span style="color:#565f89">●</span> Gmail
      <span style="color:#565f89">●</span> YouTube               <span style="color:#7dcfff">+</span>   <span style="color:#565f89">← adicionado na sessão</span>
      <span style="color:#565f89">●</span> LinkedIn              <span style="color:#7dcfff">~</span>   <span style="color:#565f89">← modificado na sessão</span>
      <span style="color:#565f89">✕</span> <span style="color:#565f89"><s>Facebook</s></span>              <span style="color:#565f89">← marcado para exclusão</span>
    ▼ Financeiro          <span style="color:#565f89">(3)</span>
      <span style="color:#565f89">●</span> Nubank
      <span style="color:#bb9af7">★</span> Bradesco
</pre>

Os indicadores `+` e `~` aparecem **no final da linha**, após o nome do segredo, alinhados à direita.
Segredos marcados para exclusão (`✕`) **não** exibem `+` ou `~` — o estado de exclusão tem prioridade visual.

| Indicador | Posição | Token | Condição |
|---|---|---|---|
| `+` | Após o nome (direita) | `semantic.info` | `sessionState == StateIncluded` |
| `~` | Após o nome (direita) | `semantic.info` | `sessionState == StateModified` |
| `✕` + strikethrough | Prefixo | `semantic.error` + `special.muted` | `sessionState == StateDeleted` |
| _(sem indicador)_ | — | — | `sessionState == StateOriginal` |

> **Prioridade de exibição:** Exclusão > Modificado > Adicionado. Um segredo não pode ter dois indicadores simultâneos.

---

### Item adicionado e selecionado

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="background:#283457;color:#a9b1d6"><strong>► YouTube                   +  </strong></span>
</pre>

O indicador `+` permanece visível no estado selecionado — fundo `special.highlight` não oculta o indicador.

---

### Item modificado e selecionado

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="background:#283457;color:#a9b1d6"><strong>● LinkedIn                  ~  </strong></span>
</pre>

---

## Modo de Reordenação

Usuário pode reordenar segredos e pastas dentro da mesma pasta pai. O modo reordenação é ativado por `Alt+↑` / `Alt+↓` com o cursor posicionado no item que se deseja mover.

> **Não existe modo reordenação explícito** — a reordenação acontece diretamente com `Alt+↑` / `Alt+↓`. O item "agarra" enquanto a tecla Alt estiver pressionada, movendo 1 posição por pressão.

---

### Item sendo reordenado

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
    ▼ Sites e Apps        <span style="color:#565f89">(5)</span>
      <span style="color:#565f89">●</span> Gmail
  <span style="background:#283457;color:#7aa2f7"><strong>  ↕ LinkedIn                 ~  </strong></span>   <span style="color:#565f89">← sendo movido</span>
      <span style="color:#565f89">●</span> YouTube               <span style="color:#7dcfff">+</span>
      <span style="color:#565f89">✕</span> <span style="color:#565f89"><s>Facebook</s></span>
</pre>

Enquanto `Alt` está pressionado e o item está sendo movido:
- Prefixo `↕` substitui `●` / `★` — indica que o item está em movimento
- Fundo `special.highlight` com texto `accent.primary` — distinção visual clara do item em modo de reordenação
- Estado de sessão (`~`, `+`) permanece visível

| Elemento | Token | Atributo |
|---|---|---|
| Prefixo `↕` (em movimento) | `accent.primary` | **bold** |
| Texto do item em movimento | `accent.primary` | **bold** |
| Fundo do item em movimento | `special.highlight` | — |

---

### Teclas de Reordenação

| Tecla | Ação |
|---|---|
| `Alt+↑` | Move item 1 posição para cima (dentro da mesma pasta) |
| `Alt+↓` | Move item 1 posição para baixo (dentro da mesma pasta) |

> **Limite:** não é possível mover um segredo para fora da pasta (isso é uma operação de mover, não de reordenar). Ao atingir a primeira ou última posição, a tecla é ignorada.

> **Cabeçalho `•`:** qualquer reordenação marca o cofre como dirty imediatamente.

---

## Tabela Mestre de Atalhos

Mapa de referência de todas as F-keys e atalhos especiais usados na aplicação. Agrupa por contexto funcional.

> **Base para o modal Ajuda:** esta tabela é a fonte canônica do que aparece no modal de ajuda.

---

### Cofre (F1–F15)

| Tecla | Contexto | Ação |
|---|---|---|
| `F1` | Boas-vindas / Qualquer | Abrir cofre (FilePicker modo open) |
| `F2` | Boas-vindas | Criar cofre (FilePicker modo save + PasswordCreate) |
| `F3` | Cofre aberto | Bloquear cofre manualmente |
| `F5` | Cofre aberto | Sair da aplicação |
| `F7` | Cofre aberto | Alterar senha mestra |
| `F9` | Cofre aberto | Salvar como (FilePicker modo save) |
| `F10` | Cofre aberto | Descartar alterações e recarregar |
| `F11` | Cofre aberto | Exportar (FilePicker modo save, sem criptografia) |
| `F12` | Qualquer | Alternar tema (Tokyo Night ↔ Cyberpunk) sem abrir menu |
| `F13` | Cofre aberto | Importar (FilePicker modo open, arquivo de intercâmbio) |

---

### Segredo (F16–F26)

| Tecla | Contexto | Ação |
|---|---|---|
| `F16` | Painel direito (segredo com campo sensível) | Revelar / ocultar campo sensível em foco |
| `F17` | Painel direito | Copiar valor do campo em foco para clipboard |
| `F18` | Painel esquerdo (segredo selecionado) | Favoritar / desfavoritar segredo |
| `F19` | Painel esquerdo (segredo selecionado) | Duplicar segredo |
| `F20` | Painel esquerdo (segredo selecionado) | Mover segredo (Select pasta de destino) |
| `F21` | Painel esquerdo | Novo segredo |
| `F22` | Painel esquerdo / direito (segredo selecionado) | Editar segredo |
| `F23` | Painel esquerdo (segredo selecionado) | Marcar segredo para exclusão (DialogAlert) |
| `F24` | Painel esquerdo (segredo excluído selecionado) | Desmarcar exclusão |
| `F25` | Painel esquerdo (segredo) | Alterar estrutura: adicionar campo |
| `F26` | Painel esquerdo (segredo, campo selecionado) | Alterar estrutura: excluir campo |

---

### Pasta (F27–F32)

| Tecla | Contexto | Ação |
|---|---|---|
| `F27` | Painel esquerdo | Criar nova pasta |
| `F28` | Painel esquerdo (pasta selecionada) | Renomear pasta |
| `F29` | Painel esquerdo (pasta selecionada) | Mover pasta (Select pasta de destino) |
| `F30` | Painel esquerdo (pasta selecionada) | Criar modelo a partir do segredo selecionado |
| `F31` | Painel esquerdo (pasta selecionada) | Excluir pasta (DialogAlert) |

---

### Modelo (F33–F40)

| Tecla | Contexto | Ação |
|---|---|---|
| `F33` | Painel esquerdo (Modelos) | Novo modelo |
| `F34` | Painel direito (Modelos) | Editar modelo |
| `F35` | Painel esquerdo (modelo selecionado) | Excluir modelo (DialogAlert) |
| `F36` | Painel direito (modelo em edição) | Adicionar campo ao modelo |
| `F37` | Painel direito (template em edição, campo) | Renomear campo selecionado |
| `F38` | Painel direito (template em edição, campo) | Reordenar campo (Alt+↑ / Alt+↓ no modo edição) |
| `F39` | Painel direito (template em edição, campo) | Excluir campo selecionado (DialogAlert) |

---

### Navegação e Sistema (Universais)

| Tecla | Ação |
|---|---|
| `↑` `↓` | Mover cursor |
| `→` / `Enter` | Expandir pasta ou selecionar item |
| `←` | Recolher pasta ou subir para pai |
| `Tab` | Foco → próximo painel |
| `Shift+Tab` | Foco → painel anterior |
| `Home` / `End` | Primeiro / último item visível |
| `Alt+↑` / `Alt+↓` | Reordenar item |
| `^F` | Buscar segredos |
| `^S` | Salvar cofre |
| `^Q` | Sair da aplicação |
| `?` | Abrir modal Ajuda |
| `Esc` | Cancelar / fechar modal / sair de modo |
| `F201` | Trocar para aba Cofre |
| `F202` | Trocar para aba Modelos |
| `F203` | Trocar para aba Configurações |

> **Sobre F201–F203:** teclas de troca de aba são atalhos lógicos — na prática implementados como `1`, `2`, `3` enquanto nenhum campo de texto está ativo, para evitar conflito com F-keys do terminal.

---

## Fluxos Visuais (continuação)

> Continuação dos Fluxos A–E definidos no documento principal.

---

### Fluxo F — Criar cofre

```
1. Boas-vindas
   ─ barra de comandos: F1 Abrir · F2 Criar · ? · ^Q Sair

2. → F2 → modal FilePicker (modo save)
   ─ campo de nome vazio com placeholder "nome-do-cofre"
   ─ .abditum será adicionado automaticamente

3. → Usuário digita nome → Enter
   → FilePicker fecha
   → modal PasswordCreate abre com título "Definir senha mestra"

4. → Usuário digita senha
   → Badge de força aparece em tempo real (Fraca / Forte)
   → Usuário Tab para Confirmação → digita novamente

5a. → Senhas não conferem → borda de confirmação vira semantic.error
    → "As senhas não coincidem" — usuário corrige e tenta novamente

5b. → Senhas conferem + senha fraca → aviso ⚠ não bloqueante visível
    → Enter → PasswordCreate fecha
    → MsgBusy "◐ Criando cofre..."

5c. → Senhas conferem + senha forte → Enter → PasswordCreate fecha
    → MsgBusy "◐ Criando cofre..."

6. → Sucesso → transição para Área de Trabalho: Cofre
   ─ Cabeçalho: "Abditum · nome-do-cofre.abditum" + aba Cofre ativa
   ─ Árvore com estrutura padrão: Favoritos, Geral (Sites e Apps, Financeiro)
   ─ MsgInfo "✓ Cofre criado" (TTL 3s)
```

---

### Fluxo G — Buscar e copiar credencial

```
1. Área de Trabalho: Cofre aberto (qualquer estado de foco)

2. → ^F → Campo de busca aparece no topo do painel esquerdo
   ─ MsgHint "• Esc para fechar a busca"

3. → Usuário digita "gmail"
   → Lista de resultados aparece em tempo real
   → Termos correspondentes destacados em `special.match` bold

4. → ↑↓ para navegar → Enter no resultado desejado
   → Painel direito mostra o segredo selecionado
   → Conector <╡ aparece na linha do resultado

5. → Tab → foco vai para o painel direito
   → F17 com foco no campo "Senha"
   → MsgInfo "✓ Senha copiada" (TTL 3s)
   → Timer de limpeza de clipboard inicia (30s default)

6. Usuário continua navegando com a busca ativa ou
   → Esc → busca fechada → árvore normal retorna
```

---

### Fluxo H — Alterar senha mestra

```
1. Área de Trabalho: Cofre aberto (modo Cofre ou Configurações)

2. → F7 → DialogAlert "⚠ Alterar senha mestra"
   ─ "Esta ação regravará o cofre imediatamente.
      Não é possível desfazê-la."
   ─ [S] Continuar (semantic.error bold)   [N] Cancelar

3. → [N] ou Esc → sem mudança

4. → [S] → modal PasswordCreate abre com título "Nova senha mestra"
   ─ Campo "Nova senha" com badge de força em tempo real
   ─ Campo "Confirmação"

5. → Usuário define nova senha → Enter
   → MsgBusy "◐ Regravando cofre..."

6. → Sucesso → MsgInfo "✓ Senha mestra alterada" (TTL 3s)
   ─ Cofre permanece aberto (autenticado com nova senha)
   ─ Cabeçalho: "•" NÃO aparece (cofre já foi salvo com nova senha)

Nota: se MsgBusy falhar → MsgError "✗ Falha ao regravar o cofre"
      cofre retorna ao estado anterior (senha anterior ainda ativa)
```

---

### Fluxo I — Descartar alterações e recarregar

```
1. Área de Trabalho: Cofre com alterações não salvas (cabeçalho mostra "•")

2. → F10 → DialogAlert "⚠ Descartar alterações"
   ─ "Todas as alterações desta sessão serão perdidas.
      O cofre será recarregado do arquivo."
   ─ [S] Descartar (semantic.error bold)   [N] Cancelar

3. → [N] ou Esc → sem mudança

4. → [S] → MsgBusy "◐ Recarregando cofre..."
   → vault.Manager.Discard() + storage.Load()
   → MsgInfo "✓ Alterações descartadas" (TTL 3s)
   ─ Cabeçalho: "•" desaparece
   ─ Árvore retorna ao estado do arquivo

Nota: a senha usada para recarregar é a senha ativa no momento do descarte
      (pode ser nova senha se ela foi alterada antes de descartar).
```

---

### Fluxo J — Exportar cofre

```
1. Área de Trabalho: Cofre aberto

2. → F11 → DialogAlert "⚠ Exportar cofre"
   ─ "O arquivo exportado NÃO é criptografado.
      Mantenha-o em local seguro."
   ─ [S] Exportar   [N] Cancelar

3. → [N] → sem mudança

4. → [S] → modal FilePicker (modo save)
   ─ Título: "Exportar cofre (não criptografado)"
   ─ Não adiciona .abditum automaticamente — extensão sugere .json

5. → Usuário escolhe destino e nome → Enter
   → MsgBusy "◐ Exportando..."

6a. → Sucesso → MsgInfo "✓ Cofre exportado em /caminho/cofre.json" (TTL 4s)

6b. → Falha → MsgError "✗ Falha ao exportar — permissão negada"
```

---

### Fluxo K — Importar cofre

```
1. Área de Trabalho: Cofre aberto

2. → F13 → DialogAlert "⚠ Importar dados"
   ─ "Segredos com nomes iguais serão substituídos.
      Esta ação não pode ser desfeita após salvar."
   ─ [S] Importar   [N] Cancelar

3. → [N] → sem mudança

4. → [S] → modal FilePicker (modo open)
   ─ Título: "Importar arquivo de intercâmbio"
   ─ Lista arquivos .json

5. → Usuário seleciona arquivo → Enter → FilePicker fecha
   → MsgBusy "◐ Importando..."

6a. → Sucesso → MsgInfo "✓ Importado: 12 segredos em 3 pastas" (TTL 4s)
   ─ Cabeçalho: "•" aparece (alterações não salvas)
    ─ Itens importados exibem indicador "+" na árvore

6b. → Arquivo inválido → MsgError "✗ Arquivo inválido ou incompatível"
```

---

### Fluxo L — Duplicar segredo

```
1. Área de Trabalho: Cofre, segredo selecionado na árvore

2. → F19 → sem confirmação (não é destrutuvo)
   → Novo segredo criado com nome "Gmail (1)" (ou próximo disponível)
   → Posicionado imediatamente após o original na lista
   → Cursor se move para o novo segredo
   → Indicador "+" exibido (sessionState = StateIncluded)
   ─ Cabeçalho: "•" aparece
   ─ MsgInfo "✓ Segredo duplicado" (TTL 2s)
```

---

### Fluxo M — Mover segredo para outra pasta

```
1. Área de Trabalho: Cofre, segredo selecionado na árvore

2. → F20 → modal Select "Mover para pasta"
   ─ Lista hierárquica de pastas (exceto a pasta atual do segredo)
   ─ Pasta Geral e subpastas (pastas virtuais como Favoritos excluídas)

3. → Usuário navega e seleciona pasta destino → Enter
   → Select fecha
   → Segredo desaparece da pasta origem
   → Segredo aparece na pasta destino (se expandida)
   ─ Cabeçalho: "•" aparece
   ─ MsgInfo "✓ Segredo movido para Financeiro" (TTL 2s)

Se pasta destino já tem segredo com mesmo nome:
   → MsgError "✗ Já existe um segredo com este nome na pasta destino"
   → Operação cancelada — nenhuma mudança
```

---

### Fluxo N — Conflito de arquivo externo ao salvar

```
1. Área de Trabalho: Cofre com alterações (cabeçalho mostra "•")

2. → ^S → storage.DetectExternalChange() retorna true
   → DialogQuestion "? Conflito de arquivo"
   ─ "O arquivo foi modificado externamente desde a abertura."
   ─ [S] Sobrescrever   [N] Salvar como   [Esc] Cancelar

3a. → Esc → operação cancelada, cofre permanece com "•"

3b. → [N] → modal FilePicker (modo save) — usuário escolhe novo destino
    → MsgBusy "◐ Salvando..." → MsgInfo "✓ Salvo em novo arquivo" (TTL 3s)

3c. → [S] → MsgBusy "◐ Salvando..."
    → storage.Save() com force
    → MsgInfo "✓ Cofre salvo" (TTL 2s), "•" desaparece
```

---

## FilePicker no Windows

No Windows, o sistema de arquivos usa letras de unidade em vez de `/` como raiz. O FilePicker precisa de tratamento especial.

---

### Seletor de unidade

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#414868">╭───────────────────────────────────────────────────────────────────────╮</span>
  <span style="color:#414868">│</span>  <strong>Abrir cofre</strong>                                                        <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">─── Diretórios ─────────────────────</span><span style="color:#414868">┬─── Selecione uma unidade ─────────</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="background:#283457;color:#bb9af7"><strong>  ► C:\                           </strong></span>  <span style="color:#414868">│</span>  <span style="color:#565f89"><em>Navegue para escolher um arquivo</em></span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">  ► D:\</span>                              <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">  ► E:\</span>                              <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>                                       <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>
  <span style="color:#414868">╰───────────────────────────────────────────────────────────────────────╯</span>

  <span style="color:#7aa2f7">→</span> Expandir · <span style="color:#7aa2f7">Tab</span> Arquivos · <span style="color:#7aa2f7">Esc</span> Cancelar
</pre>

Letras de unidade tratadas como nós raiz na árvore. `←` no topo da unidade retorna à lista de unidades.

---

### Dentro de uma unidade Windows

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">─── Diretórios ─────────────────────</span><span style="color:#414868">┬─── C:\Users\usuario\cofres ───────</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">  ▾ C:\</span>                              <span style="color:#414868">│</span>  cofre.abditum                     <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">    ▾ Users\</span>                         <span style="color:#414868">│</span>  trabalho.abditum                  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">      ▾ usuario\</span>                     <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="background:#283457;color:#bb9af7"><strong>        ► cofres\                 </strong></span>  <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">        Documents\</span>                   <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>
</pre>

O separador de caminho exibido é `\` no Windows. Internamente o código usa `filepath.Join` — nunca hardcoded.

---

## Estado NO_COLOR e Modo Monocromático

Quando `$NO_COLOR` está definido (ou o terminal informa que não suporta cores), `lipgloss` remove todas as cores. A interface deve permanecer totalmente funcional.

**Princípios de fallback monocromático:**

| Estado visual | Com cor | Fallback NO_COLOR |
|---|---|---|
| Item selecionado (cursor) | `special.highlight` + **bold** | **bold** |
| Aba ativa | `special.highlight` + **bold** | **bold** + borda `╭───╮` (já presente) |
| Badge `⚠ Fraca` | `semantic.warning` | `⚠ Fraca` (símbolo preserva semântica) |
| Badge `✓ Forte` | `semantic.success` | `✓ Forte` |
| "ativado" em configurações | `semantic.success` | `ativado` (texto preserva estado) |
| "desativado" em configurações | `semantic.off` | `desativado` (texto preserva estado) |
| `•` dirty indicator | `semantic.warning` | `•` (símbolo preserva estado) |
| Indicador `+` / `~` | `semantic.info` | `+` / `~` (símbolos preservam estado) |
| Termo correspondente na busca | `special.match` + **bold** | **bold** (sem cor, mas destaque tipográfico preservado) |
| `✕` exclusão | `semantic.error` + strikethrough | `✕` + strikethrough (duplo fallback) |
| `★` favorito | `accent.secondary` | `★` (símbolo preserva semântica) |
| Máscara `••••••••` | `text.secondary` | `text.primary` — mesma visibilidade |
| Modal border | `semantic.*` / `border.*` | `─ ╭ ╮` (borda presente — tipo distinguido por símbolo no título) |
| Campo de input | `surface.input` | `surface.base` — sem distinção visual, mas borda ainda presente |

> **Conclusão:** a interface é funcionalmente operável em NO_COLOR graças à dupla camada de comunicação — cor + símbolo para cada estado crítico. Nenhum estado depende exclusivamente de cor.

---

## Resumo de Lacunas Cobertas

| Lacuna | Bloco |
|---|---|
| Badge de força de senha | Badge de Força de Senha |
| Interface de busca (`^F`) com realce de termos | Interface de Busca |
| Estados completos de edição de modelo | Painel de Modelo: Estados de Edição |
| Indicadores de sessão (`+`, `~`) coexistindo na árvore | Indicadores de Estado de Sessão na Árvore |
| Reordenação com `Alt+↑` / `Alt+↓` | Modo de Reordenação |
| Tabela mestre de atalhos | Tabela Mestre de Atalhos |
| FilePicker no Windows (letras de unidade) | FilePicker no Windows |
| Fallback `NO_COLOR` | Estado NO_COLOR e Modo Monocromático |
| Fluxos F–N | Fluxos Visuais (continuação) |
