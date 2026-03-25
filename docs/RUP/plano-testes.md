# Plano de Testes — Abditum

| Item            | Detalhe                        |
|-----------------|--------------------------------|
| Projeto         | Abditum                       |
| Versão          | 1.0                           |
| Data            | 2026-03-25                     |
| Status          | Aprovado                       |

---

## 1. Introdução

### 1.1 Objetivo
Definir a estratégia, os tipos, os níveis e os critérios de teste para o Abditum v1, garantindo cobertura adequada de funcionalidade, segurança, confiabilidade e portabilidade.

### 1.2 Escopo
Cobre todos os requisitos funcionais (RF-01 a RF-41), requisitos não funcionais (RNF-01 a RNF-30) e regras de negócio (RN-01 a RN-28) definidos na SRS.

### 1.3 Referências
- SRS — `docs/RUP/srs.md`
- Casos de Uso — `docs/RUP/casos-de-uso.md`
- Lista de Riscos — `docs/RUP/riscos.md`
- Documento descritivo — `descricao.md`

---

## 2. Estratégia de Testes

### 2.1 Níveis de Teste

| Nível         | Descrição                                                                                     | Ferramenta          |
|---------------|-----------------------------------------------------------------------------------------------|---------------------|
| Unitário      | Testa funções e métodos isolados das camadas de domínio, criptografia e armazenamento          | `go test`           |
| Integração    | Testa interação entre camadas (domínio ↔ criptografia ↔ armazenamento) em fluxos completos   | `go test`           |
| Componente    | Testa componentes TUI isolados com entradas simuladas                                          | `teatest/v2`        |
| Visual        | Compara snapshots de telas renderizadas contra golden files de referência                      | `teatest/v2`        |
| Sistema/E2E   | Testa fluxos completos de ponta a ponta simulando interações do usuário                        | `teatest/v2`        |

### 2.2 Tipos de Teste

| Tipo             | Objetivo                                                                         |
|------------------|----------------------------------------------------------------------------------|
| Funcional        | Verificar comportamento correto dos requisitos funcionais                         |
| Segurança        | Validar proteção criptográfica, zeragem de buffers, ausência de logs sensíveis   |
| Confiabilidade   | Validar salvamento atômico, rotação de backup, rollback em caso de falha         |
| Compatibilidade  | Validar migração de formatos anteriores e rejeição de versões superiores          |
| Portabilidade    | Executar testes em Windows, macOS e Linux via CI                                 |
| Regressão        | Garantir não quebra ao adicionar funcionalidades ou alterar código existente      |

---

## 3. Itens de Teste

### 3.1 Serviço de Criptografia

| ID    | Caso de Teste                                                                      | Tipo       | Nível     | Requisitos        |
|-------|------------------------------------------------------------------------------------|------------|-----------|-------------------|
| TC-C01 | Criptografar e descriptografar payload com sucesso                                | Funcional  | Unitário  | RNF-01, RNF-24    |
| TC-C02 | Descriptografar com senha incorreta retorna erro (sem distinguir de corrupção)     | Funcional  | Unitário  | RNF-01            |
| TC-C03 | Descriptografar payload corrompido retorna erro de integridade                     | Segurança  | Unitário  | RNF-01, RNF-03    |
| TC-C04 | Validar que o cabeçalho é autenticado como AAD do GCM                              | Segurança  | Unitário  | RNF-03            |
| TC-C05 | Alterar qualquer byte do cabeçalho invalida a descriptografia                      | Segurança  | Unitário  | RNF-03            |
| TC-C06 | Derivar chave com Argon2id usando parâmetros v1 (256 MiB, 3 iterações, 4 threads) | Funcional  | Unitário  | RNF-02, RNF-27-30 |
| TC-C07 | Verificar que salt e nonce são diferentes a cada operação de criação/salvamento    | Segurança  | Unitário  | RNF-04, RNF-05    |
| TC-C08 | Validar assinatura mágica `ABDT` na leitura                                        | Funcional  | Unitário  | RNF-22            |
| TC-C09 | Rejeitar arquivo sem assinatura mágica `ABDT`                                      | Funcional  | Unitário  | RNF-22            |
| TC-C10 | Rejeitar arquivo com versão do formato superior à suportada                        | Funcional  | Unitário  | RNF-19            |

### 3.2 Serviço de Armazenamento

