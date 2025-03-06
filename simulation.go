package main

import (
	"fmt"
)

// SimulateAntMovements takes the ant count and paths and prints the ant movements turn-by-turn.
// This example assumes a single path. A more complex version could handle multiple paths.
func SimulateAntMovements(antCount int, paths [][]string) {
	// For simplicity, assume one path.
	path := paths[0]
	pathLength := len(path)

	// positions holds the current position index for each ant on the path (-1 indicates not yet entered the path).
	positions := make([]int, antCount)

	// Continue until all ants have reached the final room (index pathLength-1).
	for !allAntsAtEnd(positions, pathLength-1) {
		moves := []string{}

		// Move ants from the end of the path to the beginning to avoid collisions.
		for i := antCount - 1; i >= 0; i-- {
			if positions[i] < pathLength-1 {
				if canMove(positions, i) {
					positions[i]++
					moves = append(moves, fmt.Sprintf("L%d-%s", i+1, path[positions[i]]))
				}
			}
		}

		// Print all moves for this turn.
		if len(moves) > 0 {
			fmt.Println(joinMoves(moves))
		}
	}
}

// allAntsAtEnd checks if all ants have reached the final room.
func allAntsAtEnd(positions []int, endPos int) bool {
	for _, pos := range positions {
		if pos != endPos {
			return false
		}
	}
	return true
}

// canMove checks if ant at index antIndex can move forward (ensuring the next position is unoccupied).
func canMove(positions []int, antIndex int) bool {
	nextPos := positions[antIndex] + 1
	for j, pos := range positions {
		if j != antIndex && pos == nextPos {
			return false
		}
	}
	return true
}

// joinMoves concatenates the moves into a single string.
func joinMoves(moves []string) string {
	return fmt.Sprint(moves)
}
