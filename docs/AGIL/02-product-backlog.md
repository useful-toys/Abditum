# 02 — Product Backlog

## 02.1 Épicos

| ID | Épico | Descrição |
|---|---|---|
| E01 | Ciclo de Vida do Cofre | Criar, abrir, salvar, bloquear, configurar e encerrar o cofre — todas as operações que governam a existência e a persistência do arquivo criptografado. |
| E02 | Gerenciamento de Segredos | Criar, editar, duplicar, favoritar, excluir reversivelmente, restaurar, mover e reordenar segredos dentro da hierarquia do cofre. |
| E03 | Gerenciamento de Hierarquia | Criar, renomear, mover, reordenar e excluir pastas, mantendo a integridade da árvore de organização. |
| E04 | Modelos de Segredo | Criar, editar, excluir e derivar modelos de segredo, fornecendo estruturas reutilizáveis para tipos comuns de informação confidencial. |
| E05 | Navegação e Busca | Navegar pela hierarquia, visualizar segredos, exibir/ocultar campos sensíveis e buscar segredos em memória. |
| E06 | Área de Transferência | Copiar campos para o clipboard com limpeza automática temporizada e ao bloquear/fechar o cofre. |
| E07 | Importação e Exportação | Exportar cofre para JSON plain text e importar de JSON plain text com resolução de conflitos. |
| E08 | Segurança e Proteção | Criptografia, derivação de chave, bloqueio por inatividade, proteção contra shoulder surfing, minimização de dados em memória, privacidade de logs. |
| E09 | Interface TUI | Layout de dois painéis, barra de status, barra de ajuda contextual, file picker, feedback visual, tela inicial e fluxos modais. |

---

## 02.2 User Stories

### E01 — Ciclo de Vida do Cofre

| ID | User Story | Prioridade |
|---|---|---|
| US-0101 | Como usuário, quero **criar um novo cofre** informando caminho e senha mestra (com confirmação dupla), para começar a armazenar meus segredos de forma segura. | Must |
| US-0102 | Como usuário, quero **abrir um cofre existente** informando caminho e senha mestra, para acessar meus segredos previamente armazenados. | Must |
| US-0103 | Como usuário, quero **salvar o cofre** no caminho atual, para persistir as alterações feitas na sessão. | Must |
| US-0104 | Como usuário, quero **salvar o cofre em um novo caminho**, para criar uma cópia do cofre em outro local. | Should |
| US-0105 | Como usuário, quero **descartar alterações não salvas e recarregar o cofre**, para reverter mudanças indesejadas sem fechar a aplicação. | Should |
| US-0106 | Como usuário, quero **alterar a senha mestra** do cofre ativo, para atualizar minha credencial de acesso quando necessário. | Must |
| US-0107 | Como usuário, quero **bloquear o cofre manualmente**, para proteger meus dados imediatamente quando me afasto da máquina. | Must |
| US-0108 | Como usuário, quero que o cofre **bloqueie automaticamente por inatividade**, para que meus dados fiquem protegidos mesmo se eu esquecer a tela aberta. | Must |
| US-0109 | Como usuário, quero **configurar os tempos** de bloqueio por inatividade, reocultação de campos sensíveis e limpeza da área de transferência, para ajustar o comportamento de segurança às minhas necessidades. | Should |
| US-0110 | Como usuário, quero **sair da aplicação** a qualquer momento, sendo alertado se houver alterações não salvas, para evitar perda acidental de dados. | Must |

### E02 — Gerenciamento de Segredos

| ID | User Story | Prioridade |
|---|---|---|
| US-0201 | Como usuário, quero **criar um segredo a partir de um modelo**, para preencher rapidamente campos com uma estrutura pré-definida. | Must |
| US-0202 | Como usuário, quero **criar um segredo vazio** (sem campos iniciais), para montar minha própria estrutura de campos via edição avançada. | Must |
| US-0203 | Como usuário, quero **editar um segredo no modo padrão** (nome, observação e valores dos campos), para atualizar dados sem alterar a estrutura. | Must |
| US-0204 | Como usuário, quero **editar um segredo no modo avançado** (adicionar, renomear, excluir e reordenar campos), para alterar a estrutura do segredo conforme necessário. | Must |
| US-0205 | Como usuário, quero **duplicar um segredo existente**, para criar uma cópia com nova identidade como ponto de partida para um segredo similar. | Should |
| US-0206 | Como usuário, quero **favoritar/desfavoritar um segredo**, para destacar os segredos de uso mais frequente e acessá-los rapidamente via pasta virtual de Favoritos. | Should |
| US-0207 | Como usuário, quero **excluir reversivelmente um segredo** (soft delete para a Lixeira), para removê-lo da hierarquia com possibilidade de restauração até o próximo salvamento. | Must |
| US-0208 | Como usuário, quero **restaurar um segredo da Lixeira**, para recuperá-lo na pasta de origem (ou raiz, se a pasta não existir mais) antes do próximo salvamento. | Must |
| US-0209 | Como usuário, quero **mover um segredo para outra pasta ou raiz**, para reorganizar a hierarquia do cofre. | Should |
| US-0210 | Como usuário, quero **reordenar um segredo** entre os irmãos da mesma pasta, para controlar a ordem de exibição. | Could |
| US-0211 | Como usuário, quero **buscar segredos** por nome, nome de campo, valor de campos texto ou observação, para localizar rapidamente a informação desejada. | Must |

