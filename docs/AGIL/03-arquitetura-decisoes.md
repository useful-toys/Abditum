# 03 — Arquitetura e Decisões Técnicas

## 03.1 Visão Arquitetural

### Visão geral

Abditum é uma aplicação **single-binary, single-process** escrita em Go, com interface TUI (Bubble Tea) e criptografia forte (AES-256-GCM + Argon2id). A arquitetura segue **Domain-Driven Design (DDD)** com separação clara entre domínio, infraestrutura e apresentação.

### Diagrama de camadas

```
┌─────────────────────────────────────────────────────────┐
│                    APRESENTAÇÃO (TUI)                    │
│  Bubble Tea Models · Views · Componentes · File Picker  │
│  Barra de Status · Barra de Ajuda · Painéis · Toasts   │
├─────────────────────────────────────────────────────────┤
│                    APLICAÇÃO (Manager)                   │
│  API de manipulação do cofre · Regras de negócio        │
│  Soft delete · Importação/Exportação · Conflitos        │
│  Configurações · Busca · Favoritos · Reordenação        │
├─────────────────────────────────────────────────────────┤
│                      DOMÍNIO (Entidades)                 │
│  Cofre · Segredo · Pasta · ModeloSegredo                │
│  CampoSegredo · CampoModeloSegredo · Configurações     │
│  Entidades somente leitura para navegação externa       │
├─────────────────────────────────────────────────────────┤
│                  INFRAESTRUTURA (Serviços)               │
│  Criptografia (AES-256-GCM + Argon2id)                 │
│  Armazenamento (leitura/escrita do .abditum)            │
│  Clipboard (copiar/limpar) · Migração de formato        │
└─────────────────────────────────────────────────────────┘
```

### Componentes principais

| Componente | Responsabilidade |
|---|---|
| **Domínio (Entidades)** | Modelagem de Cofre, Segredo, Pasta, ModeloSegredo, CampoSegredo, CampoModeloSegredo. Entidades expõem acesso somente leitura; mutações são proibidas fora do Manager. |
| **Manager (Aplicação)** | API pública para toda manipulação do cofre. Centraliza regras de negócio: criação, edição, soft delete, restauração, movimentação, reordenação, importação, exportação, busca, resolução de conflitos. |
| **Serviço de Criptografia** | Derivação de chave (Argon2id), criptografia/descriptografia (AES-256-GCM), geração de salt/nonce, validação de cabeçalho (magic, versão, AAD). |
| **Serviço de Armazenamento** | Leitura e escrita do arquivo `.abditum`. Salvamento atômico (.tmp → .bak com rotação → rename). Migração de formato em memória. |
| **Serviço de Clipboard** | Cópia para área de transferência, limpeza temporizada, limpeza ao bloquear/fechar. |
| **TUI (Apresentação)** | Modelos Bubble Tea para cada tela e estado. Delegação de ações ao Manager. Renderização de árvore, formulários, toasts, modais, file picker. |

### Fluxo de dados — Abrir cofre

```
File Picker → caminho
     │
     ▼
Serviço de Armazenamento → lê bytes do arquivo
     │
     ▼
Serviço de Criptografia → valida magic + versão_formato
     │                   → seleciona perfil Argon2id (por versão)
     │                   → derive key (salt + senha mestra)
     │                   → descriptografa payload (AES-256-GCM + AAD)
     │
     ▼
Migração de formato → converte JSON histórico → modelo corrente
     │
     ▼
Manager → carrega entidades do domínio em memória
     │
     ▼
TUI → exibe hierarquia + estado "Cofre Salvo"
```

### Fluxo de dados — Salvar cofre (atômico)

```
Manager → serializa domínio em memória → JSON
     │
     ▼
Serviço de Criptografia → gera novo nonce
     │                   → criptografa payload (AES-256-GCM + AAD do cabeçalho)
     │
     ▼
Serviço de Armazenamento
     │  1. Grava em .abditum.tmp
     │  2. Se .bak existe → renomeia para .bak2
     │  3. Copia .abditum atual → .bak
     │  4. Rename .tmp → .abditum (atômico)
     │  5. Se sucesso → remove .bak2
     │  6. Se falha → restaura .bak2 → .bak
     │
     ▼
TUI → estado "Cofre Salvo"
```

