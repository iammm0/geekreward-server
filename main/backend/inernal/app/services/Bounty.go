package services

import (
	"GeekReward/inernal/app/models/dtos"
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/repositories"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"time"
)

// BountyService 定义悬赏令相关的服务接口
type BountyService interface {
	CreateBounty(input dtos.BountyDTO, userID uuid.UUID) (*tables.Bounty, error)
	GetBounty(id uuid.UUID) (*tables.Bounty, error)
	UpdateBounty(id uuid.UUID, input dtos.BountyDTO) (*tables.Bounty, error)
	DeleteBounty(id uuid.UUID) error
	LikeBounty(userID, bountyID uuid.UUID) error
	UnlikeBounty(userID, bountyID uuid.UUID) error
	RateBounty(userID, bountyID uuid.UUID, score float64) error
	IncrementViewCount(bountyID uuid.UUID) error
	GetBountiesByUserID(userID uuid.UUID) ([]tables.Bounty, error)
	GetReceivedBounties(userID uuid.UUID) ([]tables.Bounty, error)
	GetUserBountyInteraction(userID, bountyID uuid.UUID) (*dtos.BountyInteraction, error)
	SettleBountyAccounts(bountyID uuid.UUID) error
	ConfirmMilestones(bountyID uuid.UUID, userID uuid.UUID) error
	VerifyMilestones(bountyID, userID uuid.UUID) error
	ApplySettlement(bountyID, userID uuid.UUID) error
	FindBounties(filters dtos.BountyFilter) ([]tables.Bounty, error)
	PostComment(userID, bountyID uuid.UUID, content string) (*tables.Comment, error)
	GetCommentsByBountyID(bountyID uuid.UUID) ([]tables.Comment, error)

	// CancelSettlementByPublisher 发布方取消处于Settling状态的悬赏令
	CancelSettlementByPublisher(bountyID, userID uuid.UUID) error

	// CancelSettlementByReceiver 接收方取消处于Settling状态的悬赏令
	CancelSettlementByReceiver(bountyID, userID uuid.UUID) error
}

// bountyService 是 BountyService 接口的具体实现
type bountyService struct {
	userRepo         repositories.UserRepository
	bountyRepo       repositories.BountyRepository
	applicationRepo  repositories.ApplicationRepository
	notificationRepo repositories.NotificationRepository
	milestoneRepo    repositories.MilestoneRepository
}

// PostComment 用户对某个bounty发表评论
func (s *bountyService) PostComment(userID, bountyID uuid.UUID, content string) (*tables.Comment, error) {
	// 构造 Comment 对象
	comment := &tables.Comment{
		UserID:   userID,
		BountyID: bountyID,
		Content:  content,
	}

	// 1. 创建评论
	if err := s.bountyRepo.AddComment(comment); err != nil {
		return nil, err
	}

	// 2. 自增 Bounty 的 comments_count (可选)
	if err := s.bountyRepo.IncrementField(bountyID, "comments_count"); err != nil {
		return nil, err
	}

	return comment, nil
}

// GetCommentsByBountyID 获取某个bounty下的所有评论(含user信息)
func (s *bountyService) GetCommentsByBountyID(bountyID uuid.UUID) ([]tables.Comment, error) {
	// 直接调用commentRepo.GetCommentsByBountyID
	return s.bountyRepo.GetCommentsByBountyID(bountyID)
}

// NewBountyService 创建一个新的 BountyService 实例
func NewBountyService(
	userRepo repositories.UserRepository,
	bountyRepo repositories.BountyRepository,
	applicationRepo repositories.ApplicationRepository,
	notificationRepo repositories.NotificationRepository,
	milestoneRepo repositories.MilestoneRepository,
) BountyService {
	return &bountyService{
		userRepo:         userRepo,
		bountyRepo:       bountyRepo,
		applicationRepo:  applicationRepo,
		notificationRepo: notificationRepo,
		milestoneRepo:    milestoneRepo,
	}
}

// CancelSettlementByPublisher 发布方取消处于Settling状态的悬赏令
func (s *bountyService) CancelSettlementByPublisher(bountyID, userID uuid.UUID) error {
	bounty, err := s.bountyRepo.FindBountyByID(bountyID)
	if err != nil {
		return err
	}
	if bounty == nil {
		return errors.New("bounty not found")
	}

	// 发布方必须是 bounty.UserID
	if bounty.UserID != userID {
		return errors.New("you are not the publisher of this bounty")
	}

	// 仅当 bounty.Status == Settling 时，才能取消
	if bounty.Status != tables.BountyStatusSettling {
		return fmt.Errorf("bounty is not in settling state, current status: %s", bounty.Status)
	}

	// TODO: 扣除违约金 -> 资金相关处理
	// 例如: payPenalty(bounty.UserID, "publisher", bounty.Reward * 0.1)

	// 状态改为Cancelled
	bounty.Status = tables.BountyStatusCancelled
	bounty.UpdatedAt = time.Now()
	if err := s.bountyRepo.UpdateBounty(bounty); err != nil {
		return err
	}

	return nil
}

