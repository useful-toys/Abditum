# Abditum - Cofre de Senhas Portátil e Seguro

## O que é

Abditum é um cofre de senhas portátil, seguro e fácil de usar, com uma interface TUI moderna. 
Ele permite que os usuários armazenem e gerenciem suas senhas e informações confidenciais de forma organizada e protegida, sem depender de serviços em nuvem ou instalações complexas.

## Diferenciais

O cofre deve ser completamente portátil e seguro — um único arquivo executável que qualquer pessoa pode copiar e usar discretamente em qualquer lugar, sem persistir dados fora do arquivo do cofre, exceto artefatos transitórios e backups explicitamente previstos pela própria aplicação. 

O controle e a propriedade dos dados ficam inteiramente nas mãos do usuário, sem depender de terceiros ou serviços em nuvem.

O formato do segredo é flexível e personalizável, permitindo que os usuários criem seus próprios modelos de segredo com campos personalizados, além de oferecer modelos pré-definidos para tipos comuns de segredos.

## Requisitos Funcionais

### Ciclo de Vida do Cofre

- Criar novo cofre com caminho e senha mestra
- Abrir cofre existente com caminho e senha mestra
- Descartar alterações não salvas e recarregar cofre
- Salvar cofre no caminho atual
- Salvar cofre em novo caminho
- Alterar a senha mestra do cofre
- Bloquear cofre manualmente (equivale a fechar o cofre, mantendo a aplicação em execução, voltando à tela de abertura do cofre e minimizando a retenção de dados sensíveis em memória, com limpeza dos buffers controlados pela aplicação sempre que possível)
- Bloquear automaticamente após inatividade (mesmo processo do bloqueio manual, mas acionado por temporizador configurável)
- Exportar cofre para formato JSON plain text
- Importar cofre de formato JSON plain text, tratando conflitos com elementos já existentes no cofre atual
  - Regras de tratamento de conflito na importação:
    - Pastas com a mesma identidade da pasta já existente no cofre atual têm sua hierarquia mesclada
    - Se um segredo importado colidir por identidade com um segredo já existente no cofre atual, a aplicação deve criar uma nova identidade para o segredo importado, preservando seus demais dados
    - Se um segredo importado colidir por nome com outro segredo já existente na mesma pasta de destino (ou na raiz, quando o destino for a raiz), a aplicação deve ajustar seu nome com um sufixo numérico incremental para evitar ambiguidade visual. Ex: "Segredo" → "Segredo (1)", "Segredo (2)", etc.
    - Se houver conflito de segredos durante a importação, a aplicação deve informar que os segredos conflitantes por nome foram importados com nomes sufixados incrementalmente
    - Se houver conflito de pastas durante a importação, o merge da hierarquia ocorre silenciosamente
    - Modelos de segredo com a mesma identidade de um modelo já existente no cofre atual são sobrepostos pelo modelo importado
    - Se houver conflito de modelos durante a importação, a substituição pelo modelo importado ocorre silenciosamente
- Configurar o cofre
  - Tempo para bloqueio automático por inatividade
  - Tempo para ocultar valor de campo de segredo após exibição temporária
  - Tempo para limpar a área de transferência automaticamente


### Navegação da Hierarquia do Cofre (somente leitura)
- Exibir hierarquia do cofre (pastas e segredos)
- Exibir detalhes do segredo selecionado
- Exibir temporariamente o valor de um campo de segredo (ex: senha) com opção de ocultar/mostrar
  - Ocultar automaticamente o valor do campo de segredo conforme a configuração do cofre, com valor padrão sugerido de 15 segundos

### Gerenciamento de Segredos
- Criar um segredo
  - Como:
    - Usando um modelo de segredo existente
    - Começando com um segredo vazio, sem nenhum campo inicial
    - Quando o segredo é criado a partir de um modelo, ele não mantém vínculo por referência com o modelo. O nome do modelo é guardado apenas como histórico ("snapshot" do momento da criação). Alterar a estrutura (campos), renomear ou excluir o modelo posteriormente não afeta os segredos já criados.
  - Onde:
    - Na raiz do cofre ou dentro de uma pasta
    - Numa pasta da hierarquia do cofre
- Duplicar um segredo existente
- Favoritar/desfavoritar um segredo
- Editar um segredo existente:
  - Alterar dados do segredo: alterar nome do segredo, alterar os valores nos campos (dados normais e sensíveis) e alterar a observação.
  - Alterar estrutura do segredo:
        - Incluir um novo campo de segredo
            - Informar o nome do campo de segredo
            - Informar o tipo do campo de segredo
        - Alterar um campo de segredo
            - Alterar o nome do campo de segredo
          - Não suportado: alterar o tipo do campo de segredo
        - Excluir um campo de segredo
        - Alterar a posição (reordenar) de um campo de segredo 
- Remover um segredo de forma reversível até o próximo salvamento do cofre
- Restaurar um segredo excluído reversivelmente antes do próximo salvamento do cofre
- Mover um segredo para outra pasta ou para a raiz do cofre
- Mover (reordenar) segredo relativamente a outros segredos dentro da mesma pasta ou raiz do cofre
- Buscar segredos por nome, por nome de campo, por valor de campos do tipo `texto` ou por observação
  - Campos do tipo `texto sensível` nunca participam da busca
  - A observação é considerada dado não sensível e não deve ser usada para armazenar segredos
  - Todos os segredos que satisfaçam qualquer critério de busca devem ser exibidos como resultado, independentemente de haver destaque visual no nome
  - A busca ocorre apenas em memória, após o desbloqueio do cofre, sem manter qualquer índice persistido

### Gerenciamento de Hierarquia
 - Criar pasta
   - Onde:
    - Na raiz do cofre ou dentro de uma pasta
    - Numa pasta da hierarquia do cofre
 - Renomear pasta
 - Mover pasta para 
    - outra pasta
    - raiz do cofre
 - Mover (reordenar) pasta relativamente a outras pastas dentro da mesma pasta ou raiz do cofre
 - Excluir pasta, movendo seus segredos e suas subpastas filhas para a pasta pai (ou para a raiz do cofre, se a pasta excluída estiver na raiz do cofre). Os segredos promovidos são adicionados ao final da lista de segredos do pai, e as subpastas promovidas são adicionadas ao final da lista de pastas do pai.
 - Pastas pré-definidas ao criar um novo cofre, mas editáveis e removíveis pelo usuário
    - Pasta "Sites"
    - Pasta "Financeiro"
    - Pasta "Serviços"
 
