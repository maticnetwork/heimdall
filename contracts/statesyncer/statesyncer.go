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
const StatesyncerABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"states\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"vote\",\"type\":\"bytes\"},{\"name\":\"sigs\",\"type\":\"bytes\"},{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"commitState\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"validatorSet\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"vote\",\"type\":\"bytes\"},{\"name\":\"sigs\",\"type\":\"bytes\"},{\"name\":\"txBytes\",\"type\":\"bytes\"},{\"name\":\"proof\",\"type\":\"bytes\"}],\"name\":\"validateValidatorSet\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isValidatorSetContract\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"id\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"contractAddress\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"NewStateSynced\",\"type\":\"event\"}]"

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

// IsValidatorSetContract is a free data retrieval call binding the contract method 0xd79e60b7.
//
// Solidity: function isValidatorSetContract() constant returns(bool)
func (_Statesyncer *StatesyncerCaller) IsValidatorSetContract(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Statesyncer.contract.Call(opts, out, "isValidatorSetContract")
	return *ret0, err
}

// IsValidatorSetContract is a free data retrieval call binding the contract method 0xd79e60b7.
//
// Solidity: function isValidatorSetContract() constant returns(bool)
func (_Statesyncer *StatesyncerSession) IsValidatorSetContract() (bool, error) {
	return _Statesyncer.Contract.IsValidatorSetContract(&_Statesyncer.CallOpts)
}

// IsValidatorSetContract is a free data retrieval call binding the contract method 0xd79e60b7.
//
// Solidity: function isValidatorSetContract() constant returns(bool)
func (_Statesyncer *StatesyncerCallerSession) IsValidatorSetContract() (bool, error) {
	return _Statesyncer.Contract.IsValidatorSetContract(&_Statesyncer.CallOpts)
}

// States is a free data retrieval call binding the contract method 0x017a9105.
//
// Solidity: function states(uint256 ) constant returns(bool)
func (_Statesyncer *StatesyncerCaller) States(opts *bind.CallOpts, arg0 *big.Int) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Statesyncer.contract.Call(opts, out, "states", arg0)
	return *ret0, err
}

// States is a free data retrieval call binding the contract method 0x017a9105.
//
// Solidity: function states(uint256 ) constant returns(bool)
func (_Statesyncer *StatesyncerSession) States(arg0 *big.Int) (bool, error) {
	return _Statesyncer.Contract.States(&_Statesyncer.CallOpts, arg0)
}

// States is a free data retrieval call binding the contract method 0x017a9105.
//
// Solidity: function states(uint256 ) constant returns(bool)
func (_Statesyncer *StatesyncerCallerSession) States(arg0 *big.Int) (bool, error) {
	return _Statesyncer.Contract.States(&_Statesyncer.CallOpts, arg0)
}

// ValidatorSet is a free data retrieval call binding the contract method 0x9426e226.
//
// Solidity: function validatorSet() constant returns(address)
func (_Statesyncer *StatesyncerCaller) ValidatorSet(opts *bind.CallOpts) (common.Address, error) {
	var (
		ret0 = new(common.Address)
	)
	out := ret0
	err := _Statesyncer.contract.Call(opts, out, "validatorSet")
	return *ret0, err
}

// ValidatorSet is a free data retrieval call binding the contract method 0x9426e226.
//
// Solidity: function validatorSet() constant returns(address)
func (_Statesyncer *StatesyncerSession) ValidatorSet() (common.Address, error) {
	return _Statesyncer.Contract.ValidatorSet(&_Statesyncer.CallOpts)
}

// ValidatorSet is a free data retrieval call binding the contract method 0x9426e226.
//
// Solidity: function validatorSet() constant returns(address)
func (_Statesyncer *StatesyncerCallerSession) ValidatorSet() (common.Address, error) {
	return _Statesyncer.Contract.ValidatorSet(&_Statesyncer.CallOpts)
}

// CommitState is a paid mutator transaction binding the contract method 0x29050939.
//
// Solidity: function commitState(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statesyncer *StatesyncerTransactor) CommitState(opts *bind.TransactOpts, vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statesyncer.contract.Transact(opts, "commitState", vote, sigs, txBytes, proof)
}

// CommitState is a paid mutator transaction binding the contract method 0x29050939.
//
// Solidity: function commitState(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statesyncer *StatesyncerSession) CommitState(vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statesyncer.Contract.CommitState(&_Statesyncer.TransactOpts, vote, sigs, txBytes, proof)
}

