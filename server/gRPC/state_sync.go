package gRPC

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/maticnetwork/heimdall/helper"
	proto "github.com/maticnetwork/polyproto/heimdall"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Event struct {
	Id         uint64 `json:"id"`
	Contract   string `json:"contract"`
	Data       string `json:"data"`
	TxHash     string `json:"tx_hash"`
	LogIndex   uint64 `json:"log_index"`
	BorChainId string `json:"bor_chain_id"`
	RecordTime string `json:"record_time"`
}

func (h *HeimdallGRPCServer) StateSyncEvents(req *proto.StateSyncEventsRequest, reply proto.Heimdall_StateSyncEventsServer) error {
	cliCtx := cliContext.NewCLIContext().WithCodec(h.cdc)
	fromId := req.FromID

	for {
		params := map[string]string{
			"from-id": fmt.Sprint(fromId),
			"to-time": fmt.Sprint(req.ToTime),
			"limit":   fmt.Sprint(req.Limit),
		}

		result, err := helper.FetchFromAPI(cliCtx, addParamsToEndpoint(helper.GetHeimdallServerEndpoint(eventRecordList), params))
		if err != nil {
			logger.Error("Error while fetching event records", "error", err)
			return status.Errorf(codes.Internal, err.Error())
		}

		eventRecords, err := parseEvents(result.Result)
		if err != nil {
			logger.Error("Error while parsing event records", "error", err)
			return status.Errorf(codes.Internal, err.Error())
		}

		if len(eventRecords) == 0 {
			break
		}

		err = reply.Send(&proto.StateSyncEventsResponse{
			Height: fmt.Sprint(result.Height),
			Result: eventRecords,
		})
		if err != nil {
			logger.Error("Error while sending event record", "error", err)
			return status.Errorf(codes.Internal, err.Error())
		}

		fromId += req.Limit
	}

	return nil
}

func parseEvents(result json.RawMessage) ([]*proto.EventRecord, error) {
	var events []Event

	err := json.Unmarshal(result, &events)
	if err != nil {
		logger.Error("Error unmarshalling event record", "error", err)
		return nil, err
	}

	eventRecords := make([]*proto.EventRecord, len(events))

	for i, event := range events {
		eventTime, err := time.Parse(time.RFC3339, event.RecordTime)
		if err != nil {
			logger.Error("Error parsing time", "error", err)
			return nil, err
		}

		eventRecords[i] = &proto.EventRecord{
			ID:       event.Id,
			Contract: event.Contract,
			Data:     event.Data,
			TxHash:   event.TxHash,
			LogIndex: event.LogIndex,
			ChainID:  event.BorChainId,
			Time:     timestamppb.New(eventTime),
		}
	}

	return eventRecords, nil
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