### Gerenciamento de Modelos de Segredo
- Criar modelo de segredo com campos personalizados
- Editar modelo de segredo existente
    - A alteração na estrutura de um modelo (adição, modificação ou exclusão de campos) afeta apenas as criações futuras. Os segredos previamente criados a partir deste modelo permanecerão inalterados.
    - Incluir um novo campo de segredo
        - Informar o nome do campo de segredo
        - Informar o tipo do campo de segredo
    - Alterar um campo de segredo
        - Alterar o nome do campo de segredo
        - Alterar o tipo do campo de segredo
    - Excluir um campo de segredo
    - Alterar a posição (reordenar) de um campo de segredo
- Remover modelo de segredo
- Criar modelo de segredo a partir de um segredo existente, copiando os campos do segredo como estrutura inicial para o novo modelo
- Modelos de segredo pré-definidos criados no cofre por padrão, mas editáveis e removíveis pelo usuário
    - Modelo de segredo "Login": campos "URL", "Username" e "Password"
    - Modelo de segredo "Cartão de Crédito": campos "Número do Cartão", "Nome no Cartão", "Data de Validade" e "CVV"
    - Modelo de segredo "API Key": campos "Nome da API", "Chave de API"

### Área de Transferência

- Copiar qualquer campo para a área de transferência
- Limpar automaticamente a área de transferência conforme a configuração do cofre, com valor padrão sugerido de 30 segundos
- Limpar automaticamente a área de transferência ao bloquear ou fechar o cofre

### Segurança
- Ao criar ou alterar a senha mestra, exigir digitação dupla para confirmação
- Os controles de bloqueio automático por inatividade, reocultação temporizada de campos sensíveis e limpeza automática da área de transferência são obrigatórios como mecanismos de redução de exposição e seguem as regras funcionais já definidas nas seções de ciclo de vida do cofre, navegação e área de transferência
- Proteção contra brute force e ataques offline — garantida através de parâmetros rígidos do Argon2id (custo alto de memória e tempo) para atrasar tentativas de derivação de chave
- Minimizar a retenção de dados sensíveis em memória e limpar, sempre que possível, os buffers controlados pela aplicação ao bloquear ou fechar o cofre
- Proteção contra shoulder surfing — atalho para ocultar toda a interface rapidamente
- Ao exportar o cofre para formato JSON plain text, mostrar mensagem de aviso sobre os riscos de segurança envolvidos e pedir confirmação do usuário antes de prosseguir com a exportação

## Requisitos não Funcionais
- Criptografia: AES-256-GCM para criptografia dos dados e Argon2id para derivação de chave a partir da senha mestra
- Armazenamento em formato JSON criptografado
- Compatibilidade com sistemas operacionais Windows, macOS e Linux
- Formato do arquivo de cofre: extensão .abditum
- Confiabilidade (Salvamento Atômico): ao salvar o cofre no caminho atual, escrever os dados primeiro em um arquivo `.abditum.tmp` no mesmo diretório do cofre. Somente após a gravação bem-sucedida, renomear o arquivo, substituindo o `.abditum` original. Em caso de falha, o arquivo `.tmp` deve ser imediatamente apagado para evitar persistência indevida de dados fora do arquivo final do cofre.
- Ao substituir um arquivo de cofre já existente, manter uma cópia de backup do arquivo anterior com extensão .abditum.bak
  - Se já existir um arquivo `.abditum.bak`, a aplicação deve renomeá-lo temporariamente para `.abditum.bak2` antes de gerar o novo backup
  - Se a substituição do arquivo de cofre for concluída com sucesso, o arquivo `.abditum.bak2` deve ser apagado, preservando apenas o novo `.abditum.bak`
  - Se a operação falhar antes de consolidar o novo backup e o novo arquivo principal, a aplicação deve restaurar o `.abditum.bak2` para `.abditum.bak` sempre que possível
- Executável "portable"
    - Não requer instalação, pode ser executado diretamente do arquivo binário
    - Permitir copiar o arquivo binário para qualquer local do sistema de arquivos e executá-lo
    - Não usar arquivos de configuração ou dados fora do arquivo do cofre
  - Não persistir dados fora do arquivo do cofre, exceto artefatos transitórios e backups explicitamente previstos pela própria aplicação
- Os modelos de segredo são armazenados dentro do cofre, permitindo que cada cofre tenha seus próprios modelos personalizados
- Os modelos de segredo pré-definidos são fornecidos populados no cofre ao criar um novo cofre, mas podem ser editados ou removidos pelo usuário
- O cofre é criado com uma hierarquia pré-definida de pastas
- Unicidade de nomes: a aplicação não impõe unicidade de nomes para segredos, pastas ou modelos de segredo. Nomes repetidos são permitidos e tratados apenas como questão de usabilidade, sem impacto na corretude ou integridade dos dados.
- Privacidade: Ausência total de logs de aplicação (stdout/stderr) que contenham caminhos de arquivos de cofre, nomes de segredos ou valores de campos.
- Compatibilidade retroativa: a aplicação de versão N deverá ser capaz de abrir arquivos de cofre criados em qualquer versão anterior do formato suportada pela aplicação. Ao abrir um cofre antigo, o payload descriptografado é migrado em memória para o modelo atual; ao salvar, o arquivo é sempre regravado no formato da versão da aplicação.

## Requisitos Inversos
- Armazenamento na nuvem
- Múltiplos cofres abertos simultaneamente
- App mobile ou web — TUI portátil por design
- Tags — pastas/grupos suficientes para v1
- Histórico de versões de segredos
- Alteração do tipo de um campo de segredo existente — para mudar o tipo de um campo, será necessário excluir o campo existente e adicionar outro campo com o tipo desejado

## Requisitos Postergados (v2)
- Suporte a TOTP (Autenticação de Dois Fatores): Calcular e exibir tokens de 6 dígitos em tempo real para campos que armazenem uma chave secreta (seed) TOTP, com indicação visual do tempo restante.
- Gerador de senhas
- Senha falsa de coação (Duress Password): funcionalidade de ter uma senha alternativa que abre o cofre ocultando os segredos mais importantes, para situações em que o usuário seja forçado a abri-lo.
- Compartilhamento via QR Code: renderizar um QR code diretamente na TUI (usando blocos ASCII) contendo o valor de um campo sensível, permitindo transferência rápida, segura e offline para um smartphone.
- Relatório de Saúde do Cofre (Auditoria): analisar localmente todas as senhas armazenadas e alertar o usuário sobre senhas fracas, reutilizadas em múltiplos segredos ou muito antigas.
- Autenticação de Dois Fatores Offline (Keyfile / Token de Hardware): permitir que o cofre exija, além da senha mestra, um arquivo físico específico (keyfile) ou interação com um token USB (ex: YubiKey) para ser descriptografado.

