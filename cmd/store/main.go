package main

import (
	"fmt"
	"os"

	"github.com/priyansh32/nebula/internal/store"
)

func main() {

	args := os.Args[1:]

	if len(args) != 1 {
		fmt.Println("Usage: store <port>")
		return
	}

	err := store.StartStoreServer(":" + args[0])
	if err != nil {
		panic(err)
	}

	fmt.Println("Server started successfully")
}
