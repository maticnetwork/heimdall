// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rootchain

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

// RootchainABI is the input ABI used to generate the binding from.
const RootchainABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"childChainContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"depositCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"reverseTokens\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"headerBlocks\",\"outputs\":[{\"name\":\"root\",\"type\":\"bytes32\"},{\"name\":\"start\",\"type\":\"uint256\"},{\"name\":\"end\",\"type\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"wethToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"stakeManager\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"withdrawSignature\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes4\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"networkId\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"validatorContracts\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"deposits\",\"outputs\":[{\"name\":\"header\",\"type\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chain\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"withdraws\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"withdrawEventSignature\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokens\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentHeaderBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"_stakeManager\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousChildChain\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newChildChain\",\"type\":\"address\"}],\"name\":\"ChildChainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"rootToken\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"childToken\",\"type\":\"address\"}],\"name\":\"TokenMapped\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"}],\"name\":\"ValidatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"}],\"name\":\"ValidatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"depositCount\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Withdraw\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"proposer\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"number\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"start\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"end\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"NewHeaderBlock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"newChildChain\",\"type\":\"address\"}],\"name\":\"setChildContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"rootToken\",\"type\":\"address\"},{\"name\":\"childToken\",\"type\":\"address\"}],\"name\":\"mapToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"setWETHToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"addValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"removeValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_stakeManager\",\"type\":\"address\"}],\"name\":\"setStakeManager\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\"},{\"name\":\"start\",\"type\":\"uint256\"},{\"name\":\"end\",\"type\":\"uint256\"},{\"name\":\"sigs\",\"type\":\"bytes\"}],\"name\":\"submitHeaderBlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentChildBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"headerNumber\",\"type\":\"uint256\"}],\"name\":\"getHeaderBlock\",\"outputs\":[{\"name\":\"root\",\"type\":\"bytes32\"},{\"name\":\"start\",\"type\":\"uint256\"},{\"name\":\"end\",\"type\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"depositCount\",\"type\":\"uint256\"}],\"name\":\"getDepositBlock\",\"outputs\":[{\"name\":\"header\",\"type\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_sender\",\"type\":\"address\"},{\"name\":\"_value\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"bytes\"}],\"name\":\"tokenFallback\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"user\",\"type\":\"address\"}],\"name\":\"depositEthers\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"depositEthers\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"user\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"headerNumber\",\"type\":\"uint256\"},{\"name\":\"headerProof\",\"type\":\"bytes\"},{\"name\":\"blockNumber\",\"type\":\"uint256\"},{\"name\":\"blockTime\",\"type\":\"uint256\"},{\"name\":\"txRoot\",\"type\":\"bytes32\"},{\"name\":\"receiptRoot\",\"type\":\"bytes32\"},{\"name\":\"path\",\"type\":\"bytes\"},{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"txProof\",\"type\":\"bytes\"},{\"name\":\"receiptBytes\",\"type\":\"bytes\"},{\"name\":\"receiptProof\",\"type\":\"bytes\"}],\"name\":\"withdraw\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"slash\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"

// Rootchain is an auto generated Go binding around an Ethereum contract.
type Rootchain struct {
	RootchainCaller     // Read-only binding to the contract
	RootchainTransactor // Write-only binding to the contract
	RootchainFilterer   // Log filterer for contract events
}

// RootchainCaller is an auto generated read-only Go binding around an Ethereum contract.
type RootchainCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootchainTransactor is an auto generated write-only Go binding around an Ethereum contract.
type RootchainTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootchainFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type RootchainFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// RootchainSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type RootchainSession struct {
	Contract     *Rootchain        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// RootchainCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type RootchainCallerSession struct {
	Contract *RootchainCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// RootchainTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type RootchainTransactorSession struct {
	Contract     *RootchainTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// RootchainRaw is an auto generated low-level Go binding around an Ethereum contract.
type RootchainRaw struct {
	Contract *Rootchain // Generic contract binding to access the raw methods on
}

// RootchainCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type RootchainCallerRaw struct {
	Contract *RootchainCaller // Generic read-only contract binding to access the raw methods on
}

// RootchainTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type RootchainTransactorRaw struct {
	Contract *RootchainTransactor // Generic write-only contract binding to access the raw methods on
}

// NewRootchain creates a new instance of Rootchain, bound to a specific deployed contract.
func NewRootchain(address common.Address, backend bind.ContractBackend) (*Rootchain, error) {
	contract, err := bindRootchain(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Rootchain{RootchainCaller: RootchainCaller{contract: contract}, RootchainTransactor: RootchainTransactor{contract: contract}, RootchainFilterer: RootchainFilterer{contract: contract}}, nil
}

// NewRootchainCaller creates a new read-only instance of Rootchain, bound to a specific deployed contract.
func NewRootchainCaller(address common.Address, caller bind.ContractCaller) (*RootchainCaller, error) {
	contract, err := bindRootchain(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &RootchainCaller{contract: contract}, nil
}

// NewRootchainTransactor creates a new write-only instance of Rootchain, bound to a specific deployed contract.
func NewRootchainTransactor(address common.Address, transactor bind.ContractTransactor) (*RootchainTransactor, error) {
	contract, err := bindRootchain(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &RootchainTransactor{contract: contract}, nil
}

// NewRootchainFilterer creates a new log filterer instance of Rootchain, bound to a specific deployed contract.
func NewRootchainFilterer(address common.Address, filterer bind.ContractFilterer) (*RootchainFilterer, error) {
	contract, err := bindRootchain(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &RootchainFilterer{contract: contract}, nil
}

// bindRootchain binds a generic wrapper to an already deployed contract.
func bindRootchain(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(RootchainABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rootchain *RootchainRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Rootchain.Contract.RootchainCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rootchain *RootchainRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rootchain.Contract.RootchainTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rootchain *RootchainRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rootchain.Contract.RootchainTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Rootchain *RootchainCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Rootchain.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Rootchain *RootchainTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rootchain.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Rootchain *RootchainTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Rootchain.Contract.contract.Transact(opts, method, params...)
}

// Chain is a free data retrieval call binding the contract method 0xc763e5a1.
//
// Solidity: function chain() constant returns(bytes32)
func (_Rootchain *RootchainCaller) Chain(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "chain")
	return *ret0, err
}

// Chain is a free data retrieval call binding the contract method 0xc763e5a1.
//
// Solidity: function chain() constant returns(bytes32)
func (_Rootchain *RootchainSession) Chain() ([32]byte, error) {
	return _Rootchain.Contract.Chain(&_Rootchain.CallOpts)
}

// Chain is a free data retrieval call binding the contract method 0xc763e5a1.
//
// Solidity: function chain() constant returns(bytes32)
func (_Rootchain *RootchainCallerSession) Chain() ([32]byte, error) {
	return _Rootchain.Contract.Chain(&_Rootchain.CallOpts)
}

// ChildChainContract is a free data retrieval call binding the contract method 0x242fbb04.
//
// Solidity: function childChainContract() constant returns(address)
func (_Rootchain *RootchainCaller) ChildChainContract(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "childChainContract")
	return *ret0, err
}

// ChildChainContract is a free data retrieval call binding the contract method 0x242fbb04.
//
// Solidity: function childChainContract() constant returns(address)
func (_Rootchain *RootchainSession) ChildChainContract() (common.Address, error) {
	return _Rootchain.Contract.ChildChainContract(&_Rootchain.CallOpts)
}

// ChildChainContract is a free data retrieval call binding the contract method 0x242fbb04.
//
// Solidity: function childChainContract() constant returns(address)
func (_Rootchain *RootchainCallerSession) ChildChainContract() (common.Address, error) {
	return _Rootchain.Contract.ChildChainContract(&_Rootchain.CallOpts)
}

// CurrentChildBlock is a free data retrieval call binding the contract method 0x7a95f1e8.
//
// Solidity: function currentChildBlock() constant returns(uint256)
func (_Rootchain *RootchainCaller) CurrentChildBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "currentChildBlock")
	return *ret0, err
}

// CurrentChildBlock is a free data retrieval call binding the contract method 0x7a95f1e8.
//
// Solidity: function currentChildBlock() constant returns(uint256)
func (_Rootchain *RootchainSession) CurrentChildBlock() (*big.Int, error) {
	return _Rootchain.Contract.CurrentChildBlock(&_Rootchain.CallOpts)
}

// CurrentChildBlock is a free data retrieval call binding the contract method 0x7a95f1e8.
//
// Solidity: function currentChildBlock() constant returns(uint256)
func (_Rootchain *RootchainCallerSession) CurrentChildBlock() (*big.Int, error) {
	return _Rootchain.Contract.CurrentChildBlock(&_Rootchain.CallOpts)
}

// CurrentHeaderBlock is a free data retrieval call binding the contract method 0xec7e4855.
//
// Solidity: function currentHeaderBlock() constant returns(uint256)
func (_Rootchain *RootchainCaller) CurrentHeaderBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "currentHeaderBlock")
	return *ret0, err
}

