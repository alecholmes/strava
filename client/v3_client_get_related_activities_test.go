package client

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/alecholmes/strava/model"
)

func TestGetRelatedActivitySummaries(t *testing.T) {
	client, rawClient := newTestClient()

	rawClient.Gets[relatedPageUrl(rawClient, 1, t)] = expectedBody([]byte(relatedActivitySummariesJson))
	rawClient.Gets[relatedPageUrl(rawClient, 2, t)] = expectedBody([]byte("[]"))

	summaries, err := client.GetRelatedActivitySummaries(model.ActivityId(123))
	if err != nil {
		t.Fatalf("Unexpected error for GetRelatedActivitySummaries. error=%s", err)
	}

	if len(summaries) != 2 {
		t.Fatalf("Expected 2 summaries but got %d", len(summaries))
	}

	expectedFirst := model.ActivitySummary{
		Id:                 model.ActivityId(9837863),
		Name:               "Levi's Gran Fondo - Missed 16k because of mechanical issues but got all the climbs in",
		Athlete:            &model.Athlete{Id: 699515},
		StartDate:          time.Date(2014, 10, 4, 15, 7, 17, 0, time.UTC),
		StartDateLocal:     time.Date(2014, 10, 4, 8, 7, 17, 0, time.UTC),
		Timezone:           "(GMT-08:00) America/Los_Angeles",
		MovingTime:         20631,
		ElapsedTime:        27557,
		Distance:           146974,
		TotalElevationGain: 2240.0,
		AverageSpeed:       7.124,
		MaxSpeed:           17.8,
	}
	if !reflect.DeepEqual(&expectedFirst, summaries[0]) {
		t.Fatalf("Summaries were not the same. expected=%s, actual=%s", &expectedFirst, summaries[0])
	}

	expectedSecond := model.ActivitySummary{
		Id:                 model.ActivityId(29823897),
		Name:               "Levi's Gran Fondo with heat wave",
		Athlete:            &model.Athlete{Id: 11235813},
		StartDate:          time.Date(2014, 10, 4, 15, 9, 5, 0, time.UTC),
		StartDateLocal:     time.Date(2014, 10, 4, 8, 9, 5, 0, time.UTC),
		Timezone:           "(GMT-08:00) America/Los_Angeles",
		MovingTime:         21085,
		ElapsedTime:        24578,
		Distance:           165652.0,
		TotalElevationGain: 2277.0,
		AverageSpeed:       7.856,
		MaxSpeed:           16.9,
	}
	if !reflect.DeepEqual(&expectedSecond, summaries[1]) {
		t.Fatalf("Summaries were not the same. expected=%s, actual=%s", &expectedSecond, summaries[1])
	}
}

func TestGetRelatedActivitySummaries_Pages(t *testing.T) {
	client, rawClient := newTestClient()

	first := model.ActivitySummary{Id: model.ActivityId(11)}
	second := model.ActivitySummary{Id: model.ActivityId(22)}

	firstJson := fmt.Sprintf("[%s]", toJson(first, t))
	secondJson := fmt.Sprintf("[%s]", toJson(second, t))

	rawClient.Gets[relatedPageUrl(rawClient, 1, t)] = expectedBody([]byte(firstJson))
	rawClient.Gets[relatedPageUrl(rawClient, 2, t)] = expectedBody([]byte(secondJson))
	rawClient.Gets[relatedPageUrl(rawClient, 3, t)] = expectedBody([]byte("[]"))

	summaries, err := client.GetRelatedActivitySummaries(model.ActivityId(123))
	if err != nil {
		t.Fatalf("Unexpected error for GetRelatedActivitySummaries. error=%s", err)
	}

	if len(summaries) != 2 {
		t.Fatalf("Expected 2 summaries but got %d", len(summaries))
	}

	// can't do use DeepEqual since times are unmarshalling in different location
	if second.Id != summaries[0].Id {
		t.Fatalf("Summaries were not the same. expected=%s, actual=%s", &second, summaries[0])
	}
	if first.Id != summaries[1].Id {
		t.Fatalf("Summaries were not the same. expected=%s, actual=%s", &first, summaries[1])
	}
}

func relatedPageUrl(rawClient HttpClient, page uint32, t *testing.T) string {
	url, err := rawClient.AbsoluteUrl(relatedActivitySummariesUrl(model.ActivityId(123)),
		map[string]interface{}{"per_page": activitySummariesPageSize, "page": page})
	if err != nil {
		t.Fatalf("Error creating test URL. error=%s", err)
	}
	return url
}

