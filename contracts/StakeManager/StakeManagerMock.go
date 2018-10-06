// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package StakeManager

import (
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// ContractsABI is the input ABI used to generate the binding from.
const ContractsABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_pubKey\",\"type\":\"string\"},{\"name\":\"_power\",\"type\":\"int256\"}],\"name\":\"addValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lastValidatorIndex\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validators\",\"outputs\":[{\"name\":\"pubkey\",\"type\":\"string\"},{\"name\":\"power\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Contracts is an auto generated Go binding around an Ethereum contract.
type Contracts struct {
	ContractsCaller     // Read-only binding to the contract
	ContractsTransactor // Write-only binding to the contract
	ContractsFilterer   // Log filterer for contract events
}

// ContractsCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractsSession struct {
	Contract     *Contracts        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractsCallerSession struct {
	Contract *ContractsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ContractsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractsTransactorSession struct {
	Contract     *ContractsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ContractsRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractsRaw struct {
	Contract *Contracts // Generic contract binding to access the raw methods on
}

// ContractsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractsCallerRaw struct {
	Contract *ContractsCaller // Generic read-only contract binding to access the raw methods on
}

// ContractsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractsTransactorRaw struct {
	Contract *ContractsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContracts creates a new instance of Contracts, bound to a specific deployed contract.
func NewContracts(address common.Address, backend bind.ContractBackend) (*Contracts, error) {
	contract, err := bindContracts(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contracts{ContractsCaller: ContractsCaller{contract: contract}, ContractsTransactor: ContractsTransactor{contract: contract}, ContractsFilterer: ContractsFilterer{contract: contract}}, nil
}

// NewContractsCaller creates a new read-only instance of Contracts, bound to a specific deployed contract.
func NewContractsCaller(address common.Address, caller bind.ContractCaller) (*ContractsCaller, error) {
	contract, err := bindContracts(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractsCaller{contract: contract}, nil
}

// NewContractsTransactor creates a new write-only instance of Contracts, bound to a specific deployed contract.
func NewContractsTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractsTransactor, error) {
	contract, err := bindContracts(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractsTransactor{contract: contract}, nil
}

// NewContractsFilterer creates a new log filterer instance of Contracts, bound to a specific deployed contract.
func NewContractsFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractsFilterer, error) {
	contract, err := bindContracts(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractsFilterer{contract: contract}, nil
}

// bindContracts binds a generic wrapper to an already deployed contract.
func bindContracts(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contracts *ContractsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contracts.Contract.ContractsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contracts *ContractsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contracts.Contract.ContractsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contracts *ContractsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contracts.Contract.ContractsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contracts *ContractsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contracts.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contracts *ContractsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contracts.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contracts *ContractsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contracts.Contract.contract.Transact(opts, method, params...)
}

// LastValidatorIndex is a free data retrieval call binding the contract method 0xe0d3e8a5.
//
// Solidity: function lastValidatorIndex() constant returns(uint256)
func (_Contracts *ContractsCaller) LastValidatorIndex(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Contracts.contract.Call(opts, out, "lastValidatorIndex")
	return *ret0, err
}

// LastValidatorIndex is a free data retrieval call binding the contract method 0xe0d3e8a5.
//
// Solidity: function lastValidatorIndex() constant returns(uint256)
func (_Contracts *ContractsSession) LastValidatorIndex() (*big.Int, error) {
	return _Contracts.Contract.LastValidatorIndex(&_Contracts.CallOpts)
}

// LastValidatorIndex is a free data retrieval call binding the contract method 0xe0d3e8a5.
//
// Solidity: function lastValidatorIndex() constant returns(uint256)
func (_Contracts *ContractsCallerSession) LastValidatorIndex() (*big.Int, error) {
	return _Contracts.Contract.LastValidatorIndex(&_Contracts.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(pubkey string, power int256)
func (_Contracts *ContractsCaller) Validators(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Pubkey string
	Power  *big.Int
}, error) {
	ret := new(struct {
		Pubkey string
		Power  *big.Int
	})
	out := ret
	err := _Contracts.contract.Call(opts, out, "validators", arg0)
	return *ret, err
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(pubkey string, power int256)
func (_Contracts *ContractsSession) Validators(arg0 *big.Int) (struct {
	Pubkey string
	Power  *big.Int
}, error) {
	return _Contracts.Contract.Validators(&_Contracts.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(pubkey string, power int256)
func (_Contracts *ContractsCallerSession) Validators(arg0 *big.Int) (struct {
	Pubkey string
	Power  *big.Int
}, error) {
	return _Contracts.Contract.Validators(&_Contracts.CallOpts, arg0)
}

// AddValidator is a paid mutator transaction binding the contract method 0x57f0eb72.
//
// Solidity: function addValidator(_pubKey string, _power int256) returns()
func (_Contracts *ContractsTransactor) AddValidator(opts *bind.TransactOpts, _pubKey string, _power *big.Int) (*types.Transaction, error) {
	return _Contracts.contract.Transact(opts, "addValidator", _pubKey, _power)
}

// AddValidator is a paid mutator transaction binding the contract method 0x57f0eb72.
//
// Solidity: function addValidator(_pubKey string, _power int256) returns()
func (_Contracts *ContractsSession) AddValidator(_pubKey string, _power *big.Int) (*types.Transaction, error) {
	return _Contracts.Contract.AddValidator(&_Contracts.TransactOpts, _pubKey, _power)
}

// AddValidator is a paid mutator transaction binding the contract method 0x57f0eb72.
//
// Solidity: function addValidator(_pubKey string, _power int256) returns()
func (_Contracts *ContractsTransactorSession) AddValidator(_pubKey string, _power *big.Int) (*types.Transaction, error) {
	return _Contracts.Contract.AddValidator(&_Contracts.TransactOpts, _pubKey, _power)
}
