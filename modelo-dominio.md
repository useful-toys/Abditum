# Modelo de Domínio — Abditum

## Princípios de Modelagem

- **Hierarquia recursiva**: a raiz do cofre é a Pasta Geral. Pastas podem conter segredos e subpastas em qualquer nível de aninhamento.
- **Ordenação por posição**: a ordem dos elementos no JSON reflete diretamente a ordem de exibição. Segredos aparecem antes de subpastas dentro de cada coleção. Campos seguem a ordem de inserção/reordenação pelo usuário.
- **Modelo como snapshot**: segredos criados a partir de modelos não mantêm vínculo por referência. O nome do modelo é guardado apenas como registro histórico. Não há distinção estrutural entre segredos criados com ou sem modelo.
- **Identidade por NanoID**: entidades com identidade persistida (Segredo, Pasta, ModeloSegredo) usam NanoID de 6 caracteres alfanuméricos. O espaço de 62⁶ (~56 bilhões) combinações garante unicidade prática sem coordenação central. O NanoID é diretamente serializável como string JSON e permanece estável através de importação, exportação, movimentação e migração de formato. O nome não é identificador — nomes repetidos são permitidos onde não explicitamente proibido.
- **Campos uniformes**: o valor de um campo é sempre uma string. String vazia representa campo existente e não preenchido — não há distinção de estado entre preenchido e vazio.
- **Observação implícita**: todo segredo possui um campo de observação que não é declarado em modelos e não pode ser removido. Ocupa sempre a última posição. É dado não sensível.
- **Busca sequencial em memória**: o cofre não mantém índices ou estruturas auxiliares de busca. Buscas são varreduras sequenciais sobre a estrutura carregada.
- **Configurações embutidas**: as configurações operacionais são armazenadas dentro do próprio arquivo do cofre, sem arquivos externos.

---

## Classificação dos Tipos

### Agregado Raiz

| Tipo   | Justificativa                                                                                                   |
|--------|-----------------------------------------------------------------------------------------------------------------|
| Cofre  | Ponto de entrada único para toda mutação do domínio. Nenhuma entidade interna é modificada fora do contexto do cofre. Toda persistência é feita sobre o cofre como unidade atômica. |

### Entidades

Têm identidade estável representada por NanoID. A identidade é independente do nome e persiste ao longo de operações de renomeação, movimentação e migração de formato.

| Entidade       | Identidade         | Observação                                                   |
|----------------|--------------------|--------------------------------------------------------------|
| Pasta          | NanoID (6 chars)   | Inclui a Pasta Geral, que é a raiz imutável da hierarquia   |
| Segredo        | NanoID (6 chars)   | Unidade principal de armazenamento de credenciais            |
| ModeloSegredo  | NanoID (6 chars)   | Estrutura reutilizável; alterações não afetam segredos existentes |

### Objetos de Valor

Sem identidade própria. Definidos inteiramente pelos seus atributos. São sempre parte de uma entidade.

| Objeto de Valor    | Pertence a     | Observação                                                              |
|--------------------|----------------|-------------------------------------------------------------------------|
| CampoSegredo       | Segredo        | Definido por nome + tipo + valor. Não tem ID. Ordenado por posição.     |
| CampoModeloSegredo | ModeloSegredo  | Definido por nome + tipo. Não tem ID. Ordenado por posição.             |
| Configuracoes      | Cofre          | Tempos de bloqueio, ocultação e limpeza de clipboard.                   |

---

## Estrutura do Domínio

### Cofre

Raiz agregada. Contém toda a estrutura persistida no arquivo `.abditum`.

| Atributo                    | Tipo                    | Descrição                                              |
|-----------------------------|-------------------------|--------------------------------------------------------|
| configuracoes               | Configuracoes           | Configurações operacionais                             |
| pasta_geral                 | Pasta                   | Raiz da hierarquia. Todo segredo vive dentro de uma Pasta. |
| modelos_segredo             | list[ModeloSegredo]     | Modelos disponíveis no cofre                           |
| data_criacao                | datetime                | Data/hora de criação do cofre                          |
| data_ultima_modificacao     | datetime                | Data/hora da última modificação persistida             |

### Configuracoes

| Atributo                                     | Tipo    | Padrão | Descrição                                    |
|----------------------------------------------|---------|--------|----------------------------------------------|
| tempo_bloqueio_inatividade_minutos           | inteiro | 5      | Tempo até bloqueio automático por inatividade |
| tempo_ocultar_segredo_segundos               | inteiro | 15     | Tempo até reocultação de campo sensível        |
| tempo_limpar_area_transferencia_segundos     | inteiro | 30     | Tempo até limpeza automática da clipboard      |

Nenhum temporizador pode ser desabilitado — todos são obrigatórios.

### Pasta

| Atributo  | Tipo          | Descrição                                                  |
|-----------|---------------|------------------------------------------------------------|
| id        | NanoID        | Identidade persistida                                      |
| nome      | string        | Nome da pasta. Único entre irmãs da mesma pasta pai.      |
| segredos  | list[Segredo] | Segredos diretamente nessa pasta, ordenados por posição    |
| pastas    | list[Pasta]   | Subpastas nessa pasta, ordenadas por posição               |

