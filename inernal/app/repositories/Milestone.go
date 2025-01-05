package repositories

import (
	"GeekReward/inernal/app/models/tables"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MilestoneRepository interface {
	FindByBountyID(bountyID uuid.UUID) ([]tables.Milestone, error)
	FindByID(id uuid.UUID) (*tables.Milestone, error)
	CreateMilestone(milestone *tables.Milestone) error
	DeleteMilestone(milestone *tables.Milestone) error
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
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
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
