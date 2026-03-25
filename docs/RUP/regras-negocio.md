# Regras de Negócio — Abditum

| Item            | Detalhe                        |
|-----------------|--------------------------------|
| Projeto         | Abditum                       |
| Versão          | 1.1                           |
| Data            | 2026-03-25                     |
| Status          | Aprovado                       |

---

## 1. Introdução

Este documento formaliza as regras de negócio do Abditum — políticas, restrições e invariantes que governam o cofre de senhas independentemente de tecnologia, interface ou implementação. Uma regra é considerada "de negócio" se existiria mesmo na ausência de uma solução de software.

### 1.1 Referências
- Documento de Visão — `docs/RUP/visao.md`
- Glossário — `docs/RUP/glossario.md`
- SRS — `docs/RUP/srs.md`
- Casos de Uso — `docs/RUP/casos-de-uso.md`

---

## 2. Cofre e Acesso

| ID     | Regra                                                                                                        |
|--------|--------------------------------------------------------------------------------------------------------------|
| <a id="rn-01"></a>RN-01  | O acesso ao conteúdo do cofre requer autenticação por senha mestra.                                         |
| <a id="rn-02"></a>RN-02  | A senha mestra é irrecuperável. Seu esquecimento resulta em perda total e definitiva dos dados (Conhecimento Zero). |
| <a id="rn-03"></a>RN-03  | A criação e a alteração da senha mestra exigem confirmação por dupla digitação.                             |
| <a id="rn-04"></a>RN-04  | O cofre é autossuficiente: todas as informações necessárias para seu uso estão contidas nele, sem dependência de recursos externos. |
| <a id="rn-05"></a>RN-05  | Nenhum dado do cofre é transmitido pela rede. Toda operação é local e offline.                              |

---

## 3. Organização e Hierarquia

| ID     | Regra                                                                                                        |
|--------|--------------------------------------------------------------------------------------------------------------|
| <a id="rn-06"></a>RN-06  | Cada segredo pertence a exatamente uma pasta ou à raiz do cofre — nunca a dois locais simultaneamente.      |
| <a id="rn-07"></a>RN-07  | Cada pasta pertence a exatamente uma pasta pai ou à raiz do cofre — nunca a dois locais simultaneamente.    |
| <a id="rn-08"></a>RN-08  | Nomes de segredos, pastas e modelos não são identificadores. Nomes repetidos são permitidos; a identidade de cada elemento é independente do nome. |

---

## 4. Classificação de Dados

| ID     | Regra                                                                                                        |
|--------|--------------------------------------------------------------------------------------------------------------|
| <a id="rn-09"></a>RN-09  | Os campos de um segredo são classificados como dados comuns (texto) ou dados sensíveis (texto sensível).    |
| <a id="rn-10"></a>RN-10  | Dados sensíveis são protegidos por padrão e nunca participam de buscas.                                     |
| <a id="rn-11"></a>RN-11  | A observação de um segredo é classificada como dado não sensível.                                           |

---

## 5. Exclusão e Restauração

| ID     | Regra                                                                                                        |
|--------|--------------------------------------------------------------------------------------------------------------|
| <a id="rn-12"></a>RN-12  | A exclusão de um segredo é reversível até a próxima persistência definitiva do cofre.                       |
| <a id="rn-13"></a>RN-13  | Após a persistência definitiva, segredos marcados para exclusão são eliminados permanentemente, sem possibilidade de recuperação. |
| <a id="rn-14"></a>RN-14  | Um segredo marcado para exclusão não pode ser editado.                                                      |
| <a id="rn-15"></a>RN-15  | Ao restaurar um segredo cuja pasta de origem não exista mais, ele retorna à raiz do cofre.                  |
| <a id="rn-16"></a>RN-16  | A exclusão de uma pasta é imediata e irreversível. Seus segredos e subpastas são promovidos ao nível hierárquico superior. |

---

## 6. Modelos de Segredo

| ID     | Regra                                                                                                        |
|--------|--------------------------------------------------------------------------------------------------------------|
| <a id="rn-17"></a>RN-17  | Modelos de segredo são templates de criação. Um segredo criado a partir de um modelo não mantém vínculo retroativo com ele. |
| <a id="rn-18"></a>RN-18  | Alterações na estrutura de um modelo afetam apenas criações futuras. Segredos existentes permanecem inalterados. |

---

