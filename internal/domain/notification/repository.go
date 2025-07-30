package notification

import (
	"context"

	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
)

type NotificationRepository interface {
	CreateNotification(ctx context.Context, notification *Notification) (*Notification, error)
	UpdateNotification(ctx context.Context, notification *Notification) (*Notification, error)
	DeleteNotification(ctx context.Context, id string) (*Notification, error)
	EnableDisableNotification(ctx context.Context, id string, enable bool) (*Notification, error)
	NotificationExists(ctx context.Context, notificationType string, forValue NotificationFor, id *string) (bool, error)

	FetchNotificationByID(ctx context.Context, id string) (*Notification, error)
	FetchNotifications(ctx context.Context, filterParam *constant.Filter) (*common_util.PaginatedResponse[[]*Notification], error)
}
