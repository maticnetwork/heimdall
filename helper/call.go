package helper

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	ethTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	lru "github.com/hashicorp/golang-lru"

	"github.com/maticnetwork/heimdall/contracts/erc20"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/contracts/slashmanager"
	"github.com/maticnetwork/heimdall/contracts/stakemanager"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/contracts/statereceiver"
	"github.com/maticnetwork/heimdall/contracts/statesender"
	"github.com/maticnetwork/heimdall/contracts/validatorset"
	"github.com/maticnetwork/heimdall/types"
)

// smart contracts' events names
const (
	newHeaderBlockEvent = "NewHeaderBlock"
	topUpFeeEvent       = "TopUpFee"
	stakedEvent         = "Staked"
	stakeUpdateEvent    = "StakeUpdate"
	UnstakeInitEvent    = "UnstakeInit"
	signerChangeEvent   = "SignerChange"
	stateSyncedEvent    = "StateSynced"
	slashedEvent        = "Slashed"
	unJailedEvent       = "UnJailed"
)

// ContractsABIsMap is a cached map holding the ABIs of the contracts
var ContractsABIsMap = make(map[string]*abi.ABI)

// IContractCaller represents contract caller
type IContractCaller interface {
	GetHeaderInfo(headerID uint64, rootChainInstance *rootchain.Rootchain, childBlockInterval uint64) (root common.Hash, start, end, createdAt uint64, proposer types.HeimdallAddress, err error)
	GetRootHash(start uint64, end uint64, checkpointLength uint64) ([]byte, error)
	GetValidatorInfo(valID types.ValidatorID, stakingInfoInstance *stakinginfo.Stakinginfo) (validator types.Validator, err error)
	GetLastChildBlock(rootChainInstance *rootchain.Rootchain) (uint64, error)
	CurrentHeaderBlock(rootChainInstance *rootchain.Rootchain, childBlockInterval uint64) (uint64, error)
	GetBalance(address common.Address) (*big.Int, error)
	SendCheckpoint(signedData []byte, sigs [][3]*big.Int, rootChainAddress common.Address, rootChainInstance *rootchain.Rootchain) (err error)
	SendTick(signedData []byte, sigs []byte, slashManagerAddress common.Address, slashManagerInstance *slashmanager.Slashmanager) (err error)
	GetCheckpointSign(txHash common.Hash) ([]byte, []byte, []byte, error)
	GetMainChainBlock(*big.Int) (*ethTypes.Header, error)
	GetMaticChainBlock(*big.Int) (*ethTypes.Header, error)
	IsTxConfirmed(common.Hash, uint64) bool
	GetConfirmedTxReceipt(common.Hash, uint64) (*ethTypes.Receipt, error)
	GetBlockNumberFromTxHash(common.Hash) (*big.Int, error)

	// decode header event
	DecodeNewHeaderBlockEvent(common.Address, *ethTypes.Receipt, uint64) (*rootchain.RootchainNewHeaderBlock, error)
	// decode validator events
	DecodeValidatorTopupFeesEvent(common.Address, *ethTypes.Receipt, uint64) (*stakinginfo.StakinginfoTopUpFee, error)
	DecodeValidatorJoinEvent(common.Address, *ethTypes.Receipt, uint64) (*stakinginfo.StakinginfoStaked, error)
	DecodeValidatorStakeUpdateEvent(common.Address, *ethTypes.Receipt, uint64) (*stakinginfo.StakinginfoStakeUpdate, error)
	DecodeValidatorExitEvent(common.Address, *ethTypes.Receipt, uint64) (*stakinginfo.StakinginfoUnstakeInit, error)
	DecodeSignerUpdateEvent(common.Address, *ethTypes.Receipt, uint64) (*stakinginfo.StakinginfoSignerChange, error)
	// decode state events
	DecodeStateSyncedEvent(common.Address, *ethTypes.Receipt, uint64) (*statesender.StatesenderStateSynced, error)

	// decode slashing events
	DecodeSlashedEvent(common.Address, *ethTypes.Receipt, uint64) (*stakinginfo.StakinginfoSlashed, error)
	DecodeUnJailedEvent(common.Address, *ethTypes.Receipt, uint64) (*stakinginfo.StakinginfoUnJailed, error)

	GetMainTxReceipt(common.Hash) (*ethTypes.Receipt, error)
	GetMaticTxReceipt(common.Hash) (*ethTypes.Receipt, error)
	ApproveTokens(*big.Int, common.Address, common.Address, *erc20.Erc20) error
	StakeFor(common.Address, *big.Int, *big.Int, bool, common.Address, *stakemanager.Stakemanager) error
	CurrentAccountStateRoot(stakingInfoInstance *stakinginfo.Stakinginfo) ([32]byte, error)

	// bor related contracts
	CurrentSpanNumber(validatorSet *validatorset.Validatorset) (Number *big.Int)
	GetSpanDetails(id *big.Int, validatorSet *validatorset.Validatorset) (*big.Int, *big.Int, *big.Int, error)
	CurrentStateCounter(stateSenderInstance *statesender.Statesender) (Number *big.Int)
	CheckIfBlocksExist(end uint64) bool

	GetRootChainInstance(rootChainAddress common.Address) (*rootchain.Rootchain, error)
	GetStakingInfoInstance(stakingInfoAddress common.Address) (*stakinginfo.Stakinginfo, error)
	GetValidatorSetInstance(validatorSetAddress common.Address) (*validatorset.Validatorset, error)
	GetStakeManagerInstance(stakingManagerAddress common.Address) (*stakemanager.Stakemanager, error)
	GetSlashManagerInstance(slashManagerAddress common.Address) (*slashmanager.Slashmanager, error)
	GetStateSenderInstance(stateSenderAddress common.Address) (*statesender.Statesender, error)
	GetStateReceiverInstance(stateReceiverAddress common.Address) (*statereceiver.Statereceiver, error)
	GetMaticTokenInstance(maticTokenAddress common.Address) (*erc20.Erc20, error)
}

