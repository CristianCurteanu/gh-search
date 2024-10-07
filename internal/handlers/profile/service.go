package profile

import (
	"github.com/CristianCurteanu/gh-search/pkg/cache"
	"github.com/CristianCurteanu/gh-search/pkg/githubapi"
)

type service struct {
	githubClient *githubapi.GithubApi
	cache        cache.Cache
}
