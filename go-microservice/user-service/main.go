package main

import (
	"log"
	"net/http"
	"os"
	"user-service/config"
	"user-service/handlers"
	"user-service/models"
	"user-service/repositories"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables dari file .env (opsional, bisa dari Docker env)
	godotenv.Load()

	// Inisialisasi koneksi database dan jalankan AutoMigrate
	config.ConnectDatabase()
	if err := config.DB.AutoMigrate(&models.User{}); err != nil {
		log.Fatal("❌ AutoMigrate failed:", err)
	}
	log.Println("✅ AutoMigrate completed")

	// Setup dependencies
	userRepo := repositories.NewUserRepository(config.DB)
	userHandler := handlers.NewUserHandler(userRepo)

	// Setup router
	r := mux.NewRouter()
	api := r.PathPrefix("/api/users").Subrouter()
	api.HandleFunc("/register", userHandler.Register).Methods(http.MethodPost)
	api.HandleFunc("/login", userHandler.Login).Methods(http.MethodPost)
	api.HandleFunc("/profile", userHandler.Profile).Methods(http.MethodGet)

	// Ambil port dari environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("🚀 User Service running on http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
