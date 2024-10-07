package repository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"net/url"

	"github.com/CristianCurteanu/gh-search/pkg/cache"
	"github.com/CristianCurteanu/gh-search/pkg/githubapi"
)

type service struct {
	githubClient *githubapi.GithubApi
	cache        cache.Cache
}

func encodedSearchQueryString(q, page string) string {
	params := url.Values{}
	params.Add("q", q)
	params.Add("type", "repositories")
	params.Add("page", page)

	return params.Encode()
}

func prepareQueryString(query url.Values) (string, string, error) {
	ownerType := query.Get("ownerType")
	ownerName := query.Get("ownerName")
	repoQuery := query.Get("repoQuery")

	page := query.Get("page")
	if page == "" {
		page = "1"
	}

	var queryStringBuf bytes.Buffer

	if ownerType != "" && ownerName != "" {
		queryStringBuf.WriteString(fmt.Sprintf("%s:%s ", ownerType, ownerName))
	} else if ownerType != "" && ownerName == "" {
		return "", page, errors.New("If you want to search using an owner type, specify the owner type")
	}

	if repoQuery == "" {
		return "", page, errors.New("Please, use set a repo query")
	}

	queryStringBuf.WriteString(repoQuery)
	return queryStringBuf.String(), page, nil
}

func (s service) searchRepositories(ctx context.Context, accessToken string, query url.Values) (githubapi.RepositorySearchResult, string, error) {
	var err error
	var res githubapi.RepositorySearchResult

	queryString, page, err := prepareQueryString(query)
	if err != nil {
		return res, page, err
	}
	queryString = encodedSearchQueryString(queryString, page)

	key := fmt.Sprintf("%s:%s", accessToken, queryString)
	var searchResults githubapi.RepositorySearchResult

	if s.cache.Exists(ctx, key) {
		err := s.cache.Get(ctx, key, &searchResults)
		if err != nil {
			searchResults, err = s.githubClient.SearchRepository(accessToken, queryString)
			if err != nil {

				return res, page, errors.New("Repositories Not found")
			}
			s.cache.Set(ctx, key, searchResults)
		}
	} else {
		searchResults, err = s.githubClient.SearchRepository(accessToken, queryString)
		if err != nil {

			return res, page, errors.New("Repositories Not found")
		}
		s.cache.Set(ctx, key, searchResults)
	}

	return searchResults, page, nil
}

type repoDetails struct {
	repoData     githubapi.Repository
	commits      githubapi.Commits
	contributors githubapi.Contributors
}

type repoDetailsInput struct {
	token string
	repo  string
	owner string
}

func (rdi repoDetailsInput) GetKey() string {
	return fmt.Sprintf("%s:%s:%s", rdi.token, rdi.owner, rdi.repo)
}

func (rdi repoDetailsInput) GetRepoFullName() string {
	return fmt.Sprintf("%s/%s", rdi.owner, rdi.repo)
}

func (s service) getRepoDetails(ctx context.Context, token string, query url.Values) (repoDetails, error) {
	var details repoDetails
	var err error

	owner := query.Get("owner")
	repo := query.Get("repo")

	if owner == "" {
		return details, errors.New("Please, define owner, it's empty now")
	}

	if repo == "" {
		return details, errors.New("Please, define repository, it's empty now")
	}

	repoInput := repoDetailsInput{
		token: token,
		repo:  repo,
		owner: owner,
	}

	details.repoData, err = s.getRepositoryData(ctx, repoInput)
	if err != nil {
		log.Printf("failed to get repository info, err=%q", err)
		return details, errors.New("Failed to get repository data")
	}

	details.commits, err = s.getRepoCommits(ctx, repoInput)
	if err != nil {
		log.Printf("failed to get repository commits info, err=%q", err)
		return details, errors.New("Failed to get repository commits data")
	}

	details.contributors, err = s.getRepoContributors(ctx, repoInput)
	if err != nil {
		log.Printf("failed to get repository contributors info, err=%q", err)
		return details, errors.New("Failed to get repository contributors data")

	}
	return details, nil
}

func (s service) getRepositoryData(ctx context.Context, i repoDetailsInput) (githubapi.Repository, error) {
	var err error
	var data githubapi.Repository
	if !s.cache.Exists(ctx, i.GetKey()) {
		data, err = s.githubClient.GetRepositoryInfo(i.token, i.GetRepoFullName())
		if err != nil {
			return data, err
		}
		s.cache.Set(ctx, i.token, data)
	} else {
		err := s.cache.Get(ctx, i.token, &data)
		if err != nil {
			data, err = s.githubClient.GetRepositoryInfo(i.token, i.GetRepoFullName())
			if err != nil {
				return data, err
			}
		}
	}

	return data, nil
}

func (s service) getRepoCommits(ctx context.Context, i repoDetailsInput) (githubapi.Commits, error) {
	var err error
	var data githubapi.Commits
	if !s.cache.Exists(ctx, i.GetKey()) {
		data, err = s.githubClient.GetRepoCommits(i.token, i.GetRepoFullName())
		if err != nil {
			return data, err
		}
		s.cache.Set(ctx, i.token, data)
	} else {
		err := s.cache.Get(ctx, i.token, &data)
		if err != nil {
			data, err = s.githubClient.GetRepoCommits(i.token, i.GetRepoFullName())
			if err != nil {
				return data, err
			}
		}
	}

	return data, nil
}

func (s service) getRepoContributors(ctx context.Context, i repoDetailsInput) (githubapi.Contributors, error) {
	var err error
	var data githubapi.Contributors
	if !s.cache.Exists(ctx, i.GetKey()) {
		data, err = s.githubClient.GetRepoContributors(i.token, i.GetRepoFullName())
		if err != nil {
			return data, err
		}
		s.cache.Set(ctx, i.token, data)
	} else {
		err := s.cache.Get(ctx, i.token, &data)
		if err != nil {
			data, err = s.githubClient.GetRepoContributors(i.token, i.GetRepoFullName())
			if err != nil {
				return data, err
			}
		}
	}

	return data, nil
}
