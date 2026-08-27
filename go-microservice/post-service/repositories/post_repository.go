package repositories

import (
	"post-service/models"

	"gorm.io/gorm"
)

// PostRepository menangani semua operasi database yang berkaitan dengan Post.
type PostRepository struct {
	DB *gorm.DB
}

// NewPostRepository membuat instance baru PostRepository.
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

// FindAll mengambil semua post dari database.
func (r *PostRepository) FindAll() ([]models.Post, error) {
	var posts []models.Post
	if err := r.DB.Order("created_at desc").Find(&posts).Error; err != nil {
		return nil, err
	}
	return posts, nil
}

// FindByID mencari satu post berdasarkan ID.
func (r *PostRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	if err := r.DB.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

// Create menyimpan post baru ke database.
func (r *PostRepository) Create(post *models.Post) error {
	return r.DB.Create(post).Error
}

// Update menyimpan perubahan post ke database.
func (r *PostRepository) Update(post *models.Post) error {
	return r.DB.Save(post).Error
}

// Delete menghapus post berdasarkan ID (soft delete karena menggunakan gorm.Model).
func (r *PostRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Post{}, id).Error
}
