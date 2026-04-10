# Design System — Abditum TUI

## Anti-padrões

Anti-padrões documentam o que **não deve ser feito** na interface do Abditum. Cada item lista o padrão incorreto, por que viola os princípios do design system, e qual consequência concreta afeta o usuário.

> **Regra de uso:** toda decisão de implementação que contradiga um item desta seção deve ser justificada explicitamente na especificação que a adota. Sem justificativa documentada, o anti-padrão prevalece como proibição.

### Segurança Visual

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Revelação Passiva de Sensível** *(Crítico)* | Campos sensíveis revelam por foco/Tab sem ação explícita | Dado sensível exposto sem percepção do usuário |
| **Máscara Apenas Visual** *(Alto)* | `••••••••` exibido mas copiável sem feedback | Proteção ilusória; dado sensível exposto em clipboard |
| **Campo Sensível Indistinguível** *(Alto)* | Campos sensíveis e comuns têm mesma aparência | Revelação acidental ou proteção ignorada |
| **Countdown Invisível** *(Médio)* | Cópia bem-sucedida mas sem indicação de TTL da clipboard | Usuário não sabe se o dado ainda está disponível |
| **Encerramento Sem Limpeza** *(Crítico)* | Ao encerrar (gracioso, lock ou sinal SIGINT/SIGQUIT), app não garante limpeza de clipboard (goroutine pode ser interrompida antes de `os.Exit`), scrollback (apenas tela visível é limpada) e memória sensível. Ver [`arquitetura.md` § Clipboard e § Clear screen](arquitetura.md) | Dados sensíveis permanecem em clipboard, histórico do terminal e memória após encerrar ou bloquear |
| **Exportação Sem Cerimônia** *(Crítico)* | Exportação (arquivo não criptografado) com tratamento de ação rotineira | Usuário exporta para local inseguro sem compreender risco |
| **Dirty State Apenas Global** *(Crítico)* | Indicador `•` só no cabeçalho, sem `✦ ✎ ✗` por item | Usuário não consegue auditar o que será salvo |

### Estado e Feedback

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Silêncio Após Operação Crítica** *(Alto)* | Salvar/bloquear/exportar sem mensagem de confirmação | Indistinguível de falha silenciosa |
| **Spinner Sem Resolução** *(Médio)* | `◐ Carregando…` nunca substituído por `✓` ou `✕` | Usuário não sabe se pode interagir |
| **Fila de Mensagens** *(Médio)* | Múltiplas mensagens enfileiradas ou sobrepostas | Falta correspondência entre ação e mensagem |
| **Contador Defasado** *(Médio)* | Contagem de segredos não atualiza em tempo real | Decisões baseadas em dados incorretos |
| **Modo Ativo Sem Indicador** *(Alto)* | Busca/edição/reordenação sem indicador persistente na barra | Usuário digita no modo errado |

### Navegação e Teclado

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Enter Polissêmico** *(Médio)* | `Enter` expande ou edita sem distinção visual clara | Edição acidental ao tentar visualizar |
| **Cursor ao Topo Após Operação** *(Médio)* | Exclusão/reordenação retorna cursor ao topo da lista | Re-navegação obrigatória; experiência frustrante |
| **Setas com Semântica Dupla Invisível** *(Baixo)* | `←`/`→` expandem pastas E navegam diálogos sem indicador | Expansão/fechamento/navegação acidental |
| **Atalho de Tecla Simples em Painel de Lista** *(Alto)* | Atalho sem modificador (letra ou dígito isolado) registrado em painel de árvore ou lista — ex: `n` para novo, `e` para editar | Impede type-to-search e qualquer modo de entrada futura no painel; usuário aciona ação involuntária ao tentar digitar uma query |

