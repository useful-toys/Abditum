# Especificação de Interface (UX/UI) — Abditum

| Item              | Detalhe                                   |
|-------------------|-------------------------------------------|
| **Projeto**       | Abditum — Cofre de Senhas Portátil        |
| **Versão**        | 1.0                                       |
| **Data**          | 2026-03-25                                |

---

## 1. Princípios de Design

| #  | Princípio                    | Descrição                                                                                    |
|----|------------------------------|----------------------------------------------------------------------------------------------|
| P1 | Interface TUI moderna        | Ocupa todo o terminal, interativa, suporte a 256 cores, alto contraste, estética Cyberpunk   |
| P2 | Acessibilidade de controles  | Navegação integral por teclado; suporte complementar a mouse (cliques em campos, nós, botões)|
| P3 | Ajuda contextual             | Barra de ajuda sempre acessível, exibindo ações e atalhos do contexto atual                  |
| P4 | Status global persistente    | Barra de status com caminho do arquivo, estado do cofre e total de segredos                  |
| P5 | Feedback proporcional        | Toast para ações rápidas, confirmação bloqueante para ações destrutivas/críticas             |
| P6 | Privacidade por padrão       | Campos sensíveis ocultos; revelação temporária explícita com reocultação automática          |

---

## 2. Modelo de Interação

A interface segue um modelo hierárquico de três níveis:

### 2.1. Hierarquia de Interação

```
Nível 1: PAINEL ATIVO
  │  Apenas um painel recebe entrada por vez
  │  Destaque visual (borda, cor) diferencia do inativo
  │
  └── Nível 2: ELEMENTO FOCADO
        │  Um único elemento realçado no painel ativo
        │  Sofre a ação imediata do usuário
        │
        └── Nível 3: CONTEXTO E AÇÕES
              Determinado por: painel ativo + elemento focado + estado da operação
              Ações disponíveis mudam dinamicamente
```

### 2.2. Classificação de Ações

| Tipo         | Escopo                                                                   | Exemplo                                |
|--------------|--------------------------------------------------------------------------|----------------------------------------|
| **Global**   | Disponível em todos os contextos, depende do estado geral da aplicação   | Salvar cofre, Sair                     |
| **Local**    | Disponível apenas no painel ativo                                        | Navegar árvore, Scroll vertical        |
| **Foco**     | Disponível apenas no elemento focado do painel ativo                     | Excluir segredo, Renomear pasta        |
| **Navegação**| Sempre aplicável, comportamento varia por painel e foco                  | Setas, Tab                             |

### 2.3. Regras de Interação

- Um painel inativo pode reagir ao contexto do painel ativo (ex: navegar na árvore atualiza o Painel do Segredo).
- Uma mesma ação pode ocorrer em painéis diferentes se o contexto for adequado (ex: Favoritar funciona na árvore e no Painel do Segredo).
- `Ctrl+C` é tratado como comando convencional, não interrompe a aplicação.

---

## 3. Layout Principal

### 3.1. Estrutura de Painéis

```
┌─────────────────────────────────────────────────────────────────────┐
│                          BARRA DE STATUS                            │
│  📁 /caminho/cofre.abditum   │  ● Cofre Modificado   │  42 segredos │
├────────────────────────┬────────────────────────────────────────────┤
│                        │                                            │
│  PAINEL DA HIERARQUIA  │        PAINEL DO SEGREDO                  │
│                        │                                            │
│  ★ Favoritos (3)       │  ┌─────────────────────────────────────┐  │
│  ▼ Sites (5)           │  │  Nome: GitHub Login          [★]   │  │
│    ▼ Redes Sociais (2) │  │  Modelo: Login                      │  │
│      ● Twitter         │  │                                     │  │
│      ● Facebook        │  │  URL:      github.com               │  │
│    ► Email (3)         │  │  Username: usuario                  │  │
│  ▼ Financeiro (4)      │  │  Password: ••••••••    [👁] [📋]   │  │
│    ● Banco X           │  │                                     │  │
│    ● Banco Y           │  │  Observação:                        │  │
│  ► Serviços (2)        │  │  Conta principal de desenvolvimento │  │
│  🗑 Lixeira (1)        │  │                                     │  │
│                        │  └─────────────────────────────────────┘  │
│                        │                                            │
├────────────────────────┴────────────────────────────────────────────┤
│                       BARRA DE AJUDA                                │
│  Ctrl+N Novo  │  Ctrl+E Editar  │  Del Excluir  │  Ctrl+Q Sair    │
└─────────────────────────────────────────────────────────────────────┘
```

