package cmd

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"

	tmos "github.com/tendermint/tendermint/libs/os"

	"github.com/tendermint/tendermint/crypto"

	"github.com/tendermint/tendermint/p2p"

	"github.com/cosmos/cosmos-sdk/codec"

	bortypes "github.com/maticnetwork/heimdall/x/bor/types"
	stakingtypes "github.com/maticnetwork/heimdall/x/staking/types"
	topuptypes "github.com/maticnetwork/heimdall/x/topup/types"

	tmtime "github.com/tendermint/tendermint/types/time"

	hmCommon "github.com/maticnetwork/heimdall/types/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	srvconfig "github.com/cosmos/cosmos-sdk/server/config"
	"github.com/cosmos/cosmos-sdk/types/module"
	hmTypes "github.com/maticnetwork/heimdall/types"
	cfg "github.com/tendermint/tendermint/config"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	tmcommontempfile "github.com/tendermint/tendermint/libs/tempfile"
	"github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
)

var (
	flagValidatorID       = "id"
	flagStartingIPAddress = "starting-ip-address"
	flagSignerDump        = "signer-dump"
)

func newCreateTestCmd(mbm module.BasicManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-testnet",
		Short: "Initialize files for a heimdall testnet",
		Long: `testnet will create "v" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.

Example:
	heimdalld new-testnet --v 4 --output-dir ./output --starting-ip-address 192.168.10.2
	`,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			outputDir, _ := cmd.Flags().GetString(flagOutputDir)
			keyringBackend, _ := cmd.Flags().GetString(flags.FlagKeyringBackend)
			chainID, _ := cmd.Flags().GetString(flags.FlagChainID)
			nodeDirPrefix, _ := cmd.Flags().GetString(flagNodeDirPrefix)
			nodeDaemonHome, _ := cmd.Flags().GetString(flagNodeDaemonHome)
			startingIPAddress, _ := cmd.Flags().GetString(flagStartingIPAddress)
			numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
			algo, _ := cmd.Flags().GetString(flags.FlagKeyAlgorithm)
			// first validators start ID
			startID, _ := cmd.Flags().GetInt64(flagValidatorID)
			if startID == 0 {
				startID = 1
			}
			return initTestnet(
				clientCtx, cmd, config, mbm, outputDir, chainID,
				nodeDirPrefix, nodeDaemonHome, startingIPAddress, keyringBackend, algo, numValidators, startID,
			)
		},
	}

	cmd.Flags().Int(flagNumValidators, 4, "Number of validators to initialize the testnet with")
	cmd.Flags().Int(flagNumNonValidators, 4, "Number of non validators to initialize the testnet with")
	cmd.Flags().StringP(flagOutputDir, "o", "./mytestnet", "Directory to store initialization data for the testnet")
	cmd.Flags().String(flagNodeDirPrefix, "node", "Prefix the directory name for each node with (node results in node0, node1, ...)")
	cmd.Flags().String(flagNodeHostPrefix, "node", "Hostname prefix (node results in persistent peers list ID0@node0:26656, ID1@node1:26656, ...)")
	cmd.Flags().String(flagNodeDaemonHome, "heimdalld", "Home directory of the node's daemon configuration")
	cmd.Flags().String(flagStartingIPAddress, "192.168.0.1", "Starting IP address (192.168.0.1 results in persistent peers list ID0@192.168.0.1:46656, ID1@192.168.0.2:46656, ...)")
	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Bool(flagSignerDump, true, "dumps all signer information in a json file")
	cmd.Flags().String(flags.FlagKeyringBackend, flags.DefaultKeyringBackend, "Select keyring's backend (os|file|test)")
	cmd.Flags().String(flags.FlagKeyAlgorithm, string(hd.Secp256k1Type), "Key signing algorithm to generate keys for")
	cmd.Flags().Int64(flagValidatorID, 1, "validator starting id")

	return cmd
}

