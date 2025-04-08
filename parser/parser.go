package parser

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"lem-in/structs"
)

// ParseInputFile reads an input file and returns:
//   - The total number of ants,
//   - A list of rooms,
//   - A list of tunnels (connections between rooms).
func ParseInputFile(filePath string) (int, []structs.Room, []structs.Tunnel, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to open file: %v", err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	var antTotal int
	var roomList []structs.Room
	var tunnelList []structs.Tunnel

	// Create a map to check for duplicate room positions.
	// We use "x,y" as a key to ensure no two rooms share the same coordinates.
	positionMap := make(map[string]bool)

	if scanner.Scan() {
		antStr := strings.TrimSpace(scanner.Text())
		antTotal, err = strconv.Atoi(antStr)
		if err != nil {
			return 0, nil, nil, errors.New("ERROR: invalid number of ants")
		}
	}

	var nextRoomIsStart, nextRoomIsEnd bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue // Skip empty lines.
		}

		// Lines beginning with '#' are either comments or commands.
		if line[0] == '#' {
			if strings.HasPrefix(line, "##start") {
				nextRoomIsStart = true
				continue
			}
			if strings.HasPrefix(line, "##end") {
				nextRoomIsEnd = true
				continue
			}
			// Ignore other comments.
			continue
		}
		// If the line has exactly 3 parts, it's a room definition.
		parts := strings.Fields(line)
		if len(parts) == 3 {
			x, errX := strconv.Atoi(parts[1])
			y, errY := strconv.Atoi(parts[2])
			if errX != nil || errY != nil {
				return 0, nil, nil, errors.New("ERROR: invalid room coordinates")
			}
			// Create a unique key for this room's position.
			posKey := fmt.Sprintf("%d,%d", x, y)
			if _, exists := positionMap[posKey]; exists {
				return 0, nil, nil, errors.New("ERROR: duplicate room coordinates")
			}
			positionMap[posKey] = true // Mark these coordinates as used.

			// Create the room with the provided data.
			newRoom := structs.Room{
				Name:    parts[0],
				X:       x,
				Y:       y,
				IsStart: nextRoomIsStart,
				IsEnd:   nextRoomIsEnd,
			}

			roomList = append(roomList, newRoom)

			nextRoomIsStart = false
			nextRoomIsEnd = false
			continue
		}
		// If the line contains a hyphen, it's a tunnel definition.
		if strings.Contains(line, "-") {
			roomNames := strings.Split(line, "-")
			if len(roomNames) != 2 {
				return 0, nil, nil, errors.New("ERROR: invalid tunnel definition")
			}
			newTunnel := structs.Tunnel{RoomA: roomNames[0], RoomB: roomNames[1]}
			tunnelList = append(tunnelList, newTunnel)
			continue
		}
	}

	return antTotal, roomList, tunnelList, nil
}
