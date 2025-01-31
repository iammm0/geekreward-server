package tables

import (
	"github.com/google/uuid"
)

// Like 表示用户点赞模型
type Like struct {
	BaseModel
	UserID   uuid.UUID `gorm:"type:uuid;not null;index"` // 点赞用户的ID
	BountyID uuid.UUID `gorm:"type:uuid;not null;index"` // 关联的悬赏令ID

	// 关联
	User   User   `gorm:"foreignKey:UserID;references:ID"`
	Bounty Bounty `gorm:"foreignKey:BountyID;references:ID"`
}
