package main

import (
	"fmt"
	"os"
	"strings"
)

// GeneratePathGrid returns a 2D grid visualization for a single simulation path.
// Each room is represented as a box; if an ant is present, its label (e.g., L4) is shown.
func GeneratePathGrid(sim PathSim) string {
	var sb strings.Builder
	for i, room := range sim.Path {
		antLabel := ""
		// Check if an ant is currently in this room.
		for j, pos := range sim.Positions {
			if pos == i {
				antLabel = fmt.Sprintf("(L%d)", sim.AntIDs[j])
				break
			}
		}
		// Build the room box with or without the ant label.
		if antLabel != "" {
			sb.WriteString(fmt.Sprintf("[ %s %s ]", room, antLabel))
		} else {
			sb.WriteString(fmt.Sprintf("[ %s ]", room))
		}
		// Add an arrow between rooms if not the last room.
		if i < len(sim.Path)-1 {
			sb.WriteString(" ---> ")
		}
	}
	return sb.String()
}

// SimulateMultiPath simulates ant movement along each path concurrently.
// It prints minimal turn move info to the terminal and writes a detailed
// simulation (with extra info and full grid visualization) to simulation_output.txt.
func SimulateMultiPath(antCount int, paths [][]string, assignment PathAssignment, extraInfo string) {
	// Build simulation state for each path.
	sims := make([]PathSim, len(paths))
	antCounter := 1
	for i, p := range paths {
		count := assignment.AntsPerPath[i]
		positions := make([]int, count)
		// Initialize all ant positions to -1 (not yet injected).
		for j := range positions {
			positions[j] = -1
		}
		antIDs := make([]int, count)
		// Assign unique global IDs to ants.
		for j := 0; j < count; j++ {
			antIDs[j] = antCounter
			antCounter++
		}
		sims[i] = PathSim{Path: p, Positions: positions, AntIDs: antIDs}
	}

	// Create (or overwrite) the simulation output file.
	outFile, err := os.Create("simulation_output.txt")
	if err != nil {
		fmt.Println("Error creating simulation_output.txt:", err)
		return
	}
	defer outFile.Close()

	var terminalBuilder strings.Builder // For minimal terminal output.
	var fileBuilder strings.Builder     // For detailed output in file.

	// Write extra information (input, summary, etc.) at the top of the output file.
	fileBuilder.WriteString(extraInfo)
	fileBuilder.WriteString("\n\n")

	turn := 0
	for {
		var turnMoves []string
		fileBuilder.WriteString(fmt.Sprintf("TURN %d\n", turn+1))

		// Process each path's simulation.
		for idx, sim := range sims {
			pathLen := len(sim.Path)
			newPos := make([]int, len(sim.Positions))
			copy(newPos, sim.Positions)

			// Process ants in reverse order (ants closer to the end move first).
			for j := len(sim.Positions) - 1; j >= 0; j-- {
				if sim.Positions[j] == -1 {
					// Inject the ant into the path if the first room (index 1) is free.
					if !isOccupied(newPos, 1) {
						newPos[j] = 1
						turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[1]))
					}
				} else if sim.Positions[j] < pathLen-1 {
					next := sim.Positions[j] + 1
					// Move the ant forward if the next room is free or it's the end.
					if next == pathLen-1 || !isOccupied(newPos, next) {
						newPos[j] = next
						turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[next]))
					}
				}
			}
			// Update simulation state for the path.
			copy(sims[idx].Positions, newPos)

			// Append the grid visualization for this path to the file output.
			gridStr := GeneratePathGrid(sim)
			fileBuilder.WriteString(gridStr + "\n")
		}

		// End simulation if no ant has moved this turn.
		if len(turnMoves) == 0 {
			break
		}

		// Append minimal turn moves to terminal output.
		terminalBuilder.WriteString(fmt.Sprintf("Turn %d: %s\n", turn+1, strings.Join(turnMoves, " ")))
		turn++
		fileBuilder.WriteString("\n")
	}

	fileBuilder.WriteString(fmt.Sprintf("Total turns: %d\n", turn))

	// Write the detailed simulation output to the file.
	_, err = outFile.WriteString(fileBuilder.String())
	if err != nil {
		fmt.Println("Error writing to simulation_output.txt:", err)
	}

	// Print minimal move info to the terminal.
	fmt.Print(terminalBuilder.String())
	fmt.Println("2D grid visualization written to simulation_output.txt")
}

// isOccupied returns true if any ant is at the specified room index.
func isOccupied(positions []int, pos int) bool {
	for _, p := range positions {
		if p == pos {
			return true
		}
	}
	return false
}
