# Plano de Iteração — Abditum v1

| Item            | Detalhe                        |
|-----------------|--------------------------------|
| Projeto         | Abditum                       |
| Versão          | 1.0                           |
| Data            | 2026-03-25                     |
| Status          | Aprovado                       |

---

## 1. Introdução

Este documento define as fases e iterações do projeto Abditum v1, mapeando requisitos, casos de uso e critérios de conclusão para cada etapa. A estrutura segue o modelo RUP adaptado ao contexto de um projeto de escopo definido com equipe reduzida.

### 1.1 Referências
- Documento de Visão — `docs/RUP/visao.md`
- SRS — `docs/RUP/srs.md`
- Casos de Uso — `docs/RUP/casos-de-uso.md`
- Plano de Testes — `docs/RUP/plano-testes.md`
- Lista de Riscos — `docs/RUP/riscos.md`

---

## 2. Visão Geral das Fases

| Fase         | Foco                                                              | Iterações |
|--------------|-------------------------------------------------------------------|-----------|
| Concepção    | Visão, requisitos, riscos, planejamento                           | 1         |
| Elaboração   | Arquitetura, criptografia, formato de arquivo, CI                 | 2         |
| Construção   | Domínio, funcionalidades, TUI, integração                         | 4         |
| Transição    | Polimento, testes E2E, release                                    | 1         |

---

## 3. Detalhamento por Iteração

---

### Fase: Concepção

#### Iteração C1 — Visão e Planejamento

**Objetivo:** Estabelecer a visão do produto, documentar requisitos e planejar o projeto.

**Entregas:**
- Documento de Visão
- Glossário
- Especificação de Requisitos (SRS)
- Especificações de Caso de Uso
- Lista de Riscos
- Plano de Testes
- Plano de Iteração
- Estrutura inicial do repositório

**Requisitos endereçados:** Nenhum implementado — apenas documentados.

**Critérios de conclusão:**
- [x] Todos os artefatos RUP produzidos e revisados
- [ ] Repositório Git inicializado
- [ ] Estrutura de diretórios do projeto Go definida

---

### Fase: Elaboração

#### Iteração E1 — Fundação Criptográfica e Formato de Arquivo

**Objetivo:** Implementar as camadas de criptografia e formato de arquivo, validando a arquitetura de segurança.

