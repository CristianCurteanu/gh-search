package profile

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CristianCurteanu/gh-search/internal/handlers/profile/pages"
	"github.com/CristianCurteanu/gh-search/internal/handlers/profile/utils"
	"github.com/CristianCurteanu/gh-search/internal/layouts"
	"github.com/CristianCurteanu/gh-search/internal/middlewares"
	"github.com/CristianCurteanu/gh-search/pkg/cache"
	"github.com/CristianCurteanu/gh-search/pkg/githubapi"
)

type ProfileHandlers struct {
	githubClient *githubapi.GithubApi
	cache        cache.Cache

	middlewares.UseMiddleware
}

func NewProfileHandlers(githubClient *githubapi.GithubApi, cache cache.Cache) *ProfileHandlers {
	return &ProfileHandlers{
		githubClient: githubClient,
		cache:        cache,
	}
}

func (ph *ProfileHandlers) GetProfilePage(w http.ResponseWriter, r *http.Request) {
	ph.Handle(w, r, func(w http.ResponseWriter, r *http.Request) {
		t := r.Context().Value(middlewares.CookieAccessTokenKey)
		token := t.(string)

		githubData, _ := utils.GetProfileData(r.Context(), token, ph.githubClient, ph.cache)

		pages.ProfilePage(mapProfileData(githubData)).Render(r.Context(), w)
	})
}

func (ph *ProfileHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	ph.Handle(w, r, func(w http.ResponseWriter, r *http.Request) {
		t := r.Context().Value(middlewares.CookieAccessTokenKey)
		token := t.(string)

		http.SetCookie(w, &http.Cookie{
			Name:    middlewares.CookieAccessTokenKey,
			Value:   "",
			Path:    "/",
			Expires: time.Unix(0, 0),
		})
		ph.cache.Delete(r.Context(), token)
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
	})

}

func mapProfileData(githubData githubapi.ProfileData) layouts.ProfileData {
	return layouts.ProfileData{
		Username:  githubData.Username,
		Id:        githubData.Username,
		AvatarURL: githubData.AvatarURL,
		Company:   githubData.Company,
		Repos:     fmt.Sprintf("%d", githubData.Repos),
		Gists:     fmt.Sprintf("%d", githubData.Gists),
		Followers: fmt.Sprintf("%d", githubData.Followers),
		Following: fmt.Sprintf("%d", githubData.Following),
	}
}