### 3.2. Componentes Estruturais

| Componente               | Posição         | Conteúdo                                                                     |
|--------------------------|-----------------|------------------------------------------------------------------------------|
| **Barra de Status**      | Topo            | Caminho do arquivo, estado (Salvo/Modificado), total de segredos             |
| **Painel da Hierarquia** | Esquerda        | Árvore de pastas e segredos, pastas virtuais (Favoritos, Lixeira)            |
| **Painel do Segredo**    | Direita         | Detalhe, criação ou edição do segredo selecionado                            |
| **Barra de Ajuda**       | Rodapé          | Ações e atalhos disponíveis no contexto atual                                |

### 3.3. Responsividade e Tamanho Mínimo

- O tamanho mínimo do terminal é definido pelas dimensões necessárias para exibir a tela inicial (ASCII art + molduras).
- Se o terminal estiver abaixo do mínimo, ocultar painéis e exibir mensagem pedindo redimensionamento.

---

## 4. Telas e Estados Visuais

### 4.1. Tela Inicial (Welcome)

```
┌─────────────────────────────────────────────────────┐
│                                                     │
│              █████  ██████  ██████  ██  ████████    │
│             ██   ██ ██   ██ ██   ██ ██     ██       │
│             ███████ ██████  ██   ██ ██     ██       │
│             ██   ██ ██   ██ ██   ██ ██     ██       │
│             ██   ██ ██████  ██████  ██     ██       │
│                     Cofre de Senhas Portátil         │
│                                                     │
│             [A] Abrir cofre existente               │
│             [C] Criar novo cofre                    │
│             [?] Ajuda                               │
│             [Q] Sair                                │
│                                                     │
└─────────────────────────────────────────────────────┘
```

- Apresentada ao iniciar a aplicação ou após bloqueio do cofre.
- Estado: `Inicial / sem cofre ativo`.

### 4.2. File Picker (Seleção de Arquivo)

- Componente TUI integrado para navegação de diretórios.
- Suporta navegação por setas, mouse e autocompletar.
- Usado nos fluxos: Criar cofre, Abrir cofre, Salvar como, Exportar, Importar.

### 4.3. Formulário de Senha Mestra

- Campos de entrada mascarados (caracteres substituídos por `•`).
- Na criação/alteração: digitação dupla para confirmação.
- Aviso Zero Knowledge na criação: "O esquecimento da senha mestra resulta em perda total dos dados."

### 4.4. Cofre Ativo — Layout Dual Panel

Conforme seção 3.1. A alternância entre painéis é feita com `Tab`.

### 4.5. Modais e Confirmações

```
┌──────────────────────────────────────┐
│  ❓ Excluir segredo "GitHub Login"?  │
│                                      │
│  Ele poderá ser restaurado até o     │
│  próximo salvamento.                 │
│                                      │
│    [Excluir]       [Voltar]          │
└──────────────────────────────────────┘
```

- Confirmações usam verbos de ação (nunca "Sim/Não/Cancelar").
- Ações destrutivas irreversíveis: excluir pasta, excluir modelo, alterar senha mestra.

---

## 5. Painel da Hierarquia — Especificação

### 5.1. Estrutura da Árvore

