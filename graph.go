package main

import (
	"errors"
	"fmt"
)

// Graph structure using an adjacency list.
type Graph struct {
	Rooms     map[string]*Room
	Neighbors map[string][]string // mapping room name to names of adjacent rooms.
}

// BuildGraph constructs a graph from rooms and tunnels.
func BuildGraph(rooms []Room, tunnels []Tunnel) (*Graph, error) {
	g := &Graph{
		Rooms:     make(map[string]*Room),
		Neighbors: make(map[string][]string),
	}

	// Add rooms to the graph.
	for i := range rooms {
		room := rooms[i]
		g.Rooms[room.Name] = &room
	}

	// Add tunnels (edges) to the graph.
	for _, tunnel := range tunnels {
		// Check if both rooms exist.
		if _, ok := g.Rooms[tunnel.RoomA]; !ok {
			return nil, fmt.Errorf("ERROR: tunnel references unknown room %s", tunnel.RoomA)
		}
		if _, ok := g.Rooms[tunnel.RoomB]; !ok {
			return nil, fmt.Errorf("ERROR: tunnel references unknown room %s", tunnel.RoomB)
		}
		g.Neighbors[tunnel.RoomA] = append(g.Neighbors[tunnel.RoomA], tunnel.RoomB)
		g.Neighbors[tunnel.RoomB] = append(g.Neighbors[tunnel.RoomB], tunnel.RoomA)
	}

	return g, nil
}

// FindQuickestPaths implements BFS to find a shortest path from start to end.
// For now, this returns a single shortest path.
func FindQuickestPaths(g *Graph) ([][]string, error) {
	var startRoom, endRoom string
	for name, room := range g.Rooms {
		if room.IsStart {
			startRoom = name
		}
		if room.IsEnd {
			endRoom = name
		}
	}

	if startRoom == "" || endRoom == "" {
		return nil, errors.New("ERROR: missing start or end room")
	}

	// BFS initialization.
	queue := []string{startRoom}
	prev := make(map[string]string)
	visited := make(map[string]bool)
	visited[startRoom] = true

	found := false
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current == endRoom {
			found = true
			break
		}
		for _, neighbor := range g.Neighbors[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				prev[neighbor] = current
				queue = append(queue, neighbor)
			}
		}
	}

	if !found {
		return nil, errors.New("ERROR: no path found")
	}

	// Reconstruct the path.
	path := []string{}
	for at := endRoom; at != ""; at = prev[at] {
		path = append([]string{at}, path...)
		if at == startRoom {
			break
		}
	}

	return [][]string{path}, nil
}
