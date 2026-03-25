# Especificação de Requisitos de Software (SRS) — Abditum

| Item            | Detalhe                        |
|-----------------|--------------------------------|
| Projeto         | Abditum                       |
| Versão          | 1.0                           |
| Data            | 2026-03-25                     |
| Status          | Aprovado                       |

---

## 1. Introdução

### 1.1 Objetivo
Este documento especifica os requisitos funcionais, não funcionais e regras de negócio do Abditum — cofre de senhas portátil e seguro. Serve como contrato entre stakeholders e equipe de desenvolvimento.

### 1.2 Escopo
Cobre todos os requisitos do Abditum v1, incluindo ciclo de vida do cofre, gerenciamento de segredos, hierarquia, modelos de segredo, segurança, importação/exportação e restrições do sistema.

### 1.3 Referências
- Documento de Visão — `docs/RUP/visao.md`
- Glossário — `docs/RUP/glossario.md`
- Especificações de Caso de Uso — `docs/RUP/casos-de-uso.md`
- Documento descritivo — `descricao.md`

---

## 2. Requisitos Funcionais

### 2.1 Ciclo de Vida do Cofre

| ID     | Requisito                                                                                                    | Prioridade |
|--------|--------------------------------------------------------------------------------------------------------------|------------|
| RF-01  | Criar novo cofre informando caminho e senha mestra com confirmação por digitação dupla                       | Crítica    |
| RF-02  | Abrir cofre existente informando caminho e senha mestra                                                      | Crítica    |
| RF-03  | Salvar cofre no caminho atual                                                                                 | Crítica    |
| RF-04  | Salvar cofre em novo caminho (Salvar Como)                                                                    | Alta       |
| RF-05  | Descartar alterações não salvas e recarregar cofre a partir do arquivo                                        | Alta       |
| RF-06  | Alterar a senha mestra do cofre com confirmação por digitação dupla                                           | Alta       |
| RF-07  | Bloquear cofre manualmente, interrompendo o acesso e exigindo nova autenticação                               | Alta       |
| RF-08  | Bloquear cofre automaticamente após inatividade configurável                                                  | Alta       |
| RF-09  | Exportar cofre para formato legível em texto claro                                                            | Média      |
| RF-10  | Importar dados de formato legível em texto claro, com tratamento de conflitos                                 | Média      |
| RF-11  | Configurar parâmetros de proteção do cofre: duração de inatividade permitida, duração de exposição de dados sensíveis revelados e duração de retenção de dados extraídos do cofre | Média |

### 2.2 Navegação da Hierarquia do Cofre

| ID     | Requisito                                                                                                    | Prioridade |
|--------|--------------------------------------------------------------------------------------------------------------|------------|
| RF-12  | Exibir hierarquia do cofre (pastas e segredos)                                                                | Crítica    |
| RF-13  | Exibir detalhes do segredo selecionado                                                                        | Crítica    |
| RF-14  | Permitir ao usuário revelar e ocultar dados sensíveis de um segredo de forma controlada                       | Crítica    |
| RF-15  | Minimizar a duração da exposição de dados sensíveis revelados                                                 | Alta       |

### 2.3 Gerenciamento de Segredos

| ID     | Requisito                                                                                                    | Prioridade |
|--------|--------------------------------------------------------------------------------------------------------------|------------|
| RF-16  | Criar segredo a partir de modelo existente ou como segredo vazio                                              | Crítica    |
| RF-17  | Criar segredo na raiz do cofre ou em uma pasta                                                                | Crítica    |
| RF-18  | Duplicar segredo existente, com nova identidade e nome sufixado incrementalmente                              | Alta       |
| RF-19  | Favoritar e desfavoritar segredos                                                                             | Média      |
| RF-20  | Editar dados do segredo: nome, valores dos campos e observação                                                | Crítica    |
| RF-21  | Editar estrutura do segredo: incluir, renomear, excluir e reordenar campos                                   | Alta       |
| RF-22  | Excluir segredo reversivelmente (mover para Lixeira)                                                          | Crítica    |
| RF-23  | Restaurar segredo da Lixeira ao local e estado originais                                                      | Alta       |
| RF-24  | Mover segredo para outra pasta ou para a raiz                                                                 | Alta       |
| RF-25  | Reordenar segredo relativamente aos irmãos na mesma coleção                                                   | Média      |
| RF-26  | Buscar segredos por nome, nome de campo, valor de campos tipo texto ou observação                             | Alta       |

