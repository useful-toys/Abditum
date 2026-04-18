# Omissão de declaração de teclas ENTER e ESC em ConfirmModal — Plano de Implementação

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Modificar KeyHandler para tratar implicitamente a primeira opção como ativada por Enter e a última opção como ativada por Esc (além das teclas declaradas em Keys), permitindo que desenvolvedores omitam a declaração de Keys nesses casos comuns e aprimorando o comportamento existente.

**Architecture:** 
1. Modificar key_handler.go na função Handle para, ao iterar sobre as opções, adicionar implicitamente Enter como tecla válida para a primeira opção e Esc como tecla válida para a última opção
2. Isso afetará todos os modais que usam KeyHandler (que é praticamente todos), mas o comportamento é particularmente útil para ConfirmModal devido ao seu padrão de uso
3. Manter total compatibilidade com código existente - quando Keys tem elementos declarados, comportamento é aprimorado (mais teclas funcionais) mas nunca removido
4. Atualizar testes para refletir o novo comportamento
5. Nenhuma mudança necessária em confirm_modal.go ou frame.go

**Tech Stack:** Go, go build, go test

---

### Task 1: Modificar key_handler.go para adicionar comportamento implícito de Enter e Esc

**Files:**
- Modify: `internal/tui/modal/key_handler.go`

- [ ] **Step 1: Entender o fluxo atual de despacho de teclas**

  Em key_handler.go, linhas 36-42, temos:
  ```go
  for _, opt := range h.Options {
      for _, k := range opt.Keys {
          if k.Matches(msg) {
              return opt.Action(), true
          }
      }
  }
  ```

  Isso itera sobre cada opção e cada tecla declarada em opt.Keys.

  Vamos modificar para:
  - Para a primeira opção (índice 0), além de verificar as teclas em opt.Keys, também verificar se a tecla recebida é Enter
  - Para a última opção (índice len-1), além de verificar as teclas em opt.Keys, também verificar se a tecla recebida é Esc
  - Para opções do meio, apenas verificar as teclas em opt.Keys (comportamento atual)
  - Se houver apenas uma opção (que é tanto primeira quanto última), verificar tanto Enter quanto Esc

  Importante: precisamos saber o índice da opção durante a iteração.

- [ ] **Step 2: Implementar a lógica de detecção implícita de teclas**

  Substituir o laço atual por algo como:

  ```go
  // 1. Despachar ações registradas.
  for i, opt := range h.Options {
      isFirst := i == 0
      isLast := i == len(h.Options)-1
      
      // Verificar teclas declaradas
      for _, k := range opt.Keys {
          if k.Matches(msg) {
              return opt.Action(), true
          }
      }
      
      // Verificar teclas implícitas
      if isFirst && design.Keys.Enter.Matches(msg) {
          return opt.Action(), true
      }
      if isLast && design.Keys.Esc.Matches(msg) {
          return opt.Action(), true
      }
      // Nota: se houver apenas uma opção, ambas as condições acima podem ser verdadeiras
      // e o primeiro match retornará (o que é fine)
  }
  ```

  Essa abordagem mantém a ordem de verificação: primeiro teclas declaradas, depois teclas implícitas.
  Alternativamente, poderíamos verificar as teclas implícitas primeiro, mas a ordem não é crítica desde que funcionem.

- [ ] **Step 3: Testar a mudança compilando**

  ```bash
  go build ./internal/tui/modal/...
  ```

### Task 2: Verificar compatibilidade com casos existentes

**Files:**
- Modify: `internal/tui/modal/confirm_modal_test.go`
- Modify: `internal/tui/modal/help_modal_test.go` (se existir)
- Modify: `internal/tui/modal/frame_test.go`
- Modify: `internal/tui/modal/key_handler_test.go`

- [ ] **Step 1: Executar testes existentes para garantir que nada quebrou**

  ```bash
  go test ./internal/tui/modal/...
  ```

  Esperamos que alguns testem falhem porque eles esperavam um certo comportamento de teclas que agora foi aprimorado (mais teclas funcionais).
  Especificamente, testes que criavam ModalOption sem Enter/Esperavam que apenas as teclas declaradas funcionassem, agora vão descobrir que Enter e/ou Esc também funcionam para primeira/última opção.

- [ ] **Step 2: Atualizar testes que estavam verificando o comportamento antigo**

  Precisamos identificar quais testes estão verificando especificamente que somente as teclas declaradas funcionavam.
  
  Exemplos de testes que provavelmente precisarão de atualização:
  - Testes que criam uma opção com Keys vazio ou nil e esperam que nenhuma tecla funcione (agora Enter ou Esc vão funcionar dependendo da posição)
  - Testes que verificam que uma tecla específica NÃO funciona em uma certa posição (precisarão ser revisados)

  Vamos olhar os testes existentes e atualizar conforme necessário.

