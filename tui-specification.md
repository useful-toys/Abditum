# Especificação Visual — Abditum TUI

> Wireframes detalhados de cada componente, estados visuais, tokens de cor, comportamentos de navegação/edição e mapa de eventos.
>
> **Documentos complementares:**
> - [`tui-design-system.md`](tui-design-system.md) — paleta, tipografia, bordas, símbolos, estados, compatibilidade de terminal
> - [`tui-elm-architecture.md`](tui-elm-architecture.md) — arquitetura de componentes (Elm pattern)

---

## Estrutura da Tela

A interface é dividida em quatro zonas horizontais empilhadas verticalmente:

| Zona | Altura | Conteúdo |
|---|---|---|
| **Header** | 2 linhas | Linha 1: nome da app, nome do cofre, indicador não salvo, abas inativas. Linha 2: separadora — aba ativa é suspensa aqui |
| **Work area** | Restante | Conteúdo do modo ativo — boas-vindas, cofre, modelos ou configurações |
| **Barra de mensagens** | 1 linha | Sobreposta à última linha da work area — não reserva linha própria |
| **Command bar** | 1 linha | Ações do contexto ativo ou atalhos do modal do topo da stack |

> **Transição de modo:** a work area não muda durante fluxos — o usuário vê o conteúdo atual com modais sobrepostos. A transição para outro modo (ex: boas-vindas → cofre) só ocorre após o fluxo concluir com sucesso.

---

## Seção 0 — Telas Completas: 4 Modos Principais

Cada tela mostra o layout completo com um cenário representativo. Use como referência visual de conjunto antes de estudar cada componente individualmente.

---

### Modo: Boas-vindas (Welcome)

