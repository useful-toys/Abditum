# Especifica��o Visual � Abditum TUI

> Wireframes, layouts de componentes e fluxos visuais concretos.
> Cada tela e componente consome os padr�es definidos no design system.
>
> **Documento de funda��o:**
> - [`tui-design-system.md`](tui-design-system.md) � princ�pios, tokens, estados, padr�es transversais

## Sum�rio

- [Atalhos da Aplica��o](#atalhos-da-aplica��o)
- [Di�logos de Decis�o](#di�logos-de-decis�o)
- [Di�logos Funcionais](#di�logos-funcionais)
  - [PasswordEntry](#passwordentry)
  - [PasswordCreate](#passwordcreate)
  - [FilePicker](#filepicker)
    - [FilePicker � Modo Open](#filepicker--modo-open)
    - [FilePicker � Modo Save](#filepicker--modo-save)
  - [Help](#help)
- [Componentes](#componentes)
  - [Cabe�alho](#cabe�alho)
  - [Barra de Mensagens](#barra-de-mensagens)
  - [Barra de Comandos](#barra-de-comandos)
  - [Painel Esquerdo: �rvore](#painel-esquerdo-�rvore)
  - [Busca de Segredos](#busca-de-segredos)
  - [Painel Direito: Detalhe do Segredo � Modo Leitura](#painel-direito-detalhe-do-segredo--modo-leitura)
  - [Painel Direito: Detalhe do Segredo � Modo Edi��o de Valores](#painel-direito-detalhe-do-segredo--modo-edi��o-de-valores)
  - [Painel Direito: Detalhe do Segredo � Modo Edi��o de Estrutura](#painel-direito-detalhe-do-segredo--modo-edi��o-de-estrutura)
- [A��es na �rvore de Segredos](#a��es-na-�rvore-de-segredos)
  - [^D � Duplicar segredo](#d--duplicar-segredo)
  - [^M � Mover para outra pasta](#m--mover-para-outra-pasta)
  - [!? / !? � Reordenar segredo na lista](#--reordenar-segredo-na-lista)
  - [^R e ^C na �rvore � Atalhos de campo sens�vel](#r-e-c-na-�rvore--atalhos-de-campo-sens�vel)
- [Telas](#telas)
  - [Boas-vindas](#boas-vindas)

---

## Atalhos da Aplica��o

Este documento detalha as atribui��es espec�ficas de teclas para os fluxos e fun��es da aplica��o. As pol�ticas transversais de teclado e o agrupamento de teclas F por categoria de a��o s�o definidos no [Design System � Mapa de Teclas](tui-design-system.md#mapa-de-teclas).

### Atalhos Globais

| Tecla | A��o (Fluxo ou Fun��o) | Escopo | Notas |
|---|---|---|---|
| `F1` | Abrir / fechar modal de Ajuda | Global | |
| `F12` | Alternar Tema | Global | |
| `Ctrl+Q` | Sair da Aplica��o (Fluxos 3, 4, 5) | Global | Gerencia todas as sa�das com as devidas confirma��es |
| `Ctrl+Alt+Shift+Q` | Bloquear Cofre (Fluxo 6) | Global | Bloqueio emergencial, descarta altera��es, sem confirma��o. Atalho "complicado" para evitar acidentes. |

### Atalhos de �rea de Trabalho (Fluxos Principais)

Os seguintes atalhos disparam os fluxos principais da aplica��o quando a �rea de trabalho tem foco (sem di�logos abertos). Eles seguem os agrupamentos de teclas F definidos no Design System.

| Tecla | A��o (Fluxo) | Notas |
|---|---|---|
| `F2` | Modo Cofre (aba) | S� com cofre aberto |
| `F3` | Modo Modelos (aba) | S� com cofre aberto |
| `F4` | Modo Configura��es (aba) | Abrange o Fluxo 14: Configurar o Cofre |
| `F5` | Criar Novo Cofre (Fluxo 2) | |
| `F6` | Abrir Cofre Existente (Fluxo 1) | |
| `Shift+F6` | Descartar Altera��es e Recarregar Cofre (Fluxo 10) | Similaridade sem�ntica com F6 |
| `F7` | Salvar Cofre no Arquivo Atual (Fluxo 8) | |
| `Shift+F7` | Salvar Cofre em Outro Arquivo (Fluxo 9) | |
| `Ctrl+F7` | Alterar Senha Mestra (Fluxo 11) | Implica salvamento |
| `F8` | (Livre) | Reservado para futuras a��es de persist�ncia |
| `F9` | Exportar Cofre (Fluxo 12) | |
| `Shift+F9` | Importar Cofre (Fluxo 13) | |
| `F10` | Busca de Segredos � abrir/fechar campo | S� com cofre aberto e foco na �rvore; toggle |
| `F11` | (Livre) | |

> **Fluxo 7 � Aviso de Bloqueio Iminente por Inatividade:** � um fluxo iniciado pelo sistema, n�o requer um atalho manual do usu�rio.

---

## Di�logos de Decis�o

Todos os di�logos de decis�o seguem a anatomia comum e os padr�es de intera��o definidos no [design system � Di�logos](tui-design-system.md#di�logos), incluindo a [Refer�ncia Visual por Severidade](tui-design-system.md#severidade) e as [Regras de A��es na Borda Inferior](tui-design-system.md#a��es-na-borda-inferior).

---

## Cat�logo de Di�logos de Decis�o

Esta se��o lista todas as inst�ncias de di�logos de decis�o da aplica��o, especificando seu contexto, t�tulo, mensagem no corpo e a��es na borda. A estrutura visual � definida na se��o [Di�logos](tui-design-system.md#di�logos) do Design System.

| A��o | Situa��o | Tipo | T�tulo | Mensagem no Corpo | A��es na Borda |
|---|---|---|---|---|---|
| **Sair** | Sem altera��es | Confirma��o � Neutro | `Sair do Abditum` | `Sair do Abditum?` | `Enter Sair`, `Esc Voltar` |
| **Sair** | Com altera��es | Confirma��o � Alerta | `Sair do Abditum` | `Cofre modificado. Salvar ou descartar?` | `S Salvar`, `D Descartar`, `Esc Voltar` |
| **Salvar** | Conflito externo | Confirma��o � Destrutivo | `Salvar cofre` | `Arquivo modificado externamente. Sobrescrever?` | `S Sobrescrever`, `Esc Voltar` |
| **Abrir cofre** | Falha (arquivo inv�lido) | Notifica��o � Erro | `Abrir cofre` | `Arquivo corrompido ou inv�lido. Necess�rio fechar.` | `Enter OK` |
| **Abrir cofre** | Modifica��es n�o salvas | Confirma��o � Alerta | `Abrir cofre` | `Cofre modificado. Salvar ou descartar?` | `S Salvar`, `D Descartar`, `Esc Voltar` |
| **Abrir cofre** | Caminho/Formato inv�lido | Notifica��o � Erro | `Abrir cofre` | `Arquivo inv�lido ou vers�o n�o suportada. Necess�rio corrigir.` | `Enter OK` |
| **Abrir cofre** | Senha incorreta | Notifica��o � Erro | `Abrir cofre` | `Senha incorreta. Necess�rio tentar novamente.` | `Enter OK` |
| **Criar novo cofre** | Modifica��es n�o salvas | Confirma��o � Alerta | `Criar novo cofre` | `Cofre modificado. Salvar ou descartar?` | `S Salvar`, `D Descartar`, `Esc Voltar` |
| **Criar novo cofre** | Arquivo de destino existente | Confirma��o � Alerta | `Criar novo cofre` | `Arquivo '[Nome]' j� existe. Sobrescrever?` | `S Sobrescrever`, `I Outro caminho`, `Esc Voltar` |
| **Criar novo cofre** | Senhas n�o coincidem | Notifica��o � Erro | `Criar novo cofre` | `Senhas n�o conferem. Necess�rio digitar novamente.` | `Enter OK` |
| **Criar novo cofre** | Senha fraca | Confirma��o � Alerta | `Criar novo cofre` | `Senha � fraca. Prosseguir ou revisar?` | `P Prosseguir`, `R Revisar`, `Esc Voltar` |
| **Salvar cofre** | Conflito externo | Confirma��o � Destrutivo | `Salvar cofre` | `Arquivo modificado externamente. Sobrescrever ou salvar como novo?` | `S Sobrescrever`, `N Salvar como novo`, `Esc Voltar` |
| **Salvar cofre como** | Destino � arquivo atual | Notifica��o � Alerta | `Salvar cofre como` | `Destino n�o pode ser o arquivo atual. Necess�rio escolher outro.` | `Enter OK` |
| **Salvar cofre como** | Arquivo de destino existente | Confirma��o � Alerta | `Salvar cofre como` | `Arquivo '[Nome]' j� existe. Sobrescrever?` | `S Sobrescrever`, `I Outro caminho`, `Esc Voltar` |
| **Descartar e recarregar** | Arquivo modificado externamente | Confirma��o � Destrutivo | `Descartar e recarregar` | `Cofre modificado externamente. Prosseguir com recarregamento?` | `P Prosseguir`, `Esc Voltar` |
| **Descartar e recarregar** | Confirma��o de descarte | Confirma��o � Destrutivo | `? Descartar e recarregar` | `Todas as altera��es ser�o descartadas. Continuar?` | `C Continuar`, `Esc Voltar` |
| **Alterar senha mestra** | Senhas n�o coincidem | Notifica��o � Erro | `Alterar senha mestra` | `Senhas n�o conferem. Necess�rio digitar novamente.` | `Enter OK` |
| **Alterar senha mestra** | Senha fraca | Confirma��o � Alerta | `Alterar senha mestra` | `Senha � fraca. Prosseguir ou revisar?` | `P Prosseguir`, `R Revisar`, `Esc Voltar` |
| **Alterar senha mestra** | Conflito externo | Confirma��o � Destrutivo | `Alterar senha mestra` | `Arquivo modificado externamente. Sobrescrever?` | `S Sobrescrever`, `Esc Voltar` |
| **Exportar cofre** | Senha incorreta (reautentica��o) | Notifica��o � Erro | `Exportar cofre` | `Senha incorreta. Necess�rio tentar novamente.` | `Enter OK` |
| **Exportar cofre** | Riscos de seguran�a (n�o criptografado) | Confirma��o � Alerta | `Exportar cofre` | `Arquivo n�o criptografado. Expor dados sens�veis?` | `E Exportar`, `Esc Voltar` |
| **Exportar cofre** | Arquivo de destino existente | Confirma��o � Alerta | `Exportar cofre` | `Arquivo '[Nome]' j� existe. Sobrescrever?` | `S Sobrescrever`, `I Outro caminho`, `Esc Voltar` |
| **Importar cofre** | Arquivo de interc�mbio inv�lido | Notifica��o � Erro | `Importar cofre` | `Arquivo inv�lido ou sem Pasta Geral. Necess�rio corrigir.` | `Enter OK` |
| **Importar cofre** | Confirma��o da pol�tica de mesclagem | Confirma��o � Informativo | `Importar cofre` | `Pastas mescladas. Conflitos substitu�dos. Confirmar?` | `C Confirmar`, `Esc Voltar` |

---

## Di�logos Funcionais

Todos os di�logos funcionais seguem a anatomia comum do [design system � Di�logos](tui-design-system.md#di�logos), sem s�mbolo sem�ntico no t�tulo. Esta se��o especifica a anatomia interna de cada um.

---

### PasswordEntry

**Contexto de uso:** entrada de senha para abrir cofre.
**Token de borda:** `border.focused`
**Largura fixa:** 50 colunas

> Nos wireframes abaixo, `�` representa a �rea com fundo `surface.input` (campo de entrada). Na implementa��o real, o campo � uma �rea de fundo rebaixado sem hachura � conforme definido em [Campos de entrada de texto](tui-design-system.md#foco-e-navega��o).

**Estado inicial (campo vazio � a��o default bloqueada):**

```
?-- Senha mestra ----------------------------?
�                                            �
�  Senha                                     �
�  ���������������������������������������� �
�                                            �
?-- Enter Confirmar ------------- Esc Cancelar --?
       ? text.disabled (bloqueado)
```

**Estado com digita��o (a��o default ativa):**

```
?-- Senha mestra ----------------------------?
�                                            �
�  Senha                                     �
�  ����������������������������������������  �
�                                            �
?-- Enter Confirmar --------- Esc Cancelar --?
       ? accent.primary + bold (desbloqueado)
```

**Estado com contador de tentativas (a partir da 2�):**

```
?-- Senha mestra ----------------------------?
�                                            �
�  Senha                                     �
�  ����������������������������������������  �
�                                            �
�  Tentativa 2 de 5                          �
?-- Enter Confirmar --------- Esc Cancelar --?
```

| Elemento | Token | Atributo |
|---|---|---|
| T�tulo `Senha mestra` | `text.primary` | **bold** |
| Label `Senha` | `accent.primary` | **bold** (campo ativo, sempre � di�logo de campo �nico) |
| �rea do campo `�` | `surface.input` | � |
| Placeholder (antes de digitar) | `text.secondary` | *italic* |
| M�scara `��������` | `text.secondary` | � |
| Cursor `�` | `text.primary` | � |
| Contador `Tentativa 2 de 5` | `text.secondary` | � |
| A��o default (bloqueada) | `text.disabled` | � |
| A��o default (desbloqueada) | `accent.primary` | **bold** |

**Estados dos componentes:**

| Componente | Estado | Condi��o |
|---|---|---|
| Campo `Senha` | sempre vis�vel, sempre com foco | Di�logo de campo �nico |
| Contador de tentativas | vis�vel | Tentativa atual = 2 |
| Contador de tentativas | oculto | Primeira tentativa |
| A��o `Enter Confirmar` | bloqueada (`text.disabled`) | Campo `Senha` vazio |
| A��o `Enter Confirmar` | ativa (`accent.primary` **bold**) | Campo `Senha` n�o vazio |
| A��o `Esc Cancelar` | sempre ativa | � |

**Mensagens:**

| Contexto | Tipo | Texto |
|---|---|---|
| Di�logo abre / foco no campo (vazio ou v�lido) | Dica de campo | `� Digite a senha para desbloquear o cofre` |
| `Enter` ? senha incorreta | Erro (5s) | `? Senha incorreta` |
| Di�logo fecha (confirma��o ou cancelamento) | � | Barra limpa *(orquestrador assume)* |

**Comportamento:**
- M�scara de comprimento fixo (8 `�`) � n�o revela o tamanho real da senha
- Campo �nico � `Tab` n�o faz nada dentro deste di�logo

**Transi��es especiais:**

| Evento | Efeito |
|---|---|
| `Enter` com senha incorreta | Campo limpo; a��o default volta para `text.disabled`; contador incrementado |
| Tentativas esgotadas | Di�logo fecha automaticamente |

---

### PasswordCreate

**Contexto de uso:** cria��o de senha (ao criar cofre ou alterar senha mestra).
**Token de borda:** `border.focused`
**Largura fixa:** 50 colunas

**Estado inicial (foco no primeiro campo � a��o default bloqueada):**

```
?-- Definir senha mestra -------------------?
�                                            �
�  Nova senha                                �
�  ���������������������������������������� �
�                                            �
�  Confirma��o                               �
�  ���������������������������������������� �
�                                            �
?-- Enter Confirmar ----------------- Esc Cancelar --?
       ? text.disabled (bloqueado)
```

**Estado com digita��o (primeiro campo ativo, medidor aparece � a��o ainda bloqueada):**

```
?-- Definir senha mestra -------------------?
�                                            �
�  Nova senha                                �
�  ���������������������������������������� �
�                                            �
�  Confirma��o                               �
�  ���������������������������������������� �
�                                            �
�  For�a: ���������� Boa                     �
�                                            �
?-- Enter Confirmar ----------------- Esc Cancelar --?
       ? text.disabled (2� campo vazio)
```

**Estado com ambos campos preenchidos e senhas conferem (a��o default desbloqueada):**

```
?-- Definir senha mestra -------------------?
�                                            �
�  Nova senha                                �
�  ���������������������������������������� �
�                                            �
�  Confirma��o                               �
�  ���������������������������������������� �
�                                            �
�  For�a: ���������� Boa                     �
�                                            �
?-- Enter Confirmar ----------------- Esc Cancelar --?
       ? accent.primary + bold (desbloqueado)
```

**Estado com senhas divergentes (a��o default bloqueada � erro no campo):**

```
?-- Definir senha mestra -------------------?
�                                            �
�  Nova senha                                �
�  ���������������������������������������� �
�                                            �
�  Confirma��o                               �
�  ���������������������������������������� �
�                                            �
�  For�a: ���������� Boa                     �
�                                            �
?-- Enter Confirmar ----------------- Esc Cancelar --?
       ? text.disabled (senhas divergem)
```

| Elemento | Token | Atributo |
|---|---|---|
| T�tulo `Definir senha mestra` | `text.primary` | **bold** |
| Label do campo ativo | `accent.primary` | **bold** |
| Label do campo inativo | `text.secondary` | � |
| �rea do campo `�` | `surface.input` | � |
| Placeholder (antes de digitar) | `text.secondary` | *italic* |
| M�scara | `text.secondary` | � |
| Cursor `�` | `text.primary` | � |
| Medidor � preenchido | `semantic.success` ou `semantic.warning` | � |
| Medidor � vazio | `text.disabled` | � |
| Label de for�a `Boa` / `Forte` | `semantic.success` | � |
| Label de for�a `Fraca` | `semantic.warning` | � |
| A��o default (bloqueada) | `text.disabled` | � |
| A��o default (desbloqueada) | `accent.primary` | **bold** |

**Estados dos componentes:**

| Componente | Estado | Condi��o |
|---|---|---|
| Campo `Nova senha` | sempre vis�vel | � |
| Campo `Confirma��o` | sempre vis�vel | � |
| Medidor de for�a | vis�vel | Campo `Nova senha` n�o vazio |
| Medidor de for�a | oculto | Campo `Nova senha` vazio |
| Linha em branco antes do medidor | vis�vel | Medidor vis�vel |
| A��o `Enter Confirmar` | bloqueada (`text.disabled`) | Campo `Nova senha` vazio **ou** campo `Confirma��o` vazio **ou** senhas divergentes |
| A��o `Enter Confirmar` | ativa (`accent.primary` **bold**) | Ambos os campos n�o vazios **e** senhas conferem |
| A��o `Esc Cancelar` | sempre ativa | � |

> **Nota:** a verifica��o de igualdade entre as senhas ocorre **em tempo real** � a cada tecla no campo `Confirma��o` e ao abandonar o campo (Tab ou mudan�a de foco). Se as senhas divergem, a a��o default fica bloqueada e a barra de mensagens exibe erro no lugar da dica de campo.

**Mensagens:**

| Contexto | Tipo | Texto |
|---|---|---|
| Di�logo abre / foco em `Nova senha` (vazio ou v�lido) | Dica de campo | `� A senha mestra protege todo o cofre � use 12+ caracteres` |
| Foco em `Confirma��o` (vazio ou v�lido) | Dica de campo | `� Redigite a senha para confirmar` |
| Foco em `Confirma��o` (senhas divergentes) | Erro (5s) | `? As senhas n�o conferem � digite novamente` |
| Digita��o em `Confirma��o` (senhas divergentes) | Erro (5s) | `? As senhas n�o conferem � digite novamente` |
| `Enter` ? senhas divergentes | Erro (5s) | `? As senhas n�o conferem � digite novamente` |
| Di�logo fecha (confirma��o ou cancelamento) | � | Barra limpa *(orquestrador assume)* |

**Comportamento:**
- `Tab` alterna entre os campos `Nova senha` e `Confirma��o`
- Medidor de for�a atualizado a cada tecla no campo `Nova senha`
- M�scara de comprimento fixo (8 `�`) � n�o revela o tamanho real da senha
- Valida��o de igualdade em tempo real: a cada tecla no campo `Confirma��o` e ao abandonar o campo (Tab)
- Senhas divergentes: a��o default bloqueada (`text.disabled`); barra de mensagens exibe erro (`?`) no lugar da dica de campo; erro permanece at� que as senhas confiram ou o campo seja limpo

**Transi��es especiais:**

| Evento | Efeito |
|---|---|
| Digita��o em `Confirma��o` torna senhas iguais | Erro na barra � substitu�do pela dica de campo; a��o default muda para `accent.primary` **bold** |
| Digita��o em `Confirma��o` torna senhas diferentes | Dica de campo � substitu�da por erro (`?`, TTL 5s); a��o default volta para `text.disabled` |
| Abandonar `Confirma��o` (Tab) com senhas divergentes | Erro exibido na barra; foco move para `Nova senha`; a��o default bloqueada |
| Abandonar `Confirma��o` (Tab) com senhas iguais | Dica exibida na barra; foco move para `Nova senha`; a��o default ativa |

---

### FilePicker

**Contexto de uso:** abrir ou salvar arquivo do cofre.
**Token de borda:** `border.focused`
**Dimensionamento:** largura m�xima do DS (70 colunas ou 80% do terminal, o menor); altura 80% do terminal. Propor��o �rvore/arquivos ~40/60.
**Diret�rio inicial:** determinado pelo fluxo orquestrador. Se n�o informado, CWD do processo. Se o CWD n�o existe ou n�o tem permiss�o de leitura, fallback para home do usu�rio (`~`).
**Nome sugerido (modo Save):** determinado pelo fluxo orquestrador. Se n�o informado, campo inicia vazio. O campo n�o possui placeholder.
**Filtro de extens�o:** apenas arquivos com a extens�o `<ext>` (par�metro `extensao`) s�o exibidos no painel de arquivos. N�o h� campo de filtro edit�vel. Arquivos e diret�rios ocultos (nome iniciado com `.`) n�o s�o exibidos. A extens�o � omitida na exibi��o dos nomes de arquivo (redundante � o filtro j� restringe ao formato).
**Padding:** 2 colunas horizontal; **0 vertical** � exce��o ao DS [Dimensionamento de di�logos](tui-design-system.md#dimensionamento-de-di�logos). Justificativa: princ�pio "O Terminal como Meio" � espa�o vertical � recurso escasso; o FilePicker � o di�logo mais denso da aplica��o (caminho + 2 pain�is + campo `Arquivo:` no modo Save). As bordas `????` e os headers internos (`Estrutura`, `Arquivos`) criam conten��o e separa��o suficientes sem padding vertical.

O FilePicker opera em dois modos � **Open** e **Save** � com wireframes e condi��es distintos. Ambos compartilham a mesma anatomia de pain�is.

> Nos wireframes abaixo, `�` representa �reas com fundo `surface.input` (campos de entrada).

> **Decis�o de layout:** o FilePicker usa separadores internos com junctions em T (`+--�`) e pain�is lado a lado � estrutura que n�o se encaixa no modelo padr�o de di�logos do DS. Esta configura��o foi documentada como **exce��o justificada** (ver [DS � Exce��es ao dimensionamento](tui-design-system.md#dimensionamento-de-di�logos)) e n�o promoveu uma subse��o no DS porque: (1) o FilePicker � o �nico di�logo com essa complexidade; (2) � um padr�o de SO consolidado, n�o um padr�o reutiliz�vel interno; (3) o mecanismo de exce��o do DS cobre o caso. Se um segundo di�logo com pain�is internos surgir, a exce��o ser� promovida a subse��o.

**Barra de comandos durante FilePicker:** enquanto o FilePicker est� ativo, a barra de comandos exibe apenas as a��es internas do di�logo (conforme regra geral de [Barra de Comandos durante di�logo ativo](#anatomia)). A��es de confirma��o/cancelamento (`Enter`/`Esc`) j� est�o na borda do di�logo � n�o s�o duplicadas na barra.

```
  Tab Painel                                                                  F1 Ajuda
```

| A��o | Tecla | Descri��o |
|---|---|---|
| Alternar painel | `Tab` | Cicla foco entre os pain�is (�rvore ? Arquivos no modo Open; �rvore ? Arquivos ? Campo Nome no modo Save) |
| Ajuda | `F1` | Abre o Help � �ncora fixa |

---

#### Contrato de entrada e sa�da

**Entrada (par�metros do orquestrador):**

| Par�metro | Tipo | Obrigat�rio | Uso |
|---|---|---|---|
| `modo` | `Open \| Save` | Sim | Define t�tulo, a��es e presen�a do campo de nome |
| `extensao` | `String` | Sim | Extens�o filtrada e adicionada automaticamente ao salvar (ex: `".abditum"`, `".json"`). Deve incluir o ponto inicial. |
| `diretorio_inicial` | `PathBuf` | N�o | Diret�rio onde o FilePicker abre. Default: CWD ? fallback `~` |
| `nome_sugerido` | `String` | N�o (modo Save) | Valor inicial do campo `Arquivo:`. Default: vazio |

**Sa�da (retorno ao orquestrador):**

| Resultado | Valor | Significado |
|---|---|---|
| Confirmado | `Some(PathBuf)` | Caminho completo do arquivo selecionado (modo Open) ou caminho de salvamento com extens�o `<ext>` garantida (modo Save) |
| Cancelado | `None` | Usu�rio abandonou o di�logo via `Esc` |

---

#### FilePicker � Modo Open

**T�tulo:** `Abrir cofre`
**Objetivo:** selecionar um arquivo `<ext>` existente.

**Wireframe (arquivo selecionado � a��o default ativa, scroll em ambos os pain�is):**

```
?-- Abrir cofre -----------------------------------------------------?
�  /home/usuario/projetos/abditum                                    �
+- Estrutura -------------------- Arquivos --------------------------�
�  ? /                         ?  ? database   25.8 MB 15/03/25 14:32?
�    ? usuario                 �  ? config       1.2 KB 02/01/25 09:15�
�      ? documentos            �  ? backup      18.4 MB 04/04/25 18:47�
�      ? projetos              �                                     �
�        ? site                �                                     �
�        ? abditum             �                                     �
�          ? docs              �                                     �
�          ? src               �                                     �
�        ? outros              �                                     �
�      ? downloads             ?                                     ?
?-- Enter Abrir ------------------------------------- Esc Cancelar --?
       ? accent.primary + bold (desbloqueado)
```

> Scroll da �rvore (`?` `�` `?`) substitui o `�` do separador entre pain�is. Scroll dos arquivos (`?` `�` `?`) substitui o `�` da borda direita do modal. O `-` na borda inferior marca a jun��o do separador com a base do di�logo. Metadados (tamanho + `dd/mm/aa HH:MM`) na mesma linha do nome.

**Wireframe (nenhum arquivo � a��o default bloqueada, sem scroll):**

```
?-- Abrir cofre -----------------------------------------------------?
�  /home/usuario/documentos                                          �
+- Estrutura -------------------- Arquivos --------------------------�
�  ? /                         �                                     �
�    ? usuario                 �  Nenhum cofre neste diret�rio       �
�      ? documentos            �                                     �
�        ? fotos               �                                     �
�        ? textos              �                                     �
�                              �                                     �
?-- Enter Abrir ------------------------------------- Esc Cancelar --?
       ? text.disabled (bloqueado)
```

| Elemento | Token | Atributo |
|---|---|---|
| T�tulo `Abrir cofre` | `text.primary` | **bold** |
| Header `Estrutura` | `text.secondary` | **bold** |
| Header `Arquivos` | `text.secondary` | **bold** |
| Separadores internos (`+`, `-`, `-`, `-`, `�`) | `border.default` | � |
| Pasta selecionada na �rvore | `accent.primary` | **bold** |
| Pasta n�o selecionada | `text.primary` | � |
| Indicador de pasta (`?` recolhida, `?` expandida, `?` vazia) | `accent.secondary` | � |
| Arquivo selecionado no painel de arquivos | `special.highlight` (fundo) + `text.primary` | **bold** |
| Arquivo n�o selecionado | `text.primary` | � |
| Indicador de arquivo `?` | `text.secondary` | � |
| Nome do arquivo (sem extens�o `<ext>`) | � | Extens�o omitida na exibi��o � redundante com o filtro |
| Metadados (tamanho, data/hora) | `text.secondary` | � |
| Texto `Nenhum cofre neste diret�rio` | `text.secondary` | � |
| Valor do caminho | `text.secondary` | � |
| A��o default (bloqueada) | `text.disabled` | � |
| A��o default (desbloqueada) | `accent.primary` | **bold** |

**Estados dos componentes:**

| Componente | Estado | Condi��o |
|---|---|---|
| Painel `Estrutura` (�rvore) | sempre vis�vel | � |
| Painel `Arquivos` (lista) | conte�do vis�vel | Pasta selecionada cont�m arquivos `<ext>` |
| Painel `Arquivos` (lista) | texto vazio | Pasta selecionada **n�o** cont�m arquivos `<ext>` |
| Caminho (valor) | sempre vis�vel, somente leitura | Atualiza ao navegar na �rvore |
| Arquivo pr�-selecionado no painel | selecionado | Primeiro `<ext>` da pasta, automaticamente ao entrar na pasta |
| A��o `Enter Abrir` | bloqueada (`text.disabled`) | Pasta sob cursor n�o cont�m arquivos `<ext>` |
| A��o `Enter Abrir` | ativa (`accent.primary` **bold**) | Pasta sob cursor cont�m `<ext>` (pr�-sele��o autom�tica habilita a a��o, mesmo com foco na �rvore) |
| A��o `Esc Cancelar` | sempre ativa | � |

**Mensagens:**

| Contexto | Tipo | Texto |
|---|---|---|
| Di�logo abre / foco na �rvore | Dica de campo | `� Navegue pelas pastas e selecione um cofre` |
| Foco no painel de arquivos (com sele��o) | Dica de campo | `� Selecione o cofre para abrir` |
| Foco no painel de arquivos (painel vazio) | Dica de campo | `� Nenhum cofre neste diret�rio � navegue para outra pasta` |
| Di�logo fecha | � | Barra limpa *(orquestrador assume)* |

**Comportamento:**

- **Carregamento lazy:** a �rvore n�o carrega todo o filesystem na abertura. Apenas o caminho at� o diret�rio inicial � expandido. O conte�do de cada pasta � lido sob demanda ao expandir � evita lentid�o em filesystems grandes
- **Foco inicial:** �rvore de diret�rios (painel esquerdo)
- **Ordem do Tab:** �rvore ? Arquivos ? volta (2 stops)
- **Scroll:** cada painel tem scroll independente com indicadores `?`/`?`/`�` na borda direita do respectivo painel
- **Painel de arquivos reflete o cursor da �rvore:** ao mover o cursor (`??`) entre pastas na �rvore, o painel de arquivos atualiza imediatamente para mostrar os `<ext>` da pasta sob o cursor � n�o apenas ao expandir. O caminho exibido e o painel de arquivos acompanham a pasta com cursor, independente de ela estar expandida ou recolhida
- **Navega��o por teclado na �rvore:** `??` navega entre pastas vis�veis; `?` expande pasta recolhida; `?` recolhe pasta expandida; `Enter` avan�a foco para o primeiro arquivo no painel de arquivos (se a pasta sob o cursor cont�m `<ext>`; sem efeito se n�o cont�m); `Home`/`End` vai ao primeiro/�ltimo item vis�vel; `PgUp`/`PgDn` scroll por p�gina
- **Navega��o por teclado nos arquivos:** `??` navega entre arquivos; `Enter` confirma sele��o (equivale � a��o default); `Home`/`End` vai ao primeiro/�ltimo arquivo vis�vel; `PgUp`/`PgDn` scroll por p�gina
- Ao navegar para uma pasta na �rvore, se ela cont�m arquivos `<ext>`, o primeiro � pr�-selecionado automaticamente no painel de arquivos
- **Indicador de pasta vazia:** pastas sem subdiret�rios vis�veis usam `?` conforme o DS � n�o s�o expans�veis. `?` n�o tem efeito sobre elas (nada a expandir). `Enter` segue a regra padr�o � avan�a foco para o painel de arquivos se a pasta cont�m `<ext>`. `?` indica aus�ncia de subdiret�rios expans�veis � n�o impede que a pasta contenha arquivos `<ext>` exibidos no painel de arquivos
- **Clique simples em pasta:** move cursor para a pasta (atualiza painel de arquivos e caminho exibido)
- **Clique simples em arquivo:** seleciona o arquivo (highlight)
- **Duplo-clique em pasta:** expande/recolhe (mesmo que `?`/`?`)
- **Duplo-clique em arquivo:** confirma sele��o (mesmo que a��o default)
- **Scroll do mouse:** afeta o painel com foco
- **Arquivos e diret�rios ocultos** (nome iniciado com `.`) n�o s�o exibidos
- **Caminho longo:** truncado no in�cio com `�` (ex: `�/projetos/abditum`)
- **Diret�rios sem permiss�o:** exibidos normalmente na �rvore; ao tentar expandir, erro na barra (`? Sem permiss�o para acessar <pasta>`) e pasta permanece recolhida
- **Fallback de CWD:** se o CWD � inacess�vel, o FilePicker navega para home do usu�rio (`~`) e exibe mensagem informativa (`? Diret�rio atual inacess�vel � navegando para home`)

**Ordena��o:**

| Painel | Crit�rio | Detalhes |
|---|---|---|
| �rvore (pastas) | Alfab�tico, case-insensitive | Ordem lexicogr�fica (`a` = `A`) |
| Arquivos | Alfab�tico, case-insensitive | Ordem lexicogr�fica pelo nome sem extens�o |

**Indenta��o da �rvore:** 2 espa�os por n�vel de profundidade.

**Formato dos metadados:**

| Campo | Formato | Exemplo |
|---|---|---|
| Tamanho | `{valor} {unidade}` � base 1024, unidades KB/MB/GB, 1 casa decimal | `25.8 MB`, `1.2 KB`, `18.4 MB` |
| Data/hora | `dd/mm/aa HH:MM` � d�gitos num�ricos, locale local | `15/03/25 14:32` |

**Alinhamento dos metadados:** no painel de arquivos, os metadados s�o alinhados em colunas � tamanho alinhado � direita, data/hora em posi��o fixa. O nome do arquivo ocupa o espa�o restante � esquerda. Isso facilita a leitura por scanning vertical.

**Comportamento na raiz:** `?` na pasta raiz (`/`) n�o tem efeito � a sele��o permanece na raiz.

**Truncamento de metadados:** em terminais estreitos, os metadados s�o os primeiros a truncar (direita ? esquerda). O nome do arquivo tem prioridade e s� trunca se n�o houver espa�o mesmo para ele.

**Transi��es especiais:**

| Evento | Efeito |
|---|---|
| Cursor move para pasta sem `<ext>` | Painel de arquivos mostra texto vazio; a��o default muda para `text.disabled` |
| Cursor move para pasta com `<ext>` | Primeiro arquivo pr�-selecionado; a��o default muda para `accent.primary` **bold** |
| `Enter` no painel de arquivos | Di�logo fecha com o arquivo selecionado |
| `Enter` na �rvore (pasta com `<ext>`) | Foco avan�a para o primeiro arquivo no painel de arquivos |
| `Enter` na �rvore (pasta sem `<ext>`) | Sem efeito |
| `?` em pasta recolhida | Pasta expandida; cursor permanece na pasta |
| `?` em pasta expandida | Pasta recolhida; cursor permanece na pasta |
| `?` em pasta `?` (vazia) | Sem efeito (nada a expandir) |
| Tentar expandir pasta sem permiss�o | Erro na barra (`? Sem permiss�o para acessar <pasta>`); pasta permanece recolhida |

---

#### FilePicker � Modo Save

**T�tulo:** `Salvar cofre`
**Objetivo:** escolher diret�rio e nome para salvar o arquivo do cofre.

**Wireframe (campo nome preenchido � a��o default ativa):**

```
?-- Salvar cofre ----------------------------------------------------?
�  /home/usuario/projetos/abditum                                    �
+- Estrutura -------------------- Arquivos --------------------------�
�  ? /                         �  ? database   25.8 MB 15/03/25 14:32�
�    ? usuario                 �  ? config       1.2 KB 02/01/25 09:15�
�      ? projetos              �                                     �
�        ? abditum             �                                     �
�          ? docs              �                                     �
�                              �                                     �
+--------------------------------------------------------------------�
�  Arquivo: �meu-cofre������������������������������������������  �
?-- Enter Salvar ----------------------------------------- Esc Cancelar --?
       ? accent.primary + bold (desbloqueado)
```

**Wireframe (campo nome vazio � a��o default bloqueada):**

```
?-- Salvar cofre ----------------------------------------------------?
�  /home/usuario/projetos                                            �
+- Estrutura -------------------- Arquivos --------------------------�
�  ? /                         �  ? database   25.8 MB 15/03/25 14:32�
�    ? usuario                 �                                     �
�      ? projetos              �                                     �
�                              �                                     �
+--------------------------------------------------------------------�
�  Arquivo: ����������������������������������������������������  �
?-- Enter Salvar ----------------------------------------- Esc Cancelar --?
       ? text.disabled (bloqueado)
```

> Tokens de estrutura (t�tulo, headers, separadores, pasta, arquivo, metadados, caminho, a��o default) id�nticos ao [Modo Open](#filepicker--modo-open). Exclusivos do Modo Save:

| Elemento | Token | Atributo |
|---|---|---|
| R�tulo `Arquivo:` (campo ativo) | `accent.primary` | **bold** |
| R�tulo `Arquivo:` (campo inativo) | `text.secondary` | � |
| �rea do campo `�` | `surface.input` | � |
| Cursor `�` | `text.primary` | � |

**Estados dos componentes:**

| Componente | Estado | Condi��o |
|---|---|---|
| Painel `Estrutura` (�rvore) | sempre vis�vel | � |
| Painel `Arquivos` (lista) | conte�do vis�vel | Pasta selecionada cont�m arquivos `<ext>` |
| Painel `Arquivos` (lista) | texto vazio | Pasta selecionada **n�o** cont�m arquivos `<ext>` |
| Caminho (valor) | sempre vis�vel, somente leitura | Atualiza ao navegar na �rvore |
| Campo `Arquivo:` | sempre vis�vel | � |
| Caracteres inv�lidos para filesystem (`/ \ : * ? " < > \|`) | bloqueados silenciosamente | Tecla n�o produz efeito � sem mensagem de erro |
| Extens�o `<ext>` | adicionada automaticamente | Se o nome digitado n�o termina em `<ext>` |
| A��o `Enter Salvar` | bloqueada (`text.disabled`) | Campo `Arquivo:` vazio |
| A��o `Enter Salvar` | ativa (`accent.primary` **bold**) | Campo `Arquivo:` n�o vazio |
| A��o `Esc Cancelar` | sempre ativa | � |

> **Nota:** a valida��o de sobrescrita (arquivo j� existe) � responsabilidade do fluxo que chamou o FilePicker, n�o do di�logo. O picker retorna o caminho completo; o fluxo abre di�logo de Confirma��o � Destrutivo se necess�rio.

**Mensagens:**

| Contexto | Tipo | Texto |
|---|---|---|
| Di�logo abre / foco na �rvore | Dica de campo | `� Navegue pelas pastas e escolha onde salvar` |
| Foco no painel de arquivos | Dica de campo | `� Arquivos existentes neste diret�rio` |
| Foco no campo `Arquivo:` (vazio) | Dica de campo | `� Digite o nome do arquivo � <ext> ser� adicionado automaticamente` |
| Foco no campo `Arquivo:` (preenchido) | Dica de campo | `� Confirme para salvar o cofre` |
| Di�logo fecha | � | Barra limpa *(orquestrador assume)* |

**Comportamento:**

- **Foco inicial:** �rvore de diret�rios (painel esquerdo)
- **Ordem do Tab:** �rvore ? Arquivos ? Campo `Arquivo:` ? volta (3 stops)
- **Scroll:** cada painel tem scroll independente com indicadores `?`/`?`/`�` na borda direita do respectivo painel
- Navega��o na �rvore e painel de arquivos id�ntica ao modo Open, com uma exce��o: **`Enter` no painel de arquivos copia o nome (sem extens�o) para o campo `Arquivo:` e move foco para o campo** � n�o confirma o di�logo. A confirma��o requer `Enter` novamente (no campo ou em qualquer contexto com a��o default ativa)
- No painel de arquivos: `??` apenas destaca o arquivo (highlight) � **n�o** copia o nome para o campo. Somente `Enter` ou clique simples no arquivo copiam o nome (sem extens�o) para o campo `Arquivo:`
- Ao navegar na �rvore, o campo `Arquivo:` **n�o � limpo** � preserva o nome digitado
- Extens�o `<ext>` � adicionada silenciosamente ao caminho de retorno, sem alterar o texto exibido no campo
- **Duplo-clique em pasta:** expande/recolhe (mesmo que `?`/`?`)
- **Duplo-clique em arquivo existente:** copia o nome para o campo `Arquivo:`
- Scroll do mouse, arquivos ocultos, caminho longo, permiss�es, fallback CWD, ordena��o, indenta��o, formato de metadados e truncamento: id�ntico ao [Modo Open](#filepicker--modo-open)

**Transi��es especiais:**

| Evento | Efeito |
|---|---|
| Clique simples em arquivo existente no painel | Nome copiado para campo `Arquivo:`; a��o default muda para `accent.primary` **bold** |
| `Enter` no painel de arquivos | Nome copiado para campo `Arquivo:`; foco move para o campo. **N�o** confirma o di�logo |
| `Enter` na �rvore (pasta com `<ext>`) | Foco avan�a para o primeiro arquivo no painel de arquivos |
| `Enter` na �rvore (pasta sem `<ext>`) | Sem efeito |
| `?` em pasta recolhida | Pasta expandida; cursor permanece na pasta |
| `?` em pasta expandida | Pasta recolhida; cursor permanece na pasta |
| Limpar campo `Arquivo:` | A��o default volta para `text.disabled` |
| `Enter` com campo preenchido | Di�logo fecha com caminho completo (diret�rio + nome + `<ext>`) |
| Tentar expandir pasta sem permiss�o | Erro na barra (`? Sem permiss�o para acessar <pasta>`); pasta permanece recolhida |

---

### Help

**Contexto de uso:** lista todas as a��es do ActionManager, agrupadas. Acionado por `F1` em qualquer contexto.
**Token de borda:** `border.default` (di�logo de consulta, n�o recebe entrada de texto)
**Dimensionamento:** largura m�xima do DS; altura at� 80% do terminal. Scroll quando o conte�do excede a viewport.

**Wireframe (exemplo: Modo Cofre � segredo selecionado, sem scroll):**

```
?-- Ajuda � Atalhos e A��es ---------------------------------------?
�                                                                  �
�  Navega��o                                                       �
�  ??          Mover cursor na lista                               �
�  ? / Enter   Expandir pasta ou selecionar segredo                �
�  ?           Recolher pasta ou subir para pasta pai              �
�  Tab         Alternar foco entre pain�is                         �
�                                                                  �
�  Segredo                                                         �
�  Ctrl+R      Revelar / ocultar campo sens�vel                    �
�  Ctrl+C      Copiar valor para �rea de transfer�ncia             �
�  Insert      Novo segredo                                        �
�  ^E          Editar segredo                                      �
�  Delete      Excluir segredo                                     �
�                                                                  �
�  Cofre                                                           �
�  ^S          Salvar cofre                                        �
�  ^Q          Sair (salva se necess�rio)                          �
�  F1          Esta ajuda                                          �
�                                                                  �
?---------------------------------------------------- Esc Fechar --?
```

**Wireframe (exemplo: scroll � in�cio do conte�do, mais abaixo):**

```
?-- Ajuda � Atalhos e A��es ---------------------------------------?
�                                                                  �
�  Navega��o                                                       �
�  ??          Mover cursor na lista                               �
�  ? / Enter   Expandir pasta ou selecionar segredo                �
�  ?           Recolher pasta ou subir para pasta pai              �
�  Tab         Alternar foco entre pain�is                         �
�                                                                  �
�  Segredo                                                         �
�  Ctrl+R      Revelar / ocultar campo sens�vel                    ?
?---------------------------------------------------- Esc Fechar --?
```

> **Nota:** os wireframes s�o snapshots ilustrativos. O conte�do real � gerado dinamicamente pelo ActionManager a partir do contexto ativo.

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| T�tulo `Ajuda � Atalhos e A��es` | `text.primary` | **bold** |
| Label do grupo (`Navega��o`, `Segredo`, `Cofre`) | `text.secondary` | **bold** |
| Tecla (ex: `Ctrl+R`, `Insert`, `^S`) | `accent.primary` | � |
| Descri��o da a��o | `text.primary` | � |
| Seta de scroll (`?` / `?` na borda direita) | `text.secondary` | � |
| Thumb de posi��o (`�` na borda direita) | `text.secondary` | � |
| Borda | `border.default` | � |

#### Estados dos componentes

| Componente | Estado | Condi��o |
|---|---|---|
| Conte�do | sem scroll | Todas as a��es cabem na viewport |
| Conte�do | com scroll | A��es excedem a viewport � indicadores `?`/`?` e thumb `�` na borda direita (ver [DS � Scroll em di�logos](tui-design-system.md#scroll-em-di�logos)) |
| `F1` na barra de comandos | oculto (`HideFromBar`) | Enquanto o Help estiver aberto |
| Barra de comandos | vazia | Help n�o registra a��es internas na barra |

#### Eventos

| Evento | Efeito |
|---|---|
| `F1` pressionado (modal fechado) | Abre o modal; barra de comandos fica vazia; `F1` oculto |
| `F1` pressionado (modal aberto) | Fecha o modal; `F1` volta vis�vel na barra |
| `Esc` | Fecha o modal; `F1` volta vis�vel na barra |
| `?` / `?` | Scroll por linha (se conte�do excede viewport) |
| `PgUp` / `PgDn` | Scroll por p�gina (viewport - 1 linhas) |
| `Home` / `End` | Vai ao in�cio / fim do conte�do |

#### Comportamento

- **Conte�do din�mico** � gerado a partir de todas as a��es registradas no ActionManager no momento da abertura
- **Agrupamento** � a��es s�o organizadas pelo atributo num�rico `Grupo`. Cada grupo tem um `Label` registrado (ex: 1 ? "Navega��o", 2 ? "Segredo"). Grupos renderizados em ordem num�rica crescente
- **Ordena��o interna** � dentro de cada grupo, a��es ordenadas por `Prioridade` (maior primeiro)
- **Scroll** � segue o padr�o transversal do DS: indicadores `?`/`?` na borda direita, navega��o por `??` / `PgUp`/`PgDn` / `Home`/`End`
- **Borda inferior** � `Esc Fechar` sempre vis�vel, independente do estado de scroll

---

## Componentes

### Cabe�alho

**Responsabilidade:** contexto global � qual aplica��o, qual cofre, se h� altera��es pendentes e qual modo est� ativo na �rea de trabalho.
**Posi��o:** linhas 1�2 da tela (zona Cabe�alho do [DS � Dimensionamento](tui-design-system.md#dimensionamento-e-layout)).
**Altura fixa:** 2 linhas.

**Anatomia:**

| Linha | Papel | Conte�do |
|---|---|---|
| **1 � T�tulo** | Contexto + navega��o | Nome da app, `�` separador, nome do cofre, `�` dirty, abas de modo � direita |
| **2 � Separadora** | Divisa cabe�alho ? �rea de trabalho | Linha `-` full-width; a aba ativa "pousa" nesta linha via `? Texto ?` |

**Dois estados estruturais:**

| Estado | Linha 1 | Linha 2 | Abas |
|---|---|---|---|
| Sem cofre (boas-vindas) | Apenas nome da app | Separador simples, sem conectores | Ocultas |
| Cofre aberto | Nome da app `�` cofre `�` + abas | Separador com aba ativa suspensa | Vis�veis (3) |
| Busca ativa | Nome da app `�` cofre `�` + abas | Campo de busca � esquerda + aba ativa suspensa � direita | Vis�veis (3) |

---

#### Sem cofre (Boas-vindas)

> Wireframes ilustrativos a 80 colunas. A largura real acompanha o terminal.

```
  Abditum
----------------------------------------------------------------------------------
```

Sem nome de cofre, sem indicador dirty, sem abas. A linha separadora � cont�nua.

---

#### Cofre aberto � anatomia base

> Estado imposs�vel em opera��o normal (sempre h� um modo ativo). Mostrado para ilustrar a posi��o de todos os elementos antes de qualquer aba estar ativa.

**Sem altera��es:**

```
  Abditum � cofre                          ? Cofre ?  ? Modelos ?  ? Config ?
----------------------------------------------------------------------------------
```

**Com altera��es n�o salvas:**

```
  Abditum � cofre �                         ? Cofre ?  ? Modelos ?  ? Config ?
----------------------------------------------------------------------------------
```

O `�` aparece imediatamente ap�s o nome do cofre, em `semantic.warning`. Desaparece ap�s salvamento bem-sucedido.

---

#### Modo Cofre ativo

```
  Abditum � cofre �                      ?-------?  ? Modelos ?  ? Config ?
-----------------------------------------? Cofre ?------------------------------
```

A aba ativa na linha 1 substitui o texto por `-` (`?-------?`), mantendo a mesma largura da vers�o inativa (`? Cofre ?`). Na linha 2, o texto desce para o gap entre `?` e `?`, que se alinham verticalmente com `?` e `?` da linha 1 � conectando visualmente a aba � �rea de trabalho abaixo.

---

#### Modo Modelos ativo

```
  Abditum � cofre                          ? Cofre ?  ?---------?  ? Config ?
------------------------------------------------------? Modelos ?----------------
```

---

#### Modo Configura��es ativo

```
  Abditum � cofre                           ? Cofre ?  ? Modelos ?  ?--------?
--------------------------------------------------------------------? Config ?--
```

A aba mais � direita pode encostar na borda do terminal � `?` ocupa a �ltima coluna, sem `-` posterior.

> **Nota:** a aba Configura��es � referida como "Config" nos wireframes por economia de espa�o. O texto completo na implementa��o � `Config`.

---

#### Modo busca ativo

Ativo enquanto o campo de busca estiver aberto (ver [Busca de Segredos](#busca-de-segredos)). Dispon�vel apenas no Modo Cofre com cofre aberto.

A linha separadora (linha 2) � substitu�da pelo campo de busca. A aba ativa permanece suspensa � direita na mesma linha, sem altera��o de posi��o ou estilo.

**Campo aberto, sem query (rec�m-ativado):**

```
  Abditum � cofre �                      ?-------?  ? Modelos ?  ? Config ?
 - Busca: --------------------------------? Cofre ?--------------------------
```

**Campo aberto, com query:**

```
  Abditum � cofre �                      ?-------?  ? Modelos ?  ? Config ?
 - Busca: gmail --------------------------? Cofre ?--------------------------
```

**Regras de layout do campo na linha separadora:**

| Elemento | Largura | Notas |
|---|---|---|
| `- Busca: ` (prefixo fixo) | 9 colunas | `-` + espa�o + `Busca:` + espa�o |
| Texto da query | vari�vel | Em `accent.primary` **bold** |
| `-` preenchimento | restante - largura da aba ativa - 2 (margem direita m�nima) | Preenche at� a aba |
| Aba ativa (`? Texto ?`) | igual ao estado normal | Posi��o e estilo inalterados |

- **Query longa:** truncada � **esquerda** com `�` � a parte mais recente da query fica sempre vis�vel
- A largura dispon�vel para a query � calculada em tempo real e recalculada a cada resize do terminal

**Tokens exclusivos do modo busca na linha separadora:**

| Elemento | Token | Atributo |
|---|---|---|
| `- Busca: ` r�tulo | `border.default` | � |
| Texto da query | `accent.primary` | **bold** |
| `-` preenchimento | `border.default` | � |

> **Exce��o de layout documentada:** a linha separadora do cabe�alho tem papel estrutural fixo no DS (divisa cabe�alho ? �rea de trabalho). Durante o modo busca, essa linha assume papel adicional de display do campo de busca. Exce��o justificada pelo princ�pio **Hierarquia da Informa��o** � o campo imediatamente acima da �rvore cria rela��o visual direta entre query e resultado � e pelo princ�pio **O Terminal como Meio** � espa�o vertical � recurso escasso. Escopo-limitada ao Modo Cofre com busca ativa.

---

#### Mec�nica visual da aba ativa

A transforma��o de aba inativa ? ativa ocorre em duas linhas simult�neas:

| Linha | Aba inativa | Aba ativa |
|---|---|---|
| **1** | `? Texto ?` (borda + texto) | `?------?` (borda + preenchimento `-`) |
| **2** | `---------` (separador cont�nuo) | `? Texto ?` (gap com texto sobre `special.highlight`) |

Regras de alinhamento:

- A largura total da aba � **id�ntica** nos estados ativo e inativo
- `?` alinha-se verticalmente com `?` da linha acima
- `?` alinha-se verticalmente com `?` da linha acima
- O conte�do entre `?` e `?` (espa�o + texto + espa�o) tem fundo `special.highlight`
- As bordas `????` e o preenchimento `-` usam sempre `border.default`, independente do estado

---

#### Truncamento do nome do cofre

O espa�o dispon�vel para o nome do cofre � limitado � as abas ocupam largura fixa � direita. O componente calcula o espa�o em tempo real.

> **Extens�o `.abditum` � omitida** � a app s� trabalha com este formato, ent�o a extens�o � redundante. O nome exibido � o radical do arquivo (ex: `cofre.abditum` ? `cofre`).

**F�rmula:**

```
prefixo  = "  Abditum � "                             (12 colunas)
dirty    = " �"  se IsDirty(), ou ""                   (2 ou 0 colunas)
abas     = bloco de abas + espa�os entre elas           (largura fixa, ~32 colunas)
padding  = m�n. 1 coluna entre nome/dirty e abas

dispon�vel = largura_terminal - prefixo - dirty - abas - padding
```

**Algoritmo:**

1. Se o nome completo (radical sem extens�o) cabe ? exibir como est�
2. Se n�o cabe ? truncar com `�`: `{nome[0..n]}�` onde `n` � calculado para caber
3. Se nem 1 caractere + `�` (2 colunas) cabe ? exibir apenas `�`

**Prioridade de cess�o de espa�o:**

| Prioridade | Elemento | Comportamento |
|---|---|---|
| 1� (cede primeiro) | Nome do cofre | Truncado conforme algoritmo acima |
| 2� | Separador `�` e indicador `�` | Preservados enquanto houver espa�o |
| 3� (nunca cede) | Abas | Largura fixa, nunca truncadas |

**Wireframe � nome truncado (terminal ~80 colunas, modo Cofre):**

```
  Abditum � meu-cofre-pessoa� �          ?-------?  ? Modelos ?  ? Config ?
-----------------------------------------? Cofre ?------------------------------
```

O radical `meu-cofre-pessoal` foi truncado para `meu-cofre-pessoa�`.

---

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| `Abditum` (nome da app) | `accent.primary` | **bold** |
| `�` separador nome/cofre | `border.default` | � |
| Nome do cofre (radical, sem `.abditum`) | `text.secondary` | � |
| `�` indicador n�o salvo | `semantic.warning` | � |
| Bordas das abas (`????-`) � ativa e inativa | `border.default` | � |
| Aba ativa � fundo (gap entre `?` e `?`) | `special.highlight` | � |
| Aba ativa � texto | `accent.primary` | **bold** |
| Aba inativa � texto | `text.secondary` | � |
| Linha separadora | `border.default` | � |

---

#### Eventos

| Evento | Mudan�a visual |
|---|---|
| Cofre aberto com sucesso | Aparece `�` nome do cofre e as 3 abas |
| Cofre fechado / bloqueado | Desaparece nome do cofre e abas; volta ao estado boas-vindas |
| Altera��o em mem�ria (`IsDirty() = true`) | Aparece `�` em `semantic.warning` |
| Salvamento bem-sucedido (`IsDirty() = false`) | Desaparece `�` |
| Navega��o entre modos (Cofre / Modelos / Config) | Aba ativa muda; nova aba suspensa na linha separadora |
| Terminal redimensionado | Nome do cofre recalcula truncamento |

---

#### Comportamento

- **Abas clic�veis** � mouse troca o modo ativo ao clicar no texto ou na borda da aba (�rea de hit inclui linhas 1 e 2 da aba)
- **Navega��o por teclado** � `F2` Cofre, `F3` Modelos, `F4` Config (escopo �rea de trabalho � s� ativas com cofre aberto)
- **Indicador dirty** � aparece/desaparece imediatamente conforme `IsDirty()`, sem anima��o
- **Truncamento din�mico** � recalculado a cada renderiza��o (resize do terminal, mudan�a de modo ativo, cofre aberto/fechado)

---

### Barra de Comandos

**Responsabilidade:** exibir as a��es dispon�veis no contexto ativo � o usu�rio nunca precisa adivinhar o que pode fazer.
**Posi��o:** �ltima linha da tela (zona Barra de comandos do [DS � Dimensionamento](tui-design-system.md#dimensionamento-e-layout)).
**Altura fixa:** 1 linha.

**Princ�pio de conte�do:** a barra exibe apenas a��es de caso de uso (F-keys, atalhos de dom�nio, `^S`). Teclas de navega��o universais � `??`, `??`, `Tab`, `Enter`, `Esc` � s�o senso comum em TUI e n�o s�o exibidas. Exce��o: di�logos exibem a��es internas espec�ficas do contexto.

---

#### Anatomia

Cada a��o na barra segue o formato: **TECLA Label** � tecla em `accent.primary` **bold**, label em `text.primary`. A��es separadas por `�` em `text.secondary`. A a��o `F1` (Ajuda) � �ncora fixa na extrema direita.

**Estado normal:**

```
  ^I Novo � ^E Editar � Del Excluir � ^S Salvar                              F1 Ajuda
```

**Com a��o desabilitada (nenhum segredo selecionado):**

```
  ^I Novo � ^E Editar � ^S Salvar                                              F1 Ajuda
```

A��es com `Enabled = false` n�o aparecem na barra � s� no modal de Ajuda. O espa�o colapsa; separadores `�` s�o re-calculados entre a��es vis�veis.

**Durante di�logo ativo (apenas a��es internas):**

```
  Tab Campos � F5 Revelar                                                    F1 Ajuda
```

A��es do ActionManager ficam ocultas. A barra mostra apenas as a��es internas do di�logo do topo da pilha. A��es de confirma��o/cancelamento (`Enter`/`Esc`) j� est�o na borda do di�logo � n�o s�o duplicadas na barra.

**Espa�o restrito:**

```
  ^I Novo                                                                    F1 Ajuda
```

A��es de menor prioridade s�o ocultadas quando n�o h� espa�o. `F1` permanece sempre vis�vel � � via Help que o usu�rio descobre as a��es ocultas.

---

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Tecla da a��o (ex: `Insert`) | `accent.primary` | **bold** |
| Label da a��o (ex: `Novo`) | `text.primary` | � |
| Separador `�` | `text.secondary` | � |
| `F1` (Ajuda) | `accent.primary` | **bold** |

---

#### Atributos das a��es

Cada a��o registrada no ActionManager possui atributos que controlam sua apresenta��o:

| Atributo | Efeito na barra | Efeito no Help |
|---|---|---|
| `Enabled = true` | Exibida com estilo normal | Listada |
| `Enabled = false` | **N�o aparece** na barra | Listada |
| `HideFromBar = true` | **N�o aparece** na barra | Listada |
| `HideFromBar = false` | Exibida (se `Enabled`) | Listada |

Al�m destes:

- **Prioridade** � valor num�rico. Maior prioridade ? mais � esquerda na barra. Quando o espa�o � insuficiente, a��es de menor prioridade s�o removidas primeiro
- **Grupo** � valor num�rico. Usado exclusivamente no modal de Ajuda para organizar a��es. Grupos renderizados em ordem num�rica crescente. Dentro de cada grupo, a��es ordenadas por `Prioridade`. N�o afeta a barra de comandos
- **Label do grupo** � string registrada por grupo (ex: grupo 1 ? "Navega��o"). Exibido como t�tulo de se��o no Help em `text.secondary` bold

---

#### Eventos

| Evento | Mudan�a na barra |
|---|---|
| Troca de foco entre pain�is | A��es do painel que recebe foco ficam ativas |
| Sele��o de item na �rvore | A��es de item (editar, excluir, revelar) ficam `Enabled = true` � aparecem na barra |
| Nenhum item selecionado | A��es de item ficam `Enabled = false` � desaparecem da barra |
| Di�logo aberto (push na pilha) | Troca para a��es internas do di�logo |
| Di�logo fechado (pop da pilha) | Volta para a��es do ActionManager |
| Terminal redimensionado | Recalcula quais a��es cabem (prioridade governa corte) |

---

#### Comportamento

- **�ncora `F1`** � reserva espa�o fixo na extrema direita. O c�lculo de espa�o dispon�vel desconta `F1 Ajuda` antes de distribuir as demais a��es
- **A��es desabilitadas desaparecem da barra** � `Enabled = false` remove a a��o da barra (n�o fica exibida como dim). A a��o continua listada no Help
- **Di�logos de decis�o** (confirma��o/Notifica��o) � tipicamente n�o t�m a��es internas; a barra pode ficar vazia (apenas `F1 Ajuda`) enquanto o di�logo estiver aberto
- **Di�logos funcionais** (PasswordEntry, FilePicker etc.) � registram a��es internas (Tab entre campos, revelar senha, etc.) que aparecem na barra
- **Truncamento** � se mesmo a a��o de maior prioridade + `F1 Ajuda` n�o cabem, a barra mostra apenas `F1 Ajuda`

---

### Barra de Mensagens

**Responsabilidade:** comunicar feedback ao usu�rio � sucesso, erro, aviso, progresso, dicas.
**Posi��o:** 1 linha fixa entre a �rea de trabalho e a barra de comandos (zona Barra de mensagens do [DS � Dimensionamento](tui-design-system.md#dimensionamento-e-layout)).
**Altura fixa:** 1 linha.
**Anatomia:** borda `-` cont�nua na largura total do terminal. Quando h� mensagem, o texto (s�mbolo + `�` espa�o + conte�do) come�a com 2 espa�os de padding � esquerda (alinhado com o texto do cabe�alho), seguido de `-` at� o fim da linha. O espa�o entre s�mbolo e texto � sempre exatamente 1 caractere.

**Anatomia (exemplo � sucesso):**

```
-- ? Gmail copiado para a �rea de transfer�ncia --------------------------------
```

Todos os tipos seguem este padr�o. Diferen�as por tipo: `?` sucesso � `?` erro (**bold**) � `?` aviso � `????` spinner � `�` dica (*italic*) � `?` informa��o � sem mensagem (borda `-` cont�nua). Mensagem longa truncada com `�` no fim.

#### Tokens

Os tokens de cada tipo de mensagem s�o definidos no [DS � Mensagens](tui-design-system.md#mensagens). Adicional:

| Elemento | Token | Atributo |
|---|---|---|
| Borda `-` (sem mensagem) | `border.default` | � |
| Borda `-` (com mensagem) | `border.default` | � |

> A cor da borda n�o muda conforme o tipo de mensagem � apenas o texto embutido usa o token sem�ntico correspondente.

#### Estados dos componentes

| Componente | Estado | Condi��o |
|---|---|---|
| Borda | vis�vel (sem texto) | Nenhuma mensagem ativa |
| Borda + mensagem | vis�vel (texto embutido) | Mensagem ativa � tipo governa s�mbolo, cor e atributo |
| Texto | truncado com `�` | Mensagem excede largura dispon�vel (terminal - 2 padding - 2 borda m�nima) |

#### Eventos

| Evento | Efeito |
|---|---|
| Opera��o conclu�da com sucesso | Exibe `?` mensagem (`semantic.success`, TTL 5s) |
| Informa��o neutra | Exibe `?` mensagem (`semantic.info`, TTL 5s) |
| Condi��o de alerta (ex: bloqueio iminente) | Exibe `?` mensagem (`semantic.warning`, permanente, desaparece com input) |
| Falha em opera��o | Exibe `?` mensagem (`semantic.error` + bold, TTL 5s) |
| Opera��o em andamento | Exibe spinner `????` (`accent.primary`, permanente at� sucesso/erro) |
| Campo recebe foco (di�logo funcional) | Exibe `�` dica de campo (`text.secondary` italic) |
| Aplica��o emite dica proativa | Exibe `�` dica de uso (`text.secondary` italic) |
| TTL expira | Mensagem desaparece � volta � borda `-` |
| Nova mensagem emitida | Substitui imediatamente a mensagem anterior |
| Di�logo fecha | Barra � limpa � volta � borda `-` |

#### Comportamento

- **Borda permanente** � a borda `-` � sempre vis�vel, funcionando como separador entre a �rea de trabalho e a barra de comandos. Contribui para a estabilidade espacial
- **Uma mensagem por vez** � nova mensagem substitui a anterior imediatamente. N�o h� fila nem pilha
- **Texto embutido** � o texto (s�mbolo + conte�do) substitui o trecho central da borda, com `-` preenchendo os lados
- **Aviso re-emitido** � mensagens de aviso s�o re-emitidas a cada tick enquanto a condi��o persistir
- **Responsabilidade do orquestrador** � mensagens p�s-fechamento de di�logo (ex: "? Cofre aberto") s�o emitidas pelo orquestrador, n�o pelo di�logo

---

### Painel Esquerdo: �rvore

**Contexto:** �rea de trabalho � Modo Cofre.
**Largura:** ~35% da �rea de trabalho.
**Responsabilidade:** Exibir a hierarquia de pastas e segredos; permitir navega��o e sele��o do item a detalhar no painel direito.

**Wireframe (Modo Cofre � scroll ativo, segredo selecionado, painel com foco):**

```
  ? Favoritos          (2) ?
      ? Bradesco              �
      ? Gmail                 �
  ? Geral              (8)  �
    ? Sites e Apps     (5)  �
      ? Gmail           <�      ? special.highlight + bold (item selecionado)
      ? YouTube              �
      ? Facebook             �
  ? Financeiro         (3)  �
    ? Nubank                 ?
```

> `?`/`?` indicam conte�do al�m da �rea vis�vel; `�` � o thumb proporcional na posi��o `�`; `<�` marca o item sendo detalhado no painel direito. `<�` e scroll (`?`/`?`/`�`) ocupam a mesma coluna � o separador entre pain�is. Quando `<�` coincide com um indicador de scroll na mesma linha, `<�` tem prioridade (o indicador de scroll � suprimido naquela linha).

**Wireframe (item marcado para exclus�o � selecionado):**

```
    ? Sites e Apps     (5)  �
      ? Gmail           <�      ? special.highlight; `semantic.warning` + strikethrough
      ? YouTube              �
```

**Wireframe (cofre vazio):**

```
  ? Geral              (0)  �   ? special.highlight (pasta raiz selecionada)
                             �
                             �
```

Painel direito exibe placeholder "Cofre vazio" centralizado quando o cofre n�o tem nenhum segredo.

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Nome de item (normal) | `text.primary` | � |
| Fundo de item selecionado | `special.highlight` | � |
| Nome de item selecionado | `text.primary` | **bold** |
| `? ? ?` � prefixos de pasta | `text.secondary` | � |
| `?` � prefixo de segredo | `text.secondary` | � |
| `?` � prefixo de segredo favoritado | `accent.secondary` | � |
| `?` � prefixo de itens dentro de `? Favoritos` | `accent.secondary` | � |
| Nome da pasta virtual `Favoritos` | `accent.primary` | **bold** |
| Contadores `(n)` | `text.secondary` | � |
| Nome de segredo marcado para exclus�o | `semantic.warning` | ~~strikethrough~~ |
| `?` � prefixo de segredo marcado para exclus�o | `semantic.warning` | � |
| Nome de segredo rec�m-criado (n�o salvo) | `semantic.warning` | � |
| `?` � prefixo de segredo rec�m-criado | `semantic.warning` | � |
| Nome de segredo modificado (n�o salvo) | `semantic.warning` | � |
| `?` � prefixo de segredo modificado | `semantic.warning` | � |
| Nome de item desabilitado | `text.disabled` | dim |
| `�` separador � painel com foco | `border.focused` | � |
| `�` separador � painel sem foco | `border.default` | � |
| `<�` conector de sele��o no separador | `accent.primary` | � |
| `?` / `?` indicadores de scroll no `�` | `text.secondary` | � |
| `�` thumb de scroll no `�` | `text.secondary` | � |

#### Estados dos componentes

| Componente | Estado | Condi��o |
|---|---|---|
| `Favoritos` | vis�vel, expand�vel (`?/?`) | = 1 segredo favoritado |
| `Favoritos` | oculta | 0 segredos favoritados |
| Pasta ou segredo | `special.highlight` + texto **bold** | Cursor posicionado sobre o item |
| Pasta com filhos, expandida | prefixo `?` em `text.secondary` | Pasta n�o-vazia, aberta |
| Pasta com filhos, recolhida | prefixo `?` em `text.secondary` | Pasta n�o-vazia, fechada |
| Pasta sem filhos | prefixo `?` em `text.secondary` | Pasta vazia |
| Segredo (folha, limpo) | prefixo `?` em `text.secondary` | Segredo sem altera��es pendentes |
| Segredo rec�m-criado | prefixo `?` em `semantic.warning` + texto `semantic.warning` | Criado em mem�ria, ainda n�o salvo em disco |
| Segredo modificado | prefixo `?` em `semantic.warning` + texto `semantic.warning` | Editado em mem�ria, ainda n�o salvo em disco |
| Segredo marcado para exclus�o | prefixo `?` em `semantic.warning` + texto `semantic.warning` + ~~strikethrough~~ | Marcado para exclus�o, ainda n�o salvo |
| `<�` no separador | vis�vel | Foco da �rvore est� sobre um segredo |
| `<�` no separador | ausente � `�` normal | Nenhum segredo exibido no painel direito |
| `?`/`?`/`�` no `�` | vis�vel | Conte�do excede a �rea vis�vel do painel |
| Painel esquerdo | placeholder "Cofre vazio" � direita | Cofre sem nenhum segredo |

> **`<�` � `�`:** quando o item selecionado coincide com a posi��o do thumb, `<�` tem prioridade � mesma regra do DS para Di�logos em bordas.

#### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Painel recebe foco | Dica de campo | `� ?? para navegar` |
| `Favoritos` (a pasta) selecionada | Dica de campo | `� Pasta virtual � segredos permanecem na localiza��o original` |

#### Eventos

**Navega��o:**

**Navega��o � movimento linear:**

| Evento | Efeito na �rvore |
|---|---|
| Cursor desce uma linha | Foco move para o pr�ximo item vis�vel (respeitando expand/collapse); se j� est� no �ltimo item, n�o move |
| Cursor sobe uma linha | Foco move para o item anterior vis�vel; se j� est� no primeiro item, n�o move |
| Cursor vai ao primeiro item | Foco move para o topo absoluto da �rvore (primeiro item da lista, independente do scroll) |
| Cursor vai ao �ltimo item | Foco move para o �ltimo item vis�vel da �rvore |
| Scroll desce uma p�gina | Janela desliza viewport - 1 linhas para baixo; cursor vai para o item no topo da nova janela se estava fora dela |
| Scroll sobe uma p�gina | Janela desliza viewport - 1 linhas para cima; cursor vai para o item no fundo da nova janela se estava fora dela |

**Navega��o � movimento hier�rquico:**

| Evento | Efeito na �rvore |
|---|---|
| Avan�ar sobre pasta recolhida (`?`) | Pasta expandida; filhos tornam-se vis�veis; prefixo `?` ? `?`; foco salta para o primeiro filho vis�vel (subpasta ou segredo) |
| Avan�ar sobre pasta expandida (`?`) | Foco desce para o primeiro filho da pasta |
| Avan�ar sobre pasta vazia (`?`) | Sem efeito � pasta vazia n�o tem filhos para expandir |
| Avan�ar sobre segredo | Sem efeito de navega��o na �rvore � painel direito j� exibe o detalhe pelo foco |
| Recuar sobre filho de pasta | Foco sobe para a pasta pai |
| Recuar sobre pasta expandida | Pasta recolhida; prefixo `?` ? `?`; foco permanece na pasta |
| Recuar sobre pasta raiz (`Geral`) recolhida | Sem efeito � sem pai dispon�vel |
| Recuar sobre pasta raiz (`Geral`) expandida | Pasta recolhida; foco permanece na pasta raiz |

**Navega��o � foco entre pain�is:**

| Evento | Efeito na �rvore |
|---|---|
| Foco alternado para painel direito | `�` muda de `border.focused` para `border.default`; barra de comandos exibe a��es do painel direito |
| Foco recebido do painel direito | `�` muda de `border.default` para `border.focused`; barra de comandos exibe a��es da �rvore; cursor de campo vai para o item que estava com foco quando a �rvore perdeu foco |

**Navega��o � scroll visual:**

| Evento | Efeito na �rvore |
|---|---|
| Item em foco sai da �rea vis�vel (scroll para cima) | Janela rola automaticamente para manter o item em foco vis�vel |
| Item em foco sai da �rea vis�vel (scroll para baixo) | Janela rola automaticamente para manter o item em foco vis�vel |
| Conte�do total cabe na �rea vis�vel | Indicadores `?`/`?`/`�` desaparecem do `�` |
| Conte�do total n�o cabe na �rea vis�vel | `?` aparece se h� conte�do acima; `?` aparece se h� conte�do abaixo; `�` posicionado proporcionalmente |

**Navega��o � mouse:**

| Evento | Efeito na �rvore |
|---|---|
| Clique em item | Foco move para o item clicado (mesmo efeito de cursor com `?`/`?`) |
| Clique no prefixo `?` ou `?` | Pasta expande/recolhe � mesmo efeito de `?`/`?` sobre pasta |
| Clique no prefixo `?` | Sem efeito |
| Scroll do mouse para cima/baixo | Janela desliza; cursor acompanha se sair da �rea vis�vel |
| Clique em item dentro de `Favoritos` | Foco move para o atalho dentro de `Favoritos`; painel direito exibe o segredo referenciado |

**Navega��o � `Favoritos`:**

| Evento | Efeito na �rvore |
|---|---|
| Foco entra em `Favoritos` (pasta virtual) | Painel direito mant�m �ltimo segredo exibido; barra exibe dica "Pasta virtual � segredos permanecem na localiza��o original" |
| `Favoritos` expandida | Atalhos dos segredos favoritados tornam-se vis�veis; prefixo `?` ? `?` |
| `Favoritos` recolhida | Atalhos ocultados; prefixo `?` ? `?` |
| Foco em atalho dentro de `Favoritos` | Painel direito exibe o detalhe do segredo referenciado; `<�` aparece na linha do atalho |

**Segredo � cria��o e duplica��o:**

| Evento | Efeito na �rvore |
|---|---|
| Novo segredo criado (foco em pasta) | N� `? <novo>` inserido no final da pasta em foco; foco salta para o novo n�; contador da pasta e ancestrais +1 |
| Novo segredo criado (foco em segredo) | N� `? <novo>` inserido imediatamente abaixo do segredo em foco; foco salta para o novo n�; contador da pasta e ancestrais +1 |
| Segredo duplicado | N� `? <nome> (2)` inserido imediatamente abaixo do segredo original; foco salta para o duplicado; contador da pasta e ancestrais +1 |

**Segredo � edi��o de conte�do:**

| Evento | Efeito na �rvore |
|---|---|
| Nome do segredo alterado | Nome do n� atualizado imediatamente; se era `?`, prefixo muda para `?`; se j� era `?`, permanece `?` |
| Campo ou observa��o editado | Prefixo muda de `?` para `?` (apenas se `EstadoOriginal`; `?` permanece `?`) |

**Segredo � exclus�o e restaura��o:**

| Evento | Efeito na �rvore |
|---|---|
| Segredo marcado para exclus�o | Prefixo ? `?`; texto `semantic.warning` + strikethrough; contador da pasta e ancestrais -1; se favoritado, some de `Favoritos` |
| Exclus�o cancelada (restaura��o) | Prefixo original restaurado (`?`, `?`, `?` ou `?`); texto normal; contador da pasta e ancestrais +1; se era favoritado, volta a `Favoritos` |

**Segredo � favorito:**

| Evento | Efeito na �rvore |
|---|---|
| Segredo favoritado | Prefixo `?` ? `?` (se limpo); se j� era `?` ou `?`, prefixo dirty mantido (ver regra de prioridade em Comportamento); `Favoritos` aparece se era a primeira marca��o; atalho inserido em `Favoritos` |
| Segredo desfavoritado | Prefixo `?` ? `?` (se limpo); atalho removido de `Favoritos`; `Favoritos` desaparece se contagem chegar a 0 |

**Segredo � reordena��o e movimenta��o:**

| Evento | Efeito na �rvore |
|---|---|
| Segredo subido uma posi��o na pasta | N� sobe uma posi��o dentro da pasta; foco acompanha |
| Segredo descido uma posi��o na pasta | N� desce uma posi��o dentro da pasta; foco acompanha |
| Segredo reposicionado para posi��o espec�fica | N� move para a nova posi��o dentro da pasta; foco acompanha |
| Segredo movido para outra pasta | N� some da pasta de origem; aparece na pasta destino na posi��o especificada; foco acompanha o n� na nova posi��o; contadores de origem (-1) e destino (+1) e respectivos ancestrais atualizados |

**Pasta � cria��o e renomea��o:**

| Evento | Efeito na �rvore |
|---|---|
| Pasta criada | N� `? <nome>` inserido na posi��o especificada dentro do pai; foco salta para o novo n� |
| Pasta renomeada | Nome do n� atualizado imediatamente |

**Pasta � reordena��o e movimenta��o:**

| Evento | Efeito na �rvore |
|---|---|
| Pasta subida uma posi��o | N� sobe uma posi��o entre os irm�os; foco acompanha |
| Pasta descida uma posi��o | N� desce uma posi��o entre os irm�os; foco acompanha |
| Pasta reposicionada para posi��o espec�fica | N� move para a nova posi��o entre os irm�os; foco acompanha |
| Pasta movida para outro pai | N� some da posi��o atual; aparece dentro do novo pai; foco acompanha; hierarquia do novo pai atualizada |

**Pasta � exclus�o:**

| Evento | Efeito na �rvore |
|---|---|
| Pasta exclu�da (sem conflitos de nome) | N� da pasta removido; subpastas e segredos promovidos ao pai na posi��o da pasta exclu�da; contadores do pai recalculados; foco vai para o primeiro filho promovido (ou para o pai, se pasta era vazia) |
| Pasta exclu�da (com conflitos de nome) | Idem acima; segredos com conflito de nome exibidos com nome renomeado (sufixo `(N)`); barra de mensagens exibe alerta com lista de renomea��es |

**Cofre � persist�ncia:**

| Evento | Efeito na �rvore |
|---|---|
| Salvo com sucesso (mesmo arquivo) | N�s `?` removidos fisicamente da �rvore; prefixos `?` e `?` voltam a `?` ou `?` conforme o flag `favorito`; contadores recalculados; foco permanece no item atual |
| Salvo como (arquivo diferente) | Efeitos id�nticos ao salvar com sucesso � a �rvore n�o distingue o destino do arquivo |
| Salvo com outra senha | Efeitos id�nticos ao salvar com sucesso � a �rvore n�o conhece a chave de cifragem |
| Reverter altera��es (recarregar do disco) | �rvore completamente reconstru�da a partir do arquivo em disco: n�s `?` removidos (n�o existem no disco); n�s `?` voltam ao nome e prefixo originais (`?` ou `?`); n�s `?` voltam ao prefixo original (`?` ou `?`); contadores recalculados; se o item em foco ainda existe, foco permanece nele; se o item em foco era `?` (deixou de existir), foco vai para a pasta pai; `Favoritos` reconstru�da a partir dos dados do disco |

#### Comportamento

- **Espelho do cofre** � a �rvore � uma representa��o visual direta e sempre atualizada do estado do cofre. Qualquer muta��o no cofre � independentemente de onde ou como foi originada � deve se refletir imediatamente na �rvore. N�o existe estado interno da �rvore que contradiga o cofre.
- **Foco persiste sobre o mesmo elemento** � quando qualquer evento atualiza a �rvore (reordena��o, renomea��o, movimenta��o, exclus�o de outro item, salvar, reverter�), o foco permanece sobre o mesmo elemento, mesmo que sua posi��o na lista tenha mudado. O scroll se ajusta automaticamente para garantir que o elemento com foco esteja vis�vel.
- **Foco ao remover o elemento focado** � se o evento for a remo��o do pr�prio elemento com foco, o foco migra automaticamente seguindo a ordem de prefer�ncia: (1) elemento imediatamente abaixo na lista vis�vel; (2) se n�o existir, elemento imediatamente acima; (3) se a lista ficou vazia, `? Geral` (pasta raiz, que nunca pode ser removida).
- **Sele��o apenas por cor** � n�o h� s�mbolo de cursor. A sele��o � indicada exclusivamente pelo fundo `special.highlight`. Os prefixos (`? ? ? ? ? ? ? ?`) s�o estruturais e n�o mudam com a sele��o
- **Detalhe autom�tico** � o painel direito exibe o segredo que est� com foco na �rvore. Quando o foco est� sobre uma pasta, o painel mant�m o �ltimo segredo exibido. O detalhe n�o precisa ser "aberto" � � atualizado continuamente conforme o foco se move
- **Nome inicial de novo segredo** � `<novo>`; � o nome provis�rio que aparece no n� at� que o usu�rio edite o campo Nome no painel de detalhes
- **Segredos com altera��es pendentes** � tr�s prefixos indicam estado n�o salvo, todos em `semantic.warning` (mesma sem�ntica do `�` dirty no cabe�alho): `?` rec�m-criado, `?` modificado, `?` marcado para exclus�o (+ strikethrough). Todos desaparecem ap�s `^S` bem-sucedido
- **`Favoritos` � posi��o e comportamento** � quando vis�vel, � sempre o primeiro item da lista; se comporta como pasta normal (`?/?`); itens internos s�o atalhos para os segredos originais (os segredos permanecem na hierarquia de origem)
- **`Favoritos` � apari��o e remo��o** � o n� aparece instantaneamente quando o primeiro segredo � favoritado; desaparece instantaneamente quando o �ltimo segredo favoritado � desfavoritado. A atualiza��o segue o princ�pio "Espelho do cofre" � a �rvore reflete o estado do cofre imediatamente ap�s a execu��o da a��o
- **Foco preservado ao inserir/remover `Favoritos`** � quando o n� `Favoritos` aparece ou desaparece, a posi��o absoluta de todos os itens na lista desloca �1. O foco permanece sobre o mesmo elemento l�gico (identificado por identidade, n�o por �ndice). O scroll se ajusta automaticamente para manter o elemento em foco vis�vel
- **Favorito com estado dirty** � o prefixo dirty (`?`, `?`, `?`) substitui o `?` dentro de `Favoritos`; o `?` s� aparece como prefixo quando o segredo est� limpo. Prioridade de prefixo: `?` > `?` > `?` > `?` > `?`. Segredo marcado para exclus�o some imediatamente de `Favoritos` � permanece na hierarquia de origem com prefixo `?`
- **Navega��o linear ignora expand/collapse** � `?`/`?` navegam apenas entre itens *vis�veis*; filhos de pastas recolhidas s�o invis�veis e portanto pulados
- **`?` sobre segredo � no-op** � segredos s�o folhas; avan�ar sobre eles n�o tem efeito (o detalhe j� foi atualizado ao receber foco)
- **`?` tem dois comportamentos** � sobre pasta expandida, recolhe a pasta e foco permanece na pasta; sobre qualquer outro item (pasta recolhida, pasta vazia, segredo), sobe o foco para a pasta pai. Sobre a pasta raiz expandida, apenas recolhe
- **Foco ao retornar ao painel** � ao receber foco via Tab, o cursor restaura a posi��o anterior (n�o vai ao topo)
- **Scroll autom�tico** � o viewport se ajusta automaticamente para manter o item em foco vis�vel; nunca h� item em foco fora da �rea vis�vel
- **Scroll no separador** � o scroll da �rvore � indicado por `?`/`?`/`�` no `�` (separador entre pain�is). `<�` e scroll ocupam a mesma coluna: `<�` tem prioridade sobre `�` em caso de coincid�ncia (ver [DS � Scroll em di�logos](tui-design-system.md#scroll-em-di�logos)). Quando `<�` coincide com `?` ou `?`, `<�` prevalece � a dire��o do scroll � impl�cita pela presen�a do outro indicador nas demais linhas
- **Indenta��o** � 2 espa�os por n�vel de aninhamento

---

### Busca de Segredos

**Contexto de uso:** filtrar a �rvore de segredos por texto livre no Modo Cofre.
**Escopo:** dispon�vel apenas no **Modo Cofre**, com cofre aberto e foco no painel esquerdo (�rvore). Nos modos Modelos e Configura��es, `^F` e `F10` n�o t�m efeito de busca. O campo de busca na linha separadora do cabe�alho **s� aparece no Modo Cofre e apenas enquanto a busca estiver ativa** � nunca em outros modos, nunca na tela de boas-vindas.
**Modelo:** type-to-search � o campo na linha separadora do cabe�alho � display-only; o foco permanece na �rvore durante toda a intera��o.

---

#### Ativa��o e sa�da

| Mecanismo | Efeito |
|---|---|
| `^F` ou `F10` com campo **fechado** | Campo abre na linha separadora; barra de mensagens exibe dica; barra de comandos muda para a��es de busca |
| `^F` ou `F10` com campo **aberto** | Toggle: campo fecha; query descartada; �rvore restaurada; barra restaurada ao estado anterior |
| `Esc` com campo aberto | Id�ntico ao toggle com campo aberto; cursor retorna ao item que estava selecionado antes da busca |

> A busca **n�o pode ser ativada** com foco no painel direito (detalhe). O foco deve estar na �rvore.

---

#### Mapa de teclas durante busca ativa

| Tecla | Efeito |
|---|---|
| Alfanum�rica / s�mbolo imprim�vel | Acrescenta caractere � query; �rvore filtra em tempo real |
| `Backspace` | Remove o �ltimo caractere da query |
| `Del` | Limpa toda a query de uma vez; campo permanece aberto e vazio; �rvore restaurada completa |
| `?` / `?` | Navega entre os resultados vis�veis na �rvore filtrada |
| `Home` / `End` | Primeiro / �ltimo resultado vis�vel |
| `PgUp` / `PgDn` | Scroll por p�gina nos resultados |
| `Enter` com segredo selecionado | Abre detalhe no painel direito; campo permanece aberto |
| `Enter` com pasta selecionada | Expande / recolhe pasta; campo permanece aberto |
| `Tab` | Foco ? painel direito (detalhe do item selecionado); campo permanece aberto e vis�vel |
| `^F` / `F10` | Toggle: fecha o campo, descarta a query, restaura a �rvore |
| `Esc` | Fecha o campo, descarta a query, restaura a �rvore; cursor retorna ao item anterior |
| `F-keys` / `^Letra` | A��es normais da �rvore (ActionManager) � **n�o alimentam a query** |

> **Regra de roteamento:** apenas teclas que produzem caracteres imprim�veis (Unicode printable) e `Backspace` s�o interceptadas pela busca enquanto o campo estiver aberto. Modificadores, F-keys e teclas de controle passam normalmente ao ActionManager.

---

#### Comportamento do filtro

- **Correspond�ncia:** substring, case-insensitive, ignorando acentua��o � conforme requisito funcional
- **Escopo da busca:** nome do segredo, nome de campo, valor de campo **comum**, observa��o
- **Exclu�do da busca:** valores de campos sens�veis (nomes de campos sens�veis participam normalmente)
- **Exclu�dos dos resultados:** segredos marcados para exclus�o (`?`)
- **�rvore compacta:** apenas pastas que cont�m = 1 resultado s�o exibidas; pastas sem resultados desaparecem completamente
- **Contadores de pasta durante filtro ativo:** formato `(N/Total)` � `N` = segredos que atendem � busca nessa pasta e subpastas; `Total` = total de segredos ativos nessa pasta e subpastas. Exemplo: `(2/6)` significa que 2 dos 6 segredos atendem � query. Quando `N = Total`, o contador volta ao formato simples `(N)` � sem barra. O formato `(N/Total)` s� aparece durante busca ativa com query n�o vazia
- **Indicador visual de filtro ativo:** o painel esquerdo exibe `Filtro ativo` em `semantic.warning` + *italic*, alinhado � direita na primeira linha da �rea de trabalho, quando h� query n�o vazia. Garante percep��o do filtro mesmo que o cabe�alho esteja fora da viewport ou o foco esteja no painel direito
- **Match highlight:** o trecho de texto correspondente � query � exibido em `special.match` + **bold**
- **Query vazia:** campo aberto sem texto � �rvore exibe tudo; contadores voltam ao formato `(N)`; indicador `Filtro ativo` n�o aparece
- **Persist�ncia:** ao fechar o campo, a query � descartada e a �rvore restaurada completa; o campo sempre abre vazio

---

#### Wireframes

**Campo aberto, sem query (rec�m-ativado):**

```
  Abditum � cofre �                      ?-------?  ? Modelos ?  ? Config ?
 - Busca: --------------------------------? Cofre ?--------------------------
  ? Favoritos          (2)  �
    ? Bradesco         <�
    ? Gmail                 �
  ? Geral              (8)  �
    ? Sites            (5)  �
      ? Gmail               �
      ? YouTube             �
 - � Digite para filtrar os segredos ----------------------------------------
  ^F Fechar � Del Limpar                                              F1 Ajuda
```

> Query vazia: �rvore completa, contadores no formato `(N)`, sem indicador `Filtro ativo`.

**Campo aberto, com query � resultados encontrados:**

```
  Abditum � cofre �                      ?-------?  ? Modelos ?  ? Config ?
 - Busca: gmail --------------------------? Cofre ?--------------------------
  ? Favoritos        (1/2)  �              ? Filtro ativo
    ? Gmail            <�       ? match em special.match + bold
  ? Geral            (2/8)  �
    ? Sites          (2/5)  �
      ? Gmail               �
      ? Gmail Pro           �
 - ? 3 resultado(s) ---------------------------------------------------------
  ^F Fechar � Del Limpar                                              F1 Ajuda
```

> `Filtro ativo` em `semantic.warning` + *italic*, alinhado � direita. `(1/2)` = 1 resultado dos 2 segredos em Favoritos. Quando `N = Total`, contador volta a `(N)`.

**Campo aberto, sem resultados:**

```
  Abditum � cofre �                      ?-------?  ? Modelos ?  ? Config ?
 - Busca: xyzxyz -------------------------? Cofre ?--------------------------
  ? Geral              (0)  �              ? Filtro ativo
                             �
                             �
 - ? Nenhum resultado -------------------------------------------------------
  ^F Fechar � Del Limpar                                              F1 Ajuda
```

> Pasta raiz sempre vis�vel, mesmo sem resultados. Indicador `Filtro ativo` permanece.

**Campo aberto, query longa (truncada � esquerda):**

```
  Abditum � cofre �                      ?-------?  ? Modelos ?  ? Config ?
 - Busca: �ail.google.com/conta ----------? Cofre ?--------------------------
```

> A parte mais recente da query (direita) fica sempre vis�vel. `�` substitui os caracteres iniciais quando a query excede o espa�o dispon�vel.

---

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| `- Busca: ` r�tulo na linha separadora | `border.default` | � |
| Texto da query | `accent.primary` | **bold** |
| `-` preenchimento na linha separadora | `border.default` | � |
| Trecho de match na �rvore | `special.match` | **bold** |
| Contador `(N/Total)` durante filtro ativo | `text.secondary` | � |
| Indicador `Filtro ativo` | `semantic.warning` | *italic* |

---

#### Estados dos componentes

| Componente | Estado | Condi��o |
|---|---|---|
| Campo de busca na linha separadora | oculto | Campo fechado � linha separadora normal |
| Campo de busca na linha separadora | vis�vel, vazio | Campo aberto, query vazia |
| Campo de busca na linha separadora | vis�vel, com texto | Query ativa (= 1 caractere) |
| Campo de busca na linha separadora | **nunca vis�vel** fora do Modo Cofre | Modos Modelos, Configura��es, Boas-vindas |
| �rvore | completa | Campo fechado **ou** campo aberto com query vazia |
| �rvore | filtrada (compacta) | Campo aberto com query = 1 caractere |
| Pasta | vis�vel | Cont�m = 1 resultado direto ou indireto |
| Pasta | oculta | N�o cont�m nenhum resultado |
| Pasta raiz | sempre vis�vel | Mesmo sem resultados � exibe `(0)` e `?` |
| Contador de pasta | formato `(N)` | Campo fechado **ou** query vazia **ou** `N = Total` |
| Contador de pasta | formato `(N/Total)` | Query ativa com = 1 caractere e `N < Total` |
| Indicador `Filtro ativo` | vis�vel, 1� linha da �rea de trabalho, alinhado � direita | Query ativa com = 1 caractere |
| Indicador `Filtro ativo` | oculto | Campo fechado ou query vazia |
| Trecho de match | `special.match` + **bold** | Substring correspondente � query |
| Barra de comandos | a��es de busca (`^F Fechar � Del Limpar`) | Campo aberto |
| Barra de comandos | a��es normais da �rvore | Campo fechado |

---

#### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Campo abre (query vazia) | Dica de uso | `� Digite para filtrar os segredos` |
| Query ativa, com resultados | Informa��o | `? N resultado(s)` |
| Query ativa, sem resultados | Informa��o | `? Nenhum resultado` |
| `Backspace` apaga �ltimo caractere � query fica vazia | Dica de uso | `� Digite para filtrar os segredos` |
| `Del` limpa a query | Dica de uso | `� Digite para filtrar os segredos` |
| Campo fecha (`Esc`, `^F`, `F10`) | � | Barra restaurada ao estado anterior � busca |

---

#### Barra de comandos durante busca ativa

```
  ^F Fechar � Del Limpar                                              F1 Ajuda
```

As a��es normais da �rvore (ActionManager) ficam ocultas na barra enquanto o campo estiver aberto � o ActionManager continua processando suas teclas (`^Letra`, `F-keys`), mas a barra reflete apenas o contexto de busca.

---

#### Transi��es especiais

| Evento | Efeito |
|---|---|
| `^F` / `F10` � campo fechado | Campo abre; separadora substitu�da; barra muda; dica exibida |
| `^F` / `F10` � campo aberto | Campo fecha; query descartada; separadora restaurada; cursor volta ao item anterior; barra restaurada |
| `Esc` � campo aberto | Id�ntico ao toggle com campo aberto |
| Digita��o � query n�o vazia | �rvore filtra em tempo real; `? N resultado(s)` atualiza a cada caractere |
| `Backspace` � query vazia ap�s apagar | �rvore restaurada completa; campo permanece aberto; dica exibida |
| `Del` | Query limpa instantaneamente; campo permanece aberto; �rvore restaurada; dica exibida |
| `Enter` � segredo selecionado | Detalhe atualizado no painel direito; campo permanece aberto |
| `Enter` � pasta selecionada | Pasta expande / recolhe; campo permanece aberto |
| `Tab` � foco na �rvore | Foco vai para painel direito; campo permanece aberto e vis�vel; type-to-search suspende at� foco retornar � �rvore |
| Foco retorna � �rvore (`Tab` / clique) | Type-to-search retoma � teclas alfanum�ricas voltam a alimentar a query |
| Terminal redimensionado | Largura dispon�vel da query recalculada; truncamento com `�` reaplicado se necess�rio |

---

## A��es na �rvore de Segredos

Esta se��o detalha as a��es dispon�veis ao interagir com a �rvore de segredos (painel esquerdo do Modo Cofre) e seus respectivos atalhos de teclado. As regras gerais de navega��o e atribui��o de teclas s�o definidas no [Design System � Mapa de Teclas](tui-design-system.md#mapa-de-teclas).

### Navega��o na �rvore (geral)

| Tecla           | A��o                                     | Notas                                            |
|-----------------|------------------------------------------|--------------------------------------------------|
| `?` / `?`       | Mover cursor na lista / �rvore           |                                                  |
| `Home` / `End`  | Mover para o primeiro / �ltimo item vis�vel |                                                  |
| `PgUp` / `PgDn` | Rolar uma p�gina para cima / baixo       |                                                  |
| `Tab`           | Alternar foco entre pain�is              | Move o foco para o painel direito (Detalhe) e vice-versa. |

### A��es em pastas

| Tecla           | A��o                                     | Notas                                                                      |
|-----------------|------------------------------------------|----------------------------------------------------------------------------|
| `?`             | Expandir pasta                           |                                                                            |
| `?`             | Recolher pasta                           |                                                                            |
| `Enter`         | Expandir / Recolher pasta                | Quando o foco est� em uma pasta, expande/contrai.                          |
| `Shift+Insert`  | Criar nova pasta                         | Cria uma nova pasta no mesmo n�vel da pasta focada ou dentro dela, se n�o houver nenhuma pasta focada. |
| `Ctrl+Shift+I`  | Criar nova pasta                         | Atalho alternativo para criar uma nova pasta.                              |
| `Delete`        | Remover pasta                            | Marca a pasta selecionada para remo��o (revers�vel at� o salvamento).      |

### A��es em segredos

A coluna **Favoritos** indica se a a��o est� dispon�vel quando o cursor est� na pasta virtual Favoritos. A��es indispon�veis ficam ocultas na barra de comandos � n�o aparecem desabilitadas.

| Tecla    | A��o                          | Favoritos | Notas                                                                      |
|----------|-------------------------------|-----------|----------------------------------------------------------------------------|
| `Enter`  | Focar no painel de detalhes   | ?         | Comporta-se de forma similar ao `Tab` quando o foco est� em um segredo.    |
| `Insert` | Novo segredo                  | �         | Indispon�vel: Favoritos � somente leitura, sem pasta real associada.       |
| `^I`     | Novo segredo                  | �         | Atalho alternativo � mesma restri��o.                                      |
| `^E`     | Editar segredo                | ?         | Opera no segredo real, independente da vis�o atual.                        |
| `^D`     | Duplicar segredo              | �         | Indispon�vel: destino amb�guo. Navegar at� a pasta real para duplicar.     |
| `^M`     | Mover para outra pasta        | �         | Indispon�vel: mover a partir de pasta somente leitura n�o � permitido.     |
| `!?`     | Mover para cima na lista      | �         | Indispon�vel: a ordem na Favoritos reflete a �rvore real.                  |
| `!?`     | Mover para baixo na lista     | �         | Indispon�vel: idem.                                                        |
| `^S`     | Desfavoritar segredo          | ? (s� ?)  | Na Favoritos, o toggle s� remove o favorito � o segredo some da lista imediatamente. Em pasta real, alterna entre favoritar e desfavoritar normalmente. |
| `^R`     | Revelar primeiro campo sens�vel | ?       | Vis�vel apenas se o segredo tiver pelo menos um campo sens�vel.            |
| `^C`     | Copiar primeiro campo sens�vel  | ?       | Vis�vel apenas se o segredo tiver pelo menos um campo sens�vel.            |
| `Delete` | Excluir segredo               | �         | Indispon�vel: exclus�o direta a partir de pasta somente leitura n�o � permitida. |

#### ^D � Duplicar segredo

**Contexto:** foco na �rvore com cursor em um segredo, em pasta real. Indispon�vel na pasta virtual Favoritos � o destino do duplicado seria amb�guo para o usu�rio; a opera��o deve ser realizada navegando at� a pasta real do segredo.

**Comportamento:**
- Cria uma c�pia do segredo com todos os campos, valores e hist�rico de modelo id�nticos ao original.
- O novo segredo recebe automaticamente um nome �nico na mesma pasta com sufixo num�rico � ex: `Gmail (1)` se `Gmail` j� existe; `Gmail (2)` se `Gmail (1)` tamb�m j� existe.
- O novo segredo � posicionado imediatamente **ap�s o original** na lista da pasta.
- O novo segredo entra em estado `incluido`.
- O cursor da �rvore permanece no segredo original ap�s a duplica��o � o usu�rio pode navegar para o novo com `?`.
- A opera��o � instant�nea, sem di�logo de confirma��o.

**Feedback:** barra de mensagens exibe `? "[Nome original]" duplicado como "[Novo nome]"`.

**Refer�ncia:** [Fluxo 19 � Duplicar segredo](fluxos.md#fluxo-19--duplicar-segredo)

---

#### ^M � Mover para outra pasta

**Contexto:** foco na �rvore com cursor em um segredo. N�o dispon�vel na pasta virtual Favoritos (a pasta Favoritos � somente leitura � mover deve ocorrer na pasta real).

**Modo de sele��o inline:**

A �rvore entra em **modo mover** � um estado visual distinto:
- O segredo em foco recebe um indicador de "em movimento" (ex: �cone `?` ou destaque diferenciado em `accent.secondary`) e o cursor passa a navegar pela estrutura de pastas como destino.
- A barra de mensagens exibe `� Navegue at� a pasta de destino e pressione Enter para confirmar`.
- A barra de comandos muda para: `Enter Mover aqui � Esc Cancelar`.
- O usu�rio navega com `????` entre as pastas vis�veis.
- Pastas que resultariam em conflito de nome (j� cont�m um segredo com o mesmo nome) s�o marcadas visualmente como inv�lidas � o cursor pode passar por elas, mas `Enter` sobre elas exibe mensagem de erro na barra e aguarda nova sele��o.
- `Enter` sobre uma pasta v�lida confirma o movimento; o segredo � movido para a pasta de destino, o modo mover � encerrado e o cursor acompanha o segredo para a nova posi��o.
- `Esc` cancela o modo mover sem efeito colateral; o cursor retorna ao segredo original.

**Refer�ncia:** [Fluxo 25 � Mover segredo para outra pasta](fluxos.md#fluxo-25--mover-segredo-para-outra-pasta)

---

#### !? / !? � Reordenar segredo na lista

**Contexto:** foco na �rvore com cursor em um segredo, dentro de uma pasta real (n�o Favoritos).

**Comportamento:**
- `!?` desloca o segredo uma posi��o acima na lista da pasta atual; `!?` desloca uma posi��o abaixo.
- A opera��o � instant�nea e pode ser repetida sucessivamente.
- O cursor acompanha o segredo � ap�s o deslocamento, o cursor permanece sobre o mesmo segredo na nova posi��o.
- M�ltiplos deslocamentos antes de salvar resultam apenas no estado final � o hist�rico de movimentos intermedi�rios � descartado.
- A opera��o n�o tem feedback de mensagem � o deslocamento visual imediato � o feedback.

**Limites:**
- `!?` n�o tem efeito quando o segredo j� est� na primeira posi��o da pasta.
- `!?` n�o tem efeito quando o segredo j� est� na �ltima posi��o da pasta.
- Ambos ficam **ocultos na barra de comandos** e inativos quando o cursor est� na pasta virtual Favoritos.

**Indicador de modo ativo:** a barra de status/cabe�alho n�o precisa de indicador de modo para reordena��o direta � a opera��o � pontual e sem estado persistente.

**Refer�ncia:** [Fluxo 26 � Reordenar segredo dentro da mesma pasta](fluxos.md#fluxo-26--reordenar-segredo-dentro-da-mesma-pasta)

---

#### Barra de comandos contextualizada (�rvore, cursor em segredo � completa)

A tabela abaixo consolida todas as varia��es da barra de comandos para segredos na �rvore, incluindo os atalhos anteriores (`^R`, `^C`) e os novos (`^D`, `^M`, `!?`, `!?`).

| Condi��o | Barra de comandos |
|---|---|
| Pasta real � segredo sem campo sens�vel � posi��o intermedi�ria | `Enter Detalhes � ^E Editar � ^D Duplicar � ^M Mover � !? !? Reordenar � ^S Favoritar � Del Excluir � F1 Ajuda` |
| Pasta real � segredo sem campo sens�vel � primeiro da lista | `Enter Detalhes � ^E Editar � ^D Duplicar � ^M Mover � !? Mover para baixo � ^S Favoritar � Del Excluir � F1 Ajuda` |
| Pasta real � segredo sem campo sens�vel � �ltimo da lista | `Enter Detalhes � ^E Editar � ^D Duplicar � ^M Mover � !? Mover para cima � ^S Favoritar � Del Excluir � F1 Ajuda` |
| Pasta real � segredo com campo sens�vel � reveal mascarado | `Enter Detalhes � ^E Editar � ^D Duplicar � ^M Mover � !? !? Reordenar � ^S Favoritar � ^R Revelar � ^C Copiar � Del Excluir � F1 Ajuda` |
| Pasta real � segredo com campo sens�vel � reveal com dica | `Enter Detalhes � ^E Editar � ^D Duplicar � ^M Mover � !? !? Reordenar � ^S Favoritar � ^R Mostrar tudo � ^C Copiar � Del Excluir � F1 Ajuda` |
| Pasta real � segredo com campo sens�vel � reveal completo | `Enter Detalhes � ^E Editar � ^D Duplicar � ^M Mover � !? !? Reordenar � ^S Favoritar � ^R Ocultar � ^C Copiar � Del Excluir � F1 Ajuda` |
| Pasta Favoritos � segredo sem campo sens�vel | `Enter Detalhes � ^E Editar � ^S Desfavoritar � F1 Ajuda` |
| Pasta Favoritos � segredo com campo sens�vel � reveal mascarado | `Enter Detalhes � ^E Editar � ^S Desfavoritar � ^R Revelar � ^C Copiar � F1 Ajuda` |
| Pasta Favoritos � segredo com campo sens�vel � reveal com dica | `Enter Detalhes � ^E Editar � ^S Desfavoritar � ^R Mostrar tudo � ^C Copiar � F1 Ajuda` |
| Pasta Favoritos � segredo com campo sens�vel � reveal completo | `Enter Detalhes � ^E Editar � ^S Desfavoritar � ^R Ocultar � ^C Copiar � F1 Ajuda` |
| Modo mover ativo (^M pressionado) | `Enter Mover aqui � Esc Cancelar` |

> **Nota sobre tamanho da barra:** as entradas acima s�o o conjunto completo de a��es dispon�veis. Em terminais estreitos, a barra de comandos trunca � direita � as a��es mais priorit�rias devem aparecer primeiro. A ordem na barra segue a frequ�ncia de uso esperada.

---

#### ^R e ^C na �rvore � Atalhos de campo sens�vel

**Contexto:** foco na �rvore com cursor em um segredo que possui pelo menos um campo sens�vel.

**Campo alvo:** sempre o **primeiro campo sens�vel** do segredo (menor �ndice de posi��o no tipo).

**Visibilidade dos atalhos:**
- `^R` e `^C` aparecem na barra de comandos **somente** quando o cursor da �rvore est� em um segredo com pelo menos um campo sens�vel.
- Quando o cursor est� em uma pasta ou em um segredo sem campos sens�veis, os atalhos s�o omitidos da barra e n�o t�m efeito.

##### Comportamento de ^R na �rvore

- `^R` cicla o estado de reveal do primeiro campo sens�vel usando o **mesmo mecanismo de 3 estados do painel de detalhe**: mascarado ? dica (3 primeiros chars + `��`) ? completo ? mascarado.
- O painel direito � aberto (ou atualizado) automaticamente exibindo o segredo com o campo sens�vel j� no estado correspondente ao toque atual:
  - **1� toque:** painel exibe o campo sens�vel em estado de dica.
  - **2� toque:** painel exibe o campo sens�vel revelado completamente.
  - **3� toque:** campo re-mascarado; painel permanece aberto.
- As mesmas regras de re-mascaramento do painel se aplicam: trocar de segredo na �rvore ou timeout expirado re-mascara o campo silenciosamente.
- A barra de comandos reflete o estado atual do reveal (igual ao painel):
  - Mascarado: `^R Revelar`
  - Dica ativa: `^R Mostrar tudo`
  - Revelado: `^R Ocultar`

##### Comportamento de ^C na �rvore

- `^C` copia o valor **completo** do primeiro campo sens�vel para a clipboard � independentemente do estado de reveal atual (n�o � necess�rio revelar antes de copiar).
- Agenda limpeza autom�tica da clipboard (mesmo comportamento do `^C` no painel de detalhe).
- O painel direito � aberto (ou atualizado) automaticamente exibindo o segredo, mas o estado de reveal do campo **n�o muda** � a c�pia n�o desencadeia reveal.
- A barra de mensagens exibe confirma��o: `? [R�tulo do campo] copiado para a �rea de transfer�ncia`.

---

### Painel Direito: Detalhe do Segredo � Modo Leitura

**Contexto:** �rea de trabalho � Modo Cofre.
**Largura:** ~65% da �rea de trabalho.
**Responsabilidade:** Exibir o nome, o caminho de pastas, os campos e a observa��o do segredo selecionado na �rvore; permitir navega��o entre campos, c�pia de valores e reveal de campos sens�veis.

> Este documento especifica apenas o **modo leitura**. O modo edi��o de valores � especificado em [Modo Edi��o de Valores](#painel-direito-detalhe-do-segredo--modo-edi��o-de-valores). O modo edi��o de estrutura � especificado em [Modo Edi��o de Estrutura](#painel-direito-detalhe-do-segredo--modo-edi��o-de-estrutura).

---

#### Anatomia do painel

```
  Nome do Segredo                          Geral � Sites � Gmail ?
  --------------------------------------------------------------  �
  R�tulo do campo 1                                               �
  Valor do campo 1                                                �
                                                                  �
  R�tulo do campo 2                                               �
  Valor do campo 2                                                �
                                                                  �
  ??????????????????????????????????????????????????????????????  ?
  Texto da observa��o...
```

**Linha 1 � cabe�alho do segredo:**
- Esquerda: nome do segredo em `text.primary` **bold**
- Direita: breadcrumb com caminho completo de pastas em `text.secondary` � formato `Pasta � Subpasta � ...`; truncado � esquerda com `�` se n�o couber na linha. `?` aparece entre o nome e o breadcrumb quando o segredo � favoritado, em `accent.secondary`
- O breadcrumb mostra o caminho at� o segredo, excluindo o nome do segredo

**Linha 2 � separador:**
- `-` em `border.default` por toda a largura do painel (exceto a coluna reservada ao scroll)

**�rea de campos:**
- Cada campo ocupa dois segmentos: **r�tulo** (linha pr�pria, `text.secondary`) e **valor** (linha(s) seguinte(s), `text.primary`)
- Uma linha em branco separa campos consecutivos
- Campos sens�veis exibem o valor mascarado com `��������` em `text.secondary`; ao serem revelados, o valor real aparece em `text.primary`
- Campos com valor vazio: o r�tulo � exibido normalmente, a linha do valor fica em branco

**Separador da Observa��o:**
- `?` (U+254C) em `border.default`, ocupando toda a largura � omitido quando a Observa��o est� vazia
- A Observa��o n�o tem r�tulo; o separador e a posi��o final comunicam o que �

**Trilha de scroll:**
- �ltima coluna do painel reservada para `?`/`?`/`�` em `text.secondary`
- Reservada mesmo quando n�o h� scroll (evita deslocamento de conte�do ao ativar)

---

#### Wireframes

**Painel sem foco � segredo com campos variados:**

```
  Gmail ?                              Geral � Sites e Apps
  ----------------------------------------------------------
  URL
  https://accounts.google.com/login

  Usu�rio
  fulano@gmail.com

  Senha
  ����������

  Token 2FA
  ��������

  ??????????????????????????????????????????????????????????
  Conta pessoal principal � criada em 2018.
```

> Sem foco: nenhum campo destacado. O `?` aparece entre o nome e o breadcrumb quando o segredo � favoritado.

**Painel com foco � cursor em campo comum:**

```
  Gmail ?                              Geral � Sites e Apps
  ----------------------------------------------------------
  URL
  https://accounts.google.com/login

  Usu�rio                                                     ? special.highlight no bloco
  fulano@gmail.com

  Senha
  ����������

  ??????????????????????????????????????????????????????????
  Conta pessoal principal.
```

> O bloco inteiro do campo em foco (r�tulo + valor + linha em branco) recebe `special.highlight`. Barra de comandos (campo comum): `Enter Editar � ^S Favoritar � ^C Copiar � Tab �rvore � F1 Ajuda`

**Painel com foco � cursor em campo sens�vel:**

```
  Gmail ?                              Geral � Sites e Apps
  ----------------------------------------------------------
  URL
  https://accounts.google.com/login

  Usu�rio
  fulano@gmail.com

  Senha                                                       ? special.highlight no bloco
  ����������

  ??????????????????????????????????????????????????????????
  Conta pessoal principal.
```

> Barra de comandos (campo sens�vel mascarado): `Enter Editar � ^S Favoritar � ^C Copiar � ^R Revelar � Tab �rvore � F1 Ajuda`

**Campo sens�vel � estado de dica (1� `^R`):**

```
  Gmail ?                              Geral � Sites e Apps
  ----------------------------------------------------------
  URL
  https://accounts.google.com/login

  Usu�rio
  fulano@gmail.com

  Senha                                                       ? special.highlight
  min�������������                                            ? 3 chars revelados + �� mascarados

  ??????????????????????????????????????????????????????????
  Conta pessoal principal.
```

> Barra de comandos (dica ativa): `Enter Editar � ^S Favoritar � ^C Copiar � ^R Mostrar tudo � Tab �rvore � F1 Ajuda`

**Campo sens�vel � revelado completamente (2� `^R`):**

```
  Gmail ?                              Geral � Sites e Apps
  ----------------------------------------------------------
  URL
  https://accounts.google.com/login

  Usu�rio
  fulano@gmail.com

  Senha                                                       ? special.highlight
  minha-senha-secreta-123

  ??????????????????????????????????????????????????????????
  Conta pessoal principal.
```

> Barra de comandos (revelado): `Enter Editar � ^S Favoritar � ^C Copiar � ^R Ocultar � Tab �rvore � F1 Ajuda`

**Scroll ativo:**

```
  Gmail ?                              Geral � Sites e Apps ?
  ----------------------------------------------------------  �
  URL                                                         �
  https://accounts.google.com/login/v2/identifier?hl=pt-BR   �
                                                              �
  Usu�rio                                                     �
  fulano@gmail.com                                            �
                                                              ?
```

> Trilha de scroll: `?` quando h� conte�do acima, `?` quando h� abaixo, `�` na posi��o proporcional do thumb. A coluna da trilha � sempre reservada � o conte�do n�o se desloca ao ativar o scroll.

**Valor longo com quebra de linha:**

```
  Passos de acesso
  1. Acesse https://accounts.google.com
  2. Clique em "Fazer login com o Google"
  3. Confirme o dispositivo no app

```

> Valores multilinha recebem word-wrap; cada linha do valor ocupa a largura dispon�vel (exceto a coluna do scroll). O campo continua sendo tratado como uma unidade de foco � o bloco inteiro recebe highlight.

**Placeholders:**

```
  (sem segredo selecionado)
  -----------------------------------------------------------------


               Selecione um segredo para ver os detalhes


```

```
  (cofre vazio)
  -----------------------------------------------------------------


                           Cofre vazio


```

> Textos em `text.secondary` *italic*, centralizados na �rea de conte�do.

**Segredo sem Observa��o (separador omitido):**

```
  API Key � Stripe                            Geral � Financeiro
  --------------------------------------------------------------
  Servi�o
  Stripe

  Chave
  ����������

```

> Quando a Observa��o est� vazia, o separador `???` � omitido. N�o h� linha em branco extra no final.

**Breadcrumb truncado (caminho longo):**

```
  Gmail ?          � � Projetos � Cliente ABC � Acessos � Gmail
  --------------------------------------------------------------
```

> O breadcrumb � truncado � esquerda com `�` quando o caminho completo n�o cabe. O nome do segredo e o `?` nunca s�o truncados.

---

#### Mapa de teclas

| Tecla | Efeito | Condi��o |
|---|---|---|
| `?` / `?` | Move cursor para o campo anterior / pr�ximo | Painel com foco |
| `Home` | Vai ao primeiro campo | Painel com foco |
| `End` | Vai ao �ltimo campo (Observa��o, se n�o vazia) | Painel com foco |
| `PgUp` / `PgDn` | Scroll por p�gina (viewport - 1 linhas) | Painel com foco |
| `Enter` | Entra no modo edi��o do campo em foco | Painel com foco |
| `^S` | Favoritar / Desfavoritar segredo | Painel com foco |
| `^R` | 1� toque: revela dica (3 primeiros chars); 2� toque: revela valor completo; 3� toque: re-mascara | Painel com foco; campo sens�vel em foco |
| `^C` | Copiar valor do campo para clipboard; agenda limpeza da clipboard se campo sens�vel | Painel com foco; qualquer campo |
| `Tab` | Foco ? painel esquerdo (�rvore) | Painel com foco |

> `^R` n�o tem efeito quando o campo em foco � comum � a barra de comandos omite a a��o `Revelar` nesses casos.

---

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Nome do segredo | `text.primary` | **bold** |
| `?` favorito | `accent.secondary` | � |
| Breadcrumb de pasta | `text.secondary` | � |
| Separador `---` cabe�alho | `border.default` | � |
| R�tulo de campo | `text.secondary` | **bold** |
| Valor de campo comum | `text.primary` | � |
| Valor de campo � URL | `text.link` | � |
| Valor de campo sens�vel � mascarado `��������` | `text.secondary` | � |
| Valor de campo sens�vel � dica (`min����`) | `text.secondary` | � |
| Fundo do campo em foco | `special.highlight` | � |
| Separador `???` da Observa��o | `border.default` | � |
| Texto da Observa��o | `text.primary` | � |
| Placeholders | `text.secondary` | *italic* |
| `�` separador vertical � painel com foco | `border.focused` | � |
| `�` separador vertical � painel sem foco | `border.default` | � |
| `?`/`?`/`�` trilha de scroll | `text.secondary` | � |

---

#### Estados dos componentes

| Componente | Estado | Condi��o |
|---|---|---|
| Painel | placeholder "Selecione�" | Cofre tem segredos; nenhum segredo foi selecionado ainda na sess�o |
| Painel | placeholder "Cofre vazio" | Cofre sem nenhum segredo |
| Painel | segredo exibido (�ltimo selecionado) | Cursor da �rvore em pasta � painel mant�m o �ltimo segredo exibido |
| Painel | segredo exibido (atual) | Cursor da �rvore em segredo |
| Cursor de campo | ausente | Painel sem foco |
| Cursor de campo | `special.highlight` no bloco do campo | Painel com foco |
| `?` | vis�vel no cabe�alho, entre nome e breadcrumb | Segredo favoritado |
| `?` | ausente | Segredo n�o favoritado |
| Campo sens�vel | mascarado `��������` | Estado inicial ao exibir qualquer segredo |
| Campo sens�vel | dica (3 primeiros chars + `��`) | 1� `^R`; campo ainda em foco; timeout n�o expirou |
| Campo sens�vel | revelado (valor completo) | 2� `^R`; campo ainda em foco; timeout n�o expirou |
| Campo sens�vel revelado | re-mascarado | Timeout expirou; segredo diferente selecionado; foco saiu do campo |
| Separador `???` | vis�vel | Observa��o n�o vazia |
| Separador `???` | omitido | Observa��o vazia |
| Trilha de scroll | `?`/`?`/`�` ativos | Conte�do excede a �rea vis�vel |
| Trilha de scroll | coluna reservada, vazia | Conte�do cabe na �rea vis�vel |

---

#### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Painel recebe foco | Dica | `� Navegue com ?? e copie com ^C` |
| Campo sens�vel selecionado | Dica | `� ^R Revelar � ^C Copiar` |
| `^C` copia valor | Sucesso (5s) | `? [R�tulo do campo] copiado para a �rea de transfer�ncia` |

---

#### Eventos

| Evento | Efeito |
|---|---|
| Segredo selecionado na �rvore | Conte�do atualizado; campos revelados re-mascarados; cursor vai ao primeiro campo; `<�` aparece no separador |
| Painel recebe foco (`Tab`) | Cursor de campo aparece no campo anteriormente ativo, ou no primeiro campo se nunca focado |
| `?` / `?` | Cursor move para o campo anterior / pr�ximo; scroll autom�tico se necess�rio |
| `Home` / `End` | Cursor vai ao primeiro / �ltimo campo; scroll autom�tico |
| `PgUp` / `PgDn` | Scroll por p�gina |
| `Enter` | Entra no modo edi��o do campo em foco |
| `^S` | Segredo favoritado ? desfavoritado (ou vice-versa); `?` no cabe�alho do painel atualiza imediatamente; �rvore atualiza em segundo plano |
| `^R` em campo sens�vel mascarado | Campo entra em estado de dica (3 primeiros chars); barra muda para `^R Mostrar tudo � ^R Ocultar` |
| `^R` em campo sens�vel com dica | Campo revelado completamente; barra muda para `^R Ocultar` |
| `^R` em campo sens�vel revelado | Campo re-mascarado; barra volta para `^R Revelar` |
| `?` / `?` saindo de campo sens�vel revelado | Campo re-mascarado silenciosamente antes de mover o cursor |
| `Tab` com campo sens�vel revelado | Campo re-mascarado silenciosamente; foco transferido para a �rvore |
| Timeout de reveal expira | Campo re-mascarado silenciosamente; sem mensagem |
| Segredo diferente selecionado | Todos os campos revelados re-mascarados; cursor vai ao primeiro campo |

---

#### Comportamento

- **Cursor somente com foco** � o cursor de campo (highlight no bloco) aparece apenas quando o painel tem foco; sem foco, o conte�do � exibido sem destaque
- **Bloco de campo** � o campo em foco compreende: linha do r�tulo + linha(s) do valor + linha em branco de separa��o; todo o bloco recebe `special.highlight`
- **`Enter` entra no modo edi��o** � dispon�vel em qualquer campo com foco; aciona o modo edi��o de valores (especificado separadamente)
- **`^R` contextual** � dispon�vel apenas com campo sens�vel em foco; cicla entre tr�s estados: mascarado ? dica (3 primeiros chars) ? completo ? mascarado. N�o aparece na barra quando o campo em foco � comum
- **Re-mascaramento ao sair do campo** � ao mover o cursor para outro campo (`?`/`?`/`Home`/`End`) ou ao transferir o foco para a �rvore (`Tab`), qualquer campo sens�vel que estiver em estado de dica ou revelado � re-mascarado silenciosamente antes da movimenta��o
- **Campos sens�veis sempre iniciam mascarados** � incluindo segredos j� visitados anteriormente na sess�o
- **Reveal timeout** � configur�vel nas Configura��es; ao expirar, o campo � re-mascarado silenciosamente (sem mensagem na barra). Ao trocar de segredo, todos os reveals s�o cancelados imediatamente
- **URLs** � valores identificados como URL recebem `text.link`, diferenciados visualmente de texto puro
- **Observa��o � word-wrap** � o texto da Observa��o quebra na largura dispon�vel (exceto a coluna do scroll); pode ocupar m�ltiplas linhas; o painel inteiro � scroll�vel
- **Scroll** � a �ltima coluna do painel � sempre reservada para a trilha de scroll, mesmo quando n�o h� overflow � o conte�do n�o se desloca ao ativar o scroll (ver [DS � Scroll em di�logos](tui-design-system.md#scroll-em-di�logos))
- **`<�` e trilha de scroll s�o independentes** � `<�` aparece no separador vertical esquerdo e indica qual item da �rvore est� sendo detalhado; a trilha de scroll aparece na margem direita e reflete o scroll do conte�do do painel. Um n�o afeta o outro
- **Posi��o do cursor ao retornar o foco** � ao receber foco via `Tab` novamente, o cursor vai ao campo que estava ativo antes de o foco sair; se nunca focado, vai ao primeiro campo
- **Breadcrumb � truncamento** � o breadcrumb � truncado � esquerda com `�` se o caminho completo n�o couber; o nome do segredo e o `?` nunca s�o truncados

---

### Painel Direito: Detalhe do Segredo � Modo Edi��o de Valores

**Contexto:** �rea de trabalho � Modo Cofre. Ativado quando o usu�rio pressiona `Enter` sobre um campo no painel de detalhe em Modo Leitura.
**Largura:** ~65% da �rea de trabalho (igual ao Modo Leitura).
**Responsabilidade:** Permitir editar o valor de cada campo do segredo individualmente, com persist�ncia imediata por campo, sem estado global pendente.

> O modo edi��o de estrutura (renomear campos, adicionar/remover campos, reordenar) � especificado em [Modo Edi��o de Estrutura](#painel-direito-detalhe-do-segredo--modo-edi��o-de-estrutura).

---

#### Anatomia do modo

O Modo Edi��o de Valores � uma camada sobre o Modo Leitura. O layout do painel (cabe�alho, separador, campos, observa��o, scroll) permanece o mesmo � o que muda s�o:

1. **Indicador de modo** � `[editando]` em `accent.primary` **bold** aparece no cabe�alho, ap�s o nome do segredo e antes do `?`/breadcrumb
2. **Cursor de campo** � continua sendo `special.highlight` no bloco, como no Modo Leitura; o input se abre sobre o campo em foco
3. **Input inline** � quando um campo est� em edi��o, o valor � substitu�do por um campo de texto edit�vel na mesma posi��o; o input ocupa a largura total do painel (exceto a coluna de scroll)
4. **Barra de comandos** � muda conforme o estado: cursor de campo sem input aberto, ou input aberto

---

#### Anatomia do cabe�alho em edi��o

```
  Gmail [editando] ?                     Geral � Sites e Apps
  ----------------------------------------------------------
```

- Nome do segredo: `text.primary` **bold** (igual ao Modo Leitura)
- `[editando]`: `accent.primary` **bold**, separado do nome por um espa�o
- `?` e breadcrumb: inalterados

---

#### Wireframes

**Cursor no campo, sem input aberto (campo comum):**

```
  Gmail [editando] ?               Geral � Sites e Apps
  ------------------------------------------------------
  URL

  Usu�rio                                                ? special.highlight no bloco
  fulano@gmail.com

  Senha
  ����������

  ?????????????????????????????????????????????????????
  Conta pessoal principal.
```

> Barra: `Enter Editar campo � ^N Renomear � ^S Favoritar � Tab �rvore � Esc Sair da edi��o � F1 Ajuda`

**Input aberto � campo comum:**

```
  Gmail [editando] ?               Geral � Sites e Apps
  ------------------------------------------------------
  URL

  Usu�rio                                                ? special.highlight no bloco
  �fulano@gmail.com������������������������������������

  Senha
  ����������

  ?????????????????????????????????????????????????????
  Conta pessoal principal.
```

> `�` marca o fundo do input (`input.background`); `�` � o cursor de texto. O input substitui visualmente a linha do valor; o r�tulo permanece acima. Barra: `Enter Confirmar � Esc Cancelar campo � F1 Ajuda`

**Input aberto � campo sens�vel (revelado automaticamente):**

```
  Gmail [editando] ?               Geral � Sites e Apps
  ------------------------------------------------------
  URL

  Usu�rio
  fulano@gmail.com

  Senha                                                  ? special.highlight no bloco
  �minha-senha-secreta-123�����������������������������

  ?????????????????????????????????????????????????????
  Conta pessoal principal.
```

> Ao abrir o input de campo sens�vel, o valor � revelado automaticamente em texto claro dentro do input. Ao fechar o input (`Enter` ou `Esc`), o campo � re-mascarado imediatamente. Barra: `Enter Confirmar � Esc Cancelar campo � F1 Ajuda`

**Renomear segredo � input no cabe�alho (`^N`):**

```
  �Gmail�����������  [editando] ?        Geral � Sites e Apps
  ----------------------------------------------------------
  URL
  https://accounts.google.com/login

  Usu�rio
  fulano@gmail.com

  Senha
  ����������
```

> O input do nome abre inline no cabe�alho, substituindo o nome do segredo; `[editando]`, `?` e breadcrumb permanecem � direita. Nenhum campo da lista est� em foco enquanto o input do nome est� aberto. Barra: `Enter Confirmar nome � Esc Cancelar � F1 Ajuda`

**Valida��o � nome duplicado:**

```
  �Gmail�����������  [editando] ?        Geral � Sites e Apps
  ----------------------------------------------------------
```

> Barra de mensagens (erro): `? J� existe um segredo com esse nome nesta pasta` � input permanece aberto; o valor n�o � persistido.

**Cursor no campo, sem input � campo sens�vel:**

```
  Gmail [editando] ?               Geral � Sites e Apps
  ------------------------------------------------------
  URL

  Usu�rio
  fulano@gmail.com

  Senha                                                  ? special.highlight no bloco
  ����������

  ?????????????????????????????????????????????????????
  Conta pessoal principal.
```

> Campo sens�vel permanece mascarado enquanto n�o h� input aberto. Barra: `Enter Editar campo � ^N Renomear � ^S Favoritar � Tab �rvore � Esc Sair da edi��o � F1 Ajuda`

---

#### Mapa de teclas

**Com cursor de campo, sem input aberto:**

| Tecla | Efeito | Condi��o |
|---|---|---|
| `?` / `?` | Move cursor para o campo anterior / pr�ximo (sem abrir input) | � |
| `Home` / `End` | Cursor vai ao primeiro / �ltimo campo | � |
| `Enter` | Abre input inline no campo em foco | � |
| `^N` | Abre input inline no cabe�alho (renomear segredo) | � |
| `^S` | Favoritar / Desfavoritar segredo | � |
| `Tab` | Foco ? �rvore; sai do modo edi��o | � |
| `Esc` | Sai do modo edi��o; retorna ao Modo Leitura | � |

**Com input de campo aberto:**

| Tecla | Efeito |
|---|---|
| Texto / Backspace / Delete | Edita o valor no input |
| `Enter` | Persiste o valor; fecha o input; cursor permanece no campo |
| `?` | Persiste o valor implicitamente; fecha o input; move cursor para o campo anterior |
| `?` | Persiste o valor implicitamente; fecha o input; move cursor para o pr�ximo campo |
| `Esc` | Cancela; restaura o valor anterior; fecha o input; cursor permanece no campo |

**Com input do nome aberto (`^N`):**

| Tecla | Efeito |
|---|---|
| Texto / Backspace / Delete | Edita o nome no input |
| `Enter` | Valida e persiste o nome; fecha o input; retorna ao cursor de campo |
| `Esc` | Cancela; restaura o nome anterior; fecha o input |

> `Tab` com input de campo aberto: persiste o valor implicitamente, fecha o input, foco vai para a �rvore e sai do modo edi��o.
> `Tab` com input do nome aberto: cancela o nome (sem persistir), foco vai para a �rvore e sai do modo edi��o.

---

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Nome do segredo (cabe�alho) | `text.primary` | **bold** |
| `[editando]` | `accent.primary` | **bold** |
| `?` favorito | `accent.secondary` | � |
| Breadcrumb de pasta | `text.secondary` | � |
| Fundo do campo em foco (sem input) | `special.highlight` | � |
| Fundo do input aberto | `input.background` | � |
| Texto dentro do input | `text.primary` | � |
| Cursor de texto no input | terminal padr�o | � |
| R�tulo de campo | `text.secondary` | **bold** |
| Valor de campo comum (sem input) | `text.primary` | � |
| Valor de campo sens�vel mascarado (sem input) | `text.secondary` | � |
| Separador `---` cabe�alho | `border.default` | � |
| Separador `???` da Observa��o | `border.default` | � |

---

#### Estados dos componentes

| Componente | Estado | Condi��o |
|---|---|---|
| Indicador `[editando]` | vis�vel no cabe�alho | Modo edi��o de valores ativo |
| Cursor de campo | `special.highlight` no bloco | Sempre (modo edi��o tem foco impl�cito) |
| Input de campo | ausente | Cursor de campo sem edi��o ativa |
| Input de campo | aberto sobre a linha do valor | `Enter` pressionado sobre o campo |
| Campo sens�vel | mascarado `��������` | Input fechado |
| Campo sens�vel | revelado (texto claro no input) | Input aberto |
| Campo sens�vel | re-mascarado | Input fechado ap�s `Enter` ou `Esc` |
| Input do nome | ausente | `^N` n�o pressionado |
| Input do nome | aberto no cabe�alho | `^N` pressionado |
| Cursor de campo da lista | ausente | Input do nome aberto |

---

#### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Modo edi��o ativado | Dica | `� Enter para editar um campo � Esc para sair` |
| Campo confirmado (`Enter` ou `?`/`?` impl�cito) | Sucesso (3s) | `? [R�tulo do campo] salvo` |
| Nome duplicado ao confirmar | Erro | `? J� existe um segredo com esse nome nesta pasta` |
| Campo confirmado � campo sens�vel | Sucesso (3s) | `? [R�tulo do campo] salvo` |

---

#### Eventos

| Evento | Efeito |
|---|---|
| `Enter` no Modo Leitura sobre um campo | Modo edi��o de valores ativado; indicador `[editando]` aparece; input abre no campo em foco |
| `?` / `?` sem input aberto | Cursor de campo move; sem efeito colateral |
| `?` / `?` com input aberto | Valor persistido implicitamente; input fechado; cursor move para o campo anterior/pr�ximo |
| `Enter` com input aberto | Valor persistido; input fechado; cursor permanece no campo; mensagem de sucesso exibida |
| `Esc` com input aberto | Valor descartado; valor anterior restaurado; input fechado; cursor permanece no campo |
| `Tab` com input aberto | Valor persistido implicitamente; input fechado; foco vai para a �rvore; modo edi��o encerrado |
| `Tab` sem input aberto | Foco vai para a �rvore; modo edi��o encerrado |
| `Esc` sem input aberto | Modo edi��o encerrado; retorna ao Modo Leitura; indicador `[editando]` removido |
| `^N` | Input do nome abre no cabe�alho; cursor de campo da lista some |
| `Enter` com input do nome aberto | Nome validado; se v�lido: persistido, input fechado, cursor de campo da lista retorna; se inv�lido: mensagem de erro, input permanece |
| `Esc` com input do nome aberto | Nome descartado; nome anterior restaurado; input fechado; cursor de campo da lista retorna |
| `Tab` com input do nome aberto | Nome descartado (sem persistir); foco vai para a �rvore; modo edi��o encerrado |
| Campo sens�vel: input abre | Valor revelado automaticamente em texto claro no input |
| Campo sens�vel: input fecha | Campo re-mascarado imediatamente |
| `^Q` (sair da aplica��o) | Modo edi��o encerrado sem di�logo de confirma��o (persist�ncia imediata por campo elimina estado pendente) |

---

#### Comportamento

- **Persist�ncia imediata por campo** � cada campo � salvo ao confirmar (`Enter` ou movimento impl�cito com `?`/`?`/`Tab`); n�o h� estado de "edi��o pendente" global. `^Q` pode sair sem di�logo de confirma��o relacionado ao modo edi��o
- **Input inline** � o input abre na mesma posi��o da linha do valor, substituindo-a visualmente; o r�tulo permanece acima; a estrutura do painel n�o se desloca
- **Campo sens�vel revelado no input** � ao abrir o input de um campo sens�vel, o valor real � exibido em texto claro para permitir edi��o; ao fechar o input (por qualquer tecla), o campo � re-mascarado imediatamente, independentemente do resultado (confirmado ou cancelado)
- **`^R` indispon�vel no modo edi��o** � o ciclo de reveal do Modo Leitura n�o se aplica; o reveal ocorre automaticamente ao abrir o input
- **`^C` indispon�vel no modo edi��o** � c�pia de campo n�o est� dispon�vel enquanto o modo edi��o est� ativo
- **Navega��o sem abrir input** � `?`/`?`/`Home`/`End` movem o cursor entre campos sem abrir o input, igual ao Modo Leitura; o input s� abre com `Enter` expl�cito
- **Input do nome (`^N`) � independente do cursor de campo da lista** � enquanto o input do nome est� aberto, nenhum campo da lista est� em foco; ao fechar o input do nome, o cursor retorna ao campo que estava em foco antes de `^N`
- **Valida��o do nome** � o nome n�o pode ser vazio; n�o pode duplicar o nome de outro segredo na mesma pasta; a valida��o ocorre ao pressionar `Enter` no input do nome; erros mant�m o input aberto
- **Sair do modo edi��o** � `Esc` sem input aberto ou `Tab` encerram o modo edi��o; o indicador `[editando]` � removido; o painel retorna ao Modo Leitura com o mesmo campo em foco
- **Scroll** � o comportamento de scroll � id�ntico ao Modo Leitura; a coluna da trilha � sempre reservada

---

### Painel Direito: Detalhe do Segredo � Modo Edi��o de Estrutura

**Contexto:** �rea de trabalho � Modo Cofre. Ativado quando o usu�rio pressiona `^E` na �rvore, no painel em Modo Leitura ou no painel em Modo Edi��o de Valores.
**Largura:** ~65% da �rea de trabalho (igual ao Modo Leitura).
**Responsabilidade:** Permitir alterar a estrutura dos campos do segredo � renomear r�tulos, inserir campos, excluir campos e reordenar campos. Valores dos campos n�o s�o editados neste modo.

> Restri��es do dom�nio que este modo deve respeitar:
> - A **Observa��o** � n�o-delet�vel, n�o-renome�vel e n�o-mov�vel � ocupa sempre a �ltima posi��o e � exclu�da da navega��o do cursor neste modo
> - O **tipo** de um campo (`texto` / `texto_sensivel`) n�o pode ser alterado ap�s cria��o � apenas na inser��o
> - Nomes de campo **n�o t�m restri��o de unicidade**

---

#### Anatomia do modo

O Modo Edi��o de Estrutura � uma camada sobre o painel de detalhe. O layout permanece o mesmo (cabe�alho, separador, campos, observa��o, scroll). O que muda:

1. **Indicador de modo** � `[estrutura]` em `accent.primary` **bold** no cabe�alho, no lugar de `[editando]`
2. **Cursor de campo** � `special.highlight` no bloco do campo em foco, como nos outros modos; o cursor navega apenas entre campos edit�veis (Observa��o exclu�da)
3. **R�tulo em destaque** � o r�tulo do campo em foco recebe �nfase adicional (`text.primary` **bold**) para comunicar que � o alvo das a��es de estrutura
4. **Input inline de r�tulo** � quando um r�tulo est� em edi��o, o texto do r�tulo � substitu�do por um input na mesma linha
5. **Barra de comandos** � exibe as a��es do modo estrutura

---

#### Anatomia do cabe�alho em modo estrutura

```
  Gmail [estrutura] ?                    Geral � Sites e Apps
  ----------------------------------------------------------
```

- Nome do segredo: `text.primary` **bold**
- `[estrutura]`: `accent.primary` **bold**, separado do nome por um espa�o
- `?` e breadcrumb: inalterados

---

#### Wireframes

**Cursor no campo, sem input aberto:**

```
  Gmail [estrutura] ?              Geral � Sites e Apps
  ------------------------------------------------------
  URL
  https://accounts.google.com/login

  Usu�rio                                                ? special.highlight no bloco; r�tulo bold
  fulano@gmail.com

  Senha
  ����������

  ?????????????????????????????????????????????????????
  Conta pessoal principal.
```

> Barra: `Enter Renomear � !? Mover cima � !? Mover baixo � !Ins Inserir � !Del Excluir � Tab �rvore � Esc Sair � F1 Ajuda`
> Observa��o n�o tem cursor de foco � est� vis�vel mas exclu�da da navega��o do modo estrutura.

**Input de r�tulo aberto (`Enter`):**

```
  Gmail [estrutura] ?              Geral � Sites e Apps
  ------------------------------------------------------
  URL
  https://accounts.google.com/login

  �Usu�rio��������������������������������������������  ? input inline na linha do r�tulo
  fulano@gmail.com

  Senha
  ����������
```

> `�` marca o fundo do input (`input.background`); `�` � o cursor de texto. O valor do campo permanece vis�vel abaixo (leitura, sem altera��o). Barra: `Enter Confirmar � Esc Cancelar � F1 Ajuda`

**Input de r�tulo aberto � campo sens�vel:**

```
  Gmail [estrutura] ?              Geral � Sites e Apps
  ------------------------------------------------------
  URL
  https://accounts.google.com/login

  Usu�rio
  fulano@gmail.com

  �Senha����������������������������������������������  ? input do r�tulo
  ����������                                             ? valor permanece mascarado
```

> Campo sens�vel permanece mascarado no modo estrutura � n�o h� reveal autom�tico ao editar o r�tulo.

**Inser��o de novo campo (`!Ins`):**

```
  Gmail [estrutura] ?              Geral � Sites e Apps
  ------------------------------------------------------
  URL
  https://accounts.google.com/login

  Usu�rio                                                ? campo com foco antes de !Ins
  fulano@gmail.com

  ����������������������������������������  [texto] ^T  ? novo campo inserido abaixo; input vazio + badge de tipo
                                                         ? valor vazio (campo novo)
  Senha
  ����������
```

> O novo campo � inserido imediatamente abaixo do campo em foco e acima da Observa��o (se o foco estiver no �ltimo campo edit�vel, o novo campo � inserido entre ele e a Observa��o). O input do r�tulo abre automaticamente com o cursor. O badge `[texto]` indica o tipo atual; `^T` alterna entre `[texto]` e `[sens�vel]` enquanto o input est� aberto. Barra: `Enter Confirmar � ^T Tipo � Esc Cancelar � F1 Ajuda`

**Badge de tipo alternado para sens�vel:**

```
  ����������������������������������������  [sens�vel] ^T
```

> Ap�s `^T`, o badge muda para `[sens�vel]`. O campo ainda n�o tem r�tulo nem valor. `Enter` confirma nome e tipo.

**Reordenar campo (`!?` / `!?`):**

```
  Gmail [estrutura] ?              Geral � Sites e Apps
  ------------------------------------------------------
  URL
  https://accounts.google.com/login

  Senha                                                  ? campo movido para cima com !? (era abaixo de Usu�rio)
  ����������

  Usu�rio                                                ? special.highlight � campo em foco, foi deslocado para baixo
  fulano@gmail.com
```

> A reordena��o � imediata e vis�vel � o bloco do campo em foco se desloca e o cursor acompanha. O foco permanece no campo que foi movido.

---

#### Mapa de teclas

**Com cursor de campo, sem input aberto:**

| Tecla | Efeito | Condi��o |
|---|---|---|
| `?` / `?` | Move cursor para o campo anterior / pr�ximo | Apenas entre campos edit�veis (Observa��o exclu�da) |
| `Home` / `End` | Cursor vai ao primeiro / �ltimo campo edit�vel | � |
| `Enter` | Abre input inline no r�tulo do campo em foco | � |
| `!?` | Move o campo em foco uma posi��o acima | Sem efeito no primeiro campo edit�vel |
| `!?` | Move o campo em foco uma posi��o abaixo | Sem efeito no �ltimo campo edit�vel (antes da Observa��o) |
| `!Ins` | Insere novo campo abaixo do campo em foco; input do r�tulo abre automaticamente | � |
| `!Del` | Exclui o campo em foco imediatamente e irreversivelmente | � |
| `Tab` | Foco ? �rvore; sai do modo estrutura | � |
| `Esc` | Sai do modo estrutura; retorna ao Modo Leitura | � |
| `^E` | � (sem efeito � j� est� no modo estrutura) | � |

**Com input de r�tulo aberto (`Enter` ou via `!Ins`):**

| Tecla | Efeito |
|---|---|
| Texto / Backspace / Delete | Edita o nome do r�tulo |
| `^T` | Alterna o tipo do campo entre `texto` e `texto_sensivel` (apenas dispon�vel em inser��o � ver nota) |
| `Enter` | Valida e persiste o r�tulo (e tipo, se inser��o); fecha input; cursor permanece no campo |
| `Esc` | Cancela; restaura o r�tulo anterior (ou descarta inser��o); fecha input |
| `?` | Persiste implicitamente; fecha input; move cursor para o campo anterior |
| `?` | Persiste implicitamente; fecha input; move cursor para o pr�ximo campo |
| `Tab` | Persiste implicitamente; fecha input; foco vai para a �rvore; sai do modo estrutura |

> **`^T` (toggle de tipo) s� est� dispon�vel durante a inser��o** (`!Ins`). Em renomea��o de campo existente, o tipo � imut�vel � `^T` n�o tem efeito e o badge de tipo n�o � exibido.

---

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Nome do segredo (cabe�alho) | `text.primary` | **bold** |
| `[estrutura]` | `accent.primary` | **bold** |
| `?` favorito | `accent.secondary` | � |
| Breadcrumb de pasta | `text.secondary` | � |
| Fundo do campo em foco (sem input) | `special.highlight` | � |
| R�tulo do campo em foco (sem input) | `text.primary` | **bold** |
| R�tulo dos campos fora do foco | `text.secondary` | **bold** |
| Fundo do input de r�tulo | `input.background` | � |
| Texto dentro do input de r�tulo | `text.primary` | � |
| Cursor de texto no input | terminal padr�o | � |
| Badge de tipo `[texto]` / `[sens�vel]` | `text.secondary` | � |
| Valores dos campos (leitura) | inalterados do Modo Leitura | � |
| Separador `---` cabe�alho | `border.default` | � |
| Separador `???` da Observa��o | `border.default` | � |
| Observa��o (texto) | `text.secondary` | *italic* (diferenciada do modo leitura para comunicar inatividade) |

> A Observa��o recebe `text.secondary` *italic* no modo estrutura para sinalizar visualmente que est� exclu�da da navega��o e das a��es.

---

#### Estados dos componentes

| Componente | Estado | Condi��o |
|---|---|---|
| Indicador `[estrutura]` | vis�vel no cabe�alho | Modo estrutura ativo |
| Cursor de campo | `special.highlight` no bloco | Sempre (modo estrutura tem foco impl�cito) |
| Cursor de campo | ausente na Observa��o | Observa��o nunca recebe foco no modo estrutura |
| R�tulo do campo em foco | `text.primary` **bold** | � |
| Input de r�tulo | ausente | `Enter` n�o pressionado |
| Input de r�tulo | aberto sobre a linha do r�tulo | `Enter` pressionado; ou `!Ins` executado |
| Badge `[texto]` / `[sens�vel]` | vis�vel � direita do input | Apenas durante inser��o (`!Ins`) |
| Badge `[texto]` / `[sens�vel]` | ausente | Renomea��o de campo existente |
| Observa��o | vis�vel, n�o foc�vel, `text.secondary` *italic* | Sempre no modo estrutura |
| Campo sens�vel | mascarado `��������` | Sempre no modo estrutura (sem reveal) |
| Campo rec�m-inserido | input do r�tulo aberto, vazio | Imediatamente ap�s `!Ins` |

---

#### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Modo estrutura ativado | Dica | `� Enter para renomear � !Ins inserir � !Del excluir � !?? mover` |
| R�tulo renomeado confirmado | Sucesso (3s) | `? Campo renomeado` |
| Campo inserido | Sucesso (3s) | `? Campo "[nome]" adicionado` |
| Campo exclu�do | Sucesso (3s) | `? Campo "[nome]" exclu�do` |
| R�tulo vazio ao confirmar | Erro | `? O nome do campo n�o pode ser vazio` |
| `!Del` no �nico campo edit�vel | Erro | `? O segredo deve ter pelo menos um campo` |
| `!?` no primeiro campo | � | Sem mensagem � a��o sem efeito silenciosa |
| `!?` no �ltimo campo edit�vel | � | Sem mensagem � a��o sem efeito silenciosa |

---

#### Eventos

| Evento | Efeito |
|---|---|
| `^E` no Modo Leitura | Modo estrutura ativado; indicador `[estrutura]` aparece; cursor vai ao primeiro campo edit�vel |
| `^E` no Modo Edi��o de Valores | Modo valores encerrado (sem persist�ncia pendente � imediata); modo estrutura ativado |
| `^E` na �rvore | Painel recebe foco; modo estrutura ativado; cursor vai ao primeiro campo edit�vel |
| `?` / `?` sem input aberto | Cursor move entre campos edit�veis (Observa��o ignorada) |
| `Enter` sem input aberto | Input do r�tulo abre no campo em foco |
| `Enter` com input aberto | R�tulo validado; se v�lido: persistido, input fechado, cursor permanece; se inv�lido (vazio): mensagem de erro, input permanece |
| `Esc` com input aberto (renomea��o) | R�tulo descartado; r�tulo anterior restaurado; input fechado |
| `Esc` com input aberto (inser��o) | Campo rec�m-inserido descartado; cursor retorna ao campo que estava em foco antes de `!Ins` |
| `!?` | Campo em foco sobe uma posi��o; cursor acompanha; persistido imediatamente |
| `!?` | Campo em foco desce uma posi��o; cursor acompanha; persistido imediatamente |
| `!Ins` | Novo campo inserido abaixo do campo em foco (tipo `texto`); input do r�tulo abre automaticamente com cursor; badge `[texto]` vis�vel |
| `^T` com input de inser��o aberto | Tipo alterna entre `texto` e `texto_sensivel`; badge atualiza imediatamente |
| `Enter` com input de inser��o | R�tulo e tipo confirmados; campo inserido definitivamente; input fechado; cursor no novo campo |
| `!Del` | Campo em foco exclu�do imediatamente; cursor vai ao campo seguinte (ou anterior se era o �ltimo edit�vel) |
| `Esc` sem input aberto | Modo estrutura encerrado; retorna ao Modo Leitura; indicador `[estrutura]` removido |
| `Tab` sem input aberto | Foco vai para a �rvore; modo estrutura encerrado |
| `Tab` com input aberto | R�tulo persistido implicitamente; input fechado; foco vai para a �rvore; modo encerrado |
| `^Q` | Sa�da da aplica��o; persiste o que j� foi confirmado (imediato por opera��o) |

---

#### Comportamento

- **Persist�ncia imediata por opera��o** � cada a��o confirmada (renomear, inserir, mover, excluir) persiste em mem�ria imediatamente; n�o h� um "cancelar tudo" ao sair do modo. `Esc` s� cancela o input atualmente aberto, n�o as opera��es j� confirmadas
- **Observa��o exclu�da da navega��o** � o cursor de campo nunca vai para a Observa��o no modo estrutura; `?`/`?`/`Home`/`End` ignoram a Observa��o; `!?` no �ltimo campo edit�vel n�o tem efeito (n�o pode ultrapassar a Observa��o)
- **Tipo imut�vel em campos existentes** � `^T` s� funciona durante a inser��o de novo campo (`!Ins`); o badge de tipo s� � exibido nesse contexto; em renomea��o, o tipo n�o � alter�vel e o badge n�o aparece
- **`!Del` � irrevers�vel** � a exclus�o ocorre imediatamente ao pressionar `!Del`, sem confirma��o; o campo e seu valor s�o descartados; se o segredo tiver apenas um campo edit�vel, a exclus�o � bloqueada com mensagem de erro
- **`!Del` move o cursor** � ap�s excluir, o cursor vai para o campo seguinte; se era o �ltimo campo edit�vel, vai para o anterior
- **Input inline de r�tulo** � o input substitui visualmente a linha do r�tulo; o valor do campo permanece vis�vel abaixo em modo leitura durante a edi��o do r�tulo (o modo estrutura n�o altera valores)
- **Campo sens�vel permanece mascarado** � no modo estrutura, campos sens�veis exibem `��������`; n�o h� reveal nem `^R`
- **Inser��o abaixo do foco, acima da Observa��o** � se o foco est� no �ltimo campo edit�vel, o novo campo � inserido imediatamente antes da Observa��o; se o foco est� em outro campo, � inserido imediatamente abaixo do campo em foco
- **Troca de modo** � `^E` no Modo Edi��o de Valores troca para o modo estrutura sem di�logo; a persist�ncia imediata do modo valores garante que n�o h� dado pendente a perder
- **Sair do modo** � `Esc` sem input aberto ou `Tab` encerram o modo estrutura; o indicador `[estrutura]` � removido; o painel retorna ao Modo Leitura
- **Scroll** � id�ntico ao Modo Leitura; a coluna da trilha � sempre reservada

---

## Telas

### Boas-vindas

**Trigger:** Aplica��o inicia sem cofre aberto, ou ap�s fechar/bloquear cofre.  
**Intera��o:** Nenhuma � tela est�tica. Toda a��o dispon�vel via barra de comandos.

**Wireframe (�rea de trabalho � terminal 80 � 24):**

```
                                                                                
                                                                                
                                                                                
                   ___    __        ___ __                                      
                  /   |  / /_  ____/ (_) /___  ______ ___                       
                 / /| | / __ \/ __  / / __/ / / / __ `__ \                     
                / ___ |/ /_/ / /_/ / / /_/ /_/ / / / / / /                     
               /_/  |_/_.___/\__,_/_/\__/\__,_/_/ /_/ /_/                      
                                                                                
                             v0.1.0                                             
                                                                                
                                                                                
```

> Logo e vers�o centralizados via `lipgloss.Place()`. As linhas do logo recebem as cores do [DS � Gradiente do logo](tui-design-system.md#gradiente-do-logo) � n�o represent�vel neste wireframe monocrom�tico.

#### Tokens

| Elemento | Token | Atributo |
|---|---|---|
| Logo (linhas 1�5) | DS � [Gradiente do logo](tui-design-system.md#gradiente-do-logo) � por linha | � |
| Vers�o (ex: `v0.1.0`) | `text.secondary` | � |

> As cores do logo n�o s�o tokens nomeados � s�o os valores hexadecimais da tabela de gradiente do DS, aplicados por linha conforme o tema ativo.

#### Estados dos componentes

| Componente | Estado | Condi��o |
|---|---|---|
| Logo + vers�o | vis�vel, centralizado | Tela ativa |
| Cabe�alho | sem abas | Nenhum cofre aberto � ver [Cabe�alho � Sem cofre](#sem-cofre-boas-vindas) |

#### Mensagens

| Contexto | Tipo | Texto |
|---|---|---|
| Tela entra em exibi��o | Dica de uso | `� Abra ou crie um cofre para come�ar` |

#### Eventos

| Evento | Efeito |
|---|---|
| Aplica��o inicia sem cofre | Modo boas-vindas exibido |
| Cofre fechado | Tela boas-vindas exibida |
| Cofre bloqueado | Tela boas-vindas exibida; arquivo permanece em disco, requer nova autentica��o |
| Terminal redimensionado | Logo e vers�o recentralizados |

#### Comportamento

- Logo e vers�o centralizados horizontal e verticalmente na �rea de trabalho via `lipgloss.Place()`
- As cores do logo acompanham o tema ativo � mudam instantaneamente com `F12`
- O cabe�alho n�o exibe abas neste modo (ver [Cabe�alho � Sem cofre](#sem-cofre-boas-vindas))
- **Vers�o din�mica** � o texto exibido vem da string injetada em tempo de build via `-ldflags "-X main.version=$(git describe --tags --always)"`. Em builds locais sem tag, exibe `dev`. O valor **nunca** � hardcoded no fonte

---

<!-- SE��ES FUTURAS � a preencher pela equipe -->

<!--
## Telas (continua��o)

### Modo Cofre
### Modo Modelos
### Modo Configura��es

## Componentes

### Painel Direito: Detalhe do Modelo

## Fluxos Visuais

### Criar cofre
### Abrir cofre
### Salvar cofre
### Bloquear cofre
### Alterar senha mestra
### Criar segredo
### Editar segredo
### Excluir segredo
### Buscar segredo
### Exportar cofre
### Importar cofre
-->
