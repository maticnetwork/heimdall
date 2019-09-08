package pier

import "math/big"

type (
	ContractCheckpoint struct {
		start              uint64
		end                uint64
		currentHeaderBlock *big.Int
		err                error
	}

	HeimdallCheckpoint struct {
		start uint64
		end   uint64
		found bool
	}
)

func NewContractCheckpoint(_start uint64, _end uint64, _currentHeaderBlock *big.Int, _err error) ContractCheckpoint {
	return ContractCheckpoint{
		start:              _start,
		end:                _end,
		currentHeaderBlock: _currentHeaderBlock,
		err:                _err,
	}
}

// Creates new heimdall checkpoint object
func NewHeimdallCheckpoint(_start uint64, _end uint64, _found bool) HeimdallCheckpoint {
	return HeimdallCheckpoint{
		start: _start,
		end:   _end,
		found: _found,
	}
}
