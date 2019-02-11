// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rootchain

import (
	"math/big"
	"strings"

	ethereum "github.com/go-ethereum"
	"github.com/go-ethereum/accounts/abi"
	"github.com/go-ethereum/accounts/abi/bind"
	"github.com/go-ethereum/common"
	"github.com/go-ethereum/core/types"
	"github.com/go-ethereum/event"
)

// RootchainABI is the input ABI used to generate the binding from.
const RootchainABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_user\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"transferAmount\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_stakeManager\",\"type\":\"address\"}],\"name\":\"setStakeManager\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_withdrawManager\",\"type\":\"address\"}],\"name\":\"setWithdrawManager\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_depositManager\",\"type\":\"address\"}],\"name\":\"setDepositManager\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"childChainContract\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"roundType\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"slash\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_depositCount\",\"type\":\"uint256\"}],\"name\":\"depositBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_user\",\"type\":\"address\"},{\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"depositERC721\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"headerBlocks\",\"outputs\":[{\"name\":\"root\",\"type\":\"bytes32\"},{\"name\":\"start\",\"type\":\"uint256\"},{\"name\":\"end\",\"type\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\"},{\"name\":\"proposer\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_nftContract\",\"type\":\"address\"}],\"name\":\"setExitNFTContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"exitId\",\"type\":\"uint256\"}],\"name\":\"deleteExit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"removeProofValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_headerNumber\",\"type\":\"uint256\"}],\"name\":\"headerBlock\",\"outputs\":[{\"name\":\"_root\",\"type\":\"bytes32\"},{\"name\":\"_start\",\"type\":\"uint256\"},{\"name\":\"_end\",\"type\":\"uint256\"},{\"name\":\"_createdAt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"depositManager\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"stakeManager\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentChildBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"voteType\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes1\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_user\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"}],\"name\":\"deposit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"networkId\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CHILD_BLOCK_INTERVAL\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"setWETHToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_user\",\"type\":\"address\"},{\"name\":\"_amount\",\"type\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"tokenFallback\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"proofValidatorContracts\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chain\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_validator\",\"type\":\"address\"}],\"name\":\"addProofValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_rootToken\",\"type\":\"address\"},{\"name\":\"_childToken\",\"type\":\"address\"},{\"name\":\"_isERC721\",\"type\":\"bool\"}],\"name\":\"mapToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"withdrawManager\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentHeaderBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"vote\",\"type\":\"bytes\"},{\"name\":\"sigs\",\"type\":\"bytes\"},{\"name\":\"extradata\",\"type\":\"bytes\"}],\"name\":\"submitHeaderBlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"depositEthers\",\"outputs\":[],\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newChildChain\",\"type\":\"address\"}],\"name\":\"setChildContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"finalizeCommit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"payable\":true,\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousChildChain\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newChildChain\",\"type\":\"address\"}],\"name\":\"ChildChainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"}],\"name\":\"ProofValidatorAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"validator\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"}],\"name\":\"ProofValidatorRemoved\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"proposer\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"number\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"start\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"end\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"root\",\"type\":\"bytes32\"}],\"name\":\"NewHeaderBlock\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"}]"

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

// CHILDBLOCKINTERVAL is a free data retrieval call binding the contract method 0xa831fa07.
//
// Solidity: function CHILD_BLOCK_INTERVAL() constant returns(uint256)
func (_Rootchain *RootchainCaller) CHILDBLOCKINTERVAL(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "CHILD_BLOCK_INTERVAL")
	return *ret0, err
}

// CHILDBLOCKINTERVAL is a free data retrieval call binding the contract method 0xa831fa07.
//
// Solidity: function CHILD_BLOCK_INTERVAL() constant returns(uint256)
func (_Rootchain *RootchainSession) CHILDBLOCKINTERVAL() (*big.Int, error) {
	return _Rootchain.Contract.CHILDBLOCKINTERVAL(&_Rootchain.CallOpts)
}

