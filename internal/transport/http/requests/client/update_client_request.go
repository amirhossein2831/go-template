package client

type UpdateClientRequest struct {
	Name *string `json:"name" validate:"omitempty,min=3,max=100"` // Optional, but validate if present

	GrantTypes    []string `json:"grant_types" validate:"omitempty,min=1,dive,oneof=authorization_code client_credentials refresh_token implicit password device_code"`
	ResponseTypes []string `json:"response_types" validate:"omitempty,min=1,dive,oneof=code token id_token"`

	RedirectURIs           []string `json:"redirect_uris" validate:"omitempty,min=1,dive,url"`
	PostLogoutRedirectURIs []string `json:"post_logout_redirect_uris" validate:"omitempty,dive,url"`
	AllowedOrigins         []string `json:"allowed_origins" validate:"omitempty,dive,uri"`

	AccessTokenLifetime  *int     `json:"access_token_lifetime" validate:"omitempty,min=60,max=36000"`
	RefreshTokenLifetime *int     `json:"refresh_token_lifetime" validate:"omitempty,min=3600,max=31536000"`
	Scopes               []string `json:"scopes" validate:"omitempty,min=1,dive"`

	RequirePKCE        *bool `json:"require_pkce"`
	AllowPlainTextPKCE *bool `json:"allow_plain_text_pkce"`

	Type *string `json:"type" validate:"omitempty,oneof=public confidential"`

	DefaultClaims map[string]interface{} `json:"default_claims"`
}
