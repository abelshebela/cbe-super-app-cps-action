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
	"path"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/config"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type DonationCompany struct {
	DonationCompanyRepo  storage.DonationCompanyRepository
	cpsService           service.CPSActionService
	logger               utils.Logger
	minio                *s3.Client
	bucketName           string
	cfg                  *config.VaultConfig
	minioEndPoint        string
	accountLookupService account_lookup.Account
}

func NewDonationCompanyService(client *mongo.Client, DonationCompanyRepo storage.DonationCompanyRepository, cpsAction service.CPSActionService, logger utils.Logger, minio *s3.Client,
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
	ok, err := core.CompanyNameExists(ctx, donationCompany.CompanyName, d.DonationCompanyRepo)
	if err != nil {
		return err
	}
	if ok {
		return errors.New(localization.ErrorCompanyNameAlreadyExists.Code)
	}

	ok, err = core.AccountNumberExists(ctx, donationCompany.AccountNumber, d.DonationCompanyRepo)
	if err != nil {
		return err
	}
	if ok {
		return errors.New(localization.ErrorAccountNumberAlreadyExists.Code)
	}

	accountDetail, err := core.ValidateAccountNumberWithExternalAPI(ctx, donationCompany.AccountNumber, d.accountLookupService)
	if err != nil {
		d.logger.Errorf("Account number validation failed: %v", err)
		return err
	}

	if donationCompany.CompanyLogo == nil {
		return errors.New(localization.ErrorLogoIsRequired.Code)
	}

	url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donationCompany.CompanyLogo, string(constants.CampanyLogo), *d.cfg, "", d.logger)
	if err != nil {
		return err
	}

	donationCompany.CompanyCode = "DON-COMPANY-" + local_util.UniqueIdGenerator()
	result := core.MapToDonationCompanyResponse(donationCompany, url)
	result.AccountHolderName = accountDetail.CustomerName
	cpsAction := lib.CpsModelBuilder("", makerData, "", result, string(constants.RequestCreateDonationCompany), constants.CREATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}
	return nil
}

func (d *DonationCompany) UpdateDonationCompany(ctx context.Context, id string, donationCompany dto.DonationCompanyRequest) (*model.DonationCompany, error) {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return nil, errors.New(localization.ErrorAccountNumberRequired.Code)
	}

	existingCompany, err := d.DonationCompanyRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingCompany == nil {
		return nil, errors.New(localization.ErrorFileNotFound.Code)
	}

	if err := core.CheckDataSimilarityAndValidation(ctx, donationCompany, existingCompany, d.DonationCompanyRepo, d.accountLookupService); err != nil {
		return nil, err
	}

	logoURL := existingCompany.CompanyLogo
	if donationCompany.CompanyLogo != nil {
		var objectkey string
		if existingCompany.CompanyLogo != "" {
			objectkey = path.Base(existingCompany.CompanyLogo)
		}

		url, err := lib.UploadFileToMinio(ctx, d.minio, d.bucketName, donationCompany.CompanyLogo, string(constants.CampanyLogo), *d.cfg, objectkey, d.logger)
		if err != nil {
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
			return nil, err
		}
		updateData.AccountHolderName = accountDetail.CustomerName
	} else {
		updateData.AccountNumber = existingCompany.AccountNumber
		updateData.AccountHolderName = existingCompany.AccountHolderName
	}

	cpsAction := lib.CpsModelBuilder(id, makerData, existingCompany, updateData, string(constants.RequestUpdateDonationCompany), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return nil, err
	}

	// Use core mapper to create return request
	return updateData, nil
}

func (d *DonationCompany) Authorize(ctx context.Context, action *model.CPSAction) (*model.CPSAction, error) {
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
	case string(constants.RequestEnableDonationCompany):

		err := d.DonationCompanyRepo.Update(ctx, action.UniqueId, donationCompoany)
		if err != nil {
			d.logger.Errorf("Failed to update donation company: %v", err)
			return nil, err
		}
	case string(constants.RequestDisableDonationCompany):

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

func (d *DonationCompany) AccountLookup(ctx context.Context, accountNumber string) (*model.AccountDetail, error) {
	accountDetail, err := core.ValidateAccountNumberWithExternalAPI(ctx, accountNumber, d.accountLookupService)
	if err != nil {
		d.logger.Errorf("Account number validation failed: %v", err)
		return nil, err
	}
	return accountDetail, nil
}

//MapToDonationCompany

func (d *DonationCompany) EnableDonationCompany(ctx context.Context, id string) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	existingDonationCompany, err := d.DonationCompanyRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingDonationCompany == nil {
		return errors.New(localization.ErrorFileNotFound.Code)
	}
	if existingDonationCompany.Enabled {
		return errors.New(localization.ErrorAlreadyEnabled.Code)
	}
	DonationCompany := core.MapToDonationCompany(existingDonationCompany, true)
	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonationCompany, DonationCompany, string(constants.RequestEnableDonationCompany), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}

	return nil

}

func (d *DonationCompany) DisableDonationCompany(ctx context.Context, id string) error {
	makerData := local_util.ExtractUserFromContext(ctx)
	if incomplet := local_util.IsIncomplete(makerData); incomplet {
		return errors.New(localization.ErrorIncompleteUserInfo.Code)
	}
	existingDonationCompany, err := d.DonationCompanyRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingDonationCompany == nil {
		return errors.New(localization.ErrorFileNotFound.Code)
	}
	if !existingDonationCompany.Enabled {
		return errors.New(localization.ErrorAlreadyDisabled.Code)
	}
	DonationCompany := core.MapToDonationCompany(existingDonationCompany, false)
	cpsAction := lib.CpsModelBuilder(id, makerData, existingDonationCompany, DonationCompany, string(constants.RequestDisableDonationCompany), constants.UPDATE)
	if err := d.cpsService.CreateCPSAction(ctx, &cpsAction); err != nil {
		return err
	}

	return nil

}
