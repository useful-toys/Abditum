# 01 — Visão do Produto

## 01.1 Visão do Produto

### Problema a ser resolvido

Gerenciadores de senhas tradicionais apresentam uma ou mais das seguintes limitações:

- **Dependência de nuvem:** exigem conta em serviço de terceiros, transferindo o controle dos dados para fora do domínio do usuário.
- **Instalação obrigatória:** requerem setup no sistema operacional, impedindo uso discreto e imediato em máquinas temporárias ou restritivas.
- **Persistência espalhada:** gravam configurações, caches ou logs em diretórios do sistema, deixando rastros fora do arquivo do cofre.
- **Formato rígido de segredos:** oferecem campos fixos (usuário/senha), sem possibilidade de o usuário definir modelos personalizados para diferentes tipos de informação confidencial.
- **Interface desktop pesada:** dependem de ambiente gráfico completo, inviabilizando uso em servidores, containers ou sessões SSH remotas.

O resultado é que usuários técnicos que precisam de autonomia, portabilidade e controle total sobre seus dados sensíveis recorrem a soluções improvisadas (arquivos de texto, planilhas, notas adesivas) ou aceitam trade-offs de privacidade e praticidade.

### Público-alvo

| Perfil | Descrição |
|---|---|
| **Profissionais de TI e DevOps** | Administradores de sistemas, desenvolvedores e engenheiros que operam em múltiplas máquinas, servidores e ambientes restritos (SSH, containers), necessitando de acesso rápido a credenciais sem dependência de GUI ou nuvem. |
| **Usuários avançados preocupados com privacidade** | Pessoas que recusam delegar o armazenamento de senhas a serviços de terceiros e exigem controle absoluto sobre o local e o formato de armazenamento dos seus dados. |
| **Usuários em ambientes corporativos restritivos** | Profissionais cujas estações de trabalho não permitem instalação de software adicional, mas que podem executar binários portáteis diretamente de pendrive ou pasta local. |

### Proposta de valor

**Abditum** é um cofre de senhas portátil, seguro e autossuficiente, que:

- Opera como um **único arquivo executável** — sem instalação, sem dependências de runtime, sem persistência externa.
- Armazena **tudo dentro de um único arquivo `.abditum`** criptografado — segredos, modelos, configurações e hierarquia.
- Oferece **interface TUI moderna** (fullscreen, 256 cores, teclado + mouse) que funciona em qualquer terminal, incluindo sessões remotas.
- Garante **segurança de nível profissional** com AES-256-GCM + Argon2id, sem concessões à conveniência.
- Permite **modelos de segredo personalizáveis**, indo além do par usuário/senha para suportar qualquer estrutura de dados confidenciais.
- Aplica o princípio de **Conhecimento Zero** — não existe mecanismo de recuperação de senha mestra. O controle é absoluto e intencional.

### Benefícios esperados

| Benefício | Descrição |
|---|---|
| **Soberania dos dados** | Nenhum dado transita por servidores de terceiros. O usuário decide onde o arquivo vive — pendrive, HD local, pasta de rede. |
| **Portabilidade total** | Funciona em Windows, macOS e Linux. Basta copiar o executável e o arquivo do cofre para qualquer máquina. |
| **Zero rastros no sistema** | Não grava configurações, logs, caches ou qualquer artefato fora do arquivo do cofre (exceto backups explícitos e arquivos temporários de salvamento). |
| **Flexibilidade de modelagem** | Modelos de segredo personalizáveis permitem ao usuário definir campos próprios para qualquer tipo de informação confidencial. |
| **Segurança em camadas** | Criptografia forte em repouso, derivação de chave custosa contra brute force, bloqueio automático por inatividade, limpeza da área de transferência, proteção contra shoulder surfing. |
| **Resiliência de dados** | Salvamento atômico com backup automático e rotação de arquivos `.bak` garante que falhas nunca corrompam o cofre sem possibilidade de recuperação. |

### Métricas de sucesso

