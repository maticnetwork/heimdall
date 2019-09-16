package helper

import (
	"context"
	"errors"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/contracts/stakemanager"
	"github.com/maticnetwork/heimdall/contracts/validatorset"
	"github.com/maticnetwork/heimdall/types"
)

type IContractCaller interface {
	GetHeaderInfo(headerID uint64) (root common.Hash, start, end, createdAt uint64, proposer types.HeimdallAddress, err error)
	GetValidatorInfo(valID types.ValidatorID) (validator types.Validator, err error)
	GetLastChildBlock() (uint64, error)
	CurrentHeaderBlock() (uint64, error)
	GetBalance(address common.Address) (*big.Int, error)
	SendCheckpoint(voteSignBytes []byte, sigs []byte, txData []byte)
	GetMainChainBlock(blockNum *big.Int) (header *ethtypes.Header, err error)
	GetMaticChainBlock(blockNum *big.Int) (header *ethtypes.Header, err error)
	IsTxConfirmed(tx common.Hash) bool
	GetBlockNoFromTxHash(tx common.Hash) (blocknumber *big.Int, err error)
	SigUpdateEvent(tx common.Hash) (id uint64, newSigner common.Address, oldSigner common.Address, err error)

	// bor related contracts
	CurrentSpanNumber() (Number *big.Int)
}

// ContractCaller contract caller
type ContractCaller struct {
	RootChainInstance    *rootchain.Rootchain
	StakeManagerInstance *stakemanager.Stakemanager
	ValidatorSetInstance *validatorset.Validatorset
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
		Logger.Error("Error creating rootchain instance while creating contract caller obj", "error", err)
		return contractCallerObj, err
	}
	stakeManagerInstance, err := GetStakeManagerInstance()
	if err != nil {
		Logger.Error("Error creating stakeManagerInstance creating contract caller obj", "error", err)
		return contractCallerObj, err
	}
	validatorSetInstance, err := GetValidatorSetInstance()
	if err != nil {
		Logger.Error("Error creating validator set instance while creating contract caller obj", "error", err)
		return contractCallerObj, err
	}
	contractCallerObj.MainChainClient = GetMainClient()
	contractCallerObj.MainChainRPC = GetMainChainRPCClient()
	contractCallerObj.StakeManagerInstance = stakeManagerInstance
	contractCallerObj.RootChainInstance = rootChainInstance
	contractCallerObj.ValidatorSetInstance = validatorSetInstance
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
func (c *ContractCaller) GetHeaderInfo(headerID uint64) (
	root common.Hash,
	start uint64,
	end uint64,
	createdAt uint64,
	proposer types.HeimdallAddress,
	err error,
) {
	// get header from rootchain
	headerBlock, err := c.RootChainInstance.HeaderBlocks(nil, big.NewInt(0).SetUint64(headerID))
	if err != nil {
		Logger.Error("Unable to fetch header block from rootchain", "headerBlockIndex", headerID)
		return root, start, end, createdAt, proposer, errors.New("Unable to fetch header block")
	}

	return headerBlock.Root,
		headerBlock.Start.Uint64(),
		headerBlock.End.Uint64(),
		headerBlock.CreatedAt.Uint64(),
		types.BytesToHeimdallAddress(headerBlock.Proposer.Bytes()),
		nil
}

// GetLastChildBlock fetch current child block
func (c *ContractCaller) GetLastChildBlock() (uint64, error) {
	GetLastChildBlock, err := c.RootChainInstance.GetLastChildBlock(nil)
	if err != nil {
		Logger.Error("Could not fetch current child block from rootchain contract", "Error", err)
		return 0, err
	}
	return GetLastChildBlock.Uint64(), nil
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
	amount, startEpoch, endEpoch, signer, status, err := c.StakeManagerInstance.GetStakerDetails(nil, big.NewInt(int64(valID)))
	if err != nil {
		Logger.Error("Error fetching validator information from stake manager", "error", err, "validatorId", valID, "status", status)
		return
	}

	decimals18 := big.NewInt(10).Exp(big.NewInt(10), big.NewInt(18), nil)
	if amount.Uint64() < decimals18.Uint64() {
		err = errors.New("amount must be more than 1 token")
		return
	}
	newAmount := amount.Div(amount, decimals18)

	// newAmount
	validator = types.Validator{
		ID:         valID,
		Power:      newAmount.Uint64(),
		StartEpoch: startEpoch.Uint64(),
		EndEpoch:   endEpoch.Uint64(),
		Signer:     types.BytesToHeimdallAddress(signer.Bytes()),
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
	latestBlock, err := GetMaticClient().HeaderByNumber(context.Background(), blockNum)
	if err != nil {
		Logger.Error("Unable to connect to matic chain", "Error", err)
		return
	}
	return latestBlock, nil
}

// GetBlockNoFromTxHash gets block number of transaction
func (c *ContractCaller) GetBlockNoFromTxHash(tx common.Hash) (*big.Int, error) {
	var rpcTx rpcTransaction
	if err := c.MainChainRPC.CallContext(context.Background(), &rpcTx, "eth_getTransactionByHash", tx); err != nil {
		return nil, err
	}

	if rpcTx.BlockNumber == nil {
		return nil, errors.New("No tx found")
	}

	blkNum := big.NewInt(0)
	blkNum, ok := blkNum.SetString(*rpcTx.BlockNumber, 0)
	if !ok {
		return nil, errors.New("unable to set string")
	}
	return blkNum, nil
}

// IsTxConfirmed is tx confirmed
func (c *ContractCaller) IsTxConfirmed(tx common.Hash) bool {
	txBlk, err := c.GetBlockNoFromTxHash(tx)
	if err != nil {
		Logger.Error("Error getting block number from txhash", "Error", err)
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

// CurrentSpanNumber get current span
func (c *ContractCaller) CurrentSpanNumber() (Number *big.Int) {
	result, err := c.ValidatorSetInstance.CurrentSpanNumber(nil)
	if err != nil {
		Logger.Error("Unable to get current span number", "Error", err)
		return nil
	}

	return result
}
