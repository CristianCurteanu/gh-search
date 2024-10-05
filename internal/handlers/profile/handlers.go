package profile

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CristianCurteanu/gh-search/internal/components"
	"github.com/CristianCurteanu/gh-search/internal/handlers/profile/pages"
	"github.com/CristianCurteanu/gh-search/internal/middlewares"
	"github.com/CristianCurteanu/gh-search/pkg/githubapi"
)

type ProfileHandlers struct {
	githubClient *githubapi.GithubApi

	middlewares.UseMiddleware
}

func NewProfileHandlers(githubClient *githubapi.GithubApi) *ProfileHandlers {
	return &ProfileHandlers{githubClient: githubClient}
}

func (ph *ProfileHandlers) GetProfilePage(w http.ResponseWriter, r *http.Request) {
	ph.Handle(w, r, func(w http.ResponseWriter, r *http.Request) {
		t := r.Context().Value(middlewares.CookieAccessTokenKey)
		token := t.(*http.Cookie)

		githubData, _ := ph.githubClient.GetProfileInfo(token.Value)

		pages.ProfilePage(components.ProfileData{
			Username:  githubData.Username,
			Id:        githubData.Username,
			AvatarURL: githubData.AvatarURL,
			Company:   githubData.Company,
			Repos:     fmt.Sprintf("%d", githubData.Repos),
			Gists:     fmt.Sprintf("%d", githubData.Gists),
			Followers: fmt.Sprintf("%d", githubData.Followers),
			Following: fmt.Sprintf("%d", githubData.Following),
		}).Render(r.Context(), w)
	})
}

func (ph *ProfileHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{
		Name:    middlewares.CookieAccessTokenKey,
		Value:   "",
		Path:    "/",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(w, c)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}
