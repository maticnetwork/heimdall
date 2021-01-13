package cli

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"

	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/gov/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group gov queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetCmdQueryProposal(),
		GetCmdQueryProposals(),
		GetCmdQueryVote(),
		GetCmdQueryVotes(),
		GetCmdQueryParam(),
		GetCmdQueryParams(),
		GetCmdQueryProposer(),
		GetCmdQueryDeposit(),
		GetCmdQueryDeposits(),
		GetCmdQueryTally(),
	)

	return cmd
}

// GetCmdQueryProposal implements the query proposal command.
func GetCmdQueryProposal() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposal [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query details of a single proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a proposal. You can find the
proposal-id by running "%s query gov proposals".
Example:
$ %s query gov proposal 1
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid uint, please input a valid proposal-id", args[0])
			}

			// Query the proposal
			res, err := queryClient.Proposal(
				context.Background(),
				&types.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(&res.Proposal)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryProposals implements a query proposals command. Command to Get a
// Proposal Information.
func GetCmdQueryProposals() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposals",
		Short: "Query proposals with optional filters",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query for a all paginated proposals that match optional filters:
Example:
$ %s query gov proposals --depositor 2
$ %s query gov proposals --voter 2
$ %s query gov proposals --status (DepositPeriod|VotingPeriod|Passed|Rejected)
`,
				version.AppName, version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			depositorID, _ := cmd.Flags().GetUint64(flagDepositor)
			voterID, _ := cmd.Flags().GetUint64(flagVoter)
			strProposalStatus, _ := cmd.Flags().GetString(flagStatus)
			numLimit, _ := cmd.Flags().GetUint64(flagNumLimit)

			var proposalStatus types.ProposalStatus
			var err error

			if len(strProposalStatus) != 0 {
				proposalStatus, err = types.ProposalStatusFromString(NormalizeProposalStatus(strProposalStatus))
				if err != nil {
					return err
				}
			}

			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err = client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Proposals(
				context.Background(),
				&types.QueryProposalsRequest{
					ProposalStatus: proposalStatus,
					Voter:          hmTypes.NewValidatorID(voterID),
					Depositor:      hmTypes.NewValidatorID(depositorID),
					NumLimit:       numLimit,
				},
			)
			if err != nil {
				return err
			}

			if len(res.GetProposals()) == 0 {
				return fmt.Errorf("no proposals found")
			}

			return clientCtx.PrintOutput(res)
		},
	}

	cmd.Flags().Uint64(flagDepositor, 0, "(optional) filter by proposals deposited on by depositor")
	cmd.Flags().Uint64(flagVoter, 0, "(optional) filter by proposals voted on by voted")
	cmd.Flags().String(flagStatus, "", "(optional) filter proposals by proposal status, status: deposit_period/voting_period/passed/rejected")
	cmd.Flags().Uint64(flagNumLimit, 0, "(optional) limit to latest [number] proposals. Defaults to all proposals")
	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryVote implements the query proposal vote command. Command to Get a
// Proposal Information.
func GetCmdQueryVote() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vote [proposal-id] [voter-id]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a single vote",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a single vote on a proposal given its identifier.
Example:
$ %s query gov vote 1 3
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			_, err = queryClient.Proposal(
				context.Background(),
				&types.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			// validate that the proposal id is a uint
			voterID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			res, err := queryClient.Vote(
				context.Background(),
				&types.QueryVoteRequest{ProposalId: proposalID, Voter: hmTypes.NewValidatorID(voterID)},
			)
			if err != nil {
				return err
			}

			vote := res.GetVote()
			if vote.Empty() {
				resByTxQuery, err := QueryVoteByTxQuery(
					clientCtx,
					types.QueryVoteRequest{ProposalId: proposalID, Voter: hmTypes.NewValidatorID(voterID)},
				)

				if err != nil {
					return err
				}

				if err := clientCtx.JSONMarshaler.UnmarshalJSON(resByTxQuery, &vote); err != nil {
					return err
				}
			}

			return clientCtx.PrintOutput(&res.Vote)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryVotes implements the command to query for proposal votes.
func GetCmdQueryVotes() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "votes [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query votes on a proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query vote details for a single proposal by its identifier.
Example:
$ %[1]s query gov votes 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			proposalRes, err := queryClient.Proposal(
				context.Background(),
				&types.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			proposalStatus := proposalRes.GetProposal().Status
			if !(proposalStatus == types.StatusVotingPeriod || proposalStatus == types.StatusDepositPeriod) {
				resByTxQuery, err := QueryVotesByTxQuery(clientCtx, types.QueryProposalRequest{ProposalId: proposalID})
				if err != nil {
					return err
				}

				var votes types.Votes
				clientCtx.JSONMarshaler.MustUnmarshalJSON(resByTxQuery, &votes)
				return clientCtx.PrintOutput(votes)

			}

			res, err := queryClient.Votes(
				context.Background(),
				&types.QueryVotesRequest{ProposalId: proposalID},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)

		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryDeposit implements the query proposal deposit command. Command to
// get a specific Deposit Information
func GetCmdQueryDeposit() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposit [proposal-id] [depositer-addr]",
		Args:  cobra.ExactArgs(2),
		Short: "Query details of a deposit",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for a single proposal deposit on a proposal by its identifier.
Example:
$ %s query gov deposit 1 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid uint, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			_, err = queryClient.Proposal(
				context.Background(),
				&types.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			depositorID, _ := cmd.Flags().GetUint64(args[1])

			res, err := queryClient.Deposit(
				context.Background(),
				&types.QueryDepositRequest{ProposalId: proposalID, Depositor: hmTypes.NewValidatorID(depositorID)},
			)
			if err != nil {
				return err
			}

			deposit := res.GetDeposit()
			if deposit.Empty() {
				resByTxQuery, err := QueryDepositByTxQuery(
					clientCtx,
					types.QueryDepositRequest{ProposalId: proposalID, Depositor: hmTypes.NewValidatorID(depositorID)},
				)
				if err != nil {
					return err
				}
				clientCtx.JSONMarshaler.MustUnmarshalJSON(resByTxQuery, &deposit)
			}

			return clientCtx.PrintOutput(&deposit)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryDeposits implements the command to query for proposal deposits.
func GetCmdQueryDeposits() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deposits [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query deposits on a proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query details for all deposits on a proposal.
You can find the proposal-id by running "%s query gov proposals".
Example:
$ %s query gov deposits 1
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid uint, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			proposalRes, err := queryClient.Proposal(
				context.Background(),
				&types.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			proposalStatus := proposalRes.GetProposal().Status
			if !(proposalStatus == types.StatusVotingPeriod || proposalStatus == types.StatusDepositPeriod) {
				resByTxQuery, err := QueryDepositsByTxQuery(clientCtx, types.QueryProposalRequest{ProposalId: proposalID})
				if err != nil {
					return err
				}

				var dep types.Deposits
				clientCtx.JSONMarshaler.MustUnmarshalJSON(resByTxQuery, &dep)
				return clientCtx.PrintOutput(&dep)
			}

			res, err := queryClient.Deposits(
				context.Background(),
				&types.QueryDepositsRequest{ProposalId: proposalID},
			)

			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryTally implements the command to query for proposal tally result.
func GetCmdQueryTally() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tally [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Get the tally of a proposal vote",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query tally of votes on a proposal. You can find
the proposal-id by running "%s query gov proposals".
Example:
$ %s query gov tally 1
`,
				version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// validate that the proposal id is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s not a valid int, please input a valid proposal-id", args[0])
			}

			// check to see if the proposal is in the store
			_, err = queryClient.Proposal(
				context.Background(),
				&types.QueryProposalRequest{ProposalId: proposalID},
			)
			if err != nil {
				return fmt.Errorf("failed to fetch proposal-id %d: %s", proposalID, err)
			}

			// Query store
			res, err := queryClient.TallyResult(
				context.Background(),
				&types.QueryTallyResultRequest{ProposalId: proposalID},
			)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(&res.Tally)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryParams implements the query params command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Short: "Query the parameters of the governance process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the all the parameters for the governance process.
