package handlers

import (
	"encoding/json"
	"net/http"
	"post-service/models"
	"strconv"
)

// Create membuat post baru.
// POST /api/posts  (perlu JWT — divalidasi di API Gateway)
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Author ID diteruskan dari API Gateway setelah JWT divalidasi
	authorIDStr := r.Header.Get("X-User-ID")
	authorID, _ := strconv.Atoi(authorIDStr)
	authorName := r.Header.Get("X-User-Name")
	authorEmail := r.Header.Get("X-User-Email")

	var req models.CreatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Content == "" {
		jsonError(w, "Title and content are required", http.StatusBadRequest)
		return
	}

	post := &models.Post{
		Title:       req.Title,
		Content:     req.Content,
		AuthorID:    uint(authorID),
		AuthorName:  authorName,
		AuthorEmail: authorEmail,
	}

	if err := h.Repo.Create(post); err != nil {
		jsonError(w, "Failed to create post", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusCreated, post)
}
