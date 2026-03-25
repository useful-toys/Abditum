# Funcionalidade: Segurança e proteção de dados

**Como** usuário do Abditum,
**quero** que meus dados sejam protegidos por mecanismos de segurança robustos,
**para** garantir que minhas informações confidenciais estejam seguras contra acessos não autorizados e espionagem.

Esta feature documenta os comportamentos de segurança transversais: proteção criptográfica, controles de exposição de dados sensíveis e mitigações contra ameaças.

---

**Contexto:**

  **Dado** que a aplicação está em execução

---

## Criptografia e Conhecimento Zero

### **Regra:** A aplicação não possui meios de acessar dados sem a senha mestra

#### **Cenário:** Irrecuperabilidade da senha mestra

  **Dado** que um cofre foi criado com uma senha mestra
  **E** o usuário esqueceu a senha mestra
  **Então** a aplicação não oferece nenhum mecanismo de recuperação de senha
  **E** os dados do cofre estão permanentemente inacessíveis

### **Regra:** O cofre é protegido por AES-256-GCM com chave derivada via Argon2id

#### **Cenário:** Proteção contra força bruta

  **Dado** que um cofre existe em disco
  **Quando** um atacante tenta derivar a chave por força bruta
  **Então** cada tentativa consome no mínimo 256 MiB de memória e 3 iterações do Argon2id
  **E** o custo computacional torna ataques offline impraticáveis

---

## Proteção contra shoulder surfing

### **Regra:** A interface pode ser ocultada rapidamente por atalho

#### **Cenário:** Ocultar interface rapidamente

  **Dado** que o cofre está aberto e dados estão visíveis na tela
  **Quando** o usuário aciona o atalho de proteção contra shoulder surfing
  **Então** toda a interface é ocultada imediatamente

---

## Aviso de Conhecimento Zero na criação do cofre

### **Regra:** O usuário é alertado sobre a irrecuperabilidade ao criar um cofre

#### **Cenário:** Aviso de irrecuperabilidade ao criar cofre

  **Quando** o usuário cria um novo cofre
  **Então** a aplicação exibe aviso categórico de que o esquecimento da senha mestra resulta em perda total dos dados

---

## Privacidade de logs

### **Regra:** Nenhum log da aplicação contém dados sensíveis

#### **Cenário:** Ausência de caminhos e nomes em logs

  **Dado** que a aplicação está em execução
  **Então** nenhuma saída em stdout ou stderr contém caminhos de arquivos de cofre, nomes de segredos ou valores de campos

---

## Minimização de dados em memória

### **Regra:** Dados sensíveis são limpos ao bloquear ou fechar o cofre

#### **Cenário:** Limpeza de buffers ao bloquear cofre

  **Quando** o cofre é bloqueado
  **Então** a aplicação limpa os buffers controlados sempre que possível
  **E** a área de transferência é limpa

#### **Cenário:** Limpeza de buffers ao fechar aplicação

  **Quando** a aplicação é encerrada
  **Então** a aplicação limpa os buffers controlados sempre que possível
  **E** a área de transferência é limpa
