package utility

import (
	"context"
	"encoding/json"
	"errors"

	// "fmt"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/service"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/storage/kafka"
	local_util "cbe-super-app-cps-action/pkgs/utils"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type utilityService struct {
	cpsService service.CPSActionService
	producer   *kafka.ClientOrchestrationProducer
	logger     utils.Logger
}

func NewUtilityService(
	cpsService service.CPSActionService,
	producer *kafka.ClientOrchestrationProducer,
	logger utils.Logger,
) service.UtilityService {
	return &utilityService{
		cpsService: cpsService,
		producer:   producer,
		logger:     logger,
	}
}

// utilityActionParams holds the CPS builder inputs resolved per Kafka action type.
type utilityActionParams struct {
	uniqueId      string
	prevAction    interface{}
	requestAction string
	actionType    string
}

// resolveActionParams translates a Kafka action string into the parameters needed
// by lib.CpsModelBuilder. Non-CREATE actions fetch the existing CPS action by
// msg.ID to carry forward UniqueId and prevAction.
func (s *utilityService) resolveActionParams(ctx context.Context, msg imodel.UtilityKafkaMessage) (*utilityActionParams, error) {
	switch msg.Action {
	case "create":
		if len(msg.UniqueTokens) == 0 {
			return nil, errors.New("unique_tokens required for create action")
		}
		return &utilityActionParams{
			uniqueId:      msg.UniqueTokens[0],
			requestAction: string(cpsaction.RequestCreateUtility),
			actionType:    string(constants.CREATE),
		}, nil

	case "update", "enable", "disable", "delete":
		p := &utilityActionParams{actionType: string(constants.UPDATE)}
		switch msg.Action {
		case "update":
			p.requestAction = string(cpsaction.RequestUpdateUtility)
		case "enable":
			p.requestAction = string(cpsaction.RequestEnableUtility)
		case "disable":
			p.requestAction = string(cpsaction.RequestDisableUtility)
		case "delete":
			p.requestAction = string(cpsaction.RequestDeleteUtility)
			p.actionType = string(constants.DELETE)
		}
		existing, err := s.cpsService.GetCPSActionByID(ctx, msg.ID, "")
		if err != nil {
			return nil, errors.New("existing CPS action not found (id=%s): %w", msg.ID, err)
		}
		p.uniqueId = existing.UniqueId
		if prev, _ := local_util.JsonUnmarshal[map[string]interface{}](existing.CurrentAction); prev != nil {
			p.prevAction = *prev
		}
		return p, nil

	default:
		return nil, errors.New("unsupported utility action: %s", msg.Action)
	}
}

func (s *utilityService) HandleKafkaMessage(ctx context.Context, msg imodel.UtilityKafkaMessage) error {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	var payload map[string]interface{}
	if err := json.Unmarshal(msg.Payload, &payload); err != nil {
		log.Errorf("[UtilitySvc][HandleKafkaMessage] invalid payload: %v", err)
		return errors.New("invalid payload: %w", err)
	}

	ctx = injectMakerContext(ctx, msg.Maker)
	makerUser := local_util.ExtractUserFromContext(ctx)

	params, err := s.resolveActionParams(ctx, msg)
	if err != nil {
		log.Errorf("[UtilitySvc][HandleKafkaMessage] resolve params failed: %v", err)
		return err
	}

	action := lib.CpsModelBuilder(params.uniqueId, makerUser, params.prevAction, payload, params.requestAction, params.actionType)
	if err := s.cpsService.CreateCPSAction(ctx, &action); err != nil {
		log.Errorf("[UtilitySvc][HandleKafkaMessage] create CPS action failed (uniqueId=%s): %v", params.uniqueId, err)
		return err
	}

	log.Infof("[UtilitySvc][HandleKafkaMessage] CPS action created uniqueId=%s action=%s", params.uniqueId, msg.Action)
	return nil
}

func (s *utilityService) Authorize(ctx context.Context, cpsAction *model.CPSAction) (*model.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	switch string(cpsAction.RequestAction) {
	case string(cpsaction.RequestCreateUtility),
		string(cpsaction.RequestUpdateUtility),
		string(cpsaction.RequestEnableUtility),
		string(cpsaction.RequestDisableUtility),
		string(cpsaction.RequestDeleteUtility):

		data, err := json.Marshal(cpsAction.CurrentAction)
		if err != nil {
			log.Errorf("[UtilitySvc][Authorize] marshal failed for token %s: %v", cpsAction.UniqueId, err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		topic := "KAFKA_UTILITY_APPROVED_TOPIC"
		if err := s.producer.PublishMessage(ctx, json.RawMessage(data), "utility_approved", topic, "UTILITY.APPROVE"); err != nil {
			log.Errorf("[UtilitySvc][Authorize] publish failed for token %s: %v", cpsAction.UniqueId, err)
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}

		log.Infof("[UtilitySvc][Authorize] approved token=%s action=%s topic=%s", cpsAction.UniqueId, cpsAction.RequestAction, topic)
		return cpsAction, nil

	default:
		return nil, errors.New("unsupported utility action: %s", cpsAction.RequestAction)
	}
}

func (s *utilityService) GetByUniqueToken(ctx context.Context, uniqueToken, department string) (map[string]interface{}, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	cpsAction, err := s.cpsService.GetCPSActionByUniqueID(ctx, uniqueToken, department)
	if err != nil {
		log.Errorf("[UtilitySvc][GetByUniqueToken] not found token=%s: %v", uniqueToken, err)
		return nil, errors.New(localization.ErrorActionNotFound.Code)
	}

	payload, err := local_util.JsonUnmarshal[map[string]interface{}](cpsAction.CurrentAction)
	if err != nil {
		log.Errorf("[UtilitySvc][GetByUniqueToken] unmarshal failed token=%s: %v", uniqueToken, err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	return *payload, nil
}

func injectMakerContext(ctx context.Context, m imodel.UtilityMakerInfo) context.Context {
	ctx = context.WithValue(ctx, constants.ContextKey("user_id"), m.UserID)
	ctx = context.WithValue(ctx, constants.ContextKey("username"), m.Username)
	ctx = context.WithValue(ctx, constants.ContextKey("full_name"), m.FullName)
	ctx = context.WithValue(ctx, constants.ContextKey("phone_number"), m.PhoneNumber)
	return ctx
}
