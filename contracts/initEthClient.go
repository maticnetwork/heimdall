package contract

import (
	"encoding/hex"
	"github.com/basecoin/contracts/StakeManager"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmtypes "github.com/tendermint/tendermint/types"
	"log"
	"math/big"
)

var (
	stakeManagerAddress = "0x8b28d78eb59c323867c43b4ab8d06e0f1efa1573"
)

func getValidatorByIndex(_index int64) abci.Validator {
	client := initKovan()
	stakeManagerInstance, err := StakeManager.NewContracts(common.HexToAddress(stakeManagerAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	validator, _ := stakeManagerInstance.Validators(nil, big.NewInt(_index))
	var _pubkey secp256k1.PubKeySecp256k1
	_pub, _ := hex.DecodeString(validator.Pubkey)
	copy(_pubkey[:], _pub[:])
	_address, _ := hex.DecodeString(_pubkey.Address().String())

	abciValidator := abci.Validator{
		Address: _address,
		Power:   validator.Power.Int64(),
		PubKey:  tmtypes.TM2PB.PubKey(_pubkey),
	}
	return abciValidator

}
func initKovan() *ethclient.Client {
	client, err := ethclient.Dial("https://kovan.infura.io")
	if err != nil {
		log.Fatal(err)
	}
	return client
}
func initMatic() *ethclient.Client {
	client, err := ethclient.Dial("https://testnet.matic.network")
	if err != nil {
		log.Fatal(err)
	}
	return client
}
