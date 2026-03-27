


## Estados e Fluxos Principais

### Estados principais

#### Estados globais da aplicação

- **Sem cofre:** estado inicial, a aplicação está em execução, mas ainda não há um cofre ativo.
- **Apresentando cofre:** estado global em que existe um cofre carregado em memória, autenticado e disponível para uso na sessão atual. 
  - **Cofre em pesquisa:** estado transitório sobreposto ao `Apresentando cofre`, no qual a interface exibe uma busca ativa, mostrando apenas os segredos que correspondem aos critérios informados. Durante a busca, o cofre preserva o estado. Enquanto a pesquisa estiver ativa, todas as ações ficam indisponíveis exceto: sair da aplicação, navegar pelo cofre e visualizar segredo. Para retomar as demais ações, o usuário deve confirmar a pesquisa (selecionando o elemento desejado, o que encerra implicitamente a pesquisa) ou cancelar a pesquisa. Em ambos os casos, o cofre retorna ao estado anterior ao início da pesquisa.
- **Tamanho insuficiente do terminal:** estado transitório em que a aplicação detecta que o terminal é pequeno demais para exibir a interface, mostrando uma mensagem de aviso e solicitando que o usuário aumente o tamanho do terminal para continuar usando a aplicação. Uma reestabelecido o tamanho mínimo, a aplicação retorna ao estado anterior.

Temporariamente, durante os estados anteriores, a aplicação pode assumir um estado transitório, retornando ao último estado válido:
  - **Aviso modal:** estado transitório sobreposto ao estado principal, para mostrar uma mensagem crítica que precisa ser confirmada
  - **Confirmação modal:** estado transitório sobreposto ao estado principal, usado para confirmações críticas e de ações destrutivas.
  - **File picker:** estado transitório para seleção de caminho e nome de arquivo, usado nos fluxos de abertura, criação e salvamento com outro caminho do cofre. Também é usado para seleção de caminho de importação e exportação do cofre em formato JSON plain text.
  - **Tela de ajuda e comandos:** estado transitório sobreposto ao estado principal, usado para apresentar informações de ajuda e comandos disponíveis.
    
#### Estado do cofre

**Quanto ao ciclo de vida**
- **Cofre Salvo:** cofre sincronizado com o arquivo corrente.
- **Cofre Modificado:** cofre com divergência entre o estado em memória e o último estado salvo.

OBS: 
- O estado do cofre `Cofre Modificado` deve refletir qualquer divergência entre memória e último salvamento persistido.
- não existe estado observável "bloqueado", pois o bloqueio é tratado como abrir novamente o cofre, exigindo nova autenticação e recarregando o estado salvo do arquivo.
- não existe cofre "novo", pois o cofre é criado com a estrutura inicial e salvo imediatamente, entrando diretamente no estado "Cofre Salvo" desde o início.

Invariante:
- Só pode existir um cofre ativo por vez.

#### Estados do segredo, pastas e modelos de segredo

**Quanto a navegação nas telas**
- **disponível:** estado padrão.
- **ativo:** segredo atualmente com foco, passível de ações do usuário.

Invariante:
- Só pode existir um segredo ativo por vez. Ou nenhum segredo ativo.
- Se um segredo estiver ativo, então a pasta que o contém também é considerada implicitamente ativa.

**Quanto ao cliclo de vida**
- **original:** segredo carregado do arquivo, sem alterações.
- **em criação:** segredo durante o fluxo de criação; pode ser cancelado (descartado) sem efeito persistente.
- **novo**: segredo criado, ainda não salvo. Ele poderá ser novamente editado em modo padrão ou avançado e continuará sendo considerado novo até o próximo salvamento.
- **em edição padrão:** segredo durante o fluxo de edição; pode ser cancelado (revertido) sem efeito persistente.
- **em edição avançada:** segredo durante o fluxo de edição avançada; pode ser cancelado (revertido) sem efeito persistente.
- **modificado:** segredo carregado do arquivo e que sofreu alterações, ainda não salvo. Ele poderá ser novamente editado em modo padrão ou avançado e continuará sendo considerado modificado até o próximo salvamento.
- **excluído reversivelmente:** segredo retirado da hierarquia principal e materializado apenas na Lixeira até o próximo salvamento. A aplicação memoriza a pasta de origem e o estado anterior do segredo para possibilitar restauração ao local e estado originais. Enquanto permanecer nesse estado, não pode ser editado.

OBS:
- Não temos um estado **restaurado:** pois, ao restaurar, o segredo retoma o estado que possuía antes da exclusão reversível.

Invariante:
- A solução deve garantir que o estado de ciclo de vida dos segredos seja consistente.
- Quando ocorrer alteração de dados de um segredo (nome, valores de campos, estrutura dos campos, observações), então o estado de ciclo de vida deve ser ajusta corretamente:
  - Se o segredo estava em estado `original`, ele passa para `modificado`.
  - Se o segredo estava em estado `novo`, ele permanece em `novo`.
  - Se o segredo estava em estado `modificado`, ele permanece em `modificado`.
  - Se o segredo estava em estado `em criação`, ele permanece em `em criação`.
  - Se o segredo estava em estado `em edição padrão` ou `em edição avançada`, ele permanece no respectivo estado de edição.
  - Se o segredo estava em estado `excluído reversivelmente`, ele não pode ser editado, portanto não pode ocorrer alteração de dados.
- O estado do segredo não pode mudar espontaneamente sem uma ação do usuário que cause a mudança. Ele deve permanecer estável até que o usuário execute uma ação que o altere.
- FAvoritar/desfavoritar um segredo não altera seu estado de ciclo de vida.
- Quando o segredo entrar em edição, a aplicação proverá um mecanismo para "reverter" a as edições caso o usuário deseje cancelar a edição.
  Para tal, a aplicação poderia criar uma copia do segredo no estado em que ele estava antes da edição. Porém, não vamos estabelecer uma solução específica de implementação.
- O segredo que estiver em criação poderá ser cancelado. NEste caso, ele não será adicionado ao cofre.

#### Estados da pasta

- **Pasta existente:** pasta presente na hierarquia, passível de renomeação, movimentação e exclusão física com promoção dos filhos.
- **Pasta ativa:** pasta atualmente selecionada para ações de edição, movimentação e exclusão. Se um segredo estiver ativo, então a pasta que o contém também é considerada implicitamente ativa.

