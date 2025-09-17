package inappnotification

import (
	"context"

	entities "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/in-app-notification"
)

type InAppNotificationRepository interface {
	CreateInAppNotification(ctx context.Context, notification *entities.InAppNotification) error
	DeleteInAppNotification(ctx context.Context, id string) error
	EnableDisableInAppNotification(ctx context.Context, id string, enable bool) error
}