// ContractCaller contract caller
type ContractCaller struct {
	MainChainClient   *ethclient.Client
	MainChainRPC      *rpc.Client
	MainChainTimeout  time.Duration
	MaticChainClient  *ethclient.Client
	MaticChainRPC     *rpc.Client
	MaticChainTimeout time.Duration

	RootChainABI     abi.ABI
	StakingInfoABI   abi.ABI
	ValidatorSetABI  abi.ABI
	StateReceiverABI abi.ABI
	StateSenderABI   abi.ABI
	StakeManagerABI  abi.ABI
	SlashManagerABI  abi.ABI
	MaticTokenABI    abi.ABI

	ReceiptCache *lru.Cache

	ContractInstanceCache map[common.Address]interface{}
}

type txExtraInfo struct {
	BlockNumber *string         `json:"blockNumber,omitempty"`
	BlockHash   *common.Hash    `json:"blockHash,omitempty"`
	From        *common.Address `json:"from,omitempty"`
}

type rpcTransaction struct {
	txExtraInfo
}

// NewContractCaller contract caller
func NewContractCaller() (contractCallerObj ContractCaller, err error) {
	config := GetConfig()
	contractCallerObj.MainChainClient = GetMainClient()
	contractCallerObj.MainChainTimeout = config.EthRPCTimeout
	contractCallerObj.MaticChainClient = GetMaticClient()
	contractCallerObj.MaticChainTimeout = config.BorRPCTimeout
	contractCallerObj.MainChainRPC = GetMainChainRPCClient()
	contractCallerObj.MaticChainRPC = GetMaticRPCClient()
	contractCallerObj.ReceiptCache, err = lru.New(1000)

	if err != nil {
		return contractCallerObj, err
	}

	// listeners and processors instance cache (address->ABI)
	contractCallerObj.ContractInstanceCache = make(map[common.Address]interface{})
	// package global cache (string->ABI)
	if err = populateABIs(&contractCallerObj); err != nil {
		return contractCallerObj, err
	}

	return contractCallerObj, nil
}

// GetRootChainInstance returns RootChain contract instance for selected base chain
func (c *ContractCaller) GetRootChainInstance(rootChainAddress common.Address) (*rootchain.Rootchain, error) {
	contractInstance, ok := c.ContractInstanceCache[rootChainAddress]
	if !ok {
		ci, err := rootchain.NewRootchain(rootChainAddress, mainChainClient)
		c.ContractInstanceCache[rootChainAddress] = ci

		return ci, err
	}

	return contractInstance.(*rootchain.Rootchain), nil
}

