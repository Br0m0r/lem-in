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

    // Basic validation: ant count must be > 0.
    if antCount <= 0 {
        fmt.Println("ERROR: invalid data format")
        os.Exit(1)
    }

    // Check that at least one start and one end room exist.
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

    // Build the graph.
    antFarmGraph, err := BuildGraph(rooms, tunnels)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // Find multiple valid vertex-disjoint paths from start to end.
    paths, err := FindMultiplePaths(antFarmGraph)
    if err != nil || len(paths) == 0 {
        fmt.Println("ERROR: invalid data format")
        os.Exit(1)
    }

    // Echo the input data (the original lines from the file).
    PrintInputData(antCount, rooms, tunnels)

    // Assign ants to available paths using a greedy algorithm with sorting.
    assignment := AssignAnts(antCount, paths)

    // Call PrintSummary (from your visualizer code) to display a textual overview in the terminal.
    // Make sure you have defined PrintSummary somewhere in your project (e.g., in visualizer.go).
    PrintSummary(antCount, rooms, tunnels, paths, assignment.Paths)

    // Simulate ant movements concurrently on all paths.
    SimulateMultiPath(antCount, paths, assignment)
}

