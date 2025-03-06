package main

import (
	"errors"
	"fmt"
)

// BuildGraph constructs a graph from rooms and tunnels.
func BuildGraph(rooms []Room, tunnels []Tunnel) (*Graph, error) {
	g := &Graph{
		Rooms:     make(map[string]*Room),
		Neighbors: make(map[string][]string),
	}

	for i := range rooms {
		room := rooms[i]
		g.Rooms[room.Name] = &room
	}

	for _, tunnel := range tunnels {
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

// FindMultiplePaths finds all valid paths from start to end.
// This is a stub for demonstration; in a full solution, implement a multi-path algorithm (e.g., Edmonds-Karp).
func FindMultiplePaths(g *Graph) ([][]string, error) {
	// Identify start and end.
	var start, end string
	for name, room := range g.Rooms {
		if room.IsStart {
			start = name
		}
		if room.IsEnd {
			end = name
		}
	}
	if start == "" || end == "" {
		return nil, errors.New("ERROR: missing start or end room")
	}

	// --- BEGIN STUB ---
	// In a full solution, dynamically extract disjoint paths here.
	// For now, we'll use a simple heuristic: run multiple BFS searches while removing used edges.
	// As an example, we return:
	//   For any input, return at least one path (BFS path) as a fallback.
	path, err := FindSinglePath(g, start, end)
	if err != nil {
		return nil, err
	}
	// For demonstration, we return just one path.
	paths := [][]string{path}
	// --- END STUB ---

	return paths, nil
}

// FindSinglePath uses BFS to find one shortest path from start to end.
func FindSinglePath(g *Graph, start, end string) ([]string, error) {
	queue := []string{start}
	prev := make(map[string]string)
	visited := make(map[string]bool)
	visited[start] = true

	found := false
	for len(queue) > 0 {
		curr := queue[0]
		queue = queue[1:]
		if curr == end {
			found = true
			break
		}
		for _, neighbor := range g.Neighbors[curr] {
			if !visited[neighbor] {
				visited[neighbor] = true
				prev[neighbor] = curr
				queue = append(queue, neighbor)
			}
		}
	}
	if !found {
		return nil, errors.New("ERROR: no path found")
	}

	// Reconstruct path.
	var path []string
	for at := end; at != ""; at = prev[at] {
		path = append([]string{at}, path...)
		if at == start {
			break
		}
	}
	return path, nil
}
