# Requirements: Abditum

**Defined:** 2026-03-27
**Core Value:** O usuário tem controle total e exclusivo sobre seus segredos — dados existem apenas no arquivo `.abditum` e na memória da sessão ativa.

## v1 Requirements

### Criptografia e Segurança

- [ ] **CRYPTO-01**: Criptografia AES-256-GCM com nonce único por operação de escrita; derivação de chave Argon2id (t=3, m=256 MiB, p=4, keyLen=32); parâmetros fixos por versão de formato, sem calibração por máquina
- [ ] **CRYPTO-02**: Dependências de criptografia exclusivamente de stdlib Go e `golang.org/x/crypto` — sem libs de terceiros
- [ ] **CRYPTO-03**: Todos os dados sensíveis (senha mestra, buffers de chave) manipulados exclusivamente como `[]byte` zeráveis — nunca como `string`
- [ ] **CRYPTO-04**: Ao bloquear ou sair, senha mestra e buffers sensíveis são sobrescritos com zeros e descartados antes de retornar ao shell
- [ ] **CRYPTO-05**: Uso de `mlock`/`VirtualLock` quando disponíveis no SO para impedir swap de memória sensível para disco; aplicação opera normalmente quando indisponível
- [ ] **CRYPTO-06**: Zero logs de stdout/stderr que contenham caminhos de arquivo de cofre, nomes de segredos ou valores de campos

### Ciclo de Vida do Cofre

- [ ] **VAULT-01**: Usuário pode criar novo cofre em arquivo com senha mestra; confirmação dupla obrigatória; avaliação de força exibida (sem bloquear operação)
- [ ] **VAULT-02**: Ao criar cofre, Pasta Geral é criada automaticamente com subpastas "Sites e Apps" e "Financeiro" e modelos padrão: Login (URL, Usuário, Senha), Cartão de Crédito (Titular, Número, Validade, CVV), Chave de API (Serviço, Chave)
- [ ] **VAULT-03**: Usuário pode abrir cofre a partir de arquivo existente com senha mestra
- [ ] **VAULT-04**: Erros de abertura classificados em 4 categorias: tipo de arquivo inválido (magic incorreto → rejeitar), versão incompatível (versão_formato > suportado → rejeitar), autenticação (tag GCM inválida → nova tentativa permitida), integridade (JSON inválido ou Pasta Geral ausente → rejeitar); mensagens sempre genéricas, sem detalhes técnicos
- [ ] **VAULT-05**: Se a Pasta Geral não existir no arquivo aberto, rejeitar com mensagem de erro (arquivo inválido) — sem tentativa de recriar
- [ ] **VAULT-06**: Usuário pode salvar cofre no arquivo atual sem re-solicitar senha; segredos marcados para exclusão são removidos permanentemente
- [ ] **VAULT-07**: Usuário pode salvar cofre em outro arquivo; arquivo de trabalho passa a ser o novo; segredos marcados removidos; não pode ser o mesmo arquivo atual
- [ ] **VAULT-08**: Se arquivo foi modificado externamente desde última leitura/salvamento, avisar usuário antes de salvar (opções: Sobrescrever / Salvar como novo arquivo / Cancelar)
- [ ] **VAULT-09**: Usuário pode descartar alterações não salvas e recarregar cofre do arquivo (usar senha ativa da sessão; avisar se arquivo foi modificado externamente)
- [ ] **VAULT-10**: Usuário pode alterar senha mestra (confirmação dupla + avaliação de força; salva imediatamente e de forma irrevogável)
- [ ] **VAULT-11**: Cofre bloqueia automaticamente após tempo configurável de inatividade (padrão: 5 min); qualquer interação reseta o timer
- [ ] **VAULT-12**: Usuário pode bloquear o cofre manualmente
- [ ] **VAULT-13**: Ao bloquear, senha mestra é sobrescrita em memória, buffers sensíveis descartados e terminal limpo (clear screen incluindo scrollback `\033[3J`)
- [ ] **VAULT-14**: Usuário pode sair da aplicação; se houver alterações pendentes, confirmação com opções: Salvar e Sair / Descartar e Sair / Cancelar; sem confirmação se não houver alterações; mesma limpeza de memória e terminal do bloqueio
- [ ] **VAULT-15**: Usuário pode exportar cofre para JSON (aviso de risco + confirmação; segredos marcados excluídos omitidos; pastas, segredos ativos e modelos incluídos; configurações de timers não exportadas)
- [ ] **VAULT-16**: Usuário pode importar cofre de JSON (regras de mesclagem por caminho de pasta; tratamento de conflitos de nome com sufixo numérico; segredos com ID duplicado inseridos como independentes; modelos com ID duplicado substituem silenciosamente o existente)
- [ ] **VAULT-17**: Usuário pode configurar: tempo de bloqueio por inatividade (padrão 5 min), tempo de ocultação de campo sensível (padrão 15 s), tempo de limpeza de clipboard (padrão 30 s); todos os timers são obrigatórios