| Métrica | Critério |
|---|---|
| **Portabilidade efetiva** | O binário executa em Windows, macOS e Linux 64-bit sem nenhuma dependência externa ou arquivo de configuração. |
| **Cobertura de testes** | Testes unitários, de integração e golden files cobrindo criptografia, armazenamento, navegação, transições de estado e fluxos de usuário. |
| **Desempenho de abertura** | Derivação de chave (Argon2id 256 MiB, 3 iterações) concluída na faixa interativa de 0,8 s a 1,5 s em hardware compatível. |
| **Integridade do salvamento** | Nenhuma perda de dados em cenários de falha de escrita — backup sempre disponível para recuperação manual. |
| **Privacidade de runtime** | Zero ocorrências de caminhos de arquivo, nomes de segredos ou valores de campos em stdout/stderr durante a execução. |
| **Compatibilidade retroativa** | A versão N da aplicação abre e migra corretamente cofres criados em todas as versões anteriores suportadas. |

---

## 01.2 Escopo do Produto

### O que está dentro do escopo (v1)

- **Ciclo de vida completo do cofre:** criar, abrir, salvar, salvar como, bloquear (manual e por inatividade), alterar senha mestra, descartar alterações e recarregar.
- **Gerenciamento de segredos:** criar (a partir de modelo ou vazio), editar (modo padrão e avançado), duplicar, favoritar/desfavoritar, excluir reversivelmente (soft delete com Lixeira), restaurar, mover entre pastas e reordenar.
- **Gerenciamento de hierarquia:** criar, renomear, mover, reordenar e excluir pastas (com promoção de filhos).
- **Modelos de segredo:** criar, editar, excluir, criar a partir de segredo existente. Modelos pré-definidos (Login, Cartão de Crédito, API Key) fornecidos na criação do cofre.
- **Navegação e busca:** árvore hierárquica navegável, busca em memória por nome/campo/valor/observação (excluindo campos sensíveis), pastas virtuais (Favoritos e Lixeira).
- **Área de transferência:** copiar campos com limpeza automática temporizada e ao bloquear/fechar o cofre.
- **Importação/Exportação:** exportar para JSON plain text (com aviso de segurança), importar de JSON plain text com tratamento de conflitos.
- **Configuração embutida:** tempo de bloqueio por inatividade, tempo de reocultação de campos sensíveis, tempo de limpeza da área de transferência — tudo persistido dentro do cofre.
- **Segurança:** AES-256-GCM + Argon2id, salvamento atômico com rotação de backup, proteção contra shoulder surfing, minimização de dados sensíveis em memória.
- **Interface TUI:** fullscreen, Bubble Tea, 256 cores, teclado + mouse, barra de status, barra de ajuda contextual, file picker integrado.
- **Multiplataforma:** Windows, macOS e Linux (binário Go compilado estaticamente).
- **Compatibilidade retroativa:** abertura e migração de formatos históricos do cofre.

### O que está fora do escopo

| Item excluído | Justificativa |
|---|---|
| Armazenamento na nuvem | Contraria o princípio de soberania e portabilidade local. |
| Múltiplos cofres abertos simultaneamente | Complexidade desnecessária para v1; um cofre ativo por vez é suficiente. |
| Aplicação mobile ou web | A TUI portátil é a proposta central do produto. |
| Sistema de tags | Pastas e hierarquia são suficientes para organização em v1. |
| Histórico de versões de segredos | Escopo de auditoria postergado para v2. |
| Alteração do tipo de um campo existente | Decisão deliberada — excluir e recriar é o caminho previsto. |
| Reautenticação para salvar | O cofre já está desbloqueado na sessão; seria redundante. |
| Gerador de senhas | Postergado para v2. |
| Suporte a TOTP | Postergado para v2. |
| Senha falsa de coação (Duress Password) | Postergado para v2. |
| Compartilhamento via QR Code | Postergado para v2. |
| Relatório de saúde do cofre (auditoria) | Postergado para v2. |
| Autenticação por keyfile / token de hardware | Postergado para v2. |

---

## 01.3 Roadmap do Produto

### Curto prazo — v1 (versão inicial)

Entrega do produto funcional completo com todas as capacidades definidas no escopo v1.