| Elemento                  | Indicador Visual                                                        |
|---------------------------|-------------------------------------------------------------------------|
| Pasta expandida           | `▼ Nome (N)` — N = total de segredos incluindo subpastas               |
| Pasta colapsada           | `► Nome (N)` — com indicador de conteúdo                               |
| Segredo normal            | `● Nome`                                                                |
| Segredo favorito          | `★ Nome` — cor ou ícone de destaque                                    |
| Segredo novo              | `● Nome` — indicador visual (cor/badge) de "novo"                      |
| Segredo modificado        | `● Nome` — indicador visual (cor/badge) de "modificado"                |
| Segredo em edição         | `● Nome` — indicador visual de "em edição"                             |
| Pasta virtual Favoritos   | `★ Favoritos (N)` — topo da raiz, visível só se houver favoritos      |
| Pasta virtual Lixeira     | `🗑 Lixeira (N)` — final da raiz, visível só se houver itens excluídos |

### 5.2. Comportamento de Navegação

| Ação                          | Tecla                  | Comportamento                                                |
|-------------------------------|------------------------|--------------------------------------------------------------|
| Mover foco para cima/baixo    | `↑` / `↓`             | Navega entre nós visíveis da árvore                          |
| Expandir pasta                | `→` (pasta colapsada)  | Revela filhos da pasta                                       |
| Entrar na pasta               | `→` (pasta expandida)  | Move foco para o primeiro filho                              |
| Colapsar pasta                | `←` (pasta expandida)  | Esconde filhos da pasta                                      |
| Ir para pasta pai             | `←` (colapsada/segredo)| Move foco para a pasta pai                                   |
| Toggle pasta                  | `Enter` (pasta)        | Expande/colapsa                                              |
| Visualizar segredo            | `Enter` (segredo)      | Exibe detalhes no Painel do Segredo, foco vai para lá        |
| Busca alfabética              | Letras `a-z`           | Pula para próximo item correspondente                        |

### 5.3. Criação Relativa ao Foco

| Foco atual          | Destino do novo segredo/pasta                                          |
|---------------------|------------------------------------------------------------------------|
| Pasta               | Final da coleção interna dessa pasta                                   |
| Segredo             | Logo abaixo do segredo, na mesma coleção pai                           |

### 5.4. Manutenção de Foco

- Ao remover/adicionar itens, o foco permanece no mesmo nó se possível.
- Se o nó foi removido: foco vai para o próximo item.
- Se não houver próximo: recua para o nó pai.
- Scroll acompanha automaticamente o foco em cofres grandes.

### 5.5. Busca

- `Ctrl+F` ativa a barra de filtragem.
- Árvore é filtrada mostrando apenas segredos que satisfazem o critério, preservando contexto de pastas.
- Casamento no nome do segredo: trecho pesquisado recebe highlight de cor.
- Durante a busca, ações indisponíveis exceto: sair, navegar, visualizar.
- Confirmar busca (selecionar elemento) ou cancelar busca encerra o modo de pesquisa.

---

## 6. Painel do Segredo — Especificação

### 6.1. Modo Visualização

```
┌──────────────────────────────────────────────────┐
│  Nome: GitHub Login                         [★]  │
│  Modelo: Login                                   │
│  Criado: 2026-01-15 10:30  Modificado: 2026-03-20│
│──────────────────────────────────────────────────│
│                                                  │
│  URL:       github.com                    [📋]   │
│  Username:  usuario                       [📋]   │
│  Password:  ••••••••          [👁] [📋]          │
│                                                  │
│──────────────────────────────────────────────────│
│  Observação:                                     │
│  Conta principal de desenvolvimento.             │
│  MFA ativado via app autenticador.               │
│                                                  │
└──────────────────────────────────────────────────┘
```

- Campos `texto_sensivel` exibidos como `••••••••` por padrão.
- `[👁]` toggle para revelar/ocultar temporariamente.
- `[📋]` copia o valor do campo para o clipboard.
- Navegação com `↑`/`↓` move foco entre campos.
- `Esc` devolve foco para o Painel da Hierarquia.
- Se não houver segredo selecionado: exibir placeholder informativo.

