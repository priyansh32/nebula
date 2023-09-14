package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	pb "github.com/priyansh32/dkvs/internal/api/coordinator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: cli <coordinator-address>")
		return
	}

	address := os.Args[1]

	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	client := pb.NewCoordinatorAPIClient(conn)

	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("NEBULA> ")
		scanner.Scan()
		command := scanner.Text()
		makerequest(command, client)
	}

}

func makerequest(command string, client pb.CoordinatorAPIClient) {
	// split the command into tokens
	tokens := strings.Split(command, " ")

	// check if the command is valid
	if len(tokens) < 1 {
		fmt.Println("Invalid command")
		return
	}

	// check if the command is a valid command
	switch tokens[0] {
	case "ADDSTORE":
		// check if the command has the correct number of arguments
		if len(tokens) != 3 {
			fmt.Println("Missing arguments: ADDSTORE <address> <name>")
			return
		}
		res, err := client.AddStore(context.Background(), &pb.AddStoreRequest{
			Address: tokens[1],
			Name:    tokens[2],
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Added store: ", tokens[2], " with status: ", res.Status)
	case "PUT":
		// check if the command has the correct number of arguments
		if len(tokens) != 3 {
			fmt.Println("Missing arguments: PUT <key> <value>")
			return
		}
		res, err := client.Put(context.Background(), &pb.PutRequest{
			Key:   tokens[1],
			Value: tokens[2],
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Status: ", res.Status)
	case "GET":
		// check if the command has the correct number of arguments
		if len(tokens) != 2 {
			fmt.Println("Missing arguments: GET <key>")
			return
		}
		res, err := client.Get(context.Background(), &pb.GetRequest{
			Key: tokens[1],
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		if res.Status == pb.StatusType_OK {
			fmt.Println("Value: ", res.Value)
		} else {
			fmt.Println("CACHE MISS")
		}
	case "DELETE":
		// check if the command has the correct number of arguments
		if len(tokens) != 2 {
			fmt.Println("Missing arguments: DELETE <key>")
			return
		}
		res, err := client.Delete(context.Background(), &pb.DeleteRequest{
			Key: tokens[1],
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Status: ", res.Status)
	case "REMOVESTORE":
		// check if the command has the correct number of arguments
		if len(tokens) != 2 {
			fmt.Println("Missing arguments: REMOVESTORE <name>")
			return
		}
		res, err := client.RemoveStore(context.Background(), &pb.RemoveStoreRequest{
			Name: tokens[1],
		})
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Status: ", res.Status)
	case "EXIT":
		os.Exit(0)
	default:
		fmt.Println("Invalid command")
	}

}
