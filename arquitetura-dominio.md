# Arquitetura da Camada de Domínio — Abditum

## Visão Geral

Este documento descreve as decisões arquiteturais tomadas para a camada de domínio do Abditum, cobrindo a estrutura do agregado `Cofre`, o padrão `Manager`, as relações pai-filho, a navegação, a mutação e o encapsulamento.

---

## 1. Encapsulamento — A Regra Fundamental

A unidade de encapsulamento em Go é o **pacote**, não a classe. Tudo que vive em `internal/vault` compartilha o mesmo espaço de visibilidade. Tudo fora de `internal/vault` só acessa o que for exportado.

Essa característica da linguagem é usada deliberadamente para garantir que **somente o `Manager` pode mutar as entidades do domínio**.

### Como funciona na prática

Todos os campos das entidades são declarados em **minúscula** (privados ao pacote):

```go
// internal/vault/segredo.go
type Segredo struct {
    nome                  string
    pasta                 *Pasta        // referência ao pai — não persistida
    favorito              bool
    campos                []CampoSegredo
    nomeModeloSegredo     string
    estadoSessao          EstadoSessao  // transiente — não persistido
    dataCriacao           time.Time
    dataUltimaModificacao time.Time
}
```

A TUI, que vive em `internal/tui`, **não pode** escrever diretamente em nenhum campo:

```go
// internal/tui — ERRO DE COMPILAÇÃO:
segredo.nome = "outro"     // campo privado ao pacote vault
segredo.favorito = true    // campo privado ao pacote vault
```

A TUI só pode **ler** via getters exportados:

```go
// internal/vault/segredo.go — getters exportados
func (s *Segredo) Nome() string              { return s.nome }
func (s *Segredo) Favorito() bool            { return s.favorito }
func (s *Segredo) Pasta() *Pasta             { return s.pasta }
func (s *Segredo) EstadoSessao() EstadoSessao { return s.estadoSessao }
func (s *Segredo) DataCriacao() time.Time    { return s.dataCriacao }


// Retorna cópia da slice para evitar mutação indireta da lista
func (s *Segredo) Campos() []CampoSegredo {
    copia := make([]CampoSegredo, len(s.campos))
    copy(copia, s.campos)
    return copia
}
```

Mutações só são possíveis de dentro de `internal/vault`, onde o `Manager` opera.

### Acesso ao valor de campos: ValorComoString()

`CampoSegredo.valor` é um `[]byte` privado ao pacote — a TUI não pode acessá-lo diretamente. Para exibir o valor na tela, `CampoSegredo` expõe um método explicitamente nomeado:

```go
// ValorComoString converte o valor do campo para string para exibição na TUI.
//
// ATENÇÃO: esta conversão é irreversível do ponto de vista de segurança.
// A string retornada é imutável em Go e não pode ser zerada da memória.
// Chame este método apenas quando o usuário explicitamente solicitar
// a visualização do valor. Nunca armazene o resultado em estado persistente
// do componente TUI — use-o diretamente na renderização e descarte.
func (c *CampoSegredo) ValorComoString() string {
    return string(c.valor)
}
```

A decisão de embutir a conversão `[]byte` → `string` na entidade, em vez de expor o `[]byte` bruto, é deliberada:

- A conversão vai acontecer de qualquer forma — a TUI precisa de string para renderizar no Bubble Tea
- Centralizar em um único método torna o ponto de exposição auditável e visível em code review
- O nome explícito `ValorComoString()` é um sinal visual claro no código da TUI de que um dado sensível está sendo revelado
- Futuras necessidades de logging, auditoria ou rate limiting têm um ponto de interceptação natural

---

## 2. Relações Pai-Filho

### Modelagem

A árvore de pastas é navegável nos **dois sentidos**. Cada `Pasta` e cada `Segredo` carregam uma referência ao seu pai:

```go
// internal/vault/pasta.go
type Pasta struct {
    nome     string
    pai      *Pasta     // nil para Pasta Geral — não persistido
    pastas   []*Pasta
    segredos []*Segredo
}

// internal/vault/segredo.go
type Segredo struct {
    nome  string
    pasta *Pasta     // não persistido
    // ...
}
```

