# Arquitetura — Abditum

## 1. Tecnologias

| Área              | Escolha              | Observação                                                                          |
|-------------------|----------------------|-------------------------------------------------------------------------------------|
| Linguagem         | Go                   | Compilado como binário único executável, sem dependências de runtime externas       |
| Interface TUI     | Bubble Tea (Charm)   | Modelo Elm de atualização de estado; componentes reutilizáveis                      |
| Testes de TUI     | teatest/v2           | Simulação de terminal para testes de interface                                      |

---

## 2. Estrutura de Pacotes

```
cmd/
  abditum/        -- ponto de entrada da aplicação (main)
internal/
  vault/          -- domínio e lógica de negócio: Manager, entidades, regras de negócio
  crypto/         -- derivação de chave (Argon2id) e criptografia/descriptografia (AES-256-GCM)
  storage/        -- leitura e escrita do arquivo .abditum: formato binário e salvamento atômico
  tui/            -- interface TUI: modelos Bubble Tea, telas, componentes e navegação
```

---

## 3. Padrão Manager

As entidades do domínio (`Cofre`, `Pasta`, `Segredo`, `ModeloSegredo`) são navegáveis somente leitura externamente. Toda mutação é realizada exclusivamente por métodos explícitos do `Manager` (`internal/vault`).

O Manager é o contrato central de manipulação do domínio: centraliza as regras de negócio, garante consistência das estruturas de dados e impede manipulações diretas potencialmente inseguras. A TUI interage com o domínio exclusivamente via Manager.

---

## 4. Integração Contínua (CI)

CI é obrigatório. Todo push deve executar automaticamente: build, lint e a suíte completa de testes. Mudanças que quebrem o build ou qualquer teste não são aceitas.

---

## 5. Convenções de Código

### Comentários

O projeto adota uma política generosa de comentários. O código deve ser acessível a leitores com menos familiaridade com Go e com as bibliotecas especializadas de criptografia e TUI. Comentários devem explicar o *porquê* das decisões de implementação — não apenas descrever o que o código faz.

Packages de criptografia e TUI merecem atenção especial: os conceitos subjacentes (derivação de chave, AEAD, modelo Elm) devem ser explicados no código com clareza suficiente para um leitor não familiarizado.

### Logs e privacidade

A aplicação não deve emitir nenhum log (stdout/stderr) que contenha caminhos de arquivos de cofre, nomes de segredos, nomes de campos ou valores de campos. Mensagens de erro exibidas ao usuário são genéricas por design — o mesmo princípio se aplica a qualquer saída de diagnóstico interna.

### Build

- **Build estático**: `CGO_ENABLED=0`. Binário auto-contido sem dependência de libC, eliminando superfície de ataque de bibliotecas C do sistema e garantindo portabilidade.
- **Isolamento de rede**: a aplicação nunca faz chamadas de rede. Nenhum import de packages de rede (`net`, `net/http`) é permitido. Violação dessa regra é defeito de build.
- **Entropia**: toda geração de valores aleatórios (salt, nonce, NanoID) deve usar exclusivamente `crypto/rand`. O uso de `math/rand` é proibido.
- **Criptografia via stdlib**: as operações criptográficas usam exclusivamente packages da standard library de Go (`crypto/aes`, `crypto/cipher`) e `golang.org/x/crypto/argon2`. Libs de criptografia de terceiros não são permitidas.
- **Dependências mínimas**: cada dependência externa é superficíe de ataque. A adição de qualquer dependência deve ser justificada e ponderada. Preferência absoluta por packages da standard library.

---

## 6. Segurança em Sessão

### Memória protegida (mlock / VirtualLock)

A aplicação deve tentar alocar a senha mestra e a chave AES derivada em memória bloqueada (`mlock` no Linux/macOS, `VirtualLock` no Windows), impedindo que o SO faça swap dessas regiões para disco. Se a chamada falhar (permissões insuficientes, limite de quota), a aplicação opera normalmente sem essa camada — a ausência de suporte não é erro fatal.

### Zeragem de dados sensíveis ao bloquear / sair

Ao bloquear o cofre ou encerrar a aplicação, a aplicação deve sobrescrever com zeros todos os buffers sob seu controle que contenham dados sensíveis: senha mestra, chave AES derivada e quaisquer valores de campos sensíveis retidos em memória. Em Go, o runtime pode mover slices durante GC, o que limita a garantia; ainda assim, a zeragem explícita reduz a janela de exposição e é implementada como melhor esforço.

### Clipboard

O valor copiado para a área de transferência deve ser apagado:
- automaticamente após o tempo configurado no cofre (padrão: 30 segundos); e
- imediatamente ao bloquear o cofre ou encerrar a aplicação.

A limpeza depende do suporte do sistema operacional. Se não for possível, a aplicação opera normalmente sem essa garantia.

### Clear screen

Ao bloquear o cofre ou encerrar a aplicação, a aplicação deve limpar o terminal (clear screen) antes de devolver o controle ao shell, evitando que dados visíveis na TUI permaneçam no buffer de rolagem.

### Colisão de nonce

Com nonce aleatório de 96 bits e chave fixa (mesmo salt + mesma senha), o birthday bound para colisão é ~2⁴⁸ cifragens. Para um cofre pessoal salvo milhares de vezes, o risco é desprezível. Essa segurança depende do salt ser substituído a cada troca de senha — se a relação salt/chave mudar em versões futuras, a análise deve ser refeita.

---

## 7. Concorrência e Acesso Externo

Não é usado arquivo de lock. Modificações externas ao arquivo do cofre são detectadas comparando timestamp e tamanho no momento do salvamento. Se divergência for detectada, o usuário é avisado e tem as opções: sobrescrever / salvar como novo arquivo / cancelar. Esta abordagem preserva portabilidade e privacidade total (sem rastros no sistema de arquivos além do próprio cofre e seus artefatos previstos).

---

## 8. Estratégia de Testes

| Categoria                | Escopo                                                                                                    |
|--------------------------|-----------------------------------------------------------------------------------------------------------|
| Serviço de criptografia  | Casos de sucesso e falha de criptografia e descriptografia: Argon2id + AES-256-GCM                       |
| Serviço de armazenamento | Casos de sucesso e falha de salvamento e carregamento do arquivo `.abditum`                               |
| Unitários white-box      | Navegação e transições de estado do domínio                                                               |
| Golden files             | Snapshot visual de cada tela em terminal 80×24; detectam regressões visuais automaticamente               |
| Testes de comandos       | Cada tela e fluxo de usuário, via teatest/v2                                                              |
| Integração               | Fluxo completo de ponta a ponta: criar cofre, criar segredo, editar segredo e demais operações principais |
