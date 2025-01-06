package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
)

func main() {
	//connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	//create router
	router := mux.NewRouter()
	router.HandleFunc("/api/users", getUsers(db)).Methods("GET")
	router.HandleFunc("/api/users/{id}", getUser(db)).Methods("GET")
	router.HandleFunc("/api/users", createUser(db)).Methods("POST")
	router.HandleFunc("/api/users/{id}", updateUser(db)).Methods("PUT")
	router.HandleFunc("/api/users/{id}", deleteUser(db)).Methods("DELETE")

	//start server
	log.Fatal(http.ListenAndServe(":8000", jsonContentTypeMiddleware(router)))
}