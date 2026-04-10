# Arquitetura — Pacote `internal/crypto`

## 1. Responsabilidade

O pacote `internal/crypto` provê todos os primitivos criptográficos necessários para o Abditum: derivação de chave a partir de senha, criptografia/descriptografia autenticada, gerenciamento seguro de memória e avaliação de força de senha. É uma camada de serviço técnico, sem dependência de nenhum outro pacote interno. Todos os outros pacotes que precisam de operações criptográficas consomem exclusivamente este pacote.

## 2. Algoritmos e Parâmetros

### 2.1 Derivação de Chave — Argon2id

O Argon2id (RFC 9106) é o algoritmo de derivação de chave. A escolha se justifica por duas propriedades combinadas: resistência a ataques de força bruta por GPU (herdada do Argon2i, que usa acesso de memória independente de dados) e resistência a ataques de canal lateral baseados em tempo (herdada do Argon2d, que usa acesso dependente de dados). Nenhum outro algoritmo da família Argon2 oferece as duas propriedades simultaneamente.

Os parâmetros adotados para o formato versão 1 são:

| Parâmetro | Valor | Significado |
|---|---|---|
| `t` (time) | 3 | Número de iterações |
| `m` (memory) | 262144 | Custo de memória em **KiB** (= 256 MiB) |
| `p` (parallelism) | 4 | Número de threads |
| `keyLen` | 32 | Tamanho da chave derivada em bytes (256 bits) |

A unidade do parâmetro de memória é **KiB, não bytes**. 256 MiB equivale a 262144 KiB. Usar o valor 256 no lugar de 262144 é um erro comum que reduz o custo de memória em um fator de 1024 — o código, a documentação e os testes repetem este aviso explicitamente.

O tempo esperado de derivação com esses parâmetros é de 200–500 ms em hardware moderno. Esse intervalo é intencional: lento o suficiente para resistir a ataques de dicionário, rápido o suficiente para ser imperceptível ao usuário durante a abertura do cofre.

### 2.2 Criptografia Autenticada — AES-256-GCM

O AES-256-GCM (Galois/Counter Mode) é o algoritmo de criptografia autenticada. A escolha oferece confidencialidade e autenticidade em um único primitivo — qualquer adulteração no texto cifrado ou nos dados adicionais autenticados (AAD) causa falha verificável na descriptografia.

O GCM exige um nonce único por par (chave, nonce). A reutilização de nonce com a mesma chave anula completamente a segurança do GCM e constitui a vulnerabilidade mais crítica da construção. Todo nonce é gerado com `crypto/rand` via `io.ReadFull`, garantindo entropia do sistema operacional.

## 3. Formatos de Saída da Criptografia

O pacote define dois formatos de saída distintos, usados em contextos diferentes:

### 3.1 Formato Padrão (Encrypt/Decrypt)

```
[ nonce: 12 bytes | ciphertext: N bytes | tag GCM: 16 bytes ]
```

O nonce é prefixado ao texto cifrado. O overhead total é de 28 bytes. O `Decrypt` extrai o nonce dos primeiros 12 bytes antes de descriptografar.

### 3.2 Formato com AAD (EncryptWithAAD / DecryptWithAAD / SealWithAAD)

Neste formato o nonce não é prefixado: é retornado separadamente e gravado no cabeçalho do arquivo `.abditum` (bytes 37–48). A saída de `EncryptWithAAD` é apenas `ciphertext + tag`. O cabeçalho completo de 49 bytes (que inclui o próprio nonce) é passado como AAD, de modo que qualquer alteração em qualquer byte do cabeçalho invalida a autenticação.

A variante `SealWithAAD` existe para o caso em que o nonce já foi gerado e gravado no cabeçalho antes da chamada de criptografia — o chamador fornece o nonce externamente em vez de ter o pacote gerá-lo internamente.

## 4. Gerenciamento do Sal

O sal tem 32 bytes, gerado com `io.ReadFull(rand.Reader, ...)`. O uso de `io.ReadFull` em vez de `rand.Read` diretamente é intencional: garante que exatamente 32 bytes sejam preenchidos, já que `rand.Read` pode retornar menos bytes em sistemas sob pressão de entropia.

O sal é gerado apenas em dois eventos:
- Criação de um novo cofre
- Alteração da senha mestra

Uma vez criado, o sal é armazenado em texto claro no cabeçalho do arquivo `.abditum` e reutilizado a cada abertura do cofre. Sal e senha juntos determinam a chave derivada — um sal diferente produz uma chave diferente mesmo com a mesma senha.

## 5. Erros Sentinela

O pacote define quatro erros sentinela:

| Erro | Semântica |
|---|---|
| `ErrAuthFailed` | Falha de autenticação GCM: chave incorreta **ou** dado corrompido |
| `ErrInsufficientEntropy` | `crypto/rand` não conseguiu prover bytes aleatórios |
| `ErrInvalidParams` | Parâmetros inválidos passados ao chamador (erro de programação) |
| `ErrMLockFailed` | Bloqueio de memória não disponível ou falhou (não fatal) |

