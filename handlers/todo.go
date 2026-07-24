package handlers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"htmx/models"
)

type TodoHandler struct {
	db   *sql.DB
	tmpl *template.Template
}

func NewTodoHandler(db *sql.DB, tmpl *template.Template) *TodoHandler {
	return &TodoHandler{
		db:   db,
		tmpl: tmpl,
	}
}

func (h *TodoHandler) getAllTodos() ([]models.Todo, error) {
	rows, err := h.db.Query("SELECT id, title, completed FROM todos ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	todos := make([]models.Todo, 0)

	for rows.Next() {
		var t models.Todo
		if err := rows.Scan(&t.ID, &t.Title, &t.Completed); err != nil {
			return nil, err
		}
		todos = append(todos, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func (h *TodoHandler) Index(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	if path != "" && path != "/dashboard" {
		http.NotFound(w, r)
		return
	}

	todos, err := h.getAllTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.tmpl.ExecuteTemplate(w, "layout.html", todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	title := r.FormValue("title")
	if strings.TrimSpace(title) == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec("INSERT INTO todos (title, completed) VALUES (?, ?)", title, false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todos, err := h.getAllTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.tmpl.ExecuteTemplate(w, "todo-list.html", todos)
	if err != nil {
		log.Printf("Erro ao renderizar template todo-list.html: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *TodoHandler) Toggle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	_, err = h.db.Exec("UPDATE todos SET completed = NOT completed WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todos, err := h.getAllTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.tmpl.ExecuteTemplate(w, "todo-list.html", todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	_, err = h.db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todos, err := h.getAllTodos()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = h.tmpl.ExecuteTemplate(w, "todo-list.html", todos)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
