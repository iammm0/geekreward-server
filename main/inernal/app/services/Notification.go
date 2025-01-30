package services

import (
	"GeekReward/main/inernal/app/models/tables"
	"GeekReward/main/inernal/app/repositories"
	"fmt"
	"github.com/google/uuid"
)

type NotificationService interface {
	CreateNotification(notification *tables.Notification) error
	GetUserNotifications(userID uuid.UUID) ([]tables.Notification, error)
	MarkNotificationAsRead(notificationID uuid.UUID) error
	DeleteNotification(notificationID uuid.UUID) error
	CreateBountyApplicationNotification(applicantID, publisherID uuid.UUID, bountyID uuid.UUID, applicantName, bountyTitle string) error
	CreateApplicationApprovedNotification(publisherID, applicantID uuid.UUID, bountyID uuid.UUID, bountyTitle string) error
	CreateApplicationRejectedNotification(publisherID, applicantID uuid.UUID, bountyID uuid.UUID, bountyTitle string) error
	CreateMilestoneConfirmedNotification(publisherID, receiverID uuid.UUID, bountyID uuid.UUID, milestoneTitle string) error
	CreateMilestoneCompletedNotification(receiverID, publisherID uuid.UUID, bountyID uuid.UUID, milestoneTitle string) error
	CreateSettlementAppliedNotification(receiverID, publisherID uuid.UUID, bountyID uuid.UUID, bountyTitle string) error
	CreateSettlementCompletedNotification(publisherID, receiverID uuid.UUID, bountyID uuid.UUID, bountyTitle string) error
	CreateBountyCancelledNotification(publisherID, receiverID uuid.UUID, bountyID uuid.UUID, bountyTitle string) error
	CreateUserRatedNotification(actorID, targetUserID uuid.UUID, rating float64, comment string) error
	CreateBountyLikeNotification(actorID, publisherID uuid.UUID, bountyID uuid.UUID, bountyTitle string) error
	CreateCommentNotification(actorID, publisherID uuid.UUID, bountyID uuid.UUID, commentContent string) error
}

type notificationService struct {
	notificationRepo repositories.NotificationRepository
}

// CreateUserRatedNotification 用于在“用户被评价”时自动构造通知
func (s *notificationService) CreateUserRatedNotification(actorID, targetUserID uuid.UUID, rating float64, comment string) error {
	notification := &tables.Notification{
		UserID:      targetUserID,
		ActorID:     &actorID,
		Type:        "UserRated",
		Title:       "你收到了一条评价",
		Description: "评分：" + fmt.Sprintf("%.1f", rating) + "，评价内容：" + comment,
		Metadata: map[string]any{
			"rating":  rating,
			"comment": comment,
		},
	}
	return s.notificationRepo.CreateNotification(notification)
}

// CreateBountyCancelledNotification 通知发布者 & 接收者
func (s *notificationService) CreateBountyCancelledNotification(publisherID, receiverID uuid.UUID, bountyID uuid.UUID, bountyTitle string) error {
	notifications := []*tables.Notification{
		{
			UserID:      publisherID,
			Type:        "BountyCancelled",
			Title:       "你的悬赏令已取消",
			Description: "悬赏令【" + bountyTitle + "】已被取消。",
			RelatedID:   &bountyID,
			RelatedType: "Bounty",
		},
		{
			UserID:      receiverID,
			Type:        "BountyCancelled",
			Title:       "你参与的悬赏令已取消",
			Description: "悬赏令【" + bountyTitle + "】已被取消。",
			RelatedID:   &bountyID,
			RelatedType: "Bounty",
		},
	}

	for _, notification := range notifications {
		if err := s.notificationRepo.CreateNotification(notification); err != nil {
			return err
		}
	}

	return nil
}

// CreateSettlementCompletedNotification 通知发布者 & 接收者
func (s *notificationService) CreateSettlementCompletedNotification(publisherID, receiverID uuid.UUID, bountyID uuid.UUID, bountyTitle string) error {
	notifications := []*tables.Notification{
		{
			UserID:      publisherID,
			Type:        "SettlementCompleted",
			Title:       "悬赏清算成功",
			Description: "悬赏令【" + bountyTitle + "】已清算完成。",
			RelatedID:   &bountyID,
			RelatedType: "Bounty",
		},
		{
			UserID:      receiverID,
			Type:        "SettlementCompleted",
			Title:       "悬赏清算成功",
			Description: "悬赏令【" + bountyTitle + "】已清算完成。",
			RelatedID:   &bountyID,
			RelatedType: "Bounty",
		},
	}

	for _, notification := range notifications {
		if err := s.notificationRepo.CreateNotification(notification); err != nil {
			return err
		}
	}

	return nil
}

// CreateSettlementAppliedNotification 用于在“接收者申请悬赏清算”时自动构造通知
func (s *notificationService) CreateSettlementAppliedNotification(receiverID, publisherID uuid.UUID, bountyID uuid.UUID, bountyTitle string) error {
	notification := &tables.Notification{
		UserID:      publisherID,
		ActorID:     &receiverID,
		Type:        "SettlementApplied",
		Title:       "悬赏清算申请",
		Description: "接收者申请清算悬赏令【" + bountyTitle + "】，请及时审核。",
		RelatedID:   &bountyID,
		RelatedType: "Bounty",
	}
	return s.notificationRepo.CreateNotification(notification)
}

