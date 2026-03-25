# Especificações de Caso de Uso — Abditum

| Item            | Detalhe                        |
|-----------------|--------------------------------|
| Projeto         | Abditum                       |
| Versão          | 1.0                           |
| Data            | 2026-03-25                     |
| Status          | Aprovado                       |

---

## Ator

| Ator             | Descrição                                                                 |
|------------------|---------------------------------------------------------------------------|
| Usuário          | Pessoa que utiliza o Abditum para armazenar e gerenciar segredos          |

---

## Índice de Casos de Uso

| ID     | Caso de Uso                                | Categoria               |
|--------|--------------------------------------------|--------------------------|
| UC-01  | Criar novo cofre                           | Ciclo de vida do cofre   |
| UC-02  | Abrir cofre existente                      | Ciclo de vida do cofre   |
| UC-03  | Salvar cofre                               | Ciclo de vida do cofre   |
| UC-04  | Salvar cofre em novo caminho               | Ciclo de vida do cofre   |
| UC-05  | Descartar alterações e recarregar cofre    | Ciclo de vida do cofre   |
| UC-06  | Alterar senha mestra                       | Ciclo de vida do cofre   |
| UC-07  | Bloquear cofre                             | Ciclo de vida do cofre   |
| UC-08  | Configurar o cofre                         | Ciclo de vida do cofre   |
| UC-09  | Exportar cofre                             | Ciclo de vida do cofre   |
| UC-10  | Importar cofre                             | Ciclo de vida do cofre   |
| UC-11  | Sair da aplicação                          | Ciclo de vida do cofre   |
| UC-12  | Navegar hierarquia do cofre                | Navegação                |
| UC-13  | Visualizar segredo                         | Navegação                |
| UC-14  | Visualizar/ocultar campo sensível          | Navegação                |
| UC-15  | Buscar segredos                            | Navegação                |
| UC-16  | Criar segredo                              | Gerenciamento de segredos|
| UC-17  | Duplicar segredo                           | Gerenciamento de segredos|
| UC-18  | Editar segredo (edição padrão)             | Gerenciamento de segredos|
| UC-19  | Editar segredo (edição avançada)           | Gerenciamento de segredos|
| UC-20  | Favoritar/desfavoritar segredo             | Gerenciamento de segredos|
| UC-21  | Excluir segredo reversivelmente            | Gerenciamento de segredos|
| UC-22  | Restaurar segredo da Lixeira               | Gerenciamento de segredos|
| UC-23  | Mover segredo                              | Gerenciamento de segredos|
| UC-24  | Reordenar segredo                          | Gerenciamento de segredos|
| UC-25  | Copiar campo de segredo                    | Área de transferência    |
| UC-26  | Criar pasta                                | Hierarquia               |
| UC-27  | Renomear pasta                             | Hierarquia               |
| UC-28  | Mover pasta                                | Hierarquia               |
| UC-29  | Reordenar pasta                            | Hierarquia               |
| UC-30  | Excluir pasta                              | Hierarquia               |
| UC-31  | Criar modelo de segredo                    | Modelos de segredo       |
| UC-32  | Editar modelo de segredo                   | Modelos de segredo       |
| UC-33  | Remover modelo de segredo                  | Modelos de segredo       |
| UC-34  | Criar modelo a partir de segredo existente | Modelos de segredo       |

---

## Especificações Detalhadas

---

### UC-01: Criar novo cofre

**Ator:** Usuário

**Pré-condições:**
- A aplicação está no estado Inicial (sem cofre ativo)

**Fluxo Principal:**
1. O Usuário inicia a ação de criar novo cofre.
2. O Usuário informa o caminho de destino para o novo cofre.
3. O Usuário informa a senha mestra.
4. O Usuário confirma a senha mestra (digitação dupla).
5. A aplicação exibe aviso categórico sobre a irrecuperabilidade da senha mestra (Conhecimento Zero).
6. A aplicação popula a estrutura inicial do cofre com pastas e modelos pré-definidos.
7. A aplicação grava o novo cofre no caminho informado.
8. O cofre entra em estado Cofre Salvo.

