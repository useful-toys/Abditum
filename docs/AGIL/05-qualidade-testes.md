# 05 — Qualidade e Testes

## 05.1 Estratégia de Testes

### Pirâmide de testes

```
           ╱╲
          ╱  ╲         Integração E2E
         ╱ E2E╲        (fluxo completo: criar cofre → operar → salvar)
        ╱──────╲
       ╱ Golden ╲      Golden Files Visuais
      ╱  Files   ╲     (snapshot 80×24 por tela/estado)
     ╱────────────╲
    ╱  Comandos /  ╲   Testes de Comandos TUI
   ╱   Fluxos UI    ╲  (interação simulada por tela e fluxo)
  ╱──────────────────╲
 ╱    Unitários       ╲ Testes Unitários White-box
╱   (Domínio, Crypto,  ╲(Manager, entidades, criptografia,
╱     Armazenamento)     ╲ armazenamento, navegação, estados)
╱─────────────────────────╲
```

### Tipos de teste

| Tipo | Escopo | Ferramenta | Foco |
|---|---|---|---|
| **Unitário** | Funções e métodos isolados | `go test` | Lógica de domínio, Manager, regras de negócio, serviço de criptografia, serviço de armazenamento |
| **Golden Files** | Renderização visual da TUI | `teatest/v2` | Snapshot 80×24 de cada tela/estado — detecta regressões visuais |
| **Comandos TUI** | Interação simulada | `teatest/v2` | Sequências de teclas → estado esperado por fluxo de usuário |
| **Integração E2E** | Fluxo completo end-to-end | `go test` | Criar cofre → criar segredo → editar → salvar → reabrir → verificar |

---

### Cobertura por camada

#### Serviço de Criptografia

| Cenário | Tipo |
|---|---|
| Criptografar e descriptografar payload com sucesso | Unitário |
| Descriptografar com senha incorreta → erro de autenticação GCM | Unitário |
| Descriptografar com cabeçalho adulterado (AAD inválido) → erro de integridade | Unitário |
| Validar magic bytes (`ABDT`) — aceitar válidos, rejeitar inválidos | Unitário |
| Validar `versão_formato` — aceitar suportadas, rejeitar superiores | Unitário |
| Seleção de perfil Argon2id por `versão_formato` | Unitário |
| Salt gerado aleatoriamente (unicidade entre execuções) | Unitário |
| Nonce regenerado a cada chamada de criptografia | Unitário |
| Derivação de chave com parâmetros corretos (256 MiB, 3 iterações, 4 threads) | Unitário |

#### Serviço de Armazenamento

| Cenário | Tipo |
|---|---|
| Salvar cofre novo em caminho inexistente → arquivo criado | Unitário |
| Salvamento atômico: .tmp → .bak (com rotação .bak2) → rename | Unitário |
| Falha na escrita do .tmp → .tmp removido, .bak intacto | Unitário |
| Falha após criação do .bak → .bak2 restaurado para .bak | Unitário |
| Carregar cofre existente válido → payload descriptografado | Unitário |
| Carregar arquivo com magic inválido → erro de formato | Unitário |
| Carregar arquivo com versão superior → erro de incompatibilidade | Unitário |
| Migração de formato: payload v0 → modelo corrente em memória | Unitário |
| Sobrescrita com confirmação — backup do arquivo existente | Unitário |

#### Manager (Domínio e Regras de Negócio)

| Cenário | Tipo |
|---|---|
| Criar segredo a partir de modelo — campos copiados como snapshot | Unitário |
| Criar segredo vazio — sem campos iniciais além de nome/observação | Unitário |
| Editar segredo — nome, valores, observação alterados, identidade preservada | Unitário |
| Edição avançada — adicionar, renomear, excluir, reordenar campos | Unitário |
| Duplicar segredo — nova identidade, nome com sufixo incremental | Unitário |
| Favoritar/desfavoritar — apenas atributo `favorito` alterado | Unitário |
| Soft delete — segredo removido da hierarquia, presente na Lixeira | Unitário |
| Restaurar segredo — reinserido na pasta de origem (ou raiz se a pasta não existir) | Unitário |
| Mover segredo para outra pasta — identidade preservada | Unitário |
| Reordenar segredo entre irmãos — posição alterada, conteúdo intacto | Unitário |
| Criar pasta — inserida no final da coleção do destino | Unitário |
| Renomear pasta — identidade e posição preservadas | Unitário |
| Mover pasta — filhos acompanham recursivamente | Unitário |
| Excluir pasta — filhos promovidos ao pai | Unitário |
| Criar modelo — campos com nome e tipo | Unitário |
| Editar modelo — não afeta segredos existentes | Unitário |
| Excluir modelo — segredos existentes não afetados | Unitário |
| Criar modelo a partir de segredo — campos copiados | Unitário |
| Busca — match em nome, nome de campo, valor de campo texto, observação | Unitário |
| Busca — campos `texto sensível` nunca participam | Unitário |
| Importação — merge de pastas por identidade | Unitário |
| Importação — segredo com ID duplicado recebe novo ID | Unitário |
| Importação — segredo com nome duplicado recebe sufixo incremental | Unitário |
| Importação — modelo com ID duplicado é sobrescrito | Unitário |
| Configurações — alteração de tempos (bloqueio, reocultação, clipboard) | Unitário |
| Transição de estado — qualquer mutação → Cofre Modificado | Unitário |
| Transição de estado — salvar → Cofre Salvo | Unitário |
| Invariante — segredo em Lixeira não pode ser editado | Unitário |