| Marco | Descrição |
|---|---|
| **M1 — Infraestrutura e criptografia** | Serviço de criptografia (AES-256-GCM + Argon2id), formato do arquivo `.abditum`, leitura/escrita com salvamento atômico e rotação de backup. Testes de criptografia e armazenamento. |
| **M2 — Domínio e regras de negócio** | Modelo de domínio (Cofre, Segredo, Pasta, Modelo de Segredo), Manager com API de manipulação, regras de soft delete, importação/exportação, resolução de conflitos. Testes unitários de domínio. |
| **M3 — Interface TUI (estrutura)** | Tela inicial, file picker, fluxos de criar/abrir cofre, layout de dois painéis, navegação na árvore, barra de status e barra de ajuda contextual. Golden files visuais. |
| **M4 — Interface TUI (operações)** | Edição padrão e avançada de segredos, gerenciamento de pastas e modelos, busca, Favoritos, Lixeira, configurações do cofre, área de transferência com limpeza temporizada. |
| **M5 — Segurança e polish** | Bloqueio por inatividade, proteção contra shoulder surfing, minimização de dados em memória, privacidade de logs, compatibilidade cross-platform, testes de integração E2E. |
| **M6 — Release v1** | Build cross-platform (Windows, macOS, Linux), CI obrigatório, validação final e publicação. |

### Médio prazo — v2 (evolução planejada)

| Funcionalidade | Descrição |
|---|---|
| Gerador de senhas | Geração configurável de senhas fortes diretamente na interface. |
| Suporte a TOTP | Cálculo e exibição de tokens 6 dígitos em tempo real para campos com chave secreta TOTP. |
| Relatório de saúde do cofre | Auditoria local de senhas fracas, reutilizadas ou antigas. |
| QR Code na TUI | Renderização ASCII de QR codes para transferência offline de campos sensíveis. |

### Longo prazo — v3+ (horizonte exploratório)

| Funcionalidade | Descrição |
|---|---|
| Senha de coação (Duress Password) | Senha alternativa que abre uma versão restrita do cofre sob ameaça. |
| Autenticação por keyfile / YubiKey | Segundo fator offline via arquivo ou token de hardware. |
| Evolução do formato de arquivo | Migrações estruturais acompanhadas de compatibilidade retroativa. |

### Dependências relevantes

| Dependência | Tipo | Impacto |
|---|---|---|
| **Go (compilador)** | Toolchain | Compilação do binário portátil para Windows, macOS e Linux. |
| **Bubble Tea / Lip Gloss** | Biblioteca | Framework TUI — define a arquitetura de toda a camada de apresentação. |
| **golang.org/x/crypto (argon2, chacha20)** | Biblioteca | Implementação de Argon2id para derivação de chave. |
| **crypto/aes + crypto/cipher (stdlib)** | Biblioteca | Implementação de AES-256-GCM para criptografia do payload. |
| **teatest/v2** | Biblioteca | Framework de testes para golden files e testes de comandos da TUI. |
| **NanoID** | Biblioteca | Geração de identificadores únicos de 6 caracteres para entidades. |

---

## 01.4 Stakeholders

### Negócio

| Stakeholder | Papel | Interesse |
|---|---|---|
| **Idealizador / Product Owner** | Define e prioriza requisitos, valida entregas | Garantir que o produto entregue a proposta de valor: portabilidade, segurança e soberania dos dados. |

### TI

| Stakeholder | Papel | Interesse |
|---|---|---|
| **Desenvolvedor(es)** | Implementação, testes e CI | Código didático (DDD, comentários generosos), arquitetura testável, build reproduzível cross-platform. |
| **Revisor de segurança** | Validação das práticas criptográficas | Parametrização correta do Argon2id, uso adequado do AES-256-GCM (nonce/salt/AAD), ausência de vazamentos em logs. |

### Operação

| Stakeholder | Papel | Interesse |
|---|---|---|
| **CI/Pipeline** | Build automatizado e validação contínua | Compilação cross-platform, execução de testes (unitários, integração, golden files) a cada commit. |
| **Distribuição** | Disponibilização do binário | Garantir que o executável seja autossuficiente, sem dependências de runtime ou instalação. |

### Usuários finais

| Stakeholder | Perfil | Expectativa |
|---|---|---|
| **Profissional de TI / DevOps** | Acesso a credenciais em múltiplos ambientes (local, SSH, servidor) | Executável portátil, abertura rápida, busca eficiente, cópia para clipboard com limpeza automática. |
| **Usuário avançado (privacidade)** | Recusa dependência de terceiros para dados sensíveis | Controle absoluto, zero nuvem, zero rastros, formato aberto para exportação. |
| **Usuário em ambiente restritivo** | Não pode instalar software na estação | Executável copiável para pendrive/pasta, sem persistência fora do arquivo do cofre. |
