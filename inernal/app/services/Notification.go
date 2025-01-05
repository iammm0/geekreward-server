package services

import (
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/repositories"
	"github.com/google/uuid"
)

type NotificationService interface {
	CreateNotification(notification *tables.Notification) error
	GetUserNotifications(userID uuid.UUID) ([]tables.Notification, error)
	MarkNotificationAsRead(notificationID uuid.UUID) error
	DeleteNotification(notificationID uuid.UUID) error
}

type notificationService struct {
	notificationRepo repositories.NotificationRepository
}

func NewNotificationService(notificationRepo repositories.NotificationRepository) NotificationService {
	return &notificationService{notificationRepo: notificationRepo}
}

func (s *notificationService) CreateNotification(notification *tables.Notification) error {
	return s.notificationRepo.Create(notification)
}

func (s *notificationService) GetUserNotifications(userID uuid.UUID) ([]tables.Notification, error) {
	return s.notificationRepo.FindByUserID(userID)
}

func (s *notificationService) MarkNotificationAsRead(notificationID uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(notificationID)
}

func (s *notificationService) DeleteNotification(notificationID uuid.UUID) error {
	return s.notificationRepo.Delete(notificationID)
}
