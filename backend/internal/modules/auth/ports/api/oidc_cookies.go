package api

import (
	"net/http"
	"time"
)

const (
	OIDCStateCookieName = "oidc_state"
	OIDCNonceCookieName = "oidc_nonce"
	OIDCCookieTTL       = 10 * time.Minute
)

// HumaOIDCCookies can be embedded in Huma Input structs to automatically bind OIDC flow cookies.
type HumaOIDCCookies struct {
	State string `cookie:"oidc_state" doc:"CSRF state verification token"`
	Nonce string `cookie:"oidc_nonce" doc:"Replay protection nonce"`
}

type OIDCCookies struct {
	secure bool
}

func NewOIDCCookies(secure bool) *OIDCCookies {
	return &OIDCCookies{secure: secure}
}

// NewState creates a short-lived cookie holding the OIDC state parameter.
func (o *OIDCCookies) NewState(state string) http.Cookie {
	expiresAt := time.Now().Add(OIDCCookieTTL)
	return http.Cookie{
		Name:     OIDCStateCookieName,
		Value:    state,
		Path:     "/api/auth/oidc/callback",
		Expires:  expiresAt,
		MaxAge:   int(OIDCCookieTTL.Seconds()),
		Secure:   o.secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

// NewNonce creates a short-lived cookie holding the OIDC nonce parameter.
func (o *OIDCCookies) NewNonce(nonce string) http.Cookie {
	expiresAt := time.Now().Add(OIDCCookieTTL)
	return http.Cookie{
		Name:     OIDCNonceCookieName,
		Value:    nonce,
		Path:     "/api/auth/oidc/callback",
		Expires:  expiresAt,
		MaxAge:   int(OIDCCookieTTL.Seconds()),
		Secure:   o.secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

// ClearState returns a cookie configured to expire and erase the state cookie.
func (o *OIDCCookies) ClearState() http.Cookie {
	return http.Cookie{
		Name:     OIDCStateCookieName,
		Value:    "",
		Path:     "/api/auth/oidc/callback",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   o.secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

// ClearNonce returns a cookie configured to expire and erase the nonce cookie.
func (o *OIDCCookies) ClearNonce() http.Cookie {
	return http.Cookie{
		Name:     OIDCNonceCookieName,
		Value:    "",
		Path:     "/api/auth/oidc/callback",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   o.secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}
