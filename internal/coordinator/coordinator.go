package coordinator

import (
	"context"
	"errors"
	"log"
	"net"

	pb_coordinator "github.com/priyansh32/nebula/internal/api/coordinator"
	pb_store "github.com/priyansh32/nebula/internal/api/store"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StoreClient struct {
	name     string
	conn     *grpc.ClientConn
	client   pb_store.KeyValueStoreClient
	nodeKeys []uint64
}

// Returns a new store client with the given connection and name
func NewStoreClient(conn *grpc.ClientConn, name string) *StoreClient {
	return &StoreClient{
		name:   name,
		conn:   conn,
		client: pb_store.NewKeyValueStoreClient(conn),
	}
}

type Coordinator struct {
	ctx          context.Context
	hashRing     *HashRing
	storeClients map[string]*StoreClient
	pb_coordinator.UnimplementedCoordinatorAPIServer
}

// returns a new coordinator with the given replication factor
func NewCoordinator(replicationFactor int) (*Coordinator, error) {

	if replicationFactor < 1 {
		return nil, errors.New("replication factor must be greater than 0")
	}

	if replicationFactor >= 1024 {
		return nil, errors.New("replication factor must be less than 1024")
	}

	return &Coordinator{
		ctx:          context.Background(),
		hashRing:     NewHashRing(replicationFactor),
		storeClients: make(map[string]*StoreClient),
	}, nil
}

func (c *Coordinator) Get(ctx context.Context, in *pb_coordinator.GetRequest) (*pb_coordinator.GetResponse, error) {
	key := in.Key

	store, err := c.hashRing.GetStore(key)
	if err != nil {
		return nil, err
	}

	res, err := store.client.Get(c.ctx, &pb_store.GetRequest{Key: key})
	if err != nil {
		return nil, err
	}

	return &pb_coordinator.GetResponse{
		Status: pb_coordinator.StatusType(res.Status),
		Value:  res.Value,
	}, nil
}

// Put puts the value for the given key in the appropriate store
func (c *Coordinator) Put(ctx context.Context, in *pb_coordinator.PutRequest) (*pb_coordinator.PutResponse, error) {

	key := in.Key
	value := in.Value

	store, err := c.hashRing.GetStore(key)
	if err != nil {
		return nil, err
	}

	_, err = store.client.Put(c.ctx, &pb_store.PutRequest{Key: key, Value: value})
	if err != nil {
		return nil, err
	}

	return &pb_coordinator.PutResponse{
		Status: pb_coordinator.StatusType_OK,
	}, nil
}

// Delete deletes the value for the given key from the appropriate store
func (c *Coordinator) Delete(ctx context.Context, in *pb_coordinator.DeleteRequest) (*pb_coordinator.DeleteResponse, error) {

	key := in.Key

	store, err := c.hashRing.GetStore(key)
	if err != nil {
		return nil, err
	}

	_, err = store.client.Delete(c.ctx, &pb_store.DeleteRequest{Key: key})
	if err != nil {
		return nil, err
	}

	return &pb_coordinator.DeleteResponse{
		Status: pb_coordinator.StatusType_OK,
	}, nil
}

// adds nodes of a store to the hash ring
func (c *Coordinator) AddStore(ctx context.Context, in *pb_coordinator.AddStoreRequest) (*pb_coordinator.AddStoreResponse, error) {

	address := in.Address
	name := in.Name

	// check if the store already exists with the given name
	if _, ok := c.storeClients[name]; ok {
		return nil, errors.New("duplicate store name")
	}

	// address is of the form <host>:<port>
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	storeClient := NewStoreClient(conn, name)
	c.storeClients[name] = storeClient
	c.hashRing.AddStoreNodes(storeClient)

	return &pb_coordinator.AddStoreResponse{
		Status: pb_coordinator.StatusType_OK,
	}, nil
}

// removes nodes of a store from the hash ring
func (c *Coordinator) RemoveStore(ctx context.Context, in *pb_coordinator.RemoveStoreRequest) (*pb_coordinator.RemoveStoreResponse, error) {
	s := c.storeClients[in.Name]

	if s == nil {
		return nil, errors.New("store " + in.Name + " does not exist")
	}

	// remove the store from the store clients
	delete(c.storeClients, in.Name)

	// close the connection
	s.conn.Close()

	// remove the nodes from the hash ring
	for _, key := range s.nodeKeys {
		delete(c.hashRing.nodes, key)
		c.hashRing.sortedKeys = removeSorted(c.hashRing.sortedKeys, key)
	}

	return &pb_coordinator.RemoveStoreResponse{
		Status: pb_coordinator.StatusType_OK,
	}, nil
}

func InitCoordinator(port string, rf int) error {
	cdr, err := NewCoordinator(rf)
	if err != nil {
		return err
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	log.Printf("coordinator listening on port: %d\n", lis.Addr().(*net.TCPAddr).Port)

	gRPCServer := grpc.NewServer()
	pb_coordinator.RegisterCoordinatorAPIServer(gRPCServer, cdr)

	if err := gRPCServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

	return nil
}
