package initiator

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage/external_call"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage/persistance/device_history"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage/persistance/hq"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage/persistance/otp"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage/persistance/reset_session"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/storage/persistance/users"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Persistence struct {
	UserPersistence              storage.UserRepository
	HQPersistence                storage.HQRepository
	OTPPersistence               storage.OTPRepository
	DeviceLinkHistoryPersistence storage.DeviceLinkHistoryRepository
	ResetSessionPersistence      storage.ResetSessionRepository
	SMSSenderApi                 external_call.SMSPersistence
}

func InitPersistanceLayer(client *mongo.Client, dbName string, logger utils.Logger) Persistence {
	return Persistence{
		UserPersistence:              users.NewUserRepository(client, dbName, "users", logger),
		HQPersistence:                hq.NewHQRepository(client, dbName, "hq", logger),
		OTPPersistence:               otp.NewOtpRepository(client, dbName, "otps", logger),
		DeviceLinkHistoryPersistence: device_history.NewDeviceLinkHistoryRepository(client, dbName, "member_device_histories", logger),
		ResetSessionPersistence:      reset_session.NewResetSessionRepository(client, dbName, "pin_reset_sessions", logger),
		SMSSenderApi:                 *external_call.NewSMSPersistence("https://devcbe.eaglelionsystems.com/api/v1.0/chatbirrapi/ldapnotif/sms/send", logger),
	}
}
