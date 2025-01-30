package tables

import (
	"github.com/google/uuid"
)

type Rating struct {
	BaseModel
	UserID   uuid.UUID `gorm:"type:uuid;index"` // 评分用户的ID
	BountyID uuid.UUID `gorm:"type:uuid;index"` // 关联的悬赏令ID
	Score    float64   `gorm:"not null"`        // 评分分数

	// 关联
	User   User   `gorm:"foreignKey:UserID;references:ID"`
	Bounty Bounty `gorm:"foreignKey:BountyID;references:ID"`
}
