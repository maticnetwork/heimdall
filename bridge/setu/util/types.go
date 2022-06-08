package util

type TendermintUnconfirmedTxs struct {
	Result struct {
		Total string   `json:"total"`
		Txs   []string `json:"txs"`
	} `json:"result"`
}
