package feedback

import (
	"cbe-super-app-cps-action/internal/constants/types"
)

type FeedbackRequest struct {
	Responses map[string]types.Response `json:"responses" bson:"responses"`
}
