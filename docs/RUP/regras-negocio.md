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
| RN-01  | O acesso ao conteúdo do cofre requer autenticação por senha mestra.                                         |
| RN-02  | A senha mestra é irrecuperável. Seu esquecimento resulta em perda total e definitiva dos dados (Conhecimento Zero). |
| RN-03  | A criação e a alteração da senha mestra exigem confirmação por dupla digitação.                             |
| RN-04  | O cofre é autossuficiente: todas as informações necessárias para seu uso estão contidas nele, sem dependência de recursos externos. |
| RN-05  | Nenhum dado do cofre é transmitido pela rede. Toda operação é local e offline.                              |

---

## 3. Organização e Hierarquia

| ID     | Regra                                                                                                        |
|--------|--------------------------------------------------------------------------------------------------------------|
| RN-06  | Cada segredo pertence a exatamente uma pasta ou à raiz do cofre — nunca a dois locais simultaneamente.      |
| RN-07  | Cada pasta pertence a exatamente uma pasta pai ou à raiz do cofre — nunca a dois locais simultaneamente.    |
| RN-08  | Nomes de segredos, pastas e modelos não são identificadores. Nomes repetidos são permitidos; a identidade de cada elemento é independente do nome. |

---

## 4. Classificação de Dados

| ID     | Regra                                                                                                        |
|--------|--------------------------------------------------------------------------------------------------------------|
| RN-09  | Os campos de um segredo são classificados como dados comuns (texto) ou dados sensíveis (texto sensível).    |
| RN-10  | Dados sensíveis são protegidos por padrão e nunca participam de buscas.                                     |
| RN-11  | A observação de um segredo é classificada como dado não sensível.                                           |

---

## 5. Exclusão e Restauração

| ID     | Regra                                                                                                        |
|--------|--------------------------------------------------------------------------------------------------------------|
| RN-12  | A exclusão de um segredo é reversível até a próxima persistência definitiva do cofre.                       |
| RN-13  | Após a persistência definitiva, segredos marcados para exclusão são eliminados permanentemente, sem possibilidade de recuperação. |
| RN-14  | Um segredo marcado para exclusão não pode ser editado.                                                      |
| RN-15  | Ao restaurar um segredo cuja pasta de origem não exista mais, ele retorna à raiz do cofre.                  |
| RN-16  | A exclusão de uma pasta é imediata e irreversível. Seus segredos e subpastas são promovidos ao nível hierárquico superior. |

---

## 6. Modelos de Segredo

| ID     | Regra                                                                                                        |
|--------|--------------------------------------------------------------------------------------------------------------|
| RN-17  | Modelos de segredo são templates de criação. Um segredo criado a partir de um modelo não mantém vínculo retroativo com ele. |
| RN-18  | Alterações na estrutura de um modelo afetam apenas criações futuras. Segredos existentes permanecem inalterados. |

---

## 7. Importação e Conflitos

| ID     | Regra                                                                                                        |
|--------|--------------------------------------------------------------------------------------------------------------|
| RN-19  | Na importação, pastas com a mesma identidade de pastas já existentes no cofre são mescladas.                |
| RN-20  | Na importação, segredos com identidade conflitante recebem nova identidade, preservando os demais dados.    |
| RN-21  | Na importação, segredos com nome conflitante na mesma pasta recebem ajuste de nome para evitar ambiguidade. |
| RN-22  | Na importação, modelos com a mesma identidade são substituídos pelo modelo importado.                       |

---

## 8. Exportação

| ID     | Regra                                                                                                        |
|--------|--------------------------------------------------------------------------------------------------------------|
| RN-23  | A exportação do cofre gera uma cópia não protegida de todos os dados, incluindo dados sensíveis.            |

---

## 9. Rastreabilidade

| Regra  | Casos de Uso relacionados                                         |
|--------|-------------------------------------------------------------------|
| RN-01  | UC-01, UC-02                                                      |
| RN-02  | UC-01, UC-06                                                      |
| RN-03  | UC-01, UC-06                                                      |
| RN-04  | UC-01, UC-08                                                      |
| RN-05  | —                                                                 |
| RN-06  | UC-23, UC-24                                                      |
| RN-07  | UC-26, UC-28, UC-29                                               |
| RN-08  | UC-16, UC-17, UC-26, UC-31                                       |
| RN-09  | UC-14, UC-16                                                      |
| RN-10  | UC-14, UC-15                                                      |
| RN-11  | UC-15                                                             |
| RN-12  | UC-21                                                             |
| RN-13  | UC-03, UC-04, UC-21                                               |
| RN-14  | UC-21                                                             |
| RN-15  | UC-22                                                             |
| RN-16  | UC-30                                                             |
| RN-17  | UC-16, UC-33, UC-34                                               |
| RN-18  | UC-32                                                             |
| RN-19  | UC-10                                                             |
| RN-20  | UC-10                                                             |
| RN-21  | UC-10                                                             |
| RN-22  | UC-10                                                             |
| RN-23  | UC-09                                                             |