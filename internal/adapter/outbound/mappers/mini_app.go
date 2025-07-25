package mappers

import (
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ToModelMiniApp(domain *miniapp.MiniApp) (*model.MiniApp, error) {
	creds := make([]model.CredentialInformation, len(domain.Credential))
	for i, cred := range domain.Credential {
		creds[i] = model.CredentialInformation{
			ID:            cred.ID,
			Environment:   model.EnvironmentType(cred.Environment),
			MerchantAppID: cred.MerchantAppID,
			FabricAppID:   cred.FabricAppID,
			ShortCode:     cred.ShortCode,
			AppSecret:     cred.AppSecret,
			PrivateKey:    cred.PrivateKey,
			PublicKey:     cred.PublicKey,
		}
	}

	productCodes := make([]model.ProductCode, len(domain.ProductCode))
	for i, pc := range domain.ProductCode {
		productCodes[i] = model.ProductCode{
			ID:             pc.ID,
			BranchType:     model.BranchType(pc.BranchType),
			ProductCode:    pc.ProductCode,
			VATCode:        pc.VATCode,
			ServiceFeeCode: pc.VATCode,
		}
	}

	var objectID bson.ObjectID
	if domain.ID != "" {
		objID, err := bson.ObjectIDFromHex(domain.ID)

		if err != nil {
			return nil, fmt.Errorf("INVALID_ID")
		}
		objectID = objID

	} else {
		objectID = bson.NewObjectID()
	}

	return &model.MiniApp{
		ID:                objectID,
		AppName:           domain.AppName,
		AppIcon:           domain.AppIcon,
		CommisonGLAccount: domain.CommissionGLAccount,
		AppType:           string(domain.AppType),
		MerchantID:        domain.MerchantID,
		ProductCode:       productCodes,
		Credential:        creds,
		IsEventMiniApp:    domain.IsEventMiniApp,
		IsThreeClick:      domain.IsThreeClick,
		Enabled:           domain.Enabled,
		IsDeleted:         domain.IsDeleted,
		CreatedAt:         domain.CreatedAt,
		LastModifiedAt:    domain.LastModifiedAt,
		DeletedAt:         domain.DeletedAt,
	}, nil
}

func ToDomainMiniApp(model model.MiniApp) miniapp.MiniApp {
	creds := make([]miniapp.CredentialInformation, len(model.Credential))
	for i, cred := range model.Credential {
		creds[i] = miniapp.CredentialInformation{
			ID:            cred.ID,
			Environment:   miniapp.EnvironmentType(cred.Environment),
			MerchantAppID: cred.MerchantAppID,
			FabricAppID:   cred.FabricAppID,
			ShortCode:     cred.ShortCode,
			AppSecret:     cred.AppSecret,
			PrivateKey:    cred.PrivateKey,
			PublicKey:     cred.PublicKey,
		}
	}

	productCodes := make([]miniapp.ProductCode, len(model.ProductCode))
	for i, pc := range model.ProductCode {
		productCodes[i] = miniapp.ProductCode{
			ID:             pc.ID,
			BranchType:     miniapp.BranchType(pc.BranchType),
			ProductCode:    pc.ProductCode,
			VATCode:        pc.VATCode,
			ServiceFeeCode: pc.ServiceFeeCode,
		}
	}

	return miniapp.MiniApp{
		ID:                  model.ID.Hex(),
		AppName:             model.AppName,
		AppIcon:             model.AppIcon,
		CommissionGLAccount: model.CommisonGLAccount,
		AppType:             miniapp.AppType(model.AppType),
		MerchantID:          model.MerchantID,
		ProductCode:         productCodes,
		Credential:          creds,
		IsEventMiniApp:      model.IsEventMiniApp,
		IsThreeClick:        model.IsThreeClick,
		Enabled:             model.Enabled,
		IsDeleted:           model.IsDeleted,
		CreatedAt:           model.CreatedAt,
		LastModifiedAt:      model.LastModifiedAt,
		DeletedAt:           model.DeletedAt,
	}
}