### Salvamento Atômico

- [ ] **ATOMIC-01**: Gravação do cofre sempre via `.abditum.tmp` no mesmo diretório; renomeação atômica substitui o arquivo original somente após gravação bem-sucedida; em falha, `.abditum.tmp` é apagado imediatamente
- [ ] **ATOMIC-02**: Ao substituir arquivo existente, backup `.abditum.bak` é mantido; se `.abditum.bak` já existe, renomear para `.abditum.bak2` antes de gerar novo backup; em falha após backup gerado, restaurar `.abditum.bak2` → `.abditum.bak` quando possível
- [ ] **ATOMIC-03**: Criação de novo cofre e salvamento em caminho vazio não usam `.abditum.tmp` — gravação direta no destino final
- [ ] **ATOMIC-04**: Renomeação atômica em Windows usa `MoveFileEx` com `MOVEFILE_REPLACE_EXISTING` (não `os.Rename` nativo)

### Compatibilidade e Portabilidade

- [ ] **COMPAT-01**: Aplicação construída como binário único executável, sem runtime externo, sem arquivos de configuração ou dados fora do cofre (exceto `.abditum.tmp`, `.abditum.bak`, `.abditum.bak2`)
- [ ] **COMPAT-02**: Suporte a Windows, macOS e Linux
- [ ] **COMPAT-03**: Compatibilidade retroativa de formato: versão N abre arquivos de versões anteriores; migração em memória; sempre salva no formato atual

### Consulta de Segredos

- [ ] **QUERY-01**: Usuário visualiza cofre com hierarquia de pastas e segredos
- [ ] **QUERY-02**: Usuário pode buscar segredos por nome, nome de campo, valor de campo comum ou observação (substring, case-insensitive, normalização de acentuação; campos sensíveis excluídos da busca; segredos marcados excluídos não aparecem)
- [ ] **QUERY-03**: Usuário visualiza segredo com nome, todos os campos e observação
- [ ] **QUERY-04**: Usuário pode revelar temporariamente o valor de campo sensível; valor ocultado automaticamente após timer configurável (padrão: 15 s)
- [ ] **QUERY-05**: Usuário pode copiar valor de qualquer campo para clipboard; clipboard limpa automaticamente ao bloquear, sair ou após timer configurável (padrão: 30 s); limpeza de clipboard depende de suporte do SO (Wayland: best-effort, requer `wl-clipboard` ou `xclip` em execução)

### Gerenciamento de Segredos

- [ ] **SEC-01**: Usuário pode criar segredo a partir de modelo existente ou sem modelo (apenas Observação); escolhe a pasta no momento da criação
- [ ] **SEC-02**: Usuário pode duplicar segredo (cópia na mesma pasta imediatamente após original; nome ajustado automaticamente: "Segredo (1)", "Segredo (2)"; histórico de modelo preservado)
- [ ] **SEC-03**: Usuário pode editar segredo: nome, valores de campos e observação (sem alterar estrutura)
- [ ] **SEC-04**: Usuário pode alterar estrutura do segredo: adicionar campo (nome + tipo), renomear campo, reordenar campos, excluir campo; tipo de campo não pode ser alterado; Observação não participa de reordenação, não pode ser renomeada/excluída/movida
- [ ] **SEC-05**: Observação existe automaticamente em todo segredo na última posição; campo comum; não pode ser renomeada, excluída ou movida
- [ ] **SEC-06**: Usuário pode favoritar/desfavoritar segredo
- [ ] **SEC-07**: Usuário pode marcar/desmarcar segredo para exclusão; segredo permanece visível sinalizado; removido permanentemente ao salvar
- [ ] **SEC-08**: Usuário pode mover segredo para outra pasta
- [ ] **SEC-09**: Usuário pode reordenar segredo dentro da mesma pasta; ordem persistida ao salvar

### Gerenciamento de Pastas

- [ ] **FOLDER-01**: Usuário pode criar pasta dentro de outra pasta; nome único dentro da pasta pai
- [ ] **FOLDER-02**: Usuário pode renomear pasta; nome único dentro da pasta pai; Pasta Geral não pode ser renomeada
- [ ] **FOLDER-03**: Usuário pode mover pasta; validação contra ciclos hierárquicos; nome único no destino; Pasta Geral não pode ser movida
- [ ] **FOLDER-04**: Usuário pode reordenar pasta dentro da mesma pasta; ordem persistida ao salvar
- [ ] **FOLDER-05**: Usuário pode excluir pasta; segredos e subpastas promovidos para pasta pai imediata; conflitos de nome em subpastas resolvidos com sufixo numérico (usuário avisado); Pasta Geral não pode ser excluída

### Gerenciamento de Modelos de Segredo

- [ ] **TPL-01**: Usuário pode criar modelo de segredo com campos personalizados (nome + tipo: comum ou sensível)
- [ ] **TPL-02**: Usuário pode renomear modelo; nome único entre modelos
- [ ] **TPL-03**: Usuário pode alterar estrutura do modelo: adicionar campo, renomear campo, alterar tipo de campo, reordenar campos, excluir campo; alterações não afetam segredos já criados
- [ ] **TPL-04**: Usuário pode excluir modelo
- [ ] **TPL-05**: Usuário pode criar modelo a partir de segredo existente; Observação automática é ignorada; campo usuário chamado "Observação" incluído normalmente

### Força de Senha

- [ ] **PWD-01**: Força avaliada como forte quando: ≥12 caracteres, ≥1 maiúscula, ≥1 minúscula, ≥1 dígito, ≥1 caractere especial; aviso exibido se fraca mas operação não bloqueada

### Integração Contínua

- [ ] **CI-01**: CI obrigatório executando build + lint + suíte completa de testes em todo push; mudanças que quebrem build ou testes não são aceitas
- [ ] **CI-02**: Matriz de CI cobrindo Windows, macOS e Linux

## v2 Requirements

### Duress Password (Senha Falsa de Coação)

- **DURESS-01**: Usuário pode configurar duress password (diferente da senha mestra; confirmação dupla)
- **DURESS-02**: Abre "versão restrita" do cofre quando duress password é usada; usuário não é informado qual senha foi aceita
- **DURESS-03**: Usuário pode alterar ou remover duress password
- **DURESS-04**: Usuário pode configurar quais segredos/pastas são visíveis na versão restrita

### Exibição Parcial de Campos Sensíveis

- **PARTIAL-01**: Usuário pode configurar exibição parcial de campo sensível (ex: últimos 4 dígitos de cartão de crédito)

### Gerador de Senhas

- **GEN-01**: A ser especificado em v2

### Compartilhamento via QR Code

- **QR-01**: Renderizar QR code na TUI com valor de um campo para transferência offline para outro dispositivo

### Relatório de Saúde do Cofre

- **HEALTH-01**: Análise local de senhas fracas, reutilizadas ou antigas