### 2.4 Gerenciamento de Hierarquia

| ID     | Requisito                                                                                                    | Prioridade |
|--------|--------------------------------------------------------------------------------------------------------------|------------|
| RF-27  | Criar pasta na raiz ou dentro de outra pasta                                                                  | Alta       |
| RF-28  | Renomear pasta                                                                                                | Alta       |
| RF-29  | Mover pasta para outra pasta ou para a raiz                                                                   | Alta       |
| RF-30  | Reordenar pasta relativamente às irmãs na mesma coleção                                                       | Média      |
| RF-31  | Excluir pasta, promovendo seus filhos (segredos e subpastas) para o nível pai                                 | Alta       |

### 2.5 Gerenciamento de Modelos de Segredo

| ID     | Requisito                                                                                                    | Prioridade |
|--------|--------------------------------------------------------------------------------------------------------------|------------|
| RF-32  | Criar modelo de segredo com campos personalizados (nome e tipo)                                               | Alta       |
| RF-33  | Editar modelo existente: alterar nome, incluir/alterar/excluir/reordenar campos                               | Alta       |
| RF-34  | Remover modelo de segredo                                                                                     | Alta       |
| RF-35  | Criar modelo a partir de segredo existente, copiando sua estrutura de campos                                  | Média      |

### 2.6 Compartilhamento de Dados

| ID     | Requisito                                                                                                    | Prioridade |
|--------|--------------------------------------------------------------------------------------------------------------|------------|
| RF-36  | Permitir ao usuário extrair dados de qualquer campo de segredo para uso externo imediato                       | Crítica    |
| RF-37  | Minimizar a duração de exposição de dados extraídos do cofre para uso externo                                  | Alta       |
| RF-38  | Garantir que dados extraídos do cofre não permaneçam expostos após o encerramento ou bloqueio do acesso        | Alta       |

### 2.7 Segurança

| ID     | Requisito                                                                                                    | Prioridade |
|--------|--------------------------------------------------------------------------------------------------------------|------------|
| RF-39  | Exigir digitação dupla ao criar ou alterar a senha mestra                                                     | Crítica    |
| RF-40  | Permitir ao usuário proteger rapidamente o conteúdo do cofre contra observação indevida                       | Alta       |
| RF-41  | Exigir ciência explícita do usuário sobre os riscos antes de gerar cópia não protegida do cofre               | Alta       |

---

## 3. Regras de Negócio