> Sem cofre aberto. Header exibe apenas o nome da aplicação. Não há abas — não existe contexto de Cofre/Modelos/Configurações. A work area mostra o logo centralizado.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <strong><span style="color:#7aa2f7">Abditum</span></strong>
 <span style="color:#414868">──────────────────────────────────────────────────────────────────────────────────────────</span>
                                                                                          
                                                                                          
             <span style="color:#9d7cd8">    ___    __        ___ __                  </span>
             <span style="color:#89ddff">   /   |  / /_  ____/ (_) /___  ______ ___   </span>
             <span style="color:#7aa2f7">  / /| | / __ \/ __  / / __/ / / / __ `__ \  </span>
             <span style="color:#7dcfff"> / ___ |/ /_/ / /_/ / / /_/ /_/ / / / / / / </span>
             <span style="color:#bb9af7">/_/  |_/_.___/\__,_/_/\__/\__,_/_/ /_/ /_/ </span>
                                                                                          
                                          <span style="color:#565f89">v0.1.0</span>
                                                                                          
                                                                                          
 <span style="color:#414868">──────────────────────────────────────────────────────────────────────────────────────────</span>
  <span style="color:#7aa2f7">F1</span> Abrir · <span style="color:#7aa2f7">F2</span> Criar · <span style="color:#7aa2f7">?</span> · <span style="color:#7aa2f7">^Q</span> Sair
</pre>

> Especificação de tokens: Header → Seção 1 · Logo/versão → Seção 4 · Command bar → Seção 3

---

### Modo: Cofre (Vault)

> Cofre aberto com segredo selecionado.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <strong><span style="color:#7aa2f7">Abditum</span></strong>  <span style="color:#414868">·</span>  <span style="color:#565f89">cofre.abditum</span>  <span style="color:#e0af68">•</span>              <span style="background:#283457;color:#7aa2f7">╭────────╮</span>  <span style="color:#565f89">╭ Modelos ╮</span>  <span style="color:#565f89">╭ Configurações ╮</span>
 <span style="color:#414868">────────────────────────────────┬────────╯</span><span style="background:#283457;color:#7aa2f7"><strong> Cofre  </strong></span><span style="color:#414868">╰──────────────────────────────────────────</span>
  <span style="color:#bb9af7">★</span> Favoritos              <span style="color:#565f89">(2)</span>  <span style="color:#7aa2f7">│</span>  Gmail                                                <span style="color:#bb9af7">★</span>
  ▼ Geral                  <span style="color:#565f89">(8)</span>  <span style="color:#7aa2f7">│</span>  <span style="color:#414868">──────────────────────────────────────────────────────</span>
    ▼ Sites e Apps         <span style="color:#565f89">(5)</span>  <span style="color:#7aa2f7">│</span>  <span style="color:#565f89">URL</span>            https://mail.google.com
      ▶ Google             <span style="color:#565f89">(2)</span>  <span style="color:#7aa2f7">│</span>  <span style="color:#565f89">Usuário</span>        fulano@gmail.com
  <span style="background:#283457;color:#a9b1d6"><strong>  ► Gmail                     </strong></span> &lt;╡  <span style="color:#565f89">Senha</span>          <span style="color:#565f89">••••••••</span>                      <span style="color:#7aa2f7">F16</span>
    <span style="color:#565f89">●</span> YouTube                   <span style="color:#7aa2f7">│</span>
    <span style="color:#565f89">●</span> Facebook                  <span style="color:#7aa2f7">│</span>  <span style="color:#565f89">Observação</span>     Conta pessoal principal
    <span style="color:#565f89">●</span> LinkedIn                  <span style="color:#7aa2f7">│</span>
  ▼ Financeiro             <span style="color:#565f89">(3)</span>  <span style="color:#7aa2f7">│</span>
    <span style="color:#565f89">●</span> Nubank                    <span style="color:#7aa2f7">│</span>
    <span style="color:#565f89">●</span> Inter                     <span style="color:#7aa2f7">│</span>
 <span style="color:#414868">───────────────────────────────┴──────────────────────────────────────────────────────</span>
  <span style="color:#7aa2f7">F21</span> Novo · <span style="color:#7aa2f7">F22</span> Editar · <span style="color:#7aa2f7">F17</span> Copiar · <span style="color:#7aa2f7">^S</span> Salvar · <span style="color:#7aa2f7">?</span>
</pre>

> Especificação de tokens de cada elemento no componente respectivo: Header → Seção 1 · Árvore → Seção 5 · Painel direito → Seção 6 · Command bar → Seção 3

---

### Modo: Modelos

> Aba "Modelos" ativa. Layout idêntico ao Cofre (35/65), mas com lista plana de modelos à esquerda e campos do modelo à direita. Sem hierarquia de pastas, sem campos sensíveis.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <strong><span style="color:#7aa2f7">Abditum</span></strong>  <span style="color:#414868">·</span>  <span style="color:#565f89">cofre.abditum</span>                 <span style="color:#565f89">╭ Cofre ╮</span>  <span style="background:#283457;color:#7aa2f7">╭──────────╮</span>  <span style="color:#565f89">╭ Configurações ╮</span>
 <span style="color:#414868">────────────────────────────────────────────────────╯</span><span style="background:#283457;color:#7aa2f7"><strong> Modelos  </strong></span><span style="color:#414868">╰──────────────────────────────</span>
  Login Padrão                    <span style="color:#7aa2f7">│</span>  Login Padrão
  <span style="background:#283457;color:#a9b1d6"><strong>  Cartão de Crédito              </strong></span> <span style="color:#7aa2f7">│</span>  <span style="color:#414868">────────────────────────────────────────────────────</span>
  Conta Bancária                  <span style="color:#7aa2f7">│</span>  <span style="color:#565f89">Campo</span>          <span style="color:#565f89">Tipo</span>          <span style="color:#565f89">Obrigatório</span>
  SSH Key                         <span style="color:#7aa2f7">│</span>  <span style="color:#414868">────────────────────────────────────────────────────</span>
  Wi-Fi                           <span style="color:#7aa2f7">│</span>  Nome           texto         <span style="color:#9ece6a">sim</span>
  API / Token                     <span style="color:#7aa2f7">│</span>  Número         texto         <span style="color:#9ece6a">sim</span>
                                  <span style="color:#7aa2f7">│</span>  Validade        texto         não
                                  <span style="color:#7aa2f7">│</span>  CVV             <span style="color:#565f89">senha</span>         não
                                  <span style="color:#7aa2f7">│</span>  Titular         texto         não
                                  <span style="color:#7aa2f7">│</span>
                                  <span style="color:#7aa2f7">│</span>
 <span style="color:#414868">───────────────────────────────┴──────────────────────────────────────────────────────</span>
  <span style="color:#7aa2f7">F33</span> Novo · <span style="color:#7aa2f7">F34</span> Editar · <span style="color:#7aa2f7">F35</span> Excluir · <span style="color:#7aa2f7">?</span>
</pre>

---

### Modo: Configurações

> Aba "Configurações" ativa. Painel único ocupando toda a largura da work area. Sem divisão esquerda/direita.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <strong><span style="color:#7aa2f7">Abditum</span></strong>  <span style="color:#414868">·</span>  <span style="color:#565f89">cofre.abditum</span>                 <span style="color:#565f89">╭ Cofre ╮</span>  <span style="color:#565f89">╭ Modelos ╮</span>  <span style="background:#283457;color:#7aa2f7">╭────────────────╮</span>
 <span style="color:#414868">─────────────────────────────────────────────────────────────────╯</span><span style="background:#283457;color:#7aa2f7"><strong> Configurações  </strong></span><span style="color:#414868">╰────</span>
                                                                                          
  <span style="color:#565f89"><strong>Segurança</strong></span>
  <span style="color:#414868">────────────────────────────────────────────────────────────────────────────────────</span>
  <span style="background:#283457;color:#a9b1d6"><strong>  Timeout de bloqueio                                              5 minutos ▸   </strong></span>
    Confirmar ao excluir                                               <span style="color:#9ece6a">ativado</span>
    Limpar área de transferência após cópia                             30 s
                                                                                          
  <span style="color:#565f89"><strong>Interface</strong></span>
  <span style="color:#414868">────────────────────────────────────────────────────────────────────────────────────</span>
    Tema                                                          Tokyo Night ▸
    Ordenação padrão da árvore                              Alfabética ▸
                                                                                          
  <span style="color:#565f89"><strong>Sobre</strong></span>
  <span style="color:#414868">────────────────────────────────────────────────────────────────────────────────────</span>
    Versão                                                               <span style="color:#565f89">v0.1.0</span>
                                                                                          
 <span style="color:#414868">──────────────────────────────────────────────────────────────────────────────────────</span>
  <span style="color:#7aa2f7">?</span> Ajuda
</pre>

---

## Seção 1 — Header

**Posição:** Linha 1 da tela. 1 linha de altura.

**Responsabilidade:** Contexto global da sessão — qual aplicação, qual cofre, se há alterações, e qual modo está ativo (via abas).

**Dois sub-elementos:**
1. **Linha de título** — nome da app, nome do cofre, indicador de não salvo, e abas (inativas) flutuantes
2. **Linha separadora** — divide visualmente header da work area; a aba ativa é suspensa nesta linha via `╯ ╰`

---

### Estado 1.1 — Sem cofre aberto

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <strong><span style="color:#7aa2f7">Abditum</span></strong>
 <span style="color:#414868">──────────────────────────────────────────────────────────────────────────────────────────</span>
</pre>

Sem cofre = sem abas. A linha separadora é simples e sem nenhum conector.

---

### Estado 1.2 — Cofre aberto, salvo

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <strong><span style="color:#7aa2f7">Abditum</span></strong>  <span style="color:#414868">·</span>  <span style="color:#565f89">cofre.abditum</span>              <span style="color:#565f89">╭ Cofre ╮</span>  <span style="color:#565f89">╭ Modelos ╮</span>  <span style="color:#565f89">╭ Configurações ╮</span>
 <span style="color:#414868">──────────────────────────────────────────────────────────────────────────────────────────</span>
</pre>

Nenhuma aba está em destaque — estado impossível em operação normal (sempre há um modo ativo quando cofre está aberto), mas mostrado aqui para ilustrar a anatomia base.

---

### Estado 1.3 — Cofre aberto, com alterações não salvas

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <strong><span style="color:#7aa2f7">Abditum</span></strong>  <span style="color:#414868">·</span>  <span style="color:#565f89">cofre.abditum</span>  <span style="color:#e0af68">•</span>          <span style="color:#565f89">╭ Cofre ╮</span>  <span style="color:#565f89">╭ Modelos ╮</span>  <span style="color:#565f89">╭ Configurações ╮</span>
 <span style="color:#414868">──────────────────────────────────────────────────────────────────────────────────────────</span>
</pre>

O `•` aparece imediatamente após o nome do cofre, em `semantic.warning`. Desaparece após `^S` bem-sucedido.

---

### Estado 1.4 — Modo Cofre ativo (aba "Cofre" suspensa)

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <strong><span style="color:#7aa2f7">Abditum</span></strong>  <span style="color:#414868">·</span>  <span style="color:#565f89">cofre.abditum</span>  <span style="color:#e0af68">•</span>          <span style="background:#283457;color:#7aa2f7">╭────────╮</span>  <span style="color:#565f89">╭ Modelos ╮</span>  <span style="color:#565f89">╭ Configurações ╮</span>
 <span style="color:#414868">─────────────────────────────────────────────────╯</span><span style="background:#283457;color:#7aa2f7"><strong> Cofre  </strong></span><span style="color:#414868">╰──────────────────────────────────</span>
</pre>

A aba ativa usa `╯ ╰` para "pousar" na linha separadora — conecta visualmente o tab à estrutura sem adicionar linha extra.

---

### Estado 1.5 — Modo Modelos ativo

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <strong><span style="color:#7aa2f7">Abditum</span></strong>  <span style="color:#414868">·</span>  <span style="color:#565f89">cofre.abditum</span>              <span style="color:#565f89">╭ Cofre ╮</span>  <span style="background:#283457;color:#7aa2f7">╭──────────╮</span>  <span style="color:#565f89">╭ Configurações ╮</span>
 <span style="color:#414868">──────────────────────────────────────────────────────────╯</span><span style="background:#283457;color:#7aa2f7"><strong> Modelos  </strong></span><span style="color:#414868">╰──────────────────────</span>
</pre>

---

### Estado 1.6 — Modo Configurações ativo

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <strong><span style="color:#7aa2f7">Abditum</span></strong>  <span style="color:#414868">·</span>  <span style="color:#565f89">cofre.abditum</span>              <span style="color:#565f89">╭ Cofre ╮</span>  <span style="color:#565f89">╭ Modelos ╮</span>  <span style="background:#283457;color:#7aa2f7">╭────────────────╮</span>
 <span style="color:#414868">───────────────────────────────────────────────────────────────────╯</span><span style="background:#283457;color:#7aa2f7"><strong> Configurações  </strong></span><span style="color:#414868">╰</span>
</pre>

---

### Truncamento do nome do cofre

O espaço disponível para o nome do arquivo na linha de título é restrito, pois as abas ocupam espaço fixo à direita. O componente deve calcular esse espaço em tempo real e truncar o nome quando necessário.

**Espaço disponível** (em colunas):

```
disponível = largura_terminal − len("  Abditum  ·  ") − len("  •") [se dirty] − len(bloco_abas)
```

O bloco de abas tem largura fixa para cada modo:

| Modo ativo | Bloco de abas (aprox.) |
|---|---|
| Cofre | `╭────────╮  ╭ Modelos ╮  ╭ Configurações ╮` |
| Modelos | `╭ Cofre ╮  ╭──────────╮  ╭ Configurações ╮` |
| Configurações | `╭ Cofre ╮  ╭ Modelos ╮  ╭────────────────╮` |

**Regras de prioridade (o que cede espaço primeiro):**

1. O nome do cofre — truncado antes de qualquer outro elemento
2. O separador `·` e o indicador `•` — preservados enquanto houver espaço
3. As abas — nunca truncadas

**Algoritmo de truncamento:**

1. Se o nome completo cabe → exibir como está
2. Se não cabe → truncar o radical (parte antes de `.abditum`), preservar a extensão:
   - `{radical[0..n]}….abditum` onde `n` é calculado para que o total caiba no espaço disponível
3. Se nem `….abditum` (9 colunas) cabe → exibir apenas `…`

**Wireframe — nome truncado (terminal estreito ~80 colunas, modo Cofre):**

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <strong><span style="color:#7aa2f7">Abditum</span></strong>  <span style="color:#414868">·</span>  <span style="color:#565f89">meu-cofre-pe….abditum</span>  <span style="color:#e0af68">•</span>  <span style="background:#283457;color:#7aa2f7">╭────────╮</span>  <span style="color:#565f89">╭ Modelos ╮</span>  <span style="color:#565f89">╭ Configurações ╮</span>
 <span style="color:#414868">─────────────────────────────────────────────────╯</span><span style="background:#283457;color:#7aa2f7"><strong> Cofre  </strong></span><span style="color:#414868">╰──────────────────────────────────</span>
</pre>

> O radical `meu-cofre-pessoal` foi truncado para `meu-cofre-pe…` mantendo `.abditum` legível.

---

### Eventos que afetam o Header

| Evento | Mudança visual |
|---|---|
| Cofre aberto com sucesso | Aparece `· cofre.abditum` e as 3 abas |
| Cofre fechado / bloqueado | Desaparece nome do cofre e abas |
| Alteração em memória (`IsDirty() = true`) | Aparece `•` em `semantic.warning` |
| Salvar com sucesso (`IsDirty() = false`) | Desaparece `•` |
| Navegação entre modos (F201 Cofre / F202 Modelos / F203 Settings) | Aba ativa muda; nova aba suspensa na linha separadora |

### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| `Abditum` | `accent.primary` | **bold** |
| `·` separador nome/cofre | `border.default` | — |
| `cofre.abditum` (nome do arquivo) | `text.secondary` | — |
| `•` indicador não salvo | `semantic.warning` | — |
| Aba ativa — borda `╭───╮` | `accent.primary` | — |
| Aba ativa — fundo e texto | `special.highlight` / `accent.primary` | **bold** |
| Aba inativa | `text.secondary` | — |
| Linha separadora principal | `border.default` | — |
| `┬` e `┴` (juntores da linha separadora) | `border.default` | — |

---

## Seção 2 — Barra de Mensagens

**Posição:** Sobreposta à última linha da work area — **não** reserva linha própria.
**Largura:** ~95% da largura do terminal. Trunca com `…` se necessário.
**Conteúdo:** Símbolo + texto. Uma mensagem por vez — nova mensagem substitui a anterior.

---

### MsgSuccess — Operação concluída

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
 <span style="color:#9ece6a">✓  Gmail copiado para a área de transferência</span>
</pre>

`semantic.success` — TTL 2-3 s, não responde a input.

---

### MsgInfo — Informação

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
 <span style="color:#7dcfff">ℹ  Cofre criado em /home/user/documentos/pessoal.abditum</span>
</pre>

`semantic.info` — TTL 3 s, não responde a input. Usado para informações neutras que não são confirmação de sucesso nem aviso.

---

### MsgWarn — Atenção requerida

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
 <span style="color:#e0af68">⚠  Cofre será bloqueado em 15 segundos</span>
</pre>

`semantic.warning` — TTL 0 (permanente), `clearOnInput = true`. Desaparece ao próximo evento de teclado ou mouse. Re-emitido a cada tick enquanto condição persistir.

---

### MsgError — Falha

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
 <span style="color:#f7768e"><strong>✗  Falha ao salvar — arquivo em uso por outro processo</strong></span>
</pre>

`semantic.error` + **bold** — TTL 5 s.

---

### MsgBusy — Operação em andamento (spinner)

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
<span style="color:#7aa2f7"> ◐  Salvando cofre...</span>   <span style="color:#565f89">← frame 0 (segundo 0)</span>
<span style="color:#7aa2f7"> ◓  Salvando cofre...</span>   <span style="color:#565f89">← frame 1 (segundo 1)</span>
<span style="color:#7aa2f7"> ◑  Salvando cofre...</span>   <span style="color:#565f89">← frame 2 (segundo 2)</span>
<span style="color:#7aa2f7"> ◒  Salvando cofre...</span>   <span style="color:#565f89">← frame 3 (segundo 3)</span>
</pre>

`accent.primary` — TTL 0, spinner avança 1 frame/segundo sincronizado com tick global. Permanece até ser substituído por MsgSuccess ou MsgError.

---

### MsgHint — Dica contextual de campo

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
 <span style="color:#565f89"><em>•  Use Tab para alternar o foco entre os painéis</em></span>
</pre>

`text.secondary` + *italic* — TTL 0, não desaparece com input. Exibido enquanto o campo estiver em foco; substituído ao navegar para outro campo.

---

### MsgTip — Dica de uso

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
 <span style="color:#565f89"><em>💡 Pressione F17 para copiar um valor sem revelar o campo</em></span>
</pre>

`text.secondary` + *italic* — TTL 0, não desaparece com input. Emitido proativamente pela aplicação para apresentar funcionalidades ao usuário.

> **Nota de layout:** `💡` é um emoji de largura dupla (2 colunas). O cálculo de truncamento deve reservar 2 colunas para o ícone, não 1.

---

### Comportamento temporal e sobreposição

Uma mensagem por vez — uma nova chamada substitui a anterior imediatamente. Não há fila nem stack.

| Tipo | TTL | Desaparece com input | Notas |
|---|---|---|---|
| MsgSuccess | 2–3 s | Não | Confirmação de operação concluída |
| MsgInfo | 3 s | Não | Informação neutra sem conotação de sucesso |
| MsgWarn | 0 (permanente) | **Sim** | Re-emitido a cada tick enquanto condição persistir (ex: bloqueio iminente) |
| MsgError | 5 s | Não | Permanece mesmo com interação |
| MsgBusy | 0 (permanente) | Não | Substituído por MsgSuccess ou MsgError ao concluir |
| MsgHint | 0 (permanente) | Não | Dica de campo; substituído ao navegar para outro campo |
| MsgTip | 0 (permanente) | Não | Dica de uso proativa; substituída pela próxima mensagem |

---

## Seção 3 — Command Bar

**Posição:** Última linha da tela, abaixo da work area.
**Formato de cada ação:** `TECLA Label` — tecla em `accent.primary` bold, label em `text.primary`. Separados por ` · ` (`text.secondary`).

**Princípio:** a command bar exibe apenas ações de caso de uso (F-keys, atalhos de domínio, `^S`). Teclas de navegação universais — `↑↓`, `←→`, `Tab`, `Enter`, `Esc` — não são exibidas porque são senso comum em qualquer TUI. **Exceção:** modais exibem suas opções explicitamente, pois o contexto muda e as escolhas não são óbvias.

---

### Estado 3.1 — Normal (painel de árvore com foco)

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#7aa2f7">F21</span> Novo · <span style="color:#7aa2f7">F22</span> Editar · <span style="color:#7aa2f7">F23</span> Excluir · <span style="color:#7aa2f7">^S</span> Salvar · <span style="color:#7aa2f7">?</span>
</pre>

---

### Estado 3.2 — Com ação desabilitada

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#7aa2f7">F21</span> Novo · <span style="color:#7aa2f7">F22</span> Editar · <span style="color:#3b4261">F23 Excluir</span> · <span style="color:#7aa2f7">^S</span> Salvar · <span style="color:#7aa2f7">?</span>
</pre>

Ação `F23 Excluir` desabilitada: `text.disabled` + dim. Ocorre quando nenhum segredo está selecionado, por exemplo.

---

### Estado 3.3 — Modal ativo (substituição total)

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#7aa2f7">Enter</span> Confirmar · <span style="color:#7aa2f7">Esc</span> Cancelar
</pre>

Durante modal ativo, a command bar exibe **exclusivamente** os `Shortcuts()` do modal do topo da stack. As ações do `ActionManager` ficam invisíveis.

---

### Estado 3.4 — Espaço restrito (terminal estreito)

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#7aa2f7">F21</span> Novo · <span style="color:#7aa2f7">?</span>
</pre>

Ações de menor prioridade (F22, F23, ^S) são ocultadas quando não há espaço. `?` (Ajuda) permanece sempre visível — é via help modal que o usuário descobre as ações ocultas.

---

### Eventos que afetam a Command Bar

| Evento | Mudança |
|---|---|
| Troca de foco entre painéis (`Tab` / `Shift+Tab`) | Ações do painel que recebe foco ficam ativas na command bar |
| `Enter` em segredo/modelo na árvore | Foco transferido para o painel direito; ações mudam para as do painel direito |
| Seleção de item na árvore | Ações de segredo (F16, F17, F22, F23) ficam habilitadas |
| Nenhum item selecionado | Ações de segredo `text.disabled` + dim |
| Modal aberto | Troca para `Shortcuts()` do modal |
| Modal fechado | Volta para ações do `ActionManager` |
| Janela redimensionada | Recalcula quais ações cabem |

---

## Seção 4 — Work Area: Boas-vindas

**Triggger:** Aplicação iniciada sem cofre aberto, ou após fechar/bloquear cofre.
**Interação:** Nenhuma — tela estática. Toda interação via command bar.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
                                                                                          
             <span style="color:#9d7cd8">    ___    __        ___ __                  </span>
             <span style="color:#89ddff">   /   |  / /_  ____/ (_) /___  ______ ___   </span>
             <span style="color:#7aa2f7">  / /| | / __ \/ __  / / __/ / / / __ `__ \  </span>
             <span style="color:#7dcfff"> / ___ |/ /_/ / /_/ / / /_/ /_/ / / / / / / </span>
             <span style="color:#bb9af7">/_/  |_/_.___/\__,_/_/\__/\__,_/_/ /_/ /_/ </span>
                                                                                          
                                          <span style="color:#565f89">v0.1.0</span>
                                                                                          
