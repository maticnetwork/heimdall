// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package depositmanager

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

// DepositmanagerABI is the input ABI used to generate the binding from.
const DepositmanagerABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"currentHeaderBlock\",\"type\":\"uint256\"}],\"name\":\"nextDepositBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"roundType\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"depositCount\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_depositCount\",\"type\":\"uint256\"}],\"name\":\"depositBlock\",\"outputs\":[{\"name\":\"_header\",\"type\":\"uint256\"},{\"name\":\"_owner\",\"type\":\"address\"},{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_amountOrTokenId\",\"type\":\"uint256\"},{\"name\":\"_createdAt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"reverseTokens\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_nftContract\",\"type\":\"address\"}],\"name\":\"setExitNFTContract\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_currentHeaderBlock\",\"type\":\"uint256\"},{\"name\":\"_token\",\"type\":\"address\"},{\"name\":\"_user\",\"type\":\"address\"},{\"name\":\"_amountOrTokenId\",\"type\":\"uint256\"}],\"name\":\"createDepositBlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"wethToken\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"voteType\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes1\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"networkId\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"rootChain\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"CHILD_BLOCK_INTERVAL\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"deposits\",\"outputs\":[{\"name\":\"header\",\"type\":\"uint256\"},{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"token\",\"type\":\"address\"},{\"name\":\"amountOrTokenId\",\"type\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_token\",\"type\":\"address\"}],\"name\":\"setWETHToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chain\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"isERC721\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_rootToken\",\"type\":\"address\"},{\"name\":\"_childToken\",\"type\":\"address\"},{\"name\":\"_isERC721\",\"type\":\"bool\"}],\"name\":\"mapToken\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"tokens\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newRootChain\",\"type\":\"address\"}],\"name\":\"changeRootChain\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_currentHeaderBlock\",\"type\":\"uint256\"}],\"name\":\"finalizeCommit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_token\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"_amountOrTokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"_depositCount\",\"type\":\"uint256\"}],\"name\":\"Deposit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousRootChain\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newRootChain\",\"type\":\"address\"}],\"name\":\"RootChainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"_rootToken\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"_childToken\",\"type\":\"address\"}],\"name\":\"TokenMapped\",\"type\":\"event\"}]"

// Depositmanager is an auto generated Go binding around an Ethereum contract.
type Depositmanager struct {
	DepositmanagerCaller     // Read-only binding to the contract
	DepositmanagerTransactor // Write-only binding to the contract
	DepositmanagerFilterer   // Log filterer for contract events
}

// DepositmanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type DepositmanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositmanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DepositmanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositmanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DepositmanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DepositmanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DepositmanagerSession struct {
	Contract     *Depositmanager   // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DepositmanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DepositmanagerCallerSession struct {
	Contract *DepositmanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts         // Call options to use throughout this session
}

// DepositmanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DepositmanagerTransactorSession struct {
	Contract     *DepositmanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts         // Transaction auth options to use throughout this session
}

// DepositmanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type DepositmanagerRaw struct {
	Contract *Depositmanager // Generic contract binding to access the raw methods on
}

// DepositmanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DepositmanagerCallerRaw struct {
	Contract *DepositmanagerCaller // Generic read-only contract binding to access the raw methods on
}

// DepositmanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DepositmanagerTransactorRaw struct {
	Contract *DepositmanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDepositmanager creates a new instance of Depositmanager, bound to a specific deployed contract.
func NewDepositmanager(address common.Address, backend bind.ContractBackend) (*Depositmanager, error) {
	contract, err := bindDepositmanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Depositmanager{DepositmanagerCaller: DepositmanagerCaller{contract: contract}, DepositmanagerTransactor: DepositmanagerTransactor{contract: contract}, DepositmanagerFilterer: DepositmanagerFilterer{contract: contract}}, nil
}

// NewDepositmanagerCaller creates a new read-only instance of Depositmanager, bound to a specific deployed contract.
func NewDepositmanagerCaller(address common.Address, caller bind.ContractCaller) (*DepositmanagerCaller, error) {
	contract, err := bindDepositmanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerCaller{contract: contract}, nil
}

// NewDepositmanagerTransactor creates a new write-only instance of Depositmanager, bound to a specific deployed contract.
func NewDepositmanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*DepositmanagerTransactor, error) {
	contract, err := bindDepositmanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerTransactor{contract: contract}, nil
}

// NewDepositmanagerFilterer creates a new log filterer instance of Depositmanager, bound to a specific deployed contract.
func NewDepositmanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*DepositmanagerFilterer, error) {
	contract, err := bindDepositmanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerFilterer{contract: contract}, nil
}

// bindDepositmanager binds a generic wrapper to an already deployed contract.
func bindDepositmanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DepositmanagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Depositmanager *DepositmanagerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Depositmanager.Contract.DepositmanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Depositmanager *DepositmanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Depositmanager.Contract.DepositmanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Depositmanager *DepositmanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Depositmanager.Contract.DepositmanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Depositmanager *DepositmanagerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Depositmanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Depositmanager *DepositmanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Depositmanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Depositmanager *DepositmanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Depositmanager.Contract.contract.Transact(opts, method, params...)
}

