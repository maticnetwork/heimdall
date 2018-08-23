# Instructions

I am using glide for project dependency management.
Clone this project in the ~GOPATH/src/github.com directory .

### Build Instructions
sudo go install -tags "netgo ledger" -ldflags "-X github.com/cosmos/cosmos-sdk/version.GitCommit=d5652d96" ./cmd/basecli<br>
sudo go install -tags "netgo ledger" -ldflags "-X github.com/cosmos/cosmos-sdk/version.GitCommit=d5652d96" ./cmd/basecoind<br>

`NOTE : Run these from basecoin directory after installing go properly`

### Steps to test checkpoints
 `NOTE: For now i have reduced checkpoint to 5 transctions so that its easier to debug and read data and stuff`
- basecoind init
- basecli keys add alice --recover ( use the phrase given in the result of above command)
- basecli keys list (shows current keys)
- basecoind start (starts tendermint )
- basecli submitBlock --from="alice" [block_hash] [tx-reciept] [receipt-root] --chain-id=[get from genesis file in ~/.basecoind/config]
Sample command
```
basecli submitBlock --from="alice" 0xa43a8ea012b6e1f4d0ea0f4eb7a2a63b7611dccaf057847b0de41ffc2a735ada 0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421 0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421 --chain-id=test-chain-ZLAVft
```
### Testing if block data has been inserted properly
- basecli tx [tx hash from above command ]
- basecli GetBlock [delegate layer block hash]

### Checkpoint

Insert blocks 5 times to see checkpoint being created .

### In progress

Making rest endpoints for block data submission from delegate layer so that the `submitBlock` command can be removed and we dont have to insert manually
Also thinking of having some tests (lower priority)

### Known Issues while installing

These errors may occur while installing , mostly import errors , glide sucks !
- Tendermint-iavl versioned tree not defined
- go-ethereum/ethereum/ethclient not found

Ping me if this happens

