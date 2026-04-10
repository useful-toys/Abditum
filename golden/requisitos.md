# Abditum - Cofre de Senhas Portátil e Seguro

## O que é
Abditum é um cofre de senhas portátil para armazenar e organizar credenciais e informações confidenciais em um único arquivo local.
 
## Diferenciais
 - Uso discreto, adequado a diferentes contextos.
 - Controle e propriedade do cofre permanecem com o usuário.
 - Organização flexível das informações, com modelos padrão e personalizados.
 
## Conceitos fundamentais
O cofre reúne pastas, segredos e campos.
Cada segredo representa uma credencial, como acesso a um site, cartão de crédito ou chave de API.
Campos comuns ficam visíveis; campos sensíveis permanecem ocultos por padrão.

## Como protegemos seus dados
 
A proteção dos dados segue princípios claros:

 - **Acesso restrito**: sem a senha mestra, o conteúdo do cofre não pode ser acessado.
 - **Privacidade local**: os dados permanecem no arquivo do cofre, sem dependência de serviços online.
 - **Redução de exposição**: campos sensíveis ficam ocultos por padrão e dados de sessão são removidos ao bloquear ou encerrar a aplicação.
 - **Implementação enxuta**: o projeto adota um conjunto reduzido de dependências externas para diminuir riscos.

## Conceitos (Glossário)

- **Senha mestra**: chave de acesso ao cofre, usada para criptografar e descriptografar os dados
- **Senha falsa de coação** *(Duress Password)*: senha mestra alternativa que abre uma versão restrita do cofre, protegendo os dados reais em situações de ameaça *(fora de escopo v1)*
- **Cofre**: arquivo criptografado que armazena os segredos do usuário
  - **Bloqueio do cofre**: interrupção do acesso ao conteúdo do cofre, exigindo nova autenticação para retomar
- **Segredo**: item individual dentro do cofre, composto por campos
  - **Segredo favorito**: segredo marcado pelo usuário como prioritário, com destaque para acesso rápido
- **Campo**: elemento individual de um segredo, com nome e valor. Existem dois tipos de campo:
  - **Campo comum**: campo com valor sempre visível (como nome do serviço ou usuário)
  - **Campo sensível**: campo com valor oculto por padrão (como senha ou chave de API)
- **Observação**: campo comum especial que existe automaticamente em todo segredo; não pode ser renomeado, excluído ou movido
- **Pasta**: estrutura que agrupa segredos e outras pastas dentro do cofre
- **Modelo de segredo**: estrutura predefinida de campos para agilizar a criação de segredos

## Requisitos Funcionais

*Nota sobre o formato: Para manter a descrição compacta, esta seção utiliza uma representação simplificada. Cada item principal (`-`) é um requisito funcional, enquanto os sub-itens (`  -`) representam as regras e condições específicas daquele requisito.*

*Questões e questionamentos recorrentes são marcados como "Nota".*

### Ciclo de Vida do Cofre
- Criar novo cofre em um arquivo com senha mestra
  - Exigir digitação dupla da senha mestra para confirmação
  - Avaliar a força da senha mestra e exibir aviso informativo caso seja fraca
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
  - Segredos marcados para exclusão são ignorados pela serialização (não gravados no arquivo) e, após sucesso do salvamento, são removidos permanentemente da árvore em memória. Se o salvamento falhar, o estado em memória é preservado
  - Se o arquivo foi modificado externamente desde a última leitura ou salvamento, avisar o usuário e oferecer as opções: Sobrescrever / Salvar como novo arquivo / Cancelar
- Salvar cofre em outro arquivo
  - O arquivo de destino não pode ser o mesmo arquivo atual do cofre. *Nota: Esta restrição é uma medida de segurança para evitar a complexidade de sobrescrever um arquivo em uso e eliminar riscos de corrupção. Para gravar no arquivo atual, a função "Salvar" deve ser utilizada.*
  - Após a operação, o arquivo de trabalho atual passa a ser o novo arquivo
  - Próximas modificações e salvamentos ocorrem sobre o novo arquivo, não o original

- Descartar alterações não salvas e recarregar o cofre
  - Descarta todas as alterações realizadas desde o último salvamento (se houver) ou desde a abertura do cofre (se nunca foi salvo)
  - O cofre é recarregado ao seu estado anterior
  - Se o arquivo foi modificado externamente desde a última leitura ou salvamento, avisar o usuário antes de recarregar
