package bpsuser

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	localization "cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	bps_user_core "cbe-super-app-cps-action/internal/service/bps_user/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	// local_model "cbe-super-app-cps-action/internal/constants/model"

	bps_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/bps"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	bpsUserDto "cbe-super-app-cps-action/internal/constants/dto/bps_user"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type bpsUserService struct {
	cpsService     service.CPSActionService
	repo           storage.BPSUserRepository
	CPSUserRepo    storage.CpsUserRepository
	Job_roles_repo storage.JobRoleRepository
	Branch_blocks  storage.AccountBlockRepository
	RoleRepository storage.RoleRepository
	logger         utils.Logger
	minioClient    *s3.Client
	bucketName     string
	cfg            config.VaultConfig
}

func NewBPSUserService(repo storage.BPSUserRepository, JobRolesRepo storage.JobRoleRepository,
	roleRepository storage.RoleRepository, cpsService service.CPSActionService, cpsUserRepo storage.CpsUserRepository, branch_blocks storage.AccountBlockRepository, minioClient *s3.Client, bucketName string, cfg config.VaultConfig, logger utils.Logger) service.BPSUserService {
	return &bpsUserService{
		cpsService:     cpsService,
		repo:           repo,
		CPSUserRepo:    cpsUserRepo,
		Job_roles_repo: JobRolesRepo,
		RoleRepository: roleRepository,
		Branch_blocks:  branch_blocks,
		logger:         logger,
		minioClient:    minioClient,
		bucketName:     bucketName,
		cfg:            cfg,
	}
}

func bpsUserExportRow(user imodel.ExportBPSUser) []string {
	createdBy := user.CreatedBy
	if createdBy == "" {
		createdBy = "SSO"
	}
	lastModified := ""
	if !user.LastModified.IsZero() {
		lastModified = user.LastModified.Format(time.RFC3339)
	}
	expiryDate := ""
	userType := "Permanent"
	if !user.ExpiryDateForDelegation.IsZero() {
		userType = "Delegation"
		expiryDate = user.ExpiryDateForDelegation.Format(time.RFC3339)
	}
	lastLogin := ""
	if !user.LastLogin.IsZero() {
		lastLogin = user.LastLogin.Format(time.RFC3339)
	}
	return []string{
		user.FirstName,
		user.PhoneNumber,
		user.Email,
		user.Branch,
		user.JobTitle,
		user.Role,
		user.UserName,
		user.CreatedAt.Format(time.RFC3339),
		lastLogin,
		fmt.Sprintf("%t", user.Enabled),
		userType,
		expiryDate,
		user.LastModificationAction,
		lastModified,
		createdBy,
		user.ApprovedBy,
	}
}

func (b *bpsUserService) ExportUsers(ctx context.Context, startDate, endDate time.Time, fileType, userName string) (string, error) {
	log := local_util.LoggerFromCtx(ctx, b.logger)

	if b.minioClient == nil {
		log.Errorf("[BPSUser/ExportUsers] minio client is not configured")
		return "", errors.New(localization.BpsUserDataExportedError.Code)
	}

	fileType = strings.ToLower(strings.TrimSpace(fileType))
	if fileType != string(lib.FileTypeCSV) && fileType != string(lib.FileTypePDF) {
		return "", errors.New(localization.ErrorInvalidRequest.Code)
	}

	data, err := b.repo.FindForExport(ctx, startDate, endDate, userName)
	if err != nil {
		log.Errorf("[BPSUser/ExportUsers] failed to fetch users: %v", err)
		return "", err
	}

	if len(data) == 0 {
		return "", errors.New(localization.BpsUserDataNotFoundInDateRange.Code)
	}

	headers := []string{"Full Name", "Phone Number", "Email", "Branch", "Job Title", "Role", "Username", "Created At", "Last Login", "Enabled", "User Type", "Expiry Date (Delegation)", "Last Modification Action", "Last Modified", "Created By", "Approved By"}

	ext := "csv"
	if fileType == string(lib.FileTypePDF) {
		ext = "pdf"
	}

	objectName := fmt.Sprintf("bps_users_%s_to_%s_%d.%s", startDate.Format("20060102"), endDate.Format("20060102"), time.Now().Unix(), ext)

	if fileType == string(lib.FileTypePDF) {
		rows := make([][]string, 0, len(data))
		for _, item := range data {
			rows = append(rows, bpsUserExportRow(item))
		}
		url, exportErr := lib.ExportPDFAndUpload(ctx, b.minioClient, b.bucketName, b.cfg, objectName, headers, rows, lib.PDFExportOptions{PageSize: "A4"}, nil, b.logger)
		if exportErr != nil {
			log.Errorf("[BPSUser/ExportUsers] pdf export failed: %v", exportErr)
			return "", errors.New(localization.BpsUserDataExportedError.Code)
		}

		return url, nil
	}

	url, exportErr := lib.ExportCSVAndUpload(ctx, b.minioClient, b.bucketName, b.cfg, objectName, headers, func(writer *csv.Writer) error {
		for _, item := range data {
			if err := writer.Write(bpsUserExportRow(item)); err != nil {
				return err
			}
		}
		return nil
	}, b.logger)
	if exportErr != nil {
		log.Errorf("[BPSUser/ExportUsers] csv export failed: %v", exportErr)
		return "", errors.New(localization.BpsUserDataExportedError.Code)
	}

	return url, nil
}

