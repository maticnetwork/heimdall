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

type Checkpoint struct {
	Proposer   hmTypes.HeimdallAddress `json:"proposer"`
	StartBlock uint64                  `json:"start_block"`
	EndBlock   uint64                  `json:"end_block"`
	RootHash   hmTypes.HeimdallHash    `json:"root_hash"`
	BorChainID string                  `json:"bor_chain_id"`
	TimeStamp  uint64                  `json:"timestamp"`
}

func (h *HeimdallGRPCServer) FetchCheckpointCount(ctx context.Context, in *emptypb.Empty) (*proto.FetchCheckpointCountResponse, error) {
	cliCtx := cliContext.NewCLIContext().WithCodec(h.cdc)

	result, err := helper.FetchFromAPI(cliCtx, helper.GetHeimdallServerEndpoint(fetchCheckpointCount))
	if err != nil {
		logger.Error("Error while fetching checkpoint count")
		return nil, err
	}

	resp := &proto.FetchCheckpointCountResponse{}
	resp.Height = fmt.Sprint(result.Height)

	if err := json.Unmarshal(result.Result, &resp.Result); err != nil {
		logger.Error("Error unmarshalling checkpoint count", "error", err)
		return nil, err
	}

	return resp, nil
}

func (h *HeimdallGRPCServer) FetchCheckpoint(ctx context.Context, in *proto.FetchCheckpointRequest) (*proto.FetchCheckpointResponse, error) {
	cliCtx := cliContext.NewCLIContext().WithCodec(h.cdc)

	url := ""
	if in.ID == -1 {
		url = fmt.Sprintf(fetchCheckpoint, "latest")
	} else {
		url = fmt.Sprintf(fetchCheckpoint, fmt.Sprint(in.ID))
	}

	result, err := helper.FetchFromAPI(cliCtx, helper.GetHeimdallServerEndpoint(url))

	if err != nil {
		logger.Error("Error while fetching checkpoint")
		return nil, err
	}

	checkPoint := &Checkpoint{}
	if err := json.Unmarshal(result.Result, &checkPoint); err != nil {
		logger.Error("Error unmarshalling checkpoint", "error", err)
		return nil, err
	}

	var hash [32]byte

	copy(hash[:], checkPoint.RootHash.Bytes())

	var address [20]byte

	copy(address[:], checkPoint.Proposer.Bytes())

	resp := &proto.FetchCheckpointResponse{}
	resp.Height = fmt.Sprint(result.Height)
	resp.Result = &proto.Checkpoint{
		StartBlock: checkPoint.StartBlock,
		EndBlock:   checkPoint.EndBlock,
		RootHash:   protoutils.ConvertHashToH256(hash),
		Proposer:   protoutils.ConvertAddressToH160(address),
		Timestamp:  timestamppb.New(time.Unix(int64(checkPoint.TimeStamp), 0)),
		BorChainID: checkPoint.BorChainID,
	}

	return resp, nil
}
