# Requirements: Abditum

**Defined:** April 9, 2026
**Core Value:** O usuário possui e controla completamente seus dados — o cofre é um arquivo que carrega consigo, protegido por criptografia forte, acessível apenas com a senha mestra, funcionando offline em qualquer sistema.

## v1 Requirements

### Ciclo de Vida do Cofre

- [ ] **VAULT-01**: Usuário pode criar novo cofre em um arquivo `.abditum` com senha mestra, com confirmação dupla e aviso de senha fraca (não bloqueante)
- [ ] **VAULT-02**: Ao criar o cofre, estrutura padrão é criada automaticamente: Pasta Geral com subpastas "Sites e Apps" e "Financeiro", e modelos padrão (Login, Cartão de Crédito, Chave de API)
- [ ] **VAULT-03**: Usuário pode abrir cofre existente a partir de um arquivo `.abditum` com senha mestra
- [ ] **VAULT-04**: Abertura valida o arquivo contra corrupção e senha incorreta — erros de autenticação permitem nova tentativa; erros de integridade encerram sem opção de recuperação
- [ ] **VAULT-05**: Usuário pode salvar cofre no arquivo atual sem re-informar senha
- [ ] **VAULT-06**: Usuário pode salvar cofre em outro arquivo (salvar como), que passa a ser o arquivo de trabalho ativo
- [ ] **VAULT-07**: Usuário pode descartar alterações não salvas e recarregar o cofre do arquivo
- [ ] **VAULT-08**: Usuário pode alterar a senha mestra — o cofre é salvo imediatamente com a nova senha; operação irrevogável
- [ ] **VAULT-09**: Usuário pode bloquear o cofre manualmente; ao bloquear, senha é removida da memória, clipboard limpa, terminal limpo (clear screen)
- [ ] **VAULT-10**: Cofre bloqueia automaticamente após tempo de inatividade configurável (padrão 5 min); qualquer interação reseta o timer
- [ ] **VAULT-11**: Usuário pode acionar bloqueio emergencial com tela disfarçada; atalho `⌃!⇧Q`; descarta alterações; sem confirmação
- [ ] **VAULT-12**: Ao sair da aplicação com alterações pendentes, usuário é perguntado: Salvar e Sair / Descartar e Sair / Cancelar
- [ ] **VAULT-13**: Usuário pode exportar cofre para arquivo de intercâmbio não criptografado, com aviso de risco e confirmação explícita
- [ ] **VAULT-14**: Usuário pode importar cofre de arquivo de intercâmbio, com mesclagem de pastas e sobrescrita de segredos/modelos conflitantes
- [ ] **VAULT-15**: Usuário pode configurar o cofre: tempo de bloqueio automático, tempo de ocultação de campo sensível, tempo de limpeza de clipboard, tema visual
- [ ] **VAULT-16**: Configurações são persistidas no payload criptografado do cofre (não em arquivo externo)

### Formato e Persistência

- [ ] **FILE-01**: Arquivo `.abditum` segue formato binário: cabeçalho fixo 49 bytes (`magic ABDT` + versão uint8 + salt 32B + nonce 12B) + payload JSON UTF-8 criptografado + tag GCM 16 bytes
- [ ] **FILE-02**: Criptografia: AES-256-GCM com chave derivada via Argon2id (m=256MiB, t=3, p=4); salt e nonce gerados via CSPRNG (`crypto/rand`)
- [ ] **FILE-03**: Cabeçalho completo autenticado como AAD do AES-256-GCM; salt substituído apenas ao alterar senha mestra; nonce regenerado a cada salvamento
- [ ] **FILE-04**: Gravação atômica: dados escritos em `.abditum.tmp` → rename substituindo arquivo original; em falha, `.abditum.tmp` apagado imediatamente
- [ ] **FILE-05**: Backup automático a cada sobrescrita: arquivo anterior preservado como `.abditum.bak` com rotação via `.abditum.bak2`; falha durante escrita preserva backup e exibe mensagem de erro com caminho do backup
- [ ] **FILE-06**: Criação de novo cofre em caminho vazio e "salvar como" em caminho vazio gravam diretamente no destino (sem `.abditum.tmp`)
- [ ] **FILE-07**: Aplicação detecta modificação externa do arquivo antes de sobrescrever e oferece: Sobrescrever / Salvar como novo arquivo / Cancelar
- [ ] **FILE-08**: Aplicação detecta acesso concorrente ao arquivo (outro processo usando o cofre) e falha com mensagem informativa
- [ ] **FILE-09**: Compatibilidade retroativa: abre qualquer versão de formato suportada; migra payload em memória; salva sempre no formato atual. Versão superior ao suportado é rejeitada com erro claro

