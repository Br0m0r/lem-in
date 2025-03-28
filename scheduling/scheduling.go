package scheduling

import "sort"

// PathAssignment holds the distribution of ants among the available paths.
type PathAssignment struct {
	Paths       [][]string // Each path represented as a slice of room names.
	AntsPerPath []int      // Number of ants assigned to each corresponding path.
}

// AssignAnts distributes ants among paths using a greedy algorithm.
// The effective cost is calculated as: (path length) + (ants already assigned) - 1.
func AssignAnts(antCount int, paths [][]string) PathAssignment {
	numPaths := len(paths)
	antsPerPath := make([]int, numPaths)

	for i := 0; i < antCount; i++ {
		type pathCost struct {
			index int
			cost  int
		}
		pathCosts := make([]pathCost, numPaths)
		for j := 0; j < numPaths; j++ {
			cost := len(paths[j]) + antsPerPath[j] - 1
			pathCosts[j] = pathCost{index: j, cost: cost}
		}
		sort.Slice(pathCosts, func(a, b int) bool {
			return pathCosts[a].cost < pathCosts[b].cost
		})
		antsPerPath[pathCosts[0].index]++
	}
	return PathAssignment{Paths: paths, AntsPerPath: antsPerPath}
}
