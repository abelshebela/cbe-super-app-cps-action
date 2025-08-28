package service

import (
	"context"
	"mime/multipart"
    cpsuser "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	permission_dto "cbe-super-app-cps-action/internal/constants/dto/permission"
	bank_dto "cbe-super-app-cps-action/internal/constants/dto/bank"
	department_dto "cbe-super-app-cps-action/internal/constants/dto/department"
	eventdto "cbe-super-app-cps-action/internal/constants/dto/event"
	fbdto "cbe-super-app-cps-action/internal/constants/dto/feedback"
	hqDto "cbe-super-app-cps-action/internal/constants/dto/hq"
	miniappdto "cbe-super-app-cps-action/internal/constants/dto/mini_app"
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
)

type CPSActionService interface {
	ApproveCPSAction(ctx context.Context, action *model.CPSAction) error
	CreateCPSAction(ctx context.Context, action *model.CPSAction) error
	RejectCPSAction(ctx context.Context, action_code string, action *model.CPSAction) error
	GetCPSActionsByDepartment(ctx context.Context, department string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error)
	GetCPSActionByID(ctx context.Context, id, department string) (*model.CPSAction, error)
	GetCPSActionByUniqueID(ctx context.Context, id, department string) (*model.CPSAction, error)
	GetCPSActionByActionCode(ctx context.Context, uniqueID, department string) (*model.CPSAction, error)
}

type BranchService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type BudgetService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateBudgetIcon(ctx context.Context, fileHeader *multipart.FileHeader, file *multipart.File) error
	BudgetFetchIcons(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Icon], error)
	BudgetUpdateIcon(ctx context.Context, id string, fileHeader *multipart.FileHeader, file *multipart.File) error
	BudgetCreateColor(ctx context.Context, color *model.Color) error
	BudgetFetchColors(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Color], error)
	BudgetUpdateColor(ctx context.Context, id string, color *model.Color) error
	BudgetCheckerApproval(ctx context.Context, actionCode string) error
}

type BulkService interface {
	GetAllBulkServices(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.APPAccessList], error)
	EnableBulkService(ctx context.Context, keys []string) (string, error)
	DisableBulkService(ctx context.Context, keys []string) (string, error)
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type CPSUserService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateUserRequest(ctx context.Context, req cpsuser.CreateUserRequest) error
	UpdateUserRequest(ctx context.Context, req cpsuser.UpdateUserRequest) error
	FetchUserByUserCode(ctx context.Context, userCode string) (*cpsuser.CPSUserDTO, error)
	GetAllCPSUsers(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*cpsuser.CPSUserDTO], error)
	DeleteUserRequest(ctx context.Context, userCode string) error
	DisableUser(ctx context.Context, userCode string) error
	EnableUser(ctx context.Context, userCode string) error
}

type CustomerService interface {
	GetCustomersDetail(ctx context.Context, kyc_level int, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.User], error)
	GetBlockedCustomer(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.User], error)
	GetCustomerByID(ctx context.Context, id string) (*model.User, error)
}

type DepartmentService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateDepartment(ctx context.Context, department department_dto.CreateDepartmentRequest) error
	UpdateDepartment(ctx context.Context, id string, department department_dto.UpdateDepartmentRequest) error
	EnableDisableDepartment(ctx context.Context, id string, enableDisable bool) error
	GetAllDepartments(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Department], error)
	GetDepartmentByID(ctx context.Context, id string) (*model.Department, error)
}

type DonationService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type EventService interface {
	CreateEvent(ctx context.Context, event eventdto.EventRequest) error
	UpdateEvent(ctx context.Context, id string, event eventdto.EventRequest) error
	DeleteEvent(ctx context.Context, id string) error
	EnableDisableEvent(ctx context.Context, id string, enable bool) error
	FetchEventByID(ctx context.Context, id string) (*model.Event, error)
	FetchEvent(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Event], error)
	Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
}

type FaydaAccountService interface {
	EnableOrDisableFayda(ctx context.Context, user_code string, isEnable bool) error
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type FeedbackService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateFeedback(ctx context.Context, req fbdto.FeedbackRequest, userID string) (*model.Feedback, error)
	GetFeedbackByID(ctx context.Context, id string) (*model.Feedback, error)
	GetFeedbacks(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Feedback], error)
}

