package repositories

import (
	"GeekReward/main/inernal/app/models/tables"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MilestoneRepository interface {
	// FindByBountyID 根据悬赏令 ID 获取所有相关的里程碑
	FindByBountyID(bountyID uuid.UUID) ([]tables.Milestone, error)
	// FindByID 根据里程碑 ID 获取单个里程碑
	FindByID(id uuid.UUID) (*tables.Milestone, error)
	// CreateMilestone 悬赏令发布者创建里程碑
	CreateMilestone(milestone *tables.Milestone) error
	// DeleteMilestone 悬赏令发布者删除某个里程碑
	DeleteMilestone(milestone *tables.Milestone) error
	// UpdateMilestone 悬赏令发布者更新某个里程碑
	UpdateMilestone(milestone *tables.Milestone) error
}

type milestoneRepository struct {
	db *gorm.DB
}

func NewMilestoneRepository(db *gorm.DB) MilestoneRepository {
	return &milestoneRepository{db: db}
}

func (r *milestoneRepository) FindByBountyID(bountyID uuid.UUID) ([]tables.Milestone, error) {
	var milestones []tables.Milestone
	err := r.db.Where("bounty_id = ?", bountyID).Find(&milestones).Error
	return milestones, err
}

func (r *milestoneRepository) FindByID(id uuid.UUID) (*tables.Milestone, error) {
	var milestone tables.Milestone
	err := r.db.First(&milestone, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// if err != nil {
	// 	if errors.Is(err, gorm.ErrRecordNotFound) {
	// 		return nil, nil
	// 	}
	// 	return nil, err
	// }

	return &milestone, nil
}

func (r *milestoneRepository) CreateMilestone(milestone *tables.Milestone) error {
	return r.db.Create(milestone).Error
}

func (r *milestoneRepository) UpdateMilestone(milestone *tables.Milestone) error {
	return r.db.Save(milestone).Error
}

func (r *milestoneRepository) DeleteMilestone(milestone *tables.Milestone) error {
	return r.db.Delete(milestone).Error
}
