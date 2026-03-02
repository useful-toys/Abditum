# ADR 001 — Linguagem de Implementação

**Status**: Aceita
**Data**: 2026-03-02

## Contexto

O Abditum precisa ser distribuído como uma aplicação desktop com as seguintes restrições não negociáveis:

- **Executável nativo**: sem necessidade de runtime externo instalado pelo usuário
- **Arquivo único**: toda a aplicação em um único binário — sem instaladores complexos, sem dependências de DLLs distribuídas separadamente
- **Footprint pequeno**: binário leve em disco e memória — o oposto de aplicações baseadas em Electron

Essas restrições foram o critério de corte principal. A escolha de linguagem foi consequência delas, não de preferência pessoal.

## Alternativas consideradas

### Node.js — Rejeitada

- Não produz executável nativo sem ferramentas adicionais como **pkg** ou **nexe**, que empacotam o runtime Node inteiro no binário
- Para UI desktop, exige **Electron** (ou Tauri como alternativa mais leve)
- Electron embute Chromium + Node — resulta em binários de **100–200 MB** e uso de memória comparável ao WhatsApp Desktop ou VS Code
- A UI necessariamente seria baseada em tecnologias web (HTML/CSS/JS), o que não é um requisito do Abditum e carrega complexidade desnecessária
- **Eliminada** pelos critérios de executável único e footprint pequeno

### Java — Rejeitada

- Requer **JVM** instalada no sistema do usuário, ou empacotar a JRE junto (aumenta drasticamente o tamanho)
- **GraalVM Native Image** permite compilar para binário nativo, mas tem limitações com reflexão e bibliotecas que dependem de carregamento dinâmico de classes — maturidade ainda limitada para aplicações complexas
- Ecossistema de UI desktop (JavaFX, Swing) está em declínio; JavaFX exige módulos separados
- Não atende ao critério de arquivo único sem esforço considerável de empacotamento
- **Eliminada** pela dependência de runtime e complexidade de distribuição

### Rust — Considerada, não adotada

- Produz binários nativos com static linking — **cumpre todos os critérios técnicos**
- Binários muito pequenos; sem runtime; excelente para aplicações de segurança (memory safety sem GC)
- Porém o ecossistema de UI desktop em Rust ainda é **imaturo**: `egui`, `iced`, `druid` estão em desenvolvimento ativo mas com APIs instáveis
- A curva de aprendizado do Rust (ownership, lifetimes) aumentaria significativamente o tempo de desenvolvimento para um projeto como este
- Para TUI, opções como `ratatui` são maduras — porém isso não diferencia Rust de Go nesse aspecto
- **Não adotada** por maturidade do ecossistema de UI e custo de desenvolvimento; permanece uma alternativa válida para v2 se os requisitos de segurança de memória se tornarem prioritários

### Go — Adotada ✓

- Produz binários nativos com **static linking** por padrão (`CGO_ENABLED=0` para TUI)
- Binários pequenos: aplicação equivalente a este escopo fica tipicamente entre **5–20 MB** sem dependências externas
- Stdlib rica para o domínio do projeto: `crypto/aes`, `crypto/cipher`, `crypto/rand`, `encoding/json` — sem dependências de terceiros para o core de segurança
- Ecossistema de UI viável para os requisitos:
  - **TUI**: `tview` (widgets prontos, incluindo TreeView) e `Bubble Tea` (arquitetura Elm, ecossistema charm.sh)
  - **GUI Windows nativo**: `walk` — bindings Go para Win32, produz UI nativa sem Electron
- Compilação simples e rápida; ferramental (`go build`, `go test`, `go mod`) sem configuração complexa
- Sem runtime externo: o binário Go é autocontido

## Consequências

- O binário final será estático e autocontido — distribuição é copiar um único arquivo
- Para GUI, CGO será necessário (implica toolchain C disponível no ambiente de build)
- Para TUI, CGO pode ser evitado completamente
- A escolha do framework de UI permanece em aberto (ver `decisions/adr-002-ui-framework.md`) e pode influenciar o tamanho final do binário
- Rust permanece como alternativa técnica válida; uma reavaliação pode ocorrer se o ecossistema de UI evoluir ou se requisitos de segurança de memória se tornarem críticos
