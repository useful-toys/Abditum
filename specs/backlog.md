# Abditum — Backlog

Ideias e requisitos futuros sem compromisso de implementação. Quando uma entrada for promovida a spec, move-se para `features/` e remove-se daqui.

---

## Segurança

- **Senha de pânico (duress password)**: digitar uma senha alternativa abre um cofre isca vazio ou com dados falsos, sem revelar a existência do cofre real. Protege contra coerção física. Referência: VeraCrypt hidden volume. Complexidade: alta — exige duplicação da estrutura de chave derivada e cuidado para que o tempo de resposta seja indistinguível entre senha real e senha de pânico.

- **Segundo fator — YubiKey / FIDO2**: exigir token hardware além da senha mestra. A chave final seria `KDF(senha) XOR segredo_do_hardware`. Ver `security-hardening.md`.

- **Modo somente leitura**: abrir o cofre sem possibilidade de salvar. Útil em máquinas não confiáveis. O usuário consulta mas não arrisca modificações acidentais.

---

## Gestão de dados

- **Validade de credenciais**: campo de data de validade por item com alerta visual para credenciais próximas do vencimento. Relevante para senhas corporativas com política de expiração.

- **Exportação/migração**: exportar o cofre decifrado para CSV ou JSON para migração para outra ferramenta. Requer decisão deliberada: exportar expõe todos os dados em texto claro — deve exigir confirmação explícita e aviso.

- **Gerador de senhas**: gerar senhas aleatórias com parâmetros configuráveis (tamanho, caracteres especiais, legibilidade). Resultado disponível para preencher atributo ou copiar diretamente.

- **Histórico de versões de itens**: manter versões anteriores de atributos alterados, permitindo recuperar um valor antigo.

---

## UX

- **Modo somente leitura**: abrir o cofre sem possibilidade de salvar. Útil em máquinas não confiáveis.

- **Busca com diacríticos**: normalizar Unicode (NFC/NFD) para que "cafe" encontre "café". Requer folding de diacríticos.

- **Busca parcial no meio da palavra**: "gmail" encontrar "mygmail-old". Requer mudança na estratégia de indexação (atualmente prefixo/substring do início).

- **Histórico de alterações por item**: manter registro das versões anteriores de atributos, com timestamp e indicador do que mudou. Permite recuperar um valor antigo (ex: senha antes de uma troca). Requer decisão sobre armazenamento: o histórico fica dentro do próprio `.abt` (aumenta o arquivo) ou em arquivo separado. Relacionado a "Histórico de versões de itens" já listado em Gestão de dados.

- **Itens recentes**: destacar itens visualizados ou editados mais recentemente, facilitando retomar o trabalho sem buscar. Requer armazenar timestamp de último acesso por item — decisão a tomar: gravar no cofre (altera o arquivo a cada visualização) ou manter apenas em memória (perde ao fechar).

- **Detecção de conflito de arquivo**: ao salvar, detectar se o `.abt` foi modificado em disco desde a abertura (relevante para cofres em pasta de nuvem sincronizada por dois dispositivos). Exibir alerta em vez de sobrescrever silenciosamente.

---

## Plataformas

- **Extensão de navegador**: preencher automaticamente credenciais em sites. Exige comunicação entre extensão e aplicação desktop (Native Messaging API).

- **Aplicativo mobile**: versão iOS/Android. Requer decisão sobre compartilhamento de formato de arquivo e estratégia de sincronização.
