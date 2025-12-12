package service

import (
	"cbe-super-app-cps-action/internal/constants"
	actionrole_dto "cbe-super-app-cps-action/internal/constants/dto/action_role"
	"cbe-super-app-cps-action/internal/constants/dto/bankvault"
	budget_category "cbe-super-app-cps-action/internal/constants/dto/budget_category"
	cpsuser "cbe-super-app-cps-action/internal/constants/dto/cps_user"
	deviceversion "cbe-super-app-cps-action/internal/constants/dto/device_version"
	"cbe-super-app-cps-action/internal/constants/dto/merchant_lookup"
	transaction_dto "cbe-super-app-cps-action/internal/constants/dto/transaction"

	fbdto "cbe-super-app-cps-action/internal/constants/dto/feedback"

	amountauthdto "cbe-super-app-cps-action/internal/constants/dto/amount_based_auth"
	actionDto "cbe-super-app-cps-action/internal/constants/dto/cps_action"
	"cbe-super-app-cps-action/internal/constants/dto/customer"
	topupDto "cbe-super-app-cps-action/internal/constants/dto/topup"
	vaultgroup "cbe-super-app-cps-action/internal/constants/dto/vaultgroup_category"

	notify "cbe-super-app-cps-action/internal/constants/dto/notification"

	bank_dto "cbe-super-app-cps-action/internal/constants/dto/bank"
	department_dto "cbe-super-app-cps-action/internal/constants/dto/department"
	donation_dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	donationCat_dto "cbe-super-app-cps-action/internal/constants/dto/donation_category"
	donationComp_dto "cbe-super-app-cps-action/internal/constants/dto/donation_company"
	eventdto "cbe-super-app-cps-action/internal/constants/dto/event"
	hqDto "cbe-super-app-cps-action/internal/constants/dto/hq"
	kyc_dto "cbe-super-app-cps-action/internal/constants/dto/kyc_verifier"
	permission_dto "cbe-super-app-cps-action/internal/constants/dto/permission"
	productcode_dto "cbe-super-app-cps-action/internal/constants/dto/productcode"
	vaultCategory_dto "cbe-super-app-cps-action/internal/constants/dto/vaultgroup_category"
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	"cbe-super-app-cps-action/internal/constants/types"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/member"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	dtoEncryption "cbe-super-app-cps-action/internal/constants/dto/encryption"
	dtoService "cbe-super-app-cps-action/internal/constants/dto/service_details"

	"context"
	"mime/multipart"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type CPSActionService interface {
	ApproveCPSAction(ctx context.Context, action *model.CPSAction) error
	CreateCPSAction(ctx context.Context, action *model.CPSAction) error
	RejectCPSAction(ctx context.Context, action_code string, action *model.CPSAction) error
	GetCPSActionsByDepartment(ctx context.Context, department string, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.CPSAction], error)
	GetActionCountsByDepartemnt(ctx context.Context, department string) (*actionDto.CPSActionCountResponse, error)
	GetCPSActionByID(ctx context.Context, id, department string) (*model.CPSAction, error)
	GetCPSActionByUniqueID(ctx context.Context, id, department string) (*model.CPSAction, error)
	GetCPSActionByActionCode(ctx context.Context, uniqueID, department string) (*model.CPSAction, error)
}

type BranchService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type BudgetCategoryService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateBudgetCategory(ctx context.Context, budgetCategory budget_category.CreateBudgetRequest) error
	FetchBudgetCategory(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]budget_category.BudgetCategoryResponse], error)
	FetchBudgetCategoryByID(ctx context.Context, id string) (*budget_category.BudgetCategoryResponse, error)
	UpdateBudgetCategory(ctx context.Context, id string, budgetCategory budget_category.UpdateBudgetRequest) error
	DeleteBudgetCategory(ctx context.Context, id string) error
	EnableOrDisableBudgetCategory(ctx context.Context, id string, enable bool) error
}

