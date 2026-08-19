package errors

import (
	"fmt"
	"net/http"
)

// Error represents a domain error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
	Status  int    `json:"status"`
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s (%s)", e.Code, e.Message, e.Details)
	}

	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Is reports whether target is a domain error carrying the same code, which lets
// callers match on the kind of failure with errors.Is without caring about the
// identifier baked into Details.
func (e *Error) Is(target error) bool {
	other, ok := target.(*Error)
	if !ok {
		return false
	}

	return e.Code == other.Code
}

// Common error codes
const (
	// Authentication errors
	ErrInvalidCredentials    = "INVALID_CREDENTIALS"
	ErrUserNotFound          = "USER_NOT_FOUND"
	ErrUserAlreadyExists     = "USER_ALREADY_EXISTS"
	ErrEmailAlreadyExists    = "EMAIL_ALREADY_EXISTS"
	ErrUsernameAlreadyExists = "USERNAME_ALREADY_EXISTS"
	ErrAccountLocked         = "ACCOUNT_LOCKED"
	ErrAccountDisabled       = "ACCOUNT_DISABLED"
	ErrEmailNotVerified      = "EMAIL_NOT_VERIFIED"
	ErrInvalidToken          = "INVALID_TOKEN"
	ErrTokenExpired          = "TOKEN_EXPIRED"
	ErrTokenRevoked          = "TOKEN_REVOKED"
	ErrSessionExpired        = "SESSION_EXPIRED"
	ErrSessionNotFound       = "SESSION_NOT_FOUND"
	ErrInvalidSession        = "INVALID_SESSION"

	// Authorization errors
	ErrUnauthorized                 = "UNAUTHORIZED"
	ErrForbidden                    = "FORBIDDEN"
	ErrInsufficientPermissions      = "INSUFFICIENT_PERMISSIONS"
	ErrRoleNotFound                 = "ROLE_NOT_FOUND"
	ErrPermissionNotFound           = "PERMISSION_NOT_FOUND"
	ErrRoleAlreadyExists            = "ROLE_ALREADY_EXISTS"
	ErrPermissionAlreadyExists      = "PERMISSION_ALREADY_EXISTS"
	ErrCannotDeleteSystemRole       = "CANNOT_DELETE_SYSTEM_ROLE"
	ErrCannotDeleteSystemPermission = "CANNOT_DELETE_SYSTEM_PERMISSION"

	// Password errors
	ErrPasswordTooWeak            = "PASSWORD_TOO_WEAK"
	ErrPasswordTooShort           = "PASSWORD_TOO_SHORT"
	ErrPasswordTooLong            = "PASSWORD_TOO_LONG"
	ErrPasswordReused             = "PASSWORD_REUSED"
	ErrPasswordInvalid            = "PASSWORD_INVALID"
	ErrPasswordMismatch           = "PASSWORD_MISMATCH"
	ErrPasswordResetRequired      = "PASSWORD_RESET_REQUIRED"
	ErrPasswordResetTokenNotFound = "PASSWORD_RESET_TOKEN_NOT_FOUND"
	ErrPasswordResetTokenExpired  = "PASSWORD_RESET_TOKEN_EXPIRED"
	ErrPasswordResetTokenUsed     = "PASSWORD_RESET_TOKEN_USED"

	// OAuth errors
	ErrOAuthProviderNotFound     = "OAUTH_PROVIDER_NOT_FOUND"
	ErrOAuthProviderDisabled     = "OAUTH_PROVIDER_DISABLED"
	ErrOAuthInvalidState         = "OAUTH_INVALID_STATE"
	ErrOAuthInvalidCode          = "OAUTH_INVALID_CODE"
	ErrOAuthTokenExpired         = "OAUTH_TOKEN_EXPIRED"
	ErrOAuthAccountNotLinked     = "OAUTH_ACCOUNT_NOT_LINKED"
	ErrOAuthAccountAlreadyLinked = "OAUTH_ACCOUNT_ALREADY_LINKED"

	// Rate limiting errors
	ErrRateLimitExceeded     = "RATE_LIMIT_EXCEEDED"
	ErrTooManyRequests       = "TOO_MANY_REQUESTS"
	ErrTooManyFailedAttempts = "TOO_MANY_FAILED_ATTEMPTS"

	// Validation errors
	ErrValidationFailed     = "VALIDATION_FAILED"
	ErrInvalidInput         = "INVALID_INPUT"
	ErrInvalidEmail         = "INVALID_EMAIL"
	ErrInvalidUsername      = "INVALID_USERNAME"
	ErrInvalidPhoneNumber   = "INVALID_PHONE_NUMBER"
	ErrRequiredFieldMissing = "REQUIRED_FIELD_MISSING"
	ErrFieldTooLong         = "FIELD_TOO_LONG"
	ErrFieldTooShort        = "FIELD_TOO_SHORT"

	// System errors
	ErrInternalServer     = "INTERNAL_SERVER_ERROR"
	ErrDatabaseError      = "DATABASE_ERROR"
	ErrCacheError         = "CACHE_ERROR"
	ErrMailerError        = "MAILER_ERROR"
	ErrConfigurationError = "CONFIGURATION_ERROR"
	ErrNotImplemented     = "NOT_IMPLEMENTED"
	ErrServiceUnavailable = "SERVICE_UNAVAILABLE"
	ErrTimeout            = "TIMEOUT"

	// Two-factor authentication errors
	ErrTwoFactorRequired       = "TWO_FACTOR_REQUIRED"
	ErrInvalidTwoFactorCode    = "INVALID_TWO_FACTOR_CODE"
	ErrTwoFactorNotEnabled     = "TWO_FACTOR_NOT_ENABLED"
	ErrTwoFactorAlreadyEnabled = "TWO_FACTOR_ALREADY_ENABLED"
	ErrInvalidBackupCode       = "INVALID_BACKUP_CODE"
	ErrBackupCodeAlreadyUsed   = "BACKUP_CODE_ALREADY_USED"

	// Invitation errors
	ErrInvitationNotFound    = "INVITATION_NOT_FOUND"
	ErrInvitationExpired     = "INVITATION_EXPIRED"
	ErrInvitationAlreadyUsed = "INVITATION_ALREADY_USED"
	ErrInvitationInvalid     = "INVITATION_INVALID"

	// Profile errors
	ErrProfileNotFound        = "PROFILE_NOT_FOUND"
	ErrAvatarTooLarge         = "AVATAR_TOO_LARGE"
	ErrUnsupportedImageFormat = "UNSUPPORTED_IMAGE_FORMAT"
)

