package repository

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/CristianCurteanu/gh-search/internal/components"
	"github.com/CristianCurteanu/gh-search/internal/handlers/repository/pages"
	"github.com/CristianCurteanu/gh-search/internal/middlewares"
	"github.com/CristianCurteanu/gh-search/pkg/githubapi"
)

const (
	DefaultAvatarURL = "https://avatars.githubusercontent.com/u/19864447?v=4"
)

type RepositoriesHandlers struct {
	githubClient *githubapi.GithubApi
	middlewares.UseMiddleware
}

func NewRepositoriesHandlers(githubClient *githubapi.GithubApi) *RepositoriesHandlers {
	return &RepositoriesHandlers{
		githubClient: githubClient,
	}
}

func (ph *RepositoriesHandlers) Search(w http.ResponseWriter, r *http.Request) {
	ph.Handle(w, r, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")

		t := r.Context().Value(middlewares.CookieAccessTokenKey)
		token := t.(*http.Cookie)

		ownerType := r.URL.Query().Get("ownerType")
		ownerName := r.URL.Query().Get("ownerName")
		repoQuery := r.URL.Query().Get("repoQuery")

		page := r.URL.Query().Get("page")
		if page == "" {
			page = "1"
		}

		var queryStringBuf bytes.Buffer

		if ownerType != "" && ownerName != "" {
			queryStringBuf.WriteString(fmt.Sprintf("%s:%s ", ownerType, ownerName))
		} else if ownerType != "" && ownerName == "" {
			pages.NoResults("If you want to search using an owner type, specify the owner type").Render(r.Context(), w)
			return
		}

		if repoQuery == "" {
			pages.NoResults("Please, use set a repo query").Render(r.Context(), w)
			return
		}

		queryStringBuf.WriteString(repoQuery)
		queryString := queryStringBuf.String()

		params := url.Values{}
		params.Add("q", queryString)
		params.Add("type", "repositories")
		params.Add("page", page)

		queryString = params.Encode()
		fmt.Printf("queryString(): %v\n", queryString)

		githubData, err := ph.githubClient.SearchRepository(token.Value, queryString)
		if err != nil {
			pages.NoResults("Repositories Not found").Render(r.Context(), w)
			return
		}

		currentPage, _ := strconv.Atoi(page)
		prevPage := currentPage - 1
		if prevPage <= 0 {
			prevPage = -1
		}
		nextPage := currentPage + 1
		if nextPage >= githubData.Total {
			nextPage = -1
		}
		total := githubData.Total / 30
		if githubData.Total%30 != 0 {
			total += 1
		}

		pages.SearchResult(mapStruct(githubData.Items, func(rd *githubapi.Repository) components.Repository {
			return mapRepository(*rd)
		}), prevPage, currentPage, nextPage, total).Render(r.Context(), w)
	})
}

func (ph *RepositoriesHandlers) GetRepositoryPage(w http.ResponseWriter, r *http.Request) {
	ph.Handle(w, r, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/html")

		t := r.Context().Value(middlewares.CookieAccessTokenKey)
		token := t.(*http.Cookie)
		githubData, err := ph.githubClient.GetProfileInfo(token.Value)
		if err != nil {
			log.Printf("failed to get profile data, err=%q", err)
			pages.WrappedNoResults(mapProfileData(githubData), "Failed to get profile data from Github").Render(r.Context(), w)
			return
		}

		owner := r.URL.Query().Get("owner")
		repo := r.URL.Query().Get("repo")

		if owner == "" {
			pages.WrappedNoResults(mapProfileData(githubData), "Please, define owner, it's empty now").Render(r.Context(), w)
			return
		}

		if repo == "" {
			pages.WrappedNoResults(mapProfileData(githubData), "Please, define repository, it's empty now").Render(r.Context(), w)
			return
		}

		repoData, err := ph.githubClient.GetRepositoryInfo(
			token.Value,
			fmt.Sprintf("%s/%s", owner, repo),
		)
		if err != nil {
			log.Printf("failed to get repository info, err=%q", err)
			pages.WrappedNoResults(mapProfileData(githubData), "Failed to get repository data").Render(r.Context(), w)
			return
		}

		commits, err := ph.githubClient.GetRepoCommits(
			token.Value,
			fmt.Sprintf("%s/%s", owner, repo),
		)
		if err != nil {
			log.Printf("failed to get repository info, err=%q", err)
			pages.WrappedNoResults(mapProfileData(githubData), "Failed to get repository commits data").Render(r.Context(), w)
			return
		}

		contributors, err := ph.githubClient.GetRepoContributors(
			token.Value,
			fmt.Sprintf("%s/%s", owner, repo),
		)
		if err != nil {
			log.Printf("failed to get repository info, err=%q", err)
			pages.WrappedNoResults(mapProfileData(githubData), "Failed to get repository commits data").Render(r.Context(), w)
			return
		}

		pages.RepositoryDetailsPage(pages.RepositoryDetails{
			Profile:      mapProfileData(githubData),
			Repo:         mapRepository(repoData),
			Commits:      mapStruct(commits, mapCommit),
			Contributors: mapStruct(contributors, mapContributors),
		}).Render(r.Context(), w)
	})
}

func mapStruct[I any, O any](input []I, cb func(I) O) []O {
	var output []O = make([]O, 0, len(input))
	for _, i := range input {
		output = append(output, cb(i))
	}

	return output
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

func mapRepository(repo githubapi.Repository) components.Repository {
	return components.Repository{
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

func mapProfileData(githubData githubapi.ProfileData) components.ProfileData {
	return components.ProfileData{
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
