package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/console/prompt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	hmModule "github.com/maticnetwork/heimdall/types/module"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/privval"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/file"
	"github.com/maticnetwork/heimdall/version"

	"github.com/maticnetwork/heimdall/app"
	authCli "github.com/maticnetwork/heimdall/auth/client/cli"
	hmTxCli "github.com/maticnetwork/heimdall/client/tx"
	"github.com/maticnetwork/heimdall/helper"
)

var logger = helper.Logger.With("module", "cmd/heimdallcli")

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "heimdallcli",
		Short: "Heimdall light-client",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if cmd.Use != version.Cmd.Use {
				// initialise config
				initTendermintViperConfig(cmd)
			}
			return nil
		},
	}
)

func initTendermintViperConfig(cmd *cobra.Command) {
	tendermintNode, _ := cmd.Flags().GetString(helper.TendermintNodeFlag)
	homeValue, _ := cmd.Flags().GetString(helper.HomeFlag)

	// set to viper
	viper.Set(helper.TendermintNodeFlag, tendermintNode)
	viper.Set(helper.HomeFlag, homeValue)

	// start heimdall config
	helper.InitHeimdallConfig("")
}

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastSync
	cliCtx.TrustNode = true

	// TODO: Setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc.

	// chain id
	rootCmd.PersistentFlags().String(client.FlagChainID, "", "Chain ID of tendermint node")

	helper.DecorateWithHeimdallFlags(rootCmd, viper.GetViper(), logger, "main")

	// add query/post commands (custom to binary)
	rootCmd.AddCommand(
		rpc.StatusCommand(),
		client.LineBreak,
		queryCmd(cdc),
		txCmd(cdc),
		client.LineBreak,
		keys.Commands(),
		exportCmd(ctx, cdc),
		convertAddressToHexCmd(cdc),
		convertHexToAddressCmd(cdc),
		generateKeystore(cdc),
		generateValidatorKey(cdc),
		client.LineBreak,
		version.Cmd,
		client.LineBreak,

		// approve and stake on mainnet
		StakeCmd(cliCtx),
		ApproveCmd(cliCtx),
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "HD", os.ExpandEnv("/var/lib/heimdall"))
	if err := executor.Execute(); err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

func queryCmd(cdc *amino.Codec) *cobra.Command {
	queryCmd := &cobra.Command{
		Use:     "query",
		Aliases: []string{"q"},
		Short:   "Querying subcommands",
	}

	queryCmd.AddCommand(
		rpc.ValidatorCommand(cdc),
		rpc.BlockCommand(),
		hmTxCli.QueryTxsByEventsCmd(cdc),
		hmTxCli.QueryTxCmd(cdc),
	)

	// add modules' query commands
	app.ModuleBasics.AddQueryCommands(queryCmd, cdc)

	return queryCmd
}

func txCmd(cdc *amino.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:   "tx",
		Short: "Transactions subcommands",
	}

	txCmd.AddCommand(
		authCli.GetSignCommand(cdc),
		hmTxCli.GetBroadcastCommand(cdc),
		hmTxCli.GetEncodeCommand(cdc),
		client.LineBreak,
	)

	// add modules' tx commands
	app.ModuleBasics.AddTxCommands(txCmd, cdc)

	return txCmd
}

func convertAddressToHexCmd(_ *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "address-to-hex [address]",
		Short: "Convert address to hex",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			key, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			fmt.Println("Hex:", ethCommon.BytesToAddress(key).String())
			return nil
		},
	}

	return client.GetCommands(cmd)[0]
}

func convertHexToAddressCmd(_ *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hex-to-address [hex]",
		Short: "Convert hex to address",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			address := ethCommon.HexToAddress(args[0])
			fmt.Println("Address:", sdk.AccAddress(address.Bytes()).String())
			return nil
		},
	}

	return client.GetCommands(cmd)[0]
}