| ID     | Regra                                                                                                        |
|--------|--------------------------------------------------------------------------------------------------------------|
| RN-01  | O acesso ao conteúdo do cofre requer autenticação por senha mestra                                          |
| RN-02  | A senha mestra é irrecuperável; seu esquecimento resulta em perda total dos dados (Conhecimento Zero)       |
| RN-03  | A criação e a alteração da senha mestra exigem confirmação por dupla digitação                              |
| RN-04  | O cofre é autossuficiente: todas as informações necessárias para seu uso estão contidas nele                |
| RN-05  | Nenhum dado do cofre é transmitido pela rede; toda operação é local e offline                               |
| RN-06  | Cada segredo pertence a exatamente uma pasta ou à raiz do cofre — nunca a dois locais simultaneamente       |
| RN-07  | Cada pasta pertence a exatamente uma pasta pai ou à raiz do cofre — nunca a dois locais simultaneamente     |
| RN-08  | Nomes de segredos, pastas e modelos não são identificadores — nomes repetidos são permitidos                |
| RN-09  | Os campos de um segredo são classificados como dados comuns (texto) ou dados sensíveis (texto sensível)     |
| RN-10  | Dados sensíveis são protegidos por padrão e nunca participam de buscas                                      |
| RN-11  | A observação de um segredo é classificada como dado não sensível                                            |
| RN-12  | A exclusão de um segredo é reversível até a próxima persistência definitiva do cofre                        |
| RN-13  | Após a persistência definitiva, segredos marcados para exclusão são eliminados permanentemente               |
| RN-14  | Um segredo marcado para exclusão não pode ser editado                                                       |
| RN-15  | Ao restaurar um segredo cuja pasta de origem não exista mais, ele retorna à raiz do cofre                   |
| RN-16  | A exclusão de uma pasta é imediata e irreversível; seus filhos são promovidos ao nível hierárquico superior |
| RN-17  | Modelos de segredo são templates de criação — segredos criados a partir de um modelo não mantêm vínculo retroativo |
| RN-18  | Alterações na estrutura de um modelo afetam apenas criações futuras de segredos                              |
| RN-19  | Na importação, pastas com mesma identidade já existente são mescladas                                       |
| RN-20  | Na importação, segredos com identidade conflitante recebem nova identidade, preservando os demais dados     |
| RN-21  | Na importação, segredos com nome conflitante na mesma pasta recebem ajuste de nome para evitar ambiguidade  |
| RN-22  | Na importação, modelos com mesma identidade são substituídos pelo modelo importado                          |
| RN-23  | A exportação do cofre gera uma cópia não protegida de todos os dados, incluindo dados sensíveis             |

---

## 4. Requisitos Não Funcionais

### 4.1 Segurança

| ID     | Requisito                                                                                                    | Prioridade |
|--------|--------------------------------------------------------------------------------------------------------------|------------|
| RNF-01 | Criptografia AES-256-GCM para proteção dos dados do cofre                                                    | Crítica    |
| RNF-02 | Derivação de chave via Argon2id com custo de memória mínimo de 256 MiB e no mínimo 3 iterações               | Crítica    |
| RNF-03 | Cabeçalho do arquivo autenticado como AAD (Additional Authenticated Data) do AES-256-GCM                     | Crítica    |
| RNF-04 | Salt gerado aleatoriamente na criação do cofre, único por arquivo                                             | Crítica    |
| RNF-05 | Nonce regenerado a cada operação de salvamento                                                                | Crítica    |
| RNF-06 | Minimizar a retenção de dados sensíveis em memória e limpar buffers controlados ao bloquear ou fechar o cofre | Crítica    |
| RNF-07 | Ausência total de logs (stdout/stderr) que contenham caminhos de cofre, nomes de segredos ou valores de campos | Crítica    |

### 4.2 Portabilidade

| ID     | Requisito                                                                                                    | Prioridade |
|--------|--------------------------------------------------------------------------------------------------------------|------------|
| RNF-08 | Executável portátil único, sem necessidade de instalação                                                      | Crítica    |
| RNF-09 | Compatibilidade com Windows, macOS e Linux (64 bits)                                                          | Crítica    |
| RNF-10 | Nenhum dado persistido fora do arquivo do cofre, exceto artefatos transitórios (`.tmp`) e backups (`.bak`)    | Crítica    |

### 4.3 Confiabilidade

| ID     | Requisito                                                                                                    | Prioridade |
|--------|--------------------------------------------------------------------------------------------------------------|------------|
| RNF-11 | Salvamento atômico: escrita em arquivo `.tmp` seguida de rename, com remoção do `.tmp` em caso de falha      | Alta       |
| RNF-12 | Backup automático do arquivo anterior com extensão `.bak` antes de cada salvamento                            | Alta       |
| RNF-13 | Se já existir `.bak`, renomear para `.bak2` durante a operação; remover `.bak2` ao final em caso de sucesso; restaurar `.bak2` para `.bak` em caso de falha | Alta |
| RNF-14 | Em falha após geração de backup, informar explicitamente que existe um backup disponível para intervenção manual | Alta     |
| RNF-15 | Alertar se o arquivo do cofre estiver bloqueado por outro processo                                             | Média      |

### 4.4 Compatibilidade

