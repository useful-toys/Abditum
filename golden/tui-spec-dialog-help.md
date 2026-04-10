# Especificação Visual — Diálogo de Ajuda

> Diálogo funcional de referência de atalhos (somente leitura).
> Parte da [Especificação Visual](tui-specification.md).
>
> **Documentos relacionados:**
> - [`tui-design-system.md`](tui-design-system.md) — fundações visuais
> - [`tui-spec-dialogos.md`](tui-spec-dialogos.md) — anatomia comum e tipos de diálogo

## Help

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

### Identidade Visual

| Elemento | Token | Atributo |
|---|---|---|
| Título `Ajuda — Atalhos e Ações` | `text.primary` | **bold** |
| Label do grupo (`Navegação`, `Segredo`, `Cofre`) | `text.secondary` | **bold** |
| Tecla (ex: `Ctrl+R`, `Insert`, `^S`) | `accent.primary` | — |
| Descrição da ação | `text.primary` | — |
| Seta de scroll (`↑` / `↓` na borda direita) | `text.secondary` | — |
| Thumb de posição (`■` na borda direita) | `text.secondary` | — |
| Borda | `border.default` | — |

### Estados

| Componente | Estado | Condição |
|---|---|---|
| Conteúdo | sem scroll | Todas as ações cabem na viewport |
| Conteúdo | com scroll | Ações excedem a viewport — indicadores `↑`/`↓` e thumb `■` na borda direita (ver [DS — Scroll em diálogos](tui-design-system.md#scroll-em-diálogos)) |
| `F1` na barra de comandos | oculto (`HideFromBar`) | Enquanto o Help estiver aberto |
| Barra de comandos | vazia | Help não registra ações internas na barra |

### Eventos

| Evento | Efeito |
|---|---|
| `F1` pressionado (modal fechado) | Abre o modal; barra de comandos fica vazia; `F1` oculto |
| `F1` pressionado (modal aberto) | Fecha o modal; `F1` volta visível na barra |
| `Esc` | Fecha o modal; `F1` volta visível na barra |
| `↑` / `↓` | Scroll por linha (se conteúdo excede viewport) |
| `PgUp` / `PgDn` | Scroll por página (viewport − 1 linhas) |
| `Home` / `End` | Vai ao início / fim do conteúdo |

### Comportamento

- **Conteúdo dinâmico** — gerado a partir de todas as ações registradas no ActionManager no momento da abertura
- **Agrupamento** — ações são organizadas pelo atributo numérico `Grupo`. Cada grupo tem um `Label` registrado (ex: 1 → "Navegação", 2 → "Segredo"). Grupos renderizados em ordem numérica crescente
- **Ordenação interna** — dentro de cada grupo, ações ordenadas por `Prioridade` (maior primeiro)
- **Scroll** — segue o padrão transversal do DS: indicadores `↑`/`↓` na borda direita, navegação por `↑↓` / `PgUp`/`PgDn` / `Home`/`End`
- **Borda inferior** — `Esc Fechar` sempre visível, independente do estado de scroll

---
