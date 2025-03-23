package main

import (
	"errors"
	"fmt"
)

// BuildGraph constructs a graph from the provided rooms and tunnels.
// It maps each room name to its Room struct and builds an adjacency list of neighbors.
func BuildGraph(rooms []Room, tunnels []Tunnel) (*Graph, error) {
	g := &Graph{
		Rooms:     make(map[string]*Room),
		Neighbors: make(map[string][]string),
	}

	// Add all rooms to the graph.
	for i := range rooms {
		room := rooms[i]
		g.Rooms[room.Name] = &room
	}

	// Add tunnels as bidirectional edges.
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

// ==================
// Max-Flow Functions for Finding Multiple Paths
// ==================

// BuildFlowNetwork creates a flow network from the given graph.
// Each edge in the graph is converted into a directed edge with a capacity of 1.
func BuildFlowNetwork(g *Graph) *FlowNetwork {
	network := &FlowNetwork{
		Adjacency: make(map[string][]*FlowEdge),
	}
	// Add each room as a node in the flow network.
	for node := range g.Rooms {
		network.Nodes = append(network.Nodes, node)
	}
	// For each neighbor relationship, add a directed edge with capacity 1.
	for room, neighbors := range g.Neighbors {
		for _, neighbor := range neighbors {
			edge := &FlowEdge{From: room, To: neighbor, Capacity: 1, Flow: 0}
			network.Adjacency[room] = append(network.Adjacency[room], edge)
		}
	}
	return network
}

// EdmondsKarp runs the max-flow algorithm to compute the maximum number of edge-disjoint paths.
// Each augmenting path found increases the flow by 1.
func EdmondsKarp(network *FlowNetwork, start, end string) int {
	maxFlow := 0
	for {
		// Perform BFS to find an augmenting path.
		parent := make(map[string]*FlowEdge)
		queue := []string{start}
		for len(queue) > 0 && parent[end] == nil {
			current := queue[0]
			queue = queue[1:]
			for _, edge := range network.Adjacency[current] {
				residual := edge.Capacity - edge.Flow
				if residual > 0 && parent[edge.To] == nil && edge.To != start {
					parent[edge.To] = edge
					queue = append(queue, edge.To)
					if edge.To == end {
						break
					}
				}
			}
		}
		// If no augmenting path is found, exit the loop.
		if parent[end] == nil {
			break
		}
		// Augment the flow along the found path by 1.
		for node := end; node != start; {
			edge := parent[node]
			edge.Flow += 1
			node = edge.From
		}
		maxFlow++
	}
	return maxFlow
}

// ExtractPaths retrieves all edge-disjoint paths that have a flow from start to end.
// As each path is extracted, the used flow is removed to avoid duplicate paths.
func ExtractPaths(network *FlowNetwork, start, end string) [][]string {
	var paths [][]string
	for {
		var path []string
		current := start
		path = append(path, current)
		for current != end {
			found := false
			for _, edge := range network.Adjacency[current] {
				if edge.Flow > 0 { // Edge is used in the flow.
					path = append(path, edge.To)
					edge.Flow = 0 // Remove flow so it isn't reused.
					current = edge.To
					found = true
					break
				}
			}
			if !found {
				break
			}
		}
		// Break if no valid path from start to end is found.
		if len(path) == 0 || path[len(path)-1] != end {
			break
		}
		paths = append(paths, path)
	}
	return paths
}

// FindMultiplePaths finds all edge-disjoint paths from the start to the end room using the max-flow approach.
// It returns an error if no valid paths are found.
func FindMultiplePaths(g *Graph) ([][]string, error) {
	var start, end string
	// Identify start and end rooms.
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

	// Build the flow network from the graph.
	network := BuildFlowNetwork(g)
	// Run the max-flow algorithm.
	maxFlow := EdmondsKarp(network, start, end)
	if maxFlow == 0 {
		return nil, errors.New("ERROR: no valid paths found")
	}
	// Extract and return the paths corresponding to the flows.
	paths := ExtractPaths(network, start, end)
	return paths, nil
}
