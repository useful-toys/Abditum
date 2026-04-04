# Modelo de Ameaças — Abditum

## Princípios de Segurança

- **Defesa em Profundidade**: Múltiplas camadas de proteção (KDF forte, Criptografia Autenticada, Proteção de Memória).
- **Conhecimento Zero Local**: O sistema é desenhado para que nem o desenvolvedor nem terceiros acessem os dados sem a senha mestra.
- **Minimização de Superfície**: Sem acesso à rede, sem dependências de C (estático), sem logs sensíveis e sem arquivos de configuração externos.
- **Segurança Pró-ativa em Sessão**: A segurança não termina na descriptografia; ela continua na gestão da memória, clipboard e interface.
- **Responsabilidade Compartilhada**: O Abditum garante a robustez algorítmica; o usuário garante a força da senha e a integridade do ambiente (SO).

## Ativos e Fronteiras

| Ativo | Estado | Fronteira de Proteção |
| :--- | :--- | :--- |
| **Segredos (Dados)** | Em repouso (Disco) | AES-256-GCM (Criptografia Autenticada) |
| **Segredos (Dados)** | Em uso (RAM) | `mlock`, Zeragem de buffers, `[]byte` mutável |
| **Senha Mestra** | Entrada (Teclado) | TUI Masking, conversão imediata para `[]byte` |
| **Senha Mestra** | Processamento | Argon2id (Derivação de chave lenta/custosa) |
| **Integridade do Cofre** | Arquivo `.abditum` | Tag GCM (AAD) + Magic Bytes |
| **Privacidade Visual** | TUI | Ocultação automática + Limpeza de Scrollback |
| **Arquivo de Intercâmbio** | Exportação (Disco) | Responsabilidade do usuário — sem proteção criptográfica após escrita |

---

## Matriz de Ameaças e Mitigações

### 1. Ataques Offline (Acesso ao arquivo)

| Ameaça | Descrição | Mitigação Abditum |
| :--- | :--- | :--- |
| **Força Bruta / Dicionário** | Tentativa de adivinhar a senha mestra por exaustão. | **Argon2id (m=256MiB, t=3)**: Torna cada tentativa extremamente lenta e cara em hardware (GPU/ASIC). |
| **Modificação de Bits** | Alterar o arquivo criptografado para causar comportamento errático. | **AES-GCM**: Qualquer alteração no ciphertext ou no cabeçalho (AAD) invalida a tag e impede a abertura. |
| **Vazamento de Metadados** | Descobrir nomes de pastas ou quantidade de segredos. | **JSON Envelopado**: Todo o modelo de domínio (pastas, nomes, notas) está dentro do payload criptografado. |
| **Arquivo de Intercâmbio Exposto** | A exportação produz um arquivo JSON sem criptografia. Se o usuário não o deletar manualmente, ele persiste no disco como superfície de ataque equivalente a um cofre sem senha — qualquer pessoa com acesso ao sistema de arquivos lê todos os segredos diretamente. | **Cerimônia obrigatória**: aviso de risco e confirmação antes de exportar. **Limitação residual**: o arquivo não possui TTL nem auto-deleção — o usuário é responsável por removê-lo após o uso. |
| **Proliferação via Sincronização em Nuvem** | Quando o cofre reside em Dropbox/OneDrive/similar, conflitos de sincronização geram cópias automáticas com sufixo (ex: `vault (Conflito de Jefferson).abditum`). Cada cópia é um alvo independente para ataque de força bruta offline — e o usuário muitas vezes nem sabe que as cópias existem. | **Argon2id**: cada cópia exige o mesmo custo computacional elevado para qualquer tentativa. **Detecção de modificação externa**: divergência de hash alerta ao salvar — mas não impede a proliferação de cópias nem as remove automaticamente. |

### 2. Ataques em Sessão (Acesso ao sistema rodando)

