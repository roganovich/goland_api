package main

import (
	"goland_api/pkg/handlers"
	"database/sql"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	_ "goland_api/docs"
)

// @title My Golang API
// @description This is a sample server.
// @version 1.0
// @host localhost:8080
// @BasePath /api
func main() {
	//connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := mux.NewRouter()
	// Регистрация маршрутов
	// Swagger
	// Устанавливаем маршрут для Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)	// Users
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

	// Media
	router.HandleFunc("/api/media/preloader", handlers.Preloader(db)).Methods("POST")

	//start server
	log.Fatal(http.ListenAndServe(":8000", handlers.JsonContentTypeMiddleware(router)))
}