package storage

import (
	ussd_merchant_dto "cbe-super-app-cps-action/internal/constants/dto/ussd_merchant"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	"context"
	"time"

	session "cbe-super-app-cps-action/grpc"
	"cbe-super-app-cps-action/internal/constants"
	access_list_segmentation_dto "cbe-super-app-cps-action/internal/constants/dto/access_list_segmentation"
	actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/action_role"
	bpsActionDto "cbe-super-app-cps-action/internal/constants/dto/bps_action"
	actionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	cps_actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/cps_action_role"
	cps_user_dto "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	cpsuser "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	customer_dto "cbe-super-app-cps-action/internal/constants/dto/customer"
	"cbe-super-app-cps-action/internal/constants/dto/donation"
	"cbe-super-app-cps-action/internal/constants/dto/donation_category"
	"cbe-super-app-cps-action/internal/constants/dto/donation_company"
	"cbe-super-app-cps-action/internal/constants/dto/feedback"
	kyc_dto "cbe-super-app-cps-action/internal/constants/dto/kyc_verifier"
	transaction_dto "cbe-super-app-cps-action/internal/constants/dto/transaction"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage/external_call"

	member "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"

	unlink_dto "cbe-super-app-cps-action/internal/constants/dto/unlink"
	imodel "cbe-super-app-cps-action/internal/constants/model"

	donation_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/donation"
	mini_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/mini_app"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	bps_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"

	bps_user_dto "cbe-super-app-cps-action/internal/constants/dto/bps_user"

	"go.mongodb.org/mongo-driver/v2/bson"
	"google.golang.org/grpc"
)

type RoleRepository interface {
	Exists(ctx context.Context, id string) (bool, error)
	ExistsMany(ctx context.Context, ids []string) (bool, error)
	Create(ctx context.Context, role *imodel.Role) error
	Update(ctx context.Context, id string, role *imodel.Role) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	SoftDelete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*imodel.Role, error)
	FindByName(ctx context.Context, name string) (*imodel.Role, error)
	FindByCode(ctx context.Context, code string) (*imodel.Role, error)
	FindByRole(ctx context.Context, jobTitle string) (*imodel.Role, error)
	FindAll(ctx context.Context) (*[]imodel.Role, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.Role], error)
	FindByFilterKey(ctx context.Context, field string, value string) (*imodel.Role, error)
}

type JobRoleRepository interface {
	Create(ctx context.Context, role *imodel.JobRole) error
	Update(ctx context.Context, id string, role *imodel.JobRole) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	SoftDelete(ctx context.Context, id string) error
	ExistsMany(ctx context.Context, codes []string) (bool, error)
	FindByID(ctx context.Context, id string) (*imodel.JobRole, error)
	FindByCode(ctx context.Context, code string) (*imodel.JobRole, error)
	Find(ctx context.Context, filter bson.M) (*imodel.JobRole, error)
	FindAll(ctx context.Context) (*[]imodel.JobRole, error)
	FindByName(ctx context.Context, name string) (*imodel.JobRole, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.JobRole], error)
}

type UnlinkAccount interface {
	GetUserByAccount(ctx context.Context, accNumber string) (*model.ArchivedUser, error)
	GetAllArchivedUser(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[*model.ArchivedUser], error)
	Delete(ctx context.Context, id string) error
	Authorize(ctx context.Context, cpsAction model.CPSAction) (model.ArchivedUser, error)
}

// ServicesRepository manages CRUD for Services catalog
type ServicesRepository interface {
	Create(ctx context.Context, service *model.Service) error
	Update(ctx context.Context, id string, service *model.Service) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Service, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Service], error)
	FindAllServiceListWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.ServiceList], error)
	CheckServiceExistence(ctx context.Context, serviceCode, serviceKey, serviceName string) (bool, error)
}

type OTPRepository interface {
	Save(ctx context.Context, otp *model.OTP) error
	Find(ctx context.Context, filte bson.M) (*model.OTP, error)
	Delete(ctx context.Context, id string) error
}

type UserRepository interface {
	Save(ctx context.Context, user *member.User) error
	GetUserByAccount(ctx context.Context, accNumber string) (*member.User, error)
	Delete(ctx context.Context, id string) error
	DeleteHard(ctx context.Context, id string) error
	FindById(ctx context.Context, id string) (*member.User, error)
	FindByUserCode(ctx context.Context, userCode string) (*member.User, error)
	FindByCustomerNumber(ctx context.Context, customerNumber string) (*member.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*member.User, error)
	FindByDeviceUUID(ctx context.Context, deviceUUID string) (*member.User, error)
	Update(ctx context.Context, id string, update *member.User) error
}

