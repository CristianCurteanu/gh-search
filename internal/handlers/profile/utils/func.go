package utils

import (
	"context"

	"github.com/CristianCurteanu/gh-search/pkg/cache"
	"github.com/CristianCurteanu/gh-search/pkg/githubapi"
)

func GetProfileData(ctx context.Context, token string, client *githubapi.GithubApi, cache cache.Cache) (githubapi.ProfileData, error) {
	var githubData githubapi.ProfileData
	var err error
	if !cache.Exists(ctx, token) {
		githubData, err = client.GetProfileInfo(token)
		if err != nil {
			return githubData, err
		}
		cache.Set(ctx, token, githubData)
	} else {
		err := cache.Get(ctx, token, &githubData)
		if err != nil {
			githubData, err = client.GetProfileInfo(token)
			if err != nil {
				return githubData, err
			}
		}
	}
	return githubData, err
}
