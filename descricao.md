# Abditum - Cofre de Senhas Portátil e Seguro

## O que é

Abditum é um cofre de senhas portátil, seguro e fácil de usar, com uma interface TUI moderna. 
Ele permite que os usuários armazenem e gerenciem suas senhas e informações confidenciais de forma organizada e protegida, sem depender de serviços em nuvem ou instalações complexas.

## Diferenciais

O cofre deve ser completamente portátil e seguro — um único arquivo executável que qualquer pessoa pode copiar e usar discretamente em qualquer lugar, sem deixar rastros no sistema. 

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
- Bloquear cofre manualmente (voltando a tela de autenticação)
- Bloquear automaticamente após inatividade
- Exportar cofre para formato JSON plain text
- Importar cofre de formato JSON plain text
 - Caso o cofre atual já possua o segredo na mesma localização, ao invés de substituir o segredo existente, importar com um sufixo numérico incremental para evitar sobrescrever o segredo existente. Ex: "Segredo" → "Segredo (1)", "Segredo (2)", etc.
- Configurar o cofre
  - Tempo para bloqueio automático por inatividade
  - Tempo para ocultar valor de campo de segredo após exibição temporária
  - Tempo para limpar a área de transferência automaticamente

### Autenticação

- Criar novo cofre com senha mestra
- Desbloquear cofre com senha mestra
- Bloquear automaticamente após inatividade
- Proteção contra brute force e ataques offline — garantida através de parâmetros rígidos do Argon2id (custo alto de memória e tempo) para atrasar tentativas de derivação de chave


### Navegação da Hierarquia do Cofre (somente leitura)
- Exibir hierarquia do cofre (pastas e segredos)
- Exibir detalhes do segredo selecionado
- Exibir temporariamente o valor de um campo de segredo (ex: senha) com opção de ocultar/mostrar
  - Ocultar automaticamente o valor do campo de segredo após 15 segundos

### Gerenciamento de Segredos
- Criar um segredo
  - Como:
    - Usando um modelo de segredo pré-definido
    - Usando um modelo de segredo personalizado
    - O segredo gerado não mantém vínculo por referência com o modelo. O nome do modelo é guardado apenas como histórico ("snapshot" do momento da criação). Alterar a estrutura (campos), renomear ou excluir o modelo posteriormente não afeta os segredos já criados.
  - Onde:
    - Na raiz do cofre ou dentro de uma pasta
    - Numa pasta da hierarquia do cofre
- Duplicar um segredo existente
- Favoritar/desfavoritar um segredo
- Editar um segredo existente
    - Alterar campos de segredo (altera apenas o valor)
    - Avançado:
        - Adicionar um novo campo de segredo
            - Informar o nome do campo de segredo
            - Informar o tipo do campo de segredo
            - Informar o valor do campo de segredo
        - Alterar o tipo do campo de segredo
        - Excluir um campo de segredo
        - Alterar a posição (reordenar) de um campo de segredo
- Remover um segredo de forma reversível até o próximo salvamento do cofre
- Restaurar um segredo removido antes do próximo salvamento do cofre
- Mover um segredo para outra pasta ou para a raiz do cofre
- Mover (reordenar) segredo relativamente a outros segredos dentro da mesma pasta ou raiz do cofre
- Buscar segredos por nome, por valor de um campo (dado não sensível) ou observação, ou pelo nome da pasta onde o segredo está localizado

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
 - Excluir pasta, movendo seus segredos e suas subpastas filhas para a pasta pai (ou para a raiz do cofre, se a pasta excluída estiver na raiz do cofre)
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
- Criar modelo de segredo a partir de um segredo existente, copiando os campos e valores para o novo modelo
- Modelos de segredo pré-definidos criados no cofre por padrão, mas editáveis e removíveis pelo usuário
    - Modelo de segredo "Login": campos "URL", "Username" e "Password"
    - Modelo de segredo "Cartão de Crédito": campos "Número do Cartão", "Nome no Cartão", "Data de Validade" e "CVV"
    - Modelo de segredo "API Key": campos "Nome da API", "Chave de API"

### Área de Transferência