OBS:
- Pastas não possuem soft delete; sua exclusão sempre remove a pasta e promove os filhos.
- Pastas não possuem estado de modificado ou novo, pois não há necessidade de feedback visual específico para alterações em modelos.
- Pastas não possuem estado de edição ou criação, pois sua edição é feita diretamente na hierarquia e tem efeito imediato, sem fluxo separado.

#### Estado do modelo de segredo
- **Modelo disponível:** modelo existente e disponível para criação de novos segredos.
- **Modelo ativo:** modelo atualmente selecionado para ações de edição e exclusão. Se um segredo criado a partir deste modelo estiver ativo, então o modelo também é considerado implicitamente ativo.
- **Modelo em criação:** modelo ainda não confirmado pelo fluxo de criação; pode ser cancelado (descartado) sem efeito persistente.
- **Modelo em edição:** modelo com alteração estrutural em andamento, afetando apenas criações futuras após confirmação.

OBS:
- Modelos de segredo não possuem soft delete; sua exclusão sempre remove o modelo.
- Modelos de segredo não possuem estado de modificado ou novo, pois não há necessidade de feedback visual específico para alterações em modelos.

#### Estados transitórios de exposição de dados sensíveis

- **Campo sensível oculto:** estado padrão de exibição para campos do tipo `texto sensível`.
- **Campo sensível exibido temporariamente:** estado temporário após ação explícita do usuário, encerrado manualmente ou por temporizador configurado.
- **Campo sensível na área de transferência temporariamente:** existe um valor copiado aguardando limpeza automática por temporizador ou por bloqueio/fechamento do cofre.

### Fluxos iniciais

**Tamanho do temrinal**
  - Quando o tamanho do terminal for reduzido a um tamanho menor que o mínimo necessário para exibir a interface
  - A aplicação exibe uma mensagem de aviso solicitando que o usuário aumente o tamanho do terminal para continuar usando a aplicação.
  - Retorna ao estado anterior assim que o tamanho do terminal for restabelecido para o mínimo necessário. 

**Abrir aplicação**
  - Ao iniciar, a aplicação mostra uma tela de welcome com ASCII art de apresentação do Abditum. Aplicação entra em estado global `Sem cofre`.
  - A tela inicial oferece as ações de criar cofre, abrir cofre, acessar ajuda e sair.
  - A partir dessa tela, a aplicação permanece no estado `Inicial / sem cofre ativo` até o usuário escolher a próxima ação.

**Criar novo cofre**
  - Aplicação está no estado global `Sem cofre`. 
  - Usuário informa caminho e senha mestra com confirmação.
  - A aplicação popula a estrutura inicial do cofre com modelos e pastas padrão.
  - Se não existir arquivo no caminho informado, a aplicação grava diretamente o novo cofre no caminho final, usando o formato da versão atual.
  - Se já existir arquivo no caminho informado:
    - A aplicação exige confirmação explícita de sobrescrita.
    - Se já existir um backup anterior com extensão `.abditum.bak`, a aplicação o renomeia temporariamente para `.abditum.bak2` antes de gerar o novo backup.
    - A aplicação gera um novo backup do arquivo existente com extensão `.abditum.bak` e então grava diretamente o novo cofre no caminho final.
    - Se a operação for concluída com sucesso, a aplicação remove o `.abditum.bak2`, preservando apenas o novo `.abditum.bak`.
    - Se a operação falhar antes da consolidação final, a aplicação restaura o `.abditum.bak2` para `.abditum.bak` sempre que possível.
    - Em caso de falha na gravação do novo arquivo após a geração do backup, a aplicação deve exibir uma mensagem de erro informando a falha e que existe um backup disponível para intervenção manual do usuário.
  - Esse fluxo não utiliza arquivo `.abditum.tmp`, pois não se trata do salvamento incremental de um cofre já aberto, e sim da criação de um novo arquivo de cofre.
  - O cofre entra em estado `Cofre Salvo`. Aplicação entra em estado global `Apresentando cofre`.

**Abrir cofre existente**
  - Aplicação está no estado global `Sem cofre`. 
  - Usuário informa caminho.
  - A aplicação valida assinatura `magic` e `versão_formato`.
  - Seleciona o perfil Argon2id histórico a partir de `versão_formato`.
  - Usuário informa senha mestra.
  - Deriva a chave, valida o payload cifrado e carrega o domínio em memória.
  - Se o payload descriptografado estiver em um formato histórico suportado, a aplicação realiza a migração dos dados em memória para o modelo corrente do domínio.
  - O cofre entra em estado `Cofre Salvo`. Aplicação entra em estado global `Apresentando cofre`.
  
**Sair da aplicação**
  - O usuário pode encerrar a aplicação a qualquer momento.
    - No estado `Sem cofre`, a aplicação encerra após solicitar confirmação do encerramento.
    - No estado `Mostrando cofre` e `Cofre Salvo`, a aplicação encerra após solicitar confirmação do encerramento.
    - No estado `Mostrando cofre` e `Cofre Modificado`, a aplicação oferece as opções de salvar, descartar alterações ou cancelar o encerramento, para evitar perda acidental de dados.
      - Em caso de salvar, a aplicação segue o fluxo de salvamento descrito anteriormente e encerra somente após salvamento bem-sucedido.
    OBS:
    - Também é possível sair quando houver fluxos em andamento (ex: criação ou edição de segredo), não havendo um aviso específico para este caso.

#### Fluxos do cofre

Pressupõe-se que a aplicação já esteja em estado global `Apresentando cofre`.

**Visualizar hierarquia do cofre**
  - O usuário navega pela árvore de pastas e segredos do cofre.
  - A aplicação apresenta a hierarquia conforme a ordem persistida no JSON, mostrando primeiro subpastas e depois segredos em cada coleção.
  - Ao focar um segredo, a aplicação torna o segredo ativo (o que implica a exibição dos detalhes do segredo, incluindo os campos e a observação, com os dados sensíveis ocultos por padrão).
  - O usuário pode expandir, colapsar e mover o foco entre os nós.
  - Enquanto o usuário navega, enquanto não focar outro segredo, o segredo ativo permanece o mesmo.
  - Esse fluxo não não altera o estado do cofre, nem das pastas, nem dos segredos, nem dos modelos de segredo. Ele é apenas de navegação e visualização, sem efeitos colaterais.

