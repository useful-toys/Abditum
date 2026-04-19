## Context

A TUI do Abditum já reserva `WorkAreaSettings` na máquina de estados do `RootModel`, exibe a aba `Config` no cabeçalho e instancia `settings.NewSettingsView`. O código atual devolve apenas a string literal `"Settings"` sem qualquer estrutura visual ou interação real.

Ao mesmo tempo, `golden/tui-spec-telas.md` tem a seção "Modo Configurações" comentada como placeholder futuro, e `golden/requisitos.md` especifica preferências configuráveis: tema visual (persistido no cofre), bloqueio automático por inatividade, ocultação de campo sensível e limpeza de clipboard.

O design system já define dois temas (`TokyoNight` e `Cyberpunk`), tokens semânticos para estados on/off (`semantic.success`, `semantic.off`) e o atalho global `F12` para alternância rápida de tema. O `RootModel` já possui `ToggleTheme()`. A tela precisa expor esses estados navegáveis sem inventar uma linguagem visual paralela.

## Goals / Non-Goals

**Goals:**
- Substituir o placeholder por uma tela estruturada em grupos com navegação por teclado (↑↓), seleção visual e edição inline de preferências.
- Cobrir o escopo de primeira entrega: grupo Aparência (tema), grupo Segurança (timers), grupo Sobre (somente leitura).
- Documentar wireframe, estados, mensagens e comportamento em `golden/tui-spec-telas.md`.
- Persistir o tema no payload do cofre (já previsto nos requisitos).
- Reagir imediatamente a mudanças de tema, inclusive quando disparadas por `F12` fora da tela.
- Deixar explícito o contrato de aplicação: para mudanças síncronas de settings, a própria `SettingsView` valida a entrada local, chama `vault.Manager` diretamente e o `Cofre` permanece como fonte de verdade persistível.

**Non-Goals:**
- Criar um framework genérico de formulários para toda a TUI antes de validar a tela de settings.
- Adicionar preferências além das já suportadas pelo modelo de domínio e código atual.
- Redesenhar cabeçalho, barra de mensagens ou barra de comandos de forma global.
- Implementar importação/exportação de preferências.

## Decisions

1. **Layout de coluna única centralizada, com título de tela e conteúdo verticalmente centrado.**
   O conteúdo ocupa uma faixa central da área de trabalho, com padding vertical simétrico de linhas em branco acima e abaixo. O título `Configurações` aparece antes do primeiro grupo. Grupos são separados por uma linha em branco; não há linhas `─────` entre o cabeçalho do grupo e seus itens — o próprio espaçamento e o peso tipográfico (bold no grupo) estabelecem a hierarquia. Essa abordagem é mais discreta, condizente com o princípio de "discrição e portabilidade" do DS, e remove ruído visual sem prejudicar a legibilidade.
   **Alternativa considerada:** painel esquerdo alinhado com separadores `─────` por grupo (proposta inicial). Descartada pelo próprio autor — mais pesada visualmente, desnecessária para uma lista de preferências simples.

2. **Estrutura interna `settingItem` com cursor de navegação por índice.**
   Cada item é uma struct com chave, valor atual, descrição e metadados de interação. O cursor aponta para o índice do item selecionado. ↑↓ movem o cursor; Enter ativa edição inline apenas nos itens numéricos; Esc cancela edição.
   **Alternativa considerada:** modelo de foco por componentes separados. Descartada — sobrecarga de indireção desnecessária para uma lista linear.

3. **Campo de tema apenas focável, sem edição própria na tela.**
   O item de tema existe para tornar o estado atual visível dentro de Configurações e para contextualizar o atalho global `F12`, que permanece sendo o mecanismo padrão de troca. Quando o foco está sobre o item, a barra de mensagens orienta o uso de `F12`.
   **Alternativa considerada:** toggle cíclico inline por Enter ou modal de seleção. Descartada — duplicaria um comportamento global já consolidado e adicionaria um modelo de interação extra sem benefício claro.

