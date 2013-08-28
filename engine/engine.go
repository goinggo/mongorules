// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package engine

import (
	"fmt"
	"github.com/goinggo/mongorules/rules"
	"labix.org/v2/mgo"
	"time"
)

//** CONSTANTS

const (
	MONGODB_HOST     = "ds035428.mongolab.com:35428"
	MONGODB_DATABASE = "goinggo"
	MONGODB_USERNAME = "guest"
	MONGODB_PASSWORD = "welcome"
)

//** INTERFACES

// Rule is implemented by rule objects
type Rule interface {
	Run()
}

//** PUBLIC FUNCTIONS

// RunRule will run the specified rule and display the outcome
func RunRule(ruleName string) {

	// Create MongoDB connectivity parameters
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{MONGODB_HOST},
		Timeout:  10 * time.Second,
		Database: MONGODB_DATABASE,
		Username: MONGODB_USERNAME,
		Password: MONGODB_PASSWORD,
	}

	// Connect to MongoDB and establish a connection
	// Only do this once in your application. There is a lot of overhead with this call.
	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {

		fmt.Printf("ERROR : %s", err)
		return
	}

	// Close the session when we are done
	defer func() {

		// Close the session to Mongo
		session.Close()
	}()

	// Capture a reference to the collection
	collection := session.DB(MONGODB_DATABASE).C("buoy_stations")

	// Reference to the rule to run
	var rule Rule

	// Create the specified rule object
	switch ruleName {

	case "tampa":

		rule = rules.NewTampaRule(collection)
		break

	default:

		fmt.Printf("Unknown Rules\n")
		return
	}

	// Run the rule
	rule.Run()
}