**Bloquear acesso ao cofre**
  - O bloqueio pode ser manual ou por inatividade.
  - A aplicação fecha logicamente o cofre, limpa buffers controlados sempre que possível e limpa a área de transferência.
  - A aplicação volta estado "Sem cofre" no fluxo "Abrir cofre existente", assumindo o mesmo caminho do cofre previamente aberto, mas exigindo nova autenticação para desbloquear.

OBS:
- Se o cofre estiver em estado `Cofre Modificado`, as alterações não salvas são descartadas silenciosamente, sem confirmação. Essa é uma decisão de projeto: o bloqueio por inatividade ocorre em sessão desassistida, e o bloqueio manual emergencial (proteção contra shoulder surfing) precisa ser imediato — em ambos os casos, confirmações comprometeriam o propósito do bloqueio.

**Salvar cofre**
  - O cofre ativo em estado `Cofre Modificado`.
  - A aplicação grava o cofre num caminho com sufixo ".abditum.tmp", usando o formato da versão atual, e atualiza a `versão_formato` do cabeçalho quando necessário, com `nonce` diferente.
  - Se já existir um backup anterior com extensão `.abditum.bak`, a aplicação o renomeia temporariamente para `.abditum.bak2` antes de gerar o novo backup.
  - Copia o arquivo atual do cofre para um novo backup com extensão `.abditum.bak`.
  - Depois renomeia o arquivo `.abditum.tmp` para o nome final do cofre, substituindo o arquivo original.
  - Se a operação for concluída com sucesso, a aplicação remove o `.abditum.bak2`, preservando apenas o novo `.abditum.bak`.
  - Se a operação falhar antes da consolidação final, a aplicação restaura o `.abditum.bak2` para `.abditum.bak` sempre que possível.
  - Em caso de falha na escrita ou substituição do arquivo final após a geração do backup, a aplicação deve exibir uma mensagem de erro informando a falha e que existe um backup disponível para intervenção manual do usuário.
  - Se a persistência for bem-sucedida, o cofre entra em estado `Cofre Salvo`. Aplicação permanece no global `Apresentando cofre`.

**Descartar alterações não salvas e recarregar cofre**
  - O cofre ativo em estado `Cofre Modificado`.
  - O usuário inicia a ação de descartar alterações e recarregar o cofre ativo.
  - A aplicação exige confirmação para descartar as alterações locais ainda não persistidas.
  - Após a confirmação, a aplicação reabre o arquivo atual, reusando a senha previamente fornecida, repetindo validação, descriptografia e eventual migração em memória.
  - Ao final, o cofre retorna ao estado `Cofre Salvo`.

**Salvar cofre em novo caminho** 
  - O usuário inicia a ação de salvar o cofre em um novo caminho.
  - A aplicação entra em estado transitório `Salvando cofre com outro caminho`.
  - A aplicação coleta o novo caminho de destino e valida a possibilidade de gravação nesse local.
  - Se não houver arquivo no destino, a aplicação grava o cofre diretamente no novo caminho, usando o formato da versão atual.
  - Se já houver arquivo no destino, a aplicação exige confirmação de sobrescrita.
  - Se já existir um backup anterior com extensão `.abditum.bak` no destino, a aplicação o renomeia temporariamente para `.abditum.bak2` antes de gerar o novo backup.
  - A aplicação gera um novo backup do arquivo existente com extensão `.abditum.bak` e então grava o cofre diretamente no caminho final.
  - Se a operação for concluída com sucesso, a aplicação remove o `.abditum.bak2`, preservando apenas o novo `.abditum.bak`.
  - Se a operação falhar antes da consolidação final, a aplicação restaura o `.abditum.bak2` para `.abditum.bak` sempre que possível.
  - Em caso de falha na gravação do novo arquivo após a geração do backup, a aplicação deve exibir uma mensagem de erro informando a falha e que existe um backup disponível para intervenção manual do usuário.
  - Esse fluxo não utiliza arquivo `.abditum.tmp`, pois não se trata do salvamento incremental sobre o caminho atual do cofre já aberto.
  - Após o salvamento bem-sucedido, o novo caminho passa a ser o caminho atual do cofre e o estado retorna para `Cofre Salvo`.
  - O cofre entra em estado `Cofre Salvo`. Aplicação entra em estado global `Apresentando cofre`.

**Alterar senha mestra do cofre**
  - O usuário inicia a ação de alteração da senha mestra sobre o cofre ativo.
  - A aplicação solicita a nova senha mestra e a confirmação da nova senha.
  - Se a confirmação da nova senha for válida, a aplicação rederiva a chave com um novo `salt` e prepara o cofre para ser persistido com a nova credencial.
  - A alteração da senha mestra não modifica o conteúdo lógico do domínio, mas exige regravação criptográfica completa do arquivo com novo `salt`, novo `nonce` e a chave derivada da nova senha mestra.
  - A partir deste ponto, a aplicação segue o fluxo de **Salvar cofre**, incluindo gravação atômica, rotação de backup e tratamento de falha.

**Configurar o cofre**
  - O usuário inicia a edição das configurações do cofre ativo.
  - A aplicação permite alterar o tempo de bloqueio automático por inatividade, o tempo de reocultação de campos sensíveis e o tempo de limpeza automática da área de transferência.
  - As alterações passam a valer para o comportamento da sessão corrente conforme aplicável e permanecem associadas ao próprio cofre.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Exportar cofre para JSON plain text**
  - Antes de exportar, a aplicação mostra aviso explícito sobre o risco de segurança de gerar uma cópia não criptografada e exige confirmação.
  - O usuário inicia a exportação do cofre ativo para formato JSON plain text.
  - A exportação serializa o estado atual do domínio em memória, incluindo eventuais alterações não salvas.
  - Se o cofre estiver em estado `Cofre Modificado`, a aplicação exibe alerta informando que a exportação incluirá alterações ainda não salvas.
  - Após a confirmação, a aplicação serializa o domínio para JSON em texto claro no destino escolhido pelo usuário.
  - Esse fluxo não altera o conteúdo lógico do cofre ativo nem seu estado persistido.

