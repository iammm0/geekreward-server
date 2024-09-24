package migrations

import (
	"GeekReward/inernal/app/models/tables"
	"gorm.io/gorm"
)

// Migrate 自动迁移数据库模型
func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&tables.User{},
		&tables.Bounty{},
		&tables.Milestone{},
		&tables.Comment{}, // 确保 comment 迁移在外键关联前
		&tables.Application{},
		&tables.Notification{},
		&tables.Like{},
		&tables.Rating{},
	)
}
