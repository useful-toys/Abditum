# Funcionalidade: Buscar segredos

**Como** usuário do Abditum,
**quero** buscar segredos por nome, nome de campo, valor de campos de texto ou observação,
**para** localizar rapidamente um segredo sem navegar manualmente pela hierarquia.

A busca é realizada por varredura sequencial em memória, sem índices persistidos. Campos do tipo texto sensível nunca participam da busca, independentemente do estado visual.

---

**Contexto:**

  **Dado** que o cofre está aberto

---

## **Regra:** A busca filtra a hierarquia preservando o contexto de pastas

### **Cenário:** Buscar segredo por nome

  **Dado** que o cofre contém o segredo "Gmail" na pasta "Sites"
  **Quando** o usuário busca por "Gmail"
  **Então** a hierarquia é filtrada exibindo apenas a pasta "Sites" com o segredo "Gmail"
  **E** o nome "Gmail" recebe destaque visual na árvore

### **Cenário:** Buscar segredo por valor de campo de texto

  **Dado** que o segredo "Gmail" possui o campo "Username" com valor "usuario@gmail.com"
  **Quando** o usuário busca por "usuario@gmail"
  **Então** o segredo "Gmail" aparece nos resultados

### **Cenário:** Buscar segredo por observação

  **Dado** que o segredo "API Interna" possui a observação "servidor de homologação"
  **Quando** o usuário busca por "homologação"
  **Então** o segredo "API Interna" aparece nos resultados

### **Cenário:** Buscar segredo por nome de campo

  **Dado** que o segredo "API Interna" possui um campo chamado "Chave de API"
  **Quando** o usuário busca por "Chave de API"
  **Então** o segredo "API Interna" aparece nos resultados

---

## **Regra:** Campos de texto sensível nunca participam da busca

### **Cenário:** Busca não encontra valor em campo sensível

  **Dado** que o segredo "Gmail" possui o campo "Password" do tipo "texto sensível" com valor "MinhaSenh@123"
  **Quando** o usuário busca por "MinhaSenh@123"
  **Então** nenhum resultado é encontrado

---

## **Regra:** Durante a busca, apenas ações de navegação e visualização estão disponíveis

### **Cenário:** Restrição de ações durante busca ativa

  **Dado** que uma busca está ativa
  **Então** as ações disponíveis são: sair da aplicação, navegar pelo cofre e visualizar segredo
  **E** as demais ações de edição, exclusão e movimentação não estão disponíveis

### **Cenário:** Encerrar busca selecionando um resultado

  **Dado** que uma busca está ativa
  **Quando** o usuário seleciona o segredo "Gmail" nos resultados
  **Então** a busca é encerrada
  **E** o cofre retorna ao estado anterior ao início da pesquisa
  **E** o foco permanece sobre o segredo "Gmail"

### **Cenário:** Cancelar busca

  **Dado** que uma busca está ativa
  **Quando** o usuário cancela a pesquisa
  **Então** a busca é encerrada
  **E** o cofre retorna ao estado anterior ao início da pesquisa
