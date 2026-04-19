## Context

`internal/crypto` concentra primitivas de criptografia, derivação de chave e limpeza de memória. A medição atual de cobertura do pacote está em 74.1%, com lacunas concentradas em ramos de validação e erro de `aead.go` — em especial `SealWithAAD` (0.0%) — e em wrappers específicos de plataforma como `mlock_windows.go`.

A mudança precisa aumentar a confiança nesses caminhos sem alterar a API pública do pacote nem o comportamento esperado em produção. Como parte das lacunas envolve helpers não exportados e falhas difíceis de provocar de forma determinística, o desenho dos testes precisa equilibrar cobertura, isolamento e estabilidade.

## Goals / Non-Goals

**Goals:**
- Cobrir ramos críticos ainda não exercitados em `Encrypt`, `Decrypt`, `SealWithAAD`, `EncryptWithAAD` e `DecryptWithAAD`.
- Validar de forma determinística cenários negativos relevantes, como chave inválida, nonce inválido e autenticação falha.
- Aumentar a cobertura total de `internal/crypto` para um patamar verificável, mantendo a suíte estável e sem efeitos colaterais externos.
- Cobrir helpers internos ou específicos de plataforma quando eles forem responsáveis por lacunas materiais de cobertura.
- Privilegiar testes sem mocks, recorrendo a doubles apenas quando forem indispensáveis ou quando a alternativa sem mock for complexa, deselegante ou contra boas práticas.

**Non-Goals:**
- Alterar contratos públicos do pacote `internal/crypto`.
- Trocar algoritmos criptográficos, parâmetros de segurança ou dependências de produção.
- Introduzir testes de benchmark ou performance como critério de aceite desta mudança.

## Decisions

1. **Expandir a cobertura primeiro com testes de caixa-preta das APIs públicas.**  
   A maior parte das lacunas está em validação de parâmetros e caminhos de erro acessíveis pelas funções exportadas. Isso mantém os testes alinhados com o comportamento público e reduz acoplamento com detalhes internos.  
   **Alternativa considerada:** testar apenas via fluxos integrados já existentes. Foi descartada porque deixa ramos de erro específicos sem cobertura direta.

2. **Adicionar testes em pacote interno quando a lacuna estiver em helpers não exportados.**  
   Para funções como `mlock` e `munlock`, testes em `package crypto` permitem exercitar o comportamento específico da plataforma sem expor novos símbolos públicos.  
   **Alternativa considerada:** exportar helpers apenas para teste. Foi descartada porque alteraria a superfície pública sem necessidade funcional.

3. **Evitar mocks e priorizar testes determinísticos sobre componentes reais.**  
   A implementação deve primeiro buscar cenários reproduzíveis com APIs públicas, helpers internos e entradas controladas. Doubles de teste só entram quando o ramo desejado não puder ser exercitado de forma confiável sem eles, ou quando a alternativa sem mock introduzir complexidade incidental, deselegância ou quebra de boas práticas.  
   **Alternativa considerada:** usar mocks amplamente para forçar falhas de `crypto/rand` ou do sistema operacional. Foi descartada porque aumentaria o acoplamento aos detalhes internos e reduziria a confiança nos fluxos reais.

4. **Introduzir seams privados somente para falhas que não possam ser reproduzidas de forma confiável sem degradar o design.**  
   Se algum ramo depender de falhas do sistema operacional ou de `crypto/rand`, a implementação poderá usar variáveis privadas redirecionáveis em teste, preservando a API pública e o comportamento padrão. Esse recurso deve ser tratado como último recurso, preferindo sempre exercícios reais do código.  
   **Alternativa considerada:** aceitar cobertura parcial para esses ramos. Foi descartada porque manteria pontos cegos justamente em caminhos de tratamento de erro.

5. **Usar cobertura do próprio Go como critério de regressão do pacote.**  
   A validação da mudança ficará ancorada em `go test -cover ./internal/crypto`, com uma meta explícita para evitar regressão silenciosa futura.  
   **Alternativa considerada:** usar apenas contagem de novos testes. Foi descartada porque quantidade de testes não garante cobertura dos ramos mais sensíveis.

## Risks / Trade-offs

- **[Testes frágeis por dependência de ambiente]** → Priorizar cenários determinísticos e limitar seams a pontos privados bem definidos.
- **[Uso excessivo de mocks gerar testes artificiais]** → Exigir justificativa explícita antes de introduzir doubles e preferir exercícios com comportamento real do pacote.
- **[Acoplamento excessivo a detalhes internos]** → Preferir testes de API pública e usar testes internos só onde helpers não exportados bloquearem cobertura relevante.
- **[Meta de cobertura sem melhorar valor real]** → Direcionar os novos testes para ramos críticos identificados pelo relatório de cobertura, não apenas para linhas fáceis.
- **[Diferenças entre plataformas]** → Isolar testes específicos por build tags ou por arquivo de plataforma para manter a suíte consistente em cada alvo.