// GetStakingInfoInstance returns stakingInfo contract instance for selected base chain
func (c *ContractCaller) GetStakingInfoInstance(stakingInfoAddress common.Address) (*stakinginfo.Stakinginfo, error) {
	contractInstance, ok := c.ContractInstanceCache[stakingInfoAddress]
	if !ok {
		ci, err := stakinginfo.NewStakinginfo(stakingInfoAddress, mainChainClient)
		c.ContractInstanceCache[stakingInfoAddress] = ci

		return ci, err
	}

	return contractInstance.(*stakinginfo.Stakinginfo), nil
}

// GetValidatorSetInstance returns stakingInfo contract instance for selected base chain
func (c *ContractCaller) GetValidatorSetInstance(validatorSetAddress common.Address) (*validatorset.Validatorset, error) {
	contractInstance, ok := c.ContractInstanceCache[validatorSetAddress]
	if !ok {
		ci, err := validatorset.NewValidatorset(validatorSetAddress, mainChainClient)
		c.ContractInstanceCache[validatorSetAddress] = ci

		return ci, err
	}

	return contractInstance.(*validatorset.Validatorset), nil
}

// GetStakeManagerInstance returns stakingInfo contract instance for selected base chain
func (c *ContractCaller) GetStakeManagerInstance(stakingManagerAddress common.Address) (*stakemanager.Stakemanager, error) {
	contractInstance, ok := c.ContractInstanceCache[stakingManagerAddress]
	if !ok {
		ci, err := stakemanager.NewStakemanager(stakingManagerAddress, mainChainClient)
		c.ContractInstanceCache[stakingManagerAddress] = ci

		return ci, err
	}

	return contractInstance.(*stakemanager.Stakemanager), nil
}

// GetSlashManagerInstance returns slashManager contract instance for selected base chain
func (c *ContractCaller) GetSlashManagerInstance(slashManagerAddress common.Address) (*slashmanager.Slashmanager, error) {
	contractInstance, ok := c.ContractInstanceCache[slashManagerAddress]
	if !ok {
		ci, err := slashmanager.NewSlashmanager(slashManagerAddress, mainChainClient)
		c.ContractInstanceCache[slashManagerAddress] = ci

		return ci, err
	}

	return contractInstance.(*slashmanager.Slashmanager), nil
}

// GetStateSenderInstance returns stakingInfo contract instance for selected base chain
func (c *ContractCaller) GetStateSenderInstance(stateSenderAddress common.Address) (*statesender.Statesender, error) {
	contractInstance, ok := c.ContractInstanceCache[stateSenderAddress]
	if !ok {
		ci, err := statesender.NewStatesender(stateSenderAddress, mainChainClient)
		c.ContractInstanceCache[stateSenderAddress] = ci

		return ci, err
	}

	return contractInstance.(*statesender.Statesender), nil
}

// GetStateReceiverInstance returns stakingInfo contract instance for selected base chain
func (c *ContractCaller) GetStateReceiverInstance(stateReceiverAddress common.Address) (*statereceiver.Statereceiver, error) {
	contractInstance, ok := c.ContractInstanceCache[stateReceiverAddress]
	if !ok {
		ci, err := statereceiver.NewStatereceiver(stateReceiverAddress, maticClient)
		c.ContractInstanceCache[stateReceiverAddress] = ci

		return ci, err
	}

	return contractInstance.(*statereceiver.Statereceiver), nil
}

// GetMaticTokenInstance returns stakingInfo contract instance for selected base chain
func (c *ContractCaller) GetMaticTokenInstance(maticTokenAddress common.Address) (*erc20.Erc20, error) {
	contractInstance, ok := c.ContractInstanceCache[maticTokenAddress]
	if !ok {
		ci, err := erc20.NewErc20(maticTokenAddress, mainChainClient)
		c.ContractInstanceCache[maticTokenAddress] = ci

		return ci, err
	}

	return contractInstance.(*erc20.Erc20), nil
}

// GetHeaderInfo get header info from checkpoint number
func (c *ContractCaller) GetHeaderInfo(number uint64, rootChainInstance *rootchain.Rootchain, childBlockInterval uint64) (
	root common.Hash,
	start uint64,
	end uint64,
	createdAt uint64,
	proposer types.HeimdallAddress,
	err error,
) {
	// get header from rootChain
	checkpointBigInt := big.NewInt(0).Mul(big.NewInt(0).SetUint64(number), big.NewInt(0).SetUint64(childBlockInterval))

	headerBlock, err := rootChainInstance.HeaderBlocks(nil, checkpointBigInt)
	if err != nil {
		return root, start, end, createdAt, proposer, errors.New("unable to fetch checkpoint block")
	}

	return headerBlock.Root,
		headerBlock.Start.Uint64(),
		headerBlock.End.Uint64(),
		headerBlock.CreatedAt.Uint64(),
		types.BytesToHeimdallAddress(headerBlock.Proposer.Bytes()),
		nil
}

