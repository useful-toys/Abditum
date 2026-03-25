# Funcionalidade: Alterar senha mestra

**Como** usuário do Abditum,
**quero** alterar a senha mestra do meu cofre,
**para** manter a segurança dos meus dados com uma nova credencial.

Alterar a senha mestra rederiva a chave criptográfica com novo salt e exige regravação completa do arquivo. Após a alteração, o cofre segue o fluxo de salvamento atômico.

---

**Contexto:**

  **Dado** que o cofre está aberto

---

## **Regra:** A nova senha mestra exige digitação dupla

### **Cenário:** Alterar senha mestra com confirmação válida

  **Quando** o usuário inicia a alteração da senha mestra
  **E** informa a nova senha mestra "Nov@Senh4!2026"
  **E** confirma a nova senha mestra "Nov@Senh4!2026"
  **Então** a aplicação rederiva a chave com um novo salt
  **E** segue o fluxo de salvamento do cofre
  **E** o cofre entra no estado "Cofre Salvo"

### **Cenário:** Falhar ao confirmar nova senha mestra

  **Quando** o usuário inicia a alteração da senha mestra
  **E** informa a nova senha mestra "Nov@Senh4!2026"
  **E** confirma a nova senha mestra "SenhaErrada"
  **Então** a aplicação exibe uma mensagem de erro indicando que as senhas não coincidem
  **E** a senha mestra permanece inalterada
