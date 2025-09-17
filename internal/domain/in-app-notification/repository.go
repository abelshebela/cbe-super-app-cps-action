package inappnotification

import (
	"context"
)

type InAppNotificationRepository interface {
	CreateInAppNotification(ctx context.Context, notification *InAppNotification) error
	DeleteInAppNotification(ctx context.Context, id string) error
	EnableDisableInAppNotification(ctx context.Context, id string, enable bool) error
}
