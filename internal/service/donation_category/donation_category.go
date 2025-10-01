package donation_category

import (
	"cbe-super-app-cps-action/internal/constants"
	dto "cbe-super-app-cps-action/internal/constants/dto/donation_category"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	core "cbe-super-app-cps-action/internal/service/donation_category/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DonationCategory struct {
	DonationCategoryRepo storage.DonationCategoryRepository
	cpsService           service.CPSActionService
	logger               utils.Logger
	minio                config.MinioClientInterface
	bucketName           string
	cfg                  *config.VaultConfig
	minioEndPoint        string
}

func NewDonationCategoryService(client *mongo.Client, DonationCategoryRepo storage.DonationCategoryRepository, cpsAction service.CPSActionService, logger utils.Logger, minio config.MinioClientInterface,
	bucketName string,
	cfg *config.VaultConfig,
	minioEndPoint string,
) service.DonationCategoryService {

	return &DonationCategory{
		DonationCategoryRepo: DonationCategoryRepo,
		cpsService:           cpsAction,
		logger:               logger,
		minio:                minio,
		bucketName:           bucketName,
		cfg:                  cfg,
		minioEndPoint:        minioEndPoint,
	}
}

func (d *DonationCategory) FetchDonationCategory(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]dto.DonationCategoryListResponse], error) {
	return d.DonationCategoryRepo.FindAllWithPagination(ctx, *filterParams)
}

func (d *DonationCategory) FetchDonationCategoryByID(ctx context.Context, id string) (*dto.DonationCategoryListResponse, error) {
	return d.DonationCategoryRepo.FindByID(ctx, id)
}

func (d *DonationCategory) CreateDonationCategory(ctx context.Context, donationCategory dto.DonationCategoryRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}
	ok, err := core.DonationNameExists(ctx, donationCategory.CategoryName, d.DonationCategoryRepo)
	if err != nil {
		return err
	}
	if ok {
		return errors.New(localization.ErrorDonationCategoryNameDuplicated.Code)
	}
	if donationCategory.Icon == nil {
		return errors.New(localization.ErrorDonationCategoryIDRequired.Code)
	}
	url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donationCategory.Icon, string(constants.DonationIcon), d.minioEndPoint, d.logger)
	if err != nil {
		return err
	}
	result := dto.DonationCategoryResponse{
		CategoryName: donationCategory.CategoryName,
		Icon:         url,
	}

	cpsAction := lib.CpsModelBuilder("", makerData, "", result, string(constants.RequestCreateDonationCategory), constants.CREATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}
	return nil

}

func (d *DonationCategory) UpdateDonationCategory(ctx context.Context, id string, donationCategory dto.DonationCategoryRequest) (dto.DonationCategoryRequest, error) {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return donationCategory, errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	// Check if donation category exists
	existingCategory, err := d.DonationCategoryRepo.FindByID(ctx, id)
	if err != nil {
		return donationCategory, err
	}
	if existingCategory == nil {
		return donationCategory, errors.New(localization.ErrorFileNotFound.Code)
	}

	// Convert DTO to model for similarity check
	existingModel := &model.DonationCategory{
		CategoryName: existingCategory.CategoryName,
		Icon:         existingCategory.Icon,
		IsDeleted:    existingCategory.IsDeleted,
	}

	// Check if data is similar to existing data
	if core.IsDataSimilar(donationCategory, existingModel) {
		return donationCategory, errors.New(localization.ErrorNoChangesToUpdate.Code)
	}

	// Check if category name is being updated and if it already exists
	if donationCategory.CategoryName != "" && donationCategory.CategoryName != existingCategory.CategoryName {
		ok, err := core.DonationNameExists(ctx, donationCategory.CategoryName, d.DonationCategoryRepo)
		if err != nil {
			return donationCategory, err
		}
		if ok {
			return donationCategory, errors.New(localization.ErrorDonationCategoryNameDuplicated.Code)
		}
	}
	iconURL := existingCategory.Icon
	if donationCategory.Icon != nil {
		url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donationCategory.Icon, string(constants.DonationIcon), d.minioEndPoint, d.logger)
		if err != nil {
			return donationCategory, err
		}
		iconURL = url
	}
	updateData := dto.DonationCategoryCPSRequest{
		ID:           id,
		CategoryName: donationCategory.CategoryName,
		Icon:         iconURL,
	}

	if donationCategory.CategoryName == "" {
		updateData.CategoryName = existingCategory.CategoryName
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, existingCategory, updateData, string(constants.RequestUpdateDonationCategory), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return donationCategory, err
	}

	return dto.DonationCategoryRequest{
		CategoryName: updateData.CategoryName,
	}, nil
}

func (d *DonationCategory) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	// requestedAction := action.RequestAction
	fmt.Println("Donation Category Authorize called")
	if action.ActionStatus != constants.Approved {
		d.logger.Errorf("Tried to authorize service action without cps action approval")
		return nil, errors.New(localization.ErrorCPSActionStatusInvalid.Code)
	}
	var donationCPS *dto.DonationCategoryCPSRequest
	bindErr := core.BindAction(action.CurrentAction, &donationCPS)
	if bindErr != nil {
		d.logger.Errorf("failed to bind current action to donation Category: %v", bindErr)
		return nil, errors.New(localization.ErrorCPSActionFailed.Code)
	}

	switch action.RequestAction {
	case string(constants.RequestCreateDonationCategory):
		donationModel := core.MapToDonationCategory(donationCPS.CategoryName, donationCPS.Icon)

		err := d.DonationCategoryRepo.Create(ctx, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to create donation category: %v", err)
			return nil, err
		}

	case string(constants.RequestUpdateDonationCategory):
		donationModel := core.MapToDonationCategory(donationCPS.CategoryName, donationCPS.Icon)
		err := d.DonationCategoryRepo.Update(ctx, donationCPS.ID, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to update donation category: %v", err)
			return nil, err
		}
	default:
		d.logger.Errorf("Unsupported action requested: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	d.logger.Infof("Service action authorization completed: %s", action.RequestAction)
	return action, nil

}
