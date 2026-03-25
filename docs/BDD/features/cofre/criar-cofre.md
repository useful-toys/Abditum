# Funcionalidade: Criar novo cofre

**Como** usuário do Abditum,
**quero** criar um novo cofre protegido por senha mestra,
**para** armazenar meus segredos com segurança offline em um arquivo portátil.

Criar um cofre é o ponto de partida do uso do Abditum. O cofre é um arquivo `.abditum` criptografado com AES-256-GCM e protegido por uma senha mestra derivada via Argon2id. Ao ser criado, o cofre já vem populado com pastas e modelos de segredo pré-definidos, pronto para uso imediato.

---

## **Regra:** O cofre é criado com estrutura inicial pré-definida

O novo cofre já nasce com pastas e modelos pré-definidos para facilitar o uso imediato. Esses elementos podem ser editados ou removidos pelo usuário posteriormente.

### **Cenário:** Criar cofre em caminho sem arquivo existente

  **Dado** que a aplicação está no estado inicial, sem cofre ativo
  **Quando** o usuário informa o caminho "/documentos/pessoal.abditum"
  **E** informa a senha mestra "S3nh@F0rte!2026"
  **E** confirma a senha mestra "S3nh@F0rte!2026"
  **Então** um novo arquivo de cofre é gravado em "/documentos/pessoal.abditum"
  **E** o cofre contém as pastas pré-definidas "Sites", "Financeiro" e "Serviços"
  **E** o cofre contém os modelos pré-definidos "Login", "Cartão de Crédito" e "API Key"
  **E** o cofre entra no estado "Cofre Salvo"

---

## **Regra:** A senha mestra exige digitação dupla para confirmação

### **Cenário:** Confirmar senha mestra com sucesso

  **Dado** que a aplicação está no estado inicial, sem cofre ativo
  **Quando** o usuário informa a senha mestra "MinhaSenh@Segura1"
  **E** confirma a senha mestra "MinhaSenh@Segura1"
  **Então** a senha mestra é aceita e a criação do cofre prossegue

### **Cenário:** Falhar ao confirmar senha mestra com valores divergentes

  **Dado** que a aplicação está no estado inicial, sem cofre ativo
  **Quando** o usuário informa a senha mestra "MinhaSenh@Segura1"
  **E** confirma a senha mestra "SenhaDiferente123"
  **Então** a aplicação exibe uma mensagem de erro indicando que as senhas não coincidem
  **E** nenhum arquivo é criado

---

## **Regra:** Sobrescrita de arquivo existente exige confirmação explícita

### **Cenário:** Criar cofre em caminho com arquivo existente, confirmando sobrescrita

  **Dado** que a aplicação está no estado inicial, sem cofre ativo
  **E** já existe um arquivo em "/documentos/antigo.abditum"
  **Quando** o usuário informa o caminho "/documentos/antigo.abditum" e uma senha mestra válida
  **Então** a aplicação solicita confirmação de sobrescrita
  **Quando** o usuário confirma a sobrescrita
  **Então** o arquivo anterior é preservado como backup em "/documentos/antigo.abditum.bak"
  **E** o novo cofre é gravado em "/documentos/antigo.abditum"
  **E** o cofre entra no estado "Cofre Salvo"

### **Cenário:** Criar cofre em caminho com arquivo existente, cancelando sobrescrita

  **Dado** que a aplicação está no estado inicial, sem cofre ativo
  **E** já existe um arquivo em "/documentos/antigo.abditum"
  **Quando** o usuário informa o caminho "/documentos/antigo.abditum" e uma senha mestra válida
  **E** a aplicação solicita confirmação de sobrescrita
  **E** o usuário cancela a sobrescrita
  **Então** nenhum arquivo é modificado
  **E** a aplicação permanece no estado inicial

---

## **Regra:** Backup anterior é rotacionado ao sobrescrever

### **Cenário:** Sobrescrever cofre quando já existe backup anterior

  **Dado** que a aplicação está no estado inicial, sem cofre ativo
  **E** já existe um arquivo em "/documentos/meu.abditum"
  **E** já existe um backup em "/documentos/meu.abditum.bak"
  **Quando** o usuário cria um cofre em "/documentos/meu.abditum" com confirmação de sobrescrita
  **Então** o backup anterior é renomeado temporariamente para "/documentos/meu.abditum.bak2"
  **E** o arquivo existente é copiado para "/documentos/meu.abditum.bak"
  **E** o novo cofre é gravado em "/documentos/meu.abditum"
  **E** o arquivo "/documentos/meu.abditum.bak2" é removido

### **Cenário:** Falha na gravação após geração de backup

  **Dado** que a aplicação está no estado inicial, sem cofre ativo
  **E** já existe um arquivo em "/documentos/meu.abditum"
  **E** já existe um backup em "/documentos/meu.abditum.bak"
  **Quando** o usuário cria um cofre em "/documentos/meu.abditum" com confirmação de sobrescrita
  **E** a gravação do novo arquivo falha
  **Então** o arquivo "/documentos/meu.abditum.bak2" é restaurado para "/documentos/meu.abditum.bak"
  **E** a aplicação exibe uma mensagem de erro informando a falha e a existência de um backup disponível

---

## **Regra:** A criação do cofre não usa arquivo temporário `.tmp`

A criação de um novo cofre grava diretamente no caminho final, pois não se trata do salvamento incremental de um cofre já aberto.

### **Cenário:** Criação grava diretamente no caminho final

  **Dado** que a aplicação está no estado inicial, sem cofre ativo
  **Quando** o usuário cria um cofre em "/documentos/novo.abditum"
  **Então** o arquivo é gravado diretamente em "/documentos/novo.abditum"
  **E** nenhum arquivo ".abditum.tmp" é gerado durante o processo
