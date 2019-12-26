// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package delegationmanager

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

// DelegationmanagerABI is the input ABI used to generate the binding from.
const DelegationmanagerABI = "[{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"delegatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"stakeRewards\",\"type\":\"bool\"}],\"name\":\"reStake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"delegatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"bond\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"delegatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"accumBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"accumSlashedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"accIndex\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"accProof\",\"type\":\"bytes\"}],\"name\":\"unstakeClaim\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MIN_DEPOSIT_SIZE\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"validatorDelegation\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"delegatorId\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"bondAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"delegatedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"commissionRate\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"isUnBonding\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"acceptsDelegation\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"unbondAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"validatorUnstake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"delegators\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"claimedRewards\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"slashedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"bondedTo\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivationEpoch\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"delegatorId\",\"type\":\"uint256\"}],\"name\":\"unBond\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"stake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"registry\",\"outputs\":[{\"internalType\":\"contractRegistry\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalStaked\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"stakerNFT\",\"outputs\":[{\"internalType\":\"contractStaker\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"validatorHopLimit\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"delegatorId\",\"type\":\"uint256\"}],\"name\":\"withdrawRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256[]\",\"name\":\"_delegators\",\"type\":\"uint256[]\"},{\"internalType\":\"uint256\",\"name\":\"slashRate\",\"type\":\"uint256\"}],\"name\":\"slash\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"locked\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"delegatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"accumBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"accumSlashedAmount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"accIndex\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"withdraw\",\"type\":\"bool\"},{\"internalType\":\"bytes\",\"name\":\"accProof\",\"type\":\"bytes\"}],\"name\":\"claimRewards\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"updateCommissionRate\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"lock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"token\",\"outputs\":[{\"internalType\":\"contractIERC20\",\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_registry\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_stakerNFT\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"delegatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"activatonEpoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"delegatorId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"Unstaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"delegatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"deactivationEpoch\",\"type\":\"uint256\"}],\"name\":\"UnstakeInit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"delegatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"Bonding\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"delegatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"UnBonding\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"rate\",\"type\":\"uint256\"}],\"name\":\"UpdateCommission\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"delegatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"oldValidatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"newValidatorId\",\"type\":\"uint256\"}],\"name\":\"ReBonding\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"delegatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"ReStaked\",\"type\":\"event\"}]"

// Delegationmanager is an auto generated Go binding around an Ethereum contract.
type Delegationmanager struct {
	DelegationmanagerCaller     // Read-only binding to the contract
	DelegationmanagerTransactor // Write-only binding to the contract
	DelegationmanagerFilterer   // Log filterer for contract events
}

// DelegationmanagerCaller is an auto generated read-only Go binding around an Ethereum contract.
type DelegationmanagerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationmanagerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DelegationmanagerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationmanagerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DelegationmanagerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DelegationmanagerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DelegationmanagerSession struct {
	Contract     *Delegationmanager // Generic contract binding to set the session for
	CallOpts     bind.CallOpts      // Call options to use throughout this session
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// DelegationmanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DelegationmanagerCallerSession struct {
	Contract *DelegationmanagerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts            // Call options to use throughout this session
}

// DelegationmanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DelegationmanagerTransactorSession struct {
	Contract     *DelegationmanagerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts            // Transaction auth options to use throughout this session
}

// DelegationmanagerRaw is an auto generated low-level Go binding around an Ethereum contract.
type DelegationmanagerRaw struct {
	Contract *Delegationmanager // Generic contract binding to access the raw methods on
}

// DelegationmanagerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DelegationmanagerCallerRaw struct {
	Contract *DelegationmanagerCaller // Generic read-only contract binding to access the raw methods on
}

// DelegationmanagerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DelegationmanagerTransactorRaw struct {
	Contract *DelegationmanagerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDelegationmanager creates a new instance of Delegationmanager, bound to a specific deployed contract.
