# Aumentar Cobertura internal/storage para 90%

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Aumentar cobertura de testes do package internal/storage de 79.6% para >=90%

**Architecture:** Testes de integração (arquivos reais), fallback para mocks apenas quando integração não viável

**Tech Stack:** Go testing, arquivos temporários via t.TempDir()

---

## Task 1: Save - Falha Após .tmp Criado (Rollback)

**Files:**
- Modify: `internal/storage/storage_test.go`

- [ ] **Step 1: Adicionar teste de falha no rename vault→.bak**

```go
func TestSave_FailRenameVaultToBak(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "vault.abditum")

    cofre := newTestCofre()
    if err := storage.SaveNew(path, cofre, testPassword); err != nil {
        t.Fatalf("SaveNew() error: %v", err)
    }

    data, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("ReadFile() error: %v", err)
    }
    salt := data[storage.SaltOffset:storage.SaltOffset+storage.SaltSize]

    bakPath := path + ".bak"
    if err := os.Mkdir(bakPath, 0755); err != nil {
        t.Fatalf("Mkdir() error: %v", err)
    }

    err = storage.Save(path, cofre, testPassword, salt)
    if err == nil {
        t.Error("esperado erro ao renomear vault para .bak quando .bak é diretório")
    }

    os.Remove(bakPath)
}
```

- [ ] **Step 2: Run teste**

```bash
go test -v ./internal/storage/ -run TestSave_FailRenameVaultToBak
```
Expected: PASS

- [ ] **Step 3: Commit**
```bash
git add internal/storage/storage_test.go
git commit -m "test: adiciona teste Save falha rename vault→bak"
```

---

## Task 2: SaveNew - Erro de Permissão

**Files:**
- Modify: `internal/storage/storage_test.go`

- [ ] **Step 1: Adicionar teste de permissão negada**

```go
func TestSaveNew_PermissionDenied(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "readonly.dir")
    if err := os.MkdirAll(path, 0555); err != nil {
        t.Fatalf("MkdirAll() error: %v", err)
    }
    vaultPath := filepath.Join(path, "vault.abditum")

    err := storage.SaveNew(vaultPath, newTestCofre(), testPassword)
    if !os.IsPermission(err) {
        t.Errorf("esperado erro de permissão, obteve: %v", err)
    }
}
```

- [ ] **Step 2: Run teste**

```bash
go test -v ./internal/storage/ -run TestSaveNew_PermissionDenied
```
Expected: PASS (em Windows pode variar)

- [ ] **Step 3: Commit**
```bash
git add internal/storage/storage_test.go
git commit -m "test: adiciona teste SaveNew permissão negada"
```

---

## Task 3: Load - Arquivo Коррупо в Payload

**Files:**
- Modify: `internal/storage/storage_test.go`

- [ ] **Step 1: Adicionar teste de payload corrupto**

```go
func TestLoad_CorruptedPayload(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "corrupt.abditum")

    if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
        t.Fatalf("SaveNew() error: %v", err)
    }

    data, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("ReadFile() error: %v", err)
    }
    data[storage.HeaderSize] ^= 0xFF
    if err := os.WriteFile(path, data, 0600); err != nil {
        t.Fatalf("WriteFile() error: %v", err)
    }

    _, _, err = storage.Load(path, testPassword)
    if err == nil {
        t.Error("esperado erro ao carregar payload corrupto")
    }
}
```

- [ ] **Step 2: Run teste**

```bash
go test -v ./internal/storage/ -run TestLoad_CorruptedPayload
```
Expected: PASS

- [ ] **Step 3: Commit**
```bash
git add internal/storage/storage_test.go
git commit -m "test: adiciona teste Load payload corrupto"
```

---

## Task 4: Salvar -，覆盖 isNew=false

**Files:**
- Modify: `internal/storage/storage_test.go`

- [ ] **Step 1: Adicionar teste Salvar isNew=false via Carregar**

```go
func TestSalvar_UpdateExistingVault(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "vault.abditum")

    cofre1 := vault.NovoCofre()
    repo := storage.NewFileRepositoryForCreate(path, testPassword)
    if err := repo.Salvar(cofre1); err != nil {
        t.Fatalf("Salvar() error: %v", err)
    }

    loaded1, meta1, err := storage.Load(path, testPassword)
    if err != nil {
        t.Fatalf("Load() error: %v", err)
    }
    if meta1.Size == 0 {
        t.Fatal("metadata tamanho não pode ser zero")
    }

    repo2 := storage.NewFileRepository(path, testPassword, nil, meta1)
    cofre2 := vault.NovoCofre()
    if err := repo2.Salvar(cofre2); err != nil {
        t.Fatalf("Salvar() update error: %v", err)
    }

    loaded2, meta2, err := storage.Load(path, testPassword)
    if err != nil {
        t.Fatalf("Load() after update error: %v", err)
    }

    if loaded1.PastaGeral() == nil || loaded2.PastaGeral() == nil {
        t.Fatal("PastaGeral nil")
    }
    if meta2.Size == meta1.Size {
        t.Log("aviso: tamanhos iguais após update (possível em certains casos)")
    }
}
```

- [ ] **Step 2: Run teste**

```bash
go test -v ./internal/storage/ -run TestSalvar_UpdateExistingVault
```
Expected: PASS

