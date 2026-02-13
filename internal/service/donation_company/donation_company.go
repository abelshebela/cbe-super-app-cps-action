package donation_company

import (
	"cbe-super-app-cps-action/internal/constants"
	dto "cbe-super-app-cps-action/internal/constants/dto/donation_company"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	core "cbe-super-app-cps-action/internal/service/donation_company/core"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"path"
	"time"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	donation_model "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/donation"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type DonationCompany struct {
	DonationCompanyRepo  storage.DonationCompanyRepository
	DonationRepo         storage.DonationRepository
	cpsService           service.CPSActionService
	logger               utils.Logger
	minio                *s3.Client
	bucketName           string
	cfg                  *config.VaultConfig
	minioEndPoint        string
	accountLookupService account_lookup.Account
}

func NewDonationCompanyService(client *mongo.Client, DonationCompanyRepo storage.DonationCompanyRepository, DonationRepo storage.DonationRepository, cpsAction service.CPSActionService, logger utils.Logger, minio *s3.Client,
	bucketName string,
	cfg *config.VaultConfig,
	minioEndPoint string,
	accountLookupService account_lookup.Account,
) service.DonationCompanyService {

	return &DonationCompany{
		DonationCompanyRepo:  DonationCompanyRepo,
		DonationRepo:         DonationRepo,
		cpsService:           cpsAction,
		logger:               logger,
		minio:                minio,
		bucketName:           bucketName,
		cfg:                  cfg,
		minioEndPoint:        minioEndPoint,
		accountLookupService: accountLookupService,
	}
}

func (d *DonationCompany) FetchDonationCompany(ctx context.Context, filterParams *types.Filter) (*types.PaginatedResponse[[]dto.DonationCompanyListResponse], error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchDonationCompany", "DonationCompany", "FetchDonationCompany")
	defer span.End()

	result, err := d.DonationCompanyRepo.FindAllWithPagination(ctx, *filterParams)
	if err != nil {
		span.AddEvent("Failed to fetch donation companies", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return nil, err
	}
	return result, nil
}

func (d *DonationCompany) FetchDonationCompanyByID(ctx context.Context, id string) (*dto.DonationCompanyListResponse, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "FetchDonationCompanyByID", "DonationCompany", "FetchDonationCompanyByID")
	defer span.End()

	result, err := d.DonationCompanyRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to fetch donation company", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	return result, nil
}

func (d *DonationCompany) CreateDonationCompany(ctx context.Context, donationCompany dto.DonationCompanyRequest) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "CreateDonationCompany", "DonationCompany", "CreateDonationCompany")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberRequired.Code),
		))
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}
	ok, err := core.CompanyNameExists(ctx, donationCompany.CompanyName, d.DonationCompanyRepo)
	if err != nil {
		span.AddEvent("Failed to check company name existence", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("company_name", donationCompany.CompanyName),
		))
		return err
	}
	if ok {
		span.AddEvent("Company name already exists", trace.WithAttributes(
			attribute.String("error", localization.ErrorCompanyNameAlreadyExists.Code),
			attribute.String("company_name", donationCompany.CompanyName),
		))
		return errors.New(localization.ErrorCompanyNameAlreadyExists.Code)
	}

	if err = core.CheckIfAccountExists(ctx, donationCompany.AccountNumber, d.DonationCompanyRepo); err != nil {
		span.AddEvent("Account already exists", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("account_number", donationCompany.AccountNumber),
		))
		return err
	}

	accountDetail, err := core.ValidateAccountNumberWithExternalAPI(ctx, donationCompany.AccountNumber, d.accountLookupService)
	if err != nil {
		d.logger.Errorf("Account number validation failed: %v", err)
		span.AddEvent("Account number validation failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("account_number", donationCompany.AccountNumber),
		))
		return err
	}

	if donationCompany.CompanyLogo == nil {
		span.AddEvent("Logo is required", trace.WithAttributes(
			attribute.String("error", localization.ErrorLogoIsRequired.Code),
		))
		return errors.New(localization.ErrorLogoIsRequired.Code)
	}

	url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donationCompany.CompanyLogo, string(constants.DonationCompanyFolderName), *d.cfg, "", d.logger)
	if err != nil {
		span.AddEvent("Failed to upload logo", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}

	donationCompany.CompanyCode = "DON-COMPANY-" + local_util.UniqueIdGenerator()
	result := core.MapToDonationCompanyResponse(donationCompany, url)
	result.AccountHolderName = accountDetail.CustomerName
	cpsAction := lib.CpsModelBuilder("", makerData, "", result, string(constants.RequestCreateDonationCompany), constants.CREATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
		))
		return err
	}
	return nil
}

