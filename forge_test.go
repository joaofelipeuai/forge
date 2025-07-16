package main

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	app := New()
	if app == nil {
		t.Fatal("New() returned nil")
	}
	
	if len(app.routes) != 0 {
		t.Errorf("Expected 0 routes, got %d", len(app.routes))
	}
	
	if len(app.middleware) != 0 {
		t.Errorf("Expected 0 middleware, got %d", len(app.middleware))
	}
}

func TestGETRoute(t *testing.T) {
	app := New()
	
	app.GET("/test", func(c *Context) error {
		return c.String(200, "Hello World")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	app.ServeHTTP(w, req)
	
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	if w.Body.String() != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", w.Body.String())
	}
}

func TestPOSTRoute(t *testing.T) {
	app := New()
	
	app.POST("/users", func(c *Context) error {
		return c.JSON(201, map[string]string{"message": "Created"})
	})
	
	req := httptest.NewRequest("POST", "/users", strings.NewReader(`{"name":"John"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	
	app.ServeHTTP(w, req)
	
	if w.Code != 201 {
		t.Errorf("Expected status 201, got %d", w.Code)
	}
}

func TestRouteParams(t *testing.T) {
	app := New()
	
	app.GET("/users/:id", func(c *Context) error {
		userID := c.Params["id"]
		return c.String(200, "User ID: "+userID)
	})
	
	req := httptest.NewRequest("GET", "/users/123", nil)
	w := httptest.NewRecorder()
	
	app.ServeHTTP(w, req)
	
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	expected := "User ID: 123"
	if w.Body.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, w.Body.String())
	}
}

func TestQueryParams(t *testing.T) {
	app := New()
	
	app.GET("/search", func(c *Context) error {
		query := c.Query["q"]
		return c.String(200, "Query: "+query)
	})
	
	req := httptest.NewRequest("GET", "/search?q=golang", nil)
	w := httptest.NewRecorder()
	
	app.ServeHTTP(w, req)
	
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	expected := "Query: golang"
	if w.Body.String() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, w.Body.String())
	}
}

func TestMiddleware(t *testing.T) {
	app := New()
	
	// Add middleware that sets a header
	app.Use(func(c *Context) error {
		c.Header("X-Test", "middleware-works")
		return c.Next()
	})
	
	app.GET("/test", func(c *Context) error {
		return c.String(200, "OK")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	app.ServeHTTP(w, req)
	
	if w.Header().Get("X-Test") != "middleware-works" {
		t.Errorf("Expected header 'X-Test: middleware-works', got '%s'", w.Header().Get("X-Test"))
	}
}

func TestCORSMiddleware(t *testing.T) {
	app := New()
	app.Use(CORS())
	
	app.GET("/test", func(c *Context) error {
		return c.String(200, "OK")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	app.ServeHTTP(w, req)
	
	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS headers not set correctly")
	}
}

func TestRateLimiter(t *testing.T) {
	app := New()
	app.Use(RateLimiter(2, time.Second))
	
	app.GET("/test", func(c *Context) error {
		return c.String(200, "OK")
	})
	
	// First request should pass
	req1 := httptest.NewRequest("GET", "/test", nil)
	req1.RemoteAddr = "127.0.0.1:8080"
	w1 := httptest.NewRecorder()
	app.ServeHTTP(w1, req1)
	
	if w1.Code != 200 {
		t.Errorf("First request should pass, got status %d", w1.Code)
	}
	
	// Second request should pass
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.RemoteAddr = "127.0.0.1:8080"
	w2 := httptest.NewRecorder()
	app.ServeHTTP(w2, req2)
	
	if w2.Code != 200 {
		t.Errorf("Second request should pass, got status %d", w2.Code)
	}
	
	// Third request should be rate limited
	req3 := httptest.NewRequest("GET", "/test", nil)
	req3.RemoteAddr = "127.0.0.1:8080"
	w3 := httptest.NewRecorder()
	app.ServeHTTP(w3, req3)
	
	if w3.Code != 429 {
		t.Errorf("Third request should be rate limited, got status %d", w3.Code)
	}
}

func TestNotFound(t *testing.T) {
	app := New()
	
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()
	
	app.ServeHTTP(w, req)
	
	if w.Code != 404 {
		t.Errorf("Expected status 404, got %d", w.Code)
	}
}

func TestContextSetGet(t *testing.T) {
	app := New()
	
	app.Use(func(c *Context) error {
		c.Set("test_key", "test_value")
		return c.Next()
	})
	
	app.GET("/test", func(c *Context) error {
		value := c.Get("test_key")
		if value == nil {
			return c.String(500, "Value not found")
		}
		return c.String(200, value.(string))
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()
	
	app.ServeHTTP(w, req)
	
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
	
	if w.Body.String() != "test_value" {
		t.Errorf("Expected 'test_value', got '%s'", w.Body.String())
	}
}

// Benchmark tests
func BenchmarkForgeSimpleRoute(b *testing.B) {
	app := New()
	app.GET("/test", func(c *Context) error {
		return c.String(200, "OK")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}

func BenchmarkForgeWithMiddleware(b *testing.B) {
	app := New()
	app.Use(Logger())
	app.Use(CORS())
	
	app.GET("/test", func(c *Context) error {
		return c.String(200, "OK")
	})
	
	req := httptest.NewRequest("GET", "/test", nil)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		app.ServeHTTP(w, req)
	}
}