</pre>

Logo e versão centralizados horizontal e verticalmente na work area via `lipgloss.Place()`.

| Elemento | Token | Atributo |
|---|---|---|
| Logo linha 1 | `#9d7cd8` | — |
| Logo linha 2 | `#89ddff` | — |
| Logo linha 3 | `#7aa2f7` | — |
| Logo linha 4 | `#7dcfff` | — |
| Logo linha 5 | `#bb9af7` | — |
| `v0.1.0` | `text.secondary` | — |

---

## Seção 5 — Painel Esquerdo: Árvore (Modo Cofre)

**Largura:** ~35% da work area.
**Comportamento de foco:** `│` separador vertical em `border.focused` quando este painel tem foco; `border.default` quando sem foco.

### Estados de itens

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#565f89">Normal:</span>
    YouTube

  <span style="color:#565f89">Selecionado (cursor):</span>
  <span style="background:#283457;color:#a9b1d6"><strong>  ► Gmail                    </strong></span>

  <span style="color:#565f89">Favorito (prefixo estrela):</span>
  <span style="color:#bb9af7">★</span> Bradesco

  <span style="color:#565f89">Favorito + Selecionado (cursor substitui prefixo):</span>
  <span style="background:#283457;color:#a9b1d6"><strong>► Bradesco                     </strong></span>

  <span style="color:#565f89">Marcado para exclusão:</span>
    <span style="color:#565f89">✕  <span style="text-decoration:line-through">Gmail</span></span>

  <span style="color:#565f89">Marcado para exclusão + Selecionado:</span>
  <span style="background:#283457;color:#565f89"><strong>✕  <span style="text-decoration:line-through">Gmail</span>                      </strong></span>

  <span style="color:#565f89">Pasta expandida (com filhos):</span>
  <span style="color:#565f89">▼</span> Sites e Apps          <span style="color:#565f89">(5)</span>

  <span style="color:#565f89">Pasta recolhida (com filhos):</span>
  <span style="color:#565f89">▶</span> Google                <span style="color:#565f89">(2)</span>

  <span style="color:#565f89">Pasta vazia:</span>
  <span style="color:#565f89">▷</span> Nova pasta             <span style="color:#565f89">(0)</span>

  <span style="color:#565f89">Segredo normal:</span>
  <span style="color:#565f89">●</span> Gmail

  <span style="color:#565f89">Segredo favoritado:</span>
  <span style="color:#bb9af7">★</span> Bradesco

  <span style="color:#565f89">Item desabilitado (indisponível no contexto):</span>
    <span style="color:#3b4261">Conta Empresa</span>