### Diálogos e Confirmações

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Fadiga de Confirmação** *(Médio)* | Toda ação pede confirmação, inclusive benignas | Usuário aprende a apertar Enter reflexivamente |
| **Uniformidade de Risco Visual** *(Alto)* | "Excluir 47 segredos" e "Renomear pasta" têm mesma aparência | Usuário não calibra gravidade da ação |
| **Pilha de Modais Sem Profundidade** *(Médio)* | Modal abre modal abre modal sem indicação | Desorientação; fechamento acidental com Esc repetido |
| **Ação Default Ausente** *(Médio)* | Ação default desaparece quando inativa | Usuário não sabe o que falta preencher |
| **Confirmação Assimétrica** *(Crítico)* | "Salvar e Sair" pede dupla confirmação; "Descartar" não | Incentivo perverso aumenta perdas de dados |
| **Ação Destrutiva/Irreversível Sem Confirmação** *(Crítico)* | Qualquer ação com perda de dados executa sem confirmação e opção de desistir — inclui: encerrar app com mudanças pendentes, fechar/bloquear cofre com mudanças pendentes, excluir itens | Dado perdido sem chance de recuperação |
| **Fluxo Sem Saída** *(Alto)* | Fluxo de múltiplos passos sem opção de desistir (cancelar) ou voltar ao passo anterior | Usuário preso; forçado a concluir ou matar o processo |
| **Borda como Menu** *(Alto)* | Diálogo com 4 ou mais ações na borda inferior — transforma a moldura num "roteador" de escolhas | Sobrecarga cognitiva; usuário não sabe qual ação escolher; a borda perde o papel de confirmação e vira menu disfarçado |
| **Fantasma na Barra** *(Alto)* | Exibir na barra de comandos uma ação que não se aplica ao contexto atual (`Enabled = false`) ou uma tecla que não está registrada/não funciona no terminal (ex: F17 quando só existem 12 F-keys) | Usuário tenta acionar e nada acontece ou recebe erro silencioso; quebra a confiança na barra como indicador de ações disponíveis |
| **Dica com Tecla Redundante** *(Médio)* | Mensagem de dica de campo ou uso menciona teclas de ação (ex: "F17 para copiar") que já estão visíveis na barra de comandos | Ruído cognitivo; dica fica datada se a tecla mudar; duplicação de informação |
| **Tecla Ambígua** *(Alto)* | Atribuir ao aplicativo combinações que o terminal ou shell intercepta antes de entregar (ex: `Ctrl+C` → SIGINT, `Ctrl+Z` → EOF/SIGTSTP, `Ctrl+D` → EOF) | Ação nunca chega ao app; comportamento imprevisível entre plataformas; pode encerrar a aplicação |

### Layout e Estrutura

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Layout Saltitante** *(Médio)* | Elementos fixos reposicionam por conteúdo variável | Perda de ancoragem após cada seleção |
| **Over-boxing** *(Baixo)* | Toda seção envolta em borda; grade de boxes | Ruído visual; ambiguidade estrutural/decoração |
| **Informação Densa** *(Baixo)* | Nomes, labels, valores com mesmo peso tipográfico | Dificulta localização e varredura rápida |
| **Truncamento Ausente** *(Alto)* | Texto longo cortado sem `…` ou transborda | Confusão de identidade; layout corrompido |
| **Pasta Virtual Indistinguível** *(Médio)* | Favoritos parecem pastas normais | Usuário tenta criar item e recebe erro inesperado |
| **Caractere Largura Dupla** *(Alto)* | Símbolos ambíguos sem contabilização de colunas | Bordas não fecham; separadores desalinhados |
| **Resize Sem Recálculo** *(Crítico)* | Layout não atualiza ao redimensionar terminal | Interface inutilizável até reiniciar |
| **Conteúdo Sem Scroll** *(Alto)* | Painel/diálogo corta conteúdo sem `↑↓` nem thumb | Campos/ações finais inacessíveis |
| **Campo Maior que Área** *(Médio)* | Valor longo sem truncamento/scroll horizontal | Sobrescrita de labels; valor ilegível |
| **Sangramento ANSI** *(Alto)* | Estilo não resetado contamina conteúdo seguinte | Cores/estilos vaza para componentes vizinhos e shell |
| **Cálculo de Largura Errado** *(Alto)* | `len(s)` em vez de largura visual (ANSI excluído) | Desalinhamento de bordas e truncamento |
| **Layout Colapsa Vazio** *(Médio)* | Painel sem conteúdo desaparece | Separador desaparece; proporção quebra ao preenchimento |
| **Indicador Causa Deslocamento** *(Médio)* | `✦ ✎ ★` não em coluna fixa | Nomes "pulam" horizontalmente ao marcar/desmarcar |
| **Spinner com Largura Variável** *(Baixo)* | Frames do spinner ocupam 1 ou 2 colunas | Mensagem pisca horizontalmente |
| **Contador com Largura Dinâmica** *(Baixo)* | Número muda de 9 para 10 dígitos | Coluna inteira se desloca |
| **Artefato de Render Anterior** *(Alto)* | Caracteres/cores do frame antigo permanecem | Bordas flutuam; campos extras visíveis |
| **Última Linha Causa Scroll** *(Médio)* | Escrever em `(linhas, colunas)` aciona scroll | Barra de comandos "cai"; layout deslocado |
| **Cursor Desalinhado** *(Alto)* | Cursor em coluna errada durante edição (bytes vs runes) | Backspace apaga caractere errado |
| **Campo Edição Sem Scroll H** *(Alto)* | Campo longo truncado ou overflow sem scroll | Usuário não vê valor completo |
| **Scroll Mal Implementado** *(Alto)* | Painel com scroll não responde a `Home`/`End`/`PgUp`/`PgDn`; não responde ao scroll do mouse; ou `PgUp`/`PgDn`/`Home`/`End` funcionam invertidos | Navegação lenta, inacessível ou desorientadora |

