# Abditum

## What This Is

Abditum é um cofre de senhas portátil, seguro e offline, distribuído como um único binário executável Go com interface TUI moderna (Bubble Tea). Permite que o usuário armazene e gerencie credenciais e informações confidenciais organizadas em pastas hierárquicas, com segredos compostos por campos comuns e sensíveis, protegidos por criptografia AES-256-GCM e derivação de chave Argon2id — sem dependências de nuvem, sem instalação e sem rastreamento.

## Core Value

O usuário tem controle total e exclusivo sobre seus segredos: os dados existem apenas no arquivo `.abditum` e na memória da sessão ativa — nenhum terceiro, serviço externo ou processo paralelo tem acesso a eles.

## Requirements

### Validated

(None yet — ship to validate)

### Active

#### Cofre — Ciclo de Vida
- [ ] **VAULT-01**: Usuário pode criar novo cofre em arquivo com senha mestra (confirmação dupla + avaliação de força)
- [ ] **VAULT-02**: Ao criar cofre, Pasta Geral é criada automaticamente com subpastas "Sites e Apps" e "Financeiro" e modelos padrão (Login, Cartão de Crédito, Chave de API)
- [ ] **VAULT-03**: Usuário pode abrir cofre a partir de arquivo existente com senha mestra
- [ ] **VAULT-04**: Erros de abertura distinguem autenticação (nova tentativa) de integridade (bloqueio total)
- [ ] **VAULT-05**: Usuário pode salvar cofre no arquivo atual (sem re-solicitar senha; segredos marcados excluídos removidos permanentemente)
- [ ] **VAULT-06**: Usuário pode salvar cofre em outro arquivo (arquivo de trabalho passa a ser o novo; segredos marcados excluídos removidos)
- [ ] **VAULT-07**: Se arquivo foi modificado externamente antes de salvar, usuário recebe aviso e opções: Sobrescrever / Salvar como novo / Cancelar
- [ ] **VAULT-08**: Usuário pode descartar alterações não salvas e recarregar cofre do arquivo
- [ ] **VAULT-09**: Usuário pode alterar senha mestra (confirmação dupla + avaliação de força; salva imediatamente e irreversivelmente)
- [ ] **VAULT-10**: Cofre bloqueia automaticamente após tempo configurável de inatividade (padrão: 5 min; qualquer interação reseta o timer)
- [ ] **VAULT-11**: Usuário pode bloquear o cofre manualmente
- [ ] **VAULT-12**: Ao bloquear, senha mestra é sobrescrita na memória e terminal é limpo (clear screen)
- [ ] **VAULT-13**: Usuário pode sair da aplicação (confirmação com opções se houver alterações pendentes; limpeza de memória e terminal ao sair)
- [ ] **VAULT-14**: Usuário pode exportar cofre para JSON (com aviso de risco e confirmação; segredos marcados excluídos não incluídos)
- [ ] **VAULT-15**: Usuário pode importar cofre de JSON (mesclagem de pastas por caminho; tratamento de conflitos de nome; segredos/modelos com ID duplicado tratados conforme regras documentadas)
- [ ] **VAULT-16**: Usuário pode configurar tempo de bloqueio por inatividade, tempo de ocultação de campo sensível e tempo de limpeza de clipboard

#### Consulta de Segredos
- [ ] **QUERY-01**: Usuário visualiza cofre com sua hierarquia de pastas e segredos
- [ ] **QUERY-02**: Usuário pode buscar segredos por nome, nome de campo, valor de campo comum ou observação (substring, case-insensitive, sem acento; campos sensíveis excluídos; segredos marcados excluídos ocultos)
- [ ] **QUERY-03**: Usuário visualiza segredo com nome, campos e observação
- [ ] **QUERY-04**: Usuário pode revelar temporariamente o valor de campo sensível (ocultação automática após timer configurável; padrão: 15 s)
- [ ] **QUERY-05**: Usuário pode copiar qualquer campo para clipboard (limpeza automática ao bloquear/sair ou após timer configurável; padrão: 30 s)
- [ ] **QUERY-06**: Segredos exibem indicadores de estado de sessão na listagem: adicionado (criado na sessão), modificado (alterado na sessão), excluído (marcado para remoção); segredos sem alteração não exibem indicador
- [ ] **QUERY-07**: Pasta virtual "Favoritos" exibida como nó irmão da Pasta Geral na árvore (acima dela); lista todos os segredos com `favorito = true`, percorridos em profundidade; somente leitura — não é possível criar, mover ou excluir segredos diretamente a partir desta vista; não pode ser renomeada, movida ou excluída

