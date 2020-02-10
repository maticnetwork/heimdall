// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stakemanager

import (
	"math/big"
	"strings"

	ethereum "github.com/maticnetwork/bor"
	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/bor/accounts/abi/bind"
	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/bor/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = abi.U256
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// StakemanagerABI is the input ABI used to generate the binding from.
const StakemanagerABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"WITHDRAWAL_DELAY\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"accountStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"stakePower\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"blockInterval\",\"type\":\"uint256\"},{\"internalType\":\"bytes32\",\"name\":\"voteHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"stateRoot\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"sigs\",\"type\":\"bytes\"}],\"name\":\"checkSignatures\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"heimdallFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"needDelegation\",\"type\":\"bool\"}],\"name\":\"confirmAuctionBid\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"delegator\",\"type\":\"address\"}],\"name\":\"delegationTransfer\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"heimdallFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"needDelegation\",\"type\":\"bool\"}],\"name\":\"stake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"heimdallFee\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"needDelegation\",\"type\":\"bool\"}],\"name\":\"stakeFor\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"startAuction\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"supportsHistory\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"totalStakedFor\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"internalType\":\"int256\",\"name\":\"amount\",\"type\":\"int256\"}],\"name\":\"updateValidatorState\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

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

// AccountStateRoot is a free data retrieval call binding the contract method 0x17c2b910.
//
// Solidity: function accountStateRoot() constant returns(bytes32)
func (_Stakemanager *StakemanagerCaller) AccountStateRoot(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "accountStateRoot")
	return *ret0, err
}

// AccountStateRoot is a free data retrieval call binding the contract method 0x17c2b910.
//
// Solidity: function accountStateRoot() constant returns(bytes32)
func (_Stakemanager *StakemanagerSession) AccountStateRoot() ([32]byte, error) {
	return _Stakemanager.Contract.AccountStateRoot(&_Stakemanager.CallOpts)
}

// AccountStateRoot is a free data retrieval call binding the contract method 0x17c2b910.
//
// Solidity: function accountStateRoot() constant returns(bytes32)
func (_Stakemanager *StakemanagerCallerSession) AccountStateRoot() ([32]byte, error) {
	return _Stakemanager.Contract.AccountStateRoot(&_Stakemanager.CallOpts)
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

// TotalStakedFor is a free data retrieval call binding the contract method 0x4b341aed.
//
// Solidity: function totalStakedFor(address addr) constant returns(uint256)
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
// Solidity: function totalStakedFor(address addr) constant returns(uint256)
func (_Stakemanager *StakemanagerSession) TotalStakedFor(addr common.Address) (*big.Int, error) {
	return _Stakemanager.Contract.TotalStakedFor(&_Stakemanager.CallOpts, addr)
}

// TotalStakedFor is a free data retrieval call binding the contract method 0x4b341aed.
//
// Solidity: function totalStakedFor(address addr) constant returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) TotalStakedFor(addr common.Address) (*big.Int, error) {
	return _Stakemanager.Contract.TotalStakedFor(&_Stakemanager.CallOpts, addr)
}

// CheckSignatures is a paid mutator transaction binding the contract method 0x5eca276b.
//
// Solidity: function checkSignatures(uint256 stakePower, uint256 blockInterval, bytes32 voteHash, bytes32 stateRoot, bytes sigs) returns(uint256)
func (_Stakemanager *StakemanagerTransactor) CheckSignatures(opts *bind.TransactOpts, stakePower *big.Int, blockInterval *big.Int, voteHash [32]byte, stateRoot [32]byte, sigs []byte) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "checkSignatures", stakePower, blockInterval, voteHash, stateRoot, sigs)
}

// CheckSignatures is a paid mutator transaction binding the contract method 0x5eca276b.
//
// Solidity: function checkSignatures(uint256 stakePower, uint256 blockInterval, bytes32 voteHash, bytes32 stateRoot, bytes sigs) returns(uint256)
func (_Stakemanager *StakemanagerSession) CheckSignatures(stakePower *big.Int, blockInterval *big.Int, voteHash [32]byte, stateRoot [32]byte, sigs []byte) (*types.Transaction, error) {
	return _Stakemanager.Contract.CheckSignatures(&_Stakemanager.TransactOpts, stakePower, blockInterval, voteHash, stateRoot, sigs)
}

