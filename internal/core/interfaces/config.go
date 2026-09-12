package interfaces

import (
	"context"
	"time"
)

// Config defines the interface for configuration management
type Config interface {
	// Get retrieves a configuration value by key
	Get(key string) (any, error)

	// GetString retrieves a string configuration value
	GetString(key string) (string, error)

	// GetInt retrieves an integer configuration value
	GetInt(key string) (int, error)

	// GetInt64 retrieves an int64 configuration value
	GetInt64(key string) (int64, error)

	// GetFloat64 retrieves a float64 configuration value
	GetFloat64(key string) (float64, error)

	// GetBool retrieves a boolean configuration value
	GetBool(key string) (bool, error)

	// GetDuration retrieves a duration configuration value
	GetDuration(key string) (time.Duration, error)

	// GetStringSlice retrieves a string slice configuration value
	GetStringSlice(key string) ([]string, error)

	// GetStringMap retrieves a string map configuration value
	GetStringMap(key string) (map[string]string, error)

	// GetStringMapString retrieves a string map configuration value
	GetStringMapString(key string) (map[string]string, error)

	// Set sets a configuration value
	Set(key string, value any) error

	// SetDefault sets a default configuration value
	SetDefault(key string, value any)

	// IsSet checks if a configuration key is set
	IsSet(key string) bool

	// UnmarshalKey unmarshals a configuration key into a struct
	UnmarshalKey(key string, dest any) error

	// Unmarshal unmarshals the entire configuration into a struct
	Unmarshal(dest any) error

	// BindEnv binds a configuration key to an environment variable
	BindEnv(key string, envKey ...string) error

	// BindPFlag binds a configuration key to a pflags flag
	BindPFlag(key string, flag any) error

	// WatchConfig watches for configuration changes
	WatchConfig(callback func()) error

	// GetAll returns all configuration values
	GetAll() map[string]any

	// Reload reloads the configuration
	Reload() error

	// Validate validates the configuration
	Validate() error
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	// JWT Configuration
	JWT JWTConfig `json:"jwt" yaml:"jwt"`

	// Session Configuration
	Session SessionConfig `json:"session" yaml:"session"`

	// Password Configuration
	Password PasswordConfig `json:"password" yaml:"password"`

	// OAuth Configuration
	OAuth OAuthConfigMap `json:"oauth" yaml:"oauth"`

	// Email Configuration
	Email EmailConfig `json:"email" yaml:"email"`

	// Database Configuration
	Database DatabaseConfig `json:"database" yaml:"database"`

	// Redis Configuration
	Redis RedisConfig `json:"redis" yaml:"redis"`

	// Rate Limiting Configuration
	RateLimit RateLimitConfig `json:"rateLimit" yaml:"rate_limit"`

	// Security Configuration
	Security SecurityConfig `json:"security" yaml:"security"`

	// Logging Configuration
	Logging LoggingConfig `json:"logging" yaml:"logging"`

	// Server Configuration
	Server ServerConfig `json:"server" yaml:"server"`

	// Features Configuration
	Features FeaturesConfig `json:"features" yaml:"features"`
}

// JWTConfig represents JWT configuration
type JWTConfig struct {
	Secret            string        `json:"secret" yaml:"secret"`
	Algorithm         string        `json:"algorithm" yaml:"algorithm"`
	Issuer            string        `json:"issuer" yaml:"issuer"`
	Audience          string        `json:"audience" yaml:"audience"`
	Expiration        time.Duration `json:"expiration" yaml:"expiration"`
	RefreshExpiration time.Duration `json:"refreshExpiration" yaml:"refresh_expiration"`
	ClockSkew         time.Duration `json:"clockSkew" yaml:"clock_skew"`
}

// SessionConfig represents session configuration
type SessionConfig struct {
	Secret       string        `json:"secret" yaml:"secret"`
	Expiration   time.Duration `json:"expiration" yaml:"expiration"`
	CookieName   string        `json:"cookieName" yaml:"cookie_name"`
	CookiePath   string        `json:"cookiePath" yaml:"cookie_path"`
	CookieDomain string        `json:"cookieDomain" yaml:"cookie_domain"`
	Secure       bool          `json:"secure" yaml:"secure"`
	HTTPOnly     bool          `json:"httpOnly" yaml:"http_only"`
	SameSite     string        `json:"sameSite" yaml:"same_site"`
	Store        string        `json:"store" yaml:"store"` // "memory", "redis", "database"
}

