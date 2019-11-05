package helper

// Note: any changes to the comments/variables/mapstructure
// must be reflected in the appropriate struct in helper/config.go
const defaultConfigTemplate = `# This is a TOML config file.
# For more information, see https://github.com/toml-lang/toml

##### main base config options #####

# RPC endpoint for ethereum chain
eth_RPC_URL = "{{ .BaseConfig.ProxyApp }}"

# RPC endpoint for matic chain
bor_RPC_URL = "{{ .BaseConfig.Moniker }}"

# AMQP URL 
amqp_URL = "{{ . }}"
`