</pre>

| Estado | Fundo | Texto | Atributo | Notas |
|---|---|---|---|---|
| Normal | `surface.base` | `text.primary` | — | — |
| Selecionado | `special.highlight` | `text.primary` | **bold** | Toda a linha preenche a largura do painel |
| Favorito | `surface.base` | `text.primary` | — | `★` em `accent.secondary` no final da linha |
| Favorito + selecionado | `special.highlight` | `text.primary` | **bold** | `★` em `accent.secondary` preservado |
| Exclusão | `surface.base` | `special.muted` | ~~strikethrough~~ | `✕` em `semantic.error` como prefixo |
| Exclusão + selecionado | `special.highlight` | `special.muted` | ~~strikethrough~~ + **bold** | — |
| Prefixos de pasta `▼▶▷` | — | `text.secondary` | — | Indentação 2 espaços por nível |
| Prefixo de segredo `●` | — | `text.secondary` | — | — |
| Prefixo de segredo favoritado `★` | — | `accent.secondary` | — | — |
| Contadores `(n)` | — | `text.secondary` | — | À direita do nome da pasta |
| Desabilitado | `surface.base` | `text.disabled` | dim | Raramente usado na árvore |

---

### Navegação no Painel Esquerdo

| Tecla | Ação |
|---|---|
| `↑` / `↓` | Move cursor entre itens visíveis |
| `Enter` ou `→` | **Em pasta:** expande/recolhe. **Em segredo/modelo:** abre detalhes no painel direito e transfere foco para ele |
| `←` | Recolhe pasta expandida; sobe para pasta pai |
| `Tab` | Foco → próximo painel focusável (painel direito) |
| `Shift+Tab` | Foco → painel focusável anterior |
| `Home` | Move cursor para o primeiro item |
| `End` | Move cursor para o último item visível |
| `F21` | Novo segredo no contexto atual |
| `F22` | Edita segredo selecionado |
| `F23` | Marca segredo selecionado para exclusão (com DialogAlert) |
| `F27` | Cria nova pasta |
| `F28` | Renomeia pasta selecionada |
| `F31` | Exclui pasta selecionada (com DialogAlert) |

**Indicadores de scroll:**

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#565f89">↑</span>  <span style="color:#565f89">← conteúdo acima</span>
  ▾ Sites e Apps   (5)
  ▸ Google         (2)
  ── Gmail
  ─── YouTube
  <span style="color:#565f89">↓</span>  <span style="color:#565f89">← conteúdo abaixo</span>
</pre>

`↑` e `↓` em `text.secondary`, exibidos quando há conteúdo além da área visível.

---

## Seção 6 — Painel Direito: Detalhe do Segredo

**Largura:** ~65% da work area.
**Comportamento de foco:** `│` separador vertical em `border.focused` quando este painel tem foco; `border.default` quando sem foco.

---

### Estado 6.1 — Nenhum segredo selecionado

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">




                      <span style="color:#565f89"><em>Selecione um segredo para ver os detalhes</em></span>



</pre>

Placeholder centralizado — `text.secondary` + *italic*.

---

### Estado 6.2 — Segredo exibido (leitura)

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  Gmail                                                                <span style="color:#bb9af7">★</span>
  <span style="color:#414868">──────────────────────────────────────────────────────────────────</span>
  <span style="color:#565f89">URL</span>            https://mail.google.com
  <span style="color:#565f89">Usuário</span>        fulano@gmail.com
  <span style="color:#565f89">Senha</span>          <span style="color:#565f89">••••••••</span>                                        <span style="color:#7aa2f7">F16</span>
  <span style="color:#565f89">Token 2FA</span>      <span style="color:#565f89">••••••</span>                                          <span style="color:#7aa2f7">F16</span>

  <span style="color:#565f89">Observação</span>     Conta pessoal principal — criada em 2018
</pre>

| Elemento | Token | Atributo |
|---|---|---|
| Título do segredo | `text.primary` | **bold** |
| `★` (favorito) | `accent.secondary` | — |
| Linha separadora `───` | `border.default` | — |
| Labels (`URL`, `Usuário`, `Senha`…) | `text.secondary` | — |
| Valores de texto | `text.primary` | — |
| Máscaras `••••••••` | `text.secondary` | — |
| `F16` (revelar) | `accent.primary` | — |

> **Comprimento da máscara:** fixo em **8 `•`**, independente do tamanho real da senha — evita vazar o comprimento como informação.

---

### Estado 6.3 — Campo sensível revelado

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  Gmail                                                                <span style="color:#bb9af7">★</span>
  <span style="color:#414868">──────────────────────────────────────────────────────────────────</span>
  <span style="color:#565f89">URL</span>            https://mail.google.com
  <span style="color:#565f89">Usuário</span>        fulano@gmail.com
  <span style="color:#565f89">Senha</span>          minha-senha-secreta-123                         <span style="color:#7aa2f7">F16</span>
  <span style="color:#565f89">Token 2FA</span>      <span style="color:#565f89">••••••</span>                                          <span style="color:#7aa2f7">F16</span>

  <span style="color:#565f89">Observação</span>     Conta pessoal principal — criada em 2018
</pre>

Apenas o campo revelado mostra o valor em `text.primary`. Os outros campos sensíveis permanecem mascarados. `F16` ainda visível para re-ocultar.

---

### Estado 6.4 — Conector de seleção `<╡`

O conector `&lt;╡` aparece no separador vertical (`│`) exatamente na linha do item selecionado na árvore, indicando qual item está sendo detalhado:

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
    ▸ Google             (2)  <span style="color:#7aa2f7">│</span>  URL            https://mail.google.com
<span style="background:#283457;color:#a9b1d6"><strong>  ► Gmail                   </strong></span><span style="color:#7aa2f7">&lt;╡</span>  Senha          <span style="color:#f7768e">••••••••••••</span>                  <span style="color:#7aa2f7">F16</span>
    YouTube                  <span style="color:#7aa2f7">│</span>
</pre>

`&lt;╡` em `accent.primary` — substitui o `│` na linha do item selecionado. O `<` aponta da árvore para o detalhe; `╡` conecta visualmente ao separador.

---