// exportCmd a state dump file
func exportCmd(ctx *server.Context, _ *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-heimdall",
		Short: "Export genesis file with state-dump",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {

			// cliCtx := context.NewCLIContext().WithCodec(cdc)
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			// create chain id
			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("heimdall-%v", common.RandStr(6))
			}

			dataDir := path.Join(viper.GetString(cli.HomeFlag), "data")
			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
			db, err := sdk.NewLevelDB("application", dataDir)
			if err != nil {
				panic(err)
			}

			happ := app.NewHeimdallApp(logger, db)

			marshaledAppState := []byte{}

			// Anonymous function just to ensure that appState is not retained in memory
			func() {
				appState, err := getAppState(happ)
				if err != nil {
					panic(err)
				}
				runtime.GC()

				sdkCtx := happ.NewContext(true, abci.Header{Height: happ.LastBlockHeight()})
				moduleManager := happ.GetModuleManager()

				for _, moduleName := range moduleManager.OrderExportGenesis {
					module, ok := moduleManager.Modules[moduleName].(hmModule.StreamedGenesisExporter)
					if !ok {
						continue
					}

					runtime.GC()

					if err := fetchModuleData(appState, module, sdkCtx); err != nil {
						panic(err)
					}
				}

				runtime.GC()

				marshaledAppState, err = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(appState)
				if err != nil {
					panic(err)
				}

				appState = nil
			}()

			runtime.GC()

			savePath := file.Rootify("dump-genesis.json", config.RootDir)
			if err := writeGenesisFile(savePath, chainID, marshaledAppState); err != nil {
				panic(err)
			}

			fmt.Println("New genesis json file created: ", savePath)

			return nil
		},
	}
	cmd.Flags().String(cli.HomeFlag, helper.DefaultNodeHome, "Node's home directory")
	cmd.Flags().String(helper.FlagClientHome, helper.DefaultCLIHome, "Client's home directory")
	cmd.Flags().String(client.FlagChainID, "", "Genesis file chain-id, if left blank will be randomly created")

	return cmd
}

// getAppState generates and returns the app state.
func getAppState(happ *app.HeimdallApp) (map[string]interface{}, error) {
	appState, _, err := happ.ExportAppStateAndValidators()
	if err != nil {
		return nil, err
	}

	appStateData, err := appState.MarshalJSON()
	if err != nil {
		return nil, err
	}

	unmarshaledData := make(map[string]interface{})
	if err := json.Unmarshal(appStateData, &unmarshaledData); err != nil {
		return nil, err
	}

	return unmarshaledData, nil
}

// fetchModuleData fetches module genesis data in streamed fashion.
func fetchModuleData(appData map[string]interface{}, module hmModule.StreamedGenesisExporter, sdkCtx sdk.Context) error {
	var lastKey []byte
	allData := []json.RawMessage{}
	allDataLength := 0
	var currAppendingPath string

	for {
		data, err := module.NextGenesisData(sdkCtx, lastKey)
		if err != nil {
			panic(err)
		}

		lastKey = data.LastKey

		if lastKey == nil {
			allData = append(allData, data.Data)
			allDataLength += len(data.Data)

			if allDataLength == 0 {
				break
			}

			combinedData, err := combineJSONArrays(allData, allDataLength)
			if err != nil {
				return err
			}

			if err := AddProperty(appData, currAppendingPath, combinedData); err != nil {
				return err
			}

			break
		}

		if currAppendingPath != "" && currAppendingPath != data.Path {
			combinedData, err := combineJSONArrays(allData, allDataLength)
			if err != nil {
				return err
			}

			if err := AddProperty(appData, currAppendingPath, combinedData); err != nil {
				return err
			}

			allData = []json.RawMessage{}
			allDataLength = 0
		}

		currAppendingPath = data.Path

		allData = append(allData, data.Data)
		allDataLength += len(data.Data)
	}

	return nil
}

// combineJSONArrays combines multiple JSON arrays into a single JSON array.
func combineJSONArrays(arrays []json.RawMessage, allArraysLength int) (json.RawMessage, error) {
	buf := bytes.NewBuffer(make([]byte, 0, allArraysLength))
	buf.WriteByte('[')
	first := true

	for _, raw := range arrays {
		if len(raw) == 0 {
			continue
		}

		if raw[0] != '[' || raw[len(raw)-1] != ']' {
			return nil, fmt.Errorf("invalid JSON array: %s", raw)
		}

		content := raw[1 : len(raw)-1]

		if !first {
			buf.WriteByte(',')
		}
		buf.Write(content)
		first = false
	}
	buf.WriteByte(']')

	combinedJSON := buf.Bytes()
	if !json.Valid(combinedJSON) {
		return nil, errors.New("combined JSON is invalid")
	}

	return json.RawMessage(combinedJSON), nil
}

