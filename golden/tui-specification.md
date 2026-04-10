# Especificação Visual — Abditum TUI

> Wireframes, layouts de componentes e fluxos visuais concretos.
> Cada tela e componente consome os padrões definidos no design system.
>
> **Documento de fundação:**
> - [`tui-design-system.md`](tui-design-system.md) — princípios, tokens, estados, padrões transversais

## Documentos

| Documento | Conteúdo | Linhas |
|---|---|---|
| [Diálogos](tui-spec-dialogos.md) | Anatomia comum, tipos (Notificação, Confirmação, Ajuda, Funcional), severidade, catálogo de decisão | ~420 |
| [Diálogos de Senha](tui-spec-dialog-senha.md) | PasswordEntry, PasswordCreate | ~230 |
| [FilePicker](tui-spec-dialog-filepicker.md) | Diálogo de seleção de arquivo (modos Open e Save) | ~300 |
| [Ajuda](tui-spec-dialog-help.md) | Diálogo de referência de atalhos | ~100 |
| [Cabeçalho](tui-spec-cabecalho.md) | Título, abas de modo, indicador dirty, busca na linha separadora | ~240 |
| [Barras](tui-spec-barras.md) | Barra de Comandos (ações, prioridade, truncamento) e Barra de Mensagens (severidade, ciclo de vida) | ~230 |
| [Árvore de Segredos](tui-spec-arvore.md) | Painel esquerdo, busca, navegação e ações na árvore | ~610 |
| [Detalhe do Segredo](tui-spec-detalhe.md) | Painel direito: modos Leitura, Edição de Valores e Edição de Estrutura | ~890 |
| [Telas](tui-spec-telas.md) | Boas-vindas e stubs de telas futuras | ~100 |

## Atalhos da Aplicação

Esta seção atribui teclas a fluxos concretos. As políticas transversais de teclado — notação, convenções semânticas, escopos e atalhos globais (`F1`, `F12`, `⌃Q`, `⌃!⇧Q`) — são definidos no [Design System — Teclado](tui-design-system.md#teclado).

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
