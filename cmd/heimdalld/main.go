package main

import (
	"context"

	"github.com/maticnetwork/heimdall/cmd/heimdalld/service"
)

func main() {
	service.NewHeimdallService(context.Background(), nil)
}
