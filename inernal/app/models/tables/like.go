package tables

import "gorm.io/gorm"

// Like 表示用户点赞模型
type Like struct {
	gorm.Model
	UserID   uint `gorm:"not null"` // 点赞用户的ID
	BountyID uint `gorm:"not null"` // 关联的悬赏令ID
}
