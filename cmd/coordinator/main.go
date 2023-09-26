package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/priyansh32/nebula/internal/coordinator"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: coordinator <port> <replication-factor>")
		os.Exit(1)
	}

	coordinatorPort := os.Args[1]
	replicationFactor, err := strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Println("Replication factor must be an integer")
		os.Exit(1)
	}

	coordinator.InitCoordinator(coordinatorPort, replicationFactor) // Blocking call
}
