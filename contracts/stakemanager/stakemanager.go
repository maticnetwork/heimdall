// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stakemanager

import (
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// StakemanagerABI is the input ABI used to generate the binding from.
const StakemanagerABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"WITHDRAWAL_DELAY\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MIN_DEPOSIT_SIZE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"validatorThreshold\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validatorState\",\"outputs\":[{\"name\":\"amount\",\"type\":\"int256\"},{\"name\":\"stakerCount\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"dynasty\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalStaked\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"stakers\",\"outputs\":[{\"name\":\"epoch\",\"type\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"activationEpoch\",\"type\":\"uint256\"},{\"name\":\"deactivationEpoch\",\"type\":\"uint256\"},{\"name\":\"signer\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"rootChain\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minLockInPeriod\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"EPOCH_LENGTH\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"signerToStaker\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"maxStakeDrop\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"locked\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newRootChain\",\"type\":\"address\"}],\"name\":\"changeRootChain\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"lock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"oldThreshold\",\"type\":\"uint256\"}],\"name\":\"ThresholdChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newDynasty\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"oldDynasty\",\"type\":\"uint256\"}],\"name\":\"DynastyValueChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"deactivationEpoch\",\"type\":\"uint256\"}],\"name\":\"UnstakeInit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newSigner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"oldSigner\",\"type\":\"address\"}],\"name\":\"SignerChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousRootChain\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newRootChain\",\"type\":\"address\"}],\"name\":\"RootChainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"}],\"name\":\"OwnershipRenounced\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"activatonEpoch\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"Unstaked\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"unstakeValidator\",\"type\":\"address\"},{\"name\":\"signer\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"user\",\"type\":\"address\"},{\"name\":\"unstakeValidator\",\"type\":\"address\"},{\"name\":\"signer\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"stakeFor\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unstake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unstakeClaim\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getCurrentValidatorSet\",\"outputs\":[{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getStakerDetails\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"totalStakedFor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"supportsHistory\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"updateValidatorThreshold\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newDynasty\",\"type\":\"uint256\"}],\"name\":\"updateDynastyValue\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_signer\",\"type\":\"address\"}],\"name\":\"updateSigner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeCommit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"epochs\",\"type\":\"uint256\"}],\"name\":\"updateMinLockInPeriod\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentValidatorSetSize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentValidatorSetTotalStake\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"user\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"voteHash\",\"type\":\"bytes32\"},{\"name\":\"sigs\",\"type\":\"bytes\"}],\"name\":\"checkSignatures\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Stakemanager is an auto generated Go binding around an Ethereum contract.
type Stakemanager struct {
	StakemanagerCaller     // Read-only binding to the contract
	StakemanagerTransactor // Write-only binding to the contract
	StakemanagerFilterer   // Log filterer for contract events
}

// StakemanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakemanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakemanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakemanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakemanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakemanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakemanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakemanagerSession struct {
	Contract     *Stakemanager     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakemanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakemanagerCallerSession struct {
	Contract *StakemanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// StakemanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakemanagerTransactorSession struct {
	Contract     *StakemanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// StakemanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakemanagerRaw struct {
	Contract *Stakemanager // Generic contract binding to access the raw methods on
}

// StakemanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakemanagerCallerRaw struct {
	Contract *StakemanagerCaller // Generic read-only contract binding to access the raw methods on
}

// StakemanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakemanagerTransactorRaw struct {
	Contract *StakemanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakemanager creates a new instance of Stakemanager, bound to a specific deployed contract.
func NewStakemanager(address common.Address, backend bind.ContractBackend) (*Stakemanager, error) {
	contract, err := bindStakemanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Stakemanager{StakemanagerCaller: StakemanagerCaller{contract: contract}, StakemanagerTransactor: StakemanagerTransactor{contract: contract}, StakemanagerFilterer: StakemanagerFilterer{contract: contract}}, nil
}

// NewStakemanagerCaller creates a new read-only instance of Stakemanager, bound to a specific deployed contract.
func NewStakemanagerCaller(address common.Address, caller bind.ContractCaller) (*StakemanagerCaller, error) {
	contract, err := bindStakemanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakemanagerCaller{contract: contract}, nil
}

// NewStakemanagerTransactor creates a new write-only instance of Stakemanager, bound to a specific deployed contract.
func NewStakemanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*StakemanagerTransactor, error) {
	contract, err := bindStakemanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakemanagerTransactor{contract: contract}, nil
}

// NewStakemanagerFilterer creates a new log filterer instance of Stakemanager, bound to a specific deployed contract.
func NewStakemanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*StakemanagerFilterer, error) {
	contract, err := bindStakemanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakemanagerFilterer{contract: contract}, nil
}

// bindStakemanager binds a generic wrapper to an already deployed contract.
func bindStakemanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakemanagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Stakemanager *StakemanagerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Stakemanager.Contract.StakemanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Stakemanager *StakemanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakemanager.Contract.StakemanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Stakemanager *StakemanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stakemanager.Contract.StakemanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Stakemanager *StakemanagerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Stakemanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Stakemanager *StakemanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakemanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Stakemanager *StakemanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stakemanager.Contract.contract.Transact(opts, method, params...)
}

// EPOCHLENGTH is a free data retrieval call binding the contract method 0xac4746ab.
//
// Solidity: function EPOCH_LENGTH() constant returns(uint256)
func (_Stakemanager *StakemanagerCaller) EPOCHLENGTH(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "EPOCH_LENGTH")
	return *ret0, err
}

// EPOCHLENGTH is a free data retrieval call binding the contract method 0xac4746ab.
//
// Solidity: function EPOCH_LENGTH() constant returns(uint256)
func (_Stakemanager *StakemanagerSession) EPOCHLENGTH() (*big.Int, error) {
	return _Stakemanager.Contract.EPOCHLENGTH(&_Stakemanager.CallOpts)
}

// EPOCHLENGTH is a free data retrieval call binding the contract method 0xac4746ab.
//
// Solidity: function EPOCH_LENGTH() constant returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) EPOCHLENGTH() (*big.Int, error) {
	return _Stakemanager.Contract.EPOCHLENGTH(&_Stakemanager.CallOpts)
}

// MINDEPOSITSIZE is a free data retrieval call binding the contract method 0x26c0817e.
//
// Solidity: function MIN_DEPOSIT_SIZE() constant returns(uint256)
func (_Stakemanager *StakemanagerCaller) MINDEPOSITSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "MIN_DEPOSIT_SIZE")
	return *ret0, err
}

// MINDEPOSITSIZE is a free data retrieval call binding the contract method 0x26c0817e.
//
// Solidity: function MIN_DEPOSIT_SIZE() constant returns(uint256)
func (_Stakemanager *StakemanagerSession) MINDEPOSITSIZE() (*big.Int, error) {
	return _Stakemanager.Contract.MINDEPOSITSIZE(&_Stakemanager.CallOpts)
}

// MINDEPOSITSIZE is a free data retrieval call binding the contract method 0x26c0817e.
//
// Solidity: function MIN_DEPOSIT_SIZE() constant returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) MINDEPOSITSIZE() (*big.Int, error) {
	return _Stakemanager.Contract.MINDEPOSITSIZE(&_Stakemanager.CallOpts)
}

