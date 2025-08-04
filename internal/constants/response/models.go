package response


type Response struct {
	Ok      bool           `json:"ok"`
	Status  int            `json:"status,omitempty"`
	Message string         `json:"message,omitempty"`
	Data    any            `json:"data,omitempty"`
	Meta    any            `json:"meta,omitempty"`
	Error   *ErrorResponse `json:"error,omitempty"`
}

type ErrorResponse struct {
	StausCode  int          `json:"status_code"`
	Message    string       `json:"message"`
	FieldError []FieldError `json:"field_error,omitempty"`
}

type FieldError struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}
