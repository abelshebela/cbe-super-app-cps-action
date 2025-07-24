package middleware

import (
	"cbe-super-app-budget/platform/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type ErrorType string

const (
	// Business logic errors
	ErrorTypeValidation   ErrorType = "validation"
	ErrorTypeNotFound     ErrorType = "not_found"
	ErrorTypeConflict     ErrorType = "conflict"
	ErrorTypeUnauthorized ErrorType = "unauthorized"
	ErrorTypeForbidden    ErrorType = "forbidden"
	ErrorTypeRateLimit    ErrorType = "rate_limit"

	// System errors
	ErrorTypeInternal    ErrorType = "internal"
	ErrorTypeTimeout     ErrorType = "timeout"
	ErrorTypeUnavailable ErrorType = "unavailable"
	ErrorTypeDatabase    ErrorType = "database"
	ErrorTypeExternal    ErrorType = "external_service"
	ErrorTypeNetwork     ErrorType = "network"

	// Request/Response errors
	ErrorTypeBadRequest      ErrorType = "bad_request"
	ErrorTypeInvalidFormat   ErrorType = "invalid_format"
	ErrorTypePayloadTooLarge ErrorType = "payload_too_large"

	// Unknown errors
	ErrorTypeUnknown ErrorType = "unknown"
)

// errorTypeStatusMap maps error types to their corresponding HTTP status codes
var errorTypeStatusMap = map[ErrorType]int{
	// Business logic errors
	ErrorTypeValidation:   http.StatusBadRequest,
	ErrorTypeNotFound:     http.StatusNotFound,
	ErrorTypeConflict:     http.StatusConflict,
	ErrorTypeUnauthorized: http.StatusUnauthorized,
	ErrorTypeForbidden:    http.StatusForbidden,
	ErrorTypeRateLimit:    http.StatusTooManyRequests,

	// System errors
	ErrorTypeInternal:    http.StatusInternalServerError,
	ErrorTypeTimeout:     http.StatusRequestTimeout,
	ErrorTypeUnavailable: http.StatusServiceUnavailable,
	ErrorTypeDatabase:    http.StatusInternalServerError,
	ErrorTypeExternal:    http.StatusServiceUnavailable,
	ErrorTypeNetwork:     http.StatusServiceUnavailable,

	// Request/Response errors
	ErrorTypeBadRequest:      http.StatusBadRequest,
	ErrorTypeInvalidFormat:   http.StatusBadRequest,
	ErrorTypePayloadTooLarge: http.StatusRequestEntityTooLarge,

	// Unknown errors
	ErrorTypeUnknown: http.StatusInternalServerError,
}

func getStatusCodeForErrorType(errorType ErrorType) int {
	if statusCode, exists := errorTypeStatusMap[errorType]; exists {
		return statusCode
	}

	return http.StatusInternalServerError
}

type AppError struct {
	Type        ErrorType              `json:"type"`
	Message     string                 `josn:"message"`
	StatusCode  int                    `json:"status_code"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Cause       error                  `json:"-"`
	RequestID   string                 `json:"request_id,omitempty"`
	Timestamp   time.Time              `json:"timestamp"`
	Service     string                 `json:"service"`
	Operation   string                 `json:"operation"`
	UserMessage string                 `json:"user_message,omitempty"`
}

type ValidationError struct {
	Field   string      `json:"field"`
	Value   interface{} `json:"value"`
	Message string      `json:"message"`
	Tag     string      `json:"tag"`
}

type ErrorResponse struct {
	Status    int         `json:"status"`
	Error     string      `json:"error"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details"`
	RequestID string      `json:"request_id,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Path      string      `json:"path,omitempty"`
	Operation string      `json:"operation"`
	Method    string      `json:"method,omitempty"`
}

func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Message, e.Cause)
	}

	return fmt.Sprintf("[%s] %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Cause
}

func (e *AppError) WithDetail(key string, value interface{}) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]interface{})
	}

	e.Details[key] = value

	return e
}

func (e *AppError) WithOperation(operation string) *AppError {
	e.Operation = operation
	return e
}

func (e *AppError) WithService(service string) *AppError {
	e.Service = service
	return e
}

func (e *AppError) WithUserMessage(message string) *AppError {
	e.UserMessage = message
	return e
}

func (e *AppError) WithCause(err error) *AppError {
	e.Cause = err
	return e
}

func NewAppError(errorType ErrorType, message string) *AppError {
	return &AppError{
		Type:       errorType,
		Message:    message,
		StatusCode: getStatusCodeForErrorType(errorType),
		Timestamp:  time.Now(),
	}
}

func NewAppErrorWithCause(errorType ErrorType, message string, cause error) *AppError {
	return &AppError{
		Type:       errorType,
		Message:    message,
		StatusCode: getStatusCodeForErrorType(errorType),
		Cause:      cause,
		Timestamp:  time.Now(),
	}
}

func NewAppErrorWithStatusCode(errorType ErrorType, message string, statusCode int) *AppError {
	return &AppError{
		Type:       errorType,
		Message:    message,
		StatusCode: statusCode,
		Timestamp:  time.Now(),
	}
}

func NewValidationError(message string, details map[string]interface{}) *AppError {
	return &AppError{
		Type:        ErrorTypeValidation,
		Message:     message,
		StatusCode:  getStatusCodeForErrorType(ErrorTypeValidation),
		Details:     details,
		Timestamp:   time.Now(),
		UserMessage: "Please check you input",
	}
}

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Type:        ErrorTypeNotFound,
		Message:     fmt.Sprintf("%s not found", resource),
		StatusCode:  getStatusCodeForErrorType(ErrorTypeNotFound),
		Timestamp:   time.Now(),
		UserMessage: fmt.Sprintf("The requested %s was not found", resource),
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Type:        ErrorTypeUnauthorized,
		Message:     message,
		StatusCode:  getStatusCodeForErrorType(ErrorTypeUnauthorized),
		Timestamp:   time.Now(),
		UserMessage: "Authentication required!",
	}
}

func NewForbidenError(message string) *AppError {
	return &AppError{
		Type:        ErrorTypeForbidden,
		Message:     message,
		StatusCode:  getStatusCodeForErrorType(ErrorTypeForbidden),
		Timestamp:   time.Now(),
		UserMessage: "Access denied!",
	}
}

func NewInternalError(message string, cause error) *AppError {
	return &AppError{
		Type:        ErrorTypeInternal,
		Message:     message,
		StatusCode:  getStatusCodeForErrorType(ErrorTypeInternal),
		Cause:       cause,
		Timestamp:   time.Now(),
		UserMessage: "An unexpected error occured. Please try again later",
	}
}

func NewBadRequestError(message string, cause error) *AppError {
	return &AppError{
		Type:        ErrorTypeBadRequest,
		Message:     message,
		StatusCode:  getStatusCodeForErrorType(ErrorTypeBadRequest),
		Cause:       cause,
		Timestamp:   time.Now(),
		UserMessage: "An unexpected error occured. Please try again later",
	}
}

func NewDatabaseError(message string, cause error) *AppError {
	return &AppError{
		Type:        ErrorTypeDatabase,
		Message:     message,
		StatusCode:  getStatusCodeForErrorType(ErrorTypeDatabase),
		Cause:       cause,
		Timestamp:   time.Now(),
		UserMessage: "Database operation failed. Please try again",
	}
}

func NewExternalServiceError(service string, cause error) *AppError {
	return &AppError{
		Type:        ErrorTypeExternal,
		Message:     fmt.Sprintf("External service %s failed", service),
		StatusCode:  getStatusCodeForErrorType(ErrorTypeExternal),
		Cause:       cause,
		Timestamp:   time.Now(),
		UserMessage: "A required service is temporarily unavailable. Please try again later",
	}
}

func NewTimeoutError(operation string) *AppError {
	return &AppError{
		Type:        ErrorTypeTimeout,
		Message:     fmt.Sprintf("Operation %v timeout", operation),
		StatusCode:  getStatusCodeForErrorType(ErrorTypeBadRequest),
		Timestamp:   time.Now(),
		UserMessage: "The request took too long to process. Please try again",
	}
}

func HandlePanic(logger logger.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					handlePanic(w, r, err, logger)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	requestID := getRequestID(r)
	appErr := convertToAppErr(err, requestID)

	writeErrorResponse(w, r, appErr)
}

func convertToAppErr(err error, requestID string) *AppError {
	if appErr, ok := err.(*AppError); ok {
		if appErr.RequestID == "" {
			appErr.RequestID = requestID
		}

		appErr.Timestamp = time.Now()
		if appErr.Cause != nil {
			if validationErrs, ok := appErr.Cause.(validator.ValidationErrors); ok {
				details := make(map[string]interface{})
				for _, fieldErr := range validationErrs {
					details[fieldErr.Field()] = fieldErr.Error()
					// details[fieldErr.Field()] =fmt.Sprintf("failed on '%s' tag", fieldErr.Tag())
				}

				appErr.Details = details
			}
		}

		return appErr
	}

	if validationErrs, ok := err.(validator.ValidationErrors); ok {
		details := make(map[string]interface{})
		for _, fieldErr := range validationErrs {
			details[fieldErr.Field()] = fieldErr.Error()
			// details[fieldErr.Field()] =fmt.Sprintf("failed on '%s' tag", fieldErr.Tag())
		}

		return NewValidationError("Validation failed", details).
			WithDetail("request_id", requestID)
	}

	if jsonErr, ok := err.(*json.SyntaxError); ok {
		return NewAppError(ErrorTypeInvalidFormat, "Invalid JSON format").
			WithCause(jsonErr).
			WithDetail("request_id", requestID).
			WithUserMessage("Invalid request format")
	}

	return NewInternalError("Internal server error", err).
		WithDetail("request_id", requestID)
}

func handlePanic(w http.ResponseWriter, r *http.Request, panicErr interface{}, logger logger.Logger) {
	requestID := getRequestID(r)

	// Create panic error
	appErr := NewInternalError("Panic occurred during request processing",
		fmt.Errorf("panic: %v", panicErr)).
		WithDetail("request_id", requestID).
		WithDetail("stack_trace", string(debug.Stack()))

		// Log the panic with stack trace
	logger.Error(r.Context(), "[Panic recovered]",
		zap.String("request_id", requestID),
		zap.String("method", r.Method),
		zap.String("path", r.URL.Path),
		zap.Any("panic", panicErr),
		zap.String("stack_trace", string(debug.Stack())),
	)

	writeErrorResponse(w, r, appErr)
}

func writeErrorResponse(w http.ResponseWriter, r *http.Request, appErr *AppError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.StatusCode)

	response := ErrorResponse{
		Status:    appErr.StatusCode,
		Error:     string(appErr.Type),
		Message:   appErr.UserMessage,
		Details:   appErr.Details,
		RequestID: appErr.RequestID,
		Timestamp: appErr.Timestamp,
		Path:      r.URL.Path,
		Method:    r.Method,
		Operation: appErr.Operation,
	}

	if response.Message == "" {
		response.Message = appErr.Message
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)

	if err := encoder.Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func getRequestID(r *http.Request) string {
	// Try to get from context first
	if requestID := r.Context().Value(middleware.RequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			return id
		}
	}

	// Generate a new one if not found
	return fmt.Sprintf("req-%d", time.Now().UnixNano())
}