## Conceitos (Glossário)
- Senha mestra: chave de acesso ao cofre, usada para criptografar e descriptografar os dados
- Cofre: arquivo criptografado que armazena as senhas e informações do usuário
  - Bloqueio do cofre: processo de proteção em que o aplicativo continua em execução, o acesso ao conteúdo do cofre é interrompido, a aplicação minimiza a retenção de dados sensíveis em memória, limpa os buffers sob seu controle sempre que possível e retorna ao fluxo de abertura do cofre, exigindo nova autenticação para retomar o acesso
- Segredo: item individual dentro do cofre, com dados comuns e dados sensíveis
    - Segredo favorito: um segredo marcado pelo usuário como prioritário ou de uso frequente, ganhando destaque para acesso rápido
- Dados: informações armazenadas em um segredo
    - Dados comuns: informações não sensíveis, como nome do serviço ou URL
    - Dados sensíveis: informações confidenciais, como senha, apikeys
    - Observação: campo de texto livre para o usuário adicionar informações adicionais sobre o segredo
- Campo de segredo: elemento individual dentro de um segredo, com nome, tipo e valor, que armazena um dado específico
    - Tipo do campo de segredo: define o tipo de dado que o campo pode armazenar (restrito a texto ou texto sensível).
- Hierarquia do cofre: organização dos segredos em pastas e subpastas dentro do cofre
    - Raiz do cofre: o nível estrutural primário e mais alto da hierarquia, que contém os segredos e pastas que não estão aninhados em outras pastas
    - Pasta: contêiner estrutural utilizado para agrupar e organizar segredos e outras subpastas na hierarquia
    - Pasta virtual: agrupamento lógico gerado pelo sistema que exibe uma visão de segredos localizados em outras pastas ou na raiz do cofre, com base em características específicas, sem alterar a localização real deles na hierarquia
- Modelo de segredo: estrutura para criar segredos com campos específicos, como login, senha, URL, etc.
    - Modelo de segredo pré-definido: modelo de segredo fornecido pelo sistema, com campos comuns para tipos de segredos populares (ex: login, cartão de crédito, apikey)
    - Modelo de segredo personalizado: modelo de segredo criado e estruturado pelo próprio usuário para atender a necessidades específicas de formato
    - Campo modelo de segredo: elemento individual dentro de um modelo de segredo que representa um campo específico a ser preenchido ao criar um segredo a partir do modelo, com nome e tipo definidos
- Senha falsa de coação (Duress Password): senha mestra alternativa configurada para abrir uma versão restrita ou falsa do cofre, protegendo os dados reais em situações de ameaça ou extorsão física
- Conhecimento Zero (Zero Knowledge): princípio de negócio e segurança em que a aplicação não possui meios de acessar ou recuperar os dados sem a senha mestra do usuário
- Exclusão reversível (Soft Delete): mecanismo pelo qual um segredo excluído reversivelmente permanece restaurável até o próximo salvamento definitivo do cofre
- Auditoria de senhas (Saúde do Cofre): processo de análise de segurança para alertar o usuário sobre a presença de senhas fracas, antigas ou reutilizadas
- TOTP (Autenticação de Dois Fatores): código numérico temporário gerado em tempo real pelo cofre a partir de uma chave secreta, servindo como reforço de segurança para acesso a serviços externos
- Shoulder Surfing: técnica de espionagem física mitigada pela aplicação, onde um indivíduo mal-intencionado observa a tela do usuário para roubar informações visíveis

## Modelagem

Decisões de modelagem:
 - **Hierarquia recursiva:** A raiz do cofre funciona como uma pasta sem nome. Pastas podem conter segredos e subpastas em qualquer nível de aninhamento.
 - **Ordenação por posição:** A ordem dos elementos no JSON reflete diretamente a ordem de exibição na interface (segredos, pastas e campos).
 - **Modelo como snapshot:** Segredos criados a partir de modelos não mantêm vínculo por referência — o nome do modelo é guardado apenas como registro histórico. Não há distinção estrutural entre segredos criados com ou sem modelo.
 - **IDs (NanoID, 6 caracteres alfanuméricos):** Segredos, pastas e modelos de segredo possuem ID. O espaço de 62⁶ (~56 bilhões) combinações garante unicidade prática.
 - **Nomes não identificadores:** O nome de pastas, segredos e modelos é apenas um atributo editável e não participa da identidade do elemento. Nomes repetidos são permitidos.
 - **Campos uniformes:** O valor de um campo pode ser string vazia (campo existente, não preenchido). Não há distinção de estado entre preenchido e vazio.
 - **Observação implícita:** Todo segredo possui um campo de observação (opcional, texto livre) que não é declarado nos modelos e não pode ser removido. A observação é tratada como dado não sensível.
 - **Busca sequencial em memória:** O cofre não mantém índices ou estruturas auxiliares de busca, persistidos ou em memória de longa duração. Após o desbloqueio, as buscas são realizadas por varredura sequencial sobre a estrutura carregada, assumindo que o volume de dados do cofre é pequeno o suficiente para esse custo.
 - **Configurações embutidas:** As configurações do cofre são armazenadas dentro do próprio arquivo, garantindo portabilidade total sem arquivos externos.

 - Cofre (Estrutura do Payload JSON Criptografado):
    - configurações:
      - tempo_bloqueio_inatividade_minutos: inteiro (persistido no cofre; padrão sugerido: 2)
      - tempo_ocultar_segredo_segundos: inteiro (persistido no cofre; padrão sugerido: 15)
      - tempo_limpar_area_transferencia_segundos: inteiro (persistido no cofre; padrão sugerido: 30)
    - segredos: list[Segredo]
    - pastas: list[Pasta]
    - modelos_segredo: list[ModeloSegredo]
    - data criação: datetime
    - data última modificação: datetime
 - Segredo:
    - id: nanoid (6 caracteres)
    - nome: string
    - nome do modelo de segredo: string (opcional)
    - campos: list[CampoSegredo]
    - favorito: booleano
    - observação: string (opcional)
    - data criação: datetime
    - data última modificação: datetime
 - Pasta:
    - id: nanoid (6 caracteres)
    - nome: string
    - segredos: list[Segredo]
    - pastas: list[Pasta]
 - Modelo de segredo:
    - id: nanoid (6 caracteres)
    - nome: string
    - campos: list[CampoModeloSegredo]
 - Campo de segredo:
    - nome: string
    - tipo: enum (texto, texto sensível)
    - valor: string (opcional)
 - Campo modelo de segredo:
    - nome: string
    - tipo: enum (texto, texto sensível)

