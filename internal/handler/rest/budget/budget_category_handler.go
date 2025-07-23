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

type budgetCategoryHandler struct {
	service service.BudgetCategoryService
	logger  logger.Logger
}

func (s *budgetCategoryHandler) DeleteOne(w http.ResponseWriter, r *http.Request) {
	spedingID := chi.URLParam(r, "budget_category_id")
	if spedingID == "" {
		appErr := middleware.NewValidationError("BudgetCategory ID is required", map[string]interface{}{
			"filed": "budget_category_id",
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
		"budget_category deleted succesfully",
		nil,
	)
}

func (s *budgetCategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	budget_categorys, err := s.service.Get(r.Context())
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	middleware.WriteJSONResponse(w,
		http.StatusOK,
		"budget_category retrived succesfully",
		budget_categorys,
	)
}

func (s *budgetCategoryHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	budget_categoryID := chi.URLParam(r, "budget_category_id")
	if budget_categoryID == "" {
		appErr := middleware.NewValidationError("BudgetCategoryID is requried", map[string]interface{}{
			"field": "budget_category_id",
		}).WithService("cbe-super-app-budget").WithOperation("get_one")
		middleware.HandleError(w, r, appErr)
		return
	}

	budget_category, err := s.service.GetOne(r.Context(), budget_categoryID)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	middleware.WriteJSONResponse(w,
		http.StatusOK,
		"BudgetCategory retirve succesfully",
		budget_category,
	)

}

func (s *budgetCategoryHandler) InsertOne(w http.ResponseWriter, r *http.Request) {
	var req dto.BudgetCategoryRequest
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
			WithService("category").
			WithOperation("insert_one")

		middleware.HandleError(w, r, appErr)
		return
	}

	budget_category, err := s.service.Add(r.Context(), req)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	middleware.WriteJSONResponse(w,
		http.StatusCreated,
		"budget_category created successfully",
		budget_category,
	)

}

func (s *budgetCategoryHandler) UpdateOne(w http.ResponseWriter, r *http.Request) {
	var req dto.BudgetCategoryRequest
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
			WithService("category").
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
		"budget_category updated successfully",
		nil,
	)
}

func NewBudgetCategoryService(service service.BudgetCategoryService, logger logger.Logger) rest.BudgetCategoryHandler {
	return &budgetCategoryHandler{
		service: service,
		logger:  logger,
	}
}