// GetRootHash get root hash from bor chain
func (c *ContractCaller) GetRootHash(start uint64, end uint64, checkpointLength uint64) ([]byte, error) {
	noOfBlock := end - start + 1

	if start > end {
		return nil, errors.New("start is greater than end")
	}

	if noOfBlock > checkpointLength {
		return nil, errors.New("number of headers requested exceeds")
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.MaticChainTimeout)
	defer cancel()

	rootHash, err := c.MaticChainClient.GetRootHash(ctx, start, end)
	if err != nil {
		Logger.Error("Could not fetch rootHash from matic chain", "error", err)
		return nil, err
	}

	return common.FromHex(rootHash), nil
}

// GetLastChildBlock fetch current child block
func (c *ContractCaller) GetLastChildBlock(rootChainInstance *rootchain.Rootchain) (uint64, error) {
	GetLastChildBlock, err := rootChainInstance.GetLastChildBlock(nil)
	if err != nil {
		Logger.Error("Could not fetch current child block from rootChain contract", "error", err)
		return 0, err
	}

	return GetLastChildBlock.Uint64(), nil
}

// CurrentHeaderBlock fetches current header block
func (c *ContractCaller) CurrentHeaderBlock(rootChainInstance *rootchain.Rootchain, childBlockInterval uint64) (uint64, error) {
	currentHeaderBlock, err := rootChainInstance.CurrentHeaderBlock(nil)
	if err != nil {
		Logger.Error("Could not fetch current header block from rootChain contract", "error", err)
		return 0, err
	}

	return currentHeaderBlock.Uint64() / childBlockInterval, nil
}

// GetBalance get balance of account (returns big.Int balance won't fit in uint64)
func (c *ContractCaller) GetBalance(address common.Address) (*big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MainChainTimeout)
	defer cancel()

	balance, err := c.MainChainClient.BalanceAt(ctx, address, nil)
	if err != nil {
		Logger.Error("Unable to fetch balance of account from root chain", "Address", address.String(), "error", err)
		return big.NewInt(0), err
	}

	return balance, nil
}

// GetValidatorInfo get validator info
func (c *ContractCaller) GetValidatorInfo(valID types.ValidatorID, stakingInfoInstance *stakinginfo.Stakinginfo) (validator types.Validator, err error) {
	// amount, startEpoch, endEpoch, signer, status, err := c.StakingInfoInstance.GetStakerDetails(nil, big.NewInt(int64(valID)))
	stakerDetails, err := stakingInfoInstance.GetStakerDetails(nil, big.NewInt(int64(valID)))
	if err != nil {
		Logger.Error("Error fetching validator information from stake manager", "validatorId", valID, "status", stakerDetails.Status, "error", err)
		return
	}

	newAmount, err := GetPowerFromAmount(stakerDetails.Amount)
	if err != nil {
		return
	}

	// newAmount
	validator = types.Validator{
		ID:          valID,
		VotingPower: newAmount.Int64(),
		StartEpoch:  stakerDetails.ActivationEpoch.Uint64(),
		EndEpoch:    stakerDetails.DeactivationEpoch.Uint64(),
		Signer:      types.BytesToHeimdallAddress(stakerDetails.Signer.Bytes()),
	}

	return validator, nil
}

// GetMainChainBlock returns main chain block header
func (c *ContractCaller) GetMainChainBlock(blockNum *big.Int) (header *ethTypes.Header, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MainChainTimeout)
	defer cancel()

	latestBlock, err := c.MainChainClient.HeaderByNumber(ctx, blockNum)
	if err != nil {
		Logger.Error("Unable to connect to main chain", "error", err)
		return
	}

	return latestBlock, nil
}

// GetMainChainFinalizedBlock returns finalized main chain block header (post-merge)
func (c *ContractCaller) GetMainChainFinalizedBlock() (header *ethTypes.Header, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MainChainTimeout)
	defer cancel()

	latestFinalizedBlock, err := c.MainChainClient.HeaderByNumber(ctx, big.NewInt(int64(rpc.FinalizedBlockNumber)))
	if err != nil {
		Logger.Error("Unable to connect to main chain", "error", err)
		return
	}

	return latestFinalizedBlock, nil
}

