package githubapi

import (
	"fmt"
	"net/http"
	"time"

	"github.com/CristianCurteanu/gh-search/internal/httpclient"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type GithubApi struct {
	clientId     string
	clientSecret string
	url          string
	httpClient   HTTPClient
}

func NewGithubClient(clientId, clientSecret string) *GithubApi {
	return &GithubApi{
		clientId:     clientId,
		clientSecret: clientSecret,
		url:          "https://api.github.com",
	}

}

func (a *GithubApi) SetHost(u string) {
	a.url = u
}

func (a *GithubApi) WithClient(httpClient HTTPClient) {
	a.httpClient = httpClient
}

func (a *GithubApi) SearchRepository(accessToken string, query string) (RepositorySearchResult, error) {
	return httpclient.NewJsonRequest[RepositorySearchResult](nil).
		SetTimeout(10*time.Second).
		SetHeader("Authorization", fmt.Sprintf("token %s", accessToken)).
		Do(
			http.MethodGet,
			fmt.Sprintf("%s/search/repositories?%s", a.url, query), nil,
		)
}

func (a *GithubApi) GetProfileInfo(accessToken string) (ProfileData, error) {
	return httpclient.NewJsonRequest[ProfileData](nil).
		SetTimeout(10*time.Second).
		SetHeader("Authorization", fmt.Sprintf("token %s", accessToken)).
		Do(
			http.MethodGet,
			fmt.Sprintf("%s/user", a.url), nil,
		)
}

func (a *GithubApi) GetRepositoryInfo(accessToken, fullRepoName string) (Repository, error) {
	return httpclient.NewJsonRequest[Repository](nil).SetTimeout(10*time.Second).
		SetHeader("Authorization", fmt.Sprintf("token %s", accessToken)).
		Do(
			http.MethodGet,
			fmt.Sprintf("%s/repos/%s", a.url, fullRepoName), nil,
		)
}

func (a *GithubApi) GetRepoContributors(accessToken, fullRepoName string) (Contributors, error) {
	contributors, err := httpclient.NewJsonRequest[*Contributors](nil).SetTimeout(10*time.Second).
		SetHeader("Authorization", fmt.Sprintf("token %s", accessToken)).
		Do(
			http.MethodGet,
			fmt.Sprintf("%s/repos/%s/contributors", a.url, fullRepoName), nil,
		)
	return *contributors, err
}

func (a *GithubApi) GetRepoCommits(accessToken, fullRepoName string) (Commits, error) {
	commits, err := httpclient.NewJsonRequest[*Commits](nil).SetTimeout(10*time.Second).
		SetHeader("Authorization", fmt.Sprintf("token %s", accessToken)).
		Do(
			http.MethodGet,
			fmt.Sprintf("%s/repos/%s/commits", a.url, fullRepoName), nil,
		)
	return *commits, err
}

func (a *GithubApi) GetGithubAccessToken(code string) (AccessTokenResponse, error) {
	req := httpclient.NewJsonRequest[AccessTokenResponse](nil).
		SetTimeout(10*time.Second).
		SetHeader("Accept", "application/json")

	fmt.Printf("api: %+v\n", a)
	resp, err := req.Do(http.MethodPost, "https://github.com/login/oauth/access_token",
		map[string]string{
			"client_id":     a.clientId,
			"client_secret": a.clientSecret,
			"code":          code,
		},
	)

	return resp, err
}
