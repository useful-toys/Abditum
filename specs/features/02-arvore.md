# Feature: Navegação e Organização da Árvore

Navegar, criar grupos, mover, reordenar e ordenar nós na estrutura hierárquica do cofre.

> O nó raiz é oculto. A aplicação exibe os nós de nível 1 diretamente como uma lista, sem mostrar um nó "Raiz" explícito.


---

## História: Navegar na árvore

**Como** usuário com cofre aberto,
**Quero** expandir e recolher grupos na árvore,
**Para que** eu possa visualizar meus dados organizados hierarquicamente.

### Critérios de Aceite


**Cenário: Expandir grupo**

- *Dado* que a árvore exibe um grupo recolhido com filhos
- *Quando* o usuário clica no grupo
- *Então* os filhos são exibidos abaixo do grupo
- *E* o ícone do grupo indica estado expandido

**Cenário: Recolher grupo**

- *Dado* que a árvore exibe um grupo expandido
- *Quando* o usuário clica no grupo
- *Então* os filhos são ocultados
- *E* o ícone do grupo indica estado recolhido

**Cenário: Selecionar item**

- *Dado* que a árvore está visível
- *Quando* o usuário clica em um nó do tipo item
- *Então* o painel de detalhes exibe os atributos do item
- *E* o nó fica visualmente destacado como selecionado

---

## História: Criar grupo

**Como** usuário com cofre aberto,
**Quero** criar grupos (pastas) na árvore,
**Para que** eu possa organizar meus itens em categorias.

### Critérios de Aceite


**Cenário: Criar grupo no nível 1**

- *Dado* que o cofre está aberto e nenhum nó está selecionado
- *Quando* o usuário aciona "Novo Grupo"
- *Então* um novo grupo é criado como filho direto da raiz (nível 1)
- *E* aparece na lista principal
- *E* o nome do grupo entra em modo de edição imediata

**Cenário: Criar grupo dentro de grupo existente**

- *Dado* que um grupo está selecionado na árvore
- *Quando* o usuário aciona "Novo Grupo"
- *Então* um novo grupo é criado como filho do grupo selecionado
- *E* o grupo pai é expandido automaticamente

**Cenário: Confirmar nome do grupo**

- *Dado* que um novo grupo está em modo de edição de nome
- *Quando* o usuário informa um nome e pressiona Enter ou perde o foco
- *Então* o grupo é salvo com o nome informado

**Cenário: Cancelar criação de grupo**

- *Dado* que um novo grupo está em modo de edição de nome
- *Quando* o usuário pressiona Esc
- *Então* o grupo é removido (criação cancelada)

---

## História: Renomear nó

**Como** usuário com cofre aberto,
**Quero** renomear grupos e itens,
**Para que** eu possa manter a organização atualizada.

### Critérios de Aceite


**Cenário: Renomear via duplo clique**

- *Dado* que um nó está visível na árvore
- *Quando* o usuário dá duplo clique no nome do nó
- *Então* o nome entra em modo de edição
- *Quando* o usuário confirma o novo nome
- *Então* o nó é exibido com o novo nome

**Cenário: Nome em branco não é permitido**

- *Dado* que um nó está em modo de edição de nome
- *Quando* o usuário limpa o nome e confirma
- *Então* o nome não é alterado
- *E* uma mensagem de erro "O nome não pode estar vazio" é exibida

---

## História: Mover nó

**Como** usuário com cofre aberto,
**Quero** mover nós entre grupos,
**Para que** eu possa reorganizar minha estrutura.

### Critérios de Aceite


**Cenário: Mover nó via arrastar e soltar (drag and drop)**

- *Dado* que dois grupos existem na árvore
- *Quando* o usuário arrasta um nó e solta sobre um grupo destino
- *Então* o nó é movido para dentro do grupo destino
- *E* o grupo de origem não contém mais o nó
- *E* o grupo destino está expandido e mostra o nó movido

**Cenário: Mover nó para o nível 1**

- *Dado* que um nó está dentro de um grupo
- *Quando* o usuário o move para a área do nível 1
- *Então* o nó passa a ser filho direto da raiz e aparece na lista principal

**Cenário: Mover grupo para dentro de si mesmo**

- *Dado* que o usuário tenta arrastar um grupo para dentro de um de seus próprios descendentes
- *Então* a operação é rejeitada
- *E* o nó retorna à posição original

---

## História: Remover nó

**Como** usuário com cofre aberto,
**Quero** remover grupos e itens,
**Para que** eu possa excluir dados que não preciso mais.

### Critérios de Aceite


**Cenário: Remover item**

- *Dado* que um item está selecionado
- *Quando* o usuário aciona "Remover"
- *Então* uma confirmação "Remover 'Nome do Item'?" é exibida
- *E* se confirmado, o item é removido da árvore e do cofre

**Cenário: Remover grupo vazio**

- *Dado* que um grupo vazio está selecionado
- *Quando* o usuário aciona "Remover"
- *Então* uma confirmação é exibida
- *E* se confirmado, o grupo é removido

**Cenário: Remover grupo com filhos**

- *Dado* que um grupo com itens e subgrupos está selecionado
- *Quando* o usuário aciona "Remover"
- *Então* uma confirmação explícita "Remover 'Nome' e todos os X itens dentro?" é exibida
- *E* se confirmado, o grupo e todos os descendentes são removidos

**Cenário: Tentar remover o nó raiz**

- *Dado* que o nó raiz é oculto e não selecionável
- *Então* a ação de remover raiz não está disponível na interface

---

## História: Ordenar nós

**Como** usuário com cofre aberto,
**Quero** controlar a ordem dos nós dentro de um grupo,
**Para que** eu possa organizar meus dados da forma que preferir.

### Critérios de Aceite

**Cenário: Reordenar manualmente via drag and drop**

- *Dado* que um grupo contém múltiplos filhos
- *Quando* o usuário arrasta um nó para outra posição dentro do mesmo grupo
- *Então* o nó é inserido na posição indicada
- *E* a nova ordem é persistida no cofre

**Cenário: Aplicar ordenação alfabética em um grupo**

- *Dado* que um grupo está selecionado
- *Quando* o usuário aciona "Ordenar Alfabeticamente"
- *Então* os filhos diretos do grupo são reordenados de A a Z pelo nome
- *E* a nova ordem é persistida no cofre
- *E* subgrupos e itens são ordenados juntos, sem separação por tipo

**Cenário: Ordem padrão é a ordem de inserção**

- *Dado* que itens foram criados em sequência num grupo
- *Quando* nenhuma ordenação foi aplicada
- *Então* os itens aparecem na ordem em que foram criados
