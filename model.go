package main

import (
	"database/sql"
	"fmt"
)

type Note struct {
	ID    int    `json:"id,omitempty"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (n *Note) GetNote(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT title, body FROM notes WHERE id=%d", n.ID)
	return db.QueryRow(statement).Scan(&n.Title, &n.Body)
}

func (n *Note) UpdateNote(db *sql.DB, id int) error {
	statement := fmt.Sprintf("UPDATE notes SET title='%s', body='%s' WHERE id=%d", n.Title, n.Body, id)
	_, err := db.Exec(statement)
	return err
}

func (n *Note) DeleteNote(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM notes WHERE id=%d", n.ID)
	_, err := db.Exec(statement)
	return err
}

func (n *Note) CreateNote(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO notes(title, body) VALUES('%s', '%s')", n.Title, n.Body)
	_, err := db.Exec(statement)
	if err != nil {
		return err
	}
	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&n.ID)
	if err != nil {
		return err
	}
	return nil
}
