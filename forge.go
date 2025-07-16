// Package forge provides a modern, type-safe web framework for Go
// Inspired by Express.js but optimized for Go's unique characteristics
package forge

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Version of the Forge framework
const Version = "1.0.0"

// Core types and interfaces
type Context struct {
	Request    *http.Request
	Response   http.ResponseWriter
	Params     map[string]string
	Query      map[string]string
	Body       []byte
	locals     map[string]interface{}
	middleware []MiddlewareFunc
	index      int
	mu         sync.RWMutex
}

type MiddlewareFunc func(*Context) error
type HandlerFunc func(*Context) error

// Route represents a single route with its pattern and handler
type Route struct {
	Method  string
	Pattern string
	Handler HandlerFunc
	Regex   *regexp.Regexp
	Keys    []string
}

// Forge is the main framework struct
type Forge struct {
	routes         []*Route
	middleware     []MiddlewareFunc
	mu             sync.RWMutex
	server         *http.Server
	templateEngine *TemplateEngine
	hotReload      *HotReload
}

// New creates a new Forge instance
func New() *Forge {
	return &Forge{
		routes:     make([]*Route, 0),
		middleware: make([]MiddlewareFunc, 0),
	}
}

// Context methods
func (c *Context) JSON(status int, data interface{}) error {
	c.Response.Header().Set("Content-Type", "application/json")
	c.Response.WriteHeader(status)
	return json.NewEncoder(c.Response).Encode(data)
}

func (c *Context) String(status int, text string) error {
	c.Response.Header().Set("Content-Type", "text/plain")
	c.Response.WriteHeader(status)
	_, err := c.Response.Write([]byte(text))
	return err
}

func (c *Context) HTML(status int, html string) error {
	c.Response.Header().Set("Content-Type", "text/html")
	c.Response.WriteHeader(status)
	_, err := c.Response.Write([]byte(html))
	return err
}

func (c *Context) Status(code int) *Context {
	c.Response.WriteHeader(code)
	return c
}

func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.locals == nil {
		c.locals = make(map[string]interface{})
	}
	c.locals[key] = value
}

func (c *Context) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.locals[key]
}

func (c *Context) Next() error {
	c.index++
	if c.index < len(c.middleware) {
		return c.middleware[c.index](c)
	}
	return nil
}

// Header sets a response header
func (c *Context) Header(key, value string) {
	c.Response.Header().Set(key, value)
}

// Cookie sets a cookie
func (c *Context) Cookie(name, value string, maxAge int) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: maxAge,
		Path:   "/",
	}
	http.SetCookie(c.Response, cookie)
}

// Forge methods
func (f *Forge) Use(middleware MiddlewareFunc) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.middleware = append(f.middleware, middleware)
}

func (f *Forge) addRoute(method, pattern string, handler HandlerFunc) {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	route := &Route{
		Method:  method,
		Pattern: pattern,
		Handler: handler,
	}
	
	// Convert Express-style routes to regex
	route.Regex, route.Keys = f.compileRoute(pattern)
	f.routes = append(f.routes, route)
}

func (f *Forge) GET(pattern string, handler HandlerFunc) {
	f.addRoute("GET", pattern, handler)
}

func (f *Forge) POST(pattern string, handler HandlerFunc) {
	f.addRoute("POST", pattern, handler)
}

func (f *Forge) PUT(pattern string, handler HandlerFunc) {
	f.addRoute("PUT", pattern, handler)
}

func (f *Forge) DELETE(pattern string, handler HandlerFunc) {
	f.addRoute("DELETE", pattern, handler)
}

func (f *Forge) PATCH(pattern string, handler HandlerFunc) {
	f.addRoute("PATCH", pattern, handler)
}

func (f *Forge) OPTIONS(pattern string, handler HandlerFunc) {
	f.addRoute("OPTIONS", pattern, handler)
}

// Route compilation (Express-style to regex)
func (f *Forge) compileRoute(pattern string) (*regexp.Regexp, []string) {
	keys := make([]string, 0)
	regexPattern := "^"
	
	// Simple parameter extraction (:param)
	paramRegex := regexp.MustCompile(`:(\w+)`)
	matches := paramRegex.FindAllStringSubmatch(pattern, -1)
	
	for _, match := range matches {
		keys = append(keys, match[1])
	}
	
	// Replace :param with capture groups
	regexPattern += paramRegex.ReplaceAllString(pattern, `([^/]+)`)
	regexPattern += "$"
	
	regex, _ := regexp.Compile(regexPattern)
	return regex, keys
}