// CHILDBLOCKINTERVAL is a free data retrieval call binding the contract method 0xa831fa07.
//
// Solidity: function CHILD_BLOCK_INTERVAL() constant returns(uint256)
func (_Rootchain *RootchainCallerSession) CHILDBLOCKINTERVAL() (*big.Int, error) {
	return _Rootchain.Contract.CHILDBLOCKINTERVAL(&_Rootchain.CallOpts)
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

// DepositBlock is a free data retrieval call binding the contract method 0x32590654.
//
// Solidity: function depositBlock(_depositCount uint256) constant returns(uint256, address, address, uint256, uint256)
func (_Rootchain *RootchainCaller) DepositBlock(opts *bind.CallOpts, _depositCount *big.Int) (*big.Int, common.Address, common.Address, *big.Int, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(common.Address)
		ret2 = new(common.Address)
		ret3 = new(*big.Int)
		ret4 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
	}
	err := _Rootchain.contract.Call(opts, out, "depositBlock", _depositCount)
	return *ret0, *ret1, *ret2, *ret3, *ret4, err
}

// DepositBlock is a free data retrieval call binding the contract method 0x32590654.
//
// Solidity: function depositBlock(_depositCount uint256) constant returns(uint256, address, address, uint256, uint256)
func (_Rootchain *RootchainSession) DepositBlock(_depositCount *big.Int) (*big.Int, common.Address, common.Address, *big.Int, *big.Int, error) {
	return _Rootchain.Contract.DepositBlock(&_Rootchain.CallOpts, _depositCount)
}

// DepositBlock is a free data retrieval call binding the contract method 0x32590654.
//
// Solidity: function depositBlock(_depositCount uint256) constant returns(uint256, address, address, uint256, uint256)
func (_Rootchain *RootchainCallerSession) DepositBlock(_depositCount *big.Int) (*big.Int, common.Address, common.Address, *big.Int, *big.Int, error) {
	return _Rootchain.Contract.DepositBlock(&_Rootchain.CallOpts, _depositCount)
}

// DepositManager is a free data retrieval call binding the contract method 0x6c7ac9d8.
//
// Solidity: function depositManager() constant returns(address)
func (_Rootchain *RootchainCaller) DepositManager(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "depositManager")
	return *ret0, err
}

// DepositManager is a free data retrieval call binding the contract method 0x6c7ac9d8.
//
// Solidity: function depositManager() constant returns(address)
func (_Rootchain *RootchainSession) DepositManager() (common.Address, error) {
	return _Rootchain.Contract.DepositManager(&_Rootchain.CallOpts)
}

// DepositManager is a free data retrieval call binding the contract method 0x6c7ac9d8.
//
// Solidity: function depositManager() constant returns(address)
func (_Rootchain *RootchainCallerSession) DepositManager() (common.Address, error) {
	return _Rootchain.Contract.DepositManager(&_Rootchain.CallOpts)
}

// HeaderBlock is a free data retrieval call binding the contract method 0x61bbd461.
//
// Solidity: function headerBlock(_headerNumber uint256) constant returns(_root bytes32, _start uint256, _end uint256, _createdAt uint256)
func (_Rootchain *RootchainCaller) HeaderBlock(opts *bind.CallOpts, _headerNumber *big.Int) (struct {
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
	err := _Rootchain.contract.Call(opts, out, "headerBlock", _headerNumber)
	return *ret, err
}

// HeaderBlock is a free data retrieval call binding the contract method 0x61bbd461.
//
// Solidity: function headerBlock(_headerNumber uint256) constant returns(_root bytes32, _start uint256, _end uint256, _createdAt uint256)
func (_Rootchain *RootchainSession) HeaderBlock(_headerNumber *big.Int) (struct {
	Root      [32]byte
	Start     *big.Int
	End       *big.Int
	CreatedAt *big.Int
}, error) {
	return _Rootchain.Contract.HeaderBlock(&_Rootchain.CallOpts, _headerNumber)
}

// HeaderBlock is a free data retrieval call binding the contract method 0x61bbd461.
//
// Solidity: function headerBlock(_headerNumber uint256) constant returns(_root bytes32, _start uint256, _end uint256, _createdAt uint256)
func (_Rootchain *RootchainCallerSession) HeaderBlock(_headerNumber *big.Int) (struct {
	Root      [32]byte
	Start     *big.Int
	End       *big.Int
	CreatedAt *big.Int
}, error) {
	return _Rootchain.Contract.HeaderBlock(&_Rootchain.CallOpts, _headerNumber)
}

// HeaderBlocks is a free data retrieval call binding the contract method 0x41539d4a.
//
// Solidity: function headerBlocks( uint256) constant returns(root bytes32, start uint256, end uint256, createdAt uint256, proposer address)
func (_Rootchain *RootchainCaller) HeaderBlocks(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Root      [32]byte
	Start     *big.Int
	End       *big.Int
	CreatedAt *big.Int
	Proposer  common.Address
}, error) {
	ret := new(struct {
		Root      [32]byte
		Start     *big.Int
		End       *big.Int
		CreatedAt *big.Int
		Proposer  common.Address
	})
	out := ret
	err := _Rootchain.contract.Call(opts, out, "headerBlocks", arg0)
	return *ret, err
}

