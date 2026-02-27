package core

import (
	donation_dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	imodel "cbe-super-app-cps-action/internal/constants/model"
	"encoding/json"

	// "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/model"
	// types "gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/entities/types"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func nonEmptyBool(newVal *bool, oldVal bool) bool {
	if newVal != nil {
		return *newVal
	}
	return oldVal
}

func MapToDonationResponse(donation *imodel.Donation) *donation_dto.DonationResponse {
	return &donation_dto.DonationResponse{
		DonationCode:        donation.DonationCode,
		CompanyID:           donation.CompanyID.Hex(),
		CategoryID:          donation.CategoryID.Hex(),
		Title:               donation.Title,
		IsFeatured:          donation.IsFeatured,
		Target:              donation.Target,
		DonationDescription: donation.DonationDescription,
		DonationImages:      ConvertToDonationImages(donation.DonationImages),
		CoverImage:          donation.CoverImage,
		EndDate:             donation.EndDate.Format(time.RFC3339),
		StartDate:           donation.StartDate.Format(time.RFC3339),
		Enabled:             donation.Enabled,
	}
}

func MapToDonationListResponse(donation *imodel.Donation) *donation_dto.DonationListResponse {
	return &donation_dto.DonationListResponse{
		ID:                  donation.ID.Hex(),
		DonationCode:        donation.DonationCode,
		Title:               donation.Title,
		IsFeatured:          donation.IsFeatured,
		Target:              donation.Target,
		DonationDescription: donation.DonationDescription,
		DonationImages:      ConvertToDonationImages(donation.DonationImages),
		CoverImage:          donation.CoverImage,
		EndDate:             donation.EndDate.Format(time.RFC3339),
		StartDate:           donation.StartDate.Format(time.RFC3339),
		IsDeleted:           donation.IsDeleted,
		CreatedAt:           donation.CreatedAt.Format(time.RFC3339),
		LastModifiedAt:      donation.LastModifiedAt.Format(time.RFC3339),
		Enabled:             donation.Enabled,
	}
}

func MapToDonationCPSRequest(donation *imodel.Donation) *donation_dto.DonationCPSRequest {
	return &donation_dto.DonationCPSRequest{
		DonationCode:        donation.DonationCode,
		CompanyID:           donation.CompanyID.Hex(),
		CategoryID:          donation.CategoryID.Hex(),
		Title:               donation.Title,
		IsFeatured:          &donation.IsFeatured,
		Target:              donation.Target,
		DonationDescription: donation.DonationDescription,
		DonationImages:      ConvertToDonationImages(donation.DonationImages),
		CoverImage:          donation.CoverImage,
		EndDate:             donation.EndDate.Format(time.RFC3339),
		StartDate:           donation.StartDate.Format(time.RFC3339),
		Enabled:             &donation.Enabled,
	}
}

func MapToDonation(donationCode, companyID, categoryID, title, donationDescription, coverImage, endDate, startDate string, isFeatured bool, target int32, enabled bool) *imodel.Donation {
	companyObjID, _ := bson.ObjectIDFromHex(companyID)
	categoryObjID, _ := bson.ObjectIDFromHex(categoryID)

	endTime, _ := time.Parse(time.RFC3339, endDate)
	startTime, _ := time.Parse(time.RFC3339, startDate)

	return &imodel.Donation{
		DonationCode:        donationCode,
		CompanyID:           companyObjID,
		CategoryID:          categoryObjID,
		Title:               title,
		IsFeatured:          isFeatured,
		Target:              string(target),
		DonationDescription: donationDescription,
		CoverImage:          coverImage,
		EndDate:             endTime,
		StartDate:           startTime,
		Enabled:             enabled,
		IsDeleted:           false,
		CreatedAt:           time.Now(),
		LastModifiedAt:      time.Now(),
	}
}

func ConvertToDonationImages(images []types.DonationImage) []types.DonationImage {
	result := make([]types.DonationImage, len(images))
	for i, img := range images {
		result[i] = types.DonationImage{
			ID:        img.ID,
			PhotoURL:  img.PhotoURL,
			CreatedAt: img.CreatedAt,
		}
	}
	return result
}

