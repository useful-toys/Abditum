## Histórico de Decisões

# ARQUITETURA & HIERARQUIA

## 1. Pasta "Geral" como raiz não-renomeável ✓

**Decisão:** A pasta padrão do cofre se chama "Geral", não pode ser renomeada e é o ponto raiz da hierarquia.

**Contexto:** 
- Eliminar o conceito técnico de "raiz" simplifica o vocabulário do produto — usuário só conhece "pasta" e "segredo"
- Pasta Geral garante destino sempre disponível para segredos órfãos
- Sem renomear, comportamento especial permanece protegido e invisível

**Justificativa:** Abordagem simplifica UX (sem nível técnico extra) e resolve o problema de sempre haver um lugar seguro para segredos sem criar conceitos abstratos.

**Consequências:** Pasta Geral é invariante não-recuperável — se ausente/corrompida, arquivo é rejeitado.



## 2. Segredos só vivem dentro de pastas ✓

**Decisão:** Todo segredo vive dentro de uma pasta — não existe segredo "solto" ou em nível raiz.

**Contexto:**
- Alternativa seria ter nível "raiz" onde segredos e pastas coexistem sem contenção
- Essa estrutura gerou dificuldade consistente de descrição no produto, sinalizando que complicava mais do que ajudava
- Papel "Geral" resolve o problema sem introduzir conceito novo

**Consequências:** Modelo conceitual mais simples; sem ambiguidade entre segredos orphan e pastas.



## 3. Hierarquia em árvore, ciclos não permitidos ✓

**Decisão:** Pastas formam hierarquia em árvore com Pasta Geral como raiz. Não é possível mover pasta para dentro de seus próprios descendentes.

**Contexto:**
- Ciclos criariam ambiguidade de navegação
- Impossibilita definição de "caminho único" para cada pasta
- Adiciona complexidade desnecessária na representação

**Justificativa:** Árvore garante clareza total — cada pasta tem exatamente um ancestral direto (exceto raiz); caminho único até raiz.

**Consequências:** Validação em tempo real ao mover pastas; feedback claro ao usuário sobre por que operação não é permitida.



## 4. Exclusão de pasta promove conteúdo um nível ✓

**Decisão:** Ao excluir pasta, seus segredos e subpastas são movidos para a pasta pai imediata (um nível acima).

**Contexto:**
- Alternativa de mover tudo para Pasta Geral seria destrutiva em estruturas profundas
- Perda de contexto local prejudicaria organização construída pelo usuário
- Se pasta excluída está dentro de Geral, comportamento geral se aplica naturalmente

**Justificativa:** Previsibilidade local + conservação de contexto.

**Consequências:** Reorganização simples; conteúdo sobe um nível, sem deslocamentos abruptos nem perda de estrutura.



## 5. Identidades Naturais: Composite Keys por Nome, Sem UUID ✓

**Decisão:** 
- Segredos identificados por composite key (pastasPai, nomePasta) — nome único dentro da pasta pai
- Pastas identificadas por composite key (pastasPai, nomePasta) — nome único entre irmãs
- Modelos identificados por nome (único globalmente no cofre)
- Não há UUID/NanoID no modelo de dados persistido

**Contexto:**
- UUID anterior criava "identificador fantasma" — theoretically imutável mas na prática não era usado como identidade de negócio
- Assimetria: Segredo tinha UUID, mas identidade verdadeira era (pastaId, nome); Pasta tinha UUID mas identidade era (parentId, nome); Modelo tinha UUID mas identidade era nome
- Essa assimetria gerava confusão em DDD — qual é realmente a identidade?
- Solução: eliminar fantasmas, usar identidades semânticas baseadas em nomes (composite keys)

**Justificativa:** 
- Simetria total: todas as entidades usam identidades naturais (nomes ou composite keys)
- DDD limpo: Equals compara o que é realmente único, não ID artificioso
- JSON legível: identidades são explícitas no documento, sem UUIDs ocultos
- Simplicidade: renomeação/movimentação é natural (identidade muda semanticamente)
- Colisão é trivial: auto-rename com sufixo numérico se duas entidades com mesmo nome tentam coexistir

