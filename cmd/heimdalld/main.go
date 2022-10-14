package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	tmTypes "github.com/tendermint/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"
	dbm "github.com/tendermint/tm-db"

	"github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/helper"
	hmserver "github.com/maticnetwork/heimdall/server"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmModule "github.com/maticnetwork/heimdall/types/module"
)

var logger = helper.Logger.With("module", "cmd/heimdalld")

var (
	flagNodeDirPrefix    = "node-dir-prefix"
	flagNumValidators    = "v"
	flagNumNonValidators = "n"
	flagOutputDir        = "output-dir"
	flagNodeDaemonHome   = "node-daemon-home"
	flagNodeCliHome      = "node-cli-home"
	flagNodeHostPrefix   = "node-host-prefix"
)

const (
	nodeDirPerm = 0755
)

var ZeroIntString = big.NewInt(0).String()

// ValidatorAccountFormatter helps to print local validator account information
type ValidatorAccountFormatter struct {
	Address string `json:"address,omitempty" yaml:"address"`
	PrivKey string `json:"priv_key,omitempty" yaml:"priv_key"`
	PubKey  string `json:"pub_key,omitempty" yaml:"pub_key"`
}

// GetSignerInfo returns signer information
func GetSignerInfo(pub crypto.PubKey, priv []byte, cdc *codec.Codec) ValidatorAccountFormatter {
	var privObject secp256k1.PrivKeySecp256k1
	cdc.MustUnmarshalBinaryBare(priv, &privObject)
	return ValidatorAccountFormatter{
		Address: ethCommon.BytesToAddress(pub.Address().Bytes()).String(),
		PubKey:  CryptoKeyToPubkey(pub).String(),
		PrivKey: "0x" + hex.EncodeToString(privObject[:]),
	}
}

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	rootCmd := &cobra.Command{
		Use:               "heimdalld",
		Short:             "Heimdall Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	// add new persistent flag for heimdall-config
	rootCmd.PersistentFlags().String(
		helper.WithHeimdallConfigFlag,
		"",
		"Heimdall config file path (default <home>/config/heimdall-config.json)",
	)

	// bind with-heimdall-config config with root cmd
	if err := viper.BindPFlag(helper.WithHeimdallConfigFlag, rootCmd.Flags().Lookup(helper.WithHeimdallConfigFlag)); err != nil {
		logger.Error("main | BindPFlag | helper.WithHeimdallConfigFlag", "Error", err)
	}
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)
	rootCmd.AddCommand(showAccountCmd())
	rootCmd.AddCommand(showPrivateKeyCmd())
	rootCmd.AddCommand(hmserver.ServeCommands(cdc, hmserver.RegisterRoutes))
	rootCmd.AddCommand(VerifyGenesis(ctx, cdc))
	rootCmd.AddCommand(initCmd(ctx, cdc))
	rootCmd.AddCommand(testnetCmd(ctx, cdc))

	// rollback cmd
	rootCmd.AddCommand(rollbackCmd(ctx))

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "HD", os.ExpandEnv("$HOME/.heimdalld"))
	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
	// init heimdall config
	helper.InitHeimdallConfig("")
	// create new heimdall app
	return app.NewHeimdallApp(logger, db, baseapp.SetHaltHeight(cast.ToUint64(viper.GetString("halt-height"))), baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))))
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, storeTracer io.Writer, height int64, forZeroHeight bool, jailWhiteList []string) (json.RawMessage, []tmTypes.GenesisValidator, error) {
	bapp := app.NewHeimdallApp(logger, db)
	return bapp.ExportAppStateAndValidators()
}

func showAccountCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-account",
		Short: "Print the account's address and public key",
		Run: func(cmd *cobra.Command, args []string) {
			// init heimdall config
			helper.InitHeimdallConfig("")

			// get public keys
			pubObject := helper.GetPubKey()

			account := &ValidatorAccountFormatter{
				Address: ethCommon.BytesToAddress(pubObject.Address().Bytes()).String(),
				PubKey:  "0x" + hex.EncodeToString(pubObject[:]),
			}

			b, err := json.MarshalIndent(account, "", "    ")
			if err != nil {
				panic(err)
			}

			// prints json info
			fmt.Printf("%s", string(b))
		},
	}
}

func showPrivateKeyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-privatekey",
		Short: "Print the account's private key",
		Run: func(cmd *cobra.Command, args []string) {
			// init heimdall config
			helper.InitHeimdallConfig("")

			// get private and public keys
			privObject := helper.GetPrivKey()

			account := &ValidatorAccountFormatter{
				PrivKey: "0x" + hex.EncodeToString(privObject[:]),
			}

			b, err := json.MarshalIndent(account, "", "    ")
			if err != nil {
				panic(err)
			}

			// prints json info
			fmt.Printf("%s", string(b))
		},
	}
}