// HeaderBlocks is a free data retrieval call binding the contract method 0x41539d4a.
//
// Solidity: function headerBlocks( uint256) constant returns(root bytes32, start uint256, end uint256, createdAt uint256, proposer address)
func (_Rootchain *RootchainSession) HeaderBlocks(arg0 *big.Int) (struct {
	Root      [32]byte
	Start     *big.Int
	End       *big.Int
	CreatedAt *big.Int
	Proposer  common.Address
}, error) {
	return _Rootchain.Contract.HeaderBlocks(&_Rootchain.CallOpts, arg0)
}

// HeaderBlocks is a free data retrieval call binding the contract method 0x41539d4a.
//
// Solidity: function headerBlocks( uint256) constant returns(root bytes32, start uint256, end uint256, createdAt uint256, proposer address)
func (_Rootchain *RootchainCallerSession) HeaderBlocks(arg0 *big.Int) (struct {
	Root      [32]byte
	Start     *big.Int
	End       *big.Int
	CreatedAt *big.Int
	Proposer  common.Address
}, error) {
	return _Rootchain.Contract.HeaderBlocks(&_Rootchain.CallOpts, arg0)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Rootchain *RootchainCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Rootchain *RootchainSession) IsOwner() (bool, error) {
	return _Rootchain.Contract.IsOwner(&_Rootchain.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Rootchain *RootchainCallerSession) IsOwner() (bool, error) {
	return _Rootchain.Contract.IsOwner(&_Rootchain.CallOpts)
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

// ProofValidatorContracts is a free data retrieval call binding the contract method 0xc4b875d3.
//
// Solidity: function proofValidatorContracts( address) constant returns(bool)
func (_Rootchain *RootchainCaller) ProofValidatorContracts(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "proofValidatorContracts", arg0)
	return *ret0, err
}

// ProofValidatorContracts is a free data retrieval call binding the contract method 0xc4b875d3.
//
// Solidity: function proofValidatorContracts( address) constant returns(bool)
func (_Rootchain *RootchainSession) ProofValidatorContracts(arg0 common.Address) (bool, error) {
	return _Rootchain.Contract.ProofValidatorContracts(&_Rootchain.CallOpts, arg0)
}

// ProofValidatorContracts is a free data retrieval call binding the contract method 0xc4b875d3.
//
// Solidity: function proofValidatorContracts( address) constant returns(bool)
func (_Rootchain *RootchainCallerSession) ProofValidatorContracts(arg0 common.Address) (bool, error) {
	return _Rootchain.Contract.ProofValidatorContracts(&_Rootchain.CallOpts, arg0)
}

// RoundType is a free data retrieval call binding the contract method 0x2c2d1a3b.
//
// Solidity: function roundType() constant returns(bytes32)
func (_Rootchain *RootchainCaller) RoundType(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "roundType")
	return *ret0, err
}

// RoundType is a free data retrieval call binding the contract method 0x2c2d1a3b.
//
// Solidity: function roundType() constant returns(bytes32)
func (_Rootchain *RootchainSession) RoundType() ([32]byte, error) {
	return _Rootchain.Contract.RoundType(&_Rootchain.CallOpts)
}

