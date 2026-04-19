## 1. Especificação visual (golden)

- [ ] 1.1 Adicionar a seção "Modo Configurações" em `golden/tui-spec-telas.md` com wireframe ASCII 80×24 mostrando grupos Aparência, Segurança e Sobre, item selecionado, valor e linha de descrição.
- [ ] 1.2 Preencher a tabela de Identidade Visual do modo Configurações (tokens por elemento).
- [ ] 1.3 Preencher a tabela de Estados (item focado, item em edição, item somente leitura, grupo sem cofre aberto).
- [ ] 1.4 Preencher a tabela de Mensagens (hint de campo focado, hint de edição, erro de valor fora do range, confirmação de mudança aplicada).
- [ ] 1.5 Preencher a tabela de Eventos (navegação, ativação, edição, ajuste rápido com `+/-`, mudança de tema via F12).
- [ ] 1.6 Preencher a lista de Comportamento (navegação cíclica, edição inline numérica, aplicação imediata de tema, persistência no cofre, somente leitura no grupo Sobre).

## 2. Modelo de dados e estrutura interna

- [x] 2.1 Definir `settingItem` (chave, label, valor atual, tipo: toggle/numeric/readonly, descrição) em `internal/tui/settings/settings_view.go`.
- [x] 2.2 Definir constantes de range e passo dos timers em segundos no mesmo pacote (`minAutoLockSeconds = 61`, `minHideSeconds = 3`, `minClipboardSeconds = 11`, `timerStepSeconds = 5`).
- [x] 2.3 Construir a lista inicial de itens no `NewSettingsView`, separados por grupo, usando os valores padrão dos requisitos.
- [x] 2.4 Adicionar campo `cursor int` e campo `editMode bool` (com buffer de edição) na `SettingsView`.
- [x] 2.5 Adicionar getters públicos em `vault.Configuracoes`: `TempoBloqueioSegundos() int`, `TempoOcultarSegundos() int`, `TempoLimparTransferenciaSegundos() int` — necessários para que `SettingsView` leia os valores atuais fora do pacote `vault`.
- [x] 2.6 Adicionar função `vault.NovasConfiguracoes(bloqueio, ocultar, limpar int) Configuracoes` para construir o struct a ser passado ao `AlterarConfiguracoes` fora do pacote.
- [x] 2.7 Adicionar campo `temaVisual string` em `vault.Configuracoes` e atualizar serialização
- [x] 2.8 Atualizar `NewSettingsView` para receber `tui.MessageController` além de `*vault.Manager`, e atualizar a chamada em `RootModel.initVaultViews` para passar `r.MessageController()`.

## 3. Renderização

- [x] 3.1 Implementar a renderização de cabeçalho de grupo (bold, sem separador `─────`) e padding de linha em branco entre grupos, usando tokens do tema ativo.
- [x] 3.2 Implementar a renderização de item normal (label + valor com espaçamento alinhado, `text.primary`/`text.secondary`).
- [x] 3.3 Implementar o estado de item selecionado (highlight `accent.primary` no label, símbolo `›` ou equivalente do DS).
- [x] 3.4 Implementar o estado de item em edição (campo inline editável com cursor visível, fundo `surface.input`).
- [x] 3.5 Implementar o estado de item somente leitura (label + valor sem destaque interativo, `text.disabled` no valor se aplicável).
- [x] 3.6 Implementar o título `Configurações` e o centramento vertical do conteúdo (padding simétrico de linhas em branco).
- [x] 3.7 Implementar a linha de descrição contextual inline, imediatamente abaixo do item com foco (mesmo alinhamento de indentação).
- [x] 3.8 Garantir que os campos numéricos exponham apenas a parte numérica como editável, com unidade fixa `s` fora do input.
- [x] 3.9 Implementar o item `Arquivo do cofre` no grupo Sobre (nome do arquivo ativo, somente leitura).
- [x] 3.10 Atualizar `Render(height, width int, theme *design.Theme)` para compor toda a tela a partir dos elementos acima.
- [x] 3.11 Manter a tela sem fallback local para `height < 24`, confiando no guard central já existente no `RootModel`.

## 4. Navegação e interação

- [x] 4.1 Implementar movimento de cursor ↑↓ com wrapping em `HandleKey`.
- [x] 4.2 Garantir que o item de tema seja apenas focável, não editável, preservando `F12` como mecanismo de troca global.
- [x] 4.3 Implementar entrada em modo de edição numérica (Enter → edit mode) em `HandleKey`.
- [x] 4.4 Implementar edição de campo numérico (somente dígitos, Backspace, navegação básica) no modo de edição em `HandleKey`.
- [x] 4.5 Implementar confirmação de edição numérica (Enter → validação de range → aplicar ou exibir erro) em `HandleKey`.
- [x] 4.6 Implementar cancelamento de edição (Esc → restaurar valor original e sair do edit mode) em `HandleKey`.
- [x] 4.7 Implementar ajuste rápido com `+/-` para campos numéricos focados fora do modo de edição, sempre em passos de 5 segundos.
- [x] 4.8 Atualizar hints contextuais da barra de mensagens ao focar e editar cada item de settings.
- [x] 4.9 Registrar as ações de teclado da tela de settings no `ActionManager` via `Actions()`.

## 5. Integração com tema e cofre

- [x] 5.1 Garantir que a alteração de tema via `F12` continue sendo coordenada no `RootModel` e refletida imediatamente na tela de settings.
- [x] 5.2 Garantir que mudança de tema via `F12` atualiza o valor exibido no item de tema da tela de settings.
- [x] 5.3 Conectar timers e tema ao `vault.Manager` / `Configuracoes`, marcando o cofre como modificado quando houver mudança aplicada, sem criar mensagens globais de propagação por campo.
- [x] 5.6 Garantir que a implementação não introduza broadcast genérico ou `tea.Msg` global por alteração de campo; consumidores futuros devem reler a configuração canônica no domínio se necessário.
- [x] 5.4 Renomear `tempoBloqueioInatividadeMinutos` para `tempoBloqueioInatividadeSegundos`
- [x] 5.5 Alinhar a serialização do cofre para incluir o campo `temaVisual` de `Configuracoes`

## 6. Testes

- [ ] 6.1 Adicionar golden tests da tela de settings para o estado estruturado padrão e para pelo menos um estado de edição numérica.
- [ ] 6.2 Adicionar teste de navegação: cursor desce e sobe com wrapping correto.
- [ ] 6.3 Adicionar teste do item de tema: recebe foco, mostra hint contextual e reflete mudança disparada por `F12`.
- [ ] 6.4 Adicionar teste de edição numérica: somente dígitos são aceitos; valor válido é aplicado; valor fora do range é rejeitado.
- [ ] 6.5 Adicionar teste de hints e feedback na barra de mensagens para foco, edição e rejeição/aceite local.
- [ ] 6.6 Adicionar teste de aplicação direta: mudança aplicada em settings chega ao domínio do cofre via `vault.Manager`, sem passar por mensagem de orquestração no `RootModel`.
- [ ] 6.7 Manter testes comportamentais complementares (navegação, aplicação direta, validação) além dos golden tests.
- [ ] 6.8 Executar todos os testes Go relevantes (`internal/tui`, `internal/vault`) e corrigir regressões introduzidas.
