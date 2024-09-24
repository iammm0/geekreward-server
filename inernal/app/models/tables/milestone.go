package tables

import (
	"gorm.io/gorm"
	"time"
)

type Milestone struct {
	gorm.Model
	Title       string `gorm:"not null"`
	Description string
	DueDate     time.Time `gorm:"index"`
	BountyID    uint      `gorm:"not null;index"`
}
