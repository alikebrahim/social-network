package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrorType represents a category of errors
type ErrorType string

// Error categories
const (
	// Generic error types
	ErrorTypeUnauthorized   ErrorType = "unauthorized"
	ErrorTypeForbidden      ErrorType = "forbidden"
	ErrorTypeNotFound       ErrorType = "not_found"
	ErrorTypeBadRequest     ErrorType = "bad_request"
	ErrorTypeConflict       ErrorType = "conflict"
	ErrorTypeInvalidInput   ErrorType = "invalid_input"
	ErrorTypeInternalServer ErrorType = "internal_server"
	ErrorTypeTimeout        ErrorType = "timeout"
	ErrorTypeUnavailable    ErrorType = "unavailable"

	// Module-specific error types
	ErrorTypeAuth        ErrorType = "auth"
	ErrorTypeProfile     ErrorType = "profile"
	ErrorTypePosts       ErrorType = "posts"
	ErrorTypeFollowing   ErrorType = "following"
	ErrorTypeGroups      ErrorType = "groups"
	ErrorTypeChat        ErrorType = "chat"
	ErrorTypeNotification ErrorType = "notification"
)

// AppError is a structured error that provides consistent error handling
type AppError struct {
	Type       ErrorType  `json:"-"` // Not shown to clients
	HTTPStatus int        `json:"-"` // Not shown to clients
	Message    string     `json:"message"`
	Code       string     `json:"code"`
	Details    any        `json:"details,omitempty"`
	Err        error      `json:"-"` // Not shown to clients
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the wrapped error
func (e *AppError) Unwrap() error {
	return e.Err
}

// WithDetails adds optional details to the error and returns it
func (e *AppError) WithDetails(details any) *AppError {
	e.Details = details
	return e
}

// HasField checks if the error has a specific field
func (e *AppError) HasField(field string) bool {
	if fields, ok := e.Details.(map[string]string); ok {
		_, exists := fields[field]
		return exists
	}
	return false
}

// AddField adds a field-specific error message
func (e *AppError) AddField(field, message string) *AppError {
	if e.Details == nil {
		e.Details = make(map[string]string)
	}
	
	if fields, ok := e.Details.(map[string]string); ok {
		fields[field] = message
		e.Details = fields
	}
	
	return e
}

// NewAppError creates a new application error
func NewAppError(errType ErrorType, httpStatus int, message, code string, err error) *AppError {
	return &AppError{
		Type:       errType,
		HTTPStatus: httpStatus,
		Message:    message,
		Code:       code,
		Err:        err,
	}
}

// Wrap an existing error with additional context
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	
	// If it's already an AppError, just update the message
	if appErr, ok := err.(*AppError); ok {
		return &AppError{
			Type:       appErr.Type,
			HTTPStatus: appErr.HTTPStatus,
			Message:    message + ": " + appErr.Message,
			Code:       appErr.Code,
			Details:    appErr.Details,
			Err:        appErr.Err,
		}
	}
	
	// Otherwise wrap the standard error
	return &AppError{
		Type:       ErrorTypeInternalServer,
		HTTPStatus: http.StatusInternalServerError,
		Message:    message,
		Code:       "internal_error",
		Err:        err,
	}
}

// HTTPStatusFromError extracts the HTTP status code from an error
func HTTPStatusFromError(err error) int {
	if err == nil {
		return http.StatusOK
	}
	
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.HTTPStatus
	}
	
	return http.StatusInternalServerError
}

// FieldErrors represents validation errors for multiple fields
type FieldErrors map[string]string

