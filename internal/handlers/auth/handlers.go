package auth

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CristianCurteanu/gh-search/internal/components"
	"github.com/CristianCurteanu/gh-search/internal/middlewares"
	"github.com/CristianCurteanu/gh-search/pkg/githubapi"
)

type AuthHandlersConfig struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

type AuthHandlers struct {
	conf         AuthHandlersConfig
	githubClient *githubapi.GithubApi

	middlewares.UseMiddleware
}

func NewAuthHandlers(conf AuthHandlersConfig, githubClient *githubapi.GithubApi) AuthHandlers {
	return AuthHandlers{conf: conf}
}

func (ah AuthHandlers) RootHandler(w http.ResponseWriter, r *http.Request) {
	loginBtn := components.Login()

	components.LoginPage(loginBtn).Render(r.Context(), w)
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

		ghAuth, _ := ah.githubClient.GetGithubAccessToken(code)

		http.SetCookie(w, &http.Cookie{
			Name:    "access_token",
			Value:   ghAuth.AccessToken,
			Expires: time.Now().Add(time.Hour * 24),
			Path:    "/",
		})

		http.Redirect(w, r, "/profile", http.StatusMovedPermanently)
	}
}
