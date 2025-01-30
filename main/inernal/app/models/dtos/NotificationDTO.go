package dtos

import (
	"github.com/google/uuid"
	"time"
)

type NotificationDTO struct {
	ID          uuid.UUID `json:"id" binding:"uuid"`
	UserID      uuid.UUID `json:"user_id" binding:"required,uuid"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	IsRead      bool      `json:"is_read"`
	CreatedAt   time.Time `json:"created_at"`
}
