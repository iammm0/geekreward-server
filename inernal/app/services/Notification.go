package services

import (
	"GeekReward/inernal/app/models/tables"
	"GeekReward/inernal/app/repositories"
)

type NotificationService interface {
	CreateNotification(notification *tables.Notification) error
	GetUserNotifications(userID uint) ([]tables.Notification, error)
	MarkNotificationAsRead(notificationID uint) error
	DeleteNotification(notificationID uint) error
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

func (s *notificationService) GetUserNotifications(userID uint) ([]tables.Notification, error) {
	return s.notificationRepo.FindByUserID(userID)
}

func (s *notificationService) MarkNotificationAsRead(notificationID uint) error {
	return s.notificationRepo.MarkAsRead(notificationID)
}

func (s *notificationService) DeleteNotification(notificationID uint) error {
	return s.notificationRepo.Delete(notificationID)
}
