package tables

import (
	"github.com/google/uuid"
)

// Comment 表示评论模型
type Comment struct {
	BaseModel
	Content  string    `gorm:"not null"`
	UserID   uuid.UUID `gorm:"type:uuid;not null;index"` // 关联用户
	BountyID uuid.UUID `gorm:"type:uuid;not null;index"` // 关联悬赏令

	// 关联
	User   User   `gorm:"foreignKey:UserID;references:ID"`
	Bounty Bounty `gorm:"foreignKey:BountyID;references:ID"`
}
