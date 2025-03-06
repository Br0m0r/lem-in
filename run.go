package main

import (
	"fmt"
	"os"
)

func Run() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <input_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// Parse the input file: retrieve ant count, room definitions, and tunnel definitions.
	antCount, rooms, tunnels, err := ParseInputFile(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Build the graph structure representing the ant farm.
	antFarmGraph, err := BuildGraph(rooms, tunnels)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Find the quickest path(s) from the start to the end room.
	paths, err := FindQuickestPaths(antFarmGraph)
	if err != nil {
		fmt.Println("ERROR: no path found")
		os.Exit(1)
	}

	// Echo the input data if required by the project specifications.
	PrintInputData(antCount, rooms, tunnels)

	// Simulate the ant movements along the found path(s).
	SimulateAntMovements(antCount, paths)
}
