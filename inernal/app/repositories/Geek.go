package repositories

import (
	"GeekReward/inernal/app/models/tables"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GeekRepository interface {
	GetTopGeeks(limit int) ([]tables.User, error)
	GetGeekByID(id uuid.UUID) (*tables.User, error)
	GetAffection(userID, geekID uuid.UUID) (*tables.Affection, error)
	CreateAffection(affection *tables.Affection) error
}

type geekRepository struct {
	db *gorm.DB
}

func NewGeekRepository(db *gorm.DB) GeekRepository {
	return &geekRepository{db: db}
}

func (r *geekRepository) GetTopGeeks(limit int) ([]tables.User, error) {
	var geeks []tables.User
	err := r.db.Order("reputation desc").Limit(limit).Find(&geeks).Error
	return geeks, err
}

func (r *geekRepository) GetGeekByID(id uuid.UUID) (*tables.User, error) {
	var geek tables.User
	err := r.db.First(&geek, "id = ?", id).Error
	return &geek, err
}

func (r *geekRepository) GetAffection(userID, geekID uuid.UUID) (*tables.Affection, error) {
	var affection tables.Affection
	if err := r.db.First(&affection, "user_id = ? AND geek_id = ?", userID, geekID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &affection, nil
}

func (r *geekRepository) CreateAffection(affection *tables.Affection) error {
	return r.db.Create(affection).Error
}
