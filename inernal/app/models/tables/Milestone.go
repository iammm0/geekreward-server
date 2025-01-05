package tables

import (
	"github.com/google/uuid"
	"time"
)

type Milestone struct {
	BaseModel
	Title       string `gorm:"not null"`
	Description string
	DueDate     time.Time `gorm:"index"`
	BountyID    uuid.UUID `gorm:"type:uuid;not null;index"`

	// 关联
	Bounty Bounty `gorm:"foreignKey:BountyID;references:ID"`
}