// Authorize implements service.BPSUserService.
func (b *bpsUserService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "BPS User", "Authorize")
	defer span.End()

	var actionMap interface{}
	marshaled, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		span.AddEvent("failed to marshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		fmt.Printf("failed to marshal CurrentAction: %v\n", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	fmt.Printf("JSON bytes: %s\n", string(marshaled))
	err = json.Unmarshal(marshaled, &actionMap)
	if err != nil {
		span.AddEvent("failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", cpsAction.UniqueId),
		))
		fmt.Printf("failed to unmarshal CurrentAction: %v\n", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	actionData := bps_user_core.BPSUser_mapper(actionMap.(map[string]interface{}))
	// if cpsAction.UniqueId != "" {
	// 	actionData.ID = cpsAction.UniqueId
	// }

	actionData.LastModifiedAt = time.Now()
	// local_actionData := bps_user_core.MapBPSUserToWithJobTitle(actionData)

	switch cpsAction.RequestAction {
	case string(constants.RequestCreateBPSUser):
		if err := b.repo.Create(ctx, actionData); err != nil {
			span.AddEvent("[Authorize] failed to create BPS user", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[BpsUserSvc][Authorize] create err: %v", err)
			return nil, err
		}
		b.logger.Infof("[BpsUserSvc][Authorize] created")
		return nil, nil
	case string(constants.RequestBpsUserUpdate):
		if err := b.repo.Update(ctx, &actionData); err != nil {
			span.AddEvent("[Authorize] failed to update BPS user", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[BpsUserSvc][Authorize] update err: %v", err)
			return nil, err
		}
		b.logger.Infof("[BpsUserSvc][Authorize] updated id: %s", cpsAction.UniqueId)
		return nil, nil
	case string(constants.RequestBpsUserDelete):
		actionData.IsDeleted = true
		if err := b.repo.Update(ctx, &actionData); err != nil {
			span.AddEvent("[Authorize] failed to soft-delete BPS user", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[BpsUserSvc][Authorize] delete err: %v", err)
			return nil, err
		}
		b.logger.Infof("[BpsUserSvc][Authorize] soft-deleted id: %s", cpsAction.UniqueId)
		return nil, nil
	case string(constants.RequestEnableBPSUser):
		b.logger.Infof("[BpsUserSvc][Authorize] enabling id: %s", cpsAction.UniqueId)
		err := b.repo.EnableDisableBPSUser(ctx, actionData.UserCode, true)
		if err != nil {
			span.AddEvent("[Authorize] failed to enable BPS user", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[BpsUserSvc][Authorize] enable err: %v", err)
			return nil, err
		}
		return nil, nil
	case string(constants.RequestDisableBPSUser):
		b.logger.Infof("[BpsUserSvc][Authorize] disabling id: %s", cpsAction.UniqueId)
		err := b.repo.EnableDisableBPSUser(ctx, actionData.UserCode, false)
		if err != nil {
			span.AddEvent("[Authorize] failed to disable BPS user", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", cpsAction.UniqueId),
			))
			b.logger.Errorf("[BpsUserSvc][Authorize] disable err: %v", err)
			return nil, err
		}
		return nil, nil

	default:
		span.AddEvent("[Authorize] unsupported action", trace.WithAttributes(attribute.String("action", cpsAction.RequestAction)))
		b.logger.Errorf("[BpsUserSvc][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}
}

// FetchUserByUserCode implements service.BPSUserService.
func (b *bpsUserService) FetchUserByUserCode(ctx context.Context, userCode string) (*bpsUserDto.BPSUserResposenDTO, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchUserByUserCode", "BPS User", "FetchUserByUserCode")
	defer span.End()

	user, err := b.repo.GetByUserCode(ctx, userCode)
	if err != nil {
		span.AddEvent("[FetchUserByUserCode] failed to fetch BPS user", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", userCode),
		))
		b.logger.Errorf("[BpsUserSvc][FetchByCode] fetch err: %v", err)
		return nil, err
	}
	b.logger.Infof("[BpsUserSvc][FetchByCode] retrieved code: %s", userCode)
	return user, nil
}

func (b *bpsUserService) FetchUserByUserName(ctx context.Context, userName string) (*imodel.BPSUser, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchUserByUserName", "BPS User", "FetchUserByUserName")
	defer span.End()

	user, err := b.repo.GetByUsername(ctx, userName)
	if err != nil {
		span.AddEvent("[FetchUserByUserName] failed to fetch BPS user", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", userName),
		))
		b.logger.Errorf("[BpsUserSvc][FetchByUserName] fetch err: %v", err)
		return nil, localization.ErrorBpsUserNotFound
	}
	if user == nil {
		span.AddEvent("[FetchUserByUserName] BPS user not found", trace.WithAttributes(attribute.String("user_code", userName)))
		b.logger.Errorf("[BpsUserSvc][FetchByUserName] not found: %s", userName)
		return nil, errors.New(localization.ErrorUserNotFound.Code)
	}

	cleanBPSUser := imodel.BPSUser{
		ID:          user.ID,
		UserCode:    user.UserCode,
		FullName:    user.FullName,
		UserName:    user.Username,
		PhoneNumber: user.PhoneNumber,
		BranchCode:  user.BranchCode,
		Email:       user.Email,
		BranchName:  user.BranchName,
		HomeBranch:  user.HomeBranch,
		JobTitle:    user.JobTitle,
		Enabled:     user.Enabled,
	}

	jobRoles, err := b.Job_roles_repo.FindByName(ctx, user.JobTitle)
	if err != nil {
		span.AddEvent("[FetchUserByUserName] failed to fetch role by job title", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("job_title", user.JobTitle),
		))
		b.logger.Errorf("[BpsUserSvc][FetchByUserName] failed to fetch role by job title: %s, err: %v", user.JobTitle, err)
		// return nil, errors.New(localization.ErrorRoleNotFound.Code)
	}
	if jobRoles != nil && jobRoles.Role != "" {
		span.AddEvent("[FetchUserByUserName] role found for job title", trace.WithAttributes(attribute.String("job_title", user.JobTitle)))
		b.logger.Errorf("[BpsUserSvc][FetchByUserName] role found for job title: %s", user.JobTitle)
		// return nil, errors.New(localization.ErrorRoleNotFound.Code)
		cleanBPSUser.Role = jobRoles.Role

		role, err := b.RoleRepository.FindByCode(ctx, jobRoles.Role)
		if err != nil {
			span.AddEvent("[FetchUserByUserName] failed to fetch role by code", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("role_code", jobRoles.Role),
			))
			b.logger.Errorf("[BpsUserSvc][FetchByUserName] failed to fetch role by code: %s, err: %v", jobRoles.Role, err)
			// return nil, errors.New(localization.ErrorRoleNotFound.Code)
		} else if role != nil && role.Name != "" {
			cleanBPSUser.RoleName = role.Name
		}
	}

	b.logger.Infof("[BpsUserSvc][FetchByUserName] retrieved username: %s", userName)
	return &cleanBPSUser, nil
}

// GetAllBPSUsers implements service.BPSUserService.
func (b *bpsUserService) GetAllBPSUsers(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]bpsUserDto.BPSUserResposenDTO], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllBPSUsers", "BPS User", "GetAllBPSUsers")
	defer span.End()

	result, err := b.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("[GetAllBPSUsers] failed to fetch BPS users", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		b.logger.Errorf("[BpsUserSvc][GetAll] fetch err: %v", err)
		return nil, err
	}
	b.logger.Infof("[BpsUserSvc][GetAll] retrieved %d", len(result.Data))
	return result, nil
}

