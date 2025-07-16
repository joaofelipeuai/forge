package main

import (
	"log"
	"time"
	
	"github.com/velocity-framework/velocity"
)

func main() {
	app := velocity.New()
	
	// Middleware
	app.Use(velocity.Recovery())
	app.Use(velocity.Logger())
	app.Use(velocity.CORS())
	app.Use(velocity.RateLimiter(100, time.Minute))
	
	// Routes
	app.GET("/", func(c *velocity.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "Welcome to Velocity Framework!",
			"version": velocity.Version,
			"docs":    "https://github.com/velocity-framework/velocity",
		})
	})
	
	app.GET("/users/:id", func(c *velocity.Context) error {
		userID := c.Params["id"]
		return c.JSON(200, map[string]interface{}{
			"user_id": userID,
			"name":    "John Doe",
			"email":   "john@example.com",
		})
	})
	
	app.POST("/users", func(c *velocity.Context) error {
		return c.JSON(201, map[string]interface{}{
			"message": "User created successfully",
			"id":      "12345",
		})
	})
	
	app.GET("/health", func(c *velocity.Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})
	
	// Custom middleware example
	app.Use(func(c *velocity.Context) error {
		c.Set("request_id", "req-123456")
		return c.Next()
	})
	
	app.GET("/request-info", func(c *velocity.Context) error {
		return c.JSON(200, map[string]interface{}{
			"request_id": c.Get("request_id"),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Query,
		})
	})
	
	log.Fatal(app.Listen(":8080"))
}