// CHILDBLOCKINTERVAL is a free data retrieval call binding the contract method 0xa831fa07.
//
// Solidity: function CHILD_BLOCK_INTERVAL() constant returns(uint256)
func (_Depositmanager *DepositmanagerCaller) CHILDBLOCKINTERVAL(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Depositmanager.contract.Call(opts, out, "CHILD_BLOCK_INTERVAL")
	return *ret0, err
}

// CHILDBLOCKINTERVAL is a free data retrieval call binding the contract method 0xa831fa07.
//
// Solidity: function CHILD_BLOCK_INTERVAL() constant returns(uint256)
func (_Depositmanager *DepositmanagerSession) CHILDBLOCKINTERVAL() (*big.Int, error) {
	return _Depositmanager.Contract.CHILDBLOCKINTERVAL(&_Depositmanager.CallOpts)
}

// CHILDBLOCKINTERVAL is a free data retrieval call binding the contract method 0xa831fa07.
//
// Solidity: function CHILD_BLOCK_INTERVAL() constant returns(uint256)
func (_Depositmanager *DepositmanagerCallerSession) CHILDBLOCKINTERVAL() (*big.Int, error) {
	return _Depositmanager.Contract.CHILDBLOCKINTERVAL(&_Depositmanager.CallOpts)
}

// Chain is a free data retrieval call binding the contract method 0xc763e5a1.
//
// Solidity: function chain() constant returns(bytes32)
func (_Depositmanager *DepositmanagerCaller) Chain(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Depositmanager.contract.Call(opts, out, "chain")
	return *ret0, err
}

// Chain is a free data retrieval call binding the contract method 0xc763e5a1.
//
// Solidity: function chain() constant returns(bytes32)
func (_Depositmanager *DepositmanagerSession) Chain() ([32]byte, error) {
	return _Depositmanager.Contract.Chain(&_Depositmanager.CallOpts)
}

// Chain is a free data retrieval call binding the contract method 0xc763e5a1.
//
// Solidity: function chain() constant returns(bytes32)
func (_Depositmanager *DepositmanagerCallerSession) Chain() ([32]byte, error) {
	return _Depositmanager.Contract.Chain(&_Depositmanager.CallOpts)
}

// DepositBlock is a free data retrieval call binding the contract method 0x32590654.
//
// Solidity: function depositBlock(uint256 _depositCount) constant returns(uint256 _header, address _owner, address _token, uint256 _amountOrTokenId, uint256 _createdAt)
func (_Depositmanager *DepositmanagerCaller) DepositBlock(opts *bind.CallOpts, _depositCount *big.Int) (struct {
	Header          *big.Int
	Owner           common.Address
	Token           common.Address
	AmountOrTokenId *big.Int
	CreatedAt       *big.Int
}, error) {
	ret := new(struct {
		Header          *big.Int
		Owner           common.Address
		Token           common.Address
		AmountOrTokenId *big.Int
		CreatedAt       *big.Int
	})
	out := ret
	err := _Depositmanager.contract.Call(opts, out, "depositBlock", _depositCount)
	return *ret, err
}

// DepositBlock is a free data retrieval call binding the contract method 0x32590654.
//
// Solidity: function depositBlock(uint256 _depositCount) constant returns(uint256 _header, address _owner, address _token, uint256 _amountOrTokenId, uint256 _createdAt)
func (_Depositmanager *DepositmanagerSession) DepositBlock(_depositCount *big.Int) (struct {
	Header          *big.Int
	Owner           common.Address
	Token           common.Address
	AmountOrTokenId *big.Int
	CreatedAt       *big.Int
}, error) {
	return _Depositmanager.Contract.DepositBlock(&_Depositmanager.CallOpts, _depositCount)
}

