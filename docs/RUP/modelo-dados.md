# Modelo de Dados e Estrutura de Payload — Abditum

| Item              | Detalhe                                   |
|-------------------|-------------------------------------------|
| **Projeto**       | Abditum — Cofre de Senhas Portátil        |
| **Versão**        | 1.0                                       |
| **Data**          | 2026-03-25                                |

---

## 1. Visão Geral

O arquivo do cofre (`.abditum`) é um stream binário composto por duas partes: um cabeçalho de criptografia em formato fixo e um payload criptografado contendo a estrutura do cofre em JSON (UTF-8). Este documento especifica ambas as partes em detalhe.

---

## 2. Formato do Arquivo (.abditum)

### 2.1. Visão Estrutural

```
┌──────────────────────────────────────────────────────┐
│                  ARQUIVO .abditum                    │
├──────────────────────────────────────────────────────┤
│  CABEÇALHO DE CRIPTOGRAFIA (fixo, não criptografado) │
│  ┌──────────────┬──────────────────────────────────┐ │
│  │ magic        │ 4 bytes — ASCII "ABDT"           │ │
│  ├──────────────┼──────────────────────────────────┤ │
│  │ versão_formato│ inteiro — versão do formato     │ │
│  ├──────────────┼──────────────────────────────────┤ │
│  │ salt         │ bytes — salt para Argon2id       │ │
│  ├──────────────┼──────────────────────────────────┤ │
│  │ nonce        │ bytes — IV para AES-256-GCM      │ │
│  └──────────────┴──────────────────────────────────┘ │
├──────────────────────────────────────────────────────┤
│  PAYLOAD CRIPTOGRAFADO                               │
│  ┌──────────────────────────────────────────────────┐ │
│  │ AES-256-GCM( JSON UTF-8 do Cofre )              │ │
│  │ + Authentication Tag                             │ │
│  └──────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────┘
```

### 2.2. Cabeçalho de Criptografia

| Campo            | Tamanho         | Descrição                                                                  |
|------------------|-----------------|----------------------------------------------------------------------------|
| `magic`          | 4 bytes         | Assinatura fixa ASCII `ABDT` — rejeição precoce de arquivos inválidos      |
| `versão_formato` | inteiro (4 B)   | Versão do formato do arquivo; seleção de perfil Argon2id histórico         |
| `salt`           | 16 bytes        | Salt aleatório gerado na criação, armazenado no cabeçalho                  |
| `nonce`          | 12 bytes        | Vetor de inicialização (IV) para AES-256-GCM, regenerado a cada salvamento |

**Propriedades do cabeçalho:**

- Todo o cabeçalho (`magic` + `versão_formato` + `salt` + `nonce`) entra como **AAD** (Additional Authenticated Data) do AES-256-GCM, garantindo integridade sem checksum adicional.
- A assinatura `magic` permite diferenciar imediatamente erro de tipo de arquivo de erro de senha ou corrupção.
- O `salt` é único por cofre (gerado aleatoriamente na criação).
- O `nonce` é único por salvamento (regenerado a cada operação de escrita).

### 2.3. Payload Criptografado

O restante do stream após o cabeçalho contém:

- A estrutura `Cofre` serializada em **JSON UTF-8**.
- Criptografada com **AES-256-GCM** usando a chave derivada por Argon2id.
- Inclui o **Authentication Tag** do GCM, que autentica tanto o payload quanto o cabeçalho (via AAD).

### 2.4. Fluxo de Derivação de Chave

```
senha_mestra (string UTF-8)
       │
       ▼
┌──────────────────────┐
│      Argon2id         │
│  • memory: 256 MiB    │
│  • iterations: ≥ 3    │
│  • parallelism: ≤ 4   │
│  • salt: do cabeçalho │
│  • output: 32 bytes   │
└──────────────────────┘
       │
       ▼
   chave AES-256 (32 bytes)
       │
       ▼
┌──────────────────────┐
│    AES-256-GCM        │
│  • nonce: do cabeçalho│
│  • AAD: cabeçalho     │
│  • plaintext: JSON    │
└──────────────────────┘
       │
       ▼
   ciphertext + auth tag
```

### 2.5. Parâmetros Argon2id (v1)

| Parâmetro       | Valor v1     | Limites                                  |
|-----------------|--------------|------------------------------------------|
| Memória         | 256 MiB      | Piso: 128 MiB · Teto referência: 512 MiB|
| Iterações       | ≥ 3          | Fixo na versão da aplicação              |
| Paralelismo     | ≤ 4 threads  | Limitado pela máquina                    |
| Output          | 32 bytes     | Chave AES-256                            |
| Salt            | 16 bytes     | Aleatório, único por cofre               |

- Parâmetros são fixos e hard-coded por versão da aplicação.
- Evolução ocorre apenas com nova versão, acompanhada de testes de regressão.
- Ao abrir cofre, o perfil Argon2id é selecionado pela `versão_formato` do cabeçalho.