// RoundType is a free data retrieval call binding the contract method 0x2c2d1a3b.
//
// Solidity: function roundType() constant returns(bytes32)
func (_Rootchain *RootchainCallerSession) RoundType() ([32]byte, error) {
	return _Rootchain.Contract.RoundType(&_Rootchain.CallOpts)
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

// VoteType is a free data retrieval call binding the contract method 0x7d1a3d37.
//
// Solidity: function voteType() constant returns(bytes1)
func (_Rootchain *RootchainCaller) VoteType(opts *bind.CallOpts) ([1]byte, error) {
	var (
		ret0 = new([1]byte)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "voteType")
	return *ret0, err
}

// VoteType is a free data retrieval call binding the contract method 0x7d1a3d37.
//
// Solidity: function voteType() constant returns(bytes1)
func (_Rootchain *RootchainSession) VoteType() ([1]byte, error) {
	return _Rootchain.Contract.VoteType(&_Rootchain.CallOpts)
}

// VoteType is a free data retrieval call binding the contract method 0x7d1a3d37.
//
// Solidity: function voteType() constant returns(bytes1)
func (_Rootchain *RootchainCallerSession) VoteType() ([1]byte, error) {
	return _Rootchain.Contract.VoteType(&_Rootchain.CallOpts)
}

// WithdrawManager is a free data retrieval call binding the contract method 0xec3e9da5.
//
// Solidity: function withdrawManager() constant returns(address)
func (_Rootchain *RootchainCaller) WithdrawManager(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Rootchain.contract.Call(opts, out, "withdrawManager")
	return *ret0, err
}

// WithdrawManager is a free data retrieval call binding the contract method 0xec3e9da5.
//
// Solidity: function withdrawManager() constant returns(address)
func (_Rootchain *RootchainSession) WithdrawManager() (common.Address, error) {
	return _Rootchain.Contract.WithdrawManager(&_Rootchain.CallOpts)
}

// WithdrawManager is a free data retrieval call binding the contract method 0xec3e9da5.
//
// Solidity: function withdrawManager() constant returns(address)
func (_Rootchain *RootchainCallerSession) WithdrawManager() (common.Address, error) {
	return _Rootchain.Contract.WithdrawManager(&_Rootchain.CallOpts)
}

// AddProofValidator is a paid mutator transaction binding the contract method 0xd060828b.
//
// Solidity: function addProofValidator(_validator address) returns()
func (_Rootchain *RootchainTransactor) AddProofValidator(opts *bind.TransactOpts, _validator common.Address) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "addProofValidator", _validator)
}

// AddProofValidator is a paid mutator transaction binding the contract method 0xd060828b.
//
// Solidity: function addProofValidator(_validator address) returns()
func (_Rootchain *RootchainSession) AddProofValidator(_validator common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.AddProofValidator(&_Rootchain.TransactOpts, _validator)
}

// AddProofValidator is a paid mutator transaction binding the contract method 0xd060828b.
//
// Solidity: function addProofValidator(_validator address) returns()
func (_Rootchain *RootchainTransactorSession) AddProofValidator(_validator common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.AddProofValidator(&_Rootchain.TransactOpts, _validator)
}

// DeleteExit is a paid mutator transaction binding the contract method 0x50c30308.
//
// Solidity: function deleteExit(exitId uint256) returns()
func (_Rootchain *RootchainTransactor) DeleteExit(opts *bind.TransactOpts, exitId *big.Int) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "deleteExit", exitId)
}

// DeleteExit is a paid mutator transaction binding the contract method 0x50c30308.
//
// Solidity: function deleteExit(exitId uint256) returns()
func (_Rootchain *RootchainSession) DeleteExit(exitId *big.Int) (*types.Transaction, error) {
	return _Rootchain.Contract.DeleteExit(&_Rootchain.TransactOpts, exitId)
}

// DeleteExit is a paid mutator transaction binding the contract method 0x50c30308.
//
// Solidity: function deleteExit(exitId uint256) returns()
func (_Rootchain *RootchainTransactorSession) DeleteExit(exitId *big.Int) (*types.Transaction, error) {
	return _Rootchain.Contract.DeleteExit(&_Rootchain.TransactOpts, exitId)
}

// Deposit is a paid mutator transaction binding the contract method 0x8340f549.
//
// Solidity: function deposit(_token address, _user address, _amount uint256) returns()
func (_Rootchain *RootchainTransactor) Deposit(opts *bind.TransactOpts, _token common.Address, _user common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "deposit", _token, _user, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x8340f549.
//
// Solidity: function deposit(_token address, _user address, _amount uint256) returns()
func (_Rootchain *RootchainSession) Deposit(_token common.Address, _user common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Rootchain.Contract.Deposit(&_Rootchain.TransactOpts, _token, _user, _amount)
}

// Deposit is a paid mutator transaction binding the contract method 0x8340f549.
//
// Solidity: function deposit(_token address, _user address, _amount uint256) returns()
func (_Rootchain *RootchainTransactorSession) Deposit(_token common.Address, _user common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Rootchain.Contract.Deposit(&_Rootchain.TransactOpts, _token, _user, _amount)
}

// DepositERC721 is a paid mutator transaction binding the contract method 0x331ded1a.
//
// Solidity: function depositERC721(_token address, _user address, _tokenId uint256) returns()
func (_Rootchain *RootchainTransactor) DepositERC721(opts *bind.TransactOpts, _token common.Address, _user common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "depositERC721", _token, _user, _tokenId)
}

