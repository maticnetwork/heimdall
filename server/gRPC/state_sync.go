package gRPC

import (
	"encoding/json"
	"fmt"
	"net/url"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/server/gRPC/proto"
)

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
			return err
		}

		var eventRecords []*proto.EventRecord

		err = json.Unmarshal(result.Result, &eventRecords)
		if err != nil {
			logger.Error("Error unmarshalling event record", "error", err)
			return nil
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
			return err
		}

		fromId += req.Limit
	}

	return nil
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