// GetMainChainBlockTime returns main chain block time
func (c *ContractCaller) GetMainChainBlockTime(ctx context.Context, blockNum uint64) (time.Time, error) {
	ctx, cancel := context.WithTimeout(ctx, c.MainChainTimeout)
	defer cancel()

	latestBlock, err := c.MainChainClient.BlockByNumber(ctx, big.NewInt(0).SetUint64(blockNum))
	if err != nil {
		Logger.Error("Unable to connect to main chain", "error", err)
		return time.Time{}, err
	}

	return time.Unix(int64(latestBlock.Time()), 0), nil
}

// GetMaticChainBlock returns child chain block header
func (c *ContractCaller) GetMaticChainBlock(blockNum *big.Int) (header *ethTypes.Header, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MaticChainTimeout)
	defer cancel()

	latestBlock, err := c.MaticChainClient.HeaderByNumber(ctx, blockNum)
	if err != nil {
		Logger.Error("Unable to connect to matic chain", "error", err)
		return
	}

	return latestBlock, nil
}

// GetBlockNumberFromTxHash gets block number of transaction
func (c *ContractCaller) GetBlockNumberFromTxHash(tx common.Hash) (*big.Int, error) {
	var rpcTx rpcTransaction
	if err := c.MainChainRPC.CallContext(context.Background(), &rpcTx, "eth_getTransactionByHash", tx); err != nil {
		return nil, err
	}

	if rpcTx.BlockNumber == nil {
		return nil, errors.New("no tx found")
	}

	blkNum := big.NewInt(0)

	blkNum, ok := blkNum.SetString(*rpcTx.BlockNumber, 0)
	if !ok {
		return nil, errors.New("unable to set string")
	}

	return blkNum, nil
}

// IsTxConfirmed is tx confirmed
func (c *ContractCaller) IsTxConfirmed(tx common.Hash, requiredConfirmations uint64) bool {
	// get main tx receipt
	receipt, err := c.GetConfirmedTxReceipt(tx, requiredConfirmations)
	if receipt == nil || err != nil {
		return false
	}

	return true
}

// GetConfirmedTxReceipt returns confirmed tx receipt
func (c *ContractCaller) GetConfirmedTxReceipt(tx common.Hash, requiredConfirmations uint64) (*ethTypes.Receipt, error) {
	var receipt *ethTypes.Receipt

	receiptCache, ok := c.ReceiptCache.Get(tx.String())
	if !ok {
		var err error

		// get main tx receipt
		receipt, err = c.GetMainTxReceipt(tx)
		if err != nil {
			Logger.Error("Error while fetching mainChain receipt", "txHash", tx.Hex(), "error", err)
			return nil, err
		}

		c.ReceiptCache.Add(tx.String(), receipt)
	} else {
		receipt, _ = receiptCache.(*ethTypes.Receipt)
	}

	receiptBlockNumber := receipt.BlockNumber.Uint64()

	Logger.Debug("Tx included in block", "block", receiptBlockNumber, "tx", tx)

	// fetch the last finalized main chain block (available post-merge)
	latestFinalizedBlock, err := c.GetMainChainFinalizedBlock()
	if err != nil {
		Logger.Error("error getting latest finalized block from main chain", "error", err)
	}

	// If latest finalized block is available, use it to check if receipt is finalized or not.
	// Else, fallback to the `requiredConfirmations` value
	if latestFinalizedBlock != nil {
		Logger.Debug("Latest finalized block on main chain obtained", "Block", latestFinalizedBlock.Number.Uint64(), "receipt block", receiptBlockNumber)

		if receiptBlockNumber > latestFinalizedBlock.Number.Uint64() {
			return nil, errors.New("not enough confirmations")
		}
	} else {
		// get current/latest main chain block
		latestBlk, err := c.GetMainChainBlock(nil)
		if err != nil {
			Logger.Error("error getting latest block from main chain", "error", err)
			return nil, err
		}

		Logger.Debug("Latest block on main chain obtained", "Block", latestBlk.Number.Uint64(), "receipt block", receiptBlockNumber)

		diff := latestBlk.Number.Uint64() - receiptBlockNumber
		if diff < requiredConfirmations {
			return nil, errors.New("not enough confirmations")
		}
	}

	return receipt, nil
}

//
// Validator decode events
//

