package types

// Parameter store key
var (
	ParamStoreKeyDepositParams = []byte("depositparams")
	ParamStoreKeyVotingParams  = []byte("votingparams")
	ParamStoreKeyTallyParams   = []byte("tallyparams")
)

// // Key declaration for parameters
// func ParamKeyTable() subspace.KeyTable {
// 	return subspace.NewKeyTable(
// 		ParamStoreKeyDepositParams, DepositParams{},
// 		ParamStoreKeyVotingParams, VotingParams{},
// 		ParamStoreKeyTallyParams, TallyParams{},
// 	)
// }

// // NewDepositParams creates a new DepositParams object
// func NewDepositParams(minDeposit Coins, maxDepositPeriod time.Duration) DepositParams {
// 	return DepositParams{
// 		MinDeposit:       minDeposit,
// 		MaxDepositPeriod: maxDepositPeriod,
// 	}
// }

// // Checks equality of DepositParams
// func (dp DepositParams) Equal(dp2 DepositParams) bool {
// 	return dp.MinDeposit.IsEqual(dp2.MinDeposit) && dp.MaxDepositPeriod == dp2.MaxDepositPeriod
// }

// // NewTallyParams creates a new TallyParams object
// func NewTallyParams(quorum, threshold, veto Dec) TallyParams {
// 	return TallyParams{
// 		Quorum:    quorum,
// 		Threshold: threshold,
// 		Veto:      veto,
// 	}
// }

// // NewVotingParams creates a new VotingParams object
// func NewVotingParams(votingPeriod time.Duration) VotingParams {
// 	return VotingParams{
// 		VotingPeriod: votingPeriod,
// 	}
// }

// Params returns all of the governance params
type Params struct {
	VotingParams  VotingParams  `json:"voting_params" yaml:"voting_params"`
	TallyParams   TallyParams   `json:"tally_params" yaml:"tally_params"`
	DepositParams DepositParams `json:"deposit_params" yaml:"deposit_parmas"`
}

// func (gp Params) String() string {
// 	return gp.VotingParams.String() + "\n" +
// 		gp.TallyParams.String() + "\n" + gp.DepositParams.String()
// }

func NewParams(vp VotingParams, tp TallyParams, dp DepositParams) Params {
	return Params{
		VotingParams:  vp,
		DepositParams: dp,
		TallyParams:   tp,
	}
}
