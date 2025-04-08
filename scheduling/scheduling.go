package scheduling

import (
	"lem-in/structs"
	"sort"
)

// AssignAnts distributes the ants among the available paths using a simple greedy algorithm.
// It calculates an "effective cost" for each path, which is defined as:
//
//	cost = (path length) + (number of ants already assigned) - 1  ( - 1 to not count the starting room as a used move)
func AssignAnts(antCount int, paths [][]string) structs.PathAssignment {

	numPaths := len(paths)

	// antsPerPath will count how many ants are assigned to each path.
	antsPerPath := make([]int, numPaths)

	// Process each ant one by one.
	for i := 0; i < antCount; i++ {

		type pathCost struct {
			index int //position of the path in the slice.
			cost  int //the computed effective cost for that path.
		}

		pathCosts := make([]pathCost, numPaths) //store the cost for each available path.

		// Calculate the effective cost for each path.
		for j := 0; j < numPaths; j++ {
			cost := len(paths[j]) + antsPerPath[j] - 1
			pathCosts[j] = pathCost{index: j, cost: cost}
		}

		// Sort the pathCosts slice so that the path with the smallest cost comes first.
		sort.Slice(pathCosts, func(a, b int) bool {
			return pathCosts[a].cost < pathCosts[b].cost
		})

		// Assign the current ant to the path with the lowest effective cost.
		antsPerPath[pathCosts[0].index]++
	}

	return structs.PathAssignment{Paths: paths, AntsPerPath: antsPerPath}
}
