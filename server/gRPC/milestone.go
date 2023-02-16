package gRPC

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"

	proto "github.com/maticnetwork/polyproto/heimdall"
	protoutils "github.com/maticnetwork/polyproto/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Milestone struct {
	Proposer   hmTypes.HeimdallAddress `json:"proposer"`
	StartBlock uint64                  `json:"start_block"`
	EndBlock   uint64                  `json:"end_block"`
	Hash       hmTypes.HeimdallHash    `json:"hash"`
	BorChainID string                  `json:"bor_chain_id"`
	TimeStamp  uint64                  `json:"timestamp"`
}

func (h *HeimdallGRPCServer) FetchMilestoneCount(ctx context.Context, in *emptypb.Empty) (*proto.FetchMilestoneCountResponse, error) {
	cliCtx := cliContext.NewCLIContext().WithCodec(h.cdc)

	result, err := helper.FetchFromAPI(cliCtx, helper.GetHeimdallServerEndpoint(fetchMilestoneCount))
	if err != nil {
		logger.Error("Error while fetching checkpoint count")
		return nil, err
	}

	resp := &proto.FetchMilestoneCountResponse{}
	resp.Height = fmt.Sprint(result.Height)

	if err := json.Unmarshal(result.Result, &resp.Result); err != nil {
		logger.Error("Error unmarshalling milestone count", "error", err)
		return nil, err
	}

	return resp, nil
}

func (h *HeimdallGRPCServer) FetchMilestone(ctx context.Context, in *emptypb.Empty) (*proto.FetchMilestoneResponse, error) {
	cliCtx := cliContext.NewCLIContext().WithCodec(h.cdc)

	result, err := helper.FetchFromAPI(cliCtx, helper.GetHeimdallServerEndpoint(fetchMilestone))

	if err != nil {
		logger.Error("Error while fetching milestone")
		return nil, err
	}

	milestone := &Milestone{}
	if err := json.Unmarshal(result.Result, &milestone); err != nil {
		logger.Error("Error unmarshalling milestone", "error", err)
		return nil, err
	}

	var hash [32]byte

	copy(hash[:], milestone.Hash.Bytes())

	var address [20]byte

	copy(address[:], milestone.Proposer.Bytes())

	resp := &proto.FetchMilestoneResponse{}
	resp.Height = fmt.Sprint(result.Height)
	resp.Result = &proto.Milestone{
		StartBlock: milestone.StartBlock,
		EndBlock:   milestone.EndBlock,
		RootHash:   protoutils.ConvertHashToH256(hash),
		Proposer:   protoutils.ConvertAddressToH160(address),
		Timestamp:  timestamppb.New(time.Unix(int64(milestone.TimeStamp), 0)),
		BorChainID: milestone.BorChainID,
	}

	return resp, nil
}

func (h *HeimdallGRPCServer) FetchLastNoAckMilestone(ctx context.Context, in *emptypb.Empty) (*proto.FetchLastNoAckMilestoneResponse, error) {
	cliCtx := cliContext.NewCLIContext().WithCodec(h.cdc)

	result, err := helper.FetchFromAPI(cliCtx, helper.GetHeimdallServerEndpoint(fetchLastNoAckMilestone))

	if err != nil {
		logger.Error("Error while fetching milestone last no ack")
		return nil, err
	}

	resp := &proto.FetchLastNoAckMilestoneResponse{}
	resp.Height = fmt.Sprint(result.Height)

	if err := json.Unmarshal(result.Result, &resp.Result); err != nil {
		logger.Error("Error unmarshalling milestone", "error", err)
		return nil, err
	}

	return resp, nil
}

func (h *HeimdallGRPCServer) FetchNoAckMilestone(ctx context.Context, in *proto.FetchMilestoneNoAckRequest) (*proto.FetchMilestoneNoAckResponse, error) {
	cliCtx := cliContext.NewCLIContext().WithCodec(h.cdc)

	url := fmt.Sprintf(fetchMilestoneNoAck, fmt.Sprint(in.MilestoneID))

	result, err := helper.FetchFromAPI(cliCtx, helper.GetHeimdallServerEndpoint(url))

	if err != nil {
		logger.Error("Error while fetching milestone no ack")
		return nil, err
	}

	resp := &proto.FetchMilestoneNoAckResponse{}
	resp.Height = fmt.Sprint(result.Height)

	if err := json.Unmarshal(result.Result, &resp.Result); err != nil {
		logger.Error("Error unmarshalling milestone", "error", err)
		return nil, err
	}

	return resp, nil
}

func (h *HeimdallGRPCServer) FetchMilestoneID(ctx context.Context, in *proto.FetchMilestoneIDRequest) (*proto.FetchMilestoneIDResponse, error) {
	cliCtx := cliContext.NewCLIContext().WithCodec(h.cdc)

	url := fmt.Sprintf(fetchMilestoneID, fmt.Sprint(in.MilestoneID))

	result, err := helper.FetchFromAPI(cliCtx, helper.GetHeimdallServerEndpoint(url))

	if err != nil {
		logger.Error("Error while fetching milestone id")
		return nil, err
	}

	resp := &proto.FetchMilestoneIDResponse{}
	resp.Height = fmt.Sprint(result.Height)

	if err := json.Unmarshal(result.Result, &resp.Result); err != nil {
		logger.Error("Error unmarshalling milestone", "error", err)
		return nil, err
	}

	return resp, nil
}
