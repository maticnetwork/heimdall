// The following directive is necessary to make the package coherent:

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
`))

var tomlConfig struct {
	NewSelectionAlgoHeight int `toml:"new_selection_algo_height"`
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
		BlockHeight int
	}{
		BlockHeight: tomlConfig.NewSelectionAlgoHeight,
	})
}

func chekcError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
