package handlers

import (
	"encoding/json"
	"net/http"
	"post-service/repositories"
	"strconv"

	"github.com/gorilla/mux"
)

// PostHandler menyimpan dependency yang dibutuhkan oleh HTTP handlers.
type PostHandler struct {
	Repo *repositories.PostRepository
}

// NewPostHandler membuat instance baru PostHandler.
func NewPostHandler(repo *repositories.PostRepository) *PostHandler {
	return &PostHandler{Repo: repo}
}

// --- Helper functions ---

func parseID(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	return strconv.Atoi(vars["id"])
}

func jsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, statusCode int) {
	jsonResponse(w, statusCode, map[string]string{"error": message})
}
