# Funcionalidade: Copiar campo de segredo

**Como** usuário do Abditum,
**quero** copiar o valor de qualquer campo de um segredo para a área de transferência,
**para** usar rapidamente a informação em outro contexto sem precisar digitá-la.

---

**Contexto:**

  **Dado** que o cofre está aberto
  **E** o usuário está visualizando um segredo

---

## **Regra:** Qualquer campo pode ser copiado, inclusive campos sensíveis ocultos

### **Cenário:** Copiar campo de texto

  **Quando** o usuário copia o campo "Username" com valor "usuario@gmail.com"
  **Então** o valor "usuario@gmail.com" é copiado para a área de transferência
  **E** a aplicação exibe feedback visual de cópia

### **Cenário:** Copiar campo de texto sensível sem precisar exibi-lo

  **Dado** que o campo "Password" está oculto
  **Quando** o usuário copia o campo "Password"
  **Então** o valor do campo é copiado para a área de transferência
  **E** o campo permanece oculto na interface

---

## **Regra:** A área de transferência é limpa automaticamente após tempo configurado

### **Cenário:** Limpeza automática da área de transferência

  **Dado** que o tempo de limpeza da área de transferência está configurado para 30 segundos
  **Quando** o usuário copia um campo
  **Então** a aplicação inicia um temporizador de 30 segundos
  **E** após 30 segundos a área de transferência é limpa automaticamente

### **Cenário:** Limpeza da área de transferência ao bloquear cofre

  **Dado** que existe um valor copiado na área de transferência
  **Quando** o cofre é bloqueado
  **Então** a área de transferência é limpa imediatamente

### **Cenário:** Limpeza da área de transferência ao fechar cofre

  **Dado** que existe um valor copiado na área de transferência
  **Quando** a aplicação é encerrada
  **Então** a área de transferência é limpa imediatamente