func (d *DonationCompany) UpdateDonationCompany(ctx context.Context, id string, donationCompany dto.DonationCompanyRequest) (*donation_model.DonationCompany, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "UpdateDonationCompany", "DonationCompany", "UpdateDonationCompany")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorAccountNumberRequired.Code),
			attribute.String("id", id),
		))
		return nil, errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	existingCompany, err := d.DonationCompanyRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find donation company", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}
	if existingCompany == nil {
		span.AddEvent("Donation company not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorFileNotFound.Code),
			attribute.String("id", id),
		))
		return nil, errors.New(localization.ErrorFileNotFound.Code)
	}

	if err := core.CheckDataSimilarityAndValidation(ctx, donationCompany, existingCompany, d.DonationCompanyRepo, d.accountLookupService); err != nil {
		span.AddEvent("Data similarity validation failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}

	logoURL := existingCompany.CompanyLogo
	if donationCompany.CompanyLogo != nil {
		var objectkey string
		if existingCompany.CompanyLogo != "" {
			objectkey = path.Base(existingCompany.CompanyLogo)
		}

		url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donationCompany.CompanyLogo, string(constants.DonationCompanyFolderName), *d.cfg, objectkey, d.logger)
		if err != nil {
			span.AddEvent("Failed to upload logo", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
			))
			return nil, err
		}
		logoURL = url
	}

	updateData := core.MapToDonationCompanyonUpdateCPSRequest(id, *existingCompany, donationCompany, logoURL)

	// Use existing values if not provided in update
	if donationCompany.CompanyName == "" {
		updateData.CompanyName = existingCompany.CompanyName
	}
	if donationCompany.AccountNumber != "" {
		accountDetail, err := core.ValidateAccountNumberWithExternalAPI(ctx, donationCompany.AccountNumber, d.accountLookupService)
		if err != nil {
			span.AddEvent("Account number validation failed", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("id", id),
				attribute.String("account_number", donationCompany.AccountNumber),
			))
			return nil, err
		}
		updateData.AccountHolderName = accountDetail.CustomerName
	} else {
		updateData.AccountNumber = existingCompany.AccountNumber
		updateData.AccountHolderName = existingCompany.AccountHolderName
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, existingCompany, updateData, string(constants.RequestUpdateDonationCompany), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return nil, err
	}

	// Use core mapper to create return request
	return updateData, nil
}