- Copiar qualquer campo para a área de transferência
- Limpar automaticamente a área de transferência após 30 segundos
- Limpar automaticamente a área de transferência ao bloquear ou fechar o cofre

### Segurança
- Zerar da memória todos os dados sensíveis ao bloquear ou fechar
- Proteção contra shoulder surfing — atalho para ocultar toda a interface rapidamente
- Ao exportar o cofre para formato JSON plain text, mostrar mensagem de aviso sobre os riscos de segurança envolvidos e pedir confirmação do usuário antes de prosseguir com a exportação

## Requisitos não Funcionais
- Criptografia: AES-256-GCM para criptografia dos dados e Argon2id para derivação de chave a partir da senha mestra
- Armazenamento em formato JSON criptografado
- Compatibilidade com sistemas operacionais Windows, macOS e Linux
- Formato do arquivo de cofre: extensão .abditum
- Confiabilidade (Salvamento Atômico): ao salvar, escrever os dados primeiro em um arquivo `.abditum.tmp` no mesmo diretório do cofre. Somente após a gravação bem-sucedida, renomear o arquivo, substituindo o `.abditum` original. Em caso de falha, o arquivo `.tmp` deve ser imediatamente apagado para evitar persistência indevida de dados fora do arquivo final do cofre.
- Ao salvar, manter uma cópia de backup do cofre anterior com extensão .abditum.bak
- Executável "portable"
    - Não requer instalação, pode ser executado diretamente do arquivo binário
    - Permitir copiar o arquivo binário para qualquer local do sistema de arquivos e executá-lo
    - Não usar arquivos de configuração ou dados fora do arquivo do cofre
  - Não persistir dados fora do arquivo do cofre, exceto artefatos transitórios e backups explicitamente previstos pela própria aplicação
- Os modelos de segredo são armazenados dentro do cofre, permitindo que cada cofre tenha seus próprios modelos personalizados
- Os modelos de segredo pré-definidos são fornecidos populados no cofre ao criar um novo cofre, mas podem ser editados ou removidos pelo usuário
- O cofre é criado com uma hierarquia pré-definida de pastas
- Privacidade: Ausência total de logs de aplicação (stdout/stderr) que contenham caminhos de arquivos de cofre, nomes de segredos ou valores de campos.
- A aplicação de versão N deverá ser capaz de abrir arquivos de cofre criados com a versão N-1, garantindo compatibilidade retroativa

## Requisitos Inversos
- Armazenamento na nuvem
- Múltiplos cofres abertos simultaneamente
- App mobile ou web — TUI portátil por design
- Tags — pastas/grupos suficientes para v1
- Histórico de versões de segredos

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
    - Cofre bloqueado: estado de proteção em que o aplicativo continua em execução, mas os dados sensíveis foram limpos da memória, exigindo nova autenticação
- Segredo: item individual dentro do cofre, com dados comuns e dados sensíveis
    - Segredo favorito: um segredo marcado pelo usuário como prioritário ou de uso frequente, ganhando destaque para acesso rápido
- Dados: informações armazenadas em um segredo
    - Dados comuns: informações não sensíveis, como nome do serviço ou URL
    - Dados sensíveis: informações confidenciais, como senha, notas privadas, apikeys
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
- Exclusão reversível (Soft Delete): mecanismo pelo qual um segredo removido permanece restaurável até o próximo salvamento definitivo do cofre
- Auditoria de senhas (Saúde do Cofre): processo de análise de segurança para alertar o usuário sobre a presença de senhas fracas, antigas ou reutilizadas
- TOTP (Autenticação de Dois Fatores): código numérico temporário gerado em tempo real pelo cofre a partir de uma chave secreta, servindo como reforço de segurança para acesso a serviços externos
- Shoulder Surfing: técnica de espionagem física mitigada pela aplicação, onde um indivíduo mal-intencionado observa a tela do usuário para roubar informações visíveis

## Modelagem

