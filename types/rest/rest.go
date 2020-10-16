// Package rest provides HTTP types and primitives for REST
// requests validation and responses handling.
package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	cosmosContext "github.com/cosmos/cosmos-sdk/client"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/types"
)

const (
	DefaultPage    = 1
	DefaultLimit   = 30             // should be consistent with tendermint/tendermint/rpc/core/pipe.go:19
	TxMinHeightKey = "tx.minheight" // Inclusive minimum height filter
	TxMaxHeightKey = "tx.maxheight" // Inclusive maximum height filter
)

var (
	ParseFloat64OrReturnBadRequest     = rest.ParseFloat64OrReturnBadRequest
	ParseQueryHeightOrReturnBadRequest = rest.ParseQueryHeightOrReturnBadRequest
	ParseUint64OrReturnBadRequest      = rest.ParseUint64OrReturnBadRequest
	ParseInt64OrReturnBadRequest       = rest.ParseInt64OrReturnBadRequest
)

// ResponseWithHeight defines a response object type that wraps an original
// response with a height.
type ResponseWithHeight struct {
	Height int64           `json:"height"`
	Result json.RawMessage `json:"result"`
}

// NewResponseWithHeight creates a new ResponseWithHeight instance
func NewResponseWithHeight(height int64, result json.RawMessage) ResponseWithHeight {
	return ResponseWithHeight{
		Height: height,
		Result: result,
	}
}

// ParseResponseWithHeight returns the raw result from a JSON-encoded
// ResponseWithHeight object.
func ParseResponseWithHeight(cdc *codec.LegacyAmino, bz []byte) ([]byte, error) {
	r := ResponseWithHeight{}
	if err := cdc.UnmarshalJSON(bz, &r); err != nil {
		return nil, err
	}

	return r.Result, nil
}

// GasEstimateResponse defines a response definition for tx gas estimation.
type GasEstimateResponse struct {
	GasEstimate uint64 `json:"gas_estimate"`
}

// BaseReq defines a structure that can be embedded in other request structures
// that all share common "base" fields.
type BaseReq struct {
	From          string       `json:"from"`
	Memo          string       `json:"memo"`
	ChainID       string       `json:"chain_id"`
	AccountNumber uint64       `json:"account_number"`
	Sequence      uint64       `json:"sequence"`
	Fees          sdk.Coins    `json:"fees"`
	GasPrices     sdk.DecCoins `json:"gas_prices"`
	Gas           string       `json:"gas"`
	GasAdjustment string       `json:"gas_adjustment"`
	Simulate      bool         `json:"simulate"`
}

// NewBaseReq creates a new basic request instance and sanitizes its values
func NewBaseReq(
	from, memo, chainID string, gas, gasAdjustment string, accNumber, seq uint64,
	fees sdk.Coins, gasPrices sdk.DecCoins, simulate bool,
) BaseReq {

	return BaseReq{
		From:          strings.TrimSpace(from),
		Memo:          strings.TrimSpace(memo),
		ChainID:       strings.TrimSpace(chainID),
		Fees:          fees,
		GasPrices:     gasPrices,
		Gas:           strings.TrimSpace(gas),
		GasAdjustment: strings.TrimSpace(gasAdjustment),
		AccountNumber: accNumber,
		Sequence:      seq,
		Simulate:      simulate,
	}
}

// Sanitize performs basic sanitization on a BaseReq object.
func (br BaseReq) Sanitize() BaseReq {
	return NewBaseReq(
		br.From, br.Memo, br.ChainID, br.Gas, br.GasAdjustment,
		br.AccountNumber, br.Sequence, br.Fees, br.GasPrices, br.Simulate,
	)
}

