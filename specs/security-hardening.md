# Abditum — Hardening de Segurança

Catálogo de medidas de segurança além da criptografia do cofre. Organizadas por categoria e fase de implementação.

A criptografia central (Argon2id + AES-256-GCM) e a escrita atômica do arquivo são **requisitos v1** e estão documentadas em `decisions/tdr-001-criptografia.md` e `architecture.md` respectivamente. As medidas aqui são v2+, salvo indicação contrária.

---

## Proteção da interface

- **Proteção contra screenshots** *(v2)*: impedir que o SO capture a janela da aplicação. Protege contra malware de captura de tela e exposição acidental em gravações de vídeo e compartilhamentos de tela.
  - Windows: `SetWindowDisplayAffinity(HWND, WDA_EXCLUDEFROMCAPTURE)`
  - macOS: `NSWindow.sharingType = .none`
  - Linux/X11: não há equivalente universal; abordagem por compositor (Wayland oferece mais controle)

- **Accessibility APIs** *(v2)*: ferramentas de automação e leitores de tela (UIAutomation no Windows, AT-SPI no Linux) conseguem ler o conteúdo de campos mesmo mascarados. Marcar campos sensíveis como excluídos dessas APIs.

- **Impressão e exportação acidental** *(v3)*: desabilitar ou exibir aviso explícito em tentativas de imprimir ou exportar conteúdo que contenha valores de atributos sensíveis.

---

## Proteção de memória

- **Páginas de memória não exportáveis** *(v2)*: marcar as regiões de memória que contêm dados decifrados e a chave derivada como não paginável e não exportável.
  - Windows: `VirtualLock`
  - Linux/macOS: `mlock`

- **Chave em memória protegida** *(v2)*: manter a chave derivada em uma alocação especializada que zera automaticamente em panic, GC ou finalização — evita que o garbage collector do Go copie o buffer para outra região. Biblioteca de referência: `awnumar/memguard`.

- **Zeragem explícita de buffers sensíveis** *(v2)*: ao fechar o cofre, zerar explicitamente os buffers que contiveram a senha mestra, a chave derivada e os dados em texto claro antes de liberar a memória.
  - Windows: `SecureZeroMemory`
  - Go puro: `runtime.KeepAlive` + loop de zeragem (o compilador não pode otimizar para fora)

- **Zeragem do campo de senha após uso** *(v2)*: limpar o buffer do campo de entrada de senha imediatamente após a derivação da chave, não aguardar o fechamento do cofre.

- **Proteção contra hibernação e swap** *(v3)*: `mlock` já impede paginação em uso normal, mas o arquivo de hibernação do SO pode conter snapshots da memória. Considerar aviso ao usuário se hibernação estiver habilitada, e investigar APIs de exclusão de regiões do snapshot de hibernação.

---

## Proteção do arquivo

- **Escrita atômica** *(v1 — já implementado)*: gravar em arquivo temporário e renomear sobre o original somente após conclusão. Protege contra corrupção do cofre em falha de I/O, queda de energia ou encerramento abrupto.

- **Bloqueio de arquivo em uso** *(v2)*: impedir que outro processo leia ou escreva o `.abt` enquanto está aberto. Relevante especialmente para cofres armazenados em pastas de rede ou sincronizadas em nuvem.
  - Windows: `LockFileEx`
  - Linux/macOS: `flock` / `fcntl`

---

## Proteção da área de transferência

- **Limpeza automática após timeout** *(v1 — já implementado)*: conteúdo sensível copiado é apagado após 30 segundos. Ver `features/06-clipboard.md`.

- **Exclusão do histórico de clipboard do Windows** *(v2)*: o Windows 10+ mantém histórico de área de transferência que pode capturar senhas copiadas. Usar `CF_CLIPBOARD_EXCLUDE_FROM_MONITORING` ou a API de exclusão de entrada do histórico ao copiar atributos sensíveis.

---

## Proteção do processo

- **Desabilitar crash dumps** *(v2)*: evitar que o SO escreva um minidump contendo a memória do processo em caso de falha — um dump pode conter a chave derivada e dados em texto claro.
  - Windows: `SetErrorMode(SEM_NOGPFAULTERRORBOX)` + `MiniDumpWriteDump` desabilitado

- **Logs e telemetria zero** *(v1 — por design)*: nenhum dado sensível deve aparecer em logs do sistema, Event Viewer, syslog ou similares. Garantir que mensagens de erro nunca incluam valores de atributos ou a chave.

- **Anti-debug** *(v3, discutível)*: detectar se um debugger está anexado e recusar operação. Aumenta a barra para análise por malware, mas pode causar falsos positivos em ambientes de desenvolvimento. Considerar apenas em builds de release.

---

## Autenticação avançada

- **Autolock por inatividade** *(v2)*: fechar o cofre da memória automaticamente após um período configurável de inatividade, exigindo nova autenticação com senha mestra.

- **Segundo fator — YubiKey / FIDO2** *(v3)*: exigir token hardware além da senha mestra. Protege contra keylogger que capture a senha. A chave final seria derivada de `KDF(senha) XOR segredo_do_hardware`, de forma que nem a senha nem o hardware sozinhos sejam suficientes.

---

## Integridade da aplicação

- **Assinatura do executável (code signing)** *(v2)*: assinar o binário com certificado reconhecido pelo SO. Permite ao sistema verificar autoria antes de executar e protege contra substituição maliciosa do executável.
  - Windows: Authenticode
  - macOS: notarização Apple

- **Verificação de integridade do próprio binário** *(v3)*: ao iniciar, calcular o hash do executável e comparar com um valor de referência. Detecta modificação do binário em disco após distribuição.

---

## Resumo por fase

| Fase | Medidas |
|------|---------|
| **v1** | Escrita atômica, limpeza de clipboard, logs zero |
| **v2** | Screenshots, accessibility APIs, memória protegida (`memguard`, `mlock`, zeragem), bloqueio de arquivo, clipboard history exclusion, crash dumps, autolock, code signing |
| **v3** | Hibernação, impressão, YubiKey/FIDO2, anti-debug, verificação de binário |
