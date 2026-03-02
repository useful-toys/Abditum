# Abditum

> *abditum* (latim) — oculto, guardado em lugar secreto.

Gerenciador pessoal de dados sensíveis com estrutura hierárquica em árvore, armazenado em um único arquivo criptografado (`.abt`).

## Características

- **Arquivo único portável** — sem banco de dados, sem servidor, sem instalação
- **Criptografia forte** — Argon2id + AES-256-GCM
- **Interface TUI** — roda em qualquer terminal, multiplataforma
- **Payload legível** — JSON puro após descriptografia, sem ferramentas adicionais
- **Sem telemetria** — nenhum dado trafega pela rede

## Plataformas

Windows · macOS · Linux

## Status

> Em desenvolvimento. As especificações estão em [`specs/`](specs/).

## Especificações

| Documento | Descrição |
|-----------|-----------|
| [Visão Geral](specs/overview.md) | Propósito, princípios e funcionalidades |
| [Arquitetura](specs/architecture.md) | Camadas, serviços e fluxos |
| [Modelo de Dados](specs/data-model.md) | Entidades, formato do arquivo e invariantes |
| [Segurança](specs/security-hardening.md) | Hardening por fase |
| [Backlog](specs/backlog.md) | Ideias futuras sem compromisso |

## Licença

[MIT](LICENSE)
