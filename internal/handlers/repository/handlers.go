package repository

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/CristianCurteanu/gh-search/internal/handlers/profile/utils"
	"github.com/CristianCurteanu/gh-search/internal/handlers/repository/pages"
	"github.com/CristianCurteanu/gh-search/internal/layouts"
	"github.com/CristianCurteanu/gh-search/internal/middlewares"
	"github.com/CristianCurteanu/gh-search/pkg/cache"
	"github.com/CristianCurteanu/gh-search/pkg/githubapi"
	"github.com/CristianCurteanu/gh-search/pkg/slices"
)

const (
	DefaultAvatarURL = "https://avatars.githubusercontent.com/u/19864447?v=4"
)

type RepositoriesHandlers struct {
	githubClient *githubapi.GithubApi
	cache        cache.Cache
	svc          service

	middlewares.UseMiddleware
}

func NewRepositoriesHandlers(githubClient *githubapi.GithubApi, cache cache.Cache) *RepositoriesHandlers {
	return &RepositoriesHandlers{
		githubClient: githubClient,
		cache:        cache,
		svc:          service{githubClient, cache},
	}
}

func (ph *RepositoriesHandlers) Search(w http.ResponseWriter, r *http.Request) {
	ph.Handle(w, r, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")

		t := r.Context().Value(middlewares.CookieAccessTokenKey)
		token := t.(string)

		searchResults, page, err := ph.svc.searchRepositories(r.Context(), token, r.URL.Query())
		if err != nil {
			pages.NoResults(err.Error()).Render(r.Context(), w)
			return
		}

		currentPage, _ := strconv.Atoi(page)
		pages.SearchResult(pages.SearchResultsData{
			Items: slices.MapSlice(searchResults.Items, func(rd *githubapi.Repository) layouts.Repository {
				return mapRepository(*rd)
			}),
			CurrentPage: currentPage,
			TotalPages:  searchResults.Total,
		}).Render(r.Context(), w)
	})
}

func (ph *RepositoriesHandlers) GetRepositoryPage(w http.ResponseWriter, r *http.Request) {
	ph.Handle(w, r, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")

		t := r.Context().Value(middlewares.CookieAccessTokenKey)
		token := t.(string)

		profileData, err := utils.GetProfileData(r.Context(), token, ph.githubClient, ph.cache)
		if err != nil {
			log.Printf("failed to get profile data, err=%q", err)
			pages.WrappedNoResults(mapProfileData(profileData), "Failed to get profile data from Github").Render(r.Context(), w)
			return
		}

		repoDetails, err := ph.svc.getRepoDetails(r.Context(), token, r.URL.Query())
		if err != nil {
			pages.WrappedNoResults(mapProfileData(profileData), err.Error()).Render(r.Context(), w)
			return
		}

		pages.RepositoryDetailsPage(pages.RepositoryDetails{
			Profile:      mapProfileData(profileData),
			Repo:         mapRepository(repoDetails.repoData),
			Commits:      slices.MapSlice(repoDetails.commits, mapCommit),
			Contributors: slices.MapSlice(repoDetails.contributors, mapContributors),
		}).Render(r.Context(), w)
	})
}

func mapCommit(commit githubapi.Commit) pages.Commit {
	res := pages.Commit{
		AuthorName: commit.Commit.Author.Name,
		Url:        commit.HTMLUrl,
		Message:    commit.Commit.Message,
		Sha:        commit.Sha[:7],
		CommitedAt: prettifyDate(commit.Commit.Author.Date),
	}

	if commit.Author != nil {
		res.AuthorAvatar = commit.Author.AvatarURL
	} else if commit.Commiter != nil {
		res.AuthorAvatar = commit.Commiter.AvatarURL
	} else {
		res.AuthorAvatar = DefaultAvatarURL
	}

	return res
}

func mapContributors(c githubapi.Contributor) pages.Contributor {
	return pages.Contributor{
		Username:  c.Username,
		AvatarURL: c.AvatarURL,
		HtmlUrl:   c.HtmlUrl,
	}
}

func mapRepository(repo githubapi.Repository) layouts.Repository {
	return layouts.Repository{
		Name:        repo.Name,
		FullName:    repo.FullName,
		Description: repo.Description,
		OwnerAvatar: repo.Owner.AvatarURL,
		OwnerName:   repo.Owner.Username,
		Stars:       fmt.Sprintf("%d", repo.Stars),
		Forks:       fmt.Sprintf("%d", repo.Forks),
		Watchers:    fmt.Sprintf("%d", repo.Watchers),
		UpdatedAt:   prettifyDate(repo.PushedAt),
		Language:    repo.Language,
		Url:         repo.HtmlUrl,
	}
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

func prettifyDate(t *time.Time) string {
	return fmt.Sprintf("%s %d, %d",
		t.Month().String()[0:3],
		t.Day(),
		t.Year(),
	)
}
