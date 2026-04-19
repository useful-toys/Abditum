## Purpose

Definir a cobertura mínima e os cenários críticos de regressão para `internal/crypto`, com foco em fluxos AEAD/AAD, validações negativas e helpers determinísticos de memory locking.

## Requirements

### Requirement: Critical crypto validation paths SHALL be covered by automated tests
The project SHALL maintain automated tests for `internal/crypto` that exercise successful and failing paths of the package's public cryptographic helpers, including AES-GCM flows with and without AAD.

#### Scenario: AEAD helpers reject invalid parameters
- **WHEN** the test suite executes invalid-input cases for `Encrypt`, `Decrypt`, `SealWithAAD`, `EncryptWithAAD`, or `DecryptWithAAD`
- **THEN** each function MUST return the documented sentinel error for invalid keys, invalid nonces, short ciphertexts, or authentication failures

#### Scenario: AAD flows remain roundtrip-safe
- **WHEN** the test suite encrypts and decrypts payloads with authenticated additional data
- **THEN** the suite MUST verify successful roundtrip behavior and rejection of tampered ciphertext, tampered AAD, and invalid nonce inputs

#### Scenario: Test strategy avoids unnecessary mocks
- **WHEN** new tests are added to cover `internal/crypto`
- **THEN** they MUST prefer deterministic execution against real code paths and MUST only use mocks or test doubles when they are indispensable or when the no-mock alternative would be excessively complex, inelegant, or contrary to good coding practices

### Requirement: Crypto package coverage SHALL exceed the current baseline
The project SHALL keep statement coverage for `internal/crypto` at or above 90% when measured with the repository's Go test tooling for that package.

#### Scenario: Coverage check runs for internal crypto
- **WHEN** `go test -cover ./internal/crypto` is executed after the change
- **THEN** the reported statement coverage for `internal/crypto` MUST be at least 90%

### Requirement: Platform-specific memory locking helpers SHALL have deterministic tests where supported
The project SHALL include deterministic tests for memory-locking helpers that materially affect `internal/crypto` coverage on supported build targets, without changing the package's public API.

#### Scenario: Platform helper behavior is exercised
- **WHEN** tests run on a supported target with platform-specific memory locking code
- **THEN** the suite MUST verify the helper behavior for empty buffers and at least one non-empty buffer path
