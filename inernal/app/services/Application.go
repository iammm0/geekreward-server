package services

import (
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/repositories"
)

// ApplicationService 处理悬赏令申请相关的业务逻辑
type ApplicationService interface {
	CreateApplication(bountyID uint, userID uint) error
	GetApplications(bountyID uint) ([]tables.Application, error)
	ApproveApplication(applicationID uint) error
	RejectApplication(applicationID uint) error
}

type applicationService struct {
	applicationRepo repositories.ApplicationRepository
}

func NewApplicationService(applicationRepo repositories.ApplicationRepository) ApplicationService {
	return &applicationService{applicationRepo: applicationRepo}
}

// CreateApplication 创建新的悬赏令申请
func (s *applicationService) CreateApplication(bountyID uint, userID uint) error {
	// 业务逻辑：检查是否已存在申请、检查悬赏令是否仍开放等
	application := &tables.Application{
		BountyID: bountyID,
		UserID:   userID,
		Status:   "pending",
	}
	return s.applicationRepo.Create(application)
}

// GetApplications 获取指定悬赏令的所有申请
func (s *applicationService) GetApplications(bountyID uint) ([]tables.Application, error) {
	return s.applicationRepo.FindAllByBountyID(bountyID)
}

// ApproveApplication 批准申请
func (s *applicationService) ApproveApplication(applicationID uint) error {
	// 业务逻辑：检查申请状态是否可以批准
	return s.applicationRepo.UpdateStatus(applicationID, "approved")
}

// RejectApplication 拒绝申请
func (s *applicationService) RejectApplication(applicationID uint) error {
	// 业务逻辑：检查申请状态是否可以拒绝
	return s.applicationRepo.UpdateStatus(applicationID, "rejected")
}
