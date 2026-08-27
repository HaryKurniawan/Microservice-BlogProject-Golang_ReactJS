package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"user-service/models"
	"user-service/repositories"
	"user-service/utils"

	"golang.org/x/crypto/bcrypt"
)

// UserHandler menyimpan dependency yang dibutuhkan oleh HTTP handlers.
type UserHandler struct {
	Repo *repositories.UserRepository
}

// NewUserHandler membuat instance baru UserHandler.
func NewUserHandler(repo *repositories.UserRepository) *UserHandler {
	return &UserHandler{Repo: repo}
}

// Register menangani pendaftaran akun baru.
// POST /api/users/register
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		jsonError(w, "Name, email, and password are required", http.StatusBadRequest)
		return
	}

	// Hash password sebelum disimpan ke database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonError(w, "Failed to process password", http.StatusInternalServerError)
		return
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := h.Repo.Create(user); err != nil {
		jsonError(w, "Email already registered", http.StatusConflict)
		return
	}

	jsonResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"user": models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

// Login menangani proses autentikasi dan mengembalikan JWT token.
// POST /api/users/login
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.Repo.FindByEmail(req.Email)
	if err != nil {
		jsonError(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Bandingkan password dengan hash di database
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		jsonError(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Name)
	if err != nil {
		jsonError(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	jsonResponse(w, http.StatusOK, models.LoginResponse{
		Token: token,
		User: models.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	})
}

// Profile mengembalikan data profil user yang sedang login.
// GET /api/users/profile  (perlu Authorization header)
func (h *UserHandler) Profile(w http.ResponseWriter, r *http.Request) {
	// Ambil user_id dari header yang diteruskan oleh API Gateway
	userIDStr := r.Header.Get("X-User-ID")
	if userIDStr == "" {
		// Fallback: validasi JWT langsung jika tidak ada header dari Gateway
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			jsonError(w, "Authorization required", http.StatusUnauthorized)
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			jsonError(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}
		claims, err := utils.ValidateToken(parts[1])
		if err != nil {
			jsonError(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}
		userIDStr = fmt.Sprintf("%d", claims.UserID)
	}

	var userID uint
	fmt.Sscan(userIDStr, &userID)

	user, err := h.Repo.FindByID(userID)
	if err != nil {
		jsonError(w, "User not found", http.StatusNotFound)
		return
	}

	jsonResponse(w, http.StatusOK, models.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

// --- Helper functions ---

func jsonResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}

func jsonError(w http.ResponseWriter, message string, statusCode int) {
	jsonResponse(w, statusCode, map[string]string{"error": message})
}
