package persistance

import (
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/external_call/merchant_lookup"
	"cbe-super-app-cps-action/internal/storage/kafka"

	"github.com/hugokessem/coreio/core"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Persistence struct {
	DeviceVersionControlPersistence storage.DeviceVersionControlRepository
	AccountLookup                   core.CBECoreAPIInterface
	MerchantLookup                  merchant_lookup.MerchantLookupAdapter
	MongoClient                     *mongo.Client
	UserPersistence                 storage.UserRepository
	UnlinkAccountPersistence        storage.UnlinkAccount
	HQPersistence                   storage.HQRepository
	OTPPersistence                  storage.OTPRepository
	DeviceLinkHistoryPersistence    storage.DeviceLinkHistoryRepository
	ResetSessionPersistence         storage.ResetSessionRepository
	SMSSenderApi                    kafka.NotificationProducer
	CPSAction                       storage.CPSActionRepository
	AdvertPersistence               storage.AdvertRepository
	AmountBasedAuthPersistence      storage.AmountBasedAuthRepository
	AccountBlockPersistence         storage.AccountBlockRepository
	PortalCardPersistence           storage.PortalCardRepository
	MiniAppPersistence              storage.MiniAppRepository
	CityPersistence                 storage.CityRepository
	RegionPersistence               storage.RegionRepository
	DistrictPersistence             storage.DistrictRepository
	BranchPersistence               storage.BranchRepository
	EventPersistence                storage.EventRepository
	CustomerService                 storage.CustomerRepository
	BulkService                     storage.BulkServiceRepository
	BudgetCategoryPersistence       storage.BudgetCategoryRepository
	RedisService                    storage.RedisRepository
	IconPersistence                 storage.IconRepository
	// Additional repositories
	AccessListPersistence            storage.AppAccessListRepository
	AvatarPersistence                storage.AvatarRepository
	BPSUserPersistence               storage.BPSUserRepository
	AdvertRepositoryPersistence      storage.AdvertRepository
	ArchivedUserPersistence          storage.ArchivedUserRepository
	ArchivedLinkedAccountPersistence storage.ArchivedLinkedAccountRepository
	AuthTierPersistence              storage.AuthTierRepository
	BankPersistence                  storage.BankRepository
	CpsUserPersistence               storage.CpsUserRepository
	DonationPersistence              storage.DonationRepository
	DonationCategoryPersistence      storage.DonationCategoryRepository
	DonationCompanyPersistence       storage.DonationCompanyRepository
	FeedbackPersistence              storage.FeedbackRepository
	LinkedAccountPersistence         storage.LinkedAccountRepository
	MiniAppMerchantPersistence       storage.MiniAppMerchantRepository
	NotificationPersistence          storage.NotificationRepository
	PasswordRulePersistence          storage.PasswordRuleRepository
	ServiceDetailsPersistence        storage.ServiceDetailsRepository
	ServicesPersistence              storage.ServicesRepository
	ValidationRulePersistence        storage.ValidationRuleRepository
	PasswordRulePersistent           storage.PasswordRuleRepository
	WalletPersistence                storage.WalletRepository
	TopupPersistence                 storage.TopupRepository
	ProductCodePersistence           storage.ProductCodeRepository
	DepartmentPersistence            storage.DepartmentRepository
	FaydaPersistence                 storage.FaydaRepository
	PermissionPersistence            storage.PermissionRepository
	ArticlePersistence               storage.ArticleRepository
	ArticleCategoryPersistence       storage.ArticleCategoryRepository
	ShortVideoPersistence            storage.ShortVideoRepository
	NewsTagPersistence               storage.NewsTagRepository
	KYCVerifierPersistence           storage.KYCVerifierRepository
	NewsCategoryPersistence          storage.NewsCategoryRepository
	NewsTagsServiceContainer         storage.NewsTagsRepository
	BPSActionRolePersistence         storage.BPSActionRoleRepository
	BPSActionApproveIndexPersistence storage.BPSActionApproveIndexRepository
	RolePersistence                  storage.RoleRepository
	MiniAppCategoryPersistence       storage.MiniAppCategoryRepository
	CPSActionRolePersistence         storage.CPSActionRoleRepository
	CPSActionApproveIndexPersistence storage.CPSActionApproveIndexRepository
	MiniAppProductCodePersistence    storage.MiniAppProductCodeRepository
}
