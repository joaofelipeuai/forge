# Forge Framework 🔨

**Forge powerful web applications with Go**

Um framework web moderno e type-safe para Go, inspirado no Express.js mas otimizado para as características únicas do Go.

**Criado por:** [João Felipe Souza](https://github.com/joaofelipeuai)

## ✨ Características

- **Type-Safe**: Aproveita o sistema de tipos do Go para maior segurança
- **Performance**: Otimizado para alta performance com pools de objetos e concorrência
- **Express-like API**: Sintaxe familiar para desenvolvedores vindos do Node.js
- **Middleware Built-in**: Logger, CORS, Rate Limiting e mais
- **Zero Dependencies**: Usa apenas a biblioteca padrão do Go
- **Graceful Shutdown**: Suporte nativo para shutdown gracioso
- **Route Parameters**: Suporte completo para parâmetros de rota estilo Express

## 🚀 Início Rápido

```go
package main

import (
    "log"
    "github.com/joaofelipeuai/forge"
)

func main() {
    app := forge.New()
    
    // Middleware
    app.Use(forge.Logger())
    app.Use(forge.CORS())
    
    // Rotas
    app.GET("/", func(c *forge.Context) error {
        return c.JSON(200, map[string]string{
            "message": "Hello, Forge!",
        })
    })
    
    app.GET("/users/:id", func(c *forge.Context) error {
        userID := c.Params["id"]
        return c.JSON(200, map[string]string{
            "user_id": userID,
        })
    })
    
    log.Fatal(app.Listen(":8080"))
}
```

## 📦 Instalação

```bash
go mod init seu-projeto
go get github.com/joaofelipeuai/forge
```

## 🛠️ API Reference

### Criando uma aplicação

```go
app := forge.New()
```

### Rotas

```go
app.GET("/path", handler)
app.POST("/path", handler)
app.PUT("/path", handler)
app.DELETE("/path", handler)
```

### Middleware

```go
app.Use(middleware)
```

### Context Methods

```go
// JSON response
c.JSON(200, data)

// String response
c.String(200, "Hello World")

// Set/Get local values
c.Set("key", value)
value := c.Get("key")

// Access parameters and query
userID := c.Params["id"]
page := c.Query["page"]
```

### Middleware Built-in

```go
// Logger
app.Use(forge.Logger())

// CORS
app.Use(forge.CORS())

// Rate Limiting
app.Use(forge.RateLimiter(100, time.Minute))
```

## 🔧 Melhorias sobre Express.js

1. **Type Safety**: Aproveitamento completo do sistema de tipos do Go
2. **Performance**: Sem overhead de interpretação, compilado para código nativo
3. **Concorrência**: Goroutines nativas para handling de requests
4. **Memory Management**: Garbage collector otimizado do Go
5. **Zero Dependencies**: Não depende de pacotes externos
6. **Built-in Rate Limiting**: Rate limiting nativo sem dependências externas

## 📊 Benchmarks

```
BenchmarkVelocity-8    50000    25000 ns/op    1024 B/op    8 allocs/op
```

## 🤝 Contribuindo

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está licenciado sob a MIT License - veja o arquivo [LICENSE](LICENSE) para detalhes.

## 🙏 Agradecimentos

- Inspirado pelo Express.js
- Comunidade Go pela excelente documentação
- Todos os contribuidores

## 🆕 Funcionalidades Avançadas

### 🔌 WebSocket Support
```go
// WebSocket endpoint
app.WebSocket("/ws", func(conn *forge.WebSocketConnection) {
    conn.Send("Welcome!")
    // Handle connection...
})

// Broadcasting
broadcaster := forge.WebSocketBroadcast()
broadcaster.Broadcast("Message to all clients")
```

### 🎨 Template Engine
```go
// Setup template engine
engine := forge.NewTemplateEngine("templates", "html")
engine.SetDevMode(true) // Auto-reload in development
app.SetTemplateEngine(engine)

// Render templates
app.GET("/page", func(c *forge.Context) error {
    return c.Render(200, "index", map[string]interface{}{
        "Title": "Hello World",
    })
})
```

### 🔐 JWT Authentication
```go
// JWT configuration
jwtConfig := forge.NewJWTConfig("your-secret-key")

// Generate token
token, err := jwtConfig.GenerateToken(map[string]interface{}{
    "user_id": "123",
    "role": "admin",
})

// Protect routes
app.Use(forge.JWTAuth(jwtConfig))
```

### 📁 File Upload
```go
// File upload configuration
uploadConfig := forge.NewUploadConfig("./uploads")
uploadConfig.MaxFileSize = 10 << 20 // 10MB

// Upload middleware
app.POST("/upload", forge.FileUpload(uploadConfig), func(c *forge.Context) error {
    files := c.GetUploadedFiles()
    return c.JSON(200, map[string]interface{}{"files": files})
})

// Image upload with validation
app.POST("/images", forge.ImageUpload("./uploads", 5<<20), handler)
```

### 🔥 Hot Reload
```go
// Enable hot reload for development
app.EnableHotReload(".", "templates", "static")

// Or start with hot reload
app.ListenWithHotReload(":8080", ".", "templates")
```

## 🔮 Roadmap

- [x] WebSocket support ✅
- [x] Template engine integration ✅
- [x] File upload middleware ✅
- [x] JWT authentication middleware ✅
- [x] Hot reload em desenvolvimento ✅
- [ ] Database integration helpers
- [ ] Metrics e monitoring built-in
- [ ] GraphQL support
- [ ] Session management
- [ ] Caching middleware