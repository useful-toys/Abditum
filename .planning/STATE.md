# Project State — Abditum

**Last updated:** 2026-03-27 (sessão de revisão de documentos)
**Current phase:** Ready for Phase 1
**Milestone:** v1.0

## Status

| Artifact | Status |
|----------|--------|
| PROJECT.md | ✓ Created |
| config.json | ✓ Created |
| REQUIREMENTS.md | ✓ 55 requirements defined, traceability mapped |
| research/ | ✓ STACK, FEATURES, ARCHITECTURE, PITFALLS, SUMMARY |
| ROADMAP.md | ✓ 11 phases, all requirements mapped |
| modelo-dominio.md | ✓ Revisado — ordem corrigida (subpastas antes de segredos) |
| arquitetura.md | ✓ Revisado — versões Bubble Tea v2 explicitadas |
| formato-arquivo-abditum.md | ✓ Revisado — Argon2id m=256 MiB incorporado ao REQUIREMENTS |
| descricao.md | ✓ Revisado — ordem corrigida em fluxo \"Visualizar hierarquia\" |

## Current Phase

None started. Run `/gsd-plan-phase 1` to begin.

## Phase History

*(empty — no phases completed yet)*

## Open Decisions

- Help overlay de teclado (`?` key / footer hints) — especificar antes da Phase 5 (TUI scaffold)

## Notes

- Bubble Tea v2 import path: `charm.land/bubbletea/v2` (não v1)
- Argon2id: m=256 MiB, t=3, p=4 (conforme formato-arquivo-abditum.md)
- Lixeira: lista in-memory no Manager, sem campo no Segredo
- Exibição dentro de pasta: subpastas primeiro, segredos depois
- NanoID: implementar internamente com `crypto/rand`
- Datetime: RFC 3339 UTC
- Favoritos: DFS na ordem do JSON