Os campos `pai` e `pasta` são **não-persistidos** — existem apenas em memória durante a sessão. Isso não é uma infração do DDD: a relação pai-filho é um fato do domínio, e que ela seja reconstituída em memória ao carregar é apenas um detalhe de implementação, não uma violação conceitual. O Abditum já adota o mesmo padrão para `estadoSessao` em `Segredo`.

### Por que não usar IDs para referenciar o pai

IDs existem para referenciar entidades através de fronteiras que apagam a estrutura: rede, banco de dados, serialização entre processos. O Abditum não atravessa nenhuma dessas fronteiras durante uma sessão — o cofre inteiro vive em memória como um grafo de ponteiros Go.

Um `*Pasta` em Go **é** um identificador: globalmente único dentro do processo, estável durante toda a sessão, com lookup O(1) garantido pelo hardware. IDs seriam uma segunda camada de indireção sobre algo que já tem indireção nativa — complexidade sem benefício.

### Serialização

O JSON do cofre é uma estrutura hierárquica aninhada que reflete diretamente a hierarquia do domínio. Pastas contêm subpastas e segredos aninhados. Os campos `pai` e `pasta` são ignorados na serialização — a relação é implícita pela estrutura do documento.

```json
{
  "pasta_geral": {
    "nome": "Geral",
    "pastas": [
      {
        "nome": "Sites e Apps",
        "pastas": [],
        "segredos": [
          { "nome": "Gmail", "campos": [...] }
        ]
      }
    ]
  }
}
```

### Reconstituição após deserialização

Após o JSON ser deserializado, uma única passagem recursiva popula todas as referências pai-filho:

```go
// internal/vault/serialization.go — chamado uma vez após deserializar
func popularReferencias(pasta *Pasta, pai *Pasta) {
    pasta.pai = pai  // nil para Pasta Geral
    for _, subpasta := range pasta.pastas {
        popularReferencias(subpasta, pasta)
    }
    for _, segredo := range pasta.segredos {
        segredo.pasta = pasta
    }
}

// Ponto de entrada:
popularReferencias(cofre.pastaGeral, nil)
```

Complexidade O(n) no total de nós. Executado uma única vez ao abrir o cofre.

---

## 3. O Agregado Cofre

`Cofre` é o agregado raiz — concentra o estado global da sessão: `modificado`, `dataUltimaModificacao`, a raiz da hierarquia de pastas e a lista de modelos. A coordenação de mutações pertence ao `Manager`, que chama métodos privados das entidades e atualiza o estado do `Cofre` diretamente.

```go
// internal/vault/cofre.go
type Cofre struct {
    configuracoes           Configuracoes
    pastaGeral              *Pasta
    modelos                 []*ModeloSegredo
    modificado              bool      // true se há alterações não salvas
    dataCriacao             time.Time
    dataUltimaModificacao   time.Time
}
```

### Divisão de responsabilidades: Manager vs. Entidades

O critério que determina onde a lógica vive:

> **O método público vive sempre no `Manager`. A lógica de negócio (validação de invariantes, construção, mutação) vive em métodos privados das entidades. Para operações que cruzam fronteiras entre entidades, a orquestração fica no próprio `Manager`.**

| Operação | Método público em | Lógica de negócio em | Invariante cruzado? |
|---|---|---|---|
| Criar segredo | `Manager` | `Pasta` (privado) | ❌ local à pasta |
| Criar segredo de modelo | `Manager` | `Pasta` (privado) | ❌ local à pasta |
| Reordenar segredo | `Manager` | `Segredo` (privado) | ❌ local à pasta |
| Reordenar pasta | `Manager` | `Pasta` (privado) | ❌ local à pasta |
| Alterar configuração | `Manager` | `Manager` | ❌ global |
| Renomear segredo | `Manager` | `Segredo` (privado) | ❌ local à pasta pai |
| Favoritar segredo | `Manager` | `Segredo` (privado) | ❌ sem invariante |
| Criar pasta | `Manager` | `Pasta` (privado) | ❌ local à pasta pai |
| Mover segredo | `Manager` | `Manager` | ✅ duas pastas |
| Mover pasta | `Manager` | `Manager` | ✅ ciclo + duas pastas |
| Excluir pasta | `Manager` | `Manager` | ✅ promove filhos para pai |
| Criar modelo | `Manager` | `Manager` | ✅ unicidade global |

