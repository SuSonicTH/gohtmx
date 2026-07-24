package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB() *sql.DB {
	db, err := sql.Open("sqlite", "todos.db")
	if err != nil {
		log.Fatalf("Erro ao abrir banco de dados: %v", err)
	}

	db.SetMaxOpenConns(1)

	query := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT 0
	);`

	_, err = db.Exec(query)
	if err != nil {
		log.Fatalf("Erro ao criar tabela: %v", err)
	}

	return db
}
