package storage

import (
	"context"
	"time"

	session "cbe-super-app-cps-action/grpc"
	"cbe-super-app-cps-action/internal/constants/dto/donation"
	"cbe-super-app-cps-action/internal/constants/dto/donation_category"
	"cbe-super-app-cps-action/internal/constants/dto/feedback"

	cps_user_dto "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	"cbe-super-app-cps-action/internal/constants/dto/donation_company"

	"cbe-super-app-cps-action/internal/constants"

	kyc_dto "cbe-super-app-cps-action/internal/constants/dto/kyc_verifier"
	miniappdto "cbe-super-app-cps-action/internal/constants/dto/mini_app"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage/external_call"

	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/grpc"
)

type UnlinkAccount interface {
	GetUserByAccount(ctx context.Context, accNumber string) (*model.ArchivedUser, error)
	GetAllArchivedUser(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[*model.ArchivedUser], error)
	Delete(ctx context.Context, id string) error
	Authorize(ctx context.Context, cpsAction model.CPSAction) (model.ArchivedUser, error)
}

type OTPRepository interface {
	Save(ctx context.Context, otp *model.OTP) error
	Find(ctx context.Context, filte bson.M) (*model.OTP, error)
	Delete(ctx context.Context, id string) error
}

type UserRepository interface {
	Save(ctx context.Context, user *model.User) error
	GetUserByAccount(ctx context.Context, accNumber string) (*model.User, error)
	Delete(ctx context.Context, id string) error
	DeleteHard(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (*model.User, error)
	FindByUserCode(ctx context.Context, userCode string) (*model.User, error)
	FindByCustomerNumber(ctx context.Context, customerNumber string) (*model.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*model.User, error)
	FindByDeviceUUID(ctx context.Context, deviceUUID string) (*model.User, error)
	Update(ctx context.Context, id string, update *model.User) error
}

type HQRepository interface {
	FindByID(ctx context.Context, id string) (*model.HQ, error)
	Find(ctx context.Context, projections ...bson.M) (*model.HQ, error)
	Update(ctx context.Context, field string, value interface{}, now time.Time) error
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.HQ], error)
}

type DeviceLinkHistoryRepository interface {
	Save(ctx context.Context, deviceLinkHistory *model.DeviceLinkHistroy) error
	Update(ctx context.Context, update *model.DeviceLinkHistroy) error
}

type AppAccessListRepository interface {
	Create(ctx context.Context, accessList *model.APPAccessList) error
	Update(ctx context.Context, id string, accessList *model.APPAccessList) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.APPAccessList, error)
	FindAllWithPagination(ctx context.Context, department string, filterParam types.Filter) (*types.PaginatedResponse[[]*model.APPAccessList], error)
}

type ResetSessionRepository interface {
	Save(ctx context.Context, session *model.PinResetSession) error
	FindById(ctx context.Context, id string) (*model.PinResetSession, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*model.PinResetSession, error)
	FindByDeviceUUID(ctx context.Context, deviceUUID string) (*model.PinResetSession, error)
	Update(ctx context.Context, id string, update *model.PinResetSession) error
	Delete(ctx context.Context, id string) error
}

type GenericRepository[T any] interface {
	Create(ctx context.Context, entity *T) error
	FindOne(ctx context.Context, filter bson.M, projection bson.M) (*T, error)
	FindAll(ctx context.Context, filter bson.M, projection bson.M) ([]*T, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*T], error)
	Update(ctx context.Context, filter bson.M, update bson.M) error
	Delete(ctx context.Context, filter bson.M) error
}

type SessionGRPCPort interface {
	CreateSession(ctx context.Context, request *session.CreateSessionRequest, opts ...grpc.CallOption) (*session.CreateSessionResponse, error)
	UpdateSession(ctx context.Context, request *session.UpdateSessionRequest, opts ...grpc.CallOption) (*session.UpdateSessionResponse, error)
	GetSession(ctx context.Context, request *session.GetSessionRequest, opts ...grpc.CallOption) (*session.GetSessionResponse, error)
	HealthCheck(ctx context.Context, opts ...grpc.CallOption) (*session.HealthCheckResponse, error)
	Close() error
}