### O padrão: público no Manager, privado na entidade

Métodos públicos vivem no `Manager` — são a única API que a TUI conhece. A lógica de negócio (validação de invariantes, construção da entidade, mutação de estado) vive em métodos privados das entidades, acessíveis dentro do pacote `internal/vault`.

O `Manager` chama os métodos privados da entidade e, após confirmação de mudança real, atualiza diretamente os campos de estado do `Cofre`.

```go
// internal/vault/manager.go — método público, orquestrador
func (m *Manager) CriarSegredo(pasta *Pasta, nome string, modelo *ModeloSegredo) (*Segredo, error) {
    if m.bloqueado {
        return nil, ErrCofreBloqueado
    }
    if err := pasta.validarCriacaoSegredo(nome, modelo); err != nil { // invariante local à pasta
        return nil, err
    }
    segredo := pasta.criarSegredo(nome, modelo)   // factory privado da entidade
    now := time.Now().UTC()
    segredo.dataUltimaModificacao = now
    m.cofre.modificado = true                     // Manager atualiza estado do Cofre diretamente
    m.cofre.dataUltimaModificacao = now
    return segredo, nil
}

// internal/vault/entities.go — método privado, factory local
func (p *Pasta) criarSegredo(nome string, modelo *ModeloSegredo) *Segredo {
    segredo := &Segredo{
        nome:         nome,
        pasta:        p,                          // referência ao pai populada aqui
        estadoSessao: EstadoIncluido,
        dataCriacao:  time.Now().UTC(),
    }
    // popula campos a partir do modelo...
    p.segredos = append(p.segredos, segredo)
    return segredo
}

```

### Por que o estado `modificado` fica no Cofre

O estado `modificado` pertence ao `Cofre` — não ao `Manager` — porque é um dado do domínio, não da sessão de aplicação. O `Cofre` é serializado; o `Manager` não. Ao persistir, o `storage` recebe o `Cofre` e o serializa — o indicador de "há alterações não salvas" precisa ser consultável a partir dele. O `Manager` atualiza esse campo após cada mutação real e o consulta quando necessário:

```go
// Manager consulta antes de sair ou bloquear
func (m *Manager) TemAlteracoesPendentes() bool {
    return m.cofre.modificado
}
```

### Quem popula referências ao construir

A responsabilidade de popular referências pai-filho (`segredo.pasta`, `subpasta.pai`) pertence a quem constrói a entidade:

- Na **deserialização**: o `storage` popula via `popularReferencias()` após carregar o JSON
- Na **criação durante a sessão**: o método privado da entidade popula no momento da construção

O Manager nunca toca nessa responsabilidade.

---

## 4. O Padrão Manager

O `Manager` é a camada de serviço de aplicação. É a única API pública que a TUI conhece. Ele **não contém lógica de negócio** — apenas orquestra o fluxo entre TUI, domínio e storage.

```go
// internal/vault/manager.go
type Manager struct {
    cofre      *Cofre
    repositorio RepositorioCofre  // interface para internal/storage
}
```

### Duas categorias de operação: mutação em memória vs. persistência explícita

O Manager distingue claramente duas categorias de operação:

**Operações de mutação** — alteram o estado em memória, marcam o cofre como modificado, mas **não salvam no arquivo**. O salvamento é responsabilidade do usuário, que aciona explicitamente quando quiser.

```
TUI → Manager.RenomearSegredo(segredo, "novo nome")
         ↓ validação do estado bloqueado
         ↓ segredo.validarRenomear("novo nome")   — invariante local à entidade
         ↓ alterado, err := segredo.renomear("novo nome")   — mutação privada
         ↓    estadoSessao atualizado internamente pela entidade
         ↓ if alterado:
         ↓    cofre.modificado = true
         ↓    cofre.dataUltimaModificacao = now
         ↓ retorna erro ou nil
      retorna erro ou nil para a TUI
      — arquivo não tocado —
```

**Operações de persistência** — acionadas explicitamente pelo usuário ou por eventos do sistema (bloqueio, saída, alteração de senha mestra). Estas sim invocam `repositorio.Salvar()`:

