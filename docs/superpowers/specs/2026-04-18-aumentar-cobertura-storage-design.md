# design: Aumentar Cobertura de Testes internal/storage para 90%

## Contexto

- **Pacote**: `internal/storage`
- **Cobertura atual**: 79.6%
- **Meta**: 90%
- **Data**: 2026-04-18

## Funções com Baixa Cobertura

| Função | Cobertura | Caminho |
|--------|----------|---------|
| `Save` | 65.9% | storage.go:93 |
| `SaveNew` | 75.0% | storage.go:28 |
| `Load` | 90.0% | storage.go:178 |
| `Salvar` | 75.0% | repository.go:101 |
| `ComputeFileMetadata` | 75.0% | detect.go:64 |
| `DetectExternalChange` | 91.7% | detect.go:28 |
| `atomicRename` | 71.4% | atomic_rename_windows.go:16 |
| `RecoverOrphans` | 80.0% | recover.go:26 |
| `readSaltFromFile` | 75.0% | repository.go:198 |

## Abordagem

Prioridade: **testes de integração** (arquivos reais), fallback para **mocks** quando integração fica complexo/confuso.

Edge cases a testar via integração:

1. `Save` — falha de write após .tmp criado (testa rollback)
2. `SaveNew` — erro de permissão (criar arquivo só-leitura)
3. `Load` — arquivo corrupto no payload (payload não é JSON válido)
4. `Salvar` — cobertura path isNew=false com update
5. `ComputeFileMetadata` — arquivo vazio (0 bytes)
6. `RecoverOrphans` — testar .bak e .bak2 stale
7. `atomicRename` — testar falha em Windows (arquivo em uso)

## Critérios deSucesso

- Cobertura >= 90% para package internal/storage
- Todos novos testes via integração (arquivos reais)
- Mocks apenas se integração não viável
- Testes não devem quebrar existentes

## Estimativa

~15-20 novos testes de integração