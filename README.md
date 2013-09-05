# Mongo Rules

Copyright 2013 Ardan Studios. All rights reserved.  
Use of this source code is governed by a BSD-style license that can be found in the LICENSE handle.

The program shows how to use Mongo to analyze data in a Go program. The program uses the aggregation framework to build rules. The program has implemented one rule that identifies and analyzes certain buoy stations to determine if we should go finishing in Tampa, FL. The analysis is very basic but it will give you an idea of how I use Mongo to analyze data and build rules.

Ardan Studios  
12973 SW 112 ST, Suite 153  
Miami, FL 33186  
bill@ardanstudios.com

GoingGo.net Post:  
http://www.goinggo.net/2013/07/analyze-data-with-mongodb-and-go.html

	-- Get, build and install the code
	export GOPATH=$HOME/goinggo
	go get github.com/goinggo/mongorules
	
	-- Run the code
	cd $GOPATH/bin
	./mongorules