```
TUI → Manager.Salvar()
         ↓
      repositorio.Salvar(cofre)
         ↓ serialização JSON (filtra excluídos)
         ↓ criptografia AES-256-GCM
         ↓ escrita atômica no arquivo
         ↓ se sucesso:
         ↓    cofre.efetivarMutacoes() // remove excluídos e reseta estados
         ↓    cofre.modificado = false
      retorna erro ou nil para a TUI
```

A ordem é crítica: o cofre em memória só é limpo (remoção dos ponteiros de segredos excluídos e reset dos indicadores de modificação) **após** a confirmação de que o arquivo foi escrito com sucesso no disco. Se o salvamento falhar, o estado em memória permanece intacto, permitindo ao usuário tentar novamente ou corrigir o problema sem perder a marcação do que seria excluído.

### Operações que salvam automaticamente

Pela spec, há exatamente uma operação que salva automaticamente sem ação explícita do usuário: **alterar a senha mestra**. A alteração é imediata e irrevogável — o cofre é regravado na hora com a nova senha, incluindo todas as alterações pendentes da sessão.

```go
// Única operação que persiste automaticamente
func (m *Manager) AlterarSenhaMestra(novaSenha []byte) error
```

Todas as demais operações de mutação apenas atualizam o estado em memória. O Manager expõe `TemAlteracoesPendentes()` para que a TUI saiba quando oferecer ou exigir salvamento:

```go
func (m *Manager) TemAlteracoesPendentes() bool {
    return m.cofre.modificado
}
```

### Operações que retornam efeitos colaterais

Algumas operações têm sucesso mas produzem efeitos que a TUI precisa comunicar ao usuário — como renomeações automáticas ao excluir uma pasta. Nesses casos, o Manager retorna os efeitos como parte do resultado, não como erro:

```go
// Renomeações são efeito colateral de sucesso — não são erro
func (m *Manager) ExcluirPasta(pasta *Pasta) ([]Renomeacao, error)

type Renomeacao struct {
    De   string
    Para string
}
```

### Erros do domínio

Erros de validação são **sentinel errors** declarados em `internal/vault/errors.go`:

```go
var (
    ErrNameConflict        = errors.New("nome já existe")
    ErrPastaGeralProtected = errors.New("pasta geral não pode ser modificada")
    ErrCycleDetected       = errors.New("operação criaria ciclo na hierarquia")
    ErrObservacaoReserved  = errors.New("nome 'Observação' é reservado")
)
```

A TUI detecta com `errors.Is()`:

```go
err := manager.CriarPasta(pai, nome)
if errors.Is(err, ErrNameConflict) {
    // exibir mensagem amigável ao usuário
}
```

---

## 5. Navegação pela TUI

A TUI recebe ponteiros diretamente para as entidades do domínio (`*Pasta`, `*Segredo`, `*ModeloSegredo`). Ela navega a árvore chamando os getters exportados:

```go
// Navegação para baixo
pasta.Subpastas()   // []*Pasta
pasta.Segredos()    // []*Segredo

// Navegação para cima
pasta.Pai()         // *Pasta — nil se for Pasta Geral
segredo.Pasta()     // *Pasta
```

Quando a TUI precisa executar uma ação do usuário, ela chama o Manager com o ponteiro da entidade:

```go
// TUI captura intenção do usuário e chama o Manager
manager.RenomearSegredo(segredo, novoNome)
manager.MoverSegredo(segredo, pastaDestino)
manager.FavoritarSegredo(segredo)
```

O Manager repassa ao agregado `Cofre`, que valida e muta. A TUI nunca muta diretamente.

---

## 6. Observação — Invariante Estrutural

A `Observação` é um `CampoSegredo` especial que todo segredo possui automaticamente. Sua imutabilidade é garantida de duas formas:

**Structuralmente impossível de violar** para os casos principais: o método de deleção e renomeação de campos opera sobre os campos do usuário — a Observação simplesmente não está nessa slice manipulável.

**Validação explícita** para o único caso que requer checagem: ao adicionar campo a um `ModeloSegredo`, o agregado rejeita imediatamente qualquer campo com nome `"Observação"`.

