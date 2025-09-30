package donation

import (
	"cbe-super-app-cps-action/internal/constants"
	// donation_dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	core "cbe-super-app-cps-action/internal/service/donation/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Donation struct {
	DonationRepo         storage.DonationRepository
	DonationCategoryRepo storage.DonationCategoryRepository
	DonationCompanyRepo  storage.DonationCompanyRepository
	cpsService           service.CPSActionService
	logger               utils.Logger
	minio                config.MinioClientInterface
	bucketName           string
	cfg                  *config.VaultConfig
	minioEndPoint        string
}

func NewDonationService(client *mongo.Client, DonationRepo storage.DonationRepository, DonationCategoryRepo storage.DonationCategoryRepository, DonationCompanyRepo storage.DonationCompanyRepository, cpsAction service.CPSActionService, logger utils.Logger, minio config.MinioClientInterface,
	bucketName string,
	cfg *config.VaultConfig,
	minioEndPoint string,
) service.DonationService {

	return &Donation{
		DonationRepo:         DonationRepo,
		DonationCategoryRepo: DonationCategoryRepo,
		DonationCompanyRepo:  DonationCompanyRepo,
		cpsService:           cpsAction,
		logger:               logger,
		minio:                minio,
		bucketName:           bucketName,
		cfg:                  cfg,
		minioEndPoint:        minioEndPoint,
	}
}

func (d *Donation) FetchDonation(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]dto.DonationListResponse], error) {
	return d.DonationRepo.FindAllWithPagination(ctx, *filterParams)
}

func (d *Donation) FetchDonationByID(ctx context.Context, id string) (*dto.DonationListResponse, error) {
	return d.DonationRepo.FindByID(ctx, id)
}

func (d *Donation) CreateDonation(ctx context.Context, donation dto.DonationRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	ok, err := core.DonationTitleExists(ctx, donation.Title, d.DonationRepo)
	if err != nil {
		return err
	}
	if ok {
		return errors.New(localization.ErrorDonationTitleDuplicated.Code)
	}

	_, err = d.DonationCategoryRepo.FindByID(ctx, donation.CategoryID)
	if err != nil {
		return errors.New(localization.ErrorDonationCategoryNotFound.Code)
	}

	_, err = d.DonationCompanyRepo.FindByID(ctx, donation.CompanyID)
	if err != nil {
		return errors.New(localization.ErrorDonationCompanyNotFound.Code)
	}

	coverImageURL := ""
	if donation.CoverImage != nil {
		url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donation.CoverImage, string(constants.DonationCoverImage), d.minioEndPoint, d.logger)
		if err != nil {
			return err
		}
		coverImageURL = url
	}

	donationImages := make([]types.DonationImage, 0)
	d.logger.Infof("Processing %d donation images", len(donation.DonationImages))
	for _, img := range donation.DonationImages {
		url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, img, string(constants.DonationImage), d.minioEndPoint, d.logger)
		if err != nil {
			d.logger.Errorf("Failed to upload donation image: %v", err)
			return err
		}
		donationImages = append(donationImages, types.DonationImage{
			ID:        bson.NewObjectID().Hex(),
			PhotoURL:  url,
			CreatedAt: time.Now(),
		})
		d.logger.Infof("Successfully uploaded donation image: %s", url)
	}

	donationCode := core.GenerateDonationCode()
	d.logger.Infof("Creating CPS request with target: %d, donation images count: %d", donation.Target, len(donationImages))
	result := dto.DonationCPSRequest{
		DonationCode:        donationCode,
		CompanyID:           donation.CompanyID,
		CategoryID:          donation.CategoryID,
		Title:               donation.Title,
		IsFeatured:          donation.IsFeatured,
		Target:              donation.Target,
		DonationDescription: donation.DonationDescription,
		DonationImages:      core.ConvertToDonationImages(donationImages),
		CoverImage:          coverImageURL,
		EndDate:             donation.EndDate.Format(time.RFC3339),
		StartDate:           donation.StartDate.Format(time.RFC3339),
		Enabled:             donation.Enabled,
	}
	d.logger.Infof("CPS request created with target: %d, donation images count: %d", result.Target, len(result.DonationImages))

	cpsAction := lib.CpsModelBuilder("", makerData, "", result, string(constants.RequestCreateDonation), constants.CREATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}
	return nil
}

