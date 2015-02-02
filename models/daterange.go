package models

import "time"

type DateRange struct {
	Start time.Time `bson:"start" json:"start"`
	End   time.Time `bson:"end" json:"end"`
}

func ParseDate(input string) time.Time {
	if input == "" {
		return time.Time{}
	}

	time, err := time.Parse(time.RFC3339, input)

	if err != nil {
		panic(err)
	}

	return time
}
