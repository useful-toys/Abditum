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
Cobre todos os requisitos funcionais ([RF-01](srs.md#rf-01) a [RF-41](srs.md#rf-41)), requisitos não funcionais ([RNF-01](srs.md#rnf-01) a [RNF-30](srs.md#rnf-30)) e regras de negócio ([RN-01](regras-negocio.md#rn-01) a [RN-23](regras-negocio.md#rn-23)) definidos na SRS.

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
| <a id="tc-c01"></a>TC-C01 | Criptografar e descriptografar payload com sucesso                                | Funcional  | Unitário  | [RNF-01](srs.md#rnf-01), [RNF-24](srs.md#rnf-24)    |
| <a id="tc-c02"></a>TC-C02 | Descriptografar com senha incorreta retorna erro (sem distinguir de corrupção)     | Funcional  | Unitário  | [RNF-01](srs.md#rnf-01)            |
| <a id="tc-c03"></a>TC-C03 | Descriptografar payload corrompido retorna erro de integridade                     | Segurança  | Unitário  | [RNF-01](srs.md#rnf-01), [RNF-03](srs.md#rnf-03)    |
| <a id="tc-c04"></a>TC-C04 | Validar que o cabeçalho é autenticado como AAD do GCM                              | Segurança  | Unitário  | [RNF-03](srs.md#rnf-03)            |
| <a id="tc-c05"></a>TC-C05 | Alterar qualquer byte do cabeçalho invalida a descriptografia                      | Segurança  | Unitário  | [RNF-03](srs.md#rnf-03)            |
| <a id="tc-c06"></a>TC-C06 | Derivar chave com Argon2id usando parâmetros v1 (256 MiB, 3 iterações, 4 threads) | Funcional  | Unitário  | [RNF-02](srs.md#rnf-02), [RNF-27](srs.md#rnf-27) a [RNF-30](srs.md#rnf-30) |
| <a id="tc-c07"></a>TC-C07 | Verificar que salt e nonce são diferentes a cada operação de criação/salvamento    | Segurança  | Unitário  | [RNF-04](srs.md#rnf-04), [RNF-05](srs.md#rnf-05)    |
| <a id="tc-c08"></a>TC-C08 | Validar assinatura mágica `ABDT` na leitura                                        | Funcional  | Unitário  | [RNF-22](srs.md#rnf-22)            |
| <a id="tc-c09"></a>TC-C09 | Rejeitar arquivo sem assinatura mágica `ABDT`                                      | Funcional  | Unitário  | [RNF-22](srs.md#rnf-22)            |
| <a id="tc-c10"></a>TC-C10 | Rejeitar arquivo com versão do formato superior à suportada                        | Funcional  | Unitário  | [RNF-19](srs.md#rnf-19)            |

### 3.2 Serviço de Armazenamento

| ID    | Caso de Teste                                                                      | Tipo           | Nível     | Requisitos       |
|-------|------------------------------------------------------------------------------------|----------------|-----------|------------------|
| <a id="tc-a01"></a>TC-A01 | Salvar cofre com sucesso via .tmp + rename atômico                                | Funcional      | Unitário  | [RNF-11](srs.md#rnf-11)           |
| <a id="tc-a02"></a>TC-A02 | Gerar backup .bak antes de substituir arquivo existente                            | Confiabilidade | Unitário  | [RNF-12](srs.md#rnf-12)           |
| <a id="tc-a03"></a>TC-A03 | Rotação .bak → .bak2 quando .bak já existe                                        | Confiabilidade | Unitário  | [RNF-13](srs.md#rnf-13)           |
| <a id="tc-a04"></a>TC-A04 | Remover .bak2 após salvamento bem-sucedido                                         | Confiabilidade | Unitário  | [RNF-13](srs.md#rnf-13)           |
| <a id="tc-a05"></a>TC-A05 | Restaurar .bak2 → .bak em caso de falha no salvamento                              | Confiabilidade | Unitário  | [RNF-13](srs.md#rnf-13)           |
| <a id="tc-a06"></a>TC-A06 | Remover .tmp em caso de falha na gravação                                          | Confiabilidade | Unitário  | [RNF-11](srs.md#rnf-11)           |
| <a id="tc-a07"></a>TC-A07 | Carregar cofre existente com sucesso                                               | Funcional      | Unitário  | [RF-02](srs.md#rf-02)            |
| <a id="tc-a08"></a>TC-A08 | Falha ao salvar em caminho sem permissão de escrita                                | Confiabilidade | Unitário  | [RNF-14](srs.md#rnf-14)           |
| <a id="tc-a09"></a>TC-A09 | Detectar arquivo bloqueado por outro processo                                      | Confiabilidade | Unitário  | [RNF-15](srs.md#rnf-15)           |
| <a id="tc-a10"></a>TC-A10 | Salvar Como (novo caminho) sem .tmp, com gravação direta                           | Funcional      | Unitário  | [RF-04](srs.md#rf-04)            |

### 3.3 Domínio e Regras de Negócio

| ID    | Caso de Teste                                                                      | Tipo       | Nível     | Requisitos            |
|-------|------------------------------------------------------------------------------------|------------|-----------|------------------------|
| <a id="tc-d01"></a>TC-D01 | Criar cofre com estrutura inicial (pastas e modelos pré-definidos)                | Funcional  | Unitário  | [RF-01](srs.md#rf-01)           |
| <a id="tc-d02"></a>TC-D02 | Criar segredo a partir de modelo (snapshot, sem vínculo)                           | Funcional  | Unitário  | [RF-16](srs.md#rf-16), [RN-17](regras-negocio.md#rn-17)           |
| <a id="tc-d03"></a>TC-D03 | Criar segredo vazio (sem campos adicionais)                                       | Funcional  | Unitário  | [RF-16](srs.md#rf-16)                  |
| <a id="tc-d04"></a>TC-D04 | Duplicar segredo com nova identidade e nome sufixado                              | Funcional  | Unitário  | [RF-18](srs.md#rf-18), [RN-08](regras-negocio.md#rn-08)           |
| <a id="tc-d05"></a>TC-D05 | Favoritar e desfavoritar segredo                                                  | Funcional  | Unitário  | [RF-19](srs.md#rf-19)                  |
| <a id="tc-d06"></a>TC-D06 | Editar dados do segredo (nome, valores, observação)                               | Funcional  | Unitário  | [RF-20](srs.md#rf-20)                  |
| <a id="tc-d07"></a>TC-D07 | Editar estrutura do segredo (incluir, renomear, excluir, reordenar campos)        | Funcional  | Unitário  | [RF-21](srs.md#rf-21)                  |
| <a id="tc-d08"></a>TC-D08 | Impedir alteração de tipo de campo existente                                      | Funcional  | Unitário  | [RI-06](srs.md#ri-06)                  |
| <a id="tc-d09"></a>TC-D09 | Exclusão reversível move segredo para Lixeira                                     | Funcional  | Unitário  | [RF-22](srs.md#rf-22), [RN-12](regras-negocio.md#rn-12)    |
| <a id="tc-d10"></a>TC-D10 | Restaurar segredo da Lixeira ao local e estado originais                          | Funcional  | Unitário  | [RF-23](srs.md#rf-23), [RN-15](regras-negocio.md#rn-15)           |
| <a id="tc-d11"></a>TC-D11 | Restaurar segredo cuja pasta de origem foi excluída → vai para raiz               | Funcional  | Unitário  | [RN-15](regras-negocio.md#rn-15)                  |
| <a id="tc-d12"></a>TC-D12 | Segredo em Lixeira não pode ser editado                                           | Funcional  | Unitário  | [RN-14](regras-negocio.md#rn-14)                  |
| <a id="tc-d13"></a>TC-D13 | Salvar cofre esvazia Lixeira permanentemente                                      | Funcional  | Unitário  | [RN-13](regras-negocio.md#rn-13)                  |
| <a id="tc-d14"></a>TC-D14 | Mover segredo entre pastas preserva identidade e conteúdo                         | Funcional  | Unitário  | [RF-24](srs.md#rf-24), [RN-06](regras-negocio.md#rn-06)           |
| <a id="tc-d15"></a>TC-D15 | Reordenar segredo altera posição sem modificar conteúdo                           | Funcional  | Unitário  | [RF-25](srs.md#rf-25)                  |
| <a id="tc-d16"></a>TC-D16 | Busca por nome, nome de campo, valor texto e observação                           | Funcional  | Unitário  | [RF-26](srs.md#rf-26), [RN-10](regras-negocio.md#rn-10), [RN-11](regras-negocio.md#rn-11)    |
| <a id="tc-d17"></a>TC-D17 | Busca NÃO inclui valores de campos tipo texto sensível                            | Segurança  | Unitário  | [RN-10](regras-negocio.md#rn-10)                  |
| <a id="tc-d18"></a>TC-D18 | Criar pasta na raiz e dentro de outra pasta                                       | Funcional  | Unitário  | [RF-27](srs.md#rf-27)                  |
| <a id="tc-d19"></a>TC-D19 | Renomear pasta sem alterar identidade ou conteúdo                                 | Funcional  | Unitário  | [RF-28](srs.md#rf-28)                  |
| <a id="tc-d20"></a>TC-D20 | Mover pasta preserva hierarquia interna recursivamente                            | Funcional  | Unitário  | [RF-29](srs.md#rf-29), [RN-07](regras-negocio.md#rn-07)           |
| <a id="tc-d21"></a>TC-D21 | Excluir pasta promove filhos ao nível pai                                         | Funcional  | Unitário  | [RF-31](srs.md#rf-31), [RN-16](regras-negocio.md#rn-16)           |
| <a id="tc-d22"></a>TC-D22 | Criar modelo de segredo com campos personalizados                                 | Funcional  | Unitário  | [RF-32](srs.md#rf-32)                  |
| <a id="tc-d23"></a>TC-D23 | Editar modelo não afeta segredos já criados                                       | Funcional  | Unitário  | [RF-33](srs.md#rf-33), [RN-18](regras-negocio.md#rn-18)           |
| <a id="tc-d24"></a>TC-D24 | Remover modelo não afeta segredos já criados                                      | Funcional  | Unitário  | [RF-34](srs.md#rf-34), [RN-17](regras-negocio.md#rn-17)           |
| <a id="tc-d25"></a>TC-D25 | Criar modelo a partir de segredo copia estrutura sem vínculo                      | Funcional  | Unitário  | [RF-35](srs.md#rf-35), [RN-17](regras-negocio.md#rn-17)           |
| <a id="tc-d26"></a>TC-D26 | Transição de estado Cofre Salvo → Cofre Modificado em qualquer mutação            | Funcional  | Unitário  | —                  |
| <a id="tc-d27"></a>TC-D27 | Nomes repetidos permitidos para segredos, pastas e modelos                        | Funcional  | Unitário  | [RN-08](regras-negocio.md#rn-08)                  |
| <a id="tc-d28"></a>TC-D28 | Configurações embutidas no cofre: alterar e persistir tempos configuráveis        | Funcional  | Unitário  | [RF-11](srs.md#rf-11), [RN-04](regras-negocio.md#rn-04)           |

### 3.4 Importação e Exportação

| ID    | Caso de Teste                                                                      | Tipo       | Nível       | Requisitos             |
|-------|------------------------------------------------------------------------------------|------------|-------------|------------------------|
| <a id="tc-ie01"></a>TC-IE01 | Exportar cofre para texto claro com conteúdo correto                             | Funcional  | Integração  | [RF-09](srs.md#rf-09)                  |
| <a id="tc-ie02"></a>TC-IE02 | Importar cofre sem conflitos                                                     | Funcional  | Integração  | [RF-10](srs.md#rf-10)                  |
| <a id="tc-ie03"></a>TC-IE03 | Importar com conflito de identidade de segredo → nova identidade                 | Funcional  | Integração  | [RF-10](srs.md#rf-10), [RN-20](regras-negocio.md#rn-20)           |
| <a id="tc-ie04"></a>TC-IE04 | Importar com conflito de nome de segredo → sufixo incremental                    | Funcional  | Integração  | [RF-10](srs.md#rf-10), [RN-21](regras-negocio.md#rn-21)           |
| <a id="tc-ie05"></a>TC-IE05 | Importar com conflito de identidade de pasta → merge silencioso                  | Funcional  | Integração  | [RF-10](srs.md#rf-10), [RN-19](regras-negocio.md#rn-19)           |
| <a id="tc-ie06"></a>TC-IE06 | Importar com conflito de identidade de modelo → sobrescrita silenciosa           | Funcional  | Integração  | [RF-10](srs.md#rf-10), [RN-22](regras-negocio.md#rn-22)           |

### 3.5 Compatibilidade de Formato

| ID    | Caso de Teste                                                                      | Tipo            | Nível       | Requisitos        |
|-------|------------------------------------------------------------------------------------|-----------------|-------------|-------------------|
| <a id="tc-f01"></a>TC-F01 | Abrir cofre de versão anterior e migrar em memória                                | Compatibilidade | Integração  | [RNF-16](srs.md#rnf-16), [RNF-17](srs.md#rnf-17)    |
| <a id="tc-f02"></a>TC-F02 | Salvar cofre migrado no formato da versão atual                                   | Compatibilidade | Integração  | [RNF-17](srs.md#rnf-17)            |
| <a id="tc-f03"></a>TC-F03 | Selecionar perfil Argon2id correto para versão histórica                          | Compatibilidade | Unitário    | [RNF-18](srs.md#rnf-18)            |
| <a id="tc-f04"></a>TC-F04 | Rejeitar cofre com versão do formato superior                                     | Compatibilidade | Unitário    | [RNF-19](srs.md#rnf-19)            |

### 3.6 Segurança Operacional

| ID    | Caso de Teste                                                                      | Tipo       | Nível       | Requisitos       |
|-------|------------------------------------------------------------------------------------|------------|-------------|------------------|
| <a id="tc-s01"></a>TC-S01 | Limpeza de área de transferência após tempo configurado                           | Segurança  | Integração  | [RF-37](srs.md#rf-37)            |
| <a id="tc-s02"></a>TC-S02 | Limpeza de área de transferência ao bloquear cofre                                | Segurança  | Integração  | [RF-38](srs.md#rf-38)            |
| <a id="tc-s03"></a>TC-S03 | Bloqueio por inatividade após tempo configurado                                   | Funcional  | Integração  | [RF-08](srs.md#rf-08)     |
| <a id="tc-s04"></a>TC-S04 | Alerta de bloqueio iminente aos 75% do tempo de inatividade                       | Funcional  | Integração  | —            |
| <a id="tc-s05"></a>TC-S05 | Atividade do usuário reinicia cronômetro de inatividade                           | Funcional  | Integração  | —            |
| <a id="tc-s06"></a>TC-S06 | Bloqueio descarta alterações não salvas silenciosamente                           | Funcional  | Integração  | —     |
| <a id="tc-s07"></a>TC-S07 | Ausência de dados sensíveis em saída padrão e saída de erro                      | Segurança  | Unitário    | [RNF-07](srs.md#rnf-07)           |

### 3.7 Interface TUI (Golden Files e Comandos)

| ID    | Caso de Teste                                                                      | Tipo       | Nível       | Requisitos       |
|-------|------------------------------------------------------------------------------------|------------|-------------|------------------|
| <a id="tc-t01"></a>TC-T01 | Golden file: tela inicial (welcome) em 80×24                                     | Visual     | Componente  | —                |
| <a id="tc-t02"></a>TC-T02 | Golden file: tela de cofre ativo com hierarquia e segredo em 80×24               | Visual     | Componente  | —                |
| <a id="tc-t03"></a>TC-T03 | Golden file: tela de criação de segredo em 80×24                                 | Visual     | Componente  | —                |
| <a id="tc-t04"></a>TC-T04 | Golden file: tela de edição de segredo (padrão e avançada) em 80×24              | Visual     | Componente  | —                |
| <a id="tc-t05"></a>TC-T05 | Golden file: diálogo de confirmação bloqueante em 80×24                           | Visual     | Componente  | —                |
| <a id="tc-t06"></a>TC-T06 | Golden file: file picker em 80×24                                                | Visual     | Componente  | —                |
| <a id="tc-t07"></a>TC-T07 | Golden file: tela de busca com resultados filtrados em 80×24                      | Visual     | Componente  | —                |
| <a id="tc-t08"></a>TC-T08 | Golden file: mensagem de terminal muito pequeno                                  | Visual     | Componente  | —                |
| <a id="tc-t09"></a>TC-T09 | Comando: navegação na árvore (setas, expandir/colapsar)                          | Funcional  | Componente  | [RF-12](srs.md#rf-12)            |
| <a id="tc-t10"></a>TC-T10 | Comando: alternância de painéis com Tab                                          | Funcional  | Componente  | —                |
| <a id="tc-t11"></a>TC-T11 | Comando: toggle de campo sensível                                                | Funcional  | Componente  | [RF-14](srs.md#rf-14)            |
| <a id="tc-t12"></a>TC-T12 | Comando: cópia de campo para área de transferência                               | Funcional  | Componente  | [RF-36](srs.md#rf-36)            |

### 3.8 Testes de Integração E2E

| ID    | Caso de Teste                                                                      | Tipo       | Nível  | Requisitos                    |
|-------|------------------------------------------------------------------------------------|------------|--------|-------------------------------|
| <a id="tc-e01"></a>TC-E01 | Fluxo completo: criar cofre → criar segredo → salvar → fechar → reabrir → verificar dados | Funcional | E2E | [RF-01](srs.md#rf-01), [RF-02](srs.md#rf-02), [RF-03](srs.md#rf-03), [RF-16](srs.md#rf-16)  |
| <a id="tc-e02"></a>TC-E02 | Fluxo completo: abrir cofre → editar segredo → salvar → verificar alterações     | Funcional  | E2E    | [RF-02](srs.md#rf-02), [RF-03](srs.md#rf-03), [RF-20](srs.md#rf-20)          |
| <a id="tc-e03"></a>TC-E03 | Fluxo completo: criar cofre → excluir segredo → restaurar → salvar               | Funcional  | E2E    | [RF-22](srs.md#rf-22), [RF-23](srs.md#rf-23), [RF-03](srs.md#rf-03)          |
| <a id="tc-e04"></a>TC-E04 | Fluxo completo: criar cofre → exportar → importar em outro cofre → verificar     | Funcional  | E2E    | [RF-09](srs.md#rf-09), [RF-10](srs.md#rf-10)                 |
| <a id="tc-e05"></a>TC-E05 | Fluxo completo: criar cofre → alterar senha mestra → fechar → reabrir com nova senha | Funcional | E2E | [RF-06](srs.md#rf-06), [RF-02](srs.md#rf-02)                 |
| <a id="tc-e06"></a>TC-E06 | Fluxo completo: criar cofre → modificar → sair com confirmação de salvamento     | Funcional  | E2E    | [UC-11](casos-de-uso.md#uc-11)                        |
| <a id="tc-e07"></a>TC-E07 | Fluxo completo: criar cofre → criar modelo → criar segredo do modelo → verificar | Funcional  | E2E    | [RF-32](srs.md#rf-32), [RF-16](srs.md#rf-16)                 |
| <a id="tc-e08"></a>TC-E08 | Fluxo completo: abrir cofre antigo → migrar formato → salvar → reabrir           | Funcional  | E2E    | [RNF-16](srs.md#rnf-16), [RNF-17](srs.md#rnf-17)               |

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
