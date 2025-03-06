package main

import (
	"fmt"
	"strings"
)

// SimulateMultiPath moves ants along multiple paths concurrently based on the assignment.
// It prints the moves turn by turn, ensuring that no intermediate room (other than start/end)
// is occupied by more than one ant per turn.
func SimulateMultiPath(antCount int, paths [][]string, assignment PathAssignment) {
	// Each simulation state represents a path's state.
	// For each path, we track:
	// - The path (slice of room names)
	// - Positions for each ant on that path (initially -1 means not injected)
	// - Global ant IDs assigned to that path.
	type pathSim struct {
		Path      []string
		Positions []int // position index for each ant in this path
		AntIDs    []int // global ant numbers for this path
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
	for {
		turnMoves := []string{}
		allDone := true // assume all ants are done unless we find otherwise

		// Process each simulation path.
		for _, sim := range sims {
			pathLen := len(sim.Path)
			// We'll update positions in a copy to simulate synchronous moves.
			newPos := make([]int, len(sim.Positions))
			copy(newPos, sim.Positions)

			// Process ants in reverse order (those closer to the end first).
			for j := len(sim.Positions) - 1; j >= 0; j-- {
				if sim.Positions[j] == -1 {
					// Injection: if ant hasn't been injected and room at index 1 is free.
					if !isOccupied(newPos, 1) {
						newPos[j] = 1
						turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[1]))
					}
				} else if sim.Positions[j] < pathLen-1 {
					next := sim.Positions[j] + 1
					// Move forward if next room is free or it's the end.
					if next == pathLen-1 || !isOccupied(newPos, next) {
						newPos[j] = next
						turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[next]))
					}
				}
			}

			// Update positions for this simulation.
			copy(sim.Positions, newPos)

			// Check if all ants on this path have reached the end.
			for _, pos := range sim.Positions {
				if pos != pathLen-1 {
					allDone = false
					break
				}
			}
		}

		// If no ant moved this turn, the simulation is complete.
		if len(turnMoves) == 0 {
			break
		}

		// Increment turn count only if moves were made.
		turn++
		fmt.Printf("Turn %d: %s\n", turn, strings.Join(turnMoves, " "))

		// If all ants reached the end, break immediately.
		if allDone {
			break
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