// WITHDRAWALDELAY is a free data retrieval call binding the contract method 0x0ebb172a.
//
// Solidity: function WITHDRAWAL_DELAY() constant returns(uint256)
func (_Stakemanager *StakemanagerCaller) WITHDRAWALDELAY(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "WITHDRAWAL_DELAY")
	return *ret0, err
}

// WITHDRAWALDELAY is a free data retrieval call binding the contract method 0x0ebb172a.
//
// Solidity: function WITHDRAWAL_DELAY() constant returns(uint256)
func (_Stakemanager *StakemanagerSession) WITHDRAWALDELAY() (*big.Int, error) {
	return _Stakemanager.Contract.WITHDRAWALDELAY(&_Stakemanager.CallOpts)
}

// WITHDRAWALDELAY is a free data retrieval call binding the contract method 0x0ebb172a.
//
// Solidity: function WITHDRAWAL_DELAY() constant returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) WITHDRAWALDELAY() (*big.Int, error) {
	return _Stakemanager.Contract.WITHDRAWALDELAY(&_Stakemanager.CallOpts)
}

// CheckSignatures is a free data retrieval call binding the contract method 0xed516d51.
//
// Solidity: function checkSignatures(voteHash bytes32, sigs bytes) constant returns(bool)
func (_Stakemanager *StakemanagerCaller) CheckSignatures(opts *bind.CallOpts, voteHash [32]byte, sigs []byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "checkSignatures", voteHash, sigs)
	return *ret0, err
}

// CheckSignatures is a free data retrieval call binding the contract method 0xed516d51.
//
// Solidity: function checkSignatures(voteHash bytes32, sigs bytes) constant returns(bool)
func (_Stakemanager *StakemanagerSession) CheckSignatures(voteHash [32]byte, sigs []byte) (bool, error) {
	return _Stakemanager.Contract.CheckSignatures(&_Stakemanager.CallOpts, voteHash, sigs)
}

// CheckSignatures is a free data retrieval call binding the contract method 0xed516d51.
//
// Solidity: function checkSignatures(voteHash bytes32, sigs bytes) constant returns(bool)
func (_Stakemanager *StakemanagerCallerSession) CheckSignatures(voteHash [32]byte, sigs []byte) (bool, error) {
	return _Stakemanager.Contract.CheckSignatures(&_Stakemanager.CallOpts, voteHash, sigs)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(uint256)
func (_Stakemanager *StakemanagerCaller) CurrentEpoch(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "currentEpoch")
	return *ret0, err
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(uint256)
func (_Stakemanager *StakemanagerSession) CurrentEpoch() (*big.Int, error) {
	return _Stakemanager.Contract.CurrentEpoch(&_Stakemanager.CallOpts)
}

// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
//
// Solidity: function currentEpoch() constant returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) CurrentEpoch() (*big.Int, error) {
	return _Stakemanager.Contract.CurrentEpoch(&_Stakemanager.CallOpts)
}

// CurrentValidatorSetSize is a free data retrieval call binding the contract method 0x7f952d95.
//
// Solidity: function currentValidatorSetSize() constant returns(uint256)
func (_Stakemanager *StakemanagerCaller) CurrentValidatorSetSize(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "currentValidatorSetSize")
	return *ret0, err
}

// CurrentValidatorSetSize is a free data retrieval call binding the contract method 0x7f952d95.
//
// Solidity: function currentValidatorSetSize() constant returns(uint256)
func (_Stakemanager *StakemanagerSession) CurrentValidatorSetSize() (*big.Int, error) {
	return _Stakemanager.Contract.CurrentValidatorSetSize(&_Stakemanager.CallOpts)
}

// CurrentValidatorSetSize is a free data retrieval call binding the contract method 0x7f952d95.
//
// Solidity: function currentValidatorSetSize() constant returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) CurrentValidatorSetSize() (*big.Int, error) {
	return _Stakemanager.Contract.CurrentValidatorSetSize(&_Stakemanager.CallOpts)
}

// CurrentValidatorSetTotalStake is a free data retrieval call binding the contract method 0xa4769071.
//
// Solidity: function currentValidatorSetTotalStake() constant returns(uint256)
func (_Stakemanager *StakemanagerCaller) CurrentValidatorSetTotalStake(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "currentValidatorSetTotalStake")
	return *ret0, err
}

// CurrentValidatorSetTotalStake is a free data retrieval call binding the contract method 0xa4769071.
//
// Solidity: function currentValidatorSetTotalStake() constant returns(uint256)
func (_Stakemanager *StakemanagerSession) CurrentValidatorSetTotalStake() (*big.Int, error) {
	return _Stakemanager.Contract.CurrentValidatorSetTotalStake(&_Stakemanager.CallOpts)
}

// CurrentValidatorSetTotalStake is a free data retrieval call binding the contract method 0xa4769071.
//
// Solidity: function currentValidatorSetTotalStake() constant returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) CurrentValidatorSetTotalStake() (*big.Int, error) {
	return _Stakemanager.Contract.CurrentValidatorSetTotalStake(&_Stakemanager.CallOpts)
}

// Dynasty is a free data retrieval call binding the contract method 0x7060054d.
//
// Solidity: function dynasty() constant returns(uint256)
func (_Stakemanager *StakemanagerCaller) Dynasty(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "dynasty")
	return *ret0, err
}

// Dynasty is a free data retrieval call binding the contract method 0x7060054d.
//
// Solidity: function dynasty() constant returns(uint256)
func (_Stakemanager *StakemanagerSession) Dynasty() (*big.Int, error) {
	return _Stakemanager.Contract.Dynasty(&_Stakemanager.CallOpts)
}

// Dynasty is a free data retrieval call binding the contract method 0x7060054d.
//
// Solidity: function dynasty() constant returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) Dynasty() (*big.Int, error) {
	return _Stakemanager.Contract.Dynasty(&_Stakemanager.CallOpts)
}

// GetCurrentValidatorSet is a free data retrieval call binding the contract method 0x0209fdd0.
//
// Solidity: function getCurrentValidatorSet() constant returns(address[])
func (_Stakemanager *StakemanagerCaller) GetCurrentValidatorSet(opts *bind.CallOpts) ([]common.Address, error) {
	var (
		ret0 = new([]common.Address)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "getCurrentValidatorSet")
	return *ret0, err
}

// GetCurrentValidatorSet is a free data retrieval call binding the contract method 0x0209fdd0.
//
// Solidity: function getCurrentValidatorSet() constant returns(address[])
func (_Stakemanager *StakemanagerSession) GetCurrentValidatorSet() ([]common.Address, error) {
	return _Stakemanager.Contract.GetCurrentValidatorSet(&_Stakemanager.CallOpts)
}

// GetCurrentValidatorSet is a free data retrieval call binding the contract method 0x0209fdd0.
//
// Solidity: function getCurrentValidatorSet() constant returns(address[])
func (_Stakemanager *StakemanagerCallerSession) GetCurrentValidatorSet() ([]common.Address, error) {
	return _Stakemanager.Contract.GetCurrentValidatorSet(&_Stakemanager.CallOpts)
}

