package repositories

import (
	"user-service/models"

	"gorm.io/gorm"
)

// UserRepository menangani semua operasi database yang berkaitan dengan User.
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository membuat instance baru UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create menyimpan user baru ke database.
func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

// FindByEmail mencari user berdasarkan email. Digunakan saat login.
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID mencari user berdasarkan ID. Digunakan untuk endpoint /profile.
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
