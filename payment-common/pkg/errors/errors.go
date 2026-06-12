package errors

type Code string

const (
	CodeOK                 Code = "OK"
	CodeInvalidArgument    Code = "INVALID_ARGUMENT"
	CodeUnauthorized       Code = "UNAUTHORIZED"
	CodeForbidden          Code = "FORBIDDEN"
	CodeDuplicateRequest   Code = "DUPLICATE_REQUEST"
	CodeInvalidTransition  Code = "INVALID_TRANSITION"
	CodeChannelUnavailable Code = "CHANNEL_UNAVAILABLE"
	CodeInternal           Code = "INTERNAL"
)

type AppError struct {
	Code    Code
	Message string
}

func (e AppError) Error() string {
	return string(e.Code) + ": " + e.Message
}