**Importar cofre de JSON plain text**
  - O usuário inicia a importação de um arquivo JSON plain text para o cofre ativo.
  - A aplicação lê o conteúdo importado e resolve conflitos por identidade conforme as regras do cofre.
  - Pastas com a mesma identidade são mescladas silenciosamente.
  - Se um segredo importado colidir por identidade com um segredo já existente, a aplicação cria um novo segredo logicamente equivalente, com identidade diferente e preservando os demais dados importados.
  - Se um segredo importado colidir por nome com outro segredo já existente na mesma pasta de destino, a aplicação ajusta seu nome com sufixo numérico incremental e informa esse ajuste ao usuário.
  - Modelos com a mesma identidade são sobrepostos silenciosamente pelo modelo importado.
  - Após a confirmação e incorporação dos dados, o cofre:
    - se o estado original do cofre era `Cofre Modificado`, permanece em `Cofre Modificado` independentemente de ter havido alterações efetivas ou não, pois o estado modificado já reflete a existência de divergências entre memória e último salvamento.
    - senão entra em estado `Cofre Modificado` caso a importação tenha resultado em alguma alteração, 
    - ou permanece em `Cofre Salvo` caso a importação não tenha introduzido nenhuma alteração efetiva.

**Visualizar segredo**
  - Um segredo torna-se ativo, seja por navegação, seja por busca ou por visualização direta.
  - A aplicação exibe os detalhes do segredo, incluindo nome, observação e campos, com os dados sensíveis ocultos por padrão.

**Criar segredo**
  - O usuário solicita a criação de um novo segredo, seja a partir da raiz do cofre, seja a partir de uma pasta específica.
  - A aplicação oferece a escolha entre usar um modelo de segredo existente ou começar com um segredo vazio, sem nenhuma estrutura inicial.
    - Caso o usuário opte por um modelo de segredo, a estrutura inicial do segredo é gerada a partir do modelo escolhido, copiando os campos como snapshot, sem manter vínculo por referência com o modelo de origem.
    - Caso o usuário opte por começar com um segredo vazio, é gerado um segredo sem campos adicionais além do nome e da observação, e os demais campos poderão ser adicionados posteriormente pela edição avançada.
  - Após a confirmação, o novo segredo assume estado `Novo` e é inserido no destino selecionado, e o cofre entra em estado `Cofre Modificado`.
    - Caso o usuário tenha optado por um modelo de segredo, a aplicação passa para o fluxo de edição padrão.
    - Caso o usuário tenha optado por um segredo vazio, a aplicação passa para o fluxo de edição avançada, para que o usuário possa adicionar os campos desejados.

#### Fluxos principais de segredos

Pressupõe-se que:
  - A aplicação está em estado global `Apresentando cofre` com o cofre ativo em estado `Cofre Salvo` ou `Cofre Modificado`.
  - Existe um segredo ativo, seja por navegação, seja por busca ou por visualização direta.

**Visualizar ou ocultar campo sensível**
  - O usuário visualiza um segredo ativo.
  - O usuário seleciona um campo do tipo `texto sensível`.
  - O usuário solicita a exibição temporária do valor do campo sensível.
  - A aplicação revela temporariamente o valor do campo.
  - O usuário solicita a exibição temporária do valor do campo sensível. Ou ocorre o encerramento automático da exibição temporária por expiração do tempo configurado no cofre.
  - A aplicação oculta o valor do campo.
  - Esse fluxo não altera o estado do segredo nem o estado do cofre.

**Copiar campo de segredo**
  - O usuário visualiza um segredo ativo.
  - O usuário seleciona um campo do tipo `texto sensível`.
  - O usuário solicita a copia temporária do valor do campo sensível para a área de transferência do sistema.
  - A aplicação copia o valor atual do campo para a área de transferência do sistema.
  - A aplicação exibe feedback visual de cópia e inicia o temporizador de limpeza automática conforme a configuração do cofre.
  - O conteúdo copiado também é limpo ao bloquear ou fechar o cofre.
  - Esse fluxo não altera o estado do segredo nem o estado do cofre.

**Duplicar segredo**
  - O usuário solicita a duplicação.
  - A aplicação cria uma nova instância com nova identidade, copiando nome, nome do modelo de segredo, observação, favorito e campos do segredo original.
    - O nome do segredo duplicado recebe um sufixo numérico incremental para evitar confusão com o segredo original. Ex: "Segredo" → "Segredo (1)", "Segredo (2)", etc.
  - O segredo duplicado assume estado `Novo`.
  - Ele é inserido logo abaixo do segredo de origem na mesma coleção
  - O cofre entra em estado `Cofre Modificado`.

**Favoritar segredo**
  - O segredo ativo não está favoritado. 
  - O usuário solicita favoritar segredo.
  - O segredo altera seu status de `favorito` para true, sem modificar identidade, conteúdo ou localização do segredo.
  - A pasta virtual de Favoritos atualiza para refletir imediatamente esse estado: A posição do segredo na pasta favorita será sua posição percorrendo a árvore em profundidade, conforme as listas de segredos e pastas em cada pasta, seguindo a ordem persistida no JSON.
  - O cofre entra em estado `Cofre Modificado`.

**Desfavoritar segredo**
  - O segredo ativo está favoritado. 
  - O usuário solicita desfavoritar segredo.
  - O segredo altera seu status de `favorito` para false, sem modificar identidade, conteúdo ou localização do segredo.
  - A pasta virtual de Favoritos atualiza para refletir imediatamente esse estado: A posição do segredo na pasta favorita será sua posição percorrendo a árvore em profundidade, conforme as listas de segredos e pastas em cada pasta, seguindo a ordem persistida no JSON.
  - O cofre entra em estado `Cofre Modificado`.

