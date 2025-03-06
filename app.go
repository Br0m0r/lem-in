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

	// Parse the input file.
	antCount, rooms, tunnels, err := ParseInputFile(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Build the graph.
	antFarmGraph, err := BuildGraph(rooms, tunnels)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Find multiple valid vertex-disjoint paths from start to end.
	paths, err := FindMultiplePaths(antFarmGraph)
	if err != nil || len(paths) == 0 {
		fmt.Println("ERROR: no valid paths found")
		os.Exit(1)
	}

	// Echo the input data.
	PrintInputData(antCount, rooms, tunnels)

	// Assign ants to available paths using a greedy algorithm with sorting.
	assignment := AssignAnts(antCount, paths)

	// Simulate ant movements concurrently on all paths.
	SimulateMultiPath(antCount, paths, assignment)
}
