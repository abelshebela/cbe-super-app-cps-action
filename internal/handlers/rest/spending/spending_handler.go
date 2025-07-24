package spending

import (
	"cbe-super-app-budget/internal/constants/dto"
	"cbe-super-app-budget/internal/handlers/middleware"
	"cbe-super-app-budget/internal/handlers/rest"
	"cbe-super-app-budget/internal/service"
	"cbe-super-app-budget/platform/logger"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type spendingHandler struct {
	service service.SpendingService
	logger  logger.Logger
}

func (s *spendingHandler) DeleteOne(w http.ResponseWriter, r *http.Request) {
	spedingID := chi.URLParam(r, "spending_id")
	if spedingID == "" {
		appErr := middleware.NewValidationError("Spending ID is required", map[string]interface{}{
			"filed": "spending_id",
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
		"spending deleted succesfully",
		nil,
	)
}

func (s *spendingHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	spendings, err := s.service.Get(r.Context())
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	middleware.WriteJSONResponse(w,
		http.StatusOK,
		"spending retrived succesfully",
		spendings,
	)
}

func (s *spendingHandler) GetOne(w http.ResponseWriter, r *http.Request) {
	spendingID := chi.URLParam(r, "spending_id")
	if spendingID == "" {
		appErr := middleware.NewValidationError("SpendingID is requried", map[string]interface{}{
			"field": "spending_id",
		}).WithService("cbe-super-app-budget").WithOperation("get_one")
		middleware.HandleError(w, r, appErr)
		return
	}

	spending, err := s.service.GetOne(r.Context(), spendingID)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	middleware.WriteJSONResponse(w,
		http.StatusOK,
		"Spending retirve succesfully",
		spending,
	)

}

func (s *spendingHandler) InsertOne(w http.ResponseWriter, r *http.Request) {
	var req dto.SpendingRequest
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

	spending, err := s.service.Add(r.Context(), req)
	if err != nil {
		middleware.HandleError(w, r, err)
		return
	}

	middleware.WriteJSONResponse(w,
		http.StatusCreated,
		"spending created successfully",
		spending,
	)

}

func (s *spendingHandler) UpdateOne(w http.ResponseWriter, r *http.Request) {
	var req dto.SpendingRequest
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
		"spending updated successfully",
		nil,
	)
}

func NewSpendingService(service service.SpendingService, logger logger.Logger) rest.SpendingHandler {
	return &spendingHandler{
		service: service,
		logger:  logger,
	}
}
