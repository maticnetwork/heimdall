// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package stakemanager


	// StakemanagerABI is the input ABI used to generate the binding from.
	const StakemanagerABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"getCurrentValidatorSet\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256[]\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"WITHDRAWAL_DELAY\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newThreshold\",\"type\":\"uint256\"}],\"name\":\"updateValidatorThreshold\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getValidatorId\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalSupply\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"isValidator\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"MIN_DEPOSIT_SIZE\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"unstake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"tokenOfOwnerByIndex\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validators\",\"outputs\":[{\"name\":\"epoch\",\"type\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"reward\",\"type\":\"uint256\"},{\"name\":\"activationEpoch\",\"type\":\"uint256\"},{\"name\":\"deactivationEpoch\",\"type\":\"uint256\"},{\"name\":\"signer\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalizeCommit\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"address\"}],\"name\":\"signerToValidator\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"DYNASTY\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"user\",\"type\":\"address\"}],\"name\":\"totalStakedFor\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"index\",\"type\":\"uint256\"}],\"name\":\"tokenByIndex\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"validatorThreshold\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"NFTCounter\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"validatorState\",\"outputs\":[{\"name\":\"amount\",\"type\":\"int256\"},{\"name\":\"stakerCount\",\"type\":\"int256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"supportsHistory\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"pure\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentEpoch\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"getStakerDetails\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\"},{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"stake\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentValidatorSetSize\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"totalStaked\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"validatorId\",\"type\":\"uint256\"},{\"name\":\"_signer\",\"type\":\"address\"}],\"name\":\"updateSigner\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"isOwner\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"rootChain\",\"outputs\":[{\"name\":\"\",\"type\":\"address\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"epochs\",\"type\":\"uint256\"}],\"name\":\"updateMinLockInPeriod\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"user\",\"type\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\"},{\"name\":\"signer\",\"type\":\"address\"}],\"name\":\"stakeFor\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"currentValidatorSetTotalStake\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"minLockInPeriod\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"unlock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"EPOCH_LENGTH\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"from\",\"type\":\"address\"},{\"name\":\"to\",\"type\":\"address\"},{\"name\":\"tokenId\",\"type\":\"uint256\"},{\"name\":\"_data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"locked\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"validatorId\",\"type\":\"uint256\"}],\"name\":\"unstakeClaim\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"UNSTAKE_DELAY\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newDynasty\",\"type\":\"uint256\"}],\"name\":\"updateDynastyValue\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newRootChain\",\"type\":\"address\"}],\"name\":\"changeRootChain\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"owner\",\"type\":\"address\"},{\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"voteHash\",\"type\":\"bytes32\"},{\"name\":\"sigs\",\"type\":\"bytes\"}],\"name\":\"checkSignatures\",\"outputs\":[{\"name\":\"\",\"type\":\"bool\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"lock\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newThreshold\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"oldThreshold\",\"type\":\"uint256\"}],\"name\":\"ThresholdChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"newDynasty\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"oldDynasty\",\"type\":\"uint256\"}],\"name\":\"DynastyValueChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"deactivationEpoch\",\"type\":\"uint256\"}],\"name\":\"UnstakeInit\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"newSigner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"oldSigner\",\"type\":\"address\"}],\"name\":\"SignerChange\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"name\":\"\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"\",\"type\":\"bytes\"},{\"indexed\":false,\"name\":\"\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"sigsevent\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousRootChain\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newRootChain\",\"type\":\"address\"}],\"name\":\"RootChainChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":true,\"name\":\"activatonEpoch\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"Staked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"user\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"validatorId\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"name\":\"total\",\"type\":\"uint256\"}],\"name\":\"Unstaked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"}]"

	

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
	  Contract     *Stakemanager        // Generic contract binding to set the session for
	  CallOpts     bind.CallOpts     // Call options to use throughout this session
	  TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
	}

	// StakemanagerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
	// with pre-set call options.
	type StakemanagerCallerSession struct {
	  Contract *StakemanagerCaller // Generic contract caller binding to set the session for
	  CallOpts bind.CallOpts    // Call options to use throughout this session
	}

	// StakemanagerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
	// with pre-set transact options.
	type StakemanagerTransactorSession struct {
	  Contract     *StakemanagerTransactor // Generic contract transactor binding to set the session for
	  TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
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
	  return &Stakemanager{ StakemanagerCaller: StakemanagerCaller{contract: contract}, StakemanagerTransactor: StakemanagerTransactor{contract: contract}, StakemanagerFilterer: StakemanagerFilterer{contract: contract} }, nil
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

	
		// DYNASTY is a free data retrieval call binding the contract method 0x485f5b5d.
		//
		// Solidity: function DYNASTY() constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) DYNASTY(opts *bind.CallOpts ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "DYNASTY" )
			return *ret0, err
		}

		// DYNASTY is a free data retrieval call binding the contract method 0x485f5b5d.
		//
		// Solidity: function DYNASTY() constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) DYNASTY() ( *big.Int,  error) {
		  return _Stakemanager.Contract.DYNASTY(&_Stakemanager.CallOpts )
		}

		// DYNASTY is a free data retrieval call binding the contract method 0x485f5b5d.
		//
		// Solidity: function DYNASTY() constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) DYNASTY() ( *big.Int,  error) {
		  return _Stakemanager.Contract.DYNASTY(&_Stakemanager.CallOpts )
		}
	
		// EPOCHLENGTH is a free data retrieval call binding the contract method 0xac4746ab.
		//
		// Solidity: function EPOCH_LENGTH() constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) EPOCHLENGTH(opts *bind.CallOpts ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "EPOCH_LENGTH" )
			return *ret0, err
		}

		// EPOCHLENGTH is a free data retrieval call binding the contract method 0xac4746ab.
		//
		// Solidity: function EPOCH_LENGTH() constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) EPOCHLENGTH() ( *big.Int,  error) {
		  return _Stakemanager.Contract.EPOCHLENGTH(&_Stakemanager.CallOpts )
		}

		// EPOCHLENGTH is a free data retrieval call binding the contract method 0xac4746ab.
		//
		// Solidity: function EPOCH_LENGTH() constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) EPOCHLENGTH() ( *big.Int,  error) {
		  return _Stakemanager.Contract.EPOCHLENGTH(&_Stakemanager.CallOpts )
		}
	
		// MINDEPOSITSIZE is a free data retrieval call binding the contract method 0x26c0817e.
		//
		// Solidity: function MIN_DEPOSIT_SIZE() constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) MINDEPOSITSIZE(opts *bind.CallOpts ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "MIN_DEPOSIT_SIZE" )
			return *ret0, err
		}

		// MINDEPOSITSIZE is a free data retrieval call binding the contract method 0x26c0817e.
		//
		// Solidity: function MIN_DEPOSIT_SIZE() constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) MINDEPOSITSIZE() ( *big.Int,  error) {
		  return _Stakemanager.Contract.MINDEPOSITSIZE(&_Stakemanager.CallOpts )
		}

		// MINDEPOSITSIZE is a free data retrieval call binding the contract method 0x26c0817e.
		//
		// Solidity: function MIN_DEPOSIT_SIZE() constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) MINDEPOSITSIZE() ( *big.Int,  error) {
		  return _Stakemanager.Contract.MINDEPOSITSIZE(&_Stakemanager.CallOpts )
		}
	
		// NFTCounter is a free data retrieval call binding the contract method 0x5508d8e1.
		//
		// Solidity: function NFTCounter() constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) NFTCounter(opts *bind.CallOpts ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "NFTCounter" )
			return *ret0, err
		}

		// NFTCounter is a free data retrieval call binding the contract method 0x5508d8e1.
		//
		// Solidity: function NFTCounter() constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) NFTCounter() ( *big.Int,  error) {
		  return _Stakemanager.Contract.NFTCounter(&_Stakemanager.CallOpts )
		}

		// NFTCounter is a free data retrieval call binding the contract method 0x5508d8e1.
		//
		// Solidity: function NFTCounter() constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) NFTCounter() ( *big.Int,  error) {
		  return _Stakemanager.Contract.NFTCounter(&_Stakemanager.CallOpts )
		}
	
		// UNSTAKEDELAY is a free data retrieval call binding the contract method 0xe35e5d84.
		//
		// Solidity: function UNSTAKE_DELAY() constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) UNSTAKEDELAY(opts *bind.CallOpts ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "UNSTAKE_DELAY" )
			return *ret0, err
		}

		// UNSTAKEDELAY is a free data retrieval call binding the contract method 0xe35e5d84.
		//
		// Solidity: function UNSTAKE_DELAY() constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) UNSTAKEDELAY() ( *big.Int,  error) {
		  return _Stakemanager.Contract.UNSTAKEDELAY(&_Stakemanager.CallOpts )
		}

		// UNSTAKEDELAY is a free data retrieval call binding the contract method 0xe35e5d84.
		//
		// Solidity: function UNSTAKE_DELAY() constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) UNSTAKEDELAY() ( *big.Int,  error) {
		  return _Stakemanager.Contract.UNSTAKEDELAY(&_Stakemanager.CallOpts )
		}
	
		// WITHDRAWALDELAY is a free data retrieval call binding the contract method 0x0ebb172a.
		//
		// Solidity: function WITHDRAWAL_DELAY() constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) WITHDRAWALDELAY(opts *bind.CallOpts ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "WITHDRAWAL_DELAY" )
			return *ret0, err
		}

		// WITHDRAWALDELAY is a free data retrieval call binding the contract method 0x0ebb172a.
		//
		// Solidity: function WITHDRAWAL_DELAY() constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) WITHDRAWALDELAY() ( *big.Int,  error) {
		  return _Stakemanager.Contract.WITHDRAWALDELAY(&_Stakemanager.CallOpts )
		}

		// WITHDRAWALDELAY is a free data retrieval call binding the contract method 0x0ebb172a.
		//
		// Solidity: function WITHDRAWAL_DELAY() constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) WITHDRAWALDELAY() ( *big.Int,  error) {
		  return _Stakemanager.Contract.WITHDRAWALDELAY(&_Stakemanager.CallOpts )
		}
	
		// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
		//
		// Solidity: function balanceOf(owner address) constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) BalanceOf(opts *bind.CallOpts , owner common.Address ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "balanceOf" , owner)
			return *ret0, err
		}

		// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
		//
		// Solidity: function balanceOf(owner address) constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) BalanceOf( owner common.Address ) ( *big.Int,  error) {
		  return _Stakemanager.Contract.BalanceOf(&_Stakemanager.CallOpts , owner)
		}

		// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
		//
		// Solidity: function balanceOf(owner address) constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) BalanceOf( owner common.Address ) ( *big.Int,  error) {
		  return _Stakemanager.Contract.BalanceOf(&_Stakemanager.CallOpts , owner)
		}
	
		// CheckSignatures is a free data retrieval call binding the contract method 0xed516d51.
		//
		// Solidity: function checkSignatures(voteHash bytes32, sigs bytes) constant returns(bool)
		func (_Stakemanager *StakemanagerCaller) CheckSignatures(opts *bind.CallOpts , voteHash [32]byte , sigs []byte ) (bool, error) {
			var (
				ret0 = new(bool)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "checkSignatures" , voteHash, sigs)
			return *ret0, err
		}

		// CheckSignatures is a free data retrieval call binding the contract method 0xed516d51.
		//
		// Solidity: function checkSignatures(voteHash bytes32, sigs bytes) constant returns(bool)
		func (_Stakemanager *StakemanagerSession) CheckSignatures( voteHash [32]byte , sigs []byte ) ( bool,  error) {
		  return _Stakemanager.Contract.CheckSignatures(&_Stakemanager.CallOpts , voteHash, sigs)
		}

		// CheckSignatures is a free data retrieval call binding the contract method 0xed516d51.
		//
		// Solidity: function checkSignatures(voteHash bytes32, sigs bytes) constant returns(bool)
		func (_Stakemanager *StakemanagerCallerSession) CheckSignatures( voteHash [32]byte , sigs []byte ) ( bool,  error) {
		  return _Stakemanager.Contract.CheckSignatures(&_Stakemanager.CallOpts , voteHash, sigs)
		}
	
		// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
		//
		// Solidity: function currentEpoch() constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) CurrentEpoch(opts *bind.CallOpts ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "currentEpoch" )
			return *ret0, err
		}

		// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
		//
		// Solidity: function currentEpoch() constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) CurrentEpoch() ( *big.Int,  error) {
		  return _Stakemanager.Contract.CurrentEpoch(&_Stakemanager.CallOpts )
		}

		// CurrentEpoch is a free data retrieval call binding the contract method 0x76671808.
		//
		// Solidity: function currentEpoch() constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) CurrentEpoch() ( *big.Int,  error) {
		  return _Stakemanager.Contract.CurrentEpoch(&_Stakemanager.CallOpts )
		}
	
		// CurrentValidatorSetSize is a free data retrieval call binding the contract method 0x7f952d95.
		//
		// Solidity: function currentValidatorSetSize() constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) CurrentValidatorSetSize(opts *bind.CallOpts ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "currentValidatorSetSize" )
			return *ret0, err
		}

		// CurrentValidatorSetSize is a free data retrieval call binding the contract method 0x7f952d95.
		//
		// Solidity: function currentValidatorSetSize() constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) CurrentValidatorSetSize() ( *big.Int,  error) {
		  return _Stakemanager.Contract.CurrentValidatorSetSize(&_Stakemanager.CallOpts )
		}

		// CurrentValidatorSetSize is a free data retrieval call binding the contract method 0x7f952d95.
		//
		// Solidity: function currentValidatorSetSize() constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) CurrentValidatorSetSize() ( *big.Int,  error) {
		  return _Stakemanager.Contract.CurrentValidatorSetSize(&_Stakemanager.CallOpts )
		}
	
		// CurrentValidatorSetTotalStake is a free data retrieval call binding the contract method 0xa4769071.
		//
		// Solidity: function currentValidatorSetTotalStake() constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) CurrentValidatorSetTotalStake(opts *bind.CallOpts ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "currentValidatorSetTotalStake" )
			return *ret0, err
		}

		// CurrentValidatorSetTotalStake is a free data retrieval call binding the contract method 0xa4769071.
		//
		// Solidity: function currentValidatorSetTotalStake() constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) CurrentValidatorSetTotalStake() ( *big.Int,  error) {
		  return _Stakemanager.Contract.CurrentValidatorSetTotalStake(&_Stakemanager.CallOpts )
		}

		// CurrentValidatorSetTotalStake is a free data retrieval call binding the contract method 0xa4769071.
		//
		// Solidity: function currentValidatorSetTotalStake() constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) CurrentValidatorSetTotalStake() ( *big.Int,  error) {
		  return _Stakemanager.Contract.CurrentValidatorSetTotalStake(&_Stakemanager.CallOpts )
		}
	
		// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
		//
		// Solidity: function getApproved(tokenId uint256) constant returns(address)
		func (_Stakemanager *StakemanagerCaller) GetApproved(opts *bind.CallOpts , tokenId *big.Int ) (common.Address, error) {
			var (
				ret0 = new(common.Address)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "getApproved" , tokenId)
			return *ret0, err
		}

		// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
		//
		// Solidity: function getApproved(tokenId uint256) constant returns(address)
		func (_Stakemanager *StakemanagerSession) GetApproved( tokenId *big.Int ) ( common.Address,  error) {
		  return _Stakemanager.Contract.GetApproved(&_Stakemanager.CallOpts , tokenId)
		}

		// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
		//
		// Solidity: function getApproved(tokenId uint256) constant returns(address)
		func (_Stakemanager *StakemanagerCallerSession) GetApproved( tokenId *big.Int ) ( common.Address,  error) {
		  return _Stakemanager.Contract.GetApproved(&_Stakemanager.CallOpts , tokenId)
		}
	
		// GetCurrentValidatorSet is a free data retrieval call binding the contract method 0x0209fdd0.
		//
		// Solidity: function getCurrentValidatorSet() constant returns(uint256[])
		func (_Stakemanager *StakemanagerCaller) GetCurrentValidatorSet(opts *bind.CallOpts ) ([]*big.Int, error) {
			var (
				ret0 = new([]*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "getCurrentValidatorSet" )
			return *ret0, err
		}

		// GetCurrentValidatorSet is a free data retrieval call binding the contract method 0x0209fdd0.
		//
		// Solidity: function getCurrentValidatorSet() constant returns(uint256[])
		func (_Stakemanager *StakemanagerSession) GetCurrentValidatorSet() ( []*big.Int,  error) {
		  return _Stakemanager.Contract.GetCurrentValidatorSet(&_Stakemanager.CallOpts )
		}

		// GetCurrentValidatorSet is a free data retrieval call binding the contract method 0x0209fdd0.
		//
		// Solidity: function getCurrentValidatorSet() constant returns(uint256[])
		func (_Stakemanager *StakemanagerCallerSession) GetCurrentValidatorSet() ( []*big.Int,  error) {
		  return _Stakemanager.Contract.GetCurrentValidatorSet(&_Stakemanager.CallOpts )
		}
	
		// GetStakerDetails is a free data retrieval call binding the contract method 0x78daaf69.
		//
		// Solidity: function getStakerDetails(validatorId uint256) constant returns(uint256, uint256, uint256, address)
		func (_Stakemanager *StakemanagerCaller) GetStakerDetails(opts *bind.CallOpts , validatorId *big.Int ) (*big.Int,*big.Int,*big.Int,common.Address, error) {
			var (
				ret0 = new(*big.Int)
				ret1 = new(*big.Int)
				ret2 = new(*big.Int)
				ret3 = new(common.Address)
				
			)
			out := &[]interface{}{
				ret0,
				ret1,
				ret2,
				ret3,
				
			}
			err := _Stakemanager.contract.Call(opts, out, "getStakerDetails" , validatorId)
			return *ret0,*ret1,*ret2,*ret3, err
		}

		// GetStakerDetails is a free data retrieval call binding the contract method 0x78daaf69.
		//
		// Solidity: function getStakerDetails(validatorId uint256) constant returns(uint256, uint256, uint256, address)
		func (_Stakemanager *StakemanagerSession) GetStakerDetails( validatorId *big.Int ) ( *big.Int,*big.Int,*big.Int,common.Address,  error) {
		  return _Stakemanager.Contract.GetStakerDetails(&_Stakemanager.CallOpts , validatorId)
		}

		// GetStakerDetails is a free data retrieval call binding the contract method 0x78daaf69.
		//
		// Solidity: function getStakerDetails(validatorId uint256) constant returns(uint256, uint256, uint256, address)
		func (_Stakemanager *StakemanagerCallerSession) GetStakerDetails( validatorId *big.Int ) ( *big.Int,*big.Int,*big.Int,common.Address,  error) {
		  return _Stakemanager.Contract.GetStakerDetails(&_Stakemanager.CallOpts , validatorId)
		}
	
		// GetValidatorId is a free data retrieval call binding the contract method 0x174e6832.
		//
		// Solidity: function getValidatorId(user address) constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) GetValidatorId(opts *bind.CallOpts , user common.Address ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "getValidatorId" , user)
			return *ret0, err
		}

		// GetValidatorId is a free data retrieval call binding the contract method 0x174e6832.
		//
		// Solidity: function getValidatorId(user address) constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) GetValidatorId( user common.Address ) ( *big.Int,  error) {
		  return _Stakemanager.Contract.GetValidatorId(&_Stakemanager.CallOpts , user)
		}

		// GetValidatorId is a free data retrieval call binding the contract method 0x174e6832.
		//
		// Solidity: function getValidatorId(user address) constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) GetValidatorId( user common.Address ) ( *big.Int,  error) {
		  return _Stakemanager.Contract.GetValidatorId(&_Stakemanager.CallOpts , user)
		}
	
		// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
		//
		// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
		func (_Stakemanager *StakemanagerCaller) IsApprovedForAll(opts *bind.CallOpts , owner common.Address , operator common.Address ) (bool, error) {
			var (
				ret0 = new(bool)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "isApprovedForAll" , owner, operator)
			return *ret0, err
		}

		// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
		//
		// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
		func (_Stakemanager *StakemanagerSession) IsApprovedForAll( owner common.Address , operator common.Address ) ( bool,  error) {
		  return _Stakemanager.Contract.IsApprovedForAll(&_Stakemanager.CallOpts , owner, operator)
		}

		// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
		//
		// Solidity: function isApprovedForAll(owner address, operator address) constant returns(bool)
		func (_Stakemanager *StakemanagerCallerSession) IsApprovedForAll( owner common.Address , operator common.Address ) ( bool,  error) {
		  return _Stakemanager.Contract.IsApprovedForAll(&_Stakemanager.CallOpts , owner, operator)
		}
	
		// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
		//
		// Solidity: function isOwner() constant returns(bool)
		func (_Stakemanager *StakemanagerCaller) IsOwner(opts *bind.CallOpts ) (bool, error) {
			var (
				ret0 = new(bool)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "isOwner" )
			return *ret0, err
		}

		// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
		//
		// Solidity: function isOwner() constant returns(bool)
		func (_Stakemanager *StakemanagerSession) IsOwner() ( bool,  error) {
		  return _Stakemanager.Contract.IsOwner(&_Stakemanager.CallOpts )
		}

		// IsOwner is a free data retrieval call binding the contract method 0x8f32d59b.
		//
		// Solidity: function isOwner() constant returns(bool)
		func (_Stakemanager *StakemanagerCallerSession) IsOwner() ( bool,  error) {
		  return _Stakemanager.Contract.IsOwner(&_Stakemanager.CallOpts )
		}
	
		// IsValidator is a free data retrieval call binding the contract method 0x2649263a.
		//
		// Solidity: function isValidator(validatorId uint256) constant returns(bool)
		func (_Stakemanager *StakemanagerCaller) IsValidator(opts *bind.CallOpts , validatorId *big.Int ) (bool, error) {
			var (
				ret0 = new(bool)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "isValidator" , validatorId)
			return *ret0, err
		}

		// IsValidator is a free data retrieval call binding the contract method 0x2649263a.
		//
		// Solidity: function isValidator(validatorId uint256) constant returns(bool)
		func (_Stakemanager *StakemanagerSession) IsValidator( validatorId *big.Int ) ( bool,  error) {
		  return _Stakemanager.Contract.IsValidator(&_Stakemanager.CallOpts , validatorId)
		}

		// IsValidator is a free data retrieval call binding the contract method 0x2649263a.
		//
		// Solidity: function isValidator(validatorId uint256) constant returns(bool)
		func (_Stakemanager *StakemanagerCallerSession) IsValidator( validatorId *big.Int ) ( bool,  error) {
		  return _Stakemanager.Contract.IsValidator(&_Stakemanager.CallOpts , validatorId)
		}
	
		// Locked is a free data retrieval call binding the contract method 0xcf309012.
		//
		// Solidity: function locked() constant returns(bool)
		func (_Stakemanager *StakemanagerCaller) Locked(opts *bind.CallOpts ) (bool, error) {
			var (
				ret0 = new(bool)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "locked" )
			return *ret0, err
		}

		// Locked is a free data retrieval call binding the contract method 0xcf309012.
		//
		// Solidity: function locked() constant returns(bool)
		func (_Stakemanager *StakemanagerSession) Locked() ( bool,  error) {
		  return _Stakemanager.Contract.Locked(&_Stakemanager.CallOpts )
		}

		// Locked is a free data retrieval call binding the contract method 0xcf309012.
		//
		// Solidity: function locked() constant returns(bool)
		func (_Stakemanager *StakemanagerCallerSession) Locked() ( bool,  error) {
		  return _Stakemanager.Contract.Locked(&_Stakemanager.CallOpts )
		}
	
		// MinLockInPeriod is a free data retrieval call binding the contract method 0xa548c547.
		//
		// Solidity: function minLockInPeriod() constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) MinLockInPeriod(opts *bind.CallOpts ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "minLockInPeriod" )
			return *ret0, err
		}

		// MinLockInPeriod is a free data retrieval call binding the contract method 0xa548c547.
		//
		// Solidity: function minLockInPeriod() constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) MinLockInPeriod() ( *big.Int,  error) {
		  return _Stakemanager.Contract.MinLockInPeriod(&_Stakemanager.CallOpts )
		}

		// MinLockInPeriod is a free data retrieval call binding the contract method 0xa548c547.
		//
		// Solidity: function minLockInPeriod() constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) MinLockInPeriod() ( *big.Int,  error) {
		  return _Stakemanager.Contract.MinLockInPeriod(&_Stakemanager.CallOpts )
		}
	
		// Name is a free data retrieval call binding the contract method 0x06fdde03.
		//
		// Solidity: function name() constant returns(string)
		func (_Stakemanager *StakemanagerCaller) Name(opts *bind.CallOpts ) (string, error) {
			var (
				ret0 = new(string)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "name" )
			return *ret0, err
		}

		// Name is a free data retrieval call binding the contract method 0x06fdde03.
		//
		// Solidity: function name() constant returns(string)
		func (_Stakemanager *StakemanagerSession) Name() ( string,  error) {
		  return _Stakemanager.Contract.Name(&_Stakemanager.CallOpts )
		}

		// Name is a free data retrieval call binding the contract method 0x06fdde03.
		//
		// Solidity: function name() constant returns(string)
		func (_Stakemanager *StakemanagerCallerSession) Name() ( string,  error) {
		  return _Stakemanager.Contract.Name(&_Stakemanager.CallOpts )
		}
	
		// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
		//
		// Solidity: function owner() constant returns(address)
		func (_Stakemanager *StakemanagerCaller) Owner(opts *bind.CallOpts ) (common.Address, error) {
			var (
				ret0 = new(common.Address)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "owner" )
			return *ret0, err
		}

		// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
		//
		// Solidity: function owner() constant returns(address)
		func (_Stakemanager *StakemanagerSession) Owner() ( common.Address,  error) {
		  return _Stakemanager.Contract.Owner(&_Stakemanager.CallOpts )
		}

		// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
		//
		// Solidity: function owner() constant returns(address)
		func (_Stakemanager *StakemanagerCallerSession) Owner() ( common.Address,  error) {
		  return _Stakemanager.Contract.Owner(&_Stakemanager.CallOpts )
		}
	
		// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
		//
		// Solidity: function ownerOf(tokenId uint256) constant returns(address)
		func (_Stakemanager *StakemanagerCaller) OwnerOf(opts *bind.CallOpts , tokenId *big.Int ) (common.Address, error) {
			var (
				ret0 = new(common.Address)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "ownerOf" , tokenId)
			return *ret0, err
		}

		// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
		//
		// Solidity: function ownerOf(tokenId uint256) constant returns(address)
		func (_Stakemanager *StakemanagerSession) OwnerOf( tokenId *big.Int ) ( common.Address,  error) {
		  return _Stakemanager.Contract.OwnerOf(&_Stakemanager.CallOpts , tokenId)
		}

		// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
		//
		// Solidity: function ownerOf(tokenId uint256) constant returns(address)
		func (_Stakemanager *StakemanagerCallerSession) OwnerOf( tokenId *big.Int ) ( common.Address,  error) {
		  return _Stakemanager.Contract.OwnerOf(&_Stakemanager.CallOpts , tokenId)
		}
	
		// RootChain is a free data retrieval call binding the contract method 0x987ab9db.
		//
		// Solidity: function rootChain() constant returns(address)
		func (_Stakemanager *StakemanagerCaller) RootChain(opts *bind.CallOpts ) (common.Address, error) {
			var (
				ret0 = new(common.Address)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "rootChain" )
			return *ret0, err
		}

		// RootChain is a free data retrieval call binding the contract method 0x987ab9db.
		//
		// Solidity: function rootChain() constant returns(address)
		func (_Stakemanager *StakemanagerSession) RootChain() ( common.Address,  error) {
		  return _Stakemanager.Contract.RootChain(&_Stakemanager.CallOpts )
		}

		// RootChain is a free data retrieval call binding the contract method 0x987ab9db.
		//
		// Solidity: function rootChain() constant returns(address)
		func (_Stakemanager *StakemanagerCallerSession) RootChain() ( common.Address,  error) {
		  return _Stakemanager.Contract.RootChain(&_Stakemanager.CallOpts )
		}
	
		// SignerToValidator is a free data retrieval call binding the contract method 0x3862da0b.
		//
		// Solidity: function signerToValidator( address) constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) SignerToValidator(opts *bind.CallOpts , arg0 common.Address ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "signerToValidator" , arg0)
			return *ret0, err
		}

		// SignerToValidator is a free data retrieval call binding the contract method 0x3862da0b.
		//
		// Solidity: function signerToValidator( address) constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) SignerToValidator( arg0 common.Address ) ( *big.Int,  error) {
		  return _Stakemanager.Contract.SignerToValidator(&_Stakemanager.CallOpts , arg0)
		}

		// SignerToValidator is a free data retrieval call binding the contract method 0x3862da0b.
		//
		// Solidity: function signerToValidator( address) constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) SignerToValidator( arg0 common.Address ) ( *big.Int,  error) {
		  return _Stakemanager.Contract.SignerToValidator(&_Stakemanager.CallOpts , arg0)
		}
	
		// SupportsHistory is a free data retrieval call binding the contract method 0x7033e4a6.
		//
		// Solidity: function supportsHistory() constant returns(bool)
		func (_Stakemanager *StakemanagerCaller) SupportsHistory(opts *bind.CallOpts ) (bool, error) {
			var (
				ret0 = new(bool)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "supportsHistory" )
			return *ret0, err
		}

		// SupportsHistory is a free data retrieval call binding the contract method 0x7033e4a6.
		//
		// Solidity: function supportsHistory() constant returns(bool)
		func (_Stakemanager *StakemanagerSession) SupportsHistory() ( bool,  error) {
		  return _Stakemanager.Contract.SupportsHistory(&_Stakemanager.CallOpts )
		}

		// SupportsHistory is a free data retrieval call binding the contract method 0x7033e4a6.
		//
		// Solidity: function supportsHistory() constant returns(bool)
		func (_Stakemanager *StakemanagerCallerSession) SupportsHistory() ( bool,  error) {
		  return _Stakemanager.Contract.SupportsHistory(&_Stakemanager.CallOpts )
		}
	
		// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
		//
		// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
		func (_Stakemanager *StakemanagerCaller) SupportsInterface(opts *bind.CallOpts , interfaceId [4]byte ) (bool, error) {
			var (
				ret0 = new(bool)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "supportsInterface" , interfaceId)
			return *ret0, err
		}

		// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
		//
		// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
		func (_Stakemanager *StakemanagerSession) SupportsInterface( interfaceId [4]byte ) ( bool,  error) {
		  return _Stakemanager.Contract.SupportsInterface(&_Stakemanager.CallOpts , interfaceId)
		}

		// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
		//
		// Solidity: function supportsInterface(interfaceId bytes4) constant returns(bool)
		func (_Stakemanager *StakemanagerCallerSession) SupportsInterface( interfaceId [4]byte ) ( bool,  error) {
		  return _Stakemanager.Contract.SupportsInterface(&_Stakemanager.CallOpts , interfaceId)
		}
	
		// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
		//
		// Solidity: function symbol() constant returns(string)
		func (_Stakemanager *StakemanagerCaller) Symbol(opts *bind.CallOpts ) (string, error) {
			var (
				ret0 = new(string)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "symbol" )
			return *ret0, err
		}

		// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
		//
		// Solidity: function symbol() constant returns(string)
		func (_Stakemanager *StakemanagerSession) Symbol() ( string,  error) {
		  return _Stakemanager.Contract.Symbol(&_Stakemanager.CallOpts )
		}

		// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
		//
		// Solidity: function symbol() constant returns(string)
		func (_Stakemanager *StakemanagerCallerSession) Symbol() ( string,  error) {
		  return _Stakemanager.Contract.Symbol(&_Stakemanager.CallOpts )
		}
	
		// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
		//
		// Solidity: function tokenByIndex(index uint256) constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) TokenByIndex(opts *bind.CallOpts , index *big.Int ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "tokenByIndex" , index)
			return *ret0, err
		}

		// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
		//
		// Solidity: function tokenByIndex(index uint256) constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) TokenByIndex( index *big.Int ) ( *big.Int,  error) {
		  return _Stakemanager.Contract.TokenByIndex(&_Stakemanager.CallOpts , index)
		}

		// TokenByIndex is a free data retrieval call binding the contract method 0x4f6ccce7.
		//
		// Solidity: function tokenByIndex(index uint256) constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) TokenByIndex( index *big.Int ) ( *big.Int,  error) {
		  return _Stakemanager.Contract.TokenByIndex(&_Stakemanager.CallOpts , index)
		}
	
		// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
		//
		// Solidity: function tokenOfOwnerByIndex(owner address, index uint256) constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) TokenOfOwnerByIndex(opts *bind.CallOpts , owner common.Address , index *big.Int ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "tokenOfOwnerByIndex" , owner, index)
			return *ret0, err
		}

		// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
		//
		// Solidity: function tokenOfOwnerByIndex(owner address, index uint256) constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) TokenOfOwnerByIndex( owner common.Address , index *big.Int ) ( *big.Int,  error) {
		  return _Stakemanager.Contract.TokenOfOwnerByIndex(&_Stakemanager.CallOpts , owner, index)
		}

		// TokenOfOwnerByIndex is a free data retrieval call binding the contract method 0x2f745c59.
		//
		// Solidity: function tokenOfOwnerByIndex(owner address, index uint256) constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) TokenOfOwnerByIndex( owner common.Address , index *big.Int ) ( *big.Int,  error) {
		  return _Stakemanager.Contract.TokenOfOwnerByIndex(&_Stakemanager.CallOpts , owner, index)
		}
	
		// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
		//
		// Solidity: function tokenURI(tokenId uint256) constant returns(string)
		func (_Stakemanager *StakemanagerCaller) TokenURI(opts *bind.CallOpts , tokenId *big.Int ) (string, error) {
			var (
				ret0 = new(string)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "tokenURI" , tokenId)
			return *ret0, err
		}

		// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
		//
		// Solidity: function tokenURI(tokenId uint256) constant returns(string)
		func (_Stakemanager *StakemanagerSession) TokenURI( tokenId *big.Int ) ( string,  error) {
		  return _Stakemanager.Contract.TokenURI(&_Stakemanager.CallOpts , tokenId)
		}

		// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
		//
		// Solidity: function tokenURI(tokenId uint256) constant returns(string)
		func (_Stakemanager *StakemanagerCallerSession) TokenURI( tokenId *big.Int ) ( string,  error) {
		  return _Stakemanager.Contract.TokenURI(&_Stakemanager.CallOpts , tokenId)
		}
	
		// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
		//
		// Solidity: function totalStaked() constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) TotalStaked(opts *bind.CallOpts ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "totalStaked" )
			return *ret0, err
		}

		// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
		//
		// Solidity: function totalStaked() constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) TotalStaked() ( *big.Int,  error) {
		  return _Stakemanager.Contract.TotalStaked(&_Stakemanager.CallOpts )
		}

		// TotalStaked is a free data retrieval call binding the contract method 0x817b1cd2.
		//
		// Solidity: function totalStaked() constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) TotalStaked() ( *big.Int,  error) {
		  return _Stakemanager.Contract.TotalStaked(&_Stakemanager.CallOpts )
		}
	
		// TotalStakedFor is a free data retrieval call binding the contract method 0x4b341aed.
		//
		// Solidity: function totalStakedFor(user address) constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) TotalStakedFor(opts *bind.CallOpts , user common.Address ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "totalStakedFor" , user)
			return *ret0, err
		}

		// TotalStakedFor is a free data retrieval call binding the contract method 0x4b341aed.
		//
		// Solidity: function totalStakedFor(user address) constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) TotalStakedFor( user common.Address ) ( *big.Int,  error) {
		  return _Stakemanager.Contract.TotalStakedFor(&_Stakemanager.CallOpts , user)
		}

		// TotalStakedFor is a free data retrieval call binding the contract method 0x4b341aed.
		//
		// Solidity: function totalStakedFor(user address) constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) TotalStakedFor( user common.Address ) ( *big.Int,  error) {
		  return _Stakemanager.Contract.TotalStakedFor(&_Stakemanager.CallOpts , user)
		}
	
		// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
		//
		// Solidity: function totalSupply() constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) TotalSupply(opts *bind.CallOpts ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "totalSupply" )
			return *ret0, err
		}

		// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
		//
		// Solidity: function totalSupply() constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) TotalSupply() ( *big.Int,  error) {
		  return _Stakemanager.Contract.TotalSupply(&_Stakemanager.CallOpts )
		}

		// TotalSupply is a free data retrieval call binding the contract method 0x18160ddd.
		//
		// Solidity: function totalSupply() constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) TotalSupply() ( *big.Int,  error) {
		  return _Stakemanager.Contract.TotalSupply(&_Stakemanager.CallOpts )
		}
	
		// ValidatorState is a free data retrieval call binding the contract method 0x5c248855.
		//
		// Solidity: function validatorState( uint256) constant returns(amount int256, stakerCount int256)
		func (_Stakemanager *StakemanagerCaller) ValidatorState(opts *bind.CallOpts , arg0 *big.Int ) (struct{ Amount *big.Int;StakerCount *big.Int; }, error) {
			ret := new(struct{
				Amount *big.Int
				StakerCount *big.Int
				
			})
			out := ret
			err := _Stakemanager.contract.Call(opts, out, "validatorState" , arg0)
			return *ret, err
		}

		// ValidatorState is a free data retrieval call binding the contract method 0x5c248855.
		//
		// Solidity: function validatorState( uint256) constant returns(amount int256, stakerCount int256)
		func (_Stakemanager *StakemanagerSession) ValidatorState( arg0 *big.Int ) (struct{ Amount *big.Int;StakerCount *big.Int; },  error) {
		  return _Stakemanager.Contract.ValidatorState(&_Stakemanager.CallOpts , arg0)
		}

		// ValidatorState is a free data retrieval call binding the contract method 0x5c248855.
		//
		// Solidity: function validatorState( uint256) constant returns(amount int256, stakerCount int256)
		func (_Stakemanager *StakemanagerCallerSession) ValidatorState( arg0 *big.Int ) (struct{ Amount *big.Int;StakerCount *big.Int; },  error) {
		  return _Stakemanager.Contract.ValidatorState(&_Stakemanager.CallOpts , arg0)
		}
	
		// ValidatorThreshold is a free data retrieval call binding the contract method 0x4fd101d7.
		//
		// Solidity: function validatorThreshold() constant returns(uint256)
		func (_Stakemanager *StakemanagerCaller) ValidatorThreshold(opts *bind.CallOpts ) (*big.Int, error) {
			var (
				ret0 = new(*big.Int)
				
			)
			out := ret0
			err := _Stakemanager.contract.Call(opts, out, "validatorThreshold" )
			return *ret0, err
		}

		// ValidatorThreshold is a free data retrieval call binding the contract method 0x4fd101d7.
		//
		// Solidity: function validatorThreshold() constant returns(uint256)
		func (_Stakemanager *StakemanagerSession) ValidatorThreshold() ( *big.Int,  error) {
		  return _Stakemanager.Contract.ValidatorThreshold(&_Stakemanager.CallOpts )
		}

		// ValidatorThreshold is a free data retrieval call binding the contract method 0x4fd101d7.
		//
		// Solidity: function validatorThreshold() constant returns(uint256)
		func (_Stakemanager *StakemanagerCallerSession) ValidatorThreshold() ( *big.Int,  error) {
		  return _Stakemanager.Contract.ValidatorThreshold(&_Stakemanager.CallOpts )
		}
	
		// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
		//
		// Solidity: function validators( uint256) constant returns(epoch uint256, amount uint256, reward uint256, activationEpoch uint256, deactivationEpoch uint256, signer address)
		func (_Stakemanager *StakemanagerCaller) Validators(opts *bind.CallOpts , arg0 *big.Int ) (struct{ Epoch *big.Int;Amount *big.Int;Reward *big.Int;ActivationEpoch *big.Int;DeactivationEpoch *big.Int;Signer common.Address; }, error) {
			ret := new(struct{
				Epoch *big.Int
				Amount *big.Int
				Reward *big.Int
				ActivationEpoch *big.Int
				DeactivationEpoch *big.Int
				Signer common.Address
				
			})
			out := ret
			err := _Stakemanager.contract.Call(opts, out, "validators" , arg0)
			return *ret, err
		}

		// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
		//
		// Solidity: function validators( uint256) constant returns(epoch uint256, amount uint256, reward uint256, activationEpoch uint256, deactivationEpoch uint256, signer address)
		func (_Stakemanager *StakemanagerSession) Validators( arg0 *big.Int ) (struct{ Epoch *big.Int;Amount *big.Int;Reward *big.Int;ActivationEpoch *big.Int;DeactivationEpoch *big.Int;Signer common.Address; },  error) {
		  return _Stakemanager.Contract.Validators(&_Stakemanager.CallOpts , arg0)
		}

		// Validators is a free data retrieval call binding the contract method 0x35aa2e44.
		//
		// Solidity: function validators( uint256) constant returns(epoch uint256, amount uint256, reward uint256, activationEpoch uint256, deactivationEpoch uint256, signer address)
		func (_Stakemanager *StakemanagerCallerSession) Validators( arg0 *big.Int ) (struct{ Epoch *big.Int;Amount *big.Int;Reward *big.Int;ActivationEpoch *big.Int;DeactivationEpoch *big.Int;Signer common.Address; },  error) {
		  return _Stakemanager.Contract.Validators(&_Stakemanager.CallOpts , arg0)
		}
	

	
		// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
		//
		// Solidity: function approve(to address, tokenId uint256) returns()
		func (_Stakemanager *StakemanagerTransactor) Approve(opts *bind.TransactOpts , to common.Address , tokenId *big.Int ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "approve" , to, tokenId)
		}

		// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
		//
		// Solidity: function approve(to address, tokenId uint256) returns()
		func (_Stakemanager *StakemanagerSession) Approve( to common.Address , tokenId *big.Int ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.Approve(&_Stakemanager.TransactOpts , to, tokenId)
		}

		// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
		//
		// Solidity: function approve(to address, tokenId uint256) returns()
		func (_Stakemanager *StakemanagerTransactorSession) Approve( to common.Address , tokenId *big.Int ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.Approve(&_Stakemanager.TransactOpts , to, tokenId)
		}
	
		// ChangeRootChain is a paid mutator transaction binding the contract method 0xe8afa8e8.
		//
		// Solidity: function changeRootChain(newRootChain address) returns()
		func (_Stakemanager *StakemanagerTransactor) ChangeRootChain(opts *bind.TransactOpts , newRootChain common.Address ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "changeRootChain" , newRootChain)
		}

		// ChangeRootChain is a paid mutator transaction binding the contract method 0xe8afa8e8.
		//
		// Solidity: function changeRootChain(newRootChain address) returns()
		func (_Stakemanager *StakemanagerSession) ChangeRootChain( newRootChain common.Address ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.ChangeRootChain(&_Stakemanager.TransactOpts , newRootChain)
		}

		// ChangeRootChain is a paid mutator transaction binding the contract method 0xe8afa8e8.
		//
		// Solidity: function changeRootChain(newRootChain address) returns()
		func (_Stakemanager *StakemanagerTransactorSession) ChangeRootChain( newRootChain common.Address ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.ChangeRootChain(&_Stakemanager.TransactOpts , newRootChain)
		}
	
		// FinalizeCommit is a paid mutator transaction binding the contract method 0x35dda498.
		//
		// Solidity: function finalizeCommit() returns()
		func (_Stakemanager *StakemanagerTransactor) FinalizeCommit(opts *bind.TransactOpts ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "finalizeCommit" )
		}

		// FinalizeCommit is a paid mutator transaction binding the contract method 0x35dda498.
		//
		// Solidity: function finalizeCommit() returns()
		func (_Stakemanager *StakemanagerSession) FinalizeCommit() (*types.Transaction, error) {
		  return _Stakemanager.Contract.FinalizeCommit(&_Stakemanager.TransactOpts )
		}

		// FinalizeCommit is a paid mutator transaction binding the contract method 0x35dda498.
		//
		// Solidity: function finalizeCommit() returns()
		func (_Stakemanager *StakemanagerTransactorSession) FinalizeCommit() (*types.Transaction, error) {
		  return _Stakemanager.Contract.FinalizeCommit(&_Stakemanager.TransactOpts )
		}
	
		// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
		//
		// Solidity: function lock() returns()
		func (_Stakemanager *StakemanagerTransactor) Lock(opts *bind.TransactOpts ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "lock" )
		}

		// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
		//
		// Solidity: function lock() returns()
		func (_Stakemanager *StakemanagerSession) Lock() (*types.Transaction, error) {
		  return _Stakemanager.Contract.Lock(&_Stakemanager.TransactOpts )
		}

		// Lock is a paid mutator transaction binding the contract method 0xf83d08ba.
		//
		// Solidity: function lock() returns()
		func (_Stakemanager *StakemanagerTransactorSession) Lock() (*types.Transaction, error) {
		  return _Stakemanager.Contract.Lock(&_Stakemanager.TransactOpts )
		}
	
		// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
		//
		// Solidity: function renounceOwnership() returns()
		func (_Stakemanager *StakemanagerTransactor) RenounceOwnership(opts *bind.TransactOpts ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "renounceOwnership" )
		}

		// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
		//
		// Solidity: function renounceOwnership() returns()
		func (_Stakemanager *StakemanagerSession) RenounceOwnership() (*types.Transaction, error) {
		  return _Stakemanager.Contract.RenounceOwnership(&_Stakemanager.TransactOpts )
		}

		// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
		//
		// Solidity: function renounceOwnership() returns()
		func (_Stakemanager *StakemanagerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
		  return _Stakemanager.Contract.RenounceOwnership(&_Stakemanager.TransactOpts )
		}
	
		// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
		//
		// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
		func (_Stakemanager *StakemanagerTransactor) SafeTransferFrom(opts *bind.TransactOpts , from common.Address , to common.Address , tokenId *big.Int , _data []byte ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "safeTransferFrom" , from, to, tokenId, _data)
		}

		// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
		//
		// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
		func (_Stakemanager *StakemanagerSession) SafeTransferFrom( from common.Address , to common.Address , tokenId *big.Int , _data []byte ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.SafeTransferFrom(&_Stakemanager.TransactOpts , from, to, tokenId, _data)
		}

		// SafeTransferFrom is a paid mutator transaction binding the contract method 0xb88d4fde.
		//
		// Solidity: function safeTransferFrom(from address, to address, tokenId uint256, _data bytes) returns()
		func (_Stakemanager *StakemanagerTransactorSession) SafeTransferFrom( from common.Address , to common.Address , tokenId *big.Int , _data []byte ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.SafeTransferFrom(&_Stakemanager.TransactOpts , from, to, tokenId, _data)
		}
	
		// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
		//
		// Solidity: function setApprovalForAll(to address, approved bool) returns()
		func (_Stakemanager *StakemanagerTransactor) SetApprovalForAll(opts *bind.TransactOpts , to common.Address , approved bool ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "setApprovalForAll" , to, approved)
		}

		// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
		//
		// Solidity: function setApprovalForAll(to address, approved bool) returns()
		func (_Stakemanager *StakemanagerSession) SetApprovalForAll( to common.Address , approved bool ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.SetApprovalForAll(&_Stakemanager.TransactOpts , to, approved)
		}

		// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
		//
		// Solidity: function setApprovalForAll(to address, approved bool) returns()
		func (_Stakemanager *StakemanagerTransactorSession) SetApprovalForAll( to common.Address , approved bool ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.SetApprovalForAll(&_Stakemanager.TransactOpts , to, approved)
		}
	
		// Stake is a paid mutator transaction binding the contract method 0x7acb7757.
		//
		// Solidity: function stake(amount uint256, signer address) returns()
		func (_Stakemanager *StakemanagerTransactor) Stake(opts *bind.TransactOpts , amount *big.Int , signer common.Address ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "stake" , amount, signer)
		}

		// Stake is a paid mutator transaction binding the contract method 0x7acb7757.
		//
		// Solidity: function stake(amount uint256, signer address) returns()
		func (_Stakemanager *StakemanagerSession) Stake( amount *big.Int , signer common.Address ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.Stake(&_Stakemanager.TransactOpts , amount, signer)
		}

		// Stake is a paid mutator transaction binding the contract method 0x7acb7757.
		//
		// Solidity: function stake(amount uint256, signer address) returns()
		func (_Stakemanager *StakemanagerTransactorSession) Stake( amount *big.Int , signer common.Address ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.Stake(&_Stakemanager.TransactOpts , amount, signer)
		}
	
		// StakeFor is a paid mutator transaction binding the contract method 0x9b8f04b7.
		//
		// Solidity: function stakeFor(user address, amount uint256, signer address) returns()
		func (_Stakemanager *StakemanagerTransactor) StakeFor(opts *bind.TransactOpts , user common.Address , amount *big.Int , signer common.Address ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "stakeFor" , user, amount, signer)
		}

		// StakeFor is a paid mutator transaction binding the contract method 0x9b8f04b7.
		//
		// Solidity: function stakeFor(user address, amount uint256, signer address) returns()
		func (_Stakemanager *StakemanagerSession) StakeFor( user common.Address , amount *big.Int , signer common.Address ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.StakeFor(&_Stakemanager.TransactOpts , user, amount, signer)
		}

		// StakeFor is a paid mutator transaction binding the contract method 0x9b8f04b7.
		//
		// Solidity: function stakeFor(user address, amount uint256, signer address) returns()
		func (_Stakemanager *StakemanagerTransactorSession) StakeFor( user common.Address , amount *big.Int , signer common.Address ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.StakeFor(&_Stakemanager.TransactOpts , user, amount, signer)
		}
	
		// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
		//
		// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
		func (_Stakemanager *StakemanagerTransactor) TransferFrom(opts *bind.TransactOpts , from common.Address , to common.Address , tokenId *big.Int ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "transferFrom" , from, to, tokenId)
		}

		// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
		//
		// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
		func (_Stakemanager *StakemanagerSession) TransferFrom( from common.Address , to common.Address , tokenId *big.Int ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.TransferFrom(&_Stakemanager.TransactOpts , from, to, tokenId)
		}

		// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
		//
		// Solidity: function transferFrom(from address, to address, tokenId uint256) returns()
		func (_Stakemanager *StakemanagerTransactorSession) TransferFrom( from common.Address , to common.Address , tokenId *big.Int ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.TransferFrom(&_Stakemanager.TransactOpts , from, to, tokenId)
		}
	
		// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
		//
		// Solidity: function transferOwnership(newOwner address) returns()
		func (_Stakemanager *StakemanagerTransactor) TransferOwnership(opts *bind.TransactOpts , newOwner common.Address ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "transferOwnership" , newOwner)
		}

		// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
		//
		// Solidity: function transferOwnership(newOwner address) returns()
		func (_Stakemanager *StakemanagerSession) TransferOwnership( newOwner common.Address ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.TransferOwnership(&_Stakemanager.TransactOpts , newOwner)
		}

		// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
		//
		// Solidity: function transferOwnership(newOwner address) returns()
		func (_Stakemanager *StakemanagerTransactorSession) TransferOwnership( newOwner common.Address ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.TransferOwnership(&_Stakemanager.TransactOpts , newOwner)
		}
	
		// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
		//
		// Solidity: function unlock() returns()
		func (_Stakemanager *StakemanagerTransactor) Unlock(opts *bind.TransactOpts ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "unlock" )
		}

		// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
		//
		// Solidity: function unlock() returns()
		func (_Stakemanager *StakemanagerSession) Unlock() (*types.Transaction, error) {
		  return _Stakemanager.Contract.Unlock(&_Stakemanager.TransactOpts )
		}

		// Unlock is a paid mutator transaction binding the contract method 0xa69df4b5.
		//
		// Solidity: function unlock() returns()
		func (_Stakemanager *StakemanagerTransactorSession) Unlock() (*types.Transaction, error) {
		  return _Stakemanager.Contract.Unlock(&_Stakemanager.TransactOpts )
		}
	
		// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
		//
		// Solidity: function unstake(validatorId uint256) returns()
		func (_Stakemanager *StakemanagerTransactor) Unstake(opts *bind.TransactOpts , validatorId *big.Int ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "unstake" , validatorId)
		}

		// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
		//
		// Solidity: function unstake(validatorId uint256) returns()
		func (_Stakemanager *StakemanagerSession) Unstake( validatorId *big.Int ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.Unstake(&_Stakemanager.TransactOpts , validatorId)
		}

		// Unstake is a paid mutator transaction binding the contract method 0x2e17de78.
		//
		// Solidity: function unstake(validatorId uint256) returns()
		func (_Stakemanager *StakemanagerTransactorSession) Unstake( validatorId *big.Int ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.Unstake(&_Stakemanager.TransactOpts , validatorId)
		}
	
		// UnstakeClaim is a paid mutator transaction binding the contract method 0xd86d53e7.
		//
		// Solidity: function unstakeClaim(validatorId uint256) returns()
		func (_Stakemanager *StakemanagerTransactor) UnstakeClaim(opts *bind.TransactOpts , validatorId *big.Int ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "unstakeClaim" , validatorId)
		}

		// UnstakeClaim is a paid mutator transaction binding the contract method 0xd86d53e7.
		//
		// Solidity: function unstakeClaim(validatorId uint256) returns()
		func (_Stakemanager *StakemanagerSession) UnstakeClaim( validatorId *big.Int ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.UnstakeClaim(&_Stakemanager.TransactOpts , validatorId)
		}

		// UnstakeClaim is a paid mutator transaction binding the contract method 0xd86d53e7.
		//
		// Solidity: function unstakeClaim(validatorId uint256) returns()
		func (_Stakemanager *StakemanagerTransactorSession) UnstakeClaim( validatorId *big.Int ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.UnstakeClaim(&_Stakemanager.TransactOpts , validatorId)
		}
	
		// UpdateDynastyValue is a paid mutator transaction binding the contract method 0xe6692f49.
		//
		// Solidity: function updateDynastyValue(newDynasty uint256) returns()
		func (_Stakemanager *StakemanagerTransactor) UpdateDynastyValue(opts *bind.TransactOpts , newDynasty *big.Int ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "updateDynastyValue" , newDynasty)
		}

		// UpdateDynastyValue is a paid mutator transaction binding the contract method 0xe6692f49.
		//
		// Solidity: function updateDynastyValue(newDynasty uint256) returns()
		func (_Stakemanager *StakemanagerSession) UpdateDynastyValue( newDynasty *big.Int ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.UpdateDynastyValue(&_Stakemanager.TransactOpts , newDynasty)
		}

		// UpdateDynastyValue is a paid mutator transaction binding the contract method 0xe6692f49.
		//
		// Solidity: function updateDynastyValue(newDynasty uint256) returns()
		func (_Stakemanager *StakemanagerTransactorSession) UpdateDynastyValue( newDynasty *big.Int ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.UpdateDynastyValue(&_Stakemanager.TransactOpts , newDynasty)
		}
	
		// UpdateMinLockInPeriod is a paid mutator transaction binding the contract method 0x98ee773b.
		//
		// Solidity: function updateMinLockInPeriod(epochs uint256) returns()
		func (_Stakemanager *StakemanagerTransactor) UpdateMinLockInPeriod(opts *bind.TransactOpts , epochs *big.Int ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "updateMinLockInPeriod" , epochs)
		}

		// UpdateMinLockInPeriod is a paid mutator transaction binding the contract method 0x98ee773b.
		//
		// Solidity: function updateMinLockInPeriod(epochs uint256) returns()
		func (_Stakemanager *StakemanagerSession) UpdateMinLockInPeriod( epochs *big.Int ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.UpdateMinLockInPeriod(&_Stakemanager.TransactOpts , epochs)
		}

		// UpdateMinLockInPeriod is a paid mutator transaction binding the contract method 0x98ee773b.
		//
		// Solidity: function updateMinLockInPeriod(epochs uint256) returns()
		func (_Stakemanager *StakemanagerTransactorSession) UpdateMinLockInPeriod( epochs *big.Int ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.UpdateMinLockInPeriod(&_Stakemanager.TransactOpts , epochs)
		}
	
		// UpdateSigner is a paid mutator transaction binding the contract method 0x8f283a86.
		//
		// Solidity: function updateSigner(validatorId uint256, _signer address) returns()
		func (_Stakemanager *StakemanagerTransactor) UpdateSigner(opts *bind.TransactOpts , validatorId *big.Int , _signer common.Address ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "updateSigner" , validatorId, _signer)
		}

		// UpdateSigner is a paid mutator transaction binding the contract method 0x8f283a86.
		//
		// Solidity: function updateSigner(validatorId uint256, _signer address) returns()
		func (_Stakemanager *StakemanagerSession) UpdateSigner( validatorId *big.Int , _signer common.Address ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.UpdateSigner(&_Stakemanager.TransactOpts , validatorId, _signer)
		}

		// UpdateSigner is a paid mutator transaction binding the contract method 0x8f283a86.
		//
		// Solidity: function updateSigner(validatorId uint256, _signer address) returns()
		func (_Stakemanager *StakemanagerTransactorSession) UpdateSigner( validatorId *big.Int , _signer common.Address ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.UpdateSigner(&_Stakemanager.TransactOpts , validatorId, _signer)
		}
	
		// UpdateValidatorThreshold is a paid mutator transaction binding the contract method 0x16827b1b.
		//
		// Solidity: function updateValidatorThreshold(newThreshold uint256) returns()
		func (_Stakemanager *StakemanagerTransactor) UpdateValidatorThreshold(opts *bind.TransactOpts , newThreshold *big.Int ) (*types.Transaction, error) {
			return _Stakemanager.contract.Transact(opts, "updateValidatorThreshold" , newThreshold)
		}

		// UpdateValidatorThreshold is a paid mutator transaction binding the contract method 0x16827b1b.
		//
		// Solidity: function updateValidatorThreshold(newThreshold uint256) returns()
		func (_Stakemanager *StakemanagerSession) UpdateValidatorThreshold( newThreshold *big.Int ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.UpdateValidatorThreshold(&_Stakemanager.TransactOpts , newThreshold)
		}

		// UpdateValidatorThreshold is a paid mutator transaction binding the contract method 0x16827b1b.
		//
		// Solidity: function updateValidatorThreshold(newThreshold uint256) returns()
		func (_Stakemanager *StakemanagerTransactorSession) UpdateValidatorThreshold( newThreshold *big.Int ) (*types.Transaction, error) {
		  return _Stakemanager.Contract.UpdateValidatorThreshold(&_Stakemanager.TransactOpts , newThreshold)
		}
	

	
		// StakemanagerApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the Stakemanager contract.
		type StakemanagerApprovalIterator struct {
			Event *StakemanagerApproval // Event containing the contract specifics and raw log

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
		func (it *StakemanagerApprovalIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(StakemanagerApproval)
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
				it.Event = new(StakemanagerApproval)
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
		func (it *StakemanagerApprovalIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *StakemanagerApprovalIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// StakemanagerApproval represents a Approval event raised by the Stakemanager contract.
		type StakemanagerApproval struct { 
			Owner common.Address; 
			Approved common.Address; 
			TokenId *big.Int; 
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
		//
		// Solidity: e Approval(owner indexed address, approved indexed address, tokenId indexed uint256)
 		func (_Stakemanager *StakemanagerFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*StakemanagerApprovalIterator, error) {
			
			var ownerRule []interface{}
			for _, ownerItem := range owner {
				ownerRule = append(ownerRule, ownerItem)
			}
			var approvedRule []interface{}
			for _, approvedItem := range approved {
				approvedRule = append(approvedRule, approvedItem)
			}
			var tokenIdRule []interface{}
			for _, tokenIdItem := range tokenId {
				tokenIdRule = append(tokenIdRule, tokenIdItem)
			}

			logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
			if err != nil {
				return nil, err
			}
			return &StakemanagerApprovalIterator{contract: _Stakemanager.contract, event: "Approval", logs: logs, sub: sub}, nil
 		}

		// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
		//
		// Solidity: e Approval(owner indexed address, approved indexed address, tokenId indexed uint256)
		func (_Stakemanager *StakemanagerFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *StakemanagerApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {
			
			var ownerRule []interface{}
			for _, ownerItem := range owner {
				ownerRule = append(ownerRule, ownerItem)
			}
			var approvedRule []interface{}
			for _, approvedItem := range approved {
				approvedRule = append(approvedRule, approvedItem)
			}
			var tokenIdRule []interface{}
			for _, tokenIdItem := range tokenId {
				tokenIdRule = append(tokenIdRule, tokenIdItem)
			}

			logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(StakemanagerApproval)
						if err := _Stakemanager.contract.UnpackLog(event, "Approval", log); err != nil {
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
 	
		// StakemanagerApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the Stakemanager contract.
		type StakemanagerApprovalForAllIterator struct {
			Event *StakemanagerApprovalForAll // Event containing the contract specifics and raw log

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
		func (it *StakemanagerApprovalForAllIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(StakemanagerApprovalForAll)
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
				it.Event = new(StakemanagerApprovalForAll)
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
		func (it *StakemanagerApprovalForAllIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *StakemanagerApprovalForAllIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// StakemanagerApprovalForAll represents a ApprovalForAll event raised by the Stakemanager contract.
		type StakemanagerApprovalForAll struct { 
			Owner common.Address; 
			Operator common.Address; 
			Approved bool; 
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
		//
		// Solidity: e ApprovalForAll(owner indexed address, operator indexed address, approved bool)
 		func (_Stakemanager *StakemanagerFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*StakemanagerApprovalForAllIterator, error) {
			
			var ownerRule []interface{}
			for _, ownerItem := range owner {
				ownerRule = append(ownerRule, ownerItem)
			}
			var operatorRule []interface{}
			for _, operatorItem := range operator {
				operatorRule = append(operatorRule, operatorItem)
			}
			

			logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
			if err != nil {
				return nil, err
			}
			return &StakemanagerApprovalForAllIterator{contract: _Stakemanager.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
 		}

		// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
		//
		// Solidity: e ApprovalForAll(owner indexed address, operator indexed address, approved bool)
		func (_Stakemanager *StakemanagerFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *StakemanagerApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {
			
			var ownerRule []interface{}
			for _, ownerItem := range owner {
				ownerRule = append(ownerRule, ownerItem)
			}
			var operatorRule []interface{}
			for _, operatorItem := range operator {
				operatorRule = append(operatorRule, operatorItem)
			}
			

			logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(StakemanagerApprovalForAll)
						if err := _Stakemanager.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
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
 	
		// StakemanagerDynastyValueChangeIterator is returned from FilterDynastyValueChange and is used to iterate over the raw logs and unpacked data for DynastyValueChange events raised by the Stakemanager contract.
		type StakemanagerDynastyValueChangeIterator struct {
			Event *StakemanagerDynastyValueChange // Event containing the contract specifics and raw log

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
		func (it *StakemanagerDynastyValueChangeIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(StakemanagerDynastyValueChange)
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
				it.Event = new(StakemanagerDynastyValueChange)
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
		func (it *StakemanagerDynastyValueChangeIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *StakemanagerDynastyValueChangeIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// StakemanagerDynastyValueChange represents a DynastyValueChange event raised by the Stakemanager contract.
		type StakemanagerDynastyValueChange struct { 
			NewDynasty *big.Int; 
			OldDynasty *big.Int; 
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterDynastyValueChange is a free log retrieval operation binding the contract event 0x9444bfcfa6aed72a15da73de1220dcc07d7864119c44abfec0037bbcacefda98.
		//
		// Solidity: e DynastyValueChange(newDynasty uint256, oldDynasty uint256)
 		func (_Stakemanager *StakemanagerFilterer) FilterDynastyValueChange(opts *bind.FilterOpts) (*StakemanagerDynastyValueChangeIterator, error) {
			
			
			

			logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "DynastyValueChange")
			if err != nil {
				return nil, err
			}
			return &StakemanagerDynastyValueChangeIterator{contract: _Stakemanager.contract, event: "DynastyValueChange", logs: logs, sub: sub}, nil
 		}

		// WatchDynastyValueChange is a free log subscription operation binding the contract event 0x9444bfcfa6aed72a15da73de1220dcc07d7864119c44abfec0037bbcacefda98.
		//
		// Solidity: e DynastyValueChange(newDynasty uint256, oldDynasty uint256)
		func (_Stakemanager *StakemanagerFilterer) WatchDynastyValueChange(opts *bind.WatchOpts, sink chan<- *StakemanagerDynastyValueChange) (event.Subscription, error) {
			
			
			

			logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "DynastyValueChange")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(StakemanagerDynastyValueChange)
						if err := _Stakemanager.contract.UnpackLog(event, "DynastyValueChange", log); err != nil {
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
 	
		// StakemanagerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the Stakemanager contract.
		type StakemanagerOwnershipTransferredIterator struct {
			Event *StakemanagerOwnershipTransferred // Event containing the contract specifics and raw log

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
		func (it *StakemanagerOwnershipTransferredIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(StakemanagerOwnershipTransferred)
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
				it.Event = new(StakemanagerOwnershipTransferred)
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
		func (it *StakemanagerOwnershipTransferredIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *StakemanagerOwnershipTransferredIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// StakemanagerOwnershipTransferred represents a OwnershipTransferred event raised by the Stakemanager contract.
		type StakemanagerOwnershipTransferred struct { 
			PreviousOwner common.Address; 
			NewOwner common.Address; 
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
		//
		// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
 		func (_Stakemanager *StakemanagerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*StakemanagerOwnershipTransferredIterator, error) {
			
			var previousOwnerRule []interface{}
			for _, previousOwnerItem := range previousOwner {
				previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
			}
			var newOwnerRule []interface{}
			for _, newOwnerItem := range newOwner {
				newOwnerRule = append(newOwnerRule, newOwnerItem)
			}

			logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
			if err != nil {
				return nil, err
			}
			return &StakemanagerOwnershipTransferredIterator{contract: _Stakemanager.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
 		}

		// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
		//
		// Solidity: e OwnershipTransferred(previousOwner indexed address, newOwner indexed address)
		func (_Stakemanager *StakemanagerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *StakemanagerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {
			
			var previousOwnerRule []interface{}
			for _, previousOwnerItem := range previousOwner {
				previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
			}
			var newOwnerRule []interface{}
			for _, newOwnerItem := range newOwner {
				newOwnerRule = append(newOwnerRule, newOwnerItem)
			}

			logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(StakemanagerOwnershipTransferred)
						if err := _Stakemanager.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
 	
		// StakemanagerRootChainChangedIterator is returned from FilterRootChainChanged and is used to iterate over the raw logs and unpacked data for RootChainChanged events raised by the Stakemanager contract.
		type StakemanagerRootChainChangedIterator struct {
			Event *StakemanagerRootChainChanged // Event containing the contract specifics and raw log

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
		func (it *StakemanagerRootChainChangedIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(StakemanagerRootChainChanged)
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
				it.Event = new(StakemanagerRootChainChanged)
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
		func (it *StakemanagerRootChainChangedIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *StakemanagerRootChainChangedIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// StakemanagerRootChainChanged represents a RootChainChanged event raised by the Stakemanager contract.
		type StakemanagerRootChainChanged struct { 
			PreviousRootChain common.Address; 
			NewRootChain common.Address; 
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterRootChainChanged is a free log retrieval operation binding the contract event 0x211c9015fc81c0dbd45bd99f0f29fc1c143bfd53442d5ffd722bbbef7a887fe9.
		//
		// Solidity: e RootChainChanged(previousRootChain indexed address, newRootChain indexed address)
 		func (_Stakemanager *StakemanagerFilterer) FilterRootChainChanged(opts *bind.FilterOpts, previousRootChain []common.Address, newRootChain []common.Address) (*StakemanagerRootChainChangedIterator, error) {
			
			var previousRootChainRule []interface{}
			for _, previousRootChainItem := range previousRootChain {
				previousRootChainRule = append(previousRootChainRule, previousRootChainItem)
			}
			var newRootChainRule []interface{}
			for _, newRootChainItem := range newRootChain {
				newRootChainRule = append(newRootChainRule, newRootChainItem)
			}

			logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "RootChainChanged", previousRootChainRule, newRootChainRule)
			if err != nil {
				return nil, err
			}
			return &StakemanagerRootChainChangedIterator{contract: _Stakemanager.contract, event: "RootChainChanged", logs: logs, sub: sub}, nil
 		}

		// WatchRootChainChanged is a free log subscription operation binding the contract event 0x211c9015fc81c0dbd45bd99f0f29fc1c143bfd53442d5ffd722bbbef7a887fe9.
		//
		// Solidity: e RootChainChanged(previousRootChain indexed address, newRootChain indexed address)
		func (_Stakemanager *StakemanagerFilterer) WatchRootChainChanged(opts *bind.WatchOpts, sink chan<- *StakemanagerRootChainChanged, previousRootChain []common.Address, newRootChain []common.Address) (event.Subscription, error) {
			
			var previousRootChainRule []interface{}
			for _, previousRootChainItem := range previousRootChain {
				previousRootChainRule = append(previousRootChainRule, previousRootChainItem)
			}
			var newRootChainRule []interface{}
			for _, newRootChainItem := range newRootChain {
				newRootChainRule = append(newRootChainRule, newRootChainItem)
			}

			logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "RootChainChanged", previousRootChainRule, newRootChainRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(StakemanagerRootChainChanged)
						if err := _Stakemanager.contract.UnpackLog(event, "RootChainChanged", log); err != nil {
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
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
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
			ValidatorId *big.Int; 
			NewSigner common.Address; 
			OldSigner common.Address; 
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterSignerChange is a free log retrieval operation binding the contract event 0x7dfd3bad1e3cac97d3b89ff06d78394523c4f08fdee4daa71a59160003240c89.
		//
		// Solidity: e SignerChange(validatorId indexed uint256, newSigner indexed address, oldSigner indexed address)
 		func (_Stakemanager *StakemanagerFilterer) FilterSignerChange(opts *bind.FilterOpts, validatorId []*big.Int, newSigner []common.Address, oldSigner []common.Address) (*StakemanagerSignerChangeIterator, error) {
			
			var validatorIdRule []interface{}
			for _, validatorIdItem := range validatorId {
				validatorIdRule = append(validatorIdRule, validatorIdItem)
			}
			var newSignerRule []interface{}
			for _, newSignerItem := range newSigner {
				newSignerRule = append(newSignerRule, newSignerItem)
			}
			var oldSignerRule []interface{}
			for _, oldSignerItem := range oldSigner {
				oldSignerRule = append(oldSignerRule, oldSignerItem)
			}

			logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "SignerChange", validatorIdRule, newSignerRule, oldSignerRule)
			if err != nil {
				return nil, err
			}
			return &StakemanagerSignerChangeIterator{contract: _Stakemanager.contract, event: "SignerChange", logs: logs, sub: sub}, nil
 		}

		// WatchSignerChange is a free log subscription operation binding the contract event 0x7dfd3bad1e3cac97d3b89ff06d78394523c4f08fdee4daa71a59160003240c89.
		//
		// Solidity: e SignerChange(validatorId indexed uint256, newSigner indexed address, oldSigner indexed address)
		func (_Stakemanager *StakemanagerFilterer) WatchSignerChange(opts *bind.WatchOpts, sink chan<- *StakemanagerSignerChange, validatorId []*big.Int, newSigner []common.Address, oldSigner []common.Address) (event.Subscription, error) {
			
			var validatorIdRule []interface{}
			for _, validatorIdItem := range validatorId {
				validatorIdRule = append(validatorIdRule, validatorIdItem)
			}
			var newSignerRule []interface{}
			for _, newSignerItem := range newSigner {
				newSignerRule = append(newSignerRule, newSignerItem)
			}
			var oldSignerRule []interface{}
			for _, oldSignerItem := range oldSigner {
				oldSignerRule = append(oldSignerRule, oldSignerItem)
			}

			logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "SignerChange", validatorIdRule, newSignerRule, oldSignerRule)
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
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
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
			User common.Address; 
			ValidatorId *big.Int; 
			ActivatonEpoch *big.Int; 
			Amount *big.Int; 
			Total *big.Int; 
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterStaked is a free log retrieval operation binding the contract event 0x9cfd25589d1eb8ad71e342a86a8524e83522e3936c0803048c08f6d9ad974f40.
		//
		// Solidity: e Staked(user indexed address, validatorId indexed uint256, activatonEpoch indexed uint256, amount uint256, total uint256)
 		func (_Stakemanager *StakemanagerFilterer) FilterStaked(opts *bind.FilterOpts, user []common.Address, validatorId []*big.Int, activatonEpoch []*big.Int) (*StakemanagerStakedIterator, error) {
			
			var userRule []interface{}
			for _, userItem := range user {
				userRule = append(userRule, userItem)
			}
			var validatorIdRule []interface{}
			for _, validatorIdItem := range validatorId {
				validatorIdRule = append(validatorIdRule, validatorIdItem)
			}
			var activatonEpochRule []interface{}
			for _, activatonEpochItem := range activatonEpoch {
				activatonEpochRule = append(activatonEpochRule, activatonEpochItem)
			}
			
			

			logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "Staked", userRule, validatorIdRule, activatonEpochRule)
			if err != nil {
				return nil, err
			}
			return &StakemanagerStakedIterator{contract: _Stakemanager.contract, event: "Staked", logs: logs, sub: sub}, nil
 		}

		// WatchStaked is a free log subscription operation binding the contract event 0x9cfd25589d1eb8ad71e342a86a8524e83522e3936c0803048c08f6d9ad974f40.
		//
		// Solidity: e Staked(user indexed address, validatorId indexed uint256, activatonEpoch indexed uint256, amount uint256, total uint256)
		func (_Stakemanager *StakemanagerFilterer) WatchStaked(opts *bind.WatchOpts, sink chan<- *StakemanagerStaked, user []common.Address, validatorId []*big.Int, activatonEpoch []*big.Int) (event.Subscription, error) {
			
			var userRule []interface{}
			for _, userItem := range user {
				userRule = append(userRule, userItem)
			}
			var validatorIdRule []interface{}
			for _, validatorIdItem := range validatorId {
				validatorIdRule = append(validatorIdRule, validatorIdItem)
			}
			var activatonEpochRule []interface{}
			for _, activatonEpochItem := range activatonEpoch {
				activatonEpochRule = append(activatonEpochRule, activatonEpochItem)
			}
			
			

			logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "Staked", userRule, validatorIdRule, activatonEpochRule)
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
 	
		// StakemanagerThresholdChangeIterator is returned from FilterThresholdChange and is used to iterate over the raw logs and unpacked data for ThresholdChange events raised by the Stakemanager contract.
		type StakemanagerThresholdChangeIterator struct {
			Event *StakemanagerThresholdChange // Event containing the contract specifics and raw log

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
		func (it *StakemanagerThresholdChangeIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(StakemanagerThresholdChange)
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
				it.Event = new(StakemanagerThresholdChange)
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
		func (it *StakemanagerThresholdChangeIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *StakemanagerThresholdChangeIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// StakemanagerThresholdChange represents a ThresholdChange event raised by the Stakemanager contract.
		type StakemanagerThresholdChange struct { 
			NewThreshold *big.Int; 
			OldThreshold *big.Int; 
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterThresholdChange is a free log retrieval operation binding the contract event 0x5d16a900896e1160c2033bc940e6b072d3dc3b6a996fefb9b3b9b9678841824c.
		//
		// Solidity: e ThresholdChange(newThreshold uint256, oldThreshold uint256)
 		func (_Stakemanager *StakemanagerFilterer) FilterThresholdChange(opts *bind.FilterOpts) (*StakemanagerThresholdChangeIterator, error) {
			
			
			

			logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "ThresholdChange")
			if err != nil {
				return nil, err
			}
			return &StakemanagerThresholdChangeIterator{contract: _Stakemanager.contract, event: "ThresholdChange", logs: logs, sub: sub}, nil
 		}

		// WatchThresholdChange is a free log subscription operation binding the contract event 0x5d16a900896e1160c2033bc940e6b072d3dc3b6a996fefb9b3b9b9678841824c.
		//
		// Solidity: e ThresholdChange(newThreshold uint256, oldThreshold uint256)
		func (_Stakemanager *StakemanagerFilterer) WatchThresholdChange(opts *bind.WatchOpts, sink chan<- *StakemanagerThresholdChange) (event.Subscription, error) {
			
			
			

			logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "ThresholdChange")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(StakemanagerThresholdChange)
						if err := _Stakemanager.contract.UnpackLog(event, "ThresholdChange", log); err != nil {
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
 	
		// StakemanagerTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the Stakemanager contract.
		type StakemanagerTransferIterator struct {
			Event *StakemanagerTransfer // Event containing the contract specifics and raw log

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
		func (it *StakemanagerTransferIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(StakemanagerTransfer)
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
				it.Event = new(StakemanagerTransfer)
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
		func (it *StakemanagerTransferIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *StakemanagerTransferIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// StakemanagerTransfer represents a Transfer event raised by the Stakemanager contract.
		type StakemanagerTransfer struct { 
			From common.Address; 
			To common.Address; 
			TokenId *big.Int; 
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
		//
		// Solidity: e Transfer(from indexed address, to indexed address, tokenId indexed uint256)
 		func (_Stakemanager *StakemanagerFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*StakemanagerTransferIterator, error) {
			
			var fromRule []interface{}
			for _, fromItem := range from {
				fromRule = append(fromRule, fromItem)
			}
			var toRule []interface{}
			for _, toItem := range to {
				toRule = append(toRule, toItem)
			}
			var tokenIdRule []interface{}
			for _, tokenIdItem := range tokenId {
				tokenIdRule = append(tokenIdRule, tokenIdItem)
			}

			logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
			if err != nil {
				return nil, err
			}
			return &StakemanagerTransferIterator{contract: _Stakemanager.contract, event: "Transfer", logs: logs, sub: sub}, nil
 		}

		// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
		//
		// Solidity: e Transfer(from indexed address, to indexed address, tokenId indexed uint256)
		func (_Stakemanager *StakemanagerFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *StakemanagerTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {
			
			var fromRule []interface{}
			for _, fromItem := range from {
				fromRule = append(fromRule, fromItem)
			}
			var toRule []interface{}
			for _, toItem := range to {
				toRule = append(toRule, toItem)
			}
			var tokenIdRule []interface{}
			for _, tokenIdItem := range tokenId {
				tokenIdRule = append(tokenIdRule, tokenIdItem)
			}

			logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(StakemanagerTransfer)
						if err := _Stakemanager.contract.UnpackLog(event, "Transfer", log); err != nil {
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
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
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
			ValidatorId *big.Int; 
			User common.Address; 
			Amount *big.Int; 
			DeactivationEpoch *big.Int; 
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterUnstakeInit is a free log retrieval operation binding the contract event 0xa279613f0232ccbb11d60294ebc6318618f532ab3e9aeac28503595c467d0680.
		//
		// Solidity: e UnstakeInit(validatorId indexed uint256, user indexed address, amount indexed uint256, deactivationEpoch uint256)
 		func (_Stakemanager *StakemanagerFilterer) FilterUnstakeInit(opts *bind.FilterOpts, validatorId []*big.Int, user []common.Address, amount []*big.Int) (*StakemanagerUnstakeInitIterator, error) {
			
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
			

			logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "UnstakeInit", validatorIdRule, userRule, amountRule)
			if err != nil {
				return nil, err
			}
			return &StakemanagerUnstakeInitIterator{contract: _Stakemanager.contract, event: "UnstakeInit", logs: logs, sub: sub}, nil
 		}

		// WatchUnstakeInit is a free log subscription operation binding the contract event 0xa279613f0232ccbb11d60294ebc6318618f532ab3e9aeac28503595c467d0680.
		//
		// Solidity: e UnstakeInit(validatorId indexed uint256, user indexed address, amount indexed uint256, deactivationEpoch uint256)
		func (_Stakemanager *StakemanagerFilterer) WatchUnstakeInit(opts *bind.WatchOpts, sink chan<- *StakemanagerUnstakeInit, validatorId []*big.Int, user []common.Address, amount []*big.Int) (event.Subscription, error) {
			
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
			

			logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "UnstakeInit", validatorIdRule, userRule, amountRule)
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
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
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
			User common.Address; 
			ValidatorId *big.Int; 
			Amount *big.Int; 
			Total *big.Int; 
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterUnstaked is a free log retrieval operation binding the contract event 0x204fccf0d92ed8d48f204adb39b2e81e92bad0dedb93f5716ca9478cfb57de00.
		//
		// Solidity: e Unstaked(user indexed address, validatorId indexed uint256, amount uint256, total uint256)
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
		// Solidity: e Unstaked(user indexed address, validatorId indexed uint256, amount uint256, total uint256)
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
 	
		// StakemanagerSigseventIterator is returned from FilterSigsevent and is used to iterate over the raw logs and unpacked data for Sigsevent events raised by the Stakemanager contract.
		type StakemanagerSigseventIterator struct {
			Event *StakemanagerSigsevent // Event containing the contract specifics and raw log

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
		func (it *StakemanagerSigseventIterator) Next() bool {
			// If the iterator failed, stop iterating
			if (it.fail != nil) {
				return false
			}
			// If the iterator completed, deliver directly whatever's available
			if (it.done) {
				select {
				case log := <-it.logs:
					it.Event = new(StakemanagerSigsevent)
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
				it.Event = new(StakemanagerSigsevent)
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
		func (it *StakemanagerSigseventIterator) Error() error {
			return it.fail
		}
		// Close terminates the iteration process, releasing any pending underlying
		// resources.
		func (it *StakemanagerSigseventIterator) Close() error {
			it.sub.Unsubscribe()
			return nil
		}

		// StakemanagerSigsevent represents a Sigsevent event raised by the Stakemanager contract.
		type StakemanagerSigsevent struct { 
			 common.Address; 
			 []byte; 
			 *big.Int; 
			 *big.Int; 
			Raw types.Log // Blockchain specific contextual infos
		}

		// FilterSigsevent is a free log retrieval operation binding the contract event 0x19fd57f3755f6494c5c7fe4f6025f43b1af4e3f3dc6cc390efcec3b13e5c30ad.
		//
		// Solidity: e sigsevent( address,  bytes,  uint256,  uint256)
 		func (_Stakemanager *StakemanagerFilterer) FilterSigsevent(opts *bind.FilterOpts) (*StakemanagerSigseventIterator, error) {
			
			
			
			
			

			logs, sub, err := _Stakemanager.contract.FilterLogs(opts, "sigsevent")
			if err != nil {
				return nil, err
			}
			return &StakemanagerSigseventIterator{contract: _Stakemanager.contract, event: "sigsevent", logs: logs, sub: sub}, nil
 		}

		// WatchSigsevent is a free log subscription operation binding the contract event 0x19fd57f3755f6494c5c7fe4f6025f43b1af4e3f3dc6cc390efcec3b13e5c30ad.
		//
		// Solidity: e sigsevent( address,  bytes,  uint256,  uint256)
		func (_Stakemanager *StakemanagerFilterer) WatchSigsevent(opts *bind.WatchOpts, sink chan<- *StakemanagerSigsevent) (event.Subscription, error) {
			
			
			
			
			

			logs, sub, err := _Stakemanager.contract.WatchLogs(opts, "sigsevent")
			if err != nil {
				return nil, err
			}
			return event.NewSubscription(func(quit <-chan struct{}) error {
				defer sub.Unsubscribe()
				for {
					select {
					case log := <-logs:
						// New log arrived, parse the event and forward to the user
						event := new(StakemanagerSigsevent)
						if err := _Stakemanager.contract.UnpackLog(event, "sigsevent", log); err != nil {
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
 	


