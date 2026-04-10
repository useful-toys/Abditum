# Modelo de Domínio — Abditum

## Princípios de Modelagem

- **Hierarquia recursiva**: a raiz do cofre é a Pasta Geral. Pastas podem conter segredos e subpastas em qualquer nível de aninhamento.
- **Ordenação por posição**: Pasta mantém duas listas separadas (subpastas e segredos); ambas preservam ordem definida pelo usuário. Campos em Segredo e ModeloSegredo também preservam ordem. Ordem é persistida e restaurada ao carregar.
- **Modelo como snapshot**: segredos criados a partir de modelos não mantêm vínculo — no segredo, o nome do modelo é apenas histórico.
- **Campos uniformes**: valor é sempre UTF-8. Em memória: `[]byte` para zeragem; persistido: string legível via marshal customizado.
- **Observação implícita**: todo segredo tem campo automático "Observação" (tipo texto, última posição, não-deletável).
- **Busca sequencial**: sem índices — varreduras sobre estrutura em memória.
- **Estado de sessão**: cada `Segredo` carrega um estado transiente (não persistido) que rastreia sua condição em relação ao arquivo: `original`, `incluido`, `modificado` ou `excluido`. Esse estado governa comportamentos do domínio (busca, exportação) e a observabilidade pelo usuário.
- **Configurações embutidas**: tempos de bloqueio, ocultação e limpeza armazenados no arquivo do cofre.

## Classificação dos Tipos

### Agregado Raiz

| Tipo   | Justificativa                                                                                                   |
|--------|-----------------------------------------------------------------------------------------------------------------|
| Cofre  | Ponto de entrada único para toda mutação do domínio. Nenhuma entidade interna é modificada fora do contexto do cofre. Toda persistência é feita sobre o cofre como unidade atômica. |

### Entidades

Têm identidade baseada em nomes (composite keys ou nome simples). Em DDD/Go, igualdade é determinada pela identidade semântica.

| Entidade       | Identidade                    |
|--------|----------------------------------|
| Pasta          | (parentId, nome)               |
| Segredo        | (pastaId, nome)                |
| ModeloSegredo  | nome                           |

### Objetos de Valor

Sem identidade própria. Definidos inteiramente pelos seus atributos. São sempre parte de uma entidade.

| Objeto de Valor    | Pertence a     | Observação                                                              |
|--------------------|----------------|-------------------------------------------------------------------------|
| CampoSegredo       | Segredo        | Identidade = posição (índice na lista). Nomes podem ser duplicados (sem restrição). |
| CampoModeloSegredo | ModeloSegredo  | Identidade = posição (índice na lista). Nomes podem ser duplicados (sem restrição). |
| Configuracoes      | Cofre          | Instância única. Tempos de bloqueio, ocultação e limpeza.                   |

---

## Regras de Identidade e Unicidade

**Pasta e Segredo** (identidade composite key):
- Nome deve ser único dentro do container pai (pasta para Segredo; pasta para Pasta)
- Renomeação muda a identidade semântica
- Mover/renomear com colisão — renomeação automática com sufixo numérico (ex: "Login (1)", "Login (2)")

**ModeloSegredo** (identidade: nome):
- Nome deve ser único globalmente no cofre
- Renomeação muda a identidade
- Renomear com colisão — renomeação automática com sufixo numérico

---

## Regras de Estado de Sessão

Todo `Segredo` carrega um `estado_sessao` transiente (não serializado no arquivo) que reflete sua condição em relação ao estado persistido.

| Estado       | Significado                                                       |
|--------------|-------------------------------------------------------------------|
| `original`   | Carregado do arquivo sem alterações na sessão                    |
| `incluido`   | Criado durante a sessão; não existe no arquivo ainda             |
| `modificado` | Existia no arquivo; foi alterado durante a sessão                |
| `excluido`   | Marcado para remoção; será suprimido do cofre após sucesso no salvamento |


**Transições:**
- Ao abrir ou descartar: todos os segredos iniciam como `original`
- Criar segredo: → `incluido`
- Alterar nome ou campos de segredo `original`: → `modificado` (favoritar não altera o estado do segredo, apenas do cofre)

- Alterar segredo `incluido`: permanece `incluido`
- Marcar para exclusão (qualquer estado): → `excluido`; estado anterior memorizado para eventual restauração
- Desmarcar exclusão: restaura o estado anterior à marcação

**Efeitos sobre o domínio:**
- `excluido`: excluído da busca; excluído da exportação; removido permanentemente da árvore em memória apenas após a confirmação de sucesso do salvamento (commit)

- `incluido`, `modificado`, `excluido`: exibem indicador visual na UI
- `original`: sem indicador

---

### Cofre

Agregado raiz que encapsula todo o cofre de senhas. Ponto de entrada único para mutação do domínio; nenhuma entidade interna é modificada fora do contexto do cofre. Toda persistência é feita sobre o cofre como unidade atômica.

