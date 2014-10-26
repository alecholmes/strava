package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/alecholmes/strava/client"
	"github.com/alecholmes/strava/model"
)

func main() {
	accessTokenFlag := flag.String("accessToken", "", "access token; required")
	afterFlag := flag.Int("afterId", 0, "beginning activity id, exclusive")
	segmentsFlag := flag.Bool("segments", false, "print segment details")
	delimiterFlag := flag.String("delimiter", ",", "output field delimiter character")
	flag.Parse()

	if *accessTokenFlag == "" {
		flag.Usage()
		return
	}

	delimiter, size := utf8.DecodeRuneInString(*delimiterFlag)
	if size == 0 || len(*delimiterFlag) != size {
		fmt.Println("Delimiter can only be one character")
		flag.Usage()
		return
	}

	client := client.NewClient(*accessTokenFlag)

	activitySummaries, err := client.GetActivitySummaries(model.ActivityId(*afterFlag))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting activity summaries: %s", err)
	}

	if *segmentsFlag {
		activities, err := getActivities(client, activitySummaries)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting activities: %s", err)
		}
		printCsv(delimiter, segmentTuples(activities))
	} else {
		printCsv(delimiter, summaryTuples(activitySummaries))
	}
}

func getActivities(client client.Client, summaries []*model.ActivitySummary) ([]*model.Activity, error) {
	activityIds := make([]model.ActivityId, len(summaries))
	for i, summary := range summaries {
		activityIds[i] = summary.Id
	}

	return client.GetActivities(activityIds)
}

func printCsv(delimiter rune, tuples [][]string) {
	writer := csv.NewWriter(os.Stdout)
	writer.Comma = delimiter
	writer.WriteAll(tuples)
}

func summaryTuples(summaries []*model.ActivitySummary) [][]string {
	tuples := make([][]string, len(summaries))
	for i, summary := range summaries {
		tuple := make([]string, 8)
		tuple[0] = fmt.Sprintf("%d", summary.Id)
		tuple[1] = summary.Name
		tuple[2] = fmt.Sprintf("%.2f", summary.Distance)
		tuple[3] = fmt.Sprintf("%.2f", summary.AverageSpeed)
		tuple[4] = fmt.Sprintf("%.2f", summary.TotalElevationGain)
		tuple[5] = fmt.Sprintf("%d", summary.MovingTime)
		tuple[6] = fmt.Sprintf("%d", summary.ElapsedTime)
		tuple[7] = summary.StartDate.String()
		tuples[i] = tuple
	}
	return tuples
}

func segmentTuples(activities []*model.Activity) [][]string {
	tuples := make([][]string, 0)
	for _, activity := range activities {
		for _, effort := range activity.SegmentEfforts {
			tuple := make([]string, 10)
			tuple[0] = fmt.Sprintf("%d", activity.Id)
			tuple[1] = activity.Name
			tuple[2] = fmt.Sprintf("%d", effort.Id)
			tuple[3] = fmt.Sprintf("%d", effort.Segment.Id)
			tuple[4] = effort.Segment.Name
			tuple[5] = fmt.Sprintf("%.2f", effort.Segment.Distance)
			tuple[6] = fmt.Sprintf("%d", effort.Segment.ClimbCategory)
			tuple[7] = fmt.Sprintf("%d", effort.ElapsedTime)
			tuple[8] = fmt.Sprintf("%d", effort.PrRank)
			tuple[9] = effort.StartDate.String()
			tuples = append(tuples, tuple)
		}
	}
	return tuples
}