## Formato do Arquivo

O arquivo do cofre (`.abditum`) é um stream binário composto por duas partes:

1.  **Cabeçalho de Criptografia (fixo):**
    *   `magic` (bytes): Assinatura fixa do formato para identificar o arquivo como um cofre Abditum antes de qualquer tentativa de descriptografia. Exemplo recomendado: ASCII `ABDT`.
    *   `versão_formato` (inteiro): Versão do formato do arquivo.
    *   `salt` (bytes): Salt para a derivação da chave Argon2id.
    *   `nonce` (bytes): Vetor de inicialização para a criptografia AES-256-GCM.
2.  **Payload Criptografado:**
    *   O restante do stream contém a estrutura do `Cofre` em formato JSON, criptografada com AES-256-GCM.

Essa estrutura garante que os metadados necessários para a descriptografia estejam disponíveis antes de ler o conteúdo do cofre, que permanece totalmente criptografado. A assinatura `magic` permite rejeitar imediatamente arquivos que não pertencem ao formato Abditum, diferenciando erro de tipo de arquivo de erro de senha ou corrupção criptográfica. Todo o cabeçalho entra como AAD (Additional Authenticated Data) do AES-256-GCM, para que sua integridade também seja validada sem a necessidade de checksum adicional da aplicação.

## Estados e Fluxos Principais

### Invariantes de estado

- Só pode existir um cofre ativo por vez.
- Um segredo não pode estar simultaneamente na hierarquia principal e na Lixeira.
- Um segredo só pode estar na raiz ou em uma pasta, nunca em ambos, nem em duas pastas ao mesmo tempo.
- Uma pasta só pode estar na raiz ou dentro de outra pasta, nunca em ambos, nem em duas pastas ao mesmo tempo.
- O estado do cofre `Cofre Modificado` deve refletir qualquer divergência entre memória e último salvamento persistido.
- A Lixeira só materializa segredos excluídos reversivelmente.
- Ao salvar, segredos na Lixeira são permanentemente excluídos, sem possibilidade de recuperação.
- Pastas não possuem soft delete; sua exclusão sempre remove a pasta e promove os filhos.
- Campos `texto sensível` nunca participam de busca, independentemente do estado visual de ocultação ou exibição.

### Estados principais

#### Estados globais da aplicação

- **Inicial / sem cofre ativo:** a aplicação está em execução, mas ainda não há um cofre ativo. As ações disponíveis se limitam a criar cofre, abrir cofre, acessar ajuda e sair.
- **Abrindo cofre:** estado transitório em que a aplicação realiza os passos para abrir o cofre (coleta o caminho do cofre e a senha mestra, valida o arquivo, deriva a chave, descriptografa o payload e carrega o domínio em memória). Este estado é necessário devido à necessidade de apresentar um subprocesso interativo para navegar pastas, selecionar arquivo, solicitar senha, validar o arquivo, validar a criptografia e carregar os dados.
- **Criando novo cofre:** estado transitório em que a aplicação realiza os passos para criar um novo cofre (coleta o caminho do cofre e a senha mestra, verifica se já existe arquivo no destino, solicita confirmação explícita em caso de sobrescrita, popula a estrutura inicial do cofre e grava o novo arquivo). Este estado é necessário devido à necessidade de apresentar um subprocesso interativo para navegar pastas, selecionar arquivo, solicitar senha, validar o destino, tratar eventual sobrescrita e gravar os dados.
- **Salvando cofre com outro caminho:** estado transitório em que a aplicação realiza os passos para salvar o cofre em um novo caminho (coleta o novo caminho, valida a possibilidade de gravação, trata eventual sobrescrita de arquivo existente, grava o arquivo do cofre no novo caminho e atualiza o caminho atual do cofre).
- **Cofre ativo:** estado global em que existe um cofre carregado, autenticado e disponível para uso na sessão atual. O cofre ativo sempre assume um dos subestados canônicos `Cofre Salvo` ou `Cofre Modificado`.
- **Cofre em pesquisa:** estado transitório sobreposto ao `Cofre ativo`, no qual a interface exibe uma busca ativa, mostrando apenas os segredos que correspondem aos critérios informados. Durante a busca, o cofre preserva o subestado canônico corrente (`Cofre Salvo` ou `Cofre Modificado`). Enquanto a pesquisa estiver ativa, todas as ações ficam indisponíveis exceto: sair da aplicação, navegar pelo cofre e visualizar segredo. Para retomar as demais ações, o usuário deve confirmar a pesquisa (selecionando o elemento desejado, o que encerra implicitamente a pesquisa) ou cancelar a pesquisa. Em ambos os casos, o cofre retorna ao estado anterior ao início da pesquisa.

OBS:
 - Não existe um estado observável de cofre "bloqueado" separado, pois o bloqueio é tratado como um retorno ao fluxo de abrir o cofre novamente, exigindo nova autenticação e recarregando o estado salvo do arquivo, minimizando a retenção de dados sensíveis em memória.
 - Os estados "Abrindo cofre", "Criando novo cofre" e "Salvando cofre com outro caminho" são transitórios, pois exigem um tratamento especial para garantir usabilidade adequada na seleção de caminho, nome de arquivo e senha mestra, com uma UX devidamente projetada para cada um desses fluxos.

Temporariamente, durante os estados anteriores, a aplicação pode assumir um estado transitório, retornando ao último estado válido:
    - **Operação modal / confirmação bloqueante:** estado transitório sobreposto ao estado principal, usado para confirmações críticas, seleção de arquivos, formulários e ações destrutivas.
    - Apresentar tela de ajuda e comandos
    
#### Subestados do cofre ativo em memória

O cofre ativo possui dois subestados canônicos em memória:

- **Cofre Salvo:** cofre sincronizado com o arquivo corrente.
- **Cofre Modificado:** cofre com divergência entre o estado em memória e o último estado salvo.

