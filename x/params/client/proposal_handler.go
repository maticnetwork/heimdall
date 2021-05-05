package client

import (
	govclient "github.com/maticnetwork/heimdall/x/gov/client"
	"github.com/maticnetwork/heimdall/x/params/client/cli"
	"github.com/maticnetwork/heimdall/x/params/client/rest"
)

// ProposalHandler is the param change proposal handler.
var ProposalHandler = govclient.NewProposalHandler(cli.NewSubmitParamChangeProposalTxCmd, rest.ProposalRESTHandler)