type BulkService interface {
	GetAllBulkServices(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.APPAccessList], error)
	EnableBulkService(ctx context.Context, keys []string) error
	DisableBulkService(ctx context.Context, keys []string) error
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type CPSUserService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateUserRequest(ctx context.Context, req cpsuser.CreateUserRequest) error
	UpdateUserRequest(ctx context.Context, req cpsuser.UpdateUserRequest) error
	FetchUserByUserCode(ctx context.Context, userCode string) (*cpsuser.CPSUserDTO, error)
	GetAllCPSUsers(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*cpsuser.CPSUserWithDepartment], error)
	DeleteUserRequest(ctx context.Context, userCode string) error
	DisableUser(ctx context.Context, userCode string) error
	EnableUser(ctx context.Context, userCode string) error
	GetPopulatedCpsUser(ctx context.Context, userCode string) (*cpsuser.CpsUserResponse, error)
	GetCpsUserDetail(ctx context.Context, userCode string) (*cpsuser.CpsUserDetail, error)
}

type NotificationService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateNotification(ctx context.Context, notification notify.NotificationRequest) (*notify.NotificationResponse, error)
	UpdateNotification(ctx context.Context, id string, notification notify.NotificationRequest) (*notify.NotificationResponse, error)
	DeleteNotification(ctx context.Context, id string) error
	EnableNotification(ctx context.Context, id string) error
	DisableNotification(ctx context.Context, id string) error
	FetchNotificationByID(ctx context.Context, id string) (*notify.NotificationResponse, error)
	FetchNotifications(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*notify.NotificationResponse], error)
}

type CustomerService interface {
	GetCustomersDetail(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*member.User], error)
	GetBlockedCustomer(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*member.User], error)
	GetCustomerByID(ctx context.Context, id string) (*member.User, error)
	EnableCustomerByID(ctx context.Context, id string, user_otp string) error
	DisableCustomerByID(ctx context.Context, id string, payload customer.CustomerDisableDTO) error
	ApproveFaydaCustomer(ctx context.Context, id string, req customer.FaydaApproveRequest) error
	GetLinkedAccount(ctx context.Context, customerNumber string) ([]*model.LinkedAccount, error)
	CreateEnableCustomerSession(ctx context.Context, id string) (string, error)
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type DepartmentService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateDepartment(ctx context.Context, department department_dto.CreateDepartmentRequest) error
	UpdateDepartment(ctx context.Context, id string, department department_dto.UpdateDepartmentRequest) error
	EnableDisableDepartment(ctx context.Context, id string, enableDisable bool) error
	GetAllDepartments(ctx context.Context, filterParams *types.Filter) (types.PaginatedResponse[[]model.Department], error)
	GetDepartmentByID(ctx context.Context, id string) (*model.Department, error)
}

type DonationService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateDonation(ctx context.Context, donation donation_dto.DonationRequest) error
	UpdateDonation(ctx context.Context, id string, donation donation_dto.DonationRequest) error
	FetchDonation(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]donation_dto.DonationListResponse], error)
	FetchDonationByID(ctx context.Context, id string) (*donation_dto.DonationListResponse, error)
	UpdateDonationImage(ctx context.Context, id string, image donation_dto.DonationImageUpdateRequest) error
	DeleteDonationImage(ctx context.Context, id string, imageID string) error
	AddDonationImage(ctx context.Context, id string, image donation_dto.DonationRequest) error
	EnableDonation(ctx context.Context, id string) error
	DisableDonation(ctx context.Context, id string) error
}

type DonationCategoryService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateDonationCategory(ctx context.Context, donation donationCat_dto.DonationCategoryRequest) error
	FetchDonationCategory(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]donationCat_dto.DonationCategoryListResponse], error)
	FetchDonationCategoryByID(ctx context.Context, id string) (*donationCat_dto.DonationCategoryListResponse, error)
	UpdateDonationCategory(ctx context.Context, id string, donation donationCat_dto.DonationCategoryRequest) (donationCat_dto.DonationCategoryRequest, error)
	EnableDonationCategory(ctx context.Context, id string) error
	DisableDonationCategory(ctx context.Context, id string) error
}