**Consequências:** 
- Renomeação muda identidade (não é refatoração "segura" em sentido tradicional, mas isso é correto semanticamente)
- Ao mover pasta/segredo para destino com colisão de nome, operação não falha — renomeia automaticamente
- Importação com colisão: pastas mesclam, segredos trocam valores (merge by name), modelos substituem
- Estrutura JSON é a fonte de verdade — não há referências cruzadas por ID técnico a ser sincronizado



# CAMPOS & MODELOS

## 6. Tipo de campo imutável após criação ✓

**Decisão:** O tipo de um campo (comum / sensível) não pode ser alterado após a criação do campo.

**Contexto:**
- Converter campo sensível→comum exporia conteúdo sem proteção
- Mudança silenciosa seria risco de segurança difícil de perceber
- Alternativa (deletar e recriar) é explícita e segura

**Justificativa:** Segurança vem de explicitação — usuário entende claramente quando ocorre mudança de proteção.

**Consequências:** Requer atrito (usuário re-cria campo), mas oferece garantia de auditabilidade completa.



## 7. Segredo desvinculado de modelo após criação ✓

**Decisão:** Segredo criado a partir de modelo não mantém vínculo dinâmico com ele — é cópia estrutural.

**Contexto:**
- Vínculo por referência significaria alterações no modelo afetarem segredos criados
- Seria destrutivo e imprevisível para usuário (dados mudam sem ação)
- Nome do modelo é armazenado apenas como histórico da criação

**Justificativa:** Segredo é snapshot imutável; modelo evolui livremente; oferece segurança e liberdade de refatoração simultâneas.

**Consequências:** Usuário não vê mudanças em modelo aplicadas retroativamente a segredos; gestor de modelos tem liberdade total de evolução.



## 8. Termo "campo" para elementos de segredo ✓

**Decisão:** Os componentes de um segredo são chamados "campos", não "dados", "detalhes" ou "facetas".

**Contexto:**
- Alternativas consideradas: "dado" (genérico), "detalhe" (diminui semanticamente dados críticos), "faceta" (sugere perspectiva, não composição)
- "Campo" já é parte do vocabulário cotidiano do usuário (qualquer um que preencheu formulário entende)

**Justificativa:** Familiaridade universal reduz fricção de onboarding.



## 9. Modelos pré-definidos e personalizados unificados ✓

**Decisão:** Não há distinção conceitual entre modelos pré-definidos e modelos personalizados — ambos são "modelos".

**Contexto:**
- Distinção é de origem (auto-criado vs. usuário-criado), não de comportamento
- Na prática, ambos funcionam da mesma forma
- Dois termos adicionava complexidade ao vocabulário sem benefício

**Justificativa:** Um conceito simplifica semântica; modelos auto-criados são apenas ponto de partida, editáveis e removíveis como qualquer outro.



## 10. Campo "Observação" especial automático ✓

**Decisão:** Todo segredo tem um campo "Observação" especial, criado automaticamente, não deletável e não renomeável.

**Contexto:**
- Usuário precisa de espaço para notas sem ocupar campo customizado do modelo
- Deve estar sempre disponível, independentemente do modelo
- Campo comum (sempre visível), protegido contra exclusão acidental

**Justificativa:** UX simples — lugar previsível para notas; sem complexidade conceitual extra (é um campo como outro, só com restrições especiais).

**Consequências:** Sempre há destino para notas; proteção contra exclusão acidental entendida naturalmente pelo usuário.



# BUSCA & INDEXAÇÃO

## 11. Campos sensíveis não participam da busca ✓

**Decisão:** Busca ocorre apenas sobre campos comuns e observações — campos sensíveis são excluídos do índice.

**Contexto:**
- Manter sensíveis descriptografados em memória durante operação de busca aumenta superfície de exposição
- Qualquer indexação exigiria chaves em memória

