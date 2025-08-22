package notification

import (
	"context"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/notification"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification *entities.Notification) (*entities.Notification, error)
	UpdateNotification(ctx context.Context, notification *entities.Notification) (*entities.Notification, error)
	DeleteNotification(ctx context.Context, id string) (*entities.Notification, error)
	EnableDisableNotification(ctx context.Context, id string, enable bool) (*entities.Notification, error)
	NotificationExists(ctx context.Context, notificationType string, forValue entities.NotificationFor, id *string) (bool, error)
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	FetchNotificationByID(ctx context.Context, id string) (*entities.Notification, error)
	FetchNotifications(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*entities.Notification], error)
}