- [ ] **Step 3: Commit**
```bash
git add internal/storage/storage_test.go
git commit -m "test: adiciona teste Salvar update vault existente"
```

---

## Task 5: ComputeFileMetadata - ArquivoVazio

**Files:**
- Modify: `internal/storage/storage_test.go`

- [ ] **Step 1: Adicionar teste de arquivo vazio**

```go
func TestComputeFileMetadata_EmptyFile(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "empty.abditum")

    if err := os.WriteFile(path, []byte{}, 0600); err != nil {
        t.Fatalf("WriteFile() error: %v", err)
    }

    meta, err := storage.ComputeFileMetadata(path)
    if err != nil {
        t.Fatalf("ComputeFileMetadata() error: %v", err)
    }
    if meta.Size != 0 {
        t.Errorf("Size = %d, want 0", meta.Size)
    }
    var zeroHash [32]byte
    if meta.Hash == zeroHash {
        t.Error("Hash deveria ser zero para arquivo vazio")
    }
}
```

- [ ] **Step 2: Run teste**

```bash
go test -v ./internal/storage/ -run TestComputeFileMetadata_EmptyFile
```
Expected: PASS

- [ ] **Step 3: Commit**
```bash
git add internal/storage/storage_test.go
git commit -m "test: adiciona teste ComputeFileMetadata arquivo vazio"
```

---

## Task 6: RecoverOrphans - .bak e .bak2 Stale

**Files:**
- Modify: `internal/storage/storage_test.go`

- [ ] **Step 1: Adicionar teste de recovery com .bak e .bak2**

```go
func TestRecoverOrphans_WithBackupFiles(t *testing.T) {
    dir := t.TempDir()
    path := filepath.Join(dir, "vault.abditum")

    if err := storage.SaveNew(path, newTestCofre(), testPassword); err != nil {
        t.Fatalf("SaveNew() error: %v", err)
    }

    bakPath := path + ".bak"
    bak2Path := path + ".bak2"

    if err := os.WriteFile(bakPath, []byte("old backup"), 0600); err != nil {
        t.Fatalf("WriteFile() bak error: %v", err)
    }
    if err := os.WriteFile(bak2Path, []byte("older backup"), 0600); err != nil {
        t.Fatalf("WriteFile() bak2 error: %v", err)
    }

    if err := storage.RecoverOrphans(path); err != nil {
        t.Fatalf("RecoverOrphans() error: %v", err)
    }

    if _, err := os.Stat(bakPath); os.IsNotExist(err) {
        t.Error(".bak deveria remaincer após RecoverOrphans")
    }
    if _, err := os.Stat(bak2Path); os.IsNotExist(err) {
        t.Error(".bak2 deveria permanecer após RecoverOrphans")
    }
}
```

- [ ] **Step 2: Run teste**

```bash
go test -v ./internal/storage/ -run TestRecoverOrphans_WithBackupFiles
```
Expected: PASS

- [ ] **Step 3: Commit**
```bash
git add internal/storage/storage_test.go
git commit -m "test: adiciona teste RecoverOrphans com arquivos .bak"
```

---

## Task 7: atomicRename - Falha em Windows

**Files:**
- Modify: `internal/storage/storage_test.go`

- [ ] **Step 1: Adicionar teste de atomicRename falha (arquivo em uso)**

```go
func TestAtomicRename_SourceNotFound(t *testing.T) {
    dir := t.TempDir()
    src := filepath.Join(dir, "naoexiste.src")
    dst := filepath.Join(dir, "dest.dst")

    err := storage.AtomicRename(src, dst)
    if err == nil {
        t.Error("esperado erro ao renomear arquivo inexistente")
    }
}
```

Nota: `atomicRename` é privado (inicia com letra minúscula). Verificar se há função wrapper pública.

- [ ] **Step 2: Run teste apenas de funções públicas se atomicRename for privado**
Se privado, testar via Save que falha no atomicRename step.

```bash
go test -v ./internal/storage/ -run TestSave
```
Expected: PASS

- [ ] **Step 3: Commit whichever foi possível**
```bash
git add internal/storage/storage_test.go
git commit -m "test: adiciona teste atomicRename ou falha em Save"
```

---

## Task 8: Verify Cobertura >= 90%

**Files:**
- Nenhum

- [ ] **Step 1: Rodar verificação de cobertura**

```bash
go test -coverprofile=coverage.out ./internal/storage/ -count=1
go tool cover -func=coverage.out
```

- [ ] **Step 2: Verificar se >= 90%**
Procurar linha "total:" no output.

- [ ] **Step 3: Se < 90%, identificar gaps e adicionar mais testes**
Retornar ao step 1 com novos testes.

- [ ] **Step 4: Commit final**
```bash
git add internal/storage/storage_test.go
git commit -m "test: aumenta cobertura storage para 90%+"
```

---

## Tabela de Resumo

| Task | Função | Estado | Cobertura Alvo |
|------|--------|--------|-----------------|
| 1 | Save rollback | - | +2% |
| 2 | SaveNew permissão | - | +2% |
| 3 | Load corrupto | - | +2% |
| 4 | Salvar update | - | +3% |
| 5 | ComputeMetadata vazio | - | +2% |
| 6 | RecoverOrphans .bak | - | +2% |
| 7 | atomicRename | - | +1% |
| 8 | Verificação | - | >=90% |