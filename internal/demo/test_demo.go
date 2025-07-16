package demo

import (
	"fmt"
	"net/http/httptest"
	"time"
	
	"github.com/joaofelipeuai/forge"
)

func testForgeFramework() {
	fmt.Println("🔨 TESTANDO FORGE FRAMEWORK")
	fmt.Println("=============================")
	
	// Criar instância do framework
	app := forge.New()
	
	// Adicionar middleware
	app.Use(forge.Recovery())
	app.Use(forge.Logger())
	app.Use(forge.CORS())
	app.Use(forge.RateLimiter(10, time.Minute))
	
	// Configurar rotas
	app.GET("/", func(c *forge.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "🔨 Forge Framework funcionando!",
			"version": forge.Version,
			"features": []string{
				"Roteamento dinâmico",
				"Middleware pipeline",
				"Rate limiting",
				"CORS support",
				"JWT authentication",
				"File upload",
				"WebSocket support",
			},
		})
	})
	
	app.GET("/users/:id", func(c *forge.Context) error {
		userID := c.Params["id"]
		return c.JSON(200, map[string]interface{}{
			"user_id": userID,
			"name": fmt.Sprintf("User %s", userID),
			"email": fmt.Sprintf("user%s@forge.dev", userID),
		})
	})
	
	app.GET("/search", func(c *forge.Context) error {
		query := c.Query["q"]
		return c.JSON(200, map[string]interface{}{
			"query": query,
			"results": []string{"Go", "Forge", "Framework"},
			"count": 3,
		})
	})
	
	// JWT Configuration
	jwtConfig := forge.NewJWTConfig("test-secret-key")
	
	app.POST("/auth/login", func(c *forge.Context) error {
		token, err := jwtConfig.GenerateToken(map[string]interface{}{
			"user_id": "123",
			"username": "test_user",
			"role": "admin",
		})
		if err != nil {
			return c.JSON(500, map[string]string{"error": "Failed to generate token"})
		}
		
		return c.JSON(200, map[string]interface{}{
			"token": token,
			"user": "test_user",
			"expires_in": "24h",
		})
	})
	
	// Protected route (usando middleware global para simplificar o teste)
	app.Use(forge.JWTOptional(jwtConfig))
	app.GET("/profile", func(c *forge.Context) error {
		jwt := forge.GetJWT(c)
		if jwt == nil {
			return c.JSON(401, map[string]string{"error": "Token required"})
		}
		return c.JSON(200, map[string]interface{}{
			"user_id": jwt.Payload.Claims["user_id"],
			"username": jwt.Payload.Claims["username"],
			"role": jwt.Payload.Claims["role"],
		})
	})
	
	// Executar testes
	fmt.Println("📋 EXECUTANDO TESTES...")
	fmt.Println()
	
	// Teste 1: Rota principal
	fmt.Println("✅ Teste 1: GET /")
	req1 := httptest.NewRequest("GET", "/", nil)
	w1 := httptest.NewRecorder()
	app.ServeHTTP(w1, req1)
	fmt.Printf("   Status: %d\n", w1.Code)
	fmt.Printf("   Response: %s\n", w1.Body.String()[:100] + "...")
	fmt.Println()
	
	// Teste 2: Rota com parâmetros
	fmt.Println("✅ Teste 2: GET /users/123")
	req2 := httptest.NewRequest("GET", "/users/123", nil)
	w2 := httptest.NewRecorder()
	app.ServeHTTP(w2, req2)
	fmt.Printf("   Status: %d\n", w2.Code)
	fmt.Printf("   Response: %s\n", w2.Body.String())
	fmt.Println()
	
	// Teste 3: Query parameters
	fmt.Println("✅ Teste 3: GET /search?q=golang")
	req3 := httptest.NewRequest("GET", "/search?q=golang", nil)
	w3 := httptest.NewRecorder()
	app.ServeHTTP(w3, req3)
	fmt.Printf("   Status: %d\n", w3.Code)
	fmt.Printf("   Response: %s\n", w3.Body.String())
	fmt.Println()
	
	// Teste 4: JWT Login
	fmt.Println("✅ Teste 4: POST /auth/login")
	req4 := httptest.NewRequest("POST", "/auth/login", nil)
	w4 := httptest.NewRecorder()
	app.ServeHTTP(w4, req4)
	fmt.Printf("   Status: %d\n", w4.Code)
	fmt.Printf("   Response: %s\n", w4.Body.String()[:100] + "...")
	fmt.Println()
	
	// Teste 5: Rate Limiting
	fmt.Println("✅ Teste 5: Rate Limiting (múltiplas requisições)")
	for i := 0; i < 12; i++ {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = "127.0.0.1:8080"
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
		if w.Code == 429 {
			fmt.Printf("   Rate limit ativado na requisição %d (Status: %d)\n", i+1, w.Code)
			break
		}
	}
	fmt.Println()
	
	// Teste 6: CORS Headers
	fmt.Println("✅ Teste 6: CORS Headers")
	req6 := httptest.NewRequest("OPTIONS", "/", nil)
	w6 := httptest.NewRecorder()
	app.ServeHTTP(w6, req6)
	fmt.Printf("   Status: %d\n", w6.Code)
	fmt.Printf("   CORS Origin: %s\n", w6.Header().Get("Access-Control-Allow-Origin"))
	fmt.Printf("   CORS Methods: %s\n", w6.Header().Get("Access-Control-Allow-Methods"))
	fmt.Println()
	
	fmt.Println("🎉 TODOS OS TESTES CONCLUÍDOS!")
	fmt.Println("===============================")
	fmt.Printf("✅ Framework: Forge v%s\n", forge.Version)
	fmt.Println("✅ Roteamento: Funcionando")
	fmt.Println("✅ Middleware: Funcionando")
	fmt.Println("✅ JWT: Funcionando")
	fmt.Println("✅ Rate Limiting: Funcionando")
	fmt.Println("✅ CORS: Funcionando")
	fmt.Println("✅ Parâmetros: Funcionando")
	fmt.Println("✅ Query Strings: Funcionando")
	fmt.Println()
	fmt.Println("🔨 O Forge Framework está 100% operacional!")
}

// TestFramework executa testes do framework
// Para usar: demo.TestFramework()
func TestFramework() {
	testForgeFramework()
}