Decisões de modelagem:
 - **Hierarquia recursiva:** A raiz do cofre funciona como uma pasta sem nome. Pastas podem conter segredos e subpastas em qualquer nível de aninhamento.
 - **Ordenação por posição:** A ordem dos elementos no JSON reflete diretamente a ordem de exibição na interface (segredos, pastas e campos).
 - **Modelo como snapshot:** Segredos criados a partir de modelos não mantêm vínculo por referência — o nome do modelo é guardado apenas como registro histórico. Não há distinção estrutural entre segredos criados com ou sem modelo.
 - **IDs (NanoID, 6 caracteres alfanuméricos):** Segredos, pastas e modelos de segredo possuem ID. O nome do segredo, da pasta e do modelo é apenas um atributo editável, não seu identificador. O espaço de 62⁶ (~56 bilhões) combinações garante unicidade prática.
 - **Campos uniformes:** O valor de um campo pode ser string vazia (campo existente, não preenchido). Não há distinção de estado entre preenchido e vazio.
 - **Observação implícita:** Todo segredo possui um campo de observação (opcional, texto livre) que não é declarado nos modelos e não pode ser removido.
 - **Configurações embutidas:** As configurações do cofre são armazenadas dentro do próprio arquivo, garantindo portabilidade total sem arquivos externos.

 - Cofre (Estrutura do Payload JSON Criptografado):
    - configurações:
      - tempo_bloqueio_inatividade_minutos: inteiro
      - tempo_ocultar_segredo_segundos: inteiro
      - tempo_limpar_area_transferencia_segundos: inteiro
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
- id: nanoid de 6 digitos alfanuméricos gerados aleatoriamente, para garantir unicidade e evitar conflitos de nomes

## Formato do Arquivo

O arquivo do cofre (`.abditum`) é um stream binário composto por duas partes:

1.  **Cabeçalho de Criptografia (fixo):**
    *   `versão_formato` (inteiro): Versão do formato do arquivo.
    *   `salt` (bytes): Salt para a derivação da chave Argon2id.
    *   `nonce` (bytes): Vetor de inicialização para a criptografia AES-256-GCM.
2.  **Payload Criptografado:**
    *   O restante do stream contém a estrutura do `Cofre` em formato JSON, criptografada com AES-256-GCM.

Essa estrutura garante que os metadados necessários para a descriptografia estejam disponíveis antes de ler o conteúdo do cofre, que permanece totalmente criptografado.

## Decisões técnicas

- Linguagem: Go, compilado como binário único executável
- Interface: TUI moderna (Bubble Tea / teatest/v2)
- Criptografia: AES-256-GCM + Argon2id
- Testes: 
 - testes do serviço de criptografia, incluindo casos de sucesso e falha para criptografia e descriptografia;
 - testes do serviço de armazenamento, incluindo casos de sucesso e falha para salvar e carregar o cofre;
 - testes unitários white-box para navegação e transições de estado; 
 - golden files visuais em 80×24 por tela; 
 - testes de comandos para cada tela e fluxo de usuário;
 - testes de integração para fluxo completo realizando todas as operações principais (criar cofre, criar segredo, editar segredo, etc.)
 - CI obrigatório
- Hardening do arquivo de cofre:
 - Verificação de integridade do arquivo
 - Salt único por cofre — já implícito no Argon2id, mas vale documentar explicitamente que o salt é gerado aleatoriamente na criação e armazenado no cabeçalho do arquivo
 - Nonce único por salvamento — o nonce do AES-GCM deve ser regenerado a cada operação de salvamento
 - Histórico de versão do formato de arquivo — um campo de versão no cabeçalho do arquivo para permitir migrações futuras do formato
- Encoding do arquivo criptografado — UTF-8 explícito para suporte a caracteres especiais em nomes e valores
- Portabilidade Extrema: Modelos de segredo e configurações da aplicação (ex: tempo de bloqueio) são armazenados internamente no próprio arquivo do cofre. Essa decisão arquitetural garante que cada arquivo seja 100% autossuficiente e portátil, dispensando o uso de arquivos de configuração externos no sistema operacional.

