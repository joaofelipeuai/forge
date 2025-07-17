package main

import (
	"log"
	"strconv"
	
	"github.com/joaofelipeuai/forge"
)

func main() {
	app := forge.New()
	
	// Middleware básico
	app.Use(forge.Recovery())
	app.Use(forge.Logger())
	app.Use(forge.CORS())
	app.Use(forge.ValidationMiddleware())
	
	// Exemplo de validação de parâmetros
	app.GET("/users/:id", func(c *forge.Context) error {
		// Validar se ID é um número inteiro
		err := c.ValidateParam("id", func(value string) error {
			validator := c.GetValidator()
			if err := validator.ValidateRequired(value, "user ID"); err != nil {
				return err
			}
			return validator.ValidateInteger(value, "user ID")
		})
		
		if err != nil {
			return c.JSON(400, map[string]string{
				"error": err.Error(),
			})
		}
		
		userID := c.Params["id"]
		id, _ := strconv.Atoi(userID) // Já validado acima
		
		return c.JSON(200, map[string]interface{}{
			"user_id": id,
			"name":    "User " + userID,
			"email":   "user" + userID + "@example.com",
		})
	})
	
	// Exemplo de validação de query parameters
	app.GET("/search", func(c *forge.Context) error {
		// Validar query parameter 'q'
		err := c.ValidateQuery("q", func(value string) error {
			validator := c.GetValidator()
			if err := validator.ValidateRequired(value, "search query"); err != nil {
				return err
			}
			return validator.ValidateLength(value, 2, 100, "search query")
		})
		
		if err != nil {
			return c.JSON(400, map[string]string{
				"error": err.Error(),
			})
		}
		
		query := c.Query["q"]
		return c.JSON(200, map[string]interface{}{
			"query":   query,
			"results": []string{"result1", "result2", "result3"},
		})
	})
	
	// Exemplo de validação de email
	app.POST("/subscribe", func(c *forge.Context) error {
		email := c.Query["email"]
		
		validator := c.GetValidator()
		if err := validator.ValidateEmail(email); err != nil {
			return c.JSON(400, map[string]string{
				"error": err.Error(),
			})
		}
		
		return c.JSON(200, map[string]interface{}{
			"message": "Subscription successful!",
			"email":   email,
		})
	})
	
	// Exemplo de múltiplas validações
	app.POST("/register", func(c *forge.Context) error {
		username := c.Query["username"]
		email := c.Query["email"]
		age := c.Query["age"]
		
		validator := c.GetValidator()
		
		// Validar username
		if err := validator.ValidateRequired(username, "username"); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		if err := validator.ValidateLength(username, 3, 20, "username"); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		if err := validator.ValidateAlphanumeric(username, "username"); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		
		// Validar email
		if err := validator.ValidateEmail(email); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		
		// Validar idade
		if err := validator.ValidateInteger(age, "age"); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		if err := validator.ValidateRange(age, 13, 120, "age"); err != nil {
			return c.JSON(400, map[string]string{"error": err.Error()})
		}
		
		return c.JSON(200, map[string]interface{}{
			"message":  "Registration successful!",
			"username": username,
			"email":    email,
			"age":      age,
		})
	})
	
	log.Println("🔨 Servidor de validação iniciando na porta :8080")
	log.Println("📋 Endpoints disponíveis:")
	log.Println("   GET  /users/:id        - Validação de parâmetro (ID deve ser inteiro)")
	log.Println("   GET  /search?q=        - Validação de query (q deve ter 2-100 chars)")
	log.Println("   POST /subscribe?email= - Validação de email")
	log.Println("   POST /register?username=&email=&age= - Múltiplas validações")
	log.Println("")
	log.Println("💡 Exemplos de teste:")
	log.Println("   curl http://localhost:8080/users/123")
	log.Println("   curl http://localhost:8080/users/abc  # Erro: deve ser inteiro")
	log.Println("   curl http://localhost:8080/search?q=golang")
	log.Println("   curl http://localhost:8080/search?q=a  # Erro: muito curto")
	log.Println("")
	
	log.Fatal(app.Listen(":8080"))
}