// ValidateBasic performs basic validation of a BaseReq. If custom validation
// logic is needed, the implementing request handler should perform those
// checks manually.
func (br BaseReq) ValidateBasic(w http.ResponseWriter) bool {
	if !br.Simulate {
		switch {
		case len(br.ChainID) == 0:
			WriteErrorResponse(w, http.StatusUnauthorized, "chain-id required but not specified")
			return false

		case !br.Fees.IsZero() && !br.GasPrices.IsZero():
			// both fees and gas prices were provided
			WriteErrorResponse(w, http.StatusBadRequest, "cannot provide both fees and gas prices")
			return false

		case !br.Fees.IsValid() && !br.GasPrices.IsValid():
			// neither fees or gas prices were provided
			WriteErrorResponse(w, http.StatusPaymentRequired, "invalid fees or gas prices provided")
			return false
		}
	}

	if types.HexToHeimdallAddress(br.From).Empty() || len(br.From) == 0 {
		WriteErrorResponse(w, http.StatusUnauthorized, fmt.Sprintf("invalid from address: %s", br.From))
		return false
	}

	return true
}

// ReadRESTReq reads and unmarshals a Request's body to the the BaseReq stuct.
// Writes an error response to ResponseWriter and returns true if errors occurred.
func ReadRESTReq(w http.ResponseWriter, r *http.Request, cdc *codec.Codec, req interface{}) bool {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return false
	}

	err = cdc.UnmarshalJSON(body, req)
	if err != nil {
		WriteErrorResponse(w, http.StatusBadRequest, fmt.Sprintf("failed to decode JSON payload: %s", err))
		return false
	}

	return true
}

// ErrorResponse defines the attributes of a JSON error response.
type ErrorResponse struct {
	Code  int    `json:"code,omitempty"`
	Error string `json:"error"`
}

// NewErrorResponse creates a new ErrorResponse instance.
func NewErrorResponse(code int, err string) ErrorResponse {
	return ErrorResponse{Code: code, Error: err}
}

// CheckError takes care of writing an error response if err is not nil.
// Returns false when err is nil; it returns true otherwise.
func CheckError(w http.ResponseWriter, status int, err error) bool {
	if err != nil {
		WriteErrorResponse(w, status, err.Error())
		return true
	}

	return false
}

// CheckBadRequestError attaches an error message to an HTTP 400 BAD REQUEST response.
// Returns false when err is nil; it returns true otherwise.
func CheckBadRequestError(w http.ResponseWriter, err error) bool {
	return CheckError(w, http.StatusBadRequest, err)
}

// CheckInternalServerError attaches an error message to an HTTP 500 INTERNAL SERVER ERROR response.
// Returns false when err is nil; it returns true otherwise.
func CheckInternalServerError(w http.ResponseWriter, err error) bool {
	return CheckError(w, http.StatusInternalServerError, err)
}

// CheckNotFoundError attaches an error message to an HTTP 404 NOT FOUND response.
// Returns false when err is nil; it returns true otherwise.
func CheckNotFoundError(w http.ResponseWriter, err error) bool {
	return CheckError(w, http.StatusNotFound, err)
}

// WriteErrorResponse prepares and writes a HTTP error
// given a status code and an error message.
func WriteErrorResponse(w http.ResponseWriter, status int, err string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(codec.Cdc.MustMarshalJSON(NewErrorResponse(0, err)))
}