// DecodeNewHeaderBlockEvent represents new header block event
func (c *ContractCaller) DecodeNewHeaderBlockEvent(contractAddress common.Address, receipt *ethTypes.Receipt, logIndex uint64) (*rootchain.RootchainNewHeaderBlock, error) {
	event := new(rootchain.RootchainNewHeaderBlock)

	found := false

	for _, vLog := range receipt.Logs {
		if uint64(vLog.Index) == logIndex && bytes.Equal(vLog.Address.Bytes(), contractAddress.Bytes()) {
			found = true

			if err := UnpackLog(&c.RootChainABI, event, newHeaderBlockEvent, vLog); err != nil {
				return nil, err
			}

			break
		}
	}

	if !found {
		return nil, errors.New("event not found")
	}

	return event, nil
}

// DecodeValidatorTopupFeesEvent represents topUp for fees tokens
func (c *ContractCaller) DecodeValidatorTopupFeesEvent(contractAddress common.Address, receipt *ethTypes.Receipt, logIndex uint64) (*stakinginfo.StakinginfoTopUpFee, error) {
	var (
		event = new(stakinginfo.StakinginfoTopUpFee)
		found = false
	)

	for _, vLog := range receipt.Logs {
		if uint64(vLog.Index) == logIndex && bytes.Equal(vLog.Address.Bytes(), contractAddress.Bytes()) {
			found = true

			if err := UnpackLog(&c.StakingInfoABI, event, topUpFeeEvent, vLog); err != nil {
				return nil, err
			}

			break
		}
	}

	if !found {
		return nil, errors.New("event not found")
	}

	return event, nil
}

// DecodeValidatorJoinEvent represents validator staked event
func (c *ContractCaller) DecodeValidatorJoinEvent(contractAddress common.Address, receipt *ethTypes.Receipt, logIndex uint64) (*stakinginfo.StakinginfoStaked, error) {
	event := new(stakinginfo.StakinginfoStaked)

	found := false

	for _, vLog := range receipt.Logs {
		if uint64(vLog.Index) == logIndex && bytes.Equal(vLog.Address.Bytes(), contractAddress.Bytes()) {
			found = true

			if err := UnpackLog(&c.StakingInfoABI, event, stakedEvent, vLog); err != nil {
				return nil, err
			}

			break
		}
	}

	if !found {
		return nil, errors.New("event not found")
	}

	return event, nil
}

// DecodeValidatorStakeUpdateEvent represents validator stake update event
func (c *ContractCaller) DecodeValidatorStakeUpdateEvent(contractAddress common.Address, receipt *ethTypes.Receipt, logIndex uint64) (*stakinginfo.StakinginfoStakeUpdate, error) {
	var (
		event = new(stakinginfo.StakinginfoStakeUpdate)
		found = false
	)

	for _, vLog := range receipt.Logs {
		if uint64(vLog.Index) == logIndex && bytes.Equal(vLog.Address.Bytes(), contractAddress.Bytes()) {
			found = true

			if err := UnpackLog(&c.StakingInfoABI, event, stakeUpdateEvent, vLog); err != nil {
				return nil, err
			}

			break
		}
	}

	if !found {
		return nil, errors.New("event not found")
	}

	return event, nil
}

// DecodeValidatorExitEvent represents validator stake unStake event
func (c *ContractCaller) DecodeValidatorExitEvent(contractAddress common.Address, receipt *ethTypes.Receipt, logIndex uint64) (*stakinginfo.StakinginfoUnstakeInit, error) {
	var (
		event = new(stakinginfo.StakinginfoUnstakeInit)
		found = false
	)

	for _, vLog := range receipt.Logs {
		if uint64(vLog.Index) == logIndex && bytes.Equal(vLog.Address.Bytes(), contractAddress.Bytes()) {
			found = true

			if err := UnpackLog(&c.StakingInfoABI, event, UnstakeInitEvent, vLog); err != nil {
				return nil, err
			}

			break
		}
	}

	if !found {
		return nil, errors.New("event not found")
	}

	return event, nil
}

// DecodeSignerUpdateEvent represents sig update event
func (c *ContractCaller) DecodeSignerUpdateEvent(contractAddress common.Address, receipt *ethTypes.Receipt, logIndex uint64) (*stakinginfo.StakinginfoSignerChange, error) {
	var (
		event = new(stakinginfo.StakinginfoSignerChange)
		found = false
	)

	for _, vLog := range receipt.Logs {
		if uint64(vLog.Index) == logIndex && bytes.Equal(vLog.Address.Bytes(), contractAddress.Bytes()) {
			found = true

			if err := UnpackLog(&c.StakingInfoABI, event, signerChangeEvent, vLog); err != nil {
				return nil, err
			}

			break
		}
	}

	if !found {
		return nil, errors.New("event not found")
	}

	return event, nil
}