type HQRepository interface {
	FindByID(ctx context.Context, id string) (*model.HQ, error)
	Find(ctx context.Context, projections ...bson.M) (*model.HQ, error)
	Update(ctx context.Context, field string, value interface{}, now time.Time) error
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.HQ], error)
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
	FindByKeys(ctx context.Context, keys []string) (map[string]string, error)
	FindAllWithPagination(ctx context.Context, department string, filterParam types.Filter) (*types.PaginatedResponse[[]model.APPAccessList], error)
	FindAllByKeys(ctx context.Context, keys []string) ([]model.APPAccessList, error)
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
type CPSActionRoleRepository interface {
	Create(ctx context.Context, actionRole *imodel.CPSActionRole) error
	UpdateByActionCode(ctx context.Context, actionCode string, actionRole *imodel.CPSActionRole) error
	EnableOrDisableByActionCode(ctx context.Context, actionCode string, enable bool) error
	FindByActionCode(ctx context.Context, actionCode string) (*cps_actionrole_dto.GetActionRoleByActionCodeRes, error)
	UpdateActionList(ctx context.Context, actionCode, portalCard string, status bool) error
	FindAllAccessListWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.CPSActionList], error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.CPSActionRoleResposne], error)
	FindByActionName(ctx context.Context, actionName string) (*imodel.CPSActionRole, error)
	FindByActionNameAndPortalCard(ctx context.Context, actionName, portalCard string) (*imodel.CPSActionRole, error)
	FindApproverByActionName(ctx context.Context, actionName, role_code string) (imodel.CPSActionApproveIndex, error)
	FindByActionCodeOne(ctx context.Context, actionCode string) (*imodel.CPSActionRole, error)
}
type DeviceVersionControlRepository interface {
	Save(ctx context.Context, deviceVersionControl model.DeviceVersionControl) error
	FindOne(ctx context.Context, platform, deviceVersion string) (model.DeviceVersionControl, error)
	FindByID(ctx context.Context, id string) (model.DeviceVersionControl, error)
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
	SanitizedFindAllWithPaginationForApprover(ctx context.Context, userID string, filterParam types.Filter, RAList []string) (*types.PaginatedResponse[[]*model.CPSAction], error)
	SanitizedFindAllWithPaginationForAuditor(ctx context.Context, userID string, filterParam types.Filter, RAList []string) (*types.PaginatedResponse[[]*model.CPSAction], error)
	SanitizedFindAllWithPaginationCPSActions(ctx context.Context, userID, role string, filterParam types.Filter, RAList []string) (*types.PaginatedResponse[[]*model.CPSAction], error)
	SanitizedFindOne(ctx context.Context, filter bson.M) (*model.CPSAction, error)
	Update(ctx context.Context, actionCode string, update model.CPSAction) (*model.CPSAction, error)
	UpdateByActionCode(ctx context.Context, actionCode string, update model.CPSAction) (*model.CPSAction, error)
	UpdateCustome(ctx context.Context, filter, update bson.M) error
	Delete(ctx context.Context, id string) error
	GetCountByDepartment(ctx context.Context, department string) (*actionDto.CPSActionCountResponse, error)
}

// Avatar persistence
type AvatarRepository interface {
	Create(ctx context.Context, avatar *model.Avatar) error
	Update(ctx context.Context, id string, avatar *model.Avatar) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Avatar, error)
	Find(ctx context.Context, filter bson.M, projection bson.M) (*model.Avatar, error)
	FindAll(ctx context.Context, filter bson.M, projection bson.M) ([]model.Avatar, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Avatar], error)
}

// BPSUser persistence
type BPSUserRepository interface {
	GetByUserCode(ctx context.Context, userCode string) (*bps_model.BPSUser, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]bps_user_dto.BPSUserResposenDTO], error)
	Update(ctx context.Context, BpsUser *bps_model.BPSUser) error
	Create(ctx context.Context, BpsUser bps_model.BPSUser) error
	FindByFilterKey(ctx context.Context, field, value string) (*bps_model.BPSUser, error)
	FindByOr(ctx context.Context, phone, email, username string) (*bps_model.BPSUser, error)
}
type BPSActionRepository interface {
	Save(ctx context.Context, cpsAction *bps_model.BPSAction) error
	FindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (types.PaginatedResponse[[]bps_model.BPSAction], error)
	FindOne(ctx context.Context, filter bson.M) (*bps_model.BPSAction, error)
	SanitizedFindAllWithPagination(ctx context.Context, filterParam types.Filter, department string) (*types.PaginatedResponse[[]*bps_model.BPSAction], error)
	SanitizedFindAllWithPaginationForApprover(ctx context.Context, userID string, filterParam types.Filter, RAList []string) (*types.PaginatedResponse[[]*bps_model.BPSAction], error)
	SanitizedFindAllWithPaginationForAuditor(ctx context.Context, userID string, filterParam types.Filter, RAList []string) (*types.PaginatedResponse[[]*bps_model.BPSAction], error)
	SanitizedFindAllWithPaginationBPSActions(ctx context.Context, userID, role string, filterParam types.Filter, RAList []string) (*types.PaginatedResponse[[]*bps_model.BPSAction], error)
	SanitizedFindOne(ctx context.Context, filter bson.M) (*bps_model.BPSAction, error)
	Update(ctx context.Context, actionCode string, update bps_model.BPSAction) (*bps_model.BPSAction, error)
	UpdateByActionCode(ctx context.Context, actionCode string, update bps_model.BPSAction) (*bps_model.BPSAction, error)
	UpdateCustome(ctx context.Context, filter, update bson.M) error
	Delete(ctx context.Context, id string) error
	GetCountByDepartment(ctx context.Context, department string) (*bpsActionDto.BPSActionCountResponse, error)
	GetBPSActionByUserID(ctx context.Context, userID string, filter types.Filter) (types.PaginatedResponse[[]bps_model.BPSAction], error)
}

