package tables

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Notification 模型，用于存储用户通知信息
type Notification struct {
	BaseModel
	UserID      uuid.UUID  `gorm:"type:uuid;index"` // 与用户表的外键关系
	ActorID     *uuid.UUID `gorm:"type:uuid;index"` // 触发者（可选）
	Type        string     `gorm:"size:100;not null"`
	Title       string     `gorm:"size:255;not null"`
	Description string     `gorm:"type:text"`

	// 扩展字段
	RelatedID   *uuid.UUID     `gorm:"type:uuid;index"` // 关联资源(可选)
	RelatedType string         `gorm:"size:100"`        // 关联资源类型(Bounty, Invitation, etc.)
	Metadata    map[string]any `gorm:"type:jsonb"`      // JSON字段(需Postgres或自行序列化)

	IsRead    bool           `gorm:"default:false"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // 软删除字段

	// 关联
	User  User  `gorm:"foreignKey:UserID;references:ID"`
	Actor *User `gorm:"foreignKey:ActorID;references:ID"` // 触发者(可选)
}
