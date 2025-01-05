package tables

import (
	"github.com/google/uuid"
	"time"
)

// WorkExperience 表示用户的工作经验模型
type WorkExperience struct {
	BaseModel
	CompanyName      string    // 公司名
	JobTitle         string    // 职位类型
	StartDate        time.Time // 开始时间
	EndDate          time.Time // 结束时间
	Responsibilities string    // 职责职责
	Achievements     string    // 成就
	UserID           uuid.UUID `gorm:"type:uuid;index"` // 用户ID外键

	// 关联
	User User `gorm:"foreignKey:UserID;references:ID"`
}
