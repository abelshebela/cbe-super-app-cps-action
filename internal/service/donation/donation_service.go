package donation

import (
	"cbe-super-app-cps-action/internal/constants"
	"encoding/csv"
	"os"
	"path"

	dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	imodel "cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	core "cbe-super-app-cps-action/internal/service/donation/core"
	"cbe-super-app-cps-action/internal/storage"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"

	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	// types "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"

	"fmt"

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
	ServicesRepo         storage.ServicesRepository
	cpsService           service.CPSActionService
	logger               utils.Logger
	minio                *s3.Client
	bucketName           string
	cfg                  *config.VaultConfig
	minioEndPoint        string
}

func NewDonationService(client *mongo.Client, DonationRepo storage.DonationRepository, DonationCategoryRepo storage.DonationCategoryRepository, DonationCompanyRepo storage.DonationCompanyRepository, servicesRepo storage.ServicesRepository, cpsAction service.CPSActionService, logger utils.Logger, minio *s3.Client,
	bucketName string,
	cfg *config.VaultConfig,
	minioEndPoint string,
) service.DonationService {

	return &Donation{
		DonationRepo:         DonationRepo,
		DonationCategoryRepo: DonationCategoryRepo,
		DonationCompanyRepo:  DonationCompanyRepo,
		ServicesRepo:         servicesRepo,
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
		d.logger.Errorf("[Donation][Create] Failed to check donation title existence: %v", err)
		return err
	}
	if ok {
		span.AddEvent("Donation title duplicated", trace.WithAttributes(
			attribute.String("error", localization.ErrorDonationTitleDuplicated.Code),
			attribute.String("title", donation.Title),
		))
		d.logger.Warnf("[Donation][Create] Donation title duplicated: %s", donation.Title)
		return errors.New(localization.ErrorDonationTitleDuplicated.Code)
	}

	if donation.ServiceID == "" {
		span.AddEvent("Service id missing", trace.WithAttributes(
			attribute.String("error", localization.ErrorServiceNotFound.Code),
		))
		return errors.New(localization.ErrorServiceNotFound.Code)
	}
	service, err := d.ServicesRepo.FindByID(ctx, donation.ServiceID)
	if err != nil {
		span.AddEvent("Service id not found", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("service_id", donation.ServiceID),
		))
		return errors.New(localization.ErrorServiceNotFound.Code)
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
	d.logger.Infof("[DonationSvc][Create] processing %d images", len(donation.DonationImages))
	for _, img := range donation.DonationImages {
		url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, img, string(constants.DonationFolderName), *d.cfg, "", d.logger)
		if err != nil {
			d.logger.Errorf("[DonationSvc][Create] upload image err: %v", err)
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

		d.logger.Infof("[DonationSvc][Create] uploaded image: %s", url)
	}

	donationCode := core.GenerateDonationCode()
	d.logger.Infof("[DonationSvc][Create] building cps target: %d, images: %d", donation.Target, len(donationImages))
	tempval := false
	result := dto.DonationCPSRequest{
		DonationCode: donationCode,
		ServiceID:    donation.ServiceID,
		Service: dto.Service{
			ServiceName: service.ServiceName,
			ServiceKey:  service.ServiceKey,
		},
		Company: dto.Company{
			ID:          donation.CompanyID,
			CompanyName: company.CompanyName,
			CompanyLogo: company.CompanyLogo,
			Enabled:     company.Enabled,
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
		DonationImages:      donationImages,
		CoverImage:          coverImageURL,
		EndDate:             donation.EndDate.Format(time.RFC3339),
		StartDate:           donation.StartDate.Format(time.RFC3339),
		Enabled:             &tempval,
	}
	d.logger.Infof("[DonationSvc][Create] cps request target: %d, images: %d", result.Target, len(result.DonationImages))

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

	prevData := *existingDonation
	existingModel := &imodel.DonationOracle{
		DonationCode:        existingDonation.DonationCode,
		CompanyID:           existingDonation.Company.ID,
		CategoryID:          existingDonation.Category.ID,
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

	// --- Service Validation ---
	serviceForCPS := existingDonation.Service
	if donation.ServiceID != "" {
		svc, err := d.ServicesRepo.FindByID(ctx, donation.ServiceID)
		if err != nil {
			span.AddEvent("Service id not found", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
				attribute.String("service_id", donation.ServiceID),
			))
			return errors.New(localization.ErrorServiceNotFound.Code)
		}
		serviceForCPS = dto.Service{
			ID:          svc.ID,
			ServiceName: svc.ServiceName,
			ServiceKey:  svc.ServiceKey,
		}
	} else {
		donation.ServiceID = existingDonation.Service.ID
	}

	// --- Category Validation ---
	categoryForCPS := existingDonation.Category
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
		categoryForCPS = dto.Category{
			ID:           category.ID,
			CategoryName: category.CategoryName,
			Icon:         category.Icon,
		}
	} else {
		donation.CategoryID = existingDonation.Category.ID
	}

	// --- Company Validation ---
	companyForCPS := existingDonation.Company
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
		companyForCPS = dto.Company{
			ID:          company.ID,
			CompanyName: company.CompanyName,
			CompanyLogo: company.CompanyLogo,
			Enabled:     company.Enabled,
		}
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

	// Deep copy to avoid mutating existingDonation.DonationImages (used as prevAction in CPS)
	NewdonationImages := make([]types.DonationImage, len(existingDonation.DonationImages))
	copy(NewdonationImages, existingDonation.DonationImages)
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
	d.logger.Infof("[DonationSvc][Update] processing %d new images", len(donation.DonationImages))
	for _, fileHeader := range donation.DonationImages {
		var objectkey string

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
			d.logger.Errorf("[DonationSvc][Update] upload image err: %v", err)
			span.AddEvent("Failed to upload donation image", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return err
		}

		NewdonationImages = append(NewdonationImages, types.DonationImage{
			ID:        bson.NewObjectID().Hex(),
			PhotoURL:  url,
			CreatedAt: time.Now(),
		})

		d.logger.Infof("[DonationSvc][Update] uploaded image: %s", url)
	}

	// --- Map Update Data ---
	updateData := core.MapDonationUpdate(
		id,
		existingModel,
		donation,
		coverImageURL,
		NewdonationImages,
		serviceForCPS,
		companyForCPS,
		categoryForCPS,
	)

	// --- CPS Action ---
	cpsAction := lib.CpsModelBuilder(
		id,
		makerData,
		&prevData,
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

	company, err := d.DonationCompanyRepo.FindByID(ctx, existingDonation.Company.ID)
	if err != nil {
		span.AddEvent("Donation company not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorDonationCompanyNotFound.Code),
			attribute.String("id", id),
			attribute.String("company_id", existingDonation.Company.ID),
		))
		return errors.New(localization.ErrorDonationCompanyNotFound.Code)
	}
	if !company.Enabled {
		span.AddEvent("Donation company not enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorDonationCompanyNotEnabled.Code),
			attribute.String("id", id),
			attribute.String("company_id", existingDonation.Company.ID),
		))
		return errors.New(localization.ErrorDonationCompanyNotEnabled.Code)
	}

	category, err := d.DonationCategoryRepo.FindByID(ctx, existingDonation.Category.ID)
	if err != nil {
		span.AddEvent("Donation category not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorDonationCategoryNotFound.Code),
			attribute.String("id", id),
			attribute.String("category_id", existingDonation.Category.ID),
		))
		return errors.New(localization.ErrorDonationCategoryNotFound.Code)
	}
	if !category.Enabled {
		span.AddEvent("Donation category not enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorDonationCategoryNotEnabled.Code),
			attribute.String("id", id),
			attribute.String("category_id", existingDonation.Category.ID),
		))
		return errors.New(localization.ErrorDonationCategoryNotEnabled.Code)
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
		d.logger.Errorf("[DonationSvc][Authorize] not approved")
		span.AddEvent("CPS action status invalid", trace.WithAttributes(
			attribute.String("error", localization.ErrorCPSActionStatusInvalid.Code),
			attribute.String("unique_id", action.UniqueId),
		))
		return nil, errors.New(localization.ErrorCPSActionStatusInvalid.Code)
	}
	var donationCPS *dto.DonationCPSRequest
	bindErr := core.BindAction(action.CurrentAction, &donationCPS)
	if bindErr != nil {
		d.logger.Errorf("[DonationSvc][Authorize] bind err: %v", bindErr)
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
			d.logger.Errorf("[DonationSvc][Authorize] create err: %v", err)
			return nil, err
		}

	case string(constants.RequestUpdateDonation):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("[DonationSvc][Authorize] find for update err: %v", err)
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
			existingDonation.Service,
			existingDonation.Company,
			existingDonation.Category,
		)
		donationModel := core.MapToDonationModel(&updateData)
		donationModel.CurrentAmount = existingDonation.CurrentAmount

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("[DonationSvc][Authorize] update err: %v", err)
			span.AddEvent("Failed to update donation", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestUpdateDonationImage):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("[DonationSvc][Authorize] find for img update err: %v", err)
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
			d.logger.Errorf("[DonationSvc][Authorize] bind img update err: %v", bindErr)
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
			d.logger.Errorf("[DonationSvc][Authorize] img %s not found in %s", imageUpdateData.ImageID, action.UniqueId)
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
			existingDonation.Service,
			existingDonation.Company,
			existingDonation.Category,
		)
		donationModel := core.MapToDonationModel(&updateData)
		donationModel.CurrentAmount = existingDonation.CurrentAmount

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("[DonationSvc][Authorize] update img err: %v", err)
			span.AddEvent("Failed to update donation image", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestDeleteDonationImage):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("[DonationSvc][Authorize] find for img delete err: %v", err)
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
			existingDonation.Service,
			existingDonation.Company,
			existingDonation.Category,
		)
		donationModel := core.MapToDonationModel(&updateData)
		donationModel.CurrentAmount = existingDonation.CurrentAmount

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("[DonationSvc][Authorize] delete img err: %v", err)
			span.AddEvent("Failed to delete donation image", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestAddDonationImage):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("[DonationSvc][Authorize] find for img add err: %v", err)
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
			existingDonation.Service,
			existingDonation.Company,
			existingDonation.Category,
		)
		donationModel := core.MapToDonationModel(&updateData)
		donationModel.CurrentAmount = existingDonation.CurrentAmount

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("[DonationSvc][Authorize] add images err: %v", err)
			span.AddEvent("Failed to add donation images", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestEnableDonation):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("[DonationSvc][Authorize] find for enable err: %v", err)
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
			existingDonation.Service,
			existingDonation.Company,
			existingDonation.Category,
		)
		donationModel := core.MapToDonationModel(&updateData)
		donationModel.CurrentAmount = existingDonation.CurrentAmount

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("[DonationSvc][Authorize] enable err: %v", err)
			span.AddEvent("Failed to enable donation", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestDisableDonation):
		existingDonation, err := d.DonationRepo.FindByID(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("[DonationSvc][Authorize] find for disable err: %v", err)
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
			existingDonation.Service,
			existingDonation.Company,
			existingDonation.Category,
		)
		donationModel := core.MapToDonationModel(&updateData)
		donationModel.CurrentAmount = existingDonation.CurrentAmount

		err = d.DonationRepo.Update(ctx, action.UniqueId, donationModel)
		if err != nil {
			d.logger.Errorf("[DonationSvc][Authorize] disable err: %v", err)
			span.AddEvent("Failed to disable donation", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	case string(constants.RequestDeleteDonation):
		err := d.DonationRepo.Delete(ctx, action.UniqueId)
		if err != nil {
			d.logger.Errorf("[DonationSvc][Authorize] delete err: %v", err)
			span.AddEvent("Failed to delete donation", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

	default:
		d.logger.Errorf("[DonationSvc][Authorize] unsupported action: %s", action.RequestAction)
		span.AddEvent("Unsupported action", trace.WithAttributes(
			attribute.String("error", localization.ErrorUnsupportedAction.Code),
			attribute.String("request_action", string(action.RequestAction)),
		))
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	d.logger.Infof("[DonationSvc][Authorize] completed: %s", action.RequestAction)
	return action, nil
}

func (d *Donation) ExportDonationData(ctx context.Context, startDate, endDate time.Time, fileType string) (string, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "ExportDonationData", "Donation", "ExportDonationData")
	defer span.End()

	if endDate.Before(startDate) {
		span.AddEvent("Invalid date range", trace.WithAttributes(
			attribute.String("start_date", startDate.Format(time.RFC3339)),
			attribute.String("end_date", endDate.Format(time.RFC3339)),
		))
		return "", errors.New(localization.DonationDataExportedError.Code)
	}

	// 1. Create temp CSV file
	tmpFile, err := os.CreateTemp("", "donations_*.csv")
	if err != nil {
		span.AddEvent("Failed to create temp file", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return "", errors.New(localization.DonationDataExportedError.Code)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	writer := csv.NewWriter(tmpFile)

	// Write UTF-8 BOM so Excel correctly recognizes comma-delimited columns
	if _, err := tmpFile.Write([]byte{0xEF, 0xBB, 0xBF}); err != nil {
		span.AddEvent("Failed to write BOM", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return "", errors.New(localization.DonationDataExportedError.Code)
	}

	// 2. Write CSV header
	if err := writer.Write(donationCSVHeader()); err != nil {
		span.AddEvent("Failed to write CSV header", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return "", errors.New(localization.DonationDataExportedError.Code)
	}

	// 3. Stream donations from repository and write rows
	rowCount := 0
	err = d.DonationRepo.StreamByDateRange(ctx, startDate, endDate,
		func(donation *imodel.DonationOracle) error {
			rowCount++
			return writer.Write(buildDonationRow(donation))
		},
	)
	if err != nil {
		span.AddEvent("Failed to stream donation data", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return "", errors.New(localization.DonationDataExportedError.Code)
	}

	// 4. Check if any data was found
	if rowCount == 0 {
		span.AddEvent("No donation data found in the specified date range")
		return "", errors.New(localization.DonationDataNotFoundInDateRange.Code)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		span.AddEvent("Failed to flush CSV writer", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return "", errors.New(localization.DonationDataExportedError.Code)
	}

	// 5. Upload CSV to MinIO
	objectKey := fmt.Sprintf(
		"exports/donations/donations_%s_to_%s_%d.csv",
		startDate.Format("20060102"),
		endDate.Format("20060102"),
		time.Now().Unix(),
	)

	if _, err := tmpFile.Seek(0, 0); err != nil {
		span.AddEvent("Failed to seek temp file", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return "", errors.New(localization.DonationDataExportedError.Code)
	}

	stat, err := tmpFile.Stat()
	if err != nil {
		span.AddEvent("Failed to stat temp file", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return "", errors.New(localization.DonationDataExportedError.Code)
	}

	publicURL, err := lib.UploadCSVToMinio(ctx, d.minio, d.bucketName, tmpFile, stat.Size(), *d.cfg, objectKey, d.logger)
	if err != nil {
		span.AddEvent("Failed to upload to MinIO", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("object_key", objectKey),
		))
		return "", errors.New(localization.DonationDataExportedError.Code)
	}

	span.AddEvent("Donation data exported successfully", trace.WithAttributes(
		attribute.String("url", publicURL),
		attribute.String("object_key", objectKey),
		attribute.Int("row_count", rowCount),
	))

	return publicURL, nil
}

func donationCSVHeader() []string {
	return []string{
		"ID",
		"DonationCode",
		"Title",
		"CompanyID",
		"CategoryID",
		"Target",
		"CurrentAmount",
		"IsFeatured",
		"Enabled",
		"StartDate",
		"EndDate",
	}
}

func buildDonationRow(donation *imodel.DonationOracle) []string {
	return []string{
		donation.ID,
		donation.DonationCode,
		donation.Title,
		donation.CompanyID,
		donation.CategoryID,
		donation.Target,
		donation.CurrentAmount,
		fmt.Sprintf("%v", donation.IsFeatured),
		fmt.Sprintf("%v", donation.Enabled),
		donation.StartDate.Format(time.RFC3339),
		donation.EndDate.Format(time.RFC3339),
	}
}

func (d *Donation) DeleteDonation(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DeleteDonation", "Donation", "DeleteDonation")
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
	if existingDonation.IsDeleted {
		span.AddEvent("Donation already deleted", trace.WithAttributes(
			attribute.String("error", localization.ErrorAlreadyDeleted.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorAlreadyDeleted.Code)
	}

	updateData := *existingDonation
	updateData.IsDeleted = true
	updateData.LastModifiedAt = time.Now().Format(time.RFC3339)

	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonation, updateData, string(constants.RequestDeleteDonation), constants.DELETE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil
}
