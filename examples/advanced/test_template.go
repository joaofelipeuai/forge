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
	
	// Template engine setup
	templateEngine := forge.NewTemplateEngine("templates", "html")
	templateEngine.SetDevMode(true)
	
	// Carregar templates
	if err := templateEngine.LoadTemplates(); err != nil {
		log.Printf("Erro ao carregar templates: %v", err)
	} else {
		log.Println("✅ Templates carregados com sucesso!")
	}
	
	app.SetTemplateEngine(templateEngine)
	
	// Rota de teste simples
	app.GET("/", func(c *forge.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "🔨 Servidor de teste de template",
			"status":  "ok",
		})
	})
	
	// Rota do template
	app.GET("/template", func(c *forge.Context) error {
		data := map[string]interface{}{
			"Title":   "Forge Framework Demo",
			"Message": "Template funcionando perfeitamente!",
			"Time":    time.Now().Format("2006-01-02 15:04:05"),
		}
		
		log.Printf("Tentando renderizar template 'index' com dados: %+v", data)
		return c.Render(200, "index", data)
	})
	
	// Rota para debug do template engine
	app.GET("/debug", func(c *forge.Context) error {
		// Verificar se o template engine está configurado
		if engine := c.Get("template_engine"); engine != nil {
			return c.JSON(200, map[string]interface{}{
				"template_engine": "configurado",
				"dev_mode":        true,
				"base_dir":        "templates",
				"extension":       "html",
			})
		} else {
			return c.JSON(500, map[string]interface{}{
				"error": "Template engine não configurado",
			})
		}
	})
	
	log.Println("🔨 Servidor de teste iniciando na porta :8082")
	log.Println("📋 Rotas disponíveis:")
	log.Println("   GET /         - Status do servidor")
	log.Println("   GET /template - Teste do template")
	log.Println("   GET /debug    - Debug do template engine")
	log.Println("")
	
	log.Fatal(app.Listen(":8082"))
}