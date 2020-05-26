package bor_test

import (
	"fmt"
	"math/big"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
	}{
		//	{
		//		codespace: "1",
		//		code:      common.CodeInvalidMsg,
		//		msg:       "error mainchain error",
		//	},
		{
			codespace: "1",
			// code:      common.CodeInvalidMsg,
			cm: []callerMethod{
				{
					name: "GetMainChainBlock",
					args: []interface{}{big.NewInt(1)},
					ret:  []interface{}{&ethTypes.Header{}, nil},
				},
			},
		},
	}

	for i, c := range tc {
		// suite.SetupTest()
		c.msg = fmt.Sprintf("i: %v, msg: %v", i, c.msg)
		if c.cm != nil {
			for _, m := range c.cm {
				suite.mockCaller.On(m.name, big.NewInt(1)).Return(&ethTypes.Header{}, nil)
			}
		}

		out := bor.SideHandleMsgSpan(suite.ctx, suite.app.BorKeeper, borTypes.MsgProposeSpan{}, &suite.mockCaller)
		// construct output
		c.out = abci.ResponseDeliverSideTx{Code: uint32(c.code), Codespace: c.codespace, Result: abci.SideTxResultType_Skip}
		suite.Equal(c.out, out, c.msg)
	}
}
