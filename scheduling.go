package main

// PathAssignment holds the assignment of ants to each path.
type PathAssignment struct {
	Paths       [][]string // each path as a slice of room names
	AntsPerPath []int      // number of ants assigned to each corresponding path
}

// AssignAnts distributes ants among the available paths using a greedy load-balancing strategy.
func AssignAnts(antCount int, paths [][]string) PathAssignment {
	numPaths := len(paths)
	antsPerPath := make([]int, numPaths)

	// Greedily assign each ant to the path that minimizes (path_length + ants_assigned - 1).
	for i := 0; i < antCount; i++ {
		best := 0
		bestScore := len(paths[0]) + antsPerPath[0] - 1
		for j := 1; j < numPaths; j++ {
			score := len(paths[j]) + antsPerPath[j] - 1
			if score < bestScore {
				bestScore = score
				best = j
			}
		}
		antsPerPath[best]++
	}
	return PathAssignment{Paths: paths, AntsPerPath: antsPerPath}
}
