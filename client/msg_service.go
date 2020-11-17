package client

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gogogrpc "github.com/gogo/protobuf/grpc"
	"google.golang.org/grpc"
)

// This file is copy-pasted from github.com/cosmos/cosmos-sdk/blob/master/x/bank/client/testutil/cli_helpers.go.
// This is pending https://github.com/cosmos/cosmos-sdk/pull/7630.

// ServiceMsgClientConn is an instance of grpc.ClientConn that is used to test building
// transactions with MsgClient's. It is intended to be replaced by the work in
// https://github.com/cosmos/cosmos-sdk/issues/7541 when that is ready.
type ServiceMsgClientConn struct {
	Msgs []sdk.Msg
}

var _ gogogrpc.ClientConn = &ServiceMsgClientConn{}

func (t *ServiceMsgClientConn) Invoke(_ context.Context, method string, args, _ interface{}, _ ...grpc.CallOption) error {
	req, ok := args.(sdk.MsgRequest)
	if !ok {
		return fmt.Errorf("%T should implement %T", args, (*sdk.MsgRequest)(nil))
	}

	err := req.ValidateBasic()
	if err != nil {
		return err
	}

	t.Msgs = append(t.Msgs, sdk.ServiceMsg{
		MethodName: method,
		Request:    req,
	})

	return nil
}

func (t *ServiceMsgClientConn) NewStream(context.Context, *grpc.StreamDesc, string, ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, fmt.Errorf("not supported")
}
