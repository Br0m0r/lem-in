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

// CopyGraph creates a deep copy of the given graph.
func CopyGraph(g *Graph) *Graph {
	newG := &Graph{
		Rooms:     make(map[string]*Room),
		Neighbors: make(map[string][]string),
	}
	for k, room := range g.Rooms {
		newG.Rooms[k] = room // Rooms are immutable here.
	}
	for k, neighbors := range g.Neighbors {
		newNeighbors := make([]string, len(neighbors))
		copy(newNeighbors, neighbors)
		newG.Neighbors[k] = newNeighbors
	}
	return newG
}

// FindMultiplePaths finds all vertex-disjoint paths from start to end by repeatedly running BFS
// and removing intermediate vertices from the copied graph.
func FindMultiplePaths(g *Graph) ([][]string, error) {
	var paths [][]string

	Gcopy := CopyGraph(g)

	// Identify start and end.
	var start, end string
	for name, room := range Gcopy.Rooms {
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

	for {
		path, err := FindSinglePath(Gcopy, start, end)
		if err != nil {
			break
		}
		paths = append(paths, path)
		// Remove intermediate vertices (except start and end) from Gcopy.
		for i := 1; i < len(path)-1; i++ {
			v := path[i]
			delete(Gcopy.Rooms, v)
			delete(Gcopy.Neighbors, v)
			// Remove v from all neighbor lists.
			for k, nList := range Gcopy.Neighbors {
				newList := []string{}
				for _, w := range nList {
					if w != v {
						newList = append(newList, w)
					}
				}
				Gcopy.Neighbors[k] = newList
			}
		}
	}

	if len(paths) == 0 {
		return nil, errors.New("ERROR: no path found")
	}
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
	var path []string
	for at := end; at != ""; at = prev[at] {
		path = append([]string{at}, path...)
		if at == start {
			break
		}
	}
	return path, nil
}
