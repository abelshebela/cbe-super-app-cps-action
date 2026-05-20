package core

import (
	"cbe-super-app-cps-action/internal/constants"
	walletDto "cbe-super-app-cps-action/internal/constants/dto/wallet"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/storage"
	"slices"
	"time"

	// "cbe-super-app-cps-action/internal/constants/types"
	local_model "cbe-super-app-cps-action/internal/constants/model"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/service"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	shared_utils "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

func NonEmptyString(s, fallback string) string {
	if s != "" {
		return s
	}
	return fallback
}

func GeneratePrefixedName(prefix, value string, logger shared_utils.Logger) (string, error) {
	logger.Infof("[WalletCore][GenPrefix] prefix: %s value: %s", prefix, value)

	if prefix == "" || value == "" {
		logger.Errorf("[WalletCore][GenPrefix] invalid input prefix: %s value: %s", prefix, value)
		return "", fmt.Errorf("prefix and value must not be empty")
	}

	value = strings.ToUpper(strings.ReplaceAll(value, " ", "_"))
	prefix = strings.ToUpper(strings.ReplaceAll(prefix, " ", "_"))
	result := strings.Join([]string{prefix, value}, "-")
	logger.Infof("[WalletCore][GenPrefix] result: %s", result)
	return result, nil

}

func ToCreateWalletDoc(req walletDto.WalletRequest, logo string, services []local_model.Service) *local_model.WalletOracle {
	// Extract service IDs from the services slice
	var selfServiceID, otherServiceID, agentServiceID string
	var selfServiceName, otherServiceName, agentServiceName string
	for _, service := range services {
		switch service.ID {
		case req.SelfServiceID:
			selfServiceID = service.ID
			selfServiceName = service.ServiceName
		case req.OtherServiceID:
			otherServiceID = service.ID
			otherServiceName = service.ServiceName
		case req.AgentServiceID:
			agentServiceID = service.ID
			agentServiceName = service.ServiceName
		}
	}

	return &local_model.WalletOracle{
		Name:             req.Name,
		UniqueCode:       strings.ToUpper(strings.TrimSpace(req.UniqueCode)),
		Avatar:           logo,
		Self:             *req.Self,
		Other:            *req.Other,
		Agent:            *req.Agent,
		SelfServiceID:    selfServiceID,
		OtherServiceID:   otherServiceID,
		AgentServiceID:   agentServiceID,
		SelfServiceName:  selfServiceName,
		OtherServiceName: otherServiceName,
		AgentServiceName: agentServiceName,
	}
}

// note: this comparision might not be needed if the existing data is first in the request form and the user update those values
func ToUpdateWalletDoc(existing local_model.WalletOracle, req local_model.WalletOracle) (local_model.WalletOracle, int) {
	wallet := existing
	changeCount := 0

	if req.Self != existing.Self {
		changeCount++
		wallet.Self = req.Self
	}

	if req.Other != existing.Other {
		changeCount++
		wallet.Other = req.Other
	}

	if req.Agent != existing.Agent {
		changeCount++
		wallet.Agent = req.Agent
	}

	if req.SelfServiceID != "" && req.SelfServiceID != existing.SelfServiceID {
		changeCount++
		wallet.SelfServiceID = req.SelfServiceID
		wallet.SelfServiceName = req.SelfServiceName
	}

	if req.OtherServiceID != "" && req.OtherServiceID != existing.OtherServiceID {
		changeCount++
		wallet.OtherServiceID = req.OtherServiceID
		wallet.OtherServiceName = req.OtherServiceName
	}

	if req.AgentServiceID != "" && req.AgentServiceID != existing.AgentServiceID {
		changeCount++
		wallet.AgentServiceID = req.AgentServiceID
		wallet.AgentServiceName = req.AgentServiceName
	}

	if req.Name != "" && req.Name != existing.Name {
		changeCount++
		wallet.Name = req.Name
	}

	if req.UniqueCode != "" && req.UniqueCode != existing.UniqueCode {
		changeCount++
		wallet.UniqueCode = req.UniqueCode
	}

	return wallet, changeCount
}

func HandleCPSAction(ctx context.Context, cpsService service.CPSActionService, uniqueID string, requestAction constants.RequestAction, curData, prevData interface{}, actionType constants.ActionType) error {
	userData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(userData); incomplet {
		log.Println("User data incomplete for CPS action", "userCode", userData.UserCode)
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	cpsAction := lib.CpsModelBuilder(uniqueID, userData, prevData, curData, string(requestAction), string(actionType))

	if err := cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		log.Println("Failed to create CPS action", "error", err)
		return err
	}
	log.Println("Successfully created CPS action", "userCode", userData.UserCode, "uniqueID", uniqueID)
	return nil
}

func CheckServices(self, other, agent *bool, selfServiceID, otherServiceID, agentServiceID string, serviceRepo storage.ServicesRepository, ctx context.Context, logger utils.Logger) ([]model.Service, error) {
	services, err := serviceRepo.CheckIfIDsExist(ctx, selfServiceID, otherServiceID, agentServiceID)
	if err != nil {
		logger.Errorf("[WalletCore][CheckServices] Error checking service IDs: %v", err)
		return nil, err
	}
	// Build a set of valid service IDs for quick lookup
	serviceIDSet := make(map[string]struct{}, len(services))
	for _, service := range services {
		serviceIDSet[service.ID] = struct{}{}
	}

	if selfServiceID != "" {
		if _, ok := serviceIDSet[selfServiceID]; !ok {
			logger.Errorf("[WalletCore][CheckServices] Self service ID does not exist: %s", selfServiceID)
			return nil, errors.New(localization.ErrorSelfServiceNotFound.Code)
		}
	}
	if otherServiceID != "" {
		if _, ok := serviceIDSet[otherServiceID]; !ok {
			logger.Errorf("[WalletCore][CheckServices] Other service ID does not exist: %s", otherServiceID)
			return nil, errors.New(localization.ErrorOtherServiceNotFound.Code)
		}
	}
	if agentServiceID != "" {
		if _, ok := serviceIDSet[agentServiceID]; !ok {
			logger.Errorf("[WalletCore][CheckServices] Agent service ID does not exist: %s", agentServiceID)
			return nil, errors.New(localization.ErrorAgentServiceNotFound.Code)
		}
	}

	return services, nil
}
func CheckServiceIDInWalletService(ctx context.Context, selfServiceID, otherServiceID, agentServiceID string, repo storage.WalletOracleRepository, prev *local_model.WalletOracle, logger utils.Logger) error {
	ids, err := repo.CheckServiceIDInWalletService(ctx, selfServiceID, otherServiceID, agentServiceID)
	if err != nil {
		logger.Errorf("[WalletCore][CheckServiceIDInWalletService] Error checking service IDs: %v", err)
		return err
	}
	if selfServiceID != "" && prev != nil && selfServiceID != prev.SelfServiceID && slices.Contains(ids, selfServiceID) {
		logger.Errorf("[WalletCore][CheckServiceIDInWalletService] Self service ID already exists: %s", selfServiceID)
		return errors.New(localization.ErrorSelfServiceAlreadyExist.Code)
	}
	if otherServiceID != "" && prev != nil && otherServiceID != prev.OtherServiceID && slices.Contains(ids, otherServiceID) {
		logger.Errorf("[WalletCore][CheckServiceIDInWalletService] Other service ID already exists: %s", otherServiceID)
		return errors.New(localization.ErrorOtherServiceAlreadyExists.Code)
	}
	if agentServiceID != "" && prev != nil && agentServiceID != prev.AgentServiceID && slices.Contains(ids, agentServiceID) {
		logger.Errorf("[WalletCore][CheckServiceIDInWalletService] Agent service ID already exists: %s", agentServiceID)
		return errors.New(localization.ErrorAgentServiceAlreadyExists.Code)
	}
	return nil
}

func MapToModel(wallet map[string]interface{}, logger utils.Logger) *local_model.WalletOracle {
	getString := func(key string) string {
		v, exists := wallet[key]
		if !exists || v == nil {
			logger.Warnf("[MapToModel] key '%s' missing or nil", key)
			return ""
		}
		s, ok := v.(string)
		if !ok {
			logger.Warnf("[MapToModel] key '%s' not a string (actual: %T)", key, v)
			return ""
		}
		return s
	}
	getInt := func(key string) int {
		v, exists := wallet[key]
		if !exists || v == nil {
			logger.Warnf("[MapToModel] key '%s' missing or nil", key)
			return 0
		}
		s, ok := v.(int)
		if !ok {
			logger.Warnf("[MapToModel] key '%s' not an int (actual: %T)", key, v)
			return 0
		}
		return s
	}
	getBool := func(key string) bool {
		v, exists := wallet[key]
		if !exists || v == nil {
			logger.Warnf("[MapToModel] key '%s' missing or nil", key)
			return false
		}
		b, ok := v.(bool)
		if !ok {
			logger.Warnf("[MapToModel] key '%s' not a bool (actual: %T)", key, v)
			return false
		}
		return b
	}
	getTime := func(key string) time.Time {
		v, exists := wallet[key]
		if !exists || v == nil {
			logger.Warnf("[MapToModel] key '%s' missing or nil", key)
			return time.Time{}
		}
		t, ok := v.(time.Time)
		if !ok {
			logger.Warnf("[MapToModel] key '%s' not a time.Time (actual: %T)", key, v)
			return time.Time{}
		}
		return t
	}
	getTimePtr := func(key string) *time.Time {
		v, exists := wallet[key]
		if !exists || v == nil {
			logger.Warnf("[MapToModel] key '%s' missing or nil (ptr)", key)
			return nil
		}
		t, ok := v.(time.Time)
		if !ok {
			logger.Warnf("[MapToModel] key '%s' not a time.Time (ptr, actual: %T)", key, v)
			return nil
		}
		return &t
	}

	return &local_model.WalletOracle{
		Name:                getString("name"),
		UniqueCode:          getString("unique_code"),
		Avatar:              getString("avatar"),
		Self:                getBool("self"),
		Other:               getBool("other"),
		Agent:               getBool("agent"),
		SelfServiceID:       getString("self_service_id"),
		OtherServiceID:      getString("other_service_id"),
		AgentServiceID:      getString("agent_service_id"),
		SelfServiceEnabled:  getInt("self_service_enabled"),
		OtherServiceEnabled: getInt("other_service_enabled"),
		AgentServiceEnabled: getInt("agent_service_enabled"),
		IsDeleted:           getBool("is_deleted"),
		CreatedAt:           getTime("created_at"),
		Enabled:             getBool("enabled"),
		SelfServiceName:     getString("self_service_name"),
		OtherServiceName:    getString("other_service_name"),
		AgentServiceName:    getString("agent_service_name"),
		LastModifiedAt:      getTime("last_modified_at"),
		DeletedAt:           getTimePtr("deleted_at"),
	}
}
