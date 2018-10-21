# Heimdall

### Installation

Build the main files in cmd/basecli and cmd/basecoind
`cd cmd/basecli; go build; cd ../basecoind; go build;cd ../../`

Run `./basecoind start` in one terminal and `./basecli rest-server` in another

Send checkpoint via a POST request to `http://localhost:1317/checkpoint/submitCheckpoint` with following data fields

```json
{
  "Root_hash": "0xd494377d4439a844214b565e1c211ea7154ca300b98e3c296f19fc9ada36db33",
  "Start_block": 4733031,
  "End_block": 4733034,
  "Proposer_address": "0x84f8a67E4d16bb05aBCa3d154091566921e0B5e9"
}
```

Your transaction to kovan contract should go through if you have kovan-ether in address given in proposer field.
