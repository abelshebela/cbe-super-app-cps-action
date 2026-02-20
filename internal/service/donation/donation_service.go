package donation

import (
	"cbe-super-app-cps-action/internal/constants"
	"path"

	dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	core "cbe-super-app-cps-action/internal/service/donation/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	imodel "cbe-super-app-cps-action/internal/constants/model"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	// types "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Donation struct {
	DonationRepo         storage.DonationRepository
	DonationCategoryRepo storage.DonationCategoryRepository
	DonationCompanyRepo  storage.DonationCompanyRepository
	cpsService           service.CPSActionService
	logger               utils.Logger
	minio                *s3.Client
	bucketName           string
	cfg                  *config.VaultConfig
	minioEndPoint        string
}

func NewDonationService(client *mongo.Client, DonationRepo storage.DonationRepository, DonationCategoryRepo storage.DonationCategoryRepository, DonationCompanyRepo storage.DonationCompanyRepository, cpsAction service.CPSActionService, logger utils.Logger, minio *s3.Client,
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
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchDonation", "Donation", "FetchDonation")
	defer span.End()

	data, err := d.DonationRepo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Failed to fetch donations", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	return data, nil
}

func (d *Donation) FetchDonationByID(ctx context.Context, id string) (*dto.DonationListResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchDonationByID", "Donation", "FetchDonationByID")
	defer span.End()

	res, err := d.DonationRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to fetch donation", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	return res, nil
}

func (d *Donation) CreateDonation(ctx context.Context, donation dto.DonationRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateDonation", "Donation", "CreateDonation")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberRequired.Code),
		))
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	ok, err := core.DonationTitleExists(ctx, donation.Title, d.DonationRepo)
	if err != nil {
		span.AddEvent("Failed to check donation title existence", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("title", donation.Title),
		))
		return err
	}
	if ok {
		span.AddEvent("Donation title duplicated", trace.WithAttributes(
			attribute.String("error", localization.ErrorDonationTitleDuplicated.Code),
			attribute.String("title", donation.Title),
		))
		return errors.New(localization.ErrorDonationTitleDuplicated.Code)
	}

	category, err := d.DonationCategoryRepo.FindByID(ctx, donation.CategoryID)
	if err != nil {
		span.AddEvent("Donation category not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorDonationCategoryNotFound.Code),
			attribute.String("category_id", donation.CategoryID),
		))
		return errors.New(localization.ErrorDonationCategoryNotFound.Code)
	}
	if !category.Enabled {
		span.AddEvent("Category is not enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorCategoryIsNotEnabled.Code),
			attribute.String("category_id", donation.CategoryID),
		))
		return errors.New(localization.ErrorCategoryIsNotEnabled.Code)
	}

	company, err := d.DonationCompanyRepo.FindByID(ctx, donation.CompanyID)
	if err != nil {
		span.AddEvent("Donation company not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorDonationCompanyNotFound.Code),
			attribute.String("company_id", donation.CompanyID),
		))
		return errors.New(localization.ErrorDonationCompanyNotFound.Code)
	}
	if !company.Enabled {
		span.AddEvent("Company is not enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorCompanyIsNotEnabled.Code),
			attribute.String("company_id", donation.CompanyID),
		))
		return errors.New(localization.ErrorCompanyIsNotEnabled.Code)
	}

	coverImageURL := ""
	if donation.CoverImage != nil {
		url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donation.CoverImage, string(constants.DonationFolderName), *d.cfg, "", d.logger)
		if err != nil {
			span.AddEvent("Failed to upload cover image", trace.WithAttributes(
				attribute.String("error", err.Error()),
			))
			return err
		}
		coverImageURL = url
	}

	donationImages := make([]types.DonationImage, 0)
	d.logger.Infof("Processing %d donation images", len(donation.DonationImages))
	for _, img := range donation.DonationImages {
		url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, img, string(constants.DonationFolderName), *d.cfg, "", d.logger)
		if err != nil {
			d.logger.Errorf("Failed to upload donation image: %v", err)
			span.AddEvent("Failed to upload donation image", trace.WithAttributes(
				attribute.String("error", err.Error()),
			))
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
	tempval := false
	result := dto.DonationCPSRequest{
		DonationCode: donationCode,
		Company: dto.Company{
			ID:            donation.CompanyID,
			CompanyName:   company.CompanyName,
			CompanyLogo:   company.CompanyLogo,
			AccountNumber: company.AccountNumber,
			Enabled:       company.Enabled,
		},
		Category: dto.Category{
			ID:           donation.CategoryID,
			CategoryName: category.CategoryName,
			Icon:         category.Icon,
		},
		Title:      donation.Title,
		IsFeatured: donation.IsFeatured,
		Target:     donation.Target,

		DonationDescription: donation.DonationDescription,
		DonationImages:      core.ConvertToDonationImages(donationImages),
		CoverImage:          coverImageURL,
		EndDate:             donation.EndDate.Format(time.RFC3339),
		StartDate:           donation.StartDate.Format(time.RFC3339),
		Enabled:             &tempval,
	}
	d.logger.Infof("CPS request created with target: %d, donation images count: %d", result.Target, len(result.DonationImages))

	cpsAction := lib.CpsModelBuilder("", makerData, "", result, string(constants.RequestCreateDonation), constants.CREATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}
	return nil
}

