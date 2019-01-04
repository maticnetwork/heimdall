package helper

import (
	"math/big"

	"context"
	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/contracts/stakemanager"
	"github.com/ethereum/go-ethereum/ethclient"
)

type ContractCaller interface {
	GetHeaderInfo(headerID uint64) (root common.Hash, start uint64, end uint64, err error)
	GetValidatorInfo(addr common.Address) (validator types.Validator, err error)
	CurrentChildBlock() (uint64, error)
	GetBalance(address common.Address) (*big.Int, error)
}

type ContractCallerObj struct {
	rootChainInstance rootchain.Rootchain
	stakeManagerInstance stakemanager.Stakemanager
	mainChainClient *ethclient.Client
}

func NewContractCallerObj()  (contractCaller ContractCaller,err error) {
	var contractCallerObj ContractCallerObj
	rootChainInstance, err := GetRootChainInstance()
	if err != nil {
		Logger.Error("Error creating rootchain instance ", "error", err)
		return contractCaller,err
	}
	stakeManagerInstance, err := GetStakeManagerInstance()
	if err != nil {
		Logger.Error("Error creating stakeManagerInstance while getting validator info", "error", err)
		return contractCaller,err
	}
	contractCallerObj.mainChainClient = GetMainClient()
	contractCallerObj.stakeManagerInstance = *stakeManagerInstance
	contractCallerObj.rootChainInstance = *rootChainInstance

	return contractCaller,nil
}

// GetHeaderInfo get header info from header id
func (c *ContractCallerObj) GetHeaderInfo(headerID uint64) (root common.Hash, start uint64, end uint64, err error) {
	// get header from rootchain
	headerIDInt := big.NewInt(0)
	headerIDInt.SetUint64(headerID)
	headerBlock, err := c.rootChainInstance.HeaderBlock(nil, headerIDInt)
	if err != nil {
		Logger.Error("Unable to fetch header block from rootchain", "headerBlockIndex", headerID)
	}

	return headerBlock.Root, headerBlock.Start.Uint64(), headerBlock.End.Uint64(), nil
}

// GetValidatorInfo get validator info
func (c *ContractCallerObj) GetValidatorInfo(addr common.Address) (validator types.Validator, err error) {
	amount, startEpoch, endEpoch, signer, err := c.stakeManagerInstance.GetStakerDetails(nil, addr)
	if err != nil {
		Logger.Error("Error fetching validator information from stakemanager", "Error", err, "ValidatorAddress", addr)
		return
	}
	validator = types.Validator{
		Address:    addr,
		Power:      amount.Uint64(),
		StartEpoch: startEpoch.Uint64(),
		EndEpoch:   endEpoch.Uint64(),
		Signer:     signer,
	}

	return validator, nil
}

// CurrentChildBlock fetch current child block
func (c *ContractCallerObj) CurrentChildBlock() (uint64, error) {
	currentChildBlock, err := c.rootChainInstance.CurrentChildBlock(nil)
	if err != nil {
		Logger.Error("Could not fetch current child block from rootchain contract", "Error", err)
		return 0, err
	}

	return currentChildBlock.Uint64(), nil
}

// get balance of account (returns big.Int balance wont fit in uint64)
func (c *ContractCallerObj)GetBalance(address common.Address) (*big.Int, error) {
	balance, err := c.mainChainClient.BalanceAt(context.Background(), address, nil)
	if err != nil {
		Logger.Error("Unable to fetch balance of account from root chain", "Error", err, "Address", address.String())
		return big.NewInt(0), err
	}

	return balance, nil
}