type DonationCompanyService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateDonationCompany(ctx context.Context, donationCompany donationComp_dto.DonationCompanyRequest) error
	FetchDonationCompany(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]donationComp_dto.DonationCompanyListResponse], error)
	FetchDonationCompanyByID(ctx context.Context, id string) (*donationComp_dto.DonationCompanyListResponse, error)
	UpdateDonationCompany(ctx context.Context, id string, donationCompany donationComp_dto.DonationCompanyRequest) (*model.DonationCompany, error)
	AccountLookup(ctx context.Context, accountNumber string) (*model.AccountDetail, error)
	EnableDonationCompany(ctx context.Context, id string) error
	DisableDonationCompany(ctx context.Context, id string) error
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

type KYCVerifierService interface {
	FetchKYCList(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*kyc_dto.KYCVerifierResponse], error)
	FetchKYCByID(ctx context.Context, id string) (*kyc_dto.KYCVerifierResponse, error)
	UpdateKYC(ctx context.Context, id string, req kyc_dto.UpdateKYCRequest) error
	ApproveKYC(ctx context.Context, id string, req kyc_dto.ApproveKYCRequest) error
	Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
}

type FaydaAccountService interface {
	EnableOrDisableFayda(ctx context.Context, user_code string, isEnable bool) error
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type FeedbackService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateFeedback(ctx context.Context, req fbdto.FeedbackRequest, userID string) (*model.Feedback, error)
	GetFeedbackByID(ctx context.Context, id string) (*fbdto.FeedbackResponse, error)
	GetFeedbacks(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponseForFeedback[[]*fbdto.FeedbackResponse], error)
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
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type MiniAppMerchantService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	Create(ctx context.Context, req *model.MiniAppMerchant) (*model.MiniAppMerchant, error)
	Update(ctx context.Context, id string, data *model.MiniAppMerchant) (*model.MiniAppMerchant, *model.MiniAppMerchant, error)
	FindAllWithPagination(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.MiniAppMerchant], error)
	FindByID(ctx context.Context, id string) (*model.MiniAppMerchant, error)
	Delete(ctx context.Context, id string) error
	EnableOrDisable(ctx context.Context, id string, enable bool) error
	MerchantLookup(ctx context.Context, merchantID string) (*merchant_lookup.MerchantLookUpResponse, error)
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
	GetPermissionGroupById(ctx context.Context, id string) (*model.PermissionGroup, error)
	GetPermissionGroups(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.PermissionGroup], error)
	GetAllPermissionCategoriesWithPermissions(ctx context.Context) ([]*model.PermissionCategory, error)
	GetPermissionCategoriesByDepartment(ctx context.Context, departmentId string) (map[string][]*model.PermissionCategory, error)
	GetPermissionGroupsByDepartment(ctx context.Context, departmentId string, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.PermissionGroup], error)
	ValidatePermissionCategories(ctx context.Context, categoryIDs []string) (bool, error)
	ValidatePermissionGroups(ctx context.Context, groupIDs []string) (bool, error)
	GetPopulatedPermissionCategories(ctx context.Context, categoryIDs []bson.ObjectID) ([]cpsuser.PermissionCategoryResponse, error)
	GetPopulatedPermissionGroups(ctx context.Context, groupIDs []bson.ObjectID) ([]cpsuser.PermissionGroupResponse, error)
}

type PortalCardService interface {
	GetAll(ctx context.Context, filterParam *types.Filter) (*types.PaginatedResponse[[]*model.Card], error)
	ValidatePortalCard(ctx context.Context, names []string) (bool, error)
}

type ProductCodeService interface {
	FetchProductCodeByID(ctx context.Context, id string) (*model.ProductCode, error)
	FetchAllProductCodes(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ProductCode], error)
	UpdateProductCode(ctx context.Context, request productcode_dto.UpdateProductCodeRequest) (*model.ProductCode, *model.ProductCode, error)
	Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
}

type DeviceVersionServiceSrv interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateDeviceVersion(ctx context.Context, deviceVersion deviceversion.CreateDeviceVersionRequest) error
	EnableDisableDeviceVersion(ctx context.Context, id string, enableDisable bool) error
	GetAllDeviceVersions(ctx context.Context, filterParams *types.Filter) (types.PaginatedResponse[[]model.DeviceVersionControl], error)
	GetDeviceVersionByID(ctx context.Context, id string) (model.DeviceVersionControl, error)
	UpdateDeviceVersion(ctx context.Context, id string, req deviceversion.UpdateDeviceVersionRequest) error
}