### 6.2. Modo Edição Padrão

Permite alterar:
- Nome do segredo
- Observação
- Valores dos campos existentes

Não permite: alterar estrutura (adicionar/remover/reordenar campos).

O usuário pode alternar para edição avançada a qualquer momento.

### 6.3. Modo Edição Avançada

Permite alterar exclusivamente a estrutura:
- Adicionar novo campo (nome + tipo)
- Renomear campo
- Excluir campo
- Reordenar campos

**Não permite** alterar o tipo de um campo existente — deve excluir e recriar.

O usuário pode alternar para edição padrão a qualquer momento.

### 6.4. Modos de Entrada por Fluxo

| Fluxo de Criação             | Modo Inicial           |
|------------------------------|------------------------|
| A partir de modelo           | Edição padrão          |
| Segredo vazio                | Edição avançada        |

### 6.5. Gerenciamento de Espaço

- Se houver muitos campos: scroll vertical no Painel do Segredo.
- Campo "Observação" é redimensionável e ocupa automaticamente o espaço livre restante.

---

## 7. Pastas Virtuais

### 7.1. Favoritos

| Aspecto         | Comportamento                                                                |
|-----------------|------------------------------------------------------------------------------|
| Posição         | Topo da raiz da hierarquia                                                   |
| Visibilidade    | Apenas quando existir ≥1 segredo favoritado                                  |
| Conteúdo        | Atalhos para segredos favoritados (sem alterar localização real)             |
| Interação       | Visualizar, editar, desfavoritar — mesmas ações do segredo na hierarquia     |

### 7.2. Lixeira

| Aspecto         | Comportamento                                                                |
|-----------------|------------------------------------------------------------------------------|
| Posição         | Final da raiz da hierarquia                                                  |
| Visibilidade    | Apenas quando existir ≥1 segredo excluído reversivelmente                   |
| Conteúdo        | Segredos excluídos reversivelmente (soft delete)                             |
| Interação       | Restaurar segredo; sem edição enquanto na Lixeira                            |
| Ciclo de vida   | Esvaziada irreversivelmente ao salvar o cofre                                |

---

## 8. Sistema de Feedback

### 8.1. Tipos de Mensagem

| Tipo              | Cor       | Ícone | Comportamento               | Exemplo                                                  |
|-------------------|-----------|-------|-----------------------------|----------------------------------------------------------|
| **Sucesso**       | Verde     | ✅    | Non-blocking, auto-oculta    | "Segredo copiado para área de transferência."             |
| **Erro**          | Vermelha  | ❌    | Non-blocking, tempo estendido| "Não foi possível salvar: caminho de arquivo inválido."   |
| **Aviso/Alerta**  | Amarela   | ⚠️    | Non-blocking, tempo estendido| "O cofre será bloqueado em 30 segundos por inatividade."  |
| **Confirmação**   | Laranja   | ❓    | Blocking                     | "Excluir segredo? Ele poderá ser restaurado até o próximo salvamento." |
| **Informativa**   | Azul/Cinza| ℹ️    | Non-blocking, auto-oculta    | "Informe os dados do novo segredo."                       |

### 8.2. Diretrizes de Redação

| Regra                                                                               |
|---------------------------------------------------------------------------------------|
| Ser direto e específico — sem frases genéricas                                        |
| Sem exclamação nem palavras-rótulo ("Sucesso!", "Erro!", "Atenção!")                  |
| Não mencionar teclas — mecanismo separado                                             |
| Não usar "com sucesso", "realizado" ou "Tem certeza que deseja..."                   |
| Confirmações indicam impacto; opções são verbos de ação (nunca "Sim/Não/Cancelar")   |
| Erros descrevem o que falhou e sugerem correção quando possível                       |
| Erros de escrita após backup: informar existência do backup para intervenção manual   |

### 8.3. Regras de Feedback por Ação

