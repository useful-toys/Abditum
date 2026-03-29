# Guia de Build e CI/CD (Implementação)

Este guia prático mostra como configurar o pipeline de Integração Contínua (CI) para o projeto Abditum usando GitHub Actions, garantindo a qualidade, segurança e portabilidade do software a cada alteração.

## 1. Configurando o GitHub Actions

O pipeline de CI/CD do Abditum é definido por arquivos YAML no diretório `.github/workflows/` do seu repositório.

### Criar o arquivo de Workflow

1.  No seu repositório GitHub, crie a pasta `.github/workflows/`.
2.  Dentro dela, crie um novo arquivo chamado `ci.yml`.

### Conteúdo do `ci.yml`

Copie o seguinte conteúdo para o arquivo `ci.yml`. Este workflow será acionado em cada push e pull request para os branches `main` e `develop`.

