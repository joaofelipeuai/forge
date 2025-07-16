package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	app := New()
	
	// Enable hot reload for development
	app.EnableHotReload(".", "templates", "static")
	
	// Basic middleware
	app.Use(Recovery())
	app.Use(Logger())
	app.Use(CORS())
	app.Use(RateLimiter(100, time.Minute))
	
	// JWT Configuration
	jwtConfig := NewJWTConfig("your-secret-key-here")
	jwtConfig.Expiration = 24 * time.Hour
	
	// File upload configuration
	uploadConfig := NewUploadConfig("./uploads")
	uploadConfig.MaxFileSize = 10 << 20 // 10MB
	
	// Template engine setup
	templateEngine := NewTemplateEngine("templates", "html")
	templateEngine.SetDevMode(true) // Auto-reload templates in development
	templateEngine.LoadTemplates()
	app.SetTemplateEngine(templateEngine)
	
	// Serve static files and uploads
	app.ServeUploads("/uploads", "./uploads")
	
	// Public routes
	app.GET("/", func(c *Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "🚀 Velocity Framework - Advanced Features Demo",
			"version": Version,
			"features": []string{
				"WebSocket Support",
				"Template Engine",
				"JWT Authentication",
				"File Upload",
				"Hot Reload",
				"Rate Limiting",
				"CORS Support",
			},
		})
	})
	
	// Authentication routes
	app.POST("/auth/login", func(c *Context) error {
		// In a real app, validate credentials here
		username := "demo_user"
		
		claims := map[string]interface{}{
			"sub":      "user123",
			"username": username,
			"role":     "admin",
		}
		
		token, err := jwtConfig.GenerateToken(claims)
		if err != nil {
			return c.JSON(500, map[string]string{"error": "Failed to generate token"})
		}
		
		return c.JSON(200, map[string]interface{}{
			"token":    token,
			"username": username,
			"expires":  time.Now().Add(jwtConfig.Expiration).Unix(),
		})
	})
	
	// Protected routes group
	protected := app // In a real implementation, you'd create route groups
	protected.Use(JWTAuth(jwtConfig))
	
	protected.GET("/profile", func(c *Context) error {
		jwt := GetJWT(c)
		userID := GetUserID(c)
		
		return c.JSON(200, map[string]interface{}{
			"user_id":  userID,
			"username": jwt.Payload.Claims["username"],
			"role":     jwt.Payload.Claims["role"],
			"exp":      jwt.Payload.ExpirationTime,
		})
	})
	
	// File upload routes
	app.POST("/upload", FileUpload(uploadConfig), func(c *Context) error {
		result := c.GetUploadResult()
		
		if !result.Success {
			return c.JSON(400, map[string]interface{}{
				"success": false,
				"errors":  result.Errors,
			})
		}
		
		return c.JSON(200, map[string]interface{}{
			"success": true,
			"files":   result.Files,
			"count":   len(result.Files),
		})
	})
	
	// Image upload with specific validation
	app.POST("/upload/image", ImageUpload("./uploads/images", 5<<20), func(c *Context) error {
		files := c.GetUploadedFiles()
		if len(files) == 0 {
			return c.JSON(400, map[string]string{"error": "No files uploaded"})
		}
		
		return c.JSON(200, map[string]interface{}{
			"message": "Images uploaded successfully",
			"files":   files,
		})
	})
	
	// WebSocket endpoint
	broadcaster := WebSocketBroadcast()
	
	app.WebSocket("/ws", func(conn *WebSocketConnection) {
		log.Println("New WebSocket connection")
		broadcaster.AddConnection(conn)
		
		// Send welcome message
		conn.Send("Welcome to Velocity WebSocket!")
		
		// Broadcast to all connections
		broadcaster.Broadcast("New user joined the chat!")
		
		// Keep connection alive (in a real app, you'd handle incoming messages)
		time.Sleep(30 * time.Second)
		
		broadcaster.RemoveConnection(conn)
		conn.Close()
	})
	
	// WebSocket broadcast endpoint
	app.POST("/broadcast", func(c *Context) error {
		message := c.Query["message"]
		if message == "" {
			return c.JSON(400, map[string]string{"error": "Message is required"})
		}
		
		broadcaster.Broadcast(message)
		return c.JSON(200, map[string]string{"message": "Broadcast sent"})
	})
	
	// Template rendering example
	app.GET("/template", func(c *Context) error {
		data := map[string]interface{}{
			"Title":   "Velocity Template Demo",
			"Message": "Hello from Velocity Framework!",
			"Time":    time.Now().Format("2006-01-02 15:04:05"),
		}
		
		return c.Render(200, "index", data)
	})
	
	// Health check with system info
	app.GET("/health", func(c *Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
			"uptime":    time.Since(time.Now()).String(),
			"features": map[string]bool{
				"websocket":   true,
				"templates":   true,
				"jwt":         true,
				"file_upload": true,
				"hot_reload":  true,
			},
		})
	})
	
	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		
		log.Println("🛑 Shutting down server...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		
		if err := app.Shutdown(ctx); err != nil {
			log.Printf("Server shutdown error: %v", err)
		}
		os.Exit(0)
	}()
	
	// Start server
	log.Println("🚀 Starting Velocity server with all advanced features...")
	log.Fatal(app.Listen(":8080"))
}