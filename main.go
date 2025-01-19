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

	router := mux.NewRouter()
	// Регистрация маршрутов
	// Swagger
	// Устанавливаем маршрут для Swagger UI
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)	// Users
	router.HandleFunc("/api/users", handlers.GetUsers()).Methods("GET")
	router.HandleFunc("/api/users/{id}", handlers.GetUser()).Methods("GET")
	router.HandleFunc("/api/auth/info", handlers.Info()).Methods("GET")
	router.HandleFunc("/api/auth/registration", handlers.Registration()).Methods("POST")
	router.HandleFunc("/api/auth/login", handlers.Login()).Methods("POST")
	router.HandleFunc("/api/auth/refresh", handlers.Refresh()).Methods("POST")

	router.HandleFunc("/api/users/{id}", handlers.UpdateUser()).Methods("PUT")
	router.HandleFunc("/api/users/{id}", handlers.DeleteUser()).Methods("DELETE")

	// Teams
	router.HandleFunc("/api/teams", handlers.GetTeams()).Methods("GET")
	router.HandleFunc("/api/teams/{id}", handlers.GetTeam()).Methods("GET")
	router.HandleFunc("/api/teams", handlers.CreateTeam()).Methods("POST")
	router.HandleFunc("/api/teams/{id}", handlers.UpdateTeam()).Methods("PUT")
	router.HandleFunc("/api/teams/{id}", handlers.DeleteTeam()).Methods("DELETE")

	// Media
	router.HandleFunc("/api/media/preloader", handlers.Preloader()).Methods("POST")

	//start server
	log.Fatal(http.ListenAndServe(":8000", handlers.JsonContentTypeMiddleware(router)))
}