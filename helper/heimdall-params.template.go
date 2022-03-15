// The following directive is necessary to make the package coherent:

//go:build ignore
// +build ignore

// This program generate heimdall-params.go, It must be invoked from make file that purpose only.

package main

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/BurntSushi/toml"
)

var packageTemplate = template.Must(template.New("").Parse(`package helper

const NewSelectionAlgoHeight = {{ .BlockHeight }}
const SpanOverrideBlockHeight = {{ .SpanOverrideBlockHeight }}
`))

var tomlConfig struct {
	NewSelectionAlgoHeight  int `toml:"new_selection_algo_height"`
	SpanOverrideBlockHeight int `toml:"span_override_height"`
}

var networks = []string{
	"mainnet",
	"mumbai",
	"local",
}

func main() {
	var network = networks[0]
	if len(os.Args) > 1 {
		networkFile := os.Args[1]
		for _, n := range networks {
			if n == networkFile {
				network = networkFile
			}
		}
	}

	filePath := fmt.Sprintf("%s.toml", network)
	toml.DecodeFile(filePath, &tomlConfig)

	f, err := os.Create("helper/heimdall-params.go")
	chekcError(err)
	defer f.Close()

	packageTemplate.Execute(f, struct {
		BlockHeight             int
		SpanOverrideBlockHeight int
	}{
		BlockHeight:             tomlConfig.NewSelectionAlgoHeight,
		SpanOverrideBlockHeight: tomlConfig.SpanOverrideBlockHeight,
	})
}

func chekcError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
