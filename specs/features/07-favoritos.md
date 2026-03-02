# Feature: Favoritos

Marcar itens como favoritos para acesso rápido e destaque visual na interface.

---

## História: Marcar item como favorito

**Como** usuário com cofre aberto,
**Quero** marcar itens que uso com frequência como favoritos,
**Para que** eu possa identificá-los visualmente com facilidade.

### Critérios de Aceite

**Cenário: Marcar item como favorito**

- *Dado* que um item está selecionado no painel de detalhes ou na árvore
- *Quando* o usuário aciona o marcador de favorito (ex: ícone de estrela)
- *Então* o item recebe a flag `favorite = true`
- *E* recebe destaque visual na árvore (ex: ícone de estrela junto ao nome)
- *E* a alteração é persistida no cofre

**Cenário: Desmarcar item favorito**

- *Dado* que um item está marcado como favorito
- *Quando* o usuário aciona novamente o marcador de favorito
- *Então* a flag `favorite` volta a `false`
- *E* o destaque visual é removido

**Cenário: Grupos não podem ser favoritados**

- *Dado* que um grupo está selecionado
- *Quando* o usuário tenta acionar o marcador de favorito
- *Então* a ação não está disponível para nós do tipo `group`

---

## História: Visualizar itens favoritos

**Como** usuário com cofre aberto,
**Quero** ver meus itens favoritos destacados,
**Para que** eu os localize rapidamente sem navegar ou buscar.

### Critérios de Aceite

**Cenário: Destaque visual de favoritos na árvore**

- *Dado* que um ou mais itens estão marcados como favoritos
- *Quando* a árvore é exibida
- *Então* itens favoritos possuem um indicador visual distinto dos demais (a ser definido na UI)
- *E* permanecem em sua posição original na árvore — a flag não move o nó

**Cenário: Favoritos preservados ao reabrir o cofre**

- *Dado* que itens foram marcados como favoritos e o cofre foi salvo
- *Quando* o cofre é fechado e reaberto
- *Então* os itens marcados continuam com `favorite = true`
- *E* o destaque visual é restaurado

> **Nota de design**: a forma de apresentação dos favoritos (seção separada no topo, lista lateral, filtro rápido) está em aberto e será decidida durante o design da UI. A spec garante apenas a flag e o destaque visual in-place na árvore.
