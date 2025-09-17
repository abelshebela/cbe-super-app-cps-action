package common

// type Exception string

// const (
// 	Unauthorized        Exception = "UNAUTHORIZED"
// 	Forbidden           Exception = "FORBIDDEN"
// 	NotFound            Exception = "NOT_FOUND"
// 	BadRequest          Exception = "BAD_REQUEST"
// 	Conflict            Exception = "CONFLICT"
// 	InternalServerError Exception = "INTERNAL_SERVER_ERROR"
// 	ServiceUnavailable  Exception = "SERVICE_UNAVAILABLE"
// )

// type ErrorInfo struct {
// 	Code       Exception
// 	StatusCode int
// 	Message    string
// }

// var ErrorCode = map[Exception]ErrorInfo{
// 	Unauthorized: {
// 		Code:       Unauthorized,
// 		StatusCode: http.StatusUnauthorized,
// 		Message:    "No authentication method provided.",
// 	},
// 	Forbidden: {
// 		Code:       Forbidden,
// 		StatusCode: http.StatusForbidden,
// 		Message:    "You do not have permission to access this resource.",
// 	},
// 	NotFound: {
// 		Code:       NotFound,
// 		StatusCode: http.StatusNotFound,
// 		Message:    "The requested resource was not found.",
// 	},
// 	BadRequest: {
// 		Code:       BadRequest,
// 		StatusCode: http.StatusBadRequest,
// 		Message:    "The request was invalid or cannot be processed.",
// 	},
// 	Conflict: {
// 		Code:       Conflict,
// 		StatusCode: http.StatusConflict,
// 		Message:    "A conflict occurred with the current state of the resource.",
// 	},
// 	InternalServerError: {
// 		Code:       InternalServerError,
// 		StatusCode: http.StatusInternalServerError,
// 		Message:    "An unexpected error occurred on the server.",
// 	},
// 	ServiceUnavailable: {
// 		Code:       ServiceUnavailable,
// 		StatusCode: http.StatusServiceUnavailable,
// 		Message:    "The service is temporarily unavailable. Please try again later.",
// 	},
// }
