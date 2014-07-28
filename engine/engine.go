// Copyright 2013 Ardan Studios. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package engine processes the rules.
package engine

import (
	"fmt"
	"time"

	"github.com/goinggo/mongorules/rules"
	"gopkg.in/mgo.v2"
)

const (
	mongoHost     = "ds035428.mongolab.com:35428"
	mongoDatabase = "goinggo"
	mongoUser     = "guest"
	mongoPassowrd = "welcome"
)

// Rule is implemented by rule objects.
type Rule interface {
	Run()
}

// RunRule will run the specified rule and display the outcome.
func RunRule(ruleName string) {
	// Create MongoDB connectivity parameters.
	dialInfo := mgo.DialInfo{
		Addrs:    []string{mongoHost},
		Timeout:  10 * time.Second,
		Database: mongoDatabase,
		Username: mongoUser,
		Password: mongoPassowrd,
	}

	// Connect to MongoDB and establish a connection.
	// Only do this once in your application. There is a lot of overhead with this call.
	session, err := mgo.DialWithInfo(&dialInfo)
	if err != nil {
		fmt.Printf("ERROR : %s", err)
		return
	}

	// Close the session when we are done.
	defer session.Close()

	// Capture a reference to the collection.
	collection := session.DB(mongoDatabase).C("buoy_stations")

	// Reference to the rule to run.
	var rule Rule

	// Create the specified rule object.
	switch ruleName {
	case "tampa":
		rule = rules.NewTampaRule(collection)
		break
	default:
		fmt.Printf("Unknown Rules\n")
		return
	}

	// Run the rule.
	rule.Run()
}
