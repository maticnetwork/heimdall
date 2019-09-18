// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package statereceiver

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

// StatereceiverABI is the input ABI used to generate the binding from.
const StatereceiverABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"states\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"vote\",\"type\":\"bytes\"},{\"name\":\"sigs\",\"type\":\"bytes\"},{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"commitState\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"validatorSet\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"vote\",\"type\":\"bytes\"},{\"name\":\"sigs\",\"type\":\"bytes\"},{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"validateValidatorSet\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isValidatorSetContract\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"contractAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"NewStateSynced\",\"type\":\"event\"}]"

// Statereceiver is an auto generated Go binding around an Ethereum contract.
type Statereceiver struct {
	StatereceiverCaller     // Read-only binding to the contract
	StatereceiverTransactor // Write-only binding to the contract
	StatereceiverFilterer   // Log filterer for contract events
}

// StatereceiverCaller is an auto generated read-only Go binding around an Ethereum contract.
type StatereceiverCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatereceiverTransactor is an auto generated write-only Go binding around an Ethereum contract.
type StatereceiverTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatereceiverFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type StatereceiverFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// StatereceiverSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type StatereceiverSession struct {
	Contract     *Statereceiver    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// StatereceiverCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type StatereceiverCallerSession struct {
	Contract *StatereceiverCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// StatereceiverTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type StatereceiverTransactorSession struct {
	Contract     *StatereceiverTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// StatereceiverRaw is an auto generated low-level Go binding around an Ethereum contract.
type StatereceiverRaw struct {
	Contract *Statereceiver // Generic contract binding to access the raw methods on
}

// StatereceiverCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type StatereceiverCallerRaw struct {
	Contract *StatereceiverCaller // Generic read-only contract binding to access the raw methods on
}

// StatereceiverTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type StatereceiverTransactorRaw struct {
	Contract *StatereceiverTransactor // Generic write-only contract binding to access the raw methods on
}

// NewStatereceiver creates a new instance of Statereceiver, bound to a specific deployed contract.
func NewStatereceiver(address common.Address, backend bind.ContractBackend) (*Statereceiver, error) {
	contract, err := bindStatereceiver(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Statereceiver{StatereceiverCaller: StatereceiverCaller{contract: contract}, StatereceiverTransactor: StatereceiverTransactor{contract: contract}, StatereceiverFilterer: StatereceiverFilterer{contract: contract}}, nil
}

// NewStatereceiverCaller creates a new read-only instance of Statereceiver, bound to a specific deployed contract.
func NewStatereceiverCaller(address common.Address, caller bind.ContractCaller) (*StatereceiverCaller, error) {
	contract, err := bindStatereceiver(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &StatereceiverCaller{contract: contract}, nil
}

// NewStatereceiverTransactor creates a new write-only instance of Statereceiver, bound to a specific deployed contract.
func NewStatereceiverTransactor(address common.Address, transactor bind.ContractTransactor) (*StatereceiverTransactor, error) {
	contract, err := bindStatereceiver(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &StatereceiverTransactor{contract: contract}, nil
}

// NewStatereceiverFilterer creates a new log filterer instance of Statereceiver, bound to a specific deployed contract.
func NewStatereceiverFilterer(address common.Address, filterer bind.ContractFilterer) (*StatereceiverFilterer, error) {
	contract, err := bindStatereceiver(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &StatereceiverFilterer{contract: contract}, nil
}

// bindStatereceiver binds a generic wrapper to an already deployed contract.
func bindStatereceiver(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(StatereceiverABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Statereceiver *StatereceiverRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Statereceiver.Contract.StatereceiverCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Statereceiver *StatereceiverRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Statereceiver.Contract.StatereceiverTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Statereceiver *StatereceiverRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Statereceiver.Contract.StatereceiverTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Statereceiver *StatereceiverCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Statereceiver.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Statereceiver *StatereceiverTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Statereceiver.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Statereceiver *StatereceiverTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Statereceiver.Contract.contract.Transact(opts, method, params...)
}

// IsValidatorSetContract is a free data retrieval call binding the contract method 0xd79e60b7.
//
// Solidity: function isValidatorSetContract() constant returns(bool)
func (_Statereceiver *StatereceiverCaller) IsValidatorSetContract(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Statereceiver.contract.Call(opts, out, "isValidatorSetContract")
	return *ret0, err
}

// IsValidatorSetContract is a free data retrieval call binding the contract method 0xd79e60b7.
//
// Solidity: function isValidatorSetContract() constant returns(bool)
func (_Statereceiver *StatereceiverSession) IsValidatorSetContract() (bool, error) {
	return _Statereceiver.Contract.IsValidatorSetContract(&_Statereceiver.CallOpts)
}

// IsValidatorSetContract is a free data retrieval call binding the contract method 0xd79e60b7.
//
// Solidity: function isValidatorSetContract() constant returns(bool)
func (_Statereceiver *StatereceiverCallerSession) IsValidatorSetContract() (bool, error) {
	return _Statereceiver.Contract.IsValidatorSetContract(&_Statereceiver.CallOpts)
}

// States is a free data retrieval call binding the contract method 0x017a9105.
//
// Solidity: function states(uint256 ) constant returns(bool)
func (_Statereceiver *StatereceiverCaller) States(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Statereceiver.contract.Call(opts, out, "states", arg0)
	return *ret0, err
}

// States is a free data retrieval call binding the contract method 0x017a9105.
//
// Solidity: function states(uint256 ) constant returns(bool)
func (_Statereceiver *StatereceiverSession) States(arg0 *big.Int) (bool, error) {
	return _Statereceiver.Contract.States(&_Statereceiver.CallOpts, arg0)
}

// States is a free data retrieval call binding the contract method 0x017a9105.
//
// Solidity: function states(uint256 ) constant returns(bool)
func (_Statereceiver *StatereceiverCallerSession) States(arg0 *big.Int) (bool, error) {
	return _Statereceiver.Contract.States(&_Statereceiver.CallOpts, arg0)
}

// ValidatorSet is a free data retrieval call binding the contract method 0x9426e226.
//
// Solidity: function validatorSet() constant returns(address)
func (_Statereceiver *StatereceiverCaller) ValidatorSet(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Statereceiver.contract.Call(opts, out, "validatorSet")
	return *ret0, err
}

// ValidatorSet is a free data retrieval call binding the contract method 0x9426e226.
//
// Solidity: function validatorSet() constant returns(address)
func (_Statereceiver *StatereceiverSession) ValidatorSet() (common.Address, error) {
	return _Statereceiver.Contract.ValidatorSet(&_Statereceiver.CallOpts)
}

// ValidatorSet is a free data retrieval call binding the contract method 0x9426e226.
//
// Solidity: function validatorSet() constant returns(address)
func (_Statereceiver *StatereceiverCallerSession) ValidatorSet() (common.Address, error) {
	return _Statereceiver.Contract.ValidatorSet(&_Statereceiver.CallOpts)
}

// CommitState is a paid mutator transaction binding the contract method 0x29050939.
//
// Solidity: function commitState(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statereceiver *StatereceiverTransactor) CommitState(opts *bind.TransactOpts, vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statereceiver.contract.Transact(opts, "commitState", vote, sigs, txBytes, proof)
}

// CommitState is a paid mutator transaction binding the contract method 0x29050939.
//
// Solidity: function commitState(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statereceiver *StatereceiverSession) CommitState(vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statereceiver.Contract.CommitState(&_Statereceiver.TransactOpts, vote, sigs, txBytes, proof)
}

// CommitState is a paid mutator transaction binding the contract method 0x29050939.
//
// Solidity: function commitState(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statereceiver *StatereceiverTransactorSession) CommitState(vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statereceiver.Contract.CommitState(&_Statereceiver.TransactOpts, vote, sigs, txBytes, proof)
}

// ValidateValidatorSet is a paid mutator transaction binding the contract method 0xd0504f89.
//
// Solidity: function validateValidatorSet(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statereceiver *StatereceiverTransactor) ValidateValidatorSet(opts *bind.TransactOpts, vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statereceiver.contract.Transact(opts, "validateValidatorSet", vote, sigs, txBytes, proof)
}

// ValidateValidatorSet is a paid mutator transaction binding the contract method 0xd0504f89.
//
// Solidity: function validateValidatorSet(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statereceiver *StatereceiverSession) ValidateValidatorSet(vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statereceiver.Contract.ValidateValidatorSet(&_Statereceiver.TransactOpts, vote, sigs, txBytes, proof)
}

