package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/url"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/server/proto"
	"github.com/maticnetwork/heimdall/types"
	tmLog "github.com/tendermint/tendermint/libs/log"
	"google.golang.org/grpc"
)

const (
	LatestSpanURL   = "/bor/latest-span"
	SpanURL         = "/bor/span/%v"
	EventRecordList = "/clerk/event-record/list"
)

var logger tmLog.Logger

func setupGRPCServer(shutDownCtx context.Context, cdc *codec.Codec, addr string, lggr tmLog.Logger) error {
	logger = lggr
	grpcServer := grpc.NewServer()
	proto.RegisterHeimdallServer(grpcServer,
		&heimdallGRPCServer{
			cdc: cdc,
		})

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error("failed to serve grpc server", "err", err)
		}

		<-shutDownCtx.Done()
		grpcServer.Stop()
		lis.Close()
	}()

	logger.Info("GRPC Server started", "addr", addr)

	return nil
}

type heimdallGRPCServer struct {
	proto.UnimplementedHeimdallServer
	cdc *codec.Codec
}

func (h *heimdallGRPCServer) GetSpan(ctx context.Context, in *proto.GetSpanRequest) (*proto.GetSpanResponse, error) {
	cliCtx := cliContext.NewCLIContext().WithCodec(h.cdc)
	result, err := helper.FetchFromAPI(cliCtx, helper.GetHeimdallServerEndpoint(fmt.Sprintf(SpanURL, in.ID)))

	if err != nil {
		logger.Error("Error while fetching span")
		return nil, err
	}

	resp := &proto.GetSpanResponse{}
	resp.Result = parseSpan(result.Result)
	resp.Height = result.Height

	return resp, nil
}

func parseSpan(result json.RawMessage) *proto.Span {
	var addr [20]byte

	span := &types.Span{}

	err := json.Unmarshal(result, span)
	if err != nil {
		logger.Error("Error unmarshalling span", "error", err)
		return nil
	}

	resp := &proto.Span{}
	resp.ID = span.ID
	resp.StartBlock = span.StartBlock
	resp.EndBlock = span.EndBlock
	resp.ChainID = span.ChainID
	resp.ValidatorSet = &proto.ValidatorSet{}
	resp.SelectedProducers = []*proto.Validator{}

	for _, v := range span.ValidatorSet.Validators {
		copy(addr[:], v.Signer.Bytes())

		resp.ValidatorSet.Validators = append(resp.ValidatorSet.Validators, &proto.Validator{
			ID:               uint64(v.ID),
			Address:          proto.ConvertAddressToH160(addr),
			VotingPower:      v.VotingPower,
			ProposerPriority: v.ProposerPriority,
		})
	}

	copy(addr[:], span.ValidatorSet.Proposer.Signer.Bytes())

	resp.ValidatorSet.Proposer = &proto.Validator{
		ID:               uint64(span.ValidatorSet.Proposer.ID),
		Address:          proto.ConvertAddressToH160(addr),
		VotingPower:      span.ValidatorSet.Proposer.VotingPower,
		ProposerPriority: span.ValidatorSet.Proposer.ProposerPriority,
	}

	for _, v := range span.SelectedProducers {
		copy(addr[:], v.Signer.Bytes())

		resp.SelectedProducers = append(resp.SelectedProducers, &proto.Validator{
			ID:               uint64(v.ID),
			Address:          proto.ConvertAddressToH160(addr),
			VotingPower:      v.VotingPower,
			ProposerPriority: v.ProposerPriority,
		})
	}

	return resp
}

func (h *heimdallGRPCServer) GetEventRecords(req *proto.GetEventRecordsRequest, reply proto.Heimdall_GetEventRecordsServer) error {
	cliCtx := cliContext.NewCLIContext().WithCodec(h.cdc)
	fromId := req.FromID

	for {
		params := map[string]string{
			"from-id": fmt.Sprintf("%d", fromId),
			"to-time": fmt.Sprintf("%d", req.ToTime),
			"limit":   fmt.Sprintf("%d", req.Limit),
		}

		result, err := helper.FetchFromAPI(cliCtx, addParamsToEndpoint(helper.GetHeimdallServerEndpoint(EventRecordList), params))
		if err != nil {
			logger.Error("Error while fetching latest span")
			return err
		}

		eventRecords := parseEventRecords(result.Result)

		if len(eventRecords) == 0 {
			break
		}

		err = reply.Send(&proto.GetEventRecordsResponse{
			Height: result.Height,
			Result: eventRecords,
		})
		if err != nil {
			logger.Error("Error while sending event record", "error", err)
			return err
		}

		fromId += req.Limit
	}

	return nil
}

func parseEventRecords(result json.RawMessage) []*proto.EventRecord {
	resp := []*proto.EventRecord{}
	err := json.Unmarshal(result, &resp)

	if err != nil {
		logger.Error("Error unmarshalling event record", "error", err)
		return nil
	}

	return resp
}

func addParamsToEndpoint(endpoint string, params map[string]string) string {
	u, _ := url.Parse(endpoint)
	q := u.Query()

	for k, v := range params {
		q.Set(k, v)
	}

	u.RawQuery = q.Encode()

	return u.String()
}
