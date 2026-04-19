## Why

A tela de Configurações (`WorkAreaSettings`) existe na arquitetura e no cabeçalho da TUI, mas hoje renderiza apenas o texto literal `"Settings"`. O usuário não consegue visualizar nem alterar nenhuma preferência pelo modo Config, tornando a aba decorativa. Esta mudança materializa a tela para que configurações reais — tema visual, timers de segurança e informações do ambiente — sejam acessíveis e editáveis diretamente na área de trabalho dedicada.

## What Changes

- Substituir o placeholder `"Settings"` por uma tela de configurações estruturada em grupos, com navegação por teclado, indicação de valor atual e descrição contextual por item.
- Implementar os grupos iniciais da tela:
  - **Aparência** — exibição do tema ativo (Tokyo Night / Cyberpunk) como campo focável, refletindo imediatamente a escolha quando o tema for alterado via `F12`.
  - **Segurança** — temporizadores configuráveis em segundos: bloqueio automático por inatividade, ocultação de campo sensível e limpeza de clipboard.
  - **Sobre** — informações somente leitura sobre a aplicação e o cofre atual.
- Documentar a seção "Modo Configurações" em `golden/tui-spec-telas.md` com wireframe, identidade visual, estados, mensagens, eventos e comportamento, alinhados ao design system.
- Persistir tema e timers como configurações do cofre, com mutação síncrona aplicada pela própria `SettingsView` via `vault.Manager` e gravação no payload criptografado no próximo salvamento do cofre.
- Exibir hints contextuais na barra de mensagens para orientar foco e edição de campos numéricos, com ajustes rápidos de `+/-` em passos de 5 segundos, sem depender de ajuda externa.
- Validar visualmente a tela com testes golden file, cobrindo ao menos os estados estruturais mais importantes do modo Configurações.

## Capabilities

### New Capabilities
- `settings-screen`: Tela de configurações da TUI com navegação por grupos, edição de timers de segurança, item focável de tema e documentação visual no golden spec.

### Modified Capabilities

- Nenhuma.

## Impact

- Código afetado: `internal/tui/settings/settings_view.go` (implementação principal), `internal/tui/root.go` (integração com tema global e montagem da tela), barra de mensagens da TUI, modelo de cofre (`Configuracoes`) e persistência do tema no domínio.
- Documentação afetada: `golden/tui-spec-telas.md` — seção "Modo Configurações" deixa de ser placeholder.
- Validação afetada: testes da TUI passam a incluir golden files para a renderização da tela de settings.
- UX afetada: navegação dentro da área de trabalho Config, observação do tema ativo com troca global via `F12`, edição de timers de segurança e hints contextuais por item.
