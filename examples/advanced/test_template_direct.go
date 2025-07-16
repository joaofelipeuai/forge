package main

import (
	"html/template"
	"log"
	"os"
	"strings"
)

func main() {
	log.Println("🔍 Testando template engine do Go diretamente...")
	
	// Teste 1: Template inline
	log.Println("\n1. Testando template inline...")
	tmpl1, err := template.New("inline").Parse("<h1>{{.Title}}</h1><p>{{.Message}}</p>")
	if err != nil {
		log.Printf("❌ Erro no template inline: %v", err)
	} else {
		log.Println("✅ Template inline criado com sucesso")
		
		data := map[string]interface{}{
			"Title":   "Teste Inline",
			"Message": "Funcionando!",
		}
		
		var buf strings.Builder
		if err := tmpl1.Execute(&buf, data); err != nil {
			log.Printf("❌ Erro ao executar template inline: %v", err)
		} else {
			log.Printf("✅ Template inline executado: %s", buf.String())
		}
	}
	
	// Teste 2: Template de arquivo simples
	log.Println("\n2. Testando template de arquivo...")
	tmpl2, err := template.ParseFiles("templates/simple.html")
	if err != nil {
		log.Printf("❌ Erro ao carregar template de arquivo: %v", err)
	} else {
		log.Println("✅ Template de arquivo carregado com sucesso")
		
		data := map[string]interface{}{
			"Title":   "Teste Arquivo",
			"Message": "Template de arquivo funcionando!",
		}
		
		var buf strings.Builder
		if err := tmpl2.Execute(&buf, data); err != nil {
			log.Printf("❌ Erro ao executar template de arquivo: %v", err)
		} else {
			log.Printf("✅ Template de arquivo executado: %s", buf.String())
		}
	}
	
	// Teste 3: Template com New + ParseFiles (como no Forge)
	log.Println("\n3. Testando template como no Forge...")
	tmpl3 := template.New("simple")
	tmpl3, err = tmpl3.ParseFiles("templates/simple.html")
	if err != nil {
		log.Printf("❌ Erro no template estilo Forge: %v", err)
	} else {
		log.Println("✅ Template estilo Forge criado com sucesso")
		
		data := map[string]interface{}{
			"Title":   "Teste Forge Style",
			"Message": "Template estilo Forge funcionando!",
		}
		
		var buf strings.Builder
		if err := tmpl3.Execute(&buf, data); err != nil {
			log.Printf("❌ Erro ao executar template estilo Forge: %v", err)
		} else {
			log.Printf("✅ Template estilo Forge executado: %s", buf.String())
		}
	}
	
	// Teste 4: Verificar conteúdo do arquivo
	log.Println("\n4. Verificando conteúdo do arquivo simple.html...")
	content, err := os.ReadFile("templates/simple.html")
	if err != nil {
		log.Printf("❌ Erro ao ler arquivo: %v", err)
	} else {
		log.Printf("✅ Conteúdo do arquivo: %s", string(content))
	}
}