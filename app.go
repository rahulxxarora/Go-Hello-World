package main

import (
    "database/sql"
    "log"
    "net/http"
    "strconv"
    "encoding/json"
    "fmt"

    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
)

type App struct {
    Router *mux.Router
    DB     *sql.DB
}

func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("%s:%s@/%s", user, password, dbname)
    var err error
    a.DB, err = sql.Open("mysql", connectionString)
    if err != nil {
        log.Fatal(err)
    }
    a.Router = mux.NewRouter()
    a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func respondWithError(w http.ResponseWriter, code int, message string) {
    respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
    response, _ := json.Marshal(payload)
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(code)
    w.Write(response)
}

func (a *App) initializeRoutes() {
    a.Router.HandleFunc("/notes", a.createNote).Methods("POST")
    a.Router.HandleFunc("/notes/{id:[0-9]+}", a.getNote).Methods("GET")
}

func (a *App) createNote(w http.ResponseWriter, r *http.Request) {
    var n note
    decoder := json.NewDecoder(r.Body)
    if err := decoder.Decode(&n); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }
    defer r.Body.Close()
    if err := n.createNote(a.DB); err != nil {
        respondWithError(w, http.StatusInternalServerError, err.Error())
        return
    }
    respondWithJSON(w, http.StatusCreated, n)
}

func (a *App) getNote(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id, err := strconv.Atoi(vars["id"])
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid note ID")
        return
    }
    n := note{ID: id}
    if err := n.getNote(a.DB); err != nil {
        switch err {
        case sql.ErrNoRows:
            respondWithError(w, http.StatusNotFound, "Note not found")
        default:
            respondWithError(w, http.StatusInternalServerError, err.Error())
        }
        return
    }
    respondWithJSON(w, http.StatusOK, n)
}