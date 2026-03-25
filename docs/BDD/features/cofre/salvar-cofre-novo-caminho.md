# Funcionalidade: Salvar cofre em novo caminho

**Como** usuário do Abditum,
**quero** salvar o cofre em um caminho diferente do atual,
**para** criar uma cópia do cofre em outro local ou renomear o arquivo.

Após a gravação bem-sucedida, o novo caminho passa a ser o caminho atual do cofre.

---

**Contexto:**

  **Dado** que o cofre está aberto

---

## **Regra:** Salvar em caminho vazio grava diretamente

### **Cenário:** Salvar cofre em novo caminho sem arquivo existente

  **Quando** o usuário inicia a ação de salvar em novo caminho
  **E** informa o caminho "/backup/cofre-copia.abditum" onde não existe arquivo
  **Então** o cofre é gravado diretamente em "/backup/cofre-copia.abditum"
  **E** o caminho atual do cofre passa a ser "/backup/cofre-copia.abditum"
  **E** o cofre entra no estado "Cofre Salvo"

---

## **Regra:** Sobrescrita em novo caminho exige confirmação

### **Cenário:** Salvar cofre em novo caminho com arquivo existente, confirmando

  **Dado** que já existe um arquivo em "/backup/outro.abditum"
  **Quando** o usuário inicia a ação de salvar em novo caminho "/backup/outro.abditum"
  **E** confirma a sobrescrita
  **Então** o arquivo existente é preservado como backup em "/backup/outro.abditum.bak"
  **E** o cofre é gravado em "/backup/outro.abditum"
  **E** o caminho atual do cofre passa a ser "/backup/outro.abditum"
  **E** o cofre entra no estado "Cofre Salvo"

### **Cenário:** Salvar cofre em novo caminho com arquivo existente, cancelando

  **Dado** que já existe um arquivo em "/backup/outro.abditum"
  **Quando** o usuário inicia a ação de salvar em novo caminho "/backup/outro.abditum"
  **E** cancela a sobrescrita
  **Então** nenhum arquivo é modificado
  **E** o caminho atual do cofre permanece inalterado

---

## **Regra:** O salvamento em novo caminho não usa arquivo temporário `.tmp`

### **Cenário:** Gravação direta no caminho final

  **Quando** o usuário salva o cofre em um novo caminho
  **Então** o arquivo é gravado diretamente no caminho final
  **E** nenhum arquivo ".abditum.tmp" é gerado durante o processo
