package helper

import (
	"math/big"

	"context"

	"errors"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/contracts/stakemanager"
	"github.com/maticnetwork/heimdall/types"

	"strings"

	"math"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type IContractCaller interface {
	GetHeaderInfo(headerID uint64) (root common.Hash, start, end, createdAt uint64, err error)
	GetValidatorInfo(valID types.ValidatorID) (validator types.Validator, err error)
	CurrentChildBlock() (uint64, error)
	CurrentHeaderBlock() (uint64, error)
	GetBalance(address common.Address) (*big.Int, error)
	SendCheckpoint(voteSignBytes []byte, sigs []byte, txData []byte)
	GetMainChainBlock(blockNum *big.Int) (header *ethtypes.Header, err error)
	GetMaticChainBlock(blockNum *big.Int) (header *ethtypes.Header, err error)
	IsTxConfirmed(tx common.Hash) bool
	GetBlockNoFromTxHash(tx common.Hash) (blocknumber big.Int, err error)
	SigUpdateEvent(tx common.Hash) (id uint64, newSigner common.Address, oldSigner common.Address, err error)
}

type ContractCaller struct {
	RootChainInstance    *rootchain.Rootchain
	StakeManagerInstance *stakemanager.Stakemanager
	MainChainClient      *ethclient.Client
	MainChainRPC         *rpc.Client
	MaticClient          *ethclient.Client
	stakeManagerABI      abi.ABI
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}

type rpcTransaction struct {
	tx *ethtypes.Transaction
	txExtraInfo
}

func NewContractCaller() (contractCallerObj ContractCaller, err error) {
	rootChainInstance, err := GetRootChainInstance()
	if err != nil {
		Logger.Error("Error creating rootchain instance", "error", err)
		return contractCallerObj, err
	}
	stakeManagerInstance, err := GetStakeManagerInstance()
	if err != nil {
		Logger.Error("Error creating stakeManagerInstance while getting validator info", "error", err)
		return contractCallerObj, err
	}
	contractCallerObj.MainChainClient = GetMainClient()
	contractCallerObj.MainChainRPC = GetMainChainRPCClient()
	contractCallerObj.StakeManagerInstance = stakeManagerInstance
	contractCallerObj.RootChainInstance = rootChainInstance
	contractCallerObj.MaticClient = GetMaticClient()

	// load stake manager abi
	stakeContract, err := abi.JSON(strings.NewReader(string(stakemanager.StakemanagerABI)))
	if err != nil {
		Logger.Error("Error creating abi for stakemanager", "Error", err)
		return contractCallerObj, err
	}
	contractCallerObj.stakeManagerABI = stakeContract

	return contractCallerObj, nil
}

// GetHeaderInfo get header info from header id
func (c *ContractCaller) GetHeaderInfo(headerID uint64) (root common.Hash, start, end, createdAt uint64, err error) {
	// get header from rootchain
	headerIDInt := big.NewInt(0)
	headerIDInt = headerIDInt.SetUint64(headerID)
	headerBlock, err := c.RootChainInstance.HeaderBlock(nil, headerIDInt)
	if err != nil {
		Logger.Error("Unable to fetch header block from rootchain", "headerBlockIndex", headerID)
		return root, start, end, createdAt, errors.New("Unable to fetch header block")
	}

	return headerBlock.Root, headerBlock.Start.Uint64(), headerBlock.End.Uint64(), headerBlock.CreatedAt.Uint64(), nil
}

// CurrentChildBlock fetch current child block
func (c *ContractCaller) CurrentChildBlock() (uint64, error) {
	currentChildBlock, err := c.RootChainInstance.CurrentChildBlock(nil)
	if err != nil {
		Logger.Error("Could not fetch current child block from rootchain contract", "Error", err)
		return 0, err
	}
	return currentChildBlock.Uint64(), nil
}

// CurrentHeaderBlock fetches current header block
func (c *ContractCaller) CurrentHeaderBlock() (uint64, error) {
	currentHeaderBlock, err := c.RootChainInstance.CurrentHeaderBlock(nil)
	if err != nil {
		Logger.Error("Could not fetch current header block from rootchain contract", "Error", err)
		return 0, err
	}
	return currentHeaderBlock.Uint64(), nil
}

// GetBalance get balance of account (returns big.Int balance wont fit in uint64)
func (c *ContractCaller) GetBalance(address common.Address) (*big.Int, error) {
	balance, err := c.MainChainClient.BalanceAt(context.Background(), address, nil)
	if err != nil {
		Logger.Error("Unable to fetch balance of account from root chain", "Error", err, "Address", address.String())
		return big.NewInt(0), err
	}

	return balance, nil
}

// GetValidatorInfo get validator info
func (c *ContractCaller) GetValidatorInfo(valID types.ValidatorID) (validator types.Validator, err error) {
	amount, startEpoch, endEpoch, signer, err := c.StakeManagerInstance.GetStakerDetails(nil, big.NewInt(int64(valID)))
	if err != nil {
		Logger.Error("Error fetching validator information from stake manager", "Error", err, "ValidatorID", valID)
		return
	}

	decimals := math.Pow10(-18)
	newAmount := decimals * float64(amount.Int64())

	validator = types.Validator{
		ID:         valID,
		Power:      uint64(newAmount),
		StartEpoch: startEpoch.Uint64(),
		EndEpoch:   endEpoch.Uint64(),
		Signer:     types.HeimdallAddress(signer),
	}

	return validator, nil
}

// get main chain block header
func (c *ContractCaller) GetMainChainBlock(blockNum *big.Int) (header *ethtypes.Header, err error) {
	latestBlock, err := GetMainClient().HeaderByNumber(context.Background(), blockNum)
	if err != nil {
		Logger.Error("Unable to connect to main chain", "Error", err)
		return
	}
	return latestBlock, nil
}

// get child chain block header
func (c *ContractCaller) GetMaticChainBlock(blockNum *big.Int) (header *ethtypes.Header, err error) {
	latestBlock, err := GetMainClient().HeaderByNumber(context.Background(), blockNum)
	if err != nil {
		Logger.Error("Unable to connect to matic chain", "Error", err)
		return
	}
	return latestBlock, nil
}

// Get block number of transaction
func (c *ContractCaller) GetBlockNoFromTxHash(tx common.Hash) (blocknumber big.Int, err error) {
	var json *rpcTransaction
	err = c.MainChainRPC.CallContext(context.Background(), &json, "eth_getTransactionByHash", tx)
	if err != nil {
		return
	}
	var blkNum big.Int
	blockNumberPtr, ok := blkNum.SetString(*json.BlockNumber, 0)
	if !ok {
		return blocknumber, errors.New("unable to set string")
	}
	return *blockNumberPtr, nil
}

// Check finality
func (c *ContractCaller) IsTxConfirmed(tx common.Hash) bool {
	txBlk, err := c.GetBlockNoFromTxHash(tx)
	if err != nil {
		Logger.Error("error getting block number from txhash", "Error", err)
		return false
	}
	Logger.Debug("Tx included in block", "block", txBlk.Uint64(), "tx", tx)

	latestBlk, err := c.GetMainChainBlock(nil)
	if err != nil {
		Logger.Error("error getting latest block from main chain", "Error", err)
		return false
	}
	Logger.Debug("Latest block on main chain obtained", "Block", latestBlk.Number.Uint64())

	diff := latestBlk.Number.Uint64() - txBlk.Uint64()
	if diff < GetConfig().ConfirmationBlocks {
		Logger.Info("Not enough confirmations", "ExpectedConf", GetConfig().ConfirmationBlocks, "ActualConf", diff)
		return false
	}

	return true
}

func (c *ContractCaller) SigUpdateEvent(tx common.Hash) (id uint64, newSigner common.Address, oldSigner common.Address, err error) {
	txReceipt, err := c.MainChainClient.TransactionReceipt(context.Background(), tx)
	if err != nil {
		Logger.Error("Unable to get transaction receipt by hash", "Error", err)
		return
	}
	for _, vLog := range txReceipt.Logs {
		oldSigner = common.BytesToAddress(vLog.Topics[2].Bytes())
		newSigner = common.BytesToAddress(vLog.Topics[3].Bytes())
		id = vLog.Topics[1].Big().Uint64()
	}
	return
}
