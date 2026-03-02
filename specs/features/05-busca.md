# Feature: Busca

Localizar itens no cofre por nome ou conteúdo de atributos.


---

## História: Buscar itens por nome

**Como** usuário com cofre aberto,
**Quero** digitar um termo de busca e ver os itens correspondentes,
**Para que** eu possa encontrar rapidamente um dado sem navegar manualmente pela árvore.

### Critérios de Aceite


**Cenário: Busca com resultados encontrados**

- *Dado* que o cofre contém itens
- *Quando* o usuário digita um termo na caixa de busca
- *Então* a árvore é filtrada exibindo apenas itens cujo nome contém o termo (case-insensitive)
- *E* os grupos que contêm os itens encontrados são exibidos como contexto
- *E* grupos sem resultados são ocultados

**Cenário: Busca sem resultados**

- *Dado* que nenhum item no cofre corresponde ao termo buscado
- *Quando* o usuário digita o termo
- *Então* a árvore exibe uma mensagem "Nenhum resultado encontrado"

**Cenário: Limpar busca**

- *Dado* que o usuário realizou uma busca com resultados filtrados
- *Quando* o usuário limpa a caixa de busca ou pressiona Esc
- *Então* a árvore volta ao estado completo, com todos os nós visíveis

---

## História: Buscar em atributos não sensíveis

**Como** usuário com cofre aberto,
**Quero** que a busca também encontre itens por valores de atributos visíveis,
**Para que** eu possa buscar por nome de banco, URL ou login sem saber o nome do item.

### Critérios de Aceite


**Cenário: Busca retorna item por valor de atributo não sensível**

- *Dado* que existe um item cujo atributo login="joao@email.com"
- *Quando* o usuário busca por "joao"
- *Então* o item é exibido nos resultados
- *E* o atributo correspondente é destacado visualmente no resultado

**Cenário: Busca não retorna resultados de atributos sensíveis**

- *Dado* que existe um item cujo atributo senha="s3cr3t"
- *Quando* o usuário busca por "s3cr3t"
- *Então* o item NÃO aparece nos resultados (senhas não são indexadas para busca)

---

## História: Navegar nos resultados de busca

**Como** usuário com cofre aberto,
**Quero** navegar pelos resultados de busca pelo teclado,
**Para que** eu possa selecionar o item desejado rapidamente sem usar o mouse.

### Critérios de Aceite


**Cenário: Navegar entre resultados com teclado**

- *Dado* que a busca retornou múltiplos resultados
- *Quando* o usuário pressiona a tecla ↓ ou ↑
- *Então* o próximo ou anterior resultado é selecionado
- *E* o painel de detalhes exibe o item selecionado

**Cenário: Selecionar resultado e abrir detalhes**

- *Dado* que um resultado de busca está selecionado
- *Quando* o usuário pressiona Enter
- *Então* o item é selecionado na árvore
- *E* o painel de detalhes exibe seus atributos
- *E* a busca é limpa, revelando a posição do item na árvore
