package middlewares

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/CristianCurteanu/gh-search/internal/auth"
)

type CookiesSession struct {
	sessionStorage auth.SessionStorage
	jwtSigner      *auth.JWTAuth
}

const (
	CookieAccessTokenKey = "access_token"
)

func NewCookieSessionHandler(sessionStorage auth.SessionStorage, jwtSigner *auth.JWTAuth) Middleware {
	return &CookiesSession{sessionStorage, jwtSigner}
}

func (sm *CookiesSession) Execute(w http.ResponseWriter, r *http.Request) error {
	token, err := r.Cookie("access_token")
	if err != nil {
		return err
	}

	jwtToken, err := sm.jwtSigner.VerifyToken(token.Value)
	if err != nil {
		log.Printf("failed to get jwt token from cookie, err=%q", err)

		return nil
	}

	sessionId, err := jwtToken.Claims.GetSubject()
	if err != nil {
		log.Printf("failed to get subject from claims, err=%q", err)

		return nil
	}

	session, err := sm.sessionStorage.GetSession(r.Context(), sessionId)
	if err != nil {
		log.Printf("failed to get session, err=%q", err)
		return nil
	}

	ctx := context.WithValue(r.Context(), CookieAccessTokenKey, session.Secret)
	*r = *r.WithContext(ctx)

	return nil
}

func (sm *CookiesSession) GetFallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:    CookieAccessTokenKey,
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
		})
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	}
}
