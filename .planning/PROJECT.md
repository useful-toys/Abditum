# Abditum

## What This Is

Abditum é um cofre de senhas portátil e seguro, distribuído como único executável TUI (Terminal User Interface) em Go. Armazena e organiza credenciais e informações confidenciais em um único arquivo local criptografado (`.abditum`), sem dependência de serviços online, sem instalação e sem rastros no sistema além do próprio arquivo do cofre.

## Core Value

O usuário possui e controla completamente seus dados — o cofre é um arquivo que carrega consigo, protegido por criptografia forte, acessível apenas com a senha mestra, funcionando offline em qualquer sistema.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] Ciclo de vida completo do cofre: criar, abrir, salvar, salvar como, descartar, alterar senha mestra, bloquear, sair
- [ ] Exportar cofre para arquivo de intercâmbio não criptografado (com confirmação de risco)
- [ ] Importar cofre de arquivo de intercâmbio com mesclagem de estrutura e sobrescrita de conflitos
- [ ] Configurações persistidas no payload criptografado: tema visual, tempo de bloqueio automático, tempo de ocultação de campos sensíveis, tempo de limpeza de clipboard
- [ ] Navegação em árvore hierárquica de pastas e segredos
- [ ] Pasta virtual "Favoritos" (somente leitura, irmã da Pasta Geral)
- [ ] Busca de segredos por nome, nome de campo, valor de campo comum e observação (substring, case-insensitive, sem acentuação)
- [ ] Exibição temporária de campos sensíveis com ocultação automática configurável (padrão 15s)
- [ ] Cópia de qualquer campo para clipboard com limpeza automática configurável (padrão 30s)
- [ ] Indicadores de estado de sessão por segredo (✦ adicionado, ✎ modificado, ✗ excluído) e dirty state global no cabeçalho
- [ ] Criação de segredo a partir de modelo ou sem campos estruturados (apenas Observação)
- [ ] Duplicação de segredo com nome único automático e posição imediata após o original
- [ ] Edição de nome, valores de campos e observação de segredo
- [ ] Edição de estrutura de segredo: adicionar/renomear/reordenar/excluir campos (sem alterar tipo)
- [ ] Exclusão reversível de segredo (marcação até salvamento; restaurável antes disso)
- [ ] Movimentação e reordenação de segredos dentro e entre pastas
- [ ] Marcação e desmarcação de segredos como favoritos
- [ ] Criação, renomeação, movimentação e reordenação de pastas
- [ ] Exclusão de pasta com promoção de conteúdo (segredos e subpastas) para pasta pai
- [ ] Gerenciamento completo de modelos de segredo: criar, editar estrutura, renomear, excluir, criar a partir de segredo
- [ ] Bloqueio automático por inatividade configurável (padrão 5 min) com limpeza de memória e clear screen
- [ ] Bloqueio emergencial com tela disfarçada (atalho `⌃!⇧Q`) e limpeza imediata de dados sensíveis
- [ ] Salvamento atômico: gravação em `.abditum.tmp` → rename → backup `.abditum.bak` com rotação
- [ ] Proteção de memória com mlock/VirtualLock para dados sensíveis quando disponível no SO
- [ ] Detecção de modificação externa do arquivo do cofre antes de sobrescrever
- [ ] Detecção de acesso concorrente ao arquivo do cofre
- [ ] Compatibilidade retroativa de formato `.abditum` com migração em memória e escrita sempre no formato atual
- [ ] Interface TUI conforme design system: 4 zonas, painéis 35/65, inventário BMP-only, NO_COLOR, temas discretos e reativos
- [ ] Binário único cross-platform sem instalação: Windows, macOS e Linux 64-bit

### Out of Scope

- Senha falsa de coação (Duress Password) — complexidade de implementação segura; avaliação futura
- TOTP / Two-Factor Authentication — excluído permanentemente
- Backup gerenciado pela aplicação — responsabilidade do usuário
- Recuperação de dados corrompidos — AES-256-GCM não permite recuperação parcial
- Autenticação por keyfile ou token de hardware (YubiKey) — excluído permanentemente
- Armazenamento em nuvem — contraria filosofia offline e portátil; excluído permanentemente
- Múltiplos cofres abertos simultaneamente — invariante de design; excluído permanentemente
- App mobile ou web — TUI portátil por design; excluído permanentemente
- Modo somente leitura especial — falha informativa ao tentar salvar; sem tratamento adicional
- Recuperação de artefatos órfãos (`.abditum.tmp`, `.abditum.bak2`) — diferido para v2

## Context

**Stack técnica:**
- Linguagem: Go
- TUI: Bubbletea + Lipgloss (arquitetura Elm, modelo reativo)
- Criptografia: AES-256-GCM + Argon2id (m=256 MiB, t=3, p=4, saída 32 bytes)
- Formato de arquivo: `.abditum` — cabeçalho binário fixo 49 bytes (magic `ABDT` + versão uint8 + salt 32B + nonce 12B) + payload JSON UTF-8 criptografado + tag GCM 16 bytes
- AAD: cabeçalho completo autenticado pelo AES-256-GCM (sem checksum adicional)
- Salt: gerado via CSPRNG (`crypto/rand`) na criação; substituído apenas ao alterar senha mestra
- Nonce: gerado via CSPRNG a cada salvamento

