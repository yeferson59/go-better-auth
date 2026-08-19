package interfaces

import (
	"context"
	"time"
)

// Logger defines the interface for logging operations
type Logger interface {
	// Debug logs a debug message
	Debug(ctx context.Context, msg string, fields ...Field)

	// Info logs an info message
	Info(ctx context.Context, msg string, fields ...Field)

	// Warn logs a warning message
	Warn(ctx context.Context, msg string, fields ...Field)

	// Error logs an error message
	Error(ctx context.Context, msg string, fields ...Field)

	// Fatal logs a fatal message and exits
	Fatal(ctx context.Context, msg string, fields ...Field)

	// Panic logs a panic message and panics
	Panic(ctx context.Context, msg string, fields ...Field)

	// WithFields creates a new logger with the given fields
	WithFields(fields ...Field) Logger

	// WithContext creates a new logger with the given context
	WithContext(ctx context.Context) Logger

	// WithError creates a new logger with the given error
	WithError(err error) Logger

	// SetLevel sets the logging level
	SetLevel(level LogLevel)

	// GetLevel returns the current logging level
	GetLevel() LogLevel

	// IsEnabled returns true if the given level is enabled
	IsEnabled(level LogLevel) bool
}

// Field represents a logging field
type Field struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

// LogLevel represents the logging level
type LogLevel string

const (
	// LevelDebug debug level
	LevelDebug LogLevel = "debug"

	// LevelInfo info level
	LevelInfo LogLevel = "info"

	// LevelWarn warning level
	LevelWarn LogLevel = "warn"

	// LevelError error level
	LevelError LogLevel = "error"

	// LevelFatal fatal level
	LevelFatal LogLevel = "fatal"

	// LevelPanic panic level
	LevelPanic LogLevel = "panic"
)

// LogFormat represents the log format
type LogFormat string

const (
	// FormatJSON JSON format
	FormatJSON LogFormat = "json"

	// FormatText text format
	FormatText LogFormat = "text"

	// FormatConsole console format (colored text)
	FormatConsole LogFormat = "console"
)

// LogConfig represents logger configuration
type LogConfig struct {
	// Level is the minimum log level
	Level LogLevel `json:"level"`

	// Format is the log format
	Format LogFormat `json:"format"`

	// Output is the output destination
	Output string `json:"output"` // "stdout", "stderr", file path, etc.

	// TimeFormat is the time format
	TimeFormat string `json:"timeFormat"`

	// EnableCaller enables caller information
	EnableCaller bool `json:"enableCaller"`

	// EnableStacktrace enables stacktrace for errors
	EnableStacktrace bool `json:"enableStacktrace"`

	// MaxSize is the maximum size of log file in MB
	MaxSize int `json:"maxSize"`

	// MaxBackups is the maximum number of backup files
	MaxBackups int `json:"maxBackups"`

	// MaxAge is the maximum age of log files in days
	MaxAge int `json:"maxAge"`

	// Compress enables compression of backup files
	Compress bool `json:"compress"`

	// SamplingInitial is the initial sampling rate
	SamplingInitial int `json:"samplingInitial"`

	// SamplingThereafter is the sampling rate after initial
	SamplingThereafter int `json:"samplingThereafter"`

	// Development enables development mode
	Development bool `json:"development"`

	// DisableCaller disables caller information
	DisableCaller bool `json:"disableCaller"`

	// DisableStacktrace disables stacktrace
	DisableStacktrace bool `json:"disableStacktrace"`

	// DisableTimestamp disables timestamp
	DisableTimestamp bool `json:"disableTimestamp"`

	// ErrorOutputPaths is the error output paths
	ErrorOutputPaths []string `json:"errorOutputPaths"`

	// OutputPaths is the output paths
	OutputPaths []string `json:"outputPaths"`

	// InitialFields is the initial fields
	InitialFields map[string]any `json:"initialFields"`
}

// DefaultLogConfig provides default log configuration
var DefaultLogConfig = LogConfig{
	Level:              LevelInfo,
	Format:             FormatJSON,
	Output:             "stdout",
	TimeFormat:         time.RFC3339,
	EnableCaller:       true,
	EnableStacktrace:   true,
	MaxSize:            100,
	MaxBackups:         3,
	MaxAge:             28,
	Compress:           true,
	SamplingInitial:    100,
	SamplingThereafter: 100,
	Development:        true,
	DisableCaller:      false,
	DisableStacktrace:  false,
	DisableTimestamp:   false,
	ErrorOutputPaths:   []string{"stderr"},
	OutputPaths:        []string{"stdout"},
	InitialFields:      make(map[string]any),
}

// AuthLogger defines the interface for authentication-specific logging
type AuthLogger interface {
	Logger

	// LogUserLogin logs user login events
	LogUserLogin(ctx context.Context, userID string, success bool, fields ...Field)

	// LogUserLogout logs user logout events
	LogUserLogout(ctx context.Context, userID string, fields ...Field)

	// LogUserRegistration logs user registration events
	LogUserRegistration(ctx context.Context, userID string, success bool, fields ...Field)

	// LogPasswordChange logs password change events
	LogPasswordChange(ctx context.Context, userID string, success bool, fields ...Field)

	// LogPasswordReset logs password reset events
	LogPasswordReset(ctx context.Context, userID string, success bool, fields ...Field)

	// LogEmailVerification logs email verification events
	LogEmailVerification(ctx context.Context, userID string, success bool, fields ...Field)

	// LogTwoFactorAuth logs two-factor authentication events
	LogTwoFactorAuth(ctx context.Context, userID string, success bool, fields ...Field)

	// LogOAuthLogin logs OAuth login events
	LogOAuthLogin(ctx context.Context, userID string, provider string, success bool, fields ...Field)

	// LogRoleAssignment logs role assignment events
	LogRoleAssignment(ctx context.Context, userID string, roleID string, assignedBy string, fields ...Field)

	// LogPermissionGrant logs permission grant events
	LogPermissionGrant(ctx context.Context, userID string, permission string, grantedBy string, fields ...Field)

	// LogSecurityEvent logs security events
	LogSecurityEvent(ctx context.Context, event string, userID string, fields ...Field)

	// LogAuditEvent logs audit events
	LogAuditEvent(ctx context.Context, event string, userID string, fields ...Field)
}

