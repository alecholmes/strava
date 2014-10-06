package model

type SegmentId int64

type Segment struct {
	Id            SegmentId `json:"id"`
	Name          string    `json:"name"`
	Distance      float32   `json:"distance"`       // Meters
	ElevationLow  float32   `json:"elevation_low"`  // Meters
	ElevationHigh float32   `json:"elevation_high"` // Meters
	AverageGrade  float32   `json:"average_grade"`
	MaximumGrade  float32   `json:"maximum_grade"`
	ClimbCategory uint8     `json:"climb_category"` // [0, 5], hard to easy
}
