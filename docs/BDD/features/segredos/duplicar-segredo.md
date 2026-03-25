# Funcionalidade: Duplicar segredo

**Como** usuário do Abditum,
**quero** duplicar um segredo existente,
**para** criar uma cópia rápida com os mesmos dados e a partir dela fazer alterações.

---

**Contexto:**

  **Dado** que o cofre está aberto
  **E** existe um segredo "Gmail" na pasta "Sites"

---

## **Regra:** A duplicação cria uma cópia com nova identidade e nome sufixado

### **Cenário:** Duplicar segredo existente

  **Quando** o usuário duplica o segredo "Gmail"
  **Então** um novo segredo é criado com nova identidade
  **E** o nome do segredo duplicado é "Gmail (1)"
  **E** os campos, observação, nome do modelo e marcação de favorito são copiados do original
  **E** o segredo duplicado é inserido logo abaixo de "Gmail" na mesma pasta
  **E** o segredo duplicado assume o estado "Segredo novo"
  **E** o cofre entra no estado "Cofre Modificado"

### **Cenário:** Duplicar segredo quando já existe cópia com sufixo

  **Dado** que já existe o segredo "Gmail (1)" na pasta "Sites"
  **Quando** o usuário duplica o segredo "Gmail"
  **Então** o segredo duplicado recebe o nome "Gmail (2)"