// This should really be in a file
const relatedActivitySummariesJson = `[
    {
        "achievement_count": 79,
        "athlete": {
            "badge_type_id": 1,
            "city": "San Francisco",
            "country": "United States",
            "created_at": "2012-11-14T07:08:28Z",
            "firstname": "Bea",
            "follower": "accepted",
            "friend": "accepted",
            "id": 11235813,
            "lastname": "Arthur",
            "premium": true,
            "resource_state": 2,
            "sex": "M",
            "state": "CA",
            "updated_at": "2015-01-11T15:04:29Z"
        },
        "athlete_count": 6,
        "average_cadence": 85.2,
        "average_heartrate": 140.7,
        "average_speed": 7.856,
        "average_temp": 21.0,
        "average_watts": 194.0,
        "comment_count": 2,
        "commute": false,
        "device_watts": false,
        "distance": 165652.0,
        "elapsed_time": 24578,
        "end_latlng": [
            38.44,
            -122.75
        ],
        "external_id": "tap-sync-f5b6b0363fb403ddbb801cabbf9d6d18-15113-53278c25cbe97e7fb1920db9.fit",
        "flagged": false,
        "gear_id": "b1083842",
        "has_kudoed": true,
        "id": 29823897,
        "kilojoules": 4091.1,
        "kudos_count": 27,
        "location_city": "Santa Rosa",
        "location_country": "United States",
        "location_state": "California",
        "manual": false,
        "map": {
            "id": "a203335389",
            "resource_state": 2,
            "summary_polyline": "blah"
        },
        "max_heartrate": 174.0,
        "max_speed": 16.9,
        "moving_time": 21085,
        "name": "Levi's Gran Fondo with heat wave",
        "photo_count": 0,
        "private": false,
        "resource_state": 2,
        "start_date": "2014-10-04T15:09:05Z",
        "start_date_local": "2014-10-04T08:09:05Z",
        "start_latitude": 38.44,
        "start_latlng": [
            38.44,
            -122.75
        ],
        "start_longitude": -122.75,
        "timezone": "(GMT-08:00) America/Los_Angeles",
        "total_elevation_gain": 2277.0,
        "trainer": false,
        "type": "Ride",
        "upload_id": 226240084
    },
    {
        "achievement_count": 56,
        "athlete": {
            "badge_type_id": 1,
            "city": "Vancouver",
            "country": "Canada",
            "created_at": "2012-06-29T15:07:05Z",
            "firstname": "Some",
            "follower": null,
            "friend": null,
            "id": 699515,
            "lastname": "Dude",
            "premium": true,
            "resource_state": 2,
            "sex": "M",
            "state": "BC",
            "updated_at": "2014-12-31T00:49:15Z"
        },
        "athlete_count": 4,
        "average_cadence": 75.8,
        "average_heartrate": 146.3,
        "average_speed": 7.124,
        "average_temp": 20.0,
        "average_watts": 145.9,
        "comment_count": 2,
        "commute": false,
        "device_watts": true,
        "distance": 146974.0,
        "elapsed_time": 27557,
        "end_latlng": [
            38.44,
            -122.75
        ],
        "external_id": "2014-10-04-08-07-17.fit",
        "flagged": false,
        "gear_id": "b616042",
        "has_kudoed": false,
        "id": 9837863,
        "kilojoules": 3009.7,
        "kudos_count": 7,
        "location_city": "Santa Rosa",
        "location_country": "United States",
        "location_state": "California",
        "manual": false,
        "map": {
            "id": "a203353614",
            "resource_state": 2,
            "summary_polyline": "xyz"
        },
        "max_heartrate": 175.0,
        "max_speed": 17.8,
        "moving_time": 20631,
        "name": "Levi's Gran Fondo - Missed 16k because of mechanical issues but got all the climbs in",
        "photo_count": 0,
        "private": false,
        "resource_state": 2,
        "start_date": "2014-10-04T15:07:17Z",
        "start_date_local": "2014-10-04T08:07:17Z",
        "start_latitude": 38.44,
        "start_latlng": [
            38.44,
            -122.75
        ],
        "start_longitude": -122.75,
        "timezone": "(GMT-08:00) America/Los_Angeles",
        "total_elevation_gain": 2240.0,
        "trainer": false,
        "type": "Ride",
        "upload_id": 226261835,
        "weighted_average_watts": 184
    }
]`