---

## 3. Estrutura do Payload JSON

### 3.1. Esquema Completo

```json
{
  "configuracoes": {
    "tempo_bloqueio_inatividade_minutos": 2,
    "tempo_ocultar_segredo_segundos": 15,
    "tempo_limpar_area_transferencia_segundos": 30
  },
  "segredos": [
    {
      "id": "aB3xY7",
      "nome": "Exemplo na raiz",
      "nome_modelo_segredo": "Login",
      "campos": [
        {
          "nome": "URL",
          "tipo": "texto",
          "valor": "https://exemplo.com"
        },
        {
          "nome": "Username",
          "tipo": "texto",
          "valor": "usuario"
        },
        {
          "nome": "Password",
          "tipo": "texto_sensivel",
          "valor": "s3nh@F0rt3"
        }
      ],
      "favorito": true,
      "observacao": "Conta principal do serviço X",
      "data_criacao": "2026-01-15T10:30:00Z",
      "data_ultima_modificacao": "2026-03-20T14:45:00Z"
    }
  ],
  "pastas": [
    {
      "id": "kL9mN2",
      "nome": "Sites",
      "segredos": [],
      "pastas": [
        {
          "id": "pQ4rS8",
          "nome": "Redes Sociais",
          "segredos": [
            {
              "id": "tU1vW5",
              "nome": "Rede Social A",
              "nome_modelo_segredo": null,
              "campos": [
                {
                  "nome": "Email",
                  "tipo": "texto",
                  "valor": "user@email.com"
                },
                {
                  "nome": "Senha",
                  "tipo": "texto_sensivel",
                  "valor": "outr@S3nha"
                }
              ],
              "favorito": false,
              "observacao": "",
              "data_criacao": "2026-02-10T08:00:00Z",
              "data_ultima_modificacao": "2026-02-10T08:00:00Z"
            }
          ],
          "pastas": []
        }
      ]
    },
    {
      "id": "xZ6aB0",
      "nome": "Financeiro",
      "segredos": [],
      "pastas": []
    },
    {
      "id": "cD3eF1",
      "nome": "Serviços",
      "segredos": [],
      "pastas": []
    }
  ],
  "modelos_segredo": [
    {
      "id": "gH8iJ4",
      "nome": "Login",
      "campos": [
        { "nome": "URL", "tipo": "texto" },
        { "nome": "Username", "tipo": "texto" },
        { "nome": "Password", "tipo": "texto_sensivel" }
      ]
    },
    {
      "id": "oP2qR6",
      "nome": "Cartão de Crédito",
      "campos": [
        { "nome": "Número do Cartão", "tipo": "texto_sensivel" },
        { "nome": "Nome no Cartão", "tipo": "texto" },
        { "nome": "Data de Validade", "tipo": "texto" },
        { "nome": "CVV", "tipo": "texto_sensivel" }
      ]
    },
    {
      "id": "sT9uV3",
      "nome": "API Key",
      "campos": [
        { "nome": "Nome da API", "tipo": "texto" },
        { "nome": "Chave de API", "tipo": "texto_sensivel" }
      ]
    }
  ],
  "data_criacao": "2026-01-15T10:00:00Z",
  "data_ultima_modificacao": "2026-03-20T14:45:00Z"
}
```

### 3.2. Tipagem dos Campos

| Campo JSON                       | Tipo JSON  | Formato / Restrição                        |
|----------------------------------|------------|---------------------------------------------|
| `id`                             | string     | NanoID, 6 caracteres alfanuméricos [a-zA-Z0-9] |
| `nome`                           | string     | Texto livre, não vazio                      |
| `nome_modelo_segredo`            | string¹    | Texto livre ou `null`                       |
| `tipo`                           | string     | Enum: `"texto"` ou `"texto_sensivel"`       |
| `valor`                          | string¹    | Texto livre; `""` = campo não preenchido    |
| `favorito`                       | boolean    | `true` ou `false`                           |
| `observacao`                     | string¹    | Texto livre ou `""`                         |
| `data_criacao`                   | string     | ISO 8601 UTC (`YYYY-MM-DDTHH:MM:SSZ`)      |
| `data_ultima_modificacao`        | string     | ISO 8601 UTC (`YYYY-MM-DDTHH:MM:SSZ`)      |
| `tempo_*`                        | number     | Inteiro positivo                            |
| `segredos`, `pastas`, `campos`   | array      | Pode ser vazio `[]`                         |

¹ Campos opcionais podem ser `null` ou string vazia conforme o contexto.

### 3.3. Regras de Ordenação

A ordem dos elementos nos arrays JSON reflete diretamente a ordem de exibição na interface:

- `cofre.segredos[]` — segredos na raiz, na ordem definida pelo usuário
- `cofre.pastas[]` — pastas na raiz, na ordem definida pelo usuário
- `pasta.segredos[]` — segredos dentro de uma pasta, na ordem definida
- `pasta.pastas[]` — subpastas, na ordem definida
- `segredo.campos[]` — campos do segredo, na ordem definida
- `modelo.campos[]` — campos do modelo, na ordem definida

> Não existe campo `posicao` ou `ordem` explícito; a posição no array é a fonte da verdade.

---

## 4. Operações de Persistência

### 4.1. Salvamento Atômico (caminho atual)

```
Estado: Cofre Modificado
       │
       ▼
  Serializar domínio → JSON UTF-8
       │
       ▼
  Gerar novo nonce
       │
       ▼
  Criptografar (AES-256-GCM, chave derivada, nonce, AAD=cabeçalho)
       │
       ▼
  Gravar cabeçalho + ciphertext → cofre.abditum.tmp
       │
       ▼
  ┌─ Existe .bak anterior?
  │   Sim → Renomear .bak → .bak2
  │
  ▼
  Copiar cofre.abditum → cofre.abditum.bak
       │
       ▼
  Renomear cofre.abditum.tmp → cofre.abditum
       │
       ▼
  ┌─ Sucesso?
  │   Sim → Apagar .bak2 (se existir) → Estado: Cofre Salvo
  │   Não → Restaurar .bak2 → .bak (se possível) → Exibir erro + info backup
```

### 4.2. Criação de Novo Cofre / Salvar Como

```
  Gravar cofre diretamente no caminho final (sem .tmp)
       │
       ▼
  ┌─ Existe arquivo no destino?
  │   Sim → Confirmar sobrescrita
  │         → Rotação de backup (.bak / .bak2)
  │         → Gravar novo cofre
  │
  Estado: Cofre Salvo
```

> Criação e Salvar Como não utilizam `.tmp` — não se trata de salvamento incremental sobre cofre já aberto.

### 4.3. Rotação de Backups

| Arquivo              | Descrição                                                  |
|----------------------|-------------------------------------------------------------|
| `cofre.abditum`      | Arquivo principal do cofre (ativo)                          |
| `cofre.abditum.tmp`  | Arquivo temporário durante salvamento atômico               |
| `cofre.abditum.bak`  | Backup da versão anterior (mantido após salvamento)         |
| `cofre.abditum.bak2` | Backup temporário do `.bak` anterior (removido após sucesso)|

---

## 5. Compatibilidade e Migração de Formato

### 5.1. Regras de Compatibilidade

| Regra                                  | Descrição                                                                                     |
|----------------------------------------|-----------------------------------------------------------------------------------------------|
| Leitura retroativa                     | A aplicação v*N* lê todos os formatos de payload de v0 até v*N*                              |
| Escrita na versão corrente             | Salvamento sempre grava no formato mais recente, atualizando `versão_formato` no cabeçalho    |
| Perfil Argon2id por versão             | Ao abrir cofre, o perfil é selecionado exclusivamente pela `versão_formato` do cabeçalho     |
| Versão futura desconhecida             | Arquivos com `versão_formato` > versão suportada falham com erro claro de incompatibilidade   |
| Migração em memória                    | Payload de formato histórico é migrado em memória para o modelo corrente do domínio           |

### 5.2. Tabela de Perfis Argon2id

| versão_formato | Memória  | Iterações | Paralelismo | Status        |
|----------------|----------|-----------|-------------|---------------|
| 1              | 256 MiB  | 3         | ≤ 4         | Atual (v1)    |

> Novas versões adicionam linhas; versões históricas permanecem para retrocompatibilidade.

---

## 6. Formato de Exportação/Importação (JSON Plain Text)

O formato de exportação serializa o domínio em memória como JSON não criptografado, com a mesma estrutura do payload (seção 3.1), sem cabeçalho binário. Regras de importação:

| Conflito                       | Tratamento                                                                        |
|--------------------------------|------------------------------------------------------------------------------------|
| Pasta por identidade           | Hierarquia mesclada silenciosamente                                                |
| Segredo por identidade         | Nova identidade para o segredo importado; demais dados preservados                 |
| Segredo por nome (mesma pasta) | Nome ajustado com sufixo numérico incremental — ex: "Segredo (1)"                 |
| Modelo por identidade          | Modelo importado substitui o existente silenciosamente                              |

---

## 7. Codificação e Restrições

| Aspecto                  | Especificação                                                           |
|--------------------------|-------------------------------------------------------------------------|
| Codificação do payload   | UTF-8 explícito                                                         |
| Extensão do arquivo      | `.abditum`                                                              |
| Portabilidade            | Sem arquivos externos; configurações e modelos dentro do cofre          |
| Privacidade              | Nenhum log (stdout/stderr) com caminhos, nomes de segredos ou valores   |
| Unicidade de nomes       | Não imposta; nomes repetidos são permitidos                             |
