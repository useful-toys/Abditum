# Abditum - Cofre de Senhas Portátil e Seguro

## O que é
Abditum é um cofre de senhas portátil, seguro e fácil de usar, com uma interface TUI moderna. 
Ele permite que os usuários armazenem e gerenciem suas senhas e informações confidenciais de forma organizada e protegida, sem depender de serviços em nuvem ou instalações complexas.
 
## Diferenciais
 - A aplicação é portátil e segura — um único arquivo executável que qualquer pessoa pode copiar e usar discretamente em qualquer lugar, sem persistir dados fora do arquivo do cofre (exceto artefatos transitórios da própria aplicação). 
 - O controle e a propriedade dos dados ficam inteiramente nas mãos do usuário, sem depender de terceiros ou serviços em nuvem.
 - O formato do segredo é flexível e personalizável, permitindo que os usuários criem seus próprios modelos de segredo com campos personalizados, além de oferecer modelos padrão para tipos comuns de segredos.
 
## Conceitos fundamentais
O cofre é o espaço seguro onde o usuário organiza seus segredos. Cada segredo representa uma credencial ou informação confidencial — como o acesso a um site, um cartão de crédito ou uma chave de API — e é composto por campos. Os campos comuns armazenam informações visíveis, como nome do serviço ou usuário. Os campos sensíveis, como senhas e códigos, permanecem ocultos por padrão.

## Como protegemos seus dados
 
O Abditum foi projetado para que seus dados nunca estejam acessíveis a ninguém além de você.
 
 - **Criptografia forte**: o cofre é protegido por criptografia AES-256, uma das mais seguras disponíveis atualmente, garantindo que seus dados permaneçam protegidos mesmo se o arquivo do cofre for comprometido.
 - **Conhecimento zero**: a aplicação não possui meios de acessar ou recuperar seus dados sem a senha mestra. Nem o desenvolvedor do produto tem acesso ao conteúdo do seu cofre.
 - **Privacidade local**: o Abditum não armazena nenhuma informação fora do arquivo do cofre e não acessa a rede ou serviços em nuvem. O arquivo do cofre fica onde você decidir — inclusive em um serviço de nuvem, se você optar por isso.
 - **Proteção em memória**: ao bloquear o cofre, a aplicação minimiza a retenção de dados sensíveis em memória e limpa os buffers sob seu controle sempre que possível. A senha mestra é sobrescrita e buffers sensíveis são descartados.
 - **Proteção visual**: campos sensíveis ficam ocultos por padrão e a interface foi pensada para não chamar atenção — ideal para uso em ambientes públicos ou compartilhados.
 - **Responsabilidade compartilhada**: a Observação é um campo sempre visível para notas do usuário. Por essa razão, **não deve ser utilizada para dados sensíveis** como senhas, tokens ou informações pessoais — use campos sensíveis do modelo para isso. O uso responsável da Observação é de responsabilidade do usuário.

## Conceitos (Glossário)

- **Senha mestra**: chave de acesso ao cofre, usada para criptografar e descriptografar os dados
- **Senha falsa de coação** *(Duress Password)*: senha mestra alternativa que abre uma versão restrita do cofre, protegendo os dados reais em situações de ameaça
- **Cofre**: arquivo criptografado que armazena os segredos do usuário
  - **Bloqueio do cofre**: interrupção do acesso ao conteúdo do cofre, exigindo nova autenticação para retomar
- **Segredo**: item individual dentro do cofre, composto por campos
  - **Segredo favorito**: segredo marcado pelo usuário como prioritário, com destaque para acesso rápido
- **Campo**: elemento individual de um segredo, com nome e valor. Existem dois tipos de campo:
  - **Campo comum**: campo com valor sempre visível, como nome do serviço ou usuário
  - **Campo sensível**: campo com valor oculto por padrão, como senha ou chave de API
- **Observação**: campo comum especial que existe automaticamente em todo segredo, não pode ser renomeado ou excluído
- **Pasta**: estrutura que agrupa segredos e outras pastas dentro do cofre
- **Modelo de segredo**: estrutura predefinida de campos para agilizar a criação de segredos

## Requisitos Funcionais

