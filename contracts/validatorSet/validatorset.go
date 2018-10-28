// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package validatorSet

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

// ValidatorSetABI is the input ABI used to generate the binding from.
const ValidatorSetABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"validator\",\"type\":\"address\"},{\"name\":\"votingPower\",\"type\":\"uint256\"},{\"name\":\"_pubkey\",\"type\":\"string\"}],\"name\":\"addValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"roundType\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"vote\",\"type\":\"bytes\"},{\"name\":\"sigs\",\"type\":\"bytes\"},{\"name\":\"extradata\",\"type\":\"bytes\"}],\"name\":\"validate\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validators\",\"outputs\":[{\"name\":\"votingPower\",\"type\":\"uint256\"},{\"name\":\"validator\",\"type\":\"address\"},{\"name\":\"pubkey\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"sig\",\"type\":\"bytes\"},{\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"ecverify\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"lowestPower\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalVotingPower\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"ecrecovery\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_chainID\",\"type\":\"string\"}],\"name\":\"setChainId\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"voteType\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"input\",\"type\":\"bytes\"}],\"name\":\"getSha256\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes20\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"hash\",\"type\":\"bytes32\"},{\"name\":\"v\",\"type\":\"uint8\"},{\"name\":\"r\",\"type\":\"bytes32\"},{\"name\":\"s\",\"type\":\"bytes32\"}],\"name\":\"ecrecovery\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"proposer\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"chainID\",\"outputs\":[{\"name\":\"\",\"type\":\"bytes\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"data\",\"type\":\"bytes\"},{\"name\":\"sig\",\"type\":\"bytes\"}],\"name\":\"ecrecoveryFromData\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getValidatorSet\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"},{\"name\":\"\",\"type\":\"address[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"getPubkey\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"_index\",\"type\":\"uint256\"}],\"name\":\"removeValidator\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"NewProposer\",\"type\":\"event\"}]"

// ValidatorSet is an auto generated Go binding around an Ethereum contract.
type ValidatorSet struct {
	ValidatorSetCaller     // Read-only binding to the contract
	ValidatorSetTransactor // Write-only binding to the contract
	ValidatorSetFilterer   // Log filterer for contract events
}

// ValidatorSetCaller is an auto generated read-only Go binding around an Ethereum contract.
type ValidatorSetCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorSetTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ValidatorSetTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorSetFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ValidatorSetFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ValidatorSetSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ValidatorSetSession struct {
	Contract     *ValidatorSet     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ValidatorSetCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ValidatorSetCallerSession struct {
	Contract *ValidatorSetCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// ValidatorSetTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ValidatorSetTransactorSession struct {
	Contract     *ValidatorSetTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// ValidatorSetRaw is an auto generated low-level Go binding around an Ethereum contract.
type ValidatorSetRaw struct {
	Contract *ValidatorSet // Generic contract binding to access the raw methods on
}

// ValidatorSetCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ValidatorSetCallerRaw struct {
	Contract *ValidatorSetCaller // Generic read-only contract binding to access the raw methods on
}

// ValidatorSetTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ValidatorSetTransactorRaw struct {
	Contract *ValidatorSetTransactor // Generic write-only contract binding to access the raw methods on
}

// NewValidatorSet creates a new instance of ValidatorSet, bound to a specific deployed contract.
func NewValidatorSet(address common.Address, backend bind.ContractBackend) (*ValidatorSet, error) {
	contract, err := bindValidatorSet(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ValidatorSet{ValidatorSetCaller: ValidatorSetCaller{contract: contract}, ValidatorSetTransactor: ValidatorSetTransactor{contract: contract}, ValidatorSetFilterer: ValidatorSetFilterer{contract: contract}}, nil
}

// NewValidatorSetCaller creates a new read-only instance of ValidatorSet, bound to a specific deployed contract.
func NewValidatorSetCaller(address common.Address, caller bind.ContractCaller) (*ValidatorSetCaller, error) {
	contract, err := bindValidatorSet(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatorSetCaller{contract: contract}, nil
}

// NewValidatorSetTransactor creates a new write-only instance of ValidatorSet, bound to a specific deployed contract.
func NewValidatorSetTransactor(address common.Address, transactor bind.ContractTransactor) (*ValidatorSetTransactor, error) {
	contract, err := bindValidatorSet(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ValidatorSetTransactor{contract: contract}, nil
}

// NewValidatorSetFilterer creates a new log filterer instance of ValidatorSet, bound to a specific deployed contract.
func NewValidatorSetFilterer(address common.Address, filterer bind.ContractFilterer) (*ValidatorSetFilterer, error) {
	contract, err := bindValidatorSet(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ValidatorSetFilterer{contract: contract}, nil
}

// bindValidatorSet binds a generic wrapper to an already deployed contract.
func bindValidatorSet(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ValidatorSetABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ValidatorSet *ValidatorSetRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ValidatorSet.Contract.ValidatorSetCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ValidatorSet *ValidatorSetRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorSet.Contract.ValidatorSetTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ValidatorSet *ValidatorSetRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ValidatorSet.Contract.ValidatorSetTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ValidatorSet *ValidatorSetCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _ValidatorSet.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ValidatorSet *ValidatorSetTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ValidatorSet.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ValidatorSet *ValidatorSetTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ValidatorSet.Contract.contract.Transact(opts, method, params...)
}

// ChainID is a free data retrieval call binding the contract method 0xadc879e9.
//
// Solidity: function chainID() constant returns(bytes)
func (_ValidatorSet *ValidatorSetCaller) ChainID(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "chainID")
	return *ret0, err
}

// ChainID is a free data retrieval call binding the contract method 0xadc879e9.
//
// Solidity: function chainID() constant returns(bytes)
func (_ValidatorSet *ValidatorSetSession) ChainID() ([]byte, error) {
	return _ValidatorSet.Contract.ChainID(&_ValidatorSet.CallOpts)
}

// ChainID is a free data retrieval call binding the contract method 0xadc879e9.
//
// Solidity: function chainID() constant returns(bytes)
func (_ValidatorSet *ValidatorSetCallerSession) ChainID() ([]byte, error) {
	return _ValidatorSet.Contract.ChainID(&_ValidatorSet.CallOpts)
}

// Ecrecovery is a free data retrieval call binding the contract method 0x98ea1c51.
//
// Solidity: function ecrecovery(hash bytes32, v uint8, r bytes32, s bytes32) constant returns(address)
func (_ValidatorSet *ValidatorSetCaller) Ecrecovery(opts *bind.CallOpts, hash [32]byte, v uint8, r [32]byte, s [32]byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "ecrecovery", hash, v, r, s)
	return *ret0, err
}

// Ecrecovery is a free data retrieval call binding the contract method 0x98ea1c51.
//
// Solidity: function ecrecovery(hash bytes32, v uint8, r bytes32, s bytes32) constant returns(address)
func (_ValidatorSet *ValidatorSetSession) Ecrecovery(hash [32]byte, v uint8, r [32]byte, s [32]byte) (common.Address, error) {
	return _ValidatorSet.Contract.Ecrecovery(&_ValidatorSet.CallOpts, hash, v, r, s)
}

// Ecrecovery is a free data retrieval call binding the contract method 0x98ea1c51.
//
// Solidity: function ecrecovery(hash bytes32, v uint8, r bytes32, s bytes32) constant returns(address)
func (_ValidatorSet *ValidatorSetCallerSession) Ecrecovery(hash [32]byte, v uint8, r [32]byte, s [32]byte) (common.Address, error) {
	return _ValidatorSet.Contract.Ecrecovery(&_ValidatorSet.CallOpts, hash, v, r, s)
}

// EcrecoveryFromData is a free data retrieval call binding the contract method 0xba0e7252.
//
// Solidity: function ecrecoveryFromData(data bytes, sig bytes) constant returns(address)
func (_ValidatorSet *ValidatorSetCaller) EcrecoveryFromData(opts *bind.CallOpts, data []byte, sig []byte) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "ecrecoveryFromData", data, sig)
	return *ret0, err
}

// EcrecoveryFromData is a free data retrieval call binding the contract method 0xba0e7252.
//
// Solidity: function ecrecoveryFromData(data bytes, sig bytes) constant returns(address)
func (_ValidatorSet *ValidatorSetSession) EcrecoveryFromData(data []byte, sig []byte) (common.Address, error) {
	return _ValidatorSet.Contract.EcrecoveryFromData(&_ValidatorSet.CallOpts, data, sig)
}

// EcrecoveryFromData is a free data retrieval call binding the contract method 0xba0e7252.
//
// Solidity: function ecrecoveryFromData(data bytes, sig bytes) constant returns(address)
func (_ValidatorSet *ValidatorSetCallerSession) EcrecoveryFromData(data []byte, sig []byte) (common.Address, error) {
	return _ValidatorSet.Contract.EcrecoveryFromData(&_ValidatorSet.CallOpts, data, sig)
}

// Ecverify is a free data retrieval call binding the contract method 0x39cdde32.
//
// Solidity: function ecverify(hash bytes32, sig bytes, signer address) constant returns(bool)
func (_ValidatorSet *ValidatorSetCaller) Ecverify(opts *bind.CallOpts, hash [32]byte, sig []byte, signer common.Address) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "ecverify", hash, sig, signer)
	return *ret0, err
}

