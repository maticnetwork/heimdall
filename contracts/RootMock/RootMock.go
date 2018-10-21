// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package rootmock

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

// ContractsABI is the input ABI used to generate the binding from.
const ContractsABI = "[{\"constant\":true,\"inputs\":[],\"name\":\"headerBockSubmitted\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"root\",\"type\":\"bytes32\"},{\"name\":\"start\",\"type\":\"uint256\"},{\"name\":\"end\",\"type\":\"uint256\"},{\"name\":\"sigs\",\"type\":\"bytes\"}],\"name\":\"submitHeaderBlock\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"root\",\"type\":\"bytes32\"},{\"indexed\":false,\"name\":\"start\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"end\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"sigs\",\"type\":\"bytes\"}],\"name\":\"HeaderBlock\",\"type\":\"event\"}]"

// Contracts is an auto generated Go binding around an Ethereum contract.
type Contracts struct {
	ContractsCaller     // Read-only binding to the contract
	ContractsTransactor // Write-only binding to the contract
	ContractsFilterer   // Log filterer for contract events
}

// ContractsCaller is an auto generated read-only Go binding around an Ethereum contract.
type ContractsCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractsTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ContractsTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractsFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ContractsFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ContractsSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ContractsSession struct {
	Contract     *Contracts        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ContractsCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ContractsCallerSession struct {
	Contract *ContractsCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// ContractsTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ContractsTransactorSession struct {
	Contract     *ContractsTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// ContractsRaw is an auto generated low-level Go binding around an Ethereum contract.
type ContractsRaw struct {
	Contract *Contracts // Generic contract binding to access the raw methods on
}

// ContractsCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ContractsCallerRaw struct {
	Contract *ContractsCaller // Generic read-only contract binding to access the raw methods on
}

// ContractsTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ContractsTransactorRaw struct {
	Contract *ContractsTransactor // Generic write-only contract binding to access the raw methods on
}