| Atributo                    | Tipo                    | Descrição                                              |
|-----------------------------|-------------------------|--------------------------------------------------------|
| configuracoes               | Configuracoes           | Configurações operacionais                             |
| pasta_geral                 | Pasta                   | Raiz da hierarquia. Todo segredo vive dentro de uma Pasta. |
| data_criacao                | datetime                | Data/hora de criação do cofre                          |
| data_ultima_modificacao     | datetime                | Data/hora da última modificação persistida             |

### Configuracoes

Objeto de valor que concentra as preferências operacionais do cofre (tempos de bloqueio, ocultação e limpeza). Instância única por cofre; imutável durante execução (valores carregados do arquivo e persistidos en bloc ao salvar).

| Atributo                                     | Tipo    | Padrão | Descrição                                    |
|----------------------------------------------|---------|--------|----------------------------------------------|
| tempo_bloqueio_inatividade_minutos           | inteiro | 5      | Tempo até bloqueio automático por inatividade |
| tempo_ocultar_segredo_segundos               | inteiro | 15     | Tempo até reocultação de campo sensível        |
| tempo_limpar_area_transferencia_segundos     | inteiro | 30     | Tempo até limpeza automática da clipboard      |

Nenhum temporizador pode ser desabilitado — todos são obrigatórios.

### Pasta

Container hierárquico que agrupa segredos e outras pastas. Identidade é (parentId, nome); nome único entre irmãs.

| Atributo  | Tipo          | Descrição                                                  |
|-----------|---------------|------------------------------------------------------------|
| nome      | string        | Nome da pasta                                              |
| parentId  | string        | Referência ao pai (nulo para Pasta Geral)                  |
| pastas    | list[Pasta]   | Subpastas (exibidas primeiro, ordem preservada)            |
| segredos  | list[Segredo] | Segredos diretos (exibidos depois, ordem preservada)       |

**Pasta Geral**: raiz (parentId nulo); não pode ser renomeada, movida ou excluída.

### Segredo

Credencial ou informação confidencial armazenada dentro de uma pasta. Identidade é (pastaId, nome); nome único dentro da pasta pai.

| Atributo                | Tipo               | Descrição                                                         |
|-------------------------|--------------------|-------------------------------------------------------------------| 
| nome                    | string             | Nome do segredo                                                   |
| nome_modelo_segredo     | string (opcional)  | Histórico: qual modelo foi usado na criação                       |
| campos                  | list[CampoSegredo] | Campos em ordem definida pelo usuário (Observação sempre última)  |
| favorito                | booleano           | Marca segredo como favorito                                        |
| data_criacao            | datetime           | Quando foi criado                                                  |
| data_ultima_modificacao | datetime           | Última alteração                                                   |
| estado_sessao           | enum (transiente)  | `original` / `incluido` / `modificado` / `excluido` — não persistido |

### CampoSegredo

Objeto de valor que representa um campo individual dentro de um Segredo. Identidade é determinada por posição (índice) na lista; nomes podem ser duplicados sem restrição. Tipo define comportamento (sensível sofre ocultação automática).

| Atributo | Tipo                         | Descrição                            |
|----------|------------------------------|--------------------------------------|
| nome     | string                       | Nome do campo (sem restrição unicidade) |
| tipo     | enum: texto, texto_sensivel  | Define comportamento (visibilidade)  |
| valor    | []byte (texto UTF-8)         | Sempre UTF-8; zerável em memória      |

**Observação**: CampoSegredo especial (nome fixo "Observação", tipo texto, última posição, não-deletável).

### ModeloSegredo

Estrutura reutilizável de campos para agilizar criação de segredos. Identidade é o nome (único globalmente no cofre).

| Atributo | Tipo                      | Descrição                                |
|----------|---------------------------|------------------------------------------|
| nome     | string                    | Nome do modelo (único globalmente)        |
| campos   | list[CampoModeloSegredo]  | Estrutura de campos (ordem preservada)   |

**Exibição**: ordem alfabética (não-reordenável).

### CampoModeloSegredo

Objeto de valor que define a estrutura de um campo no ModeloSegredo. Identidade é determinada por posição (índice) na lista; nomes podem ser duplicados. Tipo é imutável à criação (define o template).

| Atributo | Tipo                         | Descrição               |
|----------|------------------------------|-------------------------|
| nome     | string                       | Nome do campo           |
| tipo     | enum: texto, texto_sensivel  | Tipo (mutável — pode ser alterado após a criação) |

---

## Pastas Virtuais

Pastas virtuais são **vistas derivadas** do estado em memória. Não são persistidas no arquivo.

| Pasta Virtual | Definição                                                                                                                                                  |
|---------------|------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Favoritos** | Conjunto de segredos com `favorito = true`, percorrido em profundidade seguindo a ordem do JSON. Exibida na árvore como nó irmão da Pasta Geral (acima dela). Somente leitura — não suporta criação, movimentação ou exclusão de segredos diretamente a partir desta vista. Não pode ser renomeada, movida ou excluída. |

---




































































