package main

import (
    "database/sql"
    "fmt"
)

type note struct {
    ID     int    `json:"id"`
    Title  string `json:"title"`
    Body   string `json:"body"`
}

func (n *note) getNote(db *sql.DB) error {
    statement := fmt.Sprintf("SELECT title, body FROM notes WHERE id=%d", n.ID)
    return db.QueryRow(statement).Scan(&n.Title, &n.Body)
}

func (n *note) updateNote(db *sql.DB) error {
    statement := fmt.Sprintf("UPDATE notes SET title='%s', body='%s' WHERE id=%d", n.Title, n.Body, n.ID)
    _, err := db.Exec(statement)
    return err
}

func (n *note) deleteNote(db *sql.DB) error {
    statement := fmt.Sprintf("DELETE FROM notes WHERE id=%d", n.ID)
    _, err := db.Exec(statement)
    return err
}

func (n *note) createNote(db *sql.DB) error {
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

func getNotes(db *sql.DB, start, count int) ([]note, error) {
    statement := fmt.Sprintf("SELECT id, title, body FROM notes LIMIT %d OFFSET %d", count, start)
    rows, err := db.Query(statement)
    if err != nil {
        return nil, err
    }
    defer rows.Close()
    notes := []note{}
    for rows.Next() {
        var n note
        if err := rows.Scan(&n.ID, &n.Title, &n.Body); err != nil {
            return nil, err
        }
        notes = append(notes, n)
    }
    return notes, nil
}