**Editar segredo (edição padrão)**
  - O usuário visualiza o segredo ou está editando via edição avançada, solicita a edição padrão.
  - O segredo entra no estado de edição padrão.
  - A aplicação provê campos para permitir alterar nome, observação e valores dos campos existentes.
  - O usuário realiza as alterações desejadas através da edição dos campos.
  - Ao focar sobre um campo do tipo `texto sensível`, a aplicação mostra o valor atual do campo para permitir a edição, ocultando novamente o valor ao sair do foco do campo.
  - O segredo preserva sua identidade durante toda a edição.
  - O usuário pode alternar para a edição avançada caso precise alterar a estrutura do segredo.
  - Durante a edição padrão, o usuário poderá solicitar a edição avançada. As alterações de valores realizadas até o momento serão preservadas, e ficarão disponíveis para edição avançada. 
  Após a confirmação, o segredo preserva seu estado anterior se já estiver em `Segredo novo` ou `Segredo modificado`.
  - Após a confirmação, se o segredo estava em `Segredo disponível`, ele passa para `Segredo modificado`.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.
  - O segredo volta a visualização normal, mostrando os campos com os dados sensíveis ocultos por padrão.

**Editar segredo (edição avançada)**
  - O usuário visualiza o segredo ou está editando via edição padrão, solicita a edição avançada.
  - O segredo entra no modo de edição avançada.
  - Nesse modo, o usuário altera apenas a estrutura do segredo.
  - Não é permitido alterar o tipo de um campo existente. Para isso, é necessário excluir o campo e criar um novo com o tipo desejado.
  - O segredo preserva sua identidade durante toda a edição.
  - O usuário pode alternar para a edição padrão quando quiser voltar a alterar os valores dos campos.
  - Durante a edição avançada, o usuário poderá solicitar a edição padrão. As alterações de estrutura realizadas até o momento serão preservadas, e ficarão disponíveis para edição padrão.
  - Após a confirmação, o segredo preserva seu estado anterior se já estiver em `Segredo novo` ou `Segredo modificado`.
  - Após a confirmação, se o segredo estava em `Segredo disponível`, ele passa para `Segredo modificado`.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.
  - O segredo volta a visualização normal, mostrando os campos com os dados sensíveis ocultos por padrão.

**Excluir segredo reversivelmente**
  - O usuário solicita a remoção.
  - A aplicação exige confirmação (apesar da remoção ser reversível até o próximo salvamento do cofre).
  - O segredo preserva sua identidade e seu conteúdo durante a remoção.
  - A aplicação memoriza a pasta de origem do segredo (ou a raiz do cofre, caso o segredo estivesse na raiz) e o estado do segredo antes da exclusão, para permitir restauração ao local e estado originais.
  - Enquanto o segredo permanecer na Lixeira, a aplicação não permite edição desse segredo.
  - O cofre entra em estado `Cofre Modificado`, a aplicação retira o segredo da hierarquia principal e o materializa na pasta virtual Lixeira.

**Restaurar segredo excluído reversivelmente**
  - O segredo ativo é um segredo removido reversivelmente, presente na pasta virtual Lixeira.
  - O usuário solicita a restauração.
  - O segredo preserva sua identidade e seu conteúdo durante a restauração.
  - Se a pasta de origem ainda existir na hierarquia, o segredo é reinserido nessa pasta, ao final da lista de segredos.
  - Se a pasta de origem tiver sido excluída após o soft delete, o segredo é reinserido na raiz do cofre, ao final da lista de segredos, e a aplicação exibe uma mensagem informando que a pasta original não existe mais.
  - Após a restauração, o segredo retorna ao estado que possuía antes da exclusão reversível.
  - O cofre entra em estado `Cofre Modificado`.

**Mover segredo**
  - O usuário solicita a movimentação para outra pasta ou para a raiz do cofre.
  - A aplicação coleta o novo destino, que pode ser outra pasta ou a raiz do cofre.
  - O segredo é removido da coleção atual e reinserido na coleção de destino, preservando identidade, conteúdo e marcação de favorito.
  - A identidade e o conteúdo do segredo são preservados durante toda a movimentação. O estado do segredo permanece inalterado, mas sua posição na hierarquia é atualizada para refletir o novo destino.
  - O segredo é adicionado ao final da lista de segredos do destino.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Reordenar segredo**
  - O usuário solicita a reordenação relativa.
  - A aplicação altera apenas sua posição entre os segredos irmãos da mesma coleção pai.
  - A nova ordem passa a refletir diretamente a ordem persistida e a ordem de exibição.
  - A identidade e o conteúdo do segredo são preservados durante toda a movimentação. O estado do segredo permanece inalterado, mas sua posição na hierarquia é atualizada para refletir o novo destino.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Buscar segredos**
  - O usuário inicia o processo de busca.
  - A aplicação executa uma varredura sequencial em memória sobre nome do segredo, nome de campo, valores de campos do tipo `texto` e observação.
  - A hierarquia é reapresentada mostrando apenas os segredos que satisfazem o critério de busca, mas mantendo a estrutura de pastas para preservar o contexto de localização dos segredos encontrados.
  - Enquanto a pesquisa estiver ativa, todas as ações ficam indisponíveis exceto: sair da aplicação, navegar pelo cofre e visualizar segredo.
  - Quando o casamento ocorrer no nome do segredo, o segredo correspondente recebe destaque visual na árvore.
  - O usuário confirma a pesquisa selecionando o elemento desejado (o que encerra implicitamente a pesquisa) ou cancela a pesquisa. Em ambos os casos, o cofre retorna ao estado anterior ao início da pesquisa.

#### Fluxos principais de pastas

