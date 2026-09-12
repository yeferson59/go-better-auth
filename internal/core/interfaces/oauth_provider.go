package interfaces

import (
	"context"
	"net/url"

	"github.com/yeferson59/go-better-auth/internal/core/entities"
)

// OAuthProvider defines the interface for OAuth provider operations
type OAuthProvider interface {
	// GetAuthURL returns the OAuth authorization URL
	GetAuthURL(ctx context.Context, state string, scopes []string) (string, error)

	// ExchangeCode exchanges an authorization code for tokens
	ExchangeCode(ctx context.Context, code string) (*OAuthTokens, error)

	// RefreshToken refreshes an access token using a refresh token
	RefreshToken(ctx context.Context, refreshToken string) (*OAuthTokens, error)

	// GetUserInfo retrieves user information from the provider
	GetUserInfo(ctx context.Context, accessToken string) (*entities.OAuthUserInfo, error)

	// RevokeToken revokes an access token
	RevokeToken(ctx context.Context, token string) error

	// ValidateToken validates an access token
	ValidateToken(ctx context.Context, token string) (*TokenValidationResult, error)

	// GetProviderName returns the provider name
	GetProviderName() string

	// GetConfig returns the provider configuration
	GetConfig() *OAuthConfig

	// IsEnabled returns whether the provider is enabled
	IsEnabled() bool
}

// OAuthTokens represents OAuth tokens
type OAuthTokens struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken,omitempty"`
	TokenType    string `json:"tokenType"`
	ExpiresIn    int    `json:"expiresIn"`
	Scope        string `json:"scope,omitempty"`
	IDToken      string `json:"idToken,omitempty"`
}

// TokenValidationResult represents token validation result
type TokenValidationResult struct {
	Valid     bool     `json:"valid"`
	UserID    string   `json:"userId,omitempty"`
	Username  string   `json:"username,omitempty"`
	Email     string   `json:"email,omitempty"`
	ExpiresAt int64    `json:"expiresAt,omitempty"`
	Scopes    []string `json:"scopes,omitempty"`
}

