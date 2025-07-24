package mappers

import (
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	domain "github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp_merchant"
	"go.mongodb.org/mongo-driver/v2/bson"
)


func ToMiniAppMerchantModel(domain *domain.MiniAppMerchant) (*model.MiniAppMerchant, error) {
	var objectID bson.ObjectID
	if domain.ID != "" {
		objID, err := bson.ObjectIDFromHex(domain.ID)
		if err != nil {
			return nil, fmt.Errorf("INVALID_ID")
		}
		objectID = objID
	}

	return &model.MiniAppMerchant{
		ID:           objectID,
		Code:         domain.Code,
		MerchantName: domain.MerchantName,
		MerchantType: domain.MerchantType,
		KYC: model.KYC{
			Status: string(domain.KYC.Status),
			Representative: model.KYCInformation{
				Name:  domain.KYC.Representative.Name,
				Email: domain.KYC.Representative.Email,
				Phone: domain.KYC.Representative.Phone,
			},
		},
		BankAccountNumber: domain.BankAccountNumber,
		Branches:          convertBranchesDomainToModel(domain.Branches),
		Email:             domain.Email,
		PhoneNumber:       domain.PhoneNumber,
		MiniAppIDs:        domain.MiniAppIDs,
		Enabled:           domain.Enabled,
		IsDeleted:         domain.IsDeleted,
		CreatedAt:         domain.CreatedAt,
		LastUpdatedAt:     domain.LastModifiedAt,
		DeletedAt:         domain.DeletedAt,
	}, nil
}

func ToMiniAppMerchantDomain(model *model.MiniAppMerchant) *domain.MiniAppMerchant {
	return &domain.MiniAppMerchant{
		ID:           model.ID.Hex(),
		Code:         model.Code,
		MerchantName: model.MerchantName,
		MerchantType: model.MerchantType,
		KYC: domain.KYC{
			Status: domain.KYCStatus(model.KYC.Status),
			Representative: domain.KYCInformation{
				Name:  model.KYC.Representative.Name,
				Email: model.KYC.Representative.Email,
				Phone: model.KYC.Representative.Phone,
			},
		},
		BankAccountNumber: model.BankAccountNumber,
		Branches:          convertBranchesModelToDomain(model.Branches),
		Email:             model.Email,
		PhoneNumber:       model.PhoneNumber,
		MiniAppIDs:        model.MiniAppIDs,
		Enabled:           model.Enabled,
		IsDeleted:         model.IsDeleted,
		CreatedAt:         model.CreatedAt,
		LastModifiedAt:    model.LastUpdatedAt,
		DeletedAt:         model.DeletedAt,
	}
}

func convertBranchesDomainToModel(branches []domain.BranchInformation) []model.BranchInformation {
	result := make([]model.BranchInformation, len(branches))
	for i, b := range branches {
		result[i] = model.BranchInformation{
			BranchCode:          b.BranchCode,
			BranchName:          b.BranchName,
			BranchAddress:       b.BranchAddress,
			BranchOwner:         b.BranchOwner,
			BranchAccountNumber: b.BranchAccountNumber,
		}
	}
	return result
}

func convertBranchesModelToDomain(branches []model.BranchInformation) []domain.BranchInformation {
	result := make([]domain.BranchInformation, len(branches))
	for i, b := range branches {
		result[i] = domain.BranchInformation{
			BranchCode:          b.BranchCode,
			BranchName:          b.BranchName,
			BranchAddress:       b.BranchAddress,
			BranchOwner:         b.BranchOwner,
			BranchAccountNumber: b.BranchAccountNumber,
		}
	}
	return result
}
