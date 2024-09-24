package repositories

import (
	"GeekReward/inernal/app/models/tables"
	"gorm.io/gorm"
)

type GeekRepository interface {
	GetTopGeeks(limit int) ([]tables.User, error)
	GetGeekByID(id uint) (*tables.User, error)
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

func (r *geekRepository) GetGeekByID(id uint) (*tables.User, error) {
	var geek tables.User
	err := r.db.First(&geek, id).Error
	return &geek, err
}
