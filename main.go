package main

import "os"

var step = 2

func main() {

	// Register the record definitions with gob.
	gobRegisterRecords()

	// Capture data.
	err := capture(pathData, pathIndex)
	if err != nil {
		os.Exit(1)
	}

	// Analyze data.
	err = analysis(pathData, pathIndex)
	if err != nil {
		os.Exit(1)
	}
}