// HTTP handler implementation
func (f *Forge) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := &Context{
		Request:  r,
		Response: w,
		Params:   make(map[string]string),
		Query:    make(map[string]string),
		locals:   make(map[string]interface{}),
		index:    -1,
	}
	
	// Parse query parameters
	for key, values := range r.URL.Query() {
		if len(values) > 0 {
			ctx.Query[key] = values[0]
		}
	}
	
	// Find matching route
	var matchedRoute *Route
	f.mu.RLock()
	for _, route := range f.routes {
		if route.Method == r.Method && route.Regex.MatchString(r.URL.Path) {
			matchedRoute = route
			
			// Extract parameters
			matches := route.Regex.FindStringSubmatch(r.URL.Path)
			for i, key := range route.Keys {
				if i+1 < len(matches) {
					ctx.Params[key] = matches[i+1]
				}
			}
			break
		}
	}
	f.mu.RUnlock()
	
	if matchedRoute == nil {
		http.NotFound(w, r)
		return
	}
	
	// Set template engine in context if available
	if f.templateEngine != nil {
		ctx.Set("template_engine", f.templateEngine)
	}
	
	// Build middleware chain
	ctx.middleware = append(f.middleware, func(c *Context) error {
		return matchedRoute.Handler(c)
	})
	
	// Execute middleware chain
	if len(ctx.middleware) > 0 {
		ctx.index = 0
		if err := ctx.middleware[0](ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (f *Forge) Listen(addr string) error {
	f.server = &http.Server{
		Addr:         addr,
		Handler:      f,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	
	fmt.Printf("🔨 Forge v%s server running on %s\n", Version, addr)
	return f.server.ListenAndServe()
}

func (f *Forge) Shutdown(ctx context.Context) error {
	if f.server != nil {
		fmt.Println("🛑 Shutting down Forge server...")
		return f.server.Shutdown(ctx)
	}
	return nil
}

// Built-in middleware
func Logger() MiddlewareFunc {
	return func(c *Context) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)
		
		status := 200
		if err != nil {
			status = 500
		}
		
		log.Printf("[%d] %s %s - %v", status, c.Request.Method, c.Request.URL.Path, duration)
		return err
	}
}

func CORS() MiddlewareFunc {
	return func(c *Context) error {
		c.Response.Header().Set("Access-Control-Allow-Origin", "*")
		c.Response.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Response.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.Response.WriteHeader(http.StatusOK)
			return nil
		}
		
		return c.Next()
	}
}

func RateLimiter(requests int, window time.Duration) MiddlewareFunc {
	type client struct {
		requests int
		window   time.Time
	}
	
	clients := make(map[string]*client)
	mu := sync.RWMutex{}
	
	// Cleanup routine
	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()
		
		for range ticker.C {
			mu.Lock()
			now := time.Now()
			for ip, cl := range clients {
				if now.Sub(cl.window) > window {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()
	
	return func(c *Context) error {
		ip := c.Request.RemoteAddr
		if forwarded := c.Request.Header.Get("X-Forwarded-For"); forwarded != "" {
			ip = strings.Split(forwarded, ",")[0]
		}
		
		mu.Lock()
		defer mu.Unlock()
		
		now := time.Now()
		cl, exists := clients[ip]
		
		if !exists || now.Sub(cl.window) > window {
			clients[ip] = &client{requests: 1, window: now}
			return c.Next()
		}
		
		if cl.requests >= requests {
			return c.String(429, "Rate limit exceeded")
		}
		
		cl.requests++
		return c.Next()
	}
}

// Recovery middleware
func Recovery() MiddlewareFunc {
	return func(c *Context) error {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.String(500, "Internal Server Error")
			}
		}()
		
		return c.Next()
	}
}



// Hot reload integration
func (f *Forge) EnableHotReload(watchDirs ...string) {
	f.hotReload = NewHotReload()
	f.hotReload.Enable()
	
	if len(watchDirs) > 0 {
		f.hotReload.watchDirs = watchDirs
	}
	
	f.hotReload.SetOnChange(func() {
		log.Println("🔄 Files changed - restart recommended")
	})
	
	f.Use(HotReloadMiddleware(f.hotReload))
	f.hotReload.Start()
}