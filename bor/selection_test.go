package bor

import (
	"reflect"
	"testing"

	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"

	"github.com/ethereum/go-ethereum/common"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

const testValidators = `[
  {
    "ID": 3,
    "startEpoch": 0,
    "endEpoch": 0,
    "power": 10000,
    "pubKey": "0x046434e10a34ade13c4fea917346a9fd1473eac2138a0b4e2a36426871918be63188fde4edbf598457592c9a49fe3b0036dd5497079495d132e5045bf499c4bdb1",
    "signer": "0x1c4f0f054a0d6a1415382dc0fd83c6535188b220",
    "last_updated": "0",
    "accum": -40000
  },
  {
    "ID": 4,
    "startEpoch": 0,
    "endEpoch": 0,
    "power": 10000,
    "pubKey": "0x04d9d09f2afc9da3cccc164e8112eb6911a63f5ede10169768f800df83cf99c73f944411e9d4fac3543b11c5f84a82e56b36cfcd34f1d065855c1e2b27af8b5247",
    "signer": "0x461295d3d9249215e758e939a150ab180950720b",
    "last_updated": "0",
    "accum": 10000
  },
  {
    "ID": 5,
    "startEpoch": 0,
    "endEpoch": 0,
    "power": 10000,
    "pubKey": "0x04a36f6ed1f93acb0a38f4cacbe2467c72458ac41ce3b12b34d758205b2bc5d930a4e059462da7a0976c32fce766e1f7e8d73933ae72ac2af231fe161187743932",
    "signer": "0x836fe3e3dd0a5f77d9d5b0f67e48048aaafcd5a0",
    "last_updated": "0",
    "accum": 10000
  },
  {
    "ID": 1,
    "startEpoch": 0,
    "endEpoch": 0,
    "power": 10000,
    "pubKey": "0x04a312814042a6655c8e5ecf0c52cba0b6a6f3291c87cc42260a3c0222410c0d0d59b9139d1c56542e5df0ce2fce3a86ce13e93bd9bde0dc8ff664f8dd5294dead",
    "signer": "0x925a91f8003aaeabea6037103123b93c50b86ca3",
    "last_updated": "0",
    "accum": 10000
  },
  {
    "ID": 2,
    "startEpoch": 0,
    "endEpoch": 0,
    "power": 10000,
    "pubKey": "0x0469536ae98030a7e83ec5ef3baffed2d05a32e31d978e58486f6bdb0fbbf240293838325116090190c0639db03f9cbd8b9aecfd269d016f46e3a2287fbf9ad232",
    "signer": "0xc787af4624cb3e80ee23ae7faac0f2acea2be34c",
    "last_updated": "0",
    "accum": 10000
  }
]`

func TestSelectNextProducers(t *testing.T) {
	t.Parallel()

	type producerSelectionTestCase struct {
		seed            string
		producerCount   uint64
		resultSlots     int64
		resultProducers int64
	}

	testcases := []producerSelectionTestCase{
		{"0x8f5bab218b6bb34476f51ca588e9f4553a3a7ce5e13a66c660a5283e97e9a85a", 10, 5, 5},
		{"0x8f5bab218b6bb34476f51ca588e9f4553a3a7ce5e13a66c660a5283e97e9a85a", 5, 5, 5},
		{"0xe09cc356df20c7a2dd38cb85b680a16ec29bd8b3e1ecc1b20f2e5603d5e7ee85", 10, 5, 5},
		{"0xe09cc356df20c7a2dd38cb85b680a16ec29bd8b3e1ecc1b20f2e5603d5e7ee85", 5, 5, 5},
		{"0x8f5bab218b6bb34476f51ca588e9f4553a3a7ce5e13a66c660a5283e97e9a85a", 4, 4, 3},
		{"0xe09cc356df20c7a2dd38cb85b680a16ec29bd8b3e1ecc1b20f2e5603d5e7ee85", 4, 4, 4},
	}

	var validators []hmTypes.Validator
	err := jsoniter.ConfigFastest.Unmarshal([]byte(testValidators), &validators)
	require.NoError(t, err)
	require.Equal(t, 5, len(validators), "Total validators should be 5")

	for i, testcase := range testcases {
		seed := common.HexToHash(testcase.seed)

		producerIds, err := SelectNextProducers(seed, validators, testcase.producerCount)
		require.NoError(t, err, "Error should be nil")

		producers, slots := getSelectedValidatorsFromIDs(validators, producerIds)
		require.Equal(t, testcase.resultSlots, slots, "Total slots should be %v (Testcase %v)", testcase.resultSlots, i+1)
		require.Equal(t, int(testcase.resultProducers), len(producers), "Total producers should be %v (Testcase %v)", testcase.resultProducers, i+1)
	}
}

func getSelectedValidatorsFromIDs(validators []hmTypes.Validator, producerIds []uint64) ([]hmTypes.Validator, int64) {
	var vals []hmTypes.Validator

	IDToPower := make(map[uint64]uint64)
	for _, ID := range producerIds {
		IDToPower[ID] = IDToPower[ID] + 1
	}

	var slots int64

	for key, value := range IDToPower {
		if val, ok := findValidatorByID(validators, key); ok {
			val.VotingPower = int64(value)
			vals = append(vals, val)
			slots = slots + int64(value)
		}
	}

	return vals, slots
}

func findValidatorByID(validators []hmTypes.Validator, id uint64) (val hmTypes.Validator, ok bool) {
	for _, v := range validators {
		if v.ID.Uint64() == id {
			return v, true
		}
	}

	return
}

func Test_createWeightedRanges(t *testing.T) {
	t.Parallel()

	type args struct {
		vals []uint64
	}

	tests := []struct {
		name        string
		args        args
		ranges      []uint64
		totalWeight uint64
	}{
		{
			args: args{
				vals: []uint64{30, 20, 50, 50, 1},
			},
			ranges:      []uint64{30, 50, 100, 150, 151},
			totalWeight: 151,
		},
		{
			args: args{
				vals: []uint64{1, 2, 1, 2, 1},
			},
			ranges:      []uint64{1, 3, 4, 6, 7},
			totalWeight: 7,
		},
		{
			args: args{
				vals: []uint64{10, 1, 20, 1, 2},
			},
			ranges:      []uint64{10, 11, 31, 32, 34},
			totalWeight: 34,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ranges, totalWeight := createWeightedRanges(tt.args.vals)
			if !reflect.DeepEqual(ranges, tt.ranges) {
				t.Errorf("createWeightedRange() got ranges = %v, want %v", ranges, tt.ranges)
			}

			if totalWeight != tt.totalWeight {
				t.Errorf("createWeightedRange() got totalWeight = %v, want %v", totalWeight, tt.totalWeight)
			}
		})
	}
}

func Test_binarySearch(t *testing.T) {
	t.Parallel()

	type args struct {
		array  []uint64
		search uint64
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			args: args{
				array:  []uint64{},
				search: 0,
			},
			want: -1,
		},
		{
			args: args{
				array:  []uint64{1},
				search: 100,
			},
			want: 0,
		},
		{
			args: args{
				array:  []uint64{1, 1000},
				search: 100,
			},
			want: 1,
		},
		{
			args: args{
				array:  []uint64{1, 100, 1000},
				search: 2,
			},
			want: 1,
		},
		{
			args: args{
				array:  []uint64{1, 100, 1000, 1000},
				search: 1001,
			},
			want: 3,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := binarySearch(tt.args.array, tt.args.search); got != tt.want {
				t.Errorf("binarySearch() = %v, want %v", got, tt.want)
			}
		})
	}
}