// NewContracts creates a new instance of Contracts, bound to a specific deployed contract.
func NewContracts(address common.Address, backend bind.ContractBackend) (*Contracts, error) {
	contract, err := bindContracts(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Contracts{ContractsCaller: ContractsCaller{contract: contract}, ContractsTransactor: ContractsTransactor{contract: contract}, ContractsFilterer: ContractsFilterer{contract: contract}}, nil
}

// NewContractsCaller creates a new read-only instance of Contracts, bound to a specific deployed contract.
func NewContractsCaller(address common.Address, caller bind.ContractCaller) (*ContractsCaller, error) {
	contract, err := bindContracts(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ContractsCaller{contract: contract}, nil
}

// NewContractsTransactor creates a new write-only instance of Contracts, bound to a specific deployed contract.
func NewContractsTransactor(address common.Address, transactor bind.ContractTransactor) (*ContractsTransactor, error) {
	contract, err := bindContracts(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ContractsTransactor{contract: contract}, nil
}

// NewContractsFilterer creates a new log filterer instance of Contracts, bound to a specific deployed contract.
func NewContractsFilterer(address common.Address, filterer bind.ContractFilterer) (*ContractsFilterer, error) {
	contract, err := bindContracts(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ContractsFilterer{contract: contract}, nil
}

// bindContracts binds a generic wrapper to an already deployed contract.
func bindContracts(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ContractsABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contracts *ContractsRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contracts.Contract.ContractsCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contracts *ContractsRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contracts.Contract.ContractsTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contracts *ContractsRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contracts.Contract.ContractsTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Contracts *ContractsCallerRaw) Call(opts *bind.CallOpts, result interface{}, method string, params ...interface{}) error {
	return _Contracts.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Contracts *ContractsTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Contracts.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Contracts *ContractsTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Contracts.Contract.contract.Transact(opts, method, params...)
}

// HeaderBockSubmitted is a free data retrieval call binding the contract method 0x33a20dfa.
//
// Solidity: function headerBockSubmitted() constant returns(bool)
func (_Contracts *ContractsCaller) HeaderBockSubmitted(opts *bind.CallOpts) (bool, error) {
	var (
		ret0 = new(bool)
	)
	out := ret0
	err := _Contracts.contract.Call(opts, out, "headerBockSubmitted")
	return *ret0, err
}

// HeaderBockSubmitted is a free data retrieval call binding the contract method 0x33a20dfa.
//
// Solidity: function headerBockSubmitted() constant returns(bool)
func (_Contracts *ContractsSession) HeaderBockSubmitted() (bool, error) {
	return _Contracts.Contract.HeaderBockSubmitted(&_Contracts.CallOpts)
}

// HeaderBockSubmitted is a free data retrieval call binding the contract method 0x33a20dfa.
//
// Solidity: function headerBockSubmitted() constant returns(bool)
func (_Contracts *ContractsCallerSession) HeaderBockSubmitted() (bool, error) {
	return _Contracts.Contract.HeaderBockSubmitted(&_Contracts.CallOpts)
}

// SubmitHeaderBlock is a paid mutator transaction binding the contract method 0xa0ee0f19.
//
// Solidity: function submitHeaderBlock(root bytes32, start uint256, end uint256, sigs bytes) returns(bool)
func (_Contracts *ContractsTransactor) SubmitHeaderBlock(opts *bind.TransactOpts, root [32]byte, start *big.Int, end *big.Int, sigs []byte) (*types.Transaction, error) {
	return _Contracts.contract.Transact(opts, "submitHeaderBlock", root, start, end, sigs)
}

// SubmitHeaderBlock is a paid mutator transaction binding the contract method 0xa0ee0f19.
//
// Solidity: function submitHeaderBlock(root bytes32, start uint256, end uint256, sigs bytes) returns(bool)
func (_Contracts *ContractsSession) SubmitHeaderBlock(root [32]byte, start *big.Int, end *big.Int, sigs []byte) (*types.Transaction, error) {
	return _Contracts.Contract.SubmitHeaderBlock(&_Contracts.TransactOpts, root, start, end, sigs)
}

// SubmitHeaderBlock is a paid mutator transaction binding the contract method 0xa0ee0f19.
//
// Solidity: function submitHeaderBlock(root bytes32, start uint256, end uint256, sigs bytes) returns(bool)
func (_Contracts *ContractsTransactorSession) SubmitHeaderBlock(root [32]byte, start *big.Int, end *big.Int, sigs []byte) (*types.Transaction, error) {
	return _Contracts.Contract.SubmitHeaderBlock(&_Contracts.TransactOpts, root, start, end, sigs)
}

// ContractsHeaderBlockIterator is returned from FilterHeaderBlock and is used to iterate over the raw logs and unpacked data for HeaderBlock events raised by the Contracts contract.
type ContractsHeaderBlockIterator struct {
	Event *ContractsHeaderBlock // Event containing the contract specifics and raw log

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
func (it *ContractsHeaderBlockIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(ContractsHeaderBlock)
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
		it.Event = new(ContractsHeaderBlock)
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
func (it *ContractsHeaderBlockIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *ContractsHeaderBlockIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// ContractsHeaderBlock represents a HeaderBlock event raised by the Contracts contract.
type ContractsHeaderBlock struct {
	Root  [32]byte
	Start *big.Int
	End   *big.Int
	Sigs  []byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterHeaderBlock is a free log retrieval operation binding the contract event 0x4945d56550fd111adbabc9de568e9e0d084719d295b8cb6ce99c77ab71fe9254.
//
// Solidity: e HeaderBlock(root bytes32, start uint256, end uint256, sigs bytes)
func (_Contracts *ContractsFilterer) FilterHeaderBlock(opts *bind.FilterOpts) (*ContractsHeaderBlockIterator, error) {

	logs, sub, err := _Contracts.contract.FilterLogs(opts, "HeaderBlock")
	if err != nil {
		return nil, err
	}
	return &ContractsHeaderBlockIterator{contract: _Contracts.contract, event: "HeaderBlock", logs: logs, sub: sub}, nil
}

// WatchHeaderBlock is a free log subscription operation binding the contract event 0x4945d56550fd111adbabc9de568e9e0d084719d295b8cb6ce99c77ab71fe9254.
//
// Solidity: e HeaderBlock(root bytes32, start uint256, end uint256, sigs bytes)
func (_Contracts *ContractsFilterer) WatchHeaderBlock(opts *bind.WatchOpts, sink chan<- *ContractsHeaderBlock) (event.Subscription, error) {

	logs, sub, err := _Contracts.contract.WatchLogs(opts, "HeaderBlock")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(ContractsHeaderBlock)
				if err := _Contracts.contract.UnpackLog(event, "HeaderBlock", log); err != nil {
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