type DeviceVersionControlRepository interface {
	Save(ctx context.Context, deviceVersionControl model.DeviceVersionControl) error
	FindOne(ctx context.Context, filter bson.M) (model.DeviceVersionControl, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (types.PaginatedResponse[[]model.DeviceVersionControl], error)
	Update(ctx context.Context, id string, update bson.M) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
}

type CPSActionRepository interface {
	Save(ctx context.Context, cpsAction *model.CPSAction) error
	FindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (types.PaginatedResponse[[]model.CPSAction], error)
	FindOne(ctx context.Context, filter bson.M) (*model.CPSAction, error)
	SanitizedFindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (*types.PaginatedResponse[[]*model.CPSAction], error)
	SanitizedFindOne(ctx context.Context, filter bson.M) (*model.CPSAction, error)
	Update(ctx context.Context, actionCode string, update model.CPSAction) (*model.CPSAction, error)
	UpdateCustome(ctx context.Context, filter, update bson.M) error
	Delete(ctx context.Context, id string) error
}

// Avatar persistence
type AvatarRepository interface {
	Create(ctx context.Context, avatar *model.Avatar) error
	Update(ctx context.Context, id string, avatar *model.Avatar) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Avatar, error)
	Find(ctx context.Context, filter bson.M, projection bson.M) (*model.Avatar, error)
	FindAll(ctx context.Context, filter bson.M, projection bson.M) ([]*model.Avatar, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Avatar], error)
}

// BPSUser persistence
type BPSUserRepository interface {
	GetByUserCode(ctx context.Context, userCode string) (*model.BPSUser, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.BPSUser], error)
	Update(ctx context.Context, BpsUser *model.BPSUser) error
}

type BudgetCategoryRepository interface {
	CreateBudgetCategory(ctx context.Context, budgetCategory *model.BudgetCategory) error
	UpdateBudgetCategory(ctx context.Context, id string, budgetCategory *model.BudgetCategory) error
	FindBudgetCategoryByID(ctx context.Context, id string) (*model.BudgetCategory, error)
	FindAllBudgetCategories(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.BudgetCategory], error)
	DeleteBudgetCategory(ctx context.Context, id string) error
	EnableOrDisableBudgetCategory(ctx context.Context, id string, enable bool) error
}

// AmountBasedAuth persistence
type AmountBasedAuthRepository interface {
	Update(ctx context.Context, id string, update *model.AuthTier) error
	FindAll(ctx context.Context, filter bson.M, projection bson.M) ([]*model.AuthTier, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AuthTier], error)
	FindByID(ctx context.Context, id string) (*model.AuthTier, error)
}

// AccountBlock persistence
type AccountBlockRepository interface {
	GetBranchByCode(ctx context.Context, branchCode string) (*model.AccountBlock, error)
	CreateBranch(ctx context.Context, branch *model.AccountBlock) error
	UpdateBranch(ctx context.Context, id string, branch *model.AccountBlock) error
	DeleteBranch(ctx context.Context, id string) error
	EnableOrDisableBranch(ctx context.Context, code string, reason string, enabled bool) error
	FindBranchByID(ctx context.Context, id string) (*model.AccountBlock, error)
	FindAllBranchesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error)

	GetCityByCode(ctx context.Context, cityCode string) (*model.AccountBlock, error)
	CreateCity(ctx context.Context, city *model.AccountBlock) error
	UpdateCity(ctx context.Context, id string, city *model.AccountBlock) error
	DeleteCity(ctx context.Context, id string) error
	EnableOrDisableCity(ctx context.Context, code string, reason string, enabled bool) error
	FindCityByID(ctx context.Context, id string) (*model.AccountBlock, error)
	FindAllCitiesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error)

	GetRegionByCode(ctx context.Context, regionCode string) (*model.AccountBlock, error)
	CreateRegion(ctx context.Context, region *model.AccountBlock) error
	UpdateRegion(ctx context.Context, id string, region *model.AccountBlock) error
	DeleteRegion(ctx context.Context, id string) error
	EnableOrDisableRegion(ctx context.Context, code string, reason string, enabled bool) error
	FindRegionByID(ctx context.Context, id string) (*model.AccountBlock, error)
	FindAllRegionsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error)

	GetDistrictByCode(ctx context.Context, districtCode string) (*model.AccountBlock, error)
	CreateDistrict(ctx context.Context, district *model.AccountBlock) error
	UpdateDistrict(ctx context.Context, id string, district *model.AccountBlock) error
	DeleteDistrict(ctx context.Context, id string) error
	EnableOrDisableDistrict(ctx context.Context, code string, reason string, enabled bool) error
	FindDistrictByID(ctx context.Context, id string) (*model.AccountBlock, error)
	FindAllDistrictsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error)
}