**Criar pasta**
  - O usuário inicia a criação de uma nova pasta na raiz do cofre ou na pasta ativa.
  - A aplicação coleta o nome da pasta e determina o destino conforme o contexto atual.
  - Se o destino for a raiz, a nova pasta é adicionada ao final da lista de pastas da raiz.
  - Se o destino for a pasta ativa, a nova pasta é adicionada ao final da lista de subpastas dessa pasta.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Renomear pasta**
  - O usuário seleciona uma pasta existente e inicia a ação de renomeação.
  - A aplicação coleta o novo nome e altera apenas esse atributo, sem modificar a identidade, a posição ou o conteúdo da pasta.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Mover pasta**
  - O usuário seleciona uma pasta existente e inicia a ação de movimentação.
  - A aplicação coleta o novo destino, que pode ser outra pasta ou a raiz do cofre.
  - A pasta é removida da coleção atual e reinserida na coleção de destino, preservando sua identidade, seu conteúdo e sua hierarquia interna.
    - Todos os filhos (recursivamente) da pasta movida (segredos e subpastas) acompanham a pasta em seu movimento, sem alteração de identidade ou posição relativa entre eles.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Reordenar pasta**
  - O usuário seleciona uma pasta existente e inicia a ação de reordenação relativa.
  - A aplicação altera apenas a posição da pasta entre as demais pastas irmãs da mesma coleção pai.
  - A nova ordem passa a refletir diretamente a ordem persistida e a ordem de exibição.
    - Todos os filhos (recursivamente) da pasta movida (segredos e subpastas) acompanham a pasta em seu movimento, sem alteração de identidade ou posição relativa entre eles.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Excluir pasta**
  - O usuário seleciona uma pasta existente e inicia a ação de exclusão.
  - A aplicação exige confirmação, pois a exclusão da pasta é física e imediata.
  - A pasta removida deixa de existir na hierarquia.
  - Seus segredos e subpastas filhas são promovidos para a pasta pai ou para a raiz, sendo adicionados ao final das listas correspondentes.
    - Todos os filhos (recursivamente) da pasta promovida (segredos e subpastas) mantêm sua identidade e posição relativa entre eles.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

#### Fluxos principais de modelos de segredo

**Criar modelo de segredo**
  - O usuário inicia a criação de um novo modelo de segredo.
  - A aplicação coleta o nome do modelo e sua estrutura inicial de campos, cada um com nome e tipo.
  - O novo modelo recebe identidade própria e passa a ficar disponível para criações futuras de segredos.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Editar modelo de segredo**
  - O usuário seleciona um modelo existente e inicia sua edição.
  - A aplicação permite alterar o nome do modelo, incluir novos campos, alterar nome ou tipo de campos existentes, excluir campos e reordenar os campos do modelo.
  - As alterações afetam apenas criações futuras de segredos e não alteram segredos já existentes que tenham sido criados a partir desse modelo.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Remover modelo de segredo**
  - O usuário seleciona um modelo existente e inicia a ação de remoção.
  - A aplicação exige confirmação antes de excluir o modelo.
  - A remoção do modelo impede apenas seu uso futuro e não afeta segredos já criados anteriormente.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Criar modelo a partir de segredo existente**
  - O usuário seleciona um segredo existente e inicia a criação de um modelo a partir dele.
  - A aplicação copia a estrutura de campos do segredo selecionado, preservando nome e tipo de cada campo como base inicial do novo modelo.
  - O usuário informa o nome do novo modelo e confirma sua criação.
  - O modelo resultante passa a ficar disponível para novas criações de segredos, sem criar vínculo retroativo com o segredo de origem.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.


## Estratégia de testes
- Testes: 
 - testes do serviço de criptografia, incluindo casos de sucesso e falha para criptografia e descriptografia;
 - testes do serviço de armazenamento, incluindo casos de sucesso e falha para salvar e carregar o cofre;
 - testes unitários white-box para navegação e transições de estado; 
 - golden files visuais em 80×24 por tela; 
 - testes de comandos para cada tela e fluxo de usuário;
 - testes de integração para fluxo completo realizando todas as operações principais (criar cofre, criar segredo, editar segredo, etc.)

# Interface e Experiência do Usuário (TUI e UX)

## 1. Princípios de Design Visual e Interação
- **Interface TUI Moderna:** Ocupa todo o terminal, interativa e desenhada em modo texto com suporte a 256 cores, garantindo alto contraste e legibilidade (estética Cyberpunk sugerida).
- **Acessibilidade de Controles:** Suporte integral à navegação por teclado e suporte complementar a mouse (cliques em campos para edição, nós da árvore para seleção e botões para ações).
- **Menu de Ajuda de Contexto:** Tela ou barra de ajuda sempre acessível, mostrando as ações disponíveis e os respectivos atalhos para o contexto atual (ocultando ações não aplicáveis).
- **Status Global:** Manter visíveis, em uma barra de status: caminho do arquivo atual, status do cofre (`Cofre Salvo` ou `Cofre Modificado`) e total de segredos. Apenas um cofre pode ser aberto por vez.
- **Status individual de segredos:** Exibir indicadores visuais para segredos favoritos, novos/alterados desde o último salvamento e segredos atualmente em edição ou criação.
- **Feedback das operações:**
  - Ações demoradas: mostrar indicador de progresso durante a execução e resultado ao final (não bloqueante, auto-oculta).
  - Ações com resultado visualmente evidente (navegação, expandir/colapsar, toggle de campo, favoritar) dispensam feedback adicional — o indicador visual de estado já é suficiente.
  - Operações interativas: apresentar mensagem informativa e instrucional ao início.
  - Ações destrutivas, irreversíveis ou críticas: exigir confirmação bloqueante antes de executar.

### Modelo de Interação

A interface é composta por **painéis** (áreas funcionais da tela). A interação segue um modelo hierárquico de três níveis:

1. **Painel ativo:** Apenas um painel recebe entrada do teclado por vez. O painel ativo possui destaque visual claro (borda, cor) que o diferencia dos demais.
2. **Foco (elemento focado):** Dentro do painel ativo, um único elemento está realçado — é ele que sofre a ação imediata do usuário (ex: expandir pasta, editar campo). Em campos de texto, o cursor indica a posição de digitação.
3. **Contexto e ações:** O contexto é determinado pela combinação do painel ativo, do elemento focado e do estado da operação em andamento. As **ações disponíveis** mudam dinamicamente conforme o contexto e são comunicadas pela barra de ajuda e atalhos de teclado — somente ações aplicáveis ao contexto atual ficam visíveis.

**Regras:**
- Um painel não ativo pode reagir ao contexto do painel ativo. Ex: ao navegar na árvore, o Painel do Segredo exibe o segredo com foco.
- As ações podem ser:
  - globais (disponíveis em todos os contextos, dependem do estado geral da aplicação). 
  Exemplo: "Salvar cofre" só é aplicável se houver um cofre ativo e alterações não salvas. "Sair" é sempre aplicável.
  - locais (disponíveis apenas no painel ativo). 
  - foco (disponíveis apenas no elemento focado dentro do painel ativo). Exemplo: "Excluir" é aplicável apenas se o elemento focado for um segredo ou pasta.
  - navegação (setas, tab, etc.) é sempre aplicável, mas o comportamento específico depende do foco e do painel ativo. Exemplo: seta direita expande pasta na árvore, mas não tem efeito no Painel do Segredo.
