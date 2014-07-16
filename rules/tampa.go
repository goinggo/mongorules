// Copyright 2013 Ardan Studios. All rights reserved.
// Use of tampa source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rules : tampa.go implements rules for tampa.
package rules

import (
	"fmt"
	"github.com/goinggo/mongorules/data"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const (
	// The radius of the Earth in miles.
	distanceMultiplier float64 = 3963.192
)

//** NEW TYPES

// Tampa is a rule to determine if we should go finishing in Tampa.
type Tampa struct {
	Collection      *mgo.Collection
	Latitude        float64
	Longitude       float64
	MaxDistance     float64
	MaxAvgWindSpeed float64
}

//** PUBLIC FUNCTIONS

// NewTampaRule creates the new rule.
func NewTampaRule(collection *mgo.Collection) *Tampa {
	// https://maps.google.com/maps?q=27.945886,-82.798676&z=10

	// We are going to look at buoys within a 30 mile radius
	// of our location in Clearwater Florida

	// The max average wind speed allowed is 15 miles/hour

	return &Tampa{
		Collection:      collection,
		Latitude:        27.945886,
		Longitude:       -82.798676,
		MaxDistance:     (30.0 / distanceMultiplier),
		MaxAvgWindSpeed: 15.0,
	}
}

// Run executes the tampa rule.
func (tampa *Tampa) Run() {
	// Calculate the average wind speed for all buoys in the Tampa area.
	avgWindSpeed, err := tampa._CalculateAverageWindSpeed()
	if err != nil {
		fmt.Printf("\nERROR : %s\n\n", err)
		return
	}

	// Check the average windspeed is within range.
	if avgWindSpeed > tampa.MaxAvgWindSpeed {
		fmt.Printf("\n*** Stay Home, Tampa Is Not Good : Average Wind Speed Is %.2f ***\n\n", avgWindSpeed)
		return
	}

	// Find the the buoy with the current lowest wind gust.
	buoyStationWindGust, err := tampa._FindLowestWindGust()
	if err != nil {
		fmt.Printf("\nERROR : %s\n\n", err)
		return
	}

	// Find the buoy closest to our current location.
	buoyStationDistance, err := tampa._FindClosestBuoy()
	if err != nil {
		fmt.Printf("\nERROR : %s\n\n", err)
		return
	}

	// Create the extra fields for the display results.
	extraFields := make(map[string]string)
	extraFields["Avg Wind Gust"] = fmt.Sprintf("%.2f Miles Per Hour", avgWindSpeed)

	tampa._DisplayResults("Tampa Buoy With Lowest Wind Gust", buoyStationWindGust, extraFields)
	tampa._DisplayResults("Tampa Buoy Closest To Your Location", buoyStationDistance, extraFields)
}

// Checks the average wind speed against all buoys in Tampa
func (tampa *Tampa) _CalculateAverageWindSpeed() (float64, error) {
	/*
		db.buoy_stations.aggregate(
		{"$geoNear": { "near": [-82.798676,27.945886], "query": {"condition.wind_speed_milehour" : {"$ne" : null}}, "distanceField": "distance", "maxDistance": 0.00756965597428, "spherical": true, "distanceMultiplier": 3963.192 }},
		{"$project" : { "station_id" : "$station_id", "wind_speed" : "$condition.wind_speed_milehour", "_id" : 0  }},
		{"$group" : { "_id" : 1, "total_stations" : {"$sum" : 1}, "average_wind_speed" : {"$avg" : "$wind_speed"}}}
		)
	*/

	operation1 := bson.M{
		"$geoNear": bson.M{
			"near": []float64{tampa.Longitude, tampa.Latitude},
			"query": bson.M{
				"condition.wind_speed_milehour": bson.M{"$ne": nil},
			},
			"distanceField":      "distance",
			"maxDistance":        tampa.MaxDistance,
			"spherical":          true,
			"distanceMultiplier": distanceMultiplier,
		},
	}

	operation2 := bson.M{
		"$project": bson.M{
			"station_id": "$station_id",
			"wind_speed": "$condition.wind_speed_milehour", "_id": 0,
		},
	}

	operation3 := bson.M{
		"$group": bson.M{
			"_id": 1,
			"average_wind_speed": bson.M{
				"$avg": "$wind_speed",
			},
		},
	}

	operations := []bson.M{operation1, operation2, operation3}

	// Prepare the query to run in the MongoDB aggregation pipeline.
	pipe := tampa.Collection.Pipe(operations)

	// Run the queries and capture the results
	var results []bson.M
	if err := pipe.All(&results); err != nil {
		return 0, err
	}

	return results[0]["average_wind_speed"].(float64), nil
}

// _FindLowestWindGust finds the tampa buoy with the lowest wind gust.
func (tampa *Tampa) _FindLowestWindGust() (*data.BuoyStation, error) {
	/*
		db.buoy_stations.aggregate(
		{"$geoNear": { "near": [-82.798676,27.945886], "query": {"condition.wind_speed_milehour" : {"$ne" : null}}, "distanceField": "distance", "maxDistance": 0.00756965597428, "spherical": true, "distanceMultiplier": 3963.192 }},
		{"$project" : { "station_id" : "$station_id", "gust_wind_speed" : "$condition.gust_wind_speed_milehour", "distance" : "$distance", "_id" : 0  }},
		{"$sort" : { "gust_wind_speed" : 1 }},
		{"$limit" : 1 }
		)
	*/

	operation1 := bson.M{
		"$geoNear": bson.M{
			"near": []float64{tampa.Longitude, tampa.Latitude},
			"query": bson.M{
				"condition.wind_speed_milehour": bson.M{"$ne": nil},
			},
			"distanceField":      "distance",
			"maxDistance":        tampa.MaxDistance,
			"spherical":          true,
			"distanceMultiplier": distanceMultiplier,
		},
	}

	operation2 := bson.M{
		"$project": bson.M{
			"station_id":      "$station_id",
			"gust_wind_speed": "$condition.gust_wind_speed_milehour",
			"distance":        "$distance",
			"_id":             0,
		},
	}

	operation3 := bson.M{
		"$sort": bson.M{
			"gust_wind_speed": 1,
		},
	}

	operation4 := bson.M{"$limit": 1}

	operations := []bson.M{operation1, operation2, operation3, operation4}

	// Prepare the operations to run in the MongoDB aggregation pipeline
	pipe := tampa.Collection.Pipe(operations)

	// Run the operations and capture the results
	var results []bson.M
	if err := pipe.All(&results); err != nil {
		return nil, err
	}

	stationID := results[0]["station_id"].(string)
	distance := results[0]["distance"].(float64)

	// Capture the buoy station
	buoyStation, err := data.GetBuoyStation(stationID, tampa.Collection)
	if err != nil {
		return nil, err
	}

	// Set the distance
	buoyStation.Distance = distance

	return buoyStation, nil
}

// _FindClosestBuoy finds the tampa buoy closest to the current location.
func (tampa *Tampa) _FindClosestBuoy() (*data.BuoyStation, error) {
	/*
		db.buoy_stations.aggregate(
		{"$geoNear": { "near": [-82.798676,27.945886], "query": {"condition.wind_speed_milehour" : {"$ne" : null}}, "distanceField": "distance", "maxDistance": 0.00756965597428, "spherical": true, "distanceMultiplier": 3963.192 }},
		{"$project" : { "station_id" : "$station_id", "distance" : "$distance", "_id" : 0  }},
		{"$sort" : { "distance" : 1 }},
		{"$limit" : 1 }
		)
	*/

	operation1 := bson.M{
		"$geoNear": bson.M{
			"near": []float64{tampa.Longitude, tampa.Latitude},
			"query": bson.M{
				"condition.wind_speed_milehour": bson.M{"$ne": nil},
			},
			"distanceField":      "distance",
			"maxDistance":        tampa.MaxDistance,
			"spherical":          true,
			"distanceMultiplier": distanceMultiplier,
		},
	}

	operation2 := bson.M{
		"$project": bson.M{
			"station_id": "$station_id",
			"distance":   "$distance",
			"_id":        0,
		},
	}

	operation3 := bson.M{
		"$sort": bson.M{
			"distance": 1,
		},
	}

	operation4 := bson.M{"$limit": 1}

	operations := []bson.M{operation1, operation2, operation3, operation4}

	// Prepare the operations to run in the MongoDB aggregation pipeline
	pipe := tampa.Collection.Pipe(operations)

	// Run the operations and capture the results
	var results []bson.M
	if err := pipe.All(&results); err != nil {
		return nil, err
	}

	stationID := results[0]["station_id"].(string)
	distance := results[0]["distance"].(float64)

	// Capture the buoy station
	buoyStation, err := data.GetBuoyStation(stationID, tampa.Collection)
	if err != nil {
		return nil, err
	}

	// Set the distance
	buoyStation.Distance = distance

	return buoyStation, err
}

// _DisplayResults provides the final information for a successful result.
func (tampa *Tampa) _DisplayResults(title string, buoyStation *data.BuoyStation, extraFields map[string]string) {
	fmt.Printf("\n%s\n", title)
	fmt.Printf("Station Id\t\t\t: %s\n", buoyStation.StationID)
	fmt.Printf("Name\t\t\t: %s\n", buoyStation.Name)
	fmt.Printf("Location\t\t\t: %s\n", buoyStation.LocDesc)
	fmt.Printf("Latitude\t\t\t: %f\n", buoyStation.Location.Coordinates[1])
	fmt.Printf("Logitude\t\t\t: %f\n", buoyStation.Location.Coordinates[0])
	fmt.Printf("Distance\t\t\t: %f Miles\n", buoyStation.Distance)
	fmt.Printf("Wind Speed\t\t\t: %.2f Miles/Hour\n", buoyStation.Condition.WindSpeed)
	fmt.Printf("Wind Direction\t\t: %d From True North\n", buoyStation.Condition.WindDirection)
	fmt.Printf("Wind Gust\t\t\t: %.2f Miles/Hour\n", buoyStation.Condition.WindGust)

	for field, value := range extraFields {
		fmt.Printf("%s\t\t: %s\n", field, value)
	}
}
