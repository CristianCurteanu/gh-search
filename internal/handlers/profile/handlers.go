package profile

import (
	"fmt"
	"net/http"

	"github.com/CristianCurteanu/gh-search/internal/components"
	"github.com/CristianCurteanu/gh-search/internal/middlewares"
)

type ProfileHandlers struct {
	middlewares.UseMiddleware
}

func NewProfileHandlers() *ProfileHandlers {
	return &ProfileHandlers{}
}

func (ph *ProfileHandlers) GetProfilePage(w http.ResponseWriter, r *http.Request) {
	ph.Handle(w, r, func(w http.ResponseWriter, r *http.Request) {
		t := r.Context().Value(middlewares.CookieAccessTokenKey)
		token := t.(*http.Cookie)

		githubData, _ := getGithubData(token.Value)

		fmt.Printf("githubData: %+v\n", githubData)
		// profileHandler(w, r, githubData)
		components.ProfilePage(githubData).Render(r.Context(), w)
	})
}

func (ph *ProfileHandlers) Logout(w http.ResponseWriter, r *http.Request) {}
