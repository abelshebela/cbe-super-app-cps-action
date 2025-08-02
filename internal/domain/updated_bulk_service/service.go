package updatedbulkservice

import (
	"context"
	"fmt"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	cps_constant "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/constant"
	entity "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/cps_actions/entities"
	ctx_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/context"
	common_util "github.com/CBE-Super-App/cbe-super-app-cps-action/pkgs/utils"
	constant "github.com/CBE-Super-App/cbe-super-app-cps-action/utils"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BulkService interface {
	GetAllBulkServices(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*APPAccessList], error)
	EnableBulkService(ctx context.Context, keys []string) error
	DisableBulkService(ctx context.Context, keys []string) error
	Authorize(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error)
}

type bulkServiceImpl struct {
	repo   BulkServiceRespository
	logger utils.Logger
}

func NewBulkService(repo BulkServiceRespository, logger utils.Logger) BulkService {
	return &bulkServiceImpl{repo: repo, logger: logger}
}

func (b *bulkServiceImpl) GetAllBulkServices(ctx context.Context, filterParams *constant.Filter) (*common_util.PaginatedResponse[[]*APPAccessList], error) {
	return b.repo.GetAllBulkServices(ctx, filterParams)
}

func (b *bulkServiceImpl) EnableBulkService(ctx context.Context, keys []string) error {
	userPayload := ctx_util.ExtractContext(ctx)
	actionCode := utils.RandomGenerator(24)

	type currAction struct {
		Keys []string
	}

	cpsAction := model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       actionCode,
		UniqueId:         userPayload.UserCode,
		MakerID:          userPayload.UserID,
		MakerName:        userPayload.FullName,
		MakerPhoneNumber: userPayload.PhoneNumber,
		Department:       userPayload.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionEnable),
		RequestAction:    string(model.RequestBlockCity),
		PreviousAction:   nil,
		CurrentAction: currAction{
			Keys: keys,
		},
		CreatedAt:       time.Now(),
		MakerActionTime: time.Now(),
	}

	return b.repo.EnableOrDisableBulkService(ctx, keys, cpsAction, model.RequestBulkServiceEnable)
}

func (b *bulkServiceImpl) DisableBulkService(ctx context.Context, keys []string) error {
	userPayload := ctx_util.ExtractContext(ctx)
	actionCode := utils.RandomGenerator(24)

	// var actions []APPAccessList
	// for _, code := range serviceCodes {
	// 	actions = append(actions, APPAccessList{
	// 		Key:     code,
	// 		Enabled: true,
	// 	})
	// }

	type currAction struct {
		Keys []string
	}

	cpsAction := model.CPSAction{
		ID:               bson.NewObjectID(),
		ActionCode:       actionCode,
		UniqueId:         userPayload.UserCode,
		MakerID:          userPayload.UserID,
		MakerName:        userPayload.FullName,
		MakerPhoneNumber: userPayload.PhoneNumber,
		Department:       userPayload.Department,
		ActionStatus:     string(model.ActionPending),
		ActionType:       string(model.ActionDisable),
		RequestAction:    string(model.RequestBulkServiceDisable),
		PreviousAction:   nil,
		CurrentAction: currAction{
			Keys: keys,
		},
		CreatedAt:       time.Now(),
		MakerActionTime: time.Now(),
	}

	return b.repo.EnableOrDisableBulkService(ctx, keys, cpsAction, model.RequestBulkServiceDisable)
}

func (b *bulkServiceImpl) Authorize(ctx context.Context, action *entity.CPSAction) (*entity.CPSAction, error) {
	action.MakerActionTime = time.Now()
	action.LastModifiedAt = time.Now()

	switch action.ActionType {
	case cps_constant.ActionEnable:
		return b.repo.AuthorizeBulkServiceEnable(ctx, action)
	case cps_constant.ActionDisable:
		return b.repo.AuthorizeBulkServiceDisable(ctx, action)
	default:
		return nil, fmt.Errorf("UNHANDLED_SERVER_ERROR")
	}
}
