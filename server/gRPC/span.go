package gRPC

import (
	"context"
	"encoding/json"
	"fmt"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	proto "github.com/maticnetwork/polyproto/heimdall"
	protoutils "github.com/maticnetwork/polyproto/utils"
)

func (h *HeimdallGRPCServer) Span(ctx context.Context, in *proto.SpanRequest) (*proto.SpanResponse, error) {
	cliCtx := cliContext.NewCLIContext().WithCodec(h.cdc)
	result, err := helper.FetchFromAPI(cliCtx, helper.GetHeimdallServerEndpoint(fmt.Sprintf(spanURL, in.ID)))

	if err != nil {
		logger.Error("Error while fetching span")
		return nil, err
	}

	resp := &proto.SpanResponse{}
	resp.Result = parseSpan(result.Result)
	resp.Height = fmt.Sprint(result.Height)

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

	resp := &proto.Span{
		ID:                span.ID,
		StartBlock:        span.StartBlock,
		EndBlock:          span.EndBlock,
		ChainID:           span.ChainID,
		ValidatorSet:      &proto.ValidatorSet{},
		SelectedProducers: []*proto.Validator{},
	}

	for _, v := range span.ValidatorSet.Validators {
		copy(addr[:], v.Signer.Bytes())
		resp.ValidatorSet.Validators = append(resp.ValidatorSet.Validators, parseValidator(addr, v))
	}

	copy(addr[:], span.ValidatorSet.Proposer.Signer.Bytes())
	resp.ValidatorSet.Proposer = parseValidator(addr, span.ValidatorSet.Proposer)

	for i := range span.SelectedProducers {
		copy(addr[:], span.SelectedProducers[i].Signer.Bytes())
		resp.SelectedProducers = append(resp.SelectedProducers, parseValidator(addr, &span.SelectedProducers[i]))
	}

	return resp
}

func parseValidator(address [20]byte, validator *types.Validator) *proto.Validator {
	return &proto.Validator{
		ID:               uint64(validator.ID),
		Address:          protoutils.ConvertAddressToH160(address),
		VotingPower:      validator.VotingPower,
		ProposerPriority: validator.ProposerPriority,
	}
}