### Estado 6.5 — Segredo em edição (após F22)

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  Gmail  <span style="color:#e0af68">•</span>                                                            <span style="color:#bb9af7">★</span>
  <span style="color:#414868">──────────────────────────────────────────────────────────────────</span>
  <span style="color:#565f89">URL</span>            <span style="color:#414868">╭──────────────────────────────────────────────────╮</span>
                 <span style="color:#7aa2f7">│</span> https://mail.google.com                          <span style="color:#414868">│</span>
                 <span style="color:#414868">╰──────────────────────────────────────────────────╯</span>
  <span style="color:#565f89">Usuário</span>        <span style="color:#3b4261">fulano@gmail.com</span>
  <span style="color:#565f89">Senha</span>          <span style="color:#3b4261">••••••••••••</span>

  <span style="color:#565f89">Observação</span>     <span style="color:#3b4261">Conta pessoal principal — criada em 2018</span>
</pre>

No modo edição:
- Campo ativo: borda `border.focused`, cursor visível
- Label do campo ativo: `accent.primary` + **bold**
- Outros campos: `text.disabled` + dim (indicam que apenas um campo está editável por vez)
- `•` aparece no título do segredo em `semantic.warning` (indicador de alteração local)

---

### Estado 6.6 — Campo com erro de validação

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#7aa2f7"><strong>URL</strong></span>            <span style="color:#f7768e">╭──────────────────────────────────────────────────╮</span>
                 <span style="color:#f7768e">│</span> não-é-uma-url-válida                             <span style="color:#f7768e">│</span>
                 <span style="color:#f7768e">╰──────────────────────────────────────────────────╯</span>
                 <span style="color:#f7768e"><em>URL inválida — deve começar com http:// ou https://</em></span>
</pre>

Borda em `semantic.error`; mensagem de erro inline abaixo do campo em `semantic.error` + *italic*. Campo não fecha ao pressionar Enter enquanto inválido.

---

### Eventos no Painel Direito

| Tecla / Evento | Ação |
|---|---|
| `F16` | Revela/oculta campo sensível em foco |
| `F17` | Copia valor do campo em foco → MsgInfo "Copiado" |
| `F22` | Ativa modo edição — coloca cursor no primeiro campo |
| `↑` / `↓` | Navega entre campos (no modo edição) |
| `Tab` (modo edição) | Avança para o próximo campo; no último campo, foco → painel esquerdo |
| `Shift+Tab` (modo edição) | Retrocede para o campo anterior; no primeiro campo, foco → painel esquerdo |
| `Tab` (modo leitura) | Foco → próximo painel focusável |
| `Shift+Tab` (modo leitura) | Foco → painel focusável anterior |
| `^S` | Salva alterações → MsgBusy → MsgInfo/MsgError |
| `Esc` (modo edição) | Descarta alterações — abre DialogAlert se houver mudanças |
| Timeout de reveal | Campo sensível revelado volta a ser mascarado automaticamente |

---

## Seção 7 — Work Area: Modelos

Layout idêntico ao modo Cofre (35/65). Diferenças:

| Aspecto | Modo Cofre | Modo Modelos |
|---|---|---|
| Painel esquerdo | Árvore hierárquica (pastas + segredos) | Lista plana de modelos |
| Favoritos | Sim — prefixo `★` em `accent.secondary` | Não |
| Campos sensíveis | Sim — mascarados com `••••` | Não — modelos não têm dados reais |
| Conector `&lt;╡` | Sim | Não necessário (lista plana é direta) |
| Prefixos de árvore `▼▶▷●★` | Sim | Não (usar `●` para cada modelo) |
| Hierarquia | Múltiplos níveis | 1 nível (lista) |
| Painel direito | Campos do segredo selecionado | Campos do template: Nome, Tipo, Obrigatório |

---

## Seção 8 — Work Area: Configurações

Painel único, largura total da work area. Sem divisão esquerda/direita.

### Anatomia

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#565f89"><strong>Segurança</strong></span>
  <span style="color:#414868">────────────────────────────────────────────────────────────────────────────────────</span>
  <span style="background:#283457;color:#a9b1d6"><strong>  Timeout de bloqueio automático                                 5 minutos  ▸  </strong></span>
    Confirmar ao excluir segredo                                         <span style="color:#9ece6a">ativado</span>
    Limpar área de transferência após cópia                               30 s
    Revelar campo sensível por                                            15 s

  <span style="color:#565f89"><strong>Interface</strong></span>
  <span style="color:#414868">────────────────────────────────────────────────────────────────────────────────────</span>
    Tema visual                                                    Tokyo Night  ▸
    <span style="color:#3b4261">Ordenação da árvore                                            Alfabética  </span>  <span style="color:#565f89"><em>← não implementado</em></span>

  <span style="color:#565f89"><strong>Sobre</strong></span>
  <span style="color:#414868">────────────────────────────────────────────────────────────────────────────────────</span>
    Versão                                                               <span style="color:#565f89">v0.1.0</span>
    Repositório                                                <span style="color:#7aa2f7">github.com/…</span>
</pre>

| Elemento | Token | Atributo |
|---|---|---|
| Título de grupo (`Segurança`, `Interface`) | `text.secondary` | **bold** |
| Separador de grupo `────` | `border.default` | — |
| Item selecionado (cursor) | `special.highlight` | **bold** |
| Label da configuração | `text.primary` | — |
| Valor da configuração | `text.primary` | — |
| Valor booleano ativado | `semantic.success` | — |
| Valor numérico / seleção | `text.primary` | — |
| `▸` (indica mais opções) | `text.secondary` | — |
| Item desabilitado | `text.disabled` | dim |

---

## Seção 9 — Modais

Todos os modais:
- Renderizados **sobre** todo o conteúdo via `lipgloss.Place()` (centralizado)
- Estilo de borda: **Rounded** (`╭╮╰╯│─`)
- Fundo interno: `surface.raised` (`#24283b`)
- Command bar troca para `Shortcuts()` do modal enquanto aberto

---

### Modal 9.1 — DialogQuestion

> Decisão neutra. Usado para: salvar/descartar/cancelar; sobrescrever arquivo.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
                    <span style="color:#7aa2f7">╭─────────────────────────────────────────╮</span>
                    <span style="color:#7aa2f7">│</span>  <span style="background:#24283b;color:#7aa2f7"><strong>?  Alterações não salvas</strong></span>                <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  Deseja salvar antes de sair?           <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>   <span style="color:#7aa2f7"><strong>[S] Salvar</strong></span>   [N] Descartar   [Esc] Voltar  <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">╰─────────────────────────────────────────╯</span>
</pre>

- Borda: `border.focused` (`accent.primary`)
- Título: `accent.primary` + **bold** + símbolo `?`
- Opção default `[S] Salvar`: `accent.primary` + **bold** (acionada por Enter)
- Opções neutras: `text.primary`
- Opção cancel `[Esc] Voltar`: acionada por Esc

---

### Modal 9.2 — DialogAlert

> Ação destrutiva ou irreversível. Usado para: excluir segredo; excluir pasta; descartar alterações.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
                    <span style="color:#e0af68">╭─────────────────────────────────────────╮</span>
                    <span style="color:#e0af68">│</span>  <span style="background:#24283b;color:#e0af68"><strong>⚠  Excluir segredo</strong></span>                      <span style="color:#e0af68">│</span>
                    <span style="color:#e0af68">│</span>                                         <span style="color:#e0af68">│</span>
                    <span style="color:#e0af68">│</span>  <strong>Gmail</strong> será excluído permanentemente.   <span style="color:#e0af68">│</span>
                    <span style="color:#e0af68">│</span>  Esta ação não pode ser desfeita.       <span style="color:#e0af68">│</span>
                    <span style="color:#e0af68">│</span>                                         <span style="color:#e0af68">│</span>
                    <span style="color:#e0af68">│</span>   <span style="color:#f7768e"><strong>[S] Excluir</strong></span>          [N] Cancelar        <span style="color:#e0af68">│</span>
                    <span style="color:#e0af68">╰─────────────────────────────────────────╯</span>
</pre>

- Borda: `semantic.warning`
- Título: `semantic.warning` + **bold** + `⚠`
- Opção destrutiva default `[S] Excluir`: `semantic.error` + **bold**
- Opção cancel: `text.primary`

---

### Modal 9.3 — DialogInfo