| Ameaça | Descrição | Mitigação Abditum |
| :--- | :--- | :--- |
| **RAM Scrapping** | Malware lendo a memória do processo para extrair senhas. | **Zeragem de Memória**: Sobrescrita de buffers sensíveis com zeros. **mlock**: Evita que a chave vá para o swap. |
| **Shoulder Surfing** | Alguém olhando a tela enquanto o usuário visualiza um segredo. | **Reveal Timer**: Campos sensíveis voltam a ficar ocultos automaticamente após X segundos. |
| **Keylogging** | Captura da senha mestra no momento da digitação. | **Limitação**: O Abditum depende da integridade do SO. Mitigado parcialmente pelo uso de bibliotecas TUI padrão. |
| **Clipboard Hijacking** | Outro app lendo o que foi copiado para o "Colar". | **Clipboard Auto-Clear**: Limpeza síncrona após tempo curto ou ao bloquear/sair. |
| **Sessão Abandonada** | Cofre aberto em terminal desassistido — acesso físico ao teclado ou sessão remota sem supervisão permite leitura de todos os segredos descriptografados sem exigir a senha mestra. | **Auto-lock por inatividade**: bloqueio automático configurável (padrão: 5 min) zera a memória sensível e exige re-autenticação para retomar. Qualquer interação do usuário reseta o temporizador. |
| **Phishing de UI** | Malware exibe uma tela falsa de bloqueio imitando a interface do Abditum — o usuário digita a senha mestra no processo malicioso acreditando estar re-autenticando o cofre. | **Restrição de re-autenticação**: o Abditum não solicita a senha mestra durante a sessão em operações rotineiras (salvar, descartar, mover) — apenas na abertura do cofre. Usuário treinado a questionar qualquer re-auth inesperada. **Limitação**: a aplicação não pode detectar uma sobreposição de tela realizada por outro processo. |
| **Race Condition na Limpeza de Clipboard** | `os.Exit` encerra goroutines em execução antes que uma limpeza assíncrona via goroutine ou `time.AfterFunc` possa executar, deixando o dado sensível na área de transferência após o encerramento do cofre. | **Limpeza síncrona obrigatória**: a limpeza ao bloquear/encerrar deve ocorrer no fluxo principal antes de qualquer chamada de saída. Detalhado em `arquitetura.md § Clipboard`. |

### 3. Ataques de Interface e Persistência

| Ameaça | Descrição | Mitigação Abditum |
| :--- | :--- | :--- |
| **Terminal Scrollback** | Ver dados sensíveis rolando o terminal para cima após fechar o app. | **Advanced Clear Screen**: Sequência `\033[3J` para tentar limpar o buffer de scrollback do emulador de terminal. |
| **Arquivos Temporários** | Recuperar dados de restos de salvamentos falhos. | **Salvamento Atômico**: Uso de `.tmp` seguido de `Rename`. Exclusão imediata do `.tmp` em caso de erro. |
| **Coredump do Processo** | Crash grave (OOM, SIGSEGV, kill externo) faz o SO escrever um coredump contendo todo o espaço de memória residente do processo — incluindo a chave AES derivada e os valores de campos descriptografados. | **mlock/VirtualLock (parcial)**: impede swap das páginas bloqueadas, mas não inibe geração de coredump. **Limitação**: a aplicação não invoca `prctl(PR_SET_DUMPABLE, 0)` (Linux) nem equivalente nos demais SOs para desabilitar dumps do processo. |
| **Concorrência: Múltiplas Instâncias** | Duas instâncias do Abditum abrindo o mesmo arquivo simultaneamente; a última a salvar sobrescreve silenciosamente as alterações da outra — sem aviso no momento da abertura, apenas na hora de salvar. | **Detecção de modificação externa**: comparação de tamanho + SHA-256 no momento do salvamento detecta divergência e oferece: Sobrescrever / Salvar como novo arquivo / Cancelar. A ausência de arquivo de lock é escolha consciente de portabilidade. |
| **Histórico de Shell** | O caminho do arquivo do cofre registrado no histórico do terminal (bash, zsh, PowerShell) revela a localização do ativo mais sensível para qualquer pessoa com acesso ao perfil do usuário ou ao arquivo de histórico. | **FilePicker via TUI**: o caminho é selecionado por diálogo interno — nunca passado como argumento de linha de comando. O histórico de shell não registra o caminho do cofre. |
| **Captura de Tela e Gravação** | Screenshot, gravação de tela ou compartilhamento de tela durante uma conferência captura campos sensíveis revelados na janela de exibição temporária, expondo-os de forma permanente em mídia externa ao controle do Abditum. | **Reveal Timer**: ocultação automática define a janela máxima de exposição (padrão: 15s). **Auto-lock**: bloqueia o cofre em terminais que permanecem em segundo plano sem supervisão. **Limitação**: o Abditum não possui mecanismo para impedir captura de tela por processos externos. |

### 4. Supply Chain e Ambiente de Execução

