package persistance

import (
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/external_call"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Persistence struct {
	MongoClient                  *mongo.Client
	UserPersistence              storage.UserRepository
	UnlinkAccountPersistence     storage.UnlinkAccount
	HQPersistence                storage.HQRepository
	OTPPersistence               storage.OTPRepository
	DeviceLinkHistoryPersistence storage.DeviceLinkHistoryRepository
	ResetSessionPersistence      storage.ResetSessionRepository
	SMSSenderApi                 external_call.SMSPersistence
	CPSAction                    storage.CPSActionRepository
	AdvertPersistence            storage.AccountAPIPort
	AmountBasedAuthPersistence   storage.AmountBasedAuthRepository
	AccountBlockPersistence      storage.AccountBlockRepository
	PortalCardPersistence        storage.PortalCardRepository
	MiniAppPersistence           storage.MiniAppRepository
	CityPersistence              storage.CityRepository
	RegionPersistence            storage.RegionRepository
	DistrictPersistence          storage.DistrictRepository
	BranchPersistence            storage.BranchRepository
	EventPersistence             storage.EventRepository

	// Additional repositories
	AccessListPersistence            storage.AppAccessListRepository
	AvatarPersistence                storage.AvatarRepository
	BPSUserPersistence               storage.BPSUserRepository
	AdvertRepositoryPersistence      storage.AdvertRepository
	ArchivedUserPersistence          storage.ArchivedUserRepository
	ArchivedLinkedAccountPersistence storage.ArchivedLinkedAccountRepository
	AuthTierPersistence              storage.AuthTierRepository
	BankPersistence                  storage.BankRepository
	ColorPersistence                 storage.ColorRepository
	CpsUserPersistence               storage.CpsUserRepository
	DonationPersistence              storage.DonationRepository
	DonationCategoryPersistence      storage.DonationCategoryRepository
	DonationCompanyPersistence       storage.DonationCompanyRepository
	FeedbackPersistence              storage.FeedbackRepository
	IconPersistence                  storage.IconRepository
	LinkedAccountPersistence         storage.LinkedAccountRepository
	MiniAppMerchantPersistence       storage.MiniAppMerchantRepository
	NotificationPersistence          storage.NotificationRepository
	PasswordRulePersistence          storage.PasswordRuleRepository
	ServiceDetailsPersistence        storage.ServiceDetailsRepository
	ValidationRulePersistence        storage.ValidationRuleRepository
	WalletPersistence                storage.WalletRepository
	AccountValidation                storage.AccountValidationRepository
}