**Fluxos Alternativos:**

*FA-01: Arquivo já existe no caminho informado*
1. A aplicação exige confirmação explícita de sobrescrita.
2. Se o Usuário confirmar, a aplicação gera backup do arquivo existente (`.bak`) antes de sobrescrever.
3. O fluxo retorna ao passo 7 do fluxo principal.

*FA-02: Usuário cancela a criação*
1. A aplicação descarta a operação sem efeitos.
2. A aplicação retorna ao estado Inicial.

**Fluxos de Exceção:**

*FE-01: Falha na gravação do arquivo*
1. A aplicação exibe mensagem de erro descrevendo a falha.
2. Se um backup foi gerado, a aplicação informa que existe backup disponível para intervenção manual.
3. A aplicação retorna ao estado anterior.

**Pós-condições:**
- O cofre está ativo no estado Cofre Salvo, com estrutura inicial (pastas e modelos pré-definidos).

**Regras de Negócio:** RN-01, RN-22, RN-27

---

### UC-02: Abrir cofre existente

**Ator:** Usuário

**Pré-condições:**
- A aplicação está no estado Inicial ou o cofre anterior foi bloqueado

**Fluxo Principal:**
1. O Usuário informa o caminho do cofre.
2. A aplicação valida o arquivo (assinatura e versão do formato).
3. O Usuário informa a senha mestra.
4. A aplicação deriva a chave, valida o conteúdo criptografado e carrega os dados em memória.
5. O cofre entra em estado Cofre Salvo.

**Fluxos Alternativos:**

*FA-01: Cofre em formato de versão anterior*
1. A aplicação migra os dados em memória para o formato corrente.
2. O fluxo continua no passo 5.

**Fluxos de Exceção:**

*FE-01: Arquivo inválido (não é um cofre Abditum)*
1. A aplicação informa que o arquivo não é reconhecido como cofre Abditum.

*FE-02: Versão do formato superior à suportada*
1. A aplicação informa incompatibilidade de versão e orienta o Usuário a usar uma versão mais recente do Abditum.

*FE-03: Senha mestra incorreta ou dados corrompidos*
1. A aplicação informa que não foi possível abrir o cofre (sem distinguir senha incorreta de corrupção, por segurança).

**Pós-condições:**
- O cofre está ativo no estado Cofre Salvo.

**Regras de Negócio:** RN-01

---

### UC-03: Salvar cofre

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo em estado Cofre Modificado

**Fluxo Principal:**
1. O Usuário inicia a ação de salvar.
2. A aplicação grava o cofre em arquivo temporário.
3. A aplicação gera backup do arquivo atual (`.bak`).
4. A aplicação substitui o arquivo do cofre pelo temporário via renomeação atômica.
5. O cofre entra em estado Cofre Salvo.
6. Segredos na Lixeira são excluídos permanentemente.

**Fluxos de Exceção:**

*FE-01: Falha na gravação*
1. A aplicação remove o arquivo temporário e restaura o backup se necessário.
2. A aplicação informa a falha e a existência de backup disponível.

**Pós-condições:**
- O cofre está sincronizado com o arquivo no disco. A Lixeira está vazia.

**Regras de Negócio:** RN-05, RN-13

---

### UC-04: Salvar cofre em novo caminho

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo (em qualquer subestado)

**Fluxo Principal:**
1. O Usuário inicia a ação de salvar em novo caminho.
2. O Usuário informa o novo caminho de destino.
3. A aplicação valida a possibilidade de gravação.
4. A aplicação grava o cofre diretamente no novo caminho.
5. O novo caminho passa a ser o caminho atual do cofre.
6. Segredos na Lixeira são excluídos permanentemente.
7. O cofre entra em estado Cofre Salvo.

**Fluxos Alternativos:**

