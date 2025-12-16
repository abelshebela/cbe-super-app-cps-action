package topup

import (
	"cbe-super-app-cps-action/internal/constants"
	topupDto "cbe-super-app-cps-action/internal/constants/dto/topup"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	"cbe-super-app-cps-action/internal/service/topup/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"path"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type topupService struct {
	repo        storage.TopupRepository
	cpsService  service.CPSActionService
	logger      utils.Logger
	minio       *s3.Client
	bucketName  string
	minioPubUrl string
	cfg         *config.VaultConfig
}

func NewTopupService(repo storage.TopupRepository, cps service.CPSActionService, minio *s3.Client, minioPubUrl string, bucketName string, cfg *config.VaultConfig, logger utils.Logger) service.TopupService {
	return &topupService{
		repo:       repo,
		cpsService: cps,
		logger:     logger,
		minio:      minio,
		bucketName: bucketName,
		cfg:        cfg,
	}
}

func (s *topupService) CreateTopup(ctx context.Context, req topupDto.TopupRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateTopup", "topupService", "topupService")
	defer span.End()
	s.logger.Infof("Createtopup called", "topup_name", req.Name)

	exist, err := s.repo.Find(ctx, req.Code, req.Name)
	if err != nil {
		span.AddEvent("Repo find error", trace.WithAttributes(attribute.String("error", err.Error())))
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	if exist != nil {
		span.AddEvent("Topup name already exists", trace.WithAttributes(attribute.String("name", req.Name)))
		return errors.New(localization.ErrorTopupNameAlreadyExists.Code)
	}

	code, err := core.GeneratePrefixedName("TOP", req.Code, s.logger)
	if err != nil {
		span.AddEvent("GeneratePrefixedName error", trace.WithAttributes(attribute.String("error", err.Error())))
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	req.Code = code

	existCode, err := s.repo.Find(ctx, "code", req.Code)
	if err != nil {
		span.AddEvent("Repo find by code error", trace.WithAttributes(attribute.String("error", err.Error())))
		return errors.New(localization.ErrorUnhandledServer.Code)
	}
	if existCode != nil {
		span.AddEvent("Topup code already exists", trace.WithAttributes(attribute.String("code", req.Code)))
		return errors.New(localization.ErrorTopupCodeAlreadyExists.Code)
	}

	URL, err := lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Avatar, "topup", *s.cfg, "", s.logger)
	if err != nil {
		span.AddEvent("UploadFileToMinio error", trace.WithAttributes(attribute.String("error", err.Error())))
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	topup := core.ToCreateTopupDoc(req.Name, code, URL, req.Self, req.Other, req.Agent)
	topup.Enabled = true
	//here since the unique id is nil 000.. use other unique id like the code
	// if err := core.HandleCPSAction(ctx, s.cpsService, topup.ID.Hex(), constants.RequestCreatetopup, topup, nil, constants.ActionCreate); err != nil {
	if err := core.HandleCPSAction(ctx, s.cpsService, topup.ID.Hex(), constants.RequestCreateTopup, topup, nil, constants.ActionCreate); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("code", topup.Code)))
		s.logger.Errorf("CPS action failed for topup %s: %v", topup.Code, err)
		return err
	}

	span.AddEvent("Topup created", trace.WithAttributes(attribute.String("code", topup.Code)))
	return nil
}

func (s *topupService) UpdateTopup(ctx context.Context, id string, req topupDto.TopupRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateTopup", "topupService", "topupService")
	defer span.End()
	s.logger.Infof("Updatetopup called", "topup_id", id)

	prevtopup, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return errors.New(localization.ErrorTopupNotFound.Code)
	}

	if req.Name != "" {
		exist, err := s.repo.Find(ctx, "name", req.Name)
		if err != nil {
			span.AddEvent("Repo find by name error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("name", req.Name)))
			return errors.New(localization.ErrorUnhandledServer.Code)
		}

		if exist != nil && exist.ID.Hex() != id {
			span.AddEvent("Topup name already exists", trace.WithAttributes(attribute.String("name", req.Name)))
			return errors.New(localization.ErrorTopupNameAlreadyExists.Code)
		}
	}

	code, err := core.GeneratePrefixedName("TOP", req.Code, s.logger)
	if err != nil {
		span.AddEvent("GeneratePrefixedName error", trace.WithAttributes(attribute.String("error", err.Error())))
		return errors.New(localization.ErrorUnhandledServer.Code)
	}

	req.Code = code
	if req.Code != "" {
		exist, err := s.repo.Find(ctx, "code", req.Code)
		if err != nil {
			span.AddEvent("Repo find by code error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("code", req.Code)))
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
		if exist != nil && exist.ID.Hex() != id {
			span.AddEvent("Topup code already exists", trace.WithAttributes(attribute.String("code", req.Code)))
			return errors.New(localization.ErrorTopupCodeAlreadyExists.Code)
		}
	}

	var avatarURL string
	if req.Avatar != nil {
		var objectkey string
		if prevtopup.Avatar != "" {
			objectkey = path.Base(prevtopup.Avatar)
		}

		avatarURL, err = lib.UploadFileToMinio(ctx, s.minio, s.bucketName, req.Avatar, "avatar", *s.cfg, objectkey, s.logger)
		if err != nil {
			span.AddEvent("UploadFileToMinio error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
			return errors.New(localization.ErrorUnhandledServer.Code)
		}
	} else {
		avatarURL = prevtopup.Avatar
	}
	Updatetopup, change_count := core.ToUpdateTopupDoc(*prevtopup, req)
	if avatarURL != prevtopup.Avatar {
		change_count++
	}
	Updatetopup.Avatar = avatarURL

	if change_count == 0 {
		span.AddEvent("No changes detected", trace.WithAttributes(attribute.String("id", id)))
		return errors.New(localization.ErrorNoChangesDetected.Code)
	}

	if !(Updatetopup.Services.Agent || Updatetopup.Services.Self || Updatetopup.Services.Other) {
		span.AddEvent("No service option selected", trace.WithAttributes(attribute.String("id", id)))
		return errors.New(localization.ErrorTopupServiceOption.Code)
	}

	if err := core.HandleCPSAction(ctx, s.cpsService, id, constants.RequestUpdateTopup, Updatetopup, *prevtopup, constants.ActionUpdate); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("code", Updatetopup.Code)))
		s.logger.Errorf("CPS action failed for topup %s: %v", Updatetopup.Code, err)
		return err
	}

	span.AddEvent("Topup updated", trace.WithAttributes(attribute.String("id", id)))
	return nil
}

