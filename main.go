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

// parseTemplates remains unchanged...

func main() {
	// Initialize the database
	database := db.InitDB()
	defer database.Close()

	tmpl, err := parseTemplates()
	if err != nil {
		log.Fatalf("Error loading templates: %v", err)
	}

	// Inject the database into the handler
	todoHandler := handlers.NewTodoHandler(database, tmpl)

	mux := http.NewServeMux()

	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/dashboard", todoHandler.Index)
	mux.HandleFunc("/todos/create", todoHandler.Create)
	mux.HandleFunc("/todos/toggle", todoHandler.Toggle)
	mux.HandleFunc("/todos/delete", todoHandler.Delete)
	mux.HandleFunc("/", todoHandler.Index)

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}

func parseTemplates() (*template.Template, error) {
	// 1. Register the mathematical functions in the template FuncMap
	funcMap := template.FuncMap{
		"add": func(a, b int) int { return a + b },
		"sub": func(a, b int) int { return a - b },
	}

	// 2. Apply the FuncMap before parsing the files
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
