## Problemas identificados nos requisitos

1.  **Requisito de "Conhecimento zero" vs. Migração de formato:** O documento afirma "conhecimento zero" (a aplicação não possui meios de acessar ou recuperar seus dados sem a senha mestra), mas também exige "compatibilidade retroativa" (a aplicação deve ser capaz de abrir arquivos de cofre criados em qualquer versão anterior do formato suportada pela aplicação). Se o formato do arquivo for alterado em versões futuras, a aplicação precisará conter lógica para migrar dados de formatos antigos para o novo, o que pode ser interpretado como uma forma limitada de "conhecimento" sobre a estrutura interna dos dados.

2.  **Consistência da senha em operações:** O documento afirma que a senha é fornecida uma única vez ao abrir o cofre e é usada para todas as operações. No entanto, ao alterar a senha mestra, a nova senha passa a ser usada. É importante garantir que essa transição seja atômica e consistente, sem risco de usar a senha antiga em algumas operações e a nova em outras durante o processo de alteração.

3.  **"Salvar cofre em outro arquivo" vs. "O arquivo de destino não pode ser o mesmo arquivo atual do cofre":** Este requisito parece contraditório. Salvar como deveria permitir salvar no mesmo arquivo, possivelmente sobrescrevendo-o, após confirmação.

4.  **Duress Password e quais segredos/pastas serão visíveis na versão restrita:** É necessário detalhar como essa configuração será feita (interface, local de armazenamento da configuração, etc).

5. **Necessidade de detalhar como a Observação é tratada na busca:** Os requisitos indicam que campos sensíveis não participam da busca. É necessário explicitar se a "Observação", por ser um campo comum especial, participa da busca ou não. Atualmente, está implícito que participa, mas a confirmação explícita evitaria ambiguidades.

6. **Necessidade de detalhar se é possível criar pastas e segredos na lixeira:** Atualmente, os requisitos não mencionam a possibilidade de criar pastas e segredos diretamente na Lixeira. É importante definir se essa funcionalidade será permitida ou não, e como ela se encaixa no fluxo de restauração de segredos.

7. **Necessidade de detalhar o comportamento da aplicação ao tentar abrir um arquivo .abditum que não seja um cofre válido:** O requisito de "validar arquivo contra corrupção e senha incorreta" é vago quanto ao tratamento de arquivos que não são cofres válidos. É importante especificar se a aplicação deve exibir uma mensagem de erro genérica ou se deve tentar identificar o tipo de arquivo e fornecer uma mensagem mais informativa.

8. **Ações simultâneas:** Não há nenhuma menção sobre o tratamento de ações simultâneas (ex: dois usuários tentando editar o mesmo segredo ao mesmo tempo).

9. **Comportamento da aplicação em caso de falta de espaço em disco:** É necessário definir o comportamento da aplicação em caso de falta de espaço em disco durante a criação ou salvamento do cofre.

10. **Backup e arquivos temporários:** É importante garantir que os arquivos de backup e temporários sejam removidos em caso de falha na operação de salvamento, para evitar a persistência de dados sensíveis em disco.

11. **Comportamento da aplicação quando o usuário tenta criar um segredo a partir de um modelo inexistente:** É importante definir o comportamento da aplicação quando o usuário tenta criar um segredo a partir de um modelo que não existe mais (por exemplo, se o modelo foi excluído após a criação do segredo).