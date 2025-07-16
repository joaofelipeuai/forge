package main

import (
	"log"
	"time"
	
	"github.com/joaofelipeuai/forge"
)

func main() {
	app := forge.New()
	
	// Middleware básico
	app.Use(forge.Recovery())
	app.Use(forge.Logger())
	app.Use(forge.CORS())
	
	// Rota simples
	app.GET("/", func(c *forge.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "Forge Framework funcionando!",
			"version": forge.Version,
			"status":  "ok",
		})
	})
	
	// Rota com parâmetros
	app.GET("/users/:id", func(c *forge.Context) error {
		userID := c.Params["id"]
		return c.JSON(200, map[string]interface{}{
			"user_id": userID,
			"name":    "Usuário Teste",
			"email":   "teste@exemplo.com",
		})
	})
	
	// Rota de health check
	app.GET("/health", func(c *forge.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})
	
	log.Println("🔨 Servidor Forge iniciando na porta :8080")
	log.Fatal(app.Listen(":8080"))
}