// Ecverify is a free data retrieval call binding the contract method 0x39cdde32.
//
// Solidity: function ecverify(hash bytes32, sig bytes, signer address) constant returns(bool)
func (_ValidatorSet *ValidatorSetSession) Ecverify(hash [32]byte, sig []byte, signer common.Address) (bool, error) {
	return _ValidatorSet.Contract.Ecverify(&_ValidatorSet.CallOpts, hash, sig, signer)
}

// Ecverify is a free data retrieval call binding the contract method 0x39cdde32.
//
// Solidity: function ecverify(hash bytes32, sig bytes, signer address) constant returns(bool)
func (_ValidatorSet *ValidatorSetCallerSession) Ecverify(hash [32]byte, sig []byte, signer common.Address) (bool, error) {
	return _ValidatorSet.Contract.Ecverify(&_ValidatorSet.CallOpts, hash, sig, signer)
}

// GetPubkey is a free data retrieval call binding the contract method 0xef6fd878.
//
// Solidity: function getPubkey(index uint256) constant returns(string)
func (_ValidatorSet *ValidatorSetCaller) GetPubkey(opts *bind.CallOpts, index *big.Int) (string, error) {
	var (
		ret0 = new(string)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "getPubkey", index)
	return *ret0, err
}

// GetPubkey is a free data retrieval call binding the contract method 0xef6fd878.
//
// Solidity: function getPubkey(index uint256) constant returns(string)
func (_ValidatorSet *ValidatorSetSession) GetPubkey(index *big.Int) (string, error) {
	return _ValidatorSet.Contract.GetPubkey(&_ValidatorSet.CallOpts, index)
}