- Alterar a senha mestra do cofre
  - Exigir digitação dupla para confirmação
  - Avaliar a força da senha mestra e exibir aviso informativo caso seja fraca
  - A alteração é imediata: o cofre é salvo automaticamente com a nova senha ao confirmar, incluindo todas as alterações pendentes da sessão
  - Após a alteração, não é possível descartar essa operação (o arquivo já foi regravado)
- Bloquear o cofre manualmente, automaticamente após inatividade ou por bloqueio emergencial
  - Bloquear automaticamente após tempo configurável de inatividade, com valor padrão de 5 minutos. Qualquer interação do usuário com a aplicação reseta o temporizador de inatividade
  - No bloqueio por inatividade, a aplicação apenas bloqueia o cofre e solicita a senha mestra para retomar o acesso
  - No bloqueio emergencial, a aplicação deve exibir uma tela falsa disfarçada e solicitar a senha mestra para retorno ao cofre
  - O bloqueio retorna ao fluxo de abertura do cofre, exigindo nova autenticação para retomar o acesso
  - Implementação usa memória protegida (mlock/VirtualLock) quando disponível para impedir swap de dados sensíveis. Se indisponível, a aplicação opera normalmente sem essa camada de proteção
- Sair da aplicação
  - Se houver alterações não salvas, exibir confirmação com opções: Salvar e Sair / Descartar e Sair / Cancelar
  - Se não houver alterações pendentes, sair diretamente sem confirmação
- Exportar cofre para um arquivo de intercâmbio (não criptografado)
  - O arquivo exportado contém toda a estrutura ativa do cofre: pastas, segredos ativos e modelos. Configurações de timers não são exportadas. Segredos marcados para exclusão (mesma lógica de salvamento) não são incluídos
  - Exibir aviso sobre os riscos de segurança e solicitar confirmação antes de exportar
- Importar cofre de um arquivo de intercâmbio
  - O arquivo de intercâmbio deve ser válido em estrutura e conteúdo; a existência da Pasta Geral é uma premissa. Se o arquivo for inválido ou não contiver Pasta Geral, a importação falha com mensagem de erro
    - **Pastas**: A estrutura de pastas é mesclada. Se uma pasta do arquivo de importação já existe no cofre (no mesmo caminho), seu conteúdo é mesclado com a pasta correspondente no cofre. Se a pasta não existe, ela é criada.
  - **Segredos**: Dentro de uma pasta mesclada, segredos importados que têm o mesmo nome de um segredo existente são **substituídos** pelo segredo do arquivo de importação. Segredos com nomes únicos são adicionados.
  - **Modelos**: Modelos importados com o mesmo nome de um modelo existente no cofre **substituem** o modelo do cofre.
  - *Nota sobre a política de importação*: A política de importação é intencionalmente baseada na **mesclagem** de estruturas de pastas e na **sobrescrita** de segredos e modelos com nomes conflitantes. Este comportamento contrasta com operações que usam renomeação automática para evitar perda de dados não intencional. Na importação, a sobrescrita é a ação esperada.
- Configurar o cofre
  - Todos os tempos são iniciados com valor padrão ao criar o cofre e podem ser ajustados pelo usuário via configurações do cofre. Nenhum temporizador pode ser desabilitado — todos são obrigatórios
  - Configurar tempo de bloqueio automático por inatividade (padrão: 5 minutos)
  - Configurar tempo de ocultação automática de campo sensível (padrão: 15 segundos)
  - Configurar tempo de limpeza automática da clipboard (padrão: 30 segundos)
  - Persistir preferência de tema visual
    - Os temas disponíveis devem ser discretos e adequados a ambientes públicos
    - O identificador do tema deve ser gravado no payload criptografado do cofre
    - A aplicação do tema deve ser reativa (aplicada imediatamente ao alterar a configuração)

### Consulta dos Segredos
- Exibir o cofre com suas pastas e segredos
  - Cada pasta exibe a contagem total de segredos ativos (não marcados para exclusão) contidos nela e em todas as suas subpastas recursivamente