Example:
$ %s query gov params
`,
				version.AppName,
			),
		),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// Query store for all 3 params
			votingRes, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{ParamsType: "voting"},
			)
			if err != nil {
				return err
			}

			tallyRes, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{ParamsType: "tallying"},
			)
			if err != nil {
				return err
			}

			depositRes, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{ParamsType: "deposit"},
			)
			if err != nil {
				return err
			}

			params := types.NewParams(
				votingRes.GetVotingParams(),
				tallyRes.GetTallyParams(),
				depositRes.GetDepositParams(),
			)

			return clientCtx.PrintOutput(&params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryParam implements the query param command.
func GetCmdQueryParam() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "param [param-type]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the parameters (voting|tallying|deposit) of the governance process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the all the parameters for the governance process.
Example:
$ %s query gov param voting
$ %s query gov param tallying
$ %s query gov param deposit
`,
				version.AppName, version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)

			// Query store
			res, err := queryClient.Params(
				context.Background(),
				&types.QueryParamsRequest{ParamsType: args[0]},
			)
			if err != nil {
				return err
			}

			// var out fmt.Stringer
			switch args[0] {
			case "voting":
				return clientCtx.PrintOutput(&res.VotingParams)
			case "tallying":
				return clientCtx.PrintOutput(&res.TallyParams)
			case "deposit":
				return clientCtx.PrintOutput(&res.DepositParams)
			default:
				return fmt.Errorf("argument must be one of (voting|tallying|deposit), was %s", args[0])
			}
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryProposer implements the query proposer command.
func GetCmdQueryProposer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposer [proposal-id]",
		Args:  cobra.ExactArgs(1),
		Short: "Query the proposer of a governance proposal",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query which address proposed a proposal with a given ID.