// CurrentHeaderBlock is a free data retrieval call binding the contract method 0xec7e4855.
//
// Solidity: function currentHeaderBlock() constant returns(uint256)
func (_Rootchain *RootchainSession) CurrentHeaderBlock() (*big.Int, error) {
	return _Rootchain.Contract.CurrentHeaderBlock(&_Rootchain.CallOpts)
}

// CurrentHeaderBlock is a free data retrieval call binding the contract method 0xec7e4855.
//
// Solidity: function currentHeaderBlock() constant returns(uint256)
func (_Rootchain *RootchainCallerSession) CurrentHeaderBlock() (*big.Int, error) {
	return _Rootchain.Contract.CurrentHeaderBlock(&_Rootchain.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() constant returns(uint256)
func (_Rootchain *RootchainCaller) DepositCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "depositCount")
	return *ret0, err
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() constant returns(uint256)
func (_Rootchain *RootchainSession) DepositCount() (*big.Int, error) {
	return _Rootchain.Contract.DepositCount(&_Rootchain.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() constant returns(uint256)
func (_Rootchain *RootchainCallerSession) DepositCount() (*big.Int, error) {
	return _Rootchain.Contract.DepositCount(&_Rootchain.CallOpts)
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits( uint256) constant returns(header uint256, owner address, token address, amount uint256)
func (_Rootchain *RootchainCaller) Deposits(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Header *big.Int
	Owner  common.Address
	Token  common.Address
	Amount *big.Int
}, error) {
	ret := new(struct {
		Header *big.Int
		Owner  common.Address
		Token  common.Address
		Amount *big.Int
	})
	out := ret
	err := _Rootchain.contract.Call(opts, out, "deposits", arg0)
	return *ret, err
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits( uint256) constant returns(header uint256, owner address, token address, amount uint256)
func (_Rootchain *RootchainSession) Deposits(arg0 *big.Int) (struct {
	Header *big.Int
	Owner  common.Address
	Token  common.Address
	Amount *big.Int
}, error) {
	return _Rootchain.Contract.Deposits(&_Rootchain.CallOpts, arg0)
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits( uint256) constant returns(header uint256, owner address, token address, amount uint256)
func (_Rootchain *RootchainCallerSession) Deposits(arg0 *big.Int) (struct {
	Header *big.Int
	Owner  common.Address
	Token  common.Address
	Amount *big.Int
}, error) {
	return _Rootchain.Contract.Deposits(&_Rootchain.CallOpts, arg0)
}

// GetDepositBlock is a free data retrieval call binding the contract method 0x9a52d6ea.
//
// Solidity: function getDepositBlock(depositCount uint256) constant returns(header uint256, owner address, token address, amount uint256)
func (_Rootchain *RootchainCaller) GetDepositBlock(opts *bind.CallOpts, depositCount *big.Int) (struct {
	Header *big.Int
	Owner  common.Address
	Token  common.Address
	Amount *big.Int
}, error) {
	ret := new(struct {
		Header *big.Int
		Owner  common.Address
		Token  common.Address
		Amount *big.Int
	})
	out := ret
	err := _Rootchain.contract.Call(opts, out, "getDepositBlock", depositCount)
	return *ret, err
}

// GetDepositBlock is a free data retrieval call binding the contract method 0x9a52d6ea.
//
// Solidity: function getDepositBlock(depositCount uint256) constant returns(header uint256, owner address, token address, amount uint256)
func (_Rootchain *RootchainSession) GetDepositBlock(depositCount *big.Int) (struct {
	Header *big.Int
	Owner  common.Address
	Token  common.Address
	Amount *big.Int
}, error) {
	return _Rootchain.Contract.GetDepositBlock(&_Rootchain.CallOpts, depositCount)
}

// GetDepositBlock is a free data retrieval call binding the contract method 0x9a52d6ea.
//
// Solidity: function getDepositBlock(depositCount uint256) constant returns(header uint256, owner address, token address, amount uint256)
func (_Rootchain *RootchainCallerSession) GetDepositBlock(depositCount *big.Int) (struct {
	Header *big.Int
	Owner  common.Address
	Token  common.Address
	Amount *big.Int
}, error) {
	return _Rootchain.Contract.GetDepositBlock(&_Rootchain.CallOpts, depositCount)
}

// GetHeaderBlock is a free data retrieval call binding the contract method 0x313224a7.
//
// Solidity: function getHeaderBlock(headerNumber uint256) constant returns(root bytes32, start uint256, end uint256, createdAt uint256)
func (_Rootchain *RootchainCaller) GetHeaderBlock(opts *bind.CallOpts, headerNumber *big.Int) (struct {
	Root      [32]byte
	Start     *big.Int
	End       *big.Int
	CreatedAt *big.Int
}, error) {
	ret := new(struct {
		Root      [32]byte
		Start     *big.Int
		End       *big.Int
		CreatedAt *big.Int
	})
	out := ret
	err := _Rootchain.contract.Call(opts, out, "getHeaderBlock", headerNumber)
	return *ret, err
}

// GetHeaderBlock is a free data retrieval call binding the contract method 0x313224a7.
//
// Solidity: function getHeaderBlock(headerNumber uint256) constant returns(root bytes32, start uint256, end uint256, createdAt uint256)
func (_Rootchain *RootchainSession) GetHeaderBlock(headerNumber *big.Int) (struct {
	Root      [32]byte
	Start     *big.Int
	End       *big.Int
	CreatedAt *big.Int
}, error) {
	return _Rootchain.Contract.GetHeaderBlock(&_Rootchain.CallOpts, headerNumber)
}

// GetHeaderBlock is a free data retrieval call binding the contract method 0x313224a7.
//
// Solidity: function getHeaderBlock(headerNumber uint256) constant returns(root bytes32, start uint256, end uint256, createdAt uint256)
func (_Rootchain *RootchainCallerSession) GetHeaderBlock(headerNumber *big.Int) (struct {
	Root      [32]byte
	Start     *big.Int
	End       *big.Int
	CreatedAt *big.Int
}, error) {
	return _Rootchain.Contract.GetHeaderBlock(&_Rootchain.CallOpts, headerNumber)
}

// HeaderBlocks is a free data retrieval call binding the contract method 0x41539d4a.
//
// Solidity: function headerBlocks( uint256) constant returns(root bytes32, start uint256, end uint256, createdAt uint256)
func (_Rootchain *RootchainCaller) HeaderBlocks(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Root      [32]byte
	Start     *big.Int
	End       *big.Int
	CreatedAt *big.Int
}, error) {
	ret := new(struct {
		Root      [32]byte
		Start     *big.Int
		End       *big.Int
		CreatedAt *big.Int
	})
	out := ret
	err := _Rootchain.contract.Call(opts, out, "headerBlocks", arg0)
	return *ret, err
}

// HeaderBlocks is a free data retrieval call binding the contract method 0x41539d4a.
//
// Solidity: function headerBlocks( uint256) constant returns(root bytes32, start uint256, end uint256, createdAt uint256)
func (_Rootchain *RootchainSession) HeaderBlocks(arg0 *big.Int) (struct {
	Root      [32]byte
	Start     *big.Int
	End       *big.Int
	CreatedAt *big.Int
}, error) {
	return _Rootchain.Contract.HeaderBlocks(&_Rootchain.CallOpts, arg0)
}

// HeaderBlocks is a free data retrieval call binding the contract method 0x41539d4a.
//
// Solidity: function headerBlocks( uint256) constant returns(root bytes32, start uint256, end uint256, createdAt uint256)
func (_Rootchain *RootchainCallerSession) HeaderBlocks(arg0 *big.Int) (struct {
	Root      [32]byte
	Start     *big.Int
	End       *big.Int
	CreatedAt *big.Int
}, error) {
	return _Rootchain.Contract.HeaderBlocks(&_Rootchain.CallOpts, arg0)
}

// NetworkId is a free data retrieval call binding the contract method 0x9025e64c.
//
// Solidity: function networkId() constant returns(bytes)
func (_Rootchain *RootchainCaller) NetworkId(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "networkId")
	return *ret0, err
}

// NetworkId is a free data retrieval call binding the contract method 0x9025e64c.
//
// Solidity: function networkId() constant returns(bytes)
func (_Rootchain *RootchainSession) NetworkId() ([]byte, error) {
	return _Rootchain.Contract.NetworkId(&_Rootchain.CallOpts)
}

// NetworkId is a free data retrieval call binding the contract method 0x9025e64c.
//
// Solidity: function networkId() constant returns(bytes)
func (_Rootchain *RootchainCallerSession) NetworkId() ([]byte, error) {
	return _Rootchain.Contract.NetworkId(&_Rootchain.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Rootchain *RootchainCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Rootchain *RootchainSession) Owner() (common.Address, error) {
	return _Rootchain.Contract.Owner(&_Rootchain.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Rootchain *RootchainCallerSession) Owner() (common.Address, error) {
	return _Rootchain.Contract.Owner(&_Rootchain.CallOpts)
}

// ReverseTokens is a free data retrieval call binding the contract method 0x40828ebf.
//
// Solidity: function reverseTokens( address) constant returns(address)
func (_Rootchain *RootchainCaller) ReverseTokens(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "reverseTokens", arg0)
	return *ret0, err
}

// ReverseTokens is a free data retrieval call binding the contract method 0x40828ebf.
//
// Solidity: function reverseTokens( address) constant returns(address)
func (_Rootchain *RootchainSession) ReverseTokens(arg0 common.Address) (common.Address, error) {
	return _Rootchain.Contract.ReverseTokens(&_Rootchain.CallOpts, arg0)
}

// ReverseTokens is a free data retrieval call binding the contract method 0x40828ebf.
//
// Solidity: function reverseTokens( address) constant returns(address)
func (_Rootchain *RootchainCallerSession) ReverseTokens(arg0 common.Address) (common.Address, error) {
	return _Rootchain.Contract.ReverseTokens(&_Rootchain.CallOpts, arg0)
}

// StakeManager is a free data retrieval call binding the contract method 0x7542ff95.
//
// Solidity: function stakeManager() constant returns(address)
func (_Rootchain *RootchainCaller) StakeManager(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "stakeManager")
	return *ret0, err
}

// StakeManager is a free data retrieval call binding the contract method 0x7542ff95.
//
// Solidity: function stakeManager() constant returns(address)
func (_Rootchain *RootchainSession) StakeManager() (common.Address, error) {
	return _Rootchain.Contract.StakeManager(&_Rootchain.CallOpts)
}

// StakeManager is a free data retrieval call binding the contract method 0x7542ff95.
//
// Solidity: function stakeManager() constant returns(address)
func (_Rootchain *RootchainCallerSession) StakeManager() (common.Address, error) {
	return _Rootchain.Contract.StakeManager(&_Rootchain.CallOpts)
}

// Tokens is a free data retrieval call binding the contract method 0xe4860339.
//
// Solidity: function tokens( address) constant returns(address)
func (_Rootchain *RootchainCaller) Tokens(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "tokens", arg0)
	return *ret0, err
}

// Tokens is a free data retrieval call binding the contract method 0xe4860339.
//
// Solidity: function tokens( address) constant returns(address)
func (_Rootchain *RootchainSession) Tokens(arg0 common.Address) (common.Address, error) {
	return _Rootchain.Contract.Tokens(&_Rootchain.CallOpts, arg0)
}

// Tokens is a free data retrieval call binding the contract method 0xe4860339.
//
// Solidity: function tokens( address) constant returns(address)
func (_Rootchain *RootchainCallerSession) Tokens(arg0 common.Address) (common.Address, error) {
	return _Rootchain.Contract.Tokens(&_Rootchain.CallOpts, arg0)
}

// ValidatorContracts is a free data retrieval call binding the contract method 0x93d26c1a.
//
// Solidity: function validatorContracts( address) constant returns(bool)
func (_Rootchain *RootchainCaller) ValidatorContracts(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "validatorContracts", arg0)
	return *ret0, err
}

// ValidatorContracts is a free data retrieval call binding the contract method 0x93d26c1a.
//
// Solidity: function validatorContracts( address) constant returns(bool)
func (_Rootchain *RootchainSession) ValidatorContracts(arg0 common.Address) (bool, error) {
	return _Rootchain.Contract.ValidatorContracts(&_Rootchain.CallOpts, arg0)
}

// ValidatorContracts is a free data retrieval call binding the contract method 0x93d26c1a.
//
// Solidity: function validatorContracts( address) constant returns(bool)
func (_Rootchain *RootchainCallerSession) ValidatorContracts(arg0 common.Address) (bool, error) {
	return _Rootchain.Contract.ValidatorContracts(&_Rootchain.CallOpts, arg0)
}

// WethToken is a free data retrieval call binding the contract method 0x4b57b0be.
//
// Solidity: function wethToken() constant returns(address)
func (_Rootchain *RootchainCaller) WethToken(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "wethToken")
	return *ret0, err
}