| Categoria                           | Feedback                                                          |
|-------------------------------------|-------------------------------------------------------------------|
| Ações demoradas                     | Indicador de progresso + resultado ao final (auto-oculta)         |
| Ações visualmente evidentes         | Sem feedback adicional (estado do visual é suficiente)            |
| Operações interativas               | Mensagem informativa/instrucional ao início                       |
| Ações destrutivas/críticas          | Confirmação bloqueante obrigatória                                 |

---

## 9. Indicadores Visuais de Estado

### 9.1. Estado Global (Barra de Status)

| Indicador                | Quando exibido                                       |
|--------------------------|------------------------------------------------------|
| Caminho do arquivo       | Sempre que houver cofre ativo                        |
| `Cofre Salvo`            | Cofre sincronizado com arquivo                       |
| `Cofre Modificado`       | Divergência entre memória e último salvamento        |
| Total de segredos        | Sempre que houver cofre ativo                        |

### 9.2. Estado de Segredos (Árvore e Painel)

| Estado                   | Indicador Visual                                     |
|--------------------------|------------------------------------------------------|
| Favorito                 | Ícone/cor de destaque (★)                            |
| Novo (não salvo)         | Badge/cor indicando "novo"                           |
| Modificado (não salvo)   | Badge/cor indicando "modificado"                     |
| Em edição                | Badge/cor indicando "em edição"                      |
| Excluído (Lixeira)       | Visível apenas na pasta virtual Lixeira              |

### 9.3. Campos Sensíveis

| Estado                    | Exibição                                             |
|---------------------------|------------------------------------------------------|
| Oculto (padrão)           | `••••••••`                                           |
| Revelado temporariamente  | Valor real; reocultação automática após N segundos   |

### 9.4. Área de Transferência

| Estado                    | Indicador                                            |
|---------------------------|------------------------------------------------------|
| Campo copiado             | Toast de sucesso + countdown visual até limpeza      |
| Limpa (automática/manual) | Sem indicador (estado padrão)                        |

---

## 10. Clipboard — Especificação

| Aspecto                         | Comportamento                                                          |
|---------------------------------|------------------------------------------------------------------------|
| Copiar campo                    | Qualquer campo (incluindo `texto_sensivel`), sem necessidade de revelar |
| Feedback de cópia               | Toast de sucesso                                                        |
| Countdown visual                | Indicador na interface com tempo restante até limpeza                   |
| Timer de limpeza                | Configurável (padrão: 30 s)                                            |
| Limpeza ao bloquear/fechar      | Automática e imediata                                                  |

---

## 11. Bloqueio por Inatividade

| Aspecto                         | Especificação                                                          |
|---------------------------------|------------------------------------------------------------------------|
| Timer configurável              | Padrão sugerido: 2 minutos                                            |
| Alerta de bloqueio iminente     | Exibido quando transcorrerem 75% do tempo configurado                  |
| Atividade válida                | Teclado e clique de mouse (movimento sem clique não conta)             |
| Reset do timer                  | Ao término de cada ação do usuário (rápida ou demorada)                |
| Abortar bloqueio iminente       | Qualquer atividade válida cancela e reinicia o timer                   |
| Ao bloquear                     | Descartar domínio, limpar buffers, limpar clipboard, exibir tela de abertura |
| Alterações não salvas           | Descartadas silenciosamente (sem confirmação — decisão de projeto)     |

---

## 12. Proteção contra Shoulder Surfing

| Mecanismo                    | Descrição                                                              |
|------------------------------|------------------------------------------------------------------------|
| Campos sensíveis ocultos     | Padrão `••••••••`, revelação temporária por ação explícita             |
| Atalho de ocultação total    | Ocultar toda a interface rapidamente                                   |
| Reocultação automática       | Timer configurável (padrão: 15 s) para campos revelados                |

---

## 13. Mapa de Teclas de Comando

### 13.1. Global / Qualquer Contexto

