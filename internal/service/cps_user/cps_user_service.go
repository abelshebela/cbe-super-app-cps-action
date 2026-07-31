package cpsuser

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants"
	cpsuser "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/cps_user"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/lib"
	localization "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/service/cps_user/core"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	local_util "github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type cpsUserService struct {
	repo              storage.CpsUserRepository
	jobRoleRepo       storage.JobRoleRepository
	approverRepo      storage.CPSActionApproveIndexRepository
	bpsApproverRepo   storage.BPSActionApproveIndexRepository
	permissionService service.PermissionService
	departmentRepo    storage.DepartmentRepository
	roleDelegation    storage.RoleDelegationRepository
	roleRepo          storage.RoleRepository
	logger            shared_utils.Logger
	cpsService        service.CPSActionService
	bpsRepo           storage.BPSUserRepository
	minioClient       *s3.Client
	bucketName        string
	cfg               config.VaultConfig
}

func NewCPSUserService(repo storage.CpsUserRepository, JobRoleRepo storage.JobRoleRepository, roleRepo storage.RoleRepository, approverRepo storage.CPSActionApproveIndexRepository, bpsApproverRepo storage.BPSActionApproveIndexRepository, departmentRepo storage.DepartmentRepository, permission service.PermissionService, cps service.CPSActionService, bps storage.BPSUserRepository, roleDelegation storage.RoleDelegationRepository, minioClient *s3.Client, bucketName string, cfg config.VaultConfig, logger shared_utils.Logger) service.CPSUserService {
	return &cpsUserService{
		repo:              repo,
		bpsRepo:           bps,
		jobRoleRepo:       JobRoleRepo,
		roleRepo:          roleRepo,
		approverRepo:      approverRepo,
		bpsApproverRepo:   bpsApproverRepo,
		roleDelegation:    roleDelegation,
		permissionService: permission,
		cpsService:        cps,
		logger:            logger,
		departmentRepo:    departmentRepo,
		minioClient:       minioClient,
		bucketName:        bucketName,
		cfg:               cfg,
	}
}

func cpsUserExportRow(user imodel.ExportCPSUser) []string {
	expiryDate := ""
	userType := "Permanent"
	if !user.ExpiryDateForDelegation.IsZero() {
		userType = "Delegation"
		expiryDate = user.ExpiryDateForDelegation.Format(time.RFC3339)

	}
	lastModified := ""
	if !user.LastModified.IsZero() {
		lastModified = user.LastModified.Format(time.RFC3339)
	}
	lastLogin := ""
	if !user.LastLogin.IsZero() {
		lastLogin = user.LastLogin.Format(time.RFC3339)
	}
	return []string{
		user.FirstName,
		user.PhoneNumber,
		user.Email,
		user.Department,
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
		user.CreatedBy,
		user.ApprovedBy,
	}
}

func (s *cpsUserService) ExportUsers(ctx context.Context, startDate, endDate time.Time, fileType, userName string) (string, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	if s.minioClient == nil {
		log.Errorf("[CPSUser/ExportUsers] minio client is not configured")
		return "", errors.New(localization.CpsUserDataExportedError.Code)
	}

	fileType = strings.ToLower(strings.TrimSpace(fileType))
	if fileType != string(lib.FileTypeCSV) && fileType != string(lib.FileTypePDF) {
		return "", errors.New(localization.ErrorInvalidRequest.Code)
	}

	data, err := s.repo.FindForExport(ctx, startDate, endDate, userName)
	if err != nil {
		log.Errorf("[CPSUser/ExportUsers] failed to fetch users: %v", err)
		return "", err
	}

	if len(data) == 0 {
		return "", errors.New(localization.CpsUserDataNotFoundInDateRange.Code)
	}

	headers := []string{
		"Full Name", "Phone Number", "Email", "Department", "Job Title", "Role",
		"Username", "Created At", "Last Login", "Enabled",
		"User Type", "Expiry Date (Delegation)",
		"Last Modification Action", "Last Modified",
		"Created By", "Approved By",
	}

	ext := "csv"
	if fileType == string(lib.FileTypePDF) {
		ext = "pdf"
	}

	objectName := fmt.Sprintf("cps_users_%s_to_%s_%d.%s", startDate.Format("20060102"), endDate.Format("20060102"), time.Now().Unix(), ext)

	if fileType == string(lib.FileTypePDF) {
		rows := make([][]string, 0, len(data))
		for _, item := range data {
			rows = append(rows, cpsUserExportRow(item))
		}
		url, exportErr := lib.ExportPDFAndUpload(ctx, s.minioClient, s.bucketName, s.cfg, objectName, headers, rows, lib.PDFExportOptions{PageSize: "A4"}, nil, s.logger)
		if exportErr != nil {
			log.Errorf("[CPSUser/ExportUsers] pdf export failed: %v", exportErr)
			return "", errors.New(localization.CpsUserDataExportedError.Code)
		}

		return url, nil
	}

	url, exportErr := lib.ExportCSVAndUpload(ctx, s.minioClient, s.bucketName, s.cfg, objectName, headers, func(writer *csv.Writer) error {
		for _, item := range data {
			if err := writer.Write(cpsUserExportRow(item)); err != nil {
				return err
			}
		}
		return nil
	}, s.logger)
	if exportErr != nil {
		log.Errorf("[CPSUser/ExportUsers] csv export failed: %v", exportErr)
		return "", errors.New(localization.CpsUserDataExportedError.Code)
	}

	return url, nil
}