4. **Dois modos distintos para campos numéricos: navegação e edição. Todos os timers em segundos, passo de 5.**
   Todos os temporizadores são armazenados e exibidos em **segundos** — inclusive o de bloqueio (que era "minutos" no campo de domínio `tempoBloqueioInatividadeMinutos`, que precisará ser renomeado para `tempoBloqueioInatividadeSegundos`). O passo de `+`/`-` é **5 segundos** por tecla. Quando um campo numérico está **focado mas não em edição**: `+` incrementa 5, `-` decrementa 5, `Enter` entra em modo de edição. Quando está **em edição**: apenas dígitos (`0–9`) e `Backspace` são aceitos; o cursor real do terminal fica posicionado ao final do buffer; `Enter` confirma e aplica a validação de range; `Esc` restaura o valor original e sai sem salvar. A unidade `s` permanece estática à direita do campo. Erros de range são exibidos via barra de mensagens (`✕` erro), não inline.
   **Alternativa considerada:** modo exclusivo de `+/-` sem entrada de texto. Descartada — lento para valores maiores; entrada numérica direta é mais eficiente.

5. **Persistir settings via `vault.Manager`, chamados diretamente pela `SettingsView` quando a operação for síncrona.**
    O requisito especifica que o identificador do tema deve ser gravado no payload criptografado do cofre, e o modelo de domínio trata `Configuracoes` como parte do `Cofre`. A `SettingsView` exibe o estado atual, valida a entrada local, chama `vault.Manager` nas confirmações de mudança e o cofre fica marcado como modificado até o próximo salvamento.
    **Alternativa considerada:** arquivo de configuração externo. Descartada — viola o princípio de portabilidade ("nenhum arquivo fora do cofre").
    
    Para timers e demais mutações rápidas em memória, a tela não deve emitir `tea.Msg` de alteração de configuração nem depender de ida e volta via `RootModel`. Ao confirmar um valor válido, a própria `SettingsView` chama `vault.Manager.AlterarConfiguracoes(...)`, trata o `error` localmente e atualiza sua renderização/hints. O `RootModel` permanece responsável apenas por layout global, troca de tema via `F12` e roteamento da tela ativa.

6. **Hints contextuais na barra de mensagens por tipo de campo e estado.**
   Cada item tem um hint exibido na barra de mensagens (`•` dica) ao receber foco, que muda com o estado:
   - Campo de tema (focado): `• F12 para alternar tema visual`
   - Campo numérico (focado, não em edição): `• Enter edita · +/- altera o valor`
   - Campo numérico (em edição): `• Enter confirma · Esc cancela`
   A linha de descrição do item (inline, abaixo do item) é estática por campo e descreve o que a configuração faz — não instrui sobre teclas. Os dois canais são complementares: descrição = semântica, hint = instrução de uso imediato.
   **Alternativa considerada:** linha de descrição fixa na última linha da work area. Descartada — distância visual entre item e descrição reduz coesão; inline é mais legível.

7. **Cobrir a tela com testes golden file além de testes comportamentais.**
   Como a mudança introduz uma work area visual nova, a validação precisa proteger o layout renderizado contra regressões acidentais de espaçamento, agrupamento e estados visuais relevantes. Golden tests complementam os testes de interação e ajudam a manter a implementação alinhada ao wireframe documentado.
   **Alternativa considerada:** validar apenas por asserts textuais parciais. Descartada — útil para comportamento, mas fraca para regressões de composição visual da tela.

8. **A tela de settings não implementa fallback próprio para terminais abaixo do mínimo.**
   O `RootModel` já bloqueia a renderização normal quando a altura é menor que `design.MinHeight` (24 linhas). A `SettingsView` pode assumir que só será renderizada quando a área útil estiver válida e não deve duplicar uma segunda política local para tamanhos mínimos.
   **Alternativa considerada:** fallback específico da própria tela para alturas pequenas. Descartada — duplicaria responsabilidade já centralizada no root.

9. **Mudanças de settings não inauguram um mecanismo genérico de broadcast.**
   Para esta feature, não haverá barramento de eventos nem mensagem de aplicação por campo alterado. Quem precisar de configuração no futuro deve reler o estado canônico no domínio (`vault.Manager` / `Configuracoes`) quando fizer sentido. Só vale introduzir notificação explícita se surgir um consumidor real com estado derivado persistente.
   **Alternativa considerada:** emitir mensagem global a cada alteração bem-sucedida. Descartada — complexidade prematura e acoplamento sem benefício imediato.

