# Documento de Visão — Abditum

| Item            | Detalhe                        |
|-----------------|--------------------------------|
| Projeto         | Abditum                       |
| Versão          | 1.0                           |
| Data            | 2026-03-25                     |
| Status          | Aprovado                       |

---

## 1. Introdução

### 1.1 Objetivo
Este documento estabelece a visão do produto Abditum, definindo seu posicionamento, público-alvo, necessidades dos stakeholders, características principais, restrições e requisitos de qualidade. Ele serve como referência central para alinhar todas as partes envolvidas no desenvolvimento.

### 1.2 Escopo
O Abditum é um cofre de senhas portátil, seguro e fácil de usar. Permite que usuários armazenem e gerenciem suas senhas e informações confidenciais de forma organizada e protegida, sem depender de serviços em nuvem ou instalações complexas.

### 1.3 Referências
- `descricao.md` — Documento descritivo do produto Abditum

---

## 2. Posicionamento

### 2.1 Declaração do Problema

| Aspecto                  | Descrição                                                                                                  |
|--------------------------|------------------------------------------------------------------------------------------------------------|
| **O problema de**        | Gerenciar senhas e informações confidenciais de forma segura, portátil e sob controle total do usuário      |
| **Afeta**                | Usuários técnicos e não técnicos que precisam armazenar credenciais com segurança                          |
| **Cujo impacto é**       | Risco de vazamento de dados, dependência de serviços em nuvem, perda de controle sobre informações pessoais |
| **Uma solução adequada** | Um cofre de senhas offline, portátil, criptografado e autossuficiente em um único arquivo executável        |

### 2.2 Declaração de Posição do Produto

| Aspecto              | Descrição                                                                                                   |
|----------------------|-------------------------------------------------------------------------------------------------------------|
| **Para**             | Usuários que necessitam de um gerenciador de senhas seguro, offline e portátil                               |
| **Que**              | Precisam armazenar, organizar e acessar credenciais e informações confidenciais                              |
| **O Abditum**        | É um cofre de senhas portátil                                                                                |
| **Que**              | Permite gerenciar segredos com criptografia forte, sem dependência de nuvem, em um único arquivo executável  |
| **Diferente de**     | Gerenciadores de senhas baseados em nuvem ou com instalação complexa                                        |
| **Nosso produto**    | Oferece portabilidade extrema, soberania total dos dados e criptografia moderna, funcionando discretamente em qualquer computador |

---

## 3. Descrição dos Stakeholders e Usuários

### 3.1 Resumo dos Stakeholders

| Stakeholder        | Representa                            | Papel                                        |
|--------------------|---------------------------------------|----------------------------------------------|
| Desenvolvedor      | Equipe de desenvolvimento             | Implementa, testa e mantém o produto         |
| Usuário final      | Pessoas que armazenam senhas          | Utiliza o cofre no dia a dia                  |

### 3.2 Resumo dos Usuários

| Usuário             | Descrição                                                        | Stakeholder        |
|---------------------|------------------------------------------------------------------|---------------------|
| Usuário do cofre    | Pessoa que cria, abre, gerencia segredos e navega pela hierarquia | Usuário final       |

### 3.3 Ambiente do Usuário
- Sistemas operacionais: Windows, macOS e Linux
- Não requer conexão com a internet
- Não requer instalação — executável portátil de arquivo único
- Dados armazenados em arquivo local único

---

## 4. Visão Geral do Produto

### 4.1 Perspectiva do Produto
O Abditum é um produto standalone, sem dependências externas de serviços, bancos de dados ou infraestrutura de nuvem. Toda a operação ocorre localmente, com dados armazenados exclusivamente dentro do arquivo do cofre, garantindo soberania total do usuário sobre suas informações.

### 4.2 Premissas e Dependências
- O sistema de arquivos permite leitura e escrita no local onde o cofre é armazenado
- O usuário é responsável pela custódia da senha mestra (princípio Zero Knowledge — irrecuperável em caso de esquecimento)

---

## 5. Características do Produto

### 5.1 Portabilidade Extrema
Um único arquivo executável que pode ser copiado e usado em qualquer lugar — pendrive, HD externo, qualquer computador — sem instalação e sem persistir dados fora do arquivo do cofre. O controle e a propriedade dos dados ficam inteiramente nas mãos do usuário.

### 5.2 Segurança e Criptografia
- Proteção dos dados com criptografia forte e moderna
- Derivação de chave resistente a ataques de força bruta e ataques offline
- Impossibilidade de recuperação dos dados sem a senha mestra (Conhecimento Zero)
- Confirmação por digitação dupla ao criar ou alterar a senha mestra

