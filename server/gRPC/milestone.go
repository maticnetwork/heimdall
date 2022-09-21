package gRPC

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/maticnetwork/heimdall/helper"

	proto "github.com/maticnetwork/polyproto/heimdall"
	protoutils "github.com/maticnetwork/polyproto/utils"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Milestone Checkpoint

func (h *HeimdallGRPCServer) FetchMilestoneCount(ctx context.Context, in *emptypb.Empty) (*proto.FetchMilestoneCountResponse, error) {
	cliCtx := cliContext.NewCLIContext().WithCodec(h.cdc)

	result, err := helper.FetchFromAPI(cliCtx, helper.GetHeimdallServerEndpoint(fetchCheckpointCount))
	if err != nil {
		logger.Error("Error while fetching milestone count")
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

	copy(hash[:], milestone.RootHash.Bytes())

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