## 7. Importação e Conflitos

| ID     | Regra                                                                                                        |
|--------|--------------------------------------------------------------------------------------------------------------|
| <a id="rn-19"></a>RN-19  | Na importação, pastas com a mesma identidade de pastas já existentes no cofre são mescladas.                |
| <a id="rn-20"></a>RN-20  | Na importação, segredos com identidade conflitante recebem nova identidade, preservando os demais dados.    |
| <a id="rn-21"></a>RN-21  | Na importação, segredos com nome conflitante na mesma pasta recebem ajuste de nome para evitar ambiguidade. |
| <a id="rn-22"></a>RN-22  | Na importação, modelos com a mesma identidade são substituídos pelo modelo importado.                       |

---

## 8. Exportação

| ID     | Regra                                                                                                        |
|--------|--------------------------------------------------------------------------------------------------------------|
| <a id="rn-23"></a>RN-23  | A exportação do cofre gera uma cópia não protegida de todos os dados, incluindo dados sensíveis.            |

---

## 9. Rastreabilidade

| Regra  | Casos de Uso relacionados                                         |
|--------|-------------------------------------------------------------------|
| [RN-01](#rn-01)  | [UC-01](casos-de-uso.md#uc-01), [UC-02](casos-de-uso.md#uc-02)                                                      |
| [RN-02](#rn-02)  | [UC-01](casos-de-uso.md#uc-01), [UC-06](casos-de-uso.md#uc-06)                                                      |
| [RN-03](#rn-03)  | [UC-01](casos-de-uso.md#uc-01), [UC-06](casos-de-uso.md#uc-06)                                                      |
| [RN-04](#rn-04)  | [UC-01](casos-de-uso.md#uc-01), [UC-08](casos-de-uso.md#uc-08)                                                      |
| [RN-05](#rn-05)  | —                                                                 |
| [RN-06](#rn-06)  | [UC-23](casos-de-uso.md#uc-23), [UC-24](casos-de-uso.md#uc-24)                                                      |
| [RN-07](#rn-07)  | [UC-26](casos-de-uso.md#uc-26), [UC-28](casos-de-uso.md#uc-28), [UC-29](casos-de-uso.md#uc-29)                                               |
| [RN-08](#rn-08)  | [UC-16](casos-de-uso.md#uc-16), [UC-17](casos-de-uso.md#uc-17), [UC-26](casos-de-uso.md#uc-26), [UC-31](casos-de-uso.md#uc-31)                                       |
| [RN-09](#rn-09)  | [UC-14](casos-de-uso.md#uc-14), [UC-16](casos-de-uso.md#uc-16)                                                      |
| [RN-10](#rn-10)  | [UC-14](casos-de-uso.md#uc-14), [UC-15](casos-de-uso.md#uc-15)                                                      |
| [RN-11](#rn-11)  | [UC-15](casos-de-uso.md#uc-15)                                                             |
| [RN-12](#rn-12)  | [UC-21](casos-de-uso.md#uc-21)                                                             |
| [RN-13](#rn-13)  | [UC-03](casos-de-uso.md#uc-03), [UC-04](casos-de-uso.md#uc-04), [UC-21](casos-de-uso.md#uc-21)                                               |
| [RN-14](#rn-14)  | [UC-21](casos-de-uso.md#uc-21)                                                             |
| [RN-15](#rn-15)  | [UC-22](casos-de-uso.md#uc-22)                                                             |
| [RN-16](#rn-16)  | [UC-30](casos-de-uso.md#uc-30)                                                             |
| [RN-17](#rn-17)  | [UC-16](casos-de-uso.md#uc-16), [UC-33](casos-de-uso.md#uc-33), [UC-34](casos-de-uso.md#uc-34)                                               |
| [RN-18](#rn-18)  | [UC-32](casos-de-uso.md#uc-32)                                                             |
| [RN-19](#rn-19)  | [UC-10](casos-de-uso.md#uc-10)                                                             |
| [RN-20](#rn-20)  | [UC-10](casos-de-uso.md#uc-10)                                                             |
| [RN-21](#rn-21)  | [UC-10](casos-de-uso.md#uc-10)                                                             |
| [RN-22](#rn-22)  | [UC-10](casos-de-uso.md#uc-10)                                                             |
| [RN-23](#rn-23)  | [UC-09](casos-de-uso.md#uc-09)                                                             |