package githubapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CristianCurteanu/gh-search/internal/httpclient"
)

type GithubApi struct {
	clientId     string
	clientSecret string
	url          string
}

func NewGithubClient(clientId, clientSecret string) *GithubApi {
	return &GithubApi{
		clientId:     clientId,
		clientSecret: clientSecret,
		url:          "https://api.github.com",
	}

}

func (api *GithubApi) SetHost(u string) {
	api.url = u
}

func (api *GithubApi) SearchRepository(accessToken string, query string) (RepositorySearchResult, error) {
	return httpclient.NewJsonRequest[RepositorySearchResult](nil).
		SetTimeout(10*time.Second).
		SetHeader("Authorization", fmt.Sprintf("token %s", accessToken)).
		Do(
			http.MethodGet,
			fmt.Sprintf("%s/search/repositories?%s", api.url, query), nil,
		)
}

func (api *GithubApi) GetGithubAccessToken(code string) (AccessTokenResponse, error) {
	req := httpclient.NewJsonRequest[AccessTokenResponse](nil).
		SetTimeout(10*time.Second).
		SetHeader("Accept", "application/json")

	resp, err := req.Do(http.MethodPost, "https://github.com/login/oauth/access_token",
		map[string]string{
			"client_id":     api.clientId,
			"client_secret": api.clientSecret,
			"code":          code,
		},
	)

	return resp, err
}
