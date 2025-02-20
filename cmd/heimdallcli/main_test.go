//nolint:govet
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
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

	var buf bytes.Buffer

	err = generateMarshalledAppState(happ, "test-chain", 2, &buf)
	require.NoError(t, err)

	marshaledAppState := buf.Bytes()

	var unmarshaledAppState map[string]interface{}
	err = json.Unmarshal(marshaledAppState, &unmarshaledAppState)
	require.NoError(t, err)

	appState, ok := unmarshaledAppState["app_state"].(map[string]interface{})
	require.True(t, ok, "app_state should be a map")

	clerk, err := traversePath(appState, "clerk")
	require.NoError(t, err)

	eventRecords, ok := clerk["event_records"].([]interface{})
	require.True(t, ok, "event_records should be an array")
	require.Len(t, eventRecords, 6)
	for idx, record := range eventRecords {
		eventRecord, ok := record.(map[string]interface{})
		require.True(t, ok, "eventRecord should be a map")
		require.NotEmpty(t, eventRecord["id"])
		eventIDStr, ok := eventRecord["id"].(string)
		require.True(t, ok, "id should be a string")
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

	bor, err := traversePath(appState, "bor")
	require.NoError(t, err)

	spans, ok := bor["spans"].([]interface{})
	require.True(t, ok, "spans should be an array")
	require.Len(t, spans, 5)
	for idx, span := range spans {
		spanMap, ok := span.(map[string]interface{})
		require.True(t, ok, "span should be a map")
		require.NotEmpty(t, spanMap["span_id"])
		spanIDStr, ok := spanMap["span_id"].(string)
		require.True(t, ok, "span_id should be a string")
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

// traversePath traverses the path in the data map.
func traversePath(data map[string]interface{}, path string) (map[string]interface{}, error) {
	if path == "." {
		return data, nil
	}

	keys := strings.Split(path, ".")
	current := data

	for _, key := range keys {
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
			continue
		}
		return nil, fmt.Errorf("invalid path: %s", path)
	}

	return current, nil
}
