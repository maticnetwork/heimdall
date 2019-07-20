package test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/maticnetwork/heimdall/helper"
)

func TestTxDecode(t *testing.T) {
	tx, err := helper.TendermintTxDecode("wWhvHPg6AQHY1wEBlP+zHe/ZNZTQii57ULFjrJulHewY2NcBAZT/sx3v2TWU0Ioue1CxY6ybpR3sGICEXTLzJQ==")
	if err != nil {
		t.Error(err)
	} else {
		expected := "c1686f1cf83a0101d8d7010194ffb31defd93594d08a2e7b50b163ac9ba51dec18d8d7010194ffb31defd93594d08a2e7b50b163ac9ba51dec1880845d32f325"
		require.Equal(t, expected, hex.EncodeToString(tx), "Tx encoding should match")
	}
}
