# Especificação Visual — Telas

> Telas completas e fluxos visuais.
> Parte da [Especificação Visual](tui-specification.md).
>
> **Documento de fundação:**
> - [`tui-design-system.md`](tui-design-system.md) — fundações visuais

## Boas-vindas

**Trigger:** Aplicação inicia sem cofre aberto, ou após fechar/bloquear cofre.  
**Interação:** Nenhuma — tela estática. Toda ação disponível via barra de comandos.

**Wireframe (área de trabalho — terminal 80 × 24):**

```
                                                                                
                                                                                
                                                                                
                   ___    __        ___ __                                      
                  /   |  / /_  ____/ (_) /___  ______ ___                       
                 / /| | / __ \/ __  / / __/ / / / __ `__ \                     
                / ___ |/ /_/ / /_/ / / /_/ /_/ / / / / / /                     
               /_/  |_/_.___/\__,_/_/\__/\__,_/_/ /_/ /_/                      
                                                                                
                             v0.1.0                                             
                                                                                
                                                                                
```

> Logo e versão centralizados via `lipgloss.Place()`. As linhas do logo recebem as cores do [DS — Gradiente do logo](tui-design-system.md#gradiente-do-logo) — não representável neste wireframe monocromático.

### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Logo (linhas 1–5) | DS — [Gradiente do logo](tui-design-system.md#gradiente-do-logo) — por linha | — |
| Versão (ex: `v0.1.0`) | `text.secondary` | — |

> As cores do logo não são tokens nomeados — são os valores hexadecimais da tabela de gradiente do DS, aplicados por linha conforme o tema ativo.

### Estados dos componentes

| Componente | Estado | Condição |
|---|---|---|
| Logo + versão | visível, centralizado | Tela ativa |
| Cabeçalho | sem abas | Nenhum cofre aberto — ver [Cabeçalho — Sem cofre](tui-spec-cabecalho.md#sem-cofre-boas-vindas) |

### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Tela entra em exibição | Dica de uso | `• Abra ou crie um cofre para começar` |

### Eventos

| Evento | Efeito |
|---|---|
| Aplicação inicia sem cofre | Modo boas-vindas exibido |
| Cofre fechado | Tela boas-vindas exibida |
| Cofre bloqueado | Tela boas-vindas exibida; arquivo permanece em disco, requer nova autenticação |
| Terminal redimensionado | Logo e versão recentralizados |

### Comportamento

- Logo e versão centralizados horizontal e verticalmente na área de trabalho via `lipgloss.Place()`
- As cores do logo acompanham o tema ativo — mudam instantaneamente com `F12`
- O cabeçalho não exibe abas neste modo (ver [Cabeçalho — Sem cofre](tui-spec-cabecalho.md#sem-cofre-boas-vindas))
- **Versão dinâmica** — o texto exibido vem da string injetada em tempo de build via `-ldflags "-X main.version=$(git describe --tags --always)"`. Em builds locais sem tag, exibe `dev`. O valor **nunca** é hardcoded no fonte

---

<!-- SEÇÕES FUTURAS — a preencher pela equipe -->

<!--
## Telas (continuação)

### Modo Cofre
### Modo Modelos
### Modo Configurações

## Componentes

### Painel Direito: Detalhe do Modelo

## Fluxos Visuais

### Criar cofre
### Abrir cofre
### Salvar cofre
### Bloquear cofre
### Alterar senha mestra
### Criar segredo
### Editar segredo
### Excluir segredo
### Buscar segredo
### Exportar cofre
### Importar cofre
-->
