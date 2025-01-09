package main

import (
	"goland_api/pkg/handlers"

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
	// Users
	router.HandleFunc("/api/users", handlers.GetUsers(db)).Methods("GET")
	router.HandleFunc("/api/users/{id}", handlers.GetUser(db)).Methods("GET")
	router.HandleFunc("/api/users", handlers.CreateUser(db)).Methods("POST")
	router.HandleFunc("/api/users/{id}", handlers.UpdateUser(db)).Methods("PUT")
	router.HandleFunc("/api/users/{id}", handlers.DeleteUser(db)).Methods("DELETE")

	// Teams
	router.HandleFunc("/api/teams", handlers.GetTeams(db)).Methods("GET")
	router.HandleFunc("/api/teams/{id}", handlers.GetTeam(db)).Methods("GET")
	router.HandleFunc("/api/teams", handlers.CreateTeam(db)).Methods("POST")
	router.HandleFunc("/api/teams/{id}", handlers.UpdateTeam(db)).Methods("PUT")
	router.HandleFunc("/api/teams/{id}", handlers.DeleteTeam(db)).Methods("DELETE")

	//start server
	log.Fatal(http.ListenAndServe(":8000", handlers.JsonContentTypeMiddleware(router)))
}