- Uma mesma ação pode ocorrer em painéis ou campos diferentes, se o contexto for adequado. Exemplo: "Favoritar" é aplicável tanto para segredos focados na árvore quanto para o segredo exibido no Painel do Segredo.


### Diretrizes de Redação de Mensagens

**Regras gerais (aplicam-se a todos os tipos):**
- Ser direto e específico — descrever a ação ou situação concreta, sem frases genéricas.
- Não usar exclamação nem palavras-rótulo ("Sucesso!", "Erro!", "Atenção!") — a indicação visual (cor + ícone) já diferencia o tipo.
- Não mencionar teclas — elas são apresentadas por um mecanismo dedicado.
- Não usar "com sucesso", "realizado" ou "Tem certeza que deseja..." — são redundantes ou verbosos.

- **Padrão de Mensagens (Toast/Non-blocking):**
  - **Erro:** Cor vermelha, ícone ❌
  - **Aviso/Alerta:** Cor amarela, ícone ⚠️
  - **Confirmação (Blocking):** Cor laranja, ícone ❓
  - **Sucesso:** Cor verde, ícone ✅
  - **Informativa:** Cor azul/cinza, ícone ℹ️

**Por tipo:**
- **Confirmação** (bloqueante):
  - Indicar o impacto da ação. Opções devem ser verbos de ação ("Excluir", "Salvar", "Voltar"), nunca "Sim/Não/Cancelar".
  - Ex: "Excluir segredo? Ele poderá ser restaurado até o próximo salvamento."
- **Sucesso** (não bloqueante, auto-oculta):
  - Mencionar a ação executada de forma breve.
  - Ex: "Segredo copiado para área de transferência."
- **Erro** (não bloqueante, auto-oculta com tempo estendido):
  - Descrever o que falhou e, se possível, sugerir correção. Evitar jargões técnicos.
  - Ex: "Não foi possível salvar: caminho de arquivo inválido."
  - Em falhas de escrita após a geração de backup, informar explicitamente que existe um arquivo de backup disponível para intervenção manual do usuário.
- **Alerta/Aviso** (não bloqueante, auto-oculta com tempo estendido):
  - Descrever a situação e recomendar ação, sem alarmismo.
  - Ex: "O cofre será bloqueado em 30 segundos por inatividade."
- **Informação** (não bloqueante, auto-oculta):
  - Fornecer instrução ou dado relevante para o contexto atual.
  - Ex: "Informe os dados do novo segredo."
  - Em importações com conflito de segredos, exibir mensagem informando que os itens conflitantes foram importados com sufixos numéricos incrementais.
  - Em importações com conflito de pastas ou modelos, não exibir mensagem específica: o merge de pastas e a substituição de modelos ocorrem silenciosamente.

## 2. Layout e Responsividade
- A interface principal é dividida em dois painéis:
  - **Painel da Hierarquia:** Dedicado à navegação na hierarquia do cofre.
  - **Painel do Segredo:** Dedicado à visualização, criação e edição do segredo selecionado.
- **Tamanho Mínimo do Terminal:**
  - O tamanho mínimo é determinado pelas dimensões necessárias para exibir a tela inicial com o ASCII art, molduras e demais elementos visuais. O valor exato será definido posteriormente.
  - Se o terminal estiver abaixo do tamanho mínimo, a aplicação oculta os painéis e exibe apenas uma mensagem pedindo para redimensionar o terminal.

## 3. Navegação e Hierarquia (Painel da Hierarquia)
- **Comportamento da Árvore:**
  - A ordem de exibição é idêntica à ordem de armazenamento no JSON (mostrando primeiro segredos, depois subpastas).
  - A navegação com setas permite subir/descer na lista. Digitar letras avança o foco para o próximo item correspondente alfabeticamente.
  - Pastas podem ser colapsadas ou expandidas (exibindo um indicador se possuírem conteúdo).
  - Pastas exibem no nome a quantidade total de segredos que possuem (somando subpastas).
  - O scroll acompanha automaticamente a navegação do foco em cofres grandes.
  - Ao remover/adicionar itens, o foco tenta se manter no mesmo nó; se não existir mais, vai para o próximo item ou recua para o nó pai.
- **Criação Relativa ao Foco:**
  - Se o foco estiver em uma **pasta**, o novo segredo/pasta será criado no final dela.
  - Se o foco estiver em um **segredo**, o novo segredo/pasta será criado logo abaixo dele, na mesma pasta.
- **Indicadores Visuais de Nó:**
  - Destaque em ícones ou cor para segredos Favoritos, Novos/Modificados (desde o último salvamento) e itens Atualmente em Edição.
- **Pastas Virtuais (Agrupamentos Lógicos):**
  - **Favoritos:** Visível no topo da raiz (apenas se houver favoritos). Lista atalhos para os segredos favoritados, sem alterar sua localização real. Interagir aqui permite visualizar/editar normalmente.
  - **Lixeira (Materialização do Soft Delete):** Como decisão de interface para representar a exclusão reversível de segredos, a aplicação exibirá uma pasta virtual "Lixeira" no final da raiz (apenas se houver segredos excluídos). Ela lista apenas segredos excluídos reversivelmente, permite restauração e é esvaziada irreversivelmente ao salvar o cofre.
- **Mecanismo de Busca:**
  - Ao buscar, a árvore é filtrada ocultando o que não corresponde.
  - Enquanto a pesquisa estiver ativa, todas as ações ficam indisponíveis exceto: sair da aplicação, navegar pelo cofre e visualizar segredo. O usuário confirma a pesquisa selecionando o elemento desejado (encerrando-a implicitamente) ou cancela a pesquisa; em ambos os casos, o cofre retorna ao estado anterior ao início da pesquisa.
  - Quando o casamento ocorrer no nome do segredo, o trecho exato pesquisado recebe *highlight* (destaque de cor) dentro do nome do segredo encontrado na árvore.

