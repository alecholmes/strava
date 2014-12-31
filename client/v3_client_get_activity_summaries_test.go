package client

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/alecholmes/strava/model"
)

func TestGetActivitySummaries(t *testing.T) {
	client, rawClient := newTestClient()

	rawClient.Gets[pageUrl(rawClient, 1, t)] = expectedBody([]byte(activitySummariesJson))
	rawClient.Gets[pageUrl(rawClient, 2, t)] = expectedBody([]byte("[]"))

	summaries, err := client.GetActivitySummaries(Beginning)
	if err != nil {
		t.Fatalf("Unexpected error for GetActivitySummaries. error=%s", err)
	}

	if len(summaries) != 2 {
		t.Fatalf("Expected 2 summaries but got %d", len(summaries))
	}

	expectedFirst := model.ActivitySummary{
		Id:                 model.ActivityId(202315892),
		Name:               "Headlands w MC",
		Athlete:            &model.Athlete{Id: 471686},
		StartDate:          time.Date(2014, 10, 2, 13, 12, 24, 0, time.UTC),
		StartDateLocal:     time.Date(2014, 10, 2, 6, 12, 24, 0, time.UTC),
		Timezone:           "(GMT-08:00) America/Los_Angeles",
		MovingTime:         6472,
		ElapsedTime:        7752,
		Distance:           44053.2,
		TotalElevationGain: 796.3,
		AverageSpeed:       6.807,
		MaxSpeed:           14.7,
	}
	if !reflect.DeepEqual(&expectedFirst, summaries[0]) {
		t.Fatalf("Summaries were not the same. expected=%s, actual=%s", &expectedFirst, summaries[0])
	}

	expectedSecond := model.ActivitySummary{
		Id:                 model.ActivityId(203378452),
		Name:               "Gran Fondo",
		Athlete:            &model.Athlete{Id: 471686},
		StartDate:          time.Date(2014, 10, 4, 15, 7, 31, 0, time.UTC),
		StartDateLocal:     time.Date(2014, 10, 4, 8, 7, 31, 0, time.UTC),
		Timezone:           "(GMT-08:00) America/Los_Angeles",
		MovingTime:         21921,
		ElapsedTime:        27446,
		Distance:           164918,
		TotalElevationGain: 2523.8,
		AverageSpeed:       7.523,
		MaxSpeed:           16.1,
	}
	if !reflect.DeepEqual(&expectedSecond, summaries[1]) {
		t.Fatalf("Summaries were not the same. expected=%s, actual=%s", &expectedSecond, summaries[1])
	}
}

