# Design da TUI — Abditum

> Decisões visuais e de interação para o pacote `internal/tui`.  
> Complementa `tui-elm-architecture.md` (arquitetura) e `fluxos.md` (comportamento).

---

## Barra de Mensagens

### Espaço disponível

A barra de mensagens ocupa **1 linha** com aproximadamente **95% da largura do terminal**. É renderizada **sobre** um elemento existente da interface (provavelmente a última linha da work area), não como linha adicional — o layout do frame não reserva espaço dedicado.

Se o texto exceder a largura disponível, é truncado com reticências (`…`). Não há quebra de linha nem expansão vertical.

### Anatomia da mensagem

```
◐ Salvando cofre...
├─┤ └───────────────┘
 │        texto
emoji
```

Cada mensagem é composta por:
- **Emoji** — indica o tipo visualmente. Para `MsgBusy`, o emoji é animado (rotação a cada segundo).
- **Texto** — descritivo, conciso, em linguagem neutra (sem jargão técnico).

### Tipos de mensagem

| Tipo | Emoji | Cor do texto | Hex | Uso |
|---|---|---|---|---|
| `MsgInfo` | ✅ | Verde | `#9ece6a` | Operação concluída com sucesso |
| `MsgWarn` | ⚠️ | Amarelo | `#e0af68` | Atenção — bloqueio iminente, conflito detectado |
| `MsgError` | ❌ | Vermelho (bold) | `#f7768e` | Falha — salvamento, corrupção, operação impossível |
| `MsgBusy` | ◐ ◓ ◑ ◒ | Azul | `#7aa2f7` | Operação em andamento — spinner rotativo |
| `MsgHint` | 💡 | Cinza (itálico) | `#565f89` | Explicação contextual — descrição de campo, dica de uso |

A cor é aplicada ao **texto inteiro** (emoji + mensagem) via lipgloss. As cores acima pertencem à paleta Tokyo Night, consistente com o gradiente do ASCII art (`tui-elm-architecture.md` §D-14).

### Comportamento temporal

Cada mensagem tem dois parâmetros que controlam quando desaparece:

| Parâmetro | Efeito |
|---|---|
| `ttlSeconds` | Tempo em segundos até a mensagem expirar automaticamente. `0` = permanente (até ser substituída). |
| `clearOnInput` | Se `true`, a mensagem desaparece no próximo evento de teclado ou mouse. |

Combinações típicas:

| Situação | TTL | clearOnInput | Exemplo |
|---|---|---|---|
| Sucesso de operação | 2–3s | `false` | "Favoritado", "Copiado para a área de transferência" |
| Erro recuperável | 5s | `false` | "Falha ao salvar — arquivo em uso" |
| Progresso | 0 | `false` | "Salvando cofre..." — permanece até conclusão |
| Aviso de bloqueio | 0 | `true` | "Cofre será bloqueado em breve" — some ao interagir |
| Hint de campo | 0 | `false` | Descrição do campo em foco — substituído ao navegar |

### Prioridade e sobreposição

A barra exibe **uma única mensagem por vez**. Não há fila nem stack — a última chamada `Show()` sobrescreve a anterior.

Isso gera um comportamento natural quando combinado com o tick de 1 segundo:

```
Usuário inativo por 45s (75% do timeout)
    → tick: Show(MsgWarn, "Bloqueio iminente", 0, true)     ⚠️ Bloqueio iminente
Usuário copia um campo
    → HandleInput() limpa o warning (clearOnInput)
    → rootModel: Show(MsgInfo, "Copiado", 3, false)         ✅ Copiado
3 segundos depois
    → Tick() expira a mensagem                               (vazio)
1 segundo depois, se inatividade continua
    → tick: Show(MsgWarn, "Bloqueio iminente", 0, true)     ⚠️ Bloqueio iminente
```

O warning de bloqueio é **re-emitido naturalmente** pelo tick a cada segundo enquanto a condição persistir — sem mecanismo de prioridade explícito.

### Animação do spinner (`MsgBusy`)

O spinner avança **1 frame por segundo**, sincronizado com o tick global:

```
Segundo 0:  ◐ Salvando cofre...
Segundo 1:  ◓ Salvando cofre...
Segundo 2:  ◑ Salvando cofre...
Segundo 3:  ◒ Salvando cofre...
Segundo 4:  ◐ Salvando cofre...    (ciclo reinicia)
```

A frequência de 1fps é adequada para uma barra de texto. Não existe modal de progresso — `MsgBusy` é o único indicador de operação em andamento. O `activeFlow` já bloqueia input local (teclas caem no flow, que as ignora), tornando um overlay de progresso redundante.