// CancelSettlementByReceiver 接收方取消处于Settling状态的悬赏令
func (s *bountyService) CancelSettlementByReceiver(bountyID, userID uuid.UUID) error {
	bounty, err := s.bountyRepo.FindBountyByID(bountyID)
	if err != nil {
		return err
	}
	if bounty == nil {
		return errors.New("bounty not found")
	}

	// 接收方必须是 bounty.ReceiverID
	if bounty.ReceiverID == nil || *bounty.ReceiverID != userID {
		return errors.New("you are not the receiver of this bounty")
	}

	// 必须处于 Settling 状态
	if bounty.Status != tables.BountyStatusSettling {
		return fmt.Errorf("bounty is not in settling state, current status: %s", bounty.Status)
	}

	// TODO: 扣除违约金 -> 资金相关处理
	// 例如: payPenalty(*bounty.ReceiverID, "receiver", bounty.Reward * 0.05)

	// 状态改为Cancelled
	bounty.Status = tables.BountyStatusCancelled
	bounty.UpdatedAt = time.Now()
	if err := s.bountyRepo.UpdateBounty(bounty); err != nil {
		return err
	}

	return nil
}

func (s *bountyService) FindBounties(filters dtos.BountyFilter) ([]tables.Bounty, error) {
	// 可在此进行更多业务检查，如：limit过大，status是否有效枚举等
	return s.bountyRepo.FindBounties(filters)
}

// ConfirmMilestones 接受者确认提交所有里程碑
func (s *bountyService) ConfirmMilestones(bountyID uuid.UUID, userID uuid.UUID) error {
	// 1. 根据悬赏令ID查找悬赏令
	bounty, err := s.bountyRepo.FindBountyByID(bountyID)
	if err != nil {
		return err
	}
	if bounty == nil {
		return errors.New("悬赏令未找到")
	}

	// 2. 检查当前用户是否为接收者
	if bounty.ReceiverID == nil {
		return errors.New("悬赏令并未被接收")
	}
	if *bounty.ReceiverID != userID {
		return errors.New("你并不是该悬赏令的接收者")
	}

	// 3. 获取该悬赏令下的所有里程碑
	milestones, err := s.milestoneRepo.FindByBountyID(bountyID)
	if err != nil {
		return err
	}
	if len(milestones) == 0 {
		return errors.New("该悬赏令下没有对应的里程碑")
	}

	// 标记里程碑为完成
	for _, m := range milestones {
		if !m.IsCompleted {
			m.IsCompleted = true
			if err := s.milestoneRepo.UpdateMilestone(&m); err != nil {
				return err
			}
		}
	}
	// 更新 BountyStatus -> MilestonesConfirmed
	if bounty.Status == tables.BountyStatusCreated {
		bounty.Status = tables.BountyStatusMilestonesConfirmed
		bounty.UpdatedAt = time.Now()
		if err := s.bountyRepo.UpdateBounty(bounty); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("当前 bounty status 字段并不是一个可以有效转换悬赏令转台的状态字段, 当前状态: %v", bounty.Status)
	}

	return nil
}

// VerifyMilestones 发布者审核并确认所有里程碑
func (s *bountyService) VerifyMilestones(bountyID, userID uuid.UUID) error {
	// 1. 获取悬赏令
	bounty, err := s.bountyRepo.FindBountyByID(bountyID)
	if err != nil {
		return err
	}
	if bounty == nil {
		return errors.New("悬赏令未找到")
	}

	// 2. 检查当前用户是否为发布者
	if bounty.UserID != userID {
		return errors.New("你不是该悬赏令的发布者")
	}

	// 3. 获取该悬赏令下的所有里程碑
	milestones, err := s.milestoneRepo.FindByBountyID(bountyID)
	if err != nil {
		return err
	}
	if len(milestones) == 0 {
		return errors.New("该悬赏令下没有找到对应的里程碑")
	}

	// 标记里程碑为完成
	for _, m := range milestones {
		if !m.IsCompleted {
			m.IsCompleted = true
			if err := s.milestoneRepo.UpdateMilestone(&m); err != nil {
				return err
			}
		}
	}

	// 5. 所有里程碑完成 -> 发布者确认
	// 更新 BountyStatus -> MilestonesConfirmed
	if bounty.Status == tables.BountyStatusCreated {
		bounty.Status = tables.BountyStatusMilestonesConfirmed
		bounty.UpdatedAt = time.Now()
		if err := s.bountyRepo.UpdateBounty(bounty); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("当前 bounty status 字段并不是一个可以有效转换悬赏令转台的状态字段, 当前状态: %v", bounty.Status)
	}
	return nil
}

