package pier

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/maticnetwork/heimdall/helper"
	hmtypes "github.com/maticnetwork/heimdall/types"
	rest "github.com/maticnetwork/heimdall/types/rest"
	"github.com/tendermint/tendermint/libs/log"
)

const (
	ChainSyncer          = "chain-syncer"
	HeimdallCheckpointer = "heimdall-checkpointer"
	NoackService         = "checkpoint-no-ack"

	// TODO fetch port from config
	LastNoAckURL          = "http://localhost:1317/checkpoint/last-no-ack"
	ProposersURL          = "http://localhost:1317/staking/proposer/%v"
	BufferedCheckpointURL = "http://localhost:1317/checkpoint/buffer"
	LatestCheckpointURL   = "http://localhost:1317/checkpoint/latest-checkpoint"
	TendermintBlockURL    = "http://localhost:26657/block?height=%v"
	CurrentProposerURL    = "http://localhost:1317/staking/current-proposer"

	TransactionTimeout = 1 * time.Minute
	CommitTimeout      = 2 * time.Minute

	BridgeDBFlag = "bridge-db"
)

// Big batch of reflect types for topic reconstruction.
var (
	reflectHash    = reflect.TypeOf(common.Hash{})
	reflectAddress = reflect.TypeOf(common.Address{})
	reflectBigInt  = reflect.TypeOf(new(big.Int))
)

// Global logger for bridge
var Logger log.Logger

func init() {
	Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

// checks if we are proposer
func isProposer(cliCtx cliContext.CLIContext) bool {
	count := uint64(1)
	resp, err := http.Get(fmt.Sprintf(ProposersURL, strconv.FormatUint(count, 10)))
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false
		}
		// unmarshall data from buffer
		var proposers []hmtypes.Validator
		var response rest.ResponseWithHeight

		if err := cliCtx.Codec.UnmarshalJSON(body, &response); err != nil {
			return false
		}

		if err := json.Unmarshal(response.Result, &proposers); err != nil {
			return false
		}

		// no proposer found
		if len(proposers) == 0 {
			return false
		}
		// get first proposer
		proposer := proposers[0]
		if bytes.Equal(proposer.Signer.Bytes(), helper.GetAddress()) {
			return true
		}
	} else {
		Logger.Error("Error while fetching proposer", "status", resp.StatusCode)
	}
	return false
}

// UnpackLog unpacks log
func UnpackLog(abiObject *abi.ABI, out interface{}, event string, log *types.Log) error {
	if len(log.Data) > 0 {
		if err := abiObject.Unpack(out, event, log.Data); err != nil {
			return err
		}
	}
	var indexed abi.Arguments
	for _, arg := range abiObject.Events[event].Inputs {
		if arg.Indexed {
			indexed = append(indexed, arg)
		}
	}
	return parseTopics(out, indexed, log.Topics[1:])
}

// parseTopics converts the indexed topic fields into actual log field values.
//
// Note, dynamic types cannot be reconstructed since they get mapped to Keccak256
// hashes as the topic value!
func parseTopics(out interface{}, fields abi.Arguments, topics []common.Hash) error {
	// Sanity check that the fields and topics match up
	if len(fields) != len(topics) {
		return errors.New("topic/field count mismatch")
	}

	// Iterate over all the fields and reconstruct them from topics
	for _, arg := range fields {
		if !arg.Indexed {
			return errors.New("non-indexed field in topic reconstruction")
		}
		field := reflect.ValueOf(out).Elem().FieldByName(capitalise(arg.Name))

		// Try to parse the topic back into the fields based on primitive types
		switch field.Kind() {
		case reflect.Bool:
			if topics[0][common.HashLength-1] == 1 {
				field.Set(reflect.ValueOf(true))
			}
		case reflect.Int8:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int8(num.Int64())))

		case reflect.Int16:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int16(num.Int64())))

		case reflect.Int32:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(int32(num.Int64())))

		case reflect.Int64:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(num.Int64()))

		case reflect.Uint8:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint8(num.Uint64())))

		case reflect.Uint16:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint16(num.Uint64())))

		case reflect.Uint32:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(uint32(num.Uint64())))

		case reflect.Uint64:
			num := new(big.Int).SetBytes(topics[0][:])
			field.Set(reflect.ValueOf(num.Uint64()))

		default:
			// Ran out of plain primitive types, try custom types
			switch field.Type() {
			case reflectHash: // Also covers all dynamic types
				field.Set(reflect.ValueOf(topics[0]))

			case reflectAddress:
				var addr common.Address
				copy(addr[:], topics[0][common.HashLength-common.AddressLength:])
				field.Set(reflect.ValueOf(addr))

			case reflectBigInt:
				num := new(big.Int).SetBytes(topics[0][:])
				field.Set(reflect.ValueOf(num))

			default:
				// Ran out of custom types, try the crazies
				switch {
				case arg.Type.T == abi.FixedBytesTy:
					reflect.Copy(field, reflect.ValueOf(topics[0][common.HashLength-arg.Type.Size:]))

				default:
					return fmt.Errorf("unsupported indexed type: %v", arg.Type)
				}
			}
		}
		topics = topics[1:]
	}
	return nil
}

// capitalise makes a camel-case string which starts with an upper case character.
func capitalise(input string) string {
	for len(input) > 0 && input[0] == '_' {
		input = input[1:]
	}
	if len(input) == 0 {
		return ""
	}
	return toCamelCase(strings.ToUpper(input[:1]) + input[1:])
}

// decapitalise makes a camel-case string which starts with a lower case character.
func decapitalise(input string) string {
	for len(input) > 0 && input[0] == '_' {
		input = input[1:]
	}
	if len(input) == 0 {
		return ""
	}
	return toCamelCase(strings.ToLower(input[:1]) + input[1:])
}

// toCamelCase converts an under-score string to a camel-case string
func toCamelCase(input string) string {
	toupper := false

	result := ""
	for k, v := range input {
		switch {
		case k == 0:
			result = strings.ToUpper(string(input[0]))

		case toupper:
			result += strings.ToUpper(string(v))
			toupper = false

		case v == '_':
			toupper = true

		default:
			result += string(v)
		}
	}
	return result
}