func NewDelegationmanager(address common.Address, backend bind.ContractBackend) (*Delegationmanager, error) {
	contract, err := bindDelegationmanager(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Delegationmanager{DelegationmanagerCaller: DelegationmanagerCaller{contract: contract}, DelegationmanagerTransactor: DelegationmanagerTransactor{contract: contract}, DelegationmanagerFilterer: DelegationmanagerFilterer{contract: contract}}, nil
}

// NewDelegationmanagerCaller creates a new read-only instance of Delegationmanager, bound to a specific deployed contract.
func NewDelegationmanagerCaller(address common.Address, caller bind.ContractCaller) (*DelegationmanagerCaller, error) {
	contract, err := bindDelegationmanager(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DelegationmanagerCaller{contract: contract}, nil
}

// NewDelegationmanagerTransactor creates a new write-only instance of Delegationmanager, bound to a specific deployed contract.
func NewDelegationmanagerTransactor(address common.Address, transactor bind.ContractTransactor) (*DelegationmanagerTransactor, error) {
	contract, err := bindDelegationmanager(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DelegationmanagerTransactor{contract: contract}, nil
}

// NewDelegationmanagerFilterer creates a new log filterer instance of Delegationmanager, bound to a specific deployed contract.
func NewDelegationmanagerFilterer(address common.Address, filterer bind.ContractFilterer) (*DelegationmanagerFilterer, error) {
	contract, err := bindDelegationmanager(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DelegationmanagerFilterer{contract: contract}, nil
}

// bindDelegationmanager binds a generic wrapper to an already deployed contract.
func bindDelegationmanager(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DelegationmanagerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegationmanager *DelegationmanagerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Delegationmanager.Contract.DelegationmanagerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegationmanager *DelegationmanagerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegationmanager.Contract.DelegationmanagerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegationmanager *DelegationmanagerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegationmanager.Contract.DelegationmanagerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Delegationmanager *DelegationmanagerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Delegationmanager.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Delegationmanager *DelegationmanagerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegationmanager.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Delegationmanager *DelegationmanagerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Delegationmanager.Contract.contract.Transact(opts, method, params...)
}

// MINDEPOSITSIZE is a free data retrieval call binding the contract method 0x26c0817e.
//
// Solidity: function MIN_DEPOSIT_SIZE() constant returns(uint256)
func (_Delegationmanager *DelegationmanagerCaller) MINDEPOSITSIZE(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Delegationmanager.contract.Call(opts, out, "MIN_DEPOSIT_SIZE")
	return *ret0, err
}

// MINDEPOSITSIZE is a free data retrieval call binding the contract method 0x26c0817e.
//
// Solidity: function MIN_DEPOSIT_SIZE() constant returns(uint256)
func (_Delegationmanager *DelegationmanagerSession) MINDEPOSITSIZE() (*big.Int, error) {
	return _Delegationmanager.Contract.MINDEPOSITSIZE(&_Delegationmanager.CallOpts)
}

// MINDEPOSITSIZE is a free data retrieval call binding the contract method 0x26c0817e.
//
// Solidity: function MIN_DEPOSIT_SIZE() constant returns(uint256)
func (_Delegationmanager *DelegationmanagerCallerSession) MINDEPOSITSIZE() (*big.Int, error) {
	return _Delegationmanager.Contract.MINDEPOSITSIZE(&_Delegationmanager.CallOpts)
}

// AcceptsDelegation is a free data retrieval call binding the contract method 0x4f91c702.
//
// Solidity: function acceptsDelegation(uint256 validatorId) constant returns(bool)
func (_Delegationmanager *DelegationmanagerCaller) AcceptsDelegation(opts *bind.CallOpts, validatorId *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Delegationmanager.contract.Call(opts, out, "acceptsDelegation", validatorId)
	return *ret0, err
}

// AcceptsDelegation is a free data retrieval call binding the contract method 0x4f91c702.
//
// Solidity: function acceptsDelegation(uint256 validatorId) constant returns(bool)
func (_Delegationmanager *DelegationmanagerSession) AcceptsDelegation(validatorId *big.Int) (bool, error) {
	return _Delegationmanager.Contract.AcceptsDelegation(&_Delegationmanager.CallOpts, validatorId)
}

// AcceptsDelegation is a free data retrieval call binding the contract method 0x4f91c702.
//
// Solidity: function acceptsDelegation(uint256 validatorId) constant returns(bool)
func (_Delegationmanager *DelegationmanagerCallerSession) AcceptsDelegation(validatorId *big.Int) (bool, error) {
	return _Delegationmanager.Contract.AcceptsDelegation(&_Delegationmanager.CallOpts, validatorId)
}

// Delegators is a free data retrieval call binding the contract method 0x5be612c7.
//
// Solidity: function delegators(uint256 ) constant returns(uint256 amount, uint256 reward, uint256 claimedRewards, uint256 slashedAmount, uint256 bondedTo, uint256 deactivationEpoch)
func (_Delegationmanager *DelegationmanagerCaller) Delegators(opts *bind.CallOpts, arg0 *big.Int) (struct {
	Amount            *big.Int
	Reward            *big.Int
	ClaimedRewards    *big.Int
	SlashedAmount     *big.Int
	BondedTo          *big.Int
	DeactivationEpoch *big.Int
}, error) {
	ret := new(struct {
		Amount            *big.Int
		Reward            *big.Int
		ClaimedRewards    *big.Int
		SlashedAmount     *big.Int
		BondedTo          *big.Int
		DeactivationEpoch *big.Int
	})
	out := ret
	err := _Delegationmanager.contract.Call(opts, out, "delegators", arg0)
	return *ret, err
}

// Delegators is a free data retrieval call binding the contract method 0x5be612c7.
//
// Solidity: function delegators(uint256 ) constant returns(uint256 amount, uint256 reward, uint256 claimedRewards, uint256 slashedAmount, uint256 bondedTo, uint256 deactivationEpoch)
func (_Delegationmanager *DelegationmanagerSession) Delegators(arg0 *big.Int) (struct {
	Amount            *big.Int
	Reward            *big.Int
	ClaimedRewards    *big.Int
	SlashedAmount     *big.Int
	BondedTo          *big.Int
	DeactivationEpoch *big.Int
}, error) {
	return _Delegationmanager.Contract.Delegators(&_Delegationmanager.CallOpts, arg0)
}

// Delegators is a free data retrieval call binding the contract method 0x5be612c7.
//
// Solidity: function delegators(uint256 ) constant returns(uint256 amount, uint256 reward, uint256 claimedRewards, uint256 slashedAmount, uint256 bondedTo, uint256 deactivationEpoch)
func (_Delegationmanager *DelegationmanagerCallerSession) Delegators(arg0 *big.Int) (struct {
	Amount            *big.Int
	Reward            *big.Int
	ClaimedRewards    *big.Int
	SlashedAmount     *big.Int
	BondedTo          *big.Int
	DeactivationEpoch *big.Int
}, error) {
	return _Delegationmanager.Contract.Delegators(&_Delegationmanager.CallOpts, arg0)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Delegationmanager *DelegationmanagerCaller) IsOwner(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Delegationmanager.contract.Call(opts, out, "isOwner")
	return *ret0, err
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Delegationmanager *DelegationmanagerSession) IsOwner() (bool, error) {
	return _Delegationmanager.Contract.IsOwner(&_Delegationmanager.CallOpts)
}

// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
//
// Solidity: function isOwner() constant returns(bool)
func (_Delegationmanager *DelegationmanagerCallerSession) IsOwner() (bool, error) {
	return _Delegationmanager.Contract.IsOwner(&_Delegationmanager.CallOpts)
}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() constant returns(bool)
func (_Delegationmanager *DelegationmanagerCaller) Locked(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Delegationmanager.contract.Call(opts, out, "locked")
	return *ret0, err
}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() constant returns(bool)
func (_Delegationmanager *DelegationmanagerSession) Locked() (bool, error) {
	return _Delegationmanager.Contract.Locked(&_Delegationmanager.CallOpts)
}

// Locked is a free data retrieval call binding the contract method 0xcf309012.
//
// Solidity: function locked() constant returns(bool)
func (_Delegationmanager *DelegationmanagerCallerSession) Locked() (bool, error) {
	return _Delegationmanager.Contract.Locked(&_Delegationmanager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Delegationmanager *DelegationmanagerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Delegationmanager.contract.Call(opts, out, "owner")
	return *ret0, err
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Delegationmanager *DelegationmanagerSession) Owner() (common.Address, error) {
	return _Delegationmanager.Contract.Owner(&_Delegationmanager.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() constant returns(address)
func (_Delegationmanager *DelegationmanagerCallerSession) Owner() (common.Address, error) {
	return _Delegationmanager.Contract.Owner(&_Delegationmanager.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() constant returns(address)
func (_Delegationmanager *DelegationmanagerCaller) Registry(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Delegationmanager.contract.Call(opts, out, "registry")
	return *ret0, err
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() constant returns(address)
func (_Delegationmanager *DelegationmanagerSession) Registry() (common.Address, error) {
	return _Delegationmanager.Contract.Registry(&_Delegationmanager.CallOpts)
}

// Registry is a free data retrieval call binding the contract method 0x7b103999.
//
// Solidity: function registry() constant returns(address)
func (_Delegationmanager *DelegationmanagerCallerSession) Registry() (common.Address, error) {
	return _Delegationmanager.Contract.Registry(&_Delegationmanager.CallOpts)
}

// StakerNFT is a free data retrieval call binding the contract method 0x881a9d37.
//
// Solidity: function stakerNFT() constant returns(address)
func (_Delegationmanager *DelegationmanagerCaller) StakerNFT(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Delegationmanager.contract.Call(opts, out, "stakerNFT")
	return *ret0, err
}

// StakerNFT is a free data retrieval call binding the contract method 0x881a9d37.
//
// Solidity: function stakerNFT() constant returns(address)
func (_Delegationmanager *DelegationmanagerSession) StakerNFT() (common.Address, error) {
	return _Delegationmanager.Contract.StakerNFT(&_Delegationmanager.CallOpts)
}

// StakerNFT is a free data retrieval call binding the contract method 0x881a9d37.
//
// Solidity: function stakerNFT() constant returns(address)
func (_Delegationmanager *DelegationmanagerCallerSession) StakerNFT() (common.Address, error) {
	return _Delegationmanager.Contract.StakerNFT(&_Delegationmanager.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Delegationmanager *DelegationmanagerCaller) Token(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Delegationmanager.contract.Call(opts, out, "token")
	return *ret0, err
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Delegationmanager *DelegationmanagerSession) Token() (common.Address, error) {
	return _Delegationmanager.Contract.Token(&_Delegationmanager.CallOpts)
}

// Token is a free data retrieval call binding the contract method 0xfc0c546a.
//
// Solidity: function token() constant returns(address)
func (_Delegationmanager *DelegationmanagerCallerSession) Token() (common.Address, error) {
	return _Delegationmanager.Contract.Token(&_Delegationmanager.CallOpts)
}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() constant returns(uint256)
func (_Delegationmanager *DelegationmanagerCaller) TotalStaked(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Delegationmanager.contract.Call(opts, out, "totalStaked")
	return *ret0, err
}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() constant returns(uint256)
func (_Delegationmanager *DelegationmanagerSession) TotalStaked() (*big.Int, error) {
	return _Delegationmanager.Contract.TotalStaked(&_Delegationmanager.CallOpts)
}

// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
//
// Solidity: function totalStaked() constant returns(uint256)
func (_Delegationmanager *DelegationmanagerCallerSession) TotalStaked() (*big.Int, error) {
	return _Delegationmanager.Contract.TotalStaked(&_Delegationmanager.CallOpts)
}

// ValidatorDelegation is a free data retrieval call binding the contract method 0x2768cf73.
//
// Solidity: function validatorDelegation(uint256 validatorId) constant returns(uint256)
func (_Delegationmanager *DelegationmanagerCaller) ValidatorDelegation(opts *bind.CallOpts, validatorId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Delegationmanager.contract.Call(opts, out, "validatorDelegation", validatorId)
	return *ret0, err
}

// ValidatorDelegation is a free data retrieval call binding the contract method 0x2768cf73.
//
// Solidity: function validatorDelegation(uint256 validatorId) constant returns(uint256)
func (_Delegationmanager *DelegationmanagerSession) ValidatorDelegation(validatorId *big.Int) (*big.Int, error) {
	return _Delegationmanager.Contract.ValidatorDelegation(&_Delegationmanager.CallOpts, validatorId)
}

// ValidatorDelegation is a free data retrieval call binding the contract method 0x2768cf73.
//
// Solidity: function validatorDelegation(uint256 validatorId) constant returns(uint256)
func (_Delegationmanager *DelegationmanagerCallerSession) ValidatorDelegation(validatorId *big.Int) (*big.Int, error) {
	return _Delegationmanager.Contract.ValidatorDelegation(&_Delegationmanager.CallOpts, validatorId)
}

// ValidatorHopLimit is a free data retrieval call binding the contract method 0x8c9c856e.
//
// Solidity: function validatorHopLimit() constant returns(uint256)
func (_Delegationmanager *DelegationmanagerCaller) ValidatorHopLimit(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Delegationmanager.contract.Call(opts, out, "validatorHopLimit")
	return *ret0, err
}

// ValidatorHopLimit is a free data retrieval call binding the contract method 0x8c9c856e.
//
// Solidity: function validatorHopLimit() constant returns(uint256)
func (_Delegationmanager *DelegationmanagerSession) ValidatorHopLimit() (*big.Int, error) {
	return _Delegationmanager.Contract.ValidatorHopLimit(&_Delegationmanager.CallOpts)
}

// ValidatorHopLimit is a free data retrieval call binding the contract method 0x8c9c856e.
//
// Solidity: function validatorHopLimit() constant returns(uint256)
func (_Delegationmanager *DelegationmanagerCallerSession) ValidatorHopLimit() (*big.Int, error) {
	return _Delegationmanager.Contract.ValidatorHopLimit(&_Delegationmanager.CallOpts)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) constant returns(uint256 delegatedAmount, uint256 commissionRate, bool isUnBonding)
func (_Delegationmanager *DelegationmanagerCaller) Validators(opts *bind.CallOpts, arg0 *big.Int) (struct {
	DelegatedAmount *big.Int
	CommissionRate  *big.Int
	IsUnBonding     bool
}, error) {
	ret := new(struct {
		DelegatedAmount *big.Int
		CommissionRate  *big.Int
		IsUnBonding     bool
	})
	out := ret
	err := _Delegationmanager.contract.Call(opts, out, "validators", arg0)
	return *ret, err
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) constant returns(uint256 delegatedAmount, uint256 commissionRate, bool isUnBonding)
func (_Delegationmanager *DelegationmanagerSession) Validators(arg0 *big.Int) (struct {
	DelegatedAmount *big.Int
	CommissionRate  *big.Int
	IsUnBonding     bool
}, error) {
	return _Delegationmanager.Contract.Validators(&_Delegationmanager.CallOpts, arg0)
}

// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
//
// Solidity: function validators(uint256 ) constant returns(uint256 delegatedAmount, uint256 commissionRate, bool isUnBonding)
func (_Delegationmanager *DelegationmanagerCallerSession) Validators(arg0 *big.Int) (struct {
	DelegatedAmount *big.Int
	CommissionRate  *big.Int
	IsUnBonding     bool
}, error) {
	return _Delegationmanager.Contract.Validators(&_Delegationmanager.CallOpts, arg0)
}

// Bond is a paid mutator transaction binding the contract method 0x0ba74b2f.
//
// Solidity: function bond(uint256 delegatorId, uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactor) Bond(opts *bind.TransactOpts, delegatorId *big.Int, validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "bond", delegatorId, validatorId)
}

// Bond is a paid mutator transaction binding the contract method 0x0ba74b2f.
//
// Solidity: function bond(uint256 delegatorId, uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerSession) Bond(delegatorId *big.Int, validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.Bond(&_Delegationmanager.TransactOpts, delegatorId, validatorId)
}

// Bond is a paid mutator transaction binding the contract method 0x0ba74b2f.
//
// Solidity: function bond(uint256 delegatorId, uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) Bond(delegatorId *big.Int, validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.Bond(&_Delegationmanager.TransactOpts, delegatorId, validatorId)
}

// BondAll is a paid mutator transaction binding the contract method 0x2f3f4490.
//
// Solidity: function bondAll(uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactor) BondAll(opts *bind.TransactOpts, validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "bondAll", validatorId)
}

// BondAll is a paid mutator transaction binding the contract method 0x2f3f4490.
//
// Solidity: function bondAll(uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerSession) BondAll(validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.BondAll(&_Delegationmanager.TransactOpts, validatorId)
}

// BondAll is a paid mutator transaction binding the contract method 0x2f3f4490.
//
// Solidity: function bondAll(uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) BondAll(validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.BondAll(&_Delegationmanager.TransactOpts, validatorId)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0xd02cfa56.
//
// Solidity: function claimRewards(uint256 delegatorId, uint256 accumBalance, uint256 accumSlashedAmount, uint256 accIndex, bool withdraw, bytes accProof) returns()
func (_Delegationmanager *DelegationmanagerTransactor) ClaimRewards(opts *bind.TransactOpts, delegatorId *big.Int, accumBalance *big.Int, accumSlashedAmount *big.Int, accIndex *big.Int, withdraw bool, accProof []byte) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "claimRewards", delegatorId, accumBalance, accumSlashedAmount, accIndex, withdraw, accProof)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0xd02cfa56.
//
// Solidity: function claimRewards(uint256 delegatorId, uint256 accumBalance, uint256 accumSlashedAmount, uint256 accIndex, bool withdraw, bytes accProof) returns()
func (_Delegationmanager *DelegationmanagerSession) ClaimRewards(delegatorId *big.Int, accumBalance *big.Int, accumSlashedAmount *big.Int, accIndex *big.Int, withdraw bool, accProof []byte) (*types.Transaction, error) {
	return _Delegationmanager.Contract.ClaimRewards(&_Delegationmanager.TransactOpts, delegatorId, accumBalance, accumSlashedAmount, accIndex, withdraw, accProof)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0xd02cfa56.
//
// Solidity: function claimRewards(uint256 delegatorId, uint256 accumBalance, uint256 accumSlashedAmount, uint256 accIndex, bool withdraw, bytes accProof) returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) ClaimRewards(delegatorId *big.Int, accumBalance *big.Int, accumSlashedAmount *big.Int, accIndex *big.Int, withdraw bool, accProof []byte) (*types.Transaction, error) {
	return _Delegationmanager.Contract.ClaimRewards(&_Delegationmanager.TransactOpts, delegatorId, accumBalance, accumSlashedAmount, accIndex, withdraw, accProof)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_Delegationmanager *DelegationmanagerTransactor) Lock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "lock")
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_Delegationmanager *DelegationmanagerSession) Lock() (*types.Transaction, error) {
	return _Delegationmanager.Contract.Lock(&_Delegationmanager.TransactOpts)
}

// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
//
// Solidity: function lock() returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) Lock() (*types.Transaction, error) {
	return _Delegationmanager.Contract.Lock(&_Delegationmanager.TransactOpts)
}

// ReStake is a paid mutator transaction binding the contract method 0x0b350d03.
//
// Solidity: function reStake(uint256 delegatorId, uint256 amount, bool stakeRewards) returns()
func (_Delegationmanager *DelegationmanagerTransactor) ReStake(opts *bind.TransactOpts, delegatorId *big.Int, amount *big.Int, stakeRewards bool) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "reStake", delegatorId, amount, stakeRewards)
}

