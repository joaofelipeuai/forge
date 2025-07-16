package main

import (
	"log"
	"time"
	
	"github.com/joaofelipeuai/forge"
)

func main() {
	app := forge.New()
	
	// Middleware
	app.Use(forge.Recovery())
	app.Use(forge.Logger())
	app.Use(forge.CORS())
	app.Use(forge.RateLimiter(100, time.Minute))
	
	// Routes
	app.GET("/", func(c *forge.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "Welcome to Forge Framework!",
			"version": forge.Version,
			"docs":    "https://github.com/joaofelipeuai/forge",
		})
	})
	
	app.GET("/users/:id", func(c *forge.Context) error {
		userID := c.Params["id"]
		return c.JSON(200, map[string]interface{}{
			"user_id": userID,
			"name":    "John Doe",
			"email":   "john@example.com",
		})
	})
	
	app.POST("/users", func(c *forge.Context) error {
		return c.JSON(201, map[string]interface{}{
			"message": "User created successfully",
			"id":      "12345",
		})
	})
	
	app.GET("/health", func(c *forge.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})
	
	// Custom middleware example
	app.Use(func(c *forge.Context) error {
		c.Set("request_id", "req-123456")
		return c.Next()
	})
	
	app.GET("/request-info", func(c *forge.Context) error {
		return c.JSON(200, map[string]interface{}{
			"request_id": c.Get("request_id"),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Query,
		})
	})
	
	log.Fatal(app.Listen(":8080"))
}