# Decisões de implementação
- Use DDD
- Permita que as entidades sejam navegadas somente leitura. As operações de modificação devem ser realizadas exclusivamente por meio de métodos explícitos de um Manager, que estabelecerá um contrato (uma API) para manipulação segura e consistente das estruturas de dados. Isso garante que as regras de negócio sejam centralizadas e aplicadas de forma uniforme, evitando manipulações diretas e potencialmente inseguras das entidades.
- Seja generoso com comentários. Assuma que o código pode ser lido por pessoas com menos familiaridade em GO, especialmente nas libs especializadas de criptografia e TUI. Sei que a comunidade GO é acostumada a ler código sem muitos comentários, mas aqui o objetivo é criar um projeto didático e acessível, então vale a pena explicar os conceitos e decisões de forma clara no código.

# Interface e Experiência do Usuário (TUI e UX)

## 1. Princípios de Design Visual e Interação
- **Interface TUI Moderna:** Ocupa todo o terminal, interativa e desenhada em modo texto com suporte a 256 cores, garantindo alto contraste e legibilidade (estética Cyberpunk sugerida).
- **Acessibilidade de Controles:** Suporte integral à navegação por teclado e suporte complementar a mouse (cliques em campos para edição, nós da árvore para seleção e botões para ações).
- **Menu de Ajuda de Contexto:** Tela ou barra de ajuda sempre acessível, mostrando as ações disponíveis e os respectivos atalhos para o contexto atual (ocultando ações não aplicáveis).
- **Status Global:** Manter visíveis, em uma barra de status: caminho do arquivo atual, status do cofre (novo, salvo, modificado) e total de segredos. Apenas um cofre pode ser aberto por vez.
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
- Um painel não ativo pode reagir ao contexto do painel ativo. Ex: ao navegar na árvore, o painel de segredos exibe o segredo com foco.
- As ações podem ser:
  - globais (disponíveis em todos os contextos, dependem do estado geral da aplicação). 
  Exemplo: "Salvar cofre" só é aplicável se houver um cofre aberto e alterações não salvas. "Sair" é sempre aplicável.
  - locais (disponíveis apenas no painel ativo). 
  - foco (disponíveis apenas no elemento focado dentro do painel ativo). Exemplo: "Excluir" é aplicável apenas se o elemento focado for um segredo ou pasta.
  - navegação (setas, tab, etc.) é sempre aplicável, mas o comportamento específico depende do foco e do painel ativo. Exemplo: seta direita expande pasta na árvore, mas não tem efeito no painel de segredos.
- Uma mesma ação podem ocorrer em paineis ou campos diferentes, se o contexto for adequado. Exemplo: "Favoritar" é aplicável tanto para segredos focados na árvore quanto para o segredo exibido no painel de segredos.


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
- **Alerta/Aviso** (não bloqueante, auto-oculta com tempo estendido):
  - Descrever a situação e recomendar ação, sem alarmismo.
  - Ex: "O cofre será bloqueado em 30 segundos por inatividade."
- **Informação** (não bloqueante, auto-oculta):
  - Fornecer instrução ou dado relevante para o contexto atual.
  - Ex: "Informe os dados do novo segredo."

## 2. Layout e Responsividade
- A interface principal é dividida em dois painéis:
  - **Painel Lateral (Esquerdo):** Dedicado à navegação na hierarquia do cofre.
  - **Painel Principal (Direito):** Dedicado à visualização, criação e edição do segredo selecionado.
- **Comportamento Responsivo:**
  - `Menos de 5 linhas` ou `Menos de 20 colunas`: Oculta os painéis e exibe apenas uma mensagem pedindo para redimensionar o terminal.
  - `Menos de 40 colunas`: Exibe apenas o painel principal (foco em dispositivos restritos).
  - `Mais de 40 colunas`: Exibe layout completo com painel lateral e principal.

## 3. Navegação e Hierarquia (Painel Lateral)
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
  - **Lixeira (Soft Delete):** Visível no final da raiz (apenas se houver excluídos). Lista itens removidos e permite restauração. É esvaziada irreversivelmente ao salvar o cofre.
- **Mecanismo de Busca:**
  - Ao buscar, a árvore é filtrada ocultando o que não corresponde.
  - O trecho exato pesquisado recebe *highlight* (destaque de cor) dentro do nome do segredo encontrado na árvore.

