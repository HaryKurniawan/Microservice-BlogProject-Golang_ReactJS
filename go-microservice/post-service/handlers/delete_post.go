package handlers

import (
	"net/http"
	"strconv"
)

// Delete menghapus post berdasarkan ID.
// DELETE /api/posts/{id}  (perlu JWT — divalidasi di API Gateway)
func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
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

	// VALIDASI: Pastikan user yang menghapus adalah pemilik asli postingan
	authorIDStr := r.Header.Get("X-User-ID")
	userID, _ := strconv.Atoi(authorIDStr)
	if post.AuthorID != uint(userID) {
		jsonError(w, "Forbidden: You are not the owner of this post", http.StatusForbidden)
		return
	}

	if err := h.Repo.Delete(uint(id)); err != nil {
		jsonError(w, "Failed to delete post", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"message": "Post deleted successfully"})
}
