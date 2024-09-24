package services

import (
	"GeekReward/inernal/app/models/dtos"
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/repositories"
	"time"
)

type MilestoneService interface {
	GetMilestonesByBountyID(bountyID uint) ([]tables.Milestone, error)
	CreateMilestone(input *dtos.MilestoneDTO) error
	UpdateMilestone(id uint, input *dtos.MilestoneDTO) error
	DeleteMilestone(id uint) error
}

type milestoneService struct {
	milestoneRepo repositories.MilestoneRepository
}

func NewMilestoneService(milestoneRepo repositories.MilestoneRepository) MilestoneService {
	return &milestoneService{milestoneRepo: milestoneRepo}
}

func (s *milestoneService) GetMilestonesByBountyID(bountyID uint) ([]tables.Milestone, error) {
	return s.milestoneRepo.FindByBountyID(bountyID)
}

func (s *milestoneService) CreateMilestone(input *dtos.MilestoneDTO) error {
	dueDate, err := time.Parse("2006-01-02", input.DueDate)
	if err != nil {
		return err
	}

	milestone := &tables.Milestone{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     dueDate,
		BountyID:    input.BountyID,
	}

	return s.milestoneRepo.Create(milestone)
}

func (s *milestoneService) UpdateMilestone(id uint, input *dtos.MilestoneDTO) error {
	milestone, err := s.milestoneRepo.FindByID(id)
	if err != nil {
		return err
	}

	// 更新里程碑的字段
	if input.Title != "" {
		milestone.Title = input.Title
	}
	if input.Description != "" {
		milestone.Description = input.Description
	}
	if input.DueDate != "" {
		dueDate, err := time.Parse("2006-01-02", input.DueDate)
		if err != nil {
			return err
		}
		milestone.DueDate = dueDate
	}

	return s.milestoneRepo.Update(milestone)
}

func (s *milestoneService) DeleteMilestone(id uint) error {
	milestone, err := s.milestoneRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.milestoneRepo.Delete(milestone)
}
