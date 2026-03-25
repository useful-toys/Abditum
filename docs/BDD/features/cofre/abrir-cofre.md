# Funcionalidade: Abrir cofre existente

**Como** usuário do Abditum,
**quero** abrir um cofre existente informando o caminho do arquivo e minha senha mestra,
**para** acessar meus segredos armazenados com segurança.

Abrir um cofre envolve validar a assinatura do arquivo, selecionar o perfil criptográfico correto com base na versão do formato, derivar a chave a partir da senha mestra, descriptografar o payload e carregar o domínio em memória. Se o formato for de uma versão anterior, os dados são migrados em memória para o modelo corrente.

---

**Contexto:**

  **Dado** que a aplicação está no estado inicial, sem cofre ativo

---

## **Regra:** O arquivo deve ser um cofre Abditum válido

### **Cenário:** Abrir cofre válido com senha mestra correta

  **Quando** o usuário informa o caminho de um arquivo ".abditum" válido
  **E** informa a senha mestra correta
  **Então** a aplicação valida a assinatura "magic", deriva a chave, descriptografa o payload e carrega o domínio em memória
  **E** o cofre entra no estado "Cofre Salvo"

### **Cenário:** Rejeitar arquivo que não é um cofre Abditum

  **Quando** o usuário informa o caminho de um arquivo que não possui a assinatura "magic" do Abditum
  **Então** a aplicação exibe uma mensagem de erro indicando que o arquivo não é um cofre Abditum
  **E** a aplicação permanece no estado inicial

### **Cenário:** Rejeitar cofre com versão de formato superior à suportada

  **Quando** o usuário informa o caminho de um arquivo ".abditum" com versão de formato superior à suportada pela aplicação
  **Então** a aplicação exibe uma mensagem de erro de incompatibilidade de versão
  **E** a aplicação permanece no estado inicial

---

## **Regra:** Senha mestra incorreta impede a abertura

### **Cenário:** Falhar ao abrir cofre com senha mestra incorreta

  **Quando** o usuário informa o caminho de um arquivo ".abditum" válido
  **E** informa uma senha mestra incorreta
  **Então** a aplicação exibe uma mensagem de erro indicando que a senha está incorreta ou o arquivo está corrompido
  **E** a aplicação permanece no estado inicial

---

## **Regra:** Cofres de versões anteriores são migrados em memória

### **Cenário:** Abrir cofre de versão anterior suportada

  **Quando** o usuário informa o caminho de um cofre criado em uma versão anterior do formato
  **E** informa a senha mestra correta
  **Então** a aplicação seleciona o perfil Argon2id histórico correspondente à versão do arquivo
  **E** descriptografa o payload e migra os dados em memória para o modelo corrente
  **E** o cofre entra no estado "Cofre Salvo"
