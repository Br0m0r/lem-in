package graph

import (
	"errors"
	"fmt"

	"lem-in/structs"
)

// BuildGraph constructs a graph from rooms and tunnels.
// It builds a map of room names to Room structs and an adjacency list.
func BuildGraph(rooms []structs.Room, tunnels []structs.Tunnel) (*structs.Graph, error) {
	g := &structs.Graph{
		Rooms:     make(map[string]*structs.Room),
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

// ---------------------
// Max-Flow Functions
// ---------------------

// BuildFlowNetwork creates a flow network from the graph.
// Each edge gets a capacity of 1.
func BuildFlowNetwork(g *structs.Graph) *structs.FlowNetwork {
	network := &structs.FlowNetwork{
		Adjacency: make(map[string][]*structs.FlowEdge),
	}
	for node := range g.Rooms {
		network.Nodes = append(network.Nodes, node)
	}
	for room, neighbors := range g.Neighbors {
		for _, neighbor := range neighbors {
			edge := &structs.FlowEdge{From: room, To: neighbor, Capacity: 1, Flow: 0}
			network.Adjacency[room] = append(network.Adjacency[room], edge)
		}
	}
	return network
}

// EdmondsKarp runs the max-flow algorithm to find the maximum number of edge-disjoint paths.
// Each found augmenting path increases the flow by 1.
func EdmondsKarp(network *structs.FlowNetwork, start, end string) int {
	maxFlow := 0
	for {
		parent := make(map[string]*structs.FlowEdge)
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
		if parent[end] == nil {
			break
		}
		for node := end; node != start; {
			edge := parent[node]
			edge.Flow += 1
			node = edge.From
		}
		maxFlow++
	}
	return maxFlow
}

// ExtractPaths retrieves all edge-disjoint paths with flow from start to end.
// It removes used flow as paths are extracted.
func ExtractPaths(network *structs.FlowNetwork, start, end string) [][]string {
	var paths [][]string
	for {
		var path []string
		current := start
		path = append(path, current)
		for current != end {
			found := false
			for _, edge := range network.Adjacency[current] {
				if edge.Flow > 0 {
					path = append(path, edge.To)
					edge.Flow = 0
					current = edge.To
					found = true
					break
				}
			}
			if !found {
				break
			}
		}
		if len(path) == 0 || path[len(path)-1] != end {
			break
		}
		paths = append(paths, path)
	}
	return paths
}

// FindMultiplePaths finds all edge-disjoint paths from start to end using the max-flow approach.
func FindMultiplePaths(g *structs.Graph) ([][]string, error) {
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
	network := BuildFlowNetwork(g)
	maxFlow := EdmondsKarp(network, start, end)
	if maxFlow == 0 {
		return nil, errors.New("ERROR: no valid paths found")
	}
	paths := ExtractPaths(network, start, end)
	return paths, nil
}
