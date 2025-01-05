package dtos

import "github.com/google/uuid"

type CreateUserNotificationInput struct {
	UserID      uuid.UUID `json:"user_id" binding:"required,uuid"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
}
