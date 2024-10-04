package middlewares

import (
	"context"
	"net/http"
)

type CookiesSession struct {
	fallback http.HandlerFunc
}

const (
	CookieAccessTokenKey = "access_token"
)

func NewCookieSessionHandler(fallback http.HandlerFunc) Middleware {
	return &CookiesSession{fallback}
}

func (sm *CookiesSession) Execute(w http.ResponseWriter, r *http.Request) error {
	token, err := r.Cookie("access_token")
	if err != nil {
		return err
	}

	ctx := context.WithValue(r.Context(), CookieAccessTokenKey, token)
	*r = *r.WithContext(ctx)

	return nil
}

func (sm *CookiesSession) GetFallback() http.HandlerFunc {
	return sm.fallback
}