// GetPubkey is a free data retrieval call binding the contract method 0xef6fd878.
//
// Solidity: function getPubkey(index uint256) constant returns(string)
func (_ValidatorSet *ValidatorSetCallerSession) GetPubkey(index *big.Int) (string, error) {
	return _ValidatorSet.Contract.GetPubkey(&_ValidatorSet.CallOpts, index)
}

// GetValidatorSet is a free data retrieval call binding the contract method 0xcf331250.
//
// Solidity: function getValidatorSet() constant returns(uint256[], address[])
func (_ValidatorSet *ValidatorSetCaller) GetValidatorSet(opts *bind.CallOpts) ([]*big.Int, []common.Address, error) {
	var (
		ret0 = new([]*big.Int)
		ret1 = new([]common.Address)
	)
	out := &[]interface{}{
		ret0,
		ret1,
	}
	err := _ValidatorSet.contract.Call(opts, out, "getValidatorSet")
	return *ret0, *ret1, err
}

// GetValidatorSet is a free data retrieval call binding the contract method 0xcf331250.
//
// Solidity: function getValidatorSet() constant returns(uint256[], address[])
func (_ValidatorSet *ValidatorSetSession) GetValidatorSet() ([]*big.Int, []common.Address, error) {
	return _ValidatorSet.Contract.GetValidatorSet(&_ValidatorSet.CallOpts)
}

// GetValidatorSet is a free data retrieval call binding the contract method 0xcf331250.
//
// Solidity: function getValidatorSet() constant returns(uint256[], address[])
func (_ValidatorSet *ValidatorSetCallerSession) GetValidatorSet() ([]*big.Int, []common.Address, error) {
	return _ValidatorSet.Contract.GetValidatorSet(&_ValidatorSet.CallOpts)
}

// LowestPower is a free data retrieval call binding the contract method 0x5fcf1c9b.
//
// Solidity: function lowestPower() constant returns(uint256)
func (_ValidatorSet *ValidatorSetCaller) LowestPower(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "lowestPower")
	return *ret0, err
}

// LowestPower is a free data retrieval call binding the contract method 0x5fcf1c9b.
//
// Solidity: function lowestPower() constant returns(uint256)
func (_ValidatorSet *ValidatorSetSession) LowestPower() (*big.Int, error) {
	return _ValidatorSet.Contract.LowestPower(&_ValidatorSet.CallOpts)
}