// DepositERC721 is a paid mutator transaction binding the contract method 0x331ded1a.
//
// Solidity: function depositERC721(_token address, _user address, _tokenId uint256) returns()
func (_Rootchain *RootchainSession) DepositERC721(_token common.Address, _user common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Rootchain.Contract.DepositERC721(&_Rootchain.TransactOpts, _token, _user, _tokenId)
}

// DepositERC721 is a paid mutator transaction binding the contract method 0x331ded1a.
//
// Solidity: function depositERC721(_token address, _user address, _tokenId uint256) returns()
func (_Rootchain *RootchainTransactorSession) DepositERC721(_token common.Address, _user common.Address, _tokenId *big.Int) (*types.Transaction, error) {
	return _Rootchain.Contract.DepositERC721(&_Rootchain.TransactOpts, _token, _user, _tokenId)
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

// FinalizeCommit is a paid mutator transaction binding the contract method 0xfb0df30f.
//
// Solidity: function finalizeCommit( uint256) returns()
func (_Rootchain *RootchainTransactor) FinalizeCommit(opts *bind.TransactOpts, arg0 *big.Int) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "finalizeCommit", arg0)
}

// FinalizeCommit is a paid mutator transaction binding the contract method 0xfb0df30f.
//
// Solidity: function finalizeCommit( uint256) returns()
func (_Rootchain *RootchainSession) FinalizeCommit(arg0 *big.Int) (*types.Transaction, error) {
	return _Rootchain.Contract.FinalizeCommit(&_Rootchain.TransactOpts, arg0)
}

// FinalizeCommit is a paid mutator transaction binding the contract method 0xfb0df30f.
//
// Solidity: function finalizeCommit( uint256) returns()
func (_Rootchain *RootchainTransactorSession) FinalizeCommit(arg0 *big.Int) (*types.Transaction, error) {
	return _Rootchain.Contract.FinalizeCommit(&_Rootchain.TransactOpts, arg0)
}

// MapToken is a paid mutator transaction binding the contract method 0xe117694b.
//
// Solidity: function mapToken(_rootToken address, _childToken address, _isERC721 bool) returns()
func (_Rootchain *RootchainTransactor) MapToken(opts *bind.TransactOpts, _rootToken common.Address, _childToken common.Address, _isERC721 bool) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "mapToken", _rootToken, _childToken, _isERC721)
}

// MapToken is a paid mutator transaction binding the contract method 0xe117694b.
//
// Solidity: function mapToken(_rootToken address, _childToken address, _isERC721 bool) returns()
func (_Rootchain *RootchainSession) MapToken(_rootToken common.Address, _childToken common.Address, _isERC721 bool) (*types.Transaction, error) {
	return _Rootchain.Contract.MapToken(&_Rootchain.TransactOpts, _rootToken, _childToken, _isERC721)
}

// MapToken is a paid mutator transaction binding the contract method 0xe117694b.
//
// Solidity: function mapToken(_rootToken address, _childToken address, _isERC721 bool) returns()
func (_Rootchain *RootchainTransactorSession) MapToken(_rootToken common.Address, _childToken common.Address, _isERC721 bool) (*types.Transaction, error) {
	return _Rootchain.Contract.MapToken(&_Rootchain.TransactOpts, _rootToken, _childToken, _isERC721)
}

// RemoveProofValidator is a paid mutator transaction binding the contract method 0x609dc55a.
//
// Solidity: function removeProofValidator(_validator address) returns()
func (_Rootchain *RootchainTransactor) RemoveProofValidator(opts *bind.TransactOpts, _validator common.Address) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "removeProofValidator", _validator)
}

// RemoveProofValidator is a paid mutator transaction binding the contract method 0x609dc55a.
//
// Solidity: function removeProofValidator(_validator address) returns()
func (_Rootchain *RootchainSession) RemoveProofValidator(_validator common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.RemoveProofValidator(&_Rootchain.TransactOpts, _validator)
}

// RemoveProofValidator is a paid mutator transaction binding the contract method 0x609dc55a.
//
// Solidity: function removeProofValidator(_validator address) returns()
func (_Rootchain *RootchainTransactorSession) RemoveProofValidator(_validator common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.RemoveProofValidator(&_Rootchain.TransactOpts, _validator)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rootchain *RootchainTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rootchain *RootchainSession) RenounceOwnership() (*types.Transaction, error) {
	return _Rootchain.Contract.RenounceOwnership(&_Rootchain.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Rootchain *RootchainTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Rootchain.Contract.RenounceOwnership(&_Rootchain.TransactOpts)
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

// SetDepositManager is a paid mutator transaction binding the contract method 0x228d71a9.
//
// Solidity: function setDepositManager(_depositManager address) returns()
func (_Rootchain *RootchainTransactor) SetDepositManager(opts *bind.TransactOpts, _depositManager common.Address) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "setDepositManager", _depositManager)
}

