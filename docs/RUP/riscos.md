# Lista de Riscos — Abditum

| Item            | Detalhe                        |
|-----------------|--------------------------------|
| Projeto         | Abditum                       |
| Versão          | 1.0                           |
| Data            | 2026-03-25                     |
| Status          | Aprovado                       |

---

## 1. Introdução

Este documento identifica e classifica os riscos do projeto Abditum, avaliando probabilidade, impacto, exposição e estratégias de mitigação. Os riscos são revisados a cada iteração.

### Escalas

**Probabilidade:** Muito Baixa (1) | Baixa (2) | Média (3) | Alta (4) | Muito Alta (5)

**Impacto:** Insignificante (1) | Menor (2) | Moderado (3) | Grave (4) | Catastrófico (5)

**Exposição:** Probabilidade × Impacto

---

## 2. Riscos Identificados

### 2.1 Riscos de Segurança

| ID    | Risco                                                                                   | Prob. | Impacto | Expos. | Categoria    |
|-------|-----------------------------------------------------------------------------------------|-------|---------|--------|--------------|
| RS-01 | Uso incorreto de primitivas criptográficas (AES-256-GCM, Argon2id) introduzindo vulnerabilidades | 2     | 5       | 10     | Segurança    |
| RS-02 | Dados sensíveis residuais em memória após bloqueio/fechamento do cofre (Go não garante zeragem de memória gerenciada pelo GC) | 3     | 4       | 12     | Segurança    |
| RS-03 | Exposição de dados sensíveis em logs (stdout/stderr) por erro de programação              | 2     | 4       | 8      | Segurança    |
| RS-04 | Arquivo exportado em texto claro esquecido pelo usuário em local inseguro                | 3     | 4       | 12     | Segurança    |
| RS-05 | Parâmetros do Argon2id insuficientes para resistir a ataques futuros com hardware mais potente | 2     | 4       | 8      | Segurança    |
| RS-06 | Limpeza da área de transferência não funcionar em todos os sistemas operacionais/gerenciadores | 3     | 3       | 9      | Segurança    |

### 2.2 Riscos de Dados e Confiabilidade

| ID    | Risco                                                                                   | Prob. | Impacto | Expos. | Categoria       |
|-------|-----------------------------------------------------------------------------------------|-------|---------|--------|-----------------|
| RD-01 | Corrupção do arquivo do cofre durante salvamento (falha de energia, disco cheio)         | 2     | 5       | 10     | Confiabilidade  |
| RD-02 | Perda irrecuperável de dados por esquecimento da senha mestra (Zero Knowledge)           | 3     | 5       | 15     | Dados           |
| RD-03 | Falha na migração de formato ao abrir cofre de versão anterior                           | 2     | 4       | 8      | Compatibilidade |
| RD-04 | Conflitos não resolvidos corretamente durante importação, causando perda ou duplicação    | 2     | 3       | 6      | Dados           |
| RD-05 | Arquivo do cofre bloqueado por outro processo impede salvamento                           | 2     | 3       | 6      | Confiabilidade  |
| RD-06 | Rollback de backup (`.bak`/`.bak2`) falha por permissão de arquivo ou disco cheio        | 1     | 4       | 4      | Confiabilidade  |

### 2.3 Riscos de Portabilidade e Compatibilidade

| ID    | Risco                                                                                   | Prob. | Impacto | Expos. | Categoria       |
|-------|-----------------------------------------------------------------------------------------|-------|---------|--------|-----------------|
| RP-01 | Comportamento divergente de APIs do sistema de arquivos entre Windows, macOS e Linux     | 3     | 3       | 9      | Portabilidade   |
| RP-02 | Rename atômico não garantido em todos os sistemas de arquivos (ex: NFS, FAT32 em pendrives) | 3     | 4       | 12     | Portabilidade   |
| RP-03 | Terminal do usuário não suporta 256 cores ou tem dimensões insuficientes                  | 3     | 2       | 6      | Portabilidade   |
| RP-04 | Caracteres Unicode não renderizam corretamente em terminais legados                       | 2     | 2       | 4      | Portabilidade   |
| RP-05 | Argon2id com 256 MiB excede a memória disponível em máquinas com poucos recursos         | 2     | 3       | 6      | Portabilidade   |

### 2.4 Riscos de Usabilidade

| ID    | Risco                                                                                   | Prob. | Impacto | Expos. | Categoria    |
|-------|-----------------------------------------------------------------------------------------|-------|---------|--------|--------------|
| RU-01 | Complexidade da interface TUI afasta usuários não técnicos                               | 3     | 3       | 9      | Usabilidade  |
| RU-02 | Bloqueio por inatividade descartando alterações não salvas frustra o usuário              | 3     | 3       | 9      | Usabilidade  |
| RU-03 | Usuário não percebe que segredos na Lixeira serão perdidos permanentemente ao salvar      | 3     | 3       | 9      | Usabilidade  |
| RU-04 | File picker TUI menos intuitivo que file picker nativo do SO                              | 2     | 2       | 4      | Usabilidade  |

### 2.5 Riscos Técnicos e de Projeto

