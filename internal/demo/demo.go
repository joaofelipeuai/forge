package demo

import (
	"fmt"
	"log"
	"time"
	
	"github.com/joaofelipeuai/forge"
)

func runDemo() {
	fmt.Println("🚀 Forge Framework Demo")
	fmt.Println("==========================")
	
	app := forge.New()
	
	// Middleware básico
	app.Use(forge.Recovery())
	app.Use(forge.Logger())
	app.Use(forge.CORS())
	
	// Rotas de demonstração
	app.GET("/", func(c *forge.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "🚀 Forge Framework funcionando!",
			"version": forge.Version,
			"timestamp": time.Now().Unix(),
		})
	})
	
	app.GET("/users/:id", func(c *forge.Context) error {
		userID := c.Params["id"]
		return c.JSON(200, map[string]interface{}{
			"user_id": userID,
			"name": fmt.Sprintf("User %s", userID),
			"email": fmt.Sprintf("user%s@example.com", userID),
		})
	})
	
	app.GET("/search", func(c *forge.Context) error {
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
	jwtConfig := forge.NewJWTConfig("demo-secret-key")
	
	app.POST("/login", func(c *forge.Context) error {
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
	app.Use(forge.JWTOptional(jwtConfig)) // JWT opcional para não quebrar outras rotas
	app.GET("/profile", func(c *forge.Context) error {
		jwt := forge.GetJWT(c)
		if jwt == nil {
			return c.JSON(401, map[string]string{"error": "Token necessário"})
		}
		
		return c.JSON(200, map[string]interface{}{
			"user_id": jwt.Payload.Claims["user_id"],
			"username": jwt.Payload.Claims["username"],
			"message": "Perfil acessado com sucesso!",
		})
	})
	
	fmt.Println("🌐 Servidor iniciando na porta 8080...")
	fmt.Println("📋 Endpoints disponíveis:")
	fmt.Println("   GET  /           - Página inicial")
	fmt.Println("   GET  /users/:id  - Buscar usuário")
	fmt.Println("   GET  /search?q=  - Buscar")
	fmt.Println("   POST /login      - Login (gera JWT)")
	fmt.Println("   GET  /profile    - Perfil (requer JWT)")
	fmt.Println("")
	
	log.Fatal(app.Listen(":8080"))
}

// RunDemo executa uma demonstração do framework
// Para usar: demo.RunDemo()
func RunDemo() {
	runDemo()
}
