package cli

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil"
)

func TestParseSubmitProposalFlags(t *testing.T) {
	okJSON, _ := testutil.WriteToNewTempFile(t, `
{
  "title": "Test Proposal",
  "description": "My awesome proposal",
  "type": "Text",
  "deposit": "1000test"
}
`)

	badJSON, _ := testutil.WriteToNewTempFile(t, "bad json")
	fs := NewCmdSubmitProposal().Flags()

	// nonexistent json
	err := fs.Set(FlagProposal, "fileDoesNotExist")
	require.NoError(t, err)
	_, err = parseSubmitProposalFlags(fs)
	require.Error(t, err)

	// invalid json
	err = fs.Set(FlagProposal, badJSON.Name())
	require.NoError(t, err)
	_, err = parseSubmitProposalFlags(fs)
	require.Error(t, err)

	// ok json
	err = fs.Set(FlagProposal, okJSON.Name())
	require.NoError(t, err)
	proposal1, err := parseSubmitProposalFlags(fs)
	require.Nil(t, err, "unexpected error")
	require.Equal(t, "Test Proposal", proposal1.Title)
	require.Equal(t, "My awesome proposal", proposal1.Description)
	require.Equal(t, "Text", proposal1.Type)
	require.Equal(t, "1000test", proposal1.Deposit)

	// flags that can't be used with --proposal
	for _, incompatibleFlag := range ProposalFlags {
		err = fs.Set(incompatibleFlag, "some value")
		require.NoError(t, err)
		_, err := parseSubmitProposalFlags(fs)
		require.Error(t, err)
		err = fs.Set(incompatibleFlag, "")
		require.NoError(t, err)
	}

	// no --proposal, only flags
	err = fs.Set(FlagProposal, "")
	require.NoError(t, err)
	err = fs.Set(FlagTitle, proposal1.Title)
	require.NoError(t, err)
	err = fs.Set(FlagDescription, proposal1.Description)
	require.NoError(t, err)
	err = fs.Set(FlagProposalType, proposal1.Type)
	require.NoError(t, err)
	err = fs.Set(FlagDeposit, proposal1.Deposit)
	require.NoError(t, err)
	proposal2, err := parseSubmitProposalFlags(fs)

	require.Nil(t, err, "unexpected error")
	require.Equal(t, proposal1.Title, proposal2.Title)
	require.Equal(t, proposal1.Description, proposal2.Description)
	require.Equal(t, proposal1.Type, proposal2.Type)
	require.Equal(t, proposal1.Deposit, proposal2.Deposit)

	err = okJSON.Close()
	require.Nil(t, err, "unexpected error")
	err = badJSON.Close()
	require.Nil(t, err, "unexpected error")
}
