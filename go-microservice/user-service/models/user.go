package models

import "gorm.io/gorm"

// User adalah model database untuk menyimpan data akun pengguna.
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"uniqueIndex"`
	Password string `json:"-"` // Tidak pernah dikembalikan ke client
}

// RegisterRequest adalah body request untuk endpoint /register.
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest adalah body request untuk endpoint /login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse adalah response setelah berhasil login.
type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

// UserResponse adalah data user yang aman untuk dikembalikan ke client (tanpa password).
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