// LowestPower is a free data retrieval call binding the contract method 0x5fcf1c9b.
//
// Solidity: function lowestPower() constant returns(uint256)
func (_ValidatorSet *ValidatorSetCallerSession) LowestPower() (*big.Int, error) {
	return _ValidatorSet.Contract.LowestPower(&_ValidatorSet.CallOpts)
}

// Proposer is a free data retrieval call binding the contract method 0xa8e4fb90.
//
// Solidity: function proposer() constant returns(address)
func (_ValidatorSet *ValidatorSetCaller) Proposer(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "proposer")
	return *ret0, err
}

// Proposer is a free data retrieval call binding the contract method 0xa8e4fb90.
//
// Solidity: function proposer() constant returns(address)
func (_ValidatorSet *ValidatorSetSession) Proposer() (common.Address, error) {
	return _ValidatorSet.Contract.Proposer(&_ValidatorSet.CallOpts)
}

// Proposer is a free data retrieval call binding the contract method 0xa8e4fb90.
//
// Solidity: function proposer() constant returns(address)
func (_ValidatorSet *ValidatorSetCallerSession) Proposer() (common.Address, error) {
	return _ValidatorSet.Contract.Proposer(&_ValidatorSet.CallOpts)
}

// RoundType is a free data retrieval call binding the contract method 0x2c2d1a3b.
//
// Solidity: function roundType() constant returns(bytes)
func (_ValidatorSet *ValidatorSetCaller) RoundType(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "roundType")
	return *ret0, err
}

// RoundType is a free data retrieval call binding the contract method 0x2c2d1a3b.
//
// Solidity: function roundType() constant returns(bytes)
func (_ValidatorSet *ValidatorSetSession) RoundType() ([]byte, error) {
	return _ValidatorSet.Contract.RoundType(&_ValidatorSet.CallOpts)
}

// RoundType is a free data retrieval call binding the contract method 0x2c2d1a3b.
//
// Solidity: function roundType() constant returns(bytes)
func (_ValidatorSet *ValidatorSetCallerSession) RoundType() ([]byte, error) {
	return _ValidatorSet.Contract.RoundType(&_ValidatorSet.CallOpts)
}

// TotalVotingPower is a free data retrieval call binding the contract method 0x671b3793.
//
// Solidity: function totalVotingPower() constant returns(uint256)
func (_ValidatorSet *ValidatorSetCaller) TotalVotingPower(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "totalVotingPower")
	return *ret0, err
}

// TotalVotingPower is a free data retrieval call binding the contract method 0x671b3793.
//
// Solidity: function totalVotingPower() constant returns(uint256)
func (_ValidatorSet *ValidatorSetSession) TotalVotingPower() (*big.Int, error) {
	return _ValidatorSet.Contract.TotalVotingPower(&_ValidatorSet.CallOpts)
}