func (d *Donation) UpdateDonation(ctx context.Context, id string, donation dto.DonationRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	existingDonation, err := d.DonationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingDonation == nil {
		return errors.New(localization.ErrorFileNotFound.Code)
	}

	existingModel := &model.Donation{
		Title:               existingDonation.Title,
		DonationDescription: existingDonation.DonationDescription,
		Target:              existingDonation.Target,
		IsFeatured:          existingDonation.IsFeatured,
		Enabled:             existingDonation.Enabled,
		StartDate:           core.ParseTime(existingDonation.StartDate),
		EndDate:             core.ParseTime(existingDonation.EndDate),
	}

	if err := core.CheckDataSimilarityAndValidation(ctx, donation, existingModel, d.DonationRepo); err != nil {
		return err
	}

	if donation.CategoryID != "" {
		_, err = d.DonationCategoryRepo.FindByID(ctx, donation.CategoryID)
		if err != nil {
			return errors.New(localization.ErrorDonationCategoryNotFound.Code)
		}
		}else{
			donation.CategoryID=existingDonation.Category.ID
		}


	if donation.CompanyID != "" {
		_, err = d.DonationCompanyRepo.FindByID(ctx, donation.CompanyID)
		if err != nil {
			return errors.New(localization.ErrorDonationCompanyNotFound.Code)
		}
	}else{
		donation.CompanyID=existingDonation.Company.ID
	}


	coverImageURL := existingDonation.CoverImage
	if donation.CoverImage != nil {
		url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donation.CoverImage, string(constants.DonationCoverImage), d.minioEndPoint, d.logger)
		if err != nil {
			return err
		}
		coverImageURL = url
	}

	donationImages := existingDonation.DonationImages
	updateData := core.MapDonationUpdate(id, existingModel, donation, coverImageURL, donationImages)

	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonation, updateData, string(constants.RequestUpdateDonation), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}

	return nil
}

func (d *Donation) UpdateDonationImage(ctx context.Context, id string, image dto.DonationImageUpdateRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	existingDonation, err := d.DonationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingDonation == nil {
		return errors.New(localization.ErrorFileNotFound.Code)
	}

	if image.Image == nil {
		return errors.New(localization.ErrorImageRequired.Code)
	}

	
	imageExists := false
	for _, img := range existingDonation.DonationImages {
		if img.ID == image.ImageID {
			imageExists = true
			break
		}
	}
	if !imageExists {
		return errors.New(localization.ErrorFileNotFound.Code)
	}

	url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, image.Image, string(constants.DonationImage), d.minioEndPoint, d.logger)
	if err != nil {
		return err
	}

	updateData := dto.DonationImageUpdateCPSRequest{
		ID:       id,
		ImageID:  image.ImageID,
		PhotoURL: url,
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonation, updateData, string(constants.RequestUpdateDonationImage), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}

	return nil
}

func (d *Donation) DeleteDonationImage(ctx context.Context, id string, imageID string) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	existingDonation, err := d.DonationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingDonation == nil {
		return errors.New(localization.ErrorFileNotFound.Code)
	}
	

	updateData := dto.DonationImageDeleteCPSRequest{
		ImageIDToDelete: imageID,
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonation, updateData, string(constants.RequestDeleteDonationImage), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}

	return nil
}

func (d *Donation) AddDonationImage(ctx context.Context, id string, image dto.DonationRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	existingDonation, err := d.DonationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingDonation == nil {
		return errors.New(localization.ErrorFileNotFound.Code)
	}

	if len(image.DonationImages) == 0 {
		return errors.New(localization.ErrorImageRequired.Code)
	}

	donationImages := make([]dto.DonationImage, 0)
	for _, img := range image.DonationImages {
		url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, img, string(constants.DonationImage), d.minioEndPoint, d.logger)
		if err != nil {
			return err
		}
		donationImages = append(donationImages, dto.DonationImage{
			ID:        bson.NewObjectID().Hex(),
			PhotoURL:  url,
			CreatedAt: time.Now().Format(time.RFC3339),
		})
	}
	//todo unique dto for the donation image array
	updateData := dto.DonationImageAddRequest{
		DonationImages: donationImages,
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonation, updateData, string(constants.RequestAddDonationImage), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}

	return nil
}

func (d *Donation) EnableDonation(ctx context.Context, id string) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	existingDonation, err := d.DonationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingDonation == nil {
		return errors.New(localization.ErrorFileNotFound.Code)
	}

	updateData := dto.EnableDonationRequest{
		Enabled: true,
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonation, updateData, string(constants.RequestEnableDonation), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}

	return nil
}