```go
func (c *Cofre) AdicionarCampoModelo(modelo *ModeloSegredo, nome string, tipo TipoCampo) error {
    if nome == "Observação" {
        return ErrObservacaoReserved
    }
    // ...
}
```

---

## 7. Estado de Sessão e Rastreamento de Modificações

### Duas flags com semânticas distintas

O domínio mantém duas flags de estado relacionadas mas independentes:

| Flag | Onde vive | Semântica | Muda com favoritar? |
|---|---|---|---|
| `cofre.modificado` | `Cofre` | "há alterações não salvas no cofre" | ✅ sim |
| `segredo.estadoSessao` | `Segredo` | "o conteúdo deste segredo mudou na sessão" | ❌ não |

Favoritar é uma preferência de navegação, não uma edição de conteúdo. O indicador visual `modificado` no segredo comunica ao usuário que o *conteúdo* mudou — favoritar não altera conteúdo. Porém, favoritar ainda altera o estado do cofre e precisa ser salvo — por isso `cofre.modificado` é marcado.

### A regra de ouro: mudança real, não chamada de método

Nenhuma das duas flags muda pela simples chamada de um método de mutação. A flag só muda se o valor resultante for **realmente diferente** do valor atual.

Isso significa que um usuário que abre o diálogo de renomear, não altera nada e confirma não vê nenhum indicador de modificação — comportamento esperado e correto.

O mecanismo que viabiliza essa política é o retorno `bool` dos métodos privados das entidades, indicando se houve mudança real:

```go
// Entidade — método privado retorna se houve mudança real
func (s *Segredo) renomear(novoNome string) (alterado bool, err error) {
    if s.nome == novoNome {
        return false, nil          // mesmo valor — sem mudança
    }
    // validações...
    s.nome = novoNome
    return true, nil
}

// Manager — orquestra via método privado da entidade
func (m *Manager) RenomearSegredo(segredo *Segredo, novoNome string) error {
    if m.bloqueado {
        return ErrCofreBloqueado
    }
    if err := segredo.validarRenomear(novoNome); err != nil {
        return err
    }
    alterado, _ := segredo.renomear(novoNome)  // estadoSessao atualizado internamente pela entidade
    if alterado {
        now := time.Now().UTC()
        segredo.dataUltimaModificacao = now
        m.cofre.modificado = true
        m.cofre.dataUltimaModificacao = now
    }
    return nil
}

// Manager — favoritar: muda cofre mas não estadoSessao
func (m *Manager) AlternarFavoritoSegredo(segredo *Segredo) error {
    if m.bloqueado {
        return ErrCofreBloqueado
    }
    segredo.alternarFavorito()  // entidade gerencia sua própria flag
    now := time.Now().UTC()
    m.cofre.modificado = true   // cofre tem alterações não salvas
    m.cofre.dataUltimaModificacao = now
    // estadoSessao não muda — favoritar não é edição de conteúdo
    return nil
}

```

### Transições de estadoSessao

| De | Para | Quando |
|---|---|---|
| — | `original` | Ao abrir ou descartar o cofre |
| — | `incluido` | Ao criar o segredo na sessão |
| `original` | `modificado` | Qualquer mutação de conteúdo com valor realmente diferente |
| `incluido` | `incluido` | Mutações em segredo recém-criado — permanece `incluido` |
| qualquer | `excluido` | Ao marcar para exclusão; estado anterior memorizado |
| `excluido` | estado anterior | Ao desmarcar exclusão |

Segredos `incluido` permanecem `incluido` após mutações — não faz sentido marcar como `modificado` algo que ainda não existe no arquivo.

### Filtragem de segredos excluídos

`Pasta.Segredos()` retorna **todos** os segredos, incluindo os marcados como `excluido`. A TUI renderiza excluídos com visual diferente (strikethrough) — eles não são ocultados da listagem normal.

A filtragem de excluídos acontece apenas em casos de uso específicos que a spec define:

- **Busca** — excluídos não aparecem nos resultados
- **Exportação** — excluídos não são incluídos no arquivo exportado
- **Save** — segredos marcados como `excluido` são ignorados pela serialização. A remoção definitiva da árvore em memória só ocorre no passo de "Efetivação" (commit) após o sucesso da escrita em disco.