// GetStakerDetails is a free data retrieval call binding the contract method 0x8fb80c73.
//
// Solidity: function getStakerDetails(user address) constant returns(uint256, uint256, uint256, address)
func (_Stakemanager *StakemanagerCaller) GetStakerDetails(opts *bind.CallOpts, user common.Address) (*big.Int, *big.Int, *big.Int, common.Address, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
		ret3 = new(common.Address)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
	}
	err := _Stakemanager.contract.Call(opts, out, "getStakerDetails", user)
	return *ret0, *ret1, *ret2, *ret3, err
}

// GetStakerDetails is a free data retrieval call binding the contract method 0x8fb80c73.
//
// Solidity: function getStakerDetails(user address) constant returns(uint256, uint256, uint256, address)
func (_Stakemanager *StakemanagerSession) GetStakerDetails(user common.Address) (*big.Int, *big.Int, *big.Int, common.Address, error) {
	return _Stakemanager.Contract.GetStakerDetails(&_Stakemanager.CallOpts, user)
}

// GetStakerDetails is a free data retrieval call binding the contract method 0x8fb80c73.
//
// Solidity: function getStakerDetails(user address) constant returns(uint256, uint256, uint256, address)
func (_Stakemanager *StakemanagerCallerSession) GetStakerDetails(user common.Address) (*big.Int, *big.Int, *big.Int, common.Address, error) {
	return _Stakemanager.Contract.GetStakerDetails(&_Stakemanager.CallOpts, user)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(user address) constant returns(bool)
func (_Stakemanager *StakemanagerCaller) IsValidator(opts *bind.CallOpts, user common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "isValidator", user)
	return *ret0, err
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(user address) constant returns(bool)
func (_Stakemanager *StakemanagerSession) IsValidator(user common.Address) (bool, error) {
	return _Stakemanager.Contract.IsValidator(&_Stakemanager.CallOpts, user)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(user address) constant returns(bool)
func (_Stakemanager *StakemanagerCallerSession) IsValidator(user common.Address) (bool, error) {
	return _Stakemanager.Contract.IsValidator(&_Stakemanager.CallOpts, user)
}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() constant returns(bool)
func (_Stakemanager *StakemanagerCaller) Locked(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "locked")
	return *ret0, err
}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() constant returns(bool)
func (_Stakemanager *StakemanagerSession) Locked() (bool, error) {
	return _Stakemanager.Contract.Locked(&_Stakemanager.CallOpts)
}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() constant returns(bool)
func (_Stakemanager *StakemanagerCallerSession) Locked() (bool, error) {
	return _Stakemanager.Contract.Locked(&_Stakemanager.CallOpts)
}

// MaxStakeDrop is a free data retrieval call binding the contract method 0xbbce8cec.
//
// Solidity: function maxStakeDrop() constant returns(uint256)
func (_Stakemanager *StakemanagerCaller) MaxStakeDrop(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "maxStakeDrop")
	return *ret0, err
}

// MaxStakeDrop is a free data retrieval call binding the contract method 0xbbce8cec.
//
// Solidity: function maxStakeDrop() constant returns(uint256)
func (_Stakemanager *StakemanagerSession) MaxStakeDrop() (*big.Int, error) {
	return _Stakemanager.Contract.MaxStakeDrop(&_Stakemanager.CallOpts)
}

// MaxStakeDrop is a free data retrieval call binding the contract method 0xbbce8cec.
//
// Solidity: function maxStakeDrop() constant returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) MaxStakeDrop() (*big.Int, error) {
	return _Stakemanager.Contract.MaxStakeDrop(&_Stakemanager.CallOpts)
}

// MinLockInPeriod is a free data retrieval call binding the contract method 0xa548c547.
//
// Solidity: function minLockInPeriod() constant returns(uint256)
func (_Stakemanager *StakemanagerCaller) MinLockInPeriod(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "minLockInPeriod")
	return *ret0, err
}

// MinLockInPeriod is a free data retrieval call binding the contract method 0xa548c547.
//
// Solidity: function minLockInPeriod() constant returns(uint256)
func (_Stakemanager *StakemanagerSession) MinLockInPeriod() (*big.Int, error) {
	return _Stakemanager.Contract.MinLockInPeriod(&_Stakemanager.CallOpts)
}

// MinLockInPeriod is a free data retrieval call binding the contract method 0xa548c547.
//
// Solidity: function minLockInPeriod() constant returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) MinLockInPeriod() (*big.Int, error) {
	return _Stakemanager.Contract.MinLockInPeriod(&_Stakemanager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Stakemanager *StakemanagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Stakemanager *StakemanagerSession) Owner() (common.Address, error) {
	return _Stakemanager.Contract.Owner(&_Stakemanager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Stakemanager *StakemanagerCallerSession) Owner() (common.Address, error) {
	return _Stakemanager.Contract.Owner(&_Stakemanager.CallOpts)
}

// RootChain is a free data retrieval call binding the contract method 0x987ab9db.
//
// Solidity: function rootChain() constant returns(address)
func (_Stakemanager *StakemanagerCaller) RootChain(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "rootChain")
	return *ret0, err
}

// RootChain is a free data retrieval call binding the contract method 0x987ab9db.
//
// Solidity: function rootChain() constant returns(address)
func (_Stakemanager *StakemanagerSession) RootChain() (common.Address, error) {
	return _Stakemanager.Contract.RootChain(&_Stakemanager.CallOpts)
}

// RootChain is a free data retrieval call binding the contract method 0x987ab9db.
//
// Solidity: function rootChain() constant returns(address)
func (_Stakemanager *StakemanagerCallerSession) RootChain() (common.Address, error) {
	return _Stakemanager.Contract.RootChain(&_Stakemanager.CallOpts)
}

// SignerToStaker is a free data retrieval call binding the contract method 0xad5a98c5.
//
// Solidity: function signerToStaker( address) constant returns(address)
func (_Stakemanager *StakemanagerCaller) SignerToStaker(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "signerToStaker", arg0)
	return *ret0, err
}

// SignerToStaker is a free data retrieval call binding the contract method 0xad5a98c5.
//
// Solidity: function signerToStaker( address) constant returns(address)
func (_Stakemanager *StakemanagerSession) SignerToStaker(arg0 common.Address) (common.Address, error) {
	return _Stakemanager.Contract.SignerToStaker(&_Stakemanager.CallOpts, arg0)
}

// SignerToStaker is a free data retrieval call binding the contract method 0xad5a98c5.
//
// Solidity: function signerToStaker( address) constant returns(address)
func (_Stakemanager *StakemanagerCallerSession) SignerToStaker(arg0 common.Address) (common.Address, error) {
	return _Stakemanager.Contract.SignerToStaker(&_Stakemanager.CallOpts, arg0)
}