**Justificativa:** Reduz superfície de risco criptográfico; compromisso aceitável porque usuário raramente busca por valores sensíveis (de memória).

**Consequências:** Usuário busca por rótulos/observações, não por conteúdo de senhas; sensível permanece encriptado sempre.



## 12. Importação: pastas mesclam sempre, segredos sincronizam, modelos substituem ✓

**Decisão:** Na importação:
- Pastas: sempre mesclam — se subpasta importada tem o mesmo caminho (mesma pasta pai + mesmo nome), seu conteúdo é consolidado com a pasta existente
- Segredos: sincronização completa na pasta — mesma identidade (mesma pasta pai + mesmo nome) = **sobrescrita de campos**; identidade única = inserido; ausente na importação = **marcado para exclusão**
- Modelos: mesmo nome = **substituição silenciosa**

**Contexto:**
- Estrutura de pastas raramente conflita e merge é sempre seguro (aditivo)
- Segredos são dados críticos — importação atua como sincronização: arquivo importado é fonte de verdade
- Subpastas **não** são renomeadas porque pasta é container primário — merge é a operação natural
- Ausência de segredo na importação sinaliza que deve ser removido (sincronização completa)
- Sem diálogo de merge torna importação simples e previsível

**Justificativa:** Assimetria deliberada reflete papel de cada entidade: pastas são containers (mesclam sempre), segredos são dados (sincronizam por nome — criam, atualizam, deletam), modelos são templates (substituem). Importação de segredos não é "merge parcial", é **sincronização**: o cofre fica igual ao arquivo importado na estrutura de pastas selecionada.

**Consequências:** 
- Pastas importadas sempre consolidam com existentes (não há rejeição)
- Segredos podem ser criados, atualizados ou marcados para exclusão via importação
- Arquivo importado é tratado como "estado desejado" — diferenças são sincronizadas
- Modelos podem ser distribuídos/atualizados via importação



## 13. Busca: substring, case-insensitive, sem acentuação ✓

**Decisão:** Busca funciona por substring, case-insensitive, ignorando acentuação. Sem suporte a regex.

**Contexto:**
- Expectativa padrão do usuário: "pass" encontra "password", "café" encontra "cafe"
- Case-insensitive reduz fricção sem comprometer segurança (busca é só em campos comuns)
- Regex adiciona complexidade para usuário e implementação

**Justificativa:** Simplicidade + usabilidade são equilibradas.

**Consequências:** Busca mantém campos sensíveis sempre encriptados; reduz fricção de busca sem sacrificar segurança.



# SEGURANÇA & CRIPTOGRAFIA

## 14. Salvamento/descarte usam senha fornecida ao abrir ✓

**Decisão:** Ao abrir cofre, usuário fornece senha mestra uma única vez. Essa senha é usada para todas operações cripto (salvar, descartar, alterações) durante a sessão.

**Contexto:**
- Alternativa de re-pedir senha a cada operação criaria UX tedioso
- Alternativa de chaves derivadas aumentaria complexidade criptográfica para v1
- Senha permanece em memória enquanto cofre está em uso

**Justificativa:** UX fluido reduz pontos de falha; fluxo de trabalho sem interrupções.

**Consequências:** 
- Proteção: mlock/VirtualLock para impedir swap; limpeza agressiva ao bloquear (sobrescrever buffer com zeros)
- Risco residual mitigado por timeout automático + bloqueio manual
- Trade-off deliberado: segurança prática > segurança teórica perfeita para v1



## 15. Proteção de senha em memória ✓

**Decisão:** Senha mestra é mantida em memória durante toda sessão (não é descartada após abrir ou entre operações).

**Contexto:**
- Alternativa (re-pedir senha) criaria UX tedioso; alternativa (chaves derivadas) seria complexa demais para v1
- Idealmente, senha nunca deveria estar em memória
- Compromisso pragmático: proteção técnica + UX aceitável

**Justificativa:** UX fluido reduz pontos de falha (menos oportunidades de senha incorreta); fluxo sem interrupções.

