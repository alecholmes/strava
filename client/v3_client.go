package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/alecholmes/strava/model"
)

const stravaBaseUrl = "https://www.strava.com/api/v3"

const activitySummariesUrl = "/athlete/activities"
const activitySummariesPageSize = 100

// Number of concurrent get requests allowed
const getActivitiesPoolSize = 10

// A client backed by Strava's HTTP API V3
// http://strava.github.io/api/
type v3Client struct {
	httpClient HttpClient
}

// Assert v3Client implements Client
var _ Client = &v3Client{}

func (c *v3Client) GetActivitySummaries(after model.ActivityId) ([]*model.ActivitySummary, error) {
	return c.getActivitySummaries(activitySummariesUrl, after)
}

func (c *v3Client) GetActivity(activityId model.ActivityId) (*model.Activity, error) {
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

func (c *v3Client) GetActivities(activityIds []model.ActivityId) ([]*model.Activity, error) {
	var wg sync.WaitGroup

	activityIdChan := make(chan model.ActivityId, len(activityIds))
	activityChan := make(chan *model.Activity)

	// Goroutines to consume unfetched activities
	wg.Add(getActivitiesPoolSize)
	for i := 0; i < getActivitiesPoolSize; i++ {
		go func() {
			defer wg.Done()
			for activityId := range activityIdChan {
				if activity, err := c.GetActivity(activityId); err == nil {
					activityChan <- activity
				}
			}
		}()
	}

	for _, activityId := range activityIds {
		activityIdChan <- activityId
	}
	close(activityIdChan)

	// Goroutine to close activityChan once all fetching goroutines have completed
	go func() {
		wg.Wait()
		close(activityChan)
	}()

	// Index results by activityId
	activityMap := make(map[model.ActivityId]*model.Activity, len(activityIds))
	for activity := range activityChan {
		activityMap[activity.Id] = activity
	}

	// Reorder activities to match given activityIds
	activities := make([]*model.Activity, 0, len(activityMap))
	for _, activityId := range activityIds {
		activity, found := activityMap[activityId]
		if found {
			activities = append(activities, activity)
		}
	}

	return activities, nil
}

func (c *v3Client) GetRelatedActivitySummaries(activityId model.ActivityId) ([]*model.ActivitySummary, error) {
	return c.getActivitySummaries(relatedActivitySummariesUrl(activityId), Beginning)
}

func activityUrl(activityId model.ActivityId) string {
	return fmt.Sprintf("/activities/%d", activityId)
}

func relatedActivitySummariesUrl(activityId model.ActivityId) string {
	return fmt.Sprintf("%s/related", activityUrl(activityId))
}

func (c *v3Client) getActivitySummaries(url string, after model.ActivityId) ([]*model.ActivitySummary, error) {
	allActivities := make([]*model.ActivitySummary, 0)
	complete := false
	for page := 1; !complete; page++ {
		activities, err := c.getActivitySummariesPage(url, page)
		if err != nil {
			return nil, err
		}

		filtered := filterBefore(activities, after)
		allActivities = append(allActivities, filtered...)

		complete = len(activities) == 0 || len(activities) > len(filtered)
	}

	return reverse(allActivities), nil
}

func (c *v3Client) getActivitySummariesPage(url string, page int) ([]*model.ActivitySummary, error) {
	if page <= 0 {
		return nil, errors.New("page must be positive")
	}

	body, err := c.httpClient.Get(url, map[string]interface{}{"per_page": activitySummariesPageSize, "page": page})
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