// PasswordConfig represents password configuration
type PasswordConfig struct {
	Algorithm        string        `json:"algorithm" yaml:"algorithm"`
	BcryptCost       int           `json:"bcryptCost" yaml:"bcrypt_cost"`
	Argon2Memory     int           `json:"argon2Memory" yaml:"argon2_memory"`
	Argon2Time       int           `json:"argon2Time" yaml:"argon2_time"`
	Argon2Threads    int           `json:"argon2Threads" yaml:"argon2_threads"`
	Argon2SaltLength int           `json:"argon2SaltLength" yaml:"argon2_salt_length"`
	Argon2KeyLength  int           `json:"argon2KeyLength" yaml:"argon2_key_length"`
	MinLength        int           `json:"minLength" yaml:"min_length"`
	MaxLength        int           `json:"maxLength" yaml:"max_length"`
	RequireUppercase bool          `json:"requireUppercase" yaml:"require_uppercase"`
	RequireLowercase bool          `json:"requireLowercase" yaml:"require_lowercase"`
	RequireNumbers   bool          `json:"requireNumbers" yaml:"require_numbers"`
	RequireSymbols   bool          `json:"requireSymbols" yaml:"require_symbols"`
	HistoryCount     int           `json:"historyCount" yaml:"history_count"`
	ResetTokenExpiry time.Duration `json:"resetTokenExpiry" yaml:"reset_token_expiry"`
}

// OAuthConfigMap represents OAuth providers configuration
type OAuthConfigMap map[string]OAuthProviderConfig

// OAuthProviderConfig represents OAuth provider configuration
type OAuthProviderConfig struct {
	ClientID     string   `json:"clientId" yaml:"client_id"`
	ClientSecret string   `json:"clientSecret" yaml:"client_secret"`
	RedirectURL  string   `json:"redirectUrl" yaml:"redirect_url"`
	Scopes       []string `json:"scopes" yaml:"scopes"`
	Enabled      bool     `json:"enabled" yaml:"enabled"`
}

// EmailConfig represents email configuration
type EmailConfig struct {
	Provider       string        `json:"provider" yaml:"provider"`
	SMTPHost       string        `json:"smtpHost" yaml:"smtp_host"`
	SMTPPort       int           `json:"smtpPort" yaml:"smtp_port"`
	SMTPUsername   string        `json:"smtpUsername" yaml:"smtp_username"`
	SMTPPassword   string        `json:"smtpPassword" yaml:"smtp_password"`
	SMTPUseTLS     bool          `json:"smtpUseTls" yaml:"smtp_use_tls"`
	SMTPUseSSL     bool          `json:"smtpUseSsl" yaml:"smtp_use_ssl"`
	APIKey         string        `json:"apiKey" yaml:"api_key"`
	FromEmail      string        `json:"fromEmail" yaml:"from_email"`
	FromName       string        `json:"fromName" yaml:"from_name"`
	ReplyToEmail   string        `json:"replyToEmail" yaml:"reply_to_email"`
	TemplateEngine string        `json:"templateEngine" yaml:"template_engine"`
	TemplatePath   string        `json:"templatePath" yaml:"template_path"`
	RateLimit      int           `json:"rateLimit" yaml:"rate_limit"`
	Timeout        time.Duration `json:"timeout" yaml:"timeout"`
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Driver          string        `json:"driver" yaml:"driver"`
	Host            string        `json:"host" yaml:"host"`
	Port            int           `json:"port" yaml:"port"`
	Name            string        `json:"name" yaml:"name"`
	Username        string        `json:"username" yaml:"username"`
	Password        string        `json:"password" yaml:"password"`
	SSLMode         string        `json:"sslMode" yaml:"ssl_mode"`
	MaxOpenConns    int           `json:"maxOpenConns" yaml:"max_open_conns"`
	MaxIdleConns    int           `json:"maxIdleConns" yaml:"max_idle_conns"`
	ConnMaxLifetime time.Duration `json:"connMaxLifetime" yaml:"conn_max_lifetime"`
	ConnMaxIdleTime time.Duration `json:"connMaxIdleTime" yaml:"conn_max_idle_time"`
	MigrationsPath  string        `json:"migrationsPath" yaml:"migrations_path"`
	AutoMigrate     bool          `json:"autoMigrate" yaml:"auto_migrate"`
}