type HQService interface {
	GetHQDetail(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.HQ], error)
	GetHQ(ctx context.Context, id string) (hqDto.HQ, error)
	GetBlockTime(ctx context.Context) (hqDto.BlockTimeResponse, error)
	GetArchiveTime(ctx context.Context) (hqDto.ArchiveTimeResponse, error)
	GetPasswordExpiry(ctx context.Context) (hqDto.PasswordExpiryResponse, error)
	UpdateBlockTime(ctx context.Context, request hqDto.UpdateBlockTimeRequest) error
	UpdateArchiveTime(ctx context.Context, request hqDto.UpdateArchiveTimeRequest) error
	UpdatePasswordExpiry(ctx context.Context, request hqDto.UpdatePasswordExpiryRequest) error
	Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
}

type MiniAppService interface {
	CreateMiniApp(ctx context.Context, req *miniappdto.MiniAppCreateRequest) error
	UpdateMiniApp(ctx context.Context, req *miniappdto.MiniAppCreateRequest) error
	DeleteMiniApp(ctx context.Context, id string) error
	EnableDisableMiniAppByID(ctx context.Context, id string, enable bool) error
	FindByID(ctx context.Context, id string) (*model.MiniApp, error)
	ListMiniApp(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.MiniApp], error)
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type MiniAppMerchantService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	AddMiniApp(ctx context.Context, merchantID string, miniApp model.MiniApps) error
	UpdateMiniAppEnabledState(ctx context.Context, merchantID string, miniAppID string, enabled bool) error
	SoftDeleteMiniApp(ctx context.Context, merchantID string, miniAppID string) error
	Create(ctx context.Context, req *model.MiniAppMerchant) (*model.MiniAppMerchant, error)
	Update(ctx context.Context, id string, data *model.MiniAppMerchant) (*model.MiniAppMerchant, *model.MiniAppMerchant, error)
	FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.MiniAppMerchant], error)
	FindByID(ctx context.Context, id string) (*model.MiniAppMerchant, error)
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
}

type NotificationService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type PasswordRuleService interface {
	GetAllPasswordRules(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.PasswordRule], error)
	RequestPasswordRuleUpdate(ctx context.Context, id string, body model.PasswordRule) error
	CheckPasswordRule(ctx context.Context, password string) (bool, string)
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type PermissionService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreatePermissionGroup(ctx context.Context, req permission_dto.CreatePermissionGroupRequest) error
	UpdatePermissionGroup(ctx context.Context, req permission_dto.UpdatePermissionGroupRequest) error
	GetPermissionGroup(groupName string) (*model.PermissionGroup, error)
	GetPermissionGroups(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.PermissionGroup], error)
	GetAllPermissionCategoriesWithPermissions(ctx context.Context) ([]*model.PermissionCategory, error)

	// Validation methods
	ValidatePermissionCategories(ctx context.Context, categoryIDs []string) (bool, error)
	ValidatePermissionGroups(ctx context.Context, groupIDs []string) (bool, error)
}

type PortalCardService interface {
	GetAll(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.Card], error)
	ValidatePortalCard(ctx context.Context, names []string) (bool, error)
}

type ProductCodeService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type ServiceService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type UnlinkService interface {
	GetUserByAccount(ctx context.Context, accNumber string) (*model.User, error)
	GetAllArchivedUser(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ArchivedUser], error)
	UnlinkUserCif(ctx context.Context, userCode string) error
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type WalletService interface {
	CreateWallet(ctx context.Context, req walletDto.WalletRequest) error
	UpdateWallet(ctx context.Context, id string, req walletDto.WalletRequest) error
	DeleteWallet(ctx context.Context, id string) error
	EnableOrDisableWallet(ctx context.Context, id string, enable bool) error
	GetWallet(ctx context.Context, id string) (*model.Wallet, error)
	GetAllWallet(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.Wallet], error)
	Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
}

