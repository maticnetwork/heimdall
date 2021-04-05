package rest

// CommitTxProof commit tx proof
type CommitTxProof struct {
	Vote  string `json:"vote"`
	Sigs  string `json:"sigs"`
	Tx    string `json:"tx"`
	Proof string `json:"proof"`
}

// SideTxProof side-tx proof
type SideTxProof struct {
	Sigs [][3]string `json:"sigs"`
	Tx   string      `json:"tx"`
	Data string      `json:"data"`
}
