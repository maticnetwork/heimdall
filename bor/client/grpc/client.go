package grpc

import (
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/log"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	proto "github.com/maticnetwork/polyproto/bor"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type BorGRPCClient struct {
	conn   *grpc.ClientConn
	client proto.BorApiClient
}

func NewBorGRPCClient(address string) *BorGRPCClient {
	address = removePrefix(address)

	opts := []grpc_retry.CallOption{
		grpc_retry.WithMax(5),
		grpc_retry.WithBackoff(grpc_retry.BackoffLinear(1 * time.Second)),
		grpc_retry.WithCodes(codes.Internal, codes.Unavailable, codes.Aborted, codes.NotFound),
	}

	conn, err := grpc.NewClient(address,
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(opts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(opts...)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Crit("Failed to connect to Bor gRPC", "error", err)
	}

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

// removePrefix removes the http:// or https:// prefix from the address, if present.
func removePrefix(address string) string {
	if strings.HasPrefix(address, "http://") || strings.HasPrefix(address, "https://") {
		return address[strings.Index(address, "//")+2:]
	}
	return address
}