| ID    | Caso de Teste                                                                      | Tipo           | Nível     | Requisitos       |
|-------|------------------------------------------------------------------------------------|----------------|-----------|------------------|
| TC-A01 | Salvar cofre com sucesso via .tmp + rename atômico                                | Funcional      | Unitário  | RNF-11           |
| TC-A02 | Gerar backup .bak antes de substituir arquivo existente                            | Confiabilidade | Unitário  | RNF-12           |
| TC-A03 | Rotação .bak → .bak2 quando .bak já existe                                        | Confiabilidade | Unitário  | RNF-13           |
| TC-A04 | Remover .bak2 após salvamento bem-sucedido                                         | Confiabilidade | Unitário  | RNF-13           |
| TC-A05 | Restaurar .bak2 → .bak em caso de falha no salvamento                              | Confiabilidade | Unitário  | RNF-13           |
| TC-A06 | Remover .tmp em caso de falha na gravação                                          | Confiabilidade | Unitário  | RNF-11           |
| TC-A07 | Carregar cofre existente com sucesso                                               | Funcional      | Unitário  | RF-02            |
| TC-A08 | Falha ao salvar em caminho sem permissão de escrita                                | Confiabilidade | Unitário  | RNF-14           |
| TC-A09 | Detectar arquivo bloqueado por outro processo                                      | Confiabilidade | Unitário  | RNF-15           |
| TC-A10 | Salvar Como (novo caminho) sem .tmp, com gravação direta                           | Funcional      | Unitário  | RF-04            |

### 3.3 Domínio e Regras de Negócio

| ID    | Caso de Teste                                                                      | Tipo       | Nível     | Requisitos            |
|-------|------------------------------------------------------------------------------------|------------|-----------|------------------------|
| TC-D01 | Criar cofre com estrutura inicial (pastas e modelos pré-definidos)                | Funcional  | Unitário  | RF-01, RN-22           |
| TC-D02 | Criar segredo a partir de modelo (snapshot, sem vínculo)                           | Funcional  | Unitário  | RF-16, RN-10           |
| TC-D03 | Criar segredo vazio (sem campos adicionais)                                       | Funcional  | Unitário  | RF-16                  |
| TC-D04 | Duplicar segredo com nova identidade e nome sufixado                              | Funcional  | Unitário  | RF-18, RN-26           |
| TC-D05 | Favoritar e desfavoritar segredo                                                  | Funcional  | Unitário  | RF-19                  |
| TC-D06 | Editar dados do segredo (nome, valores, observação)                               | Funcional  | Unitário  | RF-20                  |
| TC-D07 | Editar estrutura do segredo (incluir, renomear, excluir, reordenar campos)        | Funcional  | Unitário  | RF-21                  |
| TC-D08 | Impedir alteração de tipo de campo existente                                      | Funcional  | Unitário  | RN-11                  |
| TC-D09 | Exclusão reversível move segredo para Lixeira                                     | Funcional  | Unitário  | RF-22, RN-04, RN-23    |
| TC-D10 | Restaurar segredo da Lixeira ao local e estado originais                          | Funcional  | Unitário  | RF-23, RN-24           |
| TC-D11 | Restaurar segredo cuja pasta de origem foi excluída → vai para raiz               | Funcional  | Unitário  | RN-24                  |
| TC-D12 | Segredo em Lixeira não pode ser editado                                           | Funcional  | Unitário  | RN-25                  |
| TC-D13 | Salvar cofre esvazia Lixeira permanentemente                                      | Funcional  | Unitário  | RN-05                  |
| TC-D14 | Mover segredo entre pastas preserva identidade e conteúdo                         | Funcional  | Unitário  | RF-24, RN-02           |
| TC-D15 | Reordenar segredo altera posição sem modificar conteúdo                           | Funcional  | Unitário  | RF-25                  |
| TC-D16 | Busca por nome, nome de campo, valor texto e observação                           | Funcional  | Unitário  | RF-26, RN-07, RN-08    |
| TC-D17 | Busca NÃO inclui valores de campos tipo texto sensível                            | Segurança  | Unitário  | RN-07                  |
| TC-D18 | Criar pasta na raiz e dentro de outra pasta                                       | Funcional  | Unitário  | RF-27                  |
| TC-D19 | Renomear pasta sem alterar identidade ou conteúdo                                 | Funcional  | Unitário  | RF-28                  |
| TC-D20 | Mover pasta preserva hierarquia interna recursivamente                            | Funcional  | Unitário  | RF-29, RN-03           |
| TC-D21 | Excluir pasta promove filhos ao nível pai                                         | Funcional  | Unitário  | RF-31, RN-06           |
| TC-D22 | Criar modelo de segredo com campos personalizados                                 | Funcional  | Unitário  | RF-32                  |
| TC-D23 | Editar modelo não afeta segredos já criados                                       | Funcional  | Unitário  | RF-33, RN-12           |
| TC-D24 | Remover modelo não afeta segredos já criados                                      | Funcional  | Unitário  | RF-34, RN-10           |
| TC-D25 | Criar modelo a partir de segredo copia estrutura sem vínculo                      | Funcional  | Unitário  | RF-35, RN-10           |
| TC-D26 | Transição de estado Cofre Salvo → Cofre Modificado em qualquer mutação            | Funcional  | Unitário  | RN-13                  |
| TC-D27 | Nomes repetidos permitidos para segredos, pastas e modelos                        | Funcional  | Unitário  | RN-09                  |
| TC-D28 | Configurações embutidas no cofre: alterar e persistir tempos configuráveis        | Funcional  | Unitário  | RF-11, RN-27           |

