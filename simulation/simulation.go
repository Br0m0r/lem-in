package simulation

import (
	"fmt"
	"os"
	"strings"

	"lem-in/scheduling"
	"lem-in/structs"
)

// GeneratePathGrid creates a 2D grid visualization for a single simulation path.
// Each room is displayed as a box; if an ant is present, its label (e.g., L4) is shown.
func GeneratePathGrid(sim structs.PathSim) string {
	var sb strings.Builder
	for i, room := range sim.Path {
		antLabel := ""
		for j, pos := range sim.Positions {
			if pos == i {
				antLabel = fmt.Sprintf("(L%d)", sim.AntIDs[j])
				break
			}
		}
		if antLabel != "" {
			sb.WriteString(fmt.Sprintf("[ %s %s ]", room, antLabel))
		} else {
			sb.WriteString(fmt.Sprintf("[ %s ]", room))
		}
		if i < len(sim.Path)-1 {
			sb.WriteString(" ---> ")
		}
	}
	return sb.String()
}

// SimulateMultiPath simulates ant movement along each path concurrently.
// It prints minimal move information to the terminal and writes a detailed simulation
// (including extra info and full grid visualization) to simulation_output.txt.
func SimulateMultiPath(antCount int, paths [][]string, assignment scheduling.PathAssignment, extraInfo string) {
	// Build simulation state for each path.
	sims := make([]structs.PathSim, len(paths))
	antCounter := 1
	for i, p := range paths {
		count := assignment.AntsPerPath[i]
		positions := make([]int, count)
		for j := range positions {
			positions[j] = -1 // Initialize as not injected.
		}
		antIDs := make([]int, count)
		for j := 0; j < count; j++ {
			antIDs[j] = antCounter
			antCounter++
		}
		sims[i] = structs.PathSim{Path: p, Positions: positions, AntIDs: antIDs}
	}

	outFile, err := os.Create("simulation_output.txt")
	if err != nil {
		fmt.Println("Error creating simulation_output.txt:", err)
		return
	}
	defer outFile.Close()

	var terminalBuilder strings.Builder
	var fileBuilder strings.Builder

	// Write extra info at the top of the output file.
	fileBuilder.WriteString(extraInfo)
	fileBuilder.WriteString("\n\n")

	turn := 0
	for {
		var turnMoves []string
		fileBuilder.WriteString(fmt.Sprintf("TURN %d\n", turn+1))
		for idx, sim := range sims {
			pathLen := len(sim.Path)
			newPos := make([]int, len(sim.Positions))
			copy(newPos, sim.Positions)
			for j := len(sim.Positions) - 1; j >= 0; j-- {
				if sim.Positions[j] == -1 {
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
			copy(sims[idx].Positions, newPos)
			gridStr := GeneratePathGrid(sim)
			fileBuilder.WriteString(gridStr + "\n")
		}
		if len(turnMoves) == 0 {
			break
		}
		terminalBuilder.WriteString(fmt.Sprintf("Turn %d: %s\n", turn+1, strings.Join(turnMoves, " ")))
		turn++
		fileBuilder.WriteString("\n")
	}

	fileBuilder.WriteString(fmt.Sprintf("Total turns: %d\n", turn))
	_, err = outFile.WriteString(fileBuilder.String())
	if err != nil {
		fmt.Println("Error writing to simulation_output.txt:", err)
	}

	fmt.Print(terminalBuilder.String())
	fmt.Println("2D grid visualization written to simulation_output.txt")
}

// Helper function: returns true if any ant is at the given room index.
func isOccupied(positions []int, pos int) bool {
	for _, p := range positions {
		if p == pos {
			return true
		}
	}
	return false
}