> Informação que requer reconhecimento. Sem opções — apenas dismiss.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
                    <span style="color:#7dcfff">╭─────────────────────────────────────────╮</span>
                    <span style="color:#7dcfff">│</span>  <span style="background:#24283b;color:#7dcfff"><strong>ℹ  Conflito detectado</strong></span>                   <span style="color:#7dcfff">│</span>
                    <span style="color:#7dcfff">│</span>                                         <span style="color:#7dcfff">│</span>
                    <span style="color:#7dcfff">│</span>  O arquivo foi modificado externamente. <span style="color:#7dcfff">│</span>
                    <span style="color:#7dcfff">│</span>  Escolha como prosseguir na próxima     <span style="color:#7dcfff">│</span>
                    <span style="color:#7dcfff">│</span>  tela.                                  <span style="color:#7dcfff">│</span>
                    <span style="color:#7dcfff">│</span>                                         <span style="color:#7dcfff">│</span>
                    <span style="color:#7dcfff">│</span>              <span style="color:#7aa2f7"><strong>[Enter] OK</strong></span>                  <span style="color:#7dcfff">│</span>
                    <span style="color:#7dcfff">╰─────────────────────────────────────────╯</span>
</pre>

- Borda: `semantic.info`
- Título: `semantic.info` + **bold** + `ℹ`
- `[Enter] OK`: `accent.primary` + **bold** (único atalho)

---

### Modal 9.4 — PasswordEntry (1 campo)

> Entrada de senha única. Usado em: F1 Abrir cofre.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
                    <span style="color:#7aa2f7">╭─────────────────────────────────────────╮</span>
                    <span style="color:#7aa2f7">│</span>  <span style="background:#24283b;color:#a9b1d6"><strong>Senha mestra</strong></span>                            <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#565f89">Senha</span>                                   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">╭───────────────────────────────────╮</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">│</span> <span style="color:#565f89">•••••••••••••••</span>▌                  <span style="color:#7aa2f7">│</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">╰───────────────────────────────────╯</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">╰─────────────────────────────────────────╯</span>

  <span style="color:#7aa2f7">Enter</span> Confirmar · <span style="color:#7aa2f7">Esc</span> Cancelar
</pre>

- Borda externa: `border.focused`
- Borda do campo de input: `border.focused`
- Caracteres digitados: `text.secondary` (máscarados com `•`)
- Cursor `▌`: `text.primary` piscante
- Label `Senha`: `text.secondary`

---

### Modal 9.5 — PasswordCreate (2 campos)

> Criação de senha com confirmação. Usado em: F2 Criar cofre, F11 Alterar senha.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
                    <span style="color:#7aa2f7">╭─────────────────────────────────────────╮</span>
                    <span style="color:#7aa2f7">│</span>  <span style="background:#24283b;color:#a9b1d6"><strong>Definir senha mestra</strong></span>                   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7"><strong>Nova senha</strong></span>                              <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">╭───────────────────────────────────╮</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">│</span> <span style="color:#565f89">••••••••••••</span>▌                      <span style="color:#7aa2f7">│</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">╰───────────────────────────────────╯</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#565f89">Confirmação</span>                             <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#414868">╭───────────────────────────────────╮</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#414868">╰───────────────────────────────────╯</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">╰─────────────────────────────────────────╯</span>
</pre>

**Estado: confirmação não confere** (usuário pressionou Enter no 2º campo):

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
                    <span style="color:#7aa2f7">│</span>  <span style="color:#565f89">Confirmação</span>                             <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#f7768e">╭───────────────────────────────────╮</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#f7768e">│</span> <span style="color:#565f89">•••••••••</span>▌                         <span style="color:#f7768e">│</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#f7768e">╰───────────────────────────────────╯</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#f7768e"><em>As senhas não coincidem</em></span>                <span style="color:#7aa2f7">│</span>
</pre>

Campo inválido: borda `semantic.error`; mensagem inline abaixo: `semantic.error` + *italic*. Modal não fecha — usuário corrige e tenta novamente.

---

### Modal 9.6 — TextInput

> Entrada de texto livre com validação. Usado em: F18 nome do segredo, F27 criar pasta, F33 renomear modelo.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
                    <span style="color:#7aa2f7">╭─────────────────────────────────────────╮</span>
                    <span style="color:#7aa2f7">│</span>  <span style="background:#24283b;color:#a9b1d6"><strong>Nome do segredo</strong></span>                        <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">╭───────────────────────────────────╮</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">│</span> <span style="color:#565f89"><em>ex: Gmail pessoal</em></span>                  <span style="color:#7aa2f7">│</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>  <span style="color:#7aa2f7">╰───────────────────────────────────╯</span>   <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">│</span>                                         <span style="color:#7aa2f7">│</span>
                    <span style="color:#7aa2f7">╰─────────────────────────────────────────╯</span>
</pre>

Placeholder: `text.secondary` + *italic*. Desaparece ao digitar.

---

### Modal 9.7 — Select (lista navegável)

> Seleção em lista. Usado em: F18 escolher modelo, F25 mover segredo, F29 mover pasta.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
                    <span style="color:#414868">╭─────────────────────────────────────────╮</span>
                    <span style="color:#414868">│</span>  <span style="background:#24283b;color:#a9b1d6"><strong>Escolher modelo</strong></span>                        <span style="color:#414868">│</span>
                    <span style="color:#414868">│</span>  <span style="color:#414868">─────────────────────────────────────</span>  <span style="color:#414868">│</span>
                    <span style="color:#414868">│</span>  Login Padrão                           <span style="color:#414868">│</span>
                    <span style="color:#414868">│</span>  <span style="background:#283457;color:#a9b1d6"><strong>  Cartão de Crédito                  </strong></span>  <span style="color:#414868">│</span>
                    <span style="color:#414868">│</span>  Conta Bancária                         <span style="color:#414868">│</span>
                    <span style="color:#414868">│</span>  SSH Key                                <span style="color:#414868">│</span>
                    <span style="color:#414868">│</span>  Wi-Fi                                  <span style="color:#414868">│</span>
                    <span style="color:#414868">╰─────────────────────────────────────────╯</span>

  <span style="color:#7aa2f7">Enter</span> Confirmar · <span style="color:#7aa2f7">Esc</span> Cancelar
</pre>

Item selecionado: `special.highlight` + **bold**. Borda: `border.default`.

---

### Modal 9.8 — FilePicker

> Navegação de diretórios split-panel: árvore de diretórios à esquerda, arquivos do diretório selecionado à direita. `Tab` alterna o foco entre os dois painéis.
>
> **Dois modos:**
> - **open** — seleciona arquivo existente (F1 Abrir cofre, F13 Importar)
> - **save** — escolhe destino e nome para escrita (F2 Criar cofre, F9 Salvar como, F12 Exportar)

---

#### 9.8.1 — Modo Open: foco na árvore de diretórios (painel esquerdo)

O foco inicial é sempre no painel esquerdo. O painel direito atualiza automaticamente ao mover o cursor na árvore.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#414868">╭───────────────────────────────────────────────────────────────────────╮</span>
  <span style="color:#414868">│</span>  <strong>Abrir cofre</strong>                                                        <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">─── Diretórios ─────────────────────</span><span style="color:#414868">┬─── /home/usuario/cofres ──────────</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">▸ /</span>                                 <span style="color:#414868">│</span>  cofre.abditum                     <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">  ▸ home/</span>                           <span style="color:#414868">│</span>  pessoal.abditum                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">    ▾ usuario/</span>                      <span style="color:#414868">│</span>  trabalho.abditum                  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="background:#283457;color:#bb9af7"><strong>      ► cofres/                   </strong></span>  <span style="color:#414868">│</span>  <span style="color:#3b4261">relatorio.pdf</span>                     <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">        backup/</span>                     <span style="color:#414868">│</span>  <span style="color:#3b4261">notas.txt</span>                         <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">      Documents/</span>                    <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">      Downloads/</span>                    <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#414868">────────────────────────────────────┴───────────────────────────────</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">╰───────────────────────────────────────────────────────────────────────╯</span>

  <span style="color:#7aa2f7">→</span> Expandir · <span style="color:#7aa2f7">←</span> Recolher · <span style="color:#7aa2f7">Tab</span> Arquivos · <span style="color:#7aa2f7">Esc</span> Cancelar
</pre>

O título do painel direito mostra o caminho do diretório selecionado atualmente na árvore — atualiza em tempo real ao navegar.

---

#### 9.8.2 — Modo Open: foco nos arquivos (painel direito)