- Exibir pasta virtual "Favoritos" como nó irmão da Pasta Geral na árvore (acima dela)
  - Lista todos os segredos com favorito = verdadeiro, percorridos em profundidade seguindo a ordem do cofre
  - A pasta virtual é somente leitura — não é possível criar, mover ou excluir segredos diretamente a partir dela
  - A pasta virtual não pode ser renomeada, movida ou excluída
- Buscar segredos por nome, nome de campo, valor de campo comum ou observação
  - A busca funciona por substring, ignorando acentuação e capitalização (case-insensitive)
  - Valores de campos sensíveis nunca participam da busca — nomes de campos sensíveis participam normalmente
  - Segredos marcados para exclusão não aparecem nos resultados de busca
- Exibir um segredo com nome, seus campos e a observação
- Exibir indicadores de estado de sessão na listagem de segredos
  - Indicador "adicionado" para segredos criados na sessão atual
  - Indicador "modificado" para segredos cujo conteúdo foi alterado na sessão atual (nome ou campos)
  - Indicador "excluído" para segredos marcados para exclusão

  - Segredos sem alterações desde o carregamento não exibem indicador
- Exibir temporariamente o valor de um campo sensível
  - Ocultar o valor automaticamente após tempo configurável, com valor padrão de 15 segundos
- Copiar temporariamente o valor de qualquer campo para a área de transferência
  - Remover o valor da área de transferência automaticamente ao bloquear ou encerrar a aplicação, ou após tempo configurável, com valor padrão de 30 segundos. A limpeza da clipboard depende do suporte do sistema operacional

### Gerenciamento de Segredos
- Criar segredo
  - A partir de um modelo existente ou como segredo sem campos de modelo — apenas com a Observação
  - O segredo pertencerá a uma pasta, escolhida no momento da criação
- Duplicar segredo existente
  - Um novo segredo é criado com o mesmo conteúdo do original
  - O novo segredo recebe automaticamente um nome único na pasta
  - O novo segredo é posicionado imediatamente após o original na lista
  - O histórico de modelo do segredo original é preservado no segredo duplicado
- Editar segredo: alterar o nome do segredo, o valor de campos e/ou observação
  - Não altera a estrutura do segredo (para alterar estrutura, use Adicionar/Renomear/Reordenar/Excluir campo)
- Alterar estrutura do segredo: adicionar campo (com nome e tipo); renomear campo; reordenar campos; excluir campo
  - Não permite alterar o tipo de um campo
  - Não permite alterar a posição, tipo ou nome da observação
- Favoritar e desfavoritar segredo
- Marcar segredo para exclusão
  - O segredo permanece na lista da pasta, visualmente sinalizado como excluído
  - Não aparece em resultados de busca enquanto marcado
  - Ao salvar com sucesso, é removido permanentemente da árvore em memória (não consta no arquivo gravado)
  - Se a pasta do segredo for excluída antes de salvar, o segredo marcado é movido junto para a pasta pai (mantendo o estado de exclusão; com renomeação automática por colisão se necessário)

- Desmarcar exclusão de segredo
  - Restaura o segredo ao estado que tinha antes de ser marcado para exclusão
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
  - Ao excluir uma pasta, seus segredos e subpastas são movidos para a pasta que a continha — inclusive segredos marcados para exclusão, que mantêm seu estado
  - Segredos movidos são adicionados ao final da lista de segredos da pasta que a continha
  - Subpastas movidas são adicionadas ao final da lista de pastas da pasta que a continha (se subpasta com mesmo nome já existe, conteúdo é mesclado)
  - Se algum segredo promovido tiver o mesmo nome de um segredo já existente na pasta pai, ele é renomeado automaticamente
  - O usuário é avisado sobre as renomeações automáticas ocorridas

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
  - A estrutura do modelo é definida apenas pelos campos do segredo, excluindo a Observação automática
  - Uma vez criado, o modelo funciona de forma independente

## Regras Transversais

### Estrutura e Pertencimento
- Segredo só pertence a uma pasta
- Pasta pode conter segredos e outras pastas
- O cofre sempre contém a pasta Geral

