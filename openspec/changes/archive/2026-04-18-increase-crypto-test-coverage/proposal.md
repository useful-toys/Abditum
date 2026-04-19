## Why

A cobertura atual de `internal/crypto` estagnou em 74.1%, deixando sem proteção automatizada alguns ramos importantes de erro e validação em primitivas criptográficas. Como esse pacote concentra operações sensíveis de criptografia, derivação de chave e limpeza de memória, vale fechar essas lacunas agora para reduzir regressões em caminhos críticos.

## What Changes

- Ampliar os testes de `internal/crypto` para cobrir ramos hoje sem cobertura, com foco em `SealWithAAD` e em validações negativas de AEAD/AAD.
- Exercitar cenários de falha e de borda que hoje não garantem comportamento estável para chaves, nonces e ciphertexts inválidos.
- Priorizar testes reais e determinísticos, evitando mocks salvo quando eles forem indispensáveis ou quando a alternativa sem mock introduzir complexidade, deselegância ou violação de boas práticas.
- Formalizar uma expectativa mínima de cobertura e de regressão para o pacote `internal/crypto` usando os testes automatizados existentes.

## Capabilities

### New Capabilities
- `crypto-test-coverage`: Define a cobertura esperada e os cenários críticos de regressão para o pacote `internal/crypto`.

### Modified Capabilities
- Nenhuma.

## Impact

- Código afetado: `internal/crypto/aead.go`, `internal/crypto/kdf.go`, `internal/crypto/mlock_windows.go` e arquivos de teste relacionados em `internal/crypto/*_test.go`.
- Fluxos afetados: criptografia AES-GCM com e sem AAD, validação de parâmetros, geração de sal e wrappers de locking em memória.
- Validação afetada: `go test -cover ./internal/crypto` e a suíte Go que já cobre o pacote.