// UpdateBpsUser implements service.BPSUserService.
func (b *bpsUserService) UpdateStatusBpsUser(ctx context.Context, userCode string, status bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateBpsUser", "BPS User", "UpdateBpsUser")
	defer span.End()

	b.logger.Infof("[BpsUserSvc][UpdateStatus] enabled: %v", status)
	makerData := local_util.ExtractUserFromContext(ctx)
	user, err := b.repo.GetByUserCode(ctx, userCode)
	if err != nil {
		span.AddEvent("[UpdateBpsUser] failed to find BPS user", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", userCode),
		))
		b.logger.Errorf("[BpsUserSvc][UpdateStatus] find err: %v", err)
		return err
	}

	if user == nil {
		span.AddEvent("[UpdateBpsUser] BPS user not found", trace.WithAttributes(attribute.String("user_code", userCode)))
		b.logger.Errorf("[BpsUserSvc][UpdateStatus] not found: %s", userCode)
		return errors.New(localization.ErrorUserNotFound.Code)
	}

	if user.Enabled && status {
		span.AddEvent("[UpdateBpsUser] BPS user already enabled", trace.WithAttributes(attribute.String("user_code", userCode)))
		b.logger.Errorf("[BpsUserSvc][UpdateStatus] already enabled: %s", userCode)
		return errors.New(localization.ErrorUserAlreadyEnabled.Code)
	}

	if !user.Enabled && !status {
		span.AddEvent("[UpdateBpsUser] BPS user already disabled", trace.WithAttributes(attribute.String("user_code", userCode)))
		b.logger.Errorf("[BpsUserSvc][UpdateStatus] already disabled: %s", userCode)
		return errors.New(localization.ErrorUserAlreadyDisabled.Code)
	}

	updatedUser := *user
	updatedUser.Enabled = status

	var requestAction string
	if status {
		requestAction = string(constants.RequestEnableBPSUser)
	} else {
		requestAction = string(constants.RequestDisableBPSUser)
	}

	cpsActionData := lib.CpsModelBuilder(user.UserCode, makerData, user, updatedUser, requestAction, constants.UPDATE)

	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		span.AddEvent("[UpdateBpsUser] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", userCode),
		))
		b.logger.Errorf("[BpsUserSvc][UpdateStatus] cps action err: %v", err)
		return err
	}
	b.logger.Infof("[BpsUserSvc][UpdateStatus] request created code: %s", userCode)
	return nil
}
func (b *bpsUserService) CreateBPSUser(ctx context.Context, req bps_model.BPSUser) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateBPSUser", "BPS User", "CreateBPSUser")
	defer span.End()
	makerData := local_util.ExtractUserFromContext(ctx)

	existing, err := b.repo.FindByOr(ctx, req.PhoneNumber, req.Email, req.Username)
	if err != nil {
		if err.Error() != localization.ErrorResourceNotFound.Code {
			b.logger.Errorf("[BpsUserSvc][Create] check existing err: %v", err)
			return err
		}
	}

	if err := bps_user_core.ExistingIdentifier(existing, req); err != nil {
		b.logger.Infof("[BpsUserSvc][Create] duplicate data: %v", err)
		return err
	}

	roles, err := b.Job_roles_repo.FindByFilterKey(ctx, "job_title", req.JobTitle)
	if err != nil {
		b.logger.Errorf("[BpsUserSvc][Create] role lookup err for job_title: %s, err: %v", req.JobTitle, err)
		return errors.New(localization.ErrorRoleNotFound.Code)
	}
	if roles == nil || roles.Role == "" {
		b.logger.Errorf("[BpsUserSvc][Create] role not found for job_title: %s", req.JobTitle)
		return errors.New(localization.ErrorRoleNotFound.Code)
	}
	req.Role = roles.Role

	is_exist_on_CPS, err := b.CPSUserRepo.FindByEmailOrPhoneNumberOrUserName(ctx, req.Email, req.PhoneNumber, req.Username)
	if err != nil {
		b.logger.Errorf("[CreateBPSUser] got error while checking user data exist on cps user error: %v", err.Error())
		if err.Error() != localization.ErrorResourceNotFound.Code {
			return errors.New(localization.ErrorInternalServerError.Code)
		}
	}

	b.logger.Infof("[CreateBPSUser] existing user on CPS: %v", is_exist_on_CPS)
	if is_exist_on_CPS != nil {
		b.logger.Infof("[CreateBPSUser] found existing user on CPS----------------: %v", is_exist_on_CPS)
		if is_exist_on_CPS.UserName != "" && is_exist_on_CPS.UserName == req.Username {
			b.logger.Errorf("[CreateBPSUser] user name already exist")
			return errors.New(localization.ErrorExistUserName.Code)
		}

		if is_exist_on_CPS.Email != "" && is_exist_on_CPS.Email == req.Email {
			b.logger.Errorf("[CreateBPSUser] email already exist")
			return errors.New(localization.ErrorExistEmail.Code)

		}

		if is_exist_on_CPS.PhoneNumber != "" && is_exist_on_CPS.PhoneNumber == req.PhoneNumber {
			b.logger.Errorf("[CreateBPSUser] phone number already exist")
			return errors.New(localization.ErrorExistPhoneNumber.Code)
		}
	}

	if b.Branch_blocks == nil {
		b.logger.Errorf("[CreateBPSUser] account block repository is not configured")
		return errors.New(localization.ErrorInternalServerError.Code)
	}
	if len(req.BranchCode) == 0 {
		return errors.New(localization.ErrorBranchCodeRequired.Code)
	}
	branchCode := strings.TrimSpace(req.BranchCode[0])
	if branchCode == "" {
		return errors.New(localization.ErrorBranchCodeRequired.Code)
	}

	branch_detail, err := b.Branch_blocks.FindByFilterKey(ctx, "code", branchCode)
	if err != nil {
		b.logger.Errorf("[CreateBPSUser] Get error while locking branch name by branch code")
		if err.Error() == localization.ErrorResourceNotFound.Code {
			return errors.New(localization.ErrorBranchNotExistWithGivenBranchCode.Code)
		}
		return err
	}

	if branch_detail == nil {
		b.logger.Warnf("[CreateBPSUser] branch not found with given branch code")
		return errors.New(localization.ErrorBranchNotExistWithGivenBranchCode.Code)
	}

	req.BranchName = branch_detail.Name
	req.UserCode = local_util.GenerateCPSUserCode() // the library only just to create the user code not portal specific

	// Build CPS action model for create
	cpsActionModel := lib.CpsModelBuilder(
		req.UserCode,                           // unique id
		makerData,                              // maker data
		nil,                                    // old data (nil for create)
		req,                                    // new data
		string(constants.RequestCreateBPSUser), // request action
		constants.CREATE,                       // action type
	)

	// Create CPS action
	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("[CreateBPSUser] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", req.UserCode),
		))
		b.logger.Errorf("[BpsUserSvc][Create] cps action err: %v", err)
		return err
	}

	b.logger.Infof("[BpsUserSvc][Create] request created code: %s", req.UserCode)
	return nil
}