### Hierarquia de Pastas
- Pastas formam uma estrutura em árvore com a Pasta Geral como raiz
- Duas pastas não podem ter o mesmo caminho completo de hierarquia (nome da pasta + caminho de seus ancestrais)
- Ciclos não são permitidos — uma pasta nunca pode ser movida para dentro de seus próprios descendentes
- Cada pasta tem exatamente um ancestral direto (exceto a Pasta Geral, que é a raiz)
- Todas as pastas devem ser navegáveis a partir da Pasta Geral — nenhuma pasta pode ficar desconectada da hierarquia

### Pasta Geral
- A pasta Geral não pode ser renomeada
- A pasta Geral não pode ser movida
- A pasta Geral não pode ser excluída
- A pasta Geral pode estar vazia
- A pasta Geral, por ser a raiz da hierarquia, é o destino natural quando segredos e subpastas de uma pasta diretamente dentro dela são movidos por exclusão. Este comportamento se aplica a qualquer nível: o destino é sempre a pasta pai imediata da pasta excluída

### Nomes e Duplicidade
- Não é permitido ter dois segredos com o mesmo nome dentro da mesma pasta pai. A identidade do segredo é a composite key (pastasPai, nomePasta).
- Não há restrição quanto a duplicidade do nome entre campos de um mesmo segredo
- Não é permitido ter duas subpastas com o mesmo nome dentro da mesma pasta pai. A identidade da pasta é a composite key (pastasPai, nomePasta).
- Não é permitido ter dois modelos de segredo com o mesmo nome globalmente. A identidade do modelo é seu nome.
- Não há restrição quanto a duplicidade do nome entre campos de um mesmo modelo de segredo

### Resolução de Conflitos de Nome (Renomeação Automática)
- Quando uma operação criaria um conflito de nome (mesmo nome já existe no destino), o novo item é renomeado automaticamente
- Padrão de renomeação: adicionar sufixo numérico ` (N)` ao nome original, onde N é o menor inteiro >= 1 que garante unicidade
  - Exemplo: se "Gmail" existe, o novo item torna-se "Gmail (1)"; se "Gmail (1)" também existe, torna-se "Gmail (2)", etc.
- Aplica-se esta regra em operações de: duplicação de segredo, exclusão de pasta (promoção de segredos)
- O usuário é avisado sobre as renomeações automáticas ocorridas


### Limites
- Não há limite de quantidade para: pastas, segredos, modelos, campos em segredo e campos em modelo
- Limites são regidos pelos recursos do sistema; a aplicação não impõe limites artificiais