| Ameaça | Descrição | Mitigação Abditum |
| :--- | :--- | :--- |
| **Binário Comprometido** | Binário substituído por uma versão maliciosa — por ataque de distribuição, build não verificado ou comprometimento do canal de entrega — que exfiltra a senha mestra ou mantém um backdoor no cofre. | **Código-fonte aberto**: auditabilidade completa por qualquer pessoa. **Build estático (`CGO_ENABLED=0`)**: elimina a superfície de ataque de `LD_PRELOAD` e DLL-hijacking. O usuário deve obter o binário de fontes confiáveis e verificar o hash SHA-256 antes da execução. |
| **Comprometimento de Dependência** | Ataque supply chain a `golang.org/x/crypto` (que implementa Argon2id) ou ao próprio compilador Go invalida as garantias criptográficas centrais sem nenhuma mudança no código do Abditum. | **`go.sum` + `go.mod`**: integridade de cada dependência verificada por hash reproduzível a cada build. Dependências externas são intencionalmente mínimas — apenas `golang.org/x/crypto` e as bibliotecas charmbracelet. |
| **Terminal Hostil** | Emulador de terminal comprometido captura keystrokes (incluindo a senha mestra durante a digitação) ou intercepta a saída ANSI antes da renderização visual. | **Limitação**: o Abditum não pode auditar a integridade do processo de terminal que o hospeda. O usuário é responsável pela escolha de um emulador confiável e pela integridade geral do SO. |

---

## Premissas e Limitações

1.  **SO Comprometido**: Se o sistema operacional possui um malware com privilégios de root/admin ou um keylogger a nível de kernel, a segurança de qualquer software user-space (incluindo o Abditum) é considerada nula.
2.  **Entropia da Senha**: O Abditum avisa sobre senhas fracas, mas não as bloqueia. A segurança contra força bruta offline depende inteiramente da escolha do usuário.
3.  **Gerenciamento de Memória do Go**: Devido ao Garbage Collector (GC), existe um risco residual de que o runtime do Go mova um buffer sensível antes que possamos zerá-lo, deixando uma "sombra" na memória. O uso de `[]byte` minimiza isso, mas não elimina 100% como em C/Rust.
4.  **Scrollback e Clipboard como Melhor Esforço**: A limpeza de scrollback (`\033[3J`) não é garantida em todos os emuladores — terminais legados, `tmux`, `screen` e emuladores embarcados podem ignorá-la. A limpeza de clipboard pode falhar silenciosamente em sessões Wayland sem suporte explícito ao protocolo, ambientes headless (SSH sem `$DISPLAY`) e SOs onde o acesso à área de transferência não está disponível. Em ambos os casos a operação é executada incondicionalmente; a efetividade depende do ambiente do usuário.
5.  **Arquivo de Intercâmbio sem Ciclo de Vida Controlado**: O arquivo exportado é plaintext e não possui mecanismo de expiração ou auto-deleção. Uma vez escrito, sua proteção depende exclusivamente das permissões do sistema de arquivos e das ações do usuário. O Abditum não controla o ciclo de vida do arquivo após a escrita.

---

## Análise de Segurança: **9.0 / 10**

### Justificativa:

O Abditum atinge o estado da arte para gerenciadores de senhas locais e portáteis.

**Pontos Fortes (Por que 9.0):**
- **Criptografia Impecável**: A escolha de Argon2id (parâmetros conservadores) e AES-GCM com AAD é a recomendação atual das melhores práticas de segurança.
- **Foco em Memória**: A decisão de evitar `string` para dados sensíveis e implementar zeragem manual e `mlock` coloca o projeto acima da média de aplicações escritas em linguagens com GC.
- **Atomicidade e Integridade**: O protocolo de salvamento e a verificação de integridade via GCM garantem que o usuário não perca dados e não abra arquivos corrompidos/adulterados.
- **Zero Dependência C**: O binário estático (`CGO_ENABLED=0`) elimina vulnerabilidades de bibliotecas dinâmicas do sistema e facilita auditorias.

**O que falta para o 10.0:**
- **Risco do GC**: Como a linguagem Go não garante o controle absoluto sobre a alocação de memória (o GC pode copiar slices), existe um vetor teórico de persistência de dados na RAM que é inerente à stack escolhida.
- **Falta de 2FA Físico**: Para ser um "10", o projeto precisaria de suporte a chaves de hardware (YubiKey/FIDO2), o que foi descartado para manter a portabilidade extrema e dependência zero.
- **Dependência do Terminal**: A limpeza de scrollback buffer (`3J`) é baseada em "melhor esforço" e não funciona em 100% dos emuladores de terminal, o que é uma limitação da interface TUI em si.
- **Ausência de Proteção Contra Coredump**: A aplicação não desabilita geração de dumps via `prctl(PR_SET_DUMPABLE, 0)` (Linux) nem equivalentes nos demais SOs — um crash grave pode produzir um coredump contendo a chave AES derivada e valores de campos descriptografados em memória.
- **Arquivo de Intercâmbio Sem Auto-deleção**: O arquivo exportado persiste indefinidamente no disco. Um TTL automático, ou um prompt de confirmação de deleção imediatamente após o export, reduziria significativamente a janela de exposição.

**Conclusão**: É uma ferramenta extremamente robusta para o perfil de ameaça proposto (usuário que quer privacidade total e portabilidade sem depender da nuvem).
