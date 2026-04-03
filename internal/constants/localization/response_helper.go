package localization

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"
)

// ──────────────────────────────────────────────────────────────────────────────
// Context-aware ResponseWriter
// ──────────────────────────────────────────────────────────────────────────────

// ContextResponseWriter wraps http.ResponseWriter and carries the request context
// so that response helper functions can access trace IDs, user info, etc.
type ContextResponseWriter struct {
	http.ResponseWriter
	ctx context.Context
}

// NewContextResponseWriter creates a new ContextResponseWriter.
func NewContextResponseWriter(w http.ResponseWriter, ctx context.Context) *ContextResponseWriter {
	return &ContextResponseWriter{ResponseWriter: w, ctx: ctx}
}

// SetContext updates the bound request context.
func (cw *ContextResponseWriter) SetContext(ctx context.Context) {
	cw.ctx = ctx
}

// GetMetadataFromWriter extracts the ContextMetadata from the response writer's bound context.
func GetMetadataFromWriter(w http.ResponseWriter) *types.ContextMetadata {
	ctx := contextFromWriter(w)
	return types.GetMetadata(ctx)
}

// Context returns the bound request context.
func (cw *ContextResponseWriter) Context() context.Context {
	return cw.ctx
}

// Unwrap returns the underlying ResponseWriter (supports http.ResponseController).
func (cw *ContextResponseWriter) Unwrap() http.ResponseWriter {
	return cw.ResponseWriter
}

type contextSetter interface {
	SetContext(context.Context)
}

type responseWriterUnwrapper interface {
	Unwrap() http.ResponseWriter
}

// UpdateWriterContext aligns a wrapped response writer with the latest request context.
func UpdateWriterContext(w http.ResponseWriter, ctx context.Context) {
	for w != nil {
		if cw, ok := w.(contextSetter); ok {
			cw.SetContext(ctx)
			return
		}

		uw, ok := w.(responseWriterUnwrapper)
		if !ok {
			return
		}

		next := uw.Unwrap()
		if next == w {
			return
		}
		w = next
	}
}

// contextFromWriter extracts the context from the writer if it is a
// ContextResponseWriter; otherwise returns context.Background().
func contextFromWriter(w http.ResponseWriter) context.Context {
	if cw, ok := w.(*ContextResponseWriter); ok {
		return cw.Context()
	}
	return context.Background()
}

// extractTraceID returns the trace_id stored in the context (if any).
func extractTraceID(ctx context.Context) string {
	if v, ok := ctx.Value(constants.ContextKey("trace_id")).(string); ok {
		return v
	}
	return ""
}

// extractRequestID returns the x-request-id stored in the context (if any).
func extractRequestID(ctx context.Context) string {
	if v, ok := ctx.Value(constants.ContextKey("x-request-id")).(string); ok {
		return v
	}
	return ""
}

func extractActionCode(ctx context.Context) string {
	if md := types.GetMetadata(ctx); md != nil && md.CPSActionCode != "" {
		return md.CPSActionCode
	}
	if v, ok := ctx.Value(constants.ContextKey("cps_action_code")).(string); ok {
		return v
	}
	return ""
}

func setActionCodeHeader(w http.ResponseWriter, ctx context.Context) {
	if actionCode := extractActionCode(ctx); actionCode != "" {
		fmt.Printf("Setting action code header: %s \n", actionCode)
		w.Header().Set("x-action-code", actionCode)
	}
}

// ApplyActionCodeHeaderFromWriter sets the action code response header using the
// context currently bound to the response writer. Intended for handler-side use
// just before sending success responses.
func ApplyActionCodeHeaderFromWriter(w http.ResponseWriter, ctx context.Context) http.ResponseWriter {
	actionCode := extractActionCode(ctx)
	if actionCode == "" {
		return w
	}
	w.Header().Set("x-action-code", actionCode)
	return w
}

// StandardResponse represents the standardized API response structure
type StandardResponse struct {
	// Ok        bool        `json:"ok"`
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	TraceID   string      `json:"trace_id,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
	// Error     *ErrorDetail `json:"error,omitempty"`
}

// ErrorDetail represents error details in the response
type ErrorDetail struct {
	Code        string                 `json:"code"`
	Message     string                 `json:"message"`
	StatusCode  int                    `json:"status_code"`
	Type        string                 `json:"type"`
	FieldErrors []FieldError           `json:"field_errors,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

// FieldError represents validation field errors
type FieldError struct {
	Field      string `json:"field"`
	Message    string `json:"message"`
	Value      string `json:"value,omitempty"`
	Constraint string `json:"constraint,omitempty"`
}

// SendSuccessResponse sends a standardized success response
func SendSuccessResponse(w http.ResponseWriter, responseCode ResponseCode, data interface{}) {
	ctx := contextFromWriter(w)
	actionCode := extractActionCode(ctx)

	fmt.Printf("Sending success response with action code: %s \n", actionCode)

	setActionCodeHeader(w, ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode.StatusCode)

	response := StandardResponse{
		// Ok:        true,
		Status:    responseCode.StatusCode,
		Message:   responseCode.Message,
		Data:      data,
		TraceID:   extractTraceID(ctx),
		RequestID: extractRequestID(ctx),
	}
	// if responseCode.Type == "error" {
	// 	response.Ok = false
	// }

	if err := json.NewEncoder(w).Encode(response); err != nil {
		SendErrorResponse(w, ErrorUnexpectedError, nil, nil)
	}
}

// SendErrorResponse sends a standardized error response
func SendErrorResponse(w http.ResponseWriter, responseCode ResponseCode, fieldErrors []FieldError, details map[string]interface{}) {
	ctx := contextFromWriter(w)
	setActionCodeHeader(w, ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode.StatusCode)

	response := StandardResponse{
		Status:    responseCode.StatusCode,
		Message:   responseCode.Message,
		TraceID:   extractTraceID(ctx),
		RequestID: extractRequestID(ctx),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		// Fallback to basic error response if encoding fails
		w.WriteHeader(StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":      false,
			"status":  StatusInternalServerError,
			"message": "Failed to encode response",
		})
	}
}

