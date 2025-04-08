package graph

import (
	"errors"
	"fmt"

	"lem-in/structs"
)

// BuildGraph creates a graph (map) of the ant farm using the list of rooms and tunnels.
func BuildGraph(roomList []structs.Room, connections []structs.Tunnel) (*structs.Graph, error) {
	graphData := &structs.Graph{
		Rooms:     make(map[string]*structs.Room), // Keys will be room names.
		Neighbors: make(map[string][]string),      // Each room name will map to a list of connected room names.
	}

	for i := range roomList {
		currentRoom := roomList[i]
		graphData.Rooms[currentRoom.Name] = &currentRoom
	}

	// Process each tunnel and update the list of connected rooms.
	for _, tunnel := range connections {
		if _, found := graphData.Rooms[tunnel.RoomA]; !found {
			return nil, fmt.Errorf("ERROR: tunnel refers to unknown room %s", tunnel.RoomA)
		}
		if _, found := graphData.Rooms[tunnel.RoomB]; !found {
			return nil, fmt.Errorf("ERROR: tunnel refers to unknown room %s", tunnel.RoomB)
		}

		graphData.Neighbors[tunnel.RoomA] = append(graphData.Neighbors[tunnel.RoomA], tunnel.RoomB)
		graphData.Neighbors[tunnel.RoomB] = append(graphData.Neighbors[tunnel.RoomB], tunnel.RoomA)
	}

	return graphData, nil
}

// FindMultiplePaths finds all separate paths (without reusing tunnels) from the start room to the end room
// using a breadth-first search (BFS) approach.
func FindMultiplePaths(graphData *structs.Graph) ([][]string, error) {
	var startRoom, endRoom string

	for name, roomData := range graphData.Rooms {
		if roomData.IsStart {
			startRoom = name
		}
		if roomData.IsEnd {
			endRoom = name
		}
	}
	if startRoom == "" || endRoom == "" {
		return nil, errors.New("ERROR: missing start or end room")
	}

	// Make a copy of the Neighbors map (connections) so that we can modify it as we remove used tunnels.
	connectionsCopy := make(map[string][]string)
	for roomName, connectedList := range graphData.Neighbors {
		newList := make([]string, len(connectedList))
		copy(newList, connectedList)
		connectionsCopy[roomName] = newList
	}

	var foundPaths [][]string

	// Repeatedly search for a new path using BFS.
	for {
		path, pathFound := bfs(connectionsCopy, startRoom, endRoom)
		if !pathFound {
			break
		}
		foundPaths = append(foundPaths, path)
		// Remove the tunnels used in this path so they cannot be used again.
		removePathEdges(connectionsCopy, path)
	}

	if len(foundPaths) == 0 {
		return nil, errors.New("ERROR: no valid paths found")
	}
	return foundPaths, nil
}

// bfs performs a breadth-first search to find one path from the startRoom to the endRoom.
func bfs(connections map[string][]string, startRoom, endRoom string) ([]string, bool) {
	queue := []string{startRoom}          // Start with the startRoom
	visited := make(map[string]bool)      //track visited rooms.
	parentRoom := make(map[string]string) // record how we reached each room.
	visited[startRoom] = true

	for len(queue) > 0 {
		currentRoom := queue[0]
		queue = queue[1:]
		// If we've reached the end room rebuild the path.
		if currentRoom == endRoom {
			var path []string
			// Reconstruct the path by going backwards from endRoom using the parentRoom map.
			for room := endRoom; room != ""; room = parentRoom[room] {
				path = append([]string{room}, path...)
			}
			return path, true
		}
		// Check every room directly connected to the current room.
		for _, nextRoom := range connections[currentRoom] {
			if !visited[nextRoom] {
				visited[nextRoom] = true
				parentRoom[nextRoom] = currentRoom
				queue = append(queue, nextRoom)
			}
		}
	}
	return nil, false
}

// removePathEdges removes the tunnels used in a given path from the connections map.
// The reverse connections are kept.
func removePathEdges(connections map[string][]string, roomPath []string) {
	for i := 0; i < len(roomPath)-1; i++ {
		fromRoom, toRoom := roomPath[i], roomPath[i+1]
		connections[fromRoom] = removeEdge(connections[fromRoom], toRoom)
	}
}

// removeEdge removes a specific room name from a list of connected rooms.
func removeEdge(connectionList []string, roomNameToRemove string) []string {
	newList := connectionList[:0] // Reuse the same underlying array.
	for _, roomName := range connectionList {
		if roomName != roomNameToRemove {
			newList = append(newList, roomName)
		}
	}
	return newList
}
