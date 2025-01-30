package services

import (
	dtos2 "GeekReward/main/inernal/app/models/dtos"
	"GeekReward/main/inernal/app/models/tables"
	repositories2 "GeekReward/main/inernal/app/repositories"
	"errors"
	"github.com/google/uuid"
	"time"
)

type MilestoneService interface {
	GetMilestonesByBountyID(bountyID uuid.UUID) ([]tables.Milestone, error)

	// CreateMilestone 悬赏令发布者创建里程碑
	CreateMilestone(bountyID uuid.UUID, input dtos2.MilestoneDTO) (*tables.Milestone, error)

	// UpdateMilestone 悬赏令（发布者）更新里程碑
	UpdateMilestone(milestoneID uuid.UUID, input dtos2.MilestoneUpdateDTO) error

	// DeleteMilestone 悬赏令（发布者）删除指定的悬赏令
	DeleteMilestone(milestoneID uuid.UUID) error

	// UpdateMilestoneByReceiver 悬赏令（接受者）更新里程碑（完成度）
	UpdateMilestoneByReceiver(bountyID, milestoneID, userID uuid.UUID, input dtos2.MilestoneUpdateDTO) error
}

type milestoneService struct {
	milestoneRepo repositories2.MilestoneRepository
	bountyRepo    repositories2.BountyRepository
	userRepo      repositories2.UserRepository
}

func NewMilestoneService(milestoneRepo repositories2.MilestoneRepository, bountyRepo repositories2.BountyRepository) MilestoneService {
	return &milestoneService{
		milestoneRepo: milestoneRepo,
		bountyRepo:    bountyRepo,
	}
}

// UpdateMilestoneByReceiver 悬赏零接收者更新里程碑
func (s *milestoneService) UpdateMilestoneByReceiver(bountyID uuid.UUID, milestoneID uuid.UUID, userID uuid.UUID, input dtos2.MilestoneUpdateDTO) error {
	// 检查是否可以找到目标里程碑
	milestone, err := s.milestoneRepo.FindByID(milestoneID)
	if err != nil {
		return err
	}
	if milestone == nil {
		return errors.New("里程碑未找到")
	}

	// 获取并检查悬赏令
	bounty, err := s.bountyRepo.FindBountyByID(bountyID)
	if err != nil {
		return err
	}
	if bounty == nil {
		return errors.New("悬赏令未找到")
	}

	// 检查悬赏令是否有 ReceiverID
	if bounty.ReceiverID == nil {
		return errors.New("该悬赏令尚未被接受")
	}

	// 检查用户是否为悬赏令的接收者
	if *bounty.ReceiverID != userID {
		return errors.New("你没有权限更新此里程碑")
	}

	// 悬赏令接受者仅仅允许更新里程碑字段
	milestone.IsCompleted = input.IsCompleted
	milestone.UpdatedAt = time.Now()

	return s.milestoneRepo.UpdateMilestone(milestone)
}

// GetMilestonesByBountyID 获取指定悬赏令的所有里程碑
func (s *milestoneService) GetMilestonesByBountyID(bountyID uuid.UUID) ([]tables.Milestone, error) {
	// 检查悬赏令是否存在
	bounty, err := s.bountyRepo.FindBountyByID(bountyID)
	if err != nil {
		return nil, err
	}
	if bounty == nil {
		return nil, errors.New("悬赏令未找到")
	}

	// 获取里程碑
	milestones, err := s.milestoneRepo.FindByBountyID(bountyID)
	if err != nil {
		return nil, err
	}
	return milestones, nil
	// return s.milestoneRepo.FindByBountyID(bountyID)
}

// CreateMilestone 创建新的里程碑
func (s *milestoneService) CreateMilestone(bountyID uuid.UUID, input dtos2.MilestoneDTO) (*tables.Milestone, error) {
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

// UpdateMilestone 发布者更新里程碑信息
func (s *milestoneService) UpdateMilestone(milestoneID uuid.UUID, input dtos2.MilestoneUpdateDTO) error {
	// 获取里程碑
	milestone, err := s.milestoneRepo.FindByID(milestoneID)
	if err != nil {
		return err
	}
	if milestone == nil {
		return errors.New("未找到悬赏令")
	}

	// 更新里程碑字段
	milestone.Title = input.Title
	milestone.Description = input.Description
	milestone.DueDate = input.DueDate

	// 发布者也能修改完成状态 (保持可选)
	// milestone.IsCompleted = input.IsCompleted

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
		return errors.New("悬赏令未找到")
	}

	// 删除里程碑
	err = s.milestoneRepo.DeleteMilestone(milestone)
	if err != nil {
		return err
	}

	return nil
}