## 4. Visualização e Edição (Painel do Segredo)
- Exibe o detalhe do item focado no Painel da Hierarquia. Se não houver foco, exibe *placeholder* informativo.
- A navegação com `Tab` alterna o controle entre o Painel da Hierarquia e o Painel do Segredo.
- **Privacidade Padrão:** Os campos do tipo "senha" (ou texto sensível) são carregados ocultos (ex: `****`), necessitando de ação explícita (toggle) para exibir o valor.
- **Modos de Edição do Segredo:**
  - A interface separa a edição de segredos em dois modos: edição padrão e edição avançada.
  - **Edição padrão:** permite alterar nome, observação e valores dos campos já existentes.
  - **Edição avançada:** permite alterar a estrutura do segredo, incluindo adicionar campos, renomear campos, excluir campos e reordenar campos.
  - A edição avançada não permite alterar o tipo de um campo existente; para isso, o campo deve ser excluído e recriado com o tipo desejado.
  - O usuário pode alternar entre edição padrão e edição avançada durante a edição do mesmo segredo.
  - Segredos criados a partir de um modelo iniciam na edição padrão; segredos criados vazios iniciam na edição avançada.
- **Gerenciamento de Espaço:**
  - Se um modelo tiver muitos campos, o Painel do Segredo deve permitir scroll vertical.
  - O campo de "Observação" é redimensionável e deve ocupar automaticamente o espaço livre restante no painel.
- **Área de Transferência Integrada:**
  - Ao copiar um campo, disparar um *toast* de sucesso.
  - Exibir um indicador visual na interface com um *countdown* baseado na configuração do cofre, usando 30 segundos como valor padrão sugerido, após o qual a área de transferência será limpa por segurança.
- **Seleção de Arquivos (File Picker):**
  - Para criar, abrir ou salvar cofres com novos nomes, usar uma janela de File Picker integrada à TUI (não exigir digitação cega de caminhos), suportando navegação por setas, mouse e autocompletar.

## 5. Feedback e Segurança (Prevenção de Erros)
- **Aviso Fundamental:** Na criação de um cofre, alertar categoricamente sobre a Irrecuperabilidade ("Zero Knowledge"): o esquecimento da Senha Mestra resulta em perda total dos dados.
- **Ações Destrutivas Irreversíveis:** Excluir pastas (mesmo vazias), excluir Modelos de Segredo e alterar a Senha Mestra exigem um pop-up de confirmação explícita.
- **Soft Delete vs Hard Delete:** Ao excluir um segredo, exibir um aviso discreto de que ele foi movido para a Lixeira.
- **Proteção de Arquivos:** 
  - Alertar imediatamente se o arquivo do cofre for bloqueado por outro processo (Lock File).
  - Ao criar novo cofre, salvar no caminho atual ou usar Salvar Como sobre um arquivo existente, gerar um novo `.bak` antes da substituição.
  - Se já existir um `.bak` anterior, renomeá-lo temporariamente para `.bak2` durante a operação e removê-lo ao final em caso de sucesso, preservando apenas o backup mais recente.
  - Se a escrita do novo arquivo falhar depois da geração do backup, exibir erro explícito informando a falha e a existência do backup para intervenção manual do usuário.
- **Sair da Aplicação:**
  - Pode ser disparado a qualquer momento, mas se houver alterações não salvas, a aplicação deve ser retida por um menu perguntando: "Salvar", "Sair sem Salvar (Descartar)" ou "Voltar".
- **Bloqueio por Inatividade:**
  - Tempo configurável pelo usuário (padrão sugerido: 2 min).
  - O alerta de "Bloqueio Iminente" deve aparecer quando transcorrerem 75% do tempo configurado de inatividade.
  - Contam como atividade: entradas de teclado e cliques de mouse. Movimento de mouse sem clique não conta como atividade.
  - O cronômetro de inatividade deve ser reiniciado ao término de cada ação iniciada pelo usuário, seja ela rápida ou demorada.
  - Após o alerta, qualquer atividade válida aborta o bloqueio iminente e reinicia o cronômetro. Se ignorado até o limite configurado, o cofre é bloqueado.

## 6. Fluxos de Usuário Específicos
- **Criação/Alteração de Senha Mestra:** Ao definir ou alterar a senha mestra, o fluxo exige digitação dupla da nova senha para prevenir erros. Não é exigida a senha atual, pois o cofre já está desbloqueado e autenticado na sessão.

## 7. Mapa de Teclas de Comando por Contexto

Ctrl+C não deve interromper abruptamente a aplicação, mas funcionar 
como um comando convencional.

**Global / Qualquer contexto:**
 - Ctrl+Q: (Sair) 
    - Confirmar sair sem modificações no cofre: novamente Ctrl+Q (Sair)
    - Confirmar sair com modificações no cofre: Ctrl+S (Salvar)
    - Confirmar sair sem salvar: Ctrl+D (Descartar)
    - Cancelar ação de sair: Esc (Voltar)

**Foco no Painel da Hierarquia:**
  - Setas para cima/baixo: mover foco visual entre linhas.
  - Seta para a Direita:
    - Em pasta colapsada: expandir a pasta.
    - Em pasta expandida: mover foco para o primeiro filho.
    - Em segredo: sem efeito.
  - Seta para a Esquerda:
    - Em pasta expandida: colapsar a pasta.
    - Em pasta colapsada (ou segredo): mover foco para a pasta pai.
  - Digitar texto (a-z): pular foco na árvore em direção ao próximo item alfabético.
  - Enter: 
    - Em pasta: expande/colapsa.
    - Em segredo: visualiza e move o foco da aplicação para o Painel do Segredo.
  - Ctrl+N: Abrir tela de "Novo Segredo" no Painel do Segredo, visando a localização do nó atual.
  - Ctrl+F: Ativar barra de filtragem da árvore.

**Foco no Painel do Segredo (Exibindo Detalhe do Segredo):**
  - Setas para cima/baixo: mover foco entre os campos preenchidos.
  - Esc: devolver o foco da interface para o Painel da Hierarquia (retornar à árvore).

**Foco no Painel do Segredo (Editando um Segredo):**
  - *(Área aberta para definição de atalhos de reordenação de campos e submissão de formulário)*