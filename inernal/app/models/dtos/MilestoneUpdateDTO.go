package dtos

import "time"

type MilestoneUpdateDTO struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	DueDate     time.Time `json:"due_date" binding:"required"`
	IsCompleted bool      `json:"is_completed"`
}
