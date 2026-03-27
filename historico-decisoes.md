## Histórico de Decisões

**A pasta padrão do cofre se chama "Geral" e não pode ser renomeada**
Eliminar o conceito técnico de "raiz" simplifica o vocabulário do produto — o usuário só precisa conhecer "pasta" e "segredo". A pasta Geral garante que sempre há um destino disponível para os segredos, sem criar o problema de um cofre sem pasta obrigatória. Não permitir renomear evita que seu comportamento especial fique invisível.

**Todo segredo vive dentro de uma pasta — não existe segredo fora de uma pasta**
A alternativa seria ter um nível "raiz" onde segredos e pastas coexistem sem estar dentro de nada. Essa estrutura gerou dificuldade consistente de descrição no documento, o que sinalizou que o conceito complicava mais do que ajudava. A pasta Geral resolve o problema sem introduzir um conceito novo.

**O tipo de um campo não pode ser alterado após a criação**
Converter um campo sensível em comum exporia seu conteúdo, que passaria a ser exibido sem proteção. Essa mudança silenciosa seria um risco de segurança difícil de perceber. A alternativa — excluir e recriar o campo — é explícita e segura.

**O segredo criado a partir de um modelo não mantém vínculo com ele**
Manter vínculo por referência significaria que alterações no modelo afetariam segredos já criados, o que pode ser destrutivo e imprevisível para o usuário. O nome do modelo é registrado apenas como histórico do momento da criação.

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