Example:
$ %s query gov proposer 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// validate that the proposalID is a uint
			proposalID, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return fmt.Errorf("proposal-id %s is not a valid uint", args[0])
			}

			prop, err := QueryProposerByTxQuery(clientCtx, proposalID)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(&prop)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

const (
	defaultPage  = 1
	defaultLimit = 30 // should be consistent with tendermint/tendermint/rpc/core/pipe.go:19
)

// // Proposer contains metadata of a governance proposal used for querying a
// // proposer.
// type Proposer struct {
// 	ProposalID uint64 `json:"proposal_id" yaml:"proposal_id"`
// 	Proposer   string `json:"proposer" yaml:"proposer"`
// }

// // NewProposer returns a new Proposer given id and proposer
// func NewProposer(proposalID uint64, proposer string) Proposer {
// 	return Proposer{proposalID, proposer}
// }

// func (p Proposer) String() string {
// 	return fmt.Sprintf("Proposal with ID %d was proposed by %s", p.ProposalID, p.Proposer)
// }

// func (*Proposer) ProtoMessage() {}
// func (m *Proposer) Reset()      { *m = Proposer{} }

// QueryDepositsByTxQuery will query for deposits via a direct txs tags query. It
// will fetch and build deposits directly from the returned txs and return a
// JSON marshalled result or any error that occurred.
//
// NOTE: SearchTxs is used to facilitate the txs query which does not currently
// support configurable pagination.
func QueryDepositsByTxQuery(clientCtx client.Context, params types.QueryProposalRequest) ([]byte, error) {
	events := []string{
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, types.TypeMsgDeposit),
		fmt.Sprintf("%s.%s='%s'", types.EventTypeProposalDeposit, types.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalId))),
	}

	// NOTE: SearchTxs is used to facilitate the txs query which does not currently
	// support configurable pagination.
	searchResult, err := authclient.QueryTxsByEvents(clientCtx, events, defaultPage, defaultLimit, "")
	if err != nil {
		return nil, err
	}

	var deposits types.Deposits

	for _, info := range searchResult.Txs {
		for _, msg := range info.GetTx().GetMsgs() {
			if msg.Type() == types.TypeMsgDeposit {
				depMsg := msg.(*types.MsgDeposit)

				deposits = append(deposits, types.Deposit{
					Depositor:  depMsg.Validator,
					ProposalId: params.ProposalId,
					Amount:     depMsg.Amount,
				})
			}
		}
	}

	bz, err := clientCtx.JSONMarshaler.MarshalJSON(deposits)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// QueryVotesByTxQuery will query for votes via a direct txs tags query. It
// will fetch and build votes directly from the returned txs and return a JSON
// marshalled result or any error that occurred.
func QueryVotesByTxQuery(clientCtx client.Context, params types.QueryProposalRequest) ([]byte, error) {
	events := []string{
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, types.TypeMsgVote),
		fmt.Sprintf("%s.%s='%s'", types.EventTypeProposalVote, types.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalId))),
	}

	// NOTE: SearchTxs is used to facilitate the txs query which does not currently
	// support configurable pagination.
	searchResult, err := authclient.QueryTxsByEvents(clientCtx, events, defaultPage, defaultLimit, "")
	if err != nil {
		return nil, err
	}

	var votes types.Votes

	for _, info := range searchResult.Txs {
		for _, msg := range info.GetTx().GetMsgs() {
			if msg.Type() == types.TypeMsgVote {
				voteMsg := msg.(*types.MsgVote)

				votes = append(votes, types.Vote{
					Voter:      voteMsg.Validator,
					ProposalId: params.ProposalId,
					Option:     voteMsg.Option,
				})
			}
		}
	}

	bz, err := clientCtx.JSONMarshaler.MarshalJSON(votes)
	if err != nil {
		return nil, err
	}

	return bz, nil
}

// QueryVoteByTxQuery will query for a single vote via a direct txs tags query.
func QueryVoteByTxQuery(clientCtx client.Context, params types.QueryVoteRequest) ([]byte, error) {
	events := []string{
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, types.TypeMsgVote),
		fmt.Sprintf("%s.%s='%s'", types.EventTypeProposalVote, types.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalId))),
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeySender, []byte(params.Voter.String())),
	}

	// NOTE: SearchTxs is used to facilitate the txs query which does not currently
	// support configurable pagination.
	searchResult, err := authclient.QueryTxsByEvents(clientCtx, events, defaultPage, defaultLimit, "")
	if err != nil {
		return nil, err
	}
	for _, info := range searchResult.Txs {
		for _, msg := range info.GetTx().GetMsgs() {
			// there should only be a single vote under the given conditions
			if msg.Type() == types.TypeMsgVote {
				voteMsg := msg.(*types.MsgVote)

				vote := types.Vote{
					Voter:      voteMsg.Validator,
					ProposalId: params.ProposalId,
					Option:     voteMsg.Option,
				}

				bz, err := clientCtx.JSONMarshaler.MarshalJSON(&vote)
				if err != nil {
					return nil, err
				}

				return bz, nil
			}
		}
	}

	return nil, fmt.Errorf("address '%s' did not vote on proposalID %d", params.Voter, params.ProposalId)
}