type BudgetCategoryRepository interface {
	CreateBudgetCategory(ctx context.Context, budgetCategory *model.BudgetCategory) error
	UpdateBudgetCategory(ctx context.Context, id string, budgetCategory *model.BudgetCategory) error
	FindBudgetCategoryByID(ctx context.Context, id string) (*model.BudgetCategory, error)
	FindAllBudgetCategories(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]model.BudgetCategory], error)
	DeleteBudgetCategory(ctx context.Context, id string) error
	EnableOrDisableBudgetCategory(ctx context.Context, id string, enable bool) error
	FindByName(ctx context.Context, name string) (*model.BudgetCategory, error)
}

// AmountBasedAuth persistence
type AmountBasedAuthRepository interface {
	Update(ctx context.Context, id string, update *model.AuthTier) error
	FindAll(ctx context.Context, filter bson.M, projection bson.M) ([]model.AuthTier, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.AuthTier], error)
	FindByID(ctx context.Context, id string) (*model.AuthTier, error)
}

// AccountBlock persistence
type AccountBlockRepository interface {
	CreateBranch(ctx context.Context, branch *model.AccountBlock) error
	DeleteBranch(ctx context.Context, id string) error
	FindAllBranchesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error)
	EnableOrDisableBranches(ctx context.Context, ids []string, reason string, enabled bool) error
	GetBranchesByIds(ctx context.Context, ids []string) ([]*model.AccountBlock, error)

	CreateCity(ctx context.Context, city *model.AccountBlock) error
	DeleteCity(ctx context.Context, id string) error
	FindAllCitiesWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error)
	EnableOrDisableCities(ctx context.Context, ids []string, reason string, enabled bool) error
	GetCitiesByIds(ctx context.Context, ids []string) ([]*model.AccountBlock, error)

	CreateRegion(ctx context.Context, region *model.AccountBlock) error
	DeleteRegion(ctx context.Context, id string) error
	FindAllRegionsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error)
	EnableOrDisableRegions(ctx context.Context, ids []string, reason string, enabled bool) error
	GetRegionsByIds(ctx context.Context, ids []string) ([]*model.AccountBlock, error)

	CreateDistrict(ctx context.Context, district *model.AccountBlock) error
	DeleteDistrict(ctx context.Context, id string) error
	FindAllDistrictsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error)
	EnableOrDisableDistricts(ctx context.Context, ids []string, reason string, enabled bool) error
	GetDistrictsByIds(ctx context.Context, ids []string) ([]*model.AccountBlock, error)
	GetAccountBlockDetails(ctx context.Context, id string) ([]*model.CPSAction, error)
}

type AdvertRepository interface {
	Create(ctx context.Context, advert *model.Advert) error
	Update(ctx context.Context, id string, advert *model.Advert) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Advert, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Advert], error)
	FindByTitle(ctx context.Context, title string) (*model.Advert, error)
}

type ArchivedUserRepository interface {
	Create(ctx context.Context, user *member.User) error
	FindByID(ctx context.Context, id string) (*model.ArchivedUser, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.ArchivedUser], error)
	FindAllArchievedUsersWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*unlink_dto.ArchivedUserResponse], error)
}

type ArchivedLinkedAccountRepository interface {
	Create(ctx context.Context, user *model.LinkedAccount) error
	FindByID(ctx context.Context, id string, isUserId bool) (*model.ArchivedLinkedAccount, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.ArchivedLinkedAccount], error)
}

type AuthTierRepository interface {
	Update(ctx context.Context, id string, user *model.AuthTier) error
	FindByID(ctx context.Context, id string) (*model.AuthTier, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.AuthTier], error)
}

