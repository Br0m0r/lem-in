package main

import (
	"fmt"
	"os"
)

// Run is the entry point for the application.
// It parses the input file, validates data, builds the graph,
// finds paths, assigns ants, gathers extra info, and starts the simulation.
func Run() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run . <input_file>")
		os.Exit(1)
	}

	inputFile := os.Args[1]

	// Parse the input file to obtain ant count, room definitions, and tunnels.
	antCount, rooms, tunnels, err := ParseInputFile(inputFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Validate ant count.
	if antCount <= 0 {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}

	// Ensure that there is at least one start and one end room.
	var startFound, endFound bool
	for _, room := range rooms {
		if room.IsStart {
			startFound = true
		}
		if room.IsEnd {
			endFound = true
		}
	}
	if !startFound || !endFound {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}

	// Build the graph from rooms and tunnels.
	antFarmGraph, err := BuildGraph(rooms, tunnels)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Find multiple paths using the max-flow algorithm.
	paths, err := FindMultiplePaths(antFarmGraph)
	if err != nil || len(paths) == 0 {
		fmt.Println("ERROR: invalid data format")
		os.Exit(1)
	}

	// Distribute ants across the available paths.
	assignment := AssignAnts(antCount, paths)

	// Gather extra information (input details and summary) for the output file.
	extraInfo := PrintExtraInfo(antCount, rooms, tunnels, paths, assignment)

	// Run the simulation. The terminal shows only minimal move info,
	// while detailed output (including grid visualization) is written to simulation_output.txt.
	SimulateMultiPath(antCount, paths, assignment, extraInfo)
}
