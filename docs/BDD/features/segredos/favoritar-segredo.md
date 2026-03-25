# Funcionalidade: Favoritar e desfavoritar segredo

**Como** usuário do Abditum,
**quero** marcar segredos como favoritos,
**para** ter acesso rápido aos segredos que uso com mais frequência.

Segredos favoritos ganham destaque visual na hierarquia e aparecem na pasta virtual Favoritos, no topo da raiz. A marcação de favorito não altera identidade, conteúdo ou localização do segredo.

---

**Contexto:**

  **Dado** que o cofre está aberto

---

## **Regra:** Favoritar altera apenas o atributo de favorito

### **Cenário:** Favoritar segredo

  **Dado** que o segredo "Gmail" não está favoritado
  **Quando** o usuário favorita o segredo "Gmail"
  **Então** o segredo "Gmail" aparece na pasta virtual Favoritos
  **E** o segredo mantém sua localização original na hierarquia
  **E** o cofre entra no estado "Cofre Modificado"

### **Cenário:** Desfavoritar segredo

  **Dado** que o segredo "Gmail" está favoritado
  **Quando** o usuário desfavorita o segredo "Gmail"
  **Então** o segredo "Gmail" é removido da pasta virtual Favoritos
  **E** o segredo mantém sua localização original na hierarquia
  **E** o cofre entra no estado "Cofre Modificado"

---

## **Regra:** A pasta virtual Favoritos é exibida apenas quando há favoritos

### **Cenário:** Pasta virtual Favoritos visível

  **Dado** que existe pelo menos um segredo favoritado
  **Então** a pasta virtual Favoritos é exibida no topo da raiz da hierarquia

### **Cenário:** Pasta virtual Favoritos oculta

  **Dado** que nenhum segredo está favoritado
  **Então** a pasta virtual Favoritos não é exibida na hierarquia
