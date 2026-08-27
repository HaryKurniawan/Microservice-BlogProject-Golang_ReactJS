package middleware

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Claims adalah payload yang dibaca dari JWT token.
type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// ValidateJWT memvalidasi JWT token dari header Authorization.
// Jika valid, meneruskan user_id dan email sebagai header ke downstream service.
// Jika tidak valid, langsung mengembalikan response 401.
// Mengembalikan true jika token valid, false jika tidak.
func ValidateJWT(w http.ResponseWriter, r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		jsonError(w, "Authorization header required", http.StatusUnauthorized)
		return false
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		jsonError(w, "Format: Bearer <token>", http.StatusUnauthorized)
		return false
	}

	secret := os.Getenv("JWT_SECRET")
	token, err := jwt.ParseWithClaims(parts[1], &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		jsonError(w, "Invalid or expired token", http.StatusUnauthorized)
		return false
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		jsonError(w, "Invalid token claims", http.StatusUnauthorized)
		return false
	}

	// Teruskan informasi user ke downstream service via header
	r.Header.Set("X-User-ID", fmt.Sprintf("%d", claims.UserID))
	r.Header.Set("X-User-Email", claims.Email)
	r.Header.Set("X-User-Name", claims.Name)
	return true
}

func jsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	fmt.Fprintf(w, `{"error": "%s"}`, message)
}