#### Gerenciamento de Segredos
- [ ] **SEC-01**: Usuário pode criar segredo a partir de modelo existente ou sem modelo (somente Observação), escolhendo a pasta
- [ ] **SEC-02**: Usuário pode duplicar segredo (cópia na mesma pasta imediatamente após original, nome ajustado automaticamente, histórico de modelo preservado)
- [ ] **SEC-03**: Usuário pode editar segredo: nome, valores de campos e observação
- [ ] **SEC-04**: Usuário pode alterar estrutura do segredo: adicionar campo (nome + tipo), renomear campo, reordenar campos, excluir campo
- [ ] **SEC-05**: Observação existe em todo segredo: não pode ser renomeada, excluída ou movida de posição — sempre na última posição; é campo comum
- [ ] **SEC-06**: Usuário pode favoritar/desfavoritar segredo
- [ ] **SEC-07**: Usuário pode marcar/desmarcar segredo para exclusão (permanece visível sinalizado; removido ao salvar)
- [ ] **SEC-08**: Usuário pode mover segredo para outra pasta
- [ ] **SEC-09**: Usuário pode reordenar segredo dentro da mesma pasta (ordem persistida ao salvar)

#### Gerenciamento de Pastas
- [ ] **FOLDER-01**: Usuário pode criar pasta dentro de outra pasta (nome único dentro da pasta pai)
- [ ] **FOLDER-02**: Usuário pode renomear pasta (nome único dentro da pasta pai; Pasta Geral não pode ser renomeada)
- [ ] **FOLDER-03**: Usuário pode mover pasta (validação contra ciclos; nome único no destino; Pasta Geral não pode ser movida)
- [ ] **FOLDER-04**: Usuário pode reordenar pasta dentro da mesma pasta (ordem persistida ao salvar)
- [ ] **FOLDER-05**: Usuário pode excluir pasta (segredos e subpastas promovidos para pasta pai; conflito de nome em segredo promovido → renomeado com sufixo numérico; conflito de nome em subpasta promovida → conteúdo mesclado; Pasta Geral não pode ser excluída)

#### Gerenciamento de Modelos
- [ ] **TPL-01**: Usuário pode criar modelo de segredo com campos personalizados (nome + tipo)
- [ ] **TPL-02**: Usuário pode renomear modelo (nome único entre modelos)
- [ ] **TPL-03**: Usuário pode alterar estrutura do modelo: adicionar campo, renomear campo, alterar tipo, reordenar campos, excluir campo (sem efeito em segredos já criados)
- [ ] **TPL-04**: Usuário pode excluir modelo
- [ ] **TPL-05**: Usuário pode criar modelo a partir de segredo existente (todos os campos com nome 'Observação' são excluídos — tanto a Observação automática quanto campos do usuário com esse nome)

#### Requisitos Não Funcionais / Técnicos
- [ ] **SEC-CRYPTO-01**: Criptografia AES-256-GCM; derivação de chave Argon2id; dependências de crypto exclusivamente de stdlib Go + `golang.org/x/crypto`
- [ ] **SEC-PRIV-01**: Zero logs de stdout/stderr com caminhos de arquivo, nomes de segredos ou valores de campos
- [ ] **COMPAT-01**: Compatibilidade retroativa de formato de arquivo: versão N abre arquivos de versões anteriores (migração em memória; sempre salva no formato atual)
- [ ] **ATOMIC-01**: Salvamento atômico via `.abditum.tmp` com rollback e backup `.abditum.bak` / `.abditum.bak2`
- [ ] **PORT-01**: Binário único executável, cross-platform: Windows, macOS, Linux
- [ ] **MEM-01**: Ao bloquear/sair, sobrescrever senha mestra e descartar buffers sensíveis; usar `mlock`/`VirtualLock` quando disponível
- [ ] **CI-01**: CI obrigatório: build + lint + testes completos em todo push

### Out of Scope