### Invariantes arquiteturais

- Um único cofre ativo por vez.
- Entidades do domínio são somente leitura para consumidores externos — toda mutação passa pelo Manager.
- Nenhum arquivo externo de configuração — tudo é embutido no `.abditum`.
- Nenhum dado sensível em stdout/stderr.
- Nonce regenerado a cada salvamento; salt gerado na criação.
- Cabeçalho do arquivo serve como AAD do GCM — integridade garantida sem checksum adicional.

---

## 03.2 ADRs (Architecture Decision Records)

### ADR-001: Linguagem Go com binário único

- **Contexto:** O produto precisa ser portátil, executável sem instalação e cross-platform (Windows, macOS, Linux).
- **Decisão:** Usar Go compilado como binário estático único.
- **Alternativas consideradas:** Rust (curva de aprendizado maior, menos acessível para fins didáticos), Python (não gera binário portátil sem empacotamento complexo).
- **Consequências:** Binário autossuficiente, compilação rápida, boa ergonomia para TUI (Bubble Tea), ecossistema criptográfico maduro (crypto/aes, x/crypto/argon2).

### ADR-002: DDD com Manager como API de mutação

- **Contexto:** Regras de negócio são numerosas (soft delete, importação com conflitos, reordenação, promoção de filhos) e precisam ser centralizadas.
- **Decisão:** Entidades expõem acesso somente leitura. Toda mutação passa por métodos explícitos de um Manager.
- **Alternativas consideradas:** Entidades com métodos de mutação direta (risco de bypass de regras), CQRS (overengineering para aplicação local).
- **Consequências:** Regras centralizadas, testabilidade alta, segurança contra manipulação direta.

### ADR-003: AES-256-GCM + Argon2id

- **Contexto:** Proteção contra acesso não autorizado em repouso e contra ataques offline de brute force.
- **Decisão:** AES-256-GCM para criptografia simétrica (payload) com Argon2id para derivação da chave a partir da senha mestra.
- **Alternativas consideradas:** XChaCha20-Poly1305 (margem de nonce maior, mas AES-GCM é mais amplamente suportado e suficiente com nonce de 96 bits regenerado por salvamento), PBKDF2 (muito rápido para derivação resistente a GPU).
- **Consequências:** Criptografia autenticada nativa (integridade + confidencialidade), custo de derivação alto (~1s com 256 MiB), sem checksum adicional necessário.

### ADR-004: NanoID de 6 caracteres como identidade

- **Contexto:** Segredos, pastas e modelos precisam de identificadores estáveis, serializáveis e independentes do nome.
- **Decisão:** NanoID alfanumérico de 6 caracteres (62⁶ ≈ 56 bilhões de combinações).
- **Alternativas consideradas:** UUID v4 (36 caracteres, verboso para JSON e UX), auto-incremento (não portável entre importação/exportação).
- **Consequências:** Baixa probabilidade prática de colisão, boa legibilidade no JSON, portabilidade entre cofres.

### ADR-005: Cabeçalho como AAD do AES-GCM

- **Contexto:** O cabeçalho do arquivo (magic, versão_formato, salt, nonce) precisa ter integridade garantida.
- **Decisão:** Usar o cabeçalho inteiro como Additional Authenticated Data (AAD) do AES-256-GCM.
- **Alternativas consideradas:** Checksum separado (SHA-256 do cabeçalho) — redundante, pois GCM já valida AAD.
- **Consequências:** Integridade do cabeçalho validada nativamente na descriptografia, sem campo adicional no formato.

### ADR-006: Salvamento atômico com rotação de backup

- **Contexto:** Falhas durante escrita (queda de energia, disco cheio, file lock) não podem corromper o cofre.
- **Decisão:** Escrita sequencial: .tmp → backup .bak (com rotação .bak2) → rename atômico.
- **Alternativas consideradas:** Write-ahead log (overengineering para arquivo único), escrita direta com checksum (risco de corrupção parcial).
- **Consequências:** Em caso de falha, sempre existe um backup íntegro; usuário é informado da existência do .bak para recuperação manual.

### ADR-007: Tudo dentro do cofre (portabilidade extrema)

