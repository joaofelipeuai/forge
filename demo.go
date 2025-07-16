package main

import (
	"fmt"
	"log"
	"time"
)

func runDemo() {
	fmt.Println("🚀 Forge Framework Demo")
	fmt.Println("==========================")
	
	app := New()
	
	// Middleware básico
	app.Use(Recovery())
	app.Use(Logger())
	app.Use(CORS())
	
	// Rotas de demonstração
	app.GET("/", func(c *Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "🚀 Forge Framework funcionando!",
			"version": Version,
			"timestamp": time.Now().Unix(),
		})
	})
	
	app.GET("/users/:id", func(c *Context) error {
		userID := c.Params["id"]
		return c.JSON(200, map[string]interface{}{
			"user_id": userID,
			"name": fmt.Sprintf("User %s", userID),
			"email": fmt.Sprintf("user%s@example.com", userID),
		})
	})
	
	app.GET("/search", func(c *Context) error {
		query := c.Query["q"]
		if query == "" {
			query = "nenhuma"
		}
		return c.JSON(200, map[string]interface{}{
			"query": query,
			"results": []string{"resultado1", "resultado2", "resultado3"},
		})
	})
	
	// JWT Demo
	jwtConfig := NewJWTConfig("demo-secret-key")
	
	app.POST("/login", func(c *Context) error {
		token, err := jwtConfig.GenerateToken(map[string]interface{}{
			"user_id": "123",
			"username": "demo_user",
		})
		if err != nil {
			return c.JSON(500, map[string]string{"error": "Erro ao gerar token"})
		}
		
		return c.JSON(200, map[string]interface{}{
			"token": token,
			"message": "Login realizado com sucesso!",
		})
	})
	
	// Rota protegida
	app.Use(JWTOptional(jwtConfig)) // JWT opcional para não quebrar outras rotas
	app.GET("/profile", func(c *Context) error {
		jwt := GetJWT(c)
		if jwt == nil {
			return c.JSON(401, map[string]string{"error": "Token necessário"})
		}
		
		return c.JSON(200, map[string]interface{}{
			"user_id": jwt.Payload.Claims["user_id"],
			"username": jwt.Payload.Claims["username"],
			"message": "Perfil acessado com sucesso!",
		})
	})
	
	fmt.Println("🌐 Servidor iniciando na porta 3000...")
	fmt.Println("📋 Endpoints disponíveis:")
	fmt.Println("   GET  /           - Página inicial")
	fmt.Println("   GET  /users/:id  - Buscar usuário")
	fmt.Println("   GET  /search?q=  - Buscar")
	fmt.Println("   POST /login      - Login (gera JWT)")
	fmt.Println("   GET  /profile    - Perfil (requer JWT)")
	fmt.Println("")
	
	log.Fatal(app.Listen(":8080"))
}

func main() {
	runDemo()
}
