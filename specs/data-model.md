# Abditum — Modelo de Dados

## Estrutura do arquivo `.abt`

O arquivo possui dois segmentos sequenciais:

```
┌──────────────────────────────────────────┐
│  HEADER (plaintext)                      │
│  - magic bytes: "ABDITUM" (7 bytes)      │
│  - versão do formato: uint8              │
│  - algoritmo KDF: uint8 (1 = argon2id)   │
│  - argon2id time:    uint32              │
│  - argon2id memory:  uint32 (KiB)        │
│  - argon2id threads: uint8               │
│  - salt:  32 bytes aleatórios            │
│  - nonce: 12 bytes aleatórios            │
├──────────────────────────────────────────┤
│  PAYLOAD (AES-256-GCM cifrado)           │
│  - ciphertext (JSON serializado)         │
│  - tag GCM: 16 bytes (autenticação)      │
└──────────────────────────────────────────┘
```

> O campo `versão do formato` no header permite que versões futuras da aplicação detectem cofres mais antigos e apliquem migração se necessário. Como o payload é JSON, campos novos adicionados com `omitempty` são ignorados por versões anteriores sem quebrar a leitura — o formato tende a ser resiliente a evoluções incrementais.

## Entidades

### Vault

Raiz do arquivo. Contém metadados, templates e a raiz da árvore.

```json
{
  "version": 1,
  "metadata": {
    "name": "Meu Cofre",
    "created_at": "2026-03-02T10:00:00Z",
    "modified_at": "2026-03-02T10:00:00Z"
  },
  "templates": [ ...Template ],
  "root": { ...Node }
}
```

### Node

Nó da árvore — pode ser um **grupo** (pasta/diretório) ou um **item** (folha com dados).

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "name": "Banco do Brasil",
  "type": "item",
  "favorite": false,
  "item": { ...Item },
  "children": []
}
```

| Campo      | Tipo              | Obrigatório | Descrição                                       |
|------------|-------------------|-------------|--------------------------------------------------|
| `id`       | string (UUID v4)  | Sim         | Identificador único no cofre                     |
| `name`     | string            | Sim         | Nome exibido na árvore (não precisa ser único)   |
| `type`     | `group` \| `item` | Sim         | Determina o papel do nó                          |
| `item`     | Item              | Se `item`   | Dados do item; ausente quando `type = group`     |
| `favorite`  | bool              | Não         | Se `true`, o nó é marcado como favorito; padrão `false` |
| `children` | []Node            | Se `group`  | Filhos; ausente ou vazio quando `type = item`    |

### Item

Conjunto de atributos associado a um Node do tipo `item`.

```json
{
  "template_id": "bank-account",
  "icon": "bank",
  "notes": "Conta PF — agência centro",
  "created_at": "2026-03-02T10:00:00Z",
  "modified_at": "2026-03-02T10:00:00Z",
  "attributes": [
    { "key": "banco",          "label": "Banco",           "value": "Banco do Brasil", "type": "text",     "sensitive": false },
    { "key": "agencia",        "label": "Agência",         "value": "1234-5",          "type": "text",     "sensitive": false },
    { "key": "conta",          "label": "Conta Corrente",  "value": "98765-4",         "type": "text",     "sensitive": false },
    { "key": "senha_internet", "label": "Senha Internet",  "value": "s3cr3t",          "type": "password", "sensitive": true  },
    { "key": "senha_cartao",   "label": "Senha Cartão",    "value": "1234",            "type": "password", "sensitive": true  }
  ]
}
```

| Campo         | Tipo     | Obrigatório | Descrição                                          |
|---------------|----------|-------------|-----------------------------------------------------|
| `template_id` | string   | Não         | Referência ao template usado; `null` se personalizado |
| `icon`        | string   | Não         | Identificador de ícone (definido pela UI)            |
| `notes`       | string   | Não         | Texto livre de observações                           |
| `created_at`  | string (ISO 8601) | Sim  | Data e hora de criação do item                      |
| `modified_at` | string (ISO 8601) | Sim  | Data e hora da última alteração nos atributos ou nome |
| `attributes`  | []Attribute | Sim      | Lista ordenada de atributos                          |

### Attribute

Par chave-valor com tipo semântico.

| Campo       | Tipo       | Obrigatório | Descrição                                                     |
|-------------|------------|-------------|----------------------------------------------------------------|
| `key`       | string     | Sim         | Identificador do campo (único dentro do item)                  |
| `label`     | string     | Sim         | Nome legível exibido na UI                                     |
| `value`     | string     | Sim         | Valor (sempre armazenado como string)                          |
| `type`      | AttributeType | Sim      | Tipo semântico (ver tabela abaixo)                             |
| `sensitive` | bool       | Sim         | Se `true`, valor é mascarado na UI por padrão                  |

#### AttributeType

| Valor      | Descrição                                      |
|------------|------------------------------------------------|
| `text`     | Texto livre                                    |
| `password` | Texto sensível, mascarado por padrão           |
| `url`      | URL, exibida como link clicável                |
| `otp`      | Semente TOTP (RFC 6238)                        |
| `number`   | Valor numérico                                 |
| `date`     | Data (ISO 8601)                                |

### Template

Define um modelo reutilizável de atributos para facilitar a criação de itens.

```json
{
  "id": "website",
  "name": "Site",
  "icon": "globe",
  "builtin": true,
  "attribute_schemas": [
    { "key": "url",   "label": "URL",   "type": "url",      "sensitive": false },
    { "key": "login", "label": "Login", "type": "text",     "sensitive": false },
    { "key": "senha", "label": "Senha", "type": "password", "sensitive": true  }
  ]
}
```

| Campo               | Tipo             | Descrição                                       |
|---------------------|------------------|-------------------------------------------------|
| `id`                | string           | Identificador único                             |
| `name`              | string           | Nome exibido                                    |
| `icon`              | string           | Ícone padrão para itens criados com o template  |
| `builtin`           | bool             | `true` = template do sistema, não pode ser removido |
| `attribute_schemas` | []AttributeSchema | Definição dos campos padrão                    |

## Templates predefinidos

| ID             | Nome           | Atributos padrão                                          |
|----------------|----------------|-----------------------------------------------------------|
| `website`      | Site           | url, login, senha                                         |
| `bank-account` | Conta Bancária | banco, agência, conta, senha internet, senha cartão       |
| `credit-card`  | Cartão         | bandeira, número, validade, CVV, senha                    |
| `secure-note`  | Nota Segura    | conteúdo (textarea, sensitive)                            |
| `generic`      | Genérico       | (sem atributos predefinidos — totalmente livre)           |

## Invariantes e restrições

- `Node.id` deve ser único em todo o cofre
- `Node.type = group` → `Node.item` deve ser nulo; `Node.children` pode ser vazio
- `Node.type = item` → `Node.item` deve estar presente; `Node.children` deve ser vazio
- `Attribute.key` deve ser único dentro de um mesmo `Item`
- `Item.created_at` é definido na criação e nunca alterado
- `Item.modified_at` é atualizado ao salvar qualquer alteração nos atributos, nome ou notas do item; mover o nó na árvore **não** atualiza `modified_at`
- O nó raiz (`root`) é sempre do tipo `group` e não pode ser removido nem movido
- `Template.id` deve ser único na lista de templates do cofre
- Templates com `builtin = true` não podem ser removidos nem editados
