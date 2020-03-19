// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stakinginfo

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

// StakinginfoABI is the input ABI used to generate the binding from.
const StakinginfoABI = "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"ClaimFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"totalAmount\",\"type\":\"uint256\"}],\"name\":\"ClaimRewards\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"newValidatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"oldValidatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"ConfirmAuction\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"rewards\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"DelClaimRewards\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newDynasty\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldDynasty\",\"type\":\"uint256\"}],\"name\":\"DynastyValueChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"exitEpoch\",\"type\":\"uint256\"}],\"name\":\"Jailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"ReStaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newReward\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldReward\",\"type\":\"uint256\"}],\"name\":\"RewardUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"ShareBurned\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"tokens\",\"type\":\"uint256\"}],\"name\":\"ShareMinted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"oldSigner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newSigner\",\"type\":\"address\"}],\"name\":\"SignerChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"StakeUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"activationEpoch\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"auctionAmount\",\"type\":\"uint256\"}],\"name\":\"StartAuction\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"newThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"oldThreshold\",\"type\":\"uint256\"}],\"name\":\"ThresholdChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"fee\",\"type\":\"uint256\"}],\"name\":\"TopUpFee\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"deactivationEpoch\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"UnstakeInit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"Unstaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"newCommissionRate\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"oldCommissionRate\",\"type\":\"uint256\"}],\"name\":\"UpdateCommissionRate\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"getAccountStateRoot\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"accountStateRoot\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"getStakerDetails\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reward\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"activationEpoch\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"deactivationEpoch\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_status\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"getValidatorContractAddress\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"ValidatorContract\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"totalValidatorStake\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"validatorStake\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

// Stakinginfo is an auto generated Go binding around an Ethereum contract.
type Stakinginfo struct {
	StakinginfoCaller     // Read-only binding to the contract
	StakinginfoTransactor // Write-only binding to the contract
	StakinginfoFilterer   // Log filterer for contract events
}

