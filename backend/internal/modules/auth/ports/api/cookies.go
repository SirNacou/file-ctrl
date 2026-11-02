package api

import (
	"net/http"
	"time"
)

type HumaSessionCookie struct {
	SessionCookie http.Cookie `cookie:"filectrl_session"`
}

const SessionCookieName = "filectrl_session"

type SessionCookies struct {
	secure bool
}

func NewSessionCookies(secure bool) *SessionCookies {
	return &SessionCookies{secure: secure}
}

func (s *SessionCookies) New(rawToken string, expiresAt time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Value:    rawToken,
		Path:     "/",
		Expires:  expiresAt,
		MaxAge:   int(time.Until(expiresAt).Seconds()),
		Secure:   s.secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}

func (s *SessionCookies) Clear() *http.Cookie {
	return &http.Cookie{
		Name:     SessionCookieName,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   s.secure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
}