- **Contexto:** O produto precisa funcionar sem nenhum arquivo externo além do executável e do .abditum.
- **Decisão:** Modelos de segredo, configurações (tempos de inatividade, reocultação, clipboard), hierarquia de pastas — tudo persistido dentro do payload JSON criptografado.
- **Alternativas consideradas:** Arquivo de configuração externo (.config), variáveis de ambiente — ambos violam a portabilidade.
- **Consequências:** Cada arquivo .abditum é 100% autossuficiente. Copiar o arquivo = copiar o cofre completo com suas personalizações.

### ADR-008: Compatibilidade retroativa com migração em memória

- **Contexto:** Evolução do formato do arquivo ao longo do tempo (novos campos, mudança de estrutura).
- **Decisão:** Campo `versão_formato` no cabeçalho. Ao abrir, a aplicação reconhece formatos históricos e migra o JSON em memória. Ao salvar, sempre grava no formato corrente.
- **Alternativas consideradas:** Migração persistida (reescrever imediatamente no formato novo ao abrir) — risco de corrupção se a migração falhar; formato estável sem evolução — impraticável a longo prazo.
- **Consequências:** Migração segura (apenas em memória), salvamento sempre no formato mais recente, tabela hard-coded de perfis Argon2id por versão.

### ADR-009: Modelo como snapshot (sem vínculo)

- **Contexto:** Segredos criados a partir de modelos poderiam manter referência ao modelo original ou copiar a estrutura.
- **Decisão:** Cópia por snapshot. O nome do modelo é salvo como registro histórico, mas não há vínculo por referência. Alterações no modelo não afetam segredos existentes.
- **Alternativas consideradas:** Vínculo por referência (complexidade de propagação de mudanças, risco de efeitos colaterais).
- **Consequências:** Simplicidade, previsibilidade, independência entre modelos e segredos.

### ADR-010: Bubble Tea como framework TUI

- **Contexto:** A interface deve ser fullscreen, interativa, com 256 cores, suporte a teclado e mouse.
- **Decisão:** Usar Bubble Tea (Elm Architecture para Go) com Lip Gloss para estilização.
- **Alternativas consideradas:** tview (API imperativa, menos idiomática), tcell direto (muito baixo nível).
- **Consequências:** Arquitetura reativa e composável, boa ergonomia para testes (teatest/v2), ecossistema ativo.

---

## 03.3 Modelo de Domínio

### Entidades e relações

```
                            Cofre
                ┌────────────┼────────────┐
                │            │            │
          Configurações   Pastas[]    ModelosSegredo[]
                         ┌──┴──┐          │
                    Segredos[] Pastas[]  CamposModelo[]
                         │     (recursivo)
                   CamposSegredo[]
```

### Detalhamento das entidades

**Cofre** (raiz do domínio)

| Atributo | Tipo | Descrição |
|---|---|---|
| configurações | Configurações | Tempos de inatividade, reocultação e clipboard |
| segredos | list[Segredo] | Segredos na raiz do cofre |
| pastas | list[Pasta] | Pastas na raiz do cofre |
| modelos_segredo | list[ModeloSegredo] | Modelos disponíveis para criação |
| data_criação | datetime | Data de criação do cofre |
| data_última_modificação | datetime | Data da última alteração |

**Segredo**

| Atributo | Tipo | Descrição |
|---|---|---|
| id | NanoID (6 chars) | Identificador único e estável |
| nome | string | Nome do segredo (editável, não identificador) |
| nome_modelo | string (opcional) | Snapshot do nome do modelo usado na criação |
| campos | list[CampoSegredo] | Campos de dados do segredo |
| favorito | booleano | Marcação de destaque |
| observação | string (opcional) | Texto livre, não sensível |
| data_criação | datetime | Data de criação |
| data_última_modificação | datetime | Data da última alteração |

**Pasta** (recursiva)

| Atributo | Tipo | Descrição |
|---|---|---|
| id | NanoID (6 chars) | Identificador único e estável |
| nome | string | Nome da pasta (editável, não identificador) |
| segredos | list[Segredo] | Segredos contidos |
| pastas | list[Pasta] | Subpastas (recursivo) |

**ModeloSegredo**