type BankRepository interface {
	Create(ctx context.Context, bank *model.Bank) error
	Update(ctx context.Context, id string, bank *model.Bank) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Bank, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Bank], error)
	FindByNameOrBIC(ctx context.Context, bic, name string) (*model.Bank, error)
	FindByBIC(ctx context.Context, bic string) (*model.Bank, error)
}

type DepartmentRepository interface {
	Create(ctx context.Context, department *model.Department) error
	Update(ctx context.Context, id string, department *model.Department) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Department, error)
	FindByName(ctx context.Context, name string) (*model.Department, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (types.PaginatedResponse[[]model.Department], error)
}

type PortalCardRepository interface {
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Card], error)
	ValidatePortalCard(ctx context.Context, names []string) (bool, error)
	ValidatePortalCardByID(ctx context.Context, ids []string) (bool, error)
}

type CpsUserRepository interface {
	Create(ctx context.Context, cpsUser *imodel.CPSUser) error
	Update(ctx context.Context, id string, cpsUser *imodel.CPSUser) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*imodel.CPSUser, error)
	FindByUsername(ctx context.Context, username string) (*imodel.CPSUser, error)
	GetPopulatedByID(ctx context.Context, id string) (*cps_user_dto.CpsUserResponse, error)
	GetPopulatedWithRole(ctx context.Context, userCode string) (*cpsuser.CpsUserPopulatedResponse, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*cps_user_dto.CPSUserWithDepartment], error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (*imodel.CPSUser, error)
	FindByEmail(ctx context.Context, email string) (*imodel.CPSUser, error)
	FindByEmailOrPhoneNumberOrUserName(ctx context.Context, email string, phoneNumber string, username string) (*imodel.CPSUser, error)
}

type BankVaultRepository interface {
	Create(ctx context.Context, bankVault *model.BankVaultProduct) (string, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.BankVaultProduct], error)
	FindAllBankLockedVaultsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.LockedVault], error)
	FindAllGroupVaultWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.GroupVault], error)
	FindBankVaultByName(ctx context.Context, name string) error
	FindByID(ctx context.Context, id string) (*model.BankVaultProduct, error)
	Update(ctx context.Context, id string, bankVault *model.BankVaultProduct) error
	Delete(ctx context.Context, id string) (string, error)
	EnableOrDisable(ctx context.Context, id string, enable bool) error
}

type VaultCategoryRepository interface {
	Create(ctx context.Context, vaultCategory *imodel.VaultCategory) (string, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.VaultCategory], error)
	FindByID(ctx context.Context, id string) (*imodel.VaultCategory, error)
	FindByName(ctx context.Context, groupName string) (*imodel.VaultCategory, error)
	Update(ctx context.Context, id string, vaultCategory *imodel.VaultCategory) error
	Delete(ctx context.Context, id string) (string, error)
	EnableOrDisable(ctx context.Context, id string, enable bool) error

	// Transaction
	FindAllTransactionsWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.VaultTransaction], error)
	FindVaultTransaction(ctx context.Context, id string) (*imodel.VaultTransaction, error)

	// withdrawal request
	CreateWithdrawalRequest(ctx context.Context, withdrawal *imodel.Withdrawal) error
	UpdateWithdrawalRequest(ctx context.Context, id string, status string) error
	GetAllWithdrawalRequests(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.Withdrawal], error)
	GetWithdrawalRequest(ctx context.Context, id string) (*imodel.Withdrawal, error)
}

type VaultAmountTierRepository interface {
	Create(ctx context.Context, vaultAmountTier *model.VaultAmountTier) (string, error)
	FindByID(ctx context.Context, id string) (*model.VaultAmountTier, error)
	Update(ctx context.Context, id string, vaultAmountTier *model.VaultAmountTier) error
	Delete(ctx context.Context, id string) (string, error)
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.VaultAmountTier], error)
}

type DonationRepository interface {
	Create(ctx context.Context, donation *imodel.Donation) error
	Update(ctx context.Context, id string, donation *imodel.Donation) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*donation.DonationListResponse, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]donation.DonationListResponse], error)
}

type DonationCategoryRepository interface {
	Create(ctx context.Context, donationCategory *donation_model.DonationCategory) error
	Update(ctx context.Context, id string, donationCategory *donation_model.DonationCategory) error
	EnableDisable(ctx context.Context, id string, enable bool) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*donation_category.DonationCategoryListResponse, error)
	FindByName(ctx context.Context, name string) (*donation_category.DonationCategoryListResponse, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]donation_category.DonationCategoryListResponse], error)
}

