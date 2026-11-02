package adapters

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/sirnacou/file-ctrl/backend/internal/common/config"
	"golang.org/x/oauth2"
)

type OIDCClaims struct {
	Subject           string `json:"sub"`
	Email             string `json:"email"`
	EmailVerified     bool   `json:"email_verified"`
	PreferredUsername string `json:"preferred_username"`
}

type OIDCService struct {
	provider *oidc.Provider
	verifier *oidc.IDTokenVerifier
	oauthCfg oauth2.Config
	enabled  bool
}

func NewOIDCService(ctx context.Context, cfg config.OIDCConfig) (*OIDCService, error) {
	if !cfg.Enabled {
		return &OIDCService{enabled: false}, nil
	}

	provider, err := oidc.NewProvider(ctx, cfg.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("failed to discover OIDC provider: %w", err)
	}

	oidcConfig := &oidc.Config{
		ClientID: cfg.ClientID,
	}
	verifier := provider.Verifier(oidcConfig)

	oauthCfg := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  cfg.RedirectURL,
		Scopes:       cfg.Scopes,
	}

	return &OIDCService{
		provider: provider,
		verifier: verifier,
		oauthCfg: oauthCfg,
		enabled:  true,
	}, nil
}

func (s *OIDCService) IsEnabled() bool {
	return s.enabled
}

func (s *OIDCService) GenerateStateAndNonce() (state, nonce string, err error) {
	bState := make([]byte, 16)
	bNonce := make([]byte, 16)
	if _, err := rand.Read(bState); err != nil {
		return "", "", err
	}
	if _, err := rand.Read(bNonce); err != nil {
		return "", "", err
	}
	return base64.RawURLEncoding.EncodeToString(
			bState,
		), base64.RawURLEncoding.EncodeToString(
			bNonce,
		), nil
}

func (s *OIDCService) AuthCodeURL(state, nonce string) string {
	return s.oauthCfg.AuthCodeURL(state, oidc.Nonce(nonce))
}

func (s *OIDCService) ExchangeAndVerify(
	ctx context.Context,
	code, expectedNonce string,
) (*OIDCClaims, error) {
	token, err := s.oauthCfg.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("missing id_token in token response")
	}

	idToken, err := s.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify id_token: %w", err)
	}

	if idToken.Nonce != expectedNonce {
		return nil, errors.New("nonce does not match")
	}

	var claims OIDCClaims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("failed to parse id_token claims: %w", err)
	}

	return &claims, nil
}
