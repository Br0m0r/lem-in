package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

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

// ParseInputFile reads the file and returns ant count, rooms, and tunnels.
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

	// Read ant count (first line)
	if scanner.Scan() {
		countStr := strings.TrimSpace(scanner.Text())
		antCount, err = strconv.Atoi(countStr)
		if err != nil {
			return 0, nil, nil, errors.New("ERROR: invalid number of ants")
		}
	}

	var isNextRoomStart, isNextRoomEnd bool
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore comments (except commands)
		if len(line) == 0 {
			continue
		}
		if line[0] == '#' {
			// Handle commands
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

		// Check if line is a room definition (room name and coordinates)
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

		// Check if line is a tunnel definition (roomA-roomB)
		if strings.Contains(line, "-") {
			roomNames := strings.Split(line, "-")
			if len(roomNames) != 2 {
				return 0, nil, nil, errors.New("ERROR: invalid tunnel definition")
			}
			tunnels = append(tunnels, Tunnel{RoomA: roomNames[0], RoomB: roomNames[1]})
			continue
		}
	}

	// Additional validations (e.g., ensuring there is exactly one start and one end room) can be added here.

	return antCount, rooms, tunnels, nil
}

// PrintInputData echoes the input data if required.
func PrintInputData(antCount int, rooms []Room, tunnels []Tunnel) {
	fmt.Println(antCount)
	for _, room := range rooms {
		if room.IsStart {
			fmt.Println("##start")
		}
		if room.IsEnd {
			fmt.Println("##end")
		}
		fmt.Printf("%s %d %d\n", room.Name, room.X, room.Y)
	}
	for _, tunnel := range tunnels {
		fmt.Printf("%s-%s\n", tunnel.RoomA, tunnel.RoomB)
	}
}
