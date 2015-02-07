package models

import "strconv"

type Location struct {
	Longitude float64 `bson:"lng" json:"lng"`
	Latitude  float64 `bson:"lat" json:"lat"`
	Name      string  `bson:"name" json:"name"`
}

func (l Location) String() string {
	return l.Name
}

func ParseFloat64(input string) float64 {
	if input == "" {
		return 0
	}

	result, err := strconv.ParseFloat(input, 64)

	if err != nil {
		panic(err)
	}

	return result
}
