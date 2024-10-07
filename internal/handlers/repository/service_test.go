package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/CristianCurteanu/gh-search/pkg/cache"
	"github.com/CristianCurteanu/gh-search/pkg/githubapi"
	"github.com/CristianCurteanu/gh-search/pkg/httpclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSearchRepositories_SuccessFromApi(t *testing.T) {
	token := uuid.NewString()
	query := url.Values{}
	query.Add("ownerType", "user")
	query.Add("ownerName", "CristianCurteanu")
	query.Add("repoQuery", "gh-search")

	ghApi := githubapi.NewGithubClient(uuid.NewString(), uuid.NewString())
	clientMock := new(httpclient.HttpClientMock)
	ghApi.WithClient(clientMock)

	searchRes := githubapi.RepositorySearchResult{
		Total: 1,
		Items: []*githubapi.Repository{
			{
				Id:       256,
				NodeID:   uuid.NewString(),
				Name:     "gh-search",
				FullName: "CristianCurteanu/gh-search",
			},
		},
	}
	respBytes, err := json.Marshal(&searchRes)
	require.NoError(t, err)

	bodyReader := io.NopCloser(bytes.NewReader(respBytes))
	clientMock.On("Do", mock.Anything).Return(&http.Response{
		Body: bodyReader,
	}, nil)

	q, page, err := prepareQueryString(query)
	require.NoError(t, err)
	key := fmt.Sprintf("%s:%s", token, encodedSearchQueryString(q, page))

	cacheMock := new(cache.CacheMock)
	cacheMock.On("Exists", mock.Anything, key).Return(false)
	cacheMock.On("Set", mock.Anything, key, searchRes).Return(nil)

	svc := service{ghApi, cacheMock}
	res, page, err := svc.searchRepositories(context.Background(), token, query)
	require.Equal(t, page, "1")
	require.NotEmpty(t, res.Items)
	require.Equal(t, res.Items[0].FullName, "CristianCurteanu/gh-search")
	require.NoError(t, err)
}

func TestSearchRepositories_FailsFromApi(t *testing.T) {
	token := uuid.NewString()
	query := url.Values{}
	query.Add("ownerType", "user")
	query.Add("ownerName", "CristianCurteanu")
	query.Add("repoQuery", "gh-search")

	ghApi := githubapi.NewGithubClient(uuid.NewString(), uuid.NewString())
	clientMock := new(httpclient.HttpClientMock)
	ghApi.WithClient(clientMock)

	clientMock.On("Do", mock.Anything).Return(&http.Response{}, errors.New("failed request"))

	q, page, err := prepareQueryString(query)
	require.NoError(t, err)
	key := fmt.Sprintf("%s:%s", token, encodedSearchQueryString(q, page))

	cacheMock := new(cache.CacheMock)
	cacheMock.On("Exists", mock.Anything, key).Return(false)

	svc := service{ghApi, cacheMock}
	res, page, err := svc.searchRepositories(context.Background(), token, query)
	require.Error(t, err)
	require.Equal(t, page, "1")
	require.Empty(t, res.Items)
}

func TestSearchRepositories_FailsIfOwnerNameNotProvided(t *testing.T) {
	query := url.Values{}
	query.Add("ownerType", "user")
	query.Add("repoQuery", "gh-search")

	ghApi := githubapi.NewGithubClient(uuid.NewString(), uuid.NewString())
	clientMock := new(httpclient.HttpClientMock)
	ghApi.WithClient(clientMock)

	clientMock.On("Do", mock.Anything).Return(&http.Response{}, errors.New("failed request"))

	_, _, err := prepareQueryString(query)
	require.Error(t, err)
}

func TestSearchRepositories_FailsIfRepoQueryNotProvided(t *testing.T) {
	token := uuid.NewString()
	query := url.Values{}

	ghApi := githubapi.NewGithubClient(uuid.NewString(), uuid.NewString())
	clientMock := new(httpclient.HttpClientMock)
	ghApi.WithClient(clientMock)

	clientMock.On("Do", mock.Anything).Return(&http.Response{}, errors.New("failed request"))

	_, _, err := prepareQueryString(query)
	require.Error(t, err)

	svc := service{ghApi, nil}
	res, page, err := svc.searchRepositories(context.Background(), token, query)
	require.Error(t, err)
	require.Equal(t, page, "1")
	require.Empty(t, res.Items)
}
