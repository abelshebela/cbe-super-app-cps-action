package role_delegation_service

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/kafka"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	shared_dto "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/notification/dto"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type roleDelegation struct {
	repo         storage.RoleDelegationRepository
	jobTitleRepo storage.JobRoleRepository
	roleRepo     storage.RoleRepository
	cpsUserRepo  storage.CpsUserRepository
	bpsUserRepo  storage.BPSUserRepository
	department   storage.DepartmentRepository
	cpsService   service.CPSActionService
	branchRepo   storage.AccountBlockRepository
	minioClient  *s3.Client
	kafkaClient  *kafka.NotificationProducer
	bucketName   string
	cfg          config.VaultConfig
	logger       utils.Logger
}

// Authorize implements [service.RoleDelegationService].
func (r *roleDelegation) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	// var asAny any
	raw, err := json.Marshal(cpsAction.CurrentAction)
	log.Infof("[RoleDelegation Service][Authorize] raw CPS action data: %s", string(raw))
	if err != nil {
		log.Errorf("[RoleDelegation Service][Authorize] marshal CurrentAction failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	// if err := json.Unmarshal(raw, &asAny); err != nil {
	// 	log.Errorf("[RoleDelegation Service][Authorize] unmarshal CurrentAction failed: %v", err)
	// 	return nil, errors.New(localization.ErrorUnexpectedError.Code)
	// }
	// log.Infof("[RoleDelegation Service][Authorize] CPS action data as any: %v", asAny)
	// raw2, _ := json.Marshal(asAny)

	roleDelegation := imodel.RoleDelegation{}
	if err := json.Unmarshal(raw, &roleDelegation); err != nil {
		log.Errorf("[RoleDelegation Service][Authorize] map to RoleDelegation failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// CurrentAction payload may send Mongo id as "_id" while model tag is "id".
	var rawAction struct {
		RawID string `json:"_id"`
	}
	if err := json.Unmarshal(raw, &rawAction); err != nil {
		log.Errorf("[RoleDelegation Service][Authorize] map raw _id failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}
	if rawAction.RawID != "" {
		oid, oidErr := bson.ObjectIDFromHex(rawAction.RawID)
		if oidErr != nil {
			log.Errorf("[RoleDelegation Service][Authorize] invalid _id in CurrentAction: %v", oidErr)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		roleDelegation.ID = oid
	}

	log.Infof("[RoleDelegation Service][Authorize] CPS action data as RoleDelegation: %+v", roleDelegation)
	if cpsAction.UniqueId != "" {
		if oid, err := bson.ObjectIDFromHex(cpsAction.UniqueId); err == nil {
			roleDelegation.ID = oid
		}
	}

	switch string(cpsAction.RequestAction) {
	case string(constants.RequestCreateRoleDelegationForExistingUser):
		roleDelegation.CreatedAt = time.Now()
		if err := r.repo.CreateWithExistingUser(ctx, &roleDelegation); err != nil {
			return nil, err
		}
		emailBody := roleDelegationEmailBody(roleDelegation)
		emailMessage := shared_dto.SendEmailRequest{
			Recipients: []shared_dto.EmailContact{{Name: roleDelegation.DelegatedUserFullName, Email: roleDelegation.DelegatedUserEmail}},
			Subject:    "Role Delegation",
			Body:       emailBody,
		}
		// send mail to notify the delegated user about the new role delegation and timeframe
		if err := r.kafkaClient.PublishMessage(ctx, emailMessage, "email", "role-delegation", "role-delegation"); err != nil {
			log.Errorf("[RoleDelegation Service][Authorize] failed to send notification email: %v", err)
		}
	case string(constants.RequestCreateRoleDelegationForNewUser):
		roleDelegation.CreatedAt = time.Now()
		if err := r.repo.CreateWithNewUser(ctx, &roleDelegation); err != nil {
			return nil, err
		}
		emailBody := roleDelegationEmailBody(roleDelegation)
		emailMessage := shared_dto.SendEmailRequest{
			Recipients: []shared_dto.EmailContact{{Name: roleDelegation.DelegatedUserFullName, Email: roleDelegation.DelegatedUserEmail}},
			Subject:    "Role Delegation",
			Body:       emailBody,
		}
		// send mail to notify the delegated user about the new role delegation and timeframe
		if err := r.kafkaClient.PublishMessage(ctx, emailMessage, "email", "role-delegation", "role-delegation"); err != nil {
			log.Errorf("[RoleDelegation Service][Authorize] failed to send notification email: %v", err)
		}
	case string(constants.RequestUpdateRoleDelegation):
		roleDelegation.UpdatedAt = time.Now()
		if err := r.repo.Update(ctx, roleDelegation.ID.Hex(), &roleDelegation); err != nil {
			return nil, err
		}
	case string(constants.RequestEnableRoleDelegation):
		if err := r.repo.EnableOrDisable(ctx, cpsAction.UniqueId, true); err != nil {
			return nil, err
		}
	case string(constants.RequestDisableRoleDelegation):
		if err := r.repo.EnableOrDisable(ctx, cpsAction.UniqueId, false); err != nil {
			return nil, err
		}
	case string(constants.RequestDeleteRoleDelegation):
		if err := r.repo.Delete(ctx, cpsAction.UniqueId); err != nil {
			return nil, err
		}
	default:
		log.Errorf("[RoleDelegation Service][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	cpsAction.CurrentAction = roleDelegation
	return cpsAction, nil
}

// Create implements [service.RoleDelegationService].
func (r *roleDelegation) CreateWithExistingUser(ctx context.Context, roleDelegation imodel.RoleDelegation) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[RoleDelegation/Create] Creating role delegation for user %s with job title %s", roleDelegation.DelegatedUserID, roleDelegation.DelegatedUserJobTitle)
	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		log.Errorf("[Role Service][Create] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	if roleDelegation.DelegatedUserUserType == "BPS" {

		user, err := r.bpsUserRepo.GetByUserCode(ctx, roleDelegation.DelegatedUserUserCode)
		if err != nil {
			r.logger.Errorf("[RoleDelegation/Create] Failed to find BPS user: %v", err)
			return errors.New(localization.ErrorUserNotFoundOrDisabled.Code)

		}
		if user == nil || !user.Enabled {
			r.logger.Warnf("[RoleDelegation/Create] BPS User not found: %s", roleDelegation.DelegatedUserID)
			return errors.New(localization.ErrorUserNotFoundOrDisabled.Code)
		}
		roleDelegation.DelegatedUserFullName = user.FullName
		roleDelegation.DelegatedUserUserCode = user.UserCode
		roleDelegation.DelegatedUserJobTitle = user.JobTitle
		roleDelegation.DelegatedUserExistingRole = user.Role
		roleDelegation.DelegatedUserPhoneNumber = user.PhoneNumber
		roleDelegation.DelegatedUserEmail = user.Email
	} else {
		// Validate user existence
		user, err := r.cpsUserRepo.GetByUserCode(ctx, roleDelegation.DelegatedUserUserCode)
		if err != nil {
			r.logger.Errorf("[RoleDelegation/Create] Failed to find user: %v", err)
			return err
		}

		if !user.Enabled {
			r.logger.Warnf("[RoleDelegation/Create] User not found: %s", roleDelegation.DelegatedUserID)
			return errors.New(localization.ErrorUserNotFoundOrDisabled.Code)
		}
		roleDelegation.DelegatedUserFullName = user.FullName
		roleDelegation.DelegatedUserUserCode = user.UserCode
		roleDelegation.DelegatedUserJobTitle = user.JobTitle
		roleDelegation.DelegatedUserExistingRole = user.Role
		roleDelegation.DelegatedUserPhoneNumber = user.PhoneNumber
		roleDelegation.DelegatedUserEmail = user.Email
	}

	role, err := r.roleRepo.FindByCode(ctx, roleDelegation.NewRoleID)
	if err != nil {
		r.logger.Errorf("[RoleDelegation/Create] Failed to find job role: %s err:%+v", roleDelegation.NewRoleID, err)
		return localization.ErrRoleDelegationForUserNotFound
	}
	if role == nil || !role.Enable {
		r.logger.Errorf("[RoleDelegation/Create] Job role not found: %s", roleDelegation.NewRoleID)
		return errors.New(localization.ErrorRoleNotFound.Code)
	}

	if roleDelegation.DelegationType == "BPS" {
		if branch, err := r.branchRepo.GetBranchByDAOCode(ctx, roleDelegation.NewDepartmentOrBranch); err != nil || branch == nil {
			r.logger.Errorf("[RoleDelegation/Create] Branch not found: %s", roleDelegation.NewDepartmentOrBranch)
			return localization.ErrorInvalidDelegationBranch
		}
	} else {
		department, err := r.department.FindByID(ctx, roleDelegation.NewDepartmentOrBranch)
		if err != nil || department == nil || !department.Enabled {
			r.logger.Errorf("[RoleDelegation/Create] Department not found: %s err: %v", roleDelegation.DelegatedUserDepartmentOrBranch, err)
			return localization.ErrorInvalidDelegationDepartment
		}
	}
	roleDelegation.Enable = true
	cpsModel := lib.CpsModelBuilder(roleDelegation.DelegatedUserUserCode, maker, nil, roleDelegation, constants.RequestCreateRoleDelegationForExistingUser, constants.CREATE)
	return r.cpsService.CreateCPSAction(ctx, &cpsModel)

}

// Create implements [service.RoleDelegationService].
func (r *roleDelegation) CreateWithNewUser(ctx context.Context, roleDelegation imodel.RoleDelegation) error {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	log.Infof("[RoleDelegation/Create] Creating role delegation for user %s with job title %s", roleDelegation.DelegatedUserUserCode, roleDelegation.DelegatedUserJobTitle)

	maker := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(maker) {
		log.Errorf("[Role Service][Create] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	// handle delegated user validation based on user type and also populate any missing delegated user info from the request if needed for later use in delegation retrieval
	if roleDelegation.DelegatedUserUserType == "BPS" {
		user, err := r.bpsUserRepo.FindByOr(ctx, roleDelegation.DelegatedUserPhoneNumber, roleDelegation.DelegatedUserEmail, roleDelegation.DelegatedUserID)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			r.logger.Errorf("[RoleDelegation/Create] Failed to find BPS user: %v", err)
			return localization.ErrorUnexpectedError
		}
		if user != nil {
			r.logger.Warnf("[RoleDelegation/Create] BPS User not found: %s", roleDelegation.DelegatedUserID)
			return errors.New(localization.ErrorBpsUserAlreadyExists.Code)
		}
		// validate delegated user's branch existence for bps user type
		if branch, err := r.branchRepo.GetBranchByDAOCode(ctx, roleDelegation.DelegatedUserDepartmentOrBranch); err != nil || branch == nil {
			r.logger.Errorf("[RoleDelegation/Create] Branch not found: %s", roleDelegation.DelegatedUserDepartmentOrBranch)
			return localization.ErrorInvalidDelegationBranch
		}
	} else {
		// Validate user existence
		user, err := r.cpsUserRepo.FindByUsername(ctx, roleDelegation.DelegatedUserID)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			r.logger.Errorf("[RoleDelegation/Create] Failed to find user: %v", err)
			return localization.ErrorUnexpectedError
		}

		if user != nil {
			r.logger.Warnf("[RoleDelegation/Create] User not found: %s", roleDelegation.DelegatedUserID)
			return errors.New(localization.ErrorCpsUserAlreadyExists.Code)
		}
		// Validate user existence
		user, err = r.cpsUserRepo.FindByPhoneNumber(ctx, roleDelegation.DelegatedUserPhoneNumber)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			r.logger.Errorf("[RoleDelegation/Create] Failed to find user: %v", err)
			return localization.ErrorUnexpectedError
		}

		if user != nil {
			r.logger.Warnf("[RoleDelegation/Create] User not found: %s", roleDelegation.DelegatedUserID)
			return errors.New(localization.ErrorCpsUserAlreadyExists.Code)
		}
		// Validate user existence
		user, err = r.cpsUserRepo.FindByEmail(ctx, roleDelegation.DelegatedUserEmail)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			r.logger.Errorf("[RoleDelegation/Create] Failed to find user: %v", err)
			return localization.ErrorUnexpectedError
		}

		if user != nil {
			r.logger.Warnf("[RoleDelegation/Create] User not found: %s", roleDelegation.DelegatedUserID)
			return errors.New(localization.ErrorCpsUserAlreadyExists.Code)
		}
		// validate delegated user's department existence for cps user type
		existingDep, err := r.department.FindByID(ctx, roleDelegation.DelegatedUserDepartmentOrBranch)
		if err != nil || existingDep == nil || !existingDep.Enabled {
			r.logger.Errorf("[RoleDelegation/Create] Department not found: %s err: %v", roleDelegation.DelegatedUserDepartmentOrBranch, err)
			return localization.ErrorInvalidDelegationDepartment
		}
		roleDelegation.DelegatedUserDepartmentOrBranch = existingDep.ID.Hex()
	}
	// handle new role validation
	role, err := r.roleRepo.FindByCode(ctx, roleDelegation.NewRoleID)
	if err != nil {
		r.logger.Errorf("[RoleDelegation/Create] Failed to find job role: %s err: %+v", roleDelegation.NewRoleID, err)
		return localization.ErrRoleDelegationForUserNotFound
	}

	if role == nil || !role.Enable {
		r.logger.Errorf("[RoleDelegation/Create] Job role not found: %s", roleDelegation.NewRoleID)
		return errors.New(localization.ErrorRoleNotFound.Code)
	}

	// handle new department/branch validation based on delegation type
	if roleDelegation.DelegationType == "BPS" {
		if branch, err := r.branchRepo.GetBranchByDAOCode(ctx, roleDelegation.NewDepartmentOrBranch); err != nil || branch == nil {
			r.logger.Errorf("[RoleDelegation/Create] Branch not found: %s", roleDelegation.NewDepartmentOrBranch)
			return localization.ErrorInvalidDelegationBranch
		}
	} else {
		department, err := r.department.FindByID(ctx, roleDelegation.NewDepartmentOrBranch)
		if err != nil || department == nil || !department.Enabled {
			r.logger.Errorf("[RoleDelegation/Create] Department not found: %s err: %v", roleDelegation.NewDepartmentOrBranch, err)
			return localization.ErrorInvalidDelegationDepartment
		}
		roleDelegation.NewDepartmentOrBranch = department.ID.Hex()
	}

	roleDelegation.Enable = true

	cpsModel := lib.CpsModelBuilder(roleDelegation.DelegatedUserUserCode, maker, nil, roleDelegation, constants.RequestCreateRoleDelegationForNewUser, constants.CREATE)
	return r.cpsService.CreateCPSAction(ctx, &cpsModel)

}

// Delete implements [service.RoleDelegationService].
func (r *roleDelegation) Delete(ctx context.Context, id string) error {
	maker := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, r.logger)

	if local_util.IsIncomplete(maker) {
		log.Errorf("[RoleDelegation/Delete] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := r.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[RoleDelegation/Delete] failed to find existing role delegation: %v", err)
		return err
	}

	cpsActionData := lib.CpsModelBuilder(existing.DelegatedUserUserCode, maker, existing, nil, constants.RequestDeleteRoleDelegation, constants.DELETE)
	return r.cpsService.CreateCPSAction(ctx, &cpsActionData)
}

// EnableOrDisable implements [service.RoleDelegationService].
func (r *roleDelegation) EnableOrDisable(ctx context.Context, id string, enable bool) error {
	maker := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, r.logger)

	if local_util.IsIncomplete(maker) {
		log.Errorf("[RoleDelegation/EnableOrDisable] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := r.repo.FindByID(ctx, id)
	if err != nil {
		log.Errorf("[RoleDelegation/EnableOrDisable] failed to find existing role delegation: %v", err)
		return err
	}

	if enable && existing.Enable {
		log.Warnf("[RoleDelegation/EnableOrDisable] role delegation %s is already enabled", id)
		return errors.New(localization.ErrorRoleDelegationAlreadyEnabled.Code)
	}

	if !enable && !existing.Enable {
		log.Warnf("[RoleDelegation/EnableOrDisable] role delegation %s is already disabled", id)
		return errors.New(localization.ErrorRoleDelegationAlreadyDisabled.Code)
	}

	updated := *existing
	updated.UpdatedAt = time.Now()

	var requestType string
	if enable {
		requestType = constants.RequestEnableRoleDelegation
	} else {
		requestType = constants.RequestDisableRoleDelegation
	}

	cpsActionData := lib.CpsModelBuilder(existing.DelegatedUserUserCode, maker, existing, updated, requestType, constants.UPDATE)
	if err := r.cpsService.CreateCPSAction(ctx, &cpsActionData); err != nil {
		log.Errorf("[RoleDelegation/EnableOrDisable] failed to create CPS action: %v", err)
		return err
	}

	return nil
}

// Export implements [service.RoleDelegationService].
func (r *roleDelegation) Export(ctx context.Context, startDate, endDate time.Time, fileType string) (string, error) {
	log := local_util.LoggerFromCtx(ctx, r.logger)

	if r.minioClient == nil {
		log.Errorf("[RoleDelegation/Export] minio client is not configured")
		return "", errors.New(localization.RoleDelegationDataExportedError.Code)
	}

	fileType = strings.ToLower(strings.TrimSpace(fileType))
	if fileType != string(lib.FileTypeCSV) && fileType != string(lib.FileTypePDF) {
		log.Warnf("[RoleDelegation/Export] unsupported file type: %s", fileType)
		return "", errors.New(localization.ErrorInvalidRequest.Code)
	}

	data, err := r.repo.FindForExport(ctx, startDate, endDate)
	if err != nil {
		log.Errorf("[RoleDelegation/Export] fetch export data failed: %v", err)
		return "", err
	}

	if len(data) == 0 {
		return "", errors.New(localization.RoleDelegationDataNotFoundInDateRange.Code)
	}

	headers := []string{
		"Delegated User ID",
		"Delegated User Code",
		"Delegated User Name",
		"Delegated User Type",
		"Delegated User Job Title",
		"Delegated User Existing Role",
		"Delegation Type",
		"Delegator User ID",
		"Delegator User Name",
		"Delegator User Role",
		"New Role",
		"New Department Or Branch",
		"Enabled",
		"Start At",
		"End At",
		"Reason",
	}

	ext := "csv"
	if fileType == string(lib.FileTypePDF) {
		ext = "pdf"
	}

	objectName := fmt.Sprintf(
		"role_delegations_%s_to_%s_%d.%s",
		startDate.Format("20060102"),
		endDate.Format("20060102"),
		time.Now().Unix(),
		ext,
	)

	if fileType == string(lib.FileTypePDF) {
		rows := make([][]string, 0, len(data))
		for _, item := range data {
			rows = append(rows, roleDelegationExportRow(item))
		}
		url, exportErr := lib.ExportPDFAndUpload(ctx, r.minioClient, r.bucketName, r.cfg, objectName, headers, rows, lib.PDFExportOptions{PageSize: "A4"}, nil, r.logger)
		if exportErr != nil {
			log.Errorf("[RoleDelegation/Export] PDF export failed: %v", exportErr)
			return "", errors.New(localization.RoleDelegationDataExportedError.Code)
		}
		return url, nil
	}

	url, exportErr := lib.ExportCSVAndUpload(ctx, r.minioClient, r.bucketName, r.cfg, objectName, headers, func(writer *csv.Writer) error {
		for _, item := range data {
			if err := writer.Write(roleDelegationExportRow(item)); err != nil {
				return err
			}
		}
		return nil
	}, r.logger)
	if exportErr != nil {
		log.Errorf("[RoleDelegation/Export] CSV export failed: %v", exportErr)
		return "", errors.New(localization.RoleDelegationDataExportedError.Code)
	}

	return url, nil
}

// FindAll implements [service.RoleDelegationService].
func (r *roleDelegation) FindAll(ctx context.Context) (*[]imodel.RoleDelegation, error) {
	return r.repo.FindAll(ctx)
}

// FindAllWithPagination implements [service.RoleDelegationService].
func (r *roleDelegation) FindAllWithPagination(ctx context.Context, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.RoleDelegation], error) {
	return r.repo.FindAllWithPagination(ctx, filterParam)
}

// FindById implements [service.RoleDelegationService].
func (r *roleDelegation) FindById(ctx context.Context, id string) (*imodel.RoleDelegation, error) {
	roleDelegation, err := r.repo.FindByID(ctx, id)
	if err != nil {
		r.logger.Errorf("[RoleDelegation/FindById] failed to find role delegation by id: %v", err)
		return nil, err
	}
	if roleDelegation == nil {
		r.logger.Warnf("[RoleDelegation/FindById] role delegation not found by id: %s", id)
		return nil, errors.New(localization.ErrorResourceNotFound.Code)
	}

	// Populate existing delegated user's department/branch name.
	if roleDelegation.DelegatedUserUserType == "BPS" {
		if branch, err := r.branchRepo.GetBranchByCode(ctx, roleDelegation.DelegatedUserDepartmentOrBranch); err == nil && branch != nil {
			roleDelegation.DelegatedUserDepartmentOrBranchName = branch.Name
		} else if err != nil {
			r.logger.Errorf("[RoleDelegation/FindById] failed to find delegated user branch by code: %v", err)
		}
	} else {
		department, err := r.department.FindByID(ctx, roleDelegation.DelegatedUserDepartmentOrBranch)
		if err != nil {
			r.logger.Errorf("[RoleDelegation/FindById] failed to find delegated user department by id: %v", err)
			return nil, err
		}
		if department == nil {
			r.logger.Errorf("[RoleDelegation/FindById] delegated user department not found for id: %s", roleDelegation.DelegatedUserDepartmentOrBranch)
			return nil, localization.ErrorInvalidDelegationDepartment
		}
		roleDelegation.DelegatedUserDepartmentOrBranchName = department.Department
	}

	// for the new passed data
	if roleDelegation.DelegationType == "BPS" {
		branch, err := r.branchRepo.GetBranchByDAOCode(ctx, roleDelegation.NewDepartmentOrBranch)
		r.logger.Infof("[RoleDelegation/FindById] found branch for code %s: %v", roleDelegation.NewDepartmentOrBranch, branch)
		if err != nil || branch == nil {
			r.logger.Errorf("[RoleDelegation/FindById] failed to find branch by code: %v", err)
			if err != nil {
				return nil, err
			}
			return nil, localization.ErrorInvalidDelegationBranch
		}
		roleDelegation.NewDepartmentOrBranchName = branch.Name

	} else {
		department, err := r.department.FindByID(ctx, roleDelegation.NewDepartmentOrBranch)
		r.logger.Infof("[RoleDelegation/FindById] found department for id %s: %v", roleDelegation.NewDepartmentOrBranch, department)
		if err != nil {
			r.logger.Errorf("[RoleDelegation/FindById] failed to find new department by id: %v", err)
			return nil, err
		}
		if department == nil {
			r.logger.Errorf("[RoleDelegation/FindById] new department not found for id: %s", roleDelegation.NewDepartmentOrBranch)
			return nil, localization.ErrorInvalidDelegationDepartment
		}
		roleDelegation.NewDepartmentOrBranchName = department.Department
	}

	if role, err := r.roleRepo.FindByCode(ctx, roleDelegation.DelegatedUserExistingRole); err == nil && role != nil {
		roleDelegation.DelegatedUserExistingRoleName = role.Name
	} else {
		r.logger.Errorf("[RoleDelegation/FindById] failed to find delegated user's existing role by code: %v", err)
	}

	if role, err := r.roleRepo.FindByCode(ctx, roleDelegation.NewRoleID); err == nil && role != nil {
		roleDelegation.NewRoleIDName = role.Name
	} else {
		r.logger.Errorf("[RoleDelegation/FindById] failed to find new role by code: %v", err)
	}

	// delegator info
	if role, err := r.roleRepo.FindByCode(ctx, roleDelegation.DelegatorUserRole); err == nil && role != nil {
		roleDelegation.DelegatorUserRoleName = role.Name
	}
	return roleDelegation, err
}

// FindByUsername implements [service.RoleDelegationService].
func (r *roleDelegation) FindByUsername(ctx context.Context, id string, filterParam types.Filter) (*types.PaginatedResponse[[]imodel.RoleDelegation], error) {
	return r.repo.FindByUsername(ctx, id, filterParam)
}

// Update implements [service.RoleDelegationService].
func (r *roleDelegation) Update(ctx context.Context, id string, update imodel.RoleDelegation) error {
	maker := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, r.logger)

	if local_util.IsIncomplete(maker) {
		log.Errorf("[RoleDelegation/Update] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prev, err := r.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	newRole := *prev

	if !update.StartAt.IsZero() {
		newRole.StartAt = update.StartAt
	}
	if !update.EndAt.IsZero() {
		newRole.EndAt = update.EndAt
	}

	newRole.UpdatedAt = time.Now()

	cpsModel := lib.CpsModelBuilder(prev.DelegatedUserUserCode, maker, prev, newRole, constants.RequestUpdateRoleDelegation, constants.UPDATE)
	return r.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func roleDelegationExportRow(item imodel.RoleDelegation) []string {
	return []string{
		item.DelegatedUserID,
		item.DelegatedUserUserCode,
		item.DelegatedUserFullName,
		item.DelegatedUserUserType,
		item.DelegatedUserJobTitle,
		roleDelegationPreferredValue(item.DelegatedUserExistingRoleName, item.DelegatedUserExistingRole),
		item.DelegationType,
		item.DelegatorUserID,
		item.DelegatorUserFullName,
		roleDelegationPreferredValue(item.DelegatorUserRoleName, item.DelegatorUserRole),
		roleDelegationPreferredValue(item.NewRoleIDName, item.NewRoleID),
		roleDelegationPreferredValue(item.NewDepartmentOrBranchName, item.NewDepartmentOrBranch),
		fmt.Sprintf("%t", item.Enable),
		local_util.FormatTime(item.StartAt),
		local_util.FormatTime(item.EndAt),
		item.Reason,
	}
}

func roleDelegationPreferredValue(preferred, fallback string) string {
	if strings.TrimSpace(preferred) != "" {
		return preferred
	}
	return fallback
}

func roleDelegationEmailBody(roleDelegation imodel.RoleDelegation) string {
	recipientName := roleDelegationPreferredValue(roleDelegation.DelegatedUserFullName, roleDelegation.DelegatedUserID)
	oldRole := strings.TrimSpace(roleDelegationPreferredValue(roleDelegation.DelegatedUserExistingRoleName, roleDelegation.DelegatedUserExistingRole))
	newRole := strings.TrimSpace(roleDelegationPreferredValue(roleDelegation.NewRoleIDName, roleDelegation.NewRoleID))
	delegatorName := strings.TrimSpace(roleDelegationPreferredValue(roleDelegation.DelegatorUserFullName, roleDelegation.DelegatorUserID))
	delegatorEmail := "Unavailable"
	if strings.Contains(roleDelegation.DelegatorUserID, "@") {
		delegatorEmail = strings.TrimSpace(roleDelegation.DelegatorUserID)
	}

	destinationLabel := "Department"
	if strings.EqualFold(strings.TrimSpace(roleDelegation.DelegationType), "BPS") {
		destinationLabel = "Branch"
	}
	destinationValue := strings.TrimSpace(roleDelegationPreferredValue(roleDelegation.NewDepartmentOrBranchName, roleDelegation.NewDepartmentOrBranch))

	if oldRole == "" {
		oldRole = "N/A"
	}
	if newRole == "" {
		newRole = "N/A"
	}
	if delegatorName == "" {
		delegatorName = "N/A"
	}
	if destinationValue == "" {
		destinationValue = "N/A"
	}

	startAt := local_util.FormatTime(roleDelegation.StartAt)
	endAt := local_util.FormatTime(roleDelegation.EndAt)

	return fmt.Sprintf(
		"Dear %s,\n\nPlease be informed that your role delegation has been changed from %s to %s, effective from %s to %s. This notification is issued by CPS Portal for your records.\n\nDelegator Information:\nName: %s\nEmail: %s\n\nNew %s: %s\n\nSincerely,\nCPS Portal",
		recipientName,
		oldRole,
		newRole,
		startAt,
		endAt,
		delegatorName,
		delegatorEmail,
		destinationLabel,
		destinationValue,
	)
}

func NewRoleDelegationService(repo storage.RoleDelegationRepository, jobTitleRepo storage.JobRoleRepository, cpsUserRepo storage.CpsUserRepository, bpsUserRepo storage.BPSUserRepository, department storage.DepartmentRepository, roleRepo storage.RoleRepository, branchRepo storage.AccountBlockRepository, cpsService service.CPSActionService, minioClient *s3.Client, bucketName string, cfg config.VaultConfig, kafkaProducer *kafka.NotificationProducer, logger utils.Logger) service.RoleDelegationService {
	return &roleDelegation{
		cpsService:   cpsService,
		repo:         repo,
		jobTitleRepo: jobTitleRepo,
		roleRepo:     roleRepo,
		cpsUserRepo:  cpsUserRepo,
		bpsUserRepo:  bpsUserRepo,
		branchRepo:   branchRepo,
		kafkaClient:  kafkaProducer,
		department:   department,
		minioClient:  minioClient,
		bucketName:   bucketName,
		cfg:          cfg,
		logger:       logger,
	}
}