// AddProperty adds a property to the data map.
func AddProperty(data map[string]interface{}, path string, value json.RawMessage) error {
	if path == "" {
		return fmt.Errorf("path cannot be empty")
	}

	keys := strings.Split(path, ".")
	lastKey := keys[len(keys)-1]
	parentPath := strings.Join(keys[:len(keys)-1], ".")

	current, err := traversePath(data, parentPath)
	if err != nil {
		return err
	}
	current[lastKey] = value
	return nil
}

// traversePath traverses the path in the data map.
func traversePath(data map[string]interface{}, path string) (map[string]interface{}, error) {
	if path == "." {
		return data, nil
	}

	keys := strings.Split(path, ".")
	current := data

	for _, key := range keys {
		if next, ok := current[key].(map[string]interface{}); ok {
			current = next
			continue
		}
		return nil, fmt.Errorf("invalid path: %s", path)
	}

	return current, nil
}

// generateKeystore generate keystore file from private key
func generateKeystore(_ *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-keystore <private-key>",
		Short: "Generates keystore file using private key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := strings.ReplaceAll(args[0], "0x", "")
			pk, err := crypto.HexToECDSA(s)
			if err != nil {
				return err
			}

			id, err := uuid.NewRandom()
			if err != nil {
				return err
			}
			key := &keystore.Key{
				Id:         id,
				Address:    crypto.PubkeyToAddress(pk.PublicKey),
				PrivateKey: pk,
			}

			passphrase, err := promptPassphrase(true)
			if err != nil {
				return err
			}

			keyjson, err := keystore.EncryptKey(key, passphrase, keystore.StandardScryptN, keystore.StandardScryptP)
			if err != nil {
				return err
			}

			// Then write the new keyfile in place of the old one.
			if err := os.WriteFile(keyFileName(key.Address), keyjson, 0600); err != nil {
				return err
			}
			return nil
		},
	}

	return client.GetCommands(cmd)[0]
}

// generateValidatorKey generate validator key
func generateValidatorKey(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate-validatorkey <private-key>",
		Short: "Generate validator key file using private key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			s := strings.ReplaceAll(args[0], "0x", "")
			ds, err := hex.DecodeString(s)
			if err != nil {
				return err
			}

			// set private object
			var privObject secp256k1.PrivKeySecp256k1
			copy(privObject[:], ds)

			// node key
			nodeKey := privval.FilePVKey{
				Address: privObject.PubKey().Address(),
				PubKey:  privObject.PubKey(),
				PrivKey: privObject,
			}

			jsonBytes, err := cdc.MarshalJSONIndent(nodeKey, "", "  ")
			if err != nil {
				return err
			}

			err = os.WriteFile("priv_validator_key.json", jsonBytes, 0600)
			if err != nil {
				return err
			}

			return nil
		},
	}

	return client.GetCommands(cmd)[0]
}

//
// Internal functions
//

func writeGenesisFile(genesisFile, chainID string, appState json.RawMessage) error {
	genDoc := tmTypes.GenesisDoc{
		ChainID:  chainID,
		AppState: appState,
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}

	return genDoc.SaveAs(genesisFile)
}

// keyFileName implements the naming convention for keyfiles:
// UTC--<created_at UTC ISO8601>-<address hex>
func keyFileName(keyAddr ethCommon.Address) string {
	ts := time.Now().UTC()
	return fmt.Sprintf("UTC--%s--%s", toISO8601(ts), hex.EncodeToString(keyAddr[:]))
}

func toISO8601(t time.Time) string {
	var tz string

	name, offset := t.Zone()
	if name == "UTC" {
		tz = "Z"
	} else {
		tz = fmt.Sprintf("%03d00", offset/3600)
	}

	return fmt.Sprintf("%04d-%02d-%02dT%02d-%02d-%02d.%09d%s",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), tz)
}

// promptPassphrase prompts the user for a passphrase.  Set confirmation to true
// to require the user to confirm the passphrase.
func promptPassphrase(confirmation bool) (string, error) {
	passphrase, err := prompt.Stdin.PromptPassword("Passphrase: ")
	if err != nil {
		return "", err
	}

	if confirmation {
		confirm, err := prompt.Stdin.PromptPassword("Repeat passphrase: ")
		if err != nil {
			return "", err
		}

		if passphrase != confirm {
			return "", errors.New("Passphrases do not match")
		}
	}

	return passphrase, nil
}