*FA-01: Já existe arquivo no destino*
1. A aplicação exige confirmação de sobrescrita.
2. A aplicação gera backup do arquivo existente no destino (`.bak`) antes de sobrescrever.
3. O fluxo retorna ao passo 4.

**Fluxos de Exceção:**

*FE-01: Falha na gravação*
1. A aplicação informa a falha e, se backup foi gerado, indica sua existência para intervenção manual.

**Pós-condições:**
- O cofre está salvo no novo caminho, que se torna o caminho atual.

**Regras de Negócio:** RN-05, RN-13

---

### UC-05: Descartar alterações e recarregar cofre

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo em estado Cofre Modificado

**Fluxo Principal:**
1. O Usuário inicia a ação de descartar alterações.
2. A aplicação exige confirmação.
3. Após confirmação, a aplicação reabre o arquivo do cofre, repete validação e descriptografia.
4. O cofre entra em estado Cofre Salvo.

**Fluxos Alternativos:**

*FA-01: Usuário cancela*
1. A aplicação mantém o estado Cofre Modificado sem alterações.

**Pós-condições:**
- O cofre reflete o conteúdo do último salvamento no disco.

**Regras de Negócio:** RN-13

---

### UC-06: Alterar senha mestra

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo (em qualquer subestado)

**Fluxo Principal:**
1. O Usuário inicia a ação de alterar a senha mestra.
2. A aplicação solicita a nova senha mestra.
3. O Usuário confirma a nova senha (digitação dupla).
4. A aplicação prepara a regravação criptográfica com a nova credencial.
5. A aplicação segue o fluxo de salvar o cofre (UC-03).

**Fluxos de Exceção:**

*FE-01: Confirmação não corresponde*
1. A aplicação informa que as senhas não coincidem e solicita nova digitação.

**Pós-condições:**
- O cofre está salvo com a nova senha mestra.

**Regras de Negócio:** RN-28

---

### UC-07: Bloquear cofre

**Ator:** Usuário / Sistema (inatividade)

**Pré-condições:**
- Cofre ativo (em qualquer subestado)

**Fluxo Principal:**
1. O bloqueio é disparado pelo Usuário (manual) ou pelo sistema (inatividade).
2. A aplicação fecha logicamente o cofre.
3. A aplicação minimiza dados sensíveis em memória e limpa a área de transferência.
4. Alterações não salvas são descartadas silenciosamente.
5. A aplicação retorna ao fluxo de abrir cofre existente (UC-02), assumindo o mesmo caminho.

**Fluxos Alternativos:**

*FA-01: Bloqueio por inatividade com alerta*
1. Ao atingir 75% do tempo configurado de inatividade, a aplicação emite alerta de bloqueio iminente.
2. Se o Usuário realizar atividade válida (teclado ou clique), o cronômetro é reiniciado.
3. Se o Usuário não realizar atividade, o bloqueio é efetivado ao atingir o tempo limite.

**Pós-condições:**
- O cofre não está mais ativo. A aplicação aguarda nova autenticação.

**Regras de Negócio:** RN-14, RN-15, RN-20, RN-21

---

### UC-08: Configurar o cofre

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo

**Fluxo Principal:**
1. O Usuário inicia a edição das configurações.
2. A aplicação apresenta as configurações atuais: tempo de bloqueio por inatividade, tempo de reocultação de campos sensíveis, tempo de limpeza da área de transferência.
3. O Usuário altera os valores desejados e confirma.
4. As alterações passam a valer para a sessão corrente.
5. O cofre entra em estado Cofre Modificado.

**Fluxos Alternativos:**

*FA-01: Usuário cancela*
1. As configurações permanecem inalteradas.

**Pós-condições:**
- As novas configurações estão em vigor e serão persistidas no próximo salvamento.

**Regras de Negócio:** RN-27

---

### UC-09: Exportar cofre

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo

**Fluxo Principal:**
1. O Usuário inicia a exportação.
2. A aplicação exibe aviso sobre o risco de gerar uma cópia não criptografada.
3. Se o cofre estiver em estado Cofre Modificado, a aplicação informa que a exportação incluirá alterações não salvas.
4. O Usuário confirma.
5. A aplicação serializa o cofre em formato legível no destino escolhido pelo Usuário.

**Fluxos Alternativos:**

*FA-01: Usuário cancela*
1. A exportação é abortada sem efeitos.

**Pós-condições:**
- Um arquivo legível foi gerado no destino. O cofre ativo não é alterado.

---

### UC-10: Importar cofre

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo

**Fluxo Principal:**
1. O Usuário inicia a importação e seleciona o arquivo de origem.
2. A aplicação lê o conteúdo e resolve conflitos:
   - Pastas com mesma identidade são mescladas silenciosamente.
   - Segredos com identidade conflitante recebem nova identidade.
   - Segredos com nome conflitante na mesma pasta recebem sufixo numérico incremental.
   - Modelos com mesma identidade são sobrepostos pelo importado.
3. Se houver conflitos de nome em segredos, a aplicação informa os ajustes realizados.
4. O Usuário confirma a incorporação.
5. O cofre entra em estado Cofre Modificado.

**Fluxos Alternativos:**

*FA-01: Usuário cancela*
1. A importação é abortada sem efeitos.

**Pós-condições:**
- Os dados importados foram incorporados ao cofre ativo.

**Regras de Negócio:** RN-16, RN-17, RN-18, RN-19

---

### UC-11: Sair da aplicação

**Ator:** Usuário

**Pré-condições:**
- Nenhuma

**Fluxo Principal:**
1. O Usuário inicia a ação de sair.
2. A aplicação solicita confirmação de encerramento.
3. A aplicação encerra.

**Fluxos Alternativos:**

*FA-01: Cofre ativo em estado Cofre Modificado*
1. A aplicação oferece ao Usuário: Salvar, Sair sem Salvar (Descartar) ou Voltar.
2. Se Salvar: a aplicação segue o fluxo UC-03 e encerra após salvamento bem-sucedido.
3. Se Descartar: a aplicação encerra sem salvar.
4. Se Voltar: a aplicação cancela o encerramento.

**Pós-condições:**
- A aplicação está encerrada. Área de transferência limpa.

---

### UC-12: Navegar hierarquia do cofre

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo

**Fluxo Principal:**
1. O Usuário navega pela árvore de pastas e segredos.
2. O Usuário expande, colapsa e move o foco entre os nós.
3. A aplicação exibe segredos antes de subpastas em cada coleção, conforme a ordem persistida.

**Pós-condições:**
- Nenhuma alteração no cofre. O foco determina as ações contextuais disponíveis.

---

### UC-13: Visualizar segredo

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; segredo disponível na hierarquia

**Fluxo Principal:**
1. O Usuário seleciona um segredo na hierarquia.
2. A aplicação exibe os detalhes: nome, campos (com sensíveis ocultos), observação.
3. O segredo fica disponível para ações: edição, favoritar, mover, copiar, excluir.

**Pós-condições:**
- Nenhuma alteração no segredo ou no cofre.

---

### UC-14: Visualizar/ocultar campo sensível

**Ator:** Usuário

**Pré-condições:**
- Segredo visível com campo do tipo texto sensível

**Fluxo Principal:**
1. O Usuário solicita a revelação de um campo sensível.
2. A aplicação revela temporariamente o valor do campo.
3. Após o tempo configurado (padrão: 15s), a aplicação reoculta automaticamente o valor.

**Fluxos Alternativos:**

*FA-01: Ocultação manual*
1. O Usuário reoculta o campo manualmente antes do tempo automático.

**Pós-condições:**
- Nenhuma alteração no segredo ou no cofre.

---

### UC-15: Buscar segredos

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo

**Fluxo Principal:**
1. O Usuário inicia a busca e informa o critério.
2. A aplicação varre os segredos em memória, comparando com nome, nome de campo, valores de campos tipo texto e observação.
3. A hierarquia é filtrada, exibindo apenas os segredos que satisfazem o critério, dentro de suas pastas para preservar contexto.
4. Quando o casamento ocorre no nome do segredo, o trecho pesquisado recebe destaque visual.
5. O Usuário seleciona um resultado (encerrando a busca) ou cancela a busca.
6. O cofre retorna ao estado anterior ao início da busca.

**Pós-condições:**
- Nenhuma alteração no cofre.

**Regras de Negócio:** RN-07, RN-08

---

### UC-16: Criar segredo

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo

**Fluxo Principal:**
1. O Usuário inicia a criação de um novo segredo na raiz ou na pasta ativa.
2. A aplicação oferece a escolha: usar modelo existente ou começar com segredo vazio.
3. Se modelo escolhido: a aplicação gera o segredo com campos copiados do modelo, sem vínculo.
4. Se segredo vazio: a aplicação gera um segredo sem campos adicionais.
5. O Usuário preenche os dados (nome, campos, observação) e confirma.
6. O segredo assume estado Novo e o cofre entra em estado Cofre Modificado.

**Fluxos Alternativos:**

*FA-01: Criação com modelo leva à edição padrão*
1. A aplicação inicia automaticamente o fluxo de edição padrão (UC-18).

*FA-02: Criação vazia leva à edição avançada*
1. A aplicação inicia automaticamente o fluxo de edição avançada (UC-19) para que o Usuário adicione os campos desejados.

*FA-03: Usuário cancela*
1. O segredo em criação é descartado sem efeito.

**Pós-condições:**
- Novo segredo inserido na posição adequada.

**Regras de Negócio:** RN-10

---

### UC-17: Duplicar segredo

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; segredo selecionado

**Fluxo Principal:**
1. O Usuário inicia a ação de duplicação.
2. A aplicação cria novo segredo com nova identidade, copiando nome (com sufixo incremental), campos, observação, favorito e nome do modelo do segredo original.
3. O segredo duplicado assume estado Novo e é inserido logo abaixo do original.
4. O cofre entra em estado Cofre Modificado.

**Pós-condições:**
- Novo segredo disponível na hierarquia.

**Regras de Negócio:** RN-26

---

### UC-18: Editar segredo (edição padrão)

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; segredo selecionado e disponível (não em Lixeira)

**Fluxo Principal:**
1. O Usuário inicia a edição padrão.
2. A aplicação permite alterar: nome, valores dos campos existentes e observação.
3. A identidade do segredo é preservada.
4. O Usuário confirma as alterações.
5. Se o segredo era Disponível, passa a Modificado. Se já era Novo ou Modificado, permanece no mesmo estado.
6. O cofre entra em estado Cofre Modificado.

**Fluxos Alternativos:**

*FA-01: Alternar para edição avançada*
1. O Usuário pode a qualquer momento alternar para a edição avançada (UC-19).

*FA-02: Usuário cancela*
1. O segredo retorna ao estado anterior sem alterações.

**Pós-condições:**
- O segredo foi atualizado com os novos valores.

**Regras de Negócio:** RN-25

---

### UC-19: Editar segredo (edição avançada)

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; segredo selecionado e disponível (não em Lixeira)

**Fluxo Principal:**
1. O Usuário inicia a edição avançada.
2. A aplicação permite alterar a estrutura: incluir campos, renomear campos, excluir campos e reordenar campos.
3. Não é permitido alterar o tipo de um campo existente.
4. O Usuário confirma as alterações.
5. Se o segredo era Disponível, passa a Modificado. Se já era Novo ou Modificado, permanece no mesmo estado.
6. O cofre entra em estado Cofre Modificado.

**Fluxos Alternativos:**

*FA-01: Alternar para edição padrão*
1. O Usuário pode a qualquer momento alternar para a edição padrão (UC-18).

*FA-02: Usuário cancela*
1. O segredo retorna ao estado anterior sem alterações.

**Pós-condições:**
- A estrutura do segredo foi atualizada.

**Regras de Negócio:** RN-11, RN-25

