package main

import (
	"fmt"
	"strings"
)

// SimulateMultiPath moves ants along multiple paths concurrently based on the assignment.
// It prints the moves turn by turn, ensuring each intermediate room holds only one ant.
func SimulateMultiPath(antCount int, paths [][]string, assignment PathAssignment) {
	// For each path, set up a simulation state.
	// Each simulation state tracks:
	// - The path (slice of room names)
	// - Positions for each ant assigned to that path (initially -1 means not injected)
	// - The ant IDs assigned to that path.
	type pathSim struct {
		Path      []string
		Positions []int // position index on the path for each ant assigned to this path
		AntIDs    []int // the global ant numbers assigned to this path
	}
	sims := make([]pathSim, len(paths))
	antCounter := 1
	for i, p := range paths {
		count := assignment.AntsPerPath[i]
		positions := make([]int, count)
		for j := range positions {
			positions[j] = -1
		}
		antIDs := make([]int, count)
		for j := 0; j < count; j++ {
			antIDs[j] = antCounter
			antCounter++
		}
		sims[i] = pathSim{Path: p, Positions: positions, AntIDs: antIDs}
	}

	turn := 0
	done := false
	for !done {
		turn++
		turnMoves := []string{}
		done = true // assume all paths are finished; will be set false if any path still has moving ants

		// Process each path simulation.
		for _, sim := range sims {
			pathLen := len(sim.Path)
			// newPositions will hold the updated positions for this turn.
			newPos := make([]int, len(sim.Positions))
			copy(newPos, sim.Positions)

			// Process ants in order for this path.
			for j := 0; j < len(sim.Positions); j++ {
				// If ant j has not yet been injected.
				if sim.Positions[j] == -1 {
					// Check newPos: if the first intermediate room (index 1) is free,
					// inject this ant.
					if !isOccupied(newPos, 1) {
						newPos[j] = 1
						turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[1]))
					}
				} else if sim.Positions[j] < pathLen-1 { // The ant is on the path but not finished.
					next := sim.Positions[j] + 1
					// For movement, check if the next room is free in newPos.
					if next == pathLen-1 || !isOccupied(newPos, next) {
						newPos[j] = next
						turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[next]))
					}
				}
			}
			// Update the simulation state for this path.
			copy(sim.Positions, newPos)
			// Check if any ant on this path is still not finished.
			for _, pos := range sim.Positions {
				if pos != pathLen-1 {
					done = false
					break
				}
			}
		}
		if len(turnMoves) > 0 {
			fmt.Printf("Turn %d: %s\n", turn, strings.Join(turnMoves, " "))
		}
	}
	fmt.Printf("Total turns: %d\n", turn)
}

// isOccupied returns true if any ant in positions is at the given position.
func isOccupied(positions []int, pos int) bool {
	for _, p := range positions {
		if p == pos {
			return true
		}
	}
	return false
}
