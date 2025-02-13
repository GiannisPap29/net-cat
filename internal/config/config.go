package parse

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"lem-in/helpers"
	"lem-in/models"
)

// ParseInput reads a file containing a graph-based representation of an ant colony
// and parses it into a Field struct. Returns the Field struct, a slice of all lines
// from the file, and an error if any.
func ParseInput(filename string) (*models.Field, []string, error) {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open file: %v", err)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open file: %v", err)
	}

	defer file.Close() // Ensure the file is closed when the function exits

	// Initialize the Field struct to store rooms, links, and ants
	field := &models.Field{
		Ants:  make([]*models.Ant, 0),
		Rooms: make([]*models.Room, 0),
	}
	fmt.Println(file)
	scanner := bufio.NewScanner(file)

	// Read the first line to get the number of ants
	if !scanner.Scan() {
		return nil, nil, fmt.Errorf("empty file")
	}
	numAnts, err := strconv.Atoi(scanner.Text())
	if err != nil || numAnts <= 0 {
		return nil, nil, fmt.Errorf("invalid number of ants")
	}

	parsingLinks := false // Flag to determine when we start parsing links
	startFound := false   // check if there are more than one ##start
	endFound := false     // check if there are more than one ##end
	// Process each line from the file
	for scanner.Scan() {
		line := scanner.Text()

		// Ignore comments that start with "#" (except special commands like "##start" and "##end")
		if strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "##") {
			continue
		}

		// Handle the start room declaration
		if line == "##start" {
			if startFound {
				return nil, nil, fmt.Errorf("invalid data format, multiple ##start commands found")
			}
			startFound = true
			if !scanner.Scan() {
				return nil, nil, fmt.Errorf("missing start room")
			}
			room, err := parseRoom(field, scanner.Text())
			if err != nil {
				return nil, nil, err
			}
			field.StartRoomName = room.Name
			field.Rooms = append(field.Rooms, room)
			continue
		}
		// Handle the end room declaration
		if line == "##end" {
			if endFound {
				return nil, nil, fmt.Errorf("invalid data format, multiple ##end commands found")
			}
			endFound = true
			if !scanner.Scan() {
				return nil, nil, fmt.Errorf("missing end room")
			}
			room, err := parseRoom(field, scanner.Text())
			if err != nil {
				return nil, nil, err
			}
			field.EndRoomName = room.Name
			field.Rooms = append(field.Rooms, room)
			continue
		}

		// If the line contains a "-", it's a link definition (connection between rooms)
		if strings.Contains(line, "-") {
			parsingLinks = true
			if err := addLink(field, line); err != nil {
				return nil, nil, err
			}
			continue
		}

		// If we haven't started parsing links, assume it's a room definition
		if !parsingLinks {
			room, err := parseRoom(field, line) // Parse room details
			if err != nil {
				return nil, nil, err
			}
			field.Rooms = append(field.Rooms, room) // Add the room to the field
		}
	}

	// Initialize ants and assign them to the start room
	for i := 1; i <= numAnts; i++ {
		field.Ants = append(field.Ants, &models.Ant{
			Id:          i,
			CurrentRoom: field.StartRoomName, // Each ant starts in the start room
			IsFinished:  false,
		})
	}

	return field, []string{string(data)}, nil
}

// parseRoom parses a single line of input into a Room object.
// Expected format: "<room_name> <x_coordinate> <y_coordinate>"
// Returns a Room struct if valid, otherwise returns an error.
func parseRoom(f *models.Field, line string) (*models.Room, error) {
	parts := strings.Fields(line) // Split the line into words
	if len(parts) != 3 {          // A valid room definition should have exactly 3 parts
		return nil, fmt.Errorf("invalid room format: %s", line)
	}

	// check if room already exists
	for _, room := range f.Rooms {
		if room.Name == parts[0] {
			return nil, fmt.Errorf("room already exists: %s", parts[0])
		}
	}

	x, err := strconv.Atoi(parts[1]) // Parse the x-coordinate
	if err != nil {
		return nil, fmt.Errorf("invalid x-coordinate: %s", parts[1])
	}
	y, err := strconv.Atoi(parts[2]) // Parse the y-coordinate
	if err != nil {
		return nil, fmt.Errorf("invalid y-coordinate: %s", parts[2])
	}

	// check if any other room has the same coordinates
	for _, room := range f.Rooms {
		if room.X == x && room.Y == y {
			return nil, fmt.Errorf("room with same coordinates already exists: %s", line)
		}
	}

	return &models.Room{
		Name:           parts[0],          // First part is the room name
		X:              x,                 // Second part is the x-coordinate
		Y:              y,                 // Third part is the y-coordinate
		ConnectedRooms: make([]string, 0), // Initialize an empty list of connected rooms
	}, nil
}

// addLink establishes a connection between two rooms by parsing a link definition line.
// Expected format: "<room1>-<room2>"
// Returns an error if the format is invalid or if the rooms do not exist.
func addLink(field *models.Field, line string) error {
	parts := strings.Split(line, "-") // Split the line by "-"
	if len(parts) != 2 {              // A valid link must contain exactly two room names
		return fmt.Errorf("invalid link format: %s", line)
	}

	room1 := FindRoom(field, parts[0]) // Find the first room
	room2 := FindRoom(field, parts[1]) // Find the second room

	// Ensure both rooms exist
	if room1 == nil || room2 == nil {
		// find which room is missing
		if room1 == nil { // room1 is missing
			return fmt.Errorf("link references non-existent room: %s", parts[0])
		}
		return fmt.Errorf("link references non-existent room: %s", parts[1]) // room2 is missing
	}

	// Check if the link already exists to prevent duplicate connections
	if helpers.Contains(room1.ConnectedRooms, room2.Name) || helpers.Contains(room2.ConnectedRooms, room1.Name) {
		return fmt.Errorf("link already exists: %s", line)
	}

	// check if link links to itself
	if room1.Name == room2.Name {
		return fmt.Errorf("link links to itself: %s", line)
	}

	// Add connections both ways (since this is an undirected graph)
	room1.ConnectedRooms = append(room1.ConnectedRooms, room2.Name)
	room2.ConnectedRooms = append(room2.ConnectedRooms, room1.Name)

	return nil
}

// findRoom searches for a room by name within the field's list of rooms.
// If the room exists, it returns a pointer to the Room struct; otherwise, it returns nil.
func FindRoom(field *models.Field, name string) *models.Room {
	for _, room := range field.Rooms {
		if room.Name == name {
			return room // Return the found room
		}
	}
	return nil // Room not found
}






package main

import (
	"fmt"
	"os"

	"lem-in/helpers"
	"lem-in/parse"
	"lem-in/path"
	"lem-in/solver"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("ERROR: invalid number of arguments")
		return
	}

	// Parse input
	field, lines, err := parse.ParseInput(os.Args[1])
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	// Find all possible paths
	paths, err := path.FindPaths(field)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		return
	}

	// Find non-conflicting paths
	nonConflictingPaths := path.FindNonConflictingPaths(paths)

	// Distribute ants to the non-conflicting paths
	distribution := path.DistributeAntsToPaths(nonConflictingPaths, len(field.Ants))

	// Print the input file
	helpers.PrintLines(lines)

	// Simulate ant movements
	solver.SimulateAntMovements(field, nonConflictingPaths, distribution)
}

