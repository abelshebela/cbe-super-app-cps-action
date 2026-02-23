package bulk

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service/bulk/core"
	"errors"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"time"

	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"fmt"

	"math/rand"

	"cbe-super-app-cps-action/internal/constants/lib"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type bulkService struct {
	cpsActionRepo              service.CPSActionService
	repo                       storage.BulkServiceRepository
	accessListSegmentationRepo storage.AccessListSegmentationRepository
	logger                     utils.Logger
}

func NewBulkService(repo storage.BulkServiceRepository, CpsActionRepo service.CPSActionService, accessListSegmentationRepo storage.AccessListSegmentationRepository, logger utils.Logger) service.BulkService {
	return &bulkService{
		cpsActionRepo:              CpsActionRepo,
		repo:                       repo,
		accessListSegmentationRepo: accessListSegmentationRepo,
		logger:                     logger,
	}
}

func (s *bulkService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Bulk Service", "Authorize")
	defer span.End()

	s.logger.Infof("[BulkSvc][Authorize] action: %s", cpsAction.RequestAction)
	Now := time.Now()
	cpsAction.MakerActionTime = Now
	cpsAction.LastModifiedAt = Now

	cur, err := local_util.JsonUnmarshal[[]model.APPAccessList](cpsAction.CurrentAction)
	if err != nil {
		s.logger.Errorf("[Authorize] failed to unmarshal current action: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	var keys []string
	for _, v := range *cur {
		keys = append(keys, v.Key)
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestBulkServiceEnable):
		if err := s.repo.Update(ctx, keys, true); err != nil {
			span.AddEvent("[Authorize] failed to enable bulk service", trace.WithAttributes(
				attribute.String("error", err.Error()),
			))
			s.logger.Errorf("[BulkSvc][Authorize] enable err: %v", err)
			return nil, err
		}
		s.logger.Infof("[BulkSvc][Authorize] enabled %d", len(keys))
	case string(constants.RequestBulkServiceDisable):
		if err := s.repo.Update(ctx, keys, false); err != nil {
			span.AddEvent("[Authorize] failed to diable bulk services", trace.WithAttributes(
				attribute.String("error", err.Error()),
			))
			s.logger.Errorf("[BulkSvc][Authorize] disable err: %v", err)
			return nil, err
		}
		s.logger.Infof("[BulkSvc][Authorize] disabled %d", len(keys))
	default:
		s.logger.Errorf("[BulkSvc][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorInvalidRequiredAction.Code)
	}

	return cpsAction, nil
}

func GetAllKeysFromMaps(maps map[string]bool) []string {
	var keys []string
	// keySet := make(map[string]struct{}) // to avoid duplicates

	for key, _ := range maps {
		keys = append(keys, key)
	}
	return keys
}

// this function change the accesslist to key and enable map user to check the validity of  key given and check the current state
func storedAccessListToMAP(AccessLists []model.APPAccessList) map[string]bool {
	validKeys := make(map[string]bool)
	for _, access := range AccessLists {
		validKeys[access.Key] = access.Enabled
		for _, sub_access := range access.SubAccessList {
			validKeys[sub_access.Key] = sub_access.Enabled
		}
	}

	return validKeys
}

func validaterAccessKey(validAccessMap map[string]bool, accessList []string, flag bool) ([]string, map[string]bool, bool) {
	var invalidKeys []string
	validKeys := make(map[string]bool)

	isActionValid := true
	for _, access := range accessList {
		if enabled, exists := validAccessMap[access]; exists {
			if enabled == flag {
				isActionValid = false
			}
			validKeys[access] = enabled
		} else {
			invalidKeys = append(invalidKeys, access)
		}
		// for _, subAccess := range access.SubAccessList {
		// 	if enabled, exists := validAccessMap[subAccess.Key]; exists {
		// 		validKeys[subAccess.Key] = enabled
		// 	} else {
		// 		invalidKeys = append(invalidKeys, subAccess.Key)
		// 	}
		// }
	}
	return invalidKeys, validKeys, isActionValid
}

func (s *bulkService) GetAllBulkServices(ctx context.Context, filterParams *types.Filter) ([]model.APPAccessList, []model.APPAccessList, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "GetAllBulkServices", "Bulk Service", "GetAllBulkServices")
	defer span.End()

	result, err := s.repo.FindAll(ctx)
	if err != nil {
		span.AddEvent("[GetAllBulkServices] failed to fetch bulk services", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		s.logger.Errorf("[GetAllBulkServices] failed to fetch bulk services: %v", err)
		return nil, nil, err
	}

	relation, err := s.accessListSegmentationRepo.FindParentChildRelationship(ctx)
	if err != nil {
		span.AddEvent("[GetAllBulkServices] failed to fetch access list segmentation relationships", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		s.logger.Errorf("[GetAllBulkServices] failed to fetch access list segmentation relationships: %v", err)
		return nil, nil, err
	}

	// s.logger.Infof("[GetAllBulkServices] successfully fetched %d bulk services and %d parent-child relationships", len(result.Data), len(relation))

	enabled, disabled := core.SplitEnabledDisabledTree(result)
	// s.logger.Infof("[GetAllBulkServices] split bulk services into %d enabled and %d disabled", len(enabled), len(disabled))

	enabled = core.MapParentChildRelationship(relation, enabled)
	// s.logger.Infof("[GetAllBulkServices] successfully mapped parent-child relationships for %d enabled bulk services", len(enabled))

	disabled = core.MapParentChildRelationship(relation, disabled)
	// s.logger.Infof("[GetAllBulkServices] successfully mapped parent-child relationships for %d disabled bulk services", len(disabled))
	return enabled, disabled, nil
}

func (s *bulkService) CheckServiceIsEnabledOrDisabled(ctx context.Context, keys []string, isEnabled bool) ([]string, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "CheckServiceIsEnabledOrDisabled", "Bulk Service", "CheckServiceIsEnabledOrDisabled")
	defer span.End()

	allAccessLists, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	validAccessList := storedAccessListToMAP(allAccessLists)
	invalidDatas, validKeys, _ := validaterAccessKey(validAccessList, keys, isEnabled)

	if len(invalidDatas) > 0 {
		span.AddEvent("[CheckServiceIsEnabledOrDisabled] One or more services are not found with these keys")
		s.logger.Errorf("[BulkSvc][CheckState] not found: %v", invalidDatas)
		return nil, errors.New(localization.ErrorInvalidBulkServiceKey.Code)
	}

	for key, value := range validKeys {
		if isEnabled {

			if value {
				span.AddEvent("[CheckServiceIsEnabledOrDisabled] Service already enabled")
				s.logger.Errorf("[BulkSvc][CheckState] already enabled: %s", key)
				return nil, errors.New(localization.ErrorBulkServiceAlreadyEnabled.Code)
			}
		} else {

			if !value {
				span.AddEvent("[CheckServiceIsEnabledOrDisabled] Service already disabled")
				s.logger.Errorf("[BulkSvc][CheckState] already disabled: %s", key)
				return nil, errors.New(localization.ErrorBulkServiceAlreadyDisabled.Code)
			}
		}

	}
	VK := GetAllKeysFromMaps(validKeys)

	return VK, nil
}

func (s *bulkService) EnableBulkService(ctx context.Context, keys []string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableBulkService", "Bulk Service", "EnableBulkService")
	defer span.End()

	if len(keys) == 0 {
		span.AddEvent("[EnableBulkService] No keys provided to bulk enable")
		s.logger.Errorf("[BulkSvc][Enable] no keys")
		return errors.New(localization.ErrorKeyRequiredForBulkService.Code)
	}

	allAccessLists, err := s.repo.FindAll(ctx)
	if err != nil {
		return err
	}
	// Initialize the map to avoid nil map panic
	allKeys := make(map[string]model.APPAccessList)
	for _, access := range allAccessLists {
		s.logger.Infof("[DisableBulkService] Access List - Key: %s, Enabled: %t", access.Key, access.Enabled)
		allKeys[access.Key] = access
	}

	var disabledKeys []model.APPAccessList
	var noneDisabledKeys []model.APPAccessList

	for _, key := range keys {
		if access, exists := allKeys[key]; exists {
			if access.Enabled {
				span.AddEvent("[DisableBulkService] service already enabled")
				s.logger.Errorf("[DisableBulkService] service already enabled: %s", key)
				return errors.New(localization.ErrorBulkServiceAlreadyEnabled.Code)
			}
			disabledKeys = append(disabledKeys, access)
			access.Enabled = true
			noneDisabledKeys = append(noneDisabledKeys, access)
		}
	}

	userPayload := local_util.ExtractUserFromContext(ctx)

	cpsAction := lib.CpsModelBuilder("", userPayload, disabledKeys, noneDisabledKeys, string(constants.RequestBulkServiceEnable), string(constants.UpdateAction))

	if err := s.cpsActionRepo.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("[EnableBulkService] Failed to create CPS action", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("[BulkSvc][Enable] cps action err: %v", err)
		return err
	}
	s.logger.Infof("[EnableBulkService] CPS action created successfully for %d services", len(keys))
	return nil
}

func (s *bulkService) DisableBulkService(ctx context.Context, keys []string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DisableBulkService", "Bulk Service", "DisableBulkService")
	defer span.End()

	if len(keys) == 0 {
		span.AddEvent("[DisableBulkService] no keys provided")
		s.logger.Errorf("[BulkSvc][Disable] no keys")
		return errors.New(localization.ErrorKeyRequiredForBulkService.Code)
	}

	allAccessLists, err := s.repo.FindAll(ctx)
	if err != nil {
		return err
	}
	// Initialize the map to avoid nil map panic
	allKeys := make(map[string]model.APPAccessList)
	for _, access := range allAccessLists {
		s.logger.Infof("[DisableBulkService] Access List - Key: %s, Enabled: %t", access.Key, access.Enabled)
		allKeys[access.Key] = access
	}

	var disabledKeys []model.APPAccessList
	var noneDisabledKeys []model.APPAccessList
	for _, key := range keys {
		if access, exists := allKeys[key]; exists {
			// accessList := access.(model.APPAccessList)
			if !access.Enabled {
				span.AddEvent("[DisableBulkService] service already disabled")
				s.logger.Errorf("[DisableBulkService] service already disabled: %s", key)
				return errors.New(localization.ErrorBulkServiceAlreadyDisabled.Code)
			}
			noneDisabledKeys = append(noneDisabledKeys, access)
			access.Enabled = false
			disabledKeys = append(disabledKeys, access)
		}
	}

	userPayload := local_util.ExtractUserFromContext(ctx)

	unicode := GenerateUnique14DigitCode()
	cpsAction := lib.CpsModelBuilder(unicode, userPayload, noneDisabledKeys, disabledKeys, string(constants.RequestBulkServiceDisable), string(constants.UpdateAction))

	if err := s.cpsActionRepo.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("[DisableBulkService] failed to create CPS action", trace.WithAttributes(attribute.String("error", err.Error())))
		s.logger.Errorf("[BulkSvc][Disable] cps action err: %v", err)
		return err
	}
	s.logger.Infof("[DisableBulkService] CPS action created successfully for %d services", len(keys))
	return nil
}

func GenerateUnique14DigitCode() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond) // 13 digits
	rand.Seed(time.Now().UnixNano())
	randomDigit := rand.Intn(10) // 1 digit
	return fmt.Sprintf("%013d%d", timestamp, randomDigit)
}