### Ciclo de Vida do Cofre
- Criar novo cofre em um arquivo com senha mestra
  - Exigir digitação dupla da senha mestra para confirmação
  - Criar automaticamente a Pasta Geral (raiz) contendo as subpastas: Sites e Apps, Financeiro (nesta ordem)
  - Criar automaticamente os modelos padrão:
    - Login: URL (comum), Usuário (comum), Senha (sensível)
    - Cartão de Crédito: Titular (comum), Número (sensível), Validade (comum), CVV (sensível)
    - Chave de API: Serviço (comum), Chave (sensível)
- Abrir cofre a partir de um arquivo existente com senha mestra
  - Validar arquivo contra corrupção e senha incorreta. Existem duas categorias de erro: erros de autenticação (permitem nova tentativa) e erros de integridade (impedem abertura)
  - Se o arquivo estiver corrompido, exibir mensagem de erro genérica (sem detalhes técnicos sobre o tipo de corrupção) e impedir abertura (sem opção de recuperação). A criptografia adotada não permite recuperação parcial — o usuário é responsável por manter cópias de segurança do arquivo do cofre
  - Se a Pasta Geral não existir no arquivo, rejeitar com mensagem de erro (arquivo inválido ou corrompido) — não tentar recriar
  - Se a senha for incorreta, exibir mensagem de erro e permitir nova tentativa
- Salvar cofre no arquivo atual
  - Garantir a integridade e evitar corrompimento do arquivo
  - Segredos marcados para exclusão são removidos permanentemente
  - Usar a senha ativa na sessão — não solicitar novamente
  - Se o arquivo foi modificado externamente desde a última leitura ou salvamento, avisar o usuário e oferecer as opções: Sobrescrever / Salvar como novo arquivo / Cancelar
- Salvar cofre em outro arquivo
  - O arquivo de destino não pode ser o mesmo arquivo atual do cofre
  - Segredos marcados para exclusão são removidos permanentemente
  - Após a operação, o arquivo de trabalho atual passa a ser o novo arquivo
  - Próximas modificações e salvamentos ocorrem sobre o novo arquivo, não o original
  - Usar a senha ativa na sessão — não solicitar novamente
- Descartar alterações não salvas e recarregar o cofre
  - Descarta todas as alterações realizadas desde o último salvamento (se houver) ou desde a abertura do cofre (se nunca foi salvo)
  - O cofre é recarregado ao seu estado anterior
  - Se o arquivo foi modificado externamente desde a última leitura ou salvamento, avisar o usuário antes de recarregar
  - Usar a senha ativa no momento do descarte — não solicitar novamente (pode ser a senha original de abertura ou a nova senha, caso a senha mestra tenha sido alterada durante a sessão)
- Alterar a senha mestra do cofre
  - Exigir digitação dupla para confirmação
  - A alteração é imediata: o cofre é salvo automaticamente com a nova senha ao confirmar, incluindo todas as alterações pendentes da sessão
  - Após a alteração, não é possível descartar essa operação (o arquivo já foi regravado)
- Bloquear o cofre manualmente ou automaticamente após inatividade
  - Bloquear automaticamente após tempo configurável de inatividade, com valor padrão de 5 minutos. Qualquer interação do usuário com a aplicação reseta o temporizador de inatividade
  - O bloqueio retorna ao fluxo de abertura do cofre, exigindo nova autenticação para retomar o acesso
  - Ao bloquear, a senha mestra é limpa da memória (sobrescrita com zeros) e buffers sensíveis são descartados
  - Implementação esforça-se por usar memória protegida (mlock/VirtualLock quando disponível) durante a sessão para impedir swap do arquivo de memória para disco. Se memória protegida não estiver disponível, a aplicação opera normalmente sem essa camada de proteção
- Sair da aplicação
  - Se houver alterações não salvas, exibir confirmação com opções: Salvar e Sair / Descartar e Sair / Cancelar
  - Se não houver alterações pendentes, sair diretamente sem confirmação
  - Ao sair, a senha mestra é limpa da memória e buffers sensíveis são descartados (mesmo comportamento do bloqueio)
- Exportar cofre para arquivo JSON
  - O arquivo exportado contém toda a estrutura do cofre: pastas, segredos ativos e modelos. Configurações de timers não são exportadas
  - Exibir aviso sobre os riscos de segurança e solicitar confirmação antes de exportar
  - Segredos marcados para exclusão não são incluídos no arquivo exportado