// StakinginfoCaller is an auto generated read-only Go binding around an Ethereum contract.
type StakinginfoCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakinginfoTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StakinginfoTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakinginfoFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StakinginfoFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StakinginfoSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StakinginfoSession struct {
	Contract     *Stakinginfo      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StakinginfoCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StakinginfoCallerSession struct {
	Contract *StakinginfoCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// StakinginfoTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StakinginfoTransactorSession struct {
	Contract     *StakinginfoTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// StakinginfoRaw is an auto generated low-level Go binding around an Ethereum contract.
type StakinginfoRaw struct {
	Contract *Stakinginfo // Generic contract binding to access the raw methods on
}

// StakinginfoCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StakinginfoCallerRaw struct {
	Contract *StakinginfoCaller // Generic read-only contract binding to access the raw methods on
}

// StakinginfoTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StakinginfoTransactorRaw struct {
	Contract *StakinginfoTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStakinginfo creates a new instance of Stakinginfo, bound to a specific deployed contract.
func NewStakinginfo(address common.Address, backend bind.ContractBackend) (*Stakinginfo, error) {
	contract, err := bindStakinginfo(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Stakinginfo{StakinginfoCaller: StakinginfoCaller{contract: contract}, StakinginfoTransactor: StakinginfoTransactor{contract: contract}, StakinginfoFilterer: StakinginfoFilterer{contract: contract}}, nil
}

// NewStakinginfoCaller creates a new read-only instance of Stakinginfo, bound to a specific deployed contract.
func NewStakinginfoCaller(address common.Address, caller bind.ContractCaller) (*StakinginfoCaller, error) {
	contract, err := bindStakinginfo(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StakinginfoCaller{contract: contract}, nil
}

// NewStakinginfoTransactor creates a new write-only instance of Stakinginfo, bound to a specific deployed contract.
func NewStakinginfoTransactor(address common.Address, transactor bind.ContractTransactor) (*StakinginfoTransactor, error) {
	contract, err := bindStakinginfo(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StakinginfoTransactor{contract: contract}, nil
}

// NewStakinginfoFilterer creates a new log filterer instance of Stakinginfo, bound to a specific deployed contract.
func NewStakinginfoFilterer(address common.Address, filterer bind.ContractFilterer) (*StakinginfoFilterer, error) {
	contract, err := bindStakinginfo(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StakinginfoFilterer{contract: contract}, nil
}

// bindStakinginfo binds a generic wrapper to an already deployed contract.
func bindStakinginfo(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StakinginfoABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Stakinginfo *StakinginfoRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Stakinginfo.Contract.StakinginfoCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Stakinginfo *StakinginfoRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakinginfo.Contract.StakinginfoTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Stakinginfo *StakinginfoRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stakinginfo.Contract.StakinginfoTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Stakinginfo *StakinginfoCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Stakinginfo.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Stakinginfo *StakinginfoTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Stakinginfo.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Stakinginfo *StakinginfoTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Stakinginfo.Contract.contract.Transact(opts, method, params...)
}

// GetAccountStateRoot is a free data retrieval call binding the contract method 0x4b6b87ce.
//
// Solidity: function getAccountStateRoot() constant returns(bytes32 accountStateRoot)
func (_Stakinginfo *StakinginfoCaller) GetAccountStateRoot(opts *bind.CallOpts) ([32]byte, error) {
	var (
		ret0 = new([32]byte)
	)
	out := ret0
	err := _Stakinginfo.contract.Call(opts, out, "getAccountStateRoot")
	return *ret0, err
}

// GetAccountStateRoot is a free data retrieval call binding the contract method 0x4b6b87ce.
//
// Solidity: function getAccountStateRoot() constant returns(bytes32 accountStateRoot)
func (_Stakinginfo *StakinginfoSession) GetAccountStateRoot() ([32]byte, error) {
	return _Stakinginfo.Contract.GetAccountStateRoot(&_Stakinginfo.CallOpts)
}

// GetAccountStateRoot is a free data retrieval call binding the contract method 0x4b6b87ce.
//
// Solidity: function getAccountStateRoot() constant returns(bytes32 accountStateRoot)
func (_Stakinginfo *StakinginfoCallerSession) GetAccountStateRoot() ([32]byte, error) {
	return _Stakinginfo.Contract.GetAccountStateRoot(&_Stakinginfo.CallOpts)
}

// GetStakerDetails is a free data retrieval call binding the contract method 0x78daaf69.
//
// Solidity: function getStakerDetails(uint256 validatorId) constant returns(uint256 amount, uint256 reward, uint256 activationEpoch, uint256 deactivationEpoch, address signer, uint256 _status)
func (_Stakinginfo *StakinginfoCaller) GetStakerDetails(opts *bind.CallOpts, validatorId *big.Int) (struct {
	Amount            *big.Int
	Reward            *big.Int
	ActivationEpoch   *big.Int
	DeactivationEpoch *big.Int
	Signer            common.Address
	Status            *big.Int
}, error) {
	ret := new(struct {
		Amount            *big.Int
		Reward            *big.Int
		ActivationEpoch   *big.Int
		DeactivationEpoch *big.Int
		Signer            common.Address
		Status            *big.Int
	})
	out := ret
	err := _Stakinginfo.contract.Call(opts, out, "getStakerDetails", validatorId)
	return *ret, err
}

// GetStakerDetails is a free data retrieval call binding the contract method 0x78daaf69.
//
// Solidity: function getStakerDetails(uint256 validatorId) constant returns(uint256 amount, uint256 reward, uint256 activationEpoch, uint256 deactivationEpoch, address signer, uint256 _status)
func (_Stakinginfo *StakinginfoSession) GetStakerDetails(validatorId *big.Int) (struct {
	Amount            *big.Int
	Reward            *big.Int
	ActivationEpoch   *big.Int
	DeactivationEpoch *big.Int
	Signer            common.Address
	Status            *big.Int
}, error) {
	return _Stakinginfo.Contract.GetStakerDetails(&_Stakinginfo.CallOpts, validatorId)
}

// GetStakerDetails is a free data retrieval call binding the contract method 0x78daaf69.
//
// Solidity: function getStakerDetails(uint256 validatorId) constant returns(uint256 amount, uint256 reward, uint256 activationEpoch, uint256 deactivationEpoch, address signer, uint256 _status)
func (_Stakinginfo *StakinginfoCallerSession) GetStakerDetails(validatorId *big.Int) (struct {
	Amount            *big.Int
	Reward            *big.Int
	ActivationEpoch   *big.Int
	DeactivationEpoch *big.Int
	Signer            common.Address
	Status            *big.Int
}, error) {
	return _Stakinginfo.Contract.GetStakerDetails(&_Stakinginfo.CallOpts, validatorId)
}

// GetValidatorContractAddress is a free data retrieval call binding the contract method 0x178d46aa.
//
// Solidity: function getValidatorContractAddress(uint256 validatorId) constant returns(address ValidatorContract)
func (_Stakinginfo *StakinginfoCaller) GetValidatorContractAddress(opts *bind.CallOpts, validatorId *big.Int) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Stakinginfo.contract.Call(opts, out, "getValidatorContractAddress", validatorId)
	return *ret0, err
}

// GetValidatorContractAddress is a free data retrieval call binding the contract method 0x178d46aa.
//
// Solidity: function getValidatorContractAddress(uint256 validatorId) constant returns(address ValidatorContract)
func (_Stakinginfo *StakinginfoSession) GetValidatorContractAddress(validatorId *big.Int) (common.Address, error) {
	return _Stakinginfo.Contract.GetValidatorContractAddress(&_Stakinginfo.CallOpts, validatorId)
}

// GetValidatorContractAddress is a free data retrieval call binding the contract method 0x178d46aa.
//
// Solidity: function getValidatorContractAddress(uint256 validatorId) constant returns(address ValidatorContract)
func (_Stakinginfo *StakinginfoCallerSession) GetValidatorContractAddress(validatorId *big.Int) (common.Address, error) {
	return _Stakinginfo.Contract.GetValidatorContractAddress(&_Stakinginfo.CallOpts, validatorId)
}

// TotalValidatorStake is a free data retrieval call binding the contract method 0xca7d34b6.
//
// Solidity: function totalValidatorStake(uint256 validatorId) constant returns(uint256 validatorStake)
func (_Stakinginfo *StakinginfoCaller) TotalValidatorStake(opts *bind.CallOpts, validatorId *big.Int) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Stakinginfo.contract.Call(opts, out, "totalValidatorStake", validatorId)
	return *ret0, err
}

// TotalValidatorStake is a free data retrieval call binding the contract method 0xca7d34b6.
//
// Solidity: function totalValidatorStake(uint256 validatorId) constant returns(uint256 validatorStake)
func (_Stakinginfo *StakinginfoSession) TotalValidatorStake(validatorId *big.Int) (*big.Int, error) {
	return _Stakinginfo.Contract.TotalValidatorStake(&_Stakinginfo.CallOpts, validatorId)
}

// TotalValidatorStake is a free data retrieval call binding the contract method 0xca7d34b6.
//
// Solidity: function totalValidatorStake(uint256 validatorId) constant returns(uint256 validatorStake)
func (_Stakinginfo *StakinginfoCallerSession) TotalValidatorStake(validatorId *big.Int) (*big.Int, error) {
	return _Stakinginfo.Contract.TotalValidatorStake(&_Stakinginfo.CallOpts, validatorId)
}

// StakinginfoClaimFeeIterator is returned from FilterClaimFee and is used to iterate over the raw logs and unpacked data for ClaimFee events raised by the Stakinginfo contract.
type StakinginfoClaimFeeIterator struct {
	Event *StakinginfoClaimFee // Event containing the contract specifics and raw log

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
func (it *StakinginfoClaimFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoClaimFee)
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
		it.Event = new(StakinginfoClaimFee)
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
func (it *StakinginfoClaimFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoClaimFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoClaimFee represents a ClaimFee event raised by the Stakinginfo contract.
type StakinginfoClaimFee struct {
	ValidatorId *big.Int
	Signer      common.Address
	Fee         *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterClaimFee is a free log retrieval operation binding the contract event 0x4a83a8cc14c510291ef5a72ae016304a3132c29fad2238eaa29d5b959136f32f.
//
// Solidity: event ClaimFee(uint256 indexed validatorId, address indexed signer, uint256 indexed fee)
func (_Stakinginfo *StakinginfoFilterer) FilterClaimFee(opts *bind.FilterOpts, validatorId []*big.Int, signer []common.Address, fee []*big.Int) (*StakinginfoClaimFeeIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}
	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "ClaimFee", validatorIdRule, signerRule, feeRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoClaimFeeIterator{contract: _Stakinginfo.contract, event: "ClaimFee", logs: logs, sub: sub}, nil
}

// WatchClaimFee is a free log subscription operation binding the contract event 0x4a83a8cc14c510291ef5a72ae016304a3132c29fad2238eaa29d5b959136f32f.
//
// Solidity: event ClaimFee(uint256 indexed validatorId, address indexed signer, uint256 indexed fee)
func (_Stakinginfo *StakinginfoFilterer) WatchClaimFee(opts *bind.WatchOpts, sink chan<- *StakinginfoClaimFee, validatorId []*big.Int, signer []common.Address, fee []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}
	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "ClaimFee", validatorIdRule, signerRule, feeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoClaimFee)
				if err := _Stakinginfo.contract.UnpackLog(event, "ClaimFee", log); err != nil {
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

// ParseClaimFee is a log parse operation binding the contract event 0x4a83a8cc14c510291ef5a72ae016304a3132c29fad2238eaa29d5b959136f32f.
//
// Solidity: event ClaimFee(uint256 indexed validatorId, address indexed signer, uint256 indexed fee)
func (_Stakinginfo *StakinginfoFilterer) ParseClaimFee(log types.Log) (*StakinginfoClaimFee, error) {
	event := new(StakinginfoClaimFee)
	if err := _Stakinginfo.contract.UnpackLog(event, "ClaimFee", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoClaimRewardsIterator is returned from FilterClaimRewards and is used to iterate over the raw logs and unpacked data for ClaimRewards events raised by the Stakinginfo contract.
type StakinginfoClaimRewardsIterator struct {
	Event *StakinginfoClaimRewards // Event containing the contract specifics and raw log

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
func (it *StakinginfoClaimRewardsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoClaimRewards)
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
		it.Event = new(StakinginfoClaimRewards)
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
func (it *StakinginfoClaimRewardsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoClaimRewardsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoClaimRewards represents a ClaimRewards event raised by the Stakinginfo contract.
type StakinginfoClaimRewards struct {
	ValidatorId *big.Int
	Amount      *big.Int
	TotalAmount *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterClaimRewards is a free log retrieval operation binding the contract event 0x41e5e4590cfcde2f03ee9281c54d03acad8adffb83f8310d66b894532470ba35.
//
// Solidity: event ClaimRewards(uint256 indexed validatorId, uint256 indexed amount, uint256 indexed totalAmount)
func (_Stakinginfo *StakinginfoFilterer) FilterClaimRewards(opts *bind.FilterOpts, validatorId []*big.Int, amount []*big.Int, totalAmount []*big.Int) (*StakinginfoClaimRewardsIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var totalAmountRule []interface{}
	for _, totalAmountItem := range totalAmount {
		totalAmountRule = append(totalAmountRule, totalAmountItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "ClaimRewards", validatorIdRule, amountRule, totalAmountRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoClaimRewardsIterator{contract: _Stakinginfo.contract, event: "ClaimRewards", logs: logs, sub: sub}, nil
}

// WatchClaimRewards is a free log subscription operation binding the contract event 0x41e5e4590cfcde2f03ee9281c54d03acad8adffb83f8310d66b894532470ba35.
//
// Solidity: event ClaimRewards(uint256 indexed validatorId, uint256 indexed amount, uint256 indexed totalAmount)
func (_Stakinginfo *StakinginfoFilterer) WatchClaimRewards(opts *bind.WatchOpts, sink chan<- *StakinginfoClaimRewards, validatorId []*big.Int, amount []*big.Int, totalAmount []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var totalAmountRule []interface{}
	for _, totalAmountItem := range totalAmount {
		totalAmountRule = append(totalAmountRule, totalAmountItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "ClaimRewards", validatorIdRule, amountRule, totalAmountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoClaimRewards)
				if err := _Stakinginfo.contract.UnpackLog(event, "ClaimRewards", log); err != nil {
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

// ParseClaimRewards is a log parse operation binding the contract event 0x41e5e4590cfcde2f03ee9281c54d03acad8adffb83f8310d66b894532470ba35.
//
// Solidity: event ClaimRewards(uint256 indexed validatorId, uint256 indexed amount, uint256 indexed totalAmount)
func (_Stakinginfo *StakinginfoFilterer) ParseClaimRewards(log types.Log) (*StakinginfoClaimRewards, error) {
	event := new(StakinginfoClaimRewards)
	if err := _Stakinginfo.contract.UnpackLog(event, "ClaimRewards", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoConfirmAuctionIterator is returned from FilterConfirmAuction and is used to iterate over the raw logs and unpacked data for ConfirmAuction events raised by the Stakinginfo contract.
type StakinginfoConfirmAuctionIterator struct {
	Event *StakinginfoConfirmAuction // Event containing the contract specifics and raw log

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
func (it *StakinginfoConfirmAuctionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoConfirmAuction)
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
		it.Event = new(StakinginfoConfirmAuction)
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
func (it *StakinginfoConfirmAuctionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoConfirmAuctionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoConfirmAuction represents a ConfirmAuction event raised by the Stakinginfo contract.
type StakinginfoConfirmAuction struct {
	NewValidatorId *big.Int
	OldValidatorId *big.Int
	Amount         *big.Int
	Raw            types.Log // Blockchain specific contextual infos
}

// FilterConfirmAuction is a free log retrieval operation binding the contract event 0x1002381ecf76700f6f0ab4c90b9f523e39df7b0482b71ec63cf62cf854120470.
//
// Solidity: event ConfirmAuction(uint256 indexed newValidatorId, uint256 indexed oldValidatorId, uint256 indexed amount)
func (_Stakinginfo *StakinginfoFilterer) FilterConfirmAuction(opts *bind.FilterOpts, newValidatorId []*big.Int, oldValidatorId []*big.Int, amount []*big.Int) (*StakinginfoConfirmAuctionIterator, error) {

	var newValidatorIdRule []interface{}
	for _, newValidatorIdItem := range newValidatorId {
		newValidatorIdRule = append(newValidatorIdRule, newValidatorIdItem)
	}
	var oldValidatorIdRule []interface{}
	for _, oldValidatorIdItem := range oldValidatorId {
		oldValidatorIdRule = append(oldValidatorIdRule, oldValidatorIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "ConfirmAuction", newValidatorIdRule, oldValidatorIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoConfirmAuctionIterator{contract: _Stakinginfo.contract, event: "ConfirmAuction", logs: logs, sub: sub}, nil
}

// WatchConfirmAuction is a free log subscription operation binding the contract event 0x1002381ecf76700f6f0ab4c90b9f523e39df7b0482b71ec63cf62cf854120470.
//
// Solidity: event ConfirmAuction(uint256 indexed newValidatorId, uint256 indexed oldValidatorId, uint256 indexed amount)
func (_Stakinginfo *StakinginfoFilterer) WatchConfirmAuction(opts *bind.WatchOpts, sink chan<- *StakinginfoConfirmAuction, newValidatorId []*big.Int, oldValidatorId []*big.Int, amount []*big.Int) (event.Subscription, error) {

	var newValidatorIdRule []interface{}
	for _, newValidatorIdItem := range newValidatorId {
		newValidatorIdRule = append(newValidatorIdRule, newValidatorIdItem)
	}
	var oldValidatorIdRule []interface{}
	for _, oldValidatorIdItem := range oldValidatorId {
		oldValidatorIdRule = append(oldValidatorIdRule, oldValidatorIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "ConfirmAuction", newValidatorIdRule, oldValidatorIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoConfirmAuction)
				if err := _Stakinginfo.contract.UnpackLog(event, "ConfirmAuction", log); err != nil {
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

// ParseConfirmAuction is a log parse operation binding the contract event 0x1002381ecf76700f6f0ab4c90b9f523e39df7b0482b71ec63cf62cf854120470.
//
// Solidity: event ConfirmAuction(uint256 indexed newValidatorId, uint256 indexed oldValidatorId, uint256 indexed amount)
func (_Stakinginfo *StakinginfoFilterer) ParseConfirmAuction(log types.Log) (*StakinginfoConfirmAuction, error) {
	event := new(StakinginfoConfirmAuction)
	if err := _Stakinginfo.contract.UnpackLog(event, "ConfirmAuction", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoDelClaimRewardsIterator is returned from FilterDelClaimRewards and is used to iterate over the raw logs and unpacked data for DelClaimRewards events raised by the Stakinginfo contract.
type StakinginfoDelClaimRewardsIterator struct {
	Event *StakinginfoDelClaimRewards // Event containing the contract specifics and raw log

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
func (it *StakinginfoDelClaimRewardsIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoDelClaimRewards)
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
		it.Event = new(StakinginfoDelClaimRewards)
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
func (it *StakinginfoDelClaimRewardsIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoDelClaimRewardsIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoDelClaimRewards represents a DelClaimRewards event raised by the Stakinginfo contract.
type StakinginfoDelClaimRewards struct {
	ValidatorId *big.Int
	Rewards     *big.Int
	Tokens      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterDelClaimRewards is a free log retrieval operation binding the contract event 0xad00ba36c13a1357bfe6943f8219d64f36ab74cf469e050a6ba2506101630ce7.
//
// Solidity: event DelClaimRewards(uint256 indexed validatorId, uint256 indexed rewards, uint256 indexed tokens)
func (_Stakinginfo *StakinginfoFilterer) FilterDelClaimRewards(opts *bind.FilterOpts, validatorId []*big.Int, rewards []*big.Int, tokens []*big.Int) (*StakinginfoDelClaimRewardsIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var rewardsRule []interface{}
	for _, rewardsItem := range rewards {
		rewardsRule = append(rewardsRule, rewardsItem)
	}
	var tokensRule []interface{}
	for _, tokensItem := range tokens {
		tokensRule = append(tokensRule, tokensItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "DelClaimRewards", validatorIdRule, rewardsRule, tokensRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoDelClaimRewardsIterator{contract: _Stakinginfo.contract, event: "DelClaimRewards", logs: logs, sub: sub}, nil
}

// WatchDelClaimRewards is a free log subscription operation binding the contract event 0xad00ba36c13a1357bfe6943f8219d64f36ab74cf469e050a6ba2506101630ce7.
//
// Solidity: event DelClaimRewards(uint256 indexed validatorId, uint256 indexed rewards, uint256 indexed tokens)
func (_Stakinginfo *StakinginfoFilterer) WatchDelClaimRewards(opts *bind.WatchOpts, sink chan<- *StakinginfoDelClaimRewards, validatorId []*big.Int, rewards []*big.Int, tokens []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var rewardsRule []interface{}
	for _, rewardsItem := range rewards {
		rewardsRule = append(rewardsRule, rewardsItem)
	}
	var tokensRule []interface{}
	for _, tokensItem := range tokens {
		tokensRule = append(tokensRule, tokensItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "DelClaimRewards", validatorIdRule, rewardsRule, tokensRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoDelClaimRewards)
				if err := _Stakinginfo.contract.UnpackLog(event, "DelClaimRewards", log); err != nil {
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

// ParseDelClaimRewards is a log parse operation binding the contract event 0xad00ba36c13a1357bfe6943f8219d64f36ab74cf469e050a6ba2506101630ce7.
//
// Solidity: event DelClaimRewards(uint256 indexed validatorId, uint256 indexed rewards, uint256 indexed tokens)
func (_Stakinginfo *StakinginfoFilterer) ParseDelClaimRewards(log types.Log) (*StakinginfoDelClaimRewards, error) {
	event := new(StakinginfoDelClaimRewards)
	if err := _Stakinginfo.contract.UnpackLog(event, "DelClaimRewards", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoDynastyValueChangeIterator is returned from FilterDynastyValueChange and is used to iterate over the raw logs and unpacked data for DynastyValueChange events raised by the Stakinginfo contract.
type StakinginfoDynastyValueChangeIterator struct {
	Event *StakinginfoDynastyValueChange // Event containing the contract specifics and raw log

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
func (it *StakinginfoDynastyValueChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoDynastyValueChange)
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
		it.Event = new(StakinginfoDynastyValueChange)
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
func (it *StakinginfoDynastyValueChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoDynastyValueChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoDynastyValueChange represents a DynastyValueChange event raised by the Stakinginfo contract.
type StakinginfoDynastyValueChange struct {
	NewDynasty *big.Int
	OldDynasty *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterDynastyValueChange is a free log retrieval operation binding the contract event 0x9444bfcfa6aed72a15da73de1220dcc07d7864119c44abfec0037bbcacefda98.
//
// Solidity: event DynastyValueChange(uint256 newDynasty, uint256 oldDynasty)
func (_Stakinginfo *StakinginfoFilterer) FilterDynastyValueChange(opts *bind.FilterOpts) (*StakinginfoDynastyValueChangeIterator, error) {

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "DynastyValueChange")
	if err != nil {
		return nil, err
	}
	return &StakinginfoDynastyValueChangeIterator{contract: _Stakinginfo.contract, event: "DynastyValueChange", logs: logs, sub: sub}, nil
}

// WatchDynastyValueChange is a free log subscription operation binding the contract event 0x9444bfcfa6aed72a15da73de1220dcc07d7864119c44abfec0037bbcacefda98.
//
// Solidity: event DynastyValueChange(uint256 newDynasty, uint256 oldDynasty)
func (_Stakinginfo *StakinginfoFilterer) WatchDynastyValueChange(opts *bind.WatchOpts, sink chan<- *StakinginfoDynastyValueChange) (event.Subscription, error) {

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "DynastyValueChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoDynastyValueChange)
				if err := _Stakinginfo.contract.UnpackLog(event, "DynastyValueChange", log); err != nil {
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

// ParseDynastyValueChange is a log parse operation binding the contract event 0x9444bfcfa6aed72a15da73de1220dcc07d7864119c44abfec0037bbcacefda98.
//
// Solidity: event DynastyValueChange(uint256 newDynasty, uint256 oldDynasty)
func (_Stakinginfo *StakinginfoFilterer) ParseDynastyValueChange(log types.Log) (*StakinginfoDynastyValueChange, error) {
	event := new(StakinginfoDynastyValueChange)
	if err := _Stakinginfo.contract.UnpackLog(event, "DynastyValueChange", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoJailedIterator is returned from FilterJailed and is used to iterate over the raw logs and unpacked data for Jailed events raised by the Stakinginfo contract.
type StakinginfoJailedIterator struct {
	Event *StakinginfoJailed // Event containing the contract specifics and raw log

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
func (it *StakinginfoJailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoJailed)
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
		it.Event = new(StakinginfoJailed)
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
func (it *StakinginfoJailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoJailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoJailed represents a Jailed event raised by the Stakinginfo contract.
type StakinginfoJailed struct {
	ValidatorId *big.Int
	ExitEpoch   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterJailed is a free log retrieval operation binding the contract event 0xa1735a3843d9467dd849a217582720a8af66b9d034712e0b21b0f5ece80670cd.
//
// Solidity: event Jailed(uint256 indexed validatorId, uint256 indexed exitEpoch)
func (_Stakinginfo *StakinginfoFilterer) FilterJailed(opts *bind.FilterOpts, validatorId []*big.Int, exitEpoch []*big.Int) (*StakinginfoJailedIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var exitEpochRule []interface{}
	for _, exitEpochItem := range exitEpoch {
		exitEpochRule = append(exitEpochRule, exitEpochItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "Jailed", validatorIdRule, exitEpochRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoJailedIterator{contract: _Stakinginfo.contract, event: "Jailed", logs: logs, sub: sub}, nil
}

// WatchJailed is a free log subscription operation binding the contract event 0xa1735a3843d9467dd849a217582720a8af66b9d034712e0b21b0f5ece80670cd.
//
// Solidity: event Jailed(uint256 indexed validatorId, uint256 indexed exitEpoch)
func (_Stakinginfo *StakinginfoFilterer) WatchJailed(opts *bind.WatchOpts, sink chan<- *StakinginfoJailed, validatorId []*big.Int, exitEpoch []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var exitEpochRule []interface{}
	for _, exitEpochItem := range exitEpoch {
		exitEpochRule = append(exitEpochRule, exitEpochItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "Jailed", validatorIdRule, exitEpochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoJailed)
				if err := _Stakinginfo.contract.UnpackLog(event, "Jailed", log); err != nil {
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

// ParseJailed is a log parse operation binding the contract event 0xa1735a3843d9467dd849a217582720a8af66b9d034712e0b21b0f5ece80670cd.
//
// Solidity: event Jailed(uint256 indexed validatorId, uint256 indexed exitEpoch)
func (_Stakinginfo *StakinginfoFilterer) ParseJailed(log types.Log) (*StakinginfoJailed, error) {
	event := new(StakinginfoJailed)
	if err := _Stakinginfo.contract.UnpackLog(event, "Jailed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoReStakedIterator is returned from FilterReStaked and is used to iterate over the raw logs and unpacked data for ReStaked events raised by the Stakinginfo contract.
type StakinginfoReStakedIterator struct {
	Event *StakinginfoReStaked // Event containing the contract specifics and raw log

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
func (it *StakinginfoReStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoReStaked)
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
		it.Event = new(StakinginfoReStaked)
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
func (it *StakinginfoReStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoReStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoReStaked represents a ReStaked event raised by the Stakinginfo contract.
type StakinginfoReStaked struct {
	ValidatorId *big.Int
	Amount      *big.Int
	Total       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterReStaked is a free log retrieval operation binding the contract event 0x9cc0e589f20d3310eb2ad571b23529003bd46048d0d1af29277dcf0aa3c398ce.
//
// Solidity: event ReStaked(uint256 indexed validatorId, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) FilterReStaked(opts *bind.FilterOpts, validatorId []*big.Int) (*StakinginfoReStakedIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "ReStaked", validatorIdRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoReStakedIterator{contract: _Stakinginfo.contract, event: "ReStaked", logs: logs, sub: sub}, nil
}

// WatchReStaked is a free log subscription operation binding the contract event 0x9cc0e589f20d3310eb2ad571b23529003bd46048d0d1af29277dcf0aa3c398ce.
//
// Solidity: event ReStaked(uint256 indexed validatorId, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) WatchReStaked(opts *bind.WatchOpts, sink chan<- *StakinginfoReStaked, validatorId []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "ReStaked", validatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoReStaked)
				if err := _Stakinginfo.contract.UnpackLog(event, "ReStaked", log); err != nil {
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
// Solidity: event ReStaked(uint256 indexed validatorId, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) ParseReStaked(log types.Log) (*StakinginfoReStaked, error) {
	event := new(StakinginfoReStaked)
	if err := _Stakinginfo.contract.UnpackLog(event, "ReStaked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoRewardUpdateIterator is returned from FilterRewardUpdate and is used to iterate over the raw logs and unpacked data for RewardUpdate events raised by the Stakinginfo contract.
type StakinginfoRewardUpdateIterator struct {
	Event *StakinginfoRewardUpdate // Event containing the contract specifics and raw log

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
func (it *StakinginfoRewardUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoRewardUpdate)
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
		it.Event = new(StakinginfoRewardUpdate)
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
func (it *StakinginfoRewardUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoRewardUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoRewardUpdate represents a RewardUpdate event raised by the Stakinginfo contract.
type StakinginfoRewardUpdate struct {
	NewReward *big.Int
	OldReward *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterRewardUpdate is a free log retrieval operation binding the contract event 0xf67f33e8589d3ea0356303c0f9a8e764873692159f777ff79e4fc523d389dfcd.
//
// Solidity: event RewardUpdate(uint256 newReward, uint256 oldReward)
func (_Stakinginfo *StakinginfoFilterer) FilterRewardUpdate(opts *bind.FilterOpts) (*StakinginfoRewardUpdateIterator, error) {

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "RewardUpdate")
	if err != nil {
		return nil, err
	}
	return &StakinginfoRewardUpdateIterator{contract: _Stakinginfo.contract, event: "RewardUpdate", logs: logs, sub: sub}, nil
}

// WatchRewardUpdate is a free log subscription operation binding the contract event 0xf67f33e8589d3ea0356303c0f9a8e764873692159f777ff79e4fc523d389dfcd.
//
// Solidity: event RewardUpdate(uint256 newReward, uint256 oldReward)
func (_Stakinginfo *StakinginfoFilterer) WatchRewardUpdate(opts *bind.WatchOpts, sink chan<- *StakinginfoRewardUpdate) (event.Subscription, error) {

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "RewardUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoRewardUpdate)
				if err := _Stakinginfo.contract.UnpackLog(event, "RewardUpdate", log); err != nil {
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

// ParseRewardUpdate is a log parse operation binding the contract event 0xf67f33e8589d3ea0356303c0f9a8e764873692159f777ff79e4fc523d389dfcd.
//
// Solidity: event RewardUpdate(uint256 newReward, uint256 oldReward)
func (_Stakinginfo *StakinginfoFilterer) ParseRewardUpdate(log types.Log) (*StakinginfoRewardUpdate, error) {
	event := new(StakinginfoRewardUpdate)
	if err := _Stakinginfo.contract.UnpackLog(event, "RewardUpdate", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoShareBurnedIterator is returned from FilterShareBurned and is used to iterate over the raw logs and unpacked data for ShareBurned events raised by the Stakinginfo contract.
type StakinginfoShareBurnedIterator struct {
	Event *StakinginfoShareBurned // Event containing the contract specifics and raw log

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
func (it *StakinginfoShareBurnedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoShareBurned)
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
		it.Event = new(StakinginfoShareBurned)
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
func (it *StakinginfoShareBurnedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoShareBurnedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoShareBurned represents a ShareBurned event raised by the Stakinginfo contract.
type StakinginfoShareBurned struct {
	ValidatorId *big.Int
	User        common.Address
	Amount      *big.Int
	Tokens      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterShareBurned is a free log retrieval operation binding the contract event 0x7e86625aa6e668407f095af342e0cc237809c4c5086b4d665a0067de122980a9.
//
// Solidity: event ShareBurned(uint256 indexed validatorId, address indexed user, uint256 indexed amount, uint256 tokens)
func (_Stakinginfo *StakinginfoFilterer) FilterShareBurned(opts *bind.FilterOpts, validatorId []*big.Int, user []common.Address, amount []*big.Int) (*StakinginfoShareBurnedIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "ShareBurned", validatorIdRule, userRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoShareBurnedIterator{contract: _Stakinginfo.contract, event: "ShareBurned", logs: logs, sub: sub}, nil
}

// WatchShareBurned is a free log subscription operation binding the contract event 0x7e86625aa6e668407f095af342e0cc237809c4c5086b4d665a0067de122980a9.
//
// Solidity: event ShareBurned(uint256 indexed validatorId, address indexed user, uint256 indexed amount, uint256 tokens)
func (_Stakinginfo *StakinginfoFilterer) WatchShareBurned(opts *bind.WatchOpts, sink chan<- *StakinginfoShareBurned, validatorId []*big.Int, user []common.Address, amount []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "ShareBurned", validatorIdRule, userRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoShareBurned)
				if err := _Stakinginfo.contract.UnpackLog(event, "ShareBurned", log); err != nil {
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

// ParseShareBurned is a log parse operation binding the contract event 0x7e86625aa6e668407f095af342e0cc237809c4c5086b4d665a0067de122980a9.
//
// Solidity: event ShareBurned(uint256 indexed validatorId, address indexed user, uint256 indexed amount, uint256 tokens)
func (_Stakinginfo *StakinginfoFilterer) ParseShareBurned(log types.Log) (*StakinginfoShareBurned, error) {
	event := new(StakinginfoShareBurned)
	if err := _Stakinginfo.contract.UnpackLog(event, "ShareBurned", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoShareMintedIterator is returned from FilterShareMinted and is used to iterate over the raw logs and unpacked data for ShareMinted events raised by the Stakinginfo contract.
type StakinginfoShareMintedIterator struct {
	Event *StakinginfoShareMinted // Event containing the contract specifics and raw log

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
func (it *StakinginfoShareMintedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoShareMinted)
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
		it.Event = new(StakinginfoShareMinted)
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
func (it *StakinginfoShareMintedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoShareMintedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoShareMinted represents a ShareMinted event raised by the Stakinginfo contract.
type StakinginfoShareMinted struct {
	ValidatorId *big.Int
	User        common.Address
	Amount      *big.Int
	Tokens      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterShareMinted is a free log retrieval operation binding the contract event 0xc9afff0972d33d68c8d330fe0ebd0e9f54491ad8c59ae17330a9206f280f0865.
//
// Solidity: event ShareMinted(uint256 indexed validatorId, address indexed user, uint256 indexed amount, uint256 tokens)
func (_Stakinginfo *StakinginfoFilterer) FilterShareMinted(opts *bind.FilterOpts, validatorId []*big.Int, user []common.Address, amount []*big.Int) (*StakinginfoShareMintedIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "ShareMinted", validatorIdRule, userRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoShareMintedIterator{contract: _Stakinginfo.contract, event: "ShareMinted", logs: logs, sub: sub}, nil
}

// WatchShareMinted is a free log subscription operation binding the contract event 0xc9afff0972d33d68c8d330fe0ebd0e9f54491ad8c59ae17330a9206f280f0865.
//
// Solidity: event ShareMinted(uint256 indexed validatorId, address indexed user, uint256 indexed amount, uint256 tokens)
func (_Stakinginfo *StakinginfoFilterer) WatchShareMinted(opts *bind.WatchOpts, sink chan<- *StakinginfoShareMinted, validatorId []*big.Int, user []common.Address, amount []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "ShareMinted", validatorIdRule, userRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoShareMinted)
				if err := _Stakinginfo.contract.UnpackLog(event, "ShareMinted", log); err != nil {
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

// ParseShareMinted is a log parse operation binding the contract event 0xc9afff0972d33d68c8d330fe0ebd0e9f54491ad8c59ae17330a9206f280f0865.
//
// Solidity: event ShareMinted(uint256 indexed validatorId, address indexed user, uint256 indexed amount, uint256 tokens)
func (_Stakinginfo *StakinginfoFilterer) ParseShareMinted(log types.Log) (*StakinginfoShareMinted, error) {
	event := new(StakinginfoShareMinted)
	if err := _Stakinginfo.contract.UnpackLog(event, "ShareMinted", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoSignerChangeIterator is returned from FilterSignerChange and is used to iterate over the raw logs and unpacked data for SignerChange events raised by the Stakinginfo contract.
type StakinginfoSignerChangeIterator struct {
	Event *StakinginfoSignerChange // Event containing the contract specifics and raw log

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
func (it *StakinginfoSignerChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoSignerChange)
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
		it.Event = new(StakinginfoSignerChange)
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
func (it *StakinginfoSignerChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoSignerChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoSignerChange represents a SignerChange event raised by the Stakinginfo contract.
type StakinginfoSignerChange struct {
	ValidatorId *big.Int
	OldSigner   common.Address
	NewSigner   common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSignerChange is a free log retrieval operation binding the contract event 0x7dfd3bad1e3cac97d3b89ff06d78394523c4f08fdee4daa71a59160003240c89.
//
// Solidity: event SignerChange(uint256 indexed validatorId, address indexed oldSigner, address indexed newSigner)
func (_Stakinginfo *StakinginfoFilterer) FilterSignerChange(opts *bind.FilterOpts, validatorId []*big.Int, oldSigner []common.Address, newSigner []common.Address) (*StakinginfoSignerChangeIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var oldSignerRule []interface{}
	for _, oldSignerItem := range oldSigner {
		oldSignerRule = append(oldSignerRule, oldSignerItem)
	}
	var newSignerRule []interface{}
	for _, newSignerItem := range newSigner {
		newSignerRule = append(newSignerRule, newSignerItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "SignerChange", validatorIdRule, oldSignerRule, newSignerRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoSignerChangeIterator{contract: _Stakinginfo.contract, event: "SignerChange", logs: logs, sub: sub}, nil
}

// WatchSignerChange is a free log subscription operation binding the contract event 0x7dfd3bad1e3cac97d3b89ff06d78394523c4f08fdee4daa71a59160003240c89.
//
// Solidity: event SignerChange(uint256 indexed validatorId, address indexed oldSigner, address indexed newSigner)
func (_Stakinginfo *StakinginfoFilterer) WatchSignerChange(opts *bind.WatchOpts, sink chan<- *StakinginfoSignerChange, validatorId []*big.Int, oldSigner []common.Address, newSigner []common.Address) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var oldSignerRule []interface{}
	for _, oldSignerItem := range oldSigner {
		oldSignerRule = append(oldSignerRule, oldSignerItem)
	}
	var newSignerRule []interface{}
	for _, newSignerItem := range newSigner {
		newSignerRule = append(newSignerRule, newSignerItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "SignerChange", validatorIdRule, oldSignerRule, newSignerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoSignerChange)
				if err := _Stakinginfo.contract.UnpackLog(event, "SignerChange", log); err != nil {
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

// ParseSignerChange is a log parse operation binding the contract event 0x7dfd3bad1e3cac97d3b89ff06d78394523c4f08fdee4daa71a59160003240c89.
//
// Solidity: event SignerChange(uint256 indexed validatorId, address indexed oldSigner, address indexed newSigner)
func (_Stakinginfo *StakinginfoFilterer) ParseSignerChange(log types.Log) (*StakinginfoSignerChange, error) {
	event := new(StakinginfoSignerChange)
	if err := _Stakinginfo.contract.UnpackLog(event, "SignerChange", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoStakeUpdateIterator is returned from FilterStakeUpdate and is used to iterate over the raw logs and unpacked data for StakeUpdate events raised by the Stakinginfo contract.
type StakinginfoStakeUpdateIterator struct {
	Event *StakinginfoStakeUpdate // Event containing the contract specifics and raw log

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
func (it *StakinginfoStakeUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoStakeUpdate)
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
		it.Event = new(StakinginfoStakeUpdate)
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
func (it *StakinginfoStakeUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoStakeUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoStakeUpdate represents a StakeUpdate event raised by the Stakinginfo contract.
type StakinginfoStakeUpdate struct {
	ValidatorId *big.Int
	NewAmount   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterStakeUpdate is a free log retrieval operation binding the contract event 0x8a9fd81e4de0eb5b1f4055c5b30b562c02c95613ea0c5259c23033a90d1de4b4.
//
// Solidity: event StakeUpdate(uint256 indexed validatorId, uint256 indexed newAmount)
func (_Stakinginfo *StakinginfoFilterer) FilterStakeUpdate(opts *bind.FilterOpts, validatorId []*big.Int, newAmount []*big.Int) (*StakinginfoStakeUpdateIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var newAmountRule []interface{}
	for _, newAmountItem := range newAmount {
		newAmountRule = append(newAmountRule, newAmountItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "StakeUpdate", validatorIdRule, newAmountRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoStakeUpdateIterator{contract: _Stakinginfo.contract, event: "StakeUpdate", logs: logs, sub: sub}, nil
}

// WatchStakeUpdate is a free log subscription operation binding the contract event 0x8a9fd81e4de0eb5b1f4055c5b30b562c02c95613ea0c5259c23033a90d1de4b4.
//
// Solidity: event StakeUpdate(uint256 indexed validatorId, uint256 indexed newAmount)
func (_Stakinginfo *StakinginfoFilterer) WatchStakeUpdate(opts *bind.WatchOpts, sink chan<- *StakinginfoStakeUpdate, validatorId []*big.Int, newAmount []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var newAmountRule []interface{}
	for _, newAmountItem := range newAmount {
		newAmountRule = append(newAmountRule, newAmountItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "StakeUpdate", validatorIdRule, newAmountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoStakeUpdate)
				if err := _Stakinginfo.contract.UnpackLog(event, "StakeUpdate", log); err != nil {
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

// ParseStakeUpdate is a log parse operation binding the contract event 0x8a9fd81e4de0eb5b1f4055c5b30b562c02c95613ea0c5259c23033a90d1de4b4.
//
// Solidity: event StakeUpdate(uint256 indexed validatorId, uint256 indexed newAmount)
func (_Stakinginfo *StakinginfoFilterer) ParseStakeUpdate(log types.Log) (*StakinginfoStakeUpdate, error) {
	event := new(StakinginfoStakeUpdate)
	if err := _Stakinginfo.contract.UnpackLog(event, "StakeUpdate", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the Stakinginfo contract.
type StakinginfoStakedIterator struct {
	Event *StakinginfoStaked // Event containing the contract specifics and raw log

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
func (it *StakinginfoStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoStaked)
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
		it.Event = new(StakinginfoStaked)
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
func (it *StakinginfoStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoStaked represents a Staked event raised by the Stakinginfo contract.
type StakinginfoStaked struct {
	Signer          common.Address
	ValidatorId     *big.Int
	ActivationEpoch *big.Int
	Amount          *big.Int
	Total           *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0x9cfd25589d1eb8ad71e342a86a8524e83522e3936c0803048c08f6d9ad974f40.
//
// Solidity: event Staked(address indexed signer, uint256 indexed validatorId, uint256 indexed activationEpoch, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) FilterStaked(opts *bind.FilterOpts, signer []common.Address, validatorId []*big.Int, activationEpoch []*big.Int) (*StakinginfoStakedIterator, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var activationEpochRule []interface{}
	for _, activationEpochItem := range activationEpoch {
		activationEpochRule = append(activationEpochRule, activationEpochItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "Staked", signerRule, validatorIdRule, activationEpochRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoStakedIterator{contract: _Stakinginfo.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0x9cfd25589d1eb8ad71e342a86a8524e83522e3936c0803048c08f6d9ad974f40.
//
// Solidity: event Staked(address indexed signer, uint256 indexed validatorId, uint256 indexed activationEpoch, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *StakinginfoStaked, signer []common.Address, validatorId []*big.Int, activationEpoch []*big.Int) (event.Subscription, error) {

	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var activationEpochRule []interface{}
	for _, activationEpochItem := range activationEpoch {
		activationEpochRule = append(activationEpochRule, activationEpochItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "Staked", signerRule, validatorIdRule, activationEpochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoStaked)
				if err := _Stakinginfo.contract.UnpackLog(event, "Staked", log); err != nil {
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
// Solidity: event Staked(address indexed signer, uint256 indexed validatorId, uint256 indexed activationEpoch, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) ParseStaked(log types.Log) (*StakinginfoStaked, error) {
	event := new(StakinginfoStaked)
	if err := _Stakinginfo.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoStartAuctionIterator is returned from FilterStartAuction and is used to iterate over the raw logs and unpacked data for StartAuction events raised by the Stakinginfo contract.
type StakinginfoStartAuctionIterator struct {
	Event *StakinginfoStartAuction // Event containing the contract specifics and raw log

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
func (it *StakinginfoStartAuctionIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoStartAuction)
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
		it.Event = new(StakinginfoStartAuction)
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
func (it *StakinginfoStartAuctionIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoStartAuctionIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoStartAuction represents a StartAuction event raised by the Stakinginfo contract.
type StakinginfoStartAuction struct {
	ValidatorId   *big.Int
	Amount        *big.Int
	AuctionAmount *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterStartAuction is a free log retrieval operation binding the contract event 0x683d0f47c7fa11331f4e9563b3f5a7fdc3d3c5b75c600357a91d991f5a13a437.
//
// Solidity: event StartAuction(uint256 indexed validatorId, uint256 indexed amount, uint256 indexed auctionAmount)
func (_Stakinginfo *StakinginfoFilterer) FilterStartAuction(opts *bind.FilterOpts, validatorId []*big.Int, amount []*big.Int, auctionAmount []*big.Int) (*StakinginfoStartAuctionIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var auctionAmountRule []interface{}
	for _, auctionAmountItem := range auctionAmount {
		auctionAmountRule = append(auctionAmountRule, auctionAmountItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "StartAuction", validatorIdRule, amountRule, auctionAmountRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoStartAuctionIterator{contract: _Stakinginfo.contract, event: "StartAuction", logs: logs, sub: sub}, nil
}

// WatchStartAuction is a free log subscription operation binding the contract event 0x683d0f47c7fa11331f4e9563b3f5a7fdc3d3c5b75c600357a91d991f5a13a437.
//
// Solidity: event StartAuction(uint256 indexed validatorId, uint256 indexed amount, uint256 indexed auctionAmount)
func (_Stakinginfo *StakinginfoFilterer) WatchStartAuction(opts *bind.WatchOpts, sink chan<- *StakinginfoStartAuction, validatorId []*big.Int, amount []*big.Int, auctionAmount []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}
	var auctionAmountRule []interface{}
	for _, auctionAmountItem := range auctionAmount {
		auctionAmountRule = append(auctionAmountRule, auctionAmountItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "StartAuction", validatorIdRule, amountRule, auctionAmountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoStartAuction)
				if err := _Stakinginfo.contract.UnpackLog(event, "StartAuction", log); err != nil {
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

// ParseStartAuction is a log parse operation binding the contract event 0x683d0f47c7fa11331f4e9563b3f5a7fdc3d3c5b75c600357a91d991f5a13a437.
//
// Solidity: event StartAuction(uint256 indexed validatorId, uint256 indexed amount, uint256 indexed auctionAmount)
func (_Stakinginfo *StakinginfoFilterer) ParseStartAuction(log types.Log) (*StakinginfoStartAuction, error) {
	event := new(StakinginfoStartAuction)
	if err := _Stakinginfo.contract.UnpackLog(event, "StartAuction", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoThresholdChangeIterator is returned from FilterThresholdChange and is used to iterate over the raw logs and unpacked data for ThresholdChange events raised by the Stakinginfo contract.
type StakinginfoThresholdChangeIterator struct {
	Event *StakinginfoThresholdChange // Event containing the contract specifics and raw log

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
func (it *StakinginfoThresholdChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoThresholdChange)
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
		it.Event = new(StakinginfoThresholdChange)
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
func (it *StakinginfoThresholdChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoThresholdChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoThresholdChange represents a ThresholdChange event raised by the Stakinginfo contract.
type StakinginfoThresholdChange struct {
	NewThreshold *big.Int
	OldThreshold *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterThresholdChange is a free log retrieval operation binding the contract event 0x5d16a900896e1160c2033bc940e6b072d3dc3b6a996fefb9b3b9b9678841824c.
//
// Solidity: event ThresholdChange(uint256 newThreshold, uint256 oldThreshold)
func (_Stakinginfo *StakinginfoFilterer) FilterThresholdChange(opts *bind.FilterOpts) (*StakinginfoThresholdChangeIterator, error) {

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "ThresholdChange")
	if err != nil {
		return nil, err
	}
	return &StakinginfoThresholdChangeIterator{contract: _Stakinginfo.contract, event: "ThresholdChange", logs: logs, sub: sub}, nil
}

// WatchThresholdChange is a free log subscription operation binding the contract event 0x5d16a900896e1160c2033bc940e6b072d3dc3b6a996fefb9b3b9b9678841824c.
//
// Solidity: event ThresholdChange(uint256 newThreshold, uint256 oldThreshold)
func (_Stakinginfo *StakinginfoFilterer) WatchThresholdChange(opts *bind.WatchOpts, sink chan<- *StakinginfoThresholdChange) (event.Subscription, error) {

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "ThresholdChange")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoThresholdChange)
				if err := _Stakinginfo.contract.UnpackLog(event, "ThresholdChange", log); err != nil {
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

// ParseThresholdChange is a log parse operation binding the contract event 0x5d16a900896e1160c2033bc940e6b072d3dc3b6a996fefb9b3b9b9678841824c.
//
// Solidity: event ThresholdChange(uint256 newThreshold, uint256 oldThreshold)
func (_Stakinginfo *StakinginfoFilterer) ParseThresholdChange(log types.Log) (*StakinginfoThresholdChange, error) {
	event := new(StakinginfoThresholdChange)
	if err := _Stakinginfo.contract.UnpackLog(event, "ThresholdChange", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoTopUpFeeIterator is returned from FilterTopUpFee and is used to iterate over the raw logs and unpacked data for TopUpFee events raised by the Stakinginfo contract.
type StakinginfoTopUpFeeIterator struct {
	Event *StakinginfoTopUpFee // Event containing the contract specifics and raw log

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
func (it *StakinginfoTopUpFeeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoTopUpFee)
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
		it.Event = new(StakinginfoTopUpFee)
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
func (it *StakinginfoTopUpFeeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoTopUpFeeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoTopUpFee represents a TopUpFee event raised by the Stakinginfo contract.
type StakinginfoTopUpFee struct {
	ValidatorId *big.Int
	Signer      common.Address
	Fee         *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTopUpFee is a free log retrieval operation binding the contract event 0xc970ee1564cdc7a9beeae09b65b33aeb17b67d0e4a620e62e6b9b82a6430927d.
//
// Solidity: event TopUpFee(uint256 indexed validatorId, address indexed signer, uint256 indexed fee)
func (_Stakinginfo *StakinginfoFilterer) FilterTopUpFee(opts *bind.FilterOpts, validatorId []*big.Int, signer []common.Address, fee []*big.Int) (*StakinginfoTopUpFeeIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}
	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "TopUpFee", validatorIdRule, signerRule, feeRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoTopUpFeeIterator{contract: _Stakinginfo.contract, event: "TopUpFee", logs: logs, sub: sub}, nil
}

// WatchTopUpFee is a free log subscription operation binding the contract event 0xc970ee1564cdc7a9beeae09b65b33aeb17b67d0e4a620e62e6b9b82a6430927d.
//
// Solidity: event TopUpFee(uint256 indexed validatorId, address indexed signer, uint256 indexed fee)
func (_Stakinginfo *StakinginfoFilterer) WatchTopUpFee(opts *bind.WatchOpts, sink chan<- *StakinginfoTopUpFee, validatorId []*big.Int, signer []common.Address, fee []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var signerRule []interface{}
	for _, signerItem := range signer {
		signerRule = append(signerRule, signerItem)
	}
	var feeRule []interface{}
	for _, feeItem := range fee {
		feeRule = append(feeRule, feeItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "TopUpFee", validatorIdRule, signerRule, feeRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoTopUpFee)
				if err := _Stakinginfo.contract.UnpackLog(event, "TopUpFee", log); err != nil {
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

// ParseTopUpFee is a log parse operation binding the contract event 0xc970ee1564cdc7a9beeae09b65b33aeb17b67d0e4a620e62e6b9b82a6430927d.
//
// Solidity: event TopUpFee(uint256 indexed validatorId, address indexed signer, uint256 indexed fee)
func (_Stakinginfo *StakinginfoFilterer) ParseTopUpFee(log types.Log) (*StakinginfoTopUpFee, error) {
	event := new(StakinginfoTopUpFee)
	if err := _Stakinginfo.contract.UnpackLog(event, "TopUpFee", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoUnstakeInitIterator is returned from FilterUnstakeInit and is used to iterate over the raw logs and unpacked data for UnstakeInit events raised by the Stakinginfo contract.
type StakinginfoUnstakeInitIterator struct {
	Event *StakinginfoUnstakeInit // Event containing the contract specifics and raw log

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
func (it *StakinginfoUnstakeInitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoUnstakeInit)
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
		it.Event = new(StakinginfoUnstakeInit)
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
func (it *StakinginfoUnstakeInitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoUnstakeInitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoUnstakeInit represents a UnstakeInit event raised by the Stakinginfo contract.
type StakinginfoUnstakeInit struct {
	User              common.Address
	ValidatorId       *big.Int
	DeactivationEpoch *big.Int
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUnstakeInit is a free log retrieval operation binding the contract event 0x51ff6d8ad000896362d1a3c9de9ad835adb6da92ec3ddee44018ef64f3c8b561.
//
// Solidity: event UnstakeInit(address indexed user, uint256 indexed validatorId, uint256 deactivationEpoch, uint256 indexed amount)
func (_Stakinginfo *StakinginfoFilterer) FilterUnstakeInit(opts *bind.FilterOpts, user []common.Address, validatorId []*big.Int, amount []*big.Int) (*StakinginfoUnstakeInitIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "UnstakeInit", userRule, validatorIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoUnstakeInitIterator{contract: _Stakinginfo.contract, event: "UnstakeInit", logs: logs, sub: sub}, nil
}

// WatchUnstakeInit is a free log subscription operation binding the contract event 0x51ff6d8ad000896362d1a3c9de9ad835adb6da92ec3ddee44018ef64f3c8b561.
//
// Solidity: event UnstakeInit(address indexed user, uint256 indexed validatorId, uint256 deactivationEpoch, uint256 indexed amount)
func (_Stakinginfo *StakinginfoFilterer) WatchUnstakeInit(opts *bind.WatchOpts, sink chan<- *StakinginfoUnstakeInit, user []common.Address, validatorId []*big.Int, amount []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "UnstakeInit", userRule, validatorIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoUnstakeInit)
				if err := _Stakinginfo.contract.UnpackLog(event, "UnstakeInit", log); err != nil {
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

// ParseUnstakeInit is a log parse operation binding the contract event 0x51ff6d8ad000896362d1a3c9de9ad835adb6da92ec3ddee44018ef64f3c8b561.
//
// Solidity: event UnstakeInit(address indexed user, uint256 indexed validatorId, uint256 deactivationEpoch, uint256 indexed amount)
func (_Stakinginfo *StakinginfoFilterer) ParseUnstakeInit(log types.Log) (*StakinginfoUnstakeInit, error) {
	event := new(StakinginfoUnstakeInit)
	if err := _Stakinginfo.contract.UnpackLog(event, "UnstakeInit", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoUnstakedIterator is returned from FilterUnstaked and is used to iterate over the raw logs and unpacked data for Unstaked events raised by the Stakinginfo contract.
type StakinginfoUnstakedIterator struct {
	Event *StakinginfoUnstaked // Event containing the contract specifics and raw log

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
func (it *StakinginfoUnstakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoUnstaked)
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
		it.Event = new(StakinginfoUnstaked)
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
func (it *StakinginfoUnstakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoUnstakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoUnstaked represents a Unstaked event raised by the Stakinginfo contract.
type StakinginfoUnstaked struct {
	User        common.Address
	ValidatorId *big.Int
	Amount      *big.Int
	Total       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUnstaked is a free log retrieval operation binding the contract event 0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00.
//
// Solidity: event Unstaked(address indexed user, uint256 indexed validatorId, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) FilterUnstaked(opts *bind.FilterOpts, user []common.Address, validatorId []*big.Int) (*StakinginfoUnstakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "Unstaked", userRule, validatorIdRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoUnstakedIterator{contract: _Stakinginfo.contract, event: "Unstaked", logs: logs, sub: sub}, nil
}

// WatchUnstaked is a free log subscription operation binding the contract event 0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00.
//
// Solidity: event Unstaked(address indexed user, uint256 indexed validatorId, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) WatchUnstaked(opts *bind.WatchOpts, sink chan<- *StakinginfoUnstaked, user []common.Address, validatorId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "Unstaked", userRule, validatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoUnstaked)
				if err := _Stakinginfo.contract.UnpackLog(event, "Unstaked", log); err != nil {
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
// Solidity: event Unstaked(address indexed user, uint256 indexed validatorId, uint256 amount, uint256 total)
func (_Stakinginfo *StakinginfoFilterer) ParseUnstaked(log types.Log) (*StakinginfoUnstaked, error) {
	event := new(StakinginfoUnstaked)
	if err := _Stakinginfo.contract.UnpackLog(event, "Unstaked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakinginfoUpdateCommissionRateIterator is returned from FilterUpdateCommissionRate and is used to iterate over the raw logs and unpacked data for UpdateCommissionRate events raised by the Stakinginfo contract.
type StakinginfoUpdateCommissionRateIterator struct {
	Event *StakinginfoUpdateCommissionRate // Event containing the contract specifics and raw log

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
func (it *StakinginfoUpdateCommissionRateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakinginfoUpdateCommissionRate)
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
		it.Event = new(StakinginfoUpdateCommissionRate)
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
func (it *StakinginfoUpdateCommissionRateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakinginfoUpdateCommissionRateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakinginfoUpdateCommissionRate represents a UpdateCommissionRate event raised by the Stakinginfo contract.
type StakinginfoUpdateCommissionRate struct {
	ValidatorId       *big.Int
	NewCommissionRate *big.Int
	OldCommissionRate *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUpdateCommissionRate is a free log retrieval operation binding the contract event 0x7d5da5ece9d43013d62ab966f4704ca376b92be29ca6fbb958154baf1c0dc17e.
//
// Solidity: event UpdateCommissionRate(uint256 indexed validatorId, uint256 indexed newCommissionRate, uint256 indexed oldCommissionRate)
func (_Stakinginfo *StakinginfoFilterer) FilterUpdateCommissionRate(opts *bind.FilterOpts, validatorId []*big.Int, newCommissionRate []*big.Int, oldCommissionRate []*big.Int) (*StakinginfoUpdateCommissionRateIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var newCommissionRateRule []interface{}
	for _, newCommissionRateItem := range newCommissionRate {
		newCommissionRateRule = append(newCommissionRateRule, newCommissionRateItem)
	}
	var oldCommissionRateRule []interface{}
	for _, oldCommissionRateItem := range oldCommissionRate {
		oldCommissionRateRule = append(oldCommissionRateRule, oldCommissionRateItem)
	}

	logs, sub, err := _Stakinginfo.contract.FilterLogs(opts, "UpdateCommissionRate", validatorIdRule, newCommissionRateRule, oldCommissionRateRule)
	if err != nil {
		return nil, err
	}
	return &StakinginfoUpdateCommissionRateIterator{contract: _Stakinginfo.contract, event: "UpdateCommissionRate", logs: logs, sub: sub}, nil
}

// WatchUpdateCommissionRate is a free log subscription operation binding the contract event 0x7d5da5ece9d43013d62ab966f4704ca376b92be29ca6fbb958154baf1c0dc17e.
//
// Solidity: event UpdateCommissionRate(uint256 indexed validatorId, uint256 indexed newCommissionRate, uint256 indexed oldCommissionRate)
func (_Stakinginfo *StakinginfoFilterer) WatchUpdateCommissionRate(opts *bind.WatchOpts, sink chan<- *StakinginfoUpdateCommissionRate, validatorId []*big.Int, newCommissionRate []*big.Int, oldCommissionRate []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var newCommissionRateRule []interface{}
	for _, newCommissionRateItem := range newCommissionRate {
		newCommissionRateRule = append(newCommissionRateRule, newCommissionRateItem)
	}
	var oldCommissionRateRule []interface{}
	for _, oldCommissionRateItem := range oldCommissionRate {
		oldCommissionRateRule = append(oldCommissionRateRule, oldCommissionRateItem)
	}

	logs, sub, err := _Stakinginfo.contract.WatchLogs(opts, "UpdateCommissionRate", validatorIdRule, newCommissionRateRule, oldCommissionRateRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakinginfoUpdateCommissionRate)
				if err := _Stakinginfo.contract.UnpackLog(event, "UpdateCommissionRate", log); err != nil {
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

// ParseUpdateCommissionRate is a log parse operation binding the contract event 0x7d5da5ece9d43013d62ab966f4704ca376b92be29ca6fbb958154baf1c0dc17e.
//
// Solidity: event UpdateCommissionRate(uint256 indexed validatorId, uint256 indexed newCommissionRate, uint256 indexed oldCommissionRate)
func (_Stakinginfo *StakinginfoFilterer) ParseUpdateCommissionRate(log types.Log) (*StakinginfoUpdateCommissionRate, error) {
	event := new(StakinginfoUpdateCommissionRate)
	if err := _Stakinginfo.contract.UnpackLog(event, "UpdateCommissionRate", log); err != nil {
		return nil, err
	}
	return event, nil
}
