package config

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/caarlos0/env/v10"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"

	"github.com/yeferson59/go-better-auth/internal/core/interfaces"
)

// EnvConfig represents environment-based configuration
type EnvConfig struct {
	// JWT Configuration
	JWTSecret            string        `env:"JWT_SECRET,required"`
	JWTAlgorithm         string        `env:"JWT_ALGORITHM" envDefault:"HS256"`
	JWTIssuer            string        `env:"JWT_ISSUER" envDefault:"go-better-auth"`
	JWTAudience          string        `env:"JWT_AUDIENCE" envDefault:"go-better-auth"`
	JWTExpiration        time.Duration `env:"JWT_EXPIRATION" envDefault:"15m"`
	JWTRefreshExpiration time.Duration `env:"JWT_REFRESH_EXPIRATION" envDefault:"168h"` // 7 days
	JWTClockSkew         time.Duration `env:"JWT_CLOCK_SKEW" envDefault:"5m"`

	// Session Configuration
	SessionSecret     string        `env:"SESSION_SECRET,required"`
	SessionExpiration time.Duration `env:"SESSION_EXPIRATION" envDefault:"24h"`
	SessionCookieName string        `env:"SESSION_COOKIE_NAME" envDefault:"session"`
	SessionSecure     bool          `env:"SESSION_SECURE" envDefault:"true"`
	SessionHttpOnly   bool          `env:"SESSION_HTTP_ONLY" envDefault:"true"`

	// Password Configuration
	PasswordAlgorithm     string `env:"PASSWORD_ALGORITHM" envDefault:"bcrypt"`
	PasswordBcryptCost    int    `env:"PASSWORD_BCRYPT_COST" envDefault:"12"`
	PasswordMinLength     int    `env:"PASSWORD_MIN_LENGTH" envDefault:"8"`
	PasswordMaxLength     int    `env:"PASSWORD_MAX_LENGTH" envDefault:"128"`
	PasswordRequireUpper  bool   `env:"PASSWORD_REQUIRE_UPPER" envDefault:"true"`
	PasswordRequireLower  bool   `env:"PASSWORD_REQUIRE_LOWER" envDefault:"true"`
	PasswordRequireDigit  bool   `env:"PASSWORD_REQUIRE_DIGIT" envDefault:"true"`
	PasswordRequireSymbol bool   `env:"PASSWORD_REQUIRE_SYMBOL" envDefault:"false"`

	// Database Configuration
	DatabaseDriver      string        `env:"DATABASE_DRIVER" envDefault:"sqlite"`
	DatabaseHost        string        `env:"DATABASE_HOST" envDefault:"localhost"`
	DatabasePort        int           `env:"DATABASE_PORT" envDefault:"5432"`
	DatabaseName        string        `env:"DATABASE_NAME" envDefault:"auth.db"`
	DatabaseUser        string        `env:"DATABASE_USER"`
	DatabasePassword    string        `env:"DATABASE_PASSWORD"`
	DatabaseSSLMode     string        `env:"DATABASE_SSL_MODE" envDefault:"require"`
	DatabaseMaxConns    int           `env:"DATABASE_MAX_CONNS" envDefault:"10"`
	DatabaseIdleTimeout time.Duration `env:"DATABASE_IDLE_TIMEOUT" envDefault:"30m"`

	// Redis Configuration
	RedisEnabled     bool          `env:"REDIS_ENABLED" envDefault:"false"`
	RedisHost        string        `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort        int           `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword    string        `env:"REDIS_PASSWORD"`
	RedisDB          int           `env:"REDIS_DB" envDefault:"0"`
	RedisPoolSize    int           `env:"REDIS_POOL_SIZE" envDefault:"10"`
	RedisDialTimeout time.Duration `env:"REDIS_DIAL_TIMEOUT" envDefault:"5s"`

	// Email Configuration
	EmailProvider    string `env:"EMAIL_PROVIDER" envDefault:"smtp"`
	EmailSMTPHost    string `env:"EMAIL_SMTP_HOST"`
	EmailSMTPPort    int    `env:"EMAIL_SMTP_PORT" envDefault:"587"`
	EmailSMTPUser    string `env:"EMAIL_SMTP_USER"`
	EmailSMTPPass    string `env:"EMAIL_SMTP_PASS"`
	EmailFromAddress string `env:"EMAIL_FROM_ADDRESS"`
	EmailFromName    string `env:"EMAIL_FROM_NAME"`

	// Security Configuration
	// EncryptionKey has no consumer yet; it stays optional until something uses it.
	EncryptionKey string `env:"ENCRYPTION_KEY"`
	// CORSAllowedOrigins defaults to empty: origins must be opted into explicitly
	// rather than defaulting to a wildcard.
	CORSAllowedOrigins string `env:"CORS_ALLOWED_ORIGINS"`
	CSRFEnabled        bool   `env:"CSRF_ENABLED" envDefault:"true"`
	RateLimitEnabled   bool   `env:"RATE_LIMIT_ENABLED" envDefault:"true"`

	// OAuth Configuration
	GoogleClientID     string `env:"GOOGLE_CLIENT_ID"`
	GoogleClientSecret string `env:"GOOGLE_CLIENT_SECRET"`
	GitHubClientID     string `env:"GITHUB_CLIENT_ID"`
	GitHubClientSecret string `env:"GITHUB_CLIENT_SECRET"`

	// Server Configuration
	ServerHost string `env:"SERVER_HOST" envDefault:"localhost"`
	ServerPort int    `env:"SERVER_PORT" envDefault:"8080"`
	ServerTLS  bool   `env:"SERVER_TLS" envDefault:"false"`

	// Feature Flags
	FeatureRegistration      bool `env:"FEATURE_REGISTRATION" envDefault:"true"`
	FeaturePasswordReset     bool `env:"FEATURE_PASSWORD_RESET" envDefault:"true"`
	FeatureEmailVerification bool `env:"FEATURE_EMAIL_VERIFICATION" envDefault:"true"`
	FeatureTwoFactorAuth     bool `env:"FEATURE_TWO_FACTOR_AUTH" envDefault:"false"`
}

