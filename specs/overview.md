# Abditum — Visão Geral

## Propósito

Abditum é uma aplicação desktop para gerenciamento de dados pessoais sensíveis. Permite armazenar, organizar e recuperar informações como credenciais de sites, dados bancários e outros registros privados, tudo protegido por criptografia forte em um único arquivo portátil.

## Objetivos

- Armazenar dados sensíveis com segurança máxima localmente
- Organizar dados em estrutura hierárquica (árvore) flexível e navegável
- Suportar modelos predefinidos (site, conta bancária, cartão) e modelos personalizados
- Funcionar completamente offline, sem dependência de serviços externos
- Gerar um único arquivo portátil e resistente

## Princípios

- **Segurança primeiro**: criptografia forte, zero transmissão de dados
- **Offline-only**: nenhum dado sai da máquina do usuário
- **Arquivo único**: todo o cofre em um único arquivo `.abt` portátil
- **Flexibilidade**: atributos livres por item, templates como conveniência, não obrigação
- **Simplicidade**: interface focada, sem funcionalidades que distraiam

## Expectativas do usuário sobre distribuição

O Abditum se propõe a ser uma ferramenta minimalista. Isso não é apenas uma escolha técnica — é uma necessidade do usuário que opta por ela. Quem escolhe armazenar seus dados mais sensíveis localmente está, em geral, rejeitando soluções complexas, dependentes de contas online ou de infraestrutura opaca.

Essa expectativa se estende à própria instalação e uso da ferramenta:

- **Sem instalador**: copiar o executável é suficiente para "instalar"
- **Sem runtime externo**: não exige Java, .NET, Node.js ou qualquer dependência prévia no sistema
- **Portável**: o executável pode ser carregado em um pen drive e executado em qualquer máquina compatível
- **Leve**: ocupa poucos megabytes em disco e pouca memória — não compete com recursos do sistema
- **Removível sem rastros**: deletar o executável e o arquivo `.abt` é suficiente para remover completamente a aplicação

Essas restrições moldaram diretamente a escolha de linguagem e ferramentas (ver `decisions/adr-001-linguagem.md`). Não são vaidades técnicas — são o reflexo do conceito da aplicação: uma ferramenta que o usuário controla completamente, que não instala serviços, não cria entradas de registro desnecessárias, não depende de nada além de si mesma.

## Features (v1)

| Feature | Arquivo |
|---------|---------|
| Gerenciamento do cofre (criar, abrir, fechar, alterar senha) | [features/01-cofre.md](features/01-cofre.md) |
| Navegação e organização da árvore | [features/02-arvore.md](features/02-arvore.md) |
| Gerenciamento de itens e atributos | [features/03-itens.md](features/03-itens.md) |
| Templates predefinidos e personalizados | [features/04-templates.md](features/04-templates.md) |
| Busca | [features/05-busca.md](features/05-busca.md) |
| Área de transferência segura | [features/06-clipboard.md](features/06-clipboard.md) |
| Favoritos | [features/07-favoritos.md](features/07-favoritos.md) |

## Fora do escopo (v1)

- Sincronização em nuvem
- Compartilhamento entre usuários
- Extensão de navegador
- Aplicativo mobile
- Gerador de senhas (v2)
- Histórico de versões de itens (v2)

## Hardening de segurança

Além da criptografia do cofre, existem medidas adicionais de proteção contra ataques físicos, de sistema operacional e de processo. Algumas são requisito v1 (escrita atômica, limpeza de clipboard, logs zero); as demais estão catalogadas com justificativa e fase de implementação em [security-hardening.md](security-hardening.md).

## Backlog

Ideias e requisitos futuros sem compromisso de implementação estão em [backlog.md](backlog.md).

## Plataformas alvo

- Windows
- macOS
- Linux
