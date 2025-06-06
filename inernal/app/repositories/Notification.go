package repositories

import (
	"GeekReward/inernal/app/models/tables"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	CreateNotification(notification *tables.Notification) error
	FindNotificationsByUserID(userID uuid.UUID) ([]tables.Notification, error)
	MarkAsRead(notificationID uuid.UUID) error
	DeleteNotification(notificationID uuid.UUID) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository(db *gorm.DB) NotificationRepository {
	return &notificationRepository{db: db}
}

func (r *notificationRepository) CreateNotification(notification *tables.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) FindNotificationsByUserID(userID uuid.UUID) ([]tables.Notification, error) {
	var notifications []tables.Notification
	err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) MarkAsRead(notificationID uuid.UUID) error {
	return r.db.Model(&tables.Notification{}).Where("id = ?", notificationID).Update("is_read", true).Error
}

func (r *notificationRepository) DeleteNotification(notificationID uuid.UUID) error {
	return r.db.Delete(&tables.Notification{}, "id = ?", notificationID).Error
}
