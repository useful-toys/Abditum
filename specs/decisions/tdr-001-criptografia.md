# TDR 001 — Esquema de Criptografia

**Status**: Aceita
**Data**: 2026-03-02

## Contexto

O cofre precisa ser protegido contra acesso não autorizado. O arquivo é armazenado localmente e deve resistir a ataques de força bruta e dicionário contra a senha mestra. A integridade do arquivo também deve ser verificável — adulteração deve ser detectada.

## Decisão

### Derivação de chave: Argon2id

- Algoritmo: **Argon2id** (vencedor do Password Hashing Competition, 2015; resistente a ataques GPU e side-channel)
- Parâmetros padrão (ajustáveis via header):
  - `time = 3` iterações
  - `memory = 65536` KiB (64 MiB)
  - `parallelism = 4`
  - `keyLen = 32` (256 bits)
- Salt: 32 bytes aleatórios, gerados com `crypto/rand`, armazenados em plaintext no header

### Cifragem: AES-256-GCM

- Algoritmo: **AES-256-GCM** (AEAD — fornece cifragem + autenticação integrada)
- Chave: 256 bits derivados pelo Argon2id
- Nonce: 12 bytes aleatórios, gerados com `crypto/rand` a cada salvamento
- Tag GCM: 16 bytes, garante integridade — qualquer adulteração do ciphertext é detectada

### Bibliotecas Go

```go
golang.org/x/crypto/argon2   // Argon2id
crypto/aes                   // stdlib — AES
crypto/cipher                // stdlib — GCM
crypto/rand                  // stdlib — geração de bytes aleatórios
```

## Alternativas consideradas

| Alternativa      | Motivo para rejeição                                                 |
|------------------|----------------------------------------------------------------------|
| bcrypt           | Limitado a 72 bytes de senha; não deriva chave de tamanho arbitrário |
| PBKDF2           | Menos resistente a ataques de hardware dedicado que Argon2id         |
| scrypt           | Argon2id é preferido por resistência a side-channel adicional        |
| ChaCha20-Poly1305| Igualmente válido; AES-GCM preferido por suporte a AES-NI em hardware|

## Consequências

- Abertura do cofre demora ~0.5–2 segundos intencionalmente (custo do Argon2id)
- Parâmetros KDF no header permitem migração futura sem quebrar cofres existentes
- Senha errada é indistinguível de arquivo corrompido (sem oracle de senha)
- Perda da senha mestra = perda permanente dos dados (sem recuperação)
- Novo nonce a cada salvamento previne reutilização de nonce com a mesma chave