Não há método separado "incluindo excluídos" — o getter já retorna tudo, e cada caso de uso que filtra implementa sua própria política internamente.

### Busca de segredos

A busca é dividida entre duas responsabilidades:

**`Segredo.AtendeCriterio(criterio string) bool`** — o segredo sabe se seu próprio conteúdo casa com o critério. Aplica busca por substring normalizada (sem acentuação, case-insensitive) em: nome do segredo, nome de todos os campos, valor de campos comuns e observação. Valores de campos sensíveis nunca participam — apenas o nome do campo sensível.

**`Manager.BuscarSegredos(query string) []*Segredo`** — aplica a política de busca completa: percorre toda a árvore, filtra excluídos, delega a correspondência ao segredo.

```go
// Segredo — responsabilidade: "meu conteúdo casa com esse critério?"
func (s *Segredo) AtendeCriterio(criterio string) bool

// Manager — responsabilidade: política de busca (filtro de excluídos, percurso da árvore)
func (m *Manager) BuscarSegredos(query string) []*Segredo {
    var resultados []*Segredo
    m.cofre.percorrerSegredos(func(s *Segredo) {
        if s.estadoSessao == excluido {
            return  // excluídos não participam da busca
        }
        if s.AtendeCriterio(query) {
            resultados = append(resultados, s)
        }
    })
    return resultados
}
```

A TUI pode chamar `segredo.AtendeCriterio(query)` diretamente para verificar um segredo específico — por exemplo, para destacar resultados enquanto o usuário navega. O Manager usa o mesmo método internamente.

---

## 8. Resumo das Responsabilidades

| Camada | Responsabilidade |
|---|---|
| `internal/tui` | Renderização, captura de input, chamadas ao Manager |
| `internal/vault/Manager` | API pública do domínio: único ponto de entrada para toda mutação; chama métodos privados das entidades; atualiza estado do Cofre; controla quando persistir |
| `internal/vault/Cofre` | Agregado raiz: mantém estado global da sessão (`modificado`, `dataUltimaModificacao`, hierarquia de pastas, modelos) |
| `internal/vault/entidades` | Lógica local privada (validação de invariantes locais, construção); leitura pública via getters |
| `internal/vault/serialization.go` | Serialização e deserialização JSON do Cofre; reconstituição de referências pai-filho |
| `internal/storage` | Formato binário `.abditum`, criptografia do payload, escrita atômica, backup chain, detecção de mudança externa |
| `internal/crypto` | Derivação de chave (Argon2id), criptografia/descriptografia (AES-256-GCM) |

### O Manager como API do domínio

Toda mutação do domínio — sem exceção — passa por um método público do Manager. A TUI nunca chama métodos do `Cofre` ou das entidades diretamente. Esse é o contrato central da arquitetura:

```
TUI  →  Manager (API pública)  →  Entidades (métodos privados de validação e mutação)
         único ponto de                 estadoSessao gerenciado pela própria entidade
         entrada para
         toda mutação          →  cofre.modificado / cofre.dataUltimaModificacao
                                       atualizados pelo Manager após mudança real
```

Consequências práticas desse contrato:

- A interface pública do Manager é a lista completa de tudo que um usuário pode fazer no sistema — não há atalho, não há backdoor
- Adicionar uma nova operação significa sempre: novo método no Manager + (se necessário) novo método privado na entidade
- A TUI é completamente substituível — outro cliente (CLI, testes) usaria exatamente a mesma API do Manager

---

## 9. Serialização JSON — Consequência do Encapsulamento

### O problema

A serialização do `Cofre` para JSON e a deserialização de volta para o grafo de domínio precisam acessar os campos privados das entidades (`nome`, `campos`, `pai`, `estadoSessao`, etc.). Em Go, campos em minúscula são privados ao **pacote** — qualquer código fora de `internal/vault` não os enxerga.

Isso significa que a lógica de serialização **não pode viver em `internal/storage`**. O `storage` não tem visibilidade dos atributos internos das entidades.