### Task 3: Validar novos comportamentos com testes adicionais

**Files:**
- Modify: `internal/tui/modal/key_handler_test.go` (adicionar aos testes existentes)

- [ ] **Step 1: Adicionar testes para confirmar o novo comportamento**

  Testes a incluir:
  1. Quando primeira opção tem Keys vazio → deveria responder a Enter
  2. Quando última opção tem Keys vazio → deveria responder a Esc  
  3. Quando há apenas uma opção com Keys vazio → deveria responder a Enter E Esc
  4. Quando primeira opção tem Keys declarados (ex: AltC) → deveria responder a AltC E Enter
  5. Quando última opção tem Keys declarados (ex: F1) → deveria responder a F1 E Esc
  6. Quando opções do meio têm Keys declarados → devem responder apenas às teclas declaradas (não a Enter/Esc)
  7. Quando há apenas uma opção com Keys declarados (ex: F2) → deveria responder a F2, Enter E Esc

- [ ] **Step 2: Implementar os testes**

  Vamos usar a tabela de testes padrão do Go no key_handler_test.go.

### Task 4: Executar teste completo do projeto

**Files:**

- [ ] **Step 1: Build completo**

  ```bash
  go build ./...
  ```
  Esperado: sem erros

- [ ] **Step 2: Testes completos**

  ```bash
  go test ./...
  ```
  Esperado: todos PASS

### Task 5: Commit das mudanças

**Files:**

- [ ] **Step 1: Commit**

  ```bash
  git add internal/tui/modal/key_handler.go internal/tui/modal/confirm_modal_test.go internal/tui/modal/frame_test.go internal/tui/modal/key_handler_test.go
  git commit -m "feat: permitir omissão de declaração de teclas ENTER/ESC e aprimorar comportamento de teclas em KeyHandler"
  ```

  Nota: Mantivemos a lógica existente de adicionar ESC como alias para uma única opção por razões de compatibilidade, embora agora seja parcialmente redundante.

- [ ] **Step 3: Testar a mudança compilando**

  ```bash
  go build ./internal/tui/modal/...
  ```

### Task 2: Verificar compatibilidade com casos existentes

**Files:**
- Modify: `internal/tui/modal/confirm_modal_test.go`

- [ ] **Step 1: Executar testes existentes para garantir que nada quebrou**

  ```bash
  go test ./internal/tui/modal/...
  ```

  Esperamos que alguns testem falhem porque eles esperavam que Certas opções tivessem Keys vazio, mas agora elas serão preenchidas automaticamente.

- [ ] **Step 2: Atualizar testes que estavam verificando o comportamento antigo**

  Especificamente, olhar para testes que criavam ModalOption com Keys vazio ou nil e esperavam que permanecessem assim.

  Na verdade, vamos manter o teste TestConfirmModal_IntentTypes_Preserved? Não, aquele foi removido no spec anterior.

  Vamos olhar o que testes existem e atualizar conforme necessário.

### Task 3: Validar novos comportamentos com testes adicionais

**Files:**
- Create: `internal/tui/modal/confirm_modal_test.go` (adicionar aos testes existentes)
- Or Modify: se preferir manter tudo em um arquivo

- [ ] **Step 1: Adicionar testes para confirmar o novo comportamento**

  Testes a incluir:
  1. Quando primeira opção tem Keys vazio → deveria responder a Enter
  2. Quando última opção tem Keys vazio → deveria responder a Esc  
  3. Quando há apenas uma opção com Keys vazio → deveria responder a Enter (conforme especificado)
  4. Quando opções têm Keys declarados → comportamento deve permanecer exatamente o mesmo
  5. Quando primeira opção tem Keys com outras teclas → deveria responder a ESAS teclas E também Enter (como antes)
  6. Quando última opção tem Keys com outras teclas → deveria responder a ESAS teclas E também Esc (como antes)

  Porém, pensando melhor, os pontos 5 e 6 já são testados implicitamente pelo comportamento existente do key_handler - se deixarmos as teclas como estão, o key_handler já processa normalmente.

  Então focamos nos pontos 1-4.

- [ ] **Step 2: Implementar os testes**

  Vamos usar a tabela de testes padrão do Go.

### Task 4: Executar teste completo do projeto

**Files:**

- [ ] **Step 1: Build completo**

  ```bash
  go build ./...
  ```
  Esperado: sem erros

- [ ] **Step 2: Testes completos**

  ```bash
  go test ./...
  ```
  Esperado: todos PASS

### Task 5: Commit das mudanças

**Files:**

- [ ] **Step 1: Commit**

  ```bash
  git add internal/tui/modal/confirm_modal.go internal/tui/modal/confirm_modal_test.go
  git commit -m "feat: permitir omissão de declaração de teclas ENTER/ESC em ConfirmModal quando vazio/nil"
  ```
