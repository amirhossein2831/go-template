package client

type CreateClientRequest struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	ClientID string `json:"client_id" validate:"required,uuid"`

	// Authentication
	GrantTypes    []string `json:"grant_types" validate:"required,min=1,dive,oneof=authorization_code client_credentials refresh_token implicit password device_code"` // At least one, valid types
	ResponseTypes []string `json:"response_types" validate:"required,min=1,dive,oneof=code token id_token"`                                                            // At least one, valid types

	// URIs
	RedirectURIs           []string `json:"redirect_uris" validate:"required,min=1,dive,url"` // At least one, must be valid URLs
	PostLogoutRedirectURIs []string `json:"post_logout_redirect_uris" validate:"dive,url"`    // Optional, must be valid URLs if present
	AllowedOrigins         []string `json:"allowed_origins" validate:"dive,uri"`              // Optional, must be valid URIs if present

	// Access configuration
	AccessTokenLifetime  int      `json:"access_token_lifetime" validate:"min=60,max=36000"`       // Optional, min 60s, max 10 hours
	RefreshTokenLifetime int      `json:"refresh_token_lifetime" validate:"min=3600,max=31536000"` // Optional, min 1 hour, max 1 year
	Scopes               []string `json:"scopes" validate:"required,min=1,dive"`                   // At least one, must be alphanumeric (example scope validation)

	// Security
	RequirePKCE        bool `json:"require_pkce"`          // Defaults to false
	AllowPlainTextPKCE bool `json:"allow_plain_text_pkce"` // Defaults to false

	// Client Type (public or confidential)
	Type string `json:"type" validate:"required,oneof=public confidential"` // Required, must be "public" or "confidential"

	// Additional claims to include in tokens - dynamic, so no specific validation other than map
	DefaultClaims map[string]interface{} `json:"default_claims"`
}
