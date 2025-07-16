# Contribuindo para o Forge Framework

Obrigado por considerar contribuir para o Forge! Este documento fornece diretrizes para contribuições.

## 🚀 Como Contribuir

### Reportando Bugs

1. Verifique se o bug já foi reportado nas [Issues](https://github.com/joaofelipeuai/forge/issues)
2. Se não encontrar, crie uma nova issue com:
   - Descrição clara do problema
   - Passos para reproduzir
   - Comportamento esperado vs atual
   - Versão do Go e do Forge
   - Sistema operacional

### Sugerindo Melhorias

1. Abra uma issue com a tag "enhancement"
2. Descreva claramente a melhoria proposta
3. Explique por que seria útil para a comunidade
4. Forneça exemplos de uso se possível

### Enviando Pull Requests

1. Fork o repositório
2. Crie uma branch para sua feature:
   ```bash
   git checkout -b feature/nome-da-feature
   ```
3. Faça suas alterações seguindo as diretrizes de código
4. Adicione testes para novas funcionalidades
5. Execute os testes:
   ```bash
   go test ./...
   ```
6. Execute o linter:
   ```bash
   golangci-lint run
   ```
7. Commit suas mudanças:
   ```bash
   git commit -m "feat: adiciona nova funcionalidade"
   ```
8. Push para sua branch:
   ```bash
   git push origin feature/nome-da-feature
   ```
9. Abra um Pull Request

## 📝 Diretrizes de Código

### Estilo de Código

- Siga as convenções padrão do Go
- Use `gofmt` para formatação
- Mantenha funções pequenas e focadas
- Adicione comentários para código complexo
- Use nomes descritivos para variáveis e funções

### Testes

- Escreva testes para todas as novas funcionalidades
- Mantenha cobertura de testes acima de 80%
- Use testes de benchmark para funcionalidades críticas
- Testes devem ser independentes e determinísticos

### Documentação

- Documente todas as funções e tipos públicos
- Atualize o README.md se necessário
- Adicione exemplos para novas funcionalidades
- Use comentários claros e concisos

### Commits

Use o padrão de commits convencionais:

- `feat:` para novas funcionalidades
- `fix:` para correções de bugs
- `docs:` para mudanças na documentação
- `test:` para adição ou modificação de testes
- `refactor:` para refatoração de código
- `perf:` para melhorias de performance

## 🧪 Executando Testes

```bash
# Executar todos os testes
go test ./...

# Executar testes com cobertura
go test -cover ./...

# Executar benchmarks
go test -bench=. ./...

# Executar testes específicos
go test -run TestNomeDaFuncao ./...
```

## 📋 Checklist para Pull Requests

- [ ] Código segue as diretrizes de estilo
- [ ] Testes foram adicionados/atualizados
- [ ] Todos os testes passam
- [ ] Documentação foi atualizada
- [ ] Commit messages seguem o padrão
- [ ] Branch está atualizada com main
- [ ] Não há conflitos de merge

## 🤝 Código de Conduta

- Seja respeitoso e inclusivo
- Aceite feedback construtivo
- Foque no que é melhor para a comunidade
- Mantenha discussões técnicas e profissionais

## 💡 Ideias para Contribuições

### Funcionalidades Desejadas

- [ ] WebSocket support
- [ ] Template engine integration
- [ ] File upload middleware
- [ ] JWT authentication middleware
- [ ] Database integration helpers
- [ ] Metrics e monitoring
- [ ] GraphQL support
- [ ] Hot reload para desenvolvimento

### Melhorias de Performance

- [ ] Connection pooling
- [ ] Response caching
- [ ] Compression middleware
- [ ] Static file serving otimizado

### Documentação

- [ ] Mais exemplos de uso
- [ ] Tutoriais passo a passo
- [ ] Comparações com outros frameworks
- [ ] Guias de migração

## 📞 Contato

- GitHub Issues: Para bugs e sugestões
- Discussions: Para perguntas gerais
- Email: joaofelipeaps@gmail.com

Obrigado por contribuir! 🚀
