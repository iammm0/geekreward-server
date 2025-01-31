package repositories

import (
	"GeekReward/inernal/app/models/tables"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ApplicationRepository interface {
	Create(application *tables.Application) error
	FindAllByBountyID(bountyID uuid.UUID) ([]tables.Application, error)
	UpdateApplicationStatus(applicationID uuid.UUID, status string) error
	GetApprovedApplicationsByBountyID(bountyID uuid.UUID) ([]*tables.Application, error)
	HasUserApplied(bountyID uuid.UUID, UserID uuid.UUID) (bool, error)
	ApproveApplication(applicationID uuid.UUID, receiverID uuid.UUID) error
	FindByID(applicationID uuid.UUID) (*tables.Application, error)
}

type applicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) FindByID(applicationID uuid.UUID) (*tables.Application, error) {
	var app tables.Application
	err := r.db.Where("id = ?", applicationID).Preload("User").First(&app).Error
	if err != nil {
		return nil, err
	}
	return &app, nil
}

func (r *applicationRepository) ApproveApplication(applicationID uuid.UUID, receiverID uuid.UUID) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. 更新申请状态为 "approved"
		if err := tx.Model(&tables.Application{}).
			Where("id = ?", applicationID).
			Update("status", "approved").Error; err != nil {
			return err
		}

		// 2. 更新对应的悬赏令 receiver_id
		// 获取申请信息，确保申请与悬赏令关联
		var app tables.Application
		if err := tx.Where("id = ?", applicationID).First(&app).Error; err != nil {
			return err
		}

		if err := tx.Model(&tables.Bounty{}).
			Where("id = ?", app.BountyID).
			Updates(map[string]interface{}{
				"receiver_id": app.UserID,
				"status":      "assigned", // 可选：更新悬赏令状态
			}).Error; err != nil {
			return err
		}

		// 可选：将其他 pending 的申请状态更新为 "rejected"
		if err := tx.Model(&tables.Application{}).
			Where("bounty_id = ? AND status = ?", app.BountyID, "pending").
			Update("status", "rejected").Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *applicationRepository) HasUserApplied(bountyID uuid.UUID, UserID uuid.UUID) (bool, error) {
	var count int64
	// 查找status in (pending,approved)都算不能再次申请
	err := r.db.Model(&tables.Application{}).
		Where("bounty_id = ? AND user_id = ? AND status IN ?", bountyID, UserID, []string{"pending", "approved"}).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *applicationRepository) Create(application *tables.Application) error {
	return r.db.Create(application).Error
}

func (r *applicationRepository) FindAllByBountyID(bountyID uuid.UUID) ([]tables.Application, error) {
	var apps []tables.Application
	// 预加载User,方便获取app.User.Username, app.User.ProfilePicture
	err := r.db.Where("bounty_id = ?", bountyID).
		Preload("User").
		Order("created_at desc").
		Find(&apps).Error
	return apps, err
}

func (r *applicationRepository) UpdateApplicationStatus(applicationID uuid.UUID, status string) error {
	return r.db.Model(&tables.Application{}).
		Where("id = ?", applicationID).
		Update("status", status).Error
}

func (r *applicationRepository) GetApprovedApplicationsByBountyID(bountyID uuid.UUID) ([]*tables.Application, error) {
	var apps []*tables.Application
	err := r.db.Where("bounty_id = ? AND status = ?", bountyID, "approved").
		Preload("User").
		Order("created_at desc").
		Find(&apps).Error
	return apps, err
}
