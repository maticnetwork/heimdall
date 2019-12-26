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
const StakemanagerABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"getStakerDetails\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"},{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"activationEpoch\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"deactivationEpoch\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"UnstakeInit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"Unstaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"oldSigner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newSigner\",\"type\":\"address\"}],\"name\":\"SignerChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"ReStaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"exitEpoch\",\"type\":\"uint256\"}],\"name\":\"Jailed\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"oldAmount\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"newAmount\",\"type\":\"uint256\"}],\"name\":\"StakeUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"amount\",\"type\":\"uint256\"}],\"name\":\"TopupFees\",\"type\":\"event\"}]"

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

// GetStakerDetails is a free data retrieval call binding the contract method 0x78daaf69.
//
// Solidity: function getStakerDetails(uint256 validatorId) constant returns(uint256, uint256, uint256, address, uint256)
func (_Stakemanager *StakemanagerCaller) GetStakerDetails(opts *bind.CallOpts, validatorId *big.Int) (*big.Int, *big.Int, *big.Int, common.Address, *big.Int, error) {
	var (
		ret0 = new(*big.Int)
		ret1 = new(*big.Int)
		ret2 = new(*big.Int)
		ret3 = new(common.Address)
		ret4 = new(*big.Int)
	)
	out := &[]interface{}{
		ret0,
		ret1,
		ret2,
		ret3,
		ret4,
	}
	err := _Stakemanager.contract.Call(opts, out, "getStakerDetails", validatorId)
	return *ret0, *ret1, *ret2, *ret3, *ret4, err
}

// GetStakerDetails is a free data retrieval call binding the contract method 0x78daaf69.
//
// Solidity: function getStakerDetails(uint256 validatorId) constant returns(uint256, uint256, uint256, address, uint256)
func (_Stakemanager *StakemanagerSession) GetStakerDetails(validatorId *big.Int) (*big.Int, *big.Int, *big.Int, common.Address, *big.Int, error) {
	return _Stakemanager.Contract.GetStakerDetails(&_Stakemanager.CallOpts, validatorId)
}

// GetStakerDetails is a free data retrieval call binding the contract method 0x78daaf69.
//
// Solidity: function getStakerDetails(uint256 validatorId) constant returns(uint256, uint256, uint256, address, uint256)
func (_Stakemanager *StakemanagerCallerSession) GetStakerDetails(validatorId *big.Int) (*big.Int, *big.Int, *big.Int, common.Address, *big.Int, error) {
	return _Stakemanager.Contract.GetStakerDetails(&_Stakemanager.CallOpts, validatorId)
}

// StakemanagerJailedIterator is returned from FilterJailed and is used to iterate over the raw logs and unpacked data for Jailed events raised by the Stakemanager contract.
type StakemanagerJailedIterator struct {
	Event *StakemanagerJailed // Event containing the contract specifics and raw log

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
func (it *StakemanagerJailedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerJailed)
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
		it.Event = new(StakemanagerJailed)
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
func (it *StakemanagerJailedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerJailedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerJailed represents a Jailed event raised by the Stakemanager contract.
type StakemanagerJailed struct {
	ValidatorId *big.Int
	ExitEpoch   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterJailed is a free log retrieval operation binding the contract event 0xa1735a3843d9467dd849a217582720a8af66b9d034712e0b21b0f5ece80670cd.
//
// Solidity: event Jailed(uint256 indexed validatorId, uint256 indexed exitEpoch)
func (_Stakemanager *StakemanagerFilterer) FilterJailed(opts *bind.FilterOpts, validatorId []*big.Int, exitEpoch []*big.Int) (*StakemanagerJailedIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var exitEpochRule []interface{}
	for _, exitEpochItem := range exitEpoch {
		exitEpochRule = append(exitEpochRule, exitEpochItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "Jailed", validatorIdRule, exitEpochRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerJailedIterator{contract: _Stakemanager.contract, event: "Jailed", logs: logs, sub: sub}, nil
}

// WatchJailed is a free log subscription operation binding the contract event 0xa1735a3843d9467dd849a217582720a8af66b9d034712e0b21b0f5ece80670cd.
//
// Solidity: event Jailed(uint256 indexed validatorId, uint256 indexed exitEpoch)
func (_Stakemanager *StakemanagerFilterer) WatchJailed(opts *bind.WatchOpts, sink chan<- *StakemanagerJailed, validatorId []*big.Int, exitEpoch []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var exitEpochRule []interface{}
	for _, exitEpochItem := range exitEpoch {
		exitEpochRule = append(exitEpochRule, exitEpochItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "Jailed", validatorIdRule, exitEpochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerJailed)
				if err := _Stakemanager.contract.UnpackLog(event, "Jailed", log); err != nil {
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
func (_Stakemanager *StakemanagerFilterer) ParseJailed(log types.Log) (*StakemanagerJailed, error) {
	event := new(StakemanagerJailed)
	if err := _Stakemanager.contract.UnpackLog(event, "Jailed", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakemanagerReStakedIterator is returned from FilterReStaked and is used to iterate over the raw logs and unpacked data for ReStaked events raised by the Stakemanager contract.
type StakemanagerReStakedIterator struct {
	Event *StakemanagerReStaked // Event containing the contract specifics and raw log

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
func (it *StakemanagerReStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerReStaked)
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
		it.Event = new(StakemanagerReStaked)
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
func (it *StakemanagerReStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerReStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerReStaked represents a ReStaked event raised by the Stakemanager contract.
type StakemanagerReStaked struct {
	ValidatorId *big.Int
	Amount      *big.Int
	Total       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterReStaked is a free log retrieval operation binding the contract event 0x9cc0e589f20d3310eb2ad571b23529003bd46048d0d1af29277dcf0aa3c398ce.
//
// Solidity: event ReStaked(uint256 indexed validatorId, uint256 amount, uint256 total)
func (_Stakemanager *StakemanagerFilterer) FilterReStaked(opts *bind.FilterOpts, validatorId []*big.Int) (*StakemanagerReStakedIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "ReStaked", validatorIdRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerReStakedIterator{contract: _Stakemanager.contract, event: "ReStaked", logs: logs, sub: sub}, nil
}

// WatchReStaked is a free log subscription operation binding the contract event 0x9cc0e589f20d3310eb2ad571b23529003bd46048d0d1af29277dcf0aa3c398ce.
//
// Solidity: event ReStaked(uint256 indexed validatorId, uint256 amount, uint256 total)
func (_Stakemanager *StakemanagerFilterer) WatchReStaked(opts *bind.WatchOpts, sink chan<- *StakemanagerReStaked, validatorId []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "ReStaked", validatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerReStaked)
				if err := _Stakemanager.contract.UnpackLog(event, "ReStaked", log); err != nil {
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
func (_Stakemanager *StakemanagerFilterer) ParseReStaked(log types.Log) (*StakemanagerReStaked, error) {
	event := new(StakemanagerReStaked)
	if err := _Stakemanager.contract.UnpackLog(event, "ReStaked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakemanagerSignerChangeIterator is returned from FilterSignerChange and is used to iterate over the raw logs and unpacked data for SignerChange events raised by the Stakemanager contract.
type StakemanagerSignerChangeIterator struct {
	Event *StakemanagerSignerChange // Event containing the contract specifics and raw log

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
func (it *StakemanagerSignerChangeIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerSignerChange)
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
		it.Event = new(StakemanagerSignerChange)
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
func (it *StakemanagerSignerChangeIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerSignerChangeIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerSignerChange represents a SignerChange event raised by the Stakemanager contract.
type StakemanagerSignerChange struct {
	ValidatorId *big.Int
	OldSigner   common.Address
	NewSigner   common.Address
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterSignerChange is a free log retrieval operation binding the contract event 0x7dfd3bad1e3cac97d3b89ff06d78394523c4f08fdee4daa71a59160003240c89.
//
// Solidity: event SignerChange(uint256 indexed validatorId, address indexed oldSigner, address indexed newSigner)
func (_Stakemanager *StakemanagerFilterer) FilterSignerChange(opts *bind.FilterOpts, validatorId []*big.Int, oldSigner []common.Address, newSigner []common.Address) (*StakemanagerSignerChangeIterator, error) {

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

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "SignerChange", validatorIdRule, oldSignerRule, newSignerRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerSignerChangeIterator{contract: _Stakemanager.contract, event: "SignerChange", logs: logs, sub: sub}, nil
}

// WatchSignerChange is a free log subscription operation binding the contract event 0x7dfd3bad1e3cac97d3b89ff06d78394523c4f08fdee4daa71a59160003240c89.
//
// Solidity: event SignerChange(uint256 indexed validatorId, address indexed oldSigner, address indexed newSigner)
func (_Stakemanager *StakemanagerFilterer) WatchSignerChange(opts *bind.WatchOpts, sink chan<- *StakemanagerSignerChange, validatorId []*big.Int, oldSigner []common.Address, newSigner []common.Address) (event.Subscription, error) {

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

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "SignerChange", validatorIdRule, oldSignerRule, newSignerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerSignerChange)
				if err := _Stakemanager.contract.UnpackLog(event, "SignerChange", log); err != nil {
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
func (_Stakemanager *StakemanagerFilterer) ParseSignerChange(log types.Log) (*StakemanagerSignerChange, error) {
	event := new(StakemanagerSignerChange)
	if err := _Stakemanager.contract.UnpackLog(event, "SignerChange", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakemanagerStakeUpdateIterator is returned from FilterStakeUpdate and is used to iterate over the raw logs and unpacked data for StakeUpdate events raised by the Stakemanager contract.
type StakemanagerStakeUpdateIterator struct {
	Event *StakemanagerStakeUpdate // Event containing the contract specifics and raw log

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
func (it *StakemanagerStakeUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerStakeUpdate)
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
		it.Event = new(StakemanagerStakeUpdate)
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
func (it *StakemanagerStakeUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerStakeUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerStakeUpdate represents a StakeUpdate event raised by the Stakemanager contract.
type StakemanagerStakeUpdate struct {
	ValidatorId *big.Int
	OldAmount   *big.Int
	NewAmount   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterStakeUpdate is a free log retrieval operation binding the contract event 0x35af9eea1f0e7b300b0a14fae90139a072470e44daa3f14b5069bebbc1265bda.
//
// Solidity: event StakeUpdate(uint256 indexed validatorId, uint256 indexed oldAmount, uint256 indexed newAmount)
func (_Stakemanager *StakemanagerFilterer) FilterStakeUpdate(opts *bind.FilterOpts, validatorId []*big.Int, oldAmount []*big.Int, newAmount []*big.Int) (*StakemanagerStakeUpdateIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var oldAmountRule []interface{}
	for _, oldAmountItem := range oldAmount {
		oldAmountRule = append(oldAmountRule, oldAmountItem)
	}
	var newAmountRule []interface{}
	for _, newAmountItem := range newAmount {
		newAmountRule = append(newAmountRule, newAmountItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "StakeUpdate", validatorIdRule, oldAmountRule, newAmountRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerStakeUpdateIterator{contract: _Stakemanager.contract, event: "StakeUpdate", logs: logs, sub: sub}, nil
}

// WatchStakeUpdate is a free log subscription operation binding the contract event 0x35af9eea1f0e7b300b0a14fae90139a072470e44daa3f14b5069bebbc1265bda.
//
// Solidity: event StakeUpdate(uint256 indexed validatorId, uint256 indexed oldAmount, uint256 indexed newAmount)
func (_Stakemanager *StakemanagerFilterer) WatchStakeUpdate(opts *bind.WatchOpts, sink chan<- *StakemanagerStakeUpdate, validatorId []*big.Int, oldAmount []*big.Int, newAmount []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var oldAmountRule []interface{}
	for _, oldAmountItem := range oldAmount {
		oldAmountRule = append(oldAmountRule, oldAmountItem)
	}
	var newAmountRule []interface{}
	for _, newAmountItem := range newAmount {
		newAmountRule = append(newAmountRule, newAmountItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "StakeUpdate", validatorIdRule, oldAmountRule, newAmountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerStakeUpdate)
				if err := _Stakemanager.contract.UnpackLog(event, "StakeUpdate", log); err != nil {
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

// ParseStakeUpdate is a log parse operation binding the contract event 0x35af9eea1f0e7b300b0a14fae90139a072470e44daa3f14b5069bebbc1265bda.
//
// Solidity: event StakeUpdate(uint256 indexed validatorId, uint256 indexed oldAmount, uint256 indexed newAmount)
func (_Stakemanager *StakemanagerFilterer) ParseStakeUpdate(log types.Log) (*StakemanagerStakeUpdate, error) {
	event := new(StakemanagerStakeUpdate)
	if err := _Stakemanager.contract.UnpackLog(event, "StakeUpdate", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakemanagerStakedIterator is returned from FilterStaked and is used to iterate over the raw logs and unpacked data for Staked events raised by the Stakemanager contract.
type StakemanagerStakedIterator struct {
	Event *StakemanagerStaked // Event containing the contract specifics and raw log

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
func (it *StakemanagerStakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerStaked)
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
		it.Event = new(StakemanagerStaked)
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
func (it *StakemanagerStakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerStakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerStaked represents a Staked event raised by the Stakemanager contract.
type StakemanagerStaked struct {
	User            common.Address
	ValidatorId     *big.Int
	ActivationEpoch *big.Int
	Amount          *big.Int
	Total           *big.Int
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterStaked is a free log retrieval operation binding the contract event 0x9cfd25589d1eb8ad71e342a86a8524e83522e3936c0803048c08f6d9ad974f40.
//
// Solidity: event Staked(address indexed user, uint256 indexed validatorId, uint256 indexed activationEpoch, uint256 amount, uint256 total)
func (_Stakemanager *StakemanagerFilterer) FilterStaked(opts *bind.FilterOpts, user []common.Address, validatorId []*big.Int, activationEpoch []*big.Int) (*StakemanagerStakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var activationEpochRule []interface{}
	for _, activationEpochItem := range activationEpoch {
		activationEpochRule = append(activationEpochRule, activationEpochItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "Staked", userRule, validatorIdRule, activationEpochRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerStakedIterator{contract: _Stakemanager.contract, event: "Staked", logs: logs, sub: sub}, nil
}

// WatchStaked is a free log subscription operation binding the contract event 0x9cfd25589d1eb8ad71e342a86a8524e83522e3936c0803048c08f6d9ad974f40.
//
// Solidity: event Staked(address indexed user, uint256 indexed validatorId, uint256 indexed activationEpoch, uint256 amount, uint256 total)
func (_Stakemanager *StakemanagerFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *StakemanagerStaked, user []common.Address, validatorId []*big.Int, activationEpoch []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var activationEpochRule []interface{}
	for _, activationEpochItem := range activationEpoch {
		activationEpochRule = append(activationEpochRule, activationEpochItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "Staked", userRule, validatorIdRule, activationEpochRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerStaked)
				if err := _Stakemanager.contract.UnpackLog(event, "Staked", log); err != nil {
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
// Solidity: event Staked(address indexed user, uint256 indexed validatorId, uint256 indexed activationEpoch, uint256 amount, uint256 total)
func (_Stakemanager *StakemanagerFilterer) ParseStaked(log types.Log) (*StakemanagerStaked, error) {
	event := new(StakemanagerStaked)
	if err := _Stakemanager.contract.UnpackLog(event, "Staked", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakemanagerTopupFeesIterator is returned from FilterTopupFees and is used to iterate over the raw logs and unpacked data for TopupFees events raised by the Stakemanager contract.
type StakemanagerTopupFeesIterator struct {
	Event *StakemanagerTopupFees // Event containing the contract specifics and raw log

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
func (it *StakemanagerTopupFeesIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerTopupFees)
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
		it.Event = new(StakemanagerTopupFees)
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
func (it *StakemanagerTopupFeesIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerTopupFeesIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerTopupFees represents a TopupFees event raised by the Stakemanager contract.
type StakemanagerTopupFees struct {
	ValidatorId *big.Int
	Amount      *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterTopupFees is a free log retrieval operation binding the contract event 0xf25f8d72205a97cbe9be1a60e7677e81fa177fb43737759544c67ba5ff80cf9b.
//
// Solidity: event TopupFees(uint256 indexed validatorId, uint256 indexed amount)
func (_Stakemanager *StakemanagerFilterer) FilterTopupFees(opts *bind.FilterOpts, validatorId []*big.Int, amount []*big.Int) (*StakemanagerTopupFeesIterator, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "TopupFees", validatorIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerTopupFeesIterator{contract: _Stakemanager.contract, event: "TopupFees", logs: logs, sub: sub}, nil
}

// WatchTopupFees is a free log subscription operation binding the contract event 0xf25f8d72205a97cbe9be1a60e7677e81fa177fb43737759544c67ba5ff80cf9b.
//
// Solidity: event TopupFees(uint256 indexed validatorId, uint256 indexed amount)
func (_Stakemanager *StakemanagerFilterer) WatchTopupFees(opts *bind.WatchOpts, sink chan<- *StakemanagerTopupFees, validatorId []*big.Int, amount []*big.Int) (event.Subscription, error) {

	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}
	var amountRule []interface{}
	for _, amountItem := range amount {
		amountRule = append(amountRule, amountItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "TopupFees", validatorIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerTopupFees)
				if err := _Stakemanager.contract.UnpackLog(event, "TopupFees", log); err != nil {
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

// ParseTopupFees is a log parse operation binding the contract event 0xf25f8d72205a97cbe9be1a60e7677e81fa177fb43737759544c67ba5ff80cf9b.
//
// Solidity: event TopupFees(uint256 indexed validatorId, uint256 indexed amount)
func (_Stakemanager *StakemanagerFilterer) ParseTopupFees(log types.Log) (*StakemanagerTopupFees, error) {
	event := new(StakemanagerTopupFees)
	if err := _Stakemanager.contract.UnpackLog(event, "TopupFees", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakemanagerUnstakeInitIterator is returned from FilterUnstakeInit and is used to iterate over the raw logs and unpacked data for UnstakeInit events raised by the Stakemanager contract.
type StakemanagerUnstakeInitIterator struct {
	Event *StakemanagerUnstakeInit // Event containing the contract specifics and raw log

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
func (it *StakemanagerUnstakeInitIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerUnstakeInit)
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
		it.Event = new(StakemanagerUnstakeInit)
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
func (it *StakemanagerUnstakeInitIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerUnstakeInitIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerUnstakeInit represents a UnstakeInit event raised by the Stakemanager contract.
type StakemanagerUnstakeInit struct {
	User              common.Address
	ValidatorId       *big.Int
	DeactivationEpoch *big.Int
	Amount            *big.Int
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterUnstakeInit is a free log retrieval operation binding the contract event 0x51ff6d8ad000896362d1a3c9de9ad835adb6da92ec3ddee44018ef64f3c8b561.
//
// Solidity: event UnstakeInit(address indexed user, uint256 indexed validatorId, uint256 deactivationEpoch, uint256 indexed amount)
func (_Stakemanager *StakemanagerFilterer) FilterUnstakeInit(opts *bind.FilterOpts, user []common.Address, validatorId []*big.Int, amount []*big.Int) (*StakemanagerUnstakeInitIterator, error) {

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

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "UnstakeInit", userRule, validatorIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerUnstakeInitIterator{contract: _Stakemanager.contract, event: "UnstakeInit", logs: logs, sub: sub}, nil
}

// WatchUnstakeInit is a free log subscription operation binding the contract event 0x51ff6d8ad000896362d1a3c9de9ad835adb6da92ec3ddee44018ef64f3c8b561.
//
// Solidity: event UnstakeInit(address indexed user, uint256 indexed validatorId, uint256 deactivationEpoch, uint256 indexed amount)
func (_Stakemanager *StakemanagerFilterer) WatchUnstakeInit(opts *bind.WatchOpts, sink chan<- *StakemanagerUnstakeInit, user []common.Address, validatorId []*big.Int, amount []*big.Int) (event.Subscription, error) {

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

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "UnstakeInit", userRule, validatorIdRule, amountRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerUnstakeInit)
				if err := _Stakemanager.contract.UnpackLog(event, "UnstakeInit", log); err != nil {
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
func (_Stakemanager *StakemanagerFilterer) ParseUnstakeInit(log types.Log) (*StakemanagerUnstakeInit, error) {
	event := new(StakemanagerUnstakeInit)
	if err := _Stakemanager.contract.UnpackLog(event, "UnstakeInit", log); err != nil {
		return nil, err
	}
	return event, nil
}

// StakemanagerUnstakedIterator is returned from FilterUnstaked and is used to iterate over the raw logs and unpacked data for Unstaked events raised by the Stakemanager contract.
type StakemanagerUnstakedIterator struct {
	Event *StakemanagerUnstaked // Event containing the contract specifics and raw log

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
func (it *StakemanagerUnstakedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StakemanagerUnstaked)
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
		it.Event = new(StakemanagerUnstaked)
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
func (it *StakemanagerUnstakedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StakemanagerUnstakedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StakemanagerUnstaked represents a Unstaked event raised by the Stakemanager contract.
type StakemanagerUnstaked struct {
	User        common.Address
	ValidatorId *big.Int
	Amount      *big.Int
	Total       *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterUnstaked is a free log retrieval operation binding the contract event 0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00.
//
// Solidity: event Unstaked(address indexed user, uint256 indexed validatorId, uint256 amount, uint256 total)
func (_Stakemanager *StakemanagerFilterer) FilterUnstaked(opts *bind.FilterOpts, user []common.Address, validatorId []*big.Int) (*StakemanagerUnstakedIterator, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "Unstaked", userRule, validatorIdRule)
	if err != nil {
		return nil, err
	}
	return &StakemanagerUnstakedIterator{contract: _Stakemanager.contract, event: "Unstaked", logs: logs, sub: sub}, nil
}

// WatchUnstaked is a free log subscription operation binding the contract event 0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00.
//
// Solidity: event Unstaked(address indexed user, uint256 indexed validatorId, uint256 amount, uint256 total)
func (_Stakemanager *StakemanagerFilterer) WatchUnstaked(opts *bind.WatchOpts, sink chan<- *StakemanagerUnstaked, user []common.Address, validatorId []*big.Int) (event.Subscription, error) {

	var userRule []interface{}
	for _, userItem := range user {
		userRule = append(userRule, userItem)
	}
	var validatorIdRule []interface{}
	for _, validatorIdItem := range validatorId {
		validatorIdRule = append(validatorIdRule, validatorIdItem)
	}

	logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "Unstaked", userRule, validatorIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StakemanagerUnstaked)
				if err := _Stakemanager.contract.UnpackLog(event, "Unstaked", log); err != nil {
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
func (_Stakemanager *StakemanagerFilterer) ParseUnstaked(log types.Log) (*StakemanagerUnstaked, error) {
	event := new(StakemanagerUnstaked)
	if err := _Stakemanager.contract.UnpackLog(event, "Unstaked", log); err != nil {
		return nil, err
	}
	return event, nil
}