- Importar cofre de arquivo JSON
  - O arquivo JSON deve ser válido em estrutura e conteúdo; a existência da Pasta Geral é uma premissa. Se o JSON for inválido ou não contiver Pasta Geral, a importação falha com mensagem de erro
  - Pastas importadas que já existem no cofre (mesmo caminho completo na hierarquia) têm seu conteúdo mesclado automaticamente; pastas com mesmo nome mas em caminhos diferentes são tratadas como pastas distintas. Se a mesclagem resultar em subpastas com nomes conflitantes dentro da mesma pasta de destino, as subpastas importadas são renomeadas automaticamente com sufixo numérico e o usuário é avisado
  - Novas pastas e novos segredos importados são inseridos ao final da lista existente na pasta de destino, após todos os elementos (incluindo segredos marcados para exclusão)
  - Segredo importado com mesmo identificador interno de um segredo já existente no cofre: o segredo existente é preservado inalterado; o segredo importado recebe um novo identificador único e é inserido como segredo independente, com os dados (nome, campos, valores, observação) vindos da importação
  - Modelo importado com mesmo identificador interno de um modelo já existente no cofre: o modelo existente é substituído silenciosamente — nome e estrutura de campos (nomes, tipos, ordem) são sobrescritos pelo modelo importado, mantendo sua posição na lista de modelos. Segredos previamente criados a partir do modelo não são afetados
- Configurar o cofre
  - Todos os tempos são iniciados com valor padrão ao criar o cofre e podem ser ajustados pelo usuário via configurações do cofre. Nenhum temporizador pode ser desabilitado — todos são obrigatórios
  - Configurar tempo de bloqueio automático por inatividade (padrão: 5 minutos)
  - Configurar tempo de ocultação automática de campo sensível (padrão: 15 segundos)
  - Configurar tempo de limpeza automática da clipboard (padrão: 30 segundos)

### Consulta dos Segredos
- Exibir o cofre com suas pastas e segredos
- Buscar segredos por nome, nome de campo, valor de campo comum ou observação
  - A busca funciona por substring, ignorando acentuação e capitalização (case-insensitive)
  - Campos sensíveis nunca participam da busca
  - Segredos marcados para exclusão não aparecem nos resultados de busca
- Exibir um segredo com nome, seus campos e a observação
- Exibir temporariamente o valor de um campo sensível
  - Ocultar o valor automaticamente após tempo configurável, com valor padrão de 15 segundos
- Copiar temporariamente o valor de qualquer campo para a área de transferência
  - Remover o valor da área de transferência automaticamente ao bloquear ou encerrar a aplicação, ou após tempo configurável, com valor padrão de 30 segundos. A limpeza da clipboard depende do suporte do sistema operacional

### Gerenciamento de Segredos
- Criar segredo
  - A partir de um modelo existente ou como segredo sem campos de modelo — apenas com a Observação
  - O segredo pertencerá a uma pasta, escolhida no momento da criação
- Duplicar segredo existente
  - O segredo duplicado é criado na mesma pasta do original, imediatamente após o original na lista
  - O segredo duplicado recebe nome ajustado automaticamente — ex: "Segredo (1)", "Segredo (2)"
  - O histórico de modelo do segredo original é preservado no segredo duplicado
- Editar segredo: alterar o nome do segredo, o valor de campos e/ou observação
  - Não altera a estrutura do segredo (para alterar estrutura, use Adicionar/Renomear/Reordenar/Excluir campo)
- Alterar estrutura do segredo: adicionar campo (com nome e tipo); renomear campo; reordenar campos; excluir campo
  - Não permite alterar o tipo de um campo
  - Não permite alterar a posição, tipo ou nome da observação
  - A Observação ocupa sempre a última posição na lista de campos — campos adicionados pelo usuário são posicionados acima dela
  - A Observação não pode ser movida — apenas campos do usuário participam da reordenação
- Favoritar e desfavoritar segredo
- Marcar e desmarcar segredo para exclusão
  - Segredo marcado para exclusão permanece na lista, apenas sinalizado visualmente
  - Segredo marcado para exclusão é removido permanentemente ao salvar o cofre
- Mover segredo para outra pasta
  - O segredo será movido para a pasta de destino escolhida
- Reordenar segredo dentro da mesma pasta
  - A ordenação é mantida em memória durante a sessão
  - Múltiplas reordenações antes de salvar resultam em estado final apenas (histórico de movimentos descartado)
  - Ao salvar, a ordem final é persistida no arquivo

