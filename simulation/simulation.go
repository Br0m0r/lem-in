package simulation

import (
	"fmt"
	"lem-in/structs"
	"lem-in/visualizer"
	"strings"
)

// initSimulation sets up the simulation state for each path.
// It creates a PathSim for every available path, initializing each ant's starting position and assigning unique IDs.
func initSimulation(pathList [][]string, assignment structs.PathAssignment) []structs.PathSim {
	simStates := make([]structs.PathSim, len(pathList))
	antIDCounter := 1 // Start assigning ant IDs from 1.

	for i, path := range pathList {
		antCountForPath := assignment.AntsPerPath[i] // Get the number of ants for the current path from the assignment.

		positions := make([]int, antCountForPath)
		for j := range positions {
			positions[j] = -1
		}

		antIDs := make([]int, antCountForPath) //for antIds
		for j := 0; j < antCountForPath; j++ {
			antIDs[j] = antIDCounter
			antIDCounter++
		}

		// Save the simulation state for this path.
		simStates[i] = structs.PathSim{
			Path:      path,
			Positions: positions,
			AntIDs:    antIDs,
		}
	}
	return simStates
}

// isRoomOccupied checks whether any ant is currently at the specified room index in the path.
func isRoomOccupied(positions []int, roomIndex int) bool {
	for _, pos := range positions {
		if pos == roomIndex {
			return true
		}
	}
	return false
}

// processTurn simulates one turn of ant movements on all paths.
// It attempts to inject new ants into the path or move ants already in the path forward.
func processTurn(simStates []structs.PathSim) (moveDescriptions []string, gridVisualization string) {
	var gridBuilder strings.Builder // To build the visual grid output.

	for idx := range simStates { // Loop over each simulation state
		simState := &simStates[idx]
		pathLength := len(simState.Path)

		newPositions := make([]int, len(simState.Positions)) // newPositions is a temporary slice to hold updated ant positions for this turn.
		copy(newPositions, simState.Positions)

		if pathLength == 2 {
			for j := 0; j < len(simState.Positions); j++ {
				if simState.Positions[j] == -1 { //inject ant if not yet injected
					newPositions[j] = 1
					moveDescriptions = append(moveDescriptions, fmt.Sprintf("L%d-%s", simState.AntIDs[j], simState.Path[1])) //record the move
					break
				}
			}
		} else {
			// process ants in reverse order.Ensures that ants ahead move first, which can free up space for those behind.
			for j := len(simState.Positions) - 1; j >= 0; j-- {
				if simState.Positions[j] == -1 {
					if !isRoomOccupied(newPositions, 1) {

						newPositions[j] = 1
						moveDescriptions = append(moveDescriptions, fmt.Sprintf("L%d-%s", simState.AntIDs[j], simState.Path[1]))
					}
				} else if simState.Positions[j] < pathLength-1 { // For ants already in the path but not yet at the end, attempt to move them forward.
					nextIndex := simState.Positions[j] + 1
					if nextIndex == pathLength-1 || !isRoomOccupied(newPositions, nextIndex) {
						newPositions[j] = nextIndex
						moveDescriptions = append(moveDescriptions, fmt.Sprintf("L%d-%s", simState.AntIDs[j], simState.Path[nextIndex]))
					}
				}
			}
		}
		//update simState.Positions with the newPositions computed.
		copy(simState.Positions, newPositions)
		gridBuilder.WriteString(visualizer.GeneratePathGrid(*simState) + "\n") //generate a grid view for the current path.
	}
	// Return the collected move descriptions and the complete grid visualization for this turn.
	return moveDescriptions, gridBuilder.String()
}

// SimulateMultiPath runs the simulation until no ant moves occur.
// It outputs a detailed grid visualization to "simulation_output.txt" and prints a summary of moves to the terminal.
func SimulateMultiPath(antTotal int, pathList [][]string, assignment structs.PathAssignment, headerInfo string) {
	simStates := initSimulation(pathList, assignment)
	var gridOutputs []string
	var moveOutputs []string
	turnCount := 0

	// Run simulation turns repeatedly until no moves are made.
	for {
		moves, grid := processTurn(simStates)
		if len(moves) == 0 {
			break
		}

		gridOutputs = append(gridOutputs, grid)
		moveOutputs = append(moveOutputs, strings.Join(moves, " "))
		turnCount++
		// Detailed note: Each loop iteration represents one simulation turn where ants are moved.
	}

	err := visualizer.WriteSimulationOutput("simulation_output.txt", headerInfo, gridOutputs, turnCount)
	if err != nil {
		fmt.Println("Error writing simulation output:", err)
	}

	visualizer.PrintTerminalOutput(moveOutputs)
	fmt.Println("2D grid visualization written to simulation_output.txt")
}
