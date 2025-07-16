package main

import (
	"log"
	"os"
	"path/filepath"
	
	"github.com/joaofelipeuai/forge"
)

func main() {
	// Verificar se os arquivos de template existem
	log.Println("🔍 Verificando arquivos de template...")
	
	// Verificar diretório templates
	if _, err := os.Stat("templates"); os.IsNotExist(err) {
		log.Fatal("❌ Diretório 'templates' não existe!")
	}
	log.Println("✅ Diretório 'templates' existe")
	
	// Verificar arquivos específicos
	files := []string{"index.html", "test.html"}
	for _, file := range files {
		path := filepath.Join("templates", file)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			log.Printf("❌ Arquivo '%s' não existe!", path)
		} else {
			log.Printf("✅ Arquivo '%s' existe", path)
		}
	}
	
	// Testar template engine
	log.Println("\n🔧 Testando Template Engine...")
	
	templateEngine := forge.NewTemplateEngine("templates", "html")
	templateEngine.SetDevMode(true)
	
	// Carregar templates
	if err := templateEngine.LoadTemplates(); err != nil {
		log.Printf("❌ Erro ao carregar templates: %v", err)
		return
	}
	log.Println("✅ Templates carregados com sucesso!")
	
	// Criar aplicação simples
	app := forge.New()
	app.Use(forge.Recovery())
	app.Use(forge.Logger())
	
	app.SetTemplateEngine(templateEngine)
	
	// Rota de debug
	app.GET("/debug", func(c *forge.Context) error {
		return c.JSON(200, map[string]interface{}{
			"message": "Debug endpoint funcionando",
			"template_engine_set": c.Get("template_engine") != nil,
		})
	})
	
	// Rota de teste simples
	app.GET("/test", func(c *forge.Context) error {
		log.Println("🎯 Tentando renderizar template 'test'...")
		
		data := map[string]interface{}{
			"Title":   "Teste Simples",
			"Message": "Template de teste",
			"Time":    "2024-01-01 12:00:00",
		}
		
		return c.Render(200, "test", data)
	})
	
	// Rota de teste ultra-simples
	app.GET("/simple", func(c *forge.Context) error {
		log.Println("🎯 Tentando renderizar template 'simple'...")
		
		data := map[string]interface{}{
			"Title":   "Template Ultra Simples",
			"Message": "Se isso funcionar, o problema está no template complexo",
		}
		
		return c.Render(200, "simple", data)
	})
	
	// Rota HTML direta (sem template)
	app.GET("/direct", func(c *forge.Context) error {
		html := `
<!DOCTYPE html>
<html>
<head>
    <title>Teste Direto</title>
</head>
<body>
    <h1>🔨 Forge Framework</h1>
    <p>HTML direto funcionando!</p>
</body>
</html>`
		return c.HTML(200, html)
	})
	
	log.Println("\n🚀 Servidor debug iniciando na porta :8084")
	log.Println("📋 Rotas disponíveis:")
	log.Println("   GET /debug  - Debug info")
	log.Println("   GET /test   - Teste template")
	log.Println("   GET /simple - Template ultra-simples")
	log.Println("   GET /direct - HTML direto")
	log.Println("")
	
	log.Fatal(app.Listen(":9001"))
}