// DecodeStateSyncedEvent decode state sync data
func (c *ContractCaller) DecodeStateSyncedEvent(contractAddress common.Address, receipt *ethTypes.Receipt, logIndex uint64) (*statesender.StatesenderStateSynced, error) {
	var (
		event = new(statesender.StatesenderStateSynced)
		found = false
	)

	for _, vLog := range receipt.Logs {
		if uint64(vLog.Index) == logIndex && bytes.Equal(vLog.Address.Bytes(), contractAddress.Bytes()) {
			found = true

			if err := UnpackLog(&c.StateSenderABI, event, stateSyncedEvent, vLog); err != nil {
				return nil, err
			}

			break
		}
	}

	if !found {
		return nil, errors.New("event not found")
	}

	return event, nil
}

// decode slashing events

// DecodeSlashedEvent represents tick ack on contract
func (c *ContractCaller) DecodeSlashedEvent(contractAddress common.Address, receipt *ethTypes.Receipt, logIndex uint64) (*stakinginfo.StakinginfoSlashed, error) {
	var (
		event = new(stakinginfo.StakinginfoSlashed)
		found = false
	)

	for _, vLog := range receipt.Logs {
		if uint64(vLog.Index) == logIndex && bytes.Equal(vLog.Address.Bytes(), contractAddress.Bytes()) {
			found = true

			if err := UnpackLog(&c.StakingInfoABI, event, slashedEvent, vLog); err != nil {
				return nil, err
			}

			break
		}
	}

	if !found {
		return nil, errors.New("event not found")
	}

	return event, nil
}

// DecodeUnJailedEvent represents unJail on contract
func (c *ContractCaller) DecodeUnJailedEvent(contractAddress common.Address, receipt *ethTypes.Receipt, logIndex uint64) (*stakinginfo.StakinginfoUnJailed, error) {
	var (
		event = new(stakinginfo.StakinginfoUnJailed)
		found = false
	)

	for _, vLog := range receipt.Logs {
		if uint64(vLog.Index) == logIndex && bytes.Equal(vLog.Address.Bytes(), contractAddress.Bytes()) {
			found = true

			if err := UnpackLog(&c.StakingInfoABI, event, unJailedEvent, vLog); err != nil {
				return nil, err
			}

			break
		}
	}

	if !found {
		return nil, errors.New("event not found")
	}

	return event, nil
}

//
// Account root related functions
//

// CurrentAccountStateRoot get current account root from on chain
func (c *ContractCaller) CurrentAccountStateRoot(stakingInfoInstance *stakinginfo.Stakinginfo) ([32]byte, error) {
	accountStateRoot, err := stakingInfoInstance.GetAccountStateRoot(nil)

	if err != nil {
		Logger.Error("Unable to get current account state root", "error", err)

		var emptyArr [32]byte

		return emptyArr, err
	}

	return accountStateRoot, nil
}

//
// Span related functions
//

// CurrentSpanNumber get current span
func (c *ContractCaller) CurrentSpanNumber(validatorSetInstance *validatorset.Validatorset) (Number *big.Int) {
	result, err := validatorSetInstance.CurrentSpanNumber(nil)
	if err != nil {
		Logger.Error("Unable to get current span number", "error", err)
		return nil
	}

	return result
}

// GetSpanDetails get span details
func (c *ContractCaller) GetSpanDetails(id *big.Int, validatorSetInstance *validatorset.Validatorset) (
	*big.Int,
	*big.Int,
	*big.Int,
	error,
) {
	d, err := validatorSetInstance.GetSpan(nil, id)
	return d.Number, d.StartBlock, d.EndBlock, err
}

// CurrentStateCounter get state counter
func (c *ContractCaller) CurrentStateCounter(stateSenderInstance *statesender.Statesender) (Number *big.Int) {
	result, err := stateSenderInstance.Counter(nil)
	if err != nil {
		Logger.Error("Unable to get current counter number", "error", err)
		return nil
	}

	return result
}

// CheckIfBlocksExist - check if the given block exists on local chain
func (c *ContractCaller) CheckIfBlocksExist(end uint64) bool {
	// Get block by number.
	var block *ethTypes.Header

	err := c.MaticChainRPC.Call(&block, "eth_getBlockByNumber", fmt.Sprintf("0x%x", end), false)
	if err != nil {
		return false
	}

	return end == block.Number.Uint64()
}

