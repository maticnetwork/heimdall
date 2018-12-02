package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	ethCommon "github.com/ethereum/go-ethereum/common"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	amino "github.com/tendermint/go-amino"
	cfg "github.com/tendermint/tendermint/config"
	tmCommon "github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/types"
	tmTime "github.com/tendermint/tendermint/types/time"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper"
)

var (
	flagNodeDirPrefix    = "node-dir-prefix"
	flagNumValidators    = "v"
	flagNumNonValidators = "n"
	flagOutputDir        = "output-dir"
	flagNodeDaemonHome   = "node-daemon-home"
	flagNodeCliHome      = "node-cli-home"
	flagNodeHostPrefix   = "node-host-prefix"
)

const nodeDirPerm = 0755

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()
	rootCmd := TestnetFilesCmd(ctx, cdc)
	viper.BindPFlags(rootCmd.Flags())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// TestnetFilesCmd returns cmd to initialize all files for tendermint testnet and application
func TestnetFilesCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "testnet",
		Short: "Initialize files for a Heimdall testnet",
		Long: `testnet will create "v" + "n" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.
Optionally, it will fill in persistent_peers list in config file using either hostnames or IPs.

Example:
testnet --v 4 --n 8 --output-dir ./output --starting-ip-address 192.168.10.2
`,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			return initTestnet(config, cdc)
		},
	}

	cmd.Flags().Int(flagNumValidators, 4,
		"Number of validators to initialize the testnet with",
	)

	cmd.Flags().Int(flagNumNonValidators, 8,
		"Number of validators to initialize the testnet with",
	)

	cmd.Flags().StringP(flagOutputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet",
	)

	cmd.Flags().String(flagNodeDirPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)",
	)

	cmd.Flags().String(flagNodeDaemonHome, "heimdalld",
		"Home directory of the node's daemon configuration",
	)

	cmd.Flags().String(flagNodeCliHome, "heimdallcli",
		"Home directory of the node's cli configuration",
	)

	cmd.Flags().String(flagNodeHostPrefix, "node",
		"Hostname prefix (node results in persistent peers list ID0@node0:26656, ID1@node1:26656, ...)")

	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")

	return cmd
}

func initTestnet(config *cfg.Config, cdc *codec.Codec) error {
	var chainID string
	outDir := viper.GetString(flagOutputDir)
	numValidators := viper.GetInt(flagNumValidators)
	numNonValidators := viper.GetInt(flagNumNonValidators)

	chainID = viper.GetString(client.FlagChainID)
	if chainID == "" {
		chainID = "heimdall-" + tmCommon.RandStr(6)
	}

	monikers := make([]string, totalValidators())
	nodeIDs := make([]string, totalValidators())
	valPubKeys := make([]*privval.FilePV, totalValidators())
	validators := make([]app.GenesisValidator, totalValidators())
	genFiles := make([]string, totalValidators())

	// generate private keys, node IDs, and initial transactions
	for i := 0; i < numValidators+numNonValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", viper.GetString(flagNodeDirPrefix), i)
		nodeDaemonHomeName := viper.GetString(flagNodeDaemonHome)
		nodeCliHomeName := viper.GetString(flagNodeCliHome)
		nodeDir := filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
		clientDir := filepath.Join(outDir, nodeDirName, nodeCliHomeName)

		config.SetRoot(nodeDir)

		err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		err = os.MkdirAll(clientDir, nodeDirPerm)
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}

		monikers[i] = nodeDirName
		config.Moniker = nodeDirName

		nodeIDs[i], valPubKeys[i] = initializeNodeValidatorFiles(config)
		genFiles[i] = config.GenesisFile()
		_, secret, err := server.GenerateCoinKey()
		if err != nil {
			_ = os.RemoveAll(outDir)
			return err
		}
		info := map[string]string{"secret": secret}
		cliPrint, err := json.Marshal(info)
		if err != nil {
			return err
		}

		// save private key seed words
		err = writeFile(fmt.Sprintf("%v.json", "key_seed"), clientDir, cliPrint)
		if err != nil {
			return err
		}

		// read or create private key
		_, pubKey := helper.GetPkObjects(valPubKeys[i].PrivKey)
		validators[i] = app.GenesisValidator{
			Address:    ethCommon.BytesToAddress(valPubKeys[i].Address),
			PubKey:     hmTypes.NewPubKey(pubKey[:]),
			StartEpoch: 0,
			Signer:     ethCommon.BytesToAddress(valPubKeys[i].Address),
			Power:      10,
		}
	}

	for i := 0; i < totalValidators(); i++ {
		populatePersistentPeersInConfigAndWriteIt(config)
	}

	if err := initGenFiles(chainID, validators, genFiles); err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Successfully initialized %d node directories\n", numValidators+numNonValidators)
	return nil
}

func initGenFiles(
	chainID string,
	validators []app.GenesisValidator,
	genFiles []string,
) error {
	appState := &app.GenesisState{
		Validators: validators,
	}

	appStateJSON, err := json.Marshal(appState)
	if err != nil {
		return err
	}

	genDoc := types.GenesisDoc{
		GenesisTime: tmTime.Now(),
		ChainID:     chainID,
		AppState:    appStateJSON,
		Validators:  nil,
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < len(validators); i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}

	return nil
}

func writeFile(name string, dir string, contents []byte) error {
	writePath := filepath.Join(dir)
	file := filepath.Join(writePath, name)

	err := tmCommon.EnsureDir(writePath, 0700)
	if err != nil {
		return err
	}

	err = tmCommon.WriteFile(file, contents, 0600)
	if err != nil {
		return err
	}

	return nil
}

// readOrCreatePrivValidator reads or creates the private key file for this config
func readOrCreatePrivValidator(privValFile string) *privval.FilePV {
	// private validator
	var privValidator *privval.FilePV
	if tmCommon.FileExists(privValFile) {
		privValidator = privval.LoadFilePV(privValFile)
	} else {
		privValidator = privval.GenFilePV(privValFile)
		privValidator.Save()
	}
	return privValidator
}

func initializeNodeValidatorFiles(
	config *cfg.Config) (nodeID string, pval *privval.FilePV,
) {
	nodeKey, _ := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	nodeID = string(nodeKey.ID())
	pval = readOrCreatePrivValidator(config.PrivValidatorFile())
	return nodeID, pval
}

func loadGenesisDoc(cdc *amino.Codec, genFile string) (genDoc types.GenesisDoc, err error) {
	genContents, err := ioutil.ReadFile(genFile)
	if err != nil {
		return genDoc, err
	}

	if err := cdc.UnmarshalJSON(genContents, &genDoc); err != nil {
		return genDoc, err
	}

	return genDoc, err
}

func hostnameOrIP(i int) string {
	return fmt.Sprintf("%s%d", viper.GetString(flagNodeHostPrefix), i)
}

func totalValidators() int {
	numValidators := viper.GetInt(flagNumValidators)
	numNonValidators := viper.GetInt(flagNumNonValidators)
	return numNonValidators + numValidators
}

func nodeDir(i int) string {
	outDir := viper.GetString(flagOutputDir)
	nodeDirName := fmt.Sprintf("%s%d", viper.GetString(flagNodeDirPrefix), i)
	nodeDaemonHomeName := viper.GetString(flagNodeDaemonHome)
	return filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
}

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