// Error constructors
func NewInvalidCredentialsError() *Error {
	return &Error{
		Code:    ErrInvalidCredentials,
		Message: "Invalid email or password",
		Status:  http.StatusUnauthorized,
	}
}

func NewUserNotFoundError(identifier string) *Error {
	return &Error{
		Code:    ErrUserNotFound,
		Message: "User not found",
		Details: fmt.Sprintf("User with identifier '%s' not found", identifier),
		Status:  http.StatusNotFound,
	}
}

func NewUserAlreadyExistsError(email string) *Error {
	return &Error{
		Code:    ErrUserAlreadyExists,
		Message: "User already exists",
		Details: fmt.Sprintf("User with email '%s' already exists", email),
		Status:  http.StatusConflict,
	}
}

func NewEmailAlreadyExistsError(email string) *Error {
	return &Error{
		Code:    ErrEmailAlreadyExists,
		Message: "Email already exists",
		Details: fmt.Sprintf("Email '%s' is already registered", email),
		Status:  http.StatusConflict,
	}
}

func NewUsernameAlreadyExistsError(username string) *Error {
	return &Error{
		Code:    ErrUsernameAlreadyExists,
		Message: "Username already exists",
		Details: fmt.Sprintf("Username '%s' is already taken", username),
		Status:  http.StatusConflict,
	}
}

func NewAccountLockedError(unlockTime string) *Error {
	return &Error{
		Code:    ErrAccountLocked,
		Message: "Account is locked",
		Details: fmt.Sprintf("Account is locked until %s", unlockTime),
		Status:  http.StatusLocked,
	}
}

func NewAccountDisabledError() *Error {
	return &Error{
		Code:    ErrAccountDisabled,
		Message: "Account is disabled",
		Status:  http.StatusForbidden,
	}
}

func NewEmailNotVerifiedError() *Error {
	return &Error{
		Code:    ErrEmailNotVerified,
		Message: "Email not verified",
		Details: "Please verify your email address before continuing",
		Status:  http.StatusForbidden,
	}
}

func NewInvalidTokenError(tokenType string) *Error {
	return &Error{
		Code:    ErrInvalidToken,
		Message: "Invalid token",
		Details: fmt.Sprintf("Invalid %s token", tokenType),
		Status:  http.StatusUnauthorized,
	}
}

func NewTokenExpiredError(tokenType string) *Error {
	return &Error{
		Code:    ErrTokenExpired,
		Message: "Token expired",
		Details: fmt.Sprintf("%s token has expired", tokenType),
		Status:  http.StatusUnauthorized,
	}
}

func NewTokenRevokedError(tokenType string) *Error {
	return &Error{
		Code:    ErrTokenRevoked,
		Message: "Token revoked",
		Details: fmt.Sprintf("%s token has been revoked", tokenType),
		Status:  http.StatusUnauthorized,
	}
}