// Stakers is a free data retrieval call binding the contract method 0x9168ae72.
//
// Solidity: function stakers( address) constant returns(epoch uint256, amount uint256, activationEpoch uint256, deactivationEpoch uint256, signer address)
func (_Stakemanager *StakemanagerCaller) Stakers(opts *bind.CallOpts, arg0 common.Address) (struct {
	Epoch             *big.Int
	Amount            *big.Int
	ActivationEpoch   *big.Int
	DeactivationEpoch *big.Int
	Signer            common.Address
}, error) {
	ret := new(struct {
		Epoch             *big.Int
		Amount            *big.Int
		ActivationEpoch   *big.Int
		DeactivationEpoch *big.Int
		Signer            common.Address
	})
	out := ret
	err := _Stakemanager.contract.Call(opts, out, "stakers", arg0)
	return *ret, err
}

// Stakers is a free data retrieval call binding the contract method 0x9168ae72.
//
// Solidity: function stakers( address) constant returns(epoch uint256, amount uint256, activationEpoch uint256, deactivationEpoch uint256, signer address)
func (_Stakemanager *StakemanagerSession) Stakers(arg0 common.Address) (struct {
	Epoch             *big.Int
	Amount            *big.Int
	ActivationEpoch   *big.Int
	DeactivationEpoch *big.Int
	Signer            common.Address
}, error) {
	return _Stakemanager.Contract.Stakers(&_Stakemanager.CallOpts, arg0)
}

// Stakers is a free data retrieval call binding the contract method 0x9168ae72.
//
// Solidity: function stakers( address) constant returns(epoch uint256, amount uint256, activationEpoch uint256, deactivationEpoch uint256, signer address)
func (_Stakemanager *StakemanagerCallerSession) Stakers(arg0 common.Address) (struct {
	Epoch             *big.Int
	Amount            *big.Int
	ActivationEpoch   *big.Int
	DeactivationEpoch *big.Int
	Signer            common.Address
}, error) {
	return _Stakemanager.Contract.Stakers(&_Stakemanager.CallOpts, arg0)
}

// SupportsHistory is a free data retrieval call binding the contract method 0x7033e4a6.
//
// Solidity: function supportsHistory() constant returns(bool)
func (_Stakemanager *StakemanagerCaller) SupportsHistory(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "supportsHistory")
	return *ret0, err
}

// SupportsHistory is a free data retrieval call binding the contract method 0x7033e4a6.
//
// Solidity: function supportsHistory() constant returns(bool)
func (_Stakemanager *StakemanagerSession) SupportsHistory() (bool, error) {
	return _Stakemanager.Contract.SupportsHistory(&_Stakemanager.CallOpts)
}

// SupportsHistory is a free data retrieval call binding the contract method 0x7033e4a6.
//
// Solidity: function supportsHistory() constant returns(bool)
func (_Stakemanager *StakemanagerCallerSession) SupportsHistory() (bool, error) {
	return _Stakemanager.Contract.SupportsHistory(&_Stakemanager.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Stakemanager *StakemanagerCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "token")
	return *ret0, err
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Stakemanager *StakemanagerSession) Token() (common.Address, error) {
	return _Stakemanager.Contract.Token(&_Stakemanager.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Stakemanager *StakemanagerCallerSession) Token() (common.Address, error) {
	return _Stakemanager.Contract.Token(&_Stakemanager.CallOpts)
}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() constant returns(uint256)
func (_Stakemanager *StakemanagerCaller) TotalStaked(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "totalStaked")
	return *ret0, err
}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() constant returns(uint256)
func (_Stakemanager *StakemanagerSession) TotalStaked() (*big.Int, error) {
	return _Stakemanager.Contract.TotalStaked(&_Stakemanager.CallOpts)
}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() constant returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) TotalStaked() (*big.Int, error) {
	return _Stakemanager.Contract.TotalStaked(&_Stakemanager.CallOpts)
}

// TotalStakedFor is a free data retrieval call binding the contract method 0x4b341aed.
//
// Solidity: function totalStakedFor(addr address) constant returns(uint256)
func (_Stakemanager *StakemanagerCaller) TotalStakedFor(opts *bind.CallOpts, addr common.Address) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "totalStakedFor", addr)
	return *ret0, err
}

// TotalStakedFor is a free data retrieval call binding the contract method 0x4b341aed.
//
// Solidity: function totalStakedFor(addr address) constant returns(uint256)
func (_Stakemanager *StakemanagerSession) TotalStakedFor(addr common.Address) (*big.Int, error) {
	return _Stakemanager.Contract.TotalStakedFor(&_Stakemanager.CallOpts, addr)
}

// TotalStakedFor is a free data retrieval call binding the contract method 0x4b341aed.
//
// Solidity: function totalStakedFor(addr address) constant returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) TotalStakedFor(addr common.Address) (*big.Int, error) {
	return _Stakemanager.Contract.TotalStakedFor(&_Stakemanager.CallOpts, addr)
}

// ValidatorState is a free data retrieval call binding the contract method 0x5c248855.
//
// Solidity: function validatorState( uint256) constant returns(amount int256, stakerCount int256)
func (_Stakemanager *StakemanagerCaller) ValidatorState(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Amount      *big.Int
	StakerCount *big.Int
}, error) {
	ret := new(struct {
		Amount      *big.Int
		StakerCount *big.Int
	})
	out := ret
	err := _Stakemanager.contract.Call(opts, out, "validatorState", arg0)
	return *ret, err
}

// ValidatorState is a free data retrieval call binding the contract method 0x5c248855.
//
// Solidity: function validatorState( uint256) constant returns(amount int256, stakerCount int256)
func (_Stakemanager *StakemanagerSession) ValidatorState(arg0 *big.Int) (struct {
	Amount      *big.Int
	StakerCount *big.Int
}, error) {
	return _Stakemanager.Contract.ValidatorState(&_Stakemanager.CallOpts, arg0)
}

// ValidatorState is a free data retrieval call binding the contract method 0x5c248855.
//
// Solidity: function validatorState( uint256) constant returns(amount int256, stakerCount int256)
func (_Stakemanager *StakemanagerCallerSession) ValidatorState(arg0 *big.Int) (struct {
	Amount      *big.Int
	StakerCount *big.Int
}, error) {
	return _Stakemanager.Contract.ValidatorState(&_Stakemanager.CallOpts, arg0)
}

// ValidatorThreshold is a free data retrieval call binding the contract method 0x4fd101d7.
//
// Solidity: function validatorThreshold() constant returns(uint256)
func (_Stakemanager *StakemanagerCaller) ValidatorThreshold(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "validatorThreshold")
	return *ret0, err
}

// ValidatorThreshold is a free data retrieval call binding the contract method 0x4fd101d7.
//
// Solidity: function validatorThreshold() constant returns(uint256)
func (_Stakemanager *StakemanagerSession) ValidatorThreshold() (*big.Int, error) {
	return _Stakemanager.Contract.ValidatorThreshold(&_Stakemanager.CallOpts)
}

// ValidatorThreshold is a free data retrieval call binding the contract method 0x4fd101d7.
//
// Solidity: function validatorThreshold() constant returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) ValidatorThreshold() (*big.Int, error) {
	return _Stakemanager.Contract.ValidatorThreshold(&_Stakemanager.CallOpts)
}

