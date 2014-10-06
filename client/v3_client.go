package client

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/alecholmes/strava/model"
)

const stravaBaseUrl = "https://www.strava.com/api/v3"

const activitySummariesUrl = "/athlete/activities"
const activitySummariesPageSize = 100

// Number of concurrent get requests allowed
const getActivitiesPoolSize = 10

// A client backed by Strava's HTTP API V3
// http://strava.github.io/api/
type V3Client struct {
	httpClient HttpClient
}

func (c *V3Client) GetActivitySummaries(after model.ActivityId) ([]*model.ActivitySummary, error) {
	allActivities := make([]*model.ActivitySummary, 0)
	complete := false
	for page := 1; !complete; page++ {
		activities, err := c.activitiesPage(page)
		if err != nil {
			return nil, err
		}

		filtered := filterBefore(activities, after)
		allActivities = append(allActivities, filtered...)

		complete = len(activities) == 0 || len(activities) > len(filtered)
	}

	return reverse(allActivities), nil
}

func (c *V3Client) GetActivity(activityId model.ActivityId) (*model.Activity, error) {
	body, err := c.httpClient.Get(activityUrl(activityId), make(map[string]interface{}))
	if err != nil {
		return nil, err
	}

	var activity model.Activity
	if err := json.Unmarshal(body, &activity); err != nil {
		return nil, err
	}

	return &activity, nil
}

func (c *V3Client) GetActivities(activityIds []model.ActivityId) ([]*model.Activity, error) {
	activityMap := make(map[model.ActivityId]*model.Activity, len(activityIds))

	activityIdsLeft := activityIds[:]
	fetchedActivities := make(chan *model.Activity)
	runningRequests := 0
	for len(activityIdsLeft) > 0 {
		if runningRequests == getActivitiesPoolSize {
			if fetched := <-fetchedActivities; fetched != nil {
				activityMap[fetched.Id] = fetched
			}
			runningRequests--
		}
		runningRequests++
		activityId := activityIdsLeft[0]
		activityIdsLeft = activityIdsLeft[1:]
		go func() {
			activity, err := c.GetActivity(activityId)
			if err != nil {
				fetchedActivities <- nil
			} else {
				fetchedActivities <- activity
			}
		}()
	}

	// Drain any outstanding requests
	for i := 0; i < runningRequests; i++ {
		if fetched := <-fetchedActivities; fetched != nil {
			activityMap[fetched.Id] = fetched
		}
	}

	// Reorder to match activityIds
	activities := make([]*model.Activity, 0, len(activityMap))
	for _, activityId := range activityIds {
		activity, found := activityMap[activityId]
		if found {
			activities = append(activities, activity)
		}
	}

	return activities, nil
}

func activityUrl(activityId model.ActivityId) string {
	return fmt.Sprintf("/activities/%d", activityId)
}

func (c *V3Client) activitiesPage(page int) ([]*model.ActivitySummary, error) {
	if page <= 0 {
		return nil, errors.New("page must be positive")
	}

	body, err := c.httpClient.Get(activitySummariesUrl, map[string]interface{}{"per_page": activitySummariesPageSize, "page": page})
	if err != nil {
		return nil, err
	}

	summaries := make([]*model.ActivitySummary, 0)
	if err := json.Unmarshal(body, &summaries); err != nil {
		return nil, err
	}

	return summaries, nil
}

// Filter out activities having an id <= the given id
func filterBefore(activities []*model.ActivitySummary, activityId model.ActivityId) []*model.ActivitySummary {
	if len(activities) == 0 {
		return activities
	}

	keepBefore := len(activities)
	for i, v := range activities {
		if v.Id <= activityId {
			keepBefore = i
			break
		}
	}
	return activities[:keepBefore]
}

// This creates a copy but could really just reverse in place (requiring writing even more boilerplate code)
func reverse(summaries []*model.ActivitySummary) []*model.ActivitySummary {
	reversed := make([]*model.ActivitySummary, len(summaries))
	for i, s := range summaries {
		reversed[len(summaries)-i-1] = s
	}

	return reversed
}