A **Pasta Geral** é a raiz da hierarquia. Não pode ser renomeada, movida ou excluída.

### Segredo

| Atributo                | Tipo               | Descrição                                                                  |
|-------------------------|--------------------|----------------------------------------------------------------------------|
| id                      | NanoID             | Identidade persistida                                                      |
| nome                    | string             | Nome do segredo. Sem restrição de unicidade.                               |
| nome_modelo_segredo     | string (opcional)  | Registro histórico do modelo usado na criação. Não é vínculo ativo.        |
| campos                  | list[CampoSegredo] | Campos do segredo, ordenados por posição. A Observação ocupa a última posição. |
| favorito                | booleano           | Indica se o segredo está favoritado                                        |
| data_criacao            | datetime           | Data/hora de criação do segredo                                            |
| data_ultima_modificacao | datetime           | Data/hora da última modificação do segredo                                 |

### CampoSegredo

| Atributo | Tipo                         | Descrição                                                    |
|----------|------------------------------|--------------------------------------------------------------|
| nome     | string                       | Nome do campo. Sem restrição de unicidade dentro do segredo. |
| tipo     | enum: texto, texto_sensivel  | Determina o comportamento de exibição                        |
| valor    | string                       | Valor do campo. String vazia = campo não preenchido.         |

O campo **Observação** é um CampoSegredo especial: tipo `texto`, nome fixo "Observação", sempre na última posição, não pode ser renomeado, movido ou excluído.

### ModeloSegredo

| Atributo | Tipo                      | Descrição                                                         |
|----------|---------------------------|-------------------------------------------------------------------|
| id       | NanoID                    | Identidade persistida                                             |
| nome     | string                    | Nome do modelo. Único entre todos os modelos do cofre.            |
| campos   | list[CampoModeloSegredo]  | Estrutura de campos do modelo, ordenados por posição              |

Modelos são exibidos em ordem alfabética — não são reordenáveis manualmente.

### CampoModeloSegredo

| Atributo | Tipo                         | Descrição                                |
|----------|------------------------------|------------------------------------------|
| nome     | string                       | Nome do campo.                           |
| tipo     | enum: texto, texto_sensivel  | Tipo do campo. Permite alteração no modelo (não no segredo). |

---

## Pastas Virtuais

Pastas virtuais são **vistas derivadas** do estado em memória. Não são persistidas no arquivo.

| Pasta Virtual | Definição                                                                                                                                                  |
|---------------|------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **Favoritos** | Conjunto de segredos com `favorito = true`, percorrido em profundidade seguindo a ordem do JSON.                                                           |
| **Lixeira**   | Conjunto de segredos pendentes de exclusão, mantidos pelo Manager durante a sessão. Não fazem parte da hierarquia persistida. Removidos permanentemente ao salvar. |

---

## Esquema Compacto

Representação concisa da estrutura persistida no payload JSON do arquivo `.abditum`.

```
Cofre:
  configuracoes:
    tempo_bloqueio_inatividade_minutos:       inteiro  (padrão: 5)
    tempo_ocultar_segredo_segundos:           inteiro  (padrão: 15)
    tempo_limpar_area_transferencia_segundos: inteiro  (padrão: 30)
  pasta_geral:            Pasta
  modelos_segredo:        list[ModeloSegredo]
  data_criacao:           datetime
  data_ultima_modificacao: datetime

Pasta:
  id:       nanoid (6 chars)
  nome:     string
  segredos: list[Segredo]
  pastas:   list[Pasta]

Segredo:
  id:                  nanoid (6 chars)
  nome:                string
  nome_modelo_segredo: string?
  campos:              list[CampoSegredo]   -- último é sempre a Observação
  favorito:            booleano
  data_criacao:        datetime
  data_ultima_modificacao: datetime

CampoSegredo:
  nome:  string
  tipo:  enum { texto | texto_sensivel }
  valor: string                            -- string vazia = não preenchido

ModeloSegredo:
  id:     nanoid (6 chars)
  nome:   string
  campos: list[CampoModeloSegredo]

CampoModeloSegredo:
  nome: string
  tipo: enum { texto | texto_sensivel }
```

---

## Invariantes

### Pertencimento único
- Um segredo pertence a exatamente uma Pasta — nunca a duas simultaneamente.
- Uma pasta pertence a exatamente uma pasta pai — nunca a duas simultaneamente.

### Hierarquia acíclica
- Ciclos não são permitidos — uma pasta nunca pode ser descendente de si mesma.
- Todas as pastas são navegáveis a partir da Pasta Geral — nenhuma pasta pode ficar desconectada.

### Pasta Geral
- Sempre presente no cofre. Não pode ser renomeada, movida ou excluída.

### Unicidade de nomes
- Duas subpastas com o mesmo nome não podem coexistir na mesma pasta pai.
- Dois modelos de segredo com o mesmo nome não podem coexistir no cofre.
- Nomes de segredos e de campos não têm restrição de unicidade.

### Observação
- Presente em todo segredo, sempre na última posição, sempre do tipo `texto`. Imutável em nome, tipo e posição.


