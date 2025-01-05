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
	GetRatingByUserAndBounty(userID, bountyID uuid.UUID, rating *tables.Rating) error
	AddRating(rating *tables.Rating) error
	UpdateRating(rating *tables.Rating) error
	GetApprovedApplicationsByBountyID(bountyID uuid.UUID) ([]*tables.Application, error)
}

type applicationRepository struct {
	db *gorm.DB
}

func (r *applicationRepository) UpdateApplication(app *tables.Application) error {
	//TODO implement me
	panic("implement me")
}

func NewApplicationRepository(db *gorm.DB) ApplicationRepository {
	return &applicationRepository{db: db}
}

func (r *applicationRepository) Create(application *tables.Application) error {
	return r.db.Create(application).Error
}

func (r *applicationRepository) FindAllByBountyID(bountyID uuid.UUID) ([]tables.Application, error) {
	var applications []tables.Application
	err := r.db.Where("bounty_id = ?", bountyID).Find(&applications).Error
	return applications, err
}

func (r *applicationRepository) UpdateApplicationStatus(applicationID uuid.UUID, status string) error {
	return r.db.Model(&tables.Application{}).Where("id = ?", applicationID).Update("status", status).Error
}

// GetRatingByUserAndBounty 假设 RatingRepository 是 ApplicationRepository 的一部分，如果不是，请根据实际情况调整
func (r *applicationRepository) GetRatingByUserAndBounty(userID, bountyID uuid.UUID, rating *tables.Rating) error {
	return r.db.Where("user_id = ? AND bounty_id = ?", userID, bountyID).First(rating).Error
}

func (r *applicationRepository) AddRating(rating *tables.Rating) error {
	return r.db.Create(rating).Error
}

func (r *applicationRepository) UpdateRating(rating *tables.Rating) error {
	return r.db.Save(rating).Error
}

func (r *applicationRepository) GetApprovedApplicationsByBountyID(bountyID uuid.UUID) ([]*tables.Application, error) {
	var applications []*tables.Application
	if err := r.db.Where("bounty_id = ? AND status = ?", bountyID, "Approved").Find(&applications).Error; err != nil {
		return nil, err
	}
	return applications, nil
}