Essa restrição é consequência direta da decisão arquitetural central do Abditum: campos privados para impedir que a TUI mute o estado do domínio (seção 1). O mesmo mecanismo que protege as entidades contra mutação externa também impede que pacotes externos as serializem.

### A solução: funções dedicadas em `vault/serialization.go`

A serialização vive em um arquivo separado dentro do pacote `vault` — `serialization.go`. Duas funções package-level concentram toda a lógica:

```go
// internal/vault/serialization.go

// SerializarCofre converte o Cofre para JSON (UTF-8).
// Acessa campos privados das entidades diretamente.
// Segredos com estadoSessao == excluido são omitidos.
func SerializarCofre(cofre *Cofre) ([]byte, error)

// DeserializarCofre reconstrói o grafo de domínio a partir de JSON.
// Popula campos privados e reconstitui referências pai-filho.
// Todos os segredos recebem estadoSessao = original.
func DeserializarCofre(data []byte) (*Cofre, error)
```

O `internal/storage` chama essas funções como parte do fluxo de persistência:

```
Save:  Manager → vault.SerializarCofre(cofre) → []byte JSON → crypto.Encrypt → storage.Write
Load:  storage.Read → crypto.Decrypt → []byte JSON → vault.DeserializarCofre(data) → *Cofre
```

### Por que não DTOs intermediários

A abordagem alternativa — criar structs DTO (Data Transfer Object) com campos exportados e tags `json:` — foi descartada por três motivos:

1. **Duplicação estrutural**: cada entidade do domínio precisaria de uma struct espelho com os mesmos campos em maiúscula. São 6 entidades (`Cofre`, `Pasta`, `Segredo`, `ModeloSegredo`, `CampoSegredo`, `Configuracoes`) × 2 structs = 12 tipos a manter sincronizados.

2. **Conversão bidirecional**: além das structs, seriam necessárias funções `toDTO` e `fromDTO` para cada entidade — código mecânico que apenas copia campo a campo e é fonte de bugs silenciosos quando um campo novo é adicionado à entidade mas esquecido no DTO.

3. **Sem benefício real**: DTOs fazem sentido quando há uma fronteira de transporte (rede, API pública) onde o formato externo diverge do modelo interno. No Abditum, o JSON é um formato de arquivo privado — não há consumidor externo. O formato JSON pode acompanhar o modelo de domínio sem fricção.

### Por que não `MarshalJSON` / `UnmarshalJSON` nas entidades

Implementar a interface `json.Marshaler` nas entidades acopla a lógica de serialização ao ciclo de vida da entidade. A entidade passaria a ter responsabilidade dupla: lógica de domínio + formato de persistência. Além disso, `encoding/json` invoca esses métodos automaticamente em qualquer `json.Marshal`, o que pode provocar serialização acidental em contextos de debug ou logging.

Funções package-level separadas mantêm a serialização **explícita e auditável** — ela só acontece quando alguém deliberadamente chama `SerializarCofre` ou `DeserializarCofre`.

### Reconstituição de referências

A deserialização reconstrói o JSON em structs com campos privados populados, mas as referências pai-filho (`pasta.pai`, `segredo.pasta`) não existem no JSON — a hierarquia é implícita pelo aninhamento. Uma passagem recursiva O(n) reconstitui todas as referências:

```go
// internal/vault/serialization.go — chamado dentro de DeserializarCofre
func popularReferencias(pasta *Pasta, pai *Pasta) {
    pasta.pai = pai
    for _, subpasta := range pasta.subpastas {
        popularReferencias(subpasta, pasta)
    }
    for _, segredo := range pasta.segredos {
        segredo.pasta = pasta
    }
}
```

Essa função vive em `serialization.go` porque precisa acessar `pasta.pai` e `segredo.pasta` — campos privados ao pacote.

---

## 10. Limitação Inerente: Zeragem de Memória em Go

Go não oferece tipos com semântica de zeragem garantida. Em Rust, `Drop` garante zeragem ao sair do escopo. Em Go, zeragem é imperativa, manual e não-garantida pelo runtime — é uma limitação estrutural da plataforma.

### O fluxo de exposição inevitável

Qualquer valor sensível que precisa ser exibido na tela percorre este caminho:

