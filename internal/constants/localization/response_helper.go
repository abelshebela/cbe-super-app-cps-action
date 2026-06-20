package localization

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"cbe-super-app-cps-action/internal/constants"
	"cbe-super-app-cps-action/internal/constants/types"
)

type ContextResponseWriter struct {
	http.ResponseWriter
	ctx context.Context
}

func NewContextResponseWriter(w http.ResponseWriter, ctx context.Context) *ContextResponseWriter {
	return &ContextResponseWriter{ResponseWriter: w, ctx: ctx}
}

func (cw *ContextResponseWriter) SetContext(ctx context.Context) {
	cw.ctx = ctx
}

func GetMetadataFromWriter(w http.ResponseWriter) *types.ContextMetadata {
	ctx := contextFromWriter(w)
	return types.GetMetadata(ctx)
}

func (cw *ContextResponseWriter) Context() context.Context {
	return cw.ctx
}

func (cw *ContextResponseWriter) Unwrap() http.ResponseWriter {
	return cw.ResponseWriter
}

type contextSetter interface {
	SetContext(context.Context)
}

type responseWriterUnwrapper interface {
	Unwrap() http.ResponseWriter
}

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

func contextFromWriter(w http.ResponseWriter) context.Context {
	if cw, ok := w.(*ContextResponseWriter); ok {
		return cw.Context()
	}
	return context.Background()
}

func extractTraceID(ctx context.Context) string {
	if v, ok := ctx.Value(constants.ContextKey("trace_id")).(string); ok {
		return v
	}
	return ""
}

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
		w.Header().Set("x-action-code", actionCode)
	}
}

func ApplyActionCodeHeaderFromWriter(w http.ResponseWriter, ctx context.Context) http.ResponseWriter {
	actionCode := extractActionCode(ctx)
	if actionCode == "" {
		return w
	}
	w.Header().Set("x-action-code", actionCode)
	return w
}

type StandardResponse struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	TraceID   string      `json:"trace_id,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

type ErrorDetail struct {
	Code        string                 `json:"code"`
	Message     string                 `json:"message"`
	StatusCode  int                    `json:"status_code"`
	Type        string                 `json:"type"`
	FieldErrors []FieldError           `json:"field_errors,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
}

type FieldError struct {
	Field      string `json:"field"`
	Message    string `json:"message"`
	Value      string `json:"value,omitempty"`
	Constraint string `json:"constraint,omitempty"`
}

func SendSuccessResponse(w http.ResponseWriter, responseCode ResponseCode, data interface{}) {
	ctx := contextFromWriter(w)
	setActionCodeHeader(w, ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode.StatusCode)

	response := StandardResponse{
		Status:    responseCode.StatusCode,
		Message:   responseCode.Message,
		Data:      data,
		TraceID:   extractTraceID(ctx),
		RequestID: extractRequestID(ctx),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		SendErrorResponse(w, ErrorUnexpectedError, nil, nil)
	}
}

type PaginatedStandardResponse struct {
	Status    int         `json:"status"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Meta      interface{} `json:"meta,omitempty"`
	TraceID   string      `json:"trace_id,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

func SendPaginatedSuccessResponse(w http.ResponseWriter, responseCode ResponseCode, data, meta interface{}) {
	ctx := contextFromWriter(w)

	setActionCodeHeader(w, ctx)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(responseCode.StatusCode)

	response := PaginatedStandardResponse{
		Status:    responseCode.StatusCode,
		Message:   responseCode.Message,
		Data:      data,
		Meta:      meta,
		TraceID:   extractTraceID(ctx),
		RequestID: extractRequestID(ctx),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		SendErrorResponse(w, ErrorUnexpectedError, nil, nil)
	}
}

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
		w.WriteHeader(StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"ok":      false,
			"status":  StatusInternalServerError,
			"message": "Failed to encode response",
		})
	}
}

func SendValidationErrorResponse(w http.ResponseWriter, fieldErrors []FieldError) {
	SendErrorResponse(w, ErrorValidationFailed, fieldErrors, nil)
}

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

func CreateFieldError(field, message, value, constraint string) FieldError {
	return FieldError{
		Field:      field,
		Message:    message,
		Value:      value,
		Constraint: constraint,
	}
}

func CreateFieldErrors(errors ...FieldError) []FieldError {
	return errors
}

func GetResponseCodeByCode(code string) (ResponseCode, bool) {
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
		newResp := ResponseCode{
			Code:       strings.ToUpper(strings.Join(strings.Split(err, " "), "_")),
			TimeStamp:  time.Now(),
			StatusCode: statusCode,
			Message:    message,
			Type:       "error",
		}
		ResponseCodesList = append(ResponseCodesList, newResp)
		resp = newResp
	}
	return resp

}