func DonationTitleExists(ctx context.Context, title string, donationRepo storage.DonationRepository) (bool, error) {
	donations, err := donationRepo.FindAllWithPagination(ctx, types.Filter{
		Search: title,
	})
	if err != nil && err.Error() != localization.ErrorResourceNotFound.Code {
		return false, err
	}

	if donations.Data != nil {
		return true, nil
	}

	return false, nil
}

func IsDataSimilar(request donation_dto.DonationRequest, existing *imodel.Donation) bool {
	if request.Title != "" && request.Title != existing.Title {
		return false
	}
	if request.DonationDescription != "" && request.DonationDescription != existing.DonationDescription {
		return false
	}
	if request.Target != "" && request.Target != existing.Target {
		return false
	}
	if request.IsFeatured != &existing.IsFeatured {
		return false
	}
	if *request.Enabled != existing.Enabled {
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
	return true
}

func BindAction(source any, target any) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}

func CheckDataSimilarityAndValidation(ctx context.Context, request donation_dto.DonationRequest, existing *imodel.Donation, donationRepo storage.DonationRepository) error {
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

func MapToDonationModel(cpsRequest *donation_dto.DonationCPSRequest) *imodel.Donation {
	companyObjID, _ := bson.ObjectIDFromHex(cpsRequest.Company.ID)
	categoryObjID, _ := bson.ObjectIDFromHex(cpsRequest.Category.ID)

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

	return &imodel.Donation{
		DonationCode:        cpsRequest.DonationCode,
		CompanyID:           companyObjID,
		CategoryID:          categoryObjID,
		Title:               cpsRequest.Title,
		IsFeatured:          *cpsRequest.IsFeatured,
		Target:              cpsRequest.Target,
		CurrentAmount:       "0",
		DonationDescription: cpsRequest.DonationDescription,
		DonationImages:      donationImages,
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

func ConvertDonationListResponseToModel(donationResponse *donation_dto.DonationListResponse) *imodel.Donation {
	companyObjID, _ := bson.ObjectIDFromHex(donationResponse.Company.ID)
	categoryObjID, _ := bson.ObjectIDFromHex(donationResponse.Category.ID)

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

	donationID, _ := bson.ObjectIDFromHex(donationResponse.ID)

	return &imodel.Donation{
		ID:                  donationID,
		DonationCode:        donationResponse.DonationCode,
		CompanyID:           companyObjID,
		CategoryID:          categoryObjID,
		Title:               donationResponse.Title,
		IsFeatured:          donationResponse.IsFeatured,
		Target:              donationResponse.Target,
		DonationDescription: donationResponse.DonationDescription,
		DonationImages:      donationImages,
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
	existing *imodel.Donation,
	update donation_dto.DonationRequest,
	coverImageURL string,
	donationImages []types.DonationImage,
	company donation_dto.Company,
	category donation_dto.Category,
) donation_dto.DonationCPSRequest {
	isFeatured := nonEmptyBool(update.IsFeatured, existing.IsFeatured)
	isEnabled := nonEmptyBool(update.Enabled, existing.Enabled)

	return donation_dto.DonationCPSRequest{
		ID:           id,
		DonationCode: existing.DonationCode,
		Company: donation_dto.Company{
			ID: GetValueOrDefault(update.CompanyID, existing.CompanyID.Hex()),
		},
		Category: donation_dto.Category{
			ID: GetValueOrDefault(update.CategoryID, existing.CategoryID.Hex()),
		},
		// CompanyID:           GetValueOrDefault(update.CompanyID, existing.CompanyID.Hex()),
		// CompanyName:         companyName,
		// CategoryID:          GetValueOrDefault(update.CategoryID, existing.CategoryID.Hex()),
		// CategoryName:        categoryName,
		Title:               GetValueOrDefault(update.Title, existing.Title),
		IsFeatured:          &isFeatured,
		Target:              GetValueOrDefault(update.Target, existing.Target),
		DonationDescription: GetValueOrDefault(update.DonationDescription, existing.DonationDescription),
		DonationImages:      donationImages,
		CoverImage:          coverImageURL,
		StartDate:           GetTimeValueOrDefault(update.StartDate, existing.StartDate).Format(time.RFC3339),
		EndDate:             GetTimeValueOrDefault(update.EndDate, existing.EndDate).Format(time.RFC3339),
		Enabled:             &isEnabled,
	}
}
