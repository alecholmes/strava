package model

import (
	"time"
)

type ActivitySummary struct {
	Id                 ActivityId `json:"id"`
	Name               string     `json:"name"`
	Athlete            *Athlete   `json:"athlete"`
	StartDate          time.Time  `json:"start_date"`
	StartDateLocal     time.Time  `json:"start_date_local"`
	Timezone           string     `json:"timezone"`
	MovingTime         uint32     `json:"moving_time"`          // Seconds
	ElapsedTime        uint32     `json:"elapsed_time"`         // Seconds
	Distance           float32    `json:"distance"`             // Meters
	TotalElevationGain float32    `json:"total_elevation_gain"` // Meters
	AverageSpeed       float32    `json:"average_speed"`        // Meters/sec
	MaxSpeed           float32    `json:"max_speed"`            // Meters/sec
}
