package client

import (
	"github.com/alecholmes/strava/model"
)

const Beginning = 0

type Client interface {
	// Get activity summaries. Oldest are returned first.
	// Only activites with ids greater than the given after are included.
	// Pass Beginning as after to get all activities.
	GetActivitySummaries(after model.ActivityId) ([]*model.ActivitySummary, error)

	// Get an activity by its id.
	GetActivity(activityId model.ActivityId) (*model.Activity, error)

	// Get multiple activities by their ids, returned in the same order.
	// Activities that could not be fetched are excluded.
	GetActivities(activityIds []model.ActivityId) ([]*model.Activity, error)

	// Get activities summaries for activities related to the given activity id.
	GetRelatedActivitySummaries(activityId model.ActivityId) ([]*model.ActivitySummary, error)
}

func NewClient(accessToken string) *v3Client {
	return &v3Client{httpClient: newHttpClientImpl(stravaBaseUrl, accessToken)}
}
