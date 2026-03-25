# 04 — UX e Interface

## 04.1 Princípios de UX

### Princípios fundamentais

| Princípio | Aplicação no Abditum |
|---|---|
| **Portabilidade da experiência** | A TUI funciona em qualquer terminal — local, remoto (SSH), container. Nenhuma dependência de ambiente gráfico. |
| **Segurança sem atrito** | Campos sensíveis ocultos por padrão, limpeza automática do clipboard, bloqueio por inatividade — tudo operando em segundo plano sem exigir ação do usuário. |
| **Contexto sempre visível** | Barra de status (caminho, estado, total de segredos) e barra de ajuda contextual (ações e atalhos válidos) sempre presentes. |
| **Feedback proporcional** | Ações com resultado visível dispensam notificação. Ações destrutivas exigem confirmação bloqueante. Erros persistem mais tempo na tela. |
| **Prevenção de erros** | Soft delete com Lixeira, confirmação para ações irreversíveis, aviso de irrecuperabilidade ao criar cofre, salvamento atômico com backup. |

### Diretrizes de interação

- **Teclado primeiro, mouse complementar:** Navegação completa pelo teclado. Mouse suporta cliques em nós da árvore, campos de edição e botões.
- **Um painel ativo por vez:** Apenas um painel recebe input. Destaque visual claro (borda/cor) indica qual é o ativo.
- **Foco único por painel:** Dentro do painel ativo, um elemento está realçado — é ele que recebe a ação.
- **Ações contextuais dinâmicas:** A barra de ajuda mostra apenas ações válidas para o contexto atual (painel + foco + estado).

### Diretrizes de redação de mensagens

| Regra | Exemplo |
|---|---|
| Ser direto e específico | "Segredo movido para a Lixeira." |
| Sem exclamação nem rótulos ("Sucesso!", "Erro!") | A cor + ícone já diferenciam o tipo |
| Sem menção a teclas nas mensagens | Teclas aparecem na barra de ajuda dedicada |
| Sem "com sucesso" ou "Tem certeza que deseja..." | Redundantes ou verbosos |
| Confirmações usam verbos de ação | "Excluir", "Salvar", "Voltar" — nunca "Sim/Não/Cancelar" |

### Padrão de mensagens

| Tipo | Cor | Ícone | Comportamento |
|---|---|---|---|
| Erro | Vermelha | ❌ | Não bloqueante, auto-oculta (tempo estendido) |
| Aviso/Alerta | Amarela | ⚠️ | Não bloqueante, auto-oculta (tempo estendido) |
| Confirmação | Laranja | ❓ | Bloqueante — opções são verbos de ação |
| Sucesso | Verde | ✅ | Não bloqueante, auto-oculta |
| Informativa | Azul/cinza | ℹ️ | Não bloqueante, auto-oculta |

---

## 04.2 Wireframes / Mockups

### Tela inicial (estado: Inicial / sem cofre ativo)

```
┌──────────────────────────────────────────────────────────────────────┐
│                                                                      │
│                                                                      │
│                     █████╗ ██████╗ ██████╗ ████████╗                 │
│                    ██╔══██╗██╔══██╗██╔══██╗╚══██╔══╝                 │
│                    ███████║██████╔╝██║  ██║   ██║                    │
│                    ██╔══██║██╔══██╗██║  ██║   ██║                    │
│                    ██║  ██║██████╔╝██████╔╝   ██║                    │
│                    ╚═╝  ╚═╝╚═════╝ ╚═════╝    ╚═╝                    │
│                                                                      │
│                    Cofre de Senhas Portátil e Seguro                 │
│                                                                      │
│                                                                      │
│                                                                      │
│                                                                      │
│                                                                      │
│                                                                      │
│                                                                      │
│                                                                      │
│                                                                      │
│                                                                      │
│                                                                      │
├──────────────────────────────────────────────────────────────────────┤
│  [C] Criar Cofre  [A] Abrir Cofre  [?] Ajuda  [Ctrl+Q] Sair        │
└──────────────────────────────────────────────────────────────────────┘
```

### Layout principal (estado: Cofre ativo)

