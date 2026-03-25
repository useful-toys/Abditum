# Funcionalidade: Configurar cofre

**Como** usuário do Abditum,
**quero** ajustar as configurações do cofre,
**para** personalizar os tempos de bloqueio, reocultação de campos sensíveis e limpeza da área de transferência.

As configurações são armazenadas dentro do próprio arquivo do cofre, garantindo portabilidade total. As alterações passam a valer para a sessão corrente imediatamente.

---

**Contexto:**

  **Dado** que o cofre está aberto

---

## **Regra:** As configurações possuem valores padrão sugeridos

### **Esquema do Cenário:** Alterar configuração do cofre

  **Quando** o usuário altera a configuração "<configuração>" para "<valor>"
  **Então** o novo valor passa a valer para a sessão corrente
  **E** o cofre entra no estado "Cofre Modificado"

  **Exemplos:**

  | configuração                            | valor |
  |-----------------------------------------|-------|
  | tempo de bloqueio automático por inatividade | 5 minutos  |
  | tempo de reocultação de campos sensíveis     | 10 segundos |
  | tempo de limpeza da área de transferência    | 15 segundos |
