package configurator

import (
	"context"
	"fmt"
	"reflect"

	sdk "github.com/cosmos/cosmos-sdk/types"
	gogogrpc "github.com/gogo/protobuf/grpc"
	"google.golang.org/grpc"
)

type testRouter struct {
	commitWrites bool
	handlers     map[string]func(ctx context.Context, args, reply interface{}) error
}

func newTestRouter(commitWrites bool) *testRouter {
	return &testRouter{
		commitWrites: commitWrites,
		handlers:     map[string]func(ctx context.Context, args interface{}, reply interface{}) error{},
	}
}

var _ gogogrpc.Server = testRouter{}
var _ grpc.ClientConnInterface = testRouter{}

func (t testRouter) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	for _, method := range sd.Methods {
		fqName := fmt.Sprintf("/%s/%s", sd.ServiceName, method.MethodName)
		handler := method.Handler
		t.handlers[fqName] = func(ctx context.Context, args, reply interface{}) error {
			res, err := handler(ss, ctx, func(i interface{}) error { return nil },
				func(ctx context.Context, _ interface{}, _ *grpc.UnaryServerInfo, unaryHandler grpc.UnaryHandler) (resp interface{}, err error) {
					return unaryHandler(ctx, args)
				})
			if err != nil {
				return err
			}

			resValue := reflect.ValueOf(res)
			if !resValue.IsZero() {
				reflect.ValueOf(reply).Elem().Set(resValue.Elem())
			}
			return nil
		}
	}
}

func (t testRouter) Invoke(ctx context.Context, method string, args, reply interface{}, _ ...grpc.CallOption) error {
	handler := t.handlers[method]
	if handler == nil {
		return fmt.Errorf("can't find handler for method %s", method)
	}

	// cache wrap the multistore so that writes are batched
	sdkCtx := sdk.UnwrapSDKContext(ctx)
	ms := sdkCtx.MultiStore()
	cacheMs := ms.CacheMultiStore()
	sdkCtx = sdkCtx.WithMultiStore(cacheMs)
	ctx = sdk.WrapSDKContext(sdkCtx)

	err := handler(ctx, args, reply)
	if err != nil {
		return err
	}

	// only commit writes if there are no errors and commitWrites is true
	if t.commitWrites {
		cacheMs.Write()
	}

	return nil
}

func (t testRouter) NewStream(context.Context, *grpc.StreamDesc, string, ...grpc.CallOption) (grpc.ClientStream, error) {
	return nil, fmt.Errorf("unsupported")
}