// WriteSimulationResponse prepares and writes an HTTP
// response for transactions simulations.
func WriteSimulationResponse(w http.ResponseWriter, cdc *codec.Codec, gas uint64) {
	gasEst := GasEstimateResponse{GasEstimate: gas}
	resp, err := cdc.MarshalJSON(gasEst)
	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

// ReturnNotFoundIfNoContent returns not found error if no content
func ReturnNotFoundIfNoContent(w http.ResponseWriter, data []byte, message string) bool {
	if len(data) == 0 {
		rest.WriteErrorResponse(w, http.StatusNotFound, errors.New(message).Error())
		return false
	}

	return true
}

// PostProcessResponseBare post processes a body similar to PostProcessResponse
// except it does not wrap the body and inject the height.
func PostProcessResponseBare(w http.ResponseWriter, ctx cosmosContext.Context, body interface{}) {
	var (
		resp []byte
		err  error
	)

	switch b := body.(type) {
	case []byte:
		resp = b

	default:
		resp, err = ctx.LegacyAmino.MarshalJSON(body)
		if CheckInternalServerError(w, err) {
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(resp)
}

// PostProcessResponse performs post processing for a REST response. The result
// returned to clients will contain two fields, the height at which the resource
// was queried at and the original result.
func PostProcessResponse(w http.ResponseWriter, cliCtx cosmosContext.Context, resp interface{}) {
	var result []byte

	if cliCtx.Height < 0 {
		WriteErrorResponse(w, http.StatusInternalServerError, fmt.Errorf("negative height in response").Error())
		return
	}

	switch resp.(type) {
	case []byte:
		result = resp.([]byte)

	default:
		var err error
		if cliCtx.Indent {
			result, err = cliCtx.Codec.MarshalJSONIndent(resp, "", "  ")
		} else {
			result, err = cliCtx.Codec.MarshalJSON(resp)
		}

		if err != nil {
			WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	wrappedResp := NewResponseWithHeight(cliCtx.Height, result)

	var (
		output []byte
		err    error
	)

	if cliCtx.Indent {
		output, err = cliCtx.Codec.MarshalJSONIndent(wrappedResp, "", "  ")
	} else {
		output, err = cliCtx.Codec.MarshalJSON(wrappedResp)
	}

	if err != nil {
		WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(output)
}

// ParseHTTPArgsWithLimit parses the request's URL and returns a slice containing
// all arguments pairs. It separates page and limit used for pagination where a
// default limit can be provided.
func ParseHTTPArgsWithLimit(r *http.Request, defaultLimit int) (tags []string, page, limit int, err error) {
	tags = make([]string, 0, len(r.Form))
	for key, values := range r.Form {
		if key == "page" || key == "limit" {
			continue
		}
		var value string
		value, err = url.QueryUnescape(values[0])
		if err != nil {
			return tags, page, limit, err
		}

		var tag string
		if key == tmTypes.TxHeightKey {
			tag = fmt.Sprintf("%s=%s", key, value)
		} else {
			tag = fmt.Sprintf("%s='%s'", key, value)
		}
		tags = append(tags, tag)
	}

	pageStr := r.FormValue("page")
	if pageStr == "" {
		page = DefaultPage
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			return tags, page, limit, err
		} else if page <= 0 {
			return tags, page, limit, errors.New("page must greater than 0")
		}
	}

	limitStr := r.FormValue("limit")
	if limitStr == "" {
		limit = defaultLimit
	} else {
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			return tags, page, limit, err
		} else if limit <= 0 {
			return tags, page, limit, errors.New("limit must greater than 0")
		}
	}

	return tags, page, limit, nil
}

// ParseHTTPArgs parses the request's URL and returns a slice containing all
// arguments pairs. It separates page and limit used for pagination.
func ParseHTTPArgs(r *http.Request) (tags []string, page, limit int, err error) {
	return ParseHTTPArgsWithLimit(r, DefaultLimit)
}

// ParseQueryParamBool parses the given param to a boolean. It returns false by
// default if the string is not parseable to bool.
func ParseQueryParamBool(r *http.Request, param string) bool {
	if value, err := strconv.ParseBool(r.FormValue(param)); err == nil {
		return value
	}

	return false
}

// GetRequest defines a wrapper around an HTTP GET request with a provided URL.
// An error is returned if the request or reading the body fails.
func GetRequest(url string) ([]byte, error) {
	res, err := http.Get(url) // nolint:gosec
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if err = res.Body.Close(); err != nil {
		return nil, err
	}

	return body, nil
}

// PostRequest defines a wrapper around an HTTP POST request with a provided URL and data.
// An error is returned if the request or reading the body fails.
func PostRequest(url string, contentType string, data []byte) ([]byte, error) {
	res, err := http.Post(url, contentType, bytes.NewBuffer(data)) // nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("error while sending post request: %w", err)
	}

	bz, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	if err = res.Body.Close(); err != nil {
		return nil, err
	}

	return bz, nil
}
