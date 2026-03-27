# Especificação do Formato de Arquivo — `.abditum`

## 1. Visão Geral

O arquivo do cofre Abditum possui extensão `.abditum` e é um stream binário composto por duas partes sequenciais: um **cabeçalho de criptografia** em formato fixo (não criptografado) e um **payload criptografado** contendo a estrutura do cofre em JSON.

---

## 2. Estrutura do Arquivo

```
┌──────────┬─────────────────────────────────────────────────┐
│  Offset  │  Campo                                          │
├──────────┼─────────────────────────────────────────────────┤
│   0– 3   │  magic          (4 bytes, ASCII "ABDT")         │
│   4– 4   │  versão_formato (1 byte, uint8, big-endian)     │
│   5–36   │  salt           (32 bytes)                      │
│  37–48   │  nonce          (12 bytes)                      │
├──────────┴─────────────────────────────────────────────────┤
│  49–…    │  payload criptografado (AES-256-GCM)            │
│          │    = JSON da estrutura Cofre (UTF-8)            │
│          │    + tag GCM de autenticação (16 bytes, fim)    │
└──────────────────────────────────────────────────────────┘
```

Tamanho total do cabeçalho: **49 bytes**.

### 2.1. Cabeçalho de Criptografia

| Campo            | Offset  | Tamanho | Tipo             | Descrição |
|------------------|---------|---------|------------------|-----------|
| `magic`          | 0       | 4 bytes | bytes (ASCII)    | Assinatura fixa `ABDT`. Identifica o arquivo como um cofre Abditum antes de qualquer tentativa de descriptografia. |
| `versão_formato` | 4       | 1 byte  | uint8, big-endian | Versão do formato do arquivo. Determina o perfil Argon2id a utilizar na derivação de chave e o esquema de migração do payload. |
| `salt`           | 5       | 32 bytes | bytes            | Salt único por arquivo, usado na derivação da chave com Argon2id. |
| `nonce`          | 37      | 12 bytes | bytes            | Vetor de inicialização (IV) para AES-256-GCM (padrão NIST). Deve ser único por operação de escrita — um novo nonce é gerado a cada salvamento. |

### 2.2. Payload Criptografado

O restante do stream é o payload, que contém a estrutura do `Cofre` serializada em JSON (UTF-8), criptografada com AES-256-GCM usando a chave derivada de Argon2id.

---

## 3. AAD (Additional Authenticated Data)

Todo o cabeçalho de criptografia (`magic` + `versão_formato` + `salt` + `nonce`) é incluído como AAD do AES-256-GCM. Isso garante que a integridade do cabeçalho seja validada criptograficamente junto ao payload, sem necessidade de checksum adicional da aplicação.

---

## 4. Derivação de Chave

A chave AES-256 (32 bytes) é derivada da senha mestra usando **Argon2id**. O perfil de parâmetros é selecionado com base no campo `versão_formato` do cabeçalho, permitindo que futuras versões adotem parâmetros mais robustos sem quebrar compatibilidade com arquivos antigos.

### Perfis por versão de formato

| `versão_formato` | Algoritmo | `m`     | `t` | `p` | Saída  |
|------------------|-----------|---------|-----|-----|--------|
| 1                | Argon2id  | 65536 KB (64 MB) | 3 | 4 | 32 bytes |

---

## 5. Sequência de Abertura

A aplicação segue esta ordem ao abrir um arquivo `.abditum`:

1. Ler `magic` — se inválido, rejeitar imediatamente com erro de tipo de arquivo (não solicitar senha)
2. Ler `versão_formato` — selecionar o perfil Argon2id correspondente
3. Ler `salt` e `nonce` do cabeçalho
4. Solicitar senha mestra ao usuário
5. Derivar a chave com Argon2id usando `salt` e o perfil selecionado
6. Descriptografar e autenticar o payload com AES-256-GCM (usando o cabeçalho como AAD)
   - Se a autenticação falhar: erro de senha incorreta ou integridade comprometida
7. Desserializar o JSON do payload para a estrutura de domínio em memória
8. Se `versão_formato` indicar formato histórico suportado: migrar os dados em memória para o modelo corrente
9. Cofre disponível para uso

### Categorias de erro na abertura

| Categoria       | Condição                                              | Comportamento |
|-----------------|-------------------------------------------------------|---------------|
| Tipo de arquivo | `magic` inválido ou ausente                           | Rejeitar antes de solicitar senha; sem nova tentativa |
| Autenticação    | Tag GCM inválida (senha incorreta)                    | Exibir mensagem genérica; permitir nova tentativa |
| Integridade     | Payload corrompido, JSON inválido, Pasta Geral ausente | Exibir mensagem genérica; bloquear abertura sem nova tentativa |

Em todos os casos, a mensagem exibida ao usuário é genérica e não revela detalhes técnicos sobre a falha.

---

## 6. Escrita

A aplicação sempre escreve o arquivo no formato da versão atual, atualizando `versão_formato` quando necessário e gerando um novo `nonce` a cada salvamento. O `salt` é gerado uma única vez na criação do cofre e substituído apenas quando a senha mestra é alterada.

---

## 7. Compatibilidade Retroativa

A aplicação de versão N é capaz de abrir arquivos criados em qualquer versão anterior do formato suportada. Ao abrir um arquivo antigo, o payload é migrado em memória para o modelo corrente do domínio. Ao salvar, o arquivo é sempre regravado no formato da versão atual da aplicação.