### 5.3 Gerenciamento Hierárquico de Segredos
- Organização de segredos em pastas e subpastas com aninhamento livre
- Segredos compostos por campos personalizáveis, com dados comuns e dados sensíveis
- Modelos de segredo para padronizar a criação de novos segredos (ex: Login, Cartão de Crédito, API Key)
- Pastas e modelos fornecidos por padrão em novos cofres, mas editáveis e removíveis pelo usuário
- Busca por nome, campo, valor ou observação
- Segredos favoritos com destaque para acesso rápido

### 5.4 Proteção Operacional
- Bloqueio automático do cofre após período de inatividade
- Minimização da exposição de dados extraídos do cofre para uso externo
- Proteção de dados sensíveis contra exposição prolongada, com revelação somente por ação explícita do usuário
- Proteção rápida contra observação indevida do conteúdo do cofre
- Exclusão reversível de segredos (Lixeira) com possibilidade de restauração até a próxima persistência definitiva

### 5.5 Confiabilidade dos Dados
- Proteção contra perda de dados em caso de falha durante a persistência
- Preservação automática de versão anterior dos dados como rede de segurança
- Continuidade de acesso a cofres criados em versões anteriores do formato

### 5.6 Importação e Exportação
- Exportação do cofre para formato legível (com aviso explícito sobre riscos de segurança)
- Importação de dados externos para o cofre, com tratamento automático de conflitos (identidade, nomes duplicados, modelos)

### 5.7 Flexibilidade de Modelos
- O formato do segredo é flexível e personalizável
- O usuário pode criar seus próprios modelos de segredo com campos personalizados
- Modelos pré-definidos para tipos comuns (Login, Cartão de Crédito, API Key)
- Modelos funcionam como templates de criação — segredos criados a partir de um modelo não mantêm vínculo: alterações no modelo não afetam segredos existentes

---

## 6. Restrições

| ID   | Restrição                                                                                                |
|------|----------------------------------------------------------------------------------------------------------|
| R-01 | A aplicação não acessa nem armazena dados em nuvem                                                       |
| R-02 | Apenas um cofre pode estar aberto por vez                                                                 |
| R-03 | Não há mecanismo de recuperação de senha mestra (Conhecimento Zero)                                       |
| R-04 | Nenhum dado é persistido fora do arquivo do cofre, exceto artefatos transitórios e backups previstos       |
| R-05 | A aplicação não produz registros (logs) que contenham caminhos do cofre, nomes de segredos ou valores      |
| R-06 | A alteração do tipo de um campo de segredo existente não é suportada — exige exclusão e recriação          |
| R-07 | Reautenticação não é exigida para salvar — o cofre já está autenticado na sessão                           |

---

## 7. Requisitos de Qualidade

| Categoria           | Requisito                                                                                          | Prioridade |
|---------------------|-----------------------------------------------------------------------------------------------------|------------|
| Segurança           | Criptografia forte e moderna dos dados do cofre                                                     | Crítica    |
| Segurança           | Derivação de chave resistente a ataques de força bruta                                              | Crítica    |
| Segurança           | Proteção de dados sensíveis ao encerrar ou bloquear o acesso ao cofre                                | Crítica    |
| Portabilidade       | Execução sem instalação em Windows, macOS e Linux                                                   | Crítica    |
| Confiabilidade      | Salvamento seguro com backup automático e recuperação em caso de falha                              | Alta       |
| Compatibilidade     | Abertura e migração automática de cofres criados em versões anteriores                              | Alta       |
| Usabilidade         | Navegação intuitiva com ajuda contextual sobre ações e atalhos disponíveis                          | Alta       |
| Privacidade         | Ausência total de dados sensíveis em registros da aplicação                                         | Crítica    |

---

## 8. Funcionalidades Excluídas do Escopo (v1)

| Funcionalidade                  | Justificativa                                               |
|---------------------------------|-------------------------------------------------------------|
| TOTP (Autenticação 2FA)         | Postergado para v2                                          |
| Gerador de senhas               | Postergado para v2                                          |
| Senha falsa de coação           | Postergado para v2                                          |
| Compartilhamento via QR Code    | Postergado para v2                                          |
| Relatório de Saúde do Cofre     | Postergado para v2                                          |
| Keyfile / Token de Hardware     | Postergado para v2                                          |
| Tags                            | Pastas e grupos são suficientes para v1                     |
| Histórico de versões de segredos| Fora do escopo                                              |
| Múltiplos cofres simultâneos    | Decisão de projeto — um cofre por vez                       |
| Armazenamento na nuvem          | Fora do escopo — dados sob controle exclusivo do usuário    |
| Aplicação mobile ou web         | Fora do escopo — portabilidade via executável               |