// ReStake is a paid mutator transaction binding the contract method 0x0b350d03.
//
// Solidity: function reStake(uint256 delegatorId, uint256 amount, bool stakeRewards) returns()
func (_Delegationmanager *DelegationmanagerSession) ReStake(delegatorId *big.Int, amount *big.Int, stakeRewards bool) (*types.Transaction, error) {
	return _Delegationmanager.Contract.ReStake(&_Delegationmanager.TransactOpts, delegatorId, amount, stakeRewards)
}

// ReStake is a paid mutator transaction binding the contract method 0x0b350d03.
//
// Solidity: function reStake(uint256 delegatorId, uint256 amount, bool stakeRewards) returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) ReStake(delegatorId *big.Int, amount *big.Int, stakeRewards bool) (*types.Transaction, error) {
	return _Delegationmanager.Contract.ReStake(&_Delegationmanager.TransactOpts, delegatorId, amount, stakeRewards)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Delegationmanager *DelegationmanagerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Delegationmanager *DelegationmanagerSession) RenounceOwnership() (*types.Transaction, error) {
	return _Delegationmanager.Contract.RenounceOwnership(&_Delegationmanager.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _Delegationmanager.Contract.RenounceOwnership(&_Delegationmanager.TransactOpts)
}

// Slash is a paid mutator transaction binding the contract method 0xce5d67d9.
//
// Solidity: function slash(uint256[] _delegators, uint256 slashRate) returns()
func (_Delegationmanager *DelegationmanagerTransactor) Slash(opts *bind.TransactOpts, _delegators []*big.Int, slashRate *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "slash", _delegators, slashRate)
}