---

### UC-20: Favoritar/desfavoritar segredo

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; segredo selecionado e disponível

**Fluxo Principal:**
1. O Usuário alterna o marcador de favorito do segredo selecionado.
2. A aplicação altera apenas o atributo favorito, sem modificar identidade, conteúdo ou localização.
3. A presença na pasta virtual Favoritos reflete imediatamente a mudança.
4. O cofre entra em estado Cofre Modificado.

**Pós-condições:**
- O estado de favorito do segredo foi alternado.

---

### UC-21: Excluir segredo reversivelmente

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; segredo disponível na hierarquia principal

**Fluxo Principal:**
1. O Usuário inicia a exclusão do segredo.
2. A aplicação exige confirmação, informando que a exclusão é reversível até o próximo salvamento.
3. A aplicação retira o segredo da hierarquia principal e o materializa na Lixeira.
4. A aplicação memoriza a pasta de origem e o estado anterior do segredo.
5. O cofre entra em estado Cofre Modificado.

**Pós-condições:**
- O segredo está na Lixeira, restaurável até o próximo salvamento.

**Regras de Negócio:** RN-04, RN-23, RN-25

---

### UC-22: Restaurar segredo da Lixeira

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; segredo presente na Lixeira

**Fluxo Principal:**
1. O Usuário seleciona um segredo na Lixeira e inicia a restauração.
2. A aplicação reinserir o segredo na pasta de origem, ao final da lista de segredos, com identidade e conteúdo preservados.
3. O segredo retorna ao estado que possuía antes da exclusão reversível.
4. O cofre entra em estado Cofre Modificado.

**Fluxos Alternativos:**

*FA-01: Pasta de origem não existe mais*
1. O segredo é reinserido na raiz do cofre.
2. A aplicação informa que a pasta original não existe mais.

**Pós-condições:**
- O segredo voltou à hierarquia principal.

**Regras de Negócio:** RN-24

---

### UC-23: Mover segredo

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; segredo selecionado e disponível

**Fluxo Principal:**
1. O Usuário inicia a movimentação e informa o novo destino (outra pasta ou raiz).
2. O segredo é removido da coleção atual e reinserido no destino, ao final da lista de segredos.
3. Identidade, conteúdo e favorito são preservados.
4. O cofre entra em estado Cofre Modificado.

**Pós-condições:**
- O segredo está no novo destino.

**Regras de Negócio:** RN-02

---

### UC-24: Reordenar segredo

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; segredo selecionado e disponível

**Fluxo Principal:**
1. O Usuário inicia a reordenação e define a nova posição entre os segredos irmãos.
2. A aplicação altera a posição do segredo na coleção, preservando identidade e conteúdo.
3. O cofre entra em estado Cofre Modificado.

**Pós-condições:**
- O segredo está na nova posição.

---

### UC-25: Copiar campo de segredo

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; segredo visível

**Fluxo Principal:**
1. O Usuário seleciona um campo (inclusive texto sensível) e inicia a cópia.
2. A aplicação copia o valor para a área de transferência do sistema.
3. A aplicação exibe confirmação e inicia o temporizador de limpeza automática (padrão: 30s).

**Pós-condições:**
- O valor está na área de transferência. Será limpo ao expirar o temporizador ou ao bloquear/fechar o cofre.

---

### UC-26: Criar pasta

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo

**Fluxo Principal:**
1. O Usuário inicia a criação de uma nova pasta.
2. A aplicação determina o destino conforme o contexto: se o foco está em uma pasta, a nova pasta é criada dentro dela; se na raiz, é criada na raiz.
3. O Usuário informa o nome da pasta.
4. A nova pasta é adicionada ao final da lista de subpastas do destino.
5. O cofre entra em estado Cofre Modificado.

**Pós-condições:**
- Nova pasta inserida na hierarquia.

**Regras de Negócio:** RN-03

---

### UC-27: Renomear pasta

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; pasta selecionada

