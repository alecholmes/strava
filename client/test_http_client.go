package client

import (
	"fmt"
)

type BodyOrError struct {
	Body  []byte
	Error error
}

func ExpectedBody(body []byte) BodyOrError {
	return BodyOrError{Body: body, Error: nil}
}

type TestHttpClient struct {
	HttpClient
	Gets map[string]BodyOrError
}

func NewTestHttpClient() *TestHttpClient {
	client := newHttpClientImpl("http://test", "fake-access-token")
	return &TestHttpClient{client, make(map[string]BodyOrError)}
}

func (client *TestHttpClient) Get(relativePath string, params map[string]interface{}) ([]byte, error) {
	absoluteUrl, err := client.AbsoluteUrl(relativePath, params)
	if err != nil {
		return nil, err
	}

	bodyOrError, ok := client.Gets[absoluteUrl]
	if !ok {
		panic(fmt.Sprintf("Gets did not contain %s", absoluteUrl))
	}

	return bodyOrError.Body, bodyOrError.Error
}
