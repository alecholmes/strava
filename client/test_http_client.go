package client

import (
	"fmt"
)

type bodyOrError struct {
	Body  []byte
	Error error
}

func expectedBody(body []byte) bodyOrError {
	return bodyOrError{Body: body, Error: nil}
}

type testHttpClient struct {
	HttpClient
	Gets map[string]bodyOrError
}

// Assert testHttpClient implements HttpClient
var _ HttpClient = &testHttpClient{}

func newTestHttpClient() *testHttpClient {
	client := newHttpClientImpl("http://test", "fake-access-token")
	return &testHttpClient{client, make(map[string]bodyOrError)}
}

func (client *testHttpClient) Get(relativePath string, params map[string]interface{}) ([]byte, error) {
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