// WethToken is a free data retrieval call binding the contract method 0x4b57b0be.
//
// Solidity: function wethToken() constant returns(address)
func (_Rootchain *RootchainSession) WethToken() (common.Address, error) {
	return _Rootchain.Contract.WethToken(&_Rootchain.CallOpts)
}

// WethToken is a free data retrieval call binding the contract method 0x4b57b0be.
//
// Solidity: function wethToken() constant returns(address)
func (_Rootchain *RootchainCallerSession) WethToken() (common.Address, error) {
	return _Rootchain.Contract.WethToken(&_Rootchain.CallOpts)
}

// WithdrawEventSignature is a free data retrieval call binding the contract method 0xe40b2775.
//
// Solidity: function withdrawEventSignature() constant returns(bytes32)
func (_Rootchain *RootchainCaller) WithdrawEventSignature(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "withdrawEventSignature")
	return *ret0, err
}

// WithdrawEventSignature is a free data retrieval call binding the contract method 0xe40b2775.
//
// Solidity: function withdrawEventSignature() constant returns(bytes32)
func (_Rootchain *RootchainSession) WithdrawEventSignature() ([32]byte, error) {
	return _Rootchain.Contract.WithdrawEventSignature(&_Rootchain.CallOpts)
}

// WithdrawEventSignature is a free data retrieval call binding the contract method 0xe40b2775.
//
// Solidity: function withdrawEventSignature() constant returns(bytes32)
func (_Rootchain *RootchainCallerSession) WithdrawEventSignature() ([32]byte, error) {
	return _Rootchain.Contract.WithdrawEventSignature(&_Rootchain.CallOpts)
}

