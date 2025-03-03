package main

import (
	"context"

	"github.com/maticnetwork/heimdall/cmd/heimdalld/service"
	"github.com/maticnetwork/heimdall/version"
)

func main() {
	version.UpdateHeimdalldInfo()
	service.NewHeimdallService(context.Background(), nil)
}
