package main

import (

	"goland_api/pkg/handlers"
	"goland_api/pkg/database"
	"log"
	"os"
	"net/http"
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
	// InitDB
	dataSourceName := os.Getenv("DATABASE_URL")
	database.InitDB(dataSourceName)

	// Регистрация маршрутов
	router := mux.NewRouter()
	// Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Участники
	router.HandleFunc("/api/users", handlers.GetUsers()).Methods("GET")
	router.HandleFunc("/api/users/{id}", handlers.GetUser()).Methods("GET")

	// Кабинет
	router.HandleFunc("/api/auth/info", handlers.AuthMiddleware(handlers.InfoUser())).Methods("GET")
	router.HandleFunc("/api/auth/create", handlers.AuthMiddleware(handlers.CreateUser())).Methods("POST")
	router.HandleFunc("/api/auth/update", handlers.AuthMiddleware(handlers.UpdateUser())).Methods("PUT")
	router.HandleFunc("/api/auth/login", handlers.Login()).Methods("POST")
	router.HandleFunc("/api/auth/refresh", handlers.AuthMiddleware(handlers.Refresh())).Methods("POST")
	//router.HandleFunc("/api/auth", handlers.DeleteUser()).Methods("DELETE")

	// Команды
	router.HandleFunc("/api/teams", handlers.GetTeams()).Methods("GET")
	router.HandleFunc("/api/teams/{id}", handlers.GetTeam()).Methods("GET")
	router.HandleFunc("/api/teams", handlers.AuthMiddleware(handlers.CreateTeam())).Methods("POST")
	router.HandleFunc("/api/teams/{id}", handlers.AuthMiddleware(handlers.UpdateTeam())).Methods("PUT")
	router.HandleFunc("/api/teams/{id}", handlers.AuthMiddleware(handlers.DeleteTeam())).Methods("DELETE")

	// Media
	router.HandleFunc("/api/media/preloader", handlers.AuthMiddleware(handlers.Preloader())).Methods("POST")

	//start server
	log.Fatal(http.ListenAndServe(":8000", handlers.JsonContentTypeMiddleware(router)))
}