# Especificação Visual — Telas

> Telas completas e fluxos visuais.
> Parte da [Especificação Visual](tui-specification.md).
>
> **Documento de fundação:**
> - [`tui-design-system.md`](tui-design-system.md) — fundações visuais

## Boas-vindas

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

> Logo e versão centralizados via `lipgloss.Place()`. As linhas do logo recebem as cores do [DS — Gradiente do logo](tui-design-system.md#gradiente-do-logo) — não representável neste wireframe monocromático.

### Identidade Visual

| Elemento | Token | Atributo |
|---|---|---|
| Logo (linhas 1–5) | DS — [Gradiente do logo](tui-design-system.md#gradiente-do-logo) — por linha | — |
| Versão (ex: `v0.1.0`) | `text.secondary` | — |

> As cores do logo não são tokens nomeados — são os valores hexadecimais da tabela de gradiente do DS, aplicados por linha conforme o tema ativo.

### Estados

| Componente | Estado | Condição |
|---|---|---|
| Logo + versão | visível, centralizado | Tela ativa |
| Cabeçalho | sem abas | Nenhum cofre aberto — ver [Cabeçalho — Sem cofre](tui-spec-cabecalho.md#sem-cofre-boas-vindas) |

### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Tela entra em exibição | Dica de uso | `• Abra ou crie um cofre para começar` |

### Eventos

| Evento | Efeito |
|---|---|
| Aplicação inicia sem cofre | Modo boas-vindas exibido |
| Cofre fechado | Tela boas-vindas exibida |
| Cofre bloqueado | Tela boas-vindas exibida; arquivo permanece em disco, requer nova autenticação |
| Terminal redimensionado | Logo e versão recentralizados |

### Comportamento

- Logo e versão centralizados horizontal e verticalmente na área de trabalho via `lipgloss.Place()`
- As cores do logo acompanham o tema ativo — mudam instantaneamente com `F12`
- O cabeçalho não exibe abas neste modo (ver [Cabeçalho — Sem cofre](tui-spec-cabecalho.md#sem-cofre-boas-vindas))
- **Versão dinâmica** — o texto exibido vem da string injetada em tempo de build via `-ldflags "-X main.version=$(git describe --tags --always)"`. Em builds locais sem tag, exibe `dev`. O valor **nunca** é hardcoded no fonte

## Modo Configurações

**Trigger:** Usuário pressiona `F4` (ou seleciona a aba `Config` no cabeçalho).
**Interação:** Navegação por teclado (↑↓), edição inline de campos numéricos, ajuste rápido com `+`/`-`.

**Wireframe — Estado 1: item de tema com foco (80 × 24):**

```
 Abditum  meu-cofre.abd  •                                   Cofre  Modelos  Config
 ──────────────────────────────────────────────────────────────────────────────────


                             Configurações

                             Aparência
                           › Tema visual                     Tokyo Night
                             Tema aplicado ao cofre atual.

                             Segurança
                             Bloqueio por inatividade        300 s
                             Ocultar campo sensível          15 s
                             Limpar área de transferência    30 s

                             Sobre
                             Versão                          v0.1.0
                             Arquivo do cofre                meu-cofre.abd




 ─── • F12 para alternar tema visual ───────────────────────────────────────────
 F1 Ajuda · F2 Cofre · F3 Modelos · F4 Config · F7 Salvar · Ctrl+Q Sair
```

**Wireframe — Estado 2: campo numérico em modo de edição (80 × 24):**

```
 Abditum  meu-cofre.abd  •                                   Cofre  Modelos  Config
 ──────────────────────────────────────────────────────────────────────────────────


                             Configurações

                             Aparência
                             Tema visual                     Tokyo Night

                             Segurança
                           › Bloqueio por inatividade       300  s
                             Tempo de bloqueio automático por inatividade.

                             Ocultar campo sensível          15 s
                             Limpar área de transferência    30 s

                             Sobre
                             Versão                          v0.1.0
                             Arquivo do cofre                meu-cofre.abd


 ─── • Enter confirma · Esc cancela ────────────────────────────────────────────
 F1 Ajuda · F2 Cofre · F3 Modelos · F4 Config · F7 Salvar · Ctrl+Q Sair
```

> `•` = indicador de cofre modificado (`semantic.warning`). `›` = item com foco (`special.highlight` + bold). O campo em edição (Estado 2) usa fundo `surface.input` — em wireframe monocromático não é distinguível por cor; na implementação é o único delimitador visual do campo. O cursor **real** do terminal fica posicionado ao final do buffer numérico; não há caractere `▌` artificial. A linha de descrição aparece apenas sob o item com foco; nos demais itens não há linha de descrição. Padding vertical simétrico (linhas em branco) centraliza o bloco na área útil.

### Identidade Visual

| Elemento | Token | Atributo |
|---|---|---|
| Título `Configurações` | `text.primary` | — |
| Cabeçalho de grupo (ex: `Aparência`) | `text.primary` | bold |
| Label de item sem foco | `text.primary` | — |
| Label de item com foco | `accent.primary` + `special.highlight` (fundo) | bold |
| Símbolo de foco `›` | `accent.primary` | — |
| Valor de item sem foco | `text.secondary` | — |
| Valor de item numérico com foco (navegação) | `accent.primary` | — |
| Campo numérico em edição | `text.primary` / `surface.input` (fundo) | — |
| Unidade `s` (fora do campo de entrada) | `text.primary` | — |
| Linha de descrição inline | `text.secondary` | — |

### Estados

| Componente | Estado | Condição |
|---|---|---|
| Item de tema | focado, não editável | Cursor sobre o item; `F12` é o mecanismo de troca |
| Item numérico | focado, navegação | Cursor sobre o item; `+`/`-` e `Enter` disponíveis |
| Item numérico | em edição | Após `Enter` sobre o campo; apenas dígitos e `Backspace` aceitos |
| Item somente leitura | focado, sem ação | Cursor sobre `Versão` ou `Arquivo do cofre`; `Enter` ignorado |
| Qualquer item | sem foco | Sem símbolo `›`; valor em `text.secondary` |
| Tela inteira | sem cofre aberto | Campos numéricos exibem `0 s`; arquivo exibe `–`; itens focáveis mas sem efeito ao aplicar |

### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Item de tema com foco | Dica de campo | `• F12 para alternar tema visual` |
| Item numérico com foco (navegação) | Dica de campo | `• Enter edita · +/- altera o valor` |
| Item numérico em edição | Dica de campo | `• Enter confirma · Esc cancela` |
| Valor fora do range mínimo ao confirmar | Erro | `✕ Mínimo: X s` |
| Item somente leitura com foco | — | (sem hint) |

### Eventos

| Evento | Efeito |
|---|---|
| ↑ / ↓ | Move o foco para o item anterior/próximo com wrapping cíclico |
| `Enter` sobre item numérico | Entra em modo de edição inline |
| `Enter` sobre item de tema ou somente leitura | Ignorado |
| Dígito (`0`–`9`) em edição | Acrescenta ao buffer |
| `Backspace` em edição | Remove último dígito do buffer |
| `Enter` em edição (valor válido) | Confirma, persiste via `vault.Manager`, sai do modo de edição |
| `Enter` em edição (valor inválido) | Exibe erro, permanece em edição |
| `Esc` em edição | Restaura valor original, sai do modo de edição sem salvar |
| `+` / `-` em navegação sobre numérico | Ajusta valor em ±5 s; aplica e persiste imediatamente |
| `F12` (global) | Alterna tema visual; refletido imediatamente no item de tema da tela |
| Abertura de cofre | Valores de configuração carregados do cofre |
| Salvamento de cofre | Configurações persistidas no arquivo criptografado |

### Comportamento

- Navegação é cíclica: ↑ no primeiro item vai para o último; ↓ no último vai para o primeiro.
- Cabeçalhos de grupo (`Aparência`, `Segurança`, `Sobre`) não recebem foco e são pulados na navegação.
- Linha de descrição contextual aparece **imediatamente abaixo** do item com foco, com indentação equivalente à dos itens.
- Campos numéricos em edição aceitam apenas dígitos e `Backspace`; `+`, `-`, ↑ e ↓ são silenciosamente ignorados durante a edição.
- Confirmação de edição valida o valor localmente antes de chamar `vault.Manager.AlterarConfiguracoes`; erro de domínio é exibido via barra de mensagens sem sair do modo de edição.
- Cancelamento com `Esc` restaura o valor exibido para o estado anterior; nenhuma chamada ao domínio é feita.
- Ajuste rápido com `+`/`-` aplica o delta em passos de 5 s; aplica imediatamente sem entrar em modo de edição.
- Quando nenhum cofre está aberto, os campos numéricos exibem `0 s` e o arquivo exibe `–`; tentativas de editar não disparam `vault.Manager`.
- O item de tema mostra o nome do tema ativo (sincronizado a cada frame via `syncTema`); mudanças via `F12` são refletidas imediatamente sem necessidade de evento explícito.
- O tema alterado via `F12` sem cofre aberto não é persistido; ao abrir um cofre, o tema salvo no arquivo sobrepõe o estado em memória.
- A tela não implementa fallback próprio para terminais abaixo de 24 linhas — depende do guard centralizado no `RootModel`.
- Mudanças de configuração não emitem mensagem global (`tea.Msg`); consumidores externos devem reler o estado canônico em `vault.Manager` quando necessário.
- **Versão dinâmica** — o valor exibido em `Versão` vem da string injetada em build via `-ldflags`; em builds locais sem tag exibe `dev`.

<!-- SEÇÕES FUTURAS — a preencher pela equipe -->

<!--
## Telas (continuação)

### Modo Cofre
### Modo Modelos

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