// ConfigManager implements the interfaces.ConfigManager interface
type ConfigManager struct {
	envConfig *EnvConfig
	viper     *viper.Viper
}

// NewConfigManager creates a new ConfigManager instance
func NewConfigManager() (*ConfigManager, error) {
	envConfig := &EnvConfig{}

	// Parse environment variables
	if err := env.Parse(envConfig); err != nil {
		return nil, fmt.Errorf("failed to parse environment variables: %w", err)
	}

	// Initialize Viper for additional configuration sources
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/go-better-auth")

	// Set environment variable prefix
	v.SetEnvPrefix("AUTH")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Try to read config file (optional)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	manager := &ConfigManager{
		envConfig: envConfig,
		viper:     v,
	}

	// Validate up front so a weak secret fails at startup rather than silently
	// producing forgeable tokens.
	if err := manager.ValidateConfig(); err != nil {
		return nil, err
	}

	return manager, nil
}

// LoadConfig loads configuration from multiple sources
func (c *ConfigManager) LoadConfig(sources ...interfaces.ConfigSource) error {
	for _, source := range sources {
		data, err := source.Load()
		if err != nil {
			return fmt.Errorf("failed to load config from %s: %w", source.GetName(), err)
		}

		for key, value := range data {
			c.viper.Set(key, value)
		}
	}

	return nil
}

