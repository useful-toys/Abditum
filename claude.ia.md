# Contexto — Sessão de Definição do Abditum

## O que foi feito

Esta sessão foi dedicada à definição do vocabulário, glossário e requisitos do Abditum — um cofre de senhas portátil com interface TUI. O trabalho seguiu uma progressão deliberada: primeiro o vocabulário, depois os conceitos, depois os requisitos.

## Abordagem adotada

O processo foi inspirado na prática de **linguagem ubíqua do DDD** — definir termos precisos e compartilhados antes de descrever funcionalidades. Não foi aplicado DDD completo, apenas sua parte mais valiosa: nomear bem as coisas antes de construí-las.

O documento de requisitos foi organizado em camadas:
- **Requisitos Funcionais** — o que a aplicação faz, com regras específicas como subitens
- **Regras Transversais** — comportamentos que cruzam múltiplos requisitos
- **Decisões de UX** — comportamentos e padrões de interface
- **Decisões Técnicas** — escolhas de implementação
- **Histórico de Decisões** — registro do raciocínio por trás das principais escolhas
- **Coisas para Repensar** — pendências em aberto

## Decisões de vocabulário

- **Campo** foi escolhido para descrever os elementos de um segredo. Alternativas descartadas: "dado" (muito amplo), "detalhe" (diminui semanticamente dados críticos como senha), "faceta" (sugere perspectiva, não composição).
- **Modelo de segredo** unificou o que antes eram "modelo pré-definido" e "modelo personalizado" — a distinção era de origem, não de comportamento.
- **Pasta Geral** substituiu o conceito técnico de "raiz" — eliminou jargão de estrutura de dados sem criar novos problemas.
- Termos evitados intencionalmente: "sistema", "caminho", "raiz", "container", "detalhes", "hierarquia".

## Decisões de design relevantes

- Todo segredo vive dentro de uma pasta — não existe segredo fora de uma pasta
- A pasta Geral é sempre presente, não pode ser renomeada nem excluída
- O tipo de um campo não pode ser alterado após a criação — converter sensível em comum exporia o conteúdo
- Segredos criados a partir de um modelo não mantêm vínculo com ele
- Campos sensíveis nunca participam da busca
- A observação é considerada dado não sensível

## Funcionalidades descartadas do escopo

- Auditoria de senhas
- TOTP (autenticação de dois fatores)

## Arquivos gerados

- `abditum-glossario.md` — visão do produto, glossário e seção de segurança
- `abditum-requisitos.md` — requisitos funcionais, regras, decisões e histórico

## Pendências em aberto

- O termo "identidade" aparece nas regras de importação de arquivo JSON mas não foi definido no glossário — é um conceito técnico interno que pode precisar de esclarecimento ou substituição por linguagem mais acessível.