| ID     | Requisito                                                                                                    | Prioridade |
|--------|--------------------------------------------------------------------------------------------------------------|------------|
| RNF-16 | A aplicação v.N abre cofres de qualquer versão anterior do formato suportada                                  | Alta       |
| RNF-17 | Ao abrir cofre antigo, migrar o payload em memória para o formato corrente; ao salvar, regravar no formato atual | Alta     |
| RNF-18 | Selecionar perfil histórico de derivação de chave a partir da versão do formato registrada no arquivo         | Alta       |
| RNF-19 | Rejeitar com erro claro cofres cuja versão do formato seja superior à suportada pela aplicação                | Alta       |

### 4.5 Formato do Arquivo

| ID     | Requisito                                                                                                    | Prioridade |
|--------|--------------------------------------------------------------------------------------------------------------|------------|
| RNF-20 | Extensão do arquivo: `.abditum`                                                                               | Crítica    |
| RNF-21 | Formato: stream binário com cabeçalho fixo (magic `ABDT`, versão do formato, salt, nonce) + payload criptografado | Crítica |
| RNF-22 | Assinatura mágica `ABDT` no início do arquivo para rejeição precoce de arquivos inválidos                    | Crítica    |
| RNF-23 | Codificação UTF-8 para suporte a caracteres especiais                                                         | Alta       |
| RNF-24 | Payload criptografado contém estrutura JSON do cofre                                                           | Crítica    |

### 4.6 Desempenho

| ID     | Requisito                                                                                                    | Prioridade |
|--------|--------------------------------------------------------------------------------------------------------------|------------|
| RNF-25 | Derivação de chave entre 0,8s e 1,5s em hardware compatível                                                  | Média      |
| RNF-26 | Paralelismo do Argon2id: até 4 threads, limitado pela quantidade disponível na máquina                        | Média      |

### 4.7 Parametrização Argon2id

| ID     | Requisito                                                                                                    | Prioridade |
|--------|--------------------------------------------------------------------------------------------------------------|------------|
| RNF-27 | Parâmetros fixos e hard-coded em v1, iguais para todos os cofres de uma mesma versão                          | Crítica    |
| RNF-28 | Custo de memória v1: 256 MiB; piso de segurança: 128 MiB; teto operacional: 512 MiB                         | Crítica    |
| RNF-29 | Mínimo de 3 iterações                                                                                         | Crítica    |
| RNF-30 | Alterações nos parâmetros somente por decisão explícita de versão, com validação de compatibilidade           | Alta       |

---

## 5. Requisitos Inversos (Fora de Escopo)

| ID     | Funcionalidade excluída                                                                                      |
|--------|--------------------------------------------------------------------------------------------------------------|
| RI-01  | Armazenamento na nuvem                                                                                       |
| RI-02  | Múltiplos cofres abertos simultaneamente                                                                      |
| RI-03  | Aplicação mobile ou web                                                                                       |
| RI-04  | Tags — pastas e grupos são suficientes para v1                                                                |
| RI-05  | Histórico de versões de segredos                                                                              |
| RI-06  | Alteração do tipo de um campo de segredo existente                                                            |
| RI-07  | Reautenticação para salvar                                                                                    |

---

## 6. Configurações Padrão

| Configuração                                    | Valor Padrão Sugerido |
|-------------------------------------------------|-----------------------|
| Tempo de bloqueio automático por inatividade     | 2 minutos             |
| Tempo de reocultação de campo sensível           | 15 segundos           |
| Tempo de limpeza automática da área de transferência | 30 segundos       |

---

## 7. Estrutura Inicial de Novos Cofres

### 7.1 Pastas Pré-definidas
- Sites
- Financeiro
- Serviços

### 7.2 Modelos de Segredo Pré-definidos

| Modelo           | Campos                                                |
|------------------|-------------------------------------------------------|
| Login            | URL (texto), Username (texto), Password (texto sensível) |
| Cartão de Crédito| Número do Cartão (texto sensível), Nome no Cartão (texto), Data de Validade (texto), CVV (texto sensível) |
| API Key          | Nome da API (texto), Chave de API (texto sensível)    |
