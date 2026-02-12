package core

import (
	"cbe-super-app-cps-action/internal/constants"
	TopupDto "cbe-super-app-cps-action/internal/constants/dto/topup"
	"cbe-super-app-cps-action/internal/constants/lib"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"time"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	topupDto "cbe-super-app-cps-action/internal/constants/dto/topup"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	// shared_types "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"

	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func NonEmptyString(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

func GeneratePrefixedName(prefix, value string, logger shared_utils.Logger) (string, error) {
	logger.Infof("Generating prefixed name", "prefix", prefix, "value", value)

	if prefix == "" || value == "" {
		logger.Errorf("Invalid input for GeneratePrefixedName", "prefix", prefix, "value", value)
		return "", fmt.Errorf("prefix and value must not be empty")
	}

	value = strings.ToUpper(strings.ReplaceAll(value, " ", "_"))
	prefix = strings.ToUpper(strings.ReplaceAll(prefix, " ", "_"))
	if strings.HasPrefix(value, prefix+"-") {
		return value, nil
	}
	result := strings.Join([]string{prefix, value}, "-")
	logger.Infof("Successfully generated prefixed name", "result", result)
	return result, nil

}

func ExistingIdentifierForUpdate(existing model.Topup, id string, req topupDto.TopupRequest) error {
	if &existing == nil {
		return nil
	}

	normalizedName := strings.TrimSpace(req.Name)
	normalizedCode := strings.TrimSpace(req.Code)

	existingID := local_util.FirstHex24(existing.ID.String())
	if normalizedName != "" && existing.Name == normalizedName && existingID != id {
		return errors.New(localization.ErrorTopupNameAlreadyExists.Code)
	}

	if normalizedCode != "" && existing.Code == normalizedCode && existingID != id {
		return errors.New(localization.ErrorTopupCodeAlreadyExists.Code)
	}
	return nil
}

func ToCreateTopupDoc(name, code, URL string) *model.Topup {
	return &model.Topup{
		Name:   name,
		Code:   code,
		Avatar: URL,
		// Services: shared_types.Services{
		// 	Self:  self,
		// 	Other: other,
		// 	Agent: agent,
		// },
	}
}

// note: this comparision might not be needed if the existing data is first in the request form and the user update those values
func ToUpdateTopupDoc(existing model.Topup, req TopupDto.TopupRequest) (*model.Topup, int) {
	// var Topup model.Topup
	Topup := existing
	Topup.LastModifiedAt = time.Now()

	change_count := 0
	// if req.Agent == existing.Services.Agent {
	// 	Topup.Services.Agent = existing.Services.Agent
	// } else {
	// 	change_count++
	// 	Topup.Services.Agent = req.Agent
	// }
	// if req.Other == existing.Services.Other {
	// 	Topup.Services.Other = existing.Services.Other
	// } else {
	// 	change_count++
	// 	Topup.Services.Other = req.Other
	// }
	// if req.Self == existing.Services.Self {
	// 	Topup.Services.Self = existing.Services.Self
	// } else {
	// 	change_count++
	// 	Topup.Services.Self = req.Self
	// }
	if req.Name != "" {
		if req.Name == existing.Name {
			Topup.Name = existing.Name
		} else {
			change_count++
			Topup.Name = req.Name
		}
	}
	if req.Code != "" {
		if req.Code == existing.Code {
			Topup.Code = existing.Code
		} else {
			change_count++
			Topup.Code = req.Code
		}
	}
	Topup.Enabled = existing.Enabled
	Topup.Avatar = existing.Avatar
	return &Topup, change_count
}

func HandleCPSAction(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
	ctx, span := local_util.TraceLogger(ctx, "core", "HandleCPSAction", "core", "core")
	defer span.End()
	log.Println("Handling CPS action", "uniqueID", uniqueID, "requestAction", requestAction)
	userData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(userData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(attribute.String("userCode", userData.UserCode)))
		log.Println("User data incomplete for CPS action", "userCode", userData.UserCode)
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	cpsAction := lib.CpsModelBuilder(uniqueID, userData, prevData, curData, string(requestAction), string(actionType))

	if err := cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("CPS action creation failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("userCode", userData.UserCode)))
		log.Println("Failed to create CPS action", "error", err)
		return err
	}
	log.Println("Successfully created CPS action", "userCode", userData.UserCode, "uniqueID", uniqueID)
	return nil
}