// SetDepositManager is a paid mutator transaction binding the contract method 0x228d71a9.
//
// Solidity: function setDepositManager(_depositManager address) returns()
func (_Rootchain *RootchainSession) SetDepositManager(_depositManager common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.SetDepositManager(&_Rootchain.TransactOpts, _depositManager)
}

// SetDepositManager is a paid mutator transaction binding the contract method 0x228d71a9.
//
// Solidity: function setDepositManager(_depositManager address) returns()
func (_Rootchain *RootchainTransactorSession) SetDepositManager(_depositManager common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.SetDepositManager(&_Rootchain.TransactOpts, _depositManager)
}

// SetExitNFTContract is a paid mutator transaction binding the contract method 0x46e11a8d.
//
// Solidity: function setExitNFTContract(_nftContract address) returns()
func (_Rootchain *RootchainTransactor) SetExitNFTContract(opts *bind.TransactOpts, _nftContract common.Address) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "setExitNFTContract", _nftContract)
}

// SetExitNFTContract is a paid mutator transaction binding the contract method 0x46e11a8d.
//
// Solidity: function setExitNFTContract(_nftContract address) returns()
func (_Rootchain *RootchainSession) SetExitNFTContract(_nftContract common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.SetExitNFTContract(&_Rootchain.TransactOpts, _nftContract)
}

// SetExitNFTContract is a paid mutator transaction binding the contract method 0x46e11a8d.
//
// Solidity: function setExitNFTContract(_nftContract address) returns()
func (_Rootchain *RootchainTransactorSession) SetExitNFTContract(_nftContract common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.SetExitNFTContract(&_Rootchain.TransactOpts, _nftContract)
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

// SetWithdrawManager is a paid mutator transaction binding the contract method 0x17e3e2e8.
//
// Solidity: function setWithdrawManager(_withdrawManager address) returns()
func (_Rootchain *RootchainTransactor) SetWithdrawManager(opts *bind.TransactOpts, _withdrawManager common.Address) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "setWithdrawManager", _withdrawManager)
}

// SetWithdrawManager is a paid mutator transaction binding the contract method 0x17e3e2e8.
//
// Solidity: function setWithdrawManager(_withdrawManager address) returns()
func (_Rootchain *RootchainSession) SetWithdrawManager(_withdrawManager common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.SetWithdrawManager(&_Rootchain.TransactOpts, _withdrawManager)
}

// SetWithdrawManager is a paid mutator transaction binding the contract method 0x17e3e2e8.
//
// Solidity: function setWithdrawManager(_withdrawManager address) returns()
func (_Rootchain *RootchainTransactorSession) SetWithdrawManager(_withdrawManager common.Address) (*types.Transaction, error) {
	return _Rootchain.Contract.SetWithdrawManager(&_Rootchain.TransactOpts, _withdrawManager)
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

// SubmitHeaderBlock is a paid mutator transaction binding the contract method 0xec83d3ba.
//
// Solidity: function submitHeaderBlock(vote bytes, sigs bytes, extradata bytes) returns()
func (_Rootchain *RootchainTransactor) SubmitHeaderBlock(opts *bind.TransactOpts, vote []byte, sigs []byte, extradata []byte) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "submitHeaderBlock", vote, sigs, extradata)
}

// SubmitHeaderBlock is a paid mutator transaction binding the contract method 0xec83d3ba.
//
// Solidity: function submitHeaderBlock(vote bytes, sigs bytes, extradata bytes) returns()
func (_Rootchain *RootchainSession) SubmitHeaderBlock(vote []byte, sigs []byte, extradata []byte) (*types.Transaction, error) {
	return _Rootchain.Contract.SubmitHeaderBlock(&_Rootchain.TransactOpts, vote, sigs, extradata)
}

// SubmitHeaderBlock is a paid mutator transaction binding the contract method 0xec83d3ba.
//
// Solidity: function submitHeaderBlock(vote bytes, sigs bytes, extradata bytes) returns()
func (_Rootchain *RootchainTransactorSession) SubmitHeaderBlock(vote []byte, sigs []byte, extradata []byte) (*types.Transaction, error) {
	return _Rootchain.Contract.SubmitHeaderBlock(&_Rootchain.TransactOpts, vote, sigs, extradata)
}

