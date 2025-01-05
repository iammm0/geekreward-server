package services

import (
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/repositories"
	"errors"
	"github.com/google/uuid"
)

type ApplicationService interface {
	CreateApplication(bountyID uuid.UUID, userID uuid.UUID) error
	GetApplications(bountyID uuid.UUID) ([]tables.Application, error)
	ApproveApplication(applicationID uuid.UUID) error
	RejectApplication(applicationID uuid.UUID) error
	GetPublicApplications(bountyID uuid.UUID) ([]*tables.Application, error)
}

type applicationService struct {
	applicationRepo repositories.ApplicationRepository
	bountyRepo      repositories.BountyRepository
}

// GetPublicApplications 获取公开的申请信息
func (s *applicationService) GetPublicApplications(bountyID uuid.UUID) ([]*tables.Application, error) {
	// 检查悬赏令是否存在
	bounty, err := s.bountyRepo.FindBountyByID(bountyID)
	if err != nil {
		return nil, err
	}
	if bounty == nil {
		return nil, errors.New("bounty not found")
	}

	// 获取所有公开的申请（假设公开的申请是已批准的）
	applications, err := s.applicationRepo.GetApprovedApplicationsByBountyID(bountyID)
	if err != nil {
		return nil, err
	}

	return applications, nil
}

func NewApplicationService(
	applicationRepo repositories.ApplicationRepository,
	bountyRepo repositories.BountyRepository,
) ApplicationService {
	return &applicationService{
		applicationRepo: applicationRepo,
		bountyRepo:      bountyRepo,
	}
}

// CreateApplication 创建新的悬赏令申请
func (s *applicationService) CreateApplication(bountyID uuid.UUID, userID uuid.UUID) error {
	// 业务逻辑：检查是否已存在申请、检查悬赏令是否仍开放等
	application := &tables.Application{
		BountyID: bountyID,
		UserID:   userID,
		Status:   "pending",
	}
	return s.applicationRepo.Create(application)
}

// GetApplications 获取指定悬赏令的所有申请
func (s *applicationService) GetApplications(bountyID uuid.UUID) ([]tables.Application, error) {
	return s.applicationRepo.FindAllByBountyID(bountyID)
}

// ApproveApplication 批准申请
func (s *applicationService) ApproveApplication(applicationID uuid.UUID) error {
	// 业务逻辑：检查申请状态是否可以批准
	return s.applicationRepo.UpdateApplicationStatus(applicationID, "approved")
}

// RejectApplication 拒绝申请
func (s *applicationService) RejectApplication(applicationID uuid.UUID) error {
	// 业务逻辑：检查申请状态是否可以拒绝
	return s.applicationRepo.UpdateApplicationStatus(applicationID, "rejected")
}