type AdvertRepository interface {
	Create(ctx context.Context, advert *model.Advert) error
	Update(ctx context.Context, id string, advert *model.Advert) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Advert, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Advert], error)
}

type ArchivedUserRepository interface {
	Create(ctx context.Context, user *model.User) error
	FindByID(ctx context.Context, id string) (*model.ArchivedUser, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ArchivedUser], error)
}

type ArchivedLinkedAccountRepository interface {
	Create(ctx context.Context, user *model.LinkedAccount) error
	FindByID(ctx context.Context, id string) (*model.ArchivedLinkedAccount, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ArchivedLinkedAccount], error)
}

type AuthTierRepository interface {
	Update(ctx context.Context, id string, user *model.AuthTier) error
	FindByID(ctx context.Context, id string) (*model.AuthTier, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AuthTier], error)
}

type BankRepository interface {
	Create(ctx context.Context, bank *model.Bank) error
	Update(ctx context.Context, id string, bank *model.Bank) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Bank, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Bank], error)
	FindByNameOrBICOrCode(ctx context.Context, bic, code, name string) (*model.Bank, error)
}

type DepartmentRepository interface {
	Create(ctx context.Context, department *model.Department) error
	Update(ctx context.Context, id string, department *model.Department) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Department, error)
	FindByName(ctx context.Context, name string) (*model.Department, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Department], error)
}

type PortalCardRepository interface {
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Card], error)
	ValidatePortalCard(ctx context.Context, names []string) (bool, error)
	ValidatePortalCardByID(ctx context.Context, ids []string) (bool, error)
}

type CpsUserRepository interface {
	Create(ctx context.Context, cpsUser *model.CPSUser) error
	Update(ctx context.Context, id string, cpsUser *model.CPSUser) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.CPSUser, error)
	FindByUsername(ctx context.Context, username string) (*model.CPSUser, error)
	GetPopulatedByID(ctx context.Context, id string) (*cps_user_dto.CpsUserResponse, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*cps_user_dto.CPSUserWithDepartment], error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*model.CPSUser, error)
	FindByEmail(ctx context.Context, email string) (*model.CPSUser, error)
}

type BankVaultRepository interface {
	Create(ctx context.Context, bankVault *model.BankVaultProduct) (string, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.BankVaultProduct], error)
	FindBankVaultByName(ctx context.Context, name string) error
	FindByID(ctx context.Context, id string) (*model.BankVaultProduct, error)
	Update(ctx context.Context, id string, bankVault *model.BankVaultProduct) error
	Delete(ctx context.Context, id string) (string, error)
	EnableOrDisable(ctx context.Context, id string, enable bool) error
}

type VaultGroupCategoryRepository interface {
	Create(ctx context.Context, vaultGroupCategory *model.VaultGroupCategory) (string, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.VaultGroupCategory], error)
	FindByID(ctx context.Context, id string) (*model.VaultGroupCategory, error)
	Update(ctx context.Context, id string, vaultGroupCategory *model.VaultGroupCategory) error
	Delete(ctx context.Context, id string) (string, error)
	EnableOrDisable(ctx context.Context, id string, enable bool) error
}

type DonationRepository interface {
	Create(ctx context.Context, donation *model.Donation) error
	Update(ctx context.Context, id string, donation *model.Donation) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*donation.DonationListResponse, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]donation.DonationListResponse], error)
}

type DonationCategoryRepository interface {
	Create(ctx context.Context, donationCategory *model.DonationCategory) error
	Update(ctx context.Context, id string, donationCategory *model.DonationCategory) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*donation_category.DonationCategoryListResponse, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]donation_category.DonationCategoryListResponse], error)
}

type DonationCompanyRepository interface {
	Create(ctx context.Context, donationCompany *model.DonationCompany) error
	Update(ctx context.Context, id string, donationCompany *model.DonationCompany) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*donation_company.DonationCompanyListResponse, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]donation_company.DonationCompanyListResponse], error)
}

type MiniAppRepository interface {
	Create(ctx context.Context, miniApp *model.MiniApp) error
	Update(ctx context.Context, id string, miniApp *model.MiniApp) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.MiniApp, error)
	Find(ctx context.Context, name string) (*model.MiniApp, error)
	FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.MiniApp], error)
	RunInTransaction(ctx context.Context, fn func(ctx context.Context) error) error
	FindByIDWithMerchant(ctx context.Context, id string) (*miniappdto.MiniAppResponse, error)
}