### E03 — Gerenciamento de Hierarquia

| ID | User Story | Prioridade |
|---|---|---|
| US-0301 | Como usuário, quero **criar uma pasta** na raiz ou dentro de outra pasta, para organizar meus segredos em categorias. | Must |
| US-0302 | Como usuário, quero **renomear uma pasta**, para corrigir ou atualizar o nome da categoria. | Should |
| US-0303 | Como usuário, quero **mover uma pasta** para outra pasta ou raiz, para reorganizar a hierarquia com todos os filhos acompanhando. | Should |
| US-0304 | Como usuário, quero **reordenar uma pasta** entre as irmãs da mesma coleção, para controlar a ordem de exibição. | Could |
| US-0305 | Como usuário, quero **excluir uma pasta**, com seus segredos e subpastas sendo promovidos ao pai, para simplificar a hierarquia sem perder conteúdo. | Should |

### E04 — Modelos de Segredo

| ID | User Story | Prioridade |
|---|---|---|
| US-0401 | Como usuário, quero **criar um modelo de segredo** com campos personalizados (nome e tipo), para padronizar a criação de segredos similares. | Must |
| US-0402 | Como usuário, quero **editar um modelo existente** (nome, campos), para ajustar a estrutura de criações futuras sem afetar segredos já existentes. | Should |
| US-0403 | Como usuário, quero **excluir um modelo de segredo**, para remover estruturas que não uso mais. | Should |
| US-0404 | Como usuário, quero **criar um modelo a partir de um segredo existente**, para reaproveitar a estrutura de campos como base de um novo modelo. | Could |
| US-0405 | Como usuário, quero que o cofre **venha com modelos pré-definidos** (Login, Cartão de Crédito, API Key) que eu possa editar ou remover. | Must |

### E05 — Navegação e Busca

| ID | User Story | Prioridade |
|---|---|---|
| US-0501 | Como usuário, quero **navegar pela árvore hierárquica** do cofre, expandindo e colapsando pastas, para encontrar meus segredos visualmente. | Must |
| US-0502 | Como usuário, quero **visualizar os detalhes de um segredo** no Painel do Segredo ao selecioná-lo na árvore. | Must |
| US-0503 | Como usuário, quero **exibir temporariamente o valor de um campo sensível** com toggle e reocultação automática, para consultar a informação sem exposição prolongada. | Must |
| US-0504 | Como usuário, quero ver **pastas virtuais de Favoritos** (topo) e **Lixeira** (rodapé) quando houver conteúdo aplicável, para acesso rápido e visualização de exclusões reversíveis. | Should |
| US-0505 | Como usuário, quero **buscar com highlight** no nome do segredo quando o casamento ocorrer no nome, para identificar visualmente o match na árvore. | Should |

### E06 — Área de Transferência

| ID | User Story | Prioridade |
|---|---|---|
| US-0601 | Como usuário, quero **copiar qualquer campo** (incluindo sensíveis) para a área de transferência, para colar em outro aplicativo rapidamente. | Must |
| US-0602 | Como usuário, quero que a área de transferência seja **limpa automaticamente** após o tempo configurado (padrão 30s), para reduzir o risco de exposição. | Must |
| US-0603 | Como usuário, quero que a área de transferência seja **limpa ao bloquear ou fechar o cofre**, para impedir vazamento de dados após o encerramento da sessão. | Must |

### E07 — Importação e Exportação

| ID | User Story | Prioridade |
|---|---|---|
| US-0701 | Como usuário, quero **exportar o cofre para JSON plain text**, com aviso de risco e confirmação, para backup ou migração manual. | Should |
| US-0702 | Como usuário, quero **importar um arquivo JSON plain text** para o cofre ativo, com resolução automática de conflitos (merge de pastas, novo ID para segredos duplicados, sufixo numérico para nomes, sobreposição de modelos), para incorporar dados de outro cofre. | Should |

### E08 — Segurança e Proteção

| ID | User Story | Prioridade |
|---|---|---|
| US-0801 | Como usuário, quero que meus dados sejam **criptografados com AES-256-GCM e chave derivada via Argon2id**, para garantir proteção contra acesso não autorizado e ataques offline. | Must |
| US-0802 | Como usuário, quero que o **salvamento seja atômico** (escrita em .tmp → backup .bak com rotação → rename), para que falhas nunca corrompam meu cofre. | Must |
| US-0803 | Como usuário, quero **proteção contra shoulder surfing** via atalho que oculta toda a interface, para proteger meus dados em ambientes compartilhados. | Should |
| US-0804 | Como usuário, quero que a aplicação **nunca registre em logs** caminhos de cofre, nomes de segredos ou valores de campos, para garantir privacidade total. | Must |
| US-0805 | Como usuário, quero que dados sensíveis em memória sejam **minimizados e limpos** ao bloquear ou fechar o cofre, para reduzir a superfície de ataque. | Must |