func (d *DonationCompany) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "Authorize", "DonationCompany", "Authorize")
	defer span.End()

	if action.ActionStatus != constants.Approved {
		d.logger.Errorf("Tried to authorize service action without cps action approval")
		span.AddEvent("CPS action status invalid", trace.WithAttributes(
			attribute.String("error", localization.ErrorCPSActionStatusInvalid.Code),
			attribute.String("unique_id", action.UniqueId),
		))
		return nil, errors.New(localization.ErrorCPSActionStatusInvalid.Code)
	}

	donationCompoany, err := local_util.JsonUnmarshal[donation_model.DonationCompany](action.CurrentAction)
	if err != nil {
		span.AddEvent("Failed to unmarshal CurrentAction", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("unique_id", action.UniqueId),
		))
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	switch action.RequestAction {
	case string(constants.RequestCreateDonationCompany):
		err := d.DonationCompanyRepo.Create(ctx, donationCompoany)
		if err != nil {
			d.logger.Errorf("Failed to create donation company: %v", err)
			span.AddEvent("Failed to create donation company", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestUpdateDonationCompany):
		err := d.DonationCompanyRepo.Update(ctx, action.UniqueId, donationCompoany)
		if err != nil {
			d.logger.Errorf("Failed to update donation company: %v", err)
			span.AddEvent("Failed to update donation company", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestEnableDonationCompany):

		err := d.DonationCompanyRepo.Update(ctx, action.UniqueId, donationCompoany)
		if err != nil {
			d.logger.Errorf("Failed to update donation company: %v", err)
			span.AddEvent("Failed to enable donation company", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}
	case string(constants.RequestDisableDonationCompany):

		err := d.DonationCompanyRepo.Update(ctx, action.UniqueId, donationCompoany)
		if err != nil {
			d.logger.Errorf("Failed to update donation company: %v", err)
			span.AddEvent("Failed to disable donation company", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

		// Disable all donation related with this donating company
		obj, err := bson.ObjectIDFromHex(action.UniqueId)
		if err != nil {
			span.AddEvent("Failed to parse ObjectID", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, errors.New(localization.ErrorUnexpectedError.Code)
		}
		donations, err := d.DonationRepo.FindAllWithPagination(ctx, types.Filter{
			Filters: map[string]interface{}{"company_id": obj}})
		if err != nil {
			if err.Error() == "mongo: no documents in result" {
				return nil, nil
			}
			span.AddEvent("Failed to find donations", trace.WithAttributes(
				attribute.String("error", err.Error()),
				attribute.String("unique_id", action.UniqueId),
			))
			return nil, err
		}

		if donations.Data != nil {
			for _, donation := range donations.Data {
				donationModel := core.ConvertDonationListResponseToModel(&donation)
				donationModel.Enabled = false
				donationModel.LastModifiedAt = time.Now()

				err := d.DonationRepo.Update(ctx, donation.ID, donationModel)
				if err != nil {
					span.AddEvent("Failed to update donation", trace.WithAttributes(
						attribute.String("error", localization.ErrorFailedToUpdateDonation.Code),
						attribute.String("unique_id", action.UniqueId),
						attribute.String("donation_id", donation.ID),
					))
					return nil, errors.New(localization.ErrorFailedToUpdateDonation.Code)
				}
			}
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

// Account lookup end point

func (d *DonationCompany) AccountLookup(ctx context.Context, accountNumber string) (*model.AccountDetail, error) {
	ctx, span := local_util.TraceLogger(ctx, "service", "AccountLookup", "DonationCompany", "AccountLookup")
	defer span.End()

	accountDetail, err := core.ValidateAccountNumberWithExternalAPI(ctx, accountNumber, d.accountLookupService)
	if err != nil {
		d.logger.Errorf("Account number validation failed: %v", err)
		span.AddEvent("Account number validation failed", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("account_number", accountNumber),
		))
		return nil, err
	}
	return accountDetail, nil
}

//MapToDonationCompany

func (d *DonationCompany) EnableDonationCompany(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "EnableDonationCompany", "DonationCompany", "EnableDonationCompany")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	existingDonationCompany, err := d.DonationCompanyRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find donation company", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if existingDonationCompany == nil {
		span.AddEvent("Donation company not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorFileNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorFileNotFound.Code)
	}
	if existingDonationCompany.Enabled {
		span.AddEvent("Donation company already enabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorAlreadyEnabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorAlreadyEnabled.Code)
	}
	// DonationCompany := core.MapToDonationCompany(existingDonationCompany, true)

	currentData := *existingDonationCompany
	currentData.Enabled = true
	currentData.LastModifiedAt = time.Now().Format(time.RFC3339)

	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonationCompany, currentData, string(constants.RequestEnableDonationCompany), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil
}

func (d *DonationCompany) DisableDonationCompany(ctx context.Context, id string) error {
	ctx, span := local_util.TraceLogger(ctx, "service", "DisableDonationCompany", "DonationCompany", "DisableDonationCompany")
	defer span.End()

	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		span.AddEvent("Incomplete user data", trace.WithAttributes(
			attribute.String("error", localization.ErrorIncompleteUserInfo.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	existingDonationCompany, err := d.DonationCompanyRepo.FindByID(ctx, id)
	if err != nil {
		span.AddEvent("Failed to find donation company", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}
	if existingDonationCompany == nil {
		span.AddEvent("Donation company not found", trace.WithAttributes(
			attribute.String("error", localization.ErrorFileNotFound.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorFileNotFound.Code)
	}
	if !existingDonationCompany.Enabled {
		span.AddEvent("Donation company already disabled", trace.WithAttributes(
			attribute.String("error", localization.ErrorAlreadyDisabled.Code),
			attribute.String("id", id),
		))
		return errors.New(localization.ErrorAlreadyDisabled.Code)
	}

	currentData := *existingDonationCompany
	currentData.Enabled = false
	currentData.LastModifiedAt = time.Now().Format(time.RFC3339)

	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonationCompany, currentData, string(constants.RequestDisableDonationCompany), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		span.AddEvent("Failed to create CPS action", trace.WithAttributes(
			attribute.String("error", err.Error()),
			attribute.String("id", id),
		))
		return err
	}

	return nil

}
