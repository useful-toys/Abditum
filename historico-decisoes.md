## Histórico de Decisões

**A pasta padrão do cofre se chama "Geral" e não pode ser renomeada**
Eliminar o conceito técnico de "raiz" simplifica o vocabulário do produto — o usuário só precisa conhecer "pasta" e "segredo". A pasta Geral garante que sempre há um destino disponível para os segredos, sem criar o problema de um cofre sem pasta obrigatória. Não permitir renomear evita que seu comportamento especial fique invisível.

**Todo segredo vive dentro de uma pasta — não existe segredo fora de uma pasta**
A alternativa seria ter um nível "raiz" onde segredos e pastas coexistem sem estar dentro de nada. Essa estrutura gerou dificuldade consistente de descrição no documento, o que sinalizou que o conceito complicava mais do que ajudava. A pasta Geral resolve o problema sem introduzir um conceito novo.

**O tipo de um campo não pode ser alterado após a criação**
Converter um campo sensível em comum exporia seu conteúdo, que passaria a ser exibido sem proteção. Essa mudança silenciosa seria um risco de segurança difícil de perceber. A alternativa — excluir e recriar o campo — é explícita e segura.

**O segredo criado a partir de um modelo não mantém vínculo com ele**
Manter vínculo por referência significaria que alterações no modelo afetariam segredos já criados, o que pode ser destrutivo e imprevisível para o usuário. A estrutura de um segredo é imutável após sua criação — o modelo fica livre para evoluir sem afetar segredos existentes. O nome do modelo é registrado apenas como histórico do momento da criação. Essa abordagem oferece ao usuário a segurança de que seus segredos não sofrerão mudanças inesperadas e ao gestor de modelos a liberdade de refatorar estruturas conforme necessário.

**Campos sensíveis não participam da busca**
Realizar busca sobre valores sensíveis exigiria mantê-los descriptografados em memória durante a operação, aumentando a superfície de exposição. A busca ocorre apenas sobre campos comuns e observações.

**O termo "campo" foi adotado para descrever os elementos de um segredo**
Alternativas como "dado", "detalhe" e "faceta" foram consideradas. "Detalhe" diminui semanticamente dados críticos como a senha. "Faceta" sugere perspectiva, não composição. "Campo" já faz parte do vocabulário cotidiano do usuário — qualquer pessoa que preencheu um formulário entende o termo sem explicação.

**Auditoria de senhas e TOTP foram descartados do escopo**
Funcionalidades identificadas durante o processo de definição, mas que não serão implementadas por estarem fora do escopo do produto na sua versão atual.

**Modelos pré-definidos e personalizados foram unificados em um único conceito**
A distinção entre modelo pré-definido e personalizado é de origem, não de comportamento — na prática ambos funcionam da mesma forma. Manter dois termos adicionava complexidade ao vocabulário sem benefício real. Os modelos criados automaticamente ao criar o cofre são apenas o ponto de partida, editáveis e removíveis como qualquer outro.

**Importação: pastas mescladas silenciosamente, segredos com nome ajustado**
Na importação de arquivo, pastas com mesmo identificador são mescladas automaticamente, enquanto segredos com conflito de nome recebem nome ajustado visível (ex: "Segredo (1)"). A diferença reflete que apenas o que realmente importa é tratado explicitamente — estrutura de pastas é raramente conflitante em casos reais, enquanto conflito de nomes de segredo é mais comum. Simplificar a interface de resolução de conflitos reduz fricção sem sacrificar segurança.

**Busca funciona por substring, ignorando acentuação e capitalização**
O algoritmo de busca por substring é a expectativa padrão dos usuários — "pass" deve encontrar "password" e "senhaQuebrada". A busca case-insensitive reduz fricção sem comprometer segurança, já que a busca ocorre apenas sobre campos comuns (não sensíveis). Ignorar acentuação evita que "café" não encontre "cafe", problema comum em buscas que limita a usabilidade. Não implementar regex mantém a simplicidade tanto para o usuário quanto para a implementação.

**Identidade de pastas, modelos e segredos é determinada por identificador único interno, não por nome**
Pastas são identificadas por nome (na mesma pasta pai, nomes devem ser únicos). Modelos e segredos são identificados por identificador único interno gerado pelo sistema — dois segredos podem ter o mesmo nome (assim como dois modelos). Isso oferece flexibilidade no produto (permitir renomear sem perder vínculo com o modelo original, importar segredos com nomes duplicados) sem criar ambiguidade na representação. Nomes são apenas metadados descritivos, não identificadores.

Essa é uma decisão deliberada para simplificar a arquitetura da aplicação e evitar complexidade de sincronização. Reconhecemos que nomes repetidos podem criar problemas de usabilidade (confusão visual, dificuldade de identificação pelo usuário), mas confiamos que a UX será capaz de mitigar isso através de contexto visual, ordenação inteligente e affordances claras. A alternativa — impor unicidade de nomes — criaria atrito desnecessário (renomear segredos para poder importar, impossibilidade de ter "Login" e "Login" em contextos diferentes).

**Sem limites técnicos: "bom senso" governa as restrições**
O sistema não impõe limites arbitrários (quantidade de segredos, profundidade de pastas, campos por modelo, etc.). Limites técnicos surgem apenas quando o hardware/sistema operacional os impõe. A experiência do usuário é governada pelo bom senso — um cofre com 100 mil segredos pode ser tecnicamente possível, mas não é um caso de uso esperado. Essa abordagem simplifica o produto (sem validações desnecessárias), confia na racionalidade do usuário, e evita a frustração de butts arbitrários. Se um limite técnico se tornar necessário na implementação, será documentado explicitamente.

**"Observação" é um campo especial criado automaticamente em todo segredo para fins de UX**
A observação oferece um espaço dedicado para notas do usuário sem ocupar um campo customizado no modelo. É um campo comum (sempre visível), não pode ser renomeado ou deletado, e está disponível em todo segredo independentemente do modelo. Essa abordagem simplifica a UX (o usuário sempre tem um lugar previsível para adicionar notas) sem criar complexidade conceitual — é um campo como qualquer outro, apenas com restrições especiais que o usuário entende naturalmente como proteção contra exclusão acidental.

**A estrutura de pastas é uma árvore — ciclos não são permitidos**
Pastas formam uma hierarquia em árvore com a Pasta Geral como raiz. Essa estrutura é explicitamente protegida contra ciclos — não é possível mover uma pasta para dentro de seus próprios descendentes. Ciclos criariam ambiguidade de navegação, impossibilidade de definir um "caminho" único para cada pasta, e complexidade desnecessária na representação. A árvore garante que cada pasta tem exatamente um ancestral direito (exceto a raiz), oferecendo clareza total na hierarquia. Essa restrição é validada em tempo real ao mover pastas, com feedback claro ao usuário sobre por que uma operação não é permitida.

**Salvamento e descarte usam a senha fornecida ao abrir, sem solicitar novamente**
Ao abrir o cofre, o usuário fornece a senha mestra. Essa mesma senha é usada para todas as operações de criptografia (salvar, descartar, alterações futuras) durante a sessão. O sistema nunca solicita a senha novamente para salvar ou descartar — a senha permanece na memória da aplicação enquanto o cofre está em uso. Se o usuário alterar a senha mestra, a nova senha passa a ser usada para próximos salvamentos. Essa abordagem simplifica a UX (fluir de trabalho sem interrupções), oferece segurança implícita (bloquear o cofre limpa a memória), e reduz pontos de falha (menos oportunidades de senha incorreta invalidar uma operação de salvamento).

