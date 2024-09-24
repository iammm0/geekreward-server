package tables

import (
	"gorm.io/gorm"
)

// Comment 表示评论模型
type Comment struct {
	gorm.Model
	Content  string `gorm:"not null"`
	UserID   uint   `gorm:"not null"` // 关联用户
	BountyID uint   `gorm:"not null"` // 关联悬赏令
}
