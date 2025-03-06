package main

// Room represents a room in the ant farm.
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

// Graph represents the ant farm as an adjacency list.
type Graph struct {
	Rooms     map[string]*Room
	Neighbors map[string][]string // mapping room name to names of adjacent rooms.
}
