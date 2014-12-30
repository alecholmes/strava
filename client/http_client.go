package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Internal HTTP client. Interface to allow for testing implementations.
type HttpClient interface {
	AbsoluteUrl(relativeUrl string, params map[string]interface{}) (string, error)

	Get(relativePath string, params map[string]interface{}) ([]byte, error)
}

type httpClientImpl struct {
	baseUrl     string
	accessToken string
	httpClient  *http.Client
}

// Create a new HTTP client that uses the given accessToken for authentication.
func newHttpClientImpl(baseUrl string, accessToken string) HttpClient {
	return &httpClientImpl{baseUrl: baseUrl, accessToken: accessToken, httpClient: http.DefaultClient}
}

func (c *httpClientImpl) AbsoluteUrl(relativePath string, params map[string]interface{}) (string, error) {
	absUrl, err := url.Parse(c.baseUrl + relativePath)
	if err != nil {
		return "", err
	}

	urlParams := url.Values{}
	for k, v := range params {
		urlParams.Add(k, fmt.Sprintf("%v", v))
	}
	absUrl.RawQuery = urlParams.Encode()

	return absUrl.String(), nil
}

func (client *httpClientImpl) Get(relativePath string, params map[string]interface{}) ([]byte, error) {
	absUrl, err := client.AbsoluteUrl(relativePath, params)
	if err != nil {
		return nil, err
	}

	request, _ := http.NewRequest("GET", absUrl, nil)
	request.Header.Set("Authorization", client.bearerToken())
	request.Header.Set("Accept", "application/json")

	response, err := client.httpClient.Do(request)
	if err != nil {
		return nil, err
	} else if response.StatusCode != 200 {
		return nil, fmt.Errorf("Unexpected HTTP response %v", response.Status)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (client *httpClientImpl) bearerToken() string {
	return fmt.Sprint("Bearer ", client.accessToken)
}
