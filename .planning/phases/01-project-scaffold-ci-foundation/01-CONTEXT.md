# Phase 1: Project Scaffold + CI Foundation - Context

**Gathered:** 2026-03-27
**Status:** Ready for planning

<domain>
## Phase Boundary

Configuração inicial do repositório Go: `go.mod` com todas as dependências, tree de diretórios canônica com stubs de pacotes, GitHub Actions CI verde no Linux, Makefile com targets padrão, e `.golangci.yml` com linters configurados. Sem lógica de negócio, sem TUI — apenas a fundação estática do projeto.

</domain>

<decisions>
## Implementation Decisions

### Module Path
- **D-01:** Module path: `github.com/useful-toys/abditum` (corresponde ao owner/repo real — não o placeholder `github.com/user/abditum` do ROADMAP).

### Versão do Go
- **D-02:** `go 1.26.1+` no `go.mod`; `go-version: '1.26.x'` no step de setup-go do CI.

### Branch do CI
- **D-03:** CI (`ci.yml`) dispara em push e pull_request para `main` (branch padrão após renomeação de `master` → `main`).

### the agent's Discretion
- Formato exato do `.golangci.yml` (thresholds de severidade, configuração de nolint para stubs vazios)
- Estrutura interna do `Makefile` (flags adicionais, phony targets)
- Nome do binário de output no target `build` do Makefile

</decisions>

<specifics>
## Specific Ideas

- O ROADMAP menciona `github.com/user/abditum` como module path — substituir por `github.com/useful-toys/abditum` (D-01).
- O ROADMAP lista `github.com/matoous/go-nanoid/v2` na Plan 1 — esta dependência NÃO será necessária (não precisaremos de nanoid nesta fase). O planejador deve remover essa referência ao criar as tasks.
- O CI deve usar `CGO_ENABLED=0` como variável de ambiente global no job (não apenas no comando de build) para que `go test` também use static linking.
- Branch padrão agora é `main` (não `master`).

</specifics>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

### Requisitos da fase
- `.planning/REQUIREMENTS.md` §COMPAT-01 — Binário único, CGO_ENABLED=0, sem runtime externo
- `.planning/REQUIREMENTS.md` §CI-01 — CI obrigatório (implícito no roadmap)
- `.planning/ROADMAP.md` §Phase 1 — Detalhamento de plans, UAT e pitfall watch

### Arquitetura e convenções do projeto
- `arquitetura.md` — Import paths canônicos Charm (`charm.land/*`), restrições de dependências (`no net`, `no math/rand`), política de comentários generosos
- `formato-arquivo-abditum.md` — Referência para entender por que certas deps são necessárias (não afeta Phase 1 diretamente, mas é contexto geral)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- Nenhum — repositório ainda sem código Go. Esta é a fase de criação do scaffold.

### Padrões estabelecidos
- Nenhum ainda — esta fase estabelece os padrões para todas as fases seguintes.

</code_context>
