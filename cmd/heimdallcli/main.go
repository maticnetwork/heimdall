package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	defaultLogger "log"
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
				defaultLogger.SetOutput(os.Stdout)
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
		Short: "Export genesis file with state-dump. It expects --home and --chain-id flags to be set",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {

			config := ctx.Config
			if viper.GetString(cli.HomeFlag) == "" {
				panic("home flag is not set")
			}
			config.SetRoot(viper.GetString(cli.HomeFlag))

			// create chain id and genesis time
			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				panic("chain-id flag is not set")
			}

			genesisTimes := map[string]string{
				"heimdall-137":   "2020-05-30T04:28:03.177054Z",    // mainnet
				"heimdall-80001": "2020-06-04T10:47:20.806665Z",    // mumbai
				"heimdall-80002": "2023-11-06T06:41:35.410487141Z", // amoy
				"devnet":         "2025-01-01T00:00:00.000000000Z", // local devnet
			}

			genesisTime, ok := genesisTimes[chainID]
			if !ok {
				panic("invalid chain-id, it must be one of: heimdall-137 (mainnet), heimdall-80001 (mumbai), heimdall-80002 (amoy), devnet (local)")
			}

			dataDir := path.Join(viper.GetString(cli.HomeFlag), "data")
			logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout))
			db, err := sdk.NewLevelDB("application", dataDir)
			if err != nil {
				panic(err)
			}

			happ := app.NewHeimdallApp(logger, db)

			savePath := file.Rootify("dump-genesis.json", config.RootDir)
			file, err := os.Create(savePath)
			if err != nil {
				panic(err)
			}
			defer file.Close()

			if err := generateMarshalledAppState(happ, chainID, genesisTime, 1000, file); err != nil {
				panic(err)
			}

			fmt.Println("New genesis json file created: ", savePath)

			return nil
		},
	}
	cmd.Flags().String(cli.HomeFlag, helper.DefaultNodeHome, "Node's home directory")
	cmd.Flags().String(helper.FlagClientHome, helper.DefaultCLIHome, "Client's home directory")
	cmd.Flags().String(client.FlagChainID, "devnet", "Genesis file chain-id, it can be "+
		"heimdall-137 (for mainnet), "+
		"heimdall-80001 (for mumbai), "+
		"heimdall-80002 (for amoy), "+
		"devnet (for any local devnet)")

	return cmd
}

// generateMarshalledAppState writes the genesis doc with app state directly to a file to minimize memory usage.
func generateMarshalledAppState(happ *app.HeimdallApp, chainID, genesisTime string, maxNextGenesisItems int, w io.Writer) error {
	sdkCtx := happ.NewContext(true, abci.Header{Height: happ.LastBlockHeight()})
	moduleManager := happ.GetModuleManager()

	if _, err := w.Write([]byte("{")); err != nil {
		return err
	}

	if _, err := w.Write([]byte(`"app_state":`)); err != nil {
		return err
	}

	if _, err := w.Write([]byte(`{`)); err != nil {
		return err
	}

	isFirst := true

	for _, moduleName := range moduleManager.OrderExportGenesis {
		runtime.GC()

		if !isFirst {
			if _, err := w.Write([]byte(`,`)); err != nil {
				return err
			}
		}

		isFirst = false

		if _, err := w.Write([]byte(`"` + moduleName + `":`)); err != nil {
			return err
		}

		module, isStreamedGenesis := moduleManager.Modules[moduleName].(hmModule.StreamedGenesisExporter)
		if isStreamedGenesis {
			partialGenesis, err := module.ExportPartialGenesis(sdkCtx)
			if err != nil {
				return err
			}

			propertyName, data, err := fetchModuleStreamedData(sdkCtx, module, maxNextGenesisItems)
			if err != nil {
				return err
			}

			// remove the closing '}'
			if _, err = w.Write(partialGenesis[0 : len(partialGenesis)-1]); err != nil {
				return err
			}

			if _, err = w.Write([]byte(`,`)); err != nil {
				return err
			}

			if _, err = w.Write([]byte(`"` + propertyName + `":`)); err != nil {
				return err
			}

			if _, err = w.Write(data); err != nil {
				return err
			}

			// add the closing '}'
			if _, err = w.Write(partialGenesis[len(partialGenesis)-1:]); err != nil {
				return err
			}

			continue
		}

		genesis := moduleManager.Modules[moduleName].ExportGenesis(sdkCtx)

		if _, err := w.Write(genesis); err != nil {
			return err
		}
	}

	if _, err := w.Write([]byte(`}`)); err != nil {
		return err
	}

	if _, err := w.Write([]byte(`,`)); err != nil {
		return err
	}

	consensusParams := tmTypes.DefaultConsensusParams()
	genesisTime := time.Now().UTC().Format(time.RFC3339Nano)

	consensusParamsData, err := tmTypes.GetCodec().MarshalJSON(consensusParams)
	if err != nil {
		return err
	}

	remainingFields := map[string]interface{}{
		"chain_id":         chainID,
		"consensus_params": json.RawMessage(consensusParamsData),
		"genesis_time":     genesisTime,
	}

	remainingFieldsData, err := json.Marshal(remainingFields)
	if err != nil {
		return err
	}

	if _, err := w.Write(remainingFieldsData[1 : len(remainingFieldsData)-1]); err != nil {
		return err
	}

	if _, err := w.Write([]byte("}")); err != nil {
		return err
	}

	return nil
}

// fetchModuleStreamedData fetches module genesis data in streamed fashion.
func fetchModuleStreamedData(sdkCtx sdk.Context, module hmModule.StreamedGenesisExporter, maxNextGenesisItems int) (string, json.RawMessage, error) {
	var lastKey []byte
	allData := []json.RawMessage{}
	allDataLength := 0

	for {
		data, err := module.NextGenesisData(sdkCtx, lastKey, maxNextGenesisItems)
		if err != nil {
			panic(err)
		}

		lastKey = data.NextKey

		if lastKey == nil {
			allData = append(allData, data.Data)
			allDataLength += len(data.Data)

			if allDataLength == 0 {
				break
			}

			combinedData, err := combineJSONArrays(allData, allDataLength)
			if err != nil {
				return "", nil, err
			}

			return data.Path, combinedData, nil
		}

		allData = append(allData, data.Data)
		allDataLength += len(data.Data)
	}

	return "", nil, errors.New("failed to iterate module genesis data")
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