```
┌──────────────────────────────────────────────────────────────────────┐
│  ABDITUM                     meu-cofre.abditum  │ ● Modificado │ 42 │
├────────────────────┬─────────────────────────────────────────────────┤
│  HIERARQUIA        │  SEGREDO                                       │
│                    │                                                 │
│  ★ Favoritos (3)   │  Nome: GitHub Personal                         │
│  ─────────────     │  Modelo: Login                                 │
│  ● GitHub Person.. │  ─────────────────────────────                 │
│    AWS Root     ★  │  URL:      https://github.com                  │
│  ▶ Sites (5)       │  Username: user@mail.com                       │
│  ▼ Financeiro (2)  │  Password: ••••••••••••  👁                    │
│    ● Banco X       │                                                 │
│    ● Cartão Y      │  Observação:                                   │
│  ▶ Serviços (8)    │  Conta pessoal, 2FA ativo                      │
│                    │                                                 │
│                    │                                                 │
│                    │                                                 │
│                    │                                                 │
│                    │                                                 │
│  🗑 Lixeira (1)    │                                                 │
│                    │                                                 │
├────────────────────┴─────────────────────────────────────────────────┤
│  [Tab] Painel │ [E] Editar │ [C] Copiar │ [F] Favoritar │ [?] Ajuda│
└──────────────────────────────────────────────────────────────────────┘
```

### Layout de edição padrão (Painel do Segredo)

```
│  SEGREDO  ✏️ Edição Padrão                     │
│                                                 │
│  Nome: [GitHub Personal____________]            │
│  ─────────────────────────────────              │
│  URL:      [https://github.com_____]            │
│  Username: [user@mail.com__________]            │
│  Password: [••••••••••••___________] 👁         │
│                                                 │
│  Observação:                                    │
│  [Conta pessoal, 2FA ativo_________]            │
│  [_________________________________]            │
│  [_________________________________]            │
│                                                 │
│  [Confirmar]  [Avançada]  [Cancelar]            │
```

### Layout de edição avançada (Painel do Segredo)

```
│  SEGREDO  🔧 Edição Avançada                   │
│                                                 │
│  Campos:                                        │
│  ┌───┬────────────┬──────────────┬─────┐        │
│  │ # │ Nome       │ Tipo         │     │        │
│  ├───┼────────────┼──────────────┼─────┤        │
│  │ 1 │ URL        │ texto        │ ✕ ↕ │        │
│  │ 2 │ Username   │ texto        │ ✕ ↕ │        │
│  │ 3 │ Password   │ texto sensív │ ✕ ↕ │        │
│  └───┴────────────┴──────────────┴─────┘        │
│                                                 │
│  [+ Adicionar Campo]                            │
│                                                 │
│  [Confirmar]  [Padrão]  [Cancelar]              │
```

### Confirmação bloqueante (modal)

```
┌──────────────────────────────────────────┐
│ ❓ Excluir segredo? Ele poderá ser       │
│    restaurado até o próximo salvamento.  │
│                                          │
│    [Excluir]              [Voltar]       │
└──────────────────────────────────────────┘
```

### Toast de sucesso (não bloqueante)

```
┌──────────────────────────────────────────┐
│ ✅ Segredo copiado para área de          │
│    transferência.            ⏱ 28s       │
└──────────────────────────────────────────┘
```

---

## 04.3 Mapa de Navegação

### Estados e transições da aplicação

```
┌──────────────┐
│ Tela Inicial │◄──────────────────────────────────────┐
│ (ASCII Art)  │                                       │
└──────┬───────┘                                       │
       │                                               │
  ┌────┴─────┐                                         │
  │          │                                         │
  ▼          ▼                                         │
Criar     Abrir                                        │
Cofre     Cofre                                        │
  │          │                                         │
  │  ┌───────┘                                         │
  │  │ File Picker + Senha                             │
  │  │ Validação + Migração                            │
  ▼  ▼                                                 │
┌─────────────────────────────────┐                    │
│          COFRE ATIVO            │                    │
│  ┌─────────┐    ┌────────────┐  │                    │
│  │  Salvo  │◄──►│ Modificado │  │    Bloquear        │
│  └─────────┘    └────────────┘  │────────────────────┘
│                                 │
│  Hierarquia ◄─Tab─► Segredo    │
│                                 │
│  Busca (overlay transitório)    │
└──────────┬──────────────────────┘
           │
    ┌──────┴──────┐
    │             │
    ▼             ▼
  Sair      Salvar Como
(confirm)   (File Picker)
```

### Fluxos modais (overlays)