**Consequências:** 
- Proteção implementada: mlock/VirtualLock para impedir swap; limpeza agressiva ao bloquear (sobrescrever buffer com zeros)
- Risco residual (exposição durante sessão) mitigado por timeout automático + bloqueio manual
- **Trade-off deliberado: segurança prática > segurança teórica perfeita para escopo de v1**



## 16. Proteção contra força bruta: confiança em criptografia ✓

**Decisão:** Não implementamos rate limiting ou lockout para tentativas de autenticação (senha mestra). Proteção contra força bruta vem apenas da configuração do algoritmo criptográfico.

**Contexto:**
- Será usado algoritmo de derivação de chave computacionalmente custoso (Argon2, PBKDF2 com iterações elevadas)
- Cada tentativa é tão cara que força bruta é economicamente inviável
- Rate limiting seria apenas proteção "teatral" se criptografia já torna tentativa impossível

**Justificativa:** 
- Proteção por matemática é mais confiável que por delay de aplicação
- Evita complexidade de rastreamento estado de tentativas
- Confiança reside na força do algoritmo, não em proteções de código

**Consequências:** 
- Simplicidade implementação
- Segurança depende absolutamente da qualidade da criptografia (sem "fallback" técnico)
- Configuração incorreta de Argon2/PBKDF2 seria crítica



# ARQUIVOS & CONCORRÊNCIA

## 17. Acesso concorrente: detecção sem lock file ✓

**Decisão:** Não usamos arquivo de lock (.lock) para sincronização. Detectamos modificações externas comparando timestamp e tamanho do arquivo no momento do salvamento.

**Contexto:**
- Lock files deixariam rastro no SO, violando privacidade/portabilidade
- Alternativa é comparação de metadados + confirmação do usuário
- Conflitos concorrentes são raros em uso normal

**Justificativa:** Privacidade total (nenhum rastro de sistema); usuário tem controle explícito em conflito.

**Consequências:** 
- Se arquivo foi modificado externamente: pedimos confirmação, oferecemos "Salvar como novo arquivo"
- Trade-off: requer user decision em conflito (não merge automático)
- Simplificação deliberada: segurança dados > sofisticação técnica



## 18. Diagnóstico de corrupção: mensagem opaca ✓

**Decisão:** Se arquivo não pode ser aberto (magic number inválido, CRC falho, JSON corrupto, etc), sistema exibe mensagem genérica e opaca.

**Contexto:**
- Revelar qual parte falhou vaza informação sobre estrutura criptográfica
- Atacante poderia usar detalhes para inferir padrões de criptografia
- Mensagem genérica oferece privacidade total

**Justificativa:** Segurança por obscuridade não é ideal em geral, mas aqui tem propósito claro — impede vazamento de estrutura.

**Consequências:** Mensagem: "Arquivo não pode ser aberto — possível corrupção ou arquivo inválido." Sem diagnóstico técnico ao usuário.



## 19. Pasta Geral não-recuperável ✓

**Decisão:** Se Pasta Geral está ausente ou corrompida no arquivo, sistema rejeita o arquivo com mensagem de erro opaca (sem tentar recriar).

**Contexto:**
- Pasta Geral é invariante estrutural — raiz da árvore, destino garantido para segredos órfãos
- Ausência sinaliza corrupção ou manipulação intencional
- Seria possível recriar vazia, mas isso entraria em conflito com filosofia de segurança

**Justificativa:** Falha segura é rejeitar, não tentar "consertar" silenciosamente. Tentar "reparar" mascararia problemas reais — usuário não perceberia alteração não autorizada.

**Consequências:** Arquivo é rejeitado com erro opaco; sinaliza ao usuário que arquivo não é confiável.



## 20. Reordenação: estado final, não histórico ✓

**Decisão:** Quando usuário reordena segredos ou pastas múltiplas vezes antes de salvar, sistema persiste apenas ordem final (descarta histórico de movimentos).

**Contexto:**
- Rastrear cada movimento serviria apenas para undo/redo (fora do escopo v1)
- Histórico complica merge em acesso concorrente
- Overhead de dados desnecessário para v1