func NewSessionExpiredError() *Error {
	return &Error{
		Code:    ErrSessionExpired,
		Message: "Session expired",
		Details: "Your session has expired, please login again",
		Status:  http.StatusUnauthorized,
	}
}

func NewSessionNotFoundError() *Error {
	return &Error{
		Code:    ErrSessionNotFound,
		Message: "Session not found",
		Status:  http.StatusUnauthorized,
	}
}

func NewUnauthorizedError() *Error {
	return &Error{
		Code:    ErrUnauthorized,
		Message: "Unauthorized",
		Status:  http.StatusUnauthorized,
	}
}

func NewForbiddenError() *Error {
	return &Error{
		Code:    ErrForbidden,
		Message: "Forbidden",
		Status:  http.StatusForbidden,
	}
}

func NewInsufficientPermissionsError(permission string) *Error {
	return &Error{
		Code:    ErrInsufficientPermissions,
		Message: "Insufficient permissions",
		Details: fmt.Sprintf("Permission '%s' is required", permission),
		Status:  http.StatusForbidden,
	}
}

func NewPasswordTooWeakError(requirements string) *Error {
	return &Error{
		Code:    ErrPasswordTooWeak,
		Message: "Password too weak",
		Details: fmt.Sprintf("Password must meet the following requirements: %s", requirements),
		Status:  http.StatusBadRequest,
	}
}

func NewPasswordTooShortError(minLength int) *Error {
	return &Error{
		Code:    ErrPasswordTooShort,
		Message: "Password too short",
		Details: fmt.Sprintf("Password must be at least %d characters long", minLength),
		Status:  http.StatusBadRequest,
	}
}

func NewPasswordTooLongError(maxLength int) *Error {
	return &Error{
		Code:    ErrPasswordTooLong,
		Message: "Password too long",
		Details: fmt.Sprintf("Password must be no more than %d characters long", maxLength),
		Status:  http.StatusBadRequest,
	}
}

func NewPasswordReusedError() *Error {
	return &Error{
		Code:    ErrPasswordReused,
		Message: "Password reused",
		Details: "This password has been used recently and cannot be reused",
		Status:  http.StatusBadRequest,
	}
}

func NewRateLimitExceededError(retryAfter string) *Error {
	return &Error{
		Code:    ErrRateLimitExceeded,
		Message: "Rate limit exceeded",
		Details: fmt.Sprintf("Too many requests. Please try again after %s", retryAfter),
		Status:  http.StatusTooManyRequests,
	}
}

func NewValidationFailedError(field, reason string) *Error {
	return &Error{
		Code:    ErrValidationFailed,
		Message: "Validation failed",
		Details: fmt.Sprintf("Field '%s': %s", field, reason),
		Status:  http.StatusBadRequest,
	}
}

func NewInvalidInputError(field string) *Error {
	return &Error{
		Code:    ErrInvalidInput,
		Message: "Invalid input",
		Details: fmt.Sprintf("Invalid value for field '%s'", field),
		Status:  http.StatusBadRequest,
	}
}

func NewInvalidEmailError(email string) *Error {
	return &Error{
		Code:    ErrInvalidEmail,
		Message: "Invalid email address",
		Details: fmt.Sprintf("Email address '%s' is not valid", email),
		Status:  http.StatusBadRequest,
	}
}

func NewInternalServerError(details string) *Error {
	return &Error{
		Code:    ErrInternalServer,
		Message: "Internal server error",
		Details: details,
		Status:  http.StatusInternalServerError,
	}
}

func NewDatabaseError(operation string, err error) *Error {
	return &Error{
		Code:    ErrDatabaseError,
		Message: "Database error",
		Details: fmt.Sprintf("Database operation '%s' failed: %s", operation, err.Error()),
		Status:  http.StatusInternalServerError,
	}
}

func NewConfigurationError(setting string, reason string) *Error {
	return &Error{
		Code:    ErrConfigurationError,
		Message: "Configuration error",
		Details: fmt.Sprintf("Setting '%s': %s", setting, reason),
		Status:  http.StatusInternalServerError,
	}
}

func NewNotImplementedError(feature string) *Error {
	return &Error{
		Code:    ErrNotImplemented,
		Message: "Not implemented",
		Details: fmt.Sprintf("Feature '%s' is not implemented", feature),
		Status:  http.StatusNotImplemented,
	}
}

func NewServiceUnavailableError(service string) *Error {
	return &Error{
		Code:    ErrServiceUnavailable,
		Message: "Service unavailable",
		Details: fmt.Sprintf("Service '%s' is currently unavailable", service),
		Status:  http.StatusServiceUnavailable,
	}
}

