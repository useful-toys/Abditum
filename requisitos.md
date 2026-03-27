# Abditum - Cofre de Senhas Portátil e Seguro

## O que é
Abditum é um cofre de senhas portátil, seguro e fácil de usar, com uma interface TUI moderna. 
Ele permite que os usuários armazenem e gerenciem suas senhas e informações confidenciais de forma organizada e protegida, sem depender de serviços em nuvem ou instalações complexas.
 
## Diferenciais
 - A aplicação é portátil e segura — um único arquivo executável que qualquer pessoa pode copiar e usar discretamente em qualquer lugar, sem persistir dados fora do arquivo do cofre (exceto artefatos transitórios e backups explicitamente previstos pela própria aplicação). 
 - O controle e a propriedade dos dados ficam inteiramente nas mãos do usuário, sem depender de terceiros ou serviços em nuvem.
 - O formato do segredo é flexível e personalizável, permitindo que os usuários criem seus próprios modelos de segredo com campos personalizados, além de oferecer modelos pré-definidos para tipos comuns de segredos.
 
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
  - Criar automaticamente as pastas padrão: Geral, Sites e Apps, Financeiro
  - Criar automaticamente os modelos padrão: Login, Cartão de Crédito, Chave de API
- Abrir cofre a partir de um arquivo existente com senha mestra
  - Validar arquivo contra corrupção e senha incorreta
  - Se o arquivo estiver corrompido, exibir mensagem de erro genérica (sem detalhes técnicos sobre o tipo de corrupção) e impedir abertura (sem opção de recuperação)
  - Se a Pasta Geral não existir no arquivo, rejeitar com mensagem de erro (arquivo inválido ou corrompido) — não tentar recriar
  - Se a senha for incorreta, exibir mensagem de erro e permitir nova tentativa
- Salvar cofre no arquivo atual
  - Garantir a integridade e evitar corrompimento do arquivo
  - Segredos marcados para exclusão são removidos permanentemente
  - Usar a senha originalmente fornecida ao abrir (ou alterá-la caso tenha sido alterada) — não solicitar novamente
- Salvar cofre em outro arquivo
  - O arquivo de destino não pode ser o mesmo arquivo atual do cofre
  - Segredos marcados para exclusão são removidos permanentemente
  - Após a operação, o arquivo de trabalho atual passa a ser o novo arquivo
  - Próximas modificações e salvamentos ocorrem sobre o novo arquivo, não o original
  - Usar a senha originalmente fornecida ao abrir (ou alterá-la caso tenha sido alterada) — não solicitar novamente
- Descartar alterações não salvas e recarregar o cofre
  - Descarta todas as alterações realizadas desde o último salvamento (se houver) ou desde a abertura do cofre (se nunca foi salvo)
  - O cofre é recarregado ao seu estado anterior
  - Usar a senha originalmente fornecida ao abrir — não solicitar novamente
- Alterar a senha mestra do cofre
  - Exigir digitação dupla para confirmação
- Bloquear o cofre manualmente ou automaticamente após inatividade
  - Bloquear automaticamente após tempo configurável de inatividade, com valor padrão sugerido de 5 minutos
  - O bloqueio retorna ao fluxo de abertura do cofre, exigindo nova autenticação para retomar o acesso
  - Ao bloquear, a senha mestra é limpa da memória (sobrescrita com zeros) e buffers sensíveis são descartados
  - Implementação esforça-se por usar memória protegida (mlock/VirtualLock quando disponível) durante a sessão para impedir swap do arquivo de memória para disco
- Exportar cofre para arquivo JSON
  - Exibir aviso sobre os riscos de segurança e solicitar confirmação antes de exportar
- Importar cofre de arquivo JSON
  - Pastas importadas que já existem no cofre (mesmo nome) têm seu conteúdo mesclado automaticamente
  - Segredo importado que já existe no cofre (mesmo identificador interno) é salvo como um novo segredo, preservando todos seus dados
  - Segredo importado com nome conflitante na mesma pasta de destino recebe nome ajustado automaticamente — ex: "Segredo (1)", "Segredo (2)"
  - Modelo importado que já existe no cofre (mesmo identificador interno) é substituído silenciosamente
- Configurar o cofre
  - Configurar tempo de bloqueio automático por inatividade (padrão: 5 minutos)
  - Configurar tempo de ocultação automática de campo sensível (padrão: 15 segundos)
  - Configurar tempo de limpeza automática da clipboard (padrão: 30 segundos)

### Consulta dos Segredos
- Exibir o cofre com suas pastas e segredos
- Buscar segredos por nome, nome de campo, valor de campo comum ou observação
  - A busca funciona por substring, ignorando acentuação e capitalização (case-insensitive)
  - Campos sensíveis nunca participam da busca
- Exibir um segredo com nome, seus campos e a observação
- Exibir temporariamente o valor de um campo sensível
  - Ocultar o valor automaticamente após tempo configurável, com valor padrão sugerido de 15 segundos
- Copiar temporariamente o valor de qualquer campo para a área de transferência
  - Remover o valor da área de transferência automaticamente ao fechar a aplicação ou após tempo configurável, com valor padrão sugerido de 30 segundos

### Gerenciamento de Segredos
- Criar segredo
  - A partir de um modelo existente ou como segredo vazio sem campos iniciais
  - O segredo pertencerá a uma pasta, escolhida no momento da criação
- Duplicar segredo existente
  - O segredo duplicado recebe nome ajustado automaticamente — ex: "Segredo (1)", "Segredo (2)"
  - O histórico de modelo do segredo original é preservado no segredo duplicado
