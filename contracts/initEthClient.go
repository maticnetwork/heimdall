package contract

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/privval"
	"github.com/xsleonard/go-merkle"

	"github.com/maticnetwork/heimdall/contracts/stakemanager"
	"github.com/maticnetwork/heimdall/contracts/validatorset"
)

var cdc = amino.NewCodec()

func init() {
	cdc.RegisterConcrete(secp256k1.PubKeySecp256k1{}, secp256k1.Secp256k1PubKeyAminoRoute, nil)
	cdc.RegisterConcrete(secp256k1.PrivKeySecp256k1{}, secp256k1.Secp256k1PrivKeyAminoRoute, nil)
}

var (
	stakeManagerAddress = "8b28d78eb59c323867c43b4ab8d06e0f1efa1573"
	rootchainAddress    = "e022d867085b1617dc9fb04b474c4de580dccf1a"
	validatorSetAddress = "295c050e82a39392799d31ecb3f7db921daa136c"
)
var kovanClient, _ = ethclient.Dial("https://kovan.infura.io")
var maticClient, _ = ethclient.Dial("https://testnet.matic.network")

func getValidatorByIndex(_index int64) abci.Validator {
	stakeManagerInstance, err := stakemanager.NewContracts(common.HexToAddress(stakeManagerAddress), kovanClient)
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
		// PubKey:  tmtypes.TM2PB.PubKey(_pubkey),
	}
	return abciValidator
}

func getLastValidator() int64 {
	stakeManagerInstance, err := stakemanager.NewContracts(common.HexToAddress(stakeManagerAddress), kovanClient)
	if err != nil {
		log.Fatal(err)
	}

	last, _ := stakeManagerInstance.LastValidatorIndex(nil)
	return last.Int64()
}

//
//// SendCheckpoint sends transaction to main chain
//func SendCheckpoint(start int, end int, sigs []byte) {
//	clientKovan := initKovan()
//	clientMatic := initMatic()
//	rootHash := getHeaders(start, end, clientMatic)
//	fmt.Printf("Root hash obtained for blocks from %v to %v is %v", start, end, rootHash)
//	rootchainClient, err := rootmock.NewContracts(common.HexToAddress(rootchainAddress), clientKovan)
//	if err != nil {
//		panic(err)
//	}
//
//	auth := GenerateAuthObj(clientKovan)
//	auth.Value = big.NewInt(0)
//
//	// Calling contract method
//	var amount big.Int
//	amount.SetUint64(0)
//	var rootHashArray [32]byte
//	copy(rootHashArray[:], rootHash)
//	tx, err := rootchainClient.SubmitHeaderBlock(auth, rootHashArray, big.NewInt(int64(start)), big.NewInt(int64(end)), sigs)
//	if err != nil {
//		fmt.Printf("Transaction unable to send error %v", err)
//	}
//	fmt.Printf("Checkpoint sent successfully %v", tx)
//
//}

// SubmitProof submit header
func SubmitProof(voteSignBytes []byte, sigs []byte, extradata []byte, start uint64, end uint64, rootHash common.Hash) {
	fmt.Printf("Root hash obtained for blocks from %v to %v is %v", start, end, rootHash)
	//auth := GenerateAuthObj(clientKovan)
	//auth.Value = big.NewInt(0)
	//todo change this to tx , right now its a call
	validatorSetInstance := getValidatorSetInstance(kovanClient)
	fmt.Printf("inputs , vote: %v , sigs: %v , extradata %v ", hex.EncodeToString(voteSignBytes), hex.EncodeToString(sigs), hex.EncodeToString(extradata))
	res, proposer, error := validatorSetInstance.Validate(nil, voteSignBytes, sigs, extradata)

	fmt.Printf("Submitted Proof Successfully %v %v %v ", res, proposer.String(), error)
}

