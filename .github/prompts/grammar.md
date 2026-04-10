---
name: Grammar Check
description: Revisa gramática, ortografia, estilo e formatação
version: 1.0
category: Documentation
usage: "Revise gramática, ortografia e estilo do arquivo [ARQUIVO]"
---

# Grammar Check

Prompt para revisão completa de gramática, ortografia, estilo e formatação de documentos Markdown.

## Quick Start

Para revisar um arquivo completo:
```
Revise gramática, ortografia e estilo do arquivo c:\git\Abditum-T2\golden\requisitos.md
```

## Variações por Escopo

### Revisão Completa (Padrão)

```
Faça uma revisão completa de gramática, ortografia e estilo do arquivo [ARQUIVO]:

**Escopo da Revisão:**

1. **Ortografia**
   - Erros de digitação e acentuação
   - Palavras mal escritas

2. **Gramática**
   - Concordância nominal e verbal
   - Pontuação (vírgulas, travessões, pontos-e-vírgula)
   - Preposições e artigos incorretos
   - Tempos verbais
   - Estrutura de frases

3. **Estilo**
   - Redundância e repetição desnecessária
   - Clareza e concisão
   - Consistência de tom e voz
   - Padronização de terminologia
   - Fluxo e legibilidade

4. **Formatação Markdown**
   - Links, destaques e listas corretos
   - Espaçamento e indentação
   - Títulos e hierarquia

**Deliverable:**
Implemente as correções diretamente no arquivo usando múltiplas substituições simultâneas. Agrupe as correções por categoria e liste as mudanças realizadas como um resumo final.

**Padrão de Resposta:**
- Corrija em silêncio (sem avisos ou explicações, a menos que haja ambiguidade)
- Se encontrar redundâncias, consolide ou remova conforme necessário
- Se encontrar termos inconsistentes, padronize
- Se encontrar pontuação faltante ou excessiva, corrija
```

### Apenas Ortografia

```
Revise APENAS ortografia e acentuação no arquivo [ARQUIVO]:
- Erros de digitação
- Acentuação incorreta
- Caracteres especiais
Implemente as correções diretamente no arquivo.
```

### Apenas Gramática

```
Revise APENAS gramática no arquivo [ARQUIVO]:
- Concordância nominal e verbal
- Pontuação
- Tempos verbais
- Estrutura de frases
Implemente as correções diretamente no arquivo.
```

### Apenas Estilo

```
Revise APENAS estilo no arquivo [ARQUIVO]:
- Redundância e clareza
- Concisão e legibilidade
- Consistência de terminologia
- Fluxo textual
Implemente as correções diretamente no arquivo.
```

### Detectar Sem Implementar

```
IDENTIFIQUE (sem implementar) erros de gramática, ortografia ou estilo no arquivo [ARQUIVO]:
- Liste os erros encontrados por categoria
- Inclua linha/seção onde cada erro foi encontrado
- Forneça a correção recomendada
```

## Opções Avançadas

### Com Restrições

```
Revise gramática, ortografia e estilo de [ARQUIVO]:

Mantenha:
- Todas as seções e estrutura
- Nomes técnicos (Abditum, TUI, GCM, etc.)
- Acrônimos (API, TOTP, etc.)
- Exemplos de código

Foco especial em:
- Coerência de terminologia entre seções
- Consistência de tempo verbal
- Clareza de explicações
```

### Múltiplos Arquivos

```
Revise gramática, ortografia e estilo dos arquivos:
- [ARQUIVO1]
- [ARQUIVO2]
- [ARQUIVO3]

Com foco em consistência entre eles.
```

## Tips

- **Um arquivo por vez**: Use em um único arquivo para melhor contexto
- **Preview**: Peça "IDENTIFIQUE" primeiro se quer ver erros antes de corrigir
- **Contexto**: Mencione domínio (técnico, formal, etc.) se relevante
- **Verificação**: Após correções, solicite "identifique redundâncias"

