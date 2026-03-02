# Feature: Gerenciamento de Itens

Criar, visualizar, editar e remover itens com seus atributos.


---

## História: Criar item

**Como** usuário com cofre aberto,
**Quero** criar um novo item em um grupo,
**Para que** eu possa registrar um novo dado pessoal.

### Critérios de Aceite


**Cenário: Criar item a partir de template predefinido**

- *Dado* que um grupo está selecionado na árvore
- *Quando* o usuário aciona "Novo Item" e seleciona o template "Site"
- *Então* um novo item é criado com os atributos padrão do template (url, login, senha)
- *E* o item entra em modo de edição
- *E* os campos são exibidos vazios aguardando preenchimento

**Cenário: Criar item genérico sem template**

- *Dado* que um grupo está selecionado na árvore
- *Quando* o usuário aciona "Novo Item" e seleciona "Genérico"
- *Então* um novo item é criado sem atributos predefinidos
- *E* o usuário pode adicionar atributos livremente

**Cenário: Salvar novo item**

- *Dado* que um novo item está em modo de edição com nome e atributos preenchidos
- *Quando* o usuário aciona "Salvar"
- *Então* o item é persistido no cofre em memória
- *E* o cofre é marcado como modificado (não salvo em disco)
- *E* o item aparece na árvore dentro do grupo selecionado

**Cenário: Cancelar criação de item**

- *Dado* que um novo item está em modo de edição
- *Quando* o usuário aciona "Cancelar"
- *Então* o item é descartado
- *E* nenhuma alteração é feita no cofre

---

## História: Visualizar item

**Como** usuário com cofre aberto,
**Quero** visualizar os atributos de um item,
**Para que** eu possa consultar os dados armazenados.

### Critérios de Aceite


**Cenário: Exibir atributos com valores sensíveis mascarados**

- *Dado* que um item com atributos sensíveis está selecionado
- *Quando* o painel de detalhes é exibido
- *Então* atributos com sensitive=true são exibidos como "••••••••"
- *E* atributos com sensitive=false são exibidos em texto claro

**Cenário: Revelar atributo sensível**

- *Dado* que um atributo sensível está mascarado no painel de detalhes
- *Quando* o usuário clica no ícone de "revelar" ao lado do atributo
- *Então* o valor é exibido em texto claro
- *E* o ícone muda para "ocultar"
- *Quando* o usuário clica em "ocultar"
- *Então* o valor volta a ser mascarado

---

## História: Editar item

**Como** usuário com cofre aberto,
**Quero** editar um item existente,
**Para que** eu possa atualizar dados desatualizados.

### Critérios de Aceite


**Cenário: Editar nome do item**

- *Dado* que um item está selecionado
- *Quando* o usuário aciona "Editar"
- *Então* o painel de edição é exibido com todos os campos preenchidos
- *Quando* o usuário altera o nome e salva
- *Então* o item exibe o novo nome na árvore

**Cenário: Editar valor de atributo existente**

- *Dado* que o painel de edição de um item está aberto
- *Quando* o usuário modifica o valor de um atributo e salva
- *Então* o novo valor é armazenado no cofre
- *E* o painel de detalhes exibe o valor atualizado

**Cenário: Adicionar novo atributo a item existente**

- *Dado* que o painel de edição de um item está aberto
- *Quando* o usuário aciona "Adicionar Atributo", informa chave, label, tipo e valor
- *Então* o novo atributo é adicionado à lista
- *E* ao salvar, o atributo faz parte do item

**Cenário: Remover atributo de item**

- *Dado* que o painel de edição de um item está aberto com pelo menos um atributo
- *Quando* o usuário aciona "Remover" ao lado de um atributo
- *Então* o atributo é removido da lista de edição
- *E* ao salvar, o atributo não faz mais parte do item

**Cenário: Reordenar atributos**

- *Dado* que o painel de edição de um item está aberto com múltiplos atributos
- *Quando* o usuário arrasta um atributo para outra posição
- *Então* a ordem dos atributos é atualizada visualmente
- *E* ao salvar, a nova ordem é persistida

---

## História: Salvar cofre em disco

**Como** usuário com alterações pendentes,
**Quero** salvar o cofre,
**Para que** as alterações sejam persistidas no arquivo.

### Critérios de Aceite


**Cenário: Salvar manualmente**

- *Dado* que o cofre possui alterações não salvas
- *Quando* o usuário aciona "Salvar" (Ctrl+S)
- *Então* o cofre é serializado, cifrado e escrito no arquivo atomicamente
- *E* o indicador de "não salvo" desaparece

**Cenário: Arquivo salvo atomicamente**

- *Dado* que o cofre está sendo salvo
- *Quando* o processo de escrita é realizado
- *Então* a escrita ocorre em um arquivo temporário primeiro
- *E* o arquivo temporário substitui o original somente após conclusão com sucesso
- *E* em caso de falha na escrita, o arquivo original é preservado intacto