// TokenFallback is a paid mutator transaction binding the contract method 0xc0ee0b8a.
//
// Solidity: function tokenFallback(_user address, _amount uint256, _data bytes) returns()
func (_Rootchain *RootchainTransactor) TokenFallback(opts *bind.TransactOpts, _user common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "tokenFallback", _user, _amount, _data)
}

// TokenFallback is a paid mutator transaction binding the contract method 0xc0ee0b8a.
//
// Solidity: function tokenFallback(_user address, _amount uint256, _data bytes) returns()
func (_Rootchain *RootchainSession) TokenFallback(_user common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Rootchain.Contract.TokenFallback(&_Rootchain.TransactOpts, _user, _amount, _data)
}

// TokenFallback is a paid mutator transaction binding the contract method 0xc0ee0b8a.
//
// Solidity: function tokenFallback(_user address, _amount uint256, _data bytes) returns()
func (_Rootchain *RootchainTransactorSession) TokenFallback(_user common.Address, _amount *big.Int, _data []byte) (*types.Transaction, error) {
	return _Rootchain.Contract.TokenFallback(&_Rootchain.TransactOpts, _user, _amount, _data)
}

// TransferAmount is a paid mutator transaction binding the contract method 0x01f47471.
//
// Solidity: function transferAmount(_token address, _user address, _amount uint256) returns(bool)
func (_Rootchain *RootchainTransactor) TransferAmount(opts *bind.TransactOpts, _token common.Address, _user common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Rootchain.contract.Transact(opts, "transferAmount", _token, _user, _amount)
}

// TransferAmount is a paid mutator transaction binding the contract method 0x01f47471.
//
// Solidity: function transferAmount(_token address, _user address, _amount uint256) returns(bool)
func (_Rootchain *RootchainSession) TransferAmount(_token common.Address, _user common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Rootchain.Contract.TransferAmount(&_Rootchain.TransactOpts, _token, _user, _amount)
}

// TransferAmount is a paid mutator transaction binding the contract method 0x01f47471.
//
// Solidity: function transferAmount(_token address, _user address, _amount uint256) returns(bool)
func (_Rootchain *RootchainTransactorSession) TransferAmount(_token common.Address, _user common.Address, _amount *big.Int) (*types.Transaction, error) {
	return _Rootchain.Contract.TransferAmount(&_Rootchain.TransactOpts, _token, _user, _amount)
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
// Solidity: e NewHeaderBlock(proposer indexed address, number indexed uint256, start uint256, end uint256, root bytes32)
func (_Rootchain *RootchainFilterer) FilterNewHeaderBlock(opts *bind.FilterOpts, proposer []common.Address, number []*big.Int) (*RootchainNewHeaderBlockIterator, error) {

	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}
	var numberRule []interface{}
	for _, numberItem := range number {
		numberRule = append(numberRule, numberItem)
	}

	logs, sub, err := _Rootchain.contract.FilterLogs(opts, "NewHeaderBlock", proposerRule, numberRule)
	if err != nil {
		return nil, err
	}
	return &RootchainNewHeaderBlockIterator{contract: _Rootchain.contract, event: "NewHeaderBlock", logs: logs, sub: sub}, nil
}

