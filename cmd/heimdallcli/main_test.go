package main

import (
	"encoding/json"
	"os"
	"strconv"
	"testing"

	heimdallApp "github.com/maticnetwork/heimdall/app"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmTypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"
)

func TestModulesStreamedGenesisExport(t *testing.T) {
	t.Parallel()

	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
	db := dbm.NewMemDB()

	happ := heimdallApp.NewHeimdallApp(logger, db)

	genDoc, err := tmTypes.GenesisDocFromFile("./testdata/dump-genesis.json")
	require.NoError(t, err)

	var genesisState heimdallApp.GenesisState
	err = json.Unmarshal(genDoc.AppState, &genesisState)
	require.NoError(t, err)

	ctx := happ.NewContext(true, abci.Header{Height: 1})
	happ.GetModuleManager().InitGenesis(ctx, genesisState)

	marshaledAppState, err := generateMarshalledAppState(happ, 2)
	require.NoError(t, err)

	unmarshaledAppState := map[string]interface{}{}
	err = json.Unmarshal(marshaledAppState, &unmarshaledAppState)
	require.NoError(t, err)

	clerk, err := traversePath(unmarshaledAppState, "clerk")
	require.NoError(t, err)
	eventRecords := clerk["event_records"].([]interface{})
	require.Len(t, eventRecords, 6)
	for idx, record := range eventRecords {
		eventRecord := record.(map[string]interface{})
		require.NotEmpty(t, eventRecord["id"])
		eventIDStr := eventRecord["id"].(string)
		eventID, err := strconv.Atoi(eventIDStr)
		require.NoError(t, err)
		require.Equal(t, idx+1, eventID)
		require.NotEmpty(t, eventRecord["contract"])
		require.NotEmpty(t, eventRecord["data"])
		require.NotEmpty(t, eventRecord["tx_hash"])
		require.NotEmpty(t, eventRecord["log_index"])
		require.NotEmpty(t, eventRecord["bor_chain_id"])
		require.NotEmpty(t, eventRecord["record_time"])
	}

	bor, err := traversePath(unmarshaledAppState, "bor")
	require.NoError(t, err)
	spans := bor["spans"].([]interface{})
	require.Len(t, spans, 5)
	for idx, span := range spans {
		spanMap := span.(map[string]interface{})
		require.NotEmpty(t, spanMap["span_id"])
		spanIDStr := spanMap["span_id"].(string)
		spanID, err := strconv.Atoi(spanIDStr)
		require.NoError(t, err)
		require.Equal(t, idx, spanID)
		require.NotEmpty(t, spanMap["start_block"])
		require.NotEmpty(t, spanMap["end_block"])
		require.NotEmpty(t, spanMap["validator_set"])
		require.NotEmpty(t, spanMap["selected_producers"])
		require.NotEmpty(t, spanMap["bor_chain_id"])
	}
}