type ServiceService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	GetAllService(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ServiceDetails], error)
	GetAllMinimumTransferCap(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*dtoService.MinimumTransferCapResponse], error)
	GetAllMaximumTransferCap(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*dtoService.MaximumTransferCapResponse], error)
	GetAllServiceFee(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*dtoService.ServiceFeeResponse], error)
	GetAllTotalTransferCap(ctx context.Context) (*dtoService.TotalTransferCapResponse, error)
	GetServiceFeeDetail(ctx context.Context, id string) (*dtoService.ServiceFeeDetailResponse, error)
	UpdateServiceFee(ctx context.Context, id string, req dtoService.ServiceFeeDetailDTO) error
	UpdateSingleMaxTransfer(ctx context.Context, id string, req dtoService.SingleMaxTransferRequest) error
	UpdateTotalMaxTransferCap(ctx context.Context, newTotalCap dtoService.TotalMaxTransferUpdateRequest) error
	UpdateMinimumTransferCap(ctx context.Context, id string, req dtoService.MinimumTransferUpdateRequest) error
	DeleteServiceFeeTire(ctx context.Context, id string) error
}

type UnlinkService interface {
	GetUserByAccount(ctx context.Context, accNumber string) (*member.User, error)
	GetAllArchivedUser(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.ArchivedUser], error)
	UnlinkUserCif(ctx context.Context, userCode string) error
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type WalletService interface {
	CreateWallet(ctx context.Context, req walletDto.WalletRequest) error
	UpdateWallet(ctx context.Context, id string, req walletDto.WalletRequest, fieldsProvided map[string]bool) error
	DeleteWallet(ctx context.Context, id string) error
	EnableOrDisableWallet(ctx context.Context, id string, enable bool) error
	GetWallet(ctx context.Context, id string) (*model.Wallet, error)
	GetAllWallet(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.Wallet], error)
	Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
}

type TopupService interface {
	CreateTopup(ctx context.Context, req topupDto.TopupRequest) error
	UpdateTopup(ctx context.Context, id string, req topupDto.TopupRequest) error
	DeleteTopup(ctx context.Context, id string) error
	EnableOrDisableTopup(ctx context.Context, id string, enable bool) error
	GetTopup(ctx context.Context, id string) (*model.Topup, error)
	GetAllTopup(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.Topup], error)
	Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error)
}

type AccountBlockService interface {
	GetBranchById(ctx context.Context, id string) (*model.AccountBlock, error)
	GetAllBranches(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error)
	GetRegionById(ctx context.Context, id string) (*model.AccountBlock, error)
	GetAllRegions(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error)
	GetDistrictById(ctx context.Context, id string) (*model.AccountBlock, error)
	GetAllDistricts(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error)
	GetCityById(ctx context.Context, Id string) (*model.AccountBlock, error)
	GetAllCities(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*model.AccountBlock], error)
	EnableOrDisableBranches(ctx context.Context, branchIds []string, reason string, enabled bool) error
	EnableOrDisableRegions(ctx context.Context, regionIds []string, reason string, enabled bool) error
	EnableOrDisableDistricts(ctx context.Context, regionIds []string, reason string, enabled bool) error
	EnableOrDisableCities(ctx context.Context, ids []string, reason string, enabled bool) error
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
	FindAllWithPagination(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.AuthTier], error)
	UpdateAmountBasedAuth(ctx context.Context, id string, method constants.Method, request amountauthdto.UpdateAmountBasedAuthRequest) error
}