func getHeaders(start int, end int, client *ethclient.Client) string {
	if (start-end+1)%2 != 0 {
		return "Not Defined , make sure difference is even "
	}
	if start > end {
		return ""
	}
	current := start
	var result [][32]byte
	for current <= end {
		//TODO run this in different goroutines and use channels to fetch results(how to maintian order)
		blockheader, err := client.HeaderByNumber(context.Background(), big.NewInt(int64(current)))
		if err != nil {
			fmt.Printf("not found")
			log.Fatal(err)
		}
		headerBytes := appendBytes32(blockheader.Number.Bytes(),
			blockheader.Time.Bytes(),
			blockheader.TxHash.Bytes(),
			blockheader.ReceiptHash.Bytes())

		header := getsha3frombyte(headerBytes)
		var arr [32]byte
		copy(arr[:], header)
		result = append(result, arr)
		current++
	}
	merkelData := convert(result)
	tree := merkle.NewTreeWithOpts(merkle.TreeOptions{EnableHashSorting: false, DisableHashLeaves: true})

	err := tree.Generate(merkelData, sha3.NewKeccak256())
	if err != nil {
		fmt.Println("*********ERROR***********")
		log.Fatal(err)
	}
	return hex.EncodeToString(tree.Root().Hash)
}
func convert(input []([32]byte)) [][]byte {
	var output [][]byte
	for _, in := range input {
		newInput := make([]byte, len(in[:]))
		copy(newInput, in[:])
		output = append(output, newInput)
	}
	return output
}

func convertTo32(input []byte) (output [32]byte, err error) {
	l := len(input)
	if l > 32 || l == 0 {
		err = fmt.Errorf("input length is greater than 32")
		return
	}
	copy(output[32-l:], input[:])
	return
}

func appendBytes32(data ...[]byte) []byte {
	var result []byte
	for _, v := range data {
		paddedV, err := convertTo32(v)
		if err == nil {
			result = append(result, paddedV[:]...)
		}
	}
	return result
}
func getsha3frombyte(input []byte) []byte {
	hash := sha3.NewKeccak256()
	var buf []byte
	hash.Write(input)
	buf = hash.Sum(buf)
	return buf
}

func getValidatorSetInstance(client *ethclient.Client) *validatorset.ValidatorSet {
	validatorSetInstance, err := validatorset.NewValidatorSet(common.HexToAddress(validatorSetAddress), client)
	if err != nil {
		log.Fatal(err)
	}
	return validatorSetInstance
}

func GetProposer() common.Address {
	validatorSetInstance := getValidatorSetInstance(kovanClient)
	currentProposer, err := validatorSetInstance.Proposer(nil)
	if err != nil {
		fmt.Printf("error getting proposer")
	}
	return currentProposer
}

func SelectProposer() {
	validatorSetInstance := getValidatorSetInstance(kovanClient)
	auth := GenerateAuthObj(kovanClient)
	tx, err := validatorSetInstance.SelectProposer(auth)
	if err != nil {
		fmt.Printf("Unable to send transaction for proposer selection ")
	}
	fmt.Printf("New Proposer Selected ! %v", tx)
	//proposer := t
}

func GenerateAuthObj(client *ethclient.Client) (auth *bind.TransactOpts) {
	privVal := privval.LoadFilePV("/Users/vc/.basecoind/config/priv_validator.json")
	var pkObject secp256k1.PrivKeySecp256k1
	cdc.MustUnmarshalBinaryBare(privVal.PrivKey.Bytes(), &pkObject)

	// create ecdsa private key
	ecdsaPrivateKey, err := crypto.ToECDSA(pkObject[:])
	if err != nil {
		panic(err)
	}

	// from address
	fromAddress := common.BytesToAddress(privVal.Address)
	// fetch gas price
	gasprice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		panic(err)
	}

	// fetch nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		panic(err)
	}

	// create auth
	auth = bind.NewKeyedTransactor(ecdsaPrivateKey)
	auth.GasPrice = gasprice
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = uint64(3000000)
	return auth
}

func GetValidators() (validators []abci.Validator) {
	validatorSetInstance := getValidatorSetInstance(kovanClient)
	powers, ValidatorAddrs, err := validatorSetInstance.GetValidatorSet(nil)
	if err != nil {
		fmt.Printf(" The error is %v", err)
	}

	for index := range powers {
		pubkey, error := validatorSetInstance.GetPubkey(nil, big.NewInt(int64(index)))
		if error != nil {
			fmt.Errorf(" Error getting pubkey for index %v", error)
		}

		var pubkeyBytes secp256k1.PubKeySecp256k1
		_pubkey, _ := hex.DecodeString(pubkey)
		copy(pubkeyBytes[:], _pubkey)
		// todo add a check to check pubkey corresponds to address
		validator := abci.Validator{
			Address: ValidatorAddrs[index].Bytes(),
			Power:   powers[index].Int64(),
			// PubKey:  tmtypes.TM2PB.PubKey(pubkeyBytes),
		}
		fmt.Printf("New Validator is %v", validator)
		validators = append(validators, validator)
		//validatorPubKeys[index] = pubkey
	}

	return validators
}