### Concorrência e Acesso ao Arquivo
- Verificar se o arquivo do cofre está sendo usado por outro processo; a abertura deve falhar com mensagem informativa se esse for o caso
- Verificar se o arquivo do cofre tem permissão de leitura; a abertura deve falhar com mensagem informativa se não tiver
- Verificar se o arquivo do cofre tem permissão de escrita; o salvamento deve falhar com mensagem informativa se não tiver
- Se o arquivo do cofre foi modificado externamente desde a última leitura ou salvamento, a aplicação notifica o usuário quando operações subsequentes são tentadas (salvar, descartar), oferecendo opções apropriadas para cada contexto
- Nomes de pastas e segredos podem conter qualquer caractere Unicode, exceto separadores de caminho (`/`, `\`) e caracteres de controle

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
- Observação é um campo automático que existe em todo segredo, independente de como foi criado (também através de modelos)
- Atributos: é um campo comum (sempre visível, não sensível), não pode ser renomeado, não excluído, não movido, e ocupa sempre a última posição na lista de campos do segredo
- A Observação automática **não faz parte** da estrutura do modelo: modelos não declaram, controlam, nem podem incluir um campo denominado "Observação" para evitar conflito com o campo automático
- Ao criar modelo a partir de um segredo existente, a Observação do segredo é sempre ignorada — não é copiada para o modelo
- O uso responsável da observação é por conta e risco do usuário — o campo não prevê ocultação nem tratamento especial

### Força da Senha Mestra
- Uma senha mestra é considerada **forte** quando satisfaz todos os critérios abaixo:
  - Comprimento mínimo de 12 caracteres
  - Contém pelo menos uma letra maiúscula (A–Z)
  - Contém pelo menos uma letra minúscula (a–z)
  - Contém pelo menos um dígito (0–9)
  - Contém pelo menos um caractere especial (qualquer caractere que não seja letra nem dígito)
- Uma senha que não satisfaça todos os critérios acima é considerada **fraca**
- O aviso de senha fraca é apenas informativo — não impede a criação do cofre nem a alteração da senha mestra. A responsabilidade pela escolha da senha é inteiramente do usuário

### Gerenciamento de Senha na Sessão
- A senha é fornecida uma única vez ao abrir o cofre
- A senha ativa é usada para todas as operações de criptografia durante a sessão (salvar, descartar)
- Não há re-solicitação de senha para salvamento, descarte ou operações de salvamento em novo arquivo
- Alterar a senha mestra é uma operação imediata e irrevogável — o cofre é regravado na hora; não é uma mudança de estado da sessão
- Ao bloquear o cofre, a senha é removida da memória e será novamente solicitada na próxima abertura

### Limpeza de Memória e Buffers
- Ao bloquear o cofre (manual, por inatividade ou emergencial) ou sair da aplicação:
  - A senha mestra é limpa da memória (sobrescrita com zeros)
  - Buffers contendo dados sensíveis (valores de campos, payloads descriptografados) são descartados
  - O terminal é limpo (clear screen) antes de devolver o controle ao shell, evitando que dados visíveis na TUI permaneçam no buffer do terminal
- Ao descartar alterações não salvas, segue o mesmo protocolo de limpeza

## Requisitos Não Funcionais

- **Criptografia**: AES-256-GCM para criptografia dos dados; Argon2id para derivação de chave a partir da senha mestra
- **Formato de armazenamento**: JSON criptografado com AES-256-GCM, encapsulado em arquivo binário com extensão `.abditum`
- **Compatibilidade**: Windows, macOS e Linux
- **Usabilidade e Design**:
  - **Interface Moderna (TUI)**: A aplicação utiliza elementos visuais ricos (gradientes, cores semânticas), mas mantém a discrição necessária para ambientes públicos.
  - **Customização Reativa**: A troca de temas deve ser instantânea e seguir as definições técnicas do Design System do projeto.
- **Privacidade**: Ausência total de logs da aplicação (stdout/stderr) que contenham caminhos de arquivos de cofre, nomes de segredos ou valores de campos
- **Portabilidade**: A aplicação é distribuída como um único arquivo binário executável, sem instalação. Pode ser copiada para qualquer local do sistema de arquivos e executada diretamente. Não utiliza arquivos de configuração, dados ou estado fora do arquivo do cofre — exceto artefatos transitórios (`.abditum.tmp`) e backups (`.abditum.bak`, `.abditum.bak2`) explicitamente previstos pela própria aplicação
- **Compatibilidade retroativa**: A aplicação de versão N é capaz de abrir arquivos de cofre criados em qualquer versão anterior do formato suportado pela aplicação. Ao abrir um arquivo antigo, o payload descriptografado é migrado em memória para o modelo atual; ao salvar, o arquivo é sempre regravado no formato da versão atual da aplicação
- **Confiabilidade (Salvamento Atômico)**: A gravação do arquivo do cofre é sempre atômica, garantindo que uma falha durante o processo não corrompa nem o arquivo principal nem o backup. O protocolo se aplica a toda operação que sobrescreve um arquivo de cofre já existente:
  - Os dados são gravados primeiramente em `.abditum.tmp`, no mesmo diretório do cofre
  - Somente após a gravação bem-sucedida, o `.abditum.tmp` é renomeado, substituindo o arquivo original
  - Em caso de falha, o `.abditum.tmp` é apagado imediatamente para evitar persistência de dados fora do arquivo final
  - Ao substituir um arquivo existente, manter cópia de backup com extensão `.abditum.bak`:
    - Se já existir um `.abditum.bak`, renomeá-lo temporariamente para `.abditum.bak2` antes de gerar o novo backup
    - Se a operação for concluída com sucesso, apagar o `.abditum.bak2`, preservando apenas o novo `.abditum.bak`
    - Se a operação falhar antes da consolidação, restaurar o `.abditum.bak2` para `.abditum.bak` sempre que possível
    - Em caso de falha após a geração do backup, exibir mensagem de erro informando que existe um backup disponível para intervenção manual do usuário
  - Criação de novo cofre em caminho vazio e salvamento em novo caminho vazio não utilizam `.abditum.tmp` — a gravação ocorre diretamente no destino final
- **Formato do arquivo .abditum**: deve seguir rigorosamente formato-arquivo-abditum.md

## Requisitos v2

*As funcionalidades a seguir estão planejadas para uma futura versão e não serão especificadas ou detalhadas neste momento. Servem como um registro de ideias para a evolução do produto.*

### Exibição Parcial de Campos Sensíveis
- Permitir configurar exibição parcial de campos sensíveis — revelar apenas parte do valor (ex: últimos 4 dígitos de número de cartão de crédito)
- Permitir que o usuário defina, por campo sensível, uma regra de exibição parcial (ex: mostrar últimos N caracteres, primeiros N caracteres, ou padrão mascarado como "•••• •••• •••• 1234")
- A exibição parcial não substitui a revelação completa — é um modo adicional de visualização rápida

### Gerador de Senhas

A ser especificado em v2.

### Compartilhamento via QR Code

Renderizar um QR code diretamente na TUI (usando blocos ASCII/Unicode) contendo o valor de um campo, para transferência rápida, offline e segura para outro dispositivo.

#### Decisões Pendentes (v2)
- **Exibição parcial**: Regras de mascaramento são por campo individual ou por tipo de modelo? Quais padrões pré-definidos oferecer?
- **QR Code — escopo do conteúdo**: O QR code deve conter apenas o valor bruto do campo (texto plano, legível por qualquer câmera), ou há intenção de protocolo mais rico? Não existe padrão universal para importação de credenciais entre gerenciadores de senha via QR; o único protocolo padronizado no domínio (`otpauth://`) é exclusivo de TOTP.

### Relatório de Saúde do Cofre (Auditoria)

Analisar localmente todas as senhas armazenadas e alertar o usuário sobre senhas fracas, reutilizadas em múltiplos segredos ou muito antigas.

A ser especificado em v2.

### Tags

Categorização de segredos por tags, com filtragem por tag. A ser especificado em v2.

### Histórico de Versões de Segredos

Registro de versões anteriores de um segredo, com possibilidade de visualização e restauração. A ser especificado em v2.

### Recuperação de Artefatos Órfãos

Ao abrir um cofre, detectar a presença de artefatos residuais de uma operação de salvamento anterior interrompida (`.abditum.tmp` ou `.abditum.bak2` no mesmo diretório). Nesses casos:
- Alertar o usuário explicando que a última operação de salvamento pode não ter sido concluída normalmente
- Informar que existe um backup (`.abditum.bak`) disponível para inspeção manual
- Oferecer a opção de recuperar o cofre a partir do arquivo residual disponível

A ser especificado em v2.

## Fora de Escopo

Funcionalidades deliberadamente excluídas desta versão:
- **Tratamento especial do cofre em modo somente leitura**: A aplicação pode abrir um cofre em modo somente leitura, mas não há tratamento especial para esse estado além de falhar ao tentar salvar, com mensagem de erro informativa
- **Senha Falsa de Coação (Duress Password)**: Uma senha alternativa que abre uma versão restrita do cofre para proteger os dados reais em situações de ameaça. Embora valiosa, a complexidade de implementação para garantir a segurança e a usabilidade corretas a coloca fora do escopo da versão inicial.
- **TOTP (Two-Factor Authentication)**: Geração de código de autenticação de dois fatores — excluído permanentemente, sem previsão para nenhuma versão futura
- **Backup**: A aplicação não cria, gerencia nem armazena cópias de segurança do cofre. Manter cópias de segurança é responsabilidade exclusiva do usuário
- **Recuperação de dados**: A criptografia adotada não permite recuperação parcial de arquivos corrompidos. Não há mecanismo de reparo, importação forçada ou abertura em modo degradado
- **Autenticação por Keyfile / Token de Hardware**: Exigir, além da senha mestra, um arquivo físico (keyfile) ou token USB (ex: YubiKey) para descriptografar o cofre — excluído permanentemente, sem previsão para nenhuma versão futura
- **Armazenamento na nuvem**: Contraria a filosofia offline e portátil da aplicação — excluído permanentemente
- **Múltiplos cofres abertos simultaneamente**: Invariante de design — só pode existir um cofre ativo por vez — excluído permanentemente
- **App mobile ou web**: A aplicação é TUI portátil por design — excluído permanentemente



