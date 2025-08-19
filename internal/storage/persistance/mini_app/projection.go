package mini_app

import (
	"cbe-super-app-cps-action/internal/constants/model"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// MiniAppMapper maps a MiniApp model to a bson.M for updates
func MiniAppMapper(miniApp *model.MiniApp) bson.M {
	return bson.M{
		"$set": bson.M{
			"code":           miniApp.Code,
			"name":           miniApp.Name,
			"description":    miniApp.Description,
			"version":        miniApp.Version,
			"icon_url":       miniApp.IconURL,
			"banner_url":     miniApp.BannerURL,
			"category":       miniApp.Category,
			"developer_info": miniApp.DeveloperInfo,
			"api_endpoint":   miniApp.APIEndpoint,
			"webhook_url":    miniApp.WebhookURL,
			"enabled":        miniApp.Enabled,
			"updated_at":     time.Now(),
		},
	}
}