// WatchNewHeaderBlock is a free log subscription operation binding the contract event 0xf146921b854b787ba7d6045e8a8054731dc62430ae16c4bf08147539b1b6ef8f.
//
// Solidity: e NewHeaderBlock(proposer indexed address, number indexed uint256, start uint256, end uint256, root bytes32)
func (_Rootchain *RootchainFilterer) WatchNewHeaderBlock(opts *bind.WatchOpts, sink chan<- *RootchainNewHeaderBlock, proposer []common.Address, number []*big.Int) (event.Subscription, error) {

	var proposerRule []interface{}
	for _, proposerItem := range proposer {
		proposerRule = append(proposerRule, proposerItem)
	}
	var numberRule []interface{}
	for _, numberItem := range number {
		numberRule = append(numberRule, numberItem)
	}

	logs, sub, err := _Rootchain.contract.WatchLogs(opts, "NewHeaderBlock", proposerRule, numberRule)
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

// RootchainProofValidatorAddedIterator is returned from FilterProofValidatorAdded and is used to iterate over the raw logs and unpacked data for ProofValidatorAdded events raised by the Rootchain contract.
type RootchainProofValidatorAddedIterator struct {
	Event *RootchainProofValidatorAdded // Event containing the contract specifics and raw log

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
func (it *RootchainProofValidatorAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootchainProofValidatorAdded)
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
		it.Event = new(RootchainProofValidatorAdded)
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
func (it *RootchainProofValidatorAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootchainProofValidatorAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootchainProofValidatorAdded represents a ProofValidatorAdded event raised by the Rootchain contract.
type RootchainProofValidatorAdded struct {
	Validator common.Address
	From      common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterProofValidatorAdded is a free log retrieval operation binding the contract event 0x3dc12d30280bcd33917d2b84141129635923441ba7e6b388b946b41f5ace697d.
//
// Solidity: e ProofValidatorAdded(validator indexed address, from indexed address)
func (_Rootchain *RootchainFilterer) FilterProofValidatorAdded(opts *bind.FilterOpts, validator []common.Address, from []common.Address) (*RootchainProofValidatorAddedIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Rootchain.contract.FilterLogs(opts, "ProofValidatorAdded", validatorRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &RootchainProofValidatorAddedIterator{contract: _Rootchain.contract, event: "ProofValidatorAdded", logs: logs, sub: sub}, nil
}

// WatchProofValidatorAdded is a free log subscription operation binding the contract event 0x3dc12d30280bcd33917d2b84141129635923441ba7e6b388b946b41f5ace697d.
//
// Solidity: e ProofValidatorAdded(validator indexed address, from indexed address)
func (_Rootchain *RootchainFilterer) WatchProofValidatorAdded(opts *bind.WatchOpts, sink chan<- *RootchainProofValidatorAdded, validator []common.Address, from []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Rootchain.contract.WatchLogs(opts, "ProofValidatorAdded", validatorRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootchainProofValidatorAdded)
				if err := _Rootchain.contract.UnpackLog(event, "ProofValidatorAdded", log); err != nil {
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

// RootchainProofValidatorRemovedIterator is returned from FilterProofValidatorRemoved and is used to iterate over the raw logs and unpacked data for ProofValidatorRemoved events raised by the Rootchain contract.
type RootchainProofValidatorRemovedIterator struct {
	Event *RootchainProofValidatorRemoved // Event containing the contract specifics and raw log

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
func (it *RootchainProofValidatorRemovedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(RootchainProofValidatorRemoved)
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
		it.Event = new(RootchainProofValidatorRemoved)
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
func (it *RootchainProofValidatorRemovedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *RootchainProofValidatorRemovedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// RootchainProofValidatorRemoved represents a ProofValidatorRemoved event raised by the Rootchain contract.
type RootchainProofValidatorRemoved struct {
	Validator common.Address
	From      common.Address
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterProofValidatorRemoved is a free log retrieval operation binding the contract event 0x96bedef125d36a85bf369db1f6ac9d7487d9daf6d4c22539249f1bf94a11e119.
//
// Solidity: e ProofValidatorRemoved(validator indexed address, from indexed address)
func (_Rootchain *RootchainFilterer) FilterProofValidatorRemoved(opts *bind.FilterOpts, validator []common.Address, from []common.Address) (*RootchainProofValidatorRemovedIterator, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Rootchain.contract.FilterLogs(opts, "ProofValidatorRemoved", validatorRule, fromRule)
	if err != nil {
		return nil, err
	}
	return &RootchainProofValidatorRemovedIterator{contract: _Rootchain.contract, event: "ProofValidatorRemoved", logs: logs, sub: sub}, nil
}

// WatchProofValidatorRemoved is a free log subscription operation binding the contract event 0x96bedef125d36a85bf369db1f6ac9d7487d9daf6d4c22539249f1bf94a11e119.
//
// Solidity: e ProofValidatorRemoved(validator indexed address, from indexed address)
func (_Rootchain *RootchainFilterer) WatchProofValidatorRemoved(opts *bind.WatchOpts, sink chan<- *RootchainProofValidatorRemoved, validator []common.Address, from []common.Address) (event.Subscription, error) {

	var validatorRule []interface{}
	for _, validatorItem := range validator {
		validatorRule = append(validatorRule, validatorItem)
	}
	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _Rootchain.contract.WatchLogs(opts, "ProofValidatorRemoved", validatorRule, fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(RootchainProofValidatorRemoved)
				if err := _Rootchain.contract.UnpackLog(event, "ProofValidatorRemoved", log); err != nil {
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
