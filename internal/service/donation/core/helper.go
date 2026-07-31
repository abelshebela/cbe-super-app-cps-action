package core

import (
	donation_dto "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/dto/donation"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/localization"
	imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/constants/types"
	"github.com/abelshebela/cbe-super-app-cps-action/internal/storage"
	"github.com/abelshebela/cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"encoding/json"

	shared_types "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"
)

func nonEmptyBool(newVal *bool, oldVal bool) bool {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}

func ConvertToDonationImages(images []types.DonationImage) []shared_types.DonationImage {
	result := make([]shared_types.DonationImage, len(images))
	for i, img := range images {
		result[i] = shared_types.DonationImage{
			ID:        img.ID,
			PhotoURL:  img.PhotoURL,
			CreatedAt: img.CreatedAt,
		}
	}
	return result
}

func DonationTitleExists(ctx context.Context, title string, donationRepo storage.DonationRepository) (bool, error) {
	donations, err := donationRepo.FindAllWithPagination(ctx, types.Filter{
		Filters: map[string]any{"title": title},
	})
	if err != nil {
		return false, err
	}

	if len(donations.Data) != 0 {
		return true, nil
	}

	return false, nil
}

func IsDataSimilar(request donation_dto.DonationRequest, existing *imodel.DonationOracle) bool {
	if request.ServiceID != "" && request.ServiceID != existing.ServiceID {
		return false
	}
	if request.CompanyID != "" && request.CompanyID != existing.CompanyID {
		return false
	}
	if request.CategoryID != "" && request.CategoryID != existing.CategoryID {
		return false
	}
	if request.Title != "" && request.Title != existing.Title {
		return false
	}
	if request.DonationDescription != "" && request.DonationDescription != existing.DonationDescription {
		return false
	}
	if request.Target != "" && request.Target != existing.Target {
		return false
	}
	if request.IsFeatured != nil && *request.IsFeatured != existing.IsFeatured {
		return false
	}
	if request.Enabled != nil && *request.Enabled != existing.Enabled {
		return false
	}
	if !request.StartDate.IsZero() && !request.StartDate.Equal(existing.StartDate) {
		return false
	}
	if !request.EndDate.IsZero() && !request.EndDate.Equal(existing.EndDate) {
		return false
	}
	if request.CoverImage != nil {
		return false
	}
	if len(request.DonationImages) > 0 {
		return false
	}
	if len(request.RemovedImages) > 0 {
		return false
	}
	return true
}

func BindAction(source any, target any) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

func CheckDataSimilarityAndValidation(ctx context.Context, request donation_dto.DonationRequest, existing *imodel.DonationOracle, donationRepo storage.DonationRepository) error {
	if IsDataSimilar(request, existing) {
		return errors.New(localization.ErrorNoChangesToUpdate.Code)
	}
	if request.Title != "" && request.Title != existing.Title {
		ok, err := DonationTitleExists(ctx, request.Title, donationRepo)
		if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
			return err
		}
		if ok {
			return errors.New(localization.ErrorDonationTitleDuplicated.Code)
		}
	}

	return nil
}

func MapToDonationModel(cpsRequest *donation_dto.DonationCPSRequest) *imodel.DonationOracle {
	endTime, _ := time.Parse(time.RFC3339, cpsRequest.EndDate)
	startTime, _ := time.Parse(time.RFC3339, cpsRequest.StartDate)

	donationImages := make([]types.DonationImage, len(cpsRequest.DonationImages))
	for i, img := range cpsRequest.DonationImages {
		donationImages[i] = types.DonationImage{
			ID:        img.ID,
			PhotoURL:  img.PhotoURL,
			CreatedAt: img.CreatedAt,
		}
	}

	return &imodel.DonationOracle{
		DonationCode:        cpsRequest.DonationCode,
		ServiceID:           GetValueOrDefault(cpsRequest.Service.ID, cpsRequest.ServiceID),
		CompanyID:           cpsRequest.Company.ID,
		CategoryID:          cpsRequest.Category.ID,
		Title:               cpsRequest.Title,
		IsFeatured:          *cpsRequest.IsFeatured,
		Target:              cpsRequest.Target,
		CurrentAmount:       "0",
		DonationDescription: cpsRequest.DonationDescription,
		DonationImages:      ConvertToDonationImages(donationImages),
		CoverImage:          cpsRequest.CoverImage,
		EndDate:             endTime,
		StartDate:           startTime,
		Enabled:             *cpsRequest.Enabled,
		IsDeleted:           false,
		CreatedAt:           time.Now(),
		LastModifiedAt:      time.Now(),
	}
}

