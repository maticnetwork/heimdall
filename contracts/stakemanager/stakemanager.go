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
const StakemanagerABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"validator\",\"type\":\"address\"},{\"name\":\"votingPower\",\"type\":\"uint256\"},{\"name\":\"_pubkey\",\"type\":\"string\"}],\"name\":\"addValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"NewProposer\",\"type\":\"event\"},{\"constant\":false,\"inputs\":[{\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"getSha256\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes20\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"removeValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"vote\",\"type\":\"bytes\"},{\"name\":\"sigs\",\"type\":\"bytes\"},{\"name\":\"extradata\",\"type\":\"bytes\"}],\"name\":\"validate\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chain\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ecrecovery\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"ecrecovery\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"ecrecoveryFromData\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"sig\",\"type\":\"bytes\"},{\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"ecverify\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getPubkey\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getValidatorSet\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"},{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"roundType\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes32\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"startBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validators\",\"outputs\":[{\"name\":\"votingPower\",\"type\":\"uint256\"},{\"name\":\"validator\",\"type\":\"address\"},{\"name\":\"pubkey\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"voteType\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes1\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"}]"

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

// Chain is a free data retrieval call binding the contract method 0xc763e5a1.
//
// Solidity: function chain() constant returns(bytes32)
func (_Stakemanager *StakemanagerCaller) Chain(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "chain")
	return *ret0, err
}

// Chain is a free data retrieval call binding the contract method 0xc763e5a1.
//
// Solidity: function chain() constant returns(bytes32)
func (_Stakemanager *StakemanagerSession) Chain() ([32]byte, error) {
	return _Stakemanager.Contract.Chain(&_Stakemanager.CallOpts)
}

// Chain is a free data retrieval call binding the contract method 0xc763e5a1.
//
// Solidity: function chain() constant returns(bytes32)
func (_Stakemanager *StakemanagerCallerSession) Chain() ([32]byte, error) {
	return _Stakemanager.Contract.Chain(&_Stakemanager.CallOpts)
}

// Ecrecovery is a free data retrieval call binding the contract method 0x77d32e94.
//
// Solidity: function ecrecovery(hash bytes32, sig bytes) constant returns(address)
func (_Stakemanager *StakemanagerCaller) Ecrecovery(opts *bind.CallOpts, hash [32]byte, sig []byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "ecrecovery", hash, sig)
	return *ret0, err
}

// Ecrecovery is a free data retrieval call binding the contract method 0x77d32e94.
//
// Solidity: function ecrecovery(hash bytes32, sig bytes) constant returns(address)
func (_Stakemanager *StakemanagerSession) Ecrecovery(hash [32]byte, sig []byte) (common.Address, error) {
	return _Stakemanager.Contract.Ecrecovery(&_Stakemanager.CallOpts, hash, sig)
}

// Ecrecovery is a free data retrieval call binding the contract method 0x77d32e94.
//
// Solidity: function ecrecovery(hash bytes32, sig bytes) constant returns(address)
func (_Stakemanager *StakemanagerCallerSession) Ecrecovery(hash [32]byte, sig []byte) (common.Address, error) {
	return _Stakemanager.Contract.Ecrecovery(&_Stakemanager.CallOpts, hash, sig)
}

// EcrecoveryFromData is a free data retrieval call binding the contract method 0xba0e7252.
//
// Solidity: function ecrecoveryFromData(data bytes, sig bytes) constant returns(address)
func (_Stakemanager *StakemanagerCaller) EcrecoveryFromData(opts *bind.CallOpts, data []byte, sig []byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "ecrecoveryFromData", data, sig)
	return *ret0, err
}

// EcrecoveryFromData is a free data retrieval call binding the contract method 0xba0e7252.
//
// Solidity: function ecrecoveryFromData(data bytes, sig bytes) constant returns(address)
func (_Stakemanager *StakemanagerSession) EcrecoveryFromData(data []byte, sig []byte) (common.Address, error) {
	return _Stakemanager.Contract.EcrecoveryFromData(&_Stakemanager.CallOpts, data, sig)
}

// EcrecoveryFromData is a free data retrieval call binding the contract method 0xba0e7252.
//
// Solidity: function ecrecoveryFromData(data bytes, sig bytes) constant returns(address)
func (_Stakemanager *StakemanagerCallerSession) EcrecoveryFromData(data []byte, sig []byte) (common.Address, error) {
	return _Stakemanager.Contract.EcrecoveryFromData(&_Stakemanager.CallOpts, data, sig)
}

