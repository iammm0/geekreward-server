package migrations

import (
	tables2 "GeekReward/main/inernal/app/models/tables"
	"gorm.io/gorm"
)

// Migrate 自动迁移数据库模型
func Migrate(db *gorm.DB) error {

	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")

	return db.AutoMigrate(
		// 基础用户数据库表
		&tables2.User{},

		// 基础悬赏令表  为用户对象所拥有或申请
		&tables2.Bounty{},

		// 关联于悬赏令  为悬赏令阶段性的分节
		&tables2.Milestone{},

		// 连接用户与悬赏令交互的申请与通知模型
		&tables2.Application{},
		&tables2.Notification{},

		// 用户与悬赏令的交互使用的模型
		&tables2.Comment{},
		&tables2.Like{},
		&tables2.Rating{},

		// 极客与极客之间的社交活动模型
		&tables2.Affection{},
		&tables2.Invitation{},
	)
}
