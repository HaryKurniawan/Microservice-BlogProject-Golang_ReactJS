package handlers

import (
	"encoding/json"
	"net/http"
	"post-service/models"
	"strconv"
)

// Update memperbarui post berdasarkan ID.
// PUT /api/posts/{id}  (perlu JWT — divalidasi di API Gateway)
func (h *PostHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		jsonError(w, "Invalid post ID", http.StatusBadRequest)
		return
	}

	post, err := h.Repo.FindByID(uint(id))
	if err != nil {
		jsonError(w, "Post not found", http.StatusNotFound)
		return
	}

	// VALIDASI: Pastikan user yang mengedit adalah pemilik asli postingan
	authorIDStr := r.Header.Get("X-User-ID")
	userID, _ := strconv.Atoi(authorIDStr)
	if post.AuthorID != uint(userID) {
		jsonError(w, "Forbidden: You are not the owner of this post", http.StatusForbidden)
		return
	}

	var req models.UpdatePostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	if err := h.Repo.Update(post); err != nil {
		jsonError(w, "Failed to update post", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, post)
}
