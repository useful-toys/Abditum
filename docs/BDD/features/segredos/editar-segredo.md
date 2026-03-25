# Funcionalidade: Editar segredo

**Como** usuário do Abditum,
**quero** editar os dados e a estrutura dos meus segredos,
**para** manter minhas credenciais atualizadas e organizadas.

A edição de segredos possui dois modos: edição padrão (alterar valores) e edição avançada (alterar estrutura de campos). O usuário pode alternar entre os modos durante a edição.

---

**Contexto:**

  **Dado** que o cofre está aberto
  **E** existe um segredo disponível no cofre

---

## Edição padrão

### **Regra:** A edição padrão permite alterar nome, observação e valores dos campos

#### **Cenário:** Editar nome e valor de campo de um segredo

  **Quando** o usuário inicia a edição padrão do segredo
  **E** altera o nome do segredo para "Gmail Pessoal"
  **E** altera o valor do campo "Username" para "usuario@gmail.com"
  **E** confirma a edição
  **Então** o segredo é atualizado com o novo nome e valor
  **E** a identidade do segredo é preservada
  **E** o cofre entra no estado "Cofre Modificado"

#### **Cenário:** Cancelar edição padrão

  **Quando** o usuário inicia a edição padrão do segredo
  **E** altera o nome do segredo
  **E** cancela a edição
  **Então** todas as alterações locais são revertidas
  **E** o segredo permanece com os dados anteriores

---

## Edição avançada

### **Regra:** A edição avançada permite alterar a estrutura de campos

#### **Cenário:** Adicionar novo campo ao segredo

  **Quando** o usuário inicia a edição avançada do segredo
  **E** adiciona um novo campo com nome "Código de Recuperação" e tipo "texto sensível"
  **E** confirma a edição
  **Então** o segredo passa a conter o campo "Código de Recuperação"
  **E** o cofre entra no estado "Cofre Modificado"

#### **Cenário:** Renomear campo existente

  **Quando** o usuário inicia a edição avançada do segredo
  **E** altera o nome do campo "URL" para "Endereço Web"
  **E** confirma a edição
  **Então** o campo passa a se chamar "Endereço Web", preservando seu tipo e valor

#### **Cenário:** Excluir campo do segredo

  **Quando** o usuário inicia a edição avançada do segredo
  **E** exclui o campo "URL"
  **E** confirma a edição
  **Então** o campo "URL" é removido do segredo

#### **Cenário:** Reordenar campos do segredo

  **Quando** o usuário inicia a edição avançada do segredo
  **E** move o campo "Password" para a primeira posição
  **E** confirma a edição
  **Então** o campo "Password" passa a aparecer como primeiro campo do segredo

---

### **Regra:** Não é permitido alterar o tipo de um campo existente

#### **Cenário:** Tentativa de alterar tipo de campo

  **Quando** o usuário inicia a edição avançada do segredo
  **Então** a opção de alterar o tipo de um campo existente não está disponível
  **E** para mudar o tipo, o usuário deve excluir o campo e criar um novo com o tipo desejado

---

### **Regra:** O estado do segredo é preservado de forma consistente

#### **Esquema do Cenário:** Estado do segredo após edição

  **Dado** que o segredo está no estado "<estado_anterior>"
  **Quando** o usuário edita e confirma alterações no segredo
  **Então** o segredo assume o estado "<estado_posterior>"

  **Exemplos:**

  | estado_anterior      | estado_posterior     |
  |----------------------|----------------------|
  | Segredo disponível   | Segredo modificado   |
  | Segredo novo         | Segredo novo         |
  | Segredo modificado   | Segredo modificado   |
