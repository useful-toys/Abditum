# Funcionalidade: Descartar alterações e recarregar cofre

**Como** usuário do Abditum,
**quero** descartar todas as alterações não salvas e recarregar o cofre a partir do arquivo,
**para** reverter modificações indesejadas sem precisar fechar e reabrir a aplicação.

---

**Contexto:**

  **Dado** que o cofre está aberto e no estado "Cofre Modificado"

---

## **Regra:** A ação exige confirmação e só está disponível em "Cofre Modificado"

### **Cenário:** Descartar alterações com confirmação

  **Quando** o usuário inicia a ação de descartar alterações
  **E** confirma o descarte
  **Então** a aplicação reabre o arquivo atual, repete validação, descriptografia e eventual migração em memória
  **E** o cofre retorna ao estado "Cofre Salvo"

### **Cenário:** Cancelar descarte de alterações

  **Quando** o usuário inicia a ação de descartar alterações
  **E** cancela o descarte
  **Então** o cofre permanece no estado "Cofre Modificado" com todas as alterações preservadas

### **Cenário:** Ação indisponível em "Cofre Salvo"

  **Dado** que o cofre está aberto e no estado "Cofre Salvo"
  **Então** a ação de descartar alterações não está disponível