### Segurança em Sessão

- [ ] **SEC-01**: Senha mestra solicitada uma única vez por sessão; re-uso para todas as operações de criptografia; sem re-solicitação em salvamento ou descarte
- [ ] **SEC-02**: Senha mestra e buffers sensíveis removidos da memória (zeros) ao bloquear ou encerrar a aplicação
- [ ] **SEC-03**: Terminal limpo (clear screen) ao bloquear ou encerrar, evitando dados visíveis no buffer do terminal
- [ ] **SEC-04**: Dados sensíveis em memória protegida (mlock/VirtualLock) quando disponível no SO; ausência não impede operação
- [ ] **SEC-05**: Valor da clipboard removido ao bloquear, encerrar ou após tempo configurável (padrão 30s); limpeza depende de suporte do SO
- [ ] **SEC-06**: Nenhum log (stdout/stderr) contém caminhos de cofre, nomes de segredos ou valores de campos

### Navegação e Consulta

- [ ] **NAV-01**: Interface exibe árvore hierárquica do cofre com pastas e segredos; cada pasta exibe contagem total de segredos ativos (recursivo)
- [ ] **NAV-02**: Pasta virtual "Favoritos" exibida como nó irmão da Pasta Geral (acima dela); somente leitura; lista segredos favoritos em profundidade seguindo ordem do cofre
- [ ] **NAV-03**: Usuário pode buscar segredos por nome, nome de campo, valor de campo comum e observação (substring, case-insensitive, sem distinção de acentuação); campos sensíveis nunca participam da busca por valor; segredos marcados para exclusão excluídos dos resultados
- [ ] **NAV-04**: Segredo exibe nome, todos os campos (com indicação de sensível/comum) e observação
- [ ] **NAV-05**: Indicadores de estado de sessão na listagem: ✦ adicionado, ✎ modificado, ✗ excluído (segredos sem alteração não exibem indicador)
- [ ] **NAV-06**: Dirty state global do cofre sempre visível no cabeçalho da interface

### Campos Sensíveis e Clipboard

- [ ] **FIELD-01**: Usuário pode exibir temporariamente o valor de um campo sensível; valor oculto automaticamente após tempo configurável (padrão 15s)
- [ ] **FIELD-02**: Usuário pode copiar valor de qualquer campo para a clipboard; valor removido automaticamente após tempo configurável (padrão 30s) ou ao bloquear/encerrar

### Gerenciamento de Segredos

- [ ] **SECRET-01**: Usuário pode criar segredo a partir de modelo existente ou sem campos estruturados (apenas com Observação automática)
- [ ] **SECRET-02**: Usuário pode escolher a pasta de destino no momento da criação do segredo
- [ ] **SECRET-03**: Usuário pode duplicar segredo; duplicata recebe nome único automático (ex: "Gmail (1)") e é posicionada imediatamente após o original; histórico de modelo preservado
- [ ] **SECRET-04**: Usuário pode editar nome, valores de campos e observação de um segredo
- [ ] **SECRET-05**: Usuário pode editar estrutura de segredo: adicionar campo (nome + tipo), renomear campo, reordenar campos, excluir campo; não é permitido alterar o tipo de um campo existente
- [ ] **SECRET-06**: Observação é campo automático presente em todo segredo: não pode ser renomeada, excluída ou movida; sempre na última posição; campo comum (sempre visível)
- [ ] **SECRET-07**: Usuário pode favoritar e desfavoritar segredo
- [ ] **SECRET-08**: Usuário pode marcar segredo para exclusão; segredo permanece na lista com indicador ✗ e é excluído permanentemente apenas ao salvar com sucesso; restaurável antes do salvamento
- [ ] **SECRET-09**: Usuário pode desmarcar exclusão de segredo, restaurando o estado anterior
- [ ] **SECRET-10**: Usuário pode mover segredo para outra pasta
- [ ] **SECRET-11**: Usuário pode reordenar segredo dentro da mesma pasta; ordem final persistida ao salvar

