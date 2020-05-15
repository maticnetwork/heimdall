package types

import (
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// DONTCOVER

// Query endpoints supported by the slashing querier
const (
	QueryParameters        = "parameters"
	QuerySigningInfo       = "signingInfo"
	QuerySigningInfos      = "signingInfos"
	QuerySlashingInfo      = "slashingInfo"
	QuerySlashingInfos     = "slashingInfos"
	QuerySlashingInfoBytes = "slashingInfoBytes"
	QueryTickSlashingInfos = "tickSlashingInfos"
	QuerySlashingSequence  = "slashing-sequence"
	QueryTickCount         = "tick-count"
)

// QuerySigningInfoParams defines the params for the following queries:
// - 'custom/slashing/signingInfo'
type QuerySigningInfoParams struct {
	ValidatorID hmTypes.ValidatorID
}

// NewQuerySigningInfoParams creates a new QuerySigningInfoParams instance
func NewQuerySigningInfoParams(valID hmTypes.ValidatorID) QuerySigningInfoParams {
	return QuerySigningInfoParams{valID}
}

// QuerySigningInfosParams defines the params for the following queries:
// - 'custom/slashing/signingInfos'
type QuerySigningInfosParams struct {
	Page, Limit int
}

// NewQuerySigningInfosParams creates a new QuerySigningInfosParams instance
func NewQuerySigningInfosParams(page, limit int) QuerySigningInfosParams {
	return QuerySigningInfosParams{page, limit}
}

// QuerySlashingInfoParams defines the params for the following queries:
// - 'custom/slashing/slashingInfo'
type QuerySlashingInfoParams struct {
	ValidatorID hmTypes.ValidatorID
}

// NewQuerySlashingInfoParams creates a new QuerySlashingInfoParams instance
func NewQuerySlashingInfoParams(valID hmTypes.ValidatorID) QuerySlashingInfoParams {
	return QuerySlashingInfoParams{valID}
}

// QuerySlashingInfosParams defines the params for the following queries:
// - 'custom/slashing/slashingInfos'
type QuerySlashingInfosParams struct {
	Page, Limit int
}

// NewQuerySlashingInfosParams creates a new QuerySlashingInfosParams instance
func NewQuerySlashingInfosParams(page, limit int) QuerySlashingInfosParams {
	return QuerySlashingInfosParams{page, limit}
}

// QueryTickSlashingInfosParams defines the params for the following queries:
// - 'custom/slashing/tick_slash_infos'
type QueryTickSlashingInfosParams struct {
	Page, Limit int
}

// NewQueryTickSlashingInfosParams creates a new QueryTickSlashingInfosParams instance
func NewQueryTickSlashingInfosParams(page, limit int) QueryTickSlashingInfosParams {
	return QueryTickSlashingInfosParams{page, limit}
}

// QuerySlashingSequenceParams defines the params for querying an account Sequence.
type QuerySlashingSequenceParams struct {
	TxHash   string
	LogIndex uint64
}

// NewQuerySlashingSequenceParams creates a new instance of QuerySlashingSequenceParams.
func NewQuerySlashingSequenceParams(txHash string, logIndex uint64) QuerySlashingSequenceParams {
	return QuerySlashingSequenceParams{TxHash: txHash, LogIndex: logIndex}
}
