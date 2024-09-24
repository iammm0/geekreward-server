package repositories

import (
	"GeekReward/inernal/app/models/tables"
	"gorm.io/gorm"
)

// ApplicationRepository 定义申请相关的数据库操作接口
type ApplicationRepository interface {
	Create(application *tables.Application) error
	FindAllByBountyID(bountyID uint) ([]tables.Application, error)
	UpdateStatus(applicationID uint, status string) error
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(application *tables.Application) error {
	return r.db.Create(application).Error
}

func (r *applicationRepository) FindAllByBountyID(bountyID uint) ([]tables.Application, error) {
	var applications []tables.Application
	err := r.db.Where("bounty_id = ?", bountyID).Find(&applications).Error
	return applications, err
}

func (r *applicationRepository) UpdateStatus(applicationID uint, status string) error {
	return r.db.Model(&tables.Application{}).Where("id = ?", applicationID).Update("status", status).Error
}