// ApplySettlement 接收者申请悬赏令清算
func (s *bountyService) ApplySettlement(bountyID, userID uuid.UUID) error {
	// 1. 获取悬赏令
	bounty, err := s.bountyRepo.FindBountyByID(bountyID)
	if err != nil {
		return err
	}
	if bounty == nil {
		return errors.New("悬赏令未找到")
	}

	// 2. 检查当前用户是否为接收者
	if bounty.ReceiverID == nil {
		return errors.New("悬赏令并未被接收")
	}
	if *bounty.ReceiverID != userID {
		return errors.New("你不是该悬赏令的接收者")
	}

	// 检查 BountyStatus 是否已 Verify
	if bounty.Status == tables.BountyStatusMilestonesVerified {
		// 进入清算状态
		bounty.Status = tables.BountyStatusSettling
		bounty.UpdatedAt = time.Now()
		if err := s.bountyRepo.UpdateBounty(bounty); err != nil {
			return err
		}
		// 若需要实际资金结算逻辑，可在此调用
		// e.g. err = doFinancialSettlement(bounty)
	} else {
		return fmt.Errorf("当前 bounty status 字段并不是一个可以有效转换悬赏令转台的状态字段, 当前状态: %v", bounty.Status)
	}
	return nil
}

// CreateBounty 创建一个新的悬赏令
func (s *bountyService) CreateBounty(input dtos.BountyDTO, userID uuid.UUID) (*tables.Bounty, error) {
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
		UserID:                  userID,
	}

	if err := s.bountyRepo.CreateBounty(bounty); err != nil {
		return nil, err
	}
	return bounty, nil
}

// GetBounty 根据 ID 获取悬赏令
func (s *bountyService) GetBounty(id uuid.UUID) (*tables.Bounty, error) {
	return s.bountyRepo.FindBountyByID(id)
}

// UpdateBounty 更新指定 ID 的悬赏令
func (s *bountyService) UpdateBounty(id uuid.UUID, input dtos.BountyDTO) (*tables.Bounty, error) {
	bounty, err := s.bountyRepo.FindBountyByID(id)
	if err != nil {
		return nil, err
	}
	if bounty == nil {
		return nil, errors.New("bounty not found")
	}

	// 更新字段
	if input.Title != "" {
		bounty.Title = input.Title
	}
	if input.Description != "" {
		bounty.Description = input.Description
	}
	if input.Reward >= 0 {
		bounty.Reward = input.Reward
	}
	if input.Deadline != "" {
		deadline, err := time.Parse("2006-01-02", input.Deadline)
		if err != nil {
			log.Printf("Error parsing deadline: %v", err)
			return nil, err
		}
		bounty.Deadline = deadline
	}
	// 继续更新其他字段...

	bounty.UpdatedAt = time.Now()

	if err := s.bountyRepo.UpdateBounty(bounty); err != nil {
		return nil, err
	}
	return bounty, nil
}

// DeleteBounty 删除指定 ID 的悬赏令
func (s *bountyService) DeleteBounty(id uuid.UUID) error {
	bounty, err := s.bountyRepo.FindBountyByID(id)
	if err != nil {
		return err
	}
	if bounty == nil {
		return errors.New("bounty not found")
	}
	return s.bountyRepo.DeleteBounty(bounty)
}

// LikeBounty 用户点赞悬赏令
func (s *bountyService) LikeBounty(userID, bountyID uuid.UUID) error {
	like := &tables.Like{UserID: userID, BountyID: bountyID}
	if err := s.bountyRepo.AddLike(like); err != nil {
		return err
	}
	return s.bountyRepo.IncrementField(bountyID, "likes_count")
}

