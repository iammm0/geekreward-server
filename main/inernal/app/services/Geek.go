package services

import (
	tables2 "GeekReward/main/inernal/app/models/tables"
	repositories2 "GeekReward/main/inernal/app/repositories"
	"errors"
	"github.com/google/uuid"
	"time"
)

type GeekService interface {
	GetTopGeeks(limit int) ([]tables2.User, error)
	GetGeekByID(id uuid.UUID) (*tables2.User, error)
	SendInvitation(geekID uuid.UUID, inviterID uuid.UUID) error
	ExpressAffection(geekID uuid.UUID, userID uuid.UUID) error
}

type geekService struct {
	geekRepo       repositories2.GeekRepository
	invitationRepo repositories2.InvitationRepository
}

func (s *geekService) ExpressAffection(geekID uuid.UUID, userID uuid.UUID) error {
	// 检查极客是否存在
	geek, err := s.geekRepo.GetGeekByID(geekID)
	if err != nil {
		return err
	}
	if geek == nil {
		return errors.New("geek not found")
	}

	// 检查用户是否已经表达过好感
	existingLike, err := s.geekRepo.GetAffection(userID, geekID)
	if err != nil {
		return err
	}
	if existingLike != nil {
		return errors.New("you have already expressed affection to this geek")
	}

	// 创建好感记录
	affection := &tables2.Affection{
		UserID: userID,
		GeekID: geekID,
	}

	if err := s.geekRepo.CreateAffection(affection); err != nil {
		return err
	}

	return nil
}

func NewGeekService(
	geekRepo repositories2.GeekRepository,
	invitationRepo repositories2.InvitationRepository,
) GeekService {
	return &geekService{
		geekRepo:       geekRepo,
		invitationRepo: invitationRepo,
	}
}

func (s *geekService) GetTopGeeks(limit int) ([]tables2.User, error) {
	return s.geekRepo.GetTopGeeks(limit)
}

func (s *geekService) GetGeekByID(id uuid.UUID) (*tables2.User, error) {
	return s.geekRepo.GetGeekByID(id)
}

// SendInvitation 向特定极客发出组队邀请
func (s *geekService) SendInvitation(geekID uuid.UUID, inviterID uuid.UUID) error {
	// 检查极客是否存在
	geek, err := s.geekRepo.GetGeekByID(geekID)
	if err != nil {
		return err
	}
	if geek == nil {
		return errors.New("geek not found")
	}

	// 检查是否已经存在邀请
	existingInvitation, err := s.invitationRepo.GetInvitation(inviterID, geekID)
	if err != nil {
		return err
	}
	if existingInvitation != nil {
		return errors.New("invitation already exists")
	}

	// 创建邀请
	invitation := &tables2.Invitation{
		ID:        uuid.New(),
		InviterID: inviterID,
		InviteeID: geekID,
		Status:    "Pending", // Pending, Accepted, Rejected
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.invitationRepo.CreateInvitation(invitation); err != nil {
		return err
	}

	return nil
}
