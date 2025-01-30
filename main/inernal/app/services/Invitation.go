package services

import (
	repositories2 "GeekReward/main/inernal/app/repositories"
	"errors"
	"github.com/google/uuid"
)

type InvitationService interface {
	AcceptInvitation(invitationID uuid.UUID, userID uuid.UUID) error
	RejectInvitation(invitationID uuid.UUID, userID uuid.UUID) error
}

type invitationService struct {
	invitationRepo repositories2.InvitationRepository
	userRepo       repositories2.UserRepository
}

func NewInvitationService(
	invitationRepo repositories2.InvitationRepository,
	userRepo repositories2.UserRepository,
) InvitationService {
	return &invitationService{
		invitationRepo: invitationRepo,
		userRepo:       userRepo,
	}
}

// AcceptInvitation 接受组队邀请
func (s *invitationService) AcceptInvitation(invitationID uuid.UUID, userID uuid.UUID) error {
	// 获取邀请
	invitation, err := s.invitationRepo.GetInvitationByID(invitationID)
	if err != nil {
		return err
	}
	if invitation == nil {
		return errors.New("invitation not found")
	}

	// 验证邀请接受者是否为邀请的受邀者
	if invitation.InviteeID != userID {
		return errors.New("you are not authorized to accept this invitation")
	}

	// 更新邀请状态为已接受
	invitation.Status = "Accepted"
	if err := s.invitationRepo.UpdateInvitation(invitation); err != nil {
		return err
	}

	// 执行其他逻辑，如将邀请者和受邀者加入同一团队
	// team, err := s.teamRepo.GetTeamByUserID(invitation.InviterID)
	// if err != nil {
	//     return err
	// }
	// if team == nil {
	//     team = &tables.Team{
	//         ID:      uuid.New(),
	//         LeaderID: invitation.InviterID,
	//         Members: []uuid.UUID{invitation.InviterID, invitation.InviteeID},
	//     }
	//     if err := s.teamRepo.CreateTeam(team); err != nil {
	//         return err
	//     }
	// } else {
	//     team.Members = append(team.Members, invitation.InviteeID)
	//     if err := s.teamRepo.UpdateTeam(team); err != nil {
	//         return err
	//     }
	// }

	return nil
}

// RejectInvitation 拒绝组队邀请
func (s *invitationService) RejectInvitation(invitationID uuid.UUID, userID uuid.UUID) error {
	// 获取邀请
	invitation, err := s.invitationRepo.GetInvitationByID(invitationID)
	if err != nil {
		return err
	}
	if invitation == nil {
		return errors.New("invitation not found")
	}

	// 验证邀请接受者是否为邀请的受邀者
	if invitation.InviteeID != userID {
		return errors.New("you are not authorized to reject this invitation")
	}

	// 更新邀请状态为已拒绝
	invitation.Status = "Rejected"
	if err := s.invitationRepo.UpdateInvitation(invitation); err != nil {
		return err
	}

	// 发送通知给邀请者（假设有通知服务）
	// notification := &tables.Notification{
	//     UserID:    invitation.InviterID,
	//     Type:      "invitation_rejected",
	//     Message:   "Your invitation has been rejected.",
	//     IsRead:    false,
	//     CreatedAt: time.Now(),
	//     UpdatedAt: time.Now(),
	// }
	// if err := s.notificationRepo.CreateNotification(notification); err != nil {
	//     return err
	// }

	return nil
}
