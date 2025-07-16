# Middleware Guide - Forge Framework

## 🛡️ Built-in Middleware

### Recovery
```go
app.Use(forge.Recovery())
```

### Logger
```go
app.Use(forge.Logger())
```

### CORS
```go
app.Use(forge.CORS())
```

### Rate Limiting
```go
app.Use(forge.RateLimiter(100, time.Minute))
```

### JWT Authentication
```go
jwtConfig := forge.NewJWTConfig("secret-key")
app.Use(forge.JWTAuth(jwtConfig))
```

## 🔧 Custom Middleware

```go
func CustomMiddleware() forge.MiddlewareFunc {
    return func(c *forge.Context) error {
        // Your middleware logic here
        return c.Next()
    }
}

app.Use(CustomMiddleware())
```