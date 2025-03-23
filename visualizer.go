package main

import (
	"fmt"
	"strings"
)

// PrintExtraInfo generates a string containing extra information
// (such as input details, summary, and found/selected paths) that is
// written at the top of the simulation output file.
func PrintExtraInfo(antCount int, rooms []Room, tunnels []Tunnel, paths [][]string, assignment PathAssignment) string {
	var sb strings.Builder

	// Echo input data.
	sb.WriteString(fmt.Sprintf("%d\n", antCount))
	for _, room := range rooms {
		if room.IsStart {
			sb.WriteString("##start\n")
		}
		if room.IsEnd {
			sb.WriteString("##end\n")
		}
		sb.WriteString(fmt.Sprintf("%s %d %d\n", room.Name, room.X, room.Y))
	}
	for _, tunnel := range tunnels {
		sb.WriteString(fmt.Sprintf("%s-%s\n", tunnel.RoomA, tunnel.RoomB))
	}
	sb.WriteString("\n")

	// Include the summary.
	sb.WriteString("----------- Summary -----------\n")
	sb.WriteString(fmt.Sprintf("Number of ants: %d\n", antCount))
	sb.WriteString(fmt.Sprintf("Number of rooms: %d\n", len(rooms)))
	sb.WriteString(fmt.Sprintf("Number of tunnels: %d\n", len(tunnels)))
	var startRoom, endRoom string
	for _, room := range rooms {
		if room.IsStart {
			startRoom = room.Name
		}
		if room.IsEnd {
			endRoom = room.Name
		}
	}
	sb.WriteString(fmt.Sprintf("Start room: %s\n", startRoom))
	sb.WriteString(fmt.Sprintf("End room: %s\n", endRoom))
	sb.WriteString("\n")

	// List all found paths.
	sb.WriteString("---------- All Found Paths ----------\n")
	sb.WriteString(fmt.Sprintf("Number of possible paths: %d\n", len(paths)))
	for i, p := range paths {
		sb.WriteString(fmt.Sprintf("%d) %s\n", i+1, strings.Join(p, " -> ")))
	}
	sb.WriteString("\n")

	// List selected paths based on the assignment.
	sb.WriteString("---------- Selected Paths ---------- \n")
	for i, p := range assignment.Paths {
		sb.WriteString(fmt.Sprintf("%d) %s\n", i+1, strings.Join(p, " -> ")))
	}
	sb.WriteString("\n")

	return sb.String()
}
