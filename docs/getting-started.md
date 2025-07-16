# Getting Started - Forge Framework

## 🚀 Quick Start

### Installation
```bash
go get github.com/joaofelipeuai/forge
```

### Basic Usage
```go
package main

import (
    "log"
    "github.com/joaofelipeuai/forge"
)

func main() {
    app := forge.New()
    
    app.GET("/", func(c *forge.Context) error {
        return c.JSON(200, map[string]string{
            "message": "Hello, Forge!",
        })
    })
    
    log.Fatal(app.Listen(":8080"))
}
```

## 📚 Next Steps

- [Middleware Guide](middleware.md)
- [API Reference](api-reference.md)
- [Examples](../examples/)

## 🔗 Links

- [GitHub Repository](https://github.com/joaofelipeuai/forge)
- [Contributing Guide](../CONTRIBUTING.md)