type DonationCompanyRepository interface {
	Create(ctx context.Context, donationCompany *donation_model.DonationCompany) error
	Update(ctx context.Context, id string, donationCompany *donation_model.DonationCompany) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*donation_company.DonationCompanyListResponse, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]donation_company.DonationCompanyListResponse], error)
	FindByAccountNumber(ctx context.Context, accountNumber string) (*donation_model.DonationCompany, error)
}

type MiniAppRepository interface {
	Create(ctx context.Context, miniApp *mini_model.MiniApp) error
	Update(ctx context.Context, id string, miniApp *mini_model.MiniApp) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	DisableManyByMerchantIDs(ctx context.Context, merchantID string) error
	DeleteManyByMerchantIDs(ctx context.Context, merchantID string) error
}

type EventRepository interface {
	Create(ctx context.Context, event *model.Event) error
	Update(ctx context.Context, id string, event *model.Event) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.Event, error)
	Find(ctx context.Context, name string) (*model.Event, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Event], error)
}

type FeedbackRepository interface {
	Create(ctx context.Context, feedback *imodel.Feedback) error
	CreateSurveyFeedback(ctx context.Context, surveyFeedback *local_model.SurveyFeedback) error
	FindByID(ctx context.Context, id string) (*feedback.FeedbackResponse, error)

	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]local_model.Feedback], error)
	FindFeedbackByID(ctx context.Context, id string) (*local_model.Feedback, error)

	FindAllCustomerFeedbacks(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.CustomerFeedback], error)
	FindCustomerFeedbackByID(ctx context.Context, id string) (*imodel.CustomerFeedback, error)
	FindAllSurveyFeedbacks(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.SurveyFeedback], error)
	FindSurveyFeedbackByID(ctx context.Context, id string) (*imodel.SurveyFeedback, error)
}

type LinkedAccountRepository interface {
	Create(ctx context.Context, account *model.LinkedAccount) error
	Update(ctx context.Context, id string, account *model.LinkedAccount) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.LinkedAccount, error)
	FindByAccountNumber(ctx context.Context, accountNumber string) (*model.LinkedAccount, error)
	FindByCustomerNumber(ctx context.Context, customer_number string) (*model.LinkedAccount, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.LinkedAccount], error)
}

type EcommerceMerchantRepository interface {
	Create(ctx context.Context, merchant *model.EcommerceMerchant) (*model.EcommerceMerchant, error)
	Update(ctx context.Context, id string, merchant *model.EcommerceMerchant) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.EcommerceMerchant, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.EcommerceMerchant], error)
	FindOne(ctx context.Context, filter bson.M) (*model.EcommerceMerchant, error)
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notification) error
	Update(ctx context.Context, id string, notification *model.Notification) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.Notification, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Notification], error)
	EnableDisableNotification(ctx context.Context, id string, enable bool) (*model.Notification, error)
	NotificationExists(ctx context.Context, notificationType string, forValue constants.NotificationFor, id *string) (bool, error)
}

type PasswordRuleRepository interface {
	Create(ctx context.Context, rule *local_model.PasswordRule) error
	Update(ctx context.Context, id string, rule *local_model.PasswordRule) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*local_model.PasswordRule, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]local_model.PasswordRule], error)
	FindCurrentRule(ctx context.Context) (*local_model.PasswordRule, error)
}

type ServiceDetailsRepository interface {
	Create(ctx context.Context, details *imodel.ServiceDetails) error
	Update(ctx context.Context, id string, details *imodel.ServiceDetails) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, projection bson.M, id string) (*imodel.ServiceDetails, error)
	FindAllWithPagination(ctx context.Context, projection bson.M, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.ServiceDetails], error)
}

type ValidationRuleRepository interface {
	FindByID(ctx context.Context, id string) (*model.ValidationRule, error)
	Update(ctx context.Context, id string, rule *model.ValidationRule) error
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.ValidationRule], error)
}

type KYCVerifierRepository interface {
	Update(ctx context.Context, id string, kyc *model.CustomerKYC) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	// Raw fetches for internal service logic
	FindByID(ctx context.Context, id string) (*model.CustomerKYC, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.CustomerKYC], error)
	// Populated responses for API
	FindByIDPopulated(ctx context.Context, id string) (*kyc_dto.KYCVerifierResponse, error)
	FindAllWithPaginationPopulated(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]kyc_dto.KYCVerifierResponse], error)
}

