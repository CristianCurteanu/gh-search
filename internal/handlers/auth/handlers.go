package auth

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/CristianCurteanu/gh-search/internal/auth"
	"github.com/CristianCurteanu/gh-search/internal/handlers/auth/pages"
	"github.com/CristianCurteanu/gh-search/internal/middlewares"
	"github.com/CristianCurteanu/gh-search/pkg/githubapi"
	"github.com/google/uuid"
)

type AuthHandlersConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

type AuthHandlers struct {
	conf           AuthHandlersConfig
	githubClient   *githubapi.GithubApi
	jwtSigner      *auth.JWTAuth
	sessionStorage auth.SessionStorage

	middlewares.UseMiddleware
}

func NewAuthHandlers(conf AuthHandlersConfig, githubClient *githubapi.GithubApi, sessionStorage auth.SessionStorage, jwtSigner *auth.JWTAuth) AuthHandlers {
	return AuthHandlers{
		conf:           conf,
		githubClient:   githubClient,
		sessionStorage: sessionStorage,
		jwtSigner:      jwtSigner,
	}
}

func (ah AuthHandlers) RootHandler(w http.ResponseWriter, r *http.Request) {
	loginBtn := pages.Login()

	pages.LoginPage(loginBtn).Render(r.Context(), w)
}

func (ah AuthHandlers) GithubLoginHandler(w http.ResponseWriter, r *http.Request) {
	ah.Handle(w, r, func(w http.ResponseWriter, r *http.Request) {
		redirectURL := fmt.Sprintf(
			"https://github.com/login/oauth/authorize?client_id=%s&redirect_uri=%s&scope=%s",
			ah.conf.ClientId,
			ah.conf.RedirectUrl,
			"repo",
		)

		http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
	})
}

func (ah AuthHandlers) GithubCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")

		fmt.Printf("ah.githubClient: %v\n", ah.githubClient)
		ghAuth, _ := ah.githubClient.GetGithubAccessToken(code)

		sessionId := uuid.NewString()
		sessionExpiresAt := time.Now().Add(time.Hour * 24).UTC()

		token, err := ah.jwtSigner.CreateToken(sessionId, &sessionExpiresAt)
		if err != nil {
			log.Printf("failed to create JWT token, err=%q", err)
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}

		err = ah.sessionStorage.StoreSession(r.Context(), sessionId, auth.Session{
			Id:        sessionId,
			Secret:    ghAuth.AccessToken,
			ExpiresAt: &sessionExpiresAt,
		})
		if err != nil {
			log.Printf("failed to store session, err=%q", err)
			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "access_token",
			Value:   token,
			Expires: sessionExpiresAt,
			Path:    "/",
		})

		http.Redirect(w, r, "/profile", http.StatusMovedPermanently)
	}
}