func initTestnet(clientCtx client.Context, cmd *cobra.Command, nodeConfig *cfg.Config, mbm module.BasicManager,
	outputDir string, chainID string, nodeDirPrefix string, nodeDaemonHome string, startingIPAddress string,
	keyringBackend string, algoStr string, numValidators int, startID int64) error {

	nodeConfig.SetRoot(clientCtx.HomeDir)

	if chainID == "" {
		chainID = "chain-" + tmrand.NewRand().Str(6)
	}

	totalNumberOfValidators := totalValidators(cmd)
	nodeIDs := make([]string, totalNumberOfValidators)
	valPubKeys := make([]crypto.PubKey, totalNumberOfValidators)
	valPrivKeys := make([]crypto.PrivKey, totalNumberOfValidators)

	genAccounts := make([]authtypes.GenesisAccount, totalNumberOfValidators)
	genFiles := make([]string, totalNumberOfValidators)
	signers := make([]ValidatorAccountFormatter, totalNumberOfValidators)

	validators := make([]*hmTypes.Validator, numValidators)
	dividendAccounts := make([]*hmTypes.DividendAccount, numValidators)

	simappConfig := srvconfig.DefaultConfig()
	simappConfig.API.Enable = true
	simappConfig.API.Swagger = true
	simappConfig.Telemetry.Enabled = true
	simappConfig.Telemetry.PrometheusRetentionTime = 60
	simappConfig.Telemetry.EnableHostnameLabel = false
	simappConfig.Telemetry.GlobalLabels = [][]string{{"chain_id", chainID}}

	inBuf := bufio.NewReader(cmd.InOrStdin())

	// generate private keys, node IDs, and initial transactions
	for i := 0; i < totalNumberOfValidators; i++ {
		nodeDirName := fmt.Sprintf("%s%d", nodeDirPrefix, i)
		nodeDir := filepath.Join(outputDir, nodeDirName, nodeDaemonHome)
		nodeConfig.SetRoot(nodeDir)

		if err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm); err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		nodeConfig.Moniker = nodeDirName

		// keyring
		kb, err := keyring.New(sdk.KeyringServiceName(), keyringBackend, nodeDir, inBuf)
		if err != nil {
			return err
		}

		keyringAlgos, _ := kb.SupportedAlgorithms()
		algo, err := keyring.NewSigningAlgoFromString(algoStr, keyringAlgos)
		if err != nil {
			return err
		}

		sdkAddress, secret, err := server.GenerateSaveCoinKey(kb, nodeDirName, true, algo)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		fmt.Println("sdkAddress ", sdkAddress)
		info := map[string]string{"secret": secret}

		cliPrint, err := json.Marshal(info)
		if err != nil {
			return err
		}

		// save private key seed words
		if err := writeFile(fmt.Sprintf("%v.json", "key_seed"), nodeDir, cliPrint); err != nil {
			return err
		}

		var tmValPubKey cryptotypes.PubKey
		nodeIDs[i], valPubKeys[i], valPrivKeys[i], err = InitializeNodeValidatorFiles(nodeConfig, secret)
		if err != nil {
			_ = os.RemoveAll(outputDir)
			return err
		}

		tmValPubKey, err = cryptocodec.FromTmPubKeyInterface(valPubKeys[i])
		if err != nil {
			return err
		}

		fmt.Println("tender val pub addr ", valPubKeys[i].Address().String())
		fmt.Println("cosmos val pub addr ", tmValPubKey.Address().String())

		genFiles[i] = nodeConfig.GenesisFile()

		// get the key info
		newPubkey := hmCommon.NewPubKey(valPubKeys[i].Bytes())

		genAccounts[i] = authtypes.NewBaseAccount(sdkAddress, nil, 0, 0)

		if i < numValidators {
			validators[i] = hmTypes.NewValidator(
				hmTypes.NewValidatorID(uint64(startID+int64(i))),
				0,
				0,
				1,
				10000,
				newPubkey,
				sdkAddress,
			)
			signer := sdkAddress
			// create dividend account for validator
			dividendAcc := hmTypes.NewDividendAccount(signer, ZeroIntString)
			dividendAccounts[i] = &dividendAcc
		}

		if err != nil {
			return err
		}

		signers[i] = newGetSignerInfo(newPubkey, valPrivKeys[i].Bytes())
		srvconfig.WriteConfigFile(filepath.Join(nodeDir, "config/app.toml"), simappConfig)
		WriteDefaultHeimdallConfig(filepath.Join(nodeDir, "config/heimdall-config.toml"))
	}

	if err := initGenFiles(clientCtx, mbm, chainID, genAccounts, genFiles, totalNumberOfValidators, validators, dividendAccounts); err != nil {
		return err
	}

	newpopulatePersistentPeersInConfigAndWriteIt(nodeConfig, cmd, totalNumberOfValidators)

	// dump signer information in a json file
	dump, _ := cmd.Flags().GetBool(flagSignerDump)
	if dump {
		signerJSON, err := json.MarshalIndent(signers, "", "  ")
		if err != nil {
			return err
		}

		if err := tmcommontempfile.WriteFileAtomic(filepath.Join(outputDir, "signer-dump.json"), signerJSON, 0600); err != nil {
			fmt.Println("Error writing writing singers info into signer-dump file ", err)
			return err
		}
	}

	cmd.PrintErrf("Successfully initialized %d node directories\n", totalNumberOfValidators)
	return nil
}