// Ecverify is a free data retrieval call binding the contract method 0x39cdde32.
//
// Solidity: function ecverify(hash bytes32, sig bytes, signer address) constant returns(bool)
func (_Stakemanager *StakemanagerCaller) Ecverify(opts *bind.CallOpts, hash [32]byte, sig []byte, signer common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "ecverify", hash, sig, signer)
	return *ret0, err
}

// Ecverify is a free data retrieval call binding the contract method 0x39cdde32.
//
// Solidity: function ecverify(hash bytes32, sig bytes, signer address) constant returns(bool)
func (_Stakemanager *StakemanagerSession) Ecverify(hash [32]byte, sig []byte, signer common.Address) (bool, error) {
	return _Stakemanager.Contract.Ecverify(&_Stakemanager.CallOpts, hash, sig, signer)
}

// Ecverify is a free data retrieval call binding the contract method 0x39cdde32.
//
// Solidity: function ecverify(hash bytes32, sig bytes, signer address) constant returns(bool)
func (_Stakemanager *StakemanagerCallerSession) Ecverify(hash [32]byte, sig []byte, signer common.Address) (bool, error) {
	return _Stakemanager.Contract.Ecverify(&_Stakemanager.CallOpts, hash, sig, signer)
}

// GetPubkey is a free data retrieval call binding the contract method 0xef6fd878.
//
// Solidity: function getPubkey(index uint256) constant returns(string)
func (_Stakemanager *StakemanagerCaller) GetPubkey(opts *bind.CallOpts, index *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "getPubkey", index)
	return *ret0, err
}

// GetPubkey is a free data retrieval call binding the contract method 0xef6fd878.
//
// Solidity: function getPubkey(index uint256) constant returns(string)
func (_Stakemanager *StakemanagerSession) GetPubkey(index *big.Int) (string, error) {
	return _Stakemanager.Contract.GetPubkey(&_Stakemanager.CallOpts, index)
}

// GetPubkey is a free data retrieval call binding the contract method 0xef6fd878.
//
// Solidity: function getPubkey(index uint256) constant returns(string)
func (_Stakemanager *StakemanagerCallerSession) GetPubkey(index *big.Int) (string, error) {
	return _Stakemanager.Contract.GetPubkey(&_Stakemanager.CallOpts, index)
}

// GetValidatorSet is a free data retrieval call binding the contract method 0xcf331250.
//
// Solidity: function getValidatorSet() constant returns(uint256[], address[])
func (_Stakemanager *StakemanagerCaller) GetValidatorSet(opts *bind.CallOpts) ([]*big.Int, []common.Address, error) {
	var (
		ret0 = new([]*big.Int)
		ret1 = new([]common.Address)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _Stakemanager.contract.Call(opts, out, "getValidatorSet")
	return *ret0, *ret1, err
}

// GetValidatorSet is a free data retrieval call binding the contract method 0xcf331250.
//
// Solidity: function getValidatorSet() constant returns(uint256[], address[])
func (_Stakemanager *StakemanagerSession) GetValidatorSet() ([]*big.Int, []common.Address, error) {
	return _Stakemanager.Contract.GetValidatorSet(&_Stakemanager.CallOpts)
}

// GetValidatorSet is a free data retrieval call binding the contract method 0xcf331250.
//
// Solidity: function getValidatorSet() constant returns(uint256[], address[])
func (_Stakemanager *StakemanagerCallerSession) GetValidatorSet() ([]*big.Int, []common.Address, error) {
	return _Stakemanager.Contract.GetValidatorSet(&_Stakemanager.CallOpts)
}

// RoundType is a free data retrieval call binding the contract method 0x2c2d1a3b.
//
// Solidity: function roundType() constant returns(bytes32)
func (_Stakemanager *StakemanagerCaller) RoundType(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "roundType")
	return *ret0, err
}

// RoundType is a free data retrieval call binding the contract method 0x2c2d1a3b.
//
// Solidity: function roundType() constant returns(bytes32)
func (_Stakemanager *StakemanagerSession) RoundType() ([32]byte, error) {
	return _Stakemanager.Contract.RoundType(&_Stakemanager.CallOpts)
}

// RoundType is a free data retrieval call binding the contract method 0x2c2d1a3b.
//
// Solidity: function roundType() constant returns(bytes32)
func (_Stakemanager *StakemanagerCallerSession) RoundType() ([32]byte, error) {
	return _Stakemanager.Contract.RoundType(&_Stakemanager.CallOpts)
}

// StartBlock is a free data retrieval call binding the contract method 0x48cd4cb1.
//
// Solidity: function startBlock() constant returns(uint256)
func (_Stakemanager *StakemanagerCaller) StartBlock(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "startBlock")
	return *ret0, err
}

