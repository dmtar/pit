package models

import (
	"fmt"
	"strconv"

	"gopkg.in/mgo.v2/bson"
)

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

// GetBSON is actually implementing the bson.Getter
// interface so this can be Marshaled into bson.
func (l Location) GetBSON() (interface{}, error) {
	res := bson.M{
		"name":        l.Name,
		"type":        "Point",
		"coordinates": []float64{l.Longitude, l.Latitude},
	}
	return res, nil
}

// SetBSON is implementing the bson.Setter interface
// so this can be Unmarshaled from bson.
// To implement the Setter interface we need to import
// the bson package, so we can have the bson.Raw type available.
// This one may panic with index out of range. We couldn't care less
// for that at the moment.
func (l *Location) SetBSON(raw bson.Raw) error {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("\n\n!!!!! PANIC !!!!!! %s ! \n\n", r)
		}
	}()

	response := struct {
		Name        string    `bson:"name"`
		Type        string    `bson:"type"`
		Coordinates []float64 `bson:"coordinates"`
	}{}

	err := raw.Unmarshal(&response)

	if err != nil {
		return err
	}

	l.Name = response.Name
	l.Longitude = response.Coordinates[0]
	l.Latitude = response.Coordinates[1]

	return nil
}
