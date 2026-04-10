# Especificação do Formato de Arquivo — `.abditum`

## 1. Estrutura do Arquivo

Stream binário com estrutura: cabeçalho fixo não-criptografado (49 bytes) + payload criptografado.

| Campo            | Offset | Tamanho  | Tipo              | Descrição |
|------------------|--------|----------|-------------------|-----------|
| `magic`          | 0      | 4 bytes  | bytes (ASCII)     | Assinatura fixa `ABDT`. |
| `versão_formato` | 4      | 1 byte   | uint8, big-endian | Seleciona o perfil criptográfico completo (KDF + AEAD) e o esquema de migração do payload. |
| `salt`           | 5      | 32 bytes | bytes             | Salt único por cofre. |
| `nonce`          | 37     | 12 bytes | bytes             | IV do AES-256-GCM conforme o padrão NIST. |
| payload          | 49–…   | variável | bytes             | Payload criptografado com AES-256-GCM contendo JSON da estrutura Cofre (UTF-8) + tag GCM (16 bytes). |

## 2. AAD

O cabeçalho completo (`magic` + `versão_formato` + `salt` + `nonce`) é autenticado como AAD do AES-256-GCM, sem necessidade de checksum adicional.

## 3. Derivação de Chave

A chave AES-256 (32 bytes) é derivada da senha mestra com **Argon2id**, usando o perfil criptográfico completo (KDF + AEAD) definido pelo `versão_formato`. O Argon2id é responsável exclusivamente pela derivação de chave; a integridade do conteúdo é garantida pela tag GCM.

### Perfis por versão de formato

| `versão_formato` | KDF      | AEAD        | `m`                  | `t` | `p` | Saída    |
|------------------|----------|-------------|----------------------|-----|-----|----------|
| 1                | Argon2id | AES-256-GCM | 262144 KB (256 MiB)  | 3   | 4   | 32 bytes |

Versões futuras serão documentadas aqui conforme forem definidas.

### Salt e Nonce — Geração e Ciclo de Vida

`salt` e `nonce` **DEVEM** ser gerados a partir de fonte criptograficamente segura (CSPRNG do sistema operacional; em Go, `crypto/rand`). O uso de geradores não-criptográficos (como `math/rand`) quebra completamente a segurança do esquema.

**Ciclo de vida:**
- **salt:** Gerado via CSPRNG na criação do arquivo. Alterado apenas quando a senha mestra é modificada (rekeying).
- **nonce:** Regenerado via CSPRNG a cada salvamento.

## 4. Sequência de Abertura

1. `magic` inválido → rejeitar sem solicitar senha
2. `versão_formato` > suportado → rejeitar com erro de incompatibilidade; caso contrário, selecionar perfil Argon2id
3. Ler `salt` e `nonce`
4. Solicitar senha mestra
5. Derivar chave com Argon2id
6. Descriptografar e autenticar com AES-256-GCM (AAD = cabeçalho) — falha: erro de autenticação
7. Desserializar JSON para estrutura de domínio em memória
8. Validar modelo: O JSON deve atender à estrutura exigida pela versão (campos obrigatórios, elementos obrigatórios); falha: erro de integridade
9. Se formato histórico: migrar dados em memória para o modelo corrente
10. Cofre pronto para uso

### Categorias de erro

| Categoria           | Condição                                                           | Comportamento |
|---------------------|--------------------------------------------------------------------|---------------|
| Tipo de arquivo     | `magic` inválido                                                   | Rejeitar imediatamente |
| Versão incompatível | `versão_formato` > suportado                                       | Rejeitar imediatamente |
| Autenticação        | Tag GCM inválida                                                   | Mensagem genérica; permitir nova tentativa |
| Integridade         | Payload corrompido, JSON inválido ou Pasta Geral ausente/inválida  | Mensagem genérica; sem nova tentativa |

Mensagens ao usuário são sempre genéricas, sem detalhes técnicos.

## 5. Escrita e Compatibilidade

A escrita sempre usa o formato mais recente, com novo `nonce` gerado via CSPRNG. O `salt` permanece do arquivo original, sendo modificado apenas em rekeying (alteração de senha mestra).

A escrita **DEVE** ser atômica (gravar em arquivo temporário e renomear) para evitar corrupção do cofre em caso de falha durante a gravação (travamento do programa, interrupção de energia).

## 6. Política de Versioning

Os parâmetros criptográficos (Argon2id, KDF, AEAD) são fixos e hard-coded para cada versão de `versão_formato`, sem calibração por máquina nem variação por arquivo.

**Quando incrementar `versão_formato`:**
- Mudanças nos parâmetros Argon2id (`m`, `t`, `p`)
- Mudanças no KDF ou no esquema AEAD
- Mudanças significativas no formato do JSON que não possam ser tratadas transparentemente pelo parser (ex: campos obrigatórios adicionados, remoção de campos, mudanças de tipo)

**Mudanças que NÃO exigem nova versão:**
- Campos opcionais adicionados ao JSON
- Novos tipos de dados adicionais
- Aplicação desserializa tolerando extensões

**Impacto de mudanças de versão:**
- Mudanças devem ser raras
- Devem ser acompanhadas de rotina de migração para formatos históricos suportados
- Testes de regressão obrigatórios para verificar compatibilidade com todas as versões anteriores
- Versões superiores ao suportado pela aplicação são rejeitadas com mensagem clara de erro

## 7. Escrita e Compatibilidade

A aplicação abre qualquer versão anterior suportada, migrando o payload em memória, e salva sempre no formato atual (conforme política de versioning definida na Seção 6).