// WithdrawSignature is a free data retrieval call binding the contract method 0x7a59a6f3.
//
// Solidity: function withdrawSignature() constant returns(bytes4)
func (_Rootchain *RootchainCaller) WithdrawSignature(opts *bind.CallOpts) ([4]byte, error) {
	var (
		ret0 = new([4]byte)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "withdrawSignature")
	return *ret0, err
}

// WithdrawSignature is a free data retrieval call binding the contract method 0x7a59a6f3.
//
// Solidity: function withdrawSignature() constant returns(bytes4)
func (_Rootchain *RootchainSession) WithdrawSignature() ([4]byte, error) {
	return _Rootchain.Contract.WithdrawSignature(&_Rootchain.CallOpts)
}

// WithdrawSignature is a free data retrieval call binding the contract method 0x7a59a6f3.
//
// Solidity: function withdrawSignature() constant returns(bytes4)
func (_Rootchain *RootchainCallerSession) WithdrawSignature() ([4]byte, error) {
	return _Rootchain.Contract.WithdrawSignature(&_Rootchain.CallOpts)
}

// Withdraws is a free data retrieval call binding the contract method 0xe09ab428.
//
// Solidity: function withdraws( bytes32) constant returns(bool)
func (_Rootchain *RootchainCaller) Withdraws(opts *bind.CallOpts, arg0 [32]byte) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "withdraws", arg0)
	return *ret0, err
}

// Withdraws is a free data retrieval call binding the contract method 0xe09ab428.
//
// Solidity: function withdraws( bytes32) constant returns(bool)
func (_Rootchain *RootchainSession) Withdraws(arg0 [32]byte) (bool, error) {
	return _Rootchain.Contract.Withdraws(&_Rootchain.CallOpts, arg0)
}

// Withdraws is a free data retrieval call binding the contract method 0xe09ab428.
//
// Solidity: function withdraws( bytes32) constant returns(bool)
func (_Rootchain *RootchainCallerSession) Withdraws(arg0 [32]byte) (bool, error) {
	return _Rootchain.Contract.Withdraws(&_Rootchain.CallOpts, arg0)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(_validator address) returns()
func (_Rootchain *RootchainTransactor) AddValidator(opts *bind.TransactOpts, _validator common.Address) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "addValidator", _validator)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(_validator address) returns()
func (_Rootchain *RootchainSession) AddValidator(_validator common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.AddValidator(&_Rootchain.TransactOpts, _validator)
}

// AddValidator is a paid mutator transaction binding the contract method 0x4d238c8e.
//
// Solidity: function addValidator(_validator address) returns()
func (_Rootchain *RootchainTransactorSession) AddValidator(_validator common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.AddValidator(&_Rootchain.TransactOpts, _validator)
}

