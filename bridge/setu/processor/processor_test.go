package processor

import (
	"os"
	"testing"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/bridge/setu/broadcaster"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

func TestMain(m *testing.M) {
	// Set the Configuration for All Tests
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	viper.Set(helper.TendermintNodeFlag, "http://localhost:26657")
	viper.Set("log_level", "info")
	helper.InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))
}

func TestBroadcastWhenTxInMempool(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		message        clerkTypes.MsgEventRecord
		expectedStatus bool
		description   string
	}{
		{
			name: "unique message",
			message: clerkTypes.MsgEventRecord{
				From:            hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
				TxHash:          hmTypes.BytesToHeimdallHash([]byte("0xa86a2f8f30d5fab1bb858e572ecf3f24d691276f833f06fc90a745cea20f4fcb")),
				LogIndex:        71,
				BlockNumber:     14475337,
				ContractAddress: hmTypes.BytesToHeimdallAddress([]byte("0x401f6c983ea34274ec46f84d70b31c151321188b")),
				Data:           []byte{},
				ID:             1897091,
				ChainID:        "15001",
			},
			expectedStatus: false,
			description:   "First message should go through as it's unique in mempool",
		},
		{
			name: "duplicate message",
			message: clerkTypes.MsgEventRecord{
				From:            hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
				TxHash:          hmTypes.BytesToHeimdallHash([]byte("0xa86a2f8f30d5fab1bb858e572ecf3f24d691276f833f06fc90a745cea20f4fcb")),
				LogIndex:        71,
				BlockNumber:     14475337,
				ContractAddress: hmTypes.BytesToHeimdallAddress([]byte("0x401f6c983ea34274ec46f84d70b31c151321188b")),
				Data:           []byte{},
				ID:             1897091,
				ChainID:        "15001",
			},
			expectedStatus: true,
			description:   "Duplicate message should fail",
		},
		{
			name: "different txhash",
			message: clerkTypes.MsgEventRecord{
				From:            hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
				TxHash:          hmTypes.BytesToHeimdallHash([]byte("0x8a83aa78a400fe959b44ccf70d926c967af4e451ba630a849b2e1dedc7e30c07")),
				LogIndex:        71,
				BlockNumber:     14475337,
				ContractAddress: hmTypes.BytesToHeimdallAddress([]byte("0x401f6c983ea34274ec46f84d70b31c151321188b")),
				Data:           []byte{},
				ID:             1897091,
				ChainID:        "15001",
			},
			expectedStatus: false,
			description:   "Message with different txhash should go through",
		},
		{
			name: "different logIndex",
			message: clerkTypes.MsgEventRecord{
				From:            hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
				TxHash:          hmTypes.BytesToHeimdallHash([]byte("0xa86a2f8f30d5fab1bb858e572ecf3f24d691276f833f06fc90a745cea20f4fcb")),
				LogIndex:        72,
				BlockNumber:     14475337,
				ContractAddress: hmTypes.BytesToHeimdallAddress([]byte("0x401f6c983ea34274ec46f84d70b31c151321188b")),
				Data:           []byte{},
				ID:             1897091,
				ChainID:        "15001",
			},
			expectedStatus: false,
			description:   "Message with different logIndex should go through",
		},
	}

	cdc := app.MakeCodec()
	txBroadcaster := broadcaster.NewTxBroadcaster(cdc)

	// create a mock clerk processor
	contractCaller, err := helper.NewContractCaller()
	require.NoError(t, err, "Error creating contract caller")

	cp := NewClerkProcessor(&contractCaller.StateSenderABI)
	cp.BaseProcessor = *NewBaseProcessor(cdc, nil, nil, nil, "clerk", cp)

	for _, tc := range testCases {
		tc := tc // capture range variable
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			inMempool, err := cp.checkTxAgainstMempool(tc.message, nil)
			require.NoError(t, err, "Error checking tx against mempool")
			
			assert.Equal(t, tc.expectedStatus, inMempool, tc.description)

			if !inMempool {
				_, err := txBroadcaster.BroadcastToHeimdall(tc.message, nil)
				assert.NoError(t, err, "Error broadcasting tx to heimdall")
			}
		})
	}
}
