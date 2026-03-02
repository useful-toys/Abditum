# Feature: Área de Transferência

Copiar atributos de forma segura para a área de transferência com limpeza automática.


---

## História: Copiar atributo para área de transferência

**Como** usuário consultando um item,
**Quero** copiar um atributo para a área de transferência com um clique,
**Para que** eu possa usar a informação (ex: senha) sem precisar revelá-la na tela.

### Critérios de Aceite


**Cenário: Copiar atributo não sensível**

- *Dado* que um item está selecionado no painel de detalhes
- *Quando* o usuário clica no ícone de copiar ao lado de um atributo não sensível
- *Então* o valor do atributo é copiado para a área de transferência
- *E* uma notificação breve "Copiado!" é exibida

**Cenário: Copiar atributo sensível (senha)**

- *Dado* que um item está selecionado no painel de detalhes
- *Quando* o usuário clica no ícone de copiar ao lado de um atributo sensível
- *Então* o valor é copiado para a área de transferência sem ser exibido na tela
- *E* uma notificação breve "Copiado! Será limpo em 30s" é exibida
- *E* um temporizador regressivo de 30 segundos é iniciado

---

## História: Limpeza automática da área de transferência

**Como** usuário que copiou uma senha,
**Quero** que a área de transferência seja limpa automaticamente após um tempo,
**Para que** a senha não fique disponível indefinidamente após o uso.

### Critérios de Aceite


**Cenário: Limpar área de transferência após timeout**

- *Dado* que um atributo sensível foi copiado para a área de transferência
- *Quando* 30 segundos se passam
- *Então* a área de transferência é limpa (conteúdo substituído por string vazia)
- *E* uma notificação "Área de transferência limpa" é exibida brevemente

**Cenário: Não limpar se conteúdo foi substituído pelo usuário**

- *Dado* que um atributo sensível foi copiado
- *Quando* o usuário copia outro conteúdo manualmente antes dos 30 segundos
- *Então* o temporizador é cancelado
- *E* a área de transferência NÃO é limpa (o novo conteúdo não é do Abditum)

**Cenário: Limpar ao fechar o cofre**

- *Dado* que a área de transferência contém um valor copiado pelo Abditum
- *Quando* o usuário fecha o cofre
- *Então* a área de transferência é limpa imediatamente, independente do temporizador

---

## História: Limpeza manual da área de transferência

**Como** usuário que copiou uma senha,
**Quero** poder limpar a área de transferência manualmente antes do timeout,
**Para que** eu tenha controle imediato sobre a segurança.

### Critérios de Aceite


**Cenário: Limpar manualmente via interface**

- *Dado* que um atributo sensível foi copiado e o temporizador está ativo
- *Quando* o usuário aciona "Limpar Área de Transferência" na barra de status
- *Então* a área de transferência é limpa imediatamente
- *E* o temporizador é cancelado
- *E* a notificação de temporizador desaparece

---

## Configuração

| Parâmetro                    | Default | Descrição                                     |
|------------------------------|---------|-----------------------------------------------|
| Timeout de limpeza (segundos) | 30      | Tempo até limpeza automática de dados sensíveis |
| Limpeza ao fechar cofre       | true    | Limpar automaticamente ao fechar               |
