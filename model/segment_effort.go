package model

import (
	"time"
)

type SegmentEffortId int64

type SegmentEffort struct {
	Id             SegmentEffortId `json:"id"`
	ElapsedTime    uint32          `json:"elapsed_time"`
	StartDate      time.Time       `json:"start_date"`
	StartDateLocal time.Time       `json:"start_date_local"`
	PrRank         uint32          `json:"pr_rank"`
	KomRank        uint32          `json:"kom_rank"`
	Segment        *Segment        `json:"segment"`
}
