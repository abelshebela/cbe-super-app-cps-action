package budget

import (
	"cbe-super-app-budget/internal/constants/dto"
	"cbe-super-app-budget/internal/handler/middleware"
	"cbe-super-app-budget/internal/handler/rest"
	"cbe-super-app-budget/internal/service"
	"cbe-super-app-budget/platform/logger"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type budgetHandler struct {
	service service.BudgetService
	logger  logger.Logger
}

func (s *budgetHandler) DeleteOne(w http.ResponseWriter, r *http.Request) {
	spedingID := chi.URLParam(r, "budget_id")
	if spedingID == "" {
		appErr := middleware.NewValidationError("Budget ID is required", map[string]interface{}{
			"filed": "budget_id",
		}).WithService("cbe-super-app-budget").WithOperation("delete_one")

		middleware.HandleError(w, r, appErr)
		return
	}

	if err := s.service.Remove(r.Context(), spedingID); err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	middleware.WriteJSONResponse(w,
		http.StatusOK,
		"budget deleted succesfully",
		nil,
	)
}

func (s *budgetHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	budgets, err := s.service.Get(r.Context())
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	middleware.WriteJSONResponse(w,
		http.StatusOK,
		"budget retrived succesfully",
		budgets,
	)
}

func (s *budgetHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	budgetID := chi.URLParam(r, "budget_id")
	if budgetID == "" {
		appErr := middleware.NewValidationError("BudgetID is requried", map[string]interface{}{
			"field": "budget_id",
		}).WithService("cbe-super-app-budget").WithOperation("get_one")
		middleware.HandleError(w, r, appErr)
		return
	}

	budget, err := s.service.GetOne(r.Context(), budgetID)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	middleware.WriteJSONResponse(w,
		http.StatusOK,
		"Budget retirve succesfully",
		budget,
	)

}

func (s *budgetHandler) InsertOne(w http.ResponseWriter, r *http.Request) {
	var req dto.BudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appErr := middleware.NewAppError(middleware.ErrorTypeValidation, "failed to decode request body").
			WithCause(err).
			WithService("cbe-super-app-budget").
			WithOperation("insert_one")
		middleware.HandleError(w, r, appErr)
		return
	}

	if err := req.Validate(); err != nil {
		appErr := middleware.NewAppError(middleware.ErrorTypeValidation, "Validation failed").
			WithCause(err).
			WithService("cbe-super-app-budget").
			WithOperation("insert_one")

		middleware.HandleError(w, r, appErr)
		return
	}

	budget, err := s.service.Add(r.Context(), req)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	middleware.WriteJSONResponse(w,
		http.StatusCreated,
		"budget created successfully",
		budget,
	)

}

func (s *budgetHandler) UpdateOne(w http.ResponseWriter, r *http.Request) {
	var req dto.BudgetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		appErr := middleware.NewAppError(middleware.ErrorTypeValidation, "failed to decode request body").
			WithCause(err).
			WithService("cbe-super-app-budget").
			WithOperation("insert_one")
		middleware.HandleError(w, r, appErr)
		return
	}

	if err := req.Validate(); err != nil {
		appErr := middleware.NewAppError(middleware.ErrorTypeValidation, "Validation failed").
			WithCause(err).
			WithService("cbe-super-app-budget").
			WithOperation("insert_one")

		middleware.HandleError(w, r, appErr)
		return
	}

	if err := s.service.Modify(r.Context(), req); err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	middleware.WriteJSONResponse(w,
		http.StatusAccepted,
		"budget updated successfully",
		nil,
	)
}

func NewBudgetService(service service.BudgetService, logger logger.Logger) rest.BudgetHandler {
	return &budgetHandler{
		service: service,
		logger:  logger,
	}
}
