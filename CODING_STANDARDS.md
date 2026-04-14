# Padrões de Código - Abditum

Este documento define os padrões de código que todos os agentes de IA devem seguir ao trabalhar neste projeto.

## Nomes de Arquivos

Usar nomes descritivos e informativos, mesmo que mais longos:

| Em vez de | Use |
|-----------|-----|
| `tree.go` | `vault_tree.go` |
| `detail.go` | `secret_detail.go` |
| `list.go` | `template_list.go` |
| `view.go` | `welcome_view.go` |

## Nomes de Métodos e Variáveis

Preferir nomes claros e descritivos a nomes muito abreviados:

| Em vez de | Use |
|-----------|-----|
| `fn` | `FetchPassword` |
| `msg` | `WindowSizeMessage` |
| `tm` | `ThemeManager` |
| `proc` | `ProcessUserInput` |

## Comentários

Os comentários devem explicar o **porquê** das decisões de implementação, não apenas descrever o que o código faz:

```go
// NewPasswordModal cria um novo modal de entrada de senha.
// O modal é configurado com título e callback para processar o resultado.
// O callback recebe a senha informada pelo usuário.
func NewPasswordModal(title string, onResult func([]byte)) ModalView
```

```go
// EncryptVault executa a criptografia do cofre.
// Escolhemos AES-GCM porque oferece autenticação
// implícita junto com confidencialidade.
func (m *Manager) EncryptVault() error
```

## Acessibilidade

O código deve ser acessível a leitores com menos familiaridade com Go:

- Evitar shortcuts obscuros
- Usar nomes que revelam intent
- Documentar estruturas complexas
- Explicar decisões não óbvias

## Aplicação

Estes padrões se aplicam a:
- Novo código escrito
- Refatorações
- Revisões de código
- Documentação

Qualquer agente de IA deve seguir estes padrões automaticamente, sem necessidade de perguntar.