// ValidateValidatorSet is a paid mutator transaction binding the contract method 0xd0504f89.
//
// Solidity: function validateValidatorSet(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statereceiver *StatereceiverTransactorSession) ValidateValidatorSet(vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statereceiver.Contract.ValidateValidatorSet(&_Statereceiver.TransactOpts, vote, sigs, txBytes, proof)
}

// StatereceiverNewStateSyncedIterator is returned from FilterNewStateSynced and is used to iterate over the raw logs and unpacked data for NewStateSynced events raised by the Statereceiver contract.
type StatereceiverNewStateSyncedIterator struct {
	Event *StatereceiverNewStateSynced // Event containing the contract specifics and raw log

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
func (it *StatereceiverNewStateSyncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StatereceiverNewStateSynced)
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
		it.Event = new(StatereceiverNewStateSynced)
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
func (it *StatereceiverNewStateSyncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StatereceiverNewStateSyncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StatereceiverNewStateSynced represents a NewStateSynced event raised by the Statereceiver contract.
type StatereceiverNewStateSynced struct {
	Id              *big.Int
	ContractAddress common.Address
	Data            []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterNewStateSynced is a free log retrieval operation binding the contract event 0xccac6e4f23d17f6731617ae0dd91240841f79d5398b423af2621441e4e1a0053.
//
// Solidity: event NewStateSynced(uint256 indexed id, address indexed contractAddress, bytes data)
func (_Statereceiver *StatereceiverFilterer) FilterNewStateSynced(opts *bind.FilterOpts, id []*big.Int, contractAddress []common.Address) (*StatereceiverNewStateSyncedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _Statereceiver.contract.FilterLogs(opts, "NewStateSynced", idRule, contractAddressRule)
	if err != nil {
		return nil, err
	}
	return &StatereceiverNewStateSyncedIterator{contract: _Statereceiver.contract, event: "NewStateSynced", logs: logs, sub: sub}, nil
}

// WatchNewStateSynced is a free log subscription operation binding the contract event 0xccac6e4f23d17f6731617ae0dd91240841f79d5398b423af2621441e4e1a0053.
//
// Solidity: event NewStateSynced(uint256 indexed id, address indexed contractAddress, bytes data)
func (_Statereceiver *StatereceiverFilterer) WatchNewStateSynced(opts *bind.WatchOpts, sink chan<- *StatereceiverNewStateSynced, id []*big.Int, contractAddress []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _Statereceiver.contract.WatchLogs(opts, "NewStateSynced", idRule, contractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StatereceiverNewStateSynced)
				if err := _Statereceiver.contract.UnpackLog(event, "NewStateSynced", log); err != nil {
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
