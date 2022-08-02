module github.com/maticnetwork/heimdall

go 1.16

require (
	github.com/RichardKnop/machinery v1.7.4
	github.com/allegro/bigcache v1.2.1 // indirect
	github.com/aristanetworks/goarista v0.0.0-20220727195344-9ed9ecb0092c // indirect
	github.com/cbergoon/merkletree v0.2.0
	github.com/cespare/cp v1.1.1 // indirect
	github.com/cosmos/cosmos-sdk v0.37.4
	github.com/elastic/gosigar v0.10.5 // indirect
	github.com/ethereum/go-ethereum v1.10.4
	github.com/gballet/go-libpcsclite v0.0.0-20191108122812-4678299bea08 // indirect
	github.com/go-kit/kit v0.10.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/huin/goupnp v1.0.3 // indirect
	github.com/json-iterator/go v1.1.12
	github.com/maticnetwork/bor v0.0.0-20191204165821-bd9cd503a1b3
	github.com/pborman/uuid v1.2.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1
	github.com/rakyll/statik v0.1.6
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/status-im/keycard-go v0.0.0-20211109104530-b0e0482ba91d // indirect
	github.com/steakknife/bloomfilter v0.0.0-20180922174646-6819c0d2a570 // indirect
	github.com/steakknife/hamming v0.0.0-20180906055917-c99c65617cd3 // indirect
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	github.com/stretchr/testify v1.7.1
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.32.7
	github.com/tendermint/tm-db v0.2.0
	github.com/tyler-smith/go-bip39 v1.1.0 // indirect
	github.com/wsddn/go-ecdh v0.0.0-20161211032359-48726bab9208 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/tendermint/tendermint => github.com/maticnetwork/tendermint v0.26.0-dev0.0.20220311091200-3ea059b213ed

replace github.com/cosmos/cosmos-sdk => github.com/maticnetwork/cosmos-sdk v0.37.5-0.20220311095845-81690c6a53e7

replace github.com/ethereum/go-ethereum => github.com/maticnetwork/bor v0.2.16