### Gerenciamento de Pastas
- Criar pasta dentro de outra pasta
  - Não é permitido criar pasta com nome igual ao de outra subpasta existente na mesma pasta pai
- Renomear pasta
  - Não é permitido renomear para um nome já utilizado por outra subpasta na mesma pasta pai
- Mover pasta para outra pasta
  - A pasta será movida para a pasta de destino escolhida
  - O sistema valida e impede movimentos que criariam ciclos na hierarquia (ex: mover Pasta A para dentro de Pasta B se B já está dentro de A)
  - Não é permitido mover para uma pasta que já contenha subpasta com o mesmo nome
- Reordenar pasta dentro da mesma pasta
  - A ordenação é mantida em memória durante a sessão
  - Múltiplas reordenações antes de salvar resultam em estado final apenas (histórico de movimentos descartado)
  - Ao salvar, a ordem final é persistida no arquivo
- Excluir pasta
  - Ao excluir uma pasta, seus segredos e subpastas são movidos para a pasta que a continha
  - Segredos movidos são adicionados ao final da lista de segredos da pasta que a continha
  - Subpastas movidas são adicionadas ao final da lista de pastas da pasta que a continha
  - Se alguma subpasta promovida tiver o mesmo nome de uma subpasta já existente na pasta pai, ela é renomeada automaticamente com sufixo numérico — ex: "Config (1)", "Config (2)". O usuário é avisado sobre as renomeações ocorridas

### Gerenciamento de Modelos de Segredo
- Criar modelo de segredo com campos personalizados
- Renomear modelo de segredo
  - Não é permitido renomear para um nome já utilizado por outro modelo existente
- Alterar estrutura do modelo de segredo: adicionar campo (com nome e tipo); renomear campo; alterar tipo de campo; reordenar campos; excluir campo
  - Campos do modelo permitem nomes duplicados entre si
  - Alterações na estrutura afetam apenas criações futuras de segredos a partir do modelo
  - Os segredos previamente criados a partir do modelo não são afetados, mantendo seus campos inalterados
- Excluir modelo de segredo
- Criar modelo a partir de um segredo existente
  - A Observação automática do segredo é sempre ignorada — não é copiada para o modelo
  - Campos criados manualmente pelo usuário com o nome "Observação" são tratados como campos comuns e incluídos normalmente no modelo gerado

## Regras Transversais

### Estrutura e Pertencimento
- Segredo só pertence a uma pasta
- Pasta pode conter segredos e outras pastas
- O cofre sempre contém a pasta Geral

### Hierarquia de Pastas
- Pastas formam uma estrutura em árvore com a Pasta Geral como raiz
- Ciclos não são permitidos — uma pasta nunca pode ser movida para dentro de seus próprios descendentes
- Cada pasta tem exatamente um ancestral direto (exceto a Pasta Geral, que é a raiz)
- Todas as pastas devem ser navegáveis a partir da Pasta Geral — nenhuma pasta pode ficar desconectada da hierarquia

### Pasta Geral
- A pasta Geral não pode ser renomeada
- A pasta Geral não pode ser movida
- A pasta Geral não pode ser excluída
- A pasta Geral pode estar vazia
- A pasta Geral, por ser a raiz da hierarquia, é o destino natural quando segredos/subpastas de uma pasta diretamente dentro dela são movidos por exclusão — o mesmo comportamento se aplica a qualquer nível: o destino é sempre a pasta pai imediata da pasta excluída

### Nomes e Duplicidade
- Não há restrição quanto a duplicidade do nome entre segredos
- Não há restrição quanto a duplicidade do nome entre campos de um mesmo segredo
- Não é permitido ter duas subpastas com o mesmo nome dentro da mesma pasta pai
- Não é permitido ter dois modelos de segredo com o mesmo nome
- Não há restrição quanto a duplicidade do nome entre campos de um mesmo modelo de segredo
- Exceção à regra de não restrição de duplicidade: Ao duplicar um segredo, se houver conflito de nome entre segredos da mesma pasta, o nome do segredo duplicado será ajustado automaticamente — ex: "Segredo (1)", "Segredo (2)"

### Limites
- Não há limite de quantidade para: pastas, segredos, modelos, campos em segredo, campos em modelo
- Limites são regidos pelo bom senso e pelos recursos do sistema

### Ordenação
- A ordenação de segredos, pastas e campos é mantida pela ação do usuário (reordenação manual)
- A ordenação é persistida no arquivo do cofre
- A ordenação inicial de novos elementos é escolhida pelo usuário no momento da criação
- Modelos de segredo são sempre exibidos em ordem alfabética — não são reordenáveis manualmente