// ChangeRootChain is a paid mutator transaction binding the contract method 0xe8afa8e8.
//
// Solidity: function changeRootChain(newRootChain address) returns()
func (_Stakemanager *StakemanagerTransactor) ChangeRootChain(opts *bind.TransactOpts, newRootChain common.Address) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "changeRootChain", newRootChain)
}

// ChangeRootChain is a paid mutator transaction binding the contract method 0xe8afa8e8.
//
// Solidity: function changeRootChain(newRootChain address) returns()
func (_Stakemanager *StakemanagerSession) ChangeRootChain(newRootChain common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.ChangeRootChain(&_Stakemanager.TransactOpts, newRootChain)
}

// ChangeRootChain is a paid mutator transaction binding the contract method 0xe8afa8e8.
//
// Solidity: function changeRootChain(newRootChain address) returns()
func (_Stakemanager *StakemanagerTransactorSession) ChangeRootChain(newRootChain common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.ChangeRootChain(&_Stakemanager.TransactOpts, newRootChain)
}

// FinalizeCommit is a paid mutator transaction binding the contract method 0x35dda498.
//
// Solidity: function finalizeCommit() returns()
func (_Stakemanager *StakemanagerTransactor) FinalizeCommit(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "finalizeCommit")
}

// FinalizeCommit is a paid mutator transaction binding the contract method 0x35dda498.
//
// Solidity: function finalizeCommit() returns()
func (_Stakemanager *StakemanagerSession) FinalizeCommit() (*types.Transaction, error) {
	return _Stakemanager.Contract.FinalizeCommit(&_Stakemanager.TransactOpts)
}

// FinalizeCommit is a paid mutator transaction binding the contract method 0x35dda498.
//
// Solidity: function finalizeCommit() returns()
func (_Stakemanager *StakemanagerTransactorSession) FinalizeCommit() (*types.Transaction, error) {
	return _Stakemanager.Contract.FinalizeCommit(&_Stakemanager.TransactOpts)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_Stakemanager *StakemanagerTransactor) Lock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "lock")
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_Stakemanager *StakemanagerSession) Lock() (*types.Transaction, error) {
	return _Stakemanager.Contract.Lock(&_Stakemanager.TransactOpts)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_Stakemanager *StakemanagerTransactorSession) Lock() (*types.Transaction, error) {
	return _Stakemanager.Contract.Lock(&_Stakemanager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Stakemanager *StakemanagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Stakemanager *StakemanagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _Stakemanager.Contract.RenounceOwnership(&_Stakemanager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Stakemanager *StakemanagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Stakemanager.Contract.RenounceOwnership(&_Stakemanager.TransactOpts)
}

// Stake is a paid mutator transaction binding the contract method 0xbf6eac2f.
//
// Solidity: function stake(unstakeValidator address, signer address, amount uint256) returns()
func (_Stakemanager *StakemanagerTransactor) Stake(opts *bind.TransactOpts, unstakeValidator common.Address, signer common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "stake", unstakeValidator, signer, amount)
}

// Stake is a paid mutator transaction binding the contract method 0xbf6eac2f.
//
// Solidity: function stake(unstakeValidator address, signer address, amount uint256) returns()
func (_Stakemanager *StakemanagerSession) Stake(unstakeValidator common.Address, signer common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.Stake(&_Stakemanager.TransactOpts, unstakeValidator, signer, amount)
}

// Stake is a paid mutator transaction binding the contract method 0xbf6eac2f.
//
// Solidity: function stake(unstakeValidator address, signer address, amount uint256) returns()
func (_Stakemanager *StakemanagerTransactorSession) Stake(unstakeValidator common.Address, signer common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.Stake(&_Stakemanager.TransactOpts, unstakeValidator, signer, amount)
}

// StakeFor is a paid mutator transaction binding the contract method 0xae3a73fe.
//
// Solidity: function stakeFor(user address, unstakeValidator address, signer address, amount uint256) returns()
func (_Stakemanager *StakemanagerTransactor) StakeFor(opts *bind.TransactOpts, user common.Address, unstakeValidator common.Address, signer common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "stakeFor", user, unstakeValidator, signer, amount)
}

// StakeFor is a paid mutator transaction binding the contract method 0xae3a73fe.
//
// Solidity: function stakeFor(user address, unstakeValidator address, signer address, amount uint256) returns()
func (_Stakemanager *StakemanagerSession) StakeFor(user common.Address, unstakeValidator common.Address, signer common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.StakeFor(&_Stakemanager.TransactOpts, user, unstakeValidator, signer, amount)
}

// StakeFor is a paid mutator transaction binding the contract method 0xae3a73fe.
//
// Solidity: function stakeFor(user address, unstakeValidator address, signer address, amount uint256) returns()
func (_Stakemanager *StakemanagerTransactorSession) StakeFor(user common.Address, unstakeValidator common.Address, signer common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.StakeFor(&_Stakemanager.TransactOpts, user, unstakeValidator, signer, amount)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Stakemanager *StakemanagerTransactor) TransferOwnership(opts *bind.TransactOpts, _newOwner common.Address) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "transferOwnership", _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Stakemanager *StakemanagerSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.TransferOwnership(&_Stakemanager.TransactOpts, _newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(_newOwner address) returns()
func (_Stakemanager *StakemanagerTransactorSession) TransferOwnership(_newOwner common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.TransferOwnership(&_Stakemanager.TransactOpts, _newOwner)
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_Stakemanager *StakemanagerTransactor) Unlock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "unlock")
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_Stakemanager *StakemanagerSession) Unlock() (*types.Transaction, error) {
	return _Stakemanager.Contract.Unlock(&_Stakemanager.TransactOpts)
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_Stakemanager *StakemanagerTransactorSession) Unlock() (*types.Transaction, error) {
	return _Stakemanager.Contract.Unlock(&_Stakemanager.TransactOpts)
}

// Unstake is a paid mutator transaction binding the contract method 0x2def6620.
//
// Solidity: function unstake() returns()
func (_Stakemanager *StakemanagerTransactor) Unstake(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "unstake")
}

// Unstake is a paid mutator transaction binding the contract method 0x2def6620.
//
// Solidity: function unstake() returns()
func (_Stakemanager *StakemanagerSession) Unstake() (*types.Transaction, error) {
	return _Stakemanager.Contract.Unstake(&_Stakemanager.TransactOpts)
}

// Unstake is a paid mutator transaction binding the contract method 0x2def6620.
//
// Solidity: function unstake() returns()
func (_Stakemanager *StakemanagerTransactorSession) Unstake() (*types.Transaction, error) {
	return _Stakemanager.Contract.Unstake(&_Stakemanager.TransactOpts)
}

// UnstakeClaim is a paid mutator transaction binding the contract method 0x12102cc9.
//
// Solidity: function unstakeClaim() returns()
func (_Stakemanager *StakemanagerTransactor) UnstakeClaim(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "unstakeClaim")
}

// UnstakeClaim is a paid mutator transaction binding the contract method 0x12102cc9.
//
// Solidity: function unstakeClaim() returns()
func (_Stakemanager *StakemanagerSession) UnstakeClaim() (*types.Transaction, error) {
	return _Stakemanager.Contract.UnstakeClaim(&_Stakemanager.TransactOpts)
}