### 3.4 Importação e Exportação

| ID    | Caso de Teste                                                                      | Tipo       | Nível       | Requisitos             |
|-------|------------------------------------------------------------------------------------|------------|-------------|------------------------|
| TC-IE01 | Exportar cofre para texto claro com conteúdo correto                             | Funcional  | Integração  | RF-09                  |
| TC-IE02 | Importar cofre sem conflitos                                                     | Funcional  | Integração  | RF-10                  |
| TC-IE03 | Importar com conflito de identidade de segredo → nova identidade                 | Funcional  | Integração  | RF-10, RN-17           |
| TC-IE04 | Importar com conflito de nome de segredo → sufixo incremental                    | Funcional  | Integração  | RF-10, RN-18           |
| TC-IE05 | Importar com conflito de identidade de pasta → merge silencioso                  | Funcional  | Integração  | RF-10, RN-16           |
| TC-IE06 | Importar com conflito de identidade de modelo → sobrescrita silenciosa           | Funcional  | Integração  | RF-10, RN-19           |

### 3.5 Compatibilidade de Formato

| ID    | Caso de Teste                                                                      | Tipo            | Nível       | Requisitos        |
|-------|------------------------------------------------------------------------------------|-----------------|-------------|-------------------|
| TC-F01 | Abrir cofre de versão anterior e migrar em memória                                | Compatibilidade | Integração  | RNF-16, RNF-17    |
| TC-F02 | Salvar cofre migrado no formato da versão atual                                   | Compatibilidade | Integração  | RNF-17            |
| TC-F03 | Selecionar perfil Argon2id correto para versão histórica                          | Compatibilidade | Unitário    | RNF-18            |
| TC-F04 | Rejeitar cofre com versão do formato superior                                     | Compatibilidade | Unitário    | RNF-19            |

### 3.6 Segurança Operacional

| ID    | Caso de Teste                                                                      | Tipo       | Nível       | Requisitos       |
|-------|------------------------------------------------------------------------------------|------------|-------------|------------------|
| TC-S01 | Limpeza de área de transferência após tempo configurado                           | Segurança  | Integração  | RF-37            |
| TC-S02 | Limpeza de área de transferência ao bloquear cofre                                | Segurança  | Integração  | RF-38            |
| TC-S03 | Bloqueio por inatividade após tempo configurado                                   | Funcional  | Integração  | RF-08, RN-14     |
| TC-S04 | Alerta de bloqueio iminente aos 75% do tempo de inatividade                       | Funcional  | Integração  | RN-20            |
| TC-S05 | Atividade do usuário reinicia cronômetro de inatividade                           | Funcional  | Integração  | RN-21            |
| TC-S06 | Bloqueio descarta alterações não salvas silenciosamente                           | Funcional  | Integração  | RN-14, RN-15     |
| TC-S07 | Ausência de dados sensíveis em saída padrão e saída de erro                      | Segurança  | Unitário    | RNF-07           |

### 3.7 Interface TUI (Golden Files e Comandos)

