package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/priyansh32/dkvs/internal/coordinator"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("Usage: coordinator <port> <replication_factor>")
		return
	}

	coordinatorPort := os.Args[1]
	replicationFactor, err := strconv.Atoi(os.Args[2])
	if err != nil {
		panic(err)
	}

	// parse the replication factor

	err = coordinator.InitCoordinator(coordinatorPort, replicationFactor) // Blocking call
	if err != nil {
		panic(err)
	}
}
