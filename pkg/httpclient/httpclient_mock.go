package httpclient

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

type HttpClientMock struct {
	mock.Mock
}

func (hc *HttpClientMock) Do(req *http.Request) (*http.Response, error) {
	args := hc.Called(req)

	return args.Get(0).(*http.Response), args.Error(1)
}
