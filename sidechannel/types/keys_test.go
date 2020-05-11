package types_test

import (
	"encoding/binary"
	"testing"

	"github.com/maticnetwork/heimdall/sidechannel/types"
	"github.com/stretchr/testify/require"
)

func TestTxStoreKey(t *testing.T) {
	hash := []byte("test-bytes")
	txStoreKey := types.TxStoreKey(120, hash)

	require.Less(t, 9, len(txStoreKey), "TxStoreKey should be enough length")
	require.Equal(t, types.TxsKeyPrefix, txStoreKey[:1], "TxStoreKey should have valid prefix")

	data := binary.BigEndian.Uint64(txStoreKey[1:9])
	require.Equal(t, uint64(120), data, "TxStoreKey should have valid height in key")

	require.Equal(t, hash, txStoreKey[9:], "TxStoreKey should have valid hash in key")
}

func TestTxsStoreKey(t *testing.T) {
	txsStoreKey := types.TxsStoreKey(120)

	require.Equal(t, 9, len(txsStoreKey), "TxsStoreKey should be enough length")
	require.Equal(t, types.TxsKeyPrefix, txsStoreKey[:1], "TxsStoreKey should have valid prefix")

	data := binary.BigEndian.Uint64(txsStoreKey[1:9])
	require.Equal(t, uint64(120), data, "TxsStoreKey should have valid height in key")
}

func TestValidatorsKey(t *testing.T) {
	validatorsKey := types.ValidatorsKey(120)

	require.Equal(t, 9, len(validatorsKey), "ValidatorsKey should be enough length")
	require.Equal(t, types.ValidatorsKeyPrefix, validatorsKey[:1], "ValidatorsKey should have valid prefix")

	data := binary.BigEndian.Uint64(validatorsKey[1:9])
	require.Equal(t, uint64(120), data, "ValidatorsKey should have valid height in key")
}
