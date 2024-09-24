package tables

import (
	"gorm.io/gorm"
)

// Notification 模型，用于存储用户通知信息
type Notification struct {
	ID          uint           `gorm:"primaryKey"`
	UserID      uint           // 与用户表的外键关系
	Title       string         `gorm:"size:255;not null"`
	Description string         `gorm:"type:text"`
	IsRead      bool           `gorm:"default:false"`
	DeletedAt   gorm.DeletedAt `gorm:"index"` // 软删除字段
}