Após `Tab`, o foco passa para o painel direito. O separador vertical troca de `border.default` para `border.focused`.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#414868">╭───────────────────────────────────────────────────────────────────────╮</span>
  <span style="color:#414868">│</span>  <strong>Abrir cofre</strong>                                                        <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#414868">─── Diretórios ─────────────────────</span><span style="color:#7aa2f7">┬─── /home/usuario/cofres ──────────</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">▸ /</span>                                 <span style="color:#7aa2f7">│</span>  <span style="background:#283457;color:#a9b1d6"><strong>  cofre.abditum                   </strong></span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">  ▸ home/</span>                           <span style="color:#7aa2f7">│</span>  pessoal.abditum                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">    ▾ usuario/</span>                      <span style="color:#7aa2f7">│</span>  trabalho.abditum                  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">      ► cofres/</span>                     <span style="color:#7aa2f7">│</span>  <span style="color:#3b4261">relatorio.pdf</span>                     <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">        backup/</span>                     <span style="color:#7aa2f7">│</span>  <span style="color:#3b4261">notas.txt</span>                         <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">      Documents/</span>                    <span style="color:#7aa2f7">│</span>                                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">      Downloads/</span>                    <span style="color:#7aa2f7">│</span>                                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#414868">────────────────────────────────────┴───────────────────────────────</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">╰───────────────────────────────────────────────────────────────────────╯</span>

  <span style="color:#7aa2f7">Enter</span> Selecionar · <span style="color:#7aa2f7">Tab</span> Diretórios · <span style="color:#7aa2f7">Esc</span> Cancelar
</pre>

O separador `┬` e a linha vertical tornam-se `border.focused` (`accent.primary`) quando o painel direito tem foco — mesmo padrão visual dos painéis do modo Cofre.

---

#### 9.8.3 — Modo Open: diretório sem arquivos .abditum

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#414868">│</span>  <span style="color:#414868">─── Diretórios ─────────────────────</span><span style="color:#414868">┬─── /home/usuario/Downloads ───────</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">  ▸ home/</span>                           <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">    ▾ usuario/</span>                      <span style="color:#414868">│</span>   <span style="color:#565f89"><em>Nenhum arquivo .abditum aqui</em></span>    <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="background:#283457;color:#bb9af7"><strong>      ► Downloads/                </strong></span>  <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>
</pre>

Placeholder no painel direito em `text.secondary` + *italic*. O painel direito permanece sem cursor e `Tab` para ele é ignorado enquanto não houver arquivos selecionáveis.

---

#### 9.8.4 — Modo Save: foco na árvore

No modo save, o rodapé do modal inclui um campo de nome abaixo dos painéis.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#414868">╭───────────────────────────────────────────────────────────────────────╮</span>
  <span style="color:#414868">│</span>  <strong>Salvar cofre como</strong>                                                  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">─── Diretórios ─────────────────────</span><span style="color:#414868">┬─── /home/usuario/cofres ──────────</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">▸ /</span>                                 <span style="color:#414868">│</span>  cofre.abditum                     <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">  ▸ home/</span>                           <span style="color:#414868">│</span>  pessoal.abditum                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">    ▾ usuario/</span>                      <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="background:#283457;color:#bb9af7"><strong>      ► cofres/                   </strong></span>  <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#bb9af7">      Documents/</span>                    <span style="color:#414868">│</span>                                   <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#414868">────────────────────────────────────┴───────────────────────────────</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#565f89">Nome do arquivo</span>                                                    <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#414868">╭───────────────────────────────────────────────────────────────╮</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#414868">│</span> <span style="color:#565f89"><em>nome-do-cofre</em></span>                                                 <span style="color:#414868">│</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#414868">╰───────────────────────────────────────────────────────────────╯</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#565f89">.abditum será adicionado automaticamente</span>                          <span style="color:#414868">│</span>
  <span style="color:#414868">╰───────────────────────────────────────────────────────────────────────╯</span>

  <span style="color:#7aa2f7">Tab</span> Área · <span style="color:#7aa2f7">Enter</span> Salvar · <span style="color:#7aa2f7">Esc</span> Cancelar
</pre>

---

#### 9.8.5 — Modo Save: foco no campo de nome

`Tab` no painel direito (ou da árvore se não houver arquivos) avança para o campo de nome. Borda troca para `border.focused`.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#414868">│</span>  <span style="color:#565f89">Nome do arquivo</span>                                                    <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">╭───────────────────────────────────────────────────────────────╮</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">│</span> cofre-pessoal▌                                                <span style="color:#7aa2f7">│</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">╰───────────────────────────────────────────────────────────────╯</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#565f89">.abditum será adicionado automaticamente</span>                          <span style="color:#414868">│</span>
</pre>

**Pre-fill:** em "Salvar como" (F9), o campo é pré-preenchido com o nome atual sem extensão. Em "Criar cofre" (F2), o campo inicia vazio com placeholder.

---

#### 9.8.6 — Modo Save: conflito de nome

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#414868">│</span>  <span style="color:#565f89">Nome do arquivo</span>                                                    <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#f7768e">╭───────────────────────────────────────────────────────────────╮</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#f7768e">│</span> cofre▌                                                        <span style="color:#f7768e">│</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#f7768e">╰───────────────────────────────────────────────────────────────╯</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#f7768e"><em>cofre.abditum já existe neste diretório</em></span>                           <span style="color:#414868">│</span>
</pre>

Borda `semantic.error`, mensagem inline `semantic.error` + *italic*. `Enter` com conflito abre `DialogAlert` "Sobrescrever arquivo?".

---

#### 9.8.7 — Tokens por elemento

| Elemento | Token | Atributo |
|---|---|---|
| Título do modal | `text.primary` | **bold** |
| Cabeçalho "Diretórios" (painel esquerdo com foco) | `border.focused` | — |
| Cabeçalho "Diretórios" (painel esquerdo sem foco) | `border.default` | — |
| Separador `┬` e `│` (painel direito com foco) | `border.focused` | — |
| Separador `┬` e `│` (painel direito sem foco) | `border.default` | — |
| Diretórios na árvore | `accent.secondary` | — |
| Prefixos de pasta `▼` `▶` `▷` na árvore | `accent.secondary` | — |
| Prefixo de segredo `●` | `accent.secondary` | — |
| Diretório selecionado (cursor) | `special.highlight` | **bold** (texto `accent.secondary`) |
| Arquivos `.abditum` | `text.primary` | — |
| Arquivos sem extensão reconhecida | `text.disabled` | dim |
| Arquivo selecionado (cursor) | `special.highlight` | **bold** |
| Caminho no cabeçalho do painel direito | `text.secondary` | — |
| Campo de nome (sem foco) | `border.default` | — |
| Campo de nome (com foco) | `border.focused` | — |
| Campo de nome (conflito) | `semantic.error` | — |
| Erro inline | `semantic.error` | *italic* |
| Hint `.abditum` | `text.secondary` | — |
| Placeholder no campo vazio | `text.secondary` | *italic* |

---

#### 9.8.8 — Ciclo de foco e navegação

**Modo open** — ciclo `Tab`: `árvore → arquivos → árvore → …`

**Modo save** — ciclo `Tab`: `árvore → arquivos → campo de nome → árvore → …`

> Se o diretório selecionado não tiver arquivos `.abditum`, o painel direito é pulado no ciclo de Tab.

| Tecla | Painel com foco | Ação |
|---|---|---|
| `↑` / `↓` | Árvore | Move cursor entre diretórios |
| `→` / `Enter` | Árvore (cursor em `▸`) | Expande diretório |
| `←` | Árvore (cursor em `▾`) | Recolhe diretório |
| `←` | Árvore (cursor em folha) | Move cursor para o pai |
| `Tab` | Árvore | Foco → painel de arquivos (se houver) ou campo de nome (modo save) |
| `↑` / `↓` | Arquivos | Move cursor entre arquivos |
| `Enter` | Arquivos (modo open) | Confirma seleção, fecha modal |
| `Enter` | Arquivos (modo open, arquivo dim) | Ignorado |
| `Tab` | Arquivos | Foco → campo de nome (modo save) ou árvore (modo open) |
| Qualquer digitação | Campo de nome | Edita o nome |
| `Enter` | Campo de nome | Confirma save (ou DialogAlert se conflito) |
| `Tab` / `Shift+Tab` | Campo de nome | Foco → árvore |
| `Home` / `End` | Árvore ou arquivos | Primeiro / último item visível |
| `Esc` | Qualquer | Cancela e fecha o modal |

---

### Modal 9.9 — Help

> Lista todas as ações do `ActionManager`, agrupadas. Acionado por `?` em qualquer contexto.

