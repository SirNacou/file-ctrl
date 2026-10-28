package api

import (
	"net/http"
	"time"
)

const SessionCookieName = "filectrl_session"

func NewSessionCookieName(rawToken string, expiresAt time.Time, isProd bool) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    rawToken,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
		Secure:   isProd,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func NewClearSessionCookie(isProd bool) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   isProd,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}
