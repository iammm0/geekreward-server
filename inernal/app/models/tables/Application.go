package tables

import (
	"github.com/google/uuid"
)

// Application 模型，用于存储用户对悬赏令的申请
type Application struct {
	BaseModel
	BountyID uuid.UUID `gorm:"type:uuid;not null;index"` // 关联的悬赏令ID
	UserID   uuid.UUID `gorm:"type:uuid;not null;index"` // 申请用户的ID
	Status   string    `gorm:"not null"`                 // 申请状态
	Note     string    // 申请备注

	// 关联
	Bounty Bounty `gorm:"foreignKey:BountyID;references:ID"`
	User   User   `gorm:"foreignKey:UserID;references:ID"`
}
