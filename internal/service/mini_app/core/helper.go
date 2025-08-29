package core

import (
	"cbe-super-app-cps-action/internal/constants"
	miniappdto "cbe-super-app-cps-action/internal/constants/dto/mini_app"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	cRand "crypto/rand"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
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

	credentials := types.CredentialInformation{
		ID:            bson.NewObjectID(),
		Environment:   constants.UatEnvironment,
		MerchantAppID: req.Credential.MerchantAppID,
		FabricAppID:   req.Credential.FabricAppID,
		ShortCode:     req.Credential.ShortCode,
		AppSecret:     req.Credential.AppSecret,
		PrivateKey:    req.Credential.PrivateKey,
		PublicKey:     req.Credential.PublicKey,
	}

	miniApp := model.MiniApp{
		ID:                  bson.NewObjectID(),
		AppName:             req.AppName,
		CommissionGLAccount: req.CommissionGLAccount,
		AppType:             req.AppType,
		MerchantID:          req.MerchantID,
		ProductCode:         productCodes,
		Credential:          credentials,
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
		return nil
	}

	_, err := merchantService.FindByID(ctx, miniApp.MerchantID)
	if err != nil {
		log.Println("Failed to get merchant details", "merchantID", miniApp.MerchantID, "error", err)
		if err.Error() == "No Mini App merchant with these merchant!" {
			return errors.New(localization.ErrorMerchantNotFound.Code)
		}
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	return nil
}
func HandleCPSAction(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
	userData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(userData); incomplet {
		log.Println("User data incomplete for CPS action", "userCode", userData.UserCode)
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	cpsAction := lib.CpsModelBuilder(uniqueID, userData, curData, prevData, string(requestAction), string(constants.ActionDelete))

	log.Println("Creating CPS action", "userCode", userData.UserCode, "uniqueID", uniqueID, "actionType", actionType)
	if err := cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Println("Failed to create CPS action", "error", err)
		return err
	}
	log.Println("Successfully created CPS action", "userCode", userData.UserCode, "uniqueID", uniqueID)
	return nil
}

func GeneratePrefixedName(prefix, value string, logger shared_utils.Logger) (string, error) {
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
