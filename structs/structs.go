package structs

// Room holds the information for a single room in the ant farm.

type Room struct {
	Name    string
	X       int
	Y       int
	IsStart bool
	IsEnd   bool
}

// Tunnel represents a connection between two rooms.
type Tunnel struct {
	RoomA string
	RoomB string
}

// Graph represents the ant farm as a whole.
type Graph struct {
	Rooms     map[string]*Room
	Neighbors map[string][]string
}

// PathAssignment holds the result of assigning ants to the available paths.
type PathAssignment struct {
	Paths       [][]string
	AntsPerPath []int
}

// PathSim represents the simulation state for a single path during the ant movement simulation.
type PathSim struct {
	Path      []string
	Positions []int
	AntIDs    []int // Unique identifiers for the ants on this path.
}
