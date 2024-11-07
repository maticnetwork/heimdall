package grpc

import (
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/log"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	proto "github.com/maticnetwork/polyproto/bor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	stateFetchLimit = 50
)

type BorGRPCClient struct {
	conn   *grpc.ClientConn
	client proto.BorApiClient
}

func NewBorGRPCClient(address string) *BorGRPCClient {
	opts := []grpc_retry.CallOption{
		grpc_retry.WithMax(10000),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(5 * time.Second)),
		grpc_retry.WithCodes(codes.Internal, codes.Unavailable, codes.Aborted, codes.NotFound),
	}

	fmt.Printf(">>>>> Connecting to Bor gRPC server at %s\n", address)

	conn, err := grpc.Dial(address,
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		fmt.Printf(">>>>> Error connecting to Bor gRPC\n")
		log.Crit("Failed to connect to Bor gRPC", "error", err)
	}

	fmt.Printf(">>>>> Connected to Bor gRPC\n")
	log.Info("Connected to Bor gRPC server", "address", address)

	return &BorGRPCClient{
		conn:   conn,
		client: proto.NewBorApiClient(conn),
	}
}

func (h *BorGRPCClient) Close() {
	log.Debug("Shutdown detected, Closing Bor gRPC client")
	h.conn.Close()
}