OBS: 
- não existe estado observável "bloqueado", pois o bloqueio é tratado como abrir novamente o cofre, exigindo nova autenticação e recarregando o estado salvo do arquivo.
- não existe cofre "novo", pois o cofre é criado com a estrutura inicial e salvo imediatamente, entrando diretamente no estado "Cofre Salvo" desde o início.

#### Estados principais de segredos, pastas e modelos de segredo

- **Segredo disponível:** segredo visível na hierarquia principal e elegível para navegação, edição, movimentação e cópia.
- **Segredo ativo:** segredo atualmente selecionado, passível de ações como edição, movimentação e exclusão, etc. Normalmente, é o segredo que está sendo mostrado no momento.
- **Segredo favorito:** segredo disponível com marcação adicional de destaque visual e presença na pasta virtual de Favoritos.
- **Segredo em criação:** segredo ainda não confirmado pelo fluxo de criação; pode ser cancelado (descartado) sem efeito persistente.
- **Segredo novo**: segredo criado na sessão atual, confirmado mas não persistido. Ele poderá ser novamente editado em modo padrão ou avançado e continuará sendo considerado novo até o próximo salvamento, inclusive se sofrer novas alterações nesse intervalo.
- **Segredo em edição:** segredo disponível com alterações locais em andamento; pode ser cancelado (revertido) sem efeito persistente.
- **Segredo modificado:** segredo previamente persistido que sofreu alteração confirmada e ainda não foi salvo novamente. Novas edições confirmadas preservam esse mesmo estado até o próximo salvamento.
- **Segredo excluído reversivelmente:** segredo retirado da hierarquia principal e materializado apenas na Lixeira até o próximo salvamento. Enquanto permanecer nesse estado, não pode ser editado.
- **Segredo restaurado:** segredo anteriormente excluído reversivelmente e reinserido na hierarquia principal antes do próximo salvamento, retornando ao estado `Segredo novo`.
- **Pasta existente:** pasta presente na hierarquia, passível de renomeação, movimentação e exclusão física com promoção dos filhos.
- **Pasta ativa:** pasta atualmente selecionada para ações de edição, movimentação e exclusão. Se um segredo estiver ativo, então a pasta que o contém também é considerada implicitamente ativa.
- **Modelo disponível:** modelo existente e disponível para criação de novos segredos.
- **Modelo em edição:** modelo com alteração estrutural em andamento, afetando apenas criações futuras após confirmação.

#### Estados transitórios de exposição de dados sensíveis

- **Campo sensível oculto:** estado padrão de exibição para campos do tipo `texto sensível`.
- **Campo sensível exibido temporariamente:** estado temporário após ação explícita do usuário, encerrado manualmente ou por temporizador configurado.
- **Área de transferência povoada:** existe um valor copiado aguardando limpeza automática por temporizador ou por bloqueio/fechamento do cofre.

### Fluxos principais

**Abrir aplicação**
  - Ao iniciar, a aplicação mostra uma tela de welcome com ASCII art de apresentação do Abditum.
  - A tela inicial oferece as ações de criar cofre, abrir cofre, acessar ajuda e sair.
  - A partir dessa tela, a aplicação permanece no estado `Inicial / sem cofre ativo` até o usuário escolher a próxima ação.

**Criar novo cofre**
  - Usuário informa caminho e senha mestra com confirmação.
  - A aplicação popula a estrutura inicial do cofre com modelos e pastas padrão.
  - Se não existir arquivo no caminho informado, a aplicação grava diretamente o novo cofre no caminho final, usando o formato da versão atual.
  - Se já existir arquivo no caminho informado, a aplicação exige confirmação explícita de sobrescrita.
  - Se já existir um backup anterior com extensão `.abditum.bak`, a aplicação o renomeia temporariamente para `.abditum.bak2` antes de gerar o novo backup.
  - A aplicação gera um novo backup do arquivo existente com extensão `.abditum.bak` e então grava diretamente o novo cofre no caminho final.
  - Se a operação for concluída com sucesso, a aplicação remove o `.abditum.bak2`, preservando apenas o novo `.abditum.bak`.
  - Se a operação falhar antes da consolidação final, a aplicação restaura o `.abditum.bak2` para `.abditum.bak` sempre que possível.
  - Em caso de falha na gravação do novo arquivo após a geração do backup, a aplicação deve exibir uma mensagem de erro informando a falha e que existe um backup disponível para intervenção manual do usuário.
  - Esse fluxo não utiliza arquivo `.abditum.tmp`, pois não se trata do salvamento incremental de um cofre já aberto, e sim da criação de um novo arquivo de cofre.
  - O cofre entra em estado `Cofre Salvo`.

**Abrir cofre existente**
  - Usuário informa caminho.
  - A aplicação valida assinatura `magic` e `versão_formato`.
  - Seleciona o perfil Argon2id histórico a partir de `versão_formato`.
  - Usuário informa senha mestra.
  - Deriva a chave, valida o payload cifrado e carrega o domínio em memória.
  - Se o payload descriptografado estiver em um formato histórico suportado, a aplicação realiza a migração dos dados em memória para o modelo corrente do domínio.
  - O cofre entra em estado `Cofre Salvo`.

**Visualizar hierarquia do cofre**
  - O usuário navega pela árvore de pastas e segredos do cofre ativo.
  - A aplicação apresenta a hierarquia conforme a ordem persistida no JSON, mostrando primeiro segredos e depois subpastas em cada coleção.
  - O usuário pode expandir, colapsar e mover o foco entre os nós, preservando o contexto estrutural do cofre.
  - Esse fluxo não altera o conteúdo persistido do cofre nem o estado do domínio.
  
**Salvar cofre em estado `Cofre Modificado`**
  - A aplicação grava o cofre num caminho com sufixo ".abditum.tmp", usando o formato da versão atual, e atualiza a `versão_formato` do cabeçalho quando necessário, com `nonce` diferente.
  - Se já existir um backup anterior com extensão `.abditum.bak`, a aplicação o renomeia temporariamente para `.abditum.bak2` antes de gerar o novo backup.
  - Copia o arquivo atual do cofre para um novo backup com extensão `.abditum.bak`.
  - Depois renomeia o arquivo `.abditum.tmp` para o nome final do cofre, substituindo o arquivo original.
  - Se a operação for concluída com sucesso, a aplicação remove o `.abditum.bak2`, preservando apenas o novo `.abditum.bak`.
  - Se a operação falhar antes da consolidação final, a aplicação restaura o `.abditum.bak2` para `.abditum.bak` sempre que possível.
  - Em caso de falha na escrita ou substituição do arquivo final após a geração do backup, a aplicação deve exibir uma mensagem de erro informando a falha e que existe um backup disponível para intervenção manual do usuário.
  - O cofre entra em estado `Cofre Salvo`.

