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
	
	// JWT Configuration
	jwtConfig := forge.NewJWTConfig("test-secret-key")
	jwtConfig.Expiration = 24 * time.Hour
	
	// Rota pública para teste
	app.GET("/", func(c *forge.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "🔨 Forge JWT Test Server",
			"version": forge.Version,
			"endpoints": []string{
				"POST /auth/login - Login e obter token",
				"GET /profile - Perfil protegido (requer token)",
				"GET /public - Rota pública",
			},
		})
	})
	
	// Rota pública
	app.GET("/public", func(c *forge.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "Esta é uma rota pública",
			"timestamp": time.Now().Unix(),
		})
	})
	
	// Login route - gera token JWT
	app.POST("/auth/login", func(c *forge.Context) error {
		// Simular validação de credenciais
		claims := map[string]interface{}{
			"sub":      "user123",
			"username": "demo_user",
			"role":     "admin",
		}
		
		token, err := jwtConfig.GenerateToken(claims)
		if err != nil {
			return c.JSON(500, map[string]string{"error": "Falha ao gerar token"})
		}
		
		return c.JSON(200, map[string]interface{}{
			"token":    token,
			"username": "demo_user",
			"expires":  time.Now().Add(jwtConfig.Expiration).Unix(),
			"message":  "Login realizado com sucesso!",
		})
	})
	
	// Rota protegida - requer JWT
	app.GET("/profile", func(c *forge.Context) error {
		// Aplicar middleware JWT inline
		jwtMiddleware := forge.JWTAuth(jwtConfig)
		if err := jwtMiddleware(c); err != nil {
			return err
		}
		
		// Se chegou até aqui, o token é válido
		jwt := forge.GetJWT(c)
		userID := forge.GetUserID(c)
		
		return c.JSON(200, map[string]interface{}{
			"user_id":  userID,
			"username": jwt.Payload.Claims["username"],
			"role":     jwt.Payload.Claims["role"],
			"exp":      jwt.Payload.ExpirationTime,
			"message":  "Perfil acessado com sucesso!",
		})
	})
	
	// Rota de teste para verificar token
	app.GET("/verify", func(c *forge.Context) error {
		// Usar JWT opcional para não falhar se não houver token
		jwtMiddleware := forge.JWTOptional(jwtConfig)
		if err := jwtMiddleware(c); err != nil {
			return err
		}
		
		jwt := forge.GetJWT(c)
		if jwt == nil {
			return c.JSON(200, map[string]interface{}{
				"authenticated": false,
				"message":       "Nenhum token fornecido",
			})
		}
		
		return c.JSON(200, map[string]interface{}{
			"authenticated": true,
			"user_id":       jwt.Payload.Subject,
			"username":      jwt.Payload.Claims["username"],
			"message":       "Token válido!",
		})
	})
	
	log.Println("🔨 Servidor JWT Test iniciando na porta :8081")
	log.Println("📋 Endpoints disponíveis:")
	log.Println("   GET  /           - Informações do servidor")
	log.Println("   GET  /public     - Rota pública")
	log.Println("   POST /auth/login - Login (gera JWT)")
	log.Println("   GET  /profile    - Perfil (requer JWT)")
	log.Println("   GET  /verify     - Verificar token (opcional)")
	log.Println("")
	log.Println("💡 Para testar:")
	log.Println("   1. POST /auth/login para obter token")
	log.Println("   2. GET /profile com header: Authorization: Bearer <token>")
	log.Println("")
	
	log.Fatal(app.Listen(":8081"))
}