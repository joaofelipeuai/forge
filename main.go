// Example usage of Velocity Framework
package main

import (
	"log"
	"time"
)

func main() {
	app := New()
	
	// Middleware
	app.Use(Recovery())
	app.Use(Logger())
	app.Use(CORS())
	app.Use(RateLimiter(100, time.Minute))
	
	// Routes
	app.GET("/", func(c *Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "Welcome to Velocity Framework!",
			"version": Version,
			"features": []string{
				"Type-safe routing",
				"Built-in middleware",
				"Express-like API",
				"High performance",
				"Rate limiting",
				"CORS support",
				"Graceful shutdown",
				"Zero dependencies",
			},
		})
	})
	
	app.GET("/users/:id", func(c *Context) error {
		userID := c.Params["id"]
		return c.JSON(200, map[string]interface{}{
			"user_id": userID,
			"name":    "John Doe",
			"email":   "john@example.com",
		})
	})
	
	app.POST("/users", func(c *Context) error {
		return c.JSON(201, map[string]interface{}{
			"message": "User created successfully",
			"id":      "12345",
		})
	})
	
	app.GET("/health", func(c *Context) error {
		return c.JSON(200, map[string]interface{}{
			"status":    "healthy",
			"timestamp": time.Now().Unix(),
		})
	})
	
	// Custom middleware example
	app.Use(func(c *Context) error {
		c.Set("request_id", "req-123456")
		return c.Next()
	})
	
	app.GET("/request-info", func(c *Context) error {
		return c.JSON(200, map[string]interface{}{
			"request_id": c.Get("request_id"),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Query,
		})
	})
	
	// Start server
	log.Fatal(app.Listen(":3000"))
}