### Gerenciamento de Pastas

- [ ] **FOLDER-01**: Usuário pode criar pasta dentro de outra pasta; nome deve ser único entre irmãos na mesma pasta pai
- [ ] **FOLDER-02**: Usuário pode renomear pasta; novo nome deve ser único entre irmãos na mesma pasta pai
- [ ] **FOLDER-03**: Usuário pode mover pasta para outra; sistema impede movimentos que criariam ciclos; impede mover para pasta que já contém irmão com o mesmo nome
- [ ] **FOLDER-04**: Usuário pode reordenar pasta dentro da mesma pasta pai; ordem final persistida ao salvar
- [ ] **FOLDER-05**: Usuário pode excluir pasta; conteúdo (segredos e subpastas) promovido para pasta pai; segredos com nome duplicado na pasta pai renomeados automaticamente com sufixo numérico; usuário avisado sobre renomeações
- [ ] **FOLDER-06**: Pasta Geral não pode ser renomeada, movida ou excluída; pode estar vazia
- [ ] **FOLDER-07**: Regras de hierarquia: sem dois segredos com mesmo nome na mesma pasta; sem duas subpastas com mesmo nome na mesma pasta pai; ciclos proibidos; todas as pastas navegáveis a partir da Pasta Geral

### Gerenciamento de Modelos

- [ ] **TMPL-01**: Usuário pode criar modelo de segredo com campos personalizados (nome + tipo)
- [ ] **TMPL-02**: Usuário pode renomear modelo; nome deve ser único globalmente
- [ ] **TMPL-03**: Usuário pode editar estrutura do modelo: adicionar campo, renomear campo, alterar tipo de campo, reordenar campos, excluir campo; campos do modelo permitem nomes duplicados entre si
- [ ] **TMPL-04**: Usuário pode excluir modelo de segredo
- [ ] **TMPL-05**: Usuário pode criar modelo a partir de segredo existente; campo Observação do segredo não é copiado para o modelo
- [ ] **TMPL-06**: Segredo criado a partir de modelo não mantém vínculo com ele — alterações no modelo não afetam segredos existentes; nome do modelo registrado apenas como histórico
- [ ] **TMPL-07**: Modelos exibidos em ordem alfabética; não são reordenáveis manualmente

### Interface TUI e Design System

- [ ] **TUI-01**: Interface TUI em Go com Bubbletea/Lipgloss, baseada em design system definido: 4 zonas verticais (Cabeçalho 2L | Área de trabalho | Barra de mensagens 1L | Barra de comandos 1L)
- [ ] **TUI-02**: Painel esquerdo (árvore/lista ≈35%) e painel direito (detalhe ≈65%); separador `│`; conector `<╡` no item selecionado
- [ ] **TUI-03**: Terminal mínimo 80×24; degradação sem crash abaixo do mínimo; truncamento com `…`
- [ ] **TUI-04**: Todos os símbolos do inventário BMP-only (U+0000–U+FFFF); sem emojis, sem Nerd Fonts
- [ ] **TUI-05**: NO_COLOR compliant: todo estado crítico usa ≥2 camadas de comunicação (cor + símbolo ou atributo tipográfico)
- [ ] **TUI-06**: Temas visuais discretos e reativos; identificador do tema gravado no payload criptografado do cofre
- [ ] **TUI-07**: Mapa de teclas conforme design system: `Enter` avança/confirma; `Esc` retrocede/cancela; `Tab`/`⇧Tab` alterna painéis ou campos; atalhos globais `F1` (ajuda), `F12` (tema), `⌃Q` (sair), `⌃!⇧Q` (bloqueio emergencial)
- [ ] **TUI-08**: Barra de mensagens com 7 tipos: Sucesso ✓, Informação ℹ, Alerta ⚠, Erro ✕ (bold), Ocupado spinner `◐◓◑◒`, Dica de campo `•` (italic), Dica de uso `•` (italic)
- [ ] **TUI-09**: Diálogos e modais centralizados com bordas arredondadas `╭╮╰╯`; ações na borda inferior; severidade governa visual (borda + símbolo + cor da ação default)
- [ ] **TUI-10**: Scrollbar em diálogos com conteúdo que excede viewport: setas `↑↓` e thumb `■` na borda direita