// Slash is a paid mutator transaction binding the contract method 0xce5d67d9.
//
// Solidity: function slash(uint256[] _delegators, uint256 slashRate) returns()
func (_Delegationmanager *DelegationmanagerSession) Slash(_delegators []*big.Int, slashRate *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.Slash(&_Delegationmanager.TransactOpts, _delegators, slashRate)
}

// Slash is a paid mutator transaction binding the contract method 0xce5d67d9.
//
// Solidity: function slash(uint256[] _delegators, uint256 slashRate) returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) Slash(_delegators []*big.Int, slashRate *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.Slash(&_Delegationmanager.TransactOpts, _delegators, slashRate)
}

// Stake is a paid mutator transaction binding the contract method 0x7b0472f0.
//
// Solidity: function stake(uint256 amount, uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactor) Stake(opts *bind.TransactOpts, amount *big.Int, validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "stake", amount, validatorId)
}

// Stake is a paid mutator transaction binding the contract method 0x7b0472f0.
//
// Solidity: function stake(uint256 amount, uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerSession) Stake(amount *big.Int, validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.Stake(&_Delegationmanager.TransactOpts, amount, validatorId)
}

// Stake is a paid mutator transaction binding the contract method 0x7b0472f0.
//
// Solidity: function stake(uint256 amount, uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) Stake(amount *big.Int, validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.Stake(&_Delegationmanager.TransactOpts, amount, validatorId)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Delegationmanager *DelegationmanagerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Delegationmanager *DelegationmanagerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Delegationmanager.Contract.TransferOwnership(&_Delegationmanager.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Delegationmanager.Contract.TransferOwnership(&_Delegationmanager.TransactOpts, newOwner)
}

// UnBond is a paid mutator transaction binding the contract method 0x7865d72f.
//
// Solidity: function unBond(uint256 delegatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactor) UnBond(opts *bind.TransactOpts, delegatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "unBond", delegatorId)
}

// UnBond is a paid mutator transaction binding the contract method 0x7865d72f.
//
// Solidity: function unBond(uint256 delegatorId) returns()
func (_Delegationmanager *DelegationmanagerSession) UnBond(delegatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.UnBond(&_Delegationmanager.TransactOpts, delegatorId)
}

// UnBond is a paid mutator transaction binding the contract method 0x7865d72f.
//
// Solidity: function unBond(uint256 delegatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) UnBond(delegatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.UnBond(&_Delegationmanager.TransactOpts, delegatorId)
}

// UnbondAll is a paid mutator transaction binding the contract method 0x5546a520.
//
// Solidity: function unbondAll(uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactor) UnbondAll(opts *bind.TransactOpts, validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "unbondAll", validatorId)
}

// UnbondAll is a paid mutator transaction binding the contract method 0x5546a520.
//
// Solidity: function unbondAll(uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerSession) UnbondAll(validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.UnbondAll(&_Delegationmanager.TransactOpts, validatorId)
}

// UnbondAll is a paid mutator transaction binding the contract method 0x5546a520.
//
// Solidity: function unbondAll(uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) UnbondAll(validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.UnbondAll(&_Delegationmanager.TransactOpts, validatorId)
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_Delegationmanager *DelegationmanagerTransactor) Unlock(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "unlock")
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_Delegationmanager *DelegationmanagerSession) Unlock() (*types.Transaction, error) {
	return _Delegationmanager.Contract.Unlock(&_Delegationmanager.TransactOpts)
}

// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
//
// Solidity: function unlock() returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) Unlock() (*types.Transaction, error) {
	return _Delegationmanager.Contract.Unlock(&_Delegationmanager.TransactOpts)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 delegatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactor) Unstake(opts *bind.TransactOpts, delegatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "unstake", delegatorId)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 delegatorId) returns()
func (_Delegationmanager *DelegationmanagerSession) Unstake(delegatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.Unstake(&_Delegationmanager.TransactOpts, delegatorId)
}

// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
//
// Solidity: function unstake(uint256 delegatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) Unstake(delegatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.Unstake(&_Delegationmanager.TransactOpts, delegatorId)
}

// UnstakeClaim is a paid mutator transaction binding the contract method 0x22f05fe7.
//
// Solidity: function unstakeClaim(uint256 delegatorId, uint256 accumBalance, uint256 accumSlashedAmount, uint256 accIndex, bytes accProof) returns()
func (_Delegationmanager *DelegationmanagerTransactor) UnstakeClaim(opts *bind.TransactOpts, delegatorId *big.Int, accumBalance *big.Int, accumSlashedAmount *big.Int, accIndex *big.Int, accProof []byte) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "unstakeClaim", delegatorId, accumBalance, accumSlashedAmount, accIndex, accProof)
}

// UnstakeClaim is a paid mutator transaction binding the contract method 0x22f05fe7.
//
// Solidity: function unstakeClaim(uint256 delegatorId, uint256 accumBalance, uint256 accumSlashedAmount, uint256 accIndex, bytes accProof) returns()
func (_Delegationmanager *DelegationmanagerSession) UnstakeClaim(delegatorId *big.Int, accumBalance *big.Int, accumSlashedAmount *big.Int, accIndex *big.Int, accProof []byte) (*types.Transaction, error) {
	return _Delegationmanager.Contract.UnstakeClaim(&_Delegationmanager.TransactOpts, delegatorId, accumBalance, accumSlashedAmount, accIndex, accProof)
}

// UnstakeClaim is a paid mutator transaction binding the contract method 0x22f05fe7.
//
// Solidity: function unstakeClaim(uint256 delegatorId, uint256 accumBalance, uint256 accumSlashedAmount, uint256 accIndex, bytes accProof) returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) UnstakeClaim(delegatorId *big.Int, accumBalance *big.Int, accumSlashedAmount *big.Int, accIndex *big.Int, accProof []byte) (*types.Transaction, error) {
	return _Delegationmanager.Contract.UnstakeClaim(&_Delegationmanager.TransactOpts, delegatorId, accumBalance, accumSlashedAmount, accIndex, accProof)
}

