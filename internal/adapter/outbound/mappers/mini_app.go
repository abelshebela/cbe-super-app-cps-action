package mappers

import (
	"fmt"

	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/adapter/outbound/model"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/internal/domain/miniapp"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func ToModelMiniApp(domain *miniapp.MiniApp) (*model.MiniApp, error) {
	cred := model.CredentialInformation{
		ID:            bson.NewObjectID(),
		Environment:   model.UatEnvironment,
		MerchantAppID: domain.Credential.MerchantAppID,
		FabricAppID:   domain.Credential.FabricAppID,
		ShortCode:     domain.Credential.ShortCode,
		AppSecret:     domain.Credential.AppSecret,
		PrivateKey:    domain.Credential.PrivateKey,
		PublicKey:     domain.Credential.PublicKey,
		Timestamp:     domain.Credential.Timestamp,
		Signature:     domain.Credential.Signature,
		MiniAppCode:   domain.Credential.MiniAppCode,
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
		BannerImage:       domain.BannerImage,
		CommisonGLAccount: domain.CommissionGLAccount,
		AppType:           string(domain.AppType),
		MerchantID:        domain.MerchantID,
		ProductCode:       productCodes,
		Credential:        cred,
		AppViewType:       string(domain.AppViewType),
		URL:               domain.URL,
		Stage:             string(domain.Stage),
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
	creds := miniapp.CredentialInformation{
		ID:            model.Credential.ID.Hex(),
		Environment:   miniapp.UatEnvironment,
		MerchantAppID: model.Credential.MerchantAppID,
		FabricAppID:   model.Credential.FabricAppID,
		ShortCode:     model.Credential.ShortCode,
		AppSecret:     model.Credential.AppSecret,
		PrivateKey:    model.Credential.PrivateKey,
		PublicKey:     model.Credential.PublicKey,
		MiniAppCode:   model.Credential.MiniAppCode,
		Signature:     model.Credential.Signature,
		Timestamp:     model.Credential.Timestamp,
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
		BannerImage:         model.BannerImage,
		CommissionGLAccount: model.CommisonGLAccount,
		AppType:             miniapp.AppType(model.AppType),
		MerchantID:          model.MerchantID,
		ProductCode:         productCodes,
		Credential:          creds,
		URL:                 model.URL,
		AppViewType:         miniapp.AppViewType(model.AppViewType),
		Stage:               miniapp.Stage(model.Stage),
		IsEventMiniApp:      model.IsEventMiniApp,
		IsThreeClick:        model.IsThreeClick,
		Enabled:             model.Enabled,
		IsDeleted:           model.IsDeleted,
		CreatedAt:           model.CreatedAt,
		LastModifiedAt:      model.LastModifiedAt,
		DeletedAt:           model.DeletedAt,
	}
}
