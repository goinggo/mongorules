// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package data

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

//** NEW TYPES

// BuoyCondition contains information for an individual station
type BuoyCondition struct {
	WindSpeed     float64 `bson:"wind_speed_milehour"`
	WindDirection int     `bson:"wind_direction_degnorth"`
	WindGust      float64 `bson:"gust_wind_speed_milehour"`
}

// BuoyLocation contains the buoys location
type BuoyLocation struct {
	Type        string    `bson:"type"`
	Coordinates []float64 `bson:"coordinates"`
}

// BuoyStation contains information for an individual station
type BuoyStation struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	StationId string        `bson:"station_id"`
	Name      string        `bson:"name"`
	LocDesc   string        `bson:"location_desc"`
	Condition BuoyCondition `bson:"condition"`
	Location  BuoyLocation  `bson:"location"`
	Distance  float64
}

//** PUBLIC FUNCTIONS

// GetBuoyStation retrieves the specified station id
//  stationId: The station id to display
//  collection: The collection to find the buoy station
func GetBuoyStation(stationId string, collection *mgo.Collection) (buoyStation *BuoyStation, err error) {

	// Find all the buoys
	query := collection.Find(bson.M{"station_id": stationId})

	// Capture the specified buoy
	buoyStation = &BuoyStation{}
	err = query.One(buoyStation)

	return buoyStation, err
}