**Conceitos do domínio:**
- **Cofre**: arquivo `.abditum` com Pasta Geral (raiz) obrigatória
- **Segredo**: item com campos (comuns + sensíveis) e campo Observação automático — imutável, sempre na última posição, não replicável em modelos
- **Campo sensível**: valor oculto por padrão (`••••`), revelável temporariamente
- **Modelo de segredo**: estrutura predefinida sem vínculo permanente com segredos criados a partir dele (nome do modelo registrado apenas como histórico)
- **Pasta virtual Favoritos**: somente leitura, percorre todos os segredos com `favorito=true` em profundidade
- **Exclusão reversível**: segredo marcado permanece visível com indicador `✗` até salvamento

**Regras de identidade:**
- Segredos: chave composta (pasta pai + nome) — únicos dentro da mesma pasta
- Pastas: chave composta (pasta pai + nome) — únicas dentro da mesma pasta pai
- Modelos: nome único globalmente
- Campos: sem restrição de duplicidade dentro do mesmo segredo ou modelo

**Segurança em sessão:**
- Senha mestra fornecida uma única vez; usada para toda criptografia da sessão; removida da memória (zeros) ao bloquear ou encerrar
- Terminal limpo (clear screen) ao bloquear ou encerrar
- Clipboard limpa ao bloquear, encerrar ou após timer configurável
- Dados sensíveis em memória protegida (mlock/VirtualLock) quando suportado pelo SO
- Nenhum log (stdout/stderr) com caminhos de cofre, nomes de segredos ou valores de campos

**Design system TUI:**
- Terminal mínimo: 80×24 (degradação sem crash abaixo do mínimo; truncamento com `…`)
- 4 zonas verticais empilhadas: Cabeçalho (2 linhas) | Área de trabalho (restante) | Barra de mensagens (1 linha) | Barra de comandos (1 linha)
- Painéis com divisão ≈35% (árvore/lista) / ≈65% (detalhe); separador `│`; conector `<╡` no item selecionado
- Inventário de símbolos BMP-only (U+0000–U+FFFF); sem emojis ou Nerd Fonts; todos com largura 1 coluna
- Bordas arredondadas (`╭╮╰╯`) para sobreposições; diálogos centralizados; ações na borda inferior
- NO_COLOR: todo estado crítico usa ≥2 camadas de comunicação (cor + símbolo ou atributo tipográfico)
- Temas discretos e reativos; identificador do tema gravado no payload criptografado

**Modelos padrão criados ao iniciar novo cofre:**
- Login: URL (comum), Usuário (comum), Senha (sensível)
- Cartão de Crédito: Titular (comum), Número (sensível), Validade (comum), CVV (sensível)
- Chave de API: Serviço (comum), Chave (sensível)

**Estrutura padrão criada ao iniciar novo cofre:**
- Pasta Geral (raiz) com subpastas: Sites e Apps, Financeiro

## Constraints

- **Tech stack**: Go + Bubbletea/Lipgloss — decisão definitiva
- **Distribuição**: binário único cross-platform; zero dependências externas em runtime; zero arquivos de config fora do cofre (exceto artefatos transitórios `.abditum.tmp`, `.abditum.bak`, `.abditum.bak2`)
- **Criptografia**: AES-256-GCM + Argon2id com parâmetros fixos hard-coded (v1: m=256MiB, t=3, p=4); sem calibração por máquina; mudanças exigem nova versão de formato
- **Plataforma**: Windows, macOS e Linux 64-bit
- **Unicode**: BMP apenas; sem emojis; sem Nerd Fonts; largura de símbolo 1 coluna (spinner `◐◓◑◒` monitorado em ambientes de largura ambígua)
- **Privacidade**: zero logs de dados sensíveis; nenhum rastro de estado fora do arquivo do cofre
- **Terminal mínimo**: 80×24; degradação graciosa sem crash

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| TUI em Go com Bubbletea/Lipgloss | Portabilidade, single binary, modelo Elm reativo, ecosystem maduro em Go | — Pending |
| AES-256-GCM + Argon2id (m=256MiB, t=3, p=4) | Padrão da indústria; parâmetros conservadores e fixos sem calibração por máquina | — Pending |
| Formato `.abditum`: cabeçalho binário 49B + JSON criptografado | Permite rejeitar arquivos inválidos antes de solicitar senha; JSON facilita migração | — Pending |
| AAD = cabeçalho completo | Autentica magic + versão + salt + nonce sem checksum adicional | — Pending |
| Salt por cofre; nonce por salvamento (ambos via CSPRNG) | Salt protege derivação de chave; nonce evita reutilização de par chave+nonce | — Pending |
| Exclusão reversível de segredos (marcação até salvamento) | Elimina confirmação modal por exclusão; reversível enquanto não persistido | — Pending |
| Observação automática em todo segredo (imutável, última posição) | Campo de notas sempre disponível sem poluir estrutura de modelos | — Pending |
| Salvamento atômico via temp file + rename + backup rotativo | Garante que falha durante gravação não corrompe cofre nem backup | — Pending |
| Temas discretos; identificador no payload criptografado | Preferência pessoal sem vazar informação fora do cofre | — Pending |
| Duress Password fora do escopo v1 | Complexidade de implementação segura supera benefício na versão inicial | — Pending |
| Exportação como arquivo não-criptografado com confirmação | Interoperabilidade consciente; usuário assume risco explicitamente | — Pending |
| Política de importação: mesclagem de pastas, sobrescrita de segredos/modelos conflitantes | Comportamento esperado e intencional para sincronização de versões | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: April 9, 2026 after initialization*