- **TOTP (autenticação de dois fatores)** — excluído permanentemente; fora do foco de gerenciamento de credenciais estáticas
- **Backup automático** — responsabilidade do usuário; a app não gerencia cópias de segurança
- **Recuperação de dados corrompidos** — a criptografia não permite recuperação parcial; design intencional
- **Keyfile / Token de hardware (YubiKey)** — excluído permanentemente; fora do modelo de segurança atual
- **Armazenamento em nuvem** — contraria filosofia offline e portátil; excluído permanentemente
- **Múltiplos cofres simultâneos** — invariante de design; só um cofre ativo por vez
- **App mobile ou web** — TUI portátil é o produto; excluído permanentemente
- **Duress password (senha falsa de coação)** — planejado para v2; não pertence a v1
- **Gerador de senhas** — planejado para v2
- **QR Code de campo** — planejado para v2
- **Tags** — planejado para v2
- **Histórico de versões de segredos** — planejado para v2
- **Relatório de saúde do cofre** — planejado para v2
- **Recuperação de artefatos órfãos** — planejado para v2

## Context

**Estado atual:** O projeto possui documentação de requisitos, arquitetura, domínio, BDD e RUP extensamente elaborada. Nenhuma implementação Go existe ainda — estamos na fase de inicialização do desenvolvimento.

**Stack definida:**
- Linguagem: **Go** — compilado como binário único, sem runtime externo
- TUI: **Bubble Tea (Charm)** — modelo Elm de atualização de estado
- Testes de TUI: **teatest/v2**
- Crypto: **AES-256-GCM + Argon2id** (stdlib Go + `golang.org/x/crypto`)

**Estrutura de pacotes planejada:**
```
cmd/abditum/          -- ponto de entrada (main)
internal/vault/       -- domínio e lógica de negócio (Manager, entidades, regras)
internal/crypto/      -- derivação de chave e criptografia/descriptografia
internal/storage/     -- leitura/escrita do arquivo .abditum (formato binário + salvamento atômico)
internal/tui/         -- interface TUI (modelos Bubble Tea, telas, componentes, navegação)
```

**Padrão arquitetural:** Manager centraliza toda mutação do domínio. A TUI interage com o domínio exclusivamente via Manager.

**Audiência:** Usuários que valorizam privacidade local e não querem depender de serviços em nuvem — podem ser técnicos ou não, desde que confortáveis com TUI.

**Filosofia de segurança:** Zero-knowledge, local-only, dependências mínimas, sem logs sensíveis.

## Constraints

- **Tech Stack**: Go + Bubble Tea + teatest/v2 — definido e não negociável
- **Crypto**: Apenas stdlib Go e `golang.org/x/crypto` para operações criptográficas — sem libs de terceiros
- **Portabilidade**: Binário único, sem instalação, sem configuração externa — Windows + macOS + Linux
- **Privacidade**: Nenhum log com dados sensíveis (caminhos, nomes, valores)
- **Offline**: Sem acesso a rede, sem cloud, sem telemetria de qualquer natureza
- **CI**: Build + lint + testes obrigatórios em todo push — mudanças que quebrem o build não são aceitas
- **Comentários**: Política generosa — código deve ser acessível a leitores menos familiarizados com Go, Bubble Tea e criptografia

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Go como linguagem | Binário único sem runtime externo; tipagem forte; performance adequada para TUI | — Pending |
| Bubble Tea (Charm) para TUI | Modelo Elm bem estabelecido; componentes reutilizáveis; suporte multiplataforma | — Pending |
| AES-256-GCM + Argon2id | Padrão de mercado para criptografia simétrica + derivação de chave resistente a brute-force | — Pending |
| Manager pattern para domínio | Centraliza regras de negócio; impede mutação direta de entidades pela TUI | — Pending |
| JSON criptografado como formato | Legível para exportação, versionável com migração em memória | — Pending |
| Salvamento atômico via .tmp + rollback | Garante que falha durante gravação não corrompe o arquivo principal | — Pending |
| Sem recuperação de dados corrompidos | A criptografia autenticada (GCM) detecta corrupção; recuperação parcial comprometeria a integridade | — Pending |
| Exclusão soft (marcar para excluir) | Dá ao usuário oportunidade de reverter antes de salvar; remoção permanente só ocorre ao persistir | — Pending |
| Observação automática em todo segredo | Campo de notas livre sempre acessível, sem ser parte do modelo de template | — Pending |

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
*Last updated: 2026-03-27 after initialization*