func (b *bpsUserService) UpdateBPSUser(ctx context.Context, userID string, updatedUser bps_model.BPSUser) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateBPSUser", "BPS User", "UpdateBPSUser")
	defer span.End()
	makerData := local_util.ExtractUserFromContext(ctx)

	curUser, err := b.repo.GetByUserID(ctx, userID)
	if err != nil {
		span.AddEvent("[UpdateBPSUser] failed to find BPS user", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_id", userID),
		))
		return err
	}

	existing, err := b.repo.FindByOr(ctx, updatedUser.PhoneNumber, updatedUser.Email, updatedUser.Username)
	if err != nil {
		if err.Error() != localization.ErrorResourceNotFound.Code {
			b.logger.Errorf("[BpsUserSvc][Update] check existing err: %v, existing: %v, userID: %s, updatedUser: %v", err, existing, userID, updatedUser)
			return err
		}
	}

	if existing != nil {
		b.logger.Infof("[BpsUserSvc][Update] found existing user: %v for update with userID: %s, updatedUser: %v", existing, userID, updatedUser)
		if err := bps_user_core.ExistingIdentifierForUpdate(*existing, userID, updatedUser); err != nil {
			b.logger.Infof("[BpsUserSvc][Update] duplicate data: %v", err)
			if !errors.Is(err, localization.ErrorUserNotFound) {
				return err
			}
		}
	}

	updatedUser.Enabled = curUser.Enabled
	is_exist_on_CPS, err := b.CPSUserRepo.FindByEmailOrPhoneNumberOrUserName(ctx, updatedUser.Email, updatedUser.PhoneNumber, updatedUser.Username)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		b.logger.Errorf("[UpdateBPSUser] got error while checking user data exist on cps user ")
		return localization.ErrorInternalServerError
	}
	if is_exist_on_CPS != nil {
		if is_exist_on_CPS.UserName != "" && is_exist_on_CPS.UserName == updatedUser.Username {
			b.logger.Errorf("[UpdateBPSUser] user name already exist")
			return errors.New(localization.ErrorExistUserName.Code)
		}

		if is_exist_on_CPS.Email != "" && is_exist_on_CPS.Email == updatedUser.Email {
			b.logger.Errorf("[UpdateBPSUser] email already exist")
			return errors.New(localization.ErrorExistEmail.Code)

		}

		if is_exist_on_CPS.PhoneNumber != "" && is_exist_on_CPS.PhoneNumber == updatedUser.PhoneNumber {
			b.logger.Errorf("[UpdateBPSUser] phone number already exist")
			return errors.New(localization.ErrorExistPhoneNumber.Code)
		}
	}

	if len(updatedUser.BranchCode) > 0 {
		branch_detail, err := b.Branch_blocks.FindByFilterKey(ctx, "code", strings.TrimSpace(updatedUser.BranchCode[0]))
		if err != nil {
			if err.Error() != localization.ErrorResourceNotFound.Code {
				b.logger.Errorf("[UpdateBPSUser] branch not found")
				return errors.New(localization.ErrorBranchNotFoundRequired.Code)
			} else {
				b.logger.Errorf("[UpdateBPSUser] Get error while locking branch name by branch code")
				return errors.New(localization.ErrorInternalServerError.Code)
			}
		}
		if branch_detail == nil {
			b.logger.Warnf("[UpdateBPSUser] branch not found with given branch code")
			return errors.New(localization.ErrorBranchNotExistWithGivenBranchCode.Code)
		}
		updatedUser.BranchName = branch_detail.Name
	}

	var roles *imodel.JobRole
	if updatedUser.JobTitle != "" {
		roles, err = b.Job_roles_repo.FindByFilterKey(ctx, "job_title", updatedUser.JobTitle)
		if err != nil {
			b.logger.Errorf("[BpsUserSvc][Create] role lookup err for job_title: %s, err: %v", updatedUser.JobTitle, err)
			return errors.New(localization.ErrorRoleNotFound.Code)
		}
		if roles == nil || roles.Role == "" {
			b.logger.Errorf("[BpsUserSvc][Create] role not found for job_title: %s", updatedUser.JobTitle)
			return errors.New(localization.ErrorRoleNotFound.Code)
		}
		updatedUser.Role = roles.Role
	}

	updatedUser = bps_user_core.BuildUpdatedBPSUser(*curUser, updatedUser)
	updatedUser.ID = curUser.ID
	updatedUser.UserCode = curUser.UserCode
	// updatedUser.Role = roles.Role
	cpsActionModel := lib.CpsModelBuilder(
		curUser.UserCode,                       // unique id
		makerData,                              // maker data
		curUser,                                // old data
		updatedUser,                            // new data
		string(constants.RequestBpsUserUpdate), // request action
		constants.UPDATE,                       // action type
	)

	// Create CPS action
	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("[UpdateBPSUser] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_code", userID),
		))
		b.logger.Errorf("[BpsUserSvc][Update] cps action err: %v", err)
		return err
	}

	b.logger.Infof("[BpsUserSvc][Update] request created id: %s", userID)
	return nil
}