// CommitState is a paid mutator transaction binding the contract method 0x29050939.
//
// Solidity: function commitState(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statesyncer *StatesyncerTransactorSession) CommitState(vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statesyncer.Contract.CommitState(&_Statesyncer.TransactOpts, vote, sigs, txBytes, proof)
}

// ValidateValidatorSet is a paid mutator transaction binding the contract method 0xd0504f89.
//
// Solidity: function validateValidatorSet(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statesyncer *StatesyncerTransactor) ValidateValidatorSet(opts *bind.TransactOpts, vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statesyncer.contract.Transact(opts, "validateValidatorSet", vote, sigs, txBytes, proof)
}

// ValidateValidatorSet is a paid mutator transaction binding the contract method 0xd0504f89.
//
// Solidity: function validateValidatorSet(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statesyncer *StatesyncerSession) ValidateValidatorSet(vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statesyncer.Contract.ValidateValidatorSet(&_Statesyncer.TransactOpts, vote, sigs, txBytes, proof)
}

// ValidateValidatorSet is a paid mutator transaction binding the contract method 0xd0504f89.
//
// Solidity: function validateValidatorSet(bytes vote, bytes sigs, bytes txBytes, bytes proof) returns()
func (_Statesyncer *StatesyncerTransactorSession) ValidateValidatorSet(vote []byte, sigs []byte, txBytes []byte, proof []byte) (*types.Transaction, error) {
	return _Statesyncer.Contract.ValidateValidatorSet(&_Statesyncer.TransactOpts, vote, sigs, txBytes, proof)
}

// StatesyncerNewStateSyncedIterator is returned from FilterNewStateSynced and is used to iterate over the raw logs and unpacked data for NewStateSynced events raised by the Statesyncer contract.
type StatesyncerNewStateSyncedIterator struct {
	Event *StatesyncerNewStateSynced // Event containing the contract specifics and raw log

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
func (it *StatesyncerNewStateSyncedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(StatesyncerNewStateSynced)
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
		it.Event = new(StatesyncerNewStateSynced)
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
func (it *StatesyncerNewStateSyncedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *StatesyncerNewStateSyncedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// StatesyncerNewStateSynced represents a NewStateSynced event raised by the Statesyncer contract.
type StatesyncerNewStateSynced struct {
	Id              *big.Int
	ContractAddress common.Address
	Data            []byte
	Raw             types.Log // Blockchain specific contextual infos
}

// FilterNewStateSynced is a free log retrieval operation binding the contract event 0xccac6e4f23d17f6731617ae0dd91240841f79d5398b423af2621441e4e1a0053.
//
// Solidity: event NewStateSynced(uint256 indexed id, address indexed contractAddress, bytes data)
func (_Statesyncer *StatesyncerFilterer) FilterNewStateSynced(opts *bind.FilterOpts, id []*big.Int, contractAddress []common.Address) (*StatesyncerNewStateSyncedIterator, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _Statesyncer.contract.FilterLogs(opts, "NewStateSynced", idRule, contractAddressRule)
	if err != nil {
		return nil, err
	}
	return &StatesyncerNewStateSyncedIterator{contract: _Statesyncer.contract, event: "NewStateSynced", logs: logs, sub: sub}, nil
}

// WatchNewStateSynced is a free log subscription operation binding the contract event 0xccac6e4f23d17f6731617ae0dd91240841f79d5398b423af2621441e4e1a0053.
//
// Solidity: event NewStateSynced(uint256 indexed id, address indexed contractAddress, bytes data)
func (_Statesyncer *StatesyncerFilterer) WatchNewStateSynced(opts *bind.WatchOpts, sink chan<- *StatesyncerNewStateSynced, id []*big.Int, contractAddress []common.Address) (event.Subscription, error) {

	var idRule []interface{}
	for _, idItem := range id {
		idRule = append(idRule, idItem)
	}
	var contractAddressRule []interface{}
	for _, contractAddressItem := range contractAddress {
		contractAddressRule = append(contractAddressRule, contractAddressItem)
	}

	logs, sub, err := _Statesyncer.contract.WatchLogs(opts, "NewStateSynced", idRule, contractAddressRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(StatesyncerNewStateSynced)
				if err := _Statesyncer.contract.UnpackLog(event, "NewStateSynced", log); err != nil {
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