// Common errors for reuse throughout the application
var (
	// Generic errors
	ErrUnauthorized = NewAppError(
		ErrorTypeUnauthorized,
		http.StatusUnauthorized,
		"You are not authorized to perform this action",
		"unauthorized",
		nil,
	)
	
	ErrForbidden = NewAppError(
		ErrorTypeForbidden,
		http.StatusForbidden,
		"You don't have permission to access this resource",
		"forbidden",
		nil,
	)
	
	ErrNotFound = NewAppError(
		ErrorTypeNotFound,
		http.StatusNotFound,
		"The requested resource was not found",
		"not_found",
		nil,
	)
	
	ErrBadRequest = NewAppError(
		ErrorTypeBadRequest,
		http.StatusBadRequest,
		"The request is invalid or malformed",
		"bad_request",
		nil,
	)
	
	ErrConflict = NewAppError(
		ErrorTypeConflict,
		http.StatusConflict,
		"The request conflicts with the current state",
		"conflict",
		nil,
	)
	
	ErrInvalidInput = NewAppError(
		ErrorTypeInvalidInput,
		http.StatusBadRequest,
		"The provided input is invalid",
		"invalid_input",
		nil,
	)
	
	ErrInternalServer = NewAppError(
		ErrorTypeInternalServer,
		http.StatusInternalServerError,
		"An internal server error occurred",
		"internal_error",
		nil,
	)
	
	// Auth errors
	ErrInvalidCredentials = NewAppError(
		ErrorTypeAuth,
		http.StatusUnauthorized,
		"Invalid email or password",
		"invalid_credentials",
		nil,
	)
	
	ErrUserAlreadyExists = NewAppError(
		ErrorTypeAuth,
		http.StatusConflict,
		"A user with this email already exists",
		"user_already_exists",
		nil,
	)
	
	ErrSessionExpired = NewAppError(
		ErrorTypeAuth,
		http.StatusUnauthorized,
		"Your session has expired, please log in again",
		"session_expired",
		nil,
	)
	
	ErrInvalidToken = NewAppError(
		ErrorTypeAuth,
		http.StatusUnauthorized,
		"The provided token is invalid",
		"invalid_token",
		nil,
	)
	
	ErrPasswordTooWeak = NewAppError(
		ErrorTypeAuth,
		http.StatusBadRequest,
		"Password does not meet the minimum requirements",
		"password_too_weak",
		nil,
	)
	
	ErrInvalidEmail = NewAppError(
		ErrorTypeAuth,
		http.StatusBadRequest,
		"The provided email address is invalid",
		"invalid_email",
		nil,
	)
	
	// Profile errors
	ErrProfileNotFound = NewAppError(
		ErrorTypeProfile,
		http.StatusNotFound,
		"The requested profile was not found",
		"profile_not_found",
		nil,
	)
	
	ErrInvalidProfileType = NewAppError(
		ErrorTypeProfile,
		http.StatusBadRequest,
		"The provided profile type is invalid",
		"invalid_profile_type",
		nil,
	)
	
	ErrProfileUpdateFailed = NewAppError(
		ErrorTypeProfile,
		http.StatusInternalServerError,
		"Failed to update profile",
		"profile_update_failed",
		nil,
	)
	
	// Posts errors
	ErrPostNotFound = NewAppError(
		ErrorTypePosts,
		http.StatusNotFound,
		"The requested post was not found",
		"post_not_found",
		nil,
	)
	
	ErrInvalidPrivacyLevel = NewAppError(
		ErrorTypePosts,
		http.StatusBadRequest,
		"The provided privacy level is invalid",
		"invalid_privacy_level",
		nil,
	)
	
	ErrUnauthorizedToViewPost = NewAppError(
		ErrorTypePosts,
		http.StatusForbidden,
		"You are not authorized to view this post",
		"unauthorized_to_view_post",
		nil,
	)
	
	ErrCommentNotFound = NewAppError(
		ErrorTypePosts,
		http.StatusNotFound,
		"The requested comment was not found",
		"comment_not_found",
		nil,
	)
	
	// Following errors
	ErrAlreadyFollowing = NewAppError(
		ErrorTypeFollowing,
		http.StatusConflict,
		"You are already following this user",
		"already_following",
		nil,
	)
	
	ErrNotFollowing = NewAppError(
		ErrorTypeFollowing,
		http.StatusBadRequest,
		"You are not following this user",
		"not_following",
		nil,
	)
	
	ErrCannotFollowSelf = NewAppError(
		ErrorTypeFollowing,
		http.StatusBadRequest,
		"You cannot follow yourself",
		"cannot_follow_self",
		nil,
	)
	
	ErrFollowRequestPending = NewAppError(
		ErrorTypeFollowing,
		http.StatusConflict,
		"A follow request is already pending",
		"follow_request_pending",
		nil,
	)
	
	// Groups errors
	ErrGroupNotFound = NewAppError(
		ErrorTypeGroups,
		http.StatusNotFound,
		"The requested group was not found",
		"group_not_found",
		nil,
	)
	
	ErrNotGroupMember = NewAppError(
		ErrorTypeGroups,
		http.StatusForbidden,
		"You are not a member of this group",
		"not_group_member",
		nil,
	)
	
	ErrNotGroupAdmin = NewAppError(
		ErrorTypeGroups,
		http.StatusForbidden,
		"This action requires group admin privileges",
		"not_group_admin",
		nil,
	)
	
	ErrEventNotFound = NewAppError(
		ErrorTypeGroups,
		http.StatusNotFound,
		"The requested event was not found",
		"event_not_found",
		nil,
	)
	
	// Chat errors
	ErrMessageNotDelivered = NewAppError(
		ErrorTypeChat,
		http.StatusInternalServerError,
		"The message could not be delivered",
		"message_not_delivered",
		nil,
	)
	
	// Notification errors
	ErrNotificationNotFound = NewAppError(
		ErrorTypeNotification,
		http.StatusNotFound,
		"The requested notification was not found",
		"notification_not_found",
		nil,
	)
)

// IsNotFound checks if the error is a NotFound error
func IsNotFound(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ErrorTypeNotFound || appErr.HTTPStatus == http.StatusNotFound
	}
	return false
}

// IsUnauthorized checks if the error is an Unauthorized error
func IsUnauthorized(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ErrorTypeUnauthorized || appErr.HTTPStatus == http.StatusUnauthorized
	}
	return false
}

// IsForbidden checks if the error is a Forbidden error
func IsForbidden(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ErrorTypeForbidden || appErr.HTTPStatus == http.StatusForbidden
	}
	return false
}

// IsConflict checks if the error is a Conflict error
func IsConflict(err error) bool {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.Type == ErrorTypeConflict || appErr.HTTPStatus == http.StatusConflict
	}
	return false
}

// NewValidationError creates a validation error with field-specific details
func NewValidationError(message string, fields FieldErrors) *AppError {
	err := &AppError{
		Type:       ErrorTypeInvalidInput,
		HTTPStatus: http.StatusBadRequest,
		Message:    message,
		Code:       "validation_error",
		Details:    fields,
	}
	return err
}