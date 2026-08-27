package main

import (
	"log"
	"net/http"
	"os"
	"post-service/config"
	"post-service/handlers"
	"post-service/models"
	"post-service/repositories"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables dari file .env (opsional, bisa dari Docker env)
	godotenv.Load()

	// Inisialisasi koneksi database dan jalankan AutoMigrate
	config.ConnectDatabase()
	if err := config.DB.AutoMigrate(&models.Post{}); err != nil {
		log.Fatal("❌ AutoMigrate failed:", err)
	}
	log.Println("✅ AutoMigrate completed")

	// Setup dependencies
	postRepo := repositories.NewPostRepository(config.DB)
	postHandler := handlers.NewPostHandler(postRepo)

	// Setup router
	r := mux.NewRouter()
	api := r.PathPrefix("/api/posts").Subrouter()
	api.HandleFunc("", postHandler.GetAll).Methods(http.MethodGet)
	api.HandleFunc("/{id}", postHandler.GetByID).Methods(http.MethodGet)
	api.HandleFunc("", postHandler.Create).Methods(http.MethodPost)
	api.HandleFunc("/{id}", postHandler.Update).Methods(http.MethodPut)
	api.HandleFunc("/{id}", postHandler.Delete).Methods(http.MethodDelete)

	// Ambil port dari environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("🚀 Post Service running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
