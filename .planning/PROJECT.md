# Abditum

## What This Is

Abditum e um cofre de senhas portatil, seguro e facil de usar, distribuido como um unico executavel com interface TUI moderna. Ele permite que a pessoa usuaria crie, abra, bloqueie, organize e gerencie segredos localmente, com controle total dos dados e sem dependencia de nuvem, instalacao ou arquivos auxiliares fora do proprio cofre.

## Core Value

Oferecer um cofre de segredos realmente portatil e offline-first, no qual dados sensiveis permanecem sob controle total da pessoa usuaria em um unico arquivo criptografado e executavel independente.

## Requirements

### Validated

(None yet - ship to validate)

### Active

- [ ] Criar e desbloquear cofres criptografados com senha mestra, bloqueio manual e bloqueio automatico por inatividade.
- [ ] Gerenciar segredos em hierarquia de pastas com busca, favoritos, duplicacao, movimentacao e reordenacao.
- [ ] Criar e editar modelos de segredo pre-definidos e personalizados, mantendo segredos como snapshots independentes.
- [ ] Operar tudo por uma TUI moderna, portatil e responsiva, com foco forte em privacidade, seguranca e usabilidade via teclado.

### Out of Scope

- Armazenamento em nuvem - contradiz a proposta de controle local e zero dependencia de terceiros.
- Multiplos cofres abertos simultaneamente - aumenta complexidade sem reforcar o valor principal do v1.
- App mobile ou web - o produto e TUI portatil por design.
- Tags - pastas e grupos sao suficientes para o v1.
- Historico de versoes de segredos - adiado para evitar ampliar escopo inicial.

## Context

O projeto nasce com foco em portabilidade extrema: um binario unico, sem instalacao, sem servicos externos e sem arquivos de configuracao fora do proprio cofre. O arquivo `.abditum` deve ser autossuficiente, armazenando dados criptografados, configuracoes do cofre, modelos de segredo e metadados necessarios para compatibilidade futura.

O cofre precisa suportar um fluxo completo de ciclo de vida: criar, abrir, salvar, salvar como, descartar alteracoes nao salvas, alterar senha mestra, exportar e importar JSON em texto puro com confirmacoes de seguranca, backups `.bak` e persistencia atomica usando arquivo temporario local. O produto tambem precisa limpar dados sensiveis da memoria ao bloquear ou fechar, limpar a area de transferencia automaticamente e reduzir risco de shoulder surfing com ocultacao rapida da interface.

A modelagem definida usa hierarquia recursiva de pastas e segredos, ordenacao por posicao preservada no JSON, segredos derivados de modelos como snapshot sem vinculo posterior, IDs curtos para pastas e segredos, observacao implicita em todo segredo e configuracoes embutidas no cofre. O produto deve abrir arquivos da versao N-1, evitando quebrar cofres ja existentes.

Na experiencia de uso, a TUI ocupa todo o terminal, trabalha com paineis lateral e principal, suporta teclado integralmente e mouse de forma complementar, exibe ajuda contextual e feedback claro para operacoes longas, destrutivas ou sensiveis. Campos sensiveis devem iniciar ocultos e exibicao temporaria deve respeitar timeout configuravel.

## Constraints

- **Tech stack**: Go + Bubble Tea/teatest v2 - decisoes explicitas do projeto para entregar binario unico e TUI moderna.
- **Security**: AES-256-GCM + Argon2id com parametros rigidos - base de criptografia e resistencia a ataques offline e brute force.
- **Storage**: JSON criptografado em arquivo `.abditum` - garante portabilidade e inspecao estruturada do payload antes da criptografia.
- **Compatibility**: Windows, macOS e Linux - o executavel precisa funcionar de forma portatil nos tres sistemas.
- **Portability**: Nenhum arquivo de configuracao, log sensivel, rastro temporario ou dependencia de instalacao fora do cofre - requisito central do produto.
- **Reliability**: Salvamento atomico com `.abditum.tmp` e backup `.abditum.bak` - evita corrupcao e perda de dados durante gravacao.
- **Architecture**: DDD com modificacoes centralizadas em Managers - protege invariantes e evita mutacao direta insegura das entidades.
- **Testing**: Cobertura forte com testes unitarios, integracao, fluxos TUI, golden files 80x24 e CI obrigatorio - confiabilidade e didatica sao requisitos explicitos.
- **UX**: Layout responsivo para diferentes tamanhos de terminal e foco em navegacao por teclado - o app precisa continuar utilizavel em ambientes restritos.

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Distribuir como binario unico portatil | Reforca uso discreto, offline e sem instalacao | - Pending |
| Manter dados, modelos e configuracoes dentro do proprio cofre | Garante portabilidade extrema e evita rastros externos | - Pending |
| Modelos de segredo geram snapshots sem referencia viva | Permite evoluir modelos sem quebrar segredos existentes | - Pending |
| Usar TUI full-screen com painel lateral e painel principal | Equilibra navegacao estrutural e edicao detalhada no terminal | - Pending |
| Adotar Go com DDD e Bubble Tea | Combina binario unico, manutencao estruturada e TUI robusta | - Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? -> Move to Out of Scope with reason
2. Requirements validated? -> Move to Validated with phase reference
3. New requirements emerged? -> Add to Active
4. Decisions to log? -> Add to Key Decisions
5. "What This Is" still accurate? -> Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check - still the right priority?
3. Audit Out of Scope - reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-03-24 after initialization*
