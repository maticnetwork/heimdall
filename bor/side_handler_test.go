package bor_test

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethereum "github.com/maticnetwork/bor"
	borCommon "github.com/maticnetwork/bor/common"
	ethTypes "github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/bor"
	borTypes "github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

type sideChHandlerSuite struct {
	suite.Suite
	app        *app.HeimdallApp
	ctx        sdk.Context
	mockCaller mocks.IContractCaller
}

func TestBorSideChHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(sideChHandlerSuite))
}

func (suite *sideChHandlerSuite) SetupTest() {
	isCheckTx := false
	suite.app = app.Setup(isCheckTx)
	suite.ctx = suite.app.BaseApp.NewContext(isCheckTx, abci.Header{})
	suite.mockCaller = mocks.IContractCaller{}
}

func (suite *sideChHandlerSuite) TestSideHandleMsgSpan() {
	var bi *big.Int
	bi = nil
	ethBlockData := `{
		"difficulty":"997888",
		"extraData":"0xd883010503846765746887676f312e372e318664617277696e",
		"gasLimit":16760833,
		"gasUsed":0,
		"hash":"0x41800b5c3f1717687d85fc9018faac0a6e90b39deaa0b99e7fe4fe796ddeb26a",
		"logsBloom":"0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
		"miner":"0xd1aeb42885a43b72b518182ef893125814811048",
		"mixHash":"0x0f98b15f1a4901a7e9204f3c500a7bd527b3fb2c3340e12176a44b83e414a69e",
		"nonce":"0x0ece08ea8c49dfd9",
		"number":1,
		"parentHash":"0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d",
		"receiptsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"sha3Uncles":"0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
		"size":536,
		"stateRoot":"0xc7b01007a10da045eacb90385887dd0c38fcb5db7393006bdde24b93873c334b",
		"timestamp":1479642530,
		"totalDifficulty":"2046464",
		"transactions":[],
		"transactionsRoot":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
		"uncles":[]
	}`
	ethBlockHash := `0x88e96d4537bea4d9c05d12549907b32561d3bf31f45aae734cdc119f13406cb6`

	var ethHeader ethTypes.Header
	suite.Nil(json.Unmarshal(borCommon.Hex2Bytes(ethBlockData), &ethHeader))

	type callerMethod struct {
		name string
		args []interface{}
		ret  []interface{}
	}
	tc := []struct {
		out       abci.ResponseDeliverSideTx
		msg       string
		codespace string
		code      common.CodeType
		cm        []callerMethod
		seed      borCommon.Hash
	}{
		{
			codespace: "1",
			code:      common.CodeInvalidMsg,
			msg:       "error mainchain error",
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethTypes.Header{}, ethereum.NotFound},
				},
			},
		},
		{
			msg:       "error msg seed bytes failure",
			codespace: "1",
			code:      common.CodeInvalidMsg,
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethTypes.Header{}, nil},
				},
			},
		},
		{
			msg:       "error failed to GetMaticChainBlock",
			codespace: "1",
			code:      common.CodeInvalidMsg,
			seed:      borCommon.HexToHash(ethBlockHash),
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethTypes.Header{}, nil},
				},
				{
					name: "GetMaticChainBlock",
					args: []interface{}{bi},
					ret:  []interface{}{ethHeader, ethereum.NotFound},
				},
			},
		},
	}

	for i, c := range tc {
		suite.SetupTest()
		c.msg = fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		if c.cm != nil {
			for _, m := range c.cm {
				suite.mockCaller.On(m.name, m.args...).Return(m.ret...)
			}
		}
		fmt.Println(c.seed.Hex())
		out := bor.SideHandleMsgSpan(suite.ctx, suite.app.BorKeeper, borTypes.MsgProposeSpan{Seed: c.seed}, &suite.mockCaller)
		// construct output
		c.out = abci.ResponseDeliverSideTx{Code: uint32(c.code), Codespace: c.codespace, Result: abci.SideTxResultType_Skip}
		suite.Equal(c.out, out, c.msg)
	}
}
