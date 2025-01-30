package tables

import (
	"github.com/google/uuid"
	"time"
)

type Invitation struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;" json:"id"`
	InviterID uuid.UUID `gorm:"type:uuid;" json:"inviter_id"`
	InviteeID uuid.UUID `gorm:"type:uuid;" json:"invitee_id"`
	Status    string    `json:"status"` // Pending, Accepted, Rejected
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
