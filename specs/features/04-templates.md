# Feature: Templates

Gerenciar modelos predefinidos e personalizados para criação padronizada de itens.


---

## História: Usar template predefinido

**Como** usuário com cofre aberto,
**Quero** criar itens a partir de templates do sistema,
**Para que** eu não precise definir manualmente os atributos comuns.

### Critérios de Aceite


**Cenário: Listar templates disponíveis ao criar item**

- *Dado* que o usuário aciona "Novo Item"
- *Quando* a lista de templates é exibida
- *Então* os templates predefinidos (Site, Conta Bancária, Cartão, Nota Segura, Genérico) aparecem na lista
- *E* templates personalizados criados pelo usuário também aparecem

**Cenário: Aplicar template "Site"**

- *Dado* que o usuário seleciona o template "Site"
- *Quando* o formulário de novo item é aberto
- *Então* os seguintes campos estão presentes e vazios:
  | key   | label | type     | sensitive |
  | url   | URL   | url      | false     |
  | login | Login | text     | false     |
  | senha | Senha | password | true      |

**Cenário: Aplicar template "Conta Bancária"**

- *Dado* que o usuário seleciona o template "Conta Bancária"
- *Quando* o formulário de novo item é aberto
- *Então* os seguintes campos estão presentes e vazios:
  | key             | label          | type     | sensitive |
  | banco           | Banco          | text     | false     |
  | agencia         | Agência        | text     | false     |
  | conta           | Conta Corrente | text     | false     |
  | senha_internet  | Senha Internet | password | true      |
  | senha_cartao    | Senha Cartão   | password | true      |

**Cenário: Aplicar template "Cartão"**

- *Dado* que o usuário seleciona o template "Cartão"
- *Quando* o formulário de novo item é aberto
- *Então* os seguintes campos estão presentes e vazios:
  | key       | label    | type     | sensitive |
  | bandeira  | Bandeira | text     | false     |
  | numero    | Número   | text     | false     |
  | validade  | Validade | date     | false     |
  | cvv       | CVV      | password | true      |
  | senha     | Senha    | password | true      |

---

## História: Criar template personalizado

**Como** usuário com cofre aberto,
**Quero** criar meus próprios templates,
**Para que** eu possa padronizar tipos de itens específicos às minhas necessidades.

### Critérios de Aceite


**Cenário: Criar template a partir do zero**

- *Dado* que o usuário acessa a tela de gerenciamento de templates
- *Quando* aciona "Novo Template" e define nome, ícone e os atributos do schema
- *Então* o template é salvo no cofre
- *E* fica disponível na lista ao criar novos itens

**Cenário: Nome de template duplicado**

- *Dado* que já existe um template com o nome "Servidor"
- *Quando* o usuário tenta criar outro template com o mesmo nome
- *Então* o template não é criado
- *E* uma mensagem "Já existe um template com esse nome" é exibida

---

## História: Editar template personalizado

**Como** usuário com cofre aberto,
**Quero** editar um template personalizado,
**Para que** eu possa ajustar o schema de campos conforme necessidade.

### Critérios de Aceite


**Cenário: Editar template personalizado com sucesso**

- *Dado* que existe um template personalizado
- *Quando* o usuário edita o schema (adiciona, remove ou reordena atributos) e salva
- *Então* o template é atualizado no cofre
- *E* itens existentes criados com esse template NÃO são afetados (template é snapshot na criação)

**Cenário: Tentar editar template predefinido**

- *Dado* que o usuário tenta editar um template com builtin=true
- *Então* a opção de editar está desabilitada
- *E* uma mensagem "Templates do sistema não podem ser editados" é exibida

---

## História: Remover template personalizado

**Como** usuário com cofre aberto,
**Quero** remover templates personalizados desnecessários,
**Para que** a lista de templates se mantenha organizada.

### Critérios de Aceite


**Cenário: Remover template personalizado**

- *Dado* que existe um template personalizado
- *Quando* o usuário aciona "Remover" no template
- *Então* uma confirmação é exibida
- *E* se confirmado, o template é removido do cofre
- *E* itens existentes criados com esse template permanecem intactos

**Cenário: Tentar remover template predefinido**

- *Dado* que o usuário tenta remover um template predefinido (builtin=true)
- *Então* a opção de remover está desabilitada
