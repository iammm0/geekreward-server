package repositories

import (
	"GeekReward/main/inernal/app/models/tables"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InvitationRepository interface {
	GetInvitation(inviterID, inviteeID uuid.UUID) (*tables.Invitation, error)
	CreateInvitation(invitation *tables.Invitation) error
	GetInvitationByID(id uuid.UUID) (*tables.Invitation, error)
	UpdateInvitation(invitation *tables.Invitation) error
}

type invitationRepository struct {
	db *gorm.DB
}

func NewInvitationRepository(db *gorm.DB) InvitationRepository {
	return &invitationRepository{db: db}
}

func (r *invitationRepository) GetInvitation(inviterID, inviteeID uuid.UUID) (*tables.Invitation, error) {
	var invitation tables.Invitation
	if err := r.db.First(&invitation, "inviter_id = ? AND invitee_id = ?", inviterID, inviteeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &invitation, nil
}

func (r *invitationRepository) CreateInvitation(invitation *tables.Invitation) error {
	return r.db.Create(invitation).Error
}

func (r *invitationRepository) GetInvitationByID(id uuid.UUID) (*tables.Invitation, error) {
	var invitation tables.Invitation
	if err := r.db.First(&invitation, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &invitation, nil
}

func (r *invitationRepository) UpdateInvitation(invitation *tables.Invitation) error {
	return r.db.Save(invitation).Error
}