### Tipografia e Cor

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Bold Inflacionado** *(Baixo)* | Bold aplicado a tudo (títulos, labels, nomes, ações) | Hierarquia colapsa; interface gritante |
| **Token Semântico Decorativo** *(Baixo)* | `semantic.success` / `warning` usado para ornamento | Usuário para de confiar nos indicadores |
| **Cor Hardcoded** *(Alto)* | Hex literal em vez de tokens de tema | Segundo tema nunca funciona corretamente |
| **Italic Sem Cor** *(Baixo)* | Hints em italic apenas, sem `text.secondary` | Indistinguível do conteúdo em terminais sem italic |

### Acessibilidade

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Estado Apenas por Cor** *(Alto)* | `✦ ✎ ✗` não usados; apenas color diferencia | Em NO_COLOR, informação de estado desaparece |
| **Área de Clique Minúscula** *(Baixo)* | `<╡` (2 colunas dinâmicas) único alvo clicável | Mouse inutilizável; cada clique requer precisão |
| **Erro Técnico Exposto** *(Médio)* | Mensagens internas: "unexpected JSON at 1247" | Exposição de caminho de arquivo; usuário confuso |

### Ciclo de Vida do Cofre

| Anti-padrão | Problema | Impacto |
|---|---|---|
| **Auto-save Silencioso** *(Alto)* | Alteração de senha salva automaticamente sem feedback | Usuário acredita que é reversível |
| **Conflito de Arquivo Minimizado** *(Crítico)* | Arquivo modificado externamente sobrescrito sem aviso | Dados de outra sessão/backup destruídos |
| **Re-autenticação Durante Sessão** *(Alto)* | Senha mestra solicitada novamente em operações rotineiras de salvamento ou descarte — a senha já está em memória e não há ganho de segurança nisso. **Exceção legítima:** exportação, onde re-autenticação é um controle válido de defesa-em-profundidade perante um cofre aberto sem dono | Fricção ilegítima nos fluxos comuns; treina o usuário a fornecer a senha sem questionar o contexto — vetor clássico de phishing de UI |
| **Exclusão Desaparece Imediatamente** *(Crítico)* | Item marcado para exclusão some sem `✗` + strikethrough | Usuário crê ter deletado permanentemente |
| **Importação Sem Prévia de Impacto** *(Crítico)* | Mesclagem executada sem mostrar o que será sobrescrito | Perda de dados não intencionada |