// CheckSignatures is a paid mutator transaction binding the contract method 0x5eca276b.
//
// Solidity: function checkSignatures(uint256 stakePower, uint256 blockInterval, bytes32 voteHash, bytes32 stateRoot, bytes sigs) returns(uint256)
func (_Stakemanager *StakemanagerTransactorSession) CheckSignatures(stakePower *big.Int, blockInterval *big.Int, voteHash [32]byte, stateRoot [32]byte, sigs []byte) (*types.Transaction, error) {
	return _Stakemanager.Contract.CheckSignatures(&_Stakemanager.TransactOpts, stakePower, blockInterval, voteHash, stateRoot, sigs)
}

// ConfirmAuctionBid is a paid mutator transaction binding the contract method 0x7d4bd608.
//
// Solidity: function confirmAuctionBid(uint256 validatorId, uint256 heimdallFee, address signer, bool needDelegation) returns()
func (_Stakemanager *StakemanagerTransactor) ConfirmAuctionBid(opts *bind.TransactOpts, validatorId *big.Int, heimdallFee *big.Int, signer common.Address, needDelegation bool) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "confirmAuctionBid", validatorId, heimdallFee, signer, needDelegation)
}

// ConfirmAuctionBid is a paid mutator transaction binding the contract method 0x7d4bd608.
//
// Solidity: function confirmAuctionBid(uint256 validatorId, uint256 heimdallFee, address signer, bool needDelegation) returns()
func (_Stakemanager *StakemanagerSession) ConfirmAuctionBid(validatorId *big.Int, heimdallFee *big.Int, signer common.Address, needDelegation bool) (*types.Transaction, error) {
	return _Stakemanager.Contract.ConfirmAuctionBid(&_Stakemanager.TransactOpts, validatorId, heimdallFee, signer, needDelegation)
}

// ConfirmAuctionBid is a paid mutator transaction binding the contract method 0x7d4bd608.
//
// Solidity: function confirmAuctionBid(uint256 validatorId, uint256 heimdallFee, address signer, bool needDelegation) returns()
func (_Stakemanager *StakemanagerTransactorSession) ConfirmAuctionBid(validatorId *big.Int, heimdallFee *big.Int, signer common.Address, needDelegation bool) (*types.Transaction, error) {
	return _Stakemanager.Contract.ConfirmAuctionBid(&_Stakemanager.TransactOpts, validatorId, heimdallFee, signer, needDelegation)
}

// DelegationTransfer is a paid mutator transaction binding the contract method 0x4777fb42.
//
// Solidity: function delegationTransfer(uint256 validatorId, uint256 amount, address delegator) returns()
func (_Stakemanager *StakemanagerTransactor) DelegationTransfer(opts *bind.TransactOpts, validatorId *big.Int, amount *big.Int, delegator common.Address) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "delegationTransfer", validatorId, amount, delegator)
}

// DelegationTransfer is a paid mutator transaction binding the contract method 0x4777fb42.
//
// Solidity: function delegationTransfer(uint256 validatorId, uint256 amount, address delegator) returns()
func (_Stakemanager *StakemanagerSession) DelegationTransfer(validatorId *big.Int, amount *big.Int, delegator common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.DelegationTransfer(&_Stakemanager.TransactOpts, validatorId, amount, delegator)
}

// DelegationTransfer is a paid mutator transaction binding the contract method 0x4777fb42.
//
// Solidity: function delegationTransfer(uint256 validatorId, uint256 amount, address delegator) returns()
func (_Stakemanager *StakemanagerTransactorSession) DelegationTransfer(validatorId *big.Int, amount *big.Int, delegator common.Address) (*types.Transaction, error) {
	return _Stakemanager.Contract.DelegationTransfer(&_Stakemanager.TransactOpts, validatorId, amount, delegator)
}

// Stake is a paid mutator transaction binding the contract method 0x06254a9c.
//
// Solidity: function stake(uint256 amount, uint256 heimdallFee, address signer, bool needDelegation) returns()
func (_Stakemanager *StakemanagerTransactor) Stake(opts *bind.TransactOpts, amount *big.Int, heimdallFee *big.Int, signer common.Address, needDelegation bool) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "stake", amount, heimdallFee, signer, needDelegation)
}

// Stake is a paid mutator transaction binding the contract method 0x06254a9c.
//
// Solidity: function stake(uint256 amount, uint256 heimdallFee, address signer, bool needDelegation) returns()
func (_Stakemanager *StakemanagerSession) Stake(amount *big.Int, heimdallFee *big.Int, signer common.Address, needDelegation bool) (*types.Transaction, error) {
	return _Stakemanager.Contract.Stake(&_Stakemanager.TransactOpts, amount, heimdallFee, signer, needDelegation)
}