// Deposit is a paid mutator transaction binding the contract method 0x8340f549.
//
// Solidity: function deposit(token address, user address, amount uint256) returns()
func (_Rootchain *RootchainTransactor) Deposit(opts *bind.TransactOpts, token common.Address, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "deposit", token, user, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x8340f549.
//
// Solidity: function deposit(token address, user address, amount uint256) returns()
func (_Rootchain *RootchainSession) Deposit(token common.Address, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rootchain.Contract.Deposit(&_Rootchain.TransactOpts, token, user, amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x8340f549.
//
// Solidity: function deposit(token address, user address, amount uint256) returns()
func (_Rootchain *RootchainTransactorSession) Deposit(token common.Address, user common.Address, amount *big.Int) (*types.Transaction, error) {
	return _Rootchain.Contract.Deposit(&_Rootchain.TransactOpts, token, user, amount)
}

// DepositEthers is a paid mutator transaction binding the contract method 0xf477a6b7.
//
// Solidity: function depositEthers() returns()
func (_Rootchain *RootchainTransactor) DepositEthers(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "depositEthers")
}

// DepositEthers is a paid mutator transaction binding the contract method 0xf477a6b7.
//
// Solidity: function depositEthers() returns()
func (_Rootchain *RootchainSession) DepositEthers() (*types.Transaction, error) {
	return _Rootchain.Contract.DepositEthers(&_Rootchain.TransactOpts)
}

// DepositEthers is a paid mutator transaction binding the contract method 0xf477a6b7.
//
// Solidity: function depositEthers() returns()
func (_Rootchain *RootchainTransactorSession) DepositEthers() (*types.Transaction, error) {
	return _Rootchain.Contract.DepositEthers(&_Rootchain.TransactOpts)
}

// MapToken is a paid mutator transaction binding the contract method 0x47400269.
//
// Solidity: function mapToken(rootToken address, childToken address) returns()
func (_Rootchain *RootchainTransactor) MapToken(opts *bind.TransactOpts, rootToken common.Address, childToken common.Address) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "mapToken", rootToken, childToken)
}

// MapToken is a paid mutator transaction binding the contract method 0x47400269.
//
// Solidity: function mapToken(rootToken address, childToken address) returns()
func (_Rootchain *RootchainSession) MapToken(rootToken common.Address, childToken common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.MapToken(&_Rootchain.TransactOpts, rootToken, childToken)
}

// MapToken is a paid mutator transaction binding the contract method 0x47400269.
//
// Solidity: function mapToken(rootToken address, childToken address) returns()
func (_Rootchain *RootchainTransactorSession) MapToken(rootToken common.Address, childToken common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.MapToken(&_Rootchain.TransactOpts, rootToken, childToken)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(_validator address) returns()
func (_Rootchain *RootchainTransactor) RemoveValidator(opts *bind.TransactOpts, _validator common.Address) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "removeValidator", _validator)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(_validator address) returns()
func (_Rootchain *RootchainSession) RemoveValidator(_validator common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.RemoveValidator(&_Rootchain.TransactOpts, _validator)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0x40a141ff.
//
// Solidity: function removeValidator(_validator address) returns()
func (_Rootchain *RootchainTransactorSession) RemoveValidator(_validator common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.RemoveValidator(&_Rootchain.TransactOpts, _validator)
}

// SetChildContract is a paid mutator transaction binding the contract method 0xf8d86e18.
//
// Solidity: function setChildContract(newChildChain address) returns()
func (_Rootchain *RootchainTransactor) SetChildContract(opts *bind.TransactOpts, newChildChain common.Address) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "setChildContract", newChildChain)
}

// SetChildContract is a paid mutator transaction binding the contract method 0xf8d86e18.
//
// Solidity: function setChildContract(newChildChain address) returns()
func (_Rootchain *RootchainSession) SetChildContract(newChildChain common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.SetChildContract(&_Rootchain.TransactOpts, newChildChain)
}

// SetChildContract is a paid mutator transaction binding the contract method 0xf8d86e18.
//
// Solidity: function setChildContract(newChildChain address) returns()
func (_Rootchain *RootchainTransactorSession) SetChildContract(newChildChain common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.SetChildContract(&_Rootchain.TransactOpts, newChildChain)
}

// SetStakeManager is a paid mutator transaction binding the contract method 0x0e7c67fc.
//
// Solidity: function setStakeManager(_stakeManager address) returns()
func (_Rootchain *RootchainTransactor) SetStakeManager(opts *bind.TransactOpts, _stakeManager common.Address) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "setStakeManager", _stakeManager)
}

// SetStakeManager is a paid mutator transaction binding the contract method 0x0e7c67fc.
//
// Solidity: function setStakeManager(_stakeManager address) returns()
func (_Rootchain *RootchainSession) SetStakeManager(_stakeManager common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.SetStakeManager(&_Rootchain.TransactOpts, _stakeManager)
}

// SetStakeManager is a paid mutator transaction binding the contract method 0x0e7c67fc.
//
// Solidity: function setStakeManager(_stakeManager address) returns()
func (_Rootchain *RootchainTransactorSession) SetStakeManager(_stakeManager common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.SetStakeManager(&_Rootchain.TransactOpts, _stakeManager)
}

// SetWETHToken is a paid mutator transaction binding the contract method 0xb45d1f68.
//
// Solidity: function setWETHToken(_token address) returns()
func (_Rootchain *RootchainTransactor) SetWETHToken(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "setWETHToken", _token)
}

// SetWETHToken is a paid mutator transaction binding the contract method 0xb45d1f68.
//
// Solidity: function setWETHToken(_token address) returns()
func (_Rootchain *RootchainSession) SetWETHToken(_token common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.SetWETHToken(&_Rootchain.TransactOpts, _token)
}

// SetWETHToken is a paid mutator transaction binding the contract method 0xb45d1f68.
//
// Solidity: function setWETHToken(_token address) returns()
func (_Rootchain *RootchainTransactorSession) SetWETHToken(_token common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.SetWETHToken(&_Rootchain.TransactOpts, _token)
}

// Slash is a paid mutator transaction binding the contract method 0x2da25de3.
//
// Solidity: function slash() returns()
func (_Rootchain *RootchainTransactor) Slash(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "slash")
}

// Slash is a paid mutator transaction binding the contract method 0x2da25de3.
//
// Solidity: function slash() returns()
func (_Rootchain *RootchainSession) Slash() (*types.Transaction, error) {
	return _Rootchain.Contract.Slash(&_Rootchain.TransactOpts)
}

// Slash is a paid mutator transaction binding the contract method 0x2da25de3.
//
// Solidity: function slash() returns()
func (_Rootchain *RootchainTransactorSession) Slash() (*types.Transaction, error) {
	return _Rootchain.Contract.Slash(&_Rootchain.TransactOpts)
}

// SubmitHeaderBlock is a paid mutator transaction binding the contract method 0xa0ee0f19.
//
// Solidity: function submitHeaderBlock(root bytes32, start uint256, end uint256, sigs bytes) returns()
func (_Rootchain *RootchainTransactor) SubmitHeaderBlock(opts *bind.TransactOpts, root [32]byte, start *big.Int, end *big.Int, sigs []byte) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "submitHeaderBlock", root, start, end, sigs)
}

// SubmitHeaderBlock is a paid mutator transaction binding the contract method 0xa0ee0f19.
//
// Solidity: function submitHeaderBlock(root bytes32, start uint256, end uint256, sigs bytes) returns()
func (_Rootchain *RootchainSession) SubmitHeaderBlock(root [32]byte, start *big.Int, end *big.Int, sigs []byte) (*types.Transaction, error) {
	return _Rootchain.Contract.SubmitHeaderBlock(&_Rootchain.TransactOpts, root, start, end, sigs)
}