**Fluxo Principal:**
1. O Usuário inicia a renomeação.
2. O Usuário informa o novo nome.
3. A aplicação altera apenas o nome, sem modificar identidade, posição ou conteúdo.
4. O cofre entra em estado Cofre Modificado.

**Pós-condições:**
- A pasta foi renomeada.

---

### UC-28: Mover pasta

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; pasta selecionada

**Fluxo Principal:**
1. O Usuário inicia a movimentação e informa o novo destino (outra pasta ou raiz).
2. A pasta é removida da coleção atual e reinserida no destino, preservando identidade, conteúdo e hierarquia interna.
3. Todos os filhos (recursivamente) acompanham a pasta.
4. O cofre entra em estado Cofre Modificado.

**Pós-condições:**
- A pasta está no novo destino com toda sua sub-hierarquia.

**Regras de Negócio:** RN-03

---

### UC-29: Reordenar pasta

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; pasta selecionada

**Fluxo Principal:**
1. O Usuário inicia a reordenação e define a nova posição entre as pastas irmãs.
2. A aplicação altera a posição da pasta na coleção, preservando identidade e hierarquia interna.
3. O cofre entra em estado Cofre Modificado.

**Pós-condições:**
- A pasta está na nova posição.

---

### UC-30: Excluir pasta

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; pasta selecionada

**Fluxo Principal:**
1. O Usuário inicia a exclusão da pasta.
2. A aplicação exige confirmação, informando que a exclusão é imediata e irreversível.
3. A pasta é removida da hierarquia.
4. Segredos e subpastas filhas são promovidos para o nível pai (ou raiz), adicionados ao final das listas correspondentes.
5. Todos os filhos (recursivamente) mantêm identidade e posição relativa entre si.
6. O cofre entra em estado Cofre Modificado.

**Pós-condições:**
- A pasta não existe mais. Seus filhos estão no nível pai.

**Regras de Negócio:** RN-06

---

### UC-31: Criar modelo de segredo

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo

**Fluxo Principal:**
1. O Usuário inicia a criação de um novo modelo.
2. O Usuário informa o nome do modelo e define a estrutura de campos (cada um com nome e tipo).
3. O modelo recebe identidade própria e fica disponível para criações futuras.
4. O cofre entra em estado Cofre Modificado.

**Pós-condições:**
- Novo modelo disponível para uso.

---

### UC-32: Editar modelo de segredo

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; modelo selecionado

**Fluxo Principal:**
1. O Usuário inicia a edição do modelo.
2. A aplicação permite: alterar nome, incluir campos, alterar nome ou tipo de campos existentes, excluir campos e reordenar campos.
3. O Usuário confirma as alterações.
4. O cofre entra em estado Cofre Modificado.

**Pós-condições:**
- O modelo foi atualizado. As alterações afetam apenas criações futuras de segredos.

**Regras de Negócio:** RN-12

---

### UC-33: Remover modelo de segredo

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; modelo selecionado

**Fluxo Principal:**
1. O Usuário inicia a remoção do modelo.
2. A aplicação exige confirmação.
3. O modelo é removido permanentemente.
4. O cofre entra em estado Cofre Modificado.

**Pós-condições:**
- O modelo não está mais disponível para uso. Segredos já criados a partir dele não são afetados.

**Regras de Negócio:** RN-10

---

### UC-34: Criar modelo a partir de segredo existente

**Ator:** Usuário

**Pré-condições:**
- Cofre ativo; segredo selecionado

**Fluxo Principal:**
1. O Usuário inicia a criação de modelo a partir do segredo selecionado.
2. A aplicação copia a estrutura de campos do segredo (nome e tipo de cada campo) como base inicial.
3. O Usuário informa o nome do novo modelo e confirma.
4. O modelo recebe identidade própria, sem vínculo retroativo com o segredo de origem.
5. O cofre entra em estado Cofre Modificado.

**Pós-condições:**
- Novo modelo disponível, baseado na estrutura do segredo.

**Regras de Negócio:** RN-10
