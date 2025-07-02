package handler

const (
	// General
	InternalError   = "INTERNAL_ERROR"
	InvalidJSON     = "INVALID_JSON"
	ValidationError = "VALIDATION_ERROR"

	// Auth
	Unauthorized     = "UNAUTHORIZED"
	TokenExpired     = "TOKEN_EXPIRED"
	UserNotFound     = "USER_NOT_FOUND"
	PermissionDenied = "PERMISSION_DENIED"

	// Todo-related
	TodoNotFound  = "TODO_NOT_FOUND"
	TitleTooShort = "TITLE_TOO_SHORT"
)
