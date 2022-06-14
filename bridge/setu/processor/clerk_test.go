package processor

import (
	"bytes"
	"encoding/json"
	"github.com/golang/mock/gomock"
	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	authTypesMocks "github.com/maticnetwork/heimdall/auth/types/mocks"
	"github.com/maticnetwork/heimdall/bridge/setu/broadcaster"
	"github.com/maticnetwork/heimdall/helper"
	helperMocks "github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

const (
	dummyTenderMintNode    = "http://dummy-localhost:26657"
	dummyHeimdallServerUrl = "https://dummy-heimdall-api.polygon.technology"

	chainManagerParamsUrl      = dummyHeimdallServerUrl + "/chainmanager/params"
	chainManagerParamsResponse = `
	{
	"height": "0",
	"result": {
		"mainchain_tx_confirmations": 6,
		"maticchain_tx_confirmations": 10,
		"chain_params": {
			"bor_chain_id": "80001",
			"matic_token_address": "0x499d11E0b6eAC7c0593d8Fb292DCBbF815Fb29Ae",
			"staking_manager_address": "0x4864d89DCE4e24b2eDF64735E014a7E4154bfA7A",
			"slash_manager_address": "0x93D8f8A1A88498b258ceb69dD82311962374269C",
			"root_chain_address": "0x2890bA17EfE978480615e330ecB65333b880928e",
			"staking_info_address": "0x318EeD65F064904Bc6E0e3842940c5972BC8E38f",
			"state_sender_address": "0xEAa852323826C71cd7920C3b4c007184234c3945",
			"state_receiver_address": "0x0000000000000000000000000000000000001001",
			"validator_set_address": "0x0000000000000000000000000000000000001000"
			}
		}
	}`

	getAccountUrl      = dummyHeimdallServerUrl + "/auth/accounts/9FB29AAC15B9A4B7F17C3385939B007540F4D791"
	getAccountResponse = `
	{
	"height": "0",
	"result": {
		"type": "auth/Account",
		"value": {
			"address": "0x5973918275c01f50555d44e92c9d9b353cadad54",
			"coins": [{
				"denom": "matic",
				"amount": "10000000000000000000"
			}],
			"public_key": null,
			"account_number": "0",
			"sequence_number": "0",
			"name": "",
			"permissions": []
			}
		}
	}`

	isOldTxUrl = dummyHeimdallServerUrl + "/clerk/isoldtx?logindex=0&txhash=0x6d428739815d7c84cf89db055158861b089e0fd649676a0243a2a2d204c1d854"
	// TODO test false case
	isOldTxResponse = `
	{
		"height": "0",
		"result": true
	}`
)

func BenchmarkSendStateSyncedToHeimdall(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		b.Logf("Executing iteration '%d' out of '%d'", i, b.N)
		// given
		prepareMockData(b)
		cp, err := prepareClerkProcessor()
		if err != nil {
			b.Fatal("Error initializing test clerk processor")
		}
		dlb, err := prepareDummyLogBytes()
		if err != nil {
			b.Fatal("Error creating test data")
		}
		// when
		b.StartTimer()
		err = cp.sendStateSyncedToHeimdall("StateSynced", dlb.String())
		// then
		if err != nil {
			b.Fatal(err)
		}
		b.Log("StateSynced sent to heimdall successfully")
	}
}

func prepareMockData(b *testing.B) {
	mockCtrl := gomock.NewController(b)
	defer mockCtrl.Finish()
	mockHttpClient := helperMocks.NewMockHTTPClient(mockCtrl)
	mockNodeQuerier := authTypesMocks.NewMockNodeQuerier(mockCtrl)

	mockHttpClient.EXPECT().Get(chainManagerParamsUrl).Return(prepareResponse(chainManagerParamsResponse), nil).AnyTimes()
	mockHttpClient.EXPECT().Get(getAccountUrl).Return(prepareResponse(getAccountResponse), nil).AnyTimes()
	mockHttpClient.EXPECT().Get(isOldTxUrl).Return(prepareResponse(isOldTxResponse), nil).AnyTimes()
	helper.Client = mockHttpClient

	mockNodeQuerier.EXPECT().QueryWithData(gomock.Any(), gomock.Any()).Return(nil, int64(0), nil).AnyTimes()
	authTypes.NQuerier = mockNodeQuerier
}

func prepareClerkProcessor() (*ClerkProcessor, error) {
	cdc := app.MakeCodec()

	viper.Set(helper.NodeFlag, dummyTenderMintNode)
	viper.Set("log_level", "debug")

	helper.InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))
	config := helper.GetConfig()
	config.HeimdallServerURL = dummyHeimdallServerUrl
	helper.SetTestConfig(config)

	txBroadcaster := broadcaster.NewTxBroadcaster(cdc)
	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		return nil, err
	}
	cp := NewClerkProcessor(&contractCaller.StateSenderABI)
	cp.BaseProcessor = *NewBaseProcessor(cdc, nil, nil, txBroadcaster, "clerk", cp)

	return cp, nil
}

func prepareDummyLogBytes() (*bytes.Buffer, error) {
	topics := append([]common.Hash{},
		common.HexToHash("0x103fed9db65eac19c4d870f49ab7520fe03b99f1838e5996caf47e9e43308392"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000001ef6e0"),
		common.HexToHash("0x000000000000000000000000a6fa4fb5f76172d178d61b04b0ecd319c5d1c0aa"))

	log := types.Log{
		Address:     common.HexToAddress("0x28e4f3a7f651294b9564800b2d01f35189a5bfbe"),
		Topics:      topics,
		Data:        common.FromHex("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000010087a7811f4bfedea3d341ad165680ae306b01aaeacc205d227629cf157dd9f821000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000004aa11c9581573571f963bda7a41b28d90c36027c000000000000000000000000eeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee0000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000b1a2bc2ec50000"),
		BlockNumber: 14702845,
		TxHash:      common.HexToHash("0x6d428739815d7c84cf89db055158861b089e0fd649676a0243a2a2d204c1d854"),
		TxIndex:     0,
		BlockHash:   common.HexToHash("0xe8370360b861be304ef4144e33a3803cf6d4e31524832444ada797e16f859438"),
		Index:       0,
		Removed:     false,
	}
	reqBodyBytes := new(bytes.Buffer)
	err := json.NewEncoder(reqBodyBytes).Encode(log)
	if err != nil {
		return nil, err
	}

	return reqBodyBytes, nil
}

func prepareResponse(body string) *http.Response {
	return &http.Response{
		Status:           "200 OK",
		StatusCode:       200,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             ioutil.NopCloser(bytes.NewReader([]byte(body))),
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}
}
