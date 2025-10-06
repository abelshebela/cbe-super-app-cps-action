package core

import (
	"cbe-super-app-cps-action/internal/constants"
	miniappdto "cbe-super-app-cps-action/internal/constants/dto/mini_app"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	cRand "crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func BuildMiniAppFromRequest(req miniappdto.MiniAppCreateRequest, withTimestamps bool) model.MiniApp {
	now := time.Now()
	productCodes := make([]types.ProductCode, 0, len(req.ProductCode))
	for _, pc := range req.ProductCode {
		productCodes = append(productCodes, types.ProductCode{
			ID:             utils.RandomGenerator(20),
			BranchType:     constants.BranchType(pc.BranchType),
			ProductCode:    pc.ProductCode,
			VATCode:        pc.VATCode,
			ServiceFeeCode: pc.ServiceFeeCode,
		})
	}

	var miniAppID bson.ObjectID
	if req.ID != "" {
		miniAppID, _ = bson.ObjectIDFromHex(req.ID)
	} else {
		miniAppID = bson.NewObjectID()
	}
	miniApp := model.MiniApp{
		ID:                  miniAppID,
		AppName:             req.AppName,
		CommissionGLAccount: req.CommissionGLAccount,
		AppType:             req.AppType,
		MerchantID:          req.MerchantID,
		ProductCode:         productCodes,
		IsEventMiniApp:      req.IsEventMiniApp,
		IsThreeClick:        req.IsThreeClick,
		URL:                 req.URL,
		Stage:               req.Stage,
		AppViewType:         req.AppViewType,
	}

	if withTimestamps {
		miniApp.CreatedAt = now
		miniApp.LastModifiedAt = now
		miniApp.DeletedAt = time.Time{}
	} else {
		miniApp.LastModifiedAt = now
	}

	return miniApp
}

func SetMerchantDetails(ctx context.Context, merchantService service.MiniAppMerchantService, miniApp *miniappdto.MiniAppCreateRequest) error {
	if miniApp.MerchantID == "" {
		return errors.New(localization.ErrorMerchantIDRequired.Code)
	}

	merchant, err := merchantService.FindByID(ctx, miniApp.MerchantID)
	if err != nil {
		log.Println("Failed to get merchant details", "merchantID", miniApp.MerchantID, "error", err)
		return errors.New(localization.ErrorMerchantNotFound.Code)
	}

	
	if merchant.IsDeleted {
		log.Println("Merchant is deleted", "merchantID", miniApp.MerchantID)
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}

	if !merchant.Enabled {
		log.Println("Merchant is disabled", "merchantID", miniApp.MerchantID)
		return errors.New(localization.ErrorMiniAppMerchantDisableFailed.Code)
	}

	return nil
}
func HandleCPSAction(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
	userData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(userData); incomplet {
		log.Println("User data incomplete for CPS action", "userCode", userData.UserCode)
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	cpsAction := lib.CpsModelBuilder(uniqueID, userData, prevData, curData, string(requestAction), string(constants.ActionDelete))

	log.Println("Creating CPS action", "userCode", userData.UserCode, "uniqueID", uniqueID, "actionType", actionType)
	if err := cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Println("Failed to create CPS action", "error", err)
		return err
	}
	log.Println("Successfully created CPS action", "userCode", userData.UserCode, "uniqueID", uniqueID)
	return nil
}

func GeneratePrefixedName(prefix, value string, logger utils.Logger) (string, error) {
	logger.Infof("Generating prefixed name", "prefix", prefix, "value", value)

	if prefix == "" || value == "" {
		logger.Errorf("Invalid input for GeneratePrefixedName", "prefix", prefix, "value", value)
		return "", fmt.Errorf("prefix and value must not be empty")
	}

	digits := "0123456789"
	max := big.NewInt(int64(len(digits)))
	code := make([]byte, 7)

	for i := range code {
		n, err := cRand.Int(cRand.Reader, max)
		if err != nil {
			logger.Errorf("Failed to generate random digit", "error", err)
			return "", fmt.Errorf("failed to generate random digit: %v", err)
		}
		code[i] = digits[n.Int64()]
	}

	value = strings.ReplaceAll(value, " ", "")
	prefix = strings.ReplaceAll(prefix, " ", "")
	result := strings.Join([]string{prefix, value, string(code)}, "-")
	logger.Infof("Successfully generated prefixed name", "result", result)
	return result, nil
}



func ValidateParentMerchantForEnableEntity(merchant *model.MiniAppMerchant) error {
	if merchant == nil {
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}
	if merchant.IsDeleted {
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}
	if !merchant.Enabled {
		return errors.New(localization.ErrorMiniAppMerchantDisableFailed.Code)
	}
	return nil
}

func ValidateParentMerchantForOperationEntity(merchant *model.MiniAppMerchant) error {
	if merchant == nil {
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}
	if merchant.IsDeleted {
		return errors.New(localization.ErrorMiniAppMerchantNotFound.Code)
	}
	return nil
}

func ValidMiniAppChecker(ctx context.Context, miniAppRepo storage.MiniAppRepository, isCreate bool, id, miniAppName string) (bool, error) {
	existedTitleFilter := types.Filter{
		Filters: map[string]interface{}{
			"app_name": miniAppName,
		},
	}

	exists, err := miniAppRepo.FindAllWithPagination(ctx, existedTitleFilter)
	if err != nil && err.Error() != localization.ErrorFileNotFound.Code {
		return false, err
	}
	if len(exists.Data) > 0 {
		return false, nil
	}

	return true, nil
}
