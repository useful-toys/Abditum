# Arquitetura — Abditum

## 1. Tecnologias

| Área              | Escolha              | Observação                                                                          |
|-------------------|----------------------|-------------------------------------------------------------------------------------|
| Linguagem         | Go 1.26+             | Compilado como binário único executável, sem dependências de runtime externas       |
| Interface TUI     | Bubble Tea v2 (`charm.land/bubbletea/v2`) | Modelo Elm de atualização de estado. **`View()` retorna `tea.View` struct, não `string`** — incompatível com exemplos v1. Import path mudou em v2. |
| Componentes TUI   | Bubbles v2 (`charm.land/bubbles/v2`) | Componentes prontos (inputs, listas, etc). Deve ser usado junto com Bubble Tea v2 — versões v1 são incompatíveis. |
| Estilo TUI        | Lip Gloss v2 (`charm.land/lipgloss/v2`) | Estilos e layout. Mesma mudança de import path da v2. |
| Testes de TUI     | teatest/v2 (`github.com/charmbracelet/x/exp/teatest/v2`) | Simulação de terminal. Usar `WithInitialTermSize(80, 24)` para golden files estáveis entre máquinas. |

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

### Representação de dados sensíveis em memória

Três categorias de dados sensíveis são representadas como `[]byte` em memória, não como `string`, para permitir zeragem explícita:

| Dado | Tipo em memória | Observação |
|---|---|---|
| Senha mestra | `[]byte` | Desde a leitura do terminal até a chamada Argon2id. Nunca convertida para `string`. |
| Chave AES-256 derivada | `[]byte` | É o que Argon2id retorna e o que AES-GCM consome. Zerada após o bloqueio/encerramento. |
| `CampoSegredo.valor` | `[]byte` | Todos os campos, independente do tipo. Zerada ao bloquear/encerrar. |

`string` em Go é imutável — uma vez criada, a memória não pode ser sobrescrita pelo programa. `[]byte` é mutável e permite zeragem com `copy(b, make([]byte, len(b)))`.

**Serialização JSON de `CampoSegredo.valor`**: `[]byte` seria serializado automaticamente como Base64 pelo `encoding/json`, quebrando compatibilidade e legibilidade. Portanto, `CampoSegredo` **deve** implementar `MarshalJSON`/`UnmarshalJSON` customizados que serializam `valor` como string UTF-8 e desserializam de volta para `[]byte`. O JSON gravado em disco é idêntico ao que seria com `string` — `[]byte` é exclusivamente um detalhe de representação em memória.

---

## 6. Segurança em Sessão

### Memória protegida (mlock / VirtualLock)

A aplicação deve tentar alocar a senha mestra e a chave AES derivada em memória bloqueada (`mlock` no Linux/macOS, `VirtualLock` no Windows), impedindo que o SO faça swap dessas regiões para disco. Se a chamada falhar (permissões insuficientes, limite de quota), a aplicação opera normalmente sem essa camada — a ausência de suporte não é erro fatal.

### Zeragem de dados sensíveis ao bloquear / sair

Ao bloquear o cofre ou encerrar a aplicação, a aplicação deve sobrescrever com zeros todos os buffers `[]byte` sob seu controle: senha mestra, chave AES derivada e valores de todos os campos (`CampoSegredo.valor`). A representação de todos esses dados como `[]byte` (e nunca como `string`) é pré-requisito para que a zeragem seja possível.

**Limitação fundamental do GC do Go:** zeragem é melhor esforço, não garantia. Quando um `[]byte` é passado para uma função ou sofre `append` que causa realocação, o runtime pode criar uma nova backing array no heap — o buffer original torna-se um órfão inacessível que não pode ser zerado. `mlock`/`VirtualLock` protege apenas as páginas da alocação viva, não cópias históricas.

Para minimizar o problema:
- **Pré-alocar com tamanho exato** os buffers de senha mestra e chave AES — nunca usar `append` neles após a alocação inicial.
- **Nunca passar buffers sensíveis como `interface{}`/`any`** — interface boxing copia os dados para uma alocação nova e inacessível.
- **Nunca interpolar valores sensíveis** em `fmt.Sprintf`, `strings.Builder` ou qualquer construção de string — todas produzem heap strings impossíveis de zerar.
- Repetir essas restrições nos comentários do código dos pacotes `internal/crypto` e `internal/vault`.

### Clipboard

O valor copiado para a área de transferência deve ser apagado:
- automaticamente após o tempo configurado no cofre (padrão: 30 segundos); e
- imediatamente ao bloquear o cofre ou encerrar a aplicação.

Regras de implementação:
- A limpeza ao bloquear/encerrar **deve ser síncrona** — nunca depender de goroutine com `time.AfterFunc` que pode ser interrompida antes de executar quando `os.Exit` é chamado.
- No Linux, `xclip`/`xsel` escrevem apenas no clipboard X11 — em sessões Wayland, o clipboard Wayland não é afetado. A biblioteca escolhida deve suportar ambos os protocolos.
- Em ambientes sem clipboard (SSH, headless), a operação deve falhar graciosamente — sem crash, sem erro fatal; apenas prosseguir sem a garantia.
- A limpeza depende do suporte do sistema operacional. Se não for possível, a aplicação opera normalmente sem essa garantia.

### Clear screen

Ao bloquear o cofre ou encerrar a aplicação, a aplicação deve limpar o terminal antes de devolver o controle ao shell, evitando que dados visíveis na TUI permaneçam no buffer de rolagem.

**Atenção:** a sequência ANSI `\033[2J` limpa apenas a área visível — o scrollback buffer permanece intacto e o usuário pode rolar para cima e ver os segredos exibidos anteriormente. A sequência correta é `\033[3J\033[2J\033[H`: `3J` limpa o scrollback, `2J` limpa a tela visível, `H` move o cursor para o topo. Funciona na maioria dos emuladores modernos (xterm, iTerm2, Windows Terminal, GNOME Terminal, Alacritty).

O Bubble Tea deve ser encerrado corretamente (via seu próprio pipeline de renderização) antes de emitir a sequência de limpeza, para que o framework não sobrescreva o clear com um último render.

A limpeza de scrollback é melhor esforço — não é garantida em todos os terminais. O objetivo é reduzir a janela de exposição nos terminais modernos mais comuns.

### Colisão de nonce

Com nonce aleatório de 96 bits e chave fixa (mesmo salt + mesma senha), o birthday bound para colisão é ~2⁴⁸ cifragens. Para um cofre pessoal salvo milhares de vezes, o risco é desprezível. Essa segurança depende do salt ser substituído a cada troca de senha — se a relação salt/chave mudar em versões futuras, a análise deve ser refeita.

Regras de implementação para manter essa garantia:
- O nonce **deve ser gerado imediatamente antes de cada chamada `gcm.Seal()`** — nunca reutilizar o mesmo slice entre chamadas.
- O nonce **nunca deve ser derivado** de contador, timestamp ou qualquer fonte determinística. Apenas `crypto/rand`.
- O nonce lido do arquivo na abertura é usado exclusivamente para descriptografia — jamais reutilizado em uma nova cifragem.
- A suíte de testes deve incluir um teste que cifra o mesmo plaintext duas vezes e asserta que os dois ciphertexts são diferentes (smoke test de unicidade de nonce).

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
