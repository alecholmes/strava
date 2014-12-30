package client

import (
	"encoding/json"
	"testing"

	"github.com/alecholmes/strava/model"
)

func newTestClient() (*v3Client, *testHttpClient) {
	testHttpClient := newTestHttpClient()
	return &v3Client{httpClient: testHttpClient}, testHttpClient
}

func toJson(obj interface{}, t *testing.T) []byte {
	jsonStr, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("Could not marshal to JSON. object=%s, error=%s", obj, err)
	}
	return jsonStr
}

func fullActivityUrl(activityId model.ActivityId, rawClient *testHttpClient, t *testing.T) string {
	url, err := rawClient.AbsoluteUrl(activityUrl(activityId), make(map[string]interface{}))
	if err != nil {
		t.Errorf("Error creating test URL. activityId=%s, error=%s", activityId, err)
	}

	return url
}