// StartBlock is a free data retrieval call binding the contract method 0x48cd4cb1.
//
// Solidity: function startBlock() constant returns(uint256)
func (_Stakemanager *StakemanagerSession) StartBlock() (*big.Int, error) {
	return _Stakemanager.Contract.StartBlock(&_Stakemanager.CallOpts)
}

// StartBlock is a free data retrieval call binding the contract method 0x48cd4cb1.
//
// Solidity: function startBlock() constant returns(uint256)
func (_Stakemanager *StakemanagerCallerSession) StartBlock() (*big.Int, error) {
	return _Stakemanager.Contract.StartBlock(&_Stakemanager.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(votingPower uint256, validator address, pubkey string)
func (_Stakemanager *StakemanagerCaller) Validators(opts *bind.CallOpts, arg0 *big.Int) (struct {
	VotingPower *big.Int
	Validator   common.Address
	Pubkey      string
}, error) {
	ret := new(struct {
		VotingPower *big.Int
		Validator   common.Address
		Pubkey      string
	})
	out := ret
	err := _Stakemanager.contract.Call(opts, out, "validators", arg0)
	return *ret, err
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(votingPower uint256, validator address, pubkey string)
func (_Stakemanager *StakemanagerSession) Validators(arg0 *big.Int) (struct {
	VotingPower *big.Int
	Validator   common.Address
	Pubkey      string
}, error) {
	return _Stakemanager.Contract.Validators(&_Stakemanager.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(votingPower uint256, validator address, pubkey string)
func (_Stakemanager *StakemanagerCallerSession) Validators(arg0 *big.Int) (struct {
	VotingPower *big.Int
	Validator   common.Address
	Pubkey      string
}, error) {
	return _Stakemanager.Contract.Validators(&_Stakemanager.CallOpts, arg0)
}

// VoteType is a free data retrieval call binding the contract method 0x7d1a3d37.
//
// Solidity: function voteType() constant returns(bytes1)
func (_Stakemanager *StakemanagerCaller) VoteType(opts *bind.CallOpts) ([1]byte, error) {
	var (
		ret0 = new([1]byte)
	)
	out := ret0
	err := _Stakemanager.contract.Call(opts, out, "voteType")
	return *ret0, err
}

// VoteType is a free data retrieval call binding the contract method 0x7d1a3d37.
//
// Solidity: function voteType() constant returns(bytes1)
func (_Stakemanager *StakemanagerSession) VoteType() ([1]byte, error) {
	return _Stakemanager.Contract.VoteType(&_Stakemanager.CallOpts)
}

// VoteType is a free data retrieval call binding the contract method 0x7d1a3d37.
//
// Solidity: function voteType() constant returns(bytes1)
func (_Stakemanager *StakemanagerCallerSession) VoteType() ([1]byte, error) {
	return _Stakemanager.Contract.VoteType(&_Stakemanager.CallOpts)
}

// AddValidator is a paid mutator transaction binding the contract method 0x01736c35.
//
// Solidity: function addValidator(validator address, votingPower uint256, _pubkey string) returns()
func (_Stakemanager *StakemanagerTransactor) AddValidator(opts *bind.TransactOpts, validator common.Address, votingPower *big.Int, _pubkey string) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "addValidator", validator, votingPower, _pubkey)
}

// AddValidator is a paid mutator transaction binding the contract method 0x01736c35.
//
// Solidity: function addValidator(validator address, votingPower uint256, _pubkey string) returns()
func (_Stakemanager *StakemanagerSession) AddValidator(validator common.Address, votingPower *big.Int, _pubkey string) (*types.Transaction, error) {
	return _Stakemanager.Contract.AddValidator(&_Stakemanager.TransactOpts, validator, votingPower, _pubkey)
}

// AddValidator is a paid mutator transaction binding the contract method 0x01736c35.
//
// Solidity: function addValidator(validator address, votingPower uint256, _pubkey string) returns()
func (_Stakemanager *StakemanagerTransactorSession) AddValidator(validator common.Address, votingPower *big.Int, _pubkey string) (*types.Transaction, error) {
	return _Stakemanager.Contract.AddValidator(&_Stakemanager.TransactOpts, validator, votingPower, _pubkey)
}

// GetSha256 is a paid mutator transaction binding the contract method 0x8b053758.
//
// Solidity: function getSha256(input bytes) returns(bytes20)
func (_Stakemanager *StakemanagerTransactor) GetSha256(opts *bind.TransactOpts, input []byte) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "getSha256", input)
}

// GetSha256 is a paid mutator transaction binding the contract method 0x8b053758.
//
// Solidity: function getSha256(input bytes) returns(bytes20)
func (_Stakemanager *StakemanagerSession) GetSha256(input []byte) (*types.Transaction, error) {
	return _Stakemanager.Contract.GetSha256(&_Stakemanager.TransactOpts, input)
}

// GetSha256 is a paid mutator transaction binding the contract method 0x8b053758.
//
// Solidity: function getSha256(input bytes) returns(bytes20)
func (_Stakemanager *StakemanagerTransactorSession) GetSha256(input []byte) (*types.Transaction, error) {
	return _Stakemanager.Contract.GetSha256(&_Stakemanager.TransactOpts, input)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0xf94e1867.
//
// Solidity: function removeValidator(_index uint256) returns()
func (_Stakemanager *StakemanagerTransactor) RemoveValidator(opts *bind.TransactOpts, _index *big.Int) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "removeValidator", _index)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0xf94e1867.
//
// Solidity: function removeValidator(_index uint256) returns()
func (_Stakemanager *StakemanagerSession) RemoveValidator(_index *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.RemoveValidator(&_Stakemanager.TransactOpts, _index)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0xf94e1867.
//
// Solidity: function removeValidator(_index uint256) returns()
func (_Stakemanager *StakemanagerTransactorSession) RemoveValidator(_index *big.Int) (*types.Transaction, error) {
	return _Stakemanager.Contract.RemoveValidator(&_Stakemanager.TransactOpts, _index)
}

// Validate is a paid mutator transaction binding the contract method 0x2d16c59c.
//
// Solidity: function validate(vote bytes, sigs bytes, extradata bytes) returns(address, address, uint256)
func (_Stakemanager *StakemanagerTransactor) Validate(opts *bind.TransactOpts, vote []byte, sigs []byte, extradata []byte) (*types.Transaction, error) {
	return _Stakemanager.contract.Transact(opts, "validate", vote, sigs, extradata)
}

// Validate is a paid mutator transaction binding the contract method 0x2d16c59c.
//
// Solidity: function validate(vote bytes, sigs bytes, extradata bytes) returns(address, address, uint256)
func (_Stakemanager *StakemanagerSession) Validate(vote []byte, sigs []byte, extradata []byte) (*types.Transaction, error) {
	return _Stakemanager.Contract.Validate(&_Stakemanager.TransactOpts, vote, sigs, extradata)
}

// Validate is a paid mutator transaction binding the contract method 0x2d16c59c.
//
// Solidity: function validate(vote bytes, sigs bytes, extradata bytes) returns(address, address, uint256)
func (_Stakemanager *StakemanagerTransactorSession) Validate(vote []byte, sigs []byte, extradata []byte) (*types.Transaction, error) {
	return _Stakemanager.Contract.Validate(&_Stakemanager.TransactOpts, vote, sigs, extradata)
}

// StakemanagerNewProposerIterator is returned from FilterNewProposer and is used to iterate over the raw logs and unpacked data for NewProposer events raised by the Stakemanager contract.
type StakemanagerNewProposerIterator struct {
	Event *StakemanagerNewProposer // Event containing the contract specifics and raw log

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
func (it *StakemanagerNewProposerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerNewProposer)
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
		it.Event = new(StakemanagerNewProposer)
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
func (it *StakemanagerNewProposerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerNewProposerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerNewProposer represents a NewProposer event raised by the Stakemanager contract.
type StakemanagerNewProposer struct {
	User common.Address
	Data []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNewProposer is a free log retrieval operation binding the contract event 0x11d5415dca2a9428143948f67030a63ac42cbd3639a4aae25e602fbbf5da38db.
//
// Solidity: e NewProposer(user indexed address, data bytes)
func (_Stakemanager *StakemanagerFilterer) FilterNewProposer(opts *bind.FilterOpts, user []common.Address) (*StakemanagerNewProposerIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "NewProposer", userRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerNewProposerIterator{contract: _Stakemanager.contract, event: "NewProposer", logs: logs, sub: sub}, nil
}

// WatchNewProposer is a free log subscription operation binding the contract event 0x11d5415dca2a9428143948f67030a63ac42cbd3639a4aae25e602fbbf5da38db.
//
// Solidity: e NewProposer(user indexed address, data bytes)
func (_Stakemanager *StakemanagerFilterer) WatchNewProposer(opts *bind.WatchOpts, sink chan<- *StakemanagerNewProposer, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "NewProposer", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerNewProposer)
				if err := _Stakemanager.contract.UnpackLog(event, "NewProposer", log); err != nil {
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
