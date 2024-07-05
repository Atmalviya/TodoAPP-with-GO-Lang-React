package main

import (
	"log"
	"net/http"
	"todo-api/db"
	"todo-api/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db.InitDB()
	defer db.Session.Close()

	router := mux.NewRouter()
	router.HandleFunc("/todos", handlers.CreateTodo).Methods("POST")
	router.HandleFunc("/todos", handlers.GetTodos).Methods("GET")
	router.HandleFunc("/todos/{id}", handlers.GetTodo).Methods("GET")
	router.HandleFunc("/todos/{id}", handlers.UpdateTodo).Methods("PUT")
	router.HandleFunc("/todos/{id}", handlers.DeleteTodo).Methods("DELETE")

	// Setup CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Replace with your frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	}).Handler(router)

	port := ":8000"
	log.Printf("Server running on port %s", port)
	log.Fatal(http.ListenAndServe(port, corsHandler))
}
