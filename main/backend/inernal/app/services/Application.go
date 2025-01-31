package services

import (
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/repositories"
	"errors"
	"github.com/google/uuid"
)

type ApplicationService interface {
	CreateApplication(bountyID uuid.UUID, userID uuid.UUID, note string) error
	GetApplications(bountyID uuid.UUID) ([]tables.Application, error)
	ApproveApplication(applicationID uuid.UUID) error
	RejectApplication(applicationID uuid.UUID) error
	GetPublicApplications(bountyID uuid.UUID) ([]*tables.Application, error)
	HasUserApplied(bountyID uuid.UUID, uid uuid.UUID) (bool, error)
}

type applicationService struct {
	applicationRepo repositories.ApplicationRepository
	bountyRepo      repositories.BountyRepository
}

func (s *applicationService) HasUserApplied(bountyID uuid.UUID, UserID uuid.UUID) (bool, error) {
	return s.applicationRepo.HasUserApplied(bountyID, UserID)
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
func (s *applicationService) CreateApplication(bountyID uuid.UUID, userID uuid.UUID, note string) error {
	// 1. 判断是否已有"pending"/"approved"申请
	hasApplied, err := s.applicationRepo.HasUserApplied(bountyID, userID)
	if err != nil {
		return err
	}
	if hasApplied {
		return errors.New("你已对该悬赏令提交过申请或已被批准，无法再次申请")
	}

	application := &tables.Application{
		BountyID: bountyID,
		UserID:   userID,
		Status:   "pending",
		Note:     note, // 可选
	}
	return s.applicationRepo.Create(application)
}

// GetPublicApplications 获取公开的申请信息
func (s *applicationService) GetPublicApplications(bountyID uuid.UUID) ([]*tables.Application, error) {
	// 只返回 "approved" 状态
	return s.applicationRepo.GetApprovedApplicationsByBountyID(bountyID)
}

// GetApplications 获取指定悬赏令的所有申请
func (s *applicationService) GetApplications(bountyID uuid.UUID) ([]tables.Application, error) {
	return s.applicationRepo.FindAllByBountyID(bountyID)
}

// ApproveApplication 批准申请
func (s *applicationService) ApproveApplication(applicationID uuid.UUID) error {
	// 获取申请信息
	app, err := s.applicationRepo.FindByID(applicationID)
	if err != nil {
		return errors.New("申请不存在")
	}

	if app.Status != "pending" {
		return errors.New("只能批准待处理的申请")
	}

	// 调用仓库层的方法，批准申请并更新悬赏令 receiver_id
	return s.applicationRepo.ApproveApplication(applicationID, app.UserID)
}

// RejectApplication 拒绝申请
func (s *applicationService) RejectApplication(applicationID uuid.UUID) error {
	// 获取申请信息
	app, err := s.applicationRepo.FindByID(applicationID)
	if err != nil {
		return errors.New("申请不存在")
	}

	if app.Status != "pending" {
		return errors.New("只能拒绝待处理的申请")
	}

	// 更新申请状态为 "rejected"
	return s.applicationRepo.UpdateApplicationStatus(applicationID, "rejected")
}
