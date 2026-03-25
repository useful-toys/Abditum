# Funcionalidade: Navegar hierarquia do cofre

**Como** usuário do Abditum,
**quero** navegar pela árvore de pastas e segredos do cofre,
**para** localizar e acessar meus segredos rapidamente.

A navegação é somente leitura — não altera o conteúdo do cofre nem seu estado persistido. A hierarquia exibe segredos primeiro e subpastas depois dentro de cada coleção, conforme a ordem persistida.

---

**Contexto:**

  **Dado** que o cofre está aberto

---

## **Regra:** A hierarquia exibe segredos primeiro, depois subpastas

### **Cenário:** Visualizar hierarquia com segredos e pastas

  **Dado** que a raiz contém os segredos "Nota Rápida" e "API interna" e as pastas "Sites" e "Financeiro"
  **Então** a hierarquia exibe na raiz: "Nota Rápida", "API interna", "Sites", "Financeiro" — nesta ordem

---

## **Regra:** Pastas podem ser expandidas e colapsadas

### **Cenário:** Expandir pasta na hierarquia

  **Dado** que a pasta "Sites" contém o segredo "Gmail" e a subpasta "Redes Sociais"
  **Quando** o usuário expande a pasta "Sites"
  **Então** os filhos "Gmail" e "Redes Sociais" são exibidos abaixo de "Sites"

### **Cenário:** Colapsar pasta na hierarquia

  **Dado** que a pasta "Sites" está expandida
  **Quando** o usuário colapsa a pasta "Sites"
  **Então** os filhos de "Sites" são ocultados

---

## Visualizar segredo

### **Regra:** Campos sensíveis são ocultos por padrão

#### **Cenário:** Selecionar segredo para visualização

  **Quando** o usuário seleciona o segredo "Gmail"
  **Então** os detalhes do segredo são exibidos no painel do segredo
  **E** campos do tipo "texto sensível" são apresentados ocultos
  **E** a visualização não altera o conteúdo do segredo

---

## Visualizar e ocultar campo sensível

### **Regra:** A reocultação é automática conforme configuração do cofre

#### **Cenário:** Exibir campo sensível temporariamente

  **Dado** que o segredo "Gmail" está sendo visualizado
  **E** o campo "Password" está oculto
  **Quando** o usuário solicita exibir o campo "Password"
  **Então** o valor do campo é revelado temporariamente

#### **Cenário:** Ocultar campo sensível manualmente

  **Dado** que o campo "Password" está exibido temporariamente
  **Quando** o usuário solicita ocultar o campo "Password"
  **Então** o valor do campo volta a ser ocultado

#### **Cenário:** Reocultação automática de campo sensível

  **Dado** que o tempo de reocultação está configurado para 15 segundos
  **E** o campo "Password" foi revelado temporariamente
  **Quando** passam 15 segundos sem que o usuário oculte manualmente
  **Então** o valor do campo é reocultado automaticamente
