package repositories

import (
	"GeekReward/inernal/app/models/tables"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BountyRepository 定义关于 Bounty 的数据访问接口
type BountyRepository interface {
	CreateBounty(bounty *tables.Bounty) error
	FindAllBounties(limit, offset int) ([]tables.Bounty, error)
	FindBountyByID(id uuid.UUID) (*tables.Bounty, error)
	UpdateBounty(bounty *tables.Bounty) error
	DeleteBounty(bounty *tables.Bounty) error
	IncrementField(bountyID uuid.UUID, fieldName string) error
	FindByUserID(userID uuid.UUID) ([]tables.Bounty, error)
	FindReceivedByUserID(userID uuid.UUID) ([]tables.Bounty, error)
	GetCommentsByBountyID(bountyID uuid.UUID) ([]tables.Comment, error)
	AddLike(like *tables.Like) error
	AddComment(comment *tables.Comment) error
	AddRating(rating *tables.Rating) error
	FindByIDWithUsers(id uuid.UUID) (*tables.Bounty, error)
	IsBountyLikedByUser(userID, bountyID uuid.UUID) (bool, error)
	GetUserBountyRating(userID, bountyID uuid.UUID) (float64, error)
	RemoveLike(userID, bountyID uuid.UUID) error
	DecrementField(bountyID uuid.UUID, field string) error
	AddOrUpdateRating(rating *tables.Rating) error
	GetAllRatingsForBounty(bountyID uuid.UUID) ([]float64, error)
	GetRatingByUserAndBounty(userID, bountyID uuid.UUID, rating *tables.Rating) error
	UpdateRating(rating *tables.Rating) error
	GetRatingsByBountyID(bountyID uuid.UUID, ratings *[]tables.Rating) error
	UpdateBountyRating(bountyID uuid.UUID, avgScore float64, reviewCount int) error
}

// 定义 bountyRepository 对象并介入全局变量 db，在接下来的数据操作方法中实现对数据库操作主体的引用
type bountyRepository struct {
	db *gorm.DB
}

// NewBountyRepository 实现创建 bountyRepository 对象的方法，返回一个基础接口类型为 BountyRepository 的 bountyRepository的对象
func NewBountyRepository(db *gorm.DB) BountyRepository {
	return &bountyRepository{db: db}
}

func (r *bountyRepository) CreateBounty(bounty *tables.Bounty) error {
	return r.db.Create(bounty).Error
}

func (r *bountyRepository) FindAllBounties(limit, offset int) ([]tables.Bounty, error) {
	var bounties []tables.Bounty
	err := r.db.Limit(limit).Offset(offset).Find(&bounties).Error
	return bounties, err
}

func (r *bountyRepository) FindBountyByID(id uuid.UUID) (*tables.Bounty, error) {
	var bounty tables.Bounty
	err := r.db.First(&bounty, id).Error
	return &bounty, err
}

func (r *bountyRepository) UpdateBounty(bounty *tables.Bounty) error {
	return r.db.Save(bounty).Error
}

func (r *bountyRepository) DeleteBounty(bounty *tables.Bounty) error {
	return r.db.Delete(bounty).Error
}

func (r *bountyRepository) IncrementField(bountyID uuid.UUID, fieldName string) error {
	return r.db.Model(&tables.Bounty{}).Where("id = ?", bountyID).Update(fieldName, gorm.Expr(fieldName+" + ?", 1)).Error
}

func (r *bountyRepository) FindByUserID(userID uuid.UUID) ([]tables.Bounty, error) {
	var bounties []tables.Bounty
	err := r.db.Where("user_id = ?", userID).Find(&bounties).Error
	return bounties, err
}

func (r *bountyRepository) FindReceivedByUserID(userID uuid.UUID) ([]tables.Bounty, error) {
	var bounties []tables.Bounty
	err := r.db.Where("receiver_id = ?", userID).Find(&bounties).Error
	return bounties, err
}

func (r *bountyRepository) GetCommentsByBountyID(bountyID uuid.UUID) ([]tables.Comment, error) {
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

func (r *bountyRepository) FindByIDWithUsers(id uuid.UUID) (*tables.Bounty, error) {
	var bounty tables.Bounty
	err := r.db.Preload("User").Preload("Receiver").First(&bounty, id).Error
	return &bounty, err
}

func (r *bountyRepository) IsBountyLikedByUser(userID, bountyID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&tables.Like{}).Where("user_id = ? AND bounty_id = ?", userID, bountyID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *bountyRepository) GetUserBountyRating(userID, bountyID uuid.UUID) (float64, error) {
	var rating tables.Rating
	err := r.db.Where("user_id = ? AND bounty_id = ?", userID, bountyID).First(&rating).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果没有找到记录，则返回默认评分0，并且不返回错误
		return 0, nil
	}
	return rating.Score, err
}

func (r *bountyRepository) RemoveLike(userID, bountyID uuid.UUID) error {
	return r.db.Where("user_id = ? AND bounty_id = ?", userID, bountyID).Delete(&tables.Like{}).Error
}

func (r *bountyRepository) DecrementField(bountyID uuid.UUID, field string) error {
	return r.db.Model(&tables.Bounty{}).Where("id = ?", bountyID).Update(field, gorm.Expr(field+" - ?", 1)).Error
}

func (r *bountyRepository) AddOrUpdateRating(rating *tables.Rating) error {
	// 如果存在则更新评分，否则插入新的评分
	return r.db.Where("user_id = ? AND bounty_id = ?", rating.UserID, rating.BountyID).
		Assign(rating).
		FirstOrCreate(rating).Error
}

func (r *bountyRepository) GetAllRatingsForBounty(bountyID uuid.UUID) ([]float64, error) {
	var scores []float64
	err := r.db.Model(&tables.Rating{}).Where("bounty_id = ?", bountyID).Pluck("score", &scores).Error
	return scores, err
}

func (r *bountyRepository) GetRatingByUserAndBounty(userID, bountyID uuid.UUID, rating *tables.Rating) error {
	return r.db.Where("user_id = ? AND bounty_id = ?", userID, bountyID).First(rating).Error
}

func (r *bountyRepository) UpdateRating(rating *tables.Rating) error {
	return r.db.Save(rating).Error
}

func (r *bountyRepository) GetRatingsByBountyID(bountyID uuid.UUID, ratings *[]tables.Rating) error {
	return r.db.Where("bounty_id = ?", bountyID).Find(ratings).Error
}

func (r *bountyRepository) UpdateBountyRating(bountyID uuid.UUID, avgScore float64, reviewCount int) error {
	return r.db.Model(&tables.Bounty{}).Where("id = ?", bountyID).Updates(map[string]interface{}{
		"average_rating": avgScore,
		"review_count":   reviewCount,
	}).Error
}