type WalletRepository interface {
	Create(ctx context.Context, wallet *local_model.Wallet) error
	Update(ctx context.Context, id string, wallet *local_model.Wallet) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error

	FindByID(ctx context.Context, id string) (*local_model.Wallet, error)
	Find(ctx context.Context, key, value string) (*local_model.Wallet, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]local_model.Wallet], error)
	FindByIDForGRPC(ctx context.Context, id string) (*local_model.GRPCWallet, error)
	FindAllWithPaginationForGRPC(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]local_model.Wallet], error)
}
type TopupRepository interface {
	Create(ctx context.Context, topup *model.Topup) error
	Update(ctx context.Context, id string, Topup *model.Topup) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByKeyValue(ctx context.Context, key string, value string) (*model.Topup, error)

	FindByID(ctx context.Context, id string) (*model.Topup, error)
	FindByOr(ctx context.Context, filter bson.M) (model.Topup, error)
	Find(ctx context.Context, key, value string) (*model.Topup, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.Topup], error)
}
type ProductCodeRepository interface {
	FetchByID(ctx context.Context, id string) (*model.ProductCode, error)
	FindAllWithPagination(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ProductCode], error)
	FetchAll(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ProductCode], error)
	Update(ctx context.Context, productCode *model.ProductCode) error
	FindByName(ctx context.Context, name string) (*model.ProductCode, error)
	FindByPRD(ctx context.Context, cbePRD, cbeIFBPRD string) ([]*model.ProductCode, error)
}

type FaydaRepository interface {
	Update(ctx context.Context, user *member.User, isEnabled bool) error
	FindByUserCode(ctx context.Context, user_code string) (*member.User, error)
}

type CustomerRepository interface {
	FindByID(ctx context.Context, id string) (*member.User, error)
	Update(ctx context.Context, id string, data member.User) error
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*customer_dto.CustomerListResponse], error)
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FetchLinkedAccount(ctx context.Context, id string) ([]model.LinkedAccount, error)
	FindCustomerDetailByID(ctx context.Context, id string) (*customer_dto.CustomerDetailResponse, error)
	SearchCustomerByCIForAccountNumber(ctx context.Context, number string) (*customer_dto.CustomerListResponse, error)
	FindCustomerByIDs(ctx context.Context, ids []string) ([]member.User, error)
	FindCustomerByID(ctx context.Context, id string) (*member.User, error)
}

type BulkServiceRepository interface {
	FindAllWithPagination(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]model.APPAccessList], error)
	Update(ctx context.Context, keys []string, state bool) error
	FindAll(ctx context.Context) ([]model.APPAccessList, error)
}
type PermissionRepository interface {
	Create(ctx context.Context, permissionGroup *model.PermissionGroup) error
	Update(ctx context.Context, id string, permissionGroup *model.PermissionGroup) error
	Delete(ctx context.Context, id string) error
	FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.PermissionGroup], error)
	FindAllGroupsWithPagination(ctx context.Context, departmentId string, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.PermissionGroup], error)
	ValidatePermissionCategories(ctx context.Context, categoryIDs []string) ([]string, error)
	GetAllPermissionCategoriesWithPermissions(ctx context.Context) ([]*model.PermissionCategory, error)
	GetAllPermissionCategories(ctx context.Context, card string) ([]*model.PermissionCategory, error)
	ValidatePermissionGroups(ctx context.Context, groupIDs []string) ([]string, error)
	CheckPermissionGroupExists(groupName string) bool
	GetPermissionGroup(groupName string) (*model.PermissionGroup, error)
	GetPermissionGroupById(ctx context.Context, id string) (*model.PermissionGroup, error)
	FindByIDPopulated(ctx context.Context, id string) (cps_user_dto.PermissionGroupResponse, error)
	GetPopulatedPermissionCategories(ctx context.Context, categoryIDs []string) ([]cps_user_dto.PermissionCategoryResponse, error)
	GetPopulatedPermissionGroups(ctx context.Context, groupIDs []string) ([]cps_user_dto.PermissionGroupResponse, error)
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
	FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]model.NewsTag], error)
	Get(ctx context.Context, id string) (*model.NewsTag, error)
	Create(ctx context.Context, tagName []string) error
	Update(ctx context.Context, id string, tagName string) error
	Delete(ctx context.Context, id string) error
	FindByNames(ctx context.Context, names []string) (*model.NewsTag, error)
}

type NewsCategoryRepository interface {
	FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]model.NewsCategory], error)
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
	FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]model.Icon], error)
}

type BPSActionRoleRepository interface {
	Create(ctx context.Context, actionRole *imodel.BPSActionRole) error
	UpdateByActionCode(ctx context.Context, actionCode string, actionRole *imodel.BPSActionRole) error
	EnableOrDisableByActionCode(ctx context.Context, actionCode string, enable bool) error
	FindByActionCode(ctx context.Context, actionCode string) (*actionrole_dto.GetActionRoleByActionCodeRes, error)
	UpdateActionList(ctx context.Context, actionCode string, status bool) error
	FindAllAccessListWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.BPSActionList], error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*imodel.BPSActionRoleResposne], error)
	FindByActionName(ctx context.Context, actionName string) (*imodel.BPSActionRole, error)
	FindApproverByActionName(ctx context.Context, actionName, role_code string) (imodel.BPSActionApproveIndex, error)
	FindByActionCodeOne(ctx context.Context, actionCode string) (*imodel.BPSActionRole, error)
}

