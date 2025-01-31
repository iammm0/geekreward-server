package tables

import (
	"github.com/google/uuid"
)

type Affection struct {
	BaseModel
	UserID uuid.UUID `gorm:"type:uuid;" json:"user_id"`
	GeekID uuid.UUID `gorm:"type:uuid;" json:"geek_id"`
}
