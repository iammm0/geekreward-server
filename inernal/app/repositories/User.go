package repositories

import (
	"GeekReward/inernal/app/models/tables"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *tables.User) error
	FindByEmail(email string) (*tables.User, error)
	FindByUsername(username string) (*tables.User, error)
	FindByUserID(id uuid.UUID) (*tables.User, error)
	UpdateUserProfile(user *tables.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *tables.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByEmail(email string) (*tables.User, error) {
	var user tables.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByUsername(username string) (*tables.User, error) {
	var user tables.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByUserID(id uuid.UUID) (*tables.User, error) {
	var user tables.User
	err := r.db.First(&user, "id = ?", id).Error
	return &user, err
}

func (r *userRepository) UpdateUserProfile(user *tables.User) error {
	return r.db.Save(user).Error
}
