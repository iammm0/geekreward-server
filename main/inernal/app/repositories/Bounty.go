package repositories

import (
	"GeekReward/main/inernal/app/models/dtos"
	tables2 "GeekReward/main/inernal/app/models/tables"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BountyRepository 定义关于 Bounty 的数据访问接口
type BountyRepository interface {
	CreateBounty(bounty *tables2.Bounty) error
	FindBounties(filters dtos.BountyFilter) ([]tables2.Bounty, error)
	FindBountyByID(id uuid.UUID) (*tables2.Bounty, error)
	UpdateBounty(bounty *tables2.Bounty) error
	DeleteBounty(bounty *tables2.Bounty) error
	IncrementField(bountyID uuid.UUID, fieldName string) error
	FindByUserID(userID uuid.UUID) ([]tables2.Bounty, error)
	FindReceivedByUserID(userID uuid.UUID) ([]tables2.Bounty, error)
	GetCommentsByBountyID(bountyID uuid.UUID) ([]tables2.Comment, error)
	AddLike(like *tables2.Like) error
	AddComment(comment *tables2.Comment) error
	AddRating(rating *tables2.Rating) error
	FindByIDWithUsers(id uuid.UUID) (*tables2.Bounty, error)
	IsBountyLikedByUser(userID, bountyID uuid.UUID) (bool, error)
	GetUserBountyRating(userID, bountyID uuid.UUID) (float64, error)
	RemoveLike(userID, bountyID uuid.UUID) error
	DecrementField(bountyID uuid.UUID, field string) error
	AddOrUpdateRating(rating *tables2.Rating) error
	GetAllRatingsForBounty(bountyID uuid.UUID) ([]float64, error)
	GetRatingByUserAndBounty(userID, bountyID uuid.UUID, rating *tables2.Rating) error
	UpdateRating(rating *tables2.Rating) error
	GetRatingsByBountyID(bountyID uuid.UUID, ratings *[]tables2.Rating) error
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

func (r *bountyRepository) FindBounties(filters dtos.BountyFilter) ([]tables2.Bounty, error) {
	var bounties []tables2.Bounty

	query := r.db.Model(&tables2.Bounty{})

	// 如果有 status
	if filters.Status != nil {
		query = query.Where("status = ?", *filters.Status)
	}
	// 如果有 publisherID
	if filters.PublisherID != nil {
		query = query.Where("user_id = ?", *filters.PublisherID)
	}
	// 如果有 receiverID
	if filters.ReceiverID != nil {
		query = query.Where("receiver_id = ?", *filters.ReceiverID)
	}

	// 设置分页
	if filters.Limit > 0 {
		query = query.Limit(filters.Limit)
	}
	if filters.Offset > 0 {
		query = query.Offset(filters.Offset)
	}

	// 执行查询
	err := query.Find(&bounties).Error
	if err != nil {
		return nil, err
	}
	return bounties, nil
}

func (r *bountyRepository) CreateBounty(bounty *tables2.Bounty) error {
	return r.db.Create(bounty).Error
}

func (r *bountyRepository) FindBountyByID(id uuid.UUID) (*tables2.Bounty, error) {
	var bounty tables2.Bounty
	err := r.db.First(&bounty, id).Error
	return &bounty, err
}

func (r *bountyRepository) UpdateBounty(bounty *tables2.Bounty) error {
	return r.db.Save(bounty).Error
}

func (r *bountyRepository) DeleteBounty(bounty *tables2.Bounty) error {
	return r.db.Delete(bounty).Error
}

func (r *bountyRepository) IncrementField(bountyID uuid.UUID, fieldName string) error {
	return r.db.Model(&tables2.Bounty{}).Where("id = ?", bountyID).Update(fieldName, gorm.Expr(fieldName+" + ?", 1)).Error
}

func (r *bountyRepository) FindByUserID(userID uuid.UUID) ([]tables2.Bounty, error) {
	var bounties []tables2.Bounty
	err := r.db.Where("user_id = ?", userID).Find(&bounties).Error
	return bounties, err
}

func (r *bountyRepository) FindReceivedByUserID(userID uuid.UUID) ([]tables2.Bounty, error) {
	var bounties []tables2.Bounty
	err := r.db.Where("receiver_id = ?", userID).Find(&bounties).Error
	return bounties, err
}

func (r *bountyRepository) GetCommentsByBountyID(bountyID uuid.UUID) ([]tables2.Comment, error) {
	var comments []tables2.Comment
	// 使用 Preload("User") 将 comment.User 一并查询
	err := r.db.Where("bounty_id = ?", bountyID).
		Preload("User").
		Order("created_at desc").
		Find(&comments).Error

	return comments, err
}

func (r *bountyRepository) AddLike(like *tables2.Like) error {
	return r.db.Create(like).Error
}

func (r *bountyRepository) AddComment(comment *tables2.Comment) error {
	return r.db.Create(comment).Error
}

func (r *bountyRepository) AddRating(rating *tables2.Rating) error {
	return r.db.Create(rating).Error
}

func (r *bountyRepository) FindByIDWithUsers(id uuid.UUID) (*tables2.Bounty, error) {
	var bounty tables2.Bounty
	err := r.db.Preload("User").Preload("Receiver").First(&bounty, id).Error
	return &bounty, err
}

func (r *bountyRepository) IsBountyLikedByUser(userID, bountyID uuid.UUID) (bool, error) {
	var count int64
	err := r.db.Model(&tables2.Like{}).Where("user_id = ? AND bounty_id = ?", userID, bountyID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *bountyRepository) GetUserBountyRating(userID, bountyID uuid.UUID) (float64, error) {
	var rating tables2.Rating
	err := r.db.Where("user_id = ? AND bounty_id = ?", userID, bountyID).First(&rating).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果没有找到记录，则返回默认评分0，并且不返回错误
		return 0, nil
	}
	return rating.Score, err
}

func (r *bountyRepository) RemoveLike(userID, bountyID uuid.UUID) error {
	return r.db.Where("user_id = ? AND bounty_id = ?", userID, bountyID).Delete(&tables2.Like{}).Error
}

func (r *bountyRepository) DecrementField(bountyID uuid.UUID, field string) error {
	return r.db.Model(&tables2.Bounty{}).Where("id = ?", bountyID).Update(field, gorm.Expr(field+" - ?", 1)).Error
}

func (r *bountyRepository) AddOrUpdateRating(rating *tables2.Rating) error {
	// 如果存在则更新评分，否则插入新的评分
	return r.db.Where("user_id = ? AND bounty_id = ?", rating.UserID, rating.BountyID).
		Assign(rating).
		FirstOrCreate(rating).Error
}

func (r *bountyRepository) GetAllRatingsForBounty(bountyID uuid.UUID) ([]float64, error) {
	var scores []float64
	err := r.db.Model(&tables2.Rating{}).Where("bounty_id = ?", bountyID).Pluck("score", &scores).Error
	return scores, err
}

func (r *bountyRepository) GetRatingByUserAndBounty(userID, bountyID uuid.UUID, rating *tables2.Rating) error {
	return r.db.Where("user_id = ? AND bounty_id = ?", userID, bountyID).First(rating).Error
}

func (r *bountyRepository) UpdateRating(rating *tables2.Rating) error {
	return r.db.Save(rating).Error
}

func (r *bountyRepository) GetRatingsByBountyID(bountyID uuid.UUID, ratings *[]tables2.Rating) error {
	return r.db.Where("bounty_id = ?", bountyID).Find(ratings).Error
}

func (r *bountyRepository) UpdateBountyRating(bountyID uuid.UUID, avgScore float64, reviewCount int) error {
	return r.db.Model(&tables2.Bounty{}).Where("id = ?", bountyID).Updates(map[string]interface{}{
		"average_rating": avgScore,
		"review_count":   reviewCount,
	}).Error
}
