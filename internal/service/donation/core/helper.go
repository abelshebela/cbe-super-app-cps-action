package core

import (
	donation_dto "cbe-super-app-cps-action/internal/constants/dto/donation"
	"cbe-super-app-cps-action/internal/constants/localization"
	"cbe-super-app-cps-action/internal/constants/model"
	"cbe-super-app-cps-action/internal/constants/types"
	"cbe-super-app-cps-action/internal/storage"
	"cbe-super-app-cps-action/pkgs/utils"
	"context"
	"errors"
	"time"

	"encoding/json"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func MapToDonationResponse(donation *model.Donation) *donation_dto.DonationResponse {
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

func MapToDonationListResponse(donation *model.Donation) *donation_dto.DonationListResponse {
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

func MapToDonationCPSRequest(donation *model.Donation) *donation_dto.DonationCPSRequest {
	return &donation_dto.DonationCPSRequest{
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

func MapToDonation(donationCode, companyID, categoryID, title, donationDescription, coverImage, endDate, startDate string, isFeatured bool, target int32, enabled bool) *model.Donation {
	companyObjID, _ := bson.ObjectIDFromHex(companyID)
	categoryObjID, _ := bson.ObjectIDFromHex(categoryID)

	endTime, _ := time.Parse(time.RFC3339, endDate)
	startTime, _ := time.Parse(time.RFC3339, startDate)

	return &model.Donation{
		DonationCode:        donationCode,
		CompanyID:           companyObjID,
		CategoryID:          categoryObjID,
		Title:               title,
		IsFeatured:          isFeatured,
		Target:              target,
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

func ConvertToDonationImages(images []types.DonationImage) []donation_dto.DonationImage {
	result := make([]donation_dto.DonationImage, len(images))
	for i, img := range images {
		result[i] = donation_dto.DonationImage{
			ID:        img.ID,
			PhotoURL:  img.PhotoURL,
			CreatedAt: img.CreatedAt.Format(time.RFC3339),
		}
	}
	return result
}

func DonationTitleExists(ctx context.Context, title string, donationRepo storage.DonationRepository) (bool, error) {
	donations, err := donationRepo.FindAllWithPagination(ctx, types.Filter{
		Search: title,
	})
	if err != nil && err != mongo.ErrNoDocuments {
		return false, errors.New(localization.ErrorDonationLookupFailed.Code)
	}

	if donations.Data != nil {
		return true, nil
	}

	return false, nil
}

func IsDataSimilar(request donation_dto.DonationRequest, existing *model.Donation) bool {
	if request.Title != "" && request.Title != existing.Title {
		return false
	}
	if request.DonationDescription != "" && request.DonationDescription != existing.DonationDescription {
		return false
	}
	if request.Target != 0 && request.Target != existing.Target {
		return false
	}
	if request.IsFeatured != existing.IsFeatured {
		return false
	}
	if request.Enabled != existing.Enabled {
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

func CheckDataSimilarityAndValidation(ctx context.Context, request donation_dto.DonationRequest, existing *model.Donation, donationRepo storage.DonationRepository) error {
	if IsDataSimilar(request, existing) {
		return errors.New(localization.ErrorNoChangesToUpdate.Code)
	}

	if request.Title != "" && request.Title != existing.Title {
		ok, err := DonationTitleExists(ctx, request.Title, donationRepo)
		if err != nil {
			return err
		}
		if ok {
			return errors.New(localization.ErrorDonationTitleDuplicated.Code)
		}
	}

	return nil
}

func MapToDonationModel(cpsRequest *donation_dto.DonationCPSRequest) *model.Donation {
	companyObjID, _ := bson.ObjectIDFromHex(cpsRequest.CompanyID)
	categoryObjID, _ := bson.ObjectIDFromHex(cpsRequest.CategoryID)

	endTime, _ := time.Parse(time.RFC3339, cpsRequest.EndDate)
	startTime, _ := time.Parse(time.RFC3339, cpsRequest.StartDate)

	donationImages := make([]types.DonationImage, len(cpsRequest.DonationImages))
	for i, img := range cpsRequest.DonationImages {
		createdAt, _ := time.Parse(time.RFC3339, img.CreatedAt)
		donationImages[i] = types.DonationImage{
			ID:        img.ID,
			PhotoURL:  img.PhotoURL,
			CreatedAt: createdAt,
		}
	}

	return &model.Donation{
		DonationCode:        cpsRequest.DonationCode,
		CompanyID:           companyObjID,
		CategoryID:          categoryObjID,
		Title:               cpsRequest.Title,
		IsFeatured:          cpsRequest.IsFeatured,
		Target:              cpsRequest.Target,
		DonationDescription: cpsRequest.DonationDescription,
		DonationImages:      donationImages,
		CoverImage:          cpsRequest.CoverImage,
		EndDate:             endTime,
		StartDate:           startTime,
		Enabled:             cpsRequest.Enabled,
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

func ConvertDonationListResponseToModel(donationResponse *donation_dto.DonationListResponse) *model.Donation {
	companyObjID, _ := bson.ObjectIDFromHex(donationResponse.Company.ID)
	categoryObjID, _ := bson.ObjectIDFromHex(donationResponse.Category.ID)

	donationImages := make([]types.DonationImage, len(donationResponse.DonationImages))
	for i, img := range donationResponse.DonationImages {
		createdAt, _ := time.Parse(time.RFC3339, img.CreatedAt)
		donationImages[i] = types.DonationImage{
			ID:        img.ID,
			PhotoURL:  img.PhotoURL,
			CreatedAt: createdAt,
		}
	}

	startTime, _ := time.Parse(time.RFC3339, donationResponse.StartDate)
	endTime, _ := time.Parse(time.RFC3339, donationResponse.EndDate)
	createdAt, _ := time.Parse(time.RFC3339, donationResponse.CreatedAt)
	lastModifiedAt, _ := time.Parse(time.RFC3339, donationResponse.LastModifiedAt)

	donationID, _ := bson.ObjectIDFromHex(donationResponse.ID)

	return &model.Donation{
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
		IsFeatured:          donationResponse.IsFeatured,
		Target:              donationResponse.Target,
		DonationDescription: donationResponse.DonationDescription,
		StartDate:           ParseTime(donationResponse.StartDate),
		EndDate:             ParseTime(donationResponse.EndDate),
	}
}

func MapDonationUpdate(id string, existing *model.Donation, update donation_dto.DonationRequest, coverImageURL string, donationImages []donation_dto.DonationImage) donation_dto.DonationCPSRequest {
	return donation_dto.DonationCPSRequest{
		ID:                  id,
		DonationCode:        existing.DonationCode,
		CompanyID:           GetValueOrDefault(update.CompanyID, existing.CompanyID.Hex()),
		CategoryID:          GetValueOrDefault(update.CategoryID, existing.CategoryID.Hex()),
		Title:               GetValueOrDefault(update.Title, existing.Title),
		IsFeatured:          update.IsFeatured,
		Target:              GetIntValueOrDefault(update.Target, existing.Target),
		DonationDescription: GetValueOrDefault(update.DonationDescription, existing.DonationDescription),
		DonationImages:      donationImages,
		CoverImage:          coverImageURL,
		StartDate:           GetTimeValueOrDefault(update.StartDate, existing.StartDate).Format(time.RFC3339),
		EndDate:             GetTimeValueOrDefault(update.EndDate, existing.EndDate).Format(time.RFC3339),
		Enabled:             existing.Enabled,
	}
}
