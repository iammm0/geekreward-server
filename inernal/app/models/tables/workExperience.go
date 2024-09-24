package tables

import (
	"gorm.io/gorm"
	"time"
)

// WorkExperience 表示用户的工作经验模型
type WorkExperience struct {
	gorm.Model
	CompanyName      string    // 公司名
	JobTitle         string    // 职位类型
	StartDate        time.Time // 开始时间
	EndDate          time.Time // 结束时间
	Responsibilities string    // 责任职责
	Achievements     string    // 成就
	UserID           uint      // 用户ID外键
}
