# Funcionalidade: Salvar cofre

**Como** usuário do Abditum,
**quero** salvar o cofre para persistir minhas alterações no arquivo,
**para** não perder as modificações feitas nos segredos, pastas ou configurações.

O salvamento persiste o estado atual do cofre no arquivo `.abditum`. O processo utiliza gravação atômica (via arquivo `.tmp`) para proteger contra corrupção, e mantém um backup rotacionado do arquivo anterior.

---

**Contexto:**

  **Dado** que o cofre está aberto e no estado "Cofre Modificado"

---

## **Regra:** O salvamento usa gravação atômica via arquivo temporário

### **Cenário:** Salvar cofre com sucesso

  **Quando** o usuário salva o cofre
  **Então** a aplicação grava os dados em um arquivo ".abditum.tmp"
  **E** copia o arquivo atual do cofre para ".abditum.bak"
  **E** renomeia o arquivo ".abditum.tmp" para o nome final do cofre
  **E** o cofre entra no estado "Cofre Salvo"

### **Cenário:** Falha na gravação do arquivo temporário

  **Quando** o usuário salva o cofre
  **E** a gravação do arquivo ".abditum.tmp" falha
  **Então** o arquivo ".abditum.tmp" é imediatamente apagado
  **E** o arquivo original do cofre permanece inalterado
  **E** o cofre permanece no estado "Cofre Modificado"

---

## **Regra:** O backup anterior é rotacionado a cada salvamento

### **Cenário:** Salvar cofre quando já existe backup anterior

  **Dado** que já existe um backup em ".abditum.bak"
  **Quando** o usuário salva o cofre
  **Então** o backup anterior é renomeado temporariamente para ".abditum.bak2"
  **E** o arquivo atual é copiado para ".abditum.bak"
  **E** o arquivo ".abditum.tmp" é renomeado para o nome final
  **E** o arquivo ".abditum.bak2" é removido

### **Cenário:** Falha após geração de backup

  **Dado** que já existe um backup em ".abditum.bak"
  **Quando** o usuário salva o cofre
  **E** a gravação ou substituição falha após a geração do backup
  **Então** o arquivo ".abditum.bak2" é restaurado para ".abditum.bak"
  **E** a aplicação exibe mensagem de erro informando a falha e a existência do backup para intervenção manual

---

## **Regra:** O salvamento regenera o nonce criptográfico

### **Cenário:** Cada salvamento usa nonce diferente

  **Quando** o usuário salva o cofre
  **Então** a aplicação gera um novo nonce para o AES-256-GCM
  **E** atualiza a versão do formato no cabeçalho quando necessário

---

## **Regra:** Segredos na Lixeira são permanentemente excluídos ao salvar

### **Cenário:** Salvar cofre com segredos na Lixeira

  **Dado** que existem segredos excluídos reversivelmente na Lixeira
  **Quando** o usuário salva o cofre
  **Então** os segredos da Lixeira são permanentemente excluídos, sem possibilidade de recuperação
  **E** a Lixeira fica vazia
  **E** o cofre entra no estado "Cofre Salvo"