**Justificativa:** Simplifica implementação; ordem final é determinística e suficiente.

**Consequências:** Sem undo/redo; sem reconstrução de intenção de movimentos; clareza total sobre resultado ao salvar.



# DESIGN & RESTRIÇÕES

## 21. Sem limites técnicos arbitrários ✓

**Decisão:** Sistema não impõe limites arbitrários (quantidade de segredos, profundidade de pastas, campos por modelo).

**Contexto:**
- Limites técnicos surgem apenas onde HW/SO os impõe
- Um cofre com 100k segredos é tecnicamente possível mas não é caso de uso esperado
- Validações desnecessárias complicam implementação

**Justificativa:** Confia em racionalidade do usuário; caso real de abuso surgir, limite será documentado explicitamente.

**Consequências:** Experiência governa-se por bom senso, não validação; risco de UX degradada em casos extremos, mas protege contra atrito normal.



# ESCOPO v1

## 22. `CampoSegredo.valor` representado como `[]byte` em memória para todos os campos ✓

**Decisão:** O atributo `valor` de `CampoSegredo` é representado como `[]byte` em memória — para campos `texto` e `texto_sensivel` igualmente — em vez de `string`. No JSON persistido, o valor continua sendo serializado como string UTF-8 legível, via `MarshalJSON`/`UnmarshalJSON` customizados em `CampoSegredo`.

**Contexto:**
- A arquitetura exige zeragem explícita de dados sensíveis ao bloquear/encerrar (sobrescrever buffer com zeros)
- `string` em Go é imutável — uma vez criada, a memória não pode ser sobrescrita pelo programa; cópias podem existir em heap locations inaccessíveis
- `[]byte` é mutável e pode ser zerado com `for i := range b { b[i] = 0 }`
- A alternativa de dois campos distintos (`valor_texto string` + `valor_sensivel []byte`) foi considerada e descartada:
  - O `tipo` do campo já é o discriminador comportamental; adicionar dois campos de valor cria redundância estrutural
  - A invariante "preencha apenas o campo correspondente ao tipo" ficaria no Manager, não na estrutura — aumenta superfície de erro
  - Zeragem seria assimétrica (só `valor_sensivel`) — campos `texto` em `string` persistiriam em memória sem possibilidade de limpeza
- Campos `texto` (não-sensíveis) também recebem `[]byte` para uniformidade: estrutura única, zeragem uniforme, sem lógica condicional no caminho de limpeza
- `encoding/json` serializa `[]byte` automaticamente como Base64, o que quebraria compatibilidade e legibilidade do arquivo — por isso `CampoSegredo` implementa marshal/unmarshal customizados que tratam `valor` como string UTF-8 no JSON

**Justificativa:** Uniformidade elimina lógica condicional na zeragem; `[]byte` é pré-requisito para a promessa de limpeza de memória ser honrável; o custo (marshal customizado + conversões `string(b)`/`[]byte(s)` nos pontos de exibição) é pequeno e bem localizado.

**Consequências:**
- `CampoSegredo` requer implementação de `MarshalJSON`/`UnmarshalJSON` — não pode usar serialização padrão
- Conversões `string(campo.valor)` necessárias nos pontos de exibição na TUI
- Zeragem ao bloquear/encerrar percorre todos os campos sem distinguir tipo — código mais simples
- O JSON gravado em disco é idêntico ao que seria com `string` — compatibilidade de formato preservada



## 23. Lixeira como lista in-memory no Manager, sem campo `marcado_exclusao` no Segredo ✓

**Decisão:** `Segredo` não possui campo `marcado_exclusao` nem qualquer flag de estado de exclusão. A Lixeira é uma lista in-memory separada, mantida exclusivamente pelo Manager durante a sessão. Ao salvar, os segredos da Lixeira são descartados — não são persistidos.