| ID    | Caso de Teste                                                                      | Tipo       | Nível       | Requisitos       |
|-------|------------------------------------------------------------------------------------|------------|-------------|------------------|
| TC-T01 | Golden file: tela inicial (welcome) em 80×24                                     | Visual     | Componente  | —                |
| TC-T02 | Golden file: tela de cofre ativo com hierarquia e segredo em 80×24               | Visual     | Componente  | —                |
| TC-T03 | Golden file: tela de criação de segredo em 80×24                                 | Visual     | Componente  | —                |
| TC-T04 | Golden file: tela de edição de segredo (padrão e avançada) em 80×24              | Visual     | Componente  | —                |
| TC-T05 | Golden file: diálogo de confirmação bloqueante em 80×24                           | Visual     | Componente  | —                |
| TC-T06 | Golden file: file picker em 80×24                                                | Visual     | Componente  | —                |
| TC-T07 | Golden file: tela de busca com resultados filtrados em 80×24                      | Visual     | Componente  | —                |
| TC-T08 | Golden file: mensagem de terminal muito pequeno                                  | Visual     | Componente  | —                |
| TC-T09 | Comando: navegação na árvore (setas, expandir/colapsar)                          | Funcional  | Componente  | RF-12            |
| TC-T10 | Comando: alternância de painéis com Tab                                          | Funcional  | Componente  | —                |
| TC-T11 | Comando: toggle de campo sensível                                                | Funcional  | Componente  | RF-14            |
| TC-T12 | Comando: cópia de campo para área de transferência                               | Funcional  | Componente  | RF-36            |

### 3.8 Testes de Integração E2E

| ID    | Caso de Teste                                                                      | Tipo       | Nível  | Requisitos                    |
|-------|------------------------------------------------------------------------------------|------------|--------|-------------------------------|
| TC-E01 | Fluxo completo: criar cofre → criar segredo → salvar → fechar → reabrir → verificar dados | Funcional | E2E | RF-01, RF-02, RF-03, RF-16  |
| TC-E02 | Fluxo completo: abrir cofre → editar segredo → salvar → verificar alterações     | Funcional  | E2E    | RF-02, RF-03, RF-20          |
| TC-E03 | Fluxo completo: criar cofre → excluir segredo → restaurar → salvar               | Funcional  | E2E    | RF-22, RF-23, RF-03          |
| TC-E04 | Fluxo completo: criar cofre → exportar → importar em outro cofre → verificar     | Funcional  | E2E    | RF-09, RF-10                 |
| TC-E05 | Fluxo completo: criar cofre → alterar senha mestra → fechar → reabrir com nova senha | Funcional | E2E | RF-06, RF-02                 |
| TC-E06 | Fluxo completo: criar cofre → modificar → sair com confirmação de salvamento     | Funcional  | E2E    | UC-11                        |
| TC-E07 | Fluxo completo: criar cofre → criar modelo → criar segredo do modelo → verificar | Funcional  | E2E    | RF-32, RF-16                 |
| TC-E08 | Fluxo completo: abrir cofre antigo → migrar formato → salvar → reabrir           | Funcional  | E2E    | RNF-16, RNF-17               |

---

## 4. Critérios de Entrada e Saída

### 4.1 Critérios de Entrada
- Código compilando sem erros em Go para as 3 plataformas alvo
- Ambiente de CI configurado e funcional
- Casos de teste implementados para a iteração em questão

### 4.2 Critérios de Saída
- Todos os testes unitários passando (100%)
- Todos os testes de integração passando (100%)
- Golden files visuais aprovados sem diferenças
- Testes executados com sucesso em Windows, macOS e Linux
- Nenhum defeito de severidade Crítica ou Alta em aberto
- Cobertura de código ≥ 80% nas camadas de criptografia e armazenamento

---

## 5. Ambiente de Testes

| Plataforma | Ambiente                       |
|------------|--------------------------------|
| Windows    | Windows 10/11, 64 bits         |
| macOS      | macOS 13+, ARM64 e x86_64      |
| Linux      | Ubuntu 22.04+, 64 bits         |
| CI         | GitHub Actions (3 plataformas) |
| Terminal   | 80×24 mínimo para golden files |

---

## 6. Riscos de Teste

| Risco                                                                 | Mitigação                                           |
|-----------------------------------------------------------------------|-----------------------------------------------------|
| Limpeza de clipboard não testável em CI headless                      | Testar manualmente; mock em unitários                |
| Golden files quebram em terminais com renderização diferente           | Padronizar dimensão e encoding no CI                 |
| Testes de inatividade são lentos por depender de timers               | Usar tempos curtos em testes; injeção de dependência |
| Rename atômico pode se comportar diferente em FS de CI vs. produção  | Testar em múltiplos FS quando possível               |