// TotalVotingPower is a free data retrieval call binding the contract method 0x671b3793.
//
// Solidity: function totalVotingPower() constant returns(uint256)
func (_ValidatorSet *ValidatorSetCallerSession) TotalVotingPower() (*big.Int, error) {
	return _ValidatorSet.Contract.TotalVotingPower(&_ValidatorSet.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(votingPower uint256, validator address, pubkey string)
func (_ValidatorSet *ValidatorSetCaller) Validators(opts *bind.CallOpts, arg0 *big.Int) (struct {
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
	err := _ValidatorSet.contract.Call(opts, out, "validators", arg0)
	return *ret, err
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(votingPower uint256, validator address, pubkey string)
func (_ValidatorSet *ValidatorSetSession) Validators(arg0 *big.Int) (struct {
	VotingPower *big.Int
	Validator   common.Address
	Pubkey      string
}, error) {
	return _ValidatorSet.Contract.Validators(&_ValidatorSet.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators( uint256) constant returns(votingPower uint256, validator address, pubkey string)
func (_ValidatorSet *ValidatorSetCallerSession) Validators(arg0 *big.Int) (struct {
	VotingPower *big.Int
	Validator   common.Address
	Pubkey      string
}, error) {
	return _ValidatorSet.Contract.Validators(&_ValidatorSet.CallOpts, arg0)
}

// VoteType is a free data retrieval call binding the contract method 0x7d1a3d37.
//
// Solidity: function voteType() constant returns(bytes)
func (_ValidatorSet *ValidatorSetCaller) VoteType(opts *bind.CallOpts) ([]byte, error) {
	var (
		ret0 = new([]byte)
	)
	out := ret0
	err := _ValidatorSet.contract.Call(opts, out, "voteType")
	return *ret0, err
}

// VoteType is a free data retrieval call binding the contract method 0x7d1a3d37.
//
// Solidity: function voteType() constant returns(bytes)
func (_ValidatorSet *ValidatorSetSession) VoteType() ([]byte, error) {
	return _ValidatorSet.Contract.VoteType(&_ValidatorSet.CallOpts)
}

// VoteType is a free data retrieval call binding the contract method 0x7d1a3d37.
//
// Solidity: function voteType() constant returns(bytes)
func (_ValidatorSet *ValidatorSetCallerSession) VoteType() ([]byte, error) {
	return _ValidatorSet.Contract.VoteType(&_ValidatorSet.CallOpts)
}

// AddValidator is a paid mutator transaction binding the contract method 0x01736c35.
//
// Solidity: function addValidator(validator address, votingPower uint256, _pubkey string) returns()
func (_ValidatorSet *ValidatorSetTransactor) AddValidator(opts *bind.TransactOpts, validator common.Address, votingPower *big.Int, _pubkey string) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "addValidator", validator, votingPower, _pubkey)
}

// AddValidator is a paid mutator transaction binding the contract method 0x01736c35.
//
// Solidity: function addValidator(validator address, votingPower uint256, _pubkey string) returns()
func (_ValidatorSet *ValidatorSetSession) AddValidator(validator common.Address, votingPower *big.Int, _pubkey string) (*types.Transaction, error) {
	return _ValidatorSet.Contract.AddValidator(&_ValidatorSet.TransactOpts, validator, votingPower, _pubkey)
}

// AddValidator is a paid mutator transaction binding the contract method 0x01736c35.
//
// Solidity: function addValidator(validator address, votingPower uint256, _pubkey string) returns()
func (_ValidatorSet *ValidatorSetTransactorSession) AddValidator(validator common.Address, votingPower *big.Int, _pubkey string) (*types.Transaction, error) {
	return _ValidatorSet.Contract.AddValidator(&_ValidatorSet.TransactOpts, validator, votingPower, _pubkey)
}

// GetSha256 is a paid mutator transaction binding the contract method 0x8b053758.
//
// Solidity: function getSha256(input bytes) returns(bytes20)
func (_ValidatorSet *ValidatorSetTransactor) GetSha256(opts *bind.TransactOpts, input []byte) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "getSha256", input)
}

// GetSha256 is a paid mutator transaction binding the contract method 0x8b053758.
//
// Solidity: function getSha256(input bytes) returns(bytes20)
func (_ValidatorSet *ValidatorSetSession) GetSha256(input []byte) (*types.Transaction, error) {
	return _ValidatorSet.Contract.GetSha256(&_ValidatorSet.TransactOpts, input)
}

// GetSha256 is a paid mutator transaction binding the contract method 0x8b053758.
//
// Solidity: function getSha256(input bytes) returns(bytes20)
func (_ValidatorSet *ValidatorSetTransactorSession) GetSha256(input []byte) (*types.Transaction, error) {
	return _ValidatorSet.Contract.GetSha256(&_ValidatorSet.TransactOpts, input)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0xf94e1867.
//
// Solidity: function removeValidator(_index uint256) returns()
func (_ValidatorSet *ValidatorSetTransactor) RemoveValidator(opts *bind.TransactOpts, _index *big.Int) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "removeValidator", _index)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0xf94e1867.
//
// Solidity: function removeValidator(_index uint256) returns()
func (_ValidatorSet *ValidatorSetSession) RemoveValidator(_index *big.Int) (*types.Transaction, error) {
	return _ValidatorSet.Contract.RemoveValidator(&_ValidatorSet.TransactOpts, _index)
}

// RemoveValidator is a paid mutator transaction binding the contract method 0xf94e1867.
//
// Solidity: function removeValidator(_index uint256) returns()
func (_ValidatorSet *ValidatorSetTransactorSession) RemoveValidator(_index *big.Int) (*types.Transaction, error) {
	return _ValidatorSet.Contract.RemoveValidator(&_ValidatorSet.TransactOpts, _index)
}

// SetChainId is a paid mutator transaction binding the contract method 0x7973be07.
//
// Solidity: function setChainId(_chainID string) returns()
func (_ValidatorSet *ValidatorSetTransactor) SetChainId(opts *bind.TransactOpts, _chainID string) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "setChainId", _chainID)
}

