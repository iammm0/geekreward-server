package services

import (
	"GeekReward/inernal/app/models/dtos"
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/repositories"
	"errors"
	"github.com/google/uuid"
	"time"
)

type MilestoneService interface {
	GetMilestonesByBountyID(bountyID uuid.UUID) ([]tables.Milestone, error)
	CreateMilestone(input dtos.MilestoneDTO, bountyID uuid.UUID) (*tables.Milestone, error)
	UpdateMilestone(milestoneID uuid.UUID, input dtos.MilestoneUpdateDTO) error
	DeleteMilestone(milestoneID uuid.UUID) error
	UpdateMilestoneByReceiver(bountyID uuid.UUID, milestoneID uuid.UUID, userID uuid.UUID, input dtos.MilestoneUpdateDTO) error
}

type milestoneService struct {
	milestoneRepo repositories.MilestoneRepository
	bountyRepo    repositories.BountyRepository
	userRepo      repositories.UserRepository
}

func NewMilestoneService(milestoneRepo repositories.MilestoneRepository, bountyRepo repositories.BountyRepository) MilestoneService {
	return &milestoneService{
		milestoneRepo: milestoneRepo,
		bountyRepo:    bountyRepo,
	}
}

// UpdateMilestoneByReceiver 悬赏零接收者更新里程碑
func (s *milestoneService) UpdateMilestoneByReceiver(bountyID uuid.UUID, milestoneID uuid.UUID, userID uuid.UUID, input dtos.MilestoneUpdateDTO) error {
	// 获取里程碑
	milestone, err := s.milestoneRepo.FindByBountyID(milestoneID)
	if err != nil {
		return err
	}
	if milestone == nil {
		return errors.New("milestone not found")
	}

	// 获取悬赏令
	bounty, err := s.bountyRepo.FindBountyByID(bountyID)
	if err != nil {
		return err
	}
	if bounty == nil {
		return errors.New("bounty not found")
	}

	// 检查悬赏令是否有 ReceiverID
	if bounty.ReceiverID == nil {
		return errors.New("bounty has no receiver assigned")
	}

	// 检查用户是否为悬赏令的接收者
	if *bounty.ReceiverID != userID {
		return errors.New("you are not authorized to update this milestone")
	}

	// 更新里程碑字段
	// milestone.Title = input.Title
	// milestone.Description = input.Description
	// milestone.DueDate = input.DueDate
	// milestone.IsCompleted = input.IsCompleted
	// milestone.UpdatedAt = time.Now()

	// 保存更新
	// err = s.milestoneRepo.UpdateMilestone(milestone)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// GetMilestonesByBountyID 获取指定悬赏令的所有里程碑
func (s *milestoneService) GetMilestonesByBountyID(bountyID uuid.UUID) ([]tables.Milestone, error) {
	// 检查悬赏令是否存在
	bounty, err := s.bountyRepo.FindBountyByID(bountyID)
	if err != nil {
		return nil, err
	}
	if bounty == nil {
		return nil, errors.New("bounty not found")
	}

	// 获取里程碑
	milestones, err := s.milestoneRepo.FindByBountyID(bountyID)
	if err != nil {
		return nil, err
	}

	return milestones, nil
}

// CreateMilestone 创建新的里程碑
func (s *milestoneService) CreateMilestone(input dtos.MilestoneDTO, bountyID uuid.UUID) (*tables.Milestone, error) {
	// 检查关联的悬赏令是否存在
	bounty, err := s.bountyRepo.FindBountyByID(bountyID)
	if err != nil {
		return nil, err
	}
	if bounty == nil {
		return nil, errors.New("bounty not found")
	}

	// 创建里程碑
	milestone := &tables.Milestone{
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
		BountyID:    bountyID,
	}

	// 保存里程碑
	err = s.milestoneRepo.CreateMilestone(milestone)
	if err != nil {
		return nil, err
	}

	return milestone, nil
}

// UpdateMilestone 更新里程碑信息
func (s *milestoneService) UpdateMilestone(milestoneID uuid.UUID, input dtos.MilestoneUpdateDTO) error {
	// 获取里程碑
	milestone, err := s.milestoneRepo.FindByID(milestoneID)
	if err != nil {
		return err
	}
	if milestone == nil {
		return errors.New("milestone not found")
	}

	// 更新里程碑字段
	milestone.Title = input.Title
	milestone.Description = input.Description
	milestone.DueDate = input.DueDate
	milestone.UpdatedAt = time.Now()

	// 保存更新
	err = s.milestoneRepo.UpdateMilestone(milestone)
	if err != nil {
		return err
	}

	return nil
}

// DeleteMilestone 删除里程碑
func (s *milestoneService) DeleteMilestone(milestoneID uuid.UUID) error {
	// 获取里程碑
	milestone, err := s.milestoneRepo.FindByID(milestoneID)
	if err != nil {
		return err
	}
	if milestone == nil {
		return errors.New("milestone not found")
	}

	// 删除里程碑
	err = s.milestoneRepo.DeleteMilestone(milestone)
	if err != nil {
		return err
	}

	return nil
}