// UnstakeClaim is a paid mutator transaction binding the contract method 0x12102cc9.
//
// Solidity: function unstakeClaim() returns()
func (_Stakemanager *StakemanagerTransactorSession) UnstakeClaim() (*types.Transaction, error) {
	return _Stakemanager.Contract.UnstakeClaim(&_Stakemanager.TransactOpts)
}

// UpdateDynastyValue is a paid mutator transaction binding the contract method 0xe6692f49.
//
// Solidity: function updateDynastyValue(newDynasty uint256) returns()
func (_Stakemanager *StakemanagerTransactor) UpdateDynastyValue(opts *bind.TransactOpts, newDynasty *big.Int) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "updateDynastyValue", newDynasty)
}

// UpdateDynastyValue is a paid mutator transaction binding the contract method 0xe6692f49.
//
// Solidity: function updateDynastyValue(newDynasty uint256) returns()
func (_Stakemanager *StakemanagerSession) UpdateDynastyValue(newDynasty *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateDynastyValue(&_Stakemanager.TransactOpts, newDynasty)
}

// UpdateDynastyValue is a paid mutator transaction binding the contract method 0xe6692f49.
//
// Solidity: function updateDynastyValue(newDynasty uint256) returns()
func (_Stakemanager *StakemanagerTransactorSession) UpdateDynastyValue(newDynasty *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateDynastyValue(&_Stakemanager.TransactOpts, newDynasty)
}

// UpdateMinLockInPeriod is a paid mutator transaction binding the contract method 0x98ee773b.
//
// Solidity: function updateMinLockInPeriod(epochs uint256) returns()
func (_Stakemanager *StakemanagerTransactor) UpdateMinLockInPeriod(opts *bind.TransactOpts, epochs *big.Int) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "updateMinLockInPeriod", epochs)
}

// UpdateMinLockInPeriod is a paid mutator transaction binding the contract method 0x98ee773b.
//
// Solidity: function updateMinLockInPeriod(epochs uint256) returns()
func (_Stakemanager *StakemanagerSession) UpdateMinLockInPeriod(epochs *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateMinLockInPeriod(&_Stakemanager.TransactOpts, epochs)
}

// UpdateMinLockInPeriod is a paid mutator transaction binding the contract method 0x98ee773b.
//
// Solidity: function updateMinLockInPeriod(epochs uint256) returns()
func (_Stakemanager *StakemanagerTransactorSession) UpdateMinLockInPeriod(epochs *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateMinLockInPeriod(&_Stakemanager.TransactOpts, epochs)
}

// UpdateSigner is a paid mutator transaction binding the contract method 0xa7ecd37e.
//
// Solidity: function updateSigner(_signer address) returns()
func (_Stakemanager *StakemanagerTransactor) UpdateSigner(opts *bind.TransactOpts, _signer common.Address) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "updateSigner", _signer)
}

// UpdateSigner is a paid mutator transaction binding the contract method 0xa7ecd37e.
//
// Solidity: function updateSigner(_signer address) returns()
func (_Stakemanager *StakemanagerSession) UpdateSigner(_signer common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateSigner(&_Stakemanager.TransactOpts, _signer)
}

// UpdateSigner is a paid mutator transaction binding the contract method 0xa7ecd37e.
//
// Solidity: function updateSigner(_signer address) returns()
func (_Stakemanager *StakemanagerTransactorSession) UpdateSigner(_signer common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateSigner(&_Stakemanager.TransactOpts, _signer)
}

// UpdateValidatorThreshold is a paid mutator transaction binding the contract method 0x16827b1b.
//
// Solidity: function updateValidatorThreshold(newThreshold uint256) returns()
func (_Stakemanager *StakemanagerTransactor) UpdateValidatorThreshold(opts *bind.TransactOpts, newThreshold *big.Int) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "updateValidatorThreshold", newThreshold)
}

// UpdateValidatorThreshold is a paid mutator transaction binding the contract method 0x16827b1b.
//
// Solidity: function updateValidatorThreshold(newThreshold uint256) returns()
func (_Stakemanager *StakemanagerSession) UpdateValidatorThreshold(newThreshold *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateValidatorThreshold(&_Stakemanager.TransactOpts, newThreshold)
}

// UpdateValidatorThreshold is a paid mutator transaction binding the contract method 0x16827b1b.
//
// Solidity: function updateValidatorThreshold(newThreshold uint256) returns()
func (_Stakemanager *StakemanagerTransactorSession) UpdateValidatorThreshold(newThreshold *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateValidatorThreshold(&_Stakemanager.TransactOpts, newThreshold)
}

