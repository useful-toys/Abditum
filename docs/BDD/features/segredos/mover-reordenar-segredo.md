# Funcionalidade: Mover e reordenar segredo

**Como** usuário do Abditum,
**quero** mover segredos entre pastas e reordená-los dentro de uma mesma pasta,
**para** manter meu cofre organizado conforme minhas preferências.

---

**Contexto:**

  **Dado** que o cofre está aberto

---

## Mover segredo

### **Regra:** O segredo movido preserva identidade, conteúdo e marcação de favorito

#### **Cenário:** Mover segredo para outra pasta

  **Dado** que o segredo "Gmail" está na pasta "Sites"
  **Quando** o usuário move o segredo "Gmail" para a pasta "Serviços"
  **Então** o segredo "Gmail" é removido da pasta "Sites"
  **E** o segredo "Gmail" é adicionado ao final da lista de segredos da pasta "Serviços"
  **E** a identidade, o conteúdo e a marcação de favorito são preservados
  **E** o cofre entra no estado "Cofre Modificado"

#### **Cenário:** Mover segredo para a raiz do cofre

  **Dado** que o segredo "Gmail" está na pasta "Sites"
  **Quando** o usuário move o segredo "Gmail" para a raiz do cofre
  **Então** o segredo "Gmail" é adicionado ao final da lista de segredos da raiz
  **E** o cofre entra no estado "Cofre Modificado"

---

## Reordenar segredo

### **Regra:** A reordenação altera apenas a posição entre segredos irmãos

#### **Cenário:** Reordenar segredo dentro da mesma pasta

  **Dado** que a pasta "Sites" contém os segredos "Gmail", "Outlook" e "Yahoo" nesta ordem
  **Quando** o usuário move "Yahoo" para a primeira posição
  **Então** a ordem dos segredos na pasta "Sites" passa a ser "Yahoo", "Gmail", "Outlook"
  **E** a identidade e o conteúdo dos segredos são preservados
  **E** o cofre entra no estado "Cofre Modificado"