type AccountBlockService interface {
	GetBranchByCode(ctx context.Context, branchCode string) (*model.Branch, error)
	GetAllBranches(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*model.Branch], error)
	GetRegionByCode(ctx context.Context, regionCode string) (*model.Region, error)
	GetAllRegions(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*model.Region], error)
	GetDistrictByCode(ctx context.Context, districtCode string) (*model.District, error)
	GetAllDistricts(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*model.District], error)
	GetCityByCode(ctx context.Context, cityCode string) (*model.City, error)
	GetAllCities(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*model.City], error)
	EnableOrDisableBranches(ctx context.Context, branchCodes []string, enabled bool) error
	EnableOrDisableRegions(ctx context.Context, regionsCode []string, enabled bool) error
	EnableOrDisableDistricts(ctx context.Context, districtsCode []string, enabled bool) error
	EnableOrDisableCities(ctx context.Context, citiesCode []string, enabled bool) error
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type AccountValidationService interface {
	Update(ctx context.Context, id string, rule *model.ValidationRule) error
	FindById(ctx context.Context, id string) (*model.ValidationRule, error)
	FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.ValidationRule], error)
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type ActionService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type AdvertService interface {
	CreateAdvert(ctx context.Context, ad *model.Advert, bannerImage *multipart.FileHeader) error
	FetchAdverts(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]*model.Advert], error)
	FetchAdvertByID(ctx context.Context, id string) (*model.Advert, error)
	UpdateAdvert(ctx context.Context, id string, ad *model.Advert, bannerImage *multipart.FileHeader) error
	DeleteAdvert(ctx context.Context, id string) error
	EnableDisableAdvert(ctx context.Context, id string, enable bool) error
	Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
}

type AmountBasedAuthService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type AvatarService interface {
	CreateAvatar(ctx context.Context, avatar *model.Avatar, fileHeader *multipart.FileHeader) error
	UpdateAvatar(ctx context.Context, id string, avatar *model.Avatar, fileHeader *multipart.FileHeader) error
	DeleteAvatar(ctx context.Context, id string) error
	FetchAllAvatar(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.Avatar], error)
	FetchAvatarById(ctx context.Context, id string) (*model.Avatar, error)
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type BankService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	GetAllBank(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Bank], error)

	GetOneBank(ctx context.Context, id string) (*model.Bank, error)

	CreateOneBank(ctx context.Context, req bank_dto.CreateBankRequest) error
	UpdateOneBank(ctx context.Context, id string, req bank_dto.UpdateBankRequest) error
	DeleteOneBank(ctx context.Context, id string) error

	EnableOrDisableBank(ctx context.Context, id string, enableDisable bool) error
	UpdateLogo(ctx context.Context, id string, logo bank_dto.UpdateLogo) error
}

type BudgetCategoryService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type BPSUserService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) error
	FetchUserByUserCode(ctx context.Context, userCode string) (*model.BPSUser, error)
	GetAllBPSUsers(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.BPSUser], error)
	UpdateBpsUser(ctx context.Context, userCode string, status bool) error
}

type AccountSearchService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type KeyGeneratorService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type ServiceContainer struct {
	AccountBlockContainer    AccountBlockService
	AccountContainer         AccountValidationService
	ActionContainer          ActionService
	AdContainer              AdvertService
	AmountBasedAuthContainer AmountBasedAuthService
	AvatarDomian             AvatarService
	BankContainer            BankService
	BPSUserContainer         BPSUserService
	BudgetCategoryContainer  BudgetCategoryService
	BudgetContainer          BudgetService
	CPSActionContainer       CPSActionService
	CPSUserContainer         CPSUserService
	CustomerContainer        CustomerService
	DepartmentContainer      DepartmentService
	EventContainer           EventService // fully not ready
	FaydaContainer           FaydaAccountService
	FeedbackContainer        FeedbackService
	HQContainer              HQService
	MiniAppContainer         MiniAppService
	PasswordRuleContainer    PasswordRuleService
	PermissionContainer      PermissionService
	PortalCardContainer      PortalCardService
	UnlinkContainer          UnlinkService
	WalletContainer          WalletService
	MiniAppMerchantContainer MiniAppMerchantService
	AccountLookup            AccountSearchService
	BulkServiceContainer     BulkService
	ServiceCheckContainer    ServiceService
	KeyGenService            KeyGeneratorService
	NotificationService      NotificationService
	ProductCodeService       ProductCodeService
	DonationContainer        DonationService
}