type AvatarService interface {
	CreateAvatar(ctx context.Context, avatar *model.Avatar, fileHeader *multipart.FileHeader) error
	UpdateAvatar(ctx context.Context, id string, avatar *model.Avatar, fileHeader *multipart.FileHeader, fromEnabledDisable bool) error
	DeleteAvatar(ctx context.Context, id string) error
	EnableDisable(ctx context.Context, id string, enable bool) error
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

type BPSUserService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
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

type ServiceLayer struct {
	EventService           EventService
	BulkService            BulkService
	CustomerService        CustomerService
	CPSAction              CPSActionService
	Feedback               FeedbackService
	Unlink                 UnlinkService
	BpsUser                BPSUserService
	BudgetCategory         BudgetCategoryService
	Bank                   BankService
	PortalCard             PortalCardService
	Advert                 AdvertService
	ValidationService      AccountValidationService
	Wallet                 WalletService
	Topup                  TopupService
	MiniAppMerchant        MiniAppMerchantService
	AccountBlock           AccountBlockService
	Department             DepartmentService
	PasswordRule           PasswordRuleService
	HQService              HQService
	MiniAppService         MiniAppService
	Fayda                  FaydaAccountService
	Avatar                 AvatarService
	AmountBasedAuth        AmountBasedAuthService
	Permission             PermissionService
	CPSUser                CPSUserService
	ServiceDetails         ServiceService
	AccountValidation      AccountValidationService
	ProductCode            ProductCodeService
	BankVault              BankVaultService
	VaultGroupCategory     VaultGroupCategoryService
	DonationCategory       DonationCategoryService
	DonationCompany        DonationCompanyService
	Donation               DonationService
	NotificationService    NotificationService
	ArticleService         ArticleService
	ArticleCategoryService ArticleCategoryService
	ShortVideoService      ShortVideoService
	NewsTagService         NewsTagService
	NewsCategoryService    NewsCategoryService
	Sitota                 SitotaService
	KYCVerifier            KYCVerifierService
	NewsTagsService        NewsTagsService
	DeviceVersion          DeviceVersionServiceSrv
	Encryption             EncryptionService
	BPSActionRole          BPSActionRoleService
	TransactionService     TransactionService
	MiniAppCategory        MiniAppCategoryService
	CPSActionRole          CPSActionRoleService
}

type ServiceContainer struct {
	AccountBlockContainer      AccountBlockService
	AccountContainer           AccountValidationService
	ActionContainer            ActionService
	AdContainer                AdvertService
	AmountBasedAuthContainer   AmountBasedAuthService
	AvatarDomian               AvatarService
	BankContainer              BankService
	BPSUserContainer           BPSUserService
	BudgetCategoryContainer    BudgetCategoryService
	CPSActionContainer         CPSActionService
	CPSUserContainer           CPSUserService
	CustomerContainer          CustomerService
	DepartmentContainer        DepartmentService
	EventContainer             EventService // fully not ready
	FaydaContainer             FaydaAccountService
	FeedbackContainer          FeedbackService
	HQContainer                HQService
	MiniAppContainer           MiniAppService
	PasswordRuleContainer      PasswordRuleService
	PermissionContainer        PermissionService
	PortalCardContainer        PortalCardService
	UnlinkContainer            UnlinkService
	WalletContainer            WalletService
	TopupContainer             TopupService
	MiniAppMerchantContainer   MiniAppMerchantService
	AccountLookup              AccountSearchService
	BulkServiceContainer       BulkService
	ServiceCheckContainer      ServiceService
	DeviceVersionContainer     DeviceVersionServiceSrv
	KeyGenService              KeyGeneratorService
	NotificationService        NotificationService
	ProductCodeService         ProductCodeService
	DonationContainer          DonationService
	Unlink                     UnlinkService
	DonationCategoryContainer  DonationCategoryService
	DonationCompanyContainer   DonationCompanyService
	ArticleContainer           ArticleService
	ArticleCategoryContainer   ArticleCategoryService
	ShortVideoServiceContainer ShortVideoService
	NewsTagContainer           NewsTagService
	NewsCategoryContainer      NewsCategoryService
	SitotaContainer            SitotaService
	KYCVerifierContainer       KYCVerifierService
	NewsTagsServiceContainer   NewsTagsService
	EncryptionContainer        EncryptionService
	BankProductContainer       BankVaultService
	VaultCategoryContainer     VaultGroupCategoryService
	BPSActionRoleContainer     BPSActionRoleService
	TransactionContainer       TransactionService
	MiniAppCategoryContainer   MiniAppCategoryService
	CPSActionRoleContainer     CPSActionRoleService
}

type BPSActionRoleService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.ActionRole], error)
	GetByActionCode(ctx context.Context, actionCode string) (*actionrole_dto.GetActionRoleByActionCodeRes, error)
	Create(ctx context.Context, req actionrole_dto.CreateActionRoleRequest) error
	Update(ctx context.Context, actionCode string, req actionrole_dto.UpdateActionRoleRequest) error
	Enable(ctx context.Context, actionCode string) error
	Disable(ctx context.Context, actionCode string) error
}