// SubmitHeaderBlock is a paid mutator transaction binding the contract method 0xa0ee0f19.
//
// Solidity: function submitHeaderBlock(root bytes32, start uint256, end uint256, sigs bytes) returns()
func (_Rootchain *RootchainTransactorSession) SubmitHeaderBlock(root [32]byte, start *big.Int, end *big.Int, sigs []byte) (*types.Transaction, error) {
	return _Rootchain.Contract.SubmitHeaderBlock(&_Rootchain.TransactOpts, root, start, end, sigs)
}

// TokenFallback is a paid mutator transaction binding the contract method 0xc0ee0b8a.
//
// Solidity: function tokenFallback(_sender address, _value uint256,  bytes) returns()
func (_Rootchain *RootchainTransactor) TokenFallback(opts *bind.TransactOpts, _sender common.Address, _value *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "tokenFallback", _sender, _value, arg2)
}

// TokenFallback is a paid mutator transaction binding the contract method 0xc0ee0b8a.
//
// Solidity: function tokenFallback(_sender address, _value uint256,  bytes) returns()
func (_Rootchain *RootchainSession) TokenFallback(_sender common.Address, _value *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _Rootchain.Contract.TokenFallback(&_Rootchain.TransactOpts, _sender, _value, arg2)
}

// TokenFallback is a paid mutator transaction binding the contract method 0xc0ee0b8a.
//
// Solidity: function tokenFallback(_sender address, _value uint256,  bytes) returns()
func (_Rootchain *RootchainTransactorSession) TokenFallback(_sender common.Address, _value *big.Int, arg2 []byte) (*types.Transaction, error) {
	return _Rootchain.Contract.TokenFallback(&_Rootchain.TransactOpts, _sender, _value, arg2)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Rootchain *RootchainTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Rootchain *RootchainSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.TransferOwnership(&_Rootchain.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(newOwner address) returns()
func (_Rootchain *RootchainTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.TransferOwnership(&_Rootchain.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0xfce5ff2e.
//
// Solidity: function withdraw(headerNumber uint256, headerProof bytes, blockNumber uint256, blockTime uint256, txRoot bytes32, receiptRoot bytes32, path bytes, txBytes bytes, txProof bytes, receiptBytes bytes, receiptProof bytes) returns()
func (_Rootchain *RootchainTransactor) Withdraw(opts *bind.TransactOpts, headerNumber *big.Int, headerProof []byte, blockNumber *big.Int, blockTime *big.Int, txRoot [32]byte, receiptRoot [32]byte, path []byte, txBytes []byte, txProof []byte, receiptBytes []byte, receiptProof []byte) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "withdraw", headerNumber, headerProof, blockNumber, blockTime, txRoot, receiptRoot, path, txBytes, txProof, receiptBytes, receiptProof)
}

// Withdraw is a paid mutator transaction binding the contract method 0xfce5ff2e.
//
// Solidity: function withdraw(headerNumber uint256, headerProof bytes, blockNumber uint256, blockTime uint256, txRoot bytes32, receiptRoot bytes32, path bytes, txBytes bytes, txProof bytes, receiptBytes bytes, receiptProof bytes) returns()
func (_Rootchain *RootchainSession) Withdraw(headerNumber *big.Int, headerProof []byte, blockNumber *big.Int, blockTime *big.Int, txRoot [32]byte, receiptRoot [32]byte, path []byte, txBytes []byte, txProof []byte, receiptBytes []byte, receiptProof []byte) (*types.Transaction, error) {
	return _Rootchain.Contract.Withdraw(&_Rootchain.TransactOpts, headerNumber, headerProof, blockNumber, blockTime, txRoot, receiptRoot, path, txBytes, txProof, receiptBytes, receiptProof)
}

// Withdraw is a paid mutator transaction binding the contract method 0xfce5ff2e.
//
// Solidity: function withdraw(headerNumber uint256, headerProof bytes, blockNumber uint256, blockTime uint256, txRoot bytes32, receiptRoot bytes32, path bytes, txBytes bytes, txProof bytes, receiptBytes bytes, receiptProof bytes) returns()
func (_Rootchain *RootchainTransactorSession) Withdraw(headerNumber *big.Int, headerProof []byte, blockNumber *big.Int, blockTime *big.Int, txRoot [32]byte, receiptRoot [32]byte, path []byte, txBytes []byte, txProof []byte, receiptBytes []byte, receiptProof []byte) (*types.Transaction, error) {
	return _Rootchain.Contract.Withdraw(&_Rootchain.TransactOpts, headerNumber, headerProof, blockNumber, blockTime, txRoot, receiptRoot, path, txBytes, txProof, receiptBytes, receiptProof)
}

// RootchainChildChainChangedIterator is returned from FilterChildChainChanged and is used to iterate over the raw logs and unpacked data for ChildChainChanged events raised by the Rootchain contract.
type RootchainChildChainChangedIterator struct {
	Event *RootchainChildChainChanged // Event containing the contract specifics and raw log

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
func (it *RootchainChildChainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootchainChildChainChanged)
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
		it.Event = new(RootchainChildChainChanged)
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
func (it *RootchainChildChainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootchainChildChainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootchainChildChainChanged represents a ChildChainChanged event raised by the Rootchain contract.
type RootchainChildChainChanged struct {
	PreviousChildChain common.Address
	NewChildChain      common.Address
	Raw                types.Log // Blockchain specific contextual infos
}

// FilterChildChainChanged is a free log retrieval operation binding the contract event 0x1f9f3556dd336016cdf20adaead7d5c73665dba664b60e8c17e9a4eb91ce1d39.
//
// Solidity: e ChildChainChanged(previousChildChain indexed address, newChildChain indexed address)
func (_Rootchain *RootchainFilterer) FilterChildChainChanged(opts *bind.FilterOpts, previousChildChain []common.Address, newChildChain []common.Address) (*RootchainChildChainChangedIterator, error) {

	var previousChildChainRule []interface{}
	for _, previousChildChainItem := range previousChildChain {
		previousChildChainRule = append(previousChildChainRule, previousChildChainItem)
	}
	var newChildChainRule []interface{}
	for _, newChildChainItem := range newChildChain {
		newChildChainRule = append(newChildChainRule, newChildChainItem)
	}

	logs, sub, err := _Rootchain.contract.FilterLogs(opts, "ChildChainChanged", previousChildChainRule, newChildChainRule)
	if err != nil {
		return nil, err
	}
	return &RootchainChildChainChangedIterator{contract: _Rootchain.contract, event: "ChildChainChanged", logs: logs, sub: sub}, nil
}

// WatchChildChainChanged is a free log subscription operation binding the contract event 0x1f9f3556dd336016cdf20adaead7d5c73665dba664b60e8c17e9a4eb91ce1d39.
//
// Solidity: e ChildChainChanged(previousChildChain indexed address, newChildChain indexed address)
func (_Rootchain *RootchainFilterer) WatchChildChainChanged(opts *bind.WatchOpts, sink chan<- *RootchainChildChainChanged, previousChildChain []common.Address, newChildChain []common.Address) (event.Subscription, error) {

	var previousChildChainRule []interface{}
	for _, previousChildChainItem := range previousChildChain {
		previousChildChainRule = append(previousChildChainRule, previousChildChainItem)
	}
	var newChildChainRule []interface{}
	for _, newChildChainItem := range newChildChain {
		newChildChainRule = append(newChildChainRule, newChildChainItem)
	}

	logs, sub, err := _Rootchain.contract.WatchLogs(opts, "ChildChainChanged", previousChildChainRule, newChildChainRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootchainChildChainChanged)
				if err := _Rootchain.contract.UnpackLog(event, "ChildChainChanged", log); err != nil {
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

// RootchainDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the Rootchain contract.
type RootchainDepositIterator struct {
	Event *RootchainDeposit // Event containing the contract specifics and raw log

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
func (it *RootchainDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootchainDeposit)
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
		it.Event = new(RootchainDeposit)
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
func (it *RootchainDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootchainDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootchainDeposit represents a Deposit event raised by the Rootchain contract.
type RootchainDeposit struct {
	User         common.Address
	Token        common.Address
	Amount       *big.Int
	DepositCount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7.
//
// Solidity: e Deposit(user indexed address, token indexed address, amount uint256, depositCount uint256)
func (_Rootchain *RootchainFilterer) FilterDeposit(opts *bind.FilterOpts, user []common.Address, token []common.Address) (*RootchainDepositIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Rootchain.contract.FilterLogs(opts, "Deposit", userRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &RootchainDepositIterator{contract: _Rootchain.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7.
//
// Solidity: e Deposit(user indexed address, token indexed address, amount uint256, depositCount uint256)
func (_Rootchain *RootchainFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *RootchainDeposit, user []common.Address, token []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Rootchain.contract.WatchLogs(opts, "Deposit", userRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootchainDeposit)
				if err := _Rootchain.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// RootchainNewHeaderBlockIterator is returned from FilterNewHeaderBlock and is used to iterate over the raw logs and unpacked data for NewHeaderBlock events raised by the Rootchain contract.
type RootchainNewHeaderBlockIterator struct {
	Event *RootchainNewHeaderBlock // Event containing the contract specifics and raw log

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
func (it *RootchainNewHeaderBlockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootchainNewHeaderBlock)
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
		it.Event = new(RootchainNewHeaderBlock)
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
func (it *RootchainNewHeaderBlockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootchainNewHeaderBlockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootchainNewHeaderBlock represents a NewHeaderBlock event raised by the Rootchain contract.
type RootchainNewHeaderBlock struct {
	Proposer common.Address
	Number   *big.Int
	Start    *big.Int
	End      *big.Int
	Root     [32]byte
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterNewHeaderBlock is a free log retrieval operation binding the contract event 0xf146921b854b787ba7d6045e8a8054731dc62430ae16c4bf08147539b1b6ef8f.
//
// Solidity: e NewHeaderBlock(proposer indexed address, number uint256, start uint256, end uint256, root bytes32)
func (_Rootchain *RootchainFilterer) FilterNewHeaderBlock(opts *bind.FilterOpts, proposer []common.Address) (*RootchainNewHeaderBlockIterator, error) {

	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _Rootchain.contract.FilterLogs(opts, "NewHeaderBlock", proposerRule)
	if err != nil {
		return nil, err
	}
	return &RootchainNewHeaderBlockIterator{contract: _Rootchain.contract, event: "NewHeaderBlock", logs: logs, sub: sub}, nil
}

// WatchNewHeaderBlock is a free log subscription operation binding the contract event 0xf146921b854b787ba7d6045e8a8054731dc62430ae16c4bf08147539b1b6ef8f.
//
// Solidity: e NewHeaderBlock(proposer indexed address, number uint256, start uint256, end uint256, root bytes32)
func (_Rootchain *RootchainFilterer) WatchNewHeaderBlock(opts *bind.WatchOpts, sink chan<- *RootchainNewHeaderBlock, proposer []common.Address) (event.Subscription, error) {

	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}

	logs, sub, err := _Rootchain.contract.WatchLogs(opts, "NewHeaderBlock", proposerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootchainNewHeaderBlock)
				if err := _Rootchain.contract.UnpackLog(event, "NewHeaderBlock", log); err != nil {
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

// RootchainOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Rootchain contract.
type RootchainOwnershipTransferredIterator struct {
	Event *RootchainOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *RootchainOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootchainOwnershipTransferred)
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
		it.Event = new(RootchainOwnershipTransferred)
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
func (it *RootchainOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootchainOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootchainOwnershipTransferred represents a OwnershipTransferred event raised by the Rootchain contract.
type RootchainOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Rootchain *RootchainFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*RootchainOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rootchain.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &RootchainOwnershipTransferredIterator{contract: _Rootchain.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
func (_Rootchain *RootchainFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *RootchainOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Rootchain.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootchainOwnershipTransferred)
				if err := _Rootchain.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// RootchainTokenMappedIterator is returned from FilterTokenMapped and is used to iterate over the raw logs and unpacked data for TokenMapped events raised by the Rootchain contract.
type RootchainTokenMappedIterator struct {
	Event *RootchainTokenMapped // Event containing the contract specifics and raw log

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
func (it *RootchainTokenMappedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootchainTokenMapped)
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
		it.Event = new(RootchainTokenMapped)
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
func (it *RootchainTokenMappedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootchainTokenMappedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootchainTokenMapped represents a TokenMapped event raised by the Rootchain contract.
type RootchainTokenMapped struct {
	RootToken  common.Address
	ChildToken common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTokenMapped is a free log retrieval operation binding the contract event 0x85920d35e6c72f6b2affffa04298b0cecfeba86e4a9f407df661f1cb8ab5e617.
//
// Solidity: e TokenMapped(rootToken indexed address, childToken indexed address)
func (_Rootchain *RootchainFilterer) FilterTokenMapped(opts *bind.FilterOpts, rootToken []common.Address, childToken []common.Address) (*RootchainTokenMappedIterator, error) {

	var rootTokenRule []interface{}
	for _, rootTokenItem := range rootToken {
		rootTokenRule = append(rootTokenRule, rootTokenItem)
	}
	var childTokenRule []interface{}
	for _, childTokenItem := range childToken {
		childTokenRule = append(childTokenRule, childTokenItem)
	}

	logs, sub, err := _Rootchain.contract.FilterLogs(opts, "TokenMapped", rootTokenRule, childTokenRule)
	if err != nil {
		return nil, err
	}
	return &RootchainTokenMappedIterator{contract: _Rootchain.contract, event: "TokenMapped", logs: logs, sub: sub}, nil
}

// WatchTokenMapped is a free log subscription operation binding the contract event 0x85920d35e6c72f6b2affffa04298b0cecfeba86e4a9f407df661f1cb8ab5e617.
//
// Solidity: e TokenMapped(rootToken indexed address, childToken indexed address)
func (_Rootchain *RootchainFilterer) WatchTokenMapped(opts *bind.WatchOpts, sink chan<- *RootchainTokenMapped, rootToken []common.Address, childToken []common.Address) (event.Subscription, error) {

	var rootTokenRule []interface{}
	for _, rootTokenItem := range rootToken {
		rootTokenRule = append(rootTokenRule, rootTokenItem)
	}
	var childTokenRule []interface{}
	for _, childTokenItem := range childToken {
		childTokenRule = append(childTokenRule, childTokenItem)
	}

	logs, sub, err := _Rootchain.contract.WatchLogs(opts, "TokenMapped", rootTokenRule, childTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootchainTokenMapped)
				if err := _Rootchain.contract.UnpackLog(event, "TokenMapped", log); err != nil {
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

// RootchainValidatorAddedIterator is returned from FilterValidatorAdded and is used to iterate over the raw logs and unpacked data for ValidatorAdded events raised by the Rootchain contract.
type RootchainValidatorAddedIterator struct {
	Event *RootchainValidatorAdded // Event containing the contract specifics and raw log

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
func (it *RootchainValidatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootchainValidatorAdded)
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
		it.Event = new(RootchainValidatorAdded)
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
func (it *RootchainValidatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootchainValidatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootchainValidatorAdded represents a ValidatorAdded event raised by the Rootchain contract.
type RootchainValidatorAdded struct {
	Validator common.Address
	From      common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterValidatorAdded is a free log retrieval operation binding the contract event 0x8064a302796c89446a96d63470b5b036212da26bd2debe5bec73e0170a9a5e83.
//
// Solidity: e ValidatorAdded(validator indexed address, from indexed address)
func (_Rootchain *RootchainFilterer) FilterValidatorAdded(opts *bind.FilterOpts, validator []common.Address, from []common.Address) (*RootchainValidatorAddedIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Rootchain.contract.FilterLogs(opts, "ValidatorAdded", validatorRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &RootchainValidatorAddedIterator{contract: _Rootchain.contract, event: "ValidatorAdded", logs: logs, sub: sub}, nil
}

// WatchValidatorAdded is a free log subscription operation binding the contract event 0x8064a302796c89446a96d63470b5b036212da26bd2debe5bec73e0170a9a5e83.
//
// Solidity: e ValidatorAdded(validator indexed address, from indexed address)
func (_Rootchain *RootchainFilterer) WatchValidatorAdded(opts *bind.WatchOpts, sink chan<- *RootchainValidatorAdded, validator []common.Address, from []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Rootchain.contract.WatchLogs(opts, "ValidatorAdded", validatorRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootchainValidatorAdded)
				if err := _Rootchain.contract.UnpackLog(event, "ValidatorAdded", log); err != nil {
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

// RootchainValidatorRemovedIterator is returned from FilterValidatorRemoved and is used to iterate over the raw logs and unpacked data for ValidatorRemoved events raised by the Rootchain contract.
type RootchainValidatorRemovedIterator struct {
	Event *RootchainValidatorRemoved // Event containing the contract specifics and raw log

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
func (it *RootchainValidatorRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootchainValidatorRemoved)
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
		it.Event = new(RootchainValidatorRemoved)
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
func (it *RootchainValidatorRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootchainValidatorRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootchainValidatorRemoved represents a ValidatorRemoved event raised by the Rootchain contract.
type RootchainValidatorRemoved struct {
	Validator common.Address
	From      common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterValidatorRemoved is a free log retrieval operation binding the contract event 0x05fc5214912b853ccaaaeaf238daf240b8c6ae70c0f18ec215d5088f8f49d781.
//
// Solidity: e ValidatorRemoved(validator indexed address, from indexed address)
func (_Rootchain *RootchainFilterer) FilterValidatorRemoved(opts *bind.FilterOpts, validator []common.Address, from []common.Address) (*RootchainValidatorRemovedIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Rootchain.contract.FilterLogs(opts, "ValidatorRemoved", validatorRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &RootchainValidatorRemovedIterator{contract: _Rootchain.contract, event: "ValidatorRemoved", logs: logs, sub: sub}, nil
}

// WatchValidatorRemoved is a free log subscription operation binding the contract event 0x05fc5214912b853ccaaaeaf238daf240b8c6ae70c0f18ec215d5088f8f49d781.
//
// Solidity: e ValidatorRemoved(validator indexed address, from indexed address)
func (_Rootchain *RootchainFilterer) WatchValidatorRemoved(opts *bind.WatchOpts, sink chan<- *RootchainValidatorRemoved, validator []common.Address, from []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Rootchain.contract.WatchLogs(opts, "ValidatorRemoved", validatorRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootchainValidatorRemoved)
				if err := _Rootchain.contract.UnpackLog(event, "ValidatorRemoved", log); err != nil {
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

// RootchainWithdrawIterator is returned from FilterWithdraw and is used to iterate over the raw logs and unpacked data for Withdraw events raised by the Rootchain contract.
type RootchainWithdrawIterator struct {
	Event *RootchainWithdraw // Event containing the contract specifics and raw log

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
func (it *RootchainWithdrawIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootchainWithdraw)
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
		it.Event = new(RootchainWithdraw)
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
func (it *RootchainWithdrawIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootchainWithdrawIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootchainWithdraw represents a Withdraw event raised by the Rootchain contract.
type RootchainWithdraw struct {
	User   common.Address
	Token  common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterWithdraw is a free log retrieval operation binding the contract event 0x9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb.
//
// Solidity: e Withdraw(user indexed address, token indexed address, amount uint256)
func (_Rootchain *RootchainFilterer) FilterWithdraw(opts *bind.FilterOpts, user []common.Address, token []common.Address) (*RootchainWithdrawIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Rootchain.contract.FilterLogs(opts, "Withdraw", userRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return &RootchainWithdrawIterator{contract: _Rootchain.contract, event: "Withdraw", logs: logs, sub: sub}, nil
}

// WatchWithdraw is a free log subscription operation binding the contract event 0x9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb.
//
// Solidity: e Withdraw(user indexed address, token indexed address, amount uint256)
func (_Rootchain *RootchainFilterer) WatchWithdraw(opts *bind.WatchOpts, sink chan<- *RootchainWithdraw, user []common.Address, token []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var tokenRule []interface{}
	for _, tokenItem := range token {
		tokenRule = append(tokenRule, tokenItem)
	}

	logs, sub, err := _Rootchain.contract.WatchLogs(opts, "Withdraw", userRule, tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootchainWithdraw)
				if err := _Rootchain.contract.UnpackLog(event, "Withdraw", log); err != nil {
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
