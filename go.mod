module github.com/maticnetwork/heimdall

go 1.16

require (
	github.com/RichardKnop/machinery v1.7.4
	github.com/cbergoon/merkletree v0.2.0
	github.com/cosmos/cosmos-sdk v0.37.4
	github.com/ethereum/go-ethereum v1.10.4
	github.com/go-kit/kit v0.10.0
	github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp v0.5.7 // indirect
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.8.0
	github.com/hashicorp/golang-lru v0.5.5-0.20210104140557-80c98217689d
	github.com/huin/goupnp v1.0.3 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prysmaticlabs/prysm v1.4.4
	github.com/rakyll/statik v0.1.6
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/status-im/keycard-go v0.0.0-20211109104530-b0e0482ba91d // indirect
	github.com/streadway/amqp v0.0.0-20200108173154-1c71cc93ed71
	github.com/stretchr/testify v1.7.1
	github.com/syndtr/goleveldb v1.0.1-0.20210819022825-2ae1ddf74ef7
	github.com/tendermint/crypto v0.0.0-20191022145703-50d29ede1e15
	github.com/tendermint/go-amino v0.15.0
	github.com/tendermint/tendermint v0.32.7
	github.com/tendermint/tm-db v0.2.0
	github.com/xdg/scram v1.0.3 // indirect
	github.com/xdg/stringprep v1.0.3 // indirect
	golang.org/x/net v0.0.0-20210917221730-978cfadd31cf // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20220615213510-4f61da869c0c // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	golang.org/x/tools v0.1.6 // indirect
	google.golang.org/genproto v0.0.0-20210921142501-181ce0d877f6 // indirect
	google.golang.org/grpc v1.48.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

replace github.com/tendermint/tendermint => github.com/maticnetwork/tendermint v0.26.0-dev0.0.20220803122020-9c20f6000900

replace github.com/cosmos/cosmos-sdk => github.com/maticnetwork/cosmos-sdk v0.37.5-0.20220311095845-81690c6a53e7

replace github.com/ethereum/go-ethereum => github.com/maticnetwork/bor v0.2.16
