package migrations

import (
	"GeekReward/inernal/app/models/tables"
	"gorm.io/gorm"
)

// Migrate 自动迁移数据库模型
func Migrate(db *gorm.DB) error {

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	return db.AutoMigrate(
		&tables.User{},
		&tables.Bounty{},
		&tables.Milestone{},
		&tables.Comment{},
		&tables.Application{},
		&tables.Notification{},
		&tables.Like{},
		&tables.Rating{},
		&tables.Invitation{},
		&tables.Affection{},
	)
}
