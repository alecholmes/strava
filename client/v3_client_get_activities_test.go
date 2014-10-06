package client

import (
	"testing"

	"github.com/alecholmes/strava/model"
)

func TestGetActivities(t *testing.T) {
	client, rawClient := newTestClient()

	count := 10
	activities := make(map[model.ActivityId]*model.Activity, count)
	activityIds := make([]model.ActivityId, 0, count)
	for i := 1; i <= count; i++ {
		activity := model.Activity{Id: model.ActivityId(i)}
		activities[activity.Id] = &activity
		activityIds = append(activityIds, activity.Id)
		rawClient.Gets[fullActivityUrl(activity.Id, rawClient, t)] = ExpectedBody([]byte(toJson(activity, t)))
	}

	fetched, err := client.GetActivities(activityIds)
	if err != nil {
		t.Fatalf("Error calling GetActivities: %s", err)
	}

	if len(fetched) != len(activities) {
		t.Fatalf("Fetched different number of activities than requested. expectedCount=%d, actualCount=%d", len(activities), len(fetched))
	}

	for i, activityId := range activityIds {
		if activityId != fetched[i].Id {
			t.Fatalf("Fetched activity was different than expected. expectedId=%d, actualId=%d", activityId, fetched[i].Id)
		}
	}
}
