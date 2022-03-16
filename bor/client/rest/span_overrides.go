package rest

var SPAN_OVERRIDES = map[string][]byte{
	"heimdall-kBBIIu": mainnetSpanJSON,
}

var mainnetSpanJSON = []byte(`[
	{
		"height": "8588755",
		"result": {
			"span_id": 1,
			"start_block": 256,
			"end_block": 6655,
			"validator_set": {
				"validators": [
				  {
					"ID": 1,
					"startEpoch": 0,
					"endEpoch": 0,
					"nonce": 1,
					"power": 1000,
					"pubKey": "0x041e1d526efb1876bb3da8fc03d4181ec06cde0fb5911b49b31b9a7d1a96366233c840de3f5341ae0e5648ee908c4a1f2c51aa9e213f31819329d2c63d0f1a3b37",
					"signer": "0x599d1c2286f18b2218bca637f0f601beda384a6f",
					"last_updated": "",
					"jailed": false,
					"accum": 0
				  }
				],
				"proposer": {
				  "ID": 1,
				  "startEpoch": 0,
				  "endEpoch": 0,
				  "nonce": 1,
				  "power": 1000,
				  "pubKey": "0x041e1d526efb1876bb3da8fc03d4181ec06cde0fb5911b49b31b9a7d1a96366233c840de3f5341ae0e5648ee908c4a1f2c51aa9e213f31819329d2c63d0f1a3b37",
				  "signer": "0x599d1c2286f18b2218bca637f0f601beda384a6f",
				  "last_updated": "",
				  "jailed": false,
				  "accum": 0
				}
			},
			"selected_producers": [
				{
				  "ID": 1,
				  "startEpoch": 0,
				  "endEpoch": 0,
				  "nonce": 1,
				  "power": 1000,
				  "pubKey": "0x041e1d526efb1876bb3da8fc03d4181ec06cde0fb5911b49b31b9a7d1a96366233c840de3f5341ae0e5648ee908c4a1f2c51aa9e213f31819329d2c63d0f1a3b37",
				  "signer": "0x599d1c2286f18b2218bca637f0f601beda384a6f",
				  "last_updated": "",
				  "jailed": false,
				  "accum": 0
				}
			],
			"bor_chain_id": "80001"
		}
	}
]`)
