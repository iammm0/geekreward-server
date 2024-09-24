package repositories

import (
	"GeekReward/inernal/app/models/tables"
	"errors"
	"gorm.io/gorm"
)

type BountyRepository interface {
	Create(bounty *tables.Bounty) error
	FindAll(limit, offset int) ([]tables.Bounty, error)
	FindByID(id uint) (*tables.Bounty, error)
	Update(bounty *tables.Bounty) error
	Delete(bounty *tables.Bounty) error
	IncrementField(bountyID uint, fieldName string) error
	FindByUserID(userID uint) ([]tables.Bounty, error)
	FindReceivedByUserID(userID uint) ([]tables.Bounty, error)
	GetCommentsByBountyID(bountyID uint) ([]tables.Comment, error)
	AddLike(like *tables.Like) error
	AddComment(comment *tables.Comment) error
	AddRating(rating *tables.Rating) error
	FindByIDWithUsers(id uint) (*tables.Bounty, error)
	IsBountyLikedByUser(userID, bountyID uint) (bool, error)
	GetUserBountyRating(userID, bountyID uint) (float64, error)
	RemoveLike(userID, bountyID uint) error
	DecrementField(bountyID uint, field string) error
	AddOrUpdateRating(rating *tables.Rating) error
	GetAllRatingsForBounty(bountyID uint) ([]float64, error)
	GetRatingByUserAndBounty(userID, bountyID uint, rating *tables.Rating) error
	UpdateRating(rating *tables.Rating) error
	GetRatingsByBountyID(bountyID uint, ratings *[]tables.Rating) error
	UpdateBountyRating(bountyID uint, avgScore float64, reviewCount int) error
}

type bountyRepository struct {
	db *gorm.DB
}

func NewBountyRepository(db *gorm.DB) BountyRepository {
	return &bountyRepository{db: db}
}

func (r *bountyRepository) Create(bounty *tables.Bounty) error {
	return r.db.Create(bounty).Error
}

func (r *bountyRepository) FindAll(limit, offset int) ([]tables.Bounty, error) {
	var bounties []tables.Bounty
	err := r.db.Limit(limit).Offset(offset).Find(&bounties).Error
	return bounties, err
}

func (r *bountyRepository) FindByID(id uint) (*tables.Bounty, error) {
	var bounty tables.Bounty
	err := r.db.First(&bounty, id).Error
	return &bounty, err
}

func (r *bountyRepository) Update(bounty *tables.Bounty) error {
	return r.db.Save(bounty).Error
}

func (r *bountyRepository) Delete(bounty *tables.Bounty) error {
	return r.db.Delete(bounty).Error
}

func (r *bountyRepository) IncrementField(bountyID uint, fieldName string) error {
	return r.db.Model(&tables.Bounty{}).Where("id = ?", bountyID).Update(fieldName, gorm.Expr(fieldName+" + ?", 1)).Error
}

func (r *bountyRepository) FindByUserID(userID uint) ([]tables.Bounty, error) {
	var bounties []tables.Bounty
	err := r.db.Where("user_id = ?", userID).Find(&bounties).Error
	return bounties, err
}

func (r *bountyRepository) FindReceivedByUserID(userID uint) ([]tables.Bounty, error) {
	var bounties []tables.Bounty
	err := r.db.Where("receiver_id = ?", userID).Find(&bounties).Error
	return bounties, err
}

func (r *bountyRepository) GetCommentsByBountyID(bountyID uint) ([]tables.Comment, error) {
	var comments []tables.Comment
	err := r.db.Where("bounty_id = ?", bountyID).Order("created_at desc").Find(&comments).Error
	return comments, err
}

func (r *bountyRepository) AddLike(like *tables.Like) error {
	return r.db.Create(like).Error
}

func (r *bountyRepository) AddComment(comment *tables.Comment) error {
	return r.db.Create(comment).Error
}

func (r *bountyRepository) AddRating(rating *tables.Rating) error {
	return r.db.Create(rating).Error
}

func (r *bountyRepository) FindByIDWithUsers(id uint) (*tables.Bounty, error) {
	var bounty tables.Bounty
	err := r.db.Preload("User").Preload("Receiver").First(&bounty, id).Error
	return &bounty, err
}

func (r *bountyRepository) IsBountyLikedByUser(userID, bountyID uint) (bool, error) {
	var count int64
	err := r.db.Model(&tables.Like{}).Where("user_id = ? AND bounty_id = ?", userID, bountyID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *bountyRepository) GetUserBountyRating(userID, bountyID uint) (float64, error) {
	var rating tables.Rating
	err := r.db.Where("user_id = ? AND bounty_id = ?", userID, bountyID).First(&rating).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果没有找到记录，则返回默认评分0，并且不返回错误
		return 0, nil
	}
	return rating.Score, err
}

func (r *bountyRepository) RemoveLike(userID, bountyID uint) error {
	return r.db.Where("user_id = ? AND bounty_id = ?", userID, bountyID).Delete(&tables.Like{}).Error
}

func (r *bountyRepository) DecrementField(bountyID uint, field string) error {
	return r.db.Model(&tables.Bounty{}).Where("id = ?", bountyID).Update(field, gorm.Expr(field+" - ?", 1)).Error
}

func (r *bountyRepository) AddOrUpdateRating(rating *tables.Rating) error {
	// 如果存在则更新评分，否则插入新的评分
	return r.db.Where("user_id = ? AND bounty_id = ?", rating.UserID, rating.BountyID).
		Assign(rating).
		FirstOrCreate(rating).Error
}

func (r *bountyRepository) GetAllRatingsForBounty(bountyID uint) ([]float64, error) {
	var scores []float64
	err := r.db.Model(&tables.Rating{}).Where("bounty_id = ?", bountyID).Pluck("score", &scores).Error
	return scores, err
}

func (r *bountyRepository) GetRatingByUserAndBounty(userID, bountyID uint, rating *tables.Rating) error {
	return r.db.Where("user_id = ? AND bounty_id = ?", userID, bountyID).First(rating).Error
}

func (r *bountyRepository) UpdateRating(rating *tables.Rating) error {
	return r.db.Save(rating).Error
}

func (r *bountyRepository) GetRatingsByBountyID(bountyID uint, ratings *[]tables.Rating) error {
	return r.db.Where("bounty_id = ?", bountyID).Find(ratings).Error
}

func (r *bountyRepository) UpdateBountyRating(bountyID uint, avgScore float64, reviewCount int) error {
	return r.db.Model(&tables.Bounty{}).Where("id = ?", bountyID).Updates(map[string]interface{}{
		"average_rating": avgScore,
		"review_count":   reviewCount,
	}).Error
}
