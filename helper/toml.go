package helper

import (
	"bytes"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/viper"
	cmn "github.com/tendermint/tendermint/libs/common"
)

// Note: any changes to the comments/variables/mapstructure
// must be reflected in the appropriate struct in helper/config.go
const defaultConfigTemplate = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

##### RPC and REST configs #####

# RPC endpoint for ethereum chain
eth_rpc_url = "{{ .EthRPCUrl }}"

# RPC endpoint for bor chain
bor_rpc_url = "{{ .BorRPCUrl }}"

# RPC endpoint for tendermint
tendermint_rpc_url = "{{ .TendermintRPCUrl }}"

# Heimdall REST server endpoint
heimdall_rest_server = "{{ .HeimdallServerURL }}"

#### Bridge configs ####

# AMQP endpoint
amqp_url = "{{ .AmqpURL }}"

## Poll intervals
checkpoint_poll_interval = "{{ .CheckpointerPollInterval }}"
syncer_poll_interval = "{{ .SyncerPollInterval }}"
noack_poll_interval = "{{ .NoACKPollInterval }}"
clerk_poll_interval = "{{ .ClerkPollInterval }}"
span_poll_interval = "{{ .SpanPollInterval }}"

#### gas limits ####
main_chain_gas_limit = "{{ .MainchainGasLimit }}"

#### gas price ####
main_chain_max_gas_price = "{{ .MainchainMaxGasPrice }}"

##### Timeout Config #####
no_ack_wait_time = "{{ .NoACKWaitTime }}"

##### current network chain - newSelectionAlgoHeight depends on this #####
network_chain = "{{ .NetworkChain }}"
`

var configTemplate *template.Template

func init() {
	var err error
	tmpl := template.New("appConfigFileTemplate")
	if configTemplate, err = tmpl.Parse(defaultConfigTemplate); err != nil {
		panic(err)
	}
}

// ParseConfig retrieves the default environment configuration for the
// application.
func ParseConfig() (*Configuration, error) {
	conf := GetDefaultHeimdallConfig()
	err := viper.Unmarshal(conf)
	return &conf, err
}

// WriteConfigFile renders config using the template and writes it to
// configFilePath.
func WriteConfigFile(configFilePath string, config *Configuration) {
	var buffer bytes.Buffer

	if err := configTemplate.Execute(&buffer, config); err != nil {
		panic(err)
	}

	cmn.MustWriteFile(configFilePath, buffer.Bytes(), 0644)
}

func readNetworkChainToml(configDir string, network string) int64 {
	if strings.Compare(network, "local") == 0 {
		return 0
	}
	v := viper.New()
	v.SetConfigType("toml")
	v.SetConfigName(network)
	v.AddConfigPath(configDir)
	if err := v.ReadInConfig(); err != nil {
		log.Fatal(err)
	}

	result := v.GetInt64("new_selection_algo_height")
	return result
}

func CopyNetworkChainTomlsToConfigDir(configDir string, networkChains []string) (err error) {
	var input []byte
	for _, chain := range networkChains {
		if strings.Compare(chain, "local") == 0 {
			continue
		}

		fileName := chain + ".toml"
		sourceDir := filepath.Join(".", fileName)
		destDir := filepath.Join(configDir, fileName)
		if input, err = ioutil.ReadFile(sourceDir); err != nil {
			return
		}
		if err = ioutil.WriteFile(destDir, input, 0644); err != nil {
			return
		}
	}
	return
}
