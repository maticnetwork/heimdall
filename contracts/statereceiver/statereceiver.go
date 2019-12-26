// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package statereceiver

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

// StatereceiverABI is the input ABI used to generate the binding from.
const StatereceiverABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"states\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"recordBytes\",\"type\":\"bytes\"}],\"name\":\"commitState\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getPendingStates\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"SYSTEM_ADDRESS\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"validatorSet\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"vote\",\"type\":\"bytes\"},{\"name\":\"sigs\",\"type\":\"bytes\"},{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"validateValidatorSet\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isValidatorSetContract\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"stateId\",\"type\":\"uint256\"}],\"name\":\"proposeState\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"isProducer\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"isValidator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Statereceiver is an auto generated Go binding around an Ethereum contract.
type Statereceiver struct {
	StatereceiverCaller     // Read-only binding to the contract
	StatereceiverTransactor // Write-only binding to the contract
	StatereceiverFilterer   // Log filterer for contract events
}

// StatereceiverCaller is an auto generated read-only Go binding around an Ethereum contract.
type StatereceiverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatereceiverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StatereceiverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatereceiverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StatereceiverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatereceiverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StatereceiverSession struct {
	Contract     *Statereceiver    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StatereceiverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StatereceiverCallerSession struct {
	Contract *StatereceiverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// StatereceiverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StatereceiverTransactorSession struct {
	Contract     *StatereceiverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// StatereceiverRaw is an auto generated low-level Go binding around an Ethereum contract.
type StatereceiverRaw struct {
	Contract *Statereceiver // Generic contract binding to access the raw methods on
}

// StatereceiverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StatereceiverCallerRaw struct {
	Contract *StatereceiverCaller // Generic read-only contract binding to access the raw methods on
}

// StatereceiverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StatereceiverTransactorRaw struct {
	Contract *StatereceiverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStatereceiver creates a new instance of Statereceiver, bound to a specific deployed contract.
func NewStatereceiver(address common.Address, backend bind.ContractBackend) (*Statereceiver, error) {
	contract, err := bindStatereceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Statereceiver{StatereceiverCaller: StatereceiverCaller{contract: contract}, StatereceiverTransactor: StatereceiverTransactor{contract: contract}, StatereceiverFilterer: StatereceiverFilterer{contract: contract}}, nil
}

// NewStatereceiverCaller creates a new read-only instance of Statereceiver, bound to a specific deployed contract.
func NewStatereceiverCaller(address common.Address, caller bind.ContractCaller) (*StatereceiverCaller, error) {
	contract, err := bindStatereceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StatereceiverCaller{contract: contract}, nil
}

// NewStatereceiverTransactor creates a new write-only instance of Statereceiver, bound to a specific deployed contract.
func NewStatereceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*StatereceiverTransactor, error) {
	contract, err := bindStatereceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StatereceiverTransactor{contract: contract}, nil
}

// NewStatereceiverFilterer creates a new log filterer instance of Statereceiver, bound to a specific deployed contract.
func NewStatereceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*StatereceiverFilterer, error) {
	contract, err := bindStatereceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StatereceiverFilterer{contract: contract}, nil
}

// bindStatereceiver binds a generic wrapper to an already deployed contract.
func bindStatereceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StatereceiverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Statereceiver *StatereceiverRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Statereceiver.Contract.StatereceiverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Statereceiver *StatereceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Statereceiver.Contract.StatereceiverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Statereceiver *StatereceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Statereceiver.Contract.StatereceiverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Statereceiver *StatereceiverCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Statereceiver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Statereceiver *StatereceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Statereceiver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Statereceiver *StatereceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Statereceiver.Contract.contract.Transact(opts, method, params...)
}

// SYSTEMADDRESS is a free data retrieval call binding the contract method 0x3434735f.
//
// Solidity: function SYSTEM_ADDRESS() constant returns(address)
func (_Statereceiver *StatereceiverCaller) SYSTEMADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Statereceiver.contract.Call(opts, out, "SYSTEM_ADDRESS")
	return *ret0, err
}

// SYSTEMADDRESS is a free data retrieval call binding the contract method 0x3434735f.
//
// Solidity: function SYSTEM_ADDRESS() constant returns(address)
func (_Statereceiver *StatereceiverSession) SYSTEMADDRESS() (common.Address, error) {
	return _Statereceiver.Contract.SYSTEMADDRESS(&_Statereceiver.CallOpts)
}

// SYSTEMADDRESS is a free data retrieval call binding the contract method 0x3434735f.
//
// Solidity: function SYSTEM_ADDRESS() constant returns(address)
func (_Statereceiver *StatereceiverCallerSession) SYSTEMADDRESS() (common.Address, error) {
	return _Statereceiver.Contract.SYSTEMADDRESS(&_Statereceiver.CallOpts)
}