// DepositBlock is a free data retrieval call binding the contract method 0x32590654.
//
// Solidity: function depositBlock(uint256 _depositCount) constant returns(uint256 _header, address _owner, address _token, uint256 _amountOrTokenId, uint256 _createdAt)
func (_Depositmanager *DepositmanagerCallerSession) DepositBlock(_depositCount *big.Int) (struct {
	Header          *big.Int
	Owner           common.Address
	Token           common.Address
	AmountOrTokenId *big.Int
	CreatedAt       *big.Int
}, error) {
	return _Depositmanager.Contract.DepositBlock(&_Depositmanager.CallOpts, _depositCount)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() constant returns(uint256)
func (_Depositmanager *DepositmanagerCaller) DepositCount(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Depositmanager.contract.Call(opts, out, "depositCount")
	return *ret0, err
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() constant returns(uint256)
func (_Depositmanager *DepositmanagerSession) DepositCount() (*big.Int, error) {
	return _Depositmanager.Contract.DepositCount(&_Depositmanager.CallOpts)
}

// DepositCount is a free data retrieval call binding the contract method 0x2dfdf0b5.
//
// Solidity: function depositCount() constant returns(uint256)
func (_Depositmanager *DepositmanagerCallerSession) DepositCount() (*big.Int, error) {
	return _Depositmanager.Contract.DepositCount(&_Depositmanager.CallOpts)
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) constant returns(uint256 header, address owner, address token, uint256 amountOrTokenId, uint256 createdAt)
func (_Depositmanager *DepositmanagerCaller) Deposits(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Header          *big.Int
	Owner           common.Address
	Token           common.Address
	AmountOrTokenId *big.Int
	CreatedAt       *big.Int
}, error) {
	ret := new(struct {
		Header          *big.Int
		Owner           common.Address
		Token           common.Address
		AmountOrTokenId *big.Int
		CreatedAt       *big.Int
	})
	out := ret
	err := _Depositmanager.contract.Call(opts, out, "deposits", arg0)
	return *ret, err
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) constant returns(uint256 header, address owner, address token, uint256 amountOrTokenId, uint256 createdAt)
func (_Depositmanager *DepositmanagerSession) Deposits(arg0 *big.Int) (struct {
	Header          *big.Int
	Owner           common.Address
	Token           common.Address
	AmountOrTokenId *big.Int
	CreatedAt       *big.Int
}, error) {
	return _Depositmanager.Contract.Deposits(&_Depositmanager.CallOpts, arg0)
}

// Deposits is a free data retrieval call binding the contract method 0xb02c43d0.
//
// Solidity: function deposits(uint256 ) constant returns(uint256 header, address owner, address token, uint256 amountOrTokenId, uint256 createdAt)
func (_Depositmanager *DepositmanagerCallerSession) Deposits(arg0 *big.Int) (struct {
	Header          *big.Int
	Owner           common.Address
	Token           common.Address
	AmountOrTokenId *big.Int
	CreatedAt       *big.Int
}, error) {
	return _Depositmanager.Contract.Deposits(&_Depositmanager.CallOpts, arg0)
}

// IsERC721 is a free data retrieval call binding the contract method 0xdaa09e54.
//
// Solidity: function isERC721(address ) constant returns(bool)
func (_Depositmanager *DepositmanagerCaller) IsERC721(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Depositmanager.contract.Call(opts, out, "isERC721", arg0)
	return *ret0, err
}

// IsERC721 is a free data retrieval call binding the contract method 0xdaa09e54.
//
// Solidity: function isERC721(address ) constant returns(bool)
func (_Depositmanager *DepositmanagerSession) IsERC721(arg0 common.Address) (bool, error) {
	return _Depositmanager.Contract.IsERC721(&_Depositmanager.CallOpts, arg0)
}

// IsERC721 is a free data retrieval call binding the contract method 0xdaa09e54.
//
// Solidity: function isERC721(address ) constant returns(bool)
func (_Depositmanager *DepositmanagerCallerSession) IsERC721(arg0 common.Address) (bool, error) {
	return _Depositmanager.Contract.IsERC721(&_Depositmanager.CallOpts, arg0)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Depositmanager *DepositmanagerCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Depositmanager.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Depositmanager *DepositmanagerSession) IsOwner() (bool, error) {
	return _Depositmanager.Contract.IsOwner(&_Depositmanager.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Depositmanager *DepositmanagerCallerSession) IsOwner() (bool, error) {
	return _Depositmanager.Contract.IsOwner(&_Depositmanager.CallOpts)
}

// NetworkId is a free data retrieval call binding the contract method 0x9025e64c.
//
// Solidity: function networkId() constant returns(bytes)
func (_Depositmanager *DepositmanagerCaller) NetworkId(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _Depositmanager.contract.Call(opts, out, "networkId")
	return *ret0, err
}

// NetworkId is a free data retrieval call binding the contract method 0x9025e64c.
//
// Solidity: function networkId() constant returns(bytes)
func (_Depositmanager *DepositmanagerSession) NetworkId() ([]byte, error) {
	return _Depositmanager.Contract.NetworkId(&_Depositmanager.CallOpts)
}

// NetworkId is a free data retrieval call binding the contract method 0x9025e64c.
//
// Solidity: function networkId() constant returns(bytes)
func (_Depositmanager *DepositmanagerCallerSession) NetworkId() ([]byte, error) {
	return _Depositmanager.Contract.NetworkId(&_Depositmanager.CallOpts)
}

// NextDepositBlock is a free data retrieval call binding the contract method 0x0472bc3c.
//
// Solidity: function nextDepositBlock(uint256 currentHeaderBlock) constant returns(uint256)
func (_Depositmanager *DepositmanagerCaller) NextDepositBlock(opts *bind.CallOpts, currentHeaderBlock *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Depositmanager.contract.Call(opts, out, "nextDepositBlock", currentHeaderBlock)
	return *ret0, err
}

// NextDepositBlock is a free data retrieval call binding the contract method 0x0472bc3c.
//
// Solidity: function nextDepositBlock(uint256 currentHeaderBlock) constant returns(uint256)
func (_Depositmanager *DepositmanagerSession) NextDepositBlock(currentHeaderBlock *big.Int) (*big.Int, error) {
	return _Depositmanager.Contract.NextDepositBlock(&_Depositmanager.CallOpts, currentHeaderBlock)
}

// NextDepositBlock is a free data retrieval call binding the contract method 0x0472bc3c.
//
// Solidity: function nextDepositBlock(uint256 currentHeaderBlock) constant returns(uint256)
func (_Depositmanager *DepositmanagerCallerSession) NextDepositBlock(currentHeaderBlock *big.Int) (*big.Int, error) {
	return _Depositmanager.Contract.NextDepositBlock(&_Depositmanager.CallOpts, currentHeaderBlock)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Depositmanager *DepositmanagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Depositmanager.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Depositmanager *DepositmanagerSession) Owner() (common.Address, error) {
	return _Depositmanager.Contract.Owner(&_Depositmanager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Depositmanager *DepositmanagerCallerSession) Owner() (common.Address, error) {
	return _Depositmanager.Contract.Owner(&_Depositmanager.CallOpts)
}

// ReverseTokens is a free data retrieval call binding the contract method 0x40828ebf.
//
// Solidity: function reverseTokens(address ) constant returns(address)
func (_Depositmanager *DepositmanagerCaller) ReverseTokens(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Depositmanager.contract.Call(opts, out, "reverseTokens", arg0)
	return *ret0, err
}

// ReverseTokens is a free data retrieval call binding the contract method 0x40828ebf.
//
// Solidity: function reverseTokens(address ) constant returns(address)
func (_Depositmanager *DepositmanagerSession) ReverseTokens(arg0 common.Address) (common.Address, error) {
	return _Depositmanager.Contract.ReverseTokens(&_Depositmanager.CallOpts, arg0)
}

// ReverseTokens is a free data retrieval call binding the contract method 0x40828ebf.
//
// Solidity: function reverseTokens(address ) constant returns(address)
func (_Depositmanager *DepositmanagerCallerSession) ReverseTokens(arg0 common.Address) (common.Address, error) {
	return _Depositmanager.Contract.ReverseTokens(&_Depositmanager.CallOpts, arg0)
}

// RootChain is a free data retrieval call binding the contract method 0x987ab9db.
//
// Solidity: function rootChain() constant returns(address)
func (_Depositmanager *DepositmanagerCaller) RootChain(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Depositmanager.contract.Call(opts, out, "rootChain")
	return *ret0, err
}

// RootChain is a free data retrieval call binding the contract method 0x987ab9db.
//
// Solidity: function rootChain() constant returns(address)
func (_Depositmanager *DepositmanagerSession) RootChain() (common.Address, error) {
	return _Depositmanager.Contract.RootChain(&_Depositmanager.CallOpts)
}

// RootChain is a free data retrieval call binding the contract method 0x987ab9db.
//
// Solidity: function rootChain() constant returns(address)
func (_Depositmanager *DepositmanagerCallerSession) RootChain() (common.Address, error) {
	return _Depositmanager.Contract.RootChain(&_Depositmanager.CallOpts)
}

// RoundType is a free data retrieval call binding the contract method 0x2c2d1a3b.
//
// Solidity: function roundType() constant returns(bytes32)
func (_Depositmanager *DepositmanagerCaller) RoundType(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Depositmanager.contract.Call(opts, out, "roundType")
	return *ret0, err
}

// RoundType is a free data retrieval call binding the contract method 0x2c2d1a3b.
//
// Solidity: function roundType() constant returns(bytes32)
func (_Depositmanager *DepositmanagerSession) RoundType() ([32]byte, error) {
	return _Depositmanager.Contract.RoundType(&_Depositmanager.CallOpts)
}

// RoundType is a free data retrieval call binding the contract method 0x2c2d1a3b.
//
// Solidity: function roundType() constant returns(bytes32)
func (_Depositmanager *DepositmanagerCallerSession) RoundType() ([32]byte, error) {
	return _Depositmanager.Contract.RoundType(&_Depositmanager.CallOpts)
}

// Tokens is a free data retrieval call binding the contract method 0xe4860339.
//
// Solidity: function tokens(address ) constant returns(address)
func (_Depositmanager *DepositmanagerCaller) Tokens(opts *bind.CallOpts, arg0 common.Address) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Depositmanager.contract.Call(opts, out, "tokens", arg0)
	return *ret0, err
}

// Tokens is a free data retrieval call binding the contract method 0xe4860339.
//
// Solidity: function tokens(address ) constant returns(address)
func (_Depositmanager *DepositmanagerSession) Tokens(arg0 common.Address) (common.Address, error) {
	return _Depositmanager.Contract.Tokens(&_Depositmanager.CallOpts, arg0)
}

// Tokens is a free data retrieval call binding the contract method 0xe4860339.
//
// Solidity: function tokens(address ) constant returns(address)
func (_Depositmanager *DepositmanagerCallerSession) Tokens(arg0 common.Address) (common.Address, error) {
	return _Depositmanager.Contract.Tokens(&_Depositmanager.CallOpts, arg0)
}

// VoteType is a free data retrieval call binding the contract method 0x7d1a3d37.
//
// Solidity: function voteType() constant returns(bytes1)
func (_Depositmanager *DepositmanagerCaller) VoteType(opts *bind.CallOpts) ([1]byte, error) {
	var (
		ret0 = new([1]byte)
	)
	out := ret0
	err := _Depositmanager.contract.Call(opts, out, "voteType")
	return *ret0, err
}

// VoteType is a free data retrieval call binding the contract method 0x7d1a3d37.
//
// Solidity: function voteType() constant returns(bytes1)
func (_Depositmanager *DepositmanagerSession) VoteType() ([1]byte, error) {
	return _Depositmanager.Contract.VoteType(&_Depositmanager.CallOpts)
}

// VoteType is a free data retrieval call binding the contract method 0x7d1a3d37.
//
// Solidity: function voteType() constant returns(bytes1)
func (_Depositmanager *DepositmanagerCallerSession) VoteType() ([1]byte, error) {
	return _Depositmanager.Contract.VoteType(&_Depositmanager.CallOpts)
}

// WethToken is a free data retrieval call binding the contract method 0x4b57b0be.
//
// Solidity: function wethToken() constant returns(address)
func (_Depositmanager *DepositmanagerCaller) WethToken(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Depositmanager.contract.Call(opts, out, "wethToken")
	return *ret0, err
}

// WethToken is a free data retrieval call binding the contract method 0x4b57b0be.
//
// Solidity: function wethToken() constant returns(address)
func (_Depositmanager *DepositmanagerSession) WethToken() (common.Address, error) {
	return _Depositmanager.Contract.WethToken(&_Depositmanager.CallOpts)
}

// WethToken is a free data retrieval call binding the contract method 0x4b57b0be.
//
// Solidity: function wethToken() constant returns(address)
func (_Depositmanager *DepositmanagerCallerSession) WethToken() (common.Address, error) {
	return _Depositmanager.Contract.WethToken(&_Depositmanager.CallOpts)
}

// ChangeRootChain is a paid mutator transaction binding the contract method 0xe8afa8e8.
//
// Solidity: function changeRootChain(address newRootChain) returns()
func (_Depositmanager *DepositmanagerTransactor) ChangeRootChain(opts *bind.TransactOpts, newRootChain common.Address) (*types.Transaction, error) {
	return _Depositmanager.contract.Transact(opts, "changeRootChain", newRootChain)
}

// ChangeRootChain is a paid mutator transaction binding the contract method 0xe8afa8e8.
//
// Solidity: function changeRootChain(address newRootChain) returns()
func (_Depositmanager *DepositmanagerSession) ChangeRootChain(newRootChain common.Address) (*types.Transaction, error) {
	return _Depositmanager.Contract.ChangeRootChain(&_Depositmanager.TransactOpts, newRootChain)
}

// ChangeRootChain is a paid mutator transaction binding the contract method 0xe8afa8e8.
//
// Solidity: function changeRootChain(address newRootChain) returns()
func (_Depositmanager *DepositmanagerTransactorSession) ChangeRootChain(newRootChain common.Address) (*types.Transaction, error) {
	return _Depositmanager.Contract.ChangeRootChain(&_Depositmanager.TransactOpts, newRootChain)
}

// CreateDepositBlock is a paid mutator transaction binding the contract method 0x494f2a78.
//
// Solidity: function createDepositBlock(uint256 _currentHeaderBlock, address _token, address _user, uint256 _amountOrTokenId) returns()
func (_Depositmanager *DepositmanagerTransactor) CreateDepositBlock(opts *bind.TransactOpts, _currentHeaderBlock *big.Int, _token common.Address, _user common.Address, _amountOrTokenId *big.Int) (*types.Transaction, error) {
	return _Depositmanager.contract.Transact(opts, "createDepositBlock", _currentHeaderBlock, _token, _user, _amountOrTokenId)
}

// CreateDepositBlock is a paid mutator transaction binding the contract method 0x494f2a78.
//
// Solidity: function createDepositBlock(uint256 _currentHeaderBlock, address _token, address _user, uint256 _amountOrTokenId) returns()
func (_Depositmanager *DepositmanagerSession) CreateDepositBlock(_currentHeaderBlock *big.Int, _token common.Address, _user common.Address, _amountOrTokenId *big.Int) (*types.Transaction, error) {
	return _Depositmanager.Contract.CreateDepositBlock(&_Depositmanager.TransactOpts, _currentHeaderBlock, _token, _user, _amountOrTokenId)
}

// CreateDepositBlock is a paid mutator transaction binding the contract method 0x494f2a78.
//
// Solidity: function createDepositBlock(uint256 _currentHeaderBlock, address _token, address _user, uint256 _amountOrTokenId) returns()
func (_Depositmanager *DepositmanagerTransactorSession) CreateDepositBlock(_currentHeaderBlock *big.Int, _token common.Address, _user common.Address, _amountOrTokenId *big.Int) (*types.Transaction, error) {
	return _Depositmanager.Contract.CreateDepositBlock(&_Depositmanager.TransactOpts, _currentHeaderBlock, _token, _user, _amountOrTokenId)
}

// FinalizeCommit is a paid mutator transaction binding the contract method 0xfb0df30f.
//
// Solidity: function finalizeCommit(uint256 _currentHeaderBlock) returns()
func (_Depositmanager *DepositmanagerTransactor) FinalizeCommit(opts *bind.TransactOpts, _currentHeaderBlock *big.Int) (*types.Transaction, error) {
	return _Depositmanager.contract.Transact(opts, "finalizeCommit", _currentHeaderBlock)
}

// FinalizeCommit is a paid mutator transaction binding the contract method 0xfb0df30f.
//
// Solidity: function finalizeCommit(uint256 _currentHeaderBlock) returns()
func (_Depositmanager *DepositmanagerSession) FinalizeCommit(_currentHeaderBlock *big.Int) (*types.Transaction, error) {
	return _Depositmanager.Contract.FinalizeCommit(&_Depositmanager.TransactOpts, _currentHeaderBlock)
}

// FinalizeCommit is a paid mutator transaction binding the contract method 0xfb0df30f.
//
// Solidity: function finalizeCommit(uint256 _currentHeaderBlock) returns()
func (_Depositmanager *DepositmanagerTransactorSession) FinalizeCommit(_currentHeaderBlock *big.Int) (*types.Transaction, error) {
	return _Depositmanager.Contract.FinalizeCommit(&_Depositmanager.TransactOpts, _currentHeaderBlock)
}

// MapToken is a paid mutator transaction binding the contract method 0xe117694b.
//
// Solidity: function mapToken(address _rootToken, address _childToken, bool _isERC721) returns()
func (_Depositmanager *DepositmanagerTransactor) MapToken(opts *bind.TransactOpts, _rootToken common.Address, _childToken common.Address, _isERC721 bool) (*types.Transaction, error) {
	return _Depositmanager.contract.Transact(opts, "mapToken", _rootToken, _childToken, _isERC721)
}

// MapToken is a paid mutator transaction binding the contract method 0xe117694b.
//
// Solidity: function mapToken(address _rootToken, address _childToken, bool _isERC721) returns()
func (_Depositmanager *DepositmanagerSession) MapToken(_rootToken common.Address, _childToken common.Address, _isERC721 bool) (*types.Transaction, error) {
	return _Depositmanager.Contract.MapToken(&_Depositmanager.TransactOpts, _rootToken, _childToken, _isERC721)
}

// MapToken is a paid mutator transaction binding the contract method 0xe117694b.
//
// Solidity: function mapToken(address _rootToken, address _childToken, bool _isERC721) returns()
func (_Depositmanager *DepositmanagerTransactorSession) MapToken(_rootToken common.Address, _childToken common.Address, _isERC721 bool) (*types.Transaction, error) {
	return _Depositmanager.Contract.MapToken(&_Depositmanager.TransactOpts, _rootToken, _childToken, _isERC721)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Depositmanager *DepositmanagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Depositmanager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Depositmanager *DepositmanagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _Depositmanager.Contract.RenounceOwnership(&_Depositmanager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Depositmanager *DepositmanagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Depositmanager.Contract.RenounceOwnership(&_Depositmanager.TransactOpts)
}

// SetExitNFTContract is a paid mutator transaction binding the contract method 0x46e11a8d.
//
// Solidity: function setExitNFTContract(address _nftContract) returns()
func (_Depositmanager *DepositmanagerTransactor) SetExitNFTContract(opts *bind.TransactOpts, _nftContract common.Address) (*types.Transaction, error) {
	return _Depositmanager.contract.Transact(opts, "setExitNFTContract", _nftContract)
}

// SetExitNFTContract is a paid mutator transaction binding the contract method 0x46e11a8d.
//
// Solidity: function setExitNFTContract(address _nftContract) returns()
func (_Depositmanager *DepositmanagerSession) SetExitNFTContract(_nftContract common.Address) (*types.Transaction, error) {
	return _Depositmanager.Contract.SetExitNFTContract(&_Depositmanager.TransactOpts, _nftContract)
}

// SetExitNFTContract is a paid mutator transaction binding the contract method 0x46e11a8d.
//
// Solidity: function setExitNFTContract(address _nftContract) returns()
func (_Depositmanager *DepositmanagerTransactorSession) SetExitNFTContract(_nftContract common.Address) (*types.Transaction, error) {
	return _Depositmanager.Contract.SetExitNFTContract(&_Depositmanager.TransactOpts, _nftContract)
}

// SetWETHToken is a paid mutator transaction binding the contract method 0xb45d1f68.
//
// Solidity: function setWETHToken(address _token) returns()
func (_Depositmanager *DepositmanagerTransactor) SetWETHToken(opts *bind.TransactOpts, _token common.Address) (*types.Transaction, error) {
	return _Depositmanager.contract.Transact(opts, "setWETHToken", _token)
}

// SetWETHToken is a paid mutator transaction binding the contract method 0xb45d1f68.
//
// Solidity: function setWETHToken(address _token) returns()
func (_Depositmanager *DepositmanagerSession) SetWETHToken(_token common.Address) (*types.Transaction, error) {
	return _Depositmanager.Contract.SetWETHToken(&_Depositmanager.TransactOpts, _token)
}

// SetWETHToken is a paid mutator transaction binding the contract method 0xb45d1f68.
//
// Solidity: function setWETHToken(address _token) returns()
func (_Depositmanager *DepositmanagerTransactorSession) SetWETHToken(_token common.Address) (*types.Transaction, error) {
	return _Depositmanager.Contract.SetWETHToken(&_Depositmanager.TransactOpts, _token)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Depositmanager *DepositmanagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Depositmanager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Depositmanager *DepositmanagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Depositmanager.Contract.TransferOwnership(&_Depositmanager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Depositmanager *DepositmanagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Depositmanager.Contract.TransferOwnership(&_Depositmanager.TransactOpts, newOwner)
}

// DepositmanagerDepositIterator is returned from FilterDeposit and is used to iterate over the raw logs and unpacked data for Deposit events raised by the Depositmanager contract.
type DepositmanagerDepositIterator struct {
	Event *DepositmanagerDeposit // Event containing the contract specifics and raw log

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
func (it *DepositmanagerDepositIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositmanagerDeposit)
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
		it.Event = new(DepositmanagerDeposit)
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
func (it *DepositmanagerDepositIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositmanagerDepositIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositmanagerDeposit represents a Deposit event raised by the Depositmanager contract.
type DepositmanagerDeposit struct {
	User            common.Address
	Token           common.Address
	AmountOrTokenId *big.Int
	DepositCount    *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterDeposit is a free log retrieval operation binding the contract event 0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7.
//
// Solidity: event Deposit(address indexed _user, address indexed _token, uint256 _amountOrTokenId, uint256 _depositCount)
func (_Depositmanager *DepositmanagerFilterer) FilterDeposit(opts *bind.FilterOpts, _user []common.Address, _token []common.Address) (*DepositmanagerDepositIterator, error) {

	var _userRule []interface{}
	for _, _userItem := range _user {
		_userRule = append(_userRule, _userItem)
	}
	var _tokenRule []interface{}
	for _, _tokenItem := range _token {
		_tokenRule = append(_tokenRule, _tokenItem)
	}

	logs, sub, err := _Depositmanager.contract.FilterLogs(opts, "Deposit", _userRule, _tokenRule)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerDepositIterator{contract: _Depositmanager.contract, event: "Deposit", logs: logs, sub: sub}, nil
}

// WatchDeposit is a free log subscription operation binding the contract event 0xdcbc1c05240f31ff3ad067ef1ee35ce4997762752e3a095284754544f4c709d7.
//
// Solidity: event Deposit(address indexed _user, address indexed _token, uint256 _amountOrTokenId, uint256 _depositCount)
func (_Depositmanager *DepositmanagerFilterer) WatchDeposit(opts *bind.WatchOpts, sink chan<- *DepositmanagerDeposit, _user []common.Address, _token []common.Address) (event.Subscription, error) {

	var _userRule []interface{}
	for _, _userItem := range _user {
		_userRule = append(_userRule, _userItem)
	}
	var _tokenRule []interface{}
	for _, _tokenItem := range _token {
		_tokenRule = append(_tokenRule, _tokenItem)
	}

	logs, sub, err := _Depositmanager.contract.WatchLogs(opts, "Deposit", _userRule, _tokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositmanagerDeposit)
				if err := _Depositmanager.contract.UnpackLog(event, "Deposit", log); err != nil {
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

// DepositmanagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Depositmanager contract.
type DepositmanagerOwnershipTransferredIterator struct {
	Event *DepositmanagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DepositmanagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositmanagerOwnershipTransferred)
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
		it.Event = new(DepositmanagerOwnershipTransferred)
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
func (it *DepositmanagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositmanagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositmanagerOwnershipTransferred represents a OwnershipTransferred event raised by the Depositmanager contract.
type DepositmanagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Depositmanager *DepositmanagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DepositmanagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Depositmanager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerOwnershipTransferredIterator{contract: _Depositmanager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Depositmanager *DepositmanagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DepositmanagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Depositmanager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositmanagerOwnershipTransferred)
				if err := _Depositmanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// DepositmanagerRootChainChangedIterator is returned from FilterRootChainChanged and is used to iterate over the raw logs and unpacked data for RootChainChanged events raised by the Depositmanager contract.
type DepositmanagerRootChainChangedIterator struct {
	Event *DepositmanagerRootChainChanged // Event containing the contract specifics and raw log

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
func (it *DepositmanagerRootChainChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositmanagerRootChainChanged)
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
		it.Event = new(DepositmanagerRootChainChanged)
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
func (it *DepositmanagerRootChainChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositmanagerRootChainChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositmanagerRootChainChanged represents a RootChainChanged event raised by the Depositmanager contract.
type DepositmanagerRootChainChanged struct {
	PreviousRootChain common.Address
	NewRootChain      common.Address
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRootChainChanged is a free log retrieval operation binding the contract event 0x211c9015fc81c0dbd45bd99f0f29fc1c143bfd53442d5ffd722bbbef7a887fe9.
//
// Solidity: event RootChainChanged(address indexed previousRootChain, address indexed newRootChain)
func (_Depositmanager *DepositmanagerFilterer) FilterRootChainChanged(opts *bind.FilterOpts, previousRootChain []common.Address, newRootChain []common.Address) (*DepositmanagerRootChainChangedIterator, error) {

	var previousRootChainRule []interface{}
	for _, previousRootChainItem := range previousRootChain {
		previousRootChainRule = append(previousRootChainRule, previousRootChainItem)
	}
	var newRootChainRule []interface{}
	for _, newRootChainItem := range newRootChain {
		newRootChainRule = append(newRootChainRule, newRootChainItem)
	}

	logs, sub, err := _Depositmanager.contract.FilterLogs(opts, "RootChainChanged", previousRootChainRule, newRootChainRule)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerRootChainChangedIterator{contract: _Depositmanager.contract, event: "RootChainChanged", logs: logs, sub: sub}, nil
}

// WatchRootChainChanged is a free log subscription operation binding the contract event 0x211c9015fc81c0dbd45bd99f0f29fc1c143bfd53442d5ffd722bbbef7a887fe9.
//
// Solidity: event RootChainChanged(address indexed previousRootChain, address indexed newRootChain)
func (_Depositmanager *DepositmanagerFilterer) WatchRootChainChanged(opts *bind.WatchOpts, sink chan<- *DepositmanagerRootChainChanged, previousRootChain []common.Address, newRootChain []common.Address) (event.Subscription, error) {

	var previousRootChainRule []interface{}
	for _, previousRootChainItem := range previousRootChain {
		previousRootChainRule = append(previousRootChainRule, previousRootChainItem)
	}
	var newRootChainRule []interface{}
	for _, newRootChainItem := range newRootChain {
		newRootChainRule = append(newRootChainRule, newRootChainItem)
	}

	logs, sub, err := _Depositmanager.contract.WatchLogs(opts, "RootChainChanged", previousRootChainRule, newRootChainRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositmanagerRootChainChanged)
				if err := _Depositmanager.contract.UnpackLog(event, "RootChainChanged", log); err != nil {
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

// DepositmanagerTokenMappedIterator is returned from FilterTokenMapped and is used to iterate over the raw logs and unpacked data for TokenMapped events raised by the Depositmanager contract.
type DepositmanagerTokenMappedIterator struct {
	Event *DepositmanagerTokenMapped // Event containing the contract specifics and raw log

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
func (it *DepositmanagerTokenMappedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DepositmanagerTokenMapped)
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
		it.Event = new(DepositmanagerTokenMapped)
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
func (it *DepositmanagerTokenMappedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DepositmanagerTokenMappedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DepositmanagerTokenMapped represents a TokenMapped event raised by the Depositmanager contract.
type DepositmanagerTokenMapped struct {
	RootToken  common.Address
	ChildToken common.Address
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterTokenMapped is a free log retrieval operation binding the contract event 0x85920d35e6c72f6b2affffa04298b0cecfeba86e4a9f407df661f1cb8ab5e617.
//
// Solidity: event TokenMapped(address indexed _rootToken, address indexed _childToken)
func (_Depositmanager *DepositmanagerFilterer) FilterTokenMapped(opts *bind.FilterOpts, _rootToken []common.Address, _childToken []common.Address) (*DepositmanagerTokenMappedIterator, error) {

	var _rootTokenRule []interface{}
	for _, _rootTokenItem := range _rootToken {
		_rootTokenRule = append(_rootTokenRule, _rootTokenItem)
	}
	var _childTokenRule []interface{}
	for _, _childTokenItem := range _childToken {
		_childTokenRule = append(_childTokenRule, _childTokenItem)
	}

	logs, sub, err := _Depositmanager.contract.FilterLogs(opts, "TokenMapped", _rootTokenRule, _childTokenRule)
	if err != nil {
		return nil, err
	}
	return &DepositmanagerTokenMappedIterator{contract: _Depositmanager.contract, event: "TokenMapped", logs: logs, sub: sub}, nil
}

// WatchTokenMapped is a free log subscription operation binding the contract event 0x85920d35e6c72f6b2affffa04298b0cecfeba86e4a9f407df661f1cb8ab5e617.
//
// Solidity: event TokenMapped(address indexed _rootToken, address indexed _childToken)
func (_Depositmanager *DepositmanagerFilterer) WatchTokenMapped(opts *bind.WatchOpts, sink chan<- *DepositmanagerTokenMapped, _rootToken []common.Address, _childToken []common.Address) (event.Subscription, error) {

	var _rootTokenRule []interface{}
	for _, _rootTokenItem := range _rootToken {
		_rootTokenRule = append(_rootTokenRule, _rootTokenItem)
	}
	var _childTokenRule []interface{}
	for _, _childTokenItem := range _childToken {
		_childTokenRule = append(_childTokenRule, _childTokenItem)
	}

	logs, sub, err := _Depositmanager.contract.WatchLogs(opts, "TokenMapped", _rootTokenRule, _childTokenRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DepositmanagerTokenMapped)
				if err := _Depositmanager.contract.UnpackLog(event, "TokenMapped", log); err != nil {
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