```
domínio              TUI                  Bubble Tea
[]byte zerado   →   string([]byte)   →   string composta
ao bloquear          ↑                    no modelo Elm
                     aqui nasce a         ↑
                     string imutável      pode haver
                     não-zerável          mais cópias
```

A conversão `string([]byte)` em Go **sempre aloca uma nova string** — é uma cópia por design da linguagem, porque `[]byte` é mutável e `string` é imutável e não podem compartilhar memória. A partir desse momento existem pelo menos dois buffers com o dado sensível:

- O `[]byte` original no domínio — **zerado ao bloquear** ✅
- A `string` criada pela TUI — **não-zerável, fora de controle** ❌

O Bubble Tea pode criar cópias adicionais ao compor a `View`.

### O que o Abditum zera e o que não consegue zerar

| Buffer | Zerado ao bloquear | Observação |
|---|---|---|
| Senha mestra (`[]byte`) | ✅ sim | Desde a leitura do terminal até o bloqueio |
| Chave AES derivada (`[]byte`) | ✅ sim | Zerada junto com a senha mestra |
| `CampoSegredo.valor` (`[]byte`) | ✅ sim | Todos os campos de todos os segredos |
| Strings efêmeras criadas pela TUI | ❌ não | Imutáveis em Go — não zeráveis |
| Cópias históricas por `append`/realocação | ❌ não | Órfãos inacessíveis no heap |

### Postura adotada

Zeragem é **melhor esforço**, não garantia — conforme documentado na spec. O investimento de engenharia está em zerar o que é controlável (os buffers do domínio), não em tentar controlar cópias efêmeras que inevitavelmente existirão durante a renderização.

O modelo de ameaça relevante para zeragem é um atacante que lê memória do processo **após** o bloqueio — via core dump, hibernação ou acesso físico à RAM. A janela de exposição que a zeragem reduz é real mas estreita: se o atacante consegue ler RAM arbitrária, provavelmente consegue também capturar o processo em execução antes do bloqueio, onde todos os dados estão necessariamente em memória.

---

## 11. Princípios que Guiam as Decisões

- **Manager é a única API de mutação**: toda alteração no domínio passa por um método público do Manager — sem exceção. A TUI nunca acessa `Cofre` ou entidades diretamente para mutar. Adicionar uma nova operação significa sempre: método no Manager + método no Cofre + método privado na entidade.
- **Invariantes impossíveis > invariantes verificados**: prefira modelar o domínio de forma que a violação seja impossível pela estrutura, não apenas detectada por validação.
- **Ponteiros Go como identidade**: dentro de um processo com grafo em memória, ponteiros são identificadores naturais — IDs seriam indireção sem benefício.
- **Encapsulamento por pacote**: campos privados + getters exportados é o mecanismo idiomático de Go para proteger o estado interno de mutação externa.
- **Cofre como estado global, não coordenador**: o `Cofre` é o agregado raiz que concentra o estado da sessão (`modificado`, `dataUltimaModificacao`, hierarquia, modelos). A coordenação de operações pertence ao `Manager`. Lógica de negócio local vive nas entidades (métodos privados).
- **Mudança real, não chamada de método**: flags de estado (`cofre.modificado`, `segredo.estadoSessao`) só mudam se o valor resultante for realmente diferente do atual. Métodos privados de mutação retornam `bool` para viabilizar essa política.
- **Responsabilidades separadas na busca**: a entidade sabe se seu conteúdo casa com um critério (`AtendeCriterio`); o Manager aplica a política de busca (filtragem de excluídos, percurso da árvore).
- **Serialização no pacote do domínio**: a lógica de serialização/deserialização JSON vive em `vault/serialization.go` — consequência direta do encapsulamento por campos privados. Sem DTOs intermediários, sem `MarshalJSON` nas entidades.
- **Quem constrói popula**: a responsabilidade de popular referências (pai, pasta) pertence a quem cria a entidade — `DeserializarCofre` na carga do arquivo, método privado da entidade na criação durante a sessão.
- **Manager não valida, não conhece regras**: toda lógica de negócio vive no domínio. O Manager é um coordenador de fluxo — recebe intenção da TUI, invoca métodos privados das entidades, atualiza estado do Cofre, persiste via storage.