func NewTimeoutError(operation string) *Error {
	return &Error{
		Code:    ErrTimeout,
		Message: "Operation timeout",
		Details: fmt.Sprintf("Operation '%s' timed out", operation),
		Status:  http.StatusRequestTimeout,
	}
}

func NewOAuthProviderNotFoundError(provider string) *Error {
	return &Error{
		Code:    ErrOAuthProviderNotFound,
		Message: "OAuth provider not found",
		Details: fmt.Sprintf("OAuth provider '%s' not found", provider),
		Status:  http.StatusNotFound,
	}
}

func NewOAuthProviderDisabledError(provider string) *Error {
	return &Error{
		Code:    ErrOAuthProviderDisabled,
		Message: "OAuth provider disabled",
		Details: fmt.Sprintf("OAuth provider '%s' is disabled", provider),
		Status:  http.StatusForbidden,
	}
}

func NewTwoFactorRequiredError() *Error {
	return &Error{
		Code:    ErrTwoFactorRequired,
		Message: "Two-factor authentication required",
		Status:  http.StatusUnauthorized,
	}
}

func NewInvalidTwoFactorCodeError() *Error {
	return &Error{
		Code:    ErrInvalidTwoFactorCode,
		Message: "Invalid two-factor authentication code",
		Status:  http.StatusUnauthorized,
	}
}

func NewRoleNotFoundError(roleID string) *Error {
	return &Error{
		Code:    ErrRoleNotFound,
		Message: "Role not found",
		Details: fmt.Sprintf("Role with ID '%s' not found", roleID),
		Status:  http.StatusNotFound,
	}
}

func NewPermissionNotFoundError(permissionID string) *Error {
	return &Error{
		Code:    ErrPermissionNotFound,
		Message: "Permission not found",
		Details: fmt.Sprintf("Permission with ID '%s' not found", permissionID),
		Status:  http.StatusNotFound,
	}
}

func NewCannotDeleteSystemRoleError(roleName string) *Error {
	return &Error{
		Code:    ErrCannotDeleteSystemRole,
		Message: "Cannot delete system role",
		Details: fmt.Sprintf("System role '%s' cannot be deleted", roleName),
		Status:  http.StatusForbidden,
	}
}

func NewCannotDeleteSystemPermissionError(permissionName string) *Error {
	return &Error{
		Code:    ErrCannotDeleteSystemPermission,
		Message: "Cannot delete system permission",
		Details: fmt.Sprintf("System permission '%s' cannot be deleted", permissionName),
		Status:  http.StatusForbidden,
	}
}

// Helper functions to check error types
func IsAuthenticationError(err error) bool {
	if domainErr, ok := err.(*Error); ok {
		return domainErr.Code == ErrInvalidCredentials ||
			domainErr.Code == ErrUserNotFound ||
			domainErr.Code == ErrAccountLocked ||
			domainErr.Code == ErrAccountDisabled ||
			domainErr.Code == ErrEmailNotVerified ||
			domainErr.Code == ErrInvalidToken ||
			domainErr.Code == ErrTokenExpired ||
			domainErr.Code == ErrTokenRevoked ||
			domainErr.Code == ErrSessionExpired ||
			domainErr.Code == ErrSessionNotFound
	}

	return false
}

func IsAuthorizationError(err error) bool {
	if domainErr, ok := err.(*Error); ok {
		return domainErr.Code == ErrUnauthorized ||
			domainErr.Code == ErrForbidden ||
			domainErr.Code == ErrInsufficientPermissions
	}

	return false
}

func IsValidationError(err error) bool {
	if domainErr, ok := err.(*Error); ok {
		return domainErr.Code == ErrValidationFailed ||
			domainErr.Code == ErrInvalidInput ||
			domainErr.Code == ErrInvalidEmail ||
			domainErr.Code == ErrInvalidUsername ||
			domainErr.Code == ErrRequiredFieldMissing ||
			domainErr.Code == ErrFieldTooLong ||
			domainErr.Code == ErrFieldTooShort
	}

	return false
}

func IsSystemError(err error) bool {
	if domainErr, ok := err.(*Error); ok {
		return domainErr.Code == ErrInternalServer ||
			domainErr.Code == ErrDatabaseError ||
			domainErr.Code == ErrCacheError ||
			domainErr.Code == ErrMailerError ||
			domainErr.Code == ErrConfigurationError ||
			domainErr.Code == ErrServiceUnavailable ||
			domainErr.Code == ErrTimeout
	}

	return false
}

func IsRateLimitError(err error) bool {
	if domainErr, ok := err.(*Error); ok {
		return domainErr.Code == ErrRateLimitExceeded ||
			domainErr.Code == ErrTooManyRequests ||
			domainErr.Code == ErrTooManyFailedAttempts
	}

	return false
}