| ID    | Risco                                                                                   | Prob. | Impacto | Expos. | Categoria  |
|-------|-----------------------------------------------------------------------------------------|-------|---------|--------|------------|
| RT-01 | Complexidade da TUI (Bubble Tea) subestimada, atrasando cronograma                       | 3     | 3       | 9      | Técnico    |
| RT-02 | Cobertura de testes insuficiente nas camadas de criptografia e armazenamento             | 2     | 4       | 8      | Qualidade  |
| RT-03 | Colisão de NanoID (6 chars) em cofres com volume extremamente alto de elementos           | 1     | 3       | 3      | Técnico    |
| RT-04 | Ausência de CI/CD retarda detecção de regressões cross-platform                          | 2     | 3       | 6      | Processo   |

---

## 3. Ranking por Exposição

| Rank | ID    | Risco (resumo)                                        | Expos. |
|------|-------|------------------------------------------------------|--------|
| 1    | RD-02 | Perda de dados por esquecimento de senha mestra       | 15     |
| 2    | RS-02 | Dados sensíveis residuais em memória                  | 12     |
| 3    | RS-04 | Exportação em texto claro esquecida em local inseguro | 12     |
| 4    | RP-02 | Rename atômico não garantido em todos os FS           | 12     |
| 5    | RS-01 | Uso incorreto de primitivas criptográficas            | 10     |
| 6    | RD-01 | Corrupção do cofre durante salvamento                 | 10     |
| 7    | RS-06 | Limpeza de clipboard falha em alguns SOs              | 9      |
| 8    | RP-01 | Divergência de APIs de FS entre plataformas           | 9      |
| 9    | RU-01 | TUI complexa para não técnicos                        | 9      |
| 10   | RU-02 | Bloqueio descartando alterações frustra usuário        | 9      |

---

## 4. Plano de Mitigação

| ID    | Estratégia de Mitigação                                                                                                          |
|-------|----------------------------------------------------------------------------------------------------------------------------------|
| RS-01 | Revisão de código focada; usar bibliotecas criptográficas padrão do Go (`crypto/aes`, `crypto/cipher`, `golang.org/x/crypto/argon2`); testes de criptografia extensivos |
| RS-02 | Zeragem explícita de slices/buffers controlados pela aplicação; documentar limitações do GC; minimizar tempo de exposição         |
| RS-03 | Proibir qualquer `fmt.Print`/log com dados de negócio; revisão de código; testes que validam ausência de logs sensíveis           |
| RS-04 | Aviso explícito antes da exportação; nomear arquivo exportado com extensão identificável; recomendar exclusão após uso            |
| RS-05 | Documentar parâmetros e política de atualização; mecanismo de evolução por versão do formato; monitorar benchmarks de Argon2id    |
| RS-06 | Testar limpeza de clipboard em Windows, macOS e Linux; usar bibliotecas cross-platform; documentar limitações conhecidas          |
| RD-01 | Salvamento atômico (.tmp + rename); rotação de backup (.bak/.bak2); mensagem de erro com instrução de recuperação manual          |
| RD-02 | Aviso categórico e repetido sobre irrecuperabilidade na criação do cofre; responsabilidade explícita do usuário                   |
| RD-03 | Testes de regressão para cada versão histórica de formato; rotina de migração explícita e testada; erro claro para versão superior |
| RD-04 | Regras de conflito documentadas e testadas; sufixo incremental automático; mensagem informativa ao usuário                        |
| RD-05 | Detecção de lock file; mensagem de erro específica                                                                                |
| RD-06 | Validação de permissões antes de iniciar operação; mensagem de erro detalhada com caminho do backup                               |
| RP-01 | Testes cross-platform em CI (Windows, macOS, Linux); abstrair operações de arquivo                                                |
| RP-02 | Documentar limitação em sistemas de arquivos sem rename atômico; testar em FAT32; considerar fallback se viável                   |
| RP-03 | Detectar capacidades do terminal; fallback graceful; mensagem de redimensionamento                                                |
| RP-04 | Limitar uso de Unicode a caracteres amplamente suportados; testar em terminais comuns                                             |
| RP-05 | Documentar requisitos mínimos de hardware; piso mínimo de 128 MiB definido                                                       |
| RU-01 | Barra de ajuda contextual; atalhos intuitivos; documentação de uso                                                                |
| RU-02 | Documentar decisão de projeto; considerar alerta mais proeminente antes do bloqueio; alerta aos 75% de inatividade                |
| RU-03 | Toast ao excluir segredo informando sobre Lixeira; documentar que salvamento esvazia Lixeira                                      |
| RU-04 | Suportar mouse e autocomplete no file picker; autocompletar caminhos                                                              |
| RT-01 | Spike/prova de conceito da TUI cedo no projeto; usar componentes prontos do Bubble Tea                                            |
| RT-02 | Definir cobertura mínima de testes para criptografia e armazenamento; testes obrigatórios no CI                                   |
| RT-03 | Monitorar; NanoID 6 chars = 62⁶ (~56B combinações), risco prático desprezível para v1                                            |
| RT-04 | Implementar CI desde a primeira fase; builds e testes automatizados para 3 plataformas                                            |
