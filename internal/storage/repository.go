package storage

import (
	"context"

	session "cbe-super-app-cps-action/grpc"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"

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
	FindById(ctx context.Context, id string) (*model.User, error)
	FindByUserCode(ctx context.Context, userCode string) (*model.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*model.User, error)
	FindByDeviceUUID(ctx context.Context, deviceUUID string) (*model.User, error)
	Update(ctx context.Context, id string, update *model.User) error
}

type HQRepository interface {
	FindOne(ctx context.Context, filter bson.M) (*model.HQ, error)
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

type AccountAPIPort interface {
	LookupAccountByPhone(ctx context.Context, phoneNumber string, PhoneLookupUrl string) (bool, error)
}

type CPSActionRepository interface {
	Save(ctx context.Context, cpsAction *model.CPSAction) error
	FindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (*types.PaginatedResponse[[]*model.CPSAction], error)
	FindOne(ctx context.Context, filter model.CPSAction) (*model.CPSAction, error)
	Update(ctx context.Context, id string, update model.CPSAction) error
	Delete(ctx context.Context, id string) error
	CPSActionExists(ctx context.Context, user model.CheckCPSAction) (bool, error)
}

// Avatar persistence
type AvatarRepository interface {
	Create(ctx context.Context, avatar *model.Avatar) error
	Update(ctx context.Context, id string, avatar *model.Avatar) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Avatar, error)
	FindAll(ctx context.Context, filter bson.M, projection bson.M) ([]*model.Avatar, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Avatar], error)
}

// BPSUser persistence
type BPSUserRepository interface {
	GetByUserCode(ctx context.Context, userCode string) (*model.BPSUser, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.BPSUser], error)
	Update(ctx context.Context, BpsUser *model.BPSUser) error
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
	GetBranchByCode(ctx context.Context, branchCode string) (*model.Branch, error)
	CreateBranch(ctx context.Context, branch *model.Branch) error
	UpdateBranch(ctx context.Context, id string, branch *model.Branch) error
	DeleteBranch(ctx context.Context, id string) error
	EnableOrDisableBranch(ctx context.Context, id string, enable bool) error
	FindBranchByID(ctx context.Context, id string) (*model.Branch, error)
	FindAllBranchesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Branch], error)

	CreateCity(ctx context.Context, city *model.City) error
	UpdateCity(ctx context.Context, id string, city *model.City) error
	DeleteCity(ctx context.Context, id string) error
	EnableOrDisableCity(ctx context.Context, id string, enable bool) error
	FindCityByID(ctx context.Context, id string) (*model.City, error)
	FindAllCitiesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.City], error)

	CreateRegion(ctx context.Context, region *model.Region) error
	UpdateRegion(ctx context.Context, id string, region *model.Region) error
	DeleteRegion(ctx context.Context, id string) error
	EnableOrDisableRegion(ctx context.Context, id string, enable bool) error
	FindRegionByID(ctx context.Context, id string) (*model.Region, error)
	FindAllRegionsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Region], error)

	CreateDistrict(ctx context.Context, district *model.District) error
	UpdateDistrict(ctx context.Context, id string, district *model.District) error
	DeleteDistrict(ctx context.Context, id string) error
	EnableOrDisableDistrict(ctx context.Context, id string, enable bool) error
	FindDistrictByID(ctx context.Context, id string) (*model.District, error)
	FindAllDistrictsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.District], error)
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
}

type PortalCardRepository interface {
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Card], error)
	ValidatePortalCard(ctx context.Context, names []string) (bool, error)
}

type ColorRepository interface {
	Create(ctx context.Context, color *model.Color) error
	Update(ctx context.Context, id string, color *model.Color) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Color, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Color], error)
}

type CpsUserRepository interface {
	Create(ctx context.Context, cpsUser *model.CPSUser) error
	Update(ctx context.Context, id string, cpsUser *model.CPSUser) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.CPSUser, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.CPSUser], error)
}

type DonationRepository interface {
	Create(ctx context.Context, donation *model.Donation) error
	Update(ctx context.Context, id string, donation *model.Donation) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.Donation, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Donation], error)
}

type DonationCategoryRepository interface {
	Create(ctx context.Context, donationCategory *model.DonationCategory) error
	Update(ctx context.Context, id string, donationCategory *model.DonationCategory) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.DonationCategory, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.DonationCategory], error)
}

type DonationCompanyRepository interface {
	Create(ctx context.Context, donationCompany *model.DonationCompany) error
	Update(ctx context.Context, id string, donationCompany *model.DonationCompany) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.DonationCompany, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.DonationCompany], error)
}

type MiniAppRepository interface {
	Create(ctx context.Context, miniApp *model.MiniApp) error
	Update(ctx context.Context, miniApp *model.MiniApp) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.MiniApp, error)
	FindByCode(ctx context.Context, code string) (*model.MiniApp, error)
	FindAllWithPagination(ctx context.Context, department string, filterParam types.Filter) (*types.PaginatedResponse[[]*model.MiniApp], error)
}

type EventRepository interface {
	Create(ctx context.Context, event *model.Event) error
	Update(ctx context.Context, id string, event *model.Event) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error

	FindByID(ctx context.Context, id string) (*model.Event, error)
	FindByName(ctx context.Context, name string) (*model.Event, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Event], error)
}

type FeedbackRepository interface {
	Create(ctx context.Context, feedback *model.Feedback) error
	FindByID(ctx context.Context, id string) (*model.Feedback, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Feedback], error)
}

type IconRepository interface {
	Create(ctx context.Context, icon *model.Icon) error
	Update(ctx context.Context, id string, icon *model.Icon) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Icon, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Icon], error)
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
	Create(ctx context.Context, merchant *model.MiniAppMerchant) error
	Update(ctx context.Context, id string, merchant *model.MiniAppMerchant) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.MiniAppMerchant, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.MiniAppMerchant], error)
	DetailMiniAppByID(ctx context.Context, id string) (*model.MiniAppMerchant, error)
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notification) error
	Update(ctx context.Context, id string, notification *model.Notification) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.Notification, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Notification], error)
}

type PasswordRuleRepository interface {
	Create(ctx context.Context, rule *model.PasswordRule) error
	Update(ctx context.Context, id string, rule *model.PasswordRule) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.PasswordRule, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.PasswordRule], error)
}

type ServiceDetailsRepository interface {
	Create(ctx context.Context, details *model.ServiceDetails) error
	Update(ctx context.Context, id string, details *model.ServiceDetails) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.ServiceDetails, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ServiceDetails], error)
}

type ValidationRuleRepository interface {
	GetAccountValidationByID(ctx context.Context, id string) (*model.ValidationRule, error)
	UpdateAccountValidation(ctx context.Context, id string, rule *model.ValidationRule) error
	GetAllAccountValidation(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ValidationRule], error)
}

type WalletRepository interface {
	Create(ctx context.Context, wallet *model.Wallet) error
	Update(ctx context.Context, id string, wallet *model.Wallet) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.Wallet, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Wallet], error)
}

