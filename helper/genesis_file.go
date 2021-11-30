package helper

import (
	"embed"
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"
	tmTypes "github.com/tendermint/tendermint/types"
)

//go:embed allocs
var allocs embed.FS

func WriteGenesisFile(chain string, filePath string) (bool, error) {
	switch chain {
	case "mumbai", "mainnet":
		fn := fmt.Sprintf("allocs/%s.json", chain)
		genDoc, err := readPrealloc(fn)
		if err == nil {
			err = genDoc.SaveAs(filePath)
		}
		return err == nil, err
	default:
		return false, nil
	}
}

func readPrealloc(filename string) (result tmTypes.GenesisDoc, err error) {
	f, err := allocs.Open(filename)
	if err != nil {
		err = errors.Errorf("Could not open genesis preallocation for %s: %v", filename, err)
		return
	}
	defer f.Close()
	decoder := json.NewDecoder(f)
	err = decoder.Decode(&result)
	if err != nil {
		err = errors.Errorf("Could not parse genesis preallocation for %s: %v", filename, err)
	}
	return
}
