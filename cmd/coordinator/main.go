package main

import (
	"github.com/priyansh32/dkvs/internal/coordinator"
)

func main() {
	err := coordinator.InitCoordinator(3) // Blocking call
	if err != nil {
		panic(err)
	}
}
