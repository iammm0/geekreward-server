package repositories

import (
	"GeekReward/inernal/app/models/tables"
	"gorm.io/gorm"
)

type MilestoneRepository interface {
	FindByBountyID(bountyID uint) ([]tables.Milestone, error)
	FindByID(id uint) (*tables.Milestone, error)
	Create(milestone *tables.Milestone) error
	Update(milestone *tables.Milestone) error
	Delete(milestone *tables.Milestone) error
}

type milestoneRepository struct {
	db *gorm.DB
}

func NewMilestoneRepository(db *gorm.DB) MilestoneRepository {
	return &milestoneRepository{db: db}
}

func (r *milestoneRepository) FindByID(id uint) (*tables.Milestone, error) { // 实现的 FindByID 方法
	var milestone tables.Milestone
	err := r.db.First(&milestone, id).Error
	return &milestone, err
}

func (r *milestoneRepository) FindByBountyID(bountyID uint) ([]tables.Milestone, error) {
	var milestones []tables.Milestone
	err := r.db.Where("bounty_id = ?", bountyID).Find(&milestones).Error
	return milestones, err
}

func (r *milestoneRepository) Create(milestone *tables.Milestone) error {
	return r.db.Create(milestone).Error
}

func (r *milestoneRepository) Update(milestone *tables.Milestone) error {
	return r.db.Save(milestone).Error
}

func (r *milestoneRepository) Delete(milestone *tables.Milestone) error {
	return r.db.Delete(milestone).Error
}