**Sair da aplicação**
  - O usuário pode encerrar a aplicação a qualquer momento.
    - No estado `Inicial / sem cofre ativo`, a aplicação encerra após solicitar confirmação do encerramento.
    - No estado `Cofre Salvo`, a aplicação encerra após solicitar confirmação do encerramento.
    - No estado `Cofre Modificado`, a aplicação oferece as opções de salvar, descartar alterações ou cancelar o encerramento, para evitar perda acidental de dados.
      - Em caso de salvar, a aplicação segue o fluxo de salvamento descrito anteriormente e encerra somente após salvamento bem-sucedido.
    - Em caso de um fluxo em andamento (ex: criação ou edição de segredo), a aplicação segue a mesma lógica de tratamento.

**Bloquear acesso ao cofre**
  - O bloqueio pode ser manual ou por inatividade.
  - A aplicação fecha logicamente o cofre, limpa buffers controlados sempre que possível e limpa a área de transferência.
  - A aplicação volta para o fluxo "Abrir cofre existente", assumindo o mesmo caminho do cofre previamente aberto, mas exigindo nova autenticação para desbloquear.

#### Fluxos complementares do cofre

**Descartar alterações não salvas e recarregar cofre**
  - O usuário inicia a ação de recarregar o cofre ativo.
  - Se o cofre estiver em estado `Cofre Salvo`, a aplicação apenas recarrega o arquivo atual em memória.
  - Se o cofre estiver em estado `Cofre Modificado`, a aplicação exige confirmação para descartar as alterações locais ainda não persistidas.
  - Após a confirmação, a aplicação reabre o arquivo atual, repete validação, descriptografia e eventual migração em memória.
  - Ao final, o cofre retorna ao estado `Cofre Salvo`.

**Salvar cofre em novo caminho**
  - O usuário inicia a ação de salvar o cofre em um novo caminho.
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

