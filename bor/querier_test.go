package bor_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

type querierHandlerSuite struct {
	suite.Suite
	app        *app.HeimdallApp
	ctx        sdk.Context
	mockCaller mocks.IContractCaller
}

func TestBorQuerierHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(sideChHandlerSuite))
}

func (suite *querierHandlerSuite) SetupTest() {
	isCheckTx := false
	suite.app = app.Setup(isCheckTx)
	suite.ctx = suite.app.BaseApp.NewContext(isCheckTx, abci.Header{})
	suite.mockCaller = mocks.IContractCaller{}
}

func (suite *querierHandlerSuite) TestNewQueirer() {
	tc := []struct {
		msg string
	}{}
	for i, c := range tc {
		c.msg = fmt.Sprintf("i: %v, msg: %v", i, c.msg)
	}
}