// GetPendingStates is a free data retrieval call binding the contract method 0x21ec23b6.
//
// Solidity: function getPendingStates() constant returns(uint256[])
func (_Statereceiver *StatereceiverCaller) GetPendingStates(opts *bind.CallOpts) ([]*big.Int, error) {
	var (
		ret0 = new([]*big.Int)
	)
	out := ret0
	err := _Statereceiver.contract.Call(opts, out, "getPendingStates")
	return *ret0, err
}

// GetPendingStates is a free data retrieval call binding the contract method 0x21ec23b6.
//
// Solidity: function getPendingStates() constant returns(uint256[])
func (_Statereceiver *StatereceiverSession) GetPendingStates() ([]*big.Int, error) {
	return _Statereceiver.Contract.GetPendingStates(&_Statereceiver.CallOpts)
}

// GetPendingStates is a free data retrieval call binding the contract method 0x21ec23b6.
//
// Solidity: function getPendingStates() constant returns(uint256[])
func (_Statereceiver *StatereceiverCallerSession) GetPendingStates() ([]*big.Int, error) {
	return _Statereceiver.Contract.GetPendingStates(&_Statereceiver.CallOpts)
}

// IsProducer is a free data retrieval call binding the contract method 0xf5521022.
//
// Solidity: function isProducer(address signer) constant returns(bool)
func (_Statereceiver *StatereceiverCaller) IsProducer(opts *bind.CallOpts, signer common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Statereceiver.contract.Call(opts, out, "isProducer", signer)
	return *ret0, err
}

// IsProducer is a free data retrieval call binding the contract method 0xf5521022.
//
// Solidity: function isProducer(address signer) constant returns(bool)
func (_Statereceiver *StatereceiverSession) IsProducer(signer common.Address) (bool, error) {
	return _Statereceiver.Contract.IsProducer(&_Statereceiver.CallOpts, signer)
}

// IsProducer is a free data retrieval call binding the contract method 0xf5521022.
//
// Solidity: function isProducer(address signer) constant returns(bool)
func (_Statereceiver *StatereceiverCallerSession) IsProducer(signer common.Address) (bool, error) {
	return _Statereceiver.Contract.IsProducer(&_Statereceiver.CallOpts, signer)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address signer) constant returns(bool)
func (_Statereceiver *StatereceiverCaller) IsValidator(opts *bind.CallOpts, signer common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Statereceiver.contract.Call(opts, out, "isValidator", signer)
	return *ret0, err
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address signer) constant returns(bool)
func (_Statereceiver *StatereceiverSession) IsValidator(signer common.Address) (bool, error) {
	return _Statereceiver.Contract.IsValidator(&_Statereceiver.CallOpts, signer)
}

// IsValidator is a free data retrieval call binding the contract method 0xfacd743b.
//
// Solidity: function isValidator(address signer) constant returns(bool)
func (_Statereceiver *StatereceiverCallerSession) IsValidator(signer common.Address) (bool, error) {
	return _Statereceiver.Contract.IsValidator(&_Statereceiver.CallOpts, signer)
}

// IsValidatorSetContract is a free data retrieval call binding the contract method 0xd79e60b7.
//
// Solidity: function isValidatorSetContract() constant returns(bool)
func (_Statereceiver *StatereceiverCaller) IsValidatorSetContract(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Statereceiver.contract.Call(opts, out, "isValidatorSetContract")
	return *ret0, err
}

// IsValidatorSetContract is a free data retrieval call binding the contract method 0xd79e60b7.
//
// Solidity: function isValidatorSetContract() constant returns(bool)
func (_Statereceiver *StatereceiverSession) IsValidatorSetContract() (bool, error) {
	return _Statereceiver.Contract.IsValidatorSetContract(&_Statereceiver.CallOpts)
}

// IsValidatorSetContract is a free data retrieval call binding the contract method 0xd79e60b7.
//
// Solidity: function isValidatorSetContract() constant returns(bool)
func (_Statereceiver *StatereceiverCallerSession) IsValidatorSetContract() (bool, error) {
	return _Statereceiver.Contract.IsValidatorSetContract(&_Statereceiver.CallOpts)
}

// States is a free data retrieval call binding the contract method 0x017a9105.
//
// Solidity: function states(uint256 ) constant returns(bool)
func (_Statereceiver *StatereceiverCaller) States(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Statereceiver.contract.Call(opts, out, "states", arg0)
	return *ret0, err
}

// States is a free data retrieval call binding the contract method 0x017a9105.
//
// Solidity: function states(uint256 ) constant returns(bool)
func (_Statereceiver *StatereceiverSession) States(arg0 *big.Int) (bool, error) {
	return _Statereceiver.Contract.States(&_Statereceiver.CallOpts, arg0)
}

