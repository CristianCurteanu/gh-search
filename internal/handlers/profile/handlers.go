package profile

import (
	"fmt"
	"log"
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

		githubData, err := utils.GetProfileData(r.Context(), token, ph.githubClient, ph.cache)
		if err != nil {
			log.Printf("failed to get profile data, it should be refreshed")
			pages.WrappedNoResults(
				mapProfileData(githubData),
				"Something went wrong when accessing profile info. You can either reload page or retry login").Render(r.Context(), w)
			return
		}

		pages.ProfilePage(mapProfileData(githubData)).Render(r.Context(), w)
	})
}

func (ph *ProfileHandlers) Logout(w http.ResponseWriter, r *http.Request) {
	ph.Handle(w, r, func(w http.ResponseWriter, r *http.Request) {
		t := r.Context().Value(middlewares.CookieAccessTokenKey)
		token := t.(string)

		log.Println("logging out!!!")

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