// CreateMilestoneCompletedNotification 用于在“悬赏接收者提交了里程碑”时自动构造通知
func (s *notificationService) CreateMilestoneCompletedNotification(receiverID, publisherID uuid.UUID, bountyID uuid.UUID, milestoneTitle string) error {
	notification := &tables.Notification{
		UserID:      publisherID,
		ActorID:     &receiverID,
		Type:        "MilestoneSubmitted",
		Title:       "里程碑已提交",
		Description: "里程碑【" + milestoneTitle + "】已被接收者提交，请及时处理。",
		RelatedID:   &bountyID,
		RelatedType: "Bounty",
	}
	return s.notificationRepo.CreateNotification(notification)
}

// CreateMilestoneConfirmedNotification 用于在“悬赏发布者确认了里程碑”时自动构造通知
func (s *notificationService) CreateMilestoneConfirmedNotification(publisherID, receiverID uuid.UUID, bountyID uuid.UUID, milestoneTitle string) error {
	notification := &tables.Notification{
		UserID:      receiverID,
		ActorID:     &publisherID,
		Type:        "MilestoneConfirmed",
		Title:       "里程碑已确认",
		Description: "悬赏令的里程碑【" + milestoneTitle + "】已被发布者确认。",
		RelatedID:   &bountyID,
		RelatedType: "Bounty",
	}
	return s.notificationRepo.CreateNotification(notification)
}

// CreateApplicationRejectedNotification 用于在“我的悬赏令申请被拒绝”时自动构造通知
func (s *notificationService) CreateApplicationRejectedNotification(publisherID, applicantID uuid.UUID, bountyID uuid.UUID, bountyTitle string) error {
	notification := &tables.Notification{
		UserID:      applicantID,
		ActorID:     &publisherID,
		Type:        "ApplicationRejected",
		Title:       "你的悬赏申请被拒绝",
		Description: "很抱歉，你的悬赏申请【" + bountyTitle + "】被拒绝。",
		RelatedID:   &bountyID,
		RelatedType: "Bounty",
	}
	return s.notificationRepo.CreateNotification(notification)
}

// CreateApplicationApprovedNotification 用于在“我的悬赏令申请被批准”时自动构造通知
func (s *notificationService) CreateApplicationApprovedNotification(publisherID, applicantID uuid.UUID, bountyID uuid.UUID, bountyTitle string) error {
	notification := &tables.Notification{
		UserID:      applicantID, // 通知接收者
		ActorID:     &publisherID,
		Type:        "ApplicationApproved",
		Title:       "你的悬赏申请已被批准",
		Description: "你的悬赏申请【" + bountyTitle + "】已被批准，快去查看吧！",
		RelatedID:   &bountyID,
		RelatedType: "Bounty",
	}
	return s.notificationRepo.CreateNotification(notification)
}

// CreateBountyApplicationNotification 用于在“某人申请了我的悬赏令”时自动构造通知
func (s *notificationService) CreateBountyApplicationNotification(applicantID, publisherID uuid.UUID, bountyID uuid.UUID, applicantName, bountyTitle string) error {
	notification := &tables.Notification{
		UserID:      publisherID, // 通知接收者
		ActorID:     &applicantID,
		Type:        "BountyApplied",
		Title:       "有人申请了你的悬赏令",
		Description: "申请者【" + applicantName + "】申请了悬赏令【" + bountyTitle + "】",
		RelatedID:   &bountyID,
		RelatedType: "Bounty",
	}
	return s.notificationRepo.CreateNotification(notification)
}

// CreateBountyLikeNotification 用于在“某个用户点赞了我的悬赏令”时自动构造通知
func (s *notificationService) CreateBountyLikeNotification(actorID, publisherID uuid.UUID, bountyID uuid.UUID, bountyTitle string) error {
	notification := &tables.Notification{
		UserID:      publisherID, // 通知接收者
		ActorID:     &actorID,    // 触发者
		Type:        "BountyLiked",
		Title:       "有人点赞了你的悬赏令",
		Description: "你的悬赏令【" + bountyTitle + "】被点赞了",
		RelatedID:   &bountyID,
		RelatedType: "Bounty",
	}
	return s.notificationRepo.CreateNotification(notification)
}

// CreateCommentNotification 表示在悬赏令下评论时
func (s *notificationService) CreateCommentNotification(actorID, publisherID uuid.UUID, bountyID uuid.UUID, commentContent string) error {
	notification := &tables.Notification{
		UserID:      publisherID,
		ActorID:     &actorID,
		Type:        "BountyCommented",
		Title:       "你的悬赏令有新评论",
		Description: "评论内容：" + commentContent,
		RelatedID:   &bountyID,
		RelatedType: "Bounty",
		// 可以将评论内容也存到Metadata中
		Metadata: map[string]any{
			"comment": commentContent,
		},
	}
	return s.notificationRepo.CreateNotification(notification)
}

func NewNotificationService(notificationRepo repositories.NotificationRepository) NotificationService {
	return &notificationService{notificationRepo: notificationRepo}
}

func (s *notificationService) CreateNotification(notification *tables.Notification) error {
	return s.notificationRepo.CreateNotification(notification)
}

func (s *notificationService) GetUserNotifications(userID uuid.UUID) ([]tables.Notification, error) {
	return s.notificationRepo.FindNotificationsByUserID(userID)
}

func (s *notificationService) MarkNotificationAsRead(notificationID uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(notificationID)
}

func (s *notificationService) DeleteNotification(notificationID uuid.UUID) error {
	return s.notificationRepo.DeleteNotification(notificationID)
}
