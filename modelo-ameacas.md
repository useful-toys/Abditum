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

---

## Matriz de Ameaças e Mitigações

### 1. Ataques Offline (Acesso ao arquivo)

| Ameaça | Descrição | Mitigação Abditum |
| :--- | :--- | :--- |
| **Força Bruta / Dicionário** | Tentativa de adivinhar a senha mestra por exaustão. | **Argon2id (m=256MiB, t=3)**: Torna cada tentativa extremamente lenta e cara em hardware (GPU/ASIC). |
| **Modificação de Bits** | Alterar o arquivo criptografado para causar comportamento errático. | **AES-GCM**: Qualquer alteração no ciphertext ou no cabeçalho (AAD) invalida a tag e impede a abertura. |
| **Vazamento de Metadados** | Descobrir nomes de pastas ou quantidade de segredos. | **JSON Envelopado**: Todo o modelo de domínio (pastas, nomes, notas) está dentro do payload criptografado. |

### 2. Ataques em Sessão (Acesso ao sistema rodando)

| Ameaça | Descrição | Mitigação Abditum |
| :--- | :--- | :--- |
| **RAM Scrapping** | Malware lendo a memória do processo para extrair senhas. | **Zeragem de Memória**: Sobrescrita de buffers sensíveis com zeros. **mlock**: Evita que a chave vá para o swap. |
| **Shoulder Surfing** | Alguém olhando a tela enquanto o usuário visualiza um segredo. | **Reveal Timer**: Campos sensíveis voltam a ficar ocultos automaticamente após X segundos. |
| **Keylogging** | Captura da senha mestra no momento da digitação. | **Limitação**: O Abditum depende da integridade do SO. Mitigado parcialmente pelo uso de bibliotecas TUI padrão. |
| **Clipboard Hijacking** | Outro app lendo o que foi copiado para o "Colar". | **Clipboard Auto-Clear**: Limpeza síncrona após tempo curto ou ao bloquear/sair. |

### 3. Ataques de Interface e Persistência

| Ameaça | Descrição | Mitigação Abditum |
| :--- | :--- | :--- |
| **Terminal Scrollback** | Ver dados sensíveis rolando o terminal para cima após fechar o app. | **Advanced Clear Screen**: Sequência `\033[3J` para tentar limpar o buffer de scrollback do emulador de terminal. |
| **Arquivos Temporários** | Recuperar dados de restos de salvamentos falhos. | **Salvamento Atômico**: Uso de `.tmp` seguido de `Rename`. Exclusão imediata do `.tmp` em caso de erro. |

---

## Premissas e Limitações

1.  **SO Comprometido**: Se o sistema operacional possui um malware com privilégios de root/admin ou um keylogger a nível de kernel, a segurança de qualquer software user-space (incluindo o Abditum) é considerada nula.
2.  **Entropia da Senha**: O Abditum avisa sobre senhas fracas, mas não as bloqueia. A segurança contra força bruta offline depende inteiramente da escolha do usuário.
3.  **Gerenciamento de Memória do Go**: Devido ao Garbage Collector (GC), existe um risco residual de que o runtime do Go mova um buffer sensível antes que possamos zerá-lo, deixando uma "sombra" na memória. O uso de `[]byte` minimiza isso, mas não elimina 100% como em C/Rust.

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

**Conclusão**: É uma ferramenta extremamente robusta para o perfil de ameaça proposto (usuário que quer privacidade total e portabilidade sem depender da nuvem).