// VerifyGenesis verifies the genesis file and brings it in sync with on-chain contract
func VerifyGenesis(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-genesis",
		Short: "Verify if the genesis matches",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			helper.InitHeimdallConfig("")

			// Loading genesis doc
			genDoc, err := tmTypes.GenesisDocFromFile(filepath.Join(config.RootDir, "config/genesis.json"))
			if err != nil {
				return err
			}

			// get genesis state
			var genesisState app.GenesisState
			err = json.Unmarshal(genDoc.AppState, &genesisState)
			if err != nil {
				return err
			}

			// verify genesis
			for _, b := range app.ModuleBasics {
				m := b.(hmModule.HeimdallModuleBasic)
				if err := m.VerifyGenesis(genesisState); err != nil {
					return err
				}
			}

			return nil
		},
	}

	return cmd
}

// Total Validators to be included in the testnet
func totalValidators() int {
	numValidators := viper.GetInt(flagNumValidators)
	numNonValidators := viper.GetInt(flagNumNonValidators)
	return numNonValidators + numValidators
}

// get node directory path
func nodeDir(i int) string {
	outDir := viper.GetString(flagOutputDir)
	nodeDirName := fmt.Sprintf("%s%d", viper.GetString(flagNodeDirPrefix), i)
	nodeDaemonHomeName := viper.GetString(flagNodeDaemonHome)
	return filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
}

// hostname of ip of nodes
func hostnameOrIP(i int) string {
	return fmt.Sprintf("%s%d", viper.GetString(flagNodeHostPrefix), i)
}

// populate persistent peers in config
func populatePersistentPeersInConfigAndWriteIt(config *cfg.Config) {
	persistentPeers := make([]string, totalValidators())
	for i := 0; i < totalValidators(); i++ {
		config.SetRoot(nodeDir(i))
		nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
		if err != nil {
			return
		}
		persistentPeers[i] = p2p.IDAddressString(nodeKey.ID(), fmt.Sprintf("%s:%d", hostnameOrIP(i), 26656))
	}

	persistentPeersList := strings.Join(persistentPeers, ",")
	for i := 0; i < totalValidators(); i++ {
		config.SetRoot(nodeDir(i))
		config.P2P.PersistentPeers = persistentPeersList
		config.P2P.AddrBookStrict = false

		// overwrite default config
		cfg.WriteConfigFile(filepath.Join(nodeDir(i), "config", "config.toml"), config)
	}
}

func getGenesisAccount(address []byte) authTypes.GenesisAccount {
	acc := authTypes.NewBaseAccountWithAddress(hmTypes.BytesToHeimdallAddress(address))
	genesisBalance, _ := big.NewInt(0).SetString("1000000000000000000000", 10)
	if err := acc.SetCoins(sdk.Coins{sdk.Coin{Denom: authTypes.FeeToken, Amount: sdk.NewIntFromBigInt(genesisBalance)}}); err != nil {
		logger.Error("getGenesisAccount | SetCoins", "Error", err)
	}
	result, _ := authTypes.NewGenesisAccountI(&acc)
	return result
}

// WriteGenesisFile creates and writes the genesis configuration to disk. An
// error is returned if building or writing the configuration to file fails.
// nolint: unparam
func writeGenesisFile(genesisTime time.Time, genesisFile, chainID string, appState json.RawMessage) error {
	genDoc := tmTypes.GenesisDoc{
		GenesisTime: genesisTime,
		ChainID:     chainID,
		AppState:    appState,
	}

	if genDoc.GenesisTime.IsZero() {
		genDoc.GenesisTime = tmtime.Now()
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}

	return genDoc.SaveAs(genesisFile)
}

// InitializeNodeValidatorFiles initializes node and priv validator files
func InitializeNodeValidatorFiles(
	config *cfg.Config) (nodeID string, valPubKey crypto.PubKey, priv crypto.PrivKey, err error,
) {

	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return nodeID, valPubKey, priv, err
	}

	nodeID = string(nodeKey.ID())
	server.UpgradeOldPrivValFile(config)

	pvKeyFile := config.PrivValidatorKeyFile()
	if err := common.EnsureDir(filepath.Dir(pvKeyFile), 0777); err != nil {
		return nodeID, valPubKey, priv, err
	}

	pvStateFile := config.PrivValidatorStateFile()
	if err := common.EnsureDir(filepath.Dir(pvStateFile), 0777); err != nil {
		return nodeID, valPubKey, priv, err
	}

	FilePv := privval.LoadOrGenFilePV(pvKeyFile, pvStateFile)
	valPubKey = FilePv.GetPubKey()
	return nodeID, valPubKey, FilePv.Key.PrivKey, nil
}

// WriteDefaultHeimdallConfig writes default heimdall config to the given path
func WriteDefaultHeimdallConfig(path string, conf helper.Configuration) {
	heimdallConf := helper.GetDefaultHeimdallConfig()
	helper.WriteConfigFile(path, &heimdallConf)
}

func CryptoKeyToPubkey(key crypto.PubKey) hmTypes.PubKey {
	validatorPublicKey := helper.GetPubObjects(key)
	return hmTypes.NewPubKey(validatorPublicKey[:])
}