//
// Receipt functions
//

// GetMainTxReceipt returns main tx receipt
func (c *ContractCaller) GetMainTxReceipt(txHash common.Hash) (*ethTypes.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MainChainTimeout)
	defer cancel()

	return c.getTxReceipt(ctx, c.MainChainClient, txHash)
}

// GetMaticTxReceipt returns matic tx receipt
func (c *ContractCaller) GetMaticTxReceipt(txHash common.Hash) (*ethTypes.Receipt, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MaticChainTimeout)
	defer cancel()

	return c.getTxReceipt(ctx, c.MaticChainClient, txHash)
}

func (c *ContractCaller) getTxReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*ethTypes.Receipt, error) {
	return client.TransactionReceipt(ctx, txHash)
}

// GetCheckpointSign returns sigs input of committed checkpoint transaction
func (c *ContractCaller) GetCheckpointSign(txHash common.Hash) ([]byte, []byte, []byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), c.MainChainTimeout)
	defer cancel()

	mainChainClient := GetMainClient()

	transaction, isPending, err := mainChainClient.TransactionByHash(ctx, txHash)
	if err != nil {
		Logger.Error("Error while Fetching Transaction By hash from MainChain", "error", err)
		return []byte{}, []byte{}, []byte{}, err
	} else if isPending {
		return []byte{}, []byte{}, []byte{}, errors.New("transaction is still pending")
	}

	payload := transaction.Data()
	chainABI := c.RootChainABI

	return UnpackSigAndVotes(payload, chainABI)
}

// utility and helper methods

// populateABIs fills the package level cache for contracts' ABIs
// When called the first time, ContractsABIsMap will be filled and getABI method won't be invoked the next times
// This reduces the number of calls to json decode methods made by the contract caller
// It uses ABIs' definitions instead of contracts addresses, as the latter might not be available at init time
func populateABIs(contractCallerObj *ContractCaller) error {
	var ccAbi *abi.ABI

	var err error

	contractsABIs := [8]string{rootchain.RootchainABI, stakinginfo.StakinginfoABI, validatorset.ValidatorsetABI,
		statereceiver.StatereceiverABI, statesender.StatesenderABI, stakemanager.StakemanagerABI,
		slashmanager.SlashmanagerABI, erc20.Erc20ABI}

	// iterate over supported ABIs
	for _, contractABI := range contractsABIs {
		ccAbi, err = chooseContractCallerABI(contractCallerObj, contractABI)
		if err != nil {
			Logger.Error("Error while fetching contract caller ABI", "error", err)
			return err
		}

		if ContractsABIsMap[contractABI] == nil {
			// fills cached abi map
			if *ccAbi, err = getABI(contractABI); err != nil {
				Logger.Error("Error while getting ABI for contract caller", "name", contractABI, "error", err)
				return err
			} else {
				// init ABI
				ContractsABIsMap[contractABI] = ccAbi
			}
		} else {
			// use cached abi
			*ccAbi = *ContractsABIsMap[contractABI]
		}
	}

	return nil
}

// chooseContractCallerABI extracts and returns the abo.ABI object from the contractCallerObj based on its abi string
func chooseContractCallerABI(contractCallerObj *ContractCaller, abi string) (*abi.ABI, error) {
	switch abi {
	case rootchain.RootchainABI:
		return &contractCallerObj.RootChainABI, nil
	case stakinginfo.StakinginfoABI:
		return &contractCallerObj.StakingInfoABI, nil
	case validatorset.ValidatorsetABI:
		return &contractCallerObj.ValidatorSetABI, nil
	case statereceiver.StatereceiverABI:
		return &contractCallerObj.StateReceiverABI, nil
	case statesender.StatesenderABI:
		return &contractCallerObj.StateSenderABI, nil
	case stakemanager.StakemanagerABI:
		return &contractCallerObj.StakeManagerABI, nil
	case slashmanager.SlashmanagerABI:
		return &contractCallerObj.SlashManagerABI, nil
	case erc20.Erc20ABI:
		return &contractCallerObj.MaticTokenABI, nil
	}

	return nil, errors.New("no ABI associated with such data")
}

// getABI returns the contract's ABI struct from on its JSON representation
func getABI(data string) (abi.ABI, error) {
	return abi.JSON(strings.NewReader(data))
}
