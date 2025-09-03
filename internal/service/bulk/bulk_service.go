package bulk

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"errors"

	// "strconv"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"fmt"

	//   "fmt"
	"math/rand"
	// "time"

	// "time"

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
	fmt.Println("___________________________________________----")
	fmt.Println(" in Authorize")
	Now := time.Now()
	cpsAction.MakerActionTime = Now
	cpsAction.LastModifiedAt = Now
	fmt.Println("step 1-----------------------------")
	fmt.Printf("type: %T", cpsAction.CurrentAction)
	fmt.Println("-----------------------------")

	fmt.Printf("type: %v", cpsAction.CurrentAction)

	// updateData := cpsAction.CurrentAction.([]string)
	doc, ok := cpsAction.CurrentAction.(bson.D)
	if !ok {
		s.logger.Errorf("not a bson.D")
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
		s.logger.Errorf("not a bson.A")
		return nil, errors.New("not a bson.A")
	}

	var result []string
	for _, v := range arr {
		if s, ok := v.(string); ok {
			result = append(result, s)
		}
	}

	fmt.Println(result)
	// Output: [wallet wallettelebirr topup transfertodashen]

	allAccessLists, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	var flag bool
	switch cpsAction.RequestAction {
	case string(constants.RequestBulkServiceEnable):
		flag = true
	case string(constants.RequestBulkServiceDisable):
		flag = false
	}

	validAccessList := storedAccessListToMAP(allAccessLists)
	invalidDatas, validKeys, _ := validaterAccessKey(validAccessList, result, flag)

	if len(invalidDatas) > 0 {
		s.logger.Errorf("one or more Servive not found")
		return nil, fmt.Errorf(localization.ErrorInvalidBulkServiceKey.Code)
	}
	// if !isActionAllowed {
	// 	s.logger.Infof("you enter a key  satisfy the action")
	// 	return nil, fmt.Errorf(localization.ErrorBulkServiceActionNotSatisfied.Code)
	// }

	for key, value := range validKeys {

		switch strings.ToUpper(cpsAction.RequestAction) {
		case string(constants.RequestBulkServiceEnable):
			if value == true {
				s.logger.Infof("you enter already enabled service: %v", key)
				return nil, errors.New(localization.ErrorBulkServiceAlreadyEnabled.Code)
			}
		case string(constants.RequestBulkServiceDisable):
			if value == false {
				s.logger.Infof("you enter already disabled service: %v", key)
				return nil, errors.New(localization.ErrorBulkServiceAlreadyDisabled.Code)
			}

		}
	}
	keys := GetAllKeysFromMaps(validKeys)

	switch cpsAction.RequestAction {
	case string(constants.RequestBulkServiceEnable):
		return nil, s.repo.Update(ctx, keys, true)
	case string(constants.RequestBulkServiceDisable):
		return nil, s.repo.Update(ctx, keys, false)
	default:
		return nil, fmt.Errorf(localization.ErrorInvalidRequiredAction.Code)
	}

	// return s.repo.Update(ctx, &updateData)

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

	return s.repo.FindAllWithPagination(ctx, *filterParams)
}
func (s *bulkService) EnableBulkService(ctx context.Context, keys []string) (string, error) {
	if len(keys) == 0 {
		return "", errors.New(localization.ErrorKeyRequiredForBulkService.Code)
	}
	userPayload := local_util.ExtractUserFromContext(ctx)

	type currAction struct {
		Keys []string
	}
	cpsAction := lib.CpsModelBuilder("", userPayload, nil, currAction{
		Keys: keys,
	}, string(constants.RequestBulkServiceEnable), string(constants.UpdateAction))

	return "", s.cpsActionRepo.CreateCPSAction(ctx, &cpsAction)
}

func (s *bulkService) DisableBulkService(ctx context.Context, keys []string) (string, error) {
	if len(keys) == 0 {
		return "", errors.New(localization.ErrorKeyRequiredForBulkService.Code)
	}

	userPayload := local_util.ExtractUserFromContext(ctx)

	type currAction struct {
		Keys []string
	}
	unicode := GenerateUnique14DigitCode()
	cpsAction := lib.CpsModelBuilder(unicode, userPayload, nil, currAction{
		Keys: keys,
	}, string(constants.RequestBulkServiceDisable), string(constants.UpdateAction))

	err := s.cpsActionRepo.CreateCPSAction(ctx, &cpsAction)
	return "", err
}

func GenerateUnique14DigitCode() string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond) // 13 digits
	rand.Seed(time.Now().UnixNano())
	randomDigit := rand.Intn(10) // 1 digit
	return fmt.Sprintf("%013d%d", timestamp, randomDigit)
}
