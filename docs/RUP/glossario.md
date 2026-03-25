# Glossário — Abditum

| Item            | Detalhe                        |
|-----------------|--------------------------------|
| Projeto         | Abditum                       |
| Versão          | 1.1                           |
| Data            | 2026-03-25                     |
| Status          | Aprovado                       |

---

## Termos

### Bloqueio do Cofre
Interrupção do acesso ao conteúdo do cofre como medida de proteção, exigindo nova autenticação para retomar o acesso. Pode ser iniciado pelo próprio usuário ou ocorrer automaticamente após um período de inatividade.

### Campo de Segredo
Elemento individual dentro de um segredo que armazena um dado específico. Possui nome e classificação (dado comum ou dado sensível).

### Cofre
Repositório protegido que armazena os segredos, pastas, modelos de segredo e configurações do usuário. É completamente autossuficiente e portátil — todas as informações necessárias para seu uso estão contidas nele. Apenas um cofre pode estar em uso por vez.

### Conhecimento Zero (Zero Knowledge)
Princípio de segurança segundo o qual não existe nenhum meio de acessar ou recuperar os dados do cofre sem a senha mestra do usuário. O esquecimento da senha mestra resulta em perda total e irrecuperável dos dados.

### Dados Comuns
Informações não sensíveis armazenadas em um segredo, como nome de serviço, URL ou nome de usuário. Participam de buscas.

### Dados Sensíveis
Informações confidenciais armazenadas em um segredo, como senhas e chaves de API. São protegidas por padrão e não participam de buscas.

### Exclusão Reversível (Soft Delete)
Mecanismo pelo qual um segredo excluído não é imediatamente eliminado, mas movido para a Lixeira, de onde pode ser restaurado ao seu local e estado originais. A exclusão se torna irreversível (permanente) na próxima persistência definitiva do cofre.

### Favorito
Marcação que o usuário atribui a um segredo para indicá-lo como prioritário ou de uso frequente, conferindo-lhe destaque para acesso rápido sem alterar sua localização na hierarquia.

### Hierarquia do Cofre
Organização dos segredos em pastas e subpastas dentro do cofre. A hierarquia permite aninhamento livre de pastas, e cada segredo pertence a exatamente uma pasta ou à raiz do cofre.

### Lixeira
Agrupamento que reúne os segredos excluídos reversivelmente. Permite a restauração de segredos antes da próxima persistência definitiva, quando os segredos nela contidos são eliminados permanentemente.

### Modelo de Segredo
Estrutura reutilizável que define um conjunto de campos (com nome e classificação) usada como template para criar novos segredos de forma padronizada. Segredos criados a partir de um modelo não mantêm vínculo retroativo com ele — alterações no modelo afetam apenas criações futuras.

### Modelo de Segredo Personalizado
Modelo de segredo criado pelo próprio usuário para atender a necessidades específicas de formato, com campos definidos livremente.

### Modelo de Segredo Pré-definido
Modelo de segredo fornecido ao criar um novo cofre, com campos comuns para tipos populares de segredos. Exemplos: Login (URL, Username, Password), Cartão de Crédito (Número, Nome, Validade, CVV) e API Key (Nome da API, Chave). São editáveis e removíveis pelo usuário.

### Observação
Campo de texto livre presente em todo segredo, destinado a informações adicionais. É classificado como dado não sensível. Não pode ser removido do segredo.

### Pasta
Contêiner estrutural utilizado para agrupar e organizar segredos e outras subpastas na hierarquia do cofre. Pode conter segredos e subpastas em qualquer nível de aninhamento.

### Pasta Virtual
Agrupamento lógico que reúne segredos com base em características específicas (como favoritos ou segredos excluídos reversivelmente), sem alterar sua localização real na hierarquia.

### Raiz do Cofre
Nível estrutural mais alto da hierarquia, que contém os segredos e pastas que não estão aninhados em outras pastas.

### Segredo
Item individual dentro do cofre que armazena informações do usuário. É composto por um nome, campos (dados comuns e sensíveis), uma observação e atributos como favorito. Possui identidade própria independente do nome.

### Senha Mestra
Credencial única de acesso ao cofre. É a única forma de acessar o conteúdo protegido. Ao ser criada ou alterada, exige confirmação por dupla digitação. Irrecuperável em caso de esquecimento (Conhecimento Zero).

### Shoulder Surfing
Observação não autorizada de informações do usuário por terceiros presentes no mesmo ambiente físico.