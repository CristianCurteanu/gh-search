package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type SerializationFunc func(any) ([]byte, error)
type DeserializationFunc func(data []byte, v any) (bool, error)

type JSONRequest[T any] struct {
	headers     map[string]string
	client      *http.Client
	serialize   SerializationFunc
	deserialize DeserializationFunc
}

func NewJsonRequest[RESP any](client *http.Client) *JSONRequest[RESP] {
	if client == nil {
		client = http.DefaultClient
	}
	return &JSONRequest[RESP]{
		client:  client,
		headers: map[string]string{},
		serialize: func(a any) ([]byte, error) {
			return json.Marshal(a)
		},
		deserialize: func(data []byte, v any) (bool, error) {
			err := json.Unmarshal(data, v)
			return err != nil, err
		},
	}
}

func (rq *JSONRequest[RESP]) SetHeader(header, headerVal string) *JSONRequest[RESP] {
	rq.headers[header] = headerVal

	return rq
}

func (rq *JSONRequest[RESP]) SetTimeout(t time.Duration) *JSONRequest[RESP] {
	rq.client.Timeout = t

	return rq
}

func (rq *JSONRequest[RESP]) WithSerialize(sf SerializationFunc) *JSONRequest[RESP] {
	rq.serialize = sf

	return rq
}

func (rq *JSONRequest[RESP]) WithDeserialize(df DeserializationFunc) *JSONRequest[RESP] {
	rq.deserialize = df

	return rq
}

func (rq *JSONRequest[RESP]) Do(method, url string, reqBody any) (RESP, error) {
	var resp RESP
	var body []byte
	var err error

	if reqBody != nil {
		body, err = rq.serialize(reqBody)
	}
	if err != nil {
		return resp, err
	}

	var reqReader io.Reader

	if body != nil {
		reqReader = bytes.NewBuffer(body)
	}

	req, err := http.NewRequest(method, url, reqReader)
	if err != nil {
		return resp, err
	}
	req.Header.Set("Content-Type", "application/json")
	for header, headerVal := range rq.headers {
		req.Header.Set(header, headerVal)
	}

	response, err := rq.client.Do(req)
	if err != nil {
		return resp, err
	}
	defer response.Body.Close()

	body, err = io.ReadAll(response.Body)
	if err != nil {
		return resp, err
	}

	// TODO: handle different status classes

	_, err = rq.deserialize(body, &resp)

	return resp, err
}