// UpdateCommissionRate is a paid mutator transaction binding the contract method 0xdcd962b2.
//
// Solidity: function updateCommissionRate(uint256 validatorId, uint256 rate) returns()
func (_Delegationmanager *DelegationmanagerTransactor) UpdateCommissionRate(opts *bind.TransactOpts, validatorId *big.Int, rate *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "updateCommissionRate", validatorId, rate)
}

// UpdateCommissionRate is a paid mutator transaction binding the contract method 0xdcd962b2.
//
// Solidity: function updateCommissionRate(uint256 validatorId, uint256 rate) returns()
func (_Delegationmanager *DelegationmanagerSession) UpdateCommissionRate(validatorId *big.Int, rate *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.UpdateCommissionRate(&_Delegationmanager.TransactOpts, validatorId, rate)
}

// UpdateCommissionRate is a paid mutator transaction binding the contract method 0xdcd962b2.
//
// Solidity: function updateCommissionRate(uint256 validatorId, uint256 rate) returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) UpdateCommissionRate(validatorId *big.Int, rate *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.UpdateCommissionRate(&_Delegationmanager.TransactOpts, validatorId, rate)
}

// ValidatorUnstake is a paid mutator transaction binding the contract method 0x56aff179.
//
// Solidity: function validatorUnstake(uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactor) ValidatorUnstake(opts *bind.TransactOpts, validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "validatorUnstake", validatorId)
}

// ValidatorUnstake is a paid mutator transaction binding the contract method 0x56aff179.
//
// Solidity: function validatorUnstake(uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerSession) ValidatorUnstake(validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.ValidatorUnstake(&_Delegationmanager.TransactOpts, validatorId)
}

// ValidatorUnstake is a paid mutator transaction binding the contract method 0x56aff179.
//
// Solidity: function validatorUnstake(uint256 validatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) ValidatorUnstake(validatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.ValidatorUnstake(&_Delegationmanager.TransactOpts, validatorId)
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0x9342c8f4.
//
// Solidity: function withdrawRewards(uint256 delegatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactor) WithdrawRewards(opts *bind.TransactOpts, delegatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.contract.Transact(opts, "withdrawRewards", delegatorId)
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0x9342c8f4.
//
// Solidity: function withdrawRewards(uint256 delegatorId) returns()
func (_Delegationmanager *DelegationmanagerSession) WithdrawRewards(delegatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.WithdrawRewards(&_Delegationmanager.TransactOpts, delegatorId)
}

// WithdrawRewards is a paid mutator transaction binding the contract method 0x9342c8f4.
//
// Solidity: function withdrawRewards(uint256 delegatorId) returns()
func (_Delegationmanager *DelegationmanagerTransactorSession) WithdrawRewards(delegatorId *big.Int) (*types.Transaction, error) {
	return _Delegationmanager.Contract.WithdrawRewards(&_Delegationmanager.TransactOpts, delegatorId)
}

