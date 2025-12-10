package bulk

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"errors"

	"time"

	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"fmt"

	"math/rand"

	"cbe-super-app-cps-action/internal/constants/lib"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type bulkService struct {
	cpsActionRepo service.CPSActionService
	repo          storage.BulkServiceRepository
	logger        utils.Logger
}

func NewBulkService(repo storage.BulkServiceRepository, CpsActionRepo service.CPSActionService, logger utils.Logger) service.BulkService {
	return &bulkService{
		cpsActionRepo: CpsActionRepo,
		repo:          repo,
		logger:        logger,
	}
}

func (s *bulkService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	s.logger.Infof("[Authorize] authorizing bulk service action: %s", cpsAction.RequestAction)
	Now := time.Now()
	cpsAction.MakerActionTime = Now
	cpsAction.LastModifiedAt = Now
	// updateData := cpsAction.CurrentAction.([]string)
	doc, ok := cpsAction.CurrentAction.(bson.D)
	if !ok {
		s.logger.Errorf("[Authorize] current action is not a bson.D")
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	// extract the "keys" array
	var arr bson.A
	for _, elem := range doc {
		if elem.Key == "keys" {
			arr, ok = elem.Value.(bson.A)
			break
		}
	}
	if !ok {
		s.logger.Errorf("[Authorize] keys array is not a bson.A")
		return nil, errors.New("not a bson.A")
	}

	var keys []string
	for _, v := range arr {
		if s, ok := v.(string); ok {
			keys = append(keys, s)
		}
	}

	switch cpsAction.RequestAction {
	case string(constants.RequestBulkServiceEnable):
		if err := s.repo.Update(ctx, keys, true); err != nil {
			s.logger.Errorf("[Authorize] failed to enable bulk services: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] successfully enabled %d bulk services", len(keys))
		return nil, nil
	case string(constants.RequestBulkServiceDisable):
		if err := s.repo.Update(ctx, keys, false); err != nil {
			s.logger.Errorf("[Authorize] failed to disable bulk services: %v", err)
			return nil, err
		}
		s.logger.Infof("[Authorize] successfully disabled %d bulk services", len(keys))
		return nil, nil
	default:
		s.logger.Errorf("[Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorInvalidRequiredAction.Code)
	}
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
func storedAccessListToMAP(AccessLists []*model.APPAccessList) map[string]bool {
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

func (s *bulkService) GetAllBulkServices(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]*model.APPAccessList], error) {
	result, err := s.repo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		s.logger.Errorf("[GetAllBulkServices] failed to fetch bulk services: %v", err)
		return nil, err
	}
	return result, nil
}

func (s *bulkService) CheckServiceIsEnabledOrDisabled(ctx context.Context, keys []string, isEnabled bool) ([]string, error) {
	allAccessLists, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	validAccessList := storedAccessListToMAP(allAccessLists)
	invalidDatas, validKeys, _ := validaterAccessKey(validAccessList, keys, isEnabled)

	if len(invalidDatas) > 0 {
		s.logger.Errorf("[CheckServiceIsEnabledOrDisabled] one or more services not found: %v", invalidDatas)
		return nil, errors.New(localization.ErrorInvalidBulkServiceKey.Code)
	}

	for key, value := range validKeys {
		if isEnabled {

			if value {
				s.logger.Errorf("[CheckServiceIsEnabledOrDisabled] service already enabled: %s", key)
				return nil, errors.New(localization.ErrorBulkServiceAlreadyEnabled.Code)
			}
		} else {

			if !value {
				s.logger.Errorf("[CheckServiceIsEnabledOrDisabled] service already disabled: %s", key)
				return nil, errors.New(localization.ErrorBulkServiceAlreadyDisabled.Code)
			}
		}

	}
	VK := GetAllKeysFromMaps(validKeys)

	return VK, nil
}

func (s *bulkService) EnableBulkService(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		s.logger.Errorf("[EnableBulkService] no keys provided")
		return errors.New(localization.ErrorKeyRequiredForBulkService.Code)
	}

	validKeys, err := s.CheckServiceIsEnabledOrDisabled(ctx, keys, true)
	if err != nil {
		s.logger.Errorf("[EnableBulkService] validation failed: %v", err)
		return err
	}

	userPayload := local_util.ExtractUserFromContext(ctx)

	type currAction struct {
		Keys []string
	}
	cpsAction := lib.CpsModelBuilder("", userPayload, nil, currAction{
		Keys: validKeys,
	}, string(constants.RequestBulkServiceEnable), string(constants.UpdateAction))

	if err := s.cpsActionRepo.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[EnableBulkService] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[EnableBulkService] CPS action created successfully for %d services", len(validKeys))
	return nil
}

func (s *bulkService) DisableBulkService(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		s.logger.Errorf("[DisableBulkService] no keys provided")
		return errors.New(localization.ErrorKeyRequiredForBulkService.Code)
	}

	validKeys, err := s.CheckServiceIsEnabledOrDisabled(ctx, keys, false)
	if err != nil {
		s.logger.Errorf("[DisableBulkService] validation failed: %v", err)
		return err
	}

	userPayload := local_util.ExtractUserFromContext(ctx)

	type currAction struct {
		Keys []string
	}
	unicode := GenerateUnique14DigitCode()
	cpsAction := lib.CpsModelBuilder(unicode, userPayload, nil, currAction{
		Keys: validKeys,
	}, string(constants.RequestBulkServiceDisable), string(constants.UpdateAction))

	if err := s.cpsActionRepo.CreateCPSAction(ctx, &cpsAction); err != nil {
		s.logger.Errorf("[DisableBulkService] failed to create CPS action: %v", err)
		return err
	}
	s.logger.Infof("[DisableBulkService] CPS action created successfully for %d services", len(validKeys))
	return nil
}

func GenerateUnique14DigitCode() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond) // 13 digits
	rand.Seed(time.Now().UnixNano())
	randomDigit := rand.Intn(10) // 1 digit
	return fmt.Sprintf("%013d%d", timestamp, randomDigit)
}
