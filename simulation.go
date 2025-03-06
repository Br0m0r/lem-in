package main

import (
	"fmt"
	"strings"
)

// SimulateMultiPath moves ants along multiple paths concurrently based on the assignment.
// It prints the moves turn by turn, ensuring that no intermediate room is occupied by more than one ant per turn.
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
	done := false

	// Special case: if all paths are direct (length == 2), inject all ants immediately.
	allDirect := true
	for _, sim := range sims {
		if len(sim.Path) != 2 {
			allDirect = false
			break
		}
	}
	if allDirect {
		turn = 1
		turnMoves := []string{}
		for _, sim := range sims {
			for j, id := range sim.AntIDs {
				sim.Positions[j] = 1
				turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", id, sim.Path[1]))
			}
		}
		fmt.Printf("Turn %d: %s\n", turn, strings.Join(turnMoves, " "))
		fmt.Printf("Total turns: %d\n", turn)
		return
	}

	// Synchronous updating for paths with intermediate rooms.
	for !done {
		turn++
		turnMoves := []string{}
		done = true // assume finished; will set false if any path isn't complete

		for _, sim := range sims {
			pathLen := len(sim.Path)
			newPos := make([]int, len(sim.Positions))
			copy(newPos, sim.Positions)

			for j := 0; j < len(sim.Positions); j++ {
				if sim.Positions[j] == -1 {
					// Injection: if not injected and first intermediate room (index 1) is free in newPos.
					if !isOccupied(newPos, 1) {
						newPos[j] = 1
						turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[1]))
					}
				} else if sim.Positions[j] < pathLen-1 {
					next := sim.Positions[j] + 1
					if next == pathLen-1 || !isOccupied(newPos, next) {
						newPos[j] = next
						turnMoves = append(turnMoves, fmt.Sprintf("L%d-%s", sim.AntIDs[j], sim.Path[next]))
					}
				}
			}
			copy(sim.Positions, newPos)
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
