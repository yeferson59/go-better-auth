package entities

import (
	"maps"
	"time"
)

type Provider string

const (
	Google    Provider = "google"
	GitHub    Provider = "github"
	Facebook  Provider = "facebook"
	Discord   Provider = "discord"
	Twitter   Provider = "twitter"
	Microsoft Provider = "microsoft"
	LinkedIn  Provider = "linkedin"
)

// DefaultProviderConfigs provides default configurations for supported providers
var DefaultProviderConfigs = map[Provider]OAuthProvider{
	Google: {
		Name:        Google,
		DisplayName: "Google",
		AuthURL:     "https://accounts.google.com/o/oauth2/v2/auth",
		TokenURL:    "https://oauth2.googleapis.com/token",
		UserInfoURL: "https://www.googleapis.com/oauth2/v2/userinfo",
		Scopes:      []string{"openid", "email", "profile"},
	},
	GitHub: {
		Name:        GitHub,
		DisplayName: "GitHub",
		AuthURL:     "https://github.com/login/oauth/authorize",
		TokenURL:    "https://github.com/login/oauth/access_token",
		UserInfoURL: "https://api.github.com/user",
		Scopes:      []string{"user:email"},
	},
	Facebook: {
		Name:        Facebook,
		DisplayName: "Facebook",
		AuthURL:     "https://www.facebook.com/v18.0/dialog/oauth",
		TokenURL:    "https://graph.facebook.com/v18.0/oauth/access_token",
		UserInfoURL: "https://graph.facebook.com/me",
		Scopes:      []string{"email", "public_profile"},
	},
	Discord: {
		Name:        Discord,
		DisplayName: "Discord",
		AuthURL:     "https://discord.com/api/oauth2/authorize",
		TokenURL:    "https://discord.com/api/oauth2/token",
		UserInfoURL: "https://discord.com/api/users/@me",
		Scopes:      []string{"identify", "email"},
	},
	Microsoft: {
		Name:        Microsoft,
		DisplayName: "Microsoft",
		AuthURL:     "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
		TokenURL:    "https://login.microsoftonline.com/common/oauth2/v2.0/token",
		UserInfoURL: "https://graph.microsoft.com/v1.0/me",
		Scopes:      []string{"openid", "email", "profile"},
	},
}

// OAuthProvider represents an OAuth provider configuration
type OAuthProvider struct {
	ID           string            `json:"id"`
	Name         Provider          `json:"name"`        // e.g., "google", "github", "facebook"
	DisplayName  string            `json:"displayName"` // e.g., "Google", "GitHub", "Facebook"
	ClientID     string            `json:"clientId"`
	ClientSecret string            `json:"-"`
	RedirectURL  string            `json:"redirectUrl"`
	Scopes       []string          `json:"scopes"`
	AuthURL      string            `json:"authUrl"`
	TokenURL     string            `json:"tokenUrl"`
	UserInfoURL  string            `json:"userInfoUrl"`
	IsEnabled    bool              `json:"isEnabled"`
	Config       map[string]string `json:"config,omitempty"` // Additional provider-specific config
	CreatedAt    time.Time         `json:"createdAt"`
	UpdatedAt    time.Time         `json:"updatedAt"`
}

// IsSupported checks if the provider is supported
func (p *OAuthProvider) IsSupported() bool {
	_, exists := DefaultProviderConfigs[p.Name]
	return exists
}

// GetAuthURL returns the authorization URL for the provider
func (p *OAuthProvider) GetAuthURL() string {
	return p.AuthURL
}

// GetTokenURL returns the token URL for the provider
func (p *OAuthProvider) GetTokenURL() string {
	return p.TokenURL
}

// GetUserInfoURL returns the user info URL for the provider
func (p *OAuthProvider) GetUserInfoURL() string {
	return p.UserInfoURL
}

// GetScopes returns the scopes for the provider
func (p *OAuthProvider) GetScopes() []string {
	return p.Scopes
}

// Enable enables the OAuth provider
func (p *OAuthProvider) Enable() {
	p.IsEnabled = true
	p.UpdatedAt = time.Now()
}

// Disable disables the OAuth provider
func (p *OAuthProvider) Disable() {
	p.IsEnabled = false
	p.UpdatedAt = time.Now()
}

// UpdateConfig updates the provider configuration
func (p *OAuthProvider) UpdateConfig(config map[string]string) {
	if p.Config == nil {
		p.Config = make(map[string]string)
	}

	maps.Copy(p.Config, config)

	p.UpdatedAt = time.Now()
}

// GetConfig returns a configuration value
func (p *OAuthProvider) GetConfig(key string) (string, bool) {
	if p.Config == nil {
		return "", false
	}
	value, exists := p.Config[key]
	return value, exists
}

// OAuthUserInfo represents user information from OAuth provider
type OAuthUserInfo struct {
	ID        string         `json:"id"`
	Email     string         `json:"email"`
	Name      string         `json:"name"`
	FirstName string         `json:"firstName,omitempty"`
	LastName  string         `json:"lastName,omitempty"`
	Username  string         `json:"username,omitempty"`
	Avatar    string         `json:"avatar,omitempty"`
	Verified  bool           `json:"verified"`
	Locale    string         `json:"locale,omitempty"`
	Raw       map[string]any `json:"raw,omitempty"` // Raw response from provider
}

// OAuthState represents OAuth state for security
type OAuthState struct {
	State     string    `json:"state"`
	Provider  string    `json:"provider"`
	ReturnURL string    `json:"returnUrl,omitempty"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

// IsExpired checks if the OAuth state has expired
func (s *OAuthState) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// IsValid checks if the OAuth state is valid
func (s *OAuthState) IsValid() bool {
	return !s.IsExpired() && s.State != "" && s.Provider != ""
}
