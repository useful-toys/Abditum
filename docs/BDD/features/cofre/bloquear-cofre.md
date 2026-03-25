# Funcionalidade: Bloquear cofre

**Como** usuário do Abditum,
**quero** bloquear o acesso ao cofre manual ou automaticamente,
**para** proteger meus dados sensíveis quando me ausento ou em situações de risco.

O bloqueio fecha logicamente o cofre, limpa buffers controlados e a área de transferência, e retorna ao fluxo de abertura — exigindo nova autenticação. Alterações não salvas são descartadas silenciosamente, por decisão de projeto: bloqueios ocorrem em contextos onde confirmações comprometeriam o propósito da proteção.

---

**Contexto:**

  **Dado** que o cofre está aberto

---

## **Regra:** O bloqueio descarta alterações não salvas silenciosamente

### **Cenário:** Bloquear cofre manualmente em estado "Cofre Salvo"

  **Quando** o usuário bloqueia o cofre manualmente
  **Então** a aplicação fecha logicamente o cofre
  **E** limpa os buffers controlados sempre que possível
  **E** limpa a área de transferência
  **E** retorna ao fluxo de abertura do cofre assumindo o mesmo caminho do cofre previamente aberto

### **Cenário:** Bloquear cofre manualmente em estado "Cofre Modificado"

  **Dado** que o cofre está no estado "Cofre Modificado"
  **Quando** o usuário bloqueia o cofre manualmente
  **Então** as alterações não salvas são descartadas silenciosamente, sem confirmação
  **E** a aplicação fecha logicamente o cofre
  **E** limpa os buffers controlados e a área de transferência
  **E** retorna ao fluxo de abertura do cofre assumindo o mesmo caminho

---

## **Regra:** O bloqueio automático ocorre após inatividade configurada

### **Cenário:** Bloqueio automático por inatividade

  **Dado** que o tempo de bloqueio automático está configurado para 2 minutos
  **Quando** o usuário permanece inativo por 2 minutos
  **Então** a aplicação executa o mesmo processo do bloqueio manual
  **E** retorna ao fluxo de abertura do cofre

---

## **Regra:** O desbloqueio exige nova autenticação completa

### **Cenário:** Desbloquear cofre após bloqueio

  **Dado** que o cofre foi bloqueado
  **Então** a aplicação apresenta o fluxo de abertura do cofre com o caminho previamente utilizado
  **Quando** o usuário informa a senha mestra correta
  **Então** o cofre é reaberto a partir do arquivo salvo
  **E** o cofre entra no estado "Cofre Salvo"
