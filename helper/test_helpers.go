package helper

import (
	abci "github.com/tendermint/tendermint/abci/types"
)

type TestOpts struct {
	app     abci.Application
	chainId string
}

func NewTestOpts(app abci.Application, chainId string) *TestOpts {
	return &TestOpts{
		app:     app,
		chainId: chainId,
	}
}

func (t *TestOpts) SetApplication(app abci.Application) {
	t.app = app
}

func (t *TestOpts) GetApplication() abci.Application {
	return t.app
}

func (t *TestOpts) SetChainId(chainId string) {
	t.chainId = chainId
}

func (t *TestOpts) GetChainId() string {
	return t.chainId
}