// SendValidationErrorResponse sends a validation error response
func SendValidationErrorResponse(w http.ResponseWriter, fieldErrors []FieldError) {
	SendErrorResponse(w, ErrorValidationFailed, fieldErrors, nil)
}

// SendErrorByCodeResponse sends a validation error response
func SendErrorByCodeResponse(w http.ResponseWriter, code string) {
	respCode, ok := GetResponseCodeByCode(code)
	if !ok {
		SendErrorResponse(w, ErrorFormatter(code), nil, nil)
		return
	}

	if respCode.StatusCode >= http.StatusBadRequest || respCode.Type == "error" {
		SendErrorResponse(w, respCode, nil, nil)
	} else {
		SendSuccessResponse(w, respCode, nil)
	}
}

func ErrorFormatter(code string) ResponseCode {
	return ResponseCode{
		Message:    code,
		TimeStamp:  time.Now(),
		StatusCode: http.StatusBadRequest,
		Code:       strings.ToUpper(strings.ReplaceAll(code, " ", "_")),
		Type:       "error",
	}
}

// SendUnauthorizedResponse sends an unauthorized error response
func SendUnauthorizedResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = MsgUserUnauthorized
	}

	customResponseCode := ResponseCode{
		Code:       "ERROR_UNAUTHORIZED",
		StatusCode: StatusUnauthorized,
		TimeStamp:  time.Now(),
		Message:    message,
		Type:       "error",
	}

	SendErrorResponse(w, customResponseCode, nil, nil)
}

// SendForbiddenResponse sends a forbidden error response
func SendForbiddenResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = MsgUserForbidden
	}

	customResponseCode := ResponseCode{
		Code:       "ERROR_FORBIDDEN",
		TimeStamp:  time.Now(),
		StatusCode: StatusForbidden,
		Message:    message,
		Type:       "error",
	}

	SendErrorResponse(w, customResponseCode, nil, nil)
}

// SendNotFoundResponse sends a not found error response
func SendNotFoundResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = MsgUserNotFound
	}

	customResponseCode := ResponseCode{
		Code:       "ERROR_NOT_FOUND",
		TimeStamp:  time.Now(),
		StatusCode: StatusNotFound,
		Message:    message,
		Type:       "error",
	}

	SendErrorResponse(w, customResponseCode, nil, nil)
}

// SendBadRequestResponse sends a bad request error response
func SendBadRequestResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = MsgBadRequest
	}

	m := strings.Split(message, ":")
	msg := m[0]
	if len(m) > 1 {
		msg = m[1]
	}
	customResponseCode := ResponseCode{
		Code:       "ERROR_BAD_REQUEST",
		TimeStamp:  time.Now(),
		StatusCode: StatusBadRequest,
		Message:    msg,
		Type:       "error",
	}

	SendErrorResponse(w, customResponseCode, nil, nil)
}

// SendInternalServerErrorResponse sends an internal server error response
func SendInternalServerErrorResponse(w http.ResponseWriter, message string) {
	if message == "" {
		message = MsgInternalServerError
	}

	customResponseCode := ResponseCode{
		Code:       "ERROR_INTERNAL_SERVER_ERROR",
		TimeStamp:  time.Now(),
		StatusCode: StatusInternalServerError,
		Message:    message,
		Type:       "error",
	}

	SendErrorResponse(w, customResponseCode, nil, nil)
}

// CreateFieldError creates a field error for validation
func CreateFieldError(field, message, value, constraint string) FieldError {
	return FieldError{
		Field:      field,
		Message:    message,
		Value:      value,
		Constraint: constraint,
	}
}

// CreateFieldErrors creates multiple field errors
func CreateFieldErrors(errors ...FieldError) []FieldError {
	return errors
}

// GetResponseCodeByCode fetches a ResponseCode by its Code field from a predefined set of response codes.
func GetResponseCodeByCode(code string) (ResponseCode, bool) {
	// List all response codes to search through.

	for _, rc := range ResponseCodesList {
		if rc.Code == code {
			return rc, true
		}
	}
	return ResponseCode{}, false
}

func ErrorToResponseCode(err string, statusCode int, message string) ResponseCode {
	resp, ok := GetResponseCodeByCode(err)
	if !ok {
		new_resp := ResponseCode{
			Code:       strings.ToUpper(strings.Join(strings.Split(err, " "), "_")),
			TimeStamp:  time.Now(),
			StatusCode: statusCode,
			Message:    message,
			Type:       "error",
		}
		ResponseCodesList = append(ResponseCodesList, new_resp)
		resp = new_resp
	}
	return resp

}
