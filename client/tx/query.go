// nolint
package tx

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/types"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

const (
	flagTags  = "tags"
	flagPage  = "page"
	flagLimit = "limit"
)

// ----------------------------------------------------------------------------
// CLI
// ----------------------------------------------------------------------------

// QueryTxsByEventsCmd returns a command to search through tagged transactions.
func QueryTxsByEventsCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "txs",
		Short: "Search for paginated transactions that match a set of tags",
		Long: strings.TrimSpace(`
Search for transactions that match the exact given tags where results are paginated.

Example:
$ gaiacli query txs --tags '<tag1>:<value1>&<tag2>:<value2>' --page 1 --limit 30
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			tagsStr := viper.GetString(flagTags)
			tagsStr = strings.Trim(tagsStr, "'")

			var tags []string
			if strings.Contains(tagsStr, "&") {
				tags = strings.Split(tagsStr, "&")
			} else {
				tags = append(tags, tagsStr)
			}

			var tmTags []string
			for _, tag := range tags {
				if !strings.Contains(tag, ":") {
					return fmt.Errorf("%s should be of the format <key>:<value>", tagsStr)
				} else if strings.Count(tag, ":") > 1 {
					return fmt.Errorf("%s should only contain one <key>:<value> pair", tagsStr)
				}

				keyValue := strings.Split(tag, ":")
				if keyValue[0] == types.TxHeightKey {
					tag = fmt.Sprintf("%s=%s", keyValue[0], keyValue[1])
				} else {
					tag = fmt.Sprintf("%s='%s'", keyValue[0], keyValue[1])
				}
				tmTags = append(tmTags, tag)
			}

			page := viper.GetInt(flagPage)
			limit := viper.GetInt(flagLimit)

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			txs, err := helper.QueryTxsByEvents(cliCtx, tmTags, page, limit)
			if err != nil {
				return err
			}

			var output []byte
			if cliCtx.Indent {
				output, err = cdc.MarshalJSONIndent(txs, "", "  ")
			} else {
				output, err = cdc.MarshalJSON(txs)
			}

			if err != nil {
				return err
			}

			fmt.Println(string(output))

			return nil
		},
	}

	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")

	if err := viper.BindPFlag(client.FlagNode, cmd.Flags().Lookup(client.FlagNode)); err != nil {
		logger.Error("QueryTxsByEventsCmd | BindPFlag | client.FlagNode", "Error", err)
	}

	cmd.Flags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")

	if err := viper.BindPFlag(client.FlagTrustNode, cmd.Flags().Lookup(client.FlagTrustNode)); err != nil {
		logger.Error("QueryTxsByEventsCmd | BindPFlag | client.FlagTrustNode", "Error", err)
	}

	cmd.Flags().String(flagTags, "", "Tag:value list of tags that must match")
	cmd.Flags().Uint32(flagPage, rest.DefaultPage, "Query a specific page of paginated results")
	cmd.Flags().Uint32(flagLimit, rest.DefaultLimit, "Query number of transactions results per page returned")

	if err := cmd.MarkFlagRequired(flagTags); err != nil {
		logger.Error("QueryTxsByEventsCmd | MarkFlagRequired | flagTags", "Error", err)
	}

	return cmd
}

// QueryTxCmd implements the default command for a tx query.
func QueryTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tx [hash]",
		Short: "Find a transaction by hash in a committed block.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			output, err := helper.QueryTx(cliCtx, args[0])
			if err != nil {
				return err
			}

			if output.Empty() {
				return fmt.Errorf("No transaction found with hash %s", args[0])
			}

			return cliCtx.PrintOutput(output)
		},
	}

	cmd.Flags().StringP(client.FlagNode, "n", "tcp://localhost:26657", "Node to connect to")

	if err := viper.BindPFlag(client.FlagNode, cmd.Flags().Lookup(client.FlagNode)); err != nil {
		logger.Error("QueryTxCmd | BindPFlag | client.FlagNode", "Error", err)
	}

	cmd.Flags().Bool(client.FlagTrustNode, false, "Trust connected full node (don't verify proofs for responses)")

	if err := viper.BindPFlag(client.FlagTrustNode, cmd.Flags().Lookup(client.FlagTrustNode)); err != nil {
		logger.Error("QueryTxCmd | BindPFlag | client.FlagTrustNode", "Error", err)
	}

	return cmd
}

// ----------------------------------------------------------------------------
// REST
// ----------------------------------------------------------------------------

//swagger:parameters txsGET
type txsGET struct {

	//in:query
	Height int64 `json:"height"`

	//in:query
	Page int64 `json:"page"`

	//in:query
	Limit int64 `json:"limit"`
}

// swagger:route GET /txs  txs txsGET
// It returns the list of transaction based on page,limit and events specified.
// QueryTxsRequestHandlerFn implements a REST handler that searches for transactions.
// Genesis transactions are returned if the height parameter is set to zero,
// otherwise the transactions are searched for by events.
func QueryTxsRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, sdk.AppendMsgToErr("could not parse query parameters", err.Error()))
			return
		}

		// if the height query param is set to zero, query for genesis transactions
		heightStr := r.FormValue("height")
		if heightStr != "" {
			if height, err := strconv.ParseInt(heightStr, 10, 64); err == nil && height == 0 {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		}

		var (
			events      []string
			txs         []sdk.TxResponse
			page, limit int
		)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		if len(r.Form) == 0 {
			rest.PostProcessResponse(w, cliCtx, txs)
			return
		}

		events, page, limit, err = rest.ParseHTTPArgs(r)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		searchResult, err := helper.QueryTxsByEvents(cliCtx, events, page, limit)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		rest.PostProcessResponseBare(w, cliCtx, searchResult)
	}
}

//swagger:parameters txsByHash
type txsByHash struct {

	//Hash
	//required:true
	//in:path
	Hash string `json:"hash"`

	//Height
	//in:query
	Height int64 `json:"height"`
}

// swagger:route GET /txs/{hash}  txs txsByHash
// It returns the transaction by hash.
// QueryTxRequestHandlerFn implements a REST handler that queries a transaction
// by hash in a committed block.
func QueryTxRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		hashHexStr := vars["hash"]

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		output, err := helper.QueryTx(cliCtx, hashHexStr)
		if err != nil {
			if strings.Contains(err.Error(), hashHexStr) {
				rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}

			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())

			return
		}

		if output.Empty() {
			rest.WriteErrorResponse(w, http.StatusNotFound, fmt.Sprintf("no transaction found with hash %s", hashHexStr))
		}

		rest.PostProcessResponse(w, cliCtx, output)
	}
}

//swagger:parameters txsHashCommitProof
type txsHashCommitProof struct {

	//in:path
	//required:true
	Hash string `json:"hash"`

	//in:query
	Height int64 `json:"height"`
}

// swagger:route GET /txs/{hash}/commit-proof  txs txsHashCommitProof
// It returns the commit-proof for the transaction.
// QueryCommitTxRequestHandlerFn implements a REST handler that queries vote, sigs and tx bytes committed block.
func QueryCommitTxRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		hash, err := hex.DecodeString(vars["hash"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		tx, err := helper.QueryTxWithProof(cliCtx, hash)
		if err != nil {
			if strings.Contains(err.Error(), vars["hash"]) {
				rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}

			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())

			return
		}

		// get block client
		blockDetails, err := helper.GetBlock(cliCtx, tx.Height+1)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// extract signs from votes
		sigs := helper.GetVoteSigs(blockDetails.Block.LastCommit.Precommits)

		// proof
		proofList := helper.GetMerkleProofList(&tx.Proof.Proof)
		proof := helper.AppendBytes(proofList...)

		// commit tx proof
		result := hmRest.CommitTxProof{
			Vote:  hex.EncodeToString(helper.GetVoteBytes(blockDetails.Block.LastCommit.Precommits, blockDetails.Block.ChainID)),
			Sigs:  hex.EncodeToString(sigs),
			Tx:    hex.EncodeToString(tx.Tx[authTypes.PulpHashLength:]),
			Proof: hex.EncodeToString(proof),
		}

		rest.PostProcessResponse(w, cliCtx, result)
	}
}

//swagger:parameters txsSideTx
type txsSideTx struct {

	//in:path
	//required:true
	Hash string `json:"hash"`

	//in:query
	Height int64 `json:"height"`
}

// swagger:route GET /txs/{hash}/side-tx  txs txsSideTx
// It returns the side-tx bytes
// QuerySideTxRequestHandlerFn implements a REST handler that queries sigs, side-tx bytes committed block
func QuerySideTxRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		hash, err := hex.DecodeString(vars["hash"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		tx, err := helper.QueryTxWithProof(cliCtx, hash)
		if err != nil {
			if strings.Contains(err.Error(), vars["hash"]) {
				rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
				return
			}

			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		// fetch side txs sigs
		decoder := helper.GetTxDecoder(authTypes.ModuleCdc)

		stdTx, err := decoder(tx.Tx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cmsg := stdTx.GetMsgs()[0] // get first message

		sideMsg, ok := cmsg.(hmTypes.SideTxMsg)
		if !ok {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Invalid side-tx msg")
			return
		}

		// side-tx data
		sideTxData := sideMsg.GetSideSignBytes()

		// get block details
		blockDetails, err := helper.GetBlock(cliCtx, tx.Height+2) // side-tx take 2 blocks to process
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "Side-tx is not processed yet.")
			return
		}

		// extract votes from response
		preCommits := blockDetails.Block.LastCommit.Precommits

		// extract side-tx signs from votes
		sigs, err := helper.GetSideTxSigs(tx.Tx.Hash(), sideTxData, preCommits)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, "Error while fetching sigs")
			return
		}

		// commit tx proof
		result := hmRest.SideTxProof{
			Sigs: formattedSigs(sigs),
			Tx:   hex.EncodeToString(tx.Tx),
			Data: hex.EncodeToString(sideTxData),
		}

		// cli ctx with height
		cliCtx.WithHeight(tx.Height + 2)

		rest.PostProcessResponse(w, cliCtx, result)
	}
}

func formattedSigs(sigs [][3]*big.Int) (result [][3]string) {
	for _, s := range sigs {
		result = append(result, [3]string{s[0].String(), s[1].String(), s[2].String()})
	}

	return result
}
