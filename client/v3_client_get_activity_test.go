package client

import (
	"reflect"
	"testing"
	"time"

	"github.com/alecholmes/strava/model"
)

func TestGetActivity(t *testing.T) {
	client, rawClient := newTestClient()

	url := fullActivityUrl(model.ActivityId(203378452), rawClient, t)
	rawClient.Gets[url] = expectedBody([]byte(activityJson))

	activity, err := client.GetActivity(203378452)
	if err != nil {
		t.Errorf("Unexpected error for GetActivity: %s", err)
	}

	expectedSegment := model.Segment{
		Id:            7750436,
		Name:          "Graton Rd., Sullivan to Facendini",
		Distance:      6086.3,
		ElevationLow:  37.4,
		ElevationHigh: 205.2,
		AverageGrade:  2.4,
		MaximumGrade:  13.9,
		ClimbCategory: 0,
	}

	expectedSegmentEffort := model.SegmentEffort{
		Id:             4792121264,
		ElapsedTime:    877,
		StartDate:      time.Date(2014, 10, 4, 15, 38, 36, 0, time.UTC),
		StartDateLocal: time.Date(2014, 10, 4, 8, 38, 36, 0, time.UTC),
		PrRank:         1,
		KomRank:        0,
		Segment:        &expectedSegment,
	}

	expectedActivity := model.Activity{
		Id:                 203378452,
		Name:               "Gran Fondo",
		StartDate:          time.Date(2014, 10, 4, 15, 7, 31, 0, time.UTC),
		StartDateLocal:     time.Date(2014, 10, 4, 8, 7, 31, 0, time.UTC),
		Timezone:           "(GMT-08:00) America/Los_Angeles",
		MovingTime:         21921,
		ElapsedTime:        27446,
		Distance:           164918,
		TotalElevationGain: 2523.8,
		AverageSpeed:       7.523,
		MaxSpeed:           16.1,
		SegmentEfforts:     []*model.SegmentEffort{&expectedSegmentEffort},
	}
	if !reflect.DeepEqual(&expectedActivity, activity) {
		t.Errorf("Summaries were not the same. Expected %s but was %s", &expectedActivity, activity)
	}
}

// This should really be in a file
const activityJson = `
{
    "achievement_count": 84,
    "athlete": {
        "id": 471686,
        "resource_state": 1
    },
    "athlete_count": 4,
    "average_speed": 7.523,
    "average_watts": 164.0,
    "calories": 4007.9,
    "comment_count": 0,
    "commute": false,
    "description": "High 40s to 90s.",
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
    "kudos_count": 14,
    "location_city": "Santa Rosa",
    "location_country": "United States",
    "location_state": "California",
    "manual": false,
    "map": {
        "id": "a203378452",
        "polyline": "fake",
        "resource_state": 3,
        "summary_polyline": "fake"
    },
    "max_speed": 16.1,
    "moving_time": 21921,
    "name": "Gran Fondo",
    "photo_count": 0,
    "private": false,
    "resource_state": 3,
    "segment_efforts": [
        {
            "activity": {
                "id": 203378452
            },
            "athlete": {
                "id": 471686
            },
            "average_watts": 220.5,
            "distance": 6063.2,
            "elapsed_time": 877,
            "end_index": 406,
            "hidden": false,
            "id": 4792121264,
            "kom_rank": null,
            "moving_time": 877,
            "name": "Graton Rd., Sullivan to Facendini",
            "pr_rank": 1,
            "resource_state": 2,
            "segment": {
                "activity_type": "Ride",
                "average_grade": 2.4,
                "city": null,
                "climb_category": 0,
                "country": null,
                "distance": 6086.3,
                "elevation_high": 205.2,
                "elevation_low": 37.4,
                "end_latitude": 38.416422,
                "end_latlng": [
                    38.416422,
                    -122.934403
                ],
                "end_longitude": -122.934403,
                "id": 7750436,
                "maximum_grade": 13.9,
                "name": "Graton Rd., Sullivan to Facendini",
                "private": false,
                "resource_state": 2,
                "starred": false,
                "start_latitude": 38.435831,
                "start_latlng": [
                    38.435831,
                    -122.882139
                ],
                "start_longitude": -122.882139,
                "state": null
            },
            "start_date": "2014-10-04T15:38:36Z",
            "start_date_local": "2014-10-04T08:38:36Z",
            "start_index": 245
        }
    ],
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
}`
