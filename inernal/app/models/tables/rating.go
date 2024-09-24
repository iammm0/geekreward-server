package tables

import (
	"gorm.io/gorm"
)

type Rating struct {
	gorm.Model
	UserID   uint    // 评分用户的ID
	BountyID uint    // 关联的悬赏令ID
	Score    float64 // 评分分数
}