// States is a free data retrieval call binding the contract method 0x017a9105.
//
// Solidity: function states(uint256 ) constant returns(bool)
func (_Statereceiver *StatereceiverCallerSession) States(arg0 *big.Int) (bool, error) {
	return _Statereceiver.Contract.States(&_Statereceiver.CallOpts, arg0)
}

// ValidatorSet is a free data retrieval call binding the contract method 0x9426e226.
//
// Solidity: function validatorSet() constant returns(address)
func (_Statereceiver *StatereceiverCaller) ValidatorSet(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Statereceiver.contract.Call(opts, out, "validatorSet")
	return *ret0, err
}

// ValidatorSet is a free data retrieval call binding the contract method 0x9426e226.
//
// Solidity: function validatorSet() constant returns(address)
func (_Statereceiver *StatereceiverSession) ValidatorSet() (common.Address, error) {
	return _Statereceiver.Contract.ValidatorSet(&_Statereceiver.CallOpts)
}

// ValidatorSet is a free data retrieval call binding the contract method 0x9426e226.
//
// Solidity: function validatorSet() constant returns(address)
func (_Statereceiver *StatereceiverCallerSession) ValidatorSet() (common.Address, error) {
	return _Statereceiver.Contract.ValidatorSet(&_Statereceiver.CallOpts)
}

// CommitState is a paid mutator transaction binding the contract method 0x080356b7.
//
// Solidity: function commitState(bytes recordBytes) returns()
func (_Statereceiver *StatereceiverTransactor) CommitState(opts *bind.TransactOpts, recordBytes []byte) (*types.Transaction, error) {
	return _Statereceiver.contract.Transact(opts, "commitState", recordBytes)
}

// CommitState is a paid mutator transaction binding the contract method 0x080356b7.
//
// Solidity: function commitState(bytes recordBytes) returns()
func (_Statereceiver *StatereceiverSession) CommitState(recordBytes []byte) (*types.Transaction, error) {
	return _Statereceiver.Contract.CommitState(&_Statereceiver.TransactOpts, recordBytes)
}

// CommitState is a paid mutator transaction binding the contract method 0x080356b7.
//
// Solidity: function commitState(bytes recordBytes) returns()
func (_Statereceiver *StatereceiverTransactorSession) CommitState(recordBytes []byte) (*types.Transaction, error) {
	return _Statereceiver.Contract.CommitState(&_Statereceiver.TransactOpts, recordBytes)
}

// ProposeState is a paid mutator transaction binding the contract method 0xede01f17.
//
// Solidity: function proposeState(uint256 stateId) returns()
func (_Statereceiver *StatereceiverTransactor) ProposeState(opts *bind.TransactOpts, stateId *big.Int) (*types.Transaction, error) {
	return _Statereceiver.contract.Transact(opts, "proposeState", stateId)
}

// ProposeState is a paid mutator transaction binding the contract method 0xede01f17.
//
// Solidity: function proposeState(uint256 stateId) returns()
func (_Statereceiver *StatereceiverSession) ProposeState(stateId *big.Int) (*types.Transaction, error) {
	return _Statereceiver.Contract.ProposeState(&_Statereceiver.TransactOpts, stateId)
}

// ProposeState is a paid mutator transaction binding the contract method 0xede01f17.
//
// Solidity: function proposeState(uint256 stateId) returns()
func (_Statereceiver *StatereceiverTransactorSession) ProposeState(stateId *big.Int) (*types.Transaction, error) {
	return _Statereceiver.Contract.ProposeState(&_Statereceiver.TransactOpts, stateId)
}

// ValidateValidatorSet is a paid mutator transaction binding the contract method 0xd0504f89.
//
// Solidity: function validateValidatorSet(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statereceiver *StatereceiverTransactor) ValidateValidatorSet(opts *bind.TransactOpts, vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statereceiver.contract.Transact(opts, "validateValidatorSet", vote, sigs, txBytes, proof)
}

// ValidateValidatorSet is a paid mutator transaction binding the contract method 0xd0504f89.
//
// Solidity: function validateValidatorSet(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statereceiver *StatereceiverSession) ValidateValidatorSet(vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statereceiver.Contract.ValidateValidatorSet(&_Statereceiver.TransactOpts, vote, sigs, txBytes, proof)
}

// ValidateValidatorSet is a paid mutator transaction binding the contract method 0xd0504f89.
//
// Solidity: function validateValidatorSet(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statereceiver *StatereceiverTransactorSession) ValidateValidatorSet(vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statereceiver.Contract.ValidateValidatorSet(&_Statereceiver.TransactOpts, vote, sigs, txBytes, proof)
}
