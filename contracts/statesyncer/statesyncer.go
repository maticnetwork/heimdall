// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package statesyncer

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

// StatesyncerABI is the input ABI used to generate the binding from.
const StatesyncerABI = "[{\"constant\":false,\"inputs\":[{\"name\":\"contractAddress\",\"type\":\"address\"},{\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"syncState\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"counter\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"contractAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"StateSynced\",\"type\":\"event\"}]"

// Statesyncer is an auto generated Go binding around an Ethereum contract.
type Statesyncer struct {
	StatesyncerCaller     // Read-only binding to the contract
	StatesyncerTransactor // Write-only binding to the contract
	StatesyncerFilterer   // Log filterer for contract events
}

// StatesyncerCaller is an auto generated read-only Go binding around an Ethereum contract.
type StatesyncerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatesyncerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StatesyncerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatesyncerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StatesyncerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatesyncerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StatesyncerSession struct {
	Contract     *Statesyncer      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StatesyncerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StatesyncerCallerSession struct {
	Contract *StatesyncerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// StatesyncerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StatesyncerTransactorSession struct {
	Contract     *StatesyncerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// StatesyncerRaw is an auto generated low-level Go binding around an Ethereum contract.
type StatesyncerRaw struct {
	Contract *Statesyncer // Generic contract binding to access the raw methods on
}

// StatesyncerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StatesyncerCallerRaw struct {
	Contract *StatesyncerCaller // Generic read-only contract binding to access the raw methods on
}

// StatesyncerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StatesyncerTransactorRaw struct {
	Contract *StatesyncerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStatesyncer creates a new instance of Statesyncer, bound to a specific deployed contract.
func NewStatesyncer(address common.Address, backend bind.ContractBackend) (*Statesyncer, error) {
	contract, err := bindStatesyncer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Statesyncer{StatesyncerCaller: StatesyncerCaller{contract: contract}, StatesyncerTransactor: StatesyncerTransactor{contract: contract}, StatesyncerFilterer: StatesyncerFilterer{contract: contract}}, nil
}

// NewStatesyncerCaller creates a new read-only instance of Statesyncer, bound to a specific deployed contract.
func NewStatesyncerCaller(address common.Address, caller bind.ContractCaller) (*StatesyncerCaller, error) {
	contract, err := bindStatesyncer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StatesyncerCaller{contract: contract}, nil
}

// NewStatesyncerTransactor creates a new write-only instance of Statesyncer, bound to a specific deployed contract.
func NewStatesyncerTransactor(address common.Address, transactor bind.ContractTransactor) (*StatesyncerTransactor, error) {
	contract, err := bindStatesyncer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StatesyncerTransactor{contract: contract}, nil
}

// NewStatesyncerFilterer creates a new log filterer instance of Statesyncer, bound to a specific deployed contract.
func NewStatesyncerFilterer(address common.Address, filterer bind.ContractFilterer) (*StatesyncerFilterer, error) {
	contract, err := bindStatesyncer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StatesyncerFilterer{contract: contract}, nil
}

// bindStatesyncer binds a generic wrapper to an already deployed contract.
func bindStatesyncer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StatesyncerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Statesyncer *StatesyncerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Statesyncer.Contract.StatesyncerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Statesyncer *StatesyncerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Statesyncer.Contract.StatesyncerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Statesyncer *StatesyncerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Statesyncer.Contract.StatesyncerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Statesyncer *StatesyncerCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Statesyncer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Statesyncer *StatesyncerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Statesyncer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Statesyncer *StatesyncerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Statesyncer.Contract.contract.Transact(opts, method, params...)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() constant returns(uint256)
func (_Statesyncer *StatesyncerCaller) Counter(opts *bind.CallOpts) (*big.Int, error) {
	var (
		ret0 = new(*big.Int)
	)
	out := ret0
	err := _Statesyncer.contract.Call(opts, out, "counter")
	return *ret0, err
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() constant returns(uint256)
func (_Statesyncer *StatesyncerSession) Counter() (*big.Int, error) {
	return _Statesyncer.Contract.Counter(&_Statesyncer.CallOpts)
}

// Counter is a free data retrieval call binding the contract method 0x61bc221a.
//
// Solidity: function counter() constant returns(uint256)
func (_Statesyncer *StatesyncerCallerSession) Counter() (*big.Int, error) {
	return _Statesyncer.Contract.Counter(&_Statesyncer.CallOpts)
}

// SyncState is a paid mutator transaction binding the contract method 0x16f19831.
//
// Solidity: function syncState(address contractAddress, bytes data) returns()
func (_Statesyncer *StatesyncerTransactor) SyncState(opts *bind.TransactOpts, contractAddress common.Address, data []byte) (*types.Transaction, error) {
	return _Statesyncer.contract.Transact(opts, "syncState", contractAddress, data)
}

// SyncState is a paid mutator transaction binding the contract method 0x16f19831.
//
// Solidity: function syncState(address contractAddress, bytes data) returns()
func (_Statesyncer *StatesyncerSession) SyncState(contractAddress common.Address, data []byte) (*types.Transaction, error) {
	return _Statesyncer.Contract.SyncState(&_Statesyncer.TransactOpts, contractAddress, data)
}

// SyncState is a paid mutator transaction binding the contract method 0x16f19831.
//
// Solidity: function syncState(address contractAddress, bytes data) returns()
func (_Statesyncer *StatesyncerTransactorSession) SyncState(contractAddress common.Address, data []byte) (*types.Transaction, error) {
	return _Statesyncer.Contract.SyncState(&_Statesyncer.TransactOpts, contractAddress, data)
}

// StatesyncerStateSyncedIterator is returned from FilterStateSynced and is used to iterate over the raw logs and unpacked data for StateSynced events raised by the Statesyncer contract.
type StatesyncerStateSyncedIterator struct {
	Event *StatesyncerStateSynced // Event containing the contract specifics and raw log

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
func (it *StatesyncerStateSyncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StatesyncerStateSynced)
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
		it.Event = new(StatesyncerStateSynced)
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
func (it *StatesyncerStateSyncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StatesyncerStateSyncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StatesyncerStateSynced represents a StateSynced event raised by the Statesyncer contract.
type StatesyncerStateSynced struct {
	Id              *big.Int
	ContractAddress common.Address
	Data            []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterStateSynced is a free log retrieval operation binding the contract event 0x103fed9db65eac19c4d870f49ab7520fe03b99f1838e5996caf47e9e43308392.
//
// Solidity: event StateSynced(uint256 indexed id, address indexed contractAddress, bytes data)
func (_Statesyncer *StatesyncerFilterer) FilterStateSynced(opts *bind.FilterOpts, id []*big.Int, contractAddress []common.Address) (*StatesyncerStateSyncedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _Statesyncer.contract.FilterLogs(opts, "StateSynced", idRule, contractAddressRule)
	if err != nil {
		return nil, err
	}
	return &StatesyncerStateSyncedIterator{contract: _Statesyncer.contract, event: "StateSynced", logs: logs, sub: sub}, nil
}

// WatchStateSynced is a free log subscription operation binding the contract event 0x103fed9db65eac19c4d870f49ab7520fe03b99f1838e5996caf47e9e43308392.
//
// Solidity: event StateSynced(uint256 indexed id, address indexed contractAddress, bytes data)
func (_Statesyncer *StatesyncerFilterer) WatchStateSynced(opts *bind.WatchOpts, sink chan<- *StatesyncerStateSynced, id []*big.Int, contractAddress []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _Statesyncer.contract.WatchLogs(opts, "StateSynced", idRule, contractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StatesyncerStateSynced)
				if err := _Statesyncer.contract.UnpackLog(event, "StateSynced", log); err != nil {
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
