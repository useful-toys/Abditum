# Funcionalidade: Excluir e restaurar segredo

**Como** usuário do Abditum,
**quero** excluir segredos de forma reversível e poder restaurá-los antes do próximo salvamento,
**para** ter uma rede de segurança contra exclusões acidentais.

A exclusão de segredos é reversível: o segredo é movido para a Lixeira (pasta virtual) e pode ser restaurado até o próximo salvamento do cofre. Ao salvar, segredos na Lixeira são permanentemente excluídos.

---

**Contexto:**

  **Dado** que o cofre está aberto

---

## Exclusão reversível

### **Regra:** O segredo excluído é movido para a Lixeira e pode ser restaurado

#### **Cenário:** Excluir segredo reversivelmente

  **Dado** que o segredo "Gmail" está disponível na pasta "Sites"
  **Quando** o usuário exclui o segredo "Gmail"
  **E** confirma a exclusão
  **Então** o segredo "Gmail" é removido da pasta "Sites"
  **E** o segredo aparece na pasta virtual Lixeira
  **E** a identidade e o conteúdo do segredo são preservados
  **E** o cofre entra no estado "Cofre Modificado"

#### **Cenário:** Segredo na Lixeira não pode ser editado

  **Dado** que o segredo "Gmail" está na Lixeira
  **Então** a ação de editar o segredo "Gmail" não está disponível

---

## Restauração

### **Regra:** O segredo restaurado retorna à pasta de origem com seu estado anterior

#### **Cenário:** Restaurar segredo cuja pasta de origem ainda existe

  **Dado** que o segredo "Gmail" foi excluído reversivelmente da pasta "Sites"
  **E** a pasta "Sites" ainda existe na hierarquia
  **Quando** o usuário restaura o segredo "Gmail" da Lixeira
  **Então** o segredo "Gmail" é reinserido na pasta "Sites" ao final da lista de segredos
  **E** o segredo retorna ao estado que possuía antes da exclusão
  **E** o cofre entra no estado "Cofre Modificado"

#### **Cenário:** Restaurar segredo cuja pasta de origem foi excluída

  **Dado** que o segredo "Gmail" foi excluído reversivelmente da pasta "Sites"
  **E** a pasta "Sites" foi excluída após o soft delete do segredo
  **Quando** o usuário restaura o segredo "Gmail" da Lixeira
  **Então** o segredo "Gmail" é reinserido na raiz do cofre ao final da lista de segredos
  **E** a aplicação exibe mensagem informando que a pasta original não existe mais

---

## Lixeira

### **Regra:** A Lixeira é visível apenas quando há segredos excluídos reversivelmente

#### **Cenário:** Lixeira visível

  **Dado** que existe pelo menos um segredo excluído reversivelmente
  **Então** a pasta virtual Lixeira é exibida no final da raiz da hierarquia

#### **Cenário:** Lixeira oculta

  **Dado** que não existem segredos excluídos reversivelmente
  **Então** a pasta virtual Lixeira não é exibida na hierarquia
