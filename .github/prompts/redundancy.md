---
name: Redundancy Check
description: Identifica duplicidade e redundância de forma interativa
version: 1.0
category: Documentation
usage: "Verifique redundâncias em [ARQUIVO] (modo interativo)"
---

# Redundancy Check

Prompt interativo para identificar duplicidade e redundância em documentação, com sugestões iterativas para correção.

## Quick Start

Modo interativo (com aprovação para cada correção):
```
Verifique redundâncias em c:\git\Abditum-T2\golden\requisitos.md (modo interativo)
```

Modo detecção (sem implementar):
```
Detecte redundâncias em c:\git\Abditum-T2\golden\requisitos.md
```

## Modos de Operação

### Modo Interativo (Padrão)

```
MODO INTERATIVO: Identifique duplicidade e redundância em [ARQUIVO]

**Instruções:**
1. Analise o arquivo procurando por:
   - Conceitos repetidos em múltiplas seções
   - Texto literal idêntico ou muito similar
   - Regras transversais mencionadas em requisitos
   - Redundância de definições (Glossário vs Regras)
   - Padrões duplicados (ex: tratamento de erros, validações)

2. Para CADA redundância encontrada:
   - Apresente o PROBLEMA identificado
   - Mostre os TRECHOS envolvidos (com arquivo e linha)
   - Forneça 2-3 OPÇÕES de solução
   - Aguarde resposta do usuário

3. Após cada resposta do usuário:
   - Implemente a solução escolhida
   - Confirme a mudança realizada
   - Passe para a próxima redundância

4. Ao final:
   - Resuma as 5 mudanças mais impactantes
   - Pergunta: "Deseja fazer uma nova rodada de verificação? (sim/não)"
   - Se "sim", reinicia (encontrará redundâncias de segunda ordem)
   - Se "não", finaliza

**Restrições:**
- Não implemente nada sem aprovação
- Mantenha estrutura e semântica original
- Priorize redundâncias ÓBVIAS nas primeiras rodadas
```

### Modo Detecção (Sem Implementar)

```
DETECTAR (sem implementar) redundâncias em [ARQUIVO]:
- Liste TODAS as redundâncias encontradas
- Ordene por impacto (crítica/alta/média/baixa)
- Para cada uma, forneça 2 soluções possíveis
- Ao final, pergunta: "Deseja iniciar modo interativo? (sim/não)"
```

### Modo Automático (Confiança Alta)

```
CORRIGIR AUTOMATICAMENTE redundâncias óbvias em [ARQUIVO]:
- Identifique duplicações textuais LITERAIS (100% iguais)
- Implemente a consolidação mais óbvia
- Não pergunte para cada uma
- Ao final, liste as 5 mudanças feitas
- Pergunta: "Deseja modo interativo para outras? (sim/não)"
```

### Modo Específico

```
VERIFICAR [TIPO] redundâncias em [ARQUIVO]:
Exemplo: VERIFICAR "Observação" redundâncias em requisitos.md

- Procure APENAS por esse conceito no arquivo
- Apresente interativamente as variações
- Sugira unificação
```

### Modo Comparativo

```
COMPARAR redundâncias entre:
  - arquivo1.md
  - arquivo2.md
  - arquivo3.md

Identifique conceitos duplicados ENTRE os arquivos e sugira consolidações.
```

## Estrutura de Interação

Cada redundância é apresentada assim:

```
═══════════════════════════════════════════════════════════════════
🔴 REDUNDÂNCIA #[N]: [TÍTULO DO PROBLEMA]
═══════════════════════════════════════════════════════════════════

📍 LOCALIZAÇÃO:
   • Arquivo: [ARQUIVO]
   • Linhas: [L1-L2] e [L3-L4]

⚠️  PROBLEMA IDENTIFICADO:
   [Descrição clara do problema]

📜 OCORRÊNCIA 1:
   [Trecho de código/texto]
   
📜 OCORRÊNCIA 2:
   [Trecho de código/texto]

✅ OPÇÕES DE SOLUÇÃO:

   [A] [Descrição da opção A]
       Resultado: [como ficaria]
   
   [B] [Descrição da opção B]
       Resultado: [como ficaria]
   
   [C] [Descrição da opção C]
       Resultado: [como ficaria]

❓ QUE DEVO FAZER?
   Digite: A / B / C / PULAR / ou descreva outra solução
```

## Priorização

O prompt prioriza redundâncias nesta ordem:

1. **CRÍTICA** (correção obrigatória)
   - Texto 100% idêntico em 2+ lugares
   - Contradições
   - Requisito mencionando outro (deveria ser transversal)

2. **ALTA** (deve corrigir)
   - Conceitos duplicados com palavras diferentes
   - Regra em requisito + em transversal
   - Definições em glossário + em regras

3. **MÉDIA** (boa ideia)
   - Padrões similares (ex: 2 tratamentos equivalentes)
   - Exemplos duplicados

4. **BAIXA** (opcional)
   - Leves repetições de terminologia
   - Reforços de conceitos para clareza

## Rodadas Iterativas

Cada rodada remove camadas de redundância:

- **Rodada 1**: Óbvias (texto idêntico, conceitos duplicados)
- **Rodada 2**: Decorrências (referências cruzadas, simplificações)
- **Rodada 3**: Refinamento (padrões similares, organização)
- **Rodada 4+**: Ajustes finos (coesão, consistência)

Típico: 2-3 rodadas essenciais, opcional até 4ª rodada.

## Salvamento

Após cada sessão:
```
Deseja SALVAR as mudanças? (sim/não)
   - "não": Descarta tudo, reinicia análise
   - "sim": Mantém mudanças, oferece:
     ✓ Nova rodada de análise
     ✓ Sair e finalizar
```

## Tips

- **2ª rodada**: Após salvar, nova rodada encontrará redundâncias de 2ª ordem
- **Múltiplos arquivos**: Use modo comparativo para coesão entre documentos
- **Combinado**: Use após `grammar-style-check` para melhor qualidade

