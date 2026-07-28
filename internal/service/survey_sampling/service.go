package survey_sampling_service

import (
	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	cpsaction "cbe-super-app-cps-action/internal/service/cps_action"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"encoding/json"
	"errors"
	"time"

	sharedmodel "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type surveySamplingService struct {
	repo       storage.SurveySamplingRepository
	cpsService service.CPSActionService
	logger     utils.Logger
}

func NewSurveySamplingService(repo storage.SurveySamplingRepository, cpsService service.CPSActionService, logger utils.Logger) service.SurveySamplingService {
	return &surveySamplingService{
		repo:       repo,
		cpsService: cpsService,
		logger:     logger,
	}
}

func (s *surveySamplingService) Create(ctx context.Context, req imodel.SurveySamplingConfig) error {
	maker := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, s.logger)

	if local_util.IsIncomplete(maker) {
		log.Errorf("[SurveySamplingSvc][Create] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	cpsModel := lib.CpsModelBuilder(string(req.Method), maker, nil, req, string(cpsaction.RequestCreateSurveySampling), constants.CREATE)
	return s.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func (s *surveySamplingService) Update(ctx context.Context, id string, req imodel.SurveySamplingConfig) error {
	maker := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, s.logger)

	if local_util.IsIncomplete(maker) {
		log.Errorf("[SurveySamplingSvc][Update] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	prev, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	req.ID = prev.ID
	req.Method = prev.Method
	req.CreatedAt = prev.CreatedAt
	req.UpdatedAt = time.Now()

	cpsModel := lib.CpsModelBuilder(id, maker, prev, req, string(cpsaction.RequestUpdateSurveySampling), constants.UPDATE)
	return s.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func (s *surveySamplingService) Enable(ctx context.Context, id string) error {
	maker := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, s.logger)

	if local_util.IsIncomplete(maker) {
		log.Errorf("[SurveySamplingSvc][Enable] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	cpsModel := lib.CpsModelBuilder(id, maker, existing, existing, string(cpsaction.RequestEnableSurveySampling), constants.UPDATE)
	return s.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func (s *surveySamplingService) Disable(ctx context.Context, id string) error {
	maker := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, s.logger)

	if local_util.IsIncomplete(maker) {
		log.Errorf("[SurveySamplingSvc][Disable] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	cpsModel := lib.CpsModelBuilder(id, maker, existing, existing, string(cpsaction.RequestDisableSurveySampling), constants.UPDATE)
	return s.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func (s *surveySamplingService) Delete(ctx context.Context, id string) error {
	maker := local_util.ExtractUserFromContext(ctx)
	log := local_util.LoggerFromCtx(ctx, s.logger)

	if local_util.IsIncomplete(maker) {
		log.Errorf("[SurveySamplingSvc][Delete] maker data is incomplete")
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	cpsModel := lib.CpsModelBuilder(id, maker, existing, nil, string(cpsaction.RequestDeleteSurveySampling), constants.DELETE)
	return s.cpsService.CreateCPSAction(ctx, &cpsModel)
}

func (s *surveySamplingService) GetAll(ctx context.Context, filter types.Filter) (*types.PaginatedResponse[[]imodel.SurveySamplingConfig], error) {
	return s.repo.FindAllWithPagination(ctx, filter)
}

func (s *surveySamplingService) GetByID(ctx context.Context, id string) (*imodel.SurveySamplingConfig, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *surveySamplingService) Authorize(ctx context.Context, cpsAction *sharedmodel.CPSAction) (*sharedmodel.CPSAction, error) {
	log := local_util.LoggerFromCtx(ctx, s.logger)

	raw, err := json.Marshal(cpsAction.CurrentAction)
	if err != nil {
		log.Errorf("[SurveySamplingSvc][Authorize] marshal failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	var config imodel.SurveySamplingConfig
	if err := json.Unmarshal(raw, &config); err != nil {
		log.Errorf("[SurveySamplingSvc][Authorize] unmarshal failed: %v", err)
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	if cpsAction.UniqueId != "" {
		if oid, err := bson.ObjectIDFromHex(cpsAction.UniqueId); err == nil {
			config.ID = oid
		}
	}

	switch cpsAction.RequestAction {
	case string(cpsaction.RequestCreateSurveySampling):
		config.CreatedAt = time.Now()
		config.UpdatedAt = time.Now()
		if err := s.repo.Create(ctx, config); err != nil {
			log.Errorf("[SurveySamplingSvc][Authorize] create failed: %v", err)
			return nil, err
		}
	case string(cpsaction.RequestUpdateSurveySampling):
		config.UpdatedAt = time.Now()
		if err := s.repo.Update(ctx, config.ID.Hex(), config); err != nil {
			log.Errorf("[SurveySamplingSvc][Authorize] update failed: %v", err)
			return nil, err
		}
	case string(cpsaction.RequestEnableSurveySampling):
		if err := s.repo.SetEnabled(ctx, cpsAction.UniqueId, true); err != nil {
			log.Errorf("[SurveySamplingSvc][Authorize] enable failed: %v", err)
			return nil, err
		}
	case string(cpsaction.RequestDisableSurveySampling):
		if err := s.repo.SetEnabled(ctx, cpsAction.UniqueId, false); err != nil {
			log.Errorf("[SurveySamplingSvc][Authorize] disable failed: %v", err)
			return nil, err
		}
	case string(cpsaction.RequestDeleteSurveySampling):
		if err := s.repo.Delete(ctx, cpsAction.UniqueId); err != nil {
			log.Errorf("[SurveySamplingSvc][Authorize] delete failed: %v", err)
			return nil, err
		}
	default:
		log.Errorf("[SurveySamplingSvc][Authorize] unsupported action: %s", cpsAction.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	cpsAction.CurrentAction = config
	return cpsAction, nil
}