**Entregas:**
- Serviço de criptografia (AES-256-GCM + Argon2id)
- Formato do arquivo `.abditum` (cabeçalho + payload)
- Leitura e escrita do arquivo binário
- Testes do serviço de criptografia ([TC-C01](plano-testes.md#tc-c01) a [TC-C10](plano-testes.md#tc-c10))

**Requisitos endereçados:**
- [RNF-01](srs.md#rnf-01) a [RNF-05](srs.md#rnf-05) (criptografia)
- [RNF-20](srs.md#rnf-20) a [RNF-24](srs.md#rnf-24) (formato de arquivo)
- [RNF-25](srs.md#rnf-25) a [RNF-30](srs.md#rnf-30) (Argon2id)

**Riscos mitigados:** [RS-01](riscos.md#rs-01), [RS-05](riscos.md#rs-05)

**Critérios de conclusão:**
- [ ] Criptografar e descriptografar payload com sucesso
- [ ] Rejeitar arquivo sem assinatura `ABDT`
- [ ] Rejeitar versão do formato superior à suportada
- [ ] Todos os TC-C* passando
- [ ] Salt e nonce únicos verificados nos testes

#### Iteração E2 — Armazenamento Atômico e CI

**Objetivo:** Implementar persistência confiável com salvamento atômico, rotação de backup e pipeline de CI.

**Entregas:**
- Serviço de armazenamento (salvar, carregar, backup, rollback)
- Salvamento atômico (.tmp + rename)
- Rotação de backup (.bak / .bak2)
- Pipeline de CI (build + testes em Windows, macOS, Linux)
- Testes do serviço de armazenamento ([TC-A01](plano-testes.md#tc-a01) a [TC-A10](plano-testes.md#tc-a10))

**Requisitos endereçados:**
- [RNF-08](srs.md#rnf-08) a [RNF-10](srs.md#rnf-10) (portabilidade)
- [RNF-11](srs.md#rnf-11) a [RNF-15](srs.md#rnf-15) (confiabilidade)
- [RNF-16](srs.md#rnf-16) a [RNF-19](srs.md#rnf-19) (compatibilidade — estrutura base)

**Riscos mitigados:** [RD-01](riscos.md#rd-01), [RD-05](riscos.md#rd-05), [RD-06](riscos.md#rd-06), [RP-01](riscos.md#rp-01), [RP-02](riscos.md#rp-02), [RT-04](riscos.md#rt-04)

**Critérios de conclusão:**
- [ ] Salvamento atômico funcional com .tmp + rename
- [ ] Rotação .bak/.bak2 correta em cenários de sucesso e falha
- [ ] CI executando testes em 3 plataformas
- [ ] Todos os TC-A* passando

---

### Fase: Construção

#### Iteração K1 — Modelo de Domínio e Manager

**Objetivo:** Implementar as entidades do domínio e o Manager como API de manipulação, com regras de negócio centralizadas.

**Entregas:**
- Entidades: Cofre, Segredo, Pasta, ModeloSegredo, Campo
- Manager com operações CRUD sobre entidades
- Identidade por NanoID (6 caracteres)
- Hierarquia recursiva (raiz como pasta sem nome)
- Configurações embutidas no cofre
- Estrutura inicial de novos cofres (pastas e modelos pré-definidos)
- Testes do domínio ([TC-D01](plano-testes.md#tc-d01) a [TC-D28](plano-testes.md#tc-d28))

**Requisitos endereçados:**
- [RF-01](srs.md#rf-01) (estrutura inicial), [RF-16](srs.md#rf-16) a [RF-26](srs.md#rf-26) (segredos)
- [RF-27](srs.md#rf-27) a [RF-31](srs.md#rf-31) (pastas), [RF-32](srs.md#rf-32) a [RF-35](srs.md#rf-35) (modelos)
- [RF-11](srs.md#rf-11) (configurações)
- [RN-01](regras-negocio.md#rn-01) a [RN-18](regras-negocio.md#rn-18)

**Riscos mitigados:** [RT-03](riscos.md#rt-03)

**Critérios de conclusão:**
- [ ] Todas as entidades implementadas com leitura somente via getters
- [ ] Mutações apenas via métodos do Manager
- [ ] Todos os TC-D* passando
- [ ] Regras de negócio [RN-01](regras-negocio.md#rn-01) a [RN-18](regras-negocio.md#rn-18) validadas por testes

#### Iteração K2 — Fluxos do Cofre e Integração Cripto + Domínio

**Objetivo:** Integrar domínio com criptografia e armazenamento para realizar os fluxos completos do ciclo de vida do cofre.

**Entregas:**
- Criar novo cofre (fluxo completo: domínio → criptografia → arquivo)
- Abrir cofre existente (arquivo → criptografia → domínio)
- Salvar cofre (domínio → criptografia → armazenamento atômico)
- Salvar Como (novo caminho)
- Descartar alterações e recarregar
- Alterar senha mestra
- Migração de formato de versões anteriores
- Importação e exportação JSON
- Testes de integração ([TC-E01](plano-testes.md#tc-e01) a [TC-E08](plano-testes.md#tc-e08))
- Testes de compatibilidade ([TC-F01](plano-testes.md#tc-f01) a [TC-F04](plano-testes.md#tc-f04))
- Testes de importação/exportação ([TC-IE01](plano-testes.md#tc-ie01) a [TC-IE06](plano-testes.md#tc-ie06))

**Requisitos endereçados:**
- [RF-01](srs.md#rf-01) a [RF-06](srs.md#rf-06), [RF-09](srs.md#rf-09), [RF-10](srs.md#rf-10)
- [RNF-16](srs.md#rnf-16) a [RNF-19](srs.md#rnf-19) (compatibilidade completa)
- [RN-13](regras-negocio.md#rn-13), [RN-19](regras-negocio.md#rn-19) a [RN-23](regras-negocio.md#rn-23) (persistência, importação, exportação)

**Riscos mitigados:** [RD-02](riscos.md#rd-02), [RD-03](riscos.md#rd-03), [RD-04](riscos.md#rd-04)

**Critérios de conclusão:**
- [ ] Fluxo criar → salvar → fechar → reabrir funcional
- [ ] Alteração de senha mestra com regravação completa
- [ ] Migração de formato de versão anterior validada
- [ ] Importação com todos os cenários de conflito testados
- [ ] Todos os TC-E*, TC-F*, TC-IE* passando

#### Iteração K3 — Interface TUI: Estrutura e Navegação

**Objetivo:** Implementar a interface TUI com layout principal, navegação na hierarquia e visualização de segredos.

**Entregas:**
- Layout principal com dois painéis (Hierarquia e Segredo)
- Barra de status e barra de ajuda contextual
- Árvore de navegação com expandir/colapsar
- Pastas virtuais: Favoritos e Lixeira
- Visualização de segredo com campos sensíveis ocultos
- Toggle de campos sensíveis com reocultação temporizada
- File picker integrado à TUI
- Detecção de tamanho mínimo do terminal
- Tela inicial (welcome / ASCII art)
- Golden files visuais ([TC-T01](plano-testes.md#tc-t01) a [TC-T08](plano-testes.md#tc-t08))
- Testes de comandos ([TC-T09](plano-testes.md#tc-t09) a [TC-T12](plano-testes.md#tc-t12))

**Requisitos endereçados:**
- [RF-07](srs.md#rf-07), [RF-08](srs.md#rf-08), [RF-12](srs.md#rf-12) a [RF-15](srs.md#rf-15), [RF-36](srs.md#rf-36) a [RF-41](srs.md#rf-41)

**Riscos mitigados:** [RT-01](riscos.md#rt-01), [RU-01](riscos.md#ru-01), [RU-04](riscos.md#ru-04), [RP-03](riscos.md#rp-03), [RP-04](riscos.md#rp-04)

**Critérios de conclusão:**
- [ ] Layout de dois painéis funcional e responsivo
- [ ] Navegação na árvore por teclado e mouse
- [ ] Toggle de campos sensíveis com reocultação automática
- [ ] File picker funcional com navegação e autocomplete
- [ ] Golden files aprovados para todas as telas principais
- [ ] Todos os TC-T* passando

#### Iteração K4 — Interface TUI: Operações e Fluxos Interativos

**Objetivo:** Implementar todas as operações interativas sobre segredos, pastas e modelos na TUI, completando os fluxos do usuário.

**Entregas:**
- Criação de segredo (com modelo ou vazio)
- Edição padrão e edição avançada de segredos
- Duplicação, favoritação, movimentação e reordenação de segredos
- Exclusão reversível e restauração (Lixeira)
- Busca com filtragem e highlight
- Operações de pasta (criar, renomear, mover, reordenar, excluir)
- Operações de modelo (criar, editar, remover, criar a partir de segredo)
- Bloqueio manual e por inatividade com alerta
- Cópia para área de transferência com limpeza temporizada
- Configuração do cofre
- Exportação e importação com confirmações
- Diálogos de confirmação e feedback (toasts)
- Sair da aplicação com tratamento de alterações não salvas
- Ocultação rápida da interface (shoulder surfing)

**Requisitos endereçados:**
- Todos os RF restantes ([RF-16](srs.md#rf-16) a [RF-41](srs.md#rf-41) na interface)
- Todos os UC ([UC-01](casos-de-uso.md#uc-01) a [UC-34](casos-de-uso.md#uc-34) nos fluxos interativos)

**Riscos mitigados:** [RS-02](riscos.md#rs-02), [RS-03](riscos.md#rs-03), [RS-04](riscos.md#rs-04), [RS-06](riscos.md#rs-06), [RU-02](riscos.md#ru-02), [RU-03](riscos.md#ru-03)

**Critérios de conclusão:**
- [ ] Todos os 34 casos de uso exercitáveis pela TUI
- [ ] Bloqueio por inatividade com alerta funcional
- [ ] Limpeza de clipboard funcional em 3 plataformas
- [ ] Todas as confirmações bloqueantes implementadas
- [ ] Busca com filtragem e highlight funcional
- [ ] Ausência de dados sensíveis em logs verificada

---

### Fase: Transição

#### Iteração T1 — Polimento, Testes E2E e Release

**Objetivo:** Estabilizar o produto, completar testes de ponta a ponta, resolver defeitos pendentes e preparar o release de v1.

**Entregas:**
- Correção de defeitos encontrados nas iterações anteriores
- Testes E2E completos ([TC-E01](plano-testes.md#tc-e01) a [TC-E08](plano-testes.md#tc-e08)) executados e aprovados
- Validação cross-platform final em CI
- Revisão de segurança final (ausência de logs, zeragem de buffers)
- Golden files finais aprovados
- Binários compilados para Windows, macOS e Linux
- Documentação de uso básica (se aplicável)

**Requisitos endereçados:** Validação final de todos os requisitos.

**Riscos mitigados:** [RT-02](riscos.md#rt-02) (cobertura de testes)

**Critérios de conclusão:**
- [ ] Zero defeitos de severidade Crítica ou Alta em aberto
- [ ] Todos os testes (unitários, integração, visuais, E2E) passando em 3 plataformas
- [ ] Cobertura ≥ 80% nas camadas de criptografia e armazenamento
- [ ] Binários gerados e testados para Windows, macOS e Linux
- [ ] Revisão de segurança concluída sem achados críticos

---

## 4. Matriz de Rastreabilidade: Iteração × Requisitos

| Iteração | Requisitos Funcionais                        | Requisitos Não Funcionais          | Regras de Negócio                    |
|----------|----------------------------------------------|------------------------------------|--------------------------------------|
| C1       | —                                            | —                                  | —                                    |
| E1       | —                                            | [RNF-01](srs.md#rnf-01) a 05, [RNF-20](srs.md#rnf-20) a 24, [RNF-25](srs.md#rnf-25) a 30 | —                              |
| E2       | —                                            | [RNF-08](srs.md#rnf-08) a 15, [RNF-16](srs.md#rnf-16) a 19 (base)   | —                                    |
| K1       | [RF-01](srs.md#rf-01)°, [RF-11](srs.md#rf-11), [RF-16](srs.md#rf-16) a 35                   | —                                  | [RN-01](regras-negocio.md#rn-01) a [RN-18](regras-negocio.md#rn-18)              |
| K2       | [RF-01](srs.md#rf-01) a 06, [RF-09](srs.md#rf-09), [RF-10](srs.md#rf-10)                    | [RNF-16](srs.md#rnf-16) a 19 (completo)            | [RN-13](regras-negocio.md#rn-13), [RN-19](regras-negocio.md#rn-19) a [RN-23](regras-negocio.md#rn-23)                           |
| K3       | [RF-07](srs.md#rf-07), [RF-08](srs.md#rf-08), [RF-12](srs.md#rf-12) a 15, [RF-36](srs.md#rf-36) a 41        | —                                  | —                         |
| K4       | [RF-16](srs.md#rf-16) a 41 (TUI)                             | [RNF-06](srs.md#rnf-06), [RNF-07](srs.md#rnf-07)                    | —          |
| T1       | Validação final de todos                     | Validação final de todos           | Validação final de todas              |

*° [RF-01](srs.md#rf-01) em K1 cobre apenas a estrutura inicial do cofre; o fluxo completo é realizado em K2.*

---

## 5. Dependências entre Iterações

```
C1 ──→ E1 ──→ E2 ──→ K1 ──→ K2 ──→ K3 ──→ K4 ──→ T1
       │       │       │       │
       │       │       │       └─ Integração domínio+cripto+armazem.
       │       │       └─ Requer criptografia e armazenamento prontos
       │       └─ Requer criptografia pronta para testar persistência
       └─ Fundação: tudo depende da criptografia
```

- **E1** é pré-requisito absoluto — sem criptografia, nada funciona
- **E2** depende de E1 para integrar persistência com criptografia
- **K1** depende de E1+E2 para que o domínio possa ser serializado e persistido
- **K2** integra K1 com E1+E2
- **K3** e **K4** dependem de K2 para apresentar e operar dados reais
- **T1** consolida tudo