#### Navegação e Estados

| Cenário | Tipo |
|---|---|
| Estado inicial → Tela Inicial com opções corretas | Unitário |
| Criar cofre → estado Cofre Salvo | Unitário |
| Abrir cofre → estado Cofre Salvo | Unitário |
| Modificação → estado Cofre Modificado | Unitário |
| Salvar → estado Cofre Salvo | Unitário |
| Bloquear → retorno a Tela Inicial com caminho preservado | Unitário |
| Descartar → reabrir arquivo, estado Cofre Salvo | Unitário |
| Sair sem salvar → confirmação bloqueante | Unitário |
| Bloqueio por inatividade — alerta aos 75% do tempo | Unitário |

#### Golden Files Visuais (80×24)

| Tela / Estado | Golden File |
|---|---|
| Tela inicial (ASCII art) | `golden_tela_inicial.txt` |
| Cofre ativo — hierarquia com segredos e pastas | `golden_cofre_ativo.txt` |
| Detalhe de segredo (campos oculto e visível) | `golden_detalhe_segredo.txt` |
| Edição padrão de segredo | `golden_edicao_padrao.txt` |
| Edição avançada de segredo | `golden_edicao_avancada.txt` |
| Busca ativa (árvore filtrada com highlight) | `golden_busca.txt` |
| Modal de confirmação (exclusão) | `golden_confirmacao.txt` |
| Toast de sucesso (cópia para clipboard) | `golden_toast_sucesso.txt` |
| Toast de erro (falha de salvamento) | `golden_toast_erro.txt` |
| File picker integrado | `golden_file_picker.txt` |
| Terminal abaixo do mínimo (mensagem de redimensionamento) | `golden_terminal_pequeno.txt` |

#### Testes de Comandos TUI

| Fluxo | Verificação |
|---|---|
| Tela Inicial → [C] → File Picker → senha 2x → Cofre Ativo | Estado Cofre Salvo, pastas/modelos pré-definidos presentes |
| Cofre Ativo → Ctrl+N → selecionar modelo → preencher → confirmar | Segredo criado na pasta correta com estado Novo |
| Cofre Ativo → navegar → selecionar segredo → Enter | Detalhe exibido no Painel do Segredo |
| Detalhe → toggle campo sensível → aguardar tempo | Campo reocultado automaticamente |
| Detalhe → copiar campo → aguardar tempo | Clipboard limpo, toast exibido com countdown |
| Cofre Ativo → Ctrl+F → digitar termo → selecionar resultado | Busca filtrada, resultado selecionado, busca encerrada |
| Cofre Modificado → Ctrl+Q → Descartar | Aplicação encerra sem salvar |
| Cofre Modificado → Ctrl+Q → Salvar | Salvamento atômico executado, aplicação encerra |

#### Integração E2E

| Fluxo | Verificação |
|---|---|
| Criar cofre → adicionar segredos → salvar → reabrir com mesma senha → verificar conteúdo | Dados íntegros, formato atual, hierarquia correta |
| Criar cofre → modificar → salvar como novo caminho → abrir original → verificar independência | Dois arquivos independentes, backup gerado no destino se existia arquivo anterior |
| Criar cofre → alterar senha mestra → salvar → reabrir com senha nova → verificar | Senha antiga rejeitada, senha nova aceita, conteúdo íntegro |
| Criar cofre → exportar JSON → importar em novo cofre → verificar conteúdo | Dados importados presentes, conflitos resolvidos conforme regras |
| Criar cofre v0 (simulado) → abrir com aplicação atual → verificar migração | Formato migrado em memória, salvamento regrava no formato corrente |

---