func initGenFiles(clientCtx client.Context, mbm module.BasicManager, chainID string,
	genAccounts []authtypes.GenesisAccount, genFiles []string, numValidators int,
	validators []*hmTypes.Validator,
	dividendAccounts []*hmTypes.DividendAccount) error {

	validatorSet := hmTypes.NewValidatorSet(validators)
	depCdc := clientCtx.JSONMarshaler
	cdc := depCdc.(codec.Marshaler)

	appGenState := mbm.DefaultGenesis(clientCtx.JSONMarshaler)

	// set the accounts in the genesis state
	var authGenState authtypes.GenesisState
	clientCtx.JSONMarshaler.MustUnmarshalJSON(appGenState[authtypes.ModuleName], &authGenState)

	accounts, err := authtypes.PackAccounts(genAccounts)
	if err != nil {
		return err
	}

	authGenState.Accounts = accounts
	appGenState[authtypes.ModuleName] = clientCtx.JSONMarshaler.MustMarshalJSON(&authGenState)

	// staking state change
	appGenState, err = stakingtypes.SetGenesisStateToAppState(cdc, appGenState, validators, validatorSet)
	if err != nil {
		return err
	}
	// bor state change
	appGenState, err = bortypes.SetGenesisStateToAppState(appGenState, *validatorSet)
	if err != nil {
		return err
	}

	// topup state change
	appGenState, err = topuptypes.SetGenesisStateToAppState(appGenState, dividendAccounts)
	if err != nil {
		return err
	}

	appGenStateJSON, err := json.MarshalIndent(appGenState, "", "  ")
	if err != nil {
		return err
	}

	genDoc := types.GenesisDoc{
		ChainID:    chainID,
		AppState:   appGenStateJSON,
		Validators: nil,
	}

	if genDoc.GenesisTime.IsZero() {
		genDoc.GenesisTime = tmtime.Now()
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}

	// generate empty genesis files for each validator and save
	for i := 0; i < numValidators; i++ {
		if err := genDoc.SaveAs(genFiles[i]); err != nil {
			return err
		}
	}

	return nil
}

// populate persistent peers in config
func newpopulatePersistentPeersInConfigAndWriteIt(config *cfg.Config, cmd *cobra.Command, totalValidators int) {
	persistentPeers := make([]string, totalValidators)
	for i := 0; i < totalValidators; i++ {
		config.SetRoot(nodeDir(i, cmd))
		nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
		if err != nil {
			return
		}
		persistentPeers[i] = p2p.IDAddressString(nodeKey.ID(), fmt.Sprintf("%s:%d", ahostnameOrIP(i, cmd), 26656))
	}

	persistentPeersList := strings.Join(persistentPeers, ",")
	persistentPeersList = ""
	for i := 0; i < totalValidators; i++ {
		config.Moniker = ahostnameOrIP(i, cmd)
		config.SetRoot(nodeDir(i, cmd))
		config.P2P.PersistentPeers = persistentPeersList
		config.P2P.AddrBookStrict = true

		// overwrite default config
		cfg.WriteConfigFile(filepath.Join(nodeDir(i, cmd), "config", "config.toml"), config)
	}
}

// hostname of ip of nodes
func ahostnameOrIP(i int, cmd *cobra.Command) string {
	hOrIP, _ := cmd.Flags().GetString(flagNodeHostPrefix)
	return fmt.Sprintf("%s%d", hOrIP, i)
}

// GetSignerInfo returns signer information
func newGetSignerInfo(pub hmCommon.PubKey, priv []byte) ValidatorAccountFormatter {
	return ValidatorAccountFormatter{
		Address: hmCommon.AccAddressToHeimdallAddress(pub.Address().Bytes()).String(),
		PubKey:  pub.String(),
		PrivKey: "0x" + hex.EncodeToString(priv),
	}
}

func writeFile(name string, dir string, contents []byte) error {
	writePath := filepath.Join(dir)
	file := filepath.Join(writePath, name)

	err := tmos.EnsureDir(writePath, 0755)
	if err != nil {
		return err
	}

	err = tmos.WriteFile(file, contents, 0644)
	if err != nil {
		return err
	}

	return nil
}
