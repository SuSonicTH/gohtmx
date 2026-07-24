package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"htmx/db"
	"htmx/handlers"
)

// parseTemplates mantido igual...

func main() {
	// Inicializa o banco de dados
	database := db.InitDB()
	defer database.Close()

	tmpl, err := parseTemplates()
	if err != nil {
		log.Fatalf("Erro ao carregar templates: %v", err)
	}

	// Injeta o banco no Handler
	todoHandler := handlers.NewTodoHandler(database, tmpl)

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/dashboard", todoHandler.Index)
	mux.HandleFunc("/todos/create", todoHandler.Create)
	mux.HandleFunc("/todos/toggle", todoHandler.Toggle)
	mux.HandleFunc("/todos/delete", todoHandler.Delete)
	mux.HandleFunc("/", todoHandler.Index)

	fmt.Println("Servidor rodando em http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func parseTemplates() (*template.Template, error) {
	// 1. Registra as funções matemáticas no FuncMap do template
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}

	// 2. Aplica o FuncMap antes de parsear os arquivos
	tmpl := template.New("").Funcs(funcMap)

	err := filepath.Walk("templates", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".html") {
			b, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			name := filepath.Base(path)
			_, err = tmpl.New(name).Parse(string(b))
			if err != nil {
				return err
			}
		}
		return nil
	})

	return tmpl, err
}
