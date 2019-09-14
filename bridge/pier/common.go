package pier

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
	httpClient "github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/helper"
	hmtypes "github.com/maticnetwork/heimdall/types"
	rest "github.com/maticnetwork/heimdall/types/rest"
)

const (
	ChainSyncer          = "chain-syncer"
	HeimdallCheckpointer = "heimdall-checkpointer"
	NoackService         = "checkpoint-no-ack"
	SpanServiceStr       = "span-service"
	AMQPConsumerService  = "amqp-consumer-service"

	// TxsURL represents txs url
	TxsURL = "/txs"

	AccountDetailsURL     = "/auth/accounts/%v"
	LastNoAckURL          = "/checkpoint/last-no-ack"
	ProposersURL          = "/staking/proposer/%v"
	BufferedCheckpointURL = "/checkpoint/buffer"
	LatestCheckpointURL   = "/checkpoint/latest-checkpoint"
	CurrentProposerURL    = "/staking/current-proposer"
	LatestSpanURL         = "/bor/latest-span"
	SpanProposerURL       = "/bor/span-proposer"
	NextSpanInfoURL       = "/bor/prepare-next-span"

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
	var proposers []hmtypes.Validator
	count := uint64(1)
	result, err := FetchFromAPI(cliCtx,
		GetHeimdallServerEndpoint(fmt.Sprintf(ProposersURL, strconv.FormatUint(count, 10))),
	)
	if err != nil {
		Logger.Error("Error fetching proposers", "error", err)
		return false
	}

	err = json.Unmarshal(result.Result, &proposers)
	if err != nil {
		Logger.Error("error unmarshalling proposer slice", "error", err)
		return false
	}

	if bytes.Equal(proposers[0].Signer.Bytes(), helper.GetAddress()) {
		return true
	}

	return false
}

// GetHeimdallServerEndpoint returns heimdall server endpoint
func GetHeimdallServerEndpoint(endpoint string) string {
	u, _ := url.Parse(helper.GetConfig().HeimdallServerURL)
	u.Path = path.Join(u.Path, endpoint)
	return u.String()
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

// FetchFromAPI fetches data from any URL
func FetchFromAPI(cliCtx cliContext.CLIContext, URL string) (result rest.ResponseWithHeight, err error) {
	resp, err := http.Get(URL)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result, err
		}
		// unmarshall data from buffer
		// var proposers []hmtypes.Validator
		var response rest.ResponseWithHeight
		if err := cliCtx.Codec.UnmarshalJSON(body, &response); err != nil {
			return result, err
		}
		return response, nil
	} else {
		Logger.Error("Error while fetching data from URL", "status", resp.StatusCode, "URL", URL)
		return result, err
	}
}

// WaitForOneEvent subscribes to a websocket event for the given
// event time and returns upon receiving it one time, or
// when the timeout duration has expired.
//
// This handles subscribing and unsubscribing under the hood
func WaitForOneEvent(tx tmTypes.Tx, client *httpClient.HTTP) (tmTypes.TMEventData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CommitTimeout)
	defer cancel()

	// subscriber
	subscriber := hex.EncodeToString(tx.Hash())

	// query
	query := tmTypes.EventQueryTxFor(tx).String()

	// register for the next event of this type
	eventCh, err := client.Subscribe(ctx, subscriber, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe")
	}

	// make sure to unregister after the test is over
	defer client.UnsubscribeAll(ctx, subscriber)
	select {
	case event := <-eventCh:
		return event.Data.(tmTypes.TMEventData), nil
	case <-ctx.Done():
		return nil, errors.New("timed out waiting for event")
	}
}

// FetchVotes fetches votes and extracts sigs from it
func FetchVotes(
	height int64,
	client *httpClient.HTTP,
) (votes []*tmTypes.CommitSig, sigs []byte, chainID string, err error) {
	// get block client
	blockDetails, err := helper.GetBlockWithClient(client, height+1)

	if err != nil {
		return nil, nil, "", err
	}

	// extract votes from response
	preCommits := blockDetails.LastCommit.Precommits

	// extract signs from votes
	valSigs := helper.GetSigs(preCommits)

	// extract chainID
	chainID = blockDetails.ChainID

	// return
	return preCommits, valSigs, chainID, nil
}
