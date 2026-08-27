package handlers

import (
	"net/http"
)

// GetAll mengambil semua post.
// GET /api/posts  (public)
func (h *PostHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	posts, err := h.Repo.FindAll()
	if err != nil {
		jsonError(w, "Failed to fetch posts", http.StatusInternalServerError)
		return
	}
	jsonResponse(w, http.StatusOK, posts)
}

// GetByID mengambil satu post berdasarkan ID.
// GET /api/posts/{id}  (public)
func (h *PostHandler) GetByID(w http.ResponseWriter, r *http.Request) {
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

	jsonResponse(w, http.StatusOK, post)
}
