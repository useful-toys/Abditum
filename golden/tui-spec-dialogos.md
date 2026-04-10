# Especificação Visual — Diálogos

> Anatomia comum, tipos de diálogo e catálogo de diálogos de decisão.
> Parte da [Especificação Visual](tui-specification.md).
>
> **Documentos relacionados:**
> - [`tui-design-system.md`](tui-design-system.md) — fundações visuais
> - [`tui-spec-dialog-senha.md`](tui-spec-dialog-senha.md) — PasswordEntry, PasswordCreate
> - [`tui-spec-dialog-filepicker.md`](tui-spec-dialog-filepicker.md) — FilePicker
> - [`tui-spec-dialog-help.md`](tui-spec-dialog-help.md) — Ajuda

## Diálogos

Diálogos são janelas sobrepostas que capturam o foco da aplicação para uma interação isolada — uma decisão, um reconhecimento, uma entrada de dados ou uma consulta de referência. Enquanto um diálogo estiver aberto, a área de trabalho permanece visível porém inativa.

A aplicação utiliza quatro tipos de diálogo, cada um com propósito, anatomia e regras de comportamento distintos:

| Tipo | Propósito | Seção |
|---|---|---|
| [Notificação](#notificação) | Informar um fato que exige reconhecimento | [▸](#notificação) |
| [Confirmação](#confirmação) | Solicitar uma escolha explícita do usuário | [▸](#confirmação) |
| [Ajuda](#ajuda) | Exibir referência de atalhos (somente leitura) | [▸](#ajuda) |
| [Funcional](#funcional) | Capturar entrada de dados com campos interativos | [▸](#funcional) |

### Anatomia Comum

Todo diálogo é composto por três regiões estruturais, desenhadas com bordas arredondadas (`╭╮╰╯│─`):

```text
╭── <símbolo>  <título> ─────────────────────╮  ← Borda Superior
│                                            │
│  <conteúdo do corpo>                       │  ← Corpo
│                                            │
╰── <Tecla Label> ──── <Tecla Label> ────────╯  ← Borda de Ações
```

**Borda Superior (título):**
- Estrutura char a char:
  - **Sem símbolo:** `╭──` (canto + 2× borda) + ` ` (1 espaço) + título + ` ` (1 espaço) + preenchimento `─`×N + `╮` (canto).
  - **Com símbolo:** `╭──` (canto + 2× borda) + ` ` (1 espaço) + símbolo + `  ` (2 espaços) + título + ` ` (1 espaço) + preenchimento `─`×N + `╮` (canto).
- O título ocupa a borda superior a partir da 5ª coluna (após `╭── `), preservando os caracteres de canto de ambos os lados. O preenchimento `─` garante pelo menos 1 caractere de borda antes do `╮`.
- Quando a severidade é Neutro ou o tipo de diálogo não usa severidade (Ajuda, Funcional), o símbolo é omitido.
- **Truncamento:** se o título excede o espaço disponível na borda (largura máxima do diálogo − cantos − espaçamento), ele é truncado com `…`.
- O título descreve o fluxo ou ação principal (ex: `Salvar cofre`, `Senha mestra`, `Ajuda`). Capitalizado conforme o nome, sem artigos desnecessários.

**Corpo:**
- Bordas laterais `│` delimitam o conteúdo.
- Padding interno: 2 colunas horizontais (entre `│` e o texto). Padding vertical de 1 linha (acima e abaixo do conteúdo) aplica-se a Notificação, Confirmação e Ajuda; diálogos Funcionais usam **0 linhas** de padding vertical — o conteúdo denso e interativo ocupa todo o espaço disponível.
- O conteúdo varia por tipo de diálogo: texto estático (Notificação, Confirmação), tabela de atalhos (Ajuda) ou campos interativos (Funcional).
- Diálogos Funcionais podem conter **divisores internos** que segmentam o corpo em regiões. A separação pode ser horizontal, vertical ou ambas:
  - **Horizontal:** `─` conectado às bordas laterais por T junctions (`├` à esquerda, `┤` à direita).
  - **Vertical:** `│` conectado às bordas horizontais por T junctions (`┬` no topo, `┴` na base).
  - **Cruzamento:** `┼` onde um divisor horizontal e um vertical se cruzam.

**Borda de Ações (rodapé):**
- Estrutura char a char: `╰` (canto) + `─` (1× borda) + ações + preenchimento `─`×N (pelo menos 1) + `╯` (canto).
- Cada ação é representada como: ` ` (espaço) + tecla + ` ` (espaço) + label + ` ` (espaço). Exemplo: ` Enter Salvar `.
- Ações são separadas entre si por segmentos de preenchimento `─`.
- Layout varia conforme a quantidade de ações:
  - **1 ação:** alinhada à direita. Preenchimento `─` ocupa todo o espaço à esquerda.
    ```text
    ╰──────────────────────────── Enter OK ─╯
    ```
  - **2 ações:** principal à esquerda, cancelamento à direita. Preenchimento `─` entre elas.
    ```text
    ╰─ Enter Confirmar ──────── Esc Cancelar ─╯
    ```
  - **3 ações:** principal à esquerda, secundária ao centro, cancelamento à direita. Preenchimento `─` distribuído entre elas.
    ```text
    ╰─ S Sobrescrever ── N Como novo ── Esc Voltar ─╯
    ```
- Limite máximo: 3 ações. Diálogos com 4+ ações na borda são um [anti-padrão](tui-design-system-anti-patterns.md#diálogos-e-confirmações) ("Borda como Menu").

### Apresentação e Pilha

- O diálogo centraliza-se horizontal e verticalmente sobre a tela inteira.
- O conteúdo abaixo permanece visível, mas inativo (sem escurecimento de overlay).
- Apenas o elemento do topo da pilha recebe input; os inferiores permanecem montados, porém congelados.
- Ao fechar o elemento do topo, o foco retorna ao elemento imediatamente anterior na pilha.

### Dimensionamento

- **Largura mínima:** suficiente para acomodar o título e as ações da borda inferior, ou no mínimo 20 colunas.
- **Largura máxima:** até 95% da largura do terminal.
- **Largura fixa:** diálogos funcionais específicos podem definir largura fixa (ex: PasswordEntry = 50 colunas). A largura fixa é documentada na especificação de cada subtipo.
- **Altura:** determinada pelo contorno do conteúdo, sem espaços vazios exagerados.
- **Padding interno:** 2 colunas laterais. Padding vertical de 1 linha para Notificação, Confirmação e Ajuda; **0 linhas** para Funcional.

### Scroll

Quando o conteúdo do corpo excede o espaço disponível:

- **Largura excedida:** word-wrap quebra linhas mantendo integridade de palavras.
- **Altura excedida:** ativa-se scroll vertical com indicadores visuais na borda lateral direita.
- A borda superior e a borda de ações nunca participam do scroll — permanecem sempre fixas.
- **Pré-condição:** o scroll só é renderizável se o diálogo possuir pelo menos 5 linhas internas (excluindo borda superior e borda de ações). Abaixo desse mínimo, o conteúdo é truncado sem indicadores de scroll.
- **Navegação por teclado:** teclas direcionais (`↑`/`↓`), `PgUp`/`PgDn`, `Home`/`End`.

**Composição da borda lateral direita com scroll:**

A borda lateral direita do corpo é composta por 3 elementos, cada um ocupando posições fixas:

| Elemento | Posição | Caractere | Descrição |
|---|---|---|---|
| Seta superior | 1ª linha do corpo | `↑` | Indica conteúdo acima do viewport. Substitui o `│` da borda |
| Thumb | Entre a 2ª e a penúltima linha do corpo | `■` | Posição relativa do viewport no conteúdo total. Nunca sobrepõe as setas |
| Seta inferior | Última linha do corpo | `↓` | Indica conteúdo abaixo do viewport. Substitui o `│` da borda |

- As setas `↑` e `↓` ocupam sempre a primeira e a última linha do corpo, respectivamente.
- O thumb `■` é posicionado proporcionalmente entre a 2ª linha e a penúltima linha do corpo — ele **nunca** é desenhado sobre a posição de uma seta.
- Nas linhas onde nenhum indicador está presente, a borda permanece `│`.

Wireframe ilustrando o scroll ativo (5 linhas internas, com padding vertical):

```text
╭── ⚠  Título do Diálogo ────────────────────╮
│                                            ↑
│  Primeira linha do conteúdo longo.         ■
│  Segunda linha mostrando limite excedido.  │
│  Terceira linha com mais informações.      │
│                                            ↓
╰── Enter Salvar ── A Alt ──── Esc Cancelar ─╯
```

### Identidade Visual

Regras visuais padrão aplicadas a **todos** os diálogos. Cada tipo documenta apenas as variações.

> Caracteres estruturais: ver [Anatomia Comum](#anatomia-comum).

| Elemento | Token | Atributo | Observação |
|---|---|---|---|
| Bordas e cantos | Determinado pela severidade ou pelo tipo | — | Notificação/Confirmação: token da [Severidade](#severidade). Ajuda: `border.default`. Funcional: `border.focused` |
| Símbolo na borda superior | Determinado pela severidade | — | `⚠`, `✕` ou `ℹ` conforme severidade. Omitido em Neutro, Ajuda e Funcional |
| Título | `text.primary` | **bold** | Descreve o fluxo ou ação principal |
| Texto do corpo | `text.primary` | — | — |
| Tecla da ação default (`Enter`) | Token da tecla default da severidade | **bold** | Notificação/Confirmação: ver [Severidade](#severidade). Ajuda/Funcional: `accent.primary` |
| Teclas de ações secundárias e cancelamento | Segue o token de borda | — | — |

#### Severidade

Severidade governa o tratamento visual — borda, símbolo e cor da tecla default — aplicado exclusivamente aos diálogos de **Notificação** e **Confirmação**. Diálogos de Ajuda e Funcional não utilizam severidade.

| Severidade | Símbolo | Token de borda | Token da tecla default | Quando usar |
|---|---|---|---|---|
| Destrutivo | `⚠` | `semantic.warning` | `semantic.error` | Ação irreversível ou com perda de dados |
| Erro | `✕` | `semantic.error` | `accent.primary` | Falha ocorrida, condição irrecuperável |
| Alerta | `⚠` | `semantic.warning` | `accent.primary` | Situação importante mas recuperável |
| Informativo | `ℹ` | `semantic.info` | `accent.primary` | Informação que requer atenção |
| Neutro | — | `border.focused` | `accent.primary` | Operação rotineira, sem urgência |

> **Nota:** severidades Destrutivo e Alerta compartilham o símbolo `⚠` e o token de borda `semantic.warning`. A distinção visual está na tecla default: `semantic.error` (vermelho) para destrutivo, `accent.primary` para alerta. Isso reforça que o perigo está na *ação*, não apenas na *situação*.

### Teclado

Convenções de teclado aplicadas a todos os diálogos. Cada tipo documenta apenas as variações.

**Teclas implícitas — `Enter` e `Esc`:**

Todo diálogo possui duas teclas implícitas que não precisam ser declaradas pelas ações:

| Tecla | Papel implícito |
|---|---|
| `Enter` | Executa a ação **principal** (a da extrema esquerda, ou a única ação) |
| `Esc` | Executa a ação de **cancelamento** (a da extrema direita, ou a única ação) |

Comportamento conforme a quantidade de ações:

| Ações | `Enter` | `Esc` |
|---|---|---|
| 1 ação | Executa a ação única | Mesmo efeito que `Enter` — executa a ação única |
| 2 ações | Executa a ação da esquerda (principal) | Executa a ação da direita (cancelamento) |
| 3 ações | Executa a ação da esquerda (principal) | Executa a ação da direita (cancelamento) |

**Teclas explícitas — letras de atalho:**

Além das teclas implícitas, cada ação pode declarar uma tecla de atalho (tipicamente a primeira letra da label). A tecla declarada é exibida na borda de ações antes da label (ex: `S Sobrescrever`). As teclas implícitas `Enter` e `Esc` continuam funcionando mesmo quando a ação possui tecla explícita. Exemplo:

```text
╰─ S Sobrescrever ── N Como novo ── Esc Voltar ─╯
```

- `S` → Sobrescrever (tecla explícita)
- `Enter` → Sobrescrever (tecla implícita — ação principal)
- `N` → Como novo (tecla explícita)
- `Esc` → Voltar (tecla implícita — cancelamento)

A ação secundária (centro, quando presente) **sempre** precisa declarar sua tecla explícita — não possui tecla implícita.

**Diálogos Funcionais — exceção do `Enter`:**

Em diálogos funcionais, `Enter` pode ter comportamento contextual dependendo do campo em foco (ex: submeter um campo, selecionar um item em lista). Entretanto, em algum estado do diálogo o `Enter` **deve** acionar a confirmação do diálogo. `Esc` sempre cancela e fecha o diálogo, sem exceção.

### Notificação

**Quando usar:** o usuário precisa tomar ciência de um fato — uma falha, um alerta ou uma informação relevante. Não há decisão a tomar; apenas reconhecimento.

**Anatomia:**

| Região | Presença | Conteúdo |
|---|---|---|
| Borda Superior | Obrigatória | Símbolo de severidade + título |
| Corpo | Obrigatório | Apenas afirmação. Sem pergunta. Frases terminam com ponto final |
| Borda de Ações | Obrigatória | Exatamente 1 ação, alinhada à direita |

**Redação do corpo:** afirmação concisa e direta. Referências a itens específicos em aspas simples. Exemplos:
- `Arquivo corrompido ou inválido. Necessário fechar.`
- `Senhas não conferem. Necessário digitar novamente.`
- `Arquivo inválido ou versão não suportada. Necessário corrigir.`

**Variações Visuais:** sem variações — segue integralmente a [Identidade Visual](#identidade-visual) geral com severidade.

**Ações:**
- Borda de Ações: exclusivamente `Enter OK`, alinhada à direita.
- Nenhuma ação secundária. Nenhuma ação de cancelamento.

**Teclado:** sem variações — segue integralmente as convenções de [Teclado](#teclado) geral (1 ação: `Enter` e `Esc` ambos fecham).

**Barra de Comandos:** vazia. Ações do diálogo não se repetem na barra.

**Barra de Mensagens:** limpa durante toda a exibição do diálogo.

**Exemplo Visual:**

```text
╭── ✕  Arquivo corrompido ───────────────╮
│                                        │
│  Arquivo corrompido ou inválido.       │
│  Necessário fechar.                    │
│                                        │
╰────────────────────────────── Enter OK ╯
```

### Confirmação

**Quando usar:** o usuário precisa fazer uma escolha explícita que confirma, bifurca ou cancela um fluxo — salvar, descartar, sobrescrever, excluir.

**Anatomia:**

| Região | Presença | Conteúdo |
|---|---|---|
| Borda Superior | Obrigatória | Símbolo de severidade (quando não Neutro) + título |
| Corpo | Obrigatório | Afirmação de contexto (opcional, terminada em ponto) + pergunta objetiva (terminada em `?`) |
| Borda de Ações | Obrigatória | 2 ou 3 ações |

**Redação do corpo:** fato opcional seguido de pergunta concisa que apresenta as opções de decisão. A pergunta não menciona a opção `Voltar` (Esc). Referências a itens específicos em aspas simples. Exemplos:
- `Sair do Abditum?`
- `Cofre modificado. Salvar ou descartar?`
- `Arquivo modificado externamente. Sobrescrever?`
- `'Gmail' será excluído permanentemente. Continuar?`

**Variações Visuais:** sem variações — segue integralmente a [Identidade Visual](#identidade-visual) geral com severidade.

**Ações:**
- Borda de Ações: 2 a 3 ações.
  - Ação principal à esquerda (ex: `Enter Salvar`, `S Sobrescrever`).
  - Ação secundária ao centro, quando presente (ex: `N Salvar como novo`).
  - `Esc Cancelar` (ou `Esc Voltar`) sempre na extrema direita.
- Todas as ações ficam ativas simultaneamente — não há validação condicional.

**Teclado:** sem variações — segue integralmente as convenções de [Teclado](#teclado) geral (2–3 ações com teclas explícitas).

**Barra de Comandos:** vazia. Decisões ficam exclusivamente na borda de ações.

**Barra de Mensagens:** limpa durante toda a exibição do diálogo.

**Exemplo Visual:**

```text
╭── ⚠  Salvar cofre ─────────────────────────────╮
│                                                │
│  Arquivo modificado externamente.              │
│  Sobrescrever ou salvar como novo?             │
│                                                │
╰── S Sobrescrever ── N Como novo ─ Esc Voltar ──╯
```

### Ajuda

**Quando usar:** o usuário precisa consultar os atalhos de teclado disponíveis no contexto atual. Acionado por `F1`. Diálogo somente leitura, sem impacto no estado da aplicação.

**Anatomia:**

| Região | Presença | Conteúdo |
|---|---|---|
| Borda Superior | Obrigatória | Título `Ajuda` em **bold**, sem símbolo |
| Corpo | Obrigatório | Tabela de atalhos organizada por contexto (seções com cabeçalho). Scroll ativado quando o conteúdo excede a altura |
| Borda de Ações | Obrigatória | Exatamente 1 ação: `Esc Fechar`, alinhada à direita |

**Variações Visuais:**
- Não usa severidade. Borda em `border.default`.
- Nomes das teclas de atalho no corpo em `text.primary`; descrições em `text.secondary`.

**Ações:**
- Borda de Ações: exclusivamente `Esc Fechar`, alinhada à direita.

**Teclado:** sem variações — segue integralmente as convenções de [Teclado](#teclado) geral (1 ação: `Esc` fecha). Teclas de scroll (`↑`/`↓`, `PgUp`/`PgDn`, `Home`/`End`) ativas quando há conteúdo excedente.

**Barra de Comandos:** pode exibir ações auxiliares do contexto (ex: `F12` para troca de tema), sem repetir a ação da borda.

**Barra de Mensagens:** limpa durante toda a exibição do diálogo.

**Exemplo Visual:**

```text
╭── Ajuda ──────────────────────────────╮
│ Árvore                                ↑
│ F2       Renomear arquivo atual       ■
│ Ctrl+N   Novo arquivo no diretório    │
│ Ctrl+D   Marcar para exclusão         ↓
╰──────────────────────────── Esc Fechar ╯
```

### Funcional

**Quando usar:** o usuário precisa fornecer dados por meio de campos interativos — entrada de senha, seleção de arquivo, criação de senha com confirmação. Diferente dos diálogos de decisão, o Funcional captura input estruturado.

**Anatomia:**

| Região | Presença | Conteúdo |
|---|---|---|
| Borda Superior | Obrigatória | Título em **bold**, sem símbolo de severidade |
| Corpo | Obrigatório | Campos de entrada (`input`), labels, contadores e outros componentes interativos. Conteúdo varia por subtipo. Sem padding vertical (0 linhas acima e abaixo) |
| Divisores internos | Opcional | Separadores horizontais (`─` com `├` / `┤`), verticais (`│` com `┬` / `┴`) ou ambos (`┼` no cruzamento), segmentando o corpo em regiões distintas |
| Borda de Ações | Obrigatória | Ação de confirmação + ação de cancelamento (2 ações) |

**Variações Visuais:**
- Não usa severidade. Borda em `border.focused`.
- Labels de campo: `accent.primary` + **bold** quando o campo está ativo; `text.secondary` quando inativo.
- Área de campo de entrada: fundo `surface.input`.
- Máscara de senha: caracteres `●` em `text.secondary`, com comprimento fixo (não revela o tamanho real da senha).
- Cursor no campo ativo: `▌` em `text.primary`.
- Ação default (`Enter`): estado condicional — `text.disabled` / dim enquanto houver validações pendentes; `accent.primary` + **bold** ao satisfazer condições.

**Ações:**
- Borda de Ações: 2 ações.
  - Ação de confirmação à esquerda (ex: `Enter Confirmar`).
  - `Esc Cancelar` à direita.
- Estado condicional do `Enter` descrito em Variações Visuais acima.

**Teclado:**
- `Enter` e `Esc` conforme [Teclado](#teclado) geral, com a exceção de que `Enter` pode ter comportamento contextual por campo.
- `Tab` / `Shift+Tab`: navega entre campos (quando há múltiplos campos).
- Teclas de edição: digitação, `Backspace`, `Del` — comportamento padrão de campo de texto.

**Barra de Comandos:** exibe ações auxiliares específicas do diálogo (ex: `Tab Campo seguinte`, `Del Limpar linha`). Ações da borda de ações **nunca se repetem** na barra de comandos.

**Barra de Mensagens:** território de uso exclusivo do diálogo funcional.
- Ao abrir o diálogo: dica de campo (`•`) orientando a ação esperada.
- Durante a interação: dica atualizada conforme o campo em foco.
- Após erro de validação: mensagem de erro (`✕`) exibida até correção ou troca de campo.
- Ao fechar o diálogo: barra limpa (responsabilidade do orquestrador).

**Exemplo Visual:**

```text
╭── Alterar senha mestra ───────────────╮
│  Senha atual:   ••••••••              │
│  Nova senha:    ▌                     │
╰── Enter Confirmar ────── Esc Cancelar ╯
```

Exemplo com divisores internos (FilePicker simplificado):

```text
╭── Abrir cofre ──────────┬─────────────╮
│  📁 Documentos          │ cofre.abdt  │
│  📁 Projetos            │ notas.abdt  │
│    📄 cofre.abdt        │             │
├─────────────────────────┴─────────────┤
│  Arquivo: cofre.abdt                  │
╰── Enter Abrir ──────── Esc Cancelar ──╯
```

**Subtipos conhecidos:**

| Subtipo | Propósito | Campos | Referência |
|---|---|---|---|
| PasswordEntry | Entrada de senha para abrir cofre | 1 campo (senha) | [spec](tui-spec-dialog-senha.md#passwordentry) |
| PasswordCreate | Criação ou alteração de senha mestra | 2–3 campos (atual, nova, confirmação) | [spec](tui-spec-dialog-senha.md#passwordcreate) |
| FilePicker | Seleção de arquivo (abrir ou salvar) | Árvore de diretórios + campo de nome | [spec](tui-spec-dialog-filepicker.md#filepicker) |

Cada subtipo tem anatomia interna, estados e validações específicas documentadas no respectivo documento.

## Diálogos de Decisão

Todos os diálogos de decisão seguem a [Anatomia Comum](#anatomia-comum), a [Severidade](#severidade) e as convenções de [Teclado de diálogos](#teclado).

## Catálogo de Diálogos de Decisão

Esta seção lista todas as instâncias de diálogos de decisão da aplicação, especificando seu contexto, título, mensagem no corpo e ações na borda.

| Ação | Situação | Tipo | Título | Mensagem no Corpo | Ações na Borda |
|---|---|---|---|---|---|
| **Sair** | Sem alterações | Confirmação × Neutro | `Sair do Abditum` | `Sair do Abditum?` | `Enter Sair`, `Esc Voltar` |
| **Sair** | Com alterações | Confirmação × Alerta | `Sair do Abditum` | `Cofre modificado. Salvar ou descartar?` | `S Salvar`, `D Descartar`, `Esc Voltar` |
| **Salvar** | Conflito externo | Confirmação × Destrutivo | `Salvar cofre` | `Arquivo modificado externamente. Sobrescrever?` | `S Sobrescrever`, `Esc Voltar` |
| **Abrir cofre** | Falha (arquivo inválido) | Reconhecimento × Erro | `Abrir cofre` | `Arquivo corrompido ou inválido. Necessário fechar.` | `Enter OK` |
| **Abrir cofre** | Modificações não salvas | Confirmação × Alerta | `Abrir cofre` | `Cofre modificado. Salvar ou descartar?` | `S Salvar`, `D Descartar`, `Esc Voltar` |
| **Abrir cofre** | Caminho/Formato inválido | Reconhecimento × Erro | `Abrir cofre` | `Arquivo inválido ou versão não suportada. Necessário corrigir.` | `Enter OK` |
| **Abrir cofre** | Senha incorreta | Reconhecimento × Erro | `Abrir cofre` | `Senha incorreta. Necessário tentar novamente.` | `Enter OK` |
| **Criar novo cofre** | Modificações não salvas | Confirmação × Alerta | `Criar novo cofre` | `Cofre modificado. Salvar ou descartar?` | `S Salvar`, `D Descartar`, `Esc Voltar` |
| **Criar novo cofre** | Arquivo de destino existente | Confirmação × Alerta | `Criar novo cofre` | `Arquivo '[Nome]' já existe. Sobrescrever?` | `S Sobrescrever`, `I Outro caminho`, `Esc Voltar` |
| **Criar novo cofre** | Senhas não coincidem | Reconhecimento × Erro | `Criar novo cofre` | `Senhas não conferem. Necessário digitar novamente.` | `Enter OK` |
| **Criar novo cofre** | Senha fraca | Confirmação × Alerta | `Criar novo cofre` | `Senha é fraca. Prosseguir ou revisar?` | `P Prosseguir`, `R Revisar`, `Esc Voltar` |
| **Salvar cofre** | Conflito externo | Confirmação × Destrutivo | `Salvar cofre` | `Arquivo modificado externamente. Sobrescrever ou salvar como novo?` | `S Sobrescrever`, `N Salvar como novo`, `Esc Voltar` |
| **Salvar cofre como** | Destino é arquivo atual | Reconhecimento × Alerta | `Salvar cofre como` | `Destino não pode ser o arquivo atual. Necessário escolher outro.` | `Enter OK` |
| **Salvar cofre como** | Arquivo de destino existente | Confirmação × Alerta | `Salvar cofre como` | `Arquivo '[Nome]' já existe. Sobrescrever?` | `S Sobrescrever`, `I Outro caminho`, `Esc Voltar` |
| **Descartar e recarregar** | Arquivo modificado externamente | Confirmação × Destrutivo | `Descartar e recarregar` | `Cofre modificado externamente. Prosseguir com recarregamento?` | `P Prosseguir`, `Esc Voltar` |
| **Descartar e recarregar** | Confirmação de descarte | Confirmação × Destrutivo | `⚠ Descartar e recarregar` | `Todas as alterações serão descartadas. Continuar?` | `C Continuar`, `Esc Voltar` |
| **Alterar senha mestra** | Senhas não coincidem | Reconhecimento × Erro | `Alterar senha mestra` | `Senhas não conferem. Necessário digitar novamente.` | `Enter OK` |
| **Alterar senha mestra** | Senha fraca | Confirmação × Alerta | `Alterar senha mestra` | `Senha é fraca. Prosseguir ou revisar?` | `P Prosseguir`, `R Revisar`, `Esc Voltar` |
| **Alterar senha mestra** | Conflito externo | Confirmação × Destrutivo | `Alterar senha mestra` | `Arquivo modificado externamente. Sobrescrever?` | `S Sobrescrever`, `Esc Voltar` |
| **Exportar cofre** | Senha incorreta (reautenticação) | Reconhecimento × Erro | `Exportar cofre` | `Senha incorreta. Necessário tentar novamente.` | `Enter OK` |
| **Exportar cofre** | Riscos de segurança (não criptografado) | Confirmação × Alerta | `Exportar cofre` | `Arquivo não criptografado. Expor dados sensíveis?` | `E Exportar`, `Esc Voltar` |
| **Exportar cofre** | Arquivo de destino existente | Confirmação × Alerta | `Exportar cofre` | `Arquivo '[Nome]' já existe. Sobrescrever?` | `S Sobrescrever`, `I Outro caminho`, `Esc Voltar` |
| **Importar cofre** | Arquivo de intercâmbio inválido | Reconhecimento × Erro | `Importar cofre` | `Arquivo inválido ou sem Pasta Geral. Necessário corrigir.` | `Enter OK` |
| **Importar cofre** | Confirmação da política de mesclagem | Confirmação × Informativo | `Importar cofre` | `Pastas mescladas. Conflitos substituídos. Confirmar?` | `C Confirmar`, `Esc Voltar` |

---

