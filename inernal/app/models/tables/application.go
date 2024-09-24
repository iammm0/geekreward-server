package tables

import (
	"gorm.io/gorm"
)

// Application 模型，用于存储用户对悬赏令的申请
type Application struct {
	gorm.Model
	BountyID uint   `gorm:"not null"` // 关联的悬赏令ID
	UserID   uint   `gorm:"not null"` // 申请用户的ID
	Status   string `gorm:"not null"` // 申请状态
	Note     string // 申请备注
}