// Stake is a paid mutator transaction binding the contract method 0x06254a9c.
//
// Solidity: function stake(uint256 amount, uint256 heimdallFee, address signer, bool needDelegation) returns()
func (_Stakemanager *StakemanagerTransactorSession) Stake(amount *big.Int, heimdallFee *big.Int, signer common.Address, needDelegation bool) (*types.Transaction, error) {
	return _Stakemanager.Contract.Stake(&_Stakemanager.TransactOpts, amount, heimdallFee, signer, needDelegation)
}

// StakeFor is a paid mutator transaction binding the contract method 0x76324ea5.
//
// Solidity: function stakeFor(address user, uint256 amount, uint256 heimdallFee, address signer, bool needDelegation) returns()
func (_Stakemanager *StakemanagerTransactor) StakeFor(opts *bind.TransactOpts, user common.Address, amount *big.Int, heimdallFee *big.Int, signer common.Address, needDelegation bool) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "stakeFor", user, amount, heimdallFee, signer, needDelegation)
}

// StakeFor is a paid mutator transaction binding the contract method 0x76324ea5.
//
// Solidity: function stakeFor(address user, uint256 amount, uint256 heimdallFee, address signer, bool needDelegation) returns()
func (_Stakemanager *StakemanagerSession) StakeFor(user common.Address, amount *big.Int, heimdallFee *big.Int, signer common.Address, needDelegation bool) (*types.Transaction, error) {
	return _Stakemanager.Contract.StakeFor(&_Stakemanager.TransactOpts, user, amount, heimdallFee, signer, needDelegation)
}

// StakeFor is a paid mutator transaction binding the contract method 0x76324ea5.
//
// Solidity: function stakeFor(address user, uint256 amount, uint256 heimdallFee, address signer, bool needDelegation) returns()
func (_Stakemanager *StakemanagerTransactorSession) StakeFor(user common.Address, amount *big.Int, heimdallFee *big.Int, signer common.Address, needDelegation bool) (*types.Transaction, error) {
	return _Stakemanager.Contract.StakeFor(&_Stakemanager.TransactOpts, user, amount, heimdallFee, signer, needDelegation)
}

// StartAuction is a paid mutator transaction binding the contract method 0x4fee13fc.
//
// Solidity: function startAuction(uint256 validatorId, uint256 amount) returns()
func (_Stakemanager *StakemanagerTransactor) StartAuction(opts *bind.TransactOpts, validatorId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "startAuction", validatorId, amount)
}

// StartAuction is a paid mutator transaction binding the contract method 0x4fee13fc.
//
// Solidity: function startAuction(uint256 validatorId, uint256 amount) returns()
func (_Stakemanager *StakemanagerSession) StartAuction(validatorId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.StartAuction(&_Stakemanager.TransactOpts, validatorId, amount)
}

// StartAuction is a paid mutator transaction binding the contract method 0x4fee13fc.
//
// Solidity: function startAuction(uint256 validatorId, uint256 amount) returns()
func (_Stakemanager *StakemanagerTransactorSession) StartAuction(validatorId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.StartAuction(&_Stakemanager.TransactOpts, validatorId, amount)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 validatorId) returns()
func (_Stakemanager *StakemanagerTransactor) Unstake(opts *bind.TransactOpts, validatorId *big.Int) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "unstake", validatorId)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 validatorId) returns()
func (_Stakemanager *StakemanagerSession) Unstake(validatorId *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.Unstake(&_Stakemanager.TransactOpts, validatorId)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 validatorId) returns()
func (_Stakemanager *StakemanagerTransactorSession) Unstake(validatorId *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.Unstake(&_Stakemanager.TransactOpts, validatorId)
}

// UpdateValidatorState is a paid mutator transaction binding the contract method 0x9ff11500.
//
// Solidity: function updateValidatorState(uint256 validatorId, int256 amount) returns()
func (_Stakemanager *StakemanagerTransactor) UpdateValidatorState(opts *bind.TransactOpts, validatorId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "updateValidatorState", validatorId, amount)
}

// UpdateValidatorState is a paid mutator transaction binding the contract method 0x9ff11500.
//
// Solidity: function updateValidatorState(uint256 validatorId, int256 amount) returns()
func (_Stakemanager *StakemanagerSession) UpdateValidatorState(validatorId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateValidatorState(&_Stakemanager.TransactOpts, validatorId, amount)
}

// UpdateValidatorState is a paid mutator transaction binding the contract method 0x9ff11500.
//
// Solidity: function updateValidatorState(uint256 validatorId, int256 amount) returns()
func (_Stakemanager *StakemanagerTransactorSession) UpdateValidatorState(validatorId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.UpdateValidatorState(&_Stakemanager.TransactOpts, validatorId, amount)
}
