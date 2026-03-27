# Especificação do Formato de Arquivo — `.abditum`

## 1. Estrutura do Arquivo

Stream binário: cabeçalho fixo não-criptografado (49 bytes) + payload criptografado.

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

| Campo            | Offset | Tamanho  | Tipo              | Descrição |
|------------------|--------|----------|-------------------|-----------|
| `magic`          | 0      | 4 bytes  | bytes (ASCII)     | Assinatura fixa `ABDT`. Permite rejeitar arquivos inválidos antes de solicitar senha. |
| `versão_formato` | 4      | 1 byte   | uint8, big-endian | Seleciona o perfil criptográfico completo (KDF + AEAD) e o esquema de migração do payload. |
| `salt`           | 5      | 32 bytes | bytes             | Salt único por cofre; gerado via CSPRNG na criação, substituído apenas ao alterar a senha mestra. |
| `nonce`          | 37     | 12 bytes | bytes             | IV do AES-256-GCM (padrão NIST). Regenerado via CSPRNG a cada salvamento. |

O payload contém o `Cofre` serializado em JSON (UTF-8) — todos os valores de string são UTF-8 — seguido pela tag GCM de 16 bytes.

---

## 2. AAD

O cabeçalho completo (`magic` + `versão_formato` + `salt` + `nonce`) é autenticado como AAD do AES-256-GCM, sem necessidade de checksum adicional.

---

## 3. Derivação de Chave

A chave AES-256 (32 bytes) é derivada da senha mestra com **Argon2id**, usando o perfil criptográfico completo (KDF + AEAD) definido pelo `versão_formato`. O Argon2id cuida exclusivamente da derivação de chave; a integridade do conteúdo é garantida pela tag GCM.

### Perfis por versão de formato

| `versão_formato` | KDF      | AEAD        | `m`                  | `t` | `p` | Saída    |
|------------------|----------|-------------|----------------------|-----|-----|----------|
| 1                | Argon2id | AES-256-GCM | 262144 KB (256 MiB)  | 3   | 4   | 32 bytes |

Versões futuras serão documentadas aqui quando definidas.

### Política de Parametrização (v1)

- Parâmetros fixos e hard-coded; sem calibração por máquina nem variação por arquivo.
- Piso: `m` ≥ 131072 KB (128 MiB). Teto de referência: `m` ≤ 524288 KB (512 MiB).
- Mesma política em Windows, macOS e Linux 64 bits.
- Mudanças exigem decisão explícita de versão + testes de regressão.
- O perfil é selecionado **exclusivamente** pelo `versão_formato` — nunca por heurística.

### Geração de valores aleatórios

`salt` e `nonce` **DEVEM** ser gerados a partir de fonte criptograficamente segura (CSPRNG do sistema operacional — em Go, `crypto/rand`). O uso de geradores não-criptográficos (como `math/rand`) compromete completamente a segurança do esquema.

---

## 4. Sequência de Abertura

1. `magic` inválido → rejeitar sem solicitar senha
2. `versão_formato` > suportado → rejeitar com erro de incompatibilidade; senão selecionar perfil Argon2id
3. Ler `salt` e `nonce`
4. Solicitar senha mestra
5. Derivar chave com Argon2id
6. Descriptografar e autenticar com AES-256-GCM (AAD = cabeçalho) — falha: erro de autenticação
7. Desserializar JSON para estrutura de domínio em memória
8. Validar modelo: Pasta Geral deve existir com nome `"Geral"` — falha: erro de integridade
9. Se formato histórico: migrar dados em memória para o modelo corrente
10. Cofre disponível

### Categorias de erro

| Categoria           | Condição                                                           | Comportamento |
|---------------------|--------------------------------------------------------------------|---------------|
| Tipo de arquivo     | `magic` inválido                                                   | Rejeitar; sem nova tentativa |
| Versão incompatível | `versão_formato` > suportado                                       | Rejeitar; sem nova tentativa |
| Autenticação        | Tag GCM inválida                                                   | Mensagem genérica; permitir nova tentativa |
| Integridade         | Payload corrompido, JSON inválido ou Pasta Geral ausente/inválida  | Mensagem genérica; sem nova tentativa |

Mensagens ao usuário são sempre genéricas, sem detalhes técnicos.

---

## 5. Escrita e Compatibilidade

A escrita sempre usa o formato atual, com novo `nonce` gerado via CSPRNG. O `salt` é gerado na criação e substituído apenas ao alterar a senha mestra.

A escrita **DEVE** ser atômica (gravar em arquivo temporário e renomear) para evitar corrupção do cofre em caso de falha durante a gravação (crash, queda de energia).

A aplicação abre qualquer versão anterior suportada, migrando o payload em memória. Salva sempre no formato atual. Versões superiores ao suportado são rejeitadas com erro claro. Mudanças de formato devem ser raras e acompanhadas de rotina de migração e testes de regressão para todos os formatos históricos suportados.
