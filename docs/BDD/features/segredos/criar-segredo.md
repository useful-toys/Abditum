# Funcionalidade: Criar segredo

**Como** usuário do Abditum,
**quero** criar novos segredos no cofre,
**para** armazenar credenciais e informações confidenciais de forma organizada.

O usuário pode criar um segredo a partir de um modelo existente (com campos pré-definidos) ou começar com um segredo vazio. O segredo criado a partir de um modelo recebe uma cópia dos campos como snapshot, sem manter vínculo com o modelo de origem.

---

**Contexto:**

  **Dado** que o cofre está aberto

---

## **Regra:** Um segredo pode ser criado a partir de um modelo ou vazio

### **Cenário:** Criar segredo usando modelo "Login"

  **Quando** o usuário inicia a criação de um segredo na pasta "Sites"
  **E** escolhe o modelo "Login"
  **Então** um novo segredo é criado com os campos "URL", "Username" e "Password" copiados do modelo
  **E** o nome do modelo "Login" é registrado no segredo como histórico
  **E** a aplicação abre o segredo no modo de edição padrão
  **E** o segredo assume o estado "Segredo novo"
  **E** o cofre entra no estado "Cofre Modificado"

### **Cenário:** Criar segredo vazio

  **Quando** o usuário inicia a criação de um segredo na raiz do cofre
  **E** escolhe começar com um segredo vazio
  **Então** um novo segredo é criado sem campos adicionais, apenas com nome e observação
  **E** a aplicação abre o segredo no modo de edição avançada
  **E** o segredo assume o estado "Segredo novo"
  **E** o cofre entra no estado "Cofre Modificado"

---

## **Regra:** O segredo é inserido conforme o contexto de foco

### **Cenário:** Criar segredo com foco em uma pasta

  **Dado** que o foco está na pasta "Financeiro"
  **Quando** o usuário cria um novo segredo
  **Então** o segredo é inserido ao final da lista de segredos da pasta "Financeiro"

### **Cenário:** Criar segredo com foco em um segredo dentro de uma pasta

  **Dado** que o foco está no segredo "Gmail" dentro da pasta "Sites"
  **Quando** o usuário cria um novo segredo
  **Então** o segredo é inserido logo abaixo de "Gmail" na pasta "Sites"

### **Cenário:** Criar segredo com foco na raiz do cofre

  **Dado** que o foco está na raiz do cofre
  **Quando** o usuário cria um novo segredo
  **Então** o segredo é inserido ao final da lista de segredos da raiz

---

## **Regra:** O modelo é apenas um snapshot — alterações no modelo não afetam segredos já criados

### **Cenário:** Alterar modelo após criar segredo a partir dele

  **Dado** que o segredo "Meu Login" foi criado a partir do modelo "Login"
  **Quando** o usuário altera o modelo "Login" adicionando o campo "Notas"
  **E** salva o cofre
  **Então** o segredo "Meu Login" permanece com os campos originais "URL", "Username" e "Password"
  **E** o campo "Notas" não aparece no segredo "Meu Login"
