# Funcionalidade: Gerenciar modelos de segredo

**Como** usuário do Abditum,
**quero** criar, editar e remover modelos de segredo,
**para** padronizar e agilizar a criação de novos segredos com estruturas recorrentes.

Modelos de segredo definem uma estrutura de campos (nome e tipo) que serve como template na criação de segredos. Alterações em um modelo afetam apenas criações futuras — segredos já criados não são modificados.

---

**Contexto:**

  **Dado** que o cofre está aberto

---

## Criar modelo de segredo

### **Regra:** O modelo recebe identidade própria e fica disponível para uso futuro

#### **Cenário:** Criar modelo de segredo personalizado

  **Quando** o usuário cria um modelo de segredo com o nome "Conta Bancária"
  **E** adiciona o campo "Agência" do tipo "texto"
  **E** adiciona o campo "Conta" do tipo "texto"
  **E** adiciona o campo "Senha" do tipo "texto sensível"
  **E** confirma a criação
  **Então** o modelo "Conta Bancária" fica disponível para criação de novos segredos
  **E** o cofre entra no estado "Cofre Modificado"

---

## Criar modelo a partir de segredo existente

### **Regra:** A estrutura de campos do segredo é copiada como base do modelo

#### **Cenário:** Criar modelo a partir de segredo

  **Dado** que existe o segredo "Gmail" com os campos "URL" (texto), "Username" (texto) e "Password" (texto sensível)
  **Quando** o usuário cria um modelo a partir do segredo "Gmail"
  **E** informa o nome "E-mail" para o novo modelo
  **E** confirma a criação
  **Então** o modelo "E-mail" é criado com os campos "URL", "Username" e "Password" com os respectivos tipos
  **E** nenhum vínculo retroativo é criado entre o modelo e o segredo de origem
  **E** o cofre entra no estado "Cofre Modificado"

---

## Editar modelo de segredo

### **Regra:** Alterações no modelo afetam apenas criações futuras

#### **Cenário:** Adicionar campo ao modelo

  **Dado** que existe o modelo "Login" com os campos "URL", "Username" e "Password"
  **Quando** o usuário edita o modelo "Login"
  **E** adiciona o campo "Notas" do tipo "texto"
  **E** confirma a edição
  **Então** o modelo "Login" passa a ter os campos "URL", "Username", "Password" e "Notas"
  **E** segredos já criados a partir do modelo "Login" permanecem inalterados
  **E** o cofre entra no estado "Cofre Modificado"

#### **Cenário:** Alterar tipo de campo no modelo

  **Dado** que existe o modelo "Login" com o campo "URL" do tipo "texto"
  **Quando** o usuário edita o modelo "Login"
  **E** altera o tipo do campo "URL" para "texto sensível"
  **E** confirma a edição
  **Então** o campo "URL" do modelo passa a ser do tipo "texto sensível"
  **E** segredos já criados permanecem inalterados

#### **Cenário:** Excluir campo do modelo

  **Quando** o usuário edita o modelo "Login"
  **E** exclui o campo "URL"
  **E** confirma a edição
  **Então** o modelo "Login" não contém mais o campo "URL"
  **E** segredos já criados a partir do modelo permanecem com o campo "URL" intacto

#### **Cenário:** Reordenar campos do modelo

  **Quando** o usuário edita o modelo "Login"
  **E** move o campo "Password" para a primeira posição
  **E** confirma a edição
  **Então** o campo "Password" passa a ser o primeiro campo do modelo

---

## Remover modelo de segredo

### **Regra:** A remoção do modelo não afeta segredos já criados

#### **Cenário:** Remover modelo de segredo

  **Dado** que existe o modelo "API Key"
  **E** existem segredos criados a partir do modelo "API Key"
  **Quando** o usuário remove o modelo "API Key"
  **E** confirma a remoção
  **Então** o modelo "API Key" não está mais disponível para criação de novos segredos
  **E** os segredos criados anteriormente a partir dele permanecem inalterados
  **E** o cofre entra no estado "Cofre Modificado"

---

## Modelos pré-definidos

### **Regra:** O cofre é criado com modelos pré-definidos editáveis e removíveis

#### **Esquema do Cenário:** Modelos pré-definidos no cofre

  **Dado** que um novo cofre foi criado
  **Então** o cofre contém o modelo "<modelo>" com os campos "<campos>"

  **Exemplos:**

  | modelo            | campos                                                     |
  |-------------------|------------------------------------------------------------|
  | Login             | URL, Username, Password                                    |
  | Cartão de Crédito | Número do Cartão, Nome no Cartão, Data de Validade, CVV    |
  | API Key           | Nome da API, Chave de API                                  |