## Risks / Trade-offs

- **[Persistência de tema acoplada à abertura do cofre]** → Quando não há cofre aberto, o tema alterado em settings não é persistido. Documentar isso no comportamento da tela; ao abrir o cofre, o tema salvo no arquivo sobrepõe o estado em memória.
- **[Ranges divergirem entre UI e domínio]** → Fixar no pacote `settings` e no domínio os mínimos normativos já definidos em `golden/requisitos.md`: bloqueio `> 60 s`, ocultação `> 2 s`, clipboard `> 10 s`, todos com passo de ajuste de `5 s`.
- **[Contrato de aplicação ficar ambíguo]** → Formalizar na spec que mudanças síncronas da tela de settings são aplicadas localmente pela `SettingsView` via `vault.Manager`, sem mensagens de aplicação para cada campo.
- **[Responsabilidade de tamanho mínimo duplicada]** → Registrar que a tela depende do guard já existente no `RootModel` para `height < 24`, sem fallback próprio.
- **[Tela de settings acessível sem cofre aberto]** → A aba `Config` aparece no cabeçalho independentemente do estado do cofre? Verificar o comportamento atual do `RootModel` e manter consistência.
- **[Regressões visuais passarem despercebidas]** → Exigir golden tests para estados estruturais da tela de settings, além de testes comportamentais.

## Migration Plan

1. Documentar o wireframe e o comportamento do modo Configurações em `golden/tui-spec-telas.md`.
2. Implementar `settingItem`, cursor e renderização em `internal/tui/settings/settings_view.go`.
3. Adicionar ações de teclado (↑↓, Enter, Esc, `+`, `-`) registradas no `ActionManager` para o escopo de settings, preservando `F12` como atalho global de tema.
4. Conectar a persistência de settings via mutação de cofre diretamente da `SettingsView` para o `vault.Manager`.
5. Integrar hints contextuais de foco/edição à barra de mensagens.
6. Adicionar testes de renderização, navegação, hints e aplicação direta de mudanças via `vault.Manager`.
7. Adicionar golden tests para os principais estados visuais da tela de settings.

## Open Questions

- A aba `Config` deve aparecer no cabeçalho mesmo sem cofre aberto? O comportamento atual precisa ser verificado — se não aparecer, a tela de settings só é acessível com cofre aberto, o que limita onde documentar esse estado.

## Wireframes

Os wireframes abaixo representam o layout acordado em 80 colunas × 24 linhas.
O formato exato do cabeçalho segue `golden/tui-spec-cabecalho.md` — os exemplos abaixo usam a convenção ilustrativa `[Config]` para simplificar; na implementação real a aba ativa usa `╯ Config ╰`.

### Estado 1 — Item de tema com foco

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

### Estado 2 — Campo numérico em modo de edição

```
 Abditum  meu-cofre.abd  •                                   Cofre  Modelos  Config
 ──────────────────────────────────────────────────────────────────────────────────


                             Configurações

                             Aparência
                             Tema visual                     Tokyo Night

                             Segurança
                           › Bloqueio por inatividade       [300▌] s
                             Tempo de bloqueio automático por inatividade.

                             Ocultar campo sensível          15 s
                             Limpar área de transferência    30 s

                             Sobre
                             Versão                          v0.1.0
                             Arquivo do cofre                meu-cofre.abd


 ─── • Enter confirma · Esc cancela ────────────────────────────────────────────
 F1 Ajuda · F2 Cofre · F3 Modelos · F4 Config · F7 Salvar · Ctrl+Q Sair
```

**Notas sobre os wireframes:**
- `•` = indicador de cofre modificado (`semantic.warning`); aparece quando há mudanças não salvas.
- `›` = item com foco (`special.highlight` + bold); itens sem foco têm indentação equivalente sem símbolo.
- `[300▌] s` = campo em edição: colchetes delimitam a área `surface.input`; `▌` representa a posição do cursor real do terminal no wireframe (convenção ASCII art — na implementação, usa-se o cursor real do terminal, não o caractere `▌`); `s` fica fora do campo.
- A linha de descrição aparece apenas sob o item focado; nos demais itens não há linha de descrição.
- Padding vertical simétrico (linhas em branco acima e abaixo do conteúdo) centraliza o bloco na área útil.