// GetConfig returns the current configuration
func (c *ConfigManager) GetConfig() *interfaces.AuthConfig {
	return &interfaces.AuthConfig{
		JWT: interfaces.JWTConfig{
			Secret:            c.envConfig.JWTSecret,
			Algorithm:         c.envConfig.JWTAlgorithm,
			Issuer:            c.envConfig.JWTIssuer,
			Audience:          c.envConfig.JWTAudience,
			Expiration:        c.envConfig.JWTExpiration,
			RefreshExpiration: c.envConfig.JWTRefreshExpiration,
			ClockSkew:         c.envConfig.JWTClockSkew,
		},
		Session: interfaces.SessionConfig{
			Secret:       c.envConfig.SessionSecret,
			Expiration:   c.envConfig.SessionExpiration,
			CookieName:   c.envConfig.SessionCookieName,
			CookiePath:   "/",
			CookieDomain: "",
			Secure:       c.envConfig.SessionSecure,
			HttpOnly:     c.envConfig.SessionHttpOnly,
			SameSite:     "lax",
			Store:        "memory",
		},
		Password: interfaces.PasswordConfig{
			Algorithm:        c.envConfig.PasswordAlgorithm,
			BcryptCost:       c.envConfig.PasswordBcryptCost,
			MinLength:        c.envConfig.PasswordMinLength,
			MaxLength:        c.envConfig.PasswordMaxLength,
			RequireUppercase: c.envConfig.PasswordRequireUpper,
			RequireLowercase: c.envConfig.PasswordRequireLower,
			RequireNumbers:   c.envConfig.PasswordRequireDigit,
			RequireSymbols:   c.envConfig.PasswordRequireSymbol,
		},
		Database: interfaces.DatabaseConfig{
			Driver:          c.envConfig.DatabaseDriver,
			Host:            c.envConfig.DatabaseHost,
			Port:            c.envConfig.DatabasePort,
			Name:            c.envConfig.DatabaseName,
			Username:        c.envConfig.DatabaseUser,
			Password:        c.envConfig.DatabasePassword,
			SSLMode:         c.envConfig.DatabaseSSLMode,
			MaxOpenConns:    c.envConfig.DatabaseMaxConns,
			ConnMaxIdleTime: c.envConfig.DatabaseIdleTimeout,
		},
		Redis: interfaces.RedisConfig{
			Enabled:     c.envConfig.RedisEnabled,
			Host:        c.envConfig.RedisHost,
			Port:        c.envConfig.RedisPort,
			Password:    c.envConfig.RedisPassword,
			DB:          c.envConfig.RedisDB,
			PoolSize:    c.envConfig.RedisPoolSize,
			DialTimeout: c.envConfig.RedisDialTimeout,
		},
		Email: interfaces.EmailConfig{
			Provider:     c.envConfig.EmailProvider,
			SMTPHost:     c.envConfig.EmailSMTPHost,
			SMTPPort:     c.envConfig.EmailSMTPPort,
			SMTPUsername: c.envConfig.EmailSMTPUser,
			SMTPPassword: c.envConfig.EmailSMTPPass,
			FromEmail:    c.envConfig.EmailFromAddress,
			FromName:     c.envConfig.EmailFromName,
		},
		Security: interfaces.SecurityConfig{
			CSRFProtection: c.envConfig.CSRFEnabled,
			CORSEnabled:    true,
			CORSOrigins:    strings.Split(c.envConfig.CORSAllowedOrigins, ","),
		},
		Server: interfaces.ServerConfig{
			Host:       c.envConfig.ServerHost,
			Port:       c.envConfig.ServerPort,
			TLSEnabled: c.envConfig.ServerTLS,
		},
		Features: interfaces.FeaturesConfig{
			Registration:      c.envConfig.FeatureRegistration,
			PasswordReset:     c.envConfig.FeaturePasswordReset,
			EmailVerification: c.envConfig.FeatureEmailVerification,
			TwoFactorAuth:     c.envConfig.FeatureTwoFactorAuth,
		},
	}
}

// ValidateConfig validates the configuration
func (c *ConfigManager) ValidateConfig() error {
	config := c.GetConfig()

	// Validate JWT secret
	if len(config.JWT.Secret) < 32 {
		return fmt.Errorf("JWT secret must be at least 32 characters long")
	}

	// Validate session secret
	if len(config.Session.Secret) < 32 {
		return fmt.Errorf("session secret must be at least 32 characters long")
	}

	// Validate password requirements
	if config.Password.MinLength < 1 {
		return fmt.Errorf("password minimum length must be at least 1")
	}

	if config.Password.MaxLength < config.Password.MinLength {
		return fmt.Errorf("password maximum length must be greater than minimum length")
	}

	// Validate database configuration
	if config.Database.Driver == "" {
		return fmt.Errorf("database driver is required")
	}

	// Validate server configuration
	if config.Server.Port < 1 || config.Server.Port > 65535 {
		return fmt.Errorf("server port must be between 1 and 65535")
	}

	return nil
}

// ReloadConfig reloads the configuration
func (c *ConfigManager) ReloadConfig() error {
	// Re-parse environment variables
	if err := env.Parse(c.envConfig); err != nil {
		return fmt.Errorf("failed to parse environment variables: %w", err)
	}

	// Re-read config file
	if err := c.viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return nil
}

// WatchConfig watches for configuration changes
func (c *ConfigManager) WatchConfig(callback func(*interfaces.AuthConfig)) error {
	c.viper.WatchConfig()
	c.viper.OnConfigChange(func(e fsnotify.Event) {
		callback(c.GetConfig())
	})

	return nil
}

// GetConfigByKey gets configuration by key path
func (c *ConfigManager) GetConfigByKey(key string) (any, error) {
	return c.viper.Get(key), nil
}

// SetConfigByKey sets configuration by key path
func (c *ConfigManager) SetConfigByKey(key string, value any) error {
	c.viper.Set(key, value)
	return nil
}

// GetSecretFromEnv gets a secret from environment variables
func GetSecretFromEnv(key string) string {
	return os.Getenv(key)
}

// MustGetSecretFromEnv gets a secret from environment variables or panics
func MustGetSecretFromEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("required environment variable %s is not set", key))
	}
	return value
}
