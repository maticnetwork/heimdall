package client

import (
	"github.com/cosmos/cosmos-sdk/x/params/client/cli"
	"github.com/cosmos/cosmos-sdk/x/params/client/rest"
	govclient "github.com/maticnetwork/heimdall/gov/client"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