- Editar segredo: alterar o nome do segredo, o valor de campos e/ou observação
  - Não altera a estrutura do segredo (para alterar estrutura, use Adicionar/Renomear/Reordenar/Excluir campo)
- Alterar estrutura do segredo: adicionar campo (com nome e tipo); renomear campo; reordenar campos;  excluir campo
  - Não permite alterar o tipo de um campo
  - Não permite alterar a posição, tipo ou nome da observação
- Favoritar e desfavoritar segredo
- Marcar e desmarcar segredo para exclusão
  - Segredo marcado para exclusão é removido permanentemente ao salvar o cofre
- Mover segredo para outra pasta
  - O segredo será movido para a pasta de destino escolhida
- Reordenar segredo dentro da mesma pasta
  - A ordenação é mantida em memória durante a sessão
  - Múltiplas reordenações antes de salvar resultam em estado final apenas (histórico de movimentos descartado)
  - Ao salvar, a ordem final é persistida no arquivo

### Gerenciamento de Pastas
- Criar pasta dentro de outra pasta
- Renomear pasta
- Mover pasta para outra pasta
  - A pasta será movida para a pasta de destino escolhida
  - O sistema valida e impede movimentos que criariam ciclos na hierarquia (ex: mover Pasta A para dentro de Pasta B se B já está dentro de A)
- Reordenar pasta dentro da mesma pasta
  - A ordenação é mantida em memória durante a sessão
  - Múltiplas reordenações antes de salvar resultam em estado final apenas (histórico de movimentos descartado)
  - Ao salvar, a ordem final é persistida no arquivo
- Excluir pasta
  - Ao excluir uma pasta, seus segredos e subpastas são movidos para a pasta que a continha
  - Segredos movidos são adicionados ao final da lista de segredos da pasta que a continha
  - Subpastas movidas são adicionadas ao final da lista de pastas da pasta que a continha

### Gerenciamento de Modelos de Segredo
- Criar modelo de segredo com campos personalizados
- Editar modelo de segredo: adicionar, renomear, alterar tipo, reordenar e excluir campos
  - Alterações na estrutura de um modelo afetam apenas criações futuras de segredos a partir do modelo
  - Os segredos previamente criados a partir do modelo não são afetados pelas alterações no modelo, mantendo seus campos inalterados
- Excluir modelo de segredo
- Criar modelo a partir de um segredo existente

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
- A pasta Geral não pode ser reordenada
- A pasta Geral não pode ser excluída
- A pasta Geral pode estar vazia
- A pasta Geral é sempre o destino final quando segredos/subpastas são movidos por exclusão de pasta

### Nomes e Duplicidade
- Não há restrição quanto a duplicidade do nome entre segredos
- Não há restrição quanto a duplicidade do nome entre campos de um mesmo segredo
- Não há restrição quanto a duplicidade do nome entre subpastas dentro da mesma pasta
- Não há restrição quanto a duplicidade do nome entre modelos de segredo
- Não há restrição quanto a duplicidade do nome entre campos de um mesmo modelo de segredo
- Exceção à regra de não restrição de duplicidade: Ao duplicar e ao importar segredos, se houver conflito de nome entre segredos da mesma pasta, nome do segredo importado ou duplicado será ajustado automaticamente — ex: "Segredo (1)", "Segredo (2)"

### Limites
- Não há limite de quantidade para: pastas, segredos, modelos, campos em segredo, campos em modelo
- Limites são regidos pelo bom senso e pelos recursos do sistema

### Ordenação
- A ordenação de segredos, pastas e campos é mantida pela ação do usuário (reordenação manual)
- A ordenação é persistida no arquivo do cofre
- A ordenação inicial de novos elementos é determinada pela UX

### Segredos e Modelos
- O segredo criado a partir de um modelo não mantém vínculo com ele — o nome do modelo é registrado apenas como histórico
- Edições na estrutura de um segredo não alteram o modelo usado na sua criação
- Edições na estrutura de um modelo afetam apenas criações futuras de segredos a partir dele
- Edições na estrutura de um modelo não alteram os segredos previamente criados a partir do modelo

### Observação
- Observação é um campo que existe automaticamente em todo segredo
- Observação não pode ser renomeada
- Observação não pode ser excluída
- Observação é um campo comum (sempre visível, não sensível)
- Observação não é declarada no modelo de segredo
- O uso responsável da observação é por conta e risco do usuário — o campo não prevê ocultação nem tratamento especial

### Gerenciamento de Senha na Sessão
- A senha é fornecida uma única vez ao abrir o cofre (ou ao alterar a senha mestra)
- A mesma senha é usada para todas as operações de criptografia durante a sessão (salvar, descartar)
- Não há re-solicitation de senha para salvamento ou descarte
- Se a senha mestra for alterada durante a sessão, a nova senha passa a ser usada para próximas operações
- Ao bloquear o cofre, a senha é removida da memória e será novamente solicitada na próxima abertura

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

#### ⏳ Decisões Pendentes (v2)
- **Duress password após alteração de senha mestra**: Se o usuário altera a senha mestra durante uma sessão normal, qual é o relacionamento com duress password? Continua usando a senha mestra antiga ou nova? Como persiste essa mudança na próxima sessão?

## Fora de Escopo (v1)

Funcionalidades deliberadamente excluídas desta versão:
- **Auditoria de senhas**: Análise de força de senha, detecção de duplicatas, avaliação de risco
- **TOTP (Two-Factor Authentication)**: Geração de código de autenticação de dois fatores

