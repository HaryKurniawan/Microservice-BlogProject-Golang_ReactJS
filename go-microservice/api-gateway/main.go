package main

import (
	"api-gateway/middleware"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	userServiceURL string
	postServiceURL string
)

// newProxy membuat reverse proxy ke URL target.
func newProxy(target string) *httputil.ReverseProxy {
	targetURL, err := url.Parse(target)
	if err != nil {
		log.Fatal("Invalid target URL:", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Hapus header Host agar tidak menyebabkan masalah routing
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = targetURL.Host
	}
	return proxy
}

// isProtectedRoute menentukan apakah route memerlukan JWT authentication.
func isProtectedRoute(method, path string) bool {
	// Request GET ke /api/posts dan /api/posts/{id} adalah PUBLIC
	if method == http.MethodGet && (path == "/api/posts" || strings.HasPrefix(path, "/api/posts/")) {
		return false
	}

	// Route register & login juga PUBLIC
	if path == "/api/users/register" || path == "/api/users/login" {
		return false
	}

	// Semua request lain (seperti POST, PUT, DELETE posts) butuh JWT
	return strings.HasPrefix(path, "/api/users") || strings.HasPrefix(path, "/api/posts")
}

// mainHandler adalah single handler yang menangani semua routing dan proxying.
func mainHandler(w http.ResponseWriter, r *http.Request) {
	// CORS Headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight OPTIONS request
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	path := r.URL.Path

	// Validasi JWT jika route protected
	if isProtectedRoute(r.Method, path) {
		if !middleware.ValidateJWT(w, r) {
			return // Response 401 sudah dikirim oleh middleware
		}
	}

	// Route request ke service yang tepat
	switch {
	case strings.HasPrefix(path, "/api/users"):
		newProxy(userServiceURL).ServeHTTP(w, r)
	case strings.HasPrefix(path, "/api/posts"):
		newProxy(postServiceURL).ServeHTTP(w, r)
	default:
		w.Header().Set("Content-Type", "application/json")
		http.Error(w, `{"error": "Route not found"}`, http.StatusNotFound)
	}
}

func main() {
	godotenv.Load()

	// Ambil URL downstream services dari environment variable
	userServiceURL = os.Getenv("USER_SERVICE_URL")
	postServiceURL = os.Getenv("POST_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://user-service:8081"
	}
	if postServiceURL == "" {
		postServiceURL = "http://post-service:8082"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", mainHandler)

	log.Printf("🚀 API Gateway running on http://localhost:%s", port)
	log.Printf("   → User Service : %s", userServiceURL)
	log.Printf("   → Post Service : %s", postServiceURL)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("Failed to start gateway:", err)
	}
}
