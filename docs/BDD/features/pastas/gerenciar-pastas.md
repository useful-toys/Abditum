# Funcionalidade: Gerenciar pastas

**Como** usuário do Abditum,
**quero** criar, renomear, mover, reordenar e excluir pastas na hierarquia do cofre,
**para** organizar meus segredos de forma lógica e personalizada.

Pastas são contêineres estruturais que agrupam segredos e subpastas. A hierarquia permite aninhamento em qualquer nível. Pastas não possuem exclusão reversível — sua exclusão promove os filhos para a pasta pai.

---

**Contexto:**

  **Dado** que o cofre está aberto

---

## Criar pasta

### **Regra:** A pasta é criada conforme o contexto de foco

#### **Cenário:** Criar pasta na raiz do cofre

  **Dado** que o foco está na raiz do cofre
  **Quando** o usuário cria uma pasta com o nome "Redes Sociais"
  **Então** a pasta "Redes Sociais" é adicionada ao final da lista de pastas da raiz
  **E** o cofre entra no estado "Cofre Modificado"

#### **Cenário:** Criar pasta dentro de outra pasta

  **Dado** que o foco está na pasta "Financeiro"
  **Quando** o usuário cria uma pasta com o nome "Bancos"
  **Então** a pasta "Bancos" é adicionada ao final da lista de subpastas de "Financeiro"
  **E** o cofre entra no estado "Cofre Modificado"

---

## Renomear pasta

### **Regra:** Renomear altera apenas o nome, preservando identidade e posição

#### **Cenário:** Renomear pasta

  **Dado** que existe a pasta "Serviços"
  **Quando** o usuário renomeia a pasta para "Assinaturas"
  **Então** a pasta passa a se chamar "Assinaturas"
  **E** a identidade, posição e conteúdo da pasta são preservados
  **E** o cofre entra no estado "Cofre Modificado"

---

## Mover pasta

### **Regra:** A pasta movida preserva toda sua hierarquia interna

#### **Cenário:** Mover pasta para dentro de outra pasta

  **Dado** que a pasta "Bancos" está na raiz do cofre
  **E** a pasta "Bancos" contém os segredos "Itaú" e "Bradesco" e a subpasta "Contas"
  **Quando** o usuário move a pasta "Bancos" para dentro da pasta "Financeiro"
  **Então** a pasta "Bancos" é removida da raiz
  **E** a pasta "Bancos" é adicionada ao final da lista de subpastas de "Financeiro"
  **E** os segredos "Itaú" e "Bradesco" e a subpasta "Contas" permanecem dentro de "Bancos"
  **E** o cofre entra no estado "Cofre Modificado"

#### **Cenário:** Mover pasta para a raiz do cofre

  **Dado** que a pasta "Bancos" está dentro da pasta "Financeiro"
  **Quando** o usuário move a pasta "Bancos" para a raiz do cofre
  **Então** a pasta "Bancos" é adicionada ao final da lista de pastas da raiz
  **E** o cofre entra no estado "Cofre Modificado"

---

## Reordenar pasta

### **Regra:** A reordenação altera apenas a posição entre pastas irmãs

#### **Cenário:** Reordenar pasta dentro da raiz

  **Dado** que a raiz contém as pastas "Sites", "Financeiro" e "Serviços" nesta ordem
  **Quando** o usuário move "Serviços" para a primeira posição
  **Então** a ordem das pastas na raiz passa a ser "Serviços", "Sites", "Financeiro"
  **E** o conteúdo interno de cada pasta é preservado
  **E** o cofre entra no estado "Cofre Modificado"

---

## Excluir pasta

### **Regra:** A exclusão de pasta é física e promove os filhos para a pasta pai

#### **Cenário:** Excluir pasta com segredos e subpastas

  **Dado** que a pasta "Financeiro" contém os segredos "Banco X" e "Cartão Y" e a subpasta "Investimentos"
  **E** a pasta "Financeiro" está na raiz do cofre
  **Quando** o usuário exclui a pasta "Financeiro"
  **E** confirma a exclusão
  **Então** a pasta "Financeiro" é removida da hierarquia
  **E** os segredos "Banco X" e "Cartão Y" são adicionados ao final da lista de segredos da raiz
  **E** a subpasta "Investimentos" é adicionada ao final da lista de pastas da raiz
  **E** o cofre entra no estado "Cofre Modificado"

#### **Cenário:** Excluir pasta vazia

  **Dado** que a pasta "Temp" está vazia
  **Quando** o usuário exclui a pasta "Temp"
  **E** confirma a exclusão
  **Então** a pasta "Temp" é removida da hierarquia
  **E** o cofre entra no estado "Cofre Modificado"

#### **Cenário:** Excluir subpasta — filhos promovidos para pasta pai

  **Dado** que a pasta "Bancos" está dentro da pasta "Financeiro"
  **E** a pasta "Bancos" contém o segredo "Itaú"
  **Quando** o usuário exclui a pasta "Bancos"
  **E** confirma a exclusão
  **Então** o segredo "Itaú" é adicionado ao final da lista de segredos de "Financeiro"
  **E** o cofre entra no estado "Cofre Modificado"