func (s *cpsUserService) CreateUserRequest(ctx context.Context, req cpsuser.CreateUserRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "CreateUserRequest", "CPSUser", "CreateUserRequest")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)

	normalized := local_util.FormatPhoneNumber(req.PhoneNumber)
	req.PhoneNumber = normalized
	exists, err := core.UsernameExists(ctx, "", s.repo, req.UserName)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		span.AddEvent("failed to check username existence", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	if exists {
		log.Errorf("[CpsUserSvc][Create] username exists: %s", req.UserName)
		span.AddEvent("username already exists", trace.WithAttributes(attribute.String("username", req.UserName)))
		return errors.New(localization.ErrorUsernameAlreadyExists.Code)
	}
	emailCheck, err := core.EmailExists(ctx, "", s.repo, req.Email)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		span.AddEvent("failed to check email existence", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	if emailCheck {
		span.AddEvent("email already exists", trace.WithAttributes(attribute.String("email", req.Email)))
		return errors.New(localization.ErrorExistEmail.Code)
	}

	phoneCheck, err := core.PhoneNumberExists(ctx, "", s.repo, req.PhoneNumber)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		span.AddEvent("failed to check phone number existence", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	if phoneCheck {
		span.AddEvent("phone number already exists", trace.WithAttributes(attribute.String("phone_number", req.PhoneNumber)))
		return errors.New(localization.ErrorExistPhoneNumber.Code)
	}

	emailCheckBPS, err := s.bpsRepo.FindByOr(ctx, req.PhoneNumber, req.Email, req.UserName)
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		span.AddEvent("failed to check email existence in BPS", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	if emailCheckBPS != nil {
		if emailCheckBPS.Email == req.Email {
			span.AddEvent("email already exists in BPS", trace.WithAttributes(attribute.String("email", req.Email)))
			return errors.New(localization.ErrorExistEmail.Code)
		}
		if emailCheckBPS.PhoneNumber == req.PhoneNumber {
			span.AddEvent("phone number already exists in BPS", trace.WithAttributes(attribute.String("phone_number", req.PhoneNumber)))
			return errors.New(localization.ErrorExistPhoneNumber.Code)
		}
		if emailCheckBPS.Username == req.UserName {
			span.AddEvent("username already exists in BPS", trace.WithAttributes(attribute.String("username", req.UserName)))
			return errors.New(localization.ErrorUsernameAlreadyExists.Code)
		}

	}

	department, err := s.departmentRepo.FindByID(ctx, req.Department)
	if err != nil {
		span.AddEvent("failed to find department", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	if department == nil {
		span.AddEvent("department not found", trace.WithAttributes(attribute.String("department", req.Department)))
		return errors.New(localization.ErrorDepartmentNotFound.Code)
	}

	jobTitle, err := s.jobRoleRepo.FindByName(ctx, req.JobTitle)
	if err != nil {
		s.logger.Errorf("failed to find jobTitle err:%v", err)
		return err
	}
	if jobTitle == nil {
		s.logger.Errorf("unexpected error while fetching jobTitle err:%v", err)
		return localization.ErrorUnexpectedError
	}

	cpsUser := core.CPSUModel(req)
	cpsUser.JobTitle = jobTitle.JobTitle
	cpsUser.Role = jobTitle.Role

	userForAction := core.MapForActionWithDepartment(cpsUser, department)
	cpsActionModel := lib.CpsModelBuilder(cpsUser.UserCode, makerData, nil, userForAction, string(constants.RequestCpsUserCreate), constants.CREATE)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("failed to create CPS action", trace.WithAttributes(attribute.String("error", err.Error())))
		log.Errorf("[CpsUserSvc][Create] cps action err: %v", err)
		return err
	}
	log.Infof("[CpsUserSvc][Create] request created")
	return nil
}

func (s *cpsUserService) UpdateUserRequest(ctx context.Context, usercode string, req cpsuser.UpdateUserRequest) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateUserRequest", "CPSUser", "UpdateUserRequest")
	defer span.End()

	currentUser, err := s.repo.FindByID(ctx, usercode)
	if err != nil {
		span.AddEvent("failed to find user by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	if currentUser == nil {
		span.AddEvent("user not found", trace.WithAttributes(attribute.String("user_code", usercode)))
		return errors.New(localization.ErrorUserNotFound.Code)
	}

	if req.PhoneNumber != "" {
		normalized := local_util.FormatPhoneNumber(req.PhoneNumber)
		req.PhoneNumber = normalized
		if currentUser.PhoneNumber == req.PhoneNumber {
			req.PhoneNumber = ""
		}
	}
	if req.UserName != "" {
		if currentUser.UserName == req.UserName {
			req.UserName = ""
		} else {
			req.UserName = strings.ToUpper(req.UserName)
		}
	}
	if req.Email != "" {
		if currentUser.Email == req.Email {
			req.Email = ""
		}
	}
	if currentUser.JobTitle == req.JobTitle {
		req.JobTitle = ""
	}
	if req.PhoneNumber != "" || req.Email != "" || req.UserName != "" {
		foundUser, err := s.repo.FindByEmailOrPhoneNumberOrUserName(ctx, req.Email, req.PhoneNumber, req.UserName)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			span.AddEvent("failed to find user by email, phone number, or username", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[CpsUserSvc][Update] error checking user conflicts: %v", err)
			return err
		}
		log.Infof("[CpsUserSvc][Update] found user: %v", foundUser)
		if foundUser != nil && foundUser.UserCode != currentUser.UserCode {
			if foundUser.Email == req.Email {
				log.Infof("[CpsUserSvc][Update] email already exists: %s", req.Email)
				span.AddEvent("email already exists", trace.WithAttributes(attribute.String("email", req.Email)))
				return errors.New(localization.ErrorExistEmail.Code)
			}
			if foundUser.PhoneNumber == req.PhoneNumber {
				log.Infof("[CpsUserSvc][Update] phone number already exists: %s", req.PhoneNumber)
				span.AddEvent("phone number already exists", trace.WithAttributes(attribute.String("phone_number", req.PhoneNumber)))
				return errors.New(localization.ErrorExistPhoneNumber.Code)
			}
			if foundUser.UserName == req.UserName {
				log.Infof("[CpsUserSvc][Update] username already exists: %s", req.UserName)
				span.AddEvent("username already exists", trace.WithAttributes(attribute.String("username", req.UserName)))
				return errors.New(localization.ErrorUserAlreadyExists.Code)
			}
		}
		log.Infof("[CpsUserSvc][Update] no conflicts found, proceeding with update")
		bpsUser, err := s.bpsRepo.FindByOr(ctx, req.PhoneNumber, req.Email, req.UserName)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			span.AddEvent("failed to find user by email, phone number, or username", trace.WithAttributes(attribute.String("error", err.Error())))
			log.Errorf("[CpsUserSvc][Update] error checking BPS user: %v", err)
			return err
		}

		if bpsUser != nil {
			if bpsUser.Email == req.Email {
				log.Infof("[CpsUserSvc][Update] email already exists: %s", req.Email)
				span.AddEvent("email already exists", trace.WithAttributes(attribute.String("email", req.Email)))
				return errors.New(localization.ErrorExistEmail.Code)
			}
			if bpsUser.PhoneNumber == req.PhoneNumber {
				log.Infof("[CpsUserSvc][Update] phone number already exists: %s", req.PhoneNumber)
				span.AddEvent("phone number already exists", trace.WithAttributes(attribute.String("phone_number", req.PhoneNumber)))
				return errors.New(localization.ErrorExistPhoneNumber.Code)
			}
			if bpsUser.Username == req.UserName {
				log.Infof("[CpsUserSvc][Update] username already exists: %s", req.UserName)
				span.AddEvent("username already exists", trace.WithAttributes(attribute.String("username", req.UserName)))
				return errors.New(localization.ErrorUserAlreadyExists.Code)
			}
		}
	}

	var department *model.Department
	if req.Department != "" {
		department, err = s.departmentRepo.FindByID(ctx, req.Department)
		if err != nil {
			span.AddEvent("failed to find department", trace.WithAttributes(attribute.String("error", err.Error())))
			return err
		}
		if department == nil {
			span.AddEvent("department not found", trace.WithAttributes(attribute.String("department", req.Department)))
			return errors.New(localization.ErrorDepartmentNotFound.Code)
		}
	}
	changed := false
	updated := *currentUser
	if req.UserName != "" {
		changed = true
		updated.UserName = req.UserName
	}
	if req.FullName != "" {
		changed = true
		updated.FullName = req.FullName
	}
	if req.PhoneNumber != "" {
		changed = true
		updated.PhoneNumber = req.PhoneNumber
	}
	if req.Gender != "" {
		changed = true
		updated.Gender = req.Gender
	}
	if req.Email != "" {
		changed = true
		updated.Email = req.Email
	}
	if req.JobTitle != "" {
		changed = true
		updated.JobTitle = req.JobTitle
	}
	if req.Department != "" {
		changed = true
		updated.Department, err = bson.ObjectIDFromHex(req.Department)
		if err != nil {
			span.AddEvent("failed to convert department id", trace.WithAttributes(attribute.String("error", err.Error())))
			return errors.New("invalid department id")
		}
	}
	if !changed {
		span.AddEvent("no fields to update", trace.WithAttributes(attribute.String("user_code", usercode)))
		log.Infof("[CpsUserSvc][Update] no fields to update for user code: %s", usercode)
		return errors.New("no fields to update")
	}

	makerData := local_util.ExtractUserFromContext(ctx)
	userForAction := core.MapForActionWithDepartment(updated, department)

	curDpt, err := s.departmentRepo.FindByID(ctx, currentUser.Department.Hex())
	if err != nil {
		span.AddEvent("failed to find department", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	dateJoined := time.Time{}
	if currentUser.DateJoined != nil {
		dateJoined = *currentUser.DateJoined
	}

	lastModified := time.Time{}
	if currentUser.LastModified != nil {
		lastModified = *currentUser.LastModified
	}

	var departmentResp *cpsuser.DepartmentResponse
	if curDpt != nil {
		departmentResp = &cpsuser.DepartmentResponse{ID: curDpt.ID, Name: curDpt.Department}
	}

	prevUserData := cpsuser.CpsUserPopulatedResponse{
		ID:                 currentUser.ID,
		UserCode:           currentUser.UserCode,
		FullName:           currentUser.FullName,
		Role:               cpsuser.RoleResponse{Name: currentUser.Role},
		Department:         departmentResp,
		JobTitle:           currentUser.JobTitle,
		Gender:             currentUser.Gender,
		PhoneNumber:        currentUser.PhoneNumber,
		Email:              currentUser.Email,
		UserName:           currentUser.UserName,
		Realm:              currentUser.Realm,
		Enabled:            currentUser.Enabled,
		DateJoined:         dateJoined,
		LastModified:       lastModified,
		Country:            currentUser.Country,
		Region:             currentUser.Region,
		PermissionCategory: currentUser.PermissionCategory,
		LastLogin:          currentUser.LastLogin,
		PasswordDisable:    currentUser.PasswordDisable,
		IsFirstTimeLogin:   currentUser.IsFirstTimeLogin,
		CreatedAt:          currentUser.CreatedAt,
	}

	cpsActionModel := lib.CpsModelBuilder(
		usercode,
		makerData,
		prevUserData,
		userForAction,
		string(constants.RequestCpsUserUpdate),
		constants.UPDATE,
	)

	if err := s.cpsService.CreateCPSAction(ctx, &cpsActionModel); err != nil {
		span.AddEvent("failed to create CPS action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	return nil
}

func (s *cpsUserService) DeleteUserRequest(ctx context.Context, userCode string) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteUserRequest", "CPSUser", "DeleteUserRequest")
	defer span.End()

	if userCode == "" {
		span.AddEvent("user code is empty", trace.WithAttributes(attribute.String("error", "user code is empty")))
		return errors.New(localization.ErrorUserCodeRequired.Code)
	}

	maker := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(maker); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
			attribute.String("user_code", userCode),
		))
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByID(ctx, userCode)
	if err != nil {
		span.AddEvent("failed to find user by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	if existing == nil {
		span.AddEvent("CPS user not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorUserNotFound.Code),
			attribute.String("user_code", userCode),
		))
		return errors.New(localization.ErrorUserNotFound.Code)
	}
	if existing.IsDeleted {
		span.AddEvent("CPS user already deleted", trace.WithAttributes(
			attribute.String("error", localization.ErrorAlreadyDeleted.Code),
			attribute.String("user_code", userCode),
		))
		return errors.New(localization.ErrorAlreadyDeleted.Code)
	}

	updatedUser := *existing
	updatedUser.IsDeleted = true
	now := time.Now()
	updatedUser.LastModified = &now

	cpsAction := lib.CpsModelBuilder(userCode, maker, existing, updatedUser, string(constants.RequestCpsUserDelete), constants.DELETE)
	if err := s.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("failed to create cps action", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	log.Infof("[CpsUserSvc][Delete] request created code: %s", userCode)
	return nil
}

func (s *cpsUserService) EnableUser(ctx context.Context, userCode string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableUser", "CPSUser", "EnableUser")
	defer span.End()

	if userCode == "" {
		span.AddEvent("user code is empty", trace.WithAttributes(attribute.String("error", "user code is empty")))
		return errors.New(localization.ErrorUserCodeRequired.Code)
	}
	prev, err := s.repo.FindByID(ctx, userCode)
	if err != nil {
		span.AddEvent("failed to find user by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	if prev.Enabled {
		span.AddEvent("user already enabled", trace.WithAttributes(attribute.String("user_code", userCode)))
		return errors.New(localization.ErrorUserAlreadyEnabled.Code)
	}

	updated := *prev
	updated.Enabled = true
	updated.LoginAttemptCount = 0
	updated.PasswordDisable = false

	maker := local_util.ExtractUserFromContext(ctx)
	department, err := s.departmentRepo.FindByID(ctx, prev.Department.Hex())
	if err != nil {
		span.AddEvent("failed to find department", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	userForAction := core.MapForActionWithDepartment(updated, department)

	cpsAction := lib.CpsModelBuilder(userCode, maker, prev, userForAction, string(constants.RequestCpsUserEnable), constants.UPDATE)
	err = s.cpsService.CreateCPSAction(ctx, &cpsAction)
	if err != nil {
		span.AddEvent("failed to create cps action", trace.WithAttributes(attribute.String("error", err.Error())))
	}
	return err
}

func (s *cpsUserService) DisableUser(ctx context.Context, userCode string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DisableUser", "CPSUser", "DisableUser")
	defer span.End()

	if userCode == "" {
		span.AddEvent("user code is empty", trace.WithAttributes(attribute.String("error", "user code is empty")))
		return errors.New(localization.ErrorUserCodeRequired.Code)
	}

	prev, err := s.repo.FindByID(ctx, userCode)
	if err != nil {
		span.AddEvent("failed to find user by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}

	if !prev.Enabled {
		span.AddEvent("user already disabled", trace.WithAttributes(attribute.String("user_code", userCode)))
		return errors.New(localization.ErrorUserAlreadyDisabled.Code)
	}

	updated := *prev
	updated.Enabled = false
	updated.PasswordDisable = true

	maker := local_util.ExtractUserFromContext(ctx)
	department, err := s.departmentRepo.FindByID(ctx, prev.Department.Hex())
	if err != nil {
		span.AddEvent("failed to find department", trace.WithAttributes(attribute.String("error", err.Error())))
		return err
	}
	userForAction := core.MapForActionWithDepartment(updated, department)

	cpsaction := lib.CpsModelBuilder(userCode, maker, prev, userForAction, string(constants.RequestCpsUserDisable), constants.UPDATE)
	err = s.cpsService.CreateCPSAction(ctx, &cpsaction)
	if err != nil {
		span.AddEvent("failed to create cps action", trace.WithAttributes(attribute.String("error", err.Error())))
	}
	return err
}

func (s *cpsUserService) FetchUserByUserCode(ctx context.Context, userCode string) (*cpsuser.CPSUserResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchUserByUserCode", "CPSUser", "FetchUserByUserCode")
	defer span.End()

	if userCode == "" {
		span.AddEvent("user code is empty", trace.WithAttributes(attribute.String("error", "user code is empty")))
		return nil, errors.New(localization.ErrorUserCodeRequired.Code)
	}

	user, err := s.repo.FindByID(ctx, userCode)
	if err != nil {
		span.AddEvent("failed to find user by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}
	var roles *imodel.JobRole
	if user.IsDelegationActive {

		s.logger.Infof("fetching role by delegated role code: %s", user.DelegatedRoleCode)
		roles, err = s.jobRoleRepo.FindByRole(ctx, user.DelegatedRoleCode)
		if err != nil {
			s.logger.Errorf("failed to find role by delegation err:%v", err)
			span.AddEvent("failed to find role by role", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		s.logger.Infof("found role by delegation: %v", roles)
	} else {
		s.logger.Infof("fetching role by job title: %s", user.JobTitle)
		roles, err = s.jobRoleRepo.FindByName(ctx, user.JobTitle)
		if err != nil {
			s.logger.Errorf("failed to find role by job title err:%v", err)
			span.AddEvent("failed to find role by name", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		s.logger.Infof("found role by job title: %v", roles)
	}

	var makerAlloc, checkerAlloc, auditorAlloc, portalCard, bpsCheckerAlloc, bpsAuditorAlloc []string
	if roles.Enabled {
		_, makerAlloc, checkerAlloc, auditorAlloc, portalCard, err = s.approverRepo.PopulateUserApproverAllocations(ctx, roles.Role)
		if err != nil {
			span.AddEvent("failed to populate user approver allocations", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
	}

	_, _, bpsCheckerAlloc, bpsAuditorAlloc, err = s.bpsApproverRepo.PopulateUserApproverAllocations(ctx, roles.Role)
	if err != nil {
		span.AddEvent("failed to populate user approver allocations", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}

	userData, err := local_util.JsonUnmarshal[cpsuser.CpsUserResponse](user)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	return core.ConvertToDTO(portalCard, userData, makerAlloc, checkerAlloc, auditorAlloc, bpsCheckerAlloc, bpsAuditorAlloc), nil
}

func (s *cpsUserService) GetPopulatedCpsUser(ctx context.Context, userCode string) (*cpsuser.CpsUserResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetPopulatedCpsUser", "CPSUser", "GetPopulatedCpsUser")
	defer span.End()

	if userCode == "" {
		span.AddEvent("user code is empty", trace.WithAttributes(attribute.String("error", "user code is empty")))
		return nil, errors.New(localization.ErrorUserCodeRequired.Code)
	}

	user, err := s.repo.GetPopulatedByID(ctx, userCode)
	// user, err := s.repo.GetPopulatedByID(ctx, userCode)
	if err != nil {
		span.AddEvent("failed to get populated by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err

	}

	return user, nil
}

func (s *cpsUserService) GetCpsUserDetailByCode(ctx context.Context, userCode string) (imodel.CPSUser, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetPopulatedCpsUserByUserName", "CPSUser", "GetPopulatedCpsUserByUserName")
	defer span.End()

	data, err := s.repo.GetByUserCode(ctx, userCode)
	if err != nil {
		span.AddEvent("failed to get populated by id", trace.WithAttributes(attribute.String("error", err.Error())))
		return imodel.CPSUser{}, err
	}

	// fetch department
	curDpt, err := s.departmentRepo.FindByID(ctx, data.Department.Hex())
	if err != nil {
		span.AddEvent("failed to find department", trace.WithAttributes(attribute.String("error", err.Error())))
		return imodel.CPSUser{}, localization.ErrorUnexpectedError
	}
	data.DepartmentName = curDpt.Department
	return data, nil

}
func (s *cpsUserService) GetCpsUserDetail(ctx context.Context, userCode string) (*cpsuser.CpsUserPopulatedResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCpsUserDetail", "CPSUser", "GetCpsUserDetail")

	defer span.End()

	if userCode == "" {
		span.AddEvent("user code is empty", trace.WithAttributes(attribute.String("error", "user code is empty")))
		return nil, errors.New(localization.ErrorUserCodeRequired.Code)
	}

	// populated, err := s.repo.GetPopulatedByID(ctx, userCode)
	populated, err := s.repo.GetPopulatedWithRole(ctx, userCode)
	if err != nil {
		span.AddEvent("failed to get populated by id", trace.WithAttributes(attribute.String("error", err.Error())))
		if !errors.Is(err, localization.ErrorUserNotFound) {
			s.logger.Errorf("[GetCpsUserDetail] cps user not found for user code: %s in GetPopulatedWithRole", userCode)
			return nil, err
		}

		populated, err = s.roleDelegation.FindByUserCode(ctx, userCode)
		if err != nil {
			span.AddEvent("failed to find role delegation by user code", trace.WithAttributes(attribute.String("error", err.Error())))
			s.logger.Errorf("failed to find role delegation by user code err:%v", err)
			return nil, err
		}
	}

	s.logger.Infof("[GetCpsUserDetail] fetched populated user: %+v", populated)

	var makerAlloc, checkerAlloc, auditorAlloc, portalCard, bpsCheckerAlloc, bpsAuditorAlloc []string
	var roles *imodel.JobRole
	if populated.IsDelegationActive {
		s.logger.Infof("[GetCpsUserDetail]fetching role by delegated role: %s", populated.DelegatedRole)
		rolesD, err := s.roleRepo.FindByCode(ctx, populated.DelegatedRole)
		if err != nil {
			s.logger.Errorf("[GetCpsUserDetail][FindByRole] failed to find role by delegation err:%v", err)
			span.AddEvent("failed to find role by delegation", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		populated.Role.Code = rolesD.Code
		populated.Role.Name = rolesD.Name
		roles, err = s.jobRoleRepo.FindByRole(ctx, rolesD.Code)
		if err != nil {
			s.logger.Errorf("[GetCpsUserDetail][FindByRole] failed to find role by delegation err:%v", err)
			span.AddEvent("failed to find role by delegation", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, localization.ErrorRoleNotFound
		}

	} else {
		s.logger.Infof("[GetCpsUserDetail] fetching role by job title: %s", populated.JobTitle)
		if populated.JobTitle != "" {
			roles, err = s.jobRoleRepo.FindByName(ctx, populated.JobTitle)
			if err != nil {
				span.AddEvent("failed to find role by name", trace.WithAttributes(attribute.String("error", err.Error())))
				return nil, err
			}
		}
	}
	s.logger.Infof("[GetCpsUserDetail] fetched role: %v", roles)
	if roles != nil && roles.Enabled {
		_, makerAlloc, checkerAlloc, auditorAlloc, portalCard, err = s.approverRepo.PopulateUserApproverAllocations(ctx, roles.Role)
		if err != nil {
			span.AddEvent("failed to populate user approver allocations", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}

		_, _, bpsCheckerAlloc, bpsAuditorAlloc, err = s.bpsApproverRepo.PopulateUserApproverAllocations(ctx, roles.Role)
		if err != nil {
			span.AddEvent("failed to populate user approver allocations", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
	}

	return core.ConvertToResponseDTO(portalCard, populated, makerAlloc, checkerAlloc, auditorAlloc, bpsCheckerAlloc, bpsAuditorAlloc), nil
}
func (s *cpsUserService) GetCpsUserDetailByUserName(ctx context.Context, userName string) (*cpsuser.CpsUserPopulatedResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetCpsUserDetailByUserName", "CPSUser", "GetCpsUserDetailByUserName")

	defer span.End()

	if userName == "" {
		span.AddEvent("user name is empty", trace.WithAttributes(attribute.String("error", "user code is empty")))
		return nil, errors.New(localization.ErrorUserNameRequired.Code)
	}

	// populated, err := s.repo.GetPopulatedByID(ctx, userName)
	populated, err := s.repo.GetPopulatedWithRoleByUserName(ctx, userName)
	if err != nil {
		span.AddEvent("failed to get populated by id", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("[GetCpsUserDetailByUserName] cps user not found for user name: %s in GetPopulatedWithRoleByUserName", userName)
		return nil, err
	}

	var makerAlloc, checkerAlloc, auditorAlloc, portalCard, bpsCheckerAlloc, bpsAuditorAlloc []string
	var roles *imodel.JobRole
	if populated.JobTitle != "" {
		roles, err = s.jobRoleRepo.FindByName(ctx, populated.JobTitle)
		if err != nil {
			span.AddEvent("failed to find role by name", trace.WithAttributes(attribute.String("error", err.Error())))
			s.logger.Errorf("[GetCpsUserDetailByUserName][FindByName] failed to find role by job title err:%v", err)
			// return nil, err
		}
	}

	if roles != nil && roles.Enabled {
		_, makerAlloc, checkerAlloc, auditorAlloc, portalCard, err = s.approverRepo.PopulateUserApproverAllocations(ctx, roles.Role)
		if err != nil {
			span.AddEvent("failed to populate user approver allocations", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}

		_, _, bpsCheckerAlloc, bpsAuditorAlloc, err = s.bpsApproverRepo.PopulateUserApproverAllocations(ctx, roles.Role)
		if err != nil {
			span.AddEvent("failed to populate user approver allocations", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
	}

	return core.ConvertToResponseDTO(portalCard, populated, makerAlloc, checkerAlloc, auditorAlloc, bpsCheckerAlloc, bpsAuditorAlloc), nil
}

func (s *cpsUserService) GetAllCPSUsers(ctx context.Context, filter *types.Filter) (*types.PaginatedResponse[[]*cpsuser.CPSUserWithDepartment], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllCPSUsers", "CPSUser", "GetAllCPSUsers")
	defer span.End()

	if filter == nil {
		f := types.Filter{}
		filter = &f
	}

	users, err := s.repo.FindAllWithPagination(ctx, *filter)
	if err != nil {
		span.AddEvent("failed to find all with pagination", trace.WithAttributes(attribute.String("error", err.Error())))
		return nil, err
	}

	return users, nil
}

func (s *cpsUserService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "CPSUser", "Authorize")
	defer span.End()

	action.MakerActionTime = time.Now()
	action.LastModifiedAt = action.MakerActionTime

	switch action.RequestAction {
	case string(constants.RequestCpsUserCreate):

		userFromAction, err := local_util.JsonUnmarshal[cpsuser.CpsUserPopulatedResponse](action.CurrentAction)

		cur := core.MapFromPopulatedResponse(userFromAction)
		// cur, err := core.BindCPSUserFromAction(action.CurrentAction)
		if err != nil {
			span.AddEvent("failed to bind cps user from action", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.Create(ctx, cur); err != nil {
			span.AddEvent("failed to create user", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		return action, nil

	case string(constants.RequestCpsUserUpdate):
		userFromAction, err := local_util.JsonUnmarshal[cpsuser.CpsUserPopulatedResponse](action.CurrentAction)

		cur := core.MapFromPopulatedResponse(userFromAction)
		// cur, err := core.BindCPSUserUpdateFromAction(action.CurrentAction)
		if err != nil {
			span.AddEvent("failed to bind cps user update from action", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}

		if err := s.repo.Update(ctx, action.UniqueId, cur); err != nil {
			span.AddEvent("failed to update user", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		return action, nil

	case string(constants.RequestCpsUserDelete):
		_, err := core.BindCPSUserFromAction(action.CurrentAction)
		if err != nil {
			span.AddEvent("failed to bind cps user from action", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.Delete(ctx, action.UniqueId); err != nil {
			span.AddEvent("failed to delete user", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		return action, nil

	case string(constants.RequestCpsUserEnable):
		_, err := local_util.JsonUnmarshal[cpsuser.CpsUserPopulatedResponse](action.CurrentAction)

		// cur := core.MapFromPopulatedResponse(userFromAction)
		if err != nil {
			span.AddEvent("failed to bind cps user from action", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.EnableOrDisable(ctx, action.UniqueId, true); err != nil {
			span.AddEvent("failed to enable user", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		return action, nil

	case string(constants.RequestCpsUserDisable):
		// _, err := core.BindCPSUserFromAction(action.CurrentAction)
		_, err := local_util.JsonUnmarshal[cpsuser.CpsUserPopulatedResponse](action.CurrentAction)

		if err != nil {
			span.AddEvent("failed to bind cps user from action", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, errors.New(localization.ErrorInvalidRequest.Code)
		}
		if err := s.repo.EnableOrDisable(ctx, action.UniqueId, false); err != nil {
			span.AddEvent("failed to disable user", trace.WithAttributes(attribute.String("error", err.Error())))
			return nil, err
		}
		return action, nil
	}

	span.AddEvent("unhandled server error", trace.WithAttributes(attribute.String("error", "unhandled server error")))
	return nil, errors.New("UNHANDLED_SERVER_ERROR")
}
