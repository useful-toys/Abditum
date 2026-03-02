# Feature: Gerenciamento do Cofre

Gerenciar o ciclo de vida do cofre: criação, abertura, fechamento e alteração de senha mestra.


---

## História: Criar novo cofre

**Como** usuário novo,
**Quero** criar um cofre protegido por senha,
**Para que** eu possa começar a armazenar meus dados pessoais de forma segura.

### Critérios de Aceite


**Cenário: Criar cofre com sucesso**

- *Dado* que nenhum cofre está aberto
- *Quando* o usuário informa um nome de arquivo, uma senha mestra e confirma a senha
- *Então* um arquivo .abt é criado no caminho escolhido
- *E* o cofre é aberto automaticamente
- *E* quatro grupos predefinidos são criados no nível 1: Sites, E-mails, Contas Bancárias, Telefones
- *E* os grupos são exibidos como uma lista (raiz oculta)

**Cenário: Confirmar senha não confere**

- *Dado* que nenhum cofre está aberto
- *Quando* o usuário informa uma senha mestra e uma confirmação diferente
- *Então* o cofre não é criado
- *E* uma mensagem de erro "As senhas não coincidem" é exibida

**Cenário: Senha mestra muito curta**

- *Dado* que nenhum cofre está aberto
- *Quando* o usuário informa uma senha com menos de 8 caracteres
- *Então* o cofre não é criado
- *E* uma mensagem de aviso "A senha deve ter pelo menos 8 caracteres" é exibida

---

## História: Abrir cofre existente

**Como** usuário com cofre já criado,
**Quero** abrir o cofre informando minha senha,
**Para que** eu possa acessar meus dados.

### Critérios de Aceite


**Cenário: Abrir cofre com senha correta**

- *Dado* que existe um arquivo .abt válido
- *Quando* o usuário seleciona o arquivo e informa a senha mestra correta
- *Então* o cofre é decifrado e carregado em memória
- *E* a árvore é exibida na interface

**Cenário: Senha incorreta**

- *Dado* que existe um arquivo .abt válido
- *Quando* o usuário informa uma senha incorreta
- *Então* o cofre não é aberto
- *E* uma mensagem de erro "Senha incorreta ou arquivo inválido" é exibida
- *E* nenhum dado parcial é exibido

**Cenário: Arquivo corrompido ou inválido**

- *Dado* que o usuário seleciona um arquivo que não é um cofre Abditum
- *Quando* tenta abrir o arquivo
- *Então* uma mensagem de erro "Arquivo inválido ou não reconhecido" é exibida

**Cenário: Aplicação não persiste localização de cofres**

- *Dado* que o usuário fechou a aplicação após usar um cofre
- *Quando* abre a aplicação novamente
- *Então* nenhuma referência ao cofre anterior é exibida ou armazenada em disco
- *E* o usuário deve selecionar o arquivo manualmente

---

## História: Fechar cofre

**Como** usuário com cofre aberto,
**Quero** fechar o cofre,
**Para que** os dados sejam removidos da memória.

### Critérios de Aceite


**Cenário: Fechar cofre sem alterações pendentes**

- *Dado* que o cofre está aberto e sem alterações não salvas
- *Quando* o usuário aciona "Fechar Cofre"
- *Então* a árvore é limpa da interface
- *E* todos os dados são removidos da memória
- *E* a tela inicial (abrir/criar cofre) é exibida

**Cenário: Fechar cofre com alterações não salvas**

- *Dado* que o cofre está aberto com alterações não salvas
- *Quando* o usuário tenta fechar o cofre
- *Então* uma caixa de diálogo pergunta "Salvar antes de fechar?"
- *E* se o usuário confirmar, o cofre é salvo e fechado
- *E* se o usuário recusar, o cofre é fechado sem salvar
- *E* se o usuário cancelar, o cofre permanece aberto

---

## História: Alterar senha mestra

**Como** usuário com cofre aberto,
**Quero** alterar a senha mestra do cofre,
**Para que** eu possa manter a segurança ao suspeitar de comprometimento.

### Critérios de Aceite


**Cenário: Alterar senha com sucesso**

- *Dado* que o cofre está aberto
- *Quando* o usuário informa a senha atual, a nova senha e a confirmação
- *E* a senha atual está correta
- *E* a nova senha e confirmação coincidem
- *Então* o cofre é re-cifrado com a nova chave derivada da nova senha
- *E* o arquivo é sobrescrito atomicamente
- *E* uma mensagem de confirmação "Senha alterada com sucesso" é exibida

**Cenário: Senha atual incorreta ao alterar**

- *Dado* que o cofre está aberto
- *Quando* o usuário informa uma senha atual incorreta
- *Então* a senha não é alterada
- *E* uma mensagem de erro "Senha atual incorreta" é exibida