## 4. Visualização e Edição (Painel Principal)
- Exibe o detalhe do item focado no Painel Lateral. Se não houver foco, exibe *placeholder* informativo.
- A navegação com `Tab` alterna o controle entre o Painel Lateral e o Principal.
- **Privacidade Padrão:** Os campos do tipo "senha" (ou texto sensível) são carregados ocultos (ex: `****`), necessitando de ação explícita (toggle) para exibir o valor.
- **Gerenciamento de Espaço:**
  - Se um modelo tiver muitos campos, o painel principal deve permitir scroll vertical.
  - O campo de "Observação" é redimensionável e deve ocupar automaticamente o espaço livre restante no painel.
- **Área de Transferência Integrada:**
  - Ao copiar um campo, disparar um *toast* de sucesso.
  - Exibir um indicador visual na interface com um *countdown* de 30 segundos, após o qual a área de transferência será limpa por segurança.
- **Seleção de Arquivos (File Picker):**
  - Para criar, abrir ou salvar cofres com novos nomes, usar uma janela de File Picker integrada à TUI (não exigir digitação cega de caminhos), suportando navegação por setas, mouse e autocompletar.

## 5. Feedback e Segurança (Prevenção de Erros)
- **Aviso Fundamental:** Na criação de um cofre, alertar categoricamente sobre a Irrecuperabilidade ("Zero Knowledge"): o esquecimento da Senha Mestra resulta em perda total dos dados.
- **Ações Destrutivas Irreversíveis:** Excluir pastas (mesmo vazias), excluir Modelos de Segredo e alterar a Senha Mestra exigem um pop-up de confirmação explícita.
- **Soft Delete vs Hard Delete:** Ao excluir um segredo, exibir um aviso discreto de que ele foi movido para a Lixeira.
- **Proteção de Arquivos:** 
  - Alertar imediatamente se o arquivo do cofre for bloqueado por outro processo (Lock File).
  - Ao salvar sobre um cofre existente (Salvar Como / Criar Novo), confirmar sobrescrita e garantir a geração de um arquivo `.bak`.
- **Sair da Aplicação:**
  - Pode ser disparado a qualquer momento, mas se houver alterações não salvas, a aplicação deve ser retida por um menu perguntando: "Salvar", "Sair sem Salvar (Descartar)" ou "Cancelar".
- **Bloqueio por Inatividade:**
  - Tempo configurável pelo usuário (padrão 5 min).
  - Ao atingir o limite, emitir alerta de "Bloqueio Iminente". Qualquer interação (teclado/mouse) aborta o bloqueio e reseta o cronômetro. Se ignorado, o cofre é bloqueado.

## 6. Fluxos de Usuário Específicos
- **Criação/Alteração de Senha Mestra:** Ao definir a senha mestra pela primeira vez, o fluxo exige digitação dupla para prevenir erros. Ao alterar uma senha existente, exige digitar a senha atual antes da confirmação da nova.

## 7. Mapa de Teclas de Comando por Contexto

Ctrl C não deve interromper abruptamente a aplicação, mas funcionar 
como um comando convencional.

**Global / Qualquer contexto:**
 - Ctrl+Q: (Sair) 
    - Confirmar sair sem modificações no cofre: novamente Ctrl+Q (Sair)
    - Confirmar sair com modificações no cofre: Ctrl+S (Salvar)
    - Confirmar sair sem salvar: Ctrl+D (Descartar)
    - Cancelar ação de sair: Esc (Voltar)

**Foco no Painel Lateral (Hierarquia):**
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
    - Em segredo: visualiza e move o foco da aplicação para o painel principal.
  - Ctrl+N: Abrir tela de "Novo Segredo" no painel principal, visando a localização do nó atual.
  - Ctrl+F: Ativar barra de filtragem da árvore.

**Foco no Painel Principal (Exibindo Detalhe do Segredo):**
  - Setas para cima/baixo: mover foco entre os campos preenchidos.
  - Esc: devolver o foco da interface para o painel lateral (retornar à árvore).

**Foco no Painel Principal (Editando um Segredo):**
  - *(Área aberta para definição de atalhos de reordenação de campos e submissão de formulário)*