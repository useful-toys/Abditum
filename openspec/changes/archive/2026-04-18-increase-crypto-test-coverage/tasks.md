## 1. Cobertura dos fluxos públicos de AEAD

- [x] 1.1 Mapear os ramos sem cobertura em `internal/crypto/aead.go` e alinhar os casos ausentes com os requisitos da spec.
- [x] 1.2 Adicionar testes para `SealWithAAD` cobrindo roundtrip válido, chave inválida e nonce inválido.
- [x] 1.3 Completar os testes negativos de `Encrypt`, `Decrypt`, `EncryptWithAAD` e `DecryptWithAAD` para validar erros documentados com entradas inválidas.

## 2. Cobertura de helpers internos e específicos de plataforma

- [x] 2.1 Adicionar testes internos em `package crypto` para helpers não exportados que bloqueiam cobertura relevante, sem alterar a API pública.
- [x] 2.2 Priorizar cenários determinísticos sem mocks; documentar e justificar qualquer double de teste apenas quando ele for indispensável ou evitar uma solução sem mock complexa, deselegante ou contrária a boas práticas.
- [x] 2.3 Introduzir seams privados apenas onde necessário para reproduzir de forma determinística falhas de entropy ou locking específicas de plataforma, evitando mocks amplos.
- [x] 2.4 Cobrir os caminhos de buffer vazio e pelo menos um caminho não vazio dos helpers de memory locking suportados pelo alvo atual.

## 3. Validação de cobertura e regressão

- [x] 3.1 Executar `go test ./internal/crypto` e corrigir eventuais quebras introduzidas pelos novos testes.
- [x] 3.2 Executar `go test -cover ./internal/crypto` e ajustar a suíte até atingir pelo menos 90% de cobertura de statements.
- [x] 3.3 Rodar a validação Go relevante do repositório para confirmar que a ampliação de cobertura não alterou o comportamento esperado do pacote.