// RedisConfig represents Redis configuration
type RedisConfig struct {
	Host         string        `json:"host" yaml:"host"`
	Port         int           `json:"port" yaml:"port"`
	Password     string        `json:"password" yaml:"password"`
	DB           int           `json:"db" yaml:"db"`
	MaxRetries   int           `json:"maxRetries" yaml:"max_retries"`
	DialTimeout  time.Duration `json:"dialTimeout" yaml:"dial_timeout"`
	ReadTimeout  time.Duration `json:"readTimeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `json:"writeTimeout" yaml:"write_timeout"`
	PoolSize     int           `json:"poolSize" yaml:"pool_size"`
	MinIdleConns int           `json:"minIdleConns" yaml:"min_idle_conns"`
	MaxIdleConns int           `json:"maxIdleConns" yaml:"max_idle_conns"`
	TLSEnabled   bool          `json:"tlsEnabled" yaml:"tls_enabled"`
	Enabled      bool          `json:"enabled" yaml:"enabled"`
}

// RateLimitConfig represents rate limiting configuration
type RateLimitConfig struct {
	Enabled              bool          `json:"enabled" yaml:"enabled"`
	LoginAttempts        int           `json:"loginAttempts" yaml:"login_attempts"`
	LoginWindow          time.Duration `json:"loginWindow" yaml:"login_window"`
	PasswordResets       int           `json:"passwordResets" yaml:"password_resets"`
	PasswordResetWindow  time.Duration `json:"passwordResetWindow" yaml:"password_reset_window"`
	RegistrationAttempts int           `json:"registrationAttempts" yaml:"registration_attempts"`
	RegistrationWindow   time.Duration `json:"registrationWindow" yaml:"registration_window"`
	GlobalRateLimit      int           `json:"globalRateLimit" yaml:"global_rate_limit"`
	GlobalWindow         time.Duration `json:"globalWindow" yaml:"global_window"`
}

// SecurityConfig represents security configuration
type SecurityConfig struct {
	AccountLockout    AccountLockoutConfig    `json:"accountLockout" yaml:"account_lockout"`
	TwoFactorAuth     TwoFactorAuthConfig     `json:"twoFactorAuth" yaml:"two_factor_auth"`
	EmailVerification EmailVerificationConfig `json:"emailVerification" yaml:"email_verification"`
	PasswordHistory   PasswordHistoryConfig   `json:"passwordHistory" yaml:"password_history"`
	SessionSecurity   SessionSecurityConfig   `json:"sessionSecurity" yaml:"session_security"`
	IPWhitelist       []string                `json:"ipWhitelist" yaml:"ip_whitelist"`
	IPBlacklist       []string                `json:"ipBlacklist" yaml:"ip_blacklist"`
	TrustedProxies    []string                `json:"trustedProxies" yaml:"trusted_proxies"`
	CSRFProtection    bool                    `json:"csrfProtection" yaml:"csrf_protection"`
	CORSEnabled       bool                    `json:"corsEnabled" yaml:"cors_enabled"`
	CORSOrigins       []string                `json:"corsOrigins" yaml:"cors_origins"`
}

// AccountLockoutConfig represents account lockout configuration
type AccountLockoutConfig struct {
	Enabled         bool          `json:"enabled" yaml:"enabled"`
	MaxAttempts     int           `json:"maxAttempts" yaml:"max_attempts"`
	LockoutDuration time.Duration `json:"lockoutDuration" yaml:"lockout_duration"`
	ResetOnSuccess  bool          `json:"resetOnSuccess" yaml:"reset_on_success"`
}

// TwoFactorAuthConfig represents 2FA configuration
type TwoFactorAuthConfig struct {
	Enabled     bool   `json:"enabled" yaml:"enabled"`
	Issuer      string `json:"issuer" yaml:"issuer"`
	Required    bool   `json:"required" yaml:"required"`
	BackupCodes int    `json:"backupCodes" yaml:"backup_codes"`
}

// EmailVerificationConfig represents email verification configuration
type EmailVerificationConfig struct {
	Enabled     bool          `json:"enabled" yaml:"enabled"`
	Required    bool          `json:"required" yaml:"required"`
	TokenExpiry time.Duration `json:"tokenExpiry" yaml:"token_expiry"`
}

// PasswordHistoryConfig represents password history configuration
type PasswordHistoryConfig struct {
	Enabled bool `json:"enabled" yaml:"enabled"`
	Count   int  `json:"count" yaml:"count"`
}

// SessionSecurityConfig represents session security configuration
type SessionSecurityConfig struct {
	MaxConcurrentSessions int           `json:"maxConcurrentSessions" yaml:"max_concurrent_sessions"`
	SessionTimeout        time.Duration `json:"sessionTimeout" yaml:"session_timeout"`
	IdleTimeout           time.Duration `json:"idleTimeout" yaml:"idle_timeout"`
	RememberMe            bool          `json:"rememberMe" yaml:"remember_me"`
	RememberMeDuration    time.Duration `json:"rememberMeDuration" yaml:"remember_me_duration"`
}

// LoggingConfig represents logging configuration
type LoggingConfig struct {
	Level    string `json:"level" yaml:"level"`
	Format   string `json:"format" yaml:"format"`
	Output   string `json:"output" yaml:"output"`
	Audit    bool   `json:"audit" yaml:"audit"`
	Security bool   `json:"security" yaml:"security"`
}

// ServerConfig represents server configuration
type ServerConfig struct {
	Host             string        `json:"host" yaml:"host"`
	Port             int           `json:"port" yaml:"port"`
	ReadTimeout      time.Duration `json:"readTimeout" yaml:"read_timeout"`
	WriteTimeout     time.Duration `json:"writeTimeout" yaml:"write_timeout"`
	IdleTimeout      time.Duration `json:"idleTimeout" yaml:"idle_timeout"`
	MaxHeaderBytes   int           `json:"maxHeaderBytes" yaml:"max_header_bytes"`
	TLSEnabled       bool          `json:"tlsEnabled" yaml:"tls_enabled"`
	TLSCertFile      string        `json:"tlsCertFile" yaml:"tls_cert_file"`
	TLSKeyFile       string        `json:"tlsKeyFile" yaml:"tls_key_file"`
	GracefulShutdown time.Duration `json:"gracefulShutdown" yaml:"graceful_shutdown"`
}

// FeaturesConfig represents feature flags configuration
type FeaturesConfig struct {
	Registration      bool `json:"registration" yaml:"registration"`
	OAuthLogin        bool `json:"oauthLogin" yaml:"oauth_login"`
	PasswordReset     bool `json:"passwordReset" yaml:"password_reset"`
	EmailVerification bool `json:"emailVerification" yaml:"email_verification"`
	TwoFactorAuth     bool `json:"twoFactorAuth" yaml:"two_factor_auth"`
	AccountLockout    bool `json:"accountLockout" yaml:"account_lockout"`
	RateLimiting      bool `json:"rateLimiting" yaml:"rate_limiting"`
	AuditLogging      bool `json:"auditLogging" yaml:"audit_logging"`
	UserProfiles      bool `json:"userProfiles" yaml:"user_profiles"`
	RoleBasedAccess   bool `json:"roleBasedAccess" yaml:"role_based_access"`
}

// ConfigManager defines the interface for configuration management
type ConfigManager interface {
	// LoadConfig loads configuration from multiple sources
	LoadConfig(sources ...ConfigSource) error

	// GetConfig returns the current configuration
	GetConfig() *AuthConfig

	// ValidateConfig validates the configuration
	ValidateConfig() error

	// ReloadConfig reloads the configuration
	ReloadConfig() error

	// WatchConfig watches for configuration changes
	WatchConfig(callback func(*AuthConfig)) error

	// GetConfigByKey gets configuration by key path
	GetConfigByKey(key string) (any, error)

	// SetConfigByKey sets configuration by key path
	SetConfigByKey(key string, value any) error
}

// ConfigSource represents a configuration source
type ConfigSource interface {
	// Load loads configuration from the source
	Load() (map[string]any, error)

	// Watch watches for configuration changes
	Watch(callback func(map[string]any)) error

	// GetName returns the source name
	GetName() string

	// GetPriority returns the source priority
	GetPriority() int
}

// ConfigValidator defines the interface for configuration validation
type ConfigValidator interface {
	// Validate validates the configuration
	Validate(config *AuthConfig) error

	// ValidateSection validates a specific configuration section
	ValidateSection(section string, config any) error

	// GetValidationRules returns validation rules
	GetValidationRules() map[string]any
}

// ConfigEncryption defines the interface for configuration encryption
type ConfigEncryption interface {
	// Encrypt encrypts sensitive configuration values
	Encrypt(value string) (string, error)

	// Decrypt decrypts sensitive configuration values
	Decrypt(encryptedValue string) (string, error)

	// IsEncrypted checks if a value is encrypted
	IsEncrypted(value string) bool
}

// ConfigProvider defines the interface for configuration providers
type ConfigProvider interface {
	// Get gets a configuration value
	Get(ctx context.Context, key string) (any, error)

	// Set sets a configuration value
	Set(ctx context.Context, key string, value any) error

	// Delete deletes a configuration value
	Delete(ctx context.Context, key string) error

	// List lists all configuration keys
	List(ctx context.Context, prefix string) ([]string, error)

	// Watch watches for configuration changes
	Watch(ctx context.Context, key string, callback func(string, any)) error

	// Close closes the provider
	Close() error
}

// Default configurations
var (
	DefaultJWTConfig = JWTConfig{
		Algorithm:         "HS256",
		Expiration:        15 * time.Minute,
		RefreshExpiration: 7 * 24 * time.Hour,
		ClockSkew:         5 * time.Minute,
	}

	DefaultSessionConfig = SessionConfig{
		Expiration: 24 * time.Hour,
		CookieName: "session",
		CookiePath: "/",
		Secure:     false,
		HTTPOnly:   true,
		SameSite:   "lax",
		Store:      "memory",
	}

	DefaultPasswordConfig = PasswordConfig{
		Algorithm:        "bcrypt",
		BcryptCost:       12,
		MinLength:        8,
		MaxLength:        128,
		RequireUppercase: true,
		RequireLowercase: true,
		RequireNumbers:   true,
		RequireSymbols:   true,
		HistoryCount:     5,
		ResetTokenExpiry: 15 * time.Minute,
	}

	DefaultRateLimitConfig = RateLimitConfig{
		Enabled:              true,
		LoginAttempts:        5,
		LoginWindow:          15 * time.Minute,
		PasswordResets:       3,
		PasswordResetWindow:  time.Hour,
		RegistrationAttempts: 3,
		RegistrationWindow:   time.Hour,
		GlobalRateLimit:      100,
		GlobalWindow:         time.Minute,
	}

	DefaultSecurityConfig = SecurityConfig{
		AccountLockout: AccountLockoutConfig{
			Enabled:         true,
			MaxAttempts:     5,
			LockoutDuration: 30 * time.Minute,
			ResetOnSuccess:  true,
		},
		TwoFactorAuth: TwoFactorAuthConfig{
			Enabled:     false,
			Required:    false,
			BackupCodes: 10,
		},
		EmailVerification: EmailVerificationConfig{
			Enabled:     false,
			Required:    false,
			TokenExpiry: 24 * time.Hour,
		},
		PasswordHistory: PasswordHistoryConfig{
			Enabled: true,
			Count:   5,
		},
		SessionSecurity: SessionSecurityConfig{
			MaxConcurrentSessions: 3,
			SessionTimeout:        24 * time.Hour,
			IdleTimeout:           30 * time.Minute,
			RememberMe:            true,
			RememberMeDuration:    30 * 24 * time.Hour,
		},
		CSRFProtection: true,
		CORSEnabled:    true,
	}
)