// QueryDepositByTxQuery will query for a single deposit via a direct txs tags
// query.
func QueryDepositByTxQuery(clientCtx client.Context, params types.QueryDepositRequest) ([]byte, error) {
	events := []string{
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, types.TypeMsgDeposit),
		fmt.Sprintf("%s.%s='%s'", types.EventTypeProposalDeposit, types.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", params.ProposalId))),
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeySender, params.Depositor),
	}

	// NOTE: SearchTxs is used to facilitate the txs query which does not currently
	// support configurable pagination.
	searchResult, err := authclient.QueryTxsByEvents(clientCtx, events, defaultPage, defaultLimit, "")
	if err != nil {
		return nil, err
	}

	for _, info := range searchResult.Txs {
		for _, msg := range info.GetTx().GetMsgs() {
			// there should only be a single deposit under the given conditions
			if msg.Type() == types.TypeMsgDeposit {
				depMsg := msg.(*types.MsgDeposit)

				deposit := types.Deposit{
					Depositor:  depMsg.Validator,
					ProposalId: params.ProposalId,
					Amount:     depMsg.Amount,
				}

				bz, err := clientCtx.JSONMarshaler.MarshalJSON(&deposit)
				if err != nil {
					return nil, err
				}

				return bz, nil
			}
		}
	}

	return nil, fmt.Errorf("address '%s' did not deposit to proposalID %d", params.Depositor, params.ProposalId)
}

// QueryProposerByTxQuery will query for a proposer of a governance proposal by
// ID.
func QueryProposerByTxQuery(clientCtx client.Context, proposalID uint64) (types.Proposer, error) {
	events := []string{
		fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, types.TypeMsgSubmitProposal),
		fmt.Sprintf("%s.%s='%s'", types.EventTypeSubmitProposal, types.AttributeKeyProposalID, []byte(fmt.Sprintf("%d", proposalID))),
	}

	// NOTE: SearchTxs is used to facilitate the txs query which does not currently
	// support configurable pagination.
	searchResult, err := authclient.QueryTxsByEvents(clientCtx, events, defaultPage, defaultLimit, "")
	if err != nil {
		return types.Proposer{}, err
	}

	for _, info := range searchResult.Txs {
		for _, msg := range info.GetTx().GetMsgs() {
			// there should only be a single proposal under the given conditions
			if msg.Type() == types.TypeMsgSubmitProposal {
				subMsg := msg.(*types.MsgSubmitProposal)
				return types.NewProposer(proposalID, subMsg.Proposer.String()), nil
			}
		}
	}

	return types.Proposer{}, fmt.Errorf("failed to find the proposer for proposalID %d", proposalID)
}

// QueryProposalByID takes a proposalID and returns a proposal
func QueryProposalByID(proposalID uint64, clientCtx client.Context, queryRoute string) ([]byte, error) {
	params := types.QueryProposalRequest{proposalID}
	bz, err := clientCtx.JSONMarshaler.MarshalJSON(&params)
	if err != nil {
		return nil, err
	}

	res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/proposal", queryRoute), bz)
	if err != nil {
		return nil, err
	}

	return res, err
}

// NormalizeVoteOption - normalize user specified vote option
func NormalizeVoteOption(option string) string {
	switch option {
	case "Yes", "yes":
		return types.OptionYes.String()

	case "Abstain", "abstain":
		return types.OptionAbstain.String()

	case "No", "no":
		return types.OptionNo.String()

	case "NoWithVeto", "no_with_veto":
		return types.OptionNoWithVeto.String()

	default:
		return option
	}
}

//NormalizeProposalType - normalize user specified proposal type
func NormalizeProposalType(proposalType string) string {
	switch proposalType {
	case "Text", "text":
		return types.ProposalTypeText

	default:
		return ""
	}
}

//NormalizeProposalStatus - normalize user specified proposal status
func NormalizeProposalStatus(status string) string {
	switch status {
	case "DepositPeriod", "deposit_period":
		return "DepositPeriod"
	case "VotingPeriod", "voting_period":
		return "VotingPeriod"
	case "Passed", "passed":
		return "Passed"
	case "Rejected", "rejected":
		return "Rejected"
	default:
		return status
	}
}
