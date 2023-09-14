package store

import (
	"context"
	"log"
	"net"

	pb "github.com/priyansh32/dkvs/internal/api/store"
	"google.golang.org/grpc"
)

const DEFAULT_CAPACITY = 7

type store struct {
	cache *LRUCache
	pb.UnimplementedKeyValueStoreServer
}

func NewStore(capacity uint32) *store {
	return &store{
		cache: LRUConstructor(capacity),
	}
}

func (s *store) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {

	// get the value from the cache
	value, err := s.cache.Get(in.Key)
	log.Printf("Got key: %s, value: %s\n", in.Key, value)

	if err != nil {
		return &pb.GetResponse{Status: pb.StatusType_CACHE_MISS, Value: value}, nil
	}

	return &pb.GetResponse{Status: pb.StatusType_OK, Value: value}, nil
}

func (s *store) Put(ctx context.Context, in *pb.PutRequest) (*pb.PutResponse, error) {

	// put the value in the cache
	s.cache.Put(in.Key, in.Value)
	log.Printf("Cached key: %s, value: %s\n", in.Key, in.Value)
	return &pb.PutResponse{Status: pb.StatusType_OK}, nil
}

func (s *store) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {

	// delete the value from the cache
	s.cache.Remove(in.Key)
	log.Printf("Deleted key: %s\n", in.Key)
	return &pb.DeleteResponse{Status: pb.StatusType_OK}, nil
}

func StartStoreServer(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return err
	}

	gRPCServer := grpc.NewServer()
	pb.RegisterKeyValueStoreServer(gRPCServer, NewStore(DEFAULT_CAPACITY))

	log.Printf("starting gRPC server on %s", address)

	if err := gRPCServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
		return err
	}

	return nil
}