| Tecla        | Ação                                                                      |
|--------------|---------------------------------------------------------------------------|
| `Ctrl+Q`     | Iniciar fluxo de saída                                                    |
| `Ctrl+S`     | Salvar cofre (disponível se `Cofre Modificado`)                           |
| `Tab`        | Alternar painel ativo (Hierarquia ↔ Segredo)                              |

### 13.2. Confirmação de Saída

| Estado                    | Tecla         | Ação                                   |
|---------------------------|---------------|----------------------------------------|
| Sem modificações          | `Ctrl+Q`      | Confirmar saída                        |
| Com modificações          | `Ctrl+S`      | Salvar e sair                          |
| Com modificações          | `Ctrl+D`      | Descartar alterações e sair            |
| Qualquer                  | `Esc`         | Cancelar (voltar)                      |

### 13.3. Foco no Painel da Hierarquia

| Tecla         | Contexto              | Ação                                            |
|---------------|-----------------------|-------------------------------------------------|
| `↑` / `↓`    | Qualquer              | Mover foco entre nós visíveis                   |
| `→`           | Pasta colapsada       | Expandir pasta                                  |
| `→`           | Pasta expandida       | Mover foco para primeiro filho                  |
| `→`           | Segredo               | Sem efeito                                      |
| `←`           | Pasta expandida       | Colapsar pasta                                  |
| `←`           | Pasta colapsada/segredo| Mover foco para pasta pai                      |
| `Enter`       | Pasta                 | Toggle expandir/colapsar                        |
| `Enter`       | Segredo               | Visualizar no Painel do Segredo + mover foco    |
| `a-z`         | Qualquer              | Pular para próximo item alfabético              |
| `Ctrl+N`      | Qualquer              | Novo segredo (destino relativo ao foco)         |
| `Ctrl+F`      | Qualquer              | Ativar barra de busca                           |

### 13.4. Foco no Painel do Segredo (Visualização)

| Tecla         | Ação                                                      |
|---------------|-----------------------------------------------------------|
| `↑` / `↓`    | Mover foco entre campos                                   |
| `Esc`         | Devolver foco para o Painel da Hierarquia                 |

---

## 14. Fluxos de Interface

### 14.1. Fluxo de Criação de Cofre

```
Welcome → [C] Criar cofre
  │
  ▼
File Picker (selecionar caminho)
  │
  ├─ Arquivo já existe → Confirmação de sobrescrita
  │
  ▼
Formulário de Senha Mestra (dupla digitação)
  │
  ▼
Aviso Zero Knowledge
  │
  ▼
Cofre criado → Layout Dual Panel (Cofre Salvo)
```

### 14.2. Fluxo de Abertura de Cofre

```
Welcome → [A] Abrir cofre
  │
  ▼
File Picker (selecionar arquivo .abditum)
  │
  ▼
Validar magic + versão_formato
  │
  ├─ Inválido → Erro: "Arquivo não é um cofre Abditum válido"
  │
  ▼
Formulário de Senha Mestra
  │
  ▼
Derivar chave (Argon2id) → Indicador de progresso
  │
  ├─ Falha → Erro: "Senha mestra incorreta ou arquivo corrompido"
  │
  ▼
Cofre aberto → Layout Dual Panel (Cofre Salvo)
```

### 14.3. Fluxo de Criação de Segredo

```
Ctrl+N (no Painel da Hierarquia)
  │
  ▼
Escolher: [M] Usar modelo  /  [V] Segredo vazio
  │                              │
  ▼                              ▼
Selecionar modelo            Edição avançada
  │                          (adicionar campos)
  ▼
Edição padrão
(preencher valores)
  │
  ▼
Confirmar → Segredo novo na hierarquia → Cofre Modificado
```

### 14.4. Fluxo de Salvamento

```
Ctrl+S (Cofre Modificado)
  │
  ▼
Indicador de progresso
  │
  ├─ Sucesso → Toast "Cofre salvo" → Cofre Salvo
  │             Lixeira esvaziada permanentemente
  │
  └─ Falha → Erro com informação de backup disponível
```

