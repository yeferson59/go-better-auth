package interfaces

import (
	"context"
)

// Mailer defines the interface for email sending operations
type Mailer interface {
	// Send sends an email
	Send(ctx context.Context, email *Email) error

	// SendHTML sends an HTML email
	SendHTML(ctx context.Context, email *EmailHTML) error

	// SendTemplate sends an email using a template
	SendTemplate(ctx context.Context, template string, data any, to []string, subject string) error

	// SendBatch sends multiple emails in batch
	SendBatch(ctx context.Context, emails []*Email) error

	// Verify verifies the email configuration
	Verify(ctx context.Context) error

	// GetProviderInfo returns information about the email provider
	GetProviderInfo() ProviderInfo
}

// Email represents a basic email message
type Email struct {
	From        string            `json:"from"`
	To          []string          `json:"to"`
	CC          []string          `json:"cc,omitempty"`
	BCC         []string          `json:"bcc,omitempty"`
	Subject     string            `json:"subject"`
	Body        string            `json:"body"`
	ContentType string            `json:"content_type"` // "text/plain" or "text/html"
	Headers     map[string]string `json:"headers,omitempty"`
	Attachments []Attachment      `json:"attachments,omitempty"`
	Priority    Priority          `json:"priority,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// EmailHTML represents an HTML email message
type EmailHTML struct {
	From        string            `json:"from"`
	To          []string          `json:"to"`
	CC          []string          `json:"cc,omitempty"`
	BCC         []string          `json:"bcc,omitempty"`
	Subject     string            `json:"subject"`
	HTMLBody    string            `json:"htmlBody"`
	TextBody    string            `json:"textBody,omitempty"` // Fallback for non-HTML clients
	Headers     map[string]string `json:"headers,omitempty"`
	Attachments []Attachment      `json:"attachments,omitempty"`
	Priority    Priority          `json:"priority,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Attachment represents an email attachment
type Attachment struct {
	Name        string `json:"name"`
	ContentType string `json:"contentType"`
	Data        []byte `json:"data"`
	Inline      bool   `json:"inline,omitempty"`
}

// Priority represents email priority
type Priority string

const (
	PriorityLow    Priority = "low"
	PriorityNormal Priority = "normal"
	PriorityHigh   Priority = "high"
	PriorityUrgent Priority = "urgent"
)

// ProviderInfo represents email provider information
type ProviderInfo struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Description string   `json:"description"`
	Features    []string `json:"features"`
}

// EmailTemplate defines the interface for email templates
type EmailTemplate interface {
	// Render renders the template with the given data
	Render(ctx context.Context, data any) (*RenderedEmail, error)

	// GetName returns the template name
	GetName() string

	// GetSubject returns the template subject
	GetSubject() string

	// Validate validates the template
	Validate() error
}

// RenderedEmail represents a rendered email template
type RenderedEmail struct {
	Subject  string `json:"subject"`
	HTMLBody string `json:"htmlBody"`
	TextBody string `json:"textBody"`
}

// EmailTemplateEngine defines the interface for email template engines
type EmailTemplateEngine interface {
	// LoadTemplate loads a template from a file or string
	LoadTemplate(name, content string) (EmailTemplate, error)

	// LoadTemplateFromFile loads a template from a file
	LoadTemplateFromFile(name, filePath string) (EmailTemplate, error)

	// GetTemplate returns a loaded template
	GetTemplate(name string) (EmailTemplate, error)

	// ListTemplates returns all loaded templates
	ListTemplates() []string

	// DeleteTemplate removes a template
	DeleteTemplate(name string) error

	// RenderTemplate renders a template with data
	RenderTemplate(name string, data any) (*RenderedEmail, error)
}

// AuthEmailTemplates defines common authentication email templates
type AuthEmailTemplates struct {
	Welcome           string
	EmailVerification string
	PasswordReset     string
	PasswordChanged   string
	AccountLocked     string
	AccountUnlocked   string
	LoginNotification string
	TwoFactorCode     string
	InviteUser        string
}

// DefaultAuthEmailTemplates provides default template names
var DefaultAuthEmailTemplates = AuthEmailTemplates{
	Welcome:           "welcome",
	EmailVerification: "email_verification",
	PasswordReset:     "password_reset",
	PasswordChanged:   "password_changed",
	AccountLocked:     "account_locked",
	AccountUnlocked:   "account_unlocked",
	LoginNotification: "login_notification",
	TwoFactorCode:     "two_factor_code",
	InviteUser:        "invite_user",
}

// EmailQueue defines the interface for email queue operations
type EmailQueue interface {
	// Enqueue adds an email to the queue
	Enqueue(ctx context.Context, email *Email) error

	// Dequeue removes and returns an email from the queue
	Dequeue(ctx context.Context) (*Email, error)

	// Peek returns the next email without removing it
	Peek(ctx context.Context) (*Email, error)

	// Size returns the number of emails in the queue
	Size(ctx context.Context) (int, error)

	// Clear removes all emails from the queue
	Clear(ctx context.Context) error

	// GetStats returns queue statistics
	GetStats(ctx context.Context) (*QueueStats, error)
}

// QueueStats represents email queue statistics
type QueueStats struct {
	TotalEmails     int `json:"totalEmails"`
	PendingEmails   int `json:"pendingEmails"`
	ProcessedEmails int `json:"processedEmails"`
	FailedEmails    int `json:"failedEmails"`
	RetryEmails     int `json:"retryEmails"`
}

// EmailValidator defines the interface for email validation
type EmailValidator interface {
	// Validate validates an email address
	Validate(ctx context.Context, email string) error

	// IsValid checks if an email address is valid
	IsValid(ctx context.Context, email string) bool

	// IsDisposable checks if an email address is from a disposable provider
	IsDisposable(ctx context.Context, email string) bool

	// GetDomain returns the domain part of an email address
	GetDomain(email string) string

	// Normalize normalizes an email address
	Normalize(email string) string
}

// EmailProvider defines the interface for email providers
type EmailProvider interface {
	// Send sends an email through the provider
	Send(ctx context.Context, email *Email) error

	// SendHTML sends an HTML email through the provider
	SendHTML(ctx context.Context, email *EmailHTML) error

	// GetProviderInfo returns provider information
	GetProviderInfo() ProviderInfo

	// Verify verifies the provider configuration
	Verify(ctx context.Context) error

	// GetStats returns provider statistics
	GetStats(ctx context.Context) (*ProviderStats, error)
}

// ProviderStats represents email provider statistics
type ProviderStats struct {
	SentEmails      int `json:"sentEmails"`
	FailedEmails    int `json:"failedEmails"`
	DeliveredEmails int `json:"deliveredEmails"`
	BouncedEmails   int `json:"bouncedEmails"`
}
