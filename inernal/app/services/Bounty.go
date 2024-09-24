package services

import (
	"GeekReward/inernal/app/models/common"
	"GeekReward/inernal/app/models/dtos"
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/repositories"
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

type BountyService interface {
	CreateBounty(input dtos.BountyDTO, userID uint) (*tables.Bounty, error)
	GetBounties(limit, offset int) ([]tables.Bounty, error)
	GetBounty(id uint) (*tables.Bounty, error)
	UpdateBounty(id uint, input dtos.BountyDTO) (*tables.Bounty, error)
	DeleteBounty(id uint) error
	LikeBounty(userID, bountyID uint) error
	CommentOnBounty(userID, bountyID uint, content string) (*tables.Comment, error)
	RateBounty(userID, bountyID uint, score float64) error
	IncrementViewCount(bountyID uint) error
	GetBountiesByUserID(userID uint) ([]tables.Bounty, error)
	GetReceivedBounties(userID uint) ([]tables.Bounty, error)
	GetComments(bountyID uint) ([]tables.Comment, error)
	GetUserBountyInteraction(userID, bountyID uint) (*dtos.BountyInteraction, error)
	UnlikeBounty(userID, bountyID uint) error
}

type bountyService struct {
	bountyRepo repositories.BountyRepository
}

func NewBountyService(bountyRepo repositories.BountyRepository) BountyService {
	return &bountyService{bountyRepo: bountyRepo}
}

func (s *bountyService) CreateBounty(input dtos.BountyDTO, userID uint) (*tables.Bounty, error) {
	deadline, err := time.Parse("2006-01-02", input.Deadline)
	if err != nil {
		log.Printf("Error parsing deadline: %v", err)
		return nil, err
	}
	bounty := &tables.Bounty{
		Title:                   input.Title,
		Description:             input.Description,
		Reward:                  input.Reward,
		Deadline:                deadline,
		DifficultyLevel:         input.DifficultyLevel,
		Category:                input.Category,
		Tags:                    input.Tags,
		Location:                input.Location,
		AttachmentURLs:          input.AttachmentUrls,
		Anonymous:               input.Anonymous,
		Priority:                input.Priority,
		PaymentStatus:           input.PaymentStatus,
		PreferredSolutionType:   input.PreferredSolutionType,
		RequiredSkills:          input.RequiredSkills,
		RequiredExperience:      input.RequiredExperience,
		RequiredCertifications:  input.RequiredCertifications,
		Visibility:              input.Visibility,
		Confidentiality:         input.Confidentiality,
		ContractType:            input.ContractType,
		EstimatedHours:          input.EstimatedHours,
		ToolsRequired:           input.ToolsRequired,
		CommunicationPreference: input.CommunicationPreference,
		FeedbackRequired:        input.FeedbackRequired,
		CompletionCriteria:      input.CompletionCriteria,
		SubmissionGuidelines:    input.SubmissionGuidelines,
		EvaluationCriteria:      input.EvaluationCriteria,
		ReferenceMaterials:      input.ReferenceMaterials,
		ExternalLinks:           input.ExternalLinks,
		AdditionalNotes:         input.AdditionalNotes,
		NDARequired:             input.NdaRequired,
		AcceptanceCriteria:      input.AcceptanceCriteria,
		PaymentMethod:           input.PaymentMethod,
		UserID:                  common.UserID(userID),
	}
	if err := s.bountyRepo.Create(bounty); err != nil {
		return nil, err
	}
	return bounty, nil
}

func (s *bountyService) GetBounties(limit, offset int) ([]tables.Bounty, error) {
	return s.bountyRepo.FindAll(limit, offset)
}

func (s *bountyService) GetBounty(id uint) (*tables.Bounty, error) {
	return s.bountyRepo.FindByID(id)
}

func (s *bountyService) UpdateBounty(id uint, input dtos.BountyDTO) (*tables.Bounty, error) {
	bounty, err := s.bountyRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	// 更新逻辑...
	if err := s.bountyRepo.Update(bounty); err != nil {
		return nil, err
	}
	return bounty, nil
}

func (s *bountyService) DeleteBounty(id uint) error {
	bounty, err := s.bountyRepo.FindByID(id)
	if err != nil {
		return err
	}
	return s.bountyRepo.Delete(bounty)
}

func (s *bountyService) LikeBounty(userID, bountyID uint) error {
	like := &tables.Like{UserID: userID, BountyID: bountyID}
	if err := s.bountyRepo.AddLike(like); err != nil {
		return err
	}
	return s.bountyRepo.IncrementField(bountyID, "likes_count")
}

func (s *bountyService) CommentOnBounty(userID, bountyID uint, content string) (*tables.Comment, error) {
	comment := &tables.Comment{UserID: userID, BountyID: bountyID, Content: content}
	if err := s.bountyRepo.AddComment(comment); err != nil {
		return nil, err
	}
	if err := s.bountyRepo.IncrementField(bountyID, "comments_count"); err != nil {
		return nil, err
	}
	return comment, nil
}

func (s *bountyService) RateBounty(userID, bountyID uint, score float64) error {
	var rating tables.Rating

	// 检查用户是否已经评分
	err := s.bountyRepo.GetRatingByUserAndBounty(userID, bountyID, &rating)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// 如果用户已经评分，更新评分
	if rating.ID != 0 {
		rating.Score = score
		if err := s.bountyRepo.UpdateRating(&rating); err != nil {
			return err
		}
	} else {
		// 否则，添加新评分
		newRating := &tables.Rating{UserID: userID, BountyID: bountyID, Score: score}
		if err := s.bountyRepo.AddRating(newRating); err != nil {
			return err
		}
	}

	// 因为我们已经移除了 UpdateAverageRating 方法，逻辑简化至此即可
	return nil
}

func (s *bountyService) IncrementViewCount(bountyID uint) error {
	return s.bountyRepo.IncrementField(bountyID, "view_count")
}

func (s *bountyService) GetBountiesByUserID(userID uint) ([]tables.Bounty, error) {
	return s.bountyRepo.FindByUserID(userID)
}

func (s *bountyService) GetReceivedBounties(userID uint) ([]tables.Bounty, error) {
	return s.bountyRepo.FindReceivedByUserID(userID)
}

func (s *bountyService) GetComments(bountyID uint) ([]tables.Comment, error) {
	return s.bountyRepo.GetCommentsByBountyID(bountyID)
}

func (s *bountyService) GetUserBountyInteraction(userID, bountyID uint) (*dtos.BountyInteraction, error) {
	liked, err := s.bountyRepo.IsBountyLikedByUser(userID, bountyID)
	if err != nil {
		return nil, err
	}

	// 这里的 GetUserBountyRating 已经处理了 record not found 的情况
	score, err := s.bountyRepo.GetUserBountyRating(userID, bountyID)
	if err != nil {
		return nil, err
	}

	return &dtos.BountyInteraction{Liked: liked, Score: score}, nil
}

func (s *bountyService) UnlikeBounty(userID, bountyID uint) error {
	// 检查用户是否已经点赞
	liked, err := s.bountyRepo.IsBountyLikedByUser(userID, bountyID)
	if err != nil {
		return err
	}

	if !liked {
		return errors.New("用户尚未点赞该悬赏令")
	}

	// 删除点赞记录
	if err := s.bountyRepo.RemoveLike(userID, bountyID); err != nil {
		return err
	}

	// 减少点赞计数
	return s.bountyRepo.DecrementField(bountyID, "likes_count")
}