func (d *Donation) UpdateDonation(ctx context.Context, id string, donation dto.DonationRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateDonation", "Donation", "UpdateDonation")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if local_util.IsIncomplete(makerData) {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberRequired.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	existingDonation, err := d.DonationRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find donation", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if existingDonation == nil {
		span.AddEvent("Donation not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorFileNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorFileNotFound.Code)
	}

	// Prepare existing imodel for validation
	existingModel := &imodel.Donation{
		DonationCode:        existingDonation.DonationCode,
		Title:               existingDonation.Title,
		DonationDescription: existingDonation.DonationDescription,
		Target:              existingDonation.Target,
		IsFeatured:          existingDonation.IsFeatured,
		Enabled:             existingDonation.Enabled,
		StartDate:           core.ParseTime(existingDonation.StartDate),
		EndDate:             core.ParseTime(existingDonation.EndDate),
	}

	// Validate and check for data duplication
	if err := core.CheckDataSimilarityAndValidation(ctx, donation, existingModel, d.DonationRepo); err != nil {
		span.AddEvent("Data similarity validation failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	// --- Category Validation ---
	if donation.CategoryID != "" {
		category, err := d.DonationCategoryRepo.FindByID(ctx, donation.CategoryID)
		if err != nil || category == nil || !category.Enabled {
			span.AddEvent("Donation category not found or disabled", trace.WithAttributes(
				attribute.String("error", localization.ErrorDonationCategoryNotFound.Code),
				attribute.String("id", id),
				attribute.String("category_id", donation.CategoryID),
			))
			return errors.New(localization.ErrorDonationCategoryNotFound.Code)
		}
		existingDonation.Category.CategoryName = category.CategoryName
	} else {
		donation.CategoryID = existingDonation.Category.ID
	}

	// --- Company Validation ---
	if donation.CompanyID != "" {
		company, err := d.DonationCompanyRepo.FindByID(ctx, donation.CompanyID)
		if err != nil || company == nil || !company.Enabled {
			span.AddEvent("Donation company not found or disabled", trace.WithAttributes(
				attribute.String("error", localization.ErrorDonationCompanyNotFound.Code),
				attribute.String("id", id),
				attribute.String("company_id", donation.CompanyID),
			))
			return errors.New(localization.ErrorDonationCompanyNotFound.Code)
		}
		existingDonation.Company.CompanyName = company.CompanyName
	} else {
		donation.CompanyID = existingDonation.Company.ID
	}

	// --- Cover Image Handling ---
	coverImageURL := existingDonation.CoverImage
	if donation.CoverImage != nil {
		var objectkey string
		if existingDonation.CoverImage != "" {
			objectkey = path.Base(existingDonation.CoverImage)
		}

		url, err := lib.UploadFileToMinio(
			ctx,
			d.minio,
			d.bucketName,
			donation.CoverImage,
			string(constants.DonationFolderName),
			*d.cfg,
			objectkey,
			d.logger,
		)
		if err != nil {
			span.AddEvent("Failed to upload cover image", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return err
		}
		coverImageURL = url
	}

	NewdonationImages := existingDonation.DonationImages
	if len(donation.RemovedImages) > 0 {
		toRemove := make(map[string]struct{}, len(donation.RemovedImages))
		for _, id := range donation.RemovedImages {
			toRemove[id] = struct{}{}
		}
		filtered := NewdonationImages[:0]
		for _, img := range NewdonationImages {
			if _, ok := toRemove[img.ID]; !ok {
				filtered = append(filtered, img)
			}
		}
		NewdonationImages = filtered
	}

	// --- Add New Images ---
	d.logger.Infof("Processing %d new donation images", len(donation.DonationImages))
	for i, fileHeader := range donation.DonationImages {
		var objectkey string
		if len(existingDonation.DonationImages) > 0 && existingDonation.DonationImages[i].PhotoURL != "" {
			objectkey = path.Base(existingDonation.DonationImages[i].PhotoURL)
		}

		url, err := lib.UploadFileToMinio(
			ctx,
			d.minio,
			d.bucketName,
			fileHeader,
			string(constants.DonationFolderName),
			*d.cfg,
			objectkey,
			d.logger,
		)
		if err != nil {
			d.logger.Errorf("Failed to upload donation image: %v", err)
			span.AddEvent("Failed to upload donation image", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return err
		}

		NewdonationImages = append(NewdonationImages, types.DonationImage{
			ID:       bson.NewObjectID().Hex(),
			PhotoURL: url,
		})

		d.logger.Infof("Successfully uploaded donation image: %s", url)
	}

	// --- Map Update Data ---
	updateData := core.MapDonationUpdate(
		id,
		existingModel,
		donation,
		coverImageURL,
		NewdonationImages,
		existingDonation.Company,
		existingDonation.Category,
	)

	// --- CPS Action ---
	cpsAction := lib.CpsModelBuilder(
		id,
		makerData,
		existingDonation,
		updateData,
		string(constants.RequestUpdateDonation),
		constants.UPDATE,
	)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil
}

func (d *Donation) UpdateDonationImage(ctx context.Context, id string, image dto.DonationImageUpdateRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateDonationImage", "Donation", "UpdateDonationImage")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberRequired.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	existingDonation, err := d.DonationRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find donation", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if existingDonation == nil {
		span.AddEvent("Donation not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorFileNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorFileNotFound.Code)
	}

	if image.Image == nil {
		span.AddEvent("Image required", trace.WithAttributes(
			attribute.String("error", localization.ErrorImageRequired.Code),
			attribute.String("id", id),
		))
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
		span.AddEvent("Image not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorFileNotFound.Code),
			attribute.String("id", id),
			attribute.String("image_id", image.ImageID),
		))
		return errors.New(localization.ErrorFileNotFound.Code)
	}

	var objectkey string
	if len(existingDonation.DonationImages) > 0 && existingDonation.DonationImages[0].PhotoURL != "" {
		objectkey = path.Base(existingDonation.DonationImages[0].PhotoURL)
	}

	url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, image.Image, string(constants.DonationFolderName), *d.cfg, objectkey, d.logger)
	if err != nil {
		span.AddEvent("Failed to upload image", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	updateData := dto.DonationImageUpdateCPSRequest{
		ID:       id,
		ImageID:  image.ImageID,
		PhotoURL: url,
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonation, updateData, string(constants.RequestUpdateDonationImage), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil
}

func (d *Donation) DeleteDonationImage(ctx context.Context, id string, imageID string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteDonationImage", "Donation", "DeleteDonationImage")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberRequired.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	existingDonation, err := d.DonationRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find donation", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if existingDonation == nil {
		span.AddEvent("Donation not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorFileNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorFileNotFound.Code)
	}

	updateData := dto.DonationImageDeleteCPSRequest{
		ImageIDToDelete: imageID,
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonation, updateData, string(constants.RequestDeleteDonationImage), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil
}

func (d *Donation) AddDonationImage(ctx context.Context, id string, image dto.DonationRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "AddDonationImage", "Donation", "AddDonationImage")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberRequired.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	existingDonation, err := d.DonationRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find donation", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if existingDonation == nil {
		span.AddEvent("Donation not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorFileNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorFileNotFound.Code)
	}

	if len(image.DonationImages) == 0 {
		span.AddEvent("Image required", trace.WithAttributes(
			attribute.String("error", localization.ErrorImageRequired.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorImageRequired.Code)
	}

	donationImages := make([]types.DonationImage, 0)
	for i, img := range image.DonationImages {
		var objectkey string
		if len(existingDonation.DonationImages) > 0 && existingDonation.DonationImages[i].PhotoURL != "" {
			objectkey = path.Base(existingDonation.DonationImages[i].PhotoURL)
		}

		url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, img, string(constants.DonationFolderName), *d.cfg, objectkey, d.logger)
		if err != nil {
			span.AddEvent("Failed to upload image", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return err
		}
		donationImages = append(donationImages, types.DonationImage{
			ID:        bson.NewObjectID().Hex(),
			PhotoURL:  url,
			CreatedAt: time.Now(),
		})
	}
	//todo unique dto for the donation image array
	updateData := dto.DonationImageAddRequest{
		DonationImages: donationImages,
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonation, updateData, string(constants.RequestAddDonationImage), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil
}

func (d *Donation) EnableDonation(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableDonation", "Donation", "EnableDonation")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}

	existingDonation, err := d.DonationRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find donation", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if existingDonation == nil {
		span.AddEvent("Donation not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorFileNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorFileNotFound.Code)
	}
	if existingDonation.Enabled {
		span.AddEvent("Donation already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorDonationAlreadyEnabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorDonationAlreadyEnabled.Code)
	}

	updateData := *existingDonation
	updateData.Enabled = true
	updateData.LastModifiedAt = time.Now().Format(time.RFC3339)

	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonation, updateData, string(constants.RequestEnableDonation), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil
}

func (d *Donation) DisableDonation(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DisableDonation", "Donation", "DisableDonation")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberRequired.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	existingDonation, err := d.DonationRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find donation", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if existingDonation == nil {
		span.AddEvent("Donation not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorFileNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorFileNotFound.Code)
	}
	if !existingDonation.Enabled {
		span.AddEvent("Donation already disabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorDonationAlreadyDisabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorDonationAlreadyDisabled.Code)
	}

	updateData := *existingDonation
	updateData.Enabled = false
	updateData.LastModifiedAt = time.Now().Format(time.RFC3339)

	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonation, updateData, string(constants.RequestDisableDonation), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil
}

func (d *Donation) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "Donation", "Authorize")
	defer span.End()

	if action.ActionStatus != constants.Approved {
		d.logger.Errorf("Tried to authorize service action without cps action approval")
		span.AddEvent("CPS action status invalid", trace.WithAttributes(
			attribute.String("error", localization.ErrorCPSActionStatusInvalid.Code),
			attribute.String("unique_id", action.UniqueId),
		))
		return nil, errors.New(localization.ErrorCPSActionStatusInvalid.Code)
	}
	var donationCPS *dto.DonationCPSRequest
	bindErr := core.BindAction(action.CurrentAction, &donationCPS)
	if bindErr != nil {
		d.logger.Errorf("failed to bind current action to donation: %v", bindErr)
		span.AddEvent("Failed to bind current action", trace.WithAttributes(
			attribute.String("error", bindErr.Error()),
			attribute.String("unique_id", action.UniqueId),
		))
		return nil, errors.New(localization.ErrorCPSActionFailed.Code)
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
			CompanyID:           donationCPS.Company.ID,
			CategoryID:          donationCPS.Category.ID,
			Title:               donationCPS.Title,
			IsFeatured:          donationCPS.IsFeatured,
			Target:              donationCPS.Target,
			DonationDescription: donationCPS.DonationDescription,
			Enabled:             donationCPS.Enabled,
			StartDate:           core.ParseTime(donationCPS.StartDate),
			EndDate:             core.ParseTime(donationCPS.EndDate),
		}

		donationImages := make([]types.DonationImage, len(donationCPS.DonationImages))
		for i, img := range donationCPS.DonationImages {
			donationImages[i] = types.DonationImage{
				ID:        img.ID,
				PhotoURL:  img.PhotoURL,
				CreatedAt: img.CreatedAt,
			}
		}

		updateData := core.MapDonationUpdate(
			action.UniqueId,
			existingModel,
			updateRequest,
			donationCPS.CoverImage,
			donationImages,
			existingDonation.Company,
			existingDonation.Category,
		)
		donationModel := core.MapToDonationModel(&updateData)
		donationModel.CurrentAmount = existingDonation.CurrentAmount

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to update donation: %v", err)
			span.AddEvent("Failed to update donation", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestUpdateDonationImage):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("Failed to find donation for image update: %v", err)
			span.AddEvent("Failed to find donation for image update", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		if existingDonation == nil {
			span.AddEvent("Donation not found", trace.WithAttributes(
				attribute.String("error", localization.ErrorFileNotFound.Code),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}

		// Bind the current action to get the image update data
		var imageUpdateData *dto.DonationImageUpdateCPSRequest
		bindErr := core.BindAction(action.CurrentAction, &imageUpdateData)
		if bindErr != nil {
			d.logger.Errorf("failed to bind current action to donation image update: %v", bindErr)
			span.AddEvent("Failed to bind current action to donation image update", trace.WithAttributes(
				attribute.String("error", bindErr.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
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
					CreatedAt: img.CreatedAt,
				})
				imageUpdated = true
			} else {
				// Keep other images unchanged
				updatedImages = append(updatedImages, types.DonationImage{
					ID:        img.ID,
					PhotoURL:  img.PhotoURL,
					CreatedAt: img.CreatedAt,
				})
			}
		}

		if !imageUpdated {
			d.logger.Errorf("Image ID %s not found in donation %s", imageUpdateData.ImageID, action.UniqueId)
			span.AddEvent("Image ID not found", trace.WithAttributes(
				attribute.String("error", localization.ErrorFileNotFound.Code),
				attribute.String("unique_id", action.UniqueId),
				attribute.String("image_id", imageUpdateData.ImageID),
			))
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		existingModel := core.ConvertDonationListResponseToModel(existingDonation)
		updateRequest := core.ConvertDonationListResponseToRequest(existingDonation)

		donationImages := make([]types.DonationImage, len(updatedImages))
		for i, img := range updatedImages {
			donationImages[i] = types.DonationImage{
				ID:        img.ID,
				PhotoURL:  img.PhotoURL,
				CreatedAt: img.CreatedAt,
			}
		}

		updateData := core.MapDonationUpdate(
			action.UniqueId,
			existingModel,
			updateRequest,
			existingDonation.CoverImage,
			donationImages,
			existingDonation.Company,
			existingDonation.Category,
		)
		donationModel := core.MapToDonationModel(&updateData)
		donationModel.CurrentAmount = existingDonation.CurrentAmount

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to update donation image: %v", err)
			span.AddEvent("Failed to update donation image", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestDeleteDonationImage):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("Failed to find donation for image deletion: %v", err)
			span.AddEvent("Failed to find donation for image deletion", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		if existingDonation == nil {
			span.AddEvent("Donation not found", trace.WithAttributes(
				attribute.String("error", localization.ErrorFileNotFound.Code),
				attribute.String("unique_id", action.UniqueId),
			))
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
					CreatedAt: img.CreatedAt,
				})
			}
		}

		donationImages := make([]types.DonationImage, len(updatedImages))
		for i, img := range updatedImages {
			donationImages[i] = types.DonationImage{
				ID:        img.ID,
				PhotoURL:  img.PhotoURL,
				CreatedAt: img.CreatedAt,
			}
		}

		updateData := core.MapDonationUpdate(
			action.UniqueId,
			existingModel,
			updateRequest,
			existingDonation.CoverImage,
			donationImages,
			existingDonation.Company,
			existingDonation.Category,
		)
		donationModel := core.MapToDonationModel(&updateData)
		donationModel.CurrentAmount = existingDonation.CurrentAmount

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to delete donation image: %v", err)
			span.AddEvent("Failed to delete donation image", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestAddDonationImage):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("Failed to find donation for image addition: %v", err)
			span.AddEvent("Failed to find donation for image addition", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		if existingDonation == nil {
			span.AddEvent("Donation not found", trace.WithAttributes(
				attribute.String("error", localization.ErrorFileNotFound.Code),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}

		existingModel := core.ConvertDonationListResponseToModel(existingDonation)
		updateRequest := core.ConvertDonationListResponseToRequest(existingDonation)

		updatedImages := make([]types.DonationImage, len(existingDonation.DonationImages))
		for i, img := range existingDonation.DonationImages {
			updatedImages[i] = types.DonationImage{
				ID:        img.ID,
				PhotoURL:  img.PhotoURL,
				CreatedAt: img.CreatedAt,
			}
		}

		for _, newImg := range donationCPS.DonationImages {
			updatedImages = append(updatedImages, types.DonationImage{
				ID:        newImg.ID,
				PhotoURL:  newImg.PhotoURL,
				CreatedAt: time.Now(),
			})
		}

		donationImages := make([]types.DonationImage, len(updatedImages))
		for i, img := range updatedImages {
			donationImages[i] = types.DonationImage{
				ID:        img.ID,
				PhotoURL:  img.PhotoURL,
				CreatedAt: img.CreatedAt,
			}
		}

		updateData := core.MapDonationUpdate(
			action.UniqueId,
			existingModel,
			updateRequest,
			existingDonation.CoverImage,
			donationImages,
			existingDonation.Company,
			existingDonation.Category,
		)
		donationModel := core.MapToDonationModel(&updateData)
		donationModel.CurrentAmount = existingDonation.CurrentAmount

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to add donation images: %v", err)
			span.AddEvent("Failed to add donation images", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestEnableDonation):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("Failed to find donation for image addition: %v", err)
			span.AddEvent("Failed to find donation", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		if existingDonation == nil {
			span.AddEvent("Donation not found", trace.WithAttributes(
				attribute.String("error", localization.ErrorFileNotFound.Code),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		tempval := true
		existingModel := core.ConvertDonationListResponseToModel(existingDonation)
		updateRequest := dto.DonationRequest{
			Enabled: &tempval,
		}
		updateData := core.MapDonationUpdate(
			action.UniqueId,
			existingModel,
			updateRequest,
			existingDonation.CoverImage,
			existingDonation.DonationImages,
			existingDonation.Company,
			existingDonation.Category,
		)
		donationModel := core.MapToDonationModel(&updateData)
		donationModel.CurrentAmount = existingDonation.CurrentAmount

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to enable donation: %v", err)
			span.AddEvent("Failed to enable donation", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestDisableDonation):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("Failed to find donation for image addition: %v", err)
			span.AddEvent("Failed to find donation", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
		if existingDonation == nil {
			span.AddEvent("Donation not found", trace.WithAttributes(
				attribute.String("error", localization.ErrorFileNotFound.Code),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorFileNotFound.Code)
		}
		existingModel := core.ConvertDonationListResponseToModel(existingDonation)
		tempval := false
		updateRequest := dto.DonationRequest{
			Enabled: &tempval,
		}
		updateData := core.MapDonationUpdate(
			action.UniqueId,
			existingModel,
			updateRequest,
			existingDonation.CoverImage,
			existingDonation.DonationImages,
			existingDonation.Company,
			existingDonation.Category,
		)
		donationModel := core.MapToDonationModel(&updateData)
		donationModel.CurrentAmount = existingDonation.CurrentAmount

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("Failed to disable donation: %v", err)
			span.AddEvent("Failed to disable donation", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	default:
		d.logger.Errorf("Unsupported action requested: %s", action.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(action.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	d.logger.Infof("Service action authorization completed: %s", action.RequestAction)
	return action, nil
}
