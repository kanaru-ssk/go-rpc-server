package response

type ErrorCode string

const (
	// 4xx errors
	ErrInvalidRequestBody ErrorCode = "INVALID_REQUEST_BODY"
	ErrMethodNotAllowed   ErrorCode = "METHOD_NOT_ALLOWED"
	ErrNotFound           ErrorCode = "NOT_FOUND"

	// 5xx errors
	ErrInternalServerError ErrorCode = "INTERNAL_SERVER_ERROR"
)

type ErrorJson struct {
	ErrorCode ErrorCode `json:"errorCode"`
}

type ErrResMap struct {
	Err        error
	StatusCode int
	ErrorCode  ErrorCode
}