// LoggerFactory defines the interface for creating loggers
type LoggerFactory interface {
	// Create creates a new logger with the given configuration
	Create(config LogConfig) (Logger, error)

	// CreateAuth creates a new authentication logger
	CreateAuth(config LogConfig) (AuthLogger, error)

	// GetDefault returns the default logger
	GetDefault() Logger

	// GetAuthDefault returns the default authentication logger
	GetAuthDefault() AuthLogger
}

// AuditLogger defines the interface for audit logging
type AuditLogger interface {
	// LogEvent logs an audit event
	LogEvent(ctx context.Context, event *AuditEvent) error

	// LogUserEvent logs a user-related audit event
	LogUserEvent(ctx context.Context, userID string, action string, details map[string]any) error

	// LogSystemEvent logs a system-related audit event
	LogSystemEvent(ctx context.Context, action string, details map[string]any) error

	// LogSecurityEvent logs a security-related audit event
	LogSecurityEvent(ctx context.Context, userID string, action string, details map[string]any) error

	// GetEvents retrieves audit events
	GetEvents(ctx context.Context, filter *AuditEventFilter) ([]*AuditEvent, error)

	// GetUserEvents retrieves audit events for a specific user
	GetUserEvents(ctx context.Context, userID string, filter *AuditEventFilter) ([]*AuditEvent, error)

	// GetSystemEvents retrieves system audit events
	GetSystemEvents(ctx context.Context, filter *AuditEventFilter) ([]*AuditEvent, error)

	// GetSecurityEvents retrieves security audit events
	GetSecurityEvents(ctx context.Context, filter *AuditEventFilter) ([]*AuditEvent, error)
}

// AuditEvent represents an audit event
type AuditEvent struct {
	ID        string         `json:"id"`
	Timestamp time.Time      `json:"timestamp"`
	UserID    string         `json:"userId,omitempty"`
	Action    string         `json:"action"`
	Resource  string         `json:"resource,omitempty"`
	Details   map[string]any `json:"details,omitempty"`
	IPAddress string         `json:"ipAddress,omitempty"`
	UserAgent string         `json:"userAgent,omitempty"`
	SessionID string         `json:"sessionId,omitempty"`
	RequestID string         `json:"requestId,omitempty"`
	Success   bool           `json:"success"`
	Error     string         `json:"error,omitempty"`
}

// AuditEventFilter represents audit event filter
type AuditEventFilter struct {
	UserID    string    `json:"userId,omitempty"`
	Action    string    `json:"action,omitempty"`
	Resource  string    `json:"resource,omitempty"`
	StartTime time.Time `json:"startTime,omitzero"`
	EndTime   time.Time `json:"endTime,omitzero"`
	Success   *bool     `json:"success,omitempty"`
	IPAddress string    `json:"ipAddress,omitempty"`
	Limit     int       `json:"limit,omitempty"`
	Offset    int       `json:"offset,omitempty"`
}

// SecurityLogger defines the interface for security logging
type SecurityLogger interface {
	// LogFailedLogin logs failed login attempts
	LogFailedLogin(ctx context.Context, email string, ipAddress string, userAgent string, reason string)

	// LogSuspiciousActivity logs suspicious activity
	LogSuspiciousActivity(ctx context.Context, userID string, activity string, details map[string]any)

	// LogAccountLockout logs account lockout events
	LogAccountLockout(ctx context.Context, userID string, reason string, duration time.Duration)

	// LogPasswordPolicyViolation logs password policy violations
	LogPasswordPolicyViolation(ctx context.Context, userID string, policy string, details map[string]any)

	// LogRateLimitExceeded logs rate limit exceeded events
	LogRateLimitExceeded(ctx context.Context, userID string, resource string, limit int, duration time.Duration)

	// LogUnauthorizedAccess logs unauthorized access attempts
	LogUnauthorizedAccess(ctx context.Context, userID string, resource string, action string, details map[string]any)

	// LogDataBreach logs data breach events
	LogDataBreach(ctx context.Context, details map[string]any)

	// LogComplianceViolation logs compliance violations
	LogComplianceViolation(ctx context.Context, userID string, regulation string, details map[string]any)
}

// Helper functions for creating fields
func String(key, value string) Field {
	return Field{Key: key, Value: value}
}

func Int(key string, value int) Field {
	return Field{Key: key, Value: value}
}

func Int64(key string, value int64) Field {
	return Field{Key: key, Value: value}
}

func Float64(key string, value float64) Field {
	return Field{Key: key, Value: value}
}

func Bool(key string, value bool) Field {
	return Field{Key: key, Value: value}
}

func Time(key string, value time.Time) Field {
	return Field{Key: key, Value: value}
}

func Duration(key string, value time.Duration) Field {
	return Field{Key: key, Value: value}
}

func Error(err error) Field {
	return Field{Key: "error", Value: err}
}

func Any(key string, value any) Field {
	return Field{Key: key, Value: value}
}
