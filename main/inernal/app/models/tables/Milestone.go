package tables

import (
	"github.com/google/uuid"
	"time"
)

type Milestone struct {
	BaseModel
	// 里程碑内容，由里程碑发布者更新
	Title       string    `gorm:"not null" json:"title"`
	Description string    `json:"description"`
	DueDate     time.Time `gorm:"index" json:"due_date"`

	BountyID uuid.UUID `gorm:"type:uuid;not null;index" json:"bounty_id"`

	// 关联
	Bounty Bounty `gorm:"foreignKey:BountyID;references:ID" json:"bounty"`

	// 是否完成，由里程碑接收者更新
	IsCompleted bool `json:"is_completed"`
}