// UnlikeBounty 用户取消点赞悬赏令
func (s *bountyService) UnlikeBounty(userID, bountyID uuid.UUID) error {
	liked, err := s.bountyRepo.IsBountyLikedByUser(userID, bountyID)
	if err != nil {
		return err
	}
	if !liked {
		return errors.New("用户尚未点赞该悬赏令")
	}
	if err := s.bountyRepo.RemoveLike(userID, bountyID); err != nil {
		return err
	}
	return s.bountyRepo.DecrementField(bountyID, "likes_count")
}

// RateBounty 用户评分悬赏令
func (s *bountyService) RateBounty(userID, bountyID uuid.UUID, score float64) error {
	var rating tables.Rating

	// 检查用户是否已经评分
	err := s.bountyRepo.GetRatingByUserAndBounty(userID, bountyID, &rating)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if rating.BountyID != uuid.Nil {
		// 更新评分
		rating.Score = score
		if err := s.bountyRepo.UpdateRating(&rating); err != nil {
			return err
		}
	} else {
		// 添加新评分
		newRating := &tables.Rating{UserID: userID, BountyID: bountyID, Score: score}
		if err := s.bountyRepo.AddRating(newRating); err != nil {
			return err
		}
	}

	return nil
}

// IncrementViewCount 增加悬赏令的浏览次数
func (s *bountyService) IncrementViewCount(bountyID uuid.UUID) error {
	return s.bountyRepo.IncrementField(bountyID, "view_count")
}

// GetBountiesByUserID 获取用户发布的所有悬赏令
func (s *bountyService) GetBountiesByUserID(userID uuid.UUID) ([]tables.Bounty, error) {
	return s.bountyRepo.FindByUserID(userID)
}

// GetReceivedBounties 获取用户接收的所有悬赏令
func (s *bountyService) GetReceivedBounties(userID uuid.UUID) ([]tables.Bounty, error) {
	return s.bountyRepo.FindReceivedByUserID(userID)
}

// GetUserBountyInteraction 获取用户在指定悬赏令上的互动信息
func (s *bountyService) GetUserBountyInteraction(userID, bountyID uuid.UUID) (*dtos.BountyInteraction, error) {
	liked, err := s.bountyRepo.IsBountyLikedByUser(userID, bountyID)
	if err != nil {
		return nil, err
	}

	score, err := s.bountyRepo.GetUserBountyRating(userID, bountyID)
	if err != nil {
		return nil, err
	}

	return &dtos.BountyInteraction{Liked: liked, Score: score}, nil
}

// internal/app/services/bounty_service.go

func (s *bountyService) SettleBountyAccounts(bountyID uuid.UUID) error {
	// 获取悬赏令
	bounty, err := s.bountyRepo.FindBountyByID(bountyID)
	if err != nil {
		return err
	}
	if bounty == nil {
		return errors.New("bounty not found")
	}

	// 检查悬赏令状态是否允许结算
	if bounty.PaymentStatus != "Pending" {
		return errors.New("bounty is not in a state that can be settled")
	}

	// 获取所有通过的申请
	approvedApplications, err := s.applicationRepo.GetApprovedApplicationsByBountyID(bountyID)
	if err != nil {
		return err
	}

	if len(approvedApplications) == 0 {
		return errors.New("no approved applications to settle")
	}

	// 计算每个申请者应得的奖励
	totalReward := bounty.Reward
	rewardPerApplication := totalReward / float64(len(approvedApplications))

	// 分发奖励
	for _, app := range approvedApplications {
		// 更新申请状态为已结算
		app.Status = "Settled"
		if err := s.applicationRepo.UpdateApplicationStatus(app.ID, app.Status); err != nil {
			return err
		}

		// 发送通知给申请者
		notification := &tables.Notification{
			UserID:      app.UserID,
			Title:       "Bounty Settle",
			Description: "Your application for bounty '" + bounty.Title + "' has been settled. Reward: $" + fmt.Sprintf("%.2f", rewardPerApplication),
			IsRead:      false,
		}
		if err := s.notificationRepo.CreateNotification(notification); err != nil {
			return err
		}

		// 记录奖励分发（此处假设有一个函数处理资金转移）
		// err = TransferFunds(app.UserID, rewardPerApplication)
		// if err != nil {
		//     return err
		// }
	}

	// 更新悬赏令的支付状态
	bounty.PaymentStatus = "Completed"
	if err := s.bountyRepo.UpdateBounty(bounty); err != nil {
		return err
	}

	// 更新状态为Settled
	bounty.Status = tables.BountyStatusSettled
	bounty.UpdatedAt = time.Now()
	if err := s.bountyRepo.UpdateBounty(bounty); err != nil {
		return err
	}

	return nil
}