<pre style="font-family:monospace;background:#1a1b26;color:#a9b1d6;padding:1em;border-radius:6px;line-height:1.5">
  <span style="color:#414868">╭────────────────────────────────────────────────────────────────────╮</span>
  <span style="color:#414868">│</span>  <strong>Ajuda — Atalhos e Ações</strong>                                         <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#414868">──────────────────────────────────────────────────────────────────</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#565f89"><strong>Navegação</strong></span>                                                        <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">↑↓</span>          Mover cursor na lista                                <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">→ / Enter</span>   Expandir pasta ou selecionar segredo                 <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">←</span>           Recolher pasta ou subir para pasta pai               <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">Tab</span>         Alternar foco entre painéis                          <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>                                                                    <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#565f89"><strong>Segredo</strong></span>                                                          <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">F16</span>         Revelar / ocultar campo sensível                     <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">F17</span>         Copiar valor para área de transferência              <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">F21</span>         Novo segredo                                         <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">F22</span>         Editar segredo                                       <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">F23</span>         Excluir segredo                                      <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>                                                                    <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#565f89"><strong>Cofre</strong></span>                                                            <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">^S</span>          Salvar cofre                                         <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">F5</span>          Sair (salva se necessário)                           <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#7aa2f7">?</span>           Esta ajuda                                           <span style="color:#414868">│</span>
  <span style="color:#414868">│</span>  <span style="color:#414868">─────────────────────────────────────────────────── ↓ mais ──────</span>  <span style="color:#414868">│</span>
  <span style="color:#414868">╰────────────────────────────────────────────────────────────────────╯</span>

  <span style="color:#7aa2f7">↑↓</span> Scroll · <span style="color:#7aa2f7">Esc</span> Fechar
</pre>

| Elemento | Token | Atributo |
|---|---|---|
| Título do modal | `text.primary` | **bold** |
| Títulos de grupo (`Navegação`, `Segredo`…) | `text.secondary` | **bold** |
| Teclas (`↑↓`, `F16`, `^S`…) | `accent.primary` | — |
| Labels das ações | `text.primary` | — |
| `↓ mais ───` (indicador de scroll) | `text.secondary` | — |
| Borda | `border.default` | — |

---

## Seção 10 — Eventos e Reações

Mapeamento de eventos → componente afetado → mudança visual resultante.

### Eventos Globais (qualquer estado)

| Tecla / Evento | Componente | Reação |
|---|---|---|
| `?` | Modal stack | Abre Help modal; command bar troca para `[Esc] Fechar` |
| `^Q` | Root | Abre DialogAlert se cofre dirty; senão fecha o app |
| Tick 1s | Root | Reavalia timeout de bloqueio; se ≥ threshold → MsgWarn "Bloqueio iminente" |
| Tick 1s + MsgBusy ativo | Barra de Mensagens | Avança frame do spinner (◐ → ◓ → ◑ → ◒) |
| Resize do terminal | Todos | Relayout completo; command bar recalcula quais ações cabem |
| `$NO_COLOR` definido | Todos | Lipgloss desativa todas as cores — interface monocromática funcional |

### Eventos com Cofre Aberto

| Tecla / Evento | Componente | Reação |
|---|---|---|
| `^S` | Root | MsgBusy "Salvando…" → save → MsgSuccess "Salvo" ou MsgError; remove `•` do header |
| `Tab` | Painéis | Alterna foco; `│` troca de `border.focused` para `border.default` e vice-versa; command bar atualiza |
| `↑` / `↓` | Painel com foco | Move cursor; item selecionado muda (fundo `special.highlight`); painel direito atualiza |
| `Enter` / `→` na árvore | Árvore | Expande pasta OU exibe segredo no painel direito; conector `&lt;╡` aparece na linha do segredo |
| `F22` no segredo | Detalhe | Ativa modo edição; campos ficam dim exceto o ativo; `•` aparece no título e no header |
| `Esc` no modo edição | Detalhe | Se não há mudanças: volta ao modo leitura. Se há mudanças: DialogAlert |
| `F16` | Detalhe | Campo sensível alterna entre mascarado/revelado; inicia timer de auto-ocultação |
| Timeout reveal | Detalhe | Campo volta a ser mascarado automaticamente |
| `F17` | Detalhe | Copia valor para clipboard → MsgInfo "Copiado" (TTL 2s) |
| `F23` | Árvore | Abre DialogAlert; confirmado → item fica `✕` + `special.muted` + strikethrough + `•` no header |
| Bloqueio por timeout | Root | Wipe de memória sensitiva; transição para Work Area Boas-vindas; header perde nome do cofre e abas |
| Qualquer input + MsgWarn clearOnInput | Barra de Mensagens | MsgWarn desaparece imediatamente |

### Eventos em Modais

| Tecla / Evento | Componente | Reação |
|---|---|---|
| `Enter` | Modal do topo | Aciona opção `Default` |
| `Esc` | Modal do topo | Aciona opção `Cancel`; se não houver, fecha modal |
| Atalho da opção (`[S]`, `[N]`…) | Modal do topo | Aciona opção correspondente diretamente |
| Modal fechado | Command bar | Volta para `ActionManager.Visible()` do contexto anterior |
| Segundo modal aberto (ex: help sobre confirmação) | Modal stack | Novo modal empilhado; apenas o topo recebe input |

---

## Seção 11 — Fluxos Visuais

Sequências de estados para casos de uso completos.

---

### Fluxo A — Abrir cofre

```
1. Welcome
   ─ command bar: F1 Abrir · F2 Criar · ? · ^Q Sair

2. → F1 → FilePicker modal (modo open)
   ─ lista arquivos .abditum do diretório
   ─ command bar: ↑↓ Navegar · Enter Abrir · Esc Cancelar

3. → Enter sobre cofre.abditum → FilePicker fecha
   → PasswordEntry modal abre
   ─ campo de senha com máscara
   ─ command bar: Enter Confirmar · Esc Cancelar

4. → Enter com senha → PasswordEntry fecha
   → MsgBusy "◐ Abrindo cofre..." aparece sobre welcome

5. → Sucesso → MsgBusy desaparece
   → Transição para Work Area: Cofre
   ─ Header: "Abditum · cofre.abditum" + aba Cofre ativa
   ─ Árvore expandida no painel esquerdo
   ─ Painel direito: placeholder "Selecione um segredo..."
```

---

### Fluxo B — Criar e salvar segredo

```
1. Work Area: Cofre (cofre aberto)
   → F21 → TextInput "Nome do segredo"

2. → Enter com nome → TextInput fecha
   → Select "Escolher pasta" (ou raiz se não houver pastas)

3. → Enter sobre pasta → Select fecha
   → Painel direito em modo edição com campos vazios
   ─ MsgHint "• Preencha os campos e pressione ^S para salvar"
   ─ Header: "•" aparece (dirty)

4. → Usuário preenche campos → ^S
   → MsgBusy "◐ Salvando cofre..."
   → MsgInfo "✓ Segredo criado" (TTL 2s)
   ─ Header: "•" desaparece

5. → Árvore atualiza: novo segredo aparece na pasta escolhida
```

---

### Fluxo C — Revelar senha e copiar

```
1. Work Area: Cofre, segredo selecionado
   → Painel direito mostra "Senha ••••••••••••  F16"
   → Foco no painel direito (Tab)

2. → F16 → campo "Senha" revela valor em texto claro
   ─ Timer de auto-ocultação inicia (15s default)
   ─ F16 ainda visível para re-ocultar manualmente

3. → F17 com foco no campo Senha
   → MsgInfo "✓ Senha copiada para a área de transferência" (TTL 3s)

4. → 15s passados → campo volta a ser mascarado automaticamente
   ─ Nenhuma mensagem — auto-ocultação é silenciosa
```

---

### Fluxo D — Bloqueio por inatividade

```
1. Usuário inativo por 75% do timeout configurado (ex: 5 min)
   → Tick: MsgWarn "⚠ Cofre será bloqueado em 90 segundos"
   ─ clearOnInput = true — desaparece ao próximo evento

2. Usuário interage → MsgWarn desaparece

3. Usuário fica inativo até 100% do timeout
   → Tick: bloqueio
   → Wipe de memória (vault, campos, clipboard)
   → Transição forçada para Work Area: Boas-vindas
   ─ Header perde nome do cofre e abas
   ─ MsgInfo "✓ Cofre bloqueado" (TTL 2s)

4. Para reabrir: F1 → FilePicker → PasswordEntry (Fluxo A a partir do passo 3)
```

---

### Fluxo E — Excluir segredo (com cancelamento)

```
1. Work Area: Cofre, segredo selecionado na árvore
   → F23 → DialogAlert "⚠ Excluir segredo"
   ─ "Gmail será excluído permanentemente. Esta ação não pode ser desfeita."
   ─ [S] Excluir (semantic.error bold)   [N] Cancelar

2a. → [N] ou Esc → DialogAlert fecha, nenhuma mudança

2b. → [S] → DialogAlert fecha
    → Segredo na árvore: "✕ ~~Gmail~~" em special.muted
    ─ Header: "•" aparece (dirty — exclusão está pendente em memória)
    ─ MsgHint "• Pressione ^S para salvar a exclusão definitivamente"

3. → ^S → MsgBusy → MsgInfo "✓ Cofre salvo"
   → Segredo desaparece definitivamente da árvore
   ─ Header: "•" desaparece
```
