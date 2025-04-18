package errors

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	ErrBadRequest     = NewAppError(http.StatusBadRequest, "The request is invalid or malformed", "bad_request")
	ErrInvalidInput   = NewAppError(http.StatusBadRequest, "The provided input is invalid", "invalid_input")
	ErrPasswordTooWeak = NewAppError(http.StatusBadRequest, "Password does not meet the minimum requirements", "password_too_weak")
	ErrInvalidEmail   = NewAppError(http.StatusBadRequest, "The provided email address is invalid", "invalid_email")
	ErrInvalidProfileType = NewAppError(http.StatusBadRequest, "The provided profile type is invalid", "invalid_profile_type")
	ErrCannotFollowSelf = NewAppError(http.StatusBadRequest, "You cannot follow yourself", "cannot_follow_self")
	
	ErrUnauthorized   = NewAppError(http.StatusUnauthorized, "You are not authorized to perform this action", "unauthorized")
	ErrInvalidCredentials = NewAppError(http.StatusUnauthorized, "Invalid email or password", "invalid_credentials")
	ErrSessionExpired = NewAppError(http.StatusUnauthorized, "Your session has expired, please log in again", "session_expired")
	ErrInvalidToken   = NewAppError(http.StatusUnauthorized, "The provided token is invalid", "invalid_token")
	
	ErrForbidden      = NewAppError(http.StatusForbidden, "You don't have permission to access this resource", "forbidden")
	ErrNotGroupMember = NewAppError(http.StatusForbidden, "You are not a member of this group", "not_group_member")
	ErrNotGroupAdmin  = NewAppError(http.StatusForbidden, "This action requires group admin privileges", "not_group_admin")
	
	ErrNotFound       = NewAppError(http.StatusNotFound, "The requested resource was not found", "not_found")
	ErrUserNotFound   = NewAppError(http.StatusNotFound, "The requested user was not found", "user_not_found")
	ErrProfileNotFound = NewAppError(http.StatusNotFound, "The requested profile was not found", "profile_not_found")
	ErrPostNotFound   = NewAppError(http.StatusNotFound, "The requested post was not found", "post_not_found")
	ErrGroupNotFound  = NewAppError(http.StatusNotFound, "The requested group was not found", "group_not_found")
	
	ErrConflict       = NewAppError(http.StatusConflict, "The request conflicts with the current state", "conflict")
	ErrUserAlreadyExists = NewAppError(http.StatusConflict, "A user with this email already exists", "user_already_exists")
	ErrAlreadyFollowing = NewAppError(http.StatusConflict, "You are already following this user", "already_following")
	
	ErrInternalServer = NewAppError(http.StatusInternalServerError, "An internal server error occurred", "internal_error")
)

type AppError struct {
	StatusCode int         `json:"-"`
	Message    string      `json:"message"`
	Code       string      `json:"code"`
	Details    any         `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) WithDetails(details any) *AppError {
	return &AppError{
		StatusCode: e.StatusCode,
		Message:    e.Message,
		Code:       e.Code,
		Details:    details,
	}
}

func NewAppError(statusCode int, message, code string) *AppError {
	return &AppError{
		StatusCode: statusCode,
		Message:    message,
		Code:       code,
	}
}

func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	
	var appErr *AppError
	if errors.As(err, &appErr) {
		return &AppError{
			StatusCode: appErr.StatusCode,
			Message:    fmt.Sprintf("%s: %s", message, appErr.Message),
			Code:       appErr.Code,
			Details:    appErr.Details,
		}
	}
	
	return &AppError{
		StatusCode: http.StatusInternalServerError,
		Message:    fmt.Sprintf("%s: %v", message, err),
		Code:       "internal_error",
	}
}

func HTTPStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.StatusCode
	}
	
	return http.StatusInternalServerError
}

type FieldErrors map[string]string

func NewValidationError(message string, fields FieldErrors) *AppError {
	return &AppError{
		StatusCode: http.StatusBadRequest,
		Message:    message,
		Code:       "validation_error",
		Details:    fields,
	}
}

func IsNotFound(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.StatusCode == http.StatusNotFound
}

func IsUnauthorized(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.StatusCode == http.StatusUnauthorized
}

func IsForbidden(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.StatusCode == http.StatusForbidden
}

func IsConflict(err error) bool {
	var appErr *AppError
	return errors.As(err, &appErr) && appErr.StatusCode == http.StatusConflict
}