## 05.2 Cenários BDD / Gherkin

Os cenários BDD detalhados estão documentados na pasta `docs/BDD/features/`, organizados por domínio funcional. Abaixo, um resumo da cobertura:

| Domínio | Arquivo(s) | Features cobertas |
|---|---|---|
| **Cofre** | `cofre/*.md` | Criar, abrir, salvar, salvar como, alterar senha mestra, bloquear, configurar, descartar alterações, exportar/importar, sair |
| **Segredos** | `segredos/*.md` | Criar, editar, duplicar, favoritar, excluir/restaurar, mover/reordenar, copiar campo |
| **Pastas** | `pastas/*.md` | Gerenciar pastas (criar, renomear, mover, reordenar, excluir) |
| **Modelos** | `modelos/*.md` | Gerenciar modelos (criar, editar, excluir, criar a partir de segredo) |
| **Navegação** | `navegacao/*.md` | Navegar hierarquia, buscar segredos |
| **Segurança** | `seguranca/*.md` | Proteção e privacidade de dados |

Regras gerais de negócio e glossário BDD estão documentados em:
- `docs/BDD/regras-gerais.md`
- `docs/BDD/glossario.md`

---

## 05.3 Plano de Testes Não Funcionais

### Segurança

| Aspecto | Verificação | Método |
|---|---|---|
| **Criptografia em repouso** | Payload do `.abditum` é ilegível sem a senha mestra correta | Teste unitário: abrir arquivo sem senha → erro de autenticação |
| **Integridade do cabeçalho (AAD)** | Adulteração de qualquer byte do cabeçalho invalida a descriptografia | Teste unitário: alterar 1 byte de salt/nonce/magic/versão → erro GCM |
| **Anti-brute-force** | Derivação > 0,5 s com parâmetros v1 (256 MiB, 3 iter) | Teste de benchmark: medir tempo de derivação |
| **Privacidade de logs** | Nenhum dado sensível em stdout/stderr | Análise estática + teste verificando output |
| **Limpeza de memória** | Buffers controlados limpos ao bloquear/fechar | Teste unitário: verificar estado pós-bloqueio |
| **Limpeza de clipboard** | Clipboard limpo no tempo configurado e ao bloquear/fechar | Teste de comando TUI |
| **Campos sensíveis na busca** | Campos `texto sensível` nunca participam da busca | Teste unitário: buscar por valor de campo sensível → zero resultados |

### Portabilidade

| Aspecto | Verificação | Método |
|---|---|---|
| **Cross-platform** | Build e testes passam em Windows, macOS e Linux | CI pipeline com matrix de OS |
| **Binário autossuficiente** | Executa sem dependências externas ou instalação | Teste em ambiente limpo (container minimal) |
| **Zero arquivos externos** | Nenhuma leitura/escrita fora do `.abditum`, `.abditum.tmp`, `.abditum.bak` e `.abditum.bak2` | Monitoramento de syscalls de arquivo durante E2E |

### Confiabilidade

| Aspecto | Verificação | Método |
|---|---|---|
| **Salvamento atômico** | Falha no meio da escrita não corrompe arquivo original | Teste unitário com simulação de falha (mock de I/O) |
| **Rotação de backup** | `.bak2` restaurado para `.bak` em caso de falha | Teste unitário: simular falha após geração de backup |
| **Compatibilidade retroativa** | Cofre v0 abre e migra corretamente na versão atual | Teste E2E com fixture de arquivo v0 |
| **Versão futura** | Arquivo com versão superior → erro claro de incompatibilidade | Teste unitário: versão=999 → erro descritivo |

### Performance

| Aspecto | Verificação | Método |
|---|---|---|
| **Tempo de abertura** | Derivação Argon2id na faixa de 0,8 s a 1,5 s | Benchmark com parâmetros v1 |
| **Busca em memória** | Varredura de cofre com 1000 segredos em < 100 ms | Benchmark com cofre populado |
| **Renderização TUI** | Sem lag perceptível na navegação da árvore com 500+ nós | Teste manual + profiling |

### Usabilidade

| Aspecto | Verificação | Método |
|---|---|---|
| **Terminal mínimo** | Abaixo do mínimo → mensagem de redimensionamento | Golden file + teste de comando |
| **Feedback de ações** | Toasts exibidos conforme padrão (cor, ícone, tempo) | Golden files + testes de comando |
| **Confirmações bloqueantes** | Ações destrutivas exigem confirmação antes de executar | Testes de comando |
| **Ajuda contextual** | Barra de ajuda mostra apenas ações válidas para o contexto | Testes de comando por cada contexto |
