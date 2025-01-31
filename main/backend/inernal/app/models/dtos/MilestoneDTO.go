package dtos

import (
	"time"
)

type MilestoneDTO struct {
	Title       string    `json:"title" binding:"required"`       // 里程碑标题
	Description string    `json:"description" binding:"required"` // 里程碑描述
	DueDate     time.Time `json:"due_date" binding:"required"`    // 截止日期 (格式: YYYY-MM-DD)
}
