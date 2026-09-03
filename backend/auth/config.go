package auth

import "errors"

// Config contains the configuration required for
// OpenID Connect authentication.
type Config struct {
	ProviderName string
	IssuerURL    string
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// Validate checks whether the required OIDC configuration exists.
func (c Config) Validate() error {
	if c.ProviderName == "" {
		return errors.New("OIDC provider name is required")
	}

	if c.IssuerURL == "" {
		return errors.New("OIDC issuer URL is required")
	}

	if c.ClientID == "" {
		return errors.New("OIDC client ID is required")
	}

	if c.ClientSecret == "" {
		return errors.New("OIDC client secret is required")
	}

	if c.RedirectURL == "" {
		return errors.New("OIDC redirect URL is required")
	}

	return nil
}