A decisão de retornar `ErrAuthFailed` para **ambos** os casos "chave errada" e "dado corrompido" é deliberada: expor erros distintos permitiria que um atacante distinguisse os dois cenários via análise de tempo, o que constitui um canal lateral. O chamador deve tratar os dois casos da mesma forma.

`ErrMLockFailed` é o único erro não fatal. O pacote continua operando normalmente mesmo quando o bloqueio de memória não está disponível. O chamador deve registrar um aviso e continuar.

## 6. Segurança de Memória

### 6.1 Princípio: dados sensíveis como `[]byte`

Todo dado sensível (senha, chave derivada, texto simples) é representado como `[]byte`, nunca como `string`. Strings em Go são imutáveis e não podem ser zeradas — uma string contendo a senha mestra ficaria acessível na memória durante toda a vida do processo e potencialmente em core dumps ou swap. Slices de bytes permitem zeragem explícita.

### 6.2 Zeragem explícita com `Wipe`

`Wipe` percorre o slice byte a byte e atribui zero a cada posição. A chamada a `runtime.KeepAlive(data)` após o loop impede que o compilador elimine o zeroing como dead code por otimização. A zeragem é responsabilidade do chamador — o pacote não zera buffers automaticamente para dar ao chamador controle total sobre o tempo de vida dos dados.

### 6.3 Alocação segura com `SecureAllocate`

`SecureAllocate` combina três operações em uma chamada:
1. Alocação de um slice zerado do tamanho solicitado
2. Tentativa de bloqueio da memória (`mlock`/`VirtualLock`) para impedir swap para disco
3. Retorno de uma função `cleanup` que zera o buffer e desbloqueia a memória quando chamada

O erro de bloqueio é retornado junto com o buffer, não no lugar dele. O buffer é utilizável independentemente de o bloqueio ter funcionado — o chamador decide se a falha de bloqueio é aceitável para o contexto.

### 6.4 Bloqueio de memória por plataforma

A função interna `mlock` é implementada em três arquivos separados por build tags:

- `mlock_unix.go` (`!windows`): chama `unix.Mlock` via `golang.org/x/sys/unix`
- `mlock_windows.go` (`windows`): chama `windows.VirtualLock` via `golang.org/x/sys/windows`
- `mlock_other.go` (`!unix && !windows`): retorna `ErrMLockFailed` imediatamente

A simetria é mantida: `munlock` existe nas três plataformas, chamando o desbloqueio correspondente ou sendo no-op.

## 7. Avaliação de Força de Senha

`EvaluatePasswordStrength` classifica senhas em dois níveis: `StrengthWeak` e `StrengthStrong`. Não há gradações intermediárias — a avaliação é binária.

Os critérios para `StrengthStrong` são cumulativos:
- Comprimento mínimo de 12 caracteres
- Presença de pelo menos um caractere maiúsculo (A–Z)
- Presença de pelo menos um caractere minúsculo (a–z)
- Presença de pelo menos um dígito (0–9)
- Presença de pelo menos um caractere especial (qualquer outro byte)

A função opera diretamente sobre `[]byte` sem conversão para `string`, pela mesma razão que o restante do pacote: evitar cópias não zeráveis de dados sensíveis em memória.

## 8. Versionamento de Formato

A constante `FormatVersion = 1` identifica o conjunto de parâmetros criptográficos e a estrutura do arquivo `.abditum`. Ela é gravada no cabeçalho do arquivo e lida na abertura do cofre para selecionar o conjunto correto de parâmetros. Futuras versões do formato podem adotar algoritmos ou parâmetros distintos sem quebrar a compatibilidade retroativa com cofres existentes.

O tipo `ArgonParams` é uma estrutura de dados pura que parametriza a chamada ao Argon2id. Utilizá-la em vez de constantes internas permite que versões futuras do formato passem parâmetros distintos sem alterar a assinatura de `DeriveKey`.

## 9. Convenções de Testes

Os testes estão no pacote `crypto_test` (caixa preta), não em `crypto` (caixa branca). Esta convenção garante que os testes exercitem exclusivamente a API pública, validando contratos externos em vez de detalhes internos de implementação.

Cada arquivo fonte tem um arquivo de teste correspondente. O arquivo `crypto_test.go` cobre exclusivamente testes de integração que percorrem o fluxo completo: `GenerateSalt` → `DeriveKey` → `Encrypt` → `Decrypt` → `Wipe`.

Há um teste dedicado à unicidade de nonces (`TestNonceUniqueness`) que verifica que duas encriptações do mesmo texto com a mesma chave produzem textos cifrados distintos. Este teste verifica a propriedade mais crítica para a segurança do GCM.

Benchmarks de `DeriveKey` estão incluídos para verificar empiricamente que o custo computacional de derivação permanece na faixa de 200–500 ms com os parâmetros de produção.
