package survey_sampling

import imodel "github.com/abelshebela/cbe-super-app-cps-action/internal/constants/model"

// Used when creating a survey sampling configuration.
type CreateSurveySamplingRequest struct {
	Name      string                `json:"name"`
	Method    string                `json:"method"`
	Config    imodel.SamplingConfig `json:"config"`
	SurveyURL string                `json:"survey_url"`
}

// Used when updating a survey sampling configuration.
type UpdateSurveySamplingRequest struct {
	Name      string                `json:"name"`
	Config    imodel.SamplingConfig `json:"config"`
	SurveyURL string                `json:"survey_url"`
}