### Tags

- **TAGS-01**: Categorização de segredos por tags com filtragem

### Histórico de Versões

- **HIST-01**: Registro de versões anteriores de um segredo com visualização e restauração

### Recuperação de Artefatos Órfãos

- **ORPHAN-01**: Detecção e oferta de recuperação de `.abditum.tmp`/`.abditum.bak2` ao abrir cofre

## Out of Scope

| Feature | Reason |
|---------|--------|
| TOTP (autenticação de dois fatores) | Excluído permanentemente — fora do foco de credenciais estáticas |
| Backup automático | Responsabilidade do usuário — a app não gerencia cópias de segurança |
| Recuperação de dados corrompidos | Criptografia autenticada (GCM) não permite recuperação parcial — design intencional |
| Keyfile / Token de hardware (YubiKey) | Excluído permanentemente — fora do modelo de segurança atual |
| Armazenamento em nuvem / sync | Contraria filosofia offline e portátil — excluído permanentemente |
| Múltiplos cofres simultâneos | Invariante de design — só um cofre ativo por vez |
| App mobile ou web | TUI portátil é o produto — excluído permanentemente |
| Browser extension | Fora do escopo offline e single-binary |
| Compartilhamento de cofre / team | Produto individual — sem recursos de colaboração |

## Traceability

*(Preenchida durante criação do roadmap)*

| Requirement | Phase | Status |
|-------------|-------|--------|
| CRYPTO-01 | — | Pending |
| CRYPTO-02 | — | Pending |
| CRYPTO-03 | — | Pending |
| CRYPTO-04 | — | Pending |
| CRYPTO-05 | — | Pending |
| CRYPTO-06 | — | Pending |
| VAULT-01 | — | Pending |
| VAULT-02 | — | Pending |
| VAULT-03 | — | Pending |
| VAULT-04 | — | Pending |
| VAULT-05 | — | Pending |
| VAULT-06 | — | Pending |
| VAULT-07 | — | Pending |
| VAULT-08 | — | Pending |
| VAULT-09 | — | Pending |
| VAULT-10 | — | Pending |
| VAULT-11 | — | Pending |
| VAULT-12 | — | Pending |
| VAULT-13 | — | Pending |
| VAULT-14 | — | Pending |
| VAULT-15 | — | Pending |
| VAULT-16 | — | Pending |
| VAULT-17 | — | Pending |
| ATOMIC-01 | — | Pending |
| ATOMIC-02 | — | Pending |
| ATOMIC-03 | — | Pending |
| ATOMIC-04 | — | Pending |
| COMPAT-01 | — | Pending |
| COMPAT-02 | — | Pending |
| COMPAT-03 | — | Pending |
| QUERY-01 | — | Pending |
| QUERY-02 | — | Pending |
| QUERY-03 | — | Pending |
| QUERY-04 | — | Pending |
| QUERY-05 | — | Pending |
| SEC-01 | — | Pending |
| SEC-02 | — | Pending |
| SEC-03 | — | Pending |
| SEC-04 | — | Pending |
| SEC-05 | — | Pending |
| SEC-06 | — | Pending |
| SEC-07 | — | Pending |
| SEC-08 | — | Pending |
| SEC-09 | — | Pending |
| FOLDER-01 | — | Pending |
| FOLDER-02 | — | Pending |
| FOLDER-03 | — | Pending |
| FOLDER-04 | — | Pending |
| FOLDER-05 | — | Pending |
| TPL-01 | — | Pending |
| TPL-02 | — | Pending |
| TPL-03 | — | Pending |
| TPL-04 | — | Pending |
| TPL-05 | — | Pending |
| PWD-01 | — | Pending |
| CI-01 | — | Pending |
| CI-02 | — | Pending |

**Coverage:**
- v1 requirements: 55 total
- Mapped to phases: 0 (pending roadmap)
- Unmapped: 55 ⚠️

---
*Requirements defined: 2026-03-27*
*Last updated: 2026-03-27 after initialization*