```
Cofre Ativo
    │
    ├── Confirmação de exclusão (segredo/pasta/modelo)
    ├── Confirmação de sobrescrita (salvar como / criar)
    ├── Confirmação de saída com alterações não salvas
    ├── Aviso de exportação (risco de segurança)
    ├── Alerta de bloqueio iminente (75% do tempo)
    ├── Formulário de nova pasta (inline)
    ├── Formulário de configuração do cofre
    ├── Formulário de alteração de senha mestra
    ├── Seletor de modelo (criação de segredo)
    ├── Seletor de destino (mover segredo/pasta)
    └── Tela de ajuda
```

### Mapa de teclas por contexto

**Global / Qualquer contexto:**

| Tecla | Ação |
|---|---|
| Ctrl+Q | Sair (com tratamento de alterações não salvas) |
| Ctrl+S | Salvar cofre (quando Cofre Modificado) |
| Ctrl+D | Descartar alterações (quando Cofre Modificado) |
| Esc | Voltar / cancelar ação atual |
| Tab | Alternar entre Painel da Hierarquia e Painel do Segredo |

**Foco no Painel da Hierarquia:**

| Tecla | Ação |
|---|---|
| ↑ / ↓ | Mover foco entre linhas |
| → | Expandir pasta / mover foco para primeiro filho |
| ← | Colapsar pasta / mover foco para pasta pai |
| Enter | Pasta: expandir/colapsar. Segredo: visualizar no Painel do Segredo |
| a-z | Pular para próximo item alfabético |
| Ctrl+N | Novo segredo na localização do foco |
| Ctrl+F | Ativar barra de busca |

**Foco no Painel do Segredo (visualizando):**

| Tecla | Ação |
|---|---|
| ↑ / ↓ | Mover foco entre campos |
| Esc | Devolver foco ao Painel da Hierarquia |

**Foco no Painel do Segredo (editando):**

| Tecla | Ação |
|---|---|
| ↑ / ↓ / Tab | Navegar entre campos do formulário |
| Esc | Cancelar edição |

---

## 04.4 Guia de Estilo

### Paleta de cores (estética Cyberpunk, 256 cores)

| Elemento | Cor sugerida | Uso |
|---|---|---|
| Fundo principal | Preto / cinza muito escuro | Background de todos os painéis |
| Texto primário | Branco / cinza claro | Conteúdo principal, nomes de segredos |
| Texto secundário | Cinza médio | Labels, contadores, timestamps |
| Destaque (accent) | Ciano / azul elétrico | Painel ativo (borda), foco, highlight de busca |
| Favorito | Amarelo / dourado | Ícone ★ e nome do segredo favorito |
| Novo/Modificado | Verde | Indicador de estado na árvore |
| Em edição | Magenta / roxo | Indicador de segredo em edição |
| Sensível oculto | Cinza escuro | Representação visual de `••••••` |
| Erro | Vermelho | Toasts de erro, bordas de campos inválidos |
| Aviso | Amarelo | Toasts de alerta |
| Sucesso | Verde | Toasts de sucesso |
| Confirmação | Laranja | Modais bloqueantes |
| Pasta virtual (Favoritos) | Amarelo (dim) | Nome da pasta virtual |
| Pasta virtual (Lixeira) | Vermelho (dim) | Nome da pasta virtual |

### Tipografia (modo texto)

| Elemento | Estilo |
|---|---|
| Títulos de painel | MAIÚSCULAS, bold |
| Nome do segredo (árvore) | Normal, com ícone de estado à esquerda |
| Nome do segredo (detalhe) | Bold |
| Labels de campo | Normal, cor secundária |
| Valores de campo | Normal, cor primária |
| Campo sensível oculto | `••••••••••••` em cor dim |
| Barra de status | Compacta, separadores `│` |
| Barra de ajuda | `[Tecla] Ação` com separadores `│` |

### Ícones e indicadores

| Ícone | Significado |
|---|---|
| ★ | Segredo favorito |
| ● | Segredo (nó folha na árvore) |
| ▶ | Pasta colapsada (com conteúdo) |
| ▼ | Pasta expandida |
| ✏️ | Segredo em edição |
| 🗑 | Pasta virtual Lixeira |
| 👁 | Toggle de campo sensível |
| ⏱ | Countdown do clipboard |

### Bordas e molduras

| Elemento | Estilo |
|---|---|
| Painel ativo | Borda dupla ou cor accent |
| Painel inativo | Borda simples ou cor dim |
| Modal/Confirmação | Borda dupla centralizada, fundo contrastante |
| Toast | Borda arredondada, cor de fundo conforme tipo |
| Separador interno | Linha horizontal `─────` |