func (d *Donation) DisableDonation(ctx context.Context, id string) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	existingDonation, err := d.DonationRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingDonation == nil {
		return errors.New(localization.ErrorFileNotFound.Code)
	}

	updateData := dto.EnableDonationRequest{
		Enabled: false,
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonation, updateData, string(constants.RequestDisableDonation), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}

	return nil
}

func (d *Donation) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {

	if action.ActionStatus != constants.Approved {
		d.logger.Errorf("Tried to authorize service action without cps action approval")
		return nil, errors.New(localization.ErrorCPSActionStatusInvalid.Code)
	}
	var donationCPS *dto.DonationCPSRequest
	bindErr := core.BindAction(action.CurrentAction, &donationCPS)
	if bindErr != nil {
		d.logger.Errorf("failed to bind current action to donation: %v", bindErr)
		return nil, errors.New(localization.ErrorCPSActionFailed.Code)
	}
	if action.UniqueId != ""{
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("Failed to find donation for image addition: %v", err)
			return nil, err
		}
		if existingDonation == nil {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}

	}

	switch action.RequestAction {
	case string(constants.RequestCreateDonation):
		donationModel := core.MapToDonationModel(donationCPS)
		err := d.DonationRepo.Create(ctx, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to create donation: %v", err)
			return nil, err
		}

	case string(constants.RequestUpdateDonation):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("Failed to find existing donation: %v", err)
			return nil, err
		}
		if existingDonation == nil {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}

		existingModel := core.ConvertDonationListResponseToModel(existingDonation)

		updateRequest := dto.DonationRequest{
			CompanyID:           donationCPS.CompanyID,
			CategoryID:          donationCPS.CategoryID,
			Title:               donationCPS.Title,
			IsFeatured:          donationCPS.IsFeatured,
			Target:              donationCPS.Target,
			DonationDescription: donationCPS.DonationDescription,
			Enabled:             donationCPS.Enabled,
			StartDate:           core.ParseTime(donationCPS.StartDate),
			EndDate:             core.ParseTime(donationCPS.EndDate),
		}

		donationImages := make([]dto.DonationImage, len(donationCPS.DonationImages))
		for i, img := range donationCPS.DonationImages {
			donationImages[i] = dto.DonationImage{
				ID:        img.ID,
				PhotoURL:  img.PhotoURL,
				CreatedAt: img.CreatedAt,
			}
		}

		updateData := core.MapDonationUpdate(action.UniqueId, existingModel, updateRequest, donationCPS.CoverImage, donationImages)
		donationModel := core.MapToDonationModel(&updateData)

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to update donation: %v", err)
			return nil, err
		}

	case string(constants.RequestUpdateDonationImage):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("Failed to find donation for image update: %v", err)
			return nil, err
		}
		if existingDonation == nil {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}

		// Bind the current action to get the image update data
		var imageUpdateData *dto.DonationImageUpdateCPSRequest
		bindErr := core.BindAction(action.CurrentAction, &imageUpdateData)
		if bindErr != nil {
			d.logger.Errorf("failed to bind current action to donation image update: %v", bindErr)
			return nil, errors.New(localization.ErrorCPSActionFailed.Code)
		}

		// Update the specific image in the donation images array
		var updatedImages []types.DonationImage
		imageUpdated := false
		for _, img := range existingDonation.DonationImages {
			if img.ID == imageUpdateData.ImageID {
				// Update this specific image
				updatedImages = append(updatedImages, types.DonationImage{
					ID:        img.ID,
					PhotoURL:  imageUpdateData.PhotoURL,
					CreatedAt: core.ParseTime(img.CreatedAt),
				})
				imageUpdated = true
			} else {
				// Keep other images unchanged
				updatedImages = append(updatedImages, types.DonationImage{
					ID:        img.ID,
					PhotoURL:  img.PhotoURL,
					CreatedAt: core.ParseTime(img.CreatedAt),
				})
			}
		}

		if !imageUpdated {
			d.logger.Errorf("Image ID %s not found in donation %s", imageUpdateData.ImageID, action.UniqueId)
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		existingModel := core.ConvertDonationListResponseToModel(existingDonation)
		updateRequest := core.ConvertDonationListResponseToRequest(existingDonation)

		donationImages := make([]dto.DonationImage, len(updatedImages))
		for i, img := range updatedImages {
			donationImages[i] = dto.DonationImage{
				ID:        img.ID,
				PhotoURL:  img.PhotoURL,
				CreatedAt: img.CreatedAt.Format(time.RFC3339),
			}
		}

		updateData := core.MapDonationUpdate(action.UniqueId, existingModel, updateRequest, existingDonation.CoverImage, donationImages)
		donationModel := core.MapToDonationModel(&updateData)

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to update donation image: %v", err)
			return nil, err
		}

	case string(constants.RequestDeleteDonationImage):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("Failed to find donation for image deletion: %v", err)
			return nil, err
		}
		if existingDonation == nil {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}

		existingModel := core.ConvertDonationListResponseToModel(existingDonation)
		updateRequest := core.ConvertDonationListResponseToRequest(existingDonation)

		var updatedImages []types.DonationImage
		for _, img := range existingDonation.DonationImages {
			if img.ID != donationCPS.ImageIDToDelete {
				updatedImages = append(updatedImages, types.DonationImage{
					ID:        img.ID,
					PhotoURL:  img.PhotoURL,
					CreatedAt: core.ParseTime(img.CreatedAt),
				})
			}
		}

		donationImages := make([]dto.DonationImage, len(updatedImages))
		for i, img := range updatedImages {
			donationImages[i] = dto.DonationImage{
				ID:        img.ID,
				PhotoURL:  img.PhotoURL,
				CreatedAt: img.CreatedAt.Format(time.RFC3339),
			}
		}

		updateData := core.MapDonationUpdate(action.UniqueId, existingModel, updateRequest, existingDonation.CoverImage, donationImages)
		donationModel := core.MapToDonationModel(&updateData)

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to delete donation image: %v", err)
			return nil, err
		}

	case string(constants.RequestAddDonationImage):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("Failed to find donation for image addition: %v", err)
			return nil, err
		}
		if existingDonation == nil {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}

		existingModel := core.ConvertDonationListResponseToModel(existingDonation)
		updateRequest := core.ConvertDonationListResponseToRequest(existingDonation)

		updatedImages := make([]types.DonationImage, len(existingDonation.DonationImages))
		for i, img := range existingDonation.DonationImages {
			updatedImages[i] = types.DonationImage{
				ID:        img.ID,
				PhotoURL:  img.PhotoURL,
				CreatedAt: core.ParseTime(img.CreatedAt),
			}
		}

		for _, newImg := range donationCPS.DonationImages {
			updatedImages = append(updatedImages, types.DonationImage{
				ID:        newImg.ID,
				PhotoURL:  newImg.PhotoURL,
				CreatedAt: time.Now(),
			})
		}

		donationImages := make([]dto.DonationImage, len(updatedImages))
		for i, img := range updatedImages {
			donationImages[i] = dto.DonationImage{
				ID:        img.ID,
				PhotoURL:  img.PhotoURL,
				CreatedAt: img.CreatedAt.Format(time.RFC3339),
			}
		}

		updateData := core.MapDonationUpdate(action.UniqueId, existingModel, updateRequest, existingDonation.CoverImage, donationImages)
		donationModel := core.MapToDonationModel(&updateData)

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to add donation images: %v", err)
			return nil, err
		}

	case string(constants.RequestEnableDonation):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("Failed to find donation for image addition: %v", err)
			return nil, err
		}
		if existingDonation == nil {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		existingModel := core.ConvertDonationListResponseToModel(existingDonation)
		updateRequest := dto.DonationRequest{
			Enabled:             true,
			
		}
		updateData := core.MapDonationUpdate(action.UniqueId, existingModel, updateRequest, existingDonation.CoverImage, existingDonation.DonationImages)
		donationModel := core.MapToDonationModel(&updateData)


		
		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to enable donation: %v", err)
			return nil, err
		}

	case string(constants.RequestDisableDonation):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("Failed to find donation for image addition: %v", err)
			return nil, err
		}
		if existingDonation == nil {
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		existingModel := core.ConvertDonationListResponseToModel(existingDonation)
		updateRequest := dto.DonationRequest{
			Enabled:             false,
			
		}
		updateData := core.MapDonationUpdate(action.UniqueId, existingModel, updateRequest, existingDonation.CoverImage, existingDonation.DonationImages)
		donationModel := core.MapToDonationModel(&updateData)
		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to disable donation: %v", err)
			return nil, err
		}

	default:
		d.logger.Errorf("Unsupported action requested: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	d.logger.Infof("Service action authorization completed: %s", action.RequestAction)
	return action, nil
}
