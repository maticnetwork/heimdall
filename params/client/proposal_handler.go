package client

import (
	govclient "github.com/maticnetwork/heimdall/gov/client"
	"github.com/maticnetwork/heimdall/params/client/cli"
	"github.com/maticnetwork/heimdall/params/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