### 14.5. Fluxo de Exportação

```
Ação: Exportar
  │
  ▼
┌─ Cofre Modificado? → Alerta: exportação incluirá alterações não salvas
  │
  ▼
Aviso de segurança (JSON não criptografado)
  │
  ▼
Confirmação: [Exportar] / [Voltar]
  │
  ▼
File Picker (caminho de destino)
  │
  ▼
Toast de sucesso
```

### 14.6. Fluxo de Bloqueio

```
Timer de inatividade a 75% → Alerta: "Bloqueio em N segundos"
  │
  ├─ Atividade válida → Cancelar alerta, reiniciar timer
  │
  └─ Sem atividade → Bloqueio automático
                         │
Manual: atalho de bloqueio ──┤
                              │
                              ▼
                    Limpar domínio + buffers + clipboard
                              │
                              ▼
                    Tela de abertura (mesmo caminho do cofre)
                    Exigir nova autenticação
```

---

## 15. Aviso de Segurança — Textos Obrigatórios

| Contexto                          | Mensagem                                                                                         |
|-----------------------------------|--------------------------------------------------------------------------------------------------|
| Criação de cofre                  | "A senha mestra é a única forma de acessar seus dados. Em caso de esquecimento, não há recuperação possível." |
| Exportação para JSON              | "A exportação gera uma cópia não criptografada de todos os dados do cofre, incluindo senhas e informações sensíveis." |
| Falha de salvamento com backup    | "Falha ao salvar o cofre. Um backup do arquivo anterior está disponível com extensão .bak para intervenção manual." |

---

## 16. Importação — Feedback de Conflitos

| Tipo de Conflito                  | Feedback ao Usuário                                                                              |
|-----------------------------------|-------------------------------------------------------------------------------------------------|
| Segredos com nome conflitante     | Mensagem informativa: itens importados receberam sufixos numéricos incrementais                  |
| Pastas com identidade conflitante | Silencioso — merge automático                                                                    |
| Modelos com identidade conflitante| Silencioso — substituição pelo modelo importado                                                  |

---

## 17. Rastreabilidade

| Seção da Especificação              | Requisitos Associados                                         |
|--------------------------------------|--------------------------------------------------------------|
| Modelo de Interação                  | [RF-12](srs.md#rf-12), [RF-13](srs.md#rf-13), [RF-14](srs.md#rf-14)                                          |
| Painel da Hierarquia                 | [RF-12](srs.md#rf-12), [RF-13](srs.md#rf-13), [RF-26](srs.md#rf-26), [RF-27](srs.md#rf-27) a [RF-31](srs.md#rf-31)                           |
| Painel do Segredo                    | [RF-14](srs.md#rf-14), [RF-15](srs.md#rf-15) a [RF-25](srs.md#rf-25)                                         |
| Pastas Virtuais                      | [RF-19](srs.md#rf-19), [RF-22](srs.md#rf-22), [RF-23](srs.md#rf-23)                                          |
| Sistema de Feedback                  | [RNF-25](srs.md#rnf-25), [RNF-26](srs.md#rnf-26)                                               |
| Clipboard                           | [RF-36](srs.md#rf-36), [RF-37](srs.md#rf-37), [RF-38](srs.md#rf-38)                                          |
| Bloqueio por Inatividade             | [RF-07](srs.md#rf-07), [RF-08](srs.md#rf-08)                                                 |
| Proteção Shoulder Surfing            | [RF-40](srs.md#rf-40)                                                        |
| File Picker                         | [RF-01](srs.md#rf-01), [RF-02](srs.md#rf-02), [RF-05](srs.md#rf-05), [RF-06](srs.md#rf-06)                                   |
| Exportação/Importação                | [RF-09](srs.md#rf-09), [RF-10](srs.md#rf-10)                                                  |
