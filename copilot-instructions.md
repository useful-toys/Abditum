<!-- GSD:project-start source:PROJECT.md -->
## Project

**Abditum**

Abditum é um cofre de senhas portátil e seguro, distribuído como único executável TUI (Terminal User Interface) em Go. Armazena e organiza credenciais e informações confidenciais em um único arquivo local criptografado (`.abditum`), sem dependência de serviços online, sem instalação e sem rastros no sistema além do próprio arquivo do cofre.

**Core Value:** O usuário possui e controla completamente seus dados — o cofre é um arquivo que carrega consigo, protegido por criptografia forte, acessível apenas com a senha mestra, funcionando offline em qualquer sistema.

### Constraints

- **Tech stack**: Go + Bubbletea/Lipgloss — decisão definitiva
- **Distribuição**: binário único cross-platform; zero dependências externas em runtime; zero arquivos de config fora do cofre (exceto artefatos transitórios `.abditum.tmp`, `.abditum.bak`, `.abditum.bak2`)
- **Criptografia**: AES-256-GCM + Argon2id com parâmetros fixos hard-coded (v1: m=256MiB, t=3, p=4); sem calibração por máquina; mudanças exigem nova versão de formato
- **Plataforma**: Windows, macOS e Linux 64-bit
- **Unicode**: BMP apenas; sem emojis; sem Nerd Fonts; largura de símbolo 1 coluna (spinner `◐◓◑◒` monitorado em ambientes de largura ambígua)
- **Privacidade**: zero logs de dados sensíveis; nenhum rastro de estado fora do arquivo do cofre
- **Terminal mínimo**: 80×24; degradação graciosa sem crash
<!-- GSD:project-end -->

<!-- GSD:stack-start source:STACK.md -->
## Technology Stack

Technology stack not yet documented. Will populate after codebase mapping or first phase.
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

Conventions not yet established. Will populate as patterns emerge during development.
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

Architecture not yet mapped. Follow existing patterns found in the codebase.
<!-- GSD:architecture-end -->

<!-- GSD:skills-start source:skills/ -->
## Project Skills

No project skills found. Add skills to any of: `.github/skills/`, `.agents/skills/`, `.cursor/skills/`, or `.github/skills/` with a `SKILL.md` index file.
<!-- GSD:skills-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd-quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd-debug` for investigation and bug fixing
- `/gsd-execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->



<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd-profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