**Alternativa descartada — campo no Segredo:**
- Adicionar `marcado_exclusao bool` ao `Segredo` tornaria o estado de exclusão parte da entidade de domínio
- Segredo passaria a ter estado misto (dados + estado operacional transitório), violando o princípio de entidade de domínio puro
- Criaria risco de o campo ser acidentalmente persistido no JSON
- Zeragem de memória ficaria mais complexa (campo extra a limpar)

**Alternativa escolhida — lista separada no Manager:**
- `Segredo` permanece imutável como entidade de domínio puro: nome, campos, favorito, datas
- Manager mantém internamente um conjunto (slice ou set) de segredos pendentes de exclusão
- O ciclo de vida da Lixeira é estritamente atrelado ao ciclo de vida do Manager — ao bloquear sem salvar, os segredos da Lixeira são zerados junto com os demais dados sensíveis
- Mais coeso: decisão de exclusão é operação do Manager, não estado da entidade

**Implicação de API obrigatória (deve ser definida antes de implementar):**
A TUI precisa saber se um segredo está na Lixeira para exibi-lo com marcação visual adequada (e para oferecer as ações de restaurar / confirmar exclusão). O Manager deve expor isso explicitamente — por exemplo:
- um método `IsNaLixeira(nomePasta, nomeSegredo string) bool`, ou
- incluir os segredos da Lixeira em um tipo de retorno anotado (ex: `SegredoComEstado`), ou
- expor a Lixeira como coleção navegável distinta

A forma exata da API é decisão de implementação, mas a interface do Manager deve resolê-la **antes** de qualquer código de TUI que renderize a Lixeira.

**Consequências:**
- `Segredo` não carrega estado transitório — domínio permanece limpo
- Lixeira não sobrevive a bloqueio sem save — comportamento esperado e desejado
- Zeragem ao bloquear deve incluir explicitamente os segredos da Lixeira
- A interface do Manager precisa de um contrato deliberado para exposição do estado de Lixeira à TUI



## 24. Segredos e pastas em listas separadas dentro de Pasta — blocos fixos, pastas antes dos segredos ✓

**Decisão:** `Pasta` mantém duas listas independentes: `pastas` e `segredos`. Na exibição, pastas aparecem sempre em bloco antes dos segredos. Reordenação é permitida dentro de cada bloco, mas não é possível intercalar um segredo entre duas pastas nem uma pasta entre dois segredos.

**Contexto:**
- A estrutura de domínio (`Pasta.pastas` + `Pasta.segredos` como listas separadas) é uma decisão de modelo que implica diretamente a UX de ordenação — essa implicação não estava explicitada nos requisitos
- Alternativa descartada: lista única com posição global (segredos e pastas intercalados por posição)
  - Exigiria campo `posição` ou ordem por índice numa lista única de tipo polimórfico
  - Aumenta complexidade de serialização e de reordenação no Manager
  - A UX de arrastar uma pasta para entre dois segredos (ou vice-versa) não traz benefício claro para cofres de senhas pessoais
- Blocos fixos (pastas antes de segredos) é o padrão de navegadores de arquivos e gerenciadores de senhas — expectativa natural do usuário

**Justificativa:** Listas separadas são mais simples de implementar, serializar e raciocinar. A restrição de blocos é transparente ao usuário — ele nunca tentará colocar uma pasta entre dois segredos porque a interface simplesmente não oferece essa opção.

**Consequências:**
- Reordenação de pastas afeta apenas `Pasta.pastas`; reordenação de segredos afeta apenas `Pasta.segredos` — sem conflito entre as duas operações
- A TUI deve renderizar sempre o bloco de pastas antes do bloco de segredos dentro de cada nível da hierarquia
- Não é possível representar uma ordem interleaved, mesmo que um futuro formato de arquivo quisesse — a estrutura impõe o modelo



## 25. Auditoria de senhas e TOTP fora do escopo v1 ✓

**Decisão:** Funcionalidades de auditoria de senhas e geração TOTP foram descartadas do escopo.

**Contexto:** Identificadas durante processo de definição, mas estão fora do foco inicial do produto.

**Consequências:** Podem ser adicionadas em iteração futura; v1 foca em cofre core funcional.

