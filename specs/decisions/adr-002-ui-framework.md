# ADR 002 — Framework de UI

**Status**: Aceita
**Data**: 2026-03-02

## Contexto

A aplicação precisa de uma interface multiplataforma em Go que suporte: árvore navegável, formulários dinâmicos, campos mascaráveis, atalhos de teclado e integração com área de transferência.

Duas categorias de UI estão em consideração: **GUI** (janela gráfica) e **TUI** (interface de terminal).

## Opções

### GUI

| Framework    | Abordagem                     | CGO | Win | Mac | Linux | Maturidade |
|--------------|-------------------------------|-----|-----|-----|-------|------------|
| **Fyne**     | Widgets próprios em OpenGL    | Sim | ✓   | ✓   | ✓     | Alta       |
| **Wails**    | WebView (HTML/CSS/JS frontend) | Sim | ✓   | ✓   | ✓     | Alta       |
| **Gio**      | Immediate mode, sem CGO       | Não | ✓   | ✓   | ✓     | Média      |

### TUI

| Framework      | Abordagem                              | CGO | Maturidade |
|----------------|----------------------------------------|-----|------------|
| **Bubble Tea** | Elm architecture, componível           | Não | Alta       |
| **tview**      | Widgets prontos (Tree, Form, Table...) | Não | Alta       |
| **Lip Gloss**  | Estilização; usado com Bubble Tea      | Não | Alta       |

## Análise

### Fyne
- API Go pura; widgets incluem `Tree`, `Form`, `Entry`, `PasswordEntry`
- Aparência consistente mas não-nativa (visual próprio)
- Boa documentação; usado em produção
- CGO necessário

### Wails
- Frontend em HTML/CSS/JS (React, Vue, Svelte, etc.) — UI ilimitada
- Backend em Go puro com bridge automático
- Requer decisão adicional sobre framework frontend
- Mais complexidade de setup

### Gio
- Sem CGO; mais portável
- Immediate mode exige mais código para UIs complexas
- Menor ecossistema de widgets prontos

### Bubble Tea + Lip Gloss
- Arquitetura Elm (Model/Update/View) — previsível e testável
- Sem CGO; binário 100% estático
- Exige composição manual de componentes (sem widget `Tree` pronto)
- Excelente para poder de usuários que vivem no terminal
- Referência: `charm.sh` — ecossistema ativo e bem mantido

### tview
- Widgets prontos: `TreeView`, `Form`, `InputField`, `TextView`
- Mais próximo do paradigma GUI em termos de componentes disponíveis
- Menos flexível visualmente que Bubble Tea + Lip Gloss
- Sem CGO

## Comparação geral

| Critério                    | Fyne | Wails | Gio | Bubble Tea | tview |
|-----------------------------|------|-------|-----|------------|-------|
| CGO necessário              | Sim  | Sim   | Não | Não        | Não   |
| Widget Tree pronto          | ✓    | ✓*    | ✗   | ✗          | ✓     |
| Formulário dinâmico         | ✓    | ✓*    | ✗   | manual     | ✓     |
| Portabilidade (sem display) | ✗    | ✗     | ✗   | ✓          | ✓     |
| Complexidade de setup       | Baixa| Alta  | Média| Baixa     | Baixa |
| Público-alvo                | Geral| Geral | Geral| Dev/Power | Dev/Power |

\* via framework frontend (React, Svelte, etc.)

## Decisão

**TUI com tview.**

### Justificativa

| Critério                         | Bubble Tea     | tview          |
|----------------------------------|---------------|----------------|
| Widget `TreeView` nativo         | ✗ (manual)    | ✓              |
| `Form` + `InputField` + password | manual        | ✓              |
| CGO necessário                   | Não           | Não            |
| Binário estático                 | ✓             | ✓              |
| Referência de produção complexa  | —             | **k9s** (Kubernetes TUI) |
| Flexibilidade visual             | Alta          | Média          |

O critério determinante foi o widget `TreeView` nativo do tview: a árvore hierárquica é o elemento central do Abditum. Construí-la do zero em Bubble Tea adicionaria complexidade sem benefício. Formulários dinâmicos (`Form`, `InputField`, campos sensíveis) também estão prontos no tview.

O **k9s** — a TUI de Kubernetes mais utilizada — é construída sobre tview, demonstrando capacidade para UIs complexas, responsivas e com múltiplos painéis em produção.

Bubble Tea fica documentada como alternativa de alto valor caso a flexibilidade visual se torne requisito futuro.

### Opções GUI descartadas

Fyne, Wails e Gio foram descartadas: requerem CGO e/ou display gráfico, violando o requisito de binário estático portável sem dependências externas.