| Atributo | Tipo | Descrição |
|---|---|---|
| id | NanoID (6 chars) | Identificador único e estável |
| nome | string | Nome do modelo |
| campos | list[CampoModeloSegredo] | Estrutura de campos do modelo |

**CampoSegredo**

| Atributo | Tipo | Descrição |
|---|---|---|
| nome | string | Nome do campo |
| tipo | enum (texto, texto_sensível) | Tipo do dado armazenado |
| valor | string (opcional) | Valor armazenado (pode ser vazio) |

**CampoModeloSegredo**

| Atributo | Tipo | Descrição |
|---|---|---|
| nome | string | Nome do campo |
| tipo | enum (texto, texto_sensível) | Tipo do dado |

**Configurações**

| Atributo | Tipo | Padrão |
|---|---|---|
| tempo_bloqueio_inatividade_minutos | inteiro | 2 |
| tempo_ocultar_segredo_segundos | inteiro | 15 |
| tempo_limpar_area_transferencia_segundos | inteiro | 30 |

### Decisões de modelagem

- **Hierarquia recursiva:** A raiz funciona como uma pasta sem nome. Pastas contêm segredos e subpastas em qualquer nível.
- **Ordenação por posição:** A ordem no JSON = ordem de exibição (segredos primeiro, depois subpastas).
- **Nomes não identificadores:** Nomes repetidos são permitidos — a identidade é pelo NanoID.
- **Campos uniformes:** Valor vazio = string vazia (campo existente, não preenchido).
- **Observação implícita:** Todo segredo possui observação — nunca declarada nos modelos, nunca removível.
- **Busca sequencial:** Sem índices persistidos ou em memória de longa duração.

---

## 03.4 Glossário

| Termo | Definição |
|---|---|
| **Senha mestra** | Chave de acesso ao cofre, usada para derivar a chave criptográfica via Argon2id. |
| **Cofre** | Arquivo `.abditum` criptografado que armazena segredos, pastas, modelos e configurações. |
| **Bloqueio do cofre** | Fechamento lógico do cofre: limpa memória, limpa clipboard e retorna ao fluxo de abertura exigindo nova autenticação. |
| **Segredo** | Item individual com campos de dados (comuns e sensíveis), nome, observação e marcação de favorito. |
| **Segredo favorito** | Segredo com destaque visual e presença na pasta virtual de Favoritos. |
| **Dados comuns** | Informações não sensíveis (ex: nome de serviço, URL). |
| **Dados sensíveis** | Informações confidenciais (ex: senhas, API keys) — campos do tipo `texto sensível`. |
| **Observação** | Campo de texto livre, não sensível, presente em todo segredo. |
| **Campo de segredo** | Elemento com nome, tipo (texto ou texto sensível) e valor. |
| **Hierarquia do cofre** | Organização em árvore de pastas e segredos com profundidade ilimitada. |
| **Raiz do cofre** | Nível mais alto da hierarquia — contém segredos e pastas não aninhadas. |
| **Pasta** | Contêiner estrutural para agrupar segredos e subpastas. |
| **Pasta virtual** | Agrupamento lógico (Favoritos, Lixeira) sem alterar localização real dos segredos. |
| **Modelo de segredo** | Estrutura reutilizável de campos para criação padronizada de segredos. |
| **Conhecimento Zero (Zero Knowledge)** | Princípio de que a aplicação não pode recuperar dados sem a senha mestra. |
| **Exclusão reversível (Soft Delete)** | Segredo removido da hierarquia mas restaurável até o próximo salvamento. |
| **Lixeira** | Pasta virtual que materializa segredos excluídos reversivelmente. |
| **Salvamento atômico** | Escrita via .tmp → backup .bak → rename, garantindo integridade em caso de falha. |
| **AAD (Additional Authenticated Data)** | Dados autenticados mas não criptografados pelo GCM — o cabeçalho do arquivo. |
| **Shoulder surfing** | Espionagem física da tela, mitigada por atalho de ocultação rápida. |
| **NanoID** | Gerador de identificadores curtos (6 chars alfanuméricos, 62⁶ ≈ 56 bilhões de combinações). |
| **Manager** | Camada da aplicação que centraliza regras de negócio e expõe a API de mutação do domínio. |