### E09 — Interface TUI

| ID | User Story | Prioridade |
|---|---|---|
| US-0901 | Como usuário, quero ver uma **tela inicial com ASCII art** e opções de criar, abrir cofre, ajuda e sair, para ter uma experiência de boas-vindas clara. | Must |
| US-0902 | Como usuário, quero um **layout de dois painéis** (Hierarquia + Segredo) com barra de status e barra de ajuda contextual, para operar o cofre de forma eficiente. | Must |
| US-0903 | Como usuário, quero um **file picker integrado à TUI** para navegar e selecionar arquivos, para não precisar digitar caminhos manualmente. | Should |
| US-0904 | Como usuário, quero **feedback visual** via toasts (sucesso, erro, alerta, informação) e confirmações bloqueantes para ações destrutivas, para entender o resultado das minhas ações. | Must |
| US-0905 | Como usuário, quero um **countdown visual da limpeza do clipboard**, para saber quanto tempo resta antes da área de transferência ser limpa. | Could |
| US-0906 | Como usuário, quero **indicadores visuais** na árvore para segredos favoritos, novos, modificados e em edição, para identificar o estado de cada item rapidamente. | Should |

---

## 02.3 Critérios de Aceitação (exemplos representativos)

### US-0101 — Criar novo cofre

- [ ] A aplicação apresenta file picker para selecionar o caminho de destino.
- [ ] A senha mestra é solicitada com digitação dupla para confirmação.
- [ ] A aplicação exibe aviso de irrecuperabilidade (Zero Knowledge) antes da criação.
- [ ] Se já existir arquivo no destino, a aplicação exige confirmação de sobrescrita e gera backup `.bak` do arquivo existente.
- [ ] O cofre é criado com pastas pré-definidas (Sites, Financeiro, Serviços) e modelos pré-definidos (Login, Cartão de Crédito, API Key).
- [ ] Após a criação, o cofre entra em estado `Cofre Salvo`.

### US-0102 — Abrir cofre existente

- [ ] A aplicação apresenta file picker para selecionar o arquivo `.abditum`.
- [ ] A assinatura `magic` (ABDT) e a `versão_formato` são validadas antes de solicitar a senha.
- [ ] A derivação de chave (Argon2id) completa na faixa de 0,8 s a 1,5 s.
- [ ] Se a senha estiver incorreta, a aplicação informa erro sem revelar detalhes criptográficos.
- [ ] Se o formato for histórico e suportado, a migração ocorre em memória de forma transparente.
- [ ] Após a abertura, o cofre entra em estado `Cofre Salvo`.

### US-0207 — Excluir segredo reversivelmente

- [ ] A aplicação exige confirmação antes da exclusão.
- [ ] O segredo desaparece da hierarquia principal e aparece na pasta virtual Lixeira.
- [ ] O segredo não pode ser editado enquanto estiver na Lixeira.
- [ ] A pasta de origem e o estado anterior são memorizados para restauração.
- [ ] O cofre entra em estado `Cofre Modificado`.

### US-0211 — Buscar segredos

- [ ] A busca varre nome do segredo, nome de campo, valores de campos `texto` e observação.
- [ ] Campos `texto sensível` nunca participam da busca.
- [ ] A árvore é filtrada mostrando apenas matches, preservando a estrutura de pastas.
- [ ] Quando o casamento ocorre no nome, o trecho pesquisado recebe highlight.
- [ ] Durante a busca, apenas sair, navegar e visualizar estão disponíveis.
- [ ] Selecionar um resultado ou cancelar encerra a busca, retornando ao estado anterior.

---

## 02.4 Definition of Ready (DoR)

Uma User Story está pronta para entrar na sprint quando:

- [ ] A descrição segue o formato "Como... quero... para..." e está clara para a equipe.
- [ ] Critérios de aceitação estão documentados e verificáveis.
- [ ] Dependências com outras stories estão identificadas.
- [ ] O escopo está delimitado — o que está dentro e fora da story é explícito.
- [ ] A story é pequena o suficiente para ser concluída em uma sprint.
- [ ] Se houver impacto em criptografia ou formato de arquivo, os detalhes técnicos estão documentados.

---

## 02.5 Definition of Done (DoD)

Uma User Story está concluída quando:

- [ ] O código implementa todos os critérios de aceitação.
- [ ] Testes unitários cobrem a lógica de domínio (Manager/entidades).
- [ ] Testes de integração cobrem o fluxo de ponta a ponta quando aplicável.
- [ ] Golden files visuais (80×24) estão atualizados para telas afetadas.
- [ ] Nenhum dado sensível (caminhos, nomes, valores) aparece em stdout/stderr.
- [ ] O código compila e os testes passam em Windows, macOS e Linux (CI).
- [ ] Comentários explicativos foram adicionados em trechos de criptografia, TUI e decisões não triviais.
- [ ] A story foi revisada por pelo menos um membro da equipe.