type CPSActionApproveIndexRepository interface {
	SaveIndices(ctx context.Context, indices []imodel.CPSActionApproveIndex) error
	SyncIndices(ctx context.Context, oldActionName string, portalCardName string, newIndices []imodel.CPSActionApproveIndex, isVersionChanged bool) error
	ExistsByRoleAndAction(ctx context.Context, roleID string, actionName string) (bool, error)
	FindMakerAllocationsByRoleID(ctx context.Context, roleID bson.ObjectID) ([]imodel.CPSActionApproveIndex, error)
	FindCheckerAllocationsByRoleID(ctx context.Context, roleID bson.ObjectID) ([]imodel.CPSActionApproveIndex, error)
	FindAuditorAllocationsByRoleID(ctx context.Context, roleID bson.ObjectID) ([]imodel.CPSActionApproveIndex, error)
	PopulateUserApproverAllocations(ctx context.Context, role_id string) ([]string, []string, []string, []string, error)
	FindByRoleAndAction(ctx context.Context, roleID string, actionName string, version int64) (*imodel.CPSActionApproveIndex, error)
	DeleteMany(ctx context.Context, makerIndex []bson.ObjectID, checkerIndex [][]bson.ObjectID, auditorIndex []bson.ObjectID, roleCode string) error
	InsertMany(ctx context.Context, makerIndex []bson.ObjectID, checkerIndex [][]bson.ObjectID, auditorIndex []bson.ObjectID, roleCode string) error
	DeleteAll(ctx context.Context, prev imodel.CPSActionRoleResposne) error
	InsertAll(ctx context.Context, new imodel.CPSActionRole) error
}

type BPSActionApproveIndexRepository interface {
	SaveIndices(ctx context.Context, indices []imodel.BPSActionApproveIndex) error
	SyncIndices(ctx context.Context, oldActionName string, newIndices []imodel.BPSActionApproveIndex) error
	ExistsByRoleAndAction(ctx context.Context, roleID string, actionName string) (bool, error)
	FindMakerAllocationsByRoleID(ctx context.Context, roleID bson.ObjectID) ([]imodel.BPSActionApproveIndex, error)
	FindCheckerAllocationsByRoleID(ctx context.Context, roleID bson.ObjectID) ([]imodel.BPSActionApproveIndex, error)
	FindAuditorAllocationsByRoleID(ctx context.Context, roleID bson.ObjectID) ([]imodel.BPSActionApproveIndex, error)
	PopulateUserApproverAllocations(ctx context.Context, role_id string) ([]string, []string, []string, error)
	FindByRoleAndAction(ctx context.Context, roleID string, actionName string) (*imodel.BPSActionApproveIndex, error)
	DeleteMany(ctx context.Context, makerIndex []bson.ObjectID, checkerIndex [][]bson.ObjectID, auditorIndex []bson.ObjectID, roleCode string) error
	InsertMany(ctx context.Context, makerIndex []bson.ObjectID, checkerIndex [][]bson.ObjectID, auditorIndex []bson.ObjectID, roleCode string) error
	DeleteAll(ctx context.Context, prev imodel.BPSActionApproveIndex) error
	InsertAll(ctx context.Context, new imodel.BPSActionApproveIndex) error
}

type SitotaRepository interface {
	FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.SitotaTransaction], error)
	Get(ctx context.Context, id string) (*model.SitotaTransaction, error)
}

type MiniAppCategoryRepository interface {
	Create(ctx context.Context, category *model.MiniAppCategory) error
	Update(ctx context.Context, category *model.MiniAppCategory, id string) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
}

type MiniAppProductCodeRepository interface {
	Create(ctx context.Context, productCode *model.MiniAppProductCode) error
	Update(ctx context.Context, productCode *model.MiniAppProductCode, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
}

type TransactionRepository interface {
	FindTransactionByID(ctx context.Context, id string) (transaction_dto.FullTransaction, error)
	FindTransactionByCifOrAccountNumberOrFT(ctx context.Context, identifier string) (transaction_dto.FullTransaction, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]transaction_dto.FullTransaction], error)
}

