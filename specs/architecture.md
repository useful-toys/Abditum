# Abditum — Arquitetura

## Visão geral

A aplicação segue uma arquitetura em camadas com separação clara entre domínio, aplicação e interface. O objetivo é manter o domínio e a lógica de criptografia completamente independentes do framework de UI.

```
┌─────────────────────────────────────┐
│            UI Layer                 │  ← GUI ou TUI (ver ADR 002)
├─────────────────────────────────────┤
│         Application Layer           │  ← Casos de uso, orquestração
├─────────────────────────────────────┤
│           Domain Layer              │  ← Entidades, regras de negócio
├─────────────────────────────────────┤
│       Infrastructure Layer          │  ← Criptografia, persistência
└─────────────────────────────────────┘
```

## Camadas

### UI Layer

Responsável pela renderização e interação com o usuário. Comunica-se apenas com a Application Layer. Não conhece entidades de domínio diretamente — usa DTOs/ViewModels quando necessário. Framework a ser definido (ver `decisions/adr-002-ui-framework.md`).

### Application Layer

Orquestra os casos de uso da aplicação:

| Serviço           | Responsabilidades                                                  |
|-------------------|--------------------------------------------------------------------|
| `VaultService`    | Criar, abrir, fechar, salvar e alterar senha do cofre              |
| `TreeService`     | Navegar, criar grupos, mover e reordenar nós                       |
| `ItemService`     | Criar, editar e remover itens e seus atributos                     |
| `SearchService`   | Busca full-text e filtrada na árvore                               |
| `TemplateService` | Listar e aplicar templates                                         |
| `ClipboardService`| Copiar atributos sensíveis e limpar área de transferência          |

### Domain Layer

Entidades e regras de negócio sem dependências externas. Testável de forma isolada.

- `Vault` — raiz do cofre, metadados, coleção de templates
- `Node` — nó da árvore, pode ser grupo ou item
- `Item` — dado pessoal com lista de atributos
- `Attribute` — par chave-valor com tipo semântico
- `Template` — modelo reutilizável de atributos

### Infrastructure Layer

- **Crypto**: derivação de chave com Argon2id + cifragem AES-256-GCM (ver `decisions/tdr-001-criptografia.md`)
- **Storage**: serialização JSON + escrita/leitura do arquivo `.abt` (ver `decisions/tdr-002-armazenamento.md`)
- **Clipboard**: integração com área de transferência do SO

## Estrutura de diretórios (proposta)

```
abditum/
├── cmd/
│   └── abditum/
│       └── main.go
├── internal/
│   ├── domain/
│   │   ├── vault.go
│   │   ├── node.go
│   │   ├── item.go
│   │   ├── attribute.go
│   │   └── template.go
│   ├── app/
│   │   ├── vault_service.go
│   │   ├── tree_service.go
│   │   ├── item_service.go
│   │   ├── search_service.go
│   │   └── clipboard_service.go
│   ├── crypto/
│   │   ├── kdf.go        ← Argon2id
│   │   └── cipher.go     ← AES-256-GCM
│   └── storage/
│       ├── encoder.go    ← JSON serialização
│       └── file.go       ← leitura/escrita do .abt
├── ui/
│   └── ...               ← depende do framework escolhido
├── specs/
│   └── ...
└── go.mod
```

## Fluxo de abertura do cofre

```
Usuário informa senha mestra
         ↓
Lê cabeçalho do arquivo (salt, nonce, parâmetros Argon2id)
         ↓
Deriva chave com Argon2id(senha, salt, params)  → 32 bytes
         ↓
Decifra payload com AES-256-GCM(chave, nonce, ciphertext)
         ↓
Valida tag GCM (detecta senha errada ou arquivo corrompido)
         ↓
Desserializa JSON → estrutura Vault em memória
         ↓
Exibe árvore na UI
```

## Fluxo de salvamento

```
Usuário realiza alteração
         ↓
Serializa Vault para JSON
         ↓
Gera novo nonce aleatório (12 bytes)
         ↓
Cifra com AES-256-GCM(chave em memória, nonce, json)
         ↓
Escreve header + ciphertext no arquivo atomicamente
(escreve em arquivo temporário, depois move/substitui)
```

## Decisões pendentes

- [ ] Framework de UI (ver `decisions/adr-002-ui-framework.md`)
- [ ] Formato de serialização interno: JSON vs MessagePack (ver `decisions/tdr-002-armazenamento.md`)
