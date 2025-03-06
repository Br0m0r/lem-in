package main

import (
	"fmt"
	"strings"
)

// PrintSummary prints a textual summary in English, similar to your screenshot.
func PrintSummary(
	antCount int,
	rooms []Room,
	tunnels []Tunnel,
	paths [][]string,         // all found paths
	selectedPaths [][]string, // chosen/used paths, e.g. after scheduling
) {
	// 1) Basic summary
	fmt.Println("----------- Summary -----------")
	fmt.Printf("Number of ants: %d\n", antCount)
	fmt.Printf("Number of rooms: %d\n", len(rooms))
	fmt.Printf("Number of tunnels: %d\n", len(tunnels))

	// Identify start/end rooms
	var startRoomName, endRoomName string
	for _, r := range rooms {
		if r.IsStart {
			startRoomName = r.Name
		}
		if r.IsEnd {
			endRoomName = r.Name
		}
	}
	fmt.Printf("Start room: %s\n", startRoomName)
	fmt.Printf("End room: %s\n", endRoomName)

	// 2) Print room IDs and names
	//    Build a map from room name -> ID
	roomIDMap := make(map[string]int)
	for i, r := range rooms {
		roomIDMap[r.Name] = i
	}

	fmt.Println("\n---------- Room IDs and Names ----------")
	for i, r := range rooms {
		tag := ""
		if r.IsStart {
			tag = " (start)"
		} else if r.IsEnd {
			tag = " (end)"
		}
		fmt.Printf("%d => %s%s\n", i, r.Name, tag)
	}

	// 3) Build an adjacency matrix (or adjacency list) to show the connections.
	n := len(rooms)
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		matrix[i] = make([]int, n)
	}
	for _, t := range tunnels {
		i := roomIDMap[t.RoomA]
		j := roomIDMap[t.RoomB]
		matrix[i][j] = 1
		matrix[j][i] = 1
	}

	fmt.Println("\n---------- Room Connections (Adjacency Matrix) ----------")
	header := "    "
	for i := 0; i < n; i++ {
		header += fmt.Sprintf("%2d ", i)
	}
	fmt.Println(header)
	for i := 0; i < n; i++ {
		rowStr := fmt.Sprintf("%2d |", i)
		for j := 0; j < n; j++ {
			rowStr += fmt.Sprintf(" %d ", matrix[i][j])
		}
		fmt.Println(rowStr)
	}

	// 4) Print all found paths
	fmt.Println("\n---------- All Found Paths ----------")
	fmt.Printf("Number of possible paths: %d\n", len(paths))
	for i, p := range paths {
		fmt.Printf("%d) %s\n", i+1, strings.Join(p, " -> "))
	}

	// 5) Print selected paths
	fmt.Println("\n---------- Selected Paths ----------")
	if len(selectedPaths) == 0 {
		fmt.Println("(No selected paths)")
	} else {
		for i, p := range selectedPaths {
			fmt.Printf("%d) %s\n", i+1, strings.Join(p, " -> "))
		}
	}
	fmt.Println()
}