**Alterar senha mestra do cofre**
  - O usuário inicia a ação de alteração da senha mestra sobre o cofre ativo.
  - A aplicação solicita a senha mestra atual, a nova senha e a confirmação da nova senha.
  - Se a autenticação da senha atual e a confirmação da nova senha forem válidas, a aplicação rederiva a chave e prepara o cofre para ser persistido com a nova credencial.
  - A alteração da senha mestra não modifica o conteúdo lógico do domínio, mas exige regravação criptográfica completa do arquivo.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado` até o próximo salvamento bem-sucedido.

**Configurar o cofre**
  - O usuário inicia a edição das configurações do cofre ativo.
  - A aplicação permite alterar o tempo de bloqueio automático por inatividade, o tempo de reocultação de campos sensíveis e o tempo de limpeza automática da área de transferência.
  - As alterações passam a valer para o comportamento da sessão corrente conforme aplicável e permanecem associadas ao próprio cofre.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Exportar cofre para JSON plain text**
  - O usuário inicia a exportação do cofre ativo para formato JSON plain text.
  - Antes de exportar, a aplicação mostra aviso explícito sobre o risco de segurança de gerar uma cópia não criptografada e exige confirmação.
  - Após a confirmação, a aplicação serializa o domínio atual para JSON em texto claro no destino escolhido pelo usuário.
  - Esse fluxo não altera o conteúdo lógico do cofre ativo nem seu estado persistido.

**Importar cofre de JSON plain text**
  - O usuário inicia a importação de um arquivo JSON plain text para o cofre ativo.
  - A aplicação lê o conteúdo importado e resolve conflitos por identidade conforme as regras do cofre.
  - Pastas com a mesma identidade são mescladas silenciosamente.
  - Se um segredo importado colidir por identidade com um segredo já existente, a aplicação cria um novo segredo logicamente equivalente, com identidade diferente e preservando os demais dados importados.
  - Se um segredo importado colidir por nome com outro segredo já existente na mesma pasta de destino, a aplicação ajusta seu nome com sufixo numérico incremental e informa esse ajuste ao usuário.
  - Modelos com a mesma identidade são sobrepostos silenciosamente pelo modelo importado.
  - Após a confirmação e incorporação dos dados, o cofre entra em estado `Cofre Modificado`.




#### Fluxos principais de segredos

**Visualizar segredo**
  - O usuário navega pela hierarquia do cofre e seleciona um segredo disponível.
  - A aplicação exibe os detalhes do segredo selecionado no Painel do Segredo.
  - Campos do tipo `texto sensível` são apresentados ocultos por padrão, exigindo ação explícita para exibição temporária.
  - A visualização do segredo não altera seu conteúdo nem o estado persistido do cofre.
  - O segredo permanece disponível para outras ações, como edição, favoritar, movimentação ou cópia, sem restrições adicionais.

**Visualizar ou ocultar campo sensível**
  - O usuário seleciona um campo do tipo `texto sensível` no detalhe de um segredo já visível.
  - A aplicação permite revelar temporariamente o valor do campo mediante ação explícita do usuário.
  - O usuário pode ocultar novamente o valor manualmente a qualquer momento.
  - Se o usuário não ocultar o valor manualmente, a aplicação o reoculta automaticamente conforme a configuração do cofre.
  - Esse fluxo não altera o conteúdo persistido do segredo nem o estado do cofre.

**Copiar campo de segredo**
  - O usuário seleciona qualquer campo de um segredo visível, inclusive campos do tipo `texto sensível`, e inicia a ação de cópia.
  - A aplicação copia o valor atual do campo para a área de transferência do sistema.
  - A aplicação exibe feedback visual de cópia e inicia o temporizador de limpeza automática conforme a configuração do cofre.
  - O conteúdo copiado também é limpo ao bloquear ou fechar o cofre.
  - Esse fluxo não altera o conteúdo persistido do segredo nem o estado do cofre.

**Criar segredo**
  - O usuário inicia a criação de um novo segredo na raiz ou dentro de uma pasta.
  - A aplicação oferece a escolha entre usar um modelo de segredo existente ou começar com um segredo vazio, sem nenhuma estrutura inicial.
  - Caso o usuário opte por um modelo de segredo, a estrutura inicial do segredo é gerada a partir do modelo escolhido, copiando os campos como snapshot, sem manter vínculo por referência com o modelo de origem.
  - Caso o usuário opte por começar com um segredo vazio, é gerado um segredo sem campos adicionais além do nome e da observação, e os demais campos poderão ser adicionados posteriormente pela edição avançada.
  - Após a confirmação, o novo segredo assume estado `Novo` e é inserido no destino selecionado, e o cofre entra em estado `Cofre Modificado`.
  - Caso o usuário tenha optado por um modelo de segredo, então a aplicação passa para o fluxo de edição padrão.
  - Caso o usuário tenha optado por um segredo vazio, então a aplicação passa para o fluxo de edição avançada, para que o usuário possa adicionar os campos desejados.

**Duplicar segredo**
  - O usuário seleciona um segredo existente e inicia a ação de duplicação.
  - A aplicação cria uma nova instância com nova identidade, copiando nome, observação, favorito e campos do segredo original.
    - O nome do segredo duplicado recebe um sufixo numérico incremental para evitar confusão com o segredo original. Ex: "Segredo" → "Segredo (1)", "Segredo (2)", etc.
  - Após a confirmação, o segredo duplicado assume estado `Novo` e é inserido na mesma coleção do segredo de origem, e o cofre entra em estado `Cofre Modificado`.

**Favoritar segredo**
  - O usuário seleciona ou visualiza um segredo disponível e não favoritado, e alterna seu marcador de favorito.
  - A aplicação altera apenas o atributo `favorito`, sem modificar identidade, conteúdo ou localização do segredo.
  - A presença do segredo na pasta virtual de Favoritos passa a refletir imediatamente esse estado.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Desfavoritar segredo**
  - O usuário seleciona ou visualiza um segredo disponível e favoritado, e alterna seu marcador de favorito.
  - A aplicação altera apenas o atributo `favorito`, sem modificar identidade, conteúdo ou localização do segredo.
  - A presença do segredo na pasta virtual de Favoritos passa a refletir imediatamente esse estado.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Editar segredo (edição padrão)**
  - O usuário seleciona ou visualiza um segredo existente e inicia a edição.
  - A aplicação abre o segredo no modo de edição padrão.
  - A aplicação permite alterar nome, observação e valores dos campos existentes.
  - A identidade do segredo é preservada durante toda a edição.
  - O usuário pode alternar para a edição avançada caso precise alterar a estrutura do segredo.
  - Após a confirmação, o segredo preserva seu estado anterior se já estiver em `Segredo novo` ou `Segredo modificado`.
  - Após a confirmação, se o segredo estava em `Segredo disponível`, ele passa para `Segredo modificado`.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Editar segredo (edição avançada)**
  - O usuário seleciona ou visualiza um segredo existente e inicia a edição avançada.
  - A aplicação abre o segredo no modo de edição avançada.
  - Nesse modo, o usuário altera apenas a estrutura do segredo.
  - Não é permitido alterar o tipo de um campo existente. Para isso, é necessário excluir o campo e criar um novo com o tipo desejado.
  - A identidade do segredo é preservada durante toda a edição.
  - O usuário pode alternar para a edição padrão quando quiser voltar a alterar os valores dos campos.
  - Após a confirmação, o segredo preserva seu estado anterior se já estiver em `Segredo novo` ou `Segredo modificado`.
  - Após a confirmação, se o segredo estava em `Segredo disponível`, ele passa para `Segredo modificado`.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Excluir segredo reversivelmente**
  - O usuário seleciona ou visualiza um segredo disponível e inicia a ação de remoção.
  - A aplicação exige confirmação, mas a remoção é reversível até o próximo salvamento do cofre.
  - A identidade e o conteúdo do segredo permanecem preservados.
  - Enquanto o segredo permanecer na Lixeira, a aplicação não permite edição desse segredo.
  - O cofre entra em estado `Cofre Modificado`, a aplicação retira o segredo da hierarquia principal e o materializa na pasta virtual Lixeira.

**Restaurar segredo excluído reversivelmente**
  - O usuário seleciona um segredo presente na Lixeira e inicia a ação de restauração.
  - A aplicação reinsere o segredo na hierarquia principal antes do próximo salvamento.
  - A identidade e o conteúdo do segredo são preservados durante a restauração.
  - Após a restauração, o segredo retorna ao estado `Segredo novo`.
  - O cofre entra em estado `Cofre Modificado`.

**Mover segredo**
  - O usuário seleciona um segredo existente e inicia a ação de movimentação.
  - A aplicação coleta o novo destino, que pode ser outra pasta ou a raiz do cofre.
  - O segredo é removido da coleção atual e reinserido na coleção de destino, preservando identidade, conteúdo e marcação de favorito.
  - A identidade e o conteúdo do segredo são preservados durante toda a movimentação. O estado do segredo permanece inalterado, mas sua posição na hierarquia é atualizada para refletir o novo destino.
  - O segredo é adicionado ao final da lista de segredos do destino.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Reordenar segredo**
  - O usuário seleciona um segredo existente e inicia a ação de reordenação relativa.
  - A aplicação altera apenas sua posição entre os segredos irmãos da mesma coleção pai.
  - A nova ordem passa a refletir diretamente a ordem persistida e a ordem de exibição.
  - A identidade e o conteúdo do segredo são preservados durante toda a movimentação. O estado do segredo permanece inalterado, mas sua posição na hierarquia é atualizada para refletir o novo destino.
  - Após a confirmação, o cofre entra em estado `Cofre Modificado`.

**Buscar segredos**
  - O usuário inicia o processo de busca.
  - A aplicação executa uma varredura sequencial em memória sobre nome do segredo, nome de campo, valores de campos do tipo `texto` e observação.
  - A hierarquia é reapresentada mostrando apenas os segredos que satisfazem o critério de busca, mas mantendo a estrutura de pastas para preservar o contexto de localização dos segredos encontrados.
  - Enquanto a pesquisa estiver ativa, todas as ações ficam indisponíveis exceto: sair da aplicação, navegar pelo cofre e visualizar segredo.
  - Quando o casamento ocorrer no nome do segredo, o segredo correspondente pode receber destaque visual na árvore.
  - O usuário confirma a pesquisa selecionando o elemento desejado (o que encerra implicitamente a pesquisa) ou cancela a pesquisa. Em ambos os casos, o cofre retorna ao estado anterior ao início da pesquisa.

#### Fluxos principais de pastas

**Criar pasta**
  - O usuário inicia a criação de uma nova pasta a partir da raiz ou de uma pasta existente.
  - A aplicação coleta o nome da pasta e determina o destino conforme o contexto atual.
  - Se o destino for a raiz, a nova pasta é adicionada ao final da lista de pastas da raiz.
  - Se o destino for uma pasta existente, a nova pasta é adicionada ao final da lista de subpastas dessa pasta.
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





## Decisões técnicas

- Linguagem: Go, compilado como binário único executável
- Interface: TUI moderna (Bubble Tea / teatest/v2)
- Criptografia: AES-256-GCM + Argon2id
- Identidade dos elementos: usar NanoID de 6 caracteres alfanuméricos para representar a identidade persistida de segredos, pastas e modelos de segredo. A escolha busca combinar baixa probabilidade prática de colisão, independência em relação ao nome do elemento, facilidade de serialização no JSON e portabilidade entre importação, exportação, movimentação e migração de formato.
 - Parametrização operacional do Argon2id:
  - Objetivo de UX e segurança: manter o desbloqueio e a abertura do cofre em uma faixa interativa de aproximadamente 0,8 s a 1,5 s em hardware compatível, sem reduzir agressivamente o custo de memória
  - Política de parametrização: em v1, os parâmetros do Argon2id são fixos e hard-coded na aplicação, iguais para todos os cofres suportados por uma mesma versão. Não há calibração por máquina nem parametrização variável por arquivo
  - Custo de memória da v1: 256 MiB por derivação
  - Piso de segurança para revisões futuras: 128 MiB por derivação; abaixo disso a aplicação não deve operar em v1
  - Teto operacional de referência: 512 MiB por derivação, para evitar degradação excessiva em máquinas com menos recursos
  - Custo de tempo da v1: no mínimo 3 iterações, definido de forma fixa pela aplicação
  - Paralelismo da v1: até 4 threads lógicas, limitado pela quantidade disponível na máquina
  - Política por plataforma: Windows, macOS e Linux 64 bits seguem a mesma política base de segurança, sem parametrização variável por arquivo
  - Evolução controlada: alterações nesses parâmetros só podem ocorrer como decisão explícita de versão da aplicação, acompanhadas por validação de compatibilidade e testes de regressão
    - Seleção histórica de perfil: ao abrir um cofre, a aplicação deve escolher o perfil Argon2id histórico exclusivamente a partir de `versão_formato` registrada no cabeçalho do arquivo
 - CI obrigatório
- Hardening do arquivo de cofre:
 - Verificação de integridade e autenticidade do payload cifrado usando o mecanismo nativo do AES-256-GCM, sem checksum adicional implementado pela aplicação
 - O arquivo deve começar com uma assinatura mágica fixa do formato, para identificação imediata de arquivos Abditum e rejeição precoce de entradas inválidas
 - O cabeçalho de criptografia (`magic`, `versão_formato`, `salt`, `nonce`) deve ser autenticado como AAD do AES-256-GCM
 - Salt único por cofre — já implícito no Argon2id, mas vale documentar explicitamente que o salt é gerado aleatoriamente na criação e armazenado no cabeçalho do arquivo
 - Nonce único por salvamento — o nonce do AES-GCM deve ser regenerado a cada operação de salvamento
 - O Argon2id é usado exclusivamente para derivação de chave a partir da senha mestra; ele não é responsável pela integridade do conteúdo cifrado
 - Histórico de versão do formato de arquivo — um campo de versão no cabeçalho do arquivo para permitir migrações futuras do formato
- Compatibilidade e migração de formato:
 - A leitura do JSON descriptografado deve compreender todos os formatos históricos de payload, da versão 0 até a versão atual suportada pela aplicação
 - A derivação da chave para abertura de cofres históricos deve usar o perfil Argon2id correspondente à `versão_formato` do arquivo, conforme a tabela hard-coded da versão da aplicação
 - A migração de formato ocorre ao abrir o cofre: o conteúdo antigo é carregado e convertido em memória para o modelo corrente do domínio
 - O salvamento sempre escreve o arquivo no formato mais recente da aplicação, atualizando a `versão_formato` do cabeçalho quando necessário
 - Mudanças de formato devem ser raras, pontuais e acompanhadas por rotinas explícitas de migração e testes de regressão para versões anteriores suportadas
 - Arquivos com `versão_formato` superior à versão suportada pela aplicação devem falhar com erro claro de incompatibilidade de versão
- Codificação do arquivo criptografado — UTF-8 explícito para suporte a caracteres especiais em nomes e valores
- Portabilidade Extrema: Modelos de segredo e configurações da aplicação (ex: tempo de bloqueio) são armazenados internamente no próprio arquivo do cofre. Essa decisão arquitetural garante que cada arquivo seja 100% autossuficiente e portátil, dispensando o uso de arquivos de configuração externos no sistema operacional.

## Estratégia de Implementação
- Use DDD
- Permita que as entidades sejam navegadas somente leitura. As operações de modificação devem ser realizadas exclusivamente por meio de métodos explícitos de um Manager, que estabelecerá um contrato (uma API) para manipulação segura e consistente das estruturas de dados. Isso garante que as regras de negócio sejam centralizadas e aplicadas de forma uniforme, evitando manipulações diretas e potencialmente inseguras das entidades.
- Seja generoso com comentários. Assuma que o código pode ser lido por pessoas com menos familiaridade em GO, especialmente nas libs especializadas de criptografia e TUI. Sei que a comunidade GO é acostumada a ler código sem muitos comentários, mas aqui o objetivo é criar um projeto didático e acessível, então vale a pena explicar os conceitos e decisões de forma clara no código.

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
  - **Painel da Hierarquia (Esquerdo):** Dedicado à navegação na hierarquia do cofre.
  - **Painel do Segredo (Direito):** Dedicado à visualização, criação e edição do segredo selecionado.
- **Comportamento Responsivo:**
  - `Menos de 5 linhas` ou `Menos de 20 colunas`: Oculta os painéis e exibe apenas uma mensagem pedindo para redimensionar o terminal.
  - `Menos de 40 colunas`: Exibe apenas o Painel do Segredo (foco em dispositivos restritos).
  - `Mais de 40 colunas`: Exibe layout completo com Painel da Hierarquia e Painel do Segredo.

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
  - Segredos criados a partir de um modelo tendem a iniciar na edição padrão; segredos criados vazios tendem a iniciar na edição avançada.
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
- **Criação/Alteração de Senha Mestra:** Ao definir a senha mestra pela primeira vez, o fluxo exige digitação dupla para prevenir erros. Ao alterar uma senha existente, exige digitar a senha atual antes da confirmação da nova.

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