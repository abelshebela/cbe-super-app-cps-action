package donation_company

import (
	"cbe-super-app-cps-action/internal/constants"
	dto "cbe-super-app-cps-action/internal/constants/dto/donation_company"
	"cbe-super-app-cps-action/internal/constants/lib"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/service"
	core "cbe-super-app-cps-action/internal/service/donation_company/core"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/internal/storage/external_call/account_lookup"
	local_util "cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"fmt"

	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DonationCompany struct {
	DonationCompanyRepo  storage.DonationCompanyRepository
	cpsService           service.CPSActionService
	logger               utils.Logger
	minio                config.MinioClientInterface
	bucketName           string
	cfg                  *config.VaultConfig
	minioEndPoint        string
	accountLookupService account_lookup.Account
}

func NewDonationCompanyService(client *mongo.Client, DonationCompanyRepo storage.DonationCompanyRepository, cpsAction service.CPSActionService, logger utils.Logger, minio config.MinioClientInterface,
	bucketName string,
	cfg *config.VaultConfig,
	minioEndPoint string,
	accountLookupService account_lookup.Account,
) service.DonationCompanyService {

	return &DonationCompany{
		DonationCompanyRepo:  DonationCompanyRepo,
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
	return d.DonationCompanyRepo.FindAllWithPagination(ctx, *filterParams)
}

func (d *DonationCompany) FetchDonationCompanyByID(ctx context.Context, id string) (*dto.DonationCompanyListResponse, error) {
	return d.DonationCompanyRepo.FindByID(ctx, id)
}

func (d *DonationCompany) CreateDonationCompany(ctx context.Context, donationCompany dto.DonationCompanyRequest) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	// Check if company name already exists
	ok, err := core.CompanyNameExists(ctx, donationCompany.CompanyName, d.DonationCompanyRepo)
	if err != nil {
		return err
	}
	if ok {
		return errors.New(localization.ErrorCompanyNameAlreadyExists.Code)
	}

	// Check if account number already exists
	ok, err = core.AccountNumberExists(ctx, donationCompany.AccountNumber, d.DonationCompanyRepo)
	if err != nil {
		return err
	}
	if ok {
		return errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}

	// Validate account number with external API
	if _, err := core.ValidateAccountNumberWithExternalAPI(ctx, donationCompany.AccountNumber, d.accountLookupService); err != nil {
		d.logger.Errorf("Account number validation failed: %v", err)
		return err
	}

	// Check if logo is provided
	if donationCompany.CompanyLogo == nil {
		return errors.New(localization.ErrorLogoIsRequired.Code)
	}

	// Upload logo to MinIO
	url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donationCompany.CompanyLogo, string(constants.CampanyLogo), d.minioEndPoint, d.logger)
	if err != nil {
		return err
	}

	d.logger.Infof("Logo uploaded successfully: %s", url)
	// Use core mapper to create response
	result := core.MapToDonationCompanyResponse(donationCompany, url)

	cpsAction := lib.CpsModelBuilder("", makerData, "", result, string(constants.RequestCreateDonationCompany), constants.CREATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}
	return nil
}

func (d *DonationCompany) UpdateDonationCompany(ctx context.Context, id string, donationCompany dto.DonationCompanyRequest) (dto.DonationCompanyRequest, error) {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return donationCompany, errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	// Check if donation company exists
	existingCompany, err := d.DonationCompanyRepo.FindByID(ctx, id)
	if err != nil {
		return donationCompany, err
	}
	if existingCompany == nil {
		return donationCompany, errors.New(localization.ErrorFileNotFound.Code)
	}

	if err := core.CheckDataSimilarityAndValidation(ctx, donationCompany, existingCompany, d.DonationCompanyRepo, d.accountLookupService); err != nil {
		return donationCompany, err
	}

	logoURL := existingCompany.CompanyLogo
	if donationCompany.CompanyLogo != nil {
		url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donationCompany.CompanyLogo, string(constants.CampanyLogo), d.minioEndPoint, d.logger)
		if err != nil {
			return donationCompany, err
		}
		logoURL = url
	}

	updateData := core.MapToDonationCompanyCPSRequest(id, donationCompany, logoURL)

	// Use existing values if not provided in update
	if donationCompany.CompanyName == "" {
		updateData.CompanyName = existingCompany.CompanyName
	}
	if donationCompany.AccountNumber == "" {
		updateData.AccountNumber = existingCompany.AccountNumber
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, existingCompany, updateData, string(constants.RequestUpdateDonationCompany), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return donationCompany, err
	}

	// Use core mapper to create return request
	return core.MapToDonationCompanyRequest(updateData), nil
}

func (d *DonationCompany) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
	fmt.Println("Donation Company Authorize called")
	if action.ActionStatus != constants.Approved {
		d.logger.Errorf("Tried to authorize service action without cps action approval")
		return nil, errors.New(localization.ErrorCPSActionStatusInvalid.Code)
	}

	donationCompoany, err := local_util.JsonUnmarshal[model.DonationCompany](action.CurrentAction)
	if err != nil {
		return nil, errors.New(localization.ErrorUnexpectedError.Code)
	}

	switch action.RequestAction {
	case string(constants.RequestCreateDonationCompany):

		err := d.DonationCompanyRepo.Create(ctx, donationCompoany)
		if err != nil {
			d.logger.Errorf("Failed to create donation company: %v", err)
			return nil, err
		}

	case string(constants.RequestUpdateDonationCompany):

		err := d.DonationCompanyRepo.Update(ctx, action.UniqueId, donationCompoany)
		if err != nil {
			d.logger.Errorf("Failed to update donation company: %v", err)
			return nil, err
		}

	default:
		d.logger.Errorf("Unsupported action requested: %s", action.RequestAction)
		return nil, errors.New(localization.ErrorUnsupportedAction.Code)
	}

	d.logger.Infof("Service action authorization completed: %s", action.RequestAction)
	return action, nil
}

// Account lookup end point

func (d *DonationCompany) AccountLookup(ctx context.Context, accountNumber string) (*model.AccountInfo, error) {
	account, err := core.ValidateAccountNumberWithExternalAPI(ctx, accountNumber, d.accountLookupService)
	if err != nil {
		d.logger.Errorf("Account number validation failed: %v", err)
		return nil, err
	}
	return account, nil
}