### Distribuição e Portabilidade

- [ ] **DIST-01**: Aplicação distribuída como binário único executável cross-platform (Windows, macOS, Linux 64-bit); sem instalação; sem dependências externas em runtime
- [ ] **DIST-02**: Zero arquivos de estado, config ou dados fora do arquivo do cofre (exceto artefatos transitórios `.abditum.tmp`, `.abditum.bak`, `.abditum.bak2`)

## v2 Requirements

*Planejados para versão futura. Não fazem parte do roadmap atual.*

- **V2-01**: Exibição parcial de campos sensíveis configurável por campo (últimos N / primeiros N caracteres, padrão mascarado)
- **V2-02**: Gerador de senhas aleatórias integrado
- **V2-03**: Compartilhamento via QR Code renderizado na TUI (blocos ASCII/Unicode) para transferência offline
- **V2-04**: Relatório de saúde do cofre: senhas fracas, reutilizadas ou antigas
- **V2-05**: Categorização de segredos por tags com filtragem
- **V2-06**: Histórico de versões de segredos com visualização e restauração
- **V2-07**: Recuperação de artefatos órfãos (`.abditum.tmp`, `.abditum.bak2`) ao abrir cofre

## Out of Scope

| Feature | Reason |
|---------|--------|
| Senha Falsa de Coação (Duress Password) | Complexidade de implementação segura supera benefício na v1; avaliação futura |
| TOTP / Two-Factor Authentication | Excluído permanentemente |
| Backup gerenciado pela aplicação | Responsabilidade do usuário |
| Recuperação de dados corrompidos | AES-256-GCM não permite recuperação parcial; sem mecanismo de reparo |
| Autenticação por Keyfile ou token de hardware (YubiKey) | Excluído permanentemente |
| Armazenamento em nuvem | Contraria filosofia offline e portátil; excluído permanentemente |
| Múltiplos cofres abertos simultaneamente | Invariante de design; excluído permanentemente |
| App mobile ou web | TUI portátil por design; excluído permanentemente |
| Modo somente leitura com tratamento especial | Falha informativa ao tentar salvar; sem estado especial |

## Traceability

| Requirement | Phase | Status |
|-------------|-------|--------|
| VAULT-01 – VAULT-16 | TBD | Pending |
| FILE-01 – FILE-09 | TBD | Pending |
| SEC-01 – SEC-06 | TBD | Pending |
| NAV-01 – NAV-06 | TBD | Pending |
| FIELD-01 – FIELD-02 | TBD | Pending |
| SECRET-01 – SECRET-11 | TBD | Pending |
| FOLDER-01 – FOLDER-07 | TBD | Pending |
| TMPL-01 – TMPL-07 | TBD | Pending |
| TUI-01 – TUI-10 | TBD | Pending |
| DIST-01 – DIST-02 | TBD | Pending |

**Coverage:**
- v1 requirements: 63 total
- Mapped to phases: 0 (pendente — mapeamento pelo roadmapper)
- Unmapped: 63 ⚠️

---
*Requirements defined: April 9, 2026*
*Last updated: April 9, 2026 after initialization*
