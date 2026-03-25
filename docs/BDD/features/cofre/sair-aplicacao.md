# Funcionalidade: Sair da aplicação

**Como** usuário do Abditum,
**quero** encerrar a aplicação de forma segura,
**para** proteger meus dados e evitar perda acidental de alterações não salvas.

---

## **Regra:** Sair com alterações não salvas oferece opções de tratamento

### **Cenário:** Sair sem cofre ativo

  **Dado** que a aplicação está no estado inicial, sem cofre ativo
  **Quando** o usuário inicia a ação de sair
  **E** confirma o encerramento
  **Então** a aplicação é encerrada

### **Cenário:** Sair com cofre em estado "Cofre Salvo"

  **Dado** que o cofre está aberto e no estado "Cofre Salvo"
  **Quando** o usuário inicia a ação de sair
  **E** confirma o encerramento
  **Então** a aplicação é encerrada

### **Cenário:** Sair com cofre em estado "Cofre Modificado", salvando

  **Dado** que o cofre está aberto e no estado "Cofre Modificado"
  **Quando** o usuário inicia a ação de sair
  **E** escolhe "Salvar"
  **Então** a aplicação salva o cofre seguindo o fluxo padrão de salvamento
  **E** a aplicação é encerrada após o salvamento bem-sucedido

### **Cenário:** Sair com cofre em estado "Cofre Modificado", descartando

  **Dado** que o cofre está aberto e no estado "Cofre Modificado"
  **Quando** o usuário inicia a ação de sair
  **E** escolhe "Sair sem Salvar"
  **Então** as alterações não salvas são descartadas
  **E** a aplicação é encerrada

### **Cenário:** Sair com cofre em estado "Cofre Modificado", cancelando

  **Dado** que o cofre está aberto e no estado "Cofre Modificado"
  **Quando** o usuário inicia a ação de sair
  **E** escolhe "Voltar"
  **Então** a aplicação retorna ao estado anterior
  **E** o cofre permanece no estado "Cofre Modificado"