// SetChainId is a paid mutator transaction binding the contract method 0x7973be07.
//
// Solidity: function setChainId(_chainID string) returns()
func (_ValidatorSet *ValidatorSetSession) SetChainId(_chainID string) (*types.Transaction, error) {
	return _ValidatorSet.Contract.SetChainId(&_ValidatorSet.TransactOpts, _chainID)
}

// SetChainId is a paid mutator transaction binding the contract method 0x7973be07.
//
// Solidity: function setChainId(_chainID string) returns()
func (_ValidatorSet *ValidatorSetTransactorSession) SetChainId(_chainID string) (*types.Transaction, error) {
	return _ValidatorSet.Contract.SetChainId(&_ValidatorSet.TransactOpts, _chainID)
}

// Validate is a paid mutator transaction binding the contract method 0x2d16c59c.
//
// Solidity: function validate(vote bytes, sigs bytes, extradata bytes) returns(address, address, uint256)
func (_ValidatorSet *ValidatorSetTransactor) Validate(opts *bind.TransactOpts, vote []byte, sigs []byte, extradata []byte) (*types.Transaction, error) {
	return _ValidatorSet.contract.Transact(opts, "validate", vote, sigs, extradata)
}

// Validate is a paid mutator transaction binding the contract method 0x2d16c59c.
//
// Solidity: function validate(vote bytes, sigs bytes, extradata bytes) returns(address, address, uint256)
func (_ValidatorSet *ValidatorSetSession) Validate(vote []byte, sigs []byte, extradata []byte) (*types.Transaction, error) {
	return _ValidatorSet.Contract.Validate(&_ValidatorSet.TransactOpts, vote, sigs, extradata)
}

// Validate is a paid mutator transaction binding the contract method 0x2d16c59c.
//
// Solidity: function validate(vote bytes, sigs bytes, extradata bytes) returns(address, address, uint256)
func (_ValidatorSet *ValidatorSetTransactorSession) Validate(vote []byte, sigs []byte, extradata []byte) (*types.Transaction, error) {
	return _ValidatorSet.Contract.Validate(&_ValidatorSet.TransactOpts, vote, sigs, extradata)
}

// ValidatorSetNewProposerIterator is returned from FilterNewProposer and is used to iterate over the raw logs and unpacked data for NewProposer events raised by the ValidatorSet contract.
type ValidatorSetNewProposerIterator struct {
	Event *ValidatorSetNewProposer // Event containing the contract specifics and raw log

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
func (it *ValidatorSetNewProposerIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ValidatorSetNewProposer)
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
		it.Event = new(ValidatorSetNewProposer)
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
func (it *ValidatorSetNewProposerIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ValidatorSetNewProposerIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ValidatorSetNewProposer represents a NewProposer event raised by the ValidatorSet contract.
type ValidatorSetNewProposer struct {
	User common.Address
	Data []byte
	Raw  types.Log // Blockchain specific contextual infos
}

// FilterNewProposer is a free log retrieval operation binding the contract event 0x11d5415dca2a9428143948f67030a63ac42cbd3639a4aae25e602fbbf5da38db.
//
// Solidity: e NewProposer(user indexed address, data bytes)
func (_ValidatorSet *ValidatorSetFilterer) FilterNewProposer(opts *bind.FilterOpts, user []common.Address) (*ValidatorSetNewProposerIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ValidatorSet.contract.FilterLogs(opts, "NewProposer", userRule)
	if err != nil {
		return nil, err
	}
	return &ValidatorSetNewProposerIterator{contract: _ValidatorSet.contract, event: "NewProposer", logs: logs, sub: sub}, nil
}

// WatchNewProposer is a free log subscription operation binding the contract event 0x11d5415dca2a9428143948f67030a63ac42cbd3639a4aae25e602fbbf5da38db.
//
// Solidity: e NewProposer(user indexed address, data bytes)
func (_ValidatorSet *ValidatorSetFilterer) WatchNewProposer(opts *bind.WatchOpts, sink chan<- *ValidatorSetNewProposer, user []common.Address) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}

	logs, sub, err := _ValidatorSet.contract.WatchLogs(opts, "NewProposer", userRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ValidatorSetNewProposer)
				if err := _ValidatorSet.contract.UnpackLog(event, "NewProposer", log); err != nil {
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