type BankVaultService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateBankVault(ctx context.Context, req *model.BankVaultProduct) (string, error)
	FindAllBankVaults(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*bankvault.BankVaultProductResponse], error)
	GetBankVault(ctx context.Context, id string) (*bankvault.BankVaultProductResponse, error)
	UpdateBankVault(ctx context.Context, id string, req *model.UpdateBankVault) (string, error)
	DeleteBankVault(ctx context.Context, id string) (string, error)
	EnableBankVault(ctx context.Context, id string) error
	DisableBankVault(ctx context.Context, id string) error
	FindAllBankLockedVaultsWithPagination(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.LockedVault], error)
	// GetTransactions(ctx context.Context, id string) (*bankvault.LockedVaultResponse, error)
	FindAllGroupVaultsWithPagination(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.GroupVault], error)
	// GetGroupVault(ctx context.Context, id string) (*bankvault.GroupVaultResponse, error)
}
type VaultGroupCategoryService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	CreateVaultGroupCategory(ctx context.Context, req *vaultCategory_dto.CreateVaultGroupCategoryRequest) (string, error)
	FindAllVaultGroupCategories(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*vaultgroup.VaultGroupCategoryResponse], error)
	GetVaultGroupCategory(ctx context.Context, id string) (*vaultgroup.VaultGroupCategoryResponse, error)
	UpdateVaultGroupCategory(ctx context.Context, id string, req *vaultCategory_dto.UpdateVaultGroupCategoryRequest) (string, error)
	DeleteVaultGroupCategory(ctx context.Context, id string) (string, error)
	EnableVaultGroupCategory(ctx context.Context, id string) error
	DisableVaultGroupCategory(ctx context.Context, id string) error
}
type ArticleService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type ArticleCategoryService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}
type ShortVideoService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type NewsTagsService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type NewsTagService interface {
	FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.NewsTag], error)
	GetNewsTagByID(ctx context.Context, id string) (*model.NewsTag, error)
	CreateNewsTags(ctx context.Context, tagName []string) error
	UpdateNewsTag(ctx context.Context, id string, tagName string) error
	DeleteNewsTag(ctx context.Context, id string) error
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}
type NewsCategoryService interface {
	FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.NewsCategory], error)
	GetNewsCategoryByID(ctx context.Context, id string) (*model.NewsCategory, error)
	CreateNewsCategory(ctx context.Context, categoryName []string) error
	UpdateNewsCategory(ctx context.Context, id string, categoryName string) error
	DeleteNewsCategory(ctx context.Context, id string) error
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type SitotaService interface {
	GetAllSitotas(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.SitotaTransaction], error)
	GetSitotaByID(ctx context.Context, id string) (*model.SitotaTransaction, error)
}

type EncryptionService interface {
	LocalEncryptPassword(req dtoEncryption.EncryptionRequest, dataType, userSalt, action string) (dtoEncryption.EncryptionResponse, string, error)
}

type MiniAppCategoryService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type TransactionService interface {
	FetchTransactionByID(ctx context.Context, id string) (transaction_dto.FullTransaction, error)
	FetchAllTransactions(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]transaction_dto.FullTransaction], error)
}
type CPSActionRoleService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	FindAllWithPagination(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]*model.CPSActionRole], error)
	GetByActionCode(ctx context.Context, actionCode string) (*actionrole_dto.GetActionRoleByActionCodeRes, error)
	Create(ctx context.Context, req actionrole_dto.CreateActionRoleRequest) error
	Update(ctx context.Context, actionCode string, req actionrole_dto.UpdateActionRoleRequest) error
	Enable(ctx context.Context, actionCode string) error
	Disable(ctx context.Context, actionCode string) error
}
