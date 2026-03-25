# Glossário — Linguagem Ubíqua do Abditum

Este glossário define os termos canônicos usados em toda a documentação BDD, no código-fonte e nas conversas sobre o produto. Todos os cenários, regras e narrativas devem usar exclusivamente estes termos.

---

## Entidades principais

| Termo | Definição |
|-------|-----------|
| **Cofre** | Arquivo criptografado (`.abditum`) que armazena segredos, pastas, modelos de segredo e configurações do usuário. Apenas um cofre pode estar ativo por vez. |
| **Segredo** | Item individual dentro do cofre, composto por nome, campos (comuns e sensíveis), observação e marcação de favorito. |
| **Pasta** | Contêiner estrutural que agrupa segredos e subpastas na hierarquia do cofre. Pastas podem conter outras pastas em qualquer nível de aninhamento. |
| **Raiz do cofre** | Nível estrutural primário e mais alto da hierarquia — contém os segredos e pastas que não estão aninhados em outras pastas. Funciona como uma pasta sem nome. |
| **Modelo de segredo** | Estrutura reutilizável que define um conjunto de campos (nome e tipo) para criação de novos segredos. Segredos criados a partir de um modelo não mantêm vínculo por referência com ele. |
| **Campo de segredo** | Elemento individual dentro de um segredo, com nome, tipo e valor. |
| **Campo de modelo de segredo** | Elemento individual dentro de um modelo de segredo, com nome e tipo. Serve como template para os campos do segredo criado a partir do modelo. |
| **Senha mestra** | Chave de acesso ao cofre, usada para derivar a chave de criptografia. Conhecida apenas pelo usuário (Conhecimento Zero). |

## Tipos de campo

| Termo | Definição |
|-------|-----------|
| **Texto** | Tipo de campo cujo valor é visível por padrão e participa da busca. |
| **Texto sensível** | Tipo de campo cujo valor é oculto por padrão, exibido apenas sob ação explícita do usuário, e **nunca** participa da busca. |

## Atributos e dados do segredo

| Termo | Definição |
|-------|-----------|
| **Observação** | Campo de texto livre associado a todo segredo. É tratado como dado não sensível e participa da busca. Não deve ser usado para armazenar segredos. |
| **Segredo favorito** | Segredo marcado pelo usuário para destaque visual e presença na pasta virtual de Favoritos. |

## Elementos virtuais da interface

| Termo | Definição |
|-------|-----------|
| **Pasta virtual** | Agrupamento lógico gerado pelo sistema que exibe segredos de outras localizações sem alterar sua posição real na hierarquia. |
| **Favoritos** | Pasta virtual visível no topo da raiz (quando houver favoritos) que lista atalhos para os segredos favoritados. |
| **Lixeira** | Pasta virtual visível no final da raiz (quando houver segredos excluídos reversivelmente) que materializa segredos aguardando exclusão definitiva. |

## Estados do cofre

| Termo | Definição |
|-------|-----------|
| **Cofre Salvo** | Cofre sincronizado com o arquivo — não há divergência entre memória e último salvamento. |
| **Cofre Modificado** | Cofre com divergência entre o estado em memória e o último estado salvo no arquivo. |
| **Bloqueio do cofre** | Processo de proteção que interrompe o acesso ao conteúdo, limpa buffers sempre que possível, e retorna ao fluxo de abertura do cofre exigindo nova autenticação. Não é um estado observável separado. |

## Estados do segredo

| Termo | Definição |
|-------|-----------|
| **Segredo disponível** | Segredo visível na hierarquia principal, elegível para navegação, edição, movimentação e cópia. |
| **Segredo novo** | Segredo criado na sessão atual, confirmado mas não persistido. Continua neste estado até o próximo salvamento. |
| **Segredo modificado** | Segredo previamente persistido que sofreu alteração confirmada e ainda não foi salvo. |
| **Segredo excluído reversivelmente** | Segredo retirado da hierarquia principal e materializado na Lixeira, restaurável até o próximo salvamento. Não pode ser editado neste estado. |
| **Segredo restaurado** | Segredo anteriormente excluído reversivelmente e reinserido na hierarquia principal, retornando ao estado que possuía antes da exclusão. |

## Estados de exposição de dados sensíveis

| Termo | Definição |
|-------|-----------|
| **Campo sensível oculto** | Estado padrão de exibição para campos do tipo texto sensível — valor mascrado. |
| **Campo sensível exibido temporariamente** | Estado temporário após ação explícita do usuário, encerrado manualmente ou por temporizador configurado. |
| **Área de transferência povoada** | Existe um valor copiado aguardando limpeza automática por temporizador ou por bloqueio/fechamento do cofre. |

## Operações

| Termo | Definição |
|-------|-----------|
| **Exclusão reversível (Soft Delete)** | Mecanismo pelo qual um segredo excluído permanece restaurável até o próximo salvamento definitivo do cofre. |
| **Exclusão física** | Remoção definitiva e irrecuperável. Aplica-se a pastas (sempre) e a segredos na Lixeira ao salvar o cofre. |
| **Promoção de filhos** | Ao excluir uma pasta, seus segredos e subpastas são movidos para a pasta pai (ou raiz). |
| **Salvamento atômico** | Gravação do cofre em arquivo `.abditum.tmp` seguida de renomeação para o nome final, garantindo que o arquivo original não seja corrompido em caso de falha. |

## Modelos pré-definidos

| Termo | Definição |
|-------|-----------|
| **Modelo "Login"** | Modelo pré-definido com campos: URL, Username, Password. |
| **Modelo "Cartão de Crédito"** | Modelo pré-definido com campos: Número do Cartão, Nome no Cartão, Data de Validade, CVV. |
| **Modelo "API Key"** | Modelo pré-definido com campos: Nome da API, Chave de API. |

## Conceitos de segurança

| Termo | Definição |
|-------|-----------|
| **Conhecimento Zero (Zero Knowledge)** | Princípio em que a aplicação não possui meios de acessar ou recuperar os dados sem a senha mestra do usuário. |
| **Shoulder Surfing** | Técnica de espionagem física onde um indivíduo observa a tela do usuário para roubar informações visíveis. Mitigada por atalho de ocultação rápida da interface. |
| **AES-256-GCM** | Algoritmo de criptografia simétrica usado para cifrar os dados do cofre. |
| **Argon2id** | Algoritmo de derivação de chave a partir da senha mestra, com custo alto de memória e tempo para proteção contra ataques de força bruta. |

## Siglas

| Sigla | Significado |
|-------|-------------|
| **TUI** | Text User Interface — interface de texto interativa que ocupa todo o terminal. |
| **AAD** | Additional Authenticated Data — dados adicionais autenticados pelo AES-256-GCM (inclui o cabeçalho do arquivo). |
| **TOTP** | Time-based One-Time Password — código numérico temporário. Funcionalidade postergada para v2. |
