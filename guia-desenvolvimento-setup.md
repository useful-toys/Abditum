# Guia de Desenvolvimento e Setup

Este guia detalha os passos para configurar o ambiente de desenvolvimento, compilar e executar o projeto Abditum localmente.

## Pré-requisitos

Certifique-se de ter o Go instalado em sua máquina. O projeto requer a **versão 1.26 ou superior**.

Para verificar sua versão do Go, execute:

```bash
go version
```

Se precisar instalar ou atualizar o Go, visite [go.dev/doc/install](https://go.dev/doc/install).

**Para usuários Windows:** Você pode instalar o Go facilmente usando o Scoop, um gerenciador de pacotes de linha de comando. Se você já tem o Scoop, execute:

```powershell
scoop install go
```
Se não tiver o Scoop, você pode instalá-lo seguindo as instruções em [scoop.sh](https://scoop.sh).

## Configuração do Ambiente

1.  **Clone o Repositório**:
    Comece clonando o repositório do Abditum para sua máquina local:

    ```bash
    git clone https://github.com/seu-usuario/abditum.git
    cd abditum
    ```

2.  **Instale as Dependências**:
    O Go Modules gerencia as dependências do projeto. Instale-as executando:

    ```bash
    go mod tidy
    ```
    Este comando garantirá que todas as dependências listadas no `go.mod` sejam baixadas e estejam prontas para uso.

## Executando a Aplicação Localmente

Para executar a aplicação diretamente do código-fonte (útil para desenvolvimento e depuração):

```bash
go run cmd/abditum/main.go
```

Isso compilará e executará o ponto de entrada principal do Abditum.

Você também pode passar um arquivo de cofre existente como argumento:

```bash
go run cmd/abditum/main.go meu_cofre.abditum
```

## Executando os Testes

O Abditum possui uma suíte de testes abrangente (unitários, integração e testes de TUI). É crucial executá-los para garantir que suas alterações não introduzam regressões.

Para executar todos os testes do projeto:

```bash
go test ./...
```

Este comando irá percorrer todos os pacotes do projeto e executar os testes neles.

### Testes de TUI com Golden Files

Os testes de TUI usam "golden files" para snapshots visuais. Se você fizer alterações na interface do usuário (TUI) que mudem a saída esperada, os testes podem falhar.

Para atualizar os golden files após uma mudança intencional na TUI:

```bash
GO_TEST_UPDATE_GOLDEN=true go test ./...
```

**Importante**: Revise sempre as mudanças nos golden files antes de fazer commit para garantir que as alterações visuais são as esperadas.

## Linting

Para manter a qualidade do código e seguir as convenções, o projeto utiliza ferramentas de linting. Recomenda-se executar o linter antes de cada commit.

Para executar o linter (geralmente `golangci-lint` ou `go vet` + outras ferramentas):

```bash
# Exemplo com go vet
go vet ./...

# Se for usado golangci-lint, instale-o primeiro:
# go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
# E então execute:
# golangci-lint run
```

Consulte a seção "CI/CD (Implementação)" para ver as ferramentas de linting exatas configuradas no pipeline.

## Cross-Compilação (Gerando Binários para Diferentes OSes)

Para gerar binários executáveis autônomos para diferentes sistemas operacionais e arquiteturas (como o produto final distribuído):

### Linux (64-bit)

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/abditum-linux-amd64 cmd/abditum/main.go
```

### macOS (64-bit Intel)

```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o bin/abditum-darwin-amd64 cmd/abditum/main.go
```

### macOS (64-bit Apple Silicon)

```bash
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o bin/abditum-darwin-arm64 cmd/abditum/main.go
```

### Windows (64-bit)

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/abditum-windows-amd64.exe cmd/abditum/main.go
```

Os binários compilados serão gerados na pasta `bin/`.

**Observação sobre `CGO_ENABLED=0`**: Esta configuração é crucial para garantir que o Go compile um binário *estaticamente linkado*, sem dependências de bibliotecas C do sistema operacional. Isso é fundamental para a portabilidade e segurança do Abditum, conforme detalhado em `arquitetura.md`.
