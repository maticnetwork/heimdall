//nolint
package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	restClient "github.com/maticnetwork/heimdall/client/rest"
	gcutils "github.com/maticnetwork/heimdall/gov/client/utils"
	"github.com/maticnetwork/heimdall/gov/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

//It represents the gov deposit parameters
//swagger:response govParametersDepositResponse
type govParametersDepositResponse struct {
	//in:body
	Output govParameterDepositStructure `json:"output"`
}

type govParameterDepositStructure struct {
	Height string              `json:"height"`
	Result govParameterDeposit `json:"result"`
}

type govParameterDeposit struct {
	MinDeposit       deposit `json:"min_deposit"`
	MaxDepositPeriod string  `json:"max_deposit_period"`
}

type deposit struct {
	Denom   string `json:"denom"`
	Account string `json:"amount"`
}

//It represents the gov tallying parameters
//swagger:response govParametersTallyingResponse
type govParametersTallyingResponse struct {
	//in:body
	Output govParameterTallyingStructure `json:"output"`
}

type govParameterTallyingStructure struct {
	Height string               `json:"height"`
	Result govParameterTallying `json:"result"`
}

type govParameterTallying struct {
	Quorum    string `json:"quorum"`
	Threshold string `json:"threshold"`
	Veto      string `json:"veto"`
}

//It represents the gov voting parameters
//swagger:response govParametersVotingResponse
type govParametersVotingResponse struct {
	//in:body
	Output govParameterVotingStructure `json:"output"`
}

type govParameterVotingStructure struct {
	Height string             `json:"height"`
	Result govParameterVoting `json:"result"`
}

type govParameterVoting struct {
	VotingPeriod string `json:"voting_period"`
}

//It represents the gov proposals
//swagger:response govProposalsResponse
type govProposalsResponse struct {
	//in:body
	Output govProposalsStructure `json:"output"`
}

type govProposalsStructure struct {
	Height string     `json:"height"`
	Result []proposal `json:"result"`
}

type proposal struct {
	Content          content     `json:"content"`
	Id               string      `json:"id"`
	ProposalStatus   string      `json:"proposal_status"`
	FinalTallyResult TallyResult `json:"final_tally_result"`
	SubmitTime       string      `json:"submit_time"`
	DepositEndTime   string      `json:"deposit_end_time"`
	TotalDeposit     deposit     `json:"total_deposit"`
	VotingStartTime  string      `json:"voting_start_time"`
	VotingEndTime    string      `json:"voting_end_time"`
}

type content struct {
	Type  string `json:"type"`
	Value value  `json:"value"`
}

type value struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Changes     []change `json:"change"`
}

type change struct {
	Subspace string `json:"subspace"`
	Key      string `json:"key"`
	Value    string `json:"value"`
}

type TallyResult struct {
	Yes        string `json:"yes"`
	Abstain    string `json:"abstain"`
	No         string `json:"no"`
	NoWithVeto string `json:"no_with_veto"`
}

//It represents the gov proposal
//swagger:response govProposalResponse
type govProposalResponse struct {
	//in:body
	Output govProposalsStructure `json:"output"`
}

type govProposalStructure struct {
	Height string   `json:"height"`
	Result proposal `json:"result"`
}

//It represents the gov proposer
//swagger:response govProposerResponse
type govProposerResponse struct {
	//in:body
	Output govProposerStructure `json:"output"`
}

type govProposerStructure struct {
	Height string   `json:"height"`
	Result proposer `json:"result"`
}

type proposer struct {
	ProposalId string `json:"proposal_id"`
	Proposer   string `json:"proposer"`
}

//It represents the gov Tally based on Id
//swagger:response govTallyResponse
type govTallyResponse struct {
	//in:body
	Output govProposalsStructure `json:"output"`
}

type govTallyStructure struct {
	Height string      `json:"height"`
	Result TallyResult `json:"result"`
}

//It represents the votes responses
//swagger:response govVotesResponse
type govVotesResponse struct {
	//in:body
	Output govVotesStructure `json:"output"`
}

type govVotesStructure struct {
	Height string `json:"height"`
	Result []vote `json:"result"`
}

type vote struct {
	ProposalId string `json:"proposal_id"`
	Voter      string `json:"voter"`
	Option     string `json:"option"`
}

//It represents the vote response
//swagger:response govVoteResponse
type govVoteResponse struct {
	//in:body
	Output govVoteStructure `json:"output"`
}

type govVoteStructure struct {
	Height string `json:"height"`
	Result vote   `json:"result"`
}

//It represents the vote response
//swagger:response govDepositResponse
type govDepositResponse struct {
	//in:body
	Output govDepositStructure `json:"output"`
}

type govDepositStructure struct {
	Height string    `json:"height"`
	Result deposited `json:"result"`
}

type deposited struct {
	ProposalId int64     `json:"proposal_id"`
	Depositor  string    `json:"depositor"`
	Amount     []deposit `json:"amount"`
}

//It represents the vote response
//swagger:response govDepositsResponse
type govDepositsResponse struct {
	//in:body
	Output govDepositsStructure `json:"output"`
}

type govDepositsStructure struct {
	Height string      `json:"height"`
	Result []deposited `json:"result"`
}

// REST Variable names
// nolint
const (
	RestParamsType     = "type"
	RestProposalID     = "proposal-id"
	RestDepositor      = "depositor"
	RestVoter          = "voter"
	RestProposalStatus = "status"
	RestNumLimit       = "limit"
)

// ProposalRESTHandler defines a REST handler implemented in another module. The
// sub-route is mounted on the governance REST handler.
type ProposalRESTHandler struct {
	SubRoute string
	Handler  func(http.ResponseWriter, *http.Request)
}

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, phs []ProposalRESTHandler) {
	propSubRtr := r.PathPrefix("/gov/proposals").Subrouter()
	for _, ph := range phs {
		propSubRtr.HandleFunc(fmt.Sprintf("/%s", ph.SubRoute), ph.Handler).Methods("POST")
	}

	r.HandleFunc("/gov/proposals", postProposalHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/deposits", RestProposalID), depositHandlerFn(cliCtx)).Methods("POST")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/votes", RestProposalID), voteHandlerFn(cliCtx)).Methods("POST")

	r.HandleFunc(
		fmt.Sprintf("/gov/parameters/{%s}", RestParamsType),
		queryParamsHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc("/gov/proposals", queryProposalsWithParameterFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}", RestProposalID), queryProposalHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(
		fmt.Sprintf("/gov/proposals/{%s}/proposer", RestProposalID),
		queryProposerHandlerFn(cliCtx),
	).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/deposits", RestProposalID), queryDepositsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/deposits/{%s}", RestProposalID, RestDepositor), queryDepositHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/tally", RestProposalID), queryTallyOnProposalHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/votes", RestProposalID), queryVotesOnProposalHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/gov/proposals/{%s}/votes/{%s}", RestProposalID, RestVoter), queryVoteHandlerFn(cliCtx)).Methods("GET")
}

// PostProposalReq defines the properties of a proposal request's body.
type PostProposalReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title          string                  `json:"title" yaml:"title"`                     // Title of the proposal
	Description    string                  `json:"description" yaml:"description"`         // Description of the proposal
	ProposalType   string                  `json:"proposal_type" yaml:"proposal_type"`     // Type of proposal. Initial set {PlainTextProposal, SoftwareUpgradeProposal}
	Proposer       hmTypes.HeimdallAddress `json:"proposer" yaml:"proposer"`               // Address of the proposer
	Validator      hmTypes.ValidatorID     `json:"validator" yaml:"validator"`             // id of the validator
	InitialDeposit sdk.Coins               `json:"initial_deposit" yaml:"initial_deposit"` // Coins to add to the proposal's deposit
}

// DepositReq defines the properties of a deposit request's body.
type DepositReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Depositor hmTypes.HeimdallAddress `json:"depositor" yaml:"depositor"` // Address of the depositor
	Amount    sdk.Coins               `json:"amount" yaml:"amount"`       // Coins to add to the proposal's deposit
	Validator hmTypes.ValidatorID     `json:"validator" yaml:"validator"` // id of the validator
}

// VoteReq defines the properties of a vote request's body.
type VoteReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Voter     hmTypes.HeimdallAddress `json:"voter" yaml:"voter"`         // address of the voter
	Option    string                  `json:"option" yaml:"option"`       // option from OptionSet chosen by the voter
	Validator hmTypes.ValidatorID     `json:"validator" yaml:"validator"` // id of the validator
}

//swagger:parameters govProposals
type govProposalsParam struct {

	//Body
	//required:true
	//in:body
	Input govProposalsInput `json:"input"`
}

type govProposalsInput struct {
	BaseReq        BaseReq `json:"base_req"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	ProposalType   string  `json:"proposal_type"`
	Proposer       string  `json:"proposer"`
	Validator      string  `json:"validator"`
	InitialDeposit []coin  `json:"initial_deposit"`
}

type coin struct {
	Denom  string `json:"denom"`
	Amount string `json:"amount"`
}

type BaseReq struct {

	//Address of the sender
	//required:true
	//in:body
	From string `json:"address"`

	//Chain ID of Heimdall
	//required:true
	//in:body
	ChainID string `json:"chain_id"`
}

// swagger:route POST /gov/proposals gov govProposals
// It returns the prepared msg for gov Proposals
// responses:
//   200: interface{}
func postProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PostProposalReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		proposalType := gcutils.NormalizeProposalType(req.ProposalType)
		content := types.ContentFromProposalType(req.Title, req.Description, proposalType)

		msg := types.NewMsgSubmitProposal(content, req.InitialDeposit, req.Proposer, req.Validator)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

//swagger:parameters govProposalsDeposits
type govProposalsDeposits struct {

	//Proposal ID
	//required:true
	//in:path
	ProposalID int64 `json:"proposal-id"`

	//Body
	//required:true
	//in:body
	Input govProposalsDepositsInput `json:"input"`
}

type govProposalsDepositsInput struct {
	BaseReq   BaseReq `json:"base_req"`
	Depositor string  `json:"depositor"`
	Amount    []coin  `json:"amount"`
	Validator int64   `json:"validator"`
}

// swagger:route POST /gov/proposals/{proposal-id}/deposits gov govProposalsDeposits
// It returns the prepared msg for gov proposal deposit
// responses:
//   200: govProposalsDepositsResponse

func depositHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		var req DepositReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		// create the message
		msg := types.NewMsgDeposit(req.Depositor, proposalID, req.Amount, req.Validator)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

//swagger:parameters govProposalsVotes
type govProposalsVotesParam struct {

	//Proposal Id
	//required:true
	//in:path
	ProposalId int64 `json:"proposal-id"`

	//Body
	//required:true
	//in:body
	Input govProposalsVotesInput `json:"input"`
}

type govProposalsVotesInput struct {
	BaseReq   BaseReq `json:"base_req"`
	Voter     string  `json:"voter"`
	Option    string  `json:"option"`
	Validator int64   `json:"validator"`
}

// swagger:route POST /gov/proposals/{proposal-id}/votes gov govProposalsVotes
// It returns the prepared msg for gov proposal votes
// responses:
//   200: interface{}
func voteHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		var req VoteReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		voteOption, err := types.VoteOptionFromString(gcutils.NormalizeVoteOption(req.Option))
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgVote(req.Voter, proposalID, voteOption, req.Validator)
		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		restClient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}

// swagger:route GET /gov/parameters/voting gov govParametersVoting
// It returns the gov voting parameters
// responses:
//   200: govParametersVotingResponse

// swagger:route GET /gov/parameters/tallying gov govParametersTallying
// It returns the gov tallying parameters
// responses:
//   200: govParametersTallyingResponse

// swagger:route GET /gov/parameters/deposit gov govParametersDeposit
// It returns the gov deposit parameters
// responses:
//   200: govParametersDepositResponse

func queryParamsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		paramType := vars[RestParamsType]

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/gov/%s/%s", types.QueryParams, paramType), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /gov/proposals/{proposal-id} gov govProposalById
// It returns the proposal based on the proposal Id
// responses:
//   200: govProposalResponse
func queryProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		params := types.NewQueryProposalParams(proposalID)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/gov/proposal", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /gov/proposals/{proposal-id}/deposits gov govDepositByProposalId
// It returns the gov deposit parameters
// responses:
//   200: govDepositsResponse
func queryDepositsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		params := types.NewQueryProposalParams(proposalID)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData("custom/gov/proposal", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var proposal types.Proposal
		if err := cliCtx.Codec.UnmarshalJSON(res, &proposal); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// For inactive proposals we must query the txs directly to get the deposits
		// as they're no longer in state.
		propStatus := proposal.Status
		if !(propStatus == types.StatusVotingPeriod || propStatus == types.StatusDepositPeriod) {
			res, err = gcutils.QueryDepositsByTxQuery(cliCtx, params)
		} else {
			res, _, err = cliCtx.QueryWithData("custom/gov/deposits", bz)
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters govProposerByProposalId govProposalVotesByProposalId govProposalById govTallyByProposalId govDepositByProposalId
type ProposalId struct {
	//Proposal ID
	//in:path
	ProposalId int64 `json:"proposal-id"`
}

// swagger:route GET /gov/proposals/{proposal-id}/proposer gov govProposerByProposalId
// It returns the proposer based on the proposal.
// responses:
//   200: govProposerResponse
func queryProposerHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		res, err := gcutils.QueryProposerByTxQuery(cliCtx, proposalID)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters govDepositBasedOnDepositor
type DepositQuery struct {
	//Proposal ID
	//in:path
	ProposalId int64 `json:"proposal-id"`

	//Depositor ID
	//in:path
	Depositor int64 `json:"depositor"`
}

// swagger:route GET /gov/proposals/{proposal-id}/deposits/{depositor} gov govDepositBasedOnDepositor
// It returns the deposit for a particular proposal based on depositor Id
// responses:
//   200: govDepositResponse
func queryDepositHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]
		strDepositorID := vars[RestDepositor]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(strDepositorID) == 0 {
			err := errors.New("depositorId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		depositorID, ok := rest.ParseUint64OrReturnBadRequest(w, strDepositorID)
		if !ok {
			return
		}

		params := types.NewQueryDepositParams(proposalID, hmTypes.NewValidatorID(depositorID))

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData("custom/gov/deposit", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var deposit types.Deposit
		if err := cliCtx.Codec.UnmarshalJSON(res, &deposit); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// For an empty deposit, either the proposal does not exist or is inactive in
		// which case the deposit would be removed from state and should be queried
		// for directly via a txs query.
		if deposit.Empty() {
			bz, err := cliCtx.Codec.MarshalJSON(types.NewQueryProposalParams(proposalID))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			res, _, err = cliCtx.QueryWithData("custom/gov/proposal", bz)
			if err != nil || len(res) == 0 {
				err := fmt.Errorf("proposalID %d does not exist", proposalID)
				rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}

			res, err = gcutils.QueryDepositByTxQuery(cliCtx, params)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters govVotesBasedOnVoterId
type VotesQuery struct {
	//Proposal ID
	//in:path
	ProposalId int64 `json:"proposal-id"`

	//Voter ID
	//in:path
	Voter int64 `json:"voter"`
}

// swagger:route GET /gov/proposals/{proposal-id}/votes/{voter} gov govVotesBasedOnVoterId
// It returns the votes on the specific proposal Id based on voter Id.
// responses:
//   200: govVoteResponse
func queryVoteHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]
		strVoterID := vars[RestVoter]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(strVoterID) == 0 {
			err := errors.New("voterId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		voterID, ok := rest.ParseUint64OrReturnBadRequest(w, strVoterID)
		if !ok {
			return
		}

		params := types.NewQueryVoteParams(proposalID, hmTypes.NewValidatorID(voterID))

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData("custom/gov/vote", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var vote types.Vote
		if err := cliCtx.Codec.UnmarshalJSON(res, &vote); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// For an empty vote, either the proposal does not exist or is inactive in
		// which case the vote would be removed from state and should be queried for
		// directly via a txs query.
		if vote.Empty() {
			bz, err := cliCtx.Codec.MarshalJSON(types.NewQueryProposalParams(proposalID))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			res, _, err = cliCtx.QueryWithData("custom/gov/proposal", bz)
			if err != nil || len(res) == 0 {
				err := fmt.Errorf("proposalID %d does not exist", proposalID)
				rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}

			res, err = gcutils.QueryVoteByTxQuery(cliCtx, params)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /gov/proposals/{proposal-id}/votes gov govProposalVotesByProposalId
// It returns the proposal votes based on proposal id
// responses:
//   200: govVotesResponse
// todo: Split this functionality into helper functions to remove the above
func queryVotesOnProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		params := types.NewQueryProposalParams(proposalID)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, _, err := cliCtx.QueryWithData("custom/gov/proposal", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		var proposal types.Proposal
		if err := cliCtx.Codec.UnmarshalJSON(res, &proposal); err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// For inactive proposals we must query the txs directly to get the votes
		// as they're no longer in state.
		propStatus := proposal.Status
		if !(propStatus == types.StatusVotingPeriod || propStatus == types.StatusDepositPeriod) {
			res, err = gcutils.QueryVotesByTxQuery(cliCtx, params)
		} else {
			res, _, err = cliCtx.QueryWithData("custom/gov/votes", bz)
		}

		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters govProposalsGet
type ProposalsQuery struct {

	//Proposal Status [DepositPeriod,Passed,Rejected,Failed,VotingPeriod]
	//in:query
	Status string `json:"status"`

	//Limit
	//in:query
	Limit int64 `json:"limit"`

	//Voter ID
	//in:query
	Voter int64 `json:"voter"`

	//Depositor ID
	//in:query
	Depositor int64 `json:"depositor"`
}

// swagger:route GET /gov/proposals gov govProposalsGet
// It returns the gov proposals
// responses:
//   200: govProposalsResponse
// todo: Split this functionality into helper functions to remove the above
func queryProposalsWithParameterFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		strVoterID := r.URL.Query().Get(RestVoter)
		strDepositorID := r.URL.Query().Get(RestDepositor)
		strProposalStatus := r.URL.Query().Get(RestProposalStatus)
		strNumLimit := r.URL.Query().Get(RestNumLimit)

		params := types.QueryProposalsParams{}

		if len(strVoterID) != 0 {
			voterID, ok := rest.ParseUint64OrReturnBadRequest(w, strVoterID)
			if ok {
				params.Voter = hmTypes.NewValidatorID(voterID)
			}
		}

		if len(strDepositorID) != 0 {
			depositorID, ok := rest.ParseUint64OrReturnBadRequest(w, strDepositorID)
			if ok {
				params.Depositor = hmTypes.NewValidatorID(depositorID)
			}
		}

		if len(strProposalStatus) != 0 {
			proposalStatus, err := types.ProposalStatusFromString(gcutils.NormalizeProposalStatus(strProposalStatus))
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
			params.ProposalStatus = proposalStatus
		}

		if len(strNumLimit) != 0 {
			numLimit, ok := rest.ParseUint64OrReturnBadRequest(w, strNumLimit)
			if !ok {
				return
			}
			params.Limit = numLimit
		}

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/gov/proposals", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /gov/proposals/{proposal-id}/tally gov govTallyByProposalId
// It returns the tally on the proposal ID
// responses:
//   200: govTallyResponse
func queryTallyOnProposalHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)
		strProposalID := vars[RestProposalID]

		if len(strProposalID) == 0 {
			err := errors.New("proposalId required but not specified")
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		proposalID, ok := rest.ParseUint64OrReturnBadRequest(w, strProposalID)
		if !ok {
			return
		}

		params := types.NewQueryProposalParams(proposalID)

		bz, err := cliCtx.Codec.MarshalJSON(params)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, height, err := cliCtx.QueryWithData("custom/gov/tally", bz)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters govParametersDeposit govProposals
type Height struct {

	//Block Height
	//in:query
	Height string `json:"height"`
}
