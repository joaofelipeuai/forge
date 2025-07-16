# Forge Framework - Estrutura Recomendada

## 🏗️ Estrutura Atual vs. Recomendada

### ❌ Estrutura Atual (Problemática):
```
forge/
├── go.mod
├── forge.go              # Core misturado na raiz
├── jwt.go               # Módulos na raiz
├── websocket.go         # Módulos na raiz
├── upload.go            # Módulos na raiz
├── template.go          # Módulos na raiz
├── hotreload.go         # Módulos na raiz
├── demo.go              # Demo na raiz (problemático)
├── test_demo.go         # Teste na raiz (problemático)
├── forge_test.go        # OK - testes do core
├── *.exe                # Executáveis commitados (❌)
├── cmd/
│   └── demo_server.go
└── examples/
    ├── simple/
    ├── basic/
    └── advanced/
```

### ✅ Estrutura Recomendada (Go Standard Layout):
```
forge/
├── go.mod
├── README.md
├── LICENSE
├── CONTRIBUTING.md
├── .gitignore           # Para ignorar executáveis
│
├── pkg/                 # Código público reutilizável
│   └── forge/
│       ├── forge.go     # Core do framework
│       ├── context.go   # Context e métodos
│       ├── router.go    # Sistema de rotas
│       └── middleware/  # Middleware organizados
│           ├── auth/
│           │   ├── jwt.go
│           │   └── jwt_test.go
│           ├── upload/
│           │   ├── upload.go
│           │   └── upload_test.go
│           ├── websocket/
│           │   ├── websocket.go
│           │   └── websocket_test.go
│           ├── template/
│           │   ├── template.go
│           │   └── template_test.go
│           └── hotreload/
│               ├── hotreload.go
│               └── hotreload_test.go
│
├── internal/            # Código privado do projeto
│   ├── demo/
│   │   ├── demo.go      # Funções de demo
│   │   └── demo_test.go
│   └── utils/
│       └── helpers.go
│
├── cmd/                 # Aplicações executáveis
│   ├── forge-demo/
│   │   └── main.go      # Servidor de demonstração
│   └── forge-cli/
│       └── main.go      # CLI tool (futuro)
│
├── examples/            # Exemplos de uso
│   ├── simple/
│   │   ├── main.go
│   │   └── README.md
│   ├── basic/
│   │   ├── main.go
│   │   └── README.md
│   └── advanced/
│       ├── main.go
│       ├── templates/
│       ├── uploads/
│       └── README.md
│
├── docs/                # Documentação
│   ├── getting-started.md
│   ├── middleware.md
│   └── api-reference.md
│
├── scripts/             # Scripts de build/deploy
│   ├── build.sh
│   └── test.sh
│
└── test/                # Testes de integração
    ├── integration/
    └── fixtures/
```

## 🎯 Vantagens da Estrutura Recomendada:

### 1. **Separação Clara de Responsabilidades**
- `pkg/` - Código público que outros podem importar
- `internal/` - Código privado do projeto
- `cmd/` - Aplicações executáveis
- `examples/` - Exemplos de uso

### 2. **Organização por Funcionalidade**
- Cada middleware em seu próprio diretório
- Testes próximos ao código
- Documentação organizada

### 3. **Compatibilidade com Ferramentas Go**
- `go mod` funciona melhor
- IDEs reconhecem a estrutura
- Ferramentas de CI/CD padrão

### 4. **Facilita Manutenção**
- Fácil encontrar código específico
- Testes organizados
- Documentação centralizada

## 🔧 Migração Sugerida:

### Fase 1: Limpeza Básica
1. Remover executáveis (.exe)
2. Criar .gitignore
3. Mover arquivos de demo para internal/

### Fase 2: Reorganização de Middleware
1. Criar pkg/forge/middleware/
2. Mover cada middleware para seu diretório
3. Atualizar imports

### Fase 3: Estrutura Completa
1. Criar internal/ e docs/
2. Reorganizar cmd/
3. Atualizar documentação

## 📋 Benefícios Imediatos:

- ✅ Projeto mais profissional
- ✅ Fácil de navegar e entender
- ✅ Melhor para colaboração
- ✅ Compatível com padrões Go
- ✅ Preparado para crescimento