type EventMerchantRepository interface {
	FindOne(ctx context.Context, filter bson.M) (*model.EventMerchant, error)
	Update(ctx context.Context, id string, eventMerchant model.EventMerchant) error
	Create(ctx context.Context, eventMerchant model.EventMerchant) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.EventMerchant, error)
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]model.EventMerchant], error)
}
type LogisticsMerchantRepository interface {
	FindOne(ctx context.Context, filter bson.M) (*local_model.LogisticsMerchant, error)
	Update(ctx context.Context, id string, logisticsMerchant local_model.LogisticsMerchant) error
	Create(ctx context.Context, logisticsMerchant local_model.LogisticsMerchant) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*local_model.LogisticsMerchant, error)
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]local_model.LogisticsMerchant], error)
}

type UssdMerchantRepository interface {
	Create(ctx context.Context, data imodel.UssdMerchant) error
	Update(ctx context.Context, id string, update bson.M) error
	FindById(ctx context.Context, id string) (ussd_merchant_dto.UssdMerchantResponse, error)
	FindByOr(ctx context.Context, phone, email, account_number string) (imodel.UssdMerchant, error)
	Find(ctx context.Context, filter bson.M) (ussd_merchant_dto.UssdMerchantResponse, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (types.PaginatedResponse[[]ussd_merchant_dto.UssdMerchantResponse], error)
}

type AccessListSegmentationRepository interface {
	CreateAccountSegment(ctx context.Context, accessListSegmentation access_list_segmentation_dto.CreateAccessListSegmentationRequest) error
	CreateBlockSegment(ctx context.Context, accessListSegmentation access_list_segmentation_dto.CreateAccessListSegmentationRequest) error
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]local_model.AccessListSegmentation], error)
	FindByID(ctx context.Context, id string) (*local_model.AccessListSegmentation, error)
	FindBySegmentationAndServiceID(ctx context.Context, segmentationID, serviceID string) (*local_model.AccessListSegmentation, error)
	FindByAccountSegmentationAndAccessListKeys(ctx context.Context, customerSegments string, segmentKeys []string) (*local_model.AccessListSegmentation, error)
	FindByIDS(ctx context.Context, ids []string, t string) (*local_model.AccessListSegmentation, error)
	FindByIDAndType(ctx context.Context, ids string, t string) (*local_model.AccessListSegmentation, error)
	FindBySegmentIDAndAccessListKeys(ctx context.Context, id string, keys []string) (*local_model.AccessListSegmentation, error)
	Update(ctx context.Context, id string, accessListSegmentation local_model.AccessListSegmentation) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindAllBySegmentIDorSegmentCode(ctx context.Context, segmentIDorCode string) ([]local_model.AccessListSegmentation, error)
	FindAllBySegmentIDorSegmentCodeAndKeys(ctx context.Context, segmentIDorCode string, keys []string) ([]local_model.AccessListSegmentation, error)
	BulkDisable(ctx context.Context, req access_list_segmentation_dto.BulkDisableAccessListSegmentationRequest) error

	FindParentChildRelationship(ctx context.Context) ([]local_model.AccessItemRelation, error)
}

type MiniAppMerchant interface {
	Create(ctx context.Context, merchant *mini_model.MiniAppMerchant) (*mini_model.MiniAppMerchant, error)
	Update(ctx context.Context, id string, merchant *mini_model.MiniAppMerchant) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*mini_model.MiniAppMerchant, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]mini_model.MiniAppMerchant], error)
	FindOne(ctx context.Context, filter bson.M) (*mini_model.MiniAppMerchant, error)
}

type CustomerSegmentationRepository interface {
	Create(ctx context.Context, seg *imodel.CustomerSegmentation) error
	Update(ctx context.Context, id string, seg *imodel.CustomerSegmentation) error
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*imodel.CustomerSegmentation, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.CustomerSegmentation], error)
	FindByCustomerSegmentation(ctx context.Context, customerSegment string) (*imodel.CustomerSegmentation, error)
}

type CPSRolesRepository interface {
	Create(ctx context.Context, req imodel.CPSRoles) error
	Update(ctx context.Context, id string, req imodel.CPSRoles) error
	FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]imodel.CPSRoles], error)
	FindById(ctx context.Context, id string) (*imodel.CPSRoles, error)
	FindByNameOrRoleCode(ctx context.Context, name, roleCode string) (*imodel.CPSRoles, error)
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	FindByCustomerSegmentation(ctx context.Context, customerSegment string) (*imodel.CPSRoles, error)
	EnableServiceAccess(ctx context.Context, roleID string, accessListKeys []string) error
	DisableServiceAccess(ctx context.Context, roleID string, accessListKeys []string) error
	Delete(ctx context.Context, id string) error
}

type CustomerKYCRepository interface {
	Create(ctx context.Context, req *imodel.CustomerKYC) error
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.CustomerKYC], error)
	FindByID(ctx context.Context, id string) (*imodel.CustomerKYC, error)
	UpdateKYCStatus(ctx context.Context, id string, status string) error
	Delete(ctx context.Context, id string) error
}