// DelegationmanagerBondingIterator is returned from FilterBonding and is used to iterate over the raw logs and unpacked data for Bonding events raised by the Delegationmanager contract.
type DelegationmanagerBondingIterator struct {
	Event *DelegationmanagerBonding // Event containing the contract specifics and raw log

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
func (it *DelegationmanagerBondingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationmanagerBonding)
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
		it.Event = new(DelegationmanagerBonding)
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
func (it *DelegationmanagerBondingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationmanagerBondingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationmanagerBonding represents a Bonding event raised by the Delegationmanager contract.
type DelegationmanagerBonding struct {
	DelegatorId *big.Int
	ValidatorId *big.Int
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBonding is a free log retrieval operation binding the contract event 0xfb99dc5016f0b87350ecd0ea85576ccea807243ca8ba3617f382af921ed3b810.
//
// Solidity: event Bonding(uint256 indexed delegatorId, uint256 indexed validatorId, uint256 indexed amount)
func (_Delegationmanager *DelegationmanagerFilterer) FilterBonding(opts *bind.FilterOpts, delegatorId []*big.Int, validatorId []*big.Int, amount []*big.Int) (*DelegationmanagerBondingIterator, error) {

	var delegatorIdRule []interface{}
	for _, delegatorIdItem := range delegatorId {
		delegatorIdRule = append(delegatorIdRule, delegatorIdItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Delegationmanager.contract.FilterLogs(opts, "Bonding", delegatorIdRule, validatorIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &DelegationmanagerBondingIterator{contract: _Delegationmanager.contract, event: "Bonding", logs: logs, sub: sub}, nil
}

// WatchBonding is a free log subscription operation binding the contract event 0xfb99dc5016f0b87350ecd0ea85576ccea807243ca8ba3617f382af921ed3b810.
//
// Solidity: event Bonding(uint256 indexed delegatorId, uint256 indexed validatorId, uint256 indexed amount)
func (_Delegationmanager *DelegationmanagerFilterer) WatchBonding(opts *bind.WatchOpts, sink chan<- *DelegationmanagerBonding, delegatorId []*big.Int, validatorId []*big.Int, amount []*big.Int) (event.Subscription, error) {

	var delegatorIdRule []interface{}
	for _, delegatorIdItem := range delegatorId {
		delegatorIdRule = append(delegatorIdRule, delegatorIdItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Delegationmanager.contract.WatchLogs(opts, "Bonding", delegatorIdRule, validatorIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationmanagerBonding)
				if err := _Delegationmanager.contract.UnpackLog(event, "Bonding", log); err != nil {
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

// ParseBonding is a log parse operation binding the contract event 0xfb99dc5016f0b87350ecd0ea85576ccea807243ca8ba3617f382af921ed3b810.
//
// Solidity: event Bonding(uint256 indexed delegatorId, uint256 indexed validatorId, uint256 indexed amount)
func (_Delegationmanager *DelegationmanagerFilterer) ParseBonding(log types.Log) (*DelegationmanagerBonding, error) {
	event := new(DelegationmanagerBonding)
	if err := _Delegationmanager.contract.UnpackLog(event, "Bonding", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DelegationmanagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Delegationmanager contract.
type DelegationmanagerOwnershipTransferredIterator struct {
	Event *DelegationmanagerOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *DelegationmanagerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationmanagerOwnershipTransferred)
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
		it.Event = new(DelegationmanagerOwnershipTransferred)
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
func (it *DelegationmanagerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationmanagerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationmanagerOwnershipTransferred represents a OwnershipTransferred event raised by the Delegationmanager contract.
type DelegationmanagerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Delegationmanager *DelegationmanagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*DelegationmanagerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Delegationmanager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &DelegationmanagerOwnershipTransferredIterator{contract: _Delegationmanager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Delegationmanager *DelegationmanagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *DelegationmanagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _Delegationmanager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationmanagerOwnershipTransferred)
				if err := _Delegationmanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_Delegationmanager *DelegationmanagerFilterer) ParseOwnershipTransferred(log types.Log) (*DelegationmanagerOwnershipTransferred, error) {
	event := new(DelegationmanagerOwnershipTransferred)
	if err := _Delegationmanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DelegationmanagerReBondingIterator is returned from FilterReBonding and is used to iterate over the raw logs and unpacked data for ReBonding events raised by the Delegationmanager contract.
type DelegationmanagerReBondingIterator struct {
	Event *DelegationmanagerReBonding // Event containing the contract specifics and raw log

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
func (it *DelegationmanagerReBondingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationmanagerReBonding)
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
		it.Event = new(DelegationmanagerReBonding)
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
func (it *DelegationmanagerReBondingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationmanagerReBondingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationmanagerReBonding represents a ReBonding event raised by the Delegationmanager contract.
type DelegationmanagerReBonding struct {
	DelegatorId    *big.Int
	OldValidatorId *big.Int
	NewValidatorId *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterReBonding is a free log retrieval operation binding the contract event 0x28b5e6c19771cb360b738da33c22674d7d04f48a9da50ce9a39ec0ea91cb860c.
//
// Solidity: event ReBonding(uint256 indexed delegatorId, uint256 indexed oldValidatorId, uint256 indexed newValidatorId)
func (_Delegationmanager *DelegationmanagerFilterer) FilterReBonding(opts *bind.FilterOpts, delegatorId []*big.Int, oldValidatorId []*big.Int, newValidatorId []*big.Int) (*DelegationmanagerReBondingIterator, error) {

	var delegatorIdRule []interface{}
	for _, delegatorIdItem := range delegatorId {
		delegatorIdRule = append(delegatorIdRule, delegatorIdItem)
	}
	var oldValidatorIdRule []interface{}
	for _, oldValidatorIdItem := range oldValidatorId {
		oldValidatorIdRule = append(oldValidatorIdRule, oldValidatorIdItem)
	}
	var newValidatorIdRule []interface{}
	for _, newValidatorIdItem := range newValidatorId {
		newValidatorIdRule = append(newValidatorIdRule, newValidatorIdItem)
	}

	logs, sub, err := _Delegationmanager.contract.FilterLogs(opts, "ReBonding", delegatorIdRule, oldValidatorIdRule, newValidatorIdRule)
	if err != nil {
		return nil, err
	}
	return &DelegationmanagerReBondingIterator{contract: _Delegationmanager.contract, event: "ReBonding", logs: logs, sub: sub}, nil
}

// WatchReBonding is a free log subscription operation binding the contract event 0x28b5e6c19771cb360b738da33c22674d7d04f48a9da50ce9a39ec0ea91cb860c.
//
// Solidity: event ReBonding(uint256 indexed delegatorId, uint256 indexed oldValidatorId, uint256 indexed newValidatorId)
func (_Delegationmanager *DelegationmanagerFilterer) WatchReBonding(opts *bind.WatchOpts, sink chan<- *DelegationmanagerReBonding, delegatorId []*big.Int, oldValidatorId []*big.Int, newValidatorId []*big.Int) (event.Subscription, error) {

	var delegatorIdRule []interface{}
	for _, delegatorIdItem := range delegatorId {
		delegatorIdRule = append(delegatorIdRule, delegatorIdItem)
	}
	var oldValidatorIdRule []interface{}
	for _, oldValidatorIdItem := range oldValidatorId {
		oldValidatorIdRule = append(oldValidatorIdRule, oldValidatorIdItem)
	}
	var newValidatorIdRule []interface{}
	for _, newValidatorIdItem := range newValidatorId {
		newValidatorIdRule = append(newValidatorIdRule, newValidatorIdItem)
	}

	logs, sub, err := _Delegationmanager.contract.WatchLogs(opts, "ReBonding", delegatorIdRule, oldValidatorIdRule, newValidatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationmanagerReBonding)
				if err := _Delegationmanager.contract.UnpackLog(event, "ReBonding", log); err != nil {
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

// ParseReBonding is a log parse operation binding the contract event 0x28b5e6c19771cb360b738da33c22674d7d04f48a9da50ce9a39ec0ea91cb860c.
//
// Solidity: event ReBonding(uint256 indexed delegatorId, uint256 indexed oldValidatorId, uint256 indexed newValidatorId)
func (_Delegationmanager *DelegationmanagerFilterer) ParseReBonding(log types.Log) (*DelegationmanagerReBonding, error) {
	event := new(DelegationmanagerReBonding)
	if err := _Delegationmanager.contract.UnpackLog(event, "ReBonding", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DelegationmanagerReStakedIterator is returned from FilterReStaked and is used to iterate over the raw logs and unpacked data for ReStaked events raised by the Delegationmanager contract.
type DelegationmanagerReStakedIterator struct {
	Event *DelegationmanagerReStaked // Event containing the contract specifics and raw log

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
func (it *DelegationmanagerReStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationmanagerReStaked)
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
		it.Event = new(DelegationmanagerReStaked)
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
func (it *DelegationmanagerReStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationmanagerReStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationmanagerReStaked represents a ReStaked event raised by the Delegationmanager contract.
type DelegationmanagerReStaked struct {
	DelegatorId *big.Int
	Amount      *big.Int
	Total       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterReStaked is a free log retrieval operation binding the contract event 0x9cc0e589f20d3310eb2ad571b23529003bd46048d0d1af29277dcf0aa3c398ce.
//
// Solidity: event ReStaked(uint256 indexed delegatorId, uint256 indexed amount, uint256 total)
func (_Delegationmanager *DelegationmanagerFilterer) FilterReStaked(opts *bind.FilterOpts, delegatorId []*big.Int, amount []*big.Int) (*DelegationmanagerReStakedIterator, error) {

	var delegatorIdRule []interface{}
	for _, delegatorIdItem := range delegatorId {
		delegatorIdRule = append(delegatorIdRule, delegatorIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Delegationmanager.contract.FilterLogs(opts, "ReStaked", delegatorIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &DelegationmanagerReStakedIterator{contract: _Delegationmanager.contract, event: "ReStaked", logs: logs, sub: sub}, nil
}

// WatchReStaked is a free log subscription operation binding the contract event 0x9cc0e589f20d3310eb2ad571b23529003bd46048d0d1af29277dcf0aa3c398ce.
//
// Solidity: event ReStaked(uint256 indexed delegatorId, uint256 indexed amount, uint256 total)
func (_Delegationmanager *DelegationmanagerFilterer) WatchReStaked(opts *bind.WatchOpts, sink chan<- *DelegationmanagerReStaked, delegatorId []*big.Int, amount []*big.Int) (event.Subscription, error) {

	var delegatorIdRule []interface{}
	for _, delegatorIdItem := range delegatorId {
		delegatorIdRule = append(delegatorIdRule, delegatorIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Delegationmanager.contract.WatchLogs(opts, "ReStaked", delegatorIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationmanagerReStaked)
				if err := _Delegationmanager.contract.UnpackLog(event, "ReStaked", log); err != nil {
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

// ParseReStaked is a log parse operation binding the contract event 0x9cc0e589f20d3310eb2ad571b23529003bd46048d0d1af29277dcf0aa3c398ce.
//
// Solidity: event ReStaked(uint256 indexed delegatorId, uint256 indexed amount, uint256 total)
func (_Delegationmanager *DelegationmanagerFilterer) ParseReStaked(log types.Log) (*DelegationmanagerReStaked, error) {
	event := new(DelegationmanagerReStaked)
	if err := _Delegationmanager.contract.UnpackLog(event, "ReStaked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DelegationmanagerStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the Delegationmanager contract.
type DelegationmanagerStakedIterator struct {
	Event *DelegationmanagerStaked // Event containing the contract specifics and raw log

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
func (it *DelegationmanagerStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationmanagerStaked)
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
		it.Event = new(DelegationmanagerStaked)
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
func (it *DelegationmanagerStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationmanagerStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationmanagerStaked represents a Staked event raised by the Delegationmanager contract.
type DelegationmanagerStaked struct {
	User           common.Address
	DelegatorId    *big.Int
	ActivatonEpoch *big.Int
	Amount         *big.Int
	Total          *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0x9cfd25589d1eb8ad71e342a86a8524e83522e3936c0803048c08f6d9ad974f40.
//
// Solidity: event Staked(address indexed user, uint256 indexed delegatorId, uint256 indexed activatonEpoch, uint256 amount, uint256 total)
func (_Delegationmanager *DelegationmanagerFilterer) FilterStaked(opts *bind.FilterOpts, user []common.Address, delegatorId []*big.Int, activatonEpoch []*big.Int) (*DelegationmanagerStakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var delegatorIdRule []interface{}
	for _, delegatorIdItem := range delegatorId {
		delegatorIdRule = append(delegatorIdRule, delegatorIdItem)
	}
	var activatonEpochRule []interface{}
	for _, activatonEpochItem := range activatonEpoch {
		activatonEpochRule = append(activatonEpochRule, activatonEpochItem)
	}

	logs, sub, err := _Delegationmanager.contract.FilterLogs(opts, "Staked", userRule, delegatorIdRule, activatonEpochRule)
	if err != nil {
		return nil, err
	}
	return &DelegationmanagerStakedIterator{contract: _Delegationmanager.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0x9cfd25589d1eb8ad71e342a86a8524e83522e3936c0803048c08f6d9ad974f40.
//
// Solidity: event Staked(address indexed user, uint256 indexed delegatorId, uint256 indexed activatonEpoch, uint256 amount, uint256 total)
func (_Delegationmanager *DelegationmanagerFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *DelegationmanagerStaked, user []common.Address, delegatorId []*big.Int, activatonEpoch []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var delegatorIdRule []interface{}
	for _, delegatorIdItem := range delegatorId {
		delegatorIdRule = append(delegatorIdRule, delegatorIdItem)
	}
	var activatonEpochRule []interface{}
	for _, activatonEpochItem := range activatonEpoch {
		activatonEpochRule = append(activatonEpochRule, activatonEpochItem)
	}

	logs, sub, err := _Delegationmanager.contract.WatchLogs(opts, "Staked", userRule, delegatorIdRule, activatonEpochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationmanagerStaked)
				if err := _Delegationmanager.contract.UnpackLog(event, "Staked", log); err != nil {
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

// ParseStaked is a log parse operation binding the contract event 0x9cfd25589d1eb8ad71e342a86a8524e83522e3936c0803048c08f6d9ad974f40.
//
// Solidity: event Staked(address indexed user, uint256 indexed delegatorId, uint256 indexed activatonEpoch, uint256 amount, uint256 total)
func (_Delegationmanager *DelegationmanagerFilterer) ParseStaked(log types.Log) (*DelegationmanagerStaked, error) {
	event := new(DelegationmanagerStaked)
	if err := _Delegationmanager.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DelegationmanagerUnBondingIterator is returned from FilterUnBonding and is used to iterate over the raw logs and unpacked data for UnBonding events raised by the Delegationmanager contract.
type DelegationmanagerUnBondingIterator struct {
	Event *DelegationmanagerUnBonding // Event containing the contract specifics and raw log

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
func (it *DelegationmanagerUnBondingIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationmanagerUnBonding)
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
		it.Event = new(DelegationmanagerUnBonding)
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
func (it *DelegationmanagerUnBondingIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationmanagerUnBondingIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationmanagerUnBonding represents a UnBonding event raised by the Delegationmanager contract.
type DelegationmanagerUnBonding struct {
	DelegatorId *big.Int
	ValidatorId *big.Int
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUnBonding is a free log retrieval operation binding the contract event 0x77c7e66fba72362dffdbc1cebc7fd21074316ef497049af2c4ae7340291d2ce3.
//
// Solidity: event UnBonding(uint256 indexed delegatorId, uint256 indexed validatorId, uint256 indexed amount)
func (_Delegationmanager *DelegationmanagerFilterer) FilterUnBonding(opts *bind.FilterOpts, delegatorId []*big.Int, validatorId []*big.Int, amount []*big.Int) (*DelegationmanagerUnBondingIterator, error) {

	var delegatorIdRule []interface{}
	for _, delegatorIdItem := range delegatorId {
		delegatorIdRule = append(delegatorIdRule, delegatorIdItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Delegationmanager.contract.FilterLogs(opts, "UnBonding", delegatorIdRule, validatorIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &DelegationmanagerUnBondingIterator{contract: _Delegationmanager.contract, event: "UnBonding", logs: logs, sub: sub}, nil
}

// WatchUnBonding is a free log subscription operation binding the contract event 0x77c7e66fba72362dffdbc1cebc7fd21074316ef497049af2c4ae7340291d2ce3.
//
// Solidity: event UnBonding(uint256 indexed delegatorId, uint256 indexed validatorId, uint256 indexed amount)
func (_Delegationmanager *DelegationmanagerFilterer) WatchUnBonding(opts *bind.WatchOpts, sink chan<- *DelegationmanagerUnBonding, delegatorId []*big.Int, validatorId []*big.Int, amount []*big.Int) (event.Subscription, error) {

	var delegatorIdRule []interface{}
	for _, delegatorIdItem := range delegatorId {
		delegatorIdRule = append(delegatorIdRule, delegatorIdItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Delegationmanager.contract.WatchLogs(opts, "UnBonding", delegatorIdRule, validatorIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationmanagerUnBonding)
				if err := _Delegationmanager.contract.UnpackLog(event, "UnBonding", log); err != nil {
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

// ParseUnBonding is a log parse operation binding the contract event 0x77c7e66fba72362dffdbc1cebc7fd21074316ef497049af2c4ae7340291d2ce3.
//
// Solidity: event UnBonding(uint256 indexed delegatorId, uint256 indexed validatorId, uint256 indexed amount)
func (_Delegationmanager *DelegationmanagerFilterer) ParseUnBonding(log types.Log) (*DelegationmanagerUnBonding, error) {
	event := new(DelegationmanagerUnBonding)
	if err := _Delegationmanager.contract.UnpackLog(event, "UnBonding", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DelegationmanagerUnstakeInitIterator is returned from FilterUnstakeInit and is used to iterate over the raw logs and unpacked data for UnstakeInit events raised by the Delegationmanager contract.
type DelegationmanagerUnstakeInitIterator struct {
	Event *DelegationmanagerUnstakeInit // Event containing the contract specifics and raw log

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
func (it *DelegationmanagerUnstakeInitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationmanagerUnstakeInit)
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
		it.Event = new(DelegationmanagerUnstakeInit)
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
func (it *DelegationmanagerUnstakeInitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationmanagerUnstakeInitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationmanagerUnstakeInit represents a UnstakeInit event raised by the Delegationmanager contract.
type DelegationmanagerUnstakeInit struct {
	User              common.Address
	DelegatorId       *big.Int
	DeactivationEpoch *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUnstakeInit is a free log retrieval operation binding the contract event 0xcde3813e379342e7506ca1f984065c81821e879aebd4be83cb35b5ab976518f9.
//
// Solidity: event UnstakeInit(address indexed user, uint256 indexed delegatorId, uint256 indexed deactivationEpoch)
func (_Delegationmanager *DelegationmanagerFilterer) FilterUnstakeInit(opts *bind.FilterOpts, user []common.Address, delegatorId []*big.Int, deactivationEpoch []*big.Int) (*DelegationmanagerUnstakeInitIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var delegatorIdRule []interface{}
	for _, delegatorIdItem := range delegatorId {
		delegatorIdRule = append(delegatorIdRule, delegatorIdItem)
	}
	var deactivationEpochRule []interface{}
	for _, deactivationEpochItem := range deactivationEpoch {
		deactivationEpochRule = append(deactivationEpochRule, deactivationEpochItem)
	}

	logs, sub, err := _Delegationmanager.contract.FilterLogs(opts, "UnstakeInit", userRule, delegatorIdRule, deactivationEpochRule)
	if err != nil {
		return nil, err
	}
	return &DelegationmanagerUnstakeInitIterator{contract: _Delegationmanager.contract, event: "UnstakeInit", logs: logs, sub: sub}, nil
}

// WatchUnstakeInit is a free log subscription operation binding the contract event 0xcde3813e379342e7506ca1f984065c81821e879aebd4be83cb35b5ab976518f9.
//
// Solidity: event UnstakeInit(address indexed user, uint256 indexed delegatorId, uint256 indexed deactivationEpoch)
func (_Delegationmanager *DelegationmanagerFilterer) WatchUnstakeInit(opts *bind.WatchOpts, sink chan<- *DelegationmanagerUnstakeInit, user []common.Address, delegatorId []*big.Int, deactivationEpoch []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var delegatorIdRule []interface{}
	for _, delegatorIdItem := range delegatorId {
		delegatorIdRule = append(delegatorIdRule, delegatorIdItem)
	}
	var deactivationEpochRule []interface{}
	for _, deactivationEpochItem := range deactivationEpoch {
		deactivationEpochRule = append(deactivationEpochRule, deactivationEpochItem)
	}

	logs, sub, err := _Delegationmanager.contract.WatchLogs(opts, "UnstakeInit", userRule, delegatorIdRule, deactivationEpochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationmanagerUnstakeInit)
				if err := _Delegationmanager.contract.UnpackLog(event, "UnstakeInit", log); err != nil {
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

// ParseUnstakeInit is a log parse operation binding the contract event 0xcde3813e379342e7506ca1f984065c81821e879aebd4be83cb35b5ab976518f9.
//
// Solidity: event UnstakeInit(address indexed user, uint256 indexed delegatorId, uint256 indexed deactivationEpoch)
func (_Delegationmanager *DelegationmanagerFilterer) ParseUnstakeInit(log types.Log) (*DelegationmanagerUnstakeInit, error) {
	event := new(DelegationmanagerUnstakeInit)
	if err := _Delegationmanager.contract.UnpackLog(event, "UnstakeInit", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DelegationmanagerUnstakedIterator is returned from FilterUnstaked and is used to iterate over the raw logs and unpacked data for Unstaked events raised by the Delegationmanager contract.
type DelegationmanagerUnstakedIterator struct {
	Event *DelegationmanagerUnstaked // Event containing the contract specifics and raw log

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
func (it *DelegationmanagerUnstakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationmanagerUnstaked)
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
		it.Event = new(DelegationmanagerUnstaked)
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
func (it *DelegationmanagerUnstakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationmanagerUnstakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationmanagerUnstaked represents a Unstaked event raised by the Delegationmanager contract.
type DelegationmanagerUnstaked struct {
	User        common.Address
	DelegatorId *big.Int
	Amount      *big.Int
	Total       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUnstaked is a free log retrieval operation binding the contract event 0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00.
//
// Solidity: event Unstaked(address indexed user, uint256 indexed delegatorId, uint256 amount, uint256 total)
func (_Delegationmanager *DelegationmanagerFilterer) FilterUnstaked(opts *bind.FilterOpts, user []common.Address, delegatorId []*big.Int) (*DelegationmanagerUnstakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var delegatorIdRule []interface{}
	for _, delegatorIdItem := range delegatorId {
		delegatorIdRule = append(delegatorIdRule, delegatorIdItem)
	}

	logs, sub, err := _Delegationmanager.contract.FilterLogs(opts, "Unstaked", userRule, delegatorIdRule)
	if err != nil {
		return nil, err
	}
	return &DelegationmanagerUnstakedIterator{contract: _Delegationmanager.contract, event: "Unstaked", logs: logs, sub: sub}, nil
}

// WatchUnstaked is a free log subscription operation binding the contract event 0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00.
//
// Solidity: event Unstaked(address indexed user, uint256 indexed delegatorId, uint256 amount, uint256 total)
func (_Delegationmanager *DelegationmanagerFilterer) WatchUnstaked(opts *bind.WatchOpts, sink chan<- *DelegationmanagerUnstaked, user []common.Address, delegatorId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var delegatorIdRule []interface{}
	for _, delegatorIdItem := range delegatorId {
		delegatorIdRule = append(delegatorIdRule, delegatorIdItem)
	}

	logs, sub, err := _Delegationmanager.contract.WatchLogs(opts, "Unstaked", userRule, delegatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationmanagerUnstaked)
				if err := _Delegationmanager.contract.UnpackLog(event, "Unstaked", log); err != nil {
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

// ParseUnstaked is a log parse operation binding the contract event 0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00.
//
// Solidity: event Unstaked(address indexed user, uint256 indexed delegatorId, uint256 amount, uint256 total)
func (_Delegationmanager *DelegationmanagerFilterer) ParseUnstaked(log types.Log) (*DelegationmanagerUnstaked, error) {
	event := new(DelegationmanagerUnstaked)
	if err := _Delegationmanager.contract.UnpackLog(event, "Unstaked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// DelegationmanagerUpdateCommissionIterator is returned from FilterUpdateCommission and is used to iterate over the raw logs and unpacked data for UpdateCommission events raised by the Delegationmanager contract.
type DelegationmanagerUpdateCommissionIterator struct {
	Event *DelegationmanagerUpdateCommission // Event containing the contract specifics and raw log

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
func (it *DelegationmanagerUpdateCommissionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DelegationmanagerUpdateCommission)
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
		it.Event = new(DelegationmanagerUpdateCommission)
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
func (it *DelegationmanagerUpdateCommissionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DelegationmanagerUpdateCommissionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DelegationmanagerUpdateCommission represents a UpdateCommission event raised by the Delegationmanager contract.
type DelegationmanagerUpdateCommission struct {
	ValidatorId *big.Int
	Rate        *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUpdateCommission is a free log retrieval operation binding the contract event 0x0c77380d0b73c79a1f7ca53f0f9aae872bc5e0e86132ff08a539c02b8726c636.
//
// Solidity: event UpdateCommission(uint256 indexed validatorId, uint256 indexed rate)
func (_Delegationmanager *DelegationmanagerFilterer) FilterUpdateCommission(opts *bind.FilterOpts, validatorId []*big.Int, rate []*big.Int) (*DelegationmanagerUpdateCommissionIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var rateRule []interface{}
	for _, rateItem := range rate {
		rateRule = append(rateRule, rateItem)
	}

	logs, sub, err := _Delegationmanager.contract.FilterLogs(opts, "UpdateCommission", validatorIdRule, rateRule)
	if err != nil {
		return nil, err
	}
	return &DelegationmanagerUpdateCommissionIterator{contract: _Delegationmanager.contract, event: "UpdateCommission", logs: logs, sub: sub}, nil
}

// WatchUpdateCommission is a free log subscription operation binding the contract event 0x0c77380d0b73c79a1f7ca53f0f9aae872bc5e0e86132ff08a539c02b8726c636.
//
// Solidity: event UpdateCommission(uint256 indexed validatorId, uint256 indexed rate)
func (_Delegationmanager *DelegationmanagerFilterer) WatchUpdateCommission(opts *bind.WatchOpts, sink chan<- *DelegationmanagerUpdateCommission, validatorId []*big.Int, rate []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var rateRule []interface{}
	for _, rateItem := range rate {
		rateRule = append(rateRule, rateItem)
	}

	logs, sub, err := _Delegationmanager.contract.WatchLogs(opts, "UpdateCommission", validatorIdRule, rateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DelegationmanagerUpdateCommission)
				if err := _Delegationmanager.contract.UnpackLog(event, "UpdateCommission", log); err != nil {
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

// ParseUpdateCommission is a log parse operation binding the contract event 0x0c77380d0b73c79a1f7ca53f0f9aae872bc5e0e86132ff08a539c02b8726c636.
//
// Solidity: event UpdateCommission(uint256 indexed validatorId, uint256 indexed rate)
func (_Delegationmanager *DelegationmanagerFilterer) ParseUpdateCommission(log types.Log) (*DelegationmanagerUpdateCommission, error) {
	event := new(DelegationmanagerUpdateCommission)
	if err := _Delegationmanager.contract.UnpackLog(event, "UpdateCommission", log); err != nil {
		return nil, err
	}
	return event, nil
}
