package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ParseInputFile reads the input file and returns the ant count, list of rooms, tunnels, and an error if any.
// It expects the first line to be the ant count, then room definitions (with "##start" and "##end" commands),
// followed by tunnel definitions.
func ParseInputFile(filename string) (int, []Room, []Tunnel, error) {
	file, err := os.Open(filename)
	if err != nil {
		return 0, nil, nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var antCount int
	var rooms []Room
	var tunnels []Tunnel

	// Read the ant count (first line).
	if scanner.Scan() {
		countStr := strings.TrimSpace(scanner.Text())
		antCount, err = strconv.Atoi(countStr)
		if err != nil {
			return 0, nil, nil, errors.New("ERROR: invalid number of ants")
		}
	}

	var isNextRoomStart, isNextRoomEnd bool
	// Process the rest of the file line by line.
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 {
			continue
		}
		// Process commands (lines starting with "#").
		if line[0] == '#' {
			if strings.HasPrefix(line, "##start") {
				isNextRoomStart = true
				continue
			}
			if strings.HasPrefix(line, "##end") {
				isNextRoomEnd = true
				continue
			}
			continue
		}
		// Room definition: "name x y".
		parts := strings.Fields(line)
		if len(parts) == 3 {
			x, errX := strconv.Atoi(parts[1])
			y, errY := strconv.Atoi(parts[2])
			if errX != nil || errY != nil {
				return 0, nil, nil, errors.New("ERROR: invalid room coordinates")
			}
			room := Room{
				Name:    parts[0],
				X:       x,
				Y:       y,
				IsStart: isNextRoomStart,
				IsEnd:   isNextRoomEnd,
			}
			rooms = append(rooms, room)
			isNextRoomStart = false
			isNextRoomEnd = false
			continue
		}
		// Tunnel definition: "roomA-roomB".
		if strings.Contains(line, "-") {
			roomNames := strings.Split(line, "-")
			if len(roomNames) != 2 {
				return 0, nil, nil, errors.New("ERROR: invalid tunnel definition")
			}
			tunnels = append(tunnels, Tunnel{RoomA: roomNames[0], RoomB: roomNames[1]})
			continue
		}
	}

	return antCount, rooms, tunnels, nil
}
