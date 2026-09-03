package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

type OIDC struct {
	ProviderName string
	Provider     *oidc.Provider
	OAuth2Config *oauth2.Config
	Verifier     *oidc.IDTokenVerifier
}

type Claims struct {
	Subject string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
}

func NewOIDC(ctx context.Context, cfg Config) (*OIDC, error) {
	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, err
	}

	oauthConfig := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  cfg.RedirectURL,
		Scopes: []string{
			"openid",
			"email",
			"profile",
		},
	}

	verifier := provider.Verifier(
		&oidc.Config{
			ClientID: cfg.ClientID,
		},
	)

	return &OIDC{
		ProviderName: cfg.ProviderName,
		Provider:     provider,
		OAuth2Config: oauthConfig,
		Verifier:     verifier,
	}, nil
}

func (o *OIDC) AuthorizationURL(state string) string {
	return o.OAuth2Config.AuthCodeURL(
		state,
		oauth2.AccessTypeOffline,
	)
}

func (o *OIDC) ExchangeCode(ctx context.Context, code string) (*Claims, error) {
	token, err := o.OAuth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, ErrMissingIDToken
	}

	idToken, err := o.Verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	var claims Claims
	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

func GenerateState() (string, error) {
	bytes := make([]byte, 32)

	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(bytes), nil
}

var ErrMissingIDToken = errors.New("id_token missing from token response")