// OAuthConfig represents OAuth provider configuration
type OAuthConfig struct {
	ClientID     string            `json:"clientId"`
	ClientSecret string            `json:"clientSecret"`
	RedirectURL  string            `json:"redirectUrl"`
	Scopes       []string          `json:"scopes"`
	AuthURL      string            `json:"authUrl"`
	TokenURL     string            `json:"tokenUrl"`
	UserInfoURL  string            `json:"userInfoUrl"`
	RevokeURL    string            `json:"revokeUrl,omitempty"`
	Endpoints    map[string]string `json:"endpoints,omitempty"`
	Params       map[string]string `json:"params,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
}

// OAuthProviderFactory defines the interface for creating OAuth providers
type OAuthProviderFactory interface {
	// Create creates an OAuth provider instance
	Create(providerName string, config *OAuthConfig) (OAuthProvider, error)

	// GetSupportedProviders returns a list of supported providers
	GetSupportedProviders() []string

	// IsSupported checks if a provider is supported
	IsSupported(providerName string) bool

	// GetDefaultConfig returns the default configuration for a provider
	GetDefaultConfig(providerName string) (*OAuthConfig, error)
}

// OAuthStateManager defines the interface for OAuth state management
type OAuthStateManager interface {
	// Generate generates a new OAuth state
	Generate(ctx context.Context, provider string, returnURL string) (*entities.OAuthState, error)

	// Validate validates an OAuth state
	Validate(ctx context.Context, state string) (*entities.OAuthState, error)

	// Delete removes an OAuth state
	Delete(ctx context.Context, state string) error

	// Cleanup removes expired OAuth states
	Cleanup(ctx context.Context) error
}

// OAuthService defines the interface for OAuth service operations
type OAuthService interface {
	// GetAuthURL returns the OAuth authorization URL
	GetAuthURL(ctx context.Context, provider string, returnURL string) (string, error)

	// HandleCallback handles the OAuth callback
	HandleCallback(ctx context.Context, provider string, code string, state string) (*OAuthCallbackResult, error)

	// LinkAccount links an OAuth account to a user
	LinkAccount(ctx context.Context, userID string, provider string, code string) error

	// UnlinkAccount unlinks an OAuth account from a user
	UnlinkAccount(ctx context.Context, userID string, provider string) error

	// GetLinkedAccounts returns all linked accounts for a user
	GetLinkedAccounts(ctx context.Context, userID string) ([]*entities.OAuthLink, error)

	// RefreshToken refreshes an OAuth token
	RefreshToken(ctx context.Context, userID string, provider string) (*OAuthTokens, error)

	// GetEnabledProviders returns all enabled OAuth providers
	GetEnabledProviders(ctx context.Context) ([]*entities.OAuthProvider, error)
}

// OAuthCallbackResult represents the result of OAuth callback processing
type OAuthCallbackResult struct {
	User         *entities.User          `json:"user"`
	IsNewUser    bool                    `json:"isNewUser"`
	AccessToken  string                  `json:"accessToken"`
	RefreshToken string                  `json:"refreshToken,omitempty"`
	UserInfo     *entities.OAuthUserInfo `json:"userInfo"`
	Provider     string                  `json:"provider"`
	ReturnURL    string                  `json:"returnUrl,omitempty"`
}

// OAuthUserInfoMapper defines the interface for mapping OAuth user info
type OAuthUserInfoMapper interface {
	// MapUserInfo maps provider-specific user info to standard format
	MapUserInfo(ctx context.Context, provider string, rawUserInfo map[string]any) (*entities.OAuthUserInfo, error)

	// GetEmailFromUserInfo extracts email from user info
	GetEmailFromUserInfo(ctx context.Context, provider string, userInfo map[string]any) (string, error)

	// GetUsernameFromUserInfo extracts username from user info
	GetUsernameFromUserInfo(ctx context.Context, provider string, userInfo map[string]any) (string, error)

	// GetNameFromUserInfo extracts name from user info
	GetNameFromUserInfo(ctx context.Context, provider string, userInfo map[string]any) (string, error)
}

// OAuthProviderRegistry defines the interface for OAuth provider registry
type OAuthProviderRegistry interface {
	// Register registers an OAuth provider
	Register(provider OAuthProvider) error

	// Get returns an OAuth provider by name
	Get(providerName string) (OAuthProvider, error)

	// List returns all registered OAuth providers
	List() []OAuthProvider

	// GetEnabled returns all enabled OAuth providers
	GetEnabled() []OAuthProvider

	// Unregister removes an OAuth provider
	Unregister(providerName string) error

	// IsRegistered checks if a provider is registered
	IsRegistered(providerName string) bool
}

// OAuthTokenStore defines the interface for OAuth token storage
type OAuthTokenStore interface {
	// Store stores OAuth tokens
	Store(ctx context.Context, userID string, provider string, tokens *OAuthTokens) error

	// Get retrieves OAuth tokens
	Get(ctx context.Context, userID string, provider string) (*OAuthTokens, error)

	// Update updates OAuth tokens
	Update(ctx context.Context, userID string, provider string, tokens *OAuthTokens) error

	// Delete removes OAuth tokens
	Delete(ctx context.Context, userID string, provider string) error

	// List returns all OAuth tokens for a user
	List(ctx context.Context, userID string) (map[string]*OAuthTokens, error)

	// IsExpired checks if tokens are expired
	IsExpired(ctx context.Context, userID string, provider string) (bool, error)

	// Cleanup removes expired tokens
	Cleanup(ctx context.Context) error
}

// OAuthFlowHandler defines the interface for OAuth flow handling
type OAuthFlowHandler interface {
	// HandleAuthorizationRequest handles the initial authorization request
	HandleAuthorizationRequest(ctx context.Context, provider string, params url.Values) (*AuthorizationResponse, error)

	// HandleTokenRequest handles the token exchange request
	HandleTokenRequest(ctx context.Context, provider string, code string, state string) (*TokenResponse, error)

	// HandleUserInfoRequest handles the user info request
	HandleUserInfoRequest(ctx context.Context, provider string, accessToken string) (*UserInfoResponse, error)
}

// AuthorizationResponse represents the authorization response
type AuthorizationResponse struct {
	AuthURL   string `json:"auth_url"`
	State     string `json:"state"`
	Provider  string `json:"provider"`
	ReturnURL string `json:"returnUrl,omitempty"`
}

// TokenResponse represents the token response
type TokenResponse struct {
	Tokens   *OAuthTokens `json:"tokens"`
	Provider string       `json:"provider"`
	State    string       `json:"state"`
}

// UserInfoResponse represents the user info response
type UserInfoResponse struct {
	UserInfo *entities.OAuthUserInfo `json:"userInfo"`
	Provider string                  `json:"provider"`
	RawData  map[string]any          `json:"rawData,omitempty"`
}

// OAuthMiddleware defines the interface for OAuth middleware
type OAuthMiddleware interface {
	// RequireOAuth middleware that requires OAuth authentication
	RequireOAuth(provider string) func(any) any

	// OptionalOAuth middleware that optionally checks OAuth authentication
	OptionalOAuth(provider string) func(any) any

	// RequireScopes middleware that requires specific OAuth scopes
	RequireScopes(scopes []string) func(any) any

	// ExtractOAuthInfo extracts OAuth information from request
	ExtractOAuthInfo(ctx context.Context) (*OAuthInfo, error)
}

// OAuthInfo represents OAuth information from middleware
type OAuthInfo struct {
	Provider    string   `json:"provider"`
	UserID      string   `json:"userId"`
	AccessToken string   `json:"access_token"`
	Scopes      []string `json:"scopes"`
	ExpiresAt   int64    `json:"expiresAt"`
}

// OAuthEventHandler defines the interface for OAuth event handling
type OAuthEventHandler interface {
	// OnAuthorizationStart handles authorization start event
	OnAuthorizationStart(ctx context.Context, provider string, userID string) error

	// OnAuthorizationSuccess handles authorization success event
	OnAuthorizationSuccess(ctx context.Context, provider string, userID string, tokens *OAuthTokens) error

	// OnAuthorizationError handles authorization error event
	OnAuthorizationError(ctx context.Context, provider string, userID string, err error) error

	// OnTokenRefresh handles token refresh event
	OnTokenRefresh(ctx context.Context, provider string, userID string, tokens *OAuthTokens) error

	// OnAccountLink handles account link event
	OnAccountLink(ctx context.Context, userID string, provider string, oauthUserInfo *entities.OAuthUserInfo) error

	// OnAccountUnlink handles account unlink event
	OnAccountUnlink(ctx context.Context, userID string, provider string) error
}
