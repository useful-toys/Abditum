# Regras Gerais — Abditum

Regras transversais que se aplicam a múltiplas features e não se encaixam naturalmente em um único arquivo de cenários. Sempre que possível, as regras são expressas como cenários nas features correspondentes; aqui ficam apenas as que atravessam domínios.

---

## Unicidade e identidade

- A aplicação **não impõe unicidade de nomes** para segredos, pastas ou modelos de segredo. Nomes repetidos são permitidos.
- Cada segredo, pasta e modelo de segredo possui uma **identidade única** (NanoID de 6 caracteres) independente do nome.
- O nome é apenas um atributo editável — nunca participa da identidade do elemento.

## Portabilidade

- O cofre é um **arquivo único autossuficiente**: modelos de segredo, configurações e hierarquia ficam armazenados dentro dele.
- A aplicação **não persiste dados fora do arquivo do cofre**, exceto artefatos transitórios (`.abditum.tmp`) e backups (`.abditum.bak`) explicitamente previstos.
- Não há arquivos de configuração externos no sistema operacional.

## Privacidade

- A aplicação não emite **nenhum log** (stdout/stderr) que contenha caminhos de arquivos de cofre, nomes de segredos ou valores de campos.

## Compatibilidade

- A aplicação de versão N deve ser capaz de **abrir cofres criados em qualquer versão anterior** do formato suportada.
- Ao abrir um cofre antigo, o payload descriptografado é migrado em memória para o modelo atual.
- O salvamento sempre grava no formato da versão corrente da aplicação.
- Arquivos com versão de formato superior à suportada pela aplicação devem falhar com **erro claro de incompatibilidade**.

## Observação em segredos

- Todo segredo possui um campo de **observação** implícito, opcional e de texto livre.
- A observação **não é declarada nos modelos** e não pode ser removida.
- A observação é tratada como **dado não sensível** e participa da busca.

## Campos de segredo

- O valor de um campo pode ser **string vazia** (campo existente, não preenchido). Não há distinção de estado entre preenchido e vazio.
- Campos do tipo **texto sensível nunca participam da busca**, independentemente do estado visual de ocultação ou exibição.
- **Não é possível alterar o tipo** de um campo de segredo existente. Para mudar o tipo, é necessário excluir o campo e criar um novo.

## Modelo como snapshot

- Segredos criados a partir de modelos **não mantêm vínculo por referência** com o modelo de origem.
- O nome do modelo é guardado apenas como registro histórico (snapshot do momento da criação).
- Alterar, renomear ou excluir o modelo posteriormente **não afeta segredos já criados**.

## Ordenação

- A ordem de exibição dos elementos (segredos, pastas, campos) é idêntica à **ordem de armazenamento no JSON**.
- Dentro de cada coleção, segredos aparecem primeiro, depois subpastas.

## Invariantes de estado

- Só pode existir **um cofre ativo** por vez.
- Um segredo só pode estar na raiz **ou** em uma pasta — nunca em ambos, nem em duas pastas ao mesmo tempo.
- Uma pasta só pode estar na raiz **ou** dentro de outra pasta — nunca em ambos, nem em duas pastas ao mesmo tempo.
- Um segredo **não pode estar** simultaneamente na hierarquia principal e na Lixeira.
- A Lixeira só materializa segredos excluídos reversivelmente.
- Ao salvar, segredos na Lixeira são **permanentemente excluídos**, sem possibilidade de recuperação.
- Pastas **não possuem soft delete** — sua exclusão sempre remove a pasta e promove os filhos.
