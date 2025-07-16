package main

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"strings"
	"sync"
)

// TemplateEngine represents the template engine
type TemplateEngine struct {
	templates map[string]*template.Template
	baseDir   string
	extension string
	funcMap   template.FuncMap
	mu        sync.RWMutex
	devMode   bool
}

// NewTemplateEngine creates a new template engine
func NewTemplateEngine(baseDir, extension string) *TemplateEngine {
	return &TemplateEngine{
		templates: make(map[string]*template.Template),
		baseDir:   baseDir,
		extension: extension,
		funcMap:   make(template.FuncMap),
		devMode:   false,
	}
}

// SetDevMode enables/disables development mode (recompiles templates on each request)
func (te *TemplateEngine) SetDevMode(enabled bool) {
	te.devMode = enabled
}

// AddFunc adds a custom function to the template engine
func (te *TemplateEngine) AddFunc(name string, fn interface{}) {
	te.funcMap[name] = fn
}

// LoadTemplates loads all templates from the base directory
func (te *TemplateEngine) LoadTemplates() error {
	te.mu.Lock()
	defer te.mu.Unlock()

	pattern := filepath.Join(te.baseDir, "*."+te.extension)
	templates, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}

	for _, tmplPath := range templates {
		name := strings.TrimSuffix(filepath.Base(tmplPath), "."+te.extension)
		
		tmpl := template.New(name).Funcs(te.funcMap)
		tmpl, err := tmpl.ParseFiles(tmplPath)
		if err != nil {
			return err
		}
		
		te.templates[name] = tmpl
	}

	return nil
}

// Render renders a template with data
func (te *TemplateEngine) Render(w io.Writer, name string, data interface{}) error {
	if te.devMode {
		// In dev mode, reload template on each request
		te.loadSingleTemplate(name)
	}

	te.mu.RLock()
	tmpl, exists := te.templates[name]
	te.mu.RUnlock()

	if !exists {
		return fmt.Errorf("template not found: %s", name)
	}

	return tmpl.Execute(w, data)
}

// loadSingleTemplate loads a single template (used in dev mode)
func (te *TemplateEngine) loadSingleTemplate(name string) error {
	te.mu.Lock()
	defer te.mu.Unlock()

	tmplPath := filepath.Join(te.baseDir, name+"."+te.extension)
	tmpl := template.New(name).Funcs(te.funcMap)
	
	var err error
	tmpl, err = tmpl.ParseFiles(tmplPath)
	if err != nil {
		return err
	}
	
	te.templates[name] = tmpl
	return nil
}

// Template middleware for Forge
func (f *Forge) SetTemplateEngine(engine *TemplateEngine) {
	f.templateEngine = engine
}

// Add template engine to Forge struct
func init() {
	// This will be added to the main Forge struct
}

// Context method for rendering templates
func (c *Context) Render(status int, name string, data interface{}) error {
	// This assumes the template engine is set on the Forge instance
	// We'll need to modify the main forge.go to include this
	c.Response.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.Response.WriteHeader(status)
	
	// For now, return a simple implementation
	// This will be integrated with the main Forge struct
	return nil
}

// Built-in template functions
func DefaultTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
		"join":  strings.Join,
		"split": strings.Split,
		"contains": strings.Contains,
		"replace": strings.Replace,
		"trim": strings.TrimSpace,
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
		"mul": func(a, b int) int { return a * b },
		"div": func(a, b int) int { 
			if b != 0 {
				return a / b
			}
			return 0
		},
		"mod": func(a, b int) int {
			if b != 0 {
				return a % b
			}
			return 0
		},
		"eq":  func(a, b interface{}) bool { return a == b },
		"ne":  func(a, b interface{}) bool { return a != b },
		"lt":  func(a, b int) bool { return a < b },
		"le":  func(a, b int) bool { return a <= b },
		"gt":  func(a, b int) bool { return a > b },
		"ge":  func(a, b int) bool { return a >= b },
		"and": func(a, b bool) bool { return a && b },
		"or":  func(a, b bool) bool { return a || b },
		"not": func(a bool) bool { return !a },
	}
}