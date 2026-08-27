package models

import (
	"time"

	"gorm.io/gorm"
)

// Post adalah model database untuk menyimpan data artikel/posting.
type Post struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Title       string         `json:"title"`
	Content     string         `json:"content"`
	AuthorID    uint           `json:"author_id"`    // ID user dari user-service, diteruskan via API Gateway
	AuthorName  string         `json:"author_name"`  // Nama pembuat post, disalin saat submit
	AuthorEmail string         `json:"author_email"` // Email pembuat post, disalin saat submit
}

// CreatePostRequest adalah body request untuk membuat post baru.
type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdatePostRequest adalah body request untuk mengupdate post.
type UpdatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}