func GenerateDonationCode() string {
	return "DON-" + utils.UniqueIdGenerator()
}

func ParseTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t
}

func GetValueOrDefault(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}

func GetIntValueOrDefault(value, defaultValue int32) int32 {
	if value == 0 {
		return defaultValue
	}
	return value
}

func GetTimeValueOrDefault(value, defaultValue time.Time) time.Time {
	if value.IsZero() {
		return defaultValue
	}
	return value
}

func ConvertDonationListResponseToModel(donationResponse *donation_dto.DonationListResponse) *imodel.DonationOracle {
	donationImages := make([]types.DonationImage, len(donationResponse.DonationImages))
	for i, img := range donationResponse.DonationImages {
		donationImages[i] = types.DonationImage{
			ID:        img.ID,
			PhotoURL:  img.PhotoURL,
			CreatedAt: img.CreatedAt,
		}
	}

	startTime, _ := time.Parse(time.RFC3339, donationResponse.StartDate)
	endTime, _ := time.Parse(time.RFC3339, donationResponse.EndDate)
	createdAt, _ := time.Parse(time.RFC3339, donationResponse.CreatedAt)
	lastModifiedAt, _ := time.Parse(time.RFC3339, donationResponse.LastModifiedAt)

	return &imodel.DonationOracle{
		ID:                  donationResponse.ID,
		DonationCode:        donationResponse.DonationCode,
		ServiceID:           donationResponse.Service.ID,
		CompanyID:           donationResponse.Company.ID,
		CategoryID:          donationResponse.Category.ID,
		Title:               donationResponse.Title,
		IsFeatured:          donationResponse.IsFeatured,
		Target:              donationResponse.Target,
		DonationDescription: donationResponse.DonationDescription,
		DonationImages:      ConvertToDonationImages(donationImages),
		CoverImage:          donationResponse.CoverImage,
		StartDate:           startTime,
		EndDate:             endTime,
		Enabled:             donationResponse.Enabled,
		IsDeleted:           donationResponse.IsDeleted,
		CreatedAt:           createdAt,
		LastModifiedAt:      lastModifiedAt,
	}
}

func ConvertDonationListResponseToRequest(donationResponse *donation_dto.DonationListResponse) donation_dto.DonationRequest {
	return donation_dto.DonationRequest{
		CompanyID:           donationResponse.Company.ID,
		CategoryID:          donationResponse.Category.ID,
		Title:               donationResponse.Title,
		IsFeatured:          &donationResponse.IsFeatured,
		Target:              donationResponse.Target,
		DonationDescription: donationResponse.DonationDescription,
		StartDate:           ParseTime(donationResponse.StartDate),
		EndDate:             ParseTime(donationResponse.EndDate),
	}
}

func MapDonationUpdate(
	id string,
	existing *imodel.DonationOracle,
	update donation_dto.DonationRequest,
	coverImageURL string,
	donationImages []types.DonationImage,
	service donation_dto.Service,
	company donation_dto.Company,
	category donation_dto.Category,
) donation_dto.DonationCPSRequest {
	isFeatured := nonEmptyBool(update.IsFeatured, existing.IsFeatured)
	isEnabled := nonEmptyBool(update.Enabled, existing.Enabled)

	return donation_dto.DonationCPSRequest{
		ID:           id,
		DonationCode: existing.DonationCode,
		Service: donation_dto.Service{
			ID:          GetValueOrDefault(update.ServiceID, existing.ServiceID),
			ServiceName: service.ServiceName,
			ServiceKey:  service.ServiceKey,
		},
		Company: donation_dto.Company{
			ID:          GetValueOrDefault(update.CompanyID, existing.CompanyID),
			CompanyName: company.CompanyName,
			CompanyLogo: company.CompanyLogo,
		},
		Category: donation_dto.Category{
			ID:           GetValueOrDefault(update.CategoryID, existing.CategoryID),
			CategoryName: category.CategoryName,
			Icon:         category.Icon,
		},
		Title:               GetValueOrDefault(update.Title, existing.Title),
		IsFeatured:          &isFeatured,
		Target:              update.Target,
		DonationDescription: GetValueOrDefault(update.DonationDescription, existing.DonationDescription),
		DonationImages:      donationImages,
		CoverImage:          coverImageURL,
		StartDate:           GetTimeValueOrDefault(update.StartDate, existing.StartDate).Format(time.RFC3339),
		EndDate: func() string {
			if t := update.EndDate; !t.IsZero() {
				return t.Format(time.RFC3339)
			}
			return ""
		}(),
		Enabled: &isEnabled,
	}
}
