# Funcionalidade: Exportar e importar cofre

**Como** usuário do Abditum,
**quero** exportar meu cofre para JSON em texto claro e importar dados de um arquivo JSON,
**para** realizar backup ou migração de dados entre cofres.

---

## Exportar cofre para JSON plain text

A exportação serializa o estado atual do domínio em memória (incluindo alterações não salvas) para um arquivo JSON em texto claro. Por se tratar de uma operação que gera uma cópia não criptografada, exige aviso de segurança e confirmação explícita.

---

**Contexto:**

  **Dado** que o cofre está aberto

---

### **Regra:** A exportação exige aviso de segurança e confirmação

#### **Cenário:** Exportar cofre com confirmação

  **Quando** o usuário inicia a exportação do cofre
  **Então** a aplicação exibe aviso sobre o risco de segurança de gerar uma cópia não criptografada
  **Quando** o usuário confirma a exportação
  **Então** a aplicação serializa o domínio para JSON em texto claro no destino escolhido
  **E** o estado do cofre não é alterado

#### **Cenário:** Cancelar exportação

  **Quando** o usuário inicia a exportação do cofre
  **E** cancela após ver o aviso de segurança
  **Então** nenhum arquivo é gerado
  **E** o estado do cofre não é alterado

### **Regra:** Exportação em "Cofre Modificado" inclui aviso adicional

#### **Cenário:** Exportar cofre em estado "Cofre Modificado"

  **Dado** que o cofre está no estado "Cofre Modificado"
  **Quando** o usuário inicia a exportação
  **Então** a aplicação exibe alerta informando que a exportação incluirá alterações ainda não salvas
  **E** exibe o aviso padrão de segurança

---

## Importar cofre de JSON plain text

A importação incorpora dados de um arquivo JSON plain text ao cofre ativo. Conflitos de identidade e nome são resolvidos automaticamente segundo regras específicas.

---

**Contexto:**

  **Dado** que o cofre está aberto

---

### **Regra:** Pastas com mesma identidade são mescladas silenciosamente

#### **Cenário:** Importar pasta com mesma identidade de pasta existente

  **Dado** que o cofre contém a pasta "Sites" com identidade "abc123"
  **E** o arquivo JSON importado contém uma pasta com identidade "abc123"
  **Quando** o usuário importa o arquivo JSON
  **Então** as hierarquias das pastas são mescladas
  **E** nenhuma mensagem de conflito de pastas é exibida
  **E** o cofre entra no estado "Cofre Modificado"

### **Regra:** Segredos com colisão de identidade recebem nova identidade

#### **Cenário:** Importar segredo com mesma identidade de segredo existente

  **Dado** que o cofre contém o segredo "Gmail" com identidade "xyz789"
  **E** o arquivo JSON importado contém um segredo com identidade "xyz789"
  **Quando** o usuário importa o arquivo JSON
  **Então** o segredo importado recebe uma nova identidade, preservando seus demais dados
  **E** o cofre entra no estado "Cofre Modificado"

### **Regra:** Segredos com colisão de nome recebem sufixo incremental

#### **Cenário:** Importar segredo com nome duplicado na mesma pasta

  **Dado** que o cofre contém o segredo "Gmail" na pasta "Sites"
  **E** o arquivo JSON importado contém um segredo "Gmail" na mesma pasta de destino
  **Quando** o usuário importa o arquivo JSON
  **Então** o segredo importado é renomeado para "Gmail (1)"
  **E** a aplicação exibe mensagem informando que segredos conflitantes por nome foram importados com nomes sufixados
  **E** o cofre entra no estado "Cofre Modificado"

#### **Cenário:** Importar múltiplos segredos com nomes duplicados

  **Dado** que o cofre contém o segredo "API" na raiz
  **E** já existe um segredo "API (1)" na raiz
  **E** o arquivo JSON importado contém um segredo "API" na raiz
  **Quando** o usuário importa o arquivo JSON
  **Então** o segredo importado é renomeado para "API (2)"
  **E** a aplicação exibe mensagem de conflito

### **Regra:** Modelos com mesma identidade são sobrepostos silenciosamente

#### **Cenário:** Importar modelo com mesma identidade de modelo existente

  **Dado** que o cofre contém o modelo "Login" com identidade "mod001"
  **E** o arquivo JSON importado contém um modelo com identidade "mod001"
  **Quando** o usuário importa o arquivo JSON
  **Então** o modelo existente é substituído pelo modelo importado
  **E** nenhuma mensagem de conflito de modelos é exibida
  **E** o cofre entra no estado "Cofre Modificado"