### Segredos e Modelos
- O segredo criado a partir de um modelo não mantém vínculo com ele — o nome do modelo é registrado apenas como histórico
- Edições na estrutura de um segredo não alteram o modelo usado na sua criação
- Edições na estrutura de um modelo afetam apenas criações futuras de segredos a partir dele
- Edições na estrutura de um modelo não alteram os segredos previamente criados a partir do modelo

### Observação
- Observação é um campo que existe automaticamente em todo segredo, independente de como foi criado
- Observação não pode ser renomeada
- Observação não pode ser excluída
- Observação é um campo comum (sempre visível, não sensível)
- Observação ocupa sempre a última posição na lista de campos do segredo
- A Observação automática não faz parte da estrutura do modelo — modelos não declaram nem controlam a Observação
- Se o usuário criar manualmente um campo chamado "Observação" no modelo, esse campo é tratado como campo comum regular e coexistirá com a Observação automática do segredo (que permanece na última posição)
- O uso responsável da observação é por conta e risco do usuário — o campo não prevê ocultação nem tratamento especial

### Gerenciamento de Senha na Sessão
- A senha é fornecida uma única vez ao abrir o cofre
- A senha ativa é usada para todas as operações de criptografia durante a sessão (salvar, descartar)
- Não há re-solicitação de senha para salvamento ou descarte
- Alterar a senha mestra é uma operação imediata e irrevogável — o cofre é regravado na hora; não é uma mudança de estado da sessão
- Ao bloquear o cofre, a senha é removida da memória e será novamente solicitada na próxima abertura

## Requisitos Não Funcionais

- **Criptografia**: AES-256-GCM para criptografia dos dados; Argon2id para derivação de chave a partir da senha mestra
- **Formato de armazenamento**: JSON criptografado com AES-256-GCM, encapsulado em arquivo binário com extensão `.abditum`
- **Compatibilidade**: Windows, macOS e Linux

## Requisitos v2

### Duress Password (Senha Falsa de Coação)
- Criar duress password ao criar novo cofre ou ao configurar cofre existente
  - Exigir digitação dupla para confirmação
  - Duress password deve ser diferente da senha mestra
- Abrir cofre com duress password
  - Validar credenciais (tenta duress password primeiro, se falhar tenta senha mestra)
  - Usuário não é informado qual foi validada
  - Abre "versão restrita" do cofre com segredos/pastas sensíveis ocultos
- Alterar ou remover duress password durante uso normal do cofre
- Configurar quais segredos/pastas são visíveis na versão restrita

### Exibição Parcial de Campos Sensíveis
- Permitir configurar exibição parcial de campos sensíveis — revelar apenas parte do valor (ex: últimos 4 dígitos de número de cartão de crédito)
- Permitir que o usuário defina, por campo sensível, uma regra de exibição parcial (ex: mostrar últimos N caracteres, primeiros N caracteres, ou padrão mascarado como "•••• •••• •••• 1234")
- A exibição parcial não substitui a revelação completa — é um modo adicional de visualização rápida

#### Decisões Pendentes (v2)
- **Duress password — armazenamento e validação**: O design de como a duress password é armazenada, criptografada e validada sem comprometer a senha mestra será definido durante o planejamento de v2
- **Duress password após alteração de senha mestra**: Se o usuário altera a senha mestra durante uma sessão normal, qual é o relacionamento com duress password? Continua usando a senha mestra antiga ou nova? Como persiste essa mudança na próxima sessão?
- **Exibição parcial**: Regras de mascaramento são por campo individual ou por tipo de modelo? Quais padrões pré-definidos oferecer?

## Fora de Escopo (v1)

Funcionalidades deliberadamente excluídas desta versão:
- **Auditoria de senhas**: Análise de força de senha, detecção de duplicatas, avaliação de risco
- **TOTP (Two-Factor Authentication)**: Geração de código de autenticação de dois fatores
- **Backup**: A aplicação não cria, gerencia nem armazena cópias de segurança do cofre. Manter cópias de segurança é responsabilidade exclusiva do usuário
- **Recuperação de dados**: A criptografia adotada não permite recuperação parcial de arquivos corrompidos. Não há mecanismo de reparo, importação forçada ou abertura em modo degradado
