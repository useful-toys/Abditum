# TDR 002 — Formato de Armazenamento Interno

**Status**: Aceita
**Data**: 2026-03-02

## Contexto

O payload do cofre precisa ser serializado antes de ser cifrado. A escolha do formato afeta legibilidade, dependências, tamanho e facilidade de evolução do schema.

## Requisito fundamental

O formato deve ser **legível por qualquer pessoa que possua a chave**, sem necessidade de ferramentas especiais, bibliotecas ou conhecimento de schema externo. Decifrar o arquivo e abrir em qualquer editor de texto deve ser suficiente para inspecionar o conteúdo.

Esse requisito é uma extensão do princípio de controle total pelo usuário: quem tem a senha tem acesso irrestrito aos seus próprios dados — não depende da aplicação estar instalada, de uma lib estar disponível ou de um formato proprietário ser suportado.

Esse critério descarta diretamente:
- **Protocol Buffers**: requer schema `.proto` e toolchain para ser interpretado
- **MessagePack / CBOR**: binários — ilegíveis sem ferramenta de conversão
- **SQLite**: introduz camada de banco de dados, DAOs e engine de query — peso e complexidade desnecessários para uma estrutura de árvore simples
- **TOML / YAML / INI**: adequados para configuração, não para serialização de árvore recursiva de objetos heterogêneos; YAML tem spec de 80+ páginas e parsers inconsistentes

## Decisão

**JSON** com `encoding/json` da stdlib Go.

- Texto puro legível em qualquer editor após decifração
- Zero dependências externas de serialização
- Suporta naturalmente a estrutura recursiva de `Node` com filhos
- Schema evolution via campos opcionais (`omitempty`) sem quebrar versões anteriores
- Tamanho adicional é irrelevante — cofres pessoais ficam tipicamente abaixo de 100 KB de JSON

A estrutura usa `omitempty` para minimizar o payload em runtime.

## Alternativas consideradas e descartadas

| Opção             | Motivo de descarte                                                          |
|-------------------|-----------------------------------------------------------------------------|
| MessagePack       | Binário — ilegível sem ferramenta; viola requisito fundamental              |
| CBOR              | Binário — ilegível sem ferramenta; viola requisito fundamental              |
| Protocol Buffers  | Requer schema `.proto` e toolchain externo; viola requisito fundamental     |
| SQLite            | Engine de banco embutida, DAOs, peso desnecessário para árvore simples      |
| YAML              | Spec ambígua, parsers inconsistentes, inadequado para dados serializados    |
| TOML              | Verboso para árvores recursivas; sem representação natural para `Node[]`    |
| INI               | Sem tipos, sem aninhamento — inadequado para qualquer estrutura complexa    |

## Consequências

- Zero dependências extras de serialização
- Qualquer pessoa com a senha pode decifrar e inspecionar o conteúdo em texto puro
- Evolução de schema sem quebrar cofres antigos (campos novos opcionais são ignorados)
- A versão do formato no header permite detectar cofres de versões anteriores e aplicar migração se necessário