func TestGetActivitySummaries_Pages(t *testing.T) {
	client, rawClient := newTestClient()

	first := model.ActivitySummary{Id: model.ActivityId(11)}
	second := model.ActivitySummary{Id: model.ActivityId(22)}

	firstJson := fmt.Sprintf("[%s]", toJson(first, t))
	secondJson := fmt.Sprintf("[%s]", toJson(second, t))

	rawClient.Gets[pageUrl(rawClient, 1, t)] = expectedBody([]byte(firstJson))
	rawClient.Gets[pageUrl(rawClient, 2, t)] = expectedBody([]byte(secondJson))
	rawClient.Gets[pageUrl(rawClient, 3, t)] = expectedBody([]byte("[]"))

	summaries, err := client.GetActivitySummaries(Beginning)
	if err != nil {
		t.Fatalf("Unexpected error for GetActivitySummaries. error=%s", err)
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

func TestGetActivitySummaries_After(t *testing.T) {
	client, rawClient := newTestClient()

	activity11 := model.ActivitySummary{Id: model.ActivityId(11)}
	activity22 := model.ActivitySummary{Id: model.ActivityId(22)}
	activity33 := model.ActivitySummary{Id: model.ActivityId(33)}

	firstJson := fmt.Sprintf("[%s]", toJson(activity33, t))
	secondJson := fmt.Sprintf("[%s, %s]", toJson(activity22, t), toJson(activity11, t))

	// The client should never attempt to get page 3 (which would otherwise return [])
	rawClient.Gets[pageUrl(rawClient, 1, t)] = expectedBody([]byte(firstJson))
	rawClient.Gets[pageUrl(rawClient, 2, t)] = expectedBody([]byte(secondJson))

	summaries, err := client.GetActivitySummaries(activity11.Id)
	if err != nil {
		t.Fatalf("Unexpected error for GetActivitySummaries. error=%s", err)
	}

	if len(summaries) != 2 {
		t.Fatalf("Expected 2 summaries but got %d", len(summaries))
	}

	// can't do use DeepEqual since times are unmarshalling in different location
	if activity22.Id != summaries[0].Id {
		t.Fatalf("Summaries were not the same. expected=%s, actual=%s", &activity22, summaries[0])
	}
	if activity33.Id != summaries[1].Id {
		t.Fatalf("Summaries were not the same. expected=%s, actual=%s", &activity33, summaries[1])
	}
}

func pageUrl(rawClient HttpClient, page uint32, t *testing.T) string {
	url, err := rawClient.AbsoluteUrl(activitySummariesUrl,
		map[string]interface{}{"per_page": activitySummariesPageSize, "page": page})
	if err != nil {
		t.Fatalf("Error creating test URL. error=%s", err)
	}
	return url
}

// This should really be in a file
const activitySummariesJson = `[
    {
        "achievement_count": 84,
        "athlete": {
            "id": 471686,
            "resource_state": 1
        },
        "athlete_count": 4,
        "average_speed": 7.523,
        "average_watts": 164.0,
        "comment_count": 0,
        "commute": false,
        "device_watts": false,
        "distance": 164918.0,
        "elapsed_time": 27446,
        "end_latlng": [
            38.44,
            -122.75
        ],
        "external_id": "2014-10-04-15-07-30.tcx",
        "flagged": false,
        "gear_id": null,
        "has_kudoed": false,
        "id": 203378452,
        "kilojoules": 3594.5,
        "kudos_count": 15,
        "location_city": "Santa Rosa",
        "location_country": "United States",
        "location_state": "California",
        "manual": false,
        "map": {
            "id": "a203378452",
            "resource_state": 2,
            "summary_polyline": "fake"
        },
        "max_speed": 16.1,
        "moving_time": 21921,
        "name": "Gran Fondo",
        "photo_count": 0,
        "private": false,
        "resource_state": 2,
        "start_date": "2014-10-04T15:07:31Z",
        "start_date_local": "2014-10-04T08:07:31Z",
        "start_latitude": 38.44,
        "start_latlng": [
            38.44,
            -122.75
        ],
        "start_longitude": -122.75,
        "timezone": "(GMT-08:00) America/Los_Angeles",
        "total_elevation_gain": 2523.8,
        "trainer": false,
        "truncated": null,
        "type": "Ride",
        "upload_id": 226294137
    },
    {
        "achievement_count": 11,
        "athlete": {
            "id": 471686,
            "resource_state": 1
        },
        "athlete_count": 9,
        "average_speed": 6.807,
        "average_watts": 139.6,
        "comment_count": 4,
        "commute": false,
        "device_watts": false,
        "distance": 44053.2,
        "elapsed_time": 7752,
        "end_latlng": [
            37.77,
            -122.43
        ],
        "external_id": "2014-10-02-13-12-23.tcx",
        "flagged": false,
        "gear_id": null,
        "has_kudoed": false,
        "id": 202315892,
        "kilojoules": 903.3,
        "kudos_count": 13,
        "location_city": "San Francisco",
        "location_country": "United States",
        "location_state": "CA",
        "manual": false,
        "map": {
            "id": "a202315892",
            "resource_state": 2,
            "summary_polyline": "fake"
        },
        "max_speed": 14.7,
        "moving_time": 6472,
        "name": "Headlands w MC",
        "photo_count": 0,
        "private": false,
        "resource_state": 2,
        "start_date": "2014-10-02T13:12:24Z",
        "start_date_local": "2014-10-02T06:12:24Z",
        "start_latitude": 37.77,
        "start_latlng": [
            37.77,
            -122.43
        ],
        "start_longitude": -122.43,
        "timezone": "(GMT-08:00) America/Los_Angeles",
        "total_elevation_gain": 796.3,
        "trainer": false,
        "truncated": 5,
        "type": "Ride",
        "upload_id": 225048181
    }
]
`
