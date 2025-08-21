package service

import (
	"context"

	"cbe-super-app-cps-action/internal/constants/dto/feedback"
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
}

type BulkService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type CPSUserService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type CustomerService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type DepartmentService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type DonationService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type EventService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type FaydaAccountService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type FeedbackService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
	GetFeedbacks(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.Feedback], error)
	GetFeedbackByID(ctx context.Context, id string) (*model.Feedback, error)
	CreateFeedback(ctx context.Context, req feedback.FeedbackRequest, userID string) (*model.Feedback, error)
}

type HQService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type MiniAppService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type MiniAppMerchantService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type NotificationService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type PasswordRuleService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type PermissionService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type PortalCardService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
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
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type AccountBlockService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type AccountValidationService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type ActionService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type AdvertService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type AmountBasedAuthService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type AvatarService interface {
	CreateAvatar(ctx context.Context, avatar *model.Avatar) error
	UpdateAvatar(ctx context.Context, id string, avatar *model.Avatar) error
	DeleteAvatar(ctx context.Context, id string) error
	FetchAllAvatar(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.Avatar], error)
	FetchAvatarById(ctx context.Context, id string) (*model.Avatar, error)
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type BankService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type BudgetCategoryService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
}

type BPSUserService interface {
	Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error)
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
