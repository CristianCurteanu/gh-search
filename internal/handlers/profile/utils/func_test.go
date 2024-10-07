package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/CristianCurteanu/gh-search/pkg/cache"
	"github.com/CristianCurteanu/gh-search/pkg/githubapi"
	"github.com/CristianCurteanu/gh-search/pkg/httpclient"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetProfileData_SuccessGHData(t *testing.T) {
	ghApi := githubapi.NewGithubClient(uuid.NewString(), uuid.NewString())
	clientMock := new(httpclient.HttpClientMock)
	ghApi.WithClient(clientMock)

	profile := githubapi.ProfileData{
		Username: "test-username",
	}
	respBytes, err := json.Marshal(&profile)
	require.NoError(t, err)

	bodyReader := io.NopCloser(bytes.NewReader(respBytes))
	clientMock.On("Do", mock.Anything).Return(&http.Response{
		Body: bodyReader,
	}, nil)

	key := uuid.NewString()

	cacheMock := new(cache.CacheMock)
	cacheMock.On("Exists", mock.Anything, key).Return(false)
	cacheMock.On("Set", mock.Anything, key, profile).Return(nil)

	data, err := GetProfileData(context.Background(), key, ghApi, cacheMock)
	require.NoError(t, err)
	require.Equal(t, data.Username, profile.Username)
}

func TestGetProfileData_FailedGHRequest(t *testing.T) {
	ghApi := githubapi.NewGithubClient(uuid.NewString(), uuid.NewString())
	clientMock := new(httpclient.HttpClientMock)
	ghApi.WithClient(clientMock)

	clientMock.On("Do", mock.Anything).Return(&http.Response{}, errors.New("failed gh request"))

	key := uuid.NewString()

	cacheMock := new(cache.CacheMock)
	cacheMock.On("Exists", mock.Anything, key).Return(false)

	data, err := GetProfileData(context.Background(), key, ghApi, cacheMock)
	require.Error(t, err)
	require.Equal(t, data.Username, "")
}

func TestGetProfileData_SuccessCacheData(t *testing.T) {
	ghApi := githubapi.NewGithubClient(uuid.NewString(), uuid.NewString())
	clientMock := new(httpclient.HttpClientMock)
	ghApi.WithClient(clientMock)

	username := "test-username"
	key := uuid.NewString()

	cacheMock := new(cache.CacheMock)
	cacheMock.On("Exists", mock.Anything, key).Return(true)
	cacheMock.On("Get", mock.Anything, key, mock.AnythingOfType("*githubapi.ProfileData")).Return(nil).Run(func(args mock.Arguments) {
		arg := args.Get(2).(*githubapi.ProfileData)

		arg.Username = username
	})

	data, err := GetProfileData(context.Background(), key, ghApi, cacheMock)
	require.NoError(t, err)
	require.Equal(t, data.Username, username)
}

func TestGetProfileData_FailCacheFallbackGithub(t *testing.T) {
	ghApi := githubapi.NewGithubClient(uuid.NewString(), uuid.NewString())
	clientMock := new(httpclient.HttpClientMock)
	ghApi.WithClient(clientMock)

	username := "test-username"
	key := uuid.NewString()

	cacheMock := new(cache.CacheMock)
	cacheMock.On("Exists", mock.Anything, key).Return(true)
	cacheMock.On("Get", mock.Anything, key, mock.AnythingOfType("*githubapi.ProfileData")).Return(nil).Return(errors.New("failed"))

	profile := githubapi.ProfileData{
		Username: username,
	}
	respBytes, err := json.Marshal(&profile)
	require.NoError(t, err)

	bodyReader := io.NopCloser(bytes.NewReader(respBytes))
	clientMock.On("Do", mock.Anything).Return(&http.Response{
		Body: bodyReader,
	}, nil)

	data, err := GetProfileData(context.Background(), key, ghApi, cacheMock)
	require.NoError(t, err)
	require.Equal(t, data.Username, username)
}

func TestGetProfileData_FailCacheFallbackGithubFails(t *testing.T) {
	ghApi := githubapi.NewGithubClient(uuid.NewString(), uuid.NewString())
	clientMock := new(httpclient.HttpClientMock)
	ghApi.WithClient(clientMock)

	key := uuid.NewString()

	cacheMock := new(cache.CacheMock)
	cacheMock.On("Exists", mock.Anything, key).Return(true)
	cacheMock.On("Get", mock.Anything, key, mock.AnythingOfType("*githubapi.ProfileData")).Return(errors.New("failed"))

	clientMock.On("Do", mock.Anything).Return(&http.Response{}, errors.New("failed gh request"))

	data, err := GetProfileData(context.Background(), key, ghApi, cacheMock)
	require.Error(t, err)
	require.Equal(t, data.Username, "")
}