func (b *bpsUserService) DeleteBPSUser(ctx context.Context, userID string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteBPSUser", "BPS User", "DeleteBPSUser")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
			attribute.String("user_id", userID),
		))
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existingUser, err := b.repo.GetByUserID(ctx, userID)
	if err != nil {
		span.AddEvent("[DeleteBPSUser] failed to find BPS user", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_id", userID),
		))
		b.logger.Errorf("[BpsUserSvc][Delete] find err: %v", err)
		return err
	}
	if existingUser == nil {
		span.AddEvent("BPS user not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorUserNotFound.Code),
			attribute.String("user_id", userID),
		))
		return errors.New(localization.ErrorUserNotFound.Code)
	}
	if existingUser.IsDeleted {
		span.AddEvent("BPS user already deleted", trace.WithAttributes(
			attribute.String("error", localization.ErrorAlreadyDeleted.Code),
			attribute.String("user_id", userID),
		))
		return errors.New(localization.ErrorAlreadyDeleted.Code)
	}

	updatedUser := *existingUser
	updatedUser.IsDeleted = true
	updatedUser.LastModifiedAt = time.Now()
	// fmt.Println("LOLOLOLO IN FUNCTION DELETE")
	cpsActionData := lib.CpsModelBuilder(existingUser.UserCode, makerData, existingUser, updatedUser, string(constants.RequestBpsUserDelete), constants.DELETE)
	if err := b.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		span.AddEvent("[DeleteBPSUser] failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("user_id", userID),
		))
		b.logger.Errorf("[BpsUserSvc][Delete] cps action err: %v", err)
		return err
	}

	b.logger.Infof("[BpsUserSvc][Delete] request created id: %s", userID)
	return nil
}