// StakemanagerDynastyValueChangeIterator is returned from FilterDynastyValueChange and is used to iterate over the raw logs and unpacked data for DynastyValueChange events raised by the Stakemanager contract.
type StakemanagerDynastyValueChangeIterator struct {
	Event *StakemanagerDynastyValueChange // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerDynastyValueChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerDynastyValueChange)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerDynastyValueChange)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerDynastyValueChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerDynastyValueChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerDynastyValueChange represents a DynastyValueChange event raised by the Stakemanager contract.
type StakemanagerDynastyValueChange struct {
	NewDynasty *big.Int
	OldDynasty *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDynastyValueChange is a free log retrieval operation binding the contract event 0x9444bfcfa6aed72a15da73de1220dcc07d7864119c44abfec0037bbcacefda98.
//
// Solidity: e DynastyValueChange(newDynasty uint256, oldDynasty uint256)
func (_Stakemanager *StakemanagerFilterer) FilterDynastyValueChange(opts *bind.FilterOpts) (*StakemanagerDynastyValueChangeIterator, error) {

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "DynastyValueChange")
	if err != nil {
		return nil, err
	}
	return &StakemanagerDynastyValueChangeIterator{contract: _Stakemanager.contract, event: "DynastyValueChange", logs: logs, sub: sub}, nil
}

// WatchDynastyValueChange is a free log subscription operation binding the contract event 0x9444bfcfa6aed72a15da73de1220dcc07d7864119c44abfec0037bbcacefda98.
//
// Solidity: e DynastyValueChange(newDynasty uint256, oldDynasty uint256)
func (_Stakemanager *StakemanagerFilterer) WatchDynastyValueChange(opts *bind.WatchOpts, sink chan<- *StakemanagerDynastyValueChange) (event.Subscription, error) {

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "DynastyValueChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerDynastyValueChange)
				if err := _Stakemanager.contract.UnpackLog(event, "DynastyValueChange", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// StakemanagerOwnershipRenouncedIterator is returned from FilterOwnershipRenounced and is used to iterate over the raw logs and unpacked data for OwnershipRenounced events raised by the Stakemanager contract.
type StakemanagerOwnershipRenouncedIterator struct {
	Event *StakemanagerOwnershipRenounced // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerOwnershipRenouncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerOwnershipRenounced)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerOwnershipRenounced)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerOwnershipRenouncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerOwnershipRenouncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerOwnershipRenounced represents a OwnershipRenounced event raised by the Stakemanager contract.
type StakemanagerOwnershipRenounced struct {
	PreviousOwner common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipRenounced is a free log retrieval operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Stakemanager *StakemanagerFilterer) FilterOwnershipRenounced(opts *bind.FilterOpts, previousOwner []common.Address) (*StakemanagerOwnershipRenouncedIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerOwnershipRenouncedIterator{contract: _Stakemanager.contract, event: "OwnershipRenounced", logs: logs, sub: sub}, nil
}

// WatchOwnershipRenounced is a free log subscription operation binding the contract event 0xf8df31144d9c2f0f6b59d69b8b98abd5459d07f2742c4df920b25aae33c64820.
//
// Solidity: e OwnershipRenounced(previousOwner indexed address)
func (_Stakemanager *StakemanagerFilterer) WatchOwnershipRenounced(opts *bind.WatchOpts, sink chan<- *StakemanagerOwnershipRenounced, previousOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "OwnershipRenounced", previousOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerOwnershipRenounced)
				if err := _Stakemanager.contract.UnpackLog(event, "OwnershipRenounced", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// StakemanagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Stakemanager contract.
type StakemanagerOwnershipTransferredIterator struct {
	Event *StakemanagerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerOwnershipTransferred represents a OwnershipTransferred event raised by the Stakemanager contract.
type StakemanagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Stakemanager *StakemanagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StakemanagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerOwnershipTransferredIterator{contract: _Stakemanager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Stakemanager *StakemanagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StakemanagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerOwnershipTransferred)
				if err := _Stakemanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// StakemanagerRootChainChangedIterator is returned from FilterRootChainChanged and is used to iterate over the raw logs and unpacked data for RootChainChanged events raised by the Stakemanager contract.
type StakemanagerRootChainChangedIterator struct {
	Event *StakemanagerRootChainChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerRootChainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerRootChainChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerRootChainChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerRootChainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerRootChainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerRootChainChanged represents a RootChainChanged event raised by the Stakemanager contract.
type StakemanagerRootChainChanged struct {
	PreviousRootChain common.Address
	NewRootChain      common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRootChainChanged is a free log retrieval operation binding the contract event 0x211c9015fc81c0dbd45bd99f0f29fc1c143bfd53442d5ffd722bbbef7a887fe9.
//
// Solidity: e RootChainChanged(previousRootChain indexed address, newRootChain indexed address)
func (_Stakemanager *StakemanagerFilterer) FilterRootChainChanged(opts *bind.FilterOpts, previousRootChain []common.Address, newRootChain []common.Address) (*StakemanagerRootChainChangedIterator, error) {

	var previousRootChainRule []interface{}
	for _, previousRootChainItem := range previousRootChain {
		previousRootChainRule = append(previousRootChainRule, previousRootChainItem)
	}
	var newRootChainRule []interface{}
	for _, newRootChainItem := range newRootChain {
		newRootChainRule = append(newRootChainRule, newRootChainItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "RootChainChanged", previousRootChainRule, newRootChainRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerRootChainChangedIterator{contract: _Stakemanager.contract, event: "RootChainChanged", logs: logs, sub: sub}, nil
}

// WatchRootChainChanged is a free log subscription operation binding the contract event 0x211c9015fc81c0dbd45bd99f0f29fc1c143bfd53442d5ffd722bbbef7a887fe9.
//
// Solidity: e RootChainChanged(previousRootChain indexed address, newRootChain indexed address)
func (_Stakemanager *StakemanagerFilterer) WatchRootChainChanged(opts *bind.WatchOpts, sink chan<- *StakemanagerRootChainChanged, previousRootChain []common.Address, newRootChain []common.Address) (event.Subscription, error) {

	var previousRootChainRule []interface{}
	for _, previousRootChainItem := range previousRootChain {
		previousRootChainRule = append(previousRootChainRule, previousRootChainItem)
	}
	var newRootChainRule []interface{}
	for _, newRootChainItem := range newRootChain {
		newRootChainRule = append(newRootChainRule, newRootChainItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "RootChainChanged", previousRootChainRule, newRootChainRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerRootChainChanged)
				if err := _Stakemanager.contract.UnpackLog(event, "RootChainChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// StakemanagerSignerChangeIterator is returned from FilterSignerChange and is used to iterate over the raw logs and unpacked data for SignerChange events raised by the Stakemanager contract.
type StakemanagerSignerChangeIterator struct {
	Event *StakemanagerSignerChange // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerSignerChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerSignerChange)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerSignerChange)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerSignerChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerSignerChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerSignerChange represents a SignerChange event raised by the Stakemanager contract.
type StakemanagerSignerChange struct {
	Validator common.Address
	NewSigner common.Address
	OldSigner common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterSignerChange is a free log retrieval operation binding the contract event 0xf63941bc0485c6b62a7130aa6d00b03821bd7de6dffb0a92d5eb66ef34fad2a5.
//
// Solidity: e SignerChange(validator indexed address, newSigner indexed address, oldSigner indexed address)
func (_Stakemanager *StakemanagerFilterer) FilterSignerChange(opts *bind.FilterOpts, validator []common.Address, newSigner []common.Address, oldSigner []common.Address) (*StakemanagerSignerChangeIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var newSignerRule []interface{}
	for _, newSignerItem := range newSigner {
		newSignerRule = append(newSignerRule, newSignerItem)
	}
	var oldSignerRule []interface{}
	for _, oldSignerItem := range oldSigner {
		oldSignerRule = append(oldSignerRule, oldSignerItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "SignerChange", validatorRule, newSignerRule, oldSignerRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerSignerChangeIterator{contract: _Stakemanager.contract, event: "SignerChange", logs: logs, sub: sub}, nil
}

// WatchSignerChange is a free log subscription operation binding the contract event 0xf63941bc0485c6b62a7130aa6d00b03821bd7de6dffb0a92d5eb66ef34fad2a5.
//
// Solidity: e SignerChange(validator indexed address, newSigner indexed address, oldSigner indexed address)
func (_Stakemanager *StakemanagerFilterer) WatchSignerChange(opts *bind.WatchOpts, sink chan<- *StakemanagerSignerChange, validator []common.Address, newSigner []common.Address, oldSigner []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var newSignerRule []interface{}
	for _, newSignerItem := range newSigner {
		newSignerRule = append(newSignerRule, newSignerItem)
	}
	var oldSignerRule []interface{}
	for _, oldSignerItem := range oldSigner {
		oldSignerRule = append(oldSignerRule, oldSignerItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "SignerChange", validatorRule, newSignerRule, oldSignerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerSignerChange)
				if err := _Stakemanager.contract.UnpackLog(event, "SignerChange", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// StakemanagerStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the Stakemanager contract.
type StakemanagerStakedIterator struct {
	Event *StakemanagerStaked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerStaked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerStaked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerStaked represents a Staked event raised by the Stakemanager contract.
type StakemanagerStaked struct {
	User           common.Address
	Signer         common.Address
	ActivatonEpoch *big.Int
	Amount         *big.Int
	Total          *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0xad3fa07f4195b47e64892eb944ecbfc253384053c119852bb2bcae484c2fcb69.
//
// Solidity: e Staked(user indexed address, signer indexed address, activatonEpoch indexed uint256, amount uint256, total uint256)
func (_Stakemanager *StakemanagerFilterer) FilterStaked(opts *bind.FilterOpts, user []common.Address, signer []common.Address, activatonEpoch []*big.Int) (*StakemanagerStakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}
	var activatonEpochRule []interface{}
	for _, activatonEpochItem := range activatonEpoch {
		activatonEpochRule = append(activatonEpochRule, activatonEpochItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "Staked", userRule, signerRule, activatonEpochRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerStakedIterator{contract: _Stakemanager.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0xad3fa07f4195b47e64892eb944ecbfc253384053c119852bb2bcae484c2fcb69.
//
// Solidity: e Staked(user indexed address, signer indexed address, activatonEpoch indexed uint256, amount uint256, total uint256)
func (_Stakemanager *StakemanagerFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *StakemanagerStaked, user []common.Address, signer []common.Address, activatonEpoch []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}
	var activatonEpochRule []interface{}
	for _, activatonEpochItem := range activatonEpoch {
		activatonEpochRule = append(activatonEpochRule, activatonEpochItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "Staked", userRule, signerRule, activatonEpochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerStaked)
				if err := _Stakemanager.contract.UnpackLog(event, "Staked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// StakemanagerThresholdChangeIterator is returned from FilterThresholdChange and is used to iterate over the raw logs and unpacked data for ThresholdChange events raised by the Stakemanager contract.
type StakemanagerThresholdChangeIterator struct {
	Event *StakemanagerThresholdChange // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerThresholdChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerThresholdChange)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerThresholdChange)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerThresholdChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerThresholdChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerThresholdChange represents a ThresholdChange event raised by the Stakemanager contract.
type StakemanagerThresholdChange struct {
	NewThreshold *big.Int
	OldThreshold *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterThresholdChange is a free log retrieval operation binding the contract event 0x5d16a900896e1160c2033bc940e6b072d3dc3b6a996fefb9b3b9b9678841824c.
//
// Solidity: e ThresholdChange(newThreshold uint256, oldThreshold uint256)
func (_Stakemanager *StakemanagerFilterer) FilterThresholdChange(opts *bind.FilterOpts) (*StakemanagerThresholdChangeIterator, error) {

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "ThresholdChange")
	if err != nil {
		return nil, err
	}
	return &StakemanagerThresholdChangeIterator{contract: _Stakemanager.contract, event: "ThresholdChange", logs: logs, sub: sub}, nil
}

// WatchThresholdChange is a free log subscription operation binding the contract event 0x5d16a900896e1160c2033bc940e6b072d3dc3b6a996fefb9b3b9b9678841824c.
//
// Solidity: e ThresholdChange(newThreshold uint256, oldThreshold uint256)
func (_Stakemanager *StakemanagerFilterer) WatchThresholdChange(opts *bind.WatchOpts, sink chan<- *StakemanagerThresholdChange) (event.Subscription, error) {

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "ThresholdChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerThresholdChange)
				if err := _Stakemanager.contract.UnpackLog(event, "ThresholdChange", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// StakemanagerUnstakeInitIterator is returned from FilterUnstakeInit and is used to iterate over the raw logs and unpacked data for UnstakeInit events raised by the Stakemanager contract.
type StakemanagerUnstakeInitIterator struct {
	Event *StakemanagerUnstakeInit // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerUnstakeInitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerUnstakeInit)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerUnstakeInit)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerUnstakeInitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerUnstakeInitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerUnstakeInit represents a UnstakeInit event raised by the Stakemanager contract.
type StakemanagerUnstakeInit struct {
	User              common.Address
	Amount            *big.Int
	DeactivationEpoch *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUnstakeInit is a free log retrieval operation binding the contract event 0xcde3813e379342e7506ca1f984065c81821e879aebd4be83cb35b5ab976518f9.
//
// Solidity: e UnstakeInit(user indexed address, amount indexed uint256, deactivationEpoch indexed uint256)
func (_Stakemanager *StakemanagerFilterer) FilterUnstakeInit(opts *bind.FilterOpts, user []common.Address, amount []*big.Int, deactivationEpoch []*big.Int) (*StakemanagerUnstakeInitIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var deactivationEpochRule []interface{}
	for _, deactivationEpochItem := range deactivationEpoch {
		deactivationEpochRule = append(deactivationEpochRule, deactivationEpochItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "UnstakeInit", userRule, amountRule, deactivationEpochRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerUnstakeInitIterator{contract: _Stakemanager.contract, event: "UnstakeInit", logs: logs, sub: sub}, nil
}

// WatchUnstakeInit is a free log subscription operation binding the contract event 0xcde3813e379342e7506ca1f984065c81821e879aebd4be83cb35b5ab976518f9.
//
// Solidity: e UnstakeInit(user indexed address, amount indexed uint256, deactivationEpoch indexed uint256)
func (_Stakemanager *StakemanagerFilterer) WatchUnstakeInit(opts *bind.WatchOpts, sink chan<- *StakemanagerUnstakeInit, user []common.Address, amount []*big.Int, deactivationEpoch []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var deactivationEpochRule []interface{}
	for _, deactivationEpochItem := range deactivationEpoch {
		deactivationEpochRule = append(deactivationEpochRule, deactivationEpochItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "UnstakeInit", userRule, amountRule, deactivationEpochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerUnstakeInit)
				if err := _Stakemanager.contract.UnpackLog(event, "UnstakeInit", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// StakemanagerUnstakedIterator is returned from FilterUnstaked and is used to iterate over the raw logs and unpacked data for Unstaked events raised by the Stakemanager contract.
type StakemanagerUnstakedIterator struct {
	Event *StakemanagerUnstaked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *StakemanagerUnstakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerUnstaked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(StakemanagerUnstaked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *StakemanagerUnstakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerUnstakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerUnstaked represents a Unstaked event raised by the Stakemanager contract.
type StakemanagerUnstaked struct {
	User   common.Address
	Amount *big.Int
	Total  *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterUnstaked is a free log retrieval operation binding the contract event 0x7fc4727e062e336010f2c282598ef5f14facb3de68cf8195c2f23e1454b2b74e.
//
// Solidity: e Unstaked(user indexed address, amount uint256, total uint256)
func (_Stakemanager *StakemanagerFilterer) FilterUnstaked(opts *bind.FilterOpts, user []common.Address) (*StakemanagerUnstakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "Unstaked", userRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerUnstakedIterator{contract: _Stakemanager.contract, event: "Unstaked", logs: logs, sub: sub}, nil
}

// WatchUnstaked is a free log subscription operation binding the contract event 0x7fc4727e062e336010f2c282598ef5f14facb3de68cf8195c2f23e1454b2b74e.
//
// Solidity: e Unstaked(user indexed address, amount uint256, total uint256)
func (_Stakemanager *StakemanagerFilterer) WatchUnstaked(opts *bind.WatchOpts, sink chan<- *StakemanagerUnstaked, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "Unstaked", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerUnstaked)
				if err := _Stakemanager.contract.UnpackLog(event, "Unstaked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}