func (s *topupService) DeleteTopup(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteTopup", "topupService", "topupService")
	defer span.End()
	s.logger.Infof("Deletetopup called", "topup_id", id)

	prevtopup, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return errors.New(localization.ErrorTopupNotFound.Code)
	}

	now := time.Now()
	deletedtopup := *prevtopup
	deletedtopup.IsDeleted = true
	deletedtopup.DeletedAt = now

	if err := core.HandleCPSAction(ctx, s.cpsService, id, constants.RequestDeleteTopup, deletedtopup, *prevtopup, constants.ActionDelete); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("code", deletedtopup.Code)))
		s.logger.Errorf("CPS action failed for topup %s: %v", deletedtopup.Code, err)
		return err
	}

	span.AddEvent("Topup deleted", trace.WithAttributes(attribute.String("id", id)))
	return nil
}

func (s *topupService) EnableOrDisableTopup(ctx context.Context, id string, enable bool) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableOrDisableTopup", "topupService", "topupService")
	defer span.End()
	s.logger.Infof("Enable/disable called", "enable:", enable, "id:", id)

	prevtopup, err := s.repo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("FindByID error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("id", id)))
		return errors.New(localization.ErrorTopupNotFound.Code)
	}

	if enable && prevtopup.Enabled {
		span.AddEvent("Topup already enabled", trace.WithAttributes(attribute.String("id", id)))
		return errors.New(localization.ErrorTopupAlreadyEnabled.Code)
	}
	if !enable && !prevtopup.Enabled {
		span.AddEvent("Topup already disabled", trace.WithAttributes(attribute.String("id", id)))
		return errors.New(localization.ErrorTopupAlreadyDisabled.Code)
	}
	s.logger.Infof("you can enable or disable", enable, prevtopup.Enabled)

	updatedtopup := *prevtopup
	updatedtopup.Enabled = enable
	updatedtopup.LastModifiedAt = time.Now()

	var action constants.RequestAction
	if enable {
		action = constants.RequestEnableTopup
	} else {
		action = constants.RequestDisableTopup
	}

	if err := core.HandleCPSAction(ctx, s.cpsService, id, action, updatedtopup, *prevtopup, constants.ActionUpdate); err != nil {
		span.AddEvent("CPS action failed", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("code", updatedtopup.Code)))
		s.logger.Errorf("CPS action failed for topup %s: %v", updatedtopup.Code, err)
		return err
	}

	span.AddEvent("Topup enable/disable updated", trace.WithAttributes(attribute.String("id", id), attribute.Bool("enabled", enable)))
	return nil
}

func (s *topupService) GetTopup(ctx context.Context, id string) (*model.Topup, error) {
	s.logger.Infof("Fetchtopup called", "topup_id", id)

	return s.repo.FindByID(ctx, id)
}

func (s *topupService) GetAllTopup(ctx context.Context, filterParams types.Filter) (*types.PaginatedResponse[[]*model.Topup], error) {
	s.logger.Infof("FetchAllcalled")

	return s.repo.FindAllWithPagination(ctx, filterParams)
}

func (s *topupService) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorizetopup", "topupService", "topupService")
	defer span.End()
	s.logger.Infof("Authorizetopup called", "cpsaction", action.ActionCode)

	topup, err := local_util.JsonUnmarshal[model.Topup](action.CurrentAction)
	if err != nil {
		span.AddEvent("JsonUnmarshal error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("action_code", action.ActionCode)))
		return nil, errors.New(localization.ErrorInvalidActionData.Code)
	}

	switch action.RequestAction {
	case string(constants.RequestCreateTopup):
		err = s.repo.Create(ctx, topup)
	case string(constants.RequestUpdateTopup):
		err = s.repo.Update(ctx, action.UniqueId, topup)
	case string(constants.RequestDeleteTopup):
		err = s.repo.Delete(ctx, action.UniqueId)
	case string(constants.RequestEnableTopup):
		err = s.repo.EnableOrDisable(ctx, action.UniqueId, true)
	case string(constants.RequestDisableTopup):
		err = s.repo.EnableOrDisable(ctx, action.UniqueId, false)
	default:

		return nil, errors.New(localization.ErrorInvalidRequest.Code)
	}

	if err != nil {
		span.AddEvent("Topup action error", trace.WithAttributes(attribute.String("error", err.Error()), attribute.String("action_code", action.ActionCode)))
		return nil, err
	}

	action.CurrentAction = topup
	return action, nil
}