type EventRepository interface {
	Create(ctx context.Context, event *model.Event) error
	Update(ctx context.Context, id string, event *model.Event) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error

	FindByID(ctx context.Context, id string) (*model.Event, error)
	Find(ctx context.Context, name string) (*model.Event, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Event], error)
}

type FeedbackRepository interface {
	Create(ctx context.Context, feedback *model.Feedback) error
	FindByID(ctx context.Context, id string) (*feedback.FeedbackResponse, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*feedback.FeedbackResponse], error)
}

type LinkedAccountRepository interface {
	Create(ctx context.Context, account *model.LinkedAccount) error
	Update(ctx context.Context, id string, account *model.LinkedAccount) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.LinkedAccount, error)
	FindByAccountNumber(ctx context.Context, accountNumber string) (*model.LinkedAccount, error)
	FindByCustomerNumber(ctx context.Context, customer_number string) (*model.LinkedAccount, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.LinkedAccount], error)
}

type CityRepository interface {
	Create(ctx context.Context, city *model.City) error
	Update(ctx context.Context, id string, city *model.City) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.City, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.City], error)
}

type RegionRepository interface {
	Create(ctx context.Context, region *model.Region) error
	Update(ctx context.Context, id string, region *model.Region) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Region, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Region], error)
}

type DistrictRepository interface {
	Create(ctx context.Context, district *model.District) error
	Update(ctx context.Context, id string, district *model.District) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.District, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.District], error)
}

type BranchRepository interface {
	Create(ctx context.Context, branch *model.Branch) error
	Update(ctx context.Context, id string, branch *model.Branch) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Branch, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Branch], error)
}

type MiniAppMerchantRepository interface {
	Create(ctx context.Context, merchant *model.MiniAppMerchant) (*model.MiniAppMerchant, error)
	Update(ctx context.Context, id string, merchant *model.MiniAppMerchant) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.MiniAppMerchant, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.MiniAppMerchant], error)
	FindOne(ctx context.Context, filter bson.M) (*model.MiniAppMerchant, error)
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notification) error
	Update(ctx context.Context, id string, notification *model.Notification) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.Notification, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Notification], error)
	EnableDisableNotification(ctx context.Context, id string, enable bool) (*model.Notification, error)
	NotificationExists(ctx context.Context, notificationType string, forValue constants.NotificationFor, id *string) (bool, error)
}

type PasswordRuleRepository interface {
	Create(ctx context.Context, rule *model.PasswordRule) error
	Update(ctx context.Context, id string, rule *model.PasswordRule) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.PasswordRule, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.PasswordRule], error)
	FindCurrentRule(ctx context.Context) (*model.PasswordRule, error)
}

type ServiceDetailsRepository interface {
	Create(ctx context.Context, details *model.ServiceDetails) error
	Update(ctx context.Context, id string, details *model.ServiceDetails) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, projection bson.M, id string) (*model.ServiceDetails, error)
	FindAllWithPagination(ctx context.Context, projection bson.M, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ServiceDetails], error)
}

type ValidationRuleRepository interface {
	FindByID(ctx context.Context, id string) (*model.ValidationRule, error)
	Update(ctx context.Context, id string, rule *model.ValidationRule) error
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ValidationRule], error)
}

type KYCVerifierRepository interface {
	Update(ctx context.Context, id string, kyc *model.CustomerKYC) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	// Raw fetches for internal service logic
	FindByID(ctx context.Context, id string) (*model.CustomerKYC, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.CustomerKYC], error)
	// Populated responses for API
	FindByIDPopulated(ctx context.Context, id string) (*kyc_dto.KYCVerifierResponse, error)
	FindAllWithPaginationPopulated(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*kyc_dto.KYCVerifierResponse], error)
}

type WalletRepository interface {
	Create(ctx context.Context, wallet *model.Wallet) error
	Update(ctx context.Context, id string, wallet *model.Wallet) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error

	FindByID(ctx context.Context, id string) (*model.Wallet, error)
	Find(ctx context.Context, key, value string) (*model.Wallet, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Wallet], error)
}
type TopupRepository interface {
	Create(ctx context.Context, topup *model.Topup) error
	Update(ctx context.Context, id string, Topup *model.Topup) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error

	FindByID(ctx context.Context, id string) (*model.Topup, error)
	Find(ctx context.Context, key, value string) (*model.Topup, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Topup], error)
}
type ProductCodeRepository interface {
	FetchByID(ctx context.Context, id string) (*model.ProductCode, error)
	FindAllWithPagination(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ProductCode], error)
	FetchAll(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ProductCode], error)
	Update(ctx context.Context, productCode *model.ProductCode) error
}

type FaydaRepository interface {
	Update(ctx context.Context, user *model.User, isEnabled bool) error
	FindByUserCode(ctx context.Context, user_code string) (*model.User, error)
}

type CustomerRepository interface {
	FindByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, id string, data model.User) error
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.User], error)
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FetchLinkedAccount(ctx context.Context, customerNumber string) ([]*model.LinkedAccount, error)
}

type BulkServiceRepository interface {
	FindAllWithPagination(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.APPAccessList], error)
	Update(ctx context.Context, keys []string, state bool) error
	FindAll(ctx context.Context) ([]*model.APPAccessList, error)
}
type PermissionRepository interface {
	Create(ctx context.Context, permissionGroup *model.PermissionGroup) error
	Update(ctx context.Context, id string, permissionGroup *model.PermissionGroup) error
	Delete(ctx context.Context, id string) error
	// FindByID(ctx context.Context, id string) (*model.PermissionGroup, error)
	FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.PermissionGroup], error)

	// Permission category operations
	ValidatePermissionCategories(ctx context.Context, categoryIDs []string) ([]string, error)
	GetAllPermissionCategoriesWithPermissions(ctx context.Context) ([]*model.PermissionCategory, error)
	GetAllPermissionCategories(ctx context.Context, card string) ([]*model.PermissionCategory, error)
	// Permission group operations
	ValidatePermissionGroups(ctx context.Context, groupIDs []string) ([]string, error)
	CheckPermissionGroupExists(groupName string) bool
	GetPermissionGroup(groupName string) (*model.PermissionGroup, error)
	GetPermissionGroupById(ctx context.Context, id string) (*model.PermissionGroup, error)
}

// ExternalCallServices type alias for external call services
type ExternalCallServices = external_call.ExternalCallServices

type ArticleRepository interface {
	CreateArticle(ctx context.Context, article *model.NewsArticle) error
	UpdateArticle(ctx context.Context, article *model.NewsArticle, id string) error
	DeleteArticle(ctx context.Context, id string) error
	PublishUnpublishArticle(ctx context.Context, id string, isPublished bool) error
}
type ArticleCategoryRepository interface {
	CreateArticleCategory(ctx context.Context, category *model.NewsCategoryModel) error
	UpdateArticleCategory(ctx context.Context, category *model.NewsCategoryModel, id string) error
	DeleteArticleCategory(ctx context.Context, id string) error
	EnableOrDisableArticleCategory(ctx context.Context, id string, enable bool) error
}

type ShortVideoRepository interface {
	Create(ctx context.Context, shortVideo *model.ShortVideo) error
	Update(ctx context.Context, shortVideo *model.ShortVideo, id string) error
	Delete(ctx context.Context, id string) error
	PublishUnpublish(ctx context.Context, id string, isPublished bool) error
}

type NewsTagRepository interface {
	FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.NewsTag], error)
	Get(ctx context.Context, id string) (*model.NewsTag, error)
	Create(ctx context.Context, tagName []string) error
	Update(ctx context.Context, id string, tagName string) error
	Delete(ctx context.Context, id string) error
	FindByNames(ctx context.Context, names []string) (*model.NewsTag, error)
}

type NewsCategoryRepository interface {
	FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.NewsCategory], error)
	Get(ctx context.Context, id string) (*model.NewsCategory, error)
	Create(ctx context.Context, categoryName []string) error
	Update(ctx context.Context, id string, categoryName string) error
	Delete(ctx context.Context, id string) error
	FindByNames(ctx context.Context, names []string) (*model.NewsCategory, error)
}

type NewsTagsRepository interface {
	Create(ctx context.Context, newsTag *model.NewsTags) error
	Update(ctx context.Context, newsTag *model.NewsTags, id string) error
	Delete(ctx context.Context, id string) error
	EnableDisable(ctx context.Context, id string, isEnable bool) error
}

type IconRepository interface {
	Create(ctx context.Context, icon *model.Icon) error
	Update(ctx context.Context, id string, icon *model.Icon) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Icon, error